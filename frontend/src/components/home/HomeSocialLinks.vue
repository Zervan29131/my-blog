<script setup lang="ts">
import type { PublicSocialLinksConfig } from '../../types/blog'

defineProps<{ config: PublicSocialLinksConfig }>()
</script>

<template>
  <section v-if="config.links.length" class="home-dynamic-section home-social-section">
    <header class="home-module-heading"><p>CONNECT</p><h2>{{ config.title }}</h2><span v-if="config.description">{{ config.description }}</span></header>
    <div class="home-social-list">
      <a v-for="link in config.links" :key="`${link.platform}-${link.url}`" :href="link.url" :target="link.url.startsWith('mailto:') ? undefined : '_blank'" :rel="link.url.startsWith('mailto:') ? undefined : 'noopener noreferrer'">
        <span aria-hidden="true">{{ link.display_name.charAt(0).toUpperCase() }}</span>
        <strong>{{ link.display_name }}</strong>
        <small>{{ link.platform }}</small>
      </a>
    </div>
  </section>
</template>
