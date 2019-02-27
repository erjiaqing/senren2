<template>
  <el-row>
    <el-col :span="24">
      <div class="grid-content problem-title">
        <span v-if="problem">#{{ problem.uid}}. {{ problem.title }}</span>
        <span v-else-if="loading">Loading...</span>
        <span v-else-if="error">「试题加载失败」</span>
      </div>
    </el-col>
    <el-alert
      title="试题设置已被修改，需要保存以生效"
      v-if="updated"
      type="error"
    >
    </el-alert>
    <el-col
      :span="24"
      style="text-align:left"
    >
      <div
        id="baseinfo_container"
        v-if="problem"
      >
        <el-form
          label-position="left"
          inline
          class="submission_metainfo_container"
        >
          <el-form-item
            class="submission_metainfo_item"
            label="试题ID"
          >
            #<span style="font-family: monospace">{{ problem.uid }}</span>
          </el-form-item>
          <el-form-item
            class="submission_metainfo_item"
            label="线上版本"
          >
            <span v-if="problem.version != ''">
              <code>{{problem.version}}</code>
            </span>
            <span v-else>
              <el-tag type="warning">未指定 / 暂缓评测</el-tag>
            </span>
          </el-form-item>
          <el-form-item
            class="submission_metainfo_item"
            label="编辑会话"
          >
            <el-button @click="doCreateEditSession">+ 编辑试题</el-button>
          </el-form-item>
          <el-form-item
            class="submission_metainfo_item"
            label="保存试题"
          >
            <el-button @click="doSaveProblem">+ 保存试题</el-button>
          </el-form-item>
          <el-form-item
            class="submission_metainfo_item"
            label="创建密钥"
          >
            <el-button @click="doCreateKey">+ 创建密钥</el-button>
          </el-form-item>
        </el-form>
        <el-table
          :data="versions"
          style="width: 100%"
          @row-click="handleVersionClick"
          :row-class-name="versionState"
        >
          <el-table-column width="50pt">
            <template slot-scope="scope">
              <i
                class="el-icon-success"
                v-if="scope.row.version == problem.version"
              ></i>
            </template>
          </el-table-column>
          <el-table-column
            label="版本号"
            width="200pt"
          >
            <template slot-scope="scope">
              <code>{{ scope.row.version.substr(0,7) }}</code>
            </template>
          </el-table-column>
          <el-table-column label="">
            <template slot-scope="scope">
              {{ scope.row.message.split("\n")[0] }}
            </template>
          </el-table-column>
          <el-table-column label="创建时间">
            <template slot-scope="scope">
              {{ scope.row.created | moment("from", "now") }}
            </template>
          </el-table-column>
        </el-table>
        <el-table
          :data="accessKeys"
          style="width: 100%"
        >
          <el-table-column
            label="密钥"
          >
            <template slot-scope="scope">
              <code>{{ scope.row.key }}</code>
            </template>
          </el-table-column>
          <el-table-column label="权限">
            <template slot-scope="scope">
              {{ scope.row.access_control }}
            </template>
          </el-table-column>
          <el-table-column label="创建时间">
            <template slot-scope="scope">
              {{ scope.row.create_time | moment("llll") }}
            </template>
          </el-table-column>
        </el-table>
            <el-alert
      title="密钥是供外部程序（如WOJ）访问试题的一个字符串，将密钥输入到它们当中可以让它们按限制的权限访问试题，如果要供WOJ提交，则需要 .PROBLEM 权限"
      type="info"
    >
    </el-alert>
      </div>
    </el-col>
  </el-row>
</template>

<script>
import { RPC } from "../rpc.js";
import { mapState } from "vuex";

export default {
  data() {
    return {
      problem: null,
      updated: false,
      loading: false,
      error: false,
      versions: [],
      accessKeys: [],
    };
  },
  methods: {
    doCreateKey: async function() {
      this.loading = true;
      let res = await RPC.doRPC("createProblemAccessKey", {
        uid: Number(this.$route.params.uid),
        perms: ".PROBLEM",
      });
    },

    loadProblem: async function() {
      this.loading = true;
      let res = await RPC.doRPC("getProblem", {
        uid: Number(this.$route.params.uid)
      });
      this.loading = false;
      if (res == null) {
        this.error = true;
        return;
      }
      this.problem = res.problem;
      this.updated = false;
      this.loadProblemVersions();
      this.loadProblemKeys();
    },

    versionState({ row, rowIndex }) {
      if (row.state === "UNKNOWN") {
        return "warning-row";
      } else if (row.state === "CHECKED") {
        return "success-row";
      } else if (row.state === "FAILED") {
        return "error-row";
      }
      return "";
    },
    doSaveProblem: async function() {
      await RPC.doRPC("createProblem", {
        problem: this.problem
      });
      this.loadProblem();
    },
    loadProblemVersions: async function() {
      this.loading = true;
      let res = await RPC.doRPC("getProblemVersions", {
        uid: Number(this.$route.params.uid)
      });
      this.loading = false;
      if (!res) {
        return;
      }
      this.versions = res.versions;
    },
    loadProblemKeys: async function() {
      this.loading = true;
      let res = await RPC.doRPC("getProblemAccessKeys", {
        uid: Number(this.$route.params.uid)
      });
      this.loading = false;
      if (!res) {
        return;
      }
      this.accessKeys = res.keys;
    },
    handleVersionClick(e) {
      console.log(e);
      this.problem.version = e.version;
      this.updated = true;
    },
    doCreateEditSession: async function() {
      this.loading = true;
      let res = await RPC.doRPC("createProblemEditSession", {
        uid: Number(this.$route.params.uid)
      });
      this.loading = false;
      if (!res) {
        return;
      }
      const h = this.$createElement;
      this.$msgbox({
        title: "启动编辑会话",
        message: h("p", { style: "text-align:center" }, [
          h("p", null, [
            h(
              "a",
              {
                attrs: { href: "http://127.0.0.1:8084/" + res.uid },
                style: "font-size:24px"
              },
              "单击这里进入"
            )
          ]),
          h("p", null, "链接24小时内有效")
        ])
      });
      this.loadProblem();
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

.problem-title {
  font-size: 28pt;
  padding-bottom: 12pt;
}

.problem-content {
  text-align: left;
  font-size: 14pt;
}

.codeEditorPanel {
  -moz-box-shadow: 0px 0px 5px #999;
  -webkit-box-shadow: 0px 0px 5px #999;
  box-shadow: 0px 0px 5px #999;
  border-radius: 5px;

  transition: all 0.3s linear;
}

.codeEditorPanel:hover {
  -moz-box-shadow: 0px 0px 15px #999;
  -webkit-box-shadow: 0px 0px 15px #999;
  box-shadow: 0px 0px 15px #999;
  border-radius: 5px;
}

  .el-table .warning-row {
    background: oldlace;
  }

  .el-table .success-row {
    background: #f0f9eb;
  }
</style>


