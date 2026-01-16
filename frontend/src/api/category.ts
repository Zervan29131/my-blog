import request from '../utils/request'

export interface Category {
  ID: number
  name: string
}

// 获取分类列表
export function getCategoryList() {
  return request({
    url: '/categories', // 假设后端有这个接口，如果没有需参考 Tag 逻辑补充
    method: 'get'
  })
}

// 创建分类
export function createCategory(data: { name: string }) {
  return request({
    url: '/categories',
    method: 'post',
    data
  })
}

// 删除分类
export function deleteCategory(id: number) {
  return request({
    url: `/categories/${id}`,
    method: 'delete'
  })
}