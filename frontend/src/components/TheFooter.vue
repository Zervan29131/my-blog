<template>
  <footer class="site-footer">
    <div class="footer-content">
      <div class="footer-logo">My Blog</div>
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
  background-color: #fff;
  border-top: 1px solid #eaeaea;
  padding: 60px 20px 40px;
  margin-top: auto; /* 确保在 flex 布局中沉底 */
  text-align: center;
}

.footer-content {
  max-width: 800px;
  margin: 0 auto;
}

.footer-logo {
  font-size: 1.5rem;
  font-weight: 700;
  color: #333;
  margin-bottom: 10px;
  letter-spacing: -0.5px;
}

.footer-desc {
  color: #666;
  font-size: 14px;
  margin-bottom: 30px;
}

.social-links {
  display: flex;
  justify-content: center;
  gap: 20px;
  margin-bottom: 30px;
}

.social-item {
  color: #909399;
  font-size: 14px;
  font-weight: 500;
  transition: color 0.3s;
}

.social-item:hover {
  color: #409eff;
}

.copyright {
  font-size: 12px;
  color: #ccc;
}
</style>