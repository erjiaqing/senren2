<template>
  <div>
    <el-table
      :data="rankList.rank"
      style="width: 100%"
    >
      <el-table-column
        :label="v.text"
        v-for="v in rankList.title"
        :key="v.field"
        fixed
      >
        <template slot-scope="scope">
          {{ scope.row[v.field] }}
        </template>
      </el-table-column>
      <el-table-column fixed>
        <template slot-scope="scope">
          {{ scope.row.user.name }}
        </template>
      </el-table-column>
      <el-table-column
        v-for="(v, id) in rankList.problem"
        :key="id"
      >
        <template slot-scope="scope">
          <i class="submission-ac el-icon-circle-check-outline" v-if="scope.row.result[id].state == 'AC'"> {{ scope.row.result[id].time }}</i>
          <i class="submission-ac el-icon-circle-check" v-if="scope.row.result[id].state == 'AC_FIRST'"> {{ scope.row.result[id].time }}</i>
          <i class="submission-wa el-icon-circle-close-outline" v-if="scope.row.result[id].state == 'NO'"> {{ scope.row.result[id].time }}</i>
          <i class="submission-pending el-icon-remove-outline" v-if="scope.row.result[id].state == 'PENDING'"> {{ scope.row.result[id].time }}</i>
        </template>
      </el-table-column>
    </el-table>
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
    clearFilter() {
      this.$router.push(
        `/${this.$route.params.domain}/contest/${
          this.$route.params.uid
        }/submissions`
      );
    }
  },
  computed: mapState(["user", "contest"]),
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
  color: #67C23A;
}

.submission-wa {
  color: #F56C6C;
}

.submission-pending {
  color: #E6A23C;
}
</style>


