<template>
  <div class="admin-layout">
    <!-- 顶部导航栏 -->
    <header class="admin-header">
      <div class="header-left">
        <h1>寝友Go 管理员后台</h1>
      </div>
      <!-- 右上角个人中心 -->
      <div class="header-right">
        <div class="user-center" @click="toggleUserMenu">
          <img 
            src="https://via.placeholder.com/32" 
            alt="用户头像" 
            class="user-avatar"
          >
          <span class="user-name">{{ adminInfo.nickname || '管理员' }}</span>
          <i class="icon-down">▼</i>
        </div>
        <!-- 个人中心下拉菜单 -->
        <div class="user-menu" v-if="showUserMenu">
          <a @click="goToProfile">个人资料</a>
          <a @click="resetPwd">重置密码</a>
          <a @click="handleLogout" class="logout">退出登录</a>
        </div>
      </div>
    </header>

    <div class="admin-content">
      <!-- 左侧菜单（引入拆分的组件） -->
      <aside class="admin-sidebar">
        <AdminSideMenu />
      </aside>

      <!-- 中间内容区域 -->
       <main class="admin-main">
        <!-- 子路由出口 -->
        <router-view />
      </main>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { useRouter } from 'vue-router'
import AdminSideMenu from '../views/AdminSideMenu.vue'

// 响应式数据
const adminInfo = ref({})
const showUserMenu = ref(false)
const router = useRouter()

// 页面挂载时获取管理员信息
onMounted(() => {
  const info = JSON.parse(localStorage.getItem('adminInfo'))
  adminInfo.value = info || {}
})

// 切换个人中心菜单显示状态
const toggleUserMenu = () => {
  showUserMenu.value = !showUserMenu.value
}

// 点击页面其他区域关闭菜单
const closeUserMenu = (e) => {
  if (!e.target.closest('.user-center')) {
    showUserMenu.value = false
  }
}

// 个人资料页面跳转
const goToProfile = () => {
  router.push('/admin/profile')
  showUserMenu.value = false
}

// 重置密码页面跳转
const resetPwd = () => {
  router.push('/admin/reset-pwd')
  showUserMenu.value = false
}

// 退出登录
const handleLogout = () => {
  localStorage.removeItem('adminToken')
  localStorage.removeItem('adminInfo')
  router.push('/login')
}

// 监听点击事件
onMounted(() => {
  document.addEventListener('click', closeUserMenu)
})

// 组件卸载前移除事件监听
onBeforeUnmount(() => {
  document.removeEventListener('click', closeUserMenu)
})
</script>

<style scoped>
/* 布局整体样式 */
.admin-layout {
  width: 100%;
  height: 100vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

/* 顶部导航栏 */
.admin-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 30px;
  height: 60px;
  background: #2d3748;
  color: #fff;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  z-index: 10;
}
.header-left h1 {
  font-size: 18px;
  font-weight: 500;
}

/* 右上角个人中心 */
.header-right {
  display: flex;
  align-items: center;
}
.user-center {
  display: flex;
  align-items: center;
  cursor: pointer;
  padding: 5px 10px;
  border-radius: 4px;
  transition: background 0.3s;
}
.user-center:hover {
  background: rgba(255,255,255,0.1);
}
.user-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  margin-right: 8px;
  border: 2px solid rgba(255,255,255,0.3);
}
.user-name {
  font-size: 14px;
  margin-right: 8px;
}
.icon-down {
  font-size: 12px;
  opacity: 0.8;
}

/* 个人中心下拉菜单 */
.user-menu {
  position: absolute;
  top: 60px;
  right: 30px;
  width: 160px;
  background: #fff;
  border-radius: 4px;
  box-shadow: 0 4px 12px rgba(0,0,0,0.15);
  overflow: hidden;
  z-index: 100;
}
.user-menu a {
  display: block;
  padding: 12px 20px;
  color: #333;
  font-size: 14px;
  text-decoration: none;
  transition: background 0.3s;
}
.user-menu a:hover {
  background: #f5f5f5;
}
.user-menu .logout {
  color: #e53e3e;
  border-top: 1px solid #eee;
  margin-top: 5px;
}

/* 内容区域（侧边栏+主内容） */
.admin-content {
  display: flex;
  flex: 1;
  overflow: hidden;
}

/* 中间主内容区域 */
.admin-main {
  flex: 1;
  padding: 30px;
  background: #f8f9fa;
  overflow-y: auto;
}
.empty-content {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
  color: #666;
  font-size: 16px;
}
</style>