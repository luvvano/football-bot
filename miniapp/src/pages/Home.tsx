import { useNavigate } from 'react-router-dom'
import { Plus, Calendar, MapPin } from 'lucide-react'
import { useEvents } from '../hooks/useEvents'
import { useAuthStore } from '../store/auth'
import { EventCard } from '../components/EventCard'
import { Button } from '../components/Button'

export function Home() {
  const navigate = useNavigate()
  const { user, community, isAdmin } = useAuthStore()
  const { data: events, isLoading, error, refetch } = useEvents()

  const upcomingEvents = events?.filter(
    e => new Date(e.dateTime) >= new Date()
  ).sort((a, b) => 
    new Date(a.dateTime).getTime() - new Date(b.dateTime).getTime()
  ) || []

  return (
    <div className="px-4 py-4">
      {/* Header */}
      <div className="mb-6">
        <h1 className="text-2xl font-bold text-tg-text mb-1">
          {community?.name || 'Football'}
        </h1>
        <p className="text-tg-hint">
          Welcome, {user?.firstName}! 👋
        </p>
      </div>

      {/* Quick Actions */}
      <div className="grid grid-cols-2 gap-3 mb-6">
        {isAdmin() && (
          <Button
            variant="primary"
            onClick={() => navigate('/create-event')}
            className="flex items-center justify-center gap-2"
          >
            <Plus className="w-5 h-5" />
            New Event
          </Button>
        )}
        <Button
          variant="secondary"
          onClick={() => navigate('/venues')}
          className="flex items-center justify-center gap-2"
        >
          <MapPin className="w-5 h-5" />
          Venues
        </Button>
      </div>

      {/* Events Section */}
      <div>
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-lg font-semibold text-tg-text flex items-center gap-2">
            <Calendar className="w-5 h-5" />
            Upcoming Events
          </h2>
          <button
            onClick={() => refetch()}
            className="text-tg-link text-sm"
          >
            Refresh
          </button>
        </div>

        {isLoading ? (
          <div className="space-y-3">
            {[1, 2, 3].map(i => (
              <div 
                key={i}
                className="bg-tg-section-bg rounded-2xl h-32 animate-pulse"
              />
            ))}
          </div>
        ) : error ? (
          <div className="text-center py-8">
            <p className="text-tg-destructive mb-2">Failed to load events</p>
            <Button variant="outline" size="sm" onClick={() => refetch()}>
              Try Again
            </Button>
          </div>
        ) : upcomingEvents.length === 0 ? (
          <div className="text-center py-12 bg-tg-section-bg rounded-2xl">
            <Calendar className="w-12 h-12 text-tg-hint mx-auto mb-3" />
            <p className="text-tg-text font-medium mb-1">No upcoming events</p>
            <p className="text-tg-hint text-sm">
              {isAdmin() ? 'Create one to get started!' : 'Check back later'}
            </p>
          </div>
        ) : (
          <div className="space-y-3">
            {upcomingEvents.map(event => (
              <EventCard key={event.id} event={event} />
            ))}
          </div>
        )}
      </div>
    </div>
  )
}
