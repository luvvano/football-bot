import { MapPin, Check } from 'lucide-react'
import { useVenues } from '../hooks/useVenues'
import type { Venue } from '../api/client'

interface VenueSelectorProps {
  selectedId: number | null
  onSelect: (venue: Venue) => void
}

export function VenueSelector({ selectedId, onSelect }: VenueSelectorProps) {
  const { data: venues, isLoading, error } = useVenues()

  if (isLoading) {
    return (
      <div className="py-8 text-center text-tg-hint">
        Loading venues...
      </div>
    )
  }

  if (error) {
    return (
      <div className="py-8 text-center text-tg-destructive">
        Failed to load venues
      </div>
    )
  }

  if (!venues || venues.length === 0) {
    return (
      <div className="py-8 text-center text-tg-hint">
        No venues available
      </div>
    )
  }

  return (
    <div className="space-y-2">
      {venues.map((venue) => {
        const isSelected = venue.id === selectedId

        return (
          <button
            key={venue.id}
            onClick={() => onSelect(venue)}
            className={`w-full flex items-center gap-3 p-4 rounded-2xl text-left transition-all active:scale-[0.98] ${
              isSelected
                ? 'bg-tg-button/10 border-2 border-tg-button'
                : 'bg-tg-section-bg border-2 border-transparent'
            }`}
          >
            <div className={`w-10 h-10 rounded-full flex items-center justify-center flex-shrink-0 ${
              isSelected ? 'bg-tg-button' : 'bg-tg-secondary-bg'
            }`}>
              <MapPin className={`w-5 h-5 ${
                isSelected ? 'text-tg-button-text' : 'text-tg-hint'
              }`} />
            </div>
            <div className="flex-1 min-w-0">
              <p className={`font-medium truncate ${
                isSelected ? 'text-tg-button' : 'text-tg-text'
              }`}>
                {venue.name}
              </p>
              <p className="text-tg-hint text-sm truncate">
                {venue.address}
              </p>
            </div>
            {isSelected && (
              <Check className="w-5 h-5 text-tg-button flex-shrink-0" />
            )}
          </button>
        )
      })}
    </div>
  )
}
