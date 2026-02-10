import { useNavigate } from 'react-router-dom'
import { format } from 'date-fns'
import { Calendar, MapPin, Users, ChevronRight } from 'lucide-react'
import type { Event } from '../api/client'
import { useAuthStore } from '../store/auth'
import { useJoinEvent, useLeaveEvent } from '../hooks/useEvents'
import { Button } from './Button'
import { useTelegram } from '../hooks/useTelegram'

interface EventCardProps {
  event: Event
}

export function EventCard({ event }: EventCardProps) {
  const navigate = useNavigate()
  const { user } = useAuthStore()
  const { hapticNotification } = useTelegram()
  const joinEvent = useJoinEvent()
  const leaveEvent = useLeaveEvent()

  const confirmedParticipants = event.participants.filter(p => p.status === 'confirmed')
  const isParticipating = event.participants.some(
    p => p.userId === user?.id && p.status === 'confirmed'
  )
  const isFull = confirmedParticipants.length >= event.maxParticipants

  const handleJoin = async (e: React.MouseEvent) => {
    e.stopPropagation()
    try {
      await joinEvent.mutateAsync(event.id)
      hapticNotification('success')
    } catch (error) {
      hapticNotification('error')
    }
  }

  const handleLeave = async (e: React.MouseEvent) => {
    e.stopPropagation()
    try {
      await leaveEvent.mutateAsync(event.id)
      hapticNotification('success')
    } catch (error) {
      hapticNotification('error')
    }
  }

  const eventDate = new Date(event.dateTime)
  const isPast = eventDate < new Date()

  return (
    <div
      onClick={() => navigate(`/event/${event.id}`)}
      className="bg-tg-section-bg rounded-2xl p-4 shadow-sm active:bg-tg-secondary-bg transition-colors cursor-pointer"
    >
      <div className="flex justify-between items-start mb-3">
        <h3 className="text-lg font-semibold text-tg-text">{event.title}</h3>
        <ChevronRight className="w-5 h-5 text-tg-hint flex-shrink-0" />
      </div>

      <div className="space-y-2 mb-4">
        <div className="flex items-center gap-2 text-sm text-tg-subtitle">
          <Calendar className="w-4 h-4" />
          <span>{format(eventDate, 'EEE, MMM d · HH:mm')}</span>
        </div>

        <div className="flex items-center gap-2 text-sm text-tg-subtitle">
          <MapPin className="w-4 h-4" />
          <span>{event.venue.name}</span>
        </div>

        <div className="flex items-center gap-2 text-sm text-tg-subtitle">
          <Users className="w-4 h-4" />
          <span>
            {confirmedParticipants.length}/{event.maxParticipants} players
          </span>
          {isFull && !isParticipating && (
            <span className="text-tg-destructive text-xs ml-1">(Full)</span>
          )}
        </div>
      </div>

      {!isPast && (
        <div onClick={(e) => e.stopPropagation()}>
          {isParticipating ? (
            <Button
              variant="outline"
              size="sm"
              fullWidth
              onClick={handleLeave}
              isLoading={leaveEvent.isPending}
            >
              Leave
            </Button>
          ) : (
            <Button
              variant="primary"
              size="sm"
              fullWidth
              onClick={handleJoin}
              isLoading={joinEvent.isPending}
              disabled={isFull}
            >
              {isFull ? 'Full' : 'Join'}
            </Button>
          )}
        </div>
      )}
    </div>
  )
}
