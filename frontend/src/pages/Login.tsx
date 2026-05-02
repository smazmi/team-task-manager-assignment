import { useState, type FormEvent } from 'react'
import { Link } from 'react-router-dom'
import { getApiErrorMessage } from '../api/axios'
import { useLogin } from '../hooks/useAuthMutations'

export function Login() {
  const loginMutation = useLogin()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [validationError, setValidationError] = useState('')

  const errorMessage = validationError || getApiErrorMessage(loginMutation.error, '')

  const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    setValidationError('')

    if (!email.trim() || !password.trim()) {
      setValidationError('Email and password are required.')
      return
    }

    try {
      await loginMutation.mutateAsync({
        email: email.trim(),
        password,
      })
    } catch {
      return
    }
  }

  return (
    <div className="auth-shell">
      <section className="auth-card">
        <div className="auth-copy">
          <p className="eyebrow">Sign in</p>
          <h1 className="auth-title">Welcome back</h1>
          <p className="auth-description">
            Access your project dashboard, board views, and workload analytics.
          </p>
        </div>

        {errorMessage ? <div className="error-banner">{errorMessage}</div> : null}

        <form className="field-grid" onSubmit={handleSubmit}>
          <label className="field">
            <span>Email address</span>
            <input
              autoComplete="email"
              onChange={(event) => setEmail(event.target.value)}
              placeholder="jane@company.com"
              type="email"
              value={email}
            />
          </label>

          <label className="field">
            <span>Password</span>
            <input
              autoComplete="current-password"
              onChange={(event) => setPassword(event.target.value)}
              placeholder="Enter your password"
              type="password"
              value={password}
            />
          </label>

          <button className="button button-primary" disabled={loginMutation.isPending} type="submit">
            {loginMutation.isPending ? 'Signing in...' : 'Sign in'}
          </button>
        </form>

        <div className="auth-footer">
          <span>New to the workspace?</span>
          <Link className="auth-link" to="/register">
            Create an account
          </Link>
        </div>
      </section>
    </div>
  )
}
