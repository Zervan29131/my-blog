export interface ApiResponse<T> {
  data: T
  message: string
}

export interface PaginatedData<T> {
  items: T[]
  page: number
  page_size: number
  total: number
  total_pages: number
}

export interface ArticleSummary {
  id: number
  title: string
  slug: string
  summary: string
  published_at: string
  comment_count: number
}

export interface Article {
  id: number
  title: string
  slug: string
  summary: string
  content: string
  published_at: string
  created_at: string
  updated_at: string
}

export interface Comment {
  id: number
  nickname: string
  content: string
  created_at: string
}

export interface CommentInput {
  nickname: string
  email?: string
  content: string
}

export interface SubmittedComment {
  id: number
  nickname: string
  content: string
  status: 'pending'
  created_at: string
}

export type PublicLinkType = 'internal' | 'external'
export type PublicHomepageModuleType =
  | 'hero'
  | 'about'
  | 'featured_articles'
  | 'latest_articles'
  | 'tech_stack'
  | 'social_links'

export interface PublicSite {
  name: string
  short_name: string
  description: string
  title_suffix: string
  logo_url: string
  favicon_url: string
  default_share_image_url: string
}

export interface PublicNavigationItem {
  name: string
  url: string
  link_type: PublicLinkType
  open_in_new_tab: boolean
}

export interface PublicSocialLink {
  platform: string
  display_name: string
  url: string
}

export interface PublicFooter {
  copyright_name: string
  start_year: number | null
  additional_text: string
  filing_number: string
  filing_url: string
  show_technology: boolean
  technology_text: string
}

export interface PublicSiteConfig {
  site: PublicSite
  navigation: PublicNavigationItem[]
  social_links: PublicSocialLink[]
  footer: PublicFooter
}

export interface PublicHomepageButton {
  enabled: boolean
  text: string
  url: string
  link_type: PublicLinkType
  open_in_new_tab: boolean
}

export interface PublicHeroConfig {
  eyebrow: string
  title: string
  highlight_text: string
  description: string
  image_url: string
  background_image_url: string
  layout: 'left' | 'center'
  primary_button: PublicHomepageButton
  secondary_button: PublicHomepageButton
}

export interface PublicAboutConfig {
  title: string
  description: string
  content: string
  image_url: string
  image_position: 'left' | 'right' | 'none'
}

export interface PublicArticlesModuleConfig {
  title: string
  description: string
  limit: number
  articles: ArticleSummary[]
}

export interface PublicLatestArticlesConfig extends PublicArticlesModuleConfig {
  show_summary: boolean
  show_date: boolean
  show_comment_count: boolean
  show_view_all: boolean
}

export interface PublicTechItem {
  name: string
  description: string
  icon_url: string
  url: string
}

export interface PublicTechStackConfig {
  title: string
  description: string
  items: PublicTechItem[]
}

export interface PublicSocialLinksConfig {
  title: string
  description: string
  links: PublicSocialLink[]
}

export interface PublicHomepageModuleConfigMap {
  hero: PublicHeroConfig
  about: PublicAboutConfig
  featured_articles: PublicArticlesModuleConfig
  latest_articles: PublicLatestArticlesConfig
  tech_stack: PublicTechStackConfig
  social_links: PublicSocialLinksConfig
}

export type PublicHomepageModule<T extends PublicHomepageModuleType = PublicHomepageModuleType> = T extends PublicHomepageModuleType
  ? { type: T; sort_order: number; config: PublicHomepageModuleConfigMap[T] }
  : never

export interface PublicHomepageConfig {
  version: number
  modules: PublicHomepageModule[]
}
