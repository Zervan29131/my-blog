<script setup lang="ts">
import { nextTick, onBeforeUnmount, ref, watch } from 'vue'
import { RouterLink } from 'vue-router'

import { fetchArticle, fetchComments } from '../../api/blog'
import { apiErrorMessage } from '../../api/http'
import CommentForm from '../../components/CommentForm.vue'
import ErrorState from '../../components/ErrorState.vue'
import LoadingState from '../../components/LoadingState.vue'
import PaginationControls from '../../components/PaginationControls.vue'
import type { Article, Comment } from '../../types/blog'
import { usePublicConfigStore } from '../../stores/publicConfig'
import { formatDate } from '../../utils/format'
import { renderArticleMarkdown, type ArticleHeading } from '../../utils/markdown'

const props = defineProps<{
  slug: string
}>()

const configStore = usePublicConfigStore()

const article = ref<Article | null>(null)
const comments = ref<Comment[]>([])
const commentPage = ref(1)
const commentTotal = ref(0)
const commentTotalPages = ref(0)
const loading = ref(true)
const commentsLoading = ref(false)
const errorMessage = ref('')
const commentsError = ref('')
const renderedContent = ref('')
const outline = ref<ArticleHeading[]>([])
const activeHeading = ref('')
let outlineObserver: IntersectionObserver | null = null

function updateDescription(content: string) {
  let description = document.querySelector<HTMLMetaElement>('meta[name="description"]')
  if (!description) {
    description = document.createElement('meta')
    description.name = 'description'
    document.head.append(description)
  }
  description.content = content || `${configStore.site.name}博客文章`
}

function setupOutline() {
  outlineObserver?.disconnect()
  const headings = document.querySelectorAll<HTMLElement>('.article-content h2[id], .article-content h3[id]')
  if (!headings.length) return
  activeHeading.value = headings[0].id
  if (typeof IntersectionObserver === 'undefined') return
  outlineObserver = new IntersectionObserver((entries) => {
    const visible = entries.filter((entry) => entry.isIntersecting)
    if (visible.length) activeHeading.value = (visible[0].target as HTMLElement).id
  }, { rootMargin: '-80px 0px -70% 0px', threshold: 0 })
  headings.forEach((heading) => outlineObserver?.observe(heading))
}

function goToHeading(id: string) {
  const heading = document.getElementById(id)
  if (!heading) return
  heading.scrollIntoView({ behavior: 'smooth', block: 'start' })
  window.history.replaceState(null, '', `#${encodeURIComponent(id)}`)
  activeHeading.value = id
}

async function loadComments(targetPage = 1) {
  commentsLoading.value = true
  commentsError.value = ''
  try {
    const result = await fetchComments(props.slug, targetPage, 20)
    comments.value = result.items
    commentPage.value = result.page
    commentTotal.value = result.total
    commentTotalPages.value = result.total_pages
  } catch (error) {
    commentsError.value = apiErrorMessage(error, '评论加载失败，请稍后重试。')
  } finally {
    commentsLoading.value = false
  }
}

async function loadPage() {
  loading.value = true
  errorMessage.value = ''
  article.value = null
  renderedContent.value = ''
  outline.value = []
  try {
    article.value = await fetchArticle(props.slug)
    const rendered = renderArticleMarkdown(article.value.content)
    renderedContent.value = rendered.html
    outline.value = rendered.headings
    updateDescription(article.value.summary)
    await loadComments(1)
    await nextTick()
    setupOutline()
  } catch (error) {
    errorMessage.value = apiErrorMessage(error, '文章加载失败，请稍后重试。')
  } finally {
    loading.value = false
  }
}

function changeCommentPage(targetPage: number) {
  void loadComments(targetPage)
}

watch(() => props.slug, () => void loadPage(), { immediate: true })
watch(
  [() => article.value?.title, () => errorMessage.value, () => configStore.titleName],
  ([title, error, siteTitle]) => {
    if (title) document.title = `${title} | ${siteTitle}`
    else if (error) document.title = `文章加载失败 | ${siteTitle}`
  },
  { immediate: true },
)
onBeforeUnmount(() => {
  outlineObserver?.disconnect()
})
</script>

<template>
  <main class="page-container detail-page">
    <LoadingState v-if="loading" label="正在打开文章…" />
    <ErrorState v-else-if="errorMessage" :message="errorMessage" @retry="loadPage" />

    <div v-else-if="article" class="article-layout">
      <div class="article-column">
        <article class="article-detail">
          <RouterLink class="back-link" to="/archive">← 返回文章列表</RouterLink>
          <header class="article-header">
            <h1>{{ article.title }}</h1>
            <div class="article-byline">
              <time :datetime="article.published_at">{{ formatDate(article.published_at) }}</time>
              <span aria-hidden="true">·</span>
              <span>{{ commentTotal }} 条评论</span>
              <template v-if="article.updated_at !== article.created_at">
                <span aria-hidden="true">·</span>
                <span>更新于 {{ formatDate(article.updated_at) }}</span>
              </template>
            </div>
            <p v-if="article.summary" class="article-summary">{{ article.summary }}</p>
          </header>

          <div class="article-content markdown-body" v-html="renderedContent"></div>

          <footer class="article-footer">
            <p>最后更新于 {{ formatDate(article.updated_at) }}</p>
            <RouterLink to="/archive">返回文章列表 →</RouterLink>
          </footer>
        </article>

        <section class="comments-section" aria-labelledby="comments-heading">
          <div class="section-heading compact">
            <h2 id="comments-heading">评论 <span>{{ commentTotal }}</span></h2>
          </div>

          <LoadingState v-if="commentsLoading" label="正在加载评论…" />
          <ErrorState
            v-else-if="commentsError"
            :message="commentsError"
            @retry="loadComments(commentPage)"
          />
          <ol v-else-if="comments.length" class="comment-list">
            <li v-for="comment in comments" :key="comment.id" class="comment-item">
              <div class="comment-meta">
                <strong>{{ comment.nickname }}</strong>
                <time :datetime="comment.created_at">{{ formatDate(comment.created_at) }}</time>
              </div>
              <p>{{ comment.content }}</p>
            </li>
          </ol>
          <div v-else class="comments-empty">
            <strong>暂时还没有评论</strong>
            <p>欢迎留下第一条评论。</p>
          </div>

          <PaginationControls :page="commentPage" :total-pages="commentTotalPages" @change="changeCommentPage" />

          <div class="comment-form-section">
            <div>
              <h3>留下评论</h3>
              <p>评论将在审核通过后公开显示。</p>
            </div>
            <CommentForm :slug="article.slug" />
          </div>
        </section>
      </div>

      <aside v-if="outline.length" class="article-outline" aria-label="本页目录">
        <strong>本页目录</strong>
        <nav>
          <a
            v-for="heading in outline"
            :key="heading.id"
            :href="`#${heading.id}`"
            :class="[{ active: activeHeading === heading.id }, `level-${heading.level}`]"
            @click.prevent="goToHeading(heading.id)"
          >{{ heading.text }}</a>
        </nav>
      </aside>
    </div>
  </main>
</template>
