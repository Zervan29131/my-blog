<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { useRouter } from 'vue-router'

import { createArticle, fetchAdminArticle, updateArticle } from '../../api/admin'
import { apiErrorMessage } from '../../api/http'
import type { ArticleInput, ArticleStatus } from '../../types/admin'
import { renderMarkdown } from '../../utils/markdown'

const props = defineProps<{
  articleId?: number
}>()

const router = useRouter()
const form = reactive<ArticleInput>({
  title: '',
  slug: '',
  summary: '',
  content: '',
  status: 'draft',
})
const loading = ref(false)
const saving = ref(false)
const errorMessage = ref('')
const isEditing = computed(() => Number.isInteger(props.articleId) && Number(props.articleId) > 0)
const preview = computed(() => renderMarkdown(form.content || '*正文预览会显示在这里。*'))

function validate(): string {
  if (!form.title.trim() || Array.from(form.title.trim()).length > 200) return '请填写 1～200 个字符的标题。'
  if (form.slug && !/^[a-z0-9-]+$/.test(form.slug)) return 'Slug 只能包含小写字母、数字和连字符。'
  if (Array.from(form.slug).length > 200) return 'Slug 不能超过 200 个字符。'
  if (Array.from(form.summary).length > 500) return '摘要不能超过 500 个字符。'
  if (!form.content.trim()) return '请填写 Markdown 正文。'
  return ''
}

async function loadArticle() {
  if (!isEditing.value || !props.articleId) return
  loading.value = true
  errorMessage.value = ''
  try {
    const article = await fetchAdminArticle(props.articleId)
    Object.assign(form, {
      title: article.title,
      slug: article.slug,
      summary: article.summary,
      content: article.content,
      status: article.status,
    })
  } catch (error) {
    errorMessage.value = apiErrorMessage(error, '文章加载失败，请稍后重试。')
  } finally {
    loading.value = false
  }
}

async function save(status?: ArticleStatus) {
  if (status) form.status = status
  errorMessage.value = validate()
  if (errorMessage.value) return

  saving.value = true
  try {
    const input: ArticleInput = {
      title: form.title.trim(),
      slug: form.slug.trim(),
      summary: form.summary.trim(),
      content: form.content,
      status: form.status,
    }
    if (isEditing.value && props.articleId) {
      await updateArticle(props.articleId, input)
      ElMessage.success('文章已保存')
    } else {
      await createArticle(input)
      ElMessage.success(form.status === 'published' ? '文章已发布' : '草稿已保存')
    }
    await router.push('/admin/articles')
  } catch (error) {
    errorMessage.value = apiErrorMessage(error, '文章保存失败，请稍后重试。')
  } finally {
    saving.value = false
  }
}

onMounted(() => void loadArticle())
</script>

<template>
  <section>
    <div class="admin-page-heading editor-heading">
      <div>
        <p>{{ isEditing ? 'EDIT ARTICLE' : 'NEW ARTICLE' }}</p>
        <h1>{{ isEditing ? '编辑文章' : '新建文章' }}</h1>
        <span>使用 Markdown 编写正文，并在右侧确认最终效果。</span>
      </div>
      <el-button @click="router.push('/admin/articles')">返回列表</el-button>
    </div>

    <el-alert v-if="errorMessage" :title="errorMessage" type="error" show-icon :closable="false" />
    <div v-loading="loading" class="article-editor-grid">
      <el-form class="editor-form" label-position="top">
        <el-form-item label="标题" required>
          <el-input v-model="form.title" maxlength="200" show-word-limit placeholder="文章标题" />
        </el-form-item>
        <el-form-item>
          <template #label>Slug <span class="field-hint">创建时留空可由标题自动生成</span></template>
          <el-input v-model="form.slug" maxlength="200" placeholder="article-slug" />
        </el-form-item>
        <el-form-item label="摘要">
          <el-input v-model="form.summary" type="textarea" :rows="3" maxlength="500" show-word-limit placeholder="文章列表中显示的简短摘要" />
        </el-form-item>
        <el-form-item label="Markdown 正文" required>
          <el-input v-model="form.content" class="markdown-editor" type="textarea" :rows="20" placeholder="# 从这里开始写作" />
        </el-form-item>
        <el-form-item label="文章状态">
          <el-radio-group v-model="form.status">
            <el-radio-button value="draft">草稿</el-radio-button>
            <el-radio-button value="published">已发布</el-radio-button>
          </el-radio-group>
        </el-form-item>
        <div class="editor-actions">
          <el-button :disabled="saving" @click="save('draft')">保存为草稿</el-button>
          <el-button type="primary" :loading="saving" @click="save(form.status)">
            {{ form.status === 'published' ? '保存并发布' : '保存文章' }}
          </el-button>
        </div>
      </el-form>

      <aside class="editor-preview">
        <div class="preview-heading">
          <span>实时预览</span>
          <small>已安全过滤 HTML</small>
        </div>
        <article class="markdown-body" v-html="preview"></article>
      </aside>
    </div>
  </section>
</template>
