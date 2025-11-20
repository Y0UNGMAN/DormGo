<!-- src/components/PostStats.vue -->
<template>
  <div class="post-stats">
    <!-- 左侧：统计数据 -->
    <div class="stats-left">
      <span class="stat-item">
        <span class="stat-icon">👀</span>
        <span class="stat-number">{{ viewCount }}</span>
        <span class="stat-text">浏览</span>
      </span>
      <span class="stat-item">
        <span class="stat-icon">💬</span>
        <span class="stat-number">{{ commentCount }}</span>
        <span class="stat-text">评论</span>
      </span>
      <span class="stat-item">
        <span class="stat-icon">❤️</span>
        <span class="stat-number">{{ likeCount }}</span>
        <span class="stat-text">点赞</span>
      </span>
    </div>
    
    <!-- 右侧：发布时间 -->
    <div class="stats-right">
      <span class="post-time">{{ time }}</span>
    </div>
  </div>
</template>

<script setup>
import { defineProps } from 'vue'

const props = defineProps({
  viewCount: {
    type: Number,
    default: 0
  },
  commentCount: {
    type: Number,
    default: 0
  },
  likeCount: {
    type: Number,
    default: 0
  },
  time: {
    type: String,
    default: '刚刚'
  }
})
</script>

<style scoped>
.post-stats {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 0 0 0;
  border-top: 1px solid #f0f0f0;
  margin-top: 12px;
}

.stats-left {
  display: flex;
  gap: 20px;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: #999;
  cursor: default;
}

.stat-icon {
  font-size: 14px;
}

.stat-number {
  font-weight: 500;
  color: #666;
}

.stat-text {
  color: #999;
}

.stats-right {
  display: flex;
  align-items: center;
}

.post-time {
  font-size: 12px;
  color: #999;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .stats-left {
    gap: 16px;
  }
  
  .stat-item {
    font-size: 11px;
  }
}
</style>