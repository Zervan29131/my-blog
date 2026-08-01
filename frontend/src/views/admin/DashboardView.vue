<script setup lang="ts">
import { onMounted, ref } from 'vue'

import { fetchDashboardStats } from '../../api/admin'
import { apiErrorMessage } from '../../api/http'
import type { DashboardStats } from '../../types/admin'

const stats = ref<DashboardStats | null>(null)
const loading = ref(true)
const errorMessage = ref('')

const cards: Array<{ key: keyof DashboardStats; label: string; note: string; tone: string }> = [
  { key: 'article_total', label: '全部文章', note: '包含草稿与已发布内容', tone: 'ink' },
  { key: 'article_published', label: '已发布', note: '读者当前可以访问', tone: 'green' },
  { key: 'article_draft', label: '草稿', note: '仍在编辑中的文章', tone: 'sand' },
  { key: 'comment_pending', label: '待审核评论', note: '需要你的处理', tone: 'orange' },
  { key: 'comment_approved', label: '已通过评论', note: '已在博客公开显示', tone: 'blue' },
]

async function loadStats() {
  loading.value = true
  errorMessage.value = ''
  try {
    stats.value = await fetchDashboardStats()
  } catch (error) {
    errorMessage.value = apiErrorMessage(error, '统计数据加载失败，请稍后重试。')
  } finally {
    loading.value = false
  }
}

onMounted(() => void loadStats())
</script>

<template>
  <section>
    <div class="admin-page-heading">
      <div>
        <p>OVERVIEW</p>
        <h1>内容概览</h1>
        <span>快速了解博客当前的内容状态。</span>
      </div>
      <RouterLink to="/admin/articles/new">
        <el-button type="primary" size="large">写新文章</el-button>
      </RouterLink>
    </div>

    <el-skeleton v-if="loading" :rows="5" animated />
    <el-result v-else-if="errorMessage" icon="error" title="加载失败" :sub-title="errorMessage">
      <template #extra>
        <el-button type="primary" @click="loadStats">重新加载</el-button>
      </template>
    </el-result>
    <div v-else-if="stats" class="stat-grid">
      <article v-for="card in cards" :key="card.key" class="stat-card" :class="`stat-${card.tone}`">
        <span>{{ card.label }}</span>
        <strong>{{ stats[card.key] }}</strong>
        <p>{{ card.note }}</p>
      </article>
    </div>

    <div class="admin-quick-links">
      <RouterLink to="/admin/articles">
        <span>文章管理</span>
        <strong>查看、编辑与发布文章 →</strong>
      </RouterLink>
      <RouterLink to="/admin/comments">
        <span>评论管理</span>
        <strong>审核读者提交的评论 →</strong>
      </RouterLink>
    </div>
  </section>
</template>
