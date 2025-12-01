<template>
  <div class="app-container">
    <el-card shadow="never">
      <div class="toolbar">
        <el-button type="primary" @click="openAddDialog">
          <el-icon><Plus /></el-icon> 发布通知
        </el-button>
        
        <div class="search-group">
          <el-input 
            v-model="searchForm.keyword" 
            placeholder="标题关键词" 
            clearable 
            style="width: 200px; margin-right: 10px;" 
            @keyup.enter="handleSearch"
          />
          <el-select v-model="searchForm.target" placeholder="通知对象" clearable style="width: 120px; margin-right: 10px;">
            <el-option label="全体" value="all" />
            <el-option label="楼栋" value="dorm" />
            <el-option label="个人" value="user" />
          </el-select>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="resetSearch">重置</el-button>
        </div>
      </div>

      <el-table :data="noticeList" border style="margin-top: 20px" stripe v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column prop="title" label="标题" min-width="200" show-overflow-tooltip>
          <template #default="{ row }">
            <span v-if="row.is_urgent" style="color: #f56c6c; margin-right: 5px;">[紧急]</span>
            <span v-if="row.is_pinned" style="color: #409EFF; margin-right: 5px;">[置顶]</span>
            {{ row.title }}
          </template>
        </el-table-column>
        <el-table-column prop="target" label="对象" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="row.target === 'all' ? 'danger' : 'info'">
              {{ formatTarget(row) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="target_info" label="对象详情" width="150" show-overflow-tooltip />
        <el-table-column prop="created_at" label="发布时间" width="170" />
        <el-table-column label="操作" width="150" fixed="right" align="center">
          <template #default="{ row }">
            <el-button link type="primary" size="small" @click="handleView(row)">查看</el-button>
            <el-button link type="danger" size="small" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-container">
        <el-pagination
          v-model:current-page="pageInfo.page"
          v-model:page-size="pageInfo.pageSize"
          :total="total"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <el-dialog title="发布新通知" v-model="showAddDialog" width="600px" destroy-on-close>
      <el-form :model="noticeForm" label-position="top">
        <el-form-item label="标题" required>
          <el-input v-model="noticeForm.title" placeholder="请输入通知标题" />
        </el-form-item>
        
        <el-form-item label="特殊属性">
          <el-checkbox v-model="noticeForm.is_pinned">置顶显示</el-checkbox>
          <el-checkbox v-model="noticeForm.is_urgent">紧急通知</el-checkbox>
        </el-form-item>

        <el-form-item label="内容" required>
          <el-input type="textarea" :rows="5" v-model="noticeForm.content" placeholder="请输入通知详情" />
        </el-form-item>

        <el-form-item label="发送对象" required>
          <el-radio-group v-model="noticeForm.target">
            <el-radio-button label="all">全体用户</el-radio-button>
            <el-radio-button label="dorm">指定楼栋</el-radio-button>
            <el-radio-button label="user">指定个人</el-radio-button>
          </el-radio-group>
        </el-form-item>
        
        <el-form-item v-if="noticeForm.target === 'dorm'" label="选择楼栋">
          <el-select v-model="noticeForm.target_id" multiple placeholder="请选择楼栋" style="width: 100%">
            <el-option v-for="dorm in dormList" :key="dorm.id" :label="dorm.name" :value="dorm.id" />
          </el-select>
        </el-form-item>
        <el-form-item v-if="noticeForm.target === 'user'" label="用户ID">
          <el-input v-model="noticeForm.target_id" placeholder="输入用户ID，逗号分隔" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showAddDialog = false">取消</el-button>
        <el-button type="primary" @click="handleAddNotice">确认发布</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, reactive } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import request from '@/utils/request'

const noticeList = ref([])
const dormList = ref([])
const total = ref(0)
const loading = ref(false)
const showAddDialog = ref(false)

const searchForm = reactive({ keyword: '', target: '' })
const pageInfo = reactive({ page: 1, pageSize: 10 })
const noticeForm = reactive({ 
  title: '', 
  content: '', 
  target: 'all', 
  target_id: [],
  is_pinned: false,
  is_urgent: false
})

onMounted(() => {
  fetchNoticeList()
  fetchDormList()
})

const fetchNoticeList = async () => {
  loading.value = true
  try {
    const res = await request.get('/v1/admin/notices', {
      params: { ...pageInfo, ...searchForm }
    })
    noticeList.value = res.list || []
    total.value = res.total || 0
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const fetchDormList = async () => {
  try {
    const res = await request.get('/v1/admin/dorms')
    dormList.value = res.data || [] // 假设接口返回 list
  } catch (e) { console.error(e) }
}

const formatTarget = (row) => ({ all: '全体', dorm: '楼栋', user: '个人' }[row.target] || row.target)

const handleSearch = () => { pageInfo.page = 1; fetchNoticeList() }
const resetSearch = () => { searchForm.keyword = ''; searchForm.target = ''; handleSearch() }
const handleSizeChange = (val) => { pageInfo.pageSize = val; fetchNoticeList() }
const handleCurrentChange = (val) => { pageInfo.page = val; fetchNoticeList() }

const openAddDialog = () => {
  // 重置表单
  noticeForm.title = ''
  noticeForm.content = ''
  noticeForm.target = 'all'
  noticeForm.target_id = []
  noticeForm.is_pinned = false
  noticeForm.is_urgent = false
  showAddDialog.value = true
}

const handleView = (row) => {
  ElMessageBox.alert(row.content, row.title, { confirmButtonText: '关闭', dangerouslyUseHTMLString: true })
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确定删除该通知？', '提示', { type: 'warning' })
    await request.delete(`/v1/admin/notices/${row.id}`)
    ElMessage.success('删除成功')
    fetchNoticeList()
  } catch (e) {}
}

const handleAddNotice = async () => {
  if (!noticeForm.title || !noticeForm.content) return ElMessage.warning('请完善信息')
  
  // 处理 target_id 格式，如果是数组转字符串或保持数组取决于后端需求，这里假设转字符串
  const payload = { ...noticeForm }
  if (Array.isArray(payload.target_id)) {
      payload.target_id = payload.target_id.join(',')
  }

  try {
    await request.post('/v1/admin/notices', payload)
    ElMessage.success('发布成功')
    showAddDialog.value = false
    fetchNoticeList()
  } catch (e) { 
    // request.ts 已处理错误提示
  }
}
</script>

<style scoped>
.app-container { padding: 20px; }
.toolbar { display: flex; justify-content: space-between; align-items: center; }
.pagination-container { margin-top: 20px; text-align: right; }
</style>