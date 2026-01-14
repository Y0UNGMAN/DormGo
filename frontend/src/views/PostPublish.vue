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


        <div class="textarea-wrapper">
          <textarea
            v-model="form.content"
            placeholder="请输入帖子内容... (试试点击右下角的 AI 润色)"
            class="content-textarea"
            rows="8"
            maxlength="1000"
          ></textarea>
          <button 
            class="ai-polish-btn" 
            @click.prevent="handleAiPolish" 
            :disabled="isAiLoading || !form.content"
            title="让 AI 帮你写得更吸引人"
          >
            <span v-if="!isAiLoading">✨ AI 帮我润色</span>
            <span v-else>🤖 正在思考...</span>
          </button>
        </div>
  

        <div class="char-count">{{ form.content.length }}/1000</div>
      </div>

      <!-- 分类选择 -->
      <div class="form-section">
        <label class="form-label">分类</label>
        <div class="category-options">
          <button
            v-for="type in types"
            :key="type.typeid"
            class="category-option"
            :class="{ active: form.typeid === type.typeid }"
            @click="selectCategory(type.typeid)"
          >
            <span class="category-name">{{ type.typename }}</span>
          </button>
        </div>
      </div>

      <!-- 宿舍楼选择 -->
      <div class="form-section">
        <label class="form-label">宿舍楼号</label>
        <div class="dorm-options">
          <button
            v-for="dorm in dormList"
            :key="dorm.dormid"
            class="dorm-option"
            :class="{ active: form.dormid === dorm.dormid }"
            @click="selectDorm(dorm.dormid)"
          >
            {{ dorm.dormname }}
          </button>
        </div>
      </div>

      <!-- 报名信息选择 -->
      <div class="form-section">
        <div class="limit-switch-container">
          <label class="form-label">限时/报名</label>
          <div class="switch" :class="{ checked: form.isLimited }" @click="toggleLimit">
            <div class="slider"></div>
          </div>
        </div>

        <div v-if="form.isLimited" class="limit-settings">
          <div class="setting-item">
            <label class="sub-label">截止时间</label>
            <input 
              type="datetime-local" 
              v-model="form.deadline"
              class="date-input"
            >
          </div>
          <div class="setting-item">
            <label class="sub-label">最大人数 (0为不限)</label>
            <input 
              type="number" 
              v-model.number="form.maxEnrollment"
              class="number-input" 
              min="1"
            >
          </div>
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
          <span v-if="!isSubmitting">发布帖子</span>
          <span v-else class="loading-text">
            <i class="loading-icon">↻</i> 安全检测中...
          </span>
        </button>
      </div>
    </div>
    <div v-if="showModal" class="modal-overlay" @click.self="closeModal">
      <div class="modal-card" :class="modalType">
        <div class="modal-icon">
          <span v-if="modalType === 'success'">🎉</span>
          <span v-else-if="modalType === 'error'">😭</span>
          <span v-else>⚠️</span>
        </div>
        <h3 class="modal-title">{{ modalTitle }}</h3>
        <p class="modal-message">{{ modalMessage }}</p>
        <button class="modal-btn" @click="handleModalConfirm">
          {{ modalBtnText }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios';
import api from '@/api/index.ts';
import { useUserStore } from '@/stores/user';
const userStore = useUserStore();
const router = useRouter()

// 当前用户信息
const currentUser = computed(() => userStore.currentUser)
const currentUserId = computed(() => userStore.currentUserId)
// 表单数据
const form = ref({
  title: '',
  content: '',
  typeid: '',
  dormid: '',
  images: [], // 存储图片对象：{ url: string, file: File }
  isLimited: false,
  deadline: '',
  maxEnrollment: 0
})

// --- 【新增】弹窗相关状态 ---
const showModal = ref(false)
const modalType = ref('info') // success, error, warning
const modalTitle = ref('')
const modalMessage = ref('')
const modalBtnText = ref('我知道了')
let modalConfirmCallback = null

// --- 【新增】显示弹窗的辅助函数 ---
const showAlert = (message, type = 'warning', title = '提示', callback = null) => {
  modalMessage.value = message
  modalType.value = type
  modalTitle.value = title
  modalBtnText.value = type === 'success' ? '好的' : '关闭'
  modalConfirmCallback = callback
  showModal.value = true
}

// --- 【新增】关闭/确认弹窗 ---
const closeModal = () => {
  showModal.value = false
  if (modalConfirmCallback) {
    modalConfirmCallback()
    modalConfirmCallback = null
  }
}

const handleModalConfirm = () => {
  closeModal()
}

// AI Loading 状态
const isAiLoading = ref(false)
// 状态
const isSubmitting = ref(false)
const fileInput = ref(null)

// 分类选项
const types = ref([])
const fetchPostTypes = async () => {
  try {
    const response = await api.get('/api/v1/post/post_type');
    if (response.data && response.data.data) {
      // 成功获取数据，并赋值给响应式变量 types
      types.value = response.data.data;
      console.log('帖子分类列表:', types.value);
    } else {
      // 如果数据结构不符合预期
      throw new Error('接口返回数据结构异常');
    }
  } catch (error) {
    console.error('获取帖子分类失败:', error);
  }
};

// 宿舍楼选项
const dormList = ref([])
const fetchDormList = async() => {
  try {
    const response = await api.get('/api/v1/post/dorms');
    if (response.data && response.data.data) {
      // 成功获取数据，并赋值给响应式变量 dormList
      dormList.value = response.data.data;
      console.log('宿舍楼列表:', dormList.value);
    } else {
      // 如果数据结构不符合预期
      throw new Error('接口返回数据结构异常');
    }
  } catch (error) {
    console.error('获取宿舍楼列表失败:', error);
  }
};

// 计算属性：表单是否有效
const isFormValid = computed(() => {
  return form.value.title.trim() && 
         form.value.content.trim() && 
         form.value.typeid && 
         form.value.dormid
})

// 选择分类
const selectCategory = (id) => {
  form.value.typeid = id
}

// 选择宿舍楼
const selectDorm = (id) => {
  form.value.dormid = id
}

// 触发文件选择
const triggerFileInput = () => {
  fileInput.value?.click()
}

// 处理文件上传
const handleFileUpload = (event) => {
  const files = Array.from(event.target.files)
  
  if (form.value.images.length + files.length > 6) {
    alert('最多只能上传6张图片')
    return
  }

  files.forEach(file => {
    if (!file.type.startsWith('image/')) {
      alert('请上传图片文件')
      return
    }
    if (file.size > 5 * 1024 * 1024) {
      alert('图片大小不能超过5MB')
      return
    }
    const url = URL.createObjectURL(file)
    form.value.images.push({
      url,
      file
    })
  })
  event.target.value = ''
}

// 移除图片
const removeImage = (index) => {
  URL.revokeObjectURL(form.value.images[index].url)
  form.value.images.splice(index, 1)
}
// 切换开关
const toggleLimit = () => {
  form.value.isLimited = !form.value.isLimited
}

// 发布帖子
const handlePublish = async () => {
  if (!isFormValid.value || isSubmitting.value) return
  isSubmitting.value = true

  if (!userStore.isLoggedIn) {
      showAlert("请先登录才能发布帖子！", "warning", "未登录", () => router.push('/'))
      isSubmitting.value = false // 别忘了重置 loading
      return
    }
  if (form.value.isLimited) {
    if (!form.value.deadline) {
        showAlert("请选择截止时间", "warning")
        isSubmitting.value = false
        return
    }
    // 将 datetime-local 的时间转换为 RFC3339 格式 (加时区，或者让后端处理)
    // datetime-local 格式如 "2023-11-01T12:00"，可以直接传给后端，GORM通常能解析
  }
  try {
    const formData = new FormData();
    formData.append('publisherid', currentUserId.value);
    formData.append('dormid', form.value.dormid); 
    formData.append('typeid', form.value.typeid);
    formData.append('title', form.value.title);
    formData.append('content', form.value.content);
    formData.append('is_limited', form.value.isLimited);
    if (form.value.isLimited) {
        formData.append('deadline', new Date(form.value.deadline).toISOString());
        formData.append('max_enrollment', form.value.maxEnrollment);
    }
    form.value.images.forEach((image) => {
      formData.append(`images`, image.file);
    });
    const res = await api.post('/api/v1/post/create', formData,{
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
    console.log('创建帖子响应:', res.data);
    if (res.data.code === 200) {
      showAlert('发布成功！快去看看吧', 'success', '发布成功', () => {
        router.push('/dormgo')
      })
    } else {
      // 这里会显示 "内容审核未通过: xxx"
      showAlert(res.data.msg || "发布失败", 'error', '发布被拦截')
    }
    
  } catch (error) {
    console.error('网络崩溃:', error)
    showAlert("网络连接异常，请检查网络", "error", "网络错误")
    // request.ts 拦截器会处理错误提示
  } finally {
    isSubmitting.value = false
  }
}

const handleAiPolish = async () => {
  // 简单校验
  if (!form.value.content || form.value.content.trim().length < 2) {
    showAlert("请至少写几个字，AI 才能帮你润色哦~", "warning")
    return
  }

  isAiLoading.value = true
  try {
    // 调用我们在后端新写的接口
    const res = await api.post('/api/v1/post/ai_polish', 
    {
      content: form.value.content
    },
    {
        timeout: 30000 // 设置为 60000 毫秒 (60秒)，给 AI 足够的思考时间
    }
  )
    
    if (res.data.code === 200) {
      // 成功！用 AI 的结果覆盖当前内容
      form.value.content = res.data.data
    } else {
      alert(res.data.msg || "润色失败，请稍后再试")
    }
  } catch (err) {
    console.error("AI 接口调用失败:", err)
    alert("网络开小差了，AI 暂时无法连接")
  } finally {
    isAiLoading.value = false
  }
}





// 返回上一页
const goBack = () => {
  router.back()
}

// 初始化
onMounted(() => {
  // 获取当前用户信息
  const userStr = localStorage.getItem('userInfo')
  if (userStr) {
     currentUser.value = JSON.parse(userStr)
  }
  
  fetchDormList();
  fetchPostTypes();
  
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
  padding-bottom: 40px;
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

.limit-switch-container {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

/* 开关样式 */
.switch {
  width: 50px;
  height: 26px;
  background-color: #ccc;
  border-radius: 13px;
  position: relative;
  cursor: pointer;
  transition: background-color 0.3s;
}
.switch.checked {
  background-color: #1890ff;
}
.slider {
  width: 22px;
  height: 22px;
  background-color: white;
  border-radius: 50%;
  position: absolute;
  top: 2px;
  left: 2px;
  transition: transform 0.3s;
}
.switch.checked .slider {
  transform: translateX(24px);
}

.limit-settings {
  background: #f8f9fa;
  padding: 15px;
  border-radius: 8px;
  animation: slideDown 0.3s ease;
}

.setting-item {
  margin-bottom: 12px;
}

.sub-label {
  display: block;
  font-size: 14px;
  color: #666;
  margin-bottom: 6px;
}

.date-input, .number-input {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid #e8e8e8;
  border-radius: 6px;
  font-size: 14px;
}

@keyframes slideDown {
  from { opacity: 0; transform: translateY(-10px); }
  to { opacity: 1; transform: translateY(0); }
}

.textarea-wrapper {
  position: relative;
}

.ai-polish-btn {
  position: absolute;
  bottom: 15px; /* 距离底部 */
  right: 15px;  /* 距离右侧 */
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); /* 紫色渐变 */
  color: white;
  border: none;
  padding: 6px 14px;
  border-radius: 20px;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  box-shadow: 0 4px 10px rgba(118, 75, 162, 0.3);
  transition: all 0.3s ease;
  z-index: 10;
  display: flex;
  align-items: center;
  gap: 6px;
}

.ai-polish-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 15px rgba(118, 75, 162, 0.4);
}

.ai-polish-btn:disabled {
  background: #d9d9d9;
  cursor: not-allowed;
  transform: none;
  box-shadow: none;
}

/* 新增：发布按钮 Loading 动画 */
.loading-text {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.loading-icon {
  display: inline-block;
  font-style: normal;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

/* ... 原有的样式 ... */

/* --- 【新增】弹窗样式 --- */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.6); /* 半透明深色遮罩 */
  backdrop-filter: blur(4px); /* 背景模糊效果，更高级 */
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
  animation: fadeIn 0.3s ease;
}

.modal-card {
  background: white;
  width: 80%;
  max-width: 320px;
  border-radius: 16px;
  padding: 24px;
  text-align: center;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.2);
  transform: scale(0.9);
  animation: popIn 0.3s cubic-bezier(0.175, 0.885, 0.32, 1.275) forwards; /* 弹性弹出动画 */
}

/* 状态颜色边框装饰 */
.modal-card.error { border-top: 5px solid #ff4d4f; }
.modal-card.success { border-top: 5px solid #52c41a; }
.modal-card.warning { border-top: 5px solid #faad14; }

.modal-icon {
  font-size: 40px;
  margin-bottom: 12px;
}

.modal-title {
  font-size: 18px;
  font-weight: 600;
  color: #333;
  margin: 0 0 10px 0;
}

.modal-message {
  font-size: 15px;
  color: #666;
  line-height: 1.6;
  margin-bottom: 24px;
  /* 允许长文本换行 */
  white-space: pre-wrap; 
}

.modal-btn {
  width: 100%;
  padding: 12px 0;
  background: #1890ff;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 15px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s;
}

.modal-btn:hover {
  background: #40a9ff;
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(24, 144, 255, 0.3);
}

.modal-card.error .modal-btn { background: #ff4d4f; }
.modal-card.error .modal-btn:hover { background: #ff7875; box-shadow: 0 4px 12px rgba(255, 77, 79, 0.3); }

/* 动画定义 */
@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

@keyframes popIn {
  from { transform: scale(0.8); opacity: 0; }
  to { transform: scale(1); opacity: 1; }
}

</style>