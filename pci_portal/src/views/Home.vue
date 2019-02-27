<template>
  <div class="home">
    <el-row>
      <el-col :span="24">
        <div
          class="grid-content problem-title"
          style="text-align:left"
        >
          <span>
            <h3><small id="domain-alias">/pci/</small><br>Problem CI<br><small id="domain-powered">Powered 千練萬花</small></h3>
          </span>
        </div>
      </el-col>
      <el-col :span="24">
        <div
          class="grid-content problem-content"
        >
        请登陆后安全食用 / Please login.<br>
        您可能需要请求WOJ管理员以获得使用权限 / Please contact WHU ACM team for privilege to use this system.<br>
        本系统对武汉大学教师及学生开放
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script>
export default {
  name: "home"
};
</script>

<style>
.problem-title {
  font-size: 28pt;
  padding-bottom: 12pt;
}

.problem-content {
  text-align: left;
  font-size: 14pt;
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

#domain-alias {
  font-family: "Courier New", Courier, monospace;
  font-size: 20pt;
  color: #260817;
}

#domain-powered {
  font-size: 14pt;
}
</style>

