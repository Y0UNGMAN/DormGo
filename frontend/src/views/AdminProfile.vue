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
              <span class="label"><el-icon><OfficeBuilding /></el-icon> 所属部门</span>
              <span class="value">{{ profileForm.dormName || '管理中心' }}</span>
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
              <el-col :span="24">
                <el-form-item label="头像URL">
                  <el-input v-model="profileForm.avatar" :disabled="!isEditMode" placeholder="请输入图片链接" />
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
import request from '@/utils/request' // 使用封装的 request
import { ElMessage } from 'element-plus'
import { User, OfficeBuilding, Edit } from '@element-plus/icons-vue'

const router = useRouter()
const defaultAvatar = 'https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png'
const isEditMode = ref(false)

const profileForm = reactive({
  studentId: '',
  dormName: '',
  nickname: '',
  intro: '',
  avatar: ''
})

const rules = {
  nickname: [
    { required: true, message: '请输入昵称', trigger: 'blur' },
    { min: 2, max: 20, message: '长度在 2 到 20 个字符', trigger: 'blur' }
  ]
}

const formRef = ref(null)

onMounted(() => {
  fetchAdminProfile()
})

const fetchAdminProfile = async () => {
  try {
    // 修改为 admin 专属接口
    const res = await request.get('/api/v1/admin/profile')
    const data = res.data
    profileForm.studentId = data.student_id
    profileForm.dormName = data.dorm_name
    profileForm.nickname = data.nickname
    profileForm.intro = data.intro || ''
    profileForm.avatar = data.avatar || ''
  } catch (error) {
    console.error(error)
  }
}

const toggleEditMode = (val) => {
  isEditMode.value = val
  if (!val) fetchAdminProfile() 
}

const handleSave = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      try {
        await request.put('/api/v1/admin/profile', {
          nickname: profileForm.nickname,
          intro: profileForm.intro,
          avatar: profileForm.avatar
        })
        
        // 更新本地存储显示的头像和名字
        const adminInfo = JSON.parse(localStorage.getItem('adminInfo') || '{}')
        adminInfo.nickname = profileForm.nickname
        adminInfo.avatar = profileForm.avatar
        localStorage.setItem('adminInfo', JSON.stringify(adminInfo))
        
        ElMessage.success('保存成功')
        isEditMode.value = false
      } catch (error) {
        // request.ts 已处理错误提示
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
