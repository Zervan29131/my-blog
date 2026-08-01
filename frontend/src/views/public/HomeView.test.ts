import { flushPromises, mount } from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import HomeView from './HomeView.vue'
import { fetchArticles } from '../../api/blog'

vi.mock('../../api/blog', () => ({
  fetchArticles: vi.fn(),
}))

const mockedFetchArticles = vi.mocked(fetchArticles)

describe('HomeView', () => {
  beforeEach(() => {
    mockedFetchArticles.mockReset()
    window.scrollTo = vi.fn()
  })

  it('renders published article summaries and comment counts', async () => {
    mockedFetchArticles.mockResolvedValue({
      items: [
        {
          id: 1,
          title: '第一篇文章',
          slug: 'first-post',
          summary: '这是文章摘要。',
          published_at: '2026-08-01T08:00:00Z',
          comment_count: 3,
        },
      ],
      page: 1,
      page_size: 10,
      total: 1,
      total_pages: 1,
    })

    const wrapper = mount(HomeView)
    await flushPromises()

    expect(wrapper.text()).toContain('第一篇文章')
    expect(wrapper.text()).toContain('这是文章摘要。')
    expect(wrapper.text()).toContain('3 条评论')
    expect(wrapper.text()).toContain('Zervan')
    expect(wrapper.text()).toContain('全部文章')
    expect(wrapper.get('.article-card h3 a').attributes('href')).toBe('/articles/first-post')
    expect(wrapper.get('.site-overview a[href="/admin/login"]').text()).toContain('内容管理')
  })

  it('requests the next page from the pagination controls', async () => {
    mockedFetchArticles
      .mockResolvedValueOnce({ items: [], page: 1, page_size: 10, total: 12, total_pages: 2 })
      .mockResolvedValueOnce({ items: [], page: 2, page_size: 10, total: 12, total_pages: 2 })

    const wrapper = mount(HomeView)
    await flushPromises()
    await wrapper.get('button[aria-label="下一页"]').trigger('click')
    await flushPromises()

    expect(mockedFetchArticles).toHaveBeenLastCalledWith(2, 10)
  })

  it('shows a recoverable error state', async () => {
    mockedFetchArticles.mockRejectedValue(new Error('offline'))

    const wrapper = mount(HomeView)
    await flushPromises()

    expect(wrapper.text()).toContain('文章加载失败，请稍后重试。')
    expect(wrapper.get('.error-state button').text()).toBe('重新加载')
  })
})
