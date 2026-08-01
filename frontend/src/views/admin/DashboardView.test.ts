import ElementPlus from 'element-plus'
import { flushPromises, mount } from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import { fetchDashboardStats } from '../../api/admin'
import DashboardView from './DashboardView.vue'

vi.mock('../../api/admin', () => ({
  fetchDashboardStats: vi.fn(),
}))

const mockedFetchStats = vi.mocked(fetchDashboardStats)

describe('DashboardView', () => {
  beforeEach(() => mockedFetchStats.mockReset())

  it('shows all required content statistics', async () => {
    mockedFetchStats.mockResolvedValue({
      article_total: 8,
      article_published: 5,
      article_draft: 3,
      comment_pending: 4,
      comment_approved: 9,
    })

    const wrapper = mount(DashboardView, { global: { plugins: [ElementPlus] } })
    await flushPromises()

    expect(wrapper.text()).toContain('全部文章')
    expect(wrapper.text()).toContain('已发布')
    expect(wrapper.text()).toContain('草稿')
    expect(wrapper.text()).toContain('待审核评论')
    expect(wrapper.text()).toContain('已通过评论')
    expect(wrapper.findAll('.stat-card').map((card) => card.find('strong').text())).toEqual([
      '8', '5', '3', '4', '9',
    ])
  })
})
