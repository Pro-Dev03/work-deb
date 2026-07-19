// Service Worker بدون كاش - يعتمد على الشبكة فقط

// تثبيت Service Worker
self.addEventListener('install', (event) => {
  console.log('[SW] Installing new service worker')
  // تفعيل الـ service worker الجديد فوراً
  self.skipWaiting()
})

// تفعيل Service Worker وتنظيف الكاش القديم
self.addEventListener('activate', (event) => {
  console.log('[SW] Activating new service worker')
  event.waitUntil(
    caches.keys().then((cacheNames) => {
      return Promise.all(
        cacheNames.map((cacheName) => {
          console.log('[SW] Deleting cache:', cacheName)
          return caches.delete(cacheName)
        })
      )
    })
  )
  // التحكم في جميع الصفحات فوراً
  self.clients.claim()
})

// استقبال الطلبات وجلبها من الشبكة فقط (بدون كاش)
self.addEventListener('fetch', (event) => {
  event.respondWith(
    fetch(event.request)
      .then(response => {
        return response
      })
      .catch(error => {
        console.error('[SW] Fetch failed:', error)
        throw error
      })
  )
})

// معالجة طلبات التحديث عبر pull-to-refresh
self.addEventListener('message', (event) => {
  if (event.data && event.data.type === 'SKIP_WAITING') {
    console.log('[SW] Skip waiting requested')
    self.skipWaiting()
  }
})

// معالجة الإشعارات (للاستخدام المستقبلي)
self.addEventListener('push', (event) => {
  const options = {
    body: event.data ? event.data.text() : 'WorkTrack Admin',
    icon: '/favicon.ico',
    badge: '/favicon.ico',
    vibrate: [100, 50, 100],
    data: {
      dateOfArrival: Date.now(),
      primaryKey: 1
    }
  }

  event.waitUntil(
    self.registration.showNotification('WorkTrack Admin', options)
  )
})