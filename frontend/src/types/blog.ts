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
