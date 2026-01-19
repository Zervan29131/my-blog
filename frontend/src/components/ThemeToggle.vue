<template>
  <button class="theme-toggle" @click="toggleDark()" title="切换主题">
    <!-- 太阳图标 (浅色模式显示) -->
    <el-icon v-if="!isDark" class="icon sun"><Sunny /></el-icon>
    <!-- 月亮图标 (深色模式显示) -->
    <el-icon v-else class="icon moon"><Moon /></el-icon>
  </button>
</template>

<script setup lang="ts">
import { useDark, useToggle } from '@vueuse/core'
import { Sunny, Moon } from '@element-plus/icons-vue'

// 核心逻辑：useDark 会自动检测系统偏好，并持久化到 localStorage
const isDark = useDark({
  selector: 'html',
  attribute: 'class',
  valueDark: 'dark',
  valueLight: '',
})

const toggleDark = useToggle(isDark)
</script>

<style scoped>
.theme-toggle {
  background: transparent;
  border: 1px solid transparent;
  cursor: pointer;
  padding: 8px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s;
  color: #606266;
}

.theme-toggle:hover {
  background-color: rgba(0, 0, 0, 0.05);
  color: #409eff;
}

html.dark .theme-toggle {
  color: #cfd3dc;
}

html.dark .theme-toggle:hover {
  background-color: rgba(255, 255, 255, 0.1);
  color: #409eff;
}

.icon {
  font-size: 18px;
}
</style>