<template>
  <div class="space-y-5 max-w-5xl mx-auto select-none pb-8">
    <!-- Client Portal Hero Banner -->
    <div class="p-5 sm:p-6 rounded-2xl bg-gradient-to-br from-indigo-950 via-zinc-900 to-zinc-950 text-white shadow-xl relative overflow-hidden border border-indigo-900/40">
      <div class="absolute -right-8 -bottom-8 w-60 h-60 bg-indigo-500/15 rounded-full blur-3xl pointer-events-none"></div>
      <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 relative z-10">
        <div class="space-y-1">
          <div class="flex items-center space-x-2">
            <span class="px-2 py-0.5 rounded-full text-[10px] font-mono uppercase tracking-wider bg-white/15 text-indigo-200 font-semibold border border-white/10">
              Verified Client Portal
            </span>
            <span class="text-xs text-indigo-300 font-mono">v2.4 Delivery Gate</span>
          </div>
          <h1 class="text-xl sm:text-2xl font-bold font-display tracking-tight text-white">
            {{ store.currentProject?.name }}
          </h1>
          <p class="text-xs text-zinc-300 max-w-xl">
            Stakeholder overview, formal contract milestone approvals, and shared deliverable verification.
          </p>
        </div>

        <div class="flex items-center space-x-2 flex-shrink-0">
          <button
            @click="copyPublicLink"
            class="px-3 py-1.5 rounded-lg bg-white/10 hover:bg-white/20 text-white text-xs font-semibold backdrop-blur-sm border border-white/15 transition-all flex items-center gap-1.5"
          >
            <Copy class="w-3.5 h-3.5" />
            <span>Copy Portal Link</span>
          </button>

          <!-- Milestone Signoff Button (Client-only v-can) -->
          <button
            v-can="'milestone:approve'"
            @click="approveCurrentMilestone"
            :disabled="isMilestoneApproved"
            class="px-4 py-1.5 rounded-lg bg-white text-zinc-950 hover:bg-zinc-100 disabled:opacity-60 text-xs font-bold transition-all shadow-md flex items-center gap-1.5"
          >
            <CheckCircle2 class="w-4 h-4 text-emerald-600" />
            <span>{{ isMilestoneApproved ? 'Milestone Approved ✓' : 'Approve Milestone' }}</span>
          </button>
        </div>
      </div>
    </div>

    <!-- Client Deliverables Cards Grid -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
      <div class="p-4 rounded-xl bg-white dark:bg-[#18181B] border border-[#E4E4E7] dark:border-[#27272A] shadow-subtle space-y-2">
        <div class="flex items-center justify-between text-xs text-zinc-400 font-mono">
          <span>Deliverable #1</span>
          <span class="px-2 py-0.2 rounded text-[10px] bg-emerald-50 dark:bg-emerald-950/60 text-emerald-600 dark:text-emerald-400 font-medium">
            Approved
          </span>
        </div>
        <h3 class="text-sm font-semibold text-zinc-900 dark:text-white">Design Token & Component Spec</h3>
        <p class="text-xs text-zinc-500 dark:text-zinc-400">Complete responsive component specs and typography tokens delivered.</p>
      </div>

      <div class="p-4 rounded-xl bg-white dark:bg-[#18181B] border border-[#E4E4E7] dark:border-[#27272A] shadow-subtle space-y-2">
        <div class="flex items-center justify-between text-xs text-zinc-400 font-mono">
          <span>Deliverable #2</span>
          <span
            class="px-2 py-0.2 rounded text-[10px]"
            :class="isMilestoneApproved ? 'bg-emerald-50 dark:bg-emerald-950/60 text-emerald-600 dark:text-emerald-400' : 'bg-amber-50 dark:bg-amber-950/60 text-amber-600 dark:text-amber-400'"
          >
            {{ isMilestoneApproved ? 'Approved by Client' : 'Under Review' }}
          </span>
        </div>
        <h3 class="text-sm font-semibold text-zinc-900 dark:text-white">Staging TestFlight Sandbox Build</h3>
        <p class="text-xs text-zinc-500 dark:text-zinc-400">iOS sandbox with Apple Pay and Stripe sandbox keys active for QA.</p>
        <div v-if="!isMilestoneApproved && canApprove" class="pt-2 flex items-center space-x-2">
          <button
            @click="approveCurrentMilestone"
            class="px-2.5 py-1 rounded bg-brand-600 hover:bg-brand-500 text-white text-xs font-medium"
          >
            Signoff Now
          </button>
          <button
            @click="store.showToast('Revision request recorded. PM notified.')"
            class="px-2.5 py-1 rounded border border-[#E4E4E7] dark:border-[#27272A] text-zinc-600 dark:text-zinc-400 hover:text-zinc-900 text-xs"
          >
            Request Revision
          </button>
        </div>
      </div>

      <div class="p-4 rounded-xl bg-white dark:bg-[#18181B] border border-[#E4E4E7] dark:border-[#27272A] shadow-subtle space-y-2">
        <div class="flex items-center justify-between text-xs text-zinc-400 font-mono">
          <span>Deliverable #3</span>
          <span class="px-2 py-0.2 rounded text-[10px] bg-zinc-100 dark:bg-zinc-800 text-zinc-400">
            Queued
          </span>
        </div>
        <h3 class="text-sm font-semibold text-zinc-900 dark:text-white">Security & Cloudflare Audit</h3>
        <p class="text-xs text-zinc-500 dark:text-zinc-400">Scheduled for final verification signoff on Dec 04, 2026.</p>
      </div>
    </div>

    <!-- Client-Facing Milestones & Formal Signoffs Table -->
    <div class="p-5 rounded-xl bg-white dark:bg-[#18181B] border border-[#E4E4E7] dark:border-[#27272A] shadow-subtle space-y-4">
      <div class="flex items-center justify-between">
        <h2 class="text-xs font-bold uppercase tracking-wider text-zinc-400 font-mono">
          Contract Milestones & Governance
        </h2>
        <span class="text-xs font-mono text-zinc-400">{{ store.milestones.length }} Milestones</span>
      </div>

      <div class="divide-y divide-[#E4E4E7] dark:divide-[#27272A] text-xs">
        <div
          v-for="ms in store.milestones"
          :key="ms.id"
          class="py-3 flex flex-col sm:flex-row sm:items-center justify-between gap-2"
        >
          <div class="space-y-0.5">
            <div class="flex items-center space-x-2">
              <span class="font-semibold text-zinc-900 dark:text-white">{{ ms.title }}</span>
              <span
                class="px-2 py-0.2 rounded text-[10px] font-mono"
                :class="
                  ms.status === 'approved'
                    ? 'bg-emerald-50 text-emerald-600 dark:bg-emerald-950/60 dark:text-emerald-400'
                    : ms.status === 'in_review'
                    ? 'bg-amber-50 text-amber-600 dark:bg-amber-950/60 dark:text-amber-400'
                    : 'bg-zinc-100 text-zinc-500 dark:bg-zinc-800 dark:text-zinc-400'
                "
              >
                {{ ms.status.toUpperCase() }}
              </span>
            </div>
            <p class="text-zinc-500 dark:text-zinc-400 text-xs">{{ ms.description }}</p>
          </div>

          <div class="flex items-center space-x-3 text-right">
            <div class="text-[11px] font-mono text-zinc-400 tabular-nums">Due: {{ ms.due_date }}</div>
            <div v-if="ms.status !== 'approved'">
              <button
                v-can="'milestone:approve'"
                @click="store.approveMilestone(ms.id)"
                class="px-3 py-1 rounded bg-brand-600 hover:bg-brand-500 text-white font-semibold text-xs shadow-2xs transition-colors"
              >
                Signoff
              </button>
              <span
                v-if="!canApprove"
                class="text-[11px] font-mono text-amber-600 dark:text-amber-400 font-medium"
              >
                Awaiting Signoff
              </span>
            </div>
            <span v-else class="text-[11px] font-mono text-emerald-600 dark:text-emerald-400 font-semibold flex items-center gap-1">
              <Check class="w-3.5 h-3.5" /> Signed
            </span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useProjectStore } from '../stores/projectStore'
import { useAuthStore } from '../stores/authStore'
import { Copy, CheckCircle2, Check } from 'lucide-vue-next'

const store = useProjectStore()
const authStore = useAuthStore()

const canApprove = computed(() => authStore.can('milestone:approve'))

const isMilestoneApproved = computed(() => {
  const m2 = store.milestones.find((m) => m.id === 'ms_02')
  return m2?.status === 'approved'
})

function approveCurrentMilestone() {
  store.approveMilestone('ms_02')
}

function copyPublicLink() {
  navigator.clipboard?.writeText(window.location.origin + '/client-portal?guest_token=pulse_client_preview')
  store.showToast('Public Client Link copied to clipboard!')
}
</script>
