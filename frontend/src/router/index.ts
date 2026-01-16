import { createRouter, createWebHistory } from 'vue-router'

// 1. 引入布局组件
import PublicLayout from '../layouts/PublicLayout.vue'
import AdminLayout from '../layouts/AdminLayout.vue'

// 2. 引入页面组件 (使用懒加载优化性能)
// --- 前台页面 ---
const PublicHome = () => import('../views/public/Home.vue')
const PostDetail = () => import('../views/public/PostDetail.vue')

// --- 后台页面 ---
const AdminLogin = () => import('../views/admin/Login.vue')
const AdminDashboard = () => import('../views/admin/Dashboard.vue')
const AdminPostList = () => import('../views/admin/PostList.vue')
const AdminPostEdit = () => import('../views/admin/PostEdit.vue')
// 注意：请确保你已经将 CategoryTag.vue 移动到了 views/admin/ 目录下
const AdminCategoryTag = () => import('../views/admin/CategoryTag.vue')

const routes = [
  // ============================================
  // 前台路由 (面向访客)
  // ============================================
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
        path: 'post/:id',
        name: 'PostDetail',
        component: PostDetail
      },
      {
        path: 'about',
        name: 'About',
        // 如果没有 About.vue，可以暂时重定向到首页或写个简单的临时组件
        component: PublicHome 
      }
    ]
  },

  // ============================================
  // 登录页 (独立路由)
  // ============================================
  {
    path: '/login',
    name: 'Login',
    component: AdminLogin
  },

  // ============================================
  // 后台路由 (面向管理员，需鉴权)
  // ============================================
  {
    path: '/admin',
    component: AdminLayout,
    meta: { requiresAuth: true }, // 整个 admin 路由组都需要登录
    children: [
      {
        path: '',
        redirect: '/admin/dashboard'
      },
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: AdminDashboard
      },
      // 文章列表管理
      {
        path: 'posts', 
        name: 'AdminPostList',
        component: AdminPostList
      },
      // 写文章
      {
        path: 'posts/create',
        name: 'CreatePost',
        component: AdminPostEdit
      },
      // 编辑文章
      {
        path: 'posts/:id/edit',
        name: 'EditPost',
        component: AdminPostEdit
      },
      // 🟢 新增: 分类与标签管理
      {
        path: 'categories',
        name: 'CategoryManage',
        component: AdminCategoryTag
      }
    ]
  },

  // ============================================
  // 404 处理 (可选)
  // ============================================
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

// 路由守卫
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  // 检查路由元数据
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