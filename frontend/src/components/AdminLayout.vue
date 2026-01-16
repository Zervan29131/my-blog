<template>
  <div class="admin-layout">
    <!-- 左侧侧边栏 -->
    <aside class="sidebar">
      <div class="logo-area">
        <img src="/vite.svg" alt="logo" class="logo-img" />
        <span class="logo-text">Admin</span>
      </div>
      
      <el-menu 
        router 
        :default-active="$route.path" 
        background-color="#304156" 
        text-color="#bfcbd9" 
        active-text-color="#409eff"
        class="el-menu-vertical"
      >
        <el-menu-item index="/admin/dashboard">
          <el-icon><Odometer /></el-icon>
          <span>仪表盘</span>
        </el-menu-item>
        
        <el-menu-item index="/admin/posts">
          <el-icon><Document /></el-icon>
          <span>文章管理</span>
        </el-menu-item>
        
        <el-menu-item index="/admin/categories">
          <el-icon><Collection /></el-icon>
          <span>分类与标签</span>
        </el-menu-item>
      </el-menu>
    </aside>

    <!-- 右侧主体 -->
    <div class="main-container">
      <!-- 🟢 使用新的 AdminNav 组件 -->
      <AdminNav />
      
      <main class="admin-content">
        <!-- 路由出口，添加过场动画 -->
        <router-view v-slot="{ Component }">
          <transition name="fade-transform" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Odometer, Document, Collection } from '@element-plus/icons-vue'
// 引入新组件
import AdminNav from '../components/AdminNav.vue'
</script>

<style scoped>
.admin-layout {
  display: flex;
  height: 100vh;
  width: 100%;
}

.sidebar {
  width: 210px;
  background: #304156;
  color: white;
  display: flex;
  flex-direction: column;
  transition: width 0.3s;
  flex-shrink: 0;
}

.logo-area {
  height: 50px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #2b3649;
  gap: 10px;
}

.logo-img {
  width: 24px;
  height: 24px;
}

.logo-text {
  font-weight: 600;
  font-size: 16px;
}

.el-menu-vertical {
  border-right: none;
}

.main-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden; /* 防止主滚动条双重出现 */
  background-color: #f0f2f5;
}

.admin-content {
  flex: 1;
  padding: 20px;
  overflow-y: auto;
}

/* 简单的页面切换动画 */
.fade-transform-enter-active,
.fade-transform-leave-active {
  transition: all 0.3s;
}

.fade-transform-enter-from {
  opacity: 0;
  transform: translateX(-20px);
}

.fade-transform-leave-to {
  opacity: 0;
  transform: translateX(20px);
}
</style>