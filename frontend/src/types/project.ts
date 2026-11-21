export type UserRole = 'super_admin' | 'pm' | 'member' | 'client' | 'project_manager' | 'team_member'

export interface User {
  id: string
  name: string
  email: string
  avatar_url: string
  role: UserRole
  initials?: string
  current_session_id?: string
  last_active_at?: string | Date
  created_at?: string | Date
  updated_at?: string | Date
}

export type ProjectStatus = 'active' | 'completed' | 'on_hold' | 'archived'

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

export type MilestoneStatus = 'pending' | 'in_review' | 'approved' | 'revision' | 'rejected'

export interface Milestone {
  id: string
  project_id: string
  title: string
  description: string
  due_date: string
  status: MilestoneStatus
  approved_by?: string
  approved_at?: string
  feedback_notes?: string
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
  client_visible?: boolean
  created_at: string
}

export interface Attachment {
  id: string
  task_id: string
  file_name: string
  file_size: string
  file_type: string
  file_url: string
  uploaded_by?: string
  uploaded_at?: string
  created_at?: string
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
  active_projects_delta?: string
  tasks_due_this_week_delta?: string
  overdue_tasks_delta?: string
  completion_rate_delta?: string
  total_tasks?: number
  completed_tasks?: number
}

export interface ActivityItem {
  id: string
  user: any
  avatar_url?: string
  action: string
  target: string
  project?: string
  time_ago?: string
  timestamp?: string
}

export interface WikiPage {
  id: string
  project_id: string
  title: string
  slug: string
  content: string
  author_id?: string
  author_name: string
  client_visible: boolean
  created_at: string
  updated_at: string
}

export interface TimeEntry {
  id: string
  project_id: string
  task_id?: string
  user_id: string
  hours: number
  description: string
  date: string
  created_at: string
  user?: User
}

export interface TeamWorkload {
  user_id: string
  user_name: string
  user_avatar?: string
  total_hours: number
  tasks_count?: number
  task_count?: number
}

export interface NotificationItem {
  id: string
  user_id: string
  title: string
  message: string
  type: 'info' | 'success' | 'warning' | 'milestone'
  link?: string
  is_read: boolean
  created_at: string
}

export interface BudgetItem {
  id: string
  project_id: string
  category: string
  allocated: number
  spent: number
  currency: string
}

export interface StatusConfig {
  id: string
  name: string
  slug: string
  category: 'task' | 'project' | 'milestone'
  color_dot: string
  badge_bg: string
  badge_text: string
  position: number
  is_completed: boolean
  is_system: boolean
  description?: string
  usage_count?: number
  created_at?: string
  updated_at?: string
}
