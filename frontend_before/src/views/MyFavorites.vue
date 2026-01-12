<template>
  <div class="page-container">
    <!-- 顶部导航 -->
    <div class="nav-bar">
      <button class="back-btn" @click="goBack">← 返回个人中心</button>
      <span class="page-title">我的收藏</span>
      <span class="placeholder"></span>
    </div>

    <!-- 收藏列表 -->
    <div class="content-area">
      <div v-if="favorites.length > 0" class="post-list">
        <div v-for="post in favorites" :key="post.id" class="post-item" @click="goToDetail(post.id)">
          <div class="post-header">
            <span class="post-title">{{ post.title }}</span>
            <button class="unfavorite-btn" @click.stop="removeFavorite(post.id)">取消收藏</button>
          </div>
          <p class="post-content">{{ post.content }}</p>
          <div class="post-footer">
            <div class="post-author">
              <img :src="post.authorAvatar" class="author-avatar" />
              <span class="author-name">{{ post.authorName }}</span>
            </div>
            <span class="post-time">{{ formatTime(post.created_at) }}</span>
          </div>
        </div>
      </div>

      <div v-else class="empty-state">
        <div class="empty-icon">⭐</div>
        <p class="empty-text">还没有收藏任何帖子</p>
        <button class="explore-btn" @click="goHome">去逛逛</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const favorites = ref([
  // 模拟数据
  {
    id: 10,
    title: '【分享】期末复习资料整理',
    content: '整理了这学期高数和线代的复习笔记，需要的同学可以来找我...',
    authorName: '学霸小李',
    authorAvatar: 'https://dorm-go.oss-cn-guangzhou.aliyuncs.com/avator/1.jpg',
    created_at: '2025-12-06T14:30:00'
  },
  {
    id: 11,
    title: '【交易】二手台灯便宜出',
    content: '买了两个月的护眼台灯，因为换宿舍用不上了...',
    authorName: '隔壁老王',
    authorAvatar: 'https://dorm-go.oss-cn-guangzhou.aliyuncs.com/avator/2.jpg',
    created_at: '2025-12-01T09:15:00'
  }
])

const goBack = () => {
  router.push('/profile')
}

const goToDetail = (id) => {
  router.push(`/post/${id}`)
}

const goHome = () => {
  router.push('/home')
}

const removeFavorite = (id) => {
  if (confirm('确定取消收藏吗？')) {
    favorites.value = favorites.value.filter(p => p.id !== id)
  }
}

const formatTime = (timeStr) => {
  if (!timeStr) return ''
  const date = new Date(timeStr)
  return date.toLocaleDateString('zh-CN')
}

onMounted(() => {
  // 实际应调用 API 获取用户收藏列表
  // fetchFavorites()
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

.unfavorite-btn {
  background: none;
  border: 1px solid #ff4d4f;
  color: #ff4d4f;
  padding: 4px 10px;
  border-radius: 4px;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.3s;
}

.unfavorite-btn:hover {
  background: #ff4d4f;
  color: white;
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

.post-author {
  display: flex;
  align-items: center;
  gap: 8px;
}

.author-avatar {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  object-fit: cover;
}

.author-name {
  font-size: 13px;
  color: #666;
}

.post-time {
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

.explore-btn {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  padding: 10px 30px;
  border-radius: 20px;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.3s;
}

.explore-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
}
</style>
