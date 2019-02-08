<template>
  <div id="app">
    <el-container>
      <el-header id="mainHeader" :class="headerShadow ? 'scolled' : ''">
        <div id="main-nav-btns">
          <el-row :gutter="20">
            <el-col :span="4">
              <div
                class="grid-content"
                style="line-height:60px;text-align:left"
              >WOJ</div>
              <!-- TODO 根据当前用户所在的域显示对应的名称 -->
            </el-col>
            <el-col :span="16">
              <div
                class="grid-content"
                style="text-align: left; line-height:60px"
              >
                <el-button type="text" @click="$router.push('/' + ($route.params.domain ? $route.params.domain : ''))">主页</el-button>
                <el-button type="text" @click="$router.push('/' + ($route.params.domain ? $route.params.domain : 'woj') + '/problems')">题库</el-button>
                <el-button type="text" @click="$router.push('/' + ($route.params.domain ? $route.params.domain : 'woj') + '/contests')">比赛</el-button>
                <el-button type="text">小组</el-button>
              </div>
            </el-col>
            <el-col :span="4">
              <div
                class="grid-content"
                style="line-height:60px;text-align:right"
              >Login / Register</div>
              <!-- TODO 登陆&注册&权限管理 -->
            </el-col>
          </el-row>
        </div>
      </el-header>

      <el-main id="main-container-parent">
        <div id="main-container">
          <router-view></router-view>
        </div>
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
      headerShadow: false,
    };
  },
  methods: {
    handleScroll(scoll) {
      this.headerShadow = window.scrollY > 1;
    },
  },
  created () {
    window.addEventListener('scroll', this.handleScroll);
  },
  destroyed () {
    window.removeEventListener('scroll', this.handleScroll);
  }
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

#mainHeader {
  text-align: center;

  transition: all 0.3s linear;

  z-index: 99999;
  position: fixed;
  background-color: #fff;
  top: 0;
  left: 0;
  width: 100%;
}

.scolled {
  -moz-box-shadow: 0px 2px 5px #999;
  -webkit-box-shadow: 0px 2px 5px #999;
  box-shadow: 0px 2px 5px #999;
}

#main-container-parent {
  padding-top: 72px;
}

@media (max-width: 768px) {
  #main-container,
  #main-nav-btns {
    width: 100%;
  }
}

/** in 2019, who use 1366 x 768 devices? hun? */
@media (min-width: 768px) and (max-width: 992px) {
  #main-container,
  #main-nav-btns {
    max-width: 750px;
    margin: auto;
  }
}

@media (min-width: 992px) and (max-width: 1200px) {
  #main-container,
  #main-nav-btns {
    max-width: 970px;
    margin: auto;
  }
}

@media (min-width: 1200px) {
  #main-container,
  #main-nav-btns {
    max-width: 1170px;
    margin: auto;
  }
}
</style>
