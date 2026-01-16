'use client'

import { createContext, useContext, useEffect, useState, ReactNode } from 'react'
import { api, User } from '@/lib/api'

interface AuthContextType {
  user: User | null
  token: string | null
  isLoading: boolean
  login: (email: string, password: string) => Promise<void>
  register: (email: string, password: string) => Promise<void>
  logout: () => void
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [token, setToken] = useState<string | null>(null)
  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => {
    const savedToken = localStorage.getItem('auth_token')
    if (savedToken) {
      setToken(savedToken)
      api.setAuthToken(savedToken)
      api.getUser()
        .then(setUser)
        .catch(() => {
          localStorage.removeItem('auth_token')
          setToken(null)
          api.setAuthToken(null)
        })
        .finally(() => setIsLoading(false))
    } else {
      setIsLoading(false)
    }
  }, [])

  const login = async (email: string, password: string) => {
    const response = await api.login(email, password)
    setToken(response.token)
    setUser(response.user)
    localStorage.setItem('auth_token', response.token)
    api.setAuthToken(response.token)
  }

  const register = async (email: string, password: string) => {
    const response = await api.register(email, password)
    setToken(response.token)
    setUser(response.user)
    localStorage.setItem('auth_token', response.token)
    api.setAuthToken(response.token)
  }

  const logout = () => {
    setToken(null)
    setUser(null)
    localStorage.removeItem('auth_token')
    api.setAuthToken(null)
  }

  return (
    <AuthContext.Provider value={{ user, token, isLoading, login, register, logout }}>
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
