// حجز أولي لملف Service Worker — سيُستخدم لاحقاً لدعم العمل دون اتصال (offline)
self.addEventListener('install', () => self.skipWaiting())
self.addEventListener('activate', () => self.clients.claim())
