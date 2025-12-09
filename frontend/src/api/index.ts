import axios from 'axios';
import type { InternalAxiosRequestConfig, AxiosError  , AxiosInstance} from 'axios'; 

const api: AxiosInstance = axios.create({
  baseURL: 'http://127.0.0.1:8080', 
  timeout: 5000,
  headers: {
    'Content-Type': 'application/json',
  },
});
// ==========================================
// 2. 添加请求拦截器 (使用新的类型 InternalAxiosRequestConfig)
// ==========================================
api.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const token = localStorage.getItem('token');    
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error: AxiosError) => {
    return Promise.reject(error);
  }
);
export default api;