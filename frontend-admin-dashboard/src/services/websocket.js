// WebSocket service للتتبع اللحظي
class WebSocketService {
  constructor() {
    this.ws = null
    this.reconnectInterval = null
    this.listeners = []
    this.isConnected = false
  }

  connect(url) {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      return
    }

    try {
      console.log('🔌 جاري الاتصال بـ WebSocket:', url)
      this.ws = new WebSocket(url)

      this.ws.onopen = () => {
        console.log('✅ WebSocket متصل')
        this.isConnected = true
        this.notifyListeners({ type: 'connected' })
      }

      this.ws.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data)
          console.log('📡 استقبلت رسالة WebSocket:', data)
          this.notifyListeners(data)
        } catch (e) {
          console.error('❌ خطأ في تحليل رسالة WebSocket:', e)
        }
      }

      this.ws.onerror = (error) => {
        console.error('❌ خطأ WebSocket:', error)
        this.isConnected = false
      }

      this.ws.onclose = (event) => {
        console.log('🔌 WebSocket مغلق', event.code, event.reason)
        this.isConnected = false
        this.notifyListeners({ type: 'disconnected' })
        
        // إعادة الاتصال تلقائياً بعد 5 ثوانٍ
        this.reconnectInterval = setTimeout(() => {
          console.log('🔄 محاولة إعادة الاتصال...')
          this.connect(url)
        }, 5000)
      }
    } catch (e) {
      console.error('❌ فشل الاتصال بـ WebSocket:', e)
    }
  }

  disconnect() {
    if (this.reconnectInterval) {
      clearTimeout(this.reconnectInterval)
      this.reconnectInterval = null
    }

    if (this.ws) {
      this.ws.close()
      this.ws = null
    }

    this.isConnected = false
  }

  onMessage(callback) {
    this.listeners.push(callback)
  }

  removeListener(callback) {
    this.listeners = this.listeners.filter(listener => listener !== callback)
  }

  notifyListeners(data) {
    this.listeners.forEach(callback => callback(data))
  }

  send(data) {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(data))
    } else {
      console.warn('⚠️ WebSocket غير متصل')
    }
  }
}

// إنشاء instance واحد
const wsService = new WebSocketService()

export default wsService