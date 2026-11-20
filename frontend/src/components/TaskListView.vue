<template>
  <div class="space-y-4 select-none pb-6">
    <!-- Grouped Status Sections -->
    <div
      v-for="status in statusList"
      :key="status"
      class="rounded-xl border border-[#E4E4E7] dark:border-[#27272A] bg-white dark:bg-[#18181B] overflow-hidden shadow-subtle"
    >
      <!-- Group Header (Collapsible) -->
      <div
        @click="store.toggleGroupCollapse(status)"
        class="px-3.5 py-2.5 bg-[#F4F4F5]/60 dark:bg-[#27272A]/40 border-b border-[#E4E4E7] dark:border-[#27272A] flex items-center justify-between cursor-pointer hover:bg-[#F4F4F5] dark:hover:bg-[#27272A]/70 transition-colors"
      >
        <div class="flex items-center space-x-2">
          <ChevronDown
            class="w-3.5 h-3.5 text-zinc-400 transition-transform duration-200"
            :class="{ '-rotate-90': store.collapsedGroups.includes(status) }"
          />
          <span class="w-2 h-2 rounded-full" :class="statusDots[status]"></span>
          <span class="text-xs font-semibold text-zinc-900 dark:text-zinc-100 tracking-tight">{{ status }}</span>
          <span class="text-[10px] font-mono px-1.5 py-0.2 rounded-full bg-zinc-200/70 dark:bg-zinc-800 text-zinc-600 dark:text-zinc-400 font-medium tabular-nums">
            {{ store.tasksByStatus[status]?.length || 0 }}
          </span>
        </div>

        <button
          v-can="'task:create'"
          @click.stop="openGroupQuickAdd(status)"
          class="p-1 rounded text-zinc-400 hover:text-zinc-700 dark:hover:text-zinc-200 hover:bg-zinc-200/60 dark:hover:bg-zinc-800 transition-colors text-[11px] flex items-center gap-1"
        >
          <Plus class="w-3.5 h-3.5" />
          <span class="hidden sm:inline">Add</span>
        </button>
      </div>

      <!-- Group Rows Container (Shown when not collapsed) -->
      <div v-show="!store.collapsedGroups.includes(status)">
        <!-- Table-like Header row -->
        <div class="grid grid-cols-12 gap-2 px-3.5 py-2 text-[10px] font-mono uppercase text-zinc-400 border-b border-zinc-100 dark:border-zinc-800/80 bg-zinc-50/50 dark:bg-zinc-900/20">
          <div class="col-span-6 sm:col-span-5 md:col-span-6">Task Title</div>
          <div class="col-span-2 hidden md:block">Priority</div>
          <div class="col-span-3 sm:col-span-3 md:col-span-2">Assignee</div>
          <div class="col-span-3 sm:col-span-4 md:col-span-2 text-right">Due Date</div>
        </div>

        <!-- Task Rows -->
        <div class="divide-y divide-zinc-100 dark:divide-zinc-800/80">
          <div
            v-for="task in store.tasksByStatus[status]"
            :key="task.id"
            @click="store.openTaskDetail(task.id)"
            class="h-11 px-3.5 grid grid-cols-12 gap-2 items-center hover:bg-zinc-50/80 dark:hover:bg-zinc-800/40 cursor-pointer transition-colors group text-xs"
          >
            <!-- Title & Key -->
            <div class="col-span-6 sm:col-span-5 md:col-span-6 flex items-center space-x-2.5 min-w-0 pr-2">
              <span class="font-mono text-[11px] text-zinc-400 font-medium flex-shrink-0">{{ task.key }}</span>
              <span
                :class="[
                  'truncate font-medium text-zinc-900 dark:text-zinc-100',
                  task.status === 'Done' ? 'line-through text-zinc-400 dark:text-zinc-500' : ''
                ]"
              >
                {{ task.title }}
              </span>
              <span class="hidden lg:inline text-[10px] font-mono text-zinc-400 flex-shrink-0">
                {{ task.category }}
              </span>
            </div>

            <!-- Inline Priority Selector -->
            <div class="col-span-2 hidden md:flex items-center" @click.stop>
              <select
                :value="task.priority"
                @change="handlePriorityChange(task.id, ($event.target as HTMLSelectElement).value as TaskPriority)"
                class="bg-transparent text-[11px] font-medium text-zinc-700 dark:text-zinc-300 focus:outline-none cursor-pointer py-1 pr-2 rounded hover:bg-zinc-100 dark:hover:bg-zinc-800 transition-colors"
              >
                <option value="Urgent">🔴 Urgent</option>
                <option value="High">🟠 High</option>
                <option value="Medium">🔵 Medium</option>
                <option value="Low">⚪ Low</option>
              </select>
            </div>

            <!-- Inline Assignee Avatar & Name -->
            <div class="col-span-3 sm:col-span-3 md:col-span-2 flex items-center space-x-2 truncate">
              <img
                class="w-5 h-5 rounded-full object-cover ring-1 ring-zinc-200 dark:ring-zinc-700 flex-shrink-0"
                :src="task.assignee.avatar_url"
                :alt="task.assignee.name"
              />
              <span class="text-zinc-600 dark:text-zinc-300 truncate text-[11px]">{{ task.assignee.name.split(' ')[0] }}</span>
            </div>

            <!-- Due Date & Hover Quick Action Controls -->
            <div class="col-span-3 sm:col-span-4 md:col-span-2 flex items-center justify-end space-x-2">
              <!-- Regular Due Date view -->
              <span
                class="group-hover:hidden text-[11px] font-mono tabular-nums"
                :class="task.priority === 'Urgent' ? 'text-rose-500 font-medium' : 'text-zinc-400'"
              >
                {{ formatDueDate(task.due_date) }}
              </span>

              <!-- Hover-revealed quick actions -->
              <div class="hidden group-hover:flex items-center space-x-1" @click.stop>
                <button
                  @click="copyTaskLink(task.id)"
                  class="p-1 rounded text-zinc-400 hover:text-zinc-700 dark:hover:text-zinc-200 hover:bg-zinc-200/60 dark:hover:bg-zinc-700"
                  title="Copy link"
                >
                  <Link class="w-3.5 h-3.5" />
                </button>
                <button
                  v-can="'task:delete'"
                  @click="handleDelete(task.id, task.key)"
                  class="p-1 rounded text-zinc-400 hover:text-rose-600 hover:bg-rose-50 dark:hover:bg-rose-950/40"
                  title="Delete task"
                >
                  <Trash2 class="w-3.5 h-3.5" />
                </button>
                <button
                  @click="store.openTaskDetail(task.id)"
                  class="p-1 rounded text-zinc-400 hover:text-brand-600 hover:bg-brand-50 dark:hover:bg-brand-950/40"
                  title="Open details"
                >
                  <Maximize2 class="w-3.5 h-3.5" />
                </button>
              </div>
            </div>
          </div>

          <!-- Empty state row -->
          <div
            v-if="!store.tasksByStatus[status] || store.tasksByStatus[status].length === 0"
            class="py-4 text-center text-xs text-zinc-400 font-mono"
          >
            No tasks in {{ status }}
          </div>
        </div>

        <!-- Quick Add Row at bottom of list group -->
        <div v-can="'task:create'" class="p-2 border-t border-zinc-100 dark:border-zinc-800/80 bg-zinc-50/30 dark:bg-zinc-900/10">
          <div v-if="activeQuickAddGroup === status" class="flex items-center space-x-2 px-2">
            <input
              ref="groupInputRef"
              type="text"
              v-model="groupQuickAddTitle"
              :placeholder="'Add new task in ' + status + '... (Enter to save)'"
              @keydown.enter="submitGroupQuickAdd(status)"
              @keydown.esc="cancelGroupQuickAdd"
              class="flex-1 text-xs bg-white dark:bg-zinc-900 border border-brand-500 rounded-lg px-2.5 py-1.5 focus:outline-none text-zinc-900 dark:text-zinc-100 shadow-2xs"
            />
            <button
              @click="submitGroupQuickAdd(status)"
              class="px-2.5 py-1.5 rounded-lg bg-brand-600 hover:bg-brand-500 text-white text-xs font-semibold"
            >
              Add
            </button>
            <button
              @click="cancelGroupQuickAdd"
              class="px-2 py-1.5 text-xs text-zinc-400 hover:text-zinc-600 dark:hover:text-zinc-200"
            >
              Cancel
            </button>
          </div>
          <button
            v-else
            @click="openGroupQuickAdd(status)"
            class="w-full text-left px-2 py-1 text-xs text-zinc-400 hover:text-zinc-700 dark:hover:text-zinc-200 hover:bg-zinc-100/60 dark:hover:bg-zinc-800/50 rounded flex items-center space-x-1.5 transition-colors"
          >
            <Plus class="w-3.5 h-3.5 text-zinc-400" />
            <span class="text-[11px]">+ Add task to {{ status }}...</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, nextTick } from 'vue'
import { useProjectStore } from '../stores/projectStore'
import type { TaskStatus, TaskPriority } from '../types/project'
import {
  ChevronDown,
  Plus,
  Link,
  Trash2,
  Maximize2
} from 'lucide-vue-next'

const store = useProjectStore()

const statusList: TaskStatus[] = ['Backlog', 'In Progress', 'Review', 'Done']

const statusDots: Record<TaskStatus, string> = {
  'Backlog': 'bg-zinc-400',
  'In Progress': 'bg-amber-500',
  'Review': 'bg-indigo-500',
  'Done': 'bg-emerald-500'
}

const activeQuickAddGroup = ref<TaskStatus | null>(null)
const groupQuickAddTitle = ref('')
const groupInputRef = ref<HTMLInputElement | null>(null)

function openGroupQuickAdd(status: TaskStatus) {
  activeQuickAddGroup.value = status
  groupQuickAddTitle.value = ''
  nextTick(() => {
    groupInputRef.value?.focus()
  })
}

function cancelGroupQuickAdd() {
  activeQuickAddGroup.value = null
  groupQuickAddTitle.value = ''
}

function submitGroupQuickAdd(status: TaskStatus) {
  if (!groupQuickAddTitle.value.trim()) {
    cancelGroupQuickAdd()
    return
  }
  store.quickAddTask(status, groupQuickAddTitle.value.trim())
  groupQuickAddTitle.value = ''
  nextTick(() => {
    groupInputRef.value?.focus()
  })
}

function handlePriorityChange(taskId: string, priority: TaskPriority) {
  store.updateTask(taskId, { priority })
}

function formatDueDate(dateStr: string) {
  if (dateStr.includes('-')) {
    const parts = dateStr.split('-')
    const month = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'][parseInt(parts[1]) - 1]
    return `${month} ${parts[2]}`
  }
  return dateStr
}

function copyTaskLink(taskId: string) {
  navigator.clipboard?.writeText(window.location.origin + `/projects?task=${taskId}`)
  store.showToast('Task link copied to clipboard')
}

function handleDelete(taskId: string, taskKey: string) {
  if (confirm(`Are you sure you want to delete ${taskKey}?`)) {
    store.deleteTask(taskId)
  }
}
</script>
