<template>
  <div class="login-container">
    <div class="circles">
      <div v-for="n in 10" :key="n"></div>
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
        <p class="sub-title">校园生活 · 宿舍互助 · 便捷服务</p>
        <div class="illustration-text">
          <p>一站式解决宿舍生活难题</p>
          <p>连接你我，温暖互助</p>
        </div>
      </div>

      <div class="login-form-area">
        <h2 class="form-title">
          <span>{{ titleText }}</span>
          <span class="title-tip">{{ titleTip }}</span>
        </h2>

        <el-tabs 
          v-model="identity" 
          class="identity-tabs" 
          stretch
          @tab-change="handleTabChange"
        >
          <el-tab-pane label="普通用户" name="user"></el-tab-pane>
          <el-tab-pane label="管理员" name="admin"></el-tab-pane>
        </el-tabs>
        
        <el-form 
          ref="authFormRef"
          :model="formData"
          :rules="formRules"
          size="large"
          class="auth-form"
          @keyup.enter="handleSubmit"
        >
          <el-form-item prop="username">
            <el-input 
              v-model="formData.username" 
              :placeholder="identity === 'admin' ? '请输入管理员账号' : (isLoginMode ? '请输入用户名/学号' : '请输入用户名(昵称)')"
              :prefix-icon="User"
            />
          </el-form-item>

          <template v-if="identity === 'user' && !isLoginMode">
            <el-form-item prop="studentId">
              <el-input 
                v-model="formData.studentId" 
                placeholder="请输入学号"
                :prefix-icon="School"
              />
            </el-form-item>
          </template>
          
          <el-form-item prop="password">
            <el-input 
              v-model="formData.password" 
              type="password" 
              placeholder="请输入密码"
              show-password
              :prefix-icon="Lock"
            />
          </el-form-item>

          <template v-if="identity === 'user' && !isLoginMode">
            <el-form-item prop="confirmPassword">
              <el-input 
                v-model="formData.confirmPassword" 
                type="password" 
                placeholder="请确认密码"
                show-password
                :prefix-icon="Lock"
              />
            </el-form-item>
            
            <el-form-item prop="dormId">
              <el-select v-model="formData.dormId" placeholder="请选择宿舍楼" style="width: 100%">
                <el-option v-for="dorm in dormList" :key="dorm.dormid" :label="dorm.dormname" :value="dorm.dormid" />
              </el-select>
            </el-form-item>
          </template>

          <el-form-item>
            <el-button 
              type="primary" 
              class="submit-btn" 
              :loading="loading" 
              @click="handleSubmit"
              round
            >
              {{ submitBtnText }}
              <el-icon class="el-icon--right"><ArrowRight /></el-icon>
            </el-button>
          </el-form-item>

          <div class="form-footer" v-if="identity === 'user'">
            <span class="toggle-mode" @click="toggleMode">
              {{ isLoginMode ? '没有账号？去注册' : '已有账号？去登录' }}
            </span>
          </div>
        </el-form>
      </div>
    </div>
    
    <div class="login-copyright">
      © 2025 寝友Go 宿舍管理系统
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { User, Lock, School, ArrowRight } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import api from '@/api/index'
// 1. 引入 Pinia Store
import { useUserStore } from '@/stores/user';

const router = useRouter()
// 2. 初始化 Store
const userStore = useUserStore();

const authFormRef = ref(null)
const loading = ref(false)

// 状态控制
const identity = ref('user') // 'user' | 'admin'
const isLoginMode = ref(true)
const dormList = ref([])

const formData = reactive({
  username: '',
  password: '',
  confirmPassword: '',
  studentId: '', // 补充定义 studentId
  dormId: ''
})

// === 计算属性 ===

const titleText = computed(() => {
  if (identity.value === 'admin') return '管理员登录'
  return isLoginMode.value ? '欢迎登录' : '注册账号'
})

const titleTip = computed(() => {
  if (identity.value === 'admin') return 'Admin Portal'
  return isLoginMode.value ? 'Welcome Back' : 'Create Account'
})

const submitBtnText = computed(() => {
  if (identity.value === 'admin') return '登 录'
  return isLoginMode.value ? '登 录' : '注 册'
})

// 表单校验规则
const formRules = computed(() => {
  const rules = {
    username: [{ required: true, message: '请输入账号/用户名', trigger: 'blur' }],
    password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
  }
  
  if (identity.value === 'user' && !isLoginMode.value) {
    rules.confirmPassword = [
      { required: true, message: '请确认密码', trigger: 'blur' },
      { 
        validator: (rule, value, callback) => {
          if (value !== formData.password) callback(new Error('两次输入密码不一致'))
          else callback()
        }, 
        trigger: 'blur' 
      }
    ]
    rules.dormId = [{ required: true, message: '请选择宿舍楼', trigger: 'change' }]
  }
  return rules
})

// === 生命周期 ===
onMounted(() => {
  fetchDormList()
})

// === 方法 ===

const fetchDormList = async () => {
  try {
    const res = await api.get('/api/v1/post/dorms')
    dormList.value = res.data.data || []
  } catch (error) {
    console.error('获取宿舍列表失败', error)
  }
}

const handleTabChange = (tab) => {
  authFormRef.value?.resetFields()
  if (tab === 'admin') {
    isLoginMode.value = true
  }
}

const toggleMode = () => {
  isLoginMode.value = !isLoginMode.value
  authFormRef.value?.resetFields()
}

const handleSubmit = async () => {
  if (!authFormRef.value) return
  
  await authFormRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        // === 管理员登录 (保持独立逻辑，不走UserStore) ===
        if (identity.value === 'admin') {
          const res = await api.post('/api/v1/admin/login', {
            username: formData.username,
            password: formData.password
          })

          if (res.data.code === 200) { // 注意：res结构需根据拦截器确认，这里假设axios返回包含data
             // 兼容直接返回data或res.data的情况
             const responseData = res.data || res;
            const adminToken = responseData.token;
            const adminInfoObj = { 
              username: formData.username, 
              // 处理头像：后端有就用后端的，没有就用默认图
              avatar: responseData.data?.avatar || 'https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png',
              // 如果后端返回了 id，也可以加上
              id: responseData.data?.id
            }
            userStore.setAdminLogin(adminToken, adminInfoObj);
            userStore.logout();
            ElMessage.success('管理员登录成功')
            router.push('/admin/statistics')
          } else {
            ElMessage.error(res.data?.message || '登录失败')
          }
        } 
        // === 普通用户操作 (接入 Pinia) ===
        else {
          if (isLoginMode.value) {
            // -- 登录 --
            const res = await api.post('/api/v1/user/login', {
              username: formData.username,
              password: formData.password
            })
            
            if (res.data.code === 200) {
              const token = res.data.token;
              const backendUser = res.data.user || {};

              // 3. 构建符合 Store 定义的 User 对象
              // 注意：user.ts 接口定义包含 id, username, avatarurl, dormid
              const userInfoObj = {
                id: backendUser.id || 0,
                username: backendUser.username || formData.username,
                avatarurl: backendUser.avatarurl || 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png',
                dormid: backendUser.dormid || 0,
                // 如果需要把 token 也放在 userInfo 里（根据 user.ts 的 Interface），这里加上
                token: token
              };

              // 4. 调用 Pinia Action 更新状态
              // setLogin 会自动处理 state 更新和 localStorage 持久化
              userStore.setLogin(token, userInfoObj);
              
              ElMessage.success('登录成功')
              router.push('/dormgo')
            } else {
              ElMessage.error(res.data.message || '登录失败')
            }
          } else {
            // -- 注册 --
            const res = await api.post('/api/v1/user/signup', {
              username: formData.username,
              password: formData.password,
              re_password: formData.confirmPassword,
              student_id:formData.studentId,
              dorm_id: Number(formData.dormId)
            })
            
            if (res.data.code === 200) {
              ElMessage.success('注册成功，请立即登录')
              isLoginMode.value = true
              // 自动填充用户名
              formData.password = ''
              formData.confirmPassword = ''
            } else {
              ElMessage.error(res.data.message || '注册失败')
            }
          }
        }
      } catch (error) {
        console.error(error)
        ElMessage.error(error.response?.data?.message || '网络连接失败')
      } finally {
        loading.value = false
      }
    }
  })
}
</script>

<style scoped>
/* 样式保持不变，直接复用原文件的样式 */
.login-container {
  height: 100vh;
  width: 100%;
  background: linear-gradient(135deg, #3B2667 10%, #BC78EC 100%);
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  position: relative;
  overflow: hidden;
}

.login-box {
  width: 900px;
  min-height: 550px;
  background: rgba(255, 255, 255, 0.95);
  border-radius: 16px;
  box-shadow: 0 20px 50px rgba(0, 0, 0, 0.3);
  display: flex;
  overflow: hidden;
  z-index: 10;
  backdrop-filter: blur(10px);
}

.login-header {
  width: 40%;
  background: #f8fbfd;
  padding: 40px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  border-right: 1px solid #edf2f7;
  position: relative;
}

.login-header::before {
  content: '';
  position: absolute;
  top: 0; left: 0; width: 100%; height: 100%;
  background-image: radial-gradient(#4facfe 1px, transparent 1px);
  background-size: 20px 20px;
  opacity: 0.05;
}

.logo-circle {
  width: 80px; height: 80px;
  background: #fff;
  border-radius: 50%;
  display: flex; align-items: center; justify-content: center;
  box-shadow: 0 8px 20px rgba(79, 172, 254, 0.15);
  margin-bottom: 15px;
}

.app-name {
  font-size: 26px; font-weight: 800;
  background: linear-gradient(45deg, #4facfe, #00f2fe);
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
}

.header-divider {
  width: 40px; height: 4px; background: #e0e0e0; border-radius: 2px; margin: 20px 0;
}

.sub-title { font-size: 16px; color: #333; font-weight: bold; margin-bottom: 20px; }
.illustration-text { text-align: center; color: #666; font-size: 14px; line-height: 1.8; }

.login-form-area {
  flex: 1;
  padding: 40px 60px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  background: #fff;
}

.form-title {
  margin-bottom: 20px;
  display: flex; flex-direction: column;
}

.form-title span:first-child { font-size: 26px; color: #303133; font-weight: bold; }
.title-tip { font-size: 14px; color: #909399; margin-top: 5px; font-weight: normal; letter-spacing: 1px; }

.identity-tabs { margin-bottom: 20px; }

.submit-btn {
  width: 100%; height: 44px; font-size: 16px;
  background: linear-gradient(90deg, #3B2667 0%, #BC78EC 100%);
  border: none; margin-top: 10px; font-weight: 600; letter-spacing: 2px;
  transition: all 0.3s;
}
.submit-btn:hover { transform: translateY(-2px); opacity: 0.9; }

.form-footer { display: flex; justify-content: center; margin-top: 15px; font-size: 14px; }
.toggle-mode { color: #1890ff; cursor: pointer; }
.toggle-mode:hover { text-decoration: underline; }

.login-copyright {
  position: absolute; bottom: 20px; color: rgba(255,255,255,0.6); font-size: 12px; z-index: 10;
}

.circles div {
  position: absolute; display: block; list-style: none; width: 20px; height: 20px;
  background: rgba(255, 255, 255, 0.2); animation: animate 25s linear infinite;
  bottom: -150px; border-radius: 50%;
}
.circles div:nth-child(1){ left: 25%; width: 80px; height: 80px; animation-delay: 0s; }
.circles div:nth-child(2){ left: 10%; width: 20px; height: 20px; animation-delay: 2s; animation-duration: 12s; }
.circles div:nth-child(3){ left: 70%; width: 20px; height: 20px; animation-delay: 4s; }

@keyframes animate {
  0%{ transform: translateY(0) rotate(0deg); opacity: 1; border-radius: 0; }
  100%{ transform: translateY(-1000px) rotate(720deg); opacity: 0; border-radius: 50%; }
}
</style>