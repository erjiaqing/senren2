<template>
  <el-row>
    <el-col :span="24">
      <div class="grid-content problem-title">
        <span v-if="problem">{{ problem.alias == "" ? "" : problem.alias + '. '}}{{ problem.title }}</span>
        <span v-else-if="loading">Loading...</span>
        <span v-else-if="error">「试题加载失败」</span>
      </div>
    </el-col>
    <el-col :span="24">
      <div class="grid-content problem-title">
        <el-button-group>
          <el-button
            icon="el-icon-edit-outline"
            @click="codeEditor = !codeEditor"
            v-if="user"
          >{{ codeEditor ? '收起编辑器' : '代码编辑器' }}</el-button>
          <el-button
            icon="el-icon-more"
            @click="$router.push('/' + $route.params.domain + '/submissions/' + problem.uid)"
          >评测结果</el-button>
          <!-- <el-button icon="el-icon-tickets">讨论区</el-button> -->
          <el-button
            icon="el-icon-edit"
            @click="gotoEditor"
            v-if="user && (user.role == 'ADMIN' || user.role == 'ROOT')"
          >编辑试题</el-button>
          <!-- <el-button
            icon="el-icon-share"
            v-if="user && (user.role == 'ADMIN' || user.role == 'ROOT')"
          >克隆试题</el-button> -->
        </el-button-group>
      </div>
    </el-col>
    <el-col :span="24">
      <div v-if="error">
        <el-alert
          title="请求失败"
          type="error"
          description="可能的原因：服务器故障、网络问题、试题不存在或未登陆"
          show-icon
        >
        </el-alert>
      </div>
      <div
        class="grid-content problem-content"
        v-if="problem"
      >
        <div
          v-if="codeEditor"
          class="codeEditorPanel"
        >
          <el-select
            v-model="selectedLanguage"
            @change="langChange"
            placeholder="请选择语言"
          >
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
            <el-button
              icon="el-icon-upload2"
              @click="loadCode()"
            >恢复</el-button>
            <el-button
              icon="el-icon-download"
              @click="saveCode()"
            >暂存</el-button>
            <el-button
              type="success"
              icon="el-icon-upload"
              @click="submitCode"
              :disabled="selectedLanguage == ''"
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
  computed: mapState(["user"]),
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
      if (res.problem.alias && this.$route.params.uid != res.problem.alias) {
        this.$router.replace(
          `/${this.$route.params.domain}/problem/${res.problem.alias}`
        );
      }
      this.problem = res.problem;
      this.$nextTick(function() {
        MathJax.Hub.Queue(["Typeset", MathJax.Hub]);
        Han(this.$refs.problemContent).render()
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
    gotoEditor: function() {
      this.$router.push(
        `/${this.$route.params.domain}/problem/${this.$route.params.uid}/edit`
      );
    },
    submitCode: async function() {
      const loading = this.$loading({
        lock: true,
        text: "提交中",
        spinner: "el-icon-loading",
        background: "rgba(0, 0, 0, 0.7)"
      });
      this.error = false;

      let res = await RPC.doRPC("createSubmission", {
        submission: {
          problem_uid: this.problem.uid,
          domain: this.problem.domain,
          contest_uid: "0000000000000000",
          language: this.selectedLanguage,
          code: this.code
        }
      });

      let success = true;

      if (res == null || res.success != true) {
        success = false;
        this.error = true;
      }

      console.log(res);
      loading.close();
      if (success) {
        this.$router.push(
          `/${this.$route.params.domain}/submission/${res.uid}`
        );
      }
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
</style>


