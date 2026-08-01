import { mount } from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import ThemeToggle from './ThemeToggle.vue'
import { BLOG_THEME_KEY } from '../composables/useTheme'

describe('ThemeToggle', () => {
  beforeEach(() => {
    localStorage.clear()
    document.documentElement.classList.remove('dark')
    vi.stubGlobal('matchMedia', vi.fn(() => ({
      matches: false,
      addEventListener: vi.fn(),
      removeEventListener: vi.fn(),
    })))
  })

  it('cycles theme modes and persists the explicit choice', async () => {
    const wrapper = mount(ThemeToggle)

    expect(wrapper.get('button').attributes('aria-label')).toContain('跟随系统')
    await wrapper.get('button').trigger('click')

    expect(localStorage.getItem(BLOG_THEME_KEY)).toBe('light')
    expect(wrapper.get('button').attributes('aria-label')).toContain('浅色模式')
  })
})
