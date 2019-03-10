<template>
  <el-row>
    <el-col :span="24">
      <div class="grid-content problem-title">
        <span v-if="homework">{{ homework.title }}</span>
        <span v-else-if="loading">Loading...</span>
        <span v-else-if="error">「作业加载失败」</span>
      </div>
    </el-col>
    <el-col
      :span="24"
      v-if="homework"
    >
      <div class="grid-content problem-title">
        <el-button-group>
          <el-button
            icon="el-icon-more"
            @click="$router.push('/' + $route.params.domain + '/homework/' + homework.uid + '/edit')"
            v-if="user && (user.role == 'ADMIN' || user.role == 'ROOT')"
          >编辑作业</el-button>
          <el-button
            icon="el-icon-more"
            @click="$router.push('/' + $route.params.domain + '/homework/' + homework.uid + '/submissions')"
            v-if="user && (user.role == 'ADMIN' || user.role == 'ROOT')"
          >提交结果</el-button>
          <el-button icon="el-icon-tickets">讨论区</el-button>
        </el-button-group>
      </div>
    </el-col>
    <el-col :span="24">
      <div v-if="error">
        <el-alert
          title="请求失败"
          type="error"
          description="可能的原因：服务器故障、网络问题或试题不存在"
          show-icon
        >
        </el-alert>
      </div>
      <div
        class="grid-content problem-content"
        v-if="homework"
      >
        <div
          id="start-notify"
          v-if="isOverdue == -1"
        >
          <i class="el-icon-info"></i>
          作业提交开始于 {{ homework.start_time | moment("from", "now") }} ({{ homework.start_time | moment('llll')}})</div>
        <div id="ddl-notify">
          <i class="el-icon-info"></i>
          作业提交截止于 {{ homework.end_time | moment("from", "now") }} ({{ homework.end_time | moment('llll')}})</div>
        <div>
          <h3>作业描述</h3>
        </div>
        <div v-html="homework.description"></div>
        <div>
          <h3>提交项</h3>
          <div id="submit_form">
            <el-row>
              <el-col :span="12">
                <el-select
                  v-model="extData.name"
                  placeholder="选择需要提交的作业"
                >
                  <el-option
                    v-for="item in attachment_list"
                    :key="item.name"
                    :label="item.name"
                    :value="item.name"
                  >
                  </el-option>
                </el-select>
                <el-button
                  type="primary"
                  @click="generatePhoneUploadQR"
                >手机上传</el-button>
              </el-col>
              <el-col :span="12">
                <el-upload
                  class="upload-demo"
                  action="/rpc/attachments/uploadHomework"
                  :limit="1"
                  :headers="upload_head"
                  :data="extData"
                  :file-list="uplfile"
                  :on-success="onFileUploadSuccess"
                >
                  <el-button
                    type="primary"
                    :disabled="uplfile && uplfile.length > 0"
                  >点击上传</el-button>
                </el-upload>
              </el-col>
            </el-row>
          </div>
          <el-table
            :data="attachment_list"
            style="width: 100%"
            @row-click="handleHomeworkClick"
          >
            <el-table-column label="附件名称">
              <template slot-scope="scope">
                <span>{{ scope.row.name }}</span>
              </template>
            </el-table-column>
            <el-table-column label="文件类型">
              <template slot-scope="scope">
                <code>{{ scope.row.file }}</code>
              </template>
            </el-table-column>
            <el-table-column label="附件描述">
              <template slot-scope="scope">
                <span>{{ scope.row.desc }}</span>
              </template>
            </el-table-column>
            <el-table-column label="提交结果">
              <template slot-scope="scope">
                <el-tag
                  size="small"
                  type="success"
                  v-if="submission_result[scope.row.name]"
                >
                  <i class="el-icon-success"></i> {{ submission_result[scope.row.name].fullname }} , {{ parseSize(submission_result[scope.row.name].size) }}
                </el-tag>
                <el-tag
                  size="small"
                  v-else-if="isOverdue == -1"
                  type="info"
                >
                  <i class="el-icon-info"></i> 未开始
                </el-tag>
                <el-tag
                  size="small"
                  v-else-if="isOverdue == 0"
                  type="warning"
                >
                  <i class="el-icon-warning"></i> 待提交
                </el-tag>
                <el-tag
                  size="small"
                  v-else
                  type="danger"
                >
                  <i class="el-icon-error"></i> 未提交
                </el-tag>
              </template>
            </el-table-column>
          </el-table>
        </div>
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
import { mapState } from "vuex";
import QrcodeVue from "qrcode.vue";

export default {
  data() {
    return {
      homework: null,
      attachment_list: [],
      uplfile: [],
      submission_result: {},
      upload_head: {},
      extData: {
        name: ""
      },
      isOverdue: 0,
      loading: false,
      error: false
    };
  },
  components: {
    quillEditor,
    QrcodeVue
  },
  methods: {
    parseSize(size) {
      let retSuffix = ["B", "KB", "MB", "GB"];
      let sufSz = 0;
      size = Number(size);
      while (size > 0.5 * Math.pow(10, sufSz + 1)) {
        sufSz++;
        size /= 1000;
      }
      return Math.round(size * 100) / 100 + retSuffix[sufSz];
    },
    generatePhoneUploadQR: function() {
      let data =
        window.localStorage.getItem("sid") +
        "$" +
        this.homework.domain +
        ";" +
        this.homework.uid +
        ";" +
        this.extData.name;
      let qr = location.origin + "/uploadHomework/" + encodeURI(data);
      const h = this.$createElement;
      this.$msgbox({
        title: "扫描下方二维码使用手机上传作业",
        message: h("p", {style: "text-align:center"}, [
          h("qrcode-vue", { props: { value: qr, size: 256, level: "H" } }, ""),
          h("p", null, "上传完成后，请刷新该页面")
        ])
      });
    },
    onFileUploadSuccess: function(resp, file, fileList) {
      console.log(resp);
      if (resp.success) {
        this.$message({
          message: "作业上传成功",
          type: "success"
        });
      } else {
        this.$message({
          message: `作业上传失败:<br><strong>${resp.error}</strong>`,
          dangerouslyUseHTMLString: true,
          type: "error"
        });
      }
      this.loadSubmission();
    },
    loadSubmission: async function() {
      this.loading = true;
      let res = await RPC.doRPC("getHomeworkSubmission", {
        domain: this.$route.params.domain,
        uid: this.$route.params.uid
      });
      this.loading = false;
      if (res == null) {
        this.error = true;
        return;
      }
      if (res.success !== true) {
        return;
      }
      let sublst = res.homeworksubmission.attachments.split(";");
      let subres = {};
      sublst.forEach(element => {
        let l2 = element.split(",");
        if (l2.length < 5) {
          return;
        }
        subres[l2[0]] = {
          name: l2[0],
          uid: l2[2],
          origname: l2[1],
          size: Number(l2[3]),
          fullname: l2[0] + l2[4],
          extname: l2[4]
        };
      });
      this.submission_result = subres;
      console.log(res);
    },
    loadHomework: async function() {
      this.loading = true;
      let res = await RPC.doRPC("getHomework", {
        domain: this.$route.params.domain,
        uid: this.$route.params.uid
      });
      this.loading = false;
      if (res == null) {
        this.error = true;
        return;
      }
      this.homework = res.homework;
      let attachment_list = res.homework.attachments.split(";");
      let attachment_obj_list = [];
      attachment_list.forEach(element => {
        let te = element.split(",");
        if (element == "") {
          return;
        }
        attachment_obj_list.push({
          name: te[0],
          desc: te[1],
          file: te[2]
        });
      });
      let current = new Date().valueOf();
      if (current < new Date(res.homework.start_time).valueOf()) {
        this.isOverdue = -1;
      } else if (current > new Date(res.homework.end_time).valueOf()) {
        this.isOverdue = 1;
      } else {
        this.isOverdue = 0;
      }
      this.attachment_list = attachment_obj_list;
      this.upload_head = {
        UPLOAD_SESSION: window.localStorage.getItem("sid"),
        UPLOAD_DOMAIN: this.$route.params.domain,
        UPLOAD_HOMEWORK: this.$route.params.uid
      };
      this.loadSubmission();
    },
    handleHomeworkClick: function(e) {
      console.log(e);
    }
  },
  computed: mapState(["user"]),
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

#start-notify {
  padding: 8px;
  border-radius: 4px;
  color: #e6a23c;
  background-color: #fdf6ec;
}

#ddl-notify {
  padding: 8px;
  border-radius: 4px;
  color: #f56c6c;
  background-color: #fef0f0;
}

#submit_form {
  padding: 8px;
  margin-bottom: 8px;
  border-radius: 4px;
  -moz-box-shadow: 0px 0px 5px #999;
  -webkit-box-shadow: 0px 0px 5px #999;
  box-shadow: 0px 0px 5px #999;
}
</style>


