import request from '../utils/request'

export interface Tag {
  ID: number
  name: string
  slug?: string
}

// 获取标签列表
export function getTagList() {
  return request({
    url: '/tags',
    method: 'get'
  })
}

// 创建标签
export function createTag(data: { name: string; slug?: string }) {
  return request({
    url: '/tags',
    method: 'post',
    data
  })
}

// 删除标签
export function deleteTag(id: number) {
  return request({
    url: `/tags/${id}`,
    method: 'delete'
  })
}