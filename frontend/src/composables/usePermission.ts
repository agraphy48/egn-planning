import { computed } from 'vue'
import { useAuthStore } from '../stores/authStore'

export function usePermission() {
  const authStore = useAuthStore()

  const can = (permission: string): boolean => {
    return authStore.can(permission)
  }

  const hasRole = (roles: string | string[]): boolean => {
    return authStore.hasRole(roles)
  }

  const user = computed(() => authStore.user)
  const role = computed(() => authStore.currentRole)
  const permissions = computed(() => authStore.permissions)
  const isSuperAdmin = computed(() => authStore.currentRole === 'super_admin')
  const isPM = computed(() => authStore.currentRole === 'pm')
  const isMember = computed(() => authStore.currentRole === 'member')
  const isClient = computed(() => authStore.currentRole === 'client')

  return {
    can,
    hasRole,
    user,
    role,
    permissions,
    isSuperAdmin,
    isPM,
    isMember,
    isClient,
  }
}
