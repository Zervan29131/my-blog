import { createPinia, setActivePinia } from 'pinia'
import { createMemoryHistory } from 'vue-router'
import { beforeEach, describe, expect, it } from 'vitest'

import { ADMIN_TOKEN_KEY } from '../api/http'
import { createAppRouter } from './index'

describe('administrator route guard', () => {
  beforeEach(() => {
    localStorage.clear()
    setActivePinia(createPinia())
  })

  it('redirects anonymous visitors to login and preserves the target path', async () => {
    const router = createAppRouter(createMemoryHistory())

    await router.push('/admin/articles/new')

    expect(router.currentRoute.value.name).toBe('admin-login')
    expect(router.currentRoute.value.query.redirect).toBe('/admin/articles/new')
  })

  it('keeps authenticated administrators out of the login page', async () => {
    localStorage.setItem(ADMIN_TOKEN_KEY, 'signed-token')
    setActivePinia(createPinia())
    const router = createAppRouter(createMemoryHistory())

    await router.push('/admin/login')

    expect(router.currentRoute.value.name).toBe('admin-dashboard')
  })

  it('provides archive, about, and a real not-found route', async () => {
    const router = createAppRouter(createMemoryHistory())

    await router.push('/archive')
    expect(router.currentRoute.value.name).toBe('archive')
    await router.push('/about')
    expect(router.currentRoute.value.name).toBe('about')
    await router.push('/missing-page')
    expect(router.currentRoute.value.name).toBe('not-found')
  })

  it('registers the authenticated site settings page', async () => {
    localStorage.setItem(ADMIN_TOKEN_KEY, 'signed-token')
    setActivePinia(createPinia())
    const router = createAppRouter(createMemoryHistory())

    await router.push('/admin/site/settings')

    expect(router.currentRoute.value.name).toBe('admin-site-settings')
  })
})
