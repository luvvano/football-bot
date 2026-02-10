import axios from 'axios'
import WebApp from '@twa-dev/sdk'

const API_BASE_URL = import.meta.env.VITE_API_URL || '/api'

export const apiClient = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Add TWA auth header to all requests
apiClient.interceptors.request.use((config) => {
  if (WebApp.initData) {
    config.headers['X-Telegram-Init-Data'] = WebApp.initData
  }
  return config
})

// Handle response errors
apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // Handle unauthorized
      console.error('Unauthorized - invalid or expired session')
    }
    return Promise.reject(error)
  }
)

// API types
export interface User {
  id: number
  telegramId: number
  firstName: string
  lastName?: string
  username?: string
  isAdmin: boolean
}

export interface Community {
  id: number
  name: string
  telegramChatId: number
}

export interface Venue {
  id: number
  name: string
  address: string
  communityId: number
}

export interface Participant {
  id: number
  userId: number
  eventId: number
  status: 'confirmed' | 'waitlist' | 'cancelled'
  user: User
}

export interface Event {
  id: number
  title: string
  description?: string
  dateTime: string
  venueId: number
  communityId: number
  maxParticipants: number
  createdById: number
  venue: Venue
  participants: Participant[]
}

export interface CreateEventPayload {
  title: string
  description?: string
  dateTime: string
  venueId: number
  maxParticipants: number
}

export interface CreateVenuePayload {
  name: string
  address: string
}

// API functions
export const api = {
  // Events
  getEvents: () => apiClient.get<Event[]>('/events'),
  getEvent: (id: number) => apiClient.get<Event>(`/events/${id}`),
  createEvent: (data: CreateEventPayload) => apiClient.post<Event>('/events', data),
  joinEvent: (id: number) => apiClient.post(`/events/${id}/join`),
  leaveEvent: (id: number) => apiClient.post(`/events/${id}/leave`),
  
  // Venues
  getVenues: () => apiClient.get<Venue[]>('/venues'),
  createVenue: (data: CreateVenuePayload) => apiClient.post<Venue>('/venues', data),
  
  // Auth
  authenticate: (initData: string) => 
    apiClient.post<{ user: User; community: Community }>('/auth/telegram', { initData }),
}
