import ElementPlus from 'element-plus'
import { flushPromises, mount } from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import { fetchSiteSettings, updateSiteSettings } from '../../api/admin'
import type { SiteSettings, SiteSettingsInput } from '../../types/admin'
import SiteSettingsView from './SiteSettingsView.vue'

vi.mock('../../api/admin', () => ({
  fetchSiteSettings: vi.fn(),
  updateSiteSettings: vi.fn(),
}))

const mockedFetchSettings = vi.mocked(fetchSiteSettings)
const mockedUpdateSettings = vi.mocked(updateSiteSettings)

const siteSettings: SiteSettings = {
  id: 1,
  site_name: '字里行间',
  site_short_name: '字里行间',
  site_description: '记录开发、阅读与生活中的思考。',
  title_suffix: '字里行间',
  logo_url: null,
  favicon_url: null,
  default_share_image_url: null,
  copyright_name: '字里行间',
  start_year: 2024,
  additional_text: null,
  filing_number: null,
  filing_url: null,
  show_technology: true,
  technology_text: 'Built with Vue 3 and Go',
  created_at: '2026-08-01T00:00:00Z',
  updated_at: '2026-08-01T00:00:00Z',
}

describe('SiteSettingsView', () => {
  beforeEach(() => {
    mockedFetchSettings.mockReset()
    mockedUpdateSettings.mockReset()
    mockedFetchSettings.mockResolvedValue(siteSettings)
  })

  it('loads the three settings groups and resets unsaved changes', async () => {
    const wrapper = mount(SiteSettingsView, { global: { plugins: [ElementPlus] } })
    await flushPromises()

    expect(wrapper.text()).toContain('基础信息')
    expect(wrapper.text()).toContain('品牌资源')
    expect(wrapper.text()).toContain('页脚设置')

    const siteName = wrapper.get('input[placeholder="例如：字里行间"]')
    expect((siteName.element as HTMLInputElement).value).toBe('字里行间')
    await siteName.setValue('尚未保存的名称')
    expect(wrapper.text()).toContain('有未保存修改')

    const reset = wrapper.findAll('button').find((button) => button.text() === '重置')
    await reset?.trigger('click')
    expect((siteName.element as HTMLInputElement).value).toBe('字里行间')
  })

  it('validates unsafe image URLs before sending a request', async () => {
    const wrapper = mount(SiteSettingsView, { global: { plugins: [ElementPlus] } })
    await flushPromises()

    await wrapper.get('input[placeholder="https://example.com/logo.png"]').setValue('javascript:alert(1)')
    const save = wrapper.findAll('button').find((button) => button.text() === '保存设置')
    await save?.trigger('click')
    await flushPromises()

    expect(mockedUpdateSettings).not.toHaveBeenCalled()
    expect(wrapper.text()).toContain('Logo 地址必须是有效的 HTTP 或 HTTPS 地址')
  })

  it('normalizes optional fields and shows the saved state', async () => {
    mockedUpdateSettings.mockImplementation(async (input: SiteSettingsInput) => ({
      ...siteSettings,
      ...input,
      updated_at: '2026-08-01T01:00:00Z',
    }))
    const wrapper = mount(SiteSettingsView, { global: { plugins: [ElementPlus] } })
    await flushPromises()

    await wrapper.get('input[placeholder="例如：字里行间"]').setValue('  新站点  ')
    await wrapper.get('input[placeholder="用于空间有限的场景"]').setValue('   ')
    await wrapper.get('input[placeholder="https://example.com/logo.png"]').setValue('https://cdn.example.com/logo.png')
    const save = wrapper.findAll('button').find((button) => button.text() === '保存设置')
    await save?.trigger('click')
    await flushPromises()

    expect(mockedUpdateSettings).toHaveBeenCalledWith(expect.objectContaining({
      site_name: '新站点',
      site_short_name: null,
      logo_url: 'https://cdn.example.com/logo.png',
    }))
    expect(wrapper.text()).toContain('当前内容已与服务器同步')
  })
})
