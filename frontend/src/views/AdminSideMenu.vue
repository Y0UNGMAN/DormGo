<template>
  <el-menu
    :default-active="activeMenu"
    class="el-menu-vertical custom-menu"
    :collapse="collapse"
    router
    :collapse-transition="false"
  >
    <el-menu-item index="/admin/statistics">
      <el-icon><DataLine /></el-icon>
      <template #title>数据统计</template>
    </el-menu-item>

    <el-menu-item index="/admin/users">
      <el-icon><User /></el-icon>
      <template #title>用户管理</template>
    </el-menu-item>

    <el-menu-item index="/admin/dorms">
      <el-icon><OfficeBuilding /></el-icon>
      <template #title>楼栋管理</template>
    </el-menu-item>

    <el-menu-item index="/admin/audit">
      <el-icon><DocumentChecked /></el-icon>
      <template #title>内容审核</template>
    </el-menu-item>

    <el-menu-item index="/admin/notices">
      <el-icon><Bell /></el-icon>
      <template #title>通知管理</template>
    </el-menu-item>

    <el-menu-item index="/admin/violations">
      <el-icon><Warning /></el-icon>
      <template #title>违规处理</template>
    </el-menu-item>
    
    <el-sub-menu index="/admin/config-group">
      <template #title>
        <el-icon><Setting /></el-icon>
        <span>系统设置</span>
      </template>
      <el-menu-item index="/admin/sensitive-words">
        <template #title>敏感词管理</template>
      </el-menu-item>
      <el-menu-item index="/admin/post-types">
        <template #title>互助类型</template>
      </el-menu-item>
    </el-sub-menu>

    <el-menu-item index="/admin/profile">
      <el-icon><UserFilled /></el-icon>
      <template #title>个人中心</template>
    </el-menu-item>
  </el-menu>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { 
  User, UserFilled, DocumentChecked, Bell, 
  Warning, OfficeBuilding, DataLine, Setting
} from '@element-plus/icons-vue'

const route = useRoute()
defineProps({
  collapse: {
    type: Boolean,
    default: false
  }
})

// 计算当前激活菜单
const activeMenu = computed(() => {
  return route.path
})
</script>

<style scoped>
/* 菜单容器样式 */
.el-menu-vertical {
  border-right: none;
  min-height: 100vh;
  /* 使用更有质感的深色渐变背景 */
  background: linear-gradient(180deg, #1a1f3c 0%, #2c3e50 100%);
  user-select: none;
  transition: width 0.3s;
}

.el-menu-vertical:not(.el-menu--collapse) {
  width: 240px; /* 适当加宽，显得更大气 */
}

/* 菜单项基础样式 */
:deep(.el-menu-item), :deep(.el-sub-menu__title) {
  height: 54px;
  line-height: 54px;
  margin: 6px 10px; /* 增加外间距，实现悬浮卡片感 */
  border-radius: 10px; /* 圆角设计 */
  color: #a0aec0; /* 默认文字颜色为柔和的灰蓝 */
  transition: all 0.3s cubic-bezier(0.25, 0.8, 0.25, 1);
  border: 1px solid transparent; /* 预留边框位置避免抖动 */
}

/* 悬停效果 */
:deep(.el-menu-item:hover), :deep(.el-sub-menu__title:hover) {
  background-color: rgba(255, 255, 255, 0.08) !important;
  color: #fff !important;
  transform: translateX(4px); /* 悬浮时轻微右移 */
}

/* 图标样式 */
:deep(.el-icon) {
  font-size: 18px;
  margin-right: 12px;
  vertical-align: middle;
  transition: transform 0.3s;
}

/* 选中状态的高亮样式 - 核心美化点 */
:deep(.el-menu-item.is-active) {
  /* 使用品牌渐变色，与登录页风格统一 */
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #ffffff;
  box-shadow: 0 4px 15px rgba(118, 75, 162, 0.35); /* 增加发光投影 */
  font-weight: 600;
  border: none;
}

:deep(.el-menu-item.is-active .el-icon) {
  color: #fff;
  transform: scale(1.1); /* 选中时图标微放大 */
}

/* 子菜单展开背景微调 */
:deep(.el-sub-menu .el-menu) {
  background-color: transparent !important; /* 保持背景透明 */
}

:deep(.el-sub-menu .el-menu-item) {
  min-width: unset;
  margin: 4px 10px 4px 24px; /* 子菜单增加缩进 */
  height: 46px;
  line-height: 46px;
  font-size: 14px;
}

/* 折叠状态样式处理 */
:deep(.el-menu--collapse .el-menu-item), 
:deep(.el-menu--collapse .el-sub-menu__title) {
  margin: 4px 0;
  border-radius: 0;
  display: flex;
  justify-content: center;
  padding: 0 !important;
}

:deep(.el-menu--collapse .el-icon) {
  margin: 0;
}
</style>