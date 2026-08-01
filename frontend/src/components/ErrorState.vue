<script setup lang="ts">
defineProps<{
  message: string
}>()

defineEmits<{
  retry: []
}>()
</script>

<template>
  <div class="state-panel error-state" role="alert">
    <p>{{ message }}</p>
    <button class="button button-secondary" type="button" @click="$emit('retry')">重新加载</button>
  </div>
</template>
