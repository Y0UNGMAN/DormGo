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
          <img :src="currentUser.avatarurl" alt="头像" class="large-avatar" />
          <div class="edit-avatar-badge">📷</div>
        </div>
        <h2 class="user-name">{{ currentUser.username }}</h2>
        <p class="user-bio">暂无个性签名...</p>
        
        <div class="user-stats">
          <div class="stat-item">
            <span class="stat-num">{{ myPosts.length }}</span>
            <span class="stat-label">发布</span>
          </div>
          <div class="stat-item">
            <span class="stat-num">12</span>
            <span class="stat-label">获赞</span>
          </div>
          <div class="stat-item">
            <span class="stat-num">5</span>
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
        <button 
          class="tab-btn" 
          :class="{ active: activeTab === 'notifications' }" 
          @click="activeTab = 'notifications'"
        >
          我的通知 <span v-if="unreadCount > 0" class="badge"></span>
        </button>
      </div>
      <div v-if="activeTab === 'notifications'" class="notification-list">
        <div v-for="note in notifications" :key="note.id" class="notification-item" @click="handleNoteClick(note)">
          <img :src="note.sender?.avatarurl || defaultAvatar" class="note-avatar">


          <div class="note-content">
            <template v-if="note.type === 'system'">
              <p class="note-text">
                <span class="highlight" style="color: #ff4d4f">【系统通知】</span> 
                {{ note.content }}
              </p>
            </template>

            <template v-else-if="note.type === 'signup'">
              <p class="note-text">
                <span class="highlight">{{ note.sender?.username }}</span> 
                {{ note.content }}
                <span v-if="note.post" class="highlight">《{{ note.post.title }}》</span>
              </p>
            </template>
            <span class="note-time">{{ formatTime(note.created_at) }}</span>
      
          </div>


          <button v-if="note.type === 'signup'" class="check-btn" @click="showApplicantInfo(note)">
            查看用户
          </button>
          
          <button v-if="note.type === 'system'" class="check-btn primary" @click="showSystemDetail(note)">
            查看全文
          </button>
        </div>

        <div v-if="notifications.length === 0" class="empty-state">
           <div class="empty-emoji">🔕</div>
           <p>暂时没有新通知</p>
        </div>
      </div>
      <div v-if="showModal" class="modal-overlay" @click.self="showModal = false">
      </div>

      <div v-if="showSystemModal" class="modal-overlay" @click.self="showSystemModal = false">
        <div class="modal-card system-modal">
          <h3 class="modal-title">系统通知</h3>
          <div class="modal-body">
            <p>{{ currentSystemNote.content }}</p>
          </div>
          <div class="modal-footer">
            <span class="modal-time">{{ formatTime(currentSystemNote.created_at) }}</span>
            <button class="close-btn" @click="showSystemModal = false">关闭</button>
          </div>
        </div>
      </div>

      <div class="list-container" v-if="activeTab !== 'notifications'">
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
    <div v-if="showModal" class="modal-overlay" @click.self="showModal = false">
      <div class="modal-card">
        <h3>报名者信息</h3>
        <div class="applicant-info" v-if="currentApplicant.sender">
          <img :src="currentApplicant.sender.avatarurl" class="large-avatar">
          <p>姓名：{{ currentApplicant.sender.username }}</p>
          <p>学号：{{ currentApplicant.sender.studentid }}</p>
          <p>宿舍：{{ currentApplicant.sender.dorm?.dormname || '未知' }}</p>
        </div>
        <div class="modal-actions">
          <button class="contact-btn" @click="goToChat(currentApplicant.sender.id)">
            点击联系 (私信)
          </button>
          <button class="close-btn" @click="showModal = false">关闭</button>
        </div>
      </div>
    </div>
  </div>

  



</template>

<script setup>
import { ref, onMounted, computed , watch } from 'vue'
import { useRouter } from 'vue-router'
import Card from '@/components/Card.vue' // 复用你的卡片组件
import axios from 'axios'
import { useUserStore } from '@/stores/user';
import api from '@/api/index.ts';
const userStore = useUserStore(); 
const router = useRouter()
const activeTab = ref('posts') // posts 或 likes
const unreadCount = computed(() => userStore.unreadCount);
// 用户信息
const currentUser = computed(() => userStore.currentUser);

// 我的帖子列表
const myPosts = ref([])
// 我的收藏列表 
const favoritePosts = ref([])

const notifications = ref([]);
const showModal = ref(false);
const currentApplicant = ref({}); // 当前点击的那个报名通知对象
const showSystemModal = ref(false)
const currentSystemNote = ref({})

const showSystemDetail = (note) => {
  currentSystemNote.value = note
  showSystemModal.value = true
}
// 获取通知
const fetchNotifications = async () => {
    try {
        const res = await api.get('/api/v1/message/notifications');
        if(res.data.code === 200) {
            notifications.value = res.data.data;
        }
    } catch (err) {
        console.error("获取通知失败", err);
    }
}

const showApplicantInfo = (note) => {
    currentApplicant.value = note; // 这里存的是整个 notification 对象
    showModal.value = true;
}
const handleNoteClick = (note) => {
  // 如果是系统通知，点击可能不需要弹窗，或者只是标记已读
  if (note.type === 'system') {
    return; // 系统通知暂时不支持点击交互
  }
  
  // 如果是报名通知，显示详情
  if (note.type === 'signup') {
    currentApplicant.value = note;
    showModal.value = true;
  }
}
const goToChat = (targetUserId) => {
    router.push({ 
        name: 'Chat', 
        params: { id: targetUserId }
    });
}

watch(activeTab, async (newVal) => {
    if(newVal === 'notifications') {
      await fetchNotifications();
      await api.post('/api/v1/message/read_all');
      userStore.clearUnread();
    }else if (newVal === 'likes') { // 假设你把"我的收藏"tab的值设为 likes
       await fetchFavoritePosts();
    }
})

const formatTime = (timeStr) => {
  if (!timeStr) return ''
  const date = new Date(timeStr)
  return date.toLocaleString('zh-CN', {
    month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit'
  })
}

const goHome = () => {
  router.push('/dormgo')
}

const handleLogout = () => {
  userStore.logout();
  router.push('/');
}


const fetchMyPosts = async () => {
  // 这里应该调用 API 获取当前用户的帖子
  // 模拟数据:
  try {
     // 假设后端接口是 /api/v1/post/user_posts?uid=xxx
     // const res = await axios.get(...)
     
     // 暂时使用空数组或模拟数据
     myPosts.value = [
       {
         id: 101,
         title: '【转让】大四毕业出一辆捷安特山地车',
         content: '骑了两年，车况良好，海韵苑自取...',
         typeid: 1, // 交易
         Dorm: { dormname: '海韵苑' },
         User: { username: currentUser.value.username },
         created_at: '2023-10-20'
       }
     ]
  } catch (error) {
    console.error('获取个人帖子失败', error)
  }
}

// 计算当前显示的列表
const displayPosts = computed(() => {
  if (activeTab.value === 'posts') return myPosts.value
  if (activeTab.value === 'likes') return favoritePosts.value // 返回收藏数据
  return []
})

const fetchFavoritePosts = async () => {
    try {
        // 调用后端接口
        const res = await api.get('/api/v1/post/my_favorites');
        if (res.data.code === 200) {
            favoritePosts.value = res.data.data || [];
        }
    } catch (error) {
        console.error('获取收藏失败', error);
    }
}






onMounted(() => {
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

/* --- 新增：通知列表样式 --- */
.notification-list {
  background: white;
  min-height: 300px;
  padding: 0; /* 贴合边缘 */
}

.notification-item {
  display: flex;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid #f0f0f0;
  cursor: pointer;
  transition: background 0.3s;
}

.notification-item:hover {
  background: #f9f9f9;
}

.note-avatar {
  width: 44px;
  height: 44px;
  border-radius: 50%;
  object-fit: cover;
  margin-right: 12px;
  border: 1px solid #eee;
}

.note-content {
  flex: 1;
}

.note-text {
  font-size: 14px;
  color: #333;
  margin: 0 0 6px 0;
  line-height: 1.5;
}

.highlight {
  font-weight: 600;
  color: #1890ff;
}

.note-time {
  font-size: 12px;
  color: #999;
}

.check-btn {
  padding: 6px 12px;
  background: white;
  border: 1px solid #d9d9d9;
  border-radius: 4px;
  color: #666;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.3s;
}

.check-btn:hover {
  color: #1890ff;
  border-color: #1890ff;
}

/* 小红点 badge */
.badge {
  display: inline-block;
  width: 8px;
  height: 8px;
  background: #ff4d4f;
  border-radius: 50%;
  position: absolute;
  top: 12px;
  right: -5px;
}

/* --- 新增：弹窗样式 --- */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.6); /* 半透明遮罩 */
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
  animation: fadeIn 0.2s ease;
}

.modal-card {
  background: white;
  width: 90%;
  max-width: 400px;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  text-align: center;
  animation: scaleUp 0.2s ease;
}

.modal-card h3 {
  margin: 0 0 20px 0;
  font-size: 18px;
  color: #333;
}

.applicant-info {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  margin-bottom: 24px;
}

.applicant-info .large-avatar {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  object-fit: cover;
  border: 3px solid #f0f0f0;
  margin-bottom: 8px;
}

.applicant-info p {
  margin: 0;
  font-size: 15px;
  color: #555;
}

.modal-actions {
  display: flex;
  gap: 12px;
  justify-content: center;
}

.contact-btn {
  flex: 1;
  background: #1890ff;
  color: white;
  border: none;
  padding: 10px 0;
  border-radius: 6px;
  font-size: 14px;
  cursor: pointer;
}

.contact-btn:hover {
  background: #40a9ff;
}

.close-btn {
  flex: 1;
  background: #f5f5f5;
  color: #666;
  border: 1px solid #d9d9d9;
  padding: 10px 0;
  border-radius: 6px;
  font-size: 14px;
  cursor: pointer;
}

.close-btn:hover {
  background: #e6e6e6;
}

/* 动画 */
@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

@keyframes scaleUp {
  from { transform: scale(0.9); opacity: 0; }
  to { transform: scale(1); opacity: 1; }
}

.system-badge { color: #ff4d4f; font-weight: bold; margin-right: 5px; }

.system-modal { text-align: left; }
.modal-title { font-size: 18px; font-weight: bold; margin-bottom: 15px; border-bottom: 1px solid #eee; padding-bottom: 10px; color: #333; }
.modal-body { font-size: 15px; line-height: 1.6; color: #444; min-height: 80px; white-space: pre-wrap; }
.modal-footer { margin-top: 20px; display: flex; justify-content: space-between; align-items: center; border-top: 1px solid #eee; padding-top: 15px; }
.modal-time { font-size: 12px; color: #999; }
</style>