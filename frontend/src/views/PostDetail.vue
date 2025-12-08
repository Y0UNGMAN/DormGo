<!-- src/views/PostDetail.vue -->
<template>
  <div class="post-detail">
    <!-- 返回按钮 -->
    <div class="back-header">
      <button class="back-btn" @click="goBack">
        <span class="back-icon">←</span>
        返回
      </button>
    </div>

    <!-- 帖子内容 -->
    <div class="post-content">
      <!-- 标题区域 -->
      <div class="title-section">
        <h1 class="post-title">{{ post.title }}</h1>
        <div class="post-meta">
          <div class="user-info">
            <img :src="post.publisheravator" alt="用户头像" class="user-avatar">
            <span class="user-name">{{ post.publishername }}</span>
          </div>
        </div>
      </div>

      <!-- 标签区域 -->
      <div class="tags-section">
        <span class="post-category" :class="post.typeid">
          {{ post.Type?.typename }}
        </span>
        <span class="dorm-tag">
          {{ post.Dorm?.dormname }}
        </span>
      </div>

      <!-- 正文内容 -->
      <div class="content-section">
        <p class="post-content-text">{{ post.content }}</p>
        
        <!-- 图片展示 -->
        <div v-if="post.images && post.images.length > 0" class="post-images">
          <img 
            v-for="(image, index) in post.images" 
            :key="index" 
            :src="image.image_url" 
            :alt="'图片' + (index + 1)"
            class="detail-image"
            @click="previewImage(index)"
          >
        </div>
      </div>

      <!-- 使用独立的统计组件 -->
      <PostStats
        :view-count="post.view_count"
        :comment-count="post.comment_count"
        :like-count="post.like_count"
        :time="post.updated_at"
      />
    </div>

    <!-- 操作按钮 -->
    <div class="action-buttons">
      <button class="action-btn like-btn" :class="{ liked: isLiked }" @click="toggleLike">
        <span class="btn-icon">{{ isLiked ? '❤️' : '🤍' }}</span>
        <span class="btn-text">{{ isLiked ? '已点赞' : '点赞' }}</span>
      </button>
      <button class="action-btn contact-btn" @click="contactUser">
        <span class="btn-icon">💬</span>
        <span class="btn-text">联系TA</span>
      </button>
    </div>

    <!-- 评论区 -->
    <div class="comments-section">
      <h3 class="comments-title">评论 ({{ comments.length }})</h3>
      
      <!-- 评论输入 -->
      <div class="comment-input-section">
        <img :src="currentUser.avatar" alt="用户头像" class="current-user-avatar">
        <div class="comment-input-container">
          <textarea 
            v-model="newComment" 
            placeholder="写下你的评论..." 
            class="comment-input"
            rows="3"
          ></textarea>
          <button class="submit-comment-btn" @click="submitComment" :disabled="!newComment.trim()">
            发布评论
          </button>
        </div>
      </div>

      <!-- 评论列表 -->
      <div class="comments-list">
        <div v-for="comment in comments" :key="comment.id" class="comment-item">
          <img :src="comment.commenter.avatarurl" alt="用户头像" class="comment-avatar">
          <div class="comment-content">
            <div class="comment-header">
              <span class="comment-user">{{ comment.commenter.username }}</span>
              <span class="comment-time">{{ formatTime(comment.created_at) }}</span>
            </div>
            <p class="comment-text">{{ comment.content }}</p>
            <div class="comment-actions">
              <!-- <button class="comment-action-btn" @click="likeComment(comment.id)">
                <span class="action-icon">❤️</span>
                <span class="action-text">{{ comment.likeCount || 0 }}</span>
              </button> -->
              <button class="comment-action-btn">
                <span class="action-icon">❤️</span>
                <span class="action-text">赞</span>
              </button>
              <button class="comment-action-btn" @click="replyComment(comment.id)">
                <span class="action-icon">↩️</span>
                <span class="action-text">回复</span>
              </button>
            </div>
            <div v-if="comment.sub_comments && comment.sub_comments.length > 0" class="sub-comments-list">
              <div v-for="sub in comment.sub_comments" :key="sub.id" class="sub-comment-item">
                <div class="sub-comment-header">
                  <span class="sub-user">{{ sub.commenter.username }}</span>
                  <span class="sub-time">{{ formatTime(sub.created_at) }}</span>
                </div>
                <p class="sub-text">
                  <span v-if="sub.reply_to_user?.username" class="reply-target">
                    回复 @{{ sub.reply_to_user.username }}: 
                  </span>
                  {{ sub.content }}
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 图片预览模态框 -->
    <div v-if="showImagePreview" class="image-preview-modal" @click="closeImagePreview">
      <div class="modal-content">
        <button class="close-modal-btn" @click="closeImagePreview">×</button>
        <img :src="post.images[currentImageIndex]" alt="预览图片" class="preview-image">
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import PostStats from '@/components/PostStats.vue'
import axios from 'axios'

const route = useRoute()
const router = useRouter()

// 帖子数据
const post = ref({})
const isLiked = ref(false)
const newComment = ref('')
const showImagePreview = ref(false)
const currentImageIndex = ref(0)

// 当前用户信息
const currentUser = ref({
  id: '1',
  name: '当前用户',
  avatar: 'https://dorm-go.oss-cn-guangzhou.aliyuncs.com/avator/midnight.jpg'
})

// 评论数据
const comments = ref([])
const totalComments = ref(0)

// 获取评论列表
const fetchComments = async () => {
  const postId = route.params.id
  console.log('获取评论，帖子ID:', postId)
  try {
    // 对应后端接口：GET /getcomment?post_id=1&page=1
    const response = await axios.get(`http://127.0.0.1:8080/api/v1/post/getcomment`, {
      params: {
        post_id: postId,
        page: 1,
      }
    })

    // 解析后端返回的结构: { code: 200, data: { list: [...], total: 5 } }
    if (response.data && response.data.code === 200) {
      const data = response.data.data
      comments.value = data.list || [] // 赋值给 comments
      totalComments.value = data.total || 0 // 赋值给总数
      console.log('评论获取成功:', comments.value)
    }
  } catch (err) {
    console.error('获取评论失败', err)
  }
}



//根据id获取帖子详情
const fetchPost = async () =>{
  const postId = route.params.id
  try {
    const response = await axios.get(`http://127.0.0.1:8080/api/v1/post/view/${postId}`);
    if(response.data && response.data.data){
      post.value = response.data.data
      console.log('帖子数据:', post.value);
    }
  }catch (err){
      console.error('获取帖子详情失败', err)
  }
}

// 返回上一页
const goBack = () => {
  router.back()
}

// 点赞/取消点赞
const toggleLike = () => {
  isLiked.value = !isLiked.value
  if (isLiked.value) {
    post.value.like_count++
  } else {
    post.value.like_count--
  }
}

// 联系用户
const contactUser = () => {
  console.log('联系用户:', post.value.userName)
  // 这里可以打开聊天窗口
}

// 提交评论
// 修改 submitComment 函数
const submitComment = async () => {
  // 1. 基础校验：内容不能为空
  if (!newComment.value.trim()) {
    alert("请输入评论内容")
    return
  }

  // 2. 获取必要参数
  const postId = parseInt(route.params.id) // 从路由获取当前帖子ID
  
  // ⚠️ 关键点：后端需要 commenter_id。
  // 在实际应用中，你应该从 localStorage 或 Pinia Store 中获取当前登录用户的真实ID。
  // 这里暂时使用 currentUser.id，请确保 currentUser.id 是数字类型或能转为数字。
  const userId = parseInt(currentUser.value.id) 

  if (!userId) {
    alert("无法获取用户信息，请重新登录")
    // router.push('/login') // 可以选择跳转去登录
    return
  }

  try {
    // 3. 发送 POST 请求
    const response = await axios.post('http://127.0.0.1:8080/api/v1/post/comment', {
      post_id: postId,
      content: newComment.value.trim(),
      commenter_id: userId,
      parent_id: 0 // 0 代表这是顶级评论（楼主层）
    })

    // 4. 处理响应
    if (response.data && response.data.code === 200) {
      // 成功：清空输入框
      newComment.value = ''
      
      // 刷新评论列表（重新从后端拉取最新数据）
      await fetchComments()
      
      // 更新帖子统计数据中的评论数（视觉上的+1）
      if (post.value) {
        post.value.comment_count++
      }
      
      alert("评论发布成功")
    } else {
      // 业务逻辑错误
      alert(response.data.msg || "发布失败")
    }
  } catch (error) {
    console.error("发布评论接口报错:", error)
    alert("网络错误，请稍后重试")
  }
}

// 点赞评论
const likeComment = (commentId) => {
  const comment = comments.value.find(c => c.id === commentId)
  if (comment) {
    comment.likeCount = (comment.likeCount || 0) + 1
  }
}

// 回复评论
const replyComment = (commentId) => {
  const comment = comments.value.find(c => c.id === commentId)
  if (comment) {
    newComment.value = `@${comment.userName} `
  }
}

// 预览图片
const previewImage = (index) => {
  currentImageIndex.value = index
  showImagePreview.value = true
}

// 关闭图片预览
const closeImagePreview = () => {
  showImagePreview.value = false
}

// 时间格式化工具函数
const formatTime = (timeStr) => {
  if (!timeStr) return ''
  const date = new Date(timeStr)
  // 转换为本地时间格式：2025/11/28 12:04
  return date.toLocaleString('zh-CN', {
    year: 'numeric', 
    month: '2-digit', 
    day: '2-digit', 
    hour: '2-digit', 
    minute: '2-digit'
  })
}

// 初始化
onMounted(() => {
  fetchPost()
  // 每次进入详情页，浏览量+1
  post.value.view_count++
  fetchComments()
})
</script>

<style scoped>
.post-detail {
  min-height: 100vh;
  background: #f5f5f5;
  padding: 0;
}

/* 返回按钮 */
.back-header {
  background: white;
  padding: 16px 20px;
  border-bottom: 1px solid #e8e8e8;
  position: sticky;
  top: 0;
  z-index: 100;
}

.back-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  background: none;
  border: none;
  font-size: 16px;
  color: #666;
  cursor: pointer;
  padding: 8px 0;
  transition: color 0.3s;
}

.back-btn:hover {
  color: #1890ff;
}

.back-icon {
  font-size: 18px;
}

/* 帖子内容 - 移除所有居中 */
.post-content {
  background: white;
  margin: 0;
  padding: 24px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  text-align: left; /* 确保文本左对齐 */
}

/* 标题区域 */
.title-section {
  margin-bottom: 20px;
  text-align: left; /* 标题左对齐 */
}

.post-title {
  font-size: 24px;
  font-weight: 700;
  color: #333;
  margin: 0 0 16px 0;
  line-height: 1.4;
  text-align: left; /* 标题文本左对齐 */
}

.post-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  text-align: left; /* 用户信息左对齐 */
}

.user-info {
  display: flex;
  align-items: center;
  gap: 12px;
  text-align: left; /* 用户信息左对齐 */
}

.user-avatar {
  width: 44px;
  height: 44px;
  border-radius: 50%;
  object-fit: cover;
}

.user-name {
  font-size: 16px;
  font-weight: 600;
  color: #333;
}

/* 标签区域 */
.tags-section {
  display: flex;
  gap: 12px;
  margin-bottom: 24px;
  flex-wrap: wrap;
  text-align: left; /* 标签左对齐 */
  justify-content: flex-start; /* 确保标签从左边开始 */
}

.post-category {
  padding: 6px 12px;
  border-radius: 6px;
  font-size: 14px;
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

.dorm-tag {
  padding: 6px 12px;
  background: #f8f9fa;
  color: #666;
  border: 1px solid #e9ecef;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
}

/* 正文内容 */
.content-section {
  margin-bottom: 24px;
  text-align: left; /* 正文内容左对齐 */
}

.post-content-text {
  font-size: 16px;
  color: #333;
  line-height: 1.8;
  margin: 0 0 24px 0;
  white-space: pre-wrap;
  text-align: left; /* 确保正文文本左对齐 */
}

/* 图片展示 */
.post-images {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 12px;
  margin-top: 16px;
  justify-items: start; /* 图片从左边开始排列 */
}

.detail-image {
  width: 100%;
  height: 200px;
  border-radius: 8px;
  object-fit: cover;
  cursor: pointer;
  transition: transform 0.3s ease;
}

.detail-image:hover {
  transform: scale(1.02);
}

/* 操作按钮 */
.action-buttons {
  display: flex;
  gap: 12px;
  padding: 20px;
  background: white;
  margin-top: 20px;
  text-align: left; /* 操作按钮区域左对齐 */
}

.action-btn {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 12px 20px;
  border: none;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
}

.like-btn {
  background: #f5f5f5;
  color: #666;
  border: 1px solid #e8e8e8;
}

.like-btn:hover {
  background: #fff0f0;
  color: #ff4757;
  border-color: #ff4757;
}

.like-btn.liked {
  background: #fff0f0;
  color: #ff4757;
  border-color: #ff4757;
}

.contact-btn {
  background: #1890ff;
  color: white;
}

.contact-btn:hover {
  background: #40a9ff;
}

/* 评论区 - 移除居中 */
.comments-section {
  background: white;
  margin: 20px 0 0 0;
  padding: 24px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  text-align: left; /* 评论区左对齐 */
}

.comments-title {
  font-size: 20px;
  font-weight: 600;
  color: #333;
  margin: 0 0 20px 0;
  text-align: left; /* 评论标题左对齐 */
}

/* 评论输入 */
.comment-input-section {
  display: flex;
  gap: 16px;
  margin-bottom: 24px;
  text-align: left; /* 评论输入左对齐 */
}

.current-user-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  object-fit: cover;
  flex-shrink: 0;
}

.comment-input-container {
  flex: 1;
  text-align: left; /* 评论输入容器左对齐 */
}

.comment-input {
  width: 100%;
  padding: 12px;
  border: 1px solid #e8e8e8;
  border-radius: 8px;
  font-size: 14px;
  resize: vertical;
  transition: border-color 0.3s;
  text-align: left; /* 输入框文本左对齐 */
}

.comment-input:focus {
  outline: none;
  border-color: #1890ff;
}

.submit-comment-btn {
  margin-top: 12px;
  padding: 8px 16px;
  background: #1890ff;
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 14px;
  cursor: pointer;
  transition: background 0.3s;
}

.submit-comment-btn:hover:not(:disabled) {
  background: #40a9ff;
}

.submit-comment-btn:disabled {
  background: #ccc;
  cursor: not-allowed;
}

/* 评论列表 */
.comments-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
  text-align: left; /* 评论列表左对齐 */
}

.comment-item {
  display: flex;
  gap: 12px;
  padding-bottom: 20px;
  border-bottom: 1px solid #f0f0f0;
  text-align: left; /* 每个评论项左对齐 */
}

.comment-item:last-child {
  border-bottom: none;
  padding-bottom: 0;
}

.comment-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  object-fit: cover;
  flex-shrink: 0;
}

.comment-content {
  flex: 1;
  text-align: left; /* 评论内容左对齐 */
}

.comment-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
  text-align: left; /* 评论头部左对齐 */
}

.comment-user {
  font-size: 14px;
  font-weight: 600;
  color: #333;
}

.comment-time {
  font-size: 12px;
  color: #999;
}

.comment-text {
  font-size: 14px;
  color: #333;
  line-height: 1.6;
  margin: 0 0 12px 0;
  text-align: left; /* 评论文本左对齐 */
}

.comment-actions {
  display: flex;
  gap: 16px;
  text-align: left; /* 评论操作左对齐 */
}

.comment-action-btn {
  display: flex;
  align-items: center;
  gap: 4px;
  background: none;
  border: none;
  font-size: 12px;
  color: #999;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 4px;
  transition: all 0.3s;
}

.comment-action-btn:hover {
  background: #f5f5f5;
  color: #666;
}

/* 图片预览模态框 */
.image-preview-modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.8);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  position: relative;
  max-width: 90%;
  max-height: 90%;
}

.close-modal-btn {
  position: absolute;
  top: -40px;
  right: 0;
  background: none;
  border: none;
  color: white;
  font-size: 24px;
  cursor: pointer;
  padding: 8px;
}

.preview-image {
  max-width: 100%;
  max-height: 80vh;
  object-fit: contain;
  border-radius: 8px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .post-content,
  .comments-section {
    padding: 16px;
  }
  
  .post-title {
    font-size: 20px;
  }
  
  .post-meta {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }
  
  .action-buttons {
    padding: 12px;
    flex-direction: column;
  }
  
  .comment-input-section {
    flex-direction: column;
  }
  
  .current-user-avatar {
    align-self: flex-start;
  }
}

/* 新增：子评论样式 */
.sub-comments-list {
  background: #f9f9f9; /* 浅灰色背景，区分主评论 */
  padding: 12px;
  border-radius: 8px;
  margin-top: 12px;
}

.sub-comment-item {
  margin-bottom: 10px;
  border-bottom: 1px dashed #eee; /* 虚线分隔 */
  padding-bottom: 10px;
  text-align: left; /* 强制左对齐 */
}

.sub-comment-item:last-child {
  margin-bottom: 0;
  border-bottom: none;
  padding-bottom: 0;
}

.sub-comment-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 4px;
}

.sub-user {
  font-size: 13px;
  font-weight: 600;
  color: #666;
}

.sub-time {
  font-size: 12px;
  color: #bbb;
}

.sub-text {
  font-size: 13px;
  color: #444;
  margin: 0;
  text-align: left;
}

.reply-target {
  color: #1890ff; /* 蓝色的 "回复 @xxx" */
  font-weight: 500;
  margin-right: 4px;
}
</style>