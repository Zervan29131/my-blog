<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'

import { fetchSiteSettings, updateSiteSettings } from '../../api/admin'
import { apiErrorMessage } from '../../api/http'
import type { SiteSettings, SiteSettingsInput } from '../../types/admin'

interface SiteSettingsForm {
  site_name: string
  site_short_name: string
  site_description: string
  title_suffix: string
  logo_url: string
  favicon_url: string
  default_share_image_url: string
  copyright_name: string
  start_year: number | null
  additional_text: string
  filing_number: string
  filing_url: string
  show_technology: boolean
  technology_text: string
}

const currentYear = new Date().getFullYear()
const emptyForm = (): SiteSettingsForm => ({
  site_name: '',
  site_short_name: '',
  site_description: '',
  title_suffix: '',
  logo_url: '',
  favicon_url: '',
  default_share_image_url: '',
  copyright_name: '',
  start_year: null,
  additional_text: '',
  filing_number: '',
  filing_url: '',
  show_technology: true,
  technology_text: '',
})

const form = reactive<SiteSettingsForm>(emptyForm())
const savedForm = ref<SiteSettingsForm>(emptyForm())
const loading = ref(true)
const saving = ref(false)
const loadError = ref('')
const validationError = ref('')
const isDirty = computed(() => JSON.stringify(form) !== JSON.stringify(savedForm.value))

function characterCount(value: string): number {
  return Array.from(value).length
}

function isSafeOptionalURL(value: string): boolean {
  if (!value) return true
  try {
    const parsed = new URL(value)
    return parsed.protocol === 'http:' || parsed.protocol === 'https:'
  } catch {
    return false
  }
}

function validate(): string {
  if (!form.site_name.trim() || characterCount(form.site_name.trim()) > 50) return '站点名称需要填写 1～50 个字符。'
  if (characterCount(form.site_short_name.trim()) > 20) return '站点简称不能超过 20 个字符。'
  if (!form.site_description.trim() || characterCount(form.site_description.trim()) > 200) return '站点描述需要填写 1～200 个字符。'
  if (characterCount(form.title_suffix.trim()) > 50) return '标题后缀不能超过 50 个字符。'
  if (!form.copyright_name.trim() || characterCount(form.copyright_name.trim()) > 50) return '版权名称需要填写 1～50 个字符。'
  if (form.start_year !== null && (form.start_year < 1900 || form.start_year > currentYear)) return `起始年份需要在 1900～${currentYear} 之间。`
  if (characterCount(form.additional_text.trim()) > 200) return '附加说明不能超过 200 个字符。'
  if (characterCount(form.filing_number.trim()) > 100) return '备案号不能超过 100 个字符。'
  if (characterCount(form.technology_text.trim()) > 100) return '技术信息不能超过 100 个字符。'

  const urlFields: Array<[string, string]> = [
    ['Logo 地址', form.logo_url.trim()],
    ['Favicon 地址', form.favicon_url.trim()],
    ['默认分享图片地址', form.default_share_image_url.trim()],
    ['备案链接', form.filing_url.trim()],
  ]
  for (const [label, value] of urlFields) {
    if (characterCount(value) > 500 || !isSafeOptionalURL(value)) return `${label}必须是有效的 HTTP 或 HTTPS 地址。`
  }
  return ''
}

function applySettings(settings: SiteSettings) {
  const next: SiteSettingsForm = {
    site_name: settings.site_name,
    site_short_name: settings.site_short_name ?? '',
    site_description: settings.site_description,
    title_suffix: settings.title_suffix ?? '',
    logo_url: settings.logo_url ?? '',
    favicon_url: settings.favicon_url ?? '',
    default_share_image_url: settings.default_share_image_url ?? '',
    copyright_name: settings.copyright_name,
    start_year: settings.start_year,
    additional_text: settings.additional_text ?? '',
    filing_number: settings.filing_number ?? '',
    filing_url: settings.filing_url ?? '',
    show_technology: settings.show_technology,
    technology_text: settings.technology_text ?? '',
  }
  Object.assign(form, next)
  savedForm.value = { ...next }
}

function optional(value: string): string | null {
  const trimmed = value.trim()
  return trimmed || null
}

function toInput(): SiteSettingsInput {
  return {
    site_name: form.site_name.trim(),
    site_short_name: optional(form.site_short_name),
    site_description: form.site_description.trim(),
    title_suffix: optional(form.title_suffix),
    logo_url: optional(form.logo_url),
    favicon_url: optional(form.favicon_url),
    default_share_image_url: optional(form.default_share_image_url),
    copyright_name: form.copyright_name.trim(),
    start_year: form.start_year,
    additional_text: optional(form.additional_text),
    filing_number: optional(form.filing_number),
    filing_url: optional(form.filing_url),
    show_technology: form.show_technology,
    technology_text: optional(form.technology_text),
  }
}

async function loadSettings() {
  loading.value = true
  loadError.value = ''
  try {
    applySettings(await fetchSiteSettings())
  } catch (error) {
    loadError.value = apiErrorMessage(error, '站点设置加载失败，请稍后重试。')
  } finally {
    loading.value = false
  }
}

async function saveSettings() {
  validationError.value = validate()
  if (validationError.value) return

  saving.value = true
  try {
    applySettings(await updateSiteSettings(toInput()))
    ElMessage.success('站点设置已保存并立即生效')
  } catch (error) {
    validationError.value = apiErrorMessage(error, '站点设置保存失败，请稍后重试。')
  } finally {
    saving.value = false
  }
}

function resetSettings() {
  Object.assign(form, savedForm.value)
  validationError.value = ''
  ElMessage.info('已恢复到上次保存的内容')
}

onMounted(() => void loadSettings())
</script>

<template>
  <section class="settings-page">
    <div class="admin-page-heading settings-heading">
      <div>
        <p>SITE SETTINGS</p>
        <h1>站点基础设置</h1>
        <span>集中维护博客名称、品牌资源与全站页脚信息，保存后立即生效。</span>
      </div>
      <span v-if="isDirty" class="settings-dirty-badge">有未保存修改</span>
    </div>

    <el-skeleton v-if="loading" :rows="10" animated />
    <el-result v-else-if="loadError" icon="error" title="加载失败" :sub-title="loadError">
      <template #extra>
        <el-button type="primary" @click="loadSettings">重新加载</el-button>
      </template>
    </el-result>

    <el-form v-else class="settings-form" label-position="top" :disabled="saving" @submit.prevent>
      <el-alert v-if="validationError" :title="validationError" type="error" show-icon :closable="false" />

      <article class="settings-panel">
        <header class="settings-panel-heading">
          <span>01</span>
          <div>
            <h2>基础信息</h2>
            <p>用于浏览器标题、站点介绍和后台品牌识别。</p>
          </div>
        </header>
        <div class="settings-grid">
          <el-form-item label="站点名称" required>
            <el-input v-model="form.site_name" maxlength="50" show-word-limit placeholder="例如：字里行间" />
          </el-form-item>
          <el-form-item label="站点简称">
            <el-input v-model="form.site_short_name" maxlength="20" show-word-limit placeholder="用于空间有限的场景" />
          </el-form-item>
          <el-form-item class="settings-field-wide" label="站点描述" required>
            <el-input v-model="form.site_description" type="textarea" :rows="3" maxlength="200" show-word-limit placeholder="用一句话介绍这个博客" />
          </el-form-item>
          <el-form-item label="浏览器标题后缀">
            <el-input v-model="form.title_suffix" maxlength="50" show-word-limit placeholder="例如：字里行间" />
          </el-form-item>
        </div>
      </article>

      <article class="settings-panel">
        <header class="settings-panel-heading">
          <span>02</span>
          <div>
            <h2>品牌资源</h2>
            <p>首期使用安全的 HTTP/HTTPS 图片地址，不上传或执行 SVG 代码。</p>
          </div>
        </header>
        <div class="settings-grid">
          <el-form-item class="settings-field-wide" label="Logo 地址">
            <el-input v-model="form.logo_url" maxlength="500" placeholder="https://example.com/logo.png" />
          </el-form-item>
          <el-form-item class="settings-field-wide" label="Favicon 地址">
            <el-input v-model="form.favicon_url" maxlength="500" placeholder="https://example.com/favicon.ico" />
          </el-form-item>
          <el-form-item class="settings-field-wide" label="默认分享图片地址">
            <el-input v-model="form.default_share_image_url" maxlength="500" placeholder="https://example.com/share-cover.jpg" />
          </el-form-item>
        </div>
      </article>

      <article class="settings-panel">
        <header class="settings-panel-heading">
          <span>03</span>
          <div>
            <h2>页脚设置</h2>
            <p>配置版权年份、备案信息及技术说明。</p>
          </div>
        </header>
        <div class="settings-grid">
          <el-form-item label="版权名称" required>
            <el-input v-model="form.copyright_name" maxlength="50" show-word-limit placeholder="版权所有者名称" />
          </el-form-item>
          <el-form-item label="起始年份">
            <el-input-number v-model="form.start_year" :min="1900" :max="currentYear" :controls="false" placeholder="例如：2024" />
          </el-form-item>
          <el-form-item class="settings-field-wide" label="附加说明">
            <el-input v-model="form.additional_text" maxlength="200" show-word-limit placeholder="可选的页脚说明" />
          </el-form-item>
          <el-form-item label="备案号">
            <el-input v-model="form.filing_number" maxlength="100" placeholder="例如：京 ICP 备..." />
          </el-form-item>
          <el-form-item label="备案链接">
            <el-input v-model="form.filing_url" maxlength="500" placeholder="https://beian.miit.gov.cn/" />
          </el-form-item>
          <el-form-item class="settings-field-wide" label="技术信息">
            <div class="settings-switch-row">
              <div>
                <strong>显示技术信息</strong>
                <span>关闭后，前台页脚不会显示技术栈说明。</span>
              </div>
              <el-switch v-model="form.show_technology" />
            </div>
          </el-form-item>
          <el-form-item v-if="form.show_technology" class="settings-field-wide" label="技术信息文字">
            <el-input v-model="form.technology_text" maxlength="100" show-word-limit placeholder="Built with Vue 3 and Go" />
          </el-form-item>
        </div>
      </article>

      <footer class="settings-actions">
        <span>{{ isDirty ? '修改尚未保存' : '当前内容已与服务器同步' }}</span>
        <div>
          <el-button :disabled="!isDirty || saving" @click="resetSettings">重置</el-button>
          <el-button type="primary" :loading="saving" :disabled="!isDirty" @click="saveSettings">保存设置</el-button>
        </div>
      </footer>
    </el-form>
  </section>
</template>
