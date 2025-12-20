<template>
  <div class="app-container">
    <el-card shadow="never">
      <el-tabs v-model="searchForm.status" @tab-change="handleSearch">
        <el-tab-pane label="全部帖子" name="all" />
        <el-tab-pane label="违规/撤销" name="canceled" />
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
             <el-tag :type="getStatusType(row.status)" size="small" style="margin-left: 5px">
               {{ getStatusText(row.status) }}
             </el-tag>
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
            
            <template v-if="row.status !== 'canceled'">
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
              >
                违规
              </el-button>
            </template>

            <template v-else>
              <el-button 
                size="small" 
                type="success" 
                @click="handleAudit(row, 'normal')" 
              >
                恢复
              </el-button>
            </template>

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
            <el-option v-for="item in postTypes" :key="item.typeid" :label="item.typename" :value="item.typeid" />
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
// 默认选 all
const searchForm = reactive({ keyword: '', status: 'all' })
const pageInfo = reactive({ page: 1, pageSize: 10 })

// 发帖相关
const postDialogVisible = ref(false)
const postForm = reactive({ title: '', content: '', typeid: '', is_pinned: false, images: [] })
const postTypes = ref([])

onMounted(() => {
  fetchContentList()
  fetchPostTypes()
})

const fetchPostTypes = async () => {
  try {
    const res = await request.get('/api/v1/post/post_type')
    if (res.code === 200) {
      postTypes.value = res.data || []
    }
  } catch (e) { console.error(e) }
}

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

// 审核/状态变更逻辑
const handleAudit = async (row, newStatus) => {
  let confirmText = ''
  if (newStatus === 'canceled') confirmText = '确定标记为违规并下架吗？'
  if (newStatus === 'normal') confirmText = '确定恢复该帖子吗？'
  
  try {
    await ElMessageBox.confirm(confirmText, '提示', { type: 'warning' })
    // 调用后端的审核接口
    await request.put(`/api/v1/admin/contents/${row.id}/audit`, { status: newStatus })
    ElMessage.success('操作成功')
    
    // 如果当前是在“全部”标签，可能需要手动刷新列表状态
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
  return map[status] || 'info'
}

const getStatusText = (status) => {
  const map = { pending: '待审', normal: '正常', canceled: '违规', completed: '完成' }
  return map[status] || status
}

// 官方发帖逻辑
const openPostDialog = () => {
  postForm.title = ''
  postForm.content = ''
  postForm.typeid = '' 
  postForm.is_pinned = false
  postForm.images = []
  postDialogVisible.value = true
}

const handleFileChange = (e) => {
  postForm.images = Array.from(e.target.files)
}

const submitAdminPost = async () => {
  if(!postForm.title || !postForm.content) return ElMessage.warning('标题和内容必填')
  if(!postForm.typeid) return ElMessage.warning('请选择帖子类型')
  
  const formData = new FormData()
  formData.append('title', postForm.title)
  formData.append('content', postForm.content)
  formData.append('typeid', postForm.typeid)
  formData.append('is_pinned', postForm.is_pinned)
  postForm.images.forEach(file => formData.append('images', file))

  try {
    const res = await request.post('/api/v1/admin/contents/create', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
    if (res.code === 200) {
        ElMessage.success('发布成功')
        postDialogVisible.value = false
        fetchContentList()
    } else {
        ElMessage.error(res.msg || '发布失败')
    }
  } catch(e) {
    ElMessage.error('网络错误')
  }
}
</script>

<style scoped>
.app-container { padding: 20px; }
.filter-bar { display: flex; justify-content: space-between; margin-bottom: 20px; }
.left-action { display: flex; gap: 10px; }
.post-title { font-weight: bold; margin-bottom: 4px; }
.post-preview { color: #666; font-size: 13px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.pagination-container { margin-top: 20px; text-align: right; }
.content-wrapper { display: flex; flex-direction: column; }
</style>