<template>
  <div class="post-detail">
    <!-- 返回按钮 -->
    <div class="back-header">
      <button class="back-btn" @click="goBack">
        <span class="back-icon">←</span> 返回
      </button>
    </div>

    <!-- 帖子内容 -->
    <div class="post-content">
      <div class="title-section">
        <h1 class="post-title">{{ post.title }}</h1>
        <div class="post-meta">
          <div class="user-info">
            <img :src="post.publisheravator" class="user-avatar">
            <span class="user-name">{{ post.publishername }}</span>
          </div>
        </div>
      </div>

      <div class="tags-section">
        <span class="post-category" :class="post.typeid">{{ post.Type?.typename }}</span>
        <span class="dorm-tag">{{ post.Dorm?.dormname }}</span>
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
      <h3 class="comments-title">评论 ({{ post.comment_count }})</h3>

      <!-- 评论输入框 -->
      <div class="comment-input-section">
        <img :src="currentUser.avatarurl" class="current-user-avatar">
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

      <!-- 评论列表 -->
      <div class="comments-list">
        <div 
          v-for="comment in comments" 
          :key="comment.id" 
          class="comment-item"
          @click.stop="replyComment(comment)"
        >
          <img :src="comment.commenter.avatarurl" class="comment-avatar">
          <div class="comment-content">
            <div class="comment-header">
              <span class="comment-user">{{ comment.commenter.username }}</span>
              <span class="comment-time">{{ formatTime(comment.created_at) }}</span>
            </div>
            <p class="comment-text">{{ comment.content }}</p>
            <div class="comment-actions">
              <button class="comment-action-btn" @click.stop="likeComment(comment.id)">
                <span class="action-icon">❤️</span>
                <span class="action-text">{{ comment.likeCount || 0 }}</span>
              </button>
              <button class="comment-action-btn" @click.stop="replyComment(comment)">
                <span class="action-icon">↩️</span>
                <span class="action-text">回复</span>
              </button>
            </div>

            <!-- 子评论列表 -->
            <div v-if="comment.sub_comments && comment.sub_comments.length > 0" class="sub-comments-list">
              <div
                v-for="sub in getVisibleSubComments(comment)"
                :key="sub.id"
                class="sub-comment-item"
                @click.stop="replyComment(sub)"
              >
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

              <!-- 展开/折叠更多子评论 -->
              <div v-if="comment.sub_comments.length > 1" class="sub-comments-footer">
                <button class="toggle-comments-btn" @click.stop="toggleSubComments(comment.id)">
                  {{ isCommentExpanded(comment.id) ? '收起' : '查看更多' }}
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 图片预览弹窗 -->
    <div v-if="showImagePreview && post.images && post.images.length > 0" class="image-preview-modal" @click="closeImagePreview">
      <div class="modal-content" @click.stop>
        <button class="close-modal-btn" @click="closeImagePreview">×</button>
        <img :src="post.images[currentImageIndex].image_url" class="preview-image">
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from 'axios'
import PostStats from '@/components/PostStats.vue'

const route = useRoute()
const router = useRouter()

// 帖子信息
const post = ref({ comment_count: 0, like_count: 0, view_count: 0 })
const isLiked = ref(false)
const newComment = ref('')
const showImagePreview = ref(false)
const currentImageIndex = ref(0)
const isSubmitting = ref(false)

// 回复状态
const replyingToComment = ref(null)
const replyToCommentId = ref(null)
const replyToTargetUser = ref(null)
const expandedComments = reactive(new Set())

// 当前用户模拟
const currentUser = ref({
  id: '1',
  username: '当前用户',
  avatarurl: 'https://dorm-go.oss-cn-guangzhou.aliyuncs.com/avator/current_user.jpg'
})

// 评论列表
const comments = ref([])
const totalComments = ref(0)

// 获取帖子数据（模拟/真实接口）
const fetchPost = async () => {
  // 模拟
  post.value = {
    id: route.params.id,
    title: '示例帖子标题',
    content: '这里是帖子内容...',
    publishername: '张三',
    publisheravator: 'https://dorm-go.oss-cn-guangzhou.aliyuncs.com/avator/current_user.jpg',
    Type: { typename: '分享' },
    Dorm: { dormname: '宿舍A' },
    view_count: 123,
    like_count: 10,
    comment_count: 2,
    updated_at: new Date().toISOString(),
    images: [
      { image_url: 'https://dorm-go.oss-cn-guangzhou.aliyuncs.com/img1.jpg' },
      { image_url: 'https://dorm-go.oss-cn-guangzhou.aliyuncs.com/img2.jpg' }
    ]
  }
  post.value.view_count++
}

// 获取评论列表（模拟/真实接口）
const fetchComments = async () => {
  comments.value = [
    {
      id: 'c1',
      content: '这是主评论1',
      commenter: { username: '用户A', avatarurl: 'https://dorm-go.oss-cn-guangzhou.aliyuncs.com/avator/1.jpg' },
      created_at: new Date().toISOString(),
      likeCount: 1,
      sub_comments: [
        {
          id: 's1',
          content: '这是子评论1',
          commenter: { username: '用户B' , avatarurl: 'https://dorm-go.oss-cn-guangzhou.aliyuncs.com/avator/2.jpg'},
          created_at: new Date().toISOString(),
          reply_to_user: { username: '用户A' }
        }
      ]
    },
    {
      id: 'c2',
      content: '这是主评论2',
      commenter: { username: '用户C', avatarurl: 'https://dorm-go.oss-cn-guangzhou.aliyuncs.com/avator/3.jpg' },
      created_at: new Date().toISOString(),
      likeCount: 0,
      sub_comments: []
    }
  ]
}

// 进入详情页
onMounted(() => {
  fetchPost()
  fetchComments()
})

// 返回
const goBack = () => router.back()

// 点赞帖子
const toggleLike = () => {
  isLiked.value = !isLiked.value
  if (isLiked.value) post.value.like_count++
  else post.value.like_count--
}

// 联系用户
const contactUser = () => alert(`尝试联系用户: ${post.value.publishername}`)

// 提交评论
const submitComment = async () => {
  if (!newComment.value.trim() || isSubmitting.value) return
  isSubmitting.value = true
  const content = newComment.value.trim()
  const parentId = replyToCommentId.value
  const replyToUserId = replyToTargetUser.value?.id || null

  // 模拟接口
  setTimeout(() => {
    const newCommentData = {
      id: Date.now().toString(),
      post_id: post.value.id,
      parent_id: parentId,
      content: content.replace(`@${replyToTargetUser.value?.username} `, '').trim(),
      created_at: new Date().toISOString(),
      commenter: currentUser.value,
      reply_to_user: replyToTargetUser.value,
      sub_comments: []
    }
    if (parentId) {
      const parent = comments.value.find(c => c.id === parentId)
      if (!parent.sub_comments) parent.sub_comments = []
      parent.sub_comments.push(newCommentData)
      expandedComments.add(parentId)
    } else {
      comments.value.unshift(newCommentData)
    }
    post.value.comment_count++
    isSubmitting.value = false
    cancelReply()
    alert('评论发布成功！')
  }, 800)
}

// 点赞评论
const likeComment = (commentId) => {
  const comment = comments.value.find(c => c.id === commentId)
  if (comment) comment.likeCount = (comment.likeCount || 0) + 1
}

// 点击评论或子评论回复
const replyComment = (comment) => {
  replyingToComment.value = comment
  replyToCommentId.value = comment.id
  replyToTargetUser.value = comment.commenter
  newComment.value = `@${comment.commenter.username} `
  nextTick(() => {
    const textarea = document.querySelector('.comment-input')
    if (textarea) textarea.focus()
  })
}

// 取消回复
const cancelReply = () => {
  replyingToComment.value = null
  replyToCommentId.value = null
  replyToTargetUser.value = null
  newComment.value = ''
}

// 子评论展开/折叠
const toggleSubComments = (commentId) => {
  if (expandedComments.has(commentId)) expandedComments.delete(commentId)
  else expandedComments.add(commentId)
}
const isCommentExpanded = (commentId) => expandedComments.has(commentId)
const getVisibleSubComments = (comment) => {
  if (!comment.sub_comments || comment.sub_comments.length === 0) return []
  const defaultVisibleCount = 1
  return isCommentExpanded(comment.id) ? comment.sub_comments : comment.sub_comments.slice(0, defaultVisibleCount)
}

// 图片预览
const previewImage = (index) => {
  currentImageIndex.value = index
  showImagePreview.value = true
}
const closeImagePreview = () => showImagePreview.value = false

// 时间格式化
const formatTime = (timeStr) => {
  if (!timeStr) return ''
  const date = new Date(timeStr)
  return date.toLocaleString('zh-CN', {
    year: 'numeric', month: '2-digit', day: '2-digit',
    hour: '2-digit', minute: '2-digit'
  })
}
</script>

<style scoped>
/* ========== 原样保留并略微格式化你的样式 ========== */
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
.back-btn:hover { color: #1890ff; }
.back-icon { font-size: 18px; }

.post-content {
  background: white;
  margin: 0;
  padding: 24px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
  text-align: left;
}
.title-section { margin-bottom: 20px; text-align: left; }
.post-title {
  font-size: 24px; font-weight: 700; color: #333;
  margin: 0 0 16px 0; line-height: 1.4;
  text-align: left;
}
.post-meta { display:flex; justify-content:space-between; align-items:center; text-align:left; }
.user-info { display:flex; align-items:center; gap:12px; text-align:left; }
.user-avatar { width:44px; height:44px; border-radius:50%; object-fit:cover; }
.user-name { font-size:16px; font-weight:600; color:#333; }
.tags-section {
  display:flex; gap:12px; margin-bottom:24px; flex-wrap:wrap; text-align:left; justify-content:flex-start;
}
.post-category { padding:4px 8px; border-radius:6px; font-size:12px; font-weight:500; }
.dorm-tag {
  padding:6px 12px; background:#f8f9fa; color:#666; border:1px solid #e9ecef; border-radius:6px;
  font-size:14px; font-weight:500;
}
.content-section { margin-bottom:24px; text-align:left; }
.post-content-text { font-size:16px; color:#333; line-height:1.8; margin:0 0 24px 0; white-space:pre-wrap; text-align:left; }

.post-images {
  display:grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap:12px; margin-top:16px; justify-items:start;
}
.detail-image {
  width:100%; height:200px; border-radius:8px; object-fit:cover; cursor:pointer; transition:transform 0.3s ease;
}
.detail-image:hover { transform:scale(1.02); }

.action-buttons {
  display:flex; gap:12px; padding:20px; background:white; margin-top:20px; text-align:left;
}
.action-btn {
  flex:1; display:flex; align-items:center; justify-content:center; gap:8px; padding:12px 20px; border:none;
  border-radius:8px; font-size:16px; font-weight:500; cursor:pointer; transition:all 0.3s ease;
}
.like-btn { background:#f5f5f5; color:#666; border:1px solid #e8e8e8; }
.like-btn:hover { background:#fff0f0; color:#ff4757; border-color:#ff4757; }
.like-btn.liked { background:#fff0f0; color:#ff4757; border-color:#ff4757; }
.contact-btn { background:#1890ff; color:white; }
.contact-btn:hover { background:#40a9ff; }

.comments-section {
  background:white; margin:20px 0 0 0; padding:24px; box-shadow:0 2px 8px rgba(0,0,0,0.1); text-align:left;
}
.comments-title { font-size:20px; font-weight:600; color:#333; margin:0 0 20px 0; text-align:left; }
.comment-input-section { display:flex; gap:16px; margin-bottom:24px; text-align:left; }
.current-user-avatar { width:40px; height:40px; border-radius:50%; object-fit:cover; flex-shrink:0; }
.comment-input-container { flex:1; text-align:left; }
.reply-target-tip {
  display:flex; align-items:center; justify-content:space-between; padding:8px 12px; background:#e6f7ff;
  border-radius:6px 6px 0 0; border:1px solid #91d5ff; border-bottom:none; font-size:14px; color:#1890ff;
}
.cancel-reply-btn { background:none; border:none; color:#1890ff; font-size:16px; cursor:pointer; opacity:0.7; padding:0 4px; }
.cancel-reply-btn:hover { opacity:1; }

.comment-input { width:100%; padding:12px; border:1px solid #e8e8e8; border-radius:8px; font-size:14px; resize:vertical; transition:border-color 0.3s; text-align:left; }
.comment-input-container .comment-input:focus { outline:none; border-color:#1890ff; }
.submit-comment-btn { margin-top:12px; padding:8px 16px; background:#1890ff; color:white; border:none; border-radius:6px; font-size:14px; cursor:pointer; transition:background 0.3s; }
.submit-comment-btn:hover:not(:disabled) { background:#40a9ff; }
.submit-comment-btn:disabled { background:#ccc; cursor:not-allowed; }

.comments-list { display:flex; flex-direction:column; gap:20px; text-align:left; }
.comment-item { display:flex; gap:12px; padding-bottom:20px; border-bottom:1px solid #f0f0f0; text-align:left; }
.comment-item:last-child { border-bottom:none; padding-bottom:0; }
.comment-avatar { width:36px; height:36px; border-radius:50%; object-fit:cover; flex-shrink:0; }
.comment-content { flex:1; text-align:left; }
.comment-header { display:flex; justify-content:space-between; align-items:center; margin-bottom:8px; text-align:left; }
.comment-user { font-size:14px; font-weight:600; color:#333; }
.comment-time { font-size:12px; color:#999; }
.comment-text { font-size:14px; color:#333; line-height:1.6; margin:0 0 12px 0; text-align:left; }

.comment-actions { display:flex; gap:16px; text-align:left; }
.comment-action-btn {
  display:flex; align-items:center; gap:4px; background:none; border:none; font-size:12px; color:#999;
  cursor:pointer; padding:4px 8px; border-radius:4px; transition:all 0.3s;
}
.comment-action-btn:hover { background:#f5f5f5; color:#666; }

.sub-comments-list { background:#f9f9f9; padding:12px; border-radius:8px; margin-top:12px; }
.sub-comment-item { margin-bottom:10px; border-bottom:1px dashed #eee; padding-bottom:10px; text-align:left; }
.sub-comment-item:last-child { margin-bottom:0; border-bottom:none; padding-bottom:0; }
.sub-comment-header { display:flex; justify-content:space-between; margin-bottom:4px; }
.sub-user { font-size:13px; font-weight:600; color:#666; }
.sub-time { font-size:12px; color:#bbb; }
.sub-text { font-size:13px; color:#444; margin:0; text-align:left; }
.reply-target { color:#1890ff; font-weight:500; margin-right:4px; }
.sub-comments-footer { display:flex; justify-content:flex-start; gap:12px; margin-top:8px; padding-top:8px; }
.toggle-comments-btn { background:none; border:none; color:#1890ff; font-size:12px; cursor:pointer; padding:4px 0; transition:opacity 0.3s; }
.toggle-comments-btn:hover { opacity:0.8; }
.reply-sub-btn { color:#666; border-left:1px solid #e8e8e8; padding-left:12px; }

.image-preview-modal {
  position:fixed; top:0; left:0; right:0; bottom:0; background:rgba(0,0,0,0.8);
  display:flex; align-items:center; justify-content:center; z-index:1000;
}
.modal-content { position:relative; max-width:90%; max-height:90%; }
.close-modal-btn {
  position:absolute; top:-40px; right:0; background:none; border:none; color:white; font-size:24px; cursor:pointer; padding:8px;
}
.preview-image { max-width:100%; max-height:80vh; object-fit:contain; border-radius:8px; }

@media (max-width:768px) {
  .post-content, .comments-section { padding:16px; }
  .post-title { font-size:20px; }
  .post-meta { flex-direction:column; align-items:flex-start; gap:8px; }
  .action-buttons { padding:12px; flex-direction:column; }
  .comment-input-section { flex-direction:column; }
  .current-user-avatar { align-self:flex-start; }
}
</style>
