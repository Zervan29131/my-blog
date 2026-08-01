<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useRouter } from 'vue-router'

import { deleteArticle, fetchAdminArticles } from '../../api/admin'
import { apiErrorMessage } from '../../api/http'
import type { AdminArticleSummary } from '../../types/admin'
import { formatDate } from '../../utils/format'

const router = useRouter()
const articles = ref<AdminArticleSummary[]>([])
const page = ref(1)
const total = ref(0)
const loading = ref(true)
const errorMessage = ref('')

async function loadArticles(targetPage = page.value) {
  loading.value = true
  errorMessage.value = ''
  try {
    const result = await fetchAdminArticles(targetPage, 10)
    articles.value = result.items
    page.value = result.page
    total.value = result.total
  } catch (error) {
    errorMessage.value = apiErrorMessage(error, '文章列表加载失败，请稍后重试。')
  } finally {
    loading.value = false
  }
}

async function confirmDelete(article: AdminArticleSummary) {
  try {
    await ElMessageBox.confirm(
      `确定删除“${article.title}”吗？该文章下的评论也会一并删除。`,
      '删除文章',
      { confirmButtonText: '确认删除', cancelButtonText: '取消', type: 'warning' },
    )
    await deleteArticle(article.id)
    ElMessage.success('文章已删除')
    const targetPage = articles.value.length === 1 && page.value > 1 ? page.value - 1 : page.value
    await loadArticles(targetPage)
  } catch (error) {
    if (error !== 'cancel' && error !== 'close') {
      ElMessage.error(apiErrorMessage(error, '删除文章失败，请稍后重试。'))
    }
  }
}

onMounted(() => void loadArticles())
</script>

<template>
  <section>
    <div class="admin-page-heading">
      <div>
        <p>ARTICLES</p>
        <h1>文章管理</h1>
        <span>共 {{ total }} 篇文章，草稿不会出现在博客前台。</span>
      </div>
      <RouterLink to="/admin/articles/new">
        <el-button type="primary" size="large">新建文章</el-button>
      </RouterLink>
    </div>

    <el-alert v-if="errorMessage" :title="errorMessage" type="error" show-icon :closable="false">
      <template #default><el-button link type="primary" @click="loadArticles()">重新加载</el-button></template>
    </el-alert>

    <div class="admin-table-card">
      <el-table v-loading="loading" :data="articles" empty-text="暂无文章">
        <el-table-column label="文章" min-width="280">
          <template #default="{ row }">
            <div class="article-table-title">
              <strong>{{ row.title }}</strong>
              <span>/{{ row.slug }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="110">
          <template #default="{ row }">
            <el-tag :type="row.status === 'published' ? 'success' : 'info'" effect="light">
              {{ row.status === 'published' ? '已发布' : '草稿' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="更新时间" width="150">
          <template #default="{ row }">{{ formatDate(row.updated_at) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="180" align="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="router.push(`/admin/articles/${row.id}/edit`)">编辑</el-button>
            <el-button link type="danger" @click="confirmDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        v-if="total > 10"
        class="admin-pagination"
        background
        layout="prev, pager, next"
        :current-page="page"
        :page-size="10"
        :total="total"
        @current-change="loadArticles"
      />
    </div>
  </section>
</template>
