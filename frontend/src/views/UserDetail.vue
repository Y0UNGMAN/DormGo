<template>
  <div class="page-container">
    <!-- 顶部导航 -->
    <div class="nav-bar">
      <button class="back-btn" @click="goBack">← 返回个人中心</button>
      <span class="page-title">详细资料</span>
      <button class="save-btn" @click="saveProfile" :disabled="!isModified || isSaving">
        {{ isSaving ? '保存中...' : '保存' }}
      </button>
    </div>

    <!-- 头像区域 -->
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

    <!-- 表单区域 -->
    <div class="form-section">
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
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'
import api from '@/api/index'

const router = useRouter()
const userStore = useUserStore()
const fileInput = ref<HTMLInputElement | null>(null)

const defaultAvatar = 'https://dorm-go.oss-cn-guangzhou.aliyuncs.com/avator/midnight.jpg'

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

const triggerUpload = () => {
  fileInput.value?.click()
}

const handleAvatarChange = async (e: Event) => {
  const target = e.target as HTMLInputElement
  const file = target.files?.[0]
  if (!file) return

  // 检查文件大小（限制 2MB）
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
      // 更新 store 中的用户信息
      if (userStore.userInfo) {
        userStore.userInfo.username = form.value.nickname
        userStore.userInfo.avatarurl = form.value.avatar
        userStore.userInfo.dormid = form.value.dormId
        userStore.userInfo.intro = form.value.bio
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
    // 如果 API 失败，使用 store 中的数据
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

.back-btn:hover {
  color: #1890ff;
}

.save-btn {
  color: #1890ff;
  font-weight: 600;
}

.save-btn:disabled {
  color: #ccc;
  cursor: not-allowed;
}

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

.avatar {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar-overlay {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  background: rgba(0, 0, 0, 0.5);
  color: white;
  text-align: center;
  padding: 6px 0;
  font-size: 12px;
  opacity: 0;
  transition: opacity 0.3s;
}

.avatar-wrapper:hover .avatar-overlay {
  opacity: 1;
}

/* 表单区域 */
.form-section {
  background: white;
  margin: 0 12px;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.form-group {
  margin-bottom: 20px;
}

.form-group:last-child {
  margin-bottom: 0;
}

.form-label {
  display: block;
  font-size: 14px;
  color: #333;
  font-weight: 500;
  margin-bottom: 8px;
}

.form-input {
  width: 100%;
  padding: 12px 14px;
  border: 1px solid #e8e8e8;
  border-radius: 8px;
  font-size: 14px;
  transition: all 0.3s;
  box-sizing: border-box;
}

.form-input:focus {
  outline: none;
  border-color: #1890ff;
  box-shadow: 0 0 0 2px rgba(24, 144, 255, 0.1);
}

.form-input:disabled {
  background: #f5f5f5;
  color: #999;
  cursor: not-allowed;
}

.form-select {
  width: 100%;
}

.form-select :deep(.el-input__wrapper) {
  padding: 8px 14px;
  border-radius: 8px;
}

.form-textarea {
  width: 100%;
  padding: 12px 14px;
  border: 1px solid #e8e8e8;
  border-radius: 8px;
  font-size: 14px;
  resize: vertical;
  font-family: inherit;
  box-sizing: border-box;
}

.form-textarea:focus {
  outline: none;
  border-color: #1890ff;
  box-shadow: 0 0 0 2px rgba(24, 144, 255, 0.1);
}
</style>
