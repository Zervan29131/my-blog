import request from '../utils/request'

export interface Post {
  ID: number
  title: string
  content: string
  summary: string
  view_count: number
  created_at: string
  updated_at: string
  category_id: number
  category?: {
    ID: number
    name: string
  }
  tags?: Array<{
    ID: number
    name: string
  }>
  author?: {
    username: string
  }
}

export interface PostListParams {
  page: number
  page_size: number
  q?: string
  category_id?: number // 🟢 新增
}

// 获取文章列表
export function getPostList(params: PostListParams) {
  return request({
    url: '/posts',
    method: 'get',
    params
  })
}

// 获取单篇文章详情
export function getPost(id: number) {
  return request({
    url: `/posts/${id}`,
    method: 'get'
  })
}

// 创建文章
export function createPost(data: any) {
  return request({
    url: '/posts',
    method: 'post',
    data
  })
}

// 更新文章
export function updatePost(id: number, data: any) {
  return request({
    url: `/posts/${id}`,
    method: 'put',
    data
  })
}

// 删除文章
export function deletePost(id: number) {
  return request({
    url: `/posts/${id}`,
    method: 'delete'
  })
}