<template>
  <div class="page-container">
    <!-- 顶部导航 -->
    <div class="nav-bar">
      <button class="back-btn" @click="goBack">← 返回个人中心</button>
      <span class="page-title">社区规范</span>
      <span class="placeholder"></span>
    </div>

    <div class="content-area">
      <!-- 社区行为准则 -->
      <div class="section-card">
        <div class="section-header">
          <span class="section-icon">📜</span>
          <h2 class="section-title">社区行为准则</h2>
        </div>
        <div class="section-content">
          <div class="rule-item" v-for="(rule, index) in rules" :key="index">
            <span class="rule-number">{{ index + 1 }}</span>
            <span class="rule-text">{{ rule }}</span>
          </div>
        </div>
      </div>

      <!-- 奖惩机制 -->
      <div class="section-card">
        <div class="section-header">
          <span class="section-icon">⚖️</span>
          <h2 class="section-title">奖惩机制</h2>
        </div>
        <div class="section-content">
          <div class="mechanism-group">
            <h3 class="mechanism-title reward">
              <span class="title-icon">🎁</span> 奖励机制
            </h3>
            <ul class="mechanism-list">
              <li>优质内容可获得寝友币奖励</li>
              <li>帮助他人解决问题可获得感谢奖励</li>
              <li>积极参与社区活动可获得额外奖励</li>
            </ul>
          </div>

          <div class="mechanism-group">
            <h3 class="mechanism-title punish">
              <span class="title-icon">⚠️</span> 惩罚机制
            </h3>
            <ul class="mechanism-list">
              <li>违规行为将受到相应处罚</li>
              <li>发布违规内容将被删除并扣除寝友币</li>
              <li>累计警告达到3次将暂停账号</li>
              <li>严重违规将永久封禁账号</li>
            </ul>
          </div>
        </div>
      </div>

      <!-- 内容规范 -->
      <div class="section-card">
        <div class="section-header">
          <span class="section-icon">📝</span>
          <h2 class="section-title">内容规范</h2>
        </div>
        <div class="section-content">
          <div class="content-rule">
            <h4 class="rule-category forbidden">🚫 禁止发布</h4>
            <ul class="rule-details">
              <li>违法违规信息</li>
              <li>色情、暴力、血腥内容</li>
              <li>人身攻击、辱骂他人</li>
              <li>虚假、欺诈信息</li>
              <li>广告、垃圾信息</li>
              <li>泄露他人隐私</li>
            </ul>
          </div>

          <div class="content-rule">
            <h4 class="rule-category encourage">✅ 鼓励发布</h4>
            <ul class="rule-details">
              <li>真诚的互助请求</li>
              <li>有价值的分享内容</li>
              <li>友善的交流互动</li>
              <li>积极正能量的内容</li>
            </ul>
          </div>
        </div>
      </div>

      <!-- 免责声明 -->
      <div class="disclaimer">
        <p>本规范最终解释权归寝友Go运营团队所有</p>
        <p>如有疑问，请联系客服</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const rules = ref([
  '尊重他人，文明交流',
  '禁止发布违法、违规内容',
  '保护个人隐私，不泄露他人信息',
  '共同维护良好的社区环境',
  '积极举报不良内容'
])

const goBack = () => {
  router.push('/personalhome')
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

.section-card {
  background: white;
  border-radius: 12px;
  margin-bottom: 12px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.section-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 16px 20px;
  display: flex;
  align-items: center;
  gap: 10px;
}

.section-icon {
  font-size: 24px;
}

.section-title {
  font-size: 18px;
  font-weight: 600;
  color: white;
  margin: 0;
}

.section-content {
  padding: 20px;
}

/* 行为准则 */
.rule-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 12px 0;
  border-bottom: 1px dashed #f0f0f0;
}

.rule-item:last-child {
  border-bottom: none;
}

.rule-number {
  width: 24px;
  height: 24px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 600;
  flex-shrink: 0;
}

.rule-text {
  font-size: 15px;
  color: #333;
  line-height: 1.6;
}

/* 奖惩机制 */
.mechanism-group {
  margin-bottom: 20px;
}

.mechanism-group:last-child {
  margin-bottom: 0;
}

.mechanism-title {
  font-size: 16px;
  font-weight: 600;
  margin: 0 0 12px 0;
  display: flex;
  align-items: center;
  gap: 6px;
}

.mechanism-title.reward {
  color: #52c41a;
}

.mechanism-title.punish {
  color: #ff4d4f;
}

.title-icon {
  font-size: 18px;
}

.mechanism-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.mechanism-list li {
  font-size: 14px;
  color: #666;
  padding: 8px 0 8px 24px;
  position: relative;
  border-bottom: 1px solid #f5f5f5;
}

.mechanism-list li:last-child {
  border-bottom: none;
}

.mechanism-list li::before {
  content: '•';
  position: absolute;
  left: 8px;
  color: #999;
}

/* 内容规范 */
.content-rule {
  margin-bottom: 20px;
}

.content-rule:last-child {
  margin-bottom: 0;
}

.rule-category {
  font-size: 15px;
  font-weight: 600;
  margin: 0 0 10px 0;
}

.rule-category.forbidden {
  color: #ff4d4f;
}

.rule-category.encourage {
  color: #52c41a;
}

.rule-details {
  list-style: none;
  padding: 0;
  margin: 0;
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 8px;
}

.rule-details li {
  font-size: 13px;
  color: #666;
  padding: 6px 12px;
  background: #f8f9fa;
  border-radius: 6px;
}

/* 免责声明 */
.disclaimer {
  text-align: center;
  padding: 20px;
  color: #999;
  font-size: 12px;
}

.disclaimer p {
  margin: 4px 0;
}

@media (max-width: 480px) {
  .rule-details {
    grid-template-columns: 1fr;
  }
}
</style>
