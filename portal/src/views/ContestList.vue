<template>
  <div>
    <div
      class="grid-content problem-title"
      style="text-align:left"
      v-if="user && (user.role == 'ADMIN' || user.role == 'ROOT')"
    >
      <el-button-group>
        <el-button @click="$router.push('/' + $route.params.domain + '/contest/aaaaaaaaaaaaaaaa/edit')">+ 创建比赛</el-button>
      </el-button-group>
    </div>
    <el-alert
      title="以下时间均为本地时区时间"
      type="info"
      :closable="false"
    >
    </el-alert>
    <el-table
      :data="contests"
      style="width: 100%"
    >
      <el-table-column label="标题">
        <template slot-scope="scope">
          <router-link :to="'/' + $route.params.domain + '/contest/' + scope.row.uid">{{ scope.row.title }}</router-link>
        </template>
      </el-table-column>
      <el-table-column label="比赛开始时间">
        <template slot-scope="scope">
          <a :href="'https://www.timeanddate.com/worldclock/fixedtime.html?msg=Contest&iso=' + scope.row.start_time">{{ scope.row.start_time | moment("YYYY-MM-DD HH:mm:ss") }}</a>
        </template>
      </el-table-column>
      <el-table-column label="比赛长度">
        <template slot-scope="scope">
          {{ Math.floor(((new Date(scope.row.end_time)).valueOf() - (new Date(scope.row.start_time)).valueOf()) / (60 * 1000)) }} 分钟
        </template>
      </el-table-column>
    </el-table>
    <div
      class="grid-content problem-title"
      style="text-align:left"
      v-if="user && (user.role == 'ADMIN' || user.role == 'ROOT')"
    >
      <el-button-group>
        <el-button @click="$router.push('/' + $route.params.domain + '/contest/aaaaaaaaaaaaaaaa/edit')">+ 创建比赛</el-button>
      </el-button-group>
    </div>
  </div>
</template>

<script>
import { mapState } from "vuex";
import { RPC } from "../rpc.js";

export default {
  data() {
    return {
      contests: [],
      loading: false,
      error: false,
      codeEditor: false
    };
  },
  methods: {
    loadContests: async function() {
      this.loading = true;
      let res = await RPC.doRPC("getContests", {
        domain: this.$route.params.domain
      });
      this.loading = false;
      if (res == null) {
        this.error = true;
        return;
      }
      this.contests = res.contests;
    }
  },
  computed: mapState(["user"]),
  created() {
    this.loadContests();
  }
};
</script>