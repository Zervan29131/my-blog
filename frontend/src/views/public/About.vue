<template>
  <div class="about-container">
    <main class="main-wrapper">
      <el-row :gutter="20">
        <!-- 左侧：个人信息卡片 -->
        <el-col :span="8" :xs="24">
          <el-card class="profile-card" shadow="hover">
            <div class="avatar-wrapper">
              <el-avatar :size="100" src="https://i.pravatar.cc/300?u=zervan" class="avatar" />
            </div>
            <h2 class="name">Zervan</h2>
            <p class="role">全栈开发者 / 技术博主</p>
            
            <div class="social-links">
              <el-button circle size="large" class="social-btn">
                <el-icon><Platform /></el-icon>
              </el-button>
              <el-button circle size="large" class="social-btn">
                <el-icon><Message /></el-icon>
              </el-button>
            </div>

            <el-divider />

            <div class="skills-section">
              <h3 class="section-title">技能栈</h3>
              <div class="tags-wrapper">
                <el-tag effect="plain" class="skill-tag">Vue 3</el-tag>
                <el-tag effect="plain" class="skill-tag">TypeScript</el-tag>
                <el-tag effect="plain" class="skill-tag">Golang</el-tag>
                <el-tag effect="plain" class="skill-tag">MySQL</el-tag>
                <el-tag effect="plain" class="skill-tag">Docker</el-tag>
              </div>
            </div>
          </el-card>
        </el-col>

        <!-- 右侧：详细介绍 -->
        <el-col :span="16" :xs="24">
          <el-card class="content-card" shadow="hover">
            <article class="markdown-body">
              <h1>关于我</h1>
              <p>你好！我是 Zervan，一名热爱技术与分享的开发者。</p>
              
              <h2>✨ 博客初衷</h2>
              <p>
                创建这个博客是为了记录我在学习过程中的点滴，同时也希望能够帮助到遇到同样问题的朋友。
                在这里，你将看到关于 <strong>前端开发</strong>、<strong>后端架构</strong> 以及 <strong>DevOps</strong> 的相关文章。
              </p>

              <h2>🚀 我的经历</h2>
              <ul>
                <li>2020 - 至今：专注于 Web 全栈开发</li>
                <li>热衷于开源社区，喜欢折腾各种新技术</li>
              </ul>

              <h2>📫 联系我</h2>
              <p>如果你有任何问题或建议，欢迎通过邮件联系我，或者在文章下方留言。</p>
            </article>
          </el-card>
        </el-col>
      </el-row>
    </main>
  </div>
</template>

<script setup lang="ts">
import { Platform, Message } from '@element-plus/icons-vue'
// 无需额外的逻辑，纯展示页面
</script>

<style scoped>
.about-container {
  padding-top: 100px; /* 留出导航栏高度 */
  padding-bottom: 40px;
  min-height: 100vh;
  background-color: var(--bg-color);
  transition: background-color 0.3s;
}

.main-wrapper {
  max-width: 1000px;
  margin: 0 auto;
  padding: 0 20px;
}

/* 卡片通用样式 */
.profile-card, .content-card {
  background-color: var(--bg-content);
  color: var(--text-main);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  transition: all 0.3s;
  margin-bottom: 20px;
}

/* 🟢 深色模式：卡片样式增强 */
:global(html.dark) .profile-card, 
:global(html.dark) .content-card {
  background-color: #1d1e1f;
  border-color: #363637;
  /* 移除白色背景下的默认阴影，改为微弱发光或纯边框 */
  box-shadow: none;
}

:global(html.dark) .profile-card:hover, 
:global(html.dark) .content-card:hover {
  border-color: var(--primary-color);
}

/* 个人信息 */
.avatar-wrapper {
  display: flex;
  justify-content: center;
  margin-bottom: 20px;
}
.name {
  text-align: center;
  font-size: 1.8rem;
  margin: 0 0 10px;
  color: var(--text-main);
}
.role {
  text-align: center;
  color: var(--text-secondary);
  margin-bottom: 20px;
}

.social-links {
  display: flex;
  justify-content: center;
  gap: 15px;
  margin-bottom: 20px;
}
.social-btn {
  font-size: 1.2rem;
  transition: all 0.3s;
}
:global(html.dark) .social-btn {
  background-color: transparent;
  border-color: var(--border-color);
  color: var(--text-main);
}
:global(html.dark) .social-btn:hover {
  color: var(--primary-color);
  border-color: var(--primary-color);
}

/* 技能标签 */
.section-title {
  font-size: 1rem;
  margin-bottom: 15px;
  border-left: 4px solid var(--primary-color);
  padding-left: 10px;
  color: var(--text-main);
}
.tags-wrapper {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}
.skill-tag {
  border-color: var(--border-color);
  color: var(--text-regular);
  background-color: var(--bg-color);
}
:global(html.dark) .skill-tag {
  background-color: #2c2c2c;
  color: #ccc;
}

/* 文章内容区域 */
.markdown-body {
  color: var(--text-main);
  line-height: 1.8;
}

.markdown-body h1, 
.markdown-body h2 {
  color: var(--text-main);
  border-bottom: 1px solid var(--border-color);
  padding-bottom: 0.3em;
}

.markdown-body strong {
  color: var(--primary-color);
}

.markdown-body ul {
  padding-left: 20px;
  color: var(--text-regular);
}

/* 响应式 */
@media (max-width: 768px) {
  .about-container {
    padding-top: 80px;
  }
}
</style>