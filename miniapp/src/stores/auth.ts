import { create } from 'zustand'
import type { TelegramUser } from '../types'

interface AuthState {
  user: TelegramUser | null
  initData: string | null
  isAuthenticated: boolean
  init: (user: TelegramUser, initData: string) => void
  logout: () => void
}

export const useAuthStore = create<AuthState>((set) => ({
  user: null,
  initData: null,
  isAuthenticated: false,
  
  init: (user, initData) => {
    set({
      user,
      initData,
      isAuthenticated: true,
    })
  },
  
  logout: () => {
    set({
      user: null,
      initData: null,
      isAuthenticated: false,
    })
  },
}))
