<template>
  <div class="app-container">
    <el-card shadow="never">
      <el-button type="primary" @click="dialogVisible = true" style="margin-bottom: 20px">发布通知</el-button>
      <el-table :data="list" border>
        <el-table-column prop="title" label="标题" />
        <el-table-column prop="content" label="内容" show-overflow-tooltip />
        <el-table-column prop="created_at" label="发布时间" width="180">
           <template #default="{ row }">{{ new Date(row.created_at).toLocaleString() }}</template>
        </el-table-column>
        <el-table-column label="操作" width="100" align="center">
          <template #default="{ row }">
            <el-button link type="danger" @click="del(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" title="发布通知">
      <el-form :model="form">
        <el-form-item label="标题"><el-input v-model="form.title" /></el-form-item>
        <el-form-item label="内容"><el-input type="textarea" v-model="form.content" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button type="primary" @click="submit">发布</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, reactive } from 'vue'
import request from '@/utils/request'
import { ElMessage } from 'element-plus'

const list = ref([])
const dialogVisible = ref(false)
const form = reactive({ title: '', content: '' })

onMounted(() => fetchList())

const fetchList = async () => {
  const res = await request.get('/api/v1/admin/notices')
  list.value = res.list || []
}

const submit = async () => {
  await request.post('/api/v1/admin/notices', form)
  ElMessage.success('发布成功')
  dialogVisible.value = false
  fetchList()
}

const del = async (row) => {
  await request.delete(`/api/v1/admin/notices/${row.id}`)
  fetchList()
}
</script>

<style scoped>
.app-container { padding: 20px; }
</style>