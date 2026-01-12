import request from '../utils/request'

// 登录接口
export function login(data: any) {
  return request({
    url: '/login',
    method: 'post',
    data
  })
}

// 注册接口 (为了方便你测试，也加上)
export function register(data: any) {
  return request({
    url: '/register',
    method: 'post',
    data
  })
}