<template>
  <el-container>
    <system-sidebar></system-sidebar>
    <el-main class="ios-log-page">
      <div class="search-bar">
        <div class="page-title">登录日志</div>
        <div class="filter-row">
          <div class="search-filters">
            <el-select
              v-model.trim="searchParams.username"
              class="search-select"
              filterable
              clearable
              remote
              reserve-keyword
              placeholder="用户名"
              :remote-method="searchUsernames"
              :loading="usernameLoading"
              @change="search()">
              <el-option v-for="item in usernameOptions" :key="item" :label="item" :value="item"></el-option>
            </el-select>
            <el-input
              v-model.trim="searchParams.ip"
              placeholder="登录IP"
              class="search-input"
              clearable
              @input="debouncedSearch"
              @keyup.native.enter="search()">
            </el-input>
          </div>
          <div class="search-actions">
            <el-button class="ios-action-btn" type="primary" plain @click="refresh">刷新</el-button>
          </div>
        </div>
      </div>

      <div class="ios-table-card">
        <el-table :data="logs" ref="table" class="ios-data-table" style="width: 100%">
          <el-table-column prop="id" label="ID">
          </el-table-column>
          <el-table-column prop="username" label="用户名">
          </el-table-column>
          <el-table-column prop="ip" label="登录IP">
          </el-table-column>
          <el-table-column label="登录时间">
            <template slot-scope="scope">
              {{ scope.row.created | formatTime }}
            </template>
          </el-table-column>
        </el-table>
      </div>

      <div class="pagination-wrap">
        <el-pagination class="ios-pagination" background layout="prev, pager, next, sizes, total" :total="logTotal" :current-page="searchParams.page" :page-size="searchParams.page_size" @size-change="changePageSize" @current-change="changePage" @prev-click="changePage" @next-click="changePage">
        </el-pagination>
      </div>
    </el-main>
  </el-container>
</template>

<script>
import systemSidebar from './sidebar'
import systemService from '../../api/system'
import '../../styles/ios-table.css'
export default {
  name: 'login-log',
  data () {
    return {
      logs: [],
      logTotal: 0,
      searchParams: {
        page_size: 20,
        page: 1,
        username: '',
        ip: ''
      },
      usernameOptions: [],
      usernameLoading: false,
      searchTimer: null
    }
  },
  created () {
    this.searchUsernames('')
    this.search()
  },
  beforeDestroy () {
    if (this.searchTimer) clearTimeout(this.searchTimer)
  },
  components: { systemSidebar },
  methods: {
    debouncedSearch () {
      if (this.searchTimer) clearTimeout(this.searchTimer)
      this.searchTimer = setTimeout(() => {
        this.searchParams.page = 1
        this.search()
      }, 400)
    },
    searchUsernames (keyword) {
      this.usernameLoading = true
      systemService.loginLogUsers({ keyword }, (data) => {
        this.usernameOptions = Array.isArray(data) ? data : []
        this.usernameLoading = false
      }, () => {
        this.usernameLoading = false
      })
    },
    changePage (page) {
      this.searchParams.page = page
      this.search()
    },
    changePageSize (pageSize) {
      this.searchParams.page_size = pageSize
      this.search()
    },
    search (callback = null) {
      systemService.loginLogList(this.searchParams, (data) => {
        this.logs = data.data
        this.logTotal = data.total
        if (typeof callback === 'function') {
          callback()
        }
      })
    },
    refresh () {
      this.search(() => {
        this.$message.success('刷新成功')
      })
    }
  }
}
</script>

<style scoped>
.ios-log-page {
  background: #f5f5f7;
}

.search-bar {
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin-bottom: 16px;
  align-items: stretch;
}

.filter-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  flex-wrap: wrap;
}

.search-filters {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.search-select {
  width: 180px;
}

.search-input {
  width: 180px;
}

.page-title {
  font-size: 20px;
  font-weight: 600;
  color: #1d1d1f;
  margin: 0;
}

.search-actions {
  display: flex;
  gap: 8px;
  flex-shrink: 0;
}

.pagination-wrap {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}
</style>
