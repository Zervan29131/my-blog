import { createRouter, createWebHistory } from 'vue-router'

import PublicLayout from '../layouts/PublicLayout.vue'
import AdminLayout from '../layouts/AdminLayout.vue'

// 前台页面组件导入
// 使用懒加载 (Lazy Loading) 优化性能
const PublicHome = () => import('../views/public/Home.vue')
const PostDetail = () => import('../views/public/PostDetail.vue')
const About = () => import('../views/public/About.vue')
const CategoryPostList = () => import('../views/public/CategoryPostList.vue')
const SearchResults = () => import('../views/public/SearchResults.vue')
const TagPostList = () => import('../views/public/TagPostList.vue')
const Archives = () => import('../views/public/Archives.vue')
const Categories = () => import('../views/public/Categories.vue') // 🟢 新增
const Links = () => import('../views/public/Links.vue')           // 🟢 新增
const Donate = () => import('../views/public/Donate.vue')         // 🟢 新增
const NotFound = () => import('../views/public/NotFound.vue')     // 🟢 新增

// 后台页面组件导入
const AdminLogin = () => import('../views/admin/Login.vue')
const AdminDashboard = () => import('../views/admin/Dashboard.vue')
const AdminPostList = () => import('../views/admin/PostList.vue')
const AdminPostEdit = () => import('../views/admin/PostEdit.vue')
const AdminCategoryTag = () => import('../views/admin/CategoryTag.vue')

const routes = [
  // 1. 前台路由 (Public)
  {
    path: '/',
    component: PublicLayout,
    children: [
      { 
        path: '', 
        name: 'Home', 
        component: PublicHome 
      },
      { 
        path: 'search', 
        name: 'Search', 
        component: SearchResults 
      },
      { 
        path: 'archives', 
        name: 'Archives', 
        component: Archives 
      },
      { 
        path: 'categories', 
        name: 'Categories', 
        component: Categories 
      },
      { 
        path: 'links', 
        name: 'Links', 
        component: Links 
      },
      { 
        path: 'donate', 
        name: 'Donate', 
        component: Donate 
      },
      // 动态路由参数
      { 
        path: 'category/:id', 
        name: 'CategoryPostList', 
        component: CategoryPostList 
      },
      { 
        path: 'tag/:id', 
        name: 'TagPostList', 
        component: TagPostList 
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

  // 2. 登录页
  {
    path: '/login',
    name: 'Login',
    component: AdminLogin
  },

  // 3. 后台路由 (Admin) - 需要权限验证
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

  // 4. 404 错误页
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: NotFound
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  // 切换路由时自动滚动到顶部
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    } else {
      return { top: 0 }
    }
  }
})

// 简单的路由守卫 (检查 Token)
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