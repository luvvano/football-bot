import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { api, CreateEventPayload } from '../api/client'

export function useEvents() {
  return useQuery({
    queryKey: ['events'],
    queryFn: async () => {
      const response = await api.getEvents()
      return response.data
    },
  })
}

export function useEvent(id: number) {
  return useQuery({
    queryKey: ['events', id],
    queryFn: async () => {
      const response = await api.getEvent(id)
      return response.data
    },
    enabled: !!id,
  })
}

export function useCreateEvent() {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: async (data: CreateEventPayload) => {
      const response = await api.createEvent(data)
      return response.data
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['events'] })
    },
  })
}

export function useJoinEvent() {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: async (eventId: number) => {
      const response = await api.joinEvent(eventId)
      return response.data
    },
    onSuccess: (_, eventId) => {
      queryClient.invalidateQueries({ queryKey: ['events'] })
      queryClient.invalidateQueries({ queryKey: ['events', eventId] })
    },
  })
}

export function useLeaveEvent() {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: async (eventId: number) => {
      const response = await api.leaveEvent(eventId)
      return response.data
    },
    onSuccess: (_, eventId) => {
      queryClient.invalidateQueries({ queryKey: ['events'] })
      queryClient.invalidateQueries({ queryKey: ['events', eventId] })
    },
  })
}
