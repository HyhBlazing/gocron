<template>
  <el-container>
    <task-sidebar></task-sidebar>
    <el-main>
      <div class="page-header">
        <h2 class="page-title">{{ form.id ? '编辑定时' : '新增定时' }}</h2>
      </div>

      <el-form ref="form" :model="form" :rules="formRules" label-position="top" class="ios-form compact-form">
        <!-- 基础信息卡片 -->
        <div class="form-section">
          <div class="section-header">
            <i class="el-icon-info"></i>
            <span>基础信息</span>
          </div>
          <div class="section-content">
            <el-row :gutter="10">
              <el-col :span="6">
                <el-form-item label="定时名称" prop="name">
                  <el-input v-model.trim="form.name" placeholder="定时名称"></el-input>
                </el-form-item>
              </el-col>
              <el-col :span="6">
                <el-form-item class="label-with-tooltip">
                  <template slot="label">
                    <span>定时类型
                      <el-tooltip content="主定时可以配置多个子定时，当主定时执行完成后，自动执行子定时。定时类型新增后不能变更。" placement="top">
                        <i class="el-icon-question label-tooltip-icon"></i>
                      </el-tooltip>
                    </span>
                  </template>
                  <el-select v-model="form.level" :disabled="form.id !== ''" style="width: 100%">
                    <el-option v-for="item in levelList" :key="item.value" :label="item.label" :value="item.value"></el-option>
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :span="6" v-if="form.level === 1">
                <el-form-item prop="spec" class="cron-spec-form-item">
                  <template slot="label">
                    <span class="cron-label-wrap">
                      <span>Cron 表达式</span>
                      <el-tooltip content="点击查看 Cron 表达式说明和示例" placement="top">
                        <i class="el-icon-question label-tooltip-icon cron-spec-tooltip-icon" @click.stop="showCronSpecModal"></i>
                      </el-tooltip>
                    </span>
                  </template>
                  <el-input v-model.trim="form.spec" placeholder="秒 分 时 日 周 月">
                    <span slot="suffix" class="cron-next-runs-link" @click.stop="showNextRunsModal">未来执行时间</span>
                  </el-input>
                </el-form-item>
              </el-col>
              <el-col :span="form.level === 1 ? 6 : 12">
                <el-form-item label="标签">
                  <el-input v-model.trim="form.tag" placeholder="标签"></el-input>
                </el-form-item>
              </el-col>
            </el-row>
          </div>
        </div>

        <!-- 执行策略卡片 -->
        <div class="form-section">
          <div class="section-header">
            <i class="el-icon-copy-document"></i>
            <span>执行策略</span>
          </div>
          <div class="section-content">
            <el-row :gutter="10">
              <el-col :span="6">
                <el-form-item label="执行方式">
                  <el-select v-model="form.protocol" style="width: 100%">
                    <el-option v-for="item in protocolList" :key="item.value" :label="item.label" :value="item.value"></el-option>
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :span="6" v-if="form.protocol === 2">
                <el-form-item label="定时节点 (必选)">
                  <el-select :key="'host-select-' + form.protocol" v-model="selectedHosts" filterable multiple placeholder="请选择执行节点" style="width: 100%">
                    <el-option v-for="item in hosts" :key="item.id" :label="item.alias + ' - ' + item.name" :value="item.id"></el-option>
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :span="6" v-if="form.protocol === 1">
                <el-form-item label="请求方法">
                  <el-select v-model="form.http_method" style="width: 100%">
                    <el-option v-for="item in httpMethods" :key="item.value" :label="item.label" :value="item.value"></el-option>
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :span="6" v-if="form.level === 1">
                <el-form-item class="label-with-tooltip">
                  <template slot="label">
                    <span>依赖关系
                      <el-tooltip :content="dependencyTooltipContent" placement="top" popper-class="label-tooltip-popper">
                        <i class="el-icon-question label-tooltip-icon"></i>
                      </el-tooltip>
                    </span>
                  </template>
                  <el-select v-model="form.dependency_status" style="width: 100%">
                    <el-option v-for="item in dependencyStatusList" :key="item.value" :label="item.label" :value="item.value"></el-option>
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :span="6" v-if="form.level === 1">
                <el-form-item label="子定时 ID">
                  <el-input v-model.trim="form.dependency_task_id" placeholder="多个 ID 请使用逗号分隔"></el-input>
                </el-form-item>
              </el-col>
            </el-row>

            <el-form-item label="执行命令" prop="command">
              <el-input type="textarea" :rows="4" :placeholder="commandPlaceholder" v-model="form.command" class="code-editor"></el-input>
            </el-form-item>
          </div>
        </div>

        <!-- 超时设置卡片 -->
        <div class="form-section">
          <div class="section-header">
            <i class="el-icon-timer"></i>
            <span>超时设置</span>
          </div>
          <div class="section-content">
            <el-row :gutter="10">
              <el-col :span="6">
                <el-form-item label="超时限制 (秒)" prop="timeout">
                  <el-input-number v-model.number="form.timeout" :min="0" :max="86400" style="width: 100%"></el-input-number>
                </el-form-item>
              </el-col>
              <el-col :span="6">
                <el-form-item class="label-with-tooltip">
                  <template slot="label">
                    <span>单实例运行
                      <el-tooltip content="单实例运行, 前次定时未执行完成，下次定时调度时间到了是否要执行, 即是否允许多进程执行同一定时
" placement="top">
                        <i class="el-icon-question label-tooltip-icon"></i>
                      </el-tooltip>
                    </span>
                  </template>
                  <el-select v-model="form.multi" style="width: 100%">
                    <el-option v-for="item in runStatusList" :key="item.value" :label="item.label" :value="item.value"></el-option>
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :span="6">
                <el-form-item label="失败重试次数" prop="retry_times">
                  <el-input-number v-model.number="form.retry_times" :min="0" :max="10" style="width: 100%"></el-input-number>
                </el-form-item>
              </el-col>
              <el-col :span="6">
                <el-form-item label="重试间隔 (秒)" prop="retry_interval">
                  <el-input-number v-model.number="form.retry_interval" :min="0" :max="3600" style="width: 100%"></el-input-number>
                </el-form-item>
              </el-col>
            </el-row>
          </div>
        </div>

        <!-- 其他卡片 -->
        <div class="form-section">
          <div class="section-header">
            <i class="el-icon-more"></i>
            <span>其他</span>
          </div>
          <div class="section-content">
            <el-row :gutter="10">
              <el-col :span="6">
                <el-form-item label="通知策略">
                  <el-select v-model="form.notify_status" style="width: 100%">
                    <el-option v-for="item in notifyStatusList" :key="item.value" :label="item.label" :value="item.value"></el-option>
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :span="6" v-if="form.notify_status !== 1">
                <el-form-item label="通知渠道">
                  <el-select v-model="form.notify_type" style="width: 100%">
                    <el-option v-for="item in notifyTypes" :key="item.value" :label="item.label" :value="item.value"></el-option>
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :span="6" v-if="form.notify_status !== 1 && form.notify_type === 2">
                <el-form-item label="接收用户">
                  <el-select v-model="selectedMailNotifyIds" filterable multiple placeholder="选择邮件接收者" style="width: 100%">
                    <el-option v-for="item in mailUsers" :key="item.id" :label="item.username" :value="item.id"></el-option>
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :span="6" v-if="form.notify_status !== 1 && form.notify_type === 4">
                <el-form-item label="WebHook URL">
                  <el-input v-model.trim="form.notify_receiver_id" placeholder="请输入 WebHook URL 地址"></el-input>
                </el-form-item>
              </el-col>
              <el-col :span="6" v-if="form.notify_status === 4">
                <el-form-item label="匹配关键字" prop="notify_keyword">
                  <el-input v-model.trim="form.notify_keyword" placeholder="定时输出中包含此关键字时触发通知"></el-input>
                </el-form-item>
              </el-col>
            </el-row>

            <el-form-item label="备注说明">
              <el-input type="textarea" :rows="2" v-model="form.remark" placeholder="备注信息"></el-input>
            </el-form-item>
          </div>
        </div>

        <div class="form-actions">
          <el-button @click="cancel">取消</el-button>
          <el-button type="primary" @click="submit" :loading="form.id === '' && false">保存</el-button>
        </div>
      </el-form>

      <el-dialog :visible.sync="cronSpecModalVisible" width="720px" custom-class="cron-spec-dialog-wrap" :show-close="true">
        <div slot="title" class="cron-spec-header">
          <span class="cron-spec-header-title">Cron 表达式说明</span>
        </div>
        <el-tabs v-model="cronSpecActiveTab" class="cron-spec-tabs">
          <el-tab-pane label="特殊字符与预定义" name="basic">
            <div class="cron-spec-content">
              <h4>特殊字符</h4>
              <ul>
                <li><strong>*</strong> 所有有效值</li>
                <li><strong>/</strong> 步长，如 */5、3-59/15</li>
                <li><strong>,</strong> 多个值，如 1,3,5</li>
                <li><strong>-</strong> 范围，如 9-17</li>
                <li><strong>?</strong> 日或周中允许为空（二选一）</li>
              </ul>
              <h4>预定义表达式</h4>
              <table class="cron-spec-table">
                <tr>
                  <th>表达式</th>
                  <th>含义</th>
                </tr>
                <tr>
                  <td>@yearly / @annually</td>
                  <td>每年 1 月 1 日 0:00</td>
                </tr>
                <tr>
                  <td>@monthly</td>
                  <td>每月 1 日 0:00</td>
                </tr>
                <tr>
                  <td>@weekly</td>
                  <td>每周日 0:00</td>
                </tr>
                <tr>
                  <td>@daily / @midnight</td>
                  <td>每天 0:00</td>
                </tr>
                <tr>
                  <td>@hourly</td>
                  <td>每小时整点</td>
                </tr>
              </table>
              <h4>固定间隔</h4>
              <p>格式：<code>@every &lt;duration&gt;</code></p>
              <p>示例：<code>@every 1h</code>、<code>@every 30m</code>、<code>@every 1h30m10s</code></p>
            </div>
          </el-tab-pane>
          <el-tab-pane label="常用示例" name="examples">
            <div class="cron-spec-content">
              <h4>时间表达式常用示例</h4>
              <table class="cron-spec-table">
                <tr>
                  <th>示例</th>
                  <th>Cron 表达式</th>
                </tr>
                <tr>
                  <td>每天 12:00 调度</td>
                  <td>
                    <span class="cron-spec-copy-wrap">
                      <code>0 0 12 * * *</code>
                      <button
                        type="button"
                        class="cron-spec-copy-btn"
                        @click.stop="copyCronExpr('0 0 12 * * *')"
                      >复制</button>
                    </span>
                  </td>
                </tr>
                <tr>
                  <td>每天 12:30 调度</td>
                  <td>
                    <span class="cron-spec-copy-wrap">
                      <code>0 30 12 * * *</code>
                      <button
                        type="button"
                        class="cron-spec-copy-btn"
                        @click.stop="copyCronExpr('0 30 12 * * *')"
                      >复制</button>
                    </span>
                  </td>
                </tr>
                <tr>
                  <td>每小时的 26、29、33 分调度</td>
                  <td>
                    <span class="cron-spec-copy-wrap">
                      <code>0 26,29,33 * * *</code>
                      <button
                        type="button"
                        class="cron-spec-copy-btn"
                        @click.stop="copyCronExpr('0 26,29,33 * * *')"
                      >复制</button>
                    </span>
                  </td>
                </tr>
                <tr>
                  <td>周一到周五每天 12:30 调度</td>
                  <td>
                    <span class="cron-spec-copy-wrap">
                      <code>0 30 12 ? * MON-FRI</code>
                      <button
                        type="button"
                        class="cron-spec-copy-btn"
                        @click.stop="copyCronExpr('0 30 12 ? * MON-FRI')"
                      >复制</button>
                    </span>
                  </td>
                </tr>
                <tr>
                  <td>周一到周五每天 12:00～14:00 每 5 分钟调度</td>
                  <td>
                    <span class="cron-spec-copy-wrap">
                      <code>0 0/5 12-13 ? * MON-FRI</code>
                      <button
                        type="button"
                        class="cron-spec-copy-btn"
                        @click.stop="copyCronExpr('0 0/5 12-13 ? * MON-FRI')"
                      >复制</button>
                    </span>
                  </td>
                </tr>
                <tr>
                  <td>一月到四月每天 12:00 调度</td>
                  <td>
                    <span class="cron-spec-copy-wrap">
                      <code>0 0 12 * 1-4 ?</code>
                      <button
                        type="button"
                        class="cron-spec-copy-btn"
                        @click.stop="copyCronExpr('0 0 12 * 1-4 ?')"
                      >复制</button>
                    </span>
                  </td>
                </tr>
              </table>
            </div>
          </el-tab-pane>
          <el-tab-pane label="字段说明" name="fields">
            <div class="cron-spec-content">
              <h4>标准 6 字段格式</h4>
              <table class="cron-spec-table">
                <tr>
                  <th>字段</th>
                  <th>必填</th>
                  <th>取值范围</th>
                  <th>特殊字符</th>
                </tr>
                <tr>
                  <td>秒</td>
                  <td>是</td>
                  <td>0～59</td>
                  <td>不允许设置特殊字符</td>
                </tr>
                <tr>
                  <td>分</td>
                  <td>是</td>
                  <td>0～59</td>
                  <td>, - * /</td>
                </tr>
                <tr>
                  <td>小时</td>
                  <td>是</td>
                  <td>0～23</td>
                  <td>, - * /</td>
                </tr>
                <tr>
                  <td>日期</td>
                  <td>是</td>
                  <td>1～31</td>
                  <td>, - * ? /</td>
                </tr>
                <tr>
                  <td>月份</td>
                  <td>是</td>
                  <td>1～12 或 JAN～DEC</td>
                  <td>, - * /</td>
                </tr>
                <tr>
                  <td>星期</td>
                  <td>是</td>
                  <td>0～6 或 MON～SUN</td>
                  <td>, - * ?</td>
                </tr>
              </table>
            </div>
          </el-tab-pane>
        </el-tabs>
      </el-dialog>

      <el-dialog :visible.sync="nextRunsModalVisible" width="480px" class="next-runs-dialog" :show-close="true" custom-class="next-runs-dialog-wrap">
        <div slot="title" class="next-runs-header">
          <span class="next-runs-header-title">未来执行时间</span>
          <span class="next-runs-header-subtitle">未来十次执行时间为</span>
        </div>
        <div class="next-runs-content">
          <p v-if="nextRunsIsEveryType && nextRunsList.length" class="next-runs-every-tip">根据当前时间计算，可能与实际略有差异。</p>
          <div v-if="nextRunsList.length" class="next-runs-list-wrap">
            <div v-for="(item, index) in nextRunsList" :key="index" class="next-runs-item">
              <span class="next-runs-item-index">{{ index + 1 }}</span>
              <span class="next-runs-item-time">{{ item }}</span>
            </div>
          </div>
          <div v-else class="next-runs-empty">
            <span class="next-runs-empty-text">{{ nextRunsError || '暂无数据' }}</span>
          </div>
        </div>
      </el-dialog>
    </el-main>
  </el-container>
</template>

<script>
import taskSidebar from './sidebar'
import taskService from '../../api/task'
import notificationService from '../../api/notification'
import parser from 'cron-parser'

function formatCronDate (d) {
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  const h = String(d.getHours()).padStart(2, '0')
  const min = String(d.getMinutes()).padStart(2, '0')
  const s = String(d.getSeconds()).padStart(2, '0')
  return y + '-' + m + '-' + day + ' ' + h + ':' + min + ':' + s
}

// 解析 @every 1h30m 等 Go time.ParseDuration 格式，返回毫秒数
function parseEveryDuration (str) {
  let ms = 0
  const hMatch = str.match(/(\d+)h/)
  if (hMatch) ms += parseInt(hMatch[1], 10) * 3600000
  const mMatch = str.match(/(\d+)m/)
  if (mMatch) ms += parseInt(mMatch[1], 10) * 60000
  const sMatch = str.match(/(\d+)s/)
  if (sMatch) ms += parseInt(sMatch[1], 10) * 1000
  return ms
}

// 5 字段按 jakecoffman/cron 规范：秒 分 时 日 月，缺周则补 *
function normalizeCronSpec (spec) {
  const parts = spec.trim().split(/\s+/)
  if (parts.length === 5) {
    return parts.join(' ') + ' *'
  }
  return spec
}

export default {
  name: 'task-edit',
  data () {
    return {
      form: {
        id: '',
        name: '',
        tag: '',
        level: 1,
        dependency_status: 1,
        dependency_task_id: '',
        spec: '',
        protocol: 2,
        http_method: 1,
        command: '',
        host_id: '',
        timeout: 0,
        multi: 2,
        notify_status: 1,
        notify_type: 2,
        notify_receiver_id: '',
        notify_keyword: '',
        retry_times: 0,
        retry_interval: 0,
        remark: ''
      },
      cronSpecModalVisible: false,
      cronSpecActiveTab: 'basic',
      nextRunsModalVisible: false,
      nextRunsList: [],
      nextRunsError: '',
      nextRunsIsEveryType: false,
      dependencyTooltipContent: '1. 强依赖：主定时执行成功，才会运行子定时\n2. 弱依赖：无论主定时执行是否成功，都会运行子定时',
      formRules: {
        name: [
          { required: true, message: '请输入定时名称', trigger: 'blur' }
        ],
        spec: [
          { required: true, message: '请输入crontab表达式', trigger: 'blur' }
        ],
        command: [
          { required: true, message: '请输入命令', trigger: 'blur' }
        ],
        timeout: [
          { type: 'number', required: true, message: '请输入有效的定时超时时间', trigger: 'blur' }
        ],
        retry_times: [
          { type: 'number', required: true, message: '请输入有效的定时执行失败重试次数', trigger: 'blur' }
        ],
        retry_interval: [
          { type: 'number', required: true, message: '请输入有效的定时执行失败，重试间隔时间', trigger: 'blur' }
        ],
        notify_keyword: [
          { required: true, message: '请输入要匹配的定时执行输出关键字', trigger: 'blur' }
        ]
      },
      httpMethods: [
        {
          value: 1,
          label: 'get'
        },
        {
          value: 2,
          label: 'post'
        }
      ],
      protocolList: [
        {
          value: 1,
          label: 'http'
        },
        {
          value: 2,
          label: 'shell'
        }
      ],
      levelList: [
        {
          value: 1,
          label: '主定时'
        },
        {
          value: 2,
          label: '子定时'
        }
      ],
      dependencyStatusList: [
        {
          value: 1,
          label: '强依赖'
        },
        {
          value: 2,
          label: '弱依赖'
        }
      ],
      runStatusList: [
        {
          value: 2,
          label: '是'
        },
        {
          value: 1,
          label: '否'
        }
      ],
      notifyStatusList: [
        {
          value: 1,
          label: '不通知'
        },
        {
          value: 2,
          label: '失败通知'
        },
        {
          value: 3,
          label: '总是通知'
        },
        {
          value: 4,
          label: '关键字匹配通知'
        }
      ],
      notifyTypes: [
        {
          value: 2,
          label: '邮件'
        },
        {
          value: 4,
          label: 'WebHook'
        }
      ],
      hosts: [],
      mailUsers: [],
      selectedHosts: [],
      selectedMailNotifyIds: []
    }
  },
  computed: {
    commandPlaceholder () {
      if (this.form.protocol === 1) {
        return '请输入URL地址'
      }

      return '请输入shell命令'
    }
  },
  watch: {
    'form.protocol' (val) {
      if (val === 1) {
        this.selectedHosts = []
        this.form.host_id = ''
      }
    }
  },
  components: { taskSidebar },
  created () {
    const id = this.$route.params.id

    taskService.detail(id, (taskData, hosts) => {
      if (id && !taskData) {
        this.$message.error('数据不存在')
        this.cancel()
        return
      }
      this.hosts = (hosts || []).filter(item => item && item.id)
      if (!taskData) {
        return
      }
      this.form.id = taskData.id
      this.form.name = taskData.name
      this.form.tag = taskData.tag
      this.form.level = taskData.level
      if (taskData.dependency_status) {
        this.form.dependency_status = taskData.dependency_status
      }
      this.form.dependency_task_id = taskData.dependency_task_id
      this.form.spec = taskData.spec
      this.form.protocol = taskData.protocol
      if (taskData.http_method) {
        this.form.http_method = taskData.http_method
      }
      this.form.command = taskData.command
      this.form.timeout = taskData.timeout
      this.form.multi = taskData.multi ? 1 : 2
      this.form.notify_keyword = taskData.notify_keyword
      if (taskData.notify_status !== undefined) {
        this.form.notify_status = taskData.notify_status + 1
      }
      this.form.notify_receiver_id = taskData.notify_receiver_id || ''
      if (taskData.notify_type !== undefined && taskData.notify_type !== null) {
        this.form.notify_type = taskData.notify_type + 1
      }
      this.form.retry_times = taskData.retry_times
      this.form.retry_interval = taskData.retry_interval
      this.form.remark = taskData.remark
      const tempHosts = []
      if (this.form.protocol === 2 && taskData.hosts) {
        taskData.hosts.forEach((v) => {
          if (v && v.host_id) {
            tempHosts.push(v.host_id)
          }
        })
      }
      this.selectedHosts = tempHosts

      const tempMailIds = []
      if (this.form.notify_status > 1 && this.form.notify_receiver_id) {
        const notifyReceiverIds = this.form.notify_receiver_id.split(',')
        if (this.form.notify_type === 2) {
          notifyReceiverIds.forEach((v) => {
            const val = parseInt(v)
            if (!isNaN(val)) tempMailIds.push(val)
          })
        }
      }
      this.selectedMailNotifyIds = tempMailIds
    })

    notificationService.mail((data) => {
      this.mailUsers = ((data && data.mail_users) || []).filter(item => item && item.id)
    })
  },
  methods: {
    copyCronExpr (expr) {
      const text = (expr || '').trim()
      if (!text) {
        this.$message.warning('可复制内容为空')
        return
      }
      const onSuccess = () => {
        this.$message.success('复制成功')
      }
      const onError = () => {
        this.$message.error('复制失败，请手动复制')
      }

      if (navigator.clipboard && navigator.clipboard.writeText) {
        navigator.clipboard.writeText(text).then(onSuccess).catch(() => {
          try {
            const input = document.createElement('textarea')
            input.value = text
            input.setAttribute('readonly', 'readonly')
            input.style.position = 'fixed'
            input.style.top = '-9999px'
            document.body.appendChild(input)
            input.select()
            const copied = document.execCommand('copy')
            document.body.removeChild(input)
            copied ? onSuccess() : onError()
          } catch (e) {
            onError()
          }
        })
        return
      }

      try {
        const input = document.createElement('textarea')
        input.value = text
        input.setAttribute('readonly', 'readonly')
        input.style.position = 'fixed'
        input.style.top = '-9999px'
        document.body.appendChild(input)
        input.select()
        const copied = document.execCommand('copy')
        document.body.removeChild(input)
        copied ? onSuccess() : onError()
      } catch (e) {
        onError()
      }
    },
    submit () {
      this.$refs['form'].validate((valid) => {
        if (!valid) {
          return false
        }
        if (this.form.protocol === 2 && this.selectedHosts.length === 0) {
          this.$message.error('请选择定时节点')
          return false
        }
        if (this.form.notify_status > 1) {
          if (this.form.notify_type === 2 && this.selectedMailNotifyIds.length === 0) {
            this.$message.error('请选择邮件接收用户')
            return false
          }
        }

        this.save()
      })
    },
    save () {
      if (this.form.protocol === 2 && this.selectedHosts.length > 0) {
        this.form.host_id = this.selectedHosts.filter(id => id).join(',')
      } else {
        this.form.host_id = ''
      }
      if (this.form.notify_status > 1 && this.form.notify_type === 2) {
        this.form.notify_receiver_id = this.selectedMailNotifyIds.filter(id => id).join(',')
      }
      taskService.update(this.form, () => {
        window.close()
      })
    },
    cancel () {
      window.close()
    },
    showCronSpecModal () {
      this.cronSpecActiveTab = 'basic'
      this.cronSpecModalVisible = true
    },
    showNextRunsModal () {
      const spec = (this.form.spec || '').trim()
      if (!spec) {
        this.$message.warning('请先输入 cron 表达式')
        return
      }
      this.nextRunsList = []
      this.nextRunsError = ''
      this.nextRunsIsEveryType = false
      this.nextRunsModalVisible = true
      try {
        const everyMatch = spec.match(/^@every\s+(.+)$/i)
        if (everyMatch) {
          this.nextRunsIsEveryType = true
          const durationMs = parseEveryDuration(everyMatch[1])
          if (durationMs <= 0) {
            this.nextRunsError = '无效的间隔时长'
            return
          }
          const now = Date.now()
          for (let i = 1; i <= 10; i++) {
            const t = new Date(now + durationMs * i)
            this.nextRunsList.push(formatCronDate(t))
          }
          return
        }
        const normalizedSpec = normalizeCronSpec(spec)
        const interval = parser.parseExpression(normalizedSpec, { currentDate: new Date() })
        const dates = interval.iterate(10)
        this.nextRunsList = dates.map(d => formatCronDate(d))
      } catch (err) {
        this.nextRunsError = (err && err.message) ? err.message : 'cron 表达式解析失败'
      }
    }
  }
}
</script>

<style scoped>
.el-main {
  background: #f5f5f7;
  padding: 24px 40px;
  min-height: calc(100vh - 52px);
}

/* ─── Page Header ─── */
.page-header {
  margin-bottom: 24px;
  display: flex;
  align-items: center;
  gap: 16px;
}

.back-btn {
  font-size: 14px;
  color: #007aff;
  padding: 0;
}

.page-title {
  font-size: 24px;
  font-weight: 700;
  color: #1d1d1f;
  margin: 0;
  letter-spacing: -0.5px;
}

/* ─── Form Sections ─── */
.ios-form {
  max-width: 100%;
}

.form-section {
  background: #ffffff;
  border-radius: 16px;
  padding: 20px 24px;
  margin-bottom: 16px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
  border: 1px solid rgba(0, 0, 0, 0.05);
}

.section-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 10px;
  padding-bottom: 8px;
  border-bottom: 1px solid #f2f2f7;
}

.section-header i {
  font-size: 18px;
  color: #007aff;
}

.section-header span {
  font-size: 16px;
  font-weight: 600;
  color: #1d1d1f;
}

/* ─── Element Overrides ─── */
.el-form-item>>>.el-form-item__label {
  font-weight: 600;
  color: #1d1d1f;
  padding-bottom: 4px;
  line-height: 1.5;
}

.el-input>>>.el-input__inner,
.el-textarea>>>.el-textarea__inner,
.el-select>>>.el-input__inner {
  background: #f5f5f7 !important;
  border: 1px solid transparent !important;
  border-radius: 10px !important;
  color: #1d1d1f !important;
  transition: all 0.2s ease;
}

.el-input>>>.el-input__inner:focus,
.el-textarea>>>.el-textarea__inner:focus,
.el-select>>>.el-input__inner:focus {
  border-color: #007aff !important;
  box-shadow: 0 0 0 4px rgba(0, 122, 255, 0.1) !important;
}

.item-tip {
  font-size: 12px;
  color: #86868b;
  margin-top: 4px;
}

/* ─── Cron 表达式表单项 ─── */
.cron-spec-form-item>>>.el-form-item__label {
  display: flex;
  align-items: center;
  flex-wrap: nowrap;
  white-space: nowrap;
}

.cron-label-wrap {
  display: inline-flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: nowrap;
  white-space: nowrap;
  gap: 4px;
}

.cron-spec-form-item .cron-spec-tooltip-icon {
  cursor: pointer;
}

.cron-spec-form-item .cron-spec-tooltip-icon:hover {
  color: #007aff;
}

.cron-next-runs-link {
  font-size: 12px;
  color: #007aff;
  cursor: pointer;
  padding-right: 8px;
  white-space: nowrap;
}

.cron-next-runs-link:hover {
  text-decoration: underline;
}

.label-with-tooltip .label-tooltip-icon {
  margin-left: 4px;
  font-size: 14px;
  color: #86868b;
  cursor: help;
  vertical-align: middle;
}

.label-with-tooltip .label-tooltip-icon:hover {
  color: #007aff;
}

.spec-tips {
  margin-top: 8px;
  padding: 12px;
  background: rgba(0, 122, 255, 0.04);
  border-radius: 8px;
  font-size: 12px;
  color: #1d1d1f;
}

.spec-tips p {
  margin: 4px 0;
}

.spec-tips strong {
  display: inline-block;
  width: 50px;
  color: #007aff;
}

/* ─── Code Editor Style Area ─── */
.code-editor>>>.el-textarea__inner {
  font-family: "SF Mono", "Monaco", "Consolas", monospace;
  font-size: 13px;
  line-height: 1.6;
  background: #1c1c1e !important;
  color: #ffffff !important;
  padding: 12px;
}

/* ─── Form Actions ─── */
.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 32px;
  padding-top: 24px;
  border-top: 1px solid #d2d2d7;
}

.form-actions .el-button {
  padding: 12px 28px;
  border-radius: 12px;
  font-weight: 600;
  transition: all 0.2s ease;
}

.form-actions .el-button--primary {
  background: #007aff;
  box-shadow: 0 4px 12px rgba(0, 122, 255, 0.25);
}

.form-actions .el-button--primary:hover {
  background: #0071e3;
  transform: translateY(-1px);
}
</style>

<style>
/* ─── Cron 表达式说明模态框 iOS 风格 ─── */
.cron-spec-dialog-wrap {
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 12px 48px rgba(0, 0, 0, 0.12);
}

.cron-spec-dialog-wrap .el-dialog__header {
  padding: 20px 24px 16px;
  border-bottom: none;
}

.cron-spec-dialog-wrap .el-dialog__body {
  padding: 0 24px 24px;
}

.cron-spec-dialog-wrap .el-dialog__headerbtn {
  top: 18px;
  right: 20px;
  width: 28px;
  height: 28px;
  font-size: 18px;
  color: #86868b;
}

.cron-spec-dialog-wrap .el-dialog__headerbtn:hover {
  color: #1d1d1f;
}

.cron-spec-header-title {
  font-size: 20px;
  font-weight: 700;
  color: #1d1d1f;
  letter-spacing: -0.5px;
}

.cron-spec-dialog-wrap .cron-spec-tabs .el-tabs__header {
  margin: 0 0 16px;
}

.cron-spec-dialog-wrap .cron-spec-tabs .el-tabs__nav-wrap::after {
  display: none;
}

.cron-spec-dialog-wrap .cron-spec-tabs .el-tabs__item {
  font-size: 14px;
  color: #86868b;
}

.cron-spec-dialog-wrap .cron-spec-tabs .el-tabs__item.is-active {
  color: #007aff;
  font-weight: 600;
}

.cron-spec-dialog-wrap .cron-spec-tabs .el-tabs__item:hover {
  color: #1d1d1f;
}

.cron-spec-dialog-wrap .cron-spec-tabs .el-tabs__ink-bar {
  background-color: #007aff;
}

.cron-spec-dialog-wrap .cron-spec-tabs .el-tabs__nav {
  border: none;
}

.cron-spec-dialog-wrap .cron-spec-tabs .el-tabs__active-bar {
  height: 3px;
  border-radius: 2px;
}

.cron-spec-dialog-wrap .cron-spec-tabs .el-tabs__content {
  overflow: visible;
}

.cron-spec-dialog-wrap .cron-spec-content {
  padding: 4px 0;
  font-size: 14px;
  color: #1d1d1f;
  line-height: 1.6;
}

.cron-spec-dialog-wrap .cron-spec-content h4 {
  margin: 20px 0 10px;
  font-size: 15px;
  font-weight: 600;
  color: #1d1d1f;
}

.cron-spec-dialog-wrap .cron-spec-content h4:first-child {
  margin-top: 0;
}

.cron-spec-dialog-wrap .cron-spec-content ul {
  margin: 0;
  padding-left: 24px;
}

.cron-spec-dialog-wrap .cron-spec-content li {
  margin: 6px 0;
}

.cron-spec-dialog-wrap .cron-spec-content code {
  padding: 3px 8px;
  background: #f2f2f7;
  border-radius: 6px;
  font-family: 'SF Mono', Monaco, Consolas, monospace;
  font-size: 12px;
  white-space: nowrap;
  color: #1d1d1f;
}

.cron-spec-dialog-wrap .cron-spec-content p {
  margin: 8px 0;
}

.cron-spec-dialog-wrap .cron-spec-table {
  width: 100%;
  border-collapse: collapse;
  margin: 12px 0;
  font-size: 14px;
  border-radius: 12px;
  overflow: hidden;
  border: 1px solid #e5e5ea;
}

.cron-spec-dialog-wrap .cron-spec-table th,
.cron-spec-dialog-wrap .cron-spec-table td {
  padding: 12px 16px;
  border: none;
  border-bottom: 1px solid #e5e5ea;
  text-align: left;
}

.cron-spec-dialog-wrap .cron-spec-table tr:last-child td {
  border-bottom: none;
}

.cron-spec-dialog-wrap .cron-spec-table th {
  background: #f5f5f7;
  font-weight: 600;
  font-size: 13px;
  color: #1d1d1f;
  white-space: nowrap;
}

.cron-spec-dialog-wrap .cron-spec-table td {
  background: #ffffff;
}

.cron-spec-dialog-wrap .cron-spec-table td {
  color: #1d1d1f;
}

.cron-spec-dialog-wrap .cron-spec-table td code {
  padding: 2px 6px;
}

.cron-spec-dialog-wrap .cron-spec-copy-wrap {
  display: inline-flex;
  align-items: center;
  gap: 8px;
}

.cron-spec-dialog-wrap .cron-spec-copy-btn {
  border: 1px solid #d2d2d7;
  border-radius: 6px;
  background: #ffffff;
  color: #1d1d1f;
  font-size: 12px;
  line-height: 1;
  padding: 4px 8px;
  cursor: pointer;
  opacity: 0;
  visibility: hidden;
  transition: opacity 0.2s ease, visibility 0.2s ease, border-color 0.2s ease;
}

.cron-spec-dialog-wrap .cron-spec-table td:hover .cron-spec-copy-btn,
.cron-spec-dialog-wrap .cron-spec-copy-wrap:focus-within .cron-spec-copy-btn {
  opacity: 1;
  visibility: visible;
}

.cron-spec-dialog-wrap .cron-spec-copy-btn:hover {
  border-color: #007aff;
  color: #007aff;
}

.cron-spec-dialog-wrap .cron-spec-table td:nth-child(4) {
  white-space: nowrap;
}

.cron-spec-dialog-wrap .cron-spec-content ul li {
  white-space: nowrap;
}

/* ─── 未来执行时间模态框 iOS 风格 ─── */
.next-runs-dialog-wrap {
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 12px 48px rgba(0, 0, 0, 0.12);
}

.next-runs-dialog-wrap .el-dialog__header {
  padding: 20px 24px 12px;
  border-bottom: none;
}

.next-runs-dialog-wrap .el-dialog__body {
  padding: 0 24px 24px;
}

.next-runs-dialog-wrap .el-dialog__headerbtn {
  top: 18px;
  right: 20px;
  width: 28px;
  height: 28px;
  font-size: 18px;
  color: #86868b;
}

.next-runs-dialog-wrap .el-dialog__headerbtn:hover {
  color: #1d1d1f;
}

.next-runs-header {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.next-runs-header-title {
  font-size: 20px;
  font-weight: 700;
  color: #1d1d1f;
  letter-spacing: -0.5px;
}

.next-runs-header-subtitle {
  font-size: 14px;
  color: #86868b;
  font-weight: 400;
}

.next-runs-content {
  margin-top: 8px;
}

.next-runs-every-tip {
  margin: 0 0 12px;
  font-size: 13px;
  color: #86868b;
  line-height: 1.5;
}

.next-runs-list-wrap {
  background: #f5f5f7;
  border-radius: 12px;
  overflow: hidden;
}

.next-runs-item {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 14px 16px;
  background: #ffffff;
  border-bottom: 1px solid #f2f2f7;
  font-size: 15px;
  color: #1d1d1f;
}

.next-runs-item:last-child {
  border-bottom: none;
}

.next-runs-item-index {
  flex-shrink: 0;
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f2f2f7;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 600;
  color: #86868b;
}

.next-runs-item-time {
  font-family: "SF Mono", "Monaco", "Consolas", monospace;
  font-size: 14px;
  color: #1d1d1f;
}

.next-runs-empty {
  padding: 40px 24px;
  text-align: center;
  background: #f5f5f7;
  border-radius: 12px;
}

.next-runs-empty-text {
  font-size: 14px;
  color: #86868b;
}

.label-tooltip-popper {
  max-width: 280px;
  line-height: 1.6;
  white-space: pre-line;
}
</style>
