import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import { useAuthStore } from '../stores/authStore'
import DashboardView from '../views/DashboardView.vue'
import ProjectWorkspaceView from '../views/ProjectWorkspaceView.vue'
import ClientPortalView from '../views/ClientPortalView.vue'
import UnauthorizedView from '../views/UnauthorizedView.vue'
import LoginView from '../views/LoginView.vue'
import AcceptInviteView from '../views/AcceptInviteView.vue'

const routes: RouteRecordRaw[] = [
  // 1. Public Entry Point: Root directly opens /login first
  {
    path: '/',
    redirect: '/login',
  },
  {
    path: '/login',
    name: 'Login',
    component: LoginView,
    meta: { layout: 'public', requiresAuth: false },
  },
  {
    path: '/accept-invite',
    name: 'AcceptInvite',
    component: AcceptInviteView,
    meta: { layout: 'public', requiresAuth: false },
  },
  {
    path: '/unauthorized',
    name: 'Unauthorized',
    component: UnauthorizedView,
    meta: { layout: 'protected', requiresAuth: true },
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: DashboardView,
    meta: { layout: 'protected', requiresAuth: true, permission: 'project:read' },
  },
  {
    path: '/projects',
    name: 'Projects',
    component: ProjectWorkspaceView,
    meta: { layout: 'protected', requiresAuth: true, permission: 'project:read' },
  },
  {
    path: '/projects/:id',
    name: 'ProjectDetail',
    component: ProjectWorkspaceView,
    meta: { layout: 'protected', requiresAuth: true, permission: 'project:read' },
  },
  {
    path: '/client-portal',
    name: 'ClientPortal',
    component: ClientPortalView,
    meta: { layout: 'protected', requiresAuth: true, permission: 'project:read' },
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/login',
  },
]

export const router = createRouter({
  history: createWebHistory(),
  routes,
})

// Navigation Guard enforcing Login-First and RBAC
router.beforeEach(async (to, _from, next) => {
  const authStore = useAuthStore()

  // Rule 1: Opening the site root '/' strictly directs to '/login' first
  if (to.path === '/') {
    return next({ path: '/login' })
  }

  // Rule 2: User explicitly visiting public pages always proceeds
  const cleanPath = to.path.replace(/\/+$/, '')
  if (cleanPath === '/login' || cleanPath === '/accept-invite' || to.name === 'AcceptInvite' || to.name === 'Login') {
    return next()
  }

  // Rule 3: Protected routes require valid authenticated session
  const requiresAuth = to.meta.requiresAuth !== false
  if (requiresAuth && !authStore.isAuthenticated) {
    if (to.path === '/dashboard') {
      return next({ path: '/login' })
    }
    return next({ path: '/login', query: { redirect: to.fullPath } })
  }

  // Rule 4: Initialize user profile & permissions if not loaded
  if (authStore.isAuthenticated && authStore.permissions.length === 0) {
    try {
      await authStore.fetchUser()
    } catch {
      authStore.forceLogout()
      return next({ path: '/login' })
    }
  }

  // Rule 5: Check required global roles
  if (to.meta.roles) {
    const allowedRoles = to.meta.roles as string[]
    if (!authStore.hasRole(allowedRoles)) {
      return next({ path: '/unauthorized' })
    }
  }

  // Rule 6: Check required granular permission
  if (to.meta.permission) {
    const requiredPerm = to.meta.permission as string
    if (!authStore.can(requiredPerm)) {
      return next({ path: '/unauthorized' })
    }
  }

  next()
})
