<template>
  <div class="flex h-full w-full overflow-hidden bg-gradient-to-br from-slate-50 via-white to-indigo-50/20 text-slate-900 font-sans relative">
    
    <!-- Ambient Subtle Pastel Glows (Matching Login Theme) -->
    <div class="fixed -top-40 -left-40 w-96 h-96 bg-brand-500/10 rounded-full blur-[100px] pointer-events-none z-0"></div>
    <div class="fixed -bottom-40 -right-40 w-96 h-96 bg-indigo-500/10 rounded-full blur-[120px] pointer-events-none z-0"></div>
    <div class="fixed top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[700px] h-[700px] bg-purple-500/5 rounded-full blur-[150px] pointer-events-none z-0"></div>

    <!-- Clean Tech Dot Grid Pattern (Matching Login Theme) -->
    <div class="fixed inset-0 bg-[radial-gradient(#cbd5e1_1px,transparent_1px)] [background-size:24px_24px] pointer-events-none opacity-50 z-0"></div>

    <!-- Desktop Slim Sidebar (Hidden on Mobile < md, visible md and up) -->
    <div class="hidden md:flex flex-shrink-0 z-30">
      <Sidebar />
    </div>

    <!-- Mobile Slide-Over Drawer with Backdrop -->
    <Transition name="fade">
      <div
        v-if="store.isMobileSidebarOpen"
        @click="store.closeMobileSidebar()"
        class="fixed inset-0 bg-slate-900/40 backdrop-blur-xs z-40 md:hidden"
      ></div>
    </Transition>
    <Transition name="slide-sidebar">
      <div
        v-if="store.isMobileSidebarOpen"
        class="fixed inset-y-0 left-0 z-50 w-72 max-w-[85vw] bg-white shadow-2xl md:hidden flex flex-col"
      >
        <Sidebar :isMobile="true" @close="store.closeMobileSidebar()" />
      </div>
    </Transition>

    <!-- Main Workspace Container -->
    <div class="flex-1 flex flex-col min-w-0 bg-transparent overflow-hidden relative z-10">
      <!-- Minimalist Topbar -->
      <Topbar />

      <!-- View Canvas Area -->
      <main class="flex-1 overflow-y-auto overflow-x-hidden p-3.5 sm:p-5 lg:p-6 relative">
        <slot>
          <RouterView />
        </slot>
      </main>

      <!-- Mobile Bottom Navigation (Clean White Theme) -->
      <nav class="md:hidden border-t border-slate-200/90 bg-white/95 backdrop-blur-xl px-4 py-2 flex items-center justify-around text-xs font-medium z-30 flex-shrink-0 shadow-lg shadow-slate-200/50">
        <RouterLink
          to="/dashboard"
          class="flex flex-col items-center space-y-1 text-slate-500 hover:text-brand-600 transition-colors"
          active-class="text-brand-600 font-bold"
        >
          <LayoutDashboard class="w-4 h-4" />
          <span class="text-[10px]">Dashboard</span>
        </RouterLink>

        <RouterLink
          to="/projects"
          class="flex flex-col items-center space-y-1 text-slate-500 hover:text-brand-600 transition-colors"
          active-class="text-brand-600 font-bold"
        >
          <FolderKanban class="w-4 h-4" />
          <span class="text-[10px]">Projects</span>
        </RouterLink>

        <button
          v-can="'task:create'"
          @click="store.isNewTaskModalOpen = true"
          class="w-10 h-10 rounded-full bg-gradient-to-tr from-brand-600 to-indigo-600 text-white flex items-center justify-center -mt-5 shadow-lg shadow-brand-600/30 ring-4 ring-white active:scale-95 transition-all"
        >
          <Plus class="w-5 h-5 stroke-[2.5]" />
        </button>

        <RouterLink
          to="/client-portal"
          class="flex flex-col items-center space-y-1 text-slate-500 hover:text-brand-600 transition-colors"
          active-class="text-brand-600 font-bold"
        >
          <ShieldCheck class="w-4 h-4" />
          <span class="text-[10px]">Client</span>
        </RouterLink>

        <button
          @click="store.toggleMobileSidebar()"
          class="flex flex-col items-center space-y-1 text-slate-500 hover:text-slate-900 transition-colors"
        >
          <Menu class="w-4 h-4" />
          <span class="text-[10px]">Menu</span>
        </button>
      </nav>
    </div>

    <!-- Slide-over Task Inspector Drawer -->
    <TaskDetailDrawer />

    <!-- New Task Modal -->
    <NewTaskModal />

    <!-- Command Palette (Cmd+K) -->
    <CommandPalette />

    <!-- Change Password Modal -->
    <ChangePasswordModal />

    <!-- Global Toast Notification Banner -->
    <Transition name="toast">
      <div
        v-if="store.toastMessage"
        class="fixed bottom-5 right-5 z-50 px-4 py-2.5 rounded-xl bg-slate-900 text-white text-xs font-semibold shadow-xl flex items-center space-x-2 border border-slate-700/80 animate-slide-up"
      >
        <CheckCircle2 class="w-4 h-4 text-emerald-400" />
        <span>{{ store.toastMessage }}</span>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue'
import { useProjectStore } from '../stores/projectStore'
import { useAuthStore } from '../stores/authStore'
import Sidebar from '../components/Sidebar.vue'
import Topbar from '../components/Topbar.vue'
import TaskDetailDrawer from '../components/TaskDetailDrawer.vue'
import NewTaskModal from '../components/NewTaskModal.vue'
import CommandPalette from '../components/CommandPalette.vue'
import ChangePasswordModal from '../components/ChangePasswordModal.vue'
import {
  LayoutDashboard,
  FolderKanban,
  Plus,
  ShieldCheck,
  Menu,
  CheckCircle2
} from 'lucide-vue-next'

const store = useProjectStore()
const authStore = useAuthStore()

// Global Keyboard Shortcuts (Linear-like Power User Ergonomics)
const handleGlobalKeydown = (e: KeyboardEvent) => {
  // Cmd+K / Ctrl+K: Command Palette
  if ((e.metaKey || e.ctrlKey) && e.key.toLowerCase() === 'k') {
    e.preventDefault()
    store.isCommandPaletteOpen = !store.isCommandPaletteOpen
  }

  // 'C' shortcut: Open New Task Modal (only if not focused in an input/textarea)
  if (
    e.key.toLowerCase() === 'c' &&
    !['input', 'textarea'].includes((document.activeElement?.tagName || '').toLowerCase()) &&
    !e.metaKey &&
    !e.ctrlKey
  ) {
    if (authStore.can('task:create')) {
      e.preventDefault()
      store.isNewTaskModalOpen = true
    }
  }

  // Escape key closes modal & drawer & mobile sidebar
  if (e.key === 'Escape') {
    if (store.isMobileSidebarOpen) {
      store.closeMobileSidebar()
    }
    if (store.isDrawerOpen) {
      store.closeTaskDetail()
    }
    if (store.isNewTaskModalOpen) {
      store.isNewTaskModalOpen = false
    }
    if (store.isCommandPaletteOpen) {
      store.isCommandPaletteOpen = false
    }
  }
}

onMounted(() => {
  window.addEventListener('keydown', handleGlobalKeydown)
  // Fetch initial project data
  store.fetchInitialData()
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleGlobalKeydown)
})
</script>

<style scoped>
.slide-sidebar-enter-active,
.slide-sidebar-leave-active {
  transition: transform 0.25s cubic-bezier(0.16, 1, 0.3, 1);
}
.slide-sidebar-enter-from,
.slide-sidebar-leave-to {
  transform: translateX(-100%);
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.toast-enter-active,
.toast-leave-active {
  transition: all 0.25s ease;
}
.toast-enter-from,
.toast-leave-to {
  opacity: 0;
  transform: translateY(10px);
}
</style>
