<template>
  <div class="app-container">
    <!-- 1. 顶部导航栏 -->
    <header class="navbar">
      <div class="logo">My Blog</div>
      <div class="user-actions">
        <span class="username">Hello, {{ userStore.username }}</span>
        <el-button type="danger" size="small" @click="handleLogout">退出</el-button>
        <el-button type="primary" size="small" @click="goToCreate">写文章</el-button>
      </div>
    </header>

    <!-- 2. 内容区域 -->
    <main class="content">
      <div v-if="loading" class="loading-state">
        <el-skeleton :rows="5" animated />
      </div>

      <div v-else-if="posts.length === 0" class="empty-state">
        <el-empty description="暂时没有文章，快去写一篇吧！" />
      </div>

      <!-- 文章列表 -->
      <div v-else class="post-list">
        <el-card 
          v-for="post in posts" 
          :key="post.ID" 
          class="post-card" 
          shadow="hover"
        >
          <template #header>
            <div class="post-header">
              <h3 class="post-title">{{ post.title }}</h3>
              <el-tag size="small" effect="plain">{{ post.category?.name || '默认' }}</el-tag>
            </div>
          </template>
          
          <p class="post-summary">{{ post.summary || '暂无摘要' }}</p>
          
          <div class="post-footer">
            <span class="meta-info">
              <el-icon><Calendar /></el-icon> {{ formatDate(post.created_at) }}
            </span>
            <span class="meta-info">
              <el-icon><View /></el-icon> {{ post.view_count }}
            </span>
            
            <div class="actions">
              <el-button type="primary" link>阅读全文</el-button>
              <!-- 只有管理员能删除 -->
              <el-popconfirm title="确定删除这篇文章吗？" @confirm="handleDelete(post.ID)">
                <template #reference>
                  <el-button type="danger" link>删除</el-button>
                </template>
              </el-popconfirm>
            </div>
          </div>
        </el-card>
      </div>

      <!-- 3. 分页器 -->
      <div class="pagination">
        <el-pagination
          background
          layout="prev, pager, next"
          :total="total"
          :page-size="pageSize"
          @current-change="handlePageChange"
        />
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../../stores/user'
import { getPostList, deletePost, type Post } from '../../api/post'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Calendar, View } from '@element-plus/icons-vue'

const router = useRouter()
const userStore = useUserStore()

// 状态定义
const posts = ref<Post[]>([])
const loading = ref(true)
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(5)

// 格式化日期
const formatDate = (dateStr: string) => {
  return new Date(dateStr).toLocaleDateString()
}

// 获取数据
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

// 翻页
const handlePageChange = (page: number) => {
  currentPage.value = page
  fetchData()
}

// 删除文章
const handleDelete = async (id: number) => {
  try {
    await deletePost(id)
    ElMessage.success('删除成功')
    fetchData() // 刷新列表
  } catch (error) {
    console.error(error)
  }
}

// 退出登录
const handleLogout = () => {
  userStore.logout()
  router.push('/login')
}

// 跳转写文章 (还没有这个页面，暂时只是打印)

const goToCreate = () => {
  router.push('/posts/create')
}
// 初始化
onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.app-container {
  min-height: 100vh;
  background-color: #f5f7fa;
}

.navbar {
  background: white;
  padding: 0 20px;
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  box-shadow: 0 2px 4px rgba(0,0,0,0.05);
  position: sticky;
  top: 0;
  z-index: 100;
}

.logo {
  font-size: 20px;
  font-weight: bold;
  color: #409eff;
}

.user-actions {
  display: flex;
  align-items: center;
  gap: 15px;
}

.content {
  max-width: 800px;
  margin: 20px auto;
  padding: 0 20px;
}

.post-card {
  margin-bottom: 20px;
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
}

.post-summary {
  color: #666;
  line-height: 1.6;
  margin: 15px 0;
}

.post-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  color: #999;
  font-size: 13px;
}

.meta-info {
  display: flex;
  align-items: center;
  gap: 5px;
  margin-right: 15px;
}

.pagination {
  display: flex;
  justify-content: center;
  margin-top: 30px;
  padding-bottom: 30px;
}
</style>