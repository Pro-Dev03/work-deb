import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import './styles/tokens.css'
import i18n from './services/i18n'

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
