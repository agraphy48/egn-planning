<template>
  <div
    v-if="store.isCommandPaletteOpen"
    class="fixed inset-0 bg-black/40 backdrop-blur-[2px] z-50 flex items-start justify-center p-4 pt-16 sm:pt-24 select-none"
    @click.self="store.isCommandPaletteOpen = false"
  >
    <div
      class="w-full max-w-xl bg-white dark:bg-[#18181B] border border-[#E4E4E7] dark:border-[#27272A] rounded-xl shadow-palette overflow-hidden animate-in fade-in zoom-in-95 duration-150"
    >
      <!-- Search Input Header -->
      <div class="p-3 border-b border-[#E4E4E7] dark:border-[#27272A] flex items-center space-x-3">
        <Search class="w-4 h-4 text-zinc-400 flex-shrink-0" />
        <input
          ref="searchInputRef"
          type="text"
          v-model="query"
          placeholder="Search tasks, projects, or type a command..."
          class="w-full bg-transparent text-sm text-zinc-900 dark:text-zinc-100 placeholder-zinc-400 focus:outline-none"
          autofocus
          @keydown.down.prevent="navigateResults(1)"
          @keydown.up.prevent="navigateResults(-1)"
          @keydown.enter.prevent="executeSelected"
        />
        <kbd class="font-mono text-[9px] px-1.5 py-0.5 bg-[#F4F4F5] dark:bg-[#27272A] rounded text-zinc-400">ESC</kbd>
      </div>

      <!-- Scrollable List of Actions and Matching Results -->
      <div class="p-2 space-y-1 text-xs max-h-80 overflow-y-auto">
        <!-- Matching Tasks Section -->
        <template v-if="filteredMatchingTasks.length > 0">
          <div class="px-2 pt-1 pb-1 text-[10px] font-mono uppercase text-zinc-400 font-semibold tracking-wider">
            Matching Tasks ({{ filteredMatchingTasks.length }})
          </div>
          <button
            v-for="(task, idx) in filteredMatchingTasks"
            :key="task.id"
            @click="selectTask(task.id)"
            @mouseenter="selectedIndex = idx"
            :class="[
              'w-full flex items-center justify-between p-2 rounded-lg text-left transition-colors',
              selectedIndex === idx
                ? 'bg-[#F4F4F5] dark:bg-[#27272A] text-zinc-900 dark:text-white'
                : 'text-zinc-600 dark:text-zinc-300'
            ]"
          >
            <div class="flex items-center space-x-2.5 truncate">
              <span class="font-mono text-[10px] text-zinc-400 font-semibold">{{ task.key }}</span>
              <span class="truncate font-medium text-xs">{{ task.title }}</span>
            </div>
            <div class="flex items-center space-x-1.5 flex-shrink-0">
              <span class="text-[10px] font-mono px-1.5 py-0.2 rounded bg-zinc-100 dark:bg-zinc-800 text-zinc-500">
                {{ task.status }}
              </span>
            </div>
          </button>
        </template>

        <!-- Matching Projects Section -->
        <template v-if="filteredProjects.length > 0">
          <div class="px-2 pt-2 pb-1 text-[10px] font-mono uppercase text-zinc-400 font-semibold tracking-wider">
            Projects
          </div>
          <button
            v-for="(proj, pIdx) in filteredProjects"
            :key="proj.id"
            @click="selectProject(proj.id)"
            :class="[
              'w-full flex items-center justify-between p-2 rounded-lg text-left transition-colors',
              selectedIndex === filteredMatchingTasks.length + pIdx
                ? 'bg-[#F4F4F5] dark:bg-[#27272A] text-zinc-900 dark:text-white'
                : 'text-zinc-600 dark:text-zinc-300'
            ]"
          >
            <div class="flex items-center space-x-2.5 truncate">
              <span class="w-2 h-2 rounded-full" :style="{ backgroundColor: proj.color_dot }"></span>
              <span class="truncate font-medium text-xs">{{ proj.name }}</span>
            </div>
            <span class="text-[10px] font-mono text-zinc-400 tabular-nums">
              {{ proj.progress_percentage }}%
            </span>
          </button>
        </template>

        <!-- Quick Navigation & Commands -->
        <div class="px-2 pt-2 pb-1 text-[10px] font-mono uppercase text-zinc-400 font-semibold tracking-wider">
          Quick Actions
        </div>
        <button
          v-for="(cmd, cIdx) in quickCommands"
          :key="cmd.id"
          @click="cmd.action"
          :class="[
            'w-full flex items-center justify-between p-2 rounded-lg text-left transition-colors',
            selectedIndex === totalItemsBeforeCommands + cIdx
              ? 'bg-[#F4F4F5] dark:bg-[#27272A] text-zinc-900 dark:text-white'
              : 'text-zinc-600 dark:text-zinc-300'
          ]"
        >
          <span class="flex items-center gap-2">
            <component :is="cmd.icon" class="w-3.5 h-3.5 text-zinc-400" />
            <span class="text-xs">{{ cmd.label }}</span>
          </span>
          <kbd class="font-mono text-[9px] text-zinc-400 bg-zinc-100 dark:bg-zinc-800 px-1 py-0.2 rounded">{{ cmd.shortcut }}</kbd>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { useProjectStore } from '../stores/projectStore'
import {
  Search,
  PlusCircle,
  LayoutDashboard,
  Kanban,
  ListTodo,
  GitCommit,
  ShieldCheck,
  Moon
} from 'lucide-vue-next'

const store = useProjectStore()
const router = useRouter()
const query = ref('')
const searchInputRef = ref<HTMLInputElement | null>(null)
const selectedIndex = ref(0)

const filteredMatchingTasks = computed(() => {
  if (!query.value.trim()) return []
  return store.tasks.filter((t) =>
    t.title.toLowerCase().includes(query.value.toLowerCase()) ||
    t.key.toLowerCase().includes(query.value.toLowerCase())
  ).slice(0, 5)
})

const filteredProjects = computed(() => {
  if (!query.value.trim()) return []
  return store.projects.filter((p) =>
    p.name.toLowerCase().includes(query.value.toLowerCase())
  ).slice(0, 3)
})

const totalItemsBeforeCommands = computed(() => {
  return filteredMatchingTasks.value.length + filteredProjects.value.length
})

const quickCommands = [
  {
    id: 'create-task',
    label: 'Create New Task',
    icon: PlusCircle,
    shortcut: 'N',
    action: () => {
      store.isCommandPaletteOpen = false
      store.isNewTaskModalOpen = true
    }
  },
  {
    id: 'dashboard',
    label: 'Go to Executive Dashboard',
    icon: LayoutDashboard,
    shortcut: 'G D',
    action: () => {
      goTo('/dashboard')
    }
  },
  {
    id: 'kanban-board',
    label: 'Switch to Kanban Board',
    icon: Kanban,
    shortcut: 'G B',
    action: () => {
      store.setView('kanban')
      goTo('/projects?tab=kanban')
    }
  },
  {
    id: 'list-view',
    label: 'Switch to List View',
    icon: ListTodo,
    shortcut: 'G L',
    action: () => {
      store.setView('list')
      goTo('/projects?tab=list')
    }
  },
  {
    id: 'timeline-view',
    label: 'Switch to Timeline (Gantt)',
    icon: GitCommit,
    shortcut: 'G T',
    action: () => {
      store.setView('timeline')
      goTo('/projects?tab=timeline')
    }
  },
  {
    id: 'client-portal',
    label: 'Open Client Portal',
    icon: ShieldCheck,
    shortcut: 'G C',
    action: () => {
      goTo('/client-portal')
    }
  },
  {
    id: 'theme',
    label: 'Toggle Dark / Light Mode',
    icon: Moon,
    shortcut: 'T',
    action: () => {
      store.toggleTheme()
      store.isCommandPaletteOpen = false
    }
  }
]

function navigateResults(delta: number) {
  const total = totalItemsBeforeCommands.value + quickCommands.length
  if (total === 0) return
  selectedIndex.value = (selectedIndex.value + delta + total) % total
}

function executeSelected() {
  if (selectedIndex.value < filteredMatchingTasks.value.length) {
    const task = filteredMatchingTasks.value[selectedIndex.value]
    selectTask(task.id)
  } else if (selectedIndex.value < totalItemsBeforeCommands.value) {
    const pIdx = selectedIndex.value - filteredMatchingTasks.value.length
    const proj = filteredProjects.value[pIdx]
    selectProject(proj.id)
  } else {
    const cIdx = selectedIndex.value - totalItemsBeforeCommands.value
    if (quickCommands[cIdx]) {
      quickCommands[cIdx].action()
    }
  }
}

function goTo(path: string) {
  router.push(path)
  store.isCommandPaletteOpen = false
}

function selectTask(taskId: string) {
  store.isCommandPaletteOpen = false
  store.openTaskDetail(taskId)
}

function selectProject(projectId: string) {
  store.selectProject(projectId)
  router.push('/projects')
  store.isCommandPaletteOpen = false
}

watch(() => store.isCommandPaletteOpen, (open) => {
  if (open) {
    query.value = ''
    selectedIndex.value = 0
    nextTick(() => {
      searchInputRef.value?.focus()
    })
  }
})
</script>
