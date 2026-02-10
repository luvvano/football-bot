import { useCallback, useEffect, useState } from 'react'
import WebApp from '@twa-dev/sdk'

export function useTelegram() {
  const [isReady, setIsReady] = useState(false)

  useEffect(() => {
    WebApp.ready()
    setIsReady(true)
  }, [])

  const showAlert = useCallback((message: string) => {
    WebApp.showAlert(message)
  }, [])

  const showConfirm = useCallback((message: string): Promise<boolean> => {
    return new Promise((resolve) => {
      WebApp.showConfirm(message, (confirmed) => {
        resolve(confirmed)
      })
    })
  }, [])

  const showPopup = useCallback((params: Parameters<typeof WebApp.showPopup>[0]): Promise<string | null> => {
    return new Promise((resolve) => {
      WebApp.showPopup(params, (buttonId) => {
        resolve(buttonId ?? null)
      })
    })
  }, [])

  const hapticFeedback = useCallback((type: 'light' | 'medium' | 'heavy' | 'rigid' | 'soft') => {
    WebApp.HapticFeedback.impactOccurred(type)
  }, [])

  const hapticNotification = useCallback((type: 'error' | 'success' | 'warning') => {
    WebApp.HapticFeedback.notificationOccurred(type)
  }, [])

  const close = useCallback(() => {
    WebApp.close()
  }, [])

  const setBackButtonVisible = useCallback((visible: boolean) => {
    if (visible) {
      WebApp.BackButton.show()
    } else {
      WebApp.BackButton.hide()
    }
  }, [])

  const onBackButtonClick = useCallback((callback: () => void) => {
    WebApp.BackButton.onClick(callback)
    return () => WebApp.BackButton.offClick(callback)
  }, [])

  const setMainButtonParams = useCallback((params: {
    text?: string
    color?: string
    textColor?: string
    isVisible?: boolean
    isActive?: boolean
    isLoading?: boolean
  }) => {
    if (params.text) WebApp.MainButton.setText(params.text)
    if (params.color) WebApp.MainButton.setParams({ color: params.color })
    if (params.textColor) WebApp.MainButton.setParams({ text_color: params.textColor })
    if (params.isVisible !== undefined) {
      params.isVisible ? WebApp.MainButton.show() : WebApp.MainButton.hide()
    }
    if (params.isActive !== undefined) {
      params.isActive ? WebApp.MainButton.enable() : WebApp.MainButton.disable()
    }
    if (params.isLoading !== undefined) {
      params.isLoading ? WebApp.MainButton.showProgress() : WebApp.MainButton.hideProgress()
    }
  }, [])

  const onMainButtonClick = useCallback((callback: () => void) => {
    WebApp.MainButton.onClick(callback)
    return () => WebApp.MainButton.offClick(callback)
  }, [])

  return {
    isReady,
    initData: WebApp.initData,
    initDataUnsafe: WebApp.initDataUnsafe,
    themeParams: WebApp.themeParams,
    colorScheme: WebApp.colorScheme,
    platform: WebApp.platform,
    version: WebApp.version,
    showAlert,
    showConfirm,
    showPopup,
    hapticFeedback,
    hapticNotification,
    close,
    setBackButtonVisible,
    onBackButtonClick,
    setMainButtonParams,
    onMainButtonClick,
  }
}
