<template>
  <div class="min-h-[70vh] flex items-center justify-center p-6">
    <div class="max-w-md w-full text-center space-y-6 p-8 rounded-2xl bg-white dark:bg-[#131B2E] border border-slate-200 dark:border-slate-800 shadow-xl">
      <div class="w-14 h-14 mx-auto rounded-2xl bg-rose-50 dark:bg-rose-950/50 border border-rose-200 dark:border-rose-900/60 flex items-center justify-center text-rose-600 dark:text-rose-400">
        <ShieldAlert class="w-7 h-7" />
      </div>

      <div class="space-y-2">
        <span class="text-xs font-mono font-bold text-rose-500 uppercase tracking-widest">403 Access Denied</span>
        <h1 class="text-2xl font-bold font-display text-slate-900 dark:text-white">Insufficient Role Privileges</h1>
        <p class="text-xs text-slate-500 dark:text-slate-400 leading-relaxed">
          Your current account role (<span class="font-mono font-semibold text-brand-600 dark:text-brand-400 uppercase">{{ authStore.currentRole }}</span>) does not have the authorization policy required to view this resource.
        </p>
      </div>

      <div class="p-3.5 rounded-xl bg-slate-50 dark:bg-slate-900/60 border border-slate-100 dark:border-slate-800 text-left text-xs font-mono space-y-1">
        <div class="text-slate-400 text-[10px] uppercase font-bold">Policy Evaluation</div>
        <div class="text-slate-700 dark:text-slate-300">
          User: <span class="text-slate-900 dark:text-white font-semibold">{{ authStore.user?.email }}</span>
        </div>
        <div class="text-slate-700 dark:text-slate-300">
          Global Role: <span class="text-brand-600 dark:text-brand-400 font-semibold">{{ authStore.currentRole }}</span>
        </div>
      </div>

      <div class="flex items-center justify-center space-x-3 pt-2">
        <RouterLink
          to="/dashboard"
          class="px-4 py-2 rounded-lg bg-slate-900 dark:bg-white text-white dark:text-slate-900 text-xs font-semibold hover:opacity-90 transition-opacity"
        >
          Return to Dashboard
        </RouterLink>

        <button
          v-if="!isProd && authStore.allowQuickLogin"
          @click="authStore.switchUser('agus')"
          class="px-4 py-2 rounded-lg border border-slate-200 dark:border-slate-700 text-slate-700 dark:text-slate-300 text-xs font-semibold hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors"
        >
          Switch to PM (Dev)
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useAuthStore } from '../stores/authStore'
import { ShieldAlert } from 'lucide-vue-next'

const authStore = useAuthStore()
const isProd = computed(() => import.meta.env.PROD)
</script>
