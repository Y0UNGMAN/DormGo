<template>
  <div class="agent-page">
    <div class="agent-bg">
      <div class="bg-shape shape-a"></div>
      <div class="bg-shape shape-b"></div>
    </div>

    <header class="agent-header">
      <button class="icon-btn" @click="toggleSidebar" title="对话列表">☰</button>
      <button class="icon-btn" @click="router.back()" title="返回">←</button>
      <div class="title-block">
        <span class="eyebrow">DormGo Agent (ReAct)</span>
        <h1>校园任务助手</h1>
      </div>
      <div class="agent-state" :class="{ busy: isSending }">
        <span class="state-dot"></span>
        {{ isSending ? '推理中' : '在线' }}
      </div>
    </header>

    <main class="agent-shell">
      <!-- 对话列表侧边栏 -->
      <aside class="conv-sidebar" :class="{ open: sidebarOpen }">
        <div class="conv-sidebar-header">
          <span class="conv-sidebar-title">聊天记录</span>
          <button class="conv-sidebar-close" @click="toggleSidebar" title="关闭侧边栏">✕</button>
        </div>

        <button class="new-chat-btn" @click="createNewConversation">
          <span class="new-chat-icon">+</span>
          新对话
        </button>

        <div class="conv-list" ref="convListRef">
          <div
            v-for="conv in sortedConversations"
            :key="conv.id"
            class="conv-item"
            :class="{ active: conv.id === activeConvId }"
            @click="switchConversation(conv.id)"
          >
            <div class="conv-item-main">
              <span class="conv-title">{{ conv.title || '新对话' }}</span>
              <span class="conv-time">{{ formatConvTime(conv.updatedAt) }}</span>
            </div>
            <button
              class="conv-delete-btn"
              @click.stop="deleteConversation(conv.id)"
              title="删除对话"
            >
              🗑
            </button>
          </div>

          <div v-if="sortedConversations.length === 0" class="conv-empty">
            暂无对话记录
          </div>
        </div>

        <!-- 快捷工具区（折叠） -->
        <div class="conv-sidebar-footer">
          <button class="collapse-section-btn" @click="showTools = !showTools">
            常用任务 & 工具 {{ showTools ? '▾' : '▸' }}
          </button>
          <div v-if="showTools" class="tools-collapse-body">
            <div class="prompt-mini" v-for="prompt in quickPrompts" :key="prompt" @click="usePrompt(prompt)">
              {{ prompt }}
            </div>
            <div class="tools-mini-label">可用工具：search_posts, search_signup_posts, prepare_signup, confirm_signup</div>
          </div>
        </div>
      </aside>

      <!-- 侧边栏遮罩（移动端） -->
      <div v-if="sidebarOpen" class="sidebar-overlay" @click="toggleSidebar"></div>

      <section class="chat-panel">
        <div class="message-list" ref="messageListRef">
          <div
            v-for="item in displayMessages"
            :key="item.id"
            class="message-row"
            :class="item.role"
          >
            <div class="avatar-mark">
              {{ item.role === 'user' ? '我' : 'AI' }}
            </div>
            <div class="message-body">
              <div class="bubble">
                <pre>{{ item.content }}</pre>
              </div>
              <div v-if="item.steps && item.steps.length > 0" class="steps-wrapper">
                <button class="steps-toggle" @click="item._showSteps = !item._showSteps">
                  {{ item._showSteps ? '▾' : '▸' }} 推理过程 ({{ item.steps.length }} 步)
                </button>
                <div v-if="item._showSteps" class="steps-container">
                  <div
                    v-for="step in item.steps"
                    :key="step.step"
                    class="step-card"
                  >
                    <div class="step-header">
                      <span class="step-badge">Step {{ step.step }}</span>
                      <span v-if="step.action !== 'FINISH'" class="step-tool">
                        🔧 {{ step.action }}
                      </span>
                      <span v-else class="step-finish">&#10003; 完成</span>
                    </div>
                    <div class="step-body">
                      <div class="step-thought">
                        <span class="step-label">Action</span>
                        {{ summarizeStep(step) }}
                      </div>
                      <div v-if="step.observation" class="step-observation">
                        <span class="step-label">Observation</span>
                        <pre>{{ step.observation }}</pre>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div v-if="isSending && streamingSteps.length > 0" class="message-row assistant">
            <div class="avatar-mark">AI</div>
            <div class="message-body">
              <div class="bubble thinking">
                <span class="thinking-dot-pulse"></span> 推理中...
              </div>
              <div class="steps-container">
                <div
                  v-for="step in streamingSteps"
                  :key="step.step"
                  class="step-card live"
                >
                  <div class="step-header">
                    <span class="step-badge">Step {{ step.step }}</span>
                    <span v-if="step.action !== 'FINISH'" class="step-tool">
                      🔧 {{ step.action }}
                    </span>
                    <span v-else class="step-finish">&#10003; 完成</span>
                  </div>
                  <div class="step-body">
                    <div class="step-thought">
                      <span class="step-label">Action</span>
                      {{ summarizeStep(step) }}
                    </div>
                    <div v-if="step.observation" class="step-observation">
                      <span class="step-label">Observation</span>
                      <pre>{{ step.observation }}</pre>
                    </div>
                  </div>
                </div>
              </div>
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
            {{ isSending ? '推理中' : '发送' }}
          </button>
        </form>
      </section>
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, watch, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'

type AgentStep = {
  step: number
  thought: string
  action: string
  action_input?: Record<string, any>
  observation?: string
}

type ChatMessage = {
  id: number
  role: 'user' | 'assistant'
  content: string
  steps?: AgentStep[]
  traceId?: string
  _showSteps?: boolean
}

type Conversation = {
  id: string
  title: string
  messages: ChatMessage[]
  createdAt: number
  updatedAt: number
}

const router = useRouter()
const userStore = useUserStore()
const input = ref('')
const isSending = ref(false)
const messageListRef = ref<HTMLElement | null>(null)
const convListRef = ref<HTMLElement | null>(null)

// ---- 对话管理 ----
const conversations = ref<Conversation[]>([])
const activeConvId = ref<string | null>(null)
const sidebarOpen = ref(true)
const showTools = ref(false)

const convStorageKey = computed(() => `dormgo:agent:conversations:${userStore.currentUserId || 'guest'}`)
const oldStorageKey = computed(() => `dormgo:agent:messages:${userStore.currentUserId || 'guest'}`)

const defaultGreeting: ChatMessage = {
  id: Date.now(),
  role: 'assistant',
  content: '我使用 ReAct 推理框架，会先思考再行动，每一步都对你可见。'
}

const messages = ref<ChatMessage[]>([{ ...defaultGreeting, id: Date.now() }])
const streamingSteps = ref<AgentStep[]>([])
const currentTraceId = ref('')
const isPageLeaving = ref(false)
const activeAbortController = ref<AbortController | null>(null)
const interruptedMessageSaved = ref(false)

const quickPrompts = [
  '帮我找一起学习的相关帖子',
  '有没有一起看电影的相关帖子？',
  '帮我找今天下午拼车去南校的帖子',
  '报名第一个'
]

const sortedConversations = computed(() =>
  [...conversations.value].sort((a, b) => b.updatedAt - a.updatedAt)
)

const displayMessages = computed(() => messages.value)

const usePrompt = (prompt: string) => {
  input.value = prompt
}

// ---- 对话操作方法 ----
function generateId(): string {
  return Date.now().toString(36) + '-' + Math.random().toString(36).slice(2, 10)
}

function getActiveConversation(): Conversation | undefined {
  if (!activeConvId.value) return undefined
  return conversations.value.find(c => c.id === activeConvId.value)
}

function persistConversations() {
  if (!userStore.currentUserId) return
  localStorage.setItem(convStorageKey.value, JSON.stringify({
    conversations: conversations.value,
    activeConvId: activeConvId.value,
  }))
}

function loadConversations() {
  // 尝试迁移旧数据
  migrateOldData()

  const raw = localStorage.getItem(convStorageKey.value)
  if (!raw) return
  try {
    const parsed = JSON.parse(raw)
    if (Array.isArray(parsed.conversations)) {
      conversations.value = parsed.conversations
    }
    activeConvId.value = parsed.activeConvId || null
  } catch {
    console.warn('读取对话记录失败')
  }
}

function migrateOldData() {
  // 如果新格式已存在则跳过
  if (localStorage.getItem(convStorageKey.value)) return

  const oldRaw = localStorage.getItem(oldStorageKey.value)
  if (!oldRaw) return
  try {
    const oldMessages = JSON.parse(oldRaw)
    if (Array.isArray(oldMessages) && oldMessages.length > 0) {
      // 过滤掉默认问候语
      const meaningful = oldMessages.filter((m: ChatMessage) => {
        if (m.role === 'assistant' && m.content.includes('ReAct 推理框架')) return false
        return true
      })
      if (meaningful.length > 0) {
        const firstUser = meaningful.find((m: ChatMessage) => m.role === 'user')
        const title = firstUser ? firstUser.content.slice(0, 20) : '历史对话'
        conversations.value = [{
          id: generateId(),
          title,
          messages: meaningful,
          createdAt: Date.now() - 86400000,
          updatedAt: Date.now() - 86400000,
        }]
        persistConversations()
      }
    }
    // 迁移后删除旧数据
    localStorage.removeItem(oldStorageKey.value)
  } catch {
    // 忽略迁移错误
  }
}

function saveCurrentConversation() {
  const conv = getActiveConversation()
  if (!conv) return
  // 只保存有实际内容的对话
  const meaningful = messages.value.filter(m => {
    if (m.role === 'assistant' && m.content.includes('ReAct 推理框架') && messages.value.length > 1) return false
    return true
  })
  if (meaningful.length === 0) return
  conv.messages = meaningful
  conv.updatedAt = Date.now()
  // 自动标题：取第一条用户消息
  if (conv.title === '新对话' || !conv.title) {
    const firstUser = meaningful.find(m => m.role === 'user')
    if (firstUser) {
      conv.title = firstUser.content.slice(0, 20)
    }
  }
  persistConversations()
}

function createNewConversation() {
  saveCurrentConversation()
  const newConv: Conversation = {
    id: generateId(),
    title: '新对话',
    messages: [],
    createdAt: Date.now(),
    updatedAt: Date.now(),
  }
  conversations.value.push(newConv)
  activeConvId.value = newConv.id
  messages.value = [{ ...defaultGreeting, id: Date.now() }]
  streamingSteps.value = []
  currentTraceId.value = ''
  persistConversations()
  // 移动端关闭侧边栏
  if (window.innerWidth <= 768) {
    sidebarOpen.value = false
  }
}

function switchConversation(convId: string) {
  if (convId === activeConvId.value) return
  saveCurrentConversation()
  const conv = conversations.value.find(c => c.id === convId)
  if (!conv) return
  activeConvId.value = conv.id
  messages.value = conv.messages.length > 0
    ? [...conv.messages]
    : [{ ...defaultGreeting, id: Date.now() }]
  streamingSteps.value = []
  currentTraceId.value = ''
  persistConversations()
  // 移动端关闭侧边栏
  if (window.innerWidth <= 768) {
    sidebarOpen.value = false
  }
  nextTick(() => scrollToBottom())
}

function deleteConversation(convId: string) {
  const idx = conversations.value.findIndex(c => c.id === convId)
  if (idx === -1) return
  conversations.value.splice(idx, 1)

  if (activeConvId.value === convId) {
    // 切到最新对话，或新建
    if (conversations.value.length > 0) {
      const latest = conversations.value[0] // sorted by updatedAt desc
      if (!latest) return
      activeConvId.value = latest.id
      messages.value = latest.messages.length > 0
        ? [...latest.messages]
        : [{ ...defaultGreeting, id: Date.now() }]
    } else {
      activeConvId.value = null
      messages.value = [{ ...defaultGreeting, id: Date.now() }]
    }
  }
  streamingSteps.value = []
  currentTraceId.value = ''
  persistConversations()
}

function formatConvTime(ts: number): string {
  const now = Date.now()
  const diff = now - ts
  if (diff < 60_000) return '刚刚'
  if (diff < 3600_000) return `${Math.floor(diff / 60_000)}分钟前`
  if (diff < 86400_000) return `${Math.floor(diff / 3600_000)}小时前`
  const d = new Date(ts)
  return `${d.getMonth() + 1}/${d.getDate()}`
}

function toggleSidebar() {
  sidebarOpen.value = !sidebarOpen.value
}

function summarizeStep(step: AgentStep): string {
  const action = step.action
  const input = step.action_input || {}
  if (action === 'search_posts') return `正在搜索相关帖子：${input.query || '用户需求'}`
  if (action === 'search_signup_posts') return `正在搜索可报名活动：${input.query || '用户需求'}`
  if (action === 'get_post_detail') return `正在读取帖子详情：#${input.post_id || ''}`.trim()
  if (action === 'get_post_detail_from_candidate') return `正在读取第 ${input.index || 1} 个候选帖子的详情`
  if (action === 'prepare_signup') return `正在为帖子 #${input.post_id || ''} 创建待确认报名动作`.trim()
  if (action === 'prepare_signup_from_candidate') return `正在为第 ${input.index || 1} 个候选帖子创建待确认报名动作`
  if (action === 'confirm_signup') return '正在确认报名'
  if (action === 'FINISH') return '已生成最终回复'
  return `正在调用工具：${action || 'unknown'}`
}

function appendInterruptedMessage(steps: AgentStep[]) {
  if (interruptedMessageSaved.value) return
  if (steps.length === 0) return
  interruptedMessageSaved.value = true
  messages.value.push({
    id: Date.now() + 1,
    role: 'assistant',
    content: '本轮推理被页面刷新或离开中断，已保留完成的步骤。请重新发送上一条需求继续。',
    steps: [...steps],
    traceId: currentTraceId.value,
    _showSteps: true,
  })
  saveCurrentConversation()
}

function isAbortLikeError(err: unknown): boolean {
  if (isPageLeaving.value) return true
  if (!err || typeof err !== 'object') return false
  const maybe = err as { name?: string; message?: string }
  return maybe.name === 'AbortError' || /abort|cancel/i.test(maybe.message || '')
}

function handlePageUnload() {
  isPageLeaving.value = true
  appendInterruptedMessage(streamingSteps.value)
  activeAbortController.value?.abort()
}

// ---- SSE 通信（保持不变） ----
const getAuthHeaders = (): Record<string, string> => {
  const token = localStorage.getItem('token')
  return token ? { Authorization: `Bearer ${token}` } : {}
}

const sendMessage = async () => {
  const content = input.value.trim()
  if (!content || isSending.value) return

  // 如果没有活跃对话，自动创建
  if (!activeConvId.value || !conversations.value.find(c => c.id === activeConvId.value)) {
    const newConv: Conversation = {
      id: generateId(),
      title: content.slice(0, 20),
      messages: [],
      createdAt: Date.now(),
      updatedAt: Date.now(),
    }
    conversations.value.push(newConv)
    activeConvId.value = newConv.id
    messages.value = []
    persistConversations()
  }

  messages.value.push({ id: Date.now(), role: 'user', content })
  input.value = ''
  isSending.value = true
  streamingSteps.value = []
  currentTraceId.value = ''
  interruptedMessageSaved.value = false
  const controller = new AbortController()
  activeAbortController.value = controller
  await scrollToBottom()

  try {
    const resp = await fetch('http://127.0.0.1:8080/api/v1/agent/chat/stream', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        ...getAuthHeaders(),
      },
      body: JSON.stringify({ message: content }),
      signal: controller.signal,
    })

    if (!resp.ok) {
      throw new Error(`HTTP ${resp.status}`)
    }

    const reader = resp.body!.getReader()
    const decoder = new TextDecoder()
    let buffer = ''
    let finalReply = ''
    let finalRendered = false
    const steps: AgentStep[] = []

    const renderFinalReply = async (reply: string) => {
      if (!reply || finalRendered) return
      finalRendered = true
      messages.value.push({
        id: Date.now() + 1,
        role: 'assistant',
        content: reply,
        steps: [...steps],
        traceId: currentTraceId.value,
        _showSteps: false,
      })
      isSending.value = false
      streamingSteps.value = []
      currentTraceId.value = ''
      saveCurrentConversation()
      await nextTick()
      await scrollToBottom()
    }

    while (true) {
      const { done, value } = await reader.read()
      if (done) break

      buffer += decoder.decode(value, { stream: true })

      const lines = buffer.split('\n')
      buffer = lines.pop() || ''

      for (const line of lines) {
        const trimmed = line.trim()
        if (!trimmed || !trimmed.startsWith('data:')) continue

        const jsonStr = trimmed.slice(5).trim()
        try {
          const parsed = JSON.parse(jsonStr)

          if (parsed.start) {
            currentTraceId.value = parsed.start.trace_id || ''
          }
          if (parsed.step) {
            const step: AgentStep = {
              step: parsed.step.step,
              thought: parsed.step.thought,
              action: parsed.step.action,
              action_input: parsed.step.action_input,
              observation: parsed.step.observation,
            }
            steps.push(step)
            streamingSteps.value = [...steps]
            await nextTick()
            await scrollToBottom()
          }
          if (parsed.finish) {
            finalReply = parsed.finish.reply || ''
            await renderFinalReply(finalReply)
            controller.abort()
            break
          }
        } catch {
          continue
        }
      }
    }

    if (finalRendered) {
      return
    }

    if (finalReply) {
      await renderFinalReply(finalReply)
    } else {
      messages.value.push({
        id: Date.now() + 1,
        role: 'assistant',
        content: 'Agent 未返回有效结果',
        steps: [...steps],
        _showSteps: true,
      } as any)
    }
  } catch (err) {
    if (isAbortLikeError(err)) {
      return
    }
    console.error(err)
    messages.value.push({
      id: Date.now() + 1,
      role: 'assistant',
      content: 'Agent 服务暂时不可用，请确认 Go 和 Python 都已启动。',
    })
  } finally {
    if (activeAbortController.value === controller) {
      activeAbortController.value = null
    }
    isSending.value = false
    streamingSteps.value = []
    currentTraceId.value = ''
    await scrollToBottom()
  }
}

const scrollToBottom = async () => {
  await nextTick()
  if (messageListRef.value) {
    messageListRef.value.scrollTop = messageListRef.value.scrollHeight
  }
}

// ---- 持久化 ----
watch(
  messages,
  () => {
    saveCurrentConversation()
  },
  { deep: true }
)

onMounted(async () => {
  window.addEventListener('beforeunload', handlePageUnload)
  loadConversations()

  // 如果有活跃对话则恢复，否则创建新对话
  if (activeConvId.value) {
    const conv = conversations.value.find(c => c.id === activeConvId.value)
    if (conv && conv.messages.length > 0) {
      messages.value = [...conv.messages]
    } else if (conv) {
      messages.value = [{ ...defaultGreeting, id: Date.now() }]
    }
  }
  // 如果没有活跃对话，保持空白新对话状态
  // （不自动创建，等用户发消息时通过 sendMessage 创建）

  // 移动端默认折叠侧边栏
  if (window.innerWidth <= 768) {
    sidebarOpen.value = false
  }
  await scrollToBottom()
})

onBeforeUnmount(() => {
  window.removeEventListener('beforeunload', handlePageUnload)
  handlePageUnload()
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
  animation: pulse-dot 1s infinite;
}

@keyframes pulse-dot {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.3; }
}

.agent-shell {
  position: relative;
  max-width: 1120px;
  height: calc(100vh - 118px);
  margin: 0 auto;
  display: grid;
  grid-template-columns: 260px minmax(0, 1fr);
  gap: 18px;
}

/* ---- 对话侧边栏 ---- */
.conv-sidebar {
  background: rgba(255, 255, 255, 0.96);
  border: 1px solid rgba(226, 232, 240, 0.86);
  border-radius: 8px;
  box-shadow: 0 18px 42px rgba(54, 66, 96, 0.16);
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow: hidden;
  transition: transform 0.25s ease;
}

.conv-sidebar-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 16px 10px;
  border-bottom: 1px solid #e8ecf4;
}

.conv-sidebar-title {
  font-size: 15px;
  font-weight: 700;
  color: #38405f;
}

.conv-sidebar-close {
  display: none;
  width: 28px;
  height: 28px;
  border: none;
  background: none;
  color: #596276;
  cursor: pointer;
  border-radius: 6px;
  font-size: 16px;
}

.conv-sidebar-close:hover {
  background: #f0f1f6;
}

.new-chat-btn {
  margin: 12px 12px 8px;
  padding: 10px 14px;
  border: 1px dashed #764ba2;
  background: #faf6ff;
  color: #764ba2;
  border-radius: 8px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 8px;
  transition: all 0.2s ease;
}

.new-chat-btn:hover {
  background: #f0e8fb;
  border-color: #5b21b6;
}

.new-chat-icon {
  font-size: 18px;
  font-weight: 700;
  line-height: 1;
}

/* ---- 对话列表 ---- */
.conv-list {
  flex: 1;
  overflow-y: auto;
  padding: 4px 8px;
}

.conv-item {
  display: flex;
  align-items: center;
  padding: 10px 12px;
  border-radius: 8px;
  cursor: pointer;
  transition: background 0.15s ease;
  margin-bottom: 2px;
}

.conv-item:hover {
  background: #f4f6fb;
}

.conv-item.active {
  background: #f0e8fb;
  border: 1px solid rgba(118, 75, 162, 0.18);
}

.conv-item-main {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 3px;
}

.conv-title {
  font-size: 13px;
  font-weight: 600;
  color: #38405f;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.conv-item.active .conv-title {
  color: #5b21b6;
}

.conv-time {
  font-size: 11px;
  color: #8c93a8;
}

.conv-delete-btn {
  width: 26px;
  height: 26px;
  border: none;
  background: none;
  cursor: pointer;
  border-radius: 6px;
  font-size: 12px;
  opacity: 0;
  transition: opacity 0.15s ease;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
}

.conv-item:hover .conv-delete-btn {
  opacity: 0.6;
}

.conv-delete-btn:hover {
  opacity: 1 !important;
  background: #fde8e8;
}

.conv-empty {
  padding: 24px 16px;
  text-align: center;
  color: #8c93a8;
  font-size: 13px;
}

/* ---- 侧边栏底部工具区 ---- */
.conv-sidebar-footer {
  border-top: 1px solid #e8ecf4;
  padding: 8px;
}

.collapse-section-btn {
  width: 100%;
  padding: 8px 10px;
  border: none;
  background: none;
  color: #596276;
  cursor: pointer;
  font-size: 12px;
  font-weight: 600;
  text-align: left;
  border-radius: 6px;
}

.collapse-section-btn:hover {
  background: #f4f6fb;
}

.tools-collapse-body {
  padding: 6px 4px;
}

.prompt-mini {
  padding: 5px 8px;
  font-size: 11px;
  color: #596276;
  cursor: pointer;
  border-radius: 4px;
  line-height: 1.4;
}

.prompt-mini:hover {
  background: #f0e8fb;
  color: #764ba2;
}

.tools-mini-label {
  padding: 5px 8px;
  font-size: 10px;
  color: #a8afc2;
  margin-top: 4px;
}

/* ---- 侧边栏遮罩 ---- */
.sidebar-overlay {
  display: none;
}

/* ---- 聊天面板 ---- */
.chat-panel {
  background: rgba(255, 255, 255, 0.96);
  border: 1px solid rgba(226, 232, 240, 0.86);
  border-radius: 8px;
  box-shadow: 0 18px 42px rgba(54, 66, 96, 0.16);
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

.message-body {
  flex: 1;
  min-width: 0;
}

.bubble {
  max-width: min(650px, 82%);
  border-radius: 8px;
  padding: 12px 15px;
  background: #f4f6fb;
  border: 1px solid #e8ecf4;
  line-height: 1.6;
}

.bubble.thinking {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #764ba2;
  font-style: italic;
}

.thinking-dot-pulse {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #764ba2;
  animation: pulse-dot 0.8s infinite;
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

.steps-container {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-top: 10px;
}

.steps-wrapper {
  margin-top: 8px;
}

.steps-toggle {
  background: none;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  padding: 5px 10px;
  font-size: 12px;
  color: #596276;
  cursor: pointer;
  transition: all 0.15s;
}

.steps-toggle:hover {
  background: #f4f6fb;
  color: #764ba2;
  border-color: #d4bff9;
}

.step-card {
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  background: #fafbfd;
  overflow: hidden;
  font-size: 13px;
}

.step-card.live {
  border-color: #764ba2;
  background: #faf6ff;
  box-shadow: 0 0 0 1px rgba(118, 75, 162, 0.15);
}

.step-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  background: #f1f5f9;
  border-bottom: 1px solid #e2e8f0;
}

.step-card.live .step-header {
  background: #f0e8fb;
  border-color: #d4bff9;
}

.step-badge {
  font-weight: 700;
  color: #475569;
  font-size: 11px;
  background: #e2e8f0;
  padding: 2px 7px;
  border-radius: 4px;
}

.step-card.live .step-badge {
  background: #d4bff9;
  color: #5b21b6;
}

.step-tool {
  color: #5b6ee1;
  font-weight: 600;
  font-family: monospace;
}

.step-finish {
  color: #1a9b61;
  font-weight: 700;
}

.step-body {
  padding: 10px 12px;
}

.step-thought {
  color: #38405f;
  line-height: 1.55;
  margin-bottom: 6px;
}

.step-thought .step-label {
  display: inline-block;
  font-weight: 700;
  color: #764ba2;
  margin-right: 6px;
  font-size: 11px;
  text-transform: uppercase;
}

.step-observation .step-label {
  display: inline-block;
  font-weight: 700;
  color: #1a9b61;
  margin-right: 6px;
  font-size: 11px;
  text-transform: uppercase;
}

.step-observation {
  color: #596276;
  line-height: 1.55;
  margin-top: 6px;
  padding-top: 6px;
  border-top: 1px dashed #e2e8f0;
}

.step-observation pre {
  margin: 2px 0 0;
  white-space: pre-wrap;
  word-break: break-word;
  font-family: inherit;
  font-size: 12px;
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

/* ---- 响应式 ---- */
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

  .conv-sidebar {
    position: fixed;
    left: 0;
    top: 0;
    bottom: 0;
    width: 280px;
    z-index: 100;
    border-radius: 0;
    transform: translateX(-100%);
  }

  .conv-sidebar.open {
    transform: translateX(0);
  }

  .conv-sidebar-close {
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .sidebar-overlay {
    display: block;
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.35);
    z-index: 99;
  }
}
</style>
