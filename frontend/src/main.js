import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import ElementPlus from 'element-plus' // 引入element-plus组件库
import 'element-plus/dist/index.css' // 引入element-plus样式

// 创建Vue应用实例
const app = createApp(App)

// 安装插件
app.use(router)
app.use(ElementPlus)

// 挂载应用
app.mount('#app')