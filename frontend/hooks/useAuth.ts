'use client'

import { useState, useEffect } from 'react'
import { User, AuthError, Session } from '@supabase/supabase-js'
import { supabase } from '@/lib/supabase'

export interface AuthState {
  user: User | null
  session: Session | null
  loading: boolean
  error: AuthError | null
}

export function useAuth() {
  const [state, setState] = useState<AuthState>({
    user: null,
    session: null,
    loading: true,
    error: null,
  })

  useEffect(() => {
    // Get initial session
    const getSession = async () => {
      const { data: { session }, error } = await supabase.auth.getSession()
      setState({
        user: session?.user ?? null,
        session,
        loading: false,
        error,
      })
    }

    getSession()

    // Listen for auth changes
    const { data: { subscription } } = supabase.auth.onAuthStateChange(
      async (event, session) => {
        setState({
          user: session?.user ?? null,
          session,
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
      error,
      user: data.user,
      session: data.session,
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
      error,
    }))

    return { data, error }
  }

  const signOut = async () => {
    setState(prev => ({ ...prev, loading: true }))
    
    const { error } = await supabase.auth.signOut()
    
    setState({
      user: null,
      session: null,
      loading: false,
      error,
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