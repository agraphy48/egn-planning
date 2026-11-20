import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import { useAuthStore } from '../stores/authStore'
import DashboardView from '../views/DashboardView.vue'
import ProjectWorkspaceView from '../views/ProjectWorkspaceView.vue'
import ClientPortalView from '../views/ClientPortalView.vue'
import UnauthorizedView from '../views/UnauthorizedView.vue'
import LoginView from '../views/LoginView.vue'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: LoginView,
    meta: { requiresAuth: false },
  },
  {
    path: '/unauthorized',
    name: 'Unauthorized',
    component: UnauthorizedView,
    meta: { requiresAuth: true },
  },
  {
    path: '/',
    redirect: '/dashboard',
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: DashboardView,
    meta: { requiresAuth: true, permission: 'project:read' },
  },
  {
    path: '/projects',
    name: 'Projects',
    component: ProjectWorkspaceView,
    meta: { requiresAuth: true, permission: 'project:read' },
  },
  {
    path: '/projects/:id',
    name: 'ProjectDetail',
    component: ProjectWorkspaceView,
    meta: { requiresAuth: true, permission: 'project:read' },
  },
  {
    path: '/client-portal',
    name: 'ClientPortal',
    component: ClientPortalView,
    meta: { requiresAuth: true, permission: 'project:read' },
  },
]

export const router = createRouter({
  history: createWebHistory(),
  routes,
})

// Navigation Guard enforcing Auth & RBAC
router.beforeEach(async (to, _from, next) => {
  const authStore = useAuthStore()

  // Initialize permissions if not yet loaded
  if (authStore.isAuthenticated && authStore.permissions.length === 0) {
    await authStore.fetchUser()
  }

  const requiresAuth = to.meta.requiresAuth !== false

  if (requiresAuth && !authStore.isAuthenticated) {
    return next({ path: '/login', query: { redirect: to.fullPath } })
  }

  // Check required global roles
  if (to.meta.roles) {
    const allowedRoles = to.meta.roles as string[]
    if (!authStore.hasRole(allowedRoles)) {
      return next({ path: '/unauthorized' })
    }
  }

  // Check required granular permission
  if (to.meta.permission) {
    const requiredPerm = to.meta.permission as string
    if (!authStore.can(requiredPerm)) {
      return next({ path: '/unauthorized' })
    }
  }

  next()
})
