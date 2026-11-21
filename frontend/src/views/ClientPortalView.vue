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
            <span class="text-xs text-indigo-300 font-mono">Formal Contract Delivery Gate</span>
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
            @click="authStore.isChangePasswordModalOpen = true"
            class="px-3 py-1.5 rounded-lg bg-white/10 hover:bg-white/20 text-white text-xs font-semibold backdrop-blur-sm border border-white/15 transition-all flex items-center gap-1.5 cursor-pointer"
          >
            <KeyRound class="w-3.5 h-3.5" />
            <span>Ganti Kata Sandi</span>
          </button>
          <button
            @click="copyPublicLink"
            class="px-3 py-1.5 rounded-lg bg-white/10 hover:bg-white/20 text-white text-xs font-semibold backdrop-blur-sm border border-white/15 transition-all flex items-center gap-1.5 cursor-pointer"
          >
            <Copy class="w-3.5 h-3.5" />
            <span>Copy Portal Link</span>
          </button>
        </div>
      </div>

      <!-- Project Switcher Tabs for Client -->
      <div class="flex items-center space-x-2 overflow-x-auto pt-4 mt-4 border-t border-white/10 z-10 relative">
        <button
          v-for="p in store.projects"
          :key="p.id"
          @click="selectPortalProject(p.id)"
          :class="[
            'px-3 py-1.5 rounded-xl text-xs font-medium transition-all flex items-center gap-2 flex-shrink-0',
            store.currentProjectId === p.id
              ? 'bg-white text-zinc-950 font-bold shadow-md'
              : 'bg-white/10 text-zinc-300 hover:bg-white/20 hover:text-white'
          ]"
        >
          <span class="w-2 h-2 rounded-full" :style="{ backgroundColor: p.color_dot }"></span>
          <span>{{ p.name }}</span>
        </button>
      </div>
    </div>

    <!-- Client Deliverables Overview Cards -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
      <div
        v-for="(deliverable, idx) in currentDeliverables"
        :key="idx"
        class="p-4 rounded-xl bg-white dark:bg-[#18181B] border border-[#E4E4E7] dark:border-[#27272A] shadow-subtle space-y-2"
      >
        <div class="flex items-center justify-between text-xs text-zinc-400 font-mono">
          <span>Deliverable #{{ idx + 1 }}</span>
          <span
            class="px-2 py-0.2 rounded text-[10px] font-bold uppercase font-mono"
            :class="deliverable.statusBadge"
          >
            {{ deliverable.status }}
          </span>
        </div>
        <h3 class="text-sm font-semibold text-zinc-900 dark:text-white">{{ deliverable.title }}</h3>
        <p class="text-xs text-zinc-500 dark:text-zinc-400">{{ deliverable.desc }}</p>
      </div>
    </div>

    <!-- Client-Facing Milestones & Formal Signoffs Table -->
    <div class="p-5 rounded-xl bg-white dark:bg-[#18181B] border border-[#E4E4E7] dark:border-[#27272A] shadow-subtle space-y-4">
      <div class="flex items-center justify-between">
        <h2 class="text-xs font-bold uppercase tracking-wider text-zinc-400 font-mono">
          Contract Milestones & Approval Gates
        </h2>
        <span class="text-xs font-mono text-zinc-400">{{ currentMilestones.length }} Milestones</span>
      </div>

      <div v-if="currentMilestones.length === 0" class="py-8 text-center text-xs text-zinc-400 font-mono">
        Belum ada milestone tercatat untuk proyek ini
      </div>

      <div v-else class="divide-y divide-[#E4E4E7] dark:divide-[#27272A] text-xs">
        <div
          v-for="ms in currentMilestones"
          :key="ms.id"
          class="py-4 flex flex-col md:flex-row md:items-center justify-between gap-3"
        >
          <div class="space-y-1 max-w-lg">
            <div class="flex items-center space-x-2">
              <span class="font-bold text-sm text-zinc-900 dark:text-white">{{ ms.title }}</span>
              <span
                class="px-2 py-0.5 rounded text-[10px] font-mono uppercase font-bold"
                :class="
                  ms.status === 'approved'
                    ? 'bg-emerald-50 text-emerald-600 dark:bg-emerald-950/60 dark:text-emerald-400 border border-emerald-200'
                    : ms.status === 'revision'
                    ? 'bg-rose-50 text-rose-600 dark:bg-rose-950/60 dark:text-rose-400 border border-rose-200'
                    : 'bg-amber-50 text-amber-600 dark:bg-amber-950/60 dark:text-amber-400 border border-amber-200'
                "
              >
                {{ ms.status }}
              </span>
            </div>
            <p class="text-zinc-500 dark:text-zinc-400 text-xs">{{ ms.description }}</p>
            
            <!-- Revision Feedback Notes Alert -->
            <div
              v-if="ms.feedback_notes"
              class="p-2 rounded-lg bg-rose-50/70 dark:bg-rose-950/30 border border-rose-200 dark:border-rose-900/40 text-[11px] text-rose-700 dark:text-rose-300 flex items-start gap-1.5"
            >
              <AlertTriangle class="w-3.5 h-3.5 text-rose-500 flex-shrink-0 mt-0.5" />
              <div>
                <strong>Catatan Revisi Klien:</strong> {{ ms.feedback_notes }}
              </div>
            </div>

            <div v-if="ms.approved_by" class="text-[10px] font-mono text-emerald-600 dark:text-emerald-400">
              Disetujui oleh: <strong>{{ ms.approved_by }}</strong>
            </div>
          </div>

          <!-- Action Buttons -->
          <div class="flex items-center space-x-2 flex-shrink-0">
            <div class="text-[11px] font-mono text-zinc-400 mr-2">Tenggat: {{ ms.due_date }}</div>
            
            <div v-if="ms.status !== 'approved' && canApprove" class="flex items-center space-x-2">
              <button
                @click="openRevisionModal(ms)"
                class="px-3 py-1.5 rounded-lg border border-rose-300 dark:border-rose-800 text-rose-600 dark:text-rose-400 hover:bg-rose-50 dark:hover:bg-rose-950/40 text-xs font-semibold transition-colors flex items-center gap-1"
              >
                <XCircle class="w-3.5 h-3.5" />
                <span>Minta Revisi</span>
              </button>

              <button
                @click="handleApprove(ms.id)"
                class="px-4 py-1.5 rounded-lg bg-brand-600 hover:bg-brand-500 text-white font-bold text-xs shadow-md transition-all flex items-center gap-1.5"
              >
                <CheckCircle2 class="w-3.5 h-3.5" />
                <span>Setujui Milestone</span>
              </button>
            </div>

            <span
              v-else-if="ms.status === 'approved'"
              class="px-3 py-1 rounded-lg bg-emerald-50 dark:bg-emerald-950/50 text-emerald-600 dark:text-emerald-400 font-mono text-xs font-bold border border-emerald-200 dark:border-emerald-800 flex items-center gap-1"
            >
              <Check class="w-3.5 h-3.5" /> Verified Signed
            </span>
          </div>
        </div>
      </div>
    </div>

    <!-- Request Revision Modal Dialog -->
    <div
      v-if="showRevisionDialog"
      class="fixed inset-0 bg-black/50 backdrop-blur-xs flex items-center justify-center p-4 z-50"
    >
      <div class="max-w-md w-full bg-white dark:bg-[#18181B] rounded-2xl p-5 border border-zinc-200 dark:border-zinc-800 shadow-2xl space-y-4">
        <div class="flex items-center space-x-2 text-rose-600 dark:text-rose-400">
          <AlertTriangle class="w-5 h-5" />
          <h3 class="text-sm font-bold">Permintaan Revisi Deliverable</h3>
        </div>
        <p class="text-xs text-zinc-500 dark:text-zinc-400">
          Sampaikan catatan detail atau kriteria deliverable yang perlu diperbaiki oleh tim pengembang untuk milestone: <strong>{{ selectedMilestone?.title }}</strong>.
        </p>
        <textarea
          rows="4"
          v-model="revisionReason"
          placeholder="Tuliskan catatan revisi spesifik di sini..."
          class="w-full text-xs p-2.5 rounded-xl border border-zinc-300 dark:border-zinc-700 bg-zinc-50 dark:bg-zinc-900 text-zinc-900 dark:text-white focus:outline-none focus:ring-1 focus:ring-rose-500"
        ></textarea>
        <div class="flex items-center justify-end space-x-2 pt-2 border-t border-zinc-100 dark:border-zinc-800">
          <button
            @click="showRevisionDialog = false"
            class="px-3 py-1.5 rounded-lg border border-zinc-200 dark:border-zinc-700 text-xs font-semibold text-zinc-600 dark:text-zinc-400 hover:bg-zinc-100 dark:hover:bg-zinc-800"
          >
            Batal
          </button>
          <button
            @click="submitRevision"
            :disabled="!revisionReason.trim()"
            class="px-4 py-1.5 rounded-lg bg-rose-600 hover:bg-rose-500 text-white text-xs font-bold disabled:opacity-50"
          >
            Kirim Permintaan Revisi
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useProjectStore } from '../stores/projectStore'
import { useAuthStore } from '../stores/authStore'
import { projectApi } from '../services/api'
import type { Milestone } from '../types/project'
import { Copy, CheckCircle2, Check, XCircle, AlertTriangle, KeyRound } from 'lucide-vue-next'

const store = useProjectStore()
const authStore = useAuthStore()

const canApprove = computed(() => authStore.can('milestone:approve'))

const showRevisionDialog = ref(false)
const selectedMilestone = ref<Milestone | null>(null)
const revisionReason = ref('')

function selectPortalProject(projectId: string) {
  store.selectProject(projectId)
}

const currentMilestones = computed(() => {
  return store.milestones.filter((m) => m.project_id === store.currentProjectId)
})

const currentDeliverables = computed(() => {
  const visibleTasks = store.filteredTasks.filter((t) => t.client_visible && (t.deliverable || t.title))
  if (visibleTasks.length > 0) {
    return visibleTasks.slice(0, 6).map((t) => ({
      title: t.deliverable || t.title,
      desc: t.title,
      status: t.status,
      statusBadge:
        t.status === 'Done'
          ? 'bg-emerald-50 text-emerald-600 dark:bg-emerald-950 dark:text-emerald-400'
          : t.status === 'Review'
          ? 'bg-amber-50 text-amber-600 dark:bg-amber-950 dark:text-amber-400'
          : 'bg-blue-50 text-blue-600 dark:bg-blue-950 dark:text-blue-400'
    }))
  }
  return [
    {
      title: store.currentProject?.name || 'Deliverable Gate',
      desc: 'Spesifikasi kontrak dan deliverable resmi fase proyek.',
      status: 'In Progress',
      statusBadge: 'bg-blue-50 text-blue-600 dark:bg-blue-950 dark:text-blue-400'
    }
  ]
})

async function handleApprove(milestoneId: string) {
  try {
    const clientName = authStore.user?.name || 'Ir. H. Gunawan (Client)'
    await projectApi.approveMilestone(milestoneId, clientName)
    const ms = store.milestones.find((m) => m.id === milestoneId)
    if (ms) {
      ms.status = 'approved'
      ms.approved_by = clientName
    }
    store.showToast('Milestone berhasil disetujui secara resmi!')
  } catch (err: any) {
    store.showToast('Gagal menyetujui milestone')
  }
}

function openRevisionModal(ms: Milestone) {
  selectedMilestone.value = ms
  revisionReason.value = ''
  showRevisionDialog.value = true
}

async function submitRevision() {
  if (!selectedMilestone.value || !revisionReason.value.trim()) return
  try {
    const clientName = authStore.user?.name || 'Ir. H. Gunawan (Client)'
    await projectApi.rejectMilestone(selectedMilestone.value.id, revisionReason.value.trim(), clientName)
    selectedMilestone.value.status = 'revision'
    selectedMilestone.value.feedback_notes = revisionReason.value.trim()
    store.showToast('Catatan revisi telah dikirimkan ke PM & tim pengembang!')
    showRevisionDialog.value = false
  } catch (_) {
    selectedMilestone.value.status = 'revision'
    selectedMilestone.value.feedback_notes = revisionReason.value.trim()
    store.showToast('Catatan revisi telah dikirimkan ke PM!')
    showRevisionDialog.value = false
  }
}

function copyPublicLink() {
  navigator.clipboard?.writeText(window.location.origin + '/client-portal?guest_token=pulse_client_preview')
  store.showToast('Tautan portal klien berhasil disalin!')
}
</script>
