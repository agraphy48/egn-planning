import axios, { type AxiosInstance } from 'axios'
import type { Project, Task, Milestone, KPIStats, ActivityItem, TaskStatus } from '../types/project'
import {
  mockProjects,
  mockTasks,
  mockMilestones,
  mockKPIStats,
  mockActivityFeed,
  teamMembers
} from '../data/mockData'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1'
export const USE_MOCK = import.meta.env.VITE_USE_MOCK !== 'false'

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
      console.warn('[API] Unauthorized access - redirecting to auth flow')
    }
    return Promise.reject(error)
  }
)

// In-memory mutable store for offline mock mode
let inMemoryTasks: Task[] = JSON.parse(JSON.stringify(mockTasks))
let inMemoryProjects: Project[] = JSON.parse(JSON.stringify(mockProjects))
let inMemoryMilestones: Milestone[] = JSON.parse(JSON.stringify(mockMilestones))

// Helper to simulate network latency in mock mode
const mockDelay = <T>(data: T, ms = 80): Promise<T> =>
  new Promise((resolve) => setTimeout(() => resolve(data), ms))

export const projectApi = {
  // Fetch all projects
  async getProjects(): Promise<Project[]> {
    if (USE_MOCK) {
      return mockDelay([...inMemoryProjects])
    }
    const response = await apiClient.get<Project[]>('/projects')
    return response.data
  },

  // Fetch single project by ID
  async getProject(projectId: string): Promise<Project> {
    if (USE_MOCK) {
      const proj = inMemoryProjects.find((p) => p.id === projectId)
      if (!proj) throw new Error('Project not found')
      return mockDelay({ ...proj })
    }
    const response = await apiClient.get<Project>(`/projects/${projectId}`)
    return response.data
  },

  // Fetch KPI Stats for Dashboard
  async getKPIStats(): Promise<KPIStats> {
    if (USE_MOCK) {
      return mockDelay({ ...mockKPIStats })
    }
    const response = await apiClient.get<KPIStats>('/dashboard/kpis')
    return response.data
  },

  // Fetch Activity Feed
  async getActivityFeed(): Promise<ActivityItem[]> {
    if (USE_MOCK) {
      return mockDelay([...mockActivityFeed])
    }
    const response = await apiClient.get<ActivityItem[]>('/dashboard/activity')
    return response.data
  },

  // Fetch Milestones
  async getMilestones(projectId?: string): Promise<Milestone[]> {
    if (USE_MOCK) {
      if (projectId) {
        return mockDelay(inMemoryMilestones.filter((m) => m.project_id === projectId))
      }
      return mockDelay([...inMemoryMilestones])
    }
    const response = await apiClient.get<Milestone[]>('/milestones', { params: { project_id: projectId } })
    return response.data
  },

  // Fetch Tasks
  async getTasks(projectId?: string): Promise<Task[]> {
    if (USE_MOCK) {
      if (projectId) {
        return mockDelay(inMemoryTasks.filter((t) => t.project_id === projectId))
      }
      return mockDelay([...inMemoryTasks])
    }
    const response = await apiClient.get<Task[]>('/tasks', { params: { project_id: projectId } })
    return response.data
  },

  // Update Task Status (Kanban Drag-and-Drop)
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

  // Create Task
  async createTask(taskData: Partial<Task>): Promise<Task> {
    if (USE_MOCK) {
      const newTask: Task = {
        id: `task_${Date.now()}`,
        project_id: taskData.project_id || 'proj_01',
        key: `#EGN-${Math.floor(100 + Math.random() * 900)}`,
        title: taskData.title || 'Untitled Task',
        description: taskData.description || '',
        category: taskData.category || 'General',
        priority: taskData.priority || 'Medium',
        status: taskData.status || 'Backlog',
        assignee: taskData.assignee || teamMembers[1],
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

  // Toggle Subtask
  async toggleSubtask(taskId: string, subtaskId: string, isDone: boolean): Promise<Task> {
    if (USE_MOCK) {
      const task = inMemoryTasks.find((t) => t.id === taskId)
      if (!task) throw new Error('Task not found')
      const sub = task.subtasks.find((s) => s.id === subtaskId)
      if (sub) {
        sub.is_done = isDone
      }
      return mockDelay({ ...task })
    }
    const response = await apiClient.patch<Task>(`/tasks/${taskId}/subtasks/${subtaskId}`, { is_done: isDone })
    return response.data
  },

  // Add Comment
  async addComment(taskId: string, content: string): Promise<Task> {
    if (USE_MOCK) {
      const task = inMemoryTasks.find((t) => t.id === taskId)
      if (!task) throw new Error('Task not found')
      task.comments.push({
        id: `com_${Date.now()}`,
        task_id: taskId,
        user: teamMembers[0], // current user
        content,
        created_at: 'Just now'
      })
      return mockDelay({ ...task })
    }
    const response = await apiClient.post<Task>(`/tasks/${taskId}/comments`, { content })
    return response.data
  },

  // Approve Milestone (Client Portal)
  async approveMilestone(milestoneId: string, clientName: string): Promise<Milestone> {
    if (USE_MOCK) {
      const ms = inMemoryMilestones.find((m) => m.id === milestoneId)
      if (!ms) throw new Error('Milestone not found')
      ms.status = 'approved'
      ms.approved_by = clientName
      ms.approved_at = new Date().toISOString()
      return mockDelay({ ...ms })
    }
    const response = await apiClient.post<Milestone>(`/milestones/${milestoneId}/approve`, { client_name: clientName })
    return response.data
  },

  // Delete Task
  async deleteTask(taskId: string): Promise<boolean> {
    if (USE_MOCK) {
      const idx = inMemoryTasks.findIndex((t) => t.id === taskId)
      if (idx !== -1) {
        inMemoryTasks.splice(idx, 1)
      }
      return mockDelay(true)
    }
    await apiClient.delete(`/tasks/${taskId}`)
    return true
  }
}
