import { createApp } from 'vue'
import {
  ElAlert,
  ElButton,
  ElForm,
  ElFormItem,
  ElInput,
  ElInputNumber,
  ElLoading,
  ElMenu,
  ElMenuItem,
  ElOption,
  ElPagination,
  ElRadioButton,
  ElRadioGroup,
  ElResult,
  ElSelect,
  ElSkeleton,
  ElSwitch,
  ElTable,
  ElTableColumn,
  ElTag,
} from 'element-plus'
import { createPinia } from 'pinia'

import App from './App.vue'
import { setUnauthorizedHandler } from './api/http'
import router from './router'
import { useAuthStore } from './stores/auth'
import 'element-plus/dist/index.css'
import 'element-plus/theme-chalk/dark/css-vars.css'
import './style.css'
import './admin.css'

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
setUnauthorizedHandler(() => {
  useAuthStore(pinia).logout()
  if (router.currentRoute.value.meta.requiresAuth &&
      router.currentRoute.value.name !== 'admin-login') {
    void router.replace({
      name: 'admin-login',
      query: { redirect: router.currentRoute.value.fullPath },
    })
  }
})

app.use(router)
for (const component of [
  ElAlert,
  ElButton,
  ElForm,
  ElFormItem,
  ElInput,
  ElInputNumber,
  ElMenu,
  ElMenuItem,
  ElOption,
  ElPagination,
  ElRadioButton,
  ElRadioGroup,
  ElResult,
  ElSelect,
  ElSkeleton,
  ElSwitch,
  ElTable,
  ElTableColumn,
  ElTag,
]) {
  app.use(component)
}
app.use(ElLoading)
app.mount('#app')
