<template>
  <el-row>
    <el-col :span="24">
      <div v-if="error">
        <el-alert
          title="请求失败"
          type="error"
          description="可能的原因：服务器故障、网络问题或比赛不存在"
          show-icon
        >
        </el-alert>
      </div>
      <div
        class="grid-content problem-content"
        v-if="contest"
      >
        <el-table
          :data="contest.problem_list"
          style="width: 100%"
          @row-click="gotoProblem"
        >
          <el-table-column
            label=""
            width="100px"
          >
            <template slot-scope="scope">
              {{ probIndex.charAt(scope.$index) }}
            </template>
          </el-table-column>
          <el-table-column
            prop="title"
            label="标题"
          >
          </el-table-column>
        </el-table>
      </div>
    </el-col>
  </el-row>
</template>

<script>
import { RPC } from "../rpc.js";
import {mapState} from "vuex";

export default {
  data() {
    return {
      probIndex: "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
      contest: null,
      loading: false,
      error: false
    };
  },
  methods: {
    loadContest: async function() {
      this.loading = true;
      let res = await RPC.doRPC("getContest", {
        domain: this.$route.params.domain,
        uid: this.$route.params.uid
      });
      this.loading = false;
      if (res == null) {
        this.error = true;
        return;
      }
      res.contest.problem_list = JSON.parse(res.contest.problem_list);
      for (let i = 0; i < res.contest.problem_list.length; i++) {
        res.contest.problem_list[i].$index = i;
        res.contest.problem_list[i].$charid = this.probIndex.charAt(i);
      }
      this.contest = res.contest;
    },
    gotoEditor: function() {
      this.$router.push(
        `/${this.$route.params.domain}/contest/${this.$route.params.uid}/edit`
      );
    },
    gotoProblem: function(e) {
      this.$router.push(
        `/${this.$route.params.domain}/contest/${
          this.$route.params.uid
        }/problem/${e.$charid}`
      );
    }
  },
  computed: mapState(["user"]),
  watch: {
    $route: function() {
      this.loadContest();
    }
  },
  created() {
    this.loadContest();
  }
};
</script>

<style>
.problem-title {
  font-size: 28pt;
  padding-bottom: 12pt;
}

.problem-content {
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

#contest_desc {
  padding: 8px;
  margin-top: 8px;
  margin-bottom: 8px;
  border-radius: 4px;
  -moz-box-shadow: 0px 0px 5px #999;
  -webkit-box-shadow: 0px 0px 5px #999;
  box-shadow: 0px 0px 5px #999;
}
</style>


