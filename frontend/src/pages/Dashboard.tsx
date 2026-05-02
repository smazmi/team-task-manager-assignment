import { useEffect, useMemo, useState, type FormEvent } from 'react'
import {
  ArrowRight,
  CircleAlert,
  ClipboardList,
  KanbanSquare,
  ListTodo,
  Plus,
  UsersRound,
} from 'lucide-react'
import { Link, useSearchParams } from 'react-router-dom'
import { getApiErrorMessage } from '../api/axios'
import { Modal } from '../components/Modal'
import { ProjectList } from '../components/ProjectList'
import { StatCard } from '../components/StatCard'
import { useAuth } from '../context/AuthContext'
import { useDashboard } from '../hooks/useDashboard'
import { useCreateProject, useProjects } from '../hooks/useProjects'

export function Dashboard() {
  const { user } = useAuth()
  const [searchParams, setSearchParams] = useSearchParams()
  const [isCreateProjectOpen, setCreateProjectOpen] = useState(false)
  const [projectName, setProjectName] = useState('')
  const [projectDescription, setProjectDescription] = useState('')
  const [formError, setFormError] = useState('')

  const projectsQuery = useProjects()
  const createProjectMutation = useCreateProject()

  const projects = projectsQuery.data ?? []
  const selectedProjectId = Number(searchParams.get('projectId') ?? '') || null

  useEffect(() => {
    if (!selectedProjectId && projects.length > 0) {
      setSearchParams({ projectId: String(projects[0].id) }, { replace: true })
    }
  }, [projects, selectedProjectId, setSearchParams])

  const selectedProject = useMemo(
    () => projects.find((project) => project.id === selectedProjectId) ?? null,
    [projects, selectedProjectId],
  )

  const selectedRole = useMemo(() => {
    if (!selectedProject || !user) {
      return null
    }

    return (
      selectedProject.members.find((member) => member.user_id === user.id)?.role ?? null
    )
  }, [selectedProject, user])

  const statsQuery = useDashboard(selectedProject?.id ?? null)
  const stats = statsQuery.data

  const maxWorkloadCount = useMemo(() => {
    if (!stats?.tasks_per_user.length) {
      return 1
    }

    return Math.max(...stats.tasks_per_user.map((member) => member.count), 1)
  }, [stats])

  const handleCreateProject = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    setFormError('')

    if (projectName.trim().length < 3) {
      setFormError('Project name must be at least 3 characters long.')
      return
    }

    try {
      const project = await createProjectMutation.mutateAsync({
        name: projectName.trim(),
        description: projectDescription.trim(),
      })

      setSearchParams({ projectId: String(project.id) })
      setProjectName('')
      setProjectDescription('')
      setCreateProjectOpen(false)
    } catch {
      return
    }
  }

  if (projectsQuery.isLoading) {
    return (
      <section className="state-panel">
        <h2 className="state-title">Loading projects</h2>
        <p className="state-copy">Gathering your workspace projects and statistics.</p>
      </section>
    )
  }

  if (projectsQuery.error) {
    return (
      <section className="state-panel">
        <h2 className="state-title">Unable to load the dashboard</h2>
        <p className="state-copy">{getApiErrorMessage(projectsQuery.error)}</p>
      </section>
    )
  }

  return (
    <div className="page">
      <header className="page-header">
        <div className="page-heading">
          <p className="eyebrow">Workspace overview</p>
          <h1 className="page-title">Dashboard</h1>
          <p className="page-subtitle">
            Review project-level workload, overdue items, and progress distribution.
          </p>
        </div>

        <button
          className="button button-primary"
          onClick={() => setCreateProjectOpen(true)}
          type="button"
        >
          <Plus size={16} />
          New project
        </button>
      </header>

      {projects.length === 0 ? (
        <section className="empty-state">
          <h2 className="state-title">No projects yet</h2>
          <p className="state-copy">
            Start by creating your first project so tasks, assignments, and dashboards
            have a home.
          </p>
          <div className="empty-state-actions">
            <button
              className="button button-primary"
              onClick={() => setCreateProjectOpen(true)}
              type="button"
            >
              <Plus size={16} />
              Create your first project
            </button>
          </div>
        </section>
      ) : (
        <div className="dashboard-layout">
          <section className="panel">
            <div className="panel-header">
              <div>
                <h2 className="panel-title">Projects</h2>
                <p className="panel-subtitle">Select a project to inspect its workload.</p>
              </div>
            </div>

            <div className="panel-body">
              <ProjectList
                currentUserId={user?.id ?? 0}
                projects={projects}
                selectedProjectId={selectedProjectId}
              />
            </div>
          </section>

          <div className="page" id="project-metrics">
            <section className="page-header">
              <div className="page-heading">
                <div className="page-title-row">
                  <h2 className="page-title">{selectedProject?.name}</h2>
                  {selectedRole ? (
                    <span
                      className={`badge ${
                        selectedRole === 'admin' ? 'badge-admin' : 'badge-member'
                      }`}
                    >
                      {selectedRole === 'admin' ? 'Admin access' : 'Member access'}
                    </span>
                  ) : null}
                </div>

                <p className="page-subtitle">
                  {selectedProject?.description || 'No project description yet.'}
                </p>
              </div>

              {selectedProject ? (
                <Link
                  className="button button-secondary"
                  to={`/projects/${selectedProject.id}`}
                >
                  Open board
                  <ArrowRight size={16} />
                </Link>
              ) : null}
            </section>

            {statsQuery.isLoading ? (
              <section className="state-panel">
                <h2 className="state-title">Loading project metrics</h2>
                <p className="state-copy">
                  Pulling the latest task counts and assignment breakdown.
                </p>
              </section>
            ) : statsQuery.error ? (
              <section className="state-panel">
                <h2 className="state-title">Unable to load project metrics</h2>
                <p className="state-copy">{getApiErrorMessage(statsQuery.error)}</p>
              </section>
            ) : stats ? (
              <>
                <section className="stats-grid">
                  <StatCard
                    icon={ClipboardList}
                    label="Total tasks"
                    tone="accent"
                    value={stats.total_tasks}
                  />
                  <StatCard
                    icon={ListTodo}
                    label="To do"
                    tone="neutral"
                    value={stats.tasks_by_status.todo}
                  />
                  <StatCard
                    icon={KanbanSquare}
                    label="In progress"
                    tone="warning"
                    value={stats.tasks_by_status.in_progress}
                  />
                  <StatCard
                    icon={UsersRound}
                    label="Done"
                    tone="success"
                    value={stats.tasks_by_status.done}
                  />
                  <StatCard
                    icon={CircleAlert}
                    label="Overdue"
                    tone="danger"
                    value={stats.overdue_tasks}
                  />
                </section>

                <section className="dashboard-analytics">
                  <div className="panel">
                    <div className="panel-header">
                      <div>
                        <h2 className="panel-title">Status breakdown</h2>
                        <p className="panel-subtitle">
                          Progress distribution for the selected project.
                        </p>
                      </div>
                    </div>

                    <div className="panel-body">
                      <div className="bars">
                        {[
                          {
                            key: 'todo',
                            label: 'To Do',
                            value: stats.tasks_by_status.todo,
                            className: 'todo',
                          },
                          {
                            key: 'in_progress',
                            label: 'In Progress',
                            value: stats.tasks_by_status.in_progress,
                            className: 'in-progress',
                          },
                          {
                            key: 'done',
                            label: 'Done',
                            value: stats.tasks_by_status.done,
                            className: 'done',
                          },
                        ].map((entry) => {
                          const percentage =
                            stats.total_tasks === 0
                              ? 0
                              : Math.round((entry.value / stats.total_tasks) * 100)

                          return (
                            <div className="bar-row" key={entry.key}>
                              <div className="bar-topline">
                                <span>{entry.label}</span>
                                <span>
                                  {entry.value} tasks ({percentage}%)
                                </span>
                              </div>

                              <div className="bar-track">
                                <div
                                  className={`bar-fill ${entry.className}`}
                                  style={{ width: `${percentage}%` }}
                                />
                              </div>
                            </div>
                          )
                        })}
                      </div>
                    </div>
                  </div>

                  <div className="panel">
                    <div className="panel-header">
                      <div>
                        <h2 className="panel-title">Tasks per user</h2>
                        <p className="panel-subtitle">
                          Assignment load across project members.
                        </p>
                      </div>
                    </div>

                    <div className="panel-body">
                      {stats.tasks_per_user.length === 0 ? (
                        <p className="state-copy">
                          No assignments yet. Tasks become visible here once they are
                          assigned.
                        </p>
                      ) : (
                        <div className="bars">
                          {stats.tasks_per_user.map((member) => (
                            <div className="bar-row" key={member.user_id}>
                              <div className="bar-topline">
                                <span>{member.user_name}</span>
                                <span>{member.count} tasks</span>
                              </div>

                              <div className="bar-track">
                                <div
                                  className="bar-fill workload"
                                  style={{
                                    width: `${Math.max(
                                      10,
                                      Math.round((member.count / maxWorkloadCount) * 100),
                                    )}%`,
                                  }}
                                />
                              </div>
                            </div>
                          ))}
                        </div>
                      )}
                    </div>
                  </div>
                </section>
              </>
            ) : null}
          </div>
        </div>
      )}

      <Modal
        description="Create a new project and become its initial admin."
        onClose={() => {
          setCreateProjectOpen(false)
          setFormError('')
        }}
        open={isCreateProjectOpen}
        title="Create project"
      >
        <form className="inline-form" onSubmit={handleCreateProject}>
          {formError ? <div className="error-banner">{formError}</div> : null}
          {createProjectMutation.error ? (
            <div className="error-banner">
              {getApiErrorMessage(createProjectMutation.error)}
            </div>
          ) : null}

          <label className="field">
            <span>Project name</span>
            <input
              onChange={(event) => setProjectName(event.target.value)}
              placeholder="Q3 Operations Refresh"
              type="text"
              value={projectName}
            />
          </label>

          <label className="field">
            <span>Description</span>
            <textarea
              onChange={(event) => setProjectDescription(event.target.value)}
              placeholder="Scope, outcomes, and context for the project."
              value={projectDescription}
            />
          </label>

          <div className="form-actions">
            <button
              className="button button-secondary"
              onClick={() => setCreateProjectOpen(false)}
              type="button"
            >
              Cancel
            </button>
            <button
              className="button button-primary"
              disabled={createProjectMutation.isPending}
              type="submit"
            >
              {createProjectMutation.isPending ? 'Creating...' : 'Create project'}
            </button>
          </div>
        </form>
      </Modal>
    </div>
  )
}
