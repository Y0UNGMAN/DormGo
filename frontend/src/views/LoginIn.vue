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
import { useUserStore } from '@/stores/user';
import api from '@/api/index.ts';
const userStore = useUserStore();
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

  if (!isLoginMode.value) {
    if (formData.password !== formData.confirmPassword) {
          alert('两次输入的密码不一致')
          return
      }
      if (!formData.dormId) {
          alert('请选择宿舍楼')
          return
      }
  }

  try {
    // 模拟 API 调用
    // 实际开发请替换为 axios.post(...)
    
    if (isLoginMode.value) {
      console.log('执行登录:', formData)
      // const res = await axios.post('.../login', ...)
      
      // 模拟登录成功，存储用户信息
      const loginPayload = {
        username: formData.username,
        password: formData.password
      }
      const res = await api.post('/api/v1/user/login', loginPayload)
      if (res.data.code === 200) {
        console.log('登录成功，返回数据:', res.data)
        const data = res.data;
        const token = data.token; 
        const user = {
            id: data.user.id,          // 假设后端返回的 user id
            username: data.user.username,  // 假设后端返回的 username
            avatarurl: data.user.avatarurl, // 假设后端返回的 avatar url
            dormid : data.user.dormid,
            intro: data.user.intro || '',
            dorm: data.user.dorm
        };
        userStore.setLogin(token, user);
      alert('登录成功！')
      router.push('/dormgo') // 跳转回主页
      } else {
        alert(`登录失败: ${res.data.message || '用户名或密码错误'}`)
      }
  } else {
      console.log('执行注册:', formData)
      const signupPayload = {
        username: formData.username,
        password: formData.password,
        re_password: formData.confirmPassword // 对应后端需要的 re_password 字段
        // 注意：dormId 字段未包含在您提供的注册请求体中，如果需要提交，请自行添加到 payload 中
      }
      const res = await api.post('/api/v1/user/signup', signupPayload)
      if (res.data.code === 200) {
        alert('注册成功，请立即登录！')
        // 注册成功后，自动切换到登录模式
        isLoginMode.value = true 
      } else {
        // 注册失败，显示后端返回的错误信息
        alert(`注册失败: ${res.data.message || '请检查信息或重试'}`)
      }
      alert('注册成功，请登录')
      isLoginMode.value = true
    }
  }catch (error) {
    console.error('API请求错误:', error)
    const errorMessage = error.response ? error.response.data.message : '网络连接失败，请检查'
    alert(`操作失败: ${errorMessage}`)
  }
}

</script>

<style scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  /* 新增：增加内边距，防止小窗口时卡片贴边 */
  padding: 20px;
  /* 新增：确保 padding 不会撑大容器 */
  box-sizing: border-box;
}

.auth-card {
  background: white;
  /* 修改：宽度改为 100%，让它随容器缩放 */
  width: 100%;
  /* 修改：保留原设计的 400px 作为最大宽度 */
  max-width: 400px;
  /* 新增：设置最小宽度限制，防止缩得太小无法阅读 (你的诉求) */
  min-width: 300px;
  
  padding: 40px;
  border-radius: 20px;
  box-shadow: 0 10px 25px rgba(0,0,0,0.1);
  /* 新增：平滑过渡效果，让调整窗口大小时更顺滑 */
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
  font-weight: 600;
  color: #333;
  margin: 0 0 5px 0;
}

.auth-subtitle {
  color: #999;
  font-size: 14px;
}

.auth-tabs {
  display: flex;
  margin-bottom: 30px;
  border-bottom: 1px solid #f0f0f0;
}

.tab-item {
  flex: 1;
  text-align: center;
  padding: 12px 0;
  cursor: pointer;
  color: #666;
  font-size: 16px;
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
  /* 新增：确保输入框宽度计算正确 */
  box-sizing: border-box;
}

.custom-input:focus {
  outline: none;
  border-color: #1890ff;
  background: white;
  box-shadow: 0 0 0 2px rgba(24, 144, 255, 0.1);
}

.submit-btn {
  width: 100%;
  padding: 12px;
  background: #1890ff;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 16px;
  cursor: pointer;
  transition: all 0.3s;
  margin-top: 10px;
}

.submit-btn:hover {
  background: #40a9ff;
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(24, 144, 255, 0.3);
}

/* --- 响应式适配 (关键优化部分) --- */

/* 当窗口高度较小（如横屏手机或压扁的浏览器窗口）时 */
@media (max-height: 600px) {
  .login-container {
    /* 取消垂直居中，改为顶部对齐，防止内容被切掉无法滚动 */
    align-items: flex-start;
    padding-top: 40px;
    padding-bottom: 40px;
    /* 允许垂直滚动 */
    overflow-y: auto; 
  }
}

/* 当窗口宽度较小（移动端或缩小窗口）时 */
@media (max-width: 480px) {
  .auth-card {
    /* 减小内边距，腾出更多空间给内容 */
    padding: 24px;
    /* 在小屏幕上稍微减小圆角 */
    border-radius: 16px;
  }

  /* 稍微缩小标题字体 */
  .auth-title {
    font-size: 20px;
  }

  .logo-circle {
    width: 50px;
    height: 50px;
    font-size: 24px;
  }
}
</style>