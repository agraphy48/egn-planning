<template>
  <div class="min-h-screen bg-[#F8FAFC] dark:bg-[#0F172A] flex items-center justify-center p-4">
    <div class="max-w-md w-full space-y-6 p-8 rounded-2xl bg-white dark:bg-[#131B2E] border border-slate-200 dark:border-slate-800 shadow-xl">
      <div class="text-center space-y-2">
        <div class="w-10 h-10 rounded-xl bg-brand-600 mx-auto flex items-center justify-center text-white font-display font-bold text-lg shadow-sm">
          P
        </div>
        <h1 class="text-xl font-bold font-display text-slate-900 dark:text-white">Sign in to Pulse</h1>
        <p class="text-xs text-slate-500 dark:text-slate-400">High-Density Project Workspace & Client Portal</p>
      </div>

      <form @submit.prevent="handleLogin" class="space-y-4 text-xs">
        <div>
          <label class="font-semibold text-slate-700 dark:text-slate-300">Email Address</label>
          <input
            type="email"
            v-model="email"
            required
            class="w-full mt-1 px-3 py-2 rounded-lg border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-900 text-slate-900 dark:text-white focus:outline-none focus:ring-1 focus:ring-brand-500"
          />
        </div>

        <div>
          <label class="font-semibold text-slate-700 dark:text-slate-300">Password</label>
          <input
            type="password"
            v-model="password"
            required
            class="w-full mt-1 px-3 py-2 rounded-lg border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-900 text-slate-900 dark:text-white focus:outline-none focus:ring-1 focus:ring-brand-500"
          />
        </div>

        <button
          type="submit"
          class="w-full py-2.5 rounded-lg bg-brand-600 hover:bg-brand-500 text-white font-semibold text-xs shadow-sm transition-colors"
        >
          Sign In
        </button>
      </form>

      <!-- Quick 1-Click Role Logins for Testing -->
      <div class="pt-4 border-t border-slate-100 dark:border-slate-800 space-y-2">
        <div class="text-[10px] font-mono text-slate-400 uppercase font-semibold text-center">
          1-Click Test Role Accounts
        </div>
        <div class="grid grid-cols-2 gap-2 text-[11px]">
          <button
            @click="quickLogin('super_admin')"
            class="p-2 rounded-lg border border-slate-200 dark:border-slate-700 hover:bg-slate-50 dark:hover:bg-slate-800 text-left transition-colors"
          >
            <div class="font-bold text-slate-900 dark:text-white">Super Admin</div>
            <div class="text-[10px] text-slate-400">admin@pulse</div>
          </button>

          <button
            @click="quickLogin('pm')"
            class="p-2 rounded-lg border border-slate-200 dark:border-slate-700 hover:bg-slate-50 dark:hover:bg-slate-800 text-left transition-colors"
          >
            <div class="font-bold text-slate-900 dark:text-white">Project Manager</div>
            <div class="text-[10px] text-slate-400">pm@pulse</div>
          </button>

          <button
            @click="quickLogin('member')"
            class="p-2 rounded-lg border border-slate-200 dark:border-slate-700 hover:bg-slate-50 dark:hover:bg-slate-800 text-left transition-colors"
          >
            <div class="font-bold text-slate-900 dark:text-white">Team Member</div>
            <div class="text-[10px] text-slate-400">member@pulse</div>
          </button>

          <button
            @click="quickLogin('client')"
            class="p-2 rounded-lg border border-slate-200 dark:border-slate-700 hover:bg-slate-50 dark:hover:bg-slate-800 text-left transition-colors"
          >
            <div class="font-bold text-slate-900 dark:text-white">Client Stakeholder</div>
            <div class="text-[10px] text-slate-400">client@pulse</div>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore, SEED_USERS } from '../stores/authStore'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const email = ref('pm@pulse.workspace')
const password = ref('Password123!')

async function handleLogin() {
  await authStore.login(email.value, password.value)
  const redirect = (route.query.redirect as string) || '/dashboard'
  router.push(redirect)
}

function quickLogin(role: 'super_admin' | 'pm' | 'member' | 'client') {
  const seed = SEED_USERS[role]
  email.value = seed.email
  password.value = 'Password123!'
  handleLogin()
}
</script>
