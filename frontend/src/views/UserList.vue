<template>
  <div class="app-container">
    <el-card shadow="never" class="table-card">
      <template #header>
        <div class="card-header">
          <span class="title">用户管理</span>
          <div class="header-actions">
            <el-input
              v-model="searchKeyword"
              placeholder="搜索学号/用户名"
              clearable
              @clear="fetchData"
              @keyup.enter="fetchData"
              style="width: 240px; margin-right: 12px;"
            >
              <template #prefix>
                <el-icon><Search /></el-icon>
              </template>
            </el-input>
            <el-button type="primary" @click="fetchData">
              <el-icon style="margin-right: 4px"><Refresh /></el-icon> 刷新
            </el-button>
          </div>
        </div>
      </template>

      <el-table 
        :data="userList" 
        style="width: 100%" 
        v-loading="loading"
        stripe
        border
      >
        <el-table-column prop="id" label="ID" width="80" align="center" />
        
        <el-table-column prop="student_id" label="学号" min-width="120" align="center" show-overflow-tooltip />
        
        <el-table-column prop="nickname" label="用户名" min-width="120" align="center" show-overflow-tooltip />
        
        <el-table-column label="状态" width="100" align="center">
          <template #default="scope">
            <el-tag :type="scope.row.status === 1 ? 'success' : 'danger'" effect="light">
              {{ scope.row.status === 1 ? '正常' : '已封禁' }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column prop="created_at" label="注册时间" width="180" align="center">
          <template #default="scope">
            {{ formatTime(scope.row.created_at) }}
          </template>
        </el-table-column>

        <el-table-column label="操作" width="120" align="center" fixed="right">
          <template #default="scope">
            <el-button 
              :type="scope.row.status === 1 ? 'danger' : 'success'" 
              size="small" 
              plain
              @click="handleStatusChange(scope.row)"
            >
              <el-icon style="margin-right: 2px">
                <component :is="scope.row.status === 1 ? 'Lock' : 'Unlock'" />
              </el-icon>
              {{ scope.row.status === 1 ? '封禁' : '解封' }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-container">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          :total="total"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import request from '@/utils/request'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Refresh, Lock, Unlock } from '@element-plus/icons-vue'

const loading = ref(false)
const userList = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const searchKeyword = ref('')

// 初始化加载
onMounted(() => {
  fetchData()
})

// 获取数据
const fetchData = async () => {
  loading.value = true
  try {
    const res = await request.get('/api/v1/admin/users', {
      params: {
        page: currentPage.value,
        size: pageSize.value,
        keyword: searchKeyword.value
      }
    })
    userList.value = res.list || []
    total.value = res.total || 0
  } catch (error) {
    console.error('获取用户列表失败', error)
  } finally {
    loading.value = false
  }
}

// 核心逻辑：处理封禁/解封状态切换
const handleStatusChange = (row) => {
  // 判断当前行为是"封禁"还是"解封"
  const isBanAction = row.status === 1 
  const actionText = isBanAction ? '封禁' : '解封'
  
  // 根据不同操作显示不同的提示文案
  const confirmText = isBanAction 
    ? `确定要封禁用户 "${row.nickname}" (学号: ${row.student_id}) 吗？此操作将限制该用户登录。`
    : `确定要解除用户 "${row.nickname}" (学号: ${row.student_id}) 的封禁吗？用户将恢复正常登录权限。`

  ElMessageBox.confirm(
    confirmText,
    `${actionText}确认`,
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: isBanAction ? 'warning' : 'success', // 封禁使用警告色，解封使用成功色
      draggable: true,
    }
  ).then(async () => {
    try {
      // 发送请求给后端，后端会自动根据当前状态进行反转 (Toggle Logic)
      await request.post('/api/v1/admin/users/ban', { user_id: row.id })
      
      ElMessage.success(`用户${actionText}成功`)
      
      // 必须刷新列表以获取最新的 status，从而更新按钮显示
      fetchData() 
    } catch (error) {
      // 错误通常由拦截器统一处理
    }
  }).catch(() => {
    // 用户取消操作，无需处理
  })
}

// 分页大小改变
const handleSizeChange = (val) => {
  pageSize.value = val
  fetchData()
}

// 页码改变
const handleCurrentChange = (val) => {
  currentPage.value = val
  fetchData()
}

// 时间格式化
const formatTime = (timeStr) => {
  if (!timeStr) return '-'
  const date = new Date(timeStr)
  return date.toLocaleString()
}
</script>

<style scoped>
.app-container {
  padding: 20px;
  background-color: #f0f2f5;
  min-height: 100vh;
}

.table-card {
  border-radius: 8px;
  border: none;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.title {
  font-size: 18px;
  font-weight: 600;
  color: #303133;
}

.header-actions {
  display: flex;
  align-items: center;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>