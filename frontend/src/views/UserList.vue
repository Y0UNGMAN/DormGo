<template>
  <div class="app-container">
    <el-card shadow="never" class="filter-card">
      <el-form :inline="true" :model="searchForm">
        <el-form-item label="关键词">
          <el-input 
            v-model="searchForm.keyword" 
            placeholder="学号 / 昵称 / 手机号" 
            clearable 
            style="width: 220px;" 
            @keyup.enter="handleSearch"
          >
            <template #prefix><el-icon><Search /></el-icon></template>
          </el-input>
        </el-form-item>
        <el-form-item label="信用分">
          <el-select v-model="searchForm.creditStatus" placeholder="全部" clearable style="width: 120px;">
            <el-option label="高 (>4.5)" value="high" />
            <el-option label="中 (3.0-4.5)" value="mid" />
            <el-option label="低 (<3.0)" value="low" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">查询</el-button>
          <el-button @click="resetSearch">重置</el-button>
          <el-button 
            type="danger" 
            @click="handleBatchBan" 
            :disabled="!selectedIds.length"
          >
            批量封禁
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card shadow="never" class="table-card">
      <el-table 
        :data="userList" 
        stripe 
        v-loading="loading" 
        style="width: 100%"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column prop="student_id" label="学号" width="120" sortable fixed />
        
        <el-table-column label="基本信息" min-width="180">
          <template #default="{ row }">
            <div class="user-cell">
              <el-avatar :size="36" :src="row.avatar || defaultAvatar" />
              <div class="user-info">
                <div class="nickname">{{ row.nickname }}</div>
                <div class="phone">{{ maskPhone(row.phone_number) }}</div>
              </div>
            </div>
          </template>
        </el-table-column>

        <el-table-column prop="dorm_name" label="所属楼栋" width="150" />
        
        <el-table-column prop="credit_score" label="信用分" width="140" sortable align="center">
          <template #default="{ row }">
             <div class="credit-cell">
               <span :class="getCreditClass(row.credit_score)" class="credit-num">{{ row.credit_score }}</span>
               <el-button link type="primary" size="small" icon="Edit" @click="openCreditAdjust(row)" />
             </div>
          </template>
        </el-table-column>

        <el-table-column label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-switch
              v-model="row.status"
              active-value="normal"
              inactive-value="disabled"
              inline-prompt
              active-text="正常"
              inactive-text="封禁"
              :before-change="() => handleStatusChange(row)"
            />
          </template>
        </el-table-column>

        <el-table-column label="操作" width="150" fixed="right" align="center">
          <template #default="{ row }">
            <el-button link type="primary" @click="handleView(row)">详情</el-button>
            <el-button link type="primary" @click="handleTags(row)">标签</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-container">
        <el-pagination
          v-model:current-page="pageInfo.page"
          v-model:page-size="pageInfo.pageSize"
          :total="total"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="fetchUserList"
          @current-change="fetchUserList"
        />
      </div>
    </el-card>

    <UserCreditAdjust v-model="creditDialogVisible" :user-info="currentUser" @success="fetchUserList" />
  </div>
</template>

<script setup>
import { ref, onMounted, reactive } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Edit } from '@element-plus/icons-vue'
import request from '@/utils/request'
import UserCreditAdjust from '@/components/UserCreditAdjust.vue'

const userList = ref([])
const total = ref(0)
const loading = ref(false)
const selectedIds = ref([])
const defaultAvatar = 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'

const searchForm = reactive({ keyword: '', creditStatus: '' })
const pageInfo = reactive({ page: 1, pageSize: 10 })

// 信用分调整
const creditDialogVisible = ref(false)
const currentUser = ref({})

onMounted(() => fetchUserList())

const fetchUserList = async () => {
  loading.value = true
  try {
    const res = await request.get('/api/v1/admin/users', { params: { ...pageInfo, ...searchForm } })
    userList.value = res.list || []
    total.value = res.total || 0
  } catch (e) { console.error(e) } 
  finally { loading.value = false }
}

// 隐私脱敏
const maskPhone = (phone) => {
  if (!phone) return '未绑定'
  return phone.replace(/(\d{3})\d{4}(\d{4})/, '$1****$2')
}

const getCreditClass = (score) => {
  if (score >= 4.5) return 'text-success'
  if (score >= 3.0) return 'text-warning'
  return 'text-danger'
}

const handleSearch = () => { pageInfo.page = 1; fetchUserList() }
const resetSearch = () => { searchForm.keyword = ''; searchForm.creditStatus = ''; handleSearch() }

const handleSelectionChange = (val) => {
  selectedIds.value = val.map(item => item.id)
}

// 【修改点2】实现批量封禁逻辑，并添加确认弹窗
const handleBatchBan = async () => {
  if (!selectedIds.value.length) return

  try {
    await ElMessageBox.confirm(
      `确定要批量封禁选中的 ${selectedIds.value.length} 位用户吗？`,
      '警告',
      { 
        type: 'warning',
        confirmButtonText: '确定', // 明确设置为“确定”
        cancelButtonText: '取消'
      }
    )

    await request.put('/api/v1/admin/users/batch/status', { ids: selectedIds.value, status: 'disabled' })
    ElMessage.success('操作成功')
    fetchUserList()
  } catch (e) {
    if (e !== 'cancel') console.error(e)
  }
}

// 【修改点3】修改单体封禁/解封的确认按钮文案
const handleStatusChange = (user) => {
  const isBanning = user.status === 'normal'
  return ElMessageBox.confirm(
    `确定要${isBanning ? '封禁' : '解封'}该用户吗？`,
    '提示',
    { 
      type: 'warning',
      confirmButtonText: '确定', // 明确设置为“确定”
      cancelButtonText: '取消'
    }
  ).then(async () => {
    try {
      const newStatus = isBanning ? 'disabled' : 'normal'
      await request.put(`/api/v1/admin/users/${user.id}/status`, { status: newStatus })
      ElMessage.success('操作成功')
      return true
    } catch (e) { return false }
  }).catch(() => false)
}

const openCreditAdjust = (user) => {
  currentUser.value = user
  creditDialogVisible.value = true
}

const handleView = (user) => {
  ElMessage.info('查看用户详情功能开发中')
}

const handleTags = async (user) => {
  try {
    // 模拟标签数据
    ElMessageBox.alert('暂无标签数据', '用户兴趣标签', { confirmButtonText: '确定' })
  } catch (e) {}
}
</script>

<style scoped>
.app-container { padding: 20px; }
.filter-card { margin-bottom: 20px; }
.user-cell { display: flex; align-items: center; gap: 12px; }
.user-info .phone { font-size: 12px; color: #909399; }
.credit-cell { display: flex; align-items: center; justify-content: center; gap: 5px; }
.text-success { color: #67c23a; }
.text-warning { color: #e6a23c; }
.text-danger { color: #f56c6c; }
.pagination-container { margin-top: 20px; text-align: right; }
</style>