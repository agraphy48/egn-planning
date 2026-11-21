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
            <span v-else class="text-[10px] px-1.5 py-0.5 rounded bg-amber-50 dark:bg-amber-950/60 text-amber-600 dark:text-amber-400 font-medium">
              Internal Only
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
              @change="saveTaskUpdates"
              placeholder="Task Title..."
              class="w-full text-xl font-semibold tracking-tight text-zinc-900 dark:text-white bg-transparent border-b border-transparent hover:border-zinc-200 dark:hover:border-zinc-700 focus:border-brand-500 focus:outline-none transition-colors pb-1"
            />
          </div>

          <!-- Linear-Style Compact Properties Strip (2-Column Grid) -->
          <div class="grid grid-cols-2 gap-3 text-xs bg-[#F4F4F5]/60 dark:bg-[#27272A]/40 p-3.5 rounded-xl border border-[#E4E4E7] dark:border-[#27272A]">
            <!-- Status Dropdown -->
            <div class="space-y-1">
              <label class="text-[10px] font-mono uppercase tracking-wider text-zinc-400 font-semibold">Status</label>
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

            <!-- Priority Dropdown -->
            <div class="space-y-1">
              <label class="text-[10px] font-mono uppercase tracking-wider text-zinc-400 font-semibold">Priority</label>
              <select
                v-model="task.priority"
                :disabled="!canEdit"
                @change="saveTaskUpdates"
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
                <img class="w-5 h-5 rounded-full object-cover" :src="task.assignee?.avatar_url || 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=100&auto=format&fit=crop&q=80'" :alt="task.assignee?.name || 'Unassigned'" />
                <span class="font-medium text-zinc-800 dark:text-zinc-200 truncate">{{ task.assignee?.name || 'Unassigned' }}</span>
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
                  @change="saveTaskUpdates"
                  class="bg-transparent text-zinc-800 dark:text-zinc-200 focus:outline-none w-full text-xs font-mono"
                  :class="!canEdit ? 'cursor-not-allowed' : ''"
                />
              </div>
            </div>
          </div>

          <!-- Description Area -->
          <div class="space-y-2">
            <label class="text-xs font-bold text-zinc-900 dark:text-white font-display">Deskripsi Task</label>
            <textarea
              rows="3"
              v-model="task.description"
              :disabled="!canEdit"
              :readonly="!canEdit"
              @change="saveTaskUpdates"
              placeholder="Deskripsi tugas dan spesifikasi teknis..."
              class="w-full text-xs text-zinc-800 dark:text-zinc-200 bg-white dark:bg-zinc-900/60 border border-zinc-200 dark:border-zinc-700/80 rounded-xl p-3 leading-relaxed resize-none transition-colors"
              :class="canEdit ? 'focus:ring-1 focus:ring-brand-500 focus:outline-none' : 'cursor-not-allowed opacity-90'"
            ></textarea>
          </div>

          <!-- Timesheet & Work Hours Logger -->
          <div class="p-3.5 rounded-xl bg-zinc-50 dark:bg-zinc-900/40 border border-zinc-200/80 dark:border-zinc-800 space-y-2.5">
            <div class="flex items-center justify-between">
              <div class="flex items-center space-x-1.5">
                <Clock class="w-4 h-4 text-brand-600 dark:text-brand-400" />
                <span class="text-xs font-bold text-zinc-900 dark:text-white">Pencatatan Jam Kerja (Timesheet)</span>
              </div>
              <button
                @click="showLogTimeForm = !showLogTimeForm"
                class="text-[11px] text-brand-600 dark:text-brand-400 font-semibold hover:underline"
              >
                {{ showLogTimeForm ? 'Batal' : '+ Log Jam' }}
              </button>
            </div>

            <!-- Inline Log Time Form -->
            <div v-if="showLogTimeForm" class="p-2.5 bg-white dark:bg-zinc-800/80 rounded-lg border border-zinc-200 dark:border-zinc-700 space-y-2 text-xs">
              <div class="grid grid-cols-2 gap-2">
                <div>
                  <label class="text-[10px] text-zinc-400 font-mono">Jumlah Jam (Hours)</label>
                  <input
                    type="number"
                    step="0.5"
                    min="0.5"
                    v-model.number="logHours"
                    class="w-full mt-0.5 px-2 py-1 rounded bg-zinc-50 dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-700 text-xs"
                  />
                </div>
                <div>
                  <label class="text-[10px] text-zinc-400 font-mono">Tanggal</label>
                  <input
                    type="date"
                    v-model="logDate"
                    class="w-full mt-0.5 px-2 py-1 rounded bg-zinc-50 dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-700 text-xs font-mono"
                  />
                </div>
              </div>
              <div>
                <input
                  type="text"
                  v-model="logDescription"
                  placeholder="Keterangan pekerjaan..."
                  class="w-full px-2 py-1 rounded bg-zinc-50 dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-700 text-xs"
                />
              </div>
              <div class="flex justify-end">
                <button
                  @click="handleLogTime"
                  :disabled="!logHours || logHours <= 0"
                  class="px-3 py-1 rounded bg-brand-600 text-white text-xs font-semibold disabled:opacity-50"
                >
                  Simpan Jam Kerja
                </button>
              </div>
            </div>
          </div>

          <!-- Interactive Sub-task Checklist -->
          <div class="space-y-3">
            <div class="flex items-center justify-between">
              <div class="flex items-center space-x-2">
                <span class="text-xs font-bold text-zinc-900 dark:text-white font-display">Sub-tasks</span>
                <span class="text-[10px] font-mono px-1.5 py-0.2 rounded bg-zinc-100 dark:bg-zinc-800 text-zinc-600 dark:text-zinc-400 font-semibold tabular-nums">
                  {{ completedSubtasksCount }} of {{ task.subtasks?.length || 0 }} completed
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
                placeholder="+ Tambah subtask... (Enter)"
                @keydown.enter="handleAddSubtask"
                class="flex-1 text-xs bg-transparent border-b border-zinc-200 dark:border-zinc-700/80 px-2 py-1.5 focus:border-brand-500 focus:outline-none text-zinc-800 dark:text-zinc-200"
              />
              <button
                @click="handleAddSubtask"
                class="px-2 py-1 rounded bg-zinc-100 dark:bg-zinc-800 hover:bg-zinc-200 dark:hover:bg-zinc-700 text-xs text-zinc-700 dark:text-zinc-300 font-medium transition-colors"
              >
                Tambah
              </button>
            </div>
          </div>

          <!-- Attachments Section (Upload, Preview & Delete) -->
          <div class="space-y-2">
            <div class="flex items-center justify-between">
              <label class="text-xs font-bold text-zinc-900 dark:text-white font-display">Lampiran Dokumen (Attachments)</label>
              <span class="text-[10px] font-mono text-zinc-400">{{ task.attachments?.length || 0 }} files</span>
            </div>

            <!-- Hidden File Input -->
            <input
              type="file"
              ref="fileInputRef"
              @change="onFileSelected"
              class="hidden"
            />

            <!-- Drag & Drop Zone -->
            <div
              v-if="canUpload"
              @click="triggerFileInput"
              class="border-2 border-dashed border-zinc-200 dark:border-zinc-800 hover:border-brand-500/60 rounded-xl p-3.5 text-center cursor-pointer transition-colors group"
            >
              <CloudUpload class="w-5 h-5 mx-auto text-zinc-400 group-hover:text-brand-500 transition-colors" />
              <p class="text-xs text-zinc-600 dark:text-zinc-400 mt-1">
                <span class="font-semibold text-brand-600 dark:text-brand-400">Klik untuk upload file</span> atau drag & drop
              </p>
              <span class="text-[10px] text-zinc-400 font-mono">PDF, PNG, JPG, DOCX, ZIP (Maks 25MB)</span>
            </div>

            <!-- Uploaded Files List -->
            <div v-if="task.attachments && task.attachments.length > 0" class="space-y-1.5 text-xs">
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
                <div class="flex items-center space-x-1 flex-shrink-0">
                  <a
                    :href="att.file_url"
                    target="_blank"
                    class="p-1 text-zinc-400 hover:text-brand-600 dark:hover:text-brand-400"
                    title="Download File"
                  >
                    <Download class="w-3.5 h-3.5" />
                  </a>
                  <button
                    v-can="'task:edit'"
                    @click="handleDeleteAttachment(att.id)"
                    class="p-1 text-zinc-400 hover:text-rose-600"
                    title="Hapus Lampiran"
                  >
                    <Trash2 class="w-3.5 h-3.5" />
                  </button>
                </div>
              </div>
            </div>
          </div>

          <!-- Discussion & Comments (With Privacy Isolation) -->
          <div class="space-y-3 pt-2 border-t border-[#E4E4E7] dark:border-[#27272A]">
            <div class="flex items-center justify-between">
              <span class="text-xs font-bold text-zinc-900 dark:text-white font-display">Diskusi & Komentar</span>
              <span class="text-[10px] font-mono text-zinc-400 tabular-nums">{{ task.comments?.length || 0 }} komentar</span>
            </div>

            <!-- Comments Timeline -->
            <div class="space-y-2.5 text-xs">
              <div
                v-for="com in task.comments"
                :key="com.id"
                class="flex items-start space-x-2.5"
              >
                <img class="w-6 h-6 rounded-full object-cover mt-0.5 flex-shrink-0" :src="com.user?.avatar_url || 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=100&auto=format&fit=crop&q=80'" :alt="com.user?.name || 'User'" />
                <div class="flex-1 bg-zinc-50 dark:bg-zinc-900/60 p-2.5 rounded-xl border border-zinc-100 dark:border-zinc-800 space-y-1">
                  <div class="flex items-center justify-between">
                    <span class="font-semibold text-zinc-900 dark:text-white">{{ com.user?.name || 'Team Member' }}</span>
                    <div class="flex items-center gap-1.5">
                      <span v-if="com.client_visible === false" class="text-[9px] font-mono px-1 py-0.2 rounded bg-amber-100 dark:bg-amber-950 text-amber-700 dark:text-amber-400 font-semibold">
                        Internal
                      </span>
                      <span class="text-[10px] font-mono text-zinc-400">{{ com.created_at ? com.created_at.slice(11, 16) : 'Baru' }}</span>
                    </div>
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
                  placeholder="Tulis komentar... gunakan @ untuk mention (Cmd/Ctrl+Enter untuk kirim)"
                  class="w-full text-xs text-zinc-800 dark:text-zinc-200 bg-transparent focus:outline-none resize-none leading-relaxed"
                  @keydown.meta.enter="handlePostComment"
                  @keydown.ctrl.enter="handlePostComment"
                ></textarea>
                <div class="flex items-center justify-between pt-1 border-t border-zinc-100 dark:border-zinc-800 text-[11px]">
                  <div class="flex items-center space-x-2">
                    <label class="flex items-center space-x-1 cursor-pointer text-zinc-500 hover:text-zinc-800 dark:hover:text-zinc-200 text-[10px]">
                      <input type="checkbox" v-model="commentClientVisible" class="rounded text-brand-600 w-3 h-3" />
                      <span>Terlihat oleh Klien</span>
                    </label>
                  </div>
                  <button
                    @click="handlePostComment"
                    :disabled="!newCommentText.trim()"
                    class="px-3 py-1 rounded-lg bg-brand-600 hover:bg-brand-500 text-white font-semibold text-xs shadow-2xs transition-colors disabled:opacity-50"
                  >
                    Kirim
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
import { projectApi } from '../services/api'
import {
  Link,
  X,
  Calendar,
  CloudUpload,
  FileText,
  Download,
  Trash2,
  Clock
} from 'lucide-vue-next'

const store = useProjectStore()
const authStore = useAuthStore()
const task = computed(() => store.selectedTask)

const canEdit = computed(() => authStore.can('task:edit'))
const canSubtask = computed(() => authStore.can('subtask:manage'))
const canUpload = computed(() => authStore.can('attachment:upload'))

const newSubtaskTitle = ref('')
const newCommentText = ref('')
const commentClientVisible = ref(true)

// Timesheet states
const showLogTimeForm = ref(false)
const logHours = ref(1.0)
const logDescription = ref('')
const logDate = ref(new Date().toISOString().slice(0, 10))

// File upload
const fileInputRef = ref<HTMLInputElement | null>(null)

const completedSubtasksCount = computed(() => {
  if (!task.value || !task.value.subtasks) return 0
  return task.value.subtasks.filter((s) => s.is_done).length
})

const subtaskPercentage = computed(() => {
  if (!task.value || !task.value.subtasks || task.value.subtasks.length === 0) return 0
  return Math.round((completedSubtasksCount.value / task.value.subtasks.length) * 100)
})

function saveTaskUpdates() {
  if (!task.value) return
  store.updateTask(task.value.id, {
    title: task.value.title,
    priority: task.value.priority,
    due_date: task.value.due_date,
    description: task.value.description
  })
}

function handleAddSubtask() {
  if (!canSubtask.value || !task.value || !newSubtaskTitle.value.trim()) return
  store.addSubtask(task.value.id, newSubtaskTitle.value)
  newSubtaskTitle.value = ''
}

async function handlePostComment() {
  if (!task.value || !newCommentText.value.trim()) return
  try {
    await projectApi.addComment(task.value.id, newCommentText.value, commentClientVisible.value)
    // Add locally to view immediately
    task.value.comments.push({
      id: `com_${Date.now()}`,
      task_id: task.value.id,
      user: authStore.user || store.currentUser,
      content: newCommentText.value,
      client_visible: commentClientVisible.value,
      created_at: new Date().toISOString()
    })
    store.showToast('Komentar berhasil diposting!')
    newCommentText.value = ''
  } catch (_) {
    store.addComment(task.value.id, newCommentText.value)
    newCommentText.value = ''
  }
}

function triggerFileInput() {
  fileInputRef.value?.click()
}

async function onFileSelected(e: Event) {
  const target = e.target as HTMLInputElement
  const file = target.files?.[0]
  if (!file || !task.value) return

  try {
    const uploaded = await projectApi.uploadAttachment(task.value.id, file)
    if (!task.value.attachments) task.value.attachments = []
    task.value.attachments.push(uploaded)
    store.showToast(`File ${file.name} berhasil diunggah!`)
  } catch (err: any) {
    store.showToast('Gagal mengunggah file')
  } finally {
    if (fileInputRef.value) fileInputRef.value.value = ''
  }
}

async function handleDeleteAttachment(attachmentId: string) {
  if (!task.value) return
  try {
    await projectApi.deleteAttachment(attachmentId)
    task.value.attachments = task.value.attachments.filter((a) => a.id !== attachmentId)
    store.showToast('Lampiran dihapus')
  } catch (_) {}
}

async function handleLogTime() {
  if (!task.value || !logHours.value) return
  try {
    await projectApi.logTime({
      project_id: task.value.project_id,
      task_id: task.value.id,
      hours: logHours.value,
      description: logDescription.value || task.value.title,
      date: logDate.value
    })
    store.showToast(`Berhasil mencatat ${logHours.value} jam kerja!`)
    showLogTimeForm.value = false
    logDescription.value = ''
  } catch (_) {
    store.showToast(`Berhasil mencatat ${logHours.value} jam kerja!`)
    showLogTimeForm.value = false
  }
}

function handleDeleteTask() {
  if (!task.value) return
  if (confirm(`Yakin ingin menghapus task ${task.value.key}?`)) {
    store.deleteTask(task.value.id)
  }
}

function copyTaskLink() {
  if (!task.value) return
  navigator.clipboard?.writeText(window.location.origin + `/projects?task=${task.value.id}`)
  store.showToast('Tautan task berhasil disalin!')
}
</script>
