<template>
  <div class="page-container">
    <!-- 顶部导航 -->
    <div class="nav-bar">
      <button class="back-btn" @click="goBack">← 返回个人中心</button>
      <span class="page-title">发帖历史</span>
      <span class="placeholder"></span>
    </div>

    <!-- 帖子列表 -->
    <div class="content-area">
      <div v-if="posts.length > 0" class="post-list">
        <div v-for="post in posts" :key="post.id" class="post-item" @click="goToDetail(post.id)">
          <div class="post-header">
            <span class="post-title">{{ post.title }}</span>
            <span class="post-time">{{ formatTime(post.created_at) }}</span>
          </div>
          <p class="post-content">{{ post.content }}</p>
          <div class="post-footer">
            <span class="post-tag">{{ post.typeName || '分享' }}</span>
            <div class="post-stats">
              <span>👁 {{ post.view_count || 0 }}</span>
              <span>💬 {{ post.comment_count || 0 }}</span>
              <span>❤️ {{ post.like_count || 0 }}</span>
            </div>
          </div>
        </div>
      </div>

      <div v-else class="empty-state">
        <div class="empty-icon">📝</div>
        <p class="empty-text">还没有发布过帖子</p>
        <button class="publish-btn" @click="goPublish">去发布</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const posts = ref([
  // 模拟数据
  {
    id: 1,
    title: '【求助】有人能帮忙带个快递吗',
    content: '菜鸟驿站的快递，我今天有课没时间取...',
    typeName: '求助',
    created_at: '2025-12-08T10:30:00',
    view_count: 56,
    comment_count: 3,
    like_count: 2
  },
  {
    id: 2,
    title: '【约饭】今晚有人一起吃火锅吗',
    content: '想去学校后门那家重庆火锅，有没有小伙伴...',
    typeName: '约饭',
    created_at: '2025-12-05T18:20:00',
    view_count: 128,
    comment_count: 12,
    like_count: 8
  }
])

const goBack = () => {
  router.push('/profile')
}

const goToDetail = (id) => {
  router.push(`/post/${id}`)
}

const goPublish = () => {
  router.push('/publish')
}

const formatTime = (timeStr) => {
  if (!timeStr) return ''
  const date = new Date(timeStr)
  return date.toLocaleDateString('zh-CN')
}

onMounted(() => {
  // 实际应调用 API 获取用户发帖历史
  // fetchUserPosts()
})
</script>

<style scoped>
.page-container {
  min-height: 100vh;
  background: #f5f5f5;
}

.nav-bar {
  background: white;
  padding: 15px 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  box-shadow: 0 1px 4px rgba(0,0,0,0.05);
  position: sticky;
  top: 0;
  z-index: 100;
}

.back-btn {
  border: none;
  background: none;
  cursor: pointer;
  font-size: 14px;
  color: #666;
  transition: color 0.3s;
}

.back-btn:hover {
  color: #1890ff;
}

.page-title {
  font-weight: 600;
  font-size: 16px;
  color: #333;
}

.placeholder {
  width: 100px;
}

.content-area {
  padding: 12px;
  max-width: 800px;
  margin: 0 auto;
}

.post-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.post-item {
  background: white;
  border-radius: 12px;
  padding: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  cursor: pointer;
  transition: all 0.3s;
}

.post-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.post-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 8px;
}

.post-title {
  font-size: 16px;
  font-weight: 600;
  color: #333;
  flex: 1;
}

.post-time {
  font-size: 12px;
  color: #999;
  white-space: nowrap;
  margin-left: 10px;
}

.post-content {
  font-size: 14px;
  color: #666;
  line-height: 1.5;
  margin: 0 0 12px 0;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.post-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.post-tag {
  background: #e6f7ff;
  color: #1890ff;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
}

.post-stats {
  display: flex;
  gap: 12px;
  font-size: 12px;
  color: #999;
}

.empty-state {
  text-align: center;
  padding: 80px 20px;
  background: white;
  border-radius: 12px;
  margin-top: 20px;
}

.empty-icon {
  font-size: 48px;
  margin-bottom: 16px;
}

.empty-text {
  font-size: 16px;
  color: #999;
  margin-bottom: 20px;
}

.publish-btn {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  padding: 10px 30px;
  border-radius: 20px;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.3s;
}

.publish-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
}
</style>
