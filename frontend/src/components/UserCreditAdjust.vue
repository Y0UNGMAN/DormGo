<template>
  <el-dialog
    v-model="visible"
    title="调整用户信用分"
    width="450px"
    :before-close="handleClose"
  >
    <el-form :model="form" ref="formRef" :rules="rules" label-position="top">
      <el-alert 
        title="信用分将影响用户发布和接单权限，请谨慎操作" 
        type="warning" 
        :closable="false" 
        style="margin-bottom: 20px"
      />
      
      <el-form-item label="当前用户">
        <el-tag>{{ userInfo?.nickname }} ({{ userInfo?.student_id }})</el-tag>
        <span class="current-score">当前分数: {{ userInfo?.credit_score }}</span>
      </el-form-item>

      <el-form-item label="调整方式" prop="type">
        <el-radio-group v-model="form.type">
          <el-radio label="add">增加 (+)</el-radio>
          <el-radio label="deduct">扣除 (-)</el-radio>
        </el-radio-group>
      </el-form-item>

      <el-form-item label="分值" prop="amount">
        <el-input-number v-model="form.amount" :min="0.1" :max="10" :step="0.1" :precision="1" />
      </el-form-item>

      <el-form-item label="调整原因" prop="reason">
        <el-input 
          v-model="form.reason" 
          type="textarea" 
          placeholder="请输入调整原因，将通知用户" 
          rows="3"
        />
      </el-form-item>
    </el-form>

    <template #footer>
      <el-button @click="handleClose">取消</el-button>
      <el-button type="primary" @click="submit" :loading="loading">确认调整</el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import request from '@/utils/request'
import { ElMessage } from 'element-plus'

const props = defineProps({
  modelValue: Boolean,
  userInfo: Object
})

const emit = defineEmits(['update:modelValue', 'success'])

const visible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const formRef = ref(null)
const loading = ref(false)
const form = reactive({
  type: 'deduct',
  amount: 1.0,
  reason: ''
})

const rules = {
  amount: [{ required: true, message: '请输入分值', trigger: 'blur' }],
  reason: [{ required: true, message: '请输入原因', trigger: 'blur' }]
}

const handleClose = () => {
  visible.value = false
  form.type = 'deduct'
  form.amount = 1.0
  form.reason = ''
}

const submit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        const finalAmount = form.type === 'add' ? form.amount : -form.amount
        await request.put(`/v1/admin/users/${props.userInfo.id}/credit`, {
          delta: finalAmount,
          reason: form.reason
        })
        ElMessage.success('信用分调整成功')
        emit('success')
        handleClose()
      } catch (e) {
        // request.ts handles error
      } finally {
        loading.value = false
      }
    }
  })
}
</script>

<style scoped>
.current-score { margin-left: 15px; font-weight: bold; color: #409eff; }
</style>