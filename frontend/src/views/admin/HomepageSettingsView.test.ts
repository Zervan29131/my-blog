import ElementPlus, { ElMessageBox } from 'element-plus'
import { flushPromises, mount } from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import {
  fetchAdminArticles,
  fetchFeaturedArticles,
  fetchHomepageDraft,
  fetchHomepagePublished,
  fetchSocialLinks,
  publishHomepage,
  resetHomepageDraft,
  saveHomepageDraft,
} from '../../api/admin'
import type { AdminHomepageConfig, HomepageModule } from '../../types/admin'
import HomepageSettingsView from './HomepageSettingsView.vue'

vi.mock('../../api/admin', () => ({
  addFeaturedArticle: vi.fn(),
  createSocialLink: vi.fn(),
  deleteFeaturedArticle: vi.fn(),
  deleteSocialLink: vi.fn(),
  fetchAdminArticles: vi.fn(),
  fetchFeaturedArticles: vi.fn(),
  fetchHomepageDraft: vi.fn(),
  fetchHomepagePublished: vi.fn(),
  fetchSocialLinks: vi.fn(),
  publishHomepage: vi.fn(),
  reorderFeaturedArticles: vi.fn(),
  reorderSocialLinks: vi.fn(),
  resetHomepageDraft: vi.fn(),
  saveHomepageDraft: vi.fn(),
  updateFeaturedArticleVisibility: vi.fn(),
  updateSocialLink: vi.fn(),
}))

const modules: HomepageModule[] = [
  {
    type: 'hero', enabled: true, sort_order: 10,
    config: {
      eyebrow: 'WELCOME', title: '在字里行间', highlight_text: '字里行间', description: '记录值得保存的思考。',
      image_url: '', background_image_url: '', layout: 'center',
      primary_button: { enabled: true, text: '开始阅读', url: '/archive', link_type: 'internal', open_in_new_tab: false },
      secondary_button: { enabled: true, text: '关于本站', url: '/about', link_type: 'internal', open_in_new_tab: false },
    },
  },
  {
    type: 'about', enabled: true, sort_order: 20,
    config: { title: '关于我', description: '保持好奇', content: '个人简介正文', image_url: '', image_position: 'none' },
  },
  {
    type: 'featured_articles', enabled: true, sort_order: 30,
    config: { title: '推荐文章', description: '精选内容', limit: 3 },
  },
  {
    type: 'latest_articles', enabled: true, sort_order: 40,
    config: { title: '最新文章', description: '', limit: 10, show_summary: true, show_date: true, show_comment_count: true, show_view_all: true },
  },
  {
    type: 'tech_stack', enabled: false, sort_order: 50,
    config: { title: '技术栈', description: '', items: [{ name: 'Vue 3', description: '', icon_url: '', url: '', is_visible: true, sort_order: 10 }] },
  },
  {
    type: 'social_links', enabled: true, sort_order: 60,
    config: { title: '找到我', description: '' },
  },
]

const draftConfig: AdminHomepageConfig = {
  status: 'draft', version: 2, modules,
  updated_at: '2026-08-01T08:00:00Z', published_at: null,
}
const publishedConfig: AdminHomepageConfig = {
  ...draftConfig, status: 'published', version: 2,
  modules: JSON.parse(JSON.stringify(modules)) as HomepageModule[],
  published_at: '2026-08-01T07:00:00Z',
}

describe('HomepageSettingsView', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    vi.mocked(fetchHomepageDraft).mockResolvedValue(draftConfig)
    vi.mocked(fetchHomepagePublished).mockResolvedValue(publishedConfig)
    vi.mocked(fetchFeaturedArticles).mockResolvedValue([])
    vi.mocked(fetchSocialLinks).mockResolvedValue([])
    vi.mocked(fetchAdminArticles).mockResolvedValue({ items: [], page: 1, page_size: 100, total: 0, total_pages: 0 })
    vi.mocked(resetHomepageDraft).mockResolvedValue(draftConfig)
    vi.mocked(publishHomepage).mockResolvedValue({ version: 3, published_at: '2026-08-01T09:00:00Z' })
  })

  it('loads all fixed modules and tracks order changes as a draft', async () => {
    const wrapper = mount(HomepageSettingsView, { global: { plugins: [ElementPlus] } })
    await flushPromises()

    expect(wrapper.text()).toContain('Hero 欢迎区域')
    expect(wrapper.text()).toContain('个人简介')
    expect(wrapper.text()).toContain('推荐文章')
    expect(wrapper.text()).toContain('最新文章')
    expect(wrapper.text()).toContain('技术栈')
    expect(wrapper.text()).toContain('社交链接')

    await wrapper.get('button[aria-label="下移Hero 欢迎区域"]').trigger('click')
    expect(wrapper.text()).toContain('有未保存修改')

    const discard = wrapper.findAll('button').find((button) => button.text() === '放弃本地修改')
    await discard?.trigger('click')
    expect(wrapper.text()).not.toContain('有未保存修改')
  })

  it('rejects unsafe Hero image URLs before saving', async () => {
    const wrapper = mount(HomepageSettingsView, { global: { plugins: [ElementPlus] } })
    await flushPromises()

    await wrapper.get('input[placeholder="https://example.com/avatar.jpg"]').setValue('javascript:alert(1)')
    const save = wrapper.findAll('button').find((button) => button.text() === '保存草稿')
    await save?.trigger('click')
    await flushPromises()

    expect(saveHomepageDraft).not.toHaveBeenCalled()
    expect(wrapper.text()).toContain('Hero 图片地址必须是有效的 HTTP 或 HTTPS 地址')
  })

  it('saves the draft before publishing and confirms the publish action', async () => {
    vi.mocked(saveHomepageDraft).mockImplementation(async (savedModules) => ({
      ...draftConfig,
      modules: savedModules,
      updated_at: '2026-08-01T08:30:00Z',
    }))
    vi.spyOn(ElMessageBox, 'confirm').mockResolvedValue('confirm' as never)
    const wrapper = mount(HomepageSettingsView, { global: { plugins: [ElementPlus] } })
    await flushPromises()

    await wrapper.get('input[placeholder="例如：WELCOME"]').setValue('NEW WELCOME')
    const save = wrapper.findAll('button').find((button) => button.text() === '保存草稿')
    await save?.trigger('click')
    await flushPromises()

    expect(saveHomepageDraft).toHaveBeenCalledOnce()
    expect(wrapper.text()).not.toContain('有未保存修改')

    const publish = wrapper.findAll('button').find((button) => button.text() === '发布首页')
    await publish?.trigger('click')
    await flushPromises()

    expect(ElMessageBox.confirm).toHaveBeenCalledWith(
      expect.stringContaining('立即对所有访客生效'),
      '确认发布首页配置？',
      expect.any(Object),
    )
    expect(publishHomepage).toHaveBeenCalledOnce()
    expect(wrapper.text()).toContain('v3')
  })

  it('opens only a saved draft in the protected preview page', async () => {
    const open = vi.spyOn(window, 'open').mockImplementation(() => null)
    const wrapper = mount(HomepageSettingsView, { global: { plugins: [ElementPlus] } })
    await flushPromises()

    const preview = wrapper.findAll('button').find((button) => button.text() === '预览首页')
    await preview?.trigger('click')
    expect(open).toHaveBeenCalledWith('/preview/home', '_blank', 'noopener,noreferrer')

    await wrapper.get('input[placeholder="例如：WELCOME"]').setValue('LOCAL CHANGE')
    await preview?.trigger('click')

    expect(open).toHaveBeenCalledOnce()
    expect(wrapper.text()).toContain('请先保存当前修改，再预览首页草稿')
  })
})
