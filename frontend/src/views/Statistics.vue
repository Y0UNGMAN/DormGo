<template>
  <div class="app-container">
    <el-row :gutter="20" class="stat-cards">
      <el-col :span="6" v-for="(item, index) in statCards" :key="index">
        <el-card shadow="hover" :body-style="{ padding: '20px' }">
          <div class="stat-content">
            <div class="stat-text">
              <div class="stat-label">{{ item.label }}</div>
              <div class="stat-num">{{ item.value }}</div>
            </div>
            <el-icon class="stat-icon" :size="40" :color="item.color"><component :is="item.icon" /></el-icon>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px;">
      <el-col :span="16">
        <el-card shadow="never">
          <template #header>
            <div class="chart-header">
              <span>数据趋势 (近7天)</span>
            </div>
          </template>
          <div id="lineChart" style="width: 100%; height: 350px;"></div>
        </el-card>
      </el-col>
      
      <el-col :span="8">
        <el-card shadow="never">
          <template #header>
            <div class="chart-header">
              <span>内容构成</span>
            </div>
          </template>
          <div id="pieChart" style="width: 100%; height: 350px;"></div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, nextTick } from 'vue'
import request from '@/utils/request'
import { User, OfficeBuilding, DocumentCopy, Warning } from '@element-plus/icons-vue'
import * as echarts from 'echarts' // 需要 npm install echarts

const stats = ref({
  totalUsers: 0, totalDorms: 0, totalContents: 0, totalViolations: 0
})
const chartData = ref({
  dates: [],
  users: [],
  posts: [],
  pie: []
})

const statCards = computed(() => [
  { label: '总用户数', value: stats.value.totalUsers, icon: User, color: '#409EFF' },
  { label: '总楼栋数', value: stats.value.totalDorms, icon: OfficeBuilding, color: '#67C23A' },
  { label: '内容总数', value: stats.value.totalContents, icon: DocumentCopy, color: '#E6A23C' },
  { label: '违规记录', value: stats.value.totalViolations, icon: Warning, color: '#F56C6C' },
])

onMounted(async () => {
  await fetchData()
  initCharts()
})

const fetchData = async () => {
  try {
    const res = await request.get('/api/v1/admin/statistics')
    stats.value = {
      totalUsers: res.total_users,
      totalDorms: res.total_dorms,
      totalContents: res.total_contents,
      totalViolations: res.total_violations
    }
    chartData.value = {
      dates: res.dates || [],
      users: res.users_trend_data || [],
      posts: res.posts_trend_data || [],
      pie: res.pie_data || []
    }
  } catch (e) { console.error(e) }
}

const initCharts = () => {
  nextTick(() => {
    // 1. 折线图
    const lineChart = echarts.init(document.getElementById('lineChart'))
    lineChart.setOption({
      tooltip: { trigger: 'axis' },
      legend: { data: ['活跃用户', '新增内容'] },
      xAxis: { type: 'category', data: chartData.value.dates },
      yAxis: { type: 'value' },
      series: [
        { name: '活跃用户', type: 'line', data: chartData.value.users, smooth: true, color: '#409EFF' },
        { name: '新增内容', type: 'line', data: chartData.value.posts, smooth: true, color: '#E6A23C' }
      ]
    })

    // 2. 饼图
    const pieChart = echarts.init(document.getElementById('pieChart'))
    pieChart.setOption({
      tooltip: { trigger: 'item' },
      legend: { bottom: '5%', left: 'center' },
      series: [
        {
          name: '内容分布',
          type: 'pie',
          radius: ['40%', '70%'],
          avoidLabelOverlap: false,
          itemStyle: { borderRadius: 10, borderColor: '#fff', borderWidth: 2 },
          label: { show: false, position: 'center' },
          emphasis: { label: { show: true, fontSize: 20, fontWeight: 'bold' } },
          data: chartData.value.pie
        }
      ]
    })

    // 响应式
    window.addEventListener('resize', () => {
      lineChart.resize()
      pieChart.resize()
    })
  })
}
</script>

<style scoped>
.app-container { padding: 20px; }
.stat-content { display: flex; justify-content: space-between; align-items: center; }
.stat-label { color: #909399; font-size: 14px; }
.stat-num { font-size: 24px; font-weight: bold; margin-top: 5px; }
.chart-header { font-weight: bold; color: #303133; }
</style>