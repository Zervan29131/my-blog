import { computed, ref } from 'vue'
import { defineStore } from 'pinia'

import { fetchCurrentAdministrator, loginAdministrator } from '../api/admin'
import { ADMIN_TOKEN_KEY } from '../api/http'
import type { Administrator } from '../types/admin'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem(ADMIN_TOKEN_KEY) || '')
  const currentAdmin = ref<Administrator | null>(null)
  const isAuthenticated = computed(() => Boolean(token.value))

  function persistToken(value: string) {
    token.value = value
    if (value) {
      localStorage.setItem(ADMIN_TOKEN_KEY, value)
    } else {
      localStorage.removeItem(ADMIN_TOKEN_KEY)
    }
  }

  async function login(username: string, password: string): Promise<void> {
    const result = await loginAdministrator(username, password)
    persistToken(result.token)
    try {
      currentAdmin.value = await fetchCurrentAdministrator()
    } catch (error) {
      logout()
      throw error
    }
  }

  async function fetchCurrentAdmin(): Promise<Administrator> {
    const administrator = await fetchCurrentAdministrator()
    currentAdmin.value = administrator
    return administrator
  }

  function logout() {
    persistToken('')
    currentAdmin.value = null
  }

  return {
    token,
    currentAdmin,
    isAuthenticated,
    login,
    logout,
    fetchCurrentAdmin,
  }
})
