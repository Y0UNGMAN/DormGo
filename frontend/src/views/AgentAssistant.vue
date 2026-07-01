<template>
  <div class="agent-page">
    <div class="agent-bg">
      <div class="bg-shape shape-a"></div>
      <div class="bg-shape shape-b"></div>
    </div>

    <header class="agent-header">
      <button class="icon-btn" @click="router.back()" title="返回">←</button>
      <div class="title-block">
        <span class="eyebrow">DormGo Agent</span>
        <h1>校园任务助手</h1>
      </div>
      <div class="agent-state" :class="{ busy: isSending }">
        <span class="state-dot"></span>
        {{ isSending ? '处理中' : '在线' }}
      </div>
    </header>

    <main class="agent-shell">
      <aside class="agent-sidebar">
        <section class="side-section">
          <h2>常用任务</h2>
          <button
            v-for="prompt in quickPrompts"
            :key="prompt"
            class="prompt-chip"
            @click="usePrompt(prompt)"
          >
            {{ prompt }}
          </button>
        </section>

        <section class="side-section compact">
          <h2>执行边界</h2>
          <div class="guardrail">
            <span class="guardrail-icon">✓</span>
            <span>报名前需要二次确认</span>
          </div>
          <div class="guardrail">
            <span class="guardrail-icon">✓</span>
            <span>帖子检索由后端工具执行</span>
          </div>
          <div class="guardrail">
            <span class="guardrail-icon">✓</span>
            <span>上下文保存到当前账号</span>
          </div>
        </section>
      </aside>

      <section class="chat-panel">
        <div class="message-list" ref="messageListRef">
          <div
            v-for="item in messages"
            :key="item.id"
            class="message-row"
            :class="item.role"
          >
            <div class="avatar-mark">
              {{ item.role === 'user' ? '我' : 'AI' }}
            </div>
            <div class="bubble">
              <pre>{{ item.content }}</pre>
            </div>
          </div>
        </div>

        <form class="input-bar" @submit.prevent="sendMessage">
          <input
            v-model="input"
            :disabled="isSending"
            placeholder="输入你要找的帖子或要执行的任务"
          />
          <button type="submit" :disabled="isSending || !input.trim()">
            {{ isSending ? '处理中' : '发送' }}
          </button>
        </form>
      </section>
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, onMounted, watch, ref } from 'vue'
import { useRouter } from 'vue-router'
import api from '@/api/index'
import { useUserStore } from '@/stores/user'

type ChatMessage = {
  id: number
  role: 'user' | 'assistant'
  content: string
}

const router = useRouter()
const userStore = useUserStore()
const input = ref('')
const isSending = ref(false)
const messageListRef = ref<HTMLElement | null>(null)
const defaultMessages: ChatMessage[] = [
  {
    id: Date.now(),
    role: 'assistant',
    content: '你可以让我帮你查找可报名帖子。涉及报名时，我会先让你确认，再执行动作。'
  }
]
const messages = ref<ChatMessage[]>([...defaultMessages])
const storageKey = computed(() => `dormgo:agent:messages:${userStore.currentUserId || 'guest'}`)
const quickPrompts = [
  '请找一起学习的相关帖子',
  '有没有一起看电影的相关帖子？',
  '帮我找今天下午拼车去南校的帖子',
  '报名第一个'
]

const usePrompt = (prompt: string) => {
  input.value = prompt
}

const sendMessage = async () => {
  const content = input.value.trim()
  if (!content || isSending.value) return

  messages.value.push({ id: Date.now(), role: 'user', content })
  input.value = ''
  isSending.value = true
  await scrollToBottom()

  try {
    const res = await api.post('/api/v1/agent/chat', { message: content })
    if (res.data.code === 200) {
      messages.value.push({
        id: Date.now() + 1,
        role: 'assistant',
        content: res.data.data?.reply || 'Agent 没有返回内容'
      })
    } else {
      messages.value.push({
        id: Date.now() + 1,
        role: 'assistant',
        content: res.data.msg || 'Agent 请求失败'
      })
    }
  } catch (err) {
    console.error(err)
    messages.value.push({
      id: Date.now() + 1,
      role: 'assistant',
      content: 'Agent 服务暂时不可用，请确认 Go 后端和 Python agent-service 都已启动。'
    })
  } finally {
    isSending.value = false
    await scrollToBottom()
  }
}

const scrollToBottom = async () => {
  await nextTick()
  if (messageListRef.value) {
    messageListRef.value.scrollTop = messageListRef.value.scrollHeight
  }
}

const loadLocalHistory = () => {
  const raw = localStorage.getItem(storageKey.value)
  if (!raw) return
  try {
    const parsed = JSON.parse(raw)
    if (Array.isArray(parsed) && parsed.length > 0) {
      messages.value = parsed.slice(-40)
    }
  } catch (err) {
    console.warn('读取 Agent 本地聊天记录失败', err)
  }
}

watch(
  messages,
  (value) => {
    localStorage.setItem(storageKey.value, JSON.stringify(value.slice(-40)))
  },
  { deep: true }
)

onMounted(async () => {
  loadLocalHistory()
  await scrollToBottom()
})
</script>

<style scoped>
.agent-page {
  min-height: 100vh;
  background: #f8f9fc;
  color: #263043;
  padding: 22px;
  position: relative;
  overflow: hidden;
}

.agent-bg {
  position: absolute;
  inset: 0;
  height: 260px;
  background: linear-gradient(135deg, #7637b1 0%, #8ec5fc 100%);
}

.bg-shape {
  position: absolute;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.18);
}

.shape-a {
  width: 180px;
  height: 180px;
  right: 12%;
  top: 28px;
}

.shape-b {
  width: 112px;
  height: 112px;
  left: 8%;
  top: 86px;
}

.agent-header {
  position: relative;
  max-width: 1120px;
  margin: 0 auto 20px;
  display: flex;
  align-items: center;
  gap: 16px;
  color: white;
}

.icon-btn {
  width: 42px;
  height: 42px;
  border: 1px solid rgba(255, 255, 255, 0.44);
  background: rgba(255, 255, 255, 0.18);
  color: white;
  border-radius: 8px;
  cursor: pointer;
  font-size: 22px;
  backdrop-filter: blur(10px);
}

.title-block {
  flex: 1;
}

.eyebrow {
  display: block;
  font-size: 12px;
  letter-spacing: 0;
  opacity: 0.84;
  margin-bottom: 4px;
}

.agent-header h1 {
  margin: 0;
  font-size: 26px;
  font-weight: 700;
}

.agent-state {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  min-width: 78px;
  justify-content: center;
  padding: 8px 12px;
  background: rgba(255, 255, 255, 0.18);
  border: 1px solid rgba(255, 255, 255, 0.34);
  border-radius: 999px;
  font-size: 14px;
  backdrop-filter: blur(10px);
}

.state-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #46d39a;
}

.agent-state.busy .state-dot {
  background: #ffd166;
}

.agent-shell {
  position: relative;
  max-width: 1120px;
  height: calc(100vh - 118px);
  margin: 0 auto;
  display: grid;
  grid-template-columns: 280px minmax(0, 1fr);
  gap: 18px;
}

.agent-sidebar,
.chat-panel {
  background: rgba(255, 255, 255, 0.96);
  border: 1px solid rgba(226, 232, 240, 0.86);
  border-radius: 8px;
  box-shadow: 0 18px 42px rgba(54, 66, 96, 0.16);
}

.agent-sidebar {
  padding: 18px;
  height: fit-content;
}

.side-section + .side-section {
  margin-top: 22px;
}

.side-section h2 {
  margin: 0 0 12px;
  font-size: 15px;
  color: #38405f;
}

.prompt-chip {
  width: 100%;
  text-align: left;
  border: 1px solid #e1e5ef;
  background: #f7f8fc;
  color: #38405f;
  border-radius: 8px;
  padding: 10px 12px;
  cursor: pointer;
  line-height: 1.45;
  transition: all 0.2s ease;
  margin-bottom: 10px;
}

.prompt-chip:hover {
  border-color: #764ba2;
  color: #764ba2;
  background: #f7f0ff;
}

.guardrail {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #596276;
  font-size: 13px;
  line-height: 1.5;
  margin-bottom: 10px;
}

.guardrail-icon {
  width: 18px;
  height: 18px;
  flex: 0 0 18px;
  border-radius: 50%;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  background: #e8f7f0;
  color: #1a9b61;
  font-size: 12px;
}

.chat-panel {
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.message-list {
  flex: 1;
  overflow-y: auto;
  padding: 22px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.message-row {
  display: flex;
  align-items: flex-start;
  gap: 10px;
}

.message-row.user {
  justify-content: flex-end;
}

.message-row.user .avatar-mark {
  order: 2;
  background: #5b6ee1;
  color: white;
}

.avatar-mark {
  width: 34px;
  height: 34px;
  border-radius: 8px;
  flex: 0 0 34px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  background: #f0e8fb;
  color: #764ba2;
  font-size: 12px;
  font-weight: 700;
}

.bubble {
  max-width: min(650px, 82%);
  border-radius: 8px;
  padding: 12px 15px;
  background: #f4f6fb;
  border: 1px solid #e8ecf4;
  line-height: 1.6;
}

.message-row.user .bubble {
  background: #5b6ee1;
  border-color: #5b6ee1;
  color: white;
}

.bubble pre {
  margin: 0;
  white-space: pre-wrap;
  word-break: break-word;
  font-family: inherit;
}

.input-bar {
  border-top: 1px solid #e5e7eb;
  padding: 14px 16px;
  display: flex;
  gap: 10px;
  background: #fbfcff;
}

.input-bar input {
  flex: 1;
  height: 44px;
  border: 1px solid #d9dde7;
  border-radius: 8px;
  padding: 0 14px;
  font-size: 14px;
  background: white;
}

.input-bar input:focus {
  outline: none;
  border-color: #764ba2;
  box-shadow: 0 0 0 3px rgba(118, 75, 162, 0.12);
}

.input-bar button {
  width: 92px;
  border: none;
  border-radius: 8px;
  background: #764ba2;
  color: white;
  cursor: pointer;
  font-weight: 600;
}

.input-bar button:disabled {
  cursor: not-allowed;
  opacity: 0.55;
}

@media (max-width: 768px) {
  .agent-page {
    padding: 12px;
  }

  .agent-header {
    gap: 10px;
  }

  .agent-state {
    display: none;
  }

  .agent-shell {
    height: calc(100vh - 100px);
    grid-template-columns: 1fr;
  }

  .agent-sidebar {
    display: none;
  }
}
</style>
