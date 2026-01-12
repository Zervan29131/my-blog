import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'

import App from './App.vue'
import router from './router'

// import './style.css' // 如果你有全局样式的话，通常 Vite 会自动生成这个文件

const app = createApp(App)

// 1. 注册 Pinia (状态管理)
app.use(createPinia())

// 2. 注册 Router (路由)
app.use(router)

// 3. 注册 Element Plus (UI 组件库)
app.use(ElementPlus)

// 4. 注册 Element Plus 的所有图标
// 这样你就可以在组件里直接使用 <User /> 或 <Lock /> 图标了
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

app.mount('#app')