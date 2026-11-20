<template>
  <div class="flex-1 flex flex-col space-y-3 min-h-0 select-none">
    <!-- Subtle Read-Only Mode Banner if !canMove -->
    <div
      v-if="!canMove"
      class="flex items-center justify-between px-3.5 py-2 bg-amber-50/80 dark:bg-amber-950/30 border border-amber-200/60 dark:border-amber-800/40 rounded-xl text-amber-800 dark:text-amber-300 text-xs shadow-2xs flex-shrink-0"
    >
      <div class="flex items-center space-x-2">
        <Lock class="w-3.5 h-3.5 text-amber-600 dark:text-amber-400 flex-shrink-0" />
        <span class="font-semibold">Read-Only Mode:</span>
        <span class="text-zinc-600 dark:text-zinc-400">
          Kanban drag-and-drop & status modifications are disabled for your role (<strong>{{ authStore.currentRole.toUpperCase() }}</strong>).
        </span>
      </div>
      <span class="text-[10px] font-mono uppercase px-2 py-0.5 rounded bg-amber-100 dark:bg-amber-900/60 text-amber-800 dark:text-amber-200 font-semibold border border-amber-200 dark:border-amber-700">
        Read Only
      </span>
    </div>

    <!-- Kanban Board Columns Row -->
    <div class="flex-1 flex space-x-3.5 overflow-x-auto pb-4 items-start min-h-0">
      <div
        v-for="col in columns"
        :key="col.name"
        class="kanban-col w-72 lg:w-80 flex-shrink-0 bg-[#F4F4F5]/70 dark:bg-[#18181B]/70 rounded-xl p-2.5 border border-[#E4E4E7] dark:border-[#27272A] flex flex-col max-h-full transition-colors"
        :class="{ 'border-brand-500 ring-2 ring-brand-500/20 bg-brand-50/10 dark:bg-brand-950/20': canMove && dragOverColumn === col.name }"
        @dragover.prevent="canMove && onDragOver(col.name)"
        @dragleave="onDragLeave"
        @drop="canMove && onDrop(col.name, $event)"
      >
        <!-- Column Header -->
        <div class="flex items-center justify-between pb-2.5 px-1 flex-shrink-0">
          <div class="flex items-center space-x-2">
            <span class="w-2 h-2 rounded-full" :class="col.dotClass"></span>
            <span class="text-xs font-semibold text-zinc-900 dark:text-zinc-100 tracking-tight">{{ col.name }}</span>
            <span
              class="text-[10px] font-mono px-1.5 py-0.2 rounded-full font-medium tabular-nums"
              :class="col.badgeClass"
            >
              {{ store.tasksByStatus[col.name]?.length || 0 }}
            </span>
          </div>

          <button
            v-can="'task:create'"
            @click="openQuickAdd(col.name)"
            class="p-1 rounded text-zinc-400 hover:text-zinc-700 dark:hover:text-zinc-200 hover:bg-zinc-200/60 dark:hover:bg-zinc-800 transition-colors"
            title="Quick add task"
          >
            <Plus class="w-3.5 h-3.5" />
          </button>
        </div>

        <!-- Column Cards Container -->
        <div class="flex-1 overflow-y-auto space-y-2 pr-0.5 min-h-[140px]">
          <TaskCard
            v-for="task in store.tasksByStatus[col.name]"
            :key="task.id"
            :task="task"
          />

          <!-- Empty State in column -->
          <div
            v-if="!store.tasksByStatus[col.name] || store.tasksByStatus[col.name].length === 0"
            class="h-24 border border-dashed border-zinc-200 dark:border-zinc-800 rounded-lg flex flex-col items-center justify-center text-zinc-400 text-xs p-3 text-center"
          >
            <span class="text-[11px] font-mono text-zinc-400">Empty column</span>
          </div>

          <!-- Quick-Add Task Box at bottom of column -->
          <div v-can="'task:create'" class="pt-1">
            <!-- Active Quick-Add Input -->
            <div
              v-if="activeQuickAddColumn === col.name"
              class="p-2.5 rounded-lg bg-white dark:bg-[#18181B] border border-brand-500 shadow-sm space-y-2"
            >
              <textarea
                ref="quickAddInputRef"
                rows="2"
                v-model="quickAddTitle"
                placeholder="Task title... (Enter to save)"
                class="w-full text-xs text-zinc-900 dark:text-zinc-100 bg-transparent focus:outline-none resize-none leading-relaxed"
                @keydown.enter.prevent="submitQuickAdd(col.name)"
                @keydown.esc="cancelQuickAdd"
              ></textarea>
              <div class="flex items-center justify-between text-[11px] pt-1 border-t border-zinc-100 dark:border-zinc-800">
                <span class="text-zinc-400 font-mono text-[10px]">Enter ↵ to add</span>
                <div class="flex items-center space-x-1.5">
                  <button
                    @click="cancelQuickAdd"
                    class="px-2 py-0.5 rounded text-zinc-400 hover:text-zinc-600 dark:hover:text-zinc-200"
                  >
                    Cancel
                  </button>
                  <button
                    @click="submitQuickAdd(col.name)"
                    class="px-2.5 py-0.5 rounded bg-brand-600 hover:bg-brand-500 text-white font-medium shadow-2xs"
                  >
                    Add
                  </button>
                </div>
              </div>
            </div>

            <!-- Inactive "+ Add task" trigger button -->
            <button
              v-else
              @click="openQuickAdd(col.name)"
              class="w-full py-1.5 px-2 rounded-lg text-left text-xs text-zinc-400 hover:text-zinc-700 dark:hover:text-zinc-200 hover:bg-zinc-200/50 dark:hover:bg-zinc-800/50 transition-colors flex items-center space-x-1.5 group"
            >
              <Plus class="w-3.5 h-3.5 text-zinc-400 group-hover:text-brand-600 dark:group-hover:text-brand-400" />
              <span class="text-[11px]">Add task...</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, nextTick } from 'vue'
import { useProjectStore } from '../stores/projectStore'
import { useAuthStore } from '../stores/authStore'
import type { TaskStatus } from '../types/project'
import TaskCard from './TaskCard.vue'
import { Plus, Lock } from 'lucide-vue-next'

const store = useProjectStore()
const authStore = useAuthStore()
const dragOverColumn = ref<string | null>(null)

const activeQuickAddColumn = ref<TaskStatus | null>(null)
const quickAddTitle = ref('')
const quickAddInputRef = ref<HTMLTextAreaElement | null>(null)

const canMove = computed(() => authStore.can('task:move'))

interface ColumnConfig {
  name: TaskStatus
  dotClass: string
  badgeClass: string
}

const columns: ColumnConfig[] = [
  {
    name: 'Backlog',
    dotClass: 'bg-zinc-400',
    badgeClass: 'bg-zinc-200/80 dark:bg-zinc-800 text-zinc-600 dark:text-zinc-400'
  },
  {
    name: 'In Progress',
    dotClass: 'bg-amber-500',
    badgeClass: 'bg-amber-100/80 dark:bg-amber-950/60 text-amber-700 dark:text-amber-400'
  },
  {
    name: 'Review',
    dotClass: 'bg-indigo-500',
    badgeClass: 'bg-indigo-100/80 dark:bg-indigo-950/60 text-indigo-700 dark:text-indigo-400'
  },
  {
    name: 'Done',
    dotClass: 'bg-emerald-500',
    badgeClass: 'bg-emerald-100/80 dark:bg-emerald-950/60 text-emerald-700 dark:text-emerald-400'
  }
]

function onDragOver(colName: string) {
  if (!canMove.value) return
  dragOverColumn.value = colName
}

function onDragLeave() {
  dragOverColumn.value = null
}

function onDrop(colName: TaskStatus, e: DragEvent) {
  if (!canMove.value) return
  dragOverColumn.value = null
  const taskId = e.dataTransfer?.getData('text/plain')
  if (taskId) {
    store.updateTaskStatus(taskId, colName)
  }
}

function openQuickAdd(colName: TaskStatus) {
  activeQuickAddColumn.value = colName
  quickAddTitle.value = ''
  nextTick(() => {
    quickAddInputRef.value?.focus()
  })
}

function cancelQuickAdd() {
  activeQuickAddColumn.value = null
  quickAddTitle.value = ''
}

function submitQuickAdd(colName: TaskStatus) {
  if (!quickAddTitle.value.trim()) {
    cancelQuickAdd()
    return
  }
  store.quickAddTask(colName, quickAddTitle.value.trim())
  quickAddTitle.value = ''
  // keep input open for fast consecutive task creation
  nextTick(() => {
    quickAddInputRef.value?.focus()
  })
}
</script>
