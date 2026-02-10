import { useEffect } from 'react'
import { Routes, Route } from 'react-router-dom'
import WebApp from '@twa-dev/sdk'
import { Layout } from './components/Layout'
import { Home } from './pages/Home'
import { Event } from './pages/Event'
import { Venues } from './pages/Venues'
import { CreateEvent } from './pages/CreateEvent'
import { useAuthStore } from './store/auth'
import { apiClient } from './api/client'

function App() {
  const { setUser, setCommunity, setLoading, setError } = useAuthStore()

  useEffect(() => {
    // Initialize TWA
    WebApp.ready()
    WebApp.expand()
    
    // Set theme colors
    document.documentElement.style.setProperty('--tg-theme-bg-color', WebApp.themeParams.bg_color || '#ffffff')
    document.documentElement.style.setProperty('--tg-theme-text-color', WebApp.themeParams.text_color || '#000000')
    document.documentElement.style.setProperty('--tg-theme-hint-color', WebApp.themeParams.hint_color || '#999999')
    document.documentElement.style.setProperty('--tg-theme-link-color', WebApp.themeParams.link_color || '#2481cc')
    document.documentElement.style.setProperty('--tg-theme-button-color', WebApp.themeParams.button_color || '#2481cc')
    document.documentElement.style.setProperty('--tg-theme-button-text-color', WebApp.themeParams.button_text_color || '#ffffff')
    document.documentElement.style.setProperty('--tg-theme-secondary-bg-color', WebApp.themeParams.secondary_bg_color || '#f0f0f0')

    // Authenticate with backend
    const authenticate = async () => {
      try {
        setLoading(true)
        const response = await apiClient.post('/auth/telegram', {
          initData: WebApp.initData,
        })
        
        setUser(response.data.user)
        setCommunity(response.data.community)
      } catch (error) {
        console.error('Auth failed:', error)
        setError('Authentication failed')
      } finally {
        setLoading(false)
      }
    }

    if (WebApp.initData) {
      authenticate()
    } else {
      setLoading(false)
      setError('Not running in Telegram')
    }
  }, [setUser, setCommunity, setLoading, setError])

  return (
    <Layout>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/event/:id" element={<Event />} />
        <Route path="/venues" element={<Venues />} />
        <Route path="/create-event" element={<CreateEvent />} />
      </Routes>
    </Layout>
  )
}

export default App
