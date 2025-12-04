<template>
  <div class="pwd-container">
    <el-card class="pwd-card" shadow="hover">
      <template #header>
        <div class="card-header">
          <span><el-icon><Lock /></el-icon> 修改密码</span>
        </div>
      </template>
      
      <el-form 
        :model="pwdForm" 
        :rules="rules" 
        ref="formRef" 
        label-width="100px"
        status-icon
      >
        <el-form-item label="原密码" prop="oldPassword">
          <el-input 
            v-model="pwdForm.oldPassword" 
            type="password" 
            show-password 
            placeholder="请输入原密码"
          />
        </el-form-item>
        
        <el-form-item label="新密码" prop="newPassword">
          <el-input 
            v-model="pwdForm.newPassword" 
            type="password" 
            show-password 
            placeholder="6-18位字符"
          />
        </el-form-item>
        
        <el-form-item label="确认密码" prop="confirmPassword">
          <el-input 
            v-model="pwdForm.confirmPassword" 
            type="password" 
            show-password 
            placeholder="请再次输入新密码"
          />
        </el-form-item>
        
        <el-form-item>
          <el-button type="primary" @click="handleSubmit">确认修改</el-button>
          <el-button @click="goBack">返回</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import request from '@/utils/request'
import { ElMessage } from 'element-plus'
import { Lock } from '@element-plus/icons-vue'

const router = useRouter()
const formRef = ref(null)
const pwdForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const validateConfirm = (rule, value, callback) => {
  if (value !== pwdForm.newPassword) callback(new Error('两次密码输入不一致'))
  else callback()
}

const rules = {
  oldPassword: [{ required: true, message: '请输入原密码', trigger: 'blur' }],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, max: 18, message: '长度在 6 到 18 个字符', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    { validator: validateConfirm, trigger: 'blur' }
  ]
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      try {
        await request.post('/api/v1/user/reset-pwd', {
          old_password: pwdForm.oldPassword,
          new_password: pwdForm.newPassword
        })
        ElMessage.success('修改成功，请重新登录')
        localStorage.removeItem('adminToken')
        localStorage.removeItem('adminInfo')
        router.push('/login')
      } catch (error) {
        ElMessage.error(error.response?.data?.message || '修改失败')
      }
    }
  })
}

const goBack = () => router.back()
</script>

<style scoped>
.pwd-container {
  display: flex;
  justify-content: center;
  padding-top: 60px;
  background-color: #f0f2f5;
  min-height: calc(100vh - 60px);
}
.pwd-card {
  width: 500px;
  height: fit-content;
}
.card-header {
  font-weight: bold;
  display: flex;
  align-items: center;
  gap: 8px;
}
</style>