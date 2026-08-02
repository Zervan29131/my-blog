<script setup lang="ts">
import { computed, type Component } from 'vue'

import type { PublicHomepageConfig, PublicHomepageModuleType } from '../../types/blog'
import HomeAbout from './HomeAbout.vue'
import HomeFeaturedArticles from './HomeFeaturedArticles.vue'
import HomeHero from './HomeHero.vue'
import HomeLatestArticles from './HomeLatestArticles.vue'
import HomepageModuleBoundary from './HomepageModuleBoundary.vue'
import HomeSocialLinks from './HomeSocialLinks.vue'
import HomeTechStack from './HomeTechStack.vue'

const props = defineProps<{
  homepage: PublicHomepageConfig
  siteName: string
}>()

const homepageComponentMap: Record<PublicHomepageModuleType, Component> = {
  hero: HomeHero,
  about: HomeAbout,
  featured_articles: HomeFeaturedArticles,
  latest_articles: HomeLatestArticles,
  tech_stack: HomeTechStack,
  social_links: HomeSocialLinks,
}

const moduleLabels: Record<PublicHomepageModuleType, string> = {
  hero: '欢迎区域',
  about: '个人简介',
  featured_articles: '推荐文章',
  latest_articles: '最新文章',
  tech_stack: '技术栈',
  social_links: '社交链接',
}

const renderedModules = computed(() => props.homepage.modules.map((module) => ({
  module,
  component: homepageComponentMap[module.type],
  label: moduleLabels[module.type],
})))
</script>

<template>
  <div v-if="renderedModules.length" class="homepage-module-stack">
    <HomepageModuleBoundary v-for="entry in renderedModules" :key="entry.module.type" :label="entry.label">
      <component
        :is="entry.component"
        :config="entry.module.config"
        :site-name="entry.module.type === 'hero' ? siteName : undefined"
      />
    </HomepageModuleBoundary>
  </div>
  <section v-else class="home-empty-config">
    <strong>{{ siteName }}</strong>
    <p>首页内容正在整理中，你仍可以通过导航浏览文章。</p>
  </section>
</template>
