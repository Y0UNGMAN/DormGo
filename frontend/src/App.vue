<template>
  <div id="app">
    <!-- 路由出口：所有页面通过这里渲染 -->
    <router-view v-slot="{ Component }">
      <keep-alive include="DormgoHome">
        <component :is="Component" />
      </keep-alive>
  </router-view>
  </div>
</template>

<script setup lang="ts">
// Vue 3 推荐使用 <script setup> 语法糖，无需显式导出组件
</script>

<style>
/* 全局样式：清除默认边距，统一字体 */
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
  font-family: "Microsoft YaHei", sans-serif;
}

#app {
  width: 100vw;
  height: 100vh;
  /* 使用 vw/vh 更符合现代布局 */
}
</style>