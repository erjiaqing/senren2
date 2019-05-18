<template>
  <el-row>
    <el-col :span="24">
      <div class="grid-content submission-title">
        <span v-if="submission">
          提交 #<span style="font-family: monospace">{{ submission.uid }}</span>
        </span>
        <span v-else-if="loading">Loading...</span>
        <span v-else-if="error">提交加载失败</span>
      </div>
    </el-col>
    <el-col :span="24">
      <div v-if="error">
        <el-alert
          title="请求失败"
          type="error"
          description="可能的原因：服务器故障、网络问题或提交不存在"
          show-icon
        >
        </el-alert>
      </div>
      <div
        class="grid-content submission-content"
        v-if="submission"
      >
        <div id="baseinfo_container">
          <el-form
            label-position="left"
            inline
            class="submission_metainfo_container"
          >
            <el-form-item
              class="submission_metainfo_item"
              label="评测ID"
            >
              #<span style="font-family: monospace">{{ submission.uid }}</span>
            </el-form-item>
            <el-form-item
              class="submission_metainfo_item"
              label="评测状态"
            >
              <span v-if="tags['JUDGE_TAG_' + submission.verdict]">
                <el-tag :type="tags['JUDGE_TAG_' + submission.verdict][0]">
                  <img
                    v-if="submission.status == 'PENDING'"
                    src="@/assets/loading.gif"
                  />
                  {{ tags['JUDGE_TAG_' + submission.verdict][1] }}</el-tag>
              </span>
            </el-form-item>
            <el-form-item
              class="submission_metainfo_item"
              label="试题"
            >
              #<span style="font-family: monospace">{{ submission.problem_uid }}</span>
            </el-form-item>
            <el-form-item
              class="submission_metainfo_item"
              label="用户"
            >
              #<span style="font-family: monospace">{{ submission.user_uid }}</span>
            </el-form-item>
            <el-form-item
              class="submission_metainfo_item"
              label="使用时间"
            >
              {{ (submission.execute_time <
                0)
                ? '-'
                :
                ((submission.execute_time)
                + " ms"
                )
                }}
                </el-form-item>
                <el-form-item
                class="submission_metainfo_item"
                label="使用内存"
              >
                {{ (submission.execute_memory <
                  0)
                  ? '-'
                  :
                  ((submission.execute_memory)
                  + " KiB"
                  )
                  }}
                  </el-form-item>
                  <el-form-item
                  class="submission_metainfo_item"
                  label="提交时间"
                >
                  <el-tooltip
                    class="item"
                    effect="dark"
                    placement="top-start"
                  >
                    <el-button type="text">{{ submission.submit_time | moment("from", "now") }}</el-button>
                    <template slot="content">{{ submission.submit_time | moment("llll") }}</template>
                  </el-tooltip>
            </el-form-item>
            <el-form-item
              class="submission_metainfo_item"
              label="评测时间"
            >
              <el-tooltip
                class="item"
                effect="dark"
                placement="top-start"
              >
                <el-button type="text">{{ submission.judge_time | moment("from", "now") }}</el-button>
                <template slot="content">{{ submission.judge_time | moment("llll") }}</template>
              </el-tooltip>
            </el-form-item>
          </el-form>
        </div>
        <editor
          v-model="code"
          @init="editorInit"
          :options="{readOnly: true}"
          :lang="codeHighlight"
          theme="chrome"
          width="100%"
          height="300px"
        ></editor>
        <h3 v-if="submission.ce_message != ''">编译器输出</h3>
        <pre
          style="font-size: 12pt"
          v-if="submission.ce_message != ''"
        >{{ submission.ce_message }}</pre>
        <div v-if="judgerResponse">
          <h3>评测详情</h3>
          <el-steps
            direction="vertical"
            :active="3"
          >
            <el-step
              :title="d.name"
              :status="d.verdict == 'AC' ? 'success' : (d.verdict == 'IG' ? 'wait' : 'error')"
              :description="d.verdict != 'IG' ? (d.verdict + ' / ' + d.exe_time + ' s / ' + d.exe_memory + ' KiB') : 'Not judged' "
              :icon="d.icon"
              v-for="d in judgerResponse.detail"
              :key="d.name"
            ></el-step>
          </el-steps>
        </div>
      </div>
    </el-col>
  </el-row>
</template>

<script>
import { RPC } from "../rpc.js";
import { ConstString } from "../consts.js";

export default {
  data() {
    return {
      submission: null,
      loading: false,
      error: false,
      codeEditor: false,
      code: "",
      tags: ConstString,
      selectedLanguage: "",
      codeHighlight: "",
      languages: {},
      judgerResponse: {},
      loadingInterval: -1,
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
  methods: {
    loadSubmission: async function() {
      if (this.submission && this.submission.status != "PENDING") {
        clearInterval(this.loadingInterval);
        return;
      }
      this.loading = true;
      let res = await RPC.doRPC("getSubmission", {
        domain: this.$route.params.domain,
        uid: this.$route.params.uid
      });
      this.loading = false;
      if (res == null) {
        this.error = true;
        return;
      }
      this.submission = res.submission;
      this.langChange(res.submission.language);
      this.code = this.submission.code;
      let ti = JSON.parse(this.submission.judger_response);
      let compilePos = 0;
      if (ti && ti.detail) {
        for (let i = 0; i < ti.detail.length; i++) {
          if (ti.detail[i].name == "compile") {
            compilePos = i;
          } else {
            if (ti.detail[i].verdict == 'AC') {
              ti.detail[i].icon = "el-icon-success";
            } else if (ti.detail[i].verdict == 'IG') {
              ti.detail[i].icon = "el-icon-remove";
            } else {
              ti.detail[i].icon = "el-icon-error";
            }
          }
        }
        let compileDat = ti.detail[compilePos];
        compileDat.exe_memory /= 1024;
        compileDat.name = "编译";
        compileDat.icon = "el-icon-cpu";
        ti.detail.splice(compilePos, 1);
        ti.detail.unshift(compileDat);
      }
      this.judgerResponse = ti;
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
    this.loadSubmission();
    this.loadingInterval = setInterval(this.loadSubmission, 1000);
  },
  beforeDestroy() {
    clearInterval(this.loadingInterval);
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


