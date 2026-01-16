import request from '../utils/request'

export interface Comment {
  ID: number
  post_id: number
  parent_id?: number | null
  nickname: string
  email?: string
  website?: string
  content: string
  created_at: string
  children?: Comment[] // 前端递归组装时使用
}

export function getComments(postId: number) {
  return request({
    url: '/comments',
    method: 'get',
    params: { post_id: postId }
  })
}

export function createComment(data: any) {
  return request({
    url: '/comments',
    method: 'post',
    data
  })
}