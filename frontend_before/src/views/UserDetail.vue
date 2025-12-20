<template>
  <div class="page-container">
    <!-- 顶部导航 -->
    <div class="nav-bar">
      <button class="back-btn" @click="goBack">← 返回个人中心</button>
      <span class="page-title">详细资料</span>
      <button class="save-btn" @click="saveProfile" :disabled="!isModified">保存</button>
    </div>

    <!-- 头像区域 -->
    <div class="avatar-section">
      <div class="avatar-wrapper" @click="changeAvatar">
        <img :src="form.avatar || defaultAvatar" alt="头像" class="avatar" />
        <div class="avatar-overlay">
          <span>更换头像</span>
        </div>
      </div>
    </div>

    <!-- 表单区域 -->
    <div class="form-section">
      <div class="form-group">
        <label class="form-label">昵称</label>
        <input v-model="form.nickname" type="text" class="form-input" placeholder="请输入昵称" />
      </div>

      <div class="form-group">
        <label class="form-label">姓名</label>
        <input v-model="form.realName" type="text" class="form-input" placeholder="请输入真实姓名" />
      </div>

      <div class="form-group">
        <label class="form-label">学号</label>
        <input v-model="form.studentId" type="text" class="form-input" placeholder="请输入学号" disabled />
      </div>

      <div class="form-group">
        <label class="form-label">学院</label>
        <input v-model="form.college" type="text" class="form-input" placeholder="请输入学院" />
      </div>

      <div class="form-group">
        <label class="form-label">专业</label>
        <input v-model="form.major" type="text" class="form-input" placeholder="请输入专业" />
      </div>

      <div class="form-group">
        <label class="form-label">年级</label>
        <el-select v-model="form.grade" placeholder="请选择年级" class="form-select">
          <el-option label="大一" value="大一" />
          <el-option label="大二" value="大二" />
          <el-option label="大三" value="大三" />
          <el-option label="大四" value="大四" />
          <el-option label="研一" value="研一" />
          <el-option label="研二" value="研二" />
          <el-option label="研三" value="研三" />
        </el-select>
      </div>

      <div class="form-group">
        <label class="form-label">宿舍楼</label>
        <el-select v-model="form.dormBuilding" placeholder="请选择宿舍楼" class="form-select">
          <el-option v-for="i in 12" :key="i" :label="`榕园${i}栋`" :value="`榕园${i}栋`" />
        </el-select>
      </div>

      <div class="form-group">
        <label class="form-label">宿舍号</label>
        <input v-model="form.dormRoom" type="text" class="form-input" placeholder="如：501" />
      </div>

      <div class="form-group">
        <label class="form-label">手机号</label>
        <input v-model="form.phone" type="tel" class="form-input" placeholder="请输入手机号" />
      </div>

      <div class="form-group">
        <label class="form-label">邮箱</label>
        <input v-model="form.email" type="email" class="form-input" placeholder="请输入邮箱" />
      </div>

      <div class="form-group">
        <label class="form-label">个性签名</label>
        <textarea v-model="form.bio" class="form-textarea" placeholder="介绍一下自己吧..." rows="3"></textarea>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'

const router = useRouter()

const defaultAvatar = 'https://dorm-go.oss-cn-guangzhou.aliyuncs.com/avator/midnight.jpg'

const originalForm = ref({})
const form = ref({
  avatar: '',
  nickname: '',
  realName: '',
  studentId: '',
  college: '',
  major: '',
  grade: '',
  dormBuilding: '',
  dormRoom: '',
  phone: '',
  email: '',
  bio: ''
})

const isModified = computed(() => {
  return JSON.stringify(form.value) !== JSON.stringify(originalForm.value)
})

const goBack = () => {
  if (isModified.value) {
    if (confirm('有未保存的修改，确定要离开吗？')) {
      router.push('/profile')
    }
  } else {
    router.push('/profile')
  }
}

const changeAvatar = () => {
  // 实际应调用文件上传
  alert('头像更换功能开发中...')
}

const saveProfile = async () => {
  try {
    // 实际应调用 API 保存用户资料
    // await axios.put('/api/v1/user/profile', form.value)
    
    // 更新本地存储
    const userInfo = JSON.parse(localStorage.getItem('userInfo') || '{}')
    Object.assign(userInfo, form.value)
    localStorage.setItem('userInfo', JSON.stringify(userInfo))
    
    originalForm.value = { ...form.value }
    ElMessage.success('保存成功')
  } catch (error) {
    ElMessage.error('保存失败，请重试')
  }
}

const loadUserProfile = () => {
  const storedUser = localStorage.getItem('userInfo')
  if (storedUser) {
    const user = JSON.parse(storedUser)
    form.value = {
      avatar: user.avatar || '',
      nickname: user.nickname || '',
      realName: user.realName || user.name || '',
      studentId: user.studentId || user.student_id || '',
      college: user.college || '',
      major: user.major || '',
      grade: user.grade || '',
      dormBuilding: user.dormBuilding || '',
      dormRoom: user.dormRoom || '',
      phone: user.phone || '',
      email: user.email || '',
      bio: user.bio || ''
    }
    originalForm.value = { ...form.value }
  }
}

onMounted(() => {
  loadUserProfile()
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
