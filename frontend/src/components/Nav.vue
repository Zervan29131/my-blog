<template>
  <header class="site-header" :class="{ 'transparent': isTransparent }">
    <div class="nav-wrapper">
      <!-- Logo 区域 -->
      <div class="logo" @click="$router.push('/')">
        Zervan's Blog
      </div>

      <!-- 菜单链接 -->
      <nav class="nav-links">
        <router-link to="/" active-class="active">首页</router-link>
        <router-link to="/archives" active-class="active">归档</router-link>
        <router-link to="/categories" active-class="active">分类</router-link>
        <router-link to="/tags" active-class="active">标签</router-link>
        <router-link to="/about" active-class="active">关于</router-link>
        
        <!-- 新增：搜索框 -->
        <div class="search-box">
           <el-input
             v-model="searchKeyword"
             placeholder="搜索..."
             :prefix-icon="Search"
             class="search-input"
             @keyup.enter="handleSearch"
           />
        </div>

        <!-- 状态感知按钮 -->
        <div class="auth-action">
          <template v-if="userStore.isLoggedIn">
            <router-link to="/admin/dashboard" class="dashboard-btn">
              <el-icon><Odometer /></el-icon> 控制台
            </router-link>
          </template>
          
          <template v-else>
            <router-link to="/login" class="login-btn">
              登录
            </router-link>
          </template>
        </div>
      </nav>
    </div>
  </header>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '../stores/user'
import { Odometer, Search } from '@element-plus/icons-vue'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const searchKeyword = ref('')
const isTransparent = ref(false)

// 简单的搜索处理
const handleSearch = () => {
  if (searchKeyword.value.trim()) {
    // 这里暂时打印，后续可以跳转到专门的搜索页
    console.log('Search:', searchKeyword.value)
    // router.push({ path: '/search', query: { q: searchKeyword.value } })
  }
}

// 监听滚动和路由变化，控制导航栏透明度
const handleScroll = () => {
  const scrollTop = window.scrollY
  // 只有在首页且滚动距离小于 50px 时才透明
  isTransparent.value = scrollTop < 50 && route.path === '/'
}

onMounted(() => {
  window.addEventListener('scroll', handleScroll)
  // 初始化检查一次
  handleScroll()
})

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll)
})

// 监听路由变化（比如从首页跳到关于页，要立刻取消透明）
watch(() => route.path, () => {
  handleScroll()
})
</script>

<style scoped>
.site-header {
  position: sticky;
  top: 0;
  z-index: 100;
  width: 100%;
  height: 64px;
  background: rgba(255, 255, 255, 0.85);
  backdrop-filter: blur(12px);
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.02);
  transition: all 0.3s ease;
}

/* 透明模式 (用于首页 Hero 顶部) */
.site-header.transparent {
  background: transparent;
  box-shadow: none;
  border-bottom: none;
  /* 在深色背景图上，文字变白 */
  color: white; 
}

.nav-wrapper {
  max-width: 1200px;
  height: 100%;
  margin: 0 auto;
  padding: 0 20px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.logo {
  font-size: 1.25rem;
  font-weight: 800;
  /* 默认深色 */
  color: #2c3e50;
  cursor: pointer;
  letter-spacing: -0.5px;
  transition: color 0.3s;
}

.logo:hover {
  color: #409eff;
}

.nav-links {
  display: flex;
  align-items: center;
  gap: 20px;
}

.nav-links a {
  font-size: 15px;
  font-weight: 500;
  color: #606266;
  text-decoration: none;
  position: relative;
  transition: color 0.3s;
}

.nav-links a:not(.login-btn):not(.dashboard-btn):hover,
.nav-links a.active {
  color: #409eff;
}

/* 下划线动画 */
.nav-links a:not(.login-btn):not(.dashboard-btn)::after {
  content: '';
  position: absolute;
  bottom: -4px;
  left: 0;
  width: 0;
  height: 2px;
  background-color: #409eff;
  transition: width 0.3s ease;
  border-radius: 2px;
}

.nav-links a.active::after,
.nav-links a:not(.login-btn):not(.dashboard-btn):hover::after {
  width: 100%;
}

/* --- 透明模式下的文字颜色适配 --- */
.site-header.transparent .logo {
  color: rgba(255, 255, 255, 0.95);
  text-shadow: 0 2px 4px rgba(0,0,0,0.3);
}

.site-header.transparent .nav-links a {
  color: rgba(255, 255, 255, 0.85);
  text-shadow: 0 1px 2px rgba(0,0,0,0.3);
}

.site-header.transparent .nav-links a:hover,
.site-header.transparent .nav-links a.active {
  color: #fff;
}

.site-header.transparent .nav-links a::after {
  background-color: #fff;
}

/* --- 搜索框样式 --- */
.search-box {
  width: 180px;
  transition: width 0.3s;
}
.search-box:focus-within {
  width: 240px;
}

:deep(.search-input .el-input__wrapper) {
  border-radius: 20px;
  background: rgba(240, 242, 245, 0.8);
  box-shadow: none;
  padding-left: 15px;
}

/* 透明模式下搜索框稍微半透明一点 */
.site-header.transparent :deep(.search-input .el-input__wrapper) {
  background: rgba(255, 255, 255, 0.2);
  color: white;
}
.site-header.transparent :deep(.el-input__inner) {
  color: white;
}
.site-header.transparent :deep(.el-input__inner::placeholder) {
  color: rgba(255,255,255,0.7);
}

/* 按钮样式 */
.login-btn,
.dashboard-btn {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 6px 16px;
  border-radius: 20px;
  font-size: 13px !important;
  font-weight: 600;
  transition: all 0.3s ease;
}

.login-btn {
  background-color: #f0f2f5;
  color: #606266 !important;
}

.login-btn:hover {
  background-color: #e6e8eb;
  color: #303133 !important;
}

.dashboard-btn {
  background-color: #e6f7ff;
  color: #1890ff !important;
}

.dashboard-btn:hover {
  background-color: #409eff;
  color: white !important;
}
</style>