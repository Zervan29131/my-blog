import { createPinia, setActivePinia } from 'pinia'
import { flushPromises, mount } from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import { fetchHomepagePreview } from '../../api/admin'
import { fetchPublicSiteConfig } from '../../api/blog'
import type { PublicHomepageConfig, PublicSiteConfig } from '../../types/blog'
import PreviewHomeView from './PreviewHomeView.vue'

vi.mock('../../api/admin', () => ({ fetchHomepagePreview: vi.fn() }))
vi.mock('../../api/blog', () => ({ fetchPublicSiteConfig: vi.fn() }))

const siteConfig: PublicSiteConfig = {
  site: { name: '预览站点', short_name: '预览站', description: '站点描述', title_suffix: '预览站点', logo_url: '', favicon_url: '', default_share_image_url: '' },
  navigation: [], social_links: [],
  footer: { copyright_name: '预览站点', start_year: 2026, additional_text: '', filing_number: '', filing_url: '', show_technology: false, technology_text: '' },
}

const draftPreview: PublicHomepageConfig = {
  version: 8,
  modules: [{
    type: 'hero', sort_order: 10,
    config: {
      eyebrow: 'DRAFT', title: '仅在草稿中的标题', highlight_text: '草稿', description: '尚未发布的首页内容。',
      image_url: '', background_image_url: '', layout: 'center',
      primary_button: { enabled: true, text: '查看归档', url: '/archive', link_type: 'internal', open_in_new_tab: false },
      secondary_button: { enabled: false, text: '', url: '/about', link_type: 'internal', open_in_new_tab: false },
    },
  }],
}

function mountPreview() {
  const pinia = createPinia()
  setActivePinia(pinia)
  return mount(PreviewHomeView, {
    global: {
      plugins: [pinia],
      stubs: { RouterLink: { template: '<a><slot /></a>' } },
    },
  })
}

describe('PreviewHomeView', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    document.querySelector('meta[name="robots"]')?.remove()
    vi.mocked(fetchPublicSiteConfig).mockResolvedValue(siteConfig)
    vi.mocked(fetchHomepagePreview).mockResolvedValue(draftPreview)
  })

  it('renders the saved draft with a preview marker and blocks indexing', async () => {
    const wrapper = mountPreview()
    await flushPromises()

    expect(fetchHomepagePreview).toHaveBeenCalledOnce()
    expect(wrapper.text()).toContain('预览模式')
    expect(wrapper.text()).toContain('当前显示已保存的首页草稿，不会影响公开首页')
    expect(wrapper.text()).toContain('仅在草稿中的标题')
    expect(document.querySelector('meta[name="robots"]')?.getAttribute('content')).toBe('noindex,nofollow')

    wrapper.unmount()
    expect(document.querySelector('meta[name="robots"]')).toBeNull()
  })

  it('shows an authenticated-preview error and supports retrying', async () => {
    vi.mocked(fetchHomepagePreview).mockRejectedValueOnce(new Error('unauthorized'))
      .mockResolvedValueOnce(draftPreview)
    const wrapper = mountPreview()
    await flushPromises()

    expect(wrapper.text()).toContain('请确认登录仍然有效')
    await wrapper.get('.error-state button').trigger('click')
    await flushPromises()

    expect(fetchHomepagePreview).toHaveBeenCalledTimes(2)
    expect(wrapper.text()).toContain('仅在草稿中的标题')
    wrapper.unmount()
  })
})
