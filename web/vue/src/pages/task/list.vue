<template>
  <el-container>
    <task-sidebar></task-sidebar>
    <el-main class="ios-task-page">
      <div class="search-bar">
        <div class="search-filters">
          <el-input v-model.trim="searchParams.keyword" placeholder="定时ID/定时名称" class="search-input" @input="debouncedSearch" @keyup.native.enter="search()">
          </el-input>
          <el-select v-model.trim="searchParams.tag" class="search-select" filterable clearable placeholder="标签" @change="debouncedSearch">
            <el-option v-for="item in tagOptions" :key="item" :label="item" :value="item"></el-option>
          </el-select>
          <el-select v-model.trim="searchParams.protocol" class="search-select" @change="debouncedSearch">
            <el-option label="全部执行方式" value=""></el-option>
            <el-option v-for="item in protocolList" :key="item.value" :label="item.label" :value="item.value">
            </el-option>
          </el-select>
          <el-select v-model.trim="searchParams.host_id" class="search-select" @change="debouncedSearch">
            <el-option label="全部定时节点" value=""></el-option>
            <el-option v-for="item in hosts" :key="item.id" :label="item.alias + ' - ' + item.name + ':' + item.port" :value="item.id">
            </el-option>
          </el-select>
          <el-select v-model.trim="searchParams.status" class="search-select" @change="debouncedSearch">
            <el-option label="全部状态" value=""></el-option>
            <el-option v-for="item in statusList" :key="item.value" :label="item.label" :value="item.value">
            </el-option>
          </el-select>
        </div>
        <div class="search-actions">
          <el-button class="ios-action-btn" type="primary" plain @click="toEdit(null)" v-if="this.$store.getters.user.isAdmin">新增定时</el-button>
          <el-button class="ios-action-btn" type="primary" plain @click="refresh">刷新</el-button>
        </div>
      </div>
      <div class="pagination-bar">
        <span class="pagination-total">共 {{ taskTotal }} 条</span>
        <el-pagination class="ios-pagination" background layout="prev, pager, next, jumper, sizes" :total="taskTotal" :current-page="searchParams.page" :page-size="searchParams.page_size" :page-sizes="[10, 20, 50, 100]" @size-change="changePageSize" @current-change="changePage" @prev-click="changePage" @next-click="changePage">
        </el-pagination>
      </div>
      <div class="ios-table-card">
        <el-table v-loading="loading" :data="tasks" tooltip-effect="dark" style="width: 100%" class="task-list-table ios-data-table">
          <el-table-column label="定时名称" min-width="180">
            <template slot-scope="scope">
              <div class="task-name-cell">
                <div class="task-name-main">
                  <span>{{ scope.row.name }}</span>
                  <el-tooltip v-if="scope.row.tag" content="点击复制标签" placement="top">
                    <span class="ios-tag ios-tag-clickable" :style="getTagStyle(scope.row.tag)" role="button" tabindex="0" @click.stop="copyTag(scope.row.tag)" @keydown.enter.prevent="copyTag(scope.row.tag)" @keydown.space.prevent="copyTag(scope.row.tag)">{{ scope.row.tag }}</span>
                  </el-tooltip>
                </div>
                <div class="task-meta-sub">
                  <el-tooltip content="创建时间" placement="top">
                    <span class="task-time-sub">{{ formatTaskCreated(scope.row.created) }}</span>
                  </el-tooltip>
                  <el-tooltip content="定时ID" placement="top">
                    <span class="task-id-sub">{{ scope.row.id }}</span>
                  </el-tooltip>
                </div>
              </div>
            </template>
          </el-table-column>
          <el-table-column label="cron表达式" width="200">
            <template slot-scope="scope">
              <div class="cron-cell">
                <div class="cron-spec-wrap">
                  <el-tooltip content="点击查看未来 10 次执行时间" placement="top">
                    <span class="cron-spec-text" role="button" tabindex="0" @click="showCronRunsModal(scope.row)" @keydown.enter.prevent="showCronRunsModal(scope.row)" @keydown.space.prevent="showCronRunsModal(scope.row)">{{ scope.row.spec }}</span>
                  </el-tooltip>
                  <button type="button" class="cron-copy-btn" @click.stop="copyCronSpec(scope.row.spec)">复制</button>
                </div>
                <el-tooltip v-if="!isEverySpec(scope.row.spec) && isValidNextRunTime(scope.row.next_run_time)" content="下次执行时间" placement="top">
                  <span class="ios-tag ios-tag-time">{{ scope.row.next_run_time | formatTime }}</span>
                </el-tooltip>
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="protocol" :formatter="formatProtocol" width="90" label="执行方式">
          </el-table-column>
          <el-table-column label="状态" width="70" align="center" v-if="this.isAdmin">
            <template slot-scope="scope">
              <el-switch v-if="scope.row.level === 1" v-model="scope.row.status" :active-value="1" :inactive-vlaue="0" active-color="#13ce66" @change="changeStatus(scope.row)" inactive-color="#ff4949">
              </el-switch>
            </template>
          </el-table-column>
          <el-table-column label="状态" width="70" align="center" v-else>
            <template slot-scope="scope">
              <el-switch v-if="scope.row.level === 1" v-model="scope.row.status" :active-value="1" :inactive-vlaue="0" active-color="#13ce66" :disabled="true" inactive-color="#ff4949">
              </el-switch>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="360" v-if="this.isAdmin">
            <template slot-scope="scope">
              <div class="action-buttons">
                <el-button type="info" size="mini" @click="openDetail(scope.row)">详情</el-button>
                <el-button type="primary" size="mini" @click="toEdit(scope.row)">编辑</el-button>
                <el-button type="success" size="mini" @click="runTask(scope.row)">手动执行</el-button>
                <el-button type="info" size="mini" @click="jumpToLog(scope.row)">日志</el-button>
                <el-button type="danger" size="mini" @click="remove(scope.row)">删除</el-button>
              </div>
            </template>
          </el-table-column>
        </el-table>
      </div>
      <div class="pagination-bar pagination-bar-bottom">
        <span class="pagination-total">共 {{ taskTotal }} 条</span>
        <el-pagination class="ios-pagination" background layout="prev, pager, next, jumper, sizes" :total="taskTotal" :current-page="searchParams.page" :page-size="searchParams.page_size" :page-sizes="[10, 20, 50, 100]" @size-change="changePageSize" @current-change="changePage" @prev-click="changePage" @next-click="changePage">
        </el-pagination>
      </div>

      <!-- 定时详情模态窗 -->
      <el-dialog :title="'定时详情 - ' + (detailTask ? detailTask.name : '')" :visible.sync="detailDialogVisible" width="820px" custom-class="task-detail-dialog" :append-to-body="true" @close="closeDetailModal">
        <div class="task-expand-wrap" v-if="detailTask">
          <div class="task-expand-section task-expand-grid">
            <div class="task-expand-row" v-if="detailTask.created">
              <span class="task-expand-label">定时创建时间</span>
              <span class="task-expand-value">{{ detailTask.created | formatTime }}</span>
            </div>
            <div class="task-expand-row">
              <span class="task-expand-label">定时类型</span>
              <span class="task-expand-value">{{ detailTask.level | formatLevel }}</span>
            </div>
            <div class="task-expand-row">
              <span class="task-expand-label">单实例运行</span>
              <span class="task-expand-value">{{ detailTask.multi | formatMulti }}</span>
            </div>
            <div class="task-expand-row">
              <span class="task-expand-label">超时时间</span>
              <span class="task-expand-value">{{ detailTask.timeout | formatTimeout }}</span>
            </div>
            <div class="task-expand-row">
              <span class="task-expand-label">重试次数</span>
              <span class="task-expand-value">{{ detailTask.retry_times }}</span>
            </div>
            <div class="task-expand-row">
              <span class="task-expand-label">重试间隔</span>
              <span class="task-expand-value">{{ detailTask.retry_interval | formatRetryTimesInterval }}</span>
            </div>
          </div>
          <div class="task-expand-section" v-if="detailTask.protocol === 2 && (detailTask.hosts || []).length">
            <div class="task-expand-row">
              <span class="task-expand-label">定时节点</span>
              <span class="task-expand-value task-expand-value-block">
                <span v-for="item in (detailTask.hosts || [])" :key="item.host_id">{{ item.alias }} - {{ item.name }}:{{ item.port }}</span>
              </span>
            </div>
          </div>
          <div class="task-expand-section" v-if="detailTask.command">
            <div class="task-expand-row task-expand-row-full">
              <span class="task-expand-label">命令</span>
              <span class="task-expand-value task-expand-value-code">{{ detailTask.command }}</span>
            </div>
          </div>
          <div class="task-expand-section" v-if="detailTask.remark">
            <div class="task-expand-row task-expand-row-full">
              <span class="task-expand-label">备注</span>
              <span class="task-expand-value">{{ detailTask.remark }}</span>
            </div>
          </div>
        </div>
      </el-dialog>

      <!-- 定时日志模态窗 -->
      <el-dialog :title="'定时日志 - ' + (logModalTask ? logModalTask.name : '')" :visible.sync="logDialogVisible" width="1000px" custom-class="log-modal-dialog" :append-to-body="true" @close="closeLogModal">
        <div class="log-modal-content" v-if="logModalTask">
          <div class="log-modal-toolbar">
            <el-form :inline="true" size="small" @submit.native.prevent>
              <el-form-item>
                <el-select v-model="logSearchParams.status" placeholder="全部状态" clearable @change="searchLogs">
                  <el-option label="全部状态" value=""></el-option>
                  <el-option v-for="item in logStatusList" :key="item.value" :label="item.label" :value="item.value"></el-option>
                </el-select>
              </el-form-item>
              <el-form-item>
                <el-button type="info" size="small" @click="searchLogs">
                  <i class="el-icon-refresh" :class="{ 'is-loading': logLoading }"></i>
                  刷新
                </el-button>
              </el-form-item>
              <el-form-item>
                <el-button type="text" size="small" @click="openLogInNewWindow">新窗口打开</el-button>
              </el-form-item>
            </el-form>
          </div>
          <div class="log-modal-pagination">
            <span class="pagination-total">共 {{ logTotal }} 条</span>
            <el-pagination small background layout="prev, pager, next, sizes" :total="logTotal" :page-sizes="[10, 20, 50]" :page-size="logSearchParams.page_size" :current-page="logSearchParams.page" @size-change="changeLogPageSize" @current-change="changeLogPage">
            </el-pagination>
          </div>
          <el-table v-loading="logLoading" :data="logList" border size="small" class="log-modal-table" max-height="400">
            <el-table-column prop="id" label="日志ID" width="80"></el-table-column>
            <el-table-column label="执行时间" min-width="200">
              <template slot-scope="scope">
                <div class="time-info">
                  <span>{{ scope.row.start_time | formatTime }}</span>
                  <span v-if="scope.row.status !== 1">~ {{ scope.row.end_time | formatTime }}</span>
                  <div class="duration">耗时: {{ scope.row.total_time > 0 ? scope.row.total_time : 1 }}s</div>
                </div>
              </template>
            </el-table-column>
            <el-table-column label="状态" width="90">
              <template slot-scope="scope">
                <el-tag v-if="scope.row.status === 0" type="danger" size="mini">失败</el-tag>
                <el-tag v-else-if="scope.row.status === 1" type="primary" size="mini">执行</el-tag>
                <el-tag v-else-if="scope.row.status === 2" type="success" size="mini">成功</el-tag>
                <el-tag v-else-if="scope.row.status === 3" type="info" size="mini">取消</el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="240" align="center">
              <template slot-scope="scope">
                <div class="action-buttons" style="justify-content: center;">
                  <el-button size="mini" type="info" plain @click="showLogDetail(scope.row)">详情</el-button>
                  <el-button size="mini" type="primary" plain v-if="scope.row.status === 2 || scope.row.status === 0" @click="showLogResult(scope.row)">结果</el-button>
                  <el-button size="mini" type="danger" plain v-if="scope.row.status === 1 && scope.row.protocol === 2" @click="stopLogTask(scope.row)">停止定时</el-button>
                </div>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </el-dialog>

      <!-- 日志执行结果弹窗 -->
      <el-dialog title="定时执行详情" :visible.sync="logResultDialogVisible" width="800px" custom-class="result-dialog ios-result-dialog" :append-to-body="true">
        <div class="result-content">
          <div class="result-section">
            <div class="section-label">执行命令</div>
            <div class="code-container command">
              <code>{{ currentLogResult.command }}</code>
            </div>
          </div>
          <div class="result-section">
            <div class="section-label">输出结果</div>
            <div class="code-container output">
              <pre>{{ currentLogResult.result || '无输出结果' }}</pre>
            </div>
          </div>
        </div>
      </el-dialog>
      <el-dialog title="日志详情" :visible.sync="logDetailDialogVisible" width="760px" custom-class="task-detail-dialog" :append-to-body="true" @close="closeLogDetailDialog">
        <div class="task-expand-wrap" v-if="currentLogDetail">
          <div class="task-expand-section task-expand-grid">
            <div class="task-expand-row">
              <span class="task-expand-label">日志ID</span>
              <span class="task-expand-value">{{ currentLogDetail.id }}</span>
            </div>
            <div class="task-expand-row">
              <span class="task-expand-label">定时ID</span>
              <span class="task-expand-value">{{ currentLogDetail.task_id }}</span>
            </div>
            <div class="task-expand-row">
              <span class="task-expand-label">重试次数</span>
              <span class="task-expand-value">{{ currentLogDetail.retry_times }}</span>
            </div>
          </div>
          <div class="task-expand-section">
            <div class="task-expand-row task-expand-row-full">
              <span class="task-expand-label">Cron表达式</span>
              <span class="task-expand-value task-expand-value-code">{{ currentLogDetail.spec || '-' }}</span>
            </div>
          </div>
          <div class="task-expand-section">
            <div class="task-expand-row task-expand-row-full">
              <span class="task-expand-label">执行命令</span>
              <span class="task-expand-value task-expand-value-code">{{ currentLogDetail.command || '-' }}</span>
            </div>
          </div>
        </div>
      </el-dialog>

      <el-dialog :title="'未来执行时间 - ' + cronRunsTitle" :visible.sync="cronRunsModalVisible" width="520px" custom-class="task-detail-dialog" :append-to-body="true">
        <div class="cron-runs-content">
          <p v-if="cronRunsIsEveryType" class="cron-runs-tip">根据当前时间计算，可能与实际执行时间略有差异。</p>
          <div v-if="cronRunsList.length" class="cron-runs-list">
            <div v-for="(item, index) in cronRunsList" :key="index" class="cron-runs-item">
              <span class="cron-runs-index">{{ index + 1 }}</span>
              <span class="cron-runs-time">{{ item }}</span>
            </div>
          </div>
          <div v-else class="cron-runs-empty">{{ cronRunsError || '暂无数据' }}</div>
        </div>
      </el-dialog>
    </el-main>
  </el-container>
</template>

<script>
import taskSidebar from './sidebar'
import taskService from '../../api/task'
import taskLogService from '../../api/taskLog'
import parser from 'cron-parser'
import '../../styles/task-log.css'
import '../../styles/ios-table.css'

function formatCronDate (d) {
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  const h = String(d.getHours()).padStart(2, '0')
  const min = String(d.getMinutes()).padStart(2, '0')
  const s = String(d.getSeconds()).padStart(2, '0')
  return y + '-' + m + '-' + day + ' ' + h + ':' + min + ':' + s
}

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

function normalizeCronSpec (spec) {
  const parts = spec.trim().split(/\s+/)
  if (parts.length === 5) {
    return parts.join(' ') + ' *'
  }
  return spec
}

export default {
  name: 'task-list',
  data () {
    return {
      tasks: [],
      hosts: [],
      taskTotal: 0,
      loading: false,
      searchParams: {
        page_size: 20,
        page: 1,
        keyword: '',
        protocol: '',
        tag: '',
        host_id: '',
        status: ''
      },
      isAdmin: !!(this.$store.getters.user && this.$store.getters.user.isAdmin),
      protocolList: [
        {
          value: '1',
          label: 'http'
        },
        {
          value: '2',
          label: 'shell'
        }
      ],
      statusList: [
        {
          value: '2',
          label: '激活'
        },
        {
          value: '1',
          label: '停止'
        }
      ],
      detailDialogVisible: false,
      detailTask: null,
      logDialogVisible: false,
      logModalTask: null,
      logList: [],
      logTotal: 0,
      logLoading: false,
      logSearchParams: {
        page_size: 20,
        page: 1,
        task_id: '',
        status: ''
      },
      logStatusList: [
        { value: '2', label: '执行' },
        { value: '3', label: '成功' },
        { value: '1', label: '失败' },
        { value: '4', label: '取消' }
      ],
      cronRunsModalVisible: false,
      cronRunsList: [],
      cronRunsError: '',
      cronRunsTitle: '',
      cronRunsIsEveryType: false,
      logResultDialogVisible: false,
      currentLogResult: { command: '', result: '' },
      logDetailDialogVisible: false,
      currentLogDetail: null,
      tagOptions: [],
      tagLoading: false,
      searchTimer: null
    }
  },
  components: { taskSidebar },
  created () {
    const hostId = this.$route.query.host_id
    if (hostId) {
      this.searchParams.host_id = hostId
    }
    this.loadTags()
    this.search()
  },
  beforeDestroy () {
    if (this.searchTimer) clearTimeout(this.searchTimer)
  },
  filters: {
    formatLevel (value) {
      if (value === 1) {
        return '主定时'
      }
      return '子定时'
    },
    formatTimeout (value) {
      if (value > 0) {
        return value + '秒'
      }
      return '不限制'
    },
    formatRetryTimesInterval (value) {
      if (value > 0) {
        return value + '秒'
      }
      return '系统默认'
    },
    formatMulti (value) {
      if (value > 0) {
        return '否'
      }
      return '是'
    }
  },
  methods: {
    copyTag (tag) {
      const text = (tag || '').trim()
      if (!text) {
        this.$message.warning('标签内容为空')
        return
      }
      const onSuccess = () => this.$message.success('标签已复制')
      const onError = () => this.$message.error('复制失败，请手动复制')

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
    copyCronSpec (spec) {
      const text = (spec || '').trim()
      if (!text) {
        this.$message.warning('Cron 表达式为空')
        return
      }
      const onSuccess = () => this.$message.success('Cron 表达式已复制')
      const onError = () => this.$message.error('复制失败，请手动复制')

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
    showCronRunsModal (row) {
      const spec = row && row.spec ? String(row.spec).trim() : ''
      if (!spec) {
        this.$message.warning('Cron 表达式为空')
        return
      }

      this.cronRunsTitle = row && row.name ? row.name : '定时'
      this.cronRunsList = []
      this.cronRunsError = ''
      this.cronRunsIsEveryType = false
      this.cronRunsModalVisible = true

      try {
        const everyMatch = spec.match(/^@every\s+(.+)$/i)
        if (everyMatch) {
          this.cronRunsIsEveryType = true
          const durationMs = parseEveryDuration(everyMatch[1])
          if (durationMs <= 0) {
            this.cronRunsError = '无效的间隔时长'
            return
          }
          const now = Date.now()
          for (let i = 1; i <= 10; i++) {
            const t = new Date(now + durationMs * i)
            this.cronRunsList.push(formatCronDate(t))
          }
          return
        }

        const normalizedSpec = normalizeCronSpec(spec)
        const interval = parser.parseExpression(normalizedSpec, { currentDate: new Date() })
        const dates = interval.iterate(10)
        this.cronRunsList = dates.map(d => formatCronDate(d))
      } catch (err) {
        this.cronRunsError = (err && err.message) ? err.message : 'cron 表达式解析失败'
      }
    },
    getTagStyle (tag) {
      const text = (tag || '').trim()
      if (!text) {
        return {}
      }
      let hash = 0
      for (let i = 0; i < text.length; i++) {
        hash = ((hash << 5) - hash) + text.charCodeAt(i)
        hash |= 0
      }
      const positiveHash = Math.abs(hash)
      const hue = positiveHash % 360
      const saturation = 58 + (positiveHash % 18)
      const textLightness = 32 + ((positiveHash >> 3) % 8)
      const bgLightness = 94 - ((positiveHash >> 5) % 6)

      return {
        background: 'hsl(' + hue + ', ' + saturation + '%, ' + bgLightness + '%)',
        color: 'hsl(' + hue + ', ' + saturation + '%, ' + textLightness + '%)',
        border: '1px solid hsl(' + hue + ', ' + saturation + '%, ' + (bgLightness - 7) + '%)'
      }
    },
    formatTaskCreated (created) {
      if (created === null || created === undefined) return '-'

      const raw = String(created).trim()
      if (!raw || raw === '0') return '-'

      let timestamp = created
      if (/^\d+$/.test(raw)) {
        const num = Number(raw)
        if (!num) return '-'
        // 后端可能返回秒级时间戳，这里统一转毫秒
        timestamp = num < 1000000000000 ? num * 1000 : num
      }

      const date = new Date(timestamp)
      if (isNaN(date.getTime()) || date.getFullYear() < 2000) return '-'

      const fillZero = (num) => (num >= 10 ? num : '0' + num)
      return date.getFullYear() + '-' +
        fillZero(date.getMonth() + 1) + '-' +
        fillZero(date.getDate()) + ' ' +
        fillZero(date.getHours()) + ':' +
        fillZero(date.getMinutes()) + ':' +
        fillZero(date.getSeconds())
    },
    isEverySpec (spec) {
      return spec && String(spec).trim().toLowerCase().startsWith('@every')
    },
    isValidNextRunTime (time) {
      if (!time) return false
      const str = String(time)
      if (str.indexOf('0001-01-01') === 0) return false
      const date = new Date(time)
      return date.getFullYear() >= 2000
    },
    debouncedSearch () {
      if (this.searchTimer) clearTimeout(this.searchTimer)
      this.searchTimer = setTimeout(() => {
        this.searchParams.page = 1
        this.search()
      }, 400)
    },
    loadTags () {
      this.tagLoading = true
      taskService.tags({}, (data) => {
        this.tagOptions = Array.isArray(data) ? data : []
        this.tagLoading = false
      }, () => {
        this.tagLoading = false
      })
    },
    changeStatus (item) {
      const msg = item.status ? '定时已启用' : '定时已停用'
      if (item.status) {
        taskService.enable(item.id, () => {
          this.$message.success(msg)
        })
      } else {
        taskService.disable(item.id, () => {
          this.$message.success(msg)
        })
      }
    },
    formatProtocol (row, col) {
      if (!row) return ''
      if (row[col.property] === 2) {
        return 'shell'
      }
      if (row.http_method === 1) {
        return 'http-get'
      }
      return 'http-post'
    },
    changePage (page) {
      this.searchParams.page = page
      this.search()
    },
    changePageSize (pageSize) {
      this.searchParams.page_size = pageSize
      this.searchParams.page = 1
      this.search()
    },
    search (callback = null) {
      this.loading = true
      taskService.list(this.searchParams, (tasks, hosts) => {
        this.tasks = (tasks && tasks.data) ? tasks.data : []
        this.taskTotal = (tasks && tasks.total != null) ? tasks.total : 0
        this.hosts = (hosts && Array.isArray(hosts)) ? hosts : []
        this.loading = false
        if (typeof callback === 'function') {
          callback()
        }
      }, () => {
        this.loading = false
      })
    },
    runTask (item) {
      this.$appConfirm(() => {
        taskService.run(item.id, () => {
          this.$message.success('定时已开始执行')
        })
      })
    },
    remove (item) {
      this.$appConfirm(() => {
        taskService.remove(item.id, () => {
          this.refresh()
        })
      })
    },
    openDetail (item) {
      this.detailTask = item
      this.detailDialogVisible = true
    },
    closeDetailModal () {
      this.detailTask = null
    },
    jumpToLog (item) {
      this.logModalTask = item
      this.logSearchParams = {
        page_size: 20,
        page: 1,
        task_id: String(item.id),
        status: ''
      }
      this.logDialogVisible = true
      this.searchLogs()
    },
    closeLogModal () {
      this.logModalTask = null
      this.logList = []
      this.logTotal = 0
    },
    searchLogs () {
      if (!this.logModalTask) return
      this.logLoading = true
      taskLogService.list(this.logSearchParams, (data) => {
        this.logList = data.data || []
        this.logTotal = data.total || 0
        this.logLoading = false
      }, () => {
        this.logLoading = false
      })
    },
    changeLogPage (page) {
      this.logSearchParams.page = page
      this.searchLogs()
    },
    changeLogPageSize (pageSize) {
      this.logSearchParams.page_size = pageSize
      this.logSearchParams.page = 1
      this.searchLogs()
    },
    showLogResult (row) {
      this.currentLogResult = { command: row.command || '', result: row.result || '' }
      this.logResultDialogVisible = true
    },
    showLogDetail (row) {
      this.currentLogDetail = { ...row }
      this.logDetailDialogVisible = true
    },
    closeLogDetailDialog () {
      this.currentLogDetail = null
    },
    stopLogTask (row) {
      taskLogService.stop(row.id, row.task_id, () => {
        this.searchLogs()
      })
    },
    openLogInNewWindow () {
      if (!this.logModalTask) return
      const route = this.$router.resolve({ path: '/task/log', query: { task_id: this.logModalTask.id } })
      window.open(route.href, '_blank')
    },
    refresh () {
      this.search(() => {
        this.$message.success('刷新成功')
      })
    },
    toEdit (item) {
      const url = item === null ? '/task/create' : `/task/edit/${item.id}`
      const route = this.$router.resolve(url)
      window.open(route.href, '_blank')
    }
  }
}
</script>

<style scoped></style>
