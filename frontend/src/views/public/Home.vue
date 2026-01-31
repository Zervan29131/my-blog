<template>
  <div class="home-container">
    <!-- 1. Hero 头部区域 -->
    <header class="hero-header">
      <div class="hero-overlay"></div>
      <div class="hero-content">
        <h1 class="site-title">这是存在于某处的幻想世界</h1>
        <!-- 打字机效果 -->
        <div class="typing-effect">
          <span>{{ typedText }}</span><span class="cursor">|</span>
        </div>
        
        <!-- 赞赏入口 -->
        <div class="donate-wrapper">
          <el-button 
            round 
            class="donate-btn" 
            @click="handleDonateClick"
          >
            <el-icon style="margin-right: 6px;"><Coffee /></el-icon> 赞赏支持
          </el-button>
        </div>

        <!-- 向下滚动按钮 -->
        <div class="scroll-down" @click="scrollToContent">
          <el-icon class="bounce"><ArrowDown /></el-icon>
        </div>
      </div>
    </header>

    <!-- 2. 主体内容区 -->
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
                <p class="author-desc">Thinking inside the box</p>
                
                <div class="stats-row">
                   <div class="stat-block">
                     <span class="num">{{ total }}</span>
                     <span class="label">文章</span>
                   </div>
                   <div class="stat-block">
                     <span class="num">4</span>
                     <span class="label">标签</span>
                   </div>
                </div>

                <el-button type="primary" class="follow-btn" round icon="Link" @click="$router.push('/about')">
                  关于我
                </el-button>
                
                <div class="social-links">
                  <el-icon class="social-icon"><Message /></el-icon>
                  <el-icon class="social-icon"><Share /></el-icon>
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
                <p style="margin-top: 10px; font-size: 12px; opacity: 0.8">
                  <el-tag size="small" effect="plain">Update</el-tag> 适配暗黑模式
                </p>
              </div>
            </el-card>

            <!-- 🟢 新增：网站资讯 -->
            <el-card class="sidebar-card" shadow="hover">
              <template #header>
                <div class="card-header">
                  <span><el-icon><DataLine /></el-icon> 网站资讯</span>
                </div>
              </template>
              <div class="site-info-content">
                <div class="info-item">
                  <span>文章数目 :</span>
                  <span class="value">34</span>
                </div>
                <div class="info-item">
                  <span>已运行时间 :</span>
                  <span class="value">2483 天</span>
                </div>
                <div class="info-item">
                  <span>本站总字数 :</span>
                  <span class="value">86.6k</span>
                </div>
                <div class="info-item">
                  <span>本站访客数 :</span>
                  <span class="value">24694</span>
                </div>
                <div class="info-item">
                  <span>本站总访问量 :</span>
                  <span class="value">35924</span>
                </div>
                <div class="info-item">
                  <span>最后更新时间 :</span>
                  <span class="value">10 个月前</span>
                </div>
              </div>
            </el-card>

          </div>
        </el-col>
      </el-row>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { getPostList, type Post } from '../../api/post'
import { ArrowDown, Bell, Message, Share, Link, Coffee, DataLine } from '@element-plus/icons-vue'
import PostCard from '../../components/PostCard.vue'

const router = useRouter()
const route = useRoute()
const posts = ref<Post[]>([])
const loading = ref(true)
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(6)

// 打字机效果
const fullText = "Be the change you want to see in the World."
const typedText = ref('')
let typeIndex = 0

const typeWriter = () => {
  if (typeIndex < fullText.length) {
    typedText.value += fullText.charAt(typeIndex)
    typeIndex++
    setTimeout(typeWriter, 100)
  }
}

const fetchData = async () => {
  loading.value = true
  try {
    const searchQuery = route.query.q ? String(route.query.q) : ''
    const res: any = await getPostList({
      page: currentPage.value,
      page_size: pageSize.value,
      q: searchQuery
    })
    posts.value = res.data || []
    total.value = res.total || 0
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const scrollToContent = () => {
  const content = document.getElementById('content')
  if (content) {
    content.scrollIntoView({ behavior: 'smooth' })
  }
}

const handleDonateClick = () => {
  router.push('/donate')
}

const handlePageChange = (page: number) => {
  currentPage.value = page
  fetchData()
  scrollToContent()
}

const goToDetail = (id: number) => {
  router.push(`/post/${id}`)
}

watch(() => route.query, (newVal, oldVal) => {
  if (newVal.q !== oldVal?.q) {
    currentPage.value = 1
    fetchData()
  }
})

onMounted(() => {
  fetchData()
  typeWriter()
})
</script>

<style scoped>
.home-container {
  min-height: 100vh;
  background-color: var(--bg-color);
  transition: background-color 0.3s;
}

/* --- Hero Section --- */
.hero-header {
  position: relative;
  height: 100vh;
  min-height: 400px;
  max-height: 600px;
  /* 设置背景图 */
  background: url('../../assets/bg.jpg') center/cover no-repeat;
  color: #fff;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  text-align: center;
  /* 移除 z-index，避免遮挡 */
}

/* 添加一个遮罩层，确保文字清晰，且不影响点击 */
.hero-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.4);
  pointer-events: none; /* 关键：允许点击穿透 */
}

.hero-content {
  position: relative;
  z-index: 10; /* 确保内容在遮罩之上 */
  margin-top: -60px;
  pointer-events: auto; /* 确保内容可点击 */
}

.site-title {
  font-size: 3.5rem;
  font-weight: 700;
  margin: 0 0 20px;
  text-shadow: 2px 2px 4px rgba(0,0,0,0.5);
  animation: fadeInDown 1s;
}

.typing-effect {
  font-size: 1.2rem;
  font-family: 'Courier New', Courier, monospace;
  background: rgba(0,0,0,0.4);
  padding: 8px 16px;
  border-radius: 4px;
  display: inline-block;
  backdrop-filter: blur(4px);
  margin-bottom: 25px;
}

.donate-wrapper {
  animation: fadeInDown 1.5s;
}

.donate-btn {
  background: rgba(255, 255, 255, 0.2);
  border: 1px solid rgba(255, 255, 255, 0.4);
  color: #fff;
  padding: 18px 25px;
  font-weight: 500;
  transition: all 0.3s;
  cursor: pointer;
}
.donate-btn:hover {
  background: rgba(255, 255, 255, 0.9);
  color: #333;
  transform: translateY(-2px);
  box-shadow: 0 4px 15px rgba(0,0,0,0.3);
}

.cursor { animation: blink 1s infinite; }
@keyframes blink { 50% { opacity: 0; } }
@keyframes fadeInDown { from { opacity: 0; transform: translateY(-20px); } to { opacity: 1; transform: translateY(0); } }
@keyframes bounce { 0%, 20%, 50%, 80%, 100% {transform: translateY(0);} 40% {transform: translateY(-10px);} 60% {transform: translateY(-5px);} }

.scroll-down {
  position: absolute;
  bottom: 80px;
  cursor: pointer;
  font-size: 2rem;
  z-index: 10;
  animation: bounce 2s infinite;
  color: #fff;
  opacity: 0.8;
}

/* --- Main Layout --- */
.main-wrapper {
  max-width: 1200px;
  margin: 0 auto;
  padding: 40px 20px;
  /* 🟢 修复核心：确保主内容在所有背景层之上 */
  position: relative;
  z-index: 20;
  pointer-events: auto; 
}

.pagination-box {
  display: flex;
  justify-content: center;
  margin-top: 40px;
}

/* --- Sidebar --- */
.right-col { padding-left: 10px; }
.sidebar-sticky { position: sticky; top: 80px; }

.sidebar-card {
  background: var(--bg-content);
  border-radius: 8px;
  margin-bottom: 20px;
  border: 1px solid var(--border-color);
  box-shadow: var(--shadow-light);
  color: var(--text-main);
  transition: all 0.3s;
}

:global(html.dark) .sidebar-card {
  background-color: #1d1e1f;
  border-color: #363637;
}

.card-header {
  font-weight: 600;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 5px;
  color: var(--text-main);
}

.profile-card :deep(.el-card__body) { padding: 0; }
.profile-bg { height: 120px; overflow: hidden; }
.profile-bg img { width: 100%; height: 100%; object-fit: cover; }
.profile-content { text-align: center; padding: 0 20px 20px; margin-top: -35px; }

.profile-avatar {
  border: 4px solid var(--bg-content);
  background: #fff;
  transition: border-color 0.3s;
}
:global(html.dark) .profile-avatar {
  border-color: #1d1e1f;
}

.author-name { margin: 10px 0 5px; font-size: 1.2rem; color: var(--text-main); }
.author-desc { font-size: 13px; color: var(--text-secondary); margin-bottom: 15px; }

.stats-row { display: flex; justify-content: space-around; margin-bottom: 20px; }
.stat-block { display: flex; flex-direction: column; }
.stat-block .num { font-weight: bold; font-size: 16px; color: var(--text-main); }
.stat-block .label { font-size: 12px; color: var(--text-secondary); }

.follow-btn { width: 100%; margin-bottom: 15px; }

.social-links { display: flex; justify-content: center; gap: 15px; }
.social-icon { font-size: 20px; color: var(--text-secondary); cursor: pointer; transition: color 0.3s; }
.social-icon:hover { color: var(--primary-color); }

.notice-content { color: var(--text-regular); padding: 0 10px; line-height: 1.6; }

/* 网站资讯样式 */
.site-info-content {
  padding: 0 10px;
}
.info-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 14px;
  color: var(--text-regular);
  margin-bottom: 10px;
  padding-bottom: 5px;
  border-bottom: 1px dashed var(--border-color);
}
.info-item:last-child {
  margin-bottom: 0;
  border-bottom: none;
}
.info-item span:first-child {
  color: var(--text-secondary);
}
.info-item .value {
  font-weight: 500;
  color: var(--text-main);
  font-family: monospace; 
}

@media (max-width: 768px) {
  .site-title { font-size: 2.5rem; }
  .main-wrapper { padding: 20px 10px; }
}
</style>