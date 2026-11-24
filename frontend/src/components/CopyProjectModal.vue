<template>
  <div
    v-if="isOpen && project"
    class="fixed inset-0 bg-black/50 backdrop-blur-xs flex items-center justify-center p-4 z-50 select-none"
    @click.self="close"
  >
    <div class="max-w-md w-full bg-white dark:bg-[#18181B] rounded-3xl p-5 sm:p-6 border border-zinc-200 dark:border-zinc-800 shadow-2xl space-y-4 animate-in fade-in zoom-in-95 duration-150">
      
      <!-- Header -->
      <div class="flex items-center justify-between pb-3 border-b border-zinc-100 dark:border-zinc-800">
        <div class="flex items-center space-x-2.5">
          <div class="p-2 rounded-xl bg-brand-50 dark:bg-brand-950/50 text-brand-600 dark:text-brand-400 border border-brand-200 dark:border-brand-800">
            <Copy class="w-5 h-5" />
          </div>
          <div>
            <h3 class="text-sm font-bold text-zinc-900 dark:text-white">Duplikasi Proyek</h3>
            <p class="text-[11px] text-zinc-500">Salin seluruh deliverable & tugas ke proyek baru</p>
          </div>
        </div>
        <button @click="close" class="text-zinc-400 hover:text-zinc-600 dark:hover:text-zinc-200 p-1 cursor-pointer">
          <X class="w-4 h-4" />
        </button>
      </div>

      <!-- Source Project Reference Pill -->
      <div class="p-3 rounded-xl bg-slate-50 dark:bg-zinc-900/60 border border-slate-200/80 dark:border-zinc-800 text-xs flex items-center justify-between">
        <div class="flex items-center gap-2 min-w-0">
          <span class="w-2.5 h-2.5 rounded-full flex-shrink-0" :style="{ backgroundColor: project.color_dot }"></span>
          <div class="min-w-0">
            <span class="text-[10px] text-slate-400 block font-mono uppercase">Proyek Sumber</span>
            <span class="font-bold text-slate-800 dark:text-zinc-200 truncate block">{{ project.name }}</span>
          </div>
        </div>
        <span class="text-[10px] font-mono px-2 py-0.5 rounded-full bg-brand-50 text-brand-700 dark:bg-brand-950/60 dark:text-brand-300 font-semibold border border-brand-200 dark:border-brand-800 flex-shrink-0">
          {{ project.total_tasks || 0 }} deliverable
        </span>
      </div>

      <!-- Copy Form -->
      <form @submit.prevent="handleSubmit" class="space-y-3.5 text-xs">
        <!-- New Project Name -->
        <div>
          <label class="font-semibold text-zinc-700 dark:text-zinc-300 flex items-center justify-between">
            <span>Nama Proyek Baru <span class="text-rose-500">*</span></span>
            <span class="text-[10px] text-zinc-400 font-mono">Harus Unik</span>
          </label>
          <input
            v-model="name"
            required
            placeholder="Masukkan nama proyek baru..."
            class="w-full mt-1 px-3 py-2 rounded-xl border border-zinc-300 dark:border-zinc-700 bg-white dark:bg-zinc-900 text-zinc-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-brand-500 text-xs shadow-2xs"
          />
        </div>

        <!-- Client Name -->
        <div>
          <label class="font-semibold text-zinc-700 dark:text-zinc-300">Nama Klien / Stakeholder</label>
          <input
            v-model="clientName"
            placeholder="Contoh: PT Kedukan Aksara Global"
            class="w-full mt-1 px-3 py-2 rounded-xl border border-zinc-300 dark:border-zinc-700 bg-white dark:bg-zinc-900 text-zinc-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-brand-500 text-xs shadow-2xs"
          />
        </div>

        <!-- Deadline & Color Tag -->
        <div class="grid grid-cols-2 gap-2.5">
          <div>
            <label class="font-semibold text-zinc-700 dark:text-zinc-300">Tenggat Waktu Baru</label>
            <input
              type="date"
              v-model="deadline"
              required
              class="w-full mt-1 px-3 py-2 rounded-xl border border-zinc-300 dark:border-zinc-700 bg-white dark:bg-zinc-900 text-zinc-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-brand-500 font-mono text-xs shadow-2xs"
            />
          </div>
          <div>
            <label class="font-semibold text-zinc-700 dark:text-zinc-300 block mb-1">Warna Tag Proyek</label>
            <div class="flex items-center space-x-1.5 mt-1.5">
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
        </div>

        <!-- Duplication Options (Smart Checkboxes) -->
        <div class="p-3 rounded-2xl bg-indigo-50/40 dark:bg-zinc-900/60 border border-indigo-100 dark:border-zinc-800 space-y-2 text-[11px]">
          <span class="text-[10px] font-mono uppercase text-indigo-700 dark:text-indigo-400 font-bold block">Opsi Duplikasi Pekerjaan</span>
          
          <label class="flex items-start gap-2 cursor-pointer text-slate-700 dark:text-zinc-300 select-none">
            <input
              type="checkbox"
              v-model="resetTaskStatus"
              class="mt-0.5 rounded border-slate-300 text-brand-600 focus:ring-brand-500 cursor-pointer"
            />
            <div class="space-y-0.5">
              <span class="font-semibold block">Reset Status Tugas ke 'Backlog'</span>
              <span class="text-[10px] text-slate-500 block leading-tight">Mulai seluruh tugas dari awal dengan checklist subtask di-reset belum selesai.</span>
            </div>
          </label>

          <label class="flex items-start gap-2 cursor-pointer text-slate-700 dark:text-zinc-300 select-none">
            <input
              type="checkbox"
              v-model="copyMembers"
              class="mt-0.5 rounded border-slate-300 text-brand-600 focus:ring-brand-500 cursor-pointer"
            />
            <div class="space-y-0.5">
              <span class="font-semibold block">Salin Hak Akses Anggota Tim</span>
              <span class="text-[10px] text-slate-500 block leading-tight">Daftarkan seluruh anggota tim proyek sumber ke proyek baru dengan hak akses yang sama.</span>
            </div>
          </label>

          <label class="flex items-start gap-2 cursor-pointer text-slate-700 dark:text-zinc-300 select-none">
            <input
              type="checkbox"
              v-model="copyWiki"
              class="mt-0.5 rounded border-slate-300 text-brand-600 focus:ring-brand-500 cursor-pointer"
            />
            <div class="space-y-0.5">
              <span class="font-semibold block">Salin Dokumentasi Wiki Proyek</span>
              <span class="text-[10px] text-slate-500 block leading-tight">Duplikasi seluruh halaman panduan & knowledge base proyek.</span>
            </div>
          </label>
        </div>

        <!-- Action Buttons -->
        <div class="flex items-center justify-end space-x-2 pt-3 border-t border-zinc-100 dark:border-zinc-800">
          <button
            type="button"
            @click="close"
            class="px-4 py-2 rounded-xl border border-zinc-200 dark:border-zinc-700 text-xs font-semibold text-zinc-600 dark:text-zinc-400 hover:bg-zinc-100 dark:hover:bg-zinc-800 cursor-pointer"
          >
            Batal
          </button>
          <button
            type="submit"
            :disabled="isSubmitting || !name.trim()"
            class="px-4 py-2 rounded-xl bg-brand-600 hover:bg-brand-500 text-white text-xs font-bold disabled:opacity-50 flex items-center gap-2 cursor-pointer shadow-md shadow-brand-600/25 transition-all"
          >
            <Loader2 v-if="isSubmitting" class="w-4 h-4 animate-spin" />
            <Copy v-else class="w-4 h-4" />
            <span>Duplikasi Proyek Sekarang</span>
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
import type { Project } from '../types/project'
import { Copy, X, Loader2 } from 'lucide-vue-next'

const props = defineProps<{
  isOpen: boolean
  project: Project | null
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'copied', newProject: Project): void
}>()

const store = useProjectStore()

const name = ref('')
const clientName = ref('')
const deadline = ref('')
const colorDot = ref('#4F46E5')
const resetTaskStatus = ref(true)
const copyMembers = ref(true)
const copyWiki = ref(true)
const isSubmitting = ref(false)

const colorOptions = ['#4F46E5', '#10B981', '#F59E0B', '#EF4444', '#8B5CF6']

watch(
  () => props.project,
  (newProj) => {
    if (newProj) {
      name.value = `${newProj.name} (Salinan)`
      clientName.value = newProj.client_name || ''
      deadline.value = newProj.deadline || '2026-12-31'
      colorDot.value = newProj.color_dot || '#4F46E5'
      resetTaskStatus.value = true
      copyMembers.value = true
      copyWiki.value = true
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
    const res = await projectApi.copyProject(props.project.id, {
      name: name.value.trim(),
      client_name: clientName.value.trim() || undefined,
      deadline: deadline.value,
      color_dot: colorDot.value,
      copy_members: copyMembers.value,
      reset_task_status: resetTaskStatus.value,
      copy_wiki: copyWiki.value
    })

    store.projects.unshift(res.project)
    store.showToast(`Proyek "${res.project.name}" berhasil diduplikasi (${res.tasks_copied} pekerjaan disalin)!`)
    emit('copied', res.project)
    emit('close')
  } catch (err: any) {
    store.showToast(err?.response?.data?.error || 'Gagal menduplikasi proyek')
  } finally {
    isSubmitting.value = false
  }
}
</script>
