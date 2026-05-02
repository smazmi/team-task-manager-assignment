import type { AxiosResponse } from 'axios'
import { apiClient } from './axios'
import type {
  AddProjectMemberInput,
  ApiResponse,
  AuthSession,
  CreateProjectInput,
  CreateTaskInput,
  DashboardStats,
  LoginInput,
  Project,
  RegisterInput,
  Task,
  UpdateTaskInput,
  UpdateTaskStatusInput,
} from '../types'

function unwrapResponse<T>(response: AxiosResponse<ApiResponse<T>>) {
  return response.data.data
}

export async function loginRequest(payload: LoginInput) {
  const response = await apiClient.post<ApiResponse<AuthSession>>('/auth/login', payload)
  return unwrapResponse(response)
}

export async function registerRequest(payload: RegisterInput) {
  const response = await apiClient.post<ApiResponse<AuthSession>>(
    '/auth/register',
    payload,
  )
  return unwrapResponse(response)
}

export async function fetchProjects() {
  const response = await apiClient.get<ApiResponse<Project[]>>('/projects')
  return unwrapResponse(response)
}

export async function fetchProject(projectId: number) {
  const response = await apiClient.get<ApiResponse<Project>>(`/projects/${projectId}`)
  return unwrapResponse(response)
}

export async function createProject(payload: CreateProjectInput) {
  const response = await apiClient.post<ApiResponse<Project>>('/projects', payload)
  return unwrapResponse(response)
}

export async function addProjectMember(
  projectId: number,
  payload: AddProjectMemberInput,
) {
  const response = await apiClient.post<ApiResponse<Project>>(
    `/projects/${projectId}/members`,
    payload,
  )
  return unwrapResponse(response)
}

export async function removeProjectMember(projectId: number, userId: number) {
  await apiClient.delete(`/projects/${projectId}/members/${userId}`)
}

export async function fetchTasks(projectId: number) {
  const response = await apiClient.get<ApiResponse<Task[]>>('/tasks', {
    params: { project_id: projectId },
  })
  return unwrapResponse(response)
}

export async function createTask(payload: CreateTaskInput) {
  const response = await apiClient.post<ApiResponse<Task>>('/tasks', payload)
  return unwrapResponse(response)
}

export async function updateTask(taskId: number, payload: UpdateTaskInput) {
  const response = await apiClient.patch<ApiResponse<Task>>(`/tasks/${taskId}`, payload)
  return unwrapResponse(response)
}

export async function updateTaskStatus(taskId: number, payload: UpdateTaskStatusInput) {
  const response = await apiClient.patch<ApiResponse<Task>>(
    `/tasks/${taskId}/status`,
    payload,
  )
  return unwrapResponse(response)
}

export async function deleteTask(taskId: number) {
  await apiClient.delete(`/tasks/${taskId}`)
}

export async function fetchDashboardStats(projectId: number) {
  const response = await apiClient.get<ApiResponse<DashboardStats>>(
    '/dashboard/stats',
    {
      params: { project_id: projectId },
    },
  )
  return unwrapResponse(response)
}
