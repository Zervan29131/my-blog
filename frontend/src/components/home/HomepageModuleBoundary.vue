<script setup lang="ts">
import { onErrorCaptured, ref } from 'vue'

defineProps<{ label: string }>()
const failed = ref(false)

onErrorCaptured((error) => {
  failed.value = true
  console.warn('Homepage module render failed', error)
  return false
})
</script>

<template>
  <slot v-if="!failed" />
  <section v-else class="home-module-fallback" role="status">
    <strong>{{ label }}暂时无法显示</strong>
    <span>其他首页内容仍可正常浏览。</span>
  </section>
</template>
