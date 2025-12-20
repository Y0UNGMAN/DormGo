<template>
  <div class="page-container">
    <!-- 顶部导航 -->
    <div class="nav-bar">
      <button class="back-btn" @click="goBack">← 返回个人中心</button>
      <span class="page-title">联系我们</span>
      <span class="placeholder"></span>
    </div>

    <div class="content-area">
      <!-- 联系方式卡片 -->
      <div class="contact-card">
        <div class="card-header">
          <span class="card-icon">📞</span>
          <h2 class="card-title">客服联系方式</h2>
        </div>
        <div class="card-content">
          <div class="contact-item">
            <span class="contact-label">客服热线</span>
            <span class="contact-value">400-888-9999</span>
          </div>
          <div class="contact-item">
            <span class="contact-label">客服微信</span>
            <span class="contact-value">DormGo_Service</span>
          </div>
          <div class="contact-item">
            <span class="contact-label">客服邮箱</span>
            <span class="contact-value">support@dormgo.com</span>
          </div>
          <div class="contact-item">
            <span class="contact-label">服务时间</span>
            <span class="contact-value">周一至周五 9:00-18:00</span>
          </div>
        </div>
      </div>

      <!-- 反馈表单 -->
      <div class="feedback-card">
        <div class="card-header">
          <span class="card-icon">✉️</span>
          <h2 class="card-title">提交反馈</h2>
        </div>
        <div class="card-content">
          <div class="feedback-form">
            <div class="form-group">
              <label class="form-label">反馈类型</label>
              <el-select v-model="feedback.type" placeholder="请选择反馈类型" class="form-select">
                <el-option label="功能建议" value="suggestion" />
                <el-option label="问题反馈" value="bug" />
                <el-option label="投诉举报" value="complaint" />
                <el-option label="其他" value="other" />
              </el-select>
            </div>

            <div class="form-group">
              <label class="form-label">反馈内容</label>
              <textarea 
                v-model="feedback.content" 
                class="form-textarea" 
                placeholder="请详细描述您的问题或建议..."
                rows="6"
              ></textarea>
              <div class="char-count">{{ feedback.content.length }}/500</div>
            </div>

            <div class="form-group">
              <label class="form-label">联系方式（选填）</label>
              <input 
                v-model="feedback.contact" 
                type="text" 
                class="form-input" 
                placeholder="手机号或邮箱，方便我们联系您"
              />
            </div>

            <button 
              class="submit-btn" 
              :disabled="!canSubmit || isSubmitting"
              @click="submitFeedback"
            >
              {{ isSubmitting ? '提交中...' : '提交反馈' }}
            </button>
          </div>
        </div>
      </div>

      <!-- 常见问题 -->
      <div class="faq-card">
        <div class="card-header">
          <span class="card-icon">❓</span>
          <h2 class="card-title">常见问题</h2>
        </div>
        <div class="card-content">
          <div class="faq-item" v-for="(faq, index) in faqs" :key="index">
            <div class="faq-question" @click="toggleFaq(index)">
              <span>{{ faq.question }}</span>
              <span class="faq-arrow" :class="{ expanded: faq.expanded }">›</span>
            </div>
            <div class="faq-answer" v-show="faq.expanded">
              {{ faq.answer }}
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'

const router = useRouter()

const feedback = ref({
  type: '',
  content: '',
  contact: ''
})

const isSubmitting = ref(false)

const canSubmit = computed(() => {
  return feedback.value.type && feedback.value.content.trim().length >= 10
})

const faqs = ref([
  {
    question: '如何获取寝友币？',
    answer: '您可以通过发布优质帖子、每日签到、帮助他人等方式获取寝友币。详情请查看个人中心的寝友币页面。',
    expanded: false
  },
  {
    question: '发布的帖子为什么被删除了？',
    answer: '帖子可能因违反社区规范被删除，如含有违法违规内容、广告信息等。请查阅社区规范了解详情。',
    expanded: false
  },
  {
    question: '如何修改个人资料？',
    answer: '点击个人中心 → 详细资料，即可修改您的昵称、头像、联系方式等个人信息。',
    expanded: false
  },
  {
    question: '忘记密码怎么办？',
    answer: '在登录页面点击"忘记密码"，通过绑定的手机号或邮箱进行密码重置。',
    expanded: false
  }
])

const goBack = () => {
  router.push('/profile')
}

const toggleFaq = (index) => {
  faqs.value[index].expanded = !faqs.value[index].expanded
}

const submitFeedback = async () => {
  if (!canSubmit.value || isSubmitting.value) return
  
  isSubmitting.value = true
  
  try {
    // 模拟提交
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    // 实际应调用 API
    // await axios.post('/api/v1/feedback', feedback.value)
    
    ElMessage.success('反馈提交成功，感谢您的宝贵意见！')
    
    // 清空表单
    feedback.value = {
      type: '',
      content: '',
      contact: ''
    }
  } catch (error) {
    ElMessage.error('提交失败，请稍后重试')
  } finally {
    isSubmitting.value = false
  }
}
</script>

<style scoped>
.page-container {
  min-height: 100vh;
  background: #f5f5f5;
  padding-bottom: 30px;
}

.nav-bar {
  background: white;
  padding: 15px 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  box-shadow: 0 1px 4px rgba(0,0,0,0.05);
  position: sticky;
  top: 0;
  z-index: 100;
}

.back-btn {
  border: none;
  background: none;
  cursor: pointer;
  font-size: 14px;
  color: #666;
  transition: color 0.3s;
}

.back-btn:hover {
  color: #1890ff;
}

.page-title {
  font-weight: 600;
  font-size: 16px;
  color: #333;
}

.placeholder {
  width: 100px;
}

.content-area {
  padding: 12px;
  max-width: 800px;
  margin: 0 auto;
}

/* 卡片通用样式 */
.contact-card,
.feedback-card,
.faq-card {
  background: white;
  border-radius: 12px;
  margin-bottom: 12px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.card-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 16px 20px;
  display: flex;
  align-items: center;
  gap: 10px;
}

.card-icon {
  font-size: 24px;
}

.card-title {
  font-size: 18px;
  font-weight: 600;
  color: white;
  margin: 0;
}

.card-content {
  padding: 20px;
}

/* 联系方式 */
.contact-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 14px 0;
  border-bottom: 1px solid #f0f0f0;
}

.contact-item:last-child {
  border-bottom: none;
}

.contact-label {
  font-size: 14px;
  color: #666;
}

.contact-value {
  font-size: 14px;
  color: #333;
  font-weight: 500;
}

/* 反馈表单 */
.form-group {
  margin-bottom: 20px;
}

.form-label {
  display: block;
  font-size: 14px;
  color: #333;
  font-weight: 500;
  margin-bottom: 8px;
}

.form-select {
  width: 100%;
}

.form-input {
  width: 100%;
  padding: 12px 14px;
  border: 1px solid #e8e8e8;
  border-radius: 8px;
  font-size: 14px;
  transition: all 0.3s;
  box-sizing: border-box;
}

.form-input:focus {
  outline: none;
  border-color: #1890ff;
  box-shadow: 0 0 0 2px rgba(24, 144, 255, 0.1);
}

.form-textarea {
  width: 100%;
  padding: 12px 14px;
  border: 1px solid #e8e8e8;
  border-radius: 8px;
  font-size: 14px;
  resize: vertical;
  font-family: inherit;
  box-sizing: border-box;
}

.form-textarea:focus {
  outline: none;
  border-color: #1890ff;
  box-shadow: 0 0 0 2px rgba(24, 144, 255, 0.1);
}

.char-count {
  text-align: right;
  font-size: 12px;
  color: #999;
  margin-top: 4px;
}

.submit-btn {
  width: 100%;
  padding: 14px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s;
}

.submit-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
}

.submit-btn:disabled {
  background: #ccc;
  cursor: not-allowed;
  transform: none;
  box-shadow: none;
}

/* FAQ */
.faq-item {
  border-bottom: 1px solid #f0f0f0;
}

.faq-item:last-child {
  border-bottom: none;
}

.faq-question {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 14px 0;
  cursor: pointer;
  font-size: 14px;
  color: #333;
  font-weight: 500;
  transition: color 0.3s;
}

.faq-question:hover {
  color: #1890ff;
}

.faq-arrow {
  font-size: 18px;
  color: #999;
  transition: transform 0.3s;
}

.faq-arrow.expanded {
  transform: rotate(90deg);
}

.faq-answer {
  padding: 0 0 14px 0;
  font-size: 13px;
  color: #666;
  line-height: 1.6;
  background: #f9f9f9;
  padding: 12px;
  border-radius: 6px;
  margin-bottom: 10px;
}
</style>
