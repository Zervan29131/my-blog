<template>
  <div class="tags-page-container">
    <main class="main-wrapper">
      <!-- 页面头部 -->
      <div class="page-header">
        <h1 class="page-title">
          <el-icon><Collection /></el-icon> 知识图谱
        </h1>
        <p class="page-desc">探索 {{ categories.length }} 个分类与 {{ tags.length }} 个标签</p>
      </div>

      <div v-if="loading" class="loading-box">
        <el-skeleton :rows="5" animated />
      </div>

      <div v-else class="content-wrapper">
        <!-- 1. 分类区域 -->
        <section class="section-block">
          <h2 class="section-title"><el-icon><FolderOpened /></el-icon> 分类</h2>
          <div class="category-list">
            <div 
              v-for="cat in categories" 
              :key="cat.id" 
              class="category-card"
              @click="$router.push(`/category/${cat.id}`)"
            >
              <div class="cat-icon">#</div>
              <div class="cat-info">
                <span class="cat-name">{{ cat.name }}</span>
                <span class="cat-count">{{ cat.count || 0 }} 篇文章</span>
              </div>
            </div>
          </div>
        </section>

        <el-divider />

        <!-- 2. 标签云区域 -->
        <section class="section-block">
          <h2 class="section-title"><el-icon><PriceTag /></el-icon> 标签云</h2>
          <div class="tag-cloud">
            <span 
              v-for="tag in tags" 
              :key="tag.id" 
              class="tag-item"
              :style="{ fontSize: getRandomSize() + 'px', color: getRandomColor() }"
              @click="$router.push(`/tag/${tag.id}`)"
            >
              {{ tag.name }}
            </span>
          </div>
        </section>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Collection, FolderOpened, PriceTag } from '@element-plus/icons-vue'
// 如果你有真实的 API，请取消注释并使用
// import { getCategoryList } from '../../api/category'
// import { getTagList } from '../../api/tag'

const loading = ref(true)
const categories = ref<any[]>([])
const tags = ref<any[]>([])

// 模拟数据 (开发阶段使用，对接后端后请替换)
const mockData = () => {
  setTimeout(() => {
    categories.value = [
      { id: 1, name: '前端开发', count: 12 },
      { id: 2, name: '后端架构', count: 8 },
      { id: 3, name: 'DevOps', count: 5 },
      { id: 4, name: '随笔杂谈', count: 3 },
    ]
    tags.value = [
      { id: 1, name: 'Vue3' }, { id: 2, name: 'Golang' }, { id: 3, name: 'Docker' },
      { id: 4, name: 'TypeScript' }, { id: 5, name: 'MySQL' }, { id: 6, name: 'Redis' },
      { id: 7, name: 'Nginx' }, { id: 8, name: 'Linux' }, { id: 9, name: 'React' }
    ]
    loading.value = false
  }, 500)
}

const getRandomSize = () => Math.floor(Math.random() * 10) + 14 // 14px - 24px
const getRandomColor = () => {
  const colors = ['#409eff', '#67c23a', '#e6a23c', '#f56c6c', '#909399', '#a855f7']
  return colors[Math.floor(Math.random() * colors.length)]
}

onMounted(() => {
  // 真实环境请调用 API
  // Promise.all([getCategoryList(), getTagList()]).then(...)
  mockData()
})
</script>

<style scoped>
.tags-page-container {
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

.section-block { margin-bottom: 40px; }
.section-title {
  font-size: 1.5rem;
  color: var(--text-main);
  margin-bottom: 20px;
  display: flex;
  align-items: center;
  gap: 10px;
}

/* 分类卡片 Grid */
.category-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 20px;
}
.category-card {
  background: var(--bg-content);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 15px;
  cursor: pointer;
  transition: all 0.3s;
}
.category-card:hover {
  transform: translateY(-3px);
  border-color: var(--primary-color);
  box-shadow: var(--shadow-hover);
}
.cat-icon {
  width: 40px;
  height: 40px;
  background: rgba(64, 158, 255, 0.1);
  color: var(--primary-color);
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: bold;
  font-size: 1.2rem;
}
.cat-info { display: flex; flex-direction: column; }
.cat-name { color: var(--text-main); font-weight: bold; }
.cat-count { color: var(--text-secondary); font-size: 12px; margin-top: 4px; }

/* 标签云 */
.tag-cloud {
  background: var(--bg-content);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  padding: 30px;
  text-align: center;
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 20px;
}

.tag-item {
  cursor: pointer;
  transition: transform 0.3s, text-shadow 0.3s;
  display: inline-block;
}
.tag-item:hover {
  transform: scale(1.1);
  text-shadow: 0 0 10px currentColor;
}

/* 深色模式适配 */
:global(html.dark) .category-card {
  background-color: #1d1e1f;
  border-color: #363637;
}
:global(html.dark) .tag-cloud {
  background-color: #1d1e1f;
  border-color: #363637;
}
</style>