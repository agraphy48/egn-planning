<template>
  <div class="flex h-full w-full overflow-hidden bg-[#FAFAFA] dark:bg-[#09090B] text-zinc-900 dark:text-zinc-100 font-sans">
    <!-- Collapsible Slim Sidebar (240px expanded, 64px icon-only) -->
    <Sidebar />

    <!-- Main Workspace Container -->
    <div class="flex-1 flex flex-col min-w-0 bg-[#FAFAFA] dark:bg-[#09090B] overflow-hidden">
      <!-- Minimalist Topbar (52px height) -->
      <Topbar />

      <!-- View Canvas Area -->
      <main class="flex-1 overflow-y-auto overflow-x-hidden p-3.5 sm:p-5 relative">
        <RouterView />
      </main>

      <!-- Mobile Bottom Navigation (PWA Ergonomics) -->
      <nav class="md:hidden border-t border-[#E4E4E7] dark:border-[#27272A] bg-white/95 dark:bg-[#18181B]/95 backdrop-blur-md px-4 py-2 flex items-center justify-around text-xs font-medium z-30 flex-shrink-0">
        <RouterLink
          to="/dashboard"
          class="flex flex-col items-center space-y-1 text-zinc-500 hover:text-brand-600"
          active-class="text-brand-600 font-semibold"
        >
          <LayoutDashboard class="w-4 h-4" />
          <span class="text-[10px]">Dashboard</span>
        </RouterLink>

        <RouterLink
          to="/projects"
          class="flex flex-col items-center space-y-1 text-zinc-500 hover:text-brand-600"
          active-class="text-brand-600 font-semibold"
        >
          <FolderKanban class="w-4 h-4" />
          <span class="text-[10px]">Tasks</span>
        </RouterLink>

        <button
          v-can="'task:create'"
          @click="store.isNewTaskModalOpen = true"
          class="w-9 h-9 rounded-full bg-brand-600 text-white flex items-center justify-center -mt-4 shadow-lg ring-4 ring-white dark:ring-zinc-900"
        >
          <Plus class="w-5 h-5" />
        </button>

        <RouterLink
          to="/client-portal"
          class="flex flex-col items-center space-y-1 text-zinc-500 hover:text-brand-600"
          active-class="text-brand-600 font-semibold"
        >
          <ShieldCheck class="w-4 h-4" />
          <span class="text-[10px]">Client</span>
        </RouterLink>

        <button
          @click="store.toggleSidebar()"
          class="flex flex-col items-center space-y-1 text-zinc-500 hover:text-zinc-800"
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

    <!-- Global Toast Notification Banner -->
    <Transition name="toast">
      <div
        v-if="store.toastMessage"
        class="fixed bottom-5 right-5 z-50 px-3.5 py-2 rounded-xl bg-zinc-900 dark:bg-white text-white dark:text-zinc-900 text-xs font-semibold shadow-xl flex items-center space-x-2 border border-zinc-700 dark:border-zinc-200"
      >
        <CheckCircle2 class="w-4 h-4 text-emerald-400 dark:text-emerald-600" />
        <span>{{ store.toastMessage }}</span>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue'
import { useProjectStore } from '../stores/projectStore'
import Sidebar from '../components/Sidebar.vue'
import Topbar from '../components/Topbar.vue'
import TaskDetailDrawer from '../components/TaskDetailDrawer.vue'
import NewTaskModal from '../components/NewTaskModal.vue'
import CommandPalette from '../components/CommandPalette.vue'
import {
  LayoutDashboard,
  FolderKanban,
  Plus,
  ShieldCheck,
  Menu,
  CheckCircle2
} from 'lucide-vue-next'

const store = useProjectStore()

function handleKeyDown(e: KeyboardEvent) {
  // Cmd+K or Ctrl+K for search command palette
  if ((e.metaKey || e.ctrlKey) && e.key.toLowerCase() === 'k') {
    e.preventDefault()
    store.isCommandPaletteOpen = !store.isCommandPaletteOpen
  }

  // Escape closes all overlays
  if (e.key === 'Escape') {
    store.closeTaskDetail()
    store.isNewTaskModalOpen = false
    store.isCommandPaletteOpen = false
  }

  // 'N' triggers new task modal when not focused on an input or textarea
  if (
    e.key.toLowerCase() === 'n' &&
    !['INPUT', 'TEXTAREA', 'SELECT'].includes((e.target as HTMLElement).tagName)
  ) {
    e.preventDefault()
    store.isNewTaskModalOpen = true
  }
}

onMounted(() => {
  store.fetchInitialData()
  window.addEventListener('keydown', handleKeyDown)
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleKeyDown)
})
</script>

<style scoped>
.toast-enter-active,
.toast-leave-active {
  transition: all 0.25s cubic-bezier(0.16, 1, 0.3, 1);
}
.toast-enter-from,
.toast-leave-to {
  opacity: 0;
  transform: translateY(12px) scale(0.96);
}
</style>
