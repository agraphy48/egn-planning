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
            @click="authStore.logout()"
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
        <button @click="errorMessage = ''" class="text-rose-500 hover:text-rose-700 cursor-pointer">
          <X class="w-3.5 h-3.5" />
        </button>
      </div>

      <!-- Primary Secure Credentials Form -->
      <form @submit.prevent="handleManualLogin" class="space-y-4 text-xs" autocomplete="on">
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
              autocomplete="username"
              placeholder="nama@pulse.workspace"
              class="w-full pl-9 pr-3 py-2.5 rounded-xl border border-slate-300 bg-slate-50/50 text-slate-900 placeholder-slate-400 focus:bg-white focus:outline-none focus:ring-2 focus:ring-brand-500 focus:border-brand-500 transition-all text-xs shadow-xs"
            />
          </div>
        </div>

        <!-- Password Input with Toggle -->
        <div class="space-y-1.5">
          <div class="flex items-center justify-between">
            <label class="font-semibold text-slate-700">Kata Sandi (Password)</label>
            <span class="text-[10px] text-slate-400 font-mono">Autentikasi Terverifikasi</span>
          </div>
          <div class="relative">
            <Lock class="w-4 h-4 text-slate-400 absolute left-3 top-1/2 -translate-y-1/2" />
            <input
              :type="showPassword ? 'text' : 'password'"
              v-model="password"
              required
              autocomplete="current-password"
              placeholder="Masukkan kata sandi..."
              class="w-full pl-9 pr-10 py-2.5 rounded-xl border border-slate-300 bg-slate-50/50 text-slate-900 placeholder-slate-400 focus:bg-white focus:outline-none focus:ring-2 focus:ring-brand-500 focus:border-brand-500 transition-all text-xs shadow-xs"
            />
            <button
              type="button"
              @click="showPassword = !showPassword"
              class="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600 transition-colors cursor-pointer"
            >
              <EyeOff v-if="showPassword" class="w-4 h-4" />
              <Eye v-else class="w-4 h-4" />
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
import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '../stores/authStore'
import {
  Mail,
  Lock,
  Eye,
  EyeOff,
  ArrowRight,
  ShieldCheck,
  AlertCircle,
  X,
  Loader2,
  Info,
  Clock,
  ShieldAlert
} from 'lucide-vue-next'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const email = ref('')
const password = ref('')
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

async function handleManualLogin() {
  errorMessage.value = ''
  isSubmitting.value = true
  try {
    await authStore.login(email.value.trim(), password.value)
    await proceedRedirect(authStore.currentRole)
  } catch (err: any) {
    errorMessage.value = err?.response?.data?.error || err.message || 'Email atau kata sandi tidak valid'
  } finally {
    isSubmitting.value = false
  }
}
</script>
