<template>
  <div class="profile-container">
    <div class="profile-header">
      <h2>个人中心</h2>
      <button 
        class="save-btn" 
        @click="handleSave" 
        :disabled="!isEditMode"
      >
        保存修改
      </button>
    </div>

    <!-- 个人资料表单 -->
    <div class="profile-form">
      <!-- 头像上传区域 -->
      <div class="form-group avatar-group">
        <label>管理员头像</label>
        <div class="avatar-upload">
          <img 
            :src="profileForm.avatar || 'https://via.placeholder.com/100'" 
            alt="管理员头像" 
            class="avatar-img"
          >
          <input 
            type="file" 
            class="avatar-input" 
            accept="image/*" 
            @change="handleAvatarUpload"
            :disabled="!isEditMode"
          >
          <button 
            class="edit-avatar-btn" 
            @click="toggleEditMode(true)"
            v-if="!isEditMode"
          >
            更换头像
          </button>
        </div>
      </div>

      <!-- 基础信息区域 -->
      <div class="form-row">
        <div class="form-group">
          <label>管理员昵称 <span class="required">*</span></label>
          <input 
            type="text" 
            v-model="profileForm.nickname" 
            :disabled="!isEditMode"
            @blur="validateNickname"
          >
          <div class="error-tip" v-if="errors.nickname">{{ errors.nickname }}</div>
        </div>
        <div class="form-group">
          <label>所属楼栋（不可修改）</label>
          <input 
            type="text" 
            v-model="profileForm.dormName" 
            disabled
          >
        </div>
      </div>


      <div class="form-row">
        <div class="form-group">
          <label>手机号</label>
          <input 
            type="tel" 
            v-model="profileForm.phone" 
            :disabled="!isEditMode"
            @blur="validatePhone"
          >
          <div class="error-tip" v-if="errors.phone">{{ errors.phone }}</div>
        </div>
      </div>

      <div class="form-group">
        <label>个人简介</label>
        <textarea 
          v-model="profileForm.intro" 
          :disabled="!isEditMode"
          rows="3"
          placeholder="请输入个人简介（可选）"
        ></textarea>
      </div>

      <!-- 操作按钮区域 -->
      <div class="form-actions" v-if="!isEditMode">
        <button class="edit-btn" @click="toggleEditMode(true)">编辑资料</button>
        <button class="reset-pwd-btn" @click="goToResetPwd">重置密码</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'
import { ElMessage } from 'element-plus' // Vue3中常用element-plus的消息组件

// 响应式数据
const profileForm = ref({
  studentId: '',
  dormName: '',
  nickname: '',
  phone: '',
  intro: '',
  avatar: ''
})
const isEditMode = ref(false)
const errors = ref({
  nickname: '',
  phone: ''
})
const router = useRouter()

// 页面加载时获取管理员资料
onMounted(() => {
  fetchAdminProfile()
})

// 获取管理员个人资料
const fetchAdminProfile = async () => {
  try {
    const token = localStorage.getItem('adminToken')
    const response = await axios.get('/api/v1/user/profile', {
      headers: {
        Authorization: `Bearer ${token}`
      }
    })
    const { student_id, dorm_name, nickname, phone_number, intro, avatar } = response.data
    profileForm.value = {
      studentId: student_id,
      dormName: dorm_name,
      nickname: nickname,
      phone: phone_number || '',
      intro: intro || '',
      avatar: avatar || ''
    }
  } catch (error) {
    if (error.response && error.response.data && error.response.data.error) {
      ElMessage.error(error.response.data.error.message || '获取个人资料失败')
    } else {
      ElMessage.error('获取个人资料失败，请刷新重试')
    }
    console.error('Profile fetch error:', error)
  }
}

// 切换编辑模式
const toggleEditMode = (flag) => {
  isEditMode.value = flag
  if (!flag) {
    errors.value = { nickname: '', phone: '' }
  }
}

// 头像上传处理
const handleAvatarUpload = (e) => {
  const file = e.target.files[0]
  if (file) {
    if (file.size > 2 * 1024 * 1024) {
      ElMessage.warning('头像文件大小不能超过2MB')
      return
    }
    const reader = new FileReader()
    reader.onload = (res) => {
      profileForm.value.avatar = res.target.result
    }
    reader.readAsDataURL(file)
  }
}

// 昵称验证
const validateNickname = () => {
  if (!profileForm.value.nickname.trim()) {
    errors.value.nickname = '昵称不能为空'
  } else if (profileForm.value.nickname.length > 20) {
    errors.value.nickname = '昵称长度不能超过20个字符'
  } else {
    errors.value.nickname = ''
  }
}

// 手机号验证
const validatePhone = () => {
  const phone = profileForm.value.phone.trim()
  if (!phone) {
    errors.value.phone = ''
    return
  }
  const reg = /^1[3-9]\d{9}$/
  if (!reg.test(phone)) {
    errors.value.phone = '请输入正确的手机号格式'
  } else {
    errors.value.phone = ''
  }
}

// 保存修改
const handleSave = async () => {
  validateNickname()
  validatePhone()
  
  if (errors.value.nickname || errors.value.phone) {
    return
  }

  try {
    const token = localStorage.getItem('adminToken')
    await axios.put('/api/v1/user/profile', {
      nickname: profileForm.value.nickname.trim(),
      phone_number: profileForm.value.phone.trim(),
      intro: profileForm.value.intro.trim(),
      avatar: profileForm.value.avatar
    }, {
      headers: {
        Authorization: `Bearer ${token}`
      }
    })
    
    // 更新本地存储
    const adminInfo = JSON.parse(localStorage.getItem('adminInfo'))
    adminInfo.nickname = profileForm.value.nickname
    adminInfo.avatar = profileForm.value.avatar
    localStorage.setItem('adminInfo', JSON.stringify(adminInfo))
    
    toggleEditMode(false)
    ElMessage.success('个人资料修改成功')
  } catch (error) {
    ElMessage.error('保存失败，请重试')
    console.error('Profile save error:', error)
  }
}

// 跳转到重置密码页面
const goToResetPwd = () => {
  router.push('/admin/reset-pwd')
}
</script>

<style scoped>
.profile-container {
  padding: 20px 30px;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  margin: 20px;
}

.profile-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
  padding-bottom: 15px;
  border-bottom: 1px solid #f0f0f0;
}

.profile-header h2 {
  font-size: 18px;
  color: #333;
  font-weight: 500;
}

.save-btn {
  padding: 8px 16px;
  background: #42b983;
  color: #fff;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
}

.save-btn:disabled {
  background: #a3e9c4;
  cursor: not-allowed;
}

/* 表单样式 */
.profile-form {
  max-width: 800px;
}

.form-row {
  display: flex;
  gap: 30px;
  margin-bottom: 25px;
}

.form-group {
  flex: 1;
  margin-bottom: 25px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  color: #666;
  font-size: 14px;
}

.form-group input,
.form-group textarea {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
  box-sizing: border-box;
}

.form-group input:disabled,
.form-group textarea:disabled {
  background: #f5f5f5;
  color: #999;
  cursor: not-allowed;
}

.form-group textarea {
  resize: vertical;
}

.error-tip {
  margin-top: 5px;
  color: #f56c6c;
  font-size: 12px;
}

.required {
  color: #f56c6c;
}

/* 头像上传样式 */
.avatar-group {
  margin-bottom: 30px;
}

.avatar-upload {
  position: relative;
  display: inline-block;
}

.avatar-img {
  width: 100px;
  height: 100px;
  border-radius: 50%;
  object-fit: cover;
  border: 2px solid #eee;
}

.avatar-input {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  opacity: 0;
  cursor: pointer;
}

.edit-avatar-btn {
  margin-top: 10px;
  padding: 6px 12px;
  background: #f5f5f5;
  color: #333;
  border: 1px solid #ddd;
  border-radius: 4px;
  cursor: pointer;
  font-size: 13px;
}

/* 操作按钮样式 */
.form-actions {
  margin-top: 40px;
}

.edit-btn {
  padding: 8px 16px;
  background: #42b983;
  color: #fff;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  margin-right: 15px;
}

.reset-pwd-btn {
  padding: 8px 16px;
  background: #fff;
  color: #42b983;
  border: 1px solid #42b983;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
}
</style>