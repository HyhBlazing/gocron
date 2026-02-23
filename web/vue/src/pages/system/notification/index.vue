<template>
  <el-container>
    <system-sidebar></system-sidebar>
    <el-main>
      <div class="page-header">
        <h2 class="page-title">通知配置</h2>
      </div>

      <el-tabs v-model="activeTab" class="ios-tabs">
        <!-- 邮件配置 -->
        <el-tab-pane label="邮件" name="email">
          <div class="config-container">
            <el-form ref="mailForm" :model="mailForm" :rules="mailRules" label-width="120px">
              <h3 class="section-title">SMTP服务器配置</h3>
              <el-row :gutter="20">
                <el-col :span="14">
                  <el-form-item label="SMTP服务器" prop="host">
                    <el-input v-model.trim="mailForm.host" placeholder="例如: smtp.example.com"></el-input>
                  </el-form-item>
                </el-col>
                <el-col :span="10">
                  <el-form-item label="端口" prop="port">
                    <el-input v-model.number="mailForm.port" placeholder="常见: 465, 587"></el-input>
                  </el-form-item>
                </el-col>
              </el-row>
              <el-row :gutter="20">
                <el-col :span="12">
                  <el-form-item label="用户名" prop="user">
                    <el-input v-model.trim="mailForm.user" placeholder="邮箱账号"></el-input>
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="密码" prop="password">
                    <el-input v-model="mailForm.password" type="password" placeholder="邮箱密码或授权码" show-password></el-input>
                  </el-form-item>
                </el-col>
              </el-row>

              <el-alert
                title="通知模板支持html"
                type="info"
                :closable="false"
                style="margin-bottom: 20px;">
              </el-alert>

              <el-form-item label="模板" prop="template">
                <el-input type="textarea" :rows="6" v-model="mailForm.template" placeholder="输入邮件通知模板..."></el-input>
              </el-form-item>

              <el-form-item>
                <el-button type="primary" @click="saveMailConfig" :loading="saving">保存邮件配置</el-button>
              </el-form-item>

              <div class="divider"></div>

              <div class="section-header">
                <h3 class="section-title">通知用户列表</h3>
                <el-button type="primary" icon="el-icon-plus" @click="openUserModal">新增用户</el-button>
              </div>
              <div class="tag-group">
                <el-tag
                  v-for="user in mailReceivers"
                  :key="user.id"
                  closable
                  class="ios-tag"
                  @close="deleteMailUser(user)">
                  <span class="tag-main">{{ user.username }}</span>
                  <span class="tag-sub">{{ user.email }}</span>
                </el-tag>
                <div v-if="mailReceivers.length === 0" class="empty-tip">暂无通知用户</div>
              </div>
            </el-form>
          </div>
        </el-tab-pane>

        <!-- Webhook配置 -->
        <el-tab-pane label="Webhook" name="webhook">
          <div class="config-container">
            <el-form ref="webhookForm" :model="webhookForm" :rules="webhookRules" label-width="120px">
              <el-alert
                title="通知内容推送到指定URL, POST请求, 设置Header[ Content-Type: application/json]"
                type="info"
                :closable="false"
                show-icon
                class="webhook-alert">
              </el-alert>
              <el-form-item label="URL" prop="url">
                <el-input v-model.trim="webhookForm.url" placeholder="https://api.yourdomain.com/callback"></el-input>
              </el-form-item>
              <el-form-item label="模板" prop="template">
                <el-input type="textarea" :rows="8" v-model.trim="webhookForm.template" placeholder="输入JSON格式的模板..."></el-input>
              </el-form-item>
              <el-form-item>
                <el-button type="primary" @click="saveWebhookConfig" :loading="saving">保存Webhook配置</el-button>
              </el-form-item>
            </el-form>
          </div>
        </el-tab-pane>
      </el-tabs>

      <div class="variable-guide" v-pre>
        <h4 class="guide-title"><i class="el-icon-info"></i> 模板支持的变量参考</h4>
        <div class="guide-content">
          <div class="v-item"><code>&#123;&#123;.TaskId&#125;&#125;</code><span>定时ID</span></div>
          <div class="v-item"><code>&#123;&#123;.TaskName&#125;&#125;</code><span>定时名称</span></div>
          <div class="v-item"><code>&#123;&#123;.Status&#125;&#125;</code><span>执行状态 (Success/Failure)</span></div>
          <div class="v-item"><code>&#123;&#123;.Result&#125;&#125;</code><span>定时执行输出</span></div>
        </div>
      </div>

      <!-- 弹窗部分 -->
      <!-- 邮件用户弹窗 -->
      <el-dialog title="新增邮件通知用户" :visible.sync="userModalVisible" width="400px" custom-class="pwd-dialog" :append-to-body="true" @closed="resetUserForm">
        <el-form ref="userForm" :model="userForm" label-width="80px" class="pwd-form">
          <el-form-item label="用户名">
            <el-input v-model.trim="userForm.username" placeholder="用户姓名"></el-input>
          </el-form-item>
          <el-form-item label="邮箱地址">
            <el-input v-model.trim="userForm.email" placeholder="接收通知的邮箱"></el-input>
          </el-form-item>
        </el-form>
        <div slot="footer">
          <el-button @click="userModalVisible = false">取消</el-button>
          <el-button type="primary" @click="saveMailUser">确认新增</el-button>
        </div>
      </el-dialog>
    </el-main>
  </el-container>
</template>

<script>
import systemSidebar from '../sidebar'
import notificationService from '../../../api/notification'

export default {
  name: 'notification-index',
  components: { systemSidebar },
  data () {
    return {
      activeTab: 'email',
      saving: false,
      // 邮件数据
      mailForm: { host: '', port: 465, user: '', password: '', template: '' },
      mailReceivers: [],
      mailRules: {
        host: [{ required: true, message: '请输入SMTP服务器地址', trigger: 'blur' }],
        port: [{ required: true, type: 'number', message: '请输入有效的数字端口', trigger: 'blur' }],
        user: [{ required: true, message: '请输入账户邮箱', trigger: 'blur' }],
        password: [{ required: true, message: '请输入密码或授权码', trigger: 'blur' }],
        template: [{ required: true, message: '请输入通知模板内容', trigger: 'blur' }]
      },
      // Webhook数据
      webhookForm: { url: '', template: '' },
      webhookRules: {
        url: [{ type: 'url', required: true, message: '请输入有效的通知 URL', trigger: 'blur' }],
        template: [{ required: true, message: '请输入 Webhook 模板', trigger: 'blur' }]
      },
      // 弹窗表单
      userModalVisible: false,
      userForm: { username: '', email: '' }
    }
  },
  created () {
    this.initData()
  },
  methods: {
    initData () {
      // 预加载当前 tab 的数据
      this.fetchMailConfig()
      this.fetchWebhookConfig()
    },
    handleTabClick (tab) {
      // 可以在此处按需拉取数据
    },
    // --- 邮件操作 ---
    fetchMailConfig () {
      notificationService.mail((data) => {
        this.mailForm.host = data.host
        if (data.port) this.mailForm.port = data.port
        this.mailForm.user = data.user
        this.mailForm.password = data.password
        this.mailForm.template = data.template
        this.mailReceivers = data.mail_users
      })
    },
    saveMailConfig () {
      this.$refs.mailForm.validate(valid => {
        if (!valid) return
        this.saving = true
        notificationService.updateMail(this.mailForm, () => {
          this.saving = false
          this.$message.success('邮件服务器配置已保存')
          this.fetchMailConfig()
        })
      })
    },
    openUserModal () { this.userModalVisible = true },
    resetUserForm () { this.userForm = { username: '', email: '' } },
    saveMailUser () {
      if (!this.userForm.username || !this.userForm.email) {
        return this.$message.error('请填写完整的用户信息')
      }
      notificationService.createMailUser(this.userForm, () => {
        this.userModalVisible = false
        this.$message.success('通知用户已添加')
        this.fetchMailConfig()
      })
    },
    deleteMailUser (user) {
      this.$appConfirm(() => {
        notificationService.removeMailUser(user.id, () => {
          this.$message.success('用户已移除')
          this.fetchMailConfig()
        })
      })
    },
    // --- Webhook 操作 ---
    fetchWebhookConfig () {
      notificationService.webhook((data) => {
        this.webhookForm.url = data.url
        this.webhookForm.template = data.template
      })
    },
    saveWebhookConfig () {
      this.$refs.webhookForm.validate(valid => {
        if (!valid) return
        this.saving = true
        notificationService.updateWebHook(this.webhookForm, () => {
          this.saving = false
          this.$message.success('Webhook 配置已保存')
          this.fetchWebhookConfig()
        })
      })
    }
  }
}
</script>

<style scoped>

.page-title {
  font-size: 20px;
  font-weight: 600;
  color: #303133;
  margin: 0;
}

.config-container {
  background: rgba(144, 147, 153, 0.05);
  border-radius: 12px;
  padding: 24px;
  border: 1px solid rgba(144, 147, 153, 0.1);
}

.section-title {
  font-size: 16px;
  font-weight: 600;
  color: #409EFF;
  margin: 0 0 20px 0;
  letter-spacing: 0.5px;
}

.divider {
  height: 1px;
  background: rgba(144, 147, 153, 0.15);
  margin: 30px 0;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.tag-group {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  min-height: 32px;
}

.ios-tag {
  background: rgba(144, 147, 153, 0.1) !important;
  border: 1px solid rgba(144, 147, 153, 0.2) !important;
  color: #303133 !important;
  border-radius: 8px;
  display: inline-flex;
  align-items: center;
  max-width: 100%;
  overflow: hidden;
  height: 32px;
  line-height: 30px;
  padding: 0 10px;
  box-sizing: border-box;
  white-space: nowrap;
  font-weight: 500;
}

.tag-main {
  flex-shrink: 0;
  max-width: 120px;
  overflow: hidden;
  text-overflow: ellipsis;
}

.tag-sub {
  color: #909399;
  font-size: 12px;
  margin-left: 6px;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.ios-tag >>> .el-tag__close {
  margin-left: 6px;
}

.empty-tip {
  color: #909399;
  font-style: italic;
  font-size: 14px;
}

.webhook-alert {
  margin-bottom: 20px;
  background-color: rgba(64, 158, 255, 0.1) !important;
  border: 1px solid rgba(64, 158, 255, 0.2) !important;
}

/* 变量指南卡片 */
.variable-guide {
  margin-top: 32px;
  background: rgba(144, 147, 153, 0.05);
  border-radius: 12px;
  padding: 16px 20px;
  border: 1px solid rgba(144, 147, 153, 0.1);
}

.guide-title {
  margin: 0 0 12px 0;
  font-size: 14px;
  color: #606266;
  display: flex;
  align-items: center;
  gap: 8px;
}

.guide-content {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
  gap: 10px 20px;
}

.v-item {
  display: flex;
  align-items: center;
  gap: 12px;
}

.v-item code {
  background: rgba(230, 162, 60, 0.1);
  padding: 2px 6px;
  border-radius: 4px;
  color: #E6A23C;
  font-family: Menlo, Monaco, Consolas, "Courier New", monospace;
  font-size: 13px;
}

.v-item span {
  font-size: 13px;
  color: #909399;
}

/* 覆盖 tabs 样式 */
.ios-tabs >>> .el-tabs__item {
  color: #909399;
  font-size: 15px;
  font-weight: 500;
}

.ios-tabs >>> .el-tabs__item.is-active {
  color: #409EFF;
}

.ios-tabs >>> .el-tabs__active-bar {
  background-color: #409EFF;
}

.ios-tabs >>> .el-tabs__nav-wrap::after {
  background-color: rgba(144, 147, 153, 0.1);
}

</style>
