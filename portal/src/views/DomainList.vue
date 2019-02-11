<template>
  <div>
    <div class="grid-content problem-title" style="text-align:left">
      <el-button-group>
        <el-button @click="$router.push('//aaaaaaaaaaaaaaaa/edit')">+ 创建小组</el-button>
      </el-button-group>
    </div>
    <el-table
      :data="domains"
      style="width: 100%"
    >
      <el-table-column
        label=""
        width="72pt"
      >
        <template slot-scope="scope">
          <router-link :to="'/' + (scope.row.alias || scope.row.uid)">{{ scope.row.alias }}</router-link>
        </template>
      </el-table-column>
      <el-table-column label="标题">
        <template slot-scope="scope">
          <router-link :to="'/' + (scope.row.alias || scope.row.uid)">{{ scope.row.title }}</router-link>
        </template>
      </el-table-column>
    </el-table>
    <div class="grid-content problem-title" style="text-align:left">
      <el-button-group>
        <el-button @click="$router.push('/aaaaaaaaaaaaaaaa/edit')">+ 创建小组</el-button>
      </el-button-group>
    </div>
  </div>
</template>

<script>
import { RPC } from "../rpc.js";

export default {
  data() {
    return {
      domains: [],
      loading: false,
      error: false,
    };
  },
  methods: {
    loadDomain: async function() {
      this.loading = true;
      let res = await RPC.doRPC("getDomains", {
      });
      this.loading = false;
      if (res == null) {
        this.error = true;
        return;
      }
      this.domains = res.domains;
    }
  },
  created() {
    this.loadDomain();
  }
};
</script>