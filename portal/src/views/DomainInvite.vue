<template>
  <div>
    <div
      class="grid-content problem-title"
      style="text-align:left"
    >
      <el-button-group>
        <el-button
          v-if="!inviteEditor"
          @click="showInvite"
        >+ 创建邀请</el-button>
      </el-button-group>
    </div>

    <el-form
      v-if="inviteEditor"
      label-width="150px"
      style="text-align:left"
    >
      <el-form-item label="邀请UID">
        <el-input
          v-model="edit_invite.invite_uid"
          disabled
        ></el-input>
      </el-form-item>
      <el-form-item label="邀请密码">
        <el-input v-model="edit_invite.password"></el-input>
      </el-form-item>
      <el-form-item label="过期时间">
        <el-col :span="11">
          <el-date-picker
            type="datetime"
            placeholder="选择日期"
            v-model="edit_invite.valid_to"
            style="width: 100%;"
          ></el-date-picker>
        </el-col>
      </el-form-item>
      <el-form-item label="无需审核">
        <el-switch v-model="edit_invite.invite_state"></el-switch>
      </el-form-item>
      <el-form-item label="角色">
        <el-switch
          style="display: inline-block"
          active-color="#ff4949"
          inactive-color="#13ce66"
          active-text="管理员"
          inactive-text="普通用户"
          v-model="edit_invite.invite_role"
        ></el-switch>
      </el-form-item>
      <el-form-item label="文字说明">
        <el-input
          type="textarea"
          v-model="edit_invite.description"
        ></el-input>
      </el-form-item>
      <el-form-item>
        <el-button-group>
          <el-button
            type="primary"
            @click="commitInvite"
          >保存</el-button>
          <el-button @click="showInviteQR">生成二维码</el-button>
          <el-button @click="inviteEditor = false">取消</el-button>
        </el-button-group>
      </el-form-item>
    </el-form>

    <el-table
      :data="invites"
      v-if="!inviteEditor"
      @row-click="handleInviteClick"
      style="width: 100%"
    >
      <el-table-column label="UID">
        <template slot-scope="scope">
          <code>{{ scope.row.invite_uid }}</code>
        </template>
      </el-table-column>
      <el-table-column label="过期时间">
        <template slot-scope="scope">
          {{ scope.row.valid_to | moment("llll")}}
        </template>
      </el-table-column>
      <el-table-column label="密码">
        <template slot-scope="scope">
          <code v-if="scope.row.password != ''">{{ scope.row.password }}</code>
          <i v-else>no password</i>
        </template>
      </el-table-column>
      <el-table-column label="邀请后状态">
        <template slot-scope="scope">
          <span>{{ scope.row.invite_state }}</span>
        </template>
      </el-table-column>
      <el-table-column label="邀请后角色">
        <template slot-scope="scope">
          <el-tag
            size="mini"
            :type="tags['ROLE_TAG_' + scope.row.invite_role][0]"
          >{{ tags['ROLE_TAG_' + scope.row.invite_role][1] }}</el-tag>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script>
import { RPC } from "../rpc.js";
import { ConstString } from "../consts.js";
import QrcodeVue from "qrcode.vue";
import { mapState } from "vuex";

export default {
  data() {
    return {
      invites: [],
      edit_invite: {},
      loading: false,
      error: false,
      inviteEditor: false,
      tags: ConstString
    };
  },
  components: {
    QrcodeVue
  },
  methods: {
    showInviteQR: function() {
      let qr =
        location.origin +
        this.$router.resolve(
          "/" +
            (this.domain.alias || this.domain.uid) +
            "/join/" +
            this.edit_invite.invite_uid +
            (this.edit_invite.password !== ""
              ? "?p=" + encodeURI(this.edit_invite.password)
              : "")
        ).href;
      const h = this.$createElement;
      this.$msgbox({
        title: "扫描下方二维码 加入 " + this.domain.title,
        message: h("p", { style: "text-align:center" }, [
          h("qrcode-vue", { props: { value: qr, size: 256, level: "H" } }, ""),
          h("p", { style: "font-size:12px;" }, [
            "加组网址：",
            h("a", { attrs: { href: qr }, style: "font-family:monospace;" }, qr)
          ]),
          h(
            "p",
            null,
            "有效期至 " +
              this.$moment(new Date(this.edit_invite.valid_to)).format("llll")
          )
        ])
      });
    },
    loadInvites: async function() {
      this.loading = true;
      let res = await RPC.doRPC("getDomainInvites", {
        domain: this.$route.params.domain
      });
      this.loading = false;
      if (res == null) {
        this.error = true;
        return;
      }
      this.invites = res.domain_invites;
    },
    handleInviteClick(e) {
      this.edit_invite = JSON.parse(JSON.stringify(e));
      this.edit_invite.invite_state =
        this.edit_invite.invite_state == "CONFIRMED";
      this.edit_invite.invite_role = this.edit_invite.invite_role == "ADMIN";
      this.inviteEditor = true;
    },
    showInvite: async function() {
      this.edit_invite = {
        invite_state: false,
        invite_role: false,
        valid_to: new Date(new Date().valueOf() + 1000 * 24 * 3600),
        password: "",
        description: ""
      };
      this.inviteEditor = true;
    },
    commitInvite: async function() {
      let tin = JSON.parse(JSON.stringify(this.edit_invite));
      tin.valid_to = new Date(tin.valid_to);
      tin.invite_state = tin.invite_state ? "CONFIRMED" : "PENDING";
      tin.invite_role = tin.invite_role ? "ADMIN" : "USER";
      await RPC.doRPC("createDomainInvite", {
        domain: this.$route.params.domain,
        domain_invite: tin
      });
      this.loadInvites();
      this.inviteEditor = false;
    }
  },
  computed: mapState(["domain"]),
  created() {
    this.loadInvites();
  }
};
</script>