import { defineStore } from 'pinia'
import type { User } from '../types/project'
import { apiClient, USE_MOCK } from '../services/api'
import { useProjectStore } from './projectStore'

// Role Permissions Matrix matching Golang Backend 1:1
export const ROLE_PERMISSIONS: Record<string, string[]> = {
  super_admin: [
    'user:manage',
    'user:role:update',
    'project:create',
    'project:read',
    'project:edit',
    'project:delete',
    'project:archive',
    'project:member:manage',
    'milestone:create',
    'milestone:read',
    'milestone:edit',
    'milestone:approve',
    'task:create',
    'task:read',
    'task:edit',
    'task:move',
    'task:delete',
    'subtask:manage',
    'attachment:upload',
    'comment:create',
    'comment:read',
  ],
  pm: [
    'project:create',
    'project:read',
    'project:edit',
    'project:archive',
    'project:member:manage',
    'milestone:create',
    'milestone:read',
    'milestone:edit',
    'milestone:approve',
    'task:create',
    'task:read',
    'task:edit',
    'task:move',
    'task:delete',
    'subtask:manage',
    'attachment:upload',
    'comment:create',
    'comment:read',
  ],
  member: [
    'project:read',
    'milestone:read',
    'task:create',
    'task:read',
    'task:edit',
    'task:move',
    'subtask:manage',
    'attachment:upload',
    'comment:create',
    'comment:read',
  ],
  client: [
    'project:read',
    'milestone:read',
    'milestone:approve',
    'task:read',
    'comment:create',
    'comment:read',
  ],
}

// Pre-seeded User profiles for role reference & development metadata
export const SEED_USERS: Record<string, User> = {
  agus: {
    id: 'usr_agus',
    name: 'Agus Setiawan',
    email: 'agus.setiawan@pulse.workspace',
    role: 'pm' as any,
    avatar_url: 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=100&auto=format&fit=crop&q=80',
    initials: 'AS',
  },
  adi: {
    id: 'usr_adi',
    name: 'Adi Saputra',
    email: 'adi.saputra@pulse.workspace',
    role: 'pm' as any,
    avatar_url: 'https://images.unsplash.com/photo-1500648767791-00dcc994a43e?w=100&auto=format&fit=crop&q=80',
    initials: 'AS',
  },
  yudistira: {
    id: 'usr_yudistira',
    name: 'Yudistira Sira Permana',
    email: 'yudistira@pulse.workspace',
    role: 'member' as any,
    avatar_url: 'https://images.unsplash.com/photo-1519085360753-af0119f7cbe7?w=100&auto=format&fit=crop&q=80',
    initials: 'YP',
  },
  hardianto: {
    id: 'usr_hardianto',
    name: 'Hardianto',
    email: 'hardianto@pulse.workspace',
    role: 'member' as any,
    avatar_url: 'https://images.unsplash.com/photo-1506794778202-cad84cf45f1d?w=100&auto=format&fit=crop&q=80',
    initials: 'HD',
  },
  dedi: {
    id: 'usr_dedi',
    name: 'Dedi Irawan',
    email: 'dedi.irawan@pulse.workspace',
    role: 'member' as any,
    avatar_url: 'https://images.unsplash.com/photo-1522075469751-3a6694fb2f61?w=100&auto=format&fit=crop&q=80',
    initials: 'DI',
  },
  super_admin: {
    id: 'usr_super_admin',
    name: 'System Super Admin',
    email: 'admin@pulse.workspace',
    role: 'super_admin' as any,
    avatar_url: 'https://images.unsplash.com/photo-1534528741775-53994a69daeb?w=100&auto=format&fit=crop&q=80',
    initials: 'SA',
  },
  client: {
    id: 'usr_client',
    name: 'Ir. H. Gunawan (Client)',
    email: 'client@pulse.workspace',
    role: 'client' as any,
    avatar_url: 'https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=100&auto=format&fit=crop&q=80',
    initials: 'HG',
  },
  pm: {
    id: 'usr_pm',
    name: 'Alex Rivera (PM)',
    email: 'pm@pulse.workspace',
    role: 'pm' as any,
    avatar_url: 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=100&auto=format&fit=crop&q=80',
    initials: 'AR',
  },
  member: {
    id: 'usr_member',
    name: 'Sarah Connor (Member)',
    email: 'member@pulse.workspace',
    role: 'member' as any,
    avatar_url: 'https://images.unsplash.com/photo-1494790108377-be9c29b29330?w=100&auto=format&fit=crop&q=80',
    initials: 'SC',
  },
}

function getInitialUser(): User | null {
  const raw = localStorage.getItem('auth_user')
  if (raw) {
    try {
      return JSON.parse(raw)
    } catch {
      localStorage.removeItem('auth_user')
      return null
    }
  }
  // Secure default: Unauthenticated by default, never leak seed accounts
  return null
}

export const useAuthStore = defineStore('auth', {
  state: () => {
    const user = getInitialUser()
    const token = localStorage.getItem('auth_token') || ''
    const role = user?.role || ''
    return {
      user: user as User | null,
      token: token,
      permissions: (user && role && ROLE_PERMISSIONS[role] ? ROLE_PERMISSIONS[role] : []) || [],
      isLoading: false,
      isLoggingOut: false,
      allowQuickLogin: false,
      isChangePasswordModalOpen: false,
      useMock: USE_MOCK,
    }
  },

  getters: {
    isAuthenticated(state): boolean {
      return (
        !!state.token &&
        state.token.length > 20 &&
        state.token !== 'jwt-session-active' &&
        !state.token.startsWith('jwt-token-') &&
        !state.token.startsWith('mock-') &&
        !!state.user
      )
    },

    currentRole(state): string {
      return state.user?.role || ''
    },

    hasRole: (state) => (roles: string | string[]): boolean => {
      if (!state.user) return false
      if (state.user.role === 'super_admin') return true
      if (Array.isArray(roles)) {
        return roles.includes(state.user.role)
      }
      return state.user.role === roles
    },

    can: (state) => (permission: string): boolean => {
      if (!state.user) return false
      if (state.user.role === 'super_admin') return true
      return state.permissions.includes(permission)
    },
  },

  actions: {
    async fetchAuthConfig() {
      try {
        const response = await apiClient.get('/auth/config')
        this.allowQuickLogin = response.data?.allow_quick_login ?? false
      } catch {
        this.allowQuickLogin = !import.meta.env.PROD
      }
    },

    async fetchUser() {
      if (!this.token) {
        this.forceLogout()
        return
      }

      try {
        this.isLoading = true
        const response = await apiClient.get('/auth/me')
        const data = response.data
        this.user = {
          id: data.id,
          name: data.name,
          email: data.email,
          role: data.role,
          avatar_url: data.avatar_url,
          initials: data.name.substring(0, 2).toUpperCase(),
        }
        this.permissions = data.permissions || ROLE_PERMISSIONS[data.role] || []
        localStorage.setItem('auth_user', JSON.stringify(this.user))
      } catch (err: any) {
        if (err?.response?.status === 401) {
          console.warn('[Auth] Current session expired or revoked on server.')
          this.forceLogout()
          throw err
        }
        if (this.user?.role && ROLE_PERMISSIONS[this.user.role]) {
          this.permissions = ROLE_PERMISSIONS[this.user.role]
        }
      } finally {
        this.isLoading = false
      }
    },

    async login(email: string, password: string): Promise<boolean> {
      this.isLoading = true
      try {
        const response = await apiClient.post('/auth/login', { email, password })
        const { token, user, permissions } = response.data
        this.token = token
        this.user = user
        this.permissions = permissions || ROLE_PERMISSIONS[user.role] || []
        localStorage.setItem('auth_token', token)
        localStorage.setItem('auth_user', JSON.stringify(user))

        // Synchronize authenticated user with project store
        try {
          const pStore = useProjectStore()
          pStore.currentUser = user
        } catch (_) {}

        // Broadcast new login to other tabs on the same browser
        try {
          if (typeof BroadcastChannel !== 'undefined') {
            const bc = new BroadcastChannel('pulse_auth_sync')
            bc.postMessage({ type: 'NEW_LOGIN', userId: user.id, timestamp: Date.now() })
            bc.close()
          }
        } catch (_) {}

        return true
      } finally {
        this.isLoading = false
      }
    },

    // Development user switcher: authenticates via backend using real credentials
    async switchUser(userKey: string): Promise<User> {
      const selected = SEED_USERS[userKey]
      if (!selected) {
        throw new Error(`Pengguna ${userKey} tidak ditemukan`)
      }

      // Perform real credentials login to obtain signed JWT
      await this.login(selected.email, 'Password123!')
      return this.user!
    },

    async logout(): Promise<void> {
      if (this.isLoggingOut) return
      this.isLoggingOut = true

      const currentToken = this.token
      try {
        // 1. Notify backend to revoke session and blacklist JWT
        if (currentToken && currentToken.length > 20 && !currentToken.startsWith('mock')) {
          await apiClient.post('/auth/logout', {}, {
            headers: { Authorization: `Bearer ${currentToken}` }
          }).catch((err) => {
            console.warn('[Auth] Backend logout notification notice:', err?.message)
          })
        }
      } finally {
        // 2. Clear frontend reactive state
        this.token = ''
        this.user = null
        this.permissions = []
        this.isLoggingOut = false

        // 3. Clear persisted auth credentials
        localStorage.removeItem('auth_token')
        localStorage.removeItem('auth_user')

        // 4. Invalidate application domain caches in projectStore
        try {
          const pStore = useProjectStore()
          pStore.resetStore()
        } catch (_) {}

        // 5. Broadcast logout signal to all other open tabs
        try {
          if (typeof BroadcastChannel !== 'undefined') {
            const bc = new BroadcastChannel('pulse_auth_sync')
            bc.postMessage({ type: 'LOGOUT', timestamp: Date.now() })
            bc.close()
          }
        } catch (_) {}
      }
    },

    // Change user password action
    async changePassword(currentPassword: string, newPassword: string, confirmPassword: string) {
      const response = await apiClient.post('/auth/change-password', {
        current_password: currentPassword,
        new_password: newPassword,
        confirm_password: confirmPassword,
      })
      return response.data
    },

    // Force logout triggered on 401 Unauthorized or multi-tab sync
    forceLogout() {
      this.token = ''
      this.user = null
      this.permissions = []
      this.isLoggingOut = false
      localStorage.removeItem('auth_token')
      localStorage.removeItem('auth_user')
      try {
        const pStore = useProjectStore()
        pStore.resetStore()
      } catch (_) {}
    },
  },
})

// Cross-tab logout listener & idle inactivity tracker
export function setupAuthCrossTabSync(authStore: ReturnType<typeof useAuthStore>, router: any) {
  if (typeof window === 'undefined') return

  // 1. BroadcastChannel API for real-time tab synchronization
  if (typeof BroadcastChannel !== 'undefined') {
    const bc = new BroadcastChannel('pulse_auth_sync')
    bc.onmessage = (event) => {
      if (event.data?.type === 'LOGOUT') {
        console.info('[Auth] Logout message received from another tab. Clearing session...')
        authStore.forceLogout()
        if (router.currentRoute.value.path !== '/login') {
          router.replace({ path: '/login', query: { reason: 'session_expired' } })
        }
      } else if (event.data?.type === 'NEW_LOGIN') {
        if (event.data?.userId && event.data?.userId !== authStore.user?.id) {
          console.info('[Auth] Different user logged in on another tab. Synchronizing...')
          authStore.forceLogout()
          if (router.currentRoute.value.path !== '/login') {
            router.replace({ path: '/login', query: { reason: 'concurrent_device' } })
          }
        }
      }
    }
  }

  // 2. Fallback Storage Event (when auth_token removed in another tab)
  window.addEventListener('storage', (e) => {
    if (e.key === 'auth_token' && (!e.newValue || e.newValue === '')) {
      console.info('[Auth] Auth token removed in another tab. Syncing logout state...')
      authStore.forceLogout()
      if (router.currentRoute.value.path !== '/login') {
        router.replace({ path: '/login', query: { reason: 'session_expired' } })
      }
    }
  })

  // 3. Inactivity Auto-Logout Tracker (15 minutes of zero user interaction)
  const IDLE_TIMEOUT_MS = 15 * 60 * 1000 // 15 minutes (900 seconds)
  let idleTimer: any = null

  function resetIdleTimer() {
    if (!authStore.isAuthenticated) return
    if (idleTimer) clearTimeout(idleTimer)
    idleTimer = setTimeout(async () => {
      if (authStore.isAuthenticated) {
        console.warn('[Auth] Inactivity threshold reached (15m). Automatically logging out...')
        await authStore.logout()
        if (router.currentRoute.value.path !== '/login') {
          router.replace({ path: '/login', query: { reason: 'idle_timeout' } })
        }
      }
    }, IDLE_TIMEOUT_MS)
  }

  const activityEvents = ['mousemove', 'mousedown', 'keydown', 'touchstart', 'scroll', 'click']
  activityEvents.forEach((evt) => {
    window.addEventListener(evt, resetIdleTimer, { passive: true })
  })

  // Initial trigger
  resetIdleTimer()

  // 4. Background Session Heartbeat (Every 20 seconds)
  // Verifies that this session is still the only active session and hasn't been superseded on another device
  setInterval(async () => {
    if (!authStore.isAuthenticated || authStore.isLoggingOut) return
    try {
      await apiClient.get('/auth/heartbeat')
    } catch (err: any) {
      if (err?.response?.status === 401) {
        const code = err.response?.data?.code || ''
        authStore.forceLogout()
        let reason = 'session_expired'
        if (code === 'SESSION_CONCURRENT_KICKED') {
          reason = 'concurrent_device'
        } else if (code === 'SESSION_IDLE_TIMEOUT') {
          reason = 'idle_timeout'
        }
        if (router.currentRoute.value.path !== '/login') {
          router.replace({ path: '/login', query: { reason } })
        }
      }
    }
  }, 20 * 1000)
}
