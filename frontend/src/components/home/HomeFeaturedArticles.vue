<script setup lang="ts">
import type { PublicArticlesModuleConfig } from '../../types/blog'
import HomeArticleList from './HomeArticleList.vue'

defineProps<{ config: PublicArticlesModuleConfig }>()
</script>

<template>
  <section v-if="config.articles.length" class="home-dynamic-section home-articles-section">
    <header class="home-module-heading"><p>FEATURED</p><h2>{{ config.title }}</h2><span v-if="config.description">{{ config.description }}</span></header>
    <HomeArticleList :articles="config.articles" featured />
  </section>
</template>
