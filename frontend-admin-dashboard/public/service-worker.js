// Service Worker - Network Only Strategy
// This ensures all code changes appear immediately without cache delays

const CACHE_NAME = 'worktrack-admin-v2'

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
          // حذف جميع الكاش القديم
          console.log('[SW] Deleting cache:', cacheName)
          return caches.delete(cacheName)
        })
      )
    })
  )
  // التحكم في جميع الصفحات فوراً
  self.clients.claim()
})

// استراتيجية Network Only - جلب كل شيء من الشبكة مباشرة
self.addEventListener('fetch', (event) => {
  const { request } = event

  // تجاهل طلبات غير HTTP
  if (!request.url.startsWith('http')) {
    return
  }

  // جلب جميع الموارد من الشبكة مباشرة
  if (request.method === 'GET') {
    event.respondWith(
      fetch(request).catch((error) => {
        console.error('[SW] Network fetch failed:', error)
        throw error
      })
    )
  }
})

// معالجة طلبات التحديث عبر pull-to-refresh
self.addEventListener('message', (event) => {
  if (event.data && event.data.type === 'SKIP_WAITING') {
    console.log('[SW] Skip waiting requested')
    self.skipWaiting()
  }
  // مسح الكاش عند الطلب
  if (event.data && event.data.type === 'CLEAR_CACHE') {
    console.log('[SW] Clear cache requested')
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