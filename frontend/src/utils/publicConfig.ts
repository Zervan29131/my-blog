import type {
  ArticleSummary,
  PublicAboutConfig,
  PublicArticlesModuleConfig,
  PublicHeroConfig,
  PublicHomepageButton,
  PublicHomepageConfig,
  PublicHomepageModule,
  PublicHomepageModuleType,
  PublicLatestArticlesConfig,
  PublicNavigationItem,
  PublicSiteConfig,
  PublicSocialLink,
  PublicSocialLinksConfig,
  PublicTechItem,
  PublicTechStackConfig,
} from '../types/blog'

export const DEFAULT_PUBLIC_SITE_CONFIG: PublicSiteConfig = {
  site: {
    name: '字里行间',
    short_name: '字里行间',
    description: '记录开发、阅读与生活中的思考。',
    title_suffix: '字里行间',
    logo_url: '',
    favicon_url: '',
    default_share_image_url: '',
  },
  navigation: [
    { name: '首页', url: '/', link_type: 'internal', open_in_new_tab: false },
    { name: '归档', url: '/archive', link_type: 'internal', open_in_new_tab: false },
    { name: '关于', url: '/about', link_type: 'internal', open_in_new_tab: false },
  ],
  social_links: [],
  footer: {
    copyright_name: '字里行间',
    start_year: new Date().getFullYear(),
    additional_text: '',
    filing_number: '',
    filing_url: '',
    show_technology: true,
    technology_text: 'Built with Vue 3 and Go',
  },
}

const knownModuleTypes = new Set<PublicHomepageModuleType>([
  'hero', 'about', 'featured_articles', 'latest_articles', 'tech_stack', 'social_links',
])

function record(value: unknown): Record<string, unknown> | null {
  return value !== null && typeof value === 'object' && !Array.isArray(value)
    ? value as Record<string, unknown>
    : null
}

function text(value: unknown, fallback = '', maximum = 2000): string {
  return typeof value === 'string' ? Array.from(value.trim()).slice(0, maximum).join('') : fallback
}

function bool(value: unknown, fallback = false): boolean {
  return typeof value === 'boolean' ? value : fallback
}

function integer(value: unknown, fallback: number, minimum: number, maximum: number): number {
  return typeof value === 'number' && Number.isInteger(value) && value >= minimum && value <= maximum
    ? value
    : fallback
}

function httpURL(value: unknown): string {
  const candidate = text(value, '', 500)
  if (!candidate) return ''
  try {
    const parsed = new URL(candidate)
    return parsed.protocol === 'http:' || parsed.protocol === 'https:' ? candidate : ''
  } catch {
    return ''
  }
}

function internalURL(value: unknown): string {
  const candidate = text(value, '', 500)
  return candidate.startsWith('/') && !candidate.startsWith('//') && !candidate.includes('\n') && !candidate.includes('\r')
    ? candidate
    : ''
}

function publicLink(value: unknown, allowMailto = false): string {
  const candidate = text(value, '', 500)
  if (allowMailto && /^mailto:[^\s@]+@[^\s@]+\.[^\s@]+$/i.test(candidate)) return candidate
  return httpURL(candidate)
}

function normalizeNavigation(value: unknown): PublicNavigationItem[] {
  if (!Array.isArray(value)) return DEFAULT_PUBLIC_SITE_CONFIG.navigation.map((item) => ({ ...item }))
  return value.slice(0, 10).flatMap((raw) => {
    const item = record(raw)
    if (!item) return []
    const linkType = item.link_type === 'external' ? 'external' : 'internal'
    const url = linkType === 'external' ? httpURL(item.url) : internalURL(item.url)
    const name = text(item.name, '', 20)
    return name && url ? [{ name, url, link_type: linkType, open_in_new_tab: bool(item.open_in_new_tab) }] : []
  })
}

function normalizeSocialLinks(value: unknown): PublicSocialLink[] {
  if (!Array.isArray(value)) return []
  return value.slice(0, 15).flatMap((raw) => {
    const item = record(raw)
    if (!item) return []
    const platform = text(item.platform, 'custom', 30)
    const displayName = text(item.display_name, '', 30)
    const url = publicLink(item.url, platform === 'email')
    return displayName && url ? [{ platform, display_name: displayName, url }] : []
  })
}

export function normalizePublicSiteConfig(value: unknown): PublicSiteConfig {
  const input = record(value)
  if (!input) return structuredClone(DEFAULT_PUBLIC_SITE_CONFIG)
  const site = record(input.site) ?? {}
  const footer = record(input.footer) ?? {}
  const name = text(site.name, DEFAULT_PUBLIC_SITE_CONFIG.site.name, 50) || DEFAULT_PUBLIC_SITE_CONFIG.site.name
  return {
    site: {
      name,
      short_name: text(site.short_name, name, 20) || name,
      description: text(site.description, DEFAULT_PUBLIC_SITE_CONFIG.site.description, 200) || DEFAULT_PUBLIC_SITE_CONFIG.site.description,
      title_suffix: text(site.title_suffix, name, 50),
      logo_url: httpURL(site.logo_url),
      favicon_url: httpURL(site.favicon_url),
      default_share_image_url: httpURL(site.default_share_image_url),
    },
    navigation: normalizeNavigation(input.navigation),
    social_links: normalizeSocialLinks(input.social_links),
    footer: {
      copyright_name: text(footer.copyright_name, name, 50) || name,
      start_year: typeof footer.start_year === 'number' && footer.start_year >= 1900 && footer.start_year <= new Date().getFullYear()
        ? Math.trunc(footer.start_year)
        : null,
      additional_text: text(footer.additional_text, '', 200),
      filing_number: text(footer.filing_number, '', 100),
      filing_url: httpURL(footer.filing_url),
      show_technology: bool(footer.show_technology, true),
      technology_text: text(footer.technology_text, '', 100),
    },
  }
}

function normalizeButton(value: unknown, fallback: PublicHomepageButton): PublicHomepageButton {
  const input = record(value) ?? {}
  const linkType = input.link_type === 'external' ? 'external' : 'internal'
  const url = linkType === 'external' ? httpURL(input.url) : internalURL(input.url)
  return {
    enabled: bool(input.enabled, fallback.enabled) && Boolean(url),
    text: text(input.text, fallback.text, 20) || fallback.text,
    url: url || fallback.url,
    link_type: linkType,
    open_in_new_tab: bool(input.open_in_new_tab, fallback.open_in_new_tab),
  }
}

function normalizeHero(value: unknown): PublicHeroConfig {
  const input = record(value) ?? {}
  return {
    eyebrow: text(input.eyebrow, 'WELCOME', 50),
    title: text(input.title, '在字里行间，保存思想的回声。', 100) || '在字里行间，保存思想的回声。',
    highlight_text: text(input.highlight_text, '', 50),
    description: text(input.description, '记录开发、阅读与生活中的思考。', 300) || '记录开发、阅读与生活中的思考。',
    image_url: httpURL(input.image_url),
    background_image_url: httpURL(input.background_image_url),
    layout: input.layout === 'left' ? 'left' : 'center',
    primary_button: normalizeButton(input.primary_button, {
      enabled: true, text: '开始阅读', url: '/archive', link_type: 'internal', open_in_new_tab: false,
    }),
    secondary_button: normalizeButton(input.secondary_button, {
      enabled: true, text: '关于本站', url: '/about', link_type: 'internal', open_in_new_tab: false,
    }),
  }
}

function normalizeAbout(value: unknown): PublicAboutConfig {
  const input = record(value) ?? {}
  const position = input.image_position === 'left' || input.image_position === 'right' ? input.image_position : 'none'
  return {
    title: text(input.title, '关于我', 100) || '关于我',
    description: text(input.description, '', 200),
    content: text(input.content, '这里记录开发、阅读与生活中的思考。', 2000) || '这里记录开发、阅读与生活中的思考。',
    image_url: httpURL(input.image_url),
    image_position: position,
  }
}

function normalizeArticle(value: unknown): ArticleSummary | null {
  const input = record(value)
  if (!input) return null
  const id = integer(input.id, 0, 1, Number.MAX_SAFE_INTEGER)
  const title = text(input.title, '', 200)
  const slug = text(input.slug, '', 200)
  const publishedAt = text(input.published_at, '', 100)
  if (!id || !title || !slug || !publishedAt) return null
  return {
    id,
    title,
    slug,
    summary: text(input.summary, '', 1000),
    published_at: publishedAt,
    comment_count: integer(input.comment_count, 0, 0, Number.MAX_SAFE_INTEGER),
  }
}

function normalizeArticles(value: unknown): ArticleSummary[] {
  return Array.isArray(value) ? value.map(normalizeArticle).filter((item): item is ArticleSummary => Boolean(item)) : []
}

function normalizeFeatured(value: unknown): PublicArticlesModuleConfig {
  const input = record(value) ?? {}
  return {
    title: text(input.title, '推荐文章', 100) || '推荐文章',
    description: text(input.description, '', 200),
    limit: integer(input.limit, 3, 1, 10),
    articles: normalizeArticles(input.articles),
  }
}

function normalizeLatest(value: unknown): PublicLatestArticlesConfig {
  const input = record(value) ?? {}
  return {
    title: text(input.title, '最新文章', 100) || '最新文章',
    description: text(input.description, '', 200),
    limit: integer(input.limit, 10, 3, 20),
    show_summary: bool(input.show_summary, true),
    show_date: bool(input.show_date, true),
    show_comment_count: bool(input.show_comment_count, true),
    show_view_all: bool(input.show_view_all, true),
    articles: normalizeArticles(input.articles),
  }
}

function normalizeTech(value: unknown): PublicTechStackConfig {
  const input = record(value) ?? {}
  const items = Array.isArray(input.items) ? input.items.slice(0, 20).flatMap((raw): PublicTechItem[] => {
    const item = record(raw)
    const name = text(item?.name, '', 30)
    return name ? [{ name, description: text(item?.description, '', 100), icon_url: httpURL(item?.icon_url), url: httpURL(item?.url) }] : []
  }) : []
  return {
    title: text(input.title, '技术栈', 100) || '技术栈',
    description: text(input.description, '', 200),
    items,
  }
}

function normalizeSocialModule(value: unknown): PublicSocialLinksConfig {
  const input = record(value) ?? {}
  return {
    title: text(input.title, '找到我', 100) || '找到我',
    description: text(input.description, '', 200),
    links: normalizeSocialLinks(input.links),
  }
}

export function normalizePublicHomepage(value: unknown): PublicHomepageConfig {
  const input = record(value)
  const rawModules = Array.isArray(input?.modules) ? input.modules : []
  const modules: PublicHomepageModule[] = []
  const seen = new Set<PublicHomepageModuleType>()
  for (const raw of rawModules) {
    const module = record(raw)
    if (!module || typeof module.type !== 'string' || !knownModuleTypes.has(module.type as PublicHomepageModuleType)) {
      console.warn('Skipping unsupported homepage module', module?.type)
      continue
    }
    const type = module.type as PublicHomepageModuleType
    if (seen.has(type)) continue
    seen.add(type)
    const sortOrder = integer(module.sort_order, modules.length * 10 + 10, -100000, 100000)
    if (type === 'hero') modules.push({ type, sort_order: sortOrder, config: normalizeHero(module.config) })
    if (type === 'about') modules.push({ type, sort_order: sortOrder, config: normalizeAbout(module.config) })
    if (type === 'featured_articles') {
      const config = normalizeFeatured(module.config)
      if (config.articles.length) modules.push({ type, sort_order: sortOrder, config })
    }
    if (type === 'latest_articles') modules.push({ type, sort_order: sortOrder, config: normalizeLatest(module.config) })
    if (type === 'tech_stack') {
      const config = normalizeTech(module.config)
      if (config.items.length) modules.push({ type, sort_order: sortOrder, config })
    }
    if (type === 'social_links') {
      const config = normalizeSocialModule(module.config)
      if (config.links.length) modules.push({ type, sort_order: sortOrder, config })
    }
  }
  modules.sort((left, right) => left.sort_order - right.sort_order)
  return { version: integer(input?.version, 0, 0, Number.MAX_SAFE_INTEGER), modules }
}
