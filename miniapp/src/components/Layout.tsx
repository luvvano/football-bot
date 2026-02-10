import { ReactNode } from 'react'
import { useAuthStore } from '../store/auth'
import { Loader2 } from 'lucide-react'

interface LayoutProps {
  children: ReactNode
}

export function Layout({ children }: LayoutProps) {
  const { isLoading, error } = useAuthStore()

  if (isLoading) {
    return (
      <div className="min-h-screen bg-tg-bg flex items-center justify-center">
        <Loader2 className="w-8 h-8 text-tg-button animate-spin" />
      </div>
    )
  }

  if (error) {
    return (
      <div className="min-h-screen bg-tg-bg flex items-center justify-center p-4">
        <div className="text-center">
          <p className="text-tg-destructive mb-2">Error</p>
          <p className="text-tg-hint text-sm">{error}</p>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-tg-bg safe-area-bottom">
      <main className="pb-4">
        {children}
      </main>
    </div>
  )
}
