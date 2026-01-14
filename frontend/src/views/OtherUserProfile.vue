<template>
  <div class="profile-page">
    <div class="nav-bar">
      <button class="back-btn" @click="$router.back()">← 返回</button>
      <span class="page-title">用户主页</span>
      <div style="width: 40px;"></div> </div>

    <div class="profile-header">
      <div class="user-info-card">
        <div class="avatar-container">
          <img :src="userInfo.avatarurl || defaultAvatar" class="large-avatar" />
        </div>
        <h2 class="user-name">{{ userInfo.username || '加载中...' }}</h2>
        <div class="user-tags" v-if="userInfo.dorm">
          <span class="dorm-badge">
            🏠 {{ userInfo.dorm.dormname }}
          </span>
        </div>

        <p class="user-bio">{{ userInfo.intro || '这个人很懒，什么都没写' }}</p>
        
      </div>
    </div>

    <div class="profile-content">
      <div class="content-tabs">
        <div class="tab-title">TA 的发布</div>
      </div>
      
      <div class="list-container">
        <div v-if="userPosts.length > 0" class="post-list">
          <Card 
            v-for="post in userPosts" 
            :key="post.id" 
            :post="post" 
          />
        </div>
        <div v-else class="empty-state">
          <div class="empty-emoji">🍃</div>
          <p>TA 还没有发布过内容</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Card from '@/components/Card.vue'
import api from '@/api/index.ts'
import { useUserStore } from '@/stores/user'; // 用于检查登录状态

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const userId = route.params.id
const userInfo = ref({})
const userPosts = ref([])
const defaultAvatar = 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'

// 获取用户信息
const fetchUserInfo = async () => {
  try {
    const res = await api.get(`/api/v1/user/info?user_id=${userId}`)
    if (res.data.code === 200) {
      userInfo.value = res.data.data
      console.log('用户信息:', userInfo.value);
    }
  } catch (err) {
    console.error('获取用户信息失败', err)
  }
}

// 获取用户帖子
const fetchUserPosts = async () => {
  try {
    // 假设后端有获取指定用户帖子的接口，如果没有，需要后端补充
    // 例如 GET /api/v1/post/user_posts?user_id=xxx
    const res = await api.get(`/api/v1/post/user_posts`, {
        params: { user_id: userId }
    })
    if (res.data.code === 200) {
      userPosts.value = res.data.data
    }
  } catch (err) {
    console.error('获取用户帖子失败', err)
  }
}

const handleContact = () => {
  if (!userStore.isLoggedIn) {
      alert("请先登录")
      return
  }
  router.push({
    name: 'Chat',
    params: { id: userId }
  })
}

onMounted(() => {
  fetchUserInfo()
  fetchUserPosts()
})
</script>

<style scoped>
/* 复用 PersonalHome.vue 的大部分样式 */
.profile-page { min-height: 100vh; background: #f5f5f5; }
.nav-bar { background: white; padding: 15px 20px; display: flex; justify-content: space-between; align-items: center; box-shadow: 0 1px 4px rgba(0,0,0,0.05); }
.back-btn { border: none; background: none; font-size: 16px; cursor: pointer; color: #666; }
.page-title { font-weight: 600; font-size: 16px; }
.profile-header { background: white; padding: 30px 20px; margin-bottom: 12px; text-align: center; }
.large-avatar { width: 80px; height: 80px; border-radius: 50%; object-fit: cover; border: 3px solid #f0f0f0; margin-bottom: 12px; }
.user-name { font-size: 20px; color: #333; margin: 0 0 4px 0; }
.user-bio { font-size: 13px; color: #999; margin-bottom: 20px; }

.action-buttons { margin-top: 15px; }
.contact-btn { background: #1890ff; color: white; border: none; padding: 8px 30px; border-radius: 20px; font-size: 14px; cursor: pointer; box-shadow: 0 4px 12px rgba(24,144,255,0.3); }

.profile-content { max-width: 800px; margin: 0 auto; padding: 0 10px; }
.content-tabs { background: white; border-radius: 12px 12px 0 0; padding: 16px 20px; border-bottom: 1px solid #f0f0f0; }
.tab-title { font-weight: 600; color: #333; position: relative; display: inline-block; }
.tab-title::after { content: ''; position: absolute; bottom: -16px; left: 0; width: 100%; height: 3px; background: #1890ff; border-radius: 2px; }

.list-container { padding-top: 10px; }
.empty-state { text-align: center; padding: 60px 0; color: #999; }
.empty-emoji { font-size: 40px; margin-bottom: 10px; }

/* 【新增】宿舍标签样式 */
.user-tags {
  margin: 8px 0 12px 0;
  display: flex;
  justify-content: center; /* 居中显示 */
}

.dorm-badge {
  display: inline-flex;
  align-items: center;
  padding: 4px 10px;
  background-color: #f6ffed; /* 绿色背景，区分自己 */
  color: #52c41a;            /* 绿色文字 */
  border-radius: 12px;
  font-size: 13px;
  font-weight: 500;
  border: 1px solid #b7eb8f;
}
</style>