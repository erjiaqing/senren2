<template>
  <el-row :gutter="20">
    <el-col :span="24" :md="16" :lg="18">
      <div class="grid-content problem-title">
        <span v-if="problem">{{ problem.alias == "" ? "" : problem.alias + '. '}}{{ problem.title }}</span>
        <span v-else-if="loading">Loading...</span>
        <span v-else-if="error">「试题加载失败」</span>
      </div>
      <div v-if="error">
        <el-alert title="请求失败" type="error" description="可能的原因：服务器故障、网络问题或试题不存在" show-icon></el-alert>
      </div>
      <div class="grid-content problem-content" v-if="problem">
        <div v-if="codeEditor" class="codeEditorPanel">
          <el-select v-model="selectedLanguage" @change="langChange" placeholder="请选择语言">
            <el-option
              v-for="item in languagePool"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            >
              <span style="float: left">{{ item.label }}</span>
              <span style="float: right; color: #8492a6; font-size: 13px">{{ item.value }}</span>
            </el-option>
          </el-select>
          <el-button-group style="float:right;">
            <el-button icon="el-icon-upload2" @click="loadCode()">恢复</el-button>
            <el-button icon="el-icon-download" @click="saveCode()">暂存</el-button>
            <el-button
              type="success"
              icon="el-icon-upload"
              :disabled="selectedLanguage == ''"
              @click="submitCode"
            >提交</el-button>
          </el-button-group>
          <editor
            v-model="code"
            @init="editorInit"
            :lang="codeHighlight"
            theme="chrome"
            width="100%"
            height="300px"
          ></editor>
        </div>
        <div v-html="problem.description" ref="problemContent"></div>
      </div>
    </el-col>
    <el-col :span="24" :md="8" :lg="6">
      <div class="problem-sidebar">
        <el-card
          class="box-card"
          shadow="hover"
          :body-style="{padding: '0'}"
          style="margin-bottom: 20px"
        >
          <div class="problem-sidebar-item">
            <div v-if="contestTimer.state == 'PENDING'">
              <span>未开始</span>
              <div class="timerText" style="padding: 3px 0">{{ contestTimer.remain }}</div>
            </div>
            <div v-else-if="contestTimer.state == 'RUNNING'">
              <span>进行中</span>
              <el-progress :show-text="false" :percentage="contestTimer.ratio * 100"></el-progress>
              <div class="timerText" style="padding: 3px 0">{{ contestTimer.remain }}</div>
            </div>
            <div v-else-if="contestTimer.state == 'FROZEN'">
              <span>进行中 - 封榜</span>
              <el-progress :show-text="false" color="#E6A23C" :percentage="contestTimer.ratio * 100"></el-progress>
              <div
                class="timerText"
                style="padding: 3px 0"
              >{{ contestTimer.remain }}</div>
            </div>
            <div v-else-if="contestTimer.state == 'END'">
              <span>已结束</span>
            </div>
          </div>
        </el-card>

        <el-card
          class="box-card"
          shadow="hover"
          :body-style="{padding: '0'}"
          style="margin-bottom: 20px"
        >
          <div class="problem-sidebar-item">
            <div
              v-if="(new Date()).valueOf() < (new Date(contest.end_time)).valueOf()"
              @click="codeEditor = !codeEditor"
            >
              <span>{{ codeEditor ? '收起编辑器' : '代码编辑器' }}</span>
              <span style="float: right; padding: 3px 0" class="el-icon-edit-outline"></span>
            </div>
            <div
              v-else-if="(new Date()).valueOf() > (new Date(contest.end_time)).valueOf() && (new Date()).valueOf() > (new Date(problem.release)).valueOf()"
              @click="$router.push('/' + $route.params.domain + '/problem/' + problem.uid)"
            >
              <span>在题库中提交</span>
              <span style="float: right; padding: 3px 0" class="el-icon-document"></span>
            </div>
            <div v-else>
              <span>暂时无法提交本题</span>
            </div>
          </div>
          <div
            class="problem-sidebar-item"
            @click="$router.push('/' + $route.params.domain + '/contest/' + $route.params.uid + '/submissions/' + problem.uid)"
          >
            <span>评测结果</span>
            <span style="float: right; padding: 3px 0" class="el-icon-right"></span>
          </div>
        </el-card>

        <el-card
          class="box-card"
          shadow="hover"
          :body-style="{padding: '0'}"
          style="margin-bottom: 20px"
        >
          <div @click="gotoProblem">
            <div
              class="problem-sidebar-item"
              v-for="p in contest.problem_list"
              :key="p.uid"
              :data-click-url="`/${$route.params.domain}/contest/${$route.params.uid}/problem/${p.$charid}`"
            >
              <span class="problem-charid">{{ p.$charid }}.</span>
              {{ p.title }}
            </div>
          </div>
        </el-card>
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
      loading: false,
      error: false,
      codeEditor: false,
      code: "",
      selectedLanguage: "",
      codeHighlight: "",
      languages: {},
      languagePool: [
        {
          label: "C (C99 with GCC)",
          value: "c.gcc99",
          compile: "gcc main.c -O2 -lm -o main",
          highlight: "c_cpp"
        },
        {
          label: "C++ (C++98 with G++)",
          value: "cpp.gxx98",
          compile: "g++ main.cpp -O2 --std=c++98 -o main",
          highlight: "c_cpp"
        },
        {
          label: "C++ (C++11 with G++)",
          value: "cpp.gxx11",
          compile: "g++ main.c -O2 --std=c++11 -o main",
          highlight: "c_cpp"
        },
        {
          label: "Java 1.8 (OpenJDK)",
          value: "java.java18",
          compile: "javac -cp .;* {MainClass}.java",
          highlight: "java"
        },
        {
          label: "Python (CPython 3.6)",
          value: "py.py36",
          compile: "python3 main.py",
          highlight: "python"
        },
        {
          label: "C# (Linux, Mono)",
          value: "cs.mono",
          compile: "mcs Program.cs",
          highlight: "csharp"
        },
        {
          label: "Kotlin",
          value: "kotlin.default",
          compile: "kotlinc main.kt -include-runtime -d main.jar",
          highlight: "kotlin"
        },
        {
          label: "Haskell (ghc 7)",
          value: "hs.ghc7",
          compile: "ghc main.hs -o main.exe",
          highlight: "haskell"
        },
        {
          label: "PHP 7",
          value: "php.php7",
          compile: "php main.php",
          highlight: "php"
        }
      ]
    };
  },
  components: {
    editor: require("vue2-ace-editor")
  },
  computed: mapState(["user", "contest", "contestTimer"]),
  methods: {
    showTime: function(p) {
      return this.contestTimer.remain;
    },
    gotoProblem: function(e) {
      this.$router.push(e.target.dataset.clickUrl);
    },
    loadProblem: async function() {
      this.loading = true;
      let res = await RPC.doRPC("getContestProblem", {
        domain: this.$route.params.domain,
        uid: this.$route.params.uid,
        filter: this.$route.params.seq
      });
      this.loading = false;
      if (res == null) {
        this.error = true;
        return;
      }
      this.problem = res.problem;
      this.$nextTick(function() {
        MathJax.Hub.Queue(["Typeset", MathJax.Hub]);
        Han(this.$refs.problemContent).render();
        console.log("Math and Han Format Rendered");
      });
    },
    saveCode: function() {
      localStorage.setItem(
        `code-${this.problem.uid}`,
        JSON.stringify({ lang: this.selectedLanguage, code: this.code })
      );
      const h = this.$createElement;
      this.$message({
        message: "代码已暂存",
        type: "success"
      });
    },
    loadCode: function() {
      let code = localStorage.getItem(`code-${this.problem.uid}`);
      if (code) {
        let c2 = JSON.parse(code);
        this.code = c2.code;
        this.selectedLanguage = c2.lang;
      }
      this.$message({
        message: "代码已加载"
      });
    },
    submitCode: async function() {
      const loading = this.$loading({
        lock: true,
        text: "提交中",
        spinner: "el-icon-loading",
        background: "rgba(0, 0, 0, 0.7)"
      });
      this.error = false;

      let res = await RPC.doRPC("createContestSubmission", {
        submission: {
          problem_uid: this.$route.params.seq,
          domain: this.problem.domain,
          contest_uid: this.$route.params.uid,
          language: this.selectedLanguage,
          code: this.code
        }
      });

      if (res == null || res.success != true) {
        this.error = true;
      }

      console.log(res);
      loading.close();
      this.$router.push(
        `/${this.$route.params.domain}/contest/${
          this.$route.params.uid
        }/submission/${res.uid}`
      );
    },
    langChange: function(newLang) {
      console.log(`Switched to ${newLang}`);
      this.languagePool.forEach(element => {
        if (element.value == newLang) {
          this.codeHighlight = element.highlight;
        }
      });
      console.log(this.codeHighlight);
    },
    editorInit() {
      require("brace/ext/language_tools"); //language extension prerequsite...
      require("brace/mode/c_cpp");
      require("brace/mode/golang"); //language
      require("brace/mode/haskell");
      require("brace/mode/java");
      require("brace/mode/kotlin");
      require("brace/mode/pascal");
      require("brace/mode/php");
      require("brace/mode/plain_text");
      require("brace/mode/python");
      require("brace/mode/csharp");
      require("brace/theme/chrome");
      require("brace/snippets/javascript"); //snippet
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

.problem-sidebar {
  text-align: left;
}

.problem-sidebar-item {
  cursor: pointer;
  padding: 18px 20px;
}

.problem-sidebar-item:hover {
  background-color: rgb(236, 245, 255);
}

.timerText {
  font-size: 24px;
  text-align: center;
}

.problem-charid {
  width: 24px;
}
</style>


