import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import {
  createTask,
  deleteTask,
  fetchTasks,
  updateTask,
  updateTaskStatus,
} from '../api/taskManagerApi'
import type { CreateTaskInput, UpdateTaskInput, UpdateTaskStatusInput } from '../types'

export function useTasks(projectId: number) {
  return useQuery({
    queryKey: ['tasks', projectId],
    queryFn: () => fetchTasks(projectId),
    enabled: projectId > 0,
  })
}

export function useCreateTask() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (payload: CreateTaskInput) => createTask(payload),
    onSuccess: (task) => {
      void queryClient.invalidateQueries({ queryKey: ['tasks', task.project_id] })
      void queryClient.invalidateQueries({
        queryKey: ['dashboard', task.project_id],
      })
    },
  })
}

export function useUpdateTaskStatus() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({
      taskId,
      payload,
    }: {
      taskId: number
      payload: UpdateTaskStatusInput
    }) => updateTaskStatus(taskId, payload),
    onSuccess: (task) => {
      void queryClient.invalidateQueries({ queryKey: ['tasks', task.project_id] })
      void queryClient.invalidateQueries({
        queryKey: ['dashboard', task.project_id],
      })
    },
  })
}

export function useUpdateTask() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({
      taskId,
      payload,
    }: {
      taskId: number
      payload: UpdateTaskInput
    }) => updateTask(taskId, payload),
    onSuccess: (task) => {
      void queryClient.invalidateQueries({ queryKey: ['tasks', task.project_id] })
      void queryClient.invalidateQueries({
        queryKey: ['dashboard', task.project_id],
      })
    },
  })
}

export function useDeleteTask() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ taskId }: { taskId: number; projectId: number }) => deleteTask(taskId),
    onSuccess: (_data, variables) => {
      void queryClient.invalidateQueries({ queryKey: ['tasks', variables.projectId] })
      void queryClient.invalidateQueries({
        queryKey: ['dashboard', variables.projectId],
      })
    },
  })
}
