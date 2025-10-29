export function useWebSocket(playerId) {
  let ws = null
  const listeners = {
    open: [],
    close: [],
    error: [],
    message: []
  }
  let reconnectAttempts = 0
  const maxReconnectAttempts = 5
  const reconnectDelay = 3000

  const connect = () => {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const wsUrl = `${protocol}//${window.location.host}/ws?playerId=${playerId}`
    
    ws = new WebSocket(wsUrl)

    ws.onopen = () => {
      reconnectAttempts = 0
      listeners.open.forEach(listener => listener())
    }

    ws.onclose = () => {
      listeners.close.forEach(listener => listener())
      
      if (reconnectAttempts < maxReconnectAttempts) {
        setTimeout(() => {
          reconnectAttempts++
          connect()
        }, reconnectDelay)
      }
    }

    ws.onerror = () => {
      listeners.error.forEach(listener => listener())
    }

    ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data)
        listeners.message.forEach(listener => listener(data))
      } catch (error) {
        console.error('Error parsing WebSocket message:', error)
      }
    }
  }

  const send = (data) => {
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(data)
    }
  }

  const close = () => {
    if (ws) {
      ws.close()
      ws = null
    }
  }

  const isConnected = () => {
    return ws && ws.readyState === WebSocket.OPEN
  }

  const on = (event, callback) => {
    if (listeners[event]) {
      listeners[event].push(callback)
    }
  }

  const off = (event, callback) => {
    if (listeners[event]) {
      listeners[event] = listeners[event].filter(cb => cb !== callback)
    }
  }

  connect()

  return {
    send,
    close,
    isConnected,
    on,
    off
  }
}
