import axios from 'axios'
import { clearStoredSession, readStoredToken } from './tokenStorage'
import type { ApiErrorResponse } from '../types'

let unauthorizedHandler: (() => void) | undefined

export const apiClient = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api',
  headers: {
    'Content-Type': 'application/json',
  },
})

apiClient.interceptors.request.use((config) => {
  const token = readStoredToken()

  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }

  return config
})

apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    if (axios.isAxiosError<ApiErrorResponse>(error) && error.response?.status === 401) {
      clearStoredSession()
      unauthorizedHandler?.()
    }

    return Promise.reject(error)
  },
)

export function setUnauthorizedHandler(handler?: () => void) {
  unauthorizedHandler = handler
}

export function getApiErrorMessage(
  error: unknown,
  fallback = 'Something went wrong. Please try again.',
) {
  if (axios.isAxiosError<ApiErrorResponse>(error)) {
    return error.response?.data?.error.message ?? error.message ?? fallback
  }

  if (error instanceof Error) {
    return error.message
  }

  return fallback
}
