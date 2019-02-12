<template>
  <div>
    <div
      class="grid-content problem-title"
      style="text-align:left"
    >
      <el-button-group>
        <el-button @click="$router.push('/' + $route.params.domain + '/homework/aaaaaaaaaaaaaaaa/edit')">+ 创建作业</el-button>
      </el-button-group>
    </div>
    <el-alert
      title="以下时间均为本地时区时间"
      type="info"
      :closable="false"
    >
    </el-alert>
    <el-table
      :data="homeworks"
      style="width: 100%"
    >
      <el-table-column label="标题">
        <template slot-scope="scope">
          <router-link :to="'/' + $route.params.domain + '/homework/' + scope.row.uid">{{ scope.row.title }}</router-link>
        </template>
      </el-table-column>
      <el-table-column label="作业截止时间">
        <template slot-scope="scope">
          <a :href="'https://www.timeanddate.com/worldclock/fixedtime.html?msg=Homework&iso=' + scope.row.end_time">{{ scope.row.end_time | moment("llll") }}</a>
        </template>
      </el-table-column>
    </el-table>
    <div
      class="grid-content problem-title"
      style="text-align:left"
    >
      <el-button-group>
        <el-button @click="$router.push('/' + $route.params.domain + '/homework/aaaaaaaaaaaaaaaa/edit')">+ 创建作业</el-button>
      </el-button-group>
    </div>
  </div>
</template>

<script>
import { RPC } from "../rpc.js";

export default {
  data() {
    return {
      homeworks: [],
      loading: false,
      error: false,
      codeEditor: false
    };
  },
  methods: {
    loadHomeworks: async function() {
      this.loading = true;
      let res = await RPC.doRPC("getHomeworks", {
        domain: this.$route.params.domain
      });
      this.loading = false;
      if (res == null) {
        this.error = true;
        return;
      }
      this.homeworks = res.homeworks;
    }
  },
  created() {
    this.loadHomeworks();
  }
};
</script>