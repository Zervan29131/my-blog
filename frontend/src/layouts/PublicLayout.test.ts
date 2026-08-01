import { mount } from '@vue/test-utils'
import { createMemoryHistory, createRouter } from 'vue-router'
import { describe, expect, it } from 'vitest'

import PublicLayout from './PublicLayout.vue'

describe('PublicLayout', () => {
  it('exposes the administrator login in desktop and mobile navigation', async () => {
    const router = createRouter({
      history: createMemoryHistory(),
      routes: [{ path: '/', component: { template: '<main />' } }],
    })
    await router.push('/')
    await router.isReady()

    const wrapper = mount(PublicLayout, {
      global: {
        plugins: [router],
        stubs: {
          RouterView: true,
          ThemeToggle: true,
        },
      },
    })

    expect(wrapper.get('.desktop-nav a[href="/admin/login"]').text()).toBe('后台')
    expect(wrapper.get('.mobile-nav a[href="/admin/login"]').text()).toContain('后台入口')
  })
})
