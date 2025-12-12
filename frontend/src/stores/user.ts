// stores/user.ts
import { defineStore } from 'pinia';
import api from '@/api/index.ts';

export interface UserInfo {
  id: number;
  username: string;
  avatarurl: string;
  token: string;
  dormid : number;
}

export const useUserStore = defineStore('user', {
  state: () => ({
    // 默认初始状态：尝试从 localStorage 读取 token 和 user data
    // 生产环境中通常只存储 token，然后用 token 去换取完整的 user info
    token: localStorage.getItem('token') || null,
    userInfo: JSON.parse(localStorage.getItem('userInfo') || 'null') as UserInfo | null,
    unreadCount: 0,
  }),
  
  getters: {
    // 快速判断用户是否登录
    isLoggedIn: (state) => !!state.token && !!state.userInfo,
    // 提供给组件使用的当前用户信息
    currentUser: (state) => state.userInfo,
    // 获取用户 ID
    currentUserId: (state) => state.userInfo?.id,
  },
  
  actions: {
    // 登录成功时调用
    setLogin(token: string, user: { id: number,  username: string, avatarurl: string  , dormid : number}) {
      this.token = token;
      this.userInfo = user as UserInfo;
      
      // 1. 将 Token 和 User Info 写入本地存储（持久化）
      localStorage.setItem('token', token);
      localStorage.setItem('userInfo', JSON.stringify(user));
    },

    // 退出登录时调用
    logout() {
      this.token = null;
      this.userInfo = null;
      
      // 2. 清除本地存储
      localStorage.removeItem('token');
      localStorage.removeItem('userInfo');
      // 刷新页面或重定向到登录页
    },

    async fetchUnreadCount() {
      if (!this.token) return;
      try {
        const res = await api.get('api/v1/message/unreadcount');
        if (res.data.code === 200) {
          this.unreadCount = res.data.data;
        }
      } catch (error) {
        console.error("轮询未读消息失败", error);
      }  
    },

    clearUnread() {
      this.unreadCount = 0;
    }

  },
});