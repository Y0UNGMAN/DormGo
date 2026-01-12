<template>
  <div class="app-container">
    <el-tabs type="border-card">
      <el-tab-pane label="运营规则配置">
        <el-form :model="basicConfig" label-width="160px" style="max-width: 600px; margin-top: 20px">
          <el-divider content-position="left">互助流程参数</el-divider>
          
          <el-form-item label="自动超时时间 (小时)">
             <el-input-number v-model="basicConfig.postTimeout" :min="12" :max="168" />
             <div class="tips">接单后未确认完成，超过该时间自动标记待评价 [cite: 66]</div>
          </el-form-item>

          <el-divider content-position="left">信用分体系 (SRS 5.3)</el-divider>
          
          <el-form-item label="初始信用分">
             <el-input-number v-model="basicConfig.initialScore" :precision="1" :step="0.1" :max="10" disabled />
             <div class="tips">默认固定为 5.0 [cite: 141]</div>
          </el-form-item>
          
          <el-form-item label="5星好评奖励">
             <el-input-number v-model="basicConfig.reward5Star" :precision="2" :step="0.05" />
             <div class="tips">每次获得5星好评增加的分数 (默认 +0.1)</div>
          </el-form-item>
          
          <el-form-item label="1星差评扣除">
             <el-input-number v-model="basicConfig.penalty1Star" :precision="2" :step="0.05" />
             <div class="tips">每次获得1星差评扣除的分数 (默认 -0.3)</div>
          </el-form-item>
          
          <el-form-item label="无故取消扣除">
             <el-input-number v-model="basicConfig.penaltyCancel" :precision="2" :step="0.05" />
             <div class="tips">无正当理由取消互助扣除分数 (默认 -0.1)</div>
          </el-form-item>

          <el-form-item>
            <el-button type="primary" @click="saveBasicConfig" :loading="saving">保存配置</el-button>
          </el-form-item>
        </el-form>
      </el-tab-pane>

      <el-tab-pane label="内容安全配置">
        <el-alert title="即使V1.0不做实时审核，也需配置敏感词库以供预过滤 " type="warning" :closable="false" show-icon style="margin-bottom: 20px" />
        
        <el-form label-position="top">
          <el-form-item label="敏感词黑名单 (用英文逗号分隔)">
            <el-input 
              v-model="securityConfig.sensitiveWords" 
              type="textarea" 
              rows="10" 
              placeholder="例如: 涉黄,暴力,作弊,代考..."
            />
          </el-form-item>
          <el-form-item>
            <el-button type="danger" @click="saveSecurityConfig" :loading="saving">更新词库</el-button>
          </el-form-item>
        </el-form>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import request from '@/utils/request'

const saving = ref(false)

// 对应 SRS 5.3 信用分计算逻辑 [cite: 140-145]
const basicConfig = reactive({
  postTimeout: 72,
  initialScore: 5.0,
  reward5Star: 0.1,
  penalty1Star: 0.3,
  penaltyCancel: 0.1
})

const securityConfig = reactive({
  sensitiveWords: ''
})

onMounted(() => {
  // 模拟从后端获取配置
  // request.get('/v1/admin/config')...
  securityConfig.sensitiveWords = '作弊,代考,刷单,赌博'
})

const saveBasicConfig = async () => {
  saving.value = true
  try {
    await request.put('/v1/admin/config/basic', basicConfig)
    ElMessage.success('运营规则已更新')
  } catch (e) {
    // error handled
  } finally {
    saving.value = false
  }
}

const saveSecurityConfig = async () => {
  saving.value = true
  try {
    // 简单处理格式
    const words = securityConfig.sensitiveWords.split(',').map(s => s.trim()).filter(s => s)
    await request.put('/v1/admin/config/security', { words })
    ElMessage.success(`敏感词库已更新，共 ${words.length} 个词条`)
  } catch (e) {
  } finally {
    saving.value = false
  }
}
</script>

<style scoped>
.app-container { padding: 20px; }
.tips { font-size: 12px; color: #909399; margin-top: 5px; }
</style>