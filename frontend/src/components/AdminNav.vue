<template>
  <header class="admin-nav">
    <!-- 左侧：面包屑导航 -->
    <div class="left-panel">
      <el-icon class="hamburger" @click="toggleSidebar"><Fold /></el-icon>
      <el-breadcrumb separator="/">
        <el-breadcrumb-item :to="{ path: '/admin/dashboard' }">控制台</el-breadcrumb-item>
        <el-breadcrumb-item>{{ currentRouteName }}</el-breadcrumb-item>
      </el-breadcrumb>
    </div>

    <!-- 右侧：功能区 -->
    <div class="right-panel">
      <el-tooltip content="访问前台首页" placement="bottom">
        <div class="action-item" @click="$router.push('/')">
          <el-icon><Monitor /></el-icon>
        </div>
      </el-tooltip>

      <el-dropdown trigger="click" @command="handleCommand">
        <div class="user-profile">
          <el-avatar :size="30" :src="avatarUrl" />
          <span class="username">{{ userStore.username || 'Administrator' }}</span>
          <el-icon class="el-icon--right"><CaretBottom /></el-icon>
        </div>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item command="home">返回前台</el-dropdown-item>
            <el-dropdown-item command="logout" divided>退出登录</el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
  </header>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'
import { Fold, Monitor, CaretBottom } from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

// 简单的路由名称映射，实际项目可以从 router meta 中读取
const routeNameMap: Record<string, string> = {
  'Dashboard': '仪表盘',
  'AdminPostList': '文章管理',
  'CreatePost': '写文章',
  'EditPost': '编辑文章',
  'CategoryManage': '分类与标签'
}

const currentRouteName = computed(() => {
  return routeNameMap[route.name as string] || '当前页面'
})

const avatarUrl = computed(() => `https://ui-avatars.com/api/?name=${userStore.username}&background=409eff&color=fff`)

const toggleSidebar = () => {
  // 这里可以扩展 Sidebar 的折叠逻辑，暂时仅作为展示
}

const handleCommand = (command: string) => {
  if (command === 'logout') {
    userStore.logout()
    router.push('/login')
  } else if (command === 'home') {
    router.push('/')
  }
}
</script>

<style scoped>
.admin-nav {
  height: 50px;
  background: #fff;
  border-bottom: 1px solid #e6e6e6;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  box-shadow: 0 1px 4px rgba(0,21,41,.08);
}

.left-panel {
  display: flex;
  align-items: center;
  gap: 15px;
}

.hamburger {
  font-size: 20px;
  cursor: pointer;
  color: #606266;
  transition: color 0.3s;
}

.hamburger:hover {
  color: #409eff;
}

.right-panel {
  display: flex;
  align-items: center;
  gap: 20px;
}

.action-item {
  height: 100%;
  display: flex;
  align-items: center;
  padding: 0 10px;
  cursor: pointer;
  color: #606266;
  transition: background .3s;
  border-radius: 4px;
}

.action-item:hover {
  background: #f6f6f6;
}

.user-profile {
  display: flex;
  align-items: center;
  cursor: pointer;
  padding: 5px;
  border-radius: 4px;
  transition: background 0.3s;
}

.user-profile:hover {
  background: #f6f6f6;
}

.username {
  margin: 0 8px;
  font-size: 14px;
  color: #333;
}
</style>