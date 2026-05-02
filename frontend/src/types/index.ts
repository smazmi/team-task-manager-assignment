export type ProjectRole = 'admin' | 'member'
export type TaskPriority = 'low' | 'medium' | 'high'
export type TaskStatus = 'todo' | 'in_progress' | 'done'

export interface ApiResponse<T> {
  message?: string
  data: T
}

export interface ApiErrorResponse {
  error: {
    code: string
    message: string
  }
}

export interface User {
  id: number
  name: string
  email: string
  created_at: string
  updated_at: string
}

export interface AuthSession {
  token: string
  user: User
}

export interface ProjectMember {
  id: number
  project_id: number
  user_id: number
  role: ProjectRole
  user: User
  created_at: string
  updated_at: string
}

export interface Project {
  id: number
  name: string
  description: string
  creator_id: number
  creator: User
  members: ProjectMember[]
  created_at: string
  updated_at: string
}

export interface Task {
  id: number
  project_id: number
  creator_id: number
  assignee_id?: number | null
  title: string
  description: string
  due_date?: string | null
  priority: TaskPriority
  status: TaskStatus
  creator: User
  assignee?: User | null
  created_at: string
  updated_at: string
}

export interface TasksPerUser {
  user_id: number
  user_name: string
  count: number
}

export interface DashboardStats {
  project_id: number
  total_tasks: number
  overdue_tasks: number
  tasks_by_status: Record<TaskStatus, number>
  tasks_per_user: TasksPerUser[]
}

export interface LoginInput {
  email: string
  password: string
}

export interface RegisterInput {
  name: string
  email: string
  password: string
}

export interface CreateProjectInput {
  name: string
  description: string
}

export interface AddProjectMemberInput {
  email: string
  role: ProjectRole
}

export interface CreateTaskInput {
  project_id: number
  title: string
  description: string
  due_date?: string
  priority: TaskPriority
  assignee_id?: number
}

export interface UpdateTaskInput {
  title?: string
  description?: string
  due_date?: string
  priority?: TaskPriority
  assignee_id?: number
}

export interface UpdateTaskStatusInput {
  status: TaskStatus
}

export const TASK_STATUS_OPTIONS: Array<{ value: TaskStatus; label: string }> = [
  { value: 'todo', label: 'To Do' },
  { value: 'in_progress', label: 'In Progress' },
  { value: 'done', label: 'Done' },
]

export const TASK_PRIORITY_OPTIONS: Array<{
  value: TaskPriority
  label: string
}> = [
  { value: 'low', label: 'Low' },
  { value: 'medium', label: 'Medium' },
  { value: 'high', label: 'High' },
]

export const PROJECT_ROLE_OPTIONS: Array<{
  value: ProjectRole
  label: string
}> = [
  { value: 'admin', label: 'Admin' },
  { value: 'member', label: 'Member' },
]
