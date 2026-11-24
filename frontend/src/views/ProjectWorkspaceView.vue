<template>
  <div class="h-full flex flex-col space-y-3 min-h-0">
    <!-- Empty State: No Project Assigned to Current User -->
    <div
      v-if="!store.currentProject"
      class="flex-1 flex flex-col items-center justify-center p-8 text-center bg-white rounded-3xl border border-slate-200/90 shadow-sm space-y-4 my-auto select-none"
    >
      <div class="w-16 h-16 rounded-2xl bg-indigo-50 border border-indigo-100 text-indigo-600 flex items-center justify-center shadow-2xs">
        <FolderKanban class="w-8 h-8" />
      </div>
      <div class="max-w-md space-y-2">
        <h3 class="text-base sm:text-lg font-bold text-slate-900">Belum Ada Proyek yang Ditugaskan</h3>
        <p class="text-xs text-slate-500 leading-relaxed">
          Akun Anda saat ini belum memiliki akses ke proyek mana pun. Silakan hubungi Administrator atau Project Lead untuk memberikan hak akses melalui menu <strong>Akses Tim</strong> pada dashboard proyek.
        </p>
      </div>
    </div>

    <!-- Active Project Content -->
    <template v-else>
      <!-- Clean Subheader: Active Project Title, Category Chips, and Share Link -->
      <section class="flex flex-wrap items-center justify-between gap-2 flex-shrink-0">
      <div class="flex items-center space-x-2">
        <h1 class="text-sm font-bold text-zinc-900 dark:text-zinc-100 font-display tracking-tight flex items-center gap-2">
          <span>{{ store.currentProject?.name }}</span>
          <span class="text-[10px] font-mono px-2 py-0.5 rounded-full bg-zinc-100 dark:bg-zinc-800 text-zinc-500 font-normal">
            {{ store.filteredTasks.length }} tasks
          </span>
        </h1>

        <!-- Team Access Button for Admin / PM -->
        <button
          v-if="canManageProjectMembers"
          @click="showProjectAccessModal = true"
          class="flex items-center space-x-1 px-2 py-0.5 rounded-md bg-zinc-100 dark:bg-zinc-800 hover:bg-zinc-200 dark:hover:bg-zinc-700 text-zinc-700 dark:text-zinc-200 text-[10px] font-semibold transition-all cursor-pointer border border-zinc-200 dark:border-zinc-700 ml-1.5"
          title="Kelola hak akses team member ke proyek ini"
        >
          <ShieldCheck class="w-3 h-3 text-brand-600" />
          <span>Akses Tim</span>
        </button>

        <!-- Edit Project Button -->
        <button
          v-if="canManageProjectMembers"
          @click="showEditProjectModal = true"
          class="flex items-center space-x-1 px-2 py-0.5 rounded-md bg-zinc-100 dark:bg-zinc-800 hover:bg-zinc-200 dark:hover:bg-zinc-700 text-zinc-700 dark:text-zinc-200 text-[10px] font-semibold transition-all cursor-pointer border border-zinc-200 dark:border-zinc-700"
          title="Edit data proyek (Nama, Klien, Status, Deadline)"
        >
          <Edit class="w-3 h-3 text-zinc-600 dark:text-zinc-400" />
          <span>Edit</span>
        </button>

        <!-- Copy Project Button -->
        <button
          v-if="canManageProjectMembers"
          @click="showCopyProjectModal = true"
          class="flex items-center space-x-1 px-2 py-0.5 rounded-md bg-indigo-50 dark:bg-indigo-950/40 hover:bg-indigo-100 dark:hover:bg-indigo-900/60 text-indigo-700 dark:text-indigo-300 text-[10px] font-semibold transition-all cursor-pointer border border-indigo-200 dark:border-indigo-800"
          title="Duplikasi seluruh pekerjaan proyek ini ke proyek baru"
        >
          <Copy class="w-3 h-3 text-indigo-600" />
          <span>Salin Proyek</span>
        </button>

        <!-- Delete Project Button -->
        <button
          v-if="canManageProjectMembers"
          @click="handleDeleteActiveProject"
          class="flex items-center space-x-1 px-2 py-0.5 rounded-md bg-rose-50 dark:bg-rose-950/40 hover:bg-rose-100 dark:hover:bg-rose-900/60 text-rose-700 dark:text-rose-300 text-[10px] font-semibold transition-all cursor-pointer border border-rose-200 dark:border-rose-800"
          title="Hapus proyek ini"
        >
          <Trash2 class="w-3 h-3 text-rose-600" />
          <span>Hapus</span>
        </button>
      </div>

      <div class="flex items-center space-x-2 text-xs">
        <!-- Quick View Tabs: Board | List | Timeline | Wiki -->
        <div class="flex items-center p-0.5 rounded-lg bg-zinc-100 dark:bg-zinc-800 border border-zinc-200 dark:border-zinc-700 text-xs">
          <button
            @click="switchView('kanban')"
            :class="[
              'px-2.5 py-1 rounded-md text-[11px] font-medium transition-all flex items-center gap-1',
              store.currentView === 'kanban'
                ? 'bg-white dark:bg-zinc-900 text-brand-600 dark:text-brand-400 font-semibold shadow-2xs'
                : 'text-zinc-500 hover:text-zinc-800 dark:hover:text-zinc-200'
            ]"
          >
            <Kanban class="w-3.5 h-3.5" />
            <span>Board</span>
          </button>

          <button
            @click="switchView('list')"
            :class="[
              'px-2.5 py-1 rounded-md text-[11px] font-medium transition-all flex items-center gap-1',
              store.currentView === 'list'
                ? 'bg-white dark:bg-zinc-900 text-brand-600 dark:text-brand-400 font-semibold shadow-2xs'
                : 'text-zinc-500 hover:text-zinc-800 dark:hover:text-zinc-200'
            ]"
          >
            <ListTodo class="w-3.5 h-3.5" />
            <span>List</span>
          </button>

          <button
            @click="switchView('timeline')"
            :class="[
              'px-2.5 py-1 rounded-md text-[11px] font-medium transition-all flex items-center gap-1',
              store.currentView === 'timeline'
                ? 'bg-white dark:bg-zinc-900 text-brand-600 dark:text-brand-400 font-semibold shadow-2xs'
                : 'text-zinc-500 hover:text-zinc-800 dark:hover:text-zinc-200'
            ]"
          >
            <GitCommit class="w-3.5 h-3.5" />
            <span>Timeline</span>
          </button>

          <button
            @click="activeSubTab = 'wiki'"
            :class="[
              'px-2.5 py-1 rounded-md text-[11px] font-medium transition-all flex items-center gap-1',
              activeSubTab === 'wiki'
                ? 'bg-white dark:bg-zinc-900 text-brand-600 dark:text-brand-400 font-semibold shadow-2xs'
                : 'text-zinc-500 hover:text-zinc-800 dark:hover:text-zinc-200'
            ]"
          >
            <BookOpen class="w-3.5 h-3.5" />
            <span>Wiki</span>
          </button>
        </div>

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

        <!-- New Task Button -->
        <button
          v-can="'task:create'"
          @click="store.isNewTaskModalOpen = true"
          class="flex items-center space-x-1 px-2.5 py-1 rounded-lg bg-brand-600 hover:bg-brand-500 text-white text-[11px] font-semibold transition-all active:scale-[0.98]"
        >
          <Plus class="w-3.5 h-3.5" />
          <span>New Task</span>
        </button>

        <!-- Export Project Tasks to CSV / MS Excel -->
        <button
          @click="handleExportProjectTasks"
          :disabled="isExportingTasks"
          class="flex items-center space-x-1 px-2.5 py-1 rounded-lg bg-white dark:bg-zinc-800 hover:bg-zinc-50 dark:hover:bg-zinc-700 border border-zinc-200 dark:border-zinc-700 text-zinc-700 dark:text-zinc-200 text-[11px] font-medium transition-all cursor-pointer shadow-2xs"
          title="Ekspor tasks proyek ini ke MS Excel / CSV"
        >
          <Download class="w-3.5 h-3.5 text-zinc-500" />
          <span>Ekspor CSV</span>
        </button>
      </div>
    </section>

    <!-- MAIN VIEW AREA -->
    <div class="flex-1 min-h-0 flex flex-col">
      <!-- 1. KANBAN BOARD -->
      <KanbanBoard v-if="activeSubTab === 'views' && store.currentView === 'kanban'" />

      <!-- 2. LIST VIEW (With Inline Editing) -->
      <TaskListView v-else-if="activeSubTab === 'views' && store.currentView === 'list'" />

      <!-- 3. TIMELINE (GANTT) VIEW — 100% Dynamic from Real Tasks -->
      <div v-else-if="activeSubTab === 'views' && store.currentView === 'timeline'" class="space-y-4 max-w-7xl mx-auto w-full overflow-y-auto">
        <div class="flex items-center justify-between">
          <div>
            <h2 class="text-sm font-bold text-zinc-900 dark:text-white font-display">Sprint & Deliverables Roadmap</h2>
            <p class="text-xs text-zinc-400">Jalur kritis, timeline tenggat waktu, dan status deliverable langsung dari database.</p>
          </div>
          <div class="flex items-center space-x-2 text-xs font-mono">
            <span class="px-2 py-0.5 rounded bg-zinc-100 dark:bg-zinc-800 text-zinc-600 dark:text-zinc-400">Total: {{ store.filteredTasks.length }} Tasks</span>
            <span class="px-2 py-0.5 rounded bg-emerald-50 dark:bg-emerald-950 text-emerald-600 dark:text-emerald-400 font-semibold border border-emerald-200/50">
              Selesai: {{ store.taskCounts['Done'] }}
            </span>
          </div>
        </div>

        <div class="rounded-xl border border-[#E4E4E7] dark:border-[#27272A] bg-white dark:bg-[#18181B] p-4 shadow-subtle overflow-x-auto">
          <div class="min-w-[700px] space-y-3">
            <div class="grid grid-cols-12 gap-1 text-[11px] font-mono text-zinc-400 border-b border-zinc-100 dark:border-zinc-800 pb-2">
              <div class="col-span-5 font-semibold text-zinc-700 dark:text-zinc-300">Deliverable / Task</div>
              <div class="col-span-2 text-center">Assignee & Priority</div>
              <div class="col-span-2 text-center">Due Date</div>
              <div class="col-span-3 text-right">Progress Status</div>
            </div>

            <!-- Dynamic Gantt Task Rows -->
            <div
              v-for="task in store.filteredTasks"
              :key="task.id"
              class="grid grid-cols-12 gap-2 items-center p-2 rounded-lg hover:bg-zinc-50 dark:hover:bg-zinc-800/40 transition-colors group cursor-pointer"
              @click="store.openTaskDetail(task.id)"
            >
              <div class="col-span-5 flex items-center space-x-2 min-w-0 pr-2">
                <span :class="['w-2 h-2 rounded-full flex-shrink-0', statusDots[task.status]]"></span>
                <span class="font-mono text-[10px] text-zinc-400 flex-shrink-0">{{ task.key }}</span>
                <span class="font-medium text-xs text-zinc-900 dark:text-white truncate">{{ task.title }}</span>
              </div>

              <div class="col-span-2 flex items-center justify-center space-x-1.5">
                <img
                  :src="task.assignee?.avatar_url || 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=100&auto=format&fit=crop&q=80'"
                  :alt="task.assignee?.name || 'Unassigned'"
                  class="w-5 h-5 rounded-full object-cover"
                />
                <span :class="['text-[9px] font-mono px-1 py-0.2 rounded font-semibold', priorityPill[task.priority]]">
                  {{ task.priority }}
                </span>
              </div>

              <div class="col-span-2 text-center font-mono text-xs text-zinc-500 tabular-nums">
                {{ task.due_date || 'Ongoing' }}
              </div>

              <div class="col-span-3 relative h-6 bg-zinc-100 dark:bg-zinc-800 rounded-md overflow-hidden flex items-center">
                <div
                  class="h-full transition-all duration-300"
                  :class="[
                    task.status === 'Done' ? 'bg-emerald-500 w-full' :
                    task.status === 'Review' ? 'bg-amber-500 w-3/4' :
                    task.status === 'In Progress' ? 'bg-blue-500 w-1/2' : 'bg-zinc-400 w-1/4'
                  ]"
                ></div>
                <span class="absolute right-2 text-[10px] font-mono font-bold text-zinc-700 dark:text-zinc-200">
                  {{ task.status }}
                </span>
              </div>
            </div>

            <!-- Milestone Marker Bar -->
            <div
              v-for="ms in currentMilestones"
              :key="ms.id"
              class="p-2.5 rounded-lg border border-dashed border-indigo-200 dark:border-indigo-900/60 bg-indigo-50/40 dark:bg-indigo-950/20 flex items-center justify-between text-xs"
            >
              <div class="flex items-center space-x-2">
                <Flag class="w-4 h-4 text-indigo-600 dark:text-indigo-400" />
                <span class="font-bold text-indigo-900 dark:text-indigo-200">{{ ms.title }}</span>
                <span class="text-[10px] text-zinc-400">({{ ms.description || 'Milestone Delivery Gate' }})</span>
              </div>
              <div class="flex items-center space-x-2 font-mono text-[11px]">
                <span class="text-zinc-500">Tenggat: {{ ms.due_date }}</span>
                <span :class="['px-2 py-0.5 rounded uppercase font-bold text-[9px]', ms.status === 'approved' ? 'bg-emerald-500 text-white' : 'bg-amber-100 text-amber-800 dark:bg-amber-900 dark:text-amber-200']">
                  {{ ms.status }}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 4. WIKI & KNOWLEDGE BASE VIEW -->
      <div v-else-if="activeSubTab === 'wiki'" class="space-y-4 max-w-5xl mx-auto w-full overflow-y-auto pb-8">
        <div class="flex items-center justify-between">
          <div>
            <h2 class="text-sm font-bold text-zinc-900 dark:text-white font-display">Project Knowledge Base & Wiki</h2>
            <p class="text-xs text-zinc-400">Dokumentasi arsitektur, SOP integrasi, dan spesifikasi deliverable proyek.</p>
          </div>
          <button
            v-can="'project:edit'"
            @click="showNewWikiModal = true"
            class="px-3 py-1.5 rounded-lg bg-brand-600 hover:bg-brand-500 text-white text-xs font-semibold flex items-center gap-1.5 shadow-2xs"
          >
            <Plus class="w-3.5 h-3.5" />
            <span>Tambah Dokumen</span>
          </button>
        </div>

        <div v-if="wikiPages.length === 0" class="p-8 text-center bg-white dark:bg-[#18181B] rounded-2xl border border-zinc-200 dark:border-zinc-800 space-y-2">
          <BookOpen class="w-8 h-8 text-zinc-400 mx-auto" />
          <h3 class="text-xs font-bold text-zinc-700 dark:text-zinc-300">Belum Ada Dokumen Wiki</h3>
          <p class="text-[11px] text-zinc-400">Tambahkan dokumentasi proyek baru untuk referensi tim dan klien.</p>
        </div>

        <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div
            v-for="page in wikiPages"
            :key="page.id"
            class="p-4 rounded-xl bg-white dark:bg-[#18181B] border border-zinc-200 dark:border-zinc-800 shadow-subtle space-y-2"
          >
            <div class="flex items-center justify-between">
              <span class="text-xs font-bold text-zinc-900 dark:text-white">{{ page.title }}</span>
              <span v-if="page.client_visible" class="text-[9px] font-mono px-1.5 py-0.2 rounded bg-indigo-50 dark:bg-indigo-950 text-indigo-600 dark:text-indigo-400 font-semibold">
                Client Shared
              </span>
              <span v-else class="text-[9px] font-mono px-1.5 py-0.2 rounded bg-zinc-100 dark:bg-zinc-800 text-zinc-500 font-semibold">
                Internal Only
              </span>
            </div>
            <p class="text-xs text-zinc-600 dark:text-zinc-300 whitespace-pre-wrap line-clamp-3">
              {{ page.content }}
            </p>
            <div class="pt-2 border-t border-zinc-100 dark:border-zinc-800 flex items-center justify-between text-[10px] text-zinc-400 font-mono">
              <span>Author: {{ page.author_name || 'Team' }}</span>
              <span>{{ page.created_at ? page.created_at.slice(0, 10) : 'Active' }}</span>
            </div>
          </div>
        </div>

        <!-- New Wiki Modal Dialog -->
        <div
          v-if="showNewWikiModal"
          class="fixed inset-0 bg-black/50 backdrop-blur-xs flex items-center justify-center p-4 z-50"
        >
          <div class="max-w-lg w-full bg-white dark:bg-[#18181B] rounded-2xl p-5 border border-zinc-200 dark:border-zinc-800 shadow-2xl space-y-4">
            <h3 class="text-sm font-bold text-zinc-900 dark:text-white">Tambah Dokumen Wiki Baru</h3>
            <div class="space-y-3 text-xs">
              <div>
                <label class="font-semibold text-zinc-700 dark:text-zinc-300">Judul Dokumen</label>
                <input
                  v-model="newWikiTitle"
                  placeholder="Contoh: Arsitektur API Gateway & Keamanan TLS"
                  class="w-full mt-1 px-3 py-2 rounded-lg border border-zinc-300 dark:border-zinc-700 bg-white dark:bg-zinc-900 text-zinc-900 dark:text-white focus:outline-none focus:ring-1 focus:ring-brand-500"
                />
              </div>
              <div>
                <label class="font-semibold text-zinc-700 dark:text-zinc-300">Konten Dokumen</label>
                <textarea
                  v-model="newWikiContent"
                  rows="4"
                  placeholder="Tulis ringkasan atau dokumentasi SOP di sini..."
                  class="w-full mt-1 px-3 py-2 rounded-lg border border-zinc-300 dark:border-zinc-700 bg-white dark:bg-zinc-900 text-zinc-900 dark:text-white focus:outline-none focus:ring-1 focus:ring-brand-500"
                ></textarea>
              </div>
              <div class="flex items-center gap-2">
                <input type="checkbox" v-model="newWikiClientVisible" id="wikiClient" class="rounded text-brand-600" />
                <label for="wikiClient" class="text-zinc-600 dark:text-zinc-400">Bagikan dokumen ini ke Klien (Client Shared)</label>
              </div>
            </div>
            <div class="flex items-center justify-end space-x-2 pt-2 border-t border-zinc-100 dark:border-zinc-800">
              <button
                @click="showNewWikiModal = false"
                class="px-3 py-1.5 rounded-lg border border-zinc-200 dark:border-zinc-700 text-xs font-semibold text-zinc-600 dark:text-zinc-400 hover:bg-zinc-100 dark:hover:bg-zinc-800"
              >
                Batal
              </button>
              <button
                @click="saveNewWiki"
                :disabled="!newWikiTitle.trim() || !newWikiContent.trim()"
                class="px-4 py-1.5 rounded-lg bg-brand-600 hover:bg-brand-500 text-white text-xs font-bold disabled:opacity-50"
              >
                Simpan Dokumen
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
    </template>

    <!-- Project Access Management Modal -->
    <ProjectAccessModal
      :is-open="showProjectAccessModal"
      :project="store.currentProject || null"
      @close="showProjectAccessModal = false"
      @updated="store.fetchProjects()"
    />

    <!-- Edit Project Modal -->
    <EditProjectModal
      :is-open="showEditProjectModal"
      :project="store.currentProject || null"
      @close="showEditProjectModal = false"
      @updated="store.fetchProjects()"
    />

    <!-- Copy Project Modal -->
    <CopyProjectModal
      :is-open="showCopyProjectModal"
      :project="store.currentProject || null"
      @close="showCopyProjectModal = false"
      @copied="onProjectCopied"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useProjectStore } from '../stores/projectStore'
import { useAuthStore } from '../stores/authStore'
import { projectApi } from '../services/api'
import type { WikiPage, Project } from '../types/project'
import KanbanBoard from '../components/KanbanBoard.vue'
import TaskListView from '../components/TaskListView.vue'
import ProjectAccessModal from '../components/ProjectAccessModal.vue'
import EditProjectModal from '../components/EditProjectModal.vue'
import CopyProjectModal from '../components/CopyProjectModal.vue'
import { Kanban, ListTodo, GitCommit, BookOpen, Flag, Plus, Download, ShieldCheck, FolderKanban, Edit, Trash2, Copy } from 'lucide-vue-next'

const store = useProjectStore()
const authStore = useAuthStore()
const route = useRoute()
const router = useRouter()

const showProjectAccessModal = ref(false)
const showEditProjectModal = ref(false)
const showCopyProjectModal = ref(false)

async function onProjectCopied(newProj: Project) {
  await store.fetchProjects()
  store.selectProject(newProj.id)
  store.showToast(`Proyek "${newProj.name}" berhasil diduplikasi dan dibuka!`)
}

const canManageProjectMembers = computed(() => {
  return authStore.hasRole('super_admin') || authStore.hasRole('pm') || authStore.can('project:member:manage')
})

async function handleDeleteActiveProject() {
  const proj = store.currentProject
  if (!proj) return
  if (!confirm(`Hapus proyek "${proj.name}"?\nPerhatian: Seluruh deliverable, tugas, dan data keanggotaan proyek ini akan ikut terhapus secara permanen.`)) {
    return
  }

  try {
    await projectApi.deleteProject(proj.id)
    store.projects = store.projects.filter(p => p.id !== proj.id)
    store.showToast(`Proyek "${proj.name}" berhasil dihapus.`)
    if (store.projects.length > 0) {
      store.selectProject(store.projects[0].id)
    } else {
      router.push('/dashboard')
    }
  } catch (err: any) {
    store.showToast(err?.response?.data?.error || 'Gagal menghapus proyek')
  }
}

const isExportingTasks = ref(false)

async function handleExportProjectTasks() {
  if (!store.currentProjectId) return
  isExportingTasks.value = true
  try {
    store.showToast('Menyiapkan file ekspor tasks proyek ini (MS Excel / CSV)...')
    const blob = await projectApi.exportTasksCSV(store.currentProjectId)
    const timestamp = new Date().toISOString().slice(0, 10).replace(/-/g, '')
    const url = window.URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `tasks_${store.currentProject?.slug || 'project'}_${timestamp}.csv`
    document.body.appendChild(a)
    a.click()
    window.URL.revokeObjectURL(url)
    document.body.removeChild(a)
    store.showToast('File CSV berhasil diunduh!')
  } catch (err: any) {
    store.showToast(err.response?.data?.error || 'Gagal mengekspor tasks.')
  } finally {
    isExportingTasks.value = false
  }
}

const activeSubTab = ref<'views' | 'wiki'>('views')
const quickCategories = ['All', 'Frontend', 'Backend', 'AI/ML', 'QA/Test', 'Infrastructure']

const currentMilestones = computed(() => {
  return store.milestones.filter((m) => m.project_id === store.currentProjectId)
})

const statusDots: Record<string, string> = {
  'Backlog': 'bg-zinc-400',
  'In Progress': 'bg-blue-500',
  'Review': 'bg-amber-500',
  'Done': 'bg-emerald-500',
}

const priorityPill: Record<string, string> = {
  'Urgent': 'bg-rose-50 text-rose-600 dark:bg-rose-950 dark:text-rose-400',
  'High': 'bg-orange-50 text-orange-600 dark:bg-orange-950 dark:text-orange-400',
  'Medium': 'bg-blue-50 text-blue-600 dark:bg-blue-950 dark:text-blue-400',
  'Low': 'bg-zinc-100 text-zinc-600 dark:bg-zinc-800 dark:text-zinc-400',
}

const wikiPages = ref<WikiPage[]>([
  {
    id: 'wiki_01',
    project_id: 'proj_sustagri',
    title: 'Spesifikasi Arsitektur Sistem & API RBAC',
    slug: 'spesifikasi-arsitektur',
    content: 'Aplikasi SustAgri menggunakan backend Golang Gin dengan enkripsi JWT HS256 dan RBAC 4 role (Super Admin, PM, Member, Client). Keamanan jaringan diperkuat dengan Cloudflare Tunnel TLS 1.3.',
    author_name: 'Agus Setiawan',
    client_visible: true,
    created_at: '2026-11-20',
    updated_at: '2026-11-20'
  },
  {
    id: 'wiki_02',
    project_id: 'proj_sustagri',
    title: 'Panduan Deployment Windows Server Host',
    slug: 'panduan-deploy',
    content: 'Daemon produksi berjalan pada host Windows lokal (Node.js gateway port 5080 + Go binary port 5081) dengan auto-restart dan logging log rotasi.',
    author_name: 'Dedi Irawan',
    client_visible: false,
    created_at: '2026-11-21',
    updated_at: '2026-11-21'
  }
])

const showNewWikiModal = ref(false)
const newWikiTitle = ref('')
const newWikiContent = ref('')
const newWikiClientVisible = ref(true)

function switchView(view: 'kanban' | 'list' | 'timeline') {
  activeSubTab.value = 'views'
  store.setView(view)
}

async function loadWikiPages() {
  try {
    const pages = await projectApi.getWikiPages(store.currentProjectId)
    wikiPages.value = pages && pages.length > 0 ? pages : []
  } catch (_) {
    wikiPages.value = []
  }
}

async function saveNewWiki() {
  if (!newWikiTitle.value.trim() || !newWikiContent.value.trim()) return
  try {
    const created = await projectApi.createWikiPage(store.currentProjectId, {
      title: newWikiTitle.value.trim(),
      content: newWikiContent.value.trim(),
      client_visible: newWikiClientVisible.value
    })
    wikiPages.value.unshift(created)
    store.showToast('Dokumen wiki berhasil ditambahkan!')
    showNewWikiModal.value = false
    newWikiTitle.value = ''
    newWikiContent.value = ''
  } catch (_) {
    // Local fallback
    wikiPages.value.unshift({
      id: `wiki_${Date.now()}`,
      project_id: store.currentProjectId,
      title: newWikiTitle.value.trim(),
      slug: newWikiTitle.value.trim(),
      content: newWikiContent.value.trim(),
      author_name: 'Agus Setiawan',
      client_visible: newWikiClientVisible.value,
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString()
    })
    store.showToast('Dokumen wiki berhasil ditambahkan!')
    showNewWikiModal.value = false
    newWikiTitle.value = ''
    newWikiContent.value = ''
  }
}

onMounted(() => {
  if (route.params.id) {
    store.selectProject(route.params.id as string)
  }
  if (route.query.tab === 'list' || route.query.tab === 'timeline' || route.query.tab === 'kanban') {
    store.setView(route.query.tab as 'kanban' | 'list' | 'timeline')
  }
  loadWikiPages()
})

watch(() => route.params.id, (newId) => {
  if (newId) {
    store.selectProject(newId as string)
  }
})

watch(() => route.query.tab, (newTab) => {
  if (newTab === 'list' || newTab === 'timeline' || newTab === 'kanban') {
    store.setView(newTab as 'kanban' | 'list' | 'timeline')
  }
})

watch(() => store.currentProjectId, () => {
  loadWikiPages()
})
</script>
