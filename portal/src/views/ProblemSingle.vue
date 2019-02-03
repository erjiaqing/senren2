<template>
  <el-row>
    <el-col :span="24">
      <div class="grid-content problem-title">
        <span v-if="problem">{{ problem.title }}</span>
        <span v-else>Loading...</span>
      </div>
    </el-col>
    <el-col :span="24">
      <div class="grid-content problem-title">
        <el-button-group>
          <el-button icon="el-icon-upload2">提交代码</el-button>
          <el-button icon="el-icon-more">评测结果</el-button>
          <el-button icon="el-icon-tickets">讨论区</el-button>
          <el-button icon="el-icon-edit">编辑试题</el-button>
          <el-button icon="el-icon-share">克隆试题</el-button>
        </el-button-group>
      </div>
    </el-col>
    <el-col :span="24">
      <div class="grid-content problem-content" v-if="problem" v-html="problem.description">
      </div>
    </el-col>
  </el-row>
</template>

<script>
import {RPC} from '../rpc.js'

export default {
  data() {
    return {
      problem: null,
      loading: false,
      error: false,
    };
  },
  methods: {
    loadProblem: async function () {
      loading: true
      let res = await RPC.doRPC("getProblem", {
        domain: this.$route.params.domain,
        uid: this.$route.params.uid,
      });
      if (res == null) {
        error = true
        loading = false
        return
      }
      this.problem = res.problem;
    }
  },
  created() {
    this.loadProblem();
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
</style>


