<template>
  <div class="admin-layout">
    <!-- 左侧侧边栏 -->
    <aside class="sidebar">
      <div class="logo">Admin Console</div>
      <el-menu router :default-active="$route.path" background-color="#304156" text-color="#fff">
        <el-menu-item index="/admin/dashboard">
          <el-icon><Odometer /></el-icon> 仪表盘
        </el-menu-item>
        <el-menu-item index="/admin/posts">
          <el-icon><Document /></el-icon> 文章管理
        </el-menu-item>
        <el-menu-item index="/admin/categories">
          <el-icon><Collection /></el-icon> 分类标签
        </el-menu-item>
      </el-menu>
    </aside>

    <!-- 右侧主体 -->
    <div class="main-container">
      <header class="admin-header">
        <span>管理员</span>
        <el-button type="text" @click="logout">退出</el-button>
      </header>
      <main class="admin-content">
        <router-view />
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'

const router = useRouter()
const userStore = useUserStore()

const logout = () => {
  userStore.logout()
  router.push('/login')
}
</script>

<style scoped>
.admin-layout { display: flex; height: 100vh; }
.sidebar { width: 220px; background: #304156; color: white; }
.main-container { flex: 1; display: flex; flex-direction: column; }
.admin-header { height: 50px; border-bottom: 1px solid #ddd; display: flex; justify-content: flex-end; align-items: center; padding: 0 20px; }
.admin-content { padding: 20px; flex: 1; overflow-y: auto; background: #f0f2f5; }
</style>