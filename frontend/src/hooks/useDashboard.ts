import { useQuery } from '@tanstack/react-query'
import { fetchDashboardStats } from '../api/taskManagerApi'

export function useDashboard(projectId: number | null) {
  return useQuery({
    queryKey: ['dashboard', projectId],
    queryFn: () => fetchDashboardStats(projectId as number),
    enabled: Boolean(projectId),
  })
}
