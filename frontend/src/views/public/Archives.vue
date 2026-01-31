<template>
  <div class="archives-container">
    <main class="main-wrapper">
      <!-- 页面头部 -->
      <div class="page-header">
        <h1 class="page-title">
          <el-icon><Collection /></el-icon> 知识归档
        </h1>
        <p class="page-desc">
          共 {{ postTotal }} 篇博客，收录于 {{ tags.length }} 个标签
        </p>
      </div>

      <!-- 标签云区域 (放在归档上方) -->
      <section class="tags-section">
        <div class="tags-header">
          <h3><el-icon><CollectionTag /></el-icon> 标签云</h3>
        </div>
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

      <el-divider>
        <span class="divider-text"><el-icon><Calendar /></el-icon> 时光轴</span>
      </el-divider>

      <!-- 时光轴区域 -->
      <div v-if="loading" class="loading-box">
        <el-skeleton :rows="5" animated />
      </div>

      <div v-else class="timeline-wrapper">
        <el-timeline>
          <div v-for="(yearGroup, year) in groupedPosts" :key="year">
            <h2 class="year-title">{{ year }}</h2>
            <el-timeline-item
              v-for="post in yearGroup"
              :key="post.ID"
              :timestamp="formatDate(post.created_at)"
              placement="top"
              type="primary"
              hollow
            >
              <el-card class="archive-card" shadow="hover" @click="goToDetail(post.ID)">
                <h3 class="card-title">{{ post.title }}</h3>
                <div class="card-meta">
                  <span v-if="post.category">
                    <el-icon><Folder /></el-icon> {{ post.category.name }}
                  </span>
                  <span>
                    <el-icon><View /></el-icon> {{ post.view_count }}
                  </span>
                </div>
              </el-card>
            </el-timeline-item>
          </div>
        </el-timeline>
        
        <div class="no-more" v-if="posts.length > 0">
           <span>— 记录点滴，连接未来 —</span>
        </div>
      </div>

    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { getPostList, type Post } from '../../api/post'
// 假设你有 API
// import { getTagList } from '../../api/tag'
import { 
  Collection, Calendar, Folder, View, 
  CollectionTag 
} from '@element-plus/icons-vue'

const router = useRouter()

// 状态
const loading = ref(true)

// 归档数据
const posts = ref<Post[]>([])
const postTotal = ref(0)

// 标签数据
const tags = ref<any[]>([])

// 初始化
const init = async () => {
  loading.value = true
  try {
    // 1. 获取文章归档
    const postRes: any = await getPostList({ page: 1, page_size: 100 })
    posts.value = postRes.data || []
    postTotal.value = postRes.total || 0

    // 2. 获取标签
    mockTagsData() 

  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

// Mock 数据逻辑
const mockTagsData = () => {
  tags.value = [
    { id: 1, name: 'Vue3' }, { id: 2, name: 'Golang' }, { id: 3, name: 'Docker' },
    { id: 4, name: 'TypeScript' }, { id: 5, name: 'MySQL' }, { id: 6, name: 'Redis' },
    { id: 7, name: 'Nginx' }, { id: 8, name: 'Linux' }, { id: 9, name: 'React' }
  ]
}

// 归档分组逻辑
const groupedPosts = computed(() => {
  const groups: Record<string, Post[]> = {}
  posts.value.forEach(post => {
    const date = new Date(post.created_at)
    const year = date.getFullYear().toString()
    if (!groups[year]) groups[year] = []
    groups[year].push(post)
  })
  return Object.keys(groups)
    .sort((a, b) => Number(b) - Number(a))
    .reduce((obj, key) => {
      obj[key] = groups[key]
      return obj
    }, {} as Record<string, Post[]>)
})

const formatDate = (dateStr: string) => {
  const date = new Date(dateStr)
  return `${date.getMonth() + 1}月${date.getDate()}日`
}

const getRandomSize = () => Math.floor(Math.random() * 10) + 14 
const getRandomColor = () => {
  const colors = ['#409eff', '#67c23a', '#e6a23c', '#f56c6c', '#909399', '#a855f7']
  return colors[Math.floor(Math.random() * colors.length)]
}

const goToDetail = (id: number) => {
  router.push(`/post/${id}`)
}

onMounted(() => {
  init()
})
</script>

<style scoped>
.archives-container {
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
  margin-bottom: 30px;
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

/* 标签云区域 */
.tags-section {
  background: var(--bg-content);
  border-radius: 12px;
  padding: 20px 30px;
  border: 1px solid var(--border-color);
  box-shadow: var(--shadow-light);
  margin-bottom: 40px;
}
:global(html.dark) .tags-section {
  background-color: #1d1e1f;
  border-color: #363637;
}

.tags-header {
  text-align: center;
  margin-bottom: 20px;
  color: var(--text-main);
}
.tags-header h3 {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  font-size: 1.2rem;
  margin: 0;
}

.tag-cloud {
  text-align: center;
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 15px;
}

.tag-item {
  cursor: pointer;
  transition: transform 0.3s, text-shadow 0.3s;
  display: inline-block;
  padding: 4px 8px;
}
.tag-item:hover { transform: scale(1.1); text-shadow: 0 0 10px currentColor; }

/* 分割线 */
.divider-text {
  font-size: 1.2rem;
  color: var(--text-secondary);
  display: flex;
  align-items: center;
  gap: 8px;
}

/* 时间轴部分 */
.timeline-wrapper { padding: 10px 20px; }
.year-title {
  font-size: 1.8rem;
  font-weight: bold;
  color: transparent;
  background: linear-gradient(to right, var(--primary-color), #a855f7);
  -webkit-background-clip: text;
  background-clip: text;
  margin: 30px 0 20px;
  font-family: 'Impact', sans-serif;
}
:global(html.dark) .year-title { text-shadow: 0 0 15px rgba(168, 85, 247, 0.4); }

.archive-card {
  cursor: pointer;
  border-radius: 8px;
  transition: all 0.3s;
  border: 1px solid var(--border-color);
  background: var(--bg-content); 
  box-shadow: var(--shadow-light);
}
:global(html.dark) .archive-card {
  background-color: #1d1e1f;
  border-color: #363637;
}
.archive-card:hover {
  transform: translateX(5px);
  border-color: var(--primary-color);
}

.card-title { margin: 0 0 10px; font-size: 1.1rem; color: var(--text-main); }
.card-meta { display: flex; gap: 15px; font-size: 12px; color: var(--text-secondary); }
.card-meta span { display: flex; align-items: center; gap: 4px; }

.no-more {
  text-align: center; margin-top: 30px; color: var(--text-secondary); opacity: 0.5; font-size: 12px;
}
</style>