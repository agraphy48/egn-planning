<template>
  <header class="h-[52px] border-b border-[#E4E4E7] dark:border-[#27272A] bg-white/90 dark:bg-[#18181B]/90 backdrop-blur-md px-4 sm:px-5 flex items-center justify-between flex-shrink-0 z-20 select-none">
    <!-- Left: Breadcrumbs -->
    <div class="flex items-center space-x-3 min-w-0">
      <button
        @click="store.toggleSidebar()"
        class="p-1 rounded-md text-zinc-400 hover:text-zinc-700 dark:hover:text-zinc-200 hover:bg-[#F4F4F5] dark:hover:bg-[#27272A] transition-colors"
        title="Toggle Sidebar"
      >
        <PanelLeft v-if="store.isSidebarCollapsed" class="w-4 h-4" />
        <PanelLeftClose v-else class="w-4 h-4" />
      </button>

      <!-- Single-row clean breadcrumbs -->
      <nav aria-label="Breadcrumb" class="flex items-center space-x-1.5 text-xs text-zinc-500 dark:text-zinc-400 truncate">
        <span class="font-medium text-zinc-700 dark:text-zinc-300">Pulse</span>
        <ChevronRight class="w-3.5 h-3.5 text-zinc-400 flex-shrink-0" />
        <span class="text-zinc-500">Projects</span>
        <ChevronRight class="w-3.5 h-3.5 text-zinc-400 flex-shrink-0" />
        <span class="font-semibold text-zinc-900 dark:text-zinc-100 truncate">
          {{ store.currentProject?.name || 'Workspace' }}
        </span>
        <span class="hidden lg:inline-flex items-center px-1.5 py-0.2 rounded text-[10px] font-mono bg-[#F4F4F5] dark:bg-[#27272A] text-zinc-500 tabular-nums ml-1">
          Sprint 38
        </span>
      </nav>
    </div>

    <!-- Center: Compact View Switcher pills (Board | List | Timeline) -->
    <div class="hidden md:flex items-center p-0.5 rounded-lg bg-[#F4F4F5] dark:bg-[#27272A] border border-[#E4E4E7]/70 dark:border-[#27272A] text-xs font-medium">
      <button
        @click="handleViewChange('kanban')"
        :class="[
          'px-2.5 py-1 rounded-md transition-all flex items-center gap-1.5 text-[11px]',
          store.currentView === 'kanban'
            ? 'bg-white dark:bg-[#18181B] text-zinc-900 dark:text-white shadow-2xs font-semibold'
            : 'text-zinc-500 hover:text-zinc-900 dark:hover:text-zinc-200'
        ]"
      >
        <Kanban class="w-3.5 h-3.5" :class="store.currentView === 'kanban' ? 'text-brand-600 dark:text-brand-400' : ''" />
        <span>Board</span>
      </button>

      <button
        @click="handleViewChange('list')"
        :class="[
          'px-2.5 py-1 rounded-md transition-all flex items-center gap-1.5 text-[11px]',
          store.currentView === 'list'
            ? 'bg-white dark:bg-[#18181B] text-zinc-900 dark:text-white shadow-2xs font-semibold'
            : 'text-zinc-500 hover:text-zinc-900 dark:hover:text-zinc-200'
        ]"
      >
        <ListTodo class="w-3.5 h-3.5" :class="store.currentView === 'list' ? 'text-brand-600 dark:text-brand-400' : ''" />
        <span>List</span>
      </button>

      <button
        @click="handleViewChange('timeline')"
        :class="[
          'px-2.5 py-1 rounded-md transition-all flex items-center gap-1.5 text-[11px]',
          store.currentView === 'timeline'
            ? 'bg-white dark:bg-[#18181B] text-zinc-900 dark:text-white shadow-2xs font-semibold'
            : 'text-zinc-500 hover:text-zinc-900 dark:hover:text-zinc-200'
        ]"
      >
        <GitCommit class="w-3.5 h-3.5" :class="store.currentView === 'timeline' ? 'text-brand-600 dark:text-brand-400' : ''" />
        <span>Timeline</span>
      </button>
    </div>

    <!-- Right: Filter Popover, Cmd+K trigger, Role Switcher, + New Task, Notifications, Avatar -->
    <div class="flex items-center space-x-2.5 flex-shrink-0">
      <!-- Filter Popover Button -->
      <div class="relative">
        <button
          @click="showFilterPopover = !showFilterPopover"
          :class="[
            'flex items-center space-x-1.5 px-2.5 py-1 rounded-lg text-xs font-medium border transition-colors',
            store.activeFilter !== 'All'
              ? 'bg-brand-50 dark:bg-brand-950/50 text-brand-600 dark:text-brand-400 border-brand-200 dark:border-brand-800'
              : 'bg-[#F4F4F5] dark:bg-[#27272A] border-transparent text-zinc-600 dark:text-zinc-400 hover:text-zinc-900 dark:hover:text-zinc-100'
          ]"
        >
          <SlidersHorizontal class="w-3.5 h-3.5" />
          <span class="hidden sm:inline text-[11px]">{{ store.activeFilter === 'All' ? 'Filter' : store.activeFilter }}</span>
        </button>

        <!-- Filter Dropdown Menu -->
        <div
          v-if="showFilterPopover"
          @click.stop
          class="absolute right-0 mt-1.5 w-48 bg-white dark:bg-[#18181B] border border-[#E4E4E7] dark:border-[#27272A] rounded-xl shadow-lg z-50 p-1.5 space-y-0.5 text-xs animate-in fade-in zoom-in-95 duration-100"
        >
          <div class="px-2 py-1 text-[10px] font-mono uppercase text-zinc-400 font-semibold">Filter by Category</div>
          <button
            v-for="cat in categories"
            :key="cat"
            @click="setCategoryFilter(cat)"
            :class="[
              'w-full flex items-center justify-between px-2 py-1.5 rounded-lg text-left transition-colors',
              store.activeFilter === cat
                ? 'bg-brand-50 dark:bg-brand-950/60 text-brand-600 dark:text-brand-400 font-semibold'
                : 'text-zinc-700 dark:text-zinc-300 hover:bg-[#F4F4F5] dark:hover:bg-[#27272A]'
            ]"
          >
            <span>{{ cat }}</span>
            <Check v-if="store.activeFilter === cat" class="w-3.5 h-3.5 text-brand-600 dark:text-brand-400" />
          </button>
        </div>
      </div>

      <!-- Global Search trigger button (Linear Cmd+K) -->
      <button
        @click="store.isCommandPaletteOpen = true"
        class="hidden sm:flex items-center space-x-2 text-xs text-zinc-400 bg-[#F4F4F5] dark:bg-[#27272A] hover:bg-zinc-200/60 dark:hover:bg-zinc-700/60 px-2.5 py-1 rounded-lg transition-colors group"
      >
        <Search class="w-3.5 h-3.5 text-zinc-400 group-hover:text-zinc-600 dark:group-hover:text-zinc-300" />
        <span class="hidden xl:inline text-[11px]">Search...</span>
        <kbd class="font-mono text-[9px] px-1 py-0.2 bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-700 rounded text-zinc-500">⌘K</kbd>
      </button>

      <!-- Dev Role Switcher Dropdown -->
      <div
        v-if="authStore.useMock"
        class="flex items-center space-x-1 bg-[#F4F4F5] dark:bg-[#27272A] px-2 py-1 rounded-lg border border-[#E4E4E7] dark:border-[#27272A]"
        title="Switch RBAC role to test permissions"
      >
        <ShieldCheck class="w-3.5 h-3.5 text-brand-600 dark:text-brand-400 flex-shrink-0" />
        <select
          :value="authStore.currentRole"
          @change="onRoleChange(($event.target as HTMLSelectElement).value)"
          class="bg-transparent text-[10px] font-bold text-brand-600 dark:text-brand-400 focus:outline-none cursor-pointer uppercase font-mono"
        >
          <option value="super_admin">Super Admin</option>
          <option value="pm">PM</option>
          <option value="member">Member</option>
          <option value="client">Client</option>
        </select>
      </div>

      <!-- Quick "+ New Task" Button -->
      <button
        v-can="'task:create'"
        @click="store.isNewTaskModalOpen = true"
        class="flex items-center space-x-1.5 px-3 py-1 rounded-lg bg-brand-600 hover:bg-brand-500 text-white text-xs font-semibold shadow-2xs transition-all active:scale-[0.98]"
      >
        <Plus class="w-3.5 h-3.5 stroke-[2.5]" />
        <span class="hidden sm:inline text-[11px]">New Task</span>
      </button>

      <!-- Notification Bell -->
      <div class="relative">
        <button
          @click="showNotifications = !showNotifications"
          class="p-1.5 rounded-lg text-zinc-400 hover:text-zinc-700 dark:hover:text-zinc-200 hover:bg-[#F4F4F5] dark:hover:bg-[#27272A] transition-colors relative"
          title="Notifications"
        >
          <Bell class="w-4 h-4" />
          <span class="absolute top-1 right-1 w-2 h-2 rounded-full bg-rose-500 ring-2 ring-white dark:ring-zinc-900"></span>
        </button>

        <!-- Dropdown Notification Popover -->
        <div
          v-if="showNotifications"
          class="absolute right-0 mt-2 w-72 bg-white dark:bg-[#18181B] border border-[#E4E4E7] dark:border-[#27272A] rounded-xl shadow-xl z-50 p-3 space-y-2 animate-in fade-in zoom-in-95 duration-100"
        >
          <div class="flex items-center justify-between pb-2 border-b border-[#E4E4E7] dark:border-[#27272A]">
            <span class="text-xs font-bold text-zinc-900 dark:text-white">Notifications</span>
            <button @click="showNotifications = false" class="text-[10px] text-brand-600 dark:text-brand-400 hover:underline">Mark all read</button>
          </div>
          <div class="space-y-2 text-xs">
            <div class="flex items-start space-x-2.5 p-2 rounded-lg bg-brand-50/50 dark:bg-brand-950/20">
              <span class="w-2 h-2 rounded-full bg-brand-600 mt-1 flex-shrink-0"></span>
              <div>
                <p class="text-zinc-800 dark:text-zinc-200 text-[11px]">
                  <span class="font-semibold">Sarah Connor</span> updated status on
                  <span class="text-brand-600 dark:text-brand-400 font-medium">Stripe Sandbox QA</span>
                </p>
                <span class="text-[10px] text-zinc-400">10m ago</span>
              </div>
            </div>
            <div class="flex items-start space-x-2.5 p-2 rounded-lg hover:bg-zinc-50 dark:hover:bg-zinc-800/40 transition-colors">
              <span class="w-2 h-2 rounded-full bg-emerald-500 mt-1 flex-shrink-0"></span>
              <div>
                <p class="text-zinc-800 dark:text-zinc-200 text-[11px]">
                  <span class="font-semibold">Client Signoff</span>: Milestone 2 approved
                </p>
                <span class="text-[10px] text-zinc-400">45m ago</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- User Profile Avatar -->
      <div class="flex items-center pl-1">
        <div class="relative rounded-full ring-1 ring-[#E4E4E7] dark:ring-[#27272A]">
          <img
            class="w-6 h-6 rounded-full object-cover"
            :src="authStore.user?.avatar_url || store.currentUser.avatar_url"
            :alt="authStore.user?.name || store.currentUser.name"
            :title="authStore.user?.name || store.currentUser.name"
          />
          <span class="absolute bottom-0 right-0 w-1.5 h-1.5 rounded-full bg-emerald-500 ring-2 ring-white dark:ring-zinc-900"></span>
        </div>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useProjectStore } from '../stores/projectStore'
import { useAuthStore } from '../stores/authStore'
import {
  PanelLeftClose,
  PanelLeft,
  ChevronRight,
  Kanban,
  ListTodo,
  GitCommit,
  SlidersHorizontal,
  Check,
  Search,
  Plus,
  Bell,
  ShieldCheck
} from 'lucide-vue-next'

const store = useProjectStore()
const authStore = useAuthStore()
const router = useRouter()

const showNotifications = ref(false)
const showFilterPopover = ref(false)

const categories = ['All', 'Frontend', 'Backend', 'UI/UX', 'QA', 'Infrastructure']

function handleViewChange(view: 'kanban' | 'list' | 'timeline') {
  store.setView(view)
  router.push({ path: '/projects', query: { tab: view } })
}

function setCategoryFilter(cat: string) {
  store.setFilter(cat)
  showFilterPopover.value = false
}

function onRoleChange(role: string) {
  authStore.switchDevRole(role as any)
  store.showToast(`Switched active role to ${role.toUpperCase()}`)
}
</script>
