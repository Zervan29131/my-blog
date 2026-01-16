<template>
  <div class="home-content">
    <!-- 加载状态 -->
    <div v-if="loading" class="loading-state">
      <el-skeleton :rows="5" animated />
    </div>

    <!-- 空状态 -->
    <div v-else-if="posts.length === 0" class="empty-state">
      <el-empty description="暂时没有文章" />
    </div>

    <!-- 文章列表 -->
    <div v-else class="post-list">
      <el-card 
        v-for="post in posts" 
        :key="post.ID" 
        class="post-card" 
        shadow="hover"
        @click="goToDetail(post.ID)"
      >
        <template #header>
          <div class="post-header">
            <h3 class="post-title">{{ post.title }}</h3>
            <div class="post-tags">
               <el-tag v-if="post.category" size="small" effect="plain">{{ post.category.name }}</el-tag>
            </div>
          </div>
        </template>
        
        <p class="post-summary">{{ post.summary || '暂无摘要' }}</p>
        
        <div class="post-footer">
          <div class="meta-group">
            <span class="meta-info">
              <el-icon><Calendar /></el-icon> {{ formatDate(post.created_at) }}
            </span>
            <span class="meta-info">
              <el-icon><View /></el-icon> {{ post.view_count }}
            </span>
          </div>
          
          <div class="actions">
            <!-- 🟢 修复点：添加点击事件，使用 .stop 阻止冒泡防止触发卡片点击 -->
            <el-button type="primary" link @click.stop="goToDetail(post.ID)">阅读全文</el-button>
          </div>
        </div>
      </el-card>
    </div>

    <!-- 分页器 -->
    <div class="pagination">
      <el-pagination
        background
        layout="prev, pager, next"
        :total="total"
        :page-size="pageSize"
        @current-change="handlePageChange"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getPostList, type Post } from '../../api/post'
import { Calendar, View } from '@element-plus/icons-vue'

const router = useRouter()
const posts = ref<Post[]>([])
const loading = ref(true)
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(5)

const formatDate = (dateStr: string) => {
  return new Date(dateStr).toLocaleDateString()
}

// 获取文章列表
const fetchData = async () => {
  loading.value = true
  try {
    const res: any = await getPostList({
      page: currentPage.value,
      page_size: pageSize.value
    })
    posts.value = res.data || []
    total.value = res.total || 0
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const handlePageChange = (page: number) => {
  currentPage.value = page
  fetchData()
}

// 🟢 跳转详情页逻辑
const goToDetail = (id: number) => {
  router.push(`/post/${id}`)
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.home-content {
  padding: 20px 0;
}
.post-card {
  margin-bottom: 20px;
  cursor: pointer;
  transition: transform 0.2s;
}
.post-card:hover {
  transform: translateY(-2px);
}
.post-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.post-title {
  margin: 0;
  font-size: 18px;
  font-weight: bold;
  color: #303133;
}
.post-summary {
  color: #606266;
  line-height: 1.6;
  margin: 15px 0;
  /* 限制摘要显示行数 */
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 3;
  overflow: hidden;
}
.post-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  color: #909399;
  font-size: 13px;
}
.meta-group {
  display: flex;
  gap: 15px;
}
.meta-info {
  display: flex;
  align-items: center;
  gap: 5px;
}
.pagination {
  display: flex;
  justify-content: center;
  margin-top: 30px;
}
</style>