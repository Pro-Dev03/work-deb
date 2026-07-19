import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import './styles/tokens.css'
import i18n from './services/i18n'

// PWA Install Handler
let deferredPrompt = null

// معالجة beforeinstallprompt event
window.addEventListener('beforeinstallprompt', (e) => {
  e.preventDefault()
  deferredPrompt = e
  console.log('PWA install prompt available')
  
  // إرسال event للتطبيق لإظهار زر التثبيت
  window.dispatchEvent(new CustomEvent('pwa-install-available', { detail: true }))
})

// معالجة appinstalled event
window.addEventListener('appinstalled', () => {
  deferredPrompt = null
  console.log('PWA was installed')
  window.dispatchEvent(new CustomEvent('pwa-install-success', { detail: true }))
})

// التأكد من تحميل الصفحة بشكل كامل قبل تعريف التطبيق
document.addEventListener('DOMContentLoaded', () => {
  const app = createApp(App)
  app.use(i18n)
  app.use(router)
  app.mount('#app')

  // تسجيل Service Worker مع التحقق من التحديثات
  if ('serviceWorker' in navigator) {
    window.addEventListener('load', () => {
      navigator.serviceWorker.register('/service-worker.js')
        .then((registration) => {
          console.log('Service Worker registered with scope:', registration.scope)

          // التحقق من التحديثات بشكل دوري
          setInterval(() => {
            registration.update()
              .then(() => {
                console.log('[SW] Update check completed')
              })
              .catch((error) => {
                console.error('[SW] Update check failed:', error)
              })
          }, 10000) // التحقق كل 10 ثواني

          // الاستماع للتحديثات الجديدة
          registration.addEventListener('updatefound', () => {
            const newWorker = registration.installing
            console.log('[SW] New service worker found')

            newWorker.addEventListener('statechange', () => {
              if (newWorker.state === 'installed' && navigator.serviceWorker.controller) {
                console.log('[SW] New service worker installed, activating')
                // تفعيل الـ service worker الجديد فوراً
                newWorker.postMessage({ type: 'SKIP_WAITING' })

                // إظهار إشعار للمستخدم بوجود تحديث
                if (confirm('🔄 تحديث متاح! هل تريد تحديث التطبيق الآن؟')) {
                  // إعادة تحميل الصفحة لتفعيل التحديث
                  window.location.reload()
                }
              }
            })
          })
        })
        .catch((error) => {
          console.error('Service Worker registration failed:', error)
        })
    })
    
    // الاستماع لتغييرات الـ service worker
    navigator.serviceWorker.addEventListener('controllerchange', () => {
      console.log('[SW] Controller changed, reloading page')
      window.location.reload()
    })
  }
})

// تصدير الدالة لاستخدامها في المكونات
window.pwaInstall = () => {
  if (deferredPrompt) {
    deferredPrompt.prompt()
    deferredPrompt.userChoice.then((choiceResult) => {
      if (choiceResult.outcome === 'accepted') {
        console.log('User accepted the install prompt')
      } else {
        console.log('User dismissed the install prompt')
      }
      deferredPrompt = null
    })
  }
}

// دالة لتحديث Service Worker يدوياً
window.forceUpdate = () => {
  if ('serviceWorker' in navigator) {
    navigator.serviceWorker.getRegistration().then((registration) => {
      if (registration) {
        registration.update()
        console.log('[SW] Manual update triggered')

        // إجبار جميع العملاء على التحديث
        if (registration.waiting) {
          registration.waiting.postMessage({ type: 'SKIP_WAITING' })
          console.log('[SW] Skipping waiting worker')
        }
      }
    })
  }
}

// دالة لإفراغ الكاش بالكامل
window.clearCache = () => {
  if ('caches' in window) {
    caches.keys().then((cacheNames) => {
      return Promise.all(
        cacheNames.map((cacheName) => {
          console.log('[SW] Deleting cache:', cacheName)
          return caches.delete(cacheName)
        })
      )
    }).then(() => {
      console.log('[SW] All caches cleared')
      window.location.reload(true) // Force reload
    })
  }
}

// إضافة إمكانية استدعاء من وحدة التحكم
console.log('🔧 Developer Tools Available:')
console.log('  - window.forceUpdate() : تحديث Service Worker')
console.log('  - window.clearCache()  : إفراغ الكاش وإعادة التحميل')
