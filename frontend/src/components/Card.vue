<!-- src/components/Card.vue -->
<template>
  <div class="card" @click="handleCardClick">
    <!-- 标题放在最上部 -->
    <h3 class="post-title">
      {{ post.title }}
    </h3>
    
    <div class="card-content">
      <!-- 用户信息：头像和名字 -->
      <div class="user-info">
        <img :src="post.userAvatar" alt="用户头像" class="user-avatar">
        <span class="user-name">
          {{ post.userName }}
        </span>
      </div>

      <!-- 分类标签和宿舍楼标签 -->
      <div class="tags-container">
        <span class="post-category" :class="post.category">
          {{ getCategoryText(post.category) }}
        </span>
        <span class="dorm-tag">
          {{ post.dormBuilding }}
        </span>
      </div>
      
      <!-- 帖子正文 - 只显示一行 -->
      <p class="post-content">
        {{ post.content }}
      </p>
      
      <!-- 图片展示 - 最多两张 -->
      <div v-if="post.images && post.images.length > 0" class="post-images">
        <img 
          v-for="(image, index) in post.images.slice(0, 2)" 
          :key="index" 
          :src="image" 
          :alt="'图片' + (index + 1)"
          class="post-image"
        >
      </div>
    </div>
    
    <!-- 使用独立的统计组件 -->
    <PostStats
      :view-count="post.viewCount"
      :comment-count="post.commentCount"
      :like-count="post.likeCount"
      :time="post.time"
    />
  </div>
</template>

<script setup>
import { defineProps } from 'vue'
import { useRouter } from 'vue-router'
import PostStats from '@/components/PostStats.vue'

const router = useRouter()
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
  // 增加浏览量
  props.post.viewCount++
  router.push(`/post/${props.post.id}`)
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

/* 标题样式 - 加粗加大 */
.post-title {
  font-size: 18px;
  font-weight: 700;
  color: #333;
  margin: 0 0 12px 0;
  line-height: 1.4;
  padding-bottom: 12px;
  border-bottom: 1px solid #f0f0f0;
}

.card-content {
  margin-bottom: 12px;
}

/* 用户信息样式 */
.user-info {
  display: flex;
  align-items: center;
  margin-bottom: 12px;
}

.user-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  margin-right: 8px;
  object-fit: cover;
}

.user-name {
  font-weight: 500;
  color: #666;
  font-size: 14px;
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

/* 帖子正文 - 只显示一行（修复兼容性警告） */
.post-content {
  font-size: 14px;
  color: #666;
  line-height: 1.5;
  margin: 0 0 12px 0;
  
  /* 修复：同时定义标准属性和带前缀的属性 */
  display: -webkit-box;
  display: box;
  -webkit-line-clamp: 1;
  line-clamp: 1;
  -webkit-box-orient: vertical;
  box-orient: vertical;
  overflow: hidden;
  min-height: 21px; /* 保持一行高度 */
}

/* 图片展示样式 */
.post-images {
  display: flex;
  gap: 8px;
  margin-top: 8px;
}

.post-image {
  width: 120px;
  height: 90px;
  border-radius: 6px;
  object-fit: cover;
  flex-shrink: 0;
}
</style>