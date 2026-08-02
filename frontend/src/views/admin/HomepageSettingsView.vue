<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'

import {
  addFeaturedArticle,
  createSocialLink,
  deleteFeaturedArticle,
  deleteSocialLink,
  fetchAdminArticles,
  fetchFeaturedArticles,
  fetchHomepageDraft,
  fetchHomepagePublished,
  fetchSocialLinks,
  publishHomepage,
  reorderFeaturedArticles,
  reorderSocialLinks,
  resetHomepageDraft,
  saveHomepageDraft,
  updateFeaturedArticleVisibility,
  updateSocialLink,
} from '../../api/admin'
import { apiErrorMessage } from '../../api/http'
import type {
  AdminArticleSummary,
  AdminFeaturedArticle,
  AdminHomepageConfig,
  AdminSocialLink,
  HomepageButton,
  HomepageModule,
  HomepageModuleType,
  SocialLinkInput,
  SocialPlatform,
  TechItem,
} from '../../types/admin'

const moduleMeta: Array<{ type: HomepageModuleType; label: string; code: string; hint: string }> = [
  { type: 'hero', label: 'Hero 欢迎区域', code: '01', hint: '首页标题、介绍和操作按钮' },
  { type: 'about', label: '个人简介', code: '02', hint: 'Markdown 简介与个人图片' },
  { type: 'featured_articles', label: '推荐文章', code: '03', hint: '精选内容与推荐关系' },
  { type: 'latest_articles', label: '最新文章', code: '04', hint: '文章数量和展示选项' },
  { type: 'tech_stack', label: '技术栈', code: '05', hint: '技能、工具与关注方向' },
  { type: 'social_links', label: '社交链接', code: '06', hint: '社交模块与链接内容' },
]

const socialPlatforms: Array<{ value: SocialPlatform; label: string }> = [
  { value: 'github', label: 'GitHub' },
  { value: 'email', label: '邮箱' },
  { value: 'linkedin', label: 'LinkedIn' },
  { value: 'x', label: 'X' },
  { value: 'weibo', label: '微博' },
  { value: 'bilibili', label: '哔哩哔哩' },
  { value: 'zhihu', label: '知乎' },
  { value: 'custom', label: '自定义' },
]

const draft = ref<AdminHomepageConfig | null>(null)
const published = ref<AdminHomepageConfig | null>(null)
const savedModules = ref<HomepageModule[]>([])
const selectedType = ref<HomepageModuleType>('hero')
const loading = ref(true)
const saving = ref(false)
const publishing = ref(false)
const resetting = ref(false)
const loadError = ref('')
const validationError = ref('')

const featuredArticles = ref<AdminFeaturedArticle[]>([])
const articleCandidates = ref<AdminArticleSummary[]>([])
const selectedArticleID = ref<number | null>(null)
const featuredWorking = ref(false)

const socialLinks = ref<AdminSocialLink[]>([])
const socialWorking = ref(false)
const emptySocialLink = (): SocialLinkInput => ({
  platform: 'github',
  display_name: '',
  url: '',
  is_visible: true,
  sort_order: (socialLinks.value.length + 1) * 10,
})
const newSocialLink = ref<SocialLinkInput>(emptySocialLink())

function cloneModules(modules: HomepageModule[]): HomepageModule[] {
  return JSON.parse(JSON.stringify(modules)) as HomepageModule[]
}

function moduleOf<T extends HomepageModuleType>(type: T): HomepageModule<T> {
  const module = draft.value?.modules.find((item) => item.type === type)
  if (!module) throw new Error(`Missing homepage module: ${type}`)
  return module as HomepageModule<T>
}

const heroModule = computed(() => moduleOf('hero'))
const aboutModule = computed(() => moduleOf('about'))
const featuredModule = computed(() => moduleOf('featured_articles'))
const latestModule = computed(() => moduleOf('latest_articles'))
const techModule = computed(() => moduleOf('tech_stack'))
const socialModule = computed(() => moduleOf('social_links'))
const orderedModules = computed(() => [...(draft.value?.modules ?? [])].sort((left, right) => left.sort_order - right.sort_order))
const isDirty = computed(() => Boolean(draft.value) && JSON.stringify(draft.value?.modules) !== JSON.stringify(savedModules.value))
const availableArticles = computed(() => {
  const selected = new Set(featuredArticles.value.map((item) => item.article_id))
  return articleCandidates.value.filter((article) => !selected.has(article.id))
})

function formatDateTime(value: string | null | undefined): string {
  if (!value) return '尚未发布'
  const date = new Date(value)
  return Number.isNaN(date.getTime()) ? '时间未知' : new Intl.DateTimeFormat('zh-CN', {
    month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit', hour12: false,
  }).format(date)
}

function applyDraft(config: AdminHomepageConfig) {
  draft.value = { ...config, modules: cloneModules(config.modules) }
  savedModules.value = cloneModules(config.modules)
  validationError.value = ''
}

async function loadPage() {
  loading.value = true
  loadError.value = ''
  try {
    const [draftConfig, publishedConfig, featured, candidates, links] = await Promise.all([
      fetchHomepageDraft(),
      fetchHomepagePublished(),
      fetchFeaturedArticles(),
      fetchAdminArticles(1, 100, 'published'),
      fetchSocialLinks(),
    ])
    applyDraft(draftConfig)
    published.value = publishedConfig
    featuredArticles.value = featured
    articleCandidates.value = candidates.items
    socialLinks.value = links
    newSocialLink.value = emptySocialLink()
  } catch (error) {
    loadError.value = apiErrorMessage(error, '首页配置加载失败，请稍后重试。')
  } finally {
    loading.value = false
  }
}

function characterCount(value: string): number {
  return Array.from(value).length
}

function requiredText(label: string, value: string, maximum: number): string {
  const trimmed = value.trim()
  if (!trimmed || characterCount(trimmed) > maximum) return `${label}需要填写 1～${maximum} 个字符。`
  return ''
}

function optionalText(label: string, value: string, maximum: number): string {
  return characterCount(value.trim()) > maximum ? `${label}不能超过 ${maximum} 个字符。` : ''
}

function isHTTPURL(value: string): boolean {
  try {
    const parsed = new URL(value)
    return parsed.protocol === 'http:' || parsed.protocol === 'https:'
  } catch {
    return false
  }
}

function optionalHTTPURL(label: string, value: string): string {
  const trimmed = value.trim()
  if (!trimmed) return ''
  return characterCount(trimmed) <= 500 && isHTTPURL(trimmed) ? '' : `${label}必须是有效的 HTTP 或 HTTPS 地址。`
}

function safeButtonURL(label: string, button: HomepageButton): string {
  if (!button.enabled) return ''
  const value = button.url.trim()
  if (button.link_type === 'internal') {
    return value.startsWith('/') && !value.startsWith('//') && characterCount(value) <= 500
      ? ''
      : `${label}必须是安全的站内路径。`
  }
  return isHTTPURL(value) && characterCount(value) <= 500 ? '' : `${label}必须是有效的 HTTP 或 HTTPS 地址。`
}

function validateButton(label: string, button: HomepageButton): string {
  const textError = button.enabled
    ? requiredText(`${label}文字`, button.text, 20)
    : optionalText(`${label}文字`, button.text, 20)
  return textError || safeButtonURL(`${label}链接`, button)
}

function normalizeModules(): HomepageModule[] {
  const modules = cloneModules(draft.value?.modules ?? [])
  for (const module of modules) {
    const config = module.config as unknown as Record<string, unknown>
    for (const key of ['eyebrow', 'title', 'highlight_text', 'description', 'image_url', 'background_image_url', 'content']) {
      if (typeof config[key] === 'string') config[key] = config[key].trim()
    }
    if (module.type === 'hero') {
      for (const button of [module.config.primary_button, module.config.secondary_button]) {
        button.text = button.text.trim()
        button.url = button.url.trim()
      }
    }
    if (module.type === 'tech_stack') {
      module.config.items.forEach((item) => {
        item.name = item.name.trim()
        item.description = item.description.trim()
        item.icon_url = item.icon_url.trim()
        item.url = item.url.trim()
      })
    }
  }
  return modules
}

function validateModules(): string {
  const hero = heroModule.value.config
  let error = optionalText('Hero 小标题', hero.eyebrow, 50)
    || requiredText('Hero 主标题', hero.title, 100)
    || optionalText('Hero 强调文字', hero.highlight_text, 50)
    || requiredText('Hero 描述', hero.description, 300)
    || optionalHTTPURL('Hero 图片地址', hero.image_url)
    || optionalHTTPURL('Hero 背景图片地址', hero.background_image_url)
    || validateButton('主按钮', hero.primary_button)
    || validateButton('次按钮', hero.secondary_button)
  if (error) return error

  const about = aboutModule.value.config
  error = requiredText('个人简介标题', about.title, 100)
    || optionalText('个人简介说明', about.description, 200)
    || requiredText('个人简介正文', about.content, 2000)
    || optionalHTTPURL('个人简介图片地址', about.image_url)
  if (error) return error

  const featured = featuredModule.value.config
  error = requiredText('推荐文章标题', featured.title, 100)
    || optionalText('推荐文章说明', featured.description, 200)
  if (error) return error
  if (featured.limit < 1 || featured.limit > 10) return '推荐文章展示数量需要在 1～10 之间。'

  const latest = latestModule.value.config
  error = requiredText('最新文章标题', latest.title, 100)
    || optionalText('最新文章说明', latest.description, 200)
  if (error) return error
  if (latest.limit < 3 || latest.limit > 20) return '最新文章展示数量需要在 3～20 之间。'

  const tech = techModule.value.config
  error = requiredText('技术栈标题', tech.title, 100) || optionalText('技术栈说明', tech.description, 200)
  if (error) return error
  if (tech.items.length > 20) return '技术项最多只能添加 20 个。'
  for (const [index, item] of tech.items.entries()) {
    error = requiredText(`第 ${index + 1} 个技术项名称`, item.name, 30)
      || optionalText(`第 ${index + 1} 个技术项描述`, item.description, 100)
      || optionalHTTPURL(`第 ${index + 1} 个技术项图标`, item.icon_url)
      || optionalHTTPURL(`第 ${index + 1} 个技术项链接`, item.url)
    if (error) return error
  }

  const social = socialModule.value.config
  return requiredText('社交链接标题', social.title, 100) || optionalText('社交链接说明', social.description, 200)
}

function selectModule(type: HomepageModuleType) {
  selectedType.value = type
  validationError.value = ''
}

function moveModule(type: HomepageModuleType, direction: -1 | 1) {
  if (!draft.value) return
  const modules = [...orderedModules.value]
  const index = modules.findIndex((module) => module.type === type)
  const target = index + direction
  if (index < 0 || target < 0 || target >= modules.length) return
  ;[modules[index], modules[target]] = [modules[target], modules[index]]
  modules.forEach((module, moduleIndex) => { module.sort_order = (moduleIndex + 1) * 10 })
  draft.value.modules = modules
}

function discardLocalChanges() {
  if (!draft.value || !isDirty.value) return
  draft.value.modules = cloneModules(savedModules.value)
  validationError.value = ''
  ElMessage.info('已放弃尚未保存的本地修改')
}

async function saveDraft() {
  validationError.value = validateModules()
  if (validationError.value) return
  const modules = normalizeModules()
  saving.value = true
  try {
    applyDraft(await saveHomepageDraft(modules))
    ElMessage.success('首页草稿已保存，不会影响当前公开首页')
  } catch (error) {
    validationError.value = apiErrorMessage(error, '首页草稿保存失败，请稍后重试。')
  } finally {
    saving.value = false
  }
}

function openPreview() {
  if (isDirty.value) {
    validationError.value = '请先保存当前修改，再预览首页草稿。'
    return
  }
  validationError.value = ''
  window.open('/preview/home', '_blank', 'noopener,noreferrer')
}

async function publishDraft() {
  if (isDirty.value) {
    validationError.value = '请先保存当前修改，再发布首页配置。'
    return
  }
  try {
    await ElMessageBox.confirm(
      '发布后，新的首页内容将立即对所有访客生效。',
      '确认发布首页配置？',
      { confirmButtonText: '确认发布', cancelButtonText: '取消', type: 'warning' },
    )
  } catch {
    return
  }

  publishing.value = true
  try {
    const result = await publishHomepage()
    if (draft.value) draft.value.version = result.version
    if (published.value) {
      published.value.version = result.version
      published.value.published_at = result.published_at
      published.value.modules = cloneModules(savedModules.value)
    }
    ElMessage.success(`首页配置已发布（版本 ${result.version}）`)
  } catch (error) {
    validationError.value = apiErrorMessage(error, '首页配置发布失败，请稍后重试。')
  } finally {
    publishing.value = false
  }
}

async function restorePublished() {
  try {
    await ElMessageBox.confirm(
      '当前草稿将被已发布版本覆盖，此操作无法撤销。',
      '恢复为当前已发布配置？',
      { confirmButtonText: '确认恢复', cancelButtonText: '取消', type: 'warning' },
    )
  } catch {
    return
  }

  resetting.value = true
  try {
    applyDraft(await resetHomepageDraft())
    ElMessage.success('草稿已恢复为当前已发布配置')
  } catch (error) {
    validationError.value = apiErrorMessage(error, '恢复草稿失败，请稍后重试。')
  } finally {
    resetting.value = false
  }
}

function addTechItem() {
  if (techModule.value.config.items.length >= 20) return
  techModule.value.config.items.push({
    name: '', description: '', icon_url: '', url: '', is_visible: true,
    sort_order: (techModule.value.config.items.length + 1) * 10,
  })
}

function moveTechItem(index: number, direction: -1 | 1) {
  const items = techModule.value.config.items
  const target = index + direction
  if (target < 0 || target >= items.length) return
  ;[items[index], items[target]] = [items[target], items[index]]
  items.forEach((item, itemIndex) => { item.sort_order = (itemIndex + 1) * 10 })
}

function removeTechItem(index: number) {
  techModule.value.config.items.splice(index, 1)
  techModule.value.config.items.forEach((item, itemIndex) => { item.sort_order = (itemIndex + 1) * 10 })
}

async function refreshFeatured() {
  featuredArticles.value = await fetchFeaturedArticles()
}

async function addFeatured() {
  if (!selectedArticleID.value) return
  featuredWorking.value = true
  try {
    await addFeaturedArticle(selectedArticleID.value, (featuredArticles.value.length + 1) * 10)
    selectedArticleID.value = null
    await refreshFeatured()
    ElMessage.success('推荐文章已添加')
  } catch (error) {
    ElMessage.error(apiErrorMessage(error, '推荐文章添加失败。'))
  } finally {
    featuredWorking.value = false
  }
}

async function setFeaturedVisibility(item: AdminFeaturedArticle, visible: boolean) {
  const previous = !visible
  featuredWorking.value = true
  try {
    await updateFeaturedArticleVisibility(item.article_id, visible)
    ElMessage.success(visible ? '推荐文章已显示' : '推荐文章已隐藏')
  } catch (error) {
    item.is_visible = previous
    ElMessage.error(apiErrorMessage(error, '推荐文章状态更新失败。'))
  } finally {
    featuredWorking.value = false
  }
}

async function moveFeatured(index: number, direction: -1 | 1) {
  const target = index + direction
  if (target < 0 || target >= featuredArticles.value.length) return
  const previous = [...featuredArticles.value]
  ;[featuredArticles.value[index], featuredArticles.value[target]] = [featuredArticles.value[target], featuredArticles.value[index]]
  featuredWorking.value = true
  try {
    await reorderFeaturedArticles(featuredArticles.value)
    featuredArticles.value.forEach((item, itemIndex) => { item.sort_order = (itemIndex + 1) * 10 })
  } catch (error) {
    featuredArticles.value = previous
    ElMessage.error(apiErrorMessage(error, '推荐文章排序失败。'))
  } finally {
    featuredWorking.value = false
  }
}

async function removeFeatured(item: AdminFeaturedArticle) {
  try {
    await ElMessageBox.confirm(`仅移除“${item.title}”的推荐关系，不会删除原文章。`, '移除推荐文章？', {
      confirmButtonText: '确认移除', cancelButtonText: '取消', type: 'warning',
    })
  } catch {
    return
  }
  featuredWorking.value = true
  try {
    await deleteFeaturedArticle(item.article_id)
    await refreshFeatured()
    ElMessage.success('推荐关系已移除')
  } catch (error) {
    ElMessage.error(apiErrorMessage(error, '推荐文章移除失败。'))
  } finally {
    featuredWorking.value = false
  }
}

function validateSocialInput(input: SocialLinkInput): string {
  let error = requiredText('显示名称', input.display_name, 30)
  if (error) return error
  const value = input.url.trim()
  if (input.platform === 'email' && value.toLowerCase().startsWith('mailto:')) {
    return /^mailto:[^\s@]+@[^\s@]+\.[^\s@]+$/i.test(value) && value.length <= 500 ? '' : '邮箱链接必须是有效的 mailto: 地址。'
  }
  return isHTTPURL(value) && value.length <= 500 ? '' : '社交链接必须是有效的 HTTP 或 HTTPS 地址。'
}

function normalizeSocialInput(input: SocialLinkInput): SocialLinkInput {
  return { ...input, display_name: input.display_name.trim(), url: input.url.trim() }
}

async function addSocial() {
  const error = validateSocialInput(newSocialLink.value)
  if (error) {
    ElMessage.error(error)
    return
  }
  socialWorking.value = true
  try {
    socialLinks.value.push(await createSocialLink(normalizeSocialInput(newSocialLink.value)))
    newSocialLink.value = emptySocialLink()
    ElMessage.success('社交链接已添加')
  } catch (requestError) {
    ElMessage.error(apiErrorMessage(requestError, '社交链接添加失败。'))
  } finally {
    socialWorking.value = false
  }
}

async function saveSocial(item: AdminSocialLink) {
  const input: SocialLinkInput = {
    platform: item.platform, display_name: item.display_name, url: item.url,
    is_visible: item.is_visible, sort_order: item.sort_order,
  }
  const error = validateSocialInput(input)
  if (error) {
    ElMessage.error(error)
    return
  }
  socialWorking.value = true
  try {
    Object.assign(item, await updateSocialLink(item.id, normalizeSocialInput(input)))
    ElMessage.success('社交链接已保存')
  } catch (requestError) {
    ElMessage.error(apiErrorMessage(requestError, '社交链接保存失败。'))
  } finally {
    socialWorking.value = false
  }
}

async function moveSocial(index: number, direction: -1 | 1) {
  const target = index + direction
  if (target < 0 || target >= socialLinks.value.length) return
  const previous = [...socialLinks.value]
  ;[socialLinks.value[index], socialLinks.value[target]] = [socialLinks.value[target], socialLinks.value[index]]
  socialWorking.value = true
  try {
    await reorderSocialLinks(socialLinks.value)
    socialLinks.value.forEach((item, itemIndex) => { item.sort_order = (itemIndex + 1) * 10 })
  } catch (error) {
    socialLinks.value = previous
    ElMessage.error(apiErrorMessage(error, '社交链接排序失败。'))
  } finally {
    socialWorking.value = false
  }
}

async function removeSocial(item: AdminSocialLink) {
  try {
    await ElMessageBox.confirm(`确认删除“${item.display_name}”？`, '删除社交链接？', {
      confirmButtonText: '确认删除', cancelButtonText: '取消', type: 'warning',
    })
  } catch {
    return
  }
  socialWorking.value = true
  try {
    await deleteSocialLink(item.id)
    socialLinks.value = socialLinks.value.filter((link) => link.id !== item.id)
    ElMessage.success('社交链接已删除')
  } catch (error) {
    ElMessage.error(apiErrorMessage(error, '社交链接删除失败。'))
  } finally {
    socialWorking.value = false
  }
}

onMounted(() => void loadPage())
</script>

<template>
  <section class="homepage-settings-page">
    <div class="admin-page-heading homepage-heading">
      <div>
        <p>HOMEPAGE WORKSPACE</p>
        <h1>首页配置</h1>
        <span>按模块维护首页内容；保存只更新草稿，发布后才会影响访客。</span>
      </div>
      <span v-if="isDirty" class="settings-dirty-badge">有未保存修改</span>
    </div>

    <el-skeleton v-if="loading" :rows="12" animated />
    <el-result v-else-if="loadError" icon="error" title="加载失败" :sub-title="loadError">
      <template #extra><el-button type="primary" @click="loadPage">重新加载</el-button></template>
    </el-result>

    <template v-else-if="draft">
      <div class="homepage-statusbar">
        <div>
          <span>草稿版本 <strong>v{{ draft.version }}</strong></span>
          <small>最后保存 {{ formatDateTime(draft.updated_at) }}</small>
        </div>
        <div>
          <span>线上版本 <strong>v{{ published?.version ?? '—' }}</strong></span>
          <small>最后发布 {{ formatDateTime(published?.published_at) }}</small>
        </div>
        <div class="homepage-top-actions">
          <el-button :loading="resetting" @click="restorePublished">恢复已发布版本</el-button>
          <el-button :disabled="!isDirty" @click="discardLocalChanges">放弃本地修改</el-button>
          <el-button @click="openPreview">预览首页</el-button>
          <el-button type="primary" plain :loading="saving" :disabled="!isDirty" @click="saveDraft">保存草稿</el-button>
          <el-button type="primary" :loading="publishing" @click="publishDraft">发布首页</el-button>
        </div>
      </div>

      <el-alert v-if="validationError" class="homepage-alert" :title="validationError" type="error" show-icon :closable="false" />

      <div class="homepage-editor-layout">
        <aside class="homepage-module-list" aria-label="首页模块列表">
          <header>
            <div>
              <strong>模块顺序</strong>
              <span>访客将按此顺序看到已启用模块</span>
            </div>
            <small>{{ orderedModules.filter((module) => module.enabled).length }}/6 启用</small>
          </header>

          <article
            v-for="(module, index) in orderedModules"
            :key="module.type"
            class="homepage-module-card"
            :class="{ active: selectedType === module.type, disabled: !module.enabled }"
            @click="selectModule(module.type)"
          >
            <span class="module-code">{{ moduleMeta.find((item) => item.type === module.type)?.code }}</span>
            <div class="module-card-copy">
              <strong>{{ moduleMeta.find((item) => item.type === module.type)?.label }}</strong>
              <span>{{ moduleMeta.find((item) => item.type === module.type)?.hint }}</span>
            </div>
            <el-switch v-model="module.enabled" :aria-label="`${moduleMeta.find((item) => item.type === module.type)?.label}开关`" @click.stop />
            <div class="module-order-actions" @click.stop>
              <button type="button" :aria-label="`上移${moduleMeta.find((item) => item.type === module.type)?.label}`" :disabled="index === 0" @click="moveModule(module.type, -1)">↑</button>
              <button type="button" :aria-label="`下移${moduleMeta.find((item) => item.type === module.type)?.label}`" :disabled="index === orderedModules.length - 1" @click="moveModule(module.type, 1)">↓</button>
            </div>
          </article>
        </aside>

        <el-form class="homepage-module-editor" label-position="top" :disabled="saving || publishing" @submit.prevent>
          <header class="module-editor-heading">
            <span>{{ moduleMeta.find((item) => item.type === selectedType)?.code }}</span>
            <div>
              <h2>{{ moduleMeta.find((item) => item.type === selectedType)?.label }}</h2>
              <p>{{ moduleMeta.find((item) => item.type === selectedType)?.hint }}</p>
            </div>
          </header>

          <div v-if="selectedType === 'hero'" class="settings-grid">
            <el-form-item label="小标题"><el-input v-model="heroModule.config.eyebrow" maxlength="50" show-word-limit placeholder="例如：WELCOME" /></el-form-item>
            <el-form-item label="布局模式">
              <el-radio-group v-model="heroModule.config.layout"><el-radio-button value="left">左对齐</el-radio-button><el-radio-button value="center">居中</el-radio-button></el-radio-group>
            </el-form-item>
            <el-form-item class="settings-field-wide" label="主标题" required><el-input v-model="heroModule.config.title" maxlength="100" show-word-limit /></el-form-item>
            <el-form-item class="settings-field-wide" label="强调文字"><el-input v-model="heroModule.config.highlight_text" maxlength="50" show-word-limit /></el-form-item>
            <el-form-item class="settings-field-wide" label="描述" required><el-input v-model="heroModule.config.description" type="textarea" :rows="4" maxlength="300" show-word-limit /></el-form-item>
            <el-form-item label="头像或插图地址"><el-input v-model="heroModule.config.image_url" maxlength="500" placeholder="https://example.com/avatar.jpg" /></el-form-item>
            <el-form-item label="背景图片地址"><el-input v-model="heroModule.config.background_image_url" maxlength="500" placeholder="https://example.com/hero.jpg" /></el-form-item>

            <section v-for="(button, buttonIndex) in [heroModule.config.primary_button, heroModule.config.secondary_button]" :key="buttonIndex" class="homepage-subpanel settings-field-wide">
              <div class="homepage-subpanel-title"><strong>{{ buttonIndex === 0 ? '主按钮' : '次按钮' }}</strong><el-switch v-model="button.enabled" /></div>
              <div class="settings-grid">
                <el-form-item label="按钮文字" :required="button.enabled"><el-input v-model="button.text" maxlength="20" show-word-limit /></el-form-item>
                <el-form-item label="链接类型"><el-select v-model="button.link_type"><el-option label="站内链接" value="internal" /><el-option label="外部链接" value="external" /></el-select></el-form-item>
                <el-form-item class="settings-field-wide" label="链接地址" :required="button.enabled"><el-input v-model="button.url" maxlength="500" :placeholder="button.link_type === 'internal' ? '/archive' : 'https://example.com'" /></el-form-item>
                <el-form-item class="settings-field-wide" label="打开方式"><div class="settings-switch-row"><div><strong>新窗口打开</strong><span>外部链接建议在新窗口打开。</span></div><el-switch v-model="button.open_in_new_tab" /></div></el-form-item>
              </div>
            </section>
          </div>

          <div v-else-if="selectedType === 'about'" class="settings-grid">
            <el-form-item label="模块标题" required><el-input v-model="aboutModule.config.title" maxlength="100" show-word-limit /></el-form-item>
            <el-form-item label="图片位置"><el-select v-model="aboutModule.config.image_position"><el-option label="图片在左" value="left" /><el-option label="图片在右" value="right" /><el-option label="不显示图片" value="none" /></el-select></el-form-item>
            <el-form-item class="settings-field-wide" label="模块说明"><el-input v-model="aboutModule.config.description" maxlength="200" show-word-limit /></el-form-item>
            <el-form-item class="settings-field-wide" label="简介正文（Markdown）" required><el-input v-model="aboutModule.config.content" class="markdown-editor" type="textarea" :rows="12" maxlength="2000" show-word-limit /></el-form-item>
            <el-form-item v-if="aboutModule.config.image_position !== 'none'" class="settings-field-wide" label="个人图片地址"><el-input v-model="aboutModule.config.image_url" maxlength="500" placeholder="https://example.com/profile.jpg" /></el-form-item>
          </div>

          <div v-else-if="selectedType === 'featured_articles'" class="settings-grid">
            <el-form-item label="模块标题" required><el-input v-model="featuredModule.config.title" maxlength="100" show-word-limit /></el-form-item>
            <el-form-item label="展示数量" required><el-input-number v-model="featuredModule.config.limit" :min="1" :max="10" /></el-form-item>
            <el-form-item class="settings-field-wide" label="模块说明"><el-input v-model="featuredModule.config.description" maxlength="200" show-word-limit /></el-form-item>

            <section class="homepage-subpanel settings-field-wide">
              <div class="homepage-subpanel-title">
                <div><strong>推荐文章内容</strong><span>只能添加已发布文章，最多 10 篇。</span></div>
                <el-tag type="info" effect="plain">已选 {{ featuredArticles.length }}/10</el-tag>
              </div>
              <div class="featured-add-row">
                <el-select v-model="selectedArticleID" filterable clearable placeholder="搜索并选择已发布文章" :disabled="featuredArticles.length >= 10">
                  <el-option v-for="article in availableArticles" :key="article.id" :label="article.title" :value="article.id" />
                </el-select>
                <el-button type="primary" plain :disabled="!selectedArticleID || featuredArticles.length >= 10" :loading="featuredWorking" @click="addFeatured">添加推荐</el-button>
              </div>
              <div v-if="featuredArticles.length" class="managed-list">
                <article v-for="(item, index) in featuredArticles" :key="item.article_id" class="managed-list-row">
                  <span class="managed-order">{{ String(index + 1).padStart(2, '0') }}</span>
                  <div class="managed-copy"><strong>{{ item.title }}</strong><span>/articles/{{ item.slug }}</span></div>
                  <el-switch v-model="item.is_visible" aria-label="推荐文章显示状态" :disabled="featuredWorking" @change="setFeaturedVisibility(item, Boolean($event))" />
                  <div class="managed-actions"><el-button text :disabled="index === 0 || featuredWorking" @click="moveFeatured(index, -1)">上移</el-button><el-button text :disabled="index === featuredArticles.length - 1 || featuredWorking" @click="moveFeatured(index, 1)">下移</el-button><el-button text type="danger" :disabled="featuredWorking" @click="removeFeatured(item)">移除</el-button></div>
                </article>
              </div>
              <p v-else class="managed-empty">尚未选择推荐文章。公开首页会安全地显示为空。</p>
            </section>
          </div>

          <div v-else-if="selectedType === 'latest_articles'" class="settings-grid">
            <el-form-item label="模块标题" required><el-input v-model="latestModule.config.title" maxlength="100" show-word-limit /></el-form-item>
            <el-form-item label="展示数量" required><el-input-number v-model="latestModule.config.limit" :min="3" :max="20" /></el-form-item>
            <el-form-item class="settings-field-wide" label="模块说明"><el-input v-model="latestModule.config.description" maxlength="200" show-word-limit /></el-form-item>
            <el-form-item class="settings-field-wide" label="展示内容">
              <div class="homepage-toggle-grid">
                <label><span><strong>文章摘要</strong><small>显示文章简介</small></span><el-switch v-model="latestModule.config.show_summary" /></label>
                <label><span><strong>发布日期</strong><small>显示发布时间</small></span><el-switch v-model="latestModule.config.show_date" /></label>
                <label><span><strong>评论数量</strong><small>显示已审核评论数</small></span><el-switch v-model="latestModule.config.show_comment_count" /></label>
                <label><span><strong>查看全部</strong><small>显示归档入口</small></span><el-switch v-model="latestModule.config.show_view_all" /></label>
              </div>
            </el-form-item>
          </div>

          <div v-else-if="selectedType === 'tech_stack'" class="settings-grid">
            <el-form-item label="模块标题" required><el-input v-model="techModule.config.title" maxlength="100" show-word-limit /></el-form-item>
            <el-form-item label="模块说明"><el-input v-model="techModule.config.description" maxlength="200" show-word-limit /></el-form-item>
            <section class="homepage-subpanel settings-field-wide">
              <div class="homepage-subpanel-title"><div><strong>技术项</strong><span>最多 20 项，图片只接受 HTTP/HTTPS 地址。</span></div><el-button type="primary" plain :disabled="techModule.config.items.length >= 20" @click="addTechItem">添加技术项</el-button></div>
              <div v-if="techModule.config.items.length" class="tech-item-list">
                <article v-for="(item, index) in techModule.config.items" :key="index" class="tech-item-card">
                  <header><span>{{ String(index + 1).padStart(2, '0') }}</span><el-switch v-model="item.is_visible" /><div><el-button text :disabled="index === 0" @click="moveTechItem(index, -1)">上移</el-button><el-button text :disabled="index === techModule.config.items.length - 1" @click="moveTechItem(index, 1)">下移</el-button><el-button text type="danger" @click="removeTechItem(index)">删除</el-button></div></header>
                  <div class="settings-grid"><el-form-item label="名称" required><el-input v-model="item.name" maxlength="30" show-word-limit /></el-form-item><el-form-item label="简短描述"><el-input v-model="item.description" maxlength="100" show-word-limit /></el-form-item><el-form-item label="图标地址"><el-input v-model="item.icon_url" maxlength="500" placeholder="https://example.com/icon.png" /></el-form-item><el-form-item label="链接地址"><el-input v-model="item.url" maxlength="500" placeholder="https://example.com" /></el-form-item></div>
                </article>
              </div>
              <p v-else class="managed-empty">尚未添加技术项。</p>
            </section>
          </div>

          <div v-else-if="selectedType === 'social_links'" class="settings-grid">
            <el-form-item label="模块标题" required><el-input v-model="socialModule.config.title" maxlength="100" show-word-limit /></el-form-item>
            <el-form-item label="模块说明"><el-input v-model="socialModule.config.description" maxlength="200" show-word-limit /></el-form-item>
            <section class="homepage-subpanel settings-field-wide">
              <div class="homepage-subpanel-title"><div><strong>社交链接内容</strong><span>链接内容独立保存，最多 15 项；邮箱可使用 mailto:。</span></div><el-tag type="info" effect="plain">{{ socialLinks.length }}/15</el-tag></div>
              <div class="social-add-grid">
                <el-select v-model="newSocialLink.platform" aria-label="新链接平台"><el-option v-for="platform in socialPlatforms" :key="platform.value" :label="platform.label" :value="platform.value" /></el-select>
                <el-input v-model="newSocialLink.display_name" maxlength="30" placeholder="显示名称" />
                <el-input v-model="newSocialLink.url" maxlength="500" placeholder="https://example.com 或 mailto:user@example.com" />
                <el-button type="primary" plain :disabled="socialLinks.length >= 15" :loading="socialWorking" @click="addSocial">添加链接</el-button>
              </div>
              <div v-if="socialLinks.length" class="social-link-list">
                <article v-for="(item, index) in socialLinks" :key="item.id" class="social-link-row">
                  <span class="managed-order">{{ String(index + 1).padStart(2, '0') }}</span>
                  <el-select v-model="item.platform" aria-label="链接平台"><el-option v-for="platform in socialPlatforms" :key="platform.value" :label="platform.label" :value="platform.value" /></el-select>
                  <el-input v-model="item.display_name" maxlength="30" aria-label="链接显示名称" />
                  <el-input v-model="item.url" maxlength="500" aria-label="链接地址" />
                  <el-switch v-model="item.is_visible" aria-label="社交链接显示状态" />
                  <div class="managed-actions"><el-button text :disabled="index === 0 || socialWorking" @click="moveSocial(index, -1)">上移</el-button><el-button text :disabled="index === socialLinks.length - 1 || socialWorking" @click="moveSocial(index, 1)">下移</el-button><el-button text :disabled="socialWorking" @click="saveSocial(item)">保存</el-button><el-button text type="danger" :disabled="socialWorking" @click="removeSocial(item)">删除</el-button></div>
                </article>
              </div>
              <p v-else class="managed-empty">尚未添加社交链接。</p>
            </section>
          </div>
        </el-form>
      </div>

      <footer class="settings-actions homepage-footer-actions">
        <span>{{ isDirty ? '草稿有未保存修改；公开首页不会变化' : '草稿已保存；只有发布后才会影响公开首页' }}</span>
        <div><el-button :disabled="!isDirty" @click="discardLocalChanges">放弃修改</el-button><el-button @click="openPreview">预览首页</el-button><el-button type="primary" plain :loading="saving" :disabled="!isDirty" @click="saveDraft">保存草稿</el-button><el-button type="primary" :loading="publishing" @click="publishDraft">发布首页</el-button></div>
      </footer>
    </template>
  </section>
</template>
