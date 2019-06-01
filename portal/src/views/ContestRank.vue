<template>
  <div>
    <div>
      <div v-if="contestTimer.state == 'PENDING'">
        <span>未开始</span>
        <div class="timerText" style="padding: 3px 0">{{ contestTimer.remain }}</div>
      </div>
      <div v-else-if="contestTimer.state == 'RUNNING'">
        <span>进行中</span>
        <el-progress :show-text="false" :percentage="contestTimer.ratio * 100"></el-progress>
        <div class="timerText" style="padding: 3px 0">{{ contestTimer.remain }}</div>
      </div>
      <div v-else-if="contestTimer.state == 'FROZEN'">
        <span>进行中 - 封榜</span>
        <el-progress :show-text="false" color="#E6A23C" :percentage="contestTimer.ratio * 100"></el-progress>
        <div class="timerText" style="padding: 3px 0">{{ contestTimer.remain }}</div>
      </div>
      <div v-else-if="contestTimer.state == 'END'">
        <span>已结束</span>
      </div>
    </div>

    <el-table
      :data="rankList.rank"
      :row-class-name="tableIsSelf"
      :summary-method="getSummaries"
      @row-click="gotoSubmissionList"
      show-summary
      size="mini"
      style="width: 100%"
    >
      <el-table-column :label="v.text" v-for="v in rankList.title" :key="v.field">
        <template slot-scope="scope">{{ scope.row[v.field] }}</template>
      </el-table-column>
      <el-table-column fixed>
        <template slot-scope="scope">{{ scope.row.user.name }}</template>
      </el-table-column>
      <el-table-column
        v-for="(v, id) in rankList.problem"
        :key="id"
        :label="contest.problem_list[id].$charid"
        align="center"
      >
        <template slot-scope="scope">
          <i
            class="submission-ac el-icon-circle-check"
            v-if="scope.row.result[id].state == 'AC'"
          >{{ scope.row.result[id].time }}</i>
          <i
            class="submission-ac el-icon-success"
            v-if="scope.row.result[id].state == 'AC_FIRST'"
          >{{ scope.row.result[id].time }}</i>
          <i
            class="submission-wa el-icon-remove-outline"
            v-if="scope.row.result[id].state == 'NO'"
          >{{ scope.row.result[id].time }}</i>
          <i
            class="submission-pending el-icon-time"
            v-if="scope.row.result[id].state == 'PENDING'"
          >{{ scope.row.result[id].time }}</i>
        </template>
      </el-table-column>
    </el-table>
    <el-alert title="单击表格某行可以查看用户提交情况。" type="info" :closable="false"></el-alert>
  </div>
</template>

<script>
import { RPC } from "../rpc.js";
import { ConstString } from "../consts.js";
import { Homework, ICPC, ICPCInd } from "../rank.js";
import { mapState } from "vuex";

export default {
  data() {
    return {
      submissions: null,
      loading: false,
      error: false,
      codeEditor: false,
      code: "",
      tags: ConstString,
      filterApplied: false,
      filter: [],
      loadingInterval: -1,
      rankList: {}
    };
  },
  methods: {
    loadSubmission: async function() {
      this.loading = true;
      this.filterApplied = false;
      let res = await RPC.doRPC("getContestSubmissions", {
        domain: this.$route.params.domain,
        filter: this.$route.params.uid + "||rank"
      });
      this.loading = false;
      if (res == null) {
        this.error = true;
        return;
      }
      this.submissions = res.submissions;
      this.calculateRank();
    },
    calculateRank() {
      let users = [];
      let userMap = {};

      this.submissions.forEach(element => {
        if (!userMap[element.user_uid]) {
          userMap[element.user_uid] = {
            name: element.user_name
          };
          users.push({
            uid: element.user_uid,
            name: element.user_name
          });
        }
      });

      this.rankList = ICPC.calc(
        users,
        this.contest.problem_list,
        this.submissions,
        this.contest
      );
    },
    selectSubmission(e) {
      if (e.uid != "") {
        this.$router.push(
          `/${this.$route.params.domain}/contest/${e.contest_uid}/submission/${
            e.uid
          }`
        );
      }
    },
    gotoSubmissionList(e) {
      if (e.uid != "") {
        this.$router.push(
          `/${this.$route.params.domain}/contest/${
            this.contest.uid
          }/submissions/;${e.user.uid}`
        );
      }
    },
    tableIsSelf({ row, rowIndex }) {
      if (row.user.uid === this.user.uid) {
        return "selfRankItem";
      }
      return "";
    },
    getSummaries(param) {
      const { columns, data } = param;
      const sums = [];
      if (!this.rankList || !this.rankList.title) {
        return sums;
      }
      columns.forEach((column, index) => {
        if (index < 1 + this.rankList.title.length) {
          sums[index] = "";
          return;
        }
        let vals = this.rankList.problem[
          index - 1 - this.rankList.title.length
        ];
        sums[index] = `${vals.ac} / ${vals.att} / ${
          vals.ac ? vals.firstblood : "---"
        }`;
      });

      return sums;
    },
    clearFilter() {
      this.$router.push(
        `/${this.$route.params.domain}/contest/${
          this.$route.params.uid
        }/submissions`
      );
    }
  },
  computed: mapState(["user", "contest", "contestTimer"]),
  watch: {
    $route: function() {
      this.loadSubmission();
    }
  },
  created() {
    this.loadSubmission();
    this.loadingInterval = setInterval(this.loadSubmission, 30000);
  },
  beforeDestroy() {
    clearInterval(this.loadingInterval);
  }
};
</script>

<style>
.submission-title {
  font-size: 28pt;
  text-align: left;
  padding-bottom: 12pt;
}

.submission-content {
  text-align: left;
  font-size: 14pt;
}

.codeEditorPanel {
  -moz-box-shadow: 0px 0px 5px #999;
  -webkit-box-shadow: 0px 0px 5px #999;
  box-shadow: 0px 0px 5px #999;
  border-radius: 5px;

  transition: all 0.3s linear;
}

.codeEditorPanel:hover {
  -moz-box-shadow: 0px 0px 15px #999;
  -webkit-box-shadow: 0px 0px 15px #999;
  box-shadow: 0px 0px 15px #999;
  border-radius: 5px;
}

#baseinfo_container .el-form-item {
  margin-right: 0;
  margin-bottom: 0;
  margin-left: 0;
  padding: 0;
  width: 50%;
}

#baseinfo_container label {
  width: 72pt;
  color: #99a9bf;
}

.submission_metainfo_container {
  padding: 0;
}

.submission-ac {
  color: #67c23a;
}

.submission-wa {
  color: #f56c6c;
}

.el-table .selfRankItem {
  background: #ecf5ff;
}

.submission-pending {
  color: #e6a23c;
}
</style>


