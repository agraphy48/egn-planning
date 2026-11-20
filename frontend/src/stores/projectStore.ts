import { defineStore } from 'pinia'
import type { Project, Task, Milestone, KPIStats, ActivityItem, TaskStatus } from '../types/project'
import { projectApi } from '../services/api'
import { currentUser } from '../data/mockData'

export const useProjectStore = defineStore('project', {
  state: () => ({
    projects: [] as Project[],
    currentProjectId: 'proj_sustagri',
    tasks: [] as Task[],
    milestones: [] as Milestone[],
    kpiStats: null as KPIStats | null,
    activityFeed: [] as ActivityItem[],
    
    // UI & Filter state
    currentView: 'kanban' as 'kanban' | 'list' | 'timeline',
    collapsedGroups: [] as string[],
    activeFilter: 'All',
    searchQuery: '',
    selectedTaskId: null as string | null,
    isDrawerOpen: false,
    isDarkMode: false,
    isSidebarCollapsed: false,
    isNewTaskModalOpen: false,
    isCommandPaletteOpen: false,
    
    // User Role switcher (PM / Team Member / Client)
    currentUser: currentUser,
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

    tasksByStatus(): Record<TaskStatus, Task[]> {
      const grouped: Record<TaskStatus, Task[]> = {
        'Backlog': [],
        'In Progress': [],
        'Review': [],
        'Done': [],
      }
      this.filteredTasks.forEach((task) => {
        if (grouped[task.status]) {
          grouped[task.status].push(task)
        }
      })
      return grouped
    },

    taskCounts(): Record<TaskStatus, number> {
      return {
        'Backlog': this.tasksByStatus['Backlog'].length,
        'In Progress': this.tasksByStatus['In Progress'].length,
        'Review': this.tasksByStatus['Review'].length,
        'Done': this.tasksByStatus['Done'].length,
      }
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
        const [projects, tasks, kpis, milestones, activities] = await Promise.all([
          projectApi.getProjects(),
          projectApi.getTasks(),
          projectApi.getKPIStats(),
          projectApi.getMilestones(),
          projectApi.getActivityFeed(),
        ])
        this.projects = projects
        this.tasks = tasks
        this.kpiStats = kpis
        this.milestones = milestones
        this.activityFeed = activities
      } catch (err) {
        console.error('Failed to load initial project data:', err)
      } finally {
        this.isLoading = false
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
      this.isSidebarCollapsed = !this.isSidebarCollapsed
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
