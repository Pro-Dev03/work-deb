// Service Worker مع استراتيجية caching محسنة
const CACHE_NAME = 'worktrack-worker-v1'
const STATIC_CACHE = 'worktrack-static-v1'
const DYNAMIC_CACHE = 'worktrack-dynamic-v1'

// الملفات الثابتة التي يجب تخزينها فوراً
const STATIC_ASSETS = [
  '/',
  '/index.html',
  '/manifest.json',
  '/favicon.ico'
]

// تثبيت Service Worker
self.addEventListener('install', (event) => {
  console.log('[SW] Installing new service worker')
  event.waitUntil(
    caches.open(STATIC_CACHE).then((cache) => {
      console.log('[SW] Caching static assets')
      return cache.addAll(STATIC_ASSETS)
    })
  )
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
          // حذف الكاش القديم فقط
          if (cacheName !== STATIC_CACHE && cacheName !== DYNAMIC_CACHE) {
            console.log('[SW] Deleting old cache:', cacheName)
            return caches.delete(cacheName)
          }
        })
      )
    })
  )
  // التحكم في جميع الصفحات فوراً
  self.clients.claim()
})

// استراتيجية caching ذكية
self.addEventListener('fetch', (event) => {
  const { request } = event
  const url = new URL(request.url)

  // تجاهل طلبات غير HTTP
  if (!request.url.startsWith('http')) {
    return
  }

  // استراتيجية حسب نوع الملف
  if (request.method === 'GET') {
    // للملفات الثابتة (CSS, JS, Images, Fonts) - Cache First
    if (isStaticAsset(url)) {
      event.respondWith(cacheFirst(request))
    }
    // للـ HTML - Network First
    else if (url.pathname.endsWith('.html') || url.pathname === '/') {
      event.respondWith(networkFirst(request))
    }
    // للـ API - Network First
    else if (url.pathname.startsWith('/api/')) {
      event.respondWith(networkFirst(request))
    }
    // باقي الطلبات - Network First
    else {
      event.respondWith(networkFirst(request))
    }
  }
})

// التحقق من أن الملف ثابت
function isStaticAsset(url) {
  const staticExtensions = ['.css', '.js', '.png', '.jpg', '.jpeg', '.gif', '.svg', '.ico', '.woff', '.woff2', '.ttf', '.eot']
  return staticExtensions.some(ext => url.pathname.endsWith(ext)) ||
         url.pathname.includes('/assets/') ||
         url.pathname.includes('/images/')
}

// استراتيجية Cache First - أسرع للملفات الثابتة
async function cacheFirst(request) {
  try {
    const cachedResponse = await caches.match(request)
    if (cachedResponse) {
      // تحديث الكاش في الخلفية
      fetchAndCache(request)
      return cachedResponse
    }
    return await fetchAndCache(request)
  } catch (error) {
    console.error('[SW] Cache First failed:', error)
    throw error
  }
}

// استراتيجية Network First - أفضل للـ HTML و API
async function networkFirst(request) {
  try {
    const networkResponse = await fetch(request)
    if (networkResponse.ok) {
      // تخزين النسخة الجديدة
      const cache = await caches.open(DYNAMIC_CACHE)
      await cache.put(request, networkResponse.clone())
    }
    return networkResponse
  } catch (error) {
    console.log('[SW] Network failed, trying cache:', error)
    const cachedResponse = await caches.match(request)
    if (cachedResponse) {
      return cachedResponse
    }
    throw error
  }
}

// جلب من الشبكة وتخزين في الكاش
async function fetchAndCache(request) {
  const response = await fetch(request)
  if (response.ok) {
    const cache = await caches.open(STATIC_CACHE)
    await cache.put(request, response.clone())
  }
  return response
}

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
    body: event.data ? event.data.text() : 'WorkTrack',
    icon: '/favicon.ico',
    badge: '/favicon.ico',
    vibrate: [100, 50, 100],
    data: {
      dateOfArrival: Date.now(),
      primaryKey: 1
    }
  }

  event.waitUntil(
    self.registration.showNotification('WorkTrack', options)
  )
})
