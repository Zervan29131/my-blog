import { createRouter, createWebHistory } from 'vue-router'
import Login from '../views/Login.vue'
// 如果还没有 PostList 组件，可以先注释掉或者随便写个临时组件
// import PostList from '../views/PostList.vue' 

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: Login
  },
  {
    path: '/',
    name: 'Home',
    // 暂时先重定向到 login，或者展示一个简单的 Home
    // component: PostList
    component: () => import('../views/Login.vue') // 临时：首页也显示登录，后续改成 PostList
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router