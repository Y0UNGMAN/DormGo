<template>
  <div class="login-container">
    <div class="auth-card">
      <div class="auth-header">
        <div class="logo-circle">🏠</div>
        <h2 class="auth-title">DormGo 宿舍集市</h2>
        <p class="auth-subtitle">校园生活，一站搞定</p>
      </div>

      <div class="auth-tabs">
        <div 
          class="tab-item" 
          :class="{ active: isLoginMode }" 
          @click="isLoginMode = true"
        >登录</div>
        <div 
          class="tab-item" 
          :class="{ active: !isLoginMode }" 
          @click="isLoginMode = false"
        >注册</div>
      </div>

      <div class="form-container">
        <div class="input-group">
          <span class="input-icon">👤</span>
          <input 
            v-model="formData.username" 
            type="text" 
            placeholder="请输入用户名"
            class="custom-input"
          />
        </div>

        <div class="input-group">
          <span class="input-icon">🔒</span>
          <input 
            v-model="formData.password" 
            type="password" 
            placeholder="请输入密码"
            class="custom-input"
          />
        </div>

        <div class="input-group" v-if="!isLoginMode">
          <span class="input-icon">🛡️</span>
          <input 
            v-model="formData.confirmPassword" 
            type="password" 
            placeholder="请确认密码"
            class="custom-input"
          />
        </div>

        <div class="input-group" v-if="!isLoginMode">
          <span class="input-icon">🏢</span>
          <select v-model="formData.dormId" class="custom-input">
            <option value="" disabled selected>请选择宿舍楼</option>
            <option value="1">海馨苑</option>
            <option value="2">海韵苑</option>
            </select>
        </div>

        <button class="submit-btn" @click="handleSubmit">
          {{ isLoginMode ? '立即登录' : '创建账号' }}
        </button>

        <div class="auth-footer">
          <span v-if="isLoginMode" class="footer-link">忘记密码？</span>
          <span class="footer-text" @click="toggleMode">
            {{ isLoginMode ? '没有账号？去注册' : '已有账号？去登录' }}
          </span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'

const router = useRouter()
const isLoginMode = ref(true)

const formData = reactive({
  username: '',
  password: '',
  confirmPassword: '',
  dormId: ''
})

const toggleMode = () => {
  isLoginMode.value = !isLoginMode.value
  // 清空表单
  formData.password = ''
  formData.confirmPassword = ''
}

const handleSubmit = async () => {
  // 1. 简单的表单校验
  if (!formData.username || !formData.password) {
    alert('请填写完整信息')
    return
  }

  if (!isLoginMode.value && formData.password !== formData.confirmPassword) {
    alert('两次输入的密码不一致')
    return
  }

  try {
    // 模拟 API 调用
    // 实际开发请替换为 axios.post(...)
    
    if (isLoginMode.value) {
      console.log('执行登录:', formData)
      // const res = await axios.post('.../login', ...)
      
      // 模拟登录成功，存储用户信息
      localStorage.setItem('token', 'mock_token_123456')
      localStorage.setItem('user', JSON.stringify({
        name: formData.username,
        avatar: 'https://dorm-go.oss-cn-guangzhou.aliyuncs.com/avator/midnight.jpg',
        id: 1
      }))
      
      alert('登录成功！')
      router.push('/') // 跳转回主页
    } else {
      console.log('执行注册:', formData)
      // const res = await axios.post('.../register', ...)
      alert('注册成功，请登录')
      isLoginMode.value = true
    }

  } catch (error) {
    console.error(error)
    alert('操作失败，请重试')
  }
}
</script>

<style scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #e6f7ff 0%, #ffffff 100%);
  padding: 20px;
}

.auth-card {
  width: 100%;
  max-width: 400px;
  background: white;
  border-radius: 20px;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.05);
  padding: 40px 30px;
  transition: all 0.3s ease;
}

.auth-header {
  text-align: center;
  margin-bottom: 30px;
}

.logo-circle {
  width: 60px;
  height: 60px;
  background: #e6f7ff;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 30px;
  margin: 0 auto 15px;
  color: #1890ff;
}

.auth-title {
  font-size: 24px;
  color: #333;
  margin-bottom: 8px;
  font-weight: 600;
}

.auth-subtitle {
  color: #999;
  font-size: 14px;
}

/* Tab 切换 */
.auth-tabs {
  display: flex;
  margin-bottom: 30px;
  border-bottom: 1px solid #f0f0f0;
}

.tab-item {
  flex: 1;
  text-align: center;
  padding: 12px;
  cursor: pointer;
  color: #666;
  font-weight: 500;
  position: relative;
  transition: color 0.3s;
}

.tab-item.active {
  color: #1890ff;
  font-weight: 600;
}

.tab-item.active::after {
  content: '';
  position: absolute;
  bottom: -1px;
  left: 0;
  width: 100%;
  height: 2px;
  background: #1890ff;
}

/* 表单样式 */
.input-group {
  position: relative;
  margin-bottom: 20px;
}

.input-icon {
  position: absolute;
  left: 15px;
  top: 50%;
  transform: translateY(-50%);
  font-size: 16px;
}

.custom-input {
  width: 100%;
  padding: 12px 15px 12px 45px;
  border: 1px solid #e8e8e8;
  border-radius: 8px;
  font-size: 14px;
  transition: all 0.3s;
  background: #fafafa;
}

.custom-input:focus {
  outline: none;
  border-color: #1890ff;
  background: white;
  box-shadow: 0 0 0 2px rgba(24, 144, 255, 0.1);
}

.submit-btn {
  width: 100%;
  padding: 14px;
  background: #1890ff;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s;
  margin-top: 10px;
}

.submit-btn:hover {
  background: #40a9ff;
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(24, 144, 255, 0.2);
}

.auth-footer {
  margin-top: 20px;
  display: flex;
  justify-content: space-between;
  font-size: 13px;
}

.footer-link {
  color: #999;
  cursor: pointer;
}

.footer-text {
  color: #1890ff;
  cursor: pointer;
  margin-left: auto;
}

.footer-text:hover {
  text-decoration: underline;
}
</style>