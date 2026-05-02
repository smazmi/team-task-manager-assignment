import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import {
  addProjectMember,
  createProject,
  fetchProject,
  fetchProjects,
  removeProjectMember,
} from '../api/taskManagerApi'
import type { AddProjectMemberInput, CreateProjectInput } from '../types'

export function useProjects() {
  return useQuery({
    queryKey: ['projects'],
    queryFn: fetchProjects,
  })
}

export function useProject(projectId: number) {
  return useQuery({
    queryKey: ['projects', projectId],
    queryFn: () => fetchProject(projectId),
    enabled: projectId > 0,
  })
}

export function useCreateProject() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (payload: CreateProjectInput) => createProject(payload),
    onSuccess: (project) => {
      queryClient.setQueryData(['projects', project.id], project)
      void queryClient.invalidateQueries({ queryKey: ['projects'] })
    },
  })
}

export function useAddProjectMember() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({
      projectId,
      payload,
    }: {
      projectId: number
      payload: AddProjectMemberInput
    }) => addProjectMember(projectId, payload),
    onSuccess: (project) => {
      queryClient.setQueryData(['projects', project.id], project)
      void queryClient.invalidateQueries({ queryKey: ['projects'] })
    },
  })
}

export function useRemoveProjectMember() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({
      projectId,
      userId,
    }: {
      projectId: number
      userId: number
    }) => removeProjectMember(projectId, userId),
    onSuccess: (_, variables) => {
      void queryClient.invalidateQueries({
        queryKey: ['projects', variables.projectId],
      })
      void queryClient.invalidateQueries({ queryKey: ['projects'] })
    },
  })
}
