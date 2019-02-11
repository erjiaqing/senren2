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
              >{{ domain ? (domain.title || domain.alias || domain.uid) : "千练万花" }}</div>
              <!-- TODO 根据当前用户所在的域显示对应的名称 -->
            </el-col>

            <el-col
              :xs="8"
              :sm="8"
              :md="10"
              style="float:right"
            >
              <div
                class="grid-content"
                style="line-height:60px;text-align:right;"
                v-if="user"
              >
                <el-dropdown
                  trigger="click"
                  @command="handleCommand"
                >
                  <span class="el-dropdown-link">
                    {{ user.nickname }}
                    <el-tag
                      size="mini"
                      :type="tags['ROLE_TAG_' + user.role][0]"
                    >{{ tags['ROLE_TAG_' + user.role][1] }}</el-tag>
                    <i class="el-icon-arrow-down el-icon--right"></i>
                  </span>
                  <el-dropdown-menu slot="dropdown">
                    <el-dropdown-item command="gotoMySubmissions">提交列表</el-dropdown-item>
                    <el-dropdown-item
                      v-if="guser"
                      disabled
                    >当前用户 {{ guser.nickname }} @ WOJ (
                      <el-tag
                        size="mini"
                        :type="tags['ROLE_TAG_' + guser.role][0]"
                      >{{ tags['ROLE_TAG_' + guser.role][1] }}</el-tag>)</el-dropdown-item>
                    <el-dropdown-item v-if="user && user.role == 'NONE'">加入 {{ domain.title }}</el-dropdown-item>
                    <el-dropdown-item
                      v-if="user && user.role != 'NONE'"
                      disabled
                    >已加入 {{ domain.title }}</el-dropdown-item>
                    <el-dropdown-item
                      command="logout"
                      divided
                    >登出</el-dropdown-item>
                    <el-dropdown-item
                      v-if="user.role == 'ROOT' || user.role == 'ADMIN' || user.role == 'VIP_ROOT' || user.role == 'VIP_ADMIN'"
                      command="editDomain"
                      divided
                    >编辑小组
                      <el-tag
                        size="mini"
                        :type="tags['ROLE_TAG_' + user.role][0]"
                      >管理员</el-tag>

                    </el-dropdown-item>
                  </el-dropdown-menu>
                </el-dropdown>
              </div>
              <div
                class="grid-content"
                style="line-height:60px;text-align:right"
                v-else
              >
                <el-button @click="$router.push('/' + ($route.params.domain ? $route.params.domain : 'woj') + '/login')">Login / Register</el-button>
              </div>
              <!-- TODO 登陆&注册&权限管理 -->
            </el-col>

            <el-col
              :xs="24"
              :sm="24"
              :md="8"
            >
              <div
                class="grid-content"
                style="text-align: left; line-height:60px"
              >
                <el-button
                  type="text"
                  @click="$router.push('/' + ($route.params.domain ? $route.params.domain : ''))"
                >主页</el-button>
                <el-button
                  type="text"
                  @click="$router.push('/' + ($route.params.domain ? $route.params.domain : 'woj') + '/problems')"
                >题库</el-button>
                <el-button
                  type="text"
                  @click="$router.push('/' + ($route.params.domain ? $route.params.domain : 'woj') + '/contests')"
                >比赛</el-button>
                <el-button
                  type="text"
                  @click="$router.push('/' + ($route.params.domain ? $route.params.domain : 'woj') + '/domains')"
                >小组</el-button>
                <el-button type="text" v-if="guser && user && user.role == 'NONE'">加入 {{ domain.title }}</el-button>
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

      <el-footer>&copy; 2018 - {{(new Date()).getFullYear()}} 武汉大学 ACM 集训队</el-footer>
    </el-container>
  </div>
</template>

<script>
import { RPC } from "./rpc.js";
import { ConstString } from "./consts.js";

export default {
  name: "app",
  data: () => {
    return {
      headerShadow: false,
      user: null,
      guser: null,
      tags: ConstString,
      currentDomain: "0000000000000000",
      currentDomainShort: "WOJ",
      domain: null
    };
  },
  methods: {
    handleScroll(scoll) {
      this.headerShadow = window.scrollY > 1;
    },
    reloadWhoAmI: async function() {
      let res = await RPC.doRPC("whoami", {
        domain: this.$route.params.domain
      });
      if (res == null || !res.success) {
        this.user = null;
        return;
      }
      this.user = res.user;
      this.guser = res.user_global;
    },
    reloadDomain: async function() {
      let res = await RPC.doRPC("getDomain", {
        filter: "nodesc",
        domain: this.$route.params.domain
      });
      if (res == null || !res.success) {
        this.user = null;
        return;
      }
      this.domain = res.domain;
      this.currentDomain = res.domain.uid;
      this.currentDomainShort = res.domain.alias;
      this.$store.commit("setDomain", res.domain);
    },
    handleCommand(cmd) {
      switch (cmd) {
        case "gotoMySubmissions": {
          this.$router.push(
            "/" +
              (this.$route.params.domain || "woj") +
              "/submissions/;" +
              this.user.uid
          );
          break;
        }
        case "logout": {
          RPC.setSession("");
          this.$bus.$emit("loginChange");
          break;
        }
        case "editDomain": {
          this.$router.push(
            "/" + (this.$route.params.domain || "woj") + "/edit"
          );
          break;
        }
      }
    }
  },
  created() {
    window.addEventListener("scroll", this.handleScroll);
    this.reloadWhoAmI();
    this.reloadDomain();
    this.$bus.$on("loginChange", () => {
      this.reloadWhoAmI();
    });
  },
  watch: {
    $route: function() {
      if (
        this.$route.params.domain != this.currentDomain &&
        this.$route.params.domain != this.currentDomainShort
      ) {
        this.reloadDomain();
        this.reloadWhoAmI();
      }
    }
  },
  destroyed() {
    window.removeEventListener("scroll", this.handleScroll);
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

  z-index: 1000;
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
</style>
