import request from '../utils/request'

// 定义文章的数据结构接口
export interface Post {
  ID: number
  title: string
  content: string
  summary: string
  view_count: number
  created_at: string
  updated_at: string
  category_id: number
  // 关联数据通常是可选的，取决于后端是否 Preload
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

// 列表查询参数接口
export interface PostListParams {
  page: number
  page_size: number
  q?: string // 🟢 新增搜索参数
}

// 1. 获取文章列表
export function getPostList(params: PostListParams) {
  return request({
    url: '/posts',
    method: 'get',
    params
  })
}

// 2. 获取单篇文章详情 (用于编辑回显)
export function getPost(id: number) {
  return request({
    url: `/posts/${id}`,
    method: 'get'
  })
}

// 3. 创建文章
export function createPost(data: any) {
  return request({
    url: '/posts',
    method: 'post',
    data
  })
}

// 4. 更新文章
export function updatePost(id: number, data: any) {
  return request({
    url: `/posts/${id}`,
    method: 'put',
    data
  })
}

// 5. 删除文章
export function deletePost(id: number) {
  return request({
    url: `/posts/${id}`,
    method: 'delete'
  })
}