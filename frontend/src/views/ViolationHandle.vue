<template>
  <div class="app-container">
    <el-card shadow="never">
      <div class="filter-bar">
        <el-input v-model="searchForm.userKeyword" placeholder="学号/昵称" style="width: 180px" clearable />
        <el-select v-model="searchForm.type" placeholder="违规类型" clearable style="width: 140px">
           <el-option label="内容违规" value="content" />
           <el-option label="行为不当" value="behavior" />
           <el-option label="恶意举报" value="report" />
        </el-select>
        <el-input v-model="searchForm.reason" placeholder="违规原因关键词" style="width: 180px" clearable />
        <el-select v-model="searchForm.status" placeholder="处理状态" clearable style="width: 120px">
          <el-option label="待处理" value="pending" />
          <el-option label="已处理" value="processed" />
        </el-select>
        <el-button type="primary" @click="handleSearch">搜索</el-button>
        <el-button @click="resetSearch">重置</el-button>
      </div>

      <el-table :data="violationList" border style="margin-top: 20px" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" align="center" />
        
        <el-table-column label="违规用户" width="220">
          <template #default="{ row }">
            <div v-if="row.user">
              <span style="font-weight: bold">{{ row.user.username }}</span>
              <br/>
              <span style="color: #999; font-size: 12px">学号: {{ row.user.studentid }}</span>
            </div>
            <span v-else style="color: #ccc">未知用户</span>
          </template>
        </el-table-column>

        <el-table-column prop="type" label="类型" width="100" align="center">
          <template #default="{ row }">
            <el-tag>{{ formatType(row.type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="reason" label="违规原因" min-width="200" />
        <el-table-column prop="status" label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 'processed' ? 'success' : 'warning'">
              {{ row.status === 'processed' ? '已处理' : '待处理' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="时间" width="180">
           <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="150" align="center">
          <template #default="{ row }">
            <el-button v-if="row.status === 'pending'" size="small" type="primary" @click="handleProcess(row)">处理</el-button>
            <el-button v-else size="small" type="info" disabled>已处理</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog title="违规处理" v-model="showProcessDialog" width="500px">
      <el-form label-position="top">
        <el-form-item label="处罚措施">
          <el-input v-model="processForm.punishment" placeholder="例如：警告、禁言3天" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input type="textarea" v-model="processForm.notes" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showProcessDialog = false">取消</el-button>
        <el-button type="primary" @click="submitProcess">提交</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, reactive } from 'vue'
import request from '@/utils/request'
import { ElMessage } from 'element-plus'

const violationList = ref([])
const loading = ref(false)
const showProcessDialog = ref(false)
const currentViolation = ref(null)
const searchForm = reactive({ userKeyword: '', status: '', type: '', reason: '' })
const processForm = reactive({ punishment: '', notes: '' })

onMounted(() => fetchList())

const fetchList = async () => {
  loading.value = true
  try {
    const res = await request.get('/api/v1/admin/violations', {
      params: { ...searchForm }
    })
    violationList.value = res.list || []
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const handleSearch = () => fetchList()
const resetSearch = () => {
  Object.assign(searchForm, { userKeyword: '', status: '', type: '', reason: '' })
  fetchList()
}

const handleProcess = (row) => {
  currentViolation.value = row
  processForm.punishment = ''
  processForm.notes = ''
  showProcessDialog.value = true
}

const submitProcess = async () => {
  try {
    await request.put(`/api/v1/admin/violations/${currentViolation.value.id}/process`, processForm)
    ElMessage.success('已处理')
    showProcessDialog.value = false
    fetchList()
  } catch (e) { }
}

const formatType = (val) => {
  const map = { content: '内容违规', behavior: '行为不当', report: '恶意举报' }
  return map[val] || val
}

const formatTime = (time) => time ? new Date(time).toLocaleString() : ''
</script>

<style scoped>
.app-container { padding: 20px; }
.filter-bar { display: flex; gap: 10px; margin-bottom: 20px; flex-wrap: wrap; }
</style>
