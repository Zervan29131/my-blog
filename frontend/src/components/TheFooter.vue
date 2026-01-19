<template>
  <footer class="site-footer">
    <div class="footer-content">
      <div class="footer-logo">Zervan's Blog</div>
      <p class="footer-desc">记录学习，分享生活，保持好奇心。</p>
      
      <div class="social-links">
        <a href="#" target="_blank" class="social-item">GitHub</a>
        <a href="#" target="_blank" class="social-item">Twitter</a>
        <a href="#" target="_blank" class="social-item">Email</a>
      </div>

      <div class="copyright">
        &copy; {{ new Date().getFullYear() }} My Blog. Powered by Vue 3 & Go.
      </div>
    </div>
  </footer>
</template>

<style scoped>
.site-footer {
  /* 🔴 修改：使用背景色变量 */
  background-color: var(--bg-content);
  /* 🔴 修改：使用边框变量 */
  border-top: 1px solid var(--border-color);
  padding: 60px 20px 40px;
  margin-top: auto;
  text-align: center;
  transition: background-color 0.3s, border-color 0.3s;
}

.footer-content {
  max-width: 800px;
  margin: 0 auto;
}

.footer-logo {
  font-size: 1.5rem;
  font-weight: 700;
  /* 🔴 修改：使用主标题颜色变量 */
  color: var(--text-main);
  margin-bottom: 10px;
  letter-spacing: -0.5px;
  transition: color 0.3s;
}

.footer-desc {
  /* 🔴 修改：使用次级文字颜色变量 */
  color: var(--text-secondary);
  font-size: 14px;
  margin-bottom: 30px;
  transition: color 0.3s;
}

.social-links {
  display: flex;
  justify-content: center;
  gap: 20px;
  margin-bottom: 30px;
}

.social-item {
  /* 🔴 修改：使用常规文字颜色变量 */
  color: var(--text-regular);
  font-size: 14px;
  font-weight: 500;
  transition: color 0.3s;
}

.social-item:hover {
  color: var(--primary-color);
}

.copyright {
  font-size: 12px;
  color: var(--text-secondary);
  transition: color 0.3s;
}
</style>