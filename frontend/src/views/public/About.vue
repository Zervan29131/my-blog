<template>
  <div class="about-container">
    <!-- 1. 个人资料头部 -->
    <div class="profile-header">
      <div class="avatar-wrapper">
        <!-- 如果你有特定的头像 URL，请在这里替换 -->
        <el-avatar :size="100" src="https://i.pravatar.cc/300?u=zervan" />
      </div>
      <h1 class="name">Zervan</h1>
      <p class="bio">天下最普通的人之一</p>
      
      <div class="social-actions">
        <el-button round @click="openLink('https://www.zhihu.com/people/zervan')">
          知乎
        </el-button>
        <el-button type="primary" round @click="openLink('mailto:nickmechlin01@outlook.com')">
          联系我
        </el-button>
      </div>
      
      <div class="quote-box">
        <p>“我对普通人类没有兴趣。你们之中如果有外星人、未来人、异世界来客或者超能力者以及喜爱一切有趣事物之人的话就来找我。”</p>
      </div>
    </div>

    <!-- 2. 内容主体 -->
    <div class="content-body">
      <el-row :gutter="40">
        <!-- 左侧：主要内容 -->
        <el-col :span="16" :xs="24">
          
          <!-- 关于本站 -->
          <section class="section-block">
            <h2 class="section-title">
              <el-icon><House /></el-icon> 关于本站
            </h2>
            <div class="text-content">
              <p>欢迎来到我的小站！</p>
              <p>在这里你将看到我的往期文章、最新动态以及其他富有生机的内容。与个人空间不同，本站将发布一些具有技术性、多元化的内容，以别样的方式承载 Zervan 的意志。</p>
              <p>如果您对我的工作感兴趣，或是对本站的运维有更好的建议，可置邮至 nickmechlin01@outlook.com。</p>
              <p>为了获得访问该站点的最佳体验，我们推荐使用 Google Chrome 或 Mozilla Firefox 浏览器。</p>
              <p>除了在本网站外，我还会在知乎更新一些个人内容，两者之间内容不完全重合。另外出于一些保密性原因，笔者部分文章无法展示地非常细致，不过如果你遇到了类似的问题，欢迎通过邮箱和我取得进一步的联系。</p>
            </div>
          </section>

          <!-- 日常 (原 Tag Cloud 改用在这里) -->
          <section class="section-block">
            <h2 class="section-title">
              <el-icon><Coffee /></el-icon> 日常
            </h2>
            <div class="skill-cloud">
              <el-tag class="skill-tag" size="large" effect="plain">读书</el-tag>
              <el-tag class="skill-tag" size="large" effect="plain">码字</el-tag>
              <el-tag class="skill-tag" size="large" effect="plain">追番</el-tag>
              <el-tag class="skill-tag" size="large" effect="plain">睡觉</el-tag>
              <el-tag class="skill-tag" size="large" effect="plain">宅</el-tag>
            </div>
          </section>

          <!-- 技能 -->
          <section class="section-block">
            <h2 class="section-title">
              <el-icon><Trophy /></el-icon> 技能
            </h2>
            <div class="text-content">
              <p>嗯，我什么也不会。</p>
            </div>
          </section>

          <!-- 关于评论 -->
          <section class="section-block">
            <h2 class="section-title">
              <el-icon><ChatDotRound /></el-icon> 关于评论
            </h2>
            <div class="text-content">
              <p>您的意见对我们十分重要。</p>
              <p>通常情况下，您可以可置邮至 nickmechlin01@outlook.com 提交评论，并进行评论审核，审核通过后，该评论会显示出来。</p>
              <p>虽然有些许麻烦，但是这样做可以有效过滤不当评论及广告内容。</p>
              <p class="tips">部分文章的评论系统也许会被上锁，这是因为这些文章已经不具备评论价值或者尚处于极不稳定的测试阶段，评论会导致误解和不必要的维护工作。</p>
            </div>
          </section>

          <!-- 碎碎念 -->
          <section class="section-block">
            <h2 class="section-title">
              <el-icon><Notebook /></el-icon> 碎片思考
            </h2>
            <div class="text-content">
              <p><strong>关于团队文档协作：</strong>我认为，一个比较好的团队文档协作工具，应该满足：在线化、实时协同、注重内容而非形式。</p>
              <p><strong>关于鼓励上升：</strong>在一个层级比较多的团队中，即使鼓励上升，也应优先上升到自己的 leader。跨团队沟通遇到矛盾时，直接拉对方 leader 是职场禁忌之一。</p>
              <p><strong>关于打破边界：</strong>员工打破边界（如客户端去做 C++）对于公司来说并不一定是好事。公司可能要承担效率损失和潜在风险。</p>
              <p><strong>工作原则分享：</strong></p>
              <ul>
                <li><strong>事不关己原则：</strong>除非 100% 把握和你有关，否则不要回复未 @ 你的 OnCall，避免无效卷入。</li>
                <li><strong>肯定别人原则：</strong>在 QA 环节做引申前，先对分享者表示肯定，避免直接否定或秀优越。</li>
                <li><strong>延迟反馈：</strong>当收到用户反馈时，可能已有很多人默默忍受。感谢发声用户的同时，也应审视为何没能更早发现问题。</li>
              </ul>
            </div>
          </section>

          <!-- 许可证 -->
          <section class="section-block">
            <h2 class="section-title">
              <el-icon><Collection /></el-icon> 许可证
            </h2>
            <div class="text-content copyright-box">
              <p><strong>CC BY-SA 4.0（署名-相同方式共享 4.0）</strong> 适用于本站所有文本。</p>
              <p>您可以自由地共享、演绎本作品，惟须遵守：<strong>署名</strong>、<strong>相同方式共享</strong>、<strong>没有附加限制</strong>。</p>
              <p class="footer-note">最后，感谢 thaumy 的文案支持。</p>
            </div>
          </section>

        </el-col>

        <!-- 右侧：侧边栏信息 -->
        <el-col :span="8" :xs="24">
          <el-card class="info-card" shadow="hover">
            <template #header>
              <div class="card-header">
                <span><el-icon><User /></el-icon> 个人情报</span>
              </div>
            </template>
            <div class="info-item">
              <p>目前可以公开的情报：</p>
            </div>
          </el-card>

          <el-card class="info-card mt-20" shadow="hover">
            <template #header>
              <div class="card-header">
                <span><el-icon><Message /></el-icon> 联系方式</span>
              </div>
            </template>
            <ul class="contact-list">
              <li>
                <span class="label">Email:</span>
                <a href="mailto:nickmechlin01@outlook.com" class="link">nickmechlin01@outlook.com</a>
              </li>
              <li>
                <span class="label">QQ:</span>
                <span class="text">1611436675</span>
              </li>
              <li>
                <span class="label">Blog:</span>
                <a href="/" class="link">My Blog</a>
              </li>
            </ul>
          </el-card>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script setup lang="ts">
import { User, Trophy, House, ChatDotRound, Coffee, Notebook, Message, Collection } from '@element-plus/icons-vue'

const openLink = (url: string) => {
  window.open(url, '_blank')
}
</script>

<style scoped>
.about-container {
  min-height: 80vh;
  background: #fff;
}

/* --- 头部 --- */
.profile-header {
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  padding: 60px 20px;
  text-align: center;
  margin-bottom: 40px;
}

.avatar-wrapper {
  margin-bottom: 20px;
}

.name {
  font-size: 2.5rem;
  font-weight: 700;
  color: #2c3e50;
  margin: 0 0 10px;
}

.bio {
  font-size: 1.1rem;
  color: #606266;
  margin-bottom: 20px;
}

.social-actions {
  display: flex;
  justify-content: center;
  gap: 15px;
  margin-bottom: 30px;
}

.quote-box {
  max-width: 600px;
  margin: 0 auto;
  font-style: italic;
  color: #888;
  font-size: 0.9rem;
  line-height: 1.6;
  background: rgba(255,255,255,0.5);
  padding: 15px;
  border-radius: 8px;
}

/* --- 内容区 --- */
.content-body {
  max-width: 1000px;
  margin: 0 auto;
  padding: 0 20px 60px;
}

.section-block {
  margin-bottom: 40px;
}

.section-title {
  font-size: 1.5rem;
  color: #303133;
  margin-bottom: 20px;
  display: flex;
  align-items: center;
  gap: 10px;
  border-bottom: 1px solid #eee;
  padding-bottom: 10px;
}

.text-content p {
  line-height: 1.8;
  color: #555;
  margin-bottom: 15px;
  font-size: 16px;
  text-align: justify;
}

.text-content ul {
  padding-left: 20px;
  margin-bottom: 15px;
  color: #555;
  line-height: 1.8;
}

.text-content li {
  margin-bottom: 8px;
}

.tips {
  font-size: 0.9rem;
  color: #909399;
  background: #f4f4f5;
  padding: 10px;
  border-radius: 4px;
}

.copyright-box {
  background: #f9f9f9;
  padding: 15px;
  border-radius: 4px;
  border-left: 4px solid #dcdfe6;
}

.footer-note {
  margin-top: 20px;
  text-align: right;
  font-size: 0.9rem;
  color: #999;
}

.skill-cloud {
  display: flex;
  flex-wrap: wrap;
  gap: 15px;
}

.skill-tag {
  font-weight: 500;
}

/* --- 侧边栏卡片 --- */
.info-card {
  border-radius: 8px;
}

.mt-20 {
  margin-top: 20px;
}

.card-header {
  display: flex;
  align-items: center;
  font-weight: 600;
  color: #333;
}

.card-header .el-icon {
  margin-right: 6px;
}

.info-item {
  color: #606266;
  font-size: 14px;
  line-height: 1.6;
}

.contact-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.contact-list li {
  margin-bottom: 10px;
  font-size: 14px;
  color: #606266;
  display: flex;
  align-items: flex-start; /* 对齐长邮箱 */
  flex-wrap: wrap;
}

.contact-list .label {
  width: 50px;
  font-weight: 500;
  flex-shrink: 0;
}

.contact-list .link {
  color: #409eff;
  text-decoration: none;
  word-break: break-all; /* 防止长邮箱撑破容器 */
}

.contact-list .link:hover {
  text-decoration: underline;
}

.contact-list .text {
  word-break: break-all;
}
</style>