<template>
  <div class="home">
    <el-row :gutter="10">
      <el-col
        :xs="24"
        :md="8"
      >
        <div>&nbsp;</div>
      </el-col>
      <el-col
        :xs="24"
        :md="8"
      >
        <div class="grid-content">
          <h3>登陆 {{ domain ? (domain.title || domain.alias || domain.uid) : '千练万花' }}</h3>
          <!-- TODO: 设成当前域的名称 -->
          <el-form
            label-position="left"
            label-width="80px"
          >
            <el-form-item label="用户名">
              <el-input
                autocomplete="off"
                v-model="loginForm.username"
              ></el-input>
            </el-form-item>
            <el-form-item label="密码">
              <el-input
                type="password"
                autocomplete="off"
                v-model="loginForm.password"
              ></el-input>
            </el-form-item>
            <el-alert
              :title="'登陆失败：' + errorText"
              v-if="!loginSuccess"
              type="error"
              :closable="false"
            />
            <el-form-item style="text-align:left">
              <el-button
                type="primary"
                @click="domainLogin"
              >登陆</el-button>
            </el-form-item>
          </el-form>
        </div>
      </el-col>
      <el-col
        :xs="24"
        :md="8"
      >
        <div>&nbsp;</div>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import { RPC } from "../rpc.js";
import { mapState } from "vuex";

export default {
  name: "home",
  data() {
    return {
      loginForm: {
        username: "",
        password: ""
      },
      loginSuccess: true,
      errorText: ""
    };
  },
  computed: mapState([
    "domain"
  ]),
  methods: {
    domainLogin: async function() {
      let res = await RPC.doRPC("authUser", {
        domain: this.$route.params.domain,
        username: this.loginForm.username,
        password: this.loginForm.password
      });

      if (res == null) {
        this.loginSuccess = false;
        this.errorText = "未知错误";
        return;
      }

      if (!res.success) {
        this.loginSuccess = false;
        this.errorText = res.error;
        return;
      }

      RPC.setSession(res.sid);

      this.$bus.$emit("loginChange");

      this.$router.go(-1);
    }
  }
};
</script>

<style>
.el-col {
  border-radius: 4px;
}

.grid-content {
  border-radius: 4px;
  min-height: 36px;
}
</style>
