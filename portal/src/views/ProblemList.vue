<template>
  <div>
    <div class="grid-content problem-title" style="text-align:left">
      <el-button-group>
        <el-button @click="$router.push('/' + $route.params.domain + '/problem/aaaaaaaaaaaaaaaa/edit')">+ 创建试题</el-button>
      </el-button-group>
    </div>
    <el-table
      :data="problems"
      style="width: 100%"
    >
      <el-table-column
        label=""
        width="72pt"
      >
        <template slot-scope="scope">
          <router-link :to="'/' + $route.params.domain + '/problem/' + scope.row.uid">{{ scope.row.alias }}</router-link>
        </template>
      </el-table-column>
      <el-table-column label="标题">
        <template slot-scope="scope">
          <router-link :to="'/' + $route.params.domain + '/problem/' + scope.row.uid">{{ scope.row.title }}</router-link>
        </template>
      </el-table-column>
    </el-table>
    <div class="grid-content problem-title" style="text-align:left">
      <el-button-group>
        <el-button @click="$router.push('/' + $route.params.domain + '/problem/aaaaaaaaaaaaaaaa/edit')">+ 创建试题</el-button>
      </el-button-group>
    </div>
  </div>
</template>

<script>
import { RPC } from "../rpc.js";

export default {
  data() {
    return {
      problems: [],
      loading: false,
      error: false,
      codeEditor: false
    };
  },
  methods: {
    loadProblem: async function() {
      this.loading = true;
      let res = await RPC.doRPC("getProblems", {
        domain: this.$route.params.domain
      });
      this.loading = false;
      if (res == null) {
        this.error = true;
        return;
      }
      this.problems = res.problems;
    }
  },
  created() {
    this.loadProblem();
  }
};
</script>