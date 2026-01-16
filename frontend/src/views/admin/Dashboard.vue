<template>
  <div class="dashboard-container">
    <!-- 1. 欢迎区域 -->
    <el-card shadow="never" class="welcome-card">
      <div class="welcome-header">
        <el-avatar :size="64" :src="avatarUrl" />
        <div class="welcome-text">
          <h2 class="title">早安，{{ userStore.username || 'Admin' }}</h2>
          <p class="subtitle">今天是 {{ currentDate }}。喝杯咖啡，开始创作吧！</p>
        </div>
      </div>
    </el-card>

    <!-- 2. 核心指标 -->
    <el-row :gutter="20" class="mt-4">
      <el-col :span="6" :xs="12">
        <div class="stat-item bg-white">
          <div class="stat-icon blue">
            <el-icon><Document /></el-icon>
          </div>
          <div class="stat-info">
            <div class="label">文章总数</div>
            <div class="count">{{ totalPosts }}</div>
          </div>
        </div>
      </el-col>
      
      <el-col :span="6" :xs="12">
        <div class="stat-item bg-white">
          <div class="stat-icon orange">
            <el-icon><View /></el-icon>
          </div>
          <div class="stat-info">
            <div class="label">总阅读量</div>
            <div class="count">{{ totalViews }}</div>
          </div>
        </div>
      </el-col>

      <el-col :span="6" :xs="12">
        <div class="stat-item bg-white">
          <div class="stat-icon green">
            <el-icon><ChatLineSquare /></el-icon>
          </div>
          <div class="stat-info">
            <div class="label">评论总数</div>
            <div class="count">0</div>
          </div>
        </div>
      </el-col>

      <el-col :span="6" :xs="12">
        <div class="stat-item bg-white">
          <div class="stat-icon purple">
            <el-icon><Monitor /></el-icon>
          </div>
          <div class="stat-info">
            <div class="label">系统状态</div>
            <div class="count text-small">运行中</div>
          </div>
        </div>
      </el-col>
    </el-row>

    <!-- 3. 最近动态 & 系统信息 -->
    <el-row :gutter="20" class="mt-4">
      <!-- 左侧：最近发布的文章 (占据更多宽度) -->
      <el-col :span="16" :xs="24">
        <el-card shadow="hover" class="recent-card">
          <template #header>
            <div class="card-header">
              <span>最近发布</span>
              <el-button link type="primary" @click="$router.push('/admin/posts')">全部文章</el-button>
            </div>
          </template>
          
          <el-table :data="recentPosts" style="width: 100%" :show-header="false">
            <el-table-column prop="title" label="标题">
              <template #default="scope">
                 <span class="post-link" @click="$router.push(`/admin/posts/${scope.row.ID}/edit`)">
                   {{ scope.row.title }}
                 </span>
              </template>
            </el-table-column>
            <el-table-column prop="category.name" label="分类" width="120">
              <template #default="scope">
                <el-tag size="small" effect="plain">{{ scope.row.category?.name || '默认' }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="created_at" label="时间" width="120" align="right">
              <template #default="scope">
                <span class="date-text">{{ formatDate(scope.row.created_at) }}</span>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
      
      <!-- 右侧：系统信息/公告 -->
      <el-col :span="8" :xs="24">
        <el-card shadow="hover">
          <template #header>
            <span>关于系统</span>
          </template>
          <div class="system-info">
            <div class="info-row">
              <span class="label">前端版本</span>
              <span class="value">Vue 3.4 + Vite 5</span>
            </div>
            <div class="info-row">
              <span class="label">后端版本</span>
              <span class="value">Go 1.21 + Gin</span>
            </div>
            <div class="info-row">
              <span class="label">数据库</span>
              <span class="value">MySQL 8.0</span>
            </div>
            <el-divider />
            <div class="quote">
              "Code is like humor. When you have to explain it, it’s bad."
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useUserStore } from '../../stores/user'
import { getPostList, type Post } from '../../api/post'
import { Document, View, ChatLineSquare, Monitor } from '@element-plus/icons-vue'

const userStore = useUserStore()
const currentDate = new Date().toLocaleDateString()

const totalPosts = ref(0)
const totalViews = ref(0)
const recentPosts = ref<Post[]>([])

const avatarUrl = computed(() => `https://ui-avatars.com/api/?name=${userStore.username}&background=409eff&color=fff`)

const formatDate = (dateStr: string) => {
  return new Date(dateStr).toLocaleDateString()
}

onMounted(async () => {
  try {
    const res: any = await getPostList({ page: 1, page_size: 6 })
    recentPosts.value = res.data || []
    totalPosts.value = res.total || 0
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
  border: none;
  background: linear-gradient(120deg, #e0c3fc 0%, #8ec5fc 100%);
  color: #fff;
}
.welcome-header {
  display: flex;
  align-items: center;
  gap: 20px;
}
.welcome-text .title {
  margin: 0 0 8px 0;
  font-size: 24px;
}
.welcome-text .subtitle {
  margin: 0;
  opacity: 0.9;
}

.mt-4 { margin-top: 20px; }

/* 统计卡片样式 */
.stat-item {
  display: flex;
  align-items: center;
  padding: 20px;
  border-radius: 8px;
  cursor: pointer;
  transition: transform 0.3s;
  box-shadow: 0 2px 12px 0 rgba(0,0,0,0.05);
}
.stat-item:hover {
  transform: translateY(-5px);
}
.bg-white { background: #fff; }

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  margin-right: 15px;
}
.blue { background: #e6f7ff; color: #1890ff; }
.orange { background: #fff7e6; color: #fa8c16; }
.green { background: #f6ffed; color: #52c41a; }
.purple { background: #f9f0ff; color: #722ed1; }

.stat-info .label { color: #909399; font-size: 13px; margin-bottom: 5px; }
.stat-info .count { font-size: 20px; font-weight: bold; color: #303133; }
.text-small { font-size: 16px; }

/* 最近发布列表 */
.card-header { display: flex; justify-content: space-between; align-items: center; }
.post-link { cursor: pointer; color: #606266; }
.post-link:hover { color: #409eff; }
.date-text { color: #909399; font-size: 12px; }

/* 系统信息 */
.system-info .info-row {
  display: flex;
  justify-content: space-between;
  padding: 10px 0;
  border-bottom: 1px solid #f0f2f5;
  font-size: 14px;
}
.system-info .label { color: #909399; }
.quote {
  margin-top: 15px;
  font-style: italic;
  color: #999;
  font-size: 12px;
}
</style>