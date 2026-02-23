<template>
  <div v-cloak>
    <el-menu :default-active="currentRoute" mode="horizontal" background-color="#545c64" text-color="#fff" active-text-color="#ffd04b" router>
      <el-row>
        <el-col :span="22"></el-col>
        <el-col :span="2" style="float:right;">
          <el-submenu v-if="this.$store.getters.user.token" index="userStatus">
            <template slot="title">{{ this.$store.getters.user.username }}</template>
            <el-menu-item index="/user/edit-my-password">修改密码</el-menu-item>
            <el-menu-item @click="logout" index="/user/logout">退出</el-menu-item>
          </el-submenu>
        </el-col>
      </el-row>
    </el-menu>
  </div>
</template>

<script>

export default {
  name: 'app-nav-menu',
  data () {
    return {}
  },
  computed: {
    currentRoute () {
      return this.$route.path
    }
  },
  methods: {
    logout () {
      this.$store.commit('logout')
      this.$router.push('/')
    }
  }
}
</script>
