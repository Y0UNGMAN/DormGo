<!-- src/views/PublishPost.vue -->
<template>
  <div class="publish-post">
    <!-- 返回按钮 -->
    <div class="back-header">
      <button class="back-btn" @click="goBack">
        <span class="back-icon">←</span>
        返回
      </button>
      <h1 class="page-title">发布帖子</h1>
    </div>

    <!-- 发布表单 -->
    <div class="publish-form">
      <!-- 标题输入 -->
      <div class="form-section">
        <label class="form-label">标题</label>
        <input
          v-model="form.title"
          type="text"
          placeholder="请输入帖子标题"
          class="title-input"
          maxlength="50"
        >
        <div class="char-count">{{ form.title.length }}/50</div>
      </div>

      <!-- 正文输入 -->
      <div class="form-section">
        <label class="form-label">正文</label>
        <textarea
          v-model="form.content"
          placeholder="请输入帖子内容..."
          class="content-textarea"
          rows="8"
          maxlength="1000"
        ></textarea>
        <div class="char-count">{{ form.content.length }}/1000</div>
      </div>

      <!-- 分类选择 -->
      <div class="form-section">
        <label class="form-label">分类</label>
        <div class="category-options">
          <button
            v-for="category in categories"
            :key="category.id"
            class="category-option"
            :class="{ active: form.category === category.id }"
            @click="selectCategory(category.id)"
          >
            <span class="category-icon">{{ category.icon }}</span>
            <span class="category-name">{{ category.name }}</span>
          </button>
        </div>
      </div>

      <!-- 宿舍楼选择 -->
      <div class="form-section">
        <label class="form-label">宿舍楼号</label>
        <div class="dorm-options">
          <button
            v-for="dorm in dormList"
            :key="dorm.id"
            class="dorm-option"
            :class="{ active: form.dormBuilding === dorm.name }"
            @click="selectDorm(dorm.name)"
          >
            {{ dorm.name }}
          </button>
        </div>
      </div>

      <!-- 图片上传 -->
      <div class="form-section">
        <label class="form-label">上传图片</label>
        <div class="upload-section">
          <!-- 图片预览 -->
          <div class="image-preview">
            <div
              v-for="(image, index) in form.images"
              :key="index"
              class="preview-item"
            >
              <img :src="image.url" :alt="'图片' + (index + 1)" class="preview-image">
              <button class="remove-image-btn" @click="removeImage(index)">×</button>
            </div>
            
            <!-- 上传按钮 -->
            <div
              v-if="form.images.length < 6"
              class="upload-btn"
              @click="triggerFileInput"
            >
              <span class="upload-icon">+</span>
              <span class="upload-text">添加图片</span>
              <span class="upload-hint">最多6张</span>
            </div>
          </div>

          <!-- 隐藏的文件输入 -->
          <input
            ref="fileInput"
            type="file"
            multiple
            accept="image/*"
            class="file-input"
            @change="handleFileUpload"
          >
        </div>
      </div>

      <!-- 发布按钮 -->
      <div class="form-actions">
        <button
          class="publish-btn"
          :class="{ disabled: !isFormValid }"
          :disabled="!isFormValid"
          @click="handlePublish"
        >
          {{ isSubmitting ? '发布中...' : '发布帖子' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

// 表单数据
const form = ref({
  title: '',
  content: '',
  category: '',
  dormBuilding: '',
  images: [] // 存储图片对象：{ url: string, file: File }
})

// 状态
const isSubmitting = ref(false)
const fileInput = ref(null)

// 分类选项
const categories = ref([
  { id: 'food', name: '约饭', icon: '🍽️' },
  { id: 'sports', name: '约球', icon: '⚽' },
  { id: 'help', name: '求助', icon: '🙋' },
  { id: 'trade', name: '交易', icon: '💰' },
  { id: 'study', name: '学习', icon: '📚' }
])

// 宿舍楼选项
const dormList = ref([
  { id: 'rong9', name: '榕园9号' },
  { id: 'rong8', name: '榕园8号' },
  { id: 'rong7', name: '榕园7号' },
  { id: 'rong6', name: '榕园6号' },
  { id: 'rong5', name: '榕园5号' },
])

// 当前用户信息
const currentUser = ref({
  id: '1',
  name: '当前用户',
  avatar: '/avatars/current-user.jpg'
})

// 计算属性：表单是否有效
const isFormValid = computed(() => {
  return form.value.title.trim() && 
         form.value.content.trim() && 
         form.value.category && 
         form.value.dormBuilding
})

// 选择分类
const selectCategory = (categoryId) => {
  form.value.category = categoryId
}

// 选择宿舍楼
const selectDorm = (dormName) => {
  form.value.dormBuilding = dormName
}

// 触发文件选择
const triggerFileInput = () => {
  fileInput.value?.click()
}

// 处理文件上传
const handleFileUpload = (event) => {
  const files = Array.from(event.target.files)
  
  // 检查图片数量限制
  if (form.value.images.length + files.length > 6) {
    alert('最多只能上传6张图片')
    return
  }

  files.forEach(file => {
    // 检查文件类型
    if (!file.type.startsWith('image/')) {
      alert('请上传图片文件')
      return
    }

    // 检查文件大小（限制为5MB）
    if (file.size > 5 * 1024 * 1024) {
      alert('图片大小不能超过5MB')
      return
    }

    // 创建预览URL
    const url = URL.createObjectURL(file)
    form.value.images.push({
      url,
      file
    })
  })

  // 清空文件输入，允许重复选择相同文件
  event.target.value = ''
}

// 移除图片
const removeImage = (index) => {
  // 释放URL对象
  URL.revokeObjectURL(form.value.images[index].url)
  form.value.images.splice(index, 1)
}

// 发布帖子
const handlePublish = async () => {
  if (!isFormValid.value || isSubmitting.value) return

  isSubmitting.value = true

  try {
    // 模拟发布过程
    await new Promise(resolve => setTimeout(resolve, 1000))

    // 创建新帖子对象
    const newPost = {
      id: Date.now().toString(),
      userName: currentUser.value.name,
      userAvatar: currentUser.value.avatar,
      time: '刚刚',
      category: form.value.category,
      dormBuilding: form.value.dormBuilding,
      title: form.value.title.trim(),
      content: form.value.content.trim(),
      images: form.value.images.map(img => img.url), // 实际开发中这里应该上传到服务器
      commentCount: 0,
      viewCount: 0,
      likeCount: 0
    }

    // 保存到本地存储（模拟添加到首页）
    savePostToLocal(newPost)

    // 显示成功提示
    alert('帖子发布成功！')

    // 跳转回首页
    router.push('/dormgo')

  } catch (error) {
    console.error('发布失败:', error)
    alert('发布失败，请重试')
  } finally {
    isSubmitting.value = false
  }
}

// 保存帖子到本地存储（模拟添加到首页）
const savePostToLocal = (post) => {
  // 从本地存储获取现有帖子
  const existingPosts = JSON.parse(localStorage.getItem('dormgo_posts') || '[]')
  
  // 添加新帖子到开头
  existingPosts.unshift(post)
  
  // 保存回本地存储
  localStorage.setItem('dormgo_posts', JSON.stringify(existingPosts))
}

// 返回上一页
const goBack = () => {
  router.back()
}

// 清理URL对象
onMounted(() => {
  // 组件卸载时清理所有创建的URL对象
  return () => {
    form.value.images.forEach(image => {
      URL.revokeObjectURL(image.url)
    })
  }
})
</script>

<style scoped>
.publish-post {
  min-height: 100vh;
  background: #f5f5f5;
  padding: 0;
}

/* 返回按钮和标题 */
.back-header {
  background: white;
  padding: 16px 20px;
  border-bottom: 1px solid #e8e8e8;
  position: sticky;
  top: 0;
  z-index: 100;
  display: flex;
  align-items: center;
  gap: 16px;
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
  padding: 8px 12px;
  border-radius: 6px;
  transition: all 0.3s;
}

.back-btn:hover {
  background: #f5f5f5;
  color: #1890ff;
}

.back-icon {
  font-size: 18px;
}

.page-title {
  font-size: 20px;
  font-weight: 600;
  color: #333;
  margin: 0;
}

/* 发布表单 */
.publish-form {
  background: white;
  margin: 20px;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.form-section {
  margin-bottom: 32px;
}

.form-label {
  display: block;
  font-size: 16px;
  font-weight: 600;
  color: #333;
  margin-bottom: 12px;
}

/* 标题输入 */
.title-input {
  width: 100%;
  padding: 12px 16px;
  border: 1px solid #e8e8e8;
  border-radius: 8px;
  font-size: 16px;
  transition: border-color 0.3s;
}

.title-input:focus {
  outline: none;
  border-color: #1890ff;
  box-shadow: 0 0 0 2px rgba(24, 144, 255, 0.1);
}

/* 正文输入 */
.content-textarea {
  width: 100%;
  padding: 12px 16px;
  border: 1px solid #e8e8e8;
  border-radius: 8px;
  font-size: 14px;
  resize: vertical;
  transition: border-color 0.3s;
  font-family: inherit;
  line-height: 1.6;
}

.content-textarea:focus {
  outline: none;
  border-color: #1890ff;
  box-shadow: 0 0 0 2px rgba(24, 144, 255, 0.1);
}

/* 字符计数 */
.char-count {
  text-align: right;
  font-size: 12px;
  color: #999;
  margin-top: 8px;
}

/* 分类选项 */
.category-options {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(120px, 1fr));
  gap: 12px;
}

.category-option {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 16px 12px;
  background: #f8f9fa;
  border: 2px solid #e9ecef;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.category-option:hover {
  background: #e9ecef;
  border-color: #adb5bd;
}

.category-option.active {
  background: #e6f7ff;
  border-color: #1890ff;
  color: #1890ff;
}

.category-icon {
  font-size: 24px;
}

.category-name {
  font-size: 14px;
  font-weight: 500;
}

/* 宿舍楼选项 */
.dorm-options {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(100px, 1fr));
  gap: 12px;
}

.dorm-option {
  padding: 12px 16px;
  background: #f8f9fa;
  border: 2px solid #e9ecef;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  color: #666;
  cursor: pointer;
  transition: all 0.3s ease;
  text-align: center;
}

.dorm-option:hover {
  background: #e9ecef;
  border-color: #adb5bd;
}

.dorm-option.active {
  background: #e6f7ff;
  border-color: #1890ff;
  color: #1890ff;
}

/* 图片上传 */
.upload-section {
  margin-top: 8px;
}

.image-preview {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(100px, 1fr));
  gap: 12px;
}

.preview-item {
  position: relative;
  width: 100px;
  height: 100px;
  border-radius: 8px;
  overflow: hidden;
}

.preview-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border: 1px solid #e8e8e8;
  border-radius: 8px;
}

.remove-image-btn {
  position: absolute;
  top: 4px;
  right: 4px;
  width: 20px;
  height: 20px;
  background: rgba(0, 0, 0, 0.7);
  color: white;
  border: none;
  border-radius: 50%;
  font-size: 14px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background 0.3s;
}

.remove-image-btn:hover {
  background: rgba(0, 0, 0, 0.9);
}

.upload-btn {
  width: 100px;
  height: 100px;
  border: 2px dashed #d9d9d9;
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 4px;
  cursor: pointer;
  transition: all 0.3s ease;
  background: #fafafa;
}

.upload-btn:hover {
  border-color: #1890ff;
  background: #f0f8ff;
}

.upload-icon {
  font-size: 24px;
  color: #666;
}

.upload-text {
  font-size: 12px;
  color: #666;
}

.upload-hint {
  font-size: 10px;
  color: #999;
}

.file-input {
  display: none;
}

/* 发布按钮 */
.form-actions {
  margin-top: 40px;
  text-align: center;
}

.publish-btn {
  width: 100%;
  max-width: 300px;
  padding: 16px 32px;
  background: #1890ff;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.publish-btn:hover:not(.disabled) {
  background: #40a9ff;
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(24, 144, 255, 0.3);
}

.publish-btn.disabled {
  background: #ccc;
  cursor: not-allowed;
  transform: none;
  box-shadow: none;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .publish-form {
    margin: 12px;
    padding: 16px;
  }
  
  .back-header {
    padding: 12px 16px;
  }
  
  .category-options {
    grid-template-columns: repeat(3, 1fr);
  }
  
  .dorm-options {
    grid-template-columns: repeat(2, 1fr);
  }
  
  .image-preview {
    grid-template-columns: repeat(3, 1fr);
  }
}
</style>