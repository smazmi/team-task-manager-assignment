import type { LucideIcon } from 'lucide-react'

interface StatCardProps {
  icon: LucideIcon
  label: string
  tone: 'accent' | 'danger' | 'neutral' | 'success' | 'warning'
  value: number | string
}

export function StatCard({ icon: Icon, label, tone, value }: StatCardProps) {
  return (
    <article className="panel stat-card">
      <div className="stat-card-header">
        <p className="stat-card-label">{label}</p>
        <span className={`stat-card-icon tone-${tone}`}>
          <Icon size={18} />
        </span>
      </div>

      <p className="stat-card-value">{value}</p>
    </article>
  )
}
