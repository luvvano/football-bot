import { useSearchParams, useNavigate } from 'react-router-dom'
import { ArrowLeft, MapPin } from 'lucide-react'
import { useVenues } from '../api/venues'
import type { Venue } from '../types'

function VenueCard({ venue }: { venue: Venue }) {
  return (
    <div className="bg-gray-800 rounded-xl p-4 mb-3">
      <div className="flex items-start gap-3">
        <div className="w-10 h-10 bg-purple-600/20 rounded-lg flex items-center justify-center">
          <MapPin className="w-5 h-5 text-purple-400" />
        </div>
        <div className="flex-1">
          <h3 className="text-white font-medium">{venue.name}</h3>
          {venue.address && (
            <p className="text-gray-400 text-sm mt-1">{venue.address}</p>
          )}
        </div>
      </div>
    </div>
  )
}

export default function VenuesPage() {
  const [searchParams] = useSearchParams()
  const navigate = useNavigate()
  const communityId = searchParams.get('community') ? Number(searchParams.get('community')) : undefined
  
  const { data: venues, isLoading, error } = useVenues(communityId)

  return (
    <div className="p-4 safe-area-top safe-area-bottom">
      {/* Header */}
      <button
        onClick={() => navigate(-1)}
        className="flex items-center gap-2 text-gray-400 hover:text-white transition-colors mb-4"
      >
        <ArrowLeft className="w-5 h-5" />
        <span>Назад</span>
      </button>

      <h1 className="text-2xl font-bold mb-6">🏟️ Площадки</h1>

      {/* Venues List */}
      {isLoading ? (
        <div className="flex items-center justify-center py-12">
          <div className="w-8 h-8 border-2 border-purple-500 border-t-transparent rounded-full animate-spin" />
        </div>
      ) : error ? (
        <div className="text-center py-12 text-gray-400">
          <p>Не удалось загрузить площадки</p>
        </div>
      ) : venues && venues.length > 0 ? (
        <div>
          {venues.map((venue) => (
            <VenueCard key={venue.id} venue={venue} />
          ))}
        </div>
      ) : (
        <div className="text-center py-12 text-gray-400">
          <p className="text-4xl mb-4">🏟️</p>
          <p>Площадки не созданы</p>
          <p className="text-sm mt-2">Создайте площадку командой /venue в группе</p>
        </div>
      )}
    </div>
  )
}
