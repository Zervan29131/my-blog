import { createPinia, setActivePinia } from 'pinia'
import { flushPromises, mount } from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import { fetchArticles, fetchHomepageConfig, fetchPublicSiteConfig } from '../../api/blog'
import type { PublicHomepageConfig, PublicSiteConfig } from '../../types/blog'
import HomeView from './HomeView.vue'

vi.mock('../../api/blog', () => ({
  fetchArticles: vi.fn(),
  fetchHomepageConfig: vi.fn(),
  fetchPublicSiteConfig: vi.fn(),
}))

const siteConfig: PublicSiteConfig = {
  site: { name: '动态博客', short_name: '动态站', description: '站点描述', title_suffix: '随笔', logo_url: '', favicon_url: '', default_share_image_url: '' },
  navigation: [], social_links: [],
  footer: { copyright_name: '动态博客', start_year: 2024, additional_text: '', filing_number: '', filing_url: '', show_technology: true, technology_text: 'Vue + Go' },
}

const homepageConfig: PublicHomepageConfig = {
  version: 4,
  modules: [
    {
      type: 'hero', sort_order: 10,
      config: {
        eyebrow: 'HELLO', title: '动态首页标题', highlight_text: '动态', description: '来自已发布配置的描述。', image_url: '', background_image_url: '', layout: 'center',
        primary_button: { enabled: true, text: '阅读文章', url: '/archive', link_type: 'internal', open_in_new_tab: false },
        secondary_button: { enabled: false, text: '', url: '/about', link_type: 'internal', open_in_new_tab: false },
      },
    },
    {
      type: 'about', sort_order: 20,
      config: { title: '关于作者', description: '简介说明', content: '<script>window.hacked=true</script>\n\n安全的 **Markdown** 简介', image_url: '', image_position: 'none' },
    },
    {
      type: 'latest_articles', sort_order: 40,
      config: {
        title: '最近更新', description: '', limit: 3, show_summary: true, show_date: true, show_comment_count: true, show_view_all: true,
        articles: [{ id: 1, title: '第一篇文章', slug: 'first-post', summary: '文章摘要', published_at: '2026-08-01T08:00:00Z', comment_count: 3 }],
      },
    },
  ],
}

describe('HomeView', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    const pinia = createPinia()
    setActivePinia(pinia)
    vi.mocked(fetchPublicSiteConfig).mockResolvedValue(siteConfig)
    vi.mocked(fetchHomepageConfig).mockResolvedValue(homepageConfig)
  })

  it('maps published modules in order and sanitizes the about Markdown', async () => {
    const pinia = createPinia()
    setActivePinia(pinia)
    const wrapper = mount(HomeView, { global: { plugins: [pinia] } })
    await flushPromises()

    expect(wrapper.text()).toContain('动态首页标题')
    expect(wrapper.text()).toContain('关于作者')
    expect(wrapper.text()).toContain('最近更新')
    expect(wrapper.text()).toContain('第一篇文章')
    expect(wrapper.text()).toContain('3 条评论')
    expect(wrapper.find('.home-about-markdown script').exists()).toBe(false)
    expect(wrapper.get('.home-about-markdown strong').text()).toBe('Markdown')
    expect(wrapper.get('.article-card h3 a').attributes('href')).toBe('/articles/first-post')

    const text = wrapper.text()
    expect(text.indexOf('动态首页标题')).toBeLessThan(text.indexOf('关于作者'))
    expect(text.indexOf('关于作者')).toBeLessThan(text.indexOf('最近更新'))
  })

  it('keeps a basic article list and retry action when homepage config fails', async () => {
    vi.mocked(fetchHomepageConfig).mockRejectedValue(new Error('offline'))
    vi.mocked(fetchArticles).mockResolvedValue({
      items: [{ id: 2, title: '降级文章', slug: 'fallback', summary: '仍可阅读', published_at: '2026-08-01T08:00:00Z', comment_count: 0 }],
      page: 1, page_size: 10, total: 1, total_pages: 1,
    })
    const pinia = createPinia()
    setActivePinia(pinia)
    const wrapper = mount(HomeView, { global: { plugins: [pinia] } })
    await flushPromises()

    expect(wrapper.text()).toContain('首页配置加载失败')
    expect(wrapper.text()).toContain('重新加载配置')
    expect(wrapper.text()).toContain('降级文章')
    expect(fetchArticles).toHaveBeenCalledWith(1, 10)
  })
})
