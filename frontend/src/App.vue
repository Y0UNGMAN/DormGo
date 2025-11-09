<template>
  <div id="app">
    <!-- 路由出口：所有页面通过这里渲染 -->
    <router-view/>
  </div>
</template>

<script>
export default {
  name: 'App'
}
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
  width: 100%;
  height: 100%;
}
</style>
