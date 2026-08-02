<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import ThemeToggle from '../components/ThemeToggle.vue'
import { useAuthStore } from '../stores/auth'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const activeMenu = computed(() => {
  if (route.path.startsWith('/admin/articles')) return '/admin/articles'
  if (route.path.startsWith('/admin/comments')) return '/admin/comments'
  if (route.path.startsWith('/admin/site/settings')) return '/admin/site/settings'
  if (route.path.startsWith('/admin/homepage')) return '/admin/homepage'
  return '/admin'
})

function logout() {
  authStore.logout()
  void router.replace({ name: 'admin-login' })
}

onMounted(async () => {
  if (!authStore.currentAdmin) {
    try {
      await authStore.fetchCurrentAdmin()
    } catch {
      // The shared 401 handler performs the redirect and token cleanup.
    }
  }
})
</script>

<template>
  <div class="admin-shell">
    <aside class="admin-sidebar">
      <RouterLink class="admin-brand" to="/admin">
        <span>阅</span>
        <div>
          <strong>字里行间</strong>
          <small>内容管理</small>
        </div>
      </RouterLink>

      <el-menu class="admin-menu" router :default-active="activeMenu">
        <el-menu-item index="/admin">概览</el-menu-item>
        <el-menu-item index="/admin/articles">文章</el-menu-item>
        <el-menu-item index="/admin/comments">评论</el-menu-item>
        <el-menu-item index="/admin/site/settings">站点设置</el-menu-item>
        <el-menu-item index="/admin/homepage">首页配置</el-menu-item>
      </el-menu>

      <RouterLink class="view-blog-link" to="/" target="_blank">查看博客 ↗</RouterLink>
    </aside>

    <div class="admin-main">
      <header class="admin-topbar">
        <div>
          <span class="admin-context">PERSONAL BLOG</span>
          <strong>{{ authStore.currentAdmin?.username || '管理员' }}</strong>
        </div>
        <div class="admin-topbar-actions">
          <ThemeToggle />
          <el-button text @click="logout">退出登录</el-button>
        </div>
      </header>

      <div class="admin-content">
        <RouterView />
      </div>
    </div>
  </div>
</template>
