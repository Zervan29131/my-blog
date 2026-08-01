<script setup lang="ts">
import { computed } from 'vue'

import { useTheme } from '../composables/useTheme'

const { theme, effectiveTheme, cycleTheme } = useTheme()
const modeLabel = computed(() => ({
  system: '跟随系统',
  light: '浅色模式',
  dark: '深色模式',
}[theme.value]))
const nextLabel = computed(() => ({
  system: '浅色模式',
  light: '深色模式',
  dark: '跟随系统',
}[theme.value]))
</script>

<template>
  <button
    class="icon-button theme-toggle"
    type="button"
    :aria-label="`当前为${modeLabel}，切换到${nextLabel}`"
    :title="`主题：${modeLabel}`"
    @click="cycleTheme"
  >
    <span class="theme-icon" aria-hidden="true">{{ effectiveTheme === 'dark' ? '☾' : '☼' }}</span>
    <span class="sr-only">{{ modeLabel }}</span>
  </button>
</template>

