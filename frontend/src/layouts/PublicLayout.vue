<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref, watch, watchEffect } from 'vue'
import { RouterLink, RouterView, useRoute } from 'vue-router'

import ThemeToggle from '../components/ThemeToggle.vue'
import { usePublicConfigStore } from '../stores/publicConfig'

const route = useRoute()
const configStore = usePublicConfigStore()
const menuOpen = ref(false)
const menuButton = ref<HTMLButtonElement | null>(null)
const logoFailed = ref(false)
const currentYear = new Date().getFullYear()
const brandName = computed(() => configStore.site.short_name || configStore.site.name)
const copyrightYears = computed(() => {
  const start = configStore.footer.start_year
  return start && start < currentYear ? `${start}–${currentYear}` : String(currentYear)
})

function closeMenu(restoreFocus = false) {
  menuOpen.value = false
  if (restoreFocus) menuButton.value?.focus()
}

function handleKeydown(event: KeyboardEvent) {
  if (event.key === 'Escape' && menuOpen.value) closeMenu(true)
}

function updateFavicon(url: string) {
  if (!url) return
  let favicon = document.querySelector<HTMLLinkElement>('link[rel="icon"]')
  if (!favicon) {
    favicon = document.createElement('link')
    favicon.rel = 'icon'
    document.head.append(favicon)
  }
  favicon.href = url
}

function pageTitle(): string {
  const suffix = configStore.titleName
  if (route.name === 'home') return suffix === configStore.site.name ? configStore.site.name : `${configStore.site.name} | ${suffix}`
  if (route.name === 'archive') return `归档 | ${suffix}`
  if (route.name === 'about') return `关于 | ${suffix}`
  if (route.name === 'preview-home') return `首页预览 | ${suffix}`
  if (route.name === 'not-found') return `页面不存在 | ${suffix}`
  return suffix
}

watch(() => route.fullPath, () => closeMenu())
watch(menuOpen, (open) => {
  document.body.classList.toggle('nav-open', open)
})
watch(() => configStore.site.logo_url, () => { logoFailed.value = false })
watch(() => configStore.site.favicon_url, updateFavicon, { immediate: true })
watchEffect(() => {
  if (route.name === 'article-detail') return
  document.title = pageTitle()
  const description = document.querySelector<HTMLMetaElement>('meta[name="description"]')
  if (description) description.content = configStore.site.description
})

onMounted(() => {
  window.addEventListener('keydown', handleKeydown)
  void configStore.load()
})
onUnmounted(() => {
  window.removeEventListener('keydown', handleKeydown)
  document.body.classList.remove('nav-open')
})
</script>

<template>
  <div class="site-shell">
    <header class="site-header">
      <div class="header-inner">
        <RouterLink class="brand" to="/" :aria-label="`${configStore.site.name}博客首页`">
          <img v-if="configStore.site.logo_url && !logoFailed" class="brand-logo" :src="configStore.site.logo_url" alt="" @error="logoFailed = true" />
          <span v-else class="brand-mark" aria-hidden="true"><i></i></span>
          <strong>{{ brandName }}</strong>
        </RouterLink>

        <nav class="desktop-nav" aria-label="主要导航">
          <template v-for="item in configStore.navigation" :key="`${item.name}-${item.url}`">
            <RouterLink v-if="item.link_type === 'internal'" :to="item.url">{{ item.name }}</RouterLink>
            <a v-else :href="item.url" :target="item.open_in_new_tab ? '_blank' : undefined" :rel="item.open_in_new_tab ? 'noopener noreferrer' : undefined">{{ item.name }}</a>
          </template>
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
      <template v-for="item in configStore.navigation" :key="`mobile-${item.name}-${item.url}`">
        <RouterLink v-if="item.link_type === 'internal'" :to="item.url">{{ item.name }} <span>Internal</span></RouterLink>
        <a v-else :href="item.url" :target="item.open_in_new_tab ? '_blank' : undefined" :rel="item.open_in_new_tab ? 'noopener noreferrer' : undefined">{{ item.name }} <span>External</span></a>
      </template>
      <RouterLink to="/admin/login">后台入口 <span>Admin</span></RouterLink>
    </nav>

    <div v-if="configStore.error" class="site-config-notice" role="status">
      <span>{{ configStore.error }}</span>
      <button type="button" @click="configStore.load(true)">重新加载</button>
    </div>

    <RouterView />

    <footer class="site-footer">
      <div>
        <p>Copyright © {{ copyrightYears }} {{ configStore.footer.copyright_name }}</p>
        <p v-if="configStore.footer.additional_text">{{ configStore.footer.additional_text }}</p>
        <a v-if="configStore.footer.filing_number && configStore.footer.filing_url" :href="configStore.footer.filing_url" target="_blank" rel="noopener noreferrer">{{ configStore.footer.filing_number }}</a>
        <p v-else-if="configStore.footer.filing_number">{{ configStore.footer.filing_number }}</p>
      </div>
      <p v-if="configStore.footer.show_technology && configStore.footer.technology_text">{{ configStore.footer.technology_text }}</p>
    </footer>
  </div>
</template>
