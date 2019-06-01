<template>
  <el-row :gutter="20">
    <el-col :span="24" :md="16" :lg="18">
    <el-alert
      title="显示的是经过筛选的提交列表"
      v-if="filterApplied"
      close-text="清除筛选"
      @close="clearFilter"
      type="info"
    >
    </el-alert>

    <el-alert
      title="比赛期间，其他人的提交状态只会为 Accepted, Wrong Answer, Compile Error, Pending 之一"
      type="info"
      :closable="false"
    >
    </el-alert>
    <el-table
      :data="submissions"
      style="width: 100%"
      :row-class-name="tableIsSelf"
      @row-click="selectSubmission"
    >
      <el-table-column
        label="评测结果"
        width="200pt"
      >
        <template slot-scope="scope">
          <span v-if="tags['JUDGE_TAG_' + scope.row.verdict]">
            <el-tag :type="tags['JUDGE_TAG_' + scope.row.verdict][0]">{{ tags['JUDGE_TAG_' + scope.row.verdict][1] }}</el-tag>
          </span>
          <span v-else>
            <el-tag>Unknown {{ scope.row.verdict }}</el-tag>
          </span>
        </template>
      </el-table-column>
      <el-table-column label="试题">
        <template slot-scope="scope">
          {{ scope.row.problem_title }}
        </template>
      </el-table-column>
      <el-table-column label="用户">
        <template slot-scope="scope">
          <span v-if="scope.row.user_name != ''">
            {{ scope.row.user_name }}
          </span>
          <span v-else><i>Unknown</i></span>
        </template>
      </el-table-column>
      <el-table-column label="提交时间">
        <template slot-scope="scope">
          {{ scope.row.submit_time | moment("from", "now") }}
        </template>
      </el-table-column>
    </el-table>
    </el-col>

    <el-col :span="24" :md="8" :lg="6">
      <div class="problem-sidebar">
        <el-card
          class="box-card"
          shadow="hover"
          :body-style="{padding: '0'}"
          style="margin-bottom: 20px"
        >
          <div class="problem-sidebar-item">
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
              <div
                class="timerText"
                style="padding: 3px 0"
              >{{ contestTimer.remain }}</div>
            </div>
            <div v-else-if="contestTimer.state == 'END'">
              <span>已结束</span>
            </div>
          </div>
        </el-card>

        <el-card
          class="box-card"
          shadow="hover"
          :body-style="{padding: '0'}"
          style="margin-bottom: 20px"
        >
          <div @click="gotoProblem">
            <div
              class="problem-sidebar-item"
              v-for="p in contest.problem_list"
              :key="p.uid"
              :data-click-url="`/${$route.params.domain}/contest/${$route.params.uid}/submissions/${p.uid}`"
            >
              <span class="problem-charid">{{ p.$charid }}.</span>
              {{ p.title }}
            </div>
          </div>
        </el-card>
      </div>
    </el-col>
  </el-row>
</template>

<script>
import { RPC } from "../rpc.js";
import { ConstString } from "../consts.js";
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
      languageAbbr: {
        "c.gcc99": "C",
        "cpp.gxx98": "C++",
        "cpp.gxx11": "C++",
        "java.java18": "Java",
        "py.py36": "Python",
        "cs.mono": "C#",
        "kotlin.default": "Kotlin",
        "hs.ghc7": "Haskell",
        "php.php7": "PHP"
      }
    };
  },
  methods: {
    loadSubmission: async function() {
      this.loading = true;
      this.filterApplied = false;
      let filter = (this.$route.params.filter || "").split(";");
      while (filter.length < 3) {
        filter.push("");
      }
      filter.forEach(element => {
        if (element != "") {
          this.filterApplied = true;
        }
      });
      filter[2] = filter[2].split(",");
      //
      let res = await RPC.doRPC("getContestSubmissions", {
        domain: this.$route.params.domain,
        filter: this.$route.params.uid + "|" + (this.$route.params.filter || "")
      });
      this.loading = false;
      if (res == null) {
        this.error = true;
        return;
      }
      this.submissions = res.submissions;
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
    tableIsSelf({ row, rowIndex }) {
      if (row.user_uid === this.user.uid) {
        return "selfRankItem";
      }
      return "";
    },
    clearFilter() {
      this.$router.push(
        `/${this.$route.params.domain}/contest/${
          this.$route.params.uid
        }/submissions`
      );
    },
    gotoProblem: function(e) {
      this.$router.push(e.target.dataset.clickUrl);
    },
  },
  computed: mapState(["user", "contest", "contestTimer"]),
  watch: {
    $route: function() {
      this.loadSubmission();
    }
  },
  created() {
    this.loadSubmission();
    this.loadingInterval = setInterval(this.loadSubmission, 3000);
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
.el-table .selfRankItem {
  background: #ecf5ff;
}

#baseinfo_container label {
  width: 72pt;
  color: #99a9bf;
}

.submission_metainfo_container {
  padding: 0;
}
</style>


