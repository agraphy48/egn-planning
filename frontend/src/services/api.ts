import axios, { type AxiosInstance } from 'axios'
import type {
  Project,
  Task,
  Milestone,
  KPIStats,
  ActivityItem,
  TaskStatus,
  Attachment,
  WikiPage,
  TimeEntry,
  TeamWorkload,
  NotificationItem,
  BudgetItem
} from '../types/project'
import {
  mockProjects,
  mockTasks,
  mockMilestones,
  mockKPIStats,
  mockActivityFeed,
  teamMembers
} from '../data/mockData'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '/api/v1'
export const USE_MOCK = import.meta.env.VITE_USE_MOCK === 'true'

export const apiClient: AxiosInstance = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
  timeout: 10000,
})

// Request Interceptor: Attach JWT Bearer Token for Golang Backend
apiClient.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('auth_token')
    if (token && config.headers) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => Promise.reject(error)
)

// Response Interceptor: Global Error Logging & 401 handling
apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      const code = error.response?.data?.code || ''
      console.warn('[API] 401 Unauthorized - invalid or revoked session, code:', code)
      localStorage.removeItem('auth_token')
      localStorage.removeItem('auth_user')

      if (typeof window !== 'undefined' && window.location.pathname !== '/login') {
        let reason = 'session_expired'
        if (code === 'SESSION_CONCURRENT_KICKED') {
          reason = 'concurrent_device'
        } else if (code === 'SESSION_IDLE_TIMEOUT') {
          reason = 'idle_timeout'
        }
        window.location.href = `/login?reason=${reason}`
      }
    }
    return Promise.reject(error)
  }
)

// In-memory mutable store for offline fallback mode
let inMemoryTasks: Task[] = JSON.parse(JSON.stringify(mockTasks))
let inMemoryProjects: Project[] = JSON.parse(JSON.stringify(mockProjects))
let inMemoryMilestones: Milestone[] = JSON.parse(JSON.stringify(mockMilestones))

const mockDelay = <T>(data: T, ms = 80): Promise<T> =>
  new Promise((resolve) => setTimeout(() => resolve(data), ms))

export const projectApi = {
  // 1. Projects
  async getProjects(): Promise<Project[]> {
    if (USE_MOCK) return mockDelay([...inMemoryProjects])
    try {
      const response = await apiClient.get<Project[]>('/projects')
      return response.data
    } catch (err: any) {
      if (err?.response?.status === 401 || err?.response?.status === 403 || import.meta.env.PROD) {
        throw err
      }
      return mockDelay([...inMemoryProjects])
    }
  },

  async getProject(projectId: string): Promise<Project> {
    if (USE_MOCK) {
      const proj = inMemoryProjects.find((p) => p.id === projectId)
      if (!proj) throw new Error('Project not found')
      return mockDelay({ ...proj })
    }
    try {
      const response = await apiClient.get<Project>(`/projects/${projectId}`)
      return response.data
    } catch (err: any) {
      if (err?.response?.status === 401 || err?.response?.status === 403 || import.meta.env.PROD) {
        throw err
      }
      const proj = inMemoryProjects.find((p) => p.id === projectId)
      return mockDelay({ ...(proj || inMemoryProjects[0]) })
    }
  },

  async createProject(projectData: Partial<Project>): Promise<Project> {
    if (USE_MOCK) {
      const newProj: Project = {
        id: `proj_${Date.now()}`,
        name: projectData.name || 'New Project',
        slug: projectData.slug || `proj-${Date.now()}`,
        client_name: projectData.client_name || 'Client Stakeholder',
        status: 'active',
        progress_percentage: 0,
        color_dot: projectData.color_dot || '#4F46E5',
        deadline: projectData.deadline || '2026-12-31',
        total_tasks: 0,
        completed_tasks: 0,
        created_at: new Date().toISOString()
      }
      inMemoryProjects.unshift(newProj)
      return mockDelay(newProj)
    }
    const response = await apiClient.post<Project>('/projects', projectData)
    return response.data
  },

  // 2. Dashboard KPIs & Activities
  async getKPIStats(): Promise<KPIStats> {
    if (USE_MOCK) return mockDelay({ ...mockKPIStats })
    try {
      const response = await apiClient.get<KPIStats>('/dashboard/kpis')
      return response.data
    } catch {
      return mockDelay({ ...mockKPIStats })
    }
  },

  async getActivityFeed(projectId?: string): Promise<ActivityItem[]> {
    if (USE_MOCK) return mockDelay([...mockActivityFeed])
    try {
      const response = await apiClient.get<ActivityItem[]>('/dashboard/activity', { params: { project_id: projectId } })
      return response.data
    } catch {
      return mockDelay([...mockActivityFeed])
    }
  },

  // 3. Milestones & Approvals
  async getMilestones(projectId?: string): Promise<Milestone[]> {
    if (USE_MOCK) {
      if (projectId) return mockDelay(inMemoryMilestones.filter((m) => m.project_id === projectId))
      return mockDelay([...inMemoryMilestones])
    }
    try {
      const response = await apiClient.get<Milestone[]>('/milestones', { params: { project_id: projectId } })
      return response.data
    } catch {
      return mockDelay([...inMemoryMilestones])
    }
  },

  async approveMilestone(milestoneId: string, clientName: string): Promise<Milestone> {
    if (USE_MOCK) {
      const ms = inMemoryMilestones.find((m) => m.id === milestoneId)
      if (!ms) throw new Error('Milestone not found')
      ms.status = 'approved'
      ms.approved_by = clientName
      ms.approved_at = new Date().toISOString()
      return mockDelay({ ...ms })
    }
    // CONTRACT FIX: Use PATCH with fallback
    const response = await apiClient.patch<any>(`/milestones/${milestoneId}/approve`, { client_name: clientName })
    return response.data.milestone || response.data
  },

  async rejectMilestone(milestoneId: string, reason: string, clientName: string): Promise<Milestone> {
    if (USE_MOCK) {
      const ms = inMemoryMilestones.find((m) => m.id === milestoneId)
      if (!ms) throw new Error('Milestone not found')
      ms.status = 'revision'
      ms.feedback_notes = reason
      return mockDelay({ ...ms })
    }
    const response = await apiClient.patch<any>(`/milestones/${milestoneId}/reject`, { reason, client_name: clientName })
    return response.data.milestone || response.data
  },

  // 4. Tasks CRUD & Multi-View
  async getTasks(projectId?: string): Promise<Task[]> {
    if (USE_MOCK) {
      if (projectId) return mockDelay(inMemoryTasks.filter((t) => t.project_id === projectId))
      return mockDelay([...inMemoryTasks])
    }
    try {
      const response = await apiClient.get<Task[]>('/tasks', { params: { project_id: projectId } })
      return response.data
    } catch {
      return mockDelay([...inMemoryTasks])
    }
  },

  async updateTaskStatus(taskId: string, status: TaskStatus): Promise<Task> {
    if (USE_MOCK) {
      const index = inMemoryTasks.findIndex((t) => t.id === taskId)
      if (index === -1) throw new Error('Task not found')
      inMemoryTasks[index].status = status
      inMemoryTasks[index].updated_at = new Date().toISOString()
      return mockDelay({ ...inMemoryTasks[index] })
    }
    const response = await apiClient.patch<Task>(`/tasks/${taskId}/status`, { status })
    return response.data
  },

  async updateTask(taskId: string, updates: Partial<Task>): Promise<Task> {
    if (USE_MOCK) {
      const index = inMemoryTasks.findIndex((t) => t.id === taskId)
      if (index === -1) throw new Error('Task not found')
      inMemoryTasks[index] = { ...inMemoryTasks[index], ...updates, updated_at: new Date().toISOString() }
      return mockDelay({ ...inMemoryTasks[index] })
    }
    const response = await apiClient.patch<Task>(`/tasks/${taskId}`, updates)
    return response.data
  },

  async createTask(taskData: Partial<Task>): Promise<Task> {
    if (USE_MOCK) {
      const newTask: Task = {
        id: `task_${Date.now()}`,
        project_id: taskData.project_id || 'proj_sustagri',
        key: `#EGN-${Math.floor(100 + Math.random() * 900)}`,
        title: taskData.title || 'Untitled Task',
        description: taskData.description || '',
        category: taskData.category || 'General',
        priority: taskData.priority || 'Medium',
        status: taskData.status || 'Backlog',
        assignee: taskData.assignee || teamMembers[0],
        due_date: taskData.due_date || '2026-11-30',
        subtasks: taskData.subtasks || [],
        comments: [],
        attachments: [],
        client_visible: taskData.client_visible ?? true,
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString()
      }
      inMemoryTasks.unshift(newTask)
      return mockDelay(newTask)
    }
    const response = await apiClient.post<Task>('/tasks', taskData)
    return response.data
  },

  async deleteTask(taskId: string): Promise<boolean> {
    if (USE_MOCK) {
      const idx = inMemoryTasks.findIndex((t) => t.id === taskId)
      if (idx !== -1) inMemoryTasks.splice(idx, 1)
      return mockDelay(true)
    }
    await apiClient.delete(`/tasks/${taskId}`)
    return true
  },

  async toggleSubtask(taskId: string, subtaskId: string, isDone: boolean): Promise<Task> {
    if (USE_MOCK) {
      const task = inMemoryTasks.find((t) => t.id === taskId)
      if (!task) throw new Error('Task not found')
      const sub = task.subtasks.find((s) => s.id === subtaskId)
      if (sub) sub.is_done = isDone
      return mockDelay({ ...task })
    }
    const response = await apiClient.patch<Task>(`/tasks/${taskId}/subtasks/${subtaskId}`, { is_done: isDone })
    return response.data
  },

  async addComment(taskId: string, content: string, clientVisible: boolean = true): Promise<Task> {
    if (USE_MOCK) {
      const task = inMemoryTasks.find((t) => t.id === taskId)
      if (!task) throw new Error('Task not found')
      task.comments.push({
        id: `com_${Date.now()}`,
        task_id: taskId,
        user: teamMembers[0],
        content,
        client_visible: clientVisible,
        created_at: 'Just now'
      })
      return mockDelay({ ...task })
    }
    const response = await apiClient.post<Task>(`/tasks/${taskId}/comments`, {
      content,
      client_visible: clientVisible
    })
    return response.data
  },

  // 5. Attachments
  async uploadAttachment(taskId: string, file: File): Promise<Attachment> {
    const formData = new FormData()
    formData.append('file', file)
    const response = await apiClient.post<Attachment>(`/tasks/${taskId}/attachments`, formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
    return response.data
  },

  async deleteAttachment(attachmentId: string): Promise<boolean> {
    await apiClient.delete(`/tasks/attachments/${attachmentId}`)
    return true
  },

  // 6. Wiki / Knowledge Base
  async getWikiPages(projectId: string): Promise<WikiPage[]> {
    const response = await apiClient.get<WikiPage[]>(`/projects/${projectId}/wiki`)
    return response.data
  },

  async createWikiPage(projectId: string, pageData: { title: string; content: string; client_visible: boolean }): Promise<WikiPage> {
    const response = await apiClient.post<WikiPage>(`/projects/${projectId}/wiki`, pageData)
    return response.data
  },

  // 7. Timesheets & Workload
  async logTime(data: { project_id: string; task_id?: string; hours: number; description: string; date?: string }): Promise<TimeEntry> {
    const response = await apiClient.post<TimeEntry>('/timesheets', data)
    return response.data
  },

  async getTimeEntries(projectId?: string): Promise<TimeEntry[]> {
    const response = await apiClient.get<TimeEntry[]>('/timesheets', { params: { project_id: projectId } })
    return response.data
  },

  async getTeamWorkload(projectId?: string): Promise<TeamWorkload[]> {
    const response = await apiClient.get<TeamWorkload[]>('/timesheets/workload', { params: { project_id: projectId } })
    return response.data
  },

  // 8. In-App Notifications
  async getNotifications(): Promise<{ notifications: NotificationItem[]; unread_count: number }> {
    try {
      const response = await apiClient.get<{ notifications: NotificationItem[]; unread_count: number }>('/notifications')
      return response.data
    } catch {
      return { notifications: [], unread_count: 0 }
    }
  },

  async markNotificationRead(id: string): Promise<void> {
    await apiClient.patch(`/notifications/${id}/read`)
  },

  async markAllNotificationsRead(): Promise<void> {
    await apiClient.post('/notifications/read-all')
  },

  // 9. Budgets
  async getBudgets(projectId: string): Promise<{ budgets: BudgetItem[]; total_allocated: number; total_spent: number; currency: string }> {
    const response = await apiClient.get<any>(`/projects/${projectId}/budgets`)
    return response.data
  },
}
