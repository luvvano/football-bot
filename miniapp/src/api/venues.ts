import { useQuery } from '@tanstack/react-query'
import { api } from './client'
import type { Venue } from '../types'

// Get venues for community
export function useVenues(communityId?: number) {
  return useQuery({
    queryKey: ['venues', communityId],
    queryFn: async () => {
      const params = communityId ? { community_id: communityId } : {}
      const { data } = await api.get<Venue[]>('/venues', { params })
      return data
    },
  })
}
