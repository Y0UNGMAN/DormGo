<template>
  <div class="profile-page">
    <!-- 顶部导航栏 -->
    <div class="nav-bar">
      <button class="back-btn" @click="goHome">← 返回首页</button>
      <span class="page-title">个人中心</span>
      <button class="logout-btn" @click="handleLogout">退出登录</button>
    </div>

    <!-- 用户信息区域 -->
    <div class="profile-header">
      <!-- 左上角宿舍楼选择 -->
      <div class="dorm-selector">
        <el-select v-model="selectedDorm" placeholder="选择宿舍楼" size="small">
          <el-option
            v-for="dorm in dormOptions"
            :key="dorm.value"
            :label="dorm.label"
            :value="dorm.value"
          />
        </el-select>
      </div>

      <!-- 用户头像和昵称 -->
      <div class="user-info-card">
        <div class="avatar-container">
          <img :src="userInfo.avatar || defaultAvatar" alt="头像" class="large-avatar" />
        </div>
        <h2 class="user-name">{{ userInfo.nickname || userInfo.name || '未设置昵称' }}</h2>
      </div>
    </div>

    <!-- 四个功能板块 -->
    <div class="function-grid">
      <div class="function-item" @click="goToPage('/profile/post-history')">
        <div class="function-icon">📝</div>
        <span class="function-text">发帖历史</span>
      </div>
      <div class="function-item" @click="goToPage('/profile/favorites')">
        <div class="function-icon">⭐</div>
        <span class="function-text">我的收藏</span>
      </div>
      <div class="function-item" @click="goToPage('/profile/notifications')">
        <div class="function-icon">🔔</div>
        <span class="function-text">消息通知</span>
      </div>
      <div class="function-item" @click="goToPage('/profile/coins')">
        <div class="function-icon">🪙</div>
        <span class="function-text">寝友币</span>
      </div>
    </div>

    <!-- 三个设置项 -->
    <div class="settings-section">
      <div class="settings-item" @click="goToPage('/profile/detail')">
        <div class="settings-left">
          <span class="settings-icon">📌</span>
          <span class="settings-text">详细资料</span>
        </div>
        <span class="settings-arrow">›</span>
      </div>
      <div class="settings-item" @click="goToPage('/profile/rules')">
        <div class="settings-left">
          <span class="settings-icon">📖</span>
          <span class="settings-text">社区规范</span>
        </div>
        <span class="settings-arrow">›</span>
      </div>
      <div class="settings-item" @click="goToPage('/profile/contact')">
        <div class="settings-left">
          <span class="settings-icon">📞</span>
          <span class="settings-text">联系我们</span>
        </div>
        <span class="settings-arrow">›</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

// 默认头像
const defaultAvatar = 'https://dorm-go.oss-cn-guangzhou.aliyuncs.com/avator/midnight.jpg'

// 用户信息
const userInfo = ref({
  name: '',
  nickname: '',
  avatar: ''
})

// 宿舍楼选择
const selectedDorm = ref('榕8栋9栋')
const dormOptions = [
  { value: '榕8栋1栋', label: '榕8栋1栋' },
  { value: '榕8栋2栋', label: '榕8栋2栋' },
  { value: '榕8栋3栋', label: '榕8栋3栋' },
  { value: '榕8栋4栋', label: '榕8栋4栋' },
  { value: '榕8栋5栋', label: '榕8栋5栋' },
  { value: '榕8栋6栋', label: '榕8栋6栋' },
  { value: '榕8栋7栋', label: '榕8栋7栋' },
  { value: '榕8栋8栋', label: '榕8栋8栋' },
  { value: '榕8栋9栋', label: '榕8栋9栋' },
  { value: '榕8栋10栋', label: '榕8栋10栋' },
  { value: '榕8栋11栋', label: '榕8栋11栋' },
  { value: '榕8栋12栋', label: '榕8栋12栋' }
]

const goHome = () => {
  router.push('/home')
}

const goToPage = (path) => {
  router.push(path)
}

const handleLogout = () => {
  if(confirm('确定要退出登录吗？')) {
    localStorage.removeItem('userToken')
    localStorage.removeItem('userInfo')
    router.push('/user/login')
  }
}

const fetchUserInfo = () => {
  const storedUser = localStorage.getItem('userInfo')
  if (storedUser) {
    userInfo.value = JSON.parse(storedUser)
  } else {
    router.push('/user/login')
  }
}

onMounted(() => {
  fetchUserInfo()
})
</script>

<style scoped>
.profile-page {
  min-height: 100vh;
  background: #f5f5f5;
}

/* 顶部导航 */
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

.back-btn, .logout-btn {
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

.logout-btn {
  color: #ff4d4f;
}

.logout-btn:hover {
  color: #ff7875;
}

.page-title {
  font-weight: 600;
  font-size: 16px;
  color: #333;
}

/* 用户信息区域 */
.profile-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 30px 20px;
  position: relative;
}

.dorm-selector {
  position: absolute;
  top: 15px;
  left: 15px;
}

.dorm-selector :deep(.el-select) {
  width: 120px;
}

.dorm-selector :deep(.el-input__wrapper) {
  background: rgba(255, 255, 255, 0.9);
  border-radius: 20px;
  box-shadow: none;
}

.user-info-card {
  text-align: center;
  padding-top: 10px;
}

.avatar-container {
  width: 80px;
  height: 80px;
  margin: 0 auto 12px;
}

.large-avatar {
  width: 100%;
  height: 100%;
  border-radius: 50%;
  object-fit: cover;
  border: 3px solid rgba(255, 255, 255, 0.8);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.user-name {
  font-size: 20px;
  color: white;
  margin: 0;
  font-weight: 600;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
}

/* 四个功能板块 */
.function-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  background: white;
  margin: 12px;
  border-radius: 12px;
  padding: 20px 10px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.function-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  padding: 10px;
  border-radius: 8px;
  transition: all 0.3s ease;
}

.function-item:hover {
  background: #f5f5f5;
  transform: translateY(-2px);
}

.function-icon {
  font-size: 28px;
}

.function-text {
  font-size: 13px;
  color: #333;
  font-weight: 500;
}

/* 设置项 */
.settings-section {
  background: white;
  margin: 0 12px 12px;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.settings-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid #f0f0f0;
  cursor: pointer;
  transition: background 0.3s;
}

.settings-item:last-child {
  border-bottom: none;
}

.settings-item:hover {
  background: #f9f9f9;
}

.settings-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.settings-icon {
  font-size: 20px;
}

.settings-text {
  font-size: 15px;
  color: #333;
}

.settings-arrow {
  font-size: 20px;
  color: #ccc;
}

/* 响应式 */
@media (max-width: 480px) {
  .function-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: 10px;
  }
  
  .function-item {
    padding: 15px 10px;
  }
}
</style>