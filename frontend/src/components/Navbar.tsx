import { FolderKanban, LayoutDashboard, LogOut } from 'lucide-react'
import { NavLink } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'

export function Navbar() {
  const { logout, user } = useAuth()

  return (
    <header className="navbar">
      <div className="navbar-inner">
        <div className="brand">
          <span className="brand-mark">
            <FolderKanban size={20} />
          </span>
          <div className="brand-copy">
            <span className="brand-title">Team Task Manager</span>
            <span className="brand-subtitle">Operational planning workspace</span>
          </div>
        </div>

        <div className="nav-actions">
          <nav className="nav-links" aria-label="Primary">
            <NavLink className="nav-link" to="/dashboard">
              <LayoutDashboard size={16} />
              Dashboard
            </NavLink>
          </nav>

          <div className="nav-user">
            <span className="nav-user-name">{user?.name}</span>
            <span className="nav-user-email">{user?.email}</span>
          </div>

          <button className="button button-secondary" onClick={logout} type="button">
            <LogOut size={16} />
            Logout
          </button>
        </div>
      </div>
    </header>
  )
}
