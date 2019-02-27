<template>
  <div>
    <div
      class="grid-content problem-title"
      style="text-align:left"
    >
      <el-form
        :inline="true"
        class="demo-form-inline"
      >
        <el-form-item label="试题标题">
          <el-input
            autocomplete="off"
            v-model="newProblemTitle"
            placeholder="用来让**你**区分不同的题"
          ></el-input>
        </el-form-item>
        <el-form-item>
          <el-button @click="doCreateProblem">+ 创建试题</el-button>
        </el-form-item>
      </el-form>
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
          <router-link :to="'/problem/' + scope.row.uid">{{ scope.row.uid }}</router-link>
        </template>
      </el-table-column>
      <el-table-column label="标题">
        <template slot-scope="scope">
          <router-link :to="'/problem/' + scope.row.uid">{{ scope.row.title }}</router-link>
        </template>
      </el-table-column>
    </el-table>
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
      newProblemTitle: ""
    };
  },
  methods: {
    doCreateProblem: async function() {
      this.loading = true;
      if (this.newProblemTitle == "") {
        return;
      }
      let res = await RPC.doRPC("createProblem", {
        problem: {
          uid: -1,
          title: this.newProblemTitle
        }
      });
      this.newProblemTitle = "";
      this.loadProblem();
    },
    loadProblem: async function() {
      this.loading = true;
      let res = await RPC.doRPC("getProblems", {});
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