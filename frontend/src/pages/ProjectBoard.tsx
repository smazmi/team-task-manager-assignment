import { useMemo, useState, type FormEvent } from 'react'
import { ArrowLeft, Plus, UsersRound } from 'lucide-react'
import { Link, useParams } from 'react-router-dom'
import { getApiErrorMessage } from '../api/axios'
import { Modal } from '../components/Modal'
import { TaskCard } from '../components/TaskCard'
import { useAuth } from '../context/AuthContext'
import { useProject, useAddProjectMember, useRemoveProjectMember } from '../hooks/useProjects'
import { useCreateTask, useDeleteTask, useTasks, useUpdateTaskStatus } from '../hooks/useTasks'
import {
  PROJECT_ROLE_OPTIONS,
  TASK_PRIORITY_OPTIONS,
  type ProjectRole,
  type TaskPriority,
  type TaskStatus,
} from '../types'

export function ProjectBoard() {
  const { projectId: rawProjectId } = useParams()
  const projectId = Number(rawProjectId ?? '')
  const { user } = useAuth()

  const projectQuery = useProject(projectId)
  const tasksQuery = useTasks(projectId)
  const createTaskMutation = useCreateTask()
  const updateTaskStatusMutation = useUpdateTaskStatus()
  const deleteTaskMutation = useDeleteTask()
  const addMemberMutation = useAddProjectMember()
  const removeMemberMutation = useRemoveProjectMember()

  const [isCreateTaskOpen, setCreateTaskOpen] = useState(false)
  const [taskTitle, setTaskTitle] = useState('')
  const [taskDescription, setTaskDescription] = useState('')
  const [taskDueDate, setTaskDueDate] = useState('')
  const [taskPriority, setTaskPriority] = useState<TaskPriority>('medium')
  const [taskAssigneeId, setTaskAssigneeId] = useState('')
  const [memberEmail, setMemberEmail] = useState('')
  const [memberRole, setMemberRole] = useState<ProjectRole>('member')
  const [taskFormError, setTaskFormError] = useState('')
  const [memberFormError, setMemberFormError] = useState('')
  const [statusUpdatingTaskId, setStatusUpdatingTaskId] = useState<number | null>(null)
  const [deletingTaskId, setDeletingTaskId] = useState<number | null>(null)
  const [pendingDeleteTaskId, setPendingDeleteTaskId] = useState<number | null>(null)
  const [removingMemberId, setRemovingMemberId] = useState<number | null>(null)

  const project = projectQuery.data
  const tasks = tasksQuery.data ?? []

  const membership = useMemo(() => {
    if (!project || !user) {
      return null
    }

    return project.members.find((member) => member.user_id === user.id) ?? null
  }, [project, user])

  const isAdmin = membership?.role === 'admin'

  const groupedTasks = useMemo(
    () => ({
      todo: tasks.filter((task) => task.status === 'todo'),
      in_progress: tasks.filter((task) => task.status === 'in_progress'),
      done: tasks.filter((task) => task.status === 'done'),
    }),
    [tasks],
  )

  const handleStatusChange = async (taskId: number, status: TaskStatus) => {
    try {
      setStatusUpdatingTaskId(taskId)
      await updateTaskStatusMutation.mutateAsync({
        taskId,
        payload: { status },
      })
    } catch {
      return
    } finally {
      setStatusUpdatingTaskId(null)
    }
  }

  const handleDeleteRequest = (taskId: number) => {
    setPendingDeleteTaskId(taskId)
  }

  const handleConfirmDelete = async () => {
    if (!pendingDeleteTaskId) {
      return
    }

    try {
      setDeletingTaskId(pendingDeleteTaskId)
      await deleteTaskMutation.mutateAsync({
        taskId: pendingDeleteTaskId,
        projectId,
      })
      setPendingDeleteTaskId(null)
    } catch {
      return
    } finally {
      setDeletingTaskId(null)
    }
  }

  const pendingDeleteTask = pendingDeleteTaskId
    ? tasks.find((task) => task.id === pendingDeleteTaskId) ?? null
    : null

  const handleDueDateChange = (value: string) => {
    if (!value) {
      setTaskDueDate('')
      return
    }

    const [datePart] = value.split('T')
    const year = datePart?.split('-')[0] ?? ''
    if (year.length > 4) {
      return
    }

    setTaskDueDate(value)
  }

  const handleCreateTask = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    setTaskFormError('')

    if (taskTitle.trim().length === 0) {
      setTaskFormError('Task title is required.')
      return
    }

    try {
      await createTaskMutation.mutateAsync({
        project_id: projectId,
        title: taskTitle.trim(),
        description: taskDescription.trim(),
        due_date: taskDueDate ? new Date(taskDueDate).toISOString() : undefined,
        priority: taskPriority,
        assignee_id: taskAssigneeId ? Number(taskAssigneeId) : undefined,
      })

      setTaskTitle('')
      setTaskDescription('')
      setTaskDueDate('')
      setTaskPriority('medium')
      setTaskAssigneeId('')
      setCreateTaskOpen(false)
    } catch {
      return
    }
  }

  const handleAddMember = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    setMemberFormError('')

    if (memberEmail.trim().length === 0) {
      setMemberFormError('A valid email address is required.')
      return
    }

    try {
      await addMemberMutation.mutateAsync({
        projectId,
        payload: {
          email: memberEmail.trim(),
          role: memberRole,
        },
      })

      setMemberEmail('')
      setMemberRole('member')
    } catch {
      return
    }
  }

  const handleRemoveMember = async (userId: number) => {
    try {
      setRemovingMemberId(userId)
      await removeMemberMutation.mutateAsync({ projectId, userId })
    } catch {
      return
    } finally {
      setRemovingMemberId(null)
    }
  }

  if (!Number.isInteger(projectId) || projectId < 1) {
    return (
      <section className="state-panel">
        <h2 className="state-title">Invalid project</h2>
        <p className="state-copy">The project route is missing a valid project identifier.</p>
      </section>
    )
  }

  if (projectQuery.isLoading || tasksQuery.isLoading) {
    return (
      <section className="state-panel">
        <h2 className="state-title">Loading project board</h2>
        <p className="state-copy">Bringing in project details, members, and task columns.</p>
      </section>
    )
  }

  if (projectQuery.error) {
    return (
      <section className="state-panel">
        <h2 className="state-title">Unable to load the project</h2>
        <p className="state-copy">{getApiErrorMessage(projectQuery.error)}</p>
      </section>
    )
  }

  if (tasksQuery.error) {
    return (
      <section className="state-panel">
        <h2 className="state-title">Unable to load tasks</h2>
        <p className="state-copy">{getApiErrorMessage(tasksQuery.error)}</p>
      </section>
    )
  }

  if (!project) {
    return null
  }

  return (
    <div className="page">
      <header className="page-header">
        <div className="page-heading">
          <div className="toolbar">
            <Link className="button button-ghost button-sm" to="/dashboard">
              <ArrowLeft size={16} />
              Back to dashboard
            </Link>
          </div>

          <div className="page-title-row">
            <h1 className="page-title">{project.name}</h1>
            {membership ? (
              <span className={`badge ${isAdmin ? 'badge-admin' : 'badge-member'}`}>
                {isAdmin ? 'Admin' : 'Member'}
              </span>
            ) : null}
          </div>

          <p className="page-subtitle">
            {project.description || 'No project description yet.'}
          </p>
        </div>

        {isAdmin ? (
          <button
            className="button button-primary"
            onClick={() => setCreateTaskOpen(true)}
            type="button"
          >
            <Plus size={16} />
            New task
          </button>
        ) : null}
      </header>

      <div className="board-layout">
        <div className="task-columns">
          {[
            { key: 'todo', title: 'To Do', tasks: groupedTasks.todo },
            { key: 'in_progress', title: 'In Progress', tasks: groupedTasks.in_progress },
            { key: 'done', title: 'Done', tasks: groupedTasks.done },
          ].map((column) => (
            <section className="task-column" key={column.key}>
              <div className="task-column-header">
                <div>
                  <h2 className="panel-title">{column.title}</h2>
                  <p className="panel-subtitle">{column.tasks.length} tasks</p>
                </div>
              </div>

              <div className="task-column-body">
                {column.tasks.length === 0 ? (
                  <p className="state-copy">No tasks in this column yet.</p>
                ) : (
                  column.tasks.map((task) => (
                    <TaskCard
                      canDelete={isAdmin}
                      isDeleting={deletingTaskId === task.id}
                      isUpdatingStatus={statusUpdatingTaskId === task.id}
                      key={task.id}
                      onDelete={handleDeleteRequest}
                      onStatusChange={handleStatusChange}
                      task={task}
                    />
                  ))
                )}
              </div>
            </section>
          ))}
        </div>

        <aside className="page">
          <section className="panel">
            <div className="panel-header">
              <div>
                <h2 className="panel-title">Project members</h2>
                <p className="panel-subtitle">
                  Assignment and access roles for this project.
                </p>
              </div>
            </div>

            <div className="panel-body">
              <div className="member-list">
                {project.members.map((member) => (
                  <article className="member-row" key={member.id}>
                    <div className="member-row-main">
                      <p className="member-row-name">{member.user.name}</p>
                      <div className="member-row-meta">
                        <span>{member.user.email}</span>
                      </div>
                    </div>

                    <div className="toolbar">
                      <span
                        className={`badge ${
                          member.role === 'admin' ? 'badge-admin' : 'badge-member'
                        }`}
                      >
                        {member.role === 'admin' ? 'Admin' : 'Member'}
                      </span>

                      {isAdmin ? (
                        <button
                          className="button button-danger button-sm"
                          disabled={removingMemberId === member.user_id}
                          onClick={() => handleRemoveMember(member.user_id)}
                          type="button"
                        >
                          {removingMemberId === member.user_id ? 'Removing...' : 'Remove'}
                        </button>
                      ) : null}
                    </div>
                  </article>
                ))}
              </div>
            </div>
          </section>

          {isAdmin ? (
            <section className="panel">
              <div className="panel-header">
                <div>
                  <h2 className="panel-title">Manage members</h2>
                  <p className="panel-subtitle">
                    Add a registered user to this project by their email address.
                  </p>
                </div>
              </div>

              <div className="panel-body">
                <form className="inline-form" onSubmit={handleAddMember}>
                  {memberFormError ? <div className="error-banner">{memberFormError}</div> : null}
                  {addMemberMutation.error ? (
                    <div className="error-banner">
                      {getApiErrorMessage(addMemberMutation.error)}
                    </div>
                  ) : null}
                  {removeMemberMutation.error ? (
                    <div className="error-banner">
                      {getApiErrorMessage(removeMemberMutation.error)}
                    </div>
                  ) : null}

                  <label className="field">
                    <span>Email Address</span>
                    <input
                      onChange={(event) => setMemberEmail(event.target.value)}
                      placeholder="user@example.com"
                      type="email"
                      value={memberEmail}
                    />
                  </label>

                  <label className="field">
                    <span>Role</span>
                    <select
                      onChange={(event) => setMemberRole(event.target.value as ProjectRole)}
                      value={memberRole}
                    >
                      {PROJECT_ROLE_OPTIONS.map((roleOption) => (
                        <option key={roleOption.value} value={roleOption.value}>
                          {roleOption.label}
                        </option>
                      ))}
                    </select>
                  </label>

                  <button
                    className="button button-primary"
                    disabled={addMemberMutation.isPending}
                    type="submit"
                  >
                    <UsersRound size={16} />
                    {addMemberMutation.isPending ? 'Adding...' : 'Add member'}
                  </button>
                </form>
              </div>
            </section>
          ) : null}
        </aside>
      </div>

      <Modal
        description="Create a task, set a due date, and optionally assign it to a project member."
        onClose={() => {
          setCreateTaskOpen(false)
          setTaskFormError('')
        }}
        open={isCreateTaskOpen}
        title="Create task"
      >
        <form className="inline-form" onSubmit={handleCreateTask}>
          {taskFormError ? <div className="error-banner">{taskFormError}</div> : null}
          {createTaskMutation.error ? (
            <div className="error-banner">
              {getApiErrorMessage(createTaskMutation.error)}
            </div>
          ) : null}

          <div className="field-grid columns-2">
            <label className="field">
              <span>Title</span>
              <input
                onChange={(event) => setTaskTitle(event.target.value)}
                placeholder="Prepare launch checklist"
                type="text"
                value={taskTitle}
              />
            </label>

            <label className="field">
              <span>Priority</span>
              <select
                onChange={(event) => setTaskPriority(event.target.value as TaskPriority)}
                value={taskPriority}
              >
                {TASK_PRIORITY_OPTIONS.map((priorityOption) => (
                  <option key={priorityOption.value} value={priorityOption.value}>
                    {priorityOption.label}
                  </option>
                ))}
              </select>
            </label>
          </div>

          <label className="field">
            <span>Description</span>
            <textarea
              onChange={(event) => setTaskDescription(event.target.value)}
              placeholder="Describe the expected deliverable and any context."
              value={taskDescription}
            />
          </label>

          <div className="field-grid columns-2">
            <label className="field">
              <span>Due date</span>
              <input
                max="9999-12-31T23:59"
                onChange={(event) => handleDueDateChange(event.target.value)}
                type="datetime-local"
                value={taskDueDate}
              />
            </label>

            <label className="field">
              <span>Assignee</span>
              <select
                onChange={(event) => setTaskAssigneeId(event.target.value)}
                value={taskAssigneeId}
              >
                <option value="">Unassigned</option>
                {project.members.map((member) => (
                  <option key={member.user_id} value={member.user_id}>
                    {member.user.name}
                  </option>
                ))}
              </select>
            </label>
          </div>

          <div className="form-actions">
            <button
              className="button button-secondary"
              onClick={() => setCreateTaskOpen(false)}
              type="button"
            >
              Cancel
            </button>
            <button
              className="button button-primary"
              disabled={createTaskMutation.isPending}
              type="submit"
            >
              {createTaskMutation.isPending ? 'Creating...' : 'Create task'}
            </button>
          </div>
        </form>
      </Modal>

      <Modal
        description={
          pendingDeleteTask
            ? `Delete "${pendingDeleteTask.title}"? This cannot be undone.`
            : 'Delete this task? This cannot be undone.'
        }
        onClose={() => setPendingDeleteTaskId(null)}
        open={pendingDeleteTaskId !== null}
        title="Confirm delete"
      >
        <div className="form-actions">
          <button
            className="button button-secondary"
            onClick={() => setPendingDeleteTaskId(null)}
            type="button"
          >
            Cancel
          </button>
          <button
            className="button button-danger"
            disabled={deletingTaskId === pendingDeleteTaskId}
            onClick={handleConfirmDelete}
            type="button"
          >
            {deletingTaskId === pendingDeleteTaskId ? 'Deleting...' : 'Delete task'}
          </button>
        </div>
      </Modal>
    </div>
  )
}
