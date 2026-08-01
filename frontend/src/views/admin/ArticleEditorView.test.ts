import ElementPlus from 'element-plus'
import { flushPromises, mount } from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import { createArticle } from '../../api/admin'
import ArticleEditorView from './ArticleEditorView.vue'

const push = vi.fn()

vi.mock('vue-router', () => ({
  useRouter: () => ({ push }),
}))

vi.mock('../../api/admin', () => ({
  createArticle: vi.fn(),
  fetchAdminArticle: vi.fn(),
  updateArticle: vi.fn(),
}))

const mockedCreateArticle = vi.mocked(createArticle)

describe('ArticleEditorView', () => {
  beforeEach(() => {
    push.mockReset()
    mockedCreateArticle.mockReset()
  })

  it('renders a sanitized preview and creates a draft article', async () => {
    mockedCreateArticle.mockResolvedValue({
      id: 1,
      title: '新文章',
      slug: 'new-post',
      summary: '摘要',
      content: '# 正文',
      status: 'draft',
      published_at: null,
      created_at: '2026-08-01T08:00:00Z',
      updated_at: '2026-08-01T08:00:00Z',
    })
    const wrapper = mount(ArticleEditorView, { global: { plugins: [ElementPlus] } })

    await wrapper.get('input[placeholder="文章标题"]').setValue('新文章')
    await wrapper.get('input[placeholder="article-slug"]').setValue('new-post')
    await wrapper.get('textarea[placeholder="文章列表中显示的简短摘要"]').setValue('摘要')
    await wrapper.get('textarea[placeholder="# 从这里开始写作"]').setValue('# 正文\n<script>alert(1)</script>')
    await flushPromises()

    expect(wrapper.get('.editor-preview h1').text()).toBe('正文')
    expect(wrapper.find('.editor-preview script').exists()).toBe(false)

    const saveDraftButton = wrapper.findAll('.editor-actions button').find((button) => button.text().includes('保存为草稿'))
    await saveDraftButton?.trigger('click')
    await flushPromises()

    expect(mockedCreateArticle).toHaveBeenCalledWith({
      title: '新文章',
      slug: 'new-post',
      summary: '摘要',
      content: '# 正文\n<script>alert(1)</script>',
      status: 'draft',
    })
    expect(push).toHaveBeenCalledWith('/admin/articles')
  })
})
