<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import { apiErrorMessage } from '../../api/http'
import { useAuthStore } from '../../stores/auth'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const form = reactive({ username: '', password: '' })
const loading = ref(false)
const errorMessage = ref('')

async function submit() {
  errorMessage.value = ''
  if (!form.username.trim() || !form.password) {
    errorMessage.value = '请输入用户名和密码。'
    return
  }

  loading.value = true
  try {
    await authStore.login(form.username.trim(), form.password)
    const redirect = typeof route.query.redirect === 'string' && route.query.redirect.startsWith('/admin')
      ? route.query.redirect
      : '/admin'
    await router.replace(redirect)
  } catch (error) {
    errorMessage.value = apiErrorMessage(error, '登录失败，请稍后重试。')
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <main class="admin-login-page">
    <section class="admin-login-panel">
      <RouterLink class="login-brand" to="/">
        <span>阅</span>
        <strong>字里行间</strong>
      </RouterLink>
      <div class="login-heading">
        <p>ADMINISTRATION</p>
        <h1>欢迎回来</h1>
        <span>登录后管理文章与评论。</span>
      </div>

      <el-form label-position="top" @submit.prevent="submit">
        <el-form-item label="用户名">
          <el-input
            v-model="form.username"
            name="username"
            autocomplete="username"
            maxlength="50"
            placeholder="管理员用户名"
            size="large"
            :disabled="loading"
            @keyup.enter="submit"
          />
        </el-form-item>
        <el-form-item label="密码">
          <el-input
            v-model="form.password"
            name="password"
            type="password"
            autocomplete="current-password"
            maxlength="72"
            placeholder="管理员密码"
            size="large"
            show-password
            :disabled="loading"
            @keyup.enter="submit"
          />
        </el-form-item>
        <p v-if="errorMessage" class="admin-form-error" role="alert">{{ errorMessage }}</p>
        <el-button class="login-button" type="primary" size="large" :loading="loading" @click="submit">
          登录
        </el-button>
      </el-form>
      <RouterLink class="login-back" to="/">← 返回博客首页</RouterLink>
    </section>
    <section class="admin-login-quote" aria-hidden="true">
      <p>“文字是时间留下的纹理。”</p>
      <span>在这里整理、修改，再把它交给读者。</span>
    </section>
  </main>
</template>
