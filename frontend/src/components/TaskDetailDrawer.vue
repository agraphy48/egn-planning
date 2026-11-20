<template>
  <div>
    <!-- Backdrop Overlay with Fade Transition -->
    <Transition name="fade">
      <div
        v-if="store.isDrawerOpen"
        @click="store.closeTaskDetail()"
        class="fixed inset-0 bg-black/40 backdrop-blur-[2px] z-40 transition-opacity"
      ></div>
    </Transition>

    <!-- Drawer Panel with Slide-over Transition -->
    <Transition name="drawer">
      <aside
        v-if="store.isDrawerOpen && task"
        class="fixed top-0 right-0 bottom-0 w-full max-w-2xl bg-white dark:bg-[#18181B] border-l border-[#E4E4E7] dark:border-[#27272A] shadow-drawer z-50 flex flex-col select-none"
      >
        <!-- Top Action Bar -->
        <div class="h-13 border-b border-[#E4E4E7] dark:border-[#27272A] px-5 flex items-center justify-between flex-shrink-0 bg-zinc-50/50 dark:bg-zinc-900/20">
          <div class="flex items-center space-x-2.5 text-xs font-mono text-zinc-500 dark:text-zinc-400">
            <span class="font-bold text-zinc-800 dark:text-zinc-200">{{ task.key }}</span>
            <span>•</span>
            <span class="px-2 py-0.5 rounded bg-zinc-100 dark:bg-zinc-800 text-zinc-600 dark:text-zinc-300">
              {{ task.category }}
            </span>
            <span v-if="task.client_visible" class="text-[10px] px-1.5 py-0.5 rounded bg-indigo-50 dark:bg-indigo-950/60 text-indigo-600 dark:text-indigo-400 font-medium">
              Client Shared
            </span>
          </div>

          <div class="flex items-center space-x-1">
            <!-- Delete Task (PM & Admin only) -->
            <button
              v-can="'task:delete'"
              @click="handleDeleteTask"
              class="p-1.5 rounded-md text-zinc-400 hover:text-rose-600 hover:bg-rose-50 dark:hover:bg-rose-950/40 transition-colors"
              title="Delete Task"
            >
              <Trash2 class="w-4 h-4" />
            </button>

            <!-- Copy Link -->
            <button
              @click="copyTaskLink"
              class="p-1.5 rounded-md text-zinc-400 hover:text-zinc-700 dark:hover:text-zinc-200 hover:bg-zinc-100 dark:hover:bg-zinc-800 transition-colors"
              title="Copy Task Link"
            >
              <Link class="w-4 h-4" />
            </button>

            <!-- Close Drawer -->
            <button
              @click="store.closeTaskDetail()"
              class="p-1.5 rounded-md text-zinc-400 hover:text-zinc-700 dark:hover:text-zinc-200 hover:bg-zinc-100 dark:hover:bg-zinc-800 transition-colors"
              title="Close (Esc)"
            >
              <X class="w-4 h-4" />
            </button>
          </div>
        </div>

        <!-- Scrollable Content -->
        <div class="flex-1 overflow-y-auto p-6 space-y-6">
          <!-- Large Editable Task Title -->
          <div>
            <input
              type="text"
              v-model="task.title"
              :disabled="!canEdit"
              :readonly="!canEdit"
              placeholder="Task Title..."
              class="w-full text-xl font-semibold tracking-tight text-zinc-900 dark:text-white bg-transparent border-b border-transparent hover:border-zinc-200 dark:hover:border-zinc-700 focus:border-brand-500 focus:outline-none transition-colors pb-1"
            />
          </div>

          <!-- Linear-Style Compact Properties Strip (2-Column Grid) -->
          <div class="grid grid-cols-2 gap-3 text-xs bg-[#F4F4F5]/60 dark:bg-[#27272A]/40 p-3.5 rounded-xl border border-[#E4E4E7] dark:border-[#27272A]">
            <!-- Status Dropdown -->
            <div class="space-y-1">
              <label class="text-[10px] font-mono uppercase tracking-wider text-zinc-400 font-semibold">Status</label>
              <div class="relative">
                <select
                  v-model="task.status"
                  :disabled="!canEdit"
                  @change="store.updateTaskStatus(task.id, task.status)"
                  class="w-full bg-white dark:bg-[#18181B] border border-zinc-200 dark:border-zinc-700 rounded-lg px-2.5 py-1.5 text-zinc-800 dark:text-zinc-200 font-medium focus:ring-1 focus:ring-brand-500 focus:outline-none"
                  :class="!canEdit ? 'opacity-80 cursor-not-allowed' : 'cursor-pointer'"
                >
                  <option value="Backlog">Backlog</option>
                  <option value="In Progress">In Progress</option>
                  <option value="Review">Review</option>
                  <option value="Done">Done</option>
                </select>
              </div>
            </div>

            <!-- Priority Dropdown -->
            <div class="space-y-1">
              <label class="text-[10px] font-mono uppercase tracking-wider text-zinc-400 font-semibold">Priority</label>
              <select
                v-model="task.priority"
                :disabled="!canEdit"
                class="w-full bg-white dark:bg-[#18181B] border border-zinc-200 dark:border-zinc-700 rounded-lg px-2.5 py-1.5 text-zinc-800 dark:text-zinc-200 font-medium focus:ring-1 focus:ring-brand-500 focus:outline-none"
                :class="!canEdit ? 'opacity-80 cursor-not-allowed' : 'cursor-pointer'"
              >
                <option value="Urgent">🔴 Urgent (P0)</option>
                <option value="High">🟠 High (P1)</option>
                <option value="Medium">🔵 Medium (P2)</option>
                <option value="Low">⚪ Low (P3)</option>
              </select>
            </div>

            <!-- Assignee Selector -->
            <div class="space-y-1">
              <label class="text-[10px] font-mono uppercase tracking-wider text-zinc-400 font-semibold">Assignee</label>
              <div class="flex items-center space-x-2 bg-white dark:bg-[#18181B] border border-zinc-200 dark:border-zinc-700 rounded-lg px-2.5 py-1.5">
                <img class="w-5 h-5 rounded-full object-cover" :src="task.assignee.avatar_url" :alt="task.assignee.name" />
                <span class="font-medium text-zinc-800 dark:text-zinc-200 truncate">{{ task.assignee.name }}</span>
              </div>
            </div>

            <!-- Due Date Picker -->
            <div class="space-y-1">
              <label class="text-[10px] font-mono uppercase tracking-wider text-zinc-400 font-semibold">Due Date</label>
              <div class="flex items-center space-x-2 bg-white dark:bg-[#18181B] border border-zinc-200 dark:border-zinc-700 rounded-lg px-2.5 py-1.5">
                <Calendar class="w-3.5 h-3.5 text-zinc-400" />
                <input
                  type="date"
                  v-model="task.due_date"
                  :disabled="!canEdit"
                  :readonly="!canEdit"
                  class="bg-transparent text-zinc-800 dark:text-zinc-200 focus:outline-none w-full text-xs font-mono"
                  :class="!canEdit ? 'cursor-not-allowed' : ''"
                />
              </div>
            </div>
          </div>

          <!-- Deliverable & Epic Module Info -->
          <div v-if="task.fase || task.deliverable" class="p-3 rounded-xl bg-brand-50/40 dark:bg-brand-950/20 border border-brand-100/60 dark:border-brand-900/40 space-y-1 text-xs">
            <div class="flex items-center justify-between text-[11px] font-mono text-zinc-500 dark:text-zinc-400">
              <span class="font-bold text-brand-700 dark:text-brand-300">{{ task.fase }} • {{ task.epic_modul }}</span>
              <span v-if="task.time_estimate" class="text-zinc-500 font-mono">
                Est: {{ task.time_estimate }}
              </span>
            </div>
            <div v-if="task.deliverable" class="text-zinc-700 dark:text-zinc-300 font-medium text-[11px]">
              <span class="text-zinc-400 font-mono">Deliverable:</span> {{ task.deliverable }}
            </div>
          </div>

          <!-- Description Area -->
          <div class="space-y-2">
            <div class="flex items-center justify-between">
              <label class="text-xs font-bold text-zinc-900 dark:text-white font-display">Description</label>
              <div v-if="canEdit" class="flex items-center space-x-1 text-zinc-400">
                <button class="p-1 hover:text-zinc-700 dark:hover:text-zinc-200 rounded text-xs font-bold" title="Bold">B</button>
                <button class="p-1 hover:text-zinc-700 dark:hover:text-zinc-200 rounded text-xs italic" title="Italic">I</button>
                <button class="p-1 hover:text-zinc-700 dark:hover:text-zinc-200 rounded font-mono text-[11px]" title="Code">&lt;/&gt;</button>
              </div>
            </div>
            <textarea
              rows="4"
              v-model="task.description"
              :disabled="!canEdit"
              :readonly="!canEdit"
              placeholder="Add detailed task description, acceptance criteria, or technical specs..."
              class="w-full text-xs text-zinc-800 dark:text-zinc-200 bg-white dark:bg-zinc-900/60 border border-zinc-200 dark:border-zinc-700/80 rounded-xl p-3 leading-relaxed resize-none transition-colors"
              :class="canEdit ? 'focus:ring-1 focus:ring-brand-500 focus:outline-none' : 'cursor-not-allowed opacity-90'"
            ></textarea>
          </div>

          <!-- Interactive Sub-task Checklist -->
          <div class="space-y-3">
            <div class="flex items-center justify-between">
              <div class="flex items-center space-x-2">
                <span class="text-xs font-bold text-zinc-900 dark:text-white font-display">Sub-tasks</span>
                <span class="text-[10px] font-mono px-1.5 py-0.2 rounded bg-zinc-100 dark:bg-zinc-800 text-zinc-600 dark:text-zinc-400 font-semibold tabular-nums">
                  {{ completedSubtasksCount }} of {{ task.subtasks.length }} completed
                </span>
              </div>
              <div class="w-24 bg-zinc-200 dark:bg-zinc-800 h-1.5 rounded-full overflow-hidden">
                <div
                  class="bg-brand-600 h-full rounded-full transition-all duration-300"
                  :style="{ width: subtaskPercentage + '%' }"
                ></div>
              </div>
            </div>

            <!-- Subtask Items -->
            <div class="space-y-1.5 text-xs">
              <label
                v-for="sub in task.subtasks"
                :key="sub.id"
                :class="[
                  'flex items-center space-x-2.5 p-2 rounded-lg bg-zinc-50 dark:bg-zinc-900/40 transition-colors',
                  canSubtask ? 'hover:bg-zinc-100/70 dark:hover:bg-zinc-800/50 cursor-pointer' : 'cursor-default opacity-85'
                ]"
              >
                <input
                  type="checkbox"
                  v-model="sub.is_done"
                  :disabled="!canSubtask"
                  @change="store.toggleSubtask(task.id, sub.id, sub.is_done)"
                  class="w-4 h-4 rounded text-brand-600 focus:ring-0"
                  :class="canSubtask ? 'cursor-pointer' : 'cursor-not-allowed'"
                />
                <span :class="sub.is_done ? 'line-through text-zinc-400 dark:text-zinc-500' : 'text-zinc-800 dark:text-zinc-200'">
                  {{ sub.title }}
                </span>
              </label>
            </div>

            <!-- Add Subtask Input -->
            <div v-if="canSubtask" class="flex items-center space-x-2">
              <input
                type="text"
                v-model="newSubtaskTitle"
                placeholder="+ Add a subtask... (Enter)"
                @keydown.enter="handleAddSubtask"
                class="flex-1 text-xs bg-transparent border-b border-zinc-200 dark:border-zinc-700/80 px-2 py-1.5 focus:border-brand-500 focus:outline-none text-zinc-800 dark:text-zinc-200"
              />
              <button
                @click="handleAddSubtask"
                class="px-2 py-1 rounded bg-zinc-100 dark:bg-zinc-800 hover:bg-zinc-200 dark:hover:bg-zinc-700 text-xs text-zinc-700 dark:text-zinc-300 font-medium transition-colors"
              >
                Add
              </button>
            </div>
          </div>

          <!-- Attachments -->
          <div class="space-y-2">
            <label class="text-xs font-bold text-zinc-900 dark:text-white font-display">Attachments</label>
            <div
              v-if="canUpload"
              class="border-2 border-dashed border-zinc-200 dark:border-zinc-800 hover:border-brand-500/60 rounded-xl p-4 text-center cursor-pointer transition-colors group"
            >
              <CloudUpload class="w-6 h-6 mx-auto text-zinc-400 group-hover:text-brand-500 transition-colors" />
              <p class="text-xs text-zinc-600 dark:text-zinc-400 mt-1">
                <span class="font-semibold text-brand-600 dark:text-brand-400">Click to upload</span> or drag and drop
              </p>
              <span class="text-[10px] text-zinc-400 font-mono">PDF, PNG, SVG up to 25MB</span>
            </div>
            <div
              v-else
              class="p-2.5 rounded-lg bg-zinc-50 dark:bg-zinc-900/40 border border-zinc-100 dark:border-zinc-800 text-center text-xs text-zinc-400 font-mono"
            >
              File uploads are restricted for your current role.
            </div>

            <div v-if="task.attachments.length > 0" class="space-y-1.5 text-xs">
              <div
                v-for="att in task.attachments"
                :key="att.id"
                class="flex items-center justify-between p-2 rounded-lg bg-zinc-50 dark:bg-zinc-900/50 border border-zinc-100 dark:border-zinc-800/80"
              >
                <div class="flex items-center space-x-2 truncate">
                  <FileText class="w-4 h-4 text-brand-500 flex-shrink-0" />
                  <span class="font-medium text-zinc-800 dark:text-zinc-200 truncate">{{ att.file_name }}</span>
                  <span class="text-[10px] font-mono text-zinc-400">{{ att.file_size }}</span>
                </div>
                <button class="p-1 text-zinc-400 hover:text-zinc-700 dark:hover:text-zinc-200">
                  <Download class="w-3.5 h-3.5" />
                </button>
              </div>
            </div>
          </div>

          <!-- Discussion & Comments -->
          <div class="space-y-3 pt-2 border-t border-[#E4E4E7] dark:border-[#27272A]">
            <div class="flex items-center justify-between">
              <span class="text-xs font-bold text-zinc-900 dark:text-white font-display">Discussion</span>
              <span class="text-[10px] font-mono text-zinc-400 tabular-nums">{{ task.comments.length }} comments</span>
            </div>

            <!-- Comments Timeline -->
            <div class="space-y-2.5 text-xs">
              <div
                v-for="com in task.comments"
                :key="com.id"
                class="flex items-start space-x-2.5"
              >
                <img class="w-6 h-6 rounded-full object-cover mt-0.5 flex-shrink-0" :src="com.user.avatar_url" :alt="com.user.name" />
                <div class="flex-1 bg-zinc-50 dark:bg-zinc-900/60 p-2.5 rounded-xl border border-zinc-100 dark:border-zinc-800 space-y-1">
                  <div class="flex items-center justify-between">
                    <span class="font-semibold text-zinc-900 dark:text-white">{{ com.user.name }}</span>
                    <span class="text-[10px] font-mono text-zinc-400">{{ com.created_at }}</span>
                  </div>
                  <p class="text-zinc-700 dark:text-zinc-300 leading-relaxed">{{ com.content }}</p>
                </div>
              </div>
            </div>

            <!-- New Comment Box -->
            <div class="pt-2">
              <div class="rounded-xl border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-900 p-2 space-y-2 focus-within:border-brand-500">
                <textarea
                  rows="2"
                  v-model="newCommentText"
                  placeholder="Write a comment... use @ to mention teammates (Cmd+Enter to post)"
                  class="w-full text-xs text-zinc-800 dark:text-zinc-200 bg-transparent focus:outline-none resize-none leading-relaxed"
                  @keydown.meta.enter="handlePostComment"
                  @keydown.ctrl.enter="handlePostComment"
                ></textarea>
                <div class="flex items-center justify-between pt-1 border-t border-zinc-100 dark:border-zinc-800 text-[11px]">
                  <div class="flex items-center space-x-1.5 text-zinc-400">
                    <button class="hover:text-zinc-600 dark:hover:text-zinc-300 p-1" title="Mention">@</button>
                    <button class="hover:text-zinc-600 dark:hover:text-zinc-300 p-1"><Paperclip class="w-3.5 h-3.5" /></button>
                  </div>
                  <button
                    @click="handlePostComment"
                    class="px-3 py-1 rounded-lg bg-brand-600 hover:bg-brand-500 text-white font-semibold text-xs shadow-2xs transition-colors"
                  >
                    Send
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </aside>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useProjectStore } from '../stores/projectStore'
import { useAuthStore } from '../stores/authStore'
import {
  Link,
  X,
  Calendar,
  CloudUpload,
  FileText,
  Download,
  Paperclip,
  Trash2
} from 'lucide-vue-next'

const store = useProjectStore()
const authStore = useAuthStore()
const task = computed(() => store.selectedTask)

const canEdit = computed(() => authStore.can('task:edit'))
const canSubtask = computed(() => authStore.can('subtask:manage'))
const canUpload = computed(() => authStore.can('attachment:upload'))

const newSubtaskTitle = ref('')
const newCommentText = ref('')

const completedSubtasksCount = computed(() => {
  if (!task.value) return 0
  return task.value.subtasks.filter((s) => s.is_done).length
})

const subtaskPercentage = computed(() => {
  if (!task.value || task.value.subtasks.length === 0) return 0
  return Math.round((completedSubtasksCount.value / task.value.subtasks.length) * 100)
})

function handleAddSubtask() {
  if (!canSubtask.value || !task.value || !newSubtaskTitle.value.trim()) return
  store.addSubtask(task.value.id, newSubtaskTitle.value)
  newSubtaskTitle.value = ''
}

function handlePostComment() {
  if (!task.value || !newCommentText.value.trim()) return
  store.addComment(task.value.id, newCommentText.value)
  newCommentText.value = ''
}

function handleDeleteTask() {
  if (!task.value) return
  if (confirm(`Are you sure you want to delete task ${task.value.key}?`)) {
    store.deleteTask(task.value.id)
  }
}

function copyTaskLink() {
  if (!task.value) return
  navigator.clipboard?.writeText(window.location.origin + `/projects?task=${task.value.id}`)
  store.showToast('Task deep link copied!')
}
</script>
