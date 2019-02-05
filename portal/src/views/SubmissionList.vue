<template>
  <el-table
    :data="submissions"
    style="width: 100%"
  >
    <el-table-column
      label="评测结果"
      width="200pt"
    >
      <template slot-scope="scope">
        <span v-if="tags['JUDGE_TAG_' + scope.row.verdict]">
          <el-tag :type="tags['JUDGE_TAG_' + scope.row.verdict][0]">{{ tags['JUDGE_TAG_' + scope.row.verdict][1] }}</el-tag>
        </span>
        <span v-else>
          <el-tag>Unknown {{ scope.row.verdict }}</el-tag>
        </span>
      </template>
    </el-table-column>
    <el-table-column label="试题">
      <template slot-scope="scope">
        <router-link :to="'/' + $route.params.domain + '/problem/' + scope.row.problem_uid">{{ scope.row.problem_title }}</router-link>
      </template>
    </el-table-column>
    <el-table-column label="用户">
      <template slot-scope="scope">
        <router-link :to="'/' + $route.params.domain + '/user/' + scope.row.user_uid">{{ scope.row.user_name }}</router-link>
      </template>
    </el-table-column>
    <el-table-column label="语言">
      <template slot-scope="scope">
        {{ languageAbbr[scope.row.language] }}
      </template>
    </el-table-column>
    <el-table-column label="运行时间">
      <template slot-scope="scope">
        {{ (scope.row.execute_time < 0) ? '-' : ((scope.row.execute_time) + " ms") }}
      </template>
    </el-table-column>
    <el-table-column label="使用内存">
      <template slot-scope="scope">
        {{ (scope.row.execute_memory < 0) ? '-' : ((scope.row.execute_memory) + " KiB") }}
      </template>
    </el-table-column>
    <el-table-column label="提交时间">
      <template slot-scope="scope">
        {{ scope.row.submit_time }}
      </template>
    </el-table-column>
  </el-table>
</template>

<script>
import { RPC } from "../rpc.js";
import { ConstString } from "../consts.js";

export default {
  data() {
    return {
      submissions: null,
      loading: false,
      error: false,
      codeEditor: false,
      code: "",
      tags: ConstString,
      languageAbbr: {
        "c.gcc99": "C",
        "cpp.gxx98": "C++",
        "cpp.gxx11": "C++",
        "java.java18": "Java",
        "py.py36": "Python",
        "cs.mono": "C#",
        "kotlin.default": "Kotlin",
        "hs.ghc7": "Haskell",
        "php.php7": "PHP"
      }
    };
  },
  methods: {
    loadSubmission: async function() {
      this.loading = true;
      let res = await RPC.doRPC("getSubmissions", {
        domain: this.$route.params.domain,
        filter: this.$route.params.filter || "",
      });
      this.loading = false;
      if (res == null) {
        this.error = true;
        return;
      }
      this.submissions = res.submissions;
    }
  },
  created() {
    this.loadSubmission();
  }
};
</script>

<style>
.submission-title {
  font-size: 28pt;
  text-align: left;
  padding-bottom: 12pt;
}

.submission-content {
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

.submission_metainfo_container {
  padding: 0;
}
</style>


