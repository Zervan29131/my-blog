import { flushPromises, mount } from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import CommentForm from './CommentForm.vue'
import { submitComment } from '../api/blog'

vi.mock('../api/blog', () => ({
  submitComment: vi.fn(),
}))

const mockedSubmitComment = vi.mocked(submitComment)

describe('CommentForm', () => {
  beforeEach(() => {
    mockedSubmitComment.mockReset()
  })

  it('submits only the visitor fields and shows the pending-review message', async () => {
    mockedSubmitComment.mockResolvedValue({
      data: {
        id: 1,
        nickname: '认真读者',
        content: '很有启发。',
        status: 'pending',
        created_at: '2026-08-01T08:00:00Z',
      },
      message: '评论已提交，审核通过后将会显示。',
    })
    const wrapper = mount(CommentForm, { props: { slug: 'hello-world' } })

    await wrapper.get('input[name="nickname"]').setValue('认真读者')
    await wrapper.get('input[name="email"]').setValue('reader@example.com')
    await wrapper.get('textarea[name="content"]').setValue('很有启发。')
    await wrapper.get('form').trigger('submit')
    await flushPromises()

    expect(mockedSubmitComment).toHaveBeenCalledWith('hello-world', {
      nickname: '认真读者',
      email: 'reader@example.com',
      content: '很有启发。',
    })
    expect(wrapper.text()).toContain('评论已提交，审核通过后将会显示。')
  })

  it('validates minimum lengths before making a request', async () => {
    const wrapper = mount(CommentForm, { props: { slug: 'hello-world' } })

    await wrapper.get('input[name="nickname"]').setValue('a')
    await wrapper.get('textarea[name="content"]').setValue('b')
    await wrapper.get('form').trigger('submit')

    expect(mockedSubmitComment).not.toHaveBeenCalled()
    expect(wrapper.text()).toContain('昵称需要填写 2～50 个字符。')
  })

  it('rejects an invalid optional email before making a request', async () => {
    const wrapper = mount(CommentForm, { props: { slug: 'hello-world' } })

    await wrapper.get('input[name="nickname"]').setValue('认真读者')
    await wrapper.get('input[name="email"]').setValue('not-an-email')
    await wrapper.get('textarea[name="content"]').setValue('有效评论内容')
    await wrapper.get('form').trigger('submit')

    expect(mockedSubmitComment).not.toHaveBeenCalled()
    expect(wrapper.text()).toContain('请输入有效的邮箱地址。')
  })
})
