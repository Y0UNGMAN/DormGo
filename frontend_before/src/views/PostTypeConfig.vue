<template>
  <div class="app-container">
    <el-card shadow="never">
      <div class="header-action">
        <h3>互助类型配置</h3>
        <el-button type="primary" @click="handleAdd">新增类型</el-button>
      </div>
      
      <el-table :data="typeList" border style="margin-top: 20px" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column prop="name" label="类型名称" />
        <el-table-column prop="icon" label="图标" width="100" align="center">
          <template #default="{ row }">
            <i :class="row.icon"></i>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-switch 
              v-model="row.status" 
              active-value="active" 
              inactive-value="disabled"
              @change="handleStatusChange(row)"
            />
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180" align="center">
          <template #default="{ row }">
            <el-button size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog :title="form.id ? '编辑类型' : '新增类型'" v-model="visible" width="400px">
      <el-form :model="form" label-width="80px">
        <el-form-item label="名称" required>
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="图标Class">
          <el-input v-model="form.icon" placeholder="例如: el-icon-basketball" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="visible = false">取消</el-button>
        <el-button type="primary" @click="submit">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import request from '@/utils/request'
import { ElMessage, ElMessageBox } from 'element-plus'

const loading = ref(false)
const typeList = ref([])
const visible = ref(false)
const form = reactive({ id: null, name: '', icon: '', status: 'active' })

onMounted(() => fetchTypes())

const fetchTypes = async () => {
  loading.value = true
  try {
    const res = await request.get('/v1/config/types') // SRS FR-PO-202
    typeList.value = res.list || []
  } catch (e) {} finally { loading.value = false }
}

const handleAdd = () => {
  form.id = null
  form.name = ''
  form.icon = ''
  visible.value = true
}

const handleEdit = (row) => {
  Object.assign(form, row)
  visible.value = true
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('删除后该类型的历史帖子可能显示异常，确认删除？')
    await request.delete(`/v1/admin/config/types/${row.id}`)
    ElMessage.success('删除成功')
    fetchTypes()
  } catch (e) {}
}

const handleStatusChange = async (row) => {
  try {
    await request.put(`/v1/admin/config/types/${row.id}/status`, { status: row.status })
    ElMessage.success('状态已更新')
  } catch (e) {
    row.status = row.status === 'active' ? 'disabled' : 'active' // 回滚
  }
}

const submit = async () => {
  if(!form.name) return ElMessage.warning('名称必填')
  try {
    if(form.id) {
      await request.put(`/v1/admin/config/types/${form.id}`, form)
    } else {
      await request.post('/v1/admin/config/types', form)
    }
    ElMessage.success('保存成功')
    visible.value = false
    fetchTypes()
  } catch (e) {}
}
</script>

<style scoped>
.app-container { padding: 20px; }
.header-action { display: flex; justify-content: space-between; align-items: center; margin-bottom: 10px; }
</style>