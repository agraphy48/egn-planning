<template>
  <aside
    :class="[
      'bg-white dark:bg-[#18181B] border-r border-[#E4E4E7] dark:border-[#27272A] flex flex-col justify-between flex-shrink-0 z-30 transition-all duration-200 ease-in-out select-none',
      store.isSidebarCollapsed ? 'w-16' : 'w-60'
    ]"
  >
    <!-- Top Workspace Section -->
    <div class="flex flex-col flex-1 min-h-0">
      <!-- Workspace Switcher -->
      <div class="p-2.5 border-b border-[#E4E4E7] dark:border-[#27272A]">
        <div
          v-if="!store.isSidebarCollapsed"
          class="flex items-center justify-between p-1.5 rounded-lg hover:bg-[#F4F4F5] dark:hover:bg-[#27272A] transition-colors group cursor-pointer"
        >
          <div class="flex items-center space-x-2.5 min-w-0">
            <div class="w-6 h-6 rounded-md bg-brand-600 flex items-center justify-center text-white font-display font-bold text-xs shadow-2xs">
              P
            </div>
            <div class="truncate">
              <div class="text-xs font-semibold tracking-tight text-zinc-900 dark:text-zinc-100 flex items-center gap-1.5">
                <span class="truncate">Pulse HQ</span>
                <span class="text-[9px] font-mono px-1 py-0.2 rounded bg-brand-50 dark:bg-brand-950/60 text-brand-600 dark:text-brand-400 font-medium border border-brand-200/50 dark:border-brand-800/40">PRO</span>
              </div>
            </div>
          </div>
          <ChevronsUpDown class="w-3.5 h-3.5 text-zinc-400 group-hover:text-zinc-600 dark:group-hover:text-zinc-300 flex-shrink-0" />
        </div>

        <!-- Collapsed Workspace Icon -->
        <div v-else class="flex justify-center py-0.5">
          <div
            @click="store.toggleSidebar()"
            class="w-8 h-8 rounded-lg bg-brand-600 flex items-center justify-center text-white font-display font-bold text-xs shadow-2xs cursor-pointer hover:bg-brand-500 transition-colors"
            title="Pulse HQ Workspace (Click to expand)"
          >
            P
          </div>
        </div>

        <!-- Quick Search Bar Trigger (Linear Cmd+K) -->
        <div class="mt-2">
          <button
            v-if="!store.isSidebarCollapsed"
            @click="store.isCommandPaletteOpen = true"
            class="w-full flex items-center justify-between px-2.5 py-1.5 rounded-lg text-xs text-zinc-500 dark:text-zinc-400 bg-[#F4F4F5] dark:bg-[#27272A] hover:bg-zinc-200/70 dark:hover:bg-zinc-700/60 transition-colors"
          >
            <div class="flex items-center space-x-2 truncate">
              <Search class="w-3.5 h-3.5 text-zinc-400" />
              <span class="text-[11px]">Search or jump...</span>
            </div>
            <kbd class="font-mono text-[9px] px-1 py-0.5 bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-700 rounded text-zinc-500">⌘K</kbd>
          </button>
          <button
            v-else
            @click="store.isCommandPaletteOpen = true"
            class="w-full flex justify-center p-2 rounded-lg text-zinc-500 hover:bg-[#F4F4F5] dark:hover:bg-[#27272A] transition-colors"
            title="Search (⌘K)"
          >
            <Search class="w-4 h-4 text-zinc-400" />
          </button>
        </div>
      </div>

      <!-- Navigation Items -->
      <div class="flex-1 overflow-y-auto px-2 py-3 space-y-4">
        <!-- Main Nav -->
        <div>
          <div
            v-if="!store.isSidebarCollapsed"
            class="px-2 mb-1.5 text-[10px] font-bold uppercase tracking-wider text-zinc-400 dark:text-zinc-500 font-mono"
          >
            Core Views
          </div>
          <nav class="space-y-0.5 text-xs font-medium">
            <!-- Inbox / Notifications -->
            <button
              @click="store.showToast('You have 2 unread project updates')"
              :class="[
                'w-full flex items-center rounded-md transition-colors text-zinc-700 dark:text-zinc-300 hover:bg-[#F4F4F5] dark:hover:bg-[#27272A] group',
                store.isSidebarCollapsed ? 'justify-center p-2' : 'justify-between px-2.5 py-1.5'
              ]"
              :title="store.isSidebarCollapsed ? 'Inbox (2 unread)' : ''"
            >
              <div class="flex items-center space-x-2.5">
                <Inbox class="w-4 h-4 text-zinc-500 group-hover:text-brand-600 dark:group-hover:text-brand-400" />
                <span v-if="!store.isSidebarCollapsed">Inbox</span>
              </div>
              <span
                v-if="!store.isSidebarCollapsed"
                class="text-[10px] font-mono px-1.5 py-0.2 rounded-full bg-brand-50 dark:bg-brand-950 text-brand-600 dark:text-brand-400 font-semibold"
              >
                2
              </span>
            </button>

            <!-- Dashboard -->
            <RouterLink
              to="/dashboard"
              v-slot="{ isActive }"
              class="block"
            >
              <div
                :class="[
                  'w-full flex items-center rounded-md transition-colors group',
                  store.isSidebarCollapsed ? 'justify-center p-2' : 'justify-between px-2.5 py-1.5',
                  isActive
                    ? 'text-brand-600 dark:text-brand-400 bg-brand-50/70 dark:bg-brand-950/40 font-semibold'
                    : 'text-zinc-700 dark:text-zinc-300 hover:bg-[#F4F4F5] dark:hover:bg-[#27272A]'
                ]"
                :title="store.isSidebarCollapsed ? 'Executive Dashboard' : ''"
              >
                <div class="flex items-center space-x-2.5">
                  <LayoutDashboard class="w-4 h-4 text-zinc-500 group-hover:text-brand-600 dark:group-hover:text-brand-400" />
                  <span v-if="!store.isSidebarCollapsed">Dashboard</span>
                </div>
              </div>
            </RouterLink>

            <!-- Projects Workspace -->
            <RouterLink
              to="/projects"
              v-slot="{ isActive }"
              class="block"
            >
              <div
                :class="[
                  'w-full flex items-center rounded-md transition-colors group',
                  store.isSidebarCollapsed ? 'justify-center p-2' : 'justify-between px-2.5 py-1.5',
                  isActive
                    ? 'text-brand-600 dark:text-brand-400 bg-brand-50/70 dark:bg-brand-950/40 font-semibold'
                    : 'text-zinc-700 dark:text-zinc-300 hover:bg-[#F4F4F5] dark:hover:bg-[#27272A]'
                ]"
                :title="store.isSidebarCollapsed ? 'Project Workspace' : ''"
              >
                <div class="flex items-center space-x-2.5">
                  <FolderKanban class="w-4 h-4 text-zinc-500 group-hover:text-brand-600 dark:group-hover:text-brand-400" />
                  <span v-if="!store.isSidebarCollapsed">Projects</span>
                </div>
                <span
                  v-if="!store.isSidebarCollapsed"
                  class="text-[10px] font-mono px-1.5 py-0.2 rounded-full bg-zinc-100 dark:bg-zinc-800 text-zinc-500 tabular-nums"
                >
                  {{ store.tasks.length }}
                </span>
              </div>
            </RouterLink>

            <!-- Client Portal -->
            <RouterLink
              to="/client-portal"
              v-slot="{ isActive }"
              class="block"
            >
              <div
                :class="[
                  'w-full flex items-center rounded-md transition-colors group',
                  store.isSidebarCollapsed ? 'justify-center p-2' : 'justify-between px-2.5 py-1.5',
                  isActive
                    ? 'text-brand-600 dark:text-brand-400 bg-brand-50/70 dark:bg-brand-950/40 font-semibold'
                    : 'text-zinc-700 dark:text-zinc-300 hover:bg-[#F4F4F5] dark:hover:bg-[#27272A]'
                ]"
                :title="store.isSidebarCollapsed ? 'Client External Portal' : ''"
              >
                <div class="flex items-center space-x-2.5">
                  <ShieldCheck class="w-4 h-4 text-indigo-500 group-hover:text-indigo-600" />
                  <span v-if="!store.isSidebarCollapsed">Client Portal</span>
                </div>
                <span
                  v-if="!store.isSidebarCollapsed"
                  class="text-[9px] font-mono px-1.5 py-0.2 rounded bg-indigo-50 dark:bg-indigo-950/60 text-indigo-600 dark:text-indigo-300 border border-indigo-200/50 dark:border-indigo-800/40"
                >
                  Portal
                </span>
              </div>
            </RouterLink>
          </nav>
        </div>

        <!-- Pinned Active Projects List -->
        <div>
          <div
            v-if="!store.isSidebarCollapsed"
            class="px-2 mb-1.5 flex items-center justify-between text-[10px] font-bold uppercase tracking-wider text-zinc-400 dark:text-zinc-500 font-mono"
          >
            <span>Projects</span>
            <button
              v-can="'task:create'"
              @click="store.isNewTaskModalOpen = true"
              class="hover:text-zinc-700 dark:hover:text-zinc-300 transition-colors"
              title="Add task or project"
            >
              <Plus class="w-3.5 h-3.5" />
            </button>
          </div>
          <div class="space-y-0.5 text-xs">
            <button
              v-for="project in store.projects"
              :key="project.id"
              @click="selectProject(project.id)"
              :class="[
                'w-full flex items-center rounded-md transition-colors group text-left',
                store.isSidebarCollapsed ? 'justify-center p-2' : 'justify-between px-2.5 py-1.5',
                store.currentProjectId === project.id
                  ? 'bg-zinc-100 dark:bg-zinc-800 text-zinc-900 dark:text-white font-medium'
                  : 'text-zinc-600 dark:text-zinc-400 hover:bg-zinc-100/60 dark:hover:bg-zinc-800/50 hover:text-zinc-900 dark:hover:text-white'
              ]"
              :title="store.isSidebarCollapsed ? project.name : ''"
            >
              <div class="flex items-center space-x-2 truncate">
                <span
                  class="w-2 h-2 rounded-full flex-shrink-0"
                  :style="{ backgroundColor: project.color_dot }"
                ></span>
                <span v-if="!store.isSidebarCollapsed" class="truncate">{{ project.name }}</span>
              </div>
              <span
                v-if="!store.isSidebarCollapsed"
                class="text-[10px] font-mono text-zinc-400 group-hover:text-zinc-600 dark:group-hover:text-zinc-300 tabular-nums"
              >
                {{ project.progress_percentage }}%
              </span>
            </button>
          </div>
        </div>

        <!-- Administration (Super Admin / PM only) -->
        <div v-if="canManageTeamOrSettings">
          <div
            v-if="!store.isSidebarCollapsed"
            class="px-2 mb-1.5 text-[10px] font-bold uppercase tracking-wider text-zinc-400 dark:text-zinc-500 font-mono"
          >
            Admin
          </div>
          <nav class="space-y-0.5 text-xs font-medium">
            <button
              @click="store.showToast('Team & Access Permissions (RBAC: ' + authStore.currentRole.toUpperCase() + ')')"
              :class="[
                'w-full flex items-center rounded-md transition-colors text-zinc-700 dark:text-zinc-300 hover:bg-[#F4F4F5] dark:hover:bg-[#27272A] group',
                store.isSidebarCollapsed ? 'justify-center p-2' : 'justify-between px-2.5 py-1.5'
              ]"
              :title="store.isSidebarCollapsed ? 'Team & Access' : ''"
            >
              <div class="flex items-center space-x-2.5">
                <Users class="w-4 h-4 text-zinc-500 group-hover:text-brand-600 dark:group-hover:text-brand-400" />
                <span v-if="!store.isSidebarCollapsed">Team & Access</span>
              </div>
            </button>

            <button
              @click="store.showToast('Workspace Settings')"
              :class="[
                'w-full flex items-center rounded-md transition-colors text-zinc-700 dark:text-zinc-300 hover:bg-[#F4F4F5] dark:hover:bg-[#27272A] group',
                store.isSidebarCollapsed ? 'justify-center p-2' : 'justify-between px-2.5 py-1.5'
              ]"
              :title="store.isSidebarCollapsed ? 'Settings' : ''"
            >
              <div class="flex items-center space-x-2.5">
                <Settings class="w-4 h-4 text-zinc-500 group-hover:text-brand-600 dark:group-hover:text-brand-400" />
                <span v-if="!store.isSidebarCollapsed">Settings</span>
              </div>
            </button>
          </nav>
        </div>
      </div>
    </div>

    <!-- Sidebar Footer (Dev Role Switcher, Theme & Collapse Toggle) -->
    <div class="p-2 border-t border-[#E4E4E7] dark:border-[#27272A] space-y-1.5">
      <!-- Dev Role Switcher -->
      <div v-if="authStore.useMock && !store.isSidebarCollapsed" class="p-1.5 rounded-lg bg-[#F4F4F5] dark:bg-[#27272A] space-y-1">
        <div class="flex items-center justify-between text-[9px] font-mono text-zinc-400">
          <span>ROLE</span>
          <span class="text-brand-600 dark:text-brand-400 font-bold uppercase">{{ authStore.currentRole }}</span>
        </div>
        <div class="grid grid-cols-4 gap-0.5 text-[9px] font-medium text-center">
          <button
            @click="switchRole('super_admin')"
            :class="authStore.currentRole === 'super_admin' ? 'bg-white dark:bg-zinc-900 text-brand-600 dark:text-brand-400 font-bold shadow-2xs' : 'text-zinc-500 hover:text-zinc-800 dark:hover:text-zinc-200'"
            class="py-0.5 rounded transition-colors"
          >
            Admin
          </button>
          <button
            @click="switchRole('pm')"
            :class="authStore.currentRole === 'pm' ? 'bg-white dark:bg-zinc-900 text-brand-600 dark:text-brand-400 font-bold shadow-2xs' : 'text-zinc-500 hover:text-zinc-800 dark:hover:text-zinc-200'"
            class="py-0.5 rounded transition-colors"
          >
            PM
          </button>
          <button
            @click="switchRole('member')"
            :class="authStore.currentRole === 'member' ? 'bg-white dark:bg-zinc-900 text-brand-600 dark:text-brand-400 font-bold shadow-2xs' : 'text-zinc-500 hover:text-zinc-800 dark:hover:text-zinc-200'"
            class="py-0.5 rounded transition-colors"
          >
            Dev
          </button>
          <button
            @click="switchRole('client')"
            :class="authStore.currentRole === 'client' ? 'bg-white dark:bg-zinc-900 text-brand-600 dark:text-brand-400 font-bold shadow-2xs' : 'text-zinc-500 hover:text-zinc-800 dark:hover:text-zinc-200'"
            class="py-0.5 rounded transition-colors"
          >
            Client
          </button>
        </div>
      </div>

      <!-- Theme Toggle & Collapse Button Row -->
      <div class="flex items-center justify-between gap-1">
        <button
          @click="store.toggleTheme()"
          class="flex-1 flex items-center justify-center p-1.5 rounded-md text-xs font-medium text-zinc-600 dark:text-zinc-400 hover:bg-[#F4F4F5] dark:hover:bg-[#27272A] transition-colors"
          :title="store.isDarkMode ? 'Switch to Light Mode' : 'Switch to Dark Mode'"
        >
          <Sun v-if="store.isDarkMode" class="w-4 h-4 text-amber-500" />
          <Moon v-else class="w-4 h-4 text-zinc-500" />
          <span v-if="!store.isSidebarCollapsed" class="ml-2 text-xs">
            {{ store.isDarkMode ? 'Light' : 'Dark' }}
          </span>
        </button>

        <button
          @click="store.toggleSidebar()"
          class="p-1.5 rounded-md text-zinc-400 hover:text-zinc-700 dark:hover:text-zinc-200 hover:bg-[#F4F4F5] dark:hover:bg-[#27272A] transition-colors"
          :title="store.isSidebarCollapsed ? 'Expand Sidebar' : 'Collapse Sidebar'"
        >
          <PanelLeftClose v-if="!store.isSidebarCollapsed" class="w-4 h-4" />
          <PanelLeft v-else class="w-4 h-4" />
        </button>
      </div>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useProjectStore } from '../stores/projectStore'
import { useAuthStore } from '../stores/authStore'
import { useRouter } from 'vue-router'
import {
  ChevronsUpDown,
  Search,
  Inbox,
  LayoutDashboard,
  FolderKanban,
  ShieldCheck,
  Plus,
  Users,
  Settings,
  Sun,
  Moon,
  PanelLeftClose,
  PanelLeft
} from 'lucide-vue-next'

const store = useProjectStore()
const authStore = useAuthStore()
const router = useRouter()

const canManageTeamOrSettings = computed(() => {
  return authStore.hasRole(['super_admin', 'pm']) || authStore.can('user:manage') || authStore.can('project:member:manage')
})

function switchRole(role: 'super_admin' | 'pm' | 'member' | 'client') {
  authStore.switchDevRole(role)
  store.showToast(`Active role switched to: ${role.toUpperCase()}`)
}

function selectProject(projectId: string) {
  store.selectProject(projectId)
  router.push('/projects')
}
</script>
