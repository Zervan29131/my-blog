<template>
  <div class="public-layout">
    <header class="site-header">
      <div class="container">
        <h1 class="logo">My Blog</h1>
        <nav>
          <router-link to="/">首页</router-link>
          <router-link to="/about">关于</router-link>
        </nav>
      </div>
    </header>

    <main class="container main-content">
      <!-- 页面内容会渲染在这里 -->
      <router-view />
    </main>

    <footer class="site-footer">
      <p>&copy; 2026 My Blog. Powered by Vue3 & Go.</p>
    </footer>
  </div>
</template>

<style scoped>
.container { max-width: 800px; margin: 0 auto; padding: 0 20px; }
.site-header { border-bottom: 1px solid #eee; padding: 20px 0; }
.site-footer { text-align: center; padding: 40px 0; color: #999; }
</style>