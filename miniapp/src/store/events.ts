import { create } from 'zustand'
import type { Event } from '../api/client'

interface EventsState {
  events: Event[]
  currentEvent: Event | null
  setEvents: (events: Event[]) => void
  setCurrentEvent: (event: Event | null) => void
  updateEvent: (event: Event) => void
}

export const useEventsStore = create<EventsState>((set) => ({
  events: [],
  currentEvent: null,
  
  setEvents: (events) => set({ events }),
  setCurrentEvent: (currentEvent) => set({ currentEvent }),
  
  updateEvent: (updatedEvent) => 
    set((state) => ({
      events: state.events.map((e) => 
        e.id === updatedEvent.id ? updatedEvent : e
      ),
      currentEvent: state.currentEvent?.id === updatedEvent.id 
        ? updatedEvent 
        : state.currentEvent,
    })),
}))
