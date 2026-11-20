<template>
  <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-3.5">
    <!-- KPI 1: Active Projects -->
    <div class="p-4 rounded-xl bg-white dark:bg-[#18181B] border border-[#E4E4E7] dark:border-[#27272A] shadow-subtle hover:shadow-card transition-shadow">
      <div class="flex items-center justify-between text-xs text-zinc-500 dark:text-zinc-400 mb-2">
        <span class="font-medium">Active Projects</span>
        <span class="p-1 rounded bg-brand-50 dark:bg-brand-950/50 text-brand-600 dark:text-brand-400">
          <FolderKanban class="w-4 h-4" />
        </span>
      </div>
      <div class="flex items-baseline space-x-2">
        <span class="text-2xl font-bold font-display text-zinc-900 dark:text-white tabular-nums tracking-tight">
          {{ store.kpiStats?.active_projects || 12 }}
        </span>
        <span class="text-xs text-emerald-600 dark:text-emerald-400 font-medium flex items-center tabular-nums">
          <TrendingUp class="w-3 h-3 mr-0.5" />
          {{ store.kpiStats?.active_projects_delta || '+2 this month' }}
        </span>
      </div>
      <div class="mt-2.5 flex items-center text-[11px] text-zinc-400 font-mono">
        <span>8 Client</span>
        <span class="mx-1.5">•</span>
        <span>4 Internal</span>
      </div>
    </div>

    <!-- KPI 2: Tasks Due This Week -->
    <div class="p-4 rounded-xl bg-white dark:bg-[#18181B] border border-[#E4E4E7] dark:border-[#27272A] shadow-subtle hover:shadow-card transition-shadow">
      <div class="flex items-center justify-between text-xs text-zinc-500 dark:text-zinc-400 mb-2">
        <span class="font-medium">Tasks Due This Week</span>
        <span class="p-1 rounded bg-amber-50 dark:bg-amber-950/50 text-amber-600 dark:text-amber-400">
          <CalendarClock class="w-4 h-4" />
        </span>
      </div>
      <div class="flex items-baseline space-x-2">
        <span class="text-2xl font-bold font-display text-zinc-900 dark:text-white tabular-nums tracking-tight">
          {{ store.kpiStats?.tasks_due_this_week || 28 }}
        </span>
        <span class="text-xs text-zinc-500 font-medium font-mono">Across 4 teams</span>
      </div>
      <div class="mt-2.5 flex items-center text-[11px] text-amber-600 dark:text-amber-400 font-mono">
        <AlertCircle class="w-3 h-3 mr-1" />
        <span>9 due in next 48h</span>
      </div>
    </div>

    <!-- KPI 3: Overdue Tasks -->
    <div class="p-4 rounded-xl bg-white dark:bg-[#18181B] border border-[#E4E4E7] dark:border-[#27272A] shadow-subtle hover:shadow-card transition-shadow">
      <div class="flex items-center justify-between text-xs text-zinc-500 dark:text-zinc-400 mb-2">
        <span class="font-medium">Overdue Tasks</span>
        <span class="p-1 rounded bg-rose-50 dark:bg-rose-950/50 text-rose-600 dark:text-rose-400">
          <Flame class="w-4 h-4" />
        </span>
      </div>
      <div class="flex items-baseline space-x-2">
        <span class="text-2xl font-bold font-display text-rose-600 dark:text-rose-500 tabular-nums tracking-tight">
          {{ store.kpiStats?.overdue_tasks || 3 }}
        </span>
        <span class="text-xs text-rose-500/80 font-medium font-mono tabular-nums">
          {{ store.kpiStats?.overdue_tasks_delta || '-1 vs yesterday' }}
        </span>
      </div>
      <div class="mt-2.5 flex items-center text-[11px] text-zinc-400 font-mono">
        <span>Requires immediate triage</span>
      </div>
    </div>

    <!-- KPI 4: Overall Velocity & Completion Rate -->
    <div class="p-4 rounded-xl bg-white dark:bg-[#18181B] border border-[#E4E4E7] dark:border-[#27272A] shadow-subtle hover:shadow-card transition-shadow">
      <div class="flex items-center justify-between text-xs text-zinc-500 dark:text-zinc-400 mb-2">
        <span class="font-medium">Overall Sprint Velocity</span>
        <span class="p-1 rounded bg-emerald-50 dark:bg-emerald-950/50 text-emerald-600 dark:text-emerald-400">
          <Gauge class="w-4 h-4" />
        </span>
      </div>
      <div class="flex items-baseline space-x-2">
        <span class="text-2xl font-bold font-display text-zinc-900 dark:text-white tabular-nums tracking-tight">
          {{ store.kpiStats?.completion_rate || 88.4 }}%
        </span>
        <span class="text-xs text-emerald-600 dark:text-emerald-400 font-medium flex items-center tabular-nums">
          <TrendingUp class="w-3 h-3 mr-0.5" />
          {{ store.kpiStats?.completion_rate_delta || '+4.2%' }}
        </span>
      </div>
      <!-- Progress Bar -->
      <div class="mt-3 w-full bg-zinc-100 dark:bg-zinc-800 h-1.5 rounded-full overflow-hidden">
        <div
          class="bg-brand-600 h-full rounded-full transition-all duration-500"
          :style="{ width: (store.kpiStats?.completion_rate || 88.4) + '%' }"
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
