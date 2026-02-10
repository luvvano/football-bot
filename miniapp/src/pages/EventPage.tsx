import { useParams, useNavigate } from 'react-router-dom'
import { ArrowLeft, Calendar, Clock, MapPin, Users, UserPlus, UserMinus } from 'lucide-react'
import WebApp from '@twa-dev/sdk'
import { useEvent, useEventParticipants, useRegister, useUnregister } from '../api/events'
import { useAuthStore } from '../stores/auth'
import { format, parseISO } from 'date-fns'
import { ru } from 'date-fns/locale'
import type { Participant } from '../types'

function ParticipantItem({ participant, index }: { participant: Participant; index: number }) {
  const name = participant.user
    ? participant.user.first_name + (participant.user.last_name ? ` ${participant.user.last_name}` : '')
    : participant.guest_name || 'Гость'
  
  const isGuest = !participant.user_id

  return (
    <div className="flex items-center gap-3 py-2">
      <span className="w-6 h-6 bg-gray-700 rounded-full flex items-center justify-center text-xs text-gray-400">
        {index + 1}
      </span>
      <span className="flex-1 text-white">
        {name}
        {isGuest && <span className="text-gray-500 text-sm ml-2">(гость)</span>}
      </span>
    </div>
  )
}

export default function EventPage() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const eventId = Number(id)
  const { user } = useAuthStore()
  
  const { data: event, isLoading: eventLoading } = useEvent(eventId)
  const { data: participants, isLoading: participantsLoading } = useEventParticipants(eventId)
  const registerMutation = useRegister()
  const unregisterMutation = useUnregister()
  
  const isLoading = eventLoading || participantsLoading
  const isRegistered = participants?.some(p => p.user?.telegram_id === user?.id)
  const maxParticipants = event ? (event.team_size + event.max_subs_per_team) * 2 : 0
  const currentCount = participants?.length || 0
  const isFull = currentCount >= maxParticipants
  const canRegister = event?.status === 'open' && !isFull && !isRegistered
  const canUnregister = event?.status === 'open' && isRegistered

  const handleRegister = async () => {
    try {
      await registerMutation.mutateAsync(eventId)
      WebApp.HapticFeedback.notificationOccurred('success')
    } catch {
      WebApp.HapticFeedback.notificationOccurred('error')
    }
  }

  const handleUnregister = async () => {
    WebApp.showConfirm('Отписаться от игры?', async (confirmed) => {
      if (confirmed) {
        try {
          await unregisterMutation.mutateAsync(eventId)
          WebApp.HapticFeedback.notificationOccurred('success')
        } catch {
          WebApp.HapticFeedback.notificationOccurred('error')
        }
      }
    })
  }

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="w-8 h-8 border-2 border-purple-500 border-t-transparent rounded-full animate-spin" />
      </div>
    )
  }

  if (!event) {
    return (
      <div className="p-4 text-center">
        <p className="text-gray-400">Игра не найдена</p>
      </div>
    )
  }

  return (
    <div className="min-h-screen pb-24">
      {/* Header */}
      <div className="p-4 safe-area-top">
        <button
          onClick={() => navigate(-1)}
          className="flex items-center gap-2 text-gray-400 hover:text-white transition-colors mb-4"
        >
          <ArrowLeft className="w-5 h-5" />
          <span>Назад</span>
        </button>

        <h1 className="text-2xl font-bold mb-4">⚽ Игра</h1>

        {/* Event Info */}
        <div className="bg-gray-800 rounded-xl p-4 space-y-3">
          <div className="flex items-center gap-3">
            <Calendar className="w-5 h-5 text-purple-400" />
            <span className="text-white">
              {format(parseISO(event.event_date), 'd MMMM, EEEE', { locale: ru })}
            </span>
          </div>

          <div className="flex items-center gap-3">
            <Clock className="w-5 h-5 text-purple-400" />
            <span className="text-white">{event.event_time.slice(0, 5)}</span>
          </div>

          {event.venue && (
            <div className="flex items-center gap-3">
              <MapPin className="w-5 h-5 text-purple-400" />
              <div>
                <span className="text-white">{event.venue.name}</span>
                {event.venue.address && (
                  <p className="text-gray-400 text-sm">{event.venue.address}</p>
                )}
              </div>
            </div>
          )}

          <div className="flex items-center gap-3">
            <Users className="w-5 h-5 text-purple-400" />
            <span className={isFull ? 'text-red-400' : 'text-green-400'}>
              {currentCount}/{maxParticipants} игроков
            </span>
            <span className="text-gray-500 text-sm">
              ({event.team_size}x{event.team_size} + {event.max_subs_per_team} замены)
            </span>
          </div>
        </div>
      </div>

      {/* Participants */}
      <div className="p-4">
        <h2 className="text-lg font-semibold mb-3 flex items-center gap-2">
          <Users className="w-5 h-5" />
          Участники
        </h2>
        
        {participants && participants.length > 0 ? (
          <div className="bg-gray-800 rounded-xl p-4 divide-y divide-gray-700">
            {participants.map((p, i) => (
              <ParticipantItem key={p.id} participant={p} index={i} />
            ))}
          </div>
        ) : (
          <div className="bg-gray-800 rounded-xl p-6 text-center text-gray-400">
            <p>Пока никого нет</p>
            <p className="text-sm mt-1">Будьте первым!</p>
          </div>
        )}
      </div>

      {/* Action Button */}
      <div className="fixed bottom-0 left-0 right-0 p-4 bg-gray-900/90 backdrop-blur safe-area-bottom">
        {canRegister && (
          <button
            onClick={handleRegister}
            disabled={registerMutation.isPending}
            className="w-full bg-purple-600 hover:bg-purple-700 text-white font-semibold py-3 px-6 rounded-xl flex items-center justify-center gap-2 btn-press disabled:opacity-50"
          >
            <UserPlus className="w-5 h-5" />
            {registerMutation.isPending ? 'Записываем...' : 'Записаться'}
          </button>
        )}

        {canUnregister && (
          <button
            onClick={handleUnregister}
            disabled={unregisterMutation.isPending}
            className="w-full bg-red-600 hover:bg-red-700 text-white font-semibold py-3 px-6 rounded-xl flex items-center justify-center gap-2 btn-press disabled:opacity-50"
          >
            <UserMinus className="w-5 h-5" />
            {unregisterMutation.isPending ? 'Отписываем...' : 'Отписаться'}
          </button>
        )}

        {isRegistered && event.status !== 'open' && (
          <div className="text-center text-gray-400">
            Вы записаны ✓
          </div>
        )}

        {isFull && !isRegistered && (
          <div className="text-center text-red-400">
            Игра заполнена
          </div>
        )}
      </div>
    </div>
  )
}
