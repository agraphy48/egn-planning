<template>
  <div
    v-if="store.isNewTaskModalOpen"
    class="fixed inset-0 bg-black/40 backdrop-blur-[2px] z-50 flex items-center justify-center p-4 select-none"
    @click.self="store.isNewTaskModalOpen = false"
  >
    <div class="w-full max-w-md bg-white dark:bg-[#18181B] border border-[#E4E4E7] dark:border-[#27272A] rounded-xl shadow-palette p-5 space-y-4 animate-in fade-in zoom-in-95 duration-150">
      <div class="flex items-center justify-between border-b border-[#E4E4E7] dark:border-[#27272A] pb-3">
        <h3 class="text-sm font-bold text-zinc-900 dark:text-white font-display">New Issue / Task</h3>
        <button
          @click="store.isNewTaskModalOpen = false"
          class="text-zinc-400 hover:text-zinc-700 dark:hover:text-zinc-200"
        >
          <X class="w-4 h-4" />
        </button>
      </div>

      <div class="space-y-3 text-xs">
        <div>
          <label class="font-semibold text-zinc-700 dark:text-zinc-300">Task Title</label>
          <input
            type="text"
            v-model="title"
            placeholder="e.g. Implement Apple Pay sandbox checkout flow"
            @keydown.enter="submitTask"
            class="w-full mt-1 px-3 py-2 rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-900 text-zinc-900 dark:text-white focus:outline-none focus:ring-1 focus:ring-brand-500"
            autofocus
          />
        </div>

        <div class="grid grid-cols-2 gap-2">
          <div>
            <label class="font-semibold text-zinc-700 dark:text-zinc-300">Status</label>
            <select
              v-model="column"
              class="w-full mt-1 px-2.5 py-1.5 rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-900 text-zinc-900 dark:text-white focus:outline-none"
            >
              <option value="Backlog">Backlog</option>
              <option value="In Progress">In Progress</option>
              <option value="Review">Review</option>
              <option value="Done">Done</option>
            </select>
          </div>

          <div>
            <label class="font-semibold text-zinc-700 dark:text-zinc-300">Priority</label>
            <select
              v-model="priority"
              class="w-full mt-1 px-2.5 py-1.5 rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-900 text-zinc-900 dark:text-white focus:outline-none"
            >
              <option value="Urgent">🔴 Urgent (P0)</option>
              <option value="High">🟠 High (P1)</option>
              <option value="Medium">🔵 Medium (P2)</option>
              <option value="Low">⚪ Low (P3)</option>
            </select>
          </div>
        </div>

        <div class="grid grid-cols-2 gap-2">
          <div>
            <label class="font-semibold text-zinc-700 dark:text-zinc-300">Category Tag</label>
            <input
              type="text"
              v-model="category"
              class="w-full mt-1 px-3 py-1.5 rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-900 text-zinc-900 dark:text-white focus:outline-none"
            />
          </div>
          <div>
            <label class="font-semibold text-zinc-700 dark:text-zinc-300">Assignee</label>
            <select
              v-model="assigneeId"
              class="w-full mt-1 px-2.5 py-1.5 rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-900 text-zinc-900 dark:text-white focus:outline-none"
            >
              <option v-for="member in teamMembers" :key="member.id" :value="member.id">
                {{ member.name }}
              </option>
            </select>
          </div>
        </div>
      </div>

      <div class="flex items-center justify-end space-x-2 pt-2 border-t border-[#E4E4E7] dark:border-[#27272A]">
        <button
          @click="store.isNewTaskModalOpen = false"
          class="px-3 py-1.5 rounded-lg border border-zinc-200 dark:border-zinc-700 text-xs font-semibold text-zinc-600 dark:text-zinc-300 hover:bg-zinc-100 dark:hover:bg-zinc-800"
        >
          Cancel
        </button>
        <button
          @click="submitTask"
          class="px-4 py-1.5 rounded-lg bg-brand-600 hover:bg-brand-500 text-white text-xs font-semibold shadow-2xs"
        >
          Create Task
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useProjectStore } from '../stores/projectStore'
import { useAuthStore } from '../stores/authStore'
import { teamMembers } from '../data/mockData'
import type { TaskPriority, TaskStatus } from '../types/project'
import { X } from 'lucide-vue-next'

const store = useProjectStore()
const authStore = useAuthStore()

const title = ref('')
const column = ref<TaskStatus>('Backlog')
const priority = ref<TaskPriority>('High')
const category = ref('Frontend')
const assigneeId = ref(teamMembers[0]?.id || 'usr_agus')

function submitTask() {
  if (!authStore.can('task:create')) {
    store.showToast('Forbidden: Your role cannot create tasks')
    store.isNewTaskModalOpen = false
    return
  }
  if (!title.value.trim()) return
  const selectedAssignee = teamMembers.find(m => m.id === assigneeId.value) || teamMembers[0]
  store.createTask({
    title: title.value.trim(),
    status: column.value,
    priority: priority.value,
    category: category.value || 'General',
    assignee: selectedAssignee,
  })
  title.value = ''
}
</script>
