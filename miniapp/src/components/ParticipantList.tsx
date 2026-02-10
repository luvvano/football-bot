import { User, Clock } from 'lucide-react'
import type { Participant } from '../api/client'
import { useAuthStore } from '../store/auth'

interface ParticipantListProps {
  participants: Participant[]
  maxParticipants: number
}

export function ParticipantList({ participants, maxParticipants }: ParticipantListProps) {
  const { user } = useAuthStore()

  const confirmed = participants.filter(p => p.status === 'confirmed')
  const waitlist = participants.filter(p => p.status === 'waitlist')

  const renderParticipant = (participant: Participant, index: number) => {
    const isCurrentUser = participant.userId === user?.id
    const displayName = participant.user.firstName + 
      (participant.user.lastName ? ` ${participant.user.lastName}` : '')

    return (
      <div
        key={participant.id}
        className={`flex items-center gap-3 py-2.5 px-3 rounded-xl ${
          isCurrentUser ? 'bg-tg-button/10' : ''
        }`}
      >
        <div className="w-8 h-8 rounded-full bg-tg-secondary-bg flex items-center justify-center flex-shrink-0">
          <User className="w-4 h-4 text-tg-hint" />
        </div>
        <div className="flex-1 min-w-0">
          <p className="text-tg-text font-medium truncate">
            {displayName}
            {isCurrentUser && (
              <span className="text-tg-button text-sm ml-1">(You)</span>
            )}
          </p>
          {participant.user.username && (
            <p className="text-tg-hint text-sm truncate">
              @{participant.user.username}
            </p>
          )}
        </div>
        <span className="text-tg-hint text-sm flex-shrink-0">
          #{index + 1}
        </span>
      </div>
    )
  }

  return (
    <div className="space-y-4">
      {/* Confirmed participants */}
      <div>
        <h4 className="text-tg-section-header text-sm font-medium uppercase mb-2 px-1">
          Players ({confirmed.length}/{maxParticipants})
        </h4>
        <div className="bg-tg-section-bg rounded-2xl divide-y divide-tg-secondary-bg">
          {confirmed.length > 0 ? (
            confirmed.map((p, i) => renderParticipant(p, i))
          ) : (
            <div className="py-8 text-center text-tg-hint">
              No players yet
            </div>
          )}
        </div>
      </div>

      {/* Waitlist */}
      {waitlist.length > 0 && (
        <div>
          <h4 className="text-tg-section-header text-sm font-medium uppercase mb-2 px-1 flex items-center gap-1">
            <Clock className="w-4 h-4" />
            Waitlist ({waitlist.length})
          </h4>
          <div className="bg-tg-section-bg rounded-2xl divide-y divide-tg-secondary-bg">
            {waitlist.map((p, i) => renderParticipant(p, confirmed.length + i))}
          </div>
        </div>
      )}
    </div>
  )
}
