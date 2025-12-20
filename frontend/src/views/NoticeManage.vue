<template>
  <div class="app-container">
    <el-card shadow="never">
      <div class="header-actions">
        <h3>通知发布与历史</h3>
        <el-button type="primary" @click="openDialog">📢 发布新通知</el-button>
      </div>

      <el-table :data="list" border style="margin-top: 20px" v-loading="loading">
        <el-table-column prop="content" label="通知内容" show-overflow-tooltip />
        <el-table-column prop="target" label="对象" width="100">
           <template #default="{ row }">
             <el-tag>{{ row.target === 'all' ? '全员' : '个人' }}</el-tag>
           </template>
        </el-table-column>
        <el-table-column prop="created_at" label="发布时间" width="180">
           <template #default="{ row }">{{ new Date(row.created_at).toLocaleString() }}</template>
        </el-table-column>
        <el-table-column label="操作" width="100" align="center">
          <template #default="{ row }">
            <el-button link type="danger" @click="del(row)">删除记录</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" title="发送系统通知" width="500px">
      <el-form :model="form" label-width="80px">
        <el-form-item label="发送对象">
          <el-radio-group v-model="form.target_type">
            <el-radio label="all">全体用户</el-radio>
            <el-radio label="specific">指定用户</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="学号" v-if="form.target_type === 'specific'">
          <el-input v-model="form.target_value" placeholder="请输入目标学号" />
        </el-form-item>

        <el-form-item label="内容">
          <el-input type="textarea" v-model="form.content" rows="4" placeholder="请输入通知正文" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submit">发送并推送</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, reactive } from 'vue'
import request from '@/utils/request'
import { ElMessage } from 'element-plus'

const list = ref([])
const loading = ref(false)
const dialogVisible = ref(false)
// form 中去掉了 title
const form = reactive({ target_type: 'all', target_value: '', content: '' })

onMounted(() => fetchList())

const fetchList = async () => {
  loading.value = true
  try {
    const res = await request.get('/api/v1/admin/notices')
    console.log(res)
    list.value = res.list || []
  } finally {
    loading.value = false
  }
}

const openDialog = () => {
  form.content = ''
  form.target_type = 'all'
  form.target_value = ''
  dialogVisible.value = true
}

const submit = async () => {
  if (!form.content) return ElMessage.warning('请输入内容')
  
  try {
    await request.post('/api/v1/admin/send_notification', form)
    ElMessage.success('发送成功')
    dialogVisible.value = false
    fetchList()
  } catch (e) {
    console.error(e)
  }
}

const del = async (row) => {
  await request.delete(`/api/v1/admin/notices/${row.id}`)
  fetchList()
}
</script>