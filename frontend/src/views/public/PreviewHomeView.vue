<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'

import { fetchHomepagePreview } from '../../api/admin'
import { apiErrorMessage } from '../../api/http'
import HomepageRenderer from '../../components/home/HomepageRenderer.vue'
import ErrorState from '../../components/ErrorState.vue'
import LoadingState from '../../components/LoadingState.vue'
import { usePublicConfigStore } from '../../stores/publicConfig'
import type { PublicHomepageConfig } from '../../types/blog'
import { normalizePublicHomepage } from '../../utils/publicConfig'

const configStore = usePublicConfigStore()
const homepage = ref<PublicHomepageConfig>({ version: 0, modules: [] })
const loading = ref(true)
const errorMessage = ref('')
let robotsMeta: HTMLMetaElement | null = null
let previousRobotsContent: string | null = null
let robotsMetaCreated = false

function protectFromIndexing() {
  robotsMeta = document.querySelector<HTMLMetaElement>('meta[name="robots"]')
  if (robotsMeta) {
    previousRobotsContent = robotsMeta.getAttribute('content')
  } else {
    robotsMeta = document.createElement('meta')
    robotsMeta.name = 'robots'
    document.head.append(robotsMeta)
    robotsMetaCreated = true
  }
  robotsMeta.content = 'noindex,nofollow'
}

function restoreRobotsMeta() {
  if (!robotsMeta) return
  if (robotsMetaCreated) robotsMeta.remove()
  else if (previousRobotsContent === null) robotsMeta.removeAttribute('content')
  else robotsMeta.content = previousRobotsContent
  robotsMeta = null
  robotsMetaCreated = false
}

async function loadPreview() {
  loading.value = true
  errorMessage.value = ''
  try {
    homepage.value = normalizePublicHomepage(await fetchHomepagePreview())
  } catch (error) {
    homepage.value = { version: 0, modules: [] }
    errorMessage.value = apiErrorMessage(error, '首页草稿预览加载失败，请确认登录仍然有效。')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  protectFromIndexing()
  void configStore.load()
  void loadPreview()
})
onBeforeUnmount(restoreRobotsMeta)
</script>

<template>
  <main class="page-container home-page dynamic-home-page preview-home-page">
    <aside class="preview-mode-banner" role="status">
      <div>
        <strong>预览模式</strong>
        <span>当前显示已保存的首页草稿，不会影响公开首页。</span>
      </div>
      <RouterLink to="/admin/homepage">返回首页配置</RouterLink>
    </aside>

    <LoadingState v-if="loading" label="正在加载首页草稿预览…" />
    <ErrorState v-else-if="errorMessage" :message="errorMessage" @retry="loadPreview" />
    <HomepageRenderer v-else :homepage="homepage" :site-name="configStore.site.name" />
  </main>
</template>
