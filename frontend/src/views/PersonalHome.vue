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
          <img :src="currentUser?.avatarurl || defaultAvatar" alt="头像" class="large-avatar" />
          <div class="edit-avatar-badge" @click="goToPage('/profile/detail')">📷</div>
        </div>
        <h2 class="user-name">{{ currentUser?.username || '未登录' }}</h2>
        <p class="user-bio">{{ userBio || '暂无个性签名...' }}</p>
        
        <div class="user-stats">
          <div class="stat-item">
            <span class="stat-num">{{ myPosts.length }}</span>
            <span class="stat-label">发布</span>
          </div>
          <div class="stat-item">
            <span class="stat-num">{{ totalLikes }}</span>
            <span class="stat-label">获赞</span>
          </div>
          <div class="stat-item">
            <span class="stat-num">{{ favoritePosts.length }}</span>
            <span class="stat-label">收藏</span>
          </div>
        </div>
      </div>
    </div>

    <div class="function-grid">
      <div class="function-item" @click="goToPage('/profile/coins')">
        <div class="function-icon">🪙</div>
        <span class="function-text">寝友币</span>
      </div>
      <div class="function-item" @click="goToPage('/profile/detail')">
        <div class="function-icon">📌</div>
        <span class="function-text">详细资料</span>
      </div>
      <div class="function-item" @click="goToPage('/profile/rules')">
        <div class="function-icon">📖</div>
        <span class="function-text">社区规范</span>
      </div>
      <div class="function-item" @click="goToPage('/contact-us')">
        <div class="function-icon">📞</div>
        <span class="function-text">联系我们</span>
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

          <button v-if="note.type === 'signup'" class="check-btn" @click.stop="showApplicantInfo(note)">
            查看用户
          </button>
          
          <button v-if="note.type === 'system'" class="check-btn primary" @click.stop="showSystemDetail(note)">
            查看全文
          </button>
        </div>

        <div v-if="notifications.length === 0" class="empty-state">
           <div class="empty-emoji">🔕</div>
           <p>暂时没有新通知</p>
        </div>
      </div>

      <div v-if="showSystemModal" class="modal-overlay" @click.self="showSystemModal = false">
        <div class="modal-card system-modal">
          <h3 class="modal-title">系统通知</h3>
          <div class="modal-body">
            <p>{{ currentSystemNote?.content }}</p>
          </div>
          <div class="modal-footer">
            <span class="modal-time">{{ formatTime(currentSystemNote?.created_at) }}</span>
            <button class="close-btn" @click="showSystemModal = false">关闭</button>
          </div>
        </div>
      </div>

      <div class="list-container" v-if="activeTab !== 'notifications'">
        <div v-if="displayPosts.length > 0" class="post-list">
          <div 
            v-for="post in displayPosts" 
            :key="post.id" 
            class="post-wrapper"
          >
            <Card :post="post" class="profile-post-card" />
            
            <button 
              v-if="activeTab === 'posts'" 
              class="delete-btn" 
              @click.stop="handleDeletePost(post.id)"
            >
              删除
            </button>
          </div>
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
        <div class="applicant-info" v-if="currentApplicant?.sender">
          <img :src="currentApplicant.sender.avatarurl || defaultAvatar" class="large-avatar">
          <p>姓名：{{ currentApplicant.sender.username }}</p>
          <p>学号：{{ currentApplicant.sender.studentid || '未填写' }}</p>
          <p>宿舍：{{ currentApplicant.sender.dorm?.dormname || '未知' }}</p>
        </div>
        <div class="modal-actions">
          <button class="contact-btn" @click="goToChat(currentApplicant?.sender?.id)">
            点击联系 (私信)
          </button>
          <button class="close-btn" @click="showModal = false">关闭</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import { ElMessageBox, ElMessage } from 'element-plus'
import { useRouter } from 'vue-router'
import Card from '@/components/Card.vue' 
import { useUserStore } from '@/stores/user'
import api from '@/api/index'

// --- 1. 类型定义 (Interfaces) ---

// 宿舍类型
interface Dorm {
  dormname: string;
}

// 用户类型
interface User {
  id: number;
  username: string;
  avatarurl: string;
  studentid?: string;
  intro?: string; // 对应 bio
  dorm?: Dorm;
}

// 简化的帖子类型 (根据 Card 组件需求)
interface Post {
  id: number;
  title: string;
  content?: string;
  [key: string]: any; // 允许其他字段
}

// 通知类型
interface Notification {
  id: number;
  type: 'system' | 'signup';
  content: string;
  created_at: string;
  sender?: User;
  post?: { title: string };
  is_read?: boolean;
}

// 统计数据响应结构
interface UserStats {
  total_likes: number;
  bio: string;
}

// 通用 API 响应结构
interface ApiResponse<T> {
  code: number;
  data: T;
  msg: string;
}

// --- 2. 状态与初始化 ---

const router = useRouter()
const userStore = useUserStore()

const defaultAvatar = 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'

// 使用具体的类型定义 tab
type TabType = 'posts' | 'likes' | 'notifications';
const activeTab = ref<TabType>('posts')

const unreadCount = computed<number>(() => userStore.unreadCount || 0)
const currentUser = computed<User | null>(() => userStore.currentUser as User | null)
const currentUserId = computed<number>(() => userStore.currentUserId || 0)

// 用户数据
const userBio = ref<string>('')
const totalLikes = ref<number>(0)

// 列表数据 (指定为数组类型)
const myPosts = ref<Post[]>([])
const favoritePosts = ref<Post[]>([])
const notifications = ref<Notification[]>([])

// 弹窗状态与数据
const showModal = ref(false)
const showSystemModal = ref(false)

// currentApplicant 可能是 Notification，也可能是 null (或者使用 Partial<Notification> 初始化为空对象)
const currentApplicant = ref<Partial<Notification>>({}) 
const currentSystemNote = ref<Partial<Notification>>({})

// --- 3. 方法定义 (Functions) ---

const goToPage = (path: string) => {
  router.push(path)
}

const goHome = () => {
  router.push('/dormgo')
}

const handleLogout = () => {
  userStore.logout()
  router.push('/')
}

const formatTime = (timeStr?: string): string => {
  if (!timeStr) return ''
  const date = new Date(timeStr)
  return date.toLocaleString('zh-CN', {
    month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit'
  })
}

// 获取通知
const fetchNotifications = async () => {
  try {
    // 显式指定 API 返回类型
    const res = await api.get<ApiResponse<Notification[]>>('/api/v1/message/notifications')
    if (res.data.code === 200) {
      notifications.value = res.data.data
    }
  } catch (err) {
    console.error("获取通知失败", err)
  }
}

const handleDeletePost = async (postId: number) => {
  try {
    // 1. 确认弹窗
    await ElMessageBox.confirm(
      '确定要删除这条帖子吗？删除后无法恢复。',
      '提示',
      {
        confirmButtonText: '确定删除',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )

    // 2. 调用 API (根据后端路由 /api/v1/post/:id)
    const res = await api.delete(`/api/v1/post/${postId}`)

    if (res.data.code === 200) {
      ElMessage.success('删除成功')
      
      // 3. 更新本地列表（不用刷新页面）
      myPosts.value = myPosts.value.filter(p => p.id !== postId)
      
      // 可选：更新统计数据
      fetchUserStats()
    } else {
      ElMessage.error(res.data.msg || '删除失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error(error)
      ElMessage.error('操作失败')
    }
  }
}

// 显示报名者详情
const showApplicantInfo = (note: Notification) => {
  currentApplicant.value = note
  showModal.value = true
}

// 显示系统通知详情
const showSystemDetail = (note: Notification) => {
  currentSystemNote.value = note
  showSystemModal.value = true
}

// 处理点击通知条目
const handleNoteClick = (note: Notification) => {
  if (note.type === 'system') {
    // 系统通知可以在这里做逻辑，比如直接查看
    // showSystemDetail(note); 
    return
  }
  
  if (note.type === 'signup') {
    currentApplicant.value = note
    showModal.value = true
  }
}

const goToChat = (targetUserId?: number) => {
  if (!targetUserId) return
  router.push({ 
    name: 'Chat', 
    params: { id: targetUserId.toString() }
  })
}

// 获取我的帖子
const fetchMyPosts = async () => {
  try {
    const res = await api.get<ApiResponse<Post[]>>('/api/v1/post/user_posts', {
      params: { user_id: currentUserId.value }
    })
    
    if (res.data.code === 200) {
      myPosts.value = res.data.data || []
    }
  } catch (error) {
    console.error('获取个人帖子失败', error)
  }
}

// 获取收藏
const fetchFavoritePosts = async () => {
  try {
    const res = await api.get<ApiResponse<Post[]>>('/api/v1/post/my_favorites')
    if (res.data.code === 200) {
      console.log('my_favorites API 返回：', res.data.data)
      favoritePosts.value = res.data.data || []
    }
  } catch (error) {
    console.error('获取收藏失败', error)
  }
}

// 获取用户统计
const fetchUserStats = async () => {
  try {
    const res = await api.get<ApiResponse<UserStats>>('/api/v1/user/stats')
    if (res.data.code === 200) {
      totalLikes.value = res.data.data?.total_likes || 0
      userBio.value = res.data.data?.bio || ''
    }
  } catch (error) {
    console.error('获取用户统计失败', error)
    // 降级策略
    if (currentUser.value && currentUser.value.intro) {
      userBio.value = currentUser.value.intro
    }
  }
}

// 计算属性：根据 Tab 显示不同列表
const displayPosts = computed<Post[]>(() => {
  if (activeTab.value === 'posts') return myPosts.value
  if (activeTab.value === 'likes') return favoritePosts.value
  return []
})

// --- 4. 生命周期与监听 ---

watch(activeTab, async (newVal) => {
  if (newVal === 'notifications') {
    await fetchNotifications()
    // 标记已读
    await api.post('/api/v1/message/read_all')
    userStore.clearUnread() // 假设 store 有这个 action
  } else if (newVal === 'likes') {
    await fetchFavoritePosts()
  }
})

onMounted(() => {
  fetchMyPosts()
  fetchFavoritePosts()
  fetchUserStats()
})
</script>

<style scoped>
/* 样式保持不变 */
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
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
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
  border: 3px solid rgba(255, 255, 255, 0.8);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
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
  color: white;
  margin: 0 0 4px 0;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
}

.user-bio {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.85);
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
  color: white;
}

.stat-label {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.85);
  margin-top: 4px;
}

/* 功能入口区 */
.function-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  background: white;
  margin: 0 12px 12px;
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

.post-wrapper {
  position: relative; /* 为绝对定位的删除按钮提供参考 */
  margin-bottom: 16px;
}

/* 调整 Card 的 margin，因为现在由 wrapper 控制间距 */
.profile-post-card {
  margin-bottom: 0 !important; 
}

.delete-btn {
  position: absolute;
  top: 15px;
  right: 15px;
  background: white;
  border: 1px solid #ff4d4f;
  color: #ff4d4f;
  font-size: 12px;
  padding: 4px 10px;
  border-radius: 4px;
  cursor: pointer;
  z-index: 10; /* 确保在卡片上方 */
  transition: all 0.3s;
  opacity: 0.8;
}

.delete-btn:hover {
  background: #ff4d4f;
  color: white;
  opacity: 1;
}
</style>