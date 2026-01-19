<template>
  <header class="site-header" :class="{ 'transparent': isTransparent }">
    <div class="nav-wrapper">
      <!-- Logo -->
      <div class="logo" @click="$router.push('/')">
        Zervan的小站
      </div>
        <!-- 🟢 升级：即时搜索框 -->
        <div class="search-box">
          <el-autocomplete
            v-model="searchKeyword"
            :fetch-suggestions="querySearchAsync"
            placeholder="搜索..."
            :prefix-icon="Search"
            class="search-input"
            :trigger-on-focus="false"
            @select="handleSelect"
            @keyup.enter="handleEnter"
            popper-class="search-dropdown"
          >
          
            <!-- 自定义下拉模板 -->
            <template #default="{ item }">
              <div class="search-item">
                <div class="search-title" v-html="highlight(item.title)"></div>
                <div class="search-meta">
                  <span class="date">{{ formatDate(item.created_at) }}</span>
                  <span class="category" v-if="item.category">
                    <el-icon><Folder /></el-icon> {{ item.category.name }}
                  </span>
                </div>
              </div>
            </template>
          </el-autocomplete>
        </div>
        <!-- 🟢 动态分类渲染 -->
        <!-- 这里的 index 设置为路由路径，点击时 el-menu 会自动跳转 -->
        <!-- <el-menu-item 
          v-for="cat in categories" 
          :key="cat.ID" 
          :index="'/category/' + cat.ID"
        >
          {{ cat.name }}
        </el-menu-item>
         -->
      <!-- 菜单链接 -->
      <nav class="nav-links">
        <router-link to="/" active-class="active">Home</router-link>
        <router-link to="/archives" active-class="active">code</router-link>
        <router-link to="/categories" active-class="active">ToRead</router-link>
        <router-link to="/tags" active-class="active">fun</router-link>
        <router-link to="/about" active-class="active">about</router-link>
        

        <!-- 登录/控制台按钮 -->
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

        <ThemeToggle />

      </nav>
    </div>
  </header>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '../stores/user'
import { Odometer, Search, Folder } from '@element-plus/icons-vue'
import ThemeToggle from './ThemeToggle.vue'
import { getPostList, type Post } from '../api/post'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const searchKeyword = ref('')
const isTransparent = ref(false)

// 🟢 核心：异步搜索逻辑
const querySearchAsync = async (queryString: string, cb: (results: any[]) => void) => {
  if (!queryString.trim()) {
    cb([])
    return
  }

  try {
    // 复用现有的文章列表接口，只取前 5 条作为预览
    const res: any = await getPostList({
      page: 1,
      page_size: 5,
      q: queryString
    })
    
    const results = res.data || []
    
    // 如果没有结果，可以返回一个特殊的提示项（可选）
    if (results.length === 0) {
      cb([{ title: '未找到相关文章', ID: 0, disabled: true }])
    } else {
      cb(results)
    }
  } catch (error) {
    console.error('Search error:', error)
    cb([])
  }
}

// 🟢 选中下拉项：直接跳转详情页
const handleSelect = (item: any) => {
  if (item.disabled) return
  router.push(`/post/${item.ID}`)
  searchKeyword.value = '' // 跳转后清空搜索框
}

// 🟢 回车键：跳转到完整搜索结果页
const handleEnter = () => {
  if (searchKeyword.value.trim()) {
    router.push({ path: '/', query: { q: searchKeyword.value } })
    // el-autocomplete 在回车时如果不手动 blur，下拉框可能不会收起，这里强制收起
    ;(document.activeElement as HTMLElement)?.blur()
  }
}

// 辅助：高亮关键词
const highlight = (title: string) => {
  if (!searchKeyword.value) return title
  const reg = new RegExp(`(${searchKeyword.value})`, 'gi')
  return title.replace(reg, '<span class="highlight">$1</span>')
}

const formatDate = (dateStr: string) => {
  return new Date(dateStr).toLocaleDateString('zh-CN', { month: '2-digit', day: '2-digit' })
}

// 滚动透明逻辑
const handleScroll = () => {
  const scrollTop = window.scrollY
  isTransparent.value = scrollTop < 50 && route.path === '/'
}

onMounted(() => {
  window.addEventListener('scroll', handleScroll)
  handleScroll()
})

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll)
})

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
  background: var(--bg-header);
  backdrop-filter: blur(12px);
  border-bottom: 1px solid var(--border-color);
  transition: all 0.3s ease;
}

.site-header.transparent {
  background: transparent;
  border-bottom: none;
  box-shadow: none;
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
  color: var(--text-main);
  cursor: pointer;
  letter-spacing: -0.5px;
  transition: color 0.3s;
}

.site-header.transparent .logo { color: white; text-shadow: 0 2px 4px rgba(0,0,0,0.3); }
.logo:hover { color: var(--primary-color); }

.nav-links {
  display: flex;
  align-items: center;
  gap: 15px;
}

.nav-links a {
  font-size: 15px;
  font-weight: 500;
  color: var(--text-regular);
  text-decoration: none;
  position: relative;
  transition: color 0.3s;
}

.site-header.transparent .nav-links a { color: rgba(255,255,255,0.85); text-shadow: 0 1px 2px rgba(0,0,0,0.3); }
.nav-links a:not(.login-btn):not(.dashboard-btn):hover,
.nav-links a.active { color: var(--primary-color); }
.site-header.transparent .nav-links a.active { color: white; }

/* 🟢 搜索框样式升级 */
.search-box {
  width: 220px;
  transition: width 0.3s;
}
.search-box:focus-within {
  width: 280px;
}

/* 穿透修改 el-autocomplete 样式 */
:deep(.el-autocomplete) {
  width: 100%;
}
:deep(.el-input__wrapper) {
  border-radius: 20px;
  background: rgba(128, 128, 128, 0.1); 
  box-shadow: none !important;
  padding-left: 15px;
}
.site-header.transparent :deep(.el-input__wrapper) {
  background: rgba(255, 255, 255, 0.2);
}
.site-header.transparent :deep(.el-input__inner) { color: white; }
.site-header.transparent :deep(.el-input__inner::placeholder) { color: rgba(255,255,255,0.7); }

/* 按钮 */
.login-btn, .dashboard-btn {
  display: flex; align-items: center; gap: 4px; padding: 6px 16px; border-radius: 20px; font-size: 13px !important; font-weight: 600; transition: all 0.3s ease;
}
.login-btn { background-color: rgba(128, 128, 128, 0.1); color: var(--text-regular) !important; }
.dashboard-btn { background-color: rgba(64, 158, 255, 0.1); color: var(--primary-color) !important; }
</style>

<!-- 🟢 全局样式：下拉框样式定制 -->
<style>
.search-dropdown {
  border-radius: 12px !important;
  border: 1px solid var(--border-color) !important;
  background-color: var(--bg-content) !important;
  box-shadow: 0 8px 24px rgba(0,0,0,0.15) !important;
}

.search-item {
  padding: 8px 0;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.search-title {
  font-weight: 600;
  color: var(--text-main);
  font-size: 14px;
  line-height: 1.4;
}

/* 高亮样式 */
.highlight {
  color: var(--primary-color);
  font-weight: bold;
}

.search-meta {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  color: var(--text-secondary);
}

.search-meta .category {
  display: flex;
  align-items: center;
  gap: 3px;
}

/* 深色模式适配 Element Plus 下拉框 */
html.dark .search-dropdown {
  border-color: #363637 !important;
  background-color: #1e1e1e !important;
}
html.dark .el-autocomplete-suggestion__list li:hover {
  background-color: rgba(255, 255, 255, 0.05) !important;
}
</style>