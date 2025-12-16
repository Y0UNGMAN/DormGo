<template>
  <div class="card" @click="handleCardClick">
    <h3 class="post-title">
      {{ post.title }}
    </h3>
    
    <div class="card-content">
      <div class="user-info">
        <img 
          :src="post.publisheravator || post.User?.avatarurl || defaultAvatar" 
          alt="用户头像" 
          class="user-avatar"
        >
        <span class="user-name">{{ post.publishername || post.User?.username || '未知用户' }}</span>
      </div>

      <div class="tags-container">
        <span class="post-category" :class="getCategoryClass(post.typeid)">
          {{ post.Type?.typename || post.PostType?.typename || '未分类' }}
        </span>
        
        <span class="dorm-tag">
          {{ post.Dorm?.dormname || '未知宿舍' }}
        </span>

        <span v-if="post.is_limited" class="limit-tag" :class="isExpired ? 'expired' : 'active'">
          {{ isExpired ? '已结束' : '正在报名' }}
        </span>
      </div>
      
      
      <div v-if="post.is_limited" class="deadline-info" :class="{ 'text-red': isExpired, 'text-green': !isExpired }">
        ⏰ 截止时间: {{ formatDetailTime(post.deadline) }}
      </div>

      <p class="post-content">{{ post.content }}</p>
      
      <div v-if="processedImages.length > 0" class="post-images">
        <img 
          v-for="(imgUrl, index) in processedImages.slice(0, 2)" 
          :key="index" 
          :src="imgUrl" 
          :alt="`图片${index + 1}`"
          class="post-image"
        >
      </div>
    </div>
    
    <PostStats
      :view-count="post.view_count || post.viewCount || 0"
      :comment-count="post.comment_count || post.commentCount || 0"
      :like-count="post.like_count || post.likeCount || 0"
      :time="formatTime(post.created_at || post.time)"
    />
  </div>
</template>

<script setup>
import { defineProps, computed } from 'vue'
import { useRouter } from 'vue-router'
import PostStats from '@/components/PostStats.vue'

const router = useRouter()
const defaultAvatar = 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png' // 默认头像

const props = defineProps({
  post: {
    type: Object,
    required: true,
    default: () => ({})
  }
})

// 计算是否过期
const isExpired = computed(() => {
  if (!props.post.deadline) return false
  return new Date() > new Date(props.post.deadline)
})

// 详细时间格式化 (用于截止时间)
const formatDetailTime = (timeStr) => {
  if (!timeStr) return ''
  const date = new Date(timeStr)
  return date.toLocaleString('zh-CN', {
    month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit'
  })
}


// 修改重点：处理图片逻辑
const processedImages = computed(() => {
  const imgs = props.post.images
  if (!imgs) return []

  // 情况1：如果是数组 (后端返回的新格式)
  if (Array.isArray(imgs)) {
    if (imgs.length === 0) return []

    // 判断数组里面是字符串还是对象
    // 如果是对象（你提供的数据格式），需要提取 image_url
    if (typeof imgs[0] === 'object' && imgs[0] !== null) {
      // 1. 拷贝数组防止影响原数据
      // 2. 按 order 字段排序 (确保图片顺序正确)
      // 3. 提取 image_url
      return [...imgs]
        .sort((a, b) => a.order - b.order)
        .map(item => item.image_url)
    }

    // 如果本身就是字符串数组 ['url1', 'url2']，直接返回
    return imgs
  }

  // 情况2：如果是字符串 "url1,url2" (兼容旧格式)
  if (typeof imgs === 'string') {
    return imgs.split(',')
  }

  return []
})

// 简单的分类样式映射
const getCategoryClass = (typeid) => {
  const classMap = {
    1: 'food',   
    2: 'sports', 
    3: 'help',   
    4: 'trade',  
    5: 'study'   
  }
  return classMap[typeid] || 'default-tag'
}

// 时间格式化
const formatTime = (timeStr) => {
  if (!timeStr) return ''
  const date = new Date(timeStr)
  return date.toLocaleDateString() 
}

const handleCardClick = () => {
  // 兼容 id 或 ID
  const postId = props.post.id || props.post.ID
  console.log('查看帖子详情:', postId)
  router.push(`/post/${postId}`)
}
</script>

<style scoped>
.card {
  background: white; border-radius: 12px; padding: 16px; margin-bottom: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1); cursor: pointer; transition: all 0.3s ease; border: 1px solid #e8e8e8;
}

.card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
}

.post-title {
  font-size: 18px;
  font-weight: 700;
  color: #333;
  margin: 0 0 12px 0;
  padding-bottom: 12px;
  border-bottom: 1px solid #f0f0f0;
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
  border: 1px solid #eee;
}

.user-name {
  font-weight: 500;
  color: #666;
  font-size: 14px;
}

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
  border: 1px solid #eee;
}

.post-category.food { background: #fff0f0; color: #ff4757; border-color: #ff4757; }
.post-category.sports { background: #f0f8ff; color: #1e90ff; border-color: #1e90ff; }
.post-category.help { background: #fff8e1; color: #ffa502; border-color: #ffa502; }
.post-category.trade { background: #f0fff0; color: #2ed573; border-color: #2ed573; }
.post-category.study { background: #f0f0ff; color: #5352ed; border-color: #5352ed; }
.post-category.default-tag { background: #f5f5f5; color: #666; border-color: #ddd; }

.dorm-tag {
  padding: 4px 8px;
  background: #f8f9fa;
  color: #666;
  border: 1px solid #e9ecef;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 500;
}

.post-content {
  font-size: 14px;
  color: #666;
  line-height: 1.5;
  margin: 0 0 12px 0;
  display: -webkit-box;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
  overflow: hidden;
  min-height: 21px;
}

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
  background-color: #f0f0f0;
}

.limit-tag {
  padding: 4px 8px;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 600;
  color: white;
}

/* 正在报名 - 绿色 */
.limit-tag.active {
  background-color: #52c41a; 
  border: 1px solid #52c41a;
}

/* 已结束 - 红色 */
.limit-tag.expired {
  background-color: #ff4d4f;
  border: 1px solid #ff4d4f;
}

.deadline-info {
  font-size: 13px;
  margin-bottom: 8px;
  font-weight: 500;
}

.text-green {
  color: #52c41a;
}

.text-red {
  color: #ff4d4f;
}


</style>