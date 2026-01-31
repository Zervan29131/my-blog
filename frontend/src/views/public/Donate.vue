<template>
  <div class="donate-container">
    <main class="main-wrapper">
      <div class="donate-card-wrapper">
        <el-card class="donate-card" shadow="hover">
          <div class="header">
            <el-icon class="icon-coffee" :size="50"><Coffee /></el-icon>
            <h1>赞赏支持</h1>
            <p class="desc">
              维护博客不易，坚持输出更难。<br>
              如果这里的文章对您有所帮助，<br>
              不妨请作者喝杯咖啡，这将是我持续更新的动力 ☕️
            </p>
          </div>

          <!-- 支付方式切换 -->
          <div class="payment-tabs">
            <div 
              class="tab-item" 
              :class="{ active: activeTab === 'wechat' }"
              @click="activeTab = 'wechat'"
            >
              <el-icon><ChatDotRound /></el-icon> 微信支付
            </div>
            <div 
              class="tab-item" 
              :class="{ active: activeTab === 'alipay' }"
              @click="activeTab = 'alipay'"
            >
              <el-icon><Wallet /></el-icon> 支付宝
            </div>
          </div>

          <!-- 二维码展示区域 -->
          <div class="qr-container">
            <transition name="fade" mode="out-in">
              <div v-if="activeTab === 'wechat'" key="wechat" class="qr-box">
                <!-- 请替换为您的真实微信收款码 -->
                <div class="qr-placeholder wechat-bg">
                  <span class="qr-text">WeChat Pay QR Code</span>
                </div>
                <p class="qr-tip">打开微信 [扫一扫]</p>
              </div>
              
              <div v-else key="alipay" class="qr-box">
                <!-- 请替换为您的真实支付宝收款码 -->
                <div class="qr-placeholder alipay-bg">
                  <span class="qr-text">Alipay QR Code</span>
                </div>
                <p class="qr-tip">打开支付宝 [扫一扫]</p>
              </div>
            </transition>
          </div>

          <el-divider>
            <span class="divider-text">Thank You</span>
          </el-divider>

          <div class="footer-actions">
            <el-button round @click="$router.go(-1)">返回上一页</el-button>
            <el-button type="primary" round @click="$router.push('/')">回到首页</el-button>
          </div>
        </el-card>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { Coffee, ChatDotRound, Wallet } from '@element-plus/icons-vue'

const activeTab = ref<'wechat' | 'alipay'>('wechat')
</script>

<style scoped>
.donate-container {
  padding-top: 80px;
  min-height: 100vh;
  background-color: var(--bg-color);
  transition: background-color 0.3s;
  display: flex;
  align-items: center;
  justify-content: center;
}

.main-wrapper {
  width: 100%;
  max-width: 500px;
  padding: 20px;
}

.donate-card {
  text-align: center;
  border-radius: 16px;
  background: var(--bg-content);
  border: 1px solid var(--border-color);
  overflow: hidden;
}

/* 深色模式适配 */
:global(html.dark) .donate-card {
  background-color: #1d1e1f;
  border-color: #363637;
  color: var(--text-main);
}

.header {
  margin-bottom: 30px;
}
.icon-coffee {
  color: #e6a23c;
  margin-bottom: 15px;
  animation: float 3s ease-in-out infinite;
}
.header h1 {
  font-size: 1.8rem;
  margin: 0 0 15px;
  color: var(--text-main);
}
.desc {
  color: var(--text-secondary);
  line-height: 1.8;
  font-size: 14px;
}

/* 支付切换 Tab */
.payment-tabs {
  display: flex;
  justify-content: center;
  gap: 20px;
  margin-bottom: 30px;
}
.tab-item {
  cursor: pointer;
  padding: 8px 20px;
  border-radius: 20px;
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 14px;
  background: var(--bg-color);
  color: var(--text-regular);
  transition: all 0.3s;
}
.tab-item:hover {
  background: var(--border-color);
}
.tab-item.active {
  background: var(--primary-color);
  color: #fff;
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.3);
}
/* 微信特定颜色 */
.tab-item.active:first-child {
  background: #07c160;
  box-shadow: 0 4px 12px rgba(7, 193, 96, 0.3);
}
/* 支付宝特定颜色 */
.tab-item.active:last-child {
  background: #1677ff;
  box-shadow: 0 4px 12px rgba(22, 119, 255, 0.3);
}

/* 二维码区域 */
.qr-container {
  height: 260px; /* 固定高度防止切换时跳动 */
  display: flex;
  justify-content: center;
  align-items: center;
}
.qr-box {
  display: flex;
  flex-direction: column;
  align-items: center;
}
.qr-placeholder {
  width: 200px;
  height: 200px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 15px;
  color: #fff;
  font-weight: bold;
}
.wechat-bg { background: linear-gradient(135deg, #07c160 0%, #05a350 100%); }
.alipay-bg { background: linear-gradient(135deg, #1677ff 0%, #0e5fd8 100%); }

.qr-tip {
  font-size: 12px;
  color: var(--text-secondary);
  margin: 0;
}

.divider-text {
  font-family: 'Dancing Script', cursive, serif;
  font-size: 1.2rem;
  color: var(--text-secondary);
}

.footer-actions {
  margin-top: 20px;
  display: flex;
  justify-content: center;
  gap: 15px;
}

/* 动画 */
@keyframes float {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-10px); }
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>