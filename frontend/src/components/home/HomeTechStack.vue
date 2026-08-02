<script setup lang="ts">
import { ref } from 'vue'

import type { PublicTechStackConfig } from '../../types/blog'

defineProps<{ config: PublicTechStackConfig }>()
const failedIcons = ref(new Set<string>())

function failIcon(url: string) {
  failedIcons.value = new Set([...failedIcons.value, url])
}
</script>

<template>
  <section v-if="config.items.length" class="home-dynamic-section home-tech-section">
    <header class="home-module-heading"><p>TOOLKIT</p><h2>{{ config.title }}</h2><span v-if="config.description">{{ config.description }}</span></header>
    <div class="home-tech-grid">
      <component :is="item.url ? 'a' : 'article'" v-for="item in config.items" :key="item.name" class="home-tech-card" :href="item.url || undefined" :target="item.url ? '_blank' : undefined" :rel="item.url ? 'noopener noreferrer' : undefined">
        <img v-if="item.icon_url && !failedIcons.has(item.icon_url)" :src="item.icon_url" alt="" @error="failIcon(item.icon_url)" />
        <span v-else aria-hidden="true">{{ item.name.charAt(0).toUpperCase() }}</span>
        <div><strong>{{ item.name }}</strong><p v-if="item.description">{{ item.description }}</p></div>
      </component>
    </div>
  </section>
</template>
