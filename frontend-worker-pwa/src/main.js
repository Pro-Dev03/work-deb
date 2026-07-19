import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import './styles/tokens.css'
import i18n from './services/i18n'

// PWA Installation Logic
let deferredPrompt;
window.addEventListener('beforeinstallprompt', (e) => {
  e.preventDefault();
  deferredPrompt = e;
  window.deferredPrompt = deferredPrompt;
  window.dispatchEvent(new Event('pwa-install-available'));
});

window.addEventListener('appinstalled', () => {
  deferredPrompt = null;
  window.deferredPrompt = null;
  window.dispatchEvent(new Event('pwa-install-success'));
});

window.pwaInstall = async () => {
  if (deferredPrompt) {
    deferredPrompt.prompt();
    const { outcome } = await deferredPrompt.userChoice;
    if (outcome === 'accepted') {
      console.log('PWA installation accepted');
    } else {
      console.log('PWA installation dismissed');
    }
    deferredPrompt = null;
    window.deferredPrompt = null;
  }
};

// Service Worker Registration
if ('serviceWorker' in navigator) {
  window.addEventListener('load', () => {
    navigator.serviceWorker.register('/service-worker.js')
      .then(registration => {
        console.log('Service Worker registered:', registration);
      })
      .catch(error => {
        console.log('Service Worker registration failed:', error);
      });
  });
}

// التأكد من تحميل الصفحة بشكل كامل قبل تعريف التطبيق
document.addEventListener('DOMContentLoaded', () => {
  const app = createApp(App)
  app.use(i18n)
  app.use(router)
  app.mount('#app')
})
