import { useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { format } from 'date-fns'
import { Calendar, MapPin, Users, ArrowLeft } from 'lucide-react'
import { useEvent, useJoinEvent, useLeaveEvent } from '../hooks/useEvents'
import { useAuthStore } from '../store/auth'
import { useTelegram } from '../hooks/useTelegram'
import { ParticipantList } from '../components/ParticipantList'
import { Button } from '../components/Button'

export function Event() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const { user } = useAuthStore()
  const { setBackButtonVisible, onBackButtonClick, hapticNotification } = useTelegram()
  const { data: event, isLoading, error } = useEvent(Number(id))
  const joinEvent = useJoinEvent()
  const leaveEvent = useLeaveEvent()

  useEffect(() => {
    setBackButtonVisible(true)
    const cleanup = onBackButtonClick(() => navigate('/'))
    return () => {
      setBackButtonVisible(false)
      cleanup()
    }
  }, [navigate, setBackButtonVisible, onBackButtonClick])

  if (isLoading) {
    return (
      <div className="px-4 py-4">
        <div className="animate-pulse space-y-4">
          <div className="h-8 bg-tg-secondary-bg rounded-xl w-3/4" />
          <div className="h-4 bg-tg-secondary-bg rounded-lg w-1/2" />
          <div className="h-32 bg-tg-secondary-bg rounded-2xl" />
        </div>
      </div>
    )
  }

  if (error || !event) {
    return (
      <div className="px-4 py-4 text-center">
        <p className="text-tg-destructive mb-4">Event not found</p>
        <Button variant="outline" onClick={() => navigate('/')}>
          Go Back
        </Button>
      </div>
    )
  }

  const confirmedParticipants = event.participants.filter(p => p.status === 'confirmed')
  const isParticipating = event.participants.some(
    p => p.userId === user?.id && p.status === 'confirmed'
  )
  const isFull = confirmedParticipants.length >= event.maxParticipants
  const eventDate = new Date(event.dateTime)
  const isPast = eventDate < new Date()

  const handleJoin = async () => {
    try {
      await joinEvent.mutateAsync(event.id)
      hapticNotification('success')
    } catch (error) {
      hapticNotification('error')
    }
  }

  const handleLeave = async () => {
    try {
      await leaveEvent.mutateAsync(event.id)
      hapticNotification('success')
    } catch (error) {
      hapticNotification('error')
    }
  }

  return (
    <div className="px-4 py-4">
      {/* Back button for non-TWA */}
      <button
        onClick={() => navigate('/')}
        className="flex items-center gap-1 text-tg-link mb-4 md:hidden"
      >
        <ArrowLeft className="w-4 h-4" />
        Back
      </button>

      {/* Event Header */}
      <div className="mb-6">
        <h1 className="text-2xl font-bold text-tg-text mb-2">
          {event.title}
        </h1>
        {event.description && (
          <p className="text-tg-subtitle">{event.description}</p>
        )}
      </div>

      {/* Event Details */}
      <div className="bg-tg-section-bg rounded-2xl p-4 mb-6 space-y-3">
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 rounded-full bg-tg-button/10 flex items-center justify-center">
            <Calendar className="w-5 h-5 text-tg-button" />
          </div>
          <div>
            <p className="text-tg-text font-medium">
              {format(eventDate, 'EEEE, MMMM d, yyyy')}
            </p>
            <p className="text-tg-hint text-sm">
              {format(eventDate, 'HH:mm')}
            </p>
          </div>
        </div>

        <div className="flex items-center gap-3">
          <div className="w-10 h-10 rounded-full bg-tg-button/10 flex items-center justify-center">
            <MapPin className="w-5 h-5 text-tg-button" />
          </div>
          <div>
            <p className="text-tg-text font-medium">{event.venue.name}</p>
            <p className="text-tg-hint text-sm">{event.venue.address}</p>
          </div>
        </div>

        <div className="flex items-center gap-3">
          <div className="w-10 h-10 rounded-full bg-tg-button/10 flex items-center justify-center">
            <Users className="w-5 h-5 text-tg-button" />
          </div>
          <div>
            <p className="text-tg-text font-medium">
              {confirmedParticipants.length}/{event.maxParticipants} players
            </p>
            <p className="text-tg-hint text-sm">
              {isFull ? 'Event is full' : `${event.maxParticipants - confirmedParticipants.length} spots left`}
            </p>
          </div>
        </div>
      </div>

      {/* Join/Leave Button */}
      {!isPast && (
        <div className="mb-6">
          {isParticipating ? (
            <Button
              variant="outline"
              fullWidth
              onClick={handleLeave}
              isLoading={leaveEvent.isPending}
            >
              Leave Event
            </Button>
          ) : (
            <Button
              variant="primary"
              fullWidth
              onClick={handleJoin}
              isLoading={joinEvent.isPending}
              disabled={isFull}
            >
              {isFull ? 'Event Full' : 'Join Event'}
            </Button>
          )}
        </div>
      )}

      {isPast && (
        <div className="mb-6 text-center py-3 bg-tg-secondary-bg rounded-xl">
          <p className="text-tg-hint">This event has ended</p>
        </div>
      )}

      {/* Participants */}
      <ParticipantList
        participants={event.participants}
        maxParticipants={event.maxParticipants}
      />
    </div>
  )
}
