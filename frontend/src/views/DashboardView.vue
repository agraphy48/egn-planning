<template>
  <div class="space-y-5 sm:space-y-6 max-w-7xl mx-auto select-none">
    <!-- Welcome & Timeframe Header -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 sm:gap-4 p-4 sm:p-5 rounded-2xl bg-white/95 backdrop-blur-xl border border-slate-200/90 shadow-lg shadow-slate-200/30 relative overflow-hidden">
      <div class="absolute top-0 left-8 right-8 h-0.5 bg-gradient-to-r from-transparent via-brand-500 to-transparent"></div>
      
      <div class="space-y-1">
        <div class="flex items-center gap-2 flex-wrap">
          <h1 class="text-xl sm:text-2xl font-extrabold font-display text-slate-900 tracking-tight">Executive Dashboard</h1>
          <div class="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full bg-emerald-50 border border-emerald-200 text-emerald-700 text-[10px] font-mono font-semibold">
            <span class="relative flex h-2 w-2">
              <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
              <span class="relative inline-flex rounded-full h-2 w-2 bg-emerald-500"></span>
            </span>
            <span>Edge Synced</span>
          </div>
        </div>
        <p class="text-xs text-slate-500 font-medium">Ringkasan operasional milestone engineering, deliverable tim, dan status realisasi proyek.</p>
      </div>

      <div class="flex items-center flex-wrap gap-2 text-xs">
        <span class="px-3 py-1.5 rounded-xl bg-slate-50 border border-slate-200 font-mono text-slate-700 font-bold tabular-nums shadow-2xs">
          {{ store.tasks.length }} Tasks • 4 Proyek Aktif
        </span>
        <button
          @click="store.showToast('Mengekspor ringkasan portfolio ke format laporan eksekutif...')"
          class="px-4 py-1.5 rounded-xl bg-gradient-to-r from-brand-600 to-indigo-600 hover:from-brand-500 hover:to-indigo-500 text-white font-bold transition-all shadow-md shadow-brand-600/20 active:scale-95 text-xs flex items-center gap-1.5"
        >
          <Download class="w-3.5 h-3.5" />
          <span>Ekspor Laporan</span>
        </button>
      </div>
    </div>

    <!-- 4 KPI Summary Cards -->
    <KPIStats />

    <!-- TWO-COLUMN SECTION: PROJECT PROGRESS & RECENT ACTIVITY FEED -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-5 sm:gap-6">
      <!-- Project Progress List (2 cols on desktop) -->
      <div class="lg:col-span-2 p-4 sm:p-6 rounded-2xl bg-white/95 backdrop-blur-xl border border-slate-200/90 shadow-xl shadow-slate-200/50 space-y-4 relative overflow-hidden">
        <div class="absolute top-0 left-6 right-6 h-0.5 bg-gradient-to-r from-transparent via-brand-500/40 to-transparent"></div>
        
        <div class="flex items-center justify-between flex-wrap gap-2">
          <div>
            <h2 class="text-xs font-bold uppercase tracking-wider text-slate-400 font-mono">Portfolio Proyek Aktif</h2>
            <p class="text-xs text-slate-500 font-medium">Tingkat kemajuan pekerjaan per fase dan deliverable</p>
          </div>
          <RouterLink to="/projects" class="text-xs font-bold text-brand-600 hover:text-brand-700 hover:underline flex items-center gap-1">
            <span>Lihat Semua {{ store.projects.length }} Proyek</span>
            <span>→</span>
          </RouterLink>
        </div>

        <!-- Project Items Rows -->
        <div class="space-y-3 pt-1">
          <div
            v-for="proj in store.projects"
            :key="proj.id"
            @click="selectProject(proj.id)"
            class="p-3.5 sm:p-4 rounded-xl border border-slate-200/80 bg-slate-50/60 hover:bg-white hover:border-brand-300 hover:shadow-md cursor-pointer transition-all duration-200 group"
          >
            <!-- Row 1: Name, Badges & Percentage -->
            <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-1.5 sm:gap-2 mb-2.5">
              <div class="flex items-center flex-wrap gap-2 min-w-0">
                <span class="w-2.5 h-2.5 rounded-full flex-shrink-0" :style="{ backgroundColor: proj.color_dot }"></span>
                <span class="text-xs sm:text-sm font-bold text-slate-900 group-hover:text-brand-600 transition-colors">{{ proj.name }}</span>
                <span
                  v-if="proj.client_name"
                  class="px-2 py-0.5 rounded-full text-[10px] font-mono bg-indigo-50 border border-indigo-200 text-indigo-700 font-semibold"
                >
                  Klien: {{ proj.client_name }}
                </span>
                <span
                  v-else
                  class="px-2 py-0.5 rounded-full text-[10px] font-mono bg-slate-200/80 text-slate-700 font-semibold"
                >
                  Internal
                </span>
              </div>
              <span class="text-xs sm:text-sm font-mono font-bold text-slate-900 tabular-nums self-end sm:self-auto">
                {{ proj.progress_percentage }}%
              </span>
            </div>

            <!-- Progress Bar -->
            <div class="w-full bg-slate-200/80 h-2 rounded-full overflow-hidden shadow-inner">
              <div
                class="h-full rounded-full transition-all duration-500 shadow-xs"
                :style="{ width: proj.progress_percentage + '%', backgroundColor: proj.color_dot }"
              ></div>
            </div>

            <!-- Row 3: Metadata Footer -->
            <div class="mt-2.5 flex flex-col sm:flex-row sm:items-center justify-between text-[11px] text-slate-500 font-mono gap-1">
              <span class="tabular-nums font-medium">{{ proj.completed_tasks }} dari {{ proj.total_tasks }} deliverable selesai</span>
              <span class="text-slate-400">Tenggat Waktu: {{ proj.deadline }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Activity Feed & Upcoming Deadlines (1 col on desktop) -->
      <div class="p-4 sm:p-6 rounded-2xl bg-white/95 backdrop-blur-xl border border-slate-200/90 shadow-xl shadow-slate-200/50 space-y-4 relative overflow-hidden">
        <div class="absolute top-0 left-6 right-6 h-0.5 bg-gradient-to-r from-transparent via-brand-500/40 to-transparent"></div>
        
        <div class="flex items-center justify-between">
          <h2 class="text-xs font-bold uppercase tracking-wider text-slate-400 font-mono">Aktivitas Terbaru</h2>
          <span class="text-[10px] font-mono text-emerald-600 font-semibold flex items-center gap-1.5 px-2 py-0.5 rounded-full bg-emerald-50 border border-emerald-200">
            <span class="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse"></span> Live
          </span>
        </div>

        <!-- Activity Stream -->
        <div class="space-y-3 text-xs max-h-[220px] overflow-y-auto pr-1">
          <div
            v-for="act in store.activityFeed"
            :key="act.id"
            class="flex items-start space-x-2.5 p-2 rounded-xl hover:bg-slate-50 transition-colors"
          >
            <img class="w-6 h-6 rounded-lg object-cover mt-0.5 ring-1 ring-slate-200 flex-shrink-0" :src="typeof act.user === 'object' && act.user?.avatar_url ? act.user.avatar_url : (act.avatar_url || 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=100&auto=format&fit=crop&q=80')" :alt="typeof act.user === 'object' && act.user?.name ? act.user.name : String(act.user)" />
            <div class="space-y-0.5 min-w-0 flex-1">
              <p class="text-slate-700 text-xs leading-snug">
                <span class="font-bold text-slate-900">{{ typeof act.user === 'object' && act.user?.name ? act.user.name : act.user }}</span>
                {{ ' ' }} {{ act.action }}
                <span class="font-semibold text-brand-600">{{ act.target }}</span>
              </p>
              <span class="text-[10px] text-slate-400 font-mono">{{ act.time_ago || act.timestamp || 'Baru saja' }}</span>
            </div>
          </div>
        </div>

        <!-- Upcoming Deadlines Micro Widget -->
        <div class="pt-4 border-t border-slate-200 space-y-2.5">
          <span class="text-xs font-bold text-slate-900 flex items-center gap-1.5 font-display">
            <Clock class="w-3.5 h-3.5 text-amber-500" />
            Milestone Kritis & Deliverable
          </span>
          <div class="space-y-1.5 text-[11px] font-mono max-h-48 overflow-y-auto pr-1">
            <div
              v-for="ms in store.milestones"
              :key="ms.id"
              class="flex items-center justify-between p-2.5 rounded-xl bg-slate-50/80 border border-slate-200/70 hover:bg-white hover:shadow-xs transition-all text-xs"
            >
              <span class="truncate flex-1 min-w-0 font-medium text-slate-700 mr-2">{{ ms.title }}</span>
              <span class="font-bold flex-shrink-0 px-2 py-0.5 rounded-md bg-amber-50 border border-amber-200 text-amber-700 text-[10px]">{{ ms.due_date }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- SECTION 3: TEAM WORKLOAD & PROJECT FINANCIALS / BUDGET -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-5 sm:gap-6 pt-1">
      <!-- Team Workload Card -->
      <div class="p-4 sm:p-6 rounded-2xl bg-white/95 backdrop-blur-xl border border-slate-200/90 shadow-xl shadow-slate-200/50 space-y-4 relative overflow-hidden">
        <div class="absolute top-0 left-6 right-6 h-0.5 bg-gradient-to-r from-transparent via-brand-500/40 to-transparent"></div>
        
        <div class="flex items-center justify-between flex-wrap gap-2">
          <div>
            <h2 class="text-xs font-bold uppercase tracking-wider text-slate-400 font-mono">Beban Kerja & Timesheet Tim</h2>
            <p class="text-xs text-slate-500 font-medium">Total jam kerja dan beban tugas aktif per personil engineering</p>
          </div>
          <span class="px-2.5 py-1 rounded-full text-[10px] font-mono bg-brand-50 border border-brand-200 text-brand-700 font-bold">
            Workload
          </span>
        </div>

        <div class="space-y-2.5">
          <div
            v-for="member in teamWorkload"
            :key="member.user_id"
            class="p-3 rounded-xl border border-slate-200/80 bg-slate-50/60 hover:bg-white hover:shadow-xs transition-all flex flex-col sm:flex-row sm:items-center justify-between gap-2 text-xs"
          >
            <div class="flex items-center space-x-2.5 min-w-0">
              <img :src="member.user_avatar || 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=100&auto=format&fit=crop&q=80'" :alt="member.user_name" class="w-8 h-8 rounded-lg object-cover ring-1 ring-slate-200 flex-shrink-0" />
              <div class="min-w-0">
                <div class="font-bold text-slate-900 truncate">{{ member.user_name }}</div>
                <div class="text-[10px] text-slate-500 font-mono">{{ member.task_count ?? member.tasks_count ?? 0 }} tugas aktif</div>
              </div>
            </div>
            <div class="text-left sm:text-right self-end sm:self-auto">
              <span class="text-xs font-mono font-bold text-brand-600">{{ member.total_hours }} jam</span>
              <div class="text-[10px] text-slate-400 font-mono">logged time</div>
            </div>
          </div>
        </div>
      </div>

      <!-- Financials & Project Budget Card -->
      <div class="p-4 sm:p-6 rounded-2xl bg-white/95 backdrop-blur-xl border border-slate-200/90 shadow-xl shadow-slate-200/50 space-y-4 relative overflow-hidden">
        <div class="absolute top-0 left-6 right-6 h-0.5 bg-gradient-to-r from-transparent via-brand-500/40 to-transparent"></div>
        
        <div class="flex items-center justify-between flex-wrap gap-2">
          <div>
            <h2 class="text-xs font-bold uppercase tracking-wider text-slate-400 font-mono">Realisasi Anggaran Proyek (Budget)</h2>
            <p class="text-xs text-slate-500 font-medium">Realisasi pengeluaran vs pagu alokasi per pos anggaran</p>
          </div>
          <span class="px-2.5 py-1 rounded-full text-[10px] font-mono bg-emerald-50 border border-emerald-200 text-emerald-700 font-bold">
            IDR
          </span>
        </div>

        <div class="space-y-3.5">
          <div
            v-for="b in budgetList"
            :key="b.id"
            class="p-3 rounded-xl border border-slate-200/80 bg-slate-50/60 hover:bg-white hover:shadow-xs transition-all space-y-2"
          >
            <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-1 text-xs font-semibold">
              <span class="text-slate-800">{{ b.category }}</span>
              <span class="font-mono text-slate-600 text-[11px]">
                Rp {{ (b.spent / 1000000).toFixed(1) }}M / Rp {{ (b.allocated / 1000000).toFixed(1) }}M
              </span>
            </div>
            
            <div class="w-full bg-slate-200/80 h-2 rounded-full overflow-hidden shadow-inner">
              <div
                class="h-full rounded-full transition-all duration-500 shadow-xs"
                :class="b.spent / b.allocated > 0.85 ? 'bg-gradient-to-r from-amber-500 to-rose-500' : 'bg-gradient-to-r from-brand-600 to-indigo-600'"
                :style="{ width: Math.min(100, Math.round((b.spent / b.allocated) * 100)) + '%' }"
              ></div>
            </div>
            
            <div class="flex justify-between text-[10px] text-slate-500 font-mono font-medium">
              <span>Terpakai: {{ Math.round((b.spent / b.allocated) * 100) }}%</span>
              <span>Sisa: Rp {{ ((b.allocated - b.spent) / 1000000).toFixed(1) }}M</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useProjectStore } from '../stores/projectStore'
import { useRouter } from 'vue-router'
import { projectApi } from '../services/api'
import type { TeamWorkload, BudgetItem } from '../types/project'
import KPIStats from '../components/KPIStats.vue'
import { Clock, Download } from 'lucide-vue-next'

const store = useProjectStore()
const router = useRouter()

const teamWorkload = ref<TeamWorkload[]>([
  {
    user_id: 'usr_yudistira',
    user_name: 'Yudistira Sira Permana',
    user_avatar: 'https://images.unsplash.com/photo-1519085360753-af0119f7cbe7?w=100&auto=format&fit=crop&q=80',
    total_hours: 38.5,
    task_count: 5
  },
  {
    user_id: 'usr_hardianto',
    user_name: 'Hardianto',
    user_avatar: 'https://images.unsplash.com/photo-1506794778202-cad84cf45f1d?w=100&auto=format&fit=crop&q=80',
    total_hours: 42.0,
    task_count: 6
  },
  {
    user_id: 'usr_dedi',
    user_name: 'Dedi Irawan',
    user_avatar: 'https://images.unsplash.com/photo-1522075469751-3a6694fb2f61?w=100&auto=format&fit=crop&q=80',
    total_hours: 34.0,
    task_count: 4
  },
  {
    user_id: 'usr_agus',
    user_name: 'Agus Setiawan',
    user_avatar: 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=100&auto=format&fit=crop&q=80',
    total_hours: 28.0,
    task_count: 3
  }
])

const budgetList = ref<BudgetItem[]>([
  {
    id: 'bgt_01',
    project_id: 'proj_sustagri',
    category: 'Engineering & Core Backend API',
    allocated: 50000000,
    spent: 32500000,
    currency: 'IDR'
  },
  {
    id: 'bgt_02',
    project_id: 'proj_sustagri',
    category: 'Cloud & TLS 1.3 Infrastructure',
    allocated: 15000000,
    spent: 4200000,
    currency: 'IDR'
  },
  {
    id: 'bgt_03',
    project_id: 'proj_sustagri',
    category: 'UI/UX Design & Frontend Micro-interactions',
    allocated: 10000000,
    spent: 9000000,
    currency: 'IDR'
  }
])

function selectProject(projectId: string) {
  store.selectProject(projectId)
  router.push(`/projects/${projectId}`)
}

onMounted(async () => {
  try {
    const wl = await projectApi.getTeamWorkload()
    if (wl && wl.length > 0) teamWorkload.value = wl
  } catch (_) {}

  try {
    const bg = await projectApi.getBudgets(store.currentProjectId || 'proj_sustagri')
    if (bg && bg.budgets && bg.budgets.length > 0) budgetList.value = bg.budgets
  } catch (_) {}
})
</script>
