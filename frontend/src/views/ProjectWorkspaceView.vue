<template>
  <div class="h-full flex flex-col space-y-3 min-h-0">
    <!-- Clean Subheader: Active Project Title, Category Chips, and Share Link -->
    <section class="flex flex-wrap items-center justify-between gap-2 flex-shrink-0">
      <div class="flex items-center space-x-2">
        <h1 class="text-sm font-bold text-zinc-900 dark:text-zinc-100 font-display tracking-tight flex items-center gap-2">
          <span>{{ store.currentProject?.name }}</span>
          <span class="text-[10px] font-mono px-2 py-0.5 rounded-full bg-zinc-100 dark:bg-zinc-800 text-zinc-500 font-normal">
            {{ store.filteredTasks.length }} tasks
          </span>
        </h1>
      </div>

      <div class="flex items-center space-x-2 text-xs">
        <!-- Quick Category Filter Chips -->
        <div class="hidden sm:flex items-center space-x-1">
          <button
            v-for="cat in quickCategories"
            :key="cat"
            @click="store.setFilter(cat)"
            :class="[
              'px-2 py-0.5 rounded-md text-[11px] font-medium transition-colors',
              store.activeFilter === cat
                ? 'bg-zinc-900 dark:bg-white text-white dark:text-zinc-900 font-semibold shadow-2xs'
                : 'text-zinc-500 hover:text-zinc-900 dark:hover:text-zinc-200 hover:bg-zinc-100 dark:hover:bg-zinc-800'
            ]"
          >
            {{ cat }}
          </button>
        </div>

        <!-- Share Workspace Button -->
        <button
          @click="store.showToast('Copied workspace share link!')"
          class="flex items-center space-x-1.5 px-2.5 py-1 rounded-lg border border-[#E4E4E7] dark:border-[#27272A] bg-white dark:bg-[#18181B] hover:bg-zinc-50 dark:hover:bg-zinc-800 text-zinc-600 dark:text-zinc-300 transition-colors text-[11px]"
        >
          <Share2 class="w-3.5 h-3.5 text-zinc-400" />
          <span class="hidden md:inline">Share</span>
        </button>
      </div>
    </section>

    <!-- MAIN VIEW AREA -->
    <div class="flex-1 min-h-0 flex flex-col">
      <!-- 1. KANBAN BOARD -->
      <KanbanBoard v-if="store.currentView === 'kanban'" />

      <!-- 2. LIST VIEW -->
      <TaskListView v-else-if="store.currentView === 'list'" />

      <!-- 3. TIMELINE (GANTT) VIEW -->
      <div v-else-if="store.currentView === 'timeline'" class="space-y-4 max-w-7xl mx-auto w-full overflow-y-auto">
        <div class="flex items-center justify-between">
          <div>
            <h2 class="text-sm font-bold text-zinc-900 dark:text-white font-display">Sprint & Deliverables Roadmap</h2>
            <p class="text-xs text-zinc-400">Critical path dependencies, milestone gates, and deliverable signoffs.</p>
          </div>
          <div class="flex items-center space-x-2 text-xs font-mono">
            <span class="px-2 py-0.5 rounded bg-zinc-100 dark:bg-zinc-800 text-zinc-600 dark:text-zinc-400">Week 47</span>
            <span class="px-2 py-0.5 rounded bg-brand-50 dark:bg-brand-950 text-brand-600 dark:text-brand-400 font-semibold border border-brand-200/50">Current Week 48</span>
            <span class="px-2 py-0.5 rounded bg-zinc-100 dark:bg-zinc-800 text-zinc-600 dark:text-zinc-400">Week 49</span>
          </div>
        </div>

        <div class="rounded-xl border border-[#E4E4E7] dark:border-[#27272A] bg-white dark:bg-[#18181B] p-4 shadow-subtle overflow-x-auto">
          <div class="min-w-[700px]">
            <div class="grid grid-cols-12 gap-1 text-[11px] font-mono text-zinc-400 border-b border-zinc-100 dark:border-zinc-800 pb-2 mb-4">
              <div class="col-span-4 font-semibold text-zinc-700 dark:text-zinc-300">Deliverable / Task</div>
              <div class="col-span-2 text-center">Nov 16 - Nov 20</div>
              <div class="col-span-3 text-center bg-brand-50/50 dark:bg-brand-950/30 rounded py-0.5 text-brand-600 dark:text-brand-400 font-bold">Nov 21 - Nov 27 (Now)</div>
              <div class="col-span-3 text-center">Nov 28 - Dec 04</div>
            </div>

            <!-- Gantt Items -->
            <div class="space-y-4 text-xs">
              <div class="grid grid-cols-12 gap-1 items-center">
                <div class="col-span-4 flex items-center space-x-2">
                  <span class="w-2 h-2 rounded-full bg-rose-500"></span>
                  <span class="font-medium text-zinc-900 dark:text-white truncate">FaceID Auth Flow</span>
                </div>
                <div class="col-span-8 relative h-7 bg-zinc-50 dark:bg-zinc-900/40 rounded-lg flex items-center">
                  <div
                    @click="store.openTaskDetail('task_01')"
                    class="absolute left-[30%] w-[35%] h-5 bg-gradient-to-r from-rose-500 to-rose-600 rounded-md shadow-2xs flex items-center justify-between px-2 text-[10px] text-white font-mono font-medium cursor-pointer hover:opacity-95"
                  >
                    <span>#EGN-104</span>
                    <span>25%</span>
                  </div>
                </div>
              </div>

              <div class="grid grid-cols-12 gap-1 items-center">
                <div class="col-span-4 flex items-center space-x-2">
                  <span class="w-2 h-2 rounded-full bg-amber-500"></span>
                  <span class="font-medium text-zinc-900 dark:text-white truncate">1-Click Apple/Google Pay</span>
                </div>
                <div class="col-span-8 relative h-7 bg-zinc-50 dark:bg-zinc-900/40 rounded-lg flex items-center">
                  <div
                    @click="store.openTaskDetail('task_04')"
                    class="absolute left-[20%] w-[55%] h-5 bg-gradient-to-r from-amber-500 to-amber-600 rounded-md shadow-2xs flex items-center justify-between px-2 text-[10px] text-white font-mono font-medium cursor-pointer hover:opacity-95"
                  >
                    <span>#EGN-99</span>
                    <span>75%</span>
                  </div>
                </div>
              </div>

              <!-- Milestone Gate Marker -->
              <div class="grid grid-cols-12 gap-1 items-center py-2 border-y border-dashed border-indigo-200 dark:border-indigo-900/60 bg-indigo-50/30 dark:bg-indigo-950/20 rounded">
                <div class="col-span-4 flex items-center space-x-2 pl-2">
                  <Flag class="w-3.5 h-3.5 text-indigo-600 dark:text-indigo-400" />
                  <span class="font-bold text-indigo-700 dark:text-indigo-300">Milestone 2 Delivery Gate</span>
                </div>
                <div class="col-span-8 relative h-5 flex items-center">
                  <div class="absolute left-[65%] w-3 h-3 rotate-45 bg-indigo-600 ring-4 ring-indigo-200 dark:ring-indigo-900"></div>
                  <span class="absolute left-[68%] text-[10px] font-mono text-indigo-600 dark:text-indigo-400 font-bold">Client Review Signoff (Nov 28)</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useProjectStore } from '../stores/projectStore'
import KanbanBoard from '../components/KanbanBoard.vue'
import TaskListView from '../components/TaskListView.vue'
import { Share2, Flag } from 'lucide-vue-next'

const store = useProjectStore()
const route = useRoute()

const quickCategories = ['All', 'Frontend', 'Backend', 'AI/ML', 'QA/Test']

onMounted(() => {
  if (route.query.tab === 'list' || route.query.tab === 'timeline' || route.query.tab === 'kanban') {
    store.setView(route.query.tab as 'kanban' | 'list' | 'timeline')
  }
})

watch(() => route.query.tab, (newTab) => {
  if (newTab === 'list' || newTab === 'timeline' || newTab === 'kanban') {
    store.setView(newTab as 'kanban' | 'list' | 'timeline')
  }
})
</script>
