<template>
  <el-row>
    <el-col :span="24">
      <div class="grid-content problem-title">
        <span v-if="homework">{{ homework.title }}</span>
        <span v-else-if="loading">Loading...</span>
        <span v-else-if="error">「作业加载失败」</span>
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
          <el-row>
            <el-col :span="8">
              <div>{{ submission_result.__nick || "你" }} 的提交</div>
              <div>@{{ submission_result.__time | moment("llll") }} ({{ submission_result.__time | moment("from", "now") }})</div>
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
                <el-table-column label="提交结果">
                  <template slot-scope="scope">
                    <el-tag
                      size="small"
                      type="success"
                      v-if="submission_result[scope.row.name]"
                    >
                      <i class="el-icon-success"></i> {{ submission_result[scope.row.name].fullname }}
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

              <el-form
                :inline="true"
                :model="formInline"
                class="demo-form-inline"
              >
                <el-form-item label="分数">
                  <el-input v-model="submissionScore" placeholder="未批改"></el-input>
                </el-form-item>
                <el-form-item>
                  <el-button
                    type="primary"
                    @click="setHomeworScore"
                  >提交</el-button>
                </el-form-item>
              </el-form>
            </el-col>
            <el-col :span="16">
              <div v-if="loading">
                正在获取信息……
              </div>
              <div v-else-if="selHomework === null">
                请在左侧选择一个作业以显示
              </div>
              <div v-else-if="selHomework === undefined">
                未提交该作业
              </div>
              <div v-else-if="selHomework.extname == ('.pdf')">
                <div>
                  <div v-if="pdfPageCount >= 15">
                    <el-input-number
                      v-model="pdfSelectedPage"
                      :min="1"
                      :max="pdfPageCount"
                    ></el-input-number>
                  </div>
                  <el-slider
                    v-else
                    v-model="pdfSelectedPage"
                    :min="1"
                    :max="pdfPageCount"
                    :step="1"
                    show-stops
                  >
                  </el-slider>
                </div>
                {{ '第' + pdfCurrentPage + '/' + pdfPageCount + '页'}}
                <pdf
                  :src="'/rpc/attachments/downloadHomework/'+displayToken+'/'+selHomework.fullname"
                  @num-pages="pdfPageCount = $event"
                  @page-loaded="pdfCurrentPage = $event"
                  :page="pdfSelectedPage"
                >
                  <template slot="loading">
                    loading content here...
                  </template>
                </pdf>
                <div>
                  <el-input-number
                    v-if="pdfPageCount >= 15"
                    v-model="pdfSelectedPage"
                    :min="1"
                    :max="pdfPageCount"
                    :label="'第' + pdfCurrentPage + '/' + pdfPageCount + '页'"
                  ></el-input-number>
                  <el-slider
                    v-else
                    v-model="pdfSelectedPage"
                    :min="1"
                    :max="pdfPageCount"
                    :step="1"
                    show-stops
                  >
                  </el-slider>
                </div>
              </div>
              <div v-else-if="selHomework.extname == '.png' || selHomework.extname == ('.jpg')  || selHomework.extname == ('.jpeg')">
                图片: {{ selHomework.uid }} / {{ selHomework.fullname }}
                display: {{ displayToken }}
              </div>
              <div v-else>
                暂不支持该文件类型，你可以单击下面的连接下载后查看
                <a :href="'/rpc/attachments/downloadHomework/'+displayToken+'/'+selHomework.fullname">链接很快失效</a>
              </div>
            </el-col>
          </el-row>
        </div>
      </div>
    </el-col>
  </el-row>
</template>

<script>
import { RPC } from "../rpc.js";
import { mapState } from "vuex";
import pdf from "vue-pdf";

export default {
  data() {
    return {
      homework: null,
      pdfPageCount: -1,
      pdfCurrentPage: -1,
      pdfSelectedPage: 1,
      attachment_list: [],
      uplfile: [],
      submission_result: {},
      displayToken: "",
      upload_head: {},
      extData: {
        name: ""
      },
      submissionScore: 0,
      submissionGetErr: false,
      selHomework: null,
      isOverdue: 0,
      loading: false,
      error: false
    };
  },
  components: {
    pdf
  },
  methods: {
    loadSingleSubmission: async function(uid, filter) {
      this.loading = true;
      let res = await RPC.doRPC("getHomeworkSubmissionKey", {
        domain: this.$route.params.domain,
        uid,
        filter
      });
      this.loading = false;
      if (!res.success) {
        this.submissionGetErr = true;
      }
      this.pdfSelectedPage = 1;
      this.pdfPageCount = -1;
      this.pdfCurrentPage = -1;
      this.displayToken = res.uid;
    },
    loadSubmission: async function() {
      this.loading = true;
      let res = await RPC.doRPC("getHomeworkSubmission", {
        domain: this.$route.params.domain,
        uid: this.$route.params.uid,
        filter: this.$route.params.user
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
      this.submissionScore = res.homeworksubmission.score;
      if (Number((this.submissionScore = res.homeworksubmission.score)) < 0) {
        this.submissionScore = "";
      }
      subres.__nick = res.homeworksubmission.nick;
      subres.__time = res.homeworksubmission.create_time;
      this.submission_result = subres;
      console.log(res);
    },
    setHomeworScore: async function() {
      this.loading = true;
      let res = await RPC.doRPC("setHomeworkScore", {
        domain: this.$route.params.domain,
        uid: this.$route.params.user,
        filter: this.submissionScore
      });
      this.loadHomework();
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
      console.log(this.submission_result[e.name]);
      this.selHomework = this.submission_result[e.name];
      this.loadSingleSubmission(
        this.$route.params.uid,
        this.selHomework.uid + "," + this.$route.params.user
      );
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


