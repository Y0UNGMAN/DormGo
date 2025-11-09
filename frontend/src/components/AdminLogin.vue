<template>
  <div class="login-container">
    <!-- 登录卡片 -->
    <div class="login-card">
      <h2 class="login-title">寝友Go - 管理员登录</h2>
      
      <!-- 登录表单 -->
      <form @submit.prevent="handleLogin">
        <!-- 账号输入 -->
        <div class="form-item">
          <label>账号（学号/管理员ID）</label>
          <input 
            type="text" 
            v-model="loginForm.username" 
            required 
            placeholder="请输入管理员账号"
          >
        </div>
        
        <!-- 密码输入 -->
        <div class="form-item">
          <label>密码</label>
          <input 
            type="password" 
            v-model="loginForm.password" 
            required 
            placeholder="请输入密码"
          >
        </div>
        
        <!-- 错误提示 -->
        <div class="error-message" v-if="errorMsg">{{ errorMsg }}</div>
        
        <!-- 登录按钮 -->
        <button type="submit" class="login-btn">管理员登录</button>
        
        <!-- 普通用户登录入口 -->
        <div class="user-switch">
          普通用户？<a @click="$router.push('/user/login')">点击登录</a>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'

// 响应式数据
const loginForm = ref({
  username: '',
  password: ''
})
const errorMsg = ref('')
const router = useRouter()

// 登录处理
const handleLogin = async () => {
  try {
    const response = await axios.post('/api/v1/auth/login', {
      student_id: loginForm.value.username,
      password: loginForm.value.password
    })
    
    const { token, user_info } = response.data
    if (user_info.is_admin) {
      localStorage.setItem('adminToken', token)
      localStorage.setItem('adminInfo', JSON.stringify(user_info))
      router.push('/')
    } else {
      errorMsg.value = '该账号不是管理员，请使用管理员账号登录'
    }
  } catch (error) {
    errorMsg.value = error.response?.data?.message || '登录失败，请检查账号密码'
  }
}
</script>

<style scoped>
/* 登录页面样式 */
.login-container {
  width: 100%;
  height: 100vh;
  background: #f5f7fa;
  display: flex;
  justify-content: center;
  align-items: center;
}

.login-card {
  width: 350px;
  padding: 30px;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
}

.login-title {
  text-align: center;
  color: #333;
  margin-bottom: 25px;
  font-size: 18px;
}

.form-item {
  margin-bottom: 20px;
}

.form-item label {
  display: block;
  margin-bottom: 8px;
  color: #666;
  font-size: 14px;
}

.form-item input {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
  box-sizing: border-box;
}

.form-item input:focus {
  outline: none;
  border-color: #42b983;
}

.error-message {
  color: #f56c6c;
  font-size: 12px;
  text-align: center;
  margin-bottom: 15px;
}

.login-btn {
  width: 100%;
  padding: 12px;
  background: #42b983;
  color: #fff;
  border: none;
  border-radius: 4px;
  font-size: 14px;
  cursor: pointer;
  transition: background 0.3s;
}

.login-btn:hover {
  background: #359469;
}

.user-switch {
  text-align: center;
  margin-top: 15px;
  font-size: 13px;
  color: #666;
}

.user-switch a {
  color: #42b983;
  text-decoration: none;
  margin-left: 5px;
}
</style>