import request from '../utils/request'

export interface Category {
  ID: number
  name: string
}

// 获取分类列表 (Nav 组件需要用到)
export function getCategoryList() {
  return request({
    url: '/categories', // 确保后端有这个接口，如果没有请使用 /tags 或者确认你的后端路由
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