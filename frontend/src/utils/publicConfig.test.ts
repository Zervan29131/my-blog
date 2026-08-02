import { describe, expect, it, vi } from 'vitest'

import { normalizePublicHomepage, normalizePublicSiteConfig } from './publicConfig'

describe('public configuration normalization', () => {
  it('rejects unsafe links and preserves safe internal and external navigation', () => {
    const result = normalizePublicSiteConfig({
      site: { name: '测试站点', logo_url: 'javascript:alert(1)' },
      navigation: [
        { name: '归档', url: '/archive', link_type: 'internal', open_in_new_tab: false },
        { name: '危险', url: 'javascript:alert(1)', link_type: 'external', open_in_new_tab: true },
        { name: '外部', url: 'https://example.com', link_type: 'external', open_in_new_tab: true },
      ],
      social_links: [{ platform: 'email', display_name: '邮箱', url: 'mailto:test@example.com' }],
      footer: {},
    })

    expect(result.site.logo_url).toBe('')
    expect(result.navigation.map((item) => item.name)).toEqual(['归档', '外部'])
    expect(result.social_links[0]?.url).toBe('mailto:test@example.com')
  })

  it('skips unknown and empty modules while sorting known modules', () => {
    const warn = vi.spyOn(console, 'warn').mockImplementation(() => undefined)
    const result = normalizePublicHomepage({
      version: 3,
      modules: [
        { type: 'unknown', sort_order: 1, config: {} },
        { type: 'social_links', sort_order: 5, config: { title: '联系', links: [] } },
        { type: 'latest_articles', sort_order: 40, config: { title: '文章', limit: 3, articles: [] } },
        { type: 'hero', sort_order: 10, config: { title: '首页', description: '说明', layout: 'center', primary_button: {}, secondary_button: {} } },
      ],
    })

    expect(result.modules.map((module) => module.type)).toEqual(['hero', 'latest_articles'])
    expect(warn).toHaveBeenCalled()
    warn.mockRestore()
  })
})
