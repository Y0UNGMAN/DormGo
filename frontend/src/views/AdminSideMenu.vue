<template>
  <div class="sidebar-container">
    <!-- 菜单标题 -->
    <div class="sidebar-title">核心功能菜单</div>
    
    <!-- 菜单列表 -->
    <ul class="menu-list">
      <!-- 个人中心菜单 -->
      <li class="menu-item">
        <router-link 
          to="/admin/profile" 
          class="menu-link"
          :class="{ 'active': $route.path === '/admin/profile' }"
        >
          <i class="menu-icon">👤</i>
          <span class="menu-text">个人中心</span>
        </router-link>
      </li>  
      <!-- 用户管理 -->
      <li class="menu-item">
        <a href="/admin/user-list" class="menu-link">
          <i class="menu-icon">👥</i>
          <span class="menu-text">用户管理</span>
        </a>
      </li>
      
      <!-- 内容审核 -->
      <li class="menu-item">
        <a href="/admin/content-audit" class="menu-link">
          <i class="menu-icon">📝</i>
          <span class="menu-text">内容审核</span>
        </a>
      </li>
      
      <!-- 通知管理 -->
      <li class="menu-item">
        <a href="/admin/notice" class="menu-link">
          <i class="menu-icon">📢</i>
          <span class="menu-text">通知管理</span>
        </a>
      </li>
      
      <!-- 违规处理 -->
      <li class="menu-item">
        <a href="/admin/violation" class="menu-link">
          <i class="menu-icon">⚠️</i>
          <span class="menu-text">违规处理</span>
        </a>
      </li>
      
      <!-- 楼栋管理 -->
      <li class="menu-item">
        <a href="/admin/dorm-manage" class="menu-link">
          <i class="menu-icon">🏢</i>
          <span class="menu-text">楼栋管理</span>
        </a>
      </li>
      
      <!-- 数据统计 -->
      <li class="menu-item">
        <a href="/admin/statistics" class="menu-link">
          <i class="menu-icon">📊</i>
          <span class="menu-text">数据统计</span>
        </a>
      </li>
    </ul>
  </div>
</template>

<script setup>
// Vue3 setup语法糖，无需额外配置
</script>

<style scoped>
/* 侧边栏容器 */
.menu-link.active {
  background: rgba(66, 185, 131, 0.2);
  color: #42b983;
  border-left: 3px solid #42b983;
}
.sidebar-container {
  width: 200px;
  height: 100%;
  background: #374151;
  color: #fff;
  overflow-y: auto;
}

/* 菜单标题 */
.sidebar-title {
  padding: 16px 20px;
  font-size: 15px;
  font-weight: 500;
  border-bottom: 1px solid rgba(255,255,255,0.1);
}

/* 菜单列表 */
.menu-list {
  list-style: none;
  margin: 0;
  padding: 0;
}

/* 菜单项 */
.menu-item {
  border-bottom: 1px solid rgba(255,255,255,0.05);
}
.menu-item:hover {
  background: rgba(255,255,255,0.08);
}

/* 菜单链接 */
.menu-link {
  display: flex;
  align-items: center;
  padding: 14px 25px;
  color: #e5e7eb;
  font-size: 14px;
  text-decoration: none;
  transition: all 0.3s;
}
.menu-link:hover {
  color: #fff;
  background: rgba(255,255,255,0.1);
}

/* 菜单图标 */
.menu-icon {
  margin-right: 12px;
  font-size: 16px;
}

/* 菜单文本 */
.menu-text {
  flex: 1;
}
</style>