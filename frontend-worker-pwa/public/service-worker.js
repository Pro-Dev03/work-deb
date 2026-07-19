// سيتم تحديث هذا الإصدار تلقائياً عند تغيير package.json
const CACHE_NAME = 'worktrack-v__APP_VERSION__'
const urlsToCache = [
  '/',
  '/manifest.json',
  '/index.html',
  '/icon-128x128.png',
  '/icon-192x192.png',
  '/icon-256x256.png',
  '/icon-512x512.png',
  '/favicon.ico'
]

// الملفات التي يجب عدم تخزينها في الكاش (JavaScript/CSS dynamic)
const DYNAMIC_CACHE_PATTERNS = [
  /\.js$/,
  /\.css$/,
  /\/assets\//i
]

// دعم offline للصفحات الرئيسية
const CACHE_PAGES = [
  '/',
  '/index.html'
]

// تثبيت Service Worker وتخزين الملفات الأساسية
self.addEventListener('install', (event) => {
  console.log('[SW] Installing new service worker')
  event.waitUntil(
    caches.open(CACHE_NAME)
      .then((cache) => {
        console.log('[SW] Opened cache:', CACHE_NAME)
        return cache.addAll(urlsToCache)
      })
      .catch((error) => {
        console.error('[SW] Cache installation failed:', error)
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
          if (cacheName !== CACHE_NAME) {
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

// استقبال الطلبات وتخزينها
self.addEventListener('fetch', (event) => {
  const url = new URL(event.request.url)

  // التحقق مما إذا كان الطلب لملف ديناميكي (JS/CSS/assets)
  const isDynamic = DYNAMIC_CACHE_PATTERNS.some(pattern => pattern.test(url.pathname))

  // للملفات الديناميكية، جلب دائماً من الشبكة (لا كاش)
  if (isDynamic) {
    event.respondWith(
      fetch(event.request)
        .then(response => {
          // التحقق من صحة الاستجابة
          if (!response || response.status !== 200) {
            return response
          }

          // استنساخ الاستجابة لأنها يمكن استخدامها مرة واحدة فقط
          const responseToCache = response.clone()

          // تخزين الاستجابة في الكاش للاستخدام offline
          caches.open(CACHE_NAME).then(cache => {
            cache.put(event.request, responseToCache)
          })

          return response
        })
        .catch(() => {
          // في حالة فشل الشبكة، محاولة العودة إلى الكاش
          return caches.match(event.request)
        })
    )
    return
  }

  // للملفات الثابتة، استخدام الكاش أولاً
  event.respondWith(
    caches.match(event.request)
      .then((response) => {
        // إذا كان الطلب موجود في الكاش، إرجاعه
        if (response) {
          return response
        }

        // خلاف ذلك، جلب الطلب من الشبكة
        return fetch(event.request).then((response) => {
          // التحقق من صحة الاستجابة
          if (!response || response.status !== 200 || response.type !== 'basic') {
            return response
          }

          // استنساخ الاستجابة لأنها يمكن استخدامها مرة واحدة فقط
          const responseToCache = response.clone()

          // تخزين الاستجابة في الكاش
          caches.open(CACHE_NAME).then((cache) => {
            cache.put(event.request, responseToCache)
          })

          return response
        }).catch(() => {
          // في حالة فشل الشبكة، محاولة العودة إلى الكاش
          return caches.match(event.request)
        })
      })
  )
})

// معالجة طلبات التحديث عبر pull-to-refresh
self.addEventListener('message', (event) => {
  if (event.data && event.data.type === 'SKIP_WAITING') {
    console.log('[SW] Skip waiting requested')
    self.skipWaiting()
  }
  
  if (event.data && event.data.type === 'CACHE_BUST') {
    console.log('[SW] Cache bust requested')
    // تنظيف الكاش القديم
    caches.keys().then((cacheNames) => {
      return Promise.all(
        cacheNames.map((cacheName) => {
          if (cacheName !== CACHE_NAME) {
            return caches.delete(cacheName)
          }
        })
      )
    })
  }
  
  if (event.data && event.data.type === 'CHECK_UPDATE') {
    console.log('[SW] Update check requested')
    // إرسال إشعار بأن هناك تحديث متاح
    event.ports[0].postMessage({ type: 'UPDATE_AVAILABLE' })
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
