<template>
  <el-row>
    <el-col :span="24">
      <div class="grid-content problem-title" style="text-align:left">
        <span v-if="domain">
          <h3><small id="domain-alias">/{{domain.alias}}/</small><br>{{ domain.title }}<br><small id="domain-powered">Powered by 千練萬花</small></h3>
        </span>
        <span v-else-if="loading">Loading...</span>
        <span v-else-if="error">「小组加载失败」</span>
      </div>
    </el-col>
    <el-col :span="24">
      <div
        class="grid-content problem-content"
        v-if="domain"
        v-html="domain.description"
      ></div>
    </el-col>
  </el-row>
</template>

<script>
import { RPC } from "../rpc.js";
import "quill/dist/quill.core.css";
import "quill/dist/quill.snow.css";
import "quill/dist/quill.bubble.css";
import { quillEditor } from "vue-quill-editor";

export default {
  data() {
    return {
      domain: null,
      loading: false,
      error: false
    };
  },
  components: {
    quillEditor
  },
  methods: {
    loadDomain: async function() {
      this.loading = true;
      let res = await RPC.doRPC("getDomain", {
        domain: this.$route.params.domain
      });
      this.loading = false;
      if (res == null) {
        this.error = true;
        return;
      }
      this.domain = res.domain;
    },
  },
  watch: {
    $route: function() {
      this.loadDomain();
    }
  },
  created() {
    this.loadDomain();
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

#baseinfo_container .el-form-item {
  margin-right: 0;
  margin-bottom: 0;
  margin-left: 0;
  padding: 0;
  width: 50%;
}

#baseinfo_container label {
  width: 72pt;
  color: #99a9bf;
}

#domain-alias {
  font-family: 'Courier New', Courier, monospace;
  font-size: 20pt;
  color: #260817;
}

#domain-powered {
  font-size: 14pt;
}
</style>


