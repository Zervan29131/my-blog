<script setup lang="ts">
import { computed, reactive, ref } from 'vue'

import { submitComment } from '../api/blog'
import { apiErrorMessage } from '../api/http'

const props = defineProps<{
  slug: string
}>()

const form = reactive({
  nickname: '',
  email: '',
  content: '',
})
const submitting = ref(false)
const errorMessage = ref('')
const successMessage = ref('')
const contentLength = computed(() => Array.from(form.content).length)

function validate(): string {
  const nicknameLength = Array.from(form.nickname.trim()).length
  const commentLength = Array.from(form.content.trim()).length
  if (nicknameLength < 2 || nicknameLength > 50) {
    return '昵称需要填写 2～50 个字符。'
  }
  if (form.email.length > 255) {
    return '邮箱不能超过 255 个字符。'
  }
  if (form.email.trim() && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.email.trim())) {
    return '请输入有效的邮箱地址。'
  }
  if (commentLength < 2 || commentLength > 1000) {
    return '评论内容需要填写 2～1000 个字符。'
  }
  return ''
}

async function handleSubmit() {
  errorMessage.value = validate()
  successMessage.value = ''
  if (errorMessage.value) {
    return
  }

  submitting.value = true
  try {
    const response = await submitComment(props.slug, {
      nickname: form.nickname.trim(),
      email: form.email.trim() || undefined,
      content: form.content.trim(),
    })
    successMessage.value = response.message || '评论已提交，审核通过后将会显示。'
    form.nickname = ''
    form.email = ''
    form.content = ''
  } catch (error) {
    errorMessage.value = apiErrorMessage(error, '评论提交失败，请稍后重试。')
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <form class="comment-form" novalidate @submit.prevent="handleSubmit">
    <div class="form-row">
      <label>
        <span>昵称 <em>必填</em></span>
        <input
          v-model="form.nickname"
          name="nickname"
          type="text"
          autocomplete="nickname"
          minlength="2"
          maxlength="50"
          placeholder="怎么称呼你"
          :disabled="submitting"
          required
        />
      </label>
      <label>
        <span>邮箱 <small>选填，不会公开</small></span>
        <input
          v-model="form.email"
          name="email"
          type="email"
          autocomplete="email"
          maxlength="255"
          placeholder="you@example.com"
          :disabled="submitting"
        />
      </label>
    </div>

    <label>
      <span>评论内容 <em>必填</em></span>
      <textarea
        v-model="form.content"
        name="content"
        rows="5"
        minlength="2"
        maxlength="1000"
        placeholder="写下你的想法…"
        :disabled="submitting"
        required
      ></textarea>
    </label>
    <div class="form-footer">
      <span class="character-count" :class="{ warning: contentLength > 900 }">
        {{ contentLength }} / 1000
      </span>
      <button class="button button-primary" type="submit" :disabled="submitting">
        {{ submitting ? '正在提交…' : '提交评论' }}
      </button>
    </div>
    <p v-if="errorMessage" class="form-message error-message" role="alert">{{ errorMessage }}</p>
    <p v-if="successMessage" class="form-message success-message" role="status">
      {{ successMessage }}
    </p>
  </form>
</template>
