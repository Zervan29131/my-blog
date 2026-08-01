<script setup lang="ts">
import { onMounted, onUnmounted, ref, watch } from 'vue'
import { RouterLink, RouterView, useRoute } from 'vue-router'

import ThemeToggle from '../components/ThemeToggle.vue'

const route = useRoute()
const menuOpen = ref(false)
const menuButton = ref<HTMLButtonElement | null>(null)
const currentYear = new Date().getFullYear()

function closeMenu(restoreFocus = false) {
  menuOpen.value = false
  if (restoreFocus) menuButton.value?.focus()
}

function handleKeydown(event: KeyboardEvent) {
  if (event.key === 'Escape' && menuOpen.value) closeMenu(true)
}

watch(() => route.fullPath, () => closeMenu())
watch(menuOpen, (open) => {
  document.body.classList.toggle('nav-open', open)
})

onMounted(() => window.addEventListener('keydown', handleKeydown))
onUnmounted(() => {
  window.removeEventListener('keydown', handleKeydown)
  document.body.classList.remove('nav-open')
})
</script>

<template>
  <div class="site-shell">
    <header class="site-header">
      <div class="header-inner">
        <RouterLink class="brand" to="/" aria-label="字里行间博客首页">
          <span class="brand-mark" aria-hidden="true"><i></i></span>
          <strong>字里行间</strong>
        </RouterLink>

        <nav class="desktop-nav" aria-label="主要导航">
          <RouterLink to="/">首页</RouterLink>
          <RouterLink to="/archive">归档</RouterLink>
          <RouterLink to="/about">关于</RouterLink>
          <RouterLink to="/admin/login">后台</RouterLink>
        </nav>

        <div class="header-actions">
          <ThemeToggle />
          <button
            ref="menuButton"
            class="icon-button menu-button"
            type="button"
            :aria-expanded="menuOpen"
            aria-controls="mobile-navigation"
            :aria-label="menuOpen ? '关闭导航菜单' : '打开导航菜单'"
            @click="menuOpen = !menuOpen"
          >
            <span aria-hidden="true">{{ menuOpen ? '×' : '☰' }}</span>
          </button>
        </div>
      </div>
    </header>

    <button v-if="menuOpen" class="nav-backdrop" type="button" aria-label="关闭导航菜单" @click="closeMenu()"></button>
    <nav id="mobile-navigation" class="mobile-nav" :class="{ open: menuOpen }" aria-label="移动端导航">
      <p>导航</p>
      <RouterLink to="/">首页 <span>Home</span></RouterLink>
      <RouterLink to="/archive">归档 <span>Archive</span></RouterLink>
      <RouterLink to="/about">关于 <span>About</span></RouterLink>
      <RouterLink to="/admin/login">后台入口 <span>Admin</span></RouterLink>
    </nav>

    <RouterView />

    <footer class="site-footer">
      <p>Copyright © {{ currentYear }} 字里行间</p>
      <p>Built with Vue 3 and Go</p>
    </footer>
  </div>
</template>
