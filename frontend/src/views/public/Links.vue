<template>
  <div class="links-container">
    <main class="main-wrapper">
      <!-- 头部信息 -->
      <div class="page-header">
        <h1 class="page-title">
          <el-icon><Link /></el-icon> 友情链接
        </h1>
        <p class="page-desc">
          海内存知己，天涯若比邻。<br>
          欢迎交换友链，共同学习进步。
        </p>
      </div>

      <!-- 申请规则 (可选) -->
      <el-card class="rule-card" shadow="hover">
        <template #header>
          <div class="card-header">
            <span>💡 友链申请说明</span>
          </div>
        </template>
        <div class="rule-content">
          <p>请确保您的网站内容积极向上，并已添加本站链接。</p>
          <p><strong>名称：</strong>Zervan's Blog</p>
          <p><strong>简介：</strong>这是存在于某处的幻想世界</p>
          <p><strong>地址：</strong>https://blog.zervan.com</p>
        </div>
      </el-card>

      <!-- 链接列表 -->
      <div v-if="loading" class="loading-box">
        <el-skeleton :rows="3" animated />
      </div>

      <div v-else class="links-grid">
        <a 
          v-for="link in links" 
          :key="link.id" 
          :href="link.url" 
          target="_blank" 
          class="link-card"
        >
          <div class="link-avatar">
            <img :src="link.avatar" :alt="link.name" loading="lazy" />
          </div>
          <div class="link-info">
            <h3 class="link-name">{{ link.name }}</h3>
            <p class="link-desc">{{ link.desc }}</p>
          </div>
        </a>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Link } from '@element-plus/icons-vue'

interface FriendlyLink {
  id: number
  name: string
  desc: string
  avatar: string
  url: string
}

const loading = ref(true)
const links = ref<FriendlyLink[]>([])

// 模拟获取数据
const fetchLinks = () => {
  setTimeout(() => {
    links.value = [
      {
        id: 1,
        name: 'Vue.js',
        desc: '渐进式 JavaScript 框架',
        avatar: 'https://vuejs.org/images/logo.png',
        url: 'https://vuejs.org/'
      },
      {
        id: 2,
        name: 'Vite',
        desc: '下一代前端开发与构建工具',
        avatar: 'https://vitejs.dev/logo.svg',
        url: 'https://vitejs.dev/'
      },
      {
        id: 3,
        name: 'Go 语言',
        desc: 'Build simple, secure, scalable systems with Go',
        avatar: 'https://go.dev/images/go-logo-blue.svg',
        url: 'https://go.dev/'
      },
      {
        id: 4,
        name: 'Element Plus',
        desc: '基于 Vue 3，面向设计师和开发者的组件库',
        avatar: 'https://element-plus.org/images/element-plus-logo.svg',
        url: 'https://element-plus.org/'
      },
      // 你可以在这里添加更多模拟数据
    ]
    loading.value = false
  }, 600)
}

onMounted(() => {
  fetchLinks()
})
</script>

<style scoped>
.links-container {
  padding-top: 80px;
  min-height: 100vh;
  background-color: var(--bg-color);
  transition: background-color 0.3s;
}

.main-wrapper {
  max-width: 900px;
  margin: 0 auto;
  padding: 20px;
}

.page-header {
  text-align: center;
  margin-bottom: 40px;
}
.page-title {
  font-size: 2rem;
  color: var(--text-main);
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
}
.page-desc {
  color: var(--text-secondary);
  margin-top: 10px;
  line-height: 1.6;
}

/* 规则卡片 */
.rule-card {
  margin-bottom: 40px;
  background-color: var(--bg-content);
  border: 1px solid var(--border-color);
  color: var(--text-main);
}
.card-header { font-weight: bold; }
.rule-content p { margin: 5px 0; font-size: 14px; color: var(--text-regular); }
.rule-content strong { color: var(--primary-color); }

/* 深色模式适配卡片 */
:global(html.dark) .rule-card {
  background-color: #1d1e1f;
  border-color: #363637;
  color: var(--text-main);
}

/* 链接 Grid */
.links-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  gap: 20px;
}

.link-card {
  display: flex;
  align-items: center;
  padding: 15px;
  background: var(--bg-content);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  text-decoration: none;
  transition: all 0.3s;
  cursor: pointer;
}

.link-card:hover {
  transform: translateY(-5px);
  box-shadow: var(--shadow-hover);
  border-color: var(--primary-color);
}

.link-avatar {
  width: 50px;
  height: 50px;
  margin-right: 15px;
  flex-shrink: 0;
  border-radius: 50%;
  overflow: hidden;
  background-color: #f5f7fa;
  display: flex;
  align-items: center;
  justify-content: center;
}

.link-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.link-info {
  flex: 1;
  overflow: hidden;
}

.link-name {
  margin: 0 0 5px;
  font-size: 16px;
  font-weight: bold;
  color: var(--text-main);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.link-desc {
  margin: 0;
  font-size: 12px;
  color: var(--text-secondary);
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  line-height: 1.4;
}

/* 深色模式适配链接卡片 */
:global(html.dark) .link-card {
  background-color: #1d1e1f;
  border-color: #363637;
}
:global(html.dark) .link-card:hover {
  background-color: #252627;
  border-color: var(--primary-color);
}
:global(html.dark) .link-avatar {
  background-color: #2c2c2c;
}
</style>