<script setup lang="ts">
import { onMounted } from 'vue'

onMounted(() => { document.title = '页面不存在 | 字里行间' })
</script>

<template>
  <main class="not-found-page">
    <p class="not-found-code">404</p>
    <h1>页面不存在</h1>
    <p>你访问的页面可能已被删除，或者地址有误。</p>
    <RouterLink class="button button-primary" to="/">返回首页</RouterLink>
  </main>
</template>
