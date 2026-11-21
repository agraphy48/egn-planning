import { defineStore } from 'pinia'
import type { Project, Task, Milestone, KPIStats, ActivityItem, TaskStatus, StatusConfig } from '../types/project'
import { projectApi, statusApi } from '../services/api'

export const useProjectStore = defineStore('project', {
  state: () => ({
    projects: [] as Project[],
    currentProjectId: 'proj_sustagri',
    tasks: [] as Task[],
    milestones: [] as Milestone[],
    kpiStats: null as KPIStats | null,
    activityFeed: [] as ActivityItem[],
    statuses: [] as StatusConfig[],
    
    // UI & Filter state
    currentView: 'kanban' as 'kanban' | 'list' | 'timeline',
    collapsedGroups: [] as string[],
    activeFilter: 'All',
    searchQuery: '',
    selectedTaskId: null as string | null,
    isDrawerOpen: false,
    isDarkMode: false,
    isSidebarCollapsed: false,
    isMobileSidebarOpen: false,
    isNewTaskModalOpen: false,
    isCommandPaletteOpen: false,
    
    // User Profile state
    currentUser: {
      id: '',
      name: 'Guest User',
      email: '',
      role: 'member' as any,
      avatar_url: '',
      initials: 'GU',
    },
    activeRole: 'project_manager' as 'project_manager' | 'team_member' | 'client',
    
    // Loading & Notification states
    isLoading: false,
    toastMessage: '' as string | null,
  }),

  getters: {
    currentProject(state): Project | undefined {
      return state.projects.find((p) => p.id === state.currentProjectId) || state.projects[0]
    },

    filteredTasks(state): Task[] {
      return state.tasks.filter((task) => {
        const matchesProject = task.project_id === state.currentProjectId
        const matchesCategory =
          state.activeFilter === 'All' ||
          task.category.toLowerCase() === state.activeFilter.toLowerCase()
        const matchesSearch =
          state.searchQuery === '' ||
          task.title.toLowerCase().includes(state.searchQuery.toLowerCase()) ||
          task.key.toLowerCase().includes(state.searchQuery.toLowerCase())
        return matchesProject && matchesCategory && matchesSearch
      })
    },

    taskStatuses(state): StatusConfig[] {
      return state.statuses
        .filter((s) => s.category === 'task')
        .sort((a, b) => a.position - b.position)
    },

    tasksByStatus(): Record<string, Task[]> {
      const grouped: Record<string, Task[]> = {}
      if (this.taskStatuses && this.taskStatuses.length > 0) {
        this.taskStatuses.forEach((s) => {
          grouped[s.name] = []
        })
      } else {
        grouped['Backlog'] = []
        grouped['In Progress'] = []
        grouped['Review'] = []
        grouped['Done'] = []
      }

      this.filteredTasks.forEach((task) => {
        if (!grouped[task.status]) {
          grouped[task.status] = []
        }
        grouped[task.status].push(task)
      })
      return grouped
    },

    taskCounts(): Record<string, number> {
      const counts: Record<string, number> = {}
      for (const key in this.tasksByStatus) {
        counts[key] = this.tasksByStatus[key].length
      }
      return counts
    },

    selectedTask(state): Task | undefined {
      if (!state.selectedTaskId) return undefined
      return state.tasks.find((t) => t.id === state.selectedTaskId)
    },
  },

  actions: {
    async fetchInitialData() {
      this.isLoading = true
      try {
        const [projects, tasks, kpis, milestones, activities, statuses] = await Promise.all([
          projectApi.getProjects(),
          projectApi.getTasks(),
          projectApi.getKPIStats(),
          projectApi.getMilestones(),
          projectApi.getActivityFeed(),
          statusApi.getStatuses().catch(() => []),
        ])
        this.projects = projects
        this.tasks = tasks
        this.kpiStats = kpis
        this.milestones = milestones
        this.activityFeed = activities
        if (statuses && statuses.length > 0) {
          this.statuses = statuses
        }
      } catch (err) {
        console.error('Failed to load initial project data:', err)
      } finally {
        this.isLoading = false
      }
    },

    async fetchStatuses() {
      try {
        const data = await statusApi.getStatuses()
        this.statuses = data
      } catch (err) {
        console.error('Failed to load statuses:', err)
      }
    },

    async fetchProjects() {
      try {
        const projects = await projectApi.getProjects()
        this.projects = projects
      } catch (err) {
        console.error('Failed to load projects:', err)
      }
    },

    selectProject(projectId: string) {
      this.currentProjectId = projectId
    },

    async updateTaskStatus(taskId: string, newStatus: TaskStatus) {
      // Optimistic Update
      const taskIndex = this.tasks.findIndex((t) => t.id === taskId)
      if (taskIndex === -1) return
      
      const previousStatus = this.tasks[taskIndex].status
      this.tasks[taskIndex].status = newStatus

      this.showToast(`Task moved to ${newStatus}`)

      try {
        await projectApi.updateTaskStatus(taskId, newStatus)
      } catch (err) {
        // Rollback on failure
        this.tasks[taskIndex].status = previousStatus
        this.showToast('Failed to update task status')
        console.error(err)
      }
    },

    setView(view: 'kanban' | 'list' | 'timeline') {
      this.currentView = view
    },

    toggleGroupCollapse(group: string) {
      const idx = this.collapsedGroups.indexOf(group)
      if (idx === -1) {
        this.collapsedGroups.push(group)
      } else {
        this.collapsedGroups.splice(idx, 1)
      }
    },

    quickAddTask(status: TaskStatus, title: string) {
      if (!title.trim()) return
      const keyNum = Math.floor(100 + Math.random() * 900)
      const newTask: Task = {
        id: `task_${Date.now()}`,
        project_id: this.currentProjectId,
        key: `#EGN-${keyNum}`,
        title: title.trim(),
        description: 'Quick added task. Click to edit description and subtasks.',
        category: this.activeFilter !== 'All' ? this.activeFilter : 'Engineering',
        priority: 'Medium',
        status: status,
        assignee: this.currentUser,
        due_date: new Date(Date.now() + 7 * 86400000).toISOString().split('T')[0],
        subtasks: [],
        comments: [],
        attachments: [],
        client_visible: true,
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString()
      }
      this.tasks.push(newTask)
      this.showToast(`Task ${newTask.key} created in ${status}`)
      try {
        projectApi.createTask(newTask)
      } catch (err) {
        console.error(err)
      }
    },

    updateTask(taskId: string, updates: Partial<Task>) {
      const idx = this.tasks.findIndex((t) => t.id === taskId)
      if (idx !== -1) {
        this.tasks[idx] = { ...this.tasks[idx], ...updates, updated_at: new Date().toISOString() }
        this.showToast(`Task ${this.tasks[idx].key} updated`)
      }
    },

    openTaskDetail(taskId: string) {
      this.selectedTaskId = taskId
      this.isDrawerOpen = true
    },

    closeTaskDetail() {
      this.isDrawerOpen = false
      setTimeout(() => {
        this.selectedTaskId = null
      }, 250)
    },

    async toggleSubtask(taskId: string, subtaskId: string, isDone: boolean) {
      const task = this.tasks.find((t) => t.id === taskId)
      if (!task) return
      const subtask = task.subtasks.find((s) => s.id === subtaskId)
      if (subtask) {
        subtask.is_done = isDone
      }
      try {
        await projectApi.toggleSubtask(taskId, subtaskId, isDone)
      } catch (err) {
        console.error('Failed to toggle subtask:', err)
      }
    },

    addSubtask(taskId: string, title: string) {
      const task = this.tasks.find((t) => t.id === taskId)
      if (!task || !title.trim()) return
      const newSubtask = {
        id: `sub_${Date.now()}`,
        task_id: taskId,
        title: title.trim(),
        is_done: false,
      }
      task.subtasks.push(newSubtask)
      this.showToast('Subtask added')
    },

    async addComment(taskId: string, content: string) {
      if (!content.trim()) return
      const task = this.tasks.find((t) => t.id === taskId)
      if (!task) return

      const newComment = {
        id: `com_${Date.now()}`,
        task_id: taskId,
        user: this.currentUser,
        content: content.trim(),
        created_at: 'Just now',
      }
      task.comments.push(newComment)
      this.showToast('Comment posted')

      try {
        await projectApi.addComment(taskId, content)
      } catch (err) {
        console.error('Failed to persist comment:', err)
      }
    },

    async createTask(taskData: Partial<Task>) {
      const created = await projectApi.createTask({
        ...taskData,
        project_id: this.currentProjectId,
      })
      this.tasks.unshift(created)
      this.isNewTaskModalOpen = false
      this.showToast(`Task ${created.key} created successfully`)
    },

    async deleteTask(taskId: string) {
      const idx = this.tasks.findIndex((t) => t.id === taskId)
      if (idx !== -1) {
        const deleted = this.tasks.splice(idx, 1)[0]
        this.closeTaskDetail()
        this.showToast(`Task ${deleted.key} deleted`)
        try {
          await projectApi.deleteTask(taskId)
        } catch (err) {
          console.error('Failed to delete task on backend:', err)
        }
      }
    },

    async approveMilestone(milestoneId: string) {
      const ms = this.milestones.find((m) => m.id === milestoneId)
      if (ms) {
        ms.status = 'approved'
        ms.approved_by = this.currentUser.name
        ms.approved_at = new Date().toISOString()
      }
      this.showToast('Milestone approved & client notification sent!')
      try {
        await projectApi.approveMilestone(milestoneId, this.currentUser.name)
      } catch (err) {
        console.error('Failed to approve milestone:', err)
      }
    },

    toggleTheme() {
      this.isDarkMode = !this.isDarkMode
      if (this.isDarkMode) {
        document.documentElement.classList.add('dark')
      } else {
        document.documentElement.classList.remove('dark')
      }
    },

    toggleSidebar() {
      if (typeof window !== 'undefined' && window.innerWidth < 768) {
        this.isMobileSidebarOpen = !this.isMobileSidebarOpen
      } else {
        this.isSidebarCollapsed = !this.isSidebarCollapsed
      }
    },

    toggleMobileSidebar() {
      this.isMobileSidebarOpen = !this.isMobileSidebarOpen
    },

    closeMobileSidebar() {
      this.isMobileSidebarOpen = false
    },

    setFilter(filter: string) {
      this.activeFilter = filter
    },

    setSearchQuery(query: string) {
      this.searchQuery = query
    },

    switchRole(role: 'project_manager' | 'team_member' | 'client') {
      this.activeRole = role
      this.showToast(`Switched view context to ${role.replace('_', ' ').toUpperCase()}`)
    },

    resetStore() {
      this.projects = []
      this.tasks = []
      this.milestones = []
      this.kpiStats = null
      this.activityFeed = []
      this.selectedTaskId = null
      this.isDrawerOpen = false
      this.isNewTaskModalOpen = false
      this.isCommandPaletteOpen = false
      this.currentUser = {
        id: '',
        name: 'Guest User',
        email: '',
        role: 'member' as any,
        avatar_url: '',
        initials: 'GU',
      }
    },

    showToast(message: string) {
      this.toastMessage = message
      setTimeout(() => {
        if (this.toastMessage === message) {
          this.toastMessage = null
        }
      }, 3000)
    },
  },
})
