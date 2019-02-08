<template>
  <el-row>
    <el-col :span="24">
      <div class="grid-content problem-title">
        <span v-if="contest">{{ contest.title }}</span>
        <span v-else-if="loading">Loading...</span>
        <span v-else-if="error">「试题加载失败」</span>
      </div>
    </el-col>
    <el-col :span="24">
      <div class="grid-content problem-title">
        <el-button-group>
          <el-button icon="el-icon-more">评测列表</el-button>
          <el-button icon="el-icon-list">榜单</el-button>
          <el-button icon="el-icon-tickets">讨论区</el-button>
          <el-button icon="el-icon-edit" @click="gotoEditor">编辑比赛</el-button>
          <el-button icon="el-icon-share">克隆比赛</el-button>
        </el-button-group>
      </div>
    </el-col>
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
        <div v-html="contest.description"></div>
      </div>
    </el-col>
  </el-row>
</template>

<script>
import { RPC } from "../rpc.js";

export default {
  data() {
    return {
      contest: null,
      loading: false,
      error: false,
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
      this.contest = res.contest;
    },
    gotoEditor: function() {
      this.$router.push(
        `/${this.$route.params.domain}/contest/${this.$route.params.uid}/edit`
      );
    },
  },
  watch: {
    '$route': function() {
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
</style>


