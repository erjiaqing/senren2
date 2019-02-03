<template>
  <div id="app">
    <el-container>
      <el-header id="mainHeader">
        <el-row :gutter="20">
          <el-col :span="4">
            <div class="grid-content" style="line-height:60px">WOJ</div>
            <!-- TODO 根据当前用户所在的域显示对应的名称 -->
          </el-col>
          <el-col :span="16">
            <div class="grid-content" style="text-align: left; line-height:60px">
              <el-button type="text">主页</el-button>
              <el-button type="text">题库</el-button>
              <el-button type="text">比赛</el-button>
              <el-button type="text">小组</el-button>
            </div>
          </el-col>
          <el-col :span="4">
            <div class="grid-content" style="line-height:60px">Login / Register</div>
            <!-- TODO 登陆&注册&权限管理 -->
          </el-col>
        </el-row>
      </el-header>

      <el-main>
        <router-view></router-view>
        <pre style="text-align:left">
        // show group description
        // TODO: show recent discussions and contests and important groups (domains, if in official group (0000000000000000))
        // Router Format:
        // /{group shortname} [or uid, transformed to uid by backend]
        //  /problems
        //  /problem/{problem-id}
        //   /edit [admin]
        //   /fork [admin]
        //  /submissions
        //  /submission/{submission-id}
        //  /contests
        //  /contest/{contest-id}
        //   /problems
        //   /problem/{problem-seq}
        //   /discussions
        //   /discussion/{topic-id}
        //   /submissions
        //   /submission/{submission-id}
        //   /rank
        //   /edit [admin]
        //  /homeworks [not available in group start with '0']
        //  /homework/{homework-id}
        //   /edit [admin]
        //  /discussions
        //  /discussion/{topic-id}
        //  /edit [admin]
        </pre>
      </el-main>
    </el-container>
  </div>
</template>

<script>
export default {
  name: "app",
  data: () => {
    return {
    };
  },
};
</script>

<style>
body {
  margin: 0;
}

#app {
  font-family: "Avenir", Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}

.headnav-btn {
  font-size: 36px;
}
</style>
