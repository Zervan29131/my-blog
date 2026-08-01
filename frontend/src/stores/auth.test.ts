import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import { fetchCurrentAdministrator, loginAdministrator } from '../api/admin'
import { ADMIN_TOKEN_KEY } from '../api/http'
import { useAuthStore } from './auth'

vi.mock('../api/admin', () => ({
  loginAdministrator: vi.fn(),
  fetchCurrentAdministrator: vi.fn(),
}))

const mockedLogin = vi.mocked(loginAdministrator)
const mockedFetchCurrent = vi.mocked(fetchCurrentAdministrator)
const administrator = {
  id: 1,
  username: 'admin',
  created_at: '2026-08-01T08:00:00Z',
  updated_at: '2026-08-01T08:00:00Z',
}

describe('auth store', () => {
  beforeEach(() => {
    localStorage.clear()
    setActivePinia(createPinia())
    mockedLogin.mockReset()
    mockedFetchCurrent.mockReset()
  })

  it('persists the JWT and loads the current administrator after login', async () => {
    mockedLogin.mockResolvedValue({ token: 'signed-token', expires_in: 86400 })
    mockedFetchCurrent.mockResolvedValue(administrator)
    const store = useAuthStore()

    await store.login('admin', 'correct-password')

    expect(store.isAuthenticated).toBe(true)
    expect(store.currentAdmin).toEqual(administrator)
    expect(localStorage.getItem(ADMIN_TOKEN_KEY)).toBe('signed-token')
  })

  it('clears all local authentication state on logout', () => {
    localStorage.setItem(ADMIN_TOKEN_KEY, 'old-token')
    setActivePinia(createPinia())
    const store = useAuthStore()

    store.logout()

    expect(store.token).toBe('')
    expect(store.currentAdmin).toBeNull()
    expect(localStorage.getItem(ADMIN_TOKEN_KEY)).toBeNull()
  })
})
