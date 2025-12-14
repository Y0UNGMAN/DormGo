<template>
  <div class="chat-page">
    <div class="chat-header">
      <button class="back-btn" @click="$router.back()">
        <span class="back-icon">←</span>
      </button>
      <div class="header-info">
        <span class="target-name">{{ targetUser.username || '加载中...' }}</span>
        <span class="target-status" v-if="targetUser.id">在线</span>
      </div>
      <button class="menu-btn">...</button>
    </div>

    <div class="message-list" ref="msgListRef">
      <div v-for="(msg, index) in messages" :key="msg.id" class="message-wrapper">
        
        <div v-if="shouldShowTime(msg, index)" class="message-time">
          {{ formatTime(msg.created_at) }}
        </div>

        <div :class="['message-row', isMe(msg.sender_id) ? 'me' : 'other']">
          <img 
            :src="isMe(msg.sender_id) ? currentUser.avatarurl : targetUser.avatarurl" 
            class="avatar"
            @click="goToProfile(msg.sender_id)"
          >
          
          <div class="bubble-container">
            <div class="bubble">
              {{ msg.content }}
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="input-area">
      <div class="input-tools">
        <button class="tool-btn">📷</button>
      </div>
      <div class="input-box-wrapper">
        <input 
          v-model="inputContent" 
          type="text" 
          placeholder="发送消息..." 
          @keyup.enter="handleSend"
          ref="inputRef"
        >
        <button 
          class="send-btn" 
          :disabled="!inputContent.trim()"
          @click="handleSend"
        >
          发送
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import api from '@/api/index.ts'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

// 状态数据
const targetId = parseInt(route.params.id) // 对方的用户ID
const targetUser = ref({}) // 对方的信息（头像、昵称）
const messages = ref([]) // 消息列表
const inputContent = ref('')
const msgListRef = ref(null) // 用于控制滚动
const inputRef = ref(null)
let pollingTimer = null // 轮询定时器

// 当前用户信息
const currentUser = computed(() => userStore.currentUser)
const currentUserId = computed(() => userStore.currentUserId)

// 判断消息是不是我发的
const isMe = (senderId) => {
  return senderId === currentUserId.value
}

// 1. 获取对方的基本信息 (头像、昵称)
const fetchTargetInfo = async () => {
  try {
    // 调用刚才后端补充的接口
    const res = await api.get(`/api/v1/user/info?user_id=${targetId}`)
    if (res.data.code === 200) {
      targetUser.value = res.data.data
    }
  } catch (err) {
    console.error("获取对方信息失败", err)
  }
}

// 2. 获取历史消息
const fetchHistory = async (isFirstLoad = false) => {
  try {
    const res = await api.get(`/api/v1/message/history`, {
      params: { target_id: targetId }
    })
    
    if (res.data.code === 200) {
      const newMessages = res.data.data || []
      
      // 简单判断：如果消息数量变了，或者最后一条消息ID变了，就更新
      // (为了体验更好，应该做去重合并，这里简单全量替换)
      if (newMessages.length !== messages.value.length) {
        messages.value = newMessages
        if (isFirstLoad) {
           scrollToBottom()
        } else {
           // 如果是轮询拉取，且用户在看底部，则自动滚；如果用户在看历史，就不滚
           // 这里简单处理：只要有新消息就滚到底部
           scrollToBottom()
        }
      }
    }
  } catch (err) {
    console.error("获取消息记录失败", err)
  }
}

// 3. 发送消息
const handleSend = async () => {
  const content = inputContent.value.trim()
  if (!content) return

  // 乐观更新：先在界面上显示出来，不用等接口返回，体验更快
  const tempMsg = {
    id: Date.now(), // 临时ID
    sender_id: currentUserId.value,
    receiver_id: targetId,
    content: content,
    created_at: new Date().toISOString()
  }
  console.log("乐观更新消息", tempMsg)
  messages.value.push(tempMsg)
  inputContent.value = ''
  scrollToBottom()

  try {
    const res = await api.post('/api/v1/message/send', {
      receiver_id: targetId,
      content: content
    })
    
    if (res.data.code !== 200) {
      alert("发送失败")
      // 实际项目中应该把刚才乐观更新的消息标红或删除
    } else {
        // 发送成功后，立即拉取一次最新状态（确保ID和时间同步）
        fetchHistory(false)
    }
  } catch (err) {
    console.error("发送报错", err)
  }
}

// 滚动到底部工具函数
const scrollToBottom = () => {
  nextTick(() => {
    if (msgListRef.value) {
      msgListRef.value.scrollTop = msgListRef.value.scrollHeight
    }
  })
}

// 时间格式化
const formatTime = (timeStr) => {
  const date = new Date(timeStr)
  return date.getHours().toString().padStart(2, '0') + ':' + 
         date.getMinutes().toString().padStart(2, '0')
}

// 简单的辅助函数：防止时间显示太密集
const shouldShowTime = (msg, index) => {
  if (index === 0) return true
  const prevTime = new Date(messages.value[index-1].created_at).getTime()
  const currTime = new Date(msg.created_at).getTime()
  // 间隔超过 5 分钟才显示时间
  return (currTime - prevTime) > 5 * 60 * 1000
}

const goToProfile = (uid) => {
    // router.push(`/user/${uid}`) // 预留功能
}

// 生命周期
onMounted(async () => {
  if (!currentUserId.value) {
      alert("请先登录")
      router.push('/')
      return
  }
  await fetchTargetInfo()
  await fetchHistory(true) // 初次加载，强制滚到底部
  
  // 开启轮询 (每3秒拉一次新消息)
  pollingTimer = setInterval(() => {
    fetchHistory(false)
  }, 2000)
  
  // 自动聚焦输入框
  inputRef.value?.focus()
})

onUnmounted(() => {
  if (pollingTimer) clearInterval(pollingTimer)
})
</script>

<style scoped>
.chat-page {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background-color: #f5f5f5;
}

/* 顶部导航 */
.chat-header {
  height: 50px;
  background: #fff;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 16px;
  border-bottom: 1px solid #e8e8e8;
  flex-shrink: 0;
}

.back-btn, .menu-btn {
  background: none;
  border: none;
  font-size: 18px;
  cursor: pointer;
  color: #333;
}

.header-info {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.target-name {
  font-size: 16px;
  font-weight: 600;
  color: #333;
}

.target-status {
  font-size: 10px;
  color: #52c41a;
}

/* 消息列表 */
.message-list {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 16px;
  scroll-behavior: smooth;
}

.message-time {
  text-align: center;
  font-size: 12px;
  color: #999;
  margin: 10px 0;
}

.message-row {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  max-width: 80%;
}

/* 别人的消息 (左侧) */
.message-row.other {
  align-self: flex-start;
}

/* 我的消息 (右侧) */
.message-row.me {
  align-self: flex-end;
  flex-direction: row-reverse; /* 关键：反转 flex 方向 */
}

.avatar {
  width: 40px;
  height: 40px;
  border-radius: 6px;
  object-fit: cover;
  flex-shrink: 0;
  cursor: pointer;
}

.bubble-container {
  display: flex;
  flex-direction: column;
}

.bubble {
  padding: 10px 14px;
  border-radius: 8px;
  font-size: 15px;
  line-height: 1.5;
  word-wrap: break-word;
  position: relative;
  box-shadow: 0 1px 2px rgba(0,0,0,0.05);
}

/* 气泡颜色区分 */
.other .bubble {
  background: #ffffff;
  color: #333;
  border-top-left-radius: 2px; /* 模拟微信气泡尖角方向 */
}

.me .bubble {
  background: #95ec69; /* 微信绿 */
  color: #000;
  border-top-right-radius: 2px;
}

/* 底部输入区 */
.input-area {
  background: #f7f7f7;
  border-top: 1px solid #e8e8e8;
  padding: 10px 16px;
  display: flex;
  align-items: center;
  gap: 10px;
  flex-shrink: 0;
}

.tool-btn {
  font-size: 24px;
  background: none;
  border: none;
  cursor: pointer;
  color: #666;
  padding: 0;
}

.input-box-wrapper {
  flex: 1;
  display: flex;
  gap: 8px;
  background: white;
  padding: 8px;
  border-radius: 6px;
}

input {
  flex: 1;
  border: none;
  outline: none;
  font-size: 15px;
  background: transparent;
}

.send-btn {
  background: #95ec69;
  color: #000; /* 微信风格是绿底黑字或白字 */
  border: none;
  padding: 4px 12px;
  border-radius: 4px;
  font-size: 13px;
  cursor: pointer;
  transition: opacity 0.2s;
}

.send-btn:disabled {
  background: #e0e0e0;
  color: #999;
  cursor: not-allowed;
}

/* 响应式适配 */
@media (max-width: 600px) {
  .message-row {
    max-width: 90%;
  }
}
</style>