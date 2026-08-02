import { http } from './http'
import type { ApiResponse, PublicHomepageConfig } from '../types/blog'
import type {
  AdminArticle,
  AdminArticlePage,
  AdminCommentPage,
  AdminFeaturedArticle,
  AdminHomepageConfig,
  AdminSocialLink,
  Administrator,
  ArticleInput,
  ArticleStatus,
  CommentStatus,
  DashboardStats,
  HomepageModule,
  HomepagePublishResult,
  LoginResult,
  SiteSettings,
  SiteSettingsInput,
  SocialLinkInput,
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

export async function fetchHomepageDraft(): Promise<AdminHomepageConfig> {
  const response = await http.get<ApiResponse<AdminHomepageConfig>>('/admin/homepage/draft')
  return response.data.data
}

export async function fetchHomepagePublished(): Promise<AdminHomepageConfig> {
  const response = await http.get<ApiResponse<AdminHomepageConfig>>('/admin/homepage/published')
  return response.data.data
}

export async function fetchHomepagePreview(): Promise<PublicHomepageConfig> {
  const response = await http.get<ApiResponse<PublicHomepageConfig>>('/admin/homepage/preview')
  return response.data.data
}

export async function saveHomepageDraft(modules: HomepageModule[]): Promise<AdminHomepageConfig> {
  const response = await http.put<ApiResponse<AdminHomepageConfig>>('/admin/homepage/draft', { modules })
  return response.data.data
}

export async function publishHomepage(): Promise<HomepagePublishResult> {
  const response = await http.post<ApiResponse<HomepagePublishResult>>('/admin/homepage/publish')
  return response.data.data
}

export async function resetHomepageDraft(): Promise<AdminHomepageConfig> {
  const response = await http.post<ApiResponse<AdminHomepageConfig>>('/admin/homepage/reset-draft')
  return response.data.data
}

export async function fetchAdminArticles(
  page = 1,
  pageSize = 10,
  status: ArticleStatus | '' = '',
): Promise<AdminArticlePage> {
  const response = await http.get<ApiResponse<AdminArticlePage>>('/admin/articles', {
    params: { page, page_size: pageSize, status: status || undefined },
  })
  return response.data.data
}

export async function fetchFeaturedArticles(): Promise<AdminFeaturedArticle[]> {
  const response = await http.get<ApiResponse<AdminFeaturedArticle[]>>('/admin/featured-articles')
  return response.data.data
}

export async function addFeaturedArticle(articleId: number, sortOrder: number): Promise<void> {
  await http.post('/admin/featured-articles', {
    article_id: articleId,
    sort_order: sortOrder,
    is_visible: true,
  })
}

export async function updateFeaturedArticleVisibility(articleId: number, isVisible: boolean): Promise<void> {
  await http.put(`/admin/featured-articles/${articleId}`, { is_visible: isVisible })
}

export async function deleteFeaturedArticle(articleId: number): Promise<void> {
  await http.delete(`/admin/featured-articles/${articleId}`)
}

export async function reorderFeaturedArticles(items: AdminFeaturedArticle[]): Promise<void> {
  await http.put('/admin/featured-articles/order', {
    items: items.map((item, index) => ({ article_id: item.article_id, sort_order: (index + 1) * 10 })),
  })
}

export async function fetchSocialLinks(): Promise<AdminSocialLink[]> {
  const response = await http.get<ApiResponse<AdminSocialLink[]>>('/admin/social-links')
  return response.data.data
}

export async function createSocialLink(input: SocialLinkInput): Promise<AdminSocialLink> {
  const response = await http.post<ApiResponse<AdminSocialLink>>('/admin/social-links', input)
  return response.data.data
}

export async function updateSocialLink(id: number, input: SocialLinkInput): Promise<AdminSocialLink> {
  const response = await http.put<ApiResponse<AdminSocialLink>>(`/admin/social-links/${id}`, input)
  return response.data.data
}

export async function deleteSocialLink(id: number): Promise<void> {
  await http.delete(`/admin/social-links/${id}`)
}

export async function reorderSocialLinks(items: AdminSocialLink[]): Promise<void> {
  await http.put('/admin/social-links/order', {
    items: items.map((item, index) => ({ id: item.id, sort_order: (index + 1) * 10 })),
  })
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
