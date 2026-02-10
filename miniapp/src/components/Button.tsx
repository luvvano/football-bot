import { ButtonHTMLAttributes, ReactNode } from 'react'
import { Loader2 } from 'lucide-react'
import { useTelegram } from '../hooks/useTelegram'

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  children: ReactNode
  variant?: 'primary' | 'secondary' | 'outline' | 'destructive'
  size?: 'sm' | 'md' | 'lg'
  isLoading?: boolean
  fullWidth?: boolean
}

export function Button({
  children,
  variant = 'primary',
  size = 'md',
  isLoading = false,
  fullWidth = false,
  disabled,
  onClick,
  className = '',
  ...props
}: ButtonProps) {
  const { hapticFeedback } = useTelegram()

  const baseClasses = 'inline-flex items-center justify-center font-medium rounded-xl transition-all active:scale-[0.98] disabled:opacity-50 disabled:cursor-not-allowed'

  const variantClasses = {
    primary: 'bg-tg-button text-tg-button-text',
    secondary: 'bg-tg-secondary-bg text-tg-text',
    outline: 'bg-transparent border-2 border-tg-button text-tg-button',
    destructive: 'bg-tg-destructive text-white',
  }

  const sizeClasses = {
    sm: 'px-3 py-1.5 text-sm',
    md: 'px-4 py-2.5 text-base',
    lg: 'px-6 py-3 text-lg',
  }

  const handleClick = (e: React.MouseEvent<HTMLButtonElement>) => {
    hapticFeedback('light')
    onClick?.(e)
  }

  return (
    <button
      className={`
        ${baseClasses}
        ${variantClasses[variant]}
        ${sizeClasses[size]}
        ${fullWidth ? 'w-full' : ''}
        ${className}
      `}
      disabled={disabled || isLoading}
      onClick={handleClick}
      {...props}
    >
      {isLoading ? (
        <>
          <Loader2 className="w-4 h-4 mr-2 animate-spin" />
          Loading...
        </>
      ) : (
        children
      )}
    </button>
  )
}
