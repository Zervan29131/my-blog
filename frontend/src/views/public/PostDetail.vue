<template>
  <div class="post-detail-wrapper">
    <!-- 1. 加载状态 -->
    <div v-if="loading" class="loading-container">
      <el-skeleton :rows="10" animated />
    </div>

    <!-- 2. 404 状态 -->
    <div v-else-if="!post" class="not-found">
      <el-empty description="抱歉，文章不存在或已被删除">
        <el-button type="primary" @click="$router.push('/')">返回首页</el-button>
      </el-empty>
    </div>

    <!-- 3. 文章主体 -->
    <article v-else class="article-content">
      <!-- 头部元数据 -->
      <header class="post-header">
        <h1 class="post-title">{{ post.title }}</h1>
        
        <div class="post-meta">
          <span class="meta-item">
            <el-icon><User /></el-icon> {{ post.author?.username || 'Admin' }}
          </span>
          <span class="meta-item">
            <el-icon><Calendar /></el-icon> {{ formatDate(post.created_at) }}
          </span>
          <span class="meta-item">
            <el-icon><Folder /></el-icon> {{ post.category?.name || '默认分类' }}
          </span>
          <span class="meta-item">
            <el-icon><View /></el-icon> {{ post.view_count }} 阅读
          </span>
        </div>

        <!-- 标签 -->
        <div v-if="post.tags && post.tags.length > 0" class="post-tags">
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
      <!-- 注意：这里使用了 v-html，样式需要通过 deep 选择器穿透 -->
      <div class="markdown-body" v-html="htmlContent"></div>
    </article>

    <!-- 4. 底部导航 (示例) -->
    <div v-if="post" class="post-footer-nav">
       <el-button @click="$router.push('/')">← 返回列表</el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { getPost, type Post } from '../../api/post'
import { User, Calendar, Folder, View } from '@element-plus/icons-vue'
import MarkdownIt from 'markdown-it'
import hljs from 'highlight.js'
// 引入一款高亮代码的 CSS 主题，你可以换成 atom-one-dark 等其他主题
import 'highlight.js/styles/github-dark.css' 

const route = useRoute()
const post = ref<Post | null>(null)
const loading = ref(true)
const htmlContent = ref('')

// --- 初始化 Markdown 解析器 ---
const md = new MarkdownIt({
  html: true,       // 允许 HTML 标签
  linkify: true,    // 自动识别 URL
  typographer: true,
  // 配置代码高亮
  highlight: function (str, lang) {
    if (lang && hljs.getLanguage(lang)) {
      try {
        return `<pre class="hljs"><code>${hljs.highlight(str, { language: lang, ignoreIllegals: true }).value}</code></pre>`
      } catch (__) {}
    }
    // 默认回退
    return `<pre class="hljs"><code>${md.utils.escapeHtml(str)}</code></pre>`
  }
})

// --- 格式化日期 ---
const formatDate = (dateStr: string) => {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}

// --- 获取数据 ---
onMounted(async () => {
  const id = Number(route.params.id)
  if (!id) {
    loading.value = false
    return
  }

  try {
    const res: any = await getPost(id)
    post.value = res.data
    // 核心步骤：将 Markdown 文本渲染成 HTML
    if (post.value && post.value.content) {
      htmlContent.value = md.render(post.value.content)
    }
  } catch (error) {
    console.error('获取文章详情失败:', error)
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.post-detail-wrapper {
  max-width: 800px;
  margin: 0 auto;
  padding: 40px 20px;
  min-height: 60vh;
  background-color: #fff;
}

.post-header {
  text-align: center;
  margin-bottom: 30px;
}

.post-title {
  font-size: 2.2rem;
  color: #333;
  margin-bottom: 20px;
  font-weight: 700;
}

.post-meta {
  display: flex;
  justify-content: center;
  flex-wrap: wrap;
  gap: 20px;
  color: #909399;
  font-size: 0.9rem;
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

.post-footer-nav {
  margin-top: 50px;
  padding-top: 30px;
  border-top: 1px solid #eee;
}

/* =========================================
   Markdown 渲染样式定制 (模拟 GitHub 风格)
   使用 :deep() 穿透 v-html 生成的内容
   ========================================= */
.markdown-body {
  color: #24292e;
  font-family: -apple-system,BlinkMacSystemFont,Segoe UI,Helvetica,Arial,sans-serif;
  font-size: 16px;
  line-height: 1.8;
}

.markdown-body :deep(h1),
.markdown-body :deep(h2),
.markdown-body :deep(h3) {
  margin-top: 24px;
  margin-bottom: 16px;
  font-weight: 600;
  line-height: 1.25;
  border-bottom: 1px solid #eaecef;
  padding-bottom: 0.3em;
}

.markdown-body :deep(p) {
  margin-top: 0;
  margin-bottom: 16px;
}

.markdown-body :deep(blockquote) {
  padding: 0 1em;
  color: #6a737d;
  border-left: 0.25em solid #dfe2e5;
  margin: 0 0 16px 0;
  background-color: #fafbfc;
}

.markdown-body :deep(ul),
.markdown-body :deep(ol) {
  padding-left: 2em;
  margin-bottom: 16px;
}

.markdown-body :deep(img) {
  max-width: 100%;
  box-sizing: border-box;
  background-color: #fff;
  border-radius: 4px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

/* 代码块样式 */
.markdown-body :deep(pre) {
  padding: 16px;
  overflow: auto;
  font-size: 85%;
  line-height: 1.45;
  background-color: #1e1e1e; /* 与 github-dark 主题匹配的背景色 */
  border-radius: 6px;
  margin-bottom: 16px;
}

.markdown-body :deep(code) {
  padding: 0.2em 0.4em;
  margin: 0;
  font-size: 85%;
  background-color: rgba(27,31,35,0.05);
  border-radius: 3px;
  font-family: SFMono-Regular,Consolas,Liberation Mono,Menlo,monospace;
}

.markdown-body :deep(pre code) {
  background-color: transparent;
  padding: 0;
  color: inherit; /* 让 highlight.js 控制颜色 */
}

/* 链接样式 */
.markdown-body :deep(a) {
  color: #0366d6;
  text-decoration: none;
}
.markdown-body :deep(a:hover) {
  text-decoration: underline;
}
</style>