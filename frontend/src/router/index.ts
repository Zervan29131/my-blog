import { createRouter, createWebHistory, type RouterHistory } from 'vue-router'

import { useAuthStore } from '../stores/auth'

export function createAppRouter(history: RouterHistory = createWebHistory(import.meta.env.BASE_URL)) {
  const router = createRouter({
    history,
    routes: [
      {
        path: '/',
        component: () => import('../layouts/PublicLayout.vue'),
        children: [
          {
            path: '',
            name: 'home',
            component: () => import('../views/public/HomeView.vue'),
          },
          {
            path: 'articles/:slug',
            name: 'article-detail',
            component: () => import('../views/public/ArticleDetailView.vue'),
            props: true,
          },
          {
            path: 'archive',
            name: 'archive',
            component: () => import('../views/public/ArchiveView.vue'),
          },
          {
            path: 'about',
            name: 'about',
            component: () => import('../views/public/AboutView.vue'),
          },
          {
            path: 'preview/home',
            name: 'preview-home',
            component: () => import('../views/public/PreviewHomeView.vue'),
            meta: { requiresAuth: true },
          },
          {
            path: ':pathMatch(.*)*',
            name: 'not-found',
            component: () => import('../views/public/NotFoundView.vue'),
          },
        ],
      },
      {
        path: '/admin/login',
        name: 'admin-login',
        component: () => import('../views/admin/LoginView.vue'),
        meta: { guestOnly: true },
      },
      {
        path: '/admin',
        component: () => import('../layouts/AdminLayout.vue'),
        meta: { requiresAuth: true },
        children: [
          {
            path: '',
            name: 'admin-dashboard',
            component: () => import('../views/admin/DashboardView.vue'),
          },
          {
            path: 'articles',
            name: 'admin-articles',
            component: () => import('../views/admin/ArticleListView.vue'),
          },
          {
            path: 'articles/new',
            name: 'admin-article-new',
            component: () => import('../views/admin/ArticleEditorView.vue'),
          },
          {
            path: 'articles/:id/edit',
            name: 'admin-article-edit',
            component: () => import('../views/admin/ArticleEditorView.vue'),
            props: (route) => ({ articleId: Number(route.params.id) }),
          },
          {
            path: 'comments',
            name: 'admin-comments',
            component: () => import('../views/admin/CommentListView.vue'),
          },
          {
            path: 'site/settings',
            name: 'admin-site-settings',
            component: () => import('../views/admin/SiteSettingsView.vue'),
          },
          {
            path: 'homepage',
            name: 'admin-homepage',
            component: () => import('../views/admin/HomepageSettingsView.vue'),
          },
        ],
      },
    ],
    scrollBehavior(_to, _from, savedPosition) {
      return savedPosition ?? { top: 0 }
    },
  })

  router.afterEach((to) => {
    const titles: Record<string, string> = {
      home: '字里行间 | 技术与生活手记',
      archive: '归档 | 字里行间',
      about: '关于 | 字里行间',
      'not-found': '页面不存在 | 字里行间',
      'admin-login': '管理后台登录 | 字里行间',
    }
    if (typeof to.name === 'string' && to.name !== 'article-detail') {
      document.title = titles[to.name] || (to.name.startsWith('admin-') ? '内容管理 | 字里行间' : '字里行间')
      const description = document.querySelector<HTMLMetaElement>('meta[name="description"]')
      if (description) description.content = '记录开发、工程实践与个人思考的技术博客。'
    }
  })

  router.beforeEach((to) => {
    const authStore = useAuthStore()
    const requiresAuth = to.matched.some((record) => record.meta.requiresAuth)
    if (requiresAuth && !authStore.isAuthenticated) {
      return {
        name: 'admin-login',
        query: { redirect: to.fullPath },
      }
    }
    if (to.meta.guestOnly && authStore.isAuthenticated) {
      return { name: 'admin-dashboard' }
    }
    return true
  })

  return router
}

const router = createAppRouter()

export default router
