<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'

import { deleteComment, fetchAdminComments, updateCommentStatus } from '../../api/admin'
import { apiErrorMessage } from '../../api/http'
import type { AdminComment, CommentStatus } from '../../types/admin'
import { formatDate } from '../../utils/format'

const comments = ref<AdminComment[]>([])
const statusFilter = ref<CommentStatus | ''>('pending')
const page = ref(1)
const total = ref(0)
const loading = ref(true)
const errorMessage = ref('')

const statusLabels: Record<CommentStatus, string> = {
  pending: '待审核',
  approved: '已通过',
  rejected: '已拒绝',
}

const statusTagTypes: Record<CommentStatus, 'warning' | 'success' | 'info'> = {
  pending: 'warning',
  approved: 'success',
  rejected: 'info',
}

function statusLabel(status: CommentStatus): string {
  return statusLabels[status]
}

function statusTagType(status: CommentStatus): 'warning' | 'success' | 'info' {
  return statusTagTypes[status]
}

async function loadComments(targetPage = page.value) {
  loading.value = true
  errorMessage.value = ''
  try {
    const result = await fetchAdminComments(statusFilter.value, targetPage, 20)
    comments.value = result.items
    page.value = result.page
    total.value = result.total
  } catch (error) {
    errorMessage.value = apiErrorMessage(error, '评论列表加载失败，请稍后重试。')
  } finally {
    loading.value = false
  }
}

function changeFilter() {
  page.value = 1
  void loadComments(1)
}

async function changeStatus(comment: AdminComment, status: CommentStatus) {
  try {
    await updateCommentStatus(comment.id, status)
    ElMessage.success(status === 'approved' ? '评论已通过审核' : status === 'rejected' ? '评论已拒绝' : '评论已恢复待审核')
    await loadComments(page.value)
  } catch (error) {
    ElMessage.error(apiErrorMessage(error, '评论状态更新失败，请稍后重试。'))
  }
}

async function confirmDelete(comment: AdminComment) {
  try {
    await ElMessageBox.confirm(`确定删除 ${comment.nickname} 的这条评论吗？`, '删除评论', {
      confirmButtonText: '确认删除', cancelButtonText: '取消', type: 'warning',
    })
    await deleteComment(comment.id)
    ElMessage.success('评论已删除')
    const targetPage = comments.value.length === 1 && page.value > 1 ? page.value - 1 : page.value
    await loadComments(targetPage)
  } catch (error) {
    if (error !== 'cancel' && error !== 'close') {
      ElMessage.error(apiErrorMessage(error, '删除评论失败，请稍后重试。'))
    }
  }
}

onMounted(() => void loadComments())
</script>

<template>
  <section>
    <div class="admin-page-heading comments-heading">
      <div>
        <p>COMMENTS</p>
        <h1>评论管理</h1>
        <span>审核读者评论，公开接口只会显示已通过内容。</span>
      </div>
      <el-select v-model="statusFilter" class="comment-filter" aria-label="按评论状态筛选" @change="changeFilter">
        <el-option label="全部评论" value="" />
        <el-option label="待审核" value="pending" />
        <el-option label="已通过" value="approved" />
        <el-option label="已拒绝" value="rejected" />
      </el-select>
    </div>

    <el-alert v-if="errorMessage" :title="errorMessage" type="error" show-icon :closable="false">
      <template #default><el-button link type="primary" @click="loadComments()">重新加载</el-button></template>
    </el-alert>

    <div class="admin-table-card comment-table-card">
      <el-table v-loading="loading" :data="comments" empty-text="当前筛选条件下没有评论">
        <el-table-column label="评论内容" min-width="330">
          <template #default="{ row }">
            <div class="comment-table-content">
              <p>{{ row.content }}</p>
              <span>{{ row.nickname }}<template v-if="row.email"> · {{ row.email }}</template></span>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="所属文章" min-width="180">
          <template #default="{ row }">
            <RouterLink class="table-article-link" :to="`/admin/articles/${row.article.id}/edit`">{{ row.article.title }}</RouterLink>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="statusTagType(row.status)">{{ statusLabel(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="提交时间" width="140">
          <template #default="{ row }">{{ formatDate(row.created_at) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="210" align="right">
          <template #default="{ row }">
            <el-button v-if="row.status !== 'approved'" link type="success" @click="changeStatus(row, 'approved')">通过</el-button>
            <el-button v-if="row.status !== 'rejected'" link type="warning" @click="changeStatus(row, 'rejected')">拒绝</el-button>
            <el-button link type="danger" @click="confirmDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        v-if="total > 20"
        class="admin-pagination"
        background
        layout="prev, pager, next"
        :current-page="page"
        :page-size="20"
        :total="total"
        @current-change="loadComments"
      />
    </div>
  </section>
</template>
