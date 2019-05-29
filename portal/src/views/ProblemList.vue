<template>
  <div>
    <div
      class="grid-content problem-title"
      style="text-align:left"
      v-if="user && (user.role == 'ADMIN' || user.role == 'ROOT')"
    >
      <el-button-group>
        <el-button
          @click="$router.push('/' + $route.params.domain + '/problem/aaaaaaaaaaaaaaaa/edit')"
        >+ 创建试题</el-button>
      </el-button-group>
    </div>

    <div class="block">
      <el-pagination
        @current-change="refreshCurrentPage"
        :current-page="currentPage"
        :page-size="countPerPage"
        layout="total, prev, pager, next"
        :total="totalCount"
      ></el-pagination>
    </div>

    <el-table :data="problems" style="width: 100%">
      <el-table-column label width="72pt">
        <template slot-scope="scope">
          <router-link
            :to="'/' + $route.params.domain + '/problem/' + scope.row.uid"
          >{{ scope.row.alias }}</router-link>
        </template>
      </el-table-column>
      <el-table-column label="标题">
        <template slot-scope="scope">
          <router-link
            :to="'/' + $route.params.domain + '/problem/' + scope.row.uid"
          >{{ scope.row.title }}</router-link>
        </template>
      </el-table-column>
    </el-table>

    <div class="block">
      <el-pagination
        @current-change="refreshCurrentPage"
        :current-page="currentPage"
        :page-size="countPerPage"
        layout="total, prev, pager, next"
        :total="totalCount"
      ></el-pagination>
    </div>

    <div
      class="grid-content problem-title"
      style="text-align:left"
      v-if="user && (user.role == 'ADMIN' || user.role == 'ROOT')"
    >
      <el-button-group>
        <el-button
          @click="$router.push('/' + $route.params.domain + '/problem/aaaaaaaaaaaaaaaa/edit')"
        >+ 创建试题</el-button>
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
      problems: [],
      totalCount: 0,
      currentPage: 1,
      countPerPage: 50,
      loading: false,
      error: false,
      codeEditor: false
    };
  },
  methods: {
    loadProblem: async function() {
      this.loading = true;
      let currPage = Number(this.$route.params ? this.$route.params.page : 1);
      this.currentPage = currPage;
      let res = await RPC.doRPC("getProblems", {
        domain: this.$route.params.domain,
        page: currPage,
        count: this.countPerPage
      });
      this.loading = false;
      if (res == null) {
        this.error = true;
        return;
      }
      this.currentPage = res.page;
      this.totalCount = res.total;
      this.problems = res.problems;
    },
    refreshCurrentPage: function(val) {
      this.$router.replace(
        `/${this.$route.params.domain}/problems/page/${val}`
      );
      this.loadProblem();
    }
  },
  computed: mapState(["user"]),
  created() {
    this.loadProblem();
  },
  watch: {
    $route(from, to) {
      if (from != to) {
        this.loadProblem();
      }
    }
  }
};
</script>