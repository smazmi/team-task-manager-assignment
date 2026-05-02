import { useMutation } from '@tanstack/react-query'
import { useNavigate } from 'react-router-dom'
import { loginRequest, registerRequest } from '../api/taskManagerApi'
import { useAuth } from '../context/AuthContext'
import type { LoginInput, RegisterInput } from '../types'

export function useLogin() {
  const navigate = useNavigate()
  const { login } = useAuth()

  return useMutation({
    mutationFn: (payload: LoginInput) => loginRequest(payload),
    onSuccess: (session) => {
      login(session)
      navigate('/dashboard', { replace: true })
    },
  })
}

export function useRegister() {
  const navigate = useNavigate()
  const { login } = useAuth()

  return useMutation({
    mutationFn: (payload: RegisterInput) => registerRequest(payload),
    onSuccess: (session) => {
      login(session)
      navigate('/dashboard', { replace: true })
    },
  })
}
