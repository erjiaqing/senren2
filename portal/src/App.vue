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
              :xs="12"
              :sm="16"
              :md="6"
            >
              <div
                class="grid-content"
                style="line-height:60px;text-align:left"
                v-if="$route.name && $route.name.startsWith('contest_')"
              >{{ contest ? contest.title : "千练万花 - 比赛" }}</div>
              <div
                class="grid-content"
                style="line-height:60px;text-align:left"
                v-else
              >{{ domain ? (domain.title || domain.alias || domain.uid) : "千练万花" }}</div>
              <!-- TODO 根据当前用户所在的域显示对应的名称 -->
            </el-col>

            <el-col
              :xs="12"
              :sm="8"
              :md="6"
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
                    >WOJ用户 {{ guser.nickname }}
                      <el-tag
                        size="mini"
                        :type="tags['ROLE_TAG_' + guser.role][0]"
                      >{{ tags['ROLE_TAG_' + guser.role][1] }}</el-tag>
                    </el-dropdown-item>
                    <el-dropdown-item v-if="domain && user && user.role == 'NONE'">加入 {{ domain.title }}</el-dropdown-item>
                    <el-dropdown-item
                      v-if="domain && user && user.role != 'NONE'"
                      disabled
                    >已加入 {{ domain.title }}</el-dropdown-item>
                    <el-dropdown-item
                      command="logout"
                      divided
                    >登出</el-dropdown-item>
                    <el-dropdown-item
                      v-if="user.role == 'ROOT' || user.role == 'ADMIN' || user.role == 'VIP_ROOT' || user.role == 'VIP_ADMIN'"
                      divided
                      disabled
                    >
                      <el-tag
                        size="mini"
                        :type="tags['ROLE_TAG_' + user.role][0]"
                      >管理员</el-tag>
                    </el-dropdown-item>
                    <el-dropdown-item
                      v-if="user.role == 'ROOT' || user.role == 'ADMIN' || user.role == 'VIP_ROOT' || user.role == 'VIP_ADMIN'"
                      command="editDomain"
                    >编辑小组信息
                    </el-dropdown-item>
                    <el-dropdown-item
                      v-if="user.role == 'ROOT' || user.role == 'ADMIN' || user.role == 'VIP_ROOT' || user.role == 'VIP_ADMIN'"
                      command="editInvite"
                    >管理邀请码
                    </el-dropdown-item>
                    <el-dropdown-item
                      v-if="user.role == 'ROOT' || user.role == 'ADMIN' || user.role == 'VIP_ROOT' || user.role == 'VIP_ADMIN'"
                      command="editInvite"
                    >管理小组成员
                    </el-dropdown-item>
                  </el-dropdown-menu>
                </el-dropdown>
              </div>
              <div
                class="grid-content"
                style="line-height:60px;text-align:right"
                v-else
              >
                <el-button @click="$router.push('/' + ($route.params.domain ? $route.params.domain : 'woj') + '/login')">登陆 / 注册 </el-button>
              </div>
              <!-- TODO 登陆&注册&权限管理 -->
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
                <span v-if="$route.name && $route.name.startsWith('contest_')">
                  <el-button
                    type="text"
                    @click="$router.push('/' + ($route.params.domain ? $route.params.domain : ''))"
                  >返回首页</el-button>
                  <el-button
                    type="text"
                    @click="$router.push('/' + $route.params.domain + '/contest/' + $route.params.uid)"
                  >试题列表</el-button>
                  <el-button
                    type="text"
                    @click="$router.push('/' + $route.params.domain + '/contest/' + $route.params.uid + '/submissions/;' + user.uid)"
                  >提交状态</el-button>
                  <el-button
                    type="text"
                    @click="$router.push('/' + $route.params.domain + '/contest/' + $route.params.uid + '/rank')"
                  >比赛榜单</el-button>
                  <el-button
                    type="text"
                  >讨论区</el-button>
                  <el-button
                    @click="$router.push('/' + $route.params.domain + '/contest/' + $route.params.uid + '/edit')"
                    type="text"
                    v-if="user && (user.role == 'ADMIN' || user.role == 'ROOT')"
                  >编辑比赛</el-button>
                </span>
                <span v-else>
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
                    v-if="$route.params.domain != 'woj' && $route.params.domain != '0000000000000000'"
                    @click="$router.push('/' + ($route.params.domain ? $route.params.domain : 'woj') + '/homeworks')"
                  >作业</el-button>
                  <el-button
                    type="text"
                    v-if="$route.params.domain == 'woj' || $route.params.domain == '0000000000000000'"
                    @click="$router.push('/' + ($route.params.domain ? $route.params.domain : 'woj') + '/domains')"
                  >小组</el-button>
                  <el-button
                    type="text"
                    v-if="$route.params.domain != 'woj' && $route.params.domain != '0000000000000000'"
                    @click="$router.push('/woj')"
                  >返回WOJ</el-button>
                  <el-button
                    type="text"
                    v-if="$route.params.domain != 'woj' && $route.params.domain != '0000000000000000' && guser && user && user.role == 'NONE'"
                  >加入 {{ domain.title }}</el-button>
                  <el-button
                    type="text"
                    v-if="$route.params.domain != 'woj' && $route.params.domain != '0000000000000000' && guser && user && user.role != 'NONE'"
                    disabled
                  >已是小组成员</el-button>
                </span>
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
import { ConstString } from "./consts.js";
import { mapState } from "vuex";

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
        res = { user: null, user_global: null };
        //return;
      }
      this.user = res.user;
      this.guser = res.user_global;
      this.$store.commit("setUser", res.user);
      this.$store.commit("setGUser", res.user_global);
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
        case "editInvite": {
          this.$router.push(
            "/" + (this.$route.params.domain || "woj") + "/invites"
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
  computed: mapState(["contest"]),
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

pre,
code {
  font-size: 12pt;
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

i[class^="el-"][class*="el-"] {
  padding-bottom: 0;
  border-bottom: none;
}

.grid-content.problem-content h1{
  font-size: 1.4em;
}
.grid-content.problem-content h2{
  font-size: 1.3em;
}
.grid-content.problem-content h3{
  font-size: 1.2em;
}
.grid-content.problem-content h4{
  font-size: 1.1em;
}
.grid-content.problem-content h5{
  font-size: 1em;
}

.grid-content.problem-title{
  font-size: 1.4em;
}
</style>
