<template>
  <div class="app-container">
    <el-card shadow="never">
      <el-tabs v-model="searchForm.status" @tab-change="handleSearch">
        <el-tab-pane label="待审核" name="pending" />
        <el-tab-pane label="进行中" name="accepted" />
        <el-tab-pane label="已完成" name="completed" />
        <el-tab-pane label="违规/撤销" name="canceled" />
      </el-tabs>

      <div class="filter-bar">
        <el-input v-model="searchForm.keyword" placeholder="内容/标题关键词" style="width: 200px" clearable />
        <el-select v-model="searchForm.type_id" placeholder="互助类型" clearable style="width: 140px">
          <el-option v-for="type in postTypes" :key="type.id" :label="type.name" :value="type.id" />
        </el-select>
        <el-button type="primary" @click="handleSearch">筛选</el-button>
      </div>

      <el-table :data="contentList" border style="margin-top: 20px" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" align="center" />
        
        <el-table-column label="互助类型" width="120" align="center">
          <template #default="{ row }">
            <el-tag effect="plain">{{ row.type_name }}</el-tag>
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

        <el-table-column prop="reward_type" label="酬劳" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="getRewardTagType(row.reward_type)">{{ row.reward_type }}</el-tag>
          </template>
        </el-table-column>

        <el-table-column label="发布者" width="150">
          <template #default="{ row }">
            <span>{{ row.publisher_name }}</span>
            <br>
            <span style="font-size:12px; color:#999">{{ row.dorm_name }}</span>
          </template>
        </el-table-column>

        <el-table-column prop="expires_at" label="有效期至" width="160" sortable>
          <template #default="{ row }">
            <span :class="{ 'text-expired': isExpired(row.expires_at) }">
              {{ formatTime(row.expires_at) }}
            </span>
          </template>
        </el-table-column>

        <el-table-column label="操作" width="200" fixed="right" align="center">
          <template #default="{ row }">
            <el-button size="small" @click="handleView(row)">详情</el-button>
            
            <template v-if="row.status === 'pending'">
               <el-button size="small" type="danger" @click="handleAudit(row, 'canceled')">违规下架</el-button>
            </template>
            <template v-else-if="row.status === 'accepted'">
               <el-button size="small" type="warning" @click="handleAudit(row, 'canceled')">强制取消</el-button>
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

    <el-dialog v-model="dialogVisible" title="互助详情" width="600px">
      <el-descriptions :column="2" border v-if="currentContent">
        <el-descriptions-item label="标题" :span="2">{{ currentContent.title }}</el-descriptions-item>
        <el-descriptions-item label="类型">{{ currentContent.type_name }}</el-descriptions-item>
        <el-descriptions-item label="酬劳">{{ currentContent.reward_type }}</el-descriptions-item>
        <el-descriptions-item label="发布者">{{ currentContent.publisher_name }}</el-descriptions-item>
        <el-descriptions-item label="发布时间">{{ currentContent.created_at }}</el-descriptions-item>
        <el-descriptions-item label="截止时间">{{ formatTime(currentContent.expires_at) }}</el-descriptions-item>
        <el-descriptions-item label="当前状态">
            <el-tag>{{ currentContent.status }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="详细内容" :span="2">
          <div style="white-space: pre-wrap;">{{ currentContent.content }}</div>
        </el-descriptions-item>
      </el-descriptions>
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
const postTypes = ref([]) // 互助类型列表
const dialogVisible = ref(false)
const currentContent = ref(null)

const searchForm = reactive({ keyword: '', type_id: '', status: 'pending' })
const pageInfo = reactive({ page: 1, pageSize: 10 })

onMounted(() => {
  fetchPostTypes()
  fetchContentList()
})

// 模拟获取互助类型，实际应调用 /api/v1/config/types [cite: 64]
const fetchPostTypes = () => {
  postTypes.value = [
    { id: 1, name: '约球' },
    { id: 2, name: '带东西' },
    { id: 3, name: '找搭子' }
  ]
}

const fetchContentList = async () => {
  loading.value = true
  try {
    const res = await request.get('/v1/admin/contents', { params: { ...pageInfo, ...searchForm } })
    contentList.value = res.list || []
    total.value = res.total || 0
  } catch (e) { console.error(e) } 
  finally { loading.value = false }
}

const getRewardTagType = (type) => {
  if (type === '无偿') return 'success'
  if (type === '现金') return 'warning'
  return ''
}

const formatTime = (time) => time ? time.replace('T', ' ').substring(0, 16) : '长期有效'
const isExpired = (time) => time && new Date(time) < new Date()

const handleSearch = () => { pageInfo.page = 1; fetchContentList() }

const handleView = (row) => {
  currentContent.value = row
  dialogVisible.value = true
}

const handleAudit = async (row, newStatus) => {
  try {
    await ElMessageBox.confirm('确定将该互助请求标记为违规/取消吗？', '警告', { type: 'warning' })
    await request.put(`/v1/admin/contents/${row.id}/audit`, { status: newStatus })
    ElMessage.success('操作成功')
    fetchContentList()
  } catch (e) {}
}
</script>

<style scoped>
.app-container { padding: 20px; }
.filter-bar { display: flex; gap: 10px; margin-bottom: 20px; }
.post-title { font-weight: bold; margin-bottom: 4px; }
.post-preview { color: #666; font-size: 13px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.text-expired { color: #f56c6c; text-decoration: line-through; }
.pagination-container { margin-top: 20px; text-align: right; }
</style>