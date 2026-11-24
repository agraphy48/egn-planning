<template>
  <div class="min-h-screen w-full bg-gradient-to-br from-slate-50 via-white to-indigo-50/30 text-slate-900 flex items-center justify-center p-4 sm:p-6 relative overflow-hidden select-none">
    <!-- Ambient Subtle Pastel Glows -->
    <div class="absolute -top-32 -left-32 w-96 h-96 bg-brand-500/10 rounded-full blur-[100px] pointer-events-none"></div>
    <div class="absolute -bottom-32 -right-32 w-96 h-96 bg-indigo-500/10 rounded-full blur-[120px] pointer-events-none"></div>
    <div class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[600px] h-[600px] bg-purple-500/5 rounded-full blur-[140px] pointer-events-none"></div>

    <!-- Clean Tech Dot Grid Pattern -->
    <div class="absolute inset-0 bg-[radial-gradient(#cbd5e1_1px,transparent_1px)] [background-size:24px_24px] pointer-events-none opacity-60"></div>

    <!-- Main Card Container -->
    <div class="max-w-xl w-full rounded-3xl bg-white/95 backdrop-blur-xl border border-slate-200/90 shadow-2xl shadow-slate-200/60 p-6 sm:p-8 space-y-6 relative z-10 transition-all duration-300">
      
      <!-- Top Decorative Glowing Line -->
      <div class="absolute top-0 left-8 right-8 h-1 bg-gradient-to-r from-transparent via-brand-600 to-transparent rounded-full"></div>

      <!-- Header Section -->
      <div class="text-center space-y-2">
        <div class="inline-flex items-center justify-center p-3 rounded-2xl bg-gradient-to-tr from-brand-600 to-indigo-600 shadow-lg shadow-brand-600/20 text-white mb-1 group hover:scale-105 transition-transform duration-300">
          <Layers class="w-6 h-6 stroke-[2.2]" />
        </div>
        <h1 class="text-2xl font-extrabold tracking-tight text-slate-900 font-display">
          Pulse Workspace
        </h1>
        <p class="text-xs text-slate-500 font-medium">
          Aktivasi & Konfirmasi Akun Kolaborasi
        </p>

        <!-- Subdomain Pill -->
        <div class="pt-1 flex items-center justify-center">
          <div class="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-emerald-50 border border-emerald-200 text-emerald-700 text-[11px] font-mono shadow-2xs">
            <span class="relative flex h-2 w-2">
              <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
              <span class="relative inline-flex rounded-full h-2 w-2 bg-emerald-500"></span>
            </span>
            <span>project.kedukanaksara.com</span>
            <span class="text-emerald-300">|</span>
            <span>Aktivasi Resmi</span>
          </div>
        </div>
      </div>

      <!-- State 1: Loading Confirmation -->
      <div v-if="isLoading" class="py-12 flex flex-col items-center justify-center space-y-4 text-center">
        <div class="relative flex items-center justify-center">
          <div class="w-12 h-12 rounded-full border-4 border-brand-100 border-t-brand-600 animate-spin"></div>
          <Mail class="w-5 h-5 text-brand-600 absolute" />
        </div>
        <div class="space-y-1">
          <h3 class="text-sm font-bold text-slate-800">Mengonfirmasi Akun Anda...</h3>
          <p class="text-xs text-slate-500">
            Sistem sedang memverifikasi tautan undangan, mengaktifkan akun, dan menyiapkan kredensial login.
          </p>
        </div>
      </div>

      <!-- State 2: Error (Invalid / Expired / Revoked) -->
      <div v-else-if="errorMessage" class="space-y-5 py-2 animate-slide-up text-xs">
        <div class="p-5 rounded-2xl bg-rose-50 border border-rose-200 text-rose-900 space-y-3">
          <div class="flex items-start gap-3">
            <div class="p-2 rounded-xl bg-rose-100 text-rose-600 flex-shrink-0 mt-0.5">
              <AlertCircle class="w-5 h-5" />
            </div>
            <div class="space-y-1 flex-1">
              <h3 class="font-bold text-sm text-rose-950">Tautan Konfirmasi Tidak Dapat Digunakan</h3>
              <p class="text-rose-800 leading-relaxed text-xs">
                {{ errorMessage }}
              </p>
            </div>
          </div>
          <div class="pl-11 text-[11px] text-rose-700/90">
            <p>Tips: Jika Anda menerima beberapa undangan, silakan pastikan Anda mengklik tautan dari email undangan yang paling baru diterima.</p>
          </div>
        </div>

        <div class="pt-2">
          <RouterLink
            to="/login"
            class="w-full py-3 px-4 rounded-xl bg-slate-100 hover:bg-slate-200 text-slate-800 font-bold text-xs flex items-center justify-center gap-2 transition-colors cursor-pointer"
          >
            <ArrowLeft class="w-4 h-4" />
            <span>Kembali ke Halaman Login</span>
          </RouterLink>
        </div>
      </div>

      <!-- State 3: Already Confirmed -->
      <div v-else-if="isAlreadyConfirmed" class="space-y-5 py-2 animate-slide-up text-xs">
        <div class="p-5 rounded-2xl bg-emerald-50 border border-emerald-200 text-emerald-950 space-y-3 text-center">
          <div class="w-12 h-12 rounded-full bg-emerald-100 text-emerald-600 flex items-center justify-center mx-auto shadow-xs">
            <CheckCircle2 class="w-7 h-7" />
          </div>
          <h3 class="text-base font-bold text-emerald-950">Akun Anda Sudah Aktif</h3>
          <p class="text-xs text-emerald-800 leading-relaxed max-w-md mx-auto">
            Akun dengan email <strong class="font-mono text-emerald-900">{{ confirmData?.email }}</strong> telah dikonfirmasi dan aktif sebelumnya.
            Username dan kata sandi telah dikirimkan ke email Anda.
          </p>
        </div>

        <div class="space-y-3 pt-2">
          <button
            @click="goToLogin"
            class="w-full py-3.5 rounded-xl bg-gradient-to-r from-brand-600 to-indigo-600 hover:from-brand-500 text-white font-bold text-xs shadow-lg shadow-brand-600/25 flex items-center justify-center gap-2 transition-all cursor-pointer"
          >
            <span>Masuk ke Halaman Login</span>
            <ArrowRight class="w-4 h-4" />
          </button>
        </div>
      </div>

      <!-- State 4: Success - Newly Confirmed with Auto-Generated Credentials -->
      <div v-else-if="confirmData" class="space-y-5 animate-slide-up text-xs">
        
        <!-- Big Success Notice Banner -->
        <div class="p-5 rounded-2xl bg-gradient-to-br from-emerald-50 to-teal-50/60 border border-emerald-200 text-emerald-950 space-y-3 text-center">
          <div class="w-12 h-12 rounded-full bg-emerald-100 text-emerald-600 flex items-center justify-center mx-auto shadow-xs">
            <CheckCircle2 class="w-7 h-7" />
          </div>
          <div class="space-y-1">
            <h3 class="text-lg font-extrabold text-emerald-950 tracking-tight">Konfirmasi Akun Berhasil!</h3>
            <p class="text-xs text-emerald-800 leading-relaxed max-w-md mx-auto">
              Akun Anda telah aktif. Username dan Password Anda telah <strong>dikirimkan secara otomatis ke email Anda</strong>:
            </p>
            <div class="inline-block mt-1 px-3 py-1 rounded-lg bg-emerald-100/80 border border-emerald-200 font-mono font-bold text-emerald-950 text-xs">
              {{ confirmData.email }}
            </div>
          </div>
        </div>

        <!-- Credential Summary Card -->
        <div class="p-4 rounded-2xl bg-slate-50/90 border border-slate-200/90 space-y-3">
          <div class="flex items-center justify-between pb-2 border-b border-slate-200">
            <span class="text-[10px] font-mono uppercase text-slate-400 font-bold tracking-wider">Kredensial Login Anda</span>
            <span :class="['px-2.5 py-0.5 rounded-full text-[10px] font-mono font-bold uppercase border', roleBadgeStyle]">
              {{ roleLabel }}
            </span>
          </div>

          <!-- Email / Username -->
          <div class="space-y-1">
            <label class="text-[10px] text-slate-400 font-mono block">Username / Email Login:</label>
            <div class="flex items-center justify-between p-2.5 rounded-xl bg-white border border-slate-200 font-mono text-xs">
              <span class="font-bold text-slate-900 select-all">{{ confirmData.email }}</span>
              <button
                @click="copyText(confirmData.email, 'email')"
                class="px-2.5 py-1 rounded-lg bg-slate-100 hover:bg-slate-200 text-slate-700 font-sans text-[11px] font-medium flex items-center gap-1 transition-colors cursor-pointer"
                title="Salin Email"
              >
                <Check v-if="copiedField === 'email'" class="w-3.5 h-3.5 text-emerald-600" />
                <Copy v-else class="w-3.5 h-3.5" />
                <span>{{ copiedField === 'email' ? 'Tersalin' : 'Salin' }}</span>
              </button>
            </div>
          </div>

          <!-- Generated Password -->
          <div v-if="confirmData.password" class="space-y-1">
            <div class="flex items-center justify-between">
              <label class="text-[10px] text-slate-400 font-mono block">Kata Sandi (Password):</label>
              <span class="text-[10px] text-amber-600 font-medium">Juga dikirimkan ke email Anda</span>
            </div>
            <div class="flex items-center justify-between p-2.5 rounded-xl bg-amber-50/70 border border-amber-200 font-mono text-xs">
              <div class="flex items-center gap-2">
                <span class="font-bold text-amber-950 text-sm tracking-wider select-all">
                  {{ showPassword ? confirmData.password : '••••••••••••' }}
                </span>
                <button
                  type="button"
                  @click="showPassword = !showPassword"
                  class="text-amber-700 hover:text-amber-900 transition-colors cursor-pointer p-0.5"
                  :title="showPassword ? 'Sembunyikan' : 'Tampilkan'"
                >
                  <EyeOff v-if="showPassword" class="w-3.5 h-3.5" />
                  <Eye v-else class="w-3.5 h-3.5" />
                </button>
              </div>
              <button
                @click="copyText(confirmData.password, 'password')"
                class="px-2.5 py-1 rounded-lg bg-white border border-amber-300 hover:bg-amber-100 text-amber-900 font-sans text-[11px] font-semibold flex items-center gap-1 transition-colors cursor-pointer"
                title="Salin Kata Sandi"
              >
                <Check v-if="copiedField === 'password'" class="w-3.5 h-3.5 text-emerald-600" />
                <Copy v-else class="w-3.5 h-3.5" />
                <span>{{ copiedField === 'password' ? 'Tersalin' : 'Salin' }}</span>
              </button>
            </div>
          </div>

          <!-- Assigned Projects -->
          <div v-if="confirmData.projects && confirmData.projects.length > 0" class="pt-2 border-t border-slate-200 space-y-1.5">
            <label class="text-[10px] text-slate-400 font-mono block">Proyek yang Diberikan Akses:</label>
            <div class="flex flex-wrap gap-1.5">
              <span
                v-for="(pName, idx) in confirmData.projects"
                :key="idx"
                class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-lg bg-white border border-slate-200 text-slate-800 text-[11px] font-medium shadow-2xs"
              >
                <span class="w-2 h-2 rounded-full bg-brand-500"></span>
                <span>{{ pName }}</span>
              </span>
            </div>
          </div>
        </div>

        <!-- Security / Info Notice -->
        <div class="p-3 rounded-xl bg-blue-50/80 border border-blue-200/80 text-blue-900 text-[11px] flex items-start gap-2.5">
          <Info class="w-4 h-4 text-blue-600 flex-shrink-0 mt-0.5" />
          <div class="space-y-0.5 leading-relaxed">
            <p><strong>Tips Keamanan:</strong> Kredensial di atas telah dikirimkan ke kotak masuk email Anda. Anda dapat mengubah kata sandi kapan saja melalui menu profil setelah berhasil login.</p>
          </div>
        </div>

        <!-- Call to Action Buttons -->
        <div class="space-y-2 pt-2">
          <button
            @click="goToLogin"
            class="w-full py-3.5 rounded-xl bg-gradient-to-r from-brand-600 to-indigo-600 hover:from-brand-500 hover:to-indigo-500 text-white font-bold text-xs shadow-lg shadow-brand-600/25 flex items-center justify-center gap-2 transition-all duration-200 active:scale-[0.98] cursor-pointer"
          >
            <span>Masuk ke Halaman Login</span>
            <ArrowRight class="w-4 h-4" />
          </button>

          <button
            v-if="confirmData.password"
            @click="copyAllCredentials"
            class="w-full py-2.5 rounded-xl bg-white hover:bg-slate-50 border border-slate-200 text-slate-700 font-semibold text-xs shadow-2xs flex items-center justify-center gap-2 transition-all cursor-pointer"
          >
            <Check v-if="copiedField === 'all'" class="w-3.5 h-3.5 text-emerald-600" />
            <Copy v-else class="w-3.5 h-3.5 text-slate-500" />
            <span>{{ copiedField === 'all' ? 'Semua Kredensial Berhasil Disalin!' : 'Salin Semua Kredensial (Email & Password)' }}</span>
          </button>
        </div>

      </div>

      <!-- Footer Info -->
      <div class="pt-4 border-t border-slate-200 flex items-center justify-between text-[11px] text-slate-500">
        <div class="flex items-center gap-1.5">
          <ShieldCheck class="w-3.5 h-3.5 text-emerald-600" />
          <span>Role-Based Access Control Hardened</span>
        </div>
        <RouterLink to="/login" class="text-brand-600 hover:underline font-medium">
          Halaman Login
        </RouterLink>
      </div>

    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { invitationApi } from '../services/api'
import {
  Layers,
  Mail,
  CheckCircle2,
  AlertCircle,
  ArrowRight,
  ArrowLeft,
  ShieldCheck,
  Copy,
  Check,
  Eye,
  EyeOff,
  Info
} from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()

const token = ref('')
const isLoading = ref(true)
const errorMessage = ref('')
const confirmData = ref<any>(null)
const isAlreadyConfirmed = ref(false)
const showPassword = ref(true)
const copiedField = ref('')

const roleLabel = computed(() => {
  const r = confirmData.value?.role
  if (r === 'client') return 'Client Portal Access'
  if (r === 'pm') return 'Project Manager'
  if (r === 'member') return 'Team Member'
  return r || 'User'
})

const roleBadgeStyle = computed(() => {
  const r = confirmData.value?.role
  if (r === 'client') return 'bg-amber-50 text-amber-800 border-amber-200'
  if (r === 'pm') return 'bg-purple-50 text-purple-800 border-purple-200'
  if (r === 'member') return 'bg-blue-50 text-blue-800 border-blue-200'
  return 'bg-slate-100 text-slate-800 border-slate-200'
})

onMounted(async () => {
  token.value = (route.query.token as string) || ''
  if (!token.value) {
    isLoading.value = false
    errorMessage.value = 'Token undangan tidak ditemukan. Pastikan Anda mengklik tautan lengkap dari email atau pesan undangan.'
    return
  }

  try {
    const res = await invitationApi.confirmInvitation(token.value)
    confirmData.value = res
    if (res.status === 'already_confirmed') {
      isAlreadyConfirmed.value = true
    }
  } catch (err: any) {
    errorMessage.value =
      err?.response?.data?.error ||
      'Tautan undangan tidak valid atau telah kedaluwarsa. Silakan hubungi Administrator untuk meminta undangan baru.'
  } finally {
    isLoading.value = false
  }
})

function goToLogin() {
  const emailParam = confirmData.value?.email ? encodeURIComponent(confirmData.value.email) : ''
  router.push(`/login?email=${emailParam}&confirmed=true`)
}

async function copyText(text: string, field: string) {
  try {
    await navigator.clipboard.writeText(text)
    copiedField.value = field
    setTimeout(() => {
      if (copiedField.value === field) copiedField.value = ''
    }, 2000)
  } catch (err) {
    console.error('Failed to copy text', err)
  }
}

async function copyAllCredentials() {
  if (!confirmData.value) return
  const text = `Username/Email: ${confirmData.value.email}\nPassword: ${confirmData.value.password || '-'}\nLogin URL: https://project.kedukanaksara.com/login`
  await copyText(text, 'all')
}
</script>
