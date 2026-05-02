import { useState, type FormEvent } from 'react'
import { Link } from 'react-router-dom'
import { getApiErrorMessage } from '../api/axios'
import { useRegister } from '../hooks/useAuthMutations'

export function Register() {
  const registerMutation = useRegister()
  const [name, setName] = useState('')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [validationError, setValidationError] = useState('')

  const errorMessage =
    validationError || getApiErrorMessage(registerMutation.error, '')

  const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    setValidationError('')

    if (name.trim().length < 2) {
      setValidationError('Name must be at least 2 characters long.')
      return
    }

    if (!email.trim() || !password.trim()) {
      setValidationError('Email and password are required.')
      return
    }

    if (password.length < 8) {
      setValidationError('Password must be at least 8 characters long.')
      return
    }

    try {
      await registerMutation.mutateAsync({
        name: name.trim(),
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
          <p className="eyebrow">Register</p>
          <h1 className="auth-title">Create your workspace access</h1>
          <p className="auth-description">
            Join projects, track assignments, and collaborate with your team.
          </p>
        </div>

        {errorMessage ? <div className="error-banner">{errorMessage}</div> : null}

        <form className="field-grid" onSubmit={handleSubmit}>
          <label className="field">
            <span>Full name</span>
            <input
              autoComplete="name"
              onChange={(event) => setName(event.target.value)}
              placeholder="Jane Doe"
              type="text"
              value={name}
            />
          </label>

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
              autoComplete="new-password"
              onChange={(event) => setPassword(event.target.value)}
              placeholder="Create a strong password"
              type="password"
              value={password}
            />
          </label>

          <button
            className="button button-primary"
            disabled={registerMutation.isPending}
            type="submit"
          >
            {registerMutation.isPending ? 'Creating account...' : 'Create account'}
          </button>
        </form>

        <div className="auth-footer">
          <span>Already registered?</span>
          <Link className="auth-link" to="/login">
            Sign in instead
          </Link>
        </div>
      </section>
    </div>
  )
}
