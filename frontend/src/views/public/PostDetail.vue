<template>
  <div class="post-detail-wrapper">
    <!-- 1. 加载状态 -->
    <div v-if="loading" class="loading-container">
      <el-skeleton :rows="10" animated />
    </div>

    <!-- 2. 404 状态：加载结束且无数据 -->
    <div v-else-if="!loading && !post" class="not-found">
      <el-empty description="抱歉，文章不存在或已被删除">
        <el-button type="primary" @click="$router.push('/')">返回首页</el-button>
      </el-empty>
    </div>

    <!-- 3. 文章主体：确保 post 存在才渲染 -->
    <div v-else-if="post">
      <article class="article-content">
        <!-- 头部元数据 -->
        <header class="post-header">
          <h1 class="post-title">{{ post?.title }}</h1>
          
          <div class="post-meta">
            <span class="meta-item">
              <el-icon><User /></el-icon> {{ post?.author?.username || 'Admin' }}
            </span>
            <span class="meta-item">
              <el-icon><Calendar /></el-icon> {{ formatDate(post?.created_at) }}
            </span>
            <span class="meta-item">
              <el-icon><Folder /></el-icon> {{ post?.category?.name || '默认分类' }}
            </span>
            <span class="meta-item">
              <el-icon><View /></el-icon> {{ post?.view_count }} 阅读
            </span>
          </div>

          <!-- 标签 -->
          <div v-if="post?.tags && post.tags.length > 0" class="post-tags">
            <el-tag 
              v-for="tag in post.tags" 
              :key="tag.ID" 
              size="small" 
              effect="plain" 
              round
              class="tag-item"
            >
              # {{ tag.name }}
            </el-tag>
          </div>
        </header>

        <el-divider />

        <!-- Markdown 内容渲染区 -->
        <div class="markdown-body" v-html="htmlContent"></div>
      </article>

      <!-- 4. 作者卡片 -->
      <div class="author-section">
        <AuthorCard 
          :name="post?.author?.username" 
          bio="热爱编程，热爱生活的全栈开发者。" 
        />
      </div>

      <!-- 5. 评论区 -->
      <Comment :post-id="post.ID" />

      <!-- 6. 底部导航 -->
      <div class="post-footer-nav">
         <el-button @click="$router.push('/')">← 返回列表</el-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getPost, type Post } from '../../api/post'
import { User, Calendar, Folder, View } from '@element-plus/icons-vue'
import MarkdownIt from 'markdown-it'
import hljs from 'highlight.js'
// 确保安装了 highlight.js: npm install highlight.js
import 'highlight.js/styles/github-dark.css' 

// 引入组件
import AuthorCard from '../../components/AuthorCard.vue'
import Comment from '../../components/Comment.vue'

const route = useRoute()
const router = useRouter()
const post = ref<Post | null>(null)
const loading = ref(true)
const htmlContent = ref('')

// --- 初始化 Markdown 解析器 ---
// 1. 先实例化，避免在配置中循环引用 md 自身
const md = new MarkdownIt({
  html: true,       // 允许 HTML 标签
  linkify: true,    // 自动识别 URL
  typographer: true
})

// 2. 后置设置高亮逻辑，并添加明确的类型注解
md.set({
  highlight: function (str: string, lang: string): string {
    if (lang && hljs.getLanguage(lang)) {
      try {
        return `<pre class="hljs"><code>${hljs.highlight(str, { language: lang, ignoreIllegals: true }).value}</code></pre>`
      } catch (__) {}
    }
    // 默认回退：现在 md 已经初始化完成，可以安全使用 md.utils
    return `<pre class="hljs"><code>${md.utils.escapeHtml(str)}</code></pre>`
  }
})

// --- 格式化日期 ---
const formatDate = (dateStr?: string) => {
  if (!dateStr) return ''
  try {
    return new Date(dateStr).toLocaleDateString('zh-CN', {
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    })
  } catch (e) {
    return dateStr
  }
}

// --- 获取数据 ---
onMounted(async () => {
  const id = Number(route.params.id)
  
  // 如果 ID 无效，停止加载
  if (!id || isNaN(id)) {
    loading.value = false
    return
  }

  try {
    const res: any = await getPost(id)
    if (res && res.data) {
      post.value = res.data
      // 核心步骤：将 Markdown 文本渲染成 HTML
      if (post.value && post.value.content) {
        htmlContent.value = md.render(post.value.content)
      }
    } else {
      post.value = null
    }
  } catch (error) {
    console.error('获取文章详情失败:', error)
    post.value = null
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.post-detail-wrapper {
  max-width: 720px;
  margin: -60px auto 40px;
  padding: 50px;
  min-height: 60vh;
  /* 🔴 1. 背景色变量化 */
  background-color: var(--bg-content);
  border-radius: 12px;
  /* 🔴 2. 阴影变量化 */
  box-shadow: var(--shadow-light);
  position: relative;
  z-index: 10;
  transition: background-color 0.3s, box-shadow 0.3s;
}

.post-header {
  text-align: center;
  margin-bottom: 30px;
}

.post-title {
  font-size: 2.5rem;
  /* 🔴 3. 主标题颜色变量化 */
  color: var(--text-main);
  margin-bottom: 25px;
  font-weight: 800;
  line-height: 1.2;
  transition: color 0.3s;
}

.post-meta {
  display: flex;
  justify-content: center;
  flex-wrap: wrap;
  gap: 20px;
  /* 🔴 4. 辅助文字颜色变量化 */
  color: var(--text-secondary);
  font-size: 0.9rem;
  margin-bottom: 40px;
  /* 🔴 5. 边框颜色变量化 */
  border-bottom: 1px solid var(--border-color);
  padding-bottom: 20px;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 6px;
}

.post-tags {
  margin-top: 20px;
  display: flex;
  justify-content: center;
  gap: 10px;
}

.author-section {
  margin-top: 40px;
}

.post-footer-nav {
  margin-top: 30px;
  padding-top: 30px;
  text-align: center;
}

@media (max-width: 768px) {
  .post-detail-wrapper {
    padding: 20px;
    margin-top: 20px;
    margin-left: 10px;
    margin-right: 10px;
  }
  .post-title {
    font-size: 1.8rem;
  }
}

/* =========================================
   Markdown 渲染样式定制 (模拟 GitHub 风格)
   ========================================= */
.markdown-body {
  font-size: 17px;
  line-height: 1.8;
  /* 🔴 6. 正文颜色变量化 */
  color: var(--text-regular);
  font-family: -apple-system,BlinkMacSystemFont,Segoe UI,Helvetica,Arial,sans-serif;
  transition: color 0.3s;
}

.markdown-body :deep(h1),
.markdown-body :deep(h2),
.markdown-body :deep(h3) {
  margin-top: 24px;
  margin-bottom: 16px;
  font-weight: 600;
  line-height: 1.25;
  /* 🔴 7. 标题颜色变量化 */
  color: var(--text-main);
  border-bottom: 1px solid var(--border-color);
  padding-bottom: 0.3em;
}

.markdown-body :deep(p) {
  margin-top: 0;
  margin-bottom: 16px;
}

.markdown-body :deep(blockquote) {
  padding: 0 1em;
  /* 🔴 8. 引用块颜色变量化 */
  color: var(--text-secondary);
  border-left: 0.25em solid var(--border-color);
  margin: 0 0 16px 0;
  background-color: transparent; 
}

.markdown-body :deep(ul),
.markdown-body :deep(ol) {
  padding-left: 2em;
  margin-bottom: 16px;
}

.markdown-body :deep(img) {
  max-width: 100%;
  box-sizing: border-box;
  /* 🔴 9. 图片背景变量化 */
  background-color: var(--bg-content);
  border-radius: 4px;
  box-shadow: var(--shadow-light);
}

/* 代码块样式 */
.markdown-body :deep(pre) {
  padding: 16px;
  overflow: auto;
  font-size: 85%;
  line-height: 1.45;
  background-color: #1e1e1e;
  border-radius: 6px;
  margin-bottom: 16px;
}

.markdown-body :deep(code) {
  padding: 0.2em 0.4em;
  margin: 0;
  font-size: 85%;
  /* 🔴 10. 行内代码背景微调 */
  background-color: rgba(128, 128, 128, 0.1); 
  border-radius: 3px;
  font-family: SFMono-Regular,Consolas,Liberation Mono,Menlo,monospace;
}

.markdown-body :deep(pre code) {
  background-color: transparent;
  padding: 0;
  color: inherit;
}

.markdown-body :deep(a) {
  color: var(--primary-color);
  text-decoration: none;
}
.markdown-body :deep(a:hover) {
  text-decoration: underline;
}
</style>