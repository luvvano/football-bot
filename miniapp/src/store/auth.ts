import { create } from 'zustand'
import type { User, Community } from '../api/client'

interface AuthState {
  user: User | null
  community: Community | null
  isLoading: boolean
  error: string | null
  setUser: (user: User | null) => void
  setCommunity: (community: Community | null) => void
  setLoading: (loading: boolean) => void
  setError: (error: string | null) => void
  isAdmin: () => boolean
}

export const useAuthStore = create<AuthState>((set, get) => ({
  user: null,
  community: null,
  isLoading: true,
  error: null,
  
  setUser: (user) => set({ user }),
  setCommunity: (community) => set({ community }),
  setLoading: (isLoading) => set({ isLoading }),
  setError: (error) => set({ error }),
  
  isAdmin: () => get().user?.isAdmin ?? false,
}))
