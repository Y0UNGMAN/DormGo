<template>
  <div class="app-container">
    <el-row :gutter="20">
      <el-col :span="8" :xs="24">
        <el-card class="box-card" shadow="hover">
          <template #header>
            <div class="clearfix">
              <span>个人信息</span>
            </div>
          </template>
          <div class="profile-header">
            <div class="avatar-container">
              <el-avatar :size="100" :src="profileForm.avatar || defaultAvatar" class="user-avatar" />
              <div v-if="isEditMode" class="avatar-mask">
                <label for="avatar-input" class="upload-label">
                  <el-icon><Camera /></el-icon> 更换
                </label>
                <input 
                  id="avatar-input" 
                  type="file" 
                  accept="image/*" 
                  @change="handleAvatarUpload" 
                  style="display:none"
                >
              </div>
            </div>
            <div class="user-name">{{ profileForm.nickname || '管理员' }}</div>
            <div class="user-role">超级管理员</div>
          </div>
          
          <el-divider />
          
          <div class="profile-detail">
            <div class="detail-item">
              <span class="label"><el-icon><User /></el-icon> 账号ID</span>
              <span class="value">{{ profileForm.studentId }}</span>
            </div>
            <div class="detail-item">
              <span class="label"><el-icon><OfficeBuilding /></el-icon> 所属楼栋</span>
              <span class="value">{{ profileForm.dormName || '全校通用' }}</span>
            </div>
            <div class="detail-item">
              <span class="label"><el-icon><Timer /></el-icon> 上次登录</span>
              <span class="value">{{ lastLoginTime || '暂无记录' }}</span>
            </div>
          </div>
        </el-card>
      </el-col>

      <el-col :span="16" :xs="24">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>编辑资料</span>
              <div class="header-btns">
                <el-button v-if="!isEditMode" type="primary" link @click="toggleEditMode(true)">
                  <el-icon><Edit /></el-icon> 修改资料
                </el-button>
                <template v-else>
                  <el-button link @click="toggleEditMode(false)">取消</el-button>
                  <el-button type="primary" @click="handleSave">保存修改</el-button>
                </template>
              </div>
            </div>
          </template>

          <el-form 
            ref="formRef" 
            :model="profileForm" 
            :rules="rules" 
            label-position="top"
            class="profile-form"
          >
            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item label="管理员昵称" prop="nickname">
                  <el-input v-model="profileForm.nickname" :disabled="!isEditMode" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="手机号码" prop="phone">
                  <el-input v-model="profileForm.phone" :disabled="!isEditMode" />
                </el-form-item>
              </el-col>
              <el-col :span="24">
                <el-form-item label="个人简介">
                  <el-input 
                    v-model="profileForm.intro" 
                    type="textarea" 
                    :rows="4" 
                    :disabled="!isEditMode" 
                    placeholder="请输入个人简介..."
                  />
                </el-form-item>
              </el-col>
            </el-row>
          </el-form>

          <el-divider content-position="left">安全设置</el-divider>
          <div class="security-settings">
            <div class="security-item">
              <div class="security-text">
                <h4>登录密码</h4>
                <p>建议定期修改密码以保障账号安全</p>
              </div>
              <el-button type="primary" plain @click="goToResetPwd">修改密码</el-button>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted, reactive } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'
import { ElMessage } from 'element-plus'
import { User, OfficeBuilding, Timer, Camera, Edit } from '@element-plus/icons-vue'

const router = useRouter()
const defaultAvatar = 'https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png'
const isEditMode = ref(false)
const lastLoginTime = ref('')

const profileForm = reactive({
  studentId: '',
  dormName: '',
  nickname: '',
  phone: '',
  intro: '',
  avatar: ''
})

const rules = {
  nickname: [
    { required: true, message: '请输入昵称', trigger: 'blur' },
    { min: 2, max: 20, message: '长度在 2 到 20 个字符', trigger: 'blur' }
  ],
  phone: [
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号码', trigger: 'blur' }
  ]
}

const formRef = ref(null)

onMounted(() => {
  fetchAdminProfile()
  fetchLoginInfo()
})

const fetchAdminProfile = async () => {
  try {
    const token = localStorage.getItem('adminToken')
    const res = await axios.get('/api/v1/user/profile', {
      headers: { Authorization: `Bearer ${token}` }
    })
    const data = res.data
    profileForm.studentId = data.student_id
    profileForm.dormName = data.dorm_name
    profileForm.nickname = data.nickname
    profileForm.phone = data.phone_number || ''
    profileForm.intro = data.intro || ''
    profileForm.avatar = data.avatar || ''
  } catch (error) {
    ElMessage.error('获取资料失败')
  }
}

const fetchLoginInfo = async () => {
  try {
    const token = localStorage.getItem('adminToken')
    const res = await axios.get('/api/v1/user/login-info', {
      headers: { Authorization: `Bearer ${token}` }
    })
    if (res.data.last_login_time) {
      lastLoginTime.value = new Date(res.data.last_login_time).toLocaleString()
    }
  } catch (e) { console.error(e) }
}

const handleAvatarUpload = (e) => {
  const file = e.target.files[0]
  if (!file) return
  if (file.size > 2 * 1024 * 1024) {
    ElMessage.warning('头像大小不能超过2MB')
    return
  }
  const reader = new FileReader()
  reader.onload = (res) => {
    profileForm.avatar = res.target.result
  }
  reader.readAsDataURL(file)
}

const toggleEditMode = (val) => {
  isEditMode.value = val
  if (!val) fetchAdminProfile() // 取消则重置
}

const handleSave = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      try {
        const token = localStorage.getItem('adminToken')
        await axios.put('/api/v1/user/profile', {
          nickname: profileForm.nickname,
          phone_number: profileForm.phone,
          intro: profileForm.intro,
          avatar: profileForm.avatar
        }, { headers: { Authorization: `Bearer ${token}` } })
        
        // 更新本地存储
        const adminInfo = JSON.parse(localStorage.getItem('adminInfo') || '{}')
        adminInfo.nickname = profileForm.nickname
        adminInfo.avatar = profileForm.avatar
        localStorage.setItem('adminInfo', JSON.stringify(adminInfo))
        
        ElMessage.success('保存成功')
        isEditMode.value = false
      } catch (error) {
        ElMessage.error('保存失败')
      }
    }
  })
}

const goToResetPwd = () => router.push('/admin/reset-pwd')
</script>

<style scoped>
.app-container { padding: 20px; }
.profile-header { text-align: center; padding: 20px 0; }
.avatar-container { position: relative; display: inline-block; margin-bottom: 15px; }
.user-avatar { border: 2px solid #fff; box-shadow: 0 2px 12px rgba(0,0,0,0.1); }
.avatar-mask {
  position: absolute; top: 0; left: 0; width: 100%; height: 100%;
  background: rgba(0,0,0,0.5); border-radius: 50%;
  display: flex; align-items: center; justify-content: center;
  color: #fff; cursor: pointer; opacity: 0; transition: opacity 0.3s;
}
.avatar-container:hover .avatar-mask { opacity: 1; }
.upload-label { cursor: pointer; display: flex; align-items: center; gap: 4px; font-size: 12px; }

.user-name { font-size: 24px; font-weight: bold; color: #303133; margin-bottom: 5px; }
.user-role { font-size: 14px; color: #909399; }

.profile-detail .detail-item {
  display: flex; justify-content: space-between; padding: 12px 0;
  border-bottom: 1px solid #f0f2f5; font-size: 14px;
}
.detail-item:last-child { border-bottom: none; }
.detail-item .label { color: #606266; display: flex; align-items: center; gap: 5px; }
.detail-item .value { color: #303133; }

.card-header { display: flex; justify-content: space-between; align-items: center; }
.security-item { display: flex; justify-content: space-between; align-items: center; margin-top: 10px; }
.security-text h4 { margin: 0 0 5px 0; color: #303133; }
.security-text p { margin: 0; color: #909399; font-size: 13px; }
</style>