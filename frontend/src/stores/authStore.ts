import { defineStore } from 'pinia'
import type { User } from '../types/project'
import { apiClient, USE_MOCK } from '../services/api'

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

// Pre-seeded User profiles for role switching
export const SEED_USERS: Record<string, User> = {
  super_admin: {
    id: 'usr_super_admin',
    name: 'System Super Admin',
    email: 'admin@pulse.workspace',
    role: 'super_admin' as any,
    avatar_url: 'https://images.unsplash.com/photo-1534528741775-53994a69daeb?w=100&auto=format&fit=crop&q=80',
    initials: 'SA',
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
  client: {
    id: 'usr_client',
    name: 'Ir. H. Gunawan (Client)',
    email: 'client@pulse.workspace',
    role: 'client' as any,
    avatar_url: 'https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=100&auto=format&fit=crop&q=80',
    initials: 'HG',
  },
}

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: (SEED_USERS['pm'] as User) as User | null,
    token: localStorage.getItem('auth_token') || 'mock-jwt-token-pm',
    permissions: ROLE_PERMISSIONS['pm'] || [],
    isLoading: false,
    useMock: USE_MOCK,
  }),

  getters: {
    isAuthenticated(state): boolean {
      return !!state.token && !!state.user
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
      // super_admin has all permissions
      if (state.user.role === 'super_admin') return true
      return state.permissions.includes(permission)
    },
  },

  actions: {
    async fetchUser() {
      if (this.useMock) {
        const role = this.user?.role || 'pm'
        this.permissions = ROLE_PERMISSIONS[role] || []
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
        this.permissions = data.permissions || []
      } catch (err) {
        console.error('Failed to fetch user permissions from backend:', err)
        // Fallback to role-based default permissions
        if (this.user?.role && ROLE_PERMISSIONS[this.user.role]) {
          this.permissions = ROLE_PERMISSIONS[this.user.role]
        }
      } finally {
        this.isLoading = false
      }
    },

    async login(email: string, password: string): Promise<boolean> {
      if (this.useMock) {
        // Find matching seed user by email or default to PM
        const matched = Object.values(SEED_USERS).find((u) => u.email === email) || SEED_USERS['pm']
        this.user = matched
        this.token = `mock-jwt-token-${matched.role}`
        localStorage.setItem('auth_token', this.token)
        this.permissions = ROLE_PERMISSIONS[matched.role] || []
        return true
      }

      try {
        this.isLoading = true
        const response = await apiClient.post('/auth/login', { email, password })
        const { token, user, permissions } = response.data
        this.token = token
        localStorage.setItem('auth_token', token)
        this.user = user
        this.permissions = permissions || []
        return true
      } catch (err) {
        console.error('Login failed:', err)
        throw err
      } finally {
        this.isLoading = false
      }
    },

    logout() {
      this.token = ''
      this.user = null
      this.permissions = []
      localStorage.removeItem('auth_token')
    },

    // Dev Role Switcher action (for testing instant role transitions)
    switchDevRole(role: 'super_admin' | 'pm' | 'member' | 'client') {
      const selected = SEED_USERS[role]
      if (selected) {
        this.user = selected
        this.token = `mock-jwt-token-${role}`
        localStorage.setItem('auth_token', this.token)
        this.permissions = ROLE_PERMISSIONS[role] || []
      }
    },
  },
})
