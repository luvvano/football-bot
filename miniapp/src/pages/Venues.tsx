import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { MapPin, Plus, ArrowLeft } from 'lucide-react'
import { useVenues, useCreateVenue } from '../hooks/useVenues'
import { useAuthStore } from '../store/auth'
import { useTelegram } from '../hooks/useTelegram'
import { Button } from '../components/Button'

export function Venues() {
  const navigate = useNavigate()
  const { isAdmin } = useAuthStore()
  const { setBackButtonVisible, onBackButtonClick, hapticNotification, showAlert } = useTelegram()
  const { data: venues, isLoading, error, refetch } = useVenues()
  const createVenue = useCreateVenue()
  
  const [isCreating, setIsCreating] = useState(false)
  const [name, setName] = useState('')
  const [address, setAddress] = useState('')

  useEffect(() => {
    setBackButtonVisible(true)
    const cleanup = onBackButtonClick(() => navigate('/'))
    return () => {
      setBackButtonVisible(false)
      cleanup()
    }
  }, [navigate, setBackButtonVisible, onBackButtonClick])

  const handleCreate = async () => {
    if (!name.trim() || !address.trim()) {
      showAlert('Please fill in all fields')
      return
    }

    try {
      await createVenue.mutateAsync({
        name: name.trim(),
        address: address.trim(),
      })
      hapticNotification('success')
      setIsCreating(false)
      setName('')
      setAddress('')
    } catch (error) {
      hapticNotification('error')
      showAlert('Failed to create venue')
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

      {/* Header */}
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold text-tg-text">Venues</h1>
        {isAdmin() && !isCreating && (
          <Button
            variant="primary"
            size="sm"
            onClick={() => setIsCreating(true)}
            className="flex items-center gap-1"
          >
            <Plus className="w-4 h-4" />
            Add
          </Button>
        )}
      </div>

      {/* Create Form */}
      {isCreating && (
        <div className="bg-tg-section-bg rounded-2xl p-4 mb-6">
          <h3 className="text-tg-text font-medium mb-4">New Venue</h3>
          
          <div className="space-y-4">
            <div>
              <label className="block text-tg-hint text-sm mb-1">Name</label>
              <input
                type="text"
                value={name}
                onChange={(e) => setName(e.target.value)}
                placeholder="Stadium name"
                className="w-full px-4 py-3 bg-tg-secondary-bg rounded-xl text-tg-text placeholder-tg-hint outline-none focus:ring-2 focus:ring-tg-button"
              />
            </div>
            
            <div>
              <label className="block text-tg-hint text-sm mb-1">Address</label>
              <input
                type="text"
                value={address}
                onChange={(e) => setAddress(e.target.value)}
                placeholder="Full address"
                className="w-full px-4 py-3 bg-tg-secondary-bg rounded-xl text-tg-text placeholder-tg-hint outline-none focus:ring-2 focus:ring-tg-button"
              />
            </div>
            
            <div className="flex gap-3">
              <Button
                variant="secondary"
                onClick={() => {
                  setIsCreating(false)
                  setName('')
                  setAddress('')
                }}
                className="flex-1"
              >
                Cancel
              </Button>
              <Button
                variant="primary"
                onClick={handleCreate}
                isLoading={createVenue.isPending}
                className="flex-1"
              >
                Create
              </Button>
            </div>
          </div>
        </div>
      )}

      {/* Venues List */}
      {isLoading ? (
        <div className="space-y-3">
          {[1, 2, 3].map(i => (
            <div 
              key={i}
              className="bg-tg-section-bg rounded-2xl h-20 animate-pulse"
            />
          ))}
        </div>
      ) : error ? (
        <div className="text-center py-8">
          <p className="text-tg-destructive mb-2">Failed to load venues</p>
          <Button variant="outline" size="sm" onClick={() => refetch()}>
            Try Again
          </Button>
        </div>
      ) : venues && venues.length === 0 ? (
        <div className="text-center py-12 bg-tg-section-bg rounded-2xl">
          <MapPin className="w-12 h-12 text-tg-hint mx-auto mb-3" />
          <p className="text-tg-text font-medium mb-1">No venues yet</p>
          <p className="text-tg-hint text-sm">
            {isAdmin() ? 'Add a venue to get started!' : 'Ask an admin to add venues'}
          </p>
        </div>
      ) : (
        <div className="space-y-3">
          {venues?.map((venue) => (
            <div
              key={venue.id}
              className="bg-tg-section-bg rounded-2xl p-4 flex items-center gap-3"
            >
              <div className="w-12 h-12 rounded-full bg-tg-button/10 flex items-center justify-center flex-shrink-0">
                <MapPin className="w-6 h-6 text-tg-button" />
              </div>
              <div className="flex-1 min-w-0">
                <p className="text-tg-text font-medium truncate">
                  {venue.name}
                </p>
                <p className="text-tg-hint text-sm truncate">
                  {venue.address}
                </p>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
