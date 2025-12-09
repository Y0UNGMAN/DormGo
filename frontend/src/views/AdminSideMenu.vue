<template>
  <el-menu
    :default-active="route.path"
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

    <el-menu-item index="/admin/user-list">
      <el-icon><User /></el-icon>
      <template #title>用户管理</template>
    </el-menu-item>

    <el-menu-item index="/admin/dorm-manage">
      <el-icon><OfficeBuilding /></el-icon>
      <template #title>楼栋管理</template>
    </el-menu-item>

    <el-menu-item index="/admin/content-audit">
      <el-icon><DocumentChecked /></el-icon>
      <template #title>内容审核</template>
    </el-menu-item>

    <el-menu-item index="/admin/notice">
      <el-icon><Bell /></el-icon>
      <template #title>通知管理</template>
    </el-menu-item>

    <el-menu-item index="/admin/violation">
      <el-icon><Warning /></el-icon>
      <template #title>违规处理</template>
    </el-menu-item>

    <el-menu-item index="/admin/system-config">
      <el-icon><Setting /></el-icon>
      <template #title>系统配置</template>
    </el-menu-item>

    <el-menu-item index="/admin/profile">
      <el-icon><UserFilled /></el-icon>
      <template #title>个人中心</template>
    </el-menu-item>
  </el-menu>
</template>

<script setup>
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