<template>
  <div class="categories-page-container">
    <main class="main-wrapper">
      <!-- 页面头部 -->
      <div class="page-header">
        <h1 class="page-title">
          <el-icon><FolderOpened /></el-icon> 文章分类
        </h1>
        <p class="page-desc">
          在这里，你可以找到 {{ categories.length }} 个不同主题的分类
        </p>
      </div>

      <div v-if="loading" class="loading-box">
        <el-skeleton :rows="5" animated />
      </div>

      <div v-else class="content-wrapper">
        <div class="category-list">
          <div 
            v-for="cat in categories" 
            :key="cat.id" 
            class="category-card"
            @click="$router.push(`/category/${cat.id}`)"
          >
            <div class="cat-icon-wrapper">
              <el-icon :size="32" class="cat-icon"><Folder /></el-icon>
            </div>
            <div class="cat-info">
              <h3 class="cat-name">{{ cat.name }}</h3>
              <p class="cat-desc" v-if="cat.description">{{ cat.description }}</p>
              <div class="cat-meta">
                <span class="cat-count">{{ cat.count || 0 }} 篇文章</span>
                <el-icon class="arrow-icon"><Right /></el-icon>
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { FolderOpened, Folder, Right } from '@element-plus/icons-vue'
// import { getCategoryList } from '../../api/category'

const loading = ref(true)
const categories = ref<any[]>([])

// 模拟数据
const mockData = () => {
  setTimeout(() => {
    categories.value = [
      { id: 1, name: '前端开发', count: 12, description: 'Vue, React, TypeScript 等前端技术栈' },
      { id: 2, name: '后端架构', count: 8, description: 'Golang, Java, Microservices, System Design' },
      { id: 3, name: 'DevOps', count: 5, description: 'Docker, K8s, CI/CD 自动化部署' },
      { id: 4, name: '随笔杂谈', count: 3, description: '生活感悟，技术思考，读书笔记' },
      { id: 5, name: '数据库', count: 6, description: 'MySQL, Redis, Mongo 等数据库调优' },
    ]
    loading.value = false
  }, 500)
}

onMounted(() => {
  mockData()
})
</script>

<style scoped>
.categories-page-container {
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
.page-desc { color: var(--text-secondary); margin-top: 10px; }

/* 分类卡片 Grid */
.category-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 20px;
}

.category-card {
  background: var(--bg-content);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  padding: 25px;
  display: flex;
  align-items: flex-start;
  gap: 20px;
  cursor: pointer;
  transition: all 0.3s;
  position: relative;
  overflow: hidden;
}

.category-card:hover {
  transform: translateY(-5px);
  border-color: var(--primary-color);
  box-shadow: var(--shadow-hover);
}

:global(html.dark) .category-card {
  background-color: #1d1e1f;
  border-color: #363637;
}

.cat-icon-wrapper {
  background: rgba(64, 158, 255, 0.1);
  width: 60px;
  height: 60px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--primary-color);
  flex-shrink: 0;
}

.cat-info {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.cat-name {
  margin: 0 0 8px;
  font-size: 1.2rem;
  color: var(--text-main);
  font-weight: 600;
}

.cat-desc {
  font-size: 13px;
  color: var(--text-secondary);
  line-height: 1.5;
  margin: 0 0 15px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.cat-meta {
  margin-top: auto;
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 12px;
  color: var(--text-regular);
}

.arrow-icon {
  transition: transform 0.3s;
}

.category-card:hover .arrow-icon {
  transform: translateX(5px);
  color: var(--primary-color);
}
</style>