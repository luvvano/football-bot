export interface User {
  id: number
  telegram_id: number
  username?: string
  first_name?: string
  last_name?: string
}

export interface Community {
  id: number
  telegram_chat_id: number
  name: string
  settings: CommunitySettings
}

export interface CommunitySettings {
  team_size: number
  max_subs_per_team: number
  timezone: string
  rating_params: string[]
}

export interface Venue {
  id: number
  community_id: number
  name: string
  address?: string
}

export interface Event {
  id: number
  community_id: number
  venue_id?: number
  event_date: string
  event_time: string
  team_size: number
  max_subs_per_team: number
  status: EventStatus
  team_a: number[]
  team_b: number[]
  subs_a: number[]
  subs_b: number[]
  venue?: Venue
  participants?: Participant[]
  participant_count?: number
}

export type EventStatus = 'open' | 'full' | 'in_progress' | 'finished' | 'cancelled' | 'not_held'

export interface Participant {
  id: number
  event_id: number
  user_id?: number
  guest_name?: string
  status: string
  registered_at: string
  user?: User
}

export interface TelegramUser {
  id: number
  first_name: string
  last_name?: string
  username?: string
  language_code?: string
  is_premium?: boolean
}
