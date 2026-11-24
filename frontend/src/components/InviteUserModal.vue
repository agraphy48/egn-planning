<template>
  <div
    v-if="isOpen"
    class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/50 backdrop-blur-xs select-none animate-in fade-in duration-150"
    @click.self="handleClose"
  >
    <div
      class="w-full max-w-lg bg-white rounded-3xl border border-slate-200/90 shadow-2xl p-6 sm:p-7 space-y-5 relative max-h-[90vh] flex flex-col overflow-hidden animate-in zoom-in-95 duration-150"
    >
      <!-- Top Decorative Accent -->
      <div class="absolute top-0 left-8 right-8 h-1 bg-gradient-to-r from-transparent via-brand-600 to-transparent rounded-full"></div>

      <!-- Header with Close Button -->
      <div class="flex items-start justify-between flex-shrink-0">
        <div class="flex items-center gap-3">
          <div class="p-2.5 rounded-2xl bg-indigo-50 border border-indigo-100 text-indigo-600 shadow-2xs">
            <Mail class="w-5 h-5 stroke-[2.2]" />
          </div>
          <div>
            <h2 class="text-base font-bold text-slate-900 font-display">Undang Klien / Pengguna</h2>
            <p class="text-xs text-slate-500 font-medium">Kirim tautan pendaftaran via email atau salin tautan undangan</p>
          </div>
        </div>
        <button
          @click="handleClose"
          class="p-1.5 rounded-xl text-slate-400 hover:text-slate-700 hover:bg-slate-100 transition-colors cursor-pointer"
        >
          <X class="w-4 h-4" />
        </button>
      </div>

      <!-- Result State: Invitation Created Successfully -->
      <div v-if="createdInvite" class="space-y-4 overflow-y-auto py-2 pr-1 flex-1">
        <div class="p-4 rounded-2xl bg-emerald-50 border border-emerald-200 text-emerald-900 text-xs space-y-3 animate-slide-up">
          <div class="flex items-start gap-2.5">
            <CheckCircle2 class="w-5 h-5 text-emerald-600 flex-shrink-0 mt-0.5" />
            <div>
              <p class="font-bold text-emerald-950">Undangan Berhasil Dibuat!</p>
              <p class="text-[11px] text-emerald-800 mt-0.5">
                Undangan untuk <strong>{{ createdInvite.email }}</strong> dengan peran <strong>{{ formatRole(createdInvite.role) }}</strong> siap digunakan.
              </p>
              <div v-if="emailDispatched" class="mt-2 inline-flex items-center gap-1.5 px-2.5 py-1 rounded-lg bg-emerald-200/60 text-emerald-950 text-[11px] font-semibold border border-emerald-300/60">
                <span>✉ Email aktivasi telah otomatis dikirimkan ke {{ createdInvite.email }}</span>
              </div>
            </div>
          </div>

          <!-- Shareable Link Box -->
          <div class="space-y-1.5 pt-1">
            <label class="text-[11px] font-semibold text-emerald-900 flex items-center justify-between">
              <span>Tautan Undangan Pendaftaran:</span>
              <span class="text-[10px] text-emerald-700 font-mono">Berlaku 7 Hari</span>
            </label>
            <div class="flex items-center gap-2">
              <input
                type="text"
                readonly
                :value="createdInvite.invite_url"
                class="flex-1 py-2 px-3 rounded-xl bg-white border border-emerald-200 text-slate-800 text-xs font-mono select-all focus:outline-none"
              />
              <button
                type="button"
                @click="copyInviteLink(createdInvite.invite_url)"
                class="py-2 px-3 rounded-xl bg-emerald-600 hover:bg-emerald-500 text-white font-bold text-xs shadow-md shadow-emerald-600/20 flex items-center gap-1.5 transition-all cursor-pointer flex-shrink-0"
              >
                <Check v-if="copied" class="w-3.5 h-3.5" />
                <Copy v-else class="w-3.5 h-3.5" />
                <span>{{ copied ? 'Tersalin!' : 'Salin Tautan' }}</span>
              </button>
            </div>
          </div>

          <p class="text-[10px] text-emerald-700 leading-relaxed">
            * Silakan kirimkan tautan di atas kepada klien melalui WhatsApp atau Email. Saat klien membuka tautan tersebut, mereka akan langsung mengonfirmasi nama dan kata sandi untuk login.
          </p>
        </div>

        <div class="flex items-center gap-2 pt-2">
          <button
            type="button"
            @click="resetForm"
            class="flex-1 py-2.5 px-3 rounded-xl bg-slate-100 hover:bg-slate-200/80 text-slate-700 font-semibold text-xs transition-colors cursor-pointer"
          >
            Undang Klien Lain
          </button>
          <button
            type="button"
            @click="handleClose"
            class="flex-1 py-2.5 px-3 rounded-xl bg-brand-600 hover:bg-brand-500 text-white font-bold text-xs shadow-md shadow-brand-600/20 transition-all cursor-pointer"
          >
            Selesai
          </button>
        </div>
      </div>

      <!-- Main Form: Create Invitation -->
      <form v-else @submit.prevent="handleSubmit" class="space-y-4 text-xs overflow-y-auto pr-1 flex-1">
        <!-- Error Alert Banner -->
        <div
          v-if="errorMessage"
          class="p-3 rounded-xl bg-rose-50 border border-rose-200 text-rose-800 text-xs flex items-center gap-2.5 animate-slide-up"
        >
          <AlertCircle class="w-4 h-4 text-rose-500 flex-shrink-0" />
          <span class="flex-1 font-medium">{{ errorMessage }}</span>
          <button type="button" @click="errorMessage = ''" class="text-rose-500 hover:text-rose-700 cursor-pointer">
            <X class="w-3.5 h-3.5" />
          </button>
        </div>

        <!-- Email Destination Input -->
        <div class="space-y-1.5">
          <label class="font-semibold text-slate-700 flex items-center justify-between">
            <span>Alamat Email Tujuan <span class="text-rose-500">*</span></span>
            <span class="text-[10px] text-slate-400 font-mono">Akan menjadi akun login</span>
          </label>
          <div class="relative">
            <Mail class="w-4 h-4 text-slate-400 absolute left-3 top-1/2 -translate-y-1/2" />
            <input
              type="email"
              v-model="email"
              required
              placeholder="klien@perusahaan.com"
              class="w-full pl-9 pr-3 py-2.5 rounded-xl border border-slate-300 bg-slate-50/50 text-slate-900 placeholder-slate-400 focus:bg-white focus:outline-none focus:ring-2 focus:ring-brand-500 focus:border-brand-500 transition-all text-xs shadow-xs"
            />
          </div>
        </div>

        <!-- Name Input (Optional) -->
        <div class="space-y-1.5">
          <label class="font-semibold text-slate-700 flex items-center justify-between">
            <span>Nama Calon Pengguna / Klien</span>
            <span class="text-[10px] text-slate-400 font-mono">Opsional (dapat diubah nanti)</span>
          </label>
          <div class="relative">
            <UserIcon class="w-4 h-4 text-slate-400 absolute left-3 top-1/2 -translate-y-1/2" />
            <input
              type="text"
              v-model="name"
              placeholder="Contoh: Bapak Ir. Bambang / PT Mitra Sejahtera"
              class="w-full pl-9 pr-3 py-2.5 rounded-xl border border-slate-300 bg-slate-50/50 text-slate-900 placeholder-slate-400 focus:bg-white focus:outline-none focus:ring-2 focus:ring-brand-500 focus:border-brand-500 transition-all text-xs shadow-xs"
            />
          </div>
        </div>

        <!-- Role / Hak Akses Selection -->
        <div class="space-y-1.5">
          <label class="font-semibold text-slate-700 flex items-center justify-between">
            <span>Hak Akses / Peran (Role) <span class="text-rose-500">*</span></span>
            <span class="text-[10px] text-slate-400 font-mono">Menentukan izin workspace</span>
          </label>
          <div class="grid grid-cols-3 gap-2">
            <!-- Client Option (Default) -->
            <label
              :class="[
                'p-3 rounded-2xl border flex flex-col justify-between cursor-pointer transition-all',
                role === 'client'
                  ? 'border-indigo-600 bg-indigo-50/60 ring-2 ring-indigo-500/20 shadow-xs'
                  : 'border-slate-200 hover:bg-slate-50'
              ]"
            >
              <input type="radio" value="client" v-model="role" class="sr-only" />
              <div class="flex items-center justify-between">
                <span class="text-xs font-bold text-slate-900">Client</span>
                <ShieldCheck class="w-4 h-4 text-indigo-600" />
              </div>
              <p class="text-[10px] text-slate-500 mt-1 leading-snug">Portal Klien, approval deliverables</p>
            </label>

            <!-- Member Option -->
            <label
              :class="[
                'p-3 rounded-2xl border flex flex-col justify-between cursor-pointer transition-all',
                role === 'member'
                  ? 'border-brand-600 bg-brand-50/60 ring-2 ring-brand-500/20 shadow-xs'
                  : 'border-slate-200 hover:bg-slate-50'
              ]"
            >
              <input type="radio" value="member" v-model="role" class="sr-only" />
              <div class="flex items-center justify-between">
                <span class="text-xs font-bold text-slate-900">Member</span>
                <Users class="w-4 h-4 text-brand-600" />
              </div>
              <p class="text-[10px] text-slate-500 mt-1 leading-snug">Pengerjaan tugas proyek tim</p>
            </label>

            <!-- PM Option -->
            <label
              :class="[
                'p-3 rounded-2xl border flex flex-col justify-between cursor-pointer transition-all',
                role === 'pm'
                  ? 'border-purple-600 bg-purple-50/60 ring-2 ring-purple-500/20 shadow-xs'
                  : 'border-slate-200 hover:bg-slate-50'
              ]"
            >
              <input type="radio" value="pm" v-model="role" class="sr-only" />
              <div class="flex items-center justify-between">
                <span class="text-xs font-bold text-slate-900">PM</span>
                <Briefcase class="w-4 h-4 text-purple-600" />
              </div>
              <p class="text-[10px] text-slate-500 mt-1 leading-snug">Kelola proyek & roadmap tim</p>
            </label>
          </div>
        </div>

        <!-- Project Assignment Checkboxes -->
        <div class="space-y-2 pt-1">
          <div class="flex items-center justify-between">
            <div>
              <label class="font-semibold text-slate-700">Tentukan Akses Proyek</label>
              <p class="text-[10px] text-slate-400">Pilih proyek yang langsung diberikan hak akses</p>
            </div>
            <div class="flex items-center gap-1.5 text-[10px]">
              <button
                type="button"
                @click="selectAllProjects"
                class="px-2 py-0.5 rounded-md text-brand-600 hover:bg-brand-50 font-medium cursor-pointer"
              >
                Pilih Semua
              </button>
              <span class="text-slate-300">|</span>
              <button
                type="button"
                @click="clearAllProjects"
                class="px-2 py-0.5 rounded-md text-slate-500 hover:bg-slate-100 font-medium cursor-pointer"
              >
                Hapus
              </button>
            </div>
          </div>

          <div class="p-2.5 rounded-2xl bg-slate-50 border border-slate-200/80 max-h-44 overflow-y-auto space-y-1">
            <div v-if="projectList.length === 0" class="py-3 text-center text-xs text-slate-400">
              Tidak ada proyek yang tersedia
            </div>
            <label
              v-for="p in projectList"
              :key="p.id"
              class="flex items-center gap-2.5 p-2 rounded-xl hover:bg-white transition-colors cursor-pointer text-xs"
              :class="selectedProjectIds.includes(p.id) ? 'bg-white shadow-2xs font-semibold' : ''"
            >
              <input
                type="checkbox"
                :value="p.id"
                v-model="selectedProjectIds"
                class="w-4 h-4 rounded text-brand-600 focus:ring-brand-500 border-slate-300"
              />
              <span
                class="w-2.5 h-2.5 rounded-full flex-shrink-0"
                :style="{ backgroundColor: p.color_dot || '#6366F1' }"
              ></span>
              <div class="flex-1 min-w-0">
                <span class="text-slate-900 truncate block">{{ p.name }}</span>
                <span v-if="p.client_name" class="text-[10px] text-slate-400 font-normal">Klien: {{ p.client_name }}</span>
              </div>
              <span class="text-[9px] font-mono px-1.5 py-0.2 rounded bg-slate-100 text-slate-500 uppercase">
                {{ p.status }}
              </span>
            </label>
          </div>
          <p class="text-[10px] text-slate-500 flex items-center gap-1 font-mono">
            <Info class="w-3 h-3 text-slate-400" />
            <span>Terpilih: {{ selectedProjectIds.length }} dari {{ projectList.length }} proyek</span>
          </p>
        </div>

        <!-- Action Buttons -->
        <div class="flex items-center gap-2 pt-3 border-t border-slate-100 flex-shrink-0">
          <button
            type="button"
            @click="handleClose"
            :disabled="isSubmitting"
            class="flex-1 py-2.5 px-3 rounded-xl bg-slate-100 hover:bg-slate-200/80 text-slate-700 font-semibold text-xs transition-colors cursor-pointer"
          >
            Batal
          </button>

          <button
            type="submit"
            :disabled="isSubmitting || !email.trim()"
            class="flex-1 py-2.5 px-3 rounded-xl bg-gradient-to-r from-brand-600 to-indigo-600 hover:from-brand-500 hover:to-indigo-500 text-white font-bold text-xs shadow-md shadow-brand-600/25 flex items-center justify-center gap-1.5 transition-all duration-200 active:scale-[0.98] disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer"
          >
            <Loader2 v-if="isSubmitting" class="w-4 h-4 animate-spin" />
            <span v-else>Buat & Kirim Undangan</span>
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { invitationApi } from '../services/api'
import type { Project, InvitationItem } from '../types/project'
import { useProjectStore } from '../stores/projectStore'
import {
  Mail,
  User as UserIcon,
  X,
  ShieldCheck,
  Users,
  Briefcase,
  Info,
  Check,
  CheckCircle2,
  Copy,
  AlertCircle,
  Loader2
} from 'lucide-vue-next'

const props = defineProps<{
  isOpen: boolean
  projects: Project[]
}>()

const emit = defineEmits(['close', 'invited'])

const projectStore = useProjectStore()

const email = ref('')
const name = ref('')
const role = ref<'client' | 'member' | 'pm'>('client')
const selectedProjectIds = ref<string[]>([])
const isSubmitting = ref(false)
const errorMessage = ref('')
const createdInvite = ref<InvitationItem | null>(null)
const copied = ref(false)

const projectList = ref<Project[]>([])

watch(
  () => props.projects,
  (val) => {
    projectList.value = val || []
  },
  { immediate: true }
)

watch(
  () => props.isOpen,
  (open) => {
    if (open) {
      resetForm()
      projectList.value = props.projects || []
    }
  }
)

function selectAllProjects() {
  selectedProjectIds.value = projectList.value.map((p) => p.id)
}

function clearAllProjects() {
  selectedProjectIds.value = []
}

const emailDispatched = ref(false)

function resetForm() {
  email.value = ''
  name.value = ''
  role.value = 'client'
  selectedProjectIds.value = []
  errorMessage.value = ''
  createdInvite.value = null
  copied.value = false
  emailDispatched.value = false
  isSubmitting.value = false
}

function handleClose() {
  resetForm()
  emit('close')
}

function formatRole(r: string) {
  if (r === 'client') return 'Client Portal'
  if (r === 'pm') return 'Project Manager'
  if (r === 'member') return 'Team Member'
  return r
}

async function copyInviteLink(url?: string) {
  if (!url) return
  try {
    await navigator.clipboard.writeText(url)
    copied.value = true
    projectStore.showToast('Tautan undangan berhasil disalin ke clipboard')
    setTimeout(() => {
      copied.value = false
    }, 3000)
  } catch (err) {
    console.error('Failed to copy link:', err)
  }
}

async function handleSubmit() {
  if (!email.value.trim()) return

  errorMessage.value = ''
  isSubmitting.value = true

  try {
    const res = await invitationApi.createInvitation({
      email: email.value.trim(),
      name: name.value.trim(),
      role: role.value,
      project_ids: selectedProjectIds.value
    })

    createdInvite.value = res.invitation
    emailDispatched.value = (res as any).email_dispatched ?? false
    projectStore.showToast('Undangan berhasil dibuat!')
    emit('invited')
  } catch (err: any) {
    errorMessage.value =
      err?.response?.data?.error || err.message || 'Gagal membuat undangan. Silakan periksa kembali data Anda.'
  } finally {
    isSubmitting.value = false
  }
}
</script>
