<template>
  <aside
    :class="[
      'bg-white/95 backdrop-blur-xl border-r border-slate-200/90 flex flex-col justify-between flex-shrink-0 z-30 transition-all duration-200 ease-in-out select-none shadow-xs',
      isMobile ? 'w-full h-full' : (store.isSidebarCollapsed ? 'w-16' : 'w-60')
    ]"
  >
    <!-- Top Workspace Section -->
    <div class="flex flex-col flex-1 min-h-0">
      <!-- Workspace Switcher -->
      <div class="p-2.5 sm:p-3 border-b border-slate-200/80">
        <div
          v-if="!store.isSidebarCollapsed || isMobile"
          class="flex items-center justify-between p-1.5 rounded-xl hover:bg-slate-100/70 transition-colors group cursor-pointer"
        >
          <div class="flex items-center space-x-2.5 min-w-0">
            <div class="w-7 h-7 rounded-lg bg-gradient-to-tr from-brand-600 to-indigo-600 flex items-center justify-center text-white font-display font-bold text-xs shadow-md shadow-brand-600/20">
              P
            </div>
            <div class="truncate">
              <div class="text-xs font-bold tracking-tight text-slate-900 flex items-center gap-1.5">
                <span class="truncate">Pulse HQ</span>
                <span class="text-[9px] font-mono px-1.5 py-0.2 rounded-full bg-brand-50 text-brand-700 font-semibold border border-brand-200">PRO</span>
              </div>
            </div>
          </div>
          
          <div class="flex items-center gap-1">
            <button
              v-if="isMobile"
              @click="closeSidebar"
              class="p-1 rounded-lg text-slate-400 hover:text-slate-700 hover:bg-slate-200/60 transition-colors"
            >
              <X class="w-4 h-4" />
            </button>
            <ChevronsUpDown v-else class="w-3.5 h-3.5 text-slate-400 group-hover:text-slate-600 flex-shrink-0" />
          </div>
        </div>

        <!-- Collapsed Workspace Icon (Desktop only) -->
        <div v-else class="flex justify-center py-0.5">
          <div
            @click="store.toggleSidebar()"
            class="w-8 h-8 rounded-lg bg-gradient-to-tr from-brand-600 to-indigo-600 flex items-center justify-center text-white font-display font-bold text-xs shadow-sm shadow-brand-600/20 cursor-pointer hover:scale-105 transition-all"
            title="Pulse HQ Workspace (Klik untuk memperluas)"
          >
            P
          </div>
        </div>

        <!-- Quick Search Bar Trigger -->
        <div class="mt-2.5">
          <button
            v-if="!store.isSidebarCollapsed || isMobile"
            @click="openCommandPalette"
            class="w-full flex items-center justify-between px-2.5 py-1.5 rounded-xl text-xs text-slate-500 bg-slate-100/80 hover:bg-slate-200/70 border border-slate-200/60 transition-colors"
          >
            <div class="flex items-center space-x-2 truncate">
              <Search class="w-3.5 h-3.5 text-slate-400" />
              <span class="text-[11px] text-slate-500">Cari atau navigasi...</span>
            </div>
            <kbd class="font-mono text-[9px] px-1.5 py-0.5 bg-white border border-slate-200 rounded-md text-slate-500 shadow-2xs">⌘K</kbd>
          </button>
          <button
            v-else
            @click="store.isCommandPaletteOpen = true"
            class="w-full flex justify-center p-2 rounded-xl text-slate-500 hover:bg-slate-100 transition-colors"
            title="Pencarian Cepat (⌘K)"
          >
            <Search class="w-4 h-4 text-slate-400" />
          </button>
        </div>
      </div>

      <!-- Navigation Items -->
      <div class="flex-1 overflow-y-auto px-2 py-3 space-y-4">
        <!-- Main Nav -->
        <div>
          <div
            v-if="!store.isSidebarCollapsed || isMobile"
            class="px-2 mb-1.5 text-[10px] font-bold uppercase tracking-wider text-slate-400 font-mono"
          >
            Navigasi Utama
          </div>
          <nav class="space-y-1 text-xs font-medium">
            <!-- Inbox / Notifications -->
            <button
              @click="handleInboxClick"
              :class="[
                'w-full flex items-center rounded-xl transition-all text-slate-600 hover:text-slate-900 hover:bg-slate-100/80 group',
                !isMobile && store.isSidebarCollapsed ? 'justify-center p-2' : 'justify-between px-2.5 py-2'
              ]"
              :title="!isMobile && store.isSidebarCollapsed ? 'Inbox (2 update)' : ''"
            >
              <div class="flex items-center space-x-2.5">
                <Inbox class="w-4 h-4 text-slate-400 group-hover:text-brand-600 transition-colors" />
                <span v-if="!store.isSidebarCollapsed || isMobile">Inbox</span>
              </div>
              <span
                v-if="!store.isSidebarCollapsed || isMobile"
                class="text-[10px] font-mono px-1.5 py-0.2 rounded-full bg-brand-50 text-brand-700 font-semibold border border-brand-200"
              >
                2
              </span>
            </button>

            <!-- Dashboard -->
            <RouterLink
              to="/dashboard"
              v-slot="{ isActive }"
              class="block"
              @click="closeSidebar"
            >
              <div
                :class="[
                  'w-full flex items-center rounded-xl transition-all group',
                  !isMobile && store.isSidebarCollapsed ? 'justify-center p-2' : 'justify-between px-2.5 py-2',
                  isActive
                    ? 'text-brand-700 bg-brand-50/90 border border-brand-200/90 font-bold shadow-2xs'
                    : 'text-slate-600 hover:text-slate-900 hover:bg-slate-100/80'
                ]"
                :title="!isMobile && store.isSidebarCollapsed ? 'Executive Dashboard' : ''"
              >
                <div class="flex items-center space-x-2.5">
                  <LayoutDashboard class="w-4 h-4" :class="isActive ? 'text-brand-600' : 'text-slate-400 group-hover:text-brand-600'" />
                  <span v-if="!store.isSidebarCollapsed || isMobile">Dashboard</span>
                </div>
              </div>
            </RouterLink>

            <!-- Projects Workspace -->
            <RouterLink
              to="/projects"
              v-slot="{ isActive }"
              class="block"
              @click="closeSidebar"
            >
              <div
                :class="[
                  'w-full flex items-center rounded-xl transition-all group',
                  !isMobile && store.isSidebarCollapsed ? 'justify-center p-2' : 'justify-between px-2.5 py-2',
                  isActive
                    ? 'text-brand-700 bg-brand-50/90 border border-brand-200/90 font-bold shadow-2xs'
                    : 'text-slate-600 hover:text-slate-900 hover:bg-slate-100/80'
                ]"
                :title="!isMobile && store.isSidebarCollapsed ? 'Project Workspace' : ''"
              >
                <div class="flex items-center space-x-2.5">
                  <FolderKanban class="w-4 h-4" :class="isActive ? 'text-brand-600' : 'text-slate-400 group-hover:text-brand-600'" />
                  <span v-if="!store.isSidebarCollapsed || isMobile">Projects</span>
                </div>
                <span
                  v-if="!store.isSidebarCollapsed || isMobile"
                  class="text-[10px] font-mono px-1.5 py-0.2 rounded-full bg-slate-100 text-slate-600 border border-slate-200/70 tabular-nums"
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
              @click="closeSidebar"
            >
              <div
                :class="[
                  'w-full flex items-center rounded-xl transition-all group',
                  !isMobile && store.isSidebarCollapsed ? 'justify-center p-2' : 'justify-between px-2.5 py-2',
                  isActive
                    ? 'text-indigo-700 bg-indigo-50/90 border border-indigo-200/90 font-bold shadow-2xs'
                    : 'text-slate-600 hover:text-slate-900 hover:bg-slate-100/80'
                ]"
                :title="!isMobile && store.isSidebarCollapsed ? 'Client External Portal' : ''"
              >
                <div class="flex items-center space-x-2.5">
                  <ShieldCheck class="w-4 h-4" :class="isActive ? 'text-indigo-600' : 'text-indigo-500 group-hover:text-indigo-600'" />
                  <span v-if="!store.isSidebarCollapsed || isMobile">Client Portal</span>
                </div>
                <span
                  v-if="!store.isSidebarCollapsed || isMobile"
                  class="text-[9px] font-mono px-1.5 py-0.2 rounded-full bg-indigo-50 text-indigo-700 font-semibold border border-indigo-200"
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
            v-if="!store.isSidebarCollapsed || isMobile"
            class="px-2 mb-1.5 flex items-center justify-between text-[10px] font-bold uppercase tracking-wider text-slate-400 font-mono"
          >
            <span>Daftar Proyek</span>
            <button
              v-can="'project:create'"
              @click="showNewProjectModal = true"
              class="hover:text-brand-600 transition-colors p-0.5 rounded"
              title="Tambah Proyek Baru"
            >
              <Plus class="w-3.5 h-3.5" />
            </button>
          </div>
          <div class="space-y-1 text-xs">
            <button
              v-for="project in store.projects"
              :key="project.id"
              @click="selectProject(project.id)"
              :class="[
                'w-full flex items-center rounded-xl transition-all group text-left',
                !isMobile && store.isSidebarCollapsed ? 'justify-center p-2' : 'justify-between px-2.5 py-1.5',
                store.currentProjectId === project.id
                  ? 'bg-slate-100 text-slate-900 font-bold border border-slate-200/80 shadow-2xs'
                  : 'text-slate-600 hover:bg-slate-50 hover:text-slate-900'
              ]"
              :title="!isMobile && store.isSidebarCollapsed ? project.name : ''"
            >
              <div class="flex items-center space-x-2 truncate">
                <span
                  class="w-2 h-2 rounded-full flex-shrink-0"
                  :style="{ backgroundColor: project.color_dot }"
                ></span>
                <span v-if="!store.isSidebarCollapsed || isMobile" class="truncate">{{ project.name }}</span>
              </div>
              <span
                v-if="!store.isSidebarCollapsed || isMobile"
                class="text-[10px] font-mono text-slate-400 group-hover:text-slate-600 tabular-nums"
              >
                {{ project.progress_percentage }}%
              </span>
            </button>
          </div>
        </div>

        <!-- Administration (Super Admin / PM only) -->
        <div v-if="canManageTeamOrSettings">
          <div
            v-if="!store.isSidebarCollapsed || isMobile"
            class="px-2 mb-1.5 text-[10px] font-bold uppercase tracking-wider text-slate-400 font-mono"
          >
            Administrasi
          </div>
          <nav class="space-y-1 text-xs font-medium">
            <button
              @click="handleTeamClick"
              :class="[
                'w-full flex items-center rounded-xl transition-all text-slate-600 hover:text-slate-900 hover:bg-slate-100/80 group',
                !isMobile && store.isSidebarCollapsed ? 'justify-center p-2' : 'justify-between px-2.5 py-2'
              ]"
              :title="!isMobile && store.isSidebarCollapsed ? 'Team & Access' : ''"
            >
              <div class="flex items-center space-x-2.5">
                <Users class="w-4 h-4 text-slate-400 group-hover:text-brand-600" />
                <span v-if="!store.isSidebarCollapsed || isMobile">Tim & Hak Akses</span>
              </div>
            </button>

            <button
              @click="handleSettingsClick"
              :class="[
                'w-full flex items-center rounded-xl transition-all text-slate-600 hover:text-slate-900 hover:bg-slate-100/80 group',
                !isMobile && store.isSidebarCollapsed ? 'justify-center p-2' : 'justify-between px-2.5 py-2'
              ]"
              :title="!isMobile && store.isSidebarCollapsed ? 'Pengaturan' : ''"
            >
              <div class="flex items-center space-x-2.5">
                <Settings class="w-4 h-4 text-slate-400 group-hover:text-brand-600" />
                <span v-if="!store.isSidebarCollapsed || isMobile">Pengaturan</span>
              </div>
            </button>
          </nav>
        </div>
      </div>
    </div>

    <!-- Sidebar Footer (Active User Card, Theme & Collapse Toggle) -->
    <div class="p-2 sm:p-2.5 border-t border-slate-200/80 space-y-2 bg-slate-50/50">
      <!-- Active User Profile Card -->
      <div
        v-if="!store.isSidebarCollapsed || isMobile"
        class="p-2 rounded-xl bg-white border border-slate-200/90 shadow-xs flex items-center justify-between gap-2"
      >
        <div class="flex items-center gap-2 min-w-0">
          <div class="relative flex-shrink-0">
            <img
              :src="authStore.user?.avatar_url || store.currentUser.avatar_url"
              :alt="authStore.user?.name || store.currentUser.name"
              class="w-7 h-7 rounded-lg object-cover ring-1 ring-slate-200"
            />
            <span class="absolute -bottom-0.5 -right-0.5 w-2 h-2 rounded-full bg-emerald-500 ring-2 ring-white"></span>
          </div>
          <div class="min-w-0">
            <div class="text-xs font-bold text-slate-900 truncate leading-tight">
              {{ authStore.user?.name || store.currentUser.name }}
            </div>
            <div class="text-[9px] font-mono text-brand-600 font-semibold uppercase leading-none mt-0.5">
              {{ authStore.currentRole }}
            </div>
          </div>
        </div>
        <button
          @click="logoutUser"
          class="p-1.5 rounded-lg text-slate-400 hover:text-rose-600 hover:bg-rose-50 transition-colors flex-shrink-0"
          title="Keluar (Logout)"
        >
          <LogOut class="w-3.5 h-3.5" />
        </button>
      </div>

      <!-- Collapse Toggle Button Row (Desktop only) -->
      <div v-if="!isMobile" class="flex items-center justify-between gap-1">
        <button
          @click="store.toggleSidebar()"
          class="w-full flex items-center justify-center p-1.5 rounded-xl text-slate-400 hover:text-slate-700 hover:bg-slate-100 transition-colors"
          :title="store.isSidebarCollapsed ? 'Perluas Sidebar' : 'Ciutkan Sidebar'"
        >
          <PanelLeftClose v-if="!store.isSidebarCollapsed" class="w-4 h-4" />
          <PanelLeft v-else class="w-4 h-4" />
        </button>
      </div>
    </div>
  </aside>

  <!-- New Project Modal Dialog -->
  <NewProjectModal
    :isOpen="showNewProjectModal"
    @close="showNewProjectModal = false"
  />
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useProjectStore } from '../stores/projectStore'
import { useAuthStore } from '../stores/authStore'
import { useRouter } from 'vue-router'
import NewProjectModal from './NewProjectModal.vue'
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
  PanelLeftClose,
  PanelLeft,
  LogOut,
  X
} from 'lucide-vue-next'

interface Props {
  isMobile?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  isMobile: false
})

const emit = defineEmits(['close'])

const store = useProjectStore()
const authStore = useAuthStore()
const router = useRouter()

const showNewProjectModal = ref(false)

const canManageTeamOrSettings = computed(() => {
  return authStore.hasRole(['super_admin', 'pm']) || authStore.can('user:manage') || authStore.can('project:member:manage')
})

function closeSidebar() {
  if (props.isMobile) {
    emit('close')
    store.closeMobileSidebar()
  }
}

function openCommandPalette() {
  store.isCommandPaletteOpen = true
  closeSidebar()
}

function handleInboxClick() {
  store.showToast('Anda memiliki 2 pembaruan tugas proyek')
  closeSidebar()
}

function handleTeamClick() {
  closeSidebar()
  router.push('/dashboard?tab=users')
}

function handleSettingsClick() {
  store.showToast('Pengaturan Workspace')
  closeSidebar()
}

async function logoutUser() {
  closeSidebar()
  await authStore.logout()
  store.showToast('Anda telah keluar dari akun')
  router.replace('/login')
}

function selectProject(projectId: string) {
  store.selectProject(projectId)
  closeSidebar()
  router.push(`/projects/${projectId}`)
}
</script>
