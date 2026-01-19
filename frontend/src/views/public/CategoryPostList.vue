<template>
  <div class="category-container">
    <!-- 头部：展示当前分类名称 -->
    <div class="category-header">
      <div class="header-content">
        <h1>
          <el-icon><FolderOpened /></el-icon> 
          <span>{{ categoryName }}</span>
        </h1>
        <p>Category Collection</p>
      </div>
    </div>

    <div class="main-wrapper">
      <el-row :gutter="20">
        <el-col :span="18" :xs="24">
          
          <!-- 加载状态 -->
          <div v-if="loading" class="loading-box">
             <el-skeleton :rows="5" animated />
          </div>

          <!-- 空状态 -->
          <div v-else-if="posts.length === 0" class="empty-box">
            <el-empty description="该分类下暂无文章" />
            <el-button type="primary" @click="$router.push('/')">返回首页</el-button>
          </div>

          <!-- 文章列表 -->
          <div v-else class="post-list">
            <article 
              v-for="post in posts" 
              :key="post.ID" 
              class="post-item"
              @click="goToDetail(post.ID)"
            >
              <div class="post-content-wrapper">
                <h2 class="post-title">{{ post.title }}</h2>
                <p class="post-summary">{{ post.summary || '暂无摘要' }}</p>
                
                <div class="post-meta">
                   <span class="meta-item"><el-icon><Calendar /></el-icon> {{ formatDate(post.created_at) }}</span>
                   <span class="meta-divider">|</span>
                   <span class="meta-item"><el-icon><View /></el-icon> {{ post.view_count }} 阅读</span>
                   <span class="meta-divider">|</span>
                   <span class="meta-item author"><el-icon><User /></el-icon> {{ post.author?.username || 'Admin' }}</span>
                </div>
              </div>
              
              <!-- 只有当文章有封面图时才显示 (这里用占位符演示) -->
              <div class="post-cover-mini">
                 <img :src="`https://picsum.photos/seed/${post.ID}/100/100`" alt="cover">
              </div>
            </article>
          </div>
          
           <!-- 分页 -->
           <div class="pagination-box" v-if="total > 0">
             <el-pagination
              background
              layout="prev, pager, next"
              :total="total"
              :page-size="pageSize"
              v-model:current-page="currentPage"
              @current-change="handlePageChange"
            />
          </div>
        </el-col>

        <!-- 侧边栏占位 -->
        <el-col :span="6" class="hidden-xs-only">
           <!-- 可以在这里复用 Sidebar 组件 -->
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getPostList, type Post } from '../../api/post'
import { FolderOpened, Calendar, View, User } from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()
const posts = ref<Post[]>([])
const loading = ref(true)
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const categoryName = ref('分类浏览')

const formatDate = (date: string) => new Date(date).toLocaleDateString()

const fetchData = async () => {
  loading.value = true
  const cId = Number(route.params.id) // 获取 URL 中的分类 ID
  
  if (!cId) {
    loading.value = false
    return
  }

  try {
    const res: any = await getPostList({
      page: currentPage.value,
      page_size: pageSize.value,
      category_id: cId // 🟢 传给后端
    })
    posts.value = res.data || []
    total.value = res.total || 0
    
    // 尝试从第一篇文章中获取分类名 (更严谨的做法是单独调 getCategoryById)
    if (posts.value.length > 0 && posts.value[0].category) {
      categoryName.value = posts.value[0].category.name
    } else {
      // 如果列表为空，我们可能无法直接知道分类名，除非有单独的 category 详情接口
      // 这里可以暂时显示为 "分类 ID: xxx"
      categoryName.value = `分类 ID: ${cId}`
    }
  } catch (err) {
    console.error(err)
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

// 监听路由参数变化（例如从 Code 分类点到了 Fun 分类）
watch(() => route.params.id, () => {
  currentPage.value = 1
  fetchData()
})

onMounted(fetchData)
</script>

<style scoped>
.category-container {
  min-height: 80vh;
  background-color: #f5f7fa;
  padding-bottom: 40px;
}

.category-header {
  background: #fff;
  padding: 40px 0;
  text-align: center;
  margin-bottom: 30px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.05);
}

.header-content h1 {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  margin: 0;
  color: #2c3e50;
  font-size: 2rem;
}

.header-content p {
  color: #909399;
  margin-top: 10px;
  letter-spacing: 1px;
}

.main-wrapper {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
}

/* Post Item Styles */
.post-item {
  background: #fff;
  padding: 25px;
  margin-bottom: 20px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  justify-content: space-between;
  align-items: center;
  border: 1px solid transparent;
}

.post-item:hover {
  transform: translateY(-3px);
  box-shadow: 0 8px 20px rgba(0,0,0,0.06);
  border-color: #ecf5ff;
}

.post-content-wrapper {
  flex: 1;
  padding-right: 20px;
}

.post-title {
  margin: 0 0 12px;
  font-size: 1.3rem;
  color: #303133;
  transition: color 0.3s;
}

.post-item:hover .post-title {
  color: #409eff;
}

.post-summary {
  color: #606266;
  font-size: 14px;
  line-height: 1.6;
  margin-bottom: 15px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.post-meta {
  display: flex;
  align-items: center;
  font-size: 13px;
  color: #909399;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 4px;
}

.meta-divider {
  margin: 0 10px;
  color: #e4e7ed;
}

.post-cover-mini img {
  width: 100px;
  height: 100px;
  border-radius: 6px;
  object-fit: cover;
}

.pagination-box {
  margin-top: 30px;
  display: flex;
  justify-content: center;
}
</style>