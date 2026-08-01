import { http } from './http'
import type { ApiResponse } from '../types/blog'
import type {
  AdminArticle,
  AdminArticlePage,
  AdminCommentPage,
  Administrator,
  ArticleInput,
  CommentStatus,
  DashboardStats,
  LoginResult,
  SiteSettings,
  SiteSettingsInput,
} from '../types/admin'

export async function loginAdministrator(username: string, password: string): Promise<LoginResult> {
  const response = await http.post<ApiResponse<LoginResult>>('/admin/auth/login', {
    username,
    password,
  })
  return response.data.data
}

export async function fetchCurrentAdministrator(): Promise<Administrator> {
  const response = await http.get<ApiResponse<Administrator>>('/admin/auth/me')
  return response.data.data
}

export async function fetchDashboardStats(): Promise<DashboardStats> {
  const response = await http.get<ApiResponse<DashboardStats>>('/admin/dashboard')
  return response.data.data
}

export async function fetchSiteSettings(): Promise<SiteSettings> {
  const response = await http.get<ApiResponse<SiteSettings>>('/admin/site/settings')
  return response.data.data
}

export async function updateSiteSettings(input: SiteSettingsInput): Promise<SiteSettings> {
  const response = await http.put<ApiResponse<SiteSettings>>('/admin/site/settings', input)
  return response.data.data
}

export async function fetchAdminArticles(page = 1, pageSize = 10): Promise<AdminArticlePage> {
  const response = await http.get<ApiResponse<AdminArticlePage>>('/admin/articles', {
    params: { page, page_size: pageSize },
  })
  return response.data.data
}

export async function fetchAdminArticle(id: number): Promise<AdminArticle> {
  const response = await http.get<ApiResponse<AdminArticle>>(`/admin/articles/${id}`)
  return response.data.data
}

export async function createArticle(input: ArticleInput): Promise<AdminArticle> {
  const response = await http.post<ApiResponse<AdminArticle>>('/admin/articles', input)
  return response.data.data
}

export async function updateArticle(id: number, input: ArticleInput): Promise<AdminArticle> {
  const response = await http.put<ApiResponse<AdminArticle>>(`/admin/articles/${id}`, input)
  return response.data.data
}

export async function deleteArticle(id: number): Promise<void> {
  await http.delete(`/admin/articles/${id}`)
}

export async function fetchAdminComments(
  status: CommentStatus | '',
  page = 1,
  pageSize = 20,
): Promise<AdminCommentPage> {
  const response = await http.get<ApiResponse<AdminCommentPage>>('/admin/comments', {
    params: { status: status || undefined, page, page_size: pageSize },
  })
  return response.data.data
}

export async function updateCommentStatus(id: number, status: CommentStatus): Promise<void> {
  await http.put(`/admin/comments/${id}/status`, { status })
}

export async function deleteComment(id: number): Promise<void> {
  await http.delete(`/admin/comments/${id}`)
}
