<template>
  <header class="h-[54px] border-b border-slate-200/90 bg-white/90 backdrop-blur-md px-3 sm:px-5 flex items-center justify-between flex-shrink-0 z-20 select-none shadow-2xs">
    <!-- Left: Hamburger & Breadcrumbs -->
    <div class="flex items-center space-x-2 sm:space-x-3 min-w-0">
      <!-- Toggle Sidebar (on mobile toggles mobile drawer, on desktop toggles collapse) -->
      <button
        @click="store.toggleSidebar()"
        class="p-1.5 rounded-xl text-slate-400 hover:text-slate-700 hover:bg-slate-100 transition-colors"
        title="Buka / Tutup Sidebar"
      >
        <Menu class="w-4 h-4 md:hidden" />
        <PanelLeft v-if="store.isSidebarCollapsed" class="w-4 h-4 hidden md:block" />
        <PanelLeftClose v-else class="w-4 h-4 hidden md:block" />
      </button>

      <!-- Single-row clean responsive breadcrumbs -->
      <nav aria-label="Breadcrumb" class="flex items-center space-x-1.5 text-xs text-slate-500 truncate">
        <span class="font-bold text-slate-800">Pulse</span>
        <ChevronRight class="w-3.5 h-3.5 text-slate-300 flex-shrink-0" />
        <span class="text-slate-500 hidden sm:inline">Workspace</span>
        <ChevronRight class="w-3.5 h-3.5 text-slate-300 flex-shrink-0 hidden sm:inline" />
        <span class="font-bold text-slate-900 truncate">
          {{ store.currentProject?.name || 'Dashboard' }}
        </span>
        <span class="hidden lg:inline-flex items-center px-2 py-0.5 rounded-full text-[10px] font-mono bg-slate-100 text-slate-600 border border-slate-200/70 tabular-nums ml-1">
          {{ store.filteredTasks.length }} Tasks
        </span>
      </nav>
    </div>

    <!-- Center: Compact View Switcher pills (Board | List | Timeline) -->
    <div class="hidden md:flex items-center p-0.5 rounded-xl bg-slate-100/90 border border-slate-200/80 text-xs font-medium">
      <button
        @click="handleViewChange('kanban')"
        :class="[
          'px-3 py-1 rounded-lg transition-all flex items-center gap-1.5 text-[11px]',
          store.currentView === 'kanban'
            ? 'bg-white text-slate-900 shadow-2xs font-bold'
            : 'text-slate-600 hover:text-slate-900'
        ]"
      >
        <Kanban class="w-3.5 h-3.5" :class="store.currentView === 'kanban' ? 'text-brand-600' : ''" />
        <span>Board</span>
      </button>

      <button
        @click="handleViewChange('list')"
        :class="[
          'px-3 py-1 rounded-lg transition-all flex items-center gap-1.5 text-[11px]',
          store.currentView === 'list'
            ? 'bg-white text-slate-900 shadow-2xs font-bold'
            : 'text-slate-600 hover:text-slate-900'
        ]"
      >
        <ListTodo class="w-3.5 h-3.5" :class="store.currentView === 'list' ? 'text-brand-600' : ''" />
        <span>List</span>
      </button>

      <button
        @click="handleViewChange('timeline')"
        :class="[
          'px-3 py-1 rounded-lg transition-all flex items-center gap-1.5 text-[11px]',
          store.currentView === 'timeline'
            ? 'bg-white text-slate-900 shadow-2xs font-bold'
            : 'text-slate-600 hover:text-slate-900'
        ]"
      >
        <GitCommit class="w-3.5 h-3.5" :class="store.currentView === 'timeline' ? 'text-brand-600' : ''" />
        <span>Timeline</span>
      </button>
    </div>

    <!-- Right: Filter Popover, Search, New Task, Notifications, Avatar Popover -->
    <div class="flex items-center space-x-1.5 sm:space-x-2.5 flex-shrink-0">
      <!-- Filter Popover Button -->
      <div class="relative">
        <button
          @click="showFilterPopover = !showFilterPopover"
          :class="[
            'flex items-center space-x-1.5 px-2.5 py-1.5 rounded-xl text-xs font-medium border transition-colors',
            store.activeFilter !== 'All'
              ? 'bg-brand-50 text-brand-700 border-brand-200 font-bold'
              : 'bg-white border-slate-200 text-slate-600 hover:text-slate-900 hover:bg-slate-50 shadow-2xs'
          ]"
        >
          <SlidersHorizontal class="w-3.5 h-3.5" />
          <span class="hidden sm:inline text-[11px]">{{ store.activeFilter === 'All' ? 'Filter' : store.activeFilter }}</span>
        </button>

        <!-- Filter Dropdown Menu -->
        <div
          v-if="showFilterPopover"
          @click.stop
          class="absolute right-0 mt-2 w-48 bg-white border border-slate-200/90 rounded-2xl shadow-xl z-50 p-1.5 space-y-0.5 text-xs animate-in fade-in zoom-in-95 duration-100"
        >
          <div class="px-2.5 py-1 text-[10px] font-mono uppercase text-slate-400 font-bold">Kategori Tugas</div>
          <button
            v-for="cat in categories"
            :key="cat"
            @click="setCategoryFilter(cat)"
            :class="[
              'w-full flex items-center justify-between px-2.5 py-1.5 rounded-xl text-left transition-colors',
              store.activeFilter === cat
                ? 'bg-brand-50 text-brand-700 font-bold'
                : 'text-slate-700 hover:bg-slate-100'
            ]"
          >
            <span>{{ cat }}</span>
            <Check v-if="store.activeFilter === cat" class="w-3.5 h-3.5 text-brand-600" />
          </button>
        </div>
      </div>

      <!-- Global Search trigger button -->
      <button
        @click="store.isCommandPaletteOpen = true"
        class="hidden sm:flex items-center space-x-2 text-xs text-slate-400 bg-slate-100/80 hover:bg-slate-200/70 px-2.5 py-1.5 rounded-xl border border-slate-200/60 transition-colors group"
      >
        <Search class="w-3.5 h-3.5 text-slate-400 group-hover:text-slate-600" />
        <span class="hidden xl:inline text-[11px] text-slate-500">Cari...</span>
        <kbd class="font-mono text-[9px] px-1 py-0.2 bg-white border border-slate-200 rounded-md text-slate-500 shadow-2xs">⌘K</kbd>
      </button>

      <!-- Quick "+ New Task" Button -->
      <button
        v-can="'task:create'"
        @click="store.isNewTaskModalOpen = true"
        class="flex items-center space-x-1.5 px-3 py-1.5 rounded-xl bg-gradient-to-r from-brand-600 to-indigo-600 hover:from-brand-500 hover:to-indigo-500 text-white text-xs font-bold shadow-md shadow-brand-600/20 transition-all active:scale-[0.98]"
      >
        <Plus class="w-3.5 h-3.5 stroke-[2.5]" />
        <span class="hidden sm:inline text-[11px]">Tugas Baru</span>
      </button>

      <!-- Notification Bell -->
      <div class="relative">
        <button
          @click="toggleNotifications"
          class="p-2 rounded-xl text-slate-400 hover:text-slate-700 hover:bg-slate-100 transition-colors relative"
          title="Notifikasi"
        >
          <Bell class="w-4 h-4" />
          <span
            v-if="unreadCount > 0"
            class="absolute top-1.5 right-1.5 w-2 h-2 rounded-full bg-rose-500 ring-2 ring-white"
          ></span>
        </button>

        <!-- Dropdown Notification Popover -->
        <div
          v-if="showNotifications"
          class="absolute right-0 mt-2 w-80 bg-white border border-slate-200/90 rounded-2xl shadow-2xl z-50 p-3 space-y-2 animate-in fade-in zoom-in-95 duration-100"
        >
          <div class="flex items-center justify-between pb-2 border-b border-slate-200">
            <div class="flex items-center gap-1.5">
              <span class="text-xs font-bold text-slate-900">Notifikasi</span>
              <span v-if="unreadCount > 0" class="px-1.5 py-0.2 rounded-full text-[9px] font-mono bg-rose-50 text-rose-700 border border-rose-200 font-bold">
                {{ unreadCount }}
              </span>
            </div>
            <button @click="markAllRead" class="text-[10px] text-brand-600 hover:underline font-medium">
              Tandai semua dibaca
            </button>
          </div>

          <div v-if="notifications.length === 0" class="py-4 text-center text-xs text-slate-400">
            Tidak ada notifikasi baru
          </div>

          <div v-else class="space-y-1.5 text-xs max-h-60 overflow-y-auto">
            <div
              v-for="notif in notifications"
              :key="notif.id"
              class="flex items-start space-x-2.5 p-2 rounded-xl transition-colors"
              :class="notif.is_read ? 'opacity-70 hover:bg-slate-50' : 'bg-brand-50/50 border border-brand-100'"
            >
              <span
                class="w-2 h-2 rounded-full mt-1 flex-shrink-0"
                :class="notif.type === 'milestone' ? 'bg-emerald-500' : 'bg-brand-600'"
              ></span>
              <div class="flex-1 min-w-0">
                <p class="text-slate-800 text-[11px] leading-snug">
                  <span class="font-bold">{{ notif.title }}</span>: {{ notif.message }}
                </p>
                <span class="text-[10px] text-slate-400 font-mono">{{ notif.created_at ? notif.created_at.slice(11, 16) : 'Baru' }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- User Profile Avatar with Interactive Dropdown Menu -->
      <div class="relative">
        <button
          @click="showUserMenu = !showUserMenu; showNotifications = false; showFilterPopover = false"
          class="flex items-center gap-2 p-1 rounded-xl hover:bg-slate-100 transition-colors cursor-pointer group"
          title="Buka Menu Akun / Ganti Pengguna"
        >
          <div class="relative rounded-full ring-2 ring-brand-500/30 group-hover:ring-brand-500 transition-all">
            <img
              class="w-7 h-7 rounded-full object-cover"
              :src="activeAvatar"
              :alt="activeName"
            />
            <span class="absolute bottom-0 right-0 w-2 h-2 rounded-full bg-emerald-500 ring-2 ring-white"></span>
          </div>
          <div class="hidden xl:flex flex-col text-left">
            <span class="text-xs font-bold text-slate-900 leading-tight truncate max-w-[110px]">
              {{ activeName }}
            </span>
            <span class="text-[9px] font-mono text-brand-600 uppercase leading-none font-semibold">
              {{ activeRoleText }}
            </span>
          </div>
          <ChevronDown class="w-3 h-3 text-slate-400 group-hover:text-slate-600 transition-transform" :class="showUserMenu ? 'rotate-180' : ''" />
        </button>

        <!-- Dropdown User Profile Popover -->
        <div
          v-if="showUserMenu"
          @click.stop
          class="absolute right-0 mt-2 w-80 bg-white border border-slate-200/90 rounded-2xl shadow-2xl z-50 p-3 space-y-3 animate-in fade-in zoom-in-95 duration-150"
        >
          <!-- Current User Profile Header -->
          <div class="flex items-center gap-3 pb-3 border-b border-slate-200">
            <div class="relative">
              <img
                :src="activeAvatar"
                :alt="activeName"
                class="w-11 h-11 rounded-xl object-cover ring-2 ring-brand-500/50"
              />
              <span class="absolute -bottom-0.5 -right-0.5 w-3 h-3 rounded-full bg-emerald-500 ring-2 ring-white"></span>
            </div>
            <div class="flex-1 min-w-0">
              <div class="flex items-center justify-between gap-1">
                <div class="text-xs font-bold text-slate-900 truncate">
                  {{ activeName }}
                </div>
                <span :class="['text-[9px] font-mono px-1.5 py-0.2 rounded font-bold uppercase', roleBadgeStyle]">
                  {{ authStore.currentRole }}
                </span>
              </div>
              <p class="text-[11px] text-slate-500 truncate mt-0.5">
                {{ authStore.user?.email || 'user@pulse.workspace' }}
              </p>
            </div>
          </div>

          <!-- Quick Account Switcher Section (Dev Only) -->
          <div v-if="!isProd && authStore.allowQuickLogin" class="space-y-1.5">
            <div class="flex items-center justify-between text-[10px] font-mono text-slate-400 uppercase font-semibold px-1">
              <span>Ganti Akun Pengguna</span>
              <span class="text-brand-600 text-[9px] font-bold">Dev Switch</span>
            </div>

            <div class="space-y-1 max-h-48 overflow-y-auto pr-0.5 text-xs">
              <button
                v-for="u in accountSwitchList"
                :key="u.key"
                @click="switchAccount(u.key)"
                :class="[
                  'w-full flex items-center justify-between p-2 rounded-xl text-left transition-all',
                  authStore.user?.email === u.email
                    ? 'bg-brand-50 text-brand-700 font-bold border border-brand-200'
                    : 'text-slate-700 hover:bg-slate-50'
                ]"
              >
                <div class="flex items-center gap-2.5 truncate">
                  <img :src="u.avatar" :alt="u.name" class="w-6 h-6 rounded-lg object-cover ring-1 ring-slate-200 flex-shrink-0" />
                  <div class="truncate">
                    <div class="text-xs font-semibold truncate">{{ u.name }}</div>
                    <div class="text-[10px] text-slate-400 font-mono leading-none">{{ u.roleLabel }}</div>
                  </div>
                </div>
                <Check v-if="authStore.user?.email === u.email" class="w-4 h-4 text-brand-600 flex-shrink-0 ml-1" />
              </button>
            </div>
          </div>

          <!-- Account Security / Change Password Action -->
          <div class="pt-2 border-t border-slate-200">
            <button
              @click="openChangePasswordModal"
              class="w-full py-2 px-3 rounded-xl text-xs font-semibold text-slate-700 hover:bg-brand-50/80 hover:text-brand-700 hover:border-brand-200 border border-slate-200/90 transition-all flex items-center justify-between cursor-pointer group shadow-2xs"
            >
              <div class="flex items-center gap-2">
                <KeyRound class="w-3.5 h-3.5 text-brand-600 group-hover:scale-110 transition-transform" />
                <span>Ganti Kata Sandi</span>
              </div>
              <ChevronRight class="w-3.5 h-3.5 text-slate-300 group-hover:text-brand-500 transition-colors" />
            </button>
          </div>

          <!-- Sign Out & Login Page Actions -->
          <div class="pt-2 border-t border-slate-200 flex items-center justify-between gap-2">
            <button
              @click="goToLoginPage"
              class="flex-1 py-2 px-2.5 rounded-xl text-xs font-semibold text-slate-700 hover:bg-slate-100 border border-slate-200/80 transition-colors flex items-center justify-center gap-1.5"
            >
              <LogIn class="w-3.5 h-3.5" />
              <span>Halaman Login</span>
            </button>

            <button
              @click="handleLogout"
              class="flex-1 py-2 px-2.5 rounded-xl text-xs font-semibold text-rose-600 hover:bg-rose-50 border border-rose-200 transition-colors flex items-center justify-center gap-1.5"
            >
              <LogOut class="w-3.5 h-3.5" />
              <span>Keluar (Logout)</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useProjectStore } from '../stores/projectStore'
import { useAuthStore, SEED_USERS } from '../stores/authStore'
import { projectApi } from '../services/api'
import type { NotificationItem } from '../types/project'
import {
  PanelLeftClose,
  PanelLeft,
  Menu,
  ChevronRight,
  ChevronDown,
  Kanban,
  ListTodo,
  GitCommit,
  SlidersHorizontal,
  Check,
  Search,
  Plus,
  Bell,
  LogOut,
  LogIn,
  KeyRound
} from 'lucide-vue-next'

const store = useProjectStore()
const authStore = useAuthStore()
const router = useRouter()
const isProd = computed(() => import.meta.env.PROD)

const showNotifications = ref(false)
const showFilterPopover = ref(false)
const showUserMenu = ref(false)

const notifications = ref<NotificationItem[]>([])
const unreadCount = computed(() => notifications.value.filter(n => !n.is_read).length)

async function loadNotifications() {
  try {
    const res = await projectApi.getNotifications()
    notifications.value = Array.isArray(res) ? res : (res?.notifications || [])
  } catch (err) {
    console.error('Failed to load notifications:', err)
  }
}

function toggleNotifications() {
  showNotifications.value = !showNotifications.value
  if (showNotifications.value) {
    loadNotifications()
  }
}

async function markAllRead() {
  try {
    await projectApi.markAllNotificationsRead()
    notifications.value = notifications.value.map(n => ({ ...n, is_read: true }))
    store.showToast('Semua notifikasi ditandai telah dibaca')
  } catch (err) {
    console.error('Failed to mark notifications read:', err)
  }
}

onMounted(() => {
  loadNotifications()
})

const categories = ['All', 'Frontend', 'Backend', 'UI/UX', 'QA', 'Infrastructure']

const activeName = computed(() => authStore.user?.name || store.currentUser.name || 'Agus Setiawan')
const activeAvatar = computed(() => authStore.user?.avatar_url || store.currentUser.avatar_url)

const activeRoleText = computed(() => {
  const role = authStore.currentRole
  if (role === 'super_admin') return 'Super Admin'
  if (role === 'pm') return 'Project Manager'
  if (role === 'member') return 'Team Member'
  if (role === 'client') return 'Client'
  return role.toUpperCase()
})

const roleBadgeStyle = computed(() => {
  const role = authStore.currentRole
  if (role === 'super_admin') return 'bg-rose-50 text-rose-700 border border-rose-200'
  if (role === 'pm') return 'bg-emerald-50 text-emerald-700 border border-emerald-200'
  if (role === 'member') return 'bg-indigo-50 text-indigo-700 border border-indigo-200'
  if (role === 'client') return 'bg-amber-50 text-amber-700 border border-amber-200'
  return 'bg-slate-100 text-slate-700 border border-slate-200'
})

const accountSwitchList = [
  { key: 'agus', name: 'Agus Setiawan', email: 'agus.setiawan@pulse.workspace', roleLabel: 'Project Manager (PM)', avatar: SEED_USERS['agus'].avatar_url },
  { key: 'adi', name: 'Adi Saputra', email: 'adi.saputra@pulse.workspace', roleLabel: 'Project Manager (PM)', avatar: SEED_USERS['adi'].avatar_url },
  { key: 'yudistira', name: 'Yudistira Sira P.', email: 'yudistira@pulse.workspace', roleLabel: 'Lead Fullstack (Member)', avatar: SEED_USERS['yudistira'].avatar_url },
  { key: 'hardianto', name: 'Hardianto', email: 'hardianto@pulse.workspace', roleLabel: 'Frontend Dev (Member)', avatar: SEED_USERS['hardianto'].avatar_url },
  { key: 'dedi', name: 'Dedi Irawan', email: 'dedi.irawan@pulse.workspace', roleLabel: 'DevOps & QA (Member)', avatar: SEED_USERS['dedi'].avatar_url },
  { key: 'super_admin', name: 'System Super Admin', email: 'admin@pulse.workspace', roleLabel: 'Root Super Admin', avatar: SEED_USERS['super_admin'].avatar_url },
  { key: 'client', name: 'Ir. H. Gunawan', email: 'client@pulse.workspace', roleLabel: 'Client External Portal', avatar: SEED_USERS['client'].avatar_url },
]

function handleViewChange(view: 'kanban' | 'list' | 'timeline') {
  store.setView(view)
  router.push({ path: '/projects', query: { tab: view } })
}

function setCategoryFilter(cat: string) {
  store.setFilter(cat)
  showFilterPopover.value = false
}

async function switchAccount(key: string) {
  showUserMenu.value = false
  const user = await authStore.switchUser(key)
  store.showToast(`Berhasil beralih ke: ${user.name} (${user.role.toUpperCase()})`)
}

function openChangePasswordModal() {
  showUserMenu.value = false
  authStore.isChangePasswordModalOpen = true
}

function goToLoginPage() {
  showUserMenu.value = false
  router.push('/login?switch=true')
}

async function handleLogout() {
  showUserMenu.value = false
  await authStore.logout()
  store.showToast('Anda telah keluar dari sistem')
  router.replace('/login')
}
</script>
