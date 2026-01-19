import { createRouter, createWebHistory } from 'vue-router'

import PublicLayout from '../layouts/PublicLayout.vue'
import AdminLayout from '../layouts/AdminLayout.vue'

// 前台页面
const PublicHome = () => import('../views/public/Home.vue')
const PostDetail = () => import('../views/public/PostDetail.vue')
const About = () => import('../views/public/About.vue')
// 🟢 新增：分类文章列表页
const CategoryPostList = () => import('../views/public/CategoryPostList.vue')

// 后台页面
const AdminLogin = () => import('../views/admin/Login.vue')
const AdminDashboard = () => import('../views/admin/Dashboard.vue')
const AdminPostList = () => import('../views/admin/PostList.vue')
const AdminPostEdit = () => import('../views/admin/PostEdit.vue')
const AdminCategoryTag = () => import('../views/admin/CategoryTag.vue')

const routes = [
  // 前台路由
  {
    path: '/',
    component: PublicLayout,
    children: [
      {
        path: '',
        name: 'Home',
        component: PublicHome
      },
      // 🟢 新增分类路由配置
      // 当你访问 /category/3 时，会渲染 CategoryPostList 组件
      {
        path: 'category/:id',
        name: 'CategoryPostList',
        component: CategoryPostList
      },
      {
        path: 'post/:id',
        name: 'PostDetail',
        component: PostDetail
      },
      {
        path: 'about',
        name: 'About',
        component: About 
      }
    ]
  },

  // 登录页
  {
    path: '/login',
    name: 'Login',
    component: AdminLogin
  },

  // 后台路由
  {
    path: '/admin',
    component: AdminLayout,
    meta: { requiresAuth: true },
    children: [
      { path: '', redirect: '/admin/dashboard' },
      { path: 'dashboard', name: 'Dashboard', component: AdminDashboard },
      { path: 'posts', name: 'AdminPostList', component: AdminPostList },
      { path: 'posts/create', name: 'CreatePost', component: AdminPostEdit },
      { path: 'posts/:id/edit', name: 'EditPost', component: AdminPostEdit },
      { path: 'categories', name: 'CategoryManage', component: AdminCategoryTag }
    ]
  },

  // 404
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    redirect: '/'
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 简单路由守卫
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  if (to.matched.some(record => record.meta.requiresAuth)) {
    if (!token) {
      next('/login')
    } else {
      next()
    }
  } else {
    next()
  }
})

export default router