import axios from 'axios'
import type { AxiosError, InternalAxiosRequestConfig } from 'axios'
import { ElMessage } from 'element-plus'

// 创建 axios 实例
const request = axios.create({
    // 【修改点1】直接连接 Go 后端地址，不再通过 /api 代理转发
    baseURL: 'http://127.0.0.1:8080',
    timeout: 10000,
    // 【修改点2】移除默认 Content-Type，让 axios 自动根据数据类型判断（解决上传图片报错问题）
    // headers: { 'Content-Type': 'application/json' } 
})

// 请求拦截器
request.interceptors.request.use(
    (config: InternalAxiosRequestConfig) => {
        // 从 localStorage 获取 Token
        // 优先获取 userToken (前台用户)，如果没有再看 adminToken (后台管理)
        // 注意：这里逻辑可以根据当前路由判断，但为了简单，优先取 UserToken 适配前台页面
        const token = localStorage.getItem('userToken') || localStorage.getItem('adminToken')

        if (token) {
            // Go 后端中间件要求格式为 "Bearer token"
            config.headers.Authorization = `Bearer ${token}`
        }
        return config
    },
    (error: AxiosError) => {
        ElMessage.error('请求发送失败')
        return Promise.reject(error)
    }
)

// 响应拦截器
request.interceptors.response.use(
    (response) => {
        // 后端返回的结构通常是 { code: 200, data: ..., msg: ... }
        // 如果后端返回 code 非 200，也可以在这里统一抛错，但为了兼容性，我们直接返回 data
        return response.data
    },
    (error: AxiosError) => {
        if (error.response) {
            const status = error.response.status
            if (status === 401) {
                ElMessage.error('登录已过期，请重新登录')
                localStorage.clear() // 清空所有缓存
                window.location.href = '/login'
            } else if (status === 404) {
                ElMessage.error('接口不存在 (404)')
            } else if (status === 500) {
                ElMessage.error('服务器内部错误 (500)')
            } else {
                ElMessage.error(`请求错误: ${status}`)
            }
        } else {
            ElMessage.error('网络连接失败，请检查后端是否启动')
        }
        return Promise.reject(error)
    }
)

export default request