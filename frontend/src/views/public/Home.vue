<template>
  <div class="home-container">
    <!-- 1. Hero 头部区域 -->
    <header class="hero-header">
      <div class="hero-content">
        <h1 class="site-title">Zervan's Blog</h1>
        <p class="site-subtitle">Thinking inside the box</p>
        <div class="typing-effect">
          <span>{{ typedText }}</span><span class="cursor">|</span>
        </div>
        
        <div class="scroll-down" @click="scrollToContent">
          <el-icon class="bounce"><ArrowDown /></el-icon>
        </div>
      </div>
      <!-- 动态波浪效果 (SVG) -->
      <div class="waves-container">
        <svg class="waves" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" viewBox="0 24 150 28" preserveAspectRatio="none" shape-rendering="auto">
          <defs>
            <path id="gentle-wave" d="M-160 44c30 0 58-18 88-18s 58 18 88 18 58-18 88-18 58 18 88 18 v44h-352z" />
          </defs>
          <g class="parallax">
            <use xlink:href="#gentle-wave" x="48" y="0" fill="rgba(255,255,255,0.7" />
            <use xlink:href="#gentle-wave" x="48" y="3" fill="rgba(255,255,255,0.5)" />
            <use xlink:href="#gentle-wave" x="48" y="5" fill="rgba(255,255,255,0.3)" />
            <use xlink:href="#gentle-wave" x="48" y="7" fill="#f4f6f8" />
          </g>
        </svg>
      </div>
    </header>

    <!-- 2. 主体内容区 (双栏布局) -->
    <main class="main-wrapper" id="content">
      <el-row :gutter="20">
        <!-- 左侧：文章列表 -->
        <el-col :span="17" :xs="24" class="left-col">
          
          <div v-if="loading" class="loading-box">
             <el-skeleton :rows="5" animated />
          </div>

          <div v-else-if="posts.length === 0" class="empty-box">
            <el-empty description="暂无文章" />
          </div>

          <div v-else class="post-list">
            <article 
              v-for="(post, index) in posts" 
              :key="post.ID" 
              class="post-item"
              :class="{ 'reverse': index % 2 !== 0 }"
              @click="goToDetail(post.ID)"
            >
              <!-- 封面图 -->
              <div class="post-cover">
                <!-- 这里使用随机图作为示例，实际可用 post.cover_url -->
                <img :src="`https://picsum.photos/seed/${post.ID}/400/250`" alt="cover" loading="lazy" />
              </div>
              
              <!-- 文章信息 -->
              <div class="post-info">
                <div class="post-meta-top">
                   <span class="post-date">
                     <el-icon><Calendar /></el-icon> {{ formatDate(post.created_at) }}
                   </span>
                </div>
                
                <h2 class="post-title">{{ post.title }}</h2>
                
                <div class="post-meta-mid">
                  <span class="meta-icon"><el-icon><View /></el-icon> {{ post.view_count }} 热度</span>
                  <span class="meta-icon"><el-icon><ChatLineRound /></el-icon> 0 评论</span>
                  <span class="meta-icon" v-if="post.category">
                    <el-icon><Folder /></el-icon> {{ post.category.name }}
                  </span>
                </div>

                <p class="post-summary">{{ post.summary || '暂无摘要，点击阅读全文...' }}</p>
                
                <div class="post-btn-wrapper">
                   <el-icon><MoreFilled /></el-icon>
                </div>
              </div>
            </article>
          </div>

          <!-- 分页 -->
          <div class="pagination-box">
             <el-pagination
              background
              layout="prev, pager, next"
              :total="total"
              :page-size="pageSize"
              @current-change="handlePageChange"
            />
          </div>
        </el-col>

        <!-- 右侧：侧边栏 -->
        <el-col :span="7" :xs="24" class="right-col hidden-xs-only">
          <div class="sidebar-sticky">
            
            <!-- 个人卡片 -->
            <el-card class="sidebar-card profile-card" shadow="hover">
              <div class="profile-bg">
                <img src="https://picsum.photos/id/1/400/200" alt="bg" />
              </div>
              <div class="profile-content">
                <el-avatar :size="70" src="https://i.pravatar.cc/300?u=zervan" class="profile-avatar" />
                <h3 class="author-name">Zervan</h3>
                <p class="author-desc">天下最普通的人之一</p>
                
                <div class="stats-row">
                   <div class="stat-block">
                     <span class="num">{{ total }}</span>
                     <span class="label">文章</span>
                   </div>
                   <div class="stat-block">
                     <span class="num">4</span>
                     <span class="label">标签</span>
                   </div>
                   <div class="stat-block">
                     <span class="num">1</span>
                     <span class="label">分类</span>
                   </div>
                </div>

                <el-button type="primary" class="follow-btn" round icon="Link" @click="$router.push('/about')">
                  关于我
                </el-button>
                
                <div class="social-links">
                  <el-tooltip content="Github" placement="top">
                    <el-icon class="social-icon"><img src="https://github.com/favicon.ico" width="20"/></el-icon>
                  </el-tooltip>
                  <el-tooltip content="Email" placement="top">
                    <el-icon class="social-icon"><Message /></el-icon>
                  </el-tooltip>
                </div>
              </div>
            </el-card>

            <!-- 公告栏 -->
            <el-card class="sidebar-card" shadow="hover">
              <template #header>
                <div class="card-header">
                  <span><el-icon><Bell /></el-icon> 公告</span>
                </div>
              </template>
              <div class="notice-content">
                <p>欢迎来到 Zervan 的小站！这里记录技术与生活。</p>
              </div>
            </el-card>

            <!-- 标签云 -->
            <el-card class="sidebar-card" shadow="hover">
              <template #header>
                <div class="card-header">
                  <span><el-icon><Collection /></el-icon> 标签</span>
                </div>
              </template>
              <div class="tag-cloud">
                <el-tag class="tag-item" type="info" effect="plain" round>Java</el-tag>
                <el-tag class="tag-item" type="success" effect="plain" round>Golang</el-tag>
                <el-tag class="tag-item" type="warning" effect="plain" round>Vue3</el-tag>
                <el-tag class="tag-item" type="danger" effect="plain" round>Life</el-tag>
                <el-tag class="tag-item" effect="plain" round>Reading</el-tag>
              </div>
            </el-card>

          </div>
        </el-col>
      </el-row>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { getPostList, type Post } from '../../api/post'
import { Calendar, View, ChatLineRound, Folder, ArrowDown, Bell, Collection, Message, MoreFilled } from '@element-plus/icons-vue'

const router = useRouter()
const posts = ref<Post[]>([])
const loading = ref(true)
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(6)

// 打字机效果
const fullText = "Welcome to Zervan's digital garden."
const typedText = ref('')
let typeIndex = 0

const typeWriter = () => {
  if (typeIndex < fullText.length) {
    typedText.value += fullText.charAt(typeIndex)
    typeIndex++
    setTimeout(typeWriter, 100)
  }
}

const formatDate = (dateStr: string) => {
  return new Date(dateStr).toLocaleDateString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit' })
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
  scrollToContent()
}

const goToDetail = (id: number) => {
  router.push(`/post/${id}`)
}

const scrollToContent = () => {
  const content = document.getElementById('content')
  if (content) {
    content.scrollIntoView({ behavior: 'smooth' })
  }
}

onMounted(() => {
  fetchData()
  typeWriter()
})
</script>

<style scoped>
/* 全局容器背景 */
.home-container {
  min-height: 100vh;
  background-color: #f4f6f8;
}

/* --- Hero Section --- */
.hero-header {
  position: relative;
  height: 100vh; /* 全屏高度 */
  min-height: 400px;
  max-height: 600px;
  background: linear-gradient(rgba(0,0,0,0.3), rgba(0,0,0,0.3)), url('https://picsum.photos/seed/bg/1920/1080') center/cover no-repeat;
  color: #fff;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  text-align: center;
}

.hero-content {
  z-index: 2;
  margin-top: -60px;
}

.site-title {
  font-size: 3.5rem;
  font-weight: 700;
  margin: 0;
  text-shadow: 2px 2px 4px rgba(0,0,0,0.5);
  animation: fadeInDown 1s;
}

.site-subtitle {
  font-size: 1.5rem;
  margin: 10px 0 20px;
  opacity: 0.9;
  font-weight: 300;
  animation: fadeInUp 1s;
}

.typing-effect {
  font-size: 1.2rem;
  font-family: monospace;
  background: rgba(0,0,0,0.3);
  padding: 5px 15px;
  border-radius: 4px;
  display: inline-block;
}

.cursor {
  animation: blink 1s infinite;
}

@keyframes blink { 50% { opacity: 0; } }
@keyframes fadeInDown { from { opacity: 0; transform: translateY(-20px); } to { opacity: 1; transform: translateY(0); } }
@keyframes fadeInUp { from { opacity: 0; transform: translateY(20px); } to { opacity: 1; transform: translateY(0); } }
@keyframes bounce { 0%, 20%, 50%, 80%, 100% {transform: translateY(0);} 40% {transform: translateY(-10px);} 60% {transform: translateY(-5px);} }

.scroll-down {
  position: absolute;
  bottom: 50px;
  cursor: pointer;
  font-size: 2rem;
  z-index: 5;
  animation: bounce 2s infinite;
}

/* --- Waves --- */
.waves-container {
  position: absolute;
  bottom: 0;
  left: 0;
  width: 100%;
  overflow: hidden;
  line-height: 0;
}
.waves {
  width: 100%;
  height: 60px;
}
.parallax > use {
  animation: move-forever 25s cubic-bezier(.55,.5,.45,.5) infinite;
}
.parallax > use:nth-child(1) { animation-delay: -2s; animation-duration: 7s; }
.parallax > use:nth-child(2) { animation-delay: -3s; animation-duration: 10s; }
.parallax > use:nth-child(3) { animation-delay: -4s; animation-duration: 13s; }
.parallax > use:nth-child(4) { animation-delay: -5s; animation-duration: 20s; }
@keyframes move-forever {
  0% { transform: translate3d(-90px,0,0); }
  100% { transform: translate3d(85px,0,0); }
}

/* --- Main Layout --- */
.main-wrapper {
  max-width: 1200px;
  margin: 0 auto;
  padding: 40px 20px;
  animation: fadeInUp 0.8s ease-in-out;
}

/* --- Post List --- */
.post-item {
  display: flex;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 4px 12px rgba(0,0,0,0.05);
  margin-bottom: 25px;
  overflow: hidden;
  height: 280px;
  cursor: pointer;
  transition: all 0.3s;
}

.post-item:hover {
  box-shadow: 0 8px 24px rgba(0,0,0,0.12);
  transform: translateY(-4px);
}

/* 封面图在左 */
.post-cover {
  width: 45%;
  overflow: hidden;
}

.post-cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.5s;
}

.post-item:hover .post-cover img {
  transform: scale(1.1);
}

/* 内容在右 */
.post-info {
  width: 55%;
  padding: 30px 40px;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.post-meta-top {
  color: #888;
  font-size: 13px;
  margin-bottom: 10px;
}

.post-title {
  font-size: 1.6rem;
  font-weight: 700;
  color: #333;
  margin: 0 0 15px;
  line-height: 1.3;
  transition: color 0.3s;
}

.post-item:hover .post-title {
  color: #409eff;
}

.post-meta-mid {
  display: flex;
  gap: 15px;
  font-size: 13px;
  color: #999;
  margin-bottom: 15px;
}

.meta-icon {
  display: flex;
  align-items: center;
  gap: 4px;
}

.post-summary {
  color: #666;
  line-height: 1.6;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  margin-bottom: 20px;
}

.post-btn-wrapper {
  margin-top: auto;
  font-size: 20px;
  color: #ccc;
}

/* 偶数个文章反转布局：图在右 */
.post-item.reverse {
  flex-direction: row-reverse;
}

/* --- Sidebar --- */
.right-col {
  padding-left: 10px;
}

.sidebar-sticky {
  position: sticky;
  top: 80px; /* 留出 header 空间 */
}

.sidebar-card {
  border-radius: 8px;
  margin-bottom: 20px;
  border: none;
  box-shadow: 0 4px 12px rgba(0,0,0,0.05);
}

.card-header {
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 5px;
}

/* 个人资料卡片 */
.profile-card :deep(.el-card__body) {
  padding: 0;
}

.profile-bg {
  height: 120px;
  overflow: hidden;
}
.profile-bg img { width: 100%; height: 100%; object-fit: cover; }

.profile-content {
  text-align: center;
  padding: 0 20px 20px;
  margin-top: -35px;
}

.profile-avatar {
  border: 4px solid #fff;
  background: #fff;
}

.author-name {
  margin: 10px 0 5px;
  font-size: 1.2rem;
}

.author-desc {
  font-size: 13px;
  color: #888;
  margin-bottom: 15px;
}

.stats-row {
  display: flex;
  justify-content: space-around;
  margin-bottom: 20px;
}

.stat-block {
  display: flex;
  flex-direction: column;
}

.stat-block .num { font-weight: bold; font-size: 16px; color: #333; }
.stat-block .label { font-size: 12px; color: #999; }

.follow-btn {
  width: 100%;
  margin-bottom: 15px;
}

.social-links {
  display: flex;
  justify-content: center;
  gap: 15px;
}

.social-icon {
  font-size: 20px;
  color: #666;
  cursor: pointer;
  transition: color 0.3s;
}

.social-icon:hover { color: #409eff; }

/* 标签云 */
.tag-cloud {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}
.tag-item { cursor: pointer; }

/* 响应式适配 */
@media (max-width: 768px) {
  .post-item {
    flex-direction: column !important;
    height: auto;
  }
  .post-cover, .post-info {
    width: 100%;
  }
  .post-cover {
    height: 200px;
  }
  .site-title {
    font-size: 2.5rem;
  }
}
</style>