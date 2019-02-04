<template>
  <el-row>
    <el-col :span="24">
      <div class="grid-content problem-title">
        <span v-if="problem">
          <el-input
            v-model="problem.title"
            placeholder="试题标题"
          ></el-input>
        </span>
        <span v-else-if="loading">Loading...</span>
        <span v-else-if="error">「试题加载失败」</span>
      </div>
    </el-col>
    <el-col :span="24">
      <div class="grid-content problem-title">
        <el-button-group>
          <el-button
            icon="el-icon-upload2"
            @click="saveProblem"
          >保存试题</el-button>
          <el-button
            icon="el-icon-download"
            @click="loadProblem"
          >恢复试题</el-button>
          <el-button icon="el-icon-share">克隆试题</el-button>
        </el-button-group>
      </div>
    </el-col>
    <el-col :span="24">
      <div
        class="grid-content problem-content"
        v-if="problem"
      >
        <div id="baseinfo_container">
          <el-form
            label-position="left"
            inline
            class="submission_metainfo_container"
          >
            <el-form-item
              class="submission_metainfo_item"
              label="公开时间"
            >
              <el-date-picker
                v-model="problem.release"
                type="datetime"
                placeholder="选择日期时间"
                align="right"
                :picker-options="releaseTimePickerOptions"
              >
              </el-date-picker>
            </el-form-item>
            <el-form-item
              class="submission_metainfo_item"
              label="试题别名"
            >
            <el-input v-model="problem.alias" placeholder=""></el-input>
            </el-form-item>
          </el-form>
        </div>
        <quill-editor
          v-if="problem"
          v-model="problem.description"
          ref="Editor"
        >
        </quill-editor>
      </div>
      <div v-if="error">
        <el-alert
          title="试题加载失败"
          type="error"
          description="可能的原因：服务器故障、网络问题或请求的试题不存在"
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
      problem: null,
      loading: false,
      error: false,
      releaseTimePickerOptions: {
        shortcuts: [
          {
            text: "一直可见",
            onClick(picker) {
              picker.$emit("pick", new Date(0));
            }
          },
          {
            text: "不可见",
            onClick(picker) {
              picker.$emit("pick", new Date(2147483647000));
            }
          },
        ]
      }
    };
  },
  components: {
    quillEditor
  },
  methods: {
    loadProblem: async function() {
      this.loading = true;
      let res = await RPC.doRPC("getProblem", {
        domain: this.$route.params.domain,
        uid: this.$route.params.uid
      });
      this.loading = false;
      if (res == null) {
        this.error = true;
        return;
      }
      res.problem.release = new Date(res.problem.release);
      this.problem = res.problem;
    },
    saveProblem: async function() {
      this.loading = true;
      this.problem.release = new Date(this.problem.release);
      let res = await RPC.doRPC("createProblem", {
        problem: this.problem
      });
      this.loading = false;
      if (res == null) {
        this.error = true;
        return;
      }
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


