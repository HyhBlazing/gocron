<template>
  <el-container>
    <task-sidebar></task-sidebar>
    <el-main class="ios-host-page">
      <div class="page-header">
        <h2 class="page-title">节点管理</h2>
      </div>

      <div class="search-bar">
        <div class="search-filters">
          <el-input v-model.trim="searchParams.alias" placeholder="节点名称" class="search-input" clearable @input="debouncedSearch" @keyup.native.enter="search()"></el-input>
          <el-input v-model.trim="searchParams.name" placeholder="主机IP" class="search-input" clearable @input="debouncedSearch" @keyup.native.enter="search()"></el-input>
        </div>
        <div class="search-actions">
          <el-button class="ios-action-btn" type="primary" plain v-if="isAdmin" @click="openHostModal()">新增节点</el-button>
          <el-button class="ios-action-btn" type="primary" plain @click="refresh">刷新</el-button>
        </div>
      </div>

      <div class="ios-table-card">
        <el-table v-loading="loading" :data="hosts" class="ios-data-table" style="width: 100%">
          <el-table-column prop="id" label="节点ID" width="90">
          </el-table-column>
          <el-table-column prop="alias" label="节点名称">
          </el-table-column>
          <el-table-column prop="name" label="主机IP">
          </el-table-column>
          <el-table-column prop="port" label="端口" width="80">
          </el-table-column>
          <el-table-column label="查看定时" width="100">
            <template slot-scope="scope">
              <el-button type="text" class="link-btn" @click="toTasks(scope.row)">查看定时</el-button>
            </template>
          </el-table-column>
          <el-table-column prop="remark" label="备注" show-overflow-tooltip>
          </el-table-column>
          <el-table-column label="操作" width="260" v-if="isAdmin">
            <template slot-scope="scope">
              <div class="action-buttons">
                <el-button type="primary" size="mini" @click="openHostModal(scope.row)">编辑</el-button>
                <el-button type="info" size="mini" @click="ping(scope.row)">测试连接</el-button>
                <el-button type="danger" size="mini" @click="remove(scope.row)">删除</el-button>
              </div>
            </template>
          </el-table-column>
        </el-table>
      </div>
      <div class="pagination-bar">
        <span class="pagination-total">共 {{ hostTotal }} 条</span>
        <el-pagination class="ios-pagination" background layout="prev, pager, next, jumper, sizes" :total="hostTotal" :current-page="searchParams.page" :page-size="searchParams.page_size" :page-sizes="[10, 20, 50, 100]" @size-change="changePageSize" @current-change="changePage">
        </el-pagination>
      </div>

      <!-- 节点新增/编辑模态窗 -->
      <el-dialog :title="hostForm.id ? '编辑节点' : '新增节点'" :visible.sync="hostModalVisible" width="480px" :close-on-click-modal="false" :append-to-body="true" custom-class="pwd-dialog" @closed="resetHostForm">
        <el-form ref="hostForm" :model="hostForm" :rules="hostRules" label-width="90px">
          <el-form-item label="节点名称" prop="alias">
            <el-input v-model.trim="hostForm.alias" placeholder="请输入节点展示名称"></el-input>
          </el-form-item>
          <el-form-item label="主机IP" prop="name">
            <el-input v-model.trim="hostForm.name" placeholder="请输入 IP 或主机名"></el-input>
          </el-form-item>
          <el-form-item label="端口" prop="port">
            <el-input-number v-model="hostForm.port" :min="1" :max="65535" label="端口"></el-input-number>
          </el-form-item>
          <el-form-item label="备注" prop="remark">
            <el-input type="textarea" v-model="hostForm.remark" placeholder="节点备注信息"></el-input>
          </el-form-item>
        </el-form>
        <span slot="footer">
          <el-button @click="hostModalVisible = false">取消</el-button>
          <el-button type="primary" @click="submitHostForm" :loading="hostLoading">保存</el-button>
        </span>
      </el-dialog>
    </el-main>
  </el-container>
</template>

<script>
import hostService from '../../api/host'
import taskSidebar from '../task/sidebar'
import '../../styles/ios-table.css'

export default {
  name: 'host-list',
  components: { taskSidebar },
  data () {
    return {
      hosts: [],
      hostTotal: 0,
      loading: false,
      searchParams: {
        page_size: 20,
        page: 1,
        name: '',
        alias: ''
      },
      isAdmin: this.$store.getters.user.isAdmin,
      // 模态窗状态
      hostModalVisible: false,
      hostLoading: false,
      hostForm: {
        id: '',
        alias: '',
        name: '',
        port: 5921,
        remark: ''
      },
      hostRules: {
        alias: [{ required: true, message: '请输入节点名称', trigger: 'blur' }],
        name: [{ required: true, message: '请输入主机名/IP', trigger: 'blur' }],
        port: [{ required: true, message: '请输入端口', trigger: 'blur' }]
      }
    }
  },
  created () {
    this.search()
  },
  beforeDestroy () {
    if (this._searchTimer) clearTimeout(this._searchTimer)
  },
  methods: {
    debouncedSearch () {
      if (this._searchTimer) clearTimeout(this._searchTimer)
      this._searchTimer = setTimeout(() => {
        this.searchParams.page = 1
        this.search()
      }, 400)
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
      hostService.list(this.searchParams, (data) => {
        this.hosts = data.data
        this.hostTotal = data.total
        this.loading = false
        if (typeof callback === 'function') {
          callback()
        }
      }, () => {
        this.loading = false
      })
    },
    remove (item) {
      this.$appConfirm(() => {
        hostService.remove(item.id, () => this.refresh())
      })
    },
    ping (item) {
      hostService.ping(item.id, () => {
        this.$message.success('连接成功')
      })
    },
    toTasks (item) {
      this.$router.push({
        path: '/task',
        query: {
          host_id: item.id
        }
      })
    },
    refresh () {
      this.search(() => {
        this.$message.success('刷新成功')
      })
    },
    // 模态窗方法
    openHostModal (item = null) {
      if (item) {
        this.hostForm.id = item.id
        this.hostForm.alias = item.alias
        this.hostForm.name = item.name
        this.hostForm.port = item.port
        this.hostForm.remark = item.remark
      } else {
        this.hostForm.id = ''
      }
      this.hostModalVisible = true
    },
    resetHostForm () {
      this.hostForm = { id: '', alias: '', name: '', port: 5921, remark: '' }
      this.$refs.hostForm && this.$refs.hostForm.resetFields()
    },
    submitHostForm () {
      this.$refs.hostForm.validate((valid) => {
        if (!valid) return
        this.hostLoading = true
        hostService.update(this.hostForm, () => {
          this.hostLoading = false
          this.hostModalVisible = false
          this.$message.success('操作成功')
          this.refresh()
        }, () => {
          this.hostLoading = false
        })
      })
    }
  }
}
</script>

<style scoped>
.ios-host-page {
  background: #f5f5f7;
}

.search-bar {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 12px;
  margin-top: 12px;
  align-items: center;
  justify-content: space-between;
}

.search-filters {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  flex: 1;
  align-items: center;
}

.search-input {
  width: 15rem;
}

.action-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.action-row-left {
  display: flex;
  gap: 8px;
}

.action-row-right {
  display: flex;
  gap: 8px;
}

.pagination-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 12px;
  flex-wrap: wrap;
  gap: 8px;
}

.pagination-bar .pagination-total {
  color: #606266;
  font-size: 14px;
}

.action-buttons {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  white-space: nowrap;
}

.action-buttons .el-button,
.action-row .el-button {
  margin-left: 0;
}

.link-btn {
  color: #007aff;
  padding: 0;
}

.link-btn:hover {
  color: #0051d5;
}

.page-title {
  font-size: 20px;
  font-weight: 600;
  color: #1d1d1f;
  margin: 0;
}
</style>
