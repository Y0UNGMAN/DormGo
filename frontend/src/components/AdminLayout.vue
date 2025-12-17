<template>
  <el-container class="admin-layout">
    <el-aside :width="isCollapse ? '64px' : '220px'" class="sidebar-container">
      <div class="sidebar-logo">
        <el-icon color="#42b983" size="24" style="margin-right: 8px"><School /></el-icon>
        <span v-show="!isCollapse">寝友Go 管理后台</span>
      </div>
      <AdminSideMenu :collapse="isCollapse" class="side-menu" />
    </el-aside>

    <el-container>
      <el-header class="admin-header">
        <div class="header-left">
          <el-icon class="trigger-icon" @click="toggleCollapse">
            <component :is="isCollapse ? Expand : Fold" />
          </el-icon>
          <el-breadcrumb separator="/">
            <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
            <el-breadcrumb-item v-if="route.meta.title">{{ route.meta.title }}</el-breadcrumb-item>
          </el-breadcrumb>
        </div>

        <div class="header-right">
          <el-dropdown @command="handleCommand" trigger="click">
            <div class="user-info">
              <el-avatar :size="32" :src="adminInfo.avatar || defaultAvatar" />
              <span class="user-name">{{ adminInfo.nickname || '管理员' }}</span>
              <el-icon><CaretBottom /></el-icon>
            </div>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="profile">个人资料</el-dropdown-item>
                <el-dropdown-item command="config">系统配置</el-dropdown-item>
                <el-dropdown-item command="password">重置密码</el-dropdown-item>
                <el-dropdown-item divided command="logout" style="color: #f56c6c;">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>

      <el-main class="admin-main">
        <router-view v-slot="{ Component }">
          <transition name="fade-transform" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </el-main>
      
      <el-footer class="admin-footer">
         © 2025 寝友Go 管理系统
      </el-footer>
    </el-container>
  </el-container>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessageBox, ElMessage } from 'element-plus'
import { Fold, Expand, CaretBottom, School } from '@element-plus/icons-vue'
import AdminSideMenu from '@/views/AdminSideMenu.vue' 

const router = useRouter()
const route = useRoute()

const isCollapse = ref(false)
const adminInfo = ref({})
const defaultAvatar = 'https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png'

onMounted(() => {
  const info = localStorage.getItem('adminInfo')
  if (info) {
    try {
      adminInfo.value = JSON.parse(info)
    } catch (e) {
      console.error('Failed to parse admin info', e)
    }
  }
})

const toggleCollapse = () => {
  isCollapse.value = !isCollapse.value
}

const handleCommand = (command) => {
  switch (command) {
    case 'profile':
      router.push('/admin/profile')
      break
    case 'config':
      router.push('/admin/system-config')
      break
    case 'password':
      router.push('/admin/reset-pwd')
      break
    case 'logout':
      handleLogout()
      break
  }
}

const handleLogout = () => {
  ElMessageBox.confirm('确定要退出登录吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    localStorage.removeItem('adminToken')
    localStorage.removeItem('adminInfo')
    ElMessage.success('已安全退出')
    router.push('/')
  }).catch(() => {})
}
</script>

<style scoped>
.admin-layout { height: 100vh; width: 100%; }
.sidebar-container { background-color: #304156; transition: width 0.3s; overflow-x: hidden; display: flex; flex-direction: column; }
.sidebar-logo { height: 60px; line-height: 60px; background: #2b3649; display: flex; align-items: center; justify-content: center; overflow: hidden; }
.sidebar-logo span { color: #fff; font-weight: 600; font-size: 16px; white-space: nowrap; }
.side-menu { flex: 1; border-right: none; }
.admin-header { background: #fff; height: 60px; display: flex; justify-content: space-between; align-items: center; box-shadow: 0 1px 4px rgba(0,21,41,0.08); padding: 0 20px; z-index: 9; }
.header-left { display: flex; align-items: center; }
.trigger-icon { font-size: 20px; cursor: pointer; margin-right: 20px; color: #606266; }
.header-right .user-info { display: flex; align-items: center; cursor: pointer; padding: 0 8px; }
.user-name { margin: 0 8px; font-size: 14px; color: #606266; }
.admin-main { background-color: #f0f2f5; padding: 20px; position: relative; }
.admin-footer { height: 40px; line-height: 40px; text-align: center; color: #999; font-size: 12px; background: #f0f2f5; }
/* 动画效果 */
.fade-transform-leave-active, .fade-transform-enter-active { transition: all 0.5s; }
.fade-transform-enter-from { opacity: 0; transform: translateX(-30px); }
.fade-transform-leave-to { opacity: 0; transform: translateX(30px); }
</style>