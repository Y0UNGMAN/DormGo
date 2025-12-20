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

        <span v-if="post.is_pinned" class="tag-badge pinned">置顶</span>
        
        <span v-if="isAdminPost" class="tag-badge official">官方</span>
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
const defaultAvatar = 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png' 

// 定义管理员的 ID，方便后续修改（例如改成 100 或其他）
const ADMIN_USER_ID = 10; 

const props = defineProps({
  post: {
    type: Object,
    required: true,
    default: () => ({})
  }
})

// 计算是否是管理员发布的帖子
const isAdminPost = computed(() => {
  // 确保类型转换一致，防止 "1" !== 1 的问题
  return parseInt(props.post.publisherid) === ADMIN_USER_ID;
})

// 计算是否过期
const isExpired = computed(() => {
  if (!props.post.deadline) return false
  return new Date() > new Date(props.post.deadline)
})

// 详细时间格式化
const formatDetailTime = (timeStr) => {
  if (!timeStr) return ''
  const date = new Date(timeStr)
  return date.toLocaleString('zh-CN', {
    month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit'
  })
}

// 处理图片逻辑 (兼容字符串和数组对象)
const processedImages = computed(() => {
  const imgs = props.post.images
  if (!imgs) return []

  if (Array.isArray(imgs)) {
    if (imgs.length === 0) return []
    // 如果是对象数组 [{image_url: '...'}, ...]
    if (typeof imgs[0] === 'object' && imgs[0] !== null) {
      return [...imgs]
        .sort((a, b) => a.order - b.order)
        .map(item => item.image_url)
    }
    // 如果已经是字符串数组
    return imgs
  }

  // 如果是逗号分隔的字符串
  if (typeof imgs === 'string') {
    return imgs.split(',')
  }

  return []
})

// 分类样式映射
const getCategoryClass = (typeid) => {
  const classMap = {
    1: 'food',   
    2: 'sports', 
    3: 'help',   
    4: 'trade',  
    5: 'study',
    99: 'official' // 假设99是系统公告，给个特殊样式
  }
  return classMap[typeid] || 'default-tag'
}

// 简单日期格式化
const formatTime = (timeStr) => {
  if (!timeStr) return ''
  const date = new Date(timeStr)
  return date.toLocaleDateString() 
}

const handleCardClick = () => {
  const postId = props.post.id || props.post.ID
  router.push(`/post/${postId}`)
}
</script>

<style scoped>
.card {
  background: white; border-radius: 12px; padding: 16px; margin-bottom: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1); cursor: pointer; transition: all 0.3s ease; border: 1px solid #e8e8e8;
  text-align: left;
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

/* 用户信息栏 */
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

/* 标签通用样式 */
.tag-badge {
  font-size: 10px;
  padding: 2px 6px;
  border-radius: 4px;
  margin-left: 6px; /* 标签之间的间距 */
  font-weight: bold;
  color: white;
  line-height: 1.4;
  display: inline-block;
  vertical-align: middle;
}

/* 置顶标签：红色 */
.pinned {
  background-color: #ff4d4f; 
  box-shadow: 0 2px 4px rgba(255, 77, 79, 0.3);
}

/* 官方标签：蓝色 */
.official {
  background-color: #1890ff; 
  box-shadow: 0 2px 4px rgba(24, 144, 255, 0.3);
}

/* 标签容器（分类、宿舍等） */
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

/* 不同分类的颜色 */
.post-category.food { background: #fff0f0; color: #ff4757; border-color: #ff4757; }
.post-category.sports { background: #f0f8ff; color: #1e90ff; border-color: #1e90ff; }
.post-category.help { background: #fff8e1; color: #ffa502; border-color: #ffa502; }
.post-category.trade { background: #f0fff0; color: #2ed573; border-color: #2ed573; }
.post-category.study { background: #f0f0ff; color: #5352ed; border-color: #5352ed; }
.post-category.official { background: #e6f7ff; color: #1890ff; border-color: #1890ff; } /* 官方公告分类样式 */
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
  -webkit-line-clamp: 2;
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

.limit-tag.active { background-color: #52c41a; border: 1px solid #52c41a; }
.limit-tag.expired { background-color: #ff4d4f; border: 1px solid #ff4d4f; }

.deadline-info {
  font-size: 13px;
  margin-bottom: 8px;
  font-weight: 500;
}

.text-green { color: #52c41a; }
.text-red { color: #ff4d4f; }
</style>