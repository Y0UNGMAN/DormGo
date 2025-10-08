import axios from 'axios'

// 创建 axios 实例
const api = axios.create({
  baseURL: 'http://localhost:8080', // 后端地址
  timeout: 5000,
})

export default api
