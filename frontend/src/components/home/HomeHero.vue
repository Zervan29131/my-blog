<script setup lang="ts">
import { computed, ref, watch } from 'vue'

import type { PublicHeroConfig } from '../../types/blog'
import HomepageLink from './HomepageLink.vue'

const props = defineProps<{ config: PublicHeroConfig; siteName: string }>()
const imageFailed = ref(false)
const backgroundFailed = ref(false)
const titleParts = computed(() => {
  const highlight = props.config.highlight_text
  const index = highlight ? props.config.title.indexOf(highlight) : -1
  return index >= 0
    ? { before: props.config.title.slice(0, index), highlight, after: props.config.title.slice(index + highlight.length) }
    : null
})
const backgroundStyle = computed(() => props.config.background_image_url && !backgroundFailed.value
  ? { backgroundImage: `linear-gradient(rgb(14 15 18 / 78%), rgb(14 15 18 / 88%)), url("${props.config.background_image_url}")` }
  : undefined)

watch(() => props.config.image_url, () => { imageFailed.value = false })
watch(() => props.config.background_image_url, () => { backgroundFailed.value = false })
</script>

<template>
  <section class="hero-section dynamic-hero" :class="`hero-layout-${config.layout}`" :style="backgroundStyle">
    <img v-if="config.background_image_url && !backgroundFailed" class="hero-background-probe" :src="config.background_image_url" alt="" @error="backgroundFailed = true" />
    <div class="hero-copy">
      <div v-if="!config.image_url || imageFailed" class="hero-monogram" aria-hidden="true">{{ siteName.charAt(0).toUpperCase() }}</div>
      <p v-if="config.eyebrow" class="eyebrow">{{ config.eyebrow }}</p>
      <h1>
        <template v-if="titleParts">{{ titleParts.before }}<span>{{ titleParts.highlight }}</span>{{ titleParts.after }}</template>
        <template v-else>{{ config.title }}</template>
      </h1>
      <p class="hero-description">{{ config.description }}</p>
      <div v-if="config.primary_button.enabled || config.secondary_button.enabled" class="hero-actions">
        <HomepageLink v-if="config.primary_button.enabled" class="button button-primary" :button="config.primary_button">{{ config.primary_button.text }}</HomepageLink>
        <HomepageLink v-if="config.secondary_button.enabled" class="button button-secondary" :button="config.secondary_button">{{ config.secondary_button.text }}</HomepageLink>
      </div>
    </div>
    <img v-if="config.image_url && !imageFailed" class="hero-visual" :src="config.image_url" alt="" @error="imageFailed = true" />
  </section>
</template>
