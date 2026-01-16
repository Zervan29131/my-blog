<template>
  <header class="site-header">
    <div class="nav-wrapper">
      <!-- Logo 区域 -->
      <div class="logo" @click="$router.push('/')">
        My Blog
      </div>

      <!-- 菜单链接 -->
      <nav class="nav-links">
        <router-link to="/" active-class="active">首页</router-link>
        
        <!-- 这里的 About 路由如果还没有，可以先保留或注释 -->
        <router-link to="/about" active-class="active">关于</router-link>

        <!-- 状态感知按钮 -->
        <div class="auth-action">
          <template v-if="userStore.isLoggedIn">
            <router-link to="/admin/dashboard" class="dashboard-btn">
              <el-icon><Odometer /></el-icon> 控制台
            </router-link>
          </template>
          
          <template v-else>
            <router-link to="/login" class="login-btn">
              登录
            </router-link>
          </template>
        </div>
      </nav>
    </div>
  </header>
</template>

<script setup lang="ts">
import { useUserStore } from '../stores/user'
import { Odometer } from '@element-plus/icons-vue'

const userStore = useUserStore()
</script>

<style scoped>
.site-header {
  position: sticky;
  top: 0;
  z-index: 100;
  width: 100%;
  height: 64px;
  background: rgba(255, 255, 255, 0.85);
  backdrop-filter: blur(12px); /* 磨砂玻璃效果 */
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.02);
  transition: all 0.3s ease;
}

.nav-wrapper {
  max-width: 1000px;
  height: 100%;
  margin: 0 auto;
  padding: 0 20px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.logo {
  font-size: 1.25rem;
  font-weight: 800;
  color: #2c3e50;
  cursor: pointer;
  letter-spacing: -0.5px;
  transition: color 0.3s;
}

.logo:hover {
  color: #409eff;
}

.nav-links {
  display: flex;
  align-items: center;
  gap: 24px;
}

.nav-links a {
  font-size: 15px;
  font-weight: 500;
  color: #606266;
  text-decoration: none;
  position: relative;
  transition: color 0.3s;
}

.nav-links a:not(.login-btn):not(.dashboard-btn):hover,
.nav-links a.active {
  color: #409eff;
}

/* 简单的下划线动画效果 */
.nav-links a:not(.login-btn):not(.dashboard-btn)::after {
  content: '';
  position: absolute;
  bottom: -4px;
  left: 0;
  width: 0;
  height: 2px;
  background-color: #409eff;
  transition: width 0.3s ease;
  border-radius: 2px;
}

.nav-links a.active::after,
.nav-links a:not(.login-btn):not(.dashboard-btn):hover::after {
  width: 100%;
}

/* 按钮样式 */
.login-btn,
.dashboard-btn {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 6px 16px;
  border-radius: 20px;
  font-size: 13px !important;
  font-weight: 600;
  transition: all 0.3s ease;
}

.login-btn {
  background-color: #f0f2f5;
  color: #606266 !important;
}

.login-btn:hover {
  background-color: #e6e8eb;
  color: #303133 !important;
}

.dashboard-btn {
  background-color: #e6f7ff;
  color: #1890ff !important;
}

.dashboard-btn:hover {
  background-color: #409eff;
  color: white !important;
}
</style>