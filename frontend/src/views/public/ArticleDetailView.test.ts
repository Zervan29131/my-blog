import { flushPromises, mount } from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import ArticleDetailView from './ArticleDetailView.vue'
import { fetchArticle, fetchComments } from '../../api/blog'

vi.mock('../../api/blog', () => ({
  fetchArticle: vi.fn(),
  fetchComments: vi.fn(),
  submitComment: vi.fn(),
}))

const mockedFetchArticle = vi.mocked(fetchArticle)
const mockedFetchComments = vi.mocked(fetchComments)

describe('ArticleDetailView', () => {
  beforeEach(() => {
    mockedFetchArticle.mockReset()
    mockedFetchComments.mockReset()
  })

  it('renders sanitized Markdown and approved comments', async () => {
    mockedFetchArticle.mockResolvedValue({
      id: 1,
      title: '安全渲染文章',
      slug: 'safe-post',
      summary: '关于安全渲染的摘要。',
      content: '# 正文标题\n\n<script>window.hacked=true</script>正文内容',
      published_at: '2026-08-01T08:00:00Z',
      created_at: '2026-08-01T08:00:00Z',
      updated_at: '2026-08-01T08:00:00Z',
    })
    mockedFetchComments.mockResolvedValue({
      items: [
        {
          id: 2,
          nickname: '审核读者',
          content: '这是一条已公开评论。',
          created_at: '2026-08-01T09:00:00Z',
        },
      ],
      page: 1,
      page_size: 20,
      total: 1,
      total_pages: 1,
    })

    const wrapper = mount(ArticleDetailView, { props: { slug: 'safe-post' } })
    await flushPromises()

    expect(wrapper.text()).toContain('安全渲染文章')
    expect(wrapper.get('.markdown-body h1').text()).toBe('正文标题')
    expect(wrapper.find('.markdown-body script').exists()).toBe(false)
    expect(wrapper.text()).toContain('审核读者')
    expect(wrapper.text()).toContain('这是一条已公开评论。')
  })
})
