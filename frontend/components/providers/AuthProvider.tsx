'use client'

import { createContext, useContext, useEffect, useState } from 'react'
import { supabase, type User } from '@/lib/supabase'
import type { AuthError, Session } from '@supabase/supabase-js'

interface AuthContextType {
  user: User | null
  session: Session | null
  loading: boolean
  error: AuthError | null
  signIn: (email: string, password: string) => Promise<{ data?: any; error?: AuthError | null }>
  signUp: (email: string, password: string) => Promise<{ data?: any; error?: AuthError | null }>
  signOut: () => Promise<{ error?: AuthError | null }>
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [session, setSession] = useState<Session | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<AuthError | null>(null)

  useEffect(() => {
    // Get initial session
    supabase.auth.getSession().then(({ data: { session } }) => {
      setSession(session)
      setUser(session?.user ?? null)
      setLoading(false)
    })

    // Listen for auth changes
    const {
      data: { subscription },
    } = supabase.auth.onAuthStateChange((_event, session) => {
      setSession(session)
      setUser(session?.user ?? null)
      setLoading(false)
    })

    return () => subscription.unsubscribe()
  }, [])

  const signIn = async (email: string, password: string) => {
    setLoading(true)
    setError(null)
    
    const { data, error } = await supabase.auth.signInWithPassword({
      email,
      password,
    })
    
    setError(error)
    setLoading(false)
    
    return { data, error }
  }

  const signUp = async (email: string, password: string) => {
    setLoading(true)
    setError(null)
    
    const { data, error } = await supabase.auth.signUp({
      email,
      password,
    })
    
    setError(error)
    setLoading(false)
    
    return { data, error }
  }

  const signOut = async () => {
    setLoading(true)
    const { error } = await supabase.auth.signOut()
    setError(error)
    setLoading(false)
    return { error }
  }

  const value = {
    user,
    session,
    loading,
    error,
    signIn,
    signUp,
    signOut,
  }

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}