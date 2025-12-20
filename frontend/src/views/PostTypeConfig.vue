<template>
  <div class="app-container">
    <el-card shadow="never" class="config-card">
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <span class="title">互助类型配置</span>
            <span class="subtitle">管理发帖时可选的互助板块分类</span>
          </div>
          <el-button type="primary" @click="handleAdd">
            <el-icon style="margin-right: 5px"><Plus /></el-icon> 新增分类
          </el-button>
        </div>
      </template>
      
      <el-table 
        :data="typeList" 
        stripe 
        border 
        style="width: 100%" 
        v-loading="loading"
        :header-cell-style="{ background: '#f5f7fa', color: '#606266' }"
      >
        <el-table-column prop="name" label="分类名称" align="center">
           <template #default="{ row }">
             <el-tag effect="light" size="large">{{ row.name }}</el-tag>
           </template>
        </el-table-column>

        <el-table-column label="操作" width="200" align="center">
          <template #default="{ row }">
            <el-button link type="primary" @click="handleEdit(row)">
               <el-icon style="margin-right: 4px"><Edit /></el-icon> 编辑
            </el-button>
            <el-button link type="danger" @click="handleDelete(row)">
               <el-icon style="margin-right: 4px"><Delete /></el-icon> 删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog 
      :title="dialogTitle" 
      v-model="visible" 
      width="400px" 
      @close="resetForm"
      destroy-on-close
    >
      <el-form :model="form" ref="formRef" :rules="rules" label-width="80px" style="margin-top: 10px;">
        <el-form-item label="分类名称" prop="name">
          <el-input 
            v-model="form.name" 
            placeholder="请输入互助类型名称（如：失物招领）" 
            @keyup.enter="submit"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="visible = false">取消</el-button>
          <el-button type="primary" @click="submit" :loading="submitting">保存</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import request from '@/utils/request'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Edit, Delete } from '@element-plus/icons-vue'

const loading = ref(false)
const submitting = ref(false)
const typeList = ref([])
const visible = ref(false)
const formRef = ref(null)
const form = reactive({ id: null, name: '' })

const rules = {
  name: [{ required: true, message: '请输入类型名称', trigger: 'blur' }]
}

const dialogTitle = computed(() => form.id ? '编辑分类' : '新增分类')

onMounted(() => {
  fetchTypes()
})

// 获取分类列表
const fetchTypes = async () => {
  loading.value = true
  try {
    // 【关键修复】修正接口路径，添加 /admin 前缀
    const res = await request.get('/api/v1/admin/config/types') 
    typeList.value = res.list || []
  } catch (e) {
    console.error('获取分类失败', e)
  } finally { 
    loading.value = false 
  }
}

// 打开新增弹窗
const handleAdd = () => {
  form.id = null
  form.name = ''
  visible.value = true
}

// 打开编辑弹窗
const handleEdit = (row) => {
  form.id = row.id
  form.name = row.name
  visible.value = true
}

// 删除分类
const handleDelete = (row) => {
  ElMessageBox.confirm(
    `确定要删除分类 "${row.name}" 吗？\n删除后属于该分类的帖子可能无法正常筛选。`,
    '删除确认',
    {
      confirmButtonText: '确定删除',
      cancelButtonText: '取消',
      type: 'warning',
      confirmButtonClass: 'el-button--danger'
    }
  ).then(async () => {
    try {
      await request.delete(`/api/v1/admin/config/types/${row.id}`)
      ElMessage.success('删除成功')
      fetchTypes()
    } catch (e) {
      // 错误已由拦截器处理
    }
  }).catch(() => {})
}

// 提交表单
const submit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      submitting.value = true
      try {
        if (form.id) {
          // 编辑
          await request.put(`/api/v1/admin/config/types/${form.id}`, { name: form.name })
          ElMessage.success('更新成功')
        } else {
          // 新增
          await request.post('/api/v1/admin/config/types', { name: form.name })
          ElMessage.success('创建成功')
        }
        visible.value = false
        fetchTypes()
      } catch (e) {
        console.error(e)
      } finally {
        submitting.value = false
      }
    }
  })
}

// 重置表单
const resetForm = () => {
  if (formRef.value) formRef.value.resetFields()
  form.id = null
  form.name = ''
}
</script>

<style scoped>
.app-container {
  padding: 20px;
  background-color: #f0f2f5;
  min-height: calc(100vh - 84px);
}

.config-card {
  border-radius: 8px;
  border: none;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-left {
  display: flex;
  flex-direction: column;
}

.title {
  font-size: 18px;
  font-weight: 600;
  color: #303133;
}

.subtitle {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
}
</style>