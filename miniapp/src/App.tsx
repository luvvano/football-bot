import { useEffect } from 'react'
import { Routes, Route } from 'react-router-dom'
import WebApp from '@twa-dev/sdk'
import { useAuthStore } from './stores/auth'
import HomePage from './pages/HomePage'
import EventPage from './pages/EventPage'
import VenuesPage from './pages/VenuesPage'

function App() {
  const { init } = useAuthStore()

  useEffect(() => {
    // Initialize Telegram WebApp
    WebApp.ready()
    WebApp.expand()
    
    // Set theme
    document.documentElement.classList.add('dark')
    
    // Initialize auth from Telegram
    if (WebApp.initDataUnsafe?.user) {
      init(WebApp.initDataUnsafe.user, WebApp.initData)
    }
  }, [init])

  return (
    <div className="min-h-screen bg-gray-900 text-white">
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/event/:id" element={<EventPage />} />
        <Route path="/venues" element={<VenuesPage />} />
      </Routes>
    </div>
  )
}

export default App
