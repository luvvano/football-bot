import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { format } from 'date-fns'
import { ArrowLeft, Calendar, Users, FileText } from 'lucide-react'
import { useCreateEvent } from '../hooks/useEvents'
import { useAuthStore } from '../store/auth'
import { useTelegram } from '../hooks/useTelegram'
import { VenueSelector } from '../components/VenueSelector'
import { Button } from '../components/Button'
import type { Venue } from '../api/client'

export function CreateEvent() {
  const navigate = useNavigate()
  const { isAdmin } = useAuthStore()
  const { setBackButtonVisible, onBackButtonClick, hapticNotification, showAlert } = useTelegram()
  const createEvent = useCreateEvent()
  
  const [title, setTitle] = useState('')
  const [description, setDescription] = useState('')
  const [date, setDate] = useState('')
  const [time, setTime] = useState('')
  const [maxParticipants, setMaxParticipants] = useState('10')
  const [selectedVenue, setSelectedVenue] = useState<Venue | null>(null)
  const [step, setStep] = useState<'details' | 'venue'>('details')

  useEffect(() => {
    // Redirect non-admins
    if (!isAdmin()) {
      navigate('/')
      return
    }

    setBackButtonVisible(true)
    const cleanup = onBackButtonClick(() => {
      if (step === 'venue') {
        setStep('details')
      } else {
        navigate('/')
      }
    })
    return () => {
      setBackButtonVisible(false)
      cleanup()
    }
  }, [navigate, setBackButtonVisible, onBackButtonClick, isAdmin, step])

  // Set default date to tomorrow
  useEffect(() => {
    const tomorrow = new Date()
    tomorrow.setDate(tomorrow.getDate() + 1)
    setDate(format(tomorrow, 'yyyy-MM-dd'))
    setTime('19:00')
  }, [])

  const handleSubmit = async () => {
    if (!title.trim()) {
      showAlert('Please enter a title')
      return
    }
    if (!date || !time) {
      showAlert('Please select date and time')
      return
    }
    if (!selectedVenue) {
      showAlert('Please select a venue')
      return
    }

    const dateTime = new Date(`${date}T${time}:00`)
    if (dateTime <= new Date()) {
      showAlert('Event must be in the future')
      return
    }

    try {
      await createEvent.mutateAsync({
        title: title.trim(),
        description: description.trim() || undefined,
        dateTime: dateTime.toISOString(),
        venueId: selectedVenue.id,
        maxParticipants: parseInt(maxParticipants) || 10,
      })
      hapticNotification('success')
      navigate('/')
    } catch (error) {
      hapticNotification('error')
      showAlert('Failed to create event')
    }
  }

  if (step === 'venue') {
    return (
      <div className="px-4 py-4">
        {/* Back button for non-TWA */}
        <button
          onClick={() => setStep('details')}
          className="flex items-center gap-1 text-tg-link mb-4 md:hidden"
        >
          <ArrowLeft className="w-4 h-4" />
          Back
        </button>

        <h1 className="text-2xl font-bold text-tg-text mb-6">Select Venue</h1>
        
        <VenueSelector
          selectedId={selectedVenue?.id ?? null}
          onSelect={(venue) => {
            setSelectedVenue(venue)
            setStep('details')
          }}
        />
      </div>
    )
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

      <h1 className="text-2xl font-bold text-tg-text mb-6">Create Event</h1>

      <div className="space-y-6">
        {/* Title */}
        <div>
          <label className="flex items-center gap-2 text-tg-text font-medium mb-2">
            <FileText className="w-4 h-4" />
            Title
          </label>
          <input
            type="text"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            placeholder="Friday Football"
            className="w-full px-4 py-3 bg-tg-section-bg rounded-xl text-tg-text placeholder-tg-hint outline-none focus:ring-2 focus:ring-tg-button"
          />
        </div>

        {/* Description */}
        <div>
          <label className="block text-tg-text font-medium mb-2">
            Description (optional)
          </label>
          <textarea
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            placeholder="Add details about the event..."
            rows={3}
            className="w-full px-4 py-3 bg-tg-section-bg rounded-xl text-tg-text placeholder-tg-hint outline-none focus:ring-2 focus:ring-tg-button resize-none"
          />
        </div>

        {/* Date & Time */}
        <div>
          <label className="flex items-center gap-2 text-tg-text font-medium mb-2">
            <Calendar className="w-4 h-4" />
            Date & Time
          </label>
          <div className="grid grid-cols-2 gap-3">
            <input
              type="date"
              value={date}
              onChange={(e) => setDate(e.target.value)}
              min={format(new Date(), 'yyyy-MM-dd')}
              className="px-4 py-3 bg-tg-section-bg rounded-xl text-tg-text outline-none focus:ring-2 focus:ring-tg-button"
            />
            <input
              type="time"
              value={time}
              onChange={(e) => setTime(e.target.value)}
              className="px-4 py-3 bg-tg-section-bg rounded-xl text-tg-text outline-none focus:ring-2 focus:ring-tg-button"
            />
          </div>
        </div>

        {/* Max Participants */}
        <div>
          <label className="flex items-center gap-2 text-tg-text font-medium mb-2">
            <Users className="w-4 h-4" />
            Max Players
          </label>
          <input
            type="number"
            value={maxParticipants}
            onChange={(e) => setMaxParticipants(e.target.value)}
            min="2"
            max="50"
            className="w-full px-4 py-3 bg-tg-section-bg rounded-xl text-tg-text outline-none focus:ring-2 focus:ring-tg-button"
          />
        </div>

        {/* Venue Selection */}
        <div>
          <label className="block text-tg-text font-medium mb-2">
            Venue
          </label>
          {selectedVenue ? (
            <button
              onClick={() => setStep('venue')}
              className="w-full p-4 bg-tg-section-bg rounded-xl text-left flex items-center justify-between"
            >
              <div>
                <p className="text-tg-text font-medium">{selectedVenue.name}</p>
                <p className="text-tg-hint text-sm">{selectedVenue.address}</p>
              </div>
              <span className="text-tg-link text-sm">Change</span>
            </button>
          ) : (
            <Button
              variant="secondary"
              fullWidth
              onClick={() => setStep('venue')}
            >
              Select Venue
            </Button>
          )}
        </div>

        {/* Submit */}
        <Button
          variant="primary"
          fullWidth
          size="lg"
          onClick={handleSubmit}
          isLoading={createEvent.isPending}
        >
          Create Event
        </Button>
      </div>
    </div>
  )
}
