import axios from 'axios'
import { useAuthStore } from '../stores/auth'

const API_URL = import.meta.env.VITE_API_URL || '/api'

export const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Add auth header to all requests
api.interceptors.request.use((config) => {
  const initData = useAuthStore.getState().initData
  if (initData) {
    config.headers['X-Telegram-Init-Data'] = initData
  }
  return config
})

// Handle errors
api.interceptors.response.use(
  (response) => response,
  (error) => {
    console.error('API Error:', error.response?.data || error.message)
    return Promise.reject(error)
  }
)
