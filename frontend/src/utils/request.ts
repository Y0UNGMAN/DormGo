// src/utils/request.ts
import axios from 'axios'
import type { AxiosError, InternalAxiosRequestConfig, AxiosResponse } from 'axios'
import { ElMessage } from 'element-plus'
import router from '@/router' // 确保这里引入的是你的路由实例
// 引入 Pinia store
import { useUserStore } from '@/stores/user'

// 创建 axios 实例
const service = axios.create({
    baseURL: import.meta.env.VITE_API_URL || 'http://127.0.0.1:8080', 
    timeout: 10000,
    headers: {
        'Content-Type': 'application/json'
    }
})

// === 请求拦截器 ===
service.interceptors.request.use(
    (config: InternalAxiosRequestConfig) => {
        // 在拦截器内部获取 Store，确保 Pinia 已经挂载
        const userStore = useUserStore()
        
        // 智能 Token 策略：
        // 1. 如果有 adminToken (管理员)，优先使用
        // 2. 如果没有，则使用 token (普通用户)
        // 3. 也可以根据 config.url 前缀判断（例如 /api/v1/admin 用 adminToken），但下面这种更简单通用
        const token = userStore.adminToken || userStore.token
        
        if (token) {
            config.headers.Authorization = `Bearer ${token}`
        }
        return config
    },
    (error: AxiosError) => {
        ElMessage.error('请求发送失败，请检查网络')
        return Promise.reject(error)
    }
)

// === 响应拦截器 ===
service.interceptors.response.use(
    (response: AxiosResponse) => {
        // 通常后端返回 { code: 200, data: ..., message: ... }
        // 这里可以根据 code 做进一步判断，比如 code !== 200 视为业务错误
        const res = response.data
        
        // 如果你的后端成功时不一定返回 code 200，或者直接返回数据，请根据实际情况调整
        // 假设标准返回格式是 { code: number, data: any, message: string }
        if (res.code && res.code !== 200) {
             ElMessage.error(res.message || '系统错误')
             // 可以在这里处理特定业务错误码
             return Promise.reject(new Error(res.message || 'Error'))
        }

        return res
    },
    (error: AxiosError) => {
        const userStore = useUserStore()
        
        if (error.response) {
            const status = error.response.status
            const data = error.response.data as any
            const message = data?.message || '服务器内部错误'

            switch (status) {
                case 401:
                    ElMessage.error('登录状态已过期，请重新登录')
                    // 调用 Pinia 的 logout 清理所有状态（含 admin 和 user）
                    // 假设你的 userStore 提供了清理所有状态的方法，或者分别调用
                    if(userStore.adminToken) userStore.adminLogout()
                    else userStore.logout()
                    
                    router.push('/login')
                    break
                case 403:
                    ElMessage.error('拒绝访问：您没有权限执行此操作')
                    break
                case 404:
                    ElMessage.error('请求的资源或接口不存在')
                    break
                case 500:
                    ElMessage.error(`服务器错误: ${message}`)
                    break
                default:
                    ElMessage.error(message)
            }
        } else if (error.message.includes('Network Error')) {
            ElMessage.error('网络连接失败，请检查后端服务是否启动')
        } else if (error.message.includes('timeout')) {
            ElMessage.error('请求超时，请稍后重试')
        }
        
        return Promise.reject(error)
    }
)

export default service