<template>
  <div
    v-if="isOpen && project"
    class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/50 backdrop-blur-xs select-none"
    @click.self="handleClose"
  >
    <div class="max-w-2xl w-full bg-white rounded-3xl border border-slate-200 shadow-2xl p-6 sm:p-7 space-y-5 relative max-h-[90vh] flex flex-col animate-in fade-in zoom-in-95 duration-150">
      <!-- Header -->
      <div class="flex items-center justify-between pb-3 border-b border-slate-100 flex-shrink-0">
        <div class="flex items-center gap-3">
          <div class="p-2.5 rounded-2xl bg-brand-50 border border-brand-200 text-brand-600 flex-shrink-0">
            <ShieldCheck class="w-5 h-5" />
          </div>
          <div>
            <h3 class="text-base sm:text-lg font-bold text-slate-900 flex items-center gap-2">
              <span>Hak Akses Proyek</span>
              <span class="text-xs px-2 py-0.5 rounded-full font-mono font-semibold bg-slate-100 border border-slate-200 text-slate-700">
                {{ project.name }}
              </span>
            </h3>
            <p class="text-xs text-slate-500">Tentukan dan kelola siapa saja anggota tim yang dapat mengakses proyek ini</p>
          </div>
        </div>
        <button
          @click="handleClose"
          class="p-2 rounded-xl text-slate-400 hover:text-slate-600 hover:bg-slate-100 transition-colors cursor-pointer"
        >
          <X class="w-5 h-5" />
        </button>
      </div>

      <!-- Scrollable Content -->
      <div class="space-y-4 overflow-y-auto pr-1 flex-1">
        <!-- Project Metadata Pill -->
        <div class="p-3 rounded-2xl bg-slate-50 border border-slate-200/80 flex flex-wrap items-center justify-between gap-2 text-xs">
          <div class="flex items-center gap-2 font-mono text-slate-600">
            <span class="w-2.5 h-2.5 rounded-full" :style="{ backgroundColor: project.color_dot }"></span>
            <span class="font-bold text-slate-900">{{ project.slug }}</span>
            <span class="text-slate-400">|</span>
            <span>Klien: {{ project.client_name || 'Internal' }}</span>
          </div>
          <div class="flex items-center gap-1.5 font-mono text-emerald-700 font-semibold bg-emerald-50 px-2 py-0.5 rounded-full border border-emerald-200 text-[11px]">
            <span class="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse"></span>
            <span>{{ members.length }} Anggota Terdaftar</span>
          </div>
        </div>

        <!-- Section 1: Grant Access Form -->
        <div class="p-4 rounded-2xl bg-gradient-to-br from-slate-50 to-indigo-50/30 border border-slate-200 space-y-3">
          <div class="flex items-center justify-between">
            <h4 class="text-xs font-bold text-slate-800 uppercase tracking-wider font-mono flex items-center gap-1.5">
              <UserPlus class="w-3.5 h-3.5 text-brand-600" />
              <span>Berikan Hak Akses ke Anggota Tim</span>
            </h4>
            <span class="text-[10px] text-slate-500 font-mono">
              {{ availableUsers.length }} pengguna tersedia
            </span>
          </div>

          <form @submit.prevent="handleAddMember" class="grid grid-cols-1 sm:grid-cols-12 gap-2.5">
            <!-- User Selection Dropdown -->
            <div class="sm:col-span-7">
              <label class="block text-[11px] font-semibold text-slate-600 mb-1">Pilih Pengguna</label>
              <select
                v-model="newMemberUserId"
                required
                :disabled="availableUsers.length === 0"
                class="w-full px-3 py-2 rounded-xl border border-slate-300 bg-white text-slate-900 text-xs focus:outline-none focus:ring-2 focus:ring-brand-500 disabled:opacity-50 disabled:bg-slate-100"
              >
                <option value="" disabled>
                  {{ availableUsers.length === 0 ? '-- Semua pengguna sudah memiliki akses --' : '-- Pilih Anggota Tim --' }}
                </option>
                <option v-for="u in availableUsers" :key="u.id" :value="u.id">
                  {{ u.name }} ({{ u.email }}) — {{ formatRole(u.role) }}
                </option>
              </select>
            </div>

            <!-- Role in Project Dropdown -->
            <div class="sm:col-span-3">
              <label class="block text-[11px] font-semibold text-slate-600 mb-1">Peran di Proyek</label>
              <select
                v-model="newMemberRole"
                class="w-full px-3 py-2 rounded-xl border border-slate-300 bg-white text-slate-900 text-xs focus:outline-none focus:ring-2 focus:ring-brand-500"
              >
                <option value="member">Team Member</option>
                <option value="pm">Project Lead (PM)</option>
                <option value="client">Client / Observer</option>
              </select>
            </div>

            <!-- Submit Button -->
            <div class="sm:col-span-2 flex items-end">
              <button
                type="submit"
                :disabled="!newMemberUserId || isSubmitting"
                class="w-full py-2 px-3 rounded-xl bg-brand-600 hover:bg-brand-500 disabled:opacity-50 disabled:cursor-not-allowed text-white text-xs font-bold shadow-md shadow-brand-600/20 flex items-center justify-center gap-1.5 transition-all cursor-pointer"
              >
                <Loader2 v-if="isSubmitting" class="w-3.5 h-3.5 animate-spin" />
                <span v-else>+ Tambah</span>
              </button>
            </div>
          </form>

          <p v-if="grantErrorMessage" class="text-[11px] text-rose-600 font-medium">
            {{ grantErrorMessage }}
          </p>
        </div>

        <!-- Section 2: Current Members List -->
        <div class="space-y-2">
          <div class="flex items-center justify-between text-xs">
            <span class="font-bold text-slate-800 flex items-center gap-1.5">
              <span>Daftar Anggota dengan Akses Aktif</span>
            </span>
            <span class="text-[11px] text-slate-400 font-mono">
              Total: {{ members.length }} orang
            </span>
          </div>

          <!-- Loading Skeleton -->
          <div v-if="isLoadingMembers" class="py-8 text-center text-xs text-slate-400 flex items-center justify-center gap-2">
            <Loader2 class="w-4 h-4 animate-spin text-brand-600" />
            <span>Memuat data keanggotaan proyek...</span>
          </div>

          <!-- Empty State -->
          <div v-else-if="members.length === 0" class="p-6 rounded-2xl border border-dashed border-slate-200 text-center text-xs text-slate-400">
            Belum ada anggota yang terdaftar di proyek ini.
          </div>

          <!-- Members Table / Card Rows -->
          <div v-else class="border border-slate-200 rounded-2xl overflow-hidden divide-y divide-slate-100 bg-white">
            <div
              v-for="m in members"
              :key="m.user_id"
              class="p-3 sm:p-3.5 flex items-center justify-between gap-3 hover:bg-slate-50/60 transition-colors"
            >
              <!-- User Info -->
              <div class="flex items-center gap-3 min-w-0">
                <img
                  :src="m.user?.avatar_url || 'https://images.unsplash.com/photo-1534528741775-53994a69daeb?w=100&auto=format&fit=crop&q=80'"
                  :alt="m.user?.name || 'User'"
                  class="w-9 h-9 rounded-xl object-cover ring-1 ring-slate-200 flex-shrink-0"
                />
                <div class="min-w-0">
                  <div class="flex items-center gap-1.5 flex-wrap">
                    <span class="text-xs font-bold text-slate-900 truncate">{{ m.user?.name || m.user_id }}</span>
                    <!-- Project Role Badge -->
                    <span
                      :class="[
                        'px-2 py-0.2 rounded-full text-[10px] font-mono font-semibold border',
                        m.project_role === 'pm'
                          ? 'bg-emerald-50 text-emerald-700 border-emerald-200'
                          : m.project_role === 'client'
                          ? 'bg-amber-50 text-amber-700 border-amber-200'
                          : 'bg-indigo-50 text-indigo-700 border-indigo-200'
                      ]"
                    >
                      {{ m.project_role === 'pm' ? 'Lead (PM)' : m.project_role === 'client' ? 'Client' : 'Member' }}
                    </span>
                  </div>
                  <div class="text-[11px] text-slate-500 font-mono truncate">{{ m.user?.email || '-' }}</div>
                </div>
              </div>

              <!-- Action: Revoke Access -->
              <div class="flex items-center gap-2 flex-shrink-0">
                <button
                  type="button"
                  @click="handleRemoveMember(m.user_id, m.user?.name || m.user_id)"
                  :disabled="isRevokingId === m.user_id"
                  class="px-2.5 py-1.5 rounded-xl bg-rose-50 hover:bg-rose-100 border border-rose-200 text-rose-700 text-xs font-semibold flex items-center gap-1 transition-colors cursor-pointer disabled:opacity-50"
                  title="Cabut hak akses anggota ini dari proyek"
                >
                  <Loader2 v-if="isRevokingId === m.user_id" class="w-3 h-3 animate-spin" />
                  <UserMinus v-else class="w-3 h-3" />
                  <span class="hidden sm:inline">Cabut Akses</span>
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- Privacy & Security Notice -->
        <div class="p-3.5 rounded-2xl bg-indigo-50/50 border border-indigo-100 text-[11px] text-indigo-900 flex items-start gap-2.5">
          <ShieldAlert class="w-4 h-4 text-indigo-600 flex-shrink-0 mt-0.5" />
          <p class="leading-relaxed">
            <strong>Proteksi Zero-Trust:</strong> Pengguna yang tidak tercantum dalam daftar akses di atas otomatis <em>tidak dapat melihat</em> proyek ini di daftar portfolio mereka dan akan ditolak (403 Forbidden) jika mencoba mengakses tugas proyek secara langsung.
          </p>
        </div>
      </div>

      <!-- Footer Actions -->
      <div class="pt-3 border-t border-slate-100 flex items-center justify-end flex-shrink-0">
        <button
          type="button"
          @click="handleClose"
          class="px-5 py-2 rounded-xl bg-slate-100 hover:bg-slate-200 text-slate-700 text-xs font-bold transition-colors cursor-pointer"
        >
          Selesai & Tutup
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { projectApi, userApi } from '../services/api'
import type { Project, ProjectMember, User } from '../types/project'
import { useProjectStore } from '../stores/projectStore'
import {
  ShieldCheck,
  ShieldAlert,
  UserPlus,
  UserMinus,
  X,
  Loader2
} from 'lucide-vue-next'

const props = defineProps<{
  isOpen: boolean
  project: Project | null
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'updated'): void
}>()

const store = useProjectStore()

const members = ref<ProjectMember[]>([])
const allUsers = ref<User[]>([])
const isLoadingMembers = ref(false)

const newMemberUserId = ref('')
const newMemberRole = ref('member')
const isSubmitting = ref(false)
const isRevokingId = ref<string | null>(null)
const grantErrorMessage = ref('')

// Compute users who are NOT yet members of this project
const availableUsers = computed(() => {
  const memberUserIds = new Set(members.value.map(m => m.user_id))
  return allUsers.value.filter(u => !memberUserIds.has(u.id))
})

function formatRole(role: string) {
  switch (role) {
    case 'super_admin': return 'Super Admin'
    case 'pm': return 'Project Manager'
    case 'member': return 'Team Member'
    case 'client': return 'Client'
    default: return role
  }
}

async function loadData() {
  if (!props.project) return
  isLoadingMembers.value = true
  grantErrorMessage.value = ''
  newMemberUserId.value = ''

  try {
    const [fetchedMembers, fetchedUsers] = await Promise.all([
      projectApi.getProjectMembers(props.project.id),
      userApi.getUsers().catch(() => [])
    ])
    members.value = fetchedMembers || []
    allUsers.value = fetchedUsers || []
  } catch (err: any) {
    store.showToast('Gagal memuat daftar anggota proyek')
  } finally {
    isLoadingMembers.value = false
  }
}

async function handleAddMember() {
  if (!props.project || !newMemberUserId.value) return
  isSubmitting.value = true
  grantErrorMessage.value = ''

  try {
    await projectApi.addProjectMember(props.project.id, {
      user_id: newMemberUserId.value,
      project_role: newMemberRole.value
    })
    store.showToast('Hak akses proyek berhasil diberikan!')
    newMemberUserId.value = ''
    await loadData()
    emit('updated')
  } catch (err: any) {
    grantErrorMessage.value = err.response?.data?.error || 'Gagal menambahkan anggota ke proyek'
  } finally {
    isSubmitting.value = false
  }
}

async function handleRemoveMember(userId: string, userName: string) {
  if (!props.project) return
  if (!confirm(`Cabut hak akses ${userName} dari proyek ${props.project.name}?`)) return

  isRevokingId.value = userId
  try {
    await projectApi.removeProjectMember(props.project.id, userId)
    store.showToast(`Hak akses ${userName} berhasil dicabut.`)
    await loadData()
    emit('updated')
  } catch (err: any) {
    store.showToast(err.response?.data?.error || 'Gagal mencabut hak akses.')
  } finally {
    isRevokingId.value = null
  }
}

function handleClose() {
  emit('close')
}

watch(() => props.isOpen, (open) => {
  if (open && props.project) {
    loadData()
  }
})
</script>
