<template>
  <div
    v-if="isOpen"
    class="fixed inset-0 bg-black/50 backdrop-blur-xs flex items-center justify-center p-4 z-50 select-none"
  >
    <div class="max-w-md w-full bg-white dark:bg-[#18181B] rounded-2xl p-5 border border-zinc-200 dark:border-zinc-800 shadow-2xl space-y-4">
      <div class="flex items-center justify-between pb-2 border-b border-zinc-100 dark:border-zinc-800">
        <div class="flex items-center space-x-2">
          <FolderPlus class="w-5 h-5 text-brand-600 dark:text-brand-400" />
          <h3 class="text-sm font-bold text-zinc-900 dark:text-white">Buat Proyek Baru</h3>
        </div>
        <button @click="close" class="text-zinc-400 hover:text-zinc-600 dark:hover:text-zinc-200">
          <X class="w-4 h-4" />
        </button>
      </div>

      <form @submit.prevent="handleSubmit" class="space-y-3 text-xs">
        <div>
          <label class="font-semibold text-zinc-700 dark:text-zinc-300">Nama Proyek</label>
          <input
            v-model="name"
            required
            placeholder="Contoh: E-Commerce AI Recommendation"
            class="w-full mt-1 px-3 py-2 rounded-lg border border-zinc-300 dark:border-zinc-700 bg-white dark:bg-zinc-900 text-zinc-900 dark:text-white focus:outline-none focus:ring-1 focus:ring-brand-500"
          />
        </div>

        <div>
          <label class="font-semibold text-zinc-700 dark:text-zinc-300">Nama Klien / Stakeholder</label>
          <input
            v-model="clientName"
            placeholder="Contoh: PT Kedukan Aksara Global"
            class="w-full mt-1 px-3 py-2 rounded-lg border border-zinc-300 dark:border-zinc-700 bg-white dark:bg-zinc-900 text-zinc-900 dark:text-white focus:outline-none focus:ring-1 focus:ring-brand-500"
          />
        </div>

        <div class="grid grid-cols-2 gap-2">
          <div>
            <label class="font-semibold text-zinc-700 dark:text-zinc-300">Tenggat Waktu</label>
            <input
              type="date"
              v-model="deadline"
              required
              class="w-full mt-1 px-3 py-2 rounded-lg border border-zinc-300 dark:border-zinc-700 bg-white dark:bg-zinc-900 text-zinc-900 dark:text-white focus:outline-none focus:ring-1 focus:ring-brand-500 font-mono"
            />
          </div>
          <div>
            <label class="font-semibold text-zinc-700 dark:text-zinc-300">Warna Tag Proyek</label>
            <div class="flex items-center space-x-2 mt-1.5">
              <button
                v-for="c in colorOptions"
                :key="c"
                type="button"
                @click="colorDot = c"
                :class="['w-6 h-6 rounded-full transition-transform', colorDot === c ? 'scale-110 ring-2 ring-offset-2 ring-brand-500' : '']"
                :style="{ backgroundColor: c }"
              ></button>
            </div>
          </div>
        </div>

        <div class="flex items-center justify-end space-x-2 pt-3 border-t border-zinc-100 dark:border-zinc-800">
          <button
            type="button"
            @click="close"
            class="px-3 py-1.5 rounded-lg border border-zinc-200 dark:border-zinc-700 text-xs font-semibold text-zinc-600 dark:text-zinc-400 hover:bg-zinc-100 dark:hover:bg-zinc-800"
          >
            Batal
          </button>
          <button
            type="submit"
            :disabled="isSubmitting || !name.trim()"
            class="px-4 py-1.5 rounded-lg bg-brand-600 hover:bg-brand-500 text-white text-xs font-bold disabled:opacity-50 flex items-center gap-1.5"
          >
            <Loader2 v-if="isSubmitting" class="w-3.5 h-3.5 animate-spin" />
            <span>Buat Proyek</span>
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useProjectStore } from '../stores/projectStore'
import { projectApi } from '../services/api'
import { FolderPlus, X, Loader2 } from 'lucide-vue-next'

const props = defineProps<{
  isOpen: boolean
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'created', project: any): void
}>()

const store = useProjectStore()

const name = ref('')
const clientName = ref('')
const deadline = ref('2026-12-31')
const colorDot = ref('#4F46E5')
const isSubmitting = ref(false)

const colorOptions = ['#4F46E5', '#10B981', '#F59E0B', '#EF4444', '#8B5CF6']

function close() {
  emit('close')
}

async function handleSubmit() {
  if (!name.value.trim()) return
  isSubmitting.value = true
  try {
    const slug = name.value.toLowerCase().replace(/[^a-z0-9]+/g, '-')
    const created = await projectApi.createProject({
      name: name.value.trim(),
      client_name: clientName.value.trim() || undefined,
      deadline: deadline.value,
      color_dot: colorDot.value,
      slug: slug,
      status: 'active'
    })
    store.projects.push(created)
    store.selectProject(created.id)
    store.showToast(`Proyek ${created.name} berhasil dibuat!`)
    emit('created', created)
    emit('close')
    name.value = ''
    clientName.value = ''
  } catch (err: any) {
    store.showToast('Gagal membuat proyek')
  } finally {
    isSubmitting.value = false
  }
}
</script>
