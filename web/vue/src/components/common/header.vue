<template>
  <div class="ios-header" v-cloak>
    <div class="ios-nav">
      <!-- Logo / Brand -->
      <div class="ios-brand" title="剧星传媒 - 定时任务系统">
        <img src="/static/logo.svg" alt="logo" class="ios-brand-logo">
      </div>

      <div class="ios-nav-spacer"></div>

      <!-- User Area -->
      <div class="ios-user-area" v-if="$store.getters.user.token">
        <div class="ios-avatar-btn" @click="toggleDropdown" ref="userBtn">
          <span class="ios-avatar-icon">{{ usernameInitial }}</span>
          <span class="ios-username">{{ $store.getters.user.username }}</span>
          <span class="ios-chevron" :class="{ open: dropdownOpen }">›</span>
        </div>
        <transition name="ios-dropdown-fade">
          <div class="ios-dropdown" v-if="dropdownOpen">
            <div class="ios-dropdown-item" @click="openPwdModal">修改密码</div>
            <div class="ios-dropdown-divider"></div>
            <div class="ios-dropdown-item danger" @click="logout">退出登录</div>
          </div>
        </transition>
      </div>
    </div>

    <!-- 修改密码 Modal -->
    <el-dialog title="修改密码" :visible.sync="pwdModalVisible" width="420px" :close-on-click-modal="false" :append-to-body="true" custom-class="pwd-dialog" @closed="resetPwdForm">
      <el-form ref="pwdForm" :model="pwdForm" :rules="pwdRules" label-width="90px" class="pwd-form">
        <el-form-item label="原密码" prop="old_password">
          <el-input v-model="pwdForm.old_password" type="password" placeholder="请输入原密码" show-password></el-input>
        </el-form-item>
        <el-form-item label="新密码" prop="new_password">
          <el-input v-model="pwdForm.new_password" type="password" placeholder="请输入新密码" show-password></el-input>
        </el-form-item>
        <el-form-item label="确认密码" prop="confirm_new_password">
          <el-input v-model="pwdForm.confirm_new_password" type="password" placeholder="请再次输入新密码" show-password></el-input>
        </el-form-item>
      </el-form>
      <span slot="footer">
        <el-button @click="pwdModalVisible = false">取消</el-button>
        <el-button type="primary" @click="submitPwd" :loading="pwdLoading">保存</el-button>
      </span>
    </el-dialog>
  </div>
</template>

<script>
import userService from '../../api/user'
export default {
  name: 'app-header',
  data () {
    return {
      dropdownOpen: false,
      pwdModalVisible: false,
      pwdLoading: false,
      pwdForm: {
        old_password: '',
        new_password: '',
        confirm_new_password: ''
      },
      pwdRules: {
        old_password: [{ required: true, message: '请输入原密码', trigger: 'blur' }],
        new_password: [{ required: true, message: '请输入新密码', trigger: 'blur' }],
        confirm_new_password: [{ required: true, message: '请再次输入新密码', trigger: 'blur' }]
      }
    }
  },
  computed: {
    usernameInitial () {
      const name = this.$store.getters.user.username || ''
      return name.charAt(0).toUpperCase()
    }
  },
  methods: {
    toggleDropdown () {
      this.dropdownOpen = !this.dropdownOpen
    },
    openPwdModal () {
      this.dropdownOpen = false
      this.pwdModalVisible = true
    },
    resetPwdForm () {
      this.pwdForm = { old_password: '', new_password: '', confirm_new_password: '' }
      this.$refs.pwdForm && this.$refs.pwdForm.resetFields()
    },
    submitPwd () {
      this.$refs.pwdForm.validate((valid) => {
        if (!valid) return
        this.pwdLoading = true
        userService.editMyPassword(this.pwdForm, () => {
          this.pwdLoading = false
          this.pwdModalVisible = false
          this.$message.success('密码修改成功')
        })
      })
    },
    logout () {
      this.dropdownOpen = false
      this.$store.commit('logout')
      this.$router.push('/user/login')
    },
    handleOutsideClick (e) {
      if (this.$refs.userBtn && !this.$refs.userBtn.contains(e.target)) {
        this.dropdownOpen = false
      }
    }
  },
  mounted () {
    document.addEventListener('click', this.handleOutsideClick)
  },
  beforeDestroy () {
    document.removeEventListener('click', this.handleOutsideClick)
  }
}
</script>

<style scoped>
/* ─── Base ─────────────────────────────────────────── */
* {
  box-sizing: border-box;
}

.ios-header {
  position: relative;
  z-index: 1000;
  width: 100%;
  background: #0f0f0f;
  backdrop-filter: saturate(180%) blur(20px);
  -webkit-backdrop-filter: saturate(180%) blur(20px);
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
  box-shadow: 0 1px 0 rgba(0, 0, 0, 0.3);
}

/* ─── Nav Row ────────────────────────────────────────  */
.ios-nav {
  display: flex;
  align-items: center;
  height: 52px;
  padding: 0 20px;
  gap: 8px;
}

/* ─── Brand ─────────────────────────────────────────── */
.ios-brand {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-right: 16px;
  text-decoration: none;
  cursor: default;
}

.ios-brand-logo {
  height: 32px;
  width: auto;
  display: block;
  filter: brightness(0) invert(1);
  transition: opacity 0.2s ease;
}

.ios-brand:hover .ios-brand-logo {
  opacity: 0.85;
}

.ios-nav-spacer {
  flex: 1;
}

/* ─── User Area ──────────────────────────────────────── */
.ios-user-area {
  position: relative;
  margin-left: auto;
}

.ios-avatar-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 5px 12px 5px 6px;
  border-radius: 20px;
  cursor: pointer;
  transition: background 0.18s ease;
  user-select: none;
}

.ios-avatar-btn:hover {
  background: rgba(255, 255, 255, 0.08);
}

.ios-avatar-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border-radius: 50%;
  background: linear-gradient(135deg, #007aff, #5856d6);
  color: #fff;
  font-size: 12px;
  font-weight: 700;
  font-family: -apple-system, BlinkMacSystemFont, sans-serif;
}

.ios-username {
  font-family: -apple-system, BlinkMacSystemFont, "SF Pro Text", "PingFang SC",
    sans-serif;
  font-size: 14px;
  font-weight: 500;
  color: rgba(255, 255, 255, 0.85);
}

.ios-chevron {
  font-size: 18px;
  color: rgba(255, 255, 255, 0.35);
  line-height: 1;
  display: inline-block;
  transform: rotate(90deg);
  transition: transform 0.2s ease;
}

.ios-chevron.open {
  transform: rotate(-90deg);
}

/* ─── Dropdown ───────────────────────────────────────── */
.ios-dropdown {
  position: absolute;
  right: 0;
  top: calc(100% + 8px);
  min-width: 160px;
  background: rgba(44, 44, 46, 0.95);
  backdrop-filter: saturate(180%) blur(24px);
  -webkit-backdrop-filter: saturate(180%) blur(24px);
  border-radius: 14px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.5), 0 2px 8px rgba(0, 0, 0, 0.3);
  overflow: hidden;
  z-index: 999;
}

.ios-dropdown-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 11px 16px;
  font-family: -apple-system, BlinkMacSystemFont, "SF Pro Text", "PingFang SC",
    sans-serif;
  font-size: 14px;
  font-weight: 400;
  color: rgba(255, 255, 255, 0.85);
  text-decoration: none;
  cursor: pointer;
  transition: background 0.14s ease;
}

.ios-dropdown-item:hover {
  background: rgba(255, 255, 255, 0.08);
}

.ios-dropdown-item.danger {
  color: #ff453a;
}

.ios-dropdown-item.danger:hover {
  background: rgba(255, 69, 58, 0.12);
}

.ios-dropdown-icon {
  font-size: 15px;
}

.ios-dropdown-divider {
  height: 1px;
  background: rgba(255, 255, 255, 0.08);
  margin: 0 12px;
}

/* ─── Dropdown Transition ────────────────────────────── */
.ios-dropdown-fade-enter-active,
.ios-dropdown-fade-leave-active {
  transition: opacity 0.18s ease, transform 0.18s ease;
  transform-origin: top right;
}

.ios-dropdown-fade-enter,
.ios-dropdown-fade-leave-to {
  opacity: 0;
  transform: scale(0.92) translateY(-6px);
}
</style>

<!-- 全局覆盖：pwd-dialog iOS 暗色风格（append-to-body 后 scoped 无效，用全局 CSS） -->
<style>
/* ─── 遮罩 ─── */
.v-modal {
  background: rgba(0, 0, 0, 0.6) !important;
  backdrop-filter: blur(6px);
  -webkit-backdrop-filter: blur(6px);
}

/* ─── Dialog 卡片 ─── */
.pwd-dialog {
  background: #1c1c1e !important;
  border-radius: 20px !important;
  border: 1px solid rgba(255, 255, 255, 0.08);
  box-shadow: 0 24px 64px rgba(0, 0, 0, 0.6) !important;
  overflow: hidden;
}

/* ─── 标题栏 ─── */
.pwd-dialog .el-dialog__header {
  background: #1c1c1e;
  border-bottom: 1px solid rgba(255, 255, 255, 0.07);
  padding: 18px 24px 14px;
}

.pwd-dialog .el-dialog__title {
  font-family: -apple-system, BlinkMacSystemFont, "SF Pro Display", "PingFang SC", sans-serif;
  font-size: 16px;
  font-weight: 600;
  color: #ffffff;
  letter-spacing: -0.2px;
}

.pwd-dialog .el-dialog__headerbtn .el-dialog__close {
  color: rgba(255, 255, 255, 0.4);
  font-size: 18px;
  transition: color 0.15s ease;
}

.pwd-dialog .el-dialog__headerbtn:hover .el-dialog__close {
  color: rgba(255, 255, 255, 0.8);
}

/* ─── 内容区 ─── */
.pwd-dialog .el-dialog__body {
  background: #1c1c1e;
  padding: 20px 24px 12px;
}

/* ─── 表单 label ─── */
.pwd-dialog .el-form-item__label {
  font-family: -apple-system, BlinkMacSystemFont, "SF Pro Text", "PingFang SC", sans-serif;
  font-size: 14px;
  color: rgba(255, 255, 255, 0.55);
}

/* ─── 输入框 ─── */
.pwd-dialog .el-input__inner {
  background: #2c2c2e !important;
  border: 1px solid rgba(255, 255, 255, 0.1) !important;
  border-radius: 10px !important;
  color: #ffffff !important;
  font-family: -apple-system, BlinkMacSystemFont, "SF Pro Text", "PingFang SC", sans-serif;
  font-size: 14px;
  height: 40px;
  transition: border-color 0.18s ease;
}

.pwd-dialog .el-input__inner:focus {
  border-color: #0a84ff !important;
  box-shadow: 0 0 0 3px rgba(10, 132, 255, 0.18) !important;
}

.pwd-dialog .el-input__inner::placeholder {
  color: rgba(255, 255, 255, 0.25);
}

/* 显示/隐藏密码图标 */
.pwd-dialog .el-input__suffix .el-input__icon {
  color: rgba(255, 255, 255, 0.35);
}

/* ─── 底部按钮区 ─── */
.pwd-dialog .el-dialog__footer {
  background: #1c1c1e;
  border-top: 1px solid rgba(255, 255, 255, 0.07);
  padding: 14px 24px 18px;
}

/* 取消按钮 */
.pwd-dialog .el-dialog__footer .el-button--default {
  background: #2c2c2e;
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 10px;
  color: rgba(255, 255, 255, 0.7);
  font-family: -apple-system, BlinkMacSystemFont, "SF Pro Text", "PingFang SC", sans-serif;
  font-size: 14px;
  transition: background 0.15s ease;
}

.pwd-dialog .el-dialog__footer .el-button--default:hover {
  background: #3a3a3c;
  color: #ffffff;
}

/* 保存按钮 */
.pwd-dialog .el-dialog__footer .el-button--primary {
  background: #0a84ff;
  border-color: #0a84ff;
  border-radius: 10px;
  font-family: -apple-system, BlinkMacSystemFont, "SF Pro Text", "PingFang SC", sans-serif;
  font-size: 14px;
  font-weight: 600;
  transition: background 0.15s ease, opacity 0.15s ease;
}

.pwd-dialog .el-dialog__footer .el-button--primary:hover {
  background: #0071e3;
  border-color: #0071e3;
}

/* ─── 校验错误提示 ─── */
.pwd-dialog .el-form-item__error {
  color: #ff453a;
  font-size: 12px;
}
</style>
