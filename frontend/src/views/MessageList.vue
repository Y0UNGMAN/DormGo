<template>
  <div class="msg-list-page">
    <div class="header">
      <button class="back-btn" @click="$router.push('/dormgo')">←</button>
      <span class="title">消息列表</span>
    </div>

    <div class="list-container">
      <div 
        v-for="item in conversations" 
        :key="item.target_id" 
        class="chat-item"
        @click="goToChat(item.target_id)"
      >
        <div class="avatar-box">
          <img :src="item.avatar_url" class="avatar" />
          <div v-if="item.unread_count > 0" class="badge">
            {{ item.unread_count > 99 ? '99+' : item.unread_count }}
          </div>
        </div>

        <div class="content-box">
          <div class="row-top">
            <span class="username">{{ item.username }}</span>
            <span class="time">{{ formatTime(item.last_time) }}</span>
          </div>
          <div class="row-bottom">
            <span class="last-msg">{{ item.last_content }}</span>
          </div>
        </div>
      </div>

      <div v-if="conversations.length === 0" class="empty">
        暂无消息
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import api from '@/api/index.js'

const router = useRouter()
const conversations = ref([])

// 获取会话列表
const fetchConversations = async () => {
  try {
    const res = await api.get('/api/v1/message/conversations')
    if (res.data.code === 200) {
      conversations.value = res.data.data
    }
  } catch (err) {
    console.error(err)
  }
}

// 跳转到私信页
const goToChat = (targetId) => {
  router.push(`/chat/${targetId}`)
}

// 时间格式化 (简单版)
const formatTime = (timeStr) => {
  const date = new Date(timeStr)
  const now = new Date()
  // 如果是今天，显示时间；如果是以前，显示日期
  if (date.toDateString() === now.toDateString()) {
    return `${date.getHours().toString().padStart(2,'0')}:${date.getMinutes().toString().padStart(2,'0')}`
  }
  return `${date.getMonth()+1}-${date.getDate()}`
}

onMounted(() => {
  fetchConversations()
})
</script>

<style scoped>
.msg-list-page {
  background: #fff;
  min-height: 100vh;
}
.header {
  height: 50px;
  border-bottom: 1px solid #f0f0f0;
  display: flex;
  align-items: center;
  padding: 0 15px;
  background: #f9f9f9;
}
.back-btn { border:none; background:none; font-size:20px; cursor:pointer; margin-right:10px;}
.title { font-weight:600; font-size:16px;}

.chat-item {
  display: flex;
  padding: 12px 15px;
  border-bottom: 1px solid #f5f5f5;
  cursor: pointer;
  transition: background 0.2s;
}
.chat-item:hover { background: #fafafa; }

.avatar-box { position: relative; margin-right: 12px; }
.avatar { width: 48px; height: 48px; border-radius: 6px; object-fit: cover; }
.badge {
  position: absolute; top: -6px; right: -6px;
  background: #ff4d4f; color: white;
  font-size: 10px; height: 16px; min-width: 16px;
  border-radius: 8px; text-align: center; line-height: 16px;
  padding: 0 4px;
}

.content-box { flex: 1; display: flex; flex-direction: column; justify-content: space-between; padding: 2px 0; }
.row-top { display: flex; justify-content: space-between; align-items: center; }
.username { font-weight: 500; font-size: 16px; color: #333; }
.time { color: #999; font-size: 12px; }
.last-msg { color: #999; font-size: 13px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; max-width: 250px; }
.empty { text-align: center; margin-top: 50px; color: #999; }
</style>