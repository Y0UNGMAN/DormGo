<template>
  <div class="profile-page">
    <div class="nav-bar">
      <button class="back-btn" @click="goHome">← 返回首页</button>
      <span class="page-title">个人中心</span>
      <button class="logout-btn" @click="handleLogout">退出登录</button>
    </div>

    <div class="profile-header">
      <div class="user-info-card">
        <div class="avatar-container">
          <img :src="userInfo.avatar || defaultAvatar" alt="头像" class="large-avatar" />
          <div class="edit-avatar-badge">📷</div>
        </div>
        <h2 class="user-name">{{ userInfo.name }}</h2>
        <p class="user-bio">暂无个性签名...</p>
        
        <div class="user-stats">
          <div class="stat-item">
            <span class="stat-num">{{ myPosts.length }}</span>
            <span class="stat-label">发布</span>
          </div>
          <div class="stat-item">
            <span class="stat-num">0</span>
            <span class="stat-label">获赞</span>
          </div>
          <div class="stat-item">
            <span class="stat-num">0</span>
            <span class="stat-label">收藏</span>
          </div>
        </div>
      </div>
    </div>

    <div class="profile-content">
      <div class="content-tabs">
        <button 
          class="tab-btn" 
          :class="{ active: activeTab === 'posts' }"
          @click="activeTab = 'posts'"
        >
          我的发布
        </button>
        <button 
          class="tab-btn" 
          :class="{ active: activeTab === 'likes' }"
          @click="activeTab = 'likes'"
        >
          我的收藏
        </button>
      </div>

      <div class="list-container">
        <div v-if="displayPosts.length > 0" class="post-list">
          <Card 
            v-for="post in displayPosts" 
            :key="post.id" 
            :post="post" 
            class="profile-post-card"
          />
        </div>
        
        <div v-else class="empty-state">
          <div class="empty-emoji">🍃</div>
          <p>这里空空如也</p>
          <button v-if="activeTab === 'posts'" class="go-publish-btn" @click="goHome">去发布</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import Card from '@/components/Card.vue' 
import request from '@/utils/request' // 【修改点】

const router = useRouter()
const activeTab = ref('posts') // posts 或 likes
const defaultAvatar = 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'

// 用户信息
const userInfo = ref({
  name: '加载中...',
  avatar: '',
  id: ''
})

// 我的帖子列表
const myPosts = ref([])
// 我的收藏列表 (暂未实现接口)
const likedPosts = ref([])

const goHome = () => {
  router.push('/dormgo')
}

const handleLogout = () => {
  if(confirm('确定要退出登录吗？')) {
    localStorage.removeItem('userToken')
    localStorage.removeItem('userInfo')
    router.push('/login')
  }
}

const fetchUserInfo = () => {
  const storedUser = localStorage.getItem('userInfo')
  if (storedUser) {
    const user = JSON.parse(storedUser)
    userInfo.value = {
      name: user.name || user.username || '用户',
      avatar: user.avatarurl || defaultAvatar,
      id: user.id
    }
  } else {
    router.push('/login')
  }
}

const fetchMyPosts = async () => {
  try {
     // 这里我们复用获取所有帖子的接口，然后在前端过滤（或者后端应提供 /api/v1/post/my_posts）
     // 临时方案：获取所有帖子并筛选
     const res = await request.get('/api/v1/post/posts')
     const allPosts = res.data || []
     // 筛选出发布者ID等于当前用户ID的帖子
     // 注意：后端返回的 publisherid 类型可能和 userInfo.id 类型不一致（string/number），用 == 比较
     myPosts.value = allPosts.filter(p => p.publisherid == userInfo.value.id)
  } catch (error) {
    console.error('获取个人帖子失败', error)
  }
}

// 计算当前显示的列表
const displayPosts = computed(() => {
  return activeTab.value === 'posts' ? myPosts.value : likedPosts.value
})

onMounted(() => {
  fetchUserInfo()
  fetchMyPosts()
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
}

.logout-btn {
  color: #ff4d4f; /* 红色警告色 */
}

.page-title {
  font-weight: 600;
  font-size: 16px;
}

/* 头部信息 */
.profile-header {
  background: white;
  padding: 30px 20px;
  margin-bottom: 12px;
  text-align: center;
}

.user-info-card {
  max-width: 600px;
  margin: 0 auto;
}

.avatar-container {
  position: relative;
  width: 80px;
  height: 80px;
  margin: 0 auto 12px;
}

.large-avatar {
  width: 100%;
  height: 100%;
  border-radius: 50%;
  object-fit: cover;
  border: 3px solid #f0f0f0;
}

.edit-avatar-badge {
  position: absolute;
  bottom: 0;
  right: 0;
  background: #1890ff;
  color: white;
  width: 24px;
  height: 24px;
  border-radius: 50%;
  font-size: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 2px solid white;
  cursor: pointer;
}

.user-name {
  font-size: 20px;
  color: #333;
  margin: 0 0 4px 0;
}

.user-bio {
  font-size: 13px;
  color: #999;
  margin-bottom: 20px;
}

.user-stats {
  display: flex;
  justify-content: center;
  gap: 40px;
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.stat-num {
  font-size: 18px;
  font-weight: 600;
  color: #333;
}

.stat-label {
  font-size: 12px;
  color: #999;
  margin-top: 4px;
}

/* 内容区域 */
.profile-content {
  max-width: 800px;
  margin: 0 auto;
  padding: 0 10px;
}

.content-tabs {
  background: white;
  border-radius: 12px 12px 0 0;
  display: flex;
  padding: 0 20px;
  margin-bottom: 2px;
}

.tab-btn {
  padding: 16px 0;
  margin-right: 30px;
  background: none;
  border: none;
  font-size: 15px;
  color: #666;
  cursor: pointer;
  position: relative;
}

.tab-btn.active {
  color: #1890ff;
  font-weight: 600;
}

.tab-btn.active::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 50%;
  transform: translateX(-50%);
  width: 20px;
  height: 3px;
  background: #1890ff;
  border-radius: 2px;
}

.list-container {
  min-height: 300px;
  padding-top: 10px;
}

/* 空状态 */
.empty-state {
  text-align: center;
  padding: 60px 0;
  color: #999;
}

.empty-emoji {
  font-size: 40px;
  margin-bottom: 10px;
}

.go-publish-btn {
  margin-top: 15px;
  background: #1890ff;
  color: white;
  border: none;
  padding: 8px 20px;
  border-radius: 20px;
  cursor: pointer;
}
</style>