<template>
  <div class="not-found-container">
    <div class="content">
      <div class="error-code">404</div>
      <h1 class="title">页面迷路了</h1>
      <p class="description">
        抱歉，我们找不到您要访问的页面。<br>
        它可能已被移动、删除，或者您输入的网址有误。
      </p>
      
      <div class="actions">
        <el-button type="primary" size="large" round @click="goHome">
          返回首页
        </el-button>
        <el-button size="large" round @click="goBack">
          返回上一页
        </el-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'

const router = useRouter()

const goHome = () => {
  router.push('/')
}

const goBack = () => {
  router.back()
}
</script>

<style scoped>
.not-found-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: var(--bg-color);
  transition: background-color 0.3s;
  padding: 20px;
  text-align: center;
}

.content {
  max-width: 600px;
  animation: fadeInUp 0.5s ease-out;
}

.error-code {
  font-size: 8rem;
  font-weight: 800;
  line-height: 1;
  color: var(--primary-color);
  opacity: 0.2;
  margin-bottom: -20px;
  user-select: none;
}

.title {
  font-size: 2.5rem;
  color: var(--text-main);
  margin-bottom: 20px;
  position: relative;
}

.description {
  font-size: 1.1rem;
  color: var(--text-secondary);
  line-height: 1.6;
  margin-bottom: 40px;
}

.actions {
  display: flex;
  gap: 20px;
  justify-content: center;
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 移动端适配 */
@media (max-width: 768px) {
  .error-code {
    font-size: 6rem;
  }
  .title {
    font-size: 2rem;
  }
}
</style>