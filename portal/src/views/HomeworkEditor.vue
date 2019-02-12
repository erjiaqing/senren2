<template>
  <el-row>
    <el-col :span="24">
      <div class="grid-content problem-title">
        <span v-if="homework">
          <el-input
            v-model="homework.title"
            placeholder="作业标题"
          ></el-input>
        </span>
        <span v-else-if="loading">Loading...</span>
        <span v-else-if="error">「作业加载失败」</span>
      </div>
    </el-col>
    <el-col :span="24">
      <div class="grid-content problem-title">
        <el-button-group>
          <el-button
            icon="el-icon-upload2"
            @click="saveHomework"
          >保存更改</el-button>
          <el-button
            icon="el-icon-download"
            @click="loadHomework"
          >撤销更改</el-button>
        </el-button-group>
      </div>
    </el-col>
    <el-col :span="24">
      <div
        class="grid-content problem-content"
        v-if="homework"
      >
        <div id="baseinfo_container">
          <el-form
            label-position="left"
            inline
            class="submission_metainfo_container"
          >
            <el-form-item
              class="submission_metainfo_item"
              label="开始时间"
            >
              <el-date-picker
                v-model="homework.start_time"
                type="datetime"
                placeholder="选择日期时间"
                align="right"
              >
              </el-date-picker>
            </el-form-item>
            <el-form-item
              class="submission_metainfo_item"
              label="结束时间"
            >
              <el-date-picker
                v-model="homework.end_time"
                type="datetime"
                placeholder="选择日期时间"
                align="right"
              >
              </el-date-picker>
            </el-form-item>
          </el-form>
        </div>
        <quill-editor
          v-if="homework"
          v-model="homework.description"
          ref="Editor"
        >
        </quill-editor>
      </div>
      <div v-if="error">
        <el-alert
          title="作业加载失败"
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
      homework: null,
      loading: false,
      error: false
    };
  },
  components: {
    quillEditor
  },
  methods: {
    loadHomework: async function() {
      this.loading = true;
      let res = await RPC.doRPC("getHomework", {
        domain: this.$route.params.domain,
        uid: this.$route.params.uid,
      });
      this.loading = false;
      if (res == null) {
        this.error = true;
        return;
      }
      this.homework = res.homework;
    },
    saveHomework: async function() {
      this.loading = true;
      this.homework.start_time = new Date(this.homework.start_time);
      this.homework.end_time = new Date(this.homework.end_time);
      this.homework.domain = this.$route.params.domain;
      let res = await RPC.doRPC("createHomework", {
        homework: this.homework
      });
      this.loading = false;
      if (res == null) {
        this.error = true;
        return;
      }
      if (this.homework.uid == "" || this.homework.uid == "aaaaaaaaaaaaaaaa") {
        this.homework.uid = res.uid;
        this.$message({
          message: "作业已创建",
          type: "success"
        });
        this.$router.push(`/${this.$route.params.domain}/homework/${res.uid}/edit`);
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
      this.loadHomework();
    }
  },
  created() {
    this.loadHomework();
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


