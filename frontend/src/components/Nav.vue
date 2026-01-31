<template>
  <!-- 移除 :class="{ 'scrolled': isScrolled }"，因为我们不再需要根据滚动改变样式 -->
  <nav class="navbar">
    <div class="nav-container">
      <!-- 1. Logo 区域 -->
      <router-link to="/" class="logo-link" @click="closeMobileMenu">
        <!-- 确保路径正确，如果没有图片会显示文字 -->
        <img src="/img/logo.png" alt="Logo" class="logo-img" />
        <span class="logo-text">Zervan's Blog</span>
      </router-link>

      <!-- 2. 桌面端菜单 (PC显示) -->
      <div class="nav-right hidden-xs">
        <div class="menu-items">
          <router-link to="/" class="nav-item" exact-active-class="active">首页</router-link>
          <router-link to="/categories" class="nav-item" active-class="active">分类</router-link>
          <router-link to="/archives" class="nav-item" active-class="active">归档</router-link>
          <router-link to="/links" class="nav-item" active-class="active">友链</router-link>
          <router-link to="/about" class="nav-item" active-class="active">关于</router-link>
        </div>

        <!-- 搜索框 -->
        <div class="search-box">
          <el-input
            v-model="searchQuery"
            placeholder="搜索文章..."
            class="search-input"
            @keyup.enter="handleSearch"
            :prefix-icon="Search"
            clearable
          />
        </div>

        <!-- 主题切换组件 -->
        <ThemeToggle class="theme-btn" />

        <!-- 管理员入口 -->
        <template v-if="userStore.token">
          <el-button 
            type="primary" 
            size="small" 
            round 
            @click="router.push('/admin')"
            class="auth-btn"
          >
            后台
          </el-button>
        </template>
        <template v-else>
          <el-button 
            text 
            bg 
            size="small" 
            round 
            @click="router.push('/login')"
            class="auth-btn"
          >
            登录
          </el-button>
        </template>
      </div>

      <!-- 3. 移动端菜单按钮 (手机显示) -->
      <div class="mobile-toggle hidden-sm-and-up" @click="toggleMobileMenu">
        <el-icon :size="24" class="toggle-icon">
          <Close v-if="isMobileMenuOpen" />
          <Menu v-else />
        </el-icon>
      </div>
    </div>

    <!-- 4. 移动端抽屉菜单 -->
    <transition name="el-zoom-in-top">
      <div v-if="isMobileMenuOpen" class="mobile-menu hidden-sm-and-up">
        <div class="mobile-search-wrapper">
           <el-input
            v-model="searchQuery"
            placeholder="搜索感兴趣的内容..."
            @keyup.enter="handleSearch"
            :prefix-icon="Search"
            size="large"
          />
        </div>
        
        <router-link to="/" class="mobile-nav-item" exact-active-class="active" @click="closeMobileMenu">
          <span>首页</span>
          <el-icon><ArrowRight /></el-icon>
        </router-link>

        <router-link to="/categories" class="mobile-nav-item" active-class="active" @click="closeMobileMenu">
          <span>分类</span>
          <el-icon><ArrowRight /></el-icon>
        </router-link>

        <router-link to="/archives" class="mobile-nav-item" active-class="active" @click="closeMobileMenu">
          <span>归档</span>
          <el-icon><ArrowRight /></el-icon>
        </router-link>
        
        <router-link to="/links" class="mobile-nav-item" active-class="active" @click="closeMobileMenu">
          <span>友链</span>
          <el-icon><ArrowRight /></el-icon>
        </router-link>

        <router-link to="/about" class="mobile-nav-item" active-class="active" @click="closeMobileMenu">
          <span>关于我</span>
          <el-icon><ArrowRight /></el-icon>
        </router-link>

        <div class="mobile-actions">
           <div class="action-row">
             <span>切换主题</span>
             <ThemeToggle />
           </div>
           
           <div class="action-row" v-if="!userStore.token" @click="router.push('/login');closeMobileMenu()">
             <span>管理员登录</span>
             <el-icon><User /></el-icon>
           </div>
           <div class="action-row" v-else @click="router.push('/admin');closeMobileMenu()">
             <span>进入后台</span>
             <el-icon><Monitor /></el-icon>
           </div>
        </div>
      </div>
    </transition>
  </nav>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'
import { Search, Menu, Close, User, Monitor, ArrowRight } from '@element-plus/icons-vue'
import ThemeToggle from './ThemeToggle.vue'

const router = useRouter()
const userStore = useUserStore()

// 移除 isScrolled 状态，因为样式不再依赖滚动
const isMobileMenuOpen = ref(false)
const searchQuery = ref('')

// 移除 handleScroll 函数

// 搜索逻辑
const handleSearch = () => {
  if (!searchQuery.value.trim()) return
  
  router.push({
    path: '/search',
    query: { q: searchQuery.value }
  })
  
  closeMobileMenu()
}

// 移动端菜单控制
const toggleMobileMenu = () => {
  isMobileMenuOpen.value = !isMobileMenuOpen.value
}

const closeMobileMenu = () => {
  isMobileMenuOpen.value = false
}

// 移除 scroll 监听
// onMounted(() => {
//   window.addEventListener('scroll', handleScroll)
// })

// onUnmounted(() => {
//   window.removeEventListener('scroll', handleScroll)
// })
</script>

<style scoped>
/* 导航栏容器 */
.navbar {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 64px;
  z-index: 1000;
  transition: all 0.3s ease-in-out;
  /* 固定背景色 */
  background: var(--bg-content, #ffffff);
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
  pointer-events: auto; 
}

/* 移除 .navbar.scrolled 相关样式 */

/* 暗色模式适配 - 始终应用深色背景 */
:global(html.dark) .navbar {
  background: #1d1e1f;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.5);
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

.nav-container {
  max-width: 1200px;
  margin: 0 auto;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
}

/* Logo */
.logo-link {
  display: flex;
  align-items: center;
  gap: 10px;
  text-decoration: none;
  color: var(--text-main, #333);
}

.logo-img {
  width: 32px;
  height: 32px;
  border-radius: 50%;
}

.logo-text {
  font-size: 1.25rem;
  font-weight: 700;
  letter-spacing: 0.5px;
  color: var(--text-main);
  transition: color 0.3s;
}

/* 右侧菜单区 */
.nav-right {
  display: flex;
  align-items: center;
  gap: 20px;
}

.menu-items {
  display: flex;
  gap: 20px;
}

.nav-item {
  text-decoration: none;
  font-size: 15px;
  color: var(--text-regular, #606266);
  font-weight: 500;
  transition: color 0.3s;
  position: relative;
}

:global(html.dark) .nav-item {
  color: #a3a6ad;
}

.nav-item:hover, .nav-item.active {
  color: var(--primary-color, #409eff);
}

.nav-item.active::after {
  content: '';
  position: absolute;
  bottom: -5px;
  left: 0;
  width: 100%;
  height: 2px;
  background-color: var(--primary-color, #409eff);
  border-radius: 2px;
}

/* 搜索框微调 */
.search-box {
  width: 200px;
}

:deep(.el-input__wrapper) {
  border-radius: 20px;
  background-color: var(--bg-content, #f5f7fa);
  box-shadow: none !important;
  border: 1px solid transparent;
}

:global(html.dark) :deep(.el-input__wrapper) {
  background-color: #2c2c2c;
  color: #fff;
}

:deep(.el-input__wrapper.is-focus) {
  background-color: var(--bg-content, #fff);
  border-color: var(--primary-color, #409eff);
}

:global(html.dark) :deep(.el-input__wrapper.is-focus) {
  background-color: #2c2c2c;
}

.auth-btn {
  font-weight: 600;
}

/* 移动端切换按钮 */
.mobile-toggle {
  cursor: pointer;
  color: var(--text-main);
  display: flex;
  align-items: center;
}

/* 移动端菜单 */
.mobile-menu {
  position: absolute;
  top: 64px;
  left: 0;
  width: 100%;
  background: var(--bg-content, #fff);
  border-top: 1px solid var(--border-color, #ebeef5);
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 15px;
}

:global(html.dark) .mobile-menu {
  background: #1d1e1f;
  border-top: 1px solid #333;
}

.mobile-search-wrapper {
  margin-bottom: 10px;
}

.mobile-nav-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 0;
  border-bottom: 1px solid var(--border-color, #f0f2f5);
  color: var(--text-main);
  font-size: 16px;
  text-decoration: none;
}

:global(html.dark) .mobile-nav-item {
  border-bottom: 1px solid #333;
  color: #E5EAF3;
}

.mobile-nav-item.active {
  color: var(--primary-color);
  font-weight: 600;
  border-bottom-color: var(--primary-color);
}

.mobile-actions {
  margin-top: 10px;
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.action-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  color: var(--text-secondary);
  font-size: 14px;
  cursor: pointer;
}

/* 响应式断点 */
@media (max-width: 768px) {
  .hidden-xs { display: none !important; }
  .hidden-sm-and-up { display: block !important; }
}
@media (min-width: 769px) {
  .hidden-xs { display: flex !important; }
  .hidden-sm-and-up { display: none !important; }
}
</style>