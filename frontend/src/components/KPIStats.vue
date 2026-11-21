<template>
  <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-3.5 sm:gap-4">
    <!-- KPI 1: Active Projects -->
    <div class="p-4 sm:p-5 rounded-2xl bg-white/95 backdrop-blur-xl border border-slate-200/90 shadow-lg shadow-slate-200/40 hover:shadow-xl hover:border-brand-300 transition-all duration-300 relative overflow-hidden group cursor-pointer">
      <div class="absolute top-0 left-6 right-6 h-0.5 bg-gradient-to-r from-transparent via-brand-500/40 to-transparent group-hover:via-brand-500 transition-all duration-500"></div>
      <div class="flex items-center justify-between text-xs text-slate-500 mb-2">
        <span class="font-semibold text-slate-600">Total Proyek Aktif</span>
        <span class="p-1.5 rounded-xl bg-brand-50 border border-brand-200 text-brand-700 shadow-2xs">
          <FolderKanban class="w-4 h-4" />
        </span>
      </div>
      <div class="flex items-baseline space-x-2">
        <span class="text-2xl sm:text-3xl font-extrabold font-display text-slate-900 tabular-nums tracking-tight">
          {{ store.kpiStats?.active_projects ?? store.projects.length }}
        </span>
        <span class="text-xs text-emerald-700 font-semibold flex items-center tabular-nums px-2 py-0.5 rounded-full bg-emerald-50 border border-emerald-200">
          <TrendingUp class="w-3 h-3 mr-1" />
          {{ store.kpiStats?.active_projects_delta || '4 Proyek Aktif' }}
        </span>
      </div>
      <div class="mt-3 flex items-center text-[11px] text-slate-400 font-mono">
        <span>SUSTAGRI (F1/F2)</span>
        <span class="mx-1.5 text-slate-300">•</span>
        <span>SMART (F1/F2)</span>
      </div>
    </div>

    <!-- KPI 2: Tasks Due This Week -->
    <div class="p-4 sm:p-5 rounded-2xl bg-white/95 backdrop-blur-xl border border-slate-200/90 shadow-lg shadow-slate-200/40 hover:shadow-xl hover:border-amber-300 transition-all duration-300 relative overflow-hidden group cursor-pointer">
      <div class="absolute top-0 left-6 right-6 h-0.5 bg-gradient-to-r from-transparent via-amber-500/40 to-transparent group-hover:via-amber-500 transition-all duration-500"></div>
      <div class="flex items-center justify-between text-xs text-slate-500 mb-2">
        <span class="font-semibold text-slate-600">Tenggat Sprint Ini</span>
        <span class="p-1.5 rounded-xl bg-amber-50 border border-amber-200 text-amber-700 shadow-2xs">
          <CalendarClock class="w-4 h-4" />
        </span>
      </div>
      <div class="flex items-baseline space-x-2">
        <span class="text-2xl sm:text-3xl font-extrabold font-display text-slate-900 tabular-nums tracking-tight">
          {{ store.kpiStats?.tasks_due_this_week ?? 0 }}
        </span>
        <span class="text-xs text-amber-700 font-semibold font-mono px-2 py-0.5 rounded-full bg-amber-50 border border-amber-200">
          Sprint berjalan
        </span>
      </div>
      <div class="mt-3 flex items-center text-[11px] text-amber-700 font-mono">
        <AlertCircle class="w-3 h-3 mr-1 flex-shrink-0" />
        <span>Total 49 Deliverable Tasks</span>
      </div>
    </div>

    <!-- KPI 3: Overdue Tasks -->
    <div class="p-4 sm:p-5 rounded-2xl bg-white/95 backdrop-blur-xl border border-slate-200/90 shadow-lg shadow-slate-200/40 hover:shadow-xl hover:border-rose-300 transition-all duration-300 relative overflow-hidden group cursor-pointer">
      <div class="absolute top-0 left-6 right-6 h-0.5 bg-gradient-to-r from-transparent via-rose-500/40 to-transparent group-hover:via-rose-500 transition-all duration-500"></div>
      <div class="flex items-center justify-between text-xs text-slate-500 mb-2">
        <span class="font-semibold text-slate-600">Tugas Terlambat</span>
        <span class="p-1.5 rounded-xl bg-rose-50 border border-rose-200 text-rose-700 shadow-2xs">
          <Flame class="w-4 h-4" />
        </span>
      </div>
      <div class="flex items-baseline space-x-2">
        <span class="text-2xl sm:text-3xl font-extrabold font-display text-rose-600 tabular-nums tracking-tight">
          {{ store.kpiStats?.overdue_tasks ?? 0 }}
        </span>
        <span class="text-xs text-rose-700 font-semibold font-mono px-2 py-0.5 rounded-full bg-rose-50 border border-rose-200">
          {{ store.kpiStats?.overdue_tasks_delta || 'Perlu akselerasi' }}
        </span>
      </div>
      <div class="mt-3 flex items-center text-[11px] text-slate-400 font-mono">
        <span>Evaluasi tenggat sprint</span>
      </div>
    </div>

    <!-- KPI 4: Overall Velocity & Completion Rate -->
    <div class="p-4 sm:p-5 rounded-2xl bg-white/95 backdrop-blur-xl border border-slate-200/90 shadow-lg shadow-slate-200/40 hover:shadow-xl hover:border-emerald-300 transition-all duration-300 relative overflow-hidden group cursor-pointer">
      <div class="absolute top-0 left-6 right-6 h-0.5 bg-gradient-to-r from-transparent via-emerald-500/40 to-transparent group-hover:via-emerald-500 transition-all duration-500"></div>
      <div class="flex items-center justify-between text-xs text-slate-500 mb-2">
        <span class="font-semibold text-slate-600">Kecepatan & Selesai</span>
        <span class="p-1.5 rounded-xl bg-emerald-50 border border-emerald-200 text-emerald-700 shadow-2xs">
          <Gauge class="w-4 h-4" />
        </span>
      </div>
      <div class="flex items-baseline space-x-2">
        <span class="text-2xl sm:text-3xl font-extrabold font-display text-slate-900 tabular-nums tracking-tight">
          {{ store.kpiStats?.completion_rate ?? 0 }}%
        </span>
        <span class="text-xs text-emerald-700 font-semibold font-mono px-2 py-0.5 rounded-full bg-emerald-50 border border-emerald-200 flex items-center tabular-nums">
          <TrendingUp class="w-3 h-3 mr-1" />
          {{ store.kpiStats?.completion_rate_delta || 'Tingkat Penyelesaian' }}
        </span>
      </div>
      <!-- Progress Bar with subtle shadow -->
      <div class="mt-3 w-full bg-slate-200/80 h-2 rounded-full overflow-hidden shadow-inner">
        <div
          class="bg-gradient-to-r from-emerald-500 to-teal-500 h-full rounded-full transition-all duration-500 shadow-xs"
          :style="{ width: (store.kpiStats?.completion_rate ?? 0) + '%' }"
        ></div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useProjectStore } from '../stores/projectStore'
import {
  FolderKanban,
  TrendingUp,
  CalendarClock,
  AlertCircle,
  Flame,
  Gauge
} from 'lucide-vue-next'

const store = useProjectStore()
</script>
