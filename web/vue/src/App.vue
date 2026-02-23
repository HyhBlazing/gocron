<template>
  <div id="app-root">
    <app-header v-if="!isLoginPage"></app-header>
    <div id="app-body">
      <router-view />
    </div>
    <app-footer v-if="!isLoginPage"></app-footer>
  </div>
</template>

<script>
import installService from './api/install'
import appHeader from './components/common/header.vue'
import appFooter from './components/common/footer.vue'

export default {
  name: 'App',
  data () {
    return {}
  },
  computed: {
    isLoginPage () {
      return this.$route.path === '/user/login'
    }
  },
  created () {
    installService.status((data) => {
      if (!data) {
        this.$router.push('/install')
      }
    })
  },
  components: {
    appHeader,
    appFooter
  }
}
</script>

<style>
/* ─── Reset ─── */
*,
*::before,
*::after {
  box-sizing: border-box;
}

html,
body {
  margin: 0;
  padding: 0;
  height: 100vh;
  overflow: hidden;
  background: #f6f5f8;
  font-family: -apple-system, BlinkMacSystemFont, "SF Pro Text", "PingFang SC", "Helvetica Neue", sans-serif;
}

[v-cloak] {
  display: none !important;
}

/* ─── App Shell ─── */
#app-root {
  height: 100vh;
  display: flex;
  flex-direction: column;
}

#app-body {
  flex: 1;
  display: flex;
  overflow: hidden;
  /* 让内部 el-container（侧边栏+主内容）横排 */
  align-items: stretch;
}

/* ─── Element UI overrides ─── */
.el-container {
  flex: 1;
  padding: 0 !important;
  margin: 0 !important;
  width: 100% !important;
  align-items: stretch;
  height: 100%;
}

.el-header {
  padding: 0 !important;
  height: auto !important;
}

.el-main {
  padding: 20px !important;
  overflow-y: auto;
  height: 100%;
}

.el-aside {
  overflow: visible !important;
}

.el-footer {
  padding: 0 !important;
  height: auto !important;
}

/* ─── 模态窗居中 ─── */
.el-dialog__wrapper {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}

.el-dialog__wrapper .el-dialog {
  margin: 0 !important;
}

/* ─── 全局模态窗外观统一（对齐定时日志结果弹窗） ─── */
.el-dialog {
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 18px 48px rgba(17, 17, 17, 0.18);
}

.el-dialog__header {
  padding: 16px 20px;
  border-bottom: 1px solid #e5e5ea;
  background: #ffffff;
}

.el-dialog__title {
  font-size: 16px;
  font-weight: 600;
  color: #1d1d1f;
}

.el-dialog__headerbtn .el-dialog__close {
  color: #8e8e93;
  transition: color 0.2s ease;
}

.el-dialog__headerbtn:hover .el-dialog__close {
  color: #007aff;
}

.el-dialog__body {
  padding: 20px 24px 28px;
  background: #f5f5f7;
}

.el-message-box__wrapper {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}

.el-message-box__wrapper::after {
  display: none;
}

.el-message-box__wrapper .el-message-box {
  vertical-align: middle;
}
</style>
