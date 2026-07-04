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

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

interface FAQ {
  question: string
  answer: string
  expanded: boolean
}

const faqs = ref<FAQ[]>([
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
  router.push('/personalhome')
}

const toggleFaq = (index: number) => {
  const faq = faqs.value[index]
  if (!faq) return
  faq.expanded = !faq.expanded
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
