'use client'

import { useState, useEffect, useCallback } from 'react'

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

interface User {
  id: number
  email: string
  name?: string
}

interface AuthError {
  message: string
}

interface AuthResponse {
  data: { user: User | null; token?: string } | null
  error: AuthError | null
}

export function useAuth() {
  const [user, setUser] = useState<User | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<AuthError | null>(null)

  useEffect(() => {
    const token = localStorage.getItem('auth_token')
    const storedUser = localStorage.getItem('auth_user')
    if (token && storedUser) {
      setUser(JSON.parse(storedUser))
    }
    setLoading(false)
  }, [])

  const signIn = useCallback(async (email: string, password: string): Promise<AuthResponse> => {
    setLoading(true)
    setError(null)
    try {
      const response = await fetch(`${API_BASE_URL}/api/v1/auth/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password }),
      })
      const data = await response.json()
      if (!response.ok) {
        const err = { message: data.error || 'Login failed' }
        setError(err)
        return { data: null, error: err }
      }
      localStorage.setItem('auth_token', data.token)
      localStorage.setItem('auth_user', JSON.stringify(data.user))
      setUser(data.user)
      return { data: { user: data.user, token: data.token }, error: null }
    } catch (e) {
      const err = { message: 'Network error' }
      setError(err)
      return { data: null, error: err }
    } finally {
      setLoading(false)
    }
  }, [])

  const signUp = useCallback(async (email: string, password: string): Promise<AuthResponse> => {
    setLoading(true)
    setError(null)
    try {
      const response = await fetch(`${API_BASE_URL}/api/v1/auth/register`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password }),
      })
      const data = await response.json()
      if (!response.ok) {
        const err = { message: data.error || 'Registration failed' }
        setError(err)
        return { data: null, error: err }
      }
      return { data: { user: data.user }, error: null }
    } catch (e) {
      const err = { message: 'Network error' }
      setError(err)
      return { data: null, error: err }
    } finally {
      setLoading(false)
    }
  }, [])

  const signOut = useCallback(async (): Promise<{ error: AuthError | null }> => {
    setLoading(true)
    try {
      const token = localStorage.getItem('auth_token')
      await fetch(`${API_BASE_URL}/api/v1/auth/logout`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
      })
      localStorage.removeItem('auth_token')
      localStorage.removeItem('auth_user')
      setUser(null)
      return { error: null }
    } catch (e) {
      return { error: { message: 'Logout failed' } }
    } finally {
      setLoading(false)
    }
  }, [])

  const resetPassword = useCallback(async (email: string): Promise<{ error: AuthError | null }> => {
    setLoading(true)
    setError(null)
    try {
      const response = await fetch(`${API_BASE_URL}/api/v1/auth/reset-password`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email }),
      })
      const data = await response.json()
      if (!response.ok) {
        const err = { message: data.error || 'Password reset request failed' }
        setError(err)
        return { error: err }
      }
      return { error: null }
    } catch (e) {
      const err = { message: 'Network error' }
      setError(err)
      return { error: err }
    } finally {
      setLoading(false)
    }
  }, [])

  const confirmPasswordReset = useCallback(async (token: string, newPassword: string): Promise<{ error: AuthError | null }> => {
    setLoading(true)
    setError(null)
    try {
      const response = await fetch(`${API_BASE_URL}/api/v1/auth/reset-password/confirm`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ token, password: newPassword }),
      })
      const data = await response.json()
      if (!response.ok) {
        const err = { message: data.error || 'Password reset failed' }
        setError(err)
        return { error: err }
      }
      return { error: null }
    } catch (e) {
      const err = { message: 'Network error' }
      setError(err)
      return { error: err }
    } finally {
      setLoading(false)
    }
  }, [])

  return {
    user,
    loading,
    error,
    signIn,
    signUp,
    signOut,
    resetPassword,
    confirmPasswordReset,
  }
}
