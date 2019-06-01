<template>
  <el-row>
    <el-col :span="24">
      <div v-if="error">
        <el-alert title="请求失败" type="error" description="可能的原因：服务器故障、网络问题或比赛不存在" show-icon></el-alert>
      </div>
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
      <div class="grid-content problem-content" v-if="contest">
        <el-table :data="contest.problem_list" style="width: 100%" @row-click="gotoProblem">
          <el-table-column label width="100px">
            <template slot-scope="scope">{{ probIndex.charAt(scope.$index) }}</template>
          </el-table-column>
          <el-table-column prop="title" label="标题"></el-table-column>
        </el-table>
      </div>
    </el-col>
  </el-row>
</template>

<script>
import { RPC } from "../rpc.js";
import { mapState } from "vuex";

export default {
  data() {
    return {
      probIndex: "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
      loading: false,
      error: false
    };
  },
  methods: {
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
  computed: mapState(["user", "contestTimer", "contest"]),
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


