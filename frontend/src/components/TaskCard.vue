<template>
  <div
    :draggable="canMove"
    @dragstart="onDragStart"
    @dragend="onDragEnd"
    @click="store.openTaskDetail(task.id)"
    :class="[
      'task-card group p-3 rounded-xl bg-white dark:bg-[#18181B] border border-[#E4E4E7] dark:border-[#27272A] shadow-card hover:shadow-card-hover transition-all duration-150 relative select-none',
      isDragging ? 'opacity-40 scale-[0.98] border-dashed border-brand-500' : '',
      canMove ? 'cursor-grab active:cursor-grabbing hover:border-zinc-300 dark:hover:border-zinc-700' : 'cursor-pointer hover:border-zinc-300 dark:hover:border-zinc-700',
      task.status === 'Done' ? 'opacity-85' : ''
    ]"
  >
    <!-- Point 1: Subtle Task ID + Minimalist Priority Indicator (No loud pills) -->
    <div class="flex items-center justify-between pb-1.5">
      <div class="flex items-center space-x-1.5">
        <!-- Minimalist Priority Icon -->
        <component
          :is="priorityIcon[task.priority]"
          class="w-3.5 h-3.5"
          :class="priorityColor[task.priority]"
          :title="'Priority: ' + task.priority"
        />
        <span class="text-[11px] font-mono text-zinc-400 group-hover:text-zinc-600 dark:group-hover:text-zinc-300 font-medium">
          {{ task.key }}
        </span>
      </div>

      <!-- Category Mini Tag (Subtle muted text) -->
      <span class="text-[10px] font-mono text-zinc-400">
        {{ task.category }}
      </span>
    </div>

    <!-- Point 2: Crisp 2-line max Task Title -->
    <h3
      :class="[
        'text-xs font-medium leading-snug line-clamp-2 text-zinc-900 dark:text-zinc-100',
        task.status === 'Done' ? 'line-through text-zinc-400 dark:text-zinc-500' : ''
      ]"
    >
      {{ task.title }}
    </h3>

    <!-- Point 3: Minimal Footer (Due Date, Subtasks, Comments, 24px Avatar) -->
    <div class="pt-2.5 mt-1 border-t border-zinc-100 dark:border-zinc-800/80 flex items-center justify-between text-[11px] text-zinc-400 font-mono">
      <div class="flex items-center space-x-2.5">
        <!-- Subtask Progress -->
        <span
          v-if="task.subtasks.length > 0"
          class="flex items-center gap-1 tabular-nums"
          :class="completedSubtasksCount === task.subtasks.length ? 'text-emerald-600 dark:text-emerald-400 font-medium' : 'text-zinc-400'"
        >
          <CheckSquare class="w-3 h-3" />
          {{ completedSubtasksCount }}/{{ task.subtasks.length }}
        </span>

        <!-- Comments Counter (Only if > 0) -->
        <span v-if="task.comments.length > 0" class="flex items-center gap-1 text-zinc-400 tabular-nums">
          <MessageSquare class="w-3 h-3" />
          {{ task.comments.length }}
        </span>

        <!-- Due Date (Rose only if overdue) -->
        <span
          class="flex items-center gap-1 font-medium tabular-nums text-[10px]"
          :class="isOverdue ? 'text-rose-500 font-semibold' : 'text-zinc-400'"
        >
          <Calendar class="w-3 h-3" />
          {{ formatDueDate(task.due_date) }}
        </span>
      </div>

      <!-- 24px Assignee Avatar -->
      <img
        class="w-6 h-6 rounded-full object-cover ring-1 ring-zinc-200 dark:ring-zinc-700 flex-shrink-0"
        :src="task.assignee.avatar_url"
        :alt="task.assignee.name"
        :title="task.assignee.name"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import type { Task, TaskPriority } from '../types/project'
import { useProjectStore } from '../stores/projectStore'
import { useAuthStore } from '../stores/authStore'
import {
  CheckSquare,
  MessageSquare,
  Calendar,
  AlertCircle,
  Flag,
  ArrowUp,
  Minus
} from 'lucide-vue-next'

const props = defineProps<{
  task: Task
}>()

const store = useProjectStore()
const authStore = useAuthStore()
const isDragging = ref(false)

const canMove = computed(() => authStore.can('task:move'))

const priorityIcon: Record<TaskPriority, any> = {
  Urgent: AlertCircle,
  High: Flag,
  Medium: ArrowUp,
  Low: Minus
}

const priorityColor: Record<TaskPriority, string> = {
  Urgent: 'text-rose-500 fill-rose-500/10 stroke-[2.5]',
  High: 'text-amber-500 fill-amber-500/10 stroke-[2]',
  Medium: 'text-blue-500 stroke-[2]',
  Low: 'text-zinc-400 stroke-[2]'
}

const completedSubtasksCount = computed(() => {
  return props.task.subtasks.filter((s) => s.is_done).length
})

const isOverdue = computed(() => {
  return props.task.priority === 'Urgent' || props.task.due_date.includes('22')
})

function formatDueDate(dateStr: string) {
  if (dateStr.includes('-')) {
    const parts = dateStr.split('-')
    const month = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'][parseInt(parts[1]) - 1]
    return `${month} ${parts[2]}`
  }
  return dateStr
}

function onDragStart(e: DragEvent) {
  if (!canMove.value) {
    e.preventDefault()
    return
  }
  isDragging.value = true
  if (e.dataTransfer) {
    e.dataTransfer.setData('text/plain', props.task.id)
    e.dataTransfer.effectAllowed = 'move'
  }
}

function onDragEnd() {
  isDragging.value = false
}
</script>
