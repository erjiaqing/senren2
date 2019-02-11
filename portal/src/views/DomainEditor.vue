<template>
  <el-row>
    <el-col :span="24">
      <div class="grid-content problem-title">
        <span v-if="domain">
          <el-input
            v-model="domain.title"
            placeholder="小组标题"
          ></el-input>
        </span>
        <span v-else-if="loading">Loading...</span>
        <span v-else-if="error">「小组加载失败」</span>
      </div>
    </el-col>
    <el-col :span="24">
      <div class="grid-content problem-title">
        <el-button-group>
          <el-button
            icon="el-icon-upload2"
            @click="saveDomain"
          >保存更改</el-button>
          <el-button
            icon="el-icon-download"
            @click="loadDomain"
          >撤销更改</el-button>
        </el-button-group>
      </div>
    </el-col>
    <el-col :span="24">
      <div
        class="grid-content problem-content"
        v-if="domain"
      >
        <div id="baseinfo_container">
          <el-form
            label-position="left"
            inline
            class="submission_metainfo_container"
          >
            <el-form-item
              class="submission_metainfo_item"
              label="小组别名"
            >
              <el-input
                v-model="domain.alias"
                placeholder=""
              ></el-input>
            </el-form-item>
            <el-form-item label="访问权限">
              <el-radio-group v-model="domain.is_public">
                <el-radio label="PUBLIC">公开</el-radio>
                <el-radio label="PROTECTED">保护</el-radio>
                <el-radio label="PRIVATE">私有</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-form>
        </div>
        <quill-editor
          v-if="domain"
          v-model="domain.description"
          ref="Editor"
        >
        </quill-editor>
      </div>
      <div v-if="error">
        <el-alert
          title="小组加载失败"
          type="error"
          description="可能的原因：服务器故障、网络问题或请求的小组不存在"
          show-icon
        >
        </el-alert>
      </div>
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
    saveDomain: async function() {
      this.loading = true;
      let res = await RPC.doRPC("createDomain", {
        domain: this.domain
      });
      this.loading = false;
      if (res == null) {
        this.error = true;
        return;
      }
      if (this.domain.uid == "" || this.domain.uid == "aaaaaaaaaaaaaaaa") {
        this.domain.uid = res.uid;
        this.$message({
          message: "小组已创建",
          type: "success"
        });
        this.$router.push(`/${res.uid}/edit`);
      } else {
        this.$message({
          message: "小组已保存",
          type: "success"
        });
      }
    }
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
</style>


