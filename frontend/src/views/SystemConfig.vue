<template>
  <div class="app-container">
    <el-tabs type="border-card">
      <el-tab-pane label="运营规则配置">
        <el-form label-width="160px" style="max-width: 600px; margin-top: 20px">
          <el-form-item label="配置保存">
            <el-button type="primary" @click="saveBasicConfig">保存配置</el-button>
          </el-form-item>
        </el-form>
      </el-tab-pane>

      <el-tab-pane label="内容安全配置">
        <el-form label-position="top">
          <el-form-item>
            <el-button type="danger" @click="saveSecurityConfig">更新安全配置</el-button>
          </el-form-item>
        </el-form>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup>
import { ElMessage } from 'element-plus'
import request from '@/utils/request'

const saveBasicConfig = async () => {
  try {
    await request.put('/api/v1/admin/config/basic', {})
    ElMessage.success('运营规则已更新')
  } catch (e) {}
}

const saveSecurityConfig = async () => {
  try {
    await request.put('/api/v1/admin/config/security', {})
    ElMessage.success('安全配置已更新')
  } catch (e) {}
}
</script>

<style scoped>
.app-container { padding: 20px; }
</style><template>
  <div class="app-container">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span class="title">帖子分类管理</span>
          <el-button type="primary" @click="openDialog()">
            <el-icon style="margin-right: 5px"><Plus /></el-icon> 新增分类
          </el-button>
        </div>
      </template>

      <el-table :data="typeList" v-loading="loading" stripe border style="width: 100%; margin-top: 10px;">
        <el-table-column prop="id" label="ID" width="80" align="center" />
        
        <el-table-column prop="name" label="分类名称" min-width="150" align="center">
          <template #default="scope">
            <el-tag effect="plain">{{ scope.row.name }}</el-tag>
          </template>
        </el-table-column>
        
        <el-table-column label="操作" width="200" align="center">
          <template #default="scope">
            <el-button size="small" type="primary" link @click="openDialog(scope.row)">
              <el-icon><Edit /></el-icon> 编辑
            </el-button>
            <el-button size="small" type="danger" link @click="handleDelete(scope.row)">
              <el-icon><Delete /></el-icon> 删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="400px"
      @close="resetForm"
    >
      <el-form :model="form" ref="formRef" :rules="rules" label-width="80px">
        <el-form-item label="分类名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入分类名称（如：失物招领）" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitForm">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import request from '@/utils/request'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Edit, Delete } from '@element-plus/icons-vue'

const loading = ref(false)
const typeList = ref([])
const dialogVisible = ref(false)
const formRef = ref(null)
const isEdit = ref(false)

const form = ref({
  id: null,
  name: ''
})

const rules = {
  name: [{ required: true, message: '请输入分类名称', trigger: 'blur' }]
}

const dialogTitle = computed(() => isEdit.value ? '编辑分类' : '新增分类')

onMounted(() => {
  fetchTypes()
})

// 获取分类列表
const fetchTypes = async () => {
  loading.value = true
  try {
    const res = await request.get('/api/v1/admin/config/types')
    typeList.value = res.list || []
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

// 打开弹窗
const openDialog = (row = null) => {
  if (row) {
    isEdit.value = true
    form.value = { id: row.id, name: row.name }
  } else {
    isEdit.value = false
    form.value = { id: null, name: '' }
  }
  dialogVisible.value = true
}

// 提交表单
const submitForm = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      try {
        if (isEdit.value) {
          // 编辑
          await request.put(`/api/v1/admin/config/types/${form.value.id}`, {
            name: form.value.name
          })
          ElMessage.success('更新成功')
        } else {
          // 新增
          await request.post('/api/v1/admin/config/types', {
            name: form.value.name
          })
          ElMessage.success('创建成功')
        }
        dialogVisible.value = false
        fetchTypes()
      } catch (e) {
        // 错误已由拦截器处理
      }
    }
  })
}

// 删除分类
const handleDelete = (row) => {
  ElMessageBox.confirm(
    `确定要删除分类 "${row.name}" 吗？删除后该分类下的帖子可能会显示异常。`,
    '删除确认',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      await request.delete(`/api/v1/admin/config/types/${row.id}`)
      ElMessage.success('删除成功')
      fetchTypes()
    } catch (e) {}
  })
}

const resetForm = () => {
  if (formRef.value) formRef.value.resetFields()
}
</script>

<style scoped>
.app-container {
  padding: 20px;
}
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.title {
  font-size: 16px;
  font-weight: bold;
}
</style>