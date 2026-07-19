// سيتم تحديث هذا الإصدار تلقائياً عند تغيير package.json
const CACHE_NAME = 'worktrack-client-v__APP_VERSION__'
const urlsToCache = [
  '/',
  '/manifest.json'
]

// تثبيت Service Worker وتخزين الملفات الأساسية
self.addEventListener('install', (event) => {
  event.waitUntil(
    caches.open(CACHE_NAME)
      .then((cache) => {
        console.log('Opened cache')
        return cache.addAll(urlsToCache)
      })
      .catch((error) => {
        console.error('Cache installation failed:', error)
      })
  )
  self.skipWaiting()
})

// تفعيل Service Worker وتنظيف الكاش القديم
self.addEventListener('activate', (event) => {
  event.waitUntil(
    caches.keys().then((cacheNames) => {
      return Promise.all(
        cacheNames.map((cacheName) => {
          if (cacheName !== CACHE_NAME) {
            console.log('Deleting old cache:', cacheName)
            return caches.delete(cacheName)
          }
        })
      )
    })
  )
  self.clients.claim()
})

// استقبال الطلبات وتخزينها
self.addEventListener('fetch', (event) => {
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

// معالجة الإشعارات (للاستخدام المستقبلي)
self.addEventListener('push', (event) => {
  const options = {
    body: event.data ? event.data.text() : 'WorkTrack Client',
    icon: '/favicon.ico',
    badge: '/favicon.ico',
    vibrate: [100, 50, 100],
    data: {
      dateOfArrival: Date.now(),
      primaryKey: 1
    }
  }

  event.waitUntil(
    self.registration.showNotification('WorkTrack Client', options)
  )
})