<script setup lang="ts">
import { RouterLink } from 'vue-router'

import type { ArticleSummary } from '../../types/blog'
import { formatDate } from '../../utils/format'

withDefaults(defineProps<{
  articles: ArticleSummary[]
  showSummary?: boolean
  showDate?: boolean
  showCommentCount?: boolean
  featured?: boolean
}>(), {
  showSummary: true,
  showDate: true,
  showCommentCount: true,
  featured: false,
})
</script>

<template>
  <div class="article-list home-module-article-list">
    <article v-for="(article, index) in articles" :key="article.id" class="article-card">
      <span v-if="featured || index === 0" class="article-card-badge">{{ featured ? 'PICK' : 'NEW' }}</span>
      <p class="article-card-type">{{ featured ? 'FEATURED' : 'ARTICLE' }}</p>
      <h3><RouterLink :to="`/articles/${article.slug}`">{{ article.title }}</RouterLink></h3>
      <p v-if="showSummary">{{ article.summary || '这篇文章暂时没有摘要。' }}</p>
      <div v-if="showDate || showCommentCount" class="article-card-meta">
        <time v-if="showDate" :datetime="article.published_at">{{ formatDate(article.published_at) }}</time>
        <span v-if="showCommentCount">{{ article.comment_count }} 条评论</span>
      </div>
      <RouterLink class="article-card-arrow" :to="`/articles/${article.slug}`" :aria-label="`阅读《${article.title}》`">→</RouterLink>
    </article>
  </div>
</template>
