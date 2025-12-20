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
                <el-form-item label="头像设置">
                  <div class="upload-wrapper">
                    <el-upload
                      class="avatar-uploader"
                      action="#"
                      :auto-upload="false"
                      :show-file-list="false"
                      :on-change="handleFileChange"
                      :disabled="!isEditMode"
                    >
                      <img v-if="profileForm.avatar" :src="profileForm.avatar" class="avatar-preview" />
                      <el-icon v-else class="avatar-uploader-icon"><Plus /></el-icon>
                      
                      <template #tip v-if="isEditMode">
                        <div class="el-upload__tip">点击图片更换头像，支持 JPG/PNG 格式</div>
                      </template>
                    </el-upload>
                  </div>
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
import request from '@/utils/request'
import { ElMessage } from 'element-plus'
import { User, OfficeBuilding, Edit, Plus } from '@element-plus/icons-vue'

const router = useRouter()
const defaultAvatar = 'https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png'
const isEditMode = ref(false)
const avatarFile = ref(null) // 暂存上传的文件

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
    const res = await request.get('/api/v1/admin/profile')
    const data = res.data
    profileForm.studentId = data.student_id
    profileForm.dormName = data.dorm_name
    profileForm.nickname = data.nickname
    profileForm.intro = data.intro || ''
    profileForm.avatar = data.avatar || defaultAvatar
  } catch (error) {
    console.error(error)
  }
}

const toggleEditMode = (val) => {
  isEditMode.value = val
  if (!val) {
    fetchAdminProfile() // 取消时还原数据
    avatarFile.value = null
  }
}

// 处理文件选择，本地预览
const handleFileChange = (file) => {
  avatarFile.value = file.raw
  profileForm.avatar = URL.createObjectURL(file.raw)
}

const handleSave = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      try {
        // 使用 FormData 构建 multipart/form-data 请求
        const formData = new FormData()
        formData.append('nickname', profileForm.nickname)
        formData.append('intro', profileForm.intro)
        
        if (avatarFile.value) {
          formData.append('avatar', avatarFile.value)
        }

        await request.put('/api/v1/admin/profile', formData, {
          headers: {
            'Content-Type': 'multipart/form-data'
          }
        })
        
        // 更新本地存储显示的头像和名字
        const adminInfo = JSON.parse(localStorage.getItem('adminInfo') || '{}')
        adminInfo.nickname = profileForm.nickname
        // 如果后端返回了新头像URL最好，这里暂时用本地的或者不更新localStorage的avatar
        // 实际上重新拉取profile即可
        localStorage.setItem('adminInfo', JSON.stringify(adminInfo))
        
        ElMessage.success('保存成功')
        isEditMode.value = false
        avatarFile.value = null
        fetchAdminProfile() // 重新获取最新数据
      } catch (error) {
        console.error(error)
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

/* 头像上传样式 */
.avatar-uploader {
  border: 1px dashed #d9d9d9;
  border-radius: 6px;
  cursor: pointer;
  position: relative;
  overflow: hidden;
  width: 100px;
  height: 100px;
  display: flex;
  justify-content: center;
  align-items: center;
  transition: border-color 0.3s;
}
.avatar-uploader:hover {
  border-color: #409EFF;
}
.avatar-uploader-icon {
  font-size: 28px;
  color: #8c939d;
  width: 100px;
  height: 100px;
  line-height: 100px;
  text-align: center;
}
.avatar-preview {
  width: 100px;
  height: 100px;
  display: block;
  object-fit: cover;
}
</style>