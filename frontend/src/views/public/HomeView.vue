<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'

import { fetchArticles } from '../../api/blog'
import { apiErrorMessage } from '../../api/http'
import ErrorState from '../../components/ErrorState.vue'
import LoadingState from '../../components/LoadingState.vue'
import PaginationControls from '../../components/PaginationControls.vue'
import type { ArticleSummary } from '../../types/blog'
import { formatDate } from '../../utils/format'

const articles = ref<ArticleSummary[]>([])
const page = ref(1)
const total = ref(0)
const totalPages = ref(0)
const loading = ref(true)
const errorMessage = ref('')

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
    errorMessage.value = apiErrorMessage(error, '文章加载失败，请稍后重试。')
  } finally {
    loading.value = false
  }
}

function changePage(targetPage: number) {
  void loadArticles(targetPage)
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

onMounted(() => {
  document.title = '字里行间 | 技术与生活手记'
  void loadArticles()
})
</script>

<template>
  <main class="page-container home-page">
    <section class="hero-section">
      <div class="hero-monogram" aria-hidden="true">Z</div>
      <p class="eyebrow">ZERVAN'S JOURNAL</p>
      <h1>在字里行间，<br /><span>保存思想的回声。</span></h1>
      <p class="hero-description">记录开发、阅读与生活中的思考。把遇到的问题讲清楚，也把值得留存的经验认真写下来。</p>
      <div class="hero-actions">
        <a class="button button-primary" href="#latest-articles">开始阅读</a>
        <RouterLink class="button button-secondary" to="/about">关于本站</RouterLink>
      </div>
      <p class="hero-quote">Be the change you want to see in the world.</p>
    </section>

    <div class="home-content-grid">
      <section class="article-section" aria-labelledby="latest-articles">
        <div class="section-heading">
          <div>
            <p class="section-kicker">RECENT POSTS</p>
            <h2 id="latest-articles">最新文章</h2>
          </div>
          <RouterLink class="section-link" to="/archive">全部文章 <span aria-hidden="true">→</span></RouterLink>
        </div>

        <LoadingState v-if="loading" label="正在整理文章…" />
        <ErrorState v-else-if="errorMessage" :message="errorMessage" @retry="loadArticles()" />
        <div v-else-if="articles.length" class="article-list">
          <article v-for="(article, index) in articles" :key="article.id" class="article-card">
            <span v-if="index === 0" class="article-card-badge">NEW</span>
            <p class="article-card-type">ARTICLE</p>
            <h3>
              <RouterLink :to="`/articles/${article.slug}`">{{ article.title }}</RouterLink>
            </h3>
            <p>{{ article.summary || '这篇文章暂时没有摘要。' }}</p>
            <div class="article-card-meta">
              <span class="meta-author">Zervan</span>
              <time :datetime="article.published_at">{{ formatDate(article.published_at) }}</time>
              <span>{{ article.comment_count }} 条评论</span>
            </div>
            <RouterLink class="article-card-arrow" :to="`/articles/${article.slug}`" :aria-label="`阅读《${article.title}》`">→</RouterLink>
          </article>
        </div>
        <div v-else class="state-panel empty-state">
          <strong>暂时还没有文章</strong>
          <p>新的内容正在准备中。</p>
        </div>

        <PaginationControls
          :page="page"
          :total-pages="totalPages"
          @change="changePage"
        />
      </section>

      <aside class="home-sidebar" aria-label="站点信息">
        <section class="sidebar-card profile-card">
          <div class="profile-avatar" aria-hidden="true">Z</div>
          <p class="profile-label">ABOUT THE AUTHOR</p>
          <h2>Zervan</h2>
          <p>保持好奇，持续记录。<br />天下最普通的人之一。</p>
          <div class="profile-links">
            <RouterLink to="/about">关于我</RouterLink>
            <RouterLink to="/archive">文章归档</RouterLink>
          </div>
        </section>

        <section class="sidebar-card site-overview">
          <div class="sidebar-card-heading">
            <span aria-hidden="true">⌁</span>
            <h2>站点导航</h2>
          </div>
          <RouterLink to="/archive"><span>全部文章</span><strong>{{ total }}</strong></RouterLink>
          <RouterLink to="/about"><span>关于本站</span><strong>↗</strong></RouterLink>
          <RouterLink to="/admin/login"><span>内容管理</span><strong>↗</strong></RouterLink>
        </section>

        <section class="sidebar-card stack-card">
          <p class="profile-label">BUILT WITH</p>
          <div><span>Vue 3</span><span>Go</span><span>PostgreSQL</span></div>
        </section>
      </aside>
    </div>
  </main>
</template>
