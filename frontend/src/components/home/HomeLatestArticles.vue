<script setup lang="ts">
import { RouterLink } from 'vue-router'

import type { PublicLatestArticlesConfig } from '../../types/blog'
import HomeArticleList from './HomeArticleList.vue'

defineProps<{ config: PublicLatestArticlesConfig }>()
</script>

<template>
  <section class="home-dynamic-section home-articles-section" aria-labelledby="latest-articles">
    <header class="home-module-heading home-article-heading">
      <div><p>RECENT POSTS</p><h2 id="latest-articles">{{ config.title }}</h2><span v-if="config.description">{{ config.description }}</span></div>
      <RouterLink v-if="config.show_view_all" class="section-link" to="/archive">全部文章 <span aria-hidden="true">→</span></RouterLink>
    </header>
    <HomeArticleList v-if="config.articles.length" :articles="config.articles" :show-summary="config.show_summary" :show-date="config.show_date" :show-comment-count="config.show_comment_count" />
    <div v-else class="state-panel empty-state"><strong>暂时还没有文章</strong><p>新的内容正在准备中。</p></div>
  </section>
</template>
