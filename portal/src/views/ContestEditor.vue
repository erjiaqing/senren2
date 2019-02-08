<template>
  <el-row>
    <el-col :span="24">
      <div class="grid-content problem-title">
        <span v-if="contest">
          <el-input
            v-model="contest.title"
            placeholder="比赛标题"
          ></el-input>
        </span>
        <span v-else-if="loading">Loading...</span>
        <span v-else-if="error">「比赛加载失败」</span>
      </div>
    </el-col>
    <el-col :span="24">
      <div class="grid-content problem-title">
        <el-button-group>
          <el-button
            icon="el-icon-upload2"
            @click="saveContest"
          >保存比赛</el-button>
          <el-button
            icon="el-icon-download"
            @click="loadContest"
          >恢复比赛</el-button>
          <el-button icon="el-icon-share">克隆比赛</el-button>
        </el-button-group>
      </div>
    </el-col>
    <el-col :span="24">
      <div
        class="grid-content problem-content"
        v-if="contest"
      >
        <div id="baseinfo_container">
          <el-alert
            title="以下时间均为本地时区时间"
            type="info"
            :closable="false"
          >
          </el-alert>
          <el-form
            label-position="left"
            inline
            class="submission_metainfo_container"
          >
            <el-form-item
              class="submission_metainfo_item"
              label="入场时间"
            >
              <el-date-picker
                v-model="contest.open_time"
                type="datetime"
                placeholder="选择日期时间"
                align="right"
              >
              </el-date-picker>
            </el-form-item>
            <el-form-item
              class="submission_metainfo_item"
              label="存档时间"
            >
              <el-date-picker
                v-model="contest.close_time"
                type="datetime"
                placeholder="选择日期时间"
                align="right"
              >
              </el-date-picker>
            </el-form-item>
            <el-form-item
              class="submission_metainfo_item"
              label="开始时间"
            >
              <el-date-picker
                v-model="contest.start_time"
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
                v-model="contest.end_time"
                type="datetime"
                placeholder="选择日期时间"
                align="right"
              >
              </el-date-picker>
            </el-form-item>
            <el-form-item
              class="submission_metainfo_item"
              label="封榜时间"
            >
              <el-date-picker
                v-model="contest.freeze_time"
                type="datetime"
                placeholder="选择日期时间"
                align="right"
              >
              </el-date-picker>
            </el-form-item>
            <el-form-item
              class="submission_metainfo_item"
              label="开榜时间"
            >
              <el-date-picker
                v-model="contest.release_time"
                type="datetime"
                placeholder="选择日期时间"
                align="right"
              >
              </el-date-picker>
            </el-form-item>
          </el-form>
        </div>
        <quill-editor
          v-if="contest"
          v-model="contest.description"
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
      contest: null,
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
      },
      freezeTimePickerOptions: {}
    };
  },
  components: {
    quillEditor
  },
  methods: {
    loadContest: async function() {
      this.loading = true;
      let res = await RPC.doRPC("getContest", {
        domain: this.$route.params.domain,
        uid: this.$route.params.uid
      });
      this.loading = false;
      if (res == null) {
        this.error = true;
        return;
      }
      this.contest = res.contest;
    },
    saveContest: async function() {
      this.loading = true;
      if (this.contest.domain == "") {
        this.contest.domain = this.$route.params.contest;
      }
      this.contest.open_time = new Date(this.contest.open_time);
      this.contest.close_time = new Date(this.contest.close_time);
      this.contest.start_time = new Date(this.contest.start_time);
      this.contest.end_time = new Date(this.contest.end_time);
      this.contest.freeze_time = new Date(this.contest.freeze_time);
      this.contest.release_time = new Date(this.contest.release_time);
      let res = await RPC.doRPC("createContest", {
        contest: this.contest
      });
      this.loading = false;
      if (res == null) {
        this.error = true;
        return;
      }
      if (this.contest.uid == "" || this.contest.uid == "aaaaaaaaaaaaaaaa") {
        this.contest.uid = res.uid;
        this.$router.push(
          `/${this.$route.params.domain}/contest/${res.uid}/edit`
        );
      }
    }
  },
  watch: {
    $route: function() {
      this.loadContest();
    }
  },
  created() {
    this.loadContest();
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


