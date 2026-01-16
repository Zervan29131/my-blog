<template>
  <div class="author-card">
    <el-avatar :size="60" :src="avatarUrl" class="author-avatar" />
    <div class="author-info">
      <h3 class="author-name">{{ name }}</h3>
      <p class="author-bio">{{ bio || '这个作者很懒，什么都没写。' }}</p>
    </div>
    <el-button type="primary" plain size="small" round>Follow</el-button>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps({
  name: {
    type: String,
    default: 'Admin'
  },
  bio: {
    type: String,
    default: ''
  }
})

// 生成随机头像用于演示
const avatarUrl = computed(() => `https://i.pravatar.cc/150?u=${props.name}`)
</script>

<style scoped>
.author-card {
  background: #f9fafb;
  border-radius: 12px;
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 20px;
  margin-top: 40px;
  border: 1px solid #edf2f7;
}

.author-name {
  margin: 0 0 5px;
  font-size: 18px;
  color: #2c3e50;
}

.author-bio {
  margin: 0;
  font-size: 14px;
  color: #606266;
}

.author-info {
  flex: 1;
}

@media (max-width: 600px) {
  .author-card {
    flex-direction: column;
    text-align: center;
  }
}
</style>