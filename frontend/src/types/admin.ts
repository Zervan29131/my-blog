import type { PaginatedData } from './blog'

export type ArticleStatus = 'draft' | 'published'
export type CommentStatus = 'pending' | 'approved' | 'rejected'

export interface Administrator {
  id: number
  username: string
  created_at: string
  updated_at: string
}

export interface LoginResult {
  token: string
  expires_in: number
}

export interface DashboardStats {
  article_total: number
  article_published: number
  article_draft: number
  comment_pending: number
  comment_approved: number
}

export interface SiteSettings {
  id: number
  site_name: string
  site_short_name: string | null
  site_description: string
  title_suffix: string | null
  logo_url: string | null
  favicon_url: string | null
  default_share_image_url: string | null
  copyright_name: string
  start_year: number | null
  additional_text: string | null
  filing_number: string | null
  filing_url: string | null
  show_technology: boolean
  technology_text: string | null
  created_at: string
  updated_at: string
}

export type SiteSettingsInput = Omit<SiteSettings, 'id' | 'created_at' | 'updated_at'>

export interface AdminArticleSummary {
  id: number
  title: string
  slug: string
  summary: string
  status: ArticleStatus
  published_at: string | null
  created_at: string
  updated_at: string
}

export interface AdminArticle extends AdminArticleSummary {
  content: string
}

export interface ArticleInput {
  title: string
  slug: string
  summary: string
  content: string
  status: ArticleStatus
}

export interface AdminComment {
  id: number
  article: {
    id: number
    title: string
    slug: string
  }
  nickname: string
  email: string | null
  content: string
  status: CommentStatus
  created_at: string
  updated_at: string
}

export type AdminArticlePage = PaginatedData<AdminArticleSummary>
export type AdminCommentPage = PaginatedData<AdminComment>
