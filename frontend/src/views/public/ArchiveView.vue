<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'

import { fetchArticles } from '../../api/blog'
import { apiErrorMessage } from '../../api/http'
import ErrorState from '../../components/ErrorState.vue'
import LoadingState from '../../components/LoadingState.vue'
import PaginationControls from '../../components/PaginationControls.vue'
import type { ArticleSummary } from '../../types/blog'

const articles = ref<ArticleSummary[]>([])
const page = ref(1)
const total = ref(0)
const totalPages = ref(0)
const loading = ref(true)
const errorMessage = ref('')
const groupedArticles = computed(() => {
  const groups = new Map<string, ArticleSummary[]>()
  for (const article of articles.value) {
    const year = new Date(article.published_at).getFullYear().toString()
    groups.set(year, [...(groups.get(year) || []), article])
  }
  return Array.from(groups.entries())
})

function monthDay(date: string): string {
  return new Intl.DateTimeFormat('zh-CN', { month: '2-digit', day: '2-digit' }).format(new Date(date))
}

async function loadArticles(targetPage = page.value) {
  loading.value = true
  errorMessage.value = ''
  try {
    const result = await fetchArticles(targetPage, 10)
    articles.value = result.items
    page.value = result.page
    total.value = result.total
    totalPages.value = result.total_pages
  } catch (error) {
    errorMessage.value = apiErrorMessage(error, '归档加载失败，请稍后重试。')
  } finally {
    loading.value = false
  }
}

function changePage(targetPage: number) {
  void loadArticles(targetPage)
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

onMounted(() => {
  void loadArticles()
})
</script>

<template>
  <main class="page-container simple-page">
    <header class="page-heading">
      <p class="eyebrow">ARCHIVE</p>
      <h1>归档</h1>
      <p>按时间浏览已发布的文章，共 {{ total }} 篇。</p>
    </header>

    <LoadingState v-if="loading" label="正在整理归档…" />
    <ErrorState v-else-if="errorMessage" :message="errorMessage" @retry="loadArticles()" />
    <div v-else-if="groupedArticles.length" class="archive-groups">
      <section v-for="[year, yearArticles] in groupedArticles" :key="year" class="archive-group">
        <h2>{{ year }}</h2>
        <ul>
          <li v-for="article in yearArticles" :key="article.id">
            <time :datetime="article.published_at">{{ monthDay(article.published_at) }}</time>
            <RouterLink :to="`/articles/${article.slug}`">{{ article.title }}</RouterLink>
          </li>
        </ul>
      </section>
    </div>
    <div v-else class="state-panel empty-state">
      <strong>暂时还没有文章</strong>
      <p>新的内容正在准备中。</p>
    </div>

    <PaginationControls :page="page" :total-pages="totalPages" @change="changePage" />
  </main>
</template>
