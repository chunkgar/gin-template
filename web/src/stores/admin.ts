import { ref } from 'vue'
import { defineStore } from 'pinia'
import type { AdminUser } from '@/utils/auth'

export const useAdminStore = defineStore('admin', () => {
  const user = ref<AdminUser | null>(null)

  function setUser(u: AdminUser) {
    user.value = u
  }

  function clearUser() {
    user.value = null
  }

  return { user, setUser, clearUser }
})
