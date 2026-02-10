import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { api } from './client'
import type { Event, Participant } from '../types'

// Get events for community
export function useEvents(communityId?: number) {
  return useQuery({
    queryKey: ['events', communityId],
    queryFn: async () => {
      const params = communityId ? { community_id: communityId } : {}
      const { data } = await api.get<Event[]>('/events', { params })
      return data
    },
    enabled: true,
  })
}

// Get single event with participants
export function useEvent(eventId: number) {
  return useQuery({
    queryKey: ['event', eventId],
    queryFn: async () => {
      const { data } = await api.get<Event>(`/events/${eventId}`)
      return data
    },
    enabled: !!eventId,
  })
}

// Get event participants
export function useEventParticipants(eventId: number) {
  return useQuery({
    queryKey: ['event', eventId, 'participants'],
    queryFn: async () => {
      const { data } = await api.get<Participant[]>(`/events/${eventId}/participants`)
      return data
    },
    enabled: !!eventId,
  })
}

// Register for event
export function useRegister() {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: async (eventId: number) => {
      const { data } = await api.post(`/events/${eventId}/register`)
      return data
    },
    onSuccess: (_, eventId) => {
      queryClient.invalidateQueries({ queryKey: ['event', eventId] })
      queryClient.invalidateQueries({ queryKey: ['events'] })
    },
  })
}

// Unregister from event
export function useUnregister() {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: async (eventId: number) => {
      const { data } = await api.post(`/events/${eventId}/unregister`)
      return data
    },
    onSuccess: (_, eventId) => {
      queryClient.invalidateQueries({ queryKey: ['event', eventId] })
      queryClient.invalidateQueries({ queryKey: ['events'] })
    },
  })
}
