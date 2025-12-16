<template>
  <div class="app-container">
    <el-card shadow="never">
      <el-tabs v-model="searchForm.status" @tab-change="handleSearch">
        <el-tab-pane label="待审核" name="pending" />
        <el-tab-pane label="进行中" name="normal" />
        <el-tab-pane label="已完成" name="completed" />
        <el-tab-pane label="违规/撤销" name="canceled" />
        <el-tab-pane label="全部" name="all" />
      </el-tabs>

      <div class="filter-bar">
        <div class="left-action">
          <el-input v-model="searchForm.keyword" placeholder="内容/标题关键词" style="width: 200px" clearable />
          <el-button type="primary" @click="handleSearch">筛选</el-button>
        </div>
        <div class="right-action">
          <el-button type="success" @click="openPostDialog">发布官方帖子</el-button>
        </div>
      </div>

      <el-table :data="contentList" border style="margin-top: 20px" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column label="状态" width="120" align="center">
          <template #default="{ row }">
             <el-tag v-if="row.is_pinned" type="danger" effect="dark" size="small">置顶</el-tag>
             <el-tag :type="getStatusType(row.status)" size="small" style="margin-left: 5px">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="标题/内容" min-width="250">
          <template #default="{ row }">
            <div class="content-wrapper">
              <div class="post-title">{{ row.title }}</div>
              <div class="post-preview">{{ row.content }}</div>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="类型" width="100" prop="type_name" align="center" />
        <el-table-column label="发布者" width="150">
          <template #default="{ row }">
            <span>{{ row.publisher_name }}</span><br>
            <span style="font-size:12px; color:#999">{{ row.dorm_name }}</span>
          </template>
        </el-table-column>
        <el-table-column label="时间" width="160" align="center">
           <template #default="{ row }">{{ new Date(row.created_at).toLocaleDateString() }}</template>
        </el-table-column>
        <el-table-column label="操作" width="220" fixed="right" align="center">
          <template #default="{ row }">
            <el-button 
              size="small" 
              :type="row.is_pinned ? 'warning' : 'primary'" 
              @click="handlePin(row)"
            >
              {{ row.is_pinned ? '取消置顶' : '置顶' }}
            </el-button>
            <el-button 
              size="small" 
              type="danger" 
              @click="handleAudit(row, 'canceled')" 
              v-if="row.status !== 'canceled'"
            >
              违规
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-container">
        <el-pagination 
          v-model:current-page="pageInfo.page" 
          :total="total" 
          layout="total, prev, pager, next"
          @current-change="fetchContentList"
        />
      </div>
    </el-card>

    <el-dialog v-model="postDialogVisible" title="发布官方/置顶帖子" width="500px">
      <el-form :model="postForm" label-width="80px">
        <el-form-item label="标题">
          <el-input v-model="postForm.title" />
        </el-form-item>
        <el-form-item label="内容">
          <el-input type="textarea" v-model="postForm.content" rows="4" />
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="postForm.typeid" placeholder="请选择">
            <el-option label="闲置交易" :value="1" />
            <el-option label="跑腿求助" :value="2" />
            <el-option label="系统公告" :value="99" /> 
            </el-select>
        </el-form-item>
        <el-form-item label="选项">
          <el-checkbox v-model="postForm.is_pinned">直接置顶</el-checkbox>
        </el-form-item>
        <el-form-item label="图片">
          <input type="file" @change="handleFileChange" accept="image/*" multiple />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="postDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitAdminPost">发布</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, reactive } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/request'

const contentList = ref([])
const total = ref(0)
const loading = ref(false)
const searchForm = reactive({ keyword: '', status: 'normal' })
const pageInfo = reactive({ page: 1, pageSize: 10 })

// 发帖相关
const postDialogVisible = ref(false)
const postForm = reactive({ title: '', content: '', typeid: 2, is_pinned: false, images: [] })

onMounted(() => fetchContentList())

const fetchContentList = async () => {
  loading.value = true
  try {
    const res = await request.get('/api/v1/admin/contents', { params: { ...pageInfo, ...searchForm } })
    contentList.value = res.list || []
    total.value = res.total || 0
  } catch (e) { console.error(e) } 
  finally { loading.value = false }
}

const handleSearch = () => { pageInfo.page = 1; fetchContentList() }

const handleAudit = async (row, newStatus) => {
  try {
    await ElMessageBox.confirm('确定标记为违规下架吗？', '警告', { type: 'warning' })
    await request.put(`/api/v1/admin/contents/${row.id}/audit`, { status: newStatus })
    ElMessage.success('操作成功')
    fetchContentList()
  } catch (e) {}
}

const handlePin = async (row) => {
  try {
    const newStatus = !row.is_pinned
    await request.put(`/api/v1/admin/contents/${row.id}/pin`, { is_pinned: newStatus })
    ElMessage.success(newStatus ? '置顶成功' : '已取消置顶')
    fetchContentList()
  } catch (e) {}
}

const getStatusType = (status) => {
  const map = { pending: 'info', normal: 'success', canceled: 'danger', completed: '' }
  return map[status] || ''
}

// 官方发帖逻辑
const openPostDialog = () => {
  postForm.title = ''
  postForm.content = ''
  postForm.images = []
  postDialogVisible.value = true
}

const handleFileChange = (e) => {
  postForm.images = Array.from(e.target.files)
}

const submitAdminPost = async () => {
  if(!postForm.title || !postForm.content) return ElMessage.warning('标题和内容必填')
  
  const formData = new FormData()
  formData.append('title', postForm.title)
  formData.append('content', postForm.content)
  formData.append('typeid', postForm.typeid)
  formData.append('is_pinned', postForm.is_pinned)
  postForm.images.forEach(file => formData.append('images', file))

  try {
    await request.post('/api/v1/admin/contents/create', formData)
    ElMessage.success('发布成功')
    postDialogVisible.value = false
    fetchContentList()
  } catch(e) {}
}
</script>

<style scoped>
.app-container { padding: 20px; }
.filter-bar { display: flex; justify-content: space-between; margin-bottom: 20px; }
.left-action { display: flex; gap: 10px; }
.post-title { font-weight: bold; margin-bottom: 4px; }
.post-preview { color: #666; font-size: 13px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.pagination-container { margin-top: 20px; text-align: right; }
</style>