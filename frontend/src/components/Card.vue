<!-- src/components/Card.vue -->
<template>
  <div class="card" @click="handleCardClick">
    <div class="card-header">
      <img :src="post.userAvatar" alt="用户头像" class="user-avatar">
      <div class="user-info">
        <span class="user-name">{{ post.userName }}</span>
        <span class="post-time">{{ post.time }}</span>
      </div>
    </div>
    
    <div class="card-content">
      <!-- 分类标签和宿舍楼标签 -->
      <div class="tags-container">
        <span class="post-category" :class="post.category">{{ getCategoryText(post.category) }}</span>
        <span class="dorm-tag">{{ post.dormBuilding }}</span>
      </div>
      
      <h3 class="post-title">{{ post.title }}</h3>
      <p class="post-content">{{ post.content }}</p>
      <div v-if="post.images && post.images.length > 0" class="post-images">
        <img 
          v-for="(image, index) in post.images.slice(0, 3)" 
          :key="index" 
          :src="image" 
          :alt="'图片' + (index + 1)"
          class="post-image"
        >
        <span v-if="post.images.length > 3" class="image-count">+{{ post.images.length - 3 }}</span>
      </div>
    </div>
    
    <div class="card-footer">
      <div class="post-stats">
        <span class="stat-item">
          <span class="icon">💬</span>
          {{ post.commentCount || 0 }}
        </span>
        <span class="stat-item">
          <span class="icon">👀</span>
          {{ post.viewCount || 0 }}
        </span>
        <span class="stat-item">
          <span class="icon">❤️</span>
          {{ post.likeCount || 0 }}
        </span>
      </div>
      <button class="contact-btn" @click.stop="contactUser">联系TA</button>
    </div>
  </div>
</template>

<script setup>
import { defineProps } from 'vue'

const props = defineProps({
  post: {
    type: Object,
    required: true,
    default: () => ({
      id: '',
      userName: '匿名用户',
      userAvatar: '/default-avatar.png',
      time: '刚刚',
      category: 'help',
      dormBuilding: '榕园9号',
      title: '',
      content: '',
      images: [],
      commentCount: 0,
      viewCount: 0,
      likeCount: 0
    })
  }
})

const categoryMap = {
  'food': '约饭',
  'sports': '约球',
  'help': '求助',
  'trade': '交易',
  'study': '学习'
}

const getCategoryText = (category) => {
  return categoryMap[category] || '其他'
}

const handleCardClick = () => {
  console.log('查看帖子详情:', props.post.id)
}

const contactUser = () => {
  console.log('联系用户:', props.post.userName)
}
</script>

<style scoped>
.card {
  background: white;
  border-radius: 12px;
  padding: 16px;
  margin-bottom: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  cursor: pointer;
  transition: all 0.3s ease;
  border: 1px solid #e8e8e8;
}

.card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
}

.card-header {
  display: flex;
  align-items: center;
  margin-bottom: 12px;
}

.user-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  margin-right: 12px;
  object-fit: cover;
}

.user-info {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.user-name {
  font-weight: 600;
  color: #333;
  font-size: 14px;
}

.post-time {
  font-size: 12px;
  color: #999;
  margin-top: 2px;
}

/* 标签容器样式 */
.tags-container {
  display: flex;
  gap: 8px;
  margin-bottom: 12px;
  flex-wrap: wrap;
}

.post-category {
  padding: 4px 8px;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 500;
}

.post-category.food {
  background: #fff0f0;
  color: #ff4757;
  border: 1px solid #ff4757;
}

.post-category.sports {
  background: #f0f8ff;
  color: #1e90ff;
  border: 1px solid #1e90ff;
}

.post-category.help {
  background: #fff8e1;
  color: #ffa502;
  border: 1px solid #ffa502;
}

.post-category.trade {
  background: #f0fff0;
  color: #2ed573;
  border: 1px solid #2ed573;
}

.post-category.study {
  background: #f0f0ff;
  color: #5352ed;
  border: 1px solid #5352ed;
}

/* 宿舍楼标签样式 */
.dorm-tag {
  padding: 4px 8px;
  background: #f8f9fa;
  color: #666;
  border: 1px solid #e9ecef;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 500;
}

.card-content {
  margin-bottom: 12px;
}

.post-title {
  font-size: 16px;
  font-weight: 600;
  color: #333;
  margin: 0 0 8px 0;
  line-height: 1.4;
}

.post-content {
  font-size: 14px;
  color: #666;
  line-height: 1.5;
  margin: 0 0 12px 0;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.post-images {
  display: flex;
  gap: 8px;
  margin-top: 8px;
}

.post-image {
  width: 80px;
  height: 80px;
  border-radius: 6px;
  object-fit: cover;
  flex-shrink: 0;
}

.image-count {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 80px;
  height: 80px;
  background: #f5f5f5;
  border-radius: 6px;
  color: #999;
  font-size: 14px;
}

.card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 12px;
  border-top: 1px solid #f0f0f0;
}

.post-stats {
  display: flex;
  gap: 16px;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: #999;
}

.icon {
  font-size: 14px;
}

.contact-btn {
  padding: 6px 12px;
  background: #1890ff;
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 12px;
  cursor: pointer;
  transition: background 0.3s;
}

.contact-btn:hover {
  background: #40a9ff;
}
</style>