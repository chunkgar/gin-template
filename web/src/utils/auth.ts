import { ref } from 'vue'

// Token 管理
export const TOKEN_KEY = 'admin_token'

export const getToken = (): string | null => {
  return localStorage.getItem(TOKEN_KEY)
}

export const setToken = (token: string): void => {
  localStorage.setItem(TOKEN_KEY, token)
}

export const removeToken = (): void => {
  localStorage.removeItem(TOKEN_KEY)
}

// 用户信息管理
export interface AdminUser {
  id: number
  username: string
  role: string
  isActive: boolean
  createdAt: string
}

// 登录状态管理
export const isLoggedIn = ref(!!getToken())

export const checkLoginStatus = (): boolean => {
  const token = getToken()
  isLoggedIn.value = !!token
  return isLoggedIn.value
}

export const logout = (): void => {
  removeToken()
  isLoggedIn.value = false
  // 重定向到登录页面
  window.location.href = '/login'
}

// 自动检查登录状态
checkLoginStatus()