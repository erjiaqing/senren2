<template>
  <div>
    <!-- <div class="grid-content problem-title">
      <span v-if="contest">{{ contest.title }}</span>
      <span v-else-if="loading">Loading...</span>
      <span v-else-if="error">「比赛加载失败」</span>
    </div>
    <div
      class="grid-content"
      v-if="contest"
    >
      <span>比赛结束于 {{ contest.end_time | moment("from", "now") }}</span>
    </div> -->
    <div v-if="error">
      <el-alert
        title="请求失败"
        type="error"
        description="可能的原因：服务器故障、网络问题或比赛不存在"
        show-icon
      >
      </el-alert>
    </div>
    <router-view></router-view>
    <div
      id="contest_desc"
      v-if="contest && contest.description && $route.name == 'contest_index'"
      style="text-align: left;margin-top: 20px"
    >
      <div v-html="contest.description"></div>
    </div>
  </div>
</template>

<script>
import { RPC } from "../rpc.js";
import { mapState } from "vuex";

export default {
  data() {
    return {
      probIndex: "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
      contest: null,
      loading: false,
      error: false
    };
  },
  methods: {
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
      res.contest.problem_list = JSON.parse(res.contest.problem_list);
      for (let i = 0; i < res.contest.problem_list.length; i++) {
        res.contest.problem_list[i].$index = i;
        res.contest.problem_list[i].$charid = this.probIndex.charAt(i);
      }
      res.contest.start_time = new Date(res.contest.start_time);
      res.contest.end_time = new Date(res.contest.end_time);
      res.contest.freeze_time = new Date(res.contest.freeze_time);
      this.$store.commit("setContest", res.contest);
      this.contest = res.contest;
    },
    gotoEditor: function() {
      this.$router.push(
        `/${this.$route.params.domain}/contest/${this.$route.params.uid}/edit`
      );
    },
    gotoProblem: function(e) {
      this.$router.push(
        `/${this.$route.params.domain}/contest/${
          this.$route.params.uid
        }/problem/${e.$charid}`
      );
    }
  },
  computed: mapState(["user"]),
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
.el-form-item.contest_metadata_item {
  margin-bottom: 8px;
  font-size: 12pt;
  line-height: 16pt;
  width: 100%;
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

#contest_desc {
  padding: 8px;
  margin-top: 8px;
  margin-bottom: 8px;
  border-radius: 4px;
  -moz-box-shadow: 0px 0px 5px #999;
  -webkit-box-shadow: 0px 0px 5px #999;
  box-shadow: 0px 0px 5px #999;
}
</style>


