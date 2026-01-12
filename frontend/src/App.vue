<template>
  <!-- 路由出口：所有页面（如 Login.vue, PostList.vue）都会渲染在这里 -->
  <router-view />
</template>

<script setup lang="ts">
// 根组件通常不需要太多逻辑，主要作为容器
</script>

<style>
/* 全局基础样式设置 */
html, body, #app {
  margin: 0;
  padding: 0;
  height: 100%;
  width: 100%;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
}
</style>