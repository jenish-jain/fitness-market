import axios from 'axios'
import { supabase } from './supabase'

const api = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1',
  timeout: 10000,
})

// Request interceptor to add auth token
api.interceptors.request.use(
  async (config) => {
    // Get the current session from Supabase
    const { data: { session } } = await supabase.auth.getSession()
    
    if (session?.access_token) {
      config.headers.Authorization = `Bearer ${session.access_token}`
    }
    
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Response interceptor for error handling
api.interceptors.response.use(
  (response) => {
    return response
  },
  async (error) => {
    if (error.response?.status === 401) {
      // Handle unauthorized access - sign out user
      await supabase.auth.signOut()
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

export default api