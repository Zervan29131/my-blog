import { describe, expect, it } from 'vitest'

import { renderArticleMarkdown, renderMarkdown } from './markdown'

describe('renderMarkdown', () => {
  it('renders common Markdown elements', () => {
    const html = renderMarkdown('# 标题\n\n- 第一项\n- 第二项\n\n`code`')

    expect(html).toContain('<h1>标题</h1>')
    expect(html).toContain('<li>第一项</li>')
    expect(html).toContain('<code>code</code>')
  })

  it('removes scripts and unsafe URLs', () => {
    const html = renderMarkdown(
      '<script>window.hacked = true</script>\n\n[危险链接](javascript:alert(1))',
    )

    expect(html).not.toContain('<script')
    expect(html).not.toContain('href="javascript:')
    expect(html).toContain('危险链接')
  })

  it('creates stable article heading anchors and an h2/h3 outline', () => {
    const result = renderArticleMarkdown('## 项目结构\n\n### 前端\n\n## 项目结构')

    expect(result.html).toContain('<h2 id="项目结构">项目结构</h2>')
    expect(result.html).toContain('<h2 id="项目结构-2">项目结构</h2>')
    expect(result.headings).toEqual([
      { id: '项目结构', text: '项目结构', level: 2 },
      { id: '前端', text: '前端', level: 3 },
      { id: '项目结构-2', text: '项目结构', level: 2 },
    ])
  })
})
