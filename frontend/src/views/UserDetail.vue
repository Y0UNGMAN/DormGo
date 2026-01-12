<template>
  <div class="page-container">
    <div class="nav-bar">
      <button class="back-btn" @click="goBack">← 返回个人中心</button>
      <span class="page-title">详细资料</span>
      <button class="save-btn" @click="saveProfile" :disabled="!isModified || isSaving">
        {{ isSaving ? '保存中...' : '保存' }}
      </button>
    </div>

    <div class="avatar-section">
      <div class="avatar-wrapper" @click="triggerUpload">
        <img :src="form.avatar || defaultAvatar" alt="头像" class="avatar" />
        <div class="avatar-overlay">
          <span>更换头像</span>
        </div>
      </div>
      <input 
        type="file" 
        ref="fileInput" 
        @change="handleAvatarChange" 
        accept="image/*" 
        style="display: none"
      />
    </div>

    <div class="form-section">
      <div class="section-title">基本资料</div>
      <div class="form-group">
        <label class="form-label">昵称</label>
        <input v-model="form.nickname" type="text" class="form-input" placeholder="请输入昵称" />
      </div>

      <div class="form-group">
        <label class="form-label">学号</label>
        <input v-model="form.studentId" type="text" class="form-input" placeholder="请输入学号" disabled />
      </div>

      <div class="form-group">
        <label class="form-label">宿舍楼</label>
        <el-select v-model="form.dormId" placeholder="请选择宿舍楼" class="form-select">
          <el-option 
            v-for="dorm in dormOptions" 
            :key="dorm.dormid" 
            :label="dorm.dormname" 
            :value="dorm.dormid" 
          />
        </el-select>
      </div>

      <div class="form-group">
        <label class="form-label">个性签名</label>
        <textarea v-model="form.bio" class="form-textarea" placeholder="介绍一下自己吧..." rows="3"></textarea>
      </div>
    </div>

    <div class="form-section mt-12">
      <div class="section-title">账号安全</div>
      <div class="security-item" @click="showPasswordModal = true">
        <span class="security-label">修改密码</span>
        <span class="security-arrow"> > </span>
      </div>
    </div>

    <div v-if="showPasswordModal" class="modal-overlay" @click.self="closePasswordModal">
      <div class="modal-card">
        <h3 class="modal-title">修改密码</h3>
        
        <div class="modal-form">
          <div class="form-group">
            <input 
              v-model="passwordForm.oldPassword" 
              type="password" 
              class="form-input" 
              placeholder="请输入旧密码"
            />
          </div>
          <div class="form-group">
            <input 
              v-model="passwordForm.newPassword" 
              type="password" 
              class="form-input" 
              placeholder="请输入新密码"
            />
          </div>
          <div class="form-group">
            <input 
              v-model="passwordForm.confirmPassword" 
              type="password" 
              class="form-input" 
              placeholder="请确认新密码"
            />
          </div>
        </div>

        <div class="modal-actions">
          <button class="cancel-btn" @click="closePasswordModal">取消</button>
          <button 
            class="confirm-btn" 
            @click="submitPasswordChange" 
            :disabled="isSubmittingPwd"
          >
            {{ isSubmittingPwd ? '提交中...' : '确认修改' }}
          </button>
        </div>
      </div>
    </div>

  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'
import api from '@/api/index'

const router = useRouter()
const userStore = useUserStore()
const fileInput = ref<HTMLInputElement | null>(null)

const defaultAvatar = 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'

// --- 个人资料相关数据 ---
const originalForm = ref<any>({})
const form = ref({
  avatar: '',
  nickname: '',
  studentId: '',
  dormId: 0,
  bio: ''
})

const dormOptions = ref<Array<{dormid: number, dormname: string}>>([])
const isSaving = ref(false)

// --- 【新增】修改密码相关数据 ---
const showPasswordModal = ref(false)
const isSubmittingPwd = ref(false)
const passwordForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const isModified = computed(() => {
  return JSON.stringify(form.value) !== JSON.stringify(originalForm.value)
})

const goBack = () => {
  if (isModified.value) {
    if (confirm('有未保存的修改，确定要离开吗？')) {
      router.push('/personalhome')
    }
  } else {
    router.push('/personalhome')
  }
}

// --- 头像上传逻辑 ---
const triggerUpload = () => {
  fileInput.value?.click()
}

const handleAvatarChange = async (e: Event) => {
  const target = e.target as HTMLInputElement
  const file = target.files?.[0]
  if (!file) return

  if (file.size > 2 * 1024 * 1024) {
    ElMessage.warning('图片大小不能超过2MB')
    return
  }

  try {
    const formData = new FormData()
    formData.append('file', file)
    
    const res = await api.post('/api/v1/user/upload/avatar', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
    
    if (res.data.code === 200) {
      form.value.avatar = res.data.data.url
      ElMessage.success('头像上传成功')
    } else {
      ElMessage.error(res.data.msg || '上传失败')
    }
  } catch (error) {
    ElMessage.error('上传失败，请重试')
  }
}

// --- 保存资料逻辑 ---
const saveProfile = async () => {
  if (!isModified.value || isSaving.value) return
  isSaving.value = true
  
  try {
    const res = await api.put('/api/v1/user/profile', {
      nickname: form.value.nickname,
      avatar: form.value.avatar,
      dorm_id: form.value.dormId,
      bio: form.value.bio
    })
    
    if (res.data.code === 200) {
      if (userStore.userInfo) {
        userStore.userInfo.username = form.value.nickname
        userStore.userInfo.avatarurl = form.value.avatar
        userStore.userInfo.dormid = form.value.dormId
        userStore.userInfo.intro = form.value.bio
        // 假设 userStore 内部会处理 localStorage，如果不会则手动存
        localStorage.setItem('userInfo', JSON.stringify(userStore.userInfo))
      }
      originalForm.value = { ...form.value }
      ElMessage.success('保存成功')
    } else {
      ElMessage.error(res.data.msg || '保存失败')
    }
  } catch (error) {
    ElMessage.error('保存失败，请重试')
  } finally {
    isSaving.value = false
  }
}

// --- 【新增】修改密码逻辑 ---
const closePasswordModal = () => {
  showPasswordModal.value = false
  // 清空表单
  passwordForm.oldPassword = ''
  passwordForm.newPassword = ''
  passwordForm.confirmPassword = ''
}

const submitPasswordChange = async () => {
  // 1. 前端校验
  if (!passwordForm.oldPassword || !passwordForm.newPassword || !passwordForm.confirmPassword) {
    ElMessage.warning('请填写所有密码字段')
    return
  }
  if (passwordForm.newPassword !== passwordForm.confirmPassword) {
    ElMessage.warning('两次输入的新密码不一致')
    return
  }
  if (passwordForm.newPassword.length < 6) {
    ElMessage.warning('新密码长度不能少于6位')
    return
  }

  isSubmittingPwd.value = true

  try {
    // 2. 调用后端接口
    // 对应 controller/user_extension_con.go 中的 ResetPassword
    const res = await api.post('/api/v1/user/reset-pwd', {
      old_password: passwordForm.oldPassword,
      new_password: passwordForm.newPassword
    })

    if (res.data.code === 200) {
      ElMessage.success('密码修改成功，请重新登录')
      closePasswordModal()
      // 修改密码后通常需要强制重新登录
      setTimeout(() => {
        userStore.logout()
        router.push('/')
      }, 1500)
    } else {
      ElMessage.error(res.data.msg || '修改失败')
    }
  } catch (error) {
    console.error(error)
    ElMessage.error('网络错误，修改失败')
  } finally {
    isSubmittingPwd.value = false
  }
}

// --- 初始化逻辑 ---
const loadUserProfile = async () => {
  try {
    const res = await api.get('/api/v1/user/profile')
    if (res.data.code === 200) {
      const data = res.data.data
      form.value = {
        avatar: data.avatar || '',
        nickname: data.nickname || '',
        studentId: data.student_id || '',
        dormId: data.dorm_id || 0,
        bio: data.intro || ''
      }
      originalForm.value = { ...form.value }
    }
  } catch (error) {
    const user = userStore.currentUser
    if (user) {
      form.value = {
        avatar: user.avatarurl || '',
        nickname: user.username || '',
        studentId: '',
        dormId: user.dormid || 0,
        bio: user.intro || ''
      }
      originalForm.value = { ...form.value }
    }
  }
}

const loadDormOptions = async () => {
  try {
    const res = await api.get('/api/v1/dorm/list')
    if (res.data.code === 200) {
      dormOptions.value = res.data.data || []
    }
  } catch (error) {
    console.error('获取宿舍列表失败', error)
  }
}

onMounted(() => {
  loadUserProfile()
  loadDormOptions()
})
</script>

<style scoped>
.page-container {
  min-height: 100vh;
  background: #f5f5f5;
  padding-bottom: 30px;
}

.nav-bar {
  background: white;
  padding: 15px 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  box-shadow: 0 1px 4px rgba(0,0,0,0.05);
  position: sticky;
  top: 0;
  z-index: 100;
}

.back-btn, .save-btn {
  border: none;
  background: none;
  cursor: pointer;
  font-size: 14px;
  color: #666;
  transition: color 0.3s;
}

.back-btn:hover { color: #1890ff; }

.save-btn {
  color: #1890ff;
  font-weight: 600;
}
.save-btn:disabled { color: #ccc; cursor: not-allowed; }

.page-title {
  font-weight: 600;
  font-size: 16px;
  color: #333;
}

/* 头像区域 */
.avatar-section {
  background: white;
  padding: 24px;
  display: flex;
  justify-content: center;
  margin-bottom: 12px;
}

.avatar-wrapper {
  position: relative;
  width: 100px;
  height: 100px;
  border-radius: 50%;
  overflow: hidden;
  cursor: pointer;
}

.avatar { width: 100%; height: 100%; object-fit: cover; }

.avatar-overlay {
  position: absolute; bottom: 0; left: 0; right: 0;
  background: rgba(0, 0, 0, 0.5); color: white;
  text-align: center; padding: 6px 0; font-size: 12px;
  opacity: 0; transition: opacity 0.3s;
}
.avatar-wrapper:hover .avatar-overlay { opacity: 1; }

/* 表单区域 */
.form-section {
  background: white;
  margin: 0 12px 12px 12px;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.section-title {
  font-size: 15px;
  font-weight: 600;
  color: #333;
  margin-bottom: 16px;
  padding-left: 8px;
  border-left: 3px solid #1890ff;
}

.form-group { margin-bottom: 20px; }
.form-group:last-child { margin-bottom: 0; }

.form-label {
  display: block; font-size: 14px; color: #666; font-weight: 500; margin-bottom: 8px;
}

.form-input {
  width: 100%; padding: 12px 14px; border: 1px solid #e8e8e8; border-radius: 8px;
  font-size: 14px; transition: all 0.3s; box-sizing: border-box;
}

.form-input:focus {
  outline: none; border-color: #1890ff; box-shadow: 0 0 0 2px rgba(24, 144, 255, 0.1);
}
.form-input:disabled { background: #f9f9f9; color: #999; cursor: not-allowed; }

.form-select { width: 100%; }
.form-select :deep(.el-input__wrapper) { padding: 8px 14px; border-radius: 8px; }

.form-textarea {
  width: 100%; padding: 12px 14px; border: 1px solid #e8e8e8; border-radius: 8px;
  font-size: 14px; resize: vertical; font-family: inherit; box-sizing: border-box;
}
.form-textarea:focus {
  outline: none; border-color: #1890ff; box-shadow: 0 0 0 2px rgba(24, 144, 255, 0.1);
}

/* 账号安全项 */
.security-item {
  display: flex; justify-content: space-between; align-items: center;
  padding: 12px 0; cursor: pointer; border-bottom: 1px solid #f5f5f5;
}
.security-item:last-child { border-bottom: none; }
.security-label { font-size: 14px; color: #333; }
.security-arrow { color: #ccc; font-family: consolas, sans-serif; }
.security-item:hover .security-arrow { color: #1890ff; }

.mt-12 { margin-top: 12px; }

/* 弹窗样式 */
.modal-overlay {
  position: fixed; top: 0; left: 0; right: 0; bottom: 0;
  background: rgba(0, 0, 0, 0.6);
  display: flex; justify-content: center; align-items: center;
  z-index: 1000;
  animation: fadeIn 0.2s ease;
}

.modal-card {
  background: white; width: 85%; max-width: 320px;
  border-radius: 12px; padding: 24px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  text-align: center;
  animation: scaleUp 0.2s ease;
}

.modal-title { font-size: 18px; color: #333; margin: 0 0 20px 0; font-weight: 600; }
.modal-form { text-align: left; }

.modal-actions {
  display: flex; gap: 12px; margin-top: 24px;
}

.cancel-btn, .confirm-btn {
  flex: 1; padding: 10px 0; border-radius: 6px; border: none; font-size: 14px; cursor: pointer;
}
.cancel-btn { background: #f5f5f5; color: #666; }
.cancel-btn:hover { background: #e8e8e8; }
.confirm-btn { background: #1890ff; color: white; }
.confirm-btn:hover { background: #40a9ff; }
.confirm-btn:disabled { background: #ccc; cursor: not-allowed; }

@keyframes fadeIn { from { opacity: 0; } to { opacity: 1; } }
@keyframes scaleUp { from { transform: scale(0.9); opacity: 0; } to { transform: scale(1); opacity: 1; } }
</style>