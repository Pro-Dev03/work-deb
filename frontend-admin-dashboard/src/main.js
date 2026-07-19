import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import './styles/tokens.css'
import i18n from './services/i18n'

import 'leaflet/dist/leaflet.css'
import L from 'leaflet'

delete L.Icon.Default.prototype._getIconUrl
L.Icon.Default.mergeOptions({
  iconRetinaUrl: new URL('leaflet/dist/images/marker-icon-2x.png', import.meta.url).href,
  iconUrl: new URL('leaflet/dist/images/marker-icon.png', import.meta.url).href,
  shadowUrl: new URL('leaflet/dist/images/marker-shadow.png', import.meta.url).href,
})

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

  // تسجيل Service Worker
  if ('serviceWorker' in navigator) {
    window.addEventListener('load', () => {
      navigator.serviceWorker.register('/service-worker.js')
        .then((registration) => {
          console.log('Service Worker registered with scope:', registration.scope)
        })
        .catch((error) => {
          console.error('Service Worker registration failed:', error)
        })
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
