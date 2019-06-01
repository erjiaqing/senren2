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
          <!-- <el-button icon="el-icon-share">克隆试题</el-button> -->
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
              <el-input
                v-model="problem.alias"
                placeholder=""
              ></el-input>
            </el-form-item>
          </el-form>
          <el-form
            label-position="left"
            inline
            class="submission_pcikey_container"
          >
            <el-form-item
              class="submission_metainfo_item"
              label="PCI提交密钥"
            >
              <el-input
                v-model="problem.problem_ci"
                placeholder=""
                style="width: 400pt"
              ></el-input>
            </el-form-item>
            <el-form-item
              class="submission_metainfo_item"
              label="PCI 操作"
            >
              <el-button
                icon="el-icon-download"
                @click="loadProblemCI"
              >从PCI导入题面</el-button>
            </el-form-item>
          </el-form>
        </div>
        <div>
          <quill-editor
            v-if="problem && problem.description.substr(0, 21) != '<!-- !!imported!! -->'"
            v-model="problem.description"
            ref="Editor"
          >
          </quill-editor>
          <div v-else-if="problem">
            <el-alert
              title="编辑导入的试题可能会破坏ProblemCI特有格式。"
              close-text="仍然编辑"
              @close="problem.description = problem.description.substr(21)"
              type="error"
              style="margin-top: 4pt;"
            >
            </el-alert>
            <div
              id="problem-description"
              v-html="problem.description"
            >
            </div>
          </div>
        </div>
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
          }
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
    loadProblemCI: async function() {
      this.loading = true;
      let res = await RPC.doRPC("getPCIDescription", {
        filter: this.problem.problem_ci
      });
      this.loading = false;
      if (res == null || !res.success) {
        this.$message({
          message: "导入错误",
          type: "danger"
        });
      }
      this.problem.description = "<!-- !!imported!! -->" + res.description;
    },
    saveProblem: async function() {
      this.loading = true;
      this.problem.release = new Date(this.problem.release);
      if (this.problem.domain == "") {
        this.problem.domain = this.$route.params.domain;
      }
      let res = await RPC.doRPC("createProblem", {
        problem: this.problem,
        domain: this.$route.params.domain
      });
      this.loading = false;
      if (res == null) {
        this.error = true;
        return;
      }
      if (this.problem.uid == "" || this.problem.uid == "0000000000000000") {
        this.problem.uid = res.uid;
        this.$message({
          message: "试题已创建",
          type: "success"
        });
        this.$router.push(
          `/${this.$route.params.domain}/problem/${res.uid}/edit`
        );
      } else {
        this.$message({
          message: "试题已保存",
          type: "success"
        });
      }
    }
  },
  watch: {
    $route: function() {
      this.loadProblem();
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

#baseinfo_container .submission_pcikey_container .el-form-item {
  margin-right: 0;
  margin-bottom: 0;
  margin-left: 0;
  padding: 0;
  width: 100%;
}

#baseinfo_container .submission_pcikey_container .el-form-item input {
  font-family: "Courier New", Courier, monospace;
}

#baseinfo_container label {
  width: 72pt;
  color: #99a9bf;
}
</style>


