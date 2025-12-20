<template>
  <div class="page-container">
    <!-- 顶部导航 -->
    <div class="nav-bar">
      <button class="back-btn" @click="goBack">← 返回个人中心</button>
      <span class="page-title">消息通知</span>
      <button class="clear-btn" @click="clearAll">全部已读</button>
    </div>

    <!-- 通知列表 -->
    <div class="content-area">
      <div v-if="notifications.length > 0" class="notification-list">
        <div 
          v-for="item in notifications" 
          :key="item.id" 
          class="notification-item"
          :class="{ unread: !item.isRead }"
          @click="markAsRead(item.id)"
        >
          <div class="notification-icon" :class="item.type">
            {{ getIcon(item.type) }}
          </div>
          <div class="notification-content">
            <div class="notification-title">{{ item.title }}</div>
            <p class="notification-text">{{ item.content }}</p>
            <span class="notification-time">{{ formatTime(item.created_at) }}</span>
          </div>
          <div v-if="!item.isRead" class="unread-dot"></div>
        </div>
      </div>

      <div v-else class="empty-state">
        <div class="empty-icon">🔔</div>
        <p class="empty-text">暂无消息通知</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const notifications = ref([
  {
    id: 1,
    type: 'system',
    title: '欢迎来到寝友Go',
    content: '欢迎加入寝友Go社区！这里是高校宿舍互助平台，让我们一起构建温馨的宿舍生活圈。',
    created_at: '2025-12-10T08:00:00',
    isRead: false
  },
  {
    id: 2,
    type: 'like',
    title: '收到新点赞',
    content: '你的帖子《有人能帮忙带个快递吗》收到了3个新点赞',
    created_at: '2025-12-09T15:30:00',
    isRead: true
  },
  {
    id: 3,
    type: 'comment',
    title: '收到新评论',
    content: '用户"热心室友"评论了你的帖子：好的，我待会帮你取！',
    created_at: '2025-12-08T20:15:00',
    isRead: true
  }
])

const goBack = () => {
  router.push('/profile')
}

const getIcon = (type) => {
  const icons = {
    system: '📢',
    like: '❤️',
    comment: '💬',
    reward: '🎁'
  }
  return icons[type] || '📌'
}

const markAsRead = (id) => {
  const item = notifications.value.find(n => n.id === id)
  if (item) {
    item.isRead = true
  }
}

const clearAll = () => {
  notifications.value.forEach(n => n.isRead = true)
}

const formatTime = (timeStr) => {
  if (!timeStr) return ''
  const date = new Date(timeStr)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  
  if (diff < 60000) return '刚刚'
  if (diff < 3600000) return Math.floor(diff / 60000) + '分钟前'
  if (diff < 86400000) return Math.floor(diff / 3600000) + '小时前'
  if (diff < 604800000) return Math.floor(diff / 86400000) + '天前'
  return date.toLocaleDateString('zh-CN')
}

onMounted(() => {
  // 实际应调用 API 获取通知列表
  // fetchNotifications()
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

.back-btn, .clear-btn {
  border: none;
  background: none;
  cursor: pointer;
  font-size: 14px;
  color: #666;
  transition: color 0.3s;
}

.back-btn:hover, .clear-btn:hover {
  color: #1890ff;
}

.page-title {
  font-weight: 600;
  font-size: 16px;
  color: #333;
}

.content-area {
  padding: 12px;
  max-width: 800px;
  margin: 0 auto;
}

.notification-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.notification-item {
  background: white;
  border-radius: 12px;
  padding: 16px;
  display: flex;
  align-items: flex-start;
  gap: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  cursor: pointer;
  transition: all 0.3s;
  position: relative;
}

.notification-item:hover {
  transform: translateX(4px);
}

.notification-item.unread {
  background: #f0f8ff;
  border-left: 3px solid #1890ff;
}

.notification-icon {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  flex-shrink: 0;
}

.notification-icon.system {
  background: #e6f7ff;
}

.notification-icon.like {
  background: #fff0f0;
}

.notification-icon.comment {
  background: #f0fff0;
}

.notification-icon.reward {
  background: #fffbe6;
}

.notification-content {
  flex: 1;
  min-width: 0;
}

.notification-title {
  font-size: 15px;
  font-weight: 600;
  color: #333;
  margin-bottom: 4px;
}

.notification-text {
  font-size: 13px;
  color: #666;
  line-height: 1.5;
  margin: 0 0 8px 0;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.notification-time {
  font-size: 12px;
  color: #999;
}

.unread-dot {
  width: 8px;
  height: 8px;
  background: #ff4d4f;
  border-radius: 50%;
  position: absolute;
  top: 16px;
  right: 16px;
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
}
</style>
