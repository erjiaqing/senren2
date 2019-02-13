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
            @click="$router.push('/' + $route.params.domain + '/homework/' + homework.uid + '/submissions')"
            v-if="user && (user.role == 'ADMIN' || user.role == 'ROOT')"
          >批量下载</el-button>
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
          <h3>提交结果</h3>
          <el-table
            :data="submission"
            border
            style="width: 100%"
          >
            <el-table-column
              prop="__nick"
              label="用户"
            >
            </el-table-column>
            <el-table-column
              v-for="(o) in attachment_list"
              :key="o.name"
              :label="o.name"
            >
              <template slot-scope="scope">
                <el-tag
                  size="small"
                  type="success"
                  v-if="scope.row[o.name]"
                >
                  <i class="el-icon-success"></i> {{ scope.row[o.name].fullname }}
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
            <el-table-column
              label="操作"
            >
              <template slot-scope="scope">
            <router-link :to="'/' + $route.params.domain + '/homework/' + $route.params.uid + '/show/' + scope.row.__uid">批改</router-link>
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
      submission: [],
      attachment_list: [],
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
    loadSubmission: async function() {
      this.loading = true;
      let res = await RPC.doRPC("getHomeworkSubmissions", {
        domain: this.$route.params.domain,
        filter: this.$route.params.uid
      });
      this.loading = false;
      if (res == null) {
        this.error = true;
        return;
      }
      if (res.success !== true) {
        return;
      }

      let sublist = [];
      res.homeworksubmissions.forEach(element => {
        let fullattachment = element.attachments.split("!!");
        if (fullattachment.length != 2) {
          return;
        }
        let nick = fullattachment[0];
        let sublst = fullattachment[1].split(";");

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
        subres["__nick"] = nick;
        subres["__uid"] = element.uid;
        sublist.push(subres);
      });
      console.log(sublist);
      this.submission = sublist;
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


