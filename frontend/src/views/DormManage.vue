<template>
  <div class="app-container">
    <el-card shadow="never">
      <div class="header-actions">
        <div class="search-box">
          <el-input 
            v-model="searchKeyword" 
            placeholder="楼栋名称" 
            clearable 
            @keyup.enter="handleSearch"
            style="width: 250px"
          >
            <template #append><el-button @click="handleSearch"><el-icon><Search /></el-icon></el-button></template>
          </el-input>
        </div>
        <el-button type="primary" @click="handleAdd">
          <el-icon><Plus /></el-icon> 新增楼栋
        </el-button>
      </div>

      <el-table :data="dormList" border stripe style="margin-top: 20px">
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column prop="name" label="楼栋名称" width="180" />
        <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
        <el-table-column prop="user_count" label="住户数" width="120" sortable align="center">
          <template #default="{ row }">
            <el-tag type="info" effect="plain">{{ row.user_count }} 人</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="180" align="center">
          <template #default="{ row }">
            <el-button link type="primary" @click="handleEdit(row)">编辑</el-button>
            <el-button link type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-container">
        <el-pagination
          v-model:current-page="pageInfo.page"
          v-model:page-size="pageInfo.pageSize"
          :total="total"
          layout="total, prev, pager, next"
          @current-change="fetchDormList"
        />
      </div>
    </el-card>

    <el-dialog :title="dialogType === 'add' ? '新增楼栋' : '编辑楼栋'" v-model="showDialog" width="500px">
      <el-form :model="dormForm" label-width="80px">
        <el-form-item label="名称" required>
          <el-input v-model="dormForm.name" placeholder="例如：南区1号楼" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input type="textarea" v-model="dormForm.description" rows="3" placeholder="楼栋简介..." />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showDialog = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, reactive } from 'vue'
import axios from 'axios'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Plus } from '@element-plus/icons-vue'

const dormList = ref([])
const total = ref(0)
const searchKeyword = ref('')
const pageInfo = reactive({ page: 1, pageSize: 10 })

const showDialog = ref(false)
const dialogType = ref('add')
const currentDormId = ref(null)
const dormForm = reactive({ name: '', description: '' })

onMounted(() => fetchDormList())

const fetchDormList = async () => {
  try {
    const token = localStorage.getItem('adminToken')
    const res = await axios.get('/api/v1/admin/dorms', {
      params: { ...pageInfo, keyword: searchKeyword.value },
      headers: { Authorization: `Bearer ${token}` }
    })
    dormList.value = res.data.list || []
    total.value = res.data.total || 0
  } catch (e) { console.error(e) }
}

const handleSearch = () => { pageInfo.page = 1; fetchDormList() }

const handleAdd = () => {
  dialogType.value = 'add'
  dormForm.name = ''
  dormForm.description = ''
  showDialog.value = true
}

const handleEdit = (row) => {
  dialogType.value = 'edit'
  currentDormId.value = row.id
  dormForm.name = row.name
  dormForm.description = row.description
  showDialog.value = true
}

const handleSubmit = async () => {
  if (!dormForm.name.trim()) return ElMessage.warning('请输入名称')
  try {
    const token = localStorage.getItem('adminToken')
    const headers = { Authorization: `Bearer ${token}` }
    if (dialogType.value === 'add') {
      await axios.post('/api/v1/admin/dorms', dormForm, { headers })
      ElMessage.success('新增成功')
    } else {
      await axios.put(`/api/v1/admin/dorms/${currentDormId.value}`, dormForm, { headers })
      ElMessage.success('更新成功')
    }
    showDialog.value = false
    fetchDormList()
  } catch (e) { ElMessage.error('操作失败') }
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(`确认删除 ${row.name} 吗？`, '提示', { type: 'warning' })
    const token = localStorage.getItem('adminToken')
    await axios.delete(`/api/v1/admin/dorms/${row.id}`, { headers: { Authorization: `Bearer ${token}` } })
    ElMessage.success('删除成功')
    fetchDormList()
  } catch (e) {}
}
</script>

<style scoped>
.app-container { padding: 20px; }
.header-actions { display: flex; justify-content: space-between; margin-bottom: 15px; }
.pagination-container { margin-top: 20px; text-align: right; }
</style>