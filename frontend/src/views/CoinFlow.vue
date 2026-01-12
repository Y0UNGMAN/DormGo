<template>
  <div class="page-container">
    <!-- 顶部导航 -->
    <div class="nav-bar">
      <button class="back-btn" @click="goBack">← 返回个人中心</button>
      <span class="page-title">寝友币</span>
      <span class="placeholder"></span>
    </div>

    <!-- 余额卡片 -->
    <div class="balance-card">
      <div class="balance-label">当前余额</div>
      <div class="balance-amount">
        <span class="coin-icon">🪙</span>
        <span class="amount">{{ balance }}</span>
      </div>
      <div class="balance-tip">寝友币可用于兑换社区权益</div>
    </div>

    <!-- 流水列表 -->
    <div class="content-area">
      <div class="section-header">
        <span class="section-title">流水记录</span>
        <el-select v-model="filterType" placeholder="全部" size="small" class="filter-select">
          <el-option label="全部" value="all" />
          <el-option label="收入" value="income" />
          <el-option label="支出" value="expense" />
        </el-select>
      </div>

      <div v-if="filteredRecords.length > 0" class="record-list">
        <div v-for="record in filteredRecords" :key="record.id" class="record-item">
          <div class="record-left">
            <div class="record-icon" :class="record.type">
              {{ record.type === 'income' ? '↑' : '↓' }}
            </div>
            <div class="record-info">
              <div class="record-title">{{ record.title }}</div>
              <div class="record-time">{{ formatTime(record.created_at) }}</div>
            </div>
          </div>
          <div class="record-amount" :class="record.type">
            {{ record.type === 'income' ? '+' : '-' }}{{ record.amount }}
          </div>
        </div>
      </div>

      <div v-else class="empty-state">
        <div class="empty-icon">🪙</div>
        <p class="empty-text">暂无流水记录</p>
      </div>
    </div>

    <!-- 获取寝友币说明 -->
    <div class="tips-card">
      <div class="tips-title">如何获取寝友币？</div>
      <ul class="tips-list">
        <li>📝 发布优质帖子 +5币</li>
        <li>💬 评论被点赞 +1币</li>
        <li>✅ 每日签到 +2币</li>
        <li>🎁 参与社区活动获得奖励</li>
      </ul>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import api from '@/api/index'

const router = useRouter()

interface CoinRecord {
  id: number
  type: 'income' | 'expense'
  title: string
  amount: number
  created_at: string
}

const balance = ref(0)
const filterType = ref('all')
const records = ref<CoinRecord[]>([])

const filteredRecords = computed(() => {
  if (filterType.value === 'all') return records.value
  return records.value.filter(r => r.type === filterType.value)
})

const goBack = () => {
  router.push('/personalhome')
}

const formatTime = (timeStr: string) => {
  if (!timeStr) return ''
  const date = new Date(timeStr)
  return date.toLocaleString('zh-CN', {
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const fetchCoinData = async () => {
  try {
    const res = await api.get('/api/v1/user/coins')
    if (res.data.code === 200) {
      balance.value = res.data.data?.balance || 0
      records.value = res.data.data?.records || []
    }
  } catch (error) {
    console.error('获取寝友币数据失败', error)
    // 使用模拟数据
    balance.value = 128
    records.value = [
      {
        id: 1,
        type: 'income',
        title: '发布帖子奖励',
        amount: 5,
        created_at: new Date().toISOString()
      },
      {
        id: 2,
        type: 'income',
        title: '新用户注册奖励',
        amount: 50,
        created_at: new Date(Date.now() - 86400000).toISOString()
      }
    ]
  }
}

onMounted(() => {
  fetchCoinData()
})
</script>

<style scoped>
.page-container {
  min-height: 100vh;
  background: #f5f5f5;
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

/* 余额卡片 */
.balance-card {
  background: linear-gradient(135deg, #f6d365 0%, #fda085 100%);
  margin: 12px;
  border-radius: 16px;
  padding: 24px;
  text-align: center;
  color: white;
  box-shadow: 0 4px 16px rgba(253, 160, 133, 0.4);
}

.balance-label {
  font-size: 14px;
  opacity: 0.9;
  margin-bottom: 8px;
}

.balance-amount {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  margin-bottom: 8px;
}

.coin-icon {
  font-size: 32px;
}

.amount {
  font-size: 42px;
  font-weight: 700;
}

.balance-tip {
  font-size: 12px;
  opacity: 0.8;
}

/* 内容区域 */
.content-area {
  margin: 0 12px;
  background: white;
  border-radius: 12px;
  padding: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.section-title {
  font-size: 16px;
  font-weight: 600;
  color: #333;
}

.filter-select {
  width: 100px;
}

.record-list {
  display: flex;
  flex-direction: column;
}

.record-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 0;
  border-bottom: 1px solid #f0f0f0;
}

.record-item:last-child {
  border-bottom: none;
}

.record-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.record-icon {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  font-weight: bold;
}

.record-icon.income {
  background: #f0fff0;
  color: #52c41a;
}

.record-icon.expense {
  background: #fff0f0;
  color: #ff4d4f;
}

.record-info {
  display: flex;
  flex-direction: column;
}

.record-title {
  font-size: 14px;
  color: #333;
  margin-bottom: 2px;
}

.record-time {
  font-size: 12px;
  color: #999;
}

.record-amount {
  font-size: 16px;
  font-weight: 600;
}

.record-amount.income {
  color: #52c41a;
}

.record-amount.expense {
  color: #ff4d4f;
}

/* 提示卡片 */
.tips-card {
  margin: 12px;
  background: white;
  border-radius: 12px;
  padding: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.tips-title {
  font-size: 15px;
  font-weight: 600;
  color: #333;
  margin-bottom: 12px;
}

.tips-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.tips-list li {
  font-size: 13px;
  color: #666;
  padding: 6px 0;
  border-bottom: 1px dashed #f0f0f0;
}

.tips-list li:last-child {
  border-bottom: none;
}

.empty-state {
  text-align: center;
  padding: 40px 20px;
}

.empty-icon {
  font-size: 40px;
  margin-bottom: 12px;
}

.empty-text {
  font-size: 14px;
  color: #999;
}
</style>
