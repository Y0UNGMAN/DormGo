<template>
  <el-dialog
    v-model="visible"
    title="通知详情"
    width="600px"
    destroy-on-close
    :before-close="handleClose"
  >
    <div v-if="noticeData" class="notice-detail">
      <div class="header">
        <h3 class="title">
          <el-tag v-if="noticeData.is_urgent" type="danger" size="small" effect="dark">紧急</el-tag>
          {{ noticeData.title }}
        </h3>
        <div class="meta">
          <span>发布人：{{ noticeData.publisher_name || '系统管理员' }}</span>
          <span>时间：{{ noticeData.created_at }}</span>
        </div>
      </div>
      
      <el-divider />
      
      <div class="content-body">
        <div v-html="noticeData.content" class="rich-text"></div>
      </div>
      
      <el-divider content-position="left">接收对象信息</el-divider>
      <div class="target-info">
        <el-descriptions :column="1" border size="small">
          <el-descriptions-item label="接收范围">
            {{ formatTarget(noticeData.target) }}
          </el-descriptions-item>
          <el-descriptions-item label="详细名单" v-if="noticeData.target !== 'all'">
            {{ noticeData.target_info || noticeData.target_id }}
          </el-descriptions-item>
          <el-descriptions-item label="阅读情况">
            <el-progress 
              :percentage="readPercentage" 
              :format="() => `${noticeData.read_count}/${noticeData.total_count} 已读`"
            />
          </el-descriptions-item>
        </el-descriptions>
      </div>
    </div>
    
    <template #footer>
      <el-button @click="copyContent">复制内容</el-button>
      <el-button type="primary" @click="handleClose">关闭</el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'
import useClipboard from 'vue-clipboard3'

const props = defineProps({
  modelValue: Boolean,
  data: Object
})

const emit = defineEmits(['update:modelValue'])
const { toClipboard } = useClipboard()

const visible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const noticeData = computed(() => props.data)

const readPercentage = computed(() => {
  if (!noticeData.value || !noticeData.value.total_count) return 0
  return Math.round((noticeData.value.read_count / noticeData.value.total_count) * 100)
})

const formatTarget = (target) => {
  const map = { all: '全体成员', dorm: '指定楼栋', user: '指定用户' }
  return map[target] || target
}

const handleClose = () => {
  visible.value = false
}

const copyContent = async () => {
  try {
    await toClipboard(noticeData.value.content)
    ElMessage.success('内容已复制到剪贴板')
  } catch (e) {
    ElMessage.error('复制失败')
  }
}
</script>

<style scoped>
.notice-detail .header { margin-bottom: 20px; }
.notice-detail .title { margin: 0 0 10px 0; display: flex; align-items: center; gap: 8px; }
.notice-detail .meta { color: #909399; font-size: 12px; display: flex; gap: 20px; }
.content-body { min-height: 100px; line-height: 1.6; font-size: 14px; color: #303133; }
</style>
