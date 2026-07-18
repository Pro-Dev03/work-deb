// WebSocket service للتتبع اللحظي
class WebSocketService {
  constructor() {
    this.ws = null
    this.reconnectInterval = null
    this.listeners = []
    this.isConnected = false
    this.useElectronAPI = window.electronAPI && window.electronAPI.websocket
  }

  connect(url) {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      return
    }

    // استخدام Electron API إذا كان متاحاً
    if (this.useElectronAPI) {
      this.connectViaElectron(url)
    } else {
      this.connectViaBrowser(url)
    }
  }

  connectViaElectron(url) {
    try {
      console.log('🔌 جاري الاتصال بـ WebSocket عبر Electron:', url)
      
      // تنظيف المستمعين القديمين
      this.cleanupElectronListeners()

      // إعداد مستمعين جدد
      this.electronMessageHandler = (data) => {
        console.log('📡 استقبلت رسالة WebSocket (Electron):', data)
        this.notifyListeners(data)
      }

      this.electronOpenHandler = () => {
        console.log('✅ WebSocket متصل (Electron)')
        this.isConnected = true
        this.notifyListeners({ type: 'connected' })
      }

      this.electronErrorHandler = (error) => {
        console.warn('⚠️ WebSocket خطأ (Electron):', error)
        this.isConnected = false
      }

      this.electronCloseHandler = (code, reason) => {
        console.log('🔌 WebSocket مغلق (Electron):', code, reason)
        this.isConnected = false
        this.notifyListeners({ type: 'disconnected' })
        
        // إعادة الاتصال تلقائياً بعد 30 ثانية
        this.reconnectInterval = setTimeout(() => {
          console.log('🔄 محاولة إعادة الاتصال...')
          this.connect(url)
        }, 30000)
      }

      // تسجيل المستمعين
      window.electronAPI.websocket.onMessage(this.electronMessageHandler)
      window.electronAPI.websocket.onOpen(this.electronOpenHandler)
      window.electronAPI.websocket.onError(this.electronErrorHandler)
      window.electronAPI.websocket.onClose(this.electronCloseHandler)

      // الاتصال
      window.electronAPI.websocket.connect(url)
    } catch (e) {
      console.warn('⚠️ فشل الاتصال عبر Electron، استخدام المتصفح:', e)
      this.connectViaBrowser(url)
    }
  }

  connectViaBrowser(url) {
    try {
      console.log('🔌 جاري الاتصال بـ WebSocket عبر المتصفح:', url)
      this.ws = new WebSocket(url)

      this.ws.onopen = () => {
        console.log('✅ WebSocket متصل (المتصفح)')
        this.isConnected = true
        this.notifyListeners({ type: 'connected' })
      }

      this.ws.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data)
          console.log('📡 استقبلت رسالة WebSocket (المتصفح):', data)
          this.notifyListeners(data)
        } catch (e) {
          console.error('❌ خطأ في تحليل رسالة WebSocket:', e)
        }
      }

      this.ws.onerror = (error) => {
        console.warn('⚠️ WebSocket غير متاح - سيتم استخدام التحديث الدوري')
        this.isConnected = false
      }

      this.ws.onclose = (event) => {
        console.log('🔌 WebSocket مغلق (المتصفح)', event.code, event.reason)
        this.isConnected = false
        this.notifyListeners({ type: 'disconnected' })
        
        // إعادة الاتصال تلقائياً بعد 30 ثانية
        this.reconnectInterval = setTimeout(() => {
          console.log('🔄 محاولة إعادة الاتصال...')
          this.connect(url)
        }, 30000)
      }
    } catch (e) {
      console.warn('⚠️ WebSocket غير متاح - سيتم استخدام التحديث الدوري')
    }
  }

  cleanupElectronListeners() {
    if (!this.useElectronAPI) return
    
    // إزالة المستمعين (تحتاج إلى implementation في preload.js)
    // حالياً لا توجد طريقة لإزالة المستمعين في Electron IPC
    // يمكن إضافتها لاحقاً إذا لزم الأمر
  }

  disconnect() {
    if (this.reconnectInterval) {
      clearTimeout(this.reconnectInterval)
      this.reconnectInterval = null
    }

    if (this.useElectronAPI) {
      window.electronAPI.websocket.disconnect()
    } else if (this.ws) {
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
    if (this.useElectronAPI) {
      window.electronAPI.websocket.send(data)
    } else if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(data))
    } else {
      console.warn('⚠️ WebSocket غير متصل')
    }
  }
}

// إنشاء instance واحد
const wsService = new WebSocketService()

export default wsService