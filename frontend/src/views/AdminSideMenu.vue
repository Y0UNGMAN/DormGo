<template>
  <el-menu
    :default-active="activeMenu"
    class="el-menu-vertical"
    :collapse="collapse"
    background-color="#304156"
    text-color="#bfcbd9"
    active-text-color="#42b983"
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
      <el-menu-item index="/admin/config">通用配置</el-menu-item>
      <el-menu-item index="/admin/sensitive-words">敏感词管理</el-menu-item>
      <el-menu-item index="/admin/post-types">互助类型</el-menu-item>
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

// 计算当前激活菜单，防止子路由不高亮
const activeMenu = computed(() => {
  const { path } = route
  return path
})
</script>

<style scoped>
.el-menu-vertical {
  border-right: none;
  min-height: 100%;
}
.el-menu-vertical:not(.el-menu--collapse) {
  width: 220px;
}
</style>