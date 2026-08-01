import axios from 'axios'

export const ADMIN_TOKEN_KEY = 'personal_blog_admin_token'

interface ApiErrorBody {
  error?: {
    code?: string
    message?: string
  }
}

export const http = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api/v1',
  timeout: 10_000,
  headers: {
    Accept: 'application/json',
  },
})

let unauthorizedHandler: (() => void) | null = null

export function setUnauthorizedHandler(handler: () => void): void {
  unauthorizedHandler = handler
}

http.interceptors.request.use((config) => {
  const token = localStorage.getItem(ADMIN_TOKEN_KEY)
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

http.interceptors.response.use(
  (response) => response,
  (error: unknown) => {
    if (axios.isAxiosError(error) && error.response?.status === 401) {
      localStorage.removeItem(ADMIN_TOKEN_KEY)
      unauthorizedHandler?.()
    }
    return Promise.reject(error)
  },
)

export function apiErrorMessage(error: unknown, fallback = '请求失败，请稍后重试。'): string {
  if (!axios.isAxiosError<ApiErrorBody>(error)) {
    return fallback
  }
  if (error.code === 'ECONNABORTED') {
    return '请求超时，请检查网络后重试。'
  }
  return error.response?.data?.error?.message || fallback
}
