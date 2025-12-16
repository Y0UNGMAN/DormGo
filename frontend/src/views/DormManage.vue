<template>
  <div class="app-container">
    <el-card shadow="never">
      <el-button type="primary" @click="handleAdd" style="margin-bottom: 20px">新增楼栋</el-button>
      <el-table :data="list" border stripe>
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column prop="dormname" label="楼栋名称" />
        <el-table-column prop="user_count" label="入住人数" align="center" />
        <el-table-column label="操作" width="150" align="center">
          <template #default="{ row }">
            <el-button link type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
    
    <el-dialog v-model="dialogVisible" title="新增楼栋" width="400px">
      <el-input v-model="form.dormname" placeholder="请输入楼栋名称 (如: C10栋)" />
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, reactive } from 'vue'
import request from '@/utils/request'
import { ElMessage, ElMessageBox } from 'element-plus'

const list = ref([])
const dialogVisible = ref(false)
// 修改点：使用 dormname 属性
const form = reactive({ dormname: '' })

onMounted(() => fetchList())

const fetchList = async () => {
  try {
    const res = await request.get('/api/v1/admin/dorms')
    list.value = res.list || []
  } catch (e) {}
}

const handleAdd = () => {
  form.dormname = ''
  dialogVisible.value = true
}

const submit = async () => {
  if(!form.dormname) return ElMessage.warning('名称不能为空')
  // 发送给后端的 JSON 字段将是 { dormname: "..." }
  await request.post('/api/v1/admin/dorms', form)
  ElMessage.success('添加成功')
  dialogVisible.value = false
  fetchList()
}

const handleDelete = async (row) => {
  await ElMessageBox.confirm('确定删除吗？')
  await request.delete(`/api/v1/admin/dorms/${row.id}`)
  ElMessage.success('删除成功')
  fetchList()
}
</script>

<style scoped>
.app-container { padding: 20px; }
</style>