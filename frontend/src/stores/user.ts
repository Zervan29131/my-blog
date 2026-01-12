import { defineStore } from 'pinia'
import { login } from '../api/auth'

export const useUserStore = defineStore('user', {
  state: () => ({
    // 初始化时尝试从本地存储读取
    token: localStorage.getItem('token') || '',
    username: localStorage.getItem('username') || ''
  }),
  
  getters: {
    isLoggedIn: (state) => !!state.token
  },

  actions: {
    // 登录动作
    async loginAction(loginForm: any) {
      try {
        // 调用 API
        const res: any = await login(loginForm)
        
        // 保存状态到 Pinia
        this.token = res.token
        this.username = loginForm.username // 后端暂时没返回用户信息，先用请求参数里的
        
        // 持久化到 LocalStorage
        localStorage.setItem('token', res.token)
        localStorage.setItem('username', loginForm.username)
        
        return true
      } catch (error) {
        throw error
      }
    },

    // 登出动作
    logout() {
      this.token = ''
      this.username = ''
      localStorage.removeItem('token')
      localStorage.removeItem('username')
    }
  }
})