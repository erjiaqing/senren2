<template>
  <div class="home">
    <el-row :gutter="10">
      <el-col
        :xs="24"
        :md="6"
      >
        <div>&nbsp;</div>
      </el-col>
      <el-col
        :xs="24"
        :md="12"
      >
        <div class="grid-content">
          <h3>加入到 {{ domain ? (domain.title || domain.alias || domain.uid) : '千练万花' }} 小组</h3>
          <!-- TODO: 设成当前域的名称 -->
          <el-form
            label-position="left"
            label-width="80px"
          >
            <!-- <el-form-item label="说明"> -->
              <pre style="text-align:left;font-size: 10pt;">{{inviteInfo.description}}</pre>
            <!-- </el-form-item> -->
            <el-form-item label="邀请码">
              <el-input
                autocomplete="off"
                v-model="loginForm.invite_code"
                disabled
              ></el-input>
            </el-form-item>
            <el-form-item
              label="邀请密码"
              v-if="inviteInfo.password != ''"
            >
              <el-input
                type="password"
                autocomplete="off"
                v-model="loginForm.invite_password"
              ></el-input>
            </el-form-item>
            <el-form-item label="组内昵称">
              <el-input
                autocomplete="off"
                v-model="loginForm.nickname"
              ></el-input>
            </el-form-item>
            <el-alert
              :title="errorText"
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
        :md="6"
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
        invite_code: "",
        invite_password: "",
        nickname: ""
      },
      inviteInfo: {},
      loginSuccess: true,
      errorText: ""
    };
  },
  computed: mapState(["domain", "guser"]),
  methods: {
    loadDomain: async function() {
      this.loading = true;
      let res = await RPC.doRPC("getDomainInvite", {
        domain: this.$route.params.domain,
        uid: this.$route.params.uid
      });
      this.loading = false;
      if (res == null) {
        this.error = true;
        this.errorText = "获取小组邀请码信息失败"
        return;
      }
      if (!res.success) {
        this.$router.push("/"+(this.$route.params.domain || "woj")+"/login");
        return
      }
      this.inviteInfo = res.domain_invite;
      this.loginForm.invite_code = res.domain_invite.invite_uid;
      this.loginForm.nickname = this.guser.nickname;
    },
    domainLogin: async function() {
      let res = await RPC.doRPC("joinDomain", {
        domain: this.$route.params.domain,
        invite_code: this.loginForm.invite_code,
        invite_password: this.loginForm.invite_password,
        nickname: this.loginForm.nickname
      });

      if (res == null) {
        this.loginSuccess = false;
        this.errorText = "未知错误";
        return;
      }

      if (!res.success) {
        this.loginSuccess = false;
        this.errorText = res.error;
        console.log(res);
        return;
      }
      this.$bus.$emit("loginChange");
      this.$router.push("/" + this.$route.params.domain);
    }
  },
  created() {
    this.loadDomain();
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
