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
        <div id="problemList">
          <el-form
            :inline="true"
            class="demo-form-inline"
            label-width="120px"
          >
            <el-form-item label="试题UID或短ID">
              <el-input v-model="addProblemSelector.uid"></el-input>
            </el-form-item>
            <el-form-item>
              <el-button
                type="primary"
                @click="queryProblemToAdd"
              >查询</el-button>
            </el-form-item>
            <el-form-item>
              <span v-if="problemToAdd.uid">
                {{ problemToAdd.alias || "" }} {{ problemToAdd.title }}
                <el-button
                  type="success"
                  icon="el-icon-check"
                  circle
                  @click="addProblem"
                ></el-button>
              </span>
            </el-form-item>
          </el-form>
          <draggable v-model="usedProblem">
            <transition-group name="flip-list">
              <div
                v-for="(item, seq) in usedProblem"
                :key="item.uid"
              >
                <button @click="deleteProblem(seq)">x</button> {{ item.alias || item.uid }}. {{ item.title }} <code><small>UID: {{ item.uid }}</small></code>
              </div>
            </transition-group>
          </draggable>
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
import draggable from "vuedraggable";

export default {
  data() {
    return {
      contest: null,
      loading: false,
      error: false,
      usedProblem: [],
      addProblemSelector: {
        uid: ""
      },
      problemToAdd: {},
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
    quillEditor,
    draggable
  },
  methods: {
    deleteProblem: async function(seq) {
      let usedProblem = JSON.parse(JSON.stringify(this.usedProblem));
      usedProblem.splice(seq, 1);
      this.usedProblem = usedProblem;
    },
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
      this.usedProblem = JSON.parse(this.contest.problem_list);
    },
    addProblem: function() {
      let t = JSON.parse(JSON.stringify(this.usedProblem));
      t.push({
        uid: this.problemToAdd.uid,
        alias: this.problemToAdd.alias,
        title: this.problemToAdd.title
      });
      this.usedProblem = t;
      this.problemToAdd = {};
    },
    queryProblemToAdd: async function() {
      this.loading = true;
      let res = await RPC.doRPC("getProblem", {
        domain: this.$route.params.domain,
        uid: this.addProblemSelector.uid,
        filter: "+title"
      });
      if (!res) {
        this.problemToAdd = {};
      }
      this.problemToAdd = res.problem;
      this.loading = false;
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
      this.contest.problem_list = JSON.stringify(this.usedProblem);
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

#problemList {
  margin-top: 12pt;
  margin-bottom: 12pt;
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


