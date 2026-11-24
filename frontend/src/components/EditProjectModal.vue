<template>
  <div
    v-if="isOpen && project"
    class="fixed inset-0 bg-black/50 backdrop-blur-xs flex items-center justify-center p-4 z-50 select-none"
    @click.self="close"
  >
    <div class="max-w-md w-full bg-white dark:bg-[#18181B] rounded-2xl p-5 border border-zinc-200 dark:border-zinc-800 shadow-2xl space-y-4 animate-in fade-in zoom-in-95 duration-150">
      <!-- Header -->
      <div class="flex items-center justify-between pb-2 border-b border-zinc-100 dark:border-zinc-800">
        <div class="flex items-center space-x-2">
          <FolderEdit class="w-5 h-5 text-brand-600 dark:text-brand-400" />
          <h3 class="text-sm font-bold text-zinc-900 dark:text-white">Edit Proyek</h3>
        </div>
        <button @click="close" class="text-zinc-400 hover:text-zinc-600 dark:hover:text-zinc-200 cursor-pointer">
          <X class="w-4 h-4" />
        </button>
      </div>

      <!-- Edit Form -->
      <form @submit.prevent="handleSubmit" class="space-y-3.5 text-xs">
        <!-- Project Name -->
        <div>
          <label class="font-semibold text-zinc-700 dark:text-zinc-300">Nama Proyek <span class="text-rose-500">*</span></label>
          <input
            v-model="name"
            required
            placeholder="Masukkan nama proyek..."
            class="w-full mt-1 px-3 py-2 rounded-lg border border-zinc-300 dark:border-zinc-700 bg-white dark:bg-zinc-900 text-zinc-900 dark:text-white focus:outline-none focus:ring-1 focus:ring-brand-500 text-xs shadow-2xs"
          />
        </div>

        <!-- Client Name -->
        <div>
          <label class="font-semibold text-zinc-700 dark:text-zinc-300">Nama Klien / Stakeholder</label>
          <input
            v-model="clientName"
            placeholder="Contoh: PT Kedukan Aksara Global"
            class="w-full mt-1 px-3 py-2 rounded-lg border border-zinc-300 dark:border-zinc-700 bg-white dark:bg-zinc-900 text-zinc-900 dark:text-white focus:outline-none focus:ring-1 focus:ring-brand-500 text-xs shadow-2xs"
          />
        </div>

        <!-- Status & Deadline Grid -->
        <div class="grid grid-cols-2 gap-2.5">
          <div>
            <label class="font-semibold text-zinc-700 dark:text-zinc-300">Status Proyek</label>
            <select
              v-model="status"
              class="w-full mt-1 px-3 py-2 rounded-lg border border-zinc-300 dark:border-zinc-700 bg-white dark:bg-zinc-900 text-zinc-900 dark:text-white focus:outline-none focus:ring-1 focus:ring-brand-500 text-xs shadow-2xs cursor-pointer font-medium"
            >
              <option value="active">Aktif (In Progress)</option>
              <option value="on_hold">Ditunda (On Hold)</option>
              <option value="completed">Selesai (Completed)</option>
              <option value="archived">Diarsipkan (Archived)</option>
            </select>
          </div>

          <div>
            <label class="font-semibold text-zinc-700 dark:text-zinc-300">Tenggat Waktu</label>
            <input
              type="date"
              v-model="deadline"
              required
              class="w-full mt-1 px-3 py-2 rounded-lg border border-zinc-300 dark:border-zinc-700 bg-white dark:bg-zinc-900 text-zinc-900 dark:text-white focus:outline-none focus:ring-1 focus:ring-brand-500 font-mono text-xs shadow-2xs"
            />
          </div>
        </div>

        <!-- Color Tag Selection -->
        <div>
          <label class="font-semibold text-zinc-700 dark:text-zinc-300 block mb-1">Warna Tag Proyek</label>
          <div class="flex items-center space-x-2">
            <button
              v-for="c in colorOptions"
              :key="c"
              type="button"
              @click="colorDot = c"
              :class="['w-6 h-6 rounded-full transition-transform cursor-pointer', colorDot === c ? 'scale-110 ring-2 ring-offset-2 ring-brand-500' : 'hover:scale-105']"
              :style="{ backgroundColor: c }"
            ></button>
          </div>
        </div>

        <!-- Action Buttons -->
        <div class="flex items-center justify-end space-x-2 pt-3 border-t border-zinc-100 dark:border-zinc-800">
          <button
            type="button"
            @click="close"
            class="px-3.5 py-1.5 rounded-lg border border-zinc-200 dark:border-zinc-700 text-xs font-semibold text-zinc-600 dark:text-zinc-400 hover:bg-zinc-100 dark:hover:bg-zinc-800 cursor-pointer"
          >
            Batal
          </button>
          <button
            type="submit"
            :disabled="isSubmitting || !name.trim()"
            class="px-4 py-1.5 rounded-lg bg-brand-600 hover:bg-brand-500 text-white text-xs font-bold disabled:opacity-50 flex items-center gap-1.5 cursor-pointer shadow-md shadow-brand-600/20"
          >
            <Loader2 v-if="isSubmitting" class="w-3.5 h-3.5 animate-spin" />
            <span>Simpan Perubahan</span>
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useProjectStore } from '../stores/projectStore'
import { projectApi } from '../services/api'
import type { Project, ProjectStatus } from '../types/project'
import { FolderEdit, X, Loader2 } from 'lucide-vue-next'

const props = defineProps<{
  isOpen: boolean
  project: Project | null
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'updated', project: Project): void
}>()

const store = useProjectStore()

const name = ref('')
const clientName = ref('')
const status = ref<ProjectStatus>('active')
const deadline = ref('')
const colorDot = ref('#4F46E5')
const isSubmitting = ref(false)

const colorOptions = ['#4F46E5', '#10B981', '#F59E0B', '#EF4444', '#8B5CF6']

watch(
  () => props.project,
  (newProj) => {
    if (newProj) {
      name.value = newProj.name || ''
      clientName.value = newProj.client_name || ''
      status.value = newProj.status || 'active'
      deadline.value = newProj.deadline || ''
      colorDot.value = newProj.color_dot || '#4F46E5'
    }
  },
  { immediate: true }
)

function close() {
  emit('close')
}

async function handleSubmit() {
  if (!props.project || !name.value.trim()) return
  isSubmitting.value = true

  try {
    const updated = await projectApi.updateProject(props.project.id, {
      name: name.value.trim(),
      client_name: clientName.value.trim() || undefined,
      status: status.value,
      deadline: deadline.value,
      color_dot: colorDot.value
    })

    // Update in local store
    const idx = store.projects.findIndex((p) => p.id === props.project!.id)
    if (idx !== -1) {
      store.projects[idx] = {
        ...store.projects[idx],
        ...updated
      }
    }

    store.showToast(`Proyek ${updated.name} berhasil diperbarui!`)
    emit('updated', updated)
    emit('close')
  } catch (err: any) {
    store.showToast(err?.response?.data?.error || 'Gagal memperbarui proyek')
  } finally {
    isSubmitting.value = false
  }
}
</script>
