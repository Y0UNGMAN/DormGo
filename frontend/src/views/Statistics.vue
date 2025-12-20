<template>
  <div class="app-container">
    <el-row :gutter="20" class="stat-cards">
      <el-col :span="6" v-for="(item, index) in statCards" :key="index">
        <el-card shadow="hover" :body-style="{ padding: '20px' }" class="box-card">
          <div class="stat-content">
            <div class="stat-text">
              <div class="stat-label">{{ item.label }}</div>
              <div class="stat-num">{{ item.value }}</div>
            </div>
            <el-icon class="stat-icon" :size="48" :color="item.color">
              <component :is="item.icon" />
            </el-icon>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px;">
      
      <el-col :span="12">
        <el-card shadow="never">
          <template #header>
            <div class="chart-header">
              <span>内容构成概览</span>
            </div>
          </template>
          <div id="pieChart" style="width: 100%; height: 400px;"></div>
        </el-card>
      </el-col>

      <el-col :span="12">
        <el-card shadow="never">
          <template #header>
            <div class="chart-header">
              <span>帖子板块分布</span>
            </div>
          </template>
          <div id="barChart" style="width: 100%; height: 400px;"></div>
        </el-card>
      </el-col>

    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, nextTick, onBeforeUnmount } from 'vue'
import request from '@/utils/request'
import { User, OfficeBuilding, DocumentCopy, Warning } from '@element-plus/icons-vue'
import * as echarts from 'echarts'

// 响应式数据
const stats = ref({
  totalUsers: 0, 
  totalDorms: 0, 
  totalContents: 0, 
  totalViolations: 0
})

const chartData = ref({
  pie: [],
  typeNames: [],
  typeValues: []
})

// 卡片配置
const statCards = computed(() => [
  { label: '注册用户', value: stats.value.totalUsers, icon: User, color: '#409EFF' },
  { label: '覆盖楼栋', value: stats.value.totalDorms, icon: OfficeBuilding, color: '#67C23A' },
  { label: '内容总数', value: stats.value.totalContents, icon: DocumentCopy, color: '#E6A23C' },
  { label: '违规记录', value: stats.value.totalViolations, icon: Warning, color: '#F56C6C' },
])

let pieChartInstance = null
let barChartInstance = null

onMounted(async () => {
  await fetchData()
  initCharts()
  window.addEventListener('resize', handleResize)
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', handleResize)
  if (pieChartInstance) pieChartInstance.dispose()
  if (barChartInstance) barChartInstance.dispose()
})

// 获取数据
const fetchData = async () => {
  try {
    const res = await request.get('/api/v1/admin/statistics')
    
    // 更新卡片数据
    stats.value = {
      totalUsers: res.total_users || 0,
      totalDorms: res.total_dorms || 0,
      totalContents: res.total_contents || 0,
      totalViolations: res.total_violations || 0
    }

    // 更新图表数据
    chartData.value = {
      pie: res.pie_data || [],
      typeNames: res.type_names || [],  // 后端返回的板块名称数组
      typeValues: res.type_values || [] // 后端返回的对应数量数组
    }
  } catch (e) { 
    console.error("获取统计数据失败:", e) 
  }
}

// 初始化图表
const initCharts = () => {
  nextTick(() => {
    // 1. 初始化饼图 (Pie Chart)
    const pieDom = document.getElementById('pieChart')
    if (pieDom) {
      pieChartInstance = echarts.init(pieDom)
      pieChartInstance.setOption({
        tooltip: { 
          trigger: 'item', 
          formatter: '{b}: {c} ({d}%)' 
        },
        legend: { 
          bottom: '5%', 
          left: 'center' 
        },
        color: ['#409EFF', '#67C23A', '#F56C6C'],
        series: [
          {
            name: '数据来源',
            type: 'pie',
            radius: ['40%', '70%'],
            avoidLabelOverlap: false,
            itemStyle: {
              borderRadius: 10,
              borderColor: '#fff',
              borderWidth: 2
            },
            label: { show: false, position: 'center' },
            emphasis: {
              label: { show: true, fontSize: 20, fontWeight: 'bold' }
            },
            data: chartData.value.pie
          }
        ]
      })
    }

    // 2. 初始化柱状图 (Bar Chart)
    const barDom = document.getElementById('barChart')
    if (barDom) {
      barChartInstance = echarts.init(barDom)
      barChartInstance.setOption({
        tooltip: { 
          trigger: 'axis', 
          axisPointer: { type: 'shadow' } 
        },
        grid: { 
          left: '3%', 
          right: '4%', 
          bottom: '3%', 
          containLabel: true 
        },
        xAxis: [
          {
            type: 'category',
            data: chartData.value.typeNames, // 使用后端修复后的数据
            axisTick: { alignWithLabel: true },
            axisLabel: { interval: 0, rotate: 30 } // 旋转防止文字重叠
          }
        ],
        yAxis: [{ 
          type: 'value',
          minInterval: 1 // 修改点：强制最小间隔为1，确保只显示整数
        }],
        series: [
          {
            name: '帖子数量',
            type: 'bar',
            barWidth: '40%',
            data: chartData.value.typeValues, // 使用后端修复后的数据
            itemStyle: {
              color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
                { offset: 0, color: '#83bff6' },
                { offset: 0.5, color: '#188df0' },
                { offset: 1, color: '#188df0' }
              ])
            }
          }
        ]
      })
    }
  })
}

// 监听窗口调整
const handleResize = () => {
  pieChartInstance?.resize()
  barChartInstance?.resize()
}
</script>

<style scoped>
.app-container {
  padding: 20px;
}

/* 统计卡片样式优化 */
.stat-cards .el-card {
  border: none;
  background-color: #ffffff;
  border-radius: 8px;
  transition: transform 0.2s;
}

.stat-cards .el-card:hover {
  transform: translateY(-5px);
}

.stat-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.stat-label {
  color: #909399;
  font-size: 14px;
  margin-bottom: 8px;
}

.stat-num {
  font-size: 28px;
  font-weight: 700;
  color: #303133;
}

.chart-header {
  font-weight: bold;
  color: #303133;
  padding-left: 10px;
  border-left: 4px solid #409EFF;
}
</style>