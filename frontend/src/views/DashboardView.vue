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
        <p class="text-xs text-slate-500 font-medium">Ringkasan operasional milestone engineering, deliverable tim, dan manajemen sistem.</p>
      </div>

      <!-- Segmented Navigation Pills: Overview vs User Management (Admin) -->
      <div class="flex items-center gap-2 flex-wrap">
        <div v-if="canManageUsers" class="p-1 rounded-xl bg-slate-100/90 border border-slate-200/80 flex items-center gap-1 text-xs font-semibold">
          <button
            @click="switchTab('overview')"
            :class="[
              'px-3 py-1.5 rounded-lg flex items-center gap-1.5 transition-all duration-200',
              activeDashboardTab === 'overview'
                ? 'bg-white text-slate-900 shadow-sm font-bold'
                : 'text-slate-600 hover:text-slate-900'
            ]"
          >
            <LayoutDashboard class="w-3.5 h-3.5 text-brand-600" />
            <span>Ringkasan Portfolio</span>
          </button>

          <button
            @click="switchTab('users')"
            :class="[
              'px-3 py-1.5 rounded-lg flex items-center gap-1.5 transition-all duration-200',
              activeDashboardTab === 'users'
                ? 'bg-white text-brand-700 shadow-sm font-bold'
                : 'text-slate-600 hover:text-slate-900'
            ]"
          >
            <Users class="w-3.5 h-3.5 text-indigo-600" />
            <span>Manajemen Pengguna (CRUD)</span>
            <span class="px-1.5 py-0.2 rounded-full text-[9px] font-mono bg-indigo-50 text-indigo-700 border border-indigo-200 ml-1">
              {{ userList.length }}
            </span>
          </button>
        </div>

        <!-- Actions for Overview: Import & Export Module -->
        <div v-if="activeDashboardTab === 'overview'" class="flex items-center gap-2">
          <!-- Button Impor Proyek -->
          <button
            v-if="canCreateProject"
            @click="openImportModal"
            class="px-3 py-1.5 rounded-xl bg-white hover:bg-slate-50 border border-slate-200 text-slate-700 font-semibold transition-all shadow-2xs hover:border-brand-300 text-xs flex items-center gap-1.5 cursor-pointer"
            title="Impor Data Proyek dari MS Excel atau CSV"
          >
            <Upload class="w-3.5 h-3.5 text-brand-600" />
            <span>Impor Proyek</span>
          </button>

          <!-- Dropdown Ekspor Data -->
          <div class="relative">
            <button
              @click="showExportDropdown = !showExportDropdown"
              class="px-3.5 py-1.5 rounded-xl bg-gradient-to-r from-brand-600 to-indigo-600 hover:from-brand-500 hover:to-indigo-500 text-white font-bold transition-all shadow-md shadow-brand-600/20 active:scale-95 text-xs flex items-center gap-1.5 cursor-pointer"
            >
              <Download class="w-3.5 h-3.5" />
              <span>Ekspor (Excel/CSV)</span>
              <ChevronDown class="w-3 h-3 ml-0.5 opacity-80" />
            </button>

            <!-- Backdrop to close dropdown -->
            <div
              v-if="showExportDropdown"
              class="fixed inset-0 z-30"
              @click="showExportDropdown = false"
            ></div>

            <!-- Dropdown Menu -->
            <div
              v-if="showExportDropdown"
              class="absolute right-0 mt-1.5 w-60 bg-white rounded-2xl border border-slate-200 shadow-xl py-1.5 z-40 text-xs animate-in fade-in zoom-in-95 duration-100"
            >
              <div class="px-3 py-1 text-[10px] font-mono uppercase tracking-wider text-slate-400 font-bold border-b border-slate-100 mb-1">
                Opsi Ekspor MS Excel & CSV
              </div>
              <button
                @click="handleExportProjects"
                :disabled="isExportingProjects"
                class="w-full text-left px-3.5 py-2 hover:bg-slate-50 text-slate-700 flex items-center gap-2.5 transition-colors cursor-pointer"
              >
                <div class="p-1 rounded-md bg-emerald-50 text-emerald-600">
                  <FileSpreadsheet class="w-3.5 h-3.5" />
                </div>
                <div>
                  <div class="font-semibold text-slate-800">Ekspor Proyek (.csv)</div>
                  <div class="text-[10px] text-slate-400 font-normal">Kompatibel langsung dengan MS Excel</div>
                </div>
              </button>

              <button
                @click="handleExportTasks"
                :disabled="isExportingTasks"
                class="w-full text-left px-3.5 py-2 hover:bg-slate-50 text-slate-700 flex items-center gap-2.5 transition-colors cursor-pointer"
              >
                <div class="p-1 rounded-md bg-indigo-50 text-indigo-600">
                  <ListTodo class="w-3.5 h-3.5" />
                </div>
                <div>
                  <div class="font-semibold text-slate-800">Ekspor Tasks / Deliverable (.csv)</div>
                  <div class="text-[10px] text-slate-400 font-normal">Daftar tugas, assignee, & status</div>
                </div>
              </button>

              <div class="border-t border-slate-100 my-1"></div>

              <button
                @click="handleDownloadTemplate"
                class="w-full text-left px-3.5 py-2 hover:bg-slate-50 text-slate-600 flex items-center gap-2.5 transition-colors cursor-pointer"
              >
                <div class="p-1 rounded-md bg-slate-100 text-slate-600">
                  <Download class="w-3.5 h-3.5" />
                </div>
                <div>
                  <div class="font-medium text-slate-700">Unduh Template Impor</div>
                  <div class="text-[10px] text-slate-400 font-normal">Contoh format header & data CSV</div>
                </div>
              </button>
            </div>
          </div>
        </div>

        <div v-else-if="canManageUsers" class="flex items-center gap-2">
          <!-- Button Undang Klien / Pengguna via Email -->
          <button
            @click="openInviteModal"
            class="px-3.5 py-1.5 rounded-xl bg-white hover:bg-indigo-50 border border-indigo-200 text-indigo-700 font-bold transition-all shadow-2xs hover:border-indigo-300 text-xs flex items-center gap-1.5 cursor-pointer"
            title="Kirim Undangan Pendaftaran ke Klien / Pengguna via Email"
          >
            <Mail class="w-3.5 h-3.5 text-indigo-600" />
            <span>Undang Klien / User</span>
          </button>

          <button
            @click="openCreateUserModal"
            class="px-4 py-1.5 rounded-xl bg-gradient-to-r from-brand-600 to-indigo-600 hover:from-brand-500 hover:to-indigo-500 text-white font-bold transition-all shadow-md shadow-brand-600/20 active:scale-95 text-xs flex items-center gap-1.5 cursor-pointer"
          >
            <UserPlus class="w-3.5 h-3.5" />
            <span>Tambah Pengguna</span>
          </button>
        </div>
      </div>
    </div>

    <!-- ==================== TAB 1: EXECUTIVE OVERVIEW ==================== -->
    <template v-if="activeDashboardTab === 'overview'">
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

              <!-- Row 3: Metadata Footer & Actions -->
              <div class="mt-2.5 flex flex-col sm:flex-row sm:items-center justify-between text-[11px] text-slate-500 font-mono gap-1.5 pt-1.5 border-t border-slate-100">
                <span class="tabular-nums font-medium">{{ proj.completed_tasks }} dari {{ proj.total_tasks }} deliverable selesai</span>
                <div class="flex items-center gap-1.5">
                  <span class="text-slate-400 mr-1 hidden sm:inline">Tenggat: {{ proj.deadline }}</span>
                  
                  <!-- Edit Proyek Button (Admin / PM) -->
                  <button
                    v-if="canManageProjectAccess"
                    type="button"
                    @click.stop="openEditProjectModal(proj)"
                    class="px-2 py-0.5 rounded-lg bg-white hover:bg-slate-100 border border-slate-200 text-slate-700 text-[10px] font-bold flex items-center gap-1 shadow-2xs transition-colors cursor-pointer"
                    title="Edit nama, klien, deadline, dan status proyek"
                  >
                    <Edit class="w-3 h-3 text-slate-600" />
                    <span>Edit</span>
                  </button>

                  <!-- Kelola Akses Tim Button -->
                  <button
                    v-if="canManageProjectAccess"
                    type="button"
                    @click.stop="openProjectAccessModal(proj)"
                    class="px-2 py-0.5 rounded-lg bg-white hover:bg-brand-50 border border-slate-200 hover:border-brand-300 text-brand-700 text-[10px] font-bold flex items-center gap-1 shadow-2xs transition-colors cursor-pointer"
                    title="Kelola hak akses team member ke proyek ini"
                  >
                    <ShieldCheck class="w-3 h-3 text-brand-600" />
                    <span>Akses</span>
                  </button>

                  <!-- Salin / Duplikasi Proyek Button (Admin / PM) -->
                  <button
                    v-if="canManageProjectAccess"
                    type="button"
                    @click.stop="openCopyProjectModal(proj)"
                    class="px-2 py-0.5 rounded-lg bg-white hover:bg-indigo-50 border border-slate-200 hover:border-indigo-300 text-indigo-700 text-[10px] font-bold flex items-center gap-1 shadow-2xs transition-colors cursor-pointer"
                    title="Duplikasi seluruh pekerjaan proyek ini ke proyek baru"
                  >
                    <Copy class="w-3 h-3 text-indigo-600" />
                    <span>Salin</span>
                  </button>

                  <!-- Hapus Proyek Button (Admin / PM) -->
                  <button
                    v-if="canManageProjectAccess"
                    type="button"
                    @click.stop="handleDeleteProject(proj)"
                    class="p-1 rounded-lg bg-white hover:bg-rose-50 border border-slate-200 hover:border-rose-300 text-rose-600 text-[10px] font-bold flex items-center shadow-2xs transition-colors cursor-pointer"
                    title="Hapus proyek ini"
                  >
                    <Trash2 class="w-3 h-3 text-rose-500" />
                  </button>
                </div>
              </div>
            </div>

            <!-- Empty state when user has no projects assigned -->
            <div
              v-if="store.projects.length === 0"
              class="py-10 px-4 text-center rounded-2xl border-2 border-dashed border-slate-200 bg-slate-50/50 space-y-3"
            >
              <div class="w-12 h-12 rounded-2xl bg-indigo-50 border border-indigo-100 text-indigo-600 flex items-center justify-center mx-auto shadow-2xs">
                <FolderKanban class="w-6 h-6" />
              </div>
              <div class="space-y-1 max-w-sm mx-auto">
                <h4 class="text-xs sm:text-sm font-bold text-slate-800">Belum Ada Proyek yang Ditugaskan</h4>
                <p class="text-xs text-slate-500 leading-relaxed">
                  Akun Anda belum didaftarkan sebagai anggota di proyek mana pun oleh Admin/PM. Silakan hubungi Administrator untuk memberikan akses proyek melalui tombol <strong>Akses Tim</strong>.
                </p>
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
    </template>

    <!-- ==================== TAB 2: USER MANAGEMENT (ADMIN CRUD) ==================== -->
    <template v-else-if="activeDashboardTab === 'users' && canManageUsers">
      <div class="space-y-5">
        <!-- Role Statistics Cards -->
        <div class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-5 gap-3">
          <div class="p-3.5 rounded-2xl bg-white/95 border border-slate-200/90 shadow-sm flex flex-col justify-between">
            <span class="text-[11px] font-mono text-slate-500 font-semibold uppercase">Total Pengguna</span>
            <div class="flex items-baseline justify-between mt-1">
              <span class="text-2xl font-extrabold text-slate-900 font-display">{{ userList.length }}</span>
              <span class="text-[10px] font-mono text-slate-400">Terdaftar</span>
            </div>
          </div>

          <div class="p-3.5 rounded-2xl bg-white/95 border border-slate-200/90 shadow-sm flex flex-col justify-between">
            <span class="text-[11px] font-mono text-rose-600 font-semibold uppercase">Super Admin</span>
            <div class="flex items-baseline justify-between mt-1">
              <span class="text-2xl font-extrabold text-rose-600 font-display">
                {{ userList.filter(u => u.role === 'super_admin').length }}
              </span>
              <span class="text-[10px] font-mono text-rose-400">Full Akses</span>
            </div>
          </div>

          <div class="p-3.5 rounded-2xl bg-white/95 border border-slate-200/90 shadow-sm flex flex-col justify-between">
            <span class="text-[11px] font-mono text-emerald-600 font-semibold uppercase">Project Mgr (PM)</span>
            <div class="flex items-baseline justify-between mt-1">
              <span class="text-2xl font-extrabold text-emerald-600 font-display">
                {{ userList.filter(u => u.role === 'pm').length }}
              </span>
              <span class="text-[10px] font-mono text-emerald-400">Lead Tim</span>
            </div>
          </div>

          <div class="p-3.5 rounded-2xl bg-white/95 border border-slate-200/90 shadow-sm flex flex-col justify-between">
            <span class="text-[11px] font-mono text-indigo-600 font-semibold uppercase">Team Member</span>
            <div class="flex items-baseline justify-between mt-1">
              <span class="text-2xl font-extrabold text-indigo-600 font-display">
                {{ userList.filter(u => u.role === 'member').length }}
              </span>
              <span class="text-[10px] font-mono text-indigo-400">Engineer</span>
            </div>
          </div>

          <div class="p-3.5 rounded-2xl bg-white/95 border border-slate-200/90 shadow-sm flex flex-col justify-between col-span-2 sm:col-span-1">
            <span class="text-[11px] font-mono text-amber-600 font-semibold uppercase">Client Portal</span>
            <div class="flex items-baseline justify-between mt-1">
              <span class="text-2xl font-extrabold text-amber-600 font-display">
                {{ userList.filter(u => u.role === 'client').length }}
              </span>
              <span class="text-[10px] font-mono text-amber-400">Reviewer</span>
            </div>
          </div>
        </div>

        <!-- Sub-Navigation Switcher: Pengguna Aktif vs Undangan Tertunda -->
        <div class="flex items-center justify-between gap-2 flex-wrap pt-1">
          <div class="p-1 rounded-xl bg-slate-100/90 border border-slate-200/80 flex items-center gap-1 text-xs font-semibold">
            <button
              @click="userSubTab = 'active'"
              :class="[
                'px-3.5 py-1.5 rounded-lg flex items-center gap-1.5 transition-all cursor-pointer',
                userSubTab === 'active'
                  ? 'bg-white text-slate-900 shadow-2xs font-bold'
                  : 'text-slate-600 hover:text-slate-900'
              ]"
            >
              <Users class="w-3.5 h-3.5 text-brand-600" />
              <span>Daftar Pengguna Aktif</span>
              <span class="px-1.5 py-0.2 rounded-full text-[9px] font-mono bg-slate-100 text-slate-700 border border-slate-200 ml-1">
                {{ userList.length }}
              </span>
            </button>

            <button
              @click="userSubTab = 'invitations'"
              :class="[
                'px-3.5 py-1.5 rounded-lg flex items-center gap-1.5 transition-all cursor-pointer',
                userSubTab === 'invitations'
                  ? 'bg-white text-indigo-700 shadow-2xs font-bold'
                  : 'text-slate-600 hover:text-slate-900'
              ]"
            >
              <Mail class="w-3.5 h-3.5 text-indigo-600" />
              <span>Undangan Klien & User</span>
              <span
                v-if="pendingInvitationsCount > 0"
                class="px-1.5 py-0.2 rounded-full text-[9px] font-mono bg-indigo-50 text-indigo-700 font-bold border border-indigo-200 ml-1"
              >
                {{ pendingInvitationsCount }}
              </span>
            </button>
          </div>

          <div class="flex items-center gap-2">
            <button
              v-if="userSubTab === 'invitations'"
              @click="openInviteModal"
              class="px-3.5 py-1.5 rounded-xl bg-indigo-600 hover:bg-indigo-500 text-white font-bold text-xs shadow-md shadow-indigo-600/20 flex items-center gap-1.5 transition-all cursor-pointer"
            >
              <MailPlus class="w-3.5 h-3.5" />
              <span>Undang Klien Baru</span>
            </button>
          </div>
        </div>

        <!-- ==================== SUB-VIEW 1: ACTIVE USERS (CRUD) ==================== -->
        <template v-if="userSubTab === 'active'">
        <!-- Filter & Search Controls -->
        <div class="p-4 rounded-2xl bg-white/95 backdrop-blur-xl border border-slate-200/90 shadow-lg shadow-slate-200/30 flex flex-col sm:flex-row items-stretch sm:items-center justify-between gap-3">
          <!-- Search Input -->
          <div class="relative flex-1 max-w-md">
            <Search class="w-4 h-4 text-slate-400 absolute left-3 top-1/2 -translate-y-1/2" />
            <input
              type="text"
              v-model="userSearchQuery"
              placeholder="Cari berdasarkan nama, email, atau ID pengguna..."
              class="w-full pl-9 pr-3 py-2 rounded-xl border border-slate-200 bg-slate-50/70 focus:bg-white text-xs text-slate-900 placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-brand-500 transition-all shadow-xs"
            />
          </div>

          <!-- Role Filter Dropdown & Refresh -->
          <div class="flex items-center gap-2">
            <select
              v-model="userRoleFilter"
              class="py-2 px-3 rounded-xl border border-slate-200 bg-slate-50/70 text-xs text-slate-700 font-medium focus:bg-white focus:outline-none focus:ring-2 focus:ring-brand-500 shadow-xs cursor-pointer"
            >
              <option value="all">Semua Peran (All Roles)</option>
              <option value="super_admin">Super Admin</option>
              <option value="pm">Project Manager (PM)</option>
              <option value="member">Team Member</option>
              <option value="client">Client Portal</option>
            </select>

            <button
              @click="loadUsers"
              :disabled="isLoadingUsers"
              class="p-2 rounded-xl border border-slate-200 hover:bg-slate-50 text-slate-600 transition-colors shadow-xs cursor-pointer"
              title="Segarkan Data Pengguna"
            >
              <RefreshCw class="w-4 h-4" :class="isLoadingUsers ? 'animate-spin text-brand-600' : ''" />
            </button>
          </div>
        </div>

        <!-- Main User Table (Desktop & Tablet) -->
        <div class="hidden md:block rounded-2xl bg-white/95 backdrop-blur-xl border border-slate-200/90 shadow-xl shadow-slate-200/50 overflow-hidden">
          <div class="overflow-x-auto">
            <table class="w-full text-left text-xs">
              <thead class="bg-slate-50/80 border-b border-slate-200/80 text-[10px] font-mono uppercase text-slate-500 font-bold tracking-wider">
                <tr>
                  <th class="py-3 px-4">Pengguna</th>
                  <th class="py-3 px-4">Peran (Role)</th>
                  <th class="py-3 px-4">Akses Proyek</th>
                  <th class="py-3 px-4">Status Sesi / Device</th>
                  <th class="py-3 px-4">Terakhir Aktif</th>
                  <th class="py-3 px-4 text-right">Aksi</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-slate-100">
                <tr
                  v-for="u in filteredUsers"
                  :key="u.id"
                  class="hover:bg-slate-50/80 transition-colors group"
                >
                  <!-- User Info Column -->
                  <td class="py-3 px-4">
                    <div class="flex items-center gap-3">
                      <div class="relative flex-shrink-0">
                        <img
                          :src="u.avatar_url || 'https://images.unsplash.com/photo-1534528741775-53994a69daeb?w=100&auto=format&fit=crop&q=80'"
                          :alt="u.name"
                          class="w-9 h-9 rounded-xl object-cover ring-1 ring-slate-200"
                        />
                        <span
                          v-if="u.current_session_id"
                          class="absolute -bottom-0.5 -right-0.5 w-2.5 h-2.5 rounded-full bg-emerald-500 ring-2 ring-white"
                          title="Sesi Aktif"
                        ></span>
                      </div>
                      <div class="min-w-0">
                        <div class="font-bold text-slate-900 flex items-center gap-1.5">
                          <span>{{ u.name }}</span>
                          <span
                            v-if="u.id === authStore.user?.id"
                            class="text-[9px] font-mono px-1.5 py-0.2 rounded-full bg-brand-50 text-brand-700 font-bold border border-brand-200"
                          >
                            Anda
                          </span>
                        </div>
                        <p class="text-[11px] text-slate-500 truncate">{{ u.email }}</p>
                        <p class="text-[9px] font-mono text-slate-400 mt-0.5">{{ u.id }}</p>
                      </div>
                    </div>
                  </td>

                  <!-- Role Column -->
                  <td class="py-3 px-4">
                    <span :class="['inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-[10px] font-mono font-bold uppercase border', getRoleBadgeClass(u.role)]">
                      <Shield class="w-3 h-3" />
                      <span>{{ formatRoleName(u.role) }}</span>
                    </span>
                  </td>

                  <!-- Project Access Column -->
                  <td class="py-3 px-4">
                    <div v-if="u.role === 'super_admin'" class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full bg-rose-50 border border-rose-200 text-rose-700 text-[10px] font-mono font-bold">
                      <ShieldCheck class="w-3 h-3" />
                      <span>Semua Proyek (Global)</span>
                    </div>
                    <div v-else-if="u.project_ids && u.project_ids.length > 0" class="space-y-1">
                      <div class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full bg-indigo-50 border border-indigo-200 text-indigo-700 text-[10px] font-mono font-semibold">
                        <FolderKanban class="w-3 h-3 text-indigo-600" />
                        <span>{{ u.project_ids.length }} Proyek</span>
                      </div>
                      <div class="text-[10px] text-slate-500 truncate max-w-[200px]" :title="getProjectNames(u.project_ids).join(', ')">
                        {{ getProjectNames(u.project_ids).join(', ') }}
                      </div>
                    </div>
                    <div v-else class="inline-flex items-center gap-1 px-2.5 py-1 rounded-full bg-amber-50 border border-amber-200 text-amber-700 text-[10px] font-mono font-medium">
                      <AlertCircle class="w-3 h-3" />
                      <span>0 Proyek (Belum ada)</span>
                    </div>
                  </td>

                  <!-- Session Status Column -->
                  <td class="py-3 px-4">
                    <div v-if="u.current_session_id" class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full bg-emerald-50 border border-emerald-200 text-emerald-700 text-[11px] font-medium">
                      <span class="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"></span>
                      <span>Sesi Terhubung</span>
                    </div>
                    <div v-else class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full bg-slate-100 border border-slate-200 text-slate-500 text-[11px] font-medium">
                      <span class="w-2 h-2 rounded-full bg-slate-300"></span>
                      <span>Offline / Tanpa Sesi</span>
                    </div>
                  </td>

                  <!-- Last Active Column -->
                  <td class="py-3 px-4 text-slate-500 font-mono text-[11px]">
                    {{ formatLastActive(u.last_active_at) }}
                  </td>

                  <!-- Actions Column -->
                  <td class="py-3 px-4 text-right">
                    <div class="inline-flex items-center gap-1.5">
                      <!-- Edit User -->
                      <button
                        @click="openEditUserModal(u)"
                        class="p-1.5 rounded-lg border border-slate-200 hover:border-brand-300 hover:bg-brand-50 text-slate-600 hover:text-brand-600 transition-colors shadow-2xs cursor-pointer"
                        title="Edit Data Pengguna & Kata Sandi"
                      >
                        <Edit class="w-3.5 h-3.5" />
                      </button>

                      <!-- Force Logout / Reset Session -->
                      <button
                        v-if="u.current_session_id"
                        @click="handleResetSession(u)"
                        class="p-1.5 rounded-lg border border-amber-200 hover:bg-amber-50 text-amber-600 transition-colors shadow-2xs cursor-pointer"
                        title="Putus Sesi / Lepas Perangkat Pengguna"
                      >
                        <ZapOff class="w-3.5 h-3.5" />
                      </button>

                      <!-- Delete User -->
                      <button
                        @click="openDeleteUserModal(u)"
                        :disabled="u.id === 'usr_super_admin' || u.email === 'admin@pulse.workspace' || u.id === authStore.user?.id"
                        :class="[
                          'p-1.5 rounded-lg border transition-colors shadow-2xs',
                          u.id === 'usr_super_admin' || u.email === 'admin@pulse.workspace' || u.id === authStore.user?.id
                            ? 'border-slate-200 text-slate-300 cursor-not-allowed opacity-50'
                            : 'border-slate-200 hover:border-rose-300 hover:bg-rose-50 text-slate-600 hover:text-rose-600 cursor-pointer'
                        ]"
                        :title="u.id === 'usr_super_admin' ? 'Akun Super Admin dilindungi' : (u.id === authStore.user?.id ? 'Tidak dapat menghapus akun sendiri' : 'Hapus Pengguna')"
                      >
                        <Trash2 class="w-3.5 h-3.5" />
                      </button>
                    </div>
                  </td>
                </tr>

                <tr v-if="filteredUsers.length === 0">
                  <td colspan="6" class="py-8 text-center text-slate-400">
                    Tidak ada pengguna yang cocok dengan kriteria pencarian.
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- Mobile User Card List -->
        <div class="md:hidden space-y-3">
          <div
            v-for="u in filteredUsers"
            :key="u.id"
            class="p-4 rounded-2xl bg-white border border-slate-200/90 shadow-md space-y-3"
          >
            <div class="flex items-start justify-between gap-3">
              <div class="flex items-center gap-3">
                <div class="relative flex-shrink-0">
                  <img
                    :src="u.avatar_url || 'https://images.unsplash.com/photo-1534528741775-53994a69daeb?w=100&auto=format&fit=crop&q=80'"
                    :alt="u.name"
                    class="w-10 h-10 rounded-xl object-cover ring-1 ring-slate-200"
                  />
                  <span
                    v-if="u.current_session_id"
                    class="absolute -bottom-0.5 -right-0.5 w-2.5 h-2.5 rounded-full bg-emerald-500 ring-2 ring-white"
                  ></span>
                </div>
                <div>
                  <div class="font-bold text-slate-900 text-xs">{{ u.name }}</div>
                  <div class="text-[11px] text-slate-500">{{ u.email }}</div>
                </div>
              </div>

              <span :class="['px-2 py-0.5 rounded-full text-[9px] font-mono font-bold uppercase border', getRoleBadgeClass(u.role)]">
                {{ formatRoleName(u.role) }}
              </span>
            </div>

            <!-- Project Access Badges for Mobile -->
            <div class="pt-2 border-t border-slate-100 flex items-center justify-between text-xs">
              <span class="text-slate-400 font-mono text-[10px]">Akses Proyek:</span>
              <div v-if="u.role === 'super_admin'" class="text-[10px] font-bold text-rose-700 font-mono">
                Semua Proyek (Global)
              </div>
              <div v-else-if="u.project_ids && u.project_ids.length > 0" class="text-[10px] font-semibold text-indigo-700 flex items-center gap-1 font-mono">
                <FolderKanban class="w-3 h-3 text-indigo-600" />
                <span>{{ u.project_ids.length }} Proyek</span>
              </div>
              <div v-else class="text-[10px] text-amber-600 font-medium flex items-center gap-1 font-mono">
                <AlertCircle class="w-3 h-3 text-amber-500" />
                <span>0 Proyek (Belum ada)</span>
              </div>
            </div>

            <div class="pt-2 border-t border-slate-100 flex items-center justify-between text-[11px] font-mono text-slate-500">
              <span v-if="u.current_session_id" class="text-emerald-700 font-semibold flex items-center gap-1">
                <span class="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse"></span>
                Sesi Aktif
              </span>
              <span v-else class="text-slate-400">Offline</span>
              <span>{{ formatLastActive(u.last_active_at) }}</span>
            </div>

            <!-- Mobile Action Buttons -->
            <div class="flex items-center gap-2 pt-1">
              <button
                @click="openEditUserModal(u)"
                class="flex-1 py-1.5 px-2 rounded-xl bg-slate-50 hover:bg-slate-100 border border-slate-200 text-slate-700 text-xs font-semibold flex items-center justify-center gap-1 cursor-pointer"
              >
                <Edit class="w-3.5 h-3.5" />
                <span>Edit</span>
              </button>

              <button
                v-if="u.current_session_id"
                @click="handleResetSession(u)"
                class="py-1.5 px-3 rounded-xl bg-amber-50 hover:bg-amber-100 border border-amber-200 text-amber-700 text-xs font-semibold flex items-center justify-center gap-1 cursor-pointer"
                title="Lepas Sesi"
              >
                <ZapOff class="w-3.5 h-3.5" />
                <span>Reset Sesi</span>
              </button>

              <button
                @click="openDeleteUserModal(u)"
                :disabled="u.id === 'usr_super_admin' || u.email === 'admin@pulse.workspace' || u.id === authStore.user?.id"
                :class="[
                  'py-1.5 px-3 rounded-xl border text-xs font-semibold flex items-center justify-center gap-1',
                  u.id === 'usr_super_admin' || u.email === 'admin@pulse.workspace' || u.id === authStore.user?.id
                    ? 'border-slate-200 text-slate-300 opacity-50 cursor-not-allowed'
                    : 'bg-rose-50 hover:bg-rose-100 border-rose-200 text-rose-700 cursor-pointer'
                ]"
              >
                <Trash2 class="w-3.5 h-3.5" />
                <span>Hapus</span>
              </button>
            </div>
          </div>
        </div>
        </template>

        <!-- ==================== SUB-VIEW 2: INVITATIONS LIST ==================== -->
        <template v-else-if="userSubTab === 'invitations'">
          <div class="space-y-4">
            <!-- Header bar for Invitations -->
            <div class="p-4 rounded-2xl bg-white/95 backdrop-blur-xl border border-slate-200/90 shadow-lg shadow-slate-200/30 flex flex-col sm:flex-row items-stretch sm:items-center justify-between gap-3">
              <div>
                <h3 class="font-bold text-slate-900 text-sm flex items-center gap-2">
                  <Mail class="w-4 h-4 text-indigo-600" />
                  <span>Daftar Undangan Klien & Pengguna</span>
                </h3>
                <p class="text-[11px] text-slate-500 mt-0.5">
                  Tautan undangan yang dapat Anda kirim via WhatsApp atau Email. Saat calon pengguna mengonfirmasi, akun mereka langsung dibuat dengan email dan hak akses proyek yang telah Anda tetapkan.
                </p>

                <!-- SMTP Service Status Pill -->
                <div class="mt-2.5 flex flex-wrap items-center gap-2">
                  <div
                    class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-[11px] font-semibold border"
                    :class="smtpStatus.configured ? 'bg-emerald-50 text-emerald-800 border-emerald-200' : 'bg-amber-50 text-amber-800 border-amber-200'"
                  >
                    <span class="w-2 h-2 rounded-full" :class="smtpStatus.configured ? 'bg-emerald-500 animate-pulse' : 'bg-amber-500'"></span>
                    <span>{{ smtpStatus.configured ? `SMTP Aktif (${smtpStatus.smtp_host} • ${smtpStatus.smtp_user})` : 'SMTP Belum Dikonfigurasi (Email Otomatis Nonaktif)' }}</span>
                  </div>
                  <button
                    v-if="canManageUsers"
                    @click="openSMTPConfigModal"
                    class="text-[11px] text-indigo-600 hover:text-indigo-800 font-bold underline flex items-center gap-1 cursor-pointer"
                  >
                    <Settings class="w-3 h-3" />
                    <span>Pengaturan & Tes SMTP</span>
                  </button>
                </div>
              </div>

              <div class="flex items-center gap-2 flex-shrink-0">
                <button
                  @click="openInviteModal"
                  class="px-3.5 py-2 rounded-xl bg-gradient-to-r from-brand-600 to-indigo-600 hover:from-brand-500 hover:to-indigo-500 text-white font-bold text-xs shadow-md shadow-brand-600/20 flex items-center gap-1.5 transition-all cursor-pointer"
                >
                  <MailPlus class="w-3.5 h-3.5" />
                  <span>Undang Klien Baru</span>
                </button>

                <button
                  @click="loadInvitations"
                  :disabled="isLoadingInvitations"
                  class="p-2 rounded-xl border border-slate-200 hover:bg-slate-50 text-slate-600 transition-colors shadow-xs cursor-pointer"
                  title="Segarkan Undangan"
                >
                  <RefreshCw class="w-4 h-4" :class="isLoadingInvitations ? 'animate-spin text-brand-600' : ''" />
                </button>
              </div>
            </div>

            <!-- Invitations Desktop Table -->
            <div class="hidden md:block rounded-2xl bg-white/95 backdrop-blur-xl border border-slate-200/90 shadow-xl shadow-slate-200/50 overflow-hidden">
              <div class="overflow-x-auto">
                <table class="w-full text-left text-xs">
                  <thead class="bg-slate-50/80 border-b border-slate-200/80 text-[10px] font-mono uppercase text-slate-500 font-bold tracking-wider">
                    <tr>
                      <th class="py-3 px-4">Calon Pengguna (Email)</th>
                      <th class="py-3 px-4">Peran / Hak Akses</th>
                      <th class="py-3 px-4">Akses Proyek Ditugaskan</th>
                      <th class="py-3 px-4">Status Undangan</th>
                      <th class="py-3 px-4">Tenggat Waktu</th>
                      <th class="py-3 px-4 text-right">Aksi</th>
                    </tr>
                  </thead>
                  <tbody class="divide-y divide-slate-100">
                    <tr v-if="invitationsList.length === 0">
                      <td colspan="6" class="py-12 text-center text-slate-400">
                        <Mail class="w-8 h-8 text-slate-300 mx-auto mb-2" />
                        <p class="font-medium text-xs text-slate-600">Belum ada undangan yang dibuat</p>
                        <p class="text-[11px] text-slate-400 mt-1">Undang klien atau rekan tim baru untuk bergabung ke proyek</p>
                        <button
                          @click="openInviteModal"
                          class="mt-3 px-3.5 py-1.5 rounded-xl bg-indigo-600 hover:bg-indigo-500 text-white font-bold text-xs shadow-md shadow-indigo-600/20 inline-flex items-center gap-1.5 cursor-pointer"
                        >
                          <MailPlus class="w-3.5 h-3.5" />
                          <span>Undang Klien Sekarang</span>
                        </button>
                      </td>
                    </tr>

                    <tr
                      v-for="inv in invitationsList"
                      :key="inv.id"
                      class="hover:bg-slate-50/80 transition-colors"
                    >
                      <!-- Email & Name -->
                      <td class="py-3 px-4">
                        <div class="font-bold text-slate-900 select-all">
                          {{ inv.email }}
                        </div>
                        <p v-if="inv.name" class="text-[11px] text-slate-600 font-medium">{{ inv.name }}</p>
                        <p class="text-[9px] font-mono text-slate-400 mt-0.5">Pengundang: {{ inv.inviter_name || 'Administrator' }}</p>
                      </td>

                      <!-- Role Badge -->
                      <td class="py-3 px-4">
                        <span
                          :class="[
                            'px-2.5 py-1 rounded-full text-[10px] font-mono font-bold uppercase border shadow-2xs inline-block',
                            inv.role === 'client' ? 'bg-amber-50 text-amber-800 border-amber-200' :
                            inv.role === 'pm' ? 'bg-purple-50 text-purple-800 border-purple-200' :
                            'bg-blue-50 text-blue-800 border-blue-200'
                          ]"
                        >
                          {{ formatRoleName(inv.role) }}
                        </span>
                      </td>

                      <!-- Project Badges -->
                      <td class="py-3 px-4">
                        <div v-if="inv.project_ids && inv.project_ids.length > 0" class="flex flex-wrap gap-1 max-w-xs">
                          <span
                            v-for="pid in inv.project_ids"
                            :key="pid"
                            class="inline-flex items-center gap-1 px-2 py-0.5 rounded-md bg-slate-100 text-slate-700 text-[10px] font-medium"
                          >
                            <span class="w-1.5 h-1.5 rounded-full" :style="{ backgroundColor: getProjectColor(pid) }"></span>
                            <span class="truncate max-w-[120px]">{{ getProjectName(pid) }}</span>
                          </span>
                        </div>
                        <span v-else class="text-[10px] text-slate-400 italic">Belum Ada Proyek Ditugaskan</span>
                      </td>

                      <!-- Status -->
                      <td class="py-3 px-4">
                        <span
                          v-if="inv.status === 'pending'"
                          class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full bg-amber-50 border border-amber-200 text-amber-800 text-[10px] font-semibold"
                        >
                          <span class="w-1.5 h-1.5 rounded-full bg-amber-500 animate-pulse"></span>
                          <span>Menunggu Konfirmasi</span>
                        </span>
                        <span
                          v-else-if="inv.status === 'accepted'"
                          class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full bg-emerald-50 border border-emerald-200 text-emerald-800 text-[10px] font-semibold"
                        >
                          <CheckCircle2 class="w-3 h-3 text-emerald-600" />
                          <span>Telah Dikonfirmasi</span>
                        </span>
                        <span
                          v-else-if="inv.status === 'expired'"
                          class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full bg-slate-100 border border-slate-200 text-slate-500 text-[10px] font-semibold"
                        >
                          <Clock class="w-3 h-3 text-slate-400" />
                          <span>Kedaluwarsa</span>
                        </span>
                        <span
                          v-else
                          class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full bg-rose-50 border border-rose-200 text-rose-700 text-[10px] font-semibold"
                        >
                          <X class="w-3.5 h-3.5 text-rose-500" />
                          <span>Dibatalkan</span>
                        </span>
                      </td>

                      <!-- Expiry Date -->
                      <td class="py-3 px-4 text-slate-500 font-mono text-[11px]">
                        {{ formatExpiryDate(inv.expires_at) }}
                      </td>

                      <!-- Actions -->
                      <td class="py-3 px-4 text-right">
                        <div class="flex items-center justify-end gap-1.5">
                          <!-- Resend Email via SMTP -->
                          <button
                            v-if="inv.status === 'pending'"
                            @click="resendInvitationEmail(inv)"
                            :disabled="isResendingId === inv.id"
                            class="py-1 px-2.5 rounded-lg border border-indigo-200 hover:bg-indigo-50 text-indigo-700 text-[11px] font-semibold flex items-center gap-1 transition-colors cursor-pointer disabled:opacity-50"
                            title="Kirim ulang email undangan langsung ke kotak masuk pengguna"
                          >
                            <Loader2 v-if="isResendingId === inv.id" class="w-3.5 h-3.5 animate-spin text-indigo-600" />
                            <Send v-else class="w-3.5 h-3.5 text-indigo-600" />
                            <span>{{ isResendingId === inv.id ? 'Mengirim...' : 'Kirim Email' }}</span>
                          </button>

                          <!-- Copy Link -->
                          <button
                            v-if="inv.status === 'pending'"
                            @click="copyInvitationLink(inv)"
                            class="py-1 px-2.5 rounded-lg border border-slate-200 hover:bg-slate-50 text-slate-700 text-[11px] font-semibold flex items-center gap-1 transition-colors cursor-pointer"
                            :title="inv.invite_url"
                          >
                            <Check v-if="copiedInviteId === inv.id" class="w-3.5 h-3.5 text-emerald-600" />
                            <Copy v-else class="w-3.5 h-3.5 text-slate-500" />
                            <span>{{ copiedInviteId === inv.id ? 'Tersalin' : 'Salin Tautan' }}</span>
                          </button>

                          <!-- Revoke Invite -->
                          <button
                            v-if="inv.status === 'pending'"
                            @click="handleRevokeInvitation(inv)"
                            class="py-1 px-2 rounded-lg border border-rose-200 hover:bg-rose-50 text-rose-600 text-[11px] font-semibold transition-colors cursor-pointer"
                            title="Batalkan Undangan"
                          >
                            <X class="w-3.5 h-3.5" />
                          </button>
                        </div>
                      </td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>

            <!-- Invitations Mobile Cards -->
            <div class="md:hidden space-y-3">
              <div v-if="invitationsList.length === 0" class="p-8 text-center text-slate-400 bg-white rounded-2xl border border-slate-200">
                <Mail class="w-8 h-8 text-slate-300 mx-auto mb-2" />
                <p class="font-medium text-xs text-slate-600">Belum ada undangan</p>
                <button
                  @click="openInviteModal"
                  class="mt-3 px-3.5 py-1.5 rounded-xl bg-indigo-600 text-white font-bold text-xs"
                >
                  Undang Klien
                </button>
              </div>

              <div
                v-for="inv in invitationsList"
                :key="inv.id"
                class="p-4 rounded-2xl bg-white border border-slate-200 shadow-sm space-y-3"
              >
                <div class="flex items-start justify-between gap-2">
                  <div>
                    <span class="font-bold text-slate-900 text-xs block">{{ inv.email }}</span>
                    <span v-if="inv.name" class="text-[11px] text-slate-500">{{ inv.name }}</span>
                  </div>
                  <span
                    :class="[
                      'px-2 py-0.5 rounded-full text-[9px] font-mono font-bold uppercase border',
                      inv.role === 'client' ? 'bg-amber-50 text-amber-800 border-amber-200' : 'bg-blue-50 text-blue-800 border-blue-200'
                    ]"
                  >
                    {{ formatRoleName(inv.role) }}
                  </span>
                </div>

                <div v-if="inv.project_ids && inv.project_ids.length > 0" class="flex flex-wrap gap-1">
                  <span
                    v-for="pid in inv.project_ids"
                    :key="pid"
                    class="inline-flex items-center gap-1 px-2 py-0.5 rounded-md bg-slate-100 text-slate-700 text-[10px]"
                  >
                    <span class="w-1.5 h-1.5 rounded-full" :style="{ backgroundColor: getProjectColor(pid) }"></span>
                    <span>{{ getProjectName(pid) }}</span>
                  </span>
                </div>

                <div class="flex items-center justify-between pt-2 border-t border-slate-100 text-[11px]">
                  <span class="text-slate-500 font-mono text-[10px]">{{ formatExpiryDate(inv.expires_at) }}</span>
                  <div class="flex items-center gap-1.5">
                    <button
                      v-if="inv.status === 'pending'"
                      @click="resendInvitationEmail(inv)"
                      :disabled="isResendingId === inv.id"
                      class="py-1 px-2 rounded-lg border border-indigo-200 text-indigo-700 font-semibold text-xs flex items-center gap-1"
                    >
                      <Loader2 v-if="isResendingId === inv.id" class="w-3 h-3 animate-spin text-indigo-600" />
                      <Send v-else class="w-3 h-3 text-indigo-600" />
                      <span>Kirim</span>
                    </button>
                    <button
                      v-if="inv.status === 'pending'"
                      @click="copyInvitationLink(inv)"
                      class="py-1 px-2.5 rounded-lg border border-slate-200 text-slate-700 font-semibold text-xs flex items-center gap-1"
                    >
                      <Check v-if="copiedInviteId === inv.id" class="w-3.5 h-3.5 text-emerald-600" />
                      <Copy v-else class="w-3.5 h-3.5 text-slate-500" />
                      <span>Salin</span>
                    </button>
                    <button
                      v-if="inv.status === 'pending'"
                      @click="handleRevokeInvitation(inv)"
                      class="py-1 px-2 rounded-lg border border-rose-200 text-rose-600"
                    >
                      <X class="w-3.5 h-3.5" />
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </template>
      </div>
    </template>

    <!-- ==================== CREATE / EDIT USER MODAL ==================== -->
    <div
      v-if="showUserModal"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/50 backdrop-blur-xs select-none"
      @click.self="showUserModal = false"
    >
      <div class="max-w-xl w-full bg-white rounded-3xl border border-slate-200 shadow-2xl p-6 sm:p-7 space-y-5 relative max-h-[90vh] overflow-y-auto animate-in fade-in zoom-in-95 duration-150">
        <!-- Header -->
        <div class="flex items-center justify-between pb-3 border-b border-slate-100">
          <div class="flex items-center gap-2.5">
            <div class="p-2 rounded-xl bg-brand-50 border border-brand-200 text-brand-700">
              <UserCheck v-if="isEditingUser" class="w-5 h-5" />
              <UserPlus v-else class="w-5 h-5" />
            </div>
            <div>
              <h3 class="text-sm sm:text-base font-bold text-slate-900">
                {{ isEditingUser ? 'Perbarui Data Pengguna' : 'Tambah Pengguna Baru' }}
              </h3>
              <p class="text-xs text-slate-500">
                {{ isEditingUser ? 'Ubah informasi profil, peran, kata sandi, dan hak akses proyek' : 'Daftarkan pengguna baru dan berikan hak akses proyek' }}
              </p>
            </div>
          </div>
          <button @click="showUserModal = false" class="p-1 rounded-lg text-slate-400 hover:text-slate-600 hover:bg-slate-100 cursor-pointer">
            <X class="w-4 h-4" />
          </button>
        </div>

        <!-- Error Alert Banner -->
        <div v-if="modalErrorMessage" class="p-3 rounded-xl bg-rose-50 border border-rose-200 text-rose-800 text-xs flex items-center gap-2">
          <AlertCircle class="w-4 h-4 text-rose-500 flex-shrink-0" />
          <span class="flex-1 font-medium">{{ modalErrorMessage }}</span>
        </div>

        <!-- Form -->
        <form @submit.prevent="submitUserForm" class="space-y-4 text-xs">
          <!-- Full Name -->
          <div class="space-y-1">
            <label class="font-semibold text-slate-700">Nama Lengkap *</label>
            <input
              type="text"
              v-model="userForm.name"
              required
              placeholder="Contoh: Rian Pratama"
              class="w-full px-3 py-2 rounded-xl border border-slate-300 bg-slate-50/50 focus:bg-white text-slate-900 placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-brand-500 transition-all"
            />
          </div>

          <!-- Email -->
          <div class="space-y-1">
            <label class="font-semibold text-slate-700">Alamat Email *</label>
            <input
              type="email"
              v-model="userForm.email"
              required
              placeholder="nama@pulse.workspace"
              class="w-full px-3 py-2 rounded-xl border border-slate-300 bg-slate-50/50 focus:bg-white text-slate-900 placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-brand-500 transition-all"
            />
          </div>

          <!-- Role Selection -->
          <div class="space-y-1.5">
            <label class="font-semibold text-slate-700">Peran Sistem (Global Role) *</label>
            <div class="grid grid-cols-2 gap-2">
              <label
                v-for="r in roleOptions"
                :key="r.value"
                :class="[
                  'p-2.5 rounded-xl border cursor-pointer transition-all flex flex-col justify-between text-left',
                  userForm.role === r.value
                    ? 'border-brand-500 bg-brand-50/60 ring-1 ring-brand-500/50'
                    : 'border-slate-200 bg-slate-50/50 hover:bg-white hover:border-slate-300'
                ]"
              >
                <div class="flex items-center justify-between">
                  <input
                    type="radio"
                    name="userRole"
                    :value="r.value"
                    v-model="userForm.role"
                    class="sr-only"
                  />
                  <span class="font-bold text-slate-900 text-xs">{{ r.label }}</span>
                  <span :class="['w-2 h-2 rounded-full', r.dotColor]"></span>
                </div>
                <p class="text-[10px] text-slate-500 mt-1 leading-snug">{{ r.desc }}</p>
              </label>
            </div>
          </div>

          <!-- Project Access Selection (Checkboxes) -->
          <div class="space-y-2 p-3.5 rounded-2xl bg-gradient-to-br from-slate-50 to-indigo-50/30 border border-slate-200">
            <div class="flex items-center justify-between flex-wrap gap-1">
              <div>
                <label class="font-bold text-slate-800 flex items-center gap-1.5 text-xs">
                  <FolderKanban class="w-3.5 h-3.5 text-brand-600" />
                  <span>Hak Akses Proyek (Pilih Proyek)</span>
                </label>
                <p class="text-[10px] text-slate-500 mt-0.5">
                  Centang proyek yang dapat langsung dibuka oleh pengguna ini (dapat diubah kapan saja)
                </p>
              </div>
              <div class="flex items-center gap-1.5 text-[10px]">
                <button
                  type="button"
                  @click="selectAllProjects"
                  class="font-semibold text-brand-600 hover:text-brand-700 hover:underline cursor-pointer"
                >
                  Pilih Semua
                </button>
                <span class="text-slate-300">|</span>
                <button
                  type="button"
                  @click="clearAllProjects"
                  class="font-semibold text-slate-500 hover:text-slate-700 hover:underline cursor-pointer"
                >
                  Kosongkan
                </button>
              </div>
            </div>

            <!-- Project Checkbox Grid -->
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-2 pt-1 max-h-48 overflow-y-auto pr-0.5">
              <label
                v-for="p in allAvailableProjects"
                :key="p.id"
                :class="[
                  'p-2.5 rounded-xl border cursor-pointer transition-all flex items-center justify-between text-left select-none',
                  selectedProjectIds.includes(p.id)
                    ? 'border-brand-500 bg-brand-50/80 ring-1 ring-brand-500/40 text-brand-900 shadow-2xs'
                    : 'border-slate-200 bg-white hover:border-slate-300 text-slate-700'
                ]"
              >
                <div class="flex items-center gap-2 min-w-0 pr-2">
                  <input
                    type="checkbox"
                    :value="p.id"
                    v-model="selectedProjectIds"
                    class="rounded text-brand-600 focus:ring-brand-500 border-slate-300 w-3.5 h-3.5 cursor-pointer"
                  />
                  <span class="w-2 h-2 rounded-full flex-shrink-0" :style="{ backgroundColor: p.color_dot }"></span>
                  <div class="truncate">
                    <div class="font-bold text-xs truncate">{{ p.name }}</div>
                    <div class="text-[9px] font-mono text-slate-400">{{ p.slug }}</div>
                  </div>
                </div>
                <span
                  v-if="selectedProjectIds.includes(p.id)"
                  class="text-[9px] font-bold font-mono px-1.5 py-0.2 rounded-full bg-brand-100 text-brand-700 flex-shrink-0"
                >
                  Aktif
                </span>
              </label>
            </div>
            
            <div v-if="allAvailableProjects.length === 0" class="text-center py-3 text-slate-400 text-[11px]">
              Belum ada proyek terdaftar di sistem.
            </div>
          </div>

          <!-- Password Field -->
          <div class="space-y-1">
            <div class="flex items-center justify-between">
              <label class="font-semibold text-slate-700">
                {{ isEditingUser ? 'Kata Sandi Baru (Opsional)' : 'Kata Sandi (Password) *' }}
              </label>
              <span v-if="isEditingUser" class="text-[10px] text-slate-400 font-mono">Kosongkan jika tidak diubah</span>
            </div>
            <div class="relative">
              <input
                :type="showModalPassword ? 'text' : 'password'"
                v-model="userForm.password"
                :required="!isEditingUser"
                placeholder="Minimal 6 karakter"
                class="w-full pl-3 pr-9 py-2 rounded-xl border border-slate-300 bg-slate-50/50 focus:bg-white text-slate-900 placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-brand-500 transition-all"
              />
              <button
                type="button"
                @click="showModalPassword = !showModalPassword"
                class="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600"
              >
                <EyeOff v-if="showModalPassword" class="w-4 h-4" />
                <Eye v-else class="w-4 h-4" />
              </button>
            </div>
          </div>

          <!-- Avatar URL with Auto Generate Preset -->
          <div class="space-y-1">
            <div class="flex items-center justify-between">
              <label class="font-semibold text-slate-700">URL Avatar / Foto Profil</label>
              <button
                type="button"
                @click="generateRandomAvatar"
                class="text-[10px] text-brand-600 hover:text-brand-700 font-medium flex items-center gap-1"
              >
                <Sparkles class="w-3 h-3" />
                <span>Preset Otomatis</span>
              </button>
            </div>
            <div class="flex items-center gap-2">
              <img
                :src="userForm.avatar_url || 'https://images.unsplash.com/photo-1534528741775-53994a69daeb?w=100&auto=format&fit=crop&q=80'"
                alt="Preview"
                class="w-9 h-9 rounded-xl object-cover ring-1 ring-slate-200 flex-shrink-0"
              />
              <input
                type="text"
                v-model="userForm.avatar_url"
                placeholder="https://images.unsplash.com/photo-..."
                class="flex-1 px-3 py-2 rounded-xl border border-slate-300 bg-slate-50/50 focus:bg-white text-slate-900 placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-brand-500 transition-all text-xs"
              />
            </div>
          </div>

          <!-- Actions -->
          <div class="pt-3 border-t border-slate-100 flex items-center justify-end gap-2">
            <button
              type="button"
              @click="showUserModal = false"
              class="px-4 py-2 rounded-xl border border-slate-200 hover:bg-slate-50 text-slate-700 font-semibold cursor-pointer"
            >
              Batal
            </button>
            <button
              type="submit"
              :disabled="isSubmittingUser"
              class="px-5 py-2 rounded-xl bg-gradient-to-r from-brand-600 to-indigo-600 hover:from-brand-500 hover:to-indigo-500 text-white font-bold shadow-md shadow-brand-600/20 disabled:opacity-50 flex items-center gap-1.5 cursor-pointer"
            >
              <Loader2 v-if="isSubmittingUser" class="w-4 h-4 animate-spin" />
              <span>{{ isEditingUser ? 'Simpan Perubahan' : 'Buat Pengguna' }}</span>
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- ==================== DELETE USER CONFIRMATION MODAL ==================== -->
    <div
      v-if="showDeleteModal && userToDelete"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/50 backdrop-blur-xs select-none"
      @click.self="showDeleteModal = false"
    >
      <div class="max-w-md w-full bg-white rounded-3xl border border-slate-200 shadow-2xl p-6 space-y-4 animate-in fade-in zoom-in-95 duration-150">
        <div class="flex items-center gap-3">
          <div class="p-3 rounded-2xl bg-rose-50 border border-rose-200 text-rose-600 flex-shrink-0">
            <Trash2 class="w-6 h-6" />
          </div>
          <div>
            <h3 class="text-sm sm:text-base font-bold text-slate-900">Hapus Akun Pengguna?</h3>
            <p class="text-xs text-slate-500">Tindakan ini tidak dapat dibatalkan</p>
          </div>
        </div>

        <div class="p-3 rounded-xl bg-slate-50 border border-slate-200/80 text-xs space-y-1">
          <div class="font-bold text-slate-900">{{ userToDelete.name }}</div>
          <div class="text-slate-500 font-mono text-[11px]">{{ userToDelete.email }}</div>
          <div class="text-[10px] text-slate-400 font-mono">ID: {{ userToDelete.id }}</div>
        </div>

        <p class="text-xs text-slate-600 leading-relaxed">
          Pengguna ini akan dihapus dari sistem, seluruh sesi perangkat akan diputus, dan keanggotaan pada proyek-proyek terkait akan dicabut secara otomatis.
        </p>

        <div class="pt-2 flex items-center justify-end gap-2">
          <button
            @click="showDeleteModal = false"
            class="px-4 py-2 rounded-xl border border-slate-200 hover:bg-slate-50 text-slate-700 text-xs font-semibold cursor-pointer"
          >
            Batal
          </button>
          <button
            @click="confirmDeleteUser"
            :disabled="isSubmittingUser"
            class="px-4 py-2 rounded-xl bg-rose-600 hover:bg-rose-500 text-white text-xs font-bold shadow-md shadow-rose-600/20 disabled:opacity-50 flex items-center gap-1.5 cursor-pointer"
          >
            <Loader2 v-if="isSubmittingUser" class="w-3.5 h-3.5 animate-spin" />
            <span>Ya, Hapus Pengguna</span>
          </button>
        </div>
      </div>
    </div>

    <!-- ==================== IMPORT PROJECT MODAL (EXCEL / CSV) ==================== -->
    <div
      v-if="showImportModal"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/50 backdrop-blur-xs select-none"
      @click.self="closeImportModal"
    >
      <div class="max-w-2xl w-full bg-white rounded-3xl border border-slate-200 shadow-2xl p-6 sm:p-7 space-y-5 relative max-h-[90vh] flex flex-col animate-in fade-in zoom-in-95 duration-150">
        <!-- Header -->
        <div class="flex items-center justify-between pb-3 border-b border-slate-100 flex-shrink-0">
          <div class="flex items-center gap-2.5">
            <div class="p-2.5 rounded-2xl bg-emerald-50 border border-emerald-200 text-emerald-600">
              <FileSpreadsheet class="w-5 h-5" />
            </div>
            <div>
              <h3 class="text-base sm:text-lg font-bold text-slate-900">Impor Proyek (MS Excel / CSV)</h3>
              <p class="text-xs text-slate-500">Unggah file spreadsheet untuk menambahkan atau memperbarui data proyek</p>
            </div>
          </div>
          <button
            @click="closeImportModal"
            class="p-2 rounded-xl text-slate-400 hover:text-slate-600 hover:bg-slate-100 transition-colors cursor-pointer"
          >
            <X class="w-5 h-5" />
          </button>
        </div>

        <!-- Scrollable Body -->
        <div class="space-y-4 overflow-y-auto pr-1 flex-1">
          <!-- Download Template Banner -->
          <div class="p-4 rounded-2xl bg-slate-50 border border-slate-200/80 flex flex-col sm:flex-row sm:items-center justify-between gap-3">
            <div class="flex items-start gap-2.5">
              <div class="p-1.5 rounded-xl bg-indigo-50 border border-indigo-200 text-indigo-600 mt-0.5">
                <FileText class="w-4 h-4" />
              </div>
              <div>
                <div class="text-xs font-bold text-slate-800">Gunakan Template Standar</div>
                <p class="text-[11px] text-slate-500 leading-tight">Format CSV dengan encoding UTF-8 BOM, kompatibel langsung dengan Microsoft Excel & Google Sheets.</p>
              </div>
            </div>
            <button
              type="button"
              @click="handleDownloadTemplate"
              class="px-3 py-1.5 rounded-xl bg-white hover:bg-slate-100 border border-slate-200 text-brand-600 text-xs font-semibold flex items-center justify-center gap-1.5 shadow-2xs whitespace-nowrap cursor-pointer transition-colors"
            >
              <Download class="w-3.5 h-3.5" />
              <span>Unduh Template CSV</span>
            </button>
          </div>

          <!-- Dropzone -->
          <div
            @dragover.prevent="isDragging = true"
            @dragleave.prevent="isDragging = false"
            @drop.prevent="handleFileDrop"
            @click="triggerFileInput"
            :class="[
              'border-2 border-dashed rounded-2xl p-6 text-center cursor-pointer transition-all flex flex-col items-center justify-center gap-2 select-none',
              isDragging
                ? 'border-brand-500 bg-brand-50/50'
                : 'border-slate-300 hover:border-brand-400 bg-slate-50/50 hover:bg-white'
            ]"
          >
            <input
              type="file"
              ref="importFileInputRef"
              class="hidden"
              accept=".csv, text/csv, application/vnd.ms-excel, text/plain"
              @change="handleFileInputChange"
            />

            <div class="p-3 rounded-2xl bg-brand-50 text-brand-600 border border-brand-100">
              <FileUp class="w-6 h-6" />
            </div>

            <div v-if="!selectedImportFile">
              <div class="text-xs sm:text-sm font-bold text-slate-800">
                Pilih file CSV atau seret ke area ini
              </div>
              <div class="text-[11px] text-slate-400 mt-0.5">
                Mendukung file .csv yang diekspor dari MS Excel (UTF-8 / ANSI)
              </div>
            </div>

            <div v-else class="flex items-center gap-2 text-xs font-semibold text-emerald-700 bg-emerald-50 px-3 py-1.5 rounded-xl border border-emerald-200">
              <CheckCircle2 class="w-4 h-4 text-emerald-600" />
              <span>{{ selectedImportFile.name }} ({{ formatFileSize(selectedImportFile.size) }})</span>
            </div>
          </div>

          <!-- Feedback Banners -->
          <div v-if="importErrorMessage" class="p-3 rounded-xl bg-rose-50 border border-rose-200 text-rose-700 text-xs flex items-start gap-2">
            <AlertCircle class="w-4 h-4 mt-0.5 flex-shrink-0" />
            <span>{{ importErrorMessage }}</span>
          </div>

          <div v-if="importSuccessMessage" class="p-3 rounded-xl bg-emerald-50 border border-emerald-200 text-emerald-800 text-xs flex items-start gap-2">
            <CheckCircle2 class="w-4 h-4 mt-0.5 flex-shrink-0 text-emerald-600" />
            <span>{{ importSuccessMessage }}</span>
          </div>

          <!-- Live Parsed Preview Table -->
          <div v-if="parsedPreviewRows.length > 0" class="space-y-2">
            <div class="flex items-center justify-between text-xs">
              <span class="font-bold text-slate-700 flex items-center gap-1.5">
                <span>Pratinjau Data Impor</span>
                <span class="px-2 py-0.5 rounded-full text-[10px] font-mono bg-brand-50 text-brand-700 border border-brand-200 font-semibold">
                  {{ parsedPreviewRows.length }} baris siap diproses
                </span>
              </span>
              <span class="text-[11px] text-slate-400">Pastikan kolom Nama Proyek terisi</span>
            </div>

            <div class="border border-slate-200 rounded-xl overflow-hidden max-h-48 overflow-y-auto">
              <table class="w-full text-left text-xs">
                <thead class="bg-slate-100 text-slate-600 font-mono text-[10px] uppercase border-b border-slate-200">
                  <tr>
                    <th class="py-2 px-3">No</th>
                    <th class="py-2 px-3">Nama Proyek</th>
                    <th class="py-2 px-3">Slug</th>
                    <th class="py-2 px-3">Klien</th>
                    <th class="py-2 px-3">Status</th>
                    <th class="py-2 px-3">Tenggat</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-slate-100 text-[11px]">
                  <tr v-for="(row, idx) in parsedPreviewRows" :key="idx" class="hover:bg-slate-50">
                    <td class="py-1.5 px-3 font-mono text-slate-400">{{ idx + 1 }}</td>
                    <td class="py-1.5 px-3 font-bold text-slate-800 flex items-center gap-1.5">
                      <span
                        class="w-2 h-2 rounded-full flex-shrink-0"
                        :style="{ backgroundColor: row.color_dot || '#10B981' }"
                      ></span>
                      <span>{{ row.name }}</span>
                    </td>
                    <td class="py-1.5 px-3 font-mono text-slate-500">{{ row.slug || '-' }}</td>
                    <td class="py-1.5 px-3 text-slate-600">{{ row.client_name || 'Internal' }}</td>
                    <td class="py-1.5 px-3">
                      <span class="px-1.5 py-0.5 rounded text-[10px] font-medium capitalize bg-slate-100 text-slate-700">
                        {{ row.status || 'active' }}
                      </span>
                    </td>
                    <td class="py-1.5 px-3 font-mono text-slate-500">{{ row.deadline || '-' }}</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>

        <!-- Footer Actions -->
        <div class="pt-3 border-t border-slate-100 flex items-center justify-between flex-shrink-0">
          <div class="text-[11px] text-slate-400">
            Proyek yang sudah ada akan otomatis diperbarui (Upsert).
          </div>
          <div class="flex items-center gap-2">
            <button
              type="button"
              @click="closeImportModal"
              class="px-4 py-2 rounded-xl border border-slate-200 hover:bg-slate-50 text-slate-700 text-xs font-semibold cursor-pointer"
            >
              Batal
            </button>
            <button
              type="button"
              @click="executeImport"
              :disabled="!selectedImportFile || isImporting"
              class="px-5 py-2 rounded-xl bg-gradient-to-r from-emerald-600 to-teal-600 hover:from-emerald-500 hover:to-teal-500 text-white text-xs font-bold shadow-md shadow-emerald-600/20 disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-1.5 cursor-pointer transition-all"
            >
              <Loader2 v-if="isImporting" class="w-3.5 h-3.5 animate-spin" />
              <Upload v-else class="w-3.5 h-3.5" />
              <span>{{ isImporting ? 'Mengimpor...' : 'Mulai Impor Proyek' }}</span>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- ==================== EDIT PROJECT MODAL ==================== -->
    <EditProjectModal
      :is-open="showEditProjectModal"
      :project="selectedProjectForEdit"
      @close="showEditProjectModal = false"
      @updated="store.fetchProjects()"
    />

    <!-- ==================== PROJECT ACCESS MANAGEMENT MODAL ==================== -->
    <ProjectAccessModal
      :is-open="showProjectAccessModal"
      :project="selectedProjectForAccess"
      @close="showProjectAccessModal = false"
      @updated="store.fetchProjects()"
    />

    <!-- ==================== COPY PROJECT MODAL ==================== -->
    <CopyProjectModal
      :is-open="showCopyProjectModal"
      :project="selectedProjectForCopy"
      @close="showCopyProjectModal = false"
      @copied="onProjectCopied"
    />

    <!-- ==================== INVITE USER MODAL ==================== -->
    <InviteUserModal
      :is-open="showInviteModal"
      :projects="store.projects"
      @close="showInviteModal = false"
      @invited="loadInvitations"
    />

    <!-- ==================== SMTP SETTINGS & TEST MODAL ==================== -->
    <div
      v-if="showSMTPModal"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/50 backdrop-blur-xs select-none"
      @click.self="showSMTPModal = false"
    >
      <div class="max-w-lg w-full bg-white rounded-3xl border border-slate-200 shadow-2xl p-6 sm:p-7 space-y-5 relative max-h-[90vh] overflow-y-auto animate-in fade-in zoom-in-95 duration-150">
        <!-- Header -->
        <div class="flex items-start justify-between">
          <div class="flex items-center gap-3">
            <div class="p-2.5 rounded-2xl bg-indigo-50 border border-indigo-100 text-indigo-600">
              <Mail class="w-5 h-5 stroke-[2.2]" />
            </div>
            <div>
              <h3 class="text-base font-bold text-slate-900 font-display">Layanan Email (SMTP)</h3>
              <p class="text-xs text-slate-500">Status & Pengujian Pengiriman Email Undangan</p>
            </div>
          </div>
          <button @click="showSMTPModal = false" class="p-1 rounded-lg text-slate-400 hover:text-slate-600 hover:bg-slate-100 cursor-pointer">
            <X class="w-4 h-4" />
          </button>
        </div>

        <!-- Status Card -->
        <div
          class="p-4 rounded-2xl border text-xs space-y-2.5"
          :class="smtpStatus.configured ? 'bg-emerald-50/70 border-emerald-200 text-emerald-950' : 'bg-amber-50/70 border-amber-200 text-amber-950'"
        >
          <div class="flex items-center justify-between">
            <span class="font-bold flex items-center gap-1.5">
              <CheckCircle2 v-if="smtpStatus.configured" class="w-4 h-4 text-emerald-600" />
              <AlertCircle v-else class="w-4 h-4 text-amber-600" />
              <span>Status Koneksi: {{ smtpStatus.configured ? 'Terkonfigurasi (Siap Kirim)' : 'Belum Dikonfigurasi Lengkap' }}</span>
            </span>
            <span class="font-mono text-[10px] px-2 py-0.5 rounded-md bg-white border border-slate-200">
              Port: {{ smtpStatus.smtp_port }}
            </span>
          </div>

          <div class="grid grid-cols-2 gap-2 pt-2 text-[11px]">
            <div>
              <span class="text-slate-500 block">Host:</span>
              <span class="font-semibold font-mono">{{ smtpStatus.smtp_host || 'smtp.gmail.com' }}</span>
            </div>
            <div>
              <span class="text-slate-500 block">Akun Pengirim:</span>
              <span class="font-semibold font-mono">{{ smtpStatus.smtp_user || '-' }}</span>
            </div>
            <div class="col-span-2">
              <span class="text-slate-500 block">Sandi Aplikasi (Masked):</span>
              <span class="font-semibold font-mono">{{ smtpStatus.has_pass ? smtpStatus.masked_pass : '(Belum Diisi di backend/.env)' }}</span>
            </div>
          </div>

          <p v-if="!smtpStatus.configured" class="text-[10px] text-amber-800 leading-relaxed pt-1 border-t border-amber-200/60">
            * Untuk mengaktifkan pengiriman email otomatis via Google, masukkan <strong>App Password 16-digit</strong> pada file <code>backend/.env</code> parameter <code>SMTP_PASS</code>.
          </p>
        </div>

        <!-- Test Email Dispatcher -->
        <div class="space-y-3 pt-1">
          <label class="font-bold text-xs text-slate-800 flex items-center gap-1.5">
            <Send class="w-3.5 h-3.5 text-indigo-600" />
            <span>Kirim Email Uji Coba (Test Send)</span>
          </label>
          <p class="text-[11px] text-slate-500">
            Kirimkan satu email uji coba untuk memverifikasi apakah akun SMTP berhasil terkoneksi ke Google/Mail Server.
          </p>

          <div class="flex items-center gap-2">
            <input
              type="email"
              v-model="testEmailTarget"
              placeholder="email-anda@gmail.com"
              class="flex-1 px-3 py-2 rounded-xl border border-slate-300 text-xs text-slate-800 focus:outline-none focus:ring-2 focus:ring-indigo-500"
            />
            <button
              @click="executeTestSMTP"
              :disabled="isTestingSMTP || !testEmailTarget"
              class="px-4 py-2 rounded-xl bg-indigo-600 hover:bg-indigo-500 text-white font-bold text-xs shadow-md shadow-indigo-600/20 disabled:opacity-50 flex items-center gap-1.5 cursor-pointer transition-all flex-shrink-0"
            >
              <Loader2 v-if="isTestingSMTP" class="w-3.5 h-3.5 animate-spin" />
              <Send v-else class="w-3.5 h-3.5" />
              <span>{{ isTestingSMTP ? 'Mengirim...' : 'Tes Kirim' }}</span>
            </button>
          </div>

          <!-- Test Result Banner -->
          <div
            v-if="testSMTPResponse"
            class="p-3 rounded-xl text-xs flex items-start gap-2 border"
            :class="testSMTPResponse.success ? 'bg-emerald-50 text-emerald-900 border-emerald-200' : 'bg-rose-50 text-rose-900 border-rose-200'"
          >
            <CheckCircle2 v-if="testSMTPResponse.success" class="w-4 h-4 text-emerald-600 flex-shrink-0 mt-0.5" />
            <AlertCircle v-else class="w-4 h-4 text-rose-600 flex-shrink-0 mt-0.5" />
            <span class="leading-relaxed">{{ testSMTPResponse.message }}</span>
          </div>
        </div>

        <!-- Footer -->
        <div class="pt-3 border-t border-slate-100 flex items-center justify-between">
          <span class="text-[10px] text-slate-400">File Konfigurasi: <code>backend/.env</code></span>
          <button
            type="button"
            @click="showSMTPModal = false"
            class="px-4 py-2 rounded-xl bg-slate-100 hover:bg-slate-200 text-slate-700 text-xs font-semibold cursor-pointer"
          >
            Tutup
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useProjectStore } from '../stores/projectStore'
import { useAuthStore } from '../stores/authStore'
import { useRouter, useRoute } from 'vue-router'
import { projectApi, userApi, invitationApi, systemApi, type UserPayload } from '../services/api'
import type { TeamWorkload, BudgetItem, User, Project, InvitationItem } from '../types/project'
import KPIStats from '../components/KPIStats.vue'
import EditProjectModal from '../components/EditProjectModal.vue'
import ProjectAccessModal from '../components/ProjectAccessModal.vue'
import InviteUserModal from '../components/InviteUserModal.vue'
import CopyProjectModal from '../components/CopyProjectModal.vue'
import {
  Clock,
  Download,
  Users,
  LayoutDashboard,
  UserPlus,
  Search,
  RefreshCw,
  Shield,
  ShieldCheck,
  Edit,
  Trash2,
  ZapOff,
  UserCheck,
  X,
  AlertCircle,
  Eye,
  EyeOff,
  Sparkles,
  Loader2,
  FileSpreadsheet,
  Upload,
  ListTodo,
  FileText,
  CheckCircle2,
  FileUp,
  ChevronDown,
  FolderKanban,
  Mail,
  MailPlus,
  Copy,
  Check,
  Send,
  Settings
} from 'lucide-vue-next'

const store = useProjectStore()
const authStore = useAuthStore()
const router = useRouter()
const route = useRoute()

// Active tab on Dashboard: 'overview' (default) or 'users' (Admin CRUD)
const activeDashboardTab = ref<'overview' | 'users'>('overview')

const canManageUsers = computed(() => {
  return authStore.hasRole('super_admin') || authStore.can('user:manage')
})

const canCreateProject = computed(() => {
  return authStore.hasRole('super_admin') || authStore.hasRole('pm') || authStore.can('project:create')
})

const canManageProjectAccess = computed(() => {
  return authStore.hasRole('super_admin') || authStore.hasRole('pm') || authStore.can('project:member:manage')
})

const showProjectAccessModal = ref(false)
const selectedProjectForAccess = ref<Project | null>(null)

const showEditProjectModal = ref(false)
const selectedProjectForEdit = ref<Project | null>(null)

const showCopyProjectModal = ref(false)
const selectedProjectForCopy = ref<Project | null>(null)

function openCopyProjectModal(project: Project) {
  selectedProjectForCopy.value = project
  showCopyProjectModal.value = true
}

async function onProjectCopied(newProject: Project) {
  await store.fetchProjects()
  store.selectProject(newProject.id)
  store.showToast(`Proyek "${newProject.name}" berhasil diduplikasi!`)
}

function openEditProjectModal(project: Project) {
  selectedProjectForEdit.value = project
  showEditProjectModal.value = true
}

async function handleDeleteProject(project: Project) {
  if (!confirm(`Hapus proyek "${project.name}"?\nPerhatian: Seluruh deliverable, tugas, dan data keanggotaan proyek ini akan ikut terhapus secara permanen.`)) {
    return
  }

  try {
    await projectApi.deleteProject(project.id)
    store.projects = store.projects.filter(p => p.id !== project.id)
    store.showToast(`Proyek "${project.name}" berhasil dihapus.`)
  } catch (err: any) {
    store.showToast(err?.response?.data?.error || 'Gagal menghapus proyek')
  }
}

function openProjectAccessModal(project: Project) {
  selectedProjectForAccess.value = project
  showProjectAccessModal.value = true
}

function switchTab(tab: 'overview' | 'users') {
  activeDashboardTab.value = tab
  router.replace({ query: { ...route.query, tab } })
}

// ==================== USER MANAGEMENT & INVITATIONS STATE ====================
const userSubTab = ref<'active' | 'invitations'>('active')
const showInviteModal = ref(false)
const invitationsList = ref<InvitationItem[]>([])
const isLoadingInvitations = ref(false)
const copiedInviteId = ref<string | null>(null)
const isResendingId = ref<string | null>(null)

// SMTP Service state
const smtpStatus = ref<{
  configured: boolean
  smtp_host: string
  smtp_port: number
  smtp_user: string
  smtp_from: string
  smtp_secure: boolean
  has_pass: boolean
  masked_pass: string
}>({
  configured: false,
  smtp_host: 'smtp.gmail.com',
  smtp_port: 587,
  smtp_user: '',
  smtp_from: '',
  smtp_secure: false,
  has_pass: false,
  masked_pass: ''
})
const showSMTPModal = ref(false)
const isTestingSMTP = ref(false)
const testEmailTarget = ref('')
const testSMTPResponse = ref<{ success?: boolean; message?: string } | null>(null)

const pendingInvitationsCount = computed(() => {
  return invitationsList.value.filter((i) => i.status === 'pending').length
})

function openInviteModal() {
  showInviteModal.value = true
}

function openSMTPConfigModal() {
  testEmailTarget.value = authStore.user?.email || 'agus48@polsri.ac.id'
  testSMTPResponse.value = null
  showSMTPModal.value = true
}

async function loadSMTPStatus() {
  if (!canManageUsers.value) return
  try {
    const data = await systemApi.getSMTPStatus()
    if (data) smtpStatus.value = data
  } catch (err: any) {
    console.error('Failed to load SMTP status:', err)
  }
}

async function resendInvitationEmail(inv: InvitationItem) {
  if (!inv || !inv.id) return
  isResendingId.value = inv.id
  try {
    const res = await invitationApi.resendInvitation(inv.id)
    store.showToast(res.message || `Email undangan berhasil dikirim ulang ke ${inv.email}`)
  } catch (err: any) {
    store.showToast('Gagal mengirim email: ' + (err?.response?.data?.error || err.message))
  } finally {
    isResendingId.value = null
  }
}

async function executeTestSMTP() {
  if (!testEmailTarget.value.trim()) return
  isTestingSMTP.value = true
  testSMTPResponse.value = null
  try {
    const res = await systemApi.testSMTP(testEmailTarget.value.trim())
    testSMTPResponse.value = { success: true, message: res.message }
    store.showToast('Email pengujian SMTP sukses terkirim!')
  } catch (err: any) {
    testSMTPResponse.value = { success: false, message: err?.response?.data?.error || err.message }
  } finally {
    isTestingSMTP.value = false
  }
}

async function loadInvitations() {
  if (!canManageUsers.value) return
  isLoadingInvitations.value = true
  try {
    const data = await invitationApi.getInvitations()
    invitationsList.value = data || []
  } catch (err: any) {
    console.error('Failed to load invitations:', err)
  } finally {
    isLoadingInvitations.value = false
  }
}

async function copyInvitationLink(inv: InvitationItem) {
  if (!inv.invite_url) return
  try {
    await navigator.clipboard.writeText(inv.invite_url)
    copiedInviteId.value = inv.id
    store.showToast(`Tautan undangan untuk ${inv.email} berhasil disalin!`)
    setTimeout(() => {
      if (copiedInviteId.value === inv.id) {
        copiedInviteId.value = null
      }
    }, 3000)
  } catch (err) {
    console.error('Failed to copy invite link:', err)
  }
}

async function handleRevokeInvitation(inv: InvitationItem) {
  if (!confirm(`Batalkan undangan untuk email ${inv.email}? Tautan yang telah dikirim tidak akan dapat digunakan lagi.`)) {
    return
  }
  try {
    await invitationApi.revokeInvitation(inv.id)
    store.showToast(`Undangan untuk ${inv.email} berhasil dibatalkan`)
    await loadInvitations()
  } catch (err: any) {
    store.showToast('Gagal membatalkan undangan: ' + (err?.response?.data?.error || err.message))
  }
}

function getProjectName(id: string): string {
  const p = store.projects.find((proj) => proj.id === id)
  return p ? p.name : id
}

function getProjectColor(id: string): string {
  const p = store.projects.find((proj) => proj.id === id)
  return p?.color_dot || '#6366F1'
}

function formatExpiryDate(dStr: string): string {
  if (!dStr) return '-'
  try {
    const d = new Date(dStr)
    return d.toLocaleDateString('id-ID', {
      day: 'numeric',
      month: 'short',
      hour: '2-digit',
      minute: '2-digit'
    })
  } catch {
    return dStr
  }
}

const userList = ref<User[]>([])
const isLoadingUsers = ref(false)
const userSearchQuery = ref('')
const userRoleFilter = ref('all')

const showUserModal = ref(false)
const isEditingUser = ref(false)
const isSubmittingUser = ref(false)
const modalErrorMessage = ref('')
const showModalPassword = ref(false)

const editingUserId = ref<string | null>(null)
const userForm = ref<UserPayload>({
  name: '',
  email: '',
  role: 'member',
  password: '',
  avatar_url: ''
})

const showDeleteModal = ref(false)
const userToDelete = ref<User | null>(null)

const roleOptions = [
  { value: 'super_admin', label: 'Super Admin', desc: 'Akses kontrol penuh ke seluruh sistem & modul pengguna', dotColor: 'bg-rose-500' },
  { value: 'pm', label: 'Project Manager', desc: 'Pengelolaan sprint, task, alokasi tim, dan milestone', dotColor: 'bg-emerald-500' },
  { value: 'member', label: 'Team Member', desc: 'Pengerjaan tugas engineering, attachment, & timesheet', dotColor: 'bg-indigo-500' },
  { value: 'client', label: 'Client Portal', desc: 'Review deliverable dan persetujuan approval milestone', dotColor: 'bg-amber-500' }
]

async function loadUsers() {
  if (!canManageUsers.value) return
  isLoadingUsers.value = true
  try {
    const data = await userApi.getUsers()
    userList.value = data
  } catch (err: any) {
    store.showToast('Gagal memuat data pengguna: ' + (err?.response?.data?.error || err.message))
  } finally {
    isLoadingUsers.value = false
  }
}

const filteredUsers = computed(() => {
  return userList.value.filter((u) => {
    const q = userSearchQuery.value.toLowerCase().trim()
    const matchesSearch =
      q === '' ||
      u.name.toLowerCase().includes(q) ||
      u.email.toLowerCase().includes(q) ||
      u.id.toLowerCase().includes(q)

    const matchesRole =
      userRoleFilter.value === 'all' || u.role === userRoleFilter.value

    return matchesSearch && matchesRole
  })
})

function formatRoleName(role: string): string {
  if (role === 'super_admin') return 'Super Admin'
  if (role === 'pm') return 'Project Mgr'
  if (role === 'member') return 'Team Member'
  if (role === 'client') return 'Client'
  return role
}

function getRoleBadgeClass(role: string): string {
  if (role === 'super_admin') return 'bg-rose-50 text-rose-700 border-rose-200'
  if (role === 'pm') return 'bg-emerald-50 text-emerald-700 border-emerald-200'
  if (role === 'member') return 'bg-indigo-50 text-indigo-700 border-indigo-200'
  if (role === 'client') return 'bg-amber-50 text-amber-700 border-amber-200'
  return 'bg-slate-100 text-slate-700 border-slate-200'
}

function formatLastActive(dateStr?: string | Date): string {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  if (isNaN(d.getTime())) return '-'
  return d.toLocaleDateString('id-ID', {
    day: 'numeric',
    month: 'short',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const allAvailableProjects = computed(() => {
  return store.projects
})

const selectedProjectIds = ref<string[]>([])

function selectAllProjects() {
  selectedProjectIds.value = allAvailableProjects.value.map(p => p.id)
}

function clearAllProjects() {
  selectedProjectIds.value = []
}

function getProjectNames(projectIds?: string[]): string[] {
  if (!projectIds || projectIds.length === 0) return []
  return projectIds.map(id => {
    const p = store.projects.find(proj => proj.id === id)
    return p ? p.name : id
  })
}

function generateRandomAvatar() {
  const seed = userForm.value.name.trim() || `user_${Date.now()}`
  userForm.value.avatar_url = `https://api.dicebear.com/7.x/avataaars/svg?seed=${encodeURIComponent(seed)}`
}

function openCreateUserModal() {
  isEditingUser.value = false
  editingUserId.value = null
  modalErrorMessage.value = ''
  showModalPassword.value = false
  selectedProjectIds.value = []
  userForm.value = {
    name: '',
    email: '',
    role: 'member',
    password: '',
    avatar_url: 'https://images.unsplash.com/photo-1534528741775-53994a69daeb?w=100&auto=format&fit=crop&q=80'
  }
  showUserModal.value = true
}

function openEditUserModal(user: User) {
  isEditingUser.value = true
  editingUserId.value = user.id
  modalErrorMessage.value = ''
  showModalPassword.value = false
  selectedProjectIds.value = user.project_ids ? [...user.project_ids] : []
  userForm.value = {
    name: user.name,
    email: user.email,
    role: user.role,
    password: '',
    avatar_url: user.avatar_url || ''
  }
  showUserModal.value = true
}

async function submitUserForm() {
  modalErrorMessage.value = ''
  isSubmittingUser.value = true
  try {
    if (isEditingUser.value && editingUserId.value) {
      const payload: Partial<UserPayload> = {
        name: userForm.value.name,
        email: userForm.value.email,
        role: userForm.value.role,
        avatar_url: userForm.value.avatar_url,
        project_ids: selectedProjectIds.value
      }
      if (userForm.value.password) {
        payload.password = userForm.value.password
      }
      await userApi.updateUser(editingUserId.value, payload)
      store.showToast(`Data pengguna ${userForm.value.name} dan hak akses proyek berhasil diperbarui`)
    } else {
      if (!userForm.value.password || userForm.value.password.length < 6) {
        modalErrorMessage.value = 'Kata sandi wajib diisi minimal 6 karakter'
        isSubmittingUser.value = false
        return
      }
      await userApi.createUser({
        ...userForm.value,
        project_ids: selectedProjectIds.value
      })
      store.showToast(`Pengguna baru ${userForm.value.name} berhasil dibuat dengan akses proyek`)
    }
    showUserModal.value = false
    await loadUsers()
    await store.fetchProjects()
  } catch (err: any) {
    modalErrorMessage.value = err?.response?.data?.error || err.message || 'Terjadi kesalahan saat menyimpan data'
  } finally {
    isSubmittingUser.value = false
  }
}

async function handleResetSession(user: User) {
  try {
    await userApi.resetUserSession(user.id)
    store.showToast(`Sesi akun ${user.name} berhasil dilepas dari semua perangkat`)
    await loadUsers()
  } catch (err: any) {
    store.showToast('Gagal melepas sesi: ' + (err?.response?.data?.error || err.message))
  }
}

function openDeleteUserModal(user: User) {
  userToDelete.value = user
  showDeleteModal.value = true
}

async function confirmDeleteUser() {
  if (!userToDelete.value) return
  isSubmittingUser.value = true
  try {
    await userApi.deleteUser(userToDelete.value.id)
    store.showToast(`Pengguna ${userToDelete.value.name} berhasil dihapus`)
    showDeleteModal.value = false
    userToDelete.value = null
    await loadUsers()
  } catch (err: any) {
    store.showToast('Gagal menghapus pengguna: ' + (err?.response?.data?.error || err.message))
  } finally {
    isSubmittingUser.value = false
  }
}

// ==================== WORKLOAD & BUDGET STATE ====================
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
  if (route.query.tab === 'users') {
    activeDashboardTab.value = 'users'
  }

  loadUsers()
  loadInvitations()
  loadSMTPStatus()

  try {
    const wl = await projectApi.getTeamWorkload()
    if (wl && wl.length > 0) teamWorkload.value = wl
  } catch (_) {}

  try {
    const bg = await projectApi.getBudgets(store.currentProjectId || 'proj_sustagri')
    if (bg && bg.budgets && bg.budgets.length > 0) budgetList.value = bg.budgets
  } catch (_) {}
})

watch(() => route.query.tab, (tab) => {
  if (tab === 'users') {
    activeDashboardTab.value = 'users'
    loadUsers()
    loadInvitations()
    loadSMTPStatus()
  } else {
    activeDashboardTab.value = 'overview'
  }
})

// ==================== PROJECT IMPORT & EXPORT (MS EXCEL / CSV) ====================
const showExportDropdown = ref(false)
const isExportingProjects = ref(false)
const isExportingTasks = ref(false)

const showImportModal = ref(false)
const selectedImportFile = ref<File | null>(null)
const parsedPreviewRows = ref<any[]>([])
const importFileInputRef = ref<HTMLInputElement | null>(null)
const isDragging = ref(false)
const isImporting = ref(false)
const importErrorMessage = ref('')
const importSuccessMessage = ref('')

function triggerBlobDownload(blob: Blob, filename: string) {
  const url = window.URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = filename
  document.body.appendChild(a)
  a.click()
  window.URL.revokeObjectURL(url)
  document.body.removeChild(a)
}

async function handleExportProjects() {
  showExportDropdown.value = false
  isExportingProjects.value = true
  try {
    store.showToast('Menyiapkan file ekspor portfolio proyek MS Excel / CSV...')
    const blob = await projectApi.exportProjectsCSV()
    const timestamp = new Date().toISOString().slice(0, 10).replace(/-/g, '')
    triggerBlobDownload(blob, `pulse_proyek_export_${timestamp}.csv`)
    store.showToast('Ekspor proyek berhasil diunduh.')
  } catch (err: any) {
    store.showToast(err.response?.data?.error || 'Gagal mengekspor proyek.')
  } finally {
    isExportingProjects.value = false
  }
}

async function handleExportTasks() {
  showExportDropdown.value = false
  isExportingTasks.value = true
  try {
    store.showToast('Menyiapkan file ekspor tasks & deliverables...')
    const blob = await projectApi.exportTasksCSV(store.currentProjectId)
    const timestamp = new Date().toISOString().slice(0, 10).replace(/-/g, '')
    triggerBlobDownload(blob, `pulse_tasks_export_${timestamp}.csv`)
    store.showToast('Ekspor tasks berhasil diunduh.')
  } catch (err: any) {
    store.showToast(err.response?.data?.error || 'Gagal mengekspor tasks.')
  } finally {
    isExportingTasks.value = false
  }
}

async function handleDownloadTemplate() {
  showExportDropdown.value = false
  try {
    store.showToast('Mengunduh template CSV resmi untuk Microsoft Excel...')
    const blob = await projectApi.downloadImportTemplate()
    triggerBlobDownload(blob, 'template_import_proyek.csv')
    store.showToast('Template import berhasil diunduh.')
  } catch (err: any) {
    store.showToast(err.response?.data?.error || 'Gagal mengunduh template.')
  }
}

function openImportModal() {
  selectedImportFile.value = null
  parsedPreviewRows.value = []
  importErrorMessage.value = ''
  importSuccessMessage.value = ''
  showImportModal.value = true
}

function closeImportModal() {
  showImportModal.value = false
  selectedImportFile.value = null
  parsedPreviewRows.value = []
  importErrorMessage.value = ''
  importSuccessMessage.value = ''
}

function triggerFileInput() {
  importFileInputRef.value?.click()
}

function handleFileInputChange(e: Event) {
  const target = e.target as HTMLInputElement
  if (target.files && target.files.length > 0) {
    processSelectedFile(target.files[0])
  }
}

function handleFileDrop(e: DragEvent) {
  isDragging.value = false
  if (e.dataTransfer?.files && e.dataTransfer.files.length > 0) {
    processSelectedFile(e.dataTransfer.files[0])
  }
}

function formatFileSize(bytes: number) {
  if (bytes < 1024) return bytes + ' B'
  return (bytes / 1024).toFixed(1) + ' KB'
}

function parseCSVClientSide(text: string) {
  // Strip UTF-8 BOM if present
  const cleanText = text.replace(/^\uFEFF/, '')
  const lines = cleanText.split(/\r?\n/).filter(line => line.trim().length > 0)
  if (lines.length < 2) return []

  const parseLine = (line: string): string[] => {
    const result: string[] = []
    let cur = ''
    let inQuotes = false
    for (let i = 0; i < line.length; i++) {
      const char = line[i]
      if (char === '"') {
        if (inQuotes && line[i + 1] === '"') {
          cur += '"'
          i++
        } else {
          inQuotes = !inQuotes
        }
      } else if (char === ',' && !inQuotes) {
        result.push(cur.trim())
        cur = ''
      } else {
        cur += char
      }
    }
    result.push(cur.trim())
    return result
  }

  const header = parseLine(lines[0])
  let nameIdx = -1, slugIdx = -1, clientIdx = -1, statusIdx = -1, deadlineIdx = -1, colorIdx = -1

  header.forEach((col, idx) => {
    const norm = col.toLowerCase().trim()
    if (norm.includes('nama') && !norm.includes('klien')) nameIdx = idx
    else if (norm.includes('slug') || norm.includes('kode')) slugIdx = idx
    else if (norm.includes('klien') || norm.includes('client')) clientIdx = idx
    else if (norm.includes('status')) statusIdx = idx
    else if (norm.includes('tenggat') || norm.includes('deadline')) deadlineIdx = idx
    else if (norm.includes('warna') || norm.includes('color')) colorIdx = idx
  })

  if (nameIdx === -1) nameIdx = 0

  const items: any[] = []
  for (let i = 1; i < lines.length; i++) {
    const cols = parseLine(lines[i])
    if (cols.length <= nameIdx || !cols[nameIdx]) continue
    items.push({
      name: cols[nameIdx] || '',
      slug: slugIdx !== -1 && cols[slugIdx] ? cols[slugIdx] : '',
      client_name: clientIdx !== -1 && cols[clientIdx] ? cols[clientIdx] : '',
      status: statusIdx !== -1 && cols[statusIdx] ? cols[statusIdx] : 'active',
      deadline: deadlineIdx !== -1 && cols[deadlineIdx] ? cols[deadlineIdx] : '',
      color_dot: colorIdx !== -1 && cols[colorIdx] ? cols[colorIdx] : '#10B981'
    })
  }
  return items
}

function processSelectedFile(file: File) {
  selectedImportFile.value = file
  importErrorMessage.value = ''
  importSuccessMessage.value = ''

  const reader = new FileReader()
  reader.onload = (evt) => {
    try {
      const content = evt.target?.result as string
      const parsed = parseCSVClientSide(content)
      parsedPreviewRows.value = parsed
      if (parsed.length === 0) {
        importErrorMessage.value = 'Tidak ada baris data proyek yang valid terdeteksi dalam file CSV.'
      }
    } catch (err: any) {
      importErrorMessage.value = 'Gagal memproses file: ' + err.message
    }
  }
  reader.readAsText(file)
}

async function executeImport() {
  if (!selectedImportFile.value) return
  isImporting.value = true
  importErrorMessage.value = ''
  importSuccessMessage.value = ''

  try {
    const res = await projectApi.importProjects(selectedImportFile.value)
    importSuccessMessage.value = `Berhasil! ${res.created_count} proyek baru dibuat, ${res.updated_count} proyek diperbarui.`
    store.showToast(importSuccessMessage.value)
    await store.fetchProjects()
    setTimeout(() => {
      closeImportModal()
    }, 1600)
  } catch (err: any) {
    importErrorMessage.value = err.response?.data?.error || 'Gagal mengimpor proyek. Periksa kembali format file.'
  } finally {
    isImporting.value = false
  }
}
</script>
