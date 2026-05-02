import { ArrowRight, Users } from 'lucide-react'
import { Link } from 'react-router-dom'
import type { Project } from '../types'

interface ProjectListProps {
  currentUserId: number
  projects: Project[]
  selectedProjectId: number | null
}

export function ProjectList({
  currentUserId,
  projects,
  selectedProjectId,
}: ProjectListProps) {
  return (
    <div className="project-list">
      {projects.map((project) => {
        const membership = project.members.find((member) => member.user_id === currentUserId)
        const isSelected = selectedProjectId === project.id

        return (
          <article
            className={`project-list-item${isSelected ? ' is-selected' : ''}`}
            key={project.id}
          >
            <header>
              <div>
                <h3>{project.name}</h3>
                <p>{project.description || 'No project description yet.'}</p>
              </div>

              <span
                className={`badge ${
                  membership?.role === 'admin' ? 'badge-admin' : 'badge-member'
                }`}
              >
                {membership?.role === 'admin' ? 'Admin' : 'Member'}
              </span>
            </header>

            <div className="project-list-meta">
              <span className="task-meta-item">
                <Users size={15} />
                {project.members.length} members
              </span>
            </div>

            <div className="project-list-actions">
              <Link
                className="button button-secondary button-sm"
                to={{
                  pathname: '/dashboard',
                  search: `?projectId=${project.id}`,
                  hash: '#project-metrics',
                }}
              >
                {isSelected ? 'Metrics selected' : 'View metrics'}
              </Link>

              <Link className="button button-ghost button-sm" to={`/projects/${project.id}`}>
                Open board
                <ArrowRight size={15} />
              </Link>
            </div>
          </article>
        )
      })}
    </div>
  )
}
