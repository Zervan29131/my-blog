<template>
  <div class="tag-container">
    <main class="main-wrapper">
      <!-- 头部信息 -->
      <div class="tag-header">
        <div class="header-content">
          <span class="sub-title">TAG ARCHIVE</span>
          <h1 class="tag-title">
             <el-icon class="icon-tag"><Collection /></el-icon>
             {{ tagName || '加载中...' }}
          </h1>
          <p class="tag-desc">与此标签相关的文章集合</p>
        </div>
      </div>

      <el-row :gutter="20">
        <el-col :span="24">
          <div v-if="loading" class="loading-box">
             <el-skeleton :rows="5" animated />
          </div>

          <div v-else-if="posts.length === 0" class="empty-box">
            <el-empty description="该标签下暂无文章" />
            <el-button @click="$router.push('/')">返回首页</el-button>
          </div>

          <div v-else class="post-list">
            <PostCard 
              v-for="(post, index) in posts" 
              :key="post.ID" 
              :post="post"
              :reverse="index % 2 !== 0"
              @click="goToDetail"
            />
          </div>

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
import { Collection } from '@element-plus/icons-vue'
import PostCard from '../../components/PostCard.vue'

const route = useRoute()
const router = useRouter()

const posts = ref<Post[]>([])
const loading = ref(true)
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const tagName = ref('') 

const fetchData = async () => {
  loading.value = true
  const tagId = Number(route.params.id)
  
  if (!tagId) {
    loading.value = false
    return
  }

  try {
    const res: any = await getPostList({
      page: currentPage.value,
      page_size: pageSize.value,
      tag_id: tagId 
    })
    
    posts.value = res.data || []
    total.value = res.total || 0
    tagName.value = `标签 #${tagId}` 

  } catch (error) {
    console.error(error)
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

watch(() => route.params.id, (newVal) => {
  if (newVal) {
    currentPage.value = 1
    tagName.value = '' 
    fetchData()
  }
})

onMounted(() => fetchData())
</script>

<style scoped>
.tag-container {
  padding-top: 80px;
  min-height: 100vh;
  background-color: var(--bg-color); /* 关键 */
  transition: background-color 0.3s;
}

.main-wrapper {
  max-width: 900px;
  margin: 0 auto;
  padding: 20px;
}

.tag-header {
  text-align: center;
  margin-bottom: 50px;
  padding: 40px 0;
  background: var(--bg-content); /* 关键 */
  border-radius: 12px;
  box-shadow: var(--shadow-light);
  border: 1px solid var(--border-color); /* 关键 */
  transition: all 0.3s;
}

/* 🟢 深色模式：Header 增强 */
:global(html.dark) .tag-header {
  background: linear-gradient(145deg, #1d1e1f, #252627);
  border-color: #363637;
}

.sub-title {
  font-size: 14px;
  color: var(--text-secondary);
  letter-spacing: 3px;
  text-transform: uppercase;
  display: block;
  margin-bottom: 10px;
}

.tag-title {
  font-size: 2.5rem;
  color: var(--text-main); /* 关键 */
  margin: 0 0 15px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 15px;
}

.icon-tag {
  color: var(--primary-color);
}

.tag-desc {
  color: var(--text-regular);
  font-size: 1rem;
}

.loading-box, .empty-box {
  padding: 40px;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.pagination-box {
  display: flex;
  justify-content: center;
  margin-top: 40px;
}

@media (max-width: 768px) {
  .tag-title { font-size: 1.8rem; }
  .tag-header { padding: 30px 10px; }
}
</style>