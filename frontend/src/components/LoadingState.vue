<script setup lang="ts">
defineProps<{
  label?: string
}>()
</script>

<template>
  <div class="state-panel loading-state" role="status" aria-live="polite">
    <span class="loading-dot" aria-hidden="true"></span>
    <p>{{ label || '正在加载…' }}</p>
  </div>
</template>
