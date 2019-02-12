<template>
  <el-row>
    <el-col :span="24">
      <div class="grid-content problem-title">
        <span v-if="homework">{{ homework.title }}</span>
        <span v-else-if="loading">Loading...</span>
        <span v-else-if="error">「作业加载失败」</span>
        </span>
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
        <div id="ddl-notify">
          <i class="el-icon-info"></i>
          作业提交截止于 {{ homework.end_time | moment("from", "now") }} ({{ homework.end_time | moment('llll')}})</div>
        <div>
          <h3>作业描述</h3>
        </div>
        <div v-html="homework.description"></div>
        <div>
          <h3>提交项</h3>
          <el-table
            :data="attachment_list"
            style="width: 100%"
          >
            <el-table-column label="附件名称">
              <template slot-scope="scope">
                <span>{{ scope.row.name }}</span>
              </template>
            </el-table-column>
            <el-table-column label="文件名">
              <template slot-scope="scope">
                <code>{{ scope.row.file }}</code>
              </template>
            </el-table-column>
            <el-table-column label="附件描述">
              <template slot-scope="scope">
                <span>{{ scope.row.desc }}</span>
              </template>
            </el-table-column>
            <el-table-column label="操作">
              <template slot-scope="scope">
                <span>
                  <el-button
                    type="text"
                    :data="'submit|' + scope.row.name"
                  >提交作业</el-button>
                  <el-button
                    type="text"
                    :data="'view|' + scope.row.name"
                  >查看提交</el-button>
                </span>
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

export default {
  data() {
    return {
      homework: null,
      attachment_list: [],
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
      this.attachment_list = attachment_obj_list;
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

#ddl-notify {
  padding: 8px;
  border-radius: 4px;
  color: #f56c6c;
  background-color: #fef0f0;
}
</style>


