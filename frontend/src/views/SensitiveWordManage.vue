<template>
  <div class="app-container">
    <el-card shadow="never">
      <div class="toolbar">
        <el-button type="primary" @click="handleAdd">新增敏感词</el-button>
      </div>

      <el-table :data="list" border style="margin-top: 20px" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="word" label="敏感词" />
        <el-table-column label="操作" width="150" align="center">
          <template #default="{ row }">
            <el-button link type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog title="新增敏感词" v-model="dialogVisible" width="400px">
      <el-form :model="form" label-width="80px">
        <el-form-item label="敏感词">
          <el-input v-model="form.word" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitForm">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/request'

const loading = ref(false)
const list = ref([])
const dialogVisible = ref(false)
const form = reactive({ word: '' })

onMounted(() => fetchList())

const fetchList = async () => {
  loading.value = true
  try {
    const res = await request.get('/api/v1/admin/config/sensitive-words')
    list.value = res.list || []
  } catch (e) {} finally { loading.value = false }
}

const handleAdd = () => {
  form.word = ''
  dialogVisible.value = true
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(`确定删除敏感词 "${row.word}" 吗？`)
    await request.delete(`/api/v1/admin/config/sensitive-words/${row.id}`)
    ElMessage.success('删除成功')
    fetchList()
  } catch (e) {}
}

const submitForm = async () => {
  if (!form.word) return ElMessage.warning('请输入内容')
  try {
    await request.post('/api/v1/admin/config/sensitive-words', form)
    ElMessage.success('保存成功')
    dialogVisible.value = false
    fetchList()
  } catch (e) {}
}
</script>

<style scoped>
.app-container { padding: 20px; }
.toolbar { margin-bottom: 20px; }
</style>