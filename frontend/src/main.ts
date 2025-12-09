import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus' // 引入ElementPlus
import 'element-plus/dist/index.css' // 引入ElementPlus样式

// 创建应用实例
const app = createApp(App)

// 安装插件
app.use(router)
app.use(createPinia())
app.use(ElementPlus) // 注册ElementPlus

// 挂载应用
app.mount('#app')