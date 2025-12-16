<template>
  <div class="card" @click="handleCardClick">
    <!-- 标题放在最上部 -->
    <h3 class="post-title">
      {{ post.title }}
    </h3>
    
    <div class="card-content">
      <!-- 用户信息：头像和名字 -->
      <div class="user-info">
        <!-- 后端返回的是 publisheravator 和 publishername (全小写) -->
        <img :src="post.publisheravator || defaultAvatar" alt="用户头像" class="user-avatar">
        <span class="user-name">{{ post.publishername || '匿名用户' }}</span>
      </div>

      <!-- 分类标签和宿舍楼标签 -->
      <div class="tags-container">
        <span class="post-category" :class="getCategoryClass(post.typeid)">
           {{ post.Type?.typename || getCategoryText(post.typeid) }}
        </span>
        <span class="dorm-tag">
          {{ post.Dorm?.dormname || '未知楼栋' }}
        </span>
      </div>
      
      <!-- 帖子正文 - 只显示一行 -->
      <p class="post-content">{{ post.content }}</p>
      
      <!-- 图片展示 - 最多两张 -->
      <div v-if="post.images && post.images.length > 0" class="post-images">
        <img 
          v-for="(image, index) in post.images.slice(0, 2)" 
          :key="index" 
          :src="image.image_url || image" 
          :alt="`图片${index + 1}`"
          class="post-image"
        >
      </div>
    </div>
    
    <!-- 使用独立的统计组件 -->
    <PostStats
      :view-count="post.view_count || 0"
      :comment-count="post.comment_count || 0"
      :like-count="post.like_count || 0"
      :time="formatTime(post.created_at || post.updated_at)"
    />
  </div>
</template>

<script setup>
import { defineProps } from 'vue'
import { useRouter } from 'vue-router'
import PostStats from '@/components/PostStats.vue'

const router = useRouter()
const defaultAvatar = 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'

const props = defineProps({
  post: {
    type: Object,
    default: () => ({})
  }
})

const getCategoryClass = (id) => {
  const map = { 1: 'trade', 2: 'help', 3: 'sports' }
  return map[id] || 'default'
}

const getCategoryText = (category) => {
  // 简单的 fallback，如果后端没返回 Type 对象
  const map = { 1: '闲置交易', 2: '跑腿求助', 3: '学习搭子', 4: '运动约球' }
  return map[category] || '其他'
}

const formatTime = (timeStr) => {
  if (!timeStr) return '刚刚'
  const date = new Date(timeStr)
  return date.toLocaleDateString()
}

const handleCardClick = () => {
  if (props.post.id) {
    router.push(`/post/detail/${props.post.id}`)
  }
}
</script>

<style scoped>
.card {
  background: white; border-radius: 12px; padding: 16px; margin-bottom: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1); cursor: pointer; transition: all 0.3s ease; border: 1px solid #e8e8e8;
}
.card:hover { transform: translateY(-2px); box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15); }
.post-title {
  font-size: 18px; font-weight: 700; color: #333; margin: 0 0 12px 0;
  line-height: 1.4; padding-bottom: 12px; border-bottom: 1px solid #f0f0f0;
}
.card-content { margin-bottom: 12px; }
.user-info { display: flex; align-items: center; margin-bottom: 12px; }
.user-avatar { width: 32px; height: 32px; border-radius: 50%; margin-right: 8px; object-fit: cover; }
.user-name { font-weight: 500; color: #666; font-size: 14px; }

.tags-container { display: flex; gap: 8px; margin-bottom: 12px; flex-wrap: wrap; }
.post-category { padding: 4px 8px; border-radius: 6px; font-size: 12px; font-weight: 500; background: #f4f4f5; color: #909399; }
.post-category.trade { background: #f0fff0; color: #2ed573; border: 1px solid #2ed573; }
.post-category.help { background: #fff8e1; color: #ffa502; border: 1px solid #ffa502; }
.post-category.sports { background: #f0f8ff; color: #1e90ff; border: 1px solid #1e90ff; }

.dorm-tag { padding: 4px 8px; background: #f8f9fa; color: #666; border: 1px solid #e9ecef; border-radius: 6px; font-size: 12px; font-weight: 500; }

.post-content {
  font-size: 14px; color: #666; line-height: 1.5; margin: 0 0 12px 0;
  
  display: -webkit-box;
  -webkit-line-clamp: 1;
  line-clamp: 1; 
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.post-images { display: flex; gap: 8px; margin-top: 8px; }
.post-image { width: 120px; height: 90px; border-radius: 6px; object-fit: cover; flex-shrink: 0; }
</style>