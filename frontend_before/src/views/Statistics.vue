<template>
  <div class="app-container">
    <el-row :gutter="20" class="data-panel">
      <el-col :xs="12" :sm="12" :lg="6" v-for="(item, index) in statCards" :key="index">
        <el-card shadow="hover" :body-style="{ padding: '20px' }">
          <div class="stat-card-content">
            <div class="stat-icon" :class="item.class">
              <el-icon><component :is="item.icon" /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-label">{{ item.label }}</div>
              <div class="stat-num">{{ item.value }}</div>
              <div class="stat-trend" :class="item.trend >= 0 ? 'up' : 'down'">
                同比 {{ Math.abs(item.trend) }}% 
                <span>{{ item.trend >= 0 ? '↑' : '↓' }}</span>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px;">
      <el-col :span="12" :xs="24">
        <el-card shadow="hover" class="chart-card">
          <template #header><span>用户增长趋势</span></template>
          <div class="chart-wrapper">
            <canvas id="userGrowthChart"></canvas>
          </div>
        </el-card>
      </el-col>
      <el-col :span="12" :xs="24">
        <el-card shadow="hover" class="chart-card">
          <template #header><span>内容发布统计</span></template>
          <div class="chart-wrapper">
            <canvas id="contentChart"></canvas>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row style="margin-top: 20px;">
      <el-col :span="24">
        <el-card shadow="hover">
          <template #header><span>楼栋住户分布</span></template>
          <el-table :data="dormUserDistribution" border stripe>
            <el-table-column prop="dorm_name" label="楼栋名称" width="180" />
            <el-table-column prop="user_count" label="住户数量" width="120" sortable />
            <el-table-column label="占比情况">
              <template #default="{ row }">
                <el-progress 
                  :percentage="row.percentage" 
                  :stroke-width="15" 
                  :color="customColorMethod"
                />
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, onUnmounted } from 'vue'
import axios from 'axios'
import Chart from 'chart.js/auto'
import { ElMessage } from 'element-plus'
import { User, OfficeBuilding, DocumentCopy, Warning } from '@element-plus/icons-vue'

// 数据源
const stats = ref({
  totalUsers: 0, userTrend: 0,
  totalDorms: 0, dormTrend: 0,
  totalContents: 0, contentTrend: 0,
  totalViolations: 0, violationTrend: 0
})
const dormUserDistribution = ref([])
let userChartInst = null
let contentChartInst = null

// 计算属性生成卡片配置
const statCards = computed(() => [
  { label: '总用户数', value: stats.value.totalUsers, trend: stats.value.userTrend, icon: User, class: 'icon-blue' },
  { label: '总楼栋数', value: stats.value.totalDorms, trend: stats.value.dormTrend, icon: OfficeBuilding, class: 'icon-green' },
  { label: '内容总数', value: stats.value.totalContents, trend: stats.value.contentTrend, icon: DocumentCopy, class: 'icon-purple' },
  { label: '违规记录', value: stats.value.totalViolations, trend: stats.value.violationTrend, icon: Warning, class: 'icon-red' },
])

const customColorMethod = (percentage) => {
  if (percentage < 30) return '#909399'
  if (percentage < 70) return '#e6a23c'
  return '#67c23a'
}

onMounted(() => {
  fetchData()
})

onUnmounted(() => {
  if (userChartInst) userChartInst.destroy()
  if (contentChartInst) contentChartInst.destroy()
})

const fetchData = async () => {
  try {
    const token = localStorage.getItem('adminToken')
    const res = await axios.get('/api/v1/admin/statistics', {
      headers: { Authorization: `Bearer ${token}` }
    })
    const data = res.data
    
    // 填充数据
    stats.value = {
      totalUsers: data.total_users, userTrend: data.user_trend,
      totalDorms: data.total_dorms, dormTrend: data.dorm_trend,
      totalContents: data.total_contents, contentTrend: data.content_trend,
      totalViolations: data.total_violations, violationTrend: data.violation_trend
    }
    dormUserDistribution.value = data.dorm_user_distribution
    
    // 初始化图表
    initCharts(data)
  } catch (e) {
    ElMessage.error('数据加载失败')
  }
}

const initCharts = (data) => {
  // 用户图表
  const ctxUser = document.getElementById('userGrowthChart')
  userChartInst = new Chart(ctxUser, {
    type: 'line',
    data: {
      labels: data.user_growth_dates,
      datasets: [{
        label: '新增用户',
        data: data.user_growth_data,
        borderColor: '#409EFF',
        backgroundColor: 'rgba(64, 158, 255, 0.1)',
        fill: true,
        tension: 0.4
      }]
    },
    options: { responsive: true, maintainAspectRatio: false }
  })

  // 内容图表
  const ctxContent = document.getElementById('contentChart')
  contentChartInst = new Chart(ctxContent, {
    type: 'bar',
    data: {
      labels: data.content_dates,
      datasets: [
        { label: '帖子', data: data.post_data, backgroundColor: '#409EFF' },
        { label: '评论', data: data.comment_data, backgroundColor: '#67C23A' }
      ]
    },
    options: { responsive: true, maintainAspectRatio: false, scales: { x: { stacked: true }, y: { stacked: true } } }
  })
}
</script>

<style scoped>
.app-container { padding: 20px; }
.stat-card-content { display: flex; align-items: center; }
.stat-icon { 
  font-size: 48px; margin-right: 20px; padding: 10px; border-radius: 8px;
}
.stat-info { flex: 1; }
.stat-label { color: #909399; font-size: 14px; }
.stat-num { font-size: 24px; font-weight: bold; color: #303133; margin: 5px 0; }
.stat-trend { font-size: 12px; display: flex; align-items: center; }
.stat-trend.up { color: #67c23a; }
.stat-trend.down { color: #f56c6c; }

/* 颜色类 */
.icon-blue { color: #409EFF; background: #ecf5ff; }
.icon-green { color: #67c23a; background: #f0f9eb; }
.icon-purple { color: #a0cfff; background: #ecf5ff; }
.icon-red { color: #f56c6c; background: #fef0f0; }

.chart-wrapper { height: 300px; position: relative; }
</style>