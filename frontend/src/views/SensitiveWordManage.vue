<template>
  <div class="app-container">
    <el-card shadow="never">
      <div class="toolbar">
        <div class="left">
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon> 新增敏感词
          </el-button>
          <el-button type="success" @click="handleImport">
            <el-icon><Upload /></el-icon> 批量导入
          </el-button>
          <el-button type="warning" @click="handleExport">
            <el-icon><Download /></el-icon> 导出词库
          </el-button>
        </div>
        <div class="right">
          <el-input v-model="keyword" placeholder="搜索敏感词" clearable @input="filterList" style="width: 200px">
            <template #prefix><el-icon><Search /></el-icon></template>
          </el-input>
        </div>
      </div>

      <el-table :data="displayList" border style="margin-top: 20px" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="word" label="敏感词" />
        <el-table-column prop="category" label="分类" width="120">
          <template #default="{ row }">
            <el-tag>{{ row.category || '通用' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="添加时间" width="180" />
        <el-table-column label="操作" width="150" align="center">
          <template #default="{ row }">
            <el-button link type="primary" @click="handleEdit(row)">编辑</el-button>
            <el-button link type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 编辑弹窗 -->
    <el-dialog :title="dialogTitle" v-model="dialogVisible" width="400px">
      <el-form :model="form" label-width="80px">
        <el-form-item label="敏感词">
          <el-input v-model="form.word" placeholder="请输入敏感词" />
        </el-form-item>
        <el-form-item label="分类">
          <el-select v-model="form.category" placeholder="请选择分类" style="width: 100%">
            <el-option label="涉黄" value="porn" />
            <el-option label="暴力" value="violence" />
            <el-option label="广告" value="ad" />
            <el-option label="其他" value="other" />
          </el-select>
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
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Upload, Download, Search } from '@element-plus/icons-vue'
import request from '@/utils/request'

const loading = ref(false)
const list = ref([])
const keyword = ref('')
const dialogVisible = ref(false)
const isEdit = ref(false)
const form = reactive({ id: null, word: '', category: 'other' })

const dialogTitle = computed(() => isEdit.value ? '编辑敏感词' : '新增敏感词')

// 前端过滤，减少后端压力
const displayList = computed(() => {
  if (!keyword.value) return list.value
  return list.value.filter(item => item.word.includes(keyword.value))
})

onMounted(() => fetchList())

const fetchList = async () => {
  loading.value = true
  try {
    const res = await request.get('/v1/admin/config/sensitive-words')
    list.value = res.list || []
  } catch (e) {} finally { loading.value = false }
}

const handleAdd = () => {
  isEdit.value = false
  form.id = null
  form.word = ''
  form.category = 'other'
  dialogVisible.value = true
}

const handleEdit = (row) => {
  isEdit.value = true
  Object.assign(form, row)
  dialogVisible.value = true
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(`确定删除敏感词 "${row.word}" 吗？`)
    await request.delete(`/v1/admin/config/sensitive-words/${row.id}`)
    ElMessage.success('删除成功')
    fetchList()
  } catch (e) {}
}

const submitForm = async () => {
  if (!form.word) return ElMessage.warning('请输入内容')
  try {
    if (isEdit.value) {
      await request.put(`/v1/admin/config/sensitive-words/${form.id}`, form)
    } else {
      await request.post('/v1/admin/config/sensitive-words', form)
    }
    ElMessage.success('保存成功')
    dialogVisible.value = false
    fetchList()
  } catch (e) {}
}

const handleImport = () => ElMessage.info('批量导入功能开发中...')
const handleExport = () => ElMessage.info('导出功能开发中...')
</script>

<style scoped>
.app-container { padding: 20px; }
.toolbar { display: flex; justify-content: space-between; margin-bottom: 20px; }
</style>