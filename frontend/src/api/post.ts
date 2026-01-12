import request from '../utils/request'

export interface Post {
  ID: number
  title: string
  content: string
  summary: string
  category: {
    name: string
  }
  created_at: string
}

export function getPostList(params: { page: number; page_size: number }) {
  return request({
    url: '/posts',
    method: 'get',
    params
  })
}

export function createPost(data: any) {
  return request({
    url: '/posts',
    method: 'post',
    data
  })
}
