import { computed, onMounted, onUnmounted, ref } from 'vue'

export type ThemeMode = 'light' | 'dark' | 'system'

export const BLOG_THEME_KEY = 'blog-theme'

const theme = ref<ThemeMode>('system')
let initialized = false

function isThemeMode(value: string | null): value is ThemeMode {
  return value === 'light' || value === 'dark' || value === 'system'
}

function prefersDark(): boolean {
  return typeof window !== 'undefined' && typeof window.matchMedia === 'function' && window.matchMedia('(prefers-color-scheme: dark)').matches
}

function applyTheme(): void {
  if (typeof document === 'undefined') return
  const dark = theme.value === 'dark' || (theme.value === 'system' && prefersDark())
  document.documentElement.classList.toggle('dark', dark)
  document.documentElement.style.colorScheme = dark ? 'dark' : 'light'
  document.querySelector<HTMLMetaElement>('meta[name="theme-color"]')?.setAttribute('content', dark ? '#1b1b1f' : '#ffffff')
}

function initializeTheme(): void {
  if (initialized || typeof window === 'undefined') return
  const stored = window.localStorage.getItem(BLOG_THEME_KEY)
  theme.value = isThemeMode(stored) ? stored : 'system'
  applyTheme()
  initialized = true
}

export function useTheme() {
  const media = typeof window === 'undefined' || typeof window.matchMedia !== 'function'
    ? null
    : window.matchMedia('(prefers-color-scheme: dark)')
  const effectiveTheme = computed<'light' | 'dark'>(() => {
    if (theme.value === 'system') return prefersDark() ? 'dark' : 'light'
    return theme.value
  })

  function setTheme(mode: ThemeMode): void {
    theme.value = mode
    window.localStorage.setItem(BLOG_THEME_KEY, mode)
    applyTheme()
  }

  function cycleTheme(): void {
    const order: ThemeMode[] = ['system', 'light', 'dark']
    setTheme(order[(order.indexOf(theme.value) + 1) % order.length])
  }

  function handleSystemThemeChange(): void {
    if (theme.value === 'system') applyTheme()
  }

  onMounted(() => {
    initializeTheme()
    media?.addEventListener('change', handleSystemThemeChange)
  })
  onUnmounted(() => media?.removeEventListener('change', handleSystemThemeChange))

  return { theme, effectiveTheme, setTheme, cycleTheme }
}
