import { useSearchParams, Link } from 'react-router-dom'
import { Calendar, Users, MapPin } from 'lucide-react'
import { useEvents } from '../api/events'
import { useAuthStore } from '../stores/auth'
import { format, parseISO } from 'date-fns'
import { ru } from 'date-fns/locale'
import type { Event } from '../types'

function EventCard({ event }: { event: Event }) {
  const maxParticipants = (event.team_size + event.max_subs_per_team) * 2
  const currentCount = event.participant_count || 0
  const isFull = currentCount >= maxParticipants
  
  const statusColors = {
    open: 'bg-green-500',
    full: 'bg-red-500',
    in_progress: 'bg-yellow-500',
    finished: 'bg-gray-500',
    cancelled: 'bg-red-800',
    not_held: 'bg-gray-600',
  }

  return (
    <Link
      to={`/event/${event.id}`}
      className="block bg-gray-800 rounded-xl p-4 mb-3 btn-press hover:bg-gray-750 transition-colors"
    >
      <div className="flex items-start justify-between">
        <div className="flex-1">
          <div className="flex items-center gap-2 mb-2">
            <Calendar className="w-4 h-4 text-purple-400" />
            <span className="text-white font-medium">
              {format(parseISO(event.event_date), 'd MMMM', { locale: ru })}
            </span>
            <span className="text-gray-400">в {event.event_time.slice(0, 5)}</span>
          </div>
          
          {event.venue && (
            <div className="flex items-center gap-2 mb-2 text-sm text-gray-400">
              <MapPin className="w-4 h-4" />
              <span>{event.venue.name}</span>
            </div>
          )}
          
          <div className="flex items-center gap-2">
            <Users className="w-4 h-4 text-gray-400" />
            <span className={isFull ? 'text-red-400' : 'text-green-400'}>
              {currentCount}/{maxParticipants}
            </span>
            <span className="text-gray-500 text-sm">игроков</span>
          </div>
        </div>
        
        <div className={`w-3 h-3 rounded-full ${statusColors[event.status]}`} />
      </div>
    </Link>
  )
}

export default function HomePage() {
  const [searchParams] = useSearchParams()
  const communityId = searchParams.get('community') ? Number(searchParams.get('community')) : undefined
  const { user } = useAuthStore()
  
  const { data: events, isLoading, error } = useEvents(communityId)

  return (
    <div className="p-4 safe-area-top safe-area-bottom">
      {/* Header */}
      <div className="mb-6">
        <h1 className="text-2xl font-bold mb-1">⚽ Игры</h1>
        {user && (
          <p className="text-gray-400 text-sm">
            Привет, {user.first_name}!
          </p>
        )}
      </div>

      {/* Events List */}
      {isLoading ? (
        <div className="flex items-center justify-center py-12">
          <div className="w-8 h-8 border-2 border-purple-500 border-t-transparent rounded-full animate-spin" />
        </div>
      ) : error ? (
        <div className="text-center py-12 text-gray-400">
          <p>Не удалось загрузить игры</p>
          <p className="text-sm mt-2">Попробуйте позже</p>
        </div>
      ) : events && events.length > 0 ? (
        <div>
          {events.map((event) => (
            <EventCard key={event.id} event={event} />
          ))}
        </div>
      ) : (
        <div className="text-center py-12 text-gray-400">
          <p className="text-4xl mb-4">🏟️</p>
          <p>Нет предстоящих игр</p>
          <p className="text-sm mt-2">Создайте игру командой /event в группе</p>
        </div>
      )}
    </div>
  )
}
