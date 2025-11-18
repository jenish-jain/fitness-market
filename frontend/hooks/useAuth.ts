'use client'

import { useEffect, useState } from 'react'
import { User, AuthError } from '@supabase/supabase-js'
import { supabase } from '@/lib/supabase'

interface AuthState {
  user: User | null
  loading: boolean
  error: AuthError | null
}

export function useAuth() {
  const [state, setState] = useState<AuthState>({
    user: null,
    loading: true,
    error: null,
  })

  useEffect(() => {
    // Get initial session
    const getInitialSession = async () => {
      try {
        const { data: { session }, error } = await supabase.auth.getSession()
        setState({
          user: session?.user ?? null,
          loading: false,
          error: error,
        })
      } catch (error) {
        setState({
          user: null,
          loading: false,
          error: error as AuthError,
        })
      }
    }

    getInitialSession()

    // Listen for auth changes
    const { data: { subscription } } = supabase.auth.onAuthStateChange(
      (event, session) => {
        setState({
          user: session?.user ?? null,
          loading: false,
          error: null,
        })
      }
    )

    return () => subscription.unsubscribe()
  }, [])

  const signIn = async (email: string, password: string) => {
    setState(prev => ({ ...prev, loading: true, error: null }))
    
    const { data, error } = await supabase.auth.signInWithPassword({
      email,
      password,
    })

    setState(prev => ({
      ...prev,
      loading: false,
      error: error,
    }))

    return { data, error }
  }

  const signUp = async (email: string, password: string) => {
    setState(prev => ({ ...prev, loading: true, error: null }))
    
    const { data, error } = await supabase.auth.signUp({
      email,
      password,
    })

    setState(prev => ({
      ...prev,
      loading: false,
      error: error,
    }))

    return { data, error }
  }

  const signOut = async () => {
    setState(prev => ({ ...prev, loading: true }))
    
    const { error } = await supabase.auth.signOut()
    
    setState({
      user: null,
      loading: false,
      error: error,
    })

    return { error }
  }

  return {
    ...state,
    signIn,
    signUp,
    signOut,
  }
}