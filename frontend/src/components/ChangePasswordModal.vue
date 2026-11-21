<template>
  <div
    v-if="authStore.isChangePasswordModalOpen"
    class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/50 backdrop-blur-xs select-none animate-in fade-in duration-150"
    @click.self="handleClose"
  >
    <div
      class="w-full max-w-md bg-white rounded-3xl border border-slate-200/90 shadow-2xl p-6 space-y-5 relative animate-in zoom-in-95 duration-150"
    >
      <!-- Top Decorative Accent -->
      <div class="absolute top-0 left-8 right-8 h-1 bg-gradient-to-r from-transparent via-brand-600 to-transparent rounded-full"></div>

      <!-- Header with Close Button -->
      <div class="flex items-start justify-between">
        <div class="flex items-center gap-3">
          <div class="p-2.5 rounded-2xl bg-brand-50 border border-brand-100 text-brand-600 shadow-2xs">
            <KeyRound class="w-5 h-5 stroke-[2.2]" />
          </div>
          <div>
            <h2 class="text-base font-bold text-slate-900 font-display">Ganti Kata Sandi</h2>
            <p class="text-xs text-slate-500 font-medium">Perbarui kata sandi akun Pulse Workspace Anda</p>
          </div>
        </div>
        <button
          @click="handleClose"
          class="p-1.5 rounded-xl text-slate-400 hover:text-slate-700 hover:bg-slate-100 transition-colors cursor-pointer"
        >
          <X class="w-4 h-4" />
        </button>
      </div>

      <!-- Success Notification State -->
      <div
        v-if="successMessage"
        class="p-4 rounded-2xl bg-emerald-50 border border-emerald-200 text-emerald-800 text-xs space-y-3 animate-slide-up"
      >
        <div class="flex items-center gap-2.5">
          <CheckCircle2 class="w-5 h-5 text-emerald-600 flex-shrink-0" />
          <span class="font-semibold leading-relaxed">{{ successMessage }}</span>
        </div>
        <button
          @click="handleClose"
          class="w-full py-2 px-3 rounded-xl bg-emerald-600 hover:bg-emerald-500 text-white font-bold text-xs shadow-md shadow-emerald-600/20 transition-all cursor-pointer"
        >
          Selesai & Tutup
        </button>
      </div>

      <!-- Change Password Form -->
      <form v-else @submit.prevent="handleSubmit" class="space-y-4 text-xs">
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

        <!-- Current Password -->
        <div class="space-y-1.5">
          <label class="font-semibold text-slate-700 flex items-center justify-between">
            <span>Kata Sandi Saat Ini</span>
            <span class="text-[10px] text-slate-400 font-mono">Verifikasi Keamanan</span>
          </label>
          <div class="relative">
            <Lock class="w-4 h-4 text-slate-400 absolute left-3 top-1/2 -translate-y-1/2" />
            <input
              :type="showCurrent ? 'text' : 'password'"
              v-model="currentPassword"
              required
              autocomplete="current-password"
              placeholder="Masukkan kata sandi saat ini..."
              class="w-full pl-9 pr-10 py-2.5 rounded-xl border border-slate-300 bg-slate-50/50 text-slate-900 placeholder-slate-400 focus:bg-white focus:outline-none focus:ring-2 focus:ring-brand-500 focus:border-brand-500 transition-all text-xs shadow-xs"
            />
            <button
              type="button"
              @click="showCurrent = !showCurrent"
              class="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600 transition-colors cursor-pointer"
            >
              <EyeOff v-if="showCurrent" class="w-4 h-4" />
              <Eye v-else class="w-4 h-4" />
            </button>
          </div>
        </div>

        <!-- New Password -->
        <div class="space-y-1.5">
          <label class="font-semibold text-slate-700 flex items-center justify-between">
            <span>Kata Sandi Baru</span>
            <span class="text-[10px] text-slate-400 font-mono">Min. 6 karakter</span>
          </label>
          <div class="relative">
            <Lock class="w-4 h-4 text-slate-400 absolute left-3 top-1/2 -translate-y-1/2" />
            <input
              :type="showNew ? 'text' : 'password'"
              v-model="newPassword"
              required
              autocomplete="new-password"
              placeholder="Minimal 6 karakter baru..."
              class="w-full pl-9 pr-10 py-2.5 rounded-xl border border-slate-300 bg-slate-50/50 text-slate-900 placeholder-slate-400 focus:bg-white focus:outline-none focus:ring-2 focus:ring-brand-500 focus:border-brand-500 transition-all text-xs shadow-xs"
            />
            <button
              type="button"
              @click="showNew = !showNew"
              class="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600 transition-colors cursor-pointer"
            >
              <EyeOff v-if="showNew" class="w-4 h-4" />
              <Eye v-else class="w-4 h-4" />
            </button>
          </div>
        </div>

        <!-- Confirm New Password -->
        <div class="space-y-1.5">
          <label class="font-semibold text-slate-700 flex items-center justify-between">
            <span>Konfirmasi Kata Sandi Baru</span>
            <span
              v-if="confirmPassword"
              :class="[
                'text-[10px] font-mono font-semibold',
                passwordsMatch ? 'text-emerald-600' : 'text-rose-500'
              ]"
            >
              {{ passwordsMatch ? 'Cocok' : 'Belum Cocok' }}
            </span>
          </label>
          <div class="relative">
            <Lock class="w-4 h-4 text-slate-400 absolute left-3 top-1/2 -translate-y-1/2" />
            <input
              :type="showConfirm ? 'text' : 'password'"
              v-model="confirmPassword"
              required
              autocomplete="new-password"
              placeholder="Ketik ulang kata sandi baru..."
              class="w-full pl-9 pr-10 py-2.5 rounded-xl border border-slate-300 bg-slate-50/50 text-slate-900 placeholder-slate-400 focus:bg-white focus:outline-none focus:ring-2 focus:ring-brand-500 focus:border-brand-500 transition-all text-xs shadow-xs"
            />
            <button
              type="button"
              @click="showConfirm = !showConfirm"
              class="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600 transition-colors cursor-pointer"
            >
              <EyeOff v-if="showConfirm" class="w-4 h-4" />
              <Eye v-else class="w-4 h-4" />
            </button>
          </div>
        </div>

        <!-- Password Requirements Checklist -->
        <div class="p-3 rounded-2xl bg-slate-50 border border-slate-200/80 space-y-1.5 text-[11px]">
          <div class="flex items-center gap-2" :class="newPassword.length >= 6 ? 'text-emerald-700' : 'text-slate-500'">
            <Check class="w-3.5 h-3.5" :class="newPassword.length >= 6 ? 'text-emerald-600' : 'text-slate-300'" />
            <span>Panjang minimal 6 karakter</span>
          </div>
          <div class="flex items-center gap-2" :class="passwordsMatch && newPassword ? 'text-emerald-700' : 'text-slate-500'">
            <Check class="w-3.5 h-3.5" :class="passwordsMatch && newPassword ? 'text-emerald-600' : 'text-slate-300'" />
            <span>Konfirmasi kata sandi cocok</span>
          </div>
          <div class="flex items-center gap-2" :class="newPassword && currentPassword && newPassword !== currentPassword ? 'text-emerald-700' : 'text-slate-500'">
            <Check class="w-3.5 h-3.5" :class="newPassword && currentPassword && newPassword !== currentPassword ? 'text-emerald-600' : 'text-slate-300'" />
            <span>Berbeda dari kata sandi lama</span>
          </div>
        </div>

        <!-- Action Buttons -->
        <div class="flex items-center gap-2 pt-2">
          <button
            type="button"
            @click="handleClose"
            :disabled="isSubmitting"
            class="flex-1 py-2.5 px-3 rounded-xl bg-slate-100 hover:bg-slate-200/80 text-slate-700 font-semibold transition-colors cursor-pointer"
          >
            Batal
          </button>

          <button
            type="submit"
            :disabled="isSubmitting || !isFormValid"
            class="flex-1 py-2.5 px-3 rounded-xl bg-gradient-to-r from-brand-600 to-indigo-600 hover:from-brand-500 hover:to-indigo-500 text-white font-bold shadow-md shadow-brand-600/25 flex items-center justify-center gap-1.5 transition-all duration-200 active:scale-[0.98] disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer"
          >
            <Loader2 v-if="isSubmitting" class="w-4 h-4 animate-spin" />
            <span v-else>Simpan Perubahan</span>
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useAuthStore } from '../stores/authStore'
import { useProjectStore } from '../stores/projectStore'
import {
  KeyRound,
  Lock,
  Eye,
  EyeOff,
  Check,
  CheckCircle2,
  AlertCircle,
  X,
  Loader2
} from 'lucide-vue-next'

const authStore = useAuthStore()
const projectStore = useProjectStore()

const currentPassword = ref('')
const newPassword = ref('')
const confirmPassword = ref('')

const showCurrent = ref(false)
const showNew = ref(false)
const showConfirm = ref(false)

const isSubmitting = ref(false)
const errorMessage = ref('')
const successMessage = ref('')

const passwordsMatch = computed(() => {
  return newPassword.value !== '' && newPassword.value === confirmPassword.value
})

const isFormValid = computed(() => {
  return (
    currentPassword.value.length > 0 &&
    newPassword.value.length >= 6 &&
    passwordsMatch.value &&
    newPassword.value !== currentPassword.value
  )
})

function handleClose() {
  authStore.isChangePasswordModalOpen = false
  currentPassword.value = ''
  newPassword.value = ''
  confirmPassword.value = ''
  errorMessage.value = ''
  successMessage.value = ''
  isSubmitting.value = false
}

async function handleSubmit() {
  if (!isFormValid.value) return

  errorMessage.value = ''
  successMessage.value = ''
  isSubmitting.value = true

  try {
    const res = await authStore.changePassword(
      currentPassword.value,
      newPassword.value,
      confirmPassword.value
    )
    successMessage.value = res?.message || 'Kata sandi berhasil diperbarui.'
    projectStore.showToast('Kata sandi berhasil diperbarui')
  } catch (err: any) {
    errorMessage.value =
      err?.response?.data?.error || err.message || 'Gagal memperbarui kata sandi. Periksa kata sandi lama Anda.'
  } finally {
    isSubmitting.value = false
  }
}
</script>
