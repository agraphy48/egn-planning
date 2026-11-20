export type UserRole = 'super_admin' | 'pm' | 'member' | 'client' | 'project_manager' | 'team_member'

export interface User {
  id: string
  name: string
  email: string
  avatar_url: string
  role: UserRole
  initials?: string
}

export type ProjectStatus = 'active' | 'completed' | 'on_hold'

export interface Project {
  id: string
  name: string
  slug: string
  client_name?: string
  status: ProjectStatus
  progress_percentage: number
  color_dot: string
  deadline: string
  total_tasks: number
  completed_tasks: number
  created_at: string
}

export type MilestoneStatus = 'pending' | 'in_review' | 'approved'

export interface Milestone {
  id: string
  project_id: string
  title: string
  description: string
  due_date: string
  status: MilestoneStatus
  approved_by?: string
  approved_at?: string
}

export interface SubTask {
  id: string
  task_id: string
  title: string
  is_done: boolean
  created_at?: string
}

export interface Comment {
  id: string
  task_id: string
  user: User
  content: string
  created_at: string
}

export interface Attachment {
  id: string
  task_id: string
  file_name: string
  file_size: string
  file_type: string
  file_url: string
  uploaded_at: string
}

export type TaskPriority = 'Urgent' | 'High' | 'Medium' | 'Low'
export type TaskStatus = 'Backlog' | 'In Progress' | 'Review' | 'Done'

export interface Task {
  id: string
  project_id: string
  key: string
  title: string
  description: string
  category: string
  priority: TaskPriority
  status: TaskStatus
  assignee: User
  due_date: string
  start_date?: string
  time_estimate?: string
  fase?: string
  epic_modul?: string
  deliverable?: string
  tags?: string[]
  subtasks: SubTask[]
  comments: Comment[]
  attachments: Attachment[]
  client_visible: boolean
  created_at: string
  updated_at: string
}

export interface KPIStats {
  active_projects: number
  tasks_due_this_week: number
  overdue_tasks: number
  completion_rate: number
  active_projects_delta: string
  overdue_tasks_delta: string
  completion_rate_delta: string
}

export interface ActivityItem {
  id: string
  user: User
  action: string
  target: string
  timestamp: string
}
