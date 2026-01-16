<template>
  <div class="dashboard-container">
    <!-- 1. 欢迎区域 -->
    <el-card shadow="never" class="welcome-card">
      <div class="welcome-header">
        <el-avatar :size="64" src="https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png" />
        <div class="welcome-text">
          <h2 class="title">早安，{{ userStore.username || 'Admin' }}，祝你今天过得开心！</h2>
          <p class="subtitle">今天是 {{ currentDate }}，准备好分享什么新知识了吗？</p>
        </div>
      </div>
    </el-card>

    <!-- 2. 数据统计区域 -->
    <el-row :gutter="20" class="mt-4">
      <el-col :span="8" :xs="24">
        <el-card shadow="hover" class="stat-card">
          <template #header>
            <div class="card-header">
              <span>文章总数</span>
              <el-tag type="success">Total</el-tag>
            </div>
          </template>
          <div class="stat-value">{{ totalPosts }}</div>
        </el-card>
      </el-col>
      
      <el-col :span="8" :xs="24">
        <el-card shadow="hover" class="stat-card">
          <template #header>
            <div class="card-header">
              <span>总阅读量</span>
              <el-tag type="warning">Views</el-tag>
            </div>
          </template>
          <div class="stat-value">{{ totalViews }}</div>
        </el-card>
      </el-col>

      <el-col :span="8" :xs="24">
        <el-card shadow="hover" class="stat-card">
          <template #header>
            <div class="card-header">
              <span>系统状态</span>
              <el-tag type="info">Status</el-tag>
            </div>
          </template>
          <div class="stat-value text-green">
            <el-icon><CircleCheckFilled /></el-icon> 运行中
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 3. 快捷操作区域 -->
    <el-row :gutter="20" class="mt-4">
      <el-col :span="12" :xs="24">
        <el-card shadow="hover" header="快捷操作">
          <div class="quick-actions">
            <el-button type="primary" size="large" @click="$router.push('/admin/posts/create')">
              <el-icon class="el-icon--left"><EditPen /></el-icon>
              写新文章
            </el-button>
            <el-button size="large" @click="$router.push('/admin/posts')">
              <el-icon class="el-icon--left"><Management /></el-icon>
              管理内容
            </el-button>
            <el-button size="large" @click="$router.push('/')">
              <el-icon class="el-icon--left"><View /></el-icon>
              查看前台
            </el-button>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="12" :xs="24">
        <el-card shadow="hover" header="最近更新">
          <el-empty v-if="recentPosts.length === 0" description="暂无最近更新" :image-size="60" />
          <ul v-else class="recent-list">
            <li v-for="post in recentPosts" :key="post.ID" class="recent-item">
              <span class="post-title">{{ post.title }}</span>
              <span class="post-date">{{ formatDate(post.created_at) }}</span>
            </li>
          </ul>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useUserStore } from '../../stores/user'
import { getPostList, type Post } from '../../api/post'
import { EditPen, Management, View, CircleCheckFilled } from '@element-plus/icons-vue'

const userStore = useUserStore()
const currentDate = new Date().toLocaleDateString()

const totalPosts = ref(0)
const totalViews = ref(0) // 后端暂时没有聚合统计接口，先保留占位
const recentPosts = ref<Post[]>([])

const formatDate = (dateStr: string) => {
  return new Date(dateStr).toLocaleDateString()
}

onMounted(async () => {
  try {
    // 获取最新的几篇文章用于展示
    const res: any = await getPostList({ page: 1, page_size: 5 })
    recentPosts.value = res.data || []
    totalPosts.value = res.total || 0
    // 这里简单累加一下前5篇的阅读量演示效果，实际应由后端返回 dashboard 数据接口
    totalViews.value = recentPosts.value.reduce((sum, post) => sum + post.view_count, 0)
  } catch (error) {
    console.error('Failed to fetch dashboard data:', error)
  }
})
</script>

<style scoped>
.dashboard-container {
  max-width: 1200px;
}

.welcome-card {
  margin-bottom: 20px;
}

.welcome-header {
  display: flex;
  align-items: center;
  gap: 20px;
}

.welcome-text .title {
  font-size: 20px;
  font-weight: 600;
  margin-bottom: 8px;
  color: #303133;
}

.welcome-text .subtitle {
  color: #909399;
  font-size: 14px;
}

.mt-4 {
  margin-top: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.stat-value {
  font-size: 24px;
  font-weight: bold;
  color: #303133;
  text-align: center;
  padding: 10px 0;
}

.text-green {
  color: #67c23a;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 5px;
}

.quick-actions {
  display: flex;
  gap: 15px;
  flex-wrap: wrap;
}

.recent-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.recent-item {
  display: flex;
  justify-content: space-between;
  padding: 10px 0;
  border-bottom: 1px solid #f0f2f5;
}

.recent-item:last-child {
  border-bottom: none;
}

.post-title {
  color: #606266;
  font-size: 14px;
}

.post-date {
  color: #909399;
  font-size: 12px;
}
</style>