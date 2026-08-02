<script setup lang="ts">
import { computed, ref, watch } from 'vue'

import type { PublicAboutConfig } from '../../types/blog'
import { renderMarkdown } from '../../utils/markdown'

const props = defineProps<{ config: PublicAboutConfig }>()
const imageFailed = ref(false)
const renderedContent = computed(() => renderMarkdown(props.config.content))
const showImage = computed(() => props.config.image_position !== 'none' && props.config.image_url && !imageFailed.value)

watch(() => props.config.image_url, () => { imageFailed.value = false })
</script>

<template>
  <section class="home-dynamic-section home-about-section" :class="`about-image-${showImage ? config.image_position : 'none'}`">
    <img v-if="showImage" class="home-about-image" :src="config.image_url" alt="" @error="imageFailed = true" />
    <div class="home-about-copy">
      <header class="home-module-heading"><p>ABOUT</p><h2>{{ config.title }}</h2><span v-if="config.description">{{ config.description }}</span></header>
      <div class="markdown-body home-about-markdown" v-html="renderedContent"></div>
    </div>
  </section>
</template>
