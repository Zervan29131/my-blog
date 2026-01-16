<template>
  <div class="public-layout">
    <!-- 🟢 关键修改：使用 <Nav /> 组件替换了原来的 <header> 代码块 -->
    <Nav />

    <main class="main-content">
      <router-view />
    </main>

    <TheFooter />
    <BackToTop />
  </div>
</template>

<script setup lang="ts">
// 1. 引入新组件
import Nav from '../components/Nav.vue'
import TheFooter from '../components/TheFooter.vue'
import BackToTop from '../components/BackToTop.vue'
</script>

<style scoped>
.public-layout {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}
/* 注意：原来的 .site-header 等样式已经移到 Nav.vue 中了，这里不再需要 */
</style>