<template>
  <div class="login-container">
    <div class="circles">
      <li></li><li></li><li></li><li></li><li></li>
      <li></li><li></li><li></li><li></li><li></li>
    </div>

    <div class="login-box">
      <div class="login-header">
        <div class="logo-area">
          <div class="logo-circle">
            <el-icon :size="48" color="#4facfe"><School /></el-icon>
          </div>
          <span class="app-name">寝友Go</span>
        </div>
        <div class="header-divider"></div>
        <p class="sub-title">高校宿舍楼栋互助平台<br>综合管理后台</p>
      </div>

      <div class="login-form-area">
        <h2 class="form-title">
          <span>欢迎登录</span>
          <span class="title-tip">Administrator Login</span>
        </h2>
        
        <el-form 
          ref="loginFormRef"
          :model="loginForm"
          :rules="loginRules"
          size="large"
          class="login-form"
          @keyup.enter="handleLogin"
        >
          <el-form-item prop="username">
            <el-input 
              v-model="loginForm.username" 
              placeholder="请输入管理员账号"
              :prefix-icon="User"
              class="custom-input"
            />
          </el-form-item>
          
          <el-form-item prop="password">
            <el-input 
              v-model="loginForm.password" 
              type="password" 
              placeholder="请输入密码"
              show-password
              :prefix-icon="Lock"
              class="custom-input"
            />
          </el-form-item>

          <el-form-item>
            <el-button 
              type="primary" 
              class="submit-btn" 
              :loading="loading" 
              @click="handleLogin"
              round
            >
              {{ loading ? '登录中...' : '立 即 登 录' }}
              <el-icon class="el-icon--right"><ArrowRight /></el-icon>
            </el-button>
          </el-form-item>

          <div class="switch-login">
            <span class="switch-text">普通用户？</span>
            <el-button type="primary" link @click="goToUserLogin">
              用户登录
              <el-icon class="el-icon--right"><Right /></el-icon>
            </el-button>
          </div>
        </el-form>
      </div>
    </div>
    
    <div class="login-copyright">
      © 2025 寝友Go 宿舍管理系统 | Powered by Vue3 & Element Plus
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { User, Lock, School, ArrowRight, Right } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import request from '@/utils/request' // 确保 request 工具路径正确

const router = useRouter()
const loginFormRef = ref(null)
const loading = ref(false)

const loginForm = reactive({
  username: '',
  password: ''
})

const loginRules = {
  username: [{ required: true, message: '请输入管理员账号', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
}

const goToUserLogin = () => {
  router.push('/user/login')
}

const handleLogin = async () => {
  if (!loginFormRef.value) return
  
  await loginFormRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        const res = await request.post('/v1/auth/login', {
          student_id: loginForm.username,
          password: loginForm.password
        })
        
        const { token, user_info } = res;
        
        localStorage.setItem('adminToken', token)
        localStorage.setItem('adminInfo', JSON.stringify(user_info))
        
        ElMessage.success('登录成功，欢迎回来')
        router.push('/')
      } catch (error) {
        // request.ts 已统一处理错误提示
        console.error('Login error:', error)
      } finally {
        loading.value = false
      }
    }
  })
}
</script>

<style scoped>
/* 1. 整体容器与渐变背景 */
.login-container {
  height: 100vh;
  width: 100%;
  /* 核心修改：使用更加现代、美观的 135度 渐变色 */
  background: linear-gradient(135deg, #3B2667 10%, #BC78EC 100%);
  /* 备选方案1 (清新蓝): background: linear-gradient(120deg, #89f7fe 0%, #66a6ff 100%); */
  /* 备选方案2 (深海蓝): background: linear-gradient(to right, #4facfe 0%, #00f2fe 100%); */
  
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  position: relative;
  overflow: hidden;
}

/* 2. 登录卡片样式 */
.login-box {
  width: 850px;
  height: 500px;
  background: rgba(255, 255, 255, 0.95); /* 轻微透明 */
  border-radius: 16px;
  /* 更柔和的投影，制造悬浮感 */
  box-shadow: 0 20px 50px rgba(0, 0, 0, 0.3);
  display: flex;
  overflow: hidden;
  z-index: 10; /* 确保在背景动画之上 */
  backdrop-filter: blur(10px); /* 毛玻璃效果 */
}

/* 3. 左侧品牌区 */
.login-header {
  width: 45%;
  background: #f8fbfd; /* 极淡的蓝色背景 */
  padding: 40px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  border-right: 1px solid #edf2f7;
  position: relative;
}

/* 装饰性背景斜纹 */
.login-header::before {
  content: '';
  position: absolute;
  top: 0; left: 0; width: 100%; height: 100%;
  background-image: radial-gradient(#4facfe 1px, transparent 1px);
  background-size: 20px 20px;
  opacity: 0.05;
}

.logo-area {
  text-align: center;
  margin-bottom: 20px;
  z-index: 1;
}

.logo-circle {
  width: 90px;
  height: 90px;
  background: #fff;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 8px 20px rgba(79, 172, 254, 0.15);
  margin: 0 auto 15px;
  transition: transform 0.3s ease;
}

.logo-circle:hover {
  transform: rotate(15deg) scale(1.05);
}

.app-name {
  font-size: 28px;
  font-weight: 800;
  background: linear-gradient(45deg, #4facfe, #00f2fe);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
  letter-spacing: 2px;
}

.header-divider {
  width: 40px;
  height: 4px;
  background: #e0e0e0;
  border-radius: 2px;
  margin: 20px 0;
}

.sub-title {
  font-size: 15px;
  color: #606266;
  text-align: center;
  line-height: 1.8;
  font-weight: 500;
  z-index: 1;
}

/* 4. 右侧表单区 */
.login-form-area {
  flex: 1;
  padding: 0 60px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  background: #fff;
}

.form-title {
  margin-bottom: 40px;
  display: flex;
  flex-direction: column;
}

.form-title span:first-child {
  font-size: 26px;
  color: #303133;
  font-weight: bold;
}

.title-tip {
  font-size: 14px;
  color: #909399;
  margin-top: 5px;
  font-weight: normal;
  text-transform: uppercase;
  letter-spacing: 1px;
}

.custom-input :deep(.el-input__wrapper) {
  background-color: #f5f7fa;
  box-shadow: none !important;
  border: 1px solid #e4e7ed;
  transition: all 0.3s;
  padding: 12px 15px;
}

.custom-input :deep(.el-input__wrapper.is-focus) {
  background-color: #fff;
  border-color: #bc78ec;
  box-shadow: 0 0 0 1px #bc78ec !important;
}

.submit-btn {
  width: 100%;
  height: 48px;
  font-size: 16px;
  background: linear-gradient(90deg, #3B2667 0%, #BC78EC 100%);
  border: none;
  margin-top: 15px;
  font-weight: 600;
  letter-spacing: 2px;
  transition: all 0.3s;
  box-shadow: 0 8px 15px rgba(188, 120, 236, 0.2);
}

.submit-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 12px 20px rgba(188, 120, 236, 0.3);
  opacity: 0.95;
}

.submit-btn:active {
  transform: translateY(1px);
}

/* 切换登录方式 */
.switch-login {
  display: flex;
  align-items: center;
  justify-content: center;
  margin-top: 20px;
  gap: 4px;
}

.switch-text {
  color: #909399;
  font-size: 14px;
}

.switch-login :deep(.el-button) {
  font-size: 14px;
  color: #bc78ec;
}

.switch-login :deep(.el-button:hover) {
  color: #3B2667;
}

.login-copyright {
  position: absolute;
  bottom: 25px;
  color: rgba(255, 255, 255, 0.7);
  font-size: 12px;
  z-index: 10;
  text-shadow: 0 1px 2px rgba(0,0,0,0.1);
}

/* 5. 动态背景气泡动画 (可选增强效果) */
.circles {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  overflow: hidden;
  z-index: 1;
}

.circles li {
  position: absolute;
  display: block;
  list-style: none;
  width: 20px;
  height: 20px;
  background: rgba(255, 255, 255, 0.2);
  animation: animate 25s linear infinite;
  bottom: -150px;
  border-radius: 50%;
}

.circles li:nth-child(1){ left: 25%; width: 80px; height: 80px; animation-delay: 0s; }
.circles li:nth-child(2){ left: 10%; width: 20px; height: 20px; animation-delay: 2s; animation-duration: 12s; }
.circles li:nth-child(3){ left: 70%; width: 20px; height: 20px; animation-delay: 4s; }
.circles li:nth-child(4){ left: 40%; width: 60px; height: 60px; animation-delay: 0s; animation-duration: 18s; }
.circles li:nth-child(5){ left: 65%; width: 20px; height: 20px; animation-delay: 0s; }
.circles li:nth-child(6){ left: 75%; width: 110px; height: 110px; animation-delay: 3s; }
.circles li:nth-child(7){ left: 35%; width: 150px; height: 150px; animation-delay: 7s; }
.circles li:nth-child(8){ left: 50%; width: 25px; height: 25px; animation-delay: 15s; animation-duration: 45s; }
.circles li:nth-child(9){ left: 20%; width: 15px; height: 15px; animation-delay: 2s; animation-duration: 35s; }
.circles li:nth-child(10){ left: 85%; width: 150px; height: 150px; animation-delay: 0s; animation-duration: 11s; }

@keyframes animate {
  0%{
    transform: translateY(0) rotate(0deg);
    opacity: 1;
    border-radius: 0;
  }
  100%{
    transform: translateY(-1000px) rotate(720deg);
    opacity: 0;
    border-radius: 50%;
  }
}

/* 响应式适配 */
@media (max-width: 900px) {
  .login-box {
    width: 90%;
    height: auto;
    flex-direction: column;
  }
  .login-header {
    width: 100%;
    height: 150px;
    padding: 20px;
    border-right: none;
    border-bottom: 1px solid #eee;
    flex-direction: row;
    justify-content: flex-start;
  }
  .logo-circle { width: 50px; height: 50px; margin: 0 15px 0 0; }
  .logo-area { display: flex; align-items: center; margin-bottom: 0; }
  .header-divider, .sub-title { display: none; }
  .login-form-area { padding: 40px 30px; }
}
</style>