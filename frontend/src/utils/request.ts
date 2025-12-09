// src/utils/request.ts
import axios from 'axios'
import type { AxiosError, InternalAxiosRequestConfig } from 'axios' // 导入内部请求配置类型
import { ElMessage } from 'element-plus'
import { useRouter } from 'vue-router'

const router = useRouter()

const request = axios.create({
    baseURL: import.meta.env.VITE_API_URL || '/api',
    timeout: 10000,
    headers: {
        'Content-Type': 'application/json'
    }
})

// 请求拦截器：使用 InternalAxiosRequestConfig 作为参数类型
request.interceptors.request.use(
    (config: InternalAxiosRequestConfig) => {
        const token = localStorage.getItem('adminToken')
        if (token) {
            config.headers.Authorization = `Bearer ${token}`
        }
        return config
    },
    (error: AxiosError) => {
        ElMessage.error('请求发送失败，请重试')
        return Promise.reject(error)
    }
)

// 响应拦截器保持不变
request.interceptors.response.use(
    (response) => {
        return response.data
    },
    (error: AxiosError) => {
        if (error.response) {
            if (error.response.status === 401) {
                ElMessage.error('登录已过期，请重新登录')
                localStorage.removeItem('adminToken')
                localStorage.removeItem('adminInfo')
                router.push('/login')
            } else if (error.response.status === 404) {
                ElMessage.error('请求的接口不存在')
            } else {
                const errorMsg = (error.response.data as { message?: string })?.message || '服务器错误'
                ElMessage.error(errorMsg)
            }
        } else if (error.message.includes('Network Error')) {
            ElMessage.error('网络连接失败，请检查服务是否启动')
        }
        return Promise.reject(error)
    }
)

export default request