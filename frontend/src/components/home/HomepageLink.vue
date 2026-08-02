<script setup lang="ts">
import { RouterLink } from 'vue-router'

import type { PublicHomepageButton } from '../../types/blog'

defineProps<{ button: PublicHomepageButton }>()
</script>

<template>
  <RouterLink v-if="button.link_type === 'internal'" :to="button.url"><slot /></RouterLink>
  <a v-else :href="button.url" :target="button.open_in_new_tab ? '_blank' : undefined" :rel="button.open_in_new_tab ? 'noopener noreferrer' : undefined"><slot /></a>
</template>
