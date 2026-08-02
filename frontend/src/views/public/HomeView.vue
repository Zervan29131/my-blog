<script setup lang="ts">
import { onMounted, ref } from 'vue'

import { fetchArticles, fetchHomepageConfig } from '../../api/blog'
import { apiErrorMessage } from '../../api/http'
import HomepageRenderer from '../../components/home/HomepageRenderer.vue'
import ErrorState from '../../components/ErrorState.vue'
import LoadingState from '../../components/LoadingState.vue'
import { usePublicConfigStore } from '../../stores/publicConfig'
import type { PublicHomepageConfig } from '../../types/blog'
import { normalizePublicHomepage } from '../../utils/publicConfig'

const configStore = usePublicConfigStore()
const homepage = ref<PublicHomepageConfig>({ version: 0, modules: [] })
const loading = ref(true)
const errorMessage = ref('')
const fallbackLoaded = ref(false)

async function loadFallbackArticles() {
  const result = await fetchArticles(1, 10)
  homepage.value = {
    version: 0,
    modules: [{
      type: 'latest_articles',
      sort_order: 40,
      config: {
        title: '最新文章', description: '', limit: 10,
        show_summary: true, show_date: true, show_comment_count: true, show_view_all: true,
        articles: result.items,
      },
    }],
  }
  fallbackLoaded.value = true
}

async function loadHomepage() {
  loading.value = true
  errorMessage.value = ''
  fallbackLoaded.value = false
  try {
    homepage.value = normalizePublicHomepage(await fetchHomepageConfig())
  } catch (error) {
    errorMessage.value = apiErrorMessage(error, '首页配置加载失败，已尝试保留文章列表。')
    try {
      await loadFallbackArticles()
    } catch {
      homepage.value = { version: 0, modules: [] }
    }
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  void configStore.load()
  void loadHomepage()
})
</script>

<template>
  <main class="page-container home-page dynamic-home-page">
    <LoadingState v-if="loading" label="正在加载首页内容…" />

    <template v-else>
      <div v-if="errorMessage && fallbackLoaded" class="homepage-error-banner" role="status">
        <span>{{ errorMessage }}</span>
        <button type="button" @click="loadHomepage">重新加载配置</button>
      </div>
      <ErrorState v-else-if="errorMessage && !homepage.modules.length" :message="errorMessage" @retry="loadHomepage" />
      <HomepageRenderer v-if="homepage.modules.length || !errorMessage" :homepage="homepage" :site-name="configStore.site.name" />
    </template>
  </main>
</template>
