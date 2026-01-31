<template>
  <article 
    class="post-item" 
    :class="{ 'reverse': reverse }" 
    @click="handleClick"
  >
    <!-- 文章封面图 -->
    <div class="post-cover">
      <img :src="coverUrl" alt="cover" loading="lazy" />
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
        <span class="meta-icon">
          <el-icon><View /></el-icon> {{ post.view_count || 0 }} 热度
        </span>
        <span class="meta-icon" v-if="post.category">
          <el-icon><Folder /></el-icon> {{ post.category.name }}
        </span>
      </div>

      <p class="post-summary">
        {{ post.summary || '暂无摘要，点击阅读全文...' }}
      </p>
      
      <div class="post-btn-wrapper">
        <el-icon><MoreFilled /></el-icon>
      </div>
    </div>
  </article>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Calendar, View, Folder, MoreFilled } from '@element-plus/icons-vue'
import type { Post } from '../api/post'

const props = defineProps<{
  post: Post
  reverse?: boolean 
}>()

const emit = defineEmits(['click'])

const formatDate = (dateStr: string) => {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleDateString('zh-CN', { 
    year: 'numeric', 
    month: '2-digit', 
    day: '2-digit' 
  })
}

const coverUrl = computed(() => {
  return `https://picsum.photos/seed/${props.post.ID}/400/250`
})

const handleClick = () => {
  emit('click', props.post.ID)
}
</script>

<style scoped>
.post-item {
  display: flex;
  background: var(--bg-content); /* 适配变量 */
  border-radius: 12px;
  box-shadow: var(--shadow-light); /* 适配变量 */
  margin-bottom: 25px;
  overflow: hidden;
  height: 280px;
  cursor: pointer;
  transition: all 0.3s;
  border: 1px solid var(--border-color); /* 适配变量 */
}

/* 🟢 深色模式专属优化：增强边框，减弱阴影 */
:global(html.dark) .post-item {
  border: 1px solid #363637;
  box-shadow: none; /* 深色模式通常不需要重阴影 */
}
:global(html.dark) .post-item:hover {
  background-color: #252627; /* 悬浮时稍微变亮 */
  border-color: var(--primary-color); /* 悬浮时边框变色 */
}

.post-item:hover {
  box-shadow: var(--shadow-hover);
  transform: translateY(-4px);
}

.post-cover {
  width: 45%;
  overflow: hidden;
  position: relative;
}

/* 🟢 深色模式下，图片可以加一层微弱遮罩，防止太刺眼 */
:global(html.dark) .post-cover::after {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0,0,0,0.2);
  pointer-events: none;
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

.post-info {
  width: 55%;
  padding: 30px 40px;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.post-meta-top {
  color: var(--text-secondary);
  font-size: 13px;
  margin-bottom: 10px;
}

.post-title {
  font-size: 1.6rem;
  font-weight: 700;
  color: var(--text-main);
  margin: 0 0 15px;
  line-height: 1.3;
  transition: color 0.3s;
}

.post-item:hover .post-title {
  color: var(--primary-color);
}

.post-meta-mid {
  display: flex;
  gap: 15px;
  font-size: 13px;
  color: var(--text-secondary);
  margin-bottom: 15px;
}

.meta-icon {
  display: flex;
  align-items: center;
  gap: 4px;
}

.post-summary {
  color: var(--text-regular);
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
  color: var(--text-secondary);
}

.post-item.reverse {
  flex-direction: row-reverse;
}

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
  .post-info {
    padding: 20px;
  }
}
</style>