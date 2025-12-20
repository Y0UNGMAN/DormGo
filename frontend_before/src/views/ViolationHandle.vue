<template>
  <div class="app-container">
    <el-card shadow="never">
      <div class="filter-bar">
        <el-input v-model="searchForm.userKeyword" placeholder="用户关键词" style="width: 200px" clearable />
        <el-select v-model="searchForm.status" placeholder="状态" clearable style="width: 120px">
          <el-option label="待处理" value="pending" />
          <el-option label="已处理" value="processed" />
        </el-select>
        <el-select v-model="searchForm.type" placeholder="违规类型" clearable style="width: 120px">
          <el-option label="内容" value="content" />
          <el-option label="行为" value="behavior" />
          <el-option label="其他" value="other" />
        </el-select>
        <el-button type="primary" @click="handleSearch">搜索</el-button>
        <el-button @click="resetSearch">重置</el-button>
      </div>

      <el-table :data="violationList" border style="margin-top: 20px">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="违规用户" width="180">
          <template #default="{ row }">
            {{ row.user.nickname }} <span style="color: #999">({{ row.user.student_id }})</span>
          </template>
        </el-table-column>
        <el-table-column prop="type" label="类型" width="100">
          <template #default="{ row }">
            <el-tag>{{ formatType(row) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="reason" label="违规原因" min-width="200" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'processed' ? 'success' : 'warning'">
              {{ row.status === 'processed' ? '已处理' : '待处理' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="时间" width="180" />
        <el-table-column label="操作" width="180">
          <template #default="{ row }">
            <el-button size="small" @click="handleView(row)">详情</el-button>
            <el-button v-if="row.status === 'pending'" size="small" type="primary" @click="handleProcess(row)">处理</el-button>
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
import axios from 'axios'
import { ElMessage, ElMessageBox } from 'element-plus'

const violationList = ref([])
const showProcessDialog = ref(false)
const currentViolation = ref(null)
const searchForm = reactive({ userKeyword: '', status: '', type: '' })
const processForm = reactive({ punishment: '', notes: '' })
const pageInfo = reactive({ page: 1, pageSize: 10 })

onMounted(() => fetchList())

const fetchList = async () => {
  try {
    const token = localStorage.getItem('adminToken')
    const res = await axios.get('/api/v1/admin/violations', {
      params: { ...pageInfo, ...searchForm },
      headers: { Authorization: `Bearer ${token}` }
    })
    violationList.value = res.data.list
  } catch (e) {}
}

const formatType = (row) => ({ content: '内容', behavior: '行为', other: '其他' }[row.type] || row.type)
const handleSearch = () => { pageInfo.page = 1; fetchList() }
const resetSearch = () => { Object.assign(searchForm, { userKeyword: '', status: '', type: '' }); handleSearch() }

const handleView = (row) => {
  const content = `
    <p>原因：${row.reason}</p>
    <p>时间：${row.violation_time}</p>
    ${row.status === 'processed' ? `<hr><p>处罚：${row.punishment}</p><p>备注：${row.notes}</p>` : ''}
  `
  ElMessageBox.alert(content, '详情', { dangerouslyUseHTMLString: true })
}

const handleProcess = (row) => {
  currentViolation.value = row
  processForm.punishment = ''
  processForm.notes = ''
  showProcessDialog.value = true
}

const submitProcess = async () => {
  try {
    const token = localStorage.getItem('adminToken')
    await axios.put(`/api/v1/admin/violations/${currentViolation.value.id}/process`, processForm, {
      headers: { Authorization: `Bearer ${token}` }
    })
    ElMessage.success('已处理')
    showProcessDialog.value = false
    fetchList()
  } catch (e) { ElMessage.error('失败') }
}
</script>

<style scoped>
.app-container { padding: 20px; }
.filter-bar { display: flex; gap: 10px; flex-wrap: wrap; }
</style>