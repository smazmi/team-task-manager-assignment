import {
  CalendarClock,
  CircleUserRound,
  Flag,
  UserRoundPen,
} from 'lucide-react'
import { TASK_STATUS_OPTIONS, type Task, type TaskStatus } from '../types'

interface TaskCardProps {
  isUpdatingStatus: boolean
  canDelete?: boolean
  isDeleting?: boolean
  onDelete?: (taskId: number) => void
  onStatusChange: (taskId: number, status: TaskStatus) => void
  task: Task
}

function formatDueDate(value?: string | null) {
  if (!value) {
    return 'No due date'
  }

  return new Intl.DateTimeFormat(undefined, {
    dateStyle: 'medium',
    timeStyle: 'short',
  }).format(new Date(value))
}

function priorityClassName(priority: Task['priority']) {
  switch (priority) {
    case 'high':
      return 'badge badge-priority-high'
    case 'medium':
      return 'badge badge-priority-medium'
    default:
      return 'badge badge-priority-low'
  }
}

function priorityLabel(priority: Task['priority']) {
  switch (priority) {
    case 'high':
      return 'High Priority'
    case 'medium':
      return 'Medium Priority'
    default:
      return 'Low Priority'
  }
}

export function TaskCard({
  isUpdatingStatus,
  canDelete = false,
  isDeleting = false,
  onDelete,
  onStatusChange,
  task,
}: TaskCardProps) {
  return (
    <article className="task-card">
      <div className="task-card-header">
        <div>
          <h3 className="task-card-title">{task.title}</h3>
          <p className="task-card-description">
            {task.description || 'No task description provided.'}
          </p>
        </div>

        <span className={priorityClassName(task.priority)}>
          <Flag size={14} />
          {priorityLabel(task.priority)}
        </span>
      </div>

      <div className="task-meta">
        <span className="task-meta-item">
          <CircleUserRound size={15} />
          {task.assignee?.name ?? 'Unassigned'}
        </span>

        <span className="task-meta-item">
          <CalendarClock size={15} />
          {formatDueDate(task.due_date)}
        </span>

        <span className="task-meta-item">
          <UserRoundPen size={15} />
          Created by {task.creator.name}
        </span>
      </div>

      <div className="task-card-footer">
        {canDelete ? (
          <button
            className="button button-danger button-sm"
            disabled={isDeleting}
            onClick={() => onDelete?.(task.id)}
            type="button"
          >
            {isDeleting ? 'Deleting...' : 'Delete task'}
          </button>
        ) : null}

        <label className="field">
          <span>Status</span>
          <select
            className="status-select"
            disabled={isUpdatingStatus}
            onChange={(event) =>
              onStatusChange(task.id, event.target.value as TaskStatus)
            }
            value={task.status}
          >
            {TASK_STATUS_OPTIONS.map((option) => (
              <option key={option.value} value={option.value}>
                {option.label}
              </option>
            ))}
          </select>
        </label>
      </div>
    </article>
  )
}
