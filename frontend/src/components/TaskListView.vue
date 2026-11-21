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
        <!-- Table Header -->
        <div class="grid grid-cols-12 gap-2 px-3.5 py-2 text-[10px] font-mono uppercase text-zinc-400 border-b border-zinc-100 dark:border-zinc-800/80 bg-zinc-50/50 dark:bg-zinc-900/20">
          <div class="col-span-5 sm:col-span-4 md:col-span-5">Task Title</div>
          <div class="col-span-2 hidden md:block">Priority (Inline)</div>
          <div class="col-span-3 sm:col-span-3 md:col-span-2">Assignee</div>
          <div class="col-span-2 hidden sm:block">Status (Inline)</div>
          <div class="col-span-4 sm:col-span-3 md:col-span-1 text-right">Due Date</div>
        </div>

        <!-- Task Rows -->
        <div class="divide-y divide-zinc-100 dark:divide-zinc-800/80">
          <div
            v-for="task in store.tasksByStatus[status]"
            :key="task.id"
            class="min-h-11 py-2 px-3.5 grid grid-cols-12 gap-2 items-center hover:bg-zinc-50/80 dark:hover:bg-zinc-800/40 transition-colors group text-xs"
          >
            <!-- Checkbox + Title & Key -->
            <div class="col-span-5 sm:col-span-4 md:col-span-5 flex items-center space-x-2.5 min-w-0 pr-2">
              <input
                type="checkbox"
                :checked="task.status === 'Done'"
                @change="toggleTaskDone(task)"
                :disabled="!canEdit"
                class="w-3.5 h-3.5 rounded border-zinc-300 dark:border-zinc-700 text-brand-600 focus:ring-brand-500 cursor-pointer disabled:opacity-50"
              />
              <span class="font-mono text-[11px] text-zinc-400 font-medium flex-shrink-0">{{ task.key }}</span>
              <span
                @click="store.openTaskDetail(task.id)"
                :class="[
                  'truncate font-medium text-zinc-900 dark:text-zinc-100 cursor-pointer hover:text-brand-600 dark:hover:text-brand-400 transition-colors',
                  task.status === 'Done' ? 'line-through text-zinc-400 dark:text-zinc-500' : ''
                ]"
                :title="task.title"
              >
                {{ task.title }}
              </span>
            </div>

            <!-- Priority Inline Select -->
            <div class="col-span-2 hidden md:flex items-center">
              <select
                :value="task.priority"
                @change="onPriorityChange(task, ($event.target as HTMLSelectElement).value)"
                :disabled="!canEdit"
                :class="[
                  'text-[10px] font-mono px-2 py-0.5 rounded font-semibold cursor-pointer border focus:outline-none transition-colors',
                  priorityBadges[task.priority] || 'bg-zinc-100 text-zinc-600 border-zinc-200'
                ]"
              >
                <option value="Urgent">Urgent</option>
                <option value="High">High</option>
                <option value="Medium">Medium</option>
                <option value="Low">Low</option>
              </select>
            </div>

            <!-- Assignee -->
            <div class="col-span-3 sm:col-span-3 md:col-span-2 flex items-center space-x-1.5 truncate">
              <img
                :src="task.assignee?.avatar_url || 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=100&auto=format&fit=crop&q=80'"
                :alt="task.assignee?.name || 'Unassigned'"
                class="w-5 h-5 rounded-full object-cover flex-shrink-0"
              />
              <span class="text-[11px] text-zinc-600 dark:text-zinc-400 truncate">
                {{ task.assignee?.name || 'Unassigned' }}
              </span>
            </div>

            <!-- Status Inline Select -->
            <div class="col-span-2 hidden sm:flex items-center">
              <select
                :value="task.status"
                @change="onStatusChange(task, ($event.target as HTMLSelectElement).value)"
                :disabled="!canEdit"
                class="bg-zinc-100 dark:bg-zinc-800 text-zinc-800 dark:text-zinc-200 text-[10px] font-medium px-2 py-0.5 rounded border border-zinc-200 dark:border-zinc-700 cursor-pointer focus:outline-none"
              >
                <option value="Backlog">Backlog</option>
                <option value="In Progress">In Progress</option>
                <option value="Review">Review</option>
                <option value="Done">Done</option>
              </select>
            </div>

            <!-- Due Date & Action -->
            <div class="col-span-4 sm:col-span-3 md:col-span-1 flex items-center justify-end space-x-1 font-mono text-[11px] text-zinc-500 tabular-nums">
              <span>{{ task.due_date ? task.due_date.slice(5) : '-' }}</span>
              <button
                @click="store.openTaskDetail(task.id)"
                class="opacity-0 group-hover:opacity-100 p-1 text-zinc-400 hover:text-zinc-700 dark:hover:text-zinc-200 transition-opacity"
                title="Buka Detail"
              >
                <ChevronRight class="w-3.5 h-3.5" />
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useProjectStore } from '../stores/projectStore'
import { useAuthStore } from '../stores/authStore'
import type { Task, TaskPriority, TaskStatus } from '../types/project'
import { ChevronDown, ChevronRight, Plus } from 'lucide-vue-next'

const store = useProjectStore()
const authStore = useAuthStore()

const canEdit = computed(() => authStore.can('task:edit'))

const statusList: TaskStatus[] = ['Backlog', 'In Progress', 'Review', 'Done']

const statusDots: Record<TaskStatus, string> = {
  'Backlog': 'bg-zinc-400',
  'In Progress': 'bg-blue-500',
  'Review': 'bg-amber-500',
  'Done': 'bg-emerald-500',
}

const priorityBadges: Record<string, string> = {
  'Urgent': 'bg-rose-50 text-rose-600 dark:bg-rose-950/60 dark:text-rose-400 border-rose-200 dark:border-rose-900',
  'High': 'bg-orange-50 text-orange-600 dark:bg-orange-950/60 dark:text-orange-400 border-orange-200 dark:border-orange-900',
  'Medium': 'bg-blue-50 text-blue-600 dark:bg-blue-950/60 dark:text-blue-400 border-blue-200 dark:border-blue-900',
  'Low': 'bg-zinc-50 text-zinc-600 dark:bg-zinc-800 dark:text-zinc-400 border-zinc-200 dark:border-zinc-700',
}

function openGroupQuickAdd(status: TaskStatus) {
  store.quickAddTask(status, 'New task from list')
}

function toggleTaskDone(task: Task) {
  const newStatus: TaskStatus = task.status === 'Done' ? 'In Progress' : 'Done'
  store.updateTaskStatus(task.id, newStatus)
}

function onStatusChange(task: Task, newStatus: string) {
  store.updateTaskStatus(task.id, newStatus as TaskStatus)
}

function onPriorityChange(task: Task, newPriority: string) {
  store.updateTask(task.id, { priority: newPriority as TaskPriority })
}
</script>
