import { http } from './http'
import type {
  ApiResponse,
  Article,
  ArticleSummary,
  Comment,
  CommentInput,
  PaginatedData,
  PublicHomepageConfig,
  PublicSiteConfig,
  SubmittedComment,
} from '../types/blog'

export async function fetchArticles(page = 1, pageSize = 10): Promise<PaginatedData<ArticleSummary>> {
  const response = await http.get<ApiResponse<PaginatedData<ArticleSummary>>>('/articles', {
    params: { page, page_size: pageSize },
  })
  return response.data.data
}

export async function fetchPublicSiteConfig(): Promise<PublicSiteConfig> {
  const response = await http.get<ApiResponse<PublicSiteConfig>>('/site/config')
  return response.data.data
}

export async function fetchHomepageConfig(): Promise<PublicHomepageConfig> {
  const response = await http.get<ApiResponse<PublicHomepageConfig>>('/homepage')
  return response.data.data
}

export async function fetchArticle(slug: string): Promise<Article> {
  const response = await http.get<ApiResponse<Article>>(`/articles/${encodeURIComponent(slug)}`)
  return response.data.data
}

export async function fetchComments(
  slug: string,
  page = 1,
  pageSize = 20,
): Promise<PaginatedData<Comment>> {
  const response = await http.get<ApiResponse<PaginatedData<Comment>>>(
    `/articles/${encodeURIComponent(slug)}/comments`,
    { params: { page, page_size: pageSize } },
  )
  return response.data.data
}

export async function submitComment(
  slug: string,
  input: CommentInput,
): Promise<ApiResponse<SubmittedComment>> {
  const response = await http.post<ApiResponse<SubmittedComment>>(
    `/articles/${encodeURIComponent(slug)}/comments`,
    input,
  )
  return response.data
}
