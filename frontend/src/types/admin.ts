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

export type HomepageModuleType =
  | 'hero'
  | 'about'
  | 'featured_articles'
  | 'latest_articles'
  | 'tech_stack'
  | 'social_links'
export type HomepageLinkType = 'internal' | 'external'
export type HeroLayout = 'left' | 'center'
export type AboutImagePosition = 'left' | 'right' | 'none'

export interface HomepageButton {
  enabled: boolean
  text: string
  url: string
  link_type: HomepageLinkType
  open_in_new_tab: boolean
}

export interface HeroModuleConfig {
  eyebrow: string
  title: string
  highlight_text: string
  description: string
  image_url: string
  background_image_url: string
  layout: HeroLayout
  primary_button: HomepageButton
  secondary_button: HomepageButton
}

export interface AboutModuleConfig {
  title: string
  description: string
  content: string
  image_url: string
  image_position: AboutImagePosition
}

export interface FeaturedArticlesModuleConfig {
  title: string
  description: string
  limit: number
}

export interface LatestArticlesModuleConfig {
  title: string
  description: string
  limit: number
  show_summary: boolean
  show_date: boolean
  show_comment_count: boolean
  show_view_all: boolean
}

export interface TechItem {
  name: string
  description: string
  icon_url: string
  url: string
  is_visible: boolean
  sort_order: number
}

export interface TechStackModuleConfig {
  title: string
  description: string
  items: TechItem[]
}

export interface SocialLinksModuleConfig {
  title: string
  description: string
}

export interface HomepageModuleConfigMap {
  hero: HeroModuleConfig
  about: AboutModuleConfig
  featured_articles: FeaturedArticlesModuleConfig
  latest_articles: LatestArticlesModuleConfig
  tech_stack: TechStackModuleConfig
  social_links: SocialLinksModuleConfig
}

export type HomepageModule<T extends HomepageModuleType = HomepageModuleType> = T extends HomepageModuleType
  ? {
      type: T
      enabled: boolean
      sort_order: number
      config: HomepageModuleConfigMap[T]
    }
  : never

export interface AdminHomepageConfig {
  status: 'draft' | 'published'
  version: number
  modules: HomepageModule[]
  updated_at: string
  published_at: string | null
}

export interface HomepagePublishResult {
  version: number
  published_at: string
}

export type SocialPlatform =
  | 'github'
  | 'email'
  | 'linkedin'
  | 'x'
  | 'weibo'
  | 'bilibili'
  | 'zhihu'
  | 'custom'

export interface AdminSocialLink {
  id: number
  platform: SocialPlatform
  display_name: string
  url: string
  is_visible: boolean
  sort_order: number
  created_at: string
  updated_at: string
}

export type SocialLinkInput = Omit<AdminSocialLink, 'id' | 'created_at' | 'updated_at'>

export interface AdminFeaturedArticle {
  article_id: number
  title: string
  slug: string
  summary: string
  status: ArticleStatus
  published_at: string | null
  sort_order: number
  is_visible: boolean
  created_at: string
  updated_at: string
}

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
