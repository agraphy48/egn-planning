<template>
  <div class="min-h-screen w-full bg-gradient-to-br from-slate-50 via-white to-indigo-50/30 text-slate-900 flex items-center justify-center p-4 sm:p-6 relative overflow-hidden select-none">
    <!-- Ambient Subtle Pastel Glows -->
    <div class="absolute -top-32 -left-32 w-96 h-96 bg-brand-500/10 rounded-full blur-[100px] pointer-events-none"></div>
    <div class="absolute -bottom-32 -right-32 w-96 h-96 bg-indigo-500/10 rounded-full blur-[120px] pointer-events-none"></div>
    <div class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[600px] h-[600px] bg-purple-500/5 rounded-full blur-[140px] pointer-events-none"></div>

    <!-- Clean Tech Dot Grid Pattern -->
    <div class="absolute inset-0 bg-[radial-gradient(#cbd5e1_1px,transparent_1px)] [background-size:24px_24px] pointer-events-none opacity-60"></div>

    <!-- Main Card Container (Pure White Theme) -->
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
          Enterprise Project Management & Client Collaboration Portal
        </p>

        <!-- Subdomain & Cloudflare Edge Status Pill -->
        <div class="pt-1 flex items-center justify-center">
          <div class="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-emerald-50 border border-emerald-200 text-emerald-700 text-[11px] font-mono shadow-2xs">
            <span class="relative flex h-2 w-2">
              <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
              <span class="relative inline-flex rounded-full h-2 w-2 bg-emerald-500"></span>
            </span>
            <span class="font-medium">project.kedukanaksara.com</span>
            <span class="text-emerald-300">|</span>
            <span class="text-emerald-700 font-medium">Edge TLS 1.3 Active</span>
          </div>
        </div>
      </div>

      <!-- Active Session Banner (When user arrives at login but is already signed in) -->
      <div
        v-if="authStore.isAuthenticated"
        class="p-4 rounded-2xl bg-indigo-50/70 border border-indigo-200/80 text-xs space-y-3 animate-slide-up shadow-2xs"
      >
        <div class="flex items-center gap-3">
          <img
            :src="authStore.user?.avatar_url"
            :alt="authStore.user?.name"
            class="w-10 h-10 rounded-xl object-cover ring-2 ring-brand-500/40"
          />
          <div class="flex-1 min-w-0">
            <div class="flex items-center gap-2">
              <span class="font-bold text-slate-900 truncate">{{ authStore.user?.name }}</span>
              <span class="text-[9px] font-mono px-1.5 py-0.2 rounded bg-brand-100 text-brand-700 font-semibold uppercase">
                {{ authStore.currentRole }}
              </span>
            </div>
            <p class="text-[11px] text-slate-500 truncate">{{ authStore.user?.email }}</p>
          </div>
        </div>

        <div class="flex items-center gap-2 pt-1">
          <button
            @click="proceedRedirect(authStore.currentRole)"
            class="flex-1 py-2.5 px-3 rounded-xl bg-brand-600 hover:bg-brand-500 text-white font-bold text-xs shadow-md shadow-brand-600/25 flex items-center justify-center gap-1.5 transition-all cursor-pointer"
          >
            <span>Lanjutkan ke Workspace</span>
            <ArrowRight class="w-3.5 h-3.5" />
          </button>

          <button
            @click="handleLogoutSession"
            class="py-2.5 px-3 rounded-xl bg-white hover:bg-slate-50 border border-slate-200 text-slate-700 font-semibold text-xs shadow-2xs transition-all cursor-pointer"
          >
            Ganti Akun
          </button>
        </div>
      </div>

      <!-- Session Notice Notification Banner (Idle Timeout / Concurrent Login / Expired) -->
      <div
        v-if="noticeInfo"
        :class="['p-3 rounded-xl border text-xs flex items-center gap-2.5 animate-slide-up shadow-2xs', noticeInfo.bgClass]"
      >
        <component :is="noticeInfo.icon" class="w-4 h-4 flex-shrink-0" :class="noticeInfo.iconClass" />
        <span class="flex-1 leading-snug font-medium">{{ noticeInfo.message }}</span>
      </div>

      <!-- Error Alert Banner -->
      <div
        v-if="errorMessage"
        class="p-3 rounded-xl bg-rose-50 border border-rose-200 text-rose-800 text-xs flex items-center gap-2.5 animate-slide-up shadow-2xs"
      >
        <AlertCircle class="w-4 h-4 text-rose-500 flex-shrink-0" />
        <span class="flex-1 font-medium">{{ errorMessage }}</span>
        <button @click="errorMessage = ''" class="text-rose-500 hover:text-rose-700">
          <X class="w-3.5 h-3.5" />
        </button>
      </div>

      <!-- Optional Tab Switcher: Only displayed in Development Mode when explicitly enabled by server -->
      <div
        v-if="showQuickLoginTab"
        class="grid grid-cols-2 p-1 rounded-xl bg-slate-100 border border-slate-200 text-xs font-semibold"
      >
        <button
          @click="activeTab = 'credentials'"
          :class="[
            'py-2 px-3 rounded-lg flex items-center justify-center gap-2 transition-all duration-200',
            activeTab === 'credentials'
              ? 'bg-white text-brand-700 shadow-sm font-bold'
              : 'text-slate-600 hover:text-slate-900'
          ]"
        >
          <KeyRound class="w-3.5 h-3.5" />
          <span>Login Kredensial</span>
        </button>

        <button
          @click="activeTab = 'quick'"
          :class="[
            'py-2 px-3 rounded-lg flex items-center justify-center gap-2 transition-all duration-200',
            activeTab === 'quick'
              ? 'bg-white text-brand-700 shadow-sm font-bold'
              : 'text-slate-600 hover:text-slate-900'
          ]"
        >
          <Zap class="w-3.5 h-3.5" :class="activeTab === 'quick' ? 'text-amber-500' : ''" />
          <span>Akses Cepat Pengguna (Dev)</span>
        </button>
      </div>

      <!-- TAB 1: Primary Secure Credentials Form -->
      <form v-if="activeTab === 'credentials'" @submit.prevent="handleManualLogin" class="space-y-4 text-xs">
        <!-- Email Input -->
        <div class="space-y-1.5">
          <label class="font-semibold text-slate-700 flex items-center justify-between">
            <span>Alamat Email</span>
            <span class="text-[10px] text-slate-400 font-mono">Format: user@pulse.workspace</span>
          </label>
          <div class="relative">
            <Mail class="w-4 h-4 text-slate-400 absolute left-3 top-1/2 -translate-y-1/2" />
            <input
              type="email"
              v-model="email"
              required
              placeholder="nama@pulse.workspace"
              class="w-full pl-9 pr-3 py-2.5 rounded-xl border border-slate-300 bg-slate-50/50 text-slate-900 placeholder-slate-400 focus:bg-white focus:outline-none focus:ring-2 focus:ring-brand-500 focus:border-brand-500 transition-all text-xs shadow-xs"
            />
          </div>
        </div>

        <!-- Password Input with Toggle -->
        <div class="space-y-1.5">
          <div class="flex items-center justify-between">
            <label class="font-semibold text-slate-700">Kata Sandi (Password)</label>
            <span class="text-[10px] text-slate-400 font-mono">Sandi Akun Terverifikasi</span>
          </div>
          <div class="relative">
            <Lock class="w-4 h-4 text-slate-400 absolute left-3 top-1/2 -translate-y-1/2" />
            <input
              :type="showPassword ? 'text' : 'password'"
              v-model="password"
              required
              placeholder="••••••••••••"
              class="w-full pl-9 pr-10 py-2.5 rounded-xl border border-slate-300 bg-slate-50/50 text-slate-900 placeholder-slate-400 focus:bg-white focus:outline-none focus:ring-2 focus:ring-brand-500 focus:border-brand-500 transition-all text-xs shadow-xs"
            />
            <button
              type="button"
              @click="showPassword = !showPassword"
              class="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600 transition-colors"
            >
              <EyeOff v-if="showPassword" class="w-4 h-4" />
              <Eye v-else class="w-4 h-4" />
            </button>
          </div>
        </div>

        <!-- Development Autofill Assistance (Dev Mode Only) -->
        <div v-if="showQuickLoginTab" class="space-y-1 pt-1">
          <span class="text-[10px] font-mono text-slate-500 font-semibold uppercase">Auto-Fill Dev:</span>
          <div class="flex flex-wrap gap-1.5">
            <button
              v-for="user in quickUsers.slice(0, 5)"
              :key="user.key"
              type="button"
              @click="autofill(user.email)"
              class="px-2 py-1 rounded-md bg-slate-100 hover:bg-slate-200 text-slate-700 border border-slate-200 text-[10px] font-medium transition-colors"
            >
              {{ user.shortName }}
            </button>
            <button
              type="button"
              @click="autofill('admin@pulse.workspace')"
              class="px-2 py-1 rounded-md bg-rose-50 hover:bg-rose-100 text-rose-700 border border-rose-200 text-[10px] font-semibold transition-colors"
            >
              Admin
            </button>
          </div>
        </div>

        <!-- Submit Button -->
        <button
          type="submit"
          :disabled="authStore.isLoading || isSubmitting"
          class="w-full py-3 rounded-xl bg-gradient-to-r from-brand-600 to-indigo-600 hover:from-brand-500 hover:to-indigo-500 text-white font-bold text-xs shadow-lg shadow-brand-600/25 flex items-center justify-center gap-2 transition-all duration-200 active:scale-[0.98] disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer"
        >
          <Loader2 v-if="authStore.isLoading || isSubmitting" class="w-4 h-4 animate-spin" />
          <span v-else>Masuk ke Workspace</span>
          <ArrowRight v-if="!authStore.isLoading && !isSubmitting" class="w-4 h-4" />
        </button>
      </form>

      <!-- TAB 2: Development Quick Login (Only when enabled by server) -->
      <div v-else-if="showQuickLoginTab" class="space-y-3">
        <div class="flex items-center justify-between text-[11px] text-slate-500 font-medium px-1">
          <span>Pilih Profil untuk Akses Pengembang:</span>
          <span class="font-mono text-slate-400 text-[10px]">Autentikasi API Riil</span>
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-2.5 max-h-[380px] overflow-y-auto pr-1">
          <button
            v-for="user in quickUsers"
            :key="user.key"
            @click="handleQuickLogin(user.key)"
            :disabled="authStore.isLoading || isSubmitting"
            :class="[
              'p-3 rounded-xl border text-left transition-all duration-200 flex items-start gap-3 relative group',
              user.borderClass,
              'bg-slate-50/70 hover:bg-white hover:shadow-md active:scale-[0.98]'
            ]"
          >
            <!-- User Avatar with Status Indicator -->
            <div class="relative flex-shrink-0 mt-0.5">
              <img
                :src="user.avatar"
                :alt="user.name"
                class="w-10 h-10 rounded-xl object-cover ring-1 ring-slate-200 group-hover:ring-brand-500 transition-all"
              />
              <span class="absolute -bottom-0.5 -right-0.5 w-2.5 h-2.5 rounded-full bg-emerald-500 ring-2 ring-white"></span>
            </div>

            <!-- User Info & Role Badge -->
            <div class="flex-1 min-w-0">
              <div class="flex items-center justify-between gap-1">
                <span class="text-xs font-bold text-slate-900 group-hover:text-brand-600 transition-colors truncate">
                  {{ user.name }}
                </span>
                <span :class="['text-[9px] font-mono px-1.5 py-0.2 rounded font-semibold uppercase', user.badgeClass]">
                  {{ user.badgeText }}
                </span>
              </div>
              <p class="text-[10px] text-slate-500 truncate mt-0.5">{{ user.email }}</p>
              <p class="text-[10px] text-slate-400 truncate mt-0.5 italic">{{ user.roleDesc }}</p>
            </div>

            <!-- Hover Slide Arrow -->
            <div class="self-center text-slate-400 group-hover:text-brand-600 group-hover:translate-x-0.5 transition-all">
              <ChevronRight class="w-4 h-4" />
            </div>
          </button>
        </div>
      </div>

      <!-- Footer Info & Security Verification -->
      <div class="pt-4 border-t border-slate-200 flex flex-col sm:flex-row items-center justify-between text-[11px] text-slate-500 gap-2">
        <div class="flex items-center gap-1.5">
          <ShieldCheck class="w-3.5 h-3.5 text-emerald-600" />
          <span>Role-Based Access Control (RBAC) Hardened</span>
        </div>
        <div class="flex items-center gap-1 font-mono text-[10px] text-slate-500">
          <span>Target Host:</span>
          <span class="text-brand-600 font-semibold">127.0.0.1:5080</span>
        </div>
      </div>

    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore, SEED_USERS } from '../stores/authStore'
import {
  Layers,
  Zap,
  KeyRound,
  Mail,
  Lock,
  Eye,
  EyeOff,
  ArrowRight,
  ShieldCheck,
  AlertCircle,
  X,
  ChevronRight,
  Loader2,
  Info,
  Clock,
  ShieldAlert
} from 'lucide-vue-next'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

// Default tab is always the real credential login form
const activeTab = ref<'credentials' | 'quick'>('credentials')
const email = ref('agus.setiawan@pulse.workspace')
const password = ref('Password123!')
const showPassword = ref(false)
const errorMessage = ref('')
const isSubmitting = ref(false)

// Session notice indicator from route query (styled for white theme)
const noticeInfo = computed(() => {
  const reason = route.query.reason
  if (reason === 'idle_timeout') {
    return {
      type: 'warning',
      message: 'Sesi Anda telah berakhir otomatis karena tidak ada aktivitas selama 15 menit. Silakan masuk kembali.',
      bgClass: 'bg-amber-50 border-amber-200 text-amber-900',
      icon: Clock,
      iconClass: 'text-amber-600',
    }
  }
  if (reason === 'concurrent_device') {
    return {
      type: 'danger',
      message: 'Akun Anda telah dibuka di perangkat atau jendela lain. Sesi pada perangkat ini telah dihentikan demi keamanan.',
      bgClass: 'bg-rose-50 border-rose-200 text-rose-900',
      icon: ShieldAlert,
      iconClass: 'text-rose-600',
    }
  }
  if (reason === 'session_expired') {
    return {
      type: 'info',
      message: 'Sesi autentikasi Anda telah berakhir. Silakan masuk kembali.',
      bgClass: 'bg-slate-100 border-slate-200 text-slate-800',
      icon: Info,
      iconClass: 'text-slate-600',
    }
  }
  return null
})

// Quick Login tab is strictly controlled: disabled in production or if server says false
const showQuickLoginTab = computed(() => {
  return !import.meta.env.PROD && authStore.allowQuickLogin
})

onMounted(async () => {
  await authStore.fetchAuthConfig()
  // Ensure we stay on credentials form if quick login is disabled
  if (!showQuickLoginTab.value) {
    activeTab.value = 'credentials'
  }
})

interface QuickUserItem {
  key: string
  name: string
  shortName: string
  email: string
  badgeText: string
  badgeClass: string
  borderClass: string
  avatar: string
  roleDesc: string
}

const quickUsers: QuickUserItem[] = [
  {
    key: 'agus',
    name: 'Agus Setiawan',
    shortName: 'Agus',
    email: 'agus.setiawan@pulse.workspace',
    badgeText: 'Project Manager',
    badgeClass: 'bg-emerald-50 text-emerald-700 border border-emerald-200',
    borderClass: 'border-slate-200 hover:border-emerald-500',
    avatar: SEED_USERS['agus'].avatar_url,
    roleDesc: 'Sprint lead, delegasi task & timeline'
  },
  {
    key: 'adi',
    name: 'Adi Saputra',
    shortName: 'Adi',
    email: 'adi.saputra@pulse.workspace',
    badgeText: 'Project Manager',
    badgeClass: 'bg-emerald-50 text-emerald-700 border border-emerald-200',
    borderClass: 'border-slate-200 hover:border-emerald-500',
    avatar: SEED_USERS['adi'].avatar_url,
    roleDesc: 'Alokasi resource & project governance'
  },
  {
    key: 'yudistira',
    name: 'Yudistira Sira P.',
    shortName: 'Yudistira',
    email: 'yudistira@pulse.workspace',
    badgeText: 'Lead Engineer',
    badgeClass: 'bg-indigo-50 text-indigo-700 border border-indigo-200',
    borderClass: 'border-slate-200 hover:border-indigo-500',
    avatar: SEED_USERS['yudistira'].avatar_url,
    roleDesc: 'Fullstack Golang Gin & Vue 3 SPA'
  },
  {
    key: 'hardianto',
    name: 'Hardianto',
    shortName: 'Hardianto',
    email: 'hardianto@pulse.workspace',
    badgeText: 'Frontend Dev',
    badgeClass: 'bg-cyan-50 text-cyan-700 border border-cyan-200',
    borderClass: 'border-slate-200 hover:border-cyan-500',
    avatar: SEED_USERS['hardianto'].avatar_url,
    roleDesc: 'UI/UX animatif & responsif Tailwind'
  },
  {
    key: 'dedi',
    name: 'Dedi Irawan',
    shortName: 'Dedi',
    email: 'dedi.irawan@pulse.workspace',
    badgeText: 'DevOps & QA',
    badgeClass: 'bg-purple-50 text-purple-700 border border-purple-200',
    borderClass: 'border-slate-200 hover:border-purple-500',
    avatar: SEED_USERS['dedi'].avatar_url,
    roleDesc: 'Cloudflare Tunnel, TLS & Security Audit'
  },
  {
    key: 'super_admin',
    name: 'Budi Santoso',
    shortName: 'Admin',
    email: 'admin@pulse.workspace',
    badgeText: 'Super Admin',
    badgeClass: 'bg-rose-50 text-rose-700 border border-rose-200',
    borderClass: 'border-slate-200 hover:border-rose-500',
    avatar: SEED_USERS['super_admin'].avatar_url,
    roleDesc: 'Kontrol penuh sistem, audit & pengguna'
  },
  {
    key: 'client',
    name: 'Ir. H. Gunawan',
    shortName: 'Client',
    email: 'client@pulse.workspace',
    badgeText: 'Client Portal',
    badgeClass: 'bg-amber-50 text-amber-700 border border-amber-200',
    borderClass: 'border-slate-200 hover:border-amber-500',
    avatar: SEED_USERS['client'].avatar_url,
    roleDesc: 'Persetujuan milestone & review deliverables'
  }
]

function autofill(userEmail: string) {
  email.value = userEmail
  password.value = 'Password123!'
  errorMessage.value = ''
}

async function proceedRedirect(role?: string) {
  if (route.query.redirect) {
    router.push(route.query.redirect as string)
    return
  }

  if (role === 'client') {
    router.push('/client-portal')
  } else {
    router.push('/dashboard')
  }
}

async function handleQuickLogin(userKey: string) {
  errorMessage.value = ''
  isSubmitting.value = true
  try {
    const user = await authStore.switchUser(userKey)
    await proceedRedirect(user.role)
  } catch (err: any) {
    errorMessage.value = err.message || 'Gagal masuk akun'
  } finally {
    isSubmitting.value = false
  }
}

async function handleLogoutSession() {
  await authStore.logout()
  errorMessage.value = ''
}

async function handleManualLogin() {
  errorMessage.value = ''
  isSubmitting.value = true
  try {
    await authStore.login(email.value, password.value)
    await proceedRedirect(authStore.currentRole)
  } catch (err: any) {
    errorMessage.value = err?.response?.data?.error || err.message || 'Email atau password salah'
  } finally {
    isSubmitting.value = false
  }
}
</script>
