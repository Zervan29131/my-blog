import { createPinia, setActivePinia } from 'pinia'
import { flushPromises, mount } from '@vue/test-utils'
import { createMemoryHistory, createRouter } from 'vue-router'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import { fetchPublicSiteConfig } from '../api/blog'
import PublicLayout from './PublicLayout.vue'

vi.mock('../api/blog', () => ({ fetchPublicSiteConfig: vi.fn() }))

describe('PublicLayout', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.mocked(fetchPublicSiteConfig).mockResolvedValue({
      site: {
        name: '动态博客', short_name: '动态站', description: '动态站点描述', title_suffix: '随笔',
        logo_url: '', favicon_url: 'https://example.com/favicon.ico', default_share_image_url: '',
      },
      navigation: [
        { name: '首页', url: '/', link_type: 'internal', open_in_new_tab: false },
        { name: 'GitHub', url: 'https://github.com/example', link_type: 'external', open_in_new_tab: true },
      ],
      social_links: [],
      footer: {
        copyright_name: '动态博客', start_year: 2024, additional_text: '保持写作',
        filing_number: '测试备案号', filing_url: 'https://example.com/filing',
        show_technology: true, technology_text: 'Built dynamically',
      },
    })
  })

  it('renders published branding, navigation, admin entry and footer settings', async () => {
    const router = createRouter({
      history: createMemoryHistory(),
      routes: [{ path: '/', name: 'home', component: { template: '<main />' } }],
    })
    await router.push('/')
    await router.isReady()
    const pinia = createPinia()
    setActivePinia(pinia)

    const wrapper = mount(PublicLayout, {
      global: { plugins: [pinia, router], stubs: { RouterView: true, ThemeToggle: true } },
    })
    await flushPromises()

    expect(wrapper.get('.brand strong').text()).toBe('动态站')
    expect(wrapper.get('.desktop-nav a[href="https://github.com/example"]').attributes('target')).toBe('_blank')
    expect(wrapper.get('.desktop-nav a[href="/admin/login"]').text()).toBe('后台')
    expect(wrapper.get('.mobile-nav a[href="/admin/login"]').text()).toContain('后台入口')
    expect(wrapper.get('.site-footer').text()).toContain('2024–2026 动态博客')
    expect(wrapper.get('.site-footer').text()).toContain('Built dynamically')
    expect(document.title).toBe('动态博客 | 随笔')
  })
})
