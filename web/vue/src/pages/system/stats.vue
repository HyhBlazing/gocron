<template>
  <el-container>
    <system-sidebar></system-sidebar>
    <el-main class="stats-page">
      <div class="stats-header">
        <div>
          <h2 class="stats-title">统计</h2>
          <p class="stats-subtitle">按时间范围聚合统计（默认最近3个月）</p>
        </div>
        <div class="stats-actions">
          <el-button class="ios-action-btn" type="primary" plain :disabled="loading" @click="handleRefresh">
            <i class="el-icon-refresh" :class="{ 'is-loading': loading }"></i>
            刷新
          </el-button>
        </div>
      </div>

      <div class="stats-overview">
        <div class="overview-item">
          <div class="overview-label">定时总数</div>
          <div class="overview-value">{{ overview.task_total }}</div>
        </div>
        <div class="overview-item">
          <div class="overview-label">定时日志</div>
          <div class="overview-value">{{ overview.task_log_total }}</div>
        </div>
        <div class="overview-item">
          <div class="overview-label">用户总数</div>
          <div class="overview-value">{{ overview.user_total }}</div>
        </div>
        <div class="overview-item">
          <div class="overview-label">登录日志</div>
          <div class="overview-value">{{ overview.login_log_total }}</div>
        </div>
      </div>

      <div class="stats-filter-bar">
        <span class="filter-label">时间筛选</span>
        <el-date-picker v-model="dateRange" type="daterange" range-separator="至" start-placeholder="开始日期" end-placeholder="结束日期" value-format="yyyy-MM-dd" :clearable="false" :picker-options="datePickerOptions" @change="handleDateRangeChange">
        </el-date-picker>
      </div>

      <section class="stats-section">
        <h3 class="section-title">所有定时</h3>
        <div class="chart-grid chart-grid-single">
          <div class="chart-card">
            <div class="card-title">标签分布</div>
            <div ref="chartTaskTag" class="chart-box"></div>
          </div>
        </div>
      </section>

      <section class="stats-section">
        <h3 class="section-title">定时日志</h3>
        <div class="chart-grid">
          <div class="chart-card">
            <div class="card-title">耗时分布</div>
            <div ref="chartTaskLogDuration" class="chart-box"></div>
          </div>
          <div class="chart-card">
            <div class="card-title">日志状态分布</div>
            <div ref="chartTaskLogStatus" class="chart-box"></div>
          </div>
        </div>
      </section>

      <section class="stats-section">
        <h3 class="section-title">用户与登录</h3>
        <div class="chart-grid">
          <div class="chart-card">
            <div class="card-title">用户登录次数</div>
            <div ref="chartLoginUser" class="chart-box"></div>
          </div>
          <div class="chart-card">
            <div class="card-title">登录IP分布</div>
            <div ref="chartLoginIp" class="chart-box"></div>
          </div>
        </div>
      </section>
    </el-main>
  </el-container>
</template>

<script>
import { use, init } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { BarChart, PieChart } from 'echarts/charts'
import { GridComponent, TooltipComponent, LegendComponent, TitleComponent, DataZoomComponent } from 'echarts/components'
import { UniversalTransition } from 'echarts/features'
import systemSidebar from './sidebar'
import systemService from '../../api/system'

use([
  CanvasRenderer,
  BarChart,
  PieChart,
  GridComponent,
  TooltipComponent,
  LegendComponent,
  TitleComponent,
  DataZoomComponent,
  UniversalTransition
])

const MAX_RANGE_DAYS = 365

function formatDate (date) {
  const y = date.getFullYear()
  const m = String(date.getMonth() + 1).padStart(2, '0')
  const d = String(date.getDate()).padStart(2, '0')
  return y + '-' + m + '-' + d
}

function getDefaultRange () {
  const end = new Date()
  const start = new Date()
  start.setDate(end.getDate() - 89)
  return [formatDate(start), formatDate(end)]
}

function topNFromList (list, topN, othersName) {
  const source = Array.isArray(list) ? list.slice() : []
  source.sort((a, b) => (b.value || 0) - (a.value || 0))
  if (source.length <= topN) return source
  const head = source.slice(0, topN)
  const others = source.slice(topN).reduce((sum, item) => sum + (item.value || 0), 0)
  if (others > 0) {
    head.push({ name: othersName || '其他', value: others })
  }
  return head
}

export default {
  name: 'system-stats',
  components: { systemSidebar },
  data () {
    return {
      loading: false,
      dateRange: getDefaultRange(),
      refreshTimer: null,
      datePickerOptions: {
        disabledDate (time) {
          return time.getTime() > Date.now()
        }
      },
      chartMap: {},
      overview: {
        task_total: 0,
        task_log_total: 0,
        user_total: 0,
        login_log_total: 0
      },
      taskTagChart: [],
      taskLogDurationChart: [],
      taskLogStatusChart: [],
      loginUserChart: [],
      loginIpChart: [],
      pieColors: ['#007aff', '#5e5ce6', '#34c759', '#ff9f0a', '#af52de', '#30b0c7', '#ff6b6b', '#6d6d72', '#00a9a5', '#4e7cff', '#b084f5', '#53b97c']
    }
  },
  mounted () {
    this.loadStats()
    this.onResize = () => this.resizeCharts()
    window.addEventListener('resize', this.onResize)
  },
  beforeDestroy () {
    window.removeEventListener('resize', this.onResize)
    if (this.refreshTimer) {
      clearTimeout(this.refreshTimer)
      this.refreshTimer = null
    }
    this.disposeCharts()
  },
  methods: {
    handleDateRangeChange (val) {
      if (!val || val.length !== 2) {
        this.dateRange = getDefaultRange()
        return
      }
      const start = new Date(val[0] + ' 00:00:00')
      const end = new Date(val[1] + ' 00:00:00')
      const diffDays = Math.floor((end - start) / (24 * 3600 * 1000)) + 1
      if (diffDays > MAX_RANGE_DAYS) {
        this.$message.warning('时间范围最多支持1年')
        const fixedEnd = new Date(val[1] + ' 00:00:00')
        const fixedStart = new Date(fixedEnd)
        fixedStart.setDate(fixedStart.getDate() - (MAX_RANGE_DAYS - 1))
        this.dateRange = [formatDate(fixedStart), formatDate(fixedEnd)]
      }
      this.loadStats()
    },
    handleRefresh () {
      if (this.loading) return
      if (this.refreshTimer) clearTimeout(this.refreshTimer)
      this.refreshTimer = setTimeout(() => this.loadStats(), 180)
    },
    initCharts () {
      const refs = ['chartTaskTag', 'chartTaskLogDuration', 'chartTaskLogStatus', 'chartLoginUser', 'chartLoginIp']
      refs.forEach((refName) => {
        const dom = this.$refs[refName]
        if (!dom) return
        if (!this.chartMap[refName]) {
          this.chartMap[refName] = init(dom)
        }
      })
    },
    disposeCharts () {
      for (const key in this.chartMap) {
        if (this.chartMap[key]) this.chartMap[key].dispose()
      }
      this.chartMap = {}
    },
    resizeCharts () {
      for (const key in this.chartMap) {
        if (this.chartMap[key]) this.chartMap[key].resize()
      }
    },
    setChartOption (refName, option) {
      if (!this.chartMap[refName]) return
      this.chartMap[refName].setOption(option, true)
    },
    buildEmptyOption (text) {
      return {
        animation: false,
        title: {
          text: text || '暂无数据',
          left: 'center',
          top: 'middle',
          textStyle: { color: '#8e8e93', fontSize: 13, fontWeight: 400 }
        }
      }
    },
    buildPieOption (data) {
      if (!data || !data.length) return this.buildEmptyOption('暂无数据')
      return {
        color: this.pieColors,
        tooltip: { trigger: 'item', formatter: '{b}: {c} ({d}%)' },
        legend: {
          type: 'scroll',
          bottom: 0,
          textStyle: { color: '#6d6d72', fontSize: 10 }
        },
        series: [{
          type: 'pie',
          radius: ['35%', '68%'],
          center: ['50%', '42%'],
          label: {
            color: '#3a3a3c',
            fontSize: 10,
            formatter: (params) => {
              const name = String(params.name || '')
              const shortName = name.length > 8 ? (name.slice(0, 8) + '…') : name
              return shortName + '\n' + params.percent + '%'
            }
          },
          itemStyle: { borderWidth: 1, borderColor: '#ffffff' },
          data: data
        }]
      }
    },
    buildBarOption (data) {
      if (!data || !data.length) return this.buildEmptyOption('暂无数据')
      const needZoom = data.length > 8
      return {
        color: ['#007aff'],
        tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } },
        dataZoom: needZoom ? [
          { type: 'inside', xAxisIndex: 0, start: 0, end: 70 },
          { type: 'slider', xAxisIndex: 0, height: 14, bottom: 30, start: 0, end: 70 }
        ] : [],
        grid: { left: 20, right: 16, top: 10, bottom: 62, containLabel: true },
        xAxis: {
          type: 'category',
          data: data.map(item => item.name),
          axisLabel: {
            interval: 0,
            rotate: 25,
            color: '#6d6d72',
            fontSize: 10,
            formatter: (value) => {
              const text = String(value || '')
              return text.length > 10 ? (text.slice(0, 10) + '…') : text
            }
          },
          axisTick: { show: false },
          axisLine: { lineStyle: { color: '#d2d2d7' } }
        },
        yAxis: {
          type: 'value',
          minInterval: 1,
          axisLabel: { color: '#6d6d72', fontSize: 10 },
          splitLine: { lineStyle: { color: '#f0f0f2' } }
        },
        series: [{
          type: 'bar',
          data: data.map(item => item.value),
          barWidth: '48%',
          itemStyle: { borderRadius: [6, 6, 0, 0] }
        }]
      }
    },
    renderCharts () {
      this.$nextTick(() => {
        this.initCharts()
        this.setChartOption('chartTaskTag', this.buildBarOption(this.taskTagChart))
        this.setChartOption('chartTaskLogDuration', this.buildPieOption(this.taskLogDurationChart))
        this.setChartOption('chartTaskLogStatus', this.buildPieOption(this.taskLogStatusChart))
        this.setChartOption('chartLoginUser', this.buildPieOption(this.loginUserChart))
        this.setChartOption('chartLoginIp', this.buildPieOption(this.loginIpChart))
        this.resizeCharts()
      })
    },
    loadStats () {
      if (!this.dateRange || this.dateRange.length !== 2) {
        this.dateRange = getDefaultRange()
      }
      this.loading = true
      systemService.stats({
        start_time: this.dateRange[0],
        end_time: this.dateRange[1]
      }, (data) => {
        this.overview = data && data.overview_all ? data.overview_all : this.overview
        this.taskTagChart = topNFromList((data && data.task_tag) || [], 9, '其他标签')
        this.taskLogDurationChart = (data && data.task_log_duration) || []
        this.taskLogStatusChart = (data && data.task_log_status) || []
        this.loginUserChart = topNFromList((data && data.login_user) || [], 9, '其他用户')
        this.loginIpChart = topNFromList((data && data.login_ip) || [], 9, '其他IP')
        this.loading = false
        this.renderCharts()
      }, () => {
        this.loading = false
      })
    }
  }
}
</script>

<style scoped>
.stats-page {
  background: #f5f5f7;
}

.stats-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 14px;
}

.stats-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.stats-title {
  margin: 0;
  font-size: 22px;
  font-weight: 700;
  color: #1d1d1f;
}

.stats-subtitle {
  margin: 4px 0 0;
  font-size: 13px;
  color: #86868b;
}

.stats-overview {
  display: grid;
  grid-template-columns: repeat(4, minmax(120px, 1fr));
  gap: 10px;
  margin-bottom: 14px;
}

.stats-filter-bar {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 10px;
}

.filter-label {
  font-size: 13px;
  color: #6d6d72;
}

.overview-item {
  background: #ffffff;
  border: 1px solid #e5e5ea;
  border-radius: 12px;
  padding: 12px 14px;
}

.overview-label {
  font-size: 12px;
  color: #86868b;
  margin-bottom: 6px;
}

.overview-value {
  font-size: 24px;
  font-weight: 700;
  color: #1d1d1f;
  line-height: 1;
}

.stats-section {
  margin-bottom: 14px;
}

.section-title {
  margin: 0 0 10px;
  font-size: 14px;
  font-weight: 600;
  color: #1d1d1f;
}

.chart-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(320px, 1fr));
  gap: 12px;
}

.chart-grid-single {
  grid-template-columns: 1fr;
}

.chart-card {
  background: #ffffff;
  border: 1px solid #e5e5ea;
  border-radius: 12px;
  padding: 12px;
}

.card-title {
  margin-bottom: 8px;
  font-size: 12px;
  font-weight: 600;
  color: #3a3a3c;
}

.chart-box {
  width: 100%;
  height: 300px;
}

.is-loading {
  animation: rotating 2s linear infinite;
}

@keyframes rotating {
  from {
    transform: rotate(0deg);
  }

  to {
    transform: rotate(360deg);
  }
}

@media (max-width: 1080px) {
  .stats-header {
    flex-direction: column;
    align-items: flex-start;
  }

  .stats-actions {
    width: 100%;
    display: grid;
    grid-template-columns: 1fr auto;
  }

  .stats-overview {
    grid-template-columns: repeat(2, minmax(120px, 1fr));
  }

  .chart-grid {
    grid-template-columns: 1fr;
  }

  .chart-box {
    height: 280px;
  }
}
</style>
