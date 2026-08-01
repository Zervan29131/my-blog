import ElementPlus from 'element-plus'
import { flushPromises, mount } from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import { fetchAdminComments, updateCommentStatus } from '../../api/admin'
import CommentListView from './CommentListView.vue'

vi.mock('../../api/admin', () => ({
  fetchAdminComments: vi.fn(),
  updateCommentStatus: vi.fn(),
  deleteComment: vi.fn(),
}))

const mockedFetchComments = vi.mocked(fetchAdminComments)
const mockedUpdateStatus = vi.mocked(updateCommentStatus)

const pendingPage = {
  items: [{
    id: 3,
    article: { id: 1, title: '关联文章', slug: 'linked-post' },
    nickname: '等候审核的读者',
    email: 'reader@example.com',
    content: '这是一条等待审核的评论。',
    status: 'pending' as const,
    created_at: '2026-08-01T08:00:00Z',
    updated_at: '2026-08-01T08:00:00Z',
  }],
  page: 1,
  page_size: 20,
  total: 1,
  total_pages: 1,
}

describe('CommentListView', () => {
  beforeEach(() => {
    mockedFetchComments.mockReset()
    mockedUpdateStatus.mockReset()
  })

  it('loads pending comments and approves one', async () => {
    mockedFetchComments.mockResolvedValue(pendingPage)
    mockedUpdateStatus.mockResolvedValue()
    const wrapper = mount(CommentListView, { global: { plugins: [ElementPlus] } })
    await flushPromises()

    expect(mockedFetchComments).toHaveBeenCalledWith('pending', 1, 20)
    expect(wrapper.text()).toContain('这是一条等待审核的评论。')
    expect(wrapper.text()).toContain('关联文章')

    const approveButton = wrapper.findAll('button').find((button) => button.text() === '通过')
    await approveButton?.trigger('click')
    await flushPromises()

    expect(mockedUpdateStatus).toHaveBeenCalledWith(3, 'approved')
  })
})
