'use client'

import { createContext, useContext, useEffect, useState } from 'react'
import { supabase } from '@/lib/supabase'
import type { User, AuthError } from '@supabase/supabase-js'

interface AuthContextType {
  user: User | null
  loading: boolean
  error: AuthError | null
  signIn: (email: string, password: string) => Promise<{ data: any; error: AuthError | null }>
  signUp: (email: string, password: string) => Promise<{ data: any; error: AuthError | null }>
  signOut: () => Promise<{ error: AuthError | null }>
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<AuthError | null>(null)

  useEffect(() => {
    // Get initial session
    supabase.auth.getSession().then(({ data: { session } }) => {
      setUser(session?.user ?? null)
      setLoading(false)
    })

    // Listen for auth changes
    const {
      data: { subscription },
    } = supabase.auth.onAuthStateChange((event, session) => {
      setUser(session?.user ?? null)
      setLoading(false)
      setError(null)
    })

    return () => subscription.unsubscribe()
  }, [])

  const signIn = async (email: string, password: string) => {
    setLoading(true)
    setError(null)
    
    const result = await supabase.auth.signInWithPassword({
      email,
      password,
    })
    
    setError(result.error)
    setLoading(false)
    
    return result
  }

  const signUp = async (email: string, password: string) => {
    setLoading(true)
    setError(null)
    
    const result = await supabase.auth.signUp({
      email,
      password,
    })
    
    setError(result.error)
    setLoading(false)
    
    return result
  }

  const signOut = async () => {
    setLoading(true)
    setError(null)
    
    const result = await supabase.auth.signOut()
    
    setError(result.error)
    setLoading(false)
    
    return result
  }

  const value = {
    user,
    loading,
    error,
    signIn,
    signUp,
    signOut,
  }

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}