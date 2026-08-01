import DOMPurify from 'dompurify'
import MarkdownIt, { type Env } from 'markdown-it'

export interface ArticleHeading {
  id: string
  text: string
  level: 2 | 3
}

interface MarkdownEnvironment extends Env {
  collectOutline?: boolean
  headings?: ArticleHeading[]
  slugCounts?: Record<string, number>
}

const markdownParser = new MarkdownIt({
  breaks: true,
  html: true,
  linkify: true,
})

function headingSlug(text: string): string {
  const normalized = text
    .trim()
    .toLowerCase()
    .replace(/<[^>]+>/g, '')
    .replace(/[^\p{Letter}\p{Number}\s-]/gu, '')
    .replace(/\s+/g, '-')
    .replace(/-+/g, '-')
    .replace(/^-|-$/g, '')
  return normalized || 'section'
}

const defaultHeadingOpen = markdownParser.renderer.rules.heading_open
markdownParser.renderer.rules.heading_open = (tokens, index, options, environment, self) => {
  const env = environment as MarkdownEnvironment
  if (!env.collectOutline) {
    return defaultHeadingOpen?.(tokens, index, options, environment, self) || self.renderToken(tokens, index, options)
  }

  const token = tokens[index]
  const level = Number(token.tag.slice(1))
  const text = tokens[index + 1]?.content || 'section'
  const base = headingSlug(text)
  env.slugCounts ||= {}
  const count = env.slugCounts[base] || 0
  env.slugCounts[base] = count + 1
  const id = count ? `${base}-${count + 1}` : base
  token.attrSet('id', id)
  if (level === 2 || level === 3) env.headings?.push({ id, text, level })
  return self.renderToken(tokens, index, options)
}

const defaultLinkOpen = markdownParser.renderer.rules.link_open
markdownParser.renderer.rules.link_open = (tokens, index, options, environment, self) => {
  const token = tokens[index]
  const href = String(token.attrGet('href') || '')
  if (/^https?:\/\//i.test(href)) {
    token.attrSet('target', '_blank')
    token.attrSet('rel', 'noopener noreferrer')
  }
  return defaultLinkOpen?.(tokens, index, options, environment, self) || self.renderToken(tokens, index, options)
}

const defaultImage = markdownParser.renderer.rules.image
markdownParser.renderer.rules.image = (tokens, index, options, environment, self) => {
  tokens[index].attrSet('loading', 'lazy')
  return defaultImage?.(tokens, index, options, environment, self) || self.renderToken(tokens, index, options)
}

markdownParser.renderer.rules.table_open = () => '<div class="table-wrapper"><table>\n'
markdownParser.renderer.rules.table_close = () => '</table></div>\n'

function sanitize(rendered: string): string {
  return DOMPurify.sanitize(rendered, { USE_PROFILES: { html: true } })
}

export function renderMarkdown(markdown: string): string {
  return sanitize(markdownParser.render(markdown))
}

export function renderArticleMarkdown(markdown: string): { html: string; headings: ArticleHeading[] } {
  const environment: MarkdownEnvironment = { collectOutline: true, headings: [], slugCounts: {} }
  const html = sanitize(markdownParser.render(markdown, environment))
  return { html, headings: environment.headings || [] }
}
