<template>
  <div class="home-container">
    <!-- 1. Hero 区域：视觉焦点 -->
    <section class="hero-section">
      <div class="hero-content">
        <h1 class="hero-title">Thinking inside the box</h1>
        <p class="hero-subtitle">分享技术、代码与生活的点滴感悟</p>
      </div>
      <div class="hero-bg"></div>
    </section>

    <div class="content-wrapper">
      <!-- 加载骨架屏 -->
      <div v-if="loading" class="loading-state">
        <el-row :gutter="20">
          <el-col :span="8" v-for="i in 3" :key="i">
            <el-skeleton style="width: 100%" animated>
              <template #template>
                <el-skeleton-item variant="image" style="width: 100%; height: 200px" />
                <div style="padding: 14px">
                  <el-skeleton-item variant="h3" style="width: 50%" />
                  <el-skeleton-item variant="text" style="margin-top: 10px; width: 80%" />
                </div>
              </template>
            </el-skeleton>
          </el-col>
        </el-row>
      </div>

      <!-- 空状态 -->
      <div v-else-if="posts.length === 0" class="empty-state">
        <el-empty description="博主很懒，暂时还没写文章..." image-size="200" />
      </div>

      <!-- 文章列表：卡片网格布局 -->
      <div v-else class="post-grid">
        <article 
          v-for="(post, index) in posts" 
          :key="post.ID" 
          class="post-card" 
          @click="goToDetail(post.ID)"
        >
          <!-- 封面图区域 (如果没有图片，使用 CSS 渐变) -->
          <div class="card-cover" :style="{ background: getGradient(index) }">
            <span class="category-badge">{{ post.category?.name || 'Uncategorized' }}</span>
          </div>

          <div class="card-body">
            <div class="meta-top">
              <span class="date">{{ formatDate(post.created_at) }}</span>
            </div>
            <h3 class="title">{{ post.title }}</h3>
            <p class="summary">{{ post.summary || '点击阅读更多内容...' }}</p>
            
            <div class="card-footer">
              <div class="author">
                <el-avatar :size="24" :src="'https://i.pravatar.cc/150?u='+post.author?.username" />
                <span>{{ post.author?.username || 'Admin' }}</span>
              </div>
              <span class="read-more">阅读全文 →</span>
            </div>
          </div>
        </article>
      </div>

      <!-- 分页器 -->
      <div class="pagination-wrapper">
        <el-pagination
          background
          layout="prev, pager, next"
          :total="total"
          :page-size="pageSize"
          @current-change="handlePageChange"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getPostList, type Post } from '../../api/post'

const router = useRouter()
const posts = ref<Post[]>([])
const loading = ref(true)
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(6) // 改成6，方便 grid 布局 (3列或2列)

// 生成随机渐变色背景，模拟封面图
const gradients = [
  'linear-gradient(135deg, #a18cd1 0%, #fbc2eb 100%)',
  'linear-gradient(120deg, #84fab0 0%, #8fd3f4 100%)',
  'linear-gradient(120deg, #e0c3fc 0%, #8ec5fc 100%)',
  'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
  'linear-gradient(to top, #fff1eb 0%, #ace0f9 100%)'
]
const getGradient = (index: number) => {
  return gradients[index % gradients.length]
}

const formatDate = (dateStr: string) => {
  return new Date(dateStr).toLocaleDateString('zh-CN', { month: 'short', day: 'numeric', year: 'numeric' })
}

const fetchData = async () => {
  loading.value = true
  try {
    const res: any = await getPostList({
      page: currentPage.value,
      page_size: pageSize.value
    })
    posts.value = res.data || []
    total.value = res.total || 0
  } finally {
    loading.value = false
  }
}

const handlePageChange = (page: number) => {
  currentPage.value = page
  fetchData()
}

const goToDetail = (id: number) => {
  router.push(`/post/${id}`)
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.home-container {
  min-height: 100vh;
}

/* --- Hero Section --- */
.hero-section {
  position: relative;
  height: 300px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  text-align: center;
  margin-bottom: 40px;
  clip-path: polygon(0 0, 100% 0, 100% 85%, 0 100%); /* 切角效果 */
}

.hero-title {
  font-size: 3rem;
  font-weight: 800;
  margin: 0;
  letter-spacing: -1px;
}

.hero-subtitle {
  font-size: 1.2rem;
  opacity: 0.9;
  margin-top: 10px;
  font-weight: 300;
}

/* --- Content Wrapper --- */
.content-wrapper {
  max-width: 1000px;
  margin: 0 auto;
  padding: 0 20px 60px;
}

/* --- Post Grid (核心卡片样式) --- */
.post-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 30px;
}

.post-card {
  background: white;
  border-radius: var(--card-radius);
  overflow: hidden;
  box-shadow: var(--shadow-light);
  transition: all 0.3s ease;
  cursor: pointer;
  display: flex;
  flex-direction: column;
}

.post-card:hover {
  transform: translateY(-8px);
  box-shadow: 0 12px 30px rgba(0, 0, 0, 0.1);
}

.card-cover {
  height: 160px;
  position: relative;
}

.category-badge {
  position: absolute;
  top: 15px;
  right: 15px;
  background: rgba(255, 255, 255, 0.9);
  color: #333;
  padding: 4px 12px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 600;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
}

.card-body {
  padding: 24px;
  flex: 1;
  display: flex;
  flex-direction: column;
}

.meta-top {
  font-size: 12px;
  color: #999;
  margin-bottom: 10px;
  text-transform: uppercase;
  letter-spacing: 1px;
}

.title {
  font-size: 1.25rem;
  margin: 0 0 10px;
  line-height: 1.4;
  color: var(--text-main);
  font-weight: 700;
}

.summary {
  font-size: 14px;
  color: var(--text-secondary);
  line-height: 1.6;
  margin-bottom: 20px;
  flex: 1;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-top: 1px solid #f0f0f0;
  padding-top: 15px;
  margin-top: auto;
}

.author {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: #666;
}

.read-more {
  font-size: 13px;
  color: var(--primary-color);
  font-weight: 600;
}

.pagination-wrapper {
  margin-top: 50px;
  display: flex;
  justify-content: center;
}
</style>