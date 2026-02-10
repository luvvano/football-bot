import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { api, CreateVenuePayload } from '../api/client'

export function useVenues() {
  return useQuery({
    queryKey: ['venues'],
    queryFn: async () => {
      const response = await api.getVenues()
      return response.data
    },
  })
}

export function useCreateVenue() {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: async (data: CreateVenuePayload) => {
      const response = await api.createVenue(data)
      return response.data
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['venues'] })
    },
  })
}
