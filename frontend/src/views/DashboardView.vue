<template>
  <div class="space-y-5 max-w-7xl mx-auto select-none">
    <!-- Welcome & Timeframe Header -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3">
      <div>
        <h1 class="text-lg font-bold font-display text-zinc-900 dark:text-white tracking-tight">Executive Dashboard</h1>
        <p class="text-xs text-zinc-500 dark:text-zinc-400">High-level operational overview across engineering milestones and client deliverables.</p>
      </div>
      <div class="flex items-center space-x-2 text-xs">
        <span class="text-zinc-400">Sprint:</span>
        <span class="px-2.5 py-1 rounded-md bg-white dark:bg-[#18181B] border border-[#E4E4E7] dark:border-[#27272A] font-mono text-zinc-700 dark:text-zinc-300 font-medium tabular-nums">
          Nov 15 - Nov 28, 2026
        </span>
        <button
          @click="store.showToast('Generating Executive summary report...')"
          class="px-3 py-1 rounded-md bg-zinc-900 dark:bg-white hover:bg-zinc-800 dark:hover:bg-zinc-100 text-white dark:text-zinc-900 font-medium transition-colors text-xs"
        >
          Export
        </button>
      </div>
    </div>

    <!-- 4 KPI Summary Cards -->
    <KPIStats />

    <!-- TWO-COLUMN SECTION: PROJECT PROGRESS & RECENT ACTIVITY FEED -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-5">
      <!-- Project Progress List (2 cols) -->
      <div class="lg:col-span-2 p-4 sm:p-5 rounded-xl bg-white dark:bg-[#18181B] border border-[#E4E4E7] dark:border-[#27272A] shadow-subtle space-y-4">
        <div class="flex items-center justify-between">
          <div>
            <h2 class="text-xs font-bold uppercase tracking-wider text-zinc-400 font-mono">Active Projects Portfolio</h2>
            <p class="text-xs text-zinc-500 dark:text-zinc-400">Deliverables completion rate and client status</p>
          </div>
          <RouterLink to="/projects" class="text-xs text-brand-600 dark:text-brand-400 hover:underline font-medium">
            View all {{ store.projects.length }} →
          </RouterLink>
        </div>

        <!-- Project Items Table-like Rows -->
        <div class="space-y-2.5 pt-1">
          <div
            v-for="proj in store.projects"
            :key="proj.id"
            @click="selectProject(proj.id)"
            class="p-3 rounded-lg border border-[#E4E4E7] dark:border-[#27272A] bg-[#F4F4F5]/40 dark:bg-[#27272A]/30 hover:bg-[#F4F4F5] dark:hover:bg-[#27272A]/70 cursor-pointer transition-colors"
          >
            <div class="flex items-center justify-between mb-2">
              <div class="flex items-center space-x-2">
                <span class="w-2.5 h-2.5 rounded-full" :style="{ backgroundColor: proj.color_dot }"></span>
                <span class="text-xs font-semibold text-zinc-900 dark:text-white">{{ proj.name }}</span>
                <span
                  v-if="proj.client_name"
                  class="px-2 py-0.5 rounded-full text-[10px] font-mono bg-indigo-50 dark:bg-indigo-950/60 text-indigo-600 dark:text-indigo-400 border border-indigo-200/40"
                >
                  Client: {{ proj.client_name }}
                </span>
                <span
                  v-else
                  class="px-2 py-0.5 rounded-full text-[10px] font-mono bg-zinc-200/80 dark:bg-zinc-800 text-zinc-700 dark:text-zinc-300"
                >
                  Internal
                </span>
              </div>
              <span class="text-xs font-mono font-semibold text-zinc-700 dark:text-zinc-300 tabular-nums">
                {{ proj.progress_percentage }}%
              </span>
            </div>

            <!-- Progress Bar -->
            <div class="w-full bg-zinc-200 dark:bg-zinc-800 h-1.5 rounded-full overflow-hidden">
              <div
                class="h-full rounded-full transition-all duration-500"
                :style="{ width: proj.progress_percentage + '%', backgroundColor: proj.color_dot }"
              ></div>
            </div>

            <div class="mt-2 flex items-center justify-between text-[11px] text-zinc-400 font-mono">
              <span class="tabular-nums">{{ proj.completed_tasks }} of {{ proj.total_tasks }} tasks completed</span>
              <span>Deadline: {{ proj.deadline }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Activity Feed & Upcoming Deadlines (1 col) -->
      <div class="p-4 sm:p-5 rounded-xl bg-white dark:bg-[#18181B] border border-[#E4E4E7] dark:border-[#27272A] shadow-subtle space-y-4">
        <div class="flex items-center justify-between">
          <h2 class="text-xs font-bold uppercase tracking-wider text-zinc-400 font-mono">Activity Feed</h2>
          <span class="text-[10px] font-mono text-emerald-500 flex items-center gap-1">
            <span class="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse"></span> Live
          </span>
        </div>

        <!-- Activity Stream -->
        <div class="space-y-3 text-xs">
          <div
            v-for="act in store.activityFeed"
            :key="act.id"
            class="flex items-start space-x-2.5"
          >
            <img class="w-5 h-5 rounded-full object-cover mt-0.5 ring-1 ring-zinc-200 dark:ring-zinc-700 flex-shrink-0" :src="act.user.avatar_url" :alt="act.user.name" />
            <div class="space-y-0.5 min-w-0">
              <p class="text-zinc-800 dark:text-zinc-200 text-xs truncate">
                <span class="font-semibold text-zinc-900 dark:text-white">{{ act.user.name }}</span>
                {{ ' ' }} {{ act.action }}
                <span class="font-medium text-brand-600 dark:text-brand-400">{{ act.target }}</span>
              </p>
              <span class="text-[10px] text-zinc-400 font-mono">{{ act.timestamp }}</span>
            </div>
          </div>
        </div>

        <!-- Upcoming Deadlines Micro Widget -->
        <div class="pt-3 border-t border-zinc-100 dark:border-zinc-800 space-y-2">
          <span class="text-xs font-bold text-zinc-900 dark:text-white flex items-center gap-1.5 font-display">
            <Clock class="w-3.5 h-3.5 text-amber-500" />
            Critical Deadlines
          </span>
          <div class="space-y-1.5 text-[11px] font-mono">
            <div class="flex items-center justify-between p-2 rounded-lg bg-rose-50/60 dark:bg-rose-950/20 text-rose-700 dark:text-rose-300">
              <span class="truncate">App Store Release</span>
              <span class="font-bold flex-shrink-0 ml-1">Tomorrow</span>
            </div>
            <div class="flex items-center justify-between p-2 rounded-lg bg-zinc-50 dark:bg-zinc-900/60 text-zinc-600 dark:text-zinc-400">
              <span class="truncate">Stripe Live Sandbox</span>
              <span class="flex-shrink-0 ml-1">In 3 days</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useProjectStore } from '../stores/projectStore'
import { useRouter } from 'vue-router'
import KPIStats from '../components/KPIStats.vue'
import { Clock } from 'lucide-vue-next'

const store = useProjectStore()
const router = useRouter()

function selectProject(projectId: string) {
  store.selectProject(projectId)
  router.push('/projects')
}
</script>
