<template>
  <div class="app-container">
    <el-tabs type="border-card">
      <el-tab-pane label="运营规则配置">
        <el-form label-width="160px" style="max-width: 600px; margin-top: 20px">
          <el-form-item label="配置保存">
            <el-button type="primary" @click="saveBasicConfig">保存配置</el-button>
          </el-form-item>
        </el-form>
      </el-tab-pane>

      <el-tab-pane label="内容安全配置">
        <el-form label-position="top">
          <el-form-item>
            <el-button type="danger" @click="saveSecurityConfig">更新安全配置</el-button>
          </el-form-item>
        </el-form>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup>
import { ElMessage } from 'element-plus'
import request from '@/utils/request'

const saveBasicConfig = async () => {
  try {
    await request.put('/api/v1/admin/config/basic', {})
    ElMessage.success('运营规则已更新')
  } catch (e) {}
}

const saveSecurityConfig = async () => {
  try {
    await request.put('/api/v1/admin/config/security', {})
    ElMessage.success('安全配置已更新')
  } catch (e) {}
}
</script>

<style scoped>
.app-container { padding: 20px; }
</style>