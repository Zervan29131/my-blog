import { defineStore } from 'pinia'
import { computed, ref } from 'vue'

import { fetchPublicSiteConfig } from '../api/blog'
import type { PublicSiteConfig } from '../types/blog'
import { DEFAULT_PUBLIC_SITE_CONFIG, normalizePublicSiteConfig } from '../utils/publicConfig'

export const usePublicConfigStore = defineStore('public-config', () => {
  const config = ref<PublicSiteConfig>(structuredClone(DEFAULT_PUBLIC_SITE_CONFIG))
  const loading = ref(false)
  const loaded = ref(false)
  const error = ref('')
  let pending: Promise<void> | null = null

  const site = computed(() => config.value.site)
  const navigation = computed(() => config.value.navigation)
  const socialLinks = computed(() => config.value.social_links)
  const footer = computed(() => config.value.footer)
  const titleName = computed(() => site.value.title_suffix || site.value.name)

  async function load(force = false): Promise<void> {
    if (loaded.value && !force) return
    if (pending && !force) return pending
    loading.value = true
    error.value = ''
    pending = fetchPublicSiteConfig()
      .then((result) => {
        config.value = normalizePublicSiteConfig(result)
        loaded.value = true
      })
      .catch(() => {
        error.value = '站点配置暂时无法加载，当前使用安全默认设置。'
        config.value = structuredClone(DEFAULT_PUBLIC_SITE_CONFIG)
      })
      .finally(() => {
        loading.value = false
        pending = null
      })
    return pending
  }

  return { config, site, navigation, socialLinks, footer, titleName, loading, loaded, error, load }
})
