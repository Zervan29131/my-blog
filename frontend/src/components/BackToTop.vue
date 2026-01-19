<template>
  <transition name="fade">
    <div v-show="visible" class="back-to-top" @click="scrollToTop">
      <el-icon size="24"><CaretTop /></el-icon>
    </div>
  </transition>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { CaretTop } from '@element-plus/icons-vue'

const visible = ref(false)

const handleScroll = () => {
  visible.value = window.scrollY > 20
}

const scrollToTop = () => {
  window.scrollTo({
    top: 0,
    behavior: 'smooth'
  })
}

onMounted(() => {
  window.addEventListener('scroll', handleScroll)
})

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll)
})
</script>

<style scoped>
.back-to-top {
  position: fixed;
  right: 40px;
  bottom: 40px;
  width: 50px;
  height: 50px;
  /* 🔴 修改：默认使用内容背景色 (在深色模式下是深灰色) */
  background-color: var(--bg-content);
  color: var(--primary-color);
  border-radius: 50%;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  z-index: 999;
  transition: all 0.3s cubic-bezier(0.25, 0.8, 0.25, 1);
  /* 🔴 新增：边框，确保在深色模式下能看清轮廓 */
  border: 1px solid var(--border-color);
}

.back-to-top:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 16px rgba(64, 158, 255, 0.3);
  background-color: var(--primary-color);
  color: white;
  border-color: var(--primary-color);
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s, transform 0.3s;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: translateY(20px);
}
</style>