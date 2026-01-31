<template>
  <div class="search-results-container">
    <main class="main-wrapper">
      <!-- 搜索状态头部 -->
      <div class="search-header">
        <h2 class="search-title">
          <el-icon class="icon-search"><Search /></el-icon> 
          <span>搜索结果: <span class="highlight">"{{ searchQuery }}"</span></span>
        </h2>
        <p class="search-meta" v-if="!loading">
          共找到 <span class="count">{{ total }}</span> 篇相关文章
        </p>
      </div>

      <el-row :gutter="20">
        <el-col :span="24">
          <!-- 加载中 -->
          <div v-if="loading" class="loading-box">
            <el-skeleton :rows="5" animated />
          </div>

          <!-- 空状态 -->
          <div v-else-if="posts.length === 0" class="empty-box">
            <el-empty :image-size="200" :description="`抱歉，没有找到与 '${searchQuery}' 相关的文章`">
              <el-button type="primary" @click="$router.push('/')">返回首页</el-button>
            </el-empty>
          </div>

          <!-- 文章列表 -->
          <div v-else class="post-list">
            <!-- 使用 PostCard 组件 -->
            <PostCard 
              v-for="(post, index) in posts" 
              :key="post.ID" 
              :post="post"
              :reverse="index % 2 !== 0"
              @click="goToDetail"
            />
          </div>

          <!-- 分页 -->
          <div class="pagination-box" v-if="total > pageSize">
            <el-pagination
              background
              layout="prev, pager, next"
              :total="total"
              :page-size="pageSize"
              :current-page="currentPage"
              @current-change="handlePageChange"
            />
          </div>
        </el-col>
      </el-row>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getPostList, type Post } from '../../api/post'
import { Search } from '@element-plus/icons-vue'
import PostCard from '../../components/PostCard.vue'

const route = useRoute()
const router = useRouter()

const posts = ref<Post[]>([])
const loading = ref(true)
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const searchQuery = ref('')

// 获取数据
const fetchData = async () => {
  // 从路由参数获取关键词，如果没有则为空
  const q = route.query.q ? String(route.query.q).trim() : ''
  searchQuery.value = q
  
  if (!q) {
    loading.value = false
    posts.value = []
    total.value = 0
    return
  }

  loading.value = true
  try {
    // 调用 API，传递 q 参数
    const res: any = await getPostList({
      page: currentPage.value,
      page_size: pageSize.value,
      q: q
    })
    posts.value = res.data || []
    total.value = res.total || 0
  } catch (error) {
    console.error('搜索失败:', error)
    posts.value = []
  } finally {
    loading.value = false
  }
}

const handlePageChange = (page: number) => {
  currentPage.value = page
  fetchData()
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

const goToDetail = (id: number) => {
  router.push(`/post/${id}`)
}

// 监听路由参数变化（例如用户在搜索结果页再次搜索其他词）
watch(
  () => route.query.q,
  (newVal) => {
    if (newVal !== undefined) {
      currentPage.value = 1
      fetchData()
    }
  }
)

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.search-results-container {
  padding-top: 80px; /* 留出导航栏高度 */
  min-height: 80vh;
  background-color: var(--bg-color);
  transition: background-color 0.3s;
}

.main-wrapper {
  max-width: 900px;
  margin: 0 auto;
  padding: 20px;
  animation: fadeInUp 0.5s ease-out;
}

/* 搜索头部 */
.search-header {
  margin-bottom: 40px;
  padding-bottom: 20px;
  border-bottom: 1px solid var(--border-color);
  text-align: center;
}

.search-title {
  font-size: 2rem;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 15px;
  color: var(--text-main);
  margin-bottom: 10px;
}

.icon-search {
  color: var(--primary-color);
}

.highlight {
  color: var(--primary-color);
  font-style: italic;
}

.search-meta {
  color: var(--text-secondary);
  font-size: 1rem;
}

.count {
  font-weight: bold;
  color: var(--text-main);
  margin: 0 4px;
}

/* 加载与空状态 */
.loading-box, .empty-box {
  padding: 40px 0;
  min-height: 300px;
  display: flex;
  justify-content: center;
  flex-direction: column;
}

/* 分页 */
.pagination-box {
  display: flex;
  justify-content: center;
  margin-top: 40px;
  margin-bottom: 40px;
}

/* 动画 */
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
  .search-title {
    font-size: 1.5rem;
  }
  .main-wrapper {
    padding: 20px 10px;
  }
}
</style>