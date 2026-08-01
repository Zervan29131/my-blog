<script setup lang="ts">
const props = defineProps<{
  page: number
  totalPages: number
}>()

const emit = defineEmits<{
  change: [page: number]
}>()

function changePage(page: number) {
  if (page >= 1 && page <= props.totalPages && page !== props.page) {
    emit('change', page)
  }
}
</script>

<template>
  <nav v-if="totalPages > 1" class="pagination" aria-label="分页导航">
    <button
      class="pagination-button"
      type="button"
      :disabled="page <= 1"
      aria-label="上一页"
      @click="changePage(page - 1)"
    >
      ← 上一页
    </button>
    <span aria-current="page">第 {{ page }} / {{ totalPages }} 页</span>
    <button
      class="pagination-button"
      type="button"
      :disabled="page >= totalPages"
      aria-label="下一页"
      @click="changePage(page + 1)"
    >
      下一页 →
    </button>
  </nav>
</template>
