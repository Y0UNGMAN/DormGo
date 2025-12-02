<template>
  <div class="post-detail">
    <div class="back-header">
      <button class="back-btn" @click="goBack">
        <span class="back-icon">←</span>
        返回
      </button>
    </div>

    <div class="post-content">
      <div class="title-section">
        <h1 class="post-title">{{ post.title }}</h1>
        <div class="post-meta">
          <div class="user-info">
            <img :src="post.publisheravator" alt="用户头像" class="user-avatar">
            <span class="user-name">{{ post.publishername }}</span>
          </div>
        </div>
      </div>

      <div class="tags-section">
        <span class="post-category" :class="post.typeid">
          {{ post.Type?.typename }}
        </span>
        <span class="dorm-tag">
          {{ post.Dorm?.dormname }}
        </span>
      </div>

      <div class="content-section">
        <p class="post-content-text">{{ post.content }}</p>

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

      <PostStats
        :view-count="post.view_count"
        :comment-count="post.comment_count"
        :like-count="post.like_count"
        :time="formatTime(post.updated_at)"
      />
    </div>

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

    <div class="comments-section">
      <h3 class="comments-title">评论 ({{ post.comment_count }})</h3>
      
      <div class="comment-input-section">
        <img :src="currentUser.avatarurl" alt="用户头像" class="current-user-avatar">
        <div class="comment-input-container">
          <div v-if="replyingToComment" class="reply-target-tip">
            正在回复：@{{ replyingToComment.commenter.username }}
            <button class="cancel-reply-btn" @click="cancelReply">×</button>
          </div>

          <textarea 
            v-model="newComment" 
            :placeholder="replyingToComment ? '回复 ' + replyingToComment.commenter.username + '...' : '写下你的评论...'" 
            class="comment-input"
            rows="3"
          ></textarea>
          <button 
            class="submit-comment-btn" 
            @click="submitComment" 
            :disabled="!newComment.trim() || isSubmitting"
          >
            {{ isSubmitting ? '发布中...' : (replyingToComment ? '回复' : '发布评论') }}
          </button>
        </div>
      </div>

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
              <button class="comment-action-btn" @click="replyComment(comment)">
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

    <div v-if="showImagePreview && post.images && post.images.length > 0" class="image-preview-modal" @click="closeImagePreview">
      <div class="modal-content" @click.stop>
        <button class="close-modal-btn" @click="closeImagePreview">×</button>
        <img :src="post.images[currentImageIndex].image_url" alt="预览图片" class="preview-image">
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, reactive } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import PostStats from '@/components/PostStats.vue'

const route = useRoute()
const router = useRouter()

// 帖子数据
const post = ref({ comment_count: 0 }) 
const isLiked = ref(false)
const newComment = ref('')
const showImagePreview = ref(false)
const currentImageIndex = ref(0)
const comments = ref([])
const isSubmitting = ref(false)

// 回复状态管理
const replyToCommentId = ref(null)
const replyToTargetUser = ref(null)
const expandedComments = reactive(new Set())

// 当前用户信息（模拟）
const currentUser = ref({
  id: '1',
  username: '当前用户',
  avatarurl: 'https://dorm-go.oss-cn-guangzhou.aliyuncs.com/avator/current_user.jpg' 
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

// ********* 评论操作函数 *********


// ⭐ 提交评论 - 使用 alert() 提示
const submitComment = () => {
  if (!newComment.value.trim() || isSubmitting.value) return

  isSubmitting.value = true
  
  // 模拟网络请求和后端处理
  setTimeout(() => {
    isSubmitting.value = false;
    
    // 获取回复目标信息
    const postId = route.params.id
    const content = newComment.value.trim()
    const parentId = replyToCommentId.value
    const replyToUserId = replyToTargetUser.value ? replyToTargetUser.value.id : null

    // 模拟后端返回的新评论对象
    const newCommentData = {
        id: Date.now().toString(), 
        post_id: postId,
        parent_id: parentId,
        content: content.replace(`@${replyToTargetUser.value?.username} `, '').trim(),
        created_at: new Date().toISOString(),
        commenter: currentUser.value,
        reply_to_user: replyToTargetUser.value,
        sub_comments: []
    }
    
    // 更新前端数据
    if (parentId) {
      const parent = comments.value.find(c => c.id === parentId)
      if (parent) {
        if (!parent.sub_comments) {
          parent.sub_comments = []
        }
        parent.sub_comments.push(newCommentData)
        expandedComments.add(parentId)
      }
    } else {
      comments.value.unshift(newCommentData)
    }

    // 更新评论总数
    post.value.comment_count = post.value.comment_count + 1

    // 清理状态和提示
    alert('评论发布成功！');
    
    cancelReply() 

  }, 1000); // 模拟 1 秒延迟
}

// ********* 子评论展开/折叠函数 *********

// 切换子评论的展开/折叠状态
const toggleSubComments = (commentId) => {
  if (expandedComments.has(commentId)) {
    expandedComments.delete(commentId)
  } else {
    expandedComments.add(commentId)
  }
}

// 检查评论是否处于展开状态
const isCommentExpanded = (commentId) => {
  return expandedComments.has(commentId)
}

// 获取应该显示的子评论列表
const getVisibleSubComments = (comment) => {
  if (!comment.sub_comments || comment.sub_comments.length === 0) {
    return []
  }
  const defaultVisibleCount = 1;

  if (isCommentExpanded(comment.id)) {
    return comment.sub_comments
  }
  return comment.sub_comments.slice(0, defaultVisibleCount)
}


// ********* 帖子数据获取与工具函数 *********

const fetchComments = () => {
  setTimeout(() => {
    comments.value = mockComments;
  }, 500);
}

const fetchPost = () =>{
  setTimeout(() => {
    post.value = mockPost;
    post.value.view_count = (post.value.view_count || 0) + 1;
  }, 300);
}

const goBack = () => {
  router.back()
}

const toggleLike = () => {
  isLiked.value = !isLiked.value
  if (isLiked.value) {
    post.value.like_count++
  } else {
    post.value.like_count--
  }
}

const contactUser = () => {
  alert(`尝试联系用户: ${post.value.publishername}`);
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
/* 样式代码与上一轮提供的完全一致，为简洁在此省略，请确保在您的项目中包含完整的样式 */
.post-detail {
  min-height: 100vh;
  background: #f5f5f5;
  padding: 0;
}
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
.post-content {
  background: white;
  margin: 0;
  padding: 24px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  text-align: left; 
}
.title-section {
  margin-bottom: 20px;
  text-align: left; 
}
.post-title {
  font-size: 24px;
  font-weight: 700;
  color: #333;
  margin: 0 0 16px 0;
  line-height: 1.4;
  text-align: left; 
}
.post-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  text-align: left; 
}
.user-info {
  display: flex;
  align-items: center;
  gap: 12px;
  text-align: left; 
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
.tags-section {
  display: flex;
  gap: 12px;
  margin-bottom: 24px;
  flex-wrap: wrap;
  text-align: left; 
  justify-content: flex-start; 
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
.content-section {
  margin-bottom: 24px;
  text-align: left; 
}
.post-content-text {
  font-size: 16px;
  color: #333;
  line-height: 1.8;
  margin: 0 0 24px 0;
  white-space: pre-wrap;
  text-align: left; 
}
.post-images {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 12px;
  margin-top: 16px;
  justify-items: start; 
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
.action-buttons {
  display: flex;
  gap: 12px;
  padding: 20px;
  background: white;
  margin-top: 20px;
  text-align: left; 
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
.comments-section {
  background: white;
  margin: 20px 0 0 0;
  padding: 24px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  text-align: left; 
}
.comments-title {
  font-size: 20px;
  font-weight: 600;
  color: #333;
  margin: 0 0 20px 0;
  text-align: left; 
}
.comment-input-section {
  display: flex;
  gap: 16px;
  margin-bottom: 24px;
  text-align: left; 
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
  text-align: left; 
}
.reply-target-tip {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 12px;
  background: #e6f7ff;
  border-radius: 6px 6px 0 0;
  border: 1px solid #91d5ff;
  border-bottom: none;
  font-size: 14px;
  color: #1890ff;
}
.cancel-reply-btn {
  background: none;
  border: none;
  color: #1890ff;
  font-size: 16px;
  cursor: pointer;
  opacity: 0.7;
  padding: 0 4px;
}
.cancel-reply-btn:hover {
  opacity: 1;
}
.comment-input {
  width: 100%;
  padding: 12px;
  border: 1px solid #e8e8e8;
  border-radius: 8px;
  font-size: 14px;
  resize: vertical;
  transition: border-color 0.3s;
  text-align: left; 
}
.comment-input-container .comment-input:focus {
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
.comments-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
  text-align: left; 
}
.comment-item {
  display: flex;
  gap: 12px;
  padding-bottom: 20px;
  border-bottom: 1px solid #f0f0f0;
  text-align: left; 
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
  text-align: left; 
}
.comment-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
  text-align: left; 
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
  text-align: left; 
}
.comment-actions {
  display: flex;
  gap: 16px;
  text-align: left; 
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
.sub-comments-list {
  background: #f9f9f9; 
  padding: 12px;
  border-radius: 8px;
  margin-top: 12px;
}
.sub-comment-item {
  margin-bottom: 10px;
  border-bottom: 1px dashed #eee; 
  padding-bottom: 10px;
  text-align: left; 
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
  color: #1890ff; 
  font-weight: 500;
  margin-right: 4px;
}
.sub-comments-footer {
  display: flex;
  justify-content: flex-start;
  gap: 12px;
  margin-top: 8px;
  padding-top: 8px;
}
.toggle-comments-btn {
  background: none;
  border: none;
  color: #1890ff;
  font-size: 12px;
  cursor: pointer;
  padding: 4px 0;
  transition: opacity 0.3s;
}
.toggle-comments-btn:hover {
  opacity: 0.8;
}
.reply-sub-btn {
    color: #666;
    border-left: 1px solid #e8e8e8;
    padding-left: 12px;
}
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