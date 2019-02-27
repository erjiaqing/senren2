<template>
  <div id="app">
    <el-container>
      <el-header
        id="mainHeader"
        :class="headerShadow ? 'scolled' : ''"
      >
        <div id="main-nav-btns">
          <el-row :gutter="20">
            <el-col
              :xs="16"
              :sm="16"
              :md="6"
            >
              <div
                class="grid-content"
                style="line-height:60px;text-align:left"
              >Problem CI</div>
            </el-col>

            <el-col
              :xs="8"
              :sm="8"
              :md="6"
              style="float:right"
            >

              <div
                class="grid-content"
                style="line-height:60px;text-align:right;"
                v-if="user"
              >
                <el-button>{{ user.username }}</el-button>
              </div>
              <div
                class="grid-content"
                style="line-height:60px;text-align:right"
                v-else
              >
                <el-button @click="$router.push('/login')">登陆</el-button>
              </div>
            </el-col>

            <el-col
              :xs="24"
              :sm="24"
              :md="12"
            >
              <div
                class="grid-content"
                style="text-align: left; line-height:60px"
              >
                <el-button
                  type="text"
                  @click="$router.push('/')"
                >主页</el-button>
                <el-button
                  type="text"
                  @click="$router.push('/problems')"
                >创建的试题</el-button>
              </div>
            </el-col>
          </el-row>
        </div>
      </el-header>

      <el-main id="main-container-parent">
        <div id="main-container">
          <router-view></router-view>
        </div>
      </el-main>

      <el-footer id="page-footer">&copy; 2018 - {{(new Date()).getFullYear()}} 武汉大学 ACM 集训队</el-footer>
    </el-container>
  </div>
</template>

<script>
import { RPC } from "./rpc.js";
//import { ConstString } from "./consts.js";

export default {
  name: "app",
  data: () => {
    return {
      headerShadow: false,
      user: null
    };
  },
  methods: {
    handleScroll(scoll) {
      this.headerShadow = window.scrollY > 1;
    },
    reloadWhoAmI: async function() {
      await RPC.refreshSession();
      let res = await RPC.doRPC("isLogin", {});
      if (res == null || !res.success) {
        this.user = null;
      }
      this.user = res.user;
      if (!this.user) {
        if (this.$route.name != "login") {
          this.$router.push("/login");
        }
        return;
      }
    }
  },
  created() {
    window.addEventListener("scroll", this.handleScroll);
    this.reloadWhoAmI();
    this.$bus.$on("loginChange", () => {
      this.reloadWhoAmI();
    });
  },
  watch: {
    $route: function() {}
  },
  destroyed() {
    window.removeEventListener("scroll", this.handleScroll);
  }
};
</script>

<style>
html {
  height: 100%;
}

body {
  margin: 0;
  min-height: 100%;
  padding: 0;
  position: relative;
}

#page-footer {
  position: absolute;
  bottom: 0;
  width: 100%;
  height: 100px;
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

  z-index: 1000;
  position: fixed;
  background-color: #fff;
  top: 0;
  left: 0;
  width: 100%;
}

#main-container-parent {
  padding-bottom: 112px;
}

.scolled {
  -moz-box-shadow: 0px 2px 5px #999;
  -webkit-box-shadow: 0px 2px 5px #999;
  box-shadow: 0px 2px 5px #999;
}

@media (max-width: 768px) {
  #main-container,
  #main-nav-btns {
    width: 100%;
  }

  #mainHeader {
    height: 120px !important;
    -moz-box-shadow: 0px 2px 5px #999;
    -webkit-box-shadow: 0px 2px 5px #999;
    box-shadow: 0px 2px 5px #999;
  }

  #main-container-parent {
    padding-top: 132px;
  }
}

/** in 2019, who use 1366 x 768 devices? hun? */
@media (min-width: 768px) and (max-width: 992px) {
  #main-container,
  #main-nav-btns {
    max-width: 750px;
    margin: auto;
  }
  #mainHeader {
    height: 120px !important;
    -moz-box-shadow: 0px 2px 5px #999;
    -webkit-box-shadow: 0px 2px 5px #999;
    box-shadow: 0px 2px 5px #999;
  }

  #main-container-parent {
    padding-top: 132px;
  }
}

@media (min-width: 992px) and (max-width: 1200px) {
  #main-container,
  #main-nav-btns {
    max-width: 970px;
    margin: auto;
  }
  #mainHeader {
    height: 60px;
  }

  #main-container-parent {
    padding-top: 72px;
  }
}

@media (min-width: 1200px) {
  #main-container,
  #main-nav-btns {
    max-width: 1170px;
    margin: auto;
  }
  #mainHeader {
    height: 60px;
  }

  #main-container-parent {
    padding-top: 72px;
  }
}

pre {
  white-space: pre-wrap; /* css-3 */
  white-space: -moz-pre-wrap; /* Mozilla, since 1999 */
  white-space: -pre-wrap; /* Opera 4-6 */
  white-space: -o-pre-wrap; /* Opera 7 */
  word-wrap: break-word; /* Internet Explorer 5.5+ */
  padding: 8px;
  border-radius: 4px;
  border-style: solid;
  border-width: 1px;
  border-color: #0c0c0c;
}

.example-io,
.example-io th,
.example-io td {
  border-width: 1px;
  border-style: solid;
  border-color: #000000;
  border-collapse: collapse;
  width: 100%;
}

.example-io td {
  padding: 4px;
}

.example-input-title,
.example-output-title {
  font-weight: bolder;
}

.example-io {
  margin-bottom: 8px;
}

.example-io pre {
  white-space: pre-wrap; /* css-3 */
  white-space: -moz-pre-wrap; /* Mozilla, since 1999 */
  white-space: -pre-wrap; /* Opera 4-6 */
  white-space: -o-pre-wrap; /* Opera 7 */
  word-wrap: break-word; /* Internet Explorer 5.5+ */
  border: none;
  padding: 0;
  margin: 0;
}
</style>
