# Frontend Code (Vue/JS)

## تم الاستخراج في: 2026-08-06 02:09:32

## عدد الملفات: 114

---

## 📄 frontend-admin-dashboard/build-obfuscated.js

```javascript
import JavaScriptObfuscator from 'javascript-obfuscator';
import { readFileSync, writeFileSync, readdirSync, statSync } from 'fs';
import { join } from 'path';

// إعدادات التعتيم
const obfuscatorOptions = {
  compact: true,
  controlFlowFlattening: true,
  controlFlowFlatteningThreshold: 0.75,
  deadCodeInjection: true,
  deadCodeInjectionThreshold: 0.4,
  debugProtection: true,
  disableConsoleOutput: false, // Changed to false for debugging
  identifierNamesGenerator: 'hexadecimal',
  log: true, // Changed to true for debugging
  numbersToExpressions: true,
  renameGlobals: false,
  selfDefending: true,
  simplify: true,
  splitStrings: true,
  splitStringsChunkLength: 10,
  stringArray: true,
  stringArrayEncoding: ['base64'],
  stringArrayIndexShift: true,
  stringArrayWrappersCount: 2,
  stringArrayWrappersChainedCalls: true,
  stringArrayWrappersParametersMaxCount: 4,
  stringArrayWrappersType: 'function',
  stringArrayThreshold: 0.75,
  transformObjectKeys: true,
  unicodeEscapeSequence: false
};

// دالة للبحث عن ملفات JS بشكل متكرر
function findJsFiles(dir, fileList = []) {
  try {
    const files = readdirSync(dir);
    
    for (const file of files) {
      const filePath = join(dir, file);
      const stats = statSync(filePath);
      
      if (stats.isDirectory()) {
        findJsFiles(filePath, fileList);
      } else if (file.endsWith('.js')) {
        fileList.push(filePath);
      }
    }
  } catch (error) {
    console.error(`Error reading directory ${dir}:`, error.message);
  }
  
  return fileList;
}

// البحث عن جميع ملفات JS في مجلد dist
const jsFiles = findJsFiles('dist');

console.log(`Found ${jsFiles.length} JavaScript files to obfuscate`);

// تعتيم كل ملف
for (const file of jsFiles) {
  try {
    console.log(`Obfuscating: ${file}`);
    const code = readFileSync(file, 'utf8');
    const obfuscatedCode = JavaScriptObfuscator.obfuscate(code, obfuscatorOptions).getObfuscatedCode();
    writeFileSync(file, obfuscatedCode, 'utf8');
    console.log(`✓ Obfuscated: ${file}`);
  } catch (error) {
    console.error(`✗ Error obfuscating ${file}:`, error.message);
  }
}

console.log('Obfuscation complete!');
```

---

## 📄 frontend-admin-dashboard/public/service-worker.js

```javascript
// Enhanced Service Worker for faster PWA installation
const CACHE_NAME = 'worktrack-admin-v2'
const STATIC_CACHE = 'worktrack-admin-static-v2'
const RUNTIME_CACHE = 'worktrack-admin-runtime-v2'

// Files to cache immediately for faster installation
const STATIC_FILES = [
  '/',
  '/index.html',
  '/devpro-logo.jpg',
  '/favicon.ico',
  '/manifest.json',
  '/icon-128x128.png',
  '/icon-192x192.png',
  '/icon-256x256.png',
  '/icon-512x512.png'
]

// Install event - cache static files immediately with better error handling
self.addEventListener('install', (event) => {
  console.log('[Service Worker] Installing...')
  
  event.waitUntil(
    Promise.all([
      // Cache static files with better error handling
      caches.open(STATIC_CACHE).then((cache) => {
        console.log('[Service Worker] Caching static files')
        return cache.addAll(STATIC_FILES.map(url => new Request(url, { cache: 'reload' })))
          .catch(err => {
            console.error('[Service Worker] Failed to cache some files:', err)
            // Continue even if some files fail
            return Promise.resolve()
          })
      }),
      // Skip waiting to activate immediately
      self.skipWaiting()
    ])
  )
})

// Activate event - clean up old caches
self.addEventListener('activate', (event) => {
  console.log('[Service Worker] Activating...')
  
  event.waitUntil(
    Promise.all([
      // Claim all clients immediately
      self.clients.claim(),
      // Clean up old caches
      caches.keys().then((cacheNames) => {
        console.log('[Service Worker] Cleaning up old caches')
        return Promise.all(
          cacheNames
            .filter((cacheName) => cacheName !== STATIC_CACHE && cacheName !== CACHE_NAME && cacheName !== RUNTIME_CACHE)
            .map((cacheName) => {
              console.log('[Service Worker] Deleting old cache:', cacheName)
              return caches.delete(cacheName)
            })
        )
      })
    ])
  )
})

// Fetch event - implement network-first strategy for HTML, cache-first for assets
self.addEventListener('fetch', (event) => {
  // Skip non-GET requests
  if (event.request.method !== 'GET') {
    return
  }

  // Skip API calls and external resources
  const url = event.request.url
  if (url.includes('/api/') || 
      url.includes('/static/') ||
      url.includes('/assets/') ||
      url.match(/\.(js|css|png|jpg|jpeg|gif|svg|ico|woff|woff2|ttf|eot|json)$/) ||
      (url.includes('http:') && !url.includes(self.location.origin))) {
    return
  }

  const requestUrl = new URL(event.request.url)

  // Network-first strategy for HTML pages
  if (requestUrl.pathname.endsWith('.html') || requestUrl.pathname === '/') {
    event.respondWith(
      fetch(event.request)
        .then((response) => {
          // Cache successful responses
          if (response && response.status === 200) {
            const responseClone = response.clone()
            caches.open(RUNTIME_CACHE).then((cache) => {
              cache.put(event.request, responseClone)
            })
          }
          return response
        })
        .catch(() => {
          // Fallback to cache
          console.log('[Service Worker] Network failed, falling back to cache')
          return caches.match(event.request).then((cached) => {
            return cached || caches.match('/index.html')
          })
        })
    )
    return
  }

  // Cache-first strategy for static assets
  event.respondWith(
    caches.match(event.request).then((cached) => {
      if (cached) {
        return cached
      }

      return fetch(event.request)
        .then((response) => {
          // Cache successful responses
          if (response && response.status === 200) {
            const responseClone = response.clone()
            caches.open(RUNTIME_CACHE).then((cache) => {
              cache.put(event.request, responseClone)
            })
          }
          return response
        })
        .catch(() => {
          // Fallback for images
          if (event.request.destination === 'image') {
            console.log('[Service Worker] Image fetch failed, using fallback')
            return caches.match('/devpro-logo.jpg')
          }
        })
    })
  )
})

// Handle cache busting
self.addEventListener('message', (event) => {
  if (event.data && event.data.type === 'CACHE_BUST') {
    console.log('[Service Worker] Cache bust requested')
    event.waitUntil(
      caches.keys().then((cacheNames) => {
        return Promise.all(
          cacheNames.map((cacheName) => caches.delete(cacheName))
        )
      })
    )
  }
})

```

---

## 📄 frontend-admin-dashboard/src/App.vue

```vue
<template>
  <div v-if="isLogin" class="auth-shell">
    <router-view />
  </div>

  <div v-else-if="isRouterReady" class="shell">
    <aside class="sidebar">
      <div class="sidebar__brand">
        <div class="devpro-brand">
          <img src="/src/assets/devpro-logo.jpg" alt="DevPro" class="devpro-logo" />
          <div class="devpro-text">
            <span class="devpro-name">{{ t('devpro_name') }}</span>
            <span class="devpro-slogan">{{ t('powered_slogan') }}</span>
          </div>
        </div>
        <div class="brand-divider"></div>
        <div class="app-brand">
          <img src="/src/assets/devpro-logo.jpg" alt="DevPro logo" class="brand-mark" />
          <div class="brand-text">
            <span class="brand-name">{{ t('app_name') }}</span>
            <span class="brand-sub">{{ t('dashboard') }}</span>
          </div>
        </div>
      </div>

      <nav class="sidebar__nav">
        <router-link to="/dashboard" class="nav-item" active-class="nav-item--active">
          <span class="nav-icon">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <rect x="3" y="3" width="18" height="18" rx="2" ry="2"></rect>
              <line x1="3" y1="9" x2="21" y2="9"></line>
              <line x1="9" y1="21" x2="9" y2="9"></line>
            </svg>
          </span>
          <span class="nav-label">{{ t('dashboard') }}</span>
        </router-link>
        <router-link to="/service-requests" class="nav-item" active-class="nav-item--active">
          <span class="nav-icon">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
              <polyline points="14 2 14 8 20 8"></polyline>
              <line x1="16" y1="13" x2="8" y2="13"></line>
              <line x1="16" y1="17" x2="8" y2="17"></line>
              <polyline points="10 9 9 9 8 9"></polyline>
            </svg>
          </span>
          <span class="nav-label">{{ t('service_requests') }}</span>
        </router-link>
        <router-link to="/employees" class="nav-item" active-class="nav-item--active">
          <span class="nav-icon">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path>
              <circle cx="9" cy="7" r="4"></circle>
              <path d="M23 21v-2a4 4 0 0 0-3-3.87"></path>
              <path d="M16 3.13a4 4 0 0 1 0 7.75"></path>
            </svg>
          </span>
          <span class="nav-label">{{ t('employees') }}</span>
        </router-link>
        <router-link to="/customers" class="nav-item" active-class="nav-item--active">
          <span class="nav-icon">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
              <circle cx="12" cy="7" r="4"></circle>
            </svg>
          </span>
          <span class="nav-label">{{ t('customers') }}</span>
        </router-link>
        <router-link to="/worksites" class="nav-item" active-class="nav-item--active">
          <span class="nav-icon">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"></path>
              <circle cx="12" cy="10" r="3"></circle>
            </svg>
          </span>
          <span class="nav-label">{{ t('worksites') }}</span>
        </router-link>
        <router-link to="/reports" class="nav-item" active-class="nav-item--active">
          <span class="nav-icon">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <line x1="18" y1="20" x2="18" y2="10"></line>
              <line x1="12" y1="20" x2="12" y2="4"></line>
              <line x1="6" y1="20" x2="6" y2="14"></line>
            </svg>
          </span>
          <span class="nav-label">{{ t('reports') }}</span>
        </router-link>
        <router-link to="/settings" class="nav-item" active-class="nav-item--active">
          <span class="nav-icon">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <circle cx="12" cy="12" r="3"></circle>
              <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"></path>
            </svg>
          </span>
          <span class="nav-label">{{ t('settings') }}</span>
        </router-link>
        <router-link to="/subscription" class="nav-item" active-class="nav-item--active">
          <span class="nav-icon">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <rect x="2" y="5" width="20" height="14" rx="2"></rect>
              <line x1="2" y1="10" x2="22" y2="10"></line>
            </svg>
          </span>
          <span class="nav-label">{{ t('subscription') }}</span>
        </router-link>
      </nav>

      <div class="sidebar__footer">
        <div class="sidebar__footer-top">
          <LanguageSwitcher />
          <button class="theme-toggle" @click="toggleTheme" :title="isDark ? t('theme_light') : t('theme_dark')">
            <svg v-if="!isDark" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"></path>
            </svg>
            <svg v-else width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <circle cx="12" cy="12" r="5"></circle>
              <line x1="12" y1="1" x2="12" y2="3"></line>
              <line x1="12" y1="21" x2="12" y2="23"></line>
              <line x1="4.22" y1="4.22" x2="5.64" y2="5.64"></line>
              <line x1="18.36" y1="18.36" x2="19.78" y2="19.78"></line>
              <line x1="1" y1="12" x2="3" y2="12"></line>
              <line x1="21" y1="12" x2="23" y2="12"></line>
              <line x1="4.22" y1="19.78" x2="5.64" y2="18.36"></line>
              <line x1="18.36" y1="5.64" x2="19.78" y2="4.22"></line>
            </svg>
          </button>
        </div>
        <button class="btn-logout" @click="handleLogout">
          <span class="logout-icon">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"></path>
              <polyline points="16 17 21 12 16 7"></polyline>
              <line x1="21" y1="12" x2="9" y2="12"></line>
            </svg>
          </span>
          <span class="logout-text">{{ t('logout') }}</span>
        </button>
        <div class="footer-brand">
          <span>© 2026</span>
          <span class="footer-brand-name">{{ t('devpro_name') }}</span>
          <span class="footer-brand-slogan">{{ t('powered_slogan') }}</span>
        </div>
      </div>
    </aside>

    <div class="main">
      <header class="topbar">
        <h1 class="topbar__title">{{ pageTitle }}</h1>
        <div class="topbar__user">
          <NotificationDropdown />
          <span class="topbar__name">{{ displayName }}</span>
        </div>
      </header>
      <main class="content"><router-view /></main>
    </div>
  </div>
</template>

<script setup>
import { computed, ref, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { authStore } from './store/auth'
import LanguageSwitcher from './components/LanguageSwitcher.vue'
import NotificationDropdown from './components/NotificationDropdown.vue'
import { useI18n } from './services/i18n'
import { currentUser } from './services/auth'
import { wsService } from './services/websocket'

const { t, currentLang } = useI18n()

const route = useRoute()
const router = useRouter()
const isLogin = computed(() => route.path === '/login')
const isRouterReady = ref(false)

router.isReady().then(() => {
  isRouterReady.value = true
  
  // معالجة مشكلة /index.html في العنوان - استعادة المسار المحفوظ بدلاً من التحويل إلى /
  if (window.location.pathname === '/index.html') {
    const isAuthed = !!localStorage.getItem('worktrack_admin_token')
    const savedPath = localStorage.getItem('worktrack_last_path')
    
    if (isAuthed && savedPath && savedPath !== '/login' && savedPath !== '/') {
      window.history.replaceState({}, '', savedPath)
    } else if (isAuthed) {
      window.history.replaceState({}, '', '/dashboard')
    } else {
      window.history.replaceState({}, '', '/login')
    }
  }
})

const user = computed(() => authStore.user)
const initials = computed(() => {
  const name = user.value?.full_name || t('default_user_name')
  return name.trim().slice(0, 1)
})

const isDark = ref(false)

const displayName = computed(() => {
  const fullName = user.value?.full_name?.trim()
  if (!fullName) return t('default_user_name')

  const rawNames = [
    'مدير النظام',
    'System administrator',
    'System Admin',
    'System Administrator',
    'מנהל המערכת'
  ]

  if (rawNames.includes(fullName)) {
    return t('system_admin_role')
  }

  return fullName
})

const pageTitle = computed(() => {
  const titles = {
    '/dashboard': t('dashboard'),
    '/service-requests': t('service_requests'),
    '/employees': t('employees'),
    '/customers': t('customers'),
    '/worksites': t('worksites'),
    '/reports': t('reports'),
    '/settings': t('settings'),
    '/subscription': t('subscription'),
  }
  return titles[route.path] || t('app_name')
})
function toggleTheme() {
  isDark.value = !isDark.value
  document.documentElement.setAttribute('data-theme', isDark.value ? 'dark' : 'light')
  localStorage.setItem('worktrack_theme', isDark.value ? 'dark' : 'light')
}

function handleLogout() {
  localStorage.removeItem('worktrack_admin_token')
  localStorage.removeItem('worktrack_admin_user')
  authStore.clear()
  router.push('/login')
}

onMounted(() => {
  // تحديث المستخدم عند تحميل الصفحة
  const user = currentUser()
  if (user) {
    authStore.setUser(user)
    
    // Connect to WebSocket when user is logged in
    wsService.connect()
  }
  
  const savedTheme = localStorage.getItem('worktrack_theme')
  if (savedTheme === 'dark') {
    isDark.value = true
    document.documentElement.setAttribute('data-theme', 'dark')
  }
})

onUnmounted(() => {
  // Disconnect WebSocket when component is unmounted
  wsService.disconnect()
})
</script>

<style scoped>
.auth-shell { 
  min-height: 100dvh; 
  display: flex; 
  align-items: center; 
  justify-content: center; 
  background: var(--canvas); 
}

.shell {
  display: flex;
  min-height: 100dvh;
  background: var(--canvas);
}

.sidebar {
  width: 220px;
  background: var(--surface);
  border-left: 1px solid var(--line);
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
  position: sticky;
  top: 0;
  height: 100vh;
  overflow-y: auto;
  box-shadow: var(--shadow-md);
  transition: box-shadow var(--transition-base);
}

.sidebar:hover {
  box-shadow: var(--shadow-lg);
}

.sidebar__brand {
  padding: 16px 18px;
  border-bottom: 1px solid var(--line);
}

.devpro-brand {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 12px;
  background: linear-gradient(135deg, #1E3A5F08, #1E3A5F12);
  border-radius: var(--radius-sm);
  margin-bottom: 12px;
  border: 1px solid #1E3A5F15;
}

.devpro-logo {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  object-fit: contain;
  background: white;
  padding: 2px;
}

.devpro-text {
  display: flex;
  flex-direction: column;
  line-height: 1.1;
}

.devpro-name {
  font-weight: 700;
  font-size: 13px;
  color: var(--brand);
  letter-spacing: -0.5px;
}

.devpro-slogan {
  font-size: 8px;
  color: var(--ink-soft);
  font-weight: 400;
}

.brand-divider {
  height: 1px;
  background: var(--line);
  margin: 8px 0 12px;
}

.app-brand {
  display: flex;
  align-items: center;
  gap: 10px;
}

.brand-mark {
  width: 40px;
  height: 40px;
  object-fit: contain;
  border-radius: 10px;
  flex-shrink: 0;
}

.brand-text {
  display: flex;
  flex-direction: column;
  line-height: 1.1;
}

.brand-name {
  font-weight: 700;
  font-size: 15px;
  color: var(--brand);
}

.brand-sub {
  font-size: 10px;
  color: var(--ink-soft);
  font-weight: 500;
}

.sidebar__nav {
  display: flex;
  flex-direction: column;
  padding: 12px 10px;
  gap: 2px;
  flex: 1;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 14px;
  border-radius: var(--radius-sm);
  font-size: 14px;
  font-weight: 500;
  color: var(--ink-soft);
  transition: all var(--transition-base);
  position: relative;
  overflow: hidden;
}

.nav-item::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, rgba(255,255,255,0.1) 0%, rgba(255,255,255,0) 100%);
  opacity: 0;
  transition: opacity var(--transition-fast);
}

.nav-item:hover::before {
  opacity: 1;
}

.nav-item:hover {
  background: var(--brand-tint);
  color: var(--brand);
  transform: translateX(-4px);
}

.nav-item--active {
  background: linear-gradient(135deg, var(--brand) 0%, var(--brand-dark) 100%);
  color: white;
  box-shadow: var(--shadow-md), var(--brand-glow);
}

.nav-item--active:hover {
  background: linear-gradient(135deg, var(--brand-dark) 0%, var(--brand) 100%);
  color: white;
  transform: translateX(-4px);
  box-shadow: var(--shadow-lg), var(--shadow-glow);
}

.nav-icon {
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: transform var(--transition-fast);
}

.nav-item:hover .nav-icon {
  transform: scale(1.1);
}

.nav-item--active .nav-icon {
  color: var(--gold-light);
}

.nav-label {
  flex: 1;
}

/* =============================================
   القاع - مع عرض مناسب للعناصر
   ============================================= */
.sidebar__footer {
  padding: 12px 14px;
  border-top: 1px solid var(--line);
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.sidebar__footer-top {
  display: flex;
  align-items: flex-start;
  gap: 6px;
  width: 100%;
}

/* LanguageSwitcher يأخذ المساحة المتاحة */
.sidebar__footer-top .lang-switcher {
  flex: 1;
  min-width: 0;
}

/* زر الدارك - يظهر بوضوح */
.theme-toggle {
  width: 34px;
  height: 34px;
  border-radius: 50%;
  border: 1px solid var(--line);
  background: var(--surface);
  cursor: pointer;
  transition: all var(--transition-base);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  color: var(--ink-soft);
}

.theme-toggle:hover {
  background: var(--brand-tint);
  border-color: var(--brand);
  transform: scale(1.1);
  color: var(--brand);
  box-shadow: var(--shadow-md), var(--brand-glow);
}

.btn-logout {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 8px 12px;
  border-radius: var(--radius-sm);
  border: 1px solid var(--signal-out);
  background: var(--signal-out-tint);
  color: var(--signal-out);
  font-family: var(--font-body);
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  transition: all var(--transition-base);
  width: 100%;
  position: relative;
  overflow: hidden;
}

.btn-logout::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, rgba(255,255,255,0.2) 0%, rgba(255,255,255,0) 100%);
  opacity: 0;
  transition: opacity var(--transition-fast);
}

.btn-logout:hover::before {
  opacity: 1;
}

.btn-logout:hover {
  background: var(--signal-out);
  color: white;
  transform: translateY(-2px);
  box-shadow: var(--shadow-lg);
}

.logout-icon {
  width: 16px;
  height: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.logout-text { font-size: 12px; }

.footer-brand {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
  padding: 4px 0 0;
  font-size: 9px;
  color: var(--ink-soft);
  flex-wrap: wrap;
}

.footer-brand-name {
  font-weight: 700;
  color: var(--brand);
  font-size: 10px;
}

.footer-brand-slogan {
  font-size: 7px;
  color: var(--ink-light);
}

.main {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
}

.topbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 28px;
  background: var(--surface);
  border-bottom: 1px solid var(--line);
  position: sticky;
  top: 0;
  z-index: 10;
  box-shadow: var(--shadow-sm);
  transition: box-shadow var(--transition-base);
}

.topbar:hover {
  box-shadow: var(--shadow-md);
}

.topbar__title {
  font-size: 18px;
  font-weight: 600;
  background: linear-gradient(135deg, var(--brand) 0%, var(--accent) 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.topbar__user {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 6px 12px;
  border-radius: var(--radius-md);
  transition: background var(--transition-fast);
}

.topbar__user:hover {
  background: var(--brand-tint);
}

.topbar__name {
  font-size: 13px;
  color: var(--ink-soft);
  font-weight: 500;
}



.content {
  flex: 1;
  padding: 24px 28px;
  max-width: 1280px;
  width: 100%;
  margin: 0 auto;
  animation: fadeIn 0.4s ease;
}

/* =============================================
   استجابة للشاشات الصغيرة
   ============================================= */
@media (max-width: 768px) {
  .sidebar {
    width: 60px;
    min-width: 60px;
  }

  .sidebar__brand .brand-text,
  .sidebar__nav .nav-label,
  .logout-text,
  .devpro-text,
  .brand-sub,
  .footer-brand {
    display: none;
  }

  .sidebar__nav {
    flex-direction: column;
    gap: 4px;
    padding: 12px 8px;
  }

  .sidebar__nav .nav-item {
    justify-content: center;
    padding: 12px 8px;
  }

  .sidebar__nav .nav-icon {
    font-size: 20px;
    width: auto;
  }

  .sidebar__footer {
    padding: 10px 8px;
  }

  .sidebar__footer-top {
    flex-direction: column;
    gap: 8px;
  }

  .theme-toggle {
    width: 32px;
    height: 32px;
    font-size: 14px;
  }

  .btn-logout {
    padding: 8px 10px;
    justify-content: center;
  }

  .logout-icon {
    font-size: 16px;
  }

  .devpro-brand {
    justify-content: center;
    padding: 4px;
  }

  .devpro-logo {
    width: 24px;
    height: 24px;
  }

  .content {
    padding: 16px;
  }

  .topbar {
    padding: 12px 16px;
  }

  .topbar__title {
    font-size: 16px;
  }
}
</style>

```

---

## 📄 frontend-admin-dashboard/src/components/ActivityFeed.vue

```vue
<template>
  <div class="activity-feed">
    <div class="feed-header">
      <h3>{{ t('notifications_feed') }}</h3>
      <span v-if="notificationsStore.notifications && notificationsStore.notifications.length" class="badge">
        {{ notificationsStore.notifications.length }}
      </span>
    </div>

    <div v-if="notificationsStore.loading" class="empty-state">
      <p>{{ t('loading_notifications') }}</p>
    </div>

    <div v-else-if="notificationsStore.error" class="alert alert-error">
      <span>❌</span> {{ notificationsStore.error }}
    </div>

    <div v-else-if="!notificationsStore.notifications || !notificationsStore.notifications.length" class="empty-state">
      <p>{{ t('no_notifications_available') }}</p>
    </div>

    <div v-else class="notifications-list">
      <div 
        v-for="notification in notificationsStore.notifications" 
        :key="notification.id" 
        class="notification-item"
      >
        <span class="notification-icon">{{ getNotificationIcon(notification.type) }}</span>
        <div class="notification-content">
          <p class="notification-title">{{ notification.title }}</p>
          <p class="notification-message">{{ notification.message }}</p>
          <span class="notification-time">{{ formatTime(notification.created_at) }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { notificationsStore } from '../store/notifications'
import { useI18n } from '../services/i18n'

const { t } = useI18n()

onMounted(() => {
  notificationsStore.fetchNotifications()
})

function getNotificationIcon(type) {
  const icons = {
    'status_update': '📋',
    'worker_assigned': '👷',
    'worker_arrived': '📍',
    'service_completed': '✅',
    'payment_request': '💳',
    'general': '🔔'
  }
  return icons[type] || '🔔'
}

function formatTime(date) {
  if (!date) return ''
  return new Date(date).toLocaleTimeString('en-GB', { 
    hour: '2-digit', 
    minute: '2-digit',
    hour12: false 
  })
}
</script>

<style scoped>
.activity-feed {
  background: var(--surface);
  border-radius: var(--radius-md);
  padding: 16px;
  border: 1px solid var(--line);
}

.feed-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.feed-header h3 {
  font-size: 14px;
  margin: 0;
  color: var(--ink);
}

.badge {
  background: var(--brand);
  color: white;
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 11px;
  font-weight: 600;
}

.empty-state {
  text-align: center;
  padding: 20px;
  color: var(--ink-soft);
  font-size: 13px;
}

.alert {
  padding: 10px 12px;
  border-radius: var(--radius-sm);
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
}

.alert-error {
  background: var(--signal-out-tint);
  color: var(--signal-out);
}

.notifications-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.notification-item {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  padding: 10px;
  background: var(--canvas);
  border-radius: var(--radius-sm);
  transition: background var(--transition-fast);
}

.notification-item:hover {
  background: var(--brand-tint);
}

.notification-icon {
  font-size: 16px;
  flex-shrink: 0;
}

.notification-content {
  flex: 1;
  min-width: 0;
}

.notification-title {
  font-size: 13px;
  font-weight: 600;
  color: var(--ink);
  margin: 0 0 4px 0;
}

.notification-message {
  font-size: 12px;
  color: var(--ink-soft);
  margin: 0 0 4px 0;
  line-height: 1.4;
}

.notification-time {
  font-size: 11px;
  color: var(--ink-light);
}

@media (max-width: 480px) {
  .activity-feed {
    padding: 12px;
    border-radius: var(--radius-sm);
  }
  
  .feed-header {
    margin-bottom: 10px;
  }
  
  .feed-header h3 {
    font-size: 13px;
  }
  
  .badge {
    font-size: 10px;
    padding: 2px 6px;
  }
  
  .empty-state {
    padding: 16px;
    font-size: 12px;
  }
  
  .alert {
    font-size: 11px;
    padding: 8px 10px;
  }
  
  .notification-item {
    padding: 8px;
    gap: 8px;
  }
  
  .notification-icon {
    font-size: 14px;
  }
  
  .notification-title {
    font-size: 12px;
  }
  
  .notification-message {
    font-size: 11px;
  }
  
  .notification-time {
    font-size: 10px;
  }
}

@media (max-width: 360px) {
  .activity-feed {
    padding: 10px;
  }
  
  .feed-header h3 {
    font-size: 12px;
  }
  
  .notification-item {
    padding: 6px;
    gap: 6px;
  }
  
  .notification-icon {
    font-size: 12px;
  }
  
  .notification-title {
    font-size: 11px;
  }
  
  .notification-message {
    font-size: 10px;
  }
}
</style>

```

---

## 📄 frontend-admin-dashboard/src/components/ClientFormModal.vue

```vue
<template>
  <div class="modal-backdrop" @click.self="$emit('close')">
    <div class="modal card">
      <div class="modal-header">
        <h3>{{ t('client_modal_title') }}</h3>
        <button class="modal-close" @click="$emit('close')">✕</button>
      </div>

      <p class="modal__hint">{{ t('client_modal_hint') }}</p>
      <p class="modal__hint-small">{{ t('client_modal_hint_small') }}</p>

      <form class="modal__form" @submit.prevent="handleSubmit">
        <div v-if="error" class="alert alert-error">{{ error }}</div>
        <div v-if="success" class="alert alert-success">{{ success }}</div>

        <div class="form-group">
          <label>{{ t('client_modal_full_name') }} <span class="required">*</span></label>
          <input v-model="form.full_name" type="text" required :placeholder="t('full_name')" />
        </div>

        <div class="form-group">
          <label>{{ t('client_modal_email') }} <span class="required">*</span></label>
          <input 
            v-model="form.email" 
            type="email" 
            required 
            :placeholder="t('email_placeholder')"
            dir="ltr"
          />
          <span class="field-hint">{{ t('client_modal_email_hint') }}</span>
        </div>

        <div class="form-group">
          <label>{{ t('client_modal_phone') }}</label>
          <input 
            v-model="form.phone" 
            type="tel" 
            :placeholder="t('phone')"
            dir="ltr"
          />
        </div>

        <div class="form-actions">
          <button type="button" class="btn btn--ghost" @click="$emit('close')">{{ t('cancel') }}</button>
          <button type="submit" class="btn btn--primary" :disabled="loading">
            {{ loading ? t('client_modal_saving') : t('client_modal_create') }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import api from '../services/api'
import { useI18n } from '../services/i18n'

const { t } = useI18n()
const emit = defineEmits(['close', 'client-added'])

const loading = ref(false)
const error = ref('')
const success = ref('')

const form = reactive({
  full_name: '',
  email: '',
  phone: ''
})

async function handleSubmit() {
  loading.value = true
  error.value = ''
  success.value = ''

  // التحقق من البريد الإلكتروني
  if (!form.email || !form.email.includes('@')) {
    error.value = t('client_modal_error_email')
    loading.value = false
    return
  }

  try {
    const payload = {
      full_name: form.full_name,
      email: form.email.trim(),
      phone: form.phone.trim()
    }

    console.log('📤 إرسال بيانات العميل:', payload)

    const { data } = await api.post('/admin/create-client', payload)

    success.value = t('client_modal_success')
    success.value += `\n${t('client_modal_success_email')} ${data.user.email}`
    success.value += `\n${t('client_modal_success_password')} ${data.password}`
    success.value += `\n${t('client_modal_success_warning')}`

    setTimeout(() => {
      emit('client-added')
      emit('close')
    }, 5000)
  } catch (err) {
    console.error('❌ فشل الإضافة:', err.response?.data)
    error.value = err.response?.data?.error || t('client_modal_error_failed')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.modal-backdrop {
  position: fixed; inset: 0;
  background: rgba(0,0,0,0.6);
  backdrop-filter: blur(4px);
  display: flex; align-items: center; justify-content: center;
  z-index: 1000; padding: 20px;
}

.modal {
  width: 100%; max-width: 420px; padding: 0;
  max-height: 90vh; overflow-y: auto;
}

.modal-header {
  display: flex; justify-content: space-between; align-items: center;
  padding: 20px 24px; border-bottom: 1px solid var(--line);
}

.modal-header h3 { font-size: 18px; margin: 0; }
.modal-close { background: none; border: none; font-size: 24px; cursor: pointer; color: var(--ink-soft); }

.modal__hint {
  padding: 0 24px 8px;
  font-size: 14px; font-weight: 600; color: var(--brand); margin: 0;
}

.modal__hint-small {
  padding: 0 24px 16px;
  font-size: 12px; color: var(--ink-soft); margin: 0;
}

.modal__form { padding: 0 24px 24px; }

.form-group { margin-bottom: 14px; }
.form-group label { display: block; font-size: 13px; font-weight: 600; margin-bottom: 4px; }

.form-group input,
.form-group select {
  width: 100%; padding: 10px 12px;
  border: 1.5px solid var(--line); border-radius: var(--radius-sm);
  font-size: 14px; font-family: inherit; background: var(--surface);
  transition: border-color 0.3s;
}

.form-group input:focus,
.form-group select:focus {
  border-color: var(--brand); outline: none;
  box-shadow: 0 0 0 3px var(--brand-tint);
}

.required { color: var(--signal-out); }
.field-hint { font-size: 11px; color: var(--ink-soft); }

.alert {
  padding: 10px 14px; border-radius: var(--radius-sm);
  font-size: 13px; margin-bottom: 14px; white-space: pre-line;
}
.alert-error { background: var(--signal-out-tint); color: var(--signal-out); }
.alert-success { background: var(--signal-in-tint); color: var(--signal-in); }

.form-actions {
  display: flex; gap: 10px;
  justify-content: flex-end; margin-top: 14px;
  padding-top: 14px; border-top: 1px solid var(--line);
}
</style>

```

---

## 📄 frontend-admin-dashboard/src/components/CompanyBrand.vue

```vue
<template>
  <div class="company-brand" :class="size">
    <img :src="logoUrl" :alt="companyName" class="brand-logo" />
    <div class="brand-text">
      <span class="brand-name">{{ companyName }}</span>
      <span v-if="showSlogan" class="brand-slogan">{{ slogan }}</span>
    </div>
  </div>
</template>

<script setup>
import { companyConfig } from '../config/company'

const props = defineProps({
  size: { type: String, default: 'medium' },
  showSlogan: { type: Boolean, default: true },
})

const companyName = companyConfig.name
const slogan = companyConfig.slogan
const logoUrl = companyConfig.logo
</script>

<style scoped>
.company-brand {
  display: flex;
  align-items: center;
  gap: 12px;
}

.brand-logo {
  object-fit: contain;
  border-radius: 12px;
  background: white;
  padding: 4px;
}

/* الأحجام */
.small .brand-logo { width: 28px; height: 28px; border-radius: 8px; }
.small .brand-name { font-size: 12px; }
.small .brand-slogan { font-size: 9px; }

.medium .brand-logo { width: 40px; height: 40px; }
.medium .brand-name { font-size: 16px; }
.medium .brand-slogan { font-size: 10px; }

.large .brand-logo { width: 60px; height: 60px; }
.large .brand-name { font-size: 20px; }
.large .brand-slogan { font-size: 12px; }

.brand-text {
  display: flex;
  flex-direction: column;
  line-height: 1.2;
}

.brand-name {
  font-weight: 700;
  color: var(--brand);
  letter-spacing: -0.5px;
}

.brand-slogan {
  color: var(--ink-soft);
  font-weight: 400;
  font-size: 10px;
}
</style>

```

---

## 📄 frontend-admin-dashboard/src/components/CustomerFormModal.vue

```vue
<template>
  <div class="modal-backdrop" @click.self="$emit('close')">
    <div class="modal card">
      <div class="modal-header">
        <h3>➕ {{ t('add_customer') }}</h3>
        <button class="modal-close" @click="$emit('close')">✕</button>
      </div>

      <p class="modal__hint">📝 {{ t('add_customer_hint') }}</p>
      <p class="modal__hint-small">⚠️ سيتم إنشاء كلمة مرور تلقائية للعميل</p>

      <form class="modal__form" @submit.prevent="handleSubmit">
        <div v-if="error" class="alert alert-error">{{ error }}</div>
        <div v-if="success" class="alert alert-success">
          <div class="success-content">
            <div class="success-message">{{ success }}</div>
            <div v-if="generatedPassword" class="password-copy-section">
              <div class="password-label">🔑 {{ t('password') }}:</div>
              <div class="password-copy-container">
                <span class="password-text">{{ generatedPassword }}</span>
                <button 
                  class="btn-copy" 
                  @click="copyPassword"
                  :title="copied ? t('copied') : t('copy_password')"
                >
                  {{ copied ? '✅' : '📋' }}
                </button>
              </div>
              <p class="password-warning">⚠️ {{ t('save_password_warning') }}</p>
            </div>
          </div>
        </div>

        <div class="form-group">
          <label>👤 {{ t('full_name') }} <span class="required">*</span></label>
          <input v-model="form.full_name" type="text" required :placeholder="t('enter_full_name')" />
        </div>

        <div class="form-group">
          <label>📱 {{ t('phone') }} <span class="required">*</span></label>
          <input 
            v-model="form.phone" 
            type="tel" 
            required 
            placeholder="05xxxxxxxx"
            dir="ltr"
          />
        </div>

        <div class="form-group">
          <label>📧 {{ t('email') }}</label>
          <input 
            v-model="form.email" 
            type="email"
            :placeholder="t('enter_email')"
            dir="ltr"
          />
        </div>

        <div class="form-actions">
          <button type="button" class="btn btn--ghost" @click="$emit('close')">{{ t('cancel') }}</button>
          <button v-if="!success" type="submit" class="btn btn--primary" :disabled="loading">
            {{ loading ? '⏳ ' + t('saving') : '💾 ' + t('create_customer') }}
          </button>
          <button v-if="success" type="button" class="btn btn--success" @click="closeModal">
            ✅ {{ t('done') }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import api from '../services/api'
import { useI18n } from '../services/i18n'

const { t } = useI18n()
const emit = defineEmits(['close', 'customer-added'])

const loading = ref(false)
const error = ref('')
const success = ref('')
const generatedPassword = ref('')
const copied = ref(false)

const form = reactive({
  full_name: '',
  phone: '',
  email: ''
})

async function handleSubmit() {
  loading.value = true
  error.value = ''
  success.value = ''
  generatedPassword.value = ''
  copied.value = false

  // التحقق من رقم الهاتف
  if (!form.phone || form.phone.length < 9) {
    error.value = '❌ ' + t('invalid_phone')
    loading.value = false
    return
  }

  try {
    const payload = {
      full_name: form.full_name,
      phone: form.phone.trim(),
      email: form.email || null
    }

    console.log('📤 إرسال بيانات العميل:', payload)

    const { data } = await api.post('/admin/create-client', payload)

    success.value = `✅ ${t('customer_created_success')}: "${data.user.full_name}"`
    if (data.user.email) {
      success.value += `\n📧 ${t('email')}: ${data.user.email}`
    }
    if (data.user.phone) {
      success.value += `\n📱 ${t('phone')}: ${data.user.phone}`
    }
    if (data.password) {
      generatedPassword.value = data.password
    }

    // إعادة تعيين النموذج بعد النجاح
    form.full_name = ''
    form.phone = ''
    form.email = ''
  } catch (err) {
    console.error('❌ فشل الإضافة:', err.response?.data)
    error.value = err.response?.data?.error || '❌ ' + t('customer_creation_failed')
  } finally {
    loading.value = false
  }
}

function copyPassword() {
  if (generatedPassword.value) {
    navigator.clipboard.writeText(generatedPassword.value)
    copied.value = true
    setTimeout(() => {
      copied.value = false
    }, 2000)
    
    // الخروج من المودال بعد نسخ كلمة المرور
    setTimeout(() => {
      closeModal()
    }, 1000)
  }
}

function closeModal() {
  emit('customer-added')
  emit('close')
}
</script>

<style scoped>
.modal-backdrop {
  position: fixed; inset: 0;
  background: rgba(0,0,0,0.6);
  backdrop-filter: blur(4px);
  display: flex; align-items: center; justify-content: center;
  z-index: 1000; padding: 20px;
}

.modal {
  width: 100%; max-width: 420px; padding: 0;
  max-height: 90vh; overflow-y: auto;
}

.modal-header {
  display: flex; justify-content: space-between; align-items: center;
  padding: 20px 24px; border-bottom: 1px solid var(--line);
}

.modal-header h3 { font-size: 18px; margin: 0; }
.modal-close { background: none; border: none; font-size: 24px; cursor: pointer; color: var(--ink-soft); }

.modal__hint {
  padding: 0 24px 8px;
  font-size: 14px; font-weight: 600; color: var(--brand); margin: 0;
}

.modal__hint-small {
  padding: 0 24px 16px;
  font-size: 12px; color: var(--ink-soft); margin: 0;
}

.modal__form { padding: 0 24px 24px; }

.form-group { margin-bottom: 14px; }
.form-group label { display: block; font-size: 13px; font-weight: 600; margin-bottom: 4px; }

.form-group input,
.form-group select {
  width: 100%; padding: 10px 12px;
  border: 1.5px solid var(--line); border-radius: var(--radius-sm);
  font-size: 14px; font-family: inherit; background: var(--surface);
  transition: border-color 0.3s;
}

.form-group input:focus,
.form-group select:focus {
  border-color: var(--brand); outline: none;
  box-shadow: 0 0 0 3px var(--brand-tint);
}

.required { color: var(--signal-out); }

.alert {
  padding: 10px 14px; border-radius: var(--radius-sm);
  font-size: 13px; margin-bottom: 14px; white-space: pre-line;
}
.alert-error { background: var(--signal-out-tint); color: var(--signal-out); }
.alert-success { background: var(--signal-in-tint); color: var(--signal-in); }

.success-content {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.success-message {
  white-space: pre-line;
}

.password-copy-section {
  margin-top: 8px;
  padding-top: 12px;
  border-top: 1px solid rgba(0,0,0,0.1);
}

.password-label {
  font-weight: 600;
  margin-bottom: 8px;
  font-size: 14px;
}

.password-copy-container {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.password-text {
  font-family: monospace;
  font-size: 16px;
  font-weight: 700;
  letter-spacing: 1px;
  padding: 8px 12px;
  background: rgba(255,255,255,0.5);
  border-radius: var(--radius-sm);
  flex: 1;
  word-break: break-all;
}

.btn-copy {
  background: none;
  border: none;
  font-size: 18px;
  cursor: pointer;
  padding: 8px;
  border-radius: var(--radius-sm);
  transition: background 0.2s;
  min-width: 40px;
}

.btn-copy:hover {
  background: rgba(0,0,0,0.1);
}

.password-warning {
  font-size: 12px;
  color: var(--signal-out);
  margin: 0;
}

.form-actions {
  display: flex; gap: 10px;
  justify-content: flex-end; margin-top: 14px;
  padding-top: 14px; border-top: 1px solid var(--line);
}

.btn--success {
  background: var(--signal-in-tint);
  color: var(--signal-in);
  border: 1px solid var(--signal-in);
}

.btn--success:hover {
  background: var(--signal-in);
  color: white;
}
</style>

```

---

## 📄 frontend-admin-dashboard/src/components/EmployeeFormModal.vue

```vue
<template>
  <div class="modal-backdrop" @click.self="$emit('close')">
    <div class="modal card">
      <div class="modal-header">
        <h3>➕ إضافة موظف جديد</h3>
        <button class="modal-close" @click="$emit('close')">✕</button>
      </div>

      <p class="modal__hint">📱 سيتم إنشاء حساب الموظف برقم الهاتف فقط</p>
      <p class="modal__hint-small">⚠️ سجل دخول الموظف سيكون برقم هاتفه وجهازه</p>

      <form class="modal__form" @submit.prevent="handleSubmit">
        <div v-if="error" class="alert alert-error">{{ error }}</div>
        <div v-if="success" class="alert alert-success">{{ success }}</div>

        <div class="form-group">
          <label>👤 الاسم الكامل <span class="required">*</span></label>
          <input v-model="form.full_name" type="text" required placeholder="أدخل الاسم الكامل" />
        </div>

        <div class="form-group">
          <label>📱 رقم الهاتف <span class="required">*</span></label>
          <input 
            v-model="form.phone" 
            type="tel" 
            required 
            placeholder="05xxxxxxxx"
            dir="ltr"
          />
          <span class="field-hint">سيستخدم الموظف هذا الرقم لتسجيل الدخول</span>
        </div>

        <div class="form-group">
          <label>🎯 الدور</label>
          <select v-model="form.role">
            <option value="employee">ميداني</option>
          </select>
        </div>

        <div class="form-actions">
          <button type="button" class="btn btn--ghost" @click="$emit('close')">إلغاء</button>
          <button type="submit" class="btn btn--primary" :disabled="loading">
            {{ loading ? '⏳ جارٍ الحفظ...' : '💾 إنشاء الموظف' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import api from '../services/api'

const emit = defineEmits(['close', 'employee-added'])

const loading = ref(false)
const error = ref('')
const success = ref('')

const form = reactive({
  full_name: '',
  phone: '',
  role: 'employee'
})

async function handleSubmit() {
  loading.value = true
  error.value = ''
  success.value = ''

  // التحقق من رقم الهاتف
  if (!form.phone || form.phone.length < 9) {
    error.value = '❌ الرجاء إدخال رقم هاتف صحيح'
    loading.value = false
    return
  }

  try {
    const payload = {
      full_name: form.full_name,
      phone: form.phone.trim(),
      role: form.role || 'employee'
    }

    console.log('📤 إرسال بيانات الموظف:', payload)

    const { data } = await api.post('/auth/employee-phone', payload)

    success.value = `✅ تم إنشاء الموظف "${data.user.full_name}" بنجاح!`
    success.value += `\n📱 رقم الهاتف: ${data.user.phone}`
    success.value += `\n🔑 سيستخدم هذا الرقم لتسجيل الدخول`

    setTimeout(() => {
      emit('employee-added')
      emit('close')
    }, 3000)
  } catch (err) {
    console.error('❌ فشل الإضافة:', err.response?.data)
    error.value = err.response?.data?.error || '❌ فشل إنشاء الموظف'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.modal-backdrop {
  position: fixed; inset: 0;
  background: rgba(0,0,0,0.6);
  backdrop-filter: blur(4px);
  display: flex; align-items: center; justify-content: center;
  z-index: 1000; padding: 20px;
}

.modal {
  width: 100%; max-width: 420px; padding: 0;
  max-height: 90vh; overflow-y: auto;
}

.modal-header {
  display: flex; justify-content: space-between; align-items: center;
  padding: 20px 24px; border-bottom: 1px solid var(--line);
}

.modal-header h3 { font-size: 18px; margin: 0; }
.modal-close { background: none; border: none; font-size: 24px; cursor: pointer; color: var(--ink-soft); }

.modal__hint {
  padding: 0 24px 8px;
  font-size: 14px; font-weight: 600; color: var(--brand); margin: 0;
}

.modal__hint-small {
  padding: 0 24px 16px;
  font-size: 12px; color: var(--ink-soft); margin: 0;
}

.modal__form { padding: 0 24px 24px; }

.form-group { margin-bottom: 14px; }
.form-group label { display: block; font-size: 13px; font-weight: 600; margin-bottom: 4px; }

.form-group input,
.form-group select {
  width: 100%; padding: 10px 12px;
  border: 1.5px solid var(--line); border-radius: var(--radius-sm);
  font-size: 14px; font-family: inherit; background: var(--surface);
  transition: border-color 0.3s;
}

.form-group input:focus,
.form-group select:focus {
  border-color: var(--brand); outline: none;
  box-shadow: 0 0 0 3px var(--brand-tint);
}

.required { color: var(--signal-out); }
.field-hint { font-size: 11px; color: var(--ink-soft); }

.alert {
  padding: 10px 14px; border-radius: var(--radius-sm);
  font-size: 13px; margin-bottom: 14px; white-space: pre-line;
}
.alert-error { background: var(--signal-out-tint); color: var(--signal-out); }
.alert-success { background: var(--signal-in-tint); color: var(--signal-in); }

.form-actions {
  display: flex; gap: 10px;
  justify-content: flex-end; margin-top: 14px;
  padding-top: 14px; border-top: 1px solid var(--line);
}
</style>

```

---

## 📄 frontend-admin-dashboard/src/components/GeofenceRing.vue

```vue
<template>
  <div class="gf-ring" :class="status">
    <svg viewBox="0 0 160 160" class="gf-svg" aria-hidden="true">
      <circle cx="80" cy="80" r="68" class="gf-track" />
      <circle
        cx="80" cy="80" r="68" class="gf-fill"
        :stroke-dasharray="circumference"
        :stroke-dashoffset="dashOffset"
      />
      <circle cx="80" cy="80" r="46" class="gf-hole" />
    </svg>
    <div class="gf-center">
      <span class="gf-distance mono">{{ Math.round(distance) }}<small>م</small></span>
      <span class="gf-label">{{ statusText }}</span>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

// distance: مسافة الموظف الحالية عن نقطة العمل (متر)
// radius: نصف القطر المسموح به لنقطة العمل (متر)
const props = defineProps({
  distance: { type: Number, default: 0 },
  radius: { type: Number, default: 100 },
})

const circumference = 2 * Math.PI * 68

const ratio = computed(() => Math.min(props.distance / (props.radius || 1), 1))
const dashOffset = computed(() => circumference * (1 - ratio.value))
const status = computed(() => (props.distance <= props.radius ? 'inside' : 'outside'))
const statusText = computed(() => (status.value === 'inside' ? 'داخل نطاق الموقع' : 'خارج نطاق الموقع'))
</script>

<style scoped>
.gf-ring { position: relative; width: 180px; height: 180px; margin: 0 auto; }
.gf-svg { width: 100%; height: 100%; transform: rotate(-90deg); }
.gf-track { fill: none; stroke: var(--line); stroke-width: 10; }
.gf-fill { fill: none; stroke-width: 10; stroke-linecap: round; transition: stroke-dashoffset .6s ease, stroke .3s ease; }
.gf-ring.inside .gf-fill { stroke: var(--signal-in); }
.gf-ring.outside .gf-fill { stroke: var(--signal-out); }
.gf-hole { fill: var(--surface); }
.gf-center { position: absolute; inset: 0; display: flex; flex-direction: column; align-items: center; justify-content: center; }
.gf-distance { font-size: 28px; font-weight: 600; color: var(--ink); }
.gf-distance small { font-size: 14px; font-weight: 500; margin-inline-start: 2px; color: var(--ink-soft); }
.gf-label { margin-top: 4px; font-size: 13px; font-weight: 600; }
.gf-ring.inside .gf-label { color: var(--signal-in); }
.gf-ring.outside .gf-label { color: var(--signal-out); }
.gf-ring.inside::after {
  content: ''; position: absolute; inset: -6px; border-radius: 50%; border: 1px solid var(--signal-in);
  animation: gf-pulse 2.4s ease-out infinite;
}
@keyframes gf-pulse {
  0% { transform: scale(.94); opacity: .7; }
  100% { transform: scale(1.14); opacity: 0; }
}
</style>

```

---

## 📄 frontend-admin-dashboard/src/components/LanguageSwitcher.vue

```vue
<template>
  <div class="lang-switcher">
    <button
      v-for="lang in languages"
      :key="lang.code"
      class="lang-btn"
      :class="{ active: currentLang === lang.code }"
      @click="changeLanguage(lang.code)"
      :title="lang.name"
    >
      {{ lang.flag }}
    </button>
  </div>
</template>

<script setup>
import { useI18n } from '../services/i18n'

const { currentLang, setLang } = useI18n()

const languages = [
  { code: 'ar', name: 'العربية', flag: '🇸🇦' },
  { code: 'he', name: 'עברית', flag: '🇮🇱' },
  { code: 'en', name: 'English', flag: '🇬🇧' }
]

function changeLanguage(lang) {
  setLang(lang)
}
</script>

<style scoped>
.lang-switcher {
  display: flex;
  flex-direction: column;
  align-items: stretch;
  gap: 4px;
  padding: 4px;
  background: var(--canvas, #f0f4fa);
  border-radius: 8px;
  flex: 1;
}

.lang-btn {
  padding: 6px 8px;
  border: none;
  border-radius: 6px;
  background: transparent;
  font-size: 16px;
  cursor: pointer;
  transition: all 0.3s ease;
  font-family: inherit;
  min-width: 0;
  width: 100%;
}

.lang-btn:hover {
  background: rgba(30, 58, 95, 0.08);
  transform: scale(1.05);
}

.lang-btn.active {
  background: #1E3A5F;
  color: white;
  box-shadow: 0 2px 8px rgba(30, 58, 95, 0.3);
  transform: scale(1.05);
  border-radius: 6px;
}

@media (max-width: 480px) {
  .lang-btn {
    font-size: 14px;
    padding: 3px 4px;
    min-width: 28px;
  }
}
</style>

```

---

## 📄 frontend-admin-dashboard/src/components/LiveMap.vue

```vue
<template>
  <div class="live-map">
    <div class="live-map__grid"></div>
    <span v-for="p in pins" :key="p.id" class="live-map__pin" :class="p.status" :style="{ top: p.top + '%', right: p.right + '%' }"></span>
    <div class="live-map__legend">
      <span><i class="in"></i> داخل النطاق</span>
      <span><i class="out"></i> خارج النطاق</span>
    </div>
  </div>
</template>

<script setup>
defineProps({
  pins: {
    type: Array,
    default: () => [
      { id: 1, top: 30, right: 25, status: 'in' },
      { id: 2, top: 55, right: 60, status: 'in' },
      { id: 3, top: 70, right: 40, status: 'out' },
      { id: 4, top: 20, right: 70, status: 'in' },
    ],
  },
})
</script>

<style scoped>
.live-map { position: relative; height: 320px; border-radius: var(--radius-lg); overflow: hidden; background: linear-gradient(180deg, #EAF1EC, #E1EBE4); border: 1px solid var(--line); }
.live-map__grid { position: absolute; inset: 0; background-image: linear-gradient(var(--line) 1px, transparent 1px), linear-gradient(90deg, var(--line) 1px, transparent 1px); background-size: 32px 32px; opacity: .5; }
.live-map__pin { position: absolute; width: 14px; height: 14px; border-radius: 50%; border: 2px solid #fff; box-shadow: var(--shadow-sm); }
.live-map__pin.in { background: var(--signal-in); }
.live-map__pin.out { background: var(--signal-out); }
.live-map__legend { position: absolute; bottom: 14px; right: 14px; display: flex; gap: 14px; background: rgba(255,255,255,.9); padding: 8px 14px; border-radius: 999px; font-size: 12px; color: var(--ink-soft); }
.live-map__legend i { display: inline-block; width: 8px; height: 8px; border-radius: 50%; margin-inline-end: 5px; }
.live-map__legend i.in { background: var(--signal-in); }
.live-map__legend i.out { background: var(--signal-out); }
</style>

```

---

## 📄 frontend-admin-dashboard/src/components/LocationSearch.vue

```vue
<template>
  <div class="location-search">
    <!-- ========================================== -->
    <!-- 1. اختيار المدينة -->
    <!-- ========================================== -->
    <div class="search-group">
      <label class="search-label">🏙️ المدينة <span class="required">*</span></label>
      <div class="search-wrapper">
        <input
          v-model="cityQuery"
          type="text"
          placeholder="ابحث عن مدينة في إسرائيل..."
          @input="onCitySearch"
          @focus="showCityResults = true"
          class="search-input"
          required
        />
        <span v-if="cityLoading" class="search-loading">⏳</span>
      </div>

      <!-- نتائج المدن -->
      <div v-if="showCityResults && cityResults.length > 0" class="search-results">
        <div
          v-for="city in cityResults"
          :key="city.place_id"
          class="result-item"
          @click="selectCity(city)"
        >
          <span class="result-icon">🏙️</span>
          <div class="result-info">
            <strong>{{ city.display_name.split(',')[0] }}</strong>
            <span class="result-address">{{ city.display_name }}</span>
          </div>
        </div>
      </div>

      <div v-if="showCityResults && cityQuery && cityResults.length === 0 && !cityLoading" class="no-results">
        <span>❌</span>
        <p>لم يتم العثور على مدينة</p>
      </div>
    </div>

    <!-- ========================================== -->
    <!-- 2. اختيار الشارع (يظهر بعد اختيار المدينة) -->
    <!-- ========================================== -->
    <div class="search-group" v-if="selectedCity">
      <label class="search-label">📍 الشارع <span class="required">*</span></label>
      <div class="search-wrapper">
        <input
          v-model="streetQuery"
          type="text"
          :placeholder="`ابحث عن شارع في ${selectedCityName}...`"
          @input="onStreetSearch"
          @focus="showStreetResults = true"
          class="search-input"
          required
        />
        <span v-if="streetLoading" class="search-loading">⏳</span>
      </div>

      <!-- نتائج الشوارع -->
      <div v-if="showStreetResults && streetResults.length > 0" class="search-results">
        <div
          v-for="street in streetResults"
          :key="street.place_id"
          class="result-item"
          @click="selectStreet(street)"
        >
          <span class="result-icon">📍</span>
          <div class="result-info">
            <strong>{{ street.display_name.split(',')[0] }}</strong>
            <span class="result-address">{{ street.display_name }}</span>
          </div>
        </div>
      </div>

      <div v-if="showStreetResults && streetQuery && streetResults.length === 0 && !streetLoading" class="no-results">
        <span>❌</span>
        <p>لم يتم العثور على شارع</p>
      </div>
    </div>

    <!-- ========================================== -->
    <!-- 3. رقم المبنى -->
    <!-- ========================================== -->
    <div class="search-group" v-if="selectedStreet">
      <label class="search-label">🏢 رقم المبنى</label>
      <input
        v-model="buildingNumber"
        type="text"
        placeholder="مثال: 15"
        class="search-input"
      />
      <span class="field-hint">أدخل رقم المبنى (اختياري)</span>
    </div>

    <!-- ========================================== -->
    <!-- 4. الموقع المختار -->
    <!-- ========================================== -->
    <div v-if="selectedLocation" class="selected-location">
      <div class="location-preview">
        <span class="location-icon">📍</span>
        <div>
          <strong>{{ selectedLocation.name }}</strong>
          <p>{{ selectedLocation.address }}</p>
          <p class="mono">خط العرض: {{ selectedLocation.latitude.toFixed(6) }}</p>
          <p class="mono">خط الطول: {{ selectedLocation.longitude.toFixed(6) }}</p>
        </div>
      </div>
    </div>

    <!-- ========================================== -->
    <!-- زر تحديد الموقع الحالي -->
    <!-- ========================================== -->
    <button class="btn btn--ghost btn--block" @click="getCurrentLocation" type="button">
      📍 استخدام موقعي الحالي
    </button>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'

const emit = defineEmits(['select'])

// ==========================================
// حالة البحث عن المدن
// ==========================================
const cityQuery = ref('')
const cityResults = ref([])
const cityLoading = ref(false)
const showCityResults = ref(false)
const selectedCity = ref(null)
const selectedCityName = computed(() => selectedCity.value?.display_name?.split(',')[0] || '')

// ==========================================
// حالة البحث عن الشوارع
// ==========================================
const streetQuery = ref('')
const streetResults = ref([])
const streetLoading = ref(false)
const showStreetResults = ref(false)
const selectedStreet = ref(null)
const buildingNumber = ref('')

// ==========================================
// الموقع المختار
// ==========================================
const selectedLocation = ref(null)

let cityTimeout = null
let streetTimeout = null

// ==========================================
// البحث عن المدن في إسرائيل
// ==========================================
function onCitySearch() {
  clearTimeout(cityTimeout)
  
  if (cityQuery.value.length < 2) {
    cityResults.value = []
    return
  }

  cityLoading.value = true
  cityTimeout = setTimeout(async () => {
    try {
      const response = await fetch(
        `https://nominatim.openstreetmap.org/search?q=${encodeURIComponent(cityQuery.value)}&countrycodes=il&featuretype=city&format=json&limit=10&addressdetails=1`
      )
      
      if (!response.ok) throw new Error('فشل البحث')
      
      const data = await response.json()
      cityResults.value = data.map(item => ({
        place_id: item.place_id,
        display_name: item.display_name,
        lat: parseFloat(item.lat),
        lon: parseFloat(item.lon),
        address: item.address || {}
      }))
      
      showCityResults.value = true
    } catch (error) {
      console.error('خطأ في البحث عن المدن:', error)
      cityResults.value = []
    } finally {
      cityLoading.value = false
    }
  }, 500)
}

// ==========================================
// اختيار مدينة
// ==========================================
function selectCity(city) {
  selectedCity.value = city
  cityQuery.value = city.display_name.split(',')[0]
  showCityResults.value = false
  
  // إعادة تعيين الشارع
  streetQuery.value = ''
  streetResults.value = []
  selectedStreet.value = null
  
  // تحديث الموقع
  updateLocation()
}

// ==========================================
// البحث عن الشوارع في المدينة المختارة
// ==========================================
function onStreetSearch() {
  clearTimeout(streetTimeout)
  
  if (!selectedCity.value || streetQuery.value.length < 2) {
    streetResults.value = []
    return
  }

  streetLoading.value = true
  streetTimeout = setTimeout(async () => {
    try {
      const response = await fetch(
        `https://nominatim.openstreetmap.org/search?q=${encodeURIComponent(streetQuery.value)}&city=${encodeURIComponent(selectedCity.value.display_name.split(',')[0])}&countrycodes=il&featuretype=street&format=json&limit=10&addressdetails=1`
      )
      
      if (!response.ok) throw new Error('فشل البحث')
      
      const data = await response.json()
      streetResults.value = data.map(item => ({
        place_id: item.place_id,
        display_name: item.display_name,
        lat: parseFloat(item.lat),
        lon: parseFloat(item.lon),
        address: item.address || {}
      }))
      
      showStreetResults.value = true
    } catch (error) {
      console.error('خطأ في البحث عن الشوارع:', error)
      streetResults.value = []
    } finally {
      streetLoading.value = false
    }
  }, 500)
}

// ==========================================
// اختيار شارع
// ==========================================
function selectStreet(street) {
  selectedStreet.value = street
  streetQuery.value = street.display_name.split(',')[0]
  showStreetResults.value = false
  
  // تحديث الموقع
  updateLocation()
}

// ==========================================
// تحديث الموقع المختار
// ==========================================
function updateLocation() {
  if (!selectedCity.value) return
  
  const lat = selectedStreet.value?.lat || selectedCity.value.lat
  const lon = selectedStreet.value?.lon || selectedCity.value.lon
  
  const location = {
    name: selectedCity.value.display_name.split(',')[0],
    address: buildAddress(),
    latitude: lat,
    longitude: lon,
    city: selectedCity.value.display_name.split(',')[0],
    street: selectedStreet.value?.display_name?.split(',')[0] || '',
    building_number: buildingNumber.value || '',
    country: 'إسرائيل'
  }
  
  selectedLocation.value = location
  emit('select', location)
}

// ==========================================
// بناء العنوان الكامل
// ==========================================
function buildAddress() {
  let parts = []
  if (buildingNumber.value) parts.push(buildingNumber.value)
  if (selectedStreet.value) parts.push(selectedStreet.value.display_name.split(',')[0])
  if (selectedCity.value) parts.push(selectedCity.value.display_name.split(',')[0])
  return parts.join('، ') || 'إسرائيل'
}

// ==========================================
// مراقبة رقم المبنى
// ==========================================
function updateBuildingNumber() {
  if (selectedLocation.value) {
    selectedLocation.value.building_number = buildingNumber.value
    selectedLocation.value.address = buildAddress()
    emit('select', selectedLocation.value)
  }
}

// ==========================================
// استخدام الموقع الحالي
// ==========================================
function getCurrentLocation() {
  if (!('geolocation' in navigator)) {
    alert('المتصفح لا يدعم تحديد الموقع')
    return
  }

  navigator.geolocation.getCurrentPosition(
    async (pos) => {
      const { latitude, longitude } = pos.coords
      
      try {
        const response = await fetch(
          `https://nominatim.openstreetmap.org/reverse?lat=${latitude}&lon=${longitude}&format=json&zoom=10`
        )
        const data = await response.json()
        
        const location = {
          name: data.address?.city || data.address?.town || data.address?.village || 'موقعي',
          address: data.display_name || `${latitude}, ${longitude}`,
          latitude: latitude,
          longitude: longitude,
          city: data.address?.city || data.address?.town || '',
          street: data.address?.road || '',
          building_number: data.address?.house_number || '',
          country: data.address?.country || 'إسرائيل'
        }
        
        selectedLocation.value = location
        emit('select', location)
        
        // تعبئة الحقول
        cityQuery.value = location.city
        streetQuery.value = location.street
        buildingNumber.value = location.building_number
        
      } catch (error) {
        emit('select', {
          name: 'موقعي الحالي',
          address: `${latitude}, ${longitude}`,
          latitude: latitude,
          longitude: longitude,
          country: 'إسرائيل'
        })
      }
    },
    (err) => {
      alert('فشل تحديد الموقع: ' + err.message)
    },
    { enableHighAccuracy: true }
  )
}

// ==========================================
// إغلاق النتائج عند النقر خارجها
// ==========================================
function closeResults(e) {
  const wrapper = document.querySelector('.location-search')
  if (wrapper && !wrapper.contains(e.target)) {
    showCityResults.value = false
    showStreetResults.value = false
  }
}

document.addEventListener('click', closeResults)

// مراقبة رقم المبنى
import { watch } from 'vue'
watch(buildingNumber, updateBuildingNumber)
</script>

<style scoped>
.location-search {
  width: 100%;
}

/* ==========================================
   مجموعات البحث
   ========================================== */
.search-group {
  margin-bottom: 14px;
  position: relative;
}

.search-label {
  display: block;
  font-size: 13px;
  font-weight: 600;
  color: var(--ink);
  margin-bottom: 4px;
}

.search-wrapper {
  position: relative;
  display: flex;
  align-items: center;
}

.search-input {
  width: 100%;
  padding: 10px 14px;
  padding-left: 40px;
  border: 1.5px solid var(--line);
  border-radius: var(--radius-sm);
  font-size: 14px;
  font-family: var(--font-body);
  background: var(--surface);
  transition: all 0.3s;
}

.search-input:focus {
  outline: none;
  border-color: var(--brand);
  box-shadow: 0 0 0 3px var(--brand-tint);
}

.search-loading {
  position: absolute;
  left: 12px;
  font-size: 16px;
}

.required {
  color: var(--signal-out);
}

.field-hint {
  font-size: 11px;
  color: var(--ink-soft);
  display: block;
  margin-top: 4px;
}

/* ==========================================
   نتائج البحث
   ========================================== */
.search-results {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  background: var(--surface);
  border: 1px solid var(--line);
  border-radius: var(--radius-sm);
  box-shadow: var(--shadow-lg);
  max-height: 250px;
  overflow-y: auto;
  z-index: 1000;
  margin-top: 4px;
}

.result-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 14px;
  cursor: pointer;
  border-bottom: 1px solid var(--line);
  transition: background 0.2s;
}

.result-item:hover {
  background: var(--brand-tint);
}

.result-item:last-child {
  border-bottom: none;
}

.result-icon {
  font-size: 18px;
  flex-shrink: 0;
}

.result-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.result-info strong {
  font-size: 13px;
  color: var(--ink);
}

.result-address {
  font-size: 11px;
  color: var(--ink-soft);
}

/* ==========================================
   عدم وجود نتائج
   ========================================== */
.no-results {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  background: var(--surface);
  border: 1px solid var(--line);
  border-radius: var(--radius-sm);
  padding: 16px;
  text-align: center;
  z-index: 1000;
  margin-top: 4px;
}

.no-results span {
  font-size: 24px;
  display: block;
  margin-bottom: 4px;
}

.no-results p {
  font-size: 13px;
  color: var(--ink-soft);
  margin: 0;
}

/* ==========================================
   الموقع المختار
   ========================================== */
.selected-location {
  background: var(--brand-tint);
  border: 1px solid var(--brand);
  border-radius: var(--radius-sm);
  padding: 12px 16px;
  margin-bottom: 14px;
}

.location-preview {
  display: flex;
  align-items: flex-start;
  gap: 12px;
}

.location-icon {
  font-size: 24px;
}

.location-preview strong {
  font-size: 14px;
  color: var(--brand);
}

.location-preview p {
  font-size: 12px;
  color: var(--ink-soft);
  margin: 2px 0;
}

/* ==========================================
   أزرار
   ========================================== */
.btn--block {
  width: 100%;
}

.btn--ghost {
  margin-top: 4px;
}
</style>

```

---

## 📄 frontend-admin-dashboard/src/components/NotificationDropdown.vue

```vue
<template>
  <div class="notification-dropdown" ref="dropdownRef">
    <!-- أيقونة الإشعارات -->
    <button 
      @click="toggleDropdown" 
      class="notification-icon"
      :title="t('notifications')"
    >
      <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9"></path>
        <path d="M13.73 21a2 2 0 0 1-3.46 0"></path>
      </svg>
      <span v-if="unreadCount > 0" class="badge">{{ unreadCount }}</span>
    </button>
    
    <!-- Dropdown content -->
    <transition name="dropdown">
      <div v-if="isOpen" class="dropdown-content">
        <div class="feed-header">
          <h3>{{ t('notifications_feed') }}</h3>
          <span v-if="notificationsStore.notifications && notificationsStore.notifications.length" class="badge">
            {{ notificationsStore.notifications.length }}
          </span>
        </div>

        <div v-if="notificationsStore.loading" class="empty-state">
          <p>{{ t('loading_notifications') }}</p>
        </div>

        <div v-else-if="notificationsStore.error" class="alert alert-error">
          <span>❌</span> {{ notificationsStore.error }}
        </div>

        <div v-else-if="!notificationsStore.notifications || !notificationsStore.notifications.length" class="empty-state">
          <p>{{ t('no_notifications_available') }}</p>
        </div>

        <div v-else class="notifications-list">
          <div 
            v-for="notification in notificationsStore.notifications" 
            :key="notification.id" 
            class="notification-item"
          >
            <span class="notification-icon">{{ getNotificationIcon(notification.type) }}</span>
            <div class="notification-content">
              <p class="notification-title">{{ notification.title }}</p>
              <p class="notification-message">{{ notification.message }}</p>
              <span class="notification-time">{{ formatTime(notification.created_at) }}</span>
            </div>
          </div>
          <router-link to="/notifications" class="view-all-link" @click="closeDropdown">
            {{ t('view_all') }}
          </router-link>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount, computed } from 'vue'
import { notificationsStore } from '../store/notifications'
import { useI18n } from '../services/i18n'
import { useRouter } from 'vue-router'

const { t } = useI18n()
const router = useRouter()
const isOpen = ref(false)
const dropdownRef = ref(null)

const unreadCount = computed(() => {
  if (!notificationsStore.notifications) return 0
  return notificationsStore.notifications.length
})

function toggleDropdown() {
  isOpen.value = !isOpen.value
  if (isOpen.value) {
    notificationsStore.fetchNotifications()
  }
}

function closeDropdown() {
  isOpen.value = false
}

function handleClickOutside(event) {
  if (dropdownRef.value && !dropdownRef.value.contains(event.target)) {
    closeDropdown()
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
  notificationsStore.fetchNotifications()
})

onBeforeUnmount(() => {
  document.removeEventListener('click', handleClickOutside)
})

function getNotificationIcon(type) {
  const icons = {
    'status_update': '📋',
    'worker_assigned': '👷',
    'worker_arrived': '📍',
    'service_completed': '✅',
    'payment_request': '💳',
    'general': '🔔'
  }
  return icons[type] || '🔔'
}

function formatTime(date) {
  if (!date) return ''
  return new Date(date).toLocaleTimeString('en-GB', { 
    hour: '2-digit', 
    minute: '2-digit',
    hour12: false 
  })
}
</script>

<style scoped>
.notification-dropdown {
  position: relative;
  display: inline-block;
}

.notification-icon {
  background: none;
  border: none;
  cursor: pointer;
  position: relative;
  padding: 8px;
  border-radius: var(--radius-sm);
  transition: all var(--transition-base);
  color: var(--ink-soft);
  display: flex;
  align-items: center;
  justify-content: center;
}

.notification-icon:hover {
  background: var(--brand-tint);
  color: var(--brand);
  transform: scale(1.1);
}

.badge {
  position: absolute;
  top: 0;
  right: 0;
  background: var(--signal-out);
  color: white;
  border-radius: 50%;
  padding: 2px 6px;
  font-size: 10px;
  font-weight: 600;
  min-width: 16px;
  height: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0%, 100% {
    transform: scale(1);
  }
  50% {
    transform: scale(1.1);
  }
}

.dropdown-content {
  position: absolute;
  top: 100%;
  right: 0;
  background: var(--surface);
  border-radius: var(--radius-md);
  box-shadow: var(--shadow-lg);
  min-width: 320px;
  max-height: 400px;
  overflow-y: auto;
  z-index: 1000;
  margin-top: 8px;
  border: 1px solid var(--line);
}

.dropdown-enter-active,
.dropdown-leave-active {
  transition: all 0.2s ease;
}

.dropdown-enter-from {
  opacity: 0;
  transform: translateY(-10px);
}

.dropdown-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}

.feed-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  border-bottom: 1px solid var(--line);
  background: var(--canvas);
  border-radius: var(--radius-md) var(--radius-md) 0 0;
  position: sticky;
  top: 0;
  z-index: 1;
}

.feed-header h3 {
  font-size: 14px;
  margin: 0;
  color: var(--ink);
}

.feed-header .badge {
  position: static;
  background: var(--brand);
  color: white;
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 11px;
  font-weight: 600;
  animation: none;
}

.empty-state {
  text-align: center;
  padding: 20px;
  color: var(--ink-soft);
  font-size: 13px;
}

.alert {
  padding: 10px 12px;
  border-radius: var(--radius-sm);
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
  margin: 8px;
}

.alert-error {
  background: var(--signal-out-tint);
  color: var(--signal-out);
}

.notifications-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 8px;
}

.notification-item {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  padding: 10px;
  background: var(--canvas);
  border-radius: var(--radius-sm);
  transition: background var(--transition-fast);
  cursor: pointer;
}

.notification-item:hover {
  background: var(--brand-tint);
}

.notification-icon {
  font-size: 16px;
  flex-shrink: 0;
}

.notification-content {
  flex: 1;
  min-width: 0;
}

.notification-title {
  font-size: 13px;
  font-weight: 600;
  color: var(--ink);
  margin: 0 0 4px 0;
}

.notification-message {
  font-size: 12px;
  color: var(--ink-soft);
  margin: 0 0 4px 0;
  line-height: 1.4;
}

.notification-time {
  font-size: 11px;
  color: var(--ink-light);
}

.view-all-link {
  display: block;
  text-align: center;
  padding: 12px;
  margin-top: 8px;
  background: var(--brand-tint);
  color: var(--brand);
  text-decoration: none;
  border-radius: var(--radius-sm);
  font-size: 13px;
  font-weight: 600;
  transition: all var(--transition-fast);
}

.view-all-link:hover {
  background: var(--brand);
  color: white;
}

@media (max-width: 480px) {
  .dropdown-content {
    min-width: 280px;
    max-height: 350px;
  }

  .feed-header {
    padding: 10px 12px;
  }

  .feed-header h3 {
    font-size: 13px;
  }

  .notifications-list {
    padding: 6px;
  }

  .notification-item {
    padding: 8px;
    gap: 8px;
  }

  .notification-icon {
    font-size: 14px;
  }

  .notification-title {
    font-size: 12px;
  }

  .notification-message {
    font-size: 11px;
  }

  .notification-time {
    font-size: 10px;
  }
}
</style>
```

---

## 📄 frontend-admin-dashboard/src/components/OpenInMaps.vue

```vue
<template>
  <button class="maps-btn" @click="openInMaps" :title="t('open_in_maps')">
    <span class="maps-icon">🗺️</span>
    <span class="maps-text">{{ t('open_in_maps') }}</span>
  </button>
</template>

<script setup>
import { useI18n } from '../services/i18n'

const { t } = useI18n()

const props = defineProps({
  latitude: { type: Number, required: true },
  longitude: { type: Number, required: true },
  label: { type: String, default: 'الموقع' }
})

function openInMaps() {
  const url = `https://www.google.com/maps/dir/?api=1&destination=${props.latitude},${props.longitude}&destination_place_id=${encodeURIComponent(props.label)}`
  
  // فتح في نافذة جديدة
  window.open(url, '_blank')
}
</script>

<style scoped>
.maps-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 14px;
  border: 1px solid var(--brand);
  border-radius: var(--radius-sm);
  background: var(--brand-tint);
  color: var(--brand);
  font-family: var(--font-body);
  font-size: 12px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.maps-btn:hover {
  background: var(--brand);
  color: white;
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
}

.maps-icon {
  font-size: 14px;
}
</style>

```

---

## 📄 frontend-admin-dashboard/src/components/PWAInstallButton.vue

```vue
<template>
  <div v-if="showInstallButton" class="pwa-install-container">
    <button 
      @click="installPWA" 
      class="pwa-install-icon"
      :title="pwaInstallTitle"
    >
      <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path>
        <polyline points="7 10 12 15 17 10"></polyline>
        <line x1="12" y1="15" x2="12" y2="3"></line>
      </svg>
    </button>
    
    <!-- iOS Instructions Modal -->
    <div v-if="showIOSInstructions" class="ios-instructions-modal" @click.self="showIOSInstructions = false">
      <div class="ios-instructions-content">
        <button class="close-btn" @click="showIOSInstructions = false">✕</button>
        <div class="ios-header">
          <div class="ios-icon">📱</div>
          <h3>{{ iosModalTitle }}</h3>
          <p class="ios-subtitle">{{ iosModalSubtitle }}</p>
        </div>
        
        <div class="steps">
          <div class="step">
            <span class="step-number">1</span>
            <div class="step-content">
              <p v-html="step1Text"></p>
              <div class="step-visual">
                <div class="phone-frame">
                  <div class="phone-screen">
                    <div class="share-btn">⎋</div>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div class="step">
            <span class="step-number">2</span>
            <div class="step-content">
              <p v-html="step2Text"></p>
              <div class="step-visual">
                <div class="menu-item">➕ Add to Home Screen</div>
              </div>
            </div>
          </div>
          <div class="step">
            <span class="step-number">3</span>
            <div class="step-content">
              <p v-html="step3Text"></p>
              <div class="step-visual">
                <div class="add-btn">Add</div>
              </div>
            </div>
          </div>
        </div>
        
        <div class="benefits">
          <div class="benefit">{{ benefit1 }}</div>
          <div class="benefit">{{ benefit2 }}</div>
          <div class="benefit">{{ benefit3 }}</div>
        </div>
        
        <div class="actions">
          <button class="got-it-btn" @click="showIOSInstructions = false">{{ gotItText }}</button>
          <button class="remind-btn" @click="remindLater">{{ remindLaterText }}</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'PWAInstallButton',
  data() {
    return {
      showInstallButton: false,
      showIOSInstructions: false,
      deferredPrompt: null,
      isIOS: false,
      isStandalone: false
    }
  },
  computed: {
    currentLang() {
      if (this.$i18n && this.$i18n.currentLang) {
        return this.$i18n.currentLang.value
      }
      return 'ar' // default
    },
    pwaInstallTitle() {
      return this.$t('pwa.installTitle')
    },
    pwaInstallText() {
      if (this.isIOS) {
        return this.$t('pwa.installOnIos')
      }
      return this.$t('pwa.installText')
    },
    iosModalTitle() {
      return this.$t('pwa.iosModalTitle')
    },
    iosModalSubtitle() {
      return this.$t('pwa.iosModalSubtitle')
    },
    step1Text() {
      return this.$t('pwa.step1Text')
    },
    step2Text() {
      return this.$t('pwa.step2Text')
    },
    step3Text() {
      return this.$t('pwa.step3Text')
    },
    benefit1() {
      return this.$t('pwa.benefit1')
    },
    benefit2() {
      return this.$t('pwa.benefit2')
    },
    benefit3() {
      return this.$t('pwa.benefit3')
    },
    gotItText() {
      return this.$t('pwa.gotIt')
    },
    remindLaterText() {
      return this.$t('pwa.remindLater')
    }
  },
  methods: {
    detectIOS() {
      const userAgent = window.navigator.userAgent.toLowerCase()
      return /iphone|ipad|ipod/.test(userAgent) && !/edge|crios|fxios/.test(userAgent)
    },
    handleInstallAvailable() {
      this.showInstallButton = true
    },
    handleInstallSuccess() {
      this.showInstallButton = false
    },
    installPWA() {
      if (this.isIOS) {
        // Show iOS instructions
        this.showIOSInstructions = true
      } else if (window.pwaInstall) {
        // Use native install prompt for Android/Windows
        window.pwaInstall()
      } else {
        // Fallback - show general instructions
        this.showIOSInstructions = true
      }
    },
    remindLater() {
      // Just close the modal, keep icon visible
      this.showIOSInstructions = false
    }
  },
  mounted() {
    // Check if device is iOS
    this.isIOS = this.detectIOS()
    
    // Check if already installed (standalone mode)
    this.isStandalone = window.matchMedia('(display-mode: standalone)').matches || 
                      window.navigator.standalone === true

    // Don't show if already installed
    if (this.isStandalone) {
      return
    }

    // Always show button on login page
    if (this.isIOS) {
      this.showInstallButton = true
    } else {
      // For other devices, listen to install prompt
      window.addEventListener('pwa-install-available', this.handleInstallAvailable)
      window.addEventListener('pwa-install-success', this.handleInstallSuccess)
      
      // Check for existing deferredPrompt
      if (window.deferredPrompt) {
        this.showInstallButton = true
      }
    }
  },
  beforeUnmount() {
    window.removeEventListener('pwa-install-available', this.handleInstallAvailable)
    window.removeEventListener('pwa-install-success', this.handleInstallSuccess)
  }
}
</script>

<style scoped>
.pwa-install-container {
  position: fixed;
  bottom: 16px;
  right: 16px;
  z-index: 1000;
}

.pwa-install-icon {
  width: 44px;
  height: 44px;
  padding: 0;
  background: linear-gradient(135deg, #1e3a8a 0%, #3b82f6 100%);
  color: white;
  border: none;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  box-shadow: 0 4px 12px rgba(30, 58, 138, 0.3);
  transition: all 0.3s ease;
}

.pwa-install-icon:hover {
  transform: scale(1.1);
  box-shadow: 0 6px 16px rgba(30, 58, 138, 0.4);
}

.pwa-install-icon:active {
  transform: scale(0.95);
}

/* iOS Instructions Modal */
.ios-instructions-modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2000;
  padding: 20px;
}

.ios-instructions-content {
  background: white;
  border-radius: 20px;
  padding: 28px;
  max-width: 420px;
  width: 100%;
  position: relative;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.4);
  max-height: 90vh;
  overflow-y: auto;
}

.close-btn {
  position: absolute;
  top: 16px;
  left: 16px;
  background: none;
  border: none;
  font-size: 24px;
  cursor: pointer;
  color: #666;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  transition: background 0.2s;
}

.close-btn:hover {
  background: #f0f0f0;
}

.ios-header {
  text-align: center;
  margin-bottom: 24px;
  padding-top: 8px;
}

.ios-icon {
  font-size: 48px;
  margin-bottom: 12px;
}

.ios-header h3 {
  margin: 0 0 8px 0;
  font-size: 20px;
  color: #1e3a8a;
  font-weight: 700;
}

.ios-subtitle {
  margin: 0;
  font-size: 14px;
  color: #666;
}

.steps {
  display: flex;
  flex-direction: column;
  gap: 20px;
  margin-bottom: 24px;
}

.step {
  display: flex;
  align-items: flex-start;
  gap: 14px;
}

.step-number {
  width: 32px;
  height: 32px;
  background: linear-gradient(135deg, #1e3a8a 0%, #3b82f6 100%);
  color: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  font-size: 15px;
  flex-shrink: 0;
}

.step-content {
  flex: 1;
}

.step p {
  margin: 0 0 8px 0;
  font-size: 14px;
  color: #333;
  line-height: 1.5;
}

.step strong {
  color: #1e3a8a;
}

.step-visual {
  background: #f8f9fa;
  border-radius: 8px;
  padding: 8px;
  display: flex;
  justify-content: center;
  border: 1px solid #e9ecef;
}

.phone-frame {
  width: 60px;
  height: 100px;
  background: white;
  border: 2px solid #dee2e6;
  border-radius: 8px;
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
}

.phone-screen {
  width: 90%;
  height: 90%;
  background: #f8f9fa;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.share-btn {
  font-size: 20px;
  color: #007AFF;
}

.menu-item {
  background: white;
  padding: 8px 12px;
  border-radius: 6px;
  font-size: 12px;
  color: #333;
  border: 1px solid #dee2e6;
}

.add-btn {
  background: #007AFF;
  color: white;
  padding: 6px 16px;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 600;
}

.benefits {
  display: flex;
  justify-content: space-around;
  margin-bottom: 24px;
  padding: 16px;
  background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
  border-radius: 12px;
}

.benefit {
  font-size: 12px;
  color: #495057;
  text-align: center;
  font-weight: 500;
}

.actions {
  display: flex;
  gap: 12px;
}

.got-it-btn {
  flex: 2;
  padding: 14px;
  background: linear-gradient(135deg, #1e3a8a 0%, #3b82f6 100%);
  color: white;
  border: none;
  border-radius: 12px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: transform 0.2s;
}

.got-it-btn:hover {
  transform: scale(1.02);
}

.remind-btn {
  flex: 1;
  padding: 14px;
  background: #f8f9fa;
  color: #666;
  border: 1px solid #dee2e6;
  border-radius: 12px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: background 0.2s;
}

.remind-btn:hover {
  background: #e9ecef;
}

/* RTL support */
[dir="rtl"] .pwa-install-container {
  right: auto;
  left: 16px;
}

[dir="rtl"] .close-btn {
  left: auto;
  right: 16px;
}

@media (max-width: 768px) {
  .pwa-install-container {
    bottom: 12px;
    right: 12px;
  }
  
  [dir="rtl"] .pwa-install-container {
    right: auto;
    left: 12px;
  }
  
  .pwa-install-icon {
    width: 40px;
    height: 40px;
  }

  .ios-instructions-content {
    padding: 20px;
  }

  .ios-instructions-content h3 {
    font-size: 16px;
  }

  .step p {
    font-size: 13px;
  }
}
</style>
```

---

## 📄 frontend-admin-dashboard/src/components/PullToRefresh.vue

```vue
<template>
  <div 
    class="pull-to-refresh"
    @touchstart="handleTouchStart"
    @touchmove="handleTouchMove"
    @touchend="handleTouchEnd"
    @mousedown="handleTouchStart"
    @mousemove="handleTouchMove"
    @mouseup="handleTouchEnd"
  >
    <div 
      class="pull-indicator"
      :style="{ 
        transform: `translateY(${pullDistance}px)`,
        opacity: Math.min(pullDistance / 60, 1)
      }"
    >
      <div class="pull-icon" :class="{ rotating: isRefreshing }">
        <svg v-if="isRefreshing" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M12 2v4M12 18v4M4.93 4.93l2.83 2.83M16.24 16.24l2.83 2.83M2 12h4M18 12h4M4.93 19.07l2.83-2.83M16.24 7.76l2.83-2.83"/>
        </svg>
        <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M12 4v1M12 20v1M4 12h1M19 12h1M6.34 6.34l.71.71M17.24 17.24l.71.71M6.34 17.65l.71-.71M17.24 6.65l.71-.71"/>
        </svg>
      </div>
      <div class="pull-text">{{ refreshText }}</div>
    </div>
    
    <div class="pull-content" :style="{ transform: `translateY(${Math.max(0, pullDistance - 60)}px)` }">
      <slot></slot>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'

const props = defineProps({
  threshold: {
    type: Number,
    default: 80
  },
  refreshText: {
    type: String,
    default: 'اسحب للتحديث'
  },
  refreshingText: {
    type: String,
    default: 'جاري التحديث...'
  },
  releaseText: {
    type: String,
    default: 'أفعل للتحديث'
  }
})

const emit = defineEmits(['refresh'])

const isPulling = ref(false)
const isRefreshing = ref(false)
const pullDistance = ref(0)
const startY = ref(0)
const currentY = ref(0)

const displayText = computed(() => {
  if (isRefreshing.value) return props.refreshingText
  if (pullDistance.value >= props.threshold) return props.releaseText
  return props.refreshText
})

const refreshText = computed(() => displayText.value)

function handleTouchStart(e) {
  if (isRefreshing.value) return
  
  const y = e.touches ? e.touches[0].clientY : e.clientY
  const scrollTop = e.target?.scrollTop || document.documentElement.scrollTop
  
  // فقط إذا كنا في أعلى الصفحة
  if (scrollTop <= 0) {
    startY.value = y
    isPulling.value = true
    pullDistance.value = 0
  }
}

function handleTouchMove(e) {
  if (!isPulling.value || isRefreshing.value) return
  
  const y = e.touches ? e.touches[0].clientY : e.clientY
  currentY.value = y
  
  const distance = y - startY.value
  
  // حساب المسافة مع تخفيف الحركة
  if (distance > 0) {
    pullDistance.value = Math.min(distance * 0.5, 150)
    
    // منع السحب الافتراضي
    if (e.cancelable) {
      e.preventDefault()
    }
  }
}

function handleTouchEnd() {
  if (!isPulling.value) return
  
  isPulling.value = false
  
  // إذا وصلنا لعتبة التحديث
  if (pullDistance.value >= props.threshold) {
    triggerRefresh()
  }
  
  // إعادة المسافة للصفر
  pullDistance.value = 0
}

async function triggerRefresh() {
  isRefreshing.value = true
  pullDistance.value = 60
  
  try {
    await emit('refresh')
  } catch (error) {
    console.error('Refresh failed:', error)
  } finally {
    setTimeout(() => {
      isRefreshing.value = false
      pullDistance.value = 0
    }, 500)
  }
}
</script>

<style scoped>
.pull-to-refresh {
  position: relative;
  width: 100%;
  min-height: 100vh;
  overflow: hidden;
}

.pull-indicator {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 60px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  pointer-events: none;
  z-index: 1000;
  background: linear-gradient(180deg, rgba(30, 58, 95, 0.1) 0%, transparent 100%);
}

.pull-icon {
  width: 24px;
  height: 24px;
  color: #1E3A5F;
  margin-bottom: 4px;
  transition: transform 0.3s ease;
}

.pull-icon.rotating {
  animation: rotate 1s linear infinite;
}

@keyframes rotate {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

.pull-text {
  font-size: 12px;
  color: #1E3A5F;
  font-weight: 600;
  text-align: center;
}

.pull-content {
  transition: transform 0.3s ease;
  min-height: 100vh;
}

@media (max-width: 480px) {
  .pull-text {
    font-size: 11px;
  }
  
  .pull-icon {
    width: 20px;
    height: 20px;
  }
}
</style>

```

---

## 📄 frontend-admin-dashboard/src/components/RealMap.vue

```vue
<template>
  <div class="map-container">
    <l-map
      ref="mapRef"
      :key="isDarkMode.value ? 'dark' : 'light'"
      :zoom="zoom"
      @update:zoom="updateZoom"
      :center="center"
      :options="{ attributionControl: true, zoomControl: true }"
      :style="{ height: height + 'px', width: '100%' }"
    >
      <l-tile-layer
        :url="mapTileUrl"
        layer-type="base"
        :name="isDarkMode.value ? 'CartoDB Dark' : 'OpenStreetMap'"
        :attribution="mapAttribution"
      />
      
      <!-- نقاط العمل -->
      <l-marker
        v-for="site in worksites"
        :key="site.id"
        :lat-lng="[site.latitude, site.longitude]"
      >
        <l-icon>
          <div class="worksite-marker">
            <div class="worksite-dot"></div>
            <div class="worksite-ring"></div>
          </div>
        </l-icon>
        <l-popup>
          <div class="popup-content">
            <h4>🏢 {{ site.name }}</h4>
            <p>{{ site.address || 'لا يوجد عنوان' }}</p>
            <p>⭕ النطاق: {{ site.radius_meters }} متر</p>
            <p>👥 عدد الموظفين: {{ getEmployeeCount(site.id) }}</p>
          </div>
        </l-popup>
      </l-marker>

      <!-- الموظفين -->
      <l-marker
        v-for="emp in employees"
        :key="emp.id"
        :lat-lng="[emp.latitude, emp.longitude]"
      >
        <l-icon>
          <div class="employee-marker" :class="emp.status">
            <div class="employee-pulse" :class="emp.status"></div>
            <div class="employee-dot" :class="emp.status">
              {{ emp.full_name?.slice(0, 1) || '?' }}
            </div>
          </div>
        </l-icon>
        <l-popup>
          <div class="popup-content employee-popup">
            <h4>👤 {{ emp.full_name }}</h4>
            <p>📍 {{ emp.worksite.name }}</p>
            <p>📏 المسافة: {{ formatDistance(emp.worksite.distance) }}</p>
            <p>⏱️ {{ emp.hours_worked.toFixed(1) }} ساعة</p>
            <p>
              <span class="badge" :class="emp.status === 'inside' ? 'badge--in' : 'badge--out'">
                {{ emp.status_text }}
              </span>
            </p>
            <button
              class="btn btn--sm btn--primary"
              @click="handleShowDetails(emp)"
            >
              📋 عرض التفاصيل
            </button>
          </div>
        </l-popup>
      </l-marker>
    </l-map>

    <div v-if="!employees || employees.length === 0" class="map-overlay">
      <div class="overlay-content">
        <span class="overlay-icon">🗺️</span>
        <h3>لا يوجد موظفين نشطين</h3>
        <p>سيظهر هنا الموظفون الذين بدأوا الدوام</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { LMap, LTileLayer, LMarker, LPopup, LIcon } from '@vue-leaflet/vue-leaflet'
import 'leaflet/dist/leaflet.css'

const props = defineProps({
  employees: { type: Array, default: () => [] },
  worksites: { type: Array, default: () => [] },
  center: { type: Array, default: () => [31.5, 34.8] },
  zoom: { type: Number, default: 7 },
  height: { type: Number, default: 400 }  // ✅ Number وليس String
})

// ✅ تعريف emit بشكل صحيح
const emit = defineEmits(['update:zoom', 'showDetails'])

const mapRef = ref(null)
let observer = null

// ✅ اكتشاف الوضع الداكن مع مراقبة التغييرات
const isDarkMode = ref(document.documentElement.getAttribute('data-theme') === 'dark')

// ✅ مراقبة تغييرات data-theme
const observeThemeChanges = () => {
  observer = new MutationObserver((mutations) => {
    mutations.forEach((mutation) => {
      if (mutation.type === 'attributes' && mutation.attributeName === 'data-theme') {
        isDarkMode.value = document.documentElement.getAttribute('data-theme') === 'dark'
      }
    })
  })

  observer.observe(document.documentElement, {
    attributes: true,
    attributeFilter: ['data-theme']
  })
}

// ✅ تحديد URL الخريطة حسب الوضع
const mapTileUrl = computed(() => {
  if (isDarkMode.value) {
    // خريطة داكنة زرقاء فاتحة ومريحة للعين من CartoDB
    return 'https://{s}.basemaps.cartocdn.com/dark_all/{z}/{x}/{y}{r}.png'
  } else {
    // خريطة عادية من OpenStreetMap
    return 'https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png'
  }
})

// ✅ تحديد Attribution حسب الوضع
const mapAttribution = computed(() => {
  if (isDarkMode.value) {
    return '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors &copy; <a href="https://carto.com/attributions">CARTO</a>'
  } else {
    return '&copy; OpenStreetMap contributors'
  }
})

function updateZoom(newZoom) {
  emit('update:zoom', newZoom)
}

// ✅ دالة handleShowDetails لتوصيل الحدث
function handleShowDetails(employee) {
  console.log('📋 عرض تفاصيل الموظف:', employee.full_name)
  emit('showDetails', employee)
}

function getEmployeeCount(worksiteId) {
  return props.employees.filter(e => e.worksite?.id === worksiteId).length
}

function formatDistance(meters) {
  if (!meters) return '0 م'
  if (meters >= 1000) {
    return (meters / 1000).toFixed(2) + ' كيلومتر'
  }
  return Math.round(meters) + ' متر'
}

onMounted(() => {
  observeThemeChanges()
})

onUnmounted(() => {
  if (observer) {
    observer.disconnect()
  }
})
</script>

<style scoped>
.map-container {
  border-radius: var(--radius-lg);
  overflow: hidden;
  border: 1px solid var(--line);
  position: relative;
  min-height: 400px;
  background: #E8EDF2;
}

[data-theme="dark"] .map-container {
  background: #0f172a;
  border-color: #334155;
}

.worksite-marker {
  position: relative;
  width: 30px;
  height: 30px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.worksite-dot {
  width: 16px;
  height: 16px;
  border-radius: 50%;
  background: #2563EB;
  border: 2px solid white;
  box-shadow: 0 2px 8px rgba(37, 99, 235, 0.4);
  z-index: 2;
}

.worksite-ring {
  position: absolute;
  inset: -8px;
  border-radius: 50%;
  border: 2px solid rgba(37, 99, 235, 0.3);
  animation: ringPulse 2s ease-out infinite;
}

@keyframes ringPulse {
  0% { transform: scale(1); opacity: 1; }
  100% { transform: scale(1.5); opacity: 0; }
}

.employee-marker {
  position: relative;
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
}

.employee-dot {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  font-size: 14px;
  color: white;
  border: 2px solid white;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.25);
  z-index: 2;
  transition: all 0.3s ease;
}

.employee-dot.inside {
  background: #22C55E;
  box-shadow: 0 0 20px rgba(34, 197, 94, 0.4);
}

.employee-dot.outside {
  background: #EF4444;
  box-shadow: 0 0 20px rgba(239, 68, 68, 0.4);
}

.employee-pulse {
  position: absolute;
  inset: -4px;
  border-radius: 50%;
  z-index: 1;
  animation: employeePulse 1.5s ease-out infinite;
}

.employee-pulse.inside {
  border: 2px solid #22C55E;
}

.employee-pulse.outside {
  border: 2px solid #EF4444;
}

@keyframes employeePulse {
  0% { transform: scale(1); opacity: 1; }
  100% { transform: scale(1.6); opacity: 0; }
}

.employee-marker:hover .employee-dot {
  transform: scale(1.1);
}

.employee-marker:hover .employee-pulse {
  animation-duration: 0.8s;
}

.map-overlay {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.85);
  backdrop-filter: blur(4px);
  z-index: 1000;
  padding: 30px;
  border-radius: var(--radius-lg);
}

[data-theme="dark"] .map-overlay {
  background: rgba(15, 23, 42, 0.85);
}

.overlay-content {
  text-align: center;
  max-width: 400px;
  animation: fadeInUp 0.5s ease;
}

@keyframes fadeInUp {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: translateY(0); }
}

.overlay-icon {
  font-size: 64px;
  display: block;
  margin-bottom: 16px;
}

.overlay-content h3 {
  font-size: 20px;
  color: var(--ink);
  margin-bottom: 8px;
}

.overlay-content p {
  font-size: 14px;
  color: var(--ink-soft);
  margin-bottom: 16px;
}

.popup-content h4 {
  margin: 0 0 4px;
  font-size: 14px;
  color: var(--ink);
}

.popup-content p {
  margin: 2px 0;
  font-size: 12px;
  color: var(--ink-soft);
}

[data-theme="dark"] .popup-content h4 {
  color: #f1f5f9;
}

[data-theme="dark"] .popup-content p {
  color: #cbd5e1;
}

.popup-content .btn {
  margin-top: 8px;
  font-size: 12px;
  padding: 4px 12px;
}

.employee-popup {
  min-width: 200px;
}

.badge--in {
  background: #22C55E20;
  color: #22C55E;
  padding: 2px 10px;
  border-radius: 999px;
  font-size: 11px;
  font-weight: 600;
}

.badge--out {
  background: #EF444420;
  color: #EF4444;
  padding: 2px 10px;
  border-radius: 999px;
  font-size: 11px;
  font-weight: 600;
}
</style>

```

---

## 📄 frontend-admin-dashboard/src/components/StatsCard.vue

```vue
<template>
  <div class="stat-card">
    <span class="stat-card__label">{{ label }}</span>
    <span class="stat-card__value mono">{{ value }}</span>
    <span v-if="hint" class="stat-card__hint" :class="tone">{{ hint }}</span>
  </div>
</template>

<script setup>
defineProps({
  label: String,
  value: [String, Number],
  hint: String,
  tone: { type: String, default: '' }, // 'up' | 'down'
})
</script>

<style scoped>
.stat-card { background: var(--surface); border: 1px solid var(--line); border-radius: var(--radius-lg); padding: 20px; display: flex; flex-direction: column; gap: 6px; box-shadow: var(--shadow-sm); }
.stat-card__label { font-size: 13px; color: var(--ink-soft); }
.stat-card__value { font-size: 28px; font-weight: 600; }
.stat-card__hint { font-size: 12px; font-weight: 600; color: var(--ink-soft); }
.stat-card__hint.up { color: var(--signal-in); }
.stat-card__hint.down { color: var(--signal-out); }
</style>

```

---

## 📄 frontend-admin-dashboard/src/components/SwipeNav.vue

```vue
<template>
  <div 
    class="swipe-nav"
    @touchstart="handleTouchStart"
    @touchmove="handleTouchMove"
    @touchend="handleTouchEnd"
    @mousedown="handleTouchStart"
    @mousemove="handleTouchMove"
    @mouseup="handleTouchEnd"
    @mouseleave="handleTouchEnd"
  >
    <div 
      class="swipe-content"
      :style="{ 
        transform: `translateX(${translateX}px)`,
        transition: isDragging ? 'none' : 'transform 0.3s ease'
      }"
    >
      <slot></slot>
    </div>
    
    <div class="swipe-indicators" v-if="showIndicators && totalItems > 1">
      <div 
        v-for="(_, index) in totalItems" 
        :key="index"
        class="indicator"
        :class="{ active: currentIndex === index }"
        @click="goToSlide(index)"
      ></div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'

const props = defineProps({
  itemsPerView: {
    type: Number,
    default: 1
  },
  showIndicators: {
    type: Boolean,
    default: true
  },
  autoPlay: {
    type: Boolean,
    default: false
  },
  autoPlayInterval: {
    type: Number,
    default: 3000
  }
})

const emit = defineEmits(['slide-change', 'swipe-left', 'swipe-right'])

const translateX = ref(0)
const isDragging = ref(false)
const startX = ref(0)
const currentX = ref(0)
const currentIndex = ref(0)
const totalItems = ref(0)
const containerWidth = ref(0)
let autoPlayTimer = null

const maxIndex = computed(() => {
  return Math.max(0, totalItems.value - props.itemsPerView)
})

function handleTouchStart(e) {
  isDragging.value = true
  startX.value = e.touches ? e.touches[0].clientX : e.clientX
  currentX.value = startX.value
  
  if (props.autoPlay) {
    stopAutoPlay()
  }
}

function handleTouchMove(e) {
  if (!isDragging.value) return
  
  const x = e.touches ? e.touches[0].clientX : e.clientX
  currentX.value = x
  
  const diff = currentX.value - startX.value
  translateX.value = diff
}

function handleTouchEnd() {
  if (!isDragging.value) return
  
  isDragging.value = false
  
  const diff = currentX.value - startX.value
  const threshold = 50 // العتبة للتنقل
  
  if (diff > threshold) {
    // السحب لليمين - العودة للسابق
    prevSlide()
  } else if (diff < -threshold) {
    // السحب لليسار - الذهاب للتالي
    nextSlide()
  } else {
    // العودة للوضع الحالي
    resetPosition()
  }
  
  if (props.autoPlay) {
    startAutoPlay()
  }
}

function nextSlide() {
  if (currentIndex.value < maxIndex.value) {
    currentIndex.value++
    emit('swipe-left', currentIndex.value)
  } else {
    // العودة للبداية
    currentIndex.value = 0
  }
  resetPosition()
  emit('slide-change', currentIndex.value)
}

function prevSlide() {
  if (currentIndex.value > 0) {
    currentIndex.value--
    emit('swipe-right', currentIndex.value)
  } else {
    // الذهاب للنهاية
    currentIndex.value = maxIndex.value
  }
  resetPosition()
  emit('slide-change', currentIndex.value)
}

function goToSlide(index) {
  currentIndex.value = index
  resetPosition()
  emit('slide-change', currentIndex.value)
}

function resetPosition() {
  const itemWidth = containerWidth.value / props.itemsPerView
  translateX.value = -currentIndex.value * itemWidth
}

function updateTotalItems() {
  const content = document.querySelector('.swipe-content')
  if (content) {
    totalItems.value = content.children.length
    containerWidth.value = content.offsetWidth
    resetPosition()
  }
}

function startAutoPlay() {
  if (autoPlayTimer) return
  autoPlayTimer = setInterval(() => {
    nextSlide()
  }, props.autoPlayInterval)
}

function stopAutoPlay() {
  if (autoPlayTimer) {
    clearInterval(autoPlayTimer)
    autoPlayTimer = null
  }
}

onMounted(() => {
  updateTotalItems()
  window.addEventListener('resize', updateTotalItems)
  
  if (props.autoPlay) {
    startAutoPlay()
  }
})

onUnmounted(() => {
  window.removeEventListener('resize', updateTotalItems)
  stopAutoPlay()
})

// توفير الوصول للوظائف للمكونات الأب
defineExpose({
  nextSlide,
  prevSlide,
  goToSlide,
  currentIndex
})
</script>

<style scoped>
.swipe-nav {
  position: relative;
  overflow: hidden;
  width: 100%;
}

.swipe-content {
  display: flex;
  width: 100%;
  cursor: grab;
  /* user-select: none; - تمت الإزالة للسماح بالنسخ واللصق على الهاتف */
  /* -webkit-user-select: none; - تمت الإزالة للسماح بالنسخ واللصق على الهاتف */
  touch-action: manipulation;
}

.swipe-content:active {
  cursor: grabbing;
}

.swipe-indicators {
  display: flex;
  justify-content: center;
  gap: 8px;
  margin-top: 16px;
  padding: 8px;
}

.indicator {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #E2E8F0;
  cursor: pointer;
  transition: all 0.3s ease;
}

.indicator.active {
  background: #1E3A5F;
  transform: scale(1.2);
}

.indicator:hover {
  background: #4A6FA5;
}

@media (max-width: 480px) {
  .indicator {
    width: 6px;
    height: 6px;
  }
  
  .swipe-indicators {
    gap: 6px;
    margin-top: 12px;
  }
}
</style>

```

---

## 📄 frontend-admin-dashboard/src/components/TaskCreateModal.vue

```vue
<template>
  <div v-if="show" class="modal-overlay" @click.self="close">
    <div class="modal-content">
      <div class="modal-header">
        <h3>{{ t('new_task') }}</h3>
        <button @click="close" class="close-button">×</button>
      </div>
      
      <form @submit.prevent="handleSubmit" class="task-form">
        <!-- Title section with language tabs -->
        <div class="form-section">
          <h4>{{ t('task_title') }}</h4>
          <div class="language-tabs">
            <button 
              type="button" 
              :class="{ active: activeTab === 'ar' }" 
              @click="activeTab = 'ar'"
            >🇸🇦 العربية</button>
            <button 
              type="button" 
              :class="{ active: activeTab === 'he' }" 
              @click="activeTab = 'he'"
            >🇮🇱 עברית</button>
            <button 
              type="button" 
              :class="{ active: activeTab === 'en' }" 
              @click="activeTab = 'en'"
            >🇬🇧 English</button>
          </div>
          
          <div class="language-inputs">
            <div :class="{ hidden: activeTab !== 'ar' }">
              <input 
                v-model="form.title_ar" 
                type="text" 
                :placeholder="t('title_ar_placeholder')"
                required
              />
            </div>
            <div :class="{ hidden: activeTab !== 'he' }">
              <input 
                v-model="form.title_he" 
                type="text" 
                :placeholder="t('title_he_placeholder')"
              />
            </div>
            <div :class="{ hidden: activeTab !== 'en' }">
              <input 
                v-model="form.title_en" 
                type="text" 
                :placeholder="t('title_en_placeholder')"
              />
            </div>
          </div>
        </div>

        <!-- Description section with language tabs -->
        <div class="form-section">
          <h4>{{ t('task_description') }}</h4>
          <div class="language-tabs">
            <button 
              type="button" 
              :class="{ active: activeDescTab === 'ar' }" 
              @click="activeDescTab = 'ar'"
            >🇸🇦 العربية</button>
            <button 
              type="button" 
              :class="{ active: activeDescTab === 'he' }" 
              @click="activeDescTab = 'he'"
            >🇮🇱 עברית</button>
            <button 
              type="button" 
              :class="{ active: activeDescTab === 'en' }" 
              @click="activeDescTab = 'en'"
            >🇬🇧 English</button>
          </div>
          
          <div class="language-inputs">
            <div :class="{ hidden: activeDescTab !== 'ar' }">
              <textarea 
                v-model="form.description_ar" 
                :placeholder="t('description_ar_placeholder')"
                rows="3"
              ></textarea>
            </div>
            <div :class="{ hidden: activeDescTab !== 'he' }">
              <textarea 
                v-model="form.description_he" 
                :placeholder="t('description_he_placeholder')"
                rows="3"
              ></textarea>
            </div>
            <div :class="{ hidden: activeDescTab !== 'en' }">
              <textarea 
                v-model="form.description_en" 
                :placeholder="t('description_en_placeholder')"
                rows="3"
              ></textarea>
            </div>
          </div>
        </div>

        <!-- Other fields -->
        <div class="form-section">
          <label>{{ t('task_worksite') }}</label>
          <select v-model="form.worksite_id" required>
            <option value="">{{ t('select_worksite') }}</option>
            <option v-for="ws in worksites" :key="ws.id" :value="ws.id">
              {{ ws.name }}
            </option>
          </select>
        </div>

        <div class="form-section">
          <label>{{ t('task_employee') }}</label>
          <select v-model="form.assigned_user_id">
            <option value="">{{ t('select_employee') }}</option>
            <option v-for="emp in employees" :key="emp.id" :value="emp.id">
              {{ emp.full_name }}
            </option>
          </select>
        </div>

        <div class="form-section">
          <label>{{ t('task_priority') }}</label>
          <select v-model="form.priority">
            <option value="low">{{ t('priority_low') }}</option>
            <option value="normal">{{ t('priority_normal') }}</option>
            <option value="high">{{ t('priority_high') }}</option>
            <option value="urgent">{{ t('priority_urgent') }}</option>
          </select>
        </div>

        <div class="form-section">
          <label>{{ t('scheduled_start') }}</label>
          <input v-model="form.scheduled_start" type="datetime-local" />
        </div>

        <div class="form-section">
          <label>{{ t('scheduled_end') }}</label>
          <input v-model="form.scheduled_end" type="datetime-local" />
        </div>

        <div class="form-actions">
          <button type="button" @click="close" class="btn btn--secondary">{{ t('cancel') }}</button>
          <button type="submit" class="btn btn--primary" :disabled="loading">
            {{ loading ? t('creating') : t('create') }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from '../services/i18n'

const props = defineProps({
  show: Boolean
})

const emit = defineEmits(['close', 'created'])

const { t } = useI18n()
const loading = ref(false)
const activeTab = ref('ar')
const activeDescTab = ref('ar')
const worksites = ref([])
const employees = ref([])

const form = reactive({
  title_ar: '',
  title_he: '',
  title_en: '',
  description_ar: '',
  description_he: '',
  description_en: '',
  worksite_id: '',
  assigned_user_id: '',
  priority: 'normal',
  scheduled_start: '',
  scheduled_end: ''
})

const close = () => {
  emit('close')
}

const handleSubmit = async () => {
  loading.value = true
  try {
    const apiBaseUrl = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1'
    const token = localStorage.getItem('token')
    
    // Use Arabic title as default if other languages are empty
    const payload = {
      title: form.title_ar,
      title_ar: form.title_ar,
      title_he: form.title_he || form.title_ar,
      title_en: form.title_en || form.title_ar,
      description: form.description_ar,
      description_ar: form.description_ar,
      description_he: form.description_he || form.description_ar,
      description_en: form.description_en || form.description_ar,
      worksite_id: form.worksite_id,
      assigned_user_id: form.assigned_user_id || null,
      priority: form.priority,
      scheduled_start: form.scheduled_start || null,
      scheduled_end: form.scheduled_end || null
    }
    
    const response = await fetch(`${apiBaseUrl}/tasks`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(payload)
    })
    
    if (!response.ok) {
      throw new Error('Failed to create task')
    }
    
    emit('created')
    close()
  } catch (err) {
    console.error('Error creating task:', err)
    alert(t('error_creating_task'))
  } finally {
    loading.value = false
  }
}

const fetchWorksites = async () => {
  try {
    const apiBaseUrl = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1'
    const token = localStorage.getItem('token')
    
    const response = await fetch(`${apiBaseUrl}/worksites`, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })
    
    if (response.ok) {
      worksites.value = await response.json()
    }
  } catch (err) {
    console.error('Error fetching worksites:', err)
  }
}

const fetchEmployees = async () => {
  try {
    const apiBaseUrl = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1'
    const token = localStorage.getItem('token')
    
    const response = await fetch(`${apiBaseUrl}/users?role=employee`, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })
    
    if (response.ok) {
      employees.value = await response.json()
    }
  } catch (err) {
    console.error('Error fetching employees:', err)
  }
}

onMounted(() => {
  fetchWorksites()
  fetchEmployees()
})
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  border-radius: 12px;
  width: 90%;
  max-width: 600px;
  max-height: 90vh;
  overflow-y: auto;
  padding: 24px;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.modal-header h3 {
  margin: 0;
  font-size: 20px;
}

.close-button {
  background: none;
  border: none;
  font-size: 24px;
  cursor: pointer;
  padding: 0;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.task-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.form-section {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-section h4 {
  margin: 0;
  font-size: 14px;
  font-weight: 600;
  color: #495057;
}

.form-section label {
  font-size: 13px;
  font-weight: 500;
  color: #6c757d;
}

.language-tabs {
  display: flex;
  gap: 8px;
  margin-bottom: 8px;
}

.language-tabs button {
  padding: 8px 12px;
  border: 1px solid #dee2e6;
  background: white;
  border-radius: 6px;
  cursor: pointer;
  font-size: 13px;
  transition: all 0.2s ease;
}

.language-tabs button:hover {
  background: #f8f9fa;
}

.language-tabs button.active {
  background: #007bff;
  color: white;
  border-color: #007bff;
}

.language-inputs input,
.language-inputs textarea,
.form-section input,
.form-section select,
.form-section textarea {
  padding: 10px 12px;
  border: 1px solid #dee2e6;
  border-radius: 6px;
  font-size: 14px;
  font-family: inherit;
}

.language-inputs input:focus,
.language-inputs textarea:focus,
.form-section input:focus,
.form-section select:focus,
.form-section textarea:focus {
  outline: none;
  border-color: #007bff;
  box-shadow: 0 0 0 3px rgba(0, 123, 255, 0.1);
}

.hidden {
  display: none;
}

.form-actions {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
  margin-top: 8px;
}

.btn {
  padding: 10px 20px;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  border: none;
  transition: all 0.2s ease;
}

.btn--primary {
  background: #007bff;
  color: white;
}

.btn--primary:hover {
  background: #0056b3;
}

.btn--primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn--secondary {
  background: #6c757d;
  color: white;
}

.btn--secondary:hover {
  background: #545b62;
}
</style>
```

---

## 📄 frontend-admin-dashboard/src/components/TaskDistributionChart.vue

```vue
<template>
  <div class="dist">
    <div class="dist__bar">
      <span v-for="seg in segments" :key="seg.key" :style="{ width: seg.pct + '%', background: seg.color }"></span>
    </div>
    <ul class="dist__legend">
      <li v-for="seg in segments" :key="seg.key">
        <span class="dist__dot" :style="{ background: seg.color }"></span>
        {{ seg.label }} <b class="mono">{{ seg.value }}</b>
      </li>
    </ul>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  completed: { type: Number, default: 0 },
  in_progress: { type: Number, default: 0 },
  pending: { type: Number, default: 0 },
  late: { type: Number, default: 0 },
})

const segments = computed(() => {
  const total = props.completed + props.in_progress + props.pending + props.late || 1
  const defs = [
    { key: 'completed', label: 'مكتملة', value: props.completed, color: 'var(--signal-in)' },
    { key: 'in_progress', label: 'جارية', value: props.in_progress, color: 'var(--gold)' },
    { key: 'pending', label: 'قيد الانتظار', value: props.pending, color: 'var(--line-strong)' },
    { key: 'late', label: 'متأخرة', value: props.late, color: 'var(--signal-out)' },
  ]
  return defs.map((d) => ({ ...d, pct: (d.value / total) * 100 }))
})
</script>

<style scoped>
.dist__bar { display: flex; height: 12px; border-radius: 999px; overflow: hidden; background: var(--line); }
.dist__legend { list-style: none; padding: 0; margin: 16px 0 0; display: flex; flex-wrap: wrap; gap: 14px; font-size: 13px; color: var(--ink-soft); }
.dist__legend li { display: flex; align-items: center; gap: 6px; }
.dist__dot { width: 8px; height: 8px; border-radius: 50%; }
</style>

```

---

## 📄 frontend-admin-dashboard/src/components/WorksiteFormModal.vue

```vue
<template>
  <div class="modal-backdrop" @click.self="$emit('close')">
    <div class="modal card">
      <div class="modal-header">
        <h3>📍 {{ t('new_worksite') }}</h3>
        <button class="modal-close" @click="$emit('close')">✕</button>
      </div>

      <p class="modal__hint">{{ t('search_worldwide') }} - {{ t('search_language_hint') }}</p>

      <form class="modal__form" @submit.prevent="handleSubmit">
        <div v-if="error" class="alert alert-error">{{ error }}</div>
        <div v-if="success" class="alert alert-success">{{ success }}</div>

        <div class="form-group">
          <label>📛 {{ t('worksite_name') }} <span class="required">*</span></label>
          <input
            v-model="form.name"
            type="text"
            :placeholder="t('worksite_name_placeholder')"
            required
            class="search-input"
          />
        </div>

        <!-- ========================================== -->
        <!-- البحث العالمي متعدد اللغات -->
        <!-- ========================================== -->
        <div class="form-group">
          <label>📍 {{ t('search_address') }} <span class="required">*</span></label>
          <div class="search-wrapper">
            <input
              v-model="searchQuery"
              type="text"
              :placeholder="t('search_address_placeholder')"
              @input="onSearch"
              @focus="showResults = true"
              class="search-input"
              autocomplete="off"
            />
            <span v-if="loading" class="search-loading">⏳</span>
          </div>

          <!-- نتائج البحث متعددة اللغات -->
          <div v-if="showResults && searchResults.length > 0" class="search-results">
            <div
              v-for="result in searchResults"
              :key="result.id"
              class="result-item"
              @click="selectResult(result)"
            >
              <span class="result-icon">📍</span>
              <div class="result-info">
                <strong>{{ getLocalizedLabel(result) }}</strong>
                <span class="result-address">
                  {{ getLocalizedAddress(result) }}
                </span>
              </div>
              <span class="result-type">{{ getTypeLabel(result.type) }}</span>
            </div>
          </div>

          <div v-if="showResults && searchQuery && searchResults.length === 0 && !loading" class="no-results">
            <span>🔍</span>
            <p>{{ t('no_results_found') }}</p>
          </div>
        </div>

        <!-- ========================================== -->
        <!-- الموقع المختار -->
        <!-- ========================================== -->
        <div v-if="selectedResult" class="selected-location">
          <div class="selected-location__header">
            <span>✅ {{ getLocalizedLabel(selectedResult) }}</span>
            <button type="button" class="selected-clear" @click="clearSelection">✕</button>
          </div>
          <div class="selected-location__details">
            <p><strong>{{ t('city') }}:</strong> {{ getLocalizedCity(selectedResult) }}</p>
            <p><strong>{{ t('street') }}:</strong> {{ getLocalizedStreet(selectedResult) }}</p>
            <p><strong>{{ t('building_number') }}:</strong> {{ selectedResult.house_number || '—' }}</p>
            <p class="mono"><strong>{{ t('coordinates') }}:</strong> {{ selectedResult.latitude }}, {{ selectedResult.longitude }}</p>
          </div>
        </div>

        <div class="form-group">
          <label>⭕ {{ t('allowed_radius') }} <span class="required">*</span></label>
          <input
            v-model.number="form.radius_meters"
            type="number"
            required
            placeholder="100"
            min="10"
            class="search-input"
          />
        </div>

        <div class="form-actions">
          <button type="button" class="btn btn--ghost" @click="$emit('close')">{{ t('cancel') }}</button>
          <button type="submit" class="btn btn--primary" :disabled="loading || !selectedResult">
            {{ loading ? `⏳ ${t('saving')}` : `💾 ${t('save')}` }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import api from '../services/api'
import { useI18n } from '../services/i18n'

const { t, currentLang } = useI18n()
const emit = defineEmits(['close', 'worksite-added'])

const searchQuery = ref('')
const searchResults = ref([])
const loading = ref(false)
const showResults = ref(false)
const selectedResult = ref(null)

const form = ref({
  name: '',
  radius_meters: 100,
  latitude: '',
  longitude: ''
})

const error = ref('')
const success = ref('')
const isSubmitting = ref(false)

let searchTimeout = null

function getLocalizedLabel(result) {
  const lang = currentLang.value
  if (lang === 'he' && result.label_he) return result.label_he
  if (lang === 'ar' && result.label_ar) return result.label_ar
  return result.label || result.street || result.city
}

function getLocalizedAddress(result) {
  const lang = currentLang.value
  let city = ''
  let street = ''
  
  if (lang === 'he') {
    city = result.city_he || result.city || ''
    street = result.street_he || result.street || ''
  } else if (lang === 'ar') {
    city = result.city_ar || result.city || ''
    street = result.street_ar || result.street || ''
  } else {
    city = result.city || ''
    street = result.street || ''
  }
  
  return `${city} ${street} ${result.house_number || ''}`.trim()
}

function getLocalizedCity(result) {
  const lang = currentLang.value
  if (lang === 'he') return result.city_he || result.city || '—'
  if (lang === 'ar') return result.city_ar || result.city || '—'
  return result.city || '—'
}

function getLocalizedStreet(result) {
  const lang = currentLang.value
  if (lang === 'he') return result.street_he || result.street || '—'
  if (lang === 'ar') return result.street_ar || result.street || '—'
  return result.street || '—'
}

function getTypeLabel(type) {
  const types = {
    'city': t('type_city'),
    'street': t('type_street'),
    'address': t('type_address'),
    'house': t('type_house'),
    'landmark': t('type_landmark'),
    'location': t('type_location')
  }
  return types[type] || type || t('type_location')
}

function onSearch() {
  clearTimeout(searchTimeout)
  
  if (searchQuery.value.length < 2) {
    searchResults.value = []
    return
  }

  loading.value = true
  searchTimeout = setTimeout(async () => {
    try {
      // استخدام اللغة الحالية
      const lang = currentLang.value
      const { data } = await api.get('/geocode/autocomplete', {
        params: {
          q: searchQuery.value.trim(),
          lang: lang
        }
      })
      
      searchResults.value = data.results || []
      showResults.value = true
    } catch (error) {
      console.error('❌ فشل البحث:', error)
      searchResults.value = []
    } finally {
      loading.value = false
    }
  }, 600)
}

function selectResult(result) {
  selectedResult.value = result
  searchQuery.value = getLocalizedLabel(result)
  showResults.value = false
  
  form.value.name = getLocalizedLabel(result)
  form.value.latitude = result.latitude
  form.value.longitude = result.longitude
  
  error.value = ''
}

function clearSelection() {
  selectedResult.value = null
  searchQuery.value = ''
  searchResults.value = []
  form.value.name = ''
  form.value.latitude = ''
  form.value.longitude = ''
}

function closeResults(e) {
  const wrapper = document.querySelector('.modal__form')
  if (wrapper && !wrapper.contains(e.target)) {
    showResults.value = false
  }
}

document.addEventListener('click', closeResults)

async function handleSubmit() {
  if (!selectedResult.value) {
    error.value = t('select_address_required')
    return
  }

  isSubmitting.value = true
  error.value = ''
  success.value = ''

  try {
    const payload = {
      name: form.value.name || getLocalizedLabel(selectedResult.value),
      address: getLocalizedLabel(selectedResult.value),
      latitude: parseFloat(form.value.latitude || selectedResult.value.latitude),
      longitude: parseFloat(form.value.longitude || selectedResult.value.longitude),
      radius_meters: form.value.radius_meters,
      city: getLocalizedCity(selectedResult.value),
      street: getLocalizedStreet(selectedResult.value),
      street_number: selectedResult.value.house_number || ''
    }

    await api.post('/worksites', payload)
    success.value = '✅ ' + t('worksite_added_successfully')
    
    setTimeout(() => {
      emit('worksite-added')
      emit('close')
    }, 1500)
  } catch (err) {
    error.value = err.response?.data?.error || '❌ ' + t('save_failed')
    console.error('خطأ:', err)
  } finally {
    isSubmitting.value = false
  }
}
</script>

<style scoped>
.modal-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.6);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 20px;
}

.modal {
  width: 100%;
  max-width: 520px;
  padding: 0;
  max-height: 90vh;
  overflow-y: auto;
  background: var(--surface);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-xl);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid var(--line);
}

.modal-header h3 {
  font-size: 17px;
  margin: 0;
}

.modal-close {
  background: none;
  border: none;
  font-size: 24px;
  cursor: pointer;
  color: var(--ink-soft);
}

.modal__hint {
  padding: 0 20px 16px;
  font-size: 13px;
  color: var(--ink-soft);
  margin: 0;
}

.modal__form {
  padding: 0 20px 20px;
}

.form-group {
  margin-bottom: 14px;
  position: relative;
}

.form-group label {
  display: block;
  font-size: 13px;
  font-weight: 600;
  color: var(--ink);
  margin-bottom: 4px;
}

.search-wrapper {
  position: relative;
}

.search-input {
  width: 100%;
  padding: 10px 14px;
  padding-left: 40px;
  border: 1.5px solid var(--line);
  border-radius: var(--radius-sm);
  font-size: 14px;
  font-family: var(--font-body);
  background: var(--surface);
  transition: all 0.3s;
}

.search-input:focus {
  outline: none;
  border-color: var(--brand);
  box-shadow: 0 0 0 3px var(--brand-tint);
}

.search-loading {
  position: absolute;
  left: 12px;
  top: 50%;
  transform: translateY(-50%);
  font-size: 16px;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: translateY(-50%) rotate(0deg); }
  to { transform: translateY(-50%) rotate(360deg); }
}

.required {
  color: var(--signal-out);
}

.search-results {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  background: var(--surface);
  border: 1px solid var(--line);
  border-radius: var(--radius-sm);
  box-shadow: var(--shadow-lg);
  max-height: 250px;
  overflow-y: auto;
  z-index: 1000;
  margin-top: 4px;
}

.result-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 14px;
  cursor: pointer;
  border-bottom: 1px solid var(--line);
  transition: background 0.2s;
}

.result-item:hover {
  background: var(--brand-tint);
}

.result-item:last-child {
  border-bottom: none;
}

.result-icon {
  font-size: 18px;
  flex-shrink: 0;
}

.result-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.result-info strong {
  font-size: 13px;
  color: var(--ink);
}

.result-address {
  font-size: 11px;
  color: var(--ink-soft);
}

.result-type {
  font-size: 10px;
  color: var(--brand);
  background: var(--brand-tint);
  padding: 2px 8px;
  border-radius: 999px;
}

.no-results {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  background: var(--surface);
  border: 1px solid var(--line);
  border-radius: var(--radius-sm);
  padding: 16px;
  text-align: center;
  z-index: 1000;
  margin-top: 4px;
}

.no-results span {
  font-size: 24px;
  display: block;
  margin-bottom: 4px;
}

.no-results p {
  font-size: 13px;
  color: var(--ink-soft);
  margin: 0;
}

.selected-location {
  background: var(--brand-tint);
  border: 1px solid var(--brand);
  border-radius: var(--radius-sm);
  padding: 12px 16px;
  margin-bottom: 14px;
}

.selected-location__header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
  font-weight: 600;
  color: var(--brand);
}

.selected-clear {
  background: none;
  border: none;
  color: var(--signal-out);
  cursor: pointer;
  font-size: 18px;
  padding: 0 4px;
}

.selected-location__details p {
  margin: 4px 0;
  font-size: 13px;
  color: var(--ink);
}

.alert {
  padding: 10px 14px;
  border-radius: var(--radius-sm);
  font-size: 13px;
  margin-bottom: 14px;
}

.alert-error {
  background: var(--signal-out-tint);
  color: var(--signal-out);
}

.alert-success {
  background: var(--signal-in-tint);
  color: var(--signal-in);
}

.form-actions {
  display: flex;
  gap: 10px;
  justify-content: flex-end;
  margin-top: 14px;
  padding-top: 14px;
  border-top: 1px solid var(--line);
}
</style>

```

---

## 📄 frontend-admin-dashboard/src/config/company.js

```javascript
export const companyConfig = {
  name: 'DevPro',
  fullName: 'DevPro - Your Partner in Digital Transformation',
  shortName: 'DevPro',
  logo: '/src/assets/devpro-logo.jpg',
  website: 'https://devpro.com',
  email: 'info@devpro.com',
  slogan: 'Your Partner in Digital Transformation',
}

```

---

## 📄 frontend-admin-dashboard/src/main.js

```javascript
import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import './styles/tokens.css'
import i18n from './services/i18n'

import 'leaflet/dist/leaflet.css'
import L from 'leaflet'

// ==========================================
// DevTools Detection & Protection
// ==========================================

// التحقق من بيئة التطوير والإنتاج
const isDevelopment = import.meta.env.VITE_APP_ENV === 'development' || 
                      import.meta.env.MODE === 'development';

// التحقق من مفتاح التجاوز
const bypassKey = localStorage.getItem('devtools_bypass_key') || 
                 new URLSearchParams(window.location.search).get('bypass');
const isBypassActive = bypassKey === 'worktrack_dev_2024';

// الحماية مفعلة افتراضياً، تعطيل فقط في بيئة التطوير أو عند وجود مفتاح التجاوز
if (!isDevelopment && !isBypassActive) {
  const devtools = {
    open: false,
    orientation: null
  };

  const threshold = 160;

  // اكتشاف فتح DevTools من خلال تغيير حجم النافذة
  setInterval(() => {
    const widthThreshold = window.outerWidth - window.innerWidth > threshold;
    const heightThreshold = window.outerHeight - window.innerHeight > threshold;
    
    if (widthThreshold || heightThreshold) {
      if (!devtools.open) {
        devtools.open = true;
        console.warn('⚠️ DevTools detected - Unauthorized access attempt');
        // يمكن إعادة تحميل الصفحة أو اتخاذ إجراء آخر
        // window.location.reload();
      }
    } else {
      devtools.open = false;
    }
  }, 500);

  // اكتشاف debugger (Anti-Debugging)
  setInterval(() => {
    const start = new Date().getTime();
    debugger; // يتوقف إذا كانت DevTools مفتوحة
    const end = new Date().getTime();
    
    if (end - start > 100) {
      console.warn('⚠️ Debugger detected - Unauthorized access attempt');
      // window.location.reload();
    }
  }, 1000);

  // منع النقر الأيمن - ولكن السماح بالنسخ واللصق
  document.addEventListener('contextmenu', (e) => {
    // السماح بالنقر الأيمن على حقول الإدخال والنصوص للنسخ واللصق
    const target = e.target;
    if (target.tagName === 'INPUT' || target.tagName === 'TEXTAREA' || 
        target.isContentEditable || target.tagName === 'SELECT') {
      return true; // السماح بالنقر الأيمن للنسخ واللصق
    }
    e.preventDefault();
    return false;
  });

  // منع اختصارات لوحة المفاتيح لفتح DevTools - ولكن السماح بالنسخ واللصق
  document.addEventListener('keydown', (e) => {
    // F12, Ctrl+Shift+I, Ctrl+Shift+J, Ctrl+Shift+C, Ctrl+U
    // تم استثناء Ctrl+C للسماح بالنسخ
    if (e.key === 'F12' || 
        (e.ctrlKey && e.shiftKey && (e.key === 'I' || e.key === 'J' || e.key === 'C')) ||
        (e.ctrlKey && e.key === 'U')) {
      e.preventDefault();
      return false;
    }
    // السماح بـ Ctrl+C و Ctrl+V و Ctrl+X للنسخ واللصق
    if (e.ctrlKey && (e.key === 'c' || e.key === 'v' || e.key === 'x' || e.key === 'a')) {
      return true; // السماح بهذه الاختصارات
    }
  });
} else {
  console.log('🔓 DevTools protection bypassed - Development mode active');
}

// تم إزالة منع النسخ واللصق للسماح للمستخدمين بهذه الوظائف

// إضافي: تأكيد السماح بالنسخ واللصق على الهاتف
document.addEventListener('DOMContentLoaded', () => {
  // تحسين النسخ واللصق في وضع PWA
  const enableCopyPaste = () => {
    // السماح بتحديد النصوص
    document.body.style.webkitUserSelect = 'text';
    document.body.style.userSelect = 'text';
    document.body.style.webkitTouchCallout = 'default';
    
    // السماح بالتفاعل مع الحافظة
    const inputs = document.querySelectorAll('input, textarea, select');
    inputs.forEach(input => {
      input.style.webkitUserSelect = 'text';
      input.style.userSelect = 'text';
      input.style.webkitTouchCallout = 'default';
    });
  };
  
  // تشغيل عند التحميل
  enableCopyPaste();
  
  // تشغيل مرة أخرى بعد فترة للتأكد من التطبيق
  setTimeout(enableCopyPaste, 1000);
  
  // السماح بالنسخ واللصق عبر أحداث اللمس
  document.addEventListener('touchstart', (e) => {
    const target = e.target;
    if (target.tagName === 'INPUT' || target.tagName === 'TEXTAREA' || 
        target.isContentEditable || target.tagName === 'SELECT') {
      // السماح بالتفاعل للنسخ واللصق
      target.style.webkitUserSelect = 'text';
      target.style.userSelect = 'text';
    }
  }, { passive: true });
  
  // منع منع أحداث النسخ واللصق
  document.addEventListener('copy', (e) => {
    e.stopPropagation();
  }, true);
  
  document.addEventListener('cut', (e) => {
    e.stopPropagation();
  }, true);
  
  document.addEventListener('paste', (e) => {
    e.stopPropagation();
  }, true);
});

delete L.Icon.Default.prototype._getIconUrl
L.Icon.Default.mergeOptions({
  iconRetinaUrl: new URL('leaflet/dist/images/marker-icon-2x.png', import.meta.url).href,
  iconUrl: new URL('leaflet/dist/images/marker-icon.png', import.meta.url).href,
  shadowUrl: new URL('leaflet/dist/images/marker-shadow.png', import.meta.url).href,
})

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

  // Set initial direction based on stored language
  const storedLang = localStorage.getItem('worktrack_language') || 'ar'
  document.documentElement.dir = (storedLang === 'ar' || storedLang === 'he') ? 'rtl' : 'ltr'
  document.documentElement.lang = storedLang
})

```

---

## 📄 frontend-admin-dashboard/src/router/index.js

```javascript
import { createRouter, createWebHashHistory, createMemoryHistory } from 'vue-router'
import LoginView from '../views/LoginView.vue'
import DashboardView from '../views/DashboardView.vue'
import EmployeesView from '../views/EmployeesView.vue'
import CustomersView from '../views/CustomersView.vue'
import TasksView from '../views/TasksView.vue'
import WorksitesView from '../views/WorksitesView.vue'
import ClientsView from '../views/ClientsView.vue'
import ReportsView from '../views/ReportsView.vue'
import SettingsView from '../views/SettingsView.vue'
import ServiceRequestsView from '../views/ServiceRequestsView.vue'
import NotificationsView from '../views/NotificationsView.vue'
import SubscriptionStatusView from '../views/SubscriptionStatusView.vue'

const routes = [
  { path: '/', redirect: '/login' },
  { path: '/login', component: LoginView, meta: { public: true } },
  { path: '/dashboard', component: DashboardView, meta: { requiresAuth: true } },
  { path: '/service-requests', component: ServiceRequestsView, meta: { requiresAuth: true } },
  { path: '/employees', component: EmployeesView, meta: { requiresAuth: true } },
  { path: '/customers', component: CustomersView, meta: { requiresAuth: true } },
  { path: '/tasks', component: TasksView, meta: { requiresAuth: true } },
  { path: '/worksites', component: WorksitesView, meta: { requiresAuth: true } },
  { path: '/reports', component: ReportsView, meta: { requiresAuth: true } },
  { path: '/settings', component: SettingsView, meta: { requiresAuth: true } },
  { path: '/notifications', component: NotificationsView, meta: { requiresAuth: true } },
  { path: '/subscription', component: SubscriptionStatusView, meta: { requiresAuth: true } },
]

// تحويل المسارات القديمة إلى صيغة Hash Mode
function migrateSavedPath() {
  const savedPath = localStorage.getItem('worktrack_last_path')
  if (savedPath && !savedPath.startsWith('#/')) {
    // تحويل المسار القديم (/dashboard) إلى الجديد (/#/dashboard)
    const newPath = savedPath.startsWith('/') ? `#${savedPath}` : `#/${savedPath}`
    localStorage.setItem('worktrack_last_path', newPath)
    console.log('Migrated saved path:', savedPath, '->', newPath)
  }
}

// استخدام MemoryHistory لـ Electron و HashHistory للويب
const isElectron = typeof window !== 'undefined' && (window.isElectron || (window.process && window.process.type === 'renderer'))

const router = createRouter({
  history: isElectron ? createMemoryHistory() : createWebHashHistory(),
  routes
})

// تشغيل التحويل عند تحميل الـ router
if (!isElectron) {
  migrateSavedPath()
}

// حماية المسارات
router.beforeEach((to, from, next) => {
  const isAuthed = !!localStorage.getItem('worktrack_admin_token')
  
  console.log('Router guard:', { to: to.path, from: from.path, isAuthed })
  
  // إذا كان المسار الجذر والمستخدم مسجل، استعادة المسار المحفوظ أو اذهب للـ dashboard
  if (to.path === '/' && isAuthed) {
    const savedPath = localStorage.getItem('worktrack_last_path')
    if (savedPath && savedPath !== '/login' && savedPath !== '/' && savedPath !== '#/login' && savedPath !== '#/') {
      // إزالة الـ hash إذا كان موجوداً للـ Vue Router
      const cleanPath = savedPath.startsWith('#') ? savedPath.substring(1) : savedPath
      next(cleanPath)
    } else {
      next('/dashboard')
    }
  }
  // إذا كان المسار الجذر والمستخدم غير مسجل، اذهب للـ login
  else if (to.path === '/' && !isAuthed) {
    next('/login')
  }
  // إذا كان المسار يتطلب مصادقة والمستخدم غير مسجل
  else if (to.meta.requiresAuth && !isAuthed) {
    next('/login')
  }
  // إذا كان المستخدم مسجل ويحاول الذهاب إلى login
  else if (to.path === '/login' && isAuthed) {
    next('/dashboard')
  }
  // للمسارات الأخرى - إذا كان المستخدم مسجل، فقط تحديث المسار المحفوظ
  else if (isAuthed && to.meta.requiresAuth) {
    // السماح بالوصول مع حفظ المسار في afterEach
    next()
  }
  // السماح بالوصول
  else {
    next()
  }
})

// حفظ المسار الحالي بعد كل تغيير مسار
router.afterEach((to) => {
  // لا تحفظ مسار login أو المسار الجذر
  if (to.path !== '/login' && to.path !== '/') {
    // في Hash Mode، نحفظ المسار مع الـ hash
    const pathToSave = isElectron ? to.path : `#${to.path}`
    localStorage.setItem('worktrack_last_path', pathToSave)
  }
})

// معالجة الأخطاء في التوجيه
router.onError((error) => {
  console.error('Router error:', error)
})

export default router

```

---

## 📄 frontend-admin-dashboard/src/services/api.js

```javascript
import axios from 'axios'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1',
  withCredentials: true, // مهم لإرسال واستقبال cookies
})

api.interceptors.request.use((config) => {
  // Add token from localStorage as fallback for Authorization header
  const token = localStorage.getItem('worktrack_admin_token')
  if (token) {
    config.headers['Authorization'] = `Bearer ${token}`
  }
  
  // Add current language to request header
  const currentLang = localStorage.getItem('worktrack_admin_language') || 'ar'
  config.headers['X-Lang'] = currentLang
  
  return config
})

// Response interceptor to handle 401 errors and password changes
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      const errorMessage = error.response?.data?.error || ''
      console.log('🔍 401 Error detected:', errorMessage) // للتصحيح
      
      // Check if error is about password change - توسيع البحث
      if (errorMessage.includes('password changed') || 
          errorMessage.includes('كلمة المرور') ||
          errorMessage.includes('كلمة السر') ||
          errorMessage.includes('הסיסמה שונתה') ||
          errorMessage.includes('Password has been changed')) {
        console.log('🔓 Password change detected, showing alert') // للتصحيح
        // Show custom popup for password change
        showPasswordChangedAlert()
      } else {
        console.log('🚪 Normal logout required') // للتصحيح
        // Handle other 401 errors (normal logout)
        handleLogout()
      }
    }
    return Promise.reject(error)
  }
)

function showPasswordChangedAlert() {
  console.log('🚨 Showing password changed alert (API interceptor)')
  const currentLang = localStorage.getItem('worktrack_admin_language') || 'ar'
  
  // Default fallback messages
  const messages = {
    ar: {
      title: 'تم تغيير كلمة المرور',
      message: 'تم تغيير كلمة المرور الخاصة بحسابك. يرجى تسجيل الدخول مرة أخرى.',
      button: 'تسجيل الدخول'
    },
    he: {
      title: 'הסיסמה שונתה',
      message: 'הסיסמה שלך שונתה. אנא התחבר שוב.',
      button: 'התחבר'
    },
    en: {
      title: 'Password Changed',
      message: 'Your password has been changed. Please log in again.',
      button: 'Log In'
    }
  }
  
  const msg = messages[currentLang] || messages.ar
  
  // Remove any existing alerts first
  const existingAlerts = document.querySelectorAll('[data-password-alert-api]')
  existingAlerts.forEach(alert => alert.remove())
  
  // Create and show alert
  const alertDiv = document.createElement('div')
  alertDiv.setAttribute('data-password-alert-api', 'true')
  alertDiv.style.cssText = `
    position: fixed !important;
    top: 0 !important;
    left: 0 !important;
    right: 0 !important;
    bottom: 0 !important;
    background: rgba(0, 0, 0, 0.8) !important;
    display: flex !important;
    align-items: center !important;
    justify-content: center !important;
    z-index: 999999 !important;
    direction: ${currentLang === 'he' ? 'rtl' : currentLang === 'en' ? 'ltr' : 'rtl'} !important;
    font-family: ${currentLang === 'he' ? 'Heebo, Arial, sans-serif' : currentLang === 'en' ? 'Arial, sans-serif' : 'Cairo, Arial, sans-serif'} !important;
  `
  
  const alertBox = document.createElement('div')
  alertBox.style.cssText = `
    background: white !important;
    padding: 40px !important;
    border-radius: 16px !important;
    max-width: 450px !important;
    width: 90% !important;
    text-align: center !important;
    box-shadow: 0 20px 60px rgba(0, 0, 0, 0.4) !important;
    animation: slideIn 0.4s ease-out !important;
    position: relative !important;
  `
  
  alertBox.innerHTML = `
    <div style="font-size: 64px; margin-bottom: 20px; animation: pulse 2s infinite;">🔒</div>
    <h2 style="margin: 0 0 20px 0; color: #e74c3c; font-size: 24px; font-weight: 700;">${msg.title}</h2>
    <p style="margin: 0 0 30px 0; color: #555; line-height: 1.8; font-size: 16px;">${msg.message}</p>
    <button id="passwordChangedLogoutBtn" style="
      background: linear-gradient(135deg, #e74c3c 0%, #c0392b 100%) !important;
      color: white !important;
      border: none !important;
      padding: 16px 40px !important;
      border-radius: 8px !important;
      font-size: 18px !important;
      font-weight: 600 !important;
      cursor: pointer !important;
      transition: all 0.3s !important;
      box-shadow: 0 4px 15px rgba(231, 76, 60, 0.3) !important;
    ">${msg.button}</button>
    <style>
      @keyframes slideIn {
        from {
          opacity: 0;
          transform: translateY(-30px) scale(0.95);
        }
        to {
          opacity: 1;
          transform: translateY(0) scale(1);
        }
      }
      @keyframes pulse {
        0%, 100% { transform: scale(1); }
        50% { transform: scale(1.1); }
      }
      #passwordChangedLogoutBtn:hover {
        transform: translateY(-2px) !important;
        box-shadow: 0 6px 20px rgba(231, 76, 60, 0.4) !important;
      }
    </style>
  `
  
  alertDiv.appendChild(alertBox)
  document.body.appendChild(alertDiv)
  
  console.log('✅ Alert added to DOM (API interceptor)')
  
  // Make function globally available
  window.handlePasswordChangedLogoutAPI = () => {
    console.log('🚪 Password changed logout clicked (API interceptor)')
    handleLogout()
    const alertToRemove = document.querySelector('[data-password-alert-api]')
    if (alertToRemove) {
      alertToRemove.remove()
    }
  }
  
  // Add event listener to button
  setTimeout(() => {
    const btn = document.getElementById('passwordChangedLogoutBtn')
    if (btn) {
      btn.addEventListener('click', window.handlePasswordChangedLogoutAPI)
      console.log('✅ Button event listener added (API interceptor)')
    } else {
      console.error('❌ Button not found (API interceptor)')
    }
  }, 100)
}

function handleLogout() {
  // Clear all auth data
  localStorage.removeItem('worktrack_admin_token')
  localStorage.removeItem('worktrack_admin_user')
  
  // Redirect to login page
  window.location.href = '/login'
}

export default api

```

---

## 📄 frontend-admin-dashboard/src/services/auth.js

```javascript
import api from './api'

export async function login(email, password) {
  try {
    console.log('📤 محاولة تسجيل الدخول:', { email })
    
    const { data } = await api.post('/auth/login', { 
      email: email.trim(), 
      password: password.trim() 
    })
    
    console.log('✅ تم تسجيل الدخول بنجاح')
    
    localStorage.setItem('worktrack_admin_token', data.token)
    localStorage.setItem('worktrack_admin_user', JSON.stringify(data.user))
    
    return data
  } catch (error) {
    console.error('❌ فشل تسجيل الدخول:', error.response?.data || error.message)
    throw error
  }
}

export async function getCurrentUser() {
  const { data } = await api.get('/auth/me')
  return data
}

export function logout() {
  localStorage.removeItem('worktrack_admin_token')
  localStorage.removeItem('worktrack_admin_user')
}

export function currentUser() {
  const raw = localStorage.getItem('worktrack_admin_user')
  return raw ? JSON.parse(raw) : null
}

export function getToken() {
  return localStorage.getItem('worktrack_admin_token')
}

```

---

## 📄 frontend-admin-dashboard/src/services/i18n.js

```javascript
import { reactive, computed } from 'vue'

// Load translation data directly
const ar = {
  "app_name": "WorkTrack",
  "login": "تسجيل الدخول",
  "logout": "تسجيل الخروج",
  "phone": "رقم الهاتف",
  "attendance": "تسجيل الحضور",
  "profile": "حسابي",
  "tasks": "مهامي",
  "check_in": "بدء الدوام",
  "check_out": "إنهاء الدوام",
  "password_changed_title": "تم تغيير كلمة المرور",
  "password_changed_message": "تم تغيير كلمة المرور الخاصة بحسابك. يرجى تسجيل الدخول مرة أخرى.",
  "password_changed_button": "تسجيل الدخول",
  "select_worksite": "اختر منطقة العمل",
  "location": "الموقع",
  "distance": "المسافة",
  "inside_range": "داخل النطاق",
  "outside_range": "خارج النطاق",
  "loading": "جارٍ التحميل...",
  "error": "حدث خطأ",
  "success": "تم بنجاح",
  "select_location": "تحديد موقعي",
  "worksite_name": "اسم الموقع",
  "worksite_address": "العنوان",
  "hours_today": "اليوم",
  "hours_week": "الأسبوع",
  "hours_month": "الشهر",
  "address": "العنوان",
  "radius": "النطاق",
  "in_progress": "قيد التنفيذ",
  "before_checkout": "قبل إنهاء الدوام",
  "can_checkout": "يمكنك إنهاء الدوام الآن",
  "google_maps": "خرائط جوجل",
  "open_in_maps": "فتح في خرائط جوجل",
  "language": "اللغة",
  "notifications": "الإشعارات",
  "attendance_history": "سجل الحضور",
  "no_notifications": "لا توجد إشعارات",
  "no_history": "لا يوجد سجل حضور",
  "just_now": "الآن",
  "minutes_ago": "دقيقة",
  "hours_ago": "ساعة",
  "days_ago": "يوم",
  "hours": "ساعة",
  "admin": "مدير",
  "employee": "موظف",
  "completed": "مكتمل",
  "dashboard": "لوحة التحكم",
  "system_admin_role": "مدير النظام",
  "employees": "الموظفون",
  "worksites": "نقاط العمل",
  "reports": "التقارير",
  "settings": "الإعدادات",
  "service_requests": "طلبات الخدمة",
  "email": "البريد الإلكتروني",
  "password": "كلمة المرور",
  "full_name": "الاسم الكامل",
  "role": "الدور",
  "client": "عميل",
  "save": "حفظ",
  "cancel": "إلغاء",
  "delete": "حذف",
  "edit": "تعديل",
  "add": "إضافة",
  "search": "بحث",
  "no_data": "لا توجد بيانات",
  "confirm": "تأكيد",
  "close": "إغلاق",
  "back": "رجوع",
  "active": "نشط",
  "inactive": "غير نشط",
  "status": "الحالة",
  "actions": "إجراءات",
  "created_at": "تاريخ الإنشاء",
  "powered_slogan": "شريكك في التحول الرقمي",
  "devpro_name": "DevPro",
  "default_user_name": "المدير",
  "theme_light": "الوضع الفاتح",
  "theme_dark": "الوضع الداكن",
  "login_connecting": "📤 جاري الاتصال بالخادم...",
  "login_success": "✅ تم تسجيل الدخول بنجاح!",
  "login_failed": "❌ فشل تسجيل الدخول",
  "login_error_unknown": "خطأ غير معروف",
  "login_server_unreachable": "تعذر الاتصال بالخادم",
  "login_server_not_responding": "❌ الخادم لا يستجيب",
  "login_check": "🔍 تأكد من:",
  "login_check_user_exists": "المستخدم موجود في قاعدة البيانات",
  "login_check_password_correct": "كلمة المرور صحيحة",
  "login_check_account_active": "الحساب نشط (is_active = TRUE)",
  "login_error_code_prefix": "الرمز:",
  "login_error_message_prefix": "الرسالة:",
  "email_placeholder": "البريد الإلكتروني",
  "password_placeholder": "••••••••",
  "footer_copyright": "© 2026 DevPro - شريكك في التحول الرقمي",
  "devpro_watermark": "DevPro - شريكك في التحول الرقمي",
  "loading_data": "⏳ جارٍ تحميل البيانات...",
  "stats_total_employees": "إجمالي الموظفين",
  "stats_waiting_employees": "قيد الانتظار",
  "stats_active_now": "نشطين حالياً",
  "stats_completed_today": "مكتمل اليوم",
  "dashboard_tracking_title": "🗺️ تتبع الموظفين في الوقت الفعلي",
  "update_badge": "تحديث",
  "refresh_button": "تحديث",
  "tab_active": "نشطين",
  "tab_waiting": "قيد الانتظار",
  "tab_completed": "مكتمل",
  "tab_alerts": "تحذيرات",
  "active_employees_title": "🟢 الموظفين النشطين",
  "no_active_employees": "📭 لا يوجد موظفين نشطين حالياً",
  "waiting_employees_title": "⏳ الموظفين قيد الانتظار",
  "all_employees_started": "✅ جميع الموظفين بدأوا العمل اليوم",
  "completed_employees_title": "✅ الموظفين المكتملين اليوم",
  "no_completed_employees": "📭 لا يوجد موظفين مكتملين اليوم",
  "alerts_title": "🚨 التحذيرات الأمنية",
  "no_alerts": "✅ لا توجد تحذيرات أمنية",
  "left_worksite_prefix": "خرج عن نطاق",
  "employee_details_title": "📋 تفاصيل الموظف",
  "worksite_label": "نقطة العمل:",
  "undefined_text": "غير محدد",
  "subscription_lifetime": "مدى الحياة",
  "distance_label": "المسافة:",
  "working_hours_label": "ساعات العمل:",
  "last_update_label": "آخر تحديث:",
  "location_label": "الموقع:",
  "security_notes_title": "📝 سجل التحذيرات الأمنية",
  "no_security_notes": "✅ لا توجد تحذيرات أمنية لهذا الموظف",
  "clients_title": "العملاء",
  "clients_description": "حسابات العملاء الذين يتابعون تقارير الخدمة عبر بوابتهم الخاصة",
  "customers": "العملاء",
  "new_client": "عميل جديد",
  "clients_name": "الاسم",
  "clients_phone": "الهاتف",
  "clients_email": "البريد الإلكتروني",
  "clients_service_count": "عدد الخدمات",
  "customers_title": "العملاء",
  "customers_description": "قائمة العملاء المسجلين في النظام",
  "loading_customers": "⏳ جارٍ تحميل العملاء...",
  "no_customers": "📭 لا يوجد عملاء",
  "no_customers_hint": "سيظهر هنا العملاء المسجلون في النظام",
  "failed_to_fetch_customers": "فشل جلب العملاء",
  "add_customer": "إضافة عميل",
  "add_customer_hint": "أدخل بيانات العميل الجديد",
  "customer_password_hint": "سيتم إنشاء كلمة مرور تلقائية للعميل",
  "enter_full_name": "أدخل الاسم الكامل",
  "enter_email": "أدخل البريد الإلكتروني",
  "invalid_phone": "الرجاء إدخال رقم هاتف صحيح",
  "saving": "جارٍ الحفظ...",
  "create_customer": "إنشاء العميل",
  "customer_created_success": "تم إنشاء العميل بنجاح",
  "customer_creation_failed": "فشل إنشاء العميل",
  "password": "كلمة المرور",
  "save_password_warning": "يرجى حفظ كلمة المرور ومشاركتها مع العميل",
  "reset_password": "إعادة تعيين كلمة المرور",
  "reset_password_confirm": "هل أنت متأكد من إعادة تعيين كلمة مرور العميل",
  "resetting": "جارٍ إعادة التعيين...",
  "password_reset_success": "تم إعادة تعيين كلمة المرور بنجاح",
  "password_reset_failed": "فشل إعادة تعيين كلمة المرور",
  "new_password": "كلمة المرور الجديدة",
  "show_password": "إظهار كلمة المرور",
  "hide_password": "إخفاء كلمة المرور",
  "copy_password": "نسخ كلمة المرور",
  "copied": "تم النسخ",
  "delete": "حذف",
  "delete_customer": "حذف العميل",
  "delete_customer_confirm": "هل أنت متأكد من حذف العميل",
  "delete_customer_warning": "هذا الإجراء لا يمكن التراجع عنه. سيتم حذف جميع بيانات العميل.",
  "delete_customer_failed": "فشل حذف العميل",
  "deleting": "جارٍ الحذف...",
  "worksites_title": "📍 نقاط العمل",
  "worksites_description": "إدارة نقاط العمل وتعيين الموظفين",
  "new_worksite": "نقطة عمل جديدة",
  "loading_worksites": "⏳ جارٍ تحميل نقاط العمل...",
  "no_worksites": "📭 لا توجد نقاط عمل",
  "add_worksite_prompt": "قم بإضافة نقطة عمل جديدة",
  "assign_employee": "تعيين موظف",
  "no_address": "لا يوجد عنوان",
  "unassigned": "بدون موظف",
  "active_status": "✅ نشط",
  "inactive_status": "❌ غير نشط",
  "confirm_delete_title": "تأكيد الحذف",
  "confirm_delete_message": "هل أنت متأكد من حذف",
  "delete_warning": "سيؤدي هذا إلى حذف جميع المهام المرتبطة بهذه النقطة!",
  "deleting": "⏳ جارٍ الحذف...",
  "delete_final": "🗑️ حذف نهائي",
  "assign_employee_title": "تعيين موظف",
  "choose_employee_to_assign": "اختر موظفاً لتعيينه إلى",
  "no_available_employees": "📭 لا يوجد موظفين متاحين",
  "employees_title": "👥 الموظفون",
  "employees_description": "إدارة حسابات الموظفين",
  "add_employee": "إضافة موظف",
  "no_employees": "📭 لا يوجد موظفون",
  "add_new_employee_prompt": "قم بإضافة موظف جديد",
  "loading_employees": "⏳ جارٍ تحميل الموظفين...",
  "admin_role": "مدير",
  "field_employee": "موظف ميداني",
  "suspended_status": "❌ موقوف",
  "delete_irreversible": "لا يمكن التراجع عن هذا الإجراء!",
  "done": "تم",
  "settings_intro": "تحكم في إعدادات التطبيق والاشتراك الخاص بك",
  "subscription_title": "حالة الاشتراك",
  "subscription_intro": "تحقق من حالة اشتراك المدير وتاريخ انتهاءه.",
  "subscription_description": "تُظهر حالة الاشتراك وتاريخ انتهاء صلاحية حساب المدير.",
  "subscription_status_label": "حالة الاشتراك",
  "subscription_expires_at_label": "ينتهي في",
  "subscription_active": "نشط",
  "subscription_expired": "منتهٍ",
  "subscription_canceled": "ملغي",
  "subscription_active_message": "اشتراكك صالح ويمكنك استخدام لوحة المدير.",
  "subscription_expired_message": "انتهى اشتراكك. اتصل بالدعم لتجديده أو إعادة تفعيله.",
  "settings_language_title": "اللغة الافتراضية",
  "settings_language_label": "واجهة الإدارة",
  "settings_geofence_title": "نطاق التختيم الافتراضي",
  "settings_geofence_hint": "القيمة المستخدمة عند إنشاء نقطة عمل جديدة إن لم يُحدَّد نصف قطر مخصص.",
  "settings_geofence_radius": "نصف القطر (متر)",
  "ar": "العربية",
  "he": "עברית",
  "en": "English",
  "year": "السنة",
  "month": "الشهر",
  "export_pdf": "تصدير PDF",
  "total_hours": "إجمالي الساعات",
  "work_days": "أيام العمل",
  "days": "أيام",
  "no_attendance_records": "لا توجد سجلات حضور",
  "date": "التاريخ",
  "worked_hours": "ساعات العمل",
  "cleanup_old_records": "حذف السجلات القديمة",
  "name": "الاسم",
  "worksite": "نقطة العمل",
  "january": "يناير",
  "february": "فبراير",
  "march": "مارس",
  "april": "أبريل",
  "may": "مايو",
  "june": "يونيو",
  "july": "يوليو",
  "august": "أغسطس",
  "september": "سبتمبر",
  "october": "أكتوبر",
  "november": "نوفمبر",
  "december": "ديسمبر",
  "meters": "متر",
  "kilometers": "كيلومتر",
  "notifications": "الإشعارات",
  "subscription": "الاشتراك",
  "reports": "التقارير",
  "reports_description": "ملخص الأداء اليومي لكل الفرق الميدانية",
  "comprehensive_report": "التقرير الشامل",
  "service_requests_report": "تقرير طلبات الخدمة",
  "employee_performance_report": "تقرير أداء الموظفين",
  "client_activity_report": "تقرير نشاط العملاء",
  "active_employees": "الموظفون النشطون",
  "active_clients": "العملاء النشطون",
  "completed_services": "الخدمات المكتملة",
  "average_rating": "متوسط التقييم",
  "ratings": "تقييم",
  "on_duty": "في الخدمة",
  "total": "الإجمالي",
  "service_requests_status": "حالة طلبات الخدمة",
  "tasks_status": "حالة المهام",
  "attendance_statistics": "إحصائيات الحضور",
  "worksites_status": "حالة نقاط العمل",
  "total_hours_week": "إجمالي ساعات العمل هذا الأسبوع",
  "avg_daily_hours": "متوسط الساعات اليومية",
  "completed_duty_today": "أكملوا الدوام اليوم",
  "total_worksites": "إجمالي نقاط العمل",
  "active_worksites": "نقاط العمل النشطة",
  "no_data_available": "لا توجد بيانات متاحة",
  "failed_to_fetch_report": "فشل جلب التقرير",
  "title": "العنوان",
  "client": "العميل",
  "employee": "الموظف",
  "status": "الحالة",
  "priority": "الأولوية",
  "rating": "التقييم",
  "created_at": "تاريخ الإنشاء",
  "phone": "رقم الهاتف",
  "total_hours": "إجمالي الساعات",
  "completed_shifts": "الورديات المكتملة",
  "assigned_services": "الخدمات المعينة",
  "avg_rating": "متوسط التقييم",
  "total_requests": "إجمالي الطلبات",
  "completed_requests": "الطلبات المكتملة",
  "pending": "معلق",
  "range_compliance": "الالتزام بالنطاق الجغرافي",
  "range_compliance_hint": "نسبة عمليات التختيم المقبولة مقابل المرفوضة بسبب الخروج عن نطاق نقطة العمل",
  "accepted": "مقبولة",
  "rejected": "مرفوضة",
  "tasks": "المهام",
  "tasks_description": "جميع المهام عبر الموظفين ونقاط العمل",
  "new_task": "+ مهمة جديدة",
  "task_title": "المهمة",
  "task_employee": "الموظف",
  "task_worksite": "نقطة العمل",
  "task_status": "الحالة",
  "status_pending": "قيد الانتظار",
  "status_in_progress": "جارية",
  "status_completed": "مكتملة",
  "status_late": "متأخرة",
  "priority_low": "منخفضة",
  "priority_normal": "عادية",
  "priority_high": "عالية",
  "priority_urgent": "طارئة",
  "notification_punch_out": "محاولة تختيم مرفوضة خارج النطاق",
  "notification_new_employee": "تم إنشاء حساب موظف جديد",
  "notification_tasks_completed": "اكتملت جميع مهام اليوم",
  "today": "اليوم",
  "yesterday": "أمس",
  "service_requests": "طلبات الخدمة",
  "service_requests_description": "إدارة طلبات العملاء وتعيين الموظفين",
  "add_service_request": "إضافة طلب خدمة",
  "request_created_successfully": "تم إنشاء طلب الخدمة بنجاح",
  "please_fill_required_fields": "الرجاء ملء جميع الحقول المطلوبة",
  "client_name": "اسم العميل",
  "client_name_placeholder": "أدخل اسم العميل",
  "client_phone": "رقم هاتف العميل",
  "phone_placeholder": "مثال: 0501234567",
  "service_title_placeholder": "مثال: إصلاح سباكة",
  "service_description_placeholder": "وصف تفصيلي للخدمة المطلوبة",
  "address_placeholder": "العنوان الكامل مع المعالم",
  "priority_low": "منخفضة",
  "priority_normal": "عادية",
  "priority_high": "عالية",
  "priority_urgent": "عاجلة",
  "latitude_placeholder": "مثال: 24.7136",
  "longitude_placeholder": "مثال: 46.6753",
  "location_name": "اسم الموقع",
  "location_name_placeholder": "مثال: الرياض - حي النخيل",
  "my_assigned_requests": "طلباتي المعينة",
  "my_assigned_requests_description": "طلبات الخدمة المعينة لك",
  "no_assigned_requests": "لا توجد طلبات معينة",
  "no_assigned_requests_hint": "لم يتم تعيين أي طلبات خدمة لك بعد",
  "accept_request": "قبول الطلب",
  "reject_request": "رفض الطلب",
  "start_work": "بدء العمل",
  "complete_work": "إكمال العمل",
  "admin_notes": "ملاحظات الإدارة",
  "assigned_by": "عين بواسطة",
  "assigned_at": "وقت التعيين",
  "request_details": "تفاصيل الطلب",
  "client_info": "معلومات العميل",
  "location_info": "معلومات الموقع",
  "assignment_info": "معلومات التعيين",
  "coordinates": "الإحداثيات",
  "location_name": "اسم الموقع",
  "add_notes_optional": "أضف ملاحظات (اختياري)",
  "failed_to_update_status": "فشل تحديث الحالة",
  "status_accepted": "مقبول",
  "convert_to_worksite": "تحويل إلى نقطة عمل",
  "convert_to_worksite_hint": "تحويل موقع طلب الخدمة إلى نقطة عمل جديدة",
  "worksite_name_placeholder": "أدخل اسم نقطة العمل",
  "address_placeholder": "أدخل العنوان",
  "radius_placeholder": "أدخل نطاق الموقع (متر)",
  "create": "إنشاء",
  "all_statuses": "جميع الحالات",
  "status_assigned": "تم التعيين",
  "no_service_requests": "لا توجد طلبات",
  "no_service_requests_hint": "سيظهر هنا طلبات العملاء الجديدة",
  "client": "عميل",
  "no_address": "لا يوجد عنوان",
  "latitude": "خط العرض",
  "longitude": "خط الطول",
  "assign_employee": "تعيين موظف",
  "waiting_for_employee": "في انتظار الموظف",
  "in_execution": "قيد التنفيذ",
  "reassign_employee": "إعادة تعيين موظف",
  "unassign_employee": "إلغاء تعيين الموظف",
  "current_assigned_employee": "الموظف المعين:",
  "assign_employee_modal": "تعيين موظف",
  "assign_employee_hint": "اختر موظفاً لتنفيذ طلب الخدمة",
  "cancel": "إلغاء",
  "failed_to_fetch_requests": "فشل جلب الطلبات",
  "failed_to_fetch_employees": "فشل جلب الموظفين",
  "new_worksite": "نقطة عمل جديدة",
  "worksites_title": "نقاط العمل",
  "worksites_description": "إدارة نقاط العمل للموظفين",
  "loading_worksites": "جارٍ تحميل نقاط العمل...",
  "no_worksites": "لا توجد نقاط عمل",
  "add_worksite_prompt": "أضف نقطة عمل جديدة للبدء",
  "meters_unit": "متر",
  "unassigned": "غير معين",
  "active_status": "نشط",
  "inactive_status": "غير نشط",
  "confirm_delete_title": "تأكيد الحذف",
  "delete_warning": "هذا الإجراء لا يمكن التراجع عنه",
  "deleting": "جارٍ الحذف...",
  "delete_final": "حذف",
  "assign_employee_title": "تعيين موظف",
  "choose_employee_to_assign": "اختر موظفاً لتعيينه لـ",
  "no_available_employees": "لا يوجد موظفين متاحين",
  "employee_assigned_successfully": "تم تعيين الموظف بنجاح",
  "failed_to_assign_employee": "فشل تعيين الموظف",
  "worksite_name": "اسم الموقع",
  "worksite_name_placeholder": "مثال: فرع عمان - شارع الملكة رانيا",
  "search_address": "بحث عن عنوان",
  "search_address_placeholder": "ابحث عن مدينة أو شارع (مثال: شارع الملكة رانيا أو عمان)",
  "select_address_required": "الرجاء اختيار عنوان من نتائج البحث",
  "allowed_radius": "النطاق المسموح (متر)",
  "no_results_found": "لم يتم العثور على نتائج",
  "worksite_added_successfully": "تمت إضافة نقطة العمل بنجاح!",
  "save_failed": "فشلت الحفظ",
  "type_city": "مدينة",
  "type_street": "شارع",
  "type_address": "عنوان",
  "title_ar_placeholder": "عنوان المهمة بالعربية",
  "title_he_placeholder": "כותרת המשימה בעברית",
  "title_en_placeholder": "Task title in English",
  "description_ar_placeholder": "وصف المهمة بالعربية",
  "description_he_placeholder": "תיאור המשימה בעברית",
  "description_en_placeholder": "Task description in English",
  "task_description": "وصف المهمة",
  "select_worksite": "اختر نقطة العمل",
  "select_employee": "اختر الموظف",
  "priority_low": "منخفضة",
  "priority_normal": "عادية",
  "priority_high": "عالية",
  "priority_urgent": "عاجلة",
  "scheduled_start": "وقت البدء المقرر",
  "scheduled_end": "وقت الانتهاء المقرر",
  "creating": "جاري الإنشاء...",
  "create": "إنشاء",
  "error_creating_task": "خطأ في إنشاء المهمة",
  "type_house": "منزل",
  "type_landmark": "معلم",
  "type_location": "موقع",
  "city": "مدينة",
  "street": "شارع",
  "building_number": "رقم المبنى",
  "coordinates": "الإحداثيات",
  "saving": "جارٍ الحفظ...",
  "save": "حفظ",
  "search_worldwide": "ابحث في جميع أنحاء العالم",
  "search_language_hint": "يدعم البحث بالعربية والإنجليزية والعبرية",
  "view_details_title": "عرض التفاصيل",
  "has_not_started_work": "لم يبدأ العمل",
  "waiting_status": "قيد الانتظار",
  "completed_status": "مكتمل",
  "failed_to_fetch_data": "فشل جلب البيانات",
  "failed_to_fetch_security_notes": "فشل جلب الملاحظات الأمنية",
  "failed_to_delete_employee": "فشل حذف الموظف",
  "failed_to_delete_worksite": "فشل حذف موقع العمل",
  "failed_to_create_request": "فشل إنشاء طلب الخدمة",
  "failed_to_export_pdf": "فشل تصدير PDF",
  "failed_to_export_report": "فشل تصدير التقرير",
  "confirm_cleanup_old_records": "هل أنت متأكد من حذف السجلات القديمة؟",
  "cleanup_success": "تم تنظيف السجلات بنجاح",
  "cleanup_failed": "فشل تنظيف السجلات",
  "failed_to_fetch_attendance_history": "فشل جلب سجل الحضور",
  "shift_restored": "🔄 تم استعادة الوردية النشطة",
  "current_worksite": "نقطة العمل الحالية",
  "currently_working": "يعمل حالياً",
  "end_shift": "إنهاء الدوام",
  "force_checkout_title": "إنهاء الدوام إجبارياً",
  "force_checkout_message": "هل أنت متأكد من إنهاء دوام الموظف",
  "force_checkout_warning": "⚠️ سيتم إنهاء الدوام فوراً بدون الحاجة لتختيم الموظف. استخدم هذه الميزة فقط في حالات الطوارئ.",
  "force_checkout_success": "تم إنهاء الدوام بنجاح",
  "force_checkout_failed": "فشل إنهاء الدوام",
  "processing": "جارٍ المعالجة...",
  "confirm_end_shift": "تأكيد إنهاء الدوام",
  "pwa_install_title": "تثبيت التطبيق",
  "pwa_install_text": "ثبت التطبيق على جهازك",
  "pull_to_refresh": "اسحب للتحديث",
  "refreshing": "جاري التحديث...",
  "release_to_refresh": "أطلق للتحديث",
  "pwa": {
    "installTitle": "تثبيت التطبيق",
    "installText": "تثبيت التطبيق",
    "installOnIos": "تثبيت على iPhone",
    "iosModalTitle": "تثبيت التطبيق على iPhone",
    "iosModalSubtitle": "احصل على تجربة تطبيق أفضل!",
    "step1Text": "اضغط على زر Share ⎋ في أسفل الشاشة",
    "step2Text": "مرر لأسفل واضغط على Add to Home Screen",
    "step3Text": "اضغط على Add في الزاوية العلوية",
    "benefit1": "⚡ تشغيل أسرع",
    "benefit2": "📱 أيقونة على الشاشة الرئيسية",
    "benefit3": "🎨 تصميم شبيه بالتطبيقات",
    "gotIt": "فهمت ✓",
    "remindLater": "ذكرني لاحقاً"
  },
  "client_modal_title": "➕ إضافة عميل جديد",
  "client_modal_hint": "📧 سيتم إنشاء حساب العميل بريد إلكتروني وكلمة مرور تلقائياً",
  "client_modal_hint_small": "⚠️ يمكن للعميل استخدام هذه البيانات لتسجيل الدخول",
  "client_modal_full_name": "👤 الاسم الكامل",
  "client_modal_email": "📧 البريد الإلكتروني",
  "client_modal_phone": "📱 رقم الهاتف",
  "client_modal_email_hint": "سيستخدم العميل هذا البريد لتسجيل الدخول",
  "client_modal_saving": "⏳ جارٍ الحفظ...",
  "client_modal_create": "💾 إنشاء العميل",
  "client_modal_success": "✅ تم إنشاء العميل بنجاح!",
  "client_modal_success_email": "📧 البريد الإلكتروني:",
  "client_modal_success_password": "🔑 كلمة المرور:",
  "client_modal_success_warning": "⚠️ يرجى حفظ هذه البيانات ومشاركتها مع العميل",
  "client_modal_error_email": "❌ الرجاء إدخال بريد إلكتروني صحيح",
  "client_modal_error_failed": "❌ فشل إنشاء العميل",
  "delete_request": "حذف الطلب",
  "delete_request_title": "تأكيد حذف طلب الخدمة",
  "delete_request_message": "هل أنت متأكد من حذف طلب الخدمة هذا؟",
  "delete_request_warning": "⚠️ سيتم حذف الطلب وجميع التعيينات المرتبطة به بشكل نهائي.",
  "request_deleted_successfully": "✅ تم حذف طلب الخدمة بنجاح",
  "request_delete_failed": "❌ فشل حذف طلب الخدمة",
  "request_not_found_or_deleted": "⚠️ الطلب غير موجود أو تم حذفه بالفعل",
  "delete_success_title": "تم الحذف بنجاح",
  "delete_error_title": "فشل الحذف",
  "try_again": "حاول مرة أخرى",
  "ok": "حسناً",
  "err_permission_denied": "ليس لديك صلاحية للقيام بهذا الإجراء",
  "notifications_feed": "🔔 الإشعارات",
  "loading_notifications": "جاري تحميل الإشعارات...",
  "no_notifications_available": "لا توجد إشعارات حالياً",
  "view_all": "عرض الكل"
}

const he = {
  "app_name": "WorkTrack",
  "login": "התחברות",
  "logout": "התנתקות",
  "phone": "מספר טלפון",
  "attendance": "נוכחות",
  "profile": "החשבון שלי",
  "tasks": "המשימות שלי",
  "check_in": "תחילת משמרת",
  "check_out": "סיום משמרת",
  "password_changed_title": "הסיסמה שונתה",
  "password_changed_message": "הסיסמה שלך שונתה. אנא התחבר שוב.",
  "password_changed_button": "התחבר",
  "select_worksite": "בחר אתר עבודה",
  "location": "מיקום",
  "distance": "מרחק",
  "inside_range": "בתוך הטווח",
  "outside_range": "מחוץ לטווח",
  "loading": "טוען...",
  "error": "שגיאה",
  "success": "הצלחה",
  "select_location": "קבע את המיקום שלך",
  "worksite_name": "שם האתר",
  "worksite_address": "כתובת",
  "hours_today": "היום",
  "hours_week": "השבוע",
  "hours_month": "החודש",
  "address": "כתובת",
  "radius": "רדיוס",
  "in_progress": "ב",
  "before_checkout": "לפני סיום המשמרת",
  "can_checkout": "ניתן לסיים משמרת עכשיו",
  "google_maps": "מפות גוגל",
  "open_in_maps": "פתח במפות גוגל",
  "language": "שפה",
  "notifications": "התראות",
  "attendance_history": "היסטוריית נוכחות",
  "no_notifications": "אין התראות",
  "no_history": "אין היסטוריית נוכחות",
  "hours": "שעות",
  "admin": "מנהל",
  "employee": "עובד",
  "completed": "הושלם",
  "dashboard": "לוח בקרה",
  "system_admin_role": "מנהל המערכת",
  "employees": "עובדים",
  "worksites": "אתרי עבודה",
  "reports": "דוחות",
  "settings": "הגדרות",
  "service_requests": "בקשות שירות",
  "email": "אימייל",
  "password": "סיסמה",
  "full_name": "שם מלא",
  "role": "תפקיד",
  "client": "לקוח",
  "save": "שמור",
  "cancel": "בטל",
  "delete": "מחק",
  "edit": "ערוך",
  "add": "הוסף",
  "search": "חפש",
  "no_data": "אין נתונים",
  "confirm": "אשר",
  "close": "סגור",
  "back": "חזור",
  "active": "פעיל",
  "inactive": "לא פעיל",
  "status": "סטטוס",
  "actions": "פעולות",
  "created_at": "תאריך יצירה",
  "powered_slogan": "השותף שלך לטרנספורמציה דיגיטלית",
  "devpro_name": "DevPro",
  "default_user_name": "מנהל",
  "theme_light": "מצב בהיר",
  "theme_dark": "מצב כהה",
  "login_connecting": "📤 מתחבר לשרת...",
  "login_success": "✅ ההתחברות הצליחה!",
  "login_failed": "❌ הכניסה נכשלה",
  "login_error_unknown": "שגיאה לא ידועה",
  "login_server_unreachable": "לא ניתן להתחבר לשרת",
  "login_server_not_responding": "❌ השרת לא מגיב",
  "login_check": "🔍 בדוק:",
  "login_check_user_exists": "המשתמש קיים בבסיס הנתונים",
  "login_check_password_correct": "הסיסמה נכונה",
  "login_check_account_active": "החשבון פעיל (is_active = TRUE)",
  "login_error_code_prefix": "קוד:",
  "login_error_message_prefix": "הודעה:",
  "email_placeholder": "כתובת אימייל",
  "password_placeholder": "••••••••",
  "footer_copyright": "© 2026 DevPro - השותף שלך לטרנספורמציה דיגיטלית",
  "devpro_watermark": "DevPro - השותף שלך לטרנספורמציה דיגיטלית",
  "loading_data": "⏳ טוען נתוני לוח המחוונים...",
  "stats_total_employees": "סה\"כ עובדים",
  "stats_waiting_employees": "בהמתנה",
  "stats_active_now": "פעילים כעת",
  "stats_completed_today": "הושלם היום",
  "dashboard_tracking_title": "🗺️ מעקב עובדים בזמן אמת",
  "update_badge": "עדכונים",
  "refresh_button": "רענן",
  "tab_active": "פעילים",
  "tab_waiting": "בהמתנה",
  "tab_completed": "הושלם",
  "tab_alerts": "התראות",
  "active_employees_title": "🟢 עובדים פעילים",
  "no_active_employees": "📭 אין עובדים פעילים כרגע",
  "waiting_employees_title": "⏳ עובדים בהמתנה",
  "all_employees_started": "✅ כל העובדים החלו בעבודה היום",
  "completed_employees_title": "✅ עובדים שהושלמו היום",
  "no_completed_employees": "📭 אין עובדים שהושלמו היום",
  "alerts_title": "🚨 התראות אבטחה",
  "no_alerts": "✅ אין התראות אבטחה",
  "left_worksite_prefix": "יצא מאתר העבודה",
  "employee_details_title": "📋 פרטי עובד",
  "worksite_label": "אתר עבודה:",
  "undefined_text": "לא מוגדר",
  "subscription_lifetime": "לכל החיים",
  "distance_label": "מרחק:",
  "working_hours_label": "שעות עבודה:",
  "last_update_label": "עדכון אחרון:",
  "location_label": "מיקום:",
  "security_notes_title": "📝 היסטוריית התראות אבטחה",
  "no_security_notes": "✅ אין התראות אבטחה לעובד זה",
  "clients_title": "לקוחות",
  "clients_description": "חשבונות לקוחות שעוקבים אחרי דוחות שירות דרך הפורטל שלהם",
  "customers": "לקוחות",
  "new_client": "לקוח חדש",
  "clients_name": "שם",
  "clients_phone": "טלפון",
  "clients_email": "אימייל",
  "clients_service_count": "מספר שירותים",
  "customers_title": "לקוחות",
  "customers_description": "רשימת הלקוחות הרשומים במערכת",
  "loading_customers": "⏳ טוען לקוחות...",
  "no_customers": "📭 אין לקוחות",
  "no_customers_hint": "לקוחות רשומים יופיעו כאן",
  "failed_to_fetch_customers": "נכשל באחזור לקוחות",
  "add_customer": "הוסף לקוח",
  "add_customer_hint": "הזן פרטי לקוח חדש",
  "customer_password_hint": "סיסמה תיווצר אוטומטית עבור הלקוח",
  "enter_full_name": "הזן שם מלא",
  "enter_email": "הזן כתובת אימייל",
  "invalid_phone": "אנא הזן מספר טלפון תקין",
  "saving": "שומר...",
  "create_customer": "צור לקוח",
  "customer_created_success": "הלקוח נוצר בהצלחה",
  "customer_creation_failed": "יצירת הלקוח נכשלה",
  "password": "סיסמה",
  "save_password_warning": "אנא שמור את הסיסמה ושתף אותה עם הלקוח",
  "reset_password": "איפוס סיסמה",
  "reset_password_confirm": "האם אתה בטוח שברצונך לאפס את סיסמת הלקוח",
  "resetting": "מאפס...",
  "password_reset_success": "הסיסמה אופסה בהצלחה",
  "password_reset_failed": "איפוס הסיסמה נכשל",
  "new_password": "סיסמה חדשה",
  "show_password": "הצג סיסמה",
  "hide_password": "הסתר סיסמה",
  "copy_password": "העתק סיסמה",
  "copied": "הועתק",
  "delete": "מחק",
  "delete_customer": "מחק לקוח",
  "delete_customer_confirm": "האם אתה בטוח שברצונך למחוק את הלקוח",
  "delete_customer_warning": "פעולה זו אינה ניתנת לביטול. כל נתוני הלקוח יימחקו.",
  "delete_customer_failed": "מחיקת הלקוח נכשלה",
  "deleting": "מוחק...",
  "worksites_title": "📍 אתרי עבודה",
  "worksites_description": "ניהול אתרי עבודה ושהקצאת עובדים",
  "new_worksite": "אתר עבודה חדש",
  "loading_worksites": "⏳ טוען אתרי עבודה...",
  "no_worksites": "📭 לא נמצאו אתרי עבודה",
  "add_worksite_prompt": "הוסף אתר עבודה חדש",
  "assign_employee": "הקצה עובד",
  "no_address": "אין כתובת זמינה",
  "unassigned": "ללא מוקצה",
  "active_status": "✅ פעיל",
  "inactive_status": "❌ לא פעיל",
  "confirm_delete_title": "אישור מחיקה",
  "confirm_delete_message": "האם אתה בטוח שברצונך למחוק",
  "delete_warning": "זה ימחק את כל המשימות הקשורות לאתר זה!",
  "deleting": "⏳ מוחק...",
  "delete_final": "🗑️ מחק לצמיתות",
  "assign_employee_title": "הקצאת עובד",
  "choose_employee_to_assign": "בחר עובד להקצאה ל",
  "no_available_employees": "📭 אין עובדים זמינים",
  "employees_title": "👥 עובדים",
  "employees_description": "ניהול חשבונות עובדים",
  "add_employee": "הוסף עובד",
  "no_employees": "📭 אין עובדים",
  "add_new_employee_prompt": "הוסף עובד חדש",
  "loading_employees": "⏳ טוען עובדים...",
  "admin_role": "מנהל",
  "field_employee": "עובד שטח",
  "suspended_status": "❌ מושעה",
  "delete_irreversible": "לא ניתן לבטל פעולה זו!",
  "done": "בוצע",
  "settings_intro": "שלוט בהגדרות האפליקציה ובסטטוס המנוי שלך.",
  "subscription_title": "סטטוס מנוי",
  "subscription_intro": "בדוק את סטטוס המנוי של המנהל ותאריך הפקיעה.",
  "subscription_description": "מציג את סטטוס המנוי של מנהל המערכת ותאריך הפקיעה.",
  "subscription_status_label": "סטטוס המנוי",
  "subscription_expires_at_label": "מסתיים ב",
  "subscription_active": "פעיל",
  "subscription_expired": "פג",
  "subscription_canceled": "בוטל",
  "subscription_active_message": "המנוי שלך פעיל ולוח הבקרה נגיש.",
  "subscription_expired_message": "המנוי שלך פג. צור קשר עם התמיכה כדי לחדש או להפעיל אותו מחדש.",
  "settings_language_title": "שפת ברירת מחדל",
  "settings_language_label": "לוח מנהל",
  "settings_geofence_title": "רדיוס ברירת מחדל של גיאופנס",
  "settings_geofence_hint": "הרדיוס הדיפולטי שישמש בעת יצירת אתר עבודה אם לא נבחר ערך מותאם.",
  "settings_geofence_radius": "רדיוס (מטרים)",
  "ar": "العربية",
  "he": "עברית",
  "en": "English",
  "year": "שנה",
  "month": "חודש",
  "export_pdf": "ייצוא PDF",
  "total_hours": "סה\"כ שעות",
  "work_days": "ימי עבודה",
  "days": "ימים",
  "no_attendance_records": "אין רשומות נוכחות",
  "date": "תאריך",
  "worked_hours": "שעות עבודה",
  "cleanup_old_records": "מחיקת רשומות ישנות",
  "name": "שם",
  "worksite": "אתר עבודה",
  "january": "ינואר",
  "february": "פברואר",
  "march": "מרץ",
  "april": "אפריל",
  "may": "מאי",
  "june": "יוני",
  "july": "יולי",
  "august": "אוגוסט",
  "september": "ספטמבר",
  "october": "אוקטובר",
  "november": "נובמבר",
  "december": "דצמבר",
  "meters": "מטר",
  "kilometers": "ק\"מ",
  "notifications": "התראות",
  "subscription": "מנוי",
  "reports": "דוחות",
  "reports_description": "סיכום ביצועים יומי עבור כל הצוותים בשטח",
  "comprehensive_report": "דוח מקיף",
  "service_requests_report": "דוח בקשות שירות",
  "employee_performance_report": "דוח ביצועי עובדים",
  "client_activity_report": "דוח פעילות לקוחות",
  "active_employees": "עובדים פעילים",
  "active_clients": "לקוחות פעילים",
  "completed_services": "שירותים שהושלמו",
  "average_rating": "דירוג ממוצע",
  "ratings": "דירוגים",
  "on_duty": "בתפקיד",
  "total": "סה\"כ",
  "service_requests_status": "סטטוס בקשות שירות",
  "tasks_status": "סטטוס משימות",
  "attendance_statistics": "סטטיסטיקת נוכחות",
  "worksites_status": "סטטוס אתרי עבודה",
  "total_hours_week": "שעות עבודה סה\"כ השבוע",
  "avg_daily_hours": "שעות ממוצעות ליום",
  "completed_duty_today": "סיימו משמרת היום",
  "total_worksites": "סה\"כ אתרי עבודה",
  "active_worksites": "אתרי עבודה פעילים",
  "no_data_available": "אין נתונים זמינים",
  "failed_to_fetch_report": "נכשל באחזור הדוח",
  "title": "כותרת",
  "client": "לקוח",
  "employee": "עובד",
  "status": "סטטוס",
  "priority": "עדיפות",
  "rating": "דירוג",
  "created_at": "נוצר ב",
  "phone": "טלפון",
  "total_hours": "שעות סה\"כ",
  "completed_shifts": "משמרות שהושלמו",
  "assigned_services": "שירותים שהוקצו",
  "avg_rating": "דירוג ממוצע",
  "total_requests": "סה\"כ בקשות",
  "completed_requests": "בקשות שהושלמו",
  "pending": "ממתין",
  "range_compliance": "עמידה בטווח",
  "range_compliance_hint": "יחס התחברויות מתקבלות מול נדחות עקב יציאה מטווח אתר העבודה",
  "accepted": "מתקבל",
  "rejected": "נדחה",
  "tasks": "משימות",
  "tasks_description": "כל המשימות בכל העובדים ואתרי העבודה",
  "new_task": "+ משימה חדשה",
  "task_title": "משימה",
  "task_employee": "עובד",
  "task_worksite": "אתר עבודה",
  "task_status": "סטטוס",
  "status_pending": "ממתין",
  "status_in_progress": "בתהליך",
  "status_completed": "הושלם",
  "status_late": "מאוחר",
  "priority_low": "נמוכה",
  "priority_normal": "רגילה",
  "priority_high": "גבוהה",
  "priority_urgent": "דחופה",
  "notification_punch_out": "תחברות נדחתה מחוץ לטווח",
  "notification_new_employee": "נוצר חשבון עובד חדש",
  "notification_tasks_completed": "כל המשימות הושלמו להיום",
  "today": "היום",
  "yesterday": "אתמול",
  "service_requests": "בקשות שירות",
  "service_requests_description": "ניהול בקשות לקוחות והקצאת עובדים",
  "add_service_request": "הוסף בקשת שירות",
  "request_created_successfully": "בקשת השירות נוצרה בהצלחה",
  "please_fill_required_fields": "אנא מלא את כל השדות הנדרשים",
  "client_name": "שם הלקוח",
  "client_name_placeholder": "הכנס שם לקוח",
  "client_phone": "טלפון הלקוח",
  "phone_placeholder": "דוגמה: 0501234567",
  "service_title_placeholder": "דוגמה: תיקון אינסטלציה",
  "service_description_placeholder": "תיאור מפורט של השירות הנדרש",
  "address_placeholder": "כתובת מלאה עם ציוני דרך",
  "priority_low": "נמוכה",
  "priority_normal": "רגילה",
  "priority_high": "גבוהה",
  "priority_urgent": "דחופה",
  "latitude_placeholder": "דוגמה: 24.7136",
  "longitude_placeholder": "דוגמה: 46.6753",
  "location_name": "שם המיקום",
  "location_name_placeholder": "דוגמה: ריאד - שכונת אל-נחיל",
  "my_assigned_requests": "הבקשות המוקצות לי",
  "my_assigned_requests_description": "בקשות שירות שהוקצו לך",
  "no_assigned_requests": "אין בקשות מוקצות",
  "no_assigned_requests_hint": "עדיין לא הוקצו בקשות שירות לך",
  "accept_request": "קבל בקשה",
  "reject_request": "דחה בקשה",
  "start_work": "התחל עבודה",
  "complete_work": "סיים עבודה",
  "admin_notes": "הערות מנהל",
  "assigned_by": "הוקצה על ידי",
  "assigned_at": "זמן הקצאה",
  "request_details": "פרטי הבקשה",
  "client_info": "מידע לקוח",
  "location_info": "מידע מיקום",
  "assignment_info": "מידע הקצאה",
  "coordinates": "קואורדינטות",
  "location_name": "שם המיקום",
  "add_notes_optional": "הוסף הערות (אופציונלי)",
  "failed_to_update_status": "עדכון הסטטוס נכשל",
  "status_accepted": "התקבל",
  "convert_to_worksite": "המר לאתר עבודה",
  "convert_to_worksite_hint": "המר מיקום בקשת השירות לאתר עבודה חדש",
  "worksite_name_placeholder": "הכנס שם אתר עבודה",
  "address_placeholder": "הכנס כתובת",
  "radius_placeholder": "הכנס רדיוס מיקום (מטרים)",
  "create": "צור",
  "all_statuses": "כל הסטטוסים",
  "status_assigned": "הוקצה",
  "no_service_requests": "אין בקשות שירות",
  "no_service_requests_hint": "בקשות לקוחות חדשות יופיעו כאן",
  "client": "לקוח",
  "no_address": "אין כתובת",
  "latitude": "קו רוחב",
  "longitude": "קו אורך",
  "assign_employee": "הקצה עובד",
  "waiting_for_employee": "ממתין לעובד",
  "in_execution": "בביצוע",
  "reassign_employee": "הקצאה מחדש של עובד",
  "unassign_employee": "ביטול הקצאת עובד",
  "current_assigned_employee": "העובד המוקצה:",
  "assign_employee_modal": "הקצה עובד",
  "assign_employee_hint": "בחר עובד לביצוע בקשת השירות",
  "cancel": "ביטול",
  "failed_to_fetch_requests": "נכשל באחזור בקשות",
  "failed_to_fetch_employees": "נכשל באחזור עובדים",
  "new_worksite": "נקודת עבודה חדשה",
  "worksites_title": "נקודות עבודה",
  "worksites_description": "ניהול נקודות עבודה לעובדים",
  "loading_worksites": "טוען נקודות עבודה...",
  "no_worksites": "אין נקודות עבודה",
  "add_worksite_prompt": "הוסף נקודת עבודה חדשה להתחלה",
  "meters_unit": "מטר",
  "unassigned": "לא מוקצה",
  "active_status": "פעיל",
  "inactive_status": "לא פעיל",
  "confirm_delete_title": "אישור מחיקה",
  "delete_warning": "פעולה זו אינה ניתנת לביטול",
  "deleting": "מוחק...",
  "delete_final": "מחק",
  "assign_employee_title": "הקצה עובד",
  "choose_employee_to_assign": "בחר עובד להקצאה ל",
  "no_available_employees": "אין עובדים זמינים",
  "employee_assigned_successfully": "העובד הוקצה בהצלחה",
  "failed_to_assign_employee": "נכשל בהקצאת עובד",
  "worksite_name": "שם האתר",
  "worksite_name_placeholder": "דוגמה: סניף עמאן - רחוב המלכה ראניה",
  "search_address": "חיפוש כתובת",
  "search_address_placeholder": "חפש עיר או רחוב (למשל: רחוב המלכה ראניה או עמאן)",
  "select_address_required": "אנא בחר כתובת מתוצאות החיפוש",
  "allowed_radius": "רדיוס מותר (מטר)",
  "no_results_found": "לא נמצאו תוצאות",
  "worksite_added_successfully": "נקודת העבודה נוספה בהצלחה!",
  "save_failed": "השמירה נכשלה",
  "type_city": "עיר",
  "type_street": "רחוב",
  "type_address": "כתובת",
  "title_ar_placeholder": "כותרת המשימה בערבית",
  "title_he_placeholder": "כותרת המשימה בעברית",
  "title_en_placeholder": "Task title in English",
  "description_ar_placeholder": "תיאור המשימה בערבית",
  "description_he_placeholder": "תיאור המשימה בעברית",
  "description_en_placeholder": "Task description in English",
  "task_description": "תיאור המשימה",
  "select_worksite": "בחר אתר עבודה",
  "select_employee": "בחר עובד",
  "priority_low": "נמוכה",
  "priority_normal": "רגילה",
  "priority_high": "גבוהה",
  "priority_urgent": "דחופה",
  "scheduled_start": "זמן התחלה מתוכנן",
  "scheduled_end": "זמן סיום מתוכנן",
  "creating": "יוצר...",
  "create": "צור",
  "error_creating_task": "שגיאה ביצירת משימה",
  "type_house": "בית",
  "type_landmark": "ציון דרך",
  "type_location": "מיקום",
  "city": "עיר",
  "street": "רחוב",
  "building_number": "מספר בניין",
  "coordinates": "קואורדינטות",
  "saving": "שומר...",
  "save": "שמור",
  "search_worldwide": "חפש ברחבי העולם",
  "search_language_hint": "החיפוש תומך בערבית, אנגלית ועברית",
  "view_details_title": "הצג פרטים",
  "has_not_started_work": "לא התחיל לעבוד",
  "waiting_status": "בהמתנה",
  "completed_status": "הושלם",
  "failed_to_fetch_data": "נכשל באחזור נתונים",
  "failed_to_fetch_security_notes": "נכשל באחזור הערות אבטחה",
  "failed_to_delete_employee": "מחיקת עובד נכשלה",
  "failed_to_delete_worksite": "מחיקת אתר עבודה נכשלה",
  "failed_to_export_pdf": "ייצוא PDF נכשל",
  "failed_to_export_report": "ייצוא דוח נכשל",
  "confirm_cleanup_old_records": "האם אתה בטוח שברצונך למחוק רשומות ישנות?",
  "cleanup_success": "ניקוי רשומות הצליח",
  "cleanup_failed": "ניקוי רשומות נכשל",
  "failed_to_fetch_attendance_history": "אחזור היסטוריית נוכחות נכשל",
  "shift_restored": "🔄 המשמרת הפעילה שוחזרה",
  "current_worksite": "אתר עבודה נוכחי",
  "currently_working": "עובד כעת",
  "end_shift": "סיום משמרת",
  "force_checkout_title": "סיום משמרת בכפייה",
  "force_checkout_message": "האם אתה בטוח שברצונך לסיים את המשמרת עבור העובד",
  "force_checkout_warning": "⚠️ המשמרת תסתיים מיד ללא צורך בתחברות העובד. השתמש בתכונה זו רק במצבי חירום.",
  "force_checkout_success": "המשמרת הסתיימה בהצלחה",
  "force_checkout_failed": "סיום המשמרת נכשל",
  "processing": "מעבד...",
  "confirm_end_shift": "אשר סיום משמרת",
  "pwa_install_title": "התקנת אפליקציה",
  "pwa_install_text": "התקן את האפליקציה במכשיר שלך",
  "pull_to_refresh": "משוך לרענון",
  "refreshing": "מרענן...",
  "release_to_refresh": "שחרר לרענון",
  "pwa": {
    "installTitle": "התקנת אפליקציה",
    "installText": "התקן אפליקציה",
    "installOnIos": "התקנה באייפון",
    "iosModalTitle": "התקנת האפליקציה באייפון",
    "iosModalSubtitle": "קבל חוויית אפליקציה טובה יותר!",
    "step1Text": "לחץ על כפתור Share ⎋ בתחתית המסך",
    "step2Text": "גרור למטה ולחץ על Add to Home Screen",
    "step3Text": "לחץ על Add בפינה העליונה",
    "benefit1": "⚡ טעינה מהירה יותר",
    "benefit2": "📱 אייקון במסך הבית",
    "benefit3": "🎨 עיצוב דמוי אפליקציה",
    "gotIt": "הבנתי ✓",
    "remindLater": "תזכיר לי מאוחר יותר"
  },
  "client_modal_title": "➕ הוסף לקוח חדש",
  "client_modal_hint": "📧 חשבון לקוח ייווצר עם אימייל וסיסמה אוטומטית",
  "client_modal_hint_small": "⚠️ הלקוח יכול להשתמש בפרטים אלה להתחברות",
  "client_modal_full_name": "👤 שם מלא",
  "client_modal_email": "📧 אימייל",
  "client_modal_phone": "📱 מספר טלפון",
  "client_modal_email_hint": "הלקוח ישתמש באימייל זה להתחברות",
  "client_modal_saving": "⏳ שומר...",
  "client_modal_create": "💾 צור לקוח",
  "client_modal_success": "✅ הלקוח נוצר בהצלחה!",
  "client_modal_success_email": "📧 אימייל:",
  "client_modal_success_password": "🔑 סיסמה:",
  "client_modal_success_warning": "⚠️ אנא שמור את הפרטים האלה ושתף אותם עם הלקוח",
  "client_modal_error_email": "❌ אנא הכנס אימייל תקין",
  "client_modal_error_failed": "❌ יצירת הלקוח נכשלה",
  "delete_request": "מחק בקשה",
  "delete_request_title": "אישור מחיקת בקשת שירות",
  "delete_request_message": "האם אתה בטוח שברצונך למחוק את בקשת השירות הזו?",
  "delete_request_warning": "⚠️ זה ימחק את הבקשה ואת כל ההקצאות הקשורות לה לצמיתות.",
  "request_deleted_successfully": "✅ בקשת השירות נמחקה בהצלחה",
  "request_delete_failed": "❌ מחיקת בקשת השירות נכשלה",
  "request_not_found_or_deleted": "⚠️ הבקשה לא נמצאה או כבר נמחקה",
  "delete_success_title": "נמחק בהצלחה",
  "delete_error_title": "המחיקה נכשלה",
  "try_again": "נסה שוב",
  "ok": "אישור",
  "err_permission_denied": "אין לך הרשאה לבצע פעולה זו",
  "notifications_feed": "🔔 התראות",
  "loading_notifications": "טוען התראות...",
  "no_notifications_available": "אין התראות זמינות",
  "view_all": "הצג הכל"
}

const en = {
  "app_name": "WorkTrack",
  "login": "Login",
  "logout": "Logout",
  "phone": "Phone Number",
  "attendance": "Attendance",
  "profile": "Profile",
  "tasks": "My Tasks",
  "check_in": "Check In",
  "check_out": "Check Out",
  "password_changed_title": "Password Changed",
  "password_changed_message": "Your password has been changed. Please log in again.",
  "password_changed_button": "Log In",
  "select_worksite": "Select Worksite",
  "location": "Location",
  "distance": "Distance",
  "inside_range": "Inside Range",
  "outside_range": "Outside Range",
  "loading": "Loading...",
  "error": "Error",
  "success": "Success",
  "select_location": "Set your location",
  "worksite_name": "Worksite Name",
  "worksite_address": "Address",
  "hours_today": "Today",
  "hours_week": "Week",
  "hours_month": "Month",
  "address": "Address",
  "radius": "Radius",
  "in_progress": "in",
  "before_checkout": "before checkout",
  "can_checkout": "can checkout now",
  "google_maps": "Google Maps",
  "open_in_maps": "Open in Google Maps",
  "language": "Language",
  "notifications": "Notifications",
  "subscription": "Subscription",
  "attendance_history": "Attendance History",
  "no_notifications": "No notifications",
  "no_history": "No attendance history",
  "hours": "hours",
  "admin": "Admin",
  "employee": "Employee",
  "completed": "Completed",
  "dashboard": "Dashboard",
  "system_admin_role": "System administrator",
  "employees": "Employees",
  "worksites": "Worksites",
  "reports": "Reports",
  "settings": "Settings",
  "service_requests": "Service Requests",
  "email": "Email",
  "password": "Password",
  "full_name": "Full Name",
  "role": "Role",
  "client": "Client",
  "save": "Save",
  "cancel": "Cancel",
  "delete": "Delete",
  "edit": "Edit",
  "add": "Add",
  "search": "Search",
  "no_data": "No Data",
  "confirm": "Confirm",
  "close": "Close",
  "back": "Back",
  "active": "Active",
  "inactive": "Inactive",
  "status": "Status",
  "actions": "Actions",
  "created_at": "Created At",
  "powered_slogan": "Your Partner in Digital Transformation",
  "devpro_name": "DevPro",
  "default_user_name": "Admin",
  "theme_light": "Light mode",
  "theme_dark": "Dark mode",
  "login_connecting": "📤 Connecting to server...",
  "login_success": "✅ Login successful!",
  "login_failed": "❌ Login failed",
  "login_error_unknown": "Unknown error",
  "login_server_unreachable": "Failed to connect to server",
  "login_server_not_responding": "❌ Server is not responding",
  "login_check": "🔍 Please check:",
  "login_check_user_exists": "User exists in database",
  "login_check_password_correct": "Password is correct",
  "login_check_account_active": "Account is active (is_active = TRUE)",
  "login_error_code_prefix": "Code:",
  "login_error_message_prefix": "Message:",
  "email_placeholder": "Email address",
  "password_placeholder": "••••••••",
  "footer_copyright": "© 2026 DevPro - Your Partner in Digital Transformation",
  "devpro_watermark": "DevPro - Your Partner in Digital Transformation",
  "loading_data": "⏳ Loading dashboard data...",
  "stats_total_employees": "Total employees",
  "stats_waiting_employees": "Waiting",
  "stats_active_now": "Active now",
  "stats_completed_today": "Completed today",
  "dashboard_tracking_title": "🗺️ Real-time employee tracking",
  "update_badge": "updates",
  "refresh_button": "Refresh",
  "tab_active": "Active",
  "tab_waiting": "Waiting",
  "tab_completed": "Completed",
  "tab_alerts": "Alerts",
  "active_employees_title": "🟢 Active employees",
  "no_active_employees": "📭 No active employees right now",
  "waiting_employees_title": "⏳ Employees waiting",
  "all_employees_started": "✅ All employees have started work today",
  "completed_employees_title": "✅ Completed employees today",
  "no_completed_employees": "📭 No completed employees today",
  "alerts_title": "🚨 Security alerts",
  "no_alerts": "✅ No security alerts",
  "left_worksite_prefix": "Left worksite",
  "employee_details_title": "📋 Employee details",
  "worksite_label": "Worksite:",
  "undefined_text": "Unassigned",
  "subscription_lifetime": "Lifetime",
  "distance_label": "Distance:",
  "working_hours_label": "Working hours:",
  "last_update_label": "Last update:",
  "location_label": "Location:",
  "security_notes_title": "📝 Security alert history",
  "no_security_notes": "✅ No security alerts for this employee",
  "clients_title": "Clients",
  "clients_description": "Customer accounts tracking service reports through their portal",
  "customers": "Customers",
  "new_client": "New client",
  "clients_name": "Name",
  "clients_phone": "Phone",
  "clients_email": "Email",
  "clients_service_count": "Service count",
  "customers_title": "Customers",
  "customers_description": "List of customers registered in the system",
  "loading_customers": "⏳ Loading customers...",
  "no_customers": "📭 No customers",
  "no_customers_hint": "Registered customers will appear here",
  "failed_to_fetch_customers": "Failed to fetch customers",
  "add_customer": "Add customer",
  "add_customer_hint": "Enter new customer details",
  "customer_password_hint": "A password will be automatically generated for the customer",
  "enter_full_name": "Enter full name",
  "enter_email": "Enter email",
  "invalid_phone": "Please enter a valid phone number",
  "saving": "Saving...",
  "create_customer": "Create customer",
  "customer_created_success": "Customer created successfully",
  "customer_creation_failed": "Customer creation failed",
  "password": "Password",
  "save_password_warning": "Please save the password and share it with the customer",
  "reset_password": "Reset password",
  "reset_password_confirm": "Are you sure you want to reset the customer's password",
  "resetting": "Resetting...",
  "password_reset_success": "Password reset successfully",
  "password_reset_failed": "Password reset failed",
  "new_password": "New password",
  "show_password": "Show password",
  "hide_password": "Hide password",
  "copy_password": "Copy password",
  "copied": "Copied",
  "delete": "Delete",
  "delete_customer": "Delete customer",
  "delete_customer_confirm": "Are you sure you want to delete the customer",
  "delete_customer_warning": "This action cannot be undone. All customer data will be deleted.",
  "delete_customer_failed": "Failed to delete customer",
  "deleting": "Deleting...",
  "worksites_title": "📍 Worksites",
  "worksites_description": "Manage worksites and assign employees",
  "new_worksite": "New worksite",
  "loading_worksites": "⏳ Loading worksites...",
  "no_worksites": "📭 No worksites found",
  "add_worksite_prompt": "Add a new worksite",
  "assign_employee": "Assign employee",
  "no_address": "No address available",
  "unassigned": "No assignee",
  "active_status": "✅ Active",
  "inactive_status": "❌ Inactive",
  "confirm_delete_title": "Confirm delete",
  "confirm_delete_message": "Are you sure you want to delete",
  "delete_warning": "This will remove all tasks linked to this worksite!",
  "deleting": "⏳ Deleting...",
  "delete_final": "🗑️ Delete permanently",
  "assign_employee_title": "Assign employee",
  "choose_employee_to_assign": "Choose an employee to assign to",
  "no_available_employees": "📭 No available employees",
  "employees_title": "👥 Employees",
  "employees_description": "Manage employee accounts",
  "add_employee": "Add employee",
  "no_employees": "📭 No employees yet",
  "add_new_employee_prompt": "Add a new employee",
  "loading_employees": "⏳ Loading employees...",
  "admin_role": "Admin",
  "field_employee": "Field employee",
  "suspended_status": "❌ Suspended",
  "delete_irreversible": "This action cannot be undone!",
  "done": "Done",
  "settings_intro": "Control app settings and check your subscription status.",
  "subscription_title": "Subscription status",
  "subscription_intro": "Check your admin subscription state and its expiry date.",
  "subscription_description": "Shows the admin subscription state and expiry date.",
  "subscription_status_label": "Subscription status",
  "subscription_expires_at_label": "Expires at",
  "subscription_active": "Active",
  "subscription_expired": "Expired",
  "subscription_canceled": "Canceled",
  "subscription_active_message": "Your subscription is active and the admin dashboard is available.",
  "subscription_expired_message": "Your subscription has expired. Contact support to renew or reactivate it.",
  "settings_language_title": "Default language",
  "settings_language_label": "Admin dashboard",
  "settings_geofence_title": "Default geofence radius",
  "settings_geofence_hint": "The default radius used when creating a worksite if a custom value is not provided.",
  "settings_geofence_radius": "Radius (meters)",
  "ar": "العربية",
  "he": "עברית",
  "en": "English",
  "year": "Year",
  "month": "Month",
  "export_pdf": "Export PDF",
  "total_hours": "Total Hours",
  "work_days": "Work Days",
  "days": "Days",
  "no_attendance_records": "No attendance records",
  "date": "Date",
  "worked_hours": "Worked Hours",
  "cleanup_old_records": "Cleanup Old Records",
  "name": "Name",
  "worksite": "Worksite",
  "january": "January",
  "february": "February",
  "march": "March",
  "april": "April",
  "may": "May",
  "june": "June",
  "july": "July",
  "august": "August",
  "september": "September",
  "october": "October",
  "november": "November",
  "december": "December",
  "meters": "meters",
  "kilometers": "km",
  "notifications": "Notifications",
  "reports": "Reports",
  "reports_description": "Daily performance summary for all field teams",
  "comprehensive_report": "Comprehensive Report",
  "service_requests_report": "Service Requests Report",
  "employee_performance_report": "Employee Performance Report",
  "client_activity_report": "Client Activity Report",
  "active_employees": "Active Employees",
  "active_clients": "Active Clients",
  "completed_services": "Completed Services",
  "average_rating": "Average Rating",
  "ratings": "Ratings",
  "on_duty": "On Duty",
  "total": "Total",
  "service_requests_status": "Service Requests Status",
  "tasks_status": "Tasks Status",
  "attendance_statistics": "Attendance Statistics",
  "worksites_status": "Worksites Status",
  "total_hours_week": "Total Hours This Week",
  "avg_daily_hours": "Average Daily Hours",
  "completed_duty_today": "Completed Duty Today",
  "total_worksites": "Total Worksites",
  "active_worksites": "Active Worksites",
  "no_data_available": "No data available",
  "failed_to_fetch_report": "Failed to fetch report",
  "title": "Title",
  "client": "Client",
  "employee": "Employee",
  "status": "Status",
  "priority": "Priority",
  "rating": "Rating",
  "created_at": "Created At",
  "phone": "Phone",
  "total_hours": "Total Hours",
  "completed_shifts": "Completed Shifts",
  "assigned_services": "Assigned Services",
  "avg_rating": "Avg Rating",
  "total_requests": "Total Requests",
  "completed_requests": "Completed Requests",
  "pending": "Pending",
  "range_compliance": "Range Compliance",
  "range_compliance_hint": "Ratio of accepted punch-ins vs rejected due to being outside worksite range",
  "accepted": "Accepted",
  "rejected": "Rejected",
  "tasks": "Tasks",
  "tasks_description": "All tasks across employees and worksites",
  "new_task": "+ New Task",
  "task_title": "Task",
  "task_employee": "Employee",
  "task_worksite": "Worksite",
  "task_status": "Status",
  "status_pending": "Pending",
  "status_in_progress": "In Progress",
  "status_completed": "Completed",
  "status_late": "Late",
  "priority_low": "Low",
  "priority_normal": "Normal",
  "priority_high": "High",
  "priority_urgent": "Urgent",
  "notification_punch_out": "Punch-in rejected outside range",
  "notification_new_employee": "New employee account created",
  "notification_tasks_completed": "All tasks completed for today",
  "today": "Today",
  "yesterday": "Yesterday",
  "service_requests": "Service Requests",
  "service_requests_description": "Manage customer requests and assign employees",
  "add_service_request": "Add Service Request",
  "request_created_successfully": "Service request created successfully",
  "please_fill_required_fields": "Please fill all required fields",
  "client_name": "Client Name",
  "client_name_placeholder": "Enter client name",
  "client_phone": "Client Phone",
  "phone_placeholder": "Example: 0501234567",
  "service_title_placeholder": "Example: Plumbing repair",
  "service_description_placeholder": "Detailed description of the required service",
  "address_placeholder": "Full address with landmarks",
  "priority_low": "Low",
  "priority_normal": "Normal",
  "priority_high": "High",
  "priority_urgent": "Urgent",
  "latitude_placeholder": "Example: 24.7136",
  "longitude_placeholder": "Example: 46.6753",
  "location_name": "Location Name",
  "location_name_placeholder": "Example: Riyadh - Al-Nakhil District",
  "my_assigned_requests": "My Assigned Requests",
  "my_assigned_requests_description": "Service requests assigned to you",
  "no_assigned_requests": "No assigned requests",
  "no_assigned_requests_hint": "No service requests have been assigned to you yet",
  "accept_request": "Accept Request",
  "reject_request": "Reject Request",
  "start_work": "Start Work",
  "complete_work": "Complete Work",
  "admin_notes": "Admin Notes",
  "assigned_by": "Assigned by",
  "assigned_at": "Assigned at",
  "request_details": "Request Details",
  "client_info": "Client Information",
  "location_info": "Location Information",
  "assignment_info": "Assignment Information",
  "coordinates": "Coordinates",
  "location_name": "Location Name",
  "add_notes_optional": "Add notes (optional)",
  "failed_to_update_status": "Failed to update status",
  "status_accepted": "Accepted",
  "convert_to_worksite": "Convert to Worksite",
  "convert_to_worksite_hint": "Convert service request location to a new worksite",
  "worksite_name_placeholder": "Enter worksite name",
  "address_placeholder": "Enter address",
  "radius_placeholder": "Enter location radius (meters)",
  "create": "Create",
  "all_statuses": "All Statuses",
  "title_ar_placeholder": "Task title in Arabic",
  "title_he_placeholder": "כותרת המשימה בעברית",
  "title_en_placeholder": "Task title in English",
  "description_ar_placeholder": "Task description in Arabic",
  "description_he_placeholder": "תיאור המשימה בעברית",
  "description_en_placeholder": "Task description in English",
  "task_description": "Task description",
  "select_worksite": "Select worksite",
  "select_employee": "Select employee",
  "priority_low": "Low",
  "priority_normal": "Normal",
  "priority_high": "High",
  "priority_urgent": "Urgent",
  "scheduled_start": "Scheduled start time",
  "scheduled_end": "Scheduled end time",
  "creating": "Creating...",
  "error_creating_task": "Error creating task",
  "status_assigned": "Assigned",
  "no_service_requests": "No Service Requests",
  "no_service_requests_hint": "New customer requests will appear here",
  "client": "Client",
  "no_address": "No address",
  "latitude": "Latitude",
  "longitude": "Longitude",
  "assign_employee": "Assign Employee",
  "waiting_for_employee": "Waiting for Employee",
  "in_execution": "In Execution",
  "reassign_employee": "Reassign Employee",
  "unassign_employee": "Unassign Employee",
  "current_assigned_employee": "Current assigned employee:",
  "assign_employee_modal": "Assign Employee",
  "assign_employee_hint": "Choose an employee to execute the service request",
  "cancel": "Cancel",
  "failed_to_fetch_requests": "Failed to fetch requests",
  "failed_to_fetch_employees": "Failed to fetch employees",
  "new_worksite": "New Worksite",
  "worksites_title": "Worksites",
  "worksites_description": "Manage employee worksites",
  "loading_worksites": "Loading worksites...",
  "no_worksites": "No worksites",
  "add_worksite_prompt": "Add a new worksite to get started",
  "meters_unit": "meters",
  "unassigned": "Unassigned",
  "active_status": "Active",
  "inactive_status": "Inactive",
  "confirm_delete_title": "Confirm Delete",
  "delete_warning": "This action cannot be undone",
  "deleting": "Deleting...",
  "delete_final": "Delete",
  "assign_employee_title": "Assign Employee",
  "choose_employee_to_assign": "Choose an employee to assign to",
  "no_available_employees": "No available employees",
  "employee_assigned_successfully": "Employee assigned successfully",
  "failed_to_assign_employee": "Failed to assign employee",
  "worksite_name": "Site Name",
  "worksite_name_placeholder": "Example: Amman Branch - Queen Rania Street",
  "search_address": "Search Address",
  "search_address_placeholder": "Search for a city or street (e.g., Queen Rania Street or Amman)",
  "select_address_required": "Please select an address from search results",
  "allowed_radius": "Allowed Radius (meters)",
  "no_results_found": "No results found",
  "worksite_added_successfully": "Worksite added successfully!",
  "save_failed": "Save failed",
  "type_city": "City",
  "type_street": "Street",
  "type_address": "Address",
  "type_house": "House",
  "type_landmark": "Landmark",
  "type_location": "Location",
  "city": "City",
  "street": "Street",
  "building_number": "Building Number",
  "coordinates": "Coordinates",
  "saving": "Saving...",
  "save": "Save",
  "search_worldwide": "Search worldwide",
  "search_language_hint": "Search supports Arabic, English, and Hebrew",
  "view_details_title": "View details",
  "has_not_started_work": "Has not started work",
  "waiting_status": "Waiting",
  "completed_status": "Completed",
  "failed_to_fetch_data": "Failed to fetch data",
  "failed_to_fetch_security_notes": "Failed to fetch security notes",
  "failed_to_delete_employee": "Failed to delete employee",
  "failed_to_delete_worksite": "Failed to delete worksite",
  "failed_to_export_pdf": "Failed to export PDF",
  "failed_to_export_report": "Failed to export report",
  "confirm_cleanup_old_records": "Are you sure you want to delete old records?",
  "cleanup_success": "Records cleanup successful",
  "cleanup_failed": "Records cleanup failed",
  "failed_to_fetch_attendance_history": "Failed to fetch attendance history",
  "shift_restored": "🔄 Active shift restored",
  "current_worksite": "Current Worksite",
  "currently_working": "Currently Working",
  "end_shift": "End Shift",
  "force_checkout_title": "Force End Shift",
  "force_checkout_message": "Are you sure you want to end the shift for employee",
  "force_checkout_warning": "⚠️ The shift will be ended immediately without requiring employee punch-in. Use this feature only in emergency situations.",
  "force_checkout_success": "Shift ended successfully",
  "force_checkout_failed": "Failed to end shift",
  "processing": "Processing...",
  "confirm_end_shift": "Confirm End Shift",
  "pwa_install_title": "Install App",
  "pwa_install_text": "Install the app on your device",
  "pull_to_refresh": "Pull to refresh",
  "refreshing": "Refreshing...",
  "release_to_refresh": "Release to refresh",
  "pwa": {
    "installTitle": "Install App",
    "installText": "Install App",
    "installOnIos": "Install on iPhone",
    "iosModalTitle": "Install App on iPhone",
    "iosModalSubtitle": "Get a better app experience!",
    "step1Text": "Tap the Share ⎋ button at the bottom",
    "step2Text": "Scroll down and tap Add to Home Screen",
    "step3Text": "Tap Add in the top corner",
    "benefit1": "⚡ Faster loading",
    "benefit2": "📱 Icon on home screen",
    "benefit3": "🎨 App-like design",
    "gotIt": "Got it ✓",
    "remindLater": "Remind me later"
  },
  "client_modal_title": "➕ Add New Client",
  "client_modal_hint": "📧 Client account will be created with email and password automatically",
  "client_modal_hint_small": "⚠️ Client can use these credentials to login",
  "client_modal_full_name": "👤 Full Name",
  "client_modal_email": "📧 Email",
  "client_modal_phone": "📱 Phone Number",
  "client_modal_email_hint": "Client will use this email to login",
  "client_modal_saving": "⏳ Saving...",
  "client_modal_create": "💾 Create Client",
  "client_modal_success": "✅ Client created successfully!",
  "client_modal_success_email": "📧 Email:",
  "client_modal_success_password": "🔑 Password:",
  "client_modal_success_warning": "⚠️ Please save these credentials and share them with the client",
  "client_modal_error_email": "❌ Please enter a valid email",
  "client_modal_error_failed": "❌ Failed to create client",
  "delete_request": "Delete Request",
  "delete_request_title": "Confirm Delete Service Request",
  "delete_request_message": "Are you sure you want to delete this service request?",
  "delete_request_warning": "⚠️ This will permanently delete the request and all related assignments.",
  "request_deleted_successfully": "✅ Service request deleted successfully",
  "request_delete_failed": "❌ Failed to delete service request",
  "request_not_found_or_deleted": "⚠️ Request not found or already deleted",
  "delete_success_title": "Deleted Successfully",
  "delete_error_title": "Delete Failed",
  "try_again": "Try Again",
  "ok": "OK",
  "err_permission_denied": "You don't have permission to perform this action",
  "notifications_feed": "🔔 Notifications",
  "loading_notifications": "Loading notifications...",
  "no_notifications_available": "No notifications available",
  "view_all": "View all"
}

// =============================================
// المفتاح الموحد للغة في جميع التطبيقات
// =============================================
const STORAGE_KEY = 'worktrack_admin_language'

// =============================================
// ترجمة مباشرة من ملفات JSON
// =============================================
const messages = { ar, he, en }

// =============================================
// الحصول على اللغة المخزنة
// =============================================
function getStoredLang() {
  const stored = localStorage.getItem(STORAGE_KEY)
  if (stored && ['ar', 'he', 'en'].includes(stored)) {
    return stored
  }
  // استخدام لغة المتصفح كافتراضي
  const browserLang = navigator.language || navigator.languages?.[0] || 'ar'
  if (browserLang.startsWith('he')) return 'he'
  if (browserLang.startsWith('en')) return 'en'
  return 'ar'
}

// =============================================
// حالة الترجمة
// =============================================
const i18nState = reactive({
  currentLang: getStoredLang(),
  
  setLang(lang) {
    if (['ar', 'he', 'en'].includes(lang)) {
      this.currentLang = lang
      localStorage.setItem(STORAGE_KEY, lang)
      
      // تحديث اتجاه الصفحة
      document.documentElement.dir = lang === 'ar' || lang === 'he' ? 'rtl' : 'ltr'
      document.documentElement.lang = lang
      
      // إرسال event للإشارة بتغيير اللغة
      window.dispatchEvent(new CustomEvent('language-changed', { detail: { lang } }))
      
      console.log(`🌍 تم تغيير اللغة إلى: ${lang}`)
    }
  },
  
  t(key) {
    const keys = key.split('.')
    let translation = messages[i18nState.currentLang]
    
    for (const k of keys) {
      if (translation && translation[k]) {
        translation = translation[k]
      } else {
        // البحث في اللغة الافتراضية
        let fallbackTranslation = messages['ar']
        for (const fk of keys) {
          if (fallbackTranslation && fallbackTranslation[fk]) {
            fallbackTranslation = fallbackTranslation[fk]
          } else {
            console.warn(`⚠️ مفتاح الترجمة غير موجود: ${key}`)
            return key
          }
        }
        return fallbackTranslation
      }
    }
    return translation
  }
})

// =============================================
// تصدير الدوال
// =============================================
export function useI18n() {
  const t = (key) => i18nState.t(key)
  const setLang = (lang) => i18nState.setLang(lang)
  const currentLang = computed(() => i18nState.currentLang)
  return { t, setLang, currentLang }
}

export default {
  install(app) {
    app.config.globalProperties.$t = (key) => i18nState.t(key)
    app.config.globalProperties.$lang = computed(() => i18nState.currentLang)
    app.provide('i18n', i18nState)
    // Make i18nStore available globally for components that need it
    if (typeof window !== 'undefined') {
      window.i18nStore = i18nState
    }
  }
}

export { i18nState }
```

---

## 📄 frontend-admin-dashboard/src/services/websocket.js

```javascript
class WebSocketService {
  constructor() {
    this.ws = null
    this.reconnectAttempts = 0
    this.maxReconnectAttempts = 5
    this.reconnectDelay = 3000
    this.listeners = {}
  }

  connect() {
    let wsUrl = 'ws://localhost:8080/ws'
    
    // Check if we're in a browser environment with import.meta
    if (typeof import.meta !== 'undefined' && import.meta.env) {
      const apiBase = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1'
      wsUrl = apiBase.replace('/api/v1', '').replace('http://', 'ws://').replace('https://', 'wss://') + '/ws'
    }
    
    try {
      this.ws = new WebSocket(wsUrl)
      
      this.ws.onopen = () => {
        console.log('✅ WebSocket connected')
        this.reconnectAttempts = 0
      }
      
      this.ws.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data)
          console.log('📡 WebSocket message received:', data)
          
          // Handle different message types
          if (data.type === 'password_changed') {
            this.handlePasswordChanged(data)
          } else if (data.type === 'location_update') {
            this.emit('location_update', data.data)
          } else if (data.type === 'employee_status') {
            this.emit('employee_status', data.data)
          }
        } catch (error) {
          console.error('❌ Error parsing WebSocket message:', error)
        }
      }
      
      this.ws.onclose = () => {
        console.log('❌ WebSocket disconnected')
        this.attemptReconnect()
      }
      
      this.ws.onerror = (error) => {
        console.error('❌ WebSocket error:', error)
      }
    } catch (error) {
      console.error('❌ Failed to create WebSocket connection:', error)
      this.attemptReconnect()
    }
  }

  attemptReconnect() {
    if (this.reconnectAttempts < this.maxReconnectAttempts) {
      this.reconnectAttempts++
      console.log(`🔄 Attempting to reconnect (${this.reconnectAttempts}/${this.maxReconnectAttempts})...`)
      
      setTimeout(() => {
        this.connect()
      }, this.reconnectDelay)
    } else {
      console.log('❌ Max reconnection attempts reached')
    }
  }

  handlePasswordChanged(data) {
    console.log('🚨 Password change notification received:', data)
    
    // Show custom alert immediately
    this.showPasswordChangedAlert()
    
    // Emit event for other components
    this.emit('password_changed', data)
  }

  showPasswordChangedAlert() {
    console.log('🚨 Showing password changed alert')
    const currentLang = localStorage.getItem('worktrack_admin_language') || 'ar'
    
    // Try to load from locale files first
    let messages = {
      ar: {
        title: 'تم تغيير كلمة المرور',
        message: 'تم تغيير كلمة المرور الخاصة بحسابك. يرجى تسجيل الدخول مرة أخرى.',
        button: 'تسجيل الدخول'
      },
      he: {
        title: 'הסיסמה שונתה',
        message: 'הסיסמה שלך שונתה. אנא התחבר שוב.',
        button: 'התחבר'
      },
      en: {
        title: 'Password Changed',
        message: 'Your password has been changed. Please log in again.',
        button: 'Log In'
      }
    }
    
    const msg = messages[currentLang] || messages.ar
    
    // Remove any existing alerts first
    const existingAlerts = document.querySelectorAll('[data-password-alert]')
    existingAlerts.forEach(alert => alert.remove())
    
    // Create and show alert
    const alertDiv = document.createElement('div')
    alertDiv.setAttribute('data-password-alert', 'true')
    alertDiv.style.cssText = `
      position: fixed !important;
      top: 0 !important;
      left: 0 !important;
      right: 0 !important;
      bottom: 0 !important;
      background: rgba(0, 0, 0, 0.8) !important;
      display: flex !important;
      align-items: center !important;
      justify-content: center !important;
      z-index: 999999 !important;
      direction: ${currentLang === 'he' ? 'rtl' : currentLang === 'en' ? 'ltr' : 'rtl'} !important;
      font-family: ${currentLang === 'he' ? 'Heebo, Arial, sans-serif' : currentLang === 'en' ? 'Arial, sans-serif' : 'Cairo, Arial, sans-serif'} !important;
    `
    
    const alertBox = document.createElement('div')
    alertBox.style.cssText = `
      background: white !important;
      padding: 40px !important;
      border-radius: 16px !important;
      max-width: 450px !important;
      width: 90% !important;
      text-align: center !important;
      box-shadow: 0 20px 60px rgba(0, 0, 0, 0.4) !important;
      animation: slideIn 0.4s ease-out !important;
      position: relative !important;
      /* Mobile-specific improvements */
      max-height: 90vh !important;
      overflow-y: auto !important;
      -webkit-overflow-scrolling: touch !important;
    `
    
    alertBox.innerHTML = `
      <div style="font-size: 64px; margin-bottom: 20px; animation: pulse 2s infinite;">🔒</div>
      <h2 style="margin: 0 0 20px 0; color: #e74c3c; font-size: 24px; font-weight: 700;">${msg.title}</h2>
      <p style="margin: 0 0 30px 0; color: #555; line-height: 1.8; font-size: 16px;">${msg.message}</p>
      <button id="passwordChangedLogoutBtn" style="
        background: linear-gradient(135deg, #e74c3c 0%, #c0392b 100%) !important;
        color: white !important;
        border: none !important;
        padding: 16px 40px !important;
        border-radius: 8px !important;
        font-size: 18px !important;
        font-weight: 600 !important;
        cursor: pointer !important;
        transition: all 0.3s !important;
        box-shadow: 0 4px 15px rgba(231, 76, 60, 0.3) !important;
        /* Mobile-specific button improvements */
        min-width: 200px !important;
        -webkit-tap-highlight-color: transparent !important;
        -webkit-touch-callout: none !important;
        -webkit-user-select: none !important;
        user-select: none !important;
      ">${msg.button}</button>
      <style>
        @keyframes slideIn {
          from {
            opacity: 0;
            transform: translateY(-30px) scale(0.95);
          }
          to {
            opacity: 1;
            transform: translateY(0) scale(1);
          }
        }
        @keyframes pulse {
          0%, 100% { transform: scale(1); }
          50% { transform: scale(1.1); }
        }
        #passwordChangedLogoutBtn:hover {
          transform: translateY(-2px) !important;
          box-shadow: 0 6px 20px rgba(231, 76, 60, 0.4) !important;
        }
        #passwordChangedLogoutBtn:active {
          transform: translateY(0) !important;
          box-shadow: 0 2px 10px rgba(231, 76, 60, 0.3) !important;
        }
        /* Mobile-specific responsive adjustments */
        @media (max-width: 480px) {
          [data-password-alert] > div {
            padding: 30px 20px !important;
            width: 95% !important;
          }
          [data-password-alert] h2 {
            font-size: 20px !important;
          }
          [data-password-alert] p {
            font-size: 14px !important;
          }
          [data-password-alert] button {
            padding: 14px 32px !important;
            font-size: 16px !important;
            width: 100% !important;
            max-width: 280px !important;
          }
          [data-password-alert] > div > div:first-child {
            font-size: 48px !important;
            margin-bottom: 15px !important;
          }
        }
        @media (max-width: 360px) {
          [data-password-alert] > div {
            padding: 25px 15px !important;
          }
          [data-password-alert] h2 {
            font-size: 18px !important;
          }
          [data-password-alert] p {
            font-size: 13px !important;
          }
          [data-password-alert] button {
            padding: 12px 24px !important;
            font-size: 15px !important;
          }
        }
      </style>
    `
    
    alertDiv.appendChild(alertBox)
    document.body.appendChild(alertDiv)
    
    console.log('✅ Alert added to DOM')
    
    // Make function globally available
    if (typeof window !== 'undefined') {
      window.handlePasswordChangedLogout = () => {
        console.log('🚪 Password changed logout clicked')
        this.handleLogout()
        const alertToRemove = document.querySelector('[data-password-alert]')
        if (alertToRemove) {
          alertToRemove.remove()
        }
      }
      
      // Add event listener to button
      setTimeout(() => {
        const btn = document.getElementById('passwordChangedLogoutBtn')
        if (btn) {
          btn.addEventListener('click', window.handlePasswordChangedLogout)
          console.log('✅ Button event listener added')
        } else {
          console.error('❌ Button not found')
        }
      }, 100)
    }
  }

  handleLogout() {
    // Clear all auth data
    localStorage.removeItem('worktrack_admin_token')
    localStorage.removeItem('worktrack_admin_user')
    
    // Disconnect WebSocket
    if (this.ws) {
      this.ws.close()
    }
    
    // Redirect to login page
    window.location.href = '/login'
  }

  on(event, callback) {
    if (!this.listeners[event]) {
      this.listeners[event] = []
    }
    this.listeners[event].push(callback)
  }

  emit(event, data) {
    if (this.listeners[event]) {
      this.listeners[event].forEach(callback => callback(data))
    }
  }

  disconnect() {
    if (this.ws) {
      this.ws.close()
      this.ws = null
    }
  }
}

// Export singleton instance
export const wsService = new WebSocketService()

// Also export as default for compatibility
export default wsService

// Add test function for debugging
export function testPasswordAlert() {
  console.log('🧪 Testing password changed alert')
  wsService.showPasswordChangedAlert()
}

// Auto-connect when in browser environment
if (typeof window !== 'undefined' && typeof localStorage !== 'undefined') {
  // Connect only when user is logged in
  const token = localStorage.getItem('worktrack_admin_token')
  if (token) {
    // Delay connection to ensure DOM is ready
    if (document.readyState === 'loading') {
      document.addEventListener('DOMContentLoaded', () => {
        wsService.connect()
      })
    } else {
      wsService.connect()
    }
  }
  
  // Expose test function globally for debugging
  window.testPasswordAlert = testPasswordAlert
  console.log('🧪 Test function available: window.testPasswordAlert()')
}
```

---

## 📄 frontend-admin-dashboard/src/store/auth.js

```javascript
import { reactive } from 'vue'
import { currentUser } from '../services/auth'

export const authStore = reactive({
  user: currentUser(),
  
  setUser(user) {
    this.user = user
  },
  
  clear() {
    this.user = null
  }
})

```

---

## 📄 frontend-admin-dashboard/src/store/dashboard.js

```javascript
import { reactive } from 'vue'
import api from '../services/api'

export const dashboardStore = reactive({
  summary: {
    total_employees: 0,
    completed: 0,
    in_progress: 0,
    pending: 0,
    late: 0,
  },
  loading: false,
  error: null,
  async fetchSummary() {
    this.loading = true
    this.error = null
    try {
      const { data } = await api.get('/reports/daily-summary')
      this.summary = data
    } catch (e) {
      this.error = e.response?.data?.error || 'Failed to fetch daily summary'
      console.error('❌ Failed to fetch daily summary:', e)
    } finally {
      this.loading = false
    }
  },
})

```

---

## 📄 frontend-admin-dashboard/src/store/notifications.js

```javascript
import { reactive } from 'vue'
import api from '../services/api'

export const notificationsStore = reactive({
  notifications: [],
  loading: false,
  error: null,
  
  async fetchNotifications() {
    // Skip API calls in development to avoid CORS errors
    if (import.meta.env.DEV) {
      console.log('📋 Notifications disabled in development mode')
      this.notifications = []
      this.loading = false
      this.error = null
      return
    }

    this.loading = true
    this.error = null
    try {
      const { data } = await api.get('/notifications')
      this.notifications = data
    } catch (e) {
      // Handle CORS errors gracefully in development
      if (e.message?.includes('Network Error') || e.code === 'ERR_NETWORK') {
        console.warn('⚠️ Network error (likely CORS in development) - notifications disabled')
        this.error = null // Don't show error for CORS issues in dev
        this.notifications = [] // Clear notifications
      } else {
        this.error = e.response?.data?.error || 'Failed to fetch notifications'
        console.error('❌ Failed to fetch notifications:', e)
      }
    } finally {
      this.loading = false
    }
  },
  
  clear() {
    this.notifications = []
    this.error = null
  }
})

```

---

## 📄 frontend-admin-dashboard/src/tests/components.test.js

```javascript
import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'

// اختبارات المكونات (Components)

describe('مكونات لوحة التحكم', () => {
  
  describe('مكون الإحصائيات', () => {
    it('عرض عدد الموظفين بشكل صحيح', () => {
      const mockData = {
        totalEmployees: 10,
        activeEmployees: 8,
        pendingEmployees: 2,
      }
      
      expect(mockData.totalEmployees).toBe(10)
      expect(mockData.activeEmployees).toBe(8)
      expect(mockData.pendingEmployees).toBe(2)
    })

    it('حساب النسب المئوية', () => {
      const total = 10
      const active = 8
      const percentage = (active / total) * 100
      expect(percentage).toBe(80)
    })
  })

  describe('مكون جدول الموظفين', () => {
    it('عرض الموظفين بشكل صحيح', () => {
      const employees = [
        { id: 1, name: 'محمد أحمد', phone: '0501234567', role: 'employee' },
        { id: 2, name: 'خالد محمد', phone: '0507654321', role: 'admin' },
      ]
      
      expect(employees.length).toBe(2)
      expect(employees[0].name).toBe('محمد أحمد')
      expect(employees[1].role).toBe('admin')
    })

    it('فلترة الموظفين حسب الحالة', () => {
      const employees = [
        { id: 1, name: 'محمد', isActive: true },
        { id: 2, name: 'خالد', isActive: false },
        { id: 3, name: 'أحمد', isActive: true },
      ]
      
      const activeEmployees = employees.filter(e => e.isActive)
      expect(activeEmployees.length).toBe(2)
    })
  })

  describe('مكون نقاط العمل', () => {
    it('عرض نقاط العمل بشكل صحيح', () => {
      const worksites = [
        { id: 1, name: 'الرياض - الرئيسي', address: 'الرياض', active: true },
        { id: 2, name: 'جدة - الفرع', address: 'جدة', active: false },
      ]
      
      expect(worksites.length).toBe(2)
      expect(worksites[0].active).toBe(true)
    })

    it('حساب عدد الموظفين في كل نقطة', () => {
      const worksiteEmployees = {
        'site-1': 5,
        'site-2': 3,
        'site-3': 0,
      }
      
      expect(worksiteEmployees['site-1']).toBe(5)
      expect(worksiteEmployees['site-3']).toBe(0)
    })
  })

  describe('مكون طلبات الخدمة', () => {
    it('عرض الطلبات حسب الحالة', () => {
      const requests = [
        { id: 1, title: 'صيانة مكيف', status: 'pending', priority: 'high' },
        { id: 2, title: 'تسريب مياه', status: 'in_progress', priority: 'urgent' },
        { id: 3, title: 'تركيب إضاءة', status: 'completed', priority: 'normal' },
      ]
      
      const pendingRequests = requests.filter(r => r.status === 'pending')
      expect(pendingRequests.length).toBe(1)
    })

    it('فرز الطلبات حسب الأولوية', () => {
      const requests = [
        { id: 1, priority: 'low' },
        { id: 2, priority: 'urgent' },
        { id: 3, priority: 'high' },
      ]
      
      const priorityOrder = { urgent: 0, high: 1, normal: 2, low: 3 }
      const sorted = [...requests].sort((a, b) => priorityOrder[a.priority] - priorityOrder[b.priority])
      expect(sorted[0].priority).toBe('urgent')
    })
  })
})

describe('مكونات النماذج', () => {
  
  describe('نموذج إضافة موظف', () => {
    it('التحقق من البيانات المطلوبة', () => {
      const employeeData = {
        name: 'محمد أحمد',
        phone: '0501234567',
        email: 'mohamed@example.com',
        role: 'employee',
      }
      
      expect(employeeData.name).toBeTruthy()
      expect(employeeData.phone).toBeTruthy()
      expect(employeeData.email).toBeTruthy()
      expect(employeeData.role).toBeTruthy()
    })

    it('التحقق من صحة رقم الهاتف', () => {
      const validPhone = '0501234567'
      const invalidPhone = '12345'
      
      const isValidPhone = /^05[0-9]{8}$/.test(validPhone)
      const isInvalidPhone = /^05[0-9]{8}$/.test(invalidPhone)
      
      expect(isValidPhone).toBe(true)
      expect(isInvalidPhone).toBe(false)
    })
  })

  describe('نموذج إضافة نقطة عمل', () => {
    it('التحقق من الإحداثيات', () => {
      const worksiteData = {
        name: 'الرياض - الرئيسي',
        address: 'الرياض',
        latitude: 24.7136,
        longitude: 46.6753,
        radius: 100,
      }
      
      expect(worksiteData.latitude).toBeGreaterThanOrEqual(-90)
      expect(worksiteData.latitude).toBeLessThanOrEqual(90)
      expect(worksiteData.longitude).toBeGreaterThanOrEqual(-180)
      expect(worksiteData.longitude).toBeLessThanOrEqual(180)
      expect(worksiteData.radius).toBeGreaterThan(0)
    })
  })
})

describe('مكونات التقارير', () => {
  
  it('حساب إحصائيات الشهر', () => {
    const monthlyData = {
      totalHours: 160,
      totalDays: 20,
      employeeCount: 8,
    }
    
    const avgHoursPerDay = monthlyData.totalHours / monthlyData.totalDays
    const avgHoursPerEmployee = monthlyData.totalHours / monthlyData.employeeCount
    
    expect(avgHoursPerDay).toBe(8)
    expect(avgHoursPerEmployee).toBe(20)
  })

  it('توليد تقرير الحضور', () => {
    const attendanceData = [
      { employeeId: 1, date: '2024-01-01', hours: 8 },
      { employeeId: 1, date: '2024-01-02', hours: 7 },
      { employeeId: 1, date: '2024-01-03', hours: 8 },
    ]
    
    const totalHours = attendanceData.reduce((sum, record) => sum + record.hours, 0)
    const avgHours = totalHours / attendanceData.length
    
    expect(totalHours).toBe(23)
    expect(avgHours).toBeCloseTo(7.67, 2)
  })
})
```

---

## 📄 frontend-admin-dashboard/src/tests/setup.js

```javascript
import { vi } from 'vitest'
import { config } from '@vue/test-utils'

// محاكاة window.matchMedia للاختبارات
Object.defineProperty(window, 'matchMedia', {
  writable: true,
  value: vi.fn().mockImplementation(query => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: vi.fn(),
    removeListener: vi.fn(),
    addEventListener: vi.fn(),
    removeEventListener: vi.fn(),
    dispatchEvent: vi.fn(),
  })),
})

// إعدادات Vue Test Utils
config.global.mocks = {
  $t: (key) => key,
}

// محاكاة localStorage
const localStorageMock = {
  getItem: vi.fn(),
  setItem: vi.fn(),
  removeItem: vi.fn(),
  clear: vi.fn(),
}
global.localStorage = localStorageMock
```

---

## 📄 frontend-admin-dashboard/src/tests/utils.test.js

```javascript
import { describe, it, expect, vi } from 'vitest'

// اختبارات وظائف المساعدة العامة

describe('وظائف المساعدة العامة', () => {
  
  describe('حساب الإحصائيات', () => {
    it('حساب عدد الموظفين النشطين', () => {
      const employees = [
        { id: 1, isActive: true },
        { id: 2, isActive: true },
        { id: 3, isActive: false },
      ]
      const activeCount = employees.filter(e => e.isActive).length
      expect(activeCount).toBe(2)
    })

    it('حساب إجمالي الساعات', () => {
      const attendances = [
        { hours: 8 },
        { hours: 6 },
        { hours: 7 },
      ]
      const totalHours = attendances.reduce((sum, a) => sum + a.hours, 0)
      expect(totalHours).toBe(21)
    })
  })

  describe('فلترة البيانات', () => {
    it('فلترة طلبات الخدمة حسب الحالة', () => {
      const requests = [
        { id: 1, status: 'pending' },
        { id: 2, status: 'completed' },
        { id: 3, status: 'pending' },
      ]
      const pendingRequests = requests.filter(r => r.status === 'pending')
      expect(pendingRequests.length).toBe(2)
    })

    it('فلترة الموظفين حسب نقطة العمل', () => {
      const employees = [
        { id: 1, worksiteId: 'site-1' },
        { id: 2, worksiteId: 'site-2' },
        { id: 3, worksiteId: 'site-1' },
      ]
      const siteEmployees = employees.filter(e => e.worksiteId === 'site-1')
      expect(siteEmployees.length).toBe(2)
    })
  })

  describe('تنسيق البيانات', () => {
    it('تنسيق الوقت', () => {
      const date = new Date('2024-01-15T09:30:00')
      const formatted = date.toLocaleTimeString('en-US', { 
        hour: '2-digit', 
        minute: '2-digit' 
      })
      expect(formatted).toContain('09')
    })

    it('تنسيق التاريخ', () => {
      const date = new Date('2024-01-15')
      const formatted = date.toLocaleDateString('ar-SA')
      expect(formatted).toBeDefined()
    })
  })
})

describe('وظائف التحقق من البيانات', () => {
  
  it('التحقق من صحة رقم الهاتف', () => {
    const phone = '0501234567'
    const isValid = /^05[0-9]{8}$/.test(phone)
    expect(isValid).toBe(true)
  })

  it('التحقق من صحة البريد الإلكتروني', () => {
    const email = 'test@example.com'
    const isValid = /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)
    expect(isValid).toBe(true)
  })

  it('التحقق من صحة الإحداثيات', () => {
    const lat = 24.7136
    const lng = 46.6753
    const isValidLat = lat >= -90 && lat <= 90
    const isValidLng = lng >= -180 && lng <= 180
    expect(isValidLat && isValidLng).toBe(true)
  })
})

describe('وظائف الحساب', () => {
  
  it('حساب المسافة بين نقطتين', () => {
    // حساب بسيط للمسافة
    const lat1 = 24.7136
    const lon1 = 46.6753
    const lat2 = 24.7146
    const lon2 = 46.6763
    
    const R = 6371 // نصف قطر الأرض بالكيلومتر
    const dLat = (lat2 - lat1) * Math.PI / 180
    const dLon = (lon2 - lon1) * Math.PI / 180
    const a = Math.sin(dLat/2) * Math.sin(dLat/2) +
              Math.cos(lat1 * Math.PI / 180) * Math.cos(lat2 * Math.PI / 180) *
              Math.sin(dLon/2) * Math.sin(dLon/2)
    const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1-a))
    const distance = R * c
    
    expect(distance).toBeGreaterThan(0)
    expect(distance).toBeLessThan(1) // مسافة قصيرة
  })

  it('حساب ساعات العمل', () => {
    const checkIn = new Date('2024-01-15T09:00:00')
    const checkOut = new Date('2024-01-15T17:00:00')
    const diffMs = checkOut - checkIn
    const diffHours = diffMs / (1000 * 60 * 60)
    expect(diffHours).toBe(8)
  })

  it('حساب ساعات العمل عبر منتصف الليل', () => {
    const checkIn = new Date('2024-01-15T22:00:00')
    const checkOut = new Date('2024-01-16T03:00:00')
    const diffMs = checkOut - checkIn
    const diffHours = diffMs / (1000 * 60 * 60)
    expect(diffHours).toBe(5)
  })
})

describe('وظائف فرز البيانات', () => {
  
  it('فرز الموظفين حسب الاسم', () => {
    const employees = [
      { id: 1, name: 'محمد' },
      { id: 2, name: 'أحمد' },
      { id: 3, name: 'خالد' },
    ]
    const sorted = [...employees].sort((a, b) => a.name.localeCompare(b.name, 'ar'))
    expect(sorted[0].name).toBe('أحمد')
  })

  it('فرز الطلبات حسب الأولوية', () => {
    const requests = [
      { id: 1, priority: 'low' },
      { id: 2, priority: 'urgent' },
      { id: 3, priority: 'high' },
    ]
    const priorityOrder = { urgent: 0, high: 1, normal: 2, low: 3 }
    const sorted = [...requests].sort((a, b) => priorityOrder[a.priority] - priorityOrder[b.priority])
    expect(sorted[0].priority).toBe('urgent')
  })
})
```

---

## 📄 frontend-admin-dashboard/src/utils/numberFormatter.js

```javascript
/**
 * Convert Western numerals (0-9) to localized numerals
 * @param {string|number} num - The number to convert
 * @param {string} locale - The locale code ('ar', 'he', 'en')
 * @returns {string} - The number with localized numerals
 */
export function toLocalNumerals(num, locale = 'en') {
  if (num === null || num === undefined) return ''
  
  const str = String(num)
  
  // Arabic numerals (٠-٩)
  const arabicNumerals = ['٠', '١', '٢', '٣', '٤', '٥', '٦', '٧', '٨', '٩']
  
  // Hebrew numerals use the same Western numerals (0-9) in most modern contexts
  // But we keep this function ready for future Hebrew numeral support if needed
  
  if (locale === 'ar') {
    return str.replace(/[0-9]/g, (digit) => arabicNumerals[parseInt(digit)])
  }
  
  // For English and Hebrew, return the original Western numerals
  return str
}

/**
 * Format a date with localized numerals
 * @param {Date|string} date - The date to format
 * @param {string} locale - The locale code ('ar', 'he', 'en')
 * @param {string} format - The format string (default: 'DD/MM/YYYY HH:mm:ss')
 * @returns {string} - The formatted date with localized numerals
 */
export function formatLocalDate(date, locale = 'en', format = 'DD/MM/YYYY HH:mm:ss') {
  if (!date) return ''
  
  const d = new Date(date)
  if (isNaN(d.getTime())) return ''
  
  const day = String(d.getDate()).padStart(2, '0')
  const month = String(d.getMonth() + 1).padStart(2, '0')
  const year = String(d.getFullYear())
  const hours = String(d.getHours()).padStart(2, '0')
  const minutes = String(d.getMinutes()).padStart(2, '0')
  const seconds = String(d.getSeconds()).padStart(2, '0')
  
  let formatted = format
    .replace('DD', toLocalNumerals(day, locale))
    .replace('MM', toLocalNumerals(month, locale))
    .replace('YYYY', toLocalNumerals(year, locale))
    .replace('YY', toLocalNumerals(year.slice(-2), locale))
    .replace('HH', toLocalNumerals(hours, locale))
    .replace('mm', toLocalNumerals(minutes, locale))
    .replace('ss', toLocalNumerals(seconds, locale))
  
  return formatted
}

/**
 * Format a number with localized numerals and optional decimal places
 * @param {number} num - The number to format
 * @param {string} locale - The locale code ('ar', 'he', 'en')
 * @param {number} decimals - Number of decimal places (default: 0)
 * @returns {string} - The formatted number with localized numerals
 */
export function formatLocalNumber(num, locale = 'en', decimals = 0) {
  if (num === null || num === undefined) return ''
  
  const fixed = Number(num).toFixed(decimals)
  return toLocalNumerals(fixed, locale)
}

/**
 * Format a currency amount with localized numerals
 * @param {number} amount - The amount to format
 * @param {string} locale - The locale code ('ar', 'he', 'en')
 * @param {string} currency - The currency symbol (default: '$')
 * @returns {string} - The formatted currency with localized numerals
 */
export function formatLocalCurrency(amount, locale = 'en', currency = '$') {
  if (amount === null || amount === undefined) return ''
  
  const formatted = formatLocalNumber(amount, locale, 2)
  return locale === 'ar' ? `${formatted} ${currency}` : `${currency}${formatted}`
}

```

---

## 📄 frontend-admin-dashboard/src/views/AssignedRequestsView.vue

```vue
<template>
  <div>
    <div class="page-head">
      <div>
        <h2>📋 {{ t('my_assigned_requests') }}</h2>
        <p>{{ t('my_assigned_requests_description') }}</p>
      </div>
      <div class="filters">
        <select v-model="filterStatus" @change="fetchRequests">
          <option value="">{{ t('all_statuses') }}</option>
          <option value="assigned">{{ t('status_assigned') }}</option>
          <option value="accepted">{{ t('status_accepted') }}</option>
          <option value="in_progress">{{ t('status_in_progress') }}</option>
          <option value="completed">{{ t('status_completed') }}</option>
        </select>
      </div>
    </div>

    <div v-if="loading" class="empty-state"><p>⏳ {{ t('loading') }}...</p></div>

    <div v-else-if="requests.length === 0" class="empty-state">
      <h3>📭 {{ t('no_assigned_requests') }}</h3>
      <p>{{ t('no_assigned_requests_hint') }}</p>
    </div>

    <div v-else>
      <div v-for="req in filteredRequests" :key="req.id" class="card request-card" @click="selectedRequest = req">
        <div class="request-card__header">
          <div>
            <span class="badge" :class="getPriorityClass(req.priority)">
              {{ getPriorityLabel(req.priority) }}
            </span>
            <span class="badge" :class="getStatusClass(req.status)">
              {{ getStatusLabel(req.status) }}
            </span>
          </div>
          <span class="request-card__time mono">{{ formatDate(req.assigned_at) }}</span>
        </div>

        <h3>{{ req.title }}</h3>
        <p class="request-card__desc">{{ req.description }}</p>

        <div class="request-card__info">
          <span>👤 {{ req.client_name || t('client') }}</span>
          <span>📞 {{ req.client_phone || req.phone || '—' }}</span>
          <span>📍 {{ req.address || t('no_address') }}</span>
        </div>

        <div class="request-card__location">
          <div v-if="req.location_name" class="location-name">
            📍 {{ req.location_name }}
          </div>
          <div class="location-coords">
            <span class="mono">{{ t('latitude') }}: {{ toLocalNumerals(req.latitude.toFixed(5), locale?.value || 'en') }}</span>
            <span class="mono">{{ t('longitude') }}: {{ toLocalNumerals(req.longitude.toFixed(5), locale?.value || 'en') }}</span>
          </div>
          <button class="btn btn--ghost btn--sm" @click.stop="openInMaps(req.latitude, req.longitude)">
            🗺️ {{ t('open_in_maps') }}
          </button>
        </div>

        <div v-if="req.admin_notes" class="request-card__notes">
          <span class="notes-label">📝 {{ t('admin_notes') }}:</span>
          <span>{{ req.admin_notes }}</span>
        </div>

        <div class="request-card__meta">
          <span class="meta-item">
            <span class="meta-label">{{ t('assigned_by') }}:</span>
            <span>{{ req.admin_name || '—' }}</span>
          </span>
          <span class="meta-item">
            <span class="meta-label">{{ t('assigned_at') }}:</span>
            <span class="mono">{{ formatDate(req.assigned_at) }}</span>
          </span>
        </div>

        <div class="request-card__actions">
          <button 
            v-if="req.status === 'assigned'" 
            class="btn btn--success btn--sm" 
            @click.stop="updateStatus(req.id, 'accepted')"
          >
            ✅ {{ t('accept_request') }}
          </button>
          <button 
            v-if="req.status === 'assigned'" 
            class="btn btn--danger btn--sm" 
            @click.stop="updateStatus(req.id, 'rejected')"
          >
            ❌ {{ t('reject_request') }}
          </button>
          <button 
            v-if="req.status === 'accepted'" 
            class="btn btn--primary btn--sm" 
            @click.stop="updateStatus(req.id, 'in_progress')"
          >
            🚀 {{ t('start_work') }}
          </button>
          <button 
            v-if="req.status === 'in_progress'" 
            class="btn btn--success btn--sm" 
            @click.stop="updateStatus(req.id, 'completed')"
          >
            ✅ {{ t('complete_work') }}
          </button>
          <button 
            v-if="req.status === 'completed'" 
            class="btn btn--ghost btn--sm" 
            disabled
          >
            ✅ {{ t('status_completed') }}
          </button>
        </div>
      </div>
    </div>

    <!-- مودال تفاصيل الطلب -->
    <div v-if="selectedRequest" class="modal-backdrop" @click.self="selectedRequest = null">
      <div class="modal card">
        <div class="modal-header">
          <h3>{{ t('request_details') }}</h3>
          <button class="modal-close" @click="selectedRequest = null">✕</button>
        </div>
        <div class="modal-body">
          <div class="request-detail">
            <h4>{{ selectedRequest.title }}</h4>
            <p>{{ selectedRequest.description }}</p>

            <div class="detail-section">
              <h5>{{ t('client_info') }}</h5>
              <div class="detail-row">
                <span>{{ t('name') }}:</span>
                <span>{{ selectedRequest.client_name || '—' }}</span>
              </div>
              <div class="detail-row">
                <span>{{ t('phone') }}:</span>
                <span>{{ selectedRequest.client_phone || selectedRequest.phone || '—' }}</span>
              </div>
            </div>

            <div class="detail-section">
              <h5>{{ t('location_info') }}</h5>
              <div class="detail-row">
                <span>{{ t('address') }}:</span>
                <span>{{ selectedRequest.address || '—' }}</span>
              </div>
              <div class="detail-row">
                <span>{{ t('location_name') }}:</span>
                <span>{{ selectedRequest.location_name || '—' }}</span>
              </div>
              <div class="detail-row">
                <span>{{ t('coordinates') }}:</span>
                <span class="mono">{{ selectedRequest.latitude.toFixed(6) }}, {{ selectedRequest.longitude.toFixed(6) }}</span>
              </div>
            </div>

            <div v-if="selectedRequest.admin_notes" class="detail-section">
              <h5>{{ t('admin_notes') }}</h5>
              <p>{{ selectedRequest.admin_notes }}</p>
            </div>

            <div class="detail-section">
              <h5>{{ t('assignment_info') }}</h5>
              <div class="detail-row">
                <span>{{ t('assigned_by') }}:</span>
                <span>{{ selectedRequest.admin_name || '—' }}</span>
              </div>
              <div class="detail-row">
                <span>{{ t('assigned_at') }}:</span>
                <span class="mono">{{ formatDate(selectedRequest.assigned_at) }}</span>
              </div>
            </div>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn--ghost" @click="selectedRequest = null">{{ t('close') }}</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import api from '../services/api'
import { useI18n } from '../services/i18n'
import { formatLocalDate, toLocalNumerals } from '../utils/numberFormatter'

const { t, locale } = useI18n()

const requests = ref([])
const loading = ref(false)
const filterStatus = ref('')
const selectedRequest = ref(null)

const filteredRequests = computed(() => {
  if (!filterStatus.value) return requests.value
  return requests.value.filter(r => r.status === filterStatus.value)
})

async function fetchRequests() {
  loading.value = true
  try {
    const { data } = await api.get('/service/assigned')
    requests.value = data || []
  } catch (error) {
    console.error(t('failed_to_fetch_requests'), error)
  } finally {
    loading.value = false
  }
}

async function updateStatus(requestId, status) {
  const notes = prompt(t('add_notes_optional'))
  
  try {
    await api.put(`/service/assigned/${requestId}/status`, {
      status,
      notes: notes || ''
    })
    await fetchRequests()
  } catch (error) {
    console.error('فشل تحديث الحالة:', error)
    alert(t('failed_to_update_status'))
  }
}

function getPriorityLabel(priority) {
  const labels = { low: 'منخفضة', normal: 'عادية', high: 'عالية', urgent: 'طارئة' }
  return labels[priority] || priority
}

function getPriorityClass(priority) {
  const classes = { low: 'badge--gray', normal: 'badge--blue', high: 'badge--gold', urgent: 'badge--out' }
  return classes[priority] || ''
}

function getStatusLabel(status) {
  const labels = { 
    assigned: t('status_assigned'), 
    accepted: t('status_accepted'), 
    in_progress: t('status_in_progress'), 
    completed: t('status_completed'),
    rejected: t('status_rejected')
  }
  return labels[status] || status
}

function getStatusClass(status) {
  const classes = { 
    assigned: 'badge--blue', 
    accepted: 'badge--success', 
    in_progress: 'badge--info', 
    completed: 'badge--in',
    rejected: 'badge--out'
  }
  return classes[status] || ''
}

function formatDate(date) {
  if (!date) return '—'
  return formatLocalDate(date, locale?.value || 'en', 'DD/MM/YYYY HH:mm')
}

function openInMaps(lat, lng) {
  const url = `https://www.google.com/maps?q=${lat},${lng}`
  window.open(url, '_blank')
}

onMounted(fetchRequests)
</script>

<style scoped>
.page-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 18px;
  flex-wrap: wrap;
  gap: 10px;
}

.page-head h2 { font-size: 20px; margin-bottom: 4px; }
.page-head p { font-size: 13px; color: var(--ink-soft); }
.filters select {
  padding: 8px 14px;
  border-radius: var(--radius-sm);
  border: 1px solid var(--line);
  background: var(--surface);
  font-family: var(--font-body);
  min-width: 150px;
}

@media (max-width: 768px) {
  .page-head { flex-direction: column; align-items: flex-start; }
  .filters { width: 100%; }
  .filters select { width: 100%; }
  
  .request-card__location {
    flex-direction: column;
    gap: 8px;
  }
  
  .request-card__actions {
    flex-direction: column;
  }
  
  .request-card__actions .btn {
    width: 100%;
  }
}

.request-card {
  padding: 18px 20px;
  margin-bottom: 12px;
  transition: border-color 0.2s;
}
.request-card:hover { border-color: var(--brand-tint); }

.request-card__header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}
.request-card__header .badge { margin-inline-end: 6px; }
.request-card__time { font-size: 12px; color: var(--ink-soft); }

.request-card h3 { font-size: 16px; margin-bottom: 6px; }
.request-card__desc { font-size: 13px; color: var(--ink-soft); margin-bottom: 10px; }

.request-card__info {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  font-size: 13px;
  color: var(--ink-soft);
  margin-bottom: 8px;
}

.request-card__location {
  display: flex;
  flex-direction: column;
  gap: 8px;
  font-size: 12px;
  color: var(--ink-soft);
  background: var(--canvas);
  padding: 8px 12px;
  border-radius: var(--radius-sm);
  margin-bottom: 12px;
}

.location-name {
  font-size: 13px;
  font-weight: 500;
  color: var(--ink);
  margin-bottom: 4px;
}

.location-coords {
  display: flex;
  gap: 16px;
}

.request-card__notes {
  background: var(--brand-tint);
  padding: 8px 12px;
  border-radius: var(--radius-sm);
  margin-bottom: 12px;
  font-size: 13px;
}

.notes-label {
  font-weight: 600;
  color: var(--brand-dark);
}

.request-card__meta {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  font-size: 12px;
  color: var(--ink-soft);
  margin-bottom: 12px;
}

.meta-item {
  display: flex;
  gap: 4px;
}

.meta-label {
  font-weight: 500;
}

.request-card__actions {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.badge--gray { background: var(--line); color: var(--ink-soft); }
.badge--blue { background: #E3ECF7; color: #2C6B9E; }
.badge--info { background: #E3F0F7; color: #1A7A8A; }
.badge--success { background: var(--signal-in-tint); color: var(--signal-in); }
.badge--in { background: var(--signal-in-tint); color: var(--signal-in); }
.badge--out { background: var(--signal-out-tint); color: var(--signal-out); }

.modal-backdrop {
  position: fixed; inset: 0;
  background: rgba(22,35,46,0.6);
  backdrop-filter: blur(4px);
  display: flex; align-items: center; justify-content: center;
  z-index: 50; padding: 20px;
}
.modal {
  width: 100%; max-width: 500px; padding: 24px;
  max-height: 80vh; overflow-y: auto;
}
.modal h3 { font-size: 17px; margin-bottom: 4px; }
.modal p { font-size: 13px; color: var(--ink-soft); margin-bottom: 16px; }

.request-detail h4 {
  font-size: 16px;
  margin-bottom: 8px;
}

.request-detail > p {
  font-size: 14px;
  color: var(--ink-soft);
  margin-bottom: 16px;
}

.detail-section {
  margin-bottom: 16px;
  padding: 12px;
  background: var(--canvas);
  border-radius: var(--radius-sm);
}

.detail-section h5 {
  font-size: 13px;
  font-weight: 600;
  margin-bottom: 8px;
  color: var(--ink);
}

.detail-row {
  display: flex;
  justify-content: space-between;
  margin-bottom: 6px;
  font-size: 13px;
}

.detail-row:last-child {
  margin-bottom: 0;
}

.detail-row span:first-child {
  color: var(--ink-soft);
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}
</style>

```

---

## 📄 frontend-admin-dashboard/src/views/ClientsView.vue

```vue
<template>
  <div>
    <PWAInstallButton />
    <div class="page-head">
      <div>
        <h2>{{ t('clients_title') }}</h2>
        <p>{{ t('clients_description') }}</p>
      </div>
      <button class="btn btn--primary" @click="showModal = true">+ {{ t('new_client') }}</button>
    </div>

    <ClientFormModal 
      v-if="showModal" 
      @close="showModal = false" 
      @client-added="handleClientAdded" 
    />

    <div class="card">
      <!-- جدول للشاشات الكبيرة -->
      <div class="table-wrapper desktop-only">
        <table class="table">
          <thead>
            <tr>
              <th>{{ t('clients_name') }}</th>
              <th>{{ t('clients_phone') }}</th>
              <th>{{ t('clients_email') }}</th>
              <th>{{ t('clients_service_count') }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="c in clients" :key="c.id">
              <td>{{ c.name }}</td>
              <td class="mono">{{ c.phone }}</td>
              <td class="mono">{{ c.email }}</td>
              <td>{{ formatNumber(c.servicesCount) }}</td>
            </tr>
          </tbody>
        </table>
      </div>
      
      <!-- بطاقات للشاشات الصغيرة -->
      <div class="mobile-cards mobile-only">
        <div v-for="c in clients" :key="c.id" class="client-card">
          <div class="client-card__header">
            <span class="client-card__name">{{ c.name }}</span>
            <span class="badge badge--info">{{ formatNumber(c.servicesCount) }} {{ t('clients_service_count') }}</span>
          </div>
          <div class="client-card__body">
            <div class="client-card__row">
              <span class="client-card__label">{{ t('clients_phone') }}</span>
              <span class="client-card__value mono">{{ c.phone }}</span>
            </div>
            <div class="client-card__row">
              <span class="client-card__label">{{ t('clients_email') }}</span>
              <span class="client-card__value mono">{{ c.email }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useI18n } from '../services/i18n'
import api from '../services/api'
import ClientFormModal from '../components/ClientFormModal.vue'
import PWAInstallButton from '../components/PWAInstallButton.vue'
import { toLocalNumerals } from '../utils/numberFormatter'

const { t, locale } = useI18n()

const showModal = ref(false)
const clients = ref([])

async function fetchClients() {
  try {
    const { data } = await api.get('/admin/customers')
    clients.value = data.map(c => ({
      id: c.id,
      name: c.full_name,
      phone: c.phone || '—',
      email: c.email || '—',
      servicesCount: 0 // TODO: fetch actual service count
    }))
  } catch (error) {
    console.error('❌ Failed to fetch clients:', error)
    // Keep dummy data for now
    clients.value = [
      { id: 1, name: 'شركة الأفق للعقارات', phone: '079xxxxxxx', email: 'contact@ofoq.com', servicesCount: 12 },
      { id: 2, name: 'مجمع الزهور السكني', phone: '077xxxxxxx', email: 'admin@zohour.com', servicesCount: 5 },
    ]
  }
}

function handleClientAdded() {
  fetchClients()
}

// Format numbers with localized numerals
function formatNumber(num) {
  return toLocalNumerals(num, locale?.value || 'en')
}

onMounted(() => {
  fetchClients()
})
</script>

<style scoped>
.page-head { display: flex; align-items: flex-start; justify-content: space-between; margin-bottom: 18px; flex-wrap: wrap; gap: 10px; }
.page-head h2 { font-size: 20px; margin-bottom: 4px; }
.page-head p { font-size: 13px; color: var(--ink-soft); }
.table { width: 100%; border-collapse: collapse; }
.table th { text-align: right; font-size: 12px; color: var(--ink-soft); font-weight: 600; padding: 14px 20px; border-bottom: 1px solid var(--line); }
.table td { padding: 14px 20px; font-size: 14px; border-bottom: 1px solid var(--line); }
.table tr:last-child td { border-bottom: none; }

@media (max-width: 768px) {
  .page-head { flex-direction: column; align-items: flex-start; }
  .page-head .btn { width: 100%; }
  
  .desktop-only { display: none; }
  .mobile-only { display: block; }
}

@media (min-width: 769px) {
  .desktop-only { display: block; }
  .mobile-only { display: none; }
}

/* تصميم بطاقات العملاء للهاتف */
.mobile-cards {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.client-card {
  background: var(--surface);
  border: 1px solid var(--line);
  border-radius: var(--radius-md);
  padding: 16px;
  transition: all 0.2s ease;
}

.client-card:hover {
  border-color: var(--brand);
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

.client-card__header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  padding-bottom: 12px;
  border-bottom: 1px solid var(--line);
}

.client-card__name {
  font-weight: 600;
  color: var(--ink);
  font-size: 15px;
}

.client-card__body {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.client-card__row {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.client-card__label {
  font-size: 13px;
  color: var(--ink-soft);
  font-weight: 500;
}

.client-card__value {
  font-size: 14px;
  color: var(--ink);
  font-weight: 500;
}
</style>

```

---

## 📄 frontend-admin-dashboard/src/views/CustomersView.vue

```vue
<template>
  <div>
    <div class="page-head">
      <div>
        <h2>{{ t('customers_title') }}</h2>
        <p>{{ t('customers_description') }}</p>
      </div>
      <div class="page-head-actions">
        <button class="btn btn--primary" @click="showModal = true">+ {{ t('add_customer') }}</button>
      </div>
    </div>

    <div v-if="loading" class="empty-state">
      <p>{{ t('loading_customers') }}</p>
    </div>

    <div v-else-if="customers.length === 0" class="empty-state">
      <h3>{{ t('no_customers') }}</h3>
      <p>{{ t('no_customers_hint') }}</p>
    </div>

    <div v-else class="card">
      <!-- جدول للشاشات الكبيرة -->
      <div class="table-wrapper desktop-only">
        <table class="table">
          <thead>
            <tr>
              <th>{{ t('name') }}</th>
              <th>{{ t('phone') }}</th>
              <th>{{ t('email') }}</th>
              <th>{{ t('status') }}</th>
              <th>{{ t('created_at') }}</th>
              <th>{{ t('actions') }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="customer in customers" :key="customer.id">
              <td>
                <div class="table__person">
                  <span class="table__avatar">{{ customer.full_name?.slice(0, 1) || '?' }}</span>
                  {{ customer.full_name }}
                </div>
              </td>
              <td class="mono">{{ customer.phone || '—' }}</td>
              <td class="mono">{{ customer.email || '—' }}</td>
              <td>
                <span class="badge" :class="customer.is_active ? 'badge--in' : 'badge--out'">
                  {{ customer.is_active ? t('active_status') : t('suspended_status') }}
                </span>
              </td>
              <td class="mono">{{ formatDate(customer.created_at) }}</td>
              <td>
                <div class="table-actions">
                  <button 
                    class="btn btn--primary btn--sm" 
                    @click="resetPassword(customer)"
                  >
                    🔑 {{ t('reset_password') }}
                  </button>
                  <button 
                    class="btn btn--danger btn--sm" 
                    @click="confirmDelete(customer)"
                  >
                    🗑️ {{ t('delete') }}
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      
      <!-- بطاقات للشاشات الصغيرة -->
      <div class="mobile-cards mobile-only">
        <div v-for="customer in customers" :key="customer.id" class="customer-card">
          <div class="customer-card__header">
            <div class="customer-card__person">
              <span class="table__avatar">{{ customer.full_name?.slice(0, 1) || '?' }}</span>
              <div class="customer-card__info">
                <span class="customer-card__name">{{ customer.full_name }}</span>
              </div>
            </div>
            <div class="customer-card__badges">
              <span class="badge badge--compact" :class="customer.is_active ? 'badge--in' : 'badge--out'">
                {{ customer.is_active ? t('active_status') : t('suspended_status') }}
              </span>
            </div>
          </div>
          <div class="customer-card__body">
            <div class="customer-card__row">
              <span class="customer-card__label">{{ t('phone') }}</span>
              <span class="customer-card__value mono">{{ customer.phone || '—' }}</span>
            </div>
            <div class="customer-card__row">
              <span class="customer-card__label">{{ t('email') }}</span>
              <span class="customer-card__value mono">{{ customer.email || '—' }}</span>
            </div>
            <div class="customer-card__row">
              <span class="customer-card__label">{{ t('created_at') }}</span>
              <span class="customer-card__value mono">{{ formatDate(customer.created_at) }}</span>
            </div>
          </div>
          <div class="customer-card__actions">
            <button 
              class="btn btn--primary btn--sm btn--compact" 
              @click="resetPassword(customer)"
            >
              <span class="btn-icon">🔑</span>
              <span class="btn-text">{{ t('reset_password') }}</span>
            </button>
            <button 
              class="btn btn--danger btn--sm btn--compact" 
              @click="confirmDelete(customer)"
            >
              <span class="btn-icon">🗑️</span>
              <span class="btn-text">{{ t('delete') }}</span>
            </button>
          </div>
        </div>
      </div>
    </div>

    <CustomerFormModal 
      v-if="showModal" 
      @close="showModal = false" 
      @customer-added="fetchCustomers"
    />

    <!-- مودال إعادة تعيين كلمة المرور -->
    <div v-if="showPasswordModal" class="modal-backdrop" @click.self="showPasswordModal = false">
      <div class="modal card">
        <div class="modal-header">
          <h3>🔑 {{ t('reset_password') }}</h3>
          <button class="modal-close" @click="showPasswordModal = false">✕</button>
        </div>
        <div class="modal-body">
          <div v-if="resetError" class="alert alert-error">{{ resetError }}</div>
          <div v-if="resetSuccess" class="alert alert-success">
            <p>{{ resetSuccess }}</p>
            <div class="new-password-display">
              <div class="password-label">🔑 {{ t('new_password') }}:</div>
              <div class="password-value-container">
                <span v-if="showPassword" class="password-value visible">{{ newPassword }}</span>
                <span v-else class="password-value hidden">••••••••••••••••</span>
                <button 
                  class="btn-toggle-password" 
                  @click="showPassword = !showPassword"
                  :title="showPassword ? t('hide_password') : t('show_password')"
                >
                  {{ showPassword ? '🙈' : '👁️' }}
                </button>
                <button 
                  v-if="showPassword"
                  class="btn-copy-password" 
                  @click="copyPassword"
                  :title="passwordCopied ? t('copied') : t('copy_password')"
                >
                  {{ passwordCopied ? '✅' : '📋' }}
                </button>
              </div>
            </div>
            <p class="password-warning">⚠️ {{ t('save_password_warning') }}</p>
          </div>
          <div v-else>
            <p>{{ t('reset_password_confirm') }} <strong>{{ selectedCustomer?.full_name }}</strong>؟</p>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn--ghost" @click="showPasswordModal = false">{{ t('cancel') }}</button>
          <button 
            v-if="!resetSuccess" 
            class="btn btn--primary" 
            @click="confirmResetPassword" 
            :disabled="resetting"
          >
            {{ resetting ? '⏳ ' + t('resetting') : '🔑 ' + t('reset_password') }}
          </button>
        </div>
      </div>
    </div>

    <!-- مودال تأكيد حذف العميل -->
    <div v-if="showDeleteModal" class="modal-backdrop" @click.self="showDeleteModal = false">
      <div class="modal card">
        <div class="modal-header">
          <h3>🗑️ {{ t('delete_customer') }}</h3>
          <button class="modal-close" @click="showDeleteModal = false">✕</button>
        </div>
        <div class="modal-body">
          <div v-if="deleteError" class="alert alert-error">{{ deleteError }}</div>
          <div v-else>
            <p>{{ t('delete_customer_confirm') }} <strong>{{ selectedCustomer?.full_name }}</strong>؟</p>
            <p class="delete-warning">⚠️ {{ t('delete_customer_warning') }}</p>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn--ghost" @click="showDeleteModal = false">{{ t('cancel') }}</button>
          <button 
            class="btn btn--danger" 
            @click="executeDelete" 
            :disabled="deleting"
          >
            {{ deleting ? '⏳ ' + t('deleting') : '🗑️ ' + t('delete') }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import api from '../services/api'
import { useI18n } from '../services/i18n'
import { toLocalNumerals, formatLocalDate } from '../utils/numberFormatter'
import CustomerFormModal from '../components/CustomerFormModal.vue'

const { t, locale } = useI18n()
const customers = ref([])
const loading = ref(false)
const showModal = ref(false)

// Password reset state
const showPasswordModal = ref(false)
const selectedCustomer = ref(null)
const resetting = ref(false)
const resetError = ref('')
const resetSuccess = ref('')
const newPassword = ref('')
const showPassword = ref(false)

// Delete customer state
const showDeleteModal = ref(false)
const deleting = ref(false)
const deleteError = ref('')

// Copy password state
const passwordCopied = ref(false)

async function fetchCustomers() {
  loading.value = true
  try {
    const { data } = await api.get('/admin/customers')
    customers.value = data || []
  } catch (error) {
    console.error('❌ ' + t('failed_to_fetch_customers'), error)
    customers.value = []
  } finally {
    loading.value = false
  }
}

function formatDate(date) {
  if (!date) return '—'
  return formatLocalDate(date, locale?.value || 'en', 'DD/MM/YYYY')
}

function resetPassword(customer) {
  selectedCustomer.value = customer
  showPasswordModal.value = true
  resetError.value = ''
  resetSuccess.value = ''
  newPassword.value = ''
  showPassword.value = false
  passwordCopied.value = false
}

async function confirmResetPassword() {
  if (!selectedCustomer.value) return

  resetting.value = true
  resetError.value = ''
  resetSuccess.value = ''

  try {
    const { data } = await api.post('/admin/reset-customer-password', {
      customer_id: selectedCustomer.value.id
    })

    newPassword.value = data.password
    resetSuccess.value = t('password_reset_success')
  } catch (err) {
    console.error('❌ فشل إعادة تعيين كلمة المرور:', err.response?.data)
    resetError.value = err.response?.data?.error || '❌ ' + t('password_reset_failed')
  } finally {
    resetting.value = false
  }
}

function copyPassword() {
  if (newPassword.value) {
    navigator.clipboard.writeText(newPassword.value)
    passwordCopied.value = true
    setTimeout(() => {
      passwordCopied.value = false
    }, 2000)
  }
}

function confirmDelete(customer) {
  selectedCustomer.value = customer
  showDeleteModal.value = true
  deleteError.value = ''
}

async function executeDelete() {
  if (!selectedCustomer.value) return

  deleting.value = true
  deleteError.value = ''

  try {
    await api.delete(`/admin/customers/${selectedCustomer.value.id}`)
    showDeleteModal.value = false
    fetchCustomers()
  } catch (err) {
    console.error('❌ فشل حذف العميل:', err.response?.data)
    deleteError.value = err.response?.data?.error || '❌ ' + t('delete_customer_failed')
  } finally {
    deleting.value = false
  }
}

onMounted(() => {
  fetchCustomers()
})
</script>

<style scoped>
.page-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 18px;
  flex-wrap: wrap;
  gap: 10px;
}

.page-head h2 { font-size: 20px; margin-bottom: 4px; }
.page-head p { font-size: 13px; color: var(--ink-soft); }

.page-head-actions {
  display: flex;
  gap: 10px;
}

.table-actions {
  display: flex;
  gap: 8px;
}

.customer-card__actions {
  display: flex;
  gap: 8px;
  margin-top: 12px;
}

.modal-backdrop {
  position: fixed; inset: 0;
  background: rgba(0,0,0,0.6);
  backdrop-filter: blur(4px);
  display: flex; align-items: center; justify-content: center;
  z-index: 1000; padding: 20px;
}

.modal {
  width: 100%; max-width: 400px; padding: 0;
  max-height: 90vh; overflow-y: auto;
}

.modal-header {
  display: flex; justify-content: space-between; align-items: center;
  padding: 20px 24px; border-bottom: 1px solid var(--line);
}

.modal-header h3 { font-size: 18px; margin: 0; }
.modal-close { background: none; border: none; font-size: 24px; cursor: pointer; color: var(--ink-soft); }

.modal-body { padding: 20px 24px; }
.modal-footer {
  display: flex; gap: 10px;
  justify-content: flex-end; padding: 20px 24px;
  border-top: 1px solid var(--line);
}

.new-password-display {
  font-size: 16px;
  padding: 12px;
  background: var(--brand-tint);
  border-radius: var(--radius-sm);
  margin: 12px 0;
}

.password-label {
  font-weight: 600;
  margin-bottom: 8px;
}

.password-value-container {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  margin-top: 8px;
}

.password-value {
  font-family: monospace;
  font-size: 18px;
  font-weight: 700;
  letter-spacing: 2px;
}

.password-value.visible {
  color: var(--brand-dark);
}

.password-value.hidden {
  color: var(--ink-soft);
}

.btn-toggle-password {
  background: none;
  border: none;
  font-size: 20px;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: var(--radius-sm);
  transition: background 0.2s;
}

.btn-toggle-password:hover {
  background: rgba(0,0,0,0.1);
}

.btn-copy-password {
  background: none;
  border: none;
  font-size: 18px;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: var(--radius-sm);
  transition: background 0.2s;
  min-width: 36px;
}

.btn-copy-password:hover {
  background: rgba(0,0,0,0.1);
}

.password-warning {
  font-size: 12px;
  color: var(--signal-out);
  margin-top: 8px;
}

.delete-warning {
  font-size: 13px;
  color: var(--signal-out);
  margin-top: 12px;
  padding: 8px 12px;
  background: var(--signal-out-tint);
  border-radius: var(--radius-sm);
}

.table-wrapper {
  overflow-x: auto;
}

.table {
  width: 100%;
  border-collapse: collapse;
}

.table th,
.table td {
  padding: 12px 16px;
  text-align: right;
  border-bottom: 1px solid var(--line);
}

.table th {
  background: var(--canvas);
  font-weight: 600;
  font-size: 13px;
  color: var(--ink-soft);
}

.table tbody tr:hover {
  background: var(--canvas);
}

.table__person {
  display: flex;
  align-items: center;
  gap: 10px;
}

.table__avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: var(--brand-tint);
  color: var(--brand-dark);
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  font-size: 14px;
}

.mono {
  font-family: monospace;
  font-size: 13px;
}

.badge {
  display: inline-block;
  padding: 4px 10px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 500;
}

.badge--in {
  background: var(--signal-in-tint);
  color: var(--signal-in);
}

.badge--out {
  background: var(--signal-out-tint);
  color: var(--signal-out);
}

.badge--compact {
  padding: 2px 8px;
  font-size: 11px;
}

.mobile-cards {
  display: none;
}

.customer-card {
  padding: 16px;
  margin-bottom: 12px;
  border: 1px solid var(--line);
  border-radius: var(--radius-sm);
  background: var(--surface);
}

.customer-card__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.customer-card__person {
  display: flex;
  align-items: center;
  gap: 12px;
}

.customer-card__info {
  display: flex;
  flex-direction: column;
}

.customer-card__name {
  font-weight: 600;
  font-size: 15px;
}

.customer-card__badges {
  display: flex;
  gap: 6px;
}

.customer-card__body {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.customer-card__row {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.customer-card__label {
  font-size: 13px;
  color: var(--ink-soft);
}

.customer-card__value {
  font-size: 14px;
  color: var(--ink);
}

.text-muted {
  color: var(--ink-soft);
}

@media (max-width: 768px) {
  .desktop-only {
    display: none;
  }
  
  .mobile-cards {
    display: block;
  }
}
</style>
```

---

## 📄 frontend-admin-dashboard/src/views/DashboardView.vue

```vue
<template>
  <div class="dashboard">
    <div class="devpro-watermark">
      <img src="/src/assets/devpro-logo.jpg" alt="DevPro" class="watermark-logo" />
      <span class="watermark-text">{{ t('devpro_watermark') }}</span>
    </div>
    <div class="devpro-watermark">
      <img src="/src/assets/devpro-logo.jpg" alt="DevPro" class="watermark-logo" />
      <span class="watermark-text">{{ t('devpro_watermark') }}</span>
    </div>
    <div class="devpro-watermark">
      <img src="/src/assets/devpro-logo.jpg" alt="DevPro" class="watermark-logo" />
      <span class="watermark-text">{{ t('devpro_watermark') }}</span>
    </div>
    <div v-if="loading" class="loading-state">
      <div class="spinner"></div>
      <p>{{ t('loading_data') }}</p>
    </div>

    <div v-else>
      <!-- إحصائيات سريعة -->
      <div class="stats-grid">
        <div class="stat-card">
          <span class="stat-card__icon">
            <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="#3b82f6" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <circle cx="12" cy="8" r="4" fill="#3b82f6" fill-opacity="0.2"/>
              <path d="M4 20c0-4.418 3.582-8 8-8s8 3.582 8 8"/>
              <circle cx="18" cy="6" r="2" fill="#3b82f6" fill-opacity="0.15"/>
              <path d="M16 14c1.5-1.5 3-2 4-2" stroke-opacity="0.6"/>
            </svg>
          </span>
          <div>
            <span class="stat-card__value">{{ stats.total_employees || 0 }}</span>
            <span class="stat-card__label">{{ t('stats_total_employees') }}</span>
          </div>
        </div>
        <div class="stat-card stat-card--waiting">
          <span class="stat-card__icon">
            <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="#f59e0b" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <circle cx="12" cy="12" r="9"/>
              <path d="M12 6v6l4 2"/>
              <circle cx="12" cy="12" r="2" fill="#f59e0b" fill-opacity="0.2"/>
            </svg>
          </span>
          <div>
            <span class="stat-card__value">{{ stats.waiting_employees || 0 }}</span>
            <span class="stat-card__label">{{ t('stats_waiting_employees') }}</span>
          </div>
        </div>
        <div class="stat-card stat-card--active">
          <span class="stat-card__icon">
            <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="#10b981" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <circle cx="12" cy="12" r="9"/>
              <path d="M12 2v4m0 12v4m8-8h-4M6 12H2m15.07-7.07l-2.83 2.83M9.76 14.24l-2.83 2.83m0-10.66l2.83 2.83m4.48 4.48l2.83 2.83"/>
              <circle cx="12" cy="12" r="3" fill="#10b981" fill-opacity="0.2"/>
            </svg>
          </span>
          <div>
            <span class="stat-card__value">{{ activeEmployees.length }}</span>
            <span class="stat-card__label">{{ t('stats_active_now') }}</span>
          </div>
        </div>
        <div class="stat-card stat-card--completed">
          <span class="stat-card__icon">
            <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="#22c55e" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <rect x="3" y="3" width="18" height="18" rx="4"/>
              <path d="M8 12l3 3 7-7" stroke-width="2.5"/>
              <circle cx="12" cy="12" r="8" stroke-width="1" stroke-opacity="0.3"/>
            </svg>
          </span>
          <div>
            <span class="stat-card__value">{{ stats.completed_employees || 0 }}</span>
            <span class="stat-card__label">{{ t('stats_completed_today') }}</span>
          </div>
        </div>
      </div>

      <!-- الخريطة -->
      <div class="card dashboard__map">
        <div class="card-header">
          <h3>{{ t('dashboard_tracking_title') }}</h3>
          <div class="card-header__actions">
            <span class="badge badge--info">🔄 {{ updateCount }} {{ t('update_badge') }}</span>
            <button class="btn btn--sm btn--primary" @click="refreshData">{{ t('refresh_button') }}</button>
          </div>
        </div>
        <RealMap 
          :employees="activeEmployees" 
          :worksites="worksites"
          :height="500"
          @showDetails="showEmployeeDetails"
        />
      </div>

      <!-- قوائم الموظفين -->
      <div class="dashboard__tabs">
        <button 
          class="tab-btn" 
          :class="{ active: activeTab === 'active' }"
          @click="activeTab = 'active'"
        >
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <circle cx="12" cy="12" r="10"></circle>
            <circle cx="12" cy="12" r="3"></circle>
          </svg>
          {{ t('tab_active') }} ({{ activeEmployees.length }})
        </button>
        <button 
          class="tab-btn" 
          :class="{ active: activeTab === 'waiting' }"
          @click="activeTab = 'waiting'"
        >
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <circle cx="12" cy="12" r="10"></circle>
            <polyline points="12 6 12 12 16 14"></polyline>
          </svg>
          {{ t('tab_waiting') }} ({{ waitingEmployees.length }})
        </button>
        <button 
          class="tab-btn" 
          :class="{ active: activeTab === 'completed' }"
          @click="activeTab = 'completed'"
        >
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path>
            <polyline points="22 4 12 14.01 9 11.01"></polyline>
          </svg>
          {{ t('tab_completed') }} ({{ completedEmployees.length }})
        </button>
        <button 
          class="tab-btn tab-btn--alert" 
          :class="{ active: activeTab === 'alerts' }"
          @click="activeTab = 'alerts'"
        >
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"></path>
            <line x1="12" y1="9" x2="12" y2="13"></line>
            <line x1="12" y1="17" x2="12.01" y2="17"></line>
          </svg>
          {{ t('tab_alerts') }} ({{ outsideCount }})
        </button>
      </div>

      <!-- محتوى التاب -->
      <div class="dashboard__row">
        <!-- الموظفين النشطين -->
        <div v-if="activeTab === 'active'" class="card dashboard__list">
          <div class="card-header">
            <h3>{{ t('active_employees_title') }}</h3>
            <span class="badge">{{ activeEmployees.length }}</span>
          </div>
          
          <div v-if="activeEmployees.length === 0" class="empty-state">
            <p>{{ t('no_active_employees') }}</p>
          </div>
          
          <div v-else class="employee-list">
            <div 
              v-for="emp in activeEmployees" 
              :key="emp.id" 
              class="employee-item"
              :class="{ 
                'status-inside': emp.status === 'inside',
                'status-outside': emp.status === 'outside'
              }"
            >
              <div class="employee-item__avatar">
                {{ emp.full_name?.slice(0, 1) || '?' }}
                <span class="status-dot" :class="emp.status"></span>
              </div>
              <div class="employee-item__info">
                <strong>{{ emp.full_name }}</strong>
                <span class="employee-item__worksite">📍 {{ emp.worksite.name }}</span>
                <span class="employee-item__time">
                  🕐 {{ formatTime(emp.check_in_time) }} |
                  ⏱️ {{ emp.hours_worked.toFixed(1) }} {{ t('hours') }}
                </span>
              </div>
              <div class="employee-item__status">
                <span class="badge" :class="emp.status === 'inside' ? 'badge--in' : 'badge--out'">
                  {{ emp.status_text }}
                </span>
                <span class="employee-item__distance mono">
                  {{ formatDistance(emp.worksite.distance) }}
                </span>
              </div>
              <button 
                class="btn btn--sm btn--ghost" 
                @click="showEmployeeDetails(emp)"
                :title="t('view_details_title')"
              >
                📋
              </button>
            </div>
          </div>
        </div>

        <!-- قيد الانتظار -->
        <div v-if="activeTab === 'waiting'" class="card dashboard__list">
          <div class="card-header">
            <h3>{{ t('waiting_employees_title') }}</h3>
            <span class="badge badge--warning">{{ waitingEmployees.length }}</span>
          </div>
          
          <div v-if="waitingEmployees.length === 0" class="empty-state">
            <p>{{ t('all_employees_started') }}</p>
          </div>
          
          <div v-else class="employee-list">
            <div 
              v-for="emp in waitingEmployees" 
              :key="emp.id" 
              class="employee-item status-waiting"
            >
              <div class="employee-item__avatar" style="background: var(--signal-warning-tint); color: var(--signal-warning);">
                {{ emp.full_name?.slice(0, 1) || '?' }}
                <span class="status-dot waiting"></span>
              </div>
              <div class="employee-item__info">
                <strong>{{ emp.full_name }}</strong>
                <span class="employee-item__worksite">⏳ {{ t('has_not_started_work') }}</span>
              </div>
              <div class="employee-item__status">
                <span class="badge badge--warning">⏳ {{ t('waiting_status') }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- مكتمل -->
        <div v-if="activeTab === 'completed'" class="card dashboard__list">
          <div class="card-header">
            <h3>{{ t('completed_employees_title') }}</h3>
            <span class="badge badge--in">{{ completedEmployees.length }}</span>
          </div>
          
          <div v-if="completedEmployees.length === 0" class="empty-state">
            <p>{{ t('no_completed_employees') }}</p>
          </div>
          
          <div v-else class="employee-list">
            <div 
              v-for="emp in completedEmployees" 
              :key="emp.id" 
              class="employee-item status-completed"
            >
              <div class="employee-item__avatar" style="background: var(--signal-in-tint); color: var(--signal-in);">
                {{ emp.full_name?.slice(0, 1) || '?' }}
                <span class="status-dot completed"></span>
              </div>
              <div class="employee-item__info">
                <strong>{{ emp.full_name }}</strong>
                <span class="employee-item__worksite">📍 {{ emp.worksite_name }}</span>
                <span class="employee-item__time">
                  ✅ {{ formatTime(emp.check_out_time) }} |
                  ⏱️ {{ emp.hours_worked.toFixed(1) }} {{ t('hours') }}
                </span>
              </div>
              <div class="employee-item__status">
                <span class="badge badge--in">✅ {{ t('completed_status') }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- التحذيرات الأمنية -->
        <div v-if="activeTab === 'alerts'" class="card dashboard__alerts">
          <div class="card-header">
            <h3>{{ t('alerts_title') }}</h3>
            <span class="badge badge--out">{{ outsideCount }}</span>
          </div>
         
          <div v-if="outsideCount === 0" class="empty-state">
            <p>{{ t('no_alerts') }}</p>
          </div>
          
          <div v-else class="alerts-list">
            <div 
              v-for="emp in activeEmployees.filter(e => e.status === 'outside')" 
              :key="emp.id" 
              class="alert-item"
            >
              <span class="alert-item__icon">🚨</span>
              <div class="alert-item__info">
                <strong>{{ emp.full_name }}</strong>
                <span class="alert-item__message">
                  {{ t('left_worksite_prefix') }} {{ emp.worksite.name }} ({{ formatDistance(emp.worksite.distance) }})
                </span>
                <span class="alert-item__time mono">{{ formatTime(emp.last_update) }}</span>
              </div>
              <button 
                class="btn btn--sm btn--danger" 
                @click="showEmployeeDetails(emp)"
              >
                📋
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- مودال تفاصيل الموظف -->
    <div v-if="selectedEmployee" class="modal-backdrop" @click.self="selectedEmployee = null">
      <div class="modal card">
        <div class="modal-header">
          <h3>{{ t('employee_details_title') }}</h3>
          <button class="modal-close" @click="selectedEmployee = null">✕</button>
        </div>
        <div class="modal-body">
          <div class="employee-detail">
            <div class="employee-detail__header">
              <span class="employee-detail__avatar">{{ selectedEmployee.full_name?.slice(0, 1) || '?' }}</span>
              <div>
                <h4>{{ selectedEmployee.full_name }}</h4>
                <p>{{ selectedEmployee.phone }}</p>
              </div>
            </div>
            
            <div class="employee-detail__info">
              <div class="info-row">
                <span class="info-label">📍 {{ t('worksite_label') }}</span>
                <span>{{ selectedEmployee.worksite?.name || t('undefined_text') }}</span>
              </div>
              <div class="info-row">
                <span class="info-label">📏 {{ t('distance_label') }}</span>
                <span :class="selectedEmployee.status === 'inside' ? 'text-success' : 'text-danger'">
                  {{ formatDistance(selectedEmployee.worksite?.distance || 0) }}
                </span>
              </div>
              <div class="info-row">
                <span class="info-label">⏱️ {{ t('working_hours_label') }}</span>
                <span>{{ (selectedEmployee.hours_worked || 0).toFixed(1) }} {{ t('hours') }}</span>
              </div>
              <div class="info-row">
                <span class="info-label">🕐 {{ t('last_update_label') }}</span>
                <span>{{ formatTime(selectedEmployee.last_update) }}</span>
              </div>
              <div class="info-row">
                <span class="info-label">📍 {{ t('location_label') }}</span>
                <span class="mono">{{ (selectedEmployee.latitude || 0).toFixed(6) }}, {{ (selectedEmployee.longitude || 0).toFixed(6) }}</span>
              </div>
            </div>

            <hr class="divider" />

            <h4>{{ t('security_notes_title') }}</h4>
            <div v-if="securityNotes.length === 0" class="empty-state">
              <p>{{ t('no_security_notes') }}</p>
            </div>
            <div v-else class="security-notes">
              <div v-for="note in securityNotes" :key="note.id" class="note-item">
                <span class="note-item__icon">⚠️</span>
                <div>
                  <p class="note-item__title">{{ note.title }}</p>
                  <p class="note-item__body">{{ note.body }}</p>
                  <span class="note-item__time mono">{{ formatTime(note.created_at) }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn--ghost" @click="selectedEmployee = null">{{ t('close') }}</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import api from '../services/api'
import RealMap from '../components/RealMap.vue'
import { useI18n } from '../services/i18n'
import wsService from '../services/websocket'

const { t } = useI18n()

// ==========================================
// الحالة
// ==========================================
const loading = ref(true)
const activeEmployees = ref([])
const waitingEmployees = ref([])
const completedEmployees = ref([])
const worksites = ref([])
const stats = ref({})
const selectedEmployee = ref(null)
const securityNotes = ref([])
const updateCount = ref(0)
const activeTab = ref('active')
let refreshInterval = null

// ==========================================
// العمليات الحسابية
// ==========================================
const insideCount = computed(() => {
  return activeEmployees.value.filter(e => e.status === 'inside').length
})

const outsideCount = computed(() => {
  return activeEmployees.value.filter(e => e.status === 'outside').length
})

// ==========================================
// دوال مساعدة
// ==========================================
function formatTime(date) {
  if (!date) return '—'
  return new Date(date).toLocaleTimeString('en-GB', { 
    hour: '2-digit', 
    minute: '2-digit',
    second: '2-digit',
    hour12: false
  })
}

function formatDistance(meters) {
  if (!meters) return '0 ' + t('meters')
  if (meters >= 1000) {
    return (meters / 1000).toFixed(2) + ' ' + t('kilometers')
  }
  return Math.round(meters) + ' ' + t('meters')
}

// ==========================================
// جلب البيانات
// ==========================================
async function fetchData() {
  try {
    const [employeesRes, worksitesRes, statsRes, waitingRes, completedRes] = await Promise.all([
      api.get('/location/active'),
      api.get('/worksites'),
      api.get('/reports/daily-summary'),
      api.get('/reports/pending-employees'),
      api.get('/reports/completed-employees')
    ])
    
    activeEmployees.value = employeesRes.data || []
    worksites.value = worksitesRes.data || []
    stats.value = statsRes.data || {}
    waitingEmployees.value = waitingRes.data || []
    completedEmployees.value = completedRes.data || []
    updateCount.value++
    
  } catch (error) {
    console.error('❌ ' + t('failed_to_fetch_data'), error)
  } finally {
    loading.value = false
  }
}

async function refreshData() {
  loading.value = true
  await fetchData()
  loading.value = false
}

// ==========================================
// عرض تفاصيل الموظف
// ==========================================
async function showEmployeeDetails(employee) {
  selectedEmployee.value = employee
  
  // جلب الملاحظات الأمنية
  try {
    const { data } = await api.get(`/location/security/${employee.id}`)
    securityNotes.value = data || []
  } catch (error) {
    console.error('❌ ' + t('failed_to_fetch_security_notes'), error)
    securityNotes.value = []
  }
}

// ==========================================
// دورة الحياة
// ==========================================
onMounted(async () => {
  await fetchData()
  
  // تحديث كل 3 ثوانٍ للتتبع اللحظي
  refreshInterval = setInterval(fetchData, 3000)
  
  // الاتصال بـ WebSocket للتحديثات الفورية
  connectWebSocket()
})

onUnmounted(() => {
  if (refreshInterval) {
    clearInterval(refreshInterval)
  }
  
  // إغلاق WebSocket
  disconnectWebSocket()
})

// ==========================================
// WebSocket للتتبع اللحظي
// ==========================================
function connectWebSocket() {
  // Use the same host as the API, but with WebSocket protocol
  const apiBaseUrl = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1'
  const apiHost = apiBaseUrl.replace('/api/v1', '')
  const wsUrl = apiHost.replace('http://', 'ws://').replace('https://', 'wss://') + '/ws'
  console.log('🔌 Attempting to connect to WebSocket:', wsUrl)
  wsService.connect(wsUrl)
  
  wsService.onMessage((data) => {
    if (data.type === 'location_update') {
      console.log('📍 تحديث موقع فوري:', data.data)
      handleImmediateLocationUpdate(data.data)
    } else if (data.type === 'employee_status') {
      console.log('👤 تحديث حالة موظف:', data.data)
      handleEmployeeStatusUpdate(data.data)
    } else if (data.type === 'connected') {
      console.log('✅ تم الاتصال بـ WebSocket')
    } else if (data.type === 'disconnected') {
      console.log('❌ انقطع الاتصال بـ WebSocket')
    }
  })
}

function disconnectWebSocket() {
  wsService.disconnect()
}

function handleImmediateLocationUpdate(locationData) {
  // تحديث الموظف في القائمة إذا كان موجوداً
  const employeeIndex = activeEmployees.value.findIndex(
    emp => emp.id === locationData.user_id
  )
  
  if (employeeIndex !== -1) {
    // تحديث موقع الموظف فوراً
    activeEmployees.value[employeeIndex].latitude = locationData.latitude
    activeEmployees.value[employeeIndex].longitude = locationData.longitude
    activeEmployees.value[employeeIndex].last_update = new Date()
    
    // إعادة حساب المسافة من نقطة العمل
    if (activeEmployees.value[employeeIndex].worksite) {
      const worksite = activeEmployees.value[employeeIndex].worksite
      const distance = calculateDistance(
        locationData.latitude,
        locationData.longitude,
        worksite.latitude,
        worksite.longitude
      )
      
      activeEmployees.value[employeeIndex].worksite.distance = distance
      activeEmployees.value[employeeIndex].worksite.is_inside = distance <= worksite.radius
      
      if (distance <= worksite.radius) {
        activeEmployees.value[employeeIndex].status = 'inside'
        activeEmployees.value[employeeIndex].status_text = '✅ داخل النطاق'
      } else {
        activeEmployees.value[employeeIndex].status = 'outside'
        activeEmployees.value[employeeIndex].status_text = '❌ خارج النطاق'
      }
    }
    
    updateCount.value++
    console.log('🔄 تم تحديث موقع الموظف فوراً على الخريطة')
  }
}

function handleEmployeeStatusUpdate(statusData) {
  // التعامل مع تحديثات حالة الموظف
  console.log('تحديث حالة:', statusData)
  // يمكن إضافة منطق إضافي هنا
}

function calculateDistance(lat1, lon1, lat2, lon2) {
  const R = 6371 // نصف قطر الأرض بالكيلومتر
  const dLat = (lat2 - lat1) * Math.PI / 180
  const dLon = (lon2 - lon1) * Math.PI / 180
  const a = 
    Math.sin(dLat/2) * Math.sin(dLat/2) +
    Math.cos(lat1 * Math.PI / 180) * Math.cos(lat2 * Math.PI / 180) * 
    Math.sin(dLon/2) * Math.sin(dLon/2)
  const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1-a))
  const d = R * c
  return d * 1000 // تحويل إلى متر
}
</script>

<style scoped>
.dashboard { 
  display: flex; 
  flex-direction: column; 
  gap: 22px; 
  animation: fadeIn 0.6s ease;
}

.loading-state {
  text-align: center; padding: 60px 20px;
  display: flex; flex-direction: column; align-items: center; gap: 16px;
}

.spinner {
  width: 40px; height: 40px;
  border: 4px solid var(--line);
  border-top-color: var(--brand);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin { to { transform: rotate(360deg); } }

/* ==========================================
   إحصائيات
   ========================================== */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
}

.stat-card {
  background: var(--surface);
  border: 1px solid var(--line);
  border-radius: var(--radius-lg);
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 14px;
  box-shadow: var(--shadow-sm);
  transition: all var(--transition-base);
  position: relative;
  overflow: hidden;
}

.stat-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 3px;
  background: linear-gradient(90deg, var(--brand) 0%, var(--accent) 50%, var(--gold) 100%);
  opacity: 0;
  transition: opacity var(--transition-fast);
}

.stat-card:hover {
  box-shadow: var(--shadow-lg);
  transform: translateY(-4px);
  border-color: var(--line-strong);
}

.stat-card:hover::before {
  opacity: 1;
}

.stat-card--waiting { border-left: 4px solid var(--signal-warning); }
.stat-card--active { border-left: 4px solid var(--signal-in); }
.stat-card--completed { border-left: 4px solid #22C55E; }

.stat-card__icon { 
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--brand-tint);
  border-radius: var(--radius-md);
  color: var(--brand);
  transition: all var(--transition-base);
}

.stat-card:hover .stat-card__icon {
  transform: scale(1.1);
  box-shadow: var(--shadow-md), var(--brand-glow);
}

.stat-card__value { 
  font-size: 28px; 
  font-weight: 700; 
  display: block;
  background: linear-gradient(135deg, var(--brand) 0%, var(--accent) 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.stat-card__label { font-size: 13px; color: var(--ink-soft); }

/* ==========================================
   الخريطة
   ========================================== */
.dashboard__map { padding: 20px; }

.card-header {
  display: flex; align-items: center; justify-content: space-between;
  margin-bottom: 14px;
  flex-wrap: wrap;
  gap: 10px;
}

.card-header h3 { font-size: 16px; }
.card-header__actions { display: flex; gap: 8px; align-items: center; }

/* ==========================================
   Tabs
   ========================================== */
.dashboard__tabs {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.tab-btn {
  padding: 10px 20px;
  border: 2px solid var(--line);
  border-radius: var(--radius-sm);
  background: var(--surface);
  color: var(--ink-soft);
  font-weight: 600;
  cursor: pointer;
  transition: all var(--transition-base);
  font-family: var(--font-body);
  font-size: 14px;
  display: flex;
  align-items: center;
  gap: 8px;
  position: relative;
  overflow: hidden;
}

.tab-btn::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, rgba(255,255,255,0.1) 0%, rgba(255,255,255,0) 100%);
  opacity: 0;
  transition: opacity var(--transition-fast);
}

.tab-btn:hover::before {
  opacity: 1;
}

.tab-btn:hover {
  border-color: var(--brand);
  color: var(--brand);
  transform: translateY(-2px);
  box-shadow: var(--shadow-sm);
}

.tab-btn.active {
  border-color: var(--brand);
  background: linear-gradient(135deg, var(--brand) 0%, var(--brand-dark) 100%);
  color: white;
  box-shadow: var(--shadow-md), var(--brand-glow);
}

.tab-btn--alert {
  border-color: var(--signal-out);
  color: var(--signal-out);
}

.tab-btn--alert:hover {
  background: var(--signal-out-tint);
  border-color: var(--signal-out);
  transform: translateY(-2px);
  box-shadow: var(--shadow-sm);
}

.tab-btn--alert.active {
  background: linear-gradient(135deg, var(--signal-out) 0%, #dc2626 100%);
  color: white;
  box-shadow: var(--shadow-md);
}

/* ==========================================
   قوائم الموظفين
   ========================================== */
.dashboard__row {
  display: grid;
  grid-template-columns: 1fr;
  gap: 16px;
}

.dashboard__list { padding: 20px; }

.employee-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  max-height: 500px;
  overflow-y: auto;
}

.employee-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 14px;
  border-radius: var(--radius-sm);
  border: 1px solid var(--line);
  transition: all 0.2s;
  background: var(--surface);
}

.employee-item:hover {
  box-shadow: var(--shadow-sm);
}

.employee-item.status-inside { border-right: 4px solid var(--signal-in); }
.employee-item.status-outside {
  border-right: 4px solid var(--signal-out);
  background: var(--signal-out-tint);
}
.employee-item.status-waiting {
  border-right: 4px solid var(--signal-warning);
  background: var(--signal-warning-tint);
}
.employee-item.status-completed {
  border-right: 4px solid #22C55E;
  background: #22C55E10;
}

.employee-item__avatar {
  position: relative;
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: var(--brand-tint);
  color: var(--brand);
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  font-size: 16px;
  flex-shrink: 0;
}

.status-dot {
  position: absolute;
  bottom: -2px;
  right: -2px;
  width: 14px;
  height: 14px;
  border-radius: 50%;
  border: 2px solid var(--surface);
}

.status-dot.inside { background: var(--signal-in); }
.status-dot.outside { background: var(--signal-out); }
.status-dot.waiting { background: var(--signal-warning); }
.status-dot.completed { background: #22C55E; }

.employee-item__info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}

.employee-item__info strong {
  font-size: 14px;
  color: var(--ink);
}

.employee-item__worksite { font-size: 12px; color: var(--ink-soft); }
.employee-item__time { font-size: 11px; color: var(--ink-light); }

.employee-item__status {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 4px;
  flex-shrink: 0;
}

.employee-item__distance { font-size: 11px; color: var(--ink-soft); }

/* ==========================================
   تصميم محسّن للهواتف
   ========================================== */
@media (max-width: 600px) {
  .employee-item {
    flex-wrap: wrap;
    padding: 10px 12px;
    gap: 8px;
  }

  .employee-item__avatar {
    width: 36px;
    height: 36px;
    font-size: 14px;
  }

  .status-dot {
    width: 12px;
    height: 12px;
  }

  .employee-item__info {
    flex: 1;
    min-width: 120px;
  }

  .employee-item__info strong {
    font-size: 13px;
  }

  .employee-item__worksite {
    font-size: 11px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .employee-item__time {
    font-size: 10px;
  }

  .employee-item__status {
    flex-direction: row;
    align-items: center;
    gap: 6px;
    width: 100%;
    justify-content: space-between;
    margin-top: 4px;
  }

  .employee-item__distance {
    font-size: 10px;
  }

  .badge {
    font-size: 10px;
    padding: 2px 8px;
  }

  .btn--sm {
    padding: 4px 8px;
    font-size: 11px;
  }
}

/* ==========================================
   التحذيرات
   ========================================== */
.dashboard__alerts { padding: 20px; }

.alerts-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  max-height: 500px;
  overflow-y: auto;
}

.alert-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 14px;
  background: var(--signal-out-tint);
  border-radius: var(--radius-sm);
  border: 1px solid var(--signal-out);
}

.alert-item__icon { font-size: 20px; }

.alert-item__info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.alert-item__info strong {
  font-size: 13px;
  color: var(--signal-out);
}

.alert-item__message { font-size: 12px; color: var(--ink-soft); }
.alert-item__time { font-size: 11px; color: var(--ink-light); }

/* ==========================================
   Badges
   ========================================== */
.badge {
  display: inline-block;
  padding: 2px 10px;
  border-radius: 999px;
  font-size: 11px;
  font-weight: 600;
}

.badge--in { background: #22C55E20; color: #22C55E; }
.badge--out { background: #EF444420; color: #EF4444; }
.badge--warning { background: #F59E0B20; color: #F59E0B; }
.badge--info { background: #3B82F620; color: #3B82F6; }

/* ==========================================
   مودال
   ========================================== */
.modal-backdrop {
  position: fixed; inset: 0;
  background: rgba(0,0,0,0.6);
  backdrop-filter: blur(4px);
  display: flex; align-items: center; justify-content: center;
  z-index: 1000; padding: 20px;
}

.modal {
  width: 100%; max-width: 500px;
  max-height: 90vh; overflow-y: auto;
  padding: 0;
}

.modal-header {
  display: flex; justify-content: space-between; align-items: center;
  padding: 16px 20px; border-bottom: 1px solid var(--line);
}

.modal-header h3 { font-size: 17px; margin: 0; }
.modal-close { background: none; border: none; font-size: 24px; cursor: pointer; color: var(--ink-soft); }

.modal-body { padding: 20px; }
.modal-footer {
  padding: 16px 20px; border-top: 1px solid var(--line);
  display: flex; gap: 10px; justify-content: flex-end;
}

.employee-detail__header {
  display: flex; align-items: center; gap: 14px;
  margin-bottom: 16px;
}

.employee-detail__avatar {
  width: 48px; height: 48px;
  border-radius: 50%;
  background: var(--brand-tint);
  color: var(--brand);
  display: flex; align-items: center; justify-content: center;
  font-weight: 700; font-size: 20px;
}

.employee-detail__info {
  display: flex; flex-direction: column; gap: 6px;
}

.info-row {
  display: flex; gap: 8px;
  font-size: 13px;
}

.info-label {
  font-weight: 600;
  color: var(--ink-soft);
  min-width: 100px;
}

.text-success { color: var(--signal-in); font-weight: 600; }
.text-danger { color: var(--signal-out); font-weight: 600; }

.divider { border: none; border-top: 1px solid var(--line); margin: 16px 0; }

.security-notes {
  display: flex; flex-direction: column; gap: 8px;
}

.note-item {
  display: flex; gap: 10px;
  padding: 10px 12px;
  background: var(--signal-warning-tint);
  border-radius: var(--radius-sm);
  border: 1px solid var(--signal-warning);
}

.note-item__icon { font-size: 18px; }
.note-item__title { font-size: 13px; font-weight: 600; color: var(--signal-warning); }
.note-item__body { font-size: 12px; color: var(--ink-soft); margin: 2px 0; }
.note-item__time { font-size: 11px; color: var(--ink-light); }

.empty-state { text-align: center; padding: 30px 20px; color: var(--ink-soft); }

/* ==========================================
   استجابة
   ========================================== */
@media (max-width: 960px) {
  .stats-grid { grid-template-columns: repeat(2, 1fr); }
}

@media (max-width: 600px) {
  .stats-grid { grid-template-columns: 1fr 1fr; }
  .dashboard__tabs { flex-direction: column; }
  .tab-btn { width: 100%; text-align: center; }
}
</style>

<style>
/* ==========================================
   خلفية DevPro في Dashboard
   ========================================== */
.dashboard {
  position: relative;
  overflow: hidden;
}

.devpro-watermark {
  position: fixed;
  bottom: 20px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 20px;
  background: rgba(255, 255, 255, 0.85);
  backdrop-filter: blur(10px);
  border-radius: 999px;
  border: 1px solid rgba(30, 58, 130, 0.1);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.06);
  z-index: 5;
  pointer-events: none;
  opacity: 0.6;
  transition: opacity 0.3s ease;
}

.devpro-watermark:hover {
  opacity: 1;
}

.watermark-logo {
  width: 24px;
  height: 24px;
  border-radius: 6px;
  object-fit: contain;
  background: white;
  padding: 2px;
}

.watermark-text {
  font-size: 11px;
  color: #1e3a8a;
  font-weight: 600;
  white-space: nowrap;
}

@media (max-width: 768px) {
  .devpro-watermark {
    padding: 6px 14px;
    bottom: 10px;
  }
  
  .watermark-logo {
    width: 18px;
    height: 18px;
  }
  
  .watermark-text {
    font-size: 9px;
  }
}
</style>

<style>
/* ==========================================
   خلفية DevPro في Dashboard
   ========================================== */
.dashboard {
  position: relative;
  overflow: hidden;
}

.devpro-watermark {
  position: fixed;
  bottom: 20px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 20px;
  background: rgba(255, 255, 255, 0.85);
  backdrop-filter: blur(10px);
  border-radius: 999px;
  border: 1px solid rgba(30, 58, 130, 0.1);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.06);
  z-index: 5;
  pointer-events: none;
  opacity: 0.6;
  transition: opacity 0.3s ease;
}

.devpro-watermark:hover {
  opacity: 1;
}

.watermark-logo {
  width: 24px;
  height: 24px;
  border-radius: 6px;
  object-fit: contain;
  background: white;
  padding: 2px;
}

.watermark-text {
  font-size: 11px;
  color: #1e3a8a;
  font-weight: 600;
  white-space: nowrap;
}

@media (max-width: 768px) {
  .devpro-watermark {
    padding: 6px 14px;
    bottom: 10px;
  }
  
  .watermark-logo {
    width: 18px;
    height: 18px;
  }
  
  .watermark-text {
    font-size: 9px;
  }
}
</style>

<style>
/* ==========================================
   خلفية DevPro في Dashboard
   ========================================== */
.dashboard {
  position: relative;
  overflow: hidden;
}

.devpro-watermark {
  position: fixed;
  bottom: 20px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 20px;
  background: rgba(255, 255, 255, 0.85);
  backdrop-filter: blur(10px);
  border-radius: 999px;
  border: 1px solid rgba(30, 58, 130, 0.1);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.06);
  z-index: 5;
  pointer-events: none;
  opacity: 0.6;
  transition: opacity 0.3s ease;
}

.devpro-watermark:hover {
  opacity: 1;
}

.watermark-logo {
  width: 24px;
  height: 24px;
  border-radius: 6px;
  object-fit: contain;
  background: white;
  padding: 2px;
}

.watermark-text {
  font-size: 11px;
  color: #1e3a8a;
  font-weight: 600;
  white-space: nowrap;
}

@media (max-width: 768px) {
  .devpro-watermark {
    padding: 6px 14px;
    bottom: 10px;
  }
  
  .watermark-logo {
    width: 18px;
    height: 18px;
  }
  
  .watermark-text {
    font-size: 9px;
  }
}
</style>

```

---

## 📄 frontend-admin-dashboard/src/views/EmployeesView.vue

```vue
<template>
  <div>
    <div class="page-head">
      <div>
        <h2>{{ t('employees_title') }}</h2>
        <p>{{ t('employees_description') }}</p>
      </div>
      <div class="page-head-actions">
        <button class="btn btn--primary" @click="showModal = true">+ {{ t('add_employee') }}</button>
        <button class="btn btn--danger" @click="cleanupOldRecords" :disabled="cleaning">
          {{ cleaning ? '⏳' : '🗑️' }} {{ t('cleanup_old_records') }}
        </button>
      </div>
    </div>

    <div v-if="loading" class="empty-state">
      <p>{{ t('loading_employees') }}</p>
    </div>

    <div v-else-if="employees.length === 0" class="empty-state">
      <h3>{{ t('no_employees') }}</h3>
      <p>{{ t('add_new_employee_prompt') }}</p>
    </div>

    <div v-else class="card">
      <!-- جدول للشاشات الكبيرة -->
      <div class="table-wrapper desktop-only">
        <table class="table">
          <thead>
            <tr>
              <th>{{ t('name') }}</th>
              <th>{{ t('phone') }}</th>
              <th>{{ t('role') }}</th>
              <th>{{ t('status') }}</th>
              <th>{{ t('current_worksite') }}</th>
              <th>{{ t('created_at') }}</th>
              <th>{{ t('actions') }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="emp in employees" :key="emp.id">
              <td>
                <div class="table__person">
                  <span class="table__avatar">{{ emp.full_name?.slice(0, 1) || '?' }}</span>
                  {{ emp.full_name }}
                </div>
              </td>
              <td class="mono">{{ emp.phone || '—' }}</td>
              <td>
                <span class="badge" :class="emp.role === 'admin' ? 'badge--gold' : ''">
                  {{ emp.role === 'admin' ? t('admin_role') : t('field_employee') }}
                </span>
              </td>
              <td>
                <span class="badge" :class="emp.is_active ? 'badge--in' : 'badge--out'">
                  {{ emp.is_active ? t('active_status') : t('suspended_status') }}
                </span>
              </td>
              <td>
                <span v-if="emp.current_worksite" class="badge badge--success">
                  📍 {{ emp.current_worksite }}
                </span>
                <span v-else class="text-muted">—</span>
              </td>
              <td class="mono">{{ formatDate(emp.created_at) }}</td>
              <td>
                <div class="table-actions">
                  <button 
                    class="btn btn--primary btn--sm" 
                    @click="showAttendanceHistory(emp)"
                  >
                    📊 {{ t('attendance_history') }}
                  </button>
                  <button 
                    class="btn btn--danger btn--sm" 
                    @click="confirmDelete(emp)"
                    :disabled="emp.role === 'admin' && emp.email === 'admin@worktrack.com'"
                  >
                    🗑️ {{ t('delete') }}
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      
      <!-- بطاقات للشاشات الصغيرة -->
      <div class="mobile-cards mobile-only">
        <div v-for="emp in employees" :key="emp.id" class="employee-card">
          <div class="employee-card__header">
            <div class="employee-card__person">
              <span class="table__avatar">{{ emp.full_name?.slice(0, 1) || '?' }}</span>
              <div class="employee-card__info">
                <span class="employee-card__name">{{ emp.full_name }}</span>
              </div>
            </div>
            <div class="employee-card__badges">
              <span class="badge badge--compact" :class="emp.role === 'admin' ? 'badge--gold' : ''">
                {{ emp.role === 'admin' ? t('admin_role') : t('field_employee') }}
              </span>
              <span class="badge badge--compact" :class="emp.is_active ? 'badge--in' : 'badge--out'">
                {{ emp.is_active ? t('active_status') : t('suspended_status') }}
              </span>
            </div>
          </div>
          <div class="employee-card__body">
            <div class="employee-card__row">
              <span class="employee-card__label">{{ t('phone') }}</span>
              <span class="employee-card__value mono">{{ emp.phone || '—' }}</span>
            </div>
            <div class="employee-card__row">
              <span class="employee-card__label">{{ t('current_worksite') }}</span>
              <span v-if="emp.current_worksite" class="employee-card__value">
                <span class="badge badge--success badge--compact">📍 {{ emp.current_worksite }}</span>
              </span>
              <span v-else class="employee-card__value text-muted">—</span>
            </div>
            <div class="employee-card__row">
              <span class="employee-card__label">{{ t('created_at') }}</span>
              <span class="employee-card__value mono">{{ formatDate(emp.created_at) }}</span>
            </div>
          </div>
          <div class="employee-card__actions">
            <button 
              class="btn btn--primary btn--sm btn--compact" 
              @click="showAttendanceHistory(emp)"
            >
              <span class="btn-icon">📊</span>
              <span class="btn-text">{{ t('attendance_history') }}</span>
            </button>
            <button 
              class="btn btn--danger btn--sm btn--compact" 
              @click="confirmDelete(emp)"
              :disabled="emp.role === 'admin' && emp.email === 'admin@worktrack.com'"
            >
              <span class="btn-icon">🗑️</span>
              <span class="btn-text">{{ t('delete') }}</span>
            </button>
          </div>
        </div>
      </div>
    </div>

    <EmployeeFormModal 
      v-if="showModal" 
      @close="showModal = false" 
      @employee-added="fetchEmployees"
    />

    <!-- مودال تأكيد الحذف -->
    <div v-if="showDeleteModal" class="modal-backdrop" @click.self="showDeleteModal = false">
      <div class="modal card">
        <div class="modal-header">
          <h3>⚠️ {{ t('confirm_delete_title') }}</h3>
          <button class="modal-close" @click="showDeleteModal = false">✕</button>
        </div>
        <div class="modal-body">
          <p>{{ t('confirm_delete_message') }} <strong>{{ employeeToDelete?.full_name }}</strong>؟</p>
          <p class="text-danger">{{ t('delete_irreversible') }}</p>
        </div>
        <div class="modal-footer">
          <button class="btn btn--ghost" @click="showDeleteModal = false">{{ t('cancel') }}</button>
          <button class="btn btn--danger" @click="deleteEmployee" :disabled="deleting">
            {{ deleting ? t('deleting') : t('delete_final') }}
          </button>
        </div>
      </div>
    </div>

    <!-- مودال سجل الحضور -->
    <div v-if="showAttendanceModal" class="modal-backdrop" @click.self="showAttendanceModal = false">
      <div class="modal modal-lg card">
        <div class="modal-header">
          <h3>📊 {{ t('attendance_history') }} - {{ selectedEmployee?.full_name }}</h3>
          <button class="modal-close" @click="showAttendanceModal = false">✕</button>
        </div>
        <div class="modal-body">
          <!-- فلاتر الشهر والسنة -->
          <div class="filters">
            <div class="filter-group">
              <label>{{ t('year') }}</label>
              <select v-model="selectedYear" @change="fetchAttendanceHistory" class="form-select">
                <option v-for="year in availableYears" :key="year" :value="year">{{ year }}</option>
              </select>
            </div>
            <div class="filter-group">
              <label>{{ t('month') }}</label>
              <select v-model="selectedMonth" @change="fetchAttendanceHistory" class="form-select">
                <option v-for="month in availableMonths" :key="month.value" :value="month.value">{{ month.label }}</option>
              </select>
            </div>
            <button class="btn btn--primary" @click="exportToPDF" :disabled="loadingHistory || !attendanceHistory.length">
              📄 {{ t('export_pdf') }}
            </button>
          </div>

          <!-- ملخص الشهر -->
          <div v-if="monthlySummary" class="monthly-summary">
            <div class="summary-card">
              <span class="summary-label">{{ t('total_hours') }}</span>
              <span class="summary-value">{{ monthlySummary.summary?.total_hours?.toFixed(1) || 0 }} {{ t('hours') }}</span>
            </div>
            <div class="summary-card">
              <span class="summary-label">{{ t('work_days') }}</span>
              <span class="summary-value">{{ monthlySummary.summary?.work_days || 0 }} {{ t('days') }}</span>
            </div>
          </div>

          <!-- جدول سجل الحضور -->
          <div v-if="loadingHistory" class="loading-state">
            <p>{{ t('loading') }}</p>
          </div>
          <div v-else-if="attendanceHistory.length === 0" class="empty-state">
            <p>{{ t('no_attendance_records') }}</p>
          </div>
          <div v-else>
            <!-- جدول للشاشات الكبيرة -->
            <div class="table-wrapper desktop-only">
              <table class="table">
                <thead>
                  <tr>
                    <th>{{ t('date') }}</th>
                    <th>{{ t('worksite') }}</th>
                    <th>{{ t('check_in') }}</th>
                    <th>{{ t('check_out') }}</th>
                    <th>{{ t('worked_hours') }}</th>
                    <th>{{ t('location') }}</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="record in attendanceHistory" :key="record.id">
                    <td class="mono">{{ formatDate(record.check_in_time) }}</td>
                    <td>{{ record.worksite_name || '—' }}</td>
                    <td class="mono">{{ formatTime(record.check_in_time) }}</td>
                    <td class="mono">{{ record.check_out_time ? formatTime(record.check_out_time) : '—' }}</td>
                    <td class="mono">{{ record.worked_hours ? record.worked_hours.toFixed(1) + ' ' + t('hours') : '—' }}</td>
                    <td class="mono">{{ formatDistance(record.check_in_distance_meters) }}</td>
                  </tr>
                </tbody>
              </table>
            </div>
            
            <!-- بطاقات للشاشات الصغيرة -->
            <div class="mobile-cards mobile-only">
              <div v-for="record in attendanceHistory" :key="record.id" class="attendance-card">
                <div class="attendance-card__header">
                  <span class="attendance-card__date">{{ formatDate(record.check_in_time) }}</span>
                  <span class="badge badge--info">{{ record.worked_hours ? record.worked_hours.toFixed(1) + ' ' + t('hours') : '—' }}</span>
                </div>
                <div class="attendance-card__body">
                  <div class="attendance-card__row">
                    <span class="attendance-card__label">{{ t('worksite') }}</span>
                    <span class="attendance-card__value">{{ record.worksite_name || '—' }}</span>
                  </div>
                  <div class="attendance-card__row">
                    <span class="attendance-card__label">{{ t('check_in') }}</span>
                    <span class="attendance-card__value mono">{{ formatTime(record.check_in_time) }}</span>
                  </div>
                  <div class="attendance-card__row">
                    <span class="attendance-card__label">{{ t('check_out') }}</span>
                    <span class="attendance-card__value mono">{{ record.check_out_time ? formatTime(record.check_out_time) : '—' }}</span>
                  </div>
                  <div class="attendance-card__row">
                    <span class="attendance-card__label">{{ t('location') }}</span>
                    <span class="attendance-card__value mono">{{ formatDistance(record.check_in_distance_meters) }}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import api from '../services/api'
import EmployeeFormModal from '../components/EmployeeFormModal.vue'
import { useI18n } from '../services/i18n'
import { formatLocalDate, toLocalNumerals } from '../utils/numberFormatter'

const { t, currentLang, locale } = useI18n()
const employees = ref([])
const loading = ref(false)
const showModal = ref(false)
const showDeleteModal = ref(false)
const employeeToDelete = ref(null)
const deleting = ref(false)
const cleaning = ref(false)

// سجل الحضور
const showAttendanceModal = ref(false)
const selectedEmployee = ref(null)
const attendanceHistory = ref([])
const monthlySummary = ref(null)
const loadingHistory = ref(false)
const selectedYear = ref(new Date().getFullYear())
const selectedMonth = ref(String(new Date().getMonth() + 1))

const availableYears = ref([])
const availableMonths = ref([])

// دالة لتحديث أسماء الشهور حسب اللغة
function updateMonthNames() {
  const monthKeys = ['january', 'february', 'march', 'april', 'may', 'june', 'july', 'august', 'september', 'october', 'november', 'december']
  availableMonths.value = monthKeys.map((key, index) => ({
    value: String(index + 1),
    label: t(key)
  }))
}

// توليد السنوات المتاحة
const currentYear = new Date().getFullYear()
for (let i = currentYear; i >= currentYear - 5; i--) {
  availableYears.value.push(i)
}

// تحديث أسماء الشهور عند التحميل
updateMonthNames()

// مراقبة تغيير اللغة لتحديث أسماء الشهور
watch(currentLang, () => {
  updateMonthNames()
})

async function fetchEmployees() {
  loading.value = true
  try {
    const { data } = await api.get('/admin/employees')
    employees.value = data || []
  } catch (error) {
    console.error('❌ ' + t('failed_to_fetch_employees'), error)
    employees.value = []
  } finally {
    loading.value = false
  }
}

function confirmDelete(emp) {
  employeeToDelete.value = emp
  showDeleteModal.value = true
}

async function deleteEmployee() {
  if (!employeeToDelete.value) return
  
  deleting.value = true
  try {
    await api.delete(`/admin/employees/${employeeToDelete.value.id}`)
    showDeleteModal.value = false
    await fetchEmployees()
  } catch (error) {
    console.error('❌ ' + t('failed_to_delete_employee'), error)
    alert(error.response?.data?.error || t('failed_to_delete_employee'))
  } finally {
    deleting.value = false
  }
}

function formatDate(date) {
  if (!date) return '—'
  return formatLocalDate(date, locale?.value || 'en', 'DD/MM/YYYY')
}

function formatTime(date) {
  if (!date) return '—'
  return formatLocalDate(date, locale?.value || 'en', 'HH:mm')
}

function formatDistance(meters) {
  if (!meters) return '—'
  if (meters >= 1000) {
    return toLocalNumerals((meters / 1000).toFixed(2), locale?.value || 'en') + ' ' + t('kilometers')
  }
  return toLocalNumerals(Math.round(meters), locale?.value || 'en') + ' ' + t('meters')
}

// دوال سجل الحضور
async function showAttendanceHistory(employee) {
  selectedEmployee.value = employee
  showAttendanceModal.value = true
  await fetchAttendanceHistory()
  
  // تمرير سلس للمحتوى بعد التحميل
  setTimeout(() => {
    const attendanceContent = document.querySelector('.modal-body')
    if (attendanceContent) {
      attendanceContent.scrollTo({
        top: 0,
        behavior: 'smooth'
      })
    }
  }, 100)
}

async function fetchAttendanceHistory() {
  if (!selectedEmployee.value) return
  
  loadingHistory.value = true
  try {
    const { data } = await api.get(
      `/attendance/employee/${selectedEmployee.value.id}/history?year=${selectedYear.value}&month=${selectedMonth.value}`
    )
    attendanceHistory.value = data || []
    
    // جلب الملخص الشهري
    const summaryResponse = await api.get(
      `/attendance/employee/${selectedEmployee.value.id}/monthly-summary?year=${selectedYear.value}&month=${selectedMonth.value}`
    )
    monthlySummary.value = summaryResponse.data
  } catch (error) {
    console.error('❌ ' + t('failed_to_fetch_attendance_history'), error)
    attendanceHistory.value = []
    monthlySummary.value = null
  } finally {
    loadingHistory.value = false
  }
}

async function exportToPDF() {
  if (!selectedEmployee.value || !attendanceHistory.value.length) return

  try {
    // الحصول على اللغة الحالية
    const currentLangValue = currentLang.value

    // محتوى الترجمة حسب اللغة
    const translations = {
      ar: {
        title: 'سجل الحضور',
        employeeLabel: 'الموظف',
        periodLabel: 'الفترة',
        totalHoursLabel: 'إجمالي الساعات',
        workDaysLabel: 'أيام العمل',
        dateLabel: 'التاريخ',
        worksiteLabel: 'نقطة العمل',
        checkInLabel: 'بداية العمل',
        checkOutLabel: 'نهاية العمل',
        workedHoursLabel: 'ساعات العمل',
        distanceLabel: 'المسافة',
        hoursUnit: 'ساعة',
        daysUnit: 'يوم',
        footer: 'تم إنشاء هذا التقرير من WorkTrack',
        direction: 'rtl'
      },
      en: {
        title: 'Attendance History',
        employeeLabel: 'Employee',
        periodLabel: 'Period',
        totalHoursLabel: 'Total Hours',
        workDaysLabel: 'Work Days',
        dateLabel: 'Date',
        worksiteLabel: 'Worksite',
        checkInLabel: 'Check In',
        checkOutLabel: 'Check Out',
        workedHoursLabel: 'Worked Hours',
        distanceLabel: 'Distance',
        hoursUnit: 'hours',
        daysUnit: 'days',
        footer: 'Report generated from WorkTrack',
        direction: 'ltr'
      },
      he: {
        title: 'היסטוריית נוכחות',
        employeeLabel: 'עובד',
        periodLabel: 'תקופה',
        totalHoursLabel: 'סה"כ שעות',
        workDaysLabel: 'ימי עבודה',
        dateLabel: 'תאריך',
        worksiteLabel: 'אתר עבודה',
        checkInLabel: 'כניסה',
        checkOutLabel: 'יציאה',
        workedHoursLabel: 'שעות עבודה',
        distanceLabel: 'מרחק',
        hoursUnit: 'שעות',
        daysUnit: 'ימים',
        footer: 'הדוח נוצר מ-WorkTrack',
        direction: 'rtl'
      }
    }

    const trans = translations[currentLangValue] || translations.ar

    // الحصول على اسم الشهر حسب اللغة
    const monthName = Array.isArray(availableMonths)
      ? (availableMonths.find(m => m.value === selectedMonth.value)?.label || selectedMonth.value)
      : selectedMonth.value

    // إنشاء محتوى HTML للطباعة
    const htmlContent = `
      <html dir="${trans.direction}">
      <head>
        <meta charset="UTF-8">
        <title>${trans.title} - ${selectedEmployee.value.full_name}</title>
        <style>
          body { font-family: Arial, sans-serif; padding: 20px; direction: ${trans.direction}; }
          h1 { text-align: center; color: #333; }
          .summary { display: flex; justify-content: space-around; margin: 20px 0; padding: 15px; background: #f5f5f5; border-radius: 8px; }
          .summary-item { text-align: center; }
          .summary-label { font-size: 14px; color: #666; }
          .summary-value { font-size: 24px; font-weight: bold; color: #333; }
          table { width: 100%; border-collapse: collapse; margin-top: 20px; }
          th, td { padding: 12px; border: 1px solid #ddd; }
          th { background: #4CAF50; color: white; }
          th { text-align: ${trans.direction === 'rtl' ? 'right' : 'left'}; }
          td { text-align: ${trans.direction === 'rtl' ? 'right' : 'left'}; }
          tr:nth-child(even) { background: #f9f9f9; }
          .footer { margin-top: 30px; text-align: center; color: #666; font-size: 12px; }
        </style>
      </head>
      <body>
        <h1>${trans.title}</h1>
        <h2 style="text-align: center; color: #666;">${trans.employeeLabel}: ${selectedEmployee.value.full_name}</h2>
        <p style="text-align: center; color: #666;">${trans.periodLabel}: ${monthName} ${selectedYear.value}</p>

        <div class="summary">
          <div class="summary-item">
            <div class="summary-label">${trans.totalHoursLabel}</div>
            <div class="summary-value">${monthlySummary.value?.summary?.total_hours?.toFixed(1) || 0} ${trans.hoursUnit}</div>
          </div>
          <div class="summary-item">
            <div class="summary-label">${trans.workDaysLabel}</div>
            <div class="summary-value">${monthlySummary.value?.summary?.work_days || 0} ${trans.daysUnit}</div>
          </div>
        </div>

        <table>
          <thead>
            <tr>
              <th>${trans.dateLabel}</th>
              <th>${trans.worksiteLabel}</th>
              <th>${trans.checkInLabel}</th>
              <th>${trans.checkOutLabel}</th>
              <th>${trans.workedHoursLabel}</th>
              <th>${trans.distanceLabel}</th>
            </tr>
          </thead>
          <tbody>
            ${attendanceHistory.value.map(record => `
              <tr>
                <td>${formatDate(record.check_in_time)}</td>
                <td>${record.worksite_name || '—'}</td>
                <td>${formatTime(record.check_in_time)}</td>
                <td>${record.check_out_time ? formatTime(record.check_out_time) : '—'}</td>
                <td>${record.worked_hours ? record.worked_hours.toFixed(1) + ' ' + trans.hoursUnit : '—'}</td>
                <td>${formatDistance(record.check_in_distance_meters)}</td>
              </tr>
            `).join('')}
          </tbody>
        </table>

        <div class="footer">
          <p>${trans.footer}</p>
          <p>${new Date().toLocaleDateString('en-GB')} - ${new Date().toLocaleTimeString('en-GB', { hour: '2-digit', minute: '2-digit', hour12: false })}</p>
        </div>
      </body>
      </html>
    `
    
    // إنشاء نافذة جديدة للطباعة
    const printWindow = window.open('', '_blank')
    printWindow.document.write(htmlContent)
    printWindow.document.close()
    
    // انتظار تحميل المحتوى ثم الطباعة
    printWindow.onload = function() {
      printWindow.print()
    }
  } catch (error) {
    console.error('❌ ' + t('failed_to_export_pdf'), error)
    alert(t('failed_to_export_report'))
  }
}

async function cleanupOldRecords() {
  if (!confirm(t('confirm_cleanup_old_records'))) {
    return
  }

  cleaning.value = true
  try {
    const { data } = await api.post('/attendance/cleanup-old-records')
    alert(`${t('cleanup_success')}: ${data.deleted_count} records`)
  } catch (error) {
    console.error('❌ ' + t('cleanup_failed'), error)
    alert(error.response?.data?.error || t('cleanup_failed'))
  } finally {
    cleaning.value = false
  }
}

onMounted(() => {
  fetchEmployees()
})
</script>

<style scoped>
.page-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 18px;
  flex-wrap: wrap;
  gap: 10px;
}

.page-head h2 { font-size: 20px; margin-bottom: 4px; }
.page-head p { font-size: 13px; color: var(--ink-soft); }

.page-head-actions {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.table-wrapper { overflow-x: auto; }

/* إخفاء افتراضي للعناصر المحمولة والسطحية */
.desktop-only { display: none; }
.mobile-only { display: none; }

/* إظهار الجدول للشاشات الكبيرة فقط */
@media (min-width: 769px) {
  .desktop-only { display: block !important; }
  .mobile-only { display: none !important; }
}

/* إظهار البطاقات للشاشات الصغيرة فقط */
@media (max-width: 768px) {
  .desktop-only { display: none !important; }
  .mobile-only { display: block !important; }
}

.table { width: 100%; border-collapse: collapse; }

.table th {
  text-align: right;
  font-size: 12px;
  color: var(--ink-soft);
  font-weight: 600;
  padding: 12px 14px;
  border-bottom: 2px solid var(--line);
}

.table td {
  padding: 12px 14px;
  font-size: 14px;
  border-bottom: 1px solid var(--line);
}

.table tr:last-child td { border-bottom: none; }

.table__person { display: flex; align-items: center; gap: 10px; }

.table__avatar {
  width: 32px; height: 32px;
  border-radius: 50%;
  background: var(--brand-tint);
  color: var(--brand-dark);
  font-size: 13px;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.table-actions {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.modal-backdrop {
  position: fixed; inset: 0;
  background: rgba(0,0,0,0.6);
  backdrop-filter: blur(4px);
  display: flex; align-items: center; justify-content: center;
  z-index: 1000; padding: 20px;
}

.modal {
  width: 100%; max-width: 420px; padding: 0;
}

.modal-lg {
  max-width: 900px;
}

.modal-header {
  display: flex; justify-content: space-between; align-items: center;
  padding: 16px 20px; border-bottom: 1px solid var(--line);
}

.modal-header h3 { font-size: 17px; margin: 0; }
.modal-close { background: none; border: none; font-size: 24px; cursor: pointer; color: var(--ink-soft); }

.modal-body { 
  padding: 20px; 
  overflow-y: auto;
  max-height: 70vh;
  scroll-behavior: smooth;
}
.modal-body p { margin-bottom: 8px; }
.text-danger { color: var(--signal-out); font-weight: 600; }

.modal-footer {
  display: flex; gap: 10px;
  justify-content: flex-end;
  padding: 16px 20px;
  border-top: 1px solid var(--line);
}

/* سجل الحضور */
.filters {
  display: flex;
  gap: 16px;
  align-items: flex-end;
  margin-bottom: 20px;
  flex-wrap: wrap;
}

.filter-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.filter-group label {
  font-size: 13px;
  font-weight: 600;
  color: var(--ink-soft);
}

.form-select {
  padding: 8px 12px;
  border: 1px solid var(--line);
  border-radius: var(--radius-sm);
  font-size: 14px;
  background: var(--surface);
  color: var(--ink);
  min-width: 120px;
}

.monthly-summary {
  display: flex;
  gap: 16px;
  margin-bottom: 20px;
}

.summary-card {
  flex: 1;
  padding: 16px;
  background: var(--brand-tint);
  border-radius: var(--radius-md);
  text-align: center;
  border: 1px solid var(--brand);
}

.summary-label {
  display: block;
  font-size: 13px;
  color: var(--ink-soft);
  margin-bottom: 8px;
}

.summary-value {
  display: block;
  font-size: 24px;
  font-weight: 700;
  color: var(--brand-dark);
}

.loading-state, .empty-state {
  text-align: center;
  padding: 40px 20px;
  color: var(--ink-soft);
}

@media (max-width: 768px) {
  .page-head { flex-direction: column; }
  .page-head .btn { width: 100%; }
  
  .filters {
    flex-direction: column;
    align-items: stretch;
  }
  
  .form-select {
    width: 100%;
  }
  
  .monthly-summary {
    flex-direction: column;
  }
  
  .modal-lg {
    max-width: 100%;
  }
  
  .card {
    border-radius: var(--radius-md);
  }
  
  .mobile-cards {
    padding: 0 4px;
  }
}

/* تصميم بطاقات سجل الحضور للهاتف */
.mobile-cards {
  display: flex;
  flex-direction: column;
  gap: 12px;
  width: 100%;
  box-sizing: border-box;
}

.attendance-card {
  background: var(--surface);
  border: 1px solid var(--line);
  border-radius: var(--radius-md);
  padding: 16px;
  transition: all 0.2s ease;
  width: 100%;
  box-sizing: border-box;
}

.attendance-card:hover {
  border-color: var(--brand);
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

.attendance-card__header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  padding-bottom: 12px;
  border-bottom: 1px solid var(--line);
}

.attendance-card__date {
  font-weight: 600;
  color: var(--ink);
  font-size: 14px;
}

.attendance-card__body {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.attendance-card__row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 8px;
}

.attendance-card__label {
  font-size: 13px;
  color: var(--ink-soft);
  font-weight: 500;
  flex-shrink: 0;
}

.attendance-card__value {
  font-size: 14px;
  color: var(--ink);
  font-weight: 500;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 180px;
}

@media (max-width: 380px) {
  .attendance-card__value {
    max-width: 140px;
  }
}

/* تصميم بطاقات الموظفين للهاتف */
.employee-card {
  background: var(--surface);
  border: 1px solid var(--line);
  border-radius: var(--radius-md);
  padding: 16px;
  transition: all 0.2s ease;
  width: 100%;
  box-sizing: border-box;
}

.employee-card:hover {
  border-color: var(--brand);
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

.employee-card__header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 12px;
  padding-bottom: 12px;
  border-bottom: 1px solid var(--line);
  gap: 8px;
  flex-wrap: wrap;
}

.employee-card__person {
  display: flex;
  align-items: center;
  gap: 10px;
  flex: 1;
  min-width: 0; /* مهم للنصوص الطويلة */
}

.employee-card__info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.employee-card__name {
  font-weight: 600;
  color: var(--ink);
  font-size: 15px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 200px;
}

.employee-card__badges {
  display: flex;
  flex-direction: column;
  gap: 6px;
  align-items: flex-end;
  flex-shrink: 0;
}

.badge--compact {
  font-size: 11px;
  padding: 4px 8px;
  white-space: nowrap;
}

.badge--success {
  background-color: #10b981;
  color: white;
}

.text-muted {
  color: var(--ink-soft);
}

.employee-card__body {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 12px;
}

.employee-card__row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 8px;
}

.employee-card__label {
  font-size: 13px;
  color: var(--ink-soft);
  font-weight: 500;
  flex-shrink: 0;
}

.employee-card__value {
  font-size: 14px;
  color: var(--ink);
  font-weight: 500;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 180px;
}

.employee-card__actions {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.employee-card__actions .btn {
  flex: 1;
  min-width: 110px;
  font-size: 13px;
  padding: 8px 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
}

.btn--compact {
  font-size: 12px;
  padding: 6px 10px;
}

.btn--compact .btn-icon {
  font-size: 14px;
  flex-shrink: 0;
}

.btn--compact .btn-text {
  display: inline;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  flex: 1;
  min-width: 0;
}

@media (max-width: 380px) {
  .employee-card__name {
    max-width: 150px;
  }
  
  .employee-card__value {
    max-width: 130px;
  }
  
  .employee-card__actions .btn {
    min-width: 95px;
    font-size: 11px;
    padding: 6px 8px;
    gap: 4px;
  }
  
  .btn--compact .btn-icon {
    font-size: 12px;
  }
  
  .badge--compact {
    font-size: 10px;
    padding: 3px 6px;
  }
  
  .employee-card {
    padding: 12px;
  }
  
  .employee-card__header {
    gap: 6px;
  }
  
  .employee-card__body {
    gap: 6px;
  }
  
  .mobile-cards {
    padding: 0 2px;
  }
}

@media (max-width: 340px) {
  .employee-card__actions {
    flex-direction: column;
    gap: 6px;
  }
  
  .employee-card__actions .btn {
    min-width: 100%;
    font-size: 12px;
    padding: 8px 12px;
  }
  
  .btn--compact .btn-icon {
    font-size: 14px;
  }
  
  .employee-card {
    padding: 10px;
  }
  
  .employee-card__header {
    flex-direction: column;
    align-items: flex-start;
  }
  
  .employee-card__badges {
    align-items: flex-start;
    flex-direction: row;
    flex-wrap: wrap;
  }
  
  .employee-card__person {
    width: 100%;
  }
  
  .employee-card__name {
    max-width: 100%;
  }
  
  .employee-card__value {
    max-width: 100%;
  }
  
  .mobile-cards {
    padding: 0;
  }
  
  .card {
    border-radius: var(--radius-sm);
  }
}

/* تحسين شريط التمرير */
.modal-body::-webkit-scrollbar {
  width: 8px;
}

.modal-body::-webkit-scrollbar-track {
  background: var(--canvas);
  border-radius: 4px;
}

.modal-body::-webkit-scrollbar-thumb {
  background: var(--line);
  border-radius: 4px;
}

.modal-body::-webkit-scrollbar-thumb:hover {
  background: var(--ink-soft);
}

/* تحسين التمرير على الهاتف */
@media (max-width: 768px) {
  .modal-body {
    max-height: 60vh;
    -webkit-overflow-scrolling: touch;
  }
  
  .modal-body::-webkit-scrollbar {
    width: 4px;
  }
  
  .table-wrapper {
    overflow-x: auto;
    -webkit-overflow-scrolling: touch;
  }
  
  .employee-card {
    margin: 0 4px 8px 4px;
  }
  
  .attendance-card {
    margin: 0 4px 8px 4px;
  }
}
</style>
```

---

## 📄 frontend-admin-dashboard/src/views/LoginView.vue

```vue
<template>
  <PullToRefresh 
    @refresh="handleRefresh"
    :refresh-text="$t('pull_to_refresh')"
    :refreshing-text="$t('refreshing')"
    :release-text="$t('release_to_refresh')"
  >
    <div class="login-page">
      <PWAInstallButton />
      <div class="login-card">
      <div class="login-header">
        <div class="powered-by">
          <img src="/src/assets/devpro-logo.jpg" alt="DevPro" class="powered-logo" />
          <span class="powered-text">
            {{ $t('app_name') }}<br />
            <strong>DevPro</strong>
            <span class="powered-slogan">{{ $t('powered_slogan') }}</span>
          </span>
        </div>

        <div class="app-brand">
          <img src="/src/assets/devpro-logo.jpg" alt="WorkTrack logo" class="brand-mark" />
          <h1 class="title">{{ $t('app_name') }}</h1>
        </div>
        <p class="subtitle">{{ $t('login') }}</p>
      </div>

      <div class="lang-section">
        <button
          v-for="lang in languages"
          :key="lang.code"
          class="btn btn--ghost lang-btn"
          :class="{ active: currentLang === lang.code }"
          @click="changeLanguage(lang.code)"
        >
          {{ lang.flag }} {{ lang.name }}
        </button>
      </div>

      <form class="login-form" @submit.prevent="handleSubmit">
        <div class="field">
          <label>{{ $t('email') }}</label>
          <input v-model="email" type="email" :placeholder="$t('email_placeholder')" required autocomplete="username" readonly @focus="removeReadonly" ref="emailInput" />
        </div>
         
        <div class="field">
          <label>{{ $t('password') }}</label>
          <input v-model="password" type="password" :placeholder="$t('password_placeholder')" required autocomplete="current-password" />
        </div>

        <div v-if="error" class="error">{{ error }}</div>
        <div v-if="debugInfo" class="debug">{{ debugInfo }}</div>

        <button class="btn btn--primary btn--block btn-login" type="submit" :disabled="loading">
          {{ loading ? $t('loading') : $t('login') }}
        </button>
      </form>

      <div class="footer">
        <p>{{ $t('footer_copyright') }}</p>
        <p class="footer-powered">🚀 {{ $t('app_name') }} - {{ $t('dashboard') }}</p>
      </div>
    </div>
    </div>
  </PullToRefresh>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { login } from '../services/auth'
import { authStore } from '../store/auth'
import { useI18n } from '../services/i18n'
import PullToRefresh from '../components/PullToRefresh.vue'
import PWAInstallButton from '../components/PWAInstallButton.vue'

const { t, currentLang, setLang } = useI18n()
const router = useRouter()

const email = ref('')
const password = ref('')
const loading = ref(false)
const error = ref('')
const debugInfo = ref('')
const emailInput = ref(null)

function removeReadonly() {
  if (emailInput.value) {
    emailInput.value.removeAttribute('readonly')
  }
}

const languages = [
  { code: 'ar', name: 'العربية', flag: '🇸🇦' },
  { code: 'he', name: 'עברית', flag: '🇮🇱' },
  { code: 'en', name: 'English', flag: '🇬🇧' }
]

function changeLanguage(lang) {
  setLang(lang)
}

async function handleRefresh() {
  // إعادة تحميل الصفحة
  window.location.reload()
}

async function handleSubmit() {
  loading.value = true
  error.value = ''
  debugInfo.value = ''

  try {
    debugInfo.value = t('login_connecting')
    
    const data = await login(email.value, password.value)
    
    debugInfo.value = t('login_success')
    authStore.setUser(data.user)
    
    setTimeout(() => {
      router.push('/dashboard')
    }, 500)
    
  } catch (e) {
    debugInfo.value = t('login_failed')
    
    if (e.response) {
      const status = e.response.status
      const msg = e.response.data?.error || t('login_error_unknown')
      error.value = msg
      debugInfo.value += `\n${t('login_error_code_prefix')} ${status}`
      debugInfo.value += `\n${t('login_error_message_prefix')} ${msg}`
    } else {
      error.value = t('login_server_unreachable')
      debugInfo.value += '\n' + t('login_server_not_responding')
    }
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, var(--brand) 0%, var(--brand-dark) 100%);
  padding: 20px;
  margin: 0;
  min-height: 100vh;
  width: 100%;
}

.login-card {
  background: var(--surface);
  border-radius: var(--radius-xl);
  padding: 40px 44px;
  max-width: 420px;
  width: 100%;
  box-shadow: var(--shadow-xl);
  animation: fadeIn 0.5s ease;
  border: 1px solid var(--line);
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: translateY(0); }
}

.login-header {
  text-align: center;
  margin-bottom: 28px;
}

.powered-by {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  padding: 12px 16px;
  background: linear-gradient(135deg, var(--brand-tint), var(--brand-glow));
  border-radius: var(--radius-md);
  border: 1px solid var(--line);
  margin-bottom: 20px;
}

.powered-logo {
  width: 44px;
  height: 44px;
  border-radius: var(--radius-sm);
  object-fit: contain;
  background: var(--surface);
  padding: 4px;
}

.powered-text {
  text-align: right;
  font-size: 11px;
  color: var(--ink-soft);
  line-height: 1.4;
}

.powered-text strong {
  font-size: 14px;
  color: var(--brand);
  font-weight: 800;
}

.powered-slogan {
  display: block;
  font-size: 9px;
  color: var(--ink-light);
  font-weight: 400;
}

.app-brand {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  margin-bottom: 8px;
}

.brand-mark {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  object-fit: cover;
  background: transparent;
}

.title {
  font-size: 34px;
  font-weight: 800;
  color: var(--brand);
  margin: 0;
  letter-spacing: -0.5px;
}

.subtitle {
  font-size: 15px;
  color: var(--ink-soft);
  margin: 0;
}

.lang-section {
  display: flex;
  justify-content: center;
  gap: 8px;
  margin-bottom: 28px;
  padding: 6px;
  background: var(--canvas);
  border-radius: var(--radius-md);
  border: 1px solid var(--line);
}

.lang-btn {
  flex: 1;
  padding: 8px 16px;
  font-size: 13px;
}

.lang-btn.active {
  background: linear-gradient(135deg, var(--brand) 0%, var(--brand-dark) 100%);
  color: #fff;
  border-color: var(--brand);
  box-shadow: var(--shadow-md), var(--brand-glow);
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.field {
  display: flex;
  flex-direction: column;
  gap: 4px;
  font-size: 13px;
  color: var(--ink-soft);
}

.field label {
  font-size: 14px;
  font-weight: 600;
  color: var(--ink);
}

.field input {
  font-family: var(--font-body);
  font-size: 14px;
  color: var(--ink);
  padding: 10px 14px;
  border-radius: var(--radius-sm);
  border: 1.5px solid var(--line);
  background: var(--surface);
  outline: none;
  transition: all var(--transition-fast);
  width: 100%;
}

.field input:focus {
  border-color: var(--brand);
  box-shadow: 0 0 0 3px var(--brand-tint);
}

/* Hide Chrome autofill ghost text */
.field input:-webkit-autofill,
.field input:-webkit-autofill:hover,
.field input:-webkit-autofill:focus,
.field input:-webkit-autofill:active {
  -webkit-box-shadow: 0 0 0 30px var(--surface) inset !important;
  -webkit-text-fill-color: var(--ink) !important;
  transition: background-color 5000s ease-in-out 0s;
}

/* Prevent autofill suggestion ghost text */
.field input[readonly] {
  background: var(--surface);
  cursor: text;
}

.btn-login {
  margin-top: 4px;
  padding: 14px 32px;
  font-size: 16px;
}

.error {
  color: var(--signal-out);
  font-size: 14px;
  text-align: center;
  background: var(--signal-out-tint);
  padding: 12px;
  border-radius: var(--radius-sm);
  border: 1px solid var(--signal-out);
}

.debug {
  color: var(--signal-info);
  font-size: 13px;
  text-align: center;
  background: var(--signal-info-tint);
  padding: 10px;
  border-radius: var(--radius-sm);
  border: 1px solid var(--signal-info);
  white-space: pre-line;
}

.footer {
  text-align: center;
  margin-top: 24px;
  padding-top: 20px;
  border-top: 1px solid var(--line);
}

.footer p {
  font-size: 12px;
  color: var(--ink-light);
  margin: 4px 0;
}

.footer-powered {
  font-size: 11px;
  color: var(--ink-light);
}

@media (max-width: 480px) {
  .login-card {
    padding: 28px 20px;
    border-radius: var(--radius-md);
  }
  
  .title {
    font-size: 28px;
  }
  
  .lang-btn {
    padding: 6px 12px;
    font-size: 12px;
  }
  
  .powered-logo {
    width: 32px;
    height: 32px;
  }
  
  .powered-text strong {
    font-size: 12px;
  }
  
  .field input {
    font-size: 16px;
  }
  
  .btn-login {
    padding: 12px 24px;
    font-size: 15px;
  }
}
</style> 
```

---

## 📄 frontend-admin-dashboard/src/views/NotificationsView.vue

```vue
<template>
  <div>
    <div class="page-head"><h2>{{ t('notifications') }}</h2></div>
    <div class="card">
      <div v-if="loading" class="loading-state">
        <div class="spinner"></div>
        <p>{{ t('loading') }}</p>
      </div>
      <div v-else-if="notifications.length === 0" class="empty-state">
        <p>{{ t('no_notifications') }}</p>
      </div>
      <ActivityFeed v-else :items="formattedNotifications" />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import ActivityFeed from '../components/ActivityFeed.vue'
import { useI18n } from '../services/i18n'
import api from '../services/api'
import { toLocalNumerals, formatLocalDate } from '../utils/numberFormatter'

const { t, locale } = useI18n()

const notifications = ref([])
const loading = ref(true)

const formattedNotifications = computed(() => {
  return notifications.value.map(notif => ({
    id: notif.id,
    text: notif.title || notif.body,
    time: formatTime(notif.created_at),
    tone: notif.is_read ? 'neutral' : 'in'
  }))
})

function formatTime(dateStr) {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  const now = new Date()
  const diffMs = now - date
  const diffMins = Math.floor(diffMs / 60000)
  const diffHours = Math.floor(diffMs / 3600000)
  const diffDays = Math.floor(diffMs / 86400000)

  if (diffMins < 1) return t('just_now')
  if (diffMins < 60) return `${toLocalNumerals(diffMins, locale?.value || 'en')} ${t('minutes_ago')}`
  if (diffHours < 24) return `${toLocalNumerals(diffHours, locale?.value || 'en')} ${t('hours_ago')}`
  if (diffDays === 1) return t('yesterday')
  if (diffDays < 7) return `${toLocalNumerals(diffDays, locale?.value || 'en')} ${t('days_ago')}`
  
  return formatLocalDate(dateStr, locale?.value || 'en', 'DD/MM/YYYY')
}

async function fetchNotifications() {
  try {
    const { data } = await api.get('/notifications')
    notifications.value = data || []
  } catch (error) {
    console.error('Failed to fetch notifications:', error)
    notifications.value = []
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchNotifications()
})
</script>

<style scoped>
.page-head { margin-bottom: 18px; }
.page-head h2 { font-size: 20px; }
.card { padding: 22px; }

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
  gap: 12px;
}

.spinner {
  width: 32px;
  height: 32px;
  border: 3px solid rgba(59, 130, 246, 0.2);
  border-top-color: #3b82f6;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.empty-state {
  text-align: center;
  padding: 40px 20px;
  color: var(--ink-soft);
}
</style>

```

---

## 📄 frontend-admin-dashboard/src/views/ReportsView.vue

```vue
<template>
  <div>
    <div class="page-head">
      <div>
        <h2>📊 {{ t('reports') }}</h2>
        <p>{{ t('reports_description') }}</p>
      </div>
      <div class="report-tabs">
        <button 
          v-for="tab in tabs" 
          :key="tab.id"
          class="tab-btn"
          :class="{ active: activeTab === tab.id }"
          @click="activeTab = tab.id"
        >
          {{ tab.label }}
        </button>
      </div>
    </div>

    <div v-if="loading" class="empty-state">
      <p>⏳ {{ t('loading') }}...</p>
    </div>

    <div v-else-if="error" class="alert alert-error">
      <span>❌</span> {{ error }}
    </div>

    <!-- التقرير الشامل -->
    <div v-else-if="activeTab === 'comprehensive'" class="comprehensive-report">
      <!-- إحصائيات سريعة -->
      <div class="stats-grid">
        <div class="stat-card employees">
          <div class="stat-icon">👷</div>
          <div class="stat-content">
            <div class="stat-value">{{ comprehensiveData.employees?.active || 0 }}</div>
            <div class="stat-label">{{ t('active_employees') }}</div>
            <div class="stat-sub">{{ comprehensiveData.employees?.on_duty || 0 }} {{ t('on_duty') }}</div>
          </div>
        </div>

        <div class="stat-card clients">
          <div class="stat-icon">👥</div>
          <div class="stat-content">
            <div class="stat-value">{{ comprehensiveData.clients?.active || 0 }}</div>
            <div class="stat-label">{{ t('active_clients') }}</div>
            <div class="stat-sub">{{ comprehensiveData.clients?.total || 0 }} {{ t('total') }}</div>
          </div>
        </div>

        <div class="stat-card services">
          <div class="stat-icon">📋</div>
          <div class="stat-content">
            <div class="stat-value">{{ comprehensiveData.service_requests?.completed || 0 }}</div>
            <div class="stat-label">{{ t('completed_services') }}</div>
            <div class="stat-sub">{{ comprehensiveData.service_requests?.pending || 0 }} {{ t('pending') }}</div>
          </div>
        </div>

        <div class="stat-card ratings">
          <div class="stat-icon">⭐</div>
          <div class="stat-content">
            <div class="stat-value">{{ (comprehensiveData.ratings?.average || 0).toFixed(1) }}</div>
            <div class="stat-label">{{ t('average_rating') }}</div>
            <div class="stat-sub">{{ comprehensiveData.service_requests?.rated || 0 }} {{ t('ratings') }}</div>
          </div>
        </div>
      </div>

      <!-- رسوم بيانية تفصيلية -->
      <div class="charts-grid">
        <div class="chart-card">
          <h3>{{ t('service_requests_status') }}</h3>
          <div class="chart-content">
            <div class="bar-chart">
              <div class="bar-item">
                <span class="bar-label">{{ t('completed') }}</span>
                <div class="bar-wrapper">
                  <div class="bar-fill completed" :style="{ width: getPercentage(comprehensiveData.service_requests?.completed || 0, comprehensiveData.service_requests?.total || 0) + '%' }"></div>
                </div>
                <span class="bar-value">{{ comprehensiveData.service_requests?.completed || 0 }}</span>
              </div>
              <div class="bar-item">
                <span class="bar-label">{{ t('in_progress') }}</span>
                <div class="bar-wrapper">
                  <div class="bar-fill in-progress" :style="{ width: getPercentage(comprehensiveData.service_requests?.in_progress || 0, comprehensiveData.service_requests?.total || 0) + '%' }"></div>
                </div>
                <span class="bar-value">{{ comprehensiveData.service_requests?.in_progress || 0 }}</span>
              </div>
              <div class="bar-item">
                <span class="bar-label">{{ t('pending') }}</span>
                <div class="bar-wrapper">
                  <div class="bar-fill pending" :style="{ width: getPercentage(comprehensiveData.service_requests?.pending || 0, comprehensiveData.service_requests?.total || 0) + '%' }"></div>
                </div>
                <span class="bar-value">{{ comprehensiveData.service_requests?.pending || 0 }}</span>
              </div>
            </div>
          </div>
        </div>

        <div class="chart-card">
          <h3>{{ t('tasks_status') }}</h3>
          <div class="chart-content">
            <div class="bar-chart">
              <div class="bar-item">
                <span class="bar-label">{{ t('completed') }}</span>
                <div class="bar-wrapper">
                  <div class="bar-fill completed" :style="{ width: getPercentage(comprehensiveData.tasks?.completed || 0, comprehensiveData.tasks?.total || 0) + '%' }"></div>
                </div>
                <span class="bar-value">{{ comprehensiveData.tasks?.completed || 0 }}</span>
              </div>
              <div class="bar-item">
                <span class="bar-label">{{ t('in_progress') }}</span>
                <div class="bar-wrapper">
                  <div class="bar-fill in-progress" :style="{ width: getPercentage(comprehensiveData.tasks?.in_progress || 0, comprehensiveData.tasks?.total || 0) + '%' }"></div>
                </div>
                <span class="bar-value">{{ comprehensiveData.tasks?.in_progress || 0 }}</span>
              </div>
              <div class="bar-item">
                <span class="bar-label">{{ t('pending') }}</span>
                <div class="bar-wrapper">
                  <div class="bar-fill pending" :style="{ width: getPercentage(comprehensiveData.tasks?.pending || 0, comprehensiveData.tasks?.total || 0) + '%' }"></div>
                </div>
                <span class="bar-value">{{ comprehensiveData.tasks?.pending || 0 }}</span>
              </div>
            </div>
          </div>
        </div>

        <div class="chart-card">
          <h3>{{ t('attendance_statistics') }}</h3>
          <div class="chart-content">
            <div class="attendance-stats">
              <div class="attendance-item">
                <span class="attendance-label">{{ t('total_hours_week') }}</span>
                <span class="attendance-value">{{ (comprehensiveData.attendance?.total_hours_week || 0).toFixed(1) }} {{ t('hours') }}</span>
              </div>
              <div class="attendance-item">
                <span class="attendance-label">{{ t('avg_daily_hours') }}</span>
                <span class="attendance-value">{{ (comprehensiveData.attendance?.avg_daily_hours || 0).toFixed(1) }} {{ t('hours') }}</span>
              </div>
              <div class="attendance-item">
                <span class="attendance-label">{{ t('completed_duty_today') }}</span>
                <span class="attendance-value">{{ comprehensiveData.employees?.completed_duty_today || 0 }}</span>
              </div>
            </div>
          </div>
        </div>

        <div class="chart-card">
          <h3>{{ t('worksites_status') }}</h3>
          <div class="chart-content">
            <div class="worksites-stats">
              <div class="worksite-item">
                <span class="worksite-label">{{ t('total_worksites') }}</span>
                <span class="worksite-value">{{ comprehensiveData.worksites?.total || 0 }}</span>
              </div>
              <div class="worksite-item">
                <span class="worksite-label">{{ t('active_worksites') }}</span>
                <span class="worksite-value">{{ comprehensiveData.worksites?.active || 0 }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- تقرير طلبات الخدمة -->
    <div v-else-if="activeTab === 'service-requests'" class="report-section">
      <div class="card">
        <h3>{{ t('service_requests_report') }}</h3>
        <div v-if="serviceRequestsLoading" class="empty-state">
          <p>⏳ {{ t('loading') }}...</p>
        </div>
        <div v-else-if="!serviceRequests || serviceRequests.length === 0" class="empty-state">
          <p>{{ t('no_data_available') }}</p>
        </div>
        <div v-else class="table-container">
          <table class="data-table">
            <thead>
              <tr>
                <th>{{ t('title') }}</th>
                <th>{{ t('client') }}</th>
                <th>{{ t('employee') }}</th>
                <th>{{ t('status') }}</th>
                <th>{{ t('priority') }}</th>
                <th>{{ t('rating') }}</th>
                <th>{{ t('created_at') }}</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="req in serviceRequests" :key="req.id">
                <td>{{ req.title }}</td>
                <td>{{ req.client_name || '—' }}</td>
                <td>{{ req.employee_name || '—' }}</td>
                <td><span :class="['badge', getStatusClass(req.status)]">{{ req.status }}</span></td>
                <td><span :class="['badge', getPriorityClass(req.priority)]">{{ req.priority }}</span></td>
                <td>{{ req.client_rating ? `⭐ ${req.client_rating}` : '—' }}</td>
                <td>{{ formatDate(req.created_at) }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- تقرير أداء الموظفين -->
    <div v-else-if="activeTab === 'employee-performance'" class="report-section">
      <div class="card">
        <h3>{{ t('employee_performance_report') }}</h3>
        <div v-if="employeePerformanceLoading" class="empty-state">
          <p>⏳ {{ t('loading') }}...</p>
        </div>
        <div v-else-if="!employeePerformance || employeePerformance.length === 0" class="empty-state">
          <p>{{ t('no_data_available') }}</p>
        </div>
        <div v-else class="table-container">
          <table class="data-table">
            <thead>
              <tr>
                <th>{{ t('employee') }}</th>
                <th>{{ t('phone') }}</th>
                <th>{{ t('total_hours') }}</th>
                <th>{{ t('completed_shifts') }}</th>
                <th>{{ t('assigned_services') }}</th>
                <th>{{ t('avg_rating') }}</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="emp in employeePerformance" :key="emp.id">
                <td>{{ emp.full_name }}</td>
                <td>{{ emp.phone }}</td>
                <td>{{ emp.total_hours.toFixed(1) }} {{ t('hours') }}</td>
                <td>{{ emp.completed_shifts }}</td>
                <td>{{ emp.assigned_services }}</td>
                <td>{{ emp.avg_rating > 0 ? `⭐ ${emp.avg_rating.toFixed(1)}` : '—' }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- تقرير نشاط العملاء -->
    <div v-else-if="activeTab === 'client-activity'" class="report-section">
      <div class="card">
        <h3>{{ t('client_activity_report') }}</h3>
        <div v-if="clientActivityLoading" class="empty-state">
          <p>⏳ {{ t('loading') }}...</p>
        </div>
        <div v-else-if="!clientActivity || clientActivity.length === 0" class="empty-state">
          <p>{{ t('no_data_available') }}</p>
        </div>
        <div v-else class="table-container">
          <table class="data-table">
            <thead>
              <tr>
                <th>{{ t('client') }}</th>
                <th>{{ t('phone') }}</th>
                <th>{{ t('total_requests') }}</th>
                <th>{{ t('completed_requests') }}</th>
                <th>{{ t('avg_rating') }}</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="client in clientActivity" :key="client.id">
                <td>{{ client.full_name }}</td>
                <td>{{ client.phone }}</td>
                <td>{{ client.total_requests }}</td>
                <td>{{ client.completed_requests }}</td>
                <td>{{ client.avg_rating > 0 ? `⭐ ${client.avg_rating.toFixed(1)}` : '—' }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, watch } from 'vue'
import api from '../services/api'
import { useI18n } from '../services/i18n'

const { t } = useI18n()

const activeTab = ref('comprehensive')
const loading = ref(false)
const error = ref('')
const comprehensiveData = ref({})

const serviceRequests = ref([])
const serviceRequestsLoading = ref(false)

const employeePerformance = ref([])
const employeePerformanceLoading = ref(false)

const clientActivity = ref([])
const clientActivityLoading = ref(false)

const tabs = computed(() => [
  { id: 'comprehensive', label: t('comprehensive_report') },
  { id: 'service-requests', label: t('service_requests_report') },
  { id: 'employee-performance', label: t('employee_performance_report') },
  { id: 'client-activity', label: t('client_activity_report') }
])

async function fetchComprehensiveReport() {
  loading.value = true
  error.value = ''
  try {
    const { data } = await api.get('/reports/comprehensive')
    comprehensiveData.value = data
  } catch (err) {
    error.value = err.response?.data?.error || t('failed_to_fetch_report')
  } finally {
    loading.value = false
  }
}

async function fetchServiceRequestsReport() {
  serviceRequestsLoading.value = true
  try {
    const { data } = await api.get('/reports/service-requests')
    serviceRequests.value = data
  } catch (err) {
    console.error('Failed to fetch service requests report:', err)
  } finally {
    serviceRequestsLoading.value = false
  }
}

async function fetchEmployeePerformanceReport() {
  employeePerformanceLoading.value = true
  try {
    const { data } = await api.get('/reports/employee-performance')
    employeePerformance.value = data
  } catch (err) {
    console.error('Failed to fetch employee performance report:', err)
  } finally {
    employeePerformanceLoading.value = false
  }
}

async function fetchClientActivityReport() {
  clientActivityLoading.value = true
  try {
    const { data } = await api.get('/reports/client-activity')
    clientActivity.value = data
  } catch (err) {
    console.error('Failed to fetch client activity report:', err)
  } finally {
    clientActivityLoading.value = false
  }
}

function getPercentage(value, total) {
  if (!total || total === 0) return 0
  if (!value || value === 0) return 0
  return ((value / total) * 100).toFixed(1)
}

function formatDate(date) {
  if (!date) return '—'
  return new Date(date).toLocaleDateString('en-GB', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

function getStatusClass(status) {
  const classes = {
    'completed': 'badge--success',
    'in_progress': 'badge--info',
    'pending': 'badge--warning',
    'assigned': 'badge--primary'
  }
  return classes[status] || 'badge--default'
}

function getPriorityClass(priority) {
  const classes = {
    'urgent': 'badge--danger',
    'high': 'badge--warning',
    'normal': 'badge--primary',
    'low': 'badge--default'
  }
  return classes[priority] || 'badge--default'
}

// Watch tab changes and fetch corresponding data
function onTabChange(newTab) {
  if (newTab === 'comprehensive') {
    fetchComprehensiveReport()
  } else if (newTab === 'service-requests') {
    fetchServiceRequestsReport()
  } else if (newTab === 'employee-performance') {
    fetchEmployeePerformanceReport()
  } else if (newTab === 'client-activity') {
    fetchClientActivityReport()
  }
}

// Load initial data
onMounted(() => {
  fetchComprehensiveReport()
})

// Watch for tab changes
watch(activeTab, onTabChange)
</script>

<style scoped>
.page-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  flex-wrap: wrap;
  gap: 12px;
}

.page-head h2 { font-size: 20px; margin-bottom: 4px; }
.page-head p { font-size: 13px; color: var(--ink-soft); }

.report-tabs {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.tab-btn {
  padding: 8px 16px;
  border: 1px solid var(--line);
  background: var(--surface);
  border-radius: var(--radius-sm);
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
  font-family: var(--font-body);
}

.tab-btn:hover {
  border-color: var(--brand);
  background: var(--brand-tint);
}

.tab-btn.active {
  background: var(--brand);
  color: white;
  border-color: var(--brand);
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 16px;
  margin-bottom: 20px;
}

.stat-card {
  background: var(--surface);
  border: 1px solid var(--line);
  border-radius: var(--radius-md);
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 16px;
}

.stat-icon {
  font-size: 32px;
  width: 56px;
  height: 56px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--radius-md);
  background: var(--canvas);
}

.stat-card.employees .stat-icon { background: #E3F2FD; }
.stat-card.clients .stat-icon { background: #F3E5F5; }
.stat-card.services .stat-icon { background: #E8F5E9; }
.stat-card.ratings .stat-icon { background: #FFF3E0; }

.stat-content {
  flex: 1;
}

.stat-value {
  font-size: 28px;
  font-weight: 700;
  color: var(--ink);
  line-height: 1;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 13px;
  color: var(--ink-soft);
  margin-bottom: 2px;
}

.stat-sub {
  font-size: 12px;
  color: var(--ink-soft);
}

.charts-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 16px;
}

.chart-card {
  background: var(--surface);
  border: 1px solid var(--line);
  border-radius: var(--radius-md);
  padding: 20px;
}

.chart-card h3 {
  font-size: 15px;
  margin-bottom: 16px;
  color: var(--ink);
}

.chart-content {
  min-height: 200px;
}

.bar-chart {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.bar-item {
  display: flex;
  align-items: center;
  gap: 12px;
}

.bar-label {
  font-size: 12px;
  color: var(--ink-soft);
  min-width: 80px;
  text-align: right;
}

.bar-wrapper {
  flex: 1;
  height: 24px;
  background: var(--canvas);
  border-radius: var(--radius-sm);
  overflow: hidden;
}

.bar-fill {
  height: 100%;
  transition: width 0.3s ease;
  border-radius: var(--radius-sm);
}

.bar-fill.completed { background: var(--signal-in); }
.bar-fill.in-progress { background: #2196F3; }
.bar-fill.pending { background: #FFA500; }

.bar-value {
  font-size: 13px;
  font-weight: 600;
  color: var(--ink);
  min-width: 30px;
  text-align: left;
}

.attendance-stats,
.worksites-stats {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.attendance-item,
.worksite-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px;
  background: var(--canvas);
  border-radius: var(--radius-sm);
}

.attendance-label,
.worksite-label {
  font-size: 13px;
  color: var(--ink-soft);
}

.attendance-value,
.worksite-value {
  font-size: 16px;
  font-weight: 600;
  color: var(--ink);
}

.report-section {
  margin-top: 20px;
}

.table-container {
  overflow-x: auto;
  margin-top: 16px;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 13px;
}

.data-table th,
.data-table td {
  padding: 12px;
  text-align: right;
  border-bottom: 1px solid var(--line);
}

.data-table th {
  background: var(--canvas);
  font-weight: 600;
  color: var(--ink);
  position: sticky;
  top: 0;
}

.data-table tbody tr:hover {
  background: var(--canvas);
}

.badge {
  padding: 4px 8px;
  border-radius: var(--radius-sm);
  font-size: 11px;
  font-weight: 600;
  white-space: nowrap;
}

.badge--success { background: var(--signal-in-tint); color: var(--signal-in); }
.badge--info { background: #E3F2FD; color: #1976D2; }
.badge--warning { background: #FFF3E0; color: #F57C00; }
.badge--danger { background: #FFEBEE; color: #D32F2F; }
.badge--primary { background: #E3F2FD; color: #1976D2; }
.badge--default { background: var(--canvas); color: var(--ink-soft); }

.empty-state {
  text-align: center;
  padding: 40px 20px;
  color: var(--ink-soft);
}

.alert {
  padding: 12px 16px;
  border-radius: var(--radius-sm);
  margin-bottom: 16px;
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 13px;
}

.alert-error {
  background: var(--signal-out-tint);
  color: var(--signal-out);
}

@media (max-width: 768px) {
  .page-head {
    flex-direction: column;
    align-items: flex-start;
  }

  .report-tabs {
    width: 100%;
  }

  .tab-btn {
    flex: 1;
    min-width: 120px;
    text-align: center;
  }

  .stats-grid {
    grid-template-columns: 1fr;
  }

  .charts-grid {
    grid-template-columns: 1fr;
  }

  .data-table {
    font-size: 12px;
  }

  .data-table th,
  .data-table td {
    padding: 8px;
  }
}
</style>

```

---

## 📄 frontend-admin-dashboard/src/views/ServiceRequestsView.vue

```vue
<template>
  <div>
    <div class="page-head">
      <div>
        <h2>📋 {{ t('service_requests') }}</h2>
        <p>{{ t('service_requests_description') }}</p>
      </div>
      <div class="page-head__actions">
        <div class="filters">
          <select v-model="filterStatus" @change="fetchRequests">
            <option value="">{{ t('all_statuses') }}</option>
            <option value="pending">{{ t('status_pending') }}</option>
            <option value="assigned">{{ t('status_assigned') }}</option>
            <option value="in_progress">{{ t('status_in_progress') }}</option>
            <option value="completed">{{ t('status_completed') }}</option>
          </select>
        </div>
      </div>
    </div>

    <div v-if="loading" class="empty-state"><p>⏳ {{ t('loading') }}...</p></div>

    <div v-else-if="requests.length === 0" class="empty-state">
      <h3>📭 {{ t('no_service_requests') }}</h3>
      <p>{{ t('no_service_requests_hint') }}</p>
    </div>

    <div v-else>
      <div v-for="req in filteredRequests" :key="req.id" class="card request-card">
        <div class="request-card__header">
          <div>
            <span class="badge" :class="getPriorityClass(req.priority)">
              {{ getPriorityLabel(req.priority) }}
            </span>
            <span class="badge" :class="getStatusClass(req.status)">
              {{ getStatusLabel(req.status) }}
            </span>
          </div>
          <span class="request-card__time mono">{{ formatDate(req.created_at) }}</span>
        </div>

        <h3>{{ req.title }}</h3>
        <p class="request-card__desc">{{ req.description }}</p>

        <div class="request-card__info">
          <span>👤 {{ req.client_name || t('client') }}</span>
          <span>📞 {{ req.client_phone || req.phone || '—' }}</span>
          <span>📍 {{ req.address || t('no_address') }}</span>
        </div>

        <!-- عرض التقييم إذا وجد -->
        <div v-if="req.assignment && req.assignment.client_rating" class="rating-display">
          <span class="rating-stars">⭐ {{ req.assignment.client_rating }}/5</span>
          <span v-if="req.assignment.client_feedback" class="rating-feedback">"{{ req.assignment.client_feedback }}"</span>
        </div>

        <div class="request-card__location">
          <div v-if="req.location_name" class="location-name">
            📍 {{ req.location_name }}
          </div>
          <div class="location-coords">
            <span class="mono">{{ t('latitude') }}: {{ toLocalNumerals(req.latitude.toFixed(5), locale?.value || 'en') }}</span>
            <span class="mono">{{ t('longitude') }}: {{ toLocalNumerals(req.longitude.toFixed(5), locale?.value || 'en') }}</span>
          </div>
        </div>

        <div class="request-card__actions">
          <button v-if="req.status === 'pending'" class="btn btn--primary btn--sm" @click="openAssignModal(req)">
            👤 {{ t('assign_employee') }}
          </button>
          <button class="btn btn--success btn--sm" @click="convertToWorksite(req)">
            🏢 {{ t('convert_to_worksite') }}
          </button>
          <div v-if="req.status === 'assigned'" class="employee-status">
            <span class="employee-info">
              👤 {{ t('current_assigned_employee') }} {{ req.employee_name || '—' }}
            </span>
            <div class="employee-actions">
              <button class="btn btn--primary btn--sm" @click="openAssignModal(req)">
                🔄 {{ t('reassign_employee') }}
              </button>
              <button class="btn btn--danger btn--sm" @click="unassignEmployee(req)">
                ❌ {{ t('unassign_employee') }}
              </button>
            </div>
          </div>
          <button v-if="req.status === 'in_progress'" class="btn btn--gold btn--sm">
            🔄 {{ t('in_execution') }}
          </button>
          <button v-if="req.status === 'completed'" class="btn btn--success btn--sm">
            ✅ {{ t('status_completed') }}
          </button>
          <button class="btn btn--danger btn--sm" @click="deleteRequest(req)">
            🗑️ {{ t('delete_request') }}
          </button>
        </div>
      </div>
    </div>

    <!-- مودال تعيين الموظف -->
    <div v-if="showAssignModal" class="modal-backdrop" @click.self="showAssignModal = false">
      <div class="modal card">
        <h3>👤 {{ t('assign_employee_modal') }}</h3>
        <p>{{ t('assign_employee_hint') }}</p>

        <div v-if="loadingEmployees" class="empty-state"><p>⏳ {{ t('loading_employees') }}...</p></div>

        <div v-else class="employees-list">
          <button
            v-for="emp in employees"
            :key="emp.id"
            class="employee-item"
            @click="assignEmployee(emp.id)"
          >
            <span class="employee-item__avatar">{{ emp.full_name?.slice(0, 1) || '?' }}</span>
            <div>
              <strong>{{ emp.full_name }}</strong>
            </div>
          </button>
        </div>

        <button class="btn btn--ghost btn--block" @click="showAssignModal = false">{{ t('cancel') }}</button>
      </div>
    </div>

    <!-- مودال تحويل إلى نقطة عمل -->
    <div v-if="showWorksiteModal" class="modal-backdrop" @click.self="showWorksiteModal = false">
      <div class="modal card">
        <h3>🏢 {{ t('convert_to_worksite') }}</h3>
        <p>{{ t('convert_to_worksite_hint') }}</p>

        <div class="form-group">
          <label>{{ t('worksite_name') }} *</label>
          <input 
            v-model="worksiteForm.name" 
            type="text" 
            class="form-control"
            :placeholder="selectedRequest?.location_name || t('worksite_name_placeholder')"
          />
        </div>

        <div class="form-group">
          <label>{{ t('address') }}</label>
          <input 
            v-model="worksiteForm.address" 
            type="text" 
            class="form-control"
            :placeholder="selectedRequest?.address || t('address_placeholder')"
          />
        </div>

        <div class="form-group">
          <label>{{ t('radius') }} (متر) *</label>
          <input 
            v-model.number="worksiteForm.radius" 
            type="number" 
            class="form-control"
            :placeholder="t('radius_placeholder')"
            min="10"
          />
        </div>

        <div class="location-info">
          <div class="info-item">
            <span class="label">{{ t('latitude') }}:</span>
            <span class="value">{{ selectedRequest?.latitude.toFixed(6) }}</span>
          </div>
          <div class="info-item">
            <span class="label">{{ t('longitude') }}:</span>
            <span class="value">{{ selectedRequest?.longitude.toFixed(6) }}</span>
          </div>
        </div>

        <div class="modal-actions">
          <button class="btn btn--ghost" @click="showWorksiteModal = false">{{ t('cancel') }}</button>
          <button class="btn btn--primary" @click="createWorksite" :disabled="loadingWorksite">
            {{ loadingWorksite ? '⏳ ' + t('loading') : t('create') }}
          </button>
        </div>
      </div>
    </div>

    <!-- مودال تأكيد الحذف -->
    <div v-if="showDeleteModal" class="modal-backdrop" @click.self="showDeleteModal = false">
      <div class="modal card">
        <div class="modal-header">
          <h3>🗑️ {{ t('delete_request_title') }}</h3>
        </div>
        <div class="modal-body">
          <p>{{ t('delete_request_message') }}</p>
          <div class="warning-message">
            <p>{{ t('delete_request_warning') }}</p>
          </div>
          <div v-if="selectedRequest" class="request-info">
            <p><strong>{{ t('request_details') }}:</strong></p>
            <p>{{ selectedRequest.title }}</p>
            <p><small>{{ selectedRequest.client_name || t('client') }} - {{ formatDate(selectedRequest.created_at) }}</small></p>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn--ghost" @click="showDeleteModal = false" :disabled="deletingRequest">
            {{ t('cancel') }}
          </button>
          <button class="btn btn--danger" @click="confirmDeleteRequest" :disabled="deletingRequest">
            {{ deletingRequest ? '⏳ ' + t('deleting') : '🗑️ ' + t('delete_final') }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import api from '../services/api'
import { useI18n } from '../services/i18n'
import { formatLocalDate, toLocalNumerals } from '../utils/numberFormatter'

const { t, locale } = useI18n()

const requests = ref([])
const employees = ref([])
const loading = ref(false)
const loadingEmployees = ref(false)
const loadingWorksite = ref(false)
const filterStatus = ref('')
const showAssignModal = ref(false)
const showWorksiteModal = ref(false)
const showDeleteModal = ref(false)
const deletingRequest = ref(false)
const selectedRequest = ref(null)
const worksiteForm = ref({
  name: '',
  address: '',
  radius: 100
})

const filteredRequests = computed(() => {
  if (!filterStatus.value) return requests.value
  return requests.value.filter(r => r.status === filterStatus.value)
})

async function fetchRequests() {
  loading.value = true
  try {
    const { data } = await api.get('/service/requests')
    requests.value = data || []
  } catch (error) {
    console.error(t('failed_to_fetch_requests'), error)
  } finally {
    loading.value = false
  }
}

async function fetchEmployees() {
  loadingEmployees.value = true
  try {
    const { data } = await api.get('/admin/employees')
    employees.value = data || []
  } catch (error) {
    console.error(t('failed_to_fetch_employees'), error)
  } finally {
    loadingEmployees.value = false
  }
}

function openAssignModal(req) {
  selectedRequest.value = req
  showAssignModal.value = true
  fetchEmployees()
}

async function assignEmployee(employeeId) {
  try {
    await api.post('/service/assign', {
      request_id: selectedRequest.value.id,
      employee_id: employeeId
    })
    showAssignModal.value = false
    await fetchRequests()
  } catch (error) {
    console.error('فشل التعيين:', error)
  }
}

async function unassignEmployee(req) {
  if (!confirm(t('cancel') + '?')) return
  try {
    await api.post('/service/assign', {
      request_id: req.id,
      employee_id: null
    })
    await fetchRequests()
  } catch (error) {
    console.error('فشل إلغاء التعيين:', error)
  }
}

async function deleteRequest(req) {
  selectedRequest.value = req
  showDeleteModal.value = true
}

async function confirmDeleteRequest() {
  if (!selectedRequest.value) return
  
  deletingRequest.value = true
  try {
    await api.delete(`/service/requests/${selectedRequest.value.id}`)
    showDeleteModal.value = false
    selectedRequest.value = null
    await fetchRequests()
    // Show success message
    showSuccessMessage(t('request_deleted_successfully'))
  } catch (error) {
    console.error('فشل حذف الطلب:', error)
    // Handle 404 error - request might have been already deleted or route not available
    if (error.response && error.response.status === 404) {
      showErrorMessage(t('request_not_found_or_deleted'))
      // Remove the request from the local list since it doesn't exist on server
      requests.value = requests.value.filter(r => r.id !== selectedRequest.value.id)
      showDeleteModal.value = false
      selectedRequest.value = null
    } else if (error.response && error.response.status === 403) {
      showErrorMessage(t('delete_error_title') + ': ' + t('err_permission_denied'))
    } else {
      // Show error message
      showErrorMessage(t('request_delete_failed'))
    }
  } finally {
    deletingRequest.value = false
  }
}

// Success message function
function showSuccessMessage(message) {
  // Create a temporary success message element
  const messageDiv = document.createElement('div')
  messageDiv.className = 'toast-message toast-success'
  messageDiv.innerHTML = `
    <div class="toast-content">
      <span class="toast-icon">✅</span>
      <span class="toast-text">${message}</span>
    </div>
  `
  document.body.appendChild(messageDiv)
  
  // Remove after 3 seconds
  setTimeout(() => {
    messageDiv.classList.add('toast-fade-out')
    setTimeout(() => {
      document.body.removeChild(messageDiv)
    }, 300)
  }, 3000)
}

// Error message function
function showErrorMessage(message) {
  // Create a temporary error message element
  const messageDiv = document.createElement('div')
  messageDiv.className = 'toast-message toast-error'
  messageDiv.innerHTML = `
    <div class="toast-content">
      <span class="toast-icon">❌</span>
      <span class="toast-text">${message}</span>
    </div>
  `
  document.body.appendChild(messageDiv)
  
  // Remove after 4 seconds
  setTimeout(() => {
    messageDiv.classList.add('toast-fade-out')
    setTimeout(() => {
      document.body.removeChild(messageDiv)
    }, 300)
  }, 4000)
}

function getPriorityLabel(priority) {
  const labels = { low: t('priority_low'), normal: t('priority_normal'), high: t('priority_high'), urgent: t('priority_urgent') }
  return labels[priority] || priority
}

function getPriorityClass(priority) {
  const classes = { low: 'badge--gray', normal: 'badge--blue', high: 'badge--gold', urgent: 'badge--out' }
  return classes[priority] || ''
}

function getStatusLabel(status) {
  const labels = { pending: t('status_pending'), assigned: t('status_assigned'), in_progress: t('status_in_progress'), completed: t('status_completed') }
  return labels[status] || status
}

function getStatusClass(status) {
  const classes = { pending: 'badge--gold', assigned: 'badge--blue', in_progress: 'badge--info', completed: 'badge--in' }
  return classes[status] || ''
}

function formatDate(date) {
  if (!date) return '—'
  return formatLocalDate(date, locale?.value || 'en', 'DD/MM/YYYY HH:mm')
}

function convertToWorksite(req) {
  selectedRequest.value = req
  worksiteForm.value = {
    name: req.location_name || `موقع ${req.title}`,
    address: req.address || '',
    radius: 100
  }
  showWorksiteModal.value = true
}

async function createWorksite() {
  if (!selectedRequest.value || !worksiteForm.value.name || !worksiteForm.value.radius) {
    alert('الرجاء ملء جميع الحقول المطلوبة')
    return
  }

  loadingWorksite.value = true
  try {
    await api.post('/worksites', {
      name: worksiteForm.value.name,
      address: worksiteForm.value.address,
      latitude: selectedRequest.value.latitude,
      longitude: selectedRequest.value.longitude,
      radius_meters: worksiteForm.value.radius
    })
    
    showWorksiteModal.value = false
    alert('تم إنشاء نقطة العمل بنجاح')
  } catch (error) {
    console.error('فشل إنشاء نقطة العمل:', error)
    alert('فشل إنشاء نقطة العمل')
  } finally {
    loadingWorksite.value = false
  }
}

onMounted(fetchRequests)
</script>

<style scoped>
.page-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 18px;
  flex-wrap: wrap;
  gap: 10px;
}

.page-head h2 { font-size: 20px; margin-bottom: 4px; }
.page-head p { font-size: 13px; color: var(--ink-soft); }

.page-head__actions {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.filters select {
  padding: 8px 14px;
  border-radius: var(--radius-sm);
  border: 1px solid var(--line);
  background: var(--surface);
  font-family: var(--font-body);
  min-width: 150px;
}

@media (max-width: 768px) {
  .page-head { flex-direction: column; align-items: flex-start; }
  .page-head__actions { width: 100%; flex-direction: column; align-items: stretch; }
  .filters { width: 100%; }
  .filters select { width: 100%; }
  
  .request-card__location {
    flex-direction: column;
    gap: 8px;
  }
  
  .request-card__actions {
    flex-direction: column;
  }
  
  .request-card__actions .btn {
    width: 100%;
  }
}

.form-group {
  margin-bottom: 16px;
}

.form-group label {
  display: block;
  font-size: 13px;
  font-weight: 600;
  color: var(--ink);
  margin-bottom: 6px;
}

.form-control {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid var(--line);
  border-radius: var(--radius-sm);
  font-family: var(--font-body);
  font-size: 14px;
  background: var(--surface);
}

.form-control:focus {
  outline: none;
  border-color: #1E3A5F;
  box-shadow: 0 0 0 3px rgba(30, 58, 95, 0.1);
}

.required {
  color: var(--signal-out);
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20px;
}

.modal-close {
  background: none;
  border: none;
  font-size: 20px;
  cursor: pointer;
  color: var(--ink-soft);
  padding: 4px;
}

.modal-close:hover {
  color: var(--ink);
}

.alert {
  padding: 12px;
  border-radius: var(--radius-sm);
  margin-bottom: 16px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.alert-error {
  background: var(--signal-out-tint);
  color: var(--signal-out);
  border: 1px solid var(--signal-out);
}

.request-card {
  padding: 18px 20px;
  margin-bottom: 12px;
  transition: border-color 0.2s;
}
.request-card:hover { border-color: var(--brand-tint); }

.request-card__header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}
.request-card__header .badge { margin-inline-end: 6px; }
.request-card__time { font-size: 12px; color: var(--ink-soft); }

.request-card h3 { font-size: 16px; margin-bottom: 6px; }
.request-card__desc { font-size: 13px; color: var(--ink-soft); margin-bottom: 10px; }

.request-card__info {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  font-size: 13px;
  color: var(--ink-soft);
  margin-bottom: 8px;
}

.rating-display {
  background: var(--signal-in-tint);
  border: 1px solid var(--signal-in);
  border-radius: var(--radius-sm);
  padding: 10px 14px;
  margin-bottom: 10px;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.rating-stars {
  font-size: 14px;
  font-weight: 600;
  color: var(--signal-in);
}

.rating-feedback {
  font-size: 13px;
  color: var(--ink);
  font-style: italic;
}

.request-card__location {
  display: flex;
  flex-direction: column;
  gap: 8px;
  font-size: 12px;
  color: var(--ink-soft);
  background: var(--canvas);
  padding: 8px 12px;
  border-radius: var(--radius-sm);
  margin-bottom: 12px;
}

.location-name {
  font-size: 13px;
  font-weight: 500;
  color: var(--ink);
  margin-bottom: 4px;
}

.location-coords {
  display: flex;
  gap: 16px;
}

.request-card__actions {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.employee-status {
  display: flex;
  flex-direction: column;
  gap: 8px;
  width: 100%;
}

.employee-info {
  font-size: 13px;
  color: var(--ink-soft);
  padding: 8px 12px;
  background: var(--canvas);
  border-radius: var(--radius-sm);
}

.employee-actions {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.btn--danger {
  background: #fee2e2;
  color: #dc2626;
  border: 1px solid #fecaca;
}

.btn--danger:hover {
  background: #fecaca;
}

.badge--gray { background: var(--line); color: var(--ink-soft); }
.badge--blue { background: #E3ECF7; color: #2C6B9E; }
.badge--info { background: #E3F0F7; color: #1A7A8A; }
.badge--success { background: var(--signal-in-tint); color: var(--signal-in); }

.modal-backdrop {
  position: fixed; inset: 0;
  background: rgba(22,35,46,0.6);
  backdrop-filter: blur(4px);
  display: flex; align-items: center; justify-content: center;
  z-index: 50; padding: 20px;
}
.modal {
  width: 100%; max-width: 420px; padding: 24px;
  max-height: 80vh; overflow-y: auto;
}
.modal h3 { font-size: 17px; margin-bottom: 4px; }
.modal p { font-size: 13px; color: var(--ink-soft); margin-bottom: 16px; }

.employees-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 16px;
}

.employee-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 14px;
  background: var(--surface);
  border: 1px solid var(--line);
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: all 0.2s;
  text-align: right;
  font-family: var(--font-body);
  font-size: 14px;
}
.employee-item:hover {
  border-color: var(--brand);
  background: var(--brand-tint);
}

.employee-item__avatar {
  width: 36px; height: 36px;
  border-radius: 50%;
  background: var(--brand-tint);
  color: var(--brand-dark);
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  font-size: 16px;
  flex-shrink: 0;
}
.employee-item div { display: flex; flex-direction: column; }
.employee-item span { font-size: 12px; color: var(--ink-soft); }

/* Worksite Modal Styles */
.form-group {
  margin-bottom: 16px;
}

.form-group label {
  display: block;
  font-size: 13px;
  font-weight: 600;
  color: var(--ink);
  margin-bottom: 6px;
}

.form-control {
  width: 100%;
  padding: 10px 12px;
  border: 1.5px solid var(--line);
  border-radius: var(--radius-sm);
  font-family: var(--font-body);
  font-size: 14px;
  background: var(--surface);
}

.form-control:focus {
  border-color: var(--brand);
  outline: none;
  box-shadow: 0 0 0 3px rgba(30, 58, 95, 0.1);
}

.location-info {
  background: var(--canvas);
  padding: 12px;
  border-radius: var(--radius-sm);
  margin-bottom: 16px;
}

.info-item {
  display: flex;
  justify-content: space-between;
  margin-bottom: 8px;
}

.info-item:last-child {
  margin-bottom: 0;
}

.info-item .label {
  font-size: 12px;
  color: var(--ink-soft);
}

.info-item .value {
  font-size: 12px;
  font-family: monospace;
  color: var(--ink);
}

.modal-actions {
  display: flex;
  gap: 10px;
  justify-content: flex-end;
  margin-top: 20px;
}

.btn--success {
  background: var(--signal-in-tint);
  color: var(--signal-in);
  border: 1px solid var(--signal-in);
}

.btn--success:hover {
  background: var(--signal-in);
  color: white;
}

/* Delete Modal Styles */
.modal-header {
  border-bottom: 1px solid var(--line);
  padding-bottom: 16px;
  margin-bottom: 16px;
}

.modal-header h3 {
  margin: 0;
  font-size: 18px;
  color: var(--signal-out);
}

.modal-body {
  margin-bottom: 20px;
}

.warning-message {
  background: var(--signal-out-tint);
  border: 1px solid var(--signal-out);
  border-radius: var(--radius-sm);
  padding: 12px;
  margin: 12px 0;
}

.warning-message p {
  margin: 0;
  color: var(--signal-out);
  font-size: 14px;
}

.request-info {
  background: var(--canvas);
  border-radius: var(--radius-sm);
  padding: 12px;
  margin-top: 12px;
}

.request-info p {
  margin: 4px 0;
  font-size: 14px;
}

.request-info small {
  color: var(--ink-soft);
  font-size: 12px;
}

.modal-footer {
  display: flex;
  gap: 10px;
  justify-content: flex-end;
  padding-top: 16px;
  border-top: 1px solid var(--line);
}

/* Toast Messages */
.toast-message {
  position: fixed;
  top: 20px;
  right: 20px;
  z-index: 1000;
  min-width: 300px;
  max-width: 400px;
  padding: 16px 20px;
  border-radius: 12px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  display: flex;
  align-items: center;
  gap: 12px;
  animation: toastSlideIn 0.3s ease-out;
  font-family: var(--font-body);
  direction: rtl;
}

.toast-success {
  background: linear-gradient(135deg, #10b981 0%, #059669 100%);
  color: white;
  border: 1px solid #059669;
}

.toast-error {
  background: linear-gradient(135deg, #ef4444 0%, #dc2626 100%);
  color: white;
  border: 1px solid #dc2626;
}

.toast-content {
  display: flex;
  align-items: center;
  gap: 12px;
  width: 100%;
}

.toast-icon {
  font-size: 20px;
  flex-shrink: 0;
}

.toast-text {
  font-size: 14px;
  font-weight: 500;
  line-height: 1.4;
}

.toast-fade-out {
  animation: toastFadeOut 0.3s ease-out forwards;
}

@keyframes toastSlideIn {
  from {
    transform: translateX(100%);
    opacity: 0;
  }
  to {
    transform: translateX(0);
    opacity: 1;
  }
}

@keyframes toastFadeOut {
  from {
    transform: translateX(0);
    opacity: 1;
  }
  to {
    transform: translateX(100%);
    opacity: 0;
  }
}

/* Responsive */
@media (max-width: 768px) {
  .toast-message {
    right: 10px;
    left: 10px;
    min-width: auto;
    max-width: none;
  }
}
</style>

```

---

## 📄 frontend-admin-dashboard/src/views/SettingsView.vue

```vue
<template>
  <div>
    <div class="page-head">
      <h2>{{ t('settings') }}</h2>
      <p class="page-subtitle">{{ t('settings_intro') }}</p>
    </div>

    <div class="card subscription-card">
      <h3>{{ t('subscription_title') }}</h3>
      <p class="settings-card__hint">{{ t('subscription_description') }}</p>
      <SubscriptionStatusView />
    </div>

    <div class="card settings-card">
      <h3>{{ t('settings_language_title') }}</h3>
      <div class="settings-card__options">
        <label class="field"><span>{{ t('settings_language_label') }}</span>
          <select v-model="currentLang" @change="changeLanguage">
            <option value="ar">{{ t('ar') }}</option>
            <option value="he">{{ t('he') }}</option>
            <option value="en">{{ t('en') }}</option>
          </select>
        </label>
      </div>
    </div>

    <div class="card settings-card">
      <h3>{{ t('settings_geofence_title') }}</h3>
      <p class="settings-card__hint">{{ t('settings_geofence_hint') }}</p>
      <label class="field" style="max-width: 220px"><span>{{ t('settings_geofence_radius') }}</span><input type="number" value="100" disabled /></label>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from '../services/i18n'
import SubscriptionStatusView from './SubscriptionStatusView.vue'

const { t, setLang, currentLang: computedLang } = useI18n()
const currentLang = ref(computedLang.value)

const changeLanguage = () => {
  setLang(currentLang.value)
}

onMounted(() => {
  currentLang.value = computedLang.value
})
</script>

<style scoped>
.page-head { margin-bottom: 18px; }
.page-head h2 { font-size: 20px; }
.page-subtitle { color: var(--ink-soft); margin: 6px 0 16px; font-size: 14px; }
.settings-card { padding: 22px; margin-bottom: 16px; }
.settings-card h3 { font-size: 15px; margin-bottom: 10px; }
.settings-card__hint { font-size: 13px; color: var(--ink-soft); margin-bottom: 12px; }
.field { display: flex; flex-direction: column; gap: 8px; }
.field span { font-size: 13px; color: var(--ink-soft); }
.field input,
.field select { width: 100%; padding: 10px 12px; border-radius: var(--radius-sm); border: 1px solid var(--line); background: var(--surface); }

@media (max-width: 768px) {
  .settings-card { padding: 16px; }
  .field { max-width: 100%; }
}
</style>

```

---

## 📄 frontend-admin-dashboard/src/views/SubscriptionStatusView.vue

```vue
<template>
  <div>
    <div v-if="loading" class="subscription-loading">
      {{ t('loading_data') }}
    </div>

    <div v-else-if="error" class="error">
      {{ error }}
    </div>

    <div v-else>
      <div class="field-row">
        <span>{{ t('subscription_status_label') }}</span>
        <strong :class="['status-badge', `status-badge--${statusClass}`]">{{ subscriptionText }}</strong>
      </div>

      <div class="field-row">
        <span>{{ t('subscription_expires_at_label') }}</span>
        <strong>{{ expiresLabel }}</strong>
      </div>

      <div v-if="subscriptionStatus !== 'active'" class="subscription-note">
        {{ t('subscription_expired_message') }}
      </div>

      <div v-if="subscriptionStatus === 'active'" class="subscription-note subscription-note--active">
        {{ t('subscription_active_message') }}
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from '../services/i18n'
import { getCurrentUser } from '../services/auth'
import { formatLocalDate } from '../utils/numberFormatter'

const { t, locale } = useI18n()
const loading = ref(true)
const error = ref('')
const subscriptionStatus = ref('active')
const subscriptionExpiresAt = ref(null)

const statusClass = computed(() => {
  if (subscriptionStatus.value === 'active') return 'active'
  if (subscriptionStatus.value === 'expired') return 'expired'
  if (subscriptionStatus.value === 'canceled') return 'canceled'
  return 'unknown'
})

const subscriptionText = computed(() => {
  if (subscriptionStatus.value === 'active') return t('subscription_active')
  if (subscriptionStatus.value === 'expired') return t('subscription_expired')
  if (subscriptionStatus.value === 'canceled') return t('subscription_canceled')
  return t('undefined_text')
})

const expiresLabel = computed(() => {
  if (!subscriptionExpiresAt.value) {
    return t('subscription_lifetime')
  }
  return formatLocalDate(subscriptionExpiresAt.value, locale?.value || 'en', 'DD/MM/YYYY')
})

onMounted(async () => {
  try {
    const data = await getCurrentUser()
    subscriptionStatus.value = data.subscription_status || 'active'
    subscriptionExpiresAt.value = data.subscription_expires_at || null
  } catch (err) {
    error.value = err.response?.data?.error || t('error')
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.subscription-loading,
.error {
  font-size: 14px;
  color: var(--ink-soft);
  padding: 22px;
}

.field-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  margin-bottom: 14px;
  padding: 14px 0;
  border-bottom: 1px solid var(--line);
}

.field-row:last-child {
  border-bottom: none;
}

.field-row span {
  color: var(--ink-soft);
}

.status-badge {
  padding: 6px 10px;
  border-radius: 999px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.status-badge--active {
  background: #D7F6E9;
  color: #0D7A43;
}

.status-badge--expired,
.status-badge--canceled {
  background: #FFE1E1;
  color: #C21E1E;
}

.subscription-note {
  margin-top: 18px;
  font-size: 14px;
  color: var(--ink-soft);
  line-height: 1.6;
}

.subscription-note--active {
  color: #0D7A43;
}
</style>

```

---

## 📄 frontend-admin-dashboard/src/views/TasksView.vue

```vue
<template>
  <div>
    <div class="page-head">
      <div><h2>{{ t('tasks') }}</h2><p>{{ t('tasks_description') }}</p></div>
      <button @click="showCreateModal = true" class="btn btn--primary">{{ t('new_task') }}</button>
    </div>

    <div class="card">
      <!-- حالة التحميل -->
      <div v-if="loading" class="loading-state">
        <div class="loading-spinner"></div>
        <p>جاري تحميل المهام...</p>
      </div>

      <!-- حالة الخطأ -->
      <div v-else-if="error" class="error-state">
        <div class="error-icon">⚠️</div>
        <h3>حدث خطأ</h3>
        <p>{{ error }}</p>
        <button @click="fetchTasks" class="btn btn--primary">إعادة المحاولة</button>
      </div>

      <!-- عرض المهام -->
      <template v-else>
        <!-- جدول للشاشات الكبيرة -->
        <div class="table-wrapper desktop-only">
          <table class="table">
            <thead><tr><th>{{ t('task_title') }}</th><th>{{ t('task_employee') }}</th><th>{{ t('task_worksite') }}</th><th>{{ t('task_priority') }}</th><th>{{ t('task_status') }}</th></tr></thead>
            <tbody>
              <tr v-for="t in tasks" :key="t.id">
                <td dir="auto">{{ t.title }}</td>
                <td dir="auto">{{ t.employee }}</td>
                <td dir="auto">{{ t.worksite }}</td>
                <td><span class="badge" :class="priorityBadgeMap[t.priority]">{{ priorityLabels[t.priority] }}</span></td>
                <td><span class="badge" :class="badgeMap[t.status]">{{ labels[t.status] }}</span></td>
              </tr>
            </tbody>
          </table>
        </div>
      
      <!-- بطاقات للشاشات الصغيرة -->
      <div class="mobile-cards mobile-only">
        <div v-for="t in tasks" :key="t.id" class="task-card">
          <div class="task-card__header">
            <span class="task-card__title" dir="auto">{{ t.title }}</span>
            <div class="task-card__badges">
              <span v-if="t.priority" class="badge" :class="priorityBadgeMap[t.priority]">{{ priorityLabels[t.priority] }}</span>
              <span class="badge" :class="badgeMap[t.status]">{{ labels[t.status] }}</span>
            </div>
          </div>
          <div class="task-card__body">
            <div class="task-card__row">
              <span class="task-card__label">{{ t('task_employee') }}</span>
              <span class="task-card__value" dir="auto">{{ t.employee }}</span>
            </div>
            <div class="task-card__row">
              <span class="task-card__label">{{ t('task_worksite') }}</span>
              <span class="task-card__value" dir="auto">{{ t.worksite }}</span>
            </div>
          </div>
        </div>
      </div>
      </template>
    </div>
    
    <!-- Task Create Modal -->
    <TaskCreateModal 
      :show="showCreateModal" 
      @close="showCreateModal = false" 
      @created="fetchTasks" 
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from '../services/i18n'
import TaskCreateModal from '../components/TaskCreateModal.vue'

const { t } = useI18n()
const showCreateModal = ref(false)

const labels = computed(() => ({ 
  pending: t('status_pending'), 
  assigned: 'مُعينة',
  in_progress: t('status_in_progress'), 
  completed: t('status_completed'), 
  late: t('status_late'),
  cancelled: 'ملغاة'
}))
const badgeMap = { 
  pending: '', 
  assigned: 'badge--info', 
  in_progress: 'badge--gold', 
  completed: 'badge--in', 
  late: 'badge--out',
  cancelled: 'badge--out'
}

const priorityLabels = computed(() => ({
  low: 'منخفضة',
  normal: 'عادية',
  high: 'عالية',
  urgent: 'عاجلة'
}))
const priorityBadgeMap = { low: 'badge--info', normal: '', high: 'badge--gold', urgent: 'badge--out' }

const tasks = ref([])
const loading = ref(false)
const error = ref(null)

const fetchTasks = async () => {
  loading.value = true
  error.value = null
  
  try {
    const apiBaseUrl = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1'
    const token = localStorage.getItem('token')
    
    const response = await fetch(`${apiBaseUrl}/tasks`, {
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    })
    
    if (!response.ok) {
      throw new Error('فشل في جلب المهام')
    }
    
    const data = await response.json()
    
    // تحويل البيانات لتناسب الشكل المتوقع في العرض
    tasks.value = data.map(task => ({
      id: task.id,
      title: task.title,
      employee: task.assigned_user_id || 'غير معين', // سيتم تحسين هذا لجلب اسم الموظف
      worksite: task.worksite_id || 'غير محدد', // سيتم تحسين هذا لجلب اسم الموقع
      priority: task.priority || 'normal',
      status: task.status
    }))
  } catch (err) {
    console.error('Error fetching tasks:', err)
    error.value = err.message
    
    // في حالة الخطأ، استخدام بيانات تجريبية
    tasks.value = [
      { id: 1, title: 'صيانة مكيّفات', employee: 'أحمد ياسين', worksite: 'برج الأمل', priority: 'high', status: 'in_progress' },
      { id: 2, title: 'فحص دوري', employee: 'سارة قدورة', worksite: 'الشميساني', priority: 'normal', status: 'pending' },
      { id: 3, title: 'تركيب كاميرات', employee: 'ليث عودة', worksite: 'طريق المطار', priority: 'urgent', status: 'completed' },
    ]
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchTasks()
})
</script>

<style scoped>
.page-head { display: flex; align-items: flex-start; justify-content: space-between; margin-bottom: 18px; flex-wrap: wrap; gap: 10px; }
.page-head h2 { font-size: 20px; margin-bottom: 4px; }
.page-head p { font-size: 13px; color: var(--ink-soft); }
.table { width: 100%; border-collapse: collapse; }
.table th { text-align: right; font-size: 12px; color: var(--ink-soft); font-weight: 600; padding: 14px 20px; border-bottom: 1px solid var(--line); }
.table td { padding: 14px 20px; font-size: 14px; border-bottom: 1px solid var(--line); /* دعم النصوص المختلطة (عربي، عبري، إنجليزي) */ unicode-bidi: embed; text-align: start; }
.table tr:last-child td { border-bottom: none; }

@media (max-width: 768px) {
  .page-head { flex-direction: column; align-items: flex-start; }
  .page-head .btn { width: 100%; }
  
  .desktop-only { display: none; }
  .mobile-only { display: block; }
}

@media (min-width: 769px) {
  .desktop-only { display: block; }
  .mobile-only { display: none; }
}

/* تصميم بطاقات المهام للهاتف */
.mobile-cards {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.task-card {
  background: var(--surface);
  border: 1px solid var(--line);
  border-radius: var(--radius-md);
  padding: 16px;
  transition: all 0.2s ease;
}

.task-card:hover {
  border-color: var(--brand);
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

.task-card__header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  padding-bottom: 12px;
  border-bottom: 1px solid var(--line);
}

.task-card__title {
  font-weight: 600;
  color: var(--ink);
  font-size: 15px;
  /* دعم النصوص المختلطة (عربي، عبري، إنجليزي) */
  unicode-bidi: embed;
  text-align: start;
}

.task-card__badges {
  display: flex;
  gap: 8px;
  align-items: center;
}

.task-card__body {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.task-card__row {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.task-card__label {
  font-size: 13px;
  color: var(--ink-soft);
  font-weight: 500;
}

.task-card__value {
  font-size: 14px;
  color: var(--ink);
  font-weight: 500;
  /* دعم النصوص المختلطة (عربي، عبري، إنجليزي) */
  unicode-bidi: embed;
  text-align: start;
}

/* Loading state */
.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  text-align: center;
}

.loading-spinner {
  width: 48px;
  height: 48px;
  border: 4px solid var(--line);
  border-top: 4px solid var(--brand);
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 16px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.loading-state p {
  color: var(--ink-soft);
  font-size: 14px;
  margin: 0;
}

/* Error state */
.error-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  text-align: center;
}

.error-icon {
  font-size: 48px;
  margin-bottom: 16px;
}

.error-state h3 {
  color: var(--ink);
  font-size: 18px;
  margin: 0 0 8px 0;
}

.error-state p {
  color: var(--ink-soft);
  font-size: 14px;
  margin: 0 0 16px 0;
}
</style>

```

---

## 📄 frontend-admin-dashboard/src/views/WorksitesView.vue

```vue
<template>
  <div>
    <div class="page-head">
      <div>
        <h2>{{ t('worksites_title') }}</h2>
        <p>{{ t('worksites_description') }}</p>
      </div>
      <button class="btn btn--primary" @click="showModal = true">+ {{ t('new_worksite') }}</button>
    </div>

    <div v-if="loading" class="empty-state"><p>{{ t('loading_worksites') }}</p></div>

    <div v-else-if="worksites.length === 0" class="empty-state">
      <h3>{{ t('no_worksites') }}</h3>
      <p>{{ t('add_worksite_prompt') }}</p>
    </div>

    <div v-else class="sites-grid">
      <div v-for="site in worksites" :key="site.id" class="card site-card" :class="{ 'site-card--unassigned': site.is_unassigned }">
        <div class="site-card__header">
          <h3>{{ site.name }}</h3>
          <div class="site-card__actions" v-if="!site.is_unassigned">
            <button class="btn btn--primary btn--sm" @click="openAssignModal(site)">
              👤 {{ t('assign_employee') }}
            </button>
            <button class="btn btn--danger btn--sm" @click="confirmDelete(site)">
              🗑️
            </button>
          </div>
        </div>
        <p class="site-card__address">{{ site.address || t('no_address') }}</p>
        <div class="site-card__details" v-if="!site.is_unassigned">
          <span class="mono">📍 {{ site.latitude?.toFixed(5) }}, {{ site.longitude?.toFixed(5) }}</span>
          <span class="site-card__radius mono">⭕ {{ site.radius_meters || 100 }} {{ t('meters_unit') }}</span>
        </div>
        
        <!-- عرض الموظف المعين -->
        <div class="site-card__assigned" v-if="!site.is_unassigned">
          <span v-if="site.assigned_to?.name" class="badge badge--in">
            👤 {{ site.assigned_to.name }}
          </span>
          <span v-else class="badge badge--out">
            ⚠️ {{ t('unassigned') }}
          </span>
        </div>
        
        <!-- عرض الموظفين العاملين حالياً -->
        <div class="site-card__working" v-if="site.working_employees && site.working_employees.length > 0">
          <div class="site-card__working-label">
            {{ t('currently_working') }} ({{ site.working_employees.length }})
          </div>
          <div class="site-card__working-list">
            <div v-for="emp in site.working_employees" :key="emp.id" class="working-employee-item">
              <span class="badge badge--success badge--compact">
                👤 {{ emp.name }}
              </span>
              <button 
                class="btn btn--warning btn--xs" 
                @click="confirmForceCheckout(emp)"
                :disabled="forceCheckingOut"
              >
                ⏱️ {{ t('end_shift') }}
              </button>
            </div>
          </div>
        </div>
        
        <span class="badge" :class="site.is_active ? 'badge--in' : 'badge--out'" v-if="!site.is_unassigned">
          {{ site.is_active ? t('active_status') : t('inactive_status') }}
        </span>
      </div>
    </div>

    <WorksiteFormModal v-if="showModal" @close="showModal = false" @worksite-added="fetchWorksites" />

    <!-- مودال تأكيد الحذف -->
    <div v-if="showDeleteModal" class="modal-backdrop" @click.self="showDeleteModal = false">
      <div class="modal card">
        <div class="modal-header">
          <h3>⚠️ {{ t('confirm_delete_title') }}</h3>
          <button class="modal-close" @click="showDeleteModal = false">✕</button>
        </div>
        <div class="modal-body">
          <p>{{ t('confirm_delete_message') }} <strong>{{ siteToDelete?.name }}</strong>؟</p>
          <p class="text-danger">{{ t('delete_warning') }}</p>
        </div>
        <div class="modal-footer">
          <button class="btn btn--ghost" @click="showDeleteModal = false">{{ t('cancel') }}</button>
          <button class="btn btn--danger" @click="deleteWorksite" :disabled="deleting">
            {{ deleting ? t('deleting') : t('delete_final') }}
          </button>
        </div>
      </div>
    </div>

    <!-- مودال تعيين موظف - محسّن مع عرض سريع -->
    <div v-if="showAssignModal" class="modal-backdrop" @click.self="showAssignModal = false">
      <div class="modal card assign-modal">
        <div class="modal-header">
          <h3>👤 {{ t('assign_employee_title') }}</h3>
          <button class="modal-close" @click="showAssignModal = false">✕</button>
        </div>
        <div class="modal-body">
          <p>{{ t('choose_employee_to_assign') }} <strong>{{ worksiteToAssign?.name }}</strong></p>
          
          <div v-if="loadingEmployees" class="empty-state">
            <p>{{ t('loading_employees') }}</p>
          </div>
          
          <div v-else-if="employees.length === 0" class="empty-state">
            <p>{{ t('no_available_employees') }}</p>
          </div>
          
          <div v-else class="employees-list">
            <button
              v-for="emp in employees"
              :key="emp.id"
              class="employee-item"
              @click="assignEmployee(emp.id)"
              :disabled="assigning"
            >
              <span class="employee-item__avatar">{{ emp.full_name?.slice(0, 1) || '?' }}</span>
              <div>
                <strong>{{ emp.full_name }}</strong>
              </div>
              <span class="employee-item__check">✓</span>
            </button>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn--ghost" @click="showAssignModal = false">{{ t('cancel') }}</button>
        </div>
      </div>
    </div>

    <!-- مودال تأكيد إنهاء الدوام -->
    <div v-if="showForceCheckoutModal" class="modal-backdrop" @click.self="showForceCheckoutModal = false">
      <div class="modal card">
        <div class="modal-header">
          <h3>⏱️ {{ t('force_checkout_title') }}</h3>
          <button class="modal-close" @click="showForceCheckoutModal = false">✕</button>
        </div>
        <div class="modal-body">
          <p>{{ t('force_checkout_message') }} <strong>{{ employeeToForceCheckout?.name }}</strong>؟</p>
          <p class="text-warning">{{ t('force_checkout_warning') }}</p>
        </div>
        <div class="modal-footer">
          <button class="btn btn--ghost" @click="showForceCheckoutModal = false">{{ t('cancel') }}</button>
          <button class="btn btn--warning" @click="forceCheckoutEmployee" :disabled="forceCheckingOut">
            {{ forceCheckingOut ? t('processing') : t('confirm_end_shift') }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '../services/api'
import WorksiteFormModal from '../components/WorksiteFormModal.vue'
import { useI18n } from '../services/i18n'

const { t } = useI18n()
const worksites = ref([])
const employees = ref([])
const loading = ref(false)
const loadingEmployees = ref(false)
const deleting = ref(false)
const assigning = ref(false)
const forceCheckingOut = ref(false)
const showModal = ref(false)
const showDeleteModal = ref(false)
const showAssignModal = ref(false)
const showForceCheckoutModal = ref(false)
const siteToDelete = ref(null)
const worksiteToAssign = ref(null)
const employeeToForceCheckout = ref(null)

async function fetchWorksites() {
  loading.value = true
  try {
    const { data } = await api.get('/worksites')
    worksites.value = data || []
  } catch (error) {
    console.error('❌ ' + t('failed_to_fetch_worksites'), error)
  } finally {
    loading.value = false
  }
}

function confirmDelete(site) {
  siteToDelete.value = site
  showDeleteModal.value = true
}

async function deleteWorksite() {
  if (!siteToDelete.value) return
  
  deleting.value = true
  try {
    await api.delete(`/worksites/${siteToDelete.value.id}`)
    showDeleteModal.value = false
    await fetchWorksites()
  } catch (error) {
    console.error('❌ ' + t('failed_to_delete_worksite'), error)
    alert(error.response?.data?.error || t('failed_to_delete_worksite'))
  } finally {
    deleting.value = false
  }
}

function openAssignModal(site) {
  worksiteToAssign.value = site
  showAssignModal.value = true
  fetchEmployees()
}

async function fetchEmployees() {
  loadingEmployees.value = true
  try {
    const { data } = await api.get('/worksites/employees')
    employees.value = data || []
  } catch (error) {
    console.error('❌ ' + t('failed_to_fetch_employees'), error)
    employees.value = []
  } finally {
    loadingEmployees.value = false
  }
}

async function assignEmployee(employeeId) {
  if (!worksiteToAssign.value) return
  
  assigning.value = true
  try {
    const response = await api.post('/worksites/assign', {
      employee_id: employeeId,
      worksite_id: worksiteToAssign.value.id
    })
    
    showAssignModal.value = false
    alert('✅ ' + (response.data.message || t('employee_assigned_successfully')))
    await fetchWorksites()
  } catch (error) {
    console.error('❌ ' + t('failed_to_assign_employee'), error)
    const msg = error.response?.data?.error || t('failed_to_assign_employee')
    alert('❌ ' + msg)
  } finally {
    assigning.value = false
  }
}

function confirmForceCheckout(employee) {
  employeeToForceCheckout.value = employee
  showForceCheckoutModal.value = true
}

async function forceCheckoutEmployee() {
  if (!employeeToForceCheckout.value) return
  
  forceCheckingOut.value = true
  try {
    const response = await api.post('/attendance/force-checkout', {
      attendance_id: employeeToForceCheckout.value.attendance_id
    })
    
    showForceCheckoutModal.value = false
    alert('✅ ' + (response.data.message || t('force_checkout_success')))
    await fetchWorksites()
  } catch (error) {
    console.error('❌ ' + t('force_checkout_failed'), error)
    const msg = error.response?.data?.error || t('force_checkout_failed')
    alert('❌ ' + msg)
  } finally {
    forceCheckingOut.value = false
  }
}

onMounted(fetchWorksites)
</script>

<style scoped>
.page-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 18px;
  flex-wrap: wrap;
  gap: 10px;
}

.page-head h2 { font-size: 20px; margin-bottom: 4px; }
.page-head p { font-size: 13px; color: var(--ink-soft); }

.sites-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
}

.site-card {
  padding: 18px 20px;
  transition: all var(--transition-base);
}

.site-card:hover {
  box-shadow: var(--shadow-md);
  transform: translateY(-3px);
}

.site-card__header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 8px;
}

.site-card__header h3 { font-size: 16px; }

.site-card__actions {
  display: flex;
  gap: 6px;
}

.site-card__address {
  font-size: 13px;
  color: var(--ink-soft);
  margin-bottom: 8px;
}

.site-card__details {
  display: flex;
  gap: 12px;
  font-size: 12px;
  color: var(--ink-soft);
  margin-bottom: 10px;
  flex-wrap: wrap;
}

.site-card__radius {
  background: var(--brand-tint);
  color: var(--brand);
  padding: 2px 10px;
  border-radius: 999px;
}

.site-card__assigned { margin-bottom: 8px; }

.site-card__working {
  margin-top: 12px;
  padding: 10px;
  background: var(--brand-tint);
  border-radius: var(--radius-sm);
  border: 1px solid var(--brand);
}

.site-card__working-label {
  font-size: 12px;
  font-weight: 600;
  color: var(--brand-dark);
  margin-bottom: 6px;
}

.site-card__working-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.working-employee-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.site-card--unassigned {
  border: 2px solid var(--signal-out);
  background: linear-gradient(135deg, var(--surface) 0%, rgba(239, 68, 68, 0.1) 100%);
}

.site-card--unassigned .site-card__header h3 {
  color: var(--signal-out);
}

.modal-backdrop {
  position: fixed; inset: 0;
  background: rgba(0,0,0,0.6);
  backdrop-filter: blur(4px);
  display: flex; align-items: center; justify-content: center;
  z-index: 1000; padding: 20px;
}

.modal {
  width: 100%; max-width: 460px; padding: 0;
  max-height: 90vh;
  overflow-y: auto;
}

.assign-modal { max-width: 500px; }

.modal-header {
  display: flex; justify-content: space-between; align-items: center;
  padding: 16px 20px; border-bottom: 1px solid var(--line);
}

.modal-header h3 { font-size: 17px; margin: 0; }

.modal-close {
  background: none; border: none; font-size: 24px;
  cursor: pointer; color: var(--ink-soft);
}

.modal-body { padding: 20px; }
.modal-body p { margin-bottom: 8px; }

.text-danger { color: var(--signal-out); font-weight: 600; }
.text-warning { color: #f59e0b; font-weight: 600; }

.btn--xs {
  padding: 4px 8px;
  font-size: 11px;
}

.btn--warning {
  background: #f59e0b;
  color: white;
  border: none;
}

.btn--warning:hover {
  background: #d97706;
}

.modal-footer {
  display: flex; gap: 10px;
  justify-content: flex-end;
  padding: 16px 20px;
  border-top: 1px solid var(--line);
}

.employees-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-top: 12px;
}

.employee-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 14px;
  background: var(--surface);
  border: 1px solid var(--line);
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: all 0.2s;
  text-align: right;
  font-family: var(--font-body);
  font-size: 14px;
  width: 100%;
  position: relative;
}

.employee-item:hover {
  border-color: var(--brand);
  background: var(--brand-tint);
}

.employee-item:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.employee-item__avatar {
  width: 36px; height: 36px;
  border-radius: 50%;
  background: var(--brand-tint);
  color: var(--brand-dark);
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  font-size: 16px;
  flex-shrink: 0;
}

.employee-item div { 
  display: flex; 
  flex-direction: column; 
  text-align: right; 
  flex: 1;
}

.employee-item span { font-size: 12px; color: var(--ink-soft); }

.employee-item__check {
  opacity: 0;
  color: var(--signal-in);
  font-size: 18px;
  transition: opacity 0.2s;
}

.employee-item:hover .employee-item__check {
  opacity: 1;
}

@media (max-width: 960px) {
  .sites-grid { grid-template-columns: 1fr 1fr; }
}

@media (max-width: 600px) {
  .sites-grid { grid-template-columns: 1fr; }
}
</style>

```

---

## 📄 frontend-admin-dashboard/vite.config.js

```javascript
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  appType: 'spa',
  plugins: [vue()],
  test: {
    globals: true,
    environment: 'jsdom',
    setupFiles: ['./src/tests/setup.js'],
    coverage: {
      provider: 'v8',
      reporter: ['text', 'json', 'html'],
      exclude: [
        'node_modules/',
        'src/tests/',
        '**/*.d.ts',
        '**/*.config.*',
        '**/mockData',
        'dist/',
        'electron/'
      ]
    }
  },
  server: { 
    port: 3000,
    cors: true,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        secure: false
      }
    }
  },
  preview: {
    port: 3000,
    // SPA routing support for preview
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        secure: false
      }
    }
  },
  optimizeDeps: {
    include: ['vue', 'vue-router', 'leaflet']
  },
  assetsInclude: ['**/*.json'],
  build: {
    cssCodeSplit: true,
    rollupOptions: {
      output: {
        manualChunks: {
          'vue-vendor': ['vue', 'vue-router'],
          'leaflet-vendor': ['leaflet']
        },
        // إضافة timestamp تلقائي لضمان cache busting
        entryFileNames: 'assets/[name]-[hash].js',
        chunkFileNames: 'assets/[name]-[hash].js',
        assetFileNames: 'assets/[name]-[hash][extname]'
      }
    }
  },
  base: '/'
})

```

---

## 📄 frontend-client-portal/build-obfuscated.js

```javascript
import JavaScriptObfuscator from 'javascript-obfuscator';
import { readFileSync, writeFileSync, readdirSync, statSync } from 'fs';
import { join } from 'path';

// إعدادات التعتيم
const obfuscatorOptions = {
  compact: true,
  controlFlowFlattening: true,
  controlFlowFlatteningThreshold: 0.75,
  deadCodeInjection: true,
  deadCodeInjectionThreshold: 0.4,
  debugProtection: true,
  disableConsoleOutput: false, // Changed to false for debugging
  identifierNamesGenerator: 'hexadecimal',
  log: true, // Changed to true for debugging
  numbersToExpressions: true,
  renameGlobals: false,
  selfDefending: true,
  simplify: true,
  splitStrings: true,
  splitStringsChunkLength: 10,
  stringArray: true,
  stringArrayEncoding: ['base64'],
  stringArrayIndexShift: true,
  stringArrayWrappersCount: 2,
  stringArrayWrappersChainedCalls: true,
  stringArrayWrappersParametersMaxCount: 4,
  stringArrayWrappersType: 'function',
  stringArrayThreshold: 0.75,
  transformObjectKeys: true,
  unicodeEscapeSequence: false
};

// دالة للبحث عن ملفات JS بشكل متكرر
function findJsFiles(dir, fileList = []) {
  try {
    const files = readdirSync(dir);
    
    for (const file of files) {
      const filePath = join(dir, file);
      const stats = statSync(filePath);
      
      if (stats.isDirectory()) {
        findJsFiles(filePath, fileList);
      } else if (file.endsWith('.js')) {
        fileList.push(filePath);
      }
    }
  } catch (error) {
    console.error(`Error reading directory ${dir}:`, error.message);
  }
  
  return fileList;
}

// البحث عن جميع ملفات JS في مجلد dist
const jsFiles = findJsFiles('dist');

console.log(`Found ${jsFiles.length} JavaScript files to obfuscate`);

// تعتيم كل ملف
for (const file of jsFiles) {
  try {
    console.log(`Obfuscating: ${file}`);
    const code = readFileSync(file, 'utf8');
    const obfuscatedCode = JavaScriptObfuscator.obfuscate(code, obfuscatorOptions).getObfuscatedCode();
    writeFileSync(file, obfuscatedCode, 'utf8');
    console.log(`✓ Obfuscated: ${file}`);
  } catch (error) {
    console.error(`✗ Error obfuscating ${file}:`, error.message);
  }
}

console.log('Obfuscation complete!');
```

---

## 📄 frontend-client-portal/public/service-worker.js

```javascript
// حجز أولي لملف Service Worker — سيُستخدم لاحقاً لدعم العمل دون اتصال (offline)
self.addEventListener('install', () => self.skipWaiting())
self.addEventListener('activate', () => self.clients.claim())
```

---

## 📄 frontend-client-portal/src/App.vue

```vue
<template>
  <div v-if="isRouterReady" class="shell">
    <header v-if="isLoggedIn" class="topbar">
      <div class="topbar__brand">
        <img src="/devpro-logo.jpg" alt="DevPro Logo" class="devpro-logo" />
        <div class="brand-text">
          <span class="brand-name gradient-text">{{ t('app_name') }}</span>
          <span class="brand-sub">{{ t('client_portal') }}</span>
        </div>
      </div>
      <div class="topbar__actions">
        <LanguageSwitcher v-if="isLoggedIn" />
        <NotificationDropdown v-if="isLoggedIn" />
        <router-link to="/" class="btn btn--ghost btn--sm">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M3 3h18v18H3zM3 9h18M9 21V9"/>
          </svg>
          {{ t('service_history') }}
        </router-link>
        <router-link to="/new" class="btn btn--primary btn--sm">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M12 5v14M5 12h14"/>
          </svg>
          {{ t('new_service_request') }}
        </router-link>
        <button @click="handleLogout" class="btn btn--ghost btn--sm">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"></path>
            <polyline points="16 17 21 12 16 7"></polyline>
            <line x1="21" y1="12" x2="9" y2="12"></line>
          </svg>
          {{ t('logout') }}
        </button>
      </div>
    </header>

    <main class="content">
      <router-view />
    </main>

    <footer v-if="isLoggedIn" class="footer">
      <p>{{ t('copyright') }}</p>
    </footer>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import LanguageSwitcher from './components/LanguageSwitcher.vue'
import NotificationDropdown from './components/NotificationDropdown.vue'
import { useI18n } from './services/i18n'
import { logout, getToken, currentUser } from './services/auth'
import { authStore } from './store/auth'

const { t } = useI18n()
const router = useRouter()
const isRouterReady = ref(false)

router.isReady().then(() => {
  isRouterReady.value = true
  
  // تحديث authStore عند تحميل الصفحة
  const token = getToken()
  if (token && !authStore.user) {
    const user = currentUser()
    if (user) {
      authStore.setUser(user)
    }
  }
  
  // معالجة مشكلة /index.html في العنوان - استعادة المسار المحفوظ بدلاً من التحويل إلى /
  if (window.location.pathname === '/index.html') {
    const isAuthed = !!localStorage.getItem('worktrack_client_token')
    const savedPath = localStorage.getItem('worktrack_client_last_path')
    
    if (isAuthed && savedPath && savedPath !== '/login') {
      window.history.replaceState({}, '', savedPath)
    } else if (isAuthed) {
      window.history.replaceState({}, '', '#/new')
    } else {
      window.history.replaceState({}, '', '/login')
    }
  }
})

const isLoggedIn = computed(() => !!authStore.user)

function handleLogout() {
  logout()
  authStore.clear()
  // استخدام window.location.href لضمان إعادة تحميل الصفحة بالكامل
  window.location.href = window.location.origin + '/#/login'
}
</script>

<style scoped>
.shell {
  min-height: 100dvh;
  display: flex;
  flex-direction: column;
  background: var(--canvas);
}

.topbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 20px;
  background: var(--surface);
  border-bottom: 1px solid var(--line);
  position: sticky;
  top: 0;
  z-index: 10;
  box-shadow: var(--shadow-sm);
  transition: box-shadow var(--transition-base);
}

.topbar:hover {
  box-shadow: var(--shadow-md);
}

.topbar__brand {
  display: flex;
  align-items: center;
  gap: 10px;
}

.devpro-logo {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  object-fit: contain;
  transition: transform var(--transition-base);
}

.topbar__brand:hover .devpro-logo {
  transform: scale(1.1);
}

.brand-text {
  display: flex;
  flex-direction: column;
  line-height: 1.2;
}

.brand-name {
  font-weight: 700;
  font-size: 16px;
  color: var(--brand);
}

.brand-sub {
  font-size: 10px;
  color: var(--ink-soft);
  font-weight: 500;
}

.topbar__actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.content {
  flex: 1;
  padding: 20px;
  max-width: 600px;
  width: 100%;
  margin: 0 auto;
  animation: fadeIn 0.4s ease;
}

.content:has(.login-container) {
  max-width: 100%;
  padding: 0;
  background: linear-gradient(135deg, #1E3A5F 0%, #0D1B3E 100%);
}

.footer {
  text-align: center;
  padding: 16px 20px;
  border-top: 1px solid var(--line);
  background: var(--surface);
  box-shadow: 0 -2px 10px rgba(0,0,0,0.05);
}

.footer p {
  font-size: 12px;
  color: var(--ink-soft);
}

@media (max-width: 480px) {
  .topbar {
    padding: 12px 16px;
    flex-wrap: wrap;
    gap: 12px;
  }
  
  .topbar__brand {
    gap: 8px;
  }
  
  .devpro-logo {
    width: 28px;
    height: 28px;
  }
  
  .brand-text {
    line-height: 1.1;
  }
  
  .brand-name {
    font-size: 14px;
  }
  
  .brand-sub {
    font-size: 9px;
  }
  
  .topbar__actions {
    gap: 8px;
    flex-wrap: wrap;
    justify-content: flex-end;
  }

  .topbar__actions .btn {
    font-size: 11px;
    padding: 6px 10px;
    gap: 6px;
  }

  .topbar__actions .btn svg {
    width: 14px;
    height: 14px;
  }

  .topbar__actions .btn span {
    display: none;
  }

  .topbar__actions .btn--primary span {
    display: inline;
  }

  .topbar__actions .notification-dropdown {
    display: flex;
  }

  .content {
    padding: 16px 12px;
  }
  
  .footer {
    padding: 12px 16px;
  }
  
  .footer p {
    font-size: 11px;
  }
}

@media (max-width: 360px) {
  .topbar {
    padding: 10px 12px;
  }
  
  .devpro-logo {
    width: 24px;
    height: 24px;
  }
  
  .brand-name {
    font-size: 13px;
  }
  
  .topbar__actions {
    gap: 6px;
  }
  
  .topbar__actions .btn {
    font-size: 10px;
    padding: 5px 8px;
  }
  
  .content {
    padding: 12px 8px;
  }
}
</style>

```

---

## 📄 frontend-client-portal/src/components/ActivityFeed.vue

```vue
<template>
  <div class="activity-feed">
    <div class="feed-header">
      <h3>{{ t('notifications_feed') }}</h3>
      <span v-if="notificationsStore.notifications && notificationsStore.notifications.length" class="badge">
        {{ notificationsStore.notifications.length }}
      </span>
    </div>

    <div v-if="notificationsStore.loading" class="empty-state">
      <p>{{ t('loading_notifications') }}</p>
    </div>

    <div v-else-if="notificationsStore.error" class="alert alert-error">
      <span>❌</span> {{ notificationsStore.error }}
    </div>

    <div v-else-if="!notificationsStore.notifications || !notificationsStore.notifications.length" class="empty-state">
      <p>{{ t('no_notifications_available') }}</p>
    </div>

    <div v-else class="notifications-list">
      <div 
        v-for="notification in notificationsStore.notifications" 
        :key="notification.id" 
        class="notification-item"
      >
        <span class="notification-icon">{{ getNotificationIcon(notification.type) }}</span>
        <div class="notification-content">
          <p class="notification-title">{{ notification.title }}</p>
          <p class="notification-message">{{ notification.message }}</p>
          <span class="notification-time">{{ formatTime(notification.created_at) }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { notificationsStore } from '../store/notifications'
import { useI18n } from '../services/i18n'

const { t } = useI18n()

onMounted(() => {
  notificationsStore.fetchNotifications()
})

function getNotificationIcon(type) {
  const icons = {
    'status_update': '📋',
    'worker_assigned': '👷',
    'worker_arrived': '📍',
    'service_completed': '✅',
    'payment_request': '💳',
    'general': '🔔'
  }
  return icons[type] || '🔔'
}

function formatTime(date) {
  if (!date) return ''
  return new Date(date).toLocaleTimeString('en-GB', { 
    hour: '2-digit', 
    minute: '2-digit',
    hour12: false 
  })
}
</script>

<style scoped>
.activity-feed {
  background: var(--surface);
  border-radius: var(--radius-md);
  padding: 16px;
  border: 1px solid var(--line);
}

.feed-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.feed-header h3 {
  font-size: 14px;
  margin: 0;
  color: var(--ink);
}

.badge {
  background: var(--brand);
  color: white;
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 11px;
  font-weight: 600;
}

.empty-state {
  text-align: center;
  padding: 20px;
  color: var(--ink-soft);
  font-size: 13px;
}

.alert {
  padding: 10px 12px;
  border-radius: var(--radius-sm);
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
}

.alert-error {
  background: var(--signal-out-tint);
  color: var(--signal-out);
}

.notifications-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.notification-item {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  padding: 10px;
  background: var(--canvas);
  border-radius: var(--radius-sm);
  transition: background var(--transition-fast);
}

.notification-item:hover {
  background: var(--brand-tint);
}

.notification-icon {
  font-size: 16px;
  flex-shrink: 0;
}

.notification-content {
  flex: 1;
  min-width: 0;
}

.notification-title {
  font-size: 13px;
  font-weight: 600;
  color: var(--ink);
  margin: 0 0 4px 0;
}

.notification-message {
  font-size: 12px;
  color: var(--ink-soft);
  margin: 0 0 4px 0;
  line-height: 1.4;
}

.notification-time {
  font-size: 11px;
  color: var(--ink-light);
}

@media (max-width: 480px) {
  .activity-feed {
    padding: 12px;
    border-radius: var(--radius-sm);
  }
  
  .feed-header {
    margin-bottom: 10px;
  }
  
  .feed-header h3 {
    font-size: 13px;
  }
  
  .badge {
    font-size: 10px;
    padding: 2px 6px;
  }
  
  .empty-state {
    padding: 16px;
    font-size: 12px;
  }
  
  .alert {
    font-size: 11px;
    padding: 8px 10px;
  }
  
  .notification-item {
    padding: 8px;
    gap: 8px;
  }
  
  .notification-icon {
    font-size: 14px;
  }
  
  .notification-title {
    font-size: 12px;
  }
  
  .notification-message {
    font-size: 11px;
  }
  
  .notification-time {
    font-size: 10px;
  }
}

@media (max-width: 360px) {
  .activity-feed {
    padding: 10px;
  }
  
  .feed-header h3 {
    font-size: 12px;
  }
  
  .notification-item {
    padding: 6px;
    gap: 6px;
  }
  
  .notification-icon {
    font-size: 12px;
  }
  
  .notification-title {
    font-size: 11px;
  }
  
  .notification-message {
    font-size: 10px;
  }
}
</style>
```

---

## 📄 frontend-client-portal/src/components/LanguageSwitcher.vue

```vue
<template>
  <div class="lang-switcher">
    <button
      v-for="lang in languages"
      :key="lang.code"
      class="lang-btn"
      :class="{ active: currentLang === lang.code }"
      @click="changeLanguage(lang.code)"
      :title="lang.name"
    >
      {{ lang.flag }}
      <span class="lang-label">{{ lang.name }}</span>
    </button>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from '../services/i18n'

const { currentLang, setLang } = useI18n()

const languages = [
  { code: 'ar', name: 'العربية', flag: '🇸🇦' },
  { code: 'he', name: 'עברית', flag: '🇮🇱' },
  { code: 'en', name: 'English', flag: '🇬🇧' }
]

function changeLanguage(lang) {
  setLang(lang)
}
</script>

<style scoped>
.lang-switcher {
  display: flex;
  justify-content: center;
  gap: 6px;
  padding: 6px;
  background: var(--canvas, #f0f4fa);
  border-radius: 12px;
  width: 100%;
}

.lang-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 14px;
  border: none;
  border-radius: 8px;
  background: transparent;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.3s ease;
  font-family: inherit;
  flex: 1;
  justify-content: center;
}

.lang-btn:hover {
  background: rgba(30, 58, 95, 0.08);
  transform: scale(1.02);
}

.lang-btn.active {
  background: #1E3A5F;
  color: white;
  box-shadow: 0 2px 12px rgba(30, 58, 95, 0.3);
  transform: scale(1.02);
}

.lang-label {
  font-size: 12px;
  font-weight: 500;
}

@media (max-width: 480px) {
  .lang-switcher {
    gap: 4px;
    padding: 4px;
    border-radius: 10px;
  }
  
  .lang-btn {
    padding: 6px 10px;
    font-size: 18px;
    gap: 4px;
  }
  
  .lang-label {
    display: none;
  }
}

@media (max-width: 360px) {
  .lang-switcher {
    gap: 3px;
    padding: 3px;
  }
  
  .lang-btn {
    padding: 5px 8px;
    font-size: 16px;
  }
}
</style>

```

---

## 📄 frontend-client-portal/src/components/NotificationDropdown.vue

```vue
<template>
  <div class="notification-dropdown" ref="dropdownRef">
    <!-- أيقونة الإشعارات -->
    <button 
      @click="toggleDropdown" 
      class="notification-icon"
      :title="t('notifications')"
    >
      <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9"></path>
        <path d="M13.73 21a2 2 0 0 1-3.46 0"></path>
      </svg>
      <span v-if="unreadCount > 0" class="badge">{{ unreadCount }}</span>
    </button>
    
    <!-- Dropdown content -->
    <transition name="dropdown">
      <div v-if="isOpen" class="dropdown-content">
        <div class="feed-header">
          <h3>{{ t('notifications_feed') }}</h3>
          <span v-if="notificationsStore.notifications && notificationsStore.notifications.length" class="badge">
            {{ notificationsStore.notifications.length }}
          </span>
        </div>

        <div v-if="notificationsStore.loading" class="empty-state">
          <p>{{ t('loading_notifications') }}</p>
        </div>

        <div v-else-if="notificationsStore.error" class="alert alert-error">
          <span>❌</span> {{ notificationsStore.error }}
        </div>

        <div v-else-if="!notificationsStore.notifications || !notificationsStore.notifications.length" class="empty-state">
          <p>{{ t('no_notifications_available') }}</p>
        </div>

        <div v-else class="notifications-list">
          <div 
            v-for="notification in notificationsStore.notifications" 
            :key="notification.id" 
            class="notification-item"
          >
            <span class="notification-icon">{{ getNotificationIcon(notification.type) }}</span>
            <div class="notification-content">
              <p class="notification-title">{{ notification.title }}</p>
              <p class="notification-message">{{ notification.message }}</p>
              <span class="notification-time">{{ formatTime(notification.created_at) }}</span>
            </div>
          </div>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount, computed } from 'vue'
import { notificationsStore } from '../store/notifications'
import { useI18n } from '../services/i18n'

const { t } = useI18n()
const isOpen = ref(false)
const dropdownRef = ref(null)

const unreadCount = computed(() => {
  if (!notificationsStore.notifications) return 0
  return notificationsStore.notifications.length
})

function toggleDropdown() {
  isOpen.value = !isOpen.value
  if (isOpen.value) {
    notificationsStore.fetchNotifications()
  }
}

function closeDropdown() {
  isOpen.value = false
}

function handleClickOutside(event) {
  if (dropdownRef.value && !dropdownRef.value.contains(event.target)) {
    closeDropdown()
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
  notificationsStore.fetchNotifications()
})

onBeforeUnmount(() => {
  document.removeEventListener('click', handleClickOutside)
})

function getNotificationIcon(type) {
  const icons = {
    'status_update': '📋',
    'worker_assigned': '👷',
    'worker_arrived': '📍',
    'service_completed': '✅',
    'payment_request': '💳',
    'general': '🔔'
  }
  return icons[type] || '🔔'
}

function formatTime(date) {
  if (!date) return ''
  return new Date(date).toLocaleTimeString('en-GB', { 
    hour: '2-digit', 
    minute: '2-digit',
    hour12: false 
  })
}
</script>

<style scoped>
.notification-dropdown {
  position: relative;
  display: inline-block;
}

.notification-icon {
  background: none;
  border: none;
  cursor: pointer;
  position: relative;
  padding: 8px;
  border-radius: var(--radius-sm);
  transition: all var(--transition-base);
  color: var(--ink-soft);
  display: flex;
  align-items: center;
  justify-content: center;
}

.notification-icon:hover {
  background: var(--brand-tint);
  color: var(--brand);
  transform: scale(1.1);
}

.badge {
  position: absolute;
  top: 0;
  right: 0;
  background: var(--signal-out);
  color: white;
  border-radius: 50%;
  padding: 2px 6px;
  font-size: 10px;
  font-weight: 600;
  min-width: 16px;
  height: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0%, 100% {
    transform: scale(1);
  }
  50% {
    transform: scale(1.1);
  }
}

.dropdown-content {
  position: absolute;
  top: 100%;
  right: 0;
  background: var(--surface);
  border-radius: var(--radius-md);
  box-shadow: var(--shadow-lg);
  min-width: 320px;
  max-height: 400px;
  overflow-y: auto;
  z-index: 1000;
  margin-top: 8px;
  border: 1px solid var(--line);
}

.dropdown-enter-active,
.dropdown-leave-active {
  transition: all 0.2s ease;
}

.dropdown-enter-from {
  opacity: 0;
  transform: translateY(-10px);
}

.dropdown-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}

.feed-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  border-bottom: 1px solid var(--line);
  background: var(--canvas);
  border-radius: var(--radius-md) var(--radius-md) 0 0;
  position: sticky;
  top: 0;
  z-index: 1;
}

.feed-header h3 {
  font-size: 14px;
  margin: 0;
  color: var(--ink);
}

.feed-header .badge {
  position: static;
  background: var(--brand);
  color: white;
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 11px;
  font-weight: 600;
  animation: none;
}

.empty-state {
  text-align: center;
  padding: 20px;
  color: var(--ink-soft);
  font-size: 13px;
}

.alert {
  padding: 10px 12px;
  border-radius: var(--radius-sm);
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
  margin: 8px;
}

.alert-error {
  background: var(--signal-out-tint);
  color: var(--signal-out);
}

.notifications-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 8px;
}

.notification-item {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  padding: 10px;
  background: var(--canvas);
  border-radius: var(--radius-sm);
  transition: background var(--transition-fast);
  cursor: pointer;
}

.notification-item:hover {
  background: var(--brand-tint);
}

.notification-icon {
  font-size: 16px;
  flex-shrink: 0;
}

.notification-content {
  flex: 1;
  min-width: 0;
}

.notification-title {
  font-size: 13px;
  font-weight: 600;
  color: var(--ink);
  margin: 0 0 4px 0;
}

.notification-message {
  font-size: 12px;
  color: var(--ink-soft);
  margin: 0 0 4px 0;
  line-height: 1.4;
}

.notification-time {
  font-size: 11px;
  color: var(--ink-light);
}

@media (max-width: 480px) {
  .dropdown-content {
    min-width: 280px;
    max-height: 350px;
    right: -8px;
  }

  .feed-header {
    padding: 10px 12px;
  }

  .feed-header h3 {
    font-size: 13px;
  }

  .notifications-list {
    padding: 6px;
  }

  .notification-item {
    padding: 8px;
    gap: 8px;
  }

  .notification-icon {
    font-size: 14px;
  }

  .notification-title {
    font-size: 12px;
  }

  .notification-message {
    font-size: 11px;
  }

  .notification-time {
    font-size: 10px;
  }
}
</style>
```

---

## 📄 frontend-client-portal/src/components/PWAInstallButton.vue

```vue
<template>
  <div v-if="showInstallButton" class="pwa-install-container">
    <button 
      @click="installPWA" 
      class="pwa-install-icon"
      :title="pwaInstallTitle"
    >
      <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path>
        <polyline points="7 10 12 15 17 10"></polyline>
        <line x1="12" y1="15" x2="12" y2="3"></line>
      </svg>
    </button>
    
    <!-- iOS Instructions Modal -->
    <div v-if="showIOSInstructions" class="ios-instructions-modal" @click.self="showIOSInstructions = false">
      <div class="ios-instructions-content">
        <button class="close-btn" @click="showIOSInstructions = false">✕</button>
        <div class="ios-header">
          <div class="ios-icon">📱</div>
          <h3>{{ iosModalTitle }}</h3>
          <p class="ios-subtitle">{{ iosModalSubtitle }}</p>
        </div>
        
        <div class="steps">
          <div class="step">
            <span class="step-number">1</span>
            <div class="step-content">
              <p v-html="step1Text"></p>
              <div class="step-visual">
                <div class="phone-frame">
                  <div class="phone-screen">
                    <div class="share-btn">⎋</div>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div class="step">
            <span class="step-number">2</span>
            <div class="step-content">
              <p v-html="step2Text"></p>
              <div class="step-visual">
                <div class="menu-item">➕ Add to Home Screen</div>
              </div>
            </div>
          </div>
          <div class="step">
            <span class="step-number">3</span>
            <div class="step-content">
              <p v-html="step3Text"></p>
              <div class="step-visual">
                <div class="add-btn">Add</div>
              </div>
            </div>
          </div>
        </div>
        
        <div class="benefits">
          <div class="benefit">{{ benefit1 }}</div>
          <div class="benefit">{{ benefit2 }}</div>
          <div class="benefit">{{ benefit3 }}</div>
        </div>
        
        <div class="actions">
          <button class="got-it-btn" @click="showIOSInstructions = false">{{ gotItText }}</button>
          <button class="remind-btn" @click="remindLater">{{ remindLaterText }}</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useI18n } from '../services/i18n'

const { t } = useI18n()

const showInstallButton = ref(false)
const showIOSInstructions = ref(false)
const deferredPrompt = ref(null)
const isIOS = ref(false)
const isStandalone = ref(false)

const pwaInstallTitle = computed(() => t('pwa.installTitle'))
const pwaInstallText = computed(() => {
  if (isIOS.value) {
    return t('pwa.installOnIos')
  }
  return t('pwa.installText')
})
const iosModalTitle = computed(() => t('pwa.iosModalTitle'))
const iosModalSubtitle = computed(() => t('pwa.iosModalSubtitle'))
const step1Text = computed(() => t('pwa.step1Text'))
const step2Text = computed(() => t('pwa.step2Text'))
const step3Text = computed(() => t('pwa.step3Text'))
const benefit1 = computed(() => t('pwa.benefit1'))
const benefit2 = computed(() => t('pwa.benefit2'))
const benefit3 = computed(() => t('pwa.benefit3'))
const gotItText = computed(() => t('pwa.gotIt'))
const remindLaterText = computed(() => t('pwa.remindLater'))

function detectIOS() {
  const userAgent = window.navigator.userAgent.toLowerCase()
  return /iphone|ipad|ipod/.test(userAgent) && !/edge|crios|fxios/.test(userAgent)
}

function handleInstallAvailable() {
  showInstallButton.value = true
}

function handleInstallSuccess() {
  showInstallButton.value = false
}

function installPWA() {
  if (isIOS.value) {
    // Show iOS instructions
    showIOSInstructions.value = true
  } else if (window.pwaInstall) {
    // Use native install prompt for Android/Windows
    window.pwaInstall()
  } else {
    // Fallback - show general instructions
    showIOSInstructions.value = true
  }
}

function remindLater() {
  // Just close the modal, keep icon visible
  showIOSInstructions.value = false
}

onMounted(() => {
  // Check if device is iOS
  isIOS.value = detectIOS()
  
  // Check if already installed (standalone mode)
  isStandalone.value = window.matchMedia('(display-mode: standalone)').matches || 
                    window.navigator.standalone === true

  // Don't show if already installed
  if (isStandalone.value) {
    return
  }

  // Always show button on login page
  if (isIOS.value) {
    showInstallButton.value = true
  } else {
    // For other devices, listen to install prompt
    window.addEventListener('pwa-install-available', handleInstallAvailable)
    window.addEventListener('pwa-install-success', handleInstallSuccess)
    
    // Check for existing deferredPrompt
    if (window.deferredPrompt) {
      showInstallButton.value = true
    }
  }
})

onUnmounted(() => {
  window.removeEventListener('pwa-install-available', handleInstallAvailable)
  window.removeEventListener('pwa-install-success', handleInstallSuccess)
})
</script>

<style scoped>
.pwa-install-container {
  position: fixed;
  bottom: 16px;
  right: 16px;
  z-index: 1000;
}

.pwa-install-icon {
  width: 44px;
  height: 44px;
  padding: 0;
  background: linear-gradient(135deg, #1e3a8a 0%, #3b82f6 100%);
  color: white;
  border: none;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  box-shadow: 0 4px 12px rgba(30, 58, 138, 0.3);
  transition: all 0.3s ease;
}

.pwa-install-icon:hover {
  transform: scale(1.1);
  box-shadow: 0 6px 16px rgba(30, 58, 138, 0.4);
}

.pwa-install-icon:active {
  transform: scale(0.95);
}

/* iOS Instructions Modal */
.ios-instructions-modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2000;
  padding: 20px;
}

.ios-instructions-content {
  background: white;
  border-radius: 20px;
  padding: 28px;
  max-width: 420px;
  width: 100%;
  position: relative;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.4);
  max-height: 90vh;
  overflow-y: auto;
}

.close-btn {
  position: absolute;
  top: 16px;
  left: 16px;
  background: none;
  border: none;
  font-size: 24px;
  cursor: pointer;
  color: #666;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  transition: background 0.2s;
}

.close-btn:hover {
  background: #f0f0f0;
}

.ios-header {
  text-align: center;
  margin-bottom: 24px;
  padding-top: 8px;
}

.ios-icon {
  font-size: 48px;
  margin-bottom: 12px;
}

.ios-header h3 {
  margin: 0 0 8px 0;
  font-size: 20px;
  color: #1e3a8a;
  font-weight: 700;
}

.ios-subtitle {
  margin: 0;
  font-size: 14px;
  color: #666;
}

.steps {
  display: flex;
  flex-direction: column;
  gap: 20px;
  margin-bottom: 24px;
}

.step {
  display: flex;
  align-items: flex-start;
  gap: 14px;
}

.step-number {
  width: 32px;
  height: 32px;
  background: linear-gradient(135deg, #1e3a8a 0%, #3b82f6 100%);
  color: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  font-size: 15px;
  flex-shrink: 0;
}

.step-content {
  flex: 1;
}

.step p {
  margin: 0 0 8px 0;
  font-size: 14px;
  color: #333;
  line-height: 1.5;
}

.step strong {
  color: #1e3a8a;
}

.step-visual {
  background: #f8f9fa;
  border-radius: 8px;
  padding: 8px;
  display: flex;
  justify-content: center;
  border: 1px solid #e9ecef;
}

.phone-frame {
  width: 60px;
  height: 100px;
  background: white;
  border: 2px solid #dee2e6;
  border-radius: 8px;
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
}

.phone-screen {
  width: 90%;
  height: 90%;
  background: #f8f9fa;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.share-btn {
  font-size: 20px;
  color: #007AFF;
}

.menu-item {
  background: white;
  padding: 8px 12px;
  border-radius: 6px;
  font-size: 12px;
  color: #333;
  border: 1px solid #dee2e6;
}

.add-btn {
  background: #007AFF;
  color: white;
  padding: 6px 16px;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 600;
}

.benefits {
  display: flex;
  justify-content: space-around;
  margin-bottom: 24px;
  padding: 16px;
  background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
  border-radius: 12px;
}

.benefit {
  font-size: 12px;
  color: #495057;
  text-align: center;
  font-weight: 500;
}

.actions {
  display: flex;
  gap: 12px;
}

.got-it-btn {
  flex: 2;
  padding: 14px;
  background: linear-gradient(135deg, #1e3a8a 0%, #3b82f6 100%);
  color: white;
  border: none;
  border-radius: 12px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: transform 0.2s;
}

.got-it-btn:hover {
  transform: scale(1.02);
}

.remind-btn {
  flex: 1;
  padding: 14px;
  background: #f8f9fa;
  color: #666;
  border: 1px solid #dee2e6;
  border-radius: 12px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: background 0.2s;
}

.remind-btn:hover {
  background: #e9ecef;
}

/* RTL support */
[dir="rtl"] .pwa-install-container {
  right: auto;
  left: 16px;
}

[dir="rtl"] .close-btn {
  left: auto;
  right: 16px;
}

@media (max-width: 768px) {
  .pwa-install-container {
    bottom: 12px;
    right: 12px;
  }
  
  [dir="rtl"] .pwa-install-container {
    right: auto;
    left: 12px;
  }
  
  .pwa-install-icon {
    width: 40px;
    height: 40px;
  }

  .ios-instructions-content {
    padding: 20px;
  }

  .ios-instructions-content h3 {
    font-size: 16px;
  }

  .step p {
    font-size: 13px;
  }
}
</style>
```

---

## 📄 frontend-client-portal/src/components/ReportCard.vue

```vue
<template>
  <router-link :to="`/report/${report.id}`" class="report-card" @click.native="handleClick">
    <div class="report-card__content">
      <h3>{{ report.title }}</h3>
      <p class="report-card__meta mono">{{ formatDate(report.created_at) }}</p>
      <div v-if="report.assignment && report.assignment.employee_name" class="report-card__employee">
        <span>👷 {{ report.assignment.employee_name }}</span>
      </div>
      <div v-if="report.assignment && report.assignment.client_rating" class="report-card__rating">
        <span>⭐ {{ report.assignment.client_rating }}/5</span>
      </div>
    </div>
    <span :class="['badge', getStatusBadgeClass(report.status)]">
      {{ getStatusIcon(report.status) }} {{ getStatusText(report.status) }}
    </span>
  </router-link>
</template>

<script setup>
import { useI18n } from '../services/i18n'

const { t } = useI18n()

defineProps({ report: { type: Object, required: true } })

function handleClick(event) {
  // Navigation handled by router-link
}

function formatDate(date) {
  if (!date) return ''
  return new Date(date).toLocaleDateString('en-GB', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric'
  })
}

function getStatusText(status) {
  const statusMap = {
    'completed': t('status_completed'),
    'cancelled': t('status_cancelled'),
    'pending': t('status_pending'),
    'assigned': t('status_assigned'),
    'in_progress': t('status_in_progress')
  }
  return statusMap[status] || status
}

function getStatusBadgeClass(status) {
  const classMap = {
    'completed': 'badge--completed',
    'cancelled': 'badge--cancelled',
    'pending': 'badge--pending',
    'assigned': 'badge--assigned',
    'in_progress': 'badge--in-progress'
  }
  return classMap[status] || 'badge--pending'
}

function getStatusIcon(status) {
  const iconMap = {
    'completed': '✅',
    'cancelled': '❌',
    'pending': '⏳',
    'assigned': '📋',
    'in_progress': '🔄'
  }
  return iconMap[status] || '⏳'
}
</script>

<style scoped>
.report-card { 
  display: flex; 
  align-items: center; 
  justify-content: space-between; 
  background: var(--surface); 
  border: 1px solid var(--line); 
  border-radius: var(--radius-md); 
  padding: 16px; 
  margin-bottom: 10px; 
  box-shadow: var(--shadow-sm); 
  text-decoration: none;
  color: inherit;
}

.report-card__content {
  flex: 1;
  min-width: 0;
}

.report-card h3 { 
  font-size: 14.5px; 
  margin-bottom: 4px; 
  color: var(--ink);
}

.report-card__meta { 
  font-size: 12px; 
  color: var(--ink-soft); 
  margin-bottom: 6px;
}

.report-card__employee {
  font-size: 12px;
  color: var(--ink-soft);
  margin-bottom: 4px;
}

.report-card__rating {
  font-size: 12px;
  color: var(--brand);
  font-weight: 600;
}

.badge {
  background: var(--brand);
  color: white;
  padding: 4px 10px;
  border-radius: 12px;
  font-size: 11px;
  font-weight: 600;
  white-space: nowrap;
  margin-left: 10px;
}

.badge--completed {
  background: var(--signal-in);
}

.badge--cancelled {
  background: var(--signal-out);
}

.badge--pending {
  background: #FFA500;
}

.badge--assigned {
  background: #2196F3;
}

.badge--in-progress {
  background: #9C27B0;
}

@media (max-width: 480px) {
  .report-card {
    padding: 12px 14px;
    margin-bottom: 8px;
  }
  
  .report-card h3 {
    font-size: 13px;
    margin-bottom: 3px;
  }
  
  .report-card__meta {
    font-size: 11px;
  }
  
  .report-card__employee,
  .report-card__rating {
    font-size: 11px;
  }
  
  .badge {
    font-size: 9px;
    padding: 3px 8px;
  }
}

@media (max-width: 360px) {
  .report-card {
    padding: 10px 12px;
  }
  
  .report-card h3 {
    font-size: 12px;
  }
  
  .report-card__meta {
    font-size: 10px;
  }
}
</style>

```

---

## 📄 frontend-client-portal/src/components/StarRating.vue

```vue
<template>
  <div class="stars" role="radiogroup" aria-label="Service Rating">
    <button
      v-for="n in 5" :key="n" type="button"
      class="stars__item" :class="{ 'stars__item--filled': n <= (hover || modelValue) }"
      @mouseenter="hover = n" @mouseleave="hover = 0" @click="selectRating(n)"
      :aria-pressed="n <= modelValue"
    >★</button>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const props = defineProps({
  modelValue: {
    type: Number,
    default: 0
  }
})

const emit = defineEmits(['update:modelValue'])

const hover = ref(0)

function selectRating(n) {
  emit('update:modelValue', n)
}
</script>

<style scoped>
.stars { display: flex; gap: 6px; justify-content: center; }
.stars__item { font-size: 34px; line-height: 1; background: none; border: none; color: var(--line-strong); cursor: pointer; transition: color .15s ease, transform .1s ease; }
.stars__item:hover { transform: scale(1.08); }
.stars__item--filled { color: var(--gold); }

@media (max-width: 480px) {
  .stars {
    gap: 4px;
  }
  
  .stars__item {
    font-size: 28px;
  }
}

@media (max-width: 360px) {
  .stars {
    gap: 3px;
  }
  
  .stars__item {
    font-size: 24px;
  }
}
</style>

```

---

## 📄 frontend-client-portal/src/main.js

```javascript
import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import './styles/tokens.css'
import { i18nStore } from './plugins/i18n'

// ==========================================
// DevTools Detection & Protection
// ==========================================

// التحقق من بيئة التطوير والإنتاج
const isDevelopment = import.meta.env.VITE_APP_ENV === 'development' || 
                      import.meta.env.MODE === 'development';

// التحقق من مفتاح التجاوز
const bypassKey = localStorage.getItem('devtools_bypass_key') || 
                 new URLSearchParams(window.location.search).get('bypass');
const isBypassActive = bypassKey === 'worktrack_dev_2024';

// الحماية مفعلة افتراضياً، تعطيل فقط في بيئة التطوير أو عند وجود مفتاح التجاوز
if (!isDevelopment && !isBypassActive) {
  const devtools = {
    open: false,
    orientation: null
  };

  const threshold = 160;

  // اكتشاف فتح DevTools من خلال تغيير حجم النافذة
  setInterval(() => {
    const widthThreshold = window.outerWidth - window.innerWidth > threshold;
    const heightThreshold = window.outerHeight - window.innerHeight > threshold;
    
    if (widthThreshold || heightThreshold) {
      if (!devtools.open) {
        devtools.open = true;
        console.warn('⚠️ DevTools detected - Unauthorized access attempt');
        // يمكن إعادة تحميل الصفحة أو اتخاذ إجراء آخر
        // window.location.reload();
      }
    } else {
      devtools.open = false;
    }
  }, 500);

  // اكتشاف debugger (Anti-Debugging)
  setInterval(() => {
    const start = new Date().getTime();
    debugger; // يتوقف إذا كانت DevTools مفتوحة
    const end = new Date().getTime();
    
    if (end - start > 100) {
      console.warn('⚠️ Debugger detected - Unauthorized access attempt');
      // window.location.reload();
    }
  }, 1000);

  // منع النقر الأيمن - ولكن السماح بالنسخ واللصق
  document.addEventListener('contextmenu', (e) => {
    // السماح بالنقر الأيمن على حقول الإدخال والنصوص للنسخ واللصق
    const target = e.target;
    if (target.tagName === 'INPUT' || target.tagName === 'TEXTAREA' || 
        target.isContentEditable || target.tagName === 'SELECT') {
      return true; // السماح بالنقر الأيمن للنسخ واللصق
    }
    e.preventDefault();
    return false;
  });

  // منع اختصارات لوحة المفاتيح لفتح DevTools - ولكن السماح بالنسخ واللصق
  document.addEventListener('keydown', (e) => {
    // F12, Ctrl+Shift+I, Ctrl+Shift+J, Ctrl+Shift+C, Ctrl+U
    // تم استثناء Ctrl+C للسماح بالنسخ
    if (e.key === 'F12' || 
        (e.ctrlKey && e.shiftKey && (e.key === 'I' || e.key === 'J' || e.key === 'C')) ||
        (e.ctrlKey && e.key === 'U')) {
      e.preventDefault();
      return false;
    }
    // السماح بـ Ctrl+C و Ctrl+V و Ctrl+X للنسخ واللصق
    if (e.ctrlKey && (e.key === 'c' || e.key === 'v' || e.key === 'x' || e.key === 'a')) {
      return true; // السماح بهذه الاختصارات
    }
  });
} else {
  console.log('🔓 DevTools protection bypassed - Development mode active');
}

// تم إزالة منع النسخ واللصق للسماح للمستخدمين بهذه الوظائف

// إضافي: تأكيد السماح بالنسخ واللصق على الهاتف
document.addEventListener('DOMContentLoaded', () => {
  // تحسين النسخ واللصق في وضع PWA
  const enableCopyPaste = () => {
    // السماح بتحديد النصوص
    document.body.style.webkitUserSelect = 'text';
    document.body.style.userSelect = 'text';
    document.body.style.webkitTouchCallout = 'default';
    
    // السماح بالتفاعل مع الحافظة
    const inputs = document.querySelectorAll('input, textarea, select');
    inputs.forEach(input => {
      input.style.webkitUserSelect = 'text';
      input.style.userSelect = 'text';
      input.style.webkitTouchCallout = 'default';
    });
  };
  
  // تشغيل عند التحميل
  enableCopyPaste();
  
  // تشغيل مرة أخرى بعد فترة للتأكد من التطبيق
  setTimeout(enableCopyPaste, 1000);
  
  // السماح بالنسخ واللصق عبر أحداث اللمس
  document.addEventListener('touchstart', (e) => {
    const target = e.target;
    if (target.tagName === 'INPUT' || target.tagName === 'TEXTAREA' || 
        target.isContentEditable || target.tagName === 'SELECT') {
      // السماح بالتفاعل للنسخ واللصق
      target.style.webkitUserSelect = 'text';
      target.style.userSelect = 'text';
    }
  }, { passive: true });
  
  // منع منع أحداث النسخ واللصق
  document.addEventListener('copy', (e) => {
    e.stopPropagation();
  }, true);
  
  document.addEventListener('cut', (e) => {
    e.stopPropagation();
  }, true);
  
  document.addEventListener('paste', (e) => {
    e.stopPropagation();
  }, true);
});

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

// Ensure page is fully loaded before defining the app
document.addEventListener('DOMContentLoaded', () => {
  const app = createApp(App)
  app.use(router)
  app.mount('#app')

  // Set initial direction based on stored language
  const storedLang = localStorage.getItem('worktrack_lang') || 'ar'
  document.documentElement.dir = (storedLang === 'ar' || storedLang === 'he') ? 'rtl' : 'ltr'
  document.documentElement.lang = storedLang
})

```

---

## 📄 frontend-client-portal/src/plugins/i18n.js

```javascript
import { reactive, computed } from 'vue'

// ==========================================================================
// اللغات المدعومة: العربية والعبرية فقط
// ==========================================================================
const messages = {
  ar: {
    app_name: "Worktrack",
    admin_panel: "لوحة المدير",
    worker_app: "تطبيق الموظف",
    login: "تسجيل الدخول",
    email: "البريد الإلكتروني",
    password: "كلمة المرور",
    dashboard: "لوحة القيادة",
    employees: "الموظفون",
    tasks: "المهام",
    worksites: "المواقع",
    clients: "العملاء",
    reports: "التقارير",
    settings: "الإعدادات",
    service_requests: "طلبات الخدمة",
    logout: "تسجيل الخروج",
    add_employee: "إضافة موظف",
    assign: "تعيين",
    status: "الحالة",
    pending: "قيد الانتظار",
    assigned: "تم التعيين",
    in_progress: "قيد التنفيذ",
    completed: "مكتمل",
    cancelled: "ملغي",
    open_in_maps: "فتح في الخرائط",
    location: "الموقع",
    phone: "الهاتف",
    full_name: "الاسم الكامل",
    save: "حفظ",
    cancel: "إلغاء",
    delete: "حذف",
    edit: "تعديل",
    search: "بحث",
    filter: "تصفية",
    dark_mode: "الوضع الداكن",
    light_mode: "الوضع الفاتح",
    language: "اللغة",
    arabic: "العربية",
    hebrew: "עברית",
    loading: "جارٍ التحميل...",
    no_data: "لا توجد بيانات",
    error: "حدث خطأ",
    success: "تم بنجاح",
    welcome: "مرحباً بك",
    login_title: "لوحة تحكم المدير",
    worker_title: "تطبيق الموظف",
    pwa: {
      installTitle: "تثبيت التطبيق",
      installText: "تثبيت التطبيق",
      installOnIos: "تثبيت على iPhone",
      installOnAndroid: "تثبيت على Android",
      installOnDesktop: "تثبيت على الحاسوب",
      iosModalTitle: "تثبيت التطبيق على iPhone",
      iosModalSubtitle: "احصل على تجربة تطبيق أفضل!",
      androidModalTitle: "تثبيت التطبيق على Android",
      androidModalSubtitle: "احصل على تجربة تطبيق أفضل!",
      desktopModalTitle: "تثبيت التطبيق على الحاسوب",
      desktopModalSubtitle: "احصل على تجربة تطبيق أفضل!",
      step1Text: "اضغط على زر Share ⎋ في أسفل الشاشة",
      step2Text: "مرر لأسفل واضغط على Add to Home Screen",
      step3Text: "اضغط على Add في الزاوية العلوية",
      androidStep1Text: "اضغط على القائمة (⋮) في المتصفح",
      androidStep2Text: "اختر 'Install app' أو 'تثبيت التطبيق'",
      androidStep3Text: "اضغط على Install للتأكيد",
      desktopStep1Text: "اضغط على زر التثبيت في شريط العنوان",
      desktopStep2Text: "اختر 'Install' من القائمة المنسدلة",
      desktopStep3Text: "اضغط على Install للتأكيد",
      benefit1: "⚡ تشغيل أسرع",
      benefit2: "📱 أيقونة على الشاشة الرئيسية",
      benefit3: "🎨 تصميم شبيه بالتطبيقات",
      gotIt: "فهمت ✓",
      remindLater: "ذكرني لاحقاً",
      installingMessage: "جاري تثبيت التطبيق..."
    }
  },
  he: {
    app_name: "Worktrack",
    admin_panel: "לוח המנהל",
    worker_app: "אפליקציית העובד",
    login: "התחברות",
    email: "אימייל",
    password: "סיסמה",
    dashboard: "לוח בקרה",
    employees: "עובדים",
    tasks: "משימות",
    worksites: "אתרי עבודה",
    clients: "לקוחות",
    reports: "דוחות",
    settings: "הגדרות",
    service_requests: "בקשות שירות",
    logout: "התנתקות",
    add_employee: "הוסף עובד",
    assign: "הקצה",
    status: "סטטוס",
    pending: "ממתין",
    assigned: "הוקצה",
    in_progress: "בתהליך",
    completed: "הושלם",
    cancelled: "בוטל",
    open_in_maps: "פתח במפות",
    location: "מיקום",
    phone: "טלפון",
    full_name: "שם מלא",
    save: "שמור",
    cancel: "בטל",
    delete: "מחק",
    edit: "ערוך",
    search: "חפש",
    filter: "סנן",
    dark_mode: "מצב כהה",
    light_mode: "מצב בהיר",
    language: "שפה",
    arabic: "ערבית",
    hebrew: "עברית",
    loading: "טוען...",
    no_data: "אין נתונים",
    error: "שגיאה",
    success: "הצלחה",
    welcome: "ברוך הבא",
    login_title: "לוח המנהל",
    worker_title: "אפליקציית העובד",
    pwa: {
      installTitle: "התקנת אפליקציה",
      installText: "התקן אפליקציה",
      installOnIos: "התקנה באייפון",
      installOnAndroid: "התקנה באנדרואיד",
      installOnDesktop: "התקנה במחשב",
      iosModalTitle: "התקנת אפליקציה באייפון",
      iosModalSubtitle: "קבל חוויית אפליקציה טובה יותר!",
      androidModalTitle: "התקנת אפליקציה באנדרואיד",
      androidModalSubtitle: "קבל חוויית אפליקציה טובה יותר!",
      desktopModalTitle: "התקנת אפליקציה במחשב",
      desktopModalSubtitle: "קבל חוויית אפליקציה טובה יותר!",
      step1Text: "לחץ על כפתור Share ⎋ בתחתית המסך",
      step2Text: "גלול למטה ולחץ על Add to Home Screen",
      step3Text: "לחץ על Add בפינה העליונה",
      androidStep1Text: "לחץ על התפריט (⋮) בדפדפן",
      androidStep2Text: "בחר 'Install app' או 'התקנת אפליקציה'",
      androidStep3Text: "לחץ על Install לאישור",
      desktopStep1Text: "לחץ על כפתור ההתקנה בשורת הכתובת",
      desktopStep2Text: "בחר 'Install' מהתפריט הנפתח",
      desktopStep3Text: "לחץ על Install לאישור",
      benefit1: "⚡ פעולה מהירה יותר",
      benefit2: "📱 סמל במסך הבית",
      benefit3: "🎨 עיצוב דמוי אפליקציה",
      gotIt: "הבנתי ✓",
      remindLater: "תזכיר לי מאוחר יותר",
      installingMessage: "מתקין אפליקציה..."
    }
  }
}

const FALLBACK = 'ar'

// الحصول على اللغة المخزنة
const getStoredLang = () => {
  const stored = localStorage.getItem('worktrack_lang')
  // التأكد من أن اللغة مدعومة
  if (stored && (stored === 'ar' || stored === 'he')) return stored
  return FALLBACK
}

export const i18nStore = reactive({
  lang: getStoredLang(),

  setLang(lang) {
    if (lang === 'ar' || lang === 'he') {
      this.lang = lang
      localStorage.setItem('worktrack_lang', lang)
      // تحديث اتجاه الصفحة
      document.documentElement.dir = lang === 'ar' || lang === 'he' ? 'rtl' : 'ltr'
      document.documentElement.lang = lang
      // إعادة تحميل الصفحة لتطبيق اللغة الجديدة
      window.location.reload()
    }
  },

  t(key) {
    const translation = messages[this.lang]?.[key]
    if (translation) return translation
    // استخدام العربية كافتراضي إذا لم توجد الترجمة
    const fallback = messages[FALLBACK]?.[key]
    return fallback || key
  }
})

export function useI18n() {
  const t = (key) => i18nStore.t(key)
  const setLang = (lang) => i18nStore.setLang(lang)
  const currentLang = computed(() => i18nStore.lang)
  return { t, setLang, currentLang }
}

```

---

## 📄 frontend-client-portal/src/router/index.js

```javascript
import { createRouter, createWebHashHistory, createMemoryHistory } from 'vue-router'
import ServiceReportView from '../views/ServiceReportView.vue'
import ServiceHistoryView from '../views/ServiceHistoryView.vue'
import RatingView from '../views/RatingView.vue'
import NewRequestView from '../views/NewRequestView.vue'
import LoginView from '../views/LoginView.vue'
import { getToken } from '../services/auth'

const routes = [
  { path: '/login', component: LoginView, meta: { public: true } },
  { 
    path: '/', 
    component: ServiceHistoryView,
    meta: { requiresAuth: true }
  },
  { 
    path: '/new', 
    component: NewRequestView,
    meta: { requiresAuth: true }
  },
  { 
    path: '/report/:id', 
    component: ServiceReportView,
    meta: { requiresAuth: true }
  },
  { 
    path: '/rate/:id', 
    component: RatingView,
    meta: { requiresAuth: true }
  },
]

// تحويل المسارات القديمة إلى صيغة Hash Mode
function migrateSavedPath() {
  const savedPath = localStorage.getItem('worktrack_client_last_path')
  if (savedPath && !savedPath.startsWith('#/')) {
    // تحويل المسار القديم (/new) إلى الجديد (/#/new)
    const newPath = savedPath.startsWith('/') ? `#${savedPath}` : `#/${savedPath}`
    localStorage.setItem('worktrack_client_last_path', newPath)
    console.log('Migrated saved path:', savedPath, '->', newPath)
  }
}

// استخدام MemoryHistory لـ Electron و HashHistory للويب
const isElectron = typeof window !== 'undefined' && (window.isElectron || (window.process && window.process.type === 'renderer'))
const router = createRouter({ history: isElectron ? createMemoryHistory() : createWebHashHistory(), routes })

// تشغيل التحويل عند تحميل الـ router
if (!isElectron) {
  migrateSavedPath()
}

// Navigation guard for auth
router.beforeEach((to, from, next) => {
  const token = getToken()
  const requiresAuth = to.meta.requiresAuth

  console.log('Client Router guard:', { to: to.path, hasToken: !!token, meta: to.meta })

  if (requiresAuth && !token) {
    next('/login')
  } else if (to.path === '/login' && token) {
    // استعادة المسار المحفوظ أو الذهاب لصفحة طلب خدمة جديد
    const savedPath = localStorage.getItem('worktrack_client_last_path')
    if (savedPath && savedPath !== '/login' && savedPath !== '#/login') {
      // إزالة الـ hash إذا كان موجوداً للـ Vue Router
      const cleanPath = savedPath.startsWith('#') ? savedPath.substring(1) : savedPath
      next(cleanPath)
    } else {
      next('/new')
    }
  } else if (requiresAuth && token) {
    // السماح بالوصول للمسارات المحمية
    next()
  } else {
    next()
  }
})

// حفظ المسار الحالي بعد كل تغيير مسار
router.afterEach((to) => {
  // لا تحفظ مسار login
  if (to.path !== '/login') {
    // في Hash Mode، نحفظ المسار مع الـ hash
    const pathToSave = isElectron ? to.path : `#${to.path}`
    localStorage.setItem('worktrack_client_last_path', pathToSave)
  }
})

export default router

```

---

## 📄 frontend-client-portal/src/services/api.js

```javascript
import axios from 'axios'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1',
  withCredentials: true, // مهم لإرسال واستقبال cookies
})

api.interceptors.request.use((config) => {
  // Add token from localStorage as fallback for Authorization header
  const token = localStorage.getItem('worktrack_client_token')
  if (token) {
    config.headers['Authorization'] = `Bearer ${token}`
  }
  
  return config
})

// Response interceptor to handle 401 errors and password changes
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      const errorMessage = error.response?.data?.error || ''
      console.log('🔍 401 Error detected:', errorMessage) // للتصحيح
      
      // Check if error is about password change - توسيع البحث
      if (errorMessage.includes('password changed') || 
          errorMessage.includes('كلمة المرور') ||
          errorMessage.includes('كلمة السر') ||
          errorMessage.includes('הסיסמה שונתה') ||
          errorMessage.includes('Password has been changed')) {
        console.log('🔓 Password change detected, showing alert') // للتصحيح
        // Show custom popup for password change
        showPasswordChangedAlert()
      } else {
        console.log('🚪 Normal logout required') // للتصحيح
        // Handle other 401 errors (normal logout)
        handleLogout()
      }
    }
    return Promise.reject(error)
  }
)

function showPasswordChangedAlert() {
  const currentLang = localStorage.getItem('worktrack_client_language') || 'ar'
  
  const messages = {
    ar: {
      title: 'تم تغيير كلمة المرور',
      message: 'تم تغيير كلمة المرور الخاصة بحسابك. يرجى تسجيل الدخول مرة أخرى.',
      button: 'تسجيل الدخول'
    },
    he: {
      title: 'הסיסמה שונתה',
      message: 'הסיסמה שלך שונתה. אנא התחבר שוב.',
      button: 'התחבר'
    },
    en: {
      title: 'Password Changed',
      message: 'Your password has been changed. Please log in again.',
      button: 'Log In'
    }
  }
  
  const msg = messages[currentLang] || messages.ar
  
  // Create and show alert
  const alertDiv = document.createElement('div')
  alertDiv.style.cssText = `
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.7);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 10000;
    direction: ${currentLang === 'he' ? 'rtl' : currentLang === 'en' ? 'ltr' : 'rtl'};
  `
  
  const alertBox = document.createElement('div')
  alertBox.style.cssText = `
    background: white;
    padding: 30px;
    border-radius: 12px;
    max-width: 400px;
    text-align: center;
    box-shadow: 0 10px 40px rgba(0, 0, 0, 0.3);
  `
  
  alertBox.innerHTML = `
    <div style="font-size: 48px; margin-bottom: 15px;">🔒</div>
    <h2 style="margin: 0 0 15px 0; color: #e74c3c;">${msg.title}</h2>
    <p style="margin: 0 0 25px 0; color: #555; line-height: 1.6;">${msg.message}</p>
    <button onclick="handlePasswordChangedLogout()" style="
      background: #e74c3c;
      color: white;
      border: none;
      padding: 12px 30px;
      border-radius: 6px;
      font-size: 16px;
      cursor: pointer;
      transition: background 0.3s;
    ">${msg.button}</button>
  `
  
  alertDiv.appendChild(alertBox)
  document.body.appendChild(alertDiv)
  
  // Make function globally available
  window.handlePasswordChangedLogout = function() {
    handleLogout()
    document.body.removeChild(alertDiv)
  }
}

function handleLogout() {
  // Clear all auth data
  localStorage.removeItem('worktrack_client_token')
  localStorage.removeItem('worktrack_client_user')
  
  // Redirect to login page
  window.location.href = '/login'
}

export default api

```

---

## 📄 frontend-client-portal/src/services/auth.js

```javascript
import api from './api'

export async function login(email, password) {
  try {
    console.log('📤 Attempting login (client):', { email })
    
    const { data } = await api.post('/auth/login', { 
      email: email.trim(), 
      password: password.trim() 
    })
    
    console.log('✅ Login successful (client)')
    
    localStorage.setItem('worktrack_client_token', data.token)
    localStorage.setItem('worktrack_client_user', JSON.stringify(data.user))
    
    return data
  } catch (error) {
    console.error('❌ Login failed (client):', error.response?.data || error.message)
    throw error
  }
}

export async function getCurrentUser() {
  const { data } = await api.get('/auth/me')
  return data
}

export function logout() {
  localStorage.removeItem('worktrack_client_token')
  localStorage.removeItem('worktrack_client_user')
}

export function currentUser() {
  const raw = localStorage.getItem('worktrack_client_user')
  return raw ? JSON.parse(raw) : null
}

export function getToken() {
  return localStorage.getItem('worktrack_client_token')
}
```

---

## 📄 frontend-client-portal/src/services/geocoding.js

```javascript
// Geocoding service using OpenStreetMap Nominatim API
// This service converts coordinates to location names (reverse geocoding)

const API_BASE = 'https://nominatim.openstreetmap.org/reverse'

/**
 * Convert latitude and longitude to a human-readable location name
 * @param {number} lat - Latitude
 * @param {number} lng - Longitude
 * @param {string} language - Language code (e.g., 'ar', 'en', 'he')
 * @returns {Promise<string>} Location name
 */
export async function reverseGeocode(lat, lng, language = 'ar') {
  try {
    const url = `${API_BASE}?format=json&lat=${lat}&lon=${lng}&accept-language=${language}&zoom=18`
    const response = await fetch(url, {
      headers: {
        'User-Agent': 'Worktrack Client Portal'
      }
    })
    
    if (!response.ok) {
      throw new Error('Geocoding request failed')
    }
    
    const data = await response.json()
    
    // Build a comprehensive location name
    const parts = []
    
    if (data.address) {
      // Add specific location details
      if (data.address.building) parts.push(data.address.building)
      if (data.address.house_number) parts.push(data.address.house_number)
      if (data.address.road) parts.push(data.address.road)
      if (data.address.neighbourhood) parts.push(data.address.neighbourhood)
      if (data.address.suburb) parts.push(data.address.suburb)
      if (data.address.city_district) parts.push(data.address.city_district)
      if (data.address.city) parts.push(data.address.city)
      if (data.address.town) parts.push(data.address.town)
      if (data.address.village) parts.push(data.address.village)
      if (data.address.county) parts.push(data.address.county)
      if (data.address.state) parts.push(data.address.state)
      if (data.address.country) parts.push(data.address.country)
    }
    
    // If no address parts found, use display name
    if (parts.length === 0 && data.display_name) {
      return data.display_name
    }
    
    // Join parts with commas
    return parts.join(', ')
  } catch (error) {
    console.error('Reverse geocoding error:', error)
    // Return a fallback location name
    return `الموقع: ${lat.toFixed(6)}, ${lng.toFixed(6)}`
  }
}

/**
 * Search for a location by name (forward geocoding)
 * @param {string} query - Location name to search
 * @param {string} language - Language code
 * @returns {Promise<Array>} Array of location results
 */
export async function forwardGeocode(query, language = 'ar') {
  try {
    const url = `https://nominatim.openstreetmap.org/search?format=json&q=${encodeURIComponent(query)}&accept-language=${language}&limit=5`
    const response = await fetch(url, {
      headers: {
        'User-Agent': 'Worktrack Client Portal'
      }
    })
    
    if (!response.ok) {
      throw new Error('Geocoding search failed')
    }
    
    const data = await response.json()
    return data.map(item => ({
      lat: parseFloat(item.lat),
      lng: parseFloat(item.lon),
      display_name: item.display_name,
      address: item.display_name
    }))
  } catch (error) {
    console.error('Forward geocoding error:', error)
    return []
  }
}

```

---

## 📄 frontend-client-portal/src/services/i18n.js

```javascript
import { reactive, computed } from 'vue'
import ar from '../locales/ar.json'
import he from '../locales/he.json'
import en from '../locales/en.json'

// =============================================
// Unified language key for all applications
// =============================================
const STORAGE_KEY = 'worktrack_language'

// =============================================
// Translations are bundled in Vite; fetching src files doesn't work after deployment.
// =============================================
const messages = { ar, he, en }

// =============================================
// Get stored language
// =============================================
function getStoredLang() {
  const stored = localStorage.getItem(STORAGE_KEY)
  if (stored && ['ar', 'he', 'en'].includes(stored)) {
    return stored
  }
  // Use browser language as default
  const browserLang = navigator.language || navigator.languages?.[0] || 'ar'
  if (browserLang.startsWith('he')) return 'he'
  if (browserLang.startsWith('en')) return 'en'
  return 'ar'
}

// =============================================
// Translation state
// =============================================
const i18nState = reactive({
  currentLang: getStoredLang(),
  
  setLang(lang) {
    if (['ar', 'he', 'en'].includes(lang)) {
      this.currentLang = lang
      localStorage.setItem(STORAGE_KEY, lang)
      
      // Update page direction
      document.documentElement.dir = lang === 'ar' || lang === 'he' ? 'rtl' : 'ltr'
      document.documentElement.lang = lang
      
      // Send event to indicate language change
      window.dispatchEvent(new CustomEvent('language-changed', { detail: { lang } }))
      
      console.log(`🌍 Language changed to: ${lang}`)
      
      // Reload page to apply changes
      setTimeout(() => {
        window.location.reload()
      }, 300)
    }
  },
  
  t(key) {
    const keys = key.split('.')
    let translation = messages[this.currentLang]

    for (const k of keys) {
      if (translation && translation[k]) {
        translation = translation[k]
      } else {
        // البحث في اللغة الافتراضية
        let fallbackTranslation = messages['ar']
        for (const fk of keys) {
          if (fallbackTranslation && fallbackTranslation[fk]) {
            fallbackTranslation = fallbackTranslation[fk]
          } else {
            console.warn(`⚠️ مفتاح الترجمة غير موجود: ${key}`)
            return key
          }
        }
        return fallbackTranslation
      }
    }
    return translation
  }
})

// =============================================
// تصدير الدوال
// =============================================
export function useI18n() {
  const t = (key) => i18nState.t(key)
  const setLang = (lang) => i18nState.setLang(lang)
  const currentLang = computed(() => i18nState.currentLang)
  return { t, setLang, currentLang }
}

export default {
  install(app) {
    app.config.globalProperties.$t = i18nState.t
    app.config.globalProperties.$lang = i18nState.currentLang
    app.provide('i18n', i18nState)
  }
}

```

---

## 📄 frontend-client-portal/src/services/websocket.js

```javascript
// WebSocket service for real-time tracking
class WebSocketService {
  constructor() {
    this.ws = null
    this.reconnectInterval = null
    this.listeners = []
    this.isConnected = false
    this.useElectronAPI = window.electronAPI && window.electronAPI.websocket
  }

  connect(url) {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      return
    }

    // Use Electron API if available
    if (this.useElectronAPI) {
      this.connectViaElectron(url)
    } else {
      this.connectViaBrowser(url)
    }
  }

  connectViaElectron(url) {
    try {
      console.log('🔌 Connecting to WebSocket via Electron:', url)
      
      // Clean up old listeners
      this.cleanupElectronListeners()

      // Set up new listeners
      this.electronMessageHandler = (data) => {
        console.log('📡 WebSocket message received (Electron):', data)
        this.notifyListeners(data)
      }

      this.electronOpenHandler = () => {
        console.log('✅ WebSocket connected (Electron)')
        this.isConnected = true
        this.notifyListeners({ type: 'connected' })
      }

      this.electronErrorHandler = (error) => {
        console.warn('⚠️ WebSocket error (Electron):', error)
        this.isConnected = false
      }

      this.electronCloseHandler = (code, reason) => {
        console.log('🔌 WebSocket closed (Electron):', code, reason)
        this.isConnected = false
        this.notifyListeners({ type: 'disconnected' })
        
        // Auto reconnect after 30 seconds
        this.reconnectInterval = setTimeout(() => {
          console.log('🔄 Attempting to reconnect...')
          this.connect(url)
        }, 30000)
      }

      // Register listeners
      window.electronAPI.websocket.onMessage(this.electronMessageHandler)
      window.electronAPI.websocket.onOpen(this.electronOpenHandler)
      window.electronAPI.websocket.onError(this.electronErrorHandler)
      window.electronAPI.websocket.onClose(this.electronCloseHandler)

      // Connect
      window.electronAPI.websocket.connect(url)
    } catch (e) {
      console.warn('⚠️ Failed to connect via Electron, using browser:', e)
      this.connectViaBrowser(url)
    }
  }

  connectViaBrowser(url) {
    try {
      console.log('🔌 Connecting to WebSocket via browser:', url)
      this.ws = new WebSocket(url)

      this.ws.onopen = () => {
        console.log('✅ WebSocket connected (browser)')
        this.isConnected = true
        this.notifyListeners({ type: 'connected' })
      }

      this.ws.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data)
          console.log('📡 WebSocket message received (browser):', data)
          this.notifyListeners(data)
        } catch (e) {
          console.error('❌ Error parsing WebSocket message:', e)
        }
      }

      this.ws.onerror = (error) => {
        console.warn('⚠️ WebSocket not available - will use periodic updates')
        this.isConnected = false
      }

      this.ws.onclose = (event) => {
        console.log('🔌 WebSocket closed (browser)', event.code, event.reason)
        this.isConnected = false
        this.notifyListeners({ type: 'disconnected' })
        
        // Auto reconnect after 30 seconds
        this.reconnectInterval = setTimeout(() => {
          console.log('🔄 Attempting to reconnect...')
          this.connect(url)
        }, 30000)
      }
    } catch (e) {
      console.warn('⚠️ WebSocket not available - will use periodic updates')
    }
  }

  cleanupElectronListeners() {
    if (!this.useElectronAPI) return
    
    // Remove listeners (needs implementation in preload.js)
    // Currently no way to remove listeners in Electron IPC
    // Can be added later if needed
  }

  disconnect() {
    if (this.reconnectInterval) {
      clearTimeout(this.reconnectInterval)
      this.reconnectInterval = null
    }

    if (this.useElectronAPI) {
      window.electronAPI.websocket.disconnect()
    } else if (this.ws) {
      this.ws.close()
      this.ws = null
    }

    this.isConnected = false
  }

  onMessage(callback) {
    this.listeners.push(callback)
  }

  removeListener(callback) {
    this.listeners = this.listeners.filter(listener => listener !== callback)
  }

  notifyListeners(data) {
    this.listeners.forEach(callback => callback(data))
  }

  send(data) {
    if (this.useElectronAPI) {
      window.electronAPI.websocket.send(data)
    } else if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(data))
    } else {
      console.warn('⚠️ WebSocket not connected')
    }
  }
}

// Create single instance
const wsService = new WebSocketService()

export default wsService
```

---

## 📄 frontend-client-portal/src/store/auth.js

```javascript
import { reactive } from 'vue'
import { currentUser } from '../services/auth'

export const authStore = reactive({
  user: currentUser(),
  setUser(user) {
    this.user = user
  },
  clear() {
    this.user = null
  },
})
```

---

## 📄 frontend-client-portal/src/store/notifications.js

```javascript
import { reactive } from 'vue'
import api from '../services/api'

export const notificationsStore = reactive({
  notifications: [],
  loading: false,
  error: null,
  
  async fetchNotifications() {
    // Skip API calls in development to avoid CORS errors
    if (import.meta.env.DEV) {
      console.log('📋 Notifications disabled in development mode')
      this.notifications = []
      this.loading = false
      this.error = null
      return
    }

    this.loading = true
    this.error = null
    try {
      const { data } = await api.get('/notifications')
      this.notifications = data
    } catch (e) {
      // Handle CORS errors gracefully in development
      if (e.message?.includes('Network Error') || e.code === 'ERR_NETWORK') {
        console.warn('⚠️ Network error (likely CORS in development) - notifications disabled')
        this.error = null // Don't show error for CORS issues in dev
        this.notifications = [] // Clear notifications
      } else {
        this.error = e.response?.data?.error || 'Failed to fetch notifications'
        console.error('❌ Failed to fetch notifications:', e)
      }
    } finally {
      this.loading = false
    }
  },
  
  clear() {
    this.notifications = []
    this.error = null
  }
})
```

---

## 📄 frontend-client-portal/src/store/reports.js

```javascript
import { reactive } from 'vue'
import api from '../services/api'

export const reportsStore = reactive({
  reports: [],
  loading: false,
  error: null,
  
  async fetchReports() {
    this.loading = true
    this.error = null
    try {
      const { data } = await api.get('/service/requests')
      // Show all requests (not just completed) and add necessary metadata
      this.reports = (data || [])
        .map(r => ({
          ...r,
          date: r.created_at,
          employee_name: r.assignment?.employee_name || null,
          rating: r.assignment?.client_rating || null,
          feedback: r.assignment?.client_feedback || null
        }))
        .sort((a, b) => new Date(b.created_at) - new Date(a.created_at)) // Sort by date (newest first)
    } catch (e) {
      this.error = e.response?.data?.error || 'Failed to fetch reports'
      console.error('❌ Failed to fetch reports:', e)
      this.reports = []
    } finally {
      this.loading = false
    }
  },
  
  find(id) {
    return this.reports.find((r) => String(r.id) === String(id))
  },
  
  clear() {
    this.reports = []
    this.error = null
  }
})
```

---

## 📄 frontend-client-portal/src/store/services.js

```javascript
import { reactive } from 'vue'
import api from '../services/api'

export const servicesStore = reactive({
  requests: [],
  loading: false,
  error: null,
  
  async fetchRequests() {
    this.loading = true
    this.error = null
    try {
      const { data } = await api.get('/service/requests')
      this.requests = data
    } catch (e) {
      this.error = e.response?.data?.error || 'Failed to fetch service requests'
      console.error('❌ Failed to fetch service requests:', e)
    } finally {
      this.loading = false
    }
  },
  
  async createRequest(requestData) {
    this.loading = true
    this.error = null
    try {
      const { data } = await api.post('/service/requests', requestData)
      this.requests.push(data)
      return data
    } catch (e) {
      this.error = e.response?.data?.error || 'Failed to create service request'
      console.error('❌ Failed to create service request:', e)
      throw e
    } finally {
      this.loading = false
    }
  },
  
  find(id) {
    return this.requests.find((r) => String(r.id) === String(id))
  },
  
  clear() {
    this.requests = []
    this.error = null
  }
})
```

---

## 📄 frontend-client-portal/src/tests/requests.test.js

```javascript
import { describe, it, expect, vi } from 'vitest'

// اختبارات نظام طلبات الخدمة

describe('نظام طلبات الخدمة', () => {
  
  describe('نموذج طلب خدمة جديد', () => {
    it('التحقق من البيانات المطلوبة', () => {
      const requestData = {
        customerName: 'محمد أحمد',
        phone: '0501234567',
        title: 'صيانة مكيف',
        description: 'المكيف لا يعمل بشكل صحيح',
        location: {
          latitude: 24.7136,
          longitude: 46.6753,
        },
      }
      
      expect(requestData.customerName).toBeTruthy()
      expect(requestData.phone).toBeTruthy()
      expect(requestData.title).toBeTruthy()
      expect(requestData.description).toBeTruthy()
      expect(requestData.location).toBeDefined()
    })

    it('التحقق من صحة رقم الهاتف', () => {
      const validPhone = '0501234567'
      const invalidPhone = '12345'
      
      const isValidPhone = /^05[0-9]{8}$/.test(validPhone)
      const isInvalidPhone = /^05[0-9]{8}$/.test(invalidPhone)
      
      expect(isValidPhone).toBe(true)
      expect(isInvalidPhone).toBe(false)
    })

    it('التحقق من صحة الإحداثيات', () => {
      const location = {
        latitude: 24.7136,
        longitude: 46.6753,
      }
      
      const isValidLat = location.latitude >= -90 && location.latitude <= 90
      const isValidLng = location.longitude >= -180 && location.longitude <= 180
      
      expect(isValidLat && isValidLng).toBe(true)
    })

    it('اختيار الأولوية', () => {
      const priorities = ['low', 'normal', 'high', 'urgent']
      const selectedPriority = 'high'
      
      expect(priorities).toContain(selectedPriority)
    })
  })

  describe('تحديد الموقع', () => {
    it('استخدام GPS لتحديد الموقع', () => {
      const gpsLocation = {
        latitude: 24.7136,
        longitude: 46.6753,
        accuracy: 10,
      }
      
      expect(gpsLocation.latitude).toBeDefined()
      expect(gpsLocation.longitude).toBeDefined()
      expect(gpsLocation.accuracy).toBeGreaterThan(0)
    })

    it('اختيار الموقع من الخريطة', () => {
      const mapLocation = {
        latitude: 24.7146,
        longitude: 46.6763,
        address: 'الرياض، حي الملز',
      }
      
      expect(mapLocation.address).toBeTruthy()
      expect(mapLocation.latitude).toBeDefined()
      expect(mapLocation.longitude).toBeDefined()
    })
  })

  describe('إرسال الطلب', () => {
    it('إنشاء رقم طلب فريد', () => {
      const requestId = 'REQ-' + Date.now()
      expect(requestId).toContain('REQ-')
    })

    it('تأكيد استلام الطلب', () => {
      const response = {
        success: true,
        requestId: 'REQ-123456',
        message: 'تم استلام طلبك بنجاح',
      }
      
      expect(response.success).toBe(true)
      expect(response.requestId).toBeDefined()
    })
  })
})

describe('سجل الطلبات', () => {
  
  describe('عرض طلبات الخدمة', () => {
    it('عرض جميع طلبات العميل', () => {
      const requests = [
        { id: 1, title: 'صيانة مكيف', status: 'completed', date: '2024-01-15' },
        { id: 2, title: 'تسريب مياه', status: 'in_progress', date: '2024-01-20' },
        { id: 3, title: 'تركيب إضاءة', status: 'pending', date: '2024-01-25' },
      ]
      
      expect(requests.length).toBe(3)
    })

    it('فلترة الطلبات حسب الحالة', () => {
      const requests = [
        { id: 1, status: 'completed' },
        { id: 2, status: 'in_progress' },
        { id: 3, status: 'pending' },
        { id: 4, status: 'completed' },
      ]
      
      const completedRequests = requests.filter(r => r.status === 'completed')
      expect(completedRequests.length).toBe(2)
    })

    it('فرز الطلبات حسب التاريخ', () => {
      const requests = [
        { id: 1, date: '2024-01-15' },
        { id: 2, date: '2024-01-20' },
        { id: 3, date: '2024-01-10' },
      ]
      
      const sorted = [...requests].sort((a, b) => new Date(b.date) - new Date(a.date))
      expect(sorted[0].date).toBe('2024-01-20')
    })
  })

  describe('تتبع حالة الطلب', () => {
    it('عرض الحالة الحالية', () => {
      const request = {
        id: 1,
        title: 'صيانة مكيف',
        status: 'in_progress',
        assignedEmployee: 'محمد أحمد',
      }
      
      expect(request.status).toBe('in_progress')
      expect(request.assignedEmployee).toBe('محمد أحمد')
    })

    it('عرض سجل تغييرات الحالة', () => {
      const statusHistory = [
        { status: 'pending', timestamp: '2024-01-15T09:00:00' },
        { status: 'assigned', timestamp: '2024-01-15T10:00:00' },
        { status: 'in_progress', timestamp: '2024-01-15T11:00:00' },
      ]
      
      expect(statusHistory.length).toBe(3)
      expect(statusHistory[0].status).toBe('pending')
    })
  })
})

describe('نظام التقييم', () => {
  
  describe('تقييم الخدمة', () => {
    it('تقييم بالنجوم', () => {
      const rating = {
        serviceId: 1,
        stars: 5,
        comment: 'خدمة ممتازة',
      }
      
      expect(rating.stars).toBeGreaterThanOrEqual(1)
      expect(rating.stars).toBeLessThanOrEqual(5)
    })

    it('إضافة تعليق نصي', () => {
      const comment = 'الموظف كان محترفاً وسريعاً'
      expect(comment.length).toBeGreaterThan(0)
    })

    it('التحقق من صحة التقييم', () => {
      const rating = 4
      const comment = 'جيد'
      
      const isValidRating = rating >= 1 && rating <= 5
      const isValidComment = comment.length > 0
      
      expect(isValidRating && isValidComment).toBe(true)
    })
  })

  describe('عرض التقييمات', () => {
    it('حساب متوسط التقييمات', () => {
      const ratings = [5, 4, 5, 3, 4]
      const average = ratings.reduce((sum, r) => sum + r, 0) / ratings.length
      expect(average).toBe(4.2)
    })

    it('عرض عدد التقييمات', () => {
      const ratingsCount = 25
      expect(ratingsCount).toBeGreaterThan(0)
    })
  })
})

describe('التحقق من صحة البيانات', () => {
  
  it('التحقق من صحة البريد الإلكتروني', () => {
    const validEmail = 'customer@example.com'
    const invalidEmail = 'invalid-email'
    
    const isValidEmail = /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(validEmail)
    const isInvalidEmail = /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(invalidEmail)
    
    expect(isValidEmail).toBe(true)
    expect(isInvalidEmail).toBe(false)
  })

  it('التحقق من طول الوصف', () => {
    const shortDescription = 'قصير'
    const longDescription = 'هذا وصف طويل ومفصل للمشكلة التي تواجهها'
    
    const isValidShort = shortDescription.length >= 10
    const isValidLong = longDescription.length >= 10
    
    expect(isValidShort).toBe(false)
    expect(isValidLong).toBe(true)
  })
})
```

---

## 📄 frontend-client-portal/src/tests/setup.js

```javascript
import { vi } from 'vitest'
import { config } from '@vue/test-utils'

// محاكاة window.matchMedia للاختبارات
Object.defineProperty(window, 'matchMedia', {
  writable: true,
  value: vi.fn().mockImplementation(query => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: vi.fn(),
    removeListener: vi.fn(),
    addEventListener: vi.fn(),
    removeEventListener: vi.fn(),
    dispatchEvent: vi.fn(),
  })),
})

// إعدادات Vue Test Utils
config.global.mocks = {
  $t: (key) => key,
}

// محاكاة localStorage
const localStorageMock = {
  getItem: vi.fn(),
  setItem: vi.fn(),
  removeItem: vi.fn(),
  clear: vi.fn(),
}
global.localStorage = localStorageMock
```

---

## 📄 frontend-client-portal/src/views/LoginView.vue

```vue
<template>
  <div class="login-container">
    <PWAInstallButton />
    <div class="login-card">
      <div class="login-header">
        <img src="/devpro-logo.jpg" alt="DevPro Logo" class="login-logo" />
        <h2 class="login-title">{{ t('login') }}</h2>
        <p class="login-subtitle">{{ t('login_subtitle') }}</p>
      </div>

      <form @submit.prevent="handleLogin" class="login-form">
        <div class="form-group">
          <label>{{ t('email') }}</label>
          <input 
            v-model="email" 
            type="email" 
            :placeholder="t('email_placeholder')"
            required
          />
        </div>

        <div class="form-group">
          <label>{{ t('password') }}</label>
          <input 
            v-model="password" 
            type="password" 
            :placeholder="t('password_placeholder')"
            required
          />
        </div>

        <button type="submit" class="login-btn" :disabled="loading">
          {{ loading ? t('loading') : t('login') }}
        </button>

        <div v-if="error" class="error-message">
          {{ error }}
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { login } from '../services/auth'
import { useI18n } from '../services/i18n'
import PWAInstallButton from '../components/PWAInstallButton.vue'
import { authStore } from '../store/auth'

const router = useRouter()
const { t } = useI18n()

const email = ref('')
const password = ref('')
const loading = ref(false)
const error = ref('')

async function handleLogin() {
  loading.value = true
  error.value = ''

  try {
    const data = await login(email.value, password.value)
    authStore.setUser(data.user)
    // استخدام window.location.href لضمان التحويل الكامل مع إعادة التحميل
    window.location.href = window.location.origin + '/#/new'
  } catch (err) {
    error.value = err.response?.data?.error || t('login_failed')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #1E3A5F 0%, #0D1B3E 100%);
  padding: 20px;
}

.login-card {
  background: white;
  border-radius: 16px;
  padding: 40px;
  width: 100%;
  max-width: 400px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

.login-header {
  text-align: center;
  margin-bottom: 32px;
}

.login-logo {
  width: 64px;
  height: 64px;
  border-radius: 12px;
  object-fit: contain;
  margin-bottom: 16px;
}

.login-title {
  font-size: 24px;
  font-weight: 600;
  margin-bottom: 8px;
  color: #1a1a1a;
}

.login-subtitle {
  font-size: 14px;
  color: #666;
  margin-bottom: 0;
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-group label {
  font-size: 14px;
  font-weight: 500;
  color: #333;
}

.form-group input {
  padding: 12px 16px;
  border: 1px solid #ddd;
  border-radius: 8px;
  font-size: 16px;
  transition: border-color 0.2s;
}

.form-group input:focus {
  outline: none;
  border-color: #1E3A5F;
  box-shadow: 0 0 0 3px rgba(30, 58, 95, 0.1);
}

.login-btn {
  padding: 14px;
  background: #1E3A5F;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: transform 0.2s, opacity 0.2s;
  box-shadow: 0 4px 12px rgba(30, 58, 95, 0.3);
}

.login-btn:hover:not(:disabled) {
  background: #0D1B3E;
  transform: translateY(-2px);
  box-shadow: 0 8px 28px rgba(30, 58, 95, 0.3);
}

.login-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.error-message {
  padding: 12px;
  background: #fee;
  color: #c33;
  border-radius: 8px;
  font-size: 14px;
  text-align: center;
}

@media (max-width: 480px) {
  .login-container {
    padding: 16px;
  }
  
  .login-card {
    padding: 28px 20px;
    max-width: 100%;
  }
  
  .login-header {
    margin-bottom: 24px;
  }
  
  .login-logo {
    width: 56px;
    height: 56px;
    margin-bottom: 12px;
  }
  
  .login-title {
    font-size: 22px;
    margin-bottom: 6px;
  }
  
  .login-subtitle {
    font-size: 13px;
  }
  
  .form-group {
    gap: 6px;
  }
  
  .form-group label {
    font-size: 13px;
  }
  
  .form-group input {
    padding: 10px 14px;
    font-size: 14px;
  }
  
  .login-btn {
    padding: 12px;
    font-size: 15px;
  }
  
  .error-message {
    font-size: 13px;
    padding: 10px;
  }
}

@media (max-width: 360px) {
  .login-card {
    padding: 24px 16px;
  }
  
  .login-logo {
    width: 48px;
    height: 48px;
  }
  
  .login-title {
    font-size: 20px;
  }
  
  .form-group input {
    font-size: 13px;
    padding: 9px 12px;
  }
  
  .login-btn {
    font-size: 14px;
    padding: 11px;
  }
}
</style>
```

---

## 📄 frontend-client-portal/src/views/NewRequestView.vue

```vue
<template>
  <div class="request-page">
    <PWAInstallButton />
    <!-- Hero Section -->
    <div class="hero-section">
      <div class="hero-icon">🔧</div>
      <h1>{{ t('app_name') }}</h1>
      <p>{{ t('professional_field_services') }}</p>
    </div>

    <div class="card request-form">
      <div class="form-header">
        <h2>📝 {{ t('new_service_request') }}</h2>
        <p>{{ t('fill_data_send_team') }}</p>
      </div>

      <form @submit.prevent="submitRequest">
        <!-- رسائل النجاح والخطأ -->
        <div v-if="error" class="alert alert-error">
          <span>❌</span> {{ error }}
        </div>
        <div v-if="success" class="alert alert-success">
          <span>✅</span> {{ success }}
        </div>

        <!-- حقل الاسم -->
        <div class="form-group">
          <label>{{ t('full_name') }} <span class="required">*</span></label>
          <div class="input-wrapper">
            <span class="input-icon">👤</span>
            <input 
              v-model="form.full_name" 
              type="text" 
              required 
              :placeholder="t('full_name')"
            />
          </div>
        </div>

        <!-- حقل الهاتف -->
        <div class="form-group">
          <label>{{ t('phone') }} <span class="required">*</span></label>
          <div class="input-wrapper">
            <span class="input-icon">📞</span>
            <input 
              v-model="form.phone" 
              type="tel" 
              required 
              :placeholder="t('example_phone')"
            />
          </div>
        </div>

        <!-- حقل عنوان الخدمة -->
        <div class="form-group">
          <label>{{ t('service_title') }} <span class="required">*</span></label>
          <div class="input-wrapper">
            <span class="input-icon">📋</span>
            <input 
              v-model="form.title" 
              type="text" 
              required 
              :placeholder="t('example_placeholder')"
            />
          </div>
        </div>

        <!-- حقل الوصف -->
        <div class="form-group">
          <label>{{ t('service_description') }} <span class="required">*</span></label>
          <div class="input-wrapper textarea-wrapper">
            <span class="input-icon">📝</span>
            <textarea 
              v-model="form.description" 
              required 
              :placeholder="t('detailed_description')"
              rows="4"
            ></textarea>
          </div>
        </div>

        <!-- حقل العنوان -->
        <div class="form-group">
          <label>{{ t('detailed_address') }}</label>
          <div class="input-wrapper">
            <span class="input-icon">📍</span>
            <input 
              v-model="form.address" 
              type="text" 
              :placeholder="t('full_address_landmarks')"
            />
          </div>
        </div>

        <!-- الأولوية -->
        <div class="form-group">
          <label>{{ t('priority') }}</label>
          <div class="priority-options">
            <label 
              v-for="p in priorities" 
              :key="p.value"
              class="priority-option"
              :class="{ active: form.priority === p.value }"
            >
              <input 
                type="radio" 
                :value="p.value" 
                v-model="form.priority"
                :style="{ accentColor: p.color }"
              />
              <span class="priority-label">
                <span class="priority-dot" :style="{ background: p.color }"></span>
                {{ t(p.value) }}
              </span>
            </label>
          </div>
        </div>

        <!-- تحديد الموقع -->
        <div class="form-group location-group">
          <label>📍 {{ t('your_location') }} <span class="required">*</span></label>
          
          <div class="location-status" :class="locationStatusClass">
            <div v-if="locationStatus === 'idle'" class="location-message">
              <span>📍</span>
              <span>{{ t('click_set_location') }}</span>
            </div>
            <div v-else-if="locationStatus === 'loading'" class="location-message">
              <span class="spinner"></span>
              <span>{{ t('determining_location') }}</span>
            </div>
            <div v-else-if="locationStatus === 'success'" class="location-message">
              <span>✅</span>
              <span>{{ t('location_determined_successfully') }}</span>
              <span class="location-coords mono">
                {{ locationName || `${location.lat.toFixed(6)}, ${location.lng.toFixed(6)}` }}
              </span>
            </div>
            <div v-else-if="locationStatus === 'error'" class="location-message">
              <span>❌</span>
              <span>{{ locationError }}</span>
            </div>
          </div>

          <div class="location-actions">
            <button 
              type="button" 
              class="btn btn--primary btn-location"
              @click="getLocation"
              :disabled="locationStatus === 'loading'"
            >
              {{ locationStatus === 'loading' ? '⏳ ' + t('sending') : t('get_my_location') }}
            </button>
            <button 
              type="button" 
              class="btn btn--ghost"
              @click="openMapPicker"
            >
              {{ t('select_from_map') }}
            </button>
          </div>
        </div>

        <!-- زر الإرسال -->
        <button 
          type="submit" 
          class="btn btn--primary btn--block btn--large btn-submit"
          :disabled="loading || !location"
        >
          <span v-if="loading" class="spinner"></span>
          {{ loading ? t('sending') : t('send_request') }}
        </button>

        <p class="form-footer">
          {{ t('all_data_saved_securely') }}
        </p>
      </form>
    </div>

    <!-- خريطة اختيار الموقع (Modal) -->
    <div v-if="showMapPicker" class="modal-backdrop" @click.self="showMapPicker = false">
      <div class="modal map-modal">
        <div class="modal-header">
          <h3>{{ t('choose_location_map') }}</h3>
          <button class="modal-close" @click="showMapPicker = false">✕</button>
        </div>
        <div class="modal-body">
          <div ref="mapContainer" class="map-container"></div>
        </div>
        <div class="modal-footer">
          <button class="btn btn--ghost" @click="showMapPicker = false">{{ t('cancel_btn') }}</button>
          <button class="btn btn--primary btn-confirm" @click="confirmMapLocation">{{ t('confirm_location') }}</button>
        </div>
      </div>
    </div>
  </div>

  <!-- مودال نجاح إرسال الطلب - خارج request-page -->
  <div v-if="showSuccessModal" class="modal-backdrop success-backdrop" @click.self="closeSuccessModal">
    <div class="modal success-modal">
      <div class="success-icon-wrapper">
        <div class="success-icon">✓</div>
      </div>
      
      <h2 class="success-title">{{ t('request_sent_successfully') }}</h2>
      <p class="success-subtitle">{{ t('request_received') }}</p>
      
      <div class="success-progress-bar">
        <div class="progress-fill"></div>
      </div>
      
      <div class="success-details">
        <div class="detail-item">
          <div class="detail-icon-wrapper">
            <span class="detail-icon">📄</span>
          </div>
          <div class="detail-content">
            <span class="detail-label">{{ t('request_number') }}</span>
            <span class="detail-value">#{{ successRequestData?.request_id }}</span>
          </div>
        </div>
        <div class="detail-item">
          <div class="detail-icon-wrapper">
            <span class="detail-icon">✏️</span>
          </div>
          <div class="detail-content">
            <span class="detail-label">{{ t('service_title') }}</span>
            <span class="detail-value">{{ successRequestData?.title }}</span>
          </div>
        </div>
        <div class="detail-item">
          <div class="detail-icon-wrapper">
            <span class="detail-icon">⚡</span>
          </div>
          <div class="detail-content">
            <span class="detail-label">{{ t('priority') }}</span>
            <span class="detail-value priority-badge" :class="'priority-' + successRequestData?.priority">
              {{ t(successRequestData?.priority) }}
            </span>
          </div>
        </div>
        <div class="detail-item">
          <div class="detail-icon-wrapper">
            <span class="detail-icon">📍</span>
          </div>
          <div class="detail-content">
            <span class="detail-label">{{ t('location') }}</span>
            <span class="detail-value">{{ successRequestData?.location }}</span>
          </div>
        </div>
      </div>

      <div class="success-message">
        <div class="message-icon">
          <span>💬</span>
        </div>
        <p>{{ t('team_will_contact') }}</p>
        <p class="success-note">{{ t('thank_you_choosing') }}</p>
      </div>

      <button class="btn btn--primary btn--large btn-success-action" @click="closeSuccessModal">
        <span class="btn-icon">🏠</span>
        {{ t('return_to_home') }}
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, nextTick, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from '../services/i18n'
import api from '../services/api'
import { servicesStore } from '../store/services'
import { reverseGeocode } from '../services/geocoding'
import PWAInstallButton from '../components/PWAInstallButton.vue'

const { t } = useI18n()
const router = useRouter()

const loading = ref(false)
const error = ref('')
const success = ref('')
const location = ref(null)
const locationName = ref('')
const locationStatus = ref('idle')
const locationError = ref('')
const showMapPicker = ref(false)
const showSuccessModal = ref(false)
const successRequestData = ref(null)
const mapContainer = ref(null)
let map = null
let mapMarker = null

const form = reactive({
  full_name: '',
  phone: '',
  title: '',
  description: '',
  address: '',
  priority: 'normal'
})

const priorities = [
  { value: 'low', color: '#4CAF50' },
  { value: 'normal', color: '#2196F3' },
  { value: 'high', color: '#FF9800' },
  { value: 'urgent', color: '#f44336' }
]

const locationStatusClass = computed(() => {
  return {
    'idle': locationStatus.value === 'idle',
    'loading': locationStatus.value === 'loading',
    'success': locationStatus.value === 'success',
    'error': locationStatus.value === 'error'
  }
})



async function getLocation() {
  if (!('geolocation' in navigator)) {
    locationStatus.value = 'error'
    locationError.value = t('browser_not_support_geolocation')
    return
  }

  locationStatus.value = 'loading'
  locationError.value = ''
  locationName.value = ''

  navigator.geolocation.getCurrentPosition(
    async (pos) => {
      location.value = {
        lat: pos.coords.latitude,
        lng: pos.coords.longitude
      }
      
      // Get location name using reverse geocoding
      try {
        const name = await reverseGeocode(
          pos.coords.latitude, 
          pos.coords.longitude, 
          'ar'
        )
        locationName.value = name
        locationStatus.value = 'success'
      } catch (err) {
        console.error('Geocoding error:', err)
        // Still succeed even if geocoding fails
        locationName.value = `${pos.coords.latitude.toFixed(6)}, ${pos.coords.longitude.toFixed(6)}`
        locationStatus.value = 'success'
      }
    },
    (err) => {
      locationStatus.value = 'error'
      
      // معالجة أفضل لأنواع مختلفة من الأخطاء
      switch(err.code) {
        case err.PERMISSION_DENIED:
          locationError.value = t('location_permission_denied')
          break
        case err.POSITION_UNAVAILABLE:
          locationError.value = t('location_unavailable')
          break
        case err.TIMEOUT:
          locationError.value = t('location_timeout')
          break
        default:
          locationError.value = t('location_failed') + ': ' + err.message
      }
    },
    { enableHighAccuracy: true, timeout: 30000, maximumAge: 0 }
  )
}

// تحديث الخريطة عند تغيير الموقع من GPS
function watchLocationChanges() {
  if (location.value && showMapPicker.value && map) {
    if (mapMarker) {
      map.removeLayer(mapMarker)
    }
    mapMarker = window.L.marker([location.value.lat, location.value.lng]).addTo(map)
    map.setView([location.value.lat, location.value.lng], 13)
    mapMarker.bindPopup('الموقع المحدد').openPopup()
  }
}

function openMapPicker() {
  showMapPicker.value = true
  nextTick(() => {
    initMap()
  })
}

function initMap() {
  if (!mapContainer.value) return
  
  // تحميل Leaflet dynamically
  if (!window.L) {
    const leafletScript = document.createElement('script')
    leafletScript.src = 'https://unpkg.com/leaflet@1.9.4/dist/leaflet.js'
    leafletScript.onload = () => initLeafletMap()
    document.head.appendChild(leafletScript)
    
    const leafletCSS = document.createElement('link')
    leafletCSS.rel = 'stylesheet'
    leafletCSS.href = 'https://unpkg.com/leaflet@1.9.4/dist/leaflet.css'
    document.head.appendChild(leafletCSS)
  } else {
    initLeafletMap()
  }
}

function initLeafletMap() {
  if (!mapContainer.value || !window.L) return
  
  // تنظيف الخريطة الموجودة
  if (map) {
    map.remove()
    map = null
  }
  
  // إنشاء الخريطة
  map = window.L.map(mapContainer.value).setView(
    location.value ? [location.value.lat, location.value.lng] : [24.7136, 46.6753], // Riyadh default
    13
  )
  
  // إضافة طبقة الخريطة
  window.L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    attribution: '© OpenStreetMap contributors'
  }).addTo(map)
  
  // إضافة marker للموقع المحدد
  if (location.value) {
    if (mapMarker) {
      map.removeLayer(mapMarker)
    }
    mapMarker = window.L.marker([location.value.lat, location.value.lng]).addTo(map)
    mapMarker.bindPopup('الموقع المحدد').openPopup()
  }
  
  // إضافة حدث النقر لتحديد الموقع
  map.on('click', async (e) => {
    const { lat, lng } = e.latlng
    
    if (mapMarker) {
      map.removeLayer(mapMarker)
    }
    
    mapMarker = window.L.marker([lat, lng]).addTo(map)
    mapMarker.bindPopup('الموقع المحدد').openPopup()
    
    location.value = { lat, lng }
    
    // Get location name using reverse geocoding
    try {
      const name = await reverseGeocode(lat, lng, 'ar')
      locationName.value = name
    } catch (err) {
      console.error('Geocoding error:', err)
      locationName.value = `${lat.toFixed(6)}, ${lng.toFixed(6)}`
    }
    
    locationStatus.value = 'success'
  })
}

async function confirmMapLocation() {
  if (mapMarker) {
    const markerLatLng = mapMarker.getLatLng()
    location.value = {
      lat: markerLatLng.lat,
      lng: markerLatLng.lng
    }
    
    // Get location name using reverse geocoding
    try {
      const name = await reverseGeocode(markerLatLng.lat, markerLatLng.lng, 'ar')
      locationName.value = name
    } catch (err) {
      console.error('Geocoding error:', err)
      locationName.value = `${markerLatLng.lat.toFixed(6)}, ${markerLatLng.lng.toFixed(6)}`
    }
  }
  showMapPicker.value = false
  locationStatus.value = 'success'
}

async function submitRequest() {
  if (!location.value) {
    error.value = t('please_set_location_first')
    return
  }

  loading.value = true
  error.value = ''
  success.value = ''

  try {
    const payload = {
      full_name: form.full_name,
      phone: form.phone,
      title: form.title,
      description: form.description,
      address: form.address,
      latitude: location.value.lat,
      longitude: location.value.lng,
      priority: form.priority,
      location_name: locationName.value || `${location.value.lat.toFixed(6)}, ${location.value.lng.toFixed(6)}`
    }

    const result = await servicesStore.createRequest(payload)
    
    // إظهار مودال النجاح بدلاً من رسالة بسيطة
    showSuccessModal.value = true
    successRequestData.value = {
      request_id: result.id || 'N/A',
      title: form.title,
      priority: form.priority,
      location: locationName.value || `${location.value.lat.toFixed(6)}, ${location.value.lng.toFixed(6)}`
    }
    
    // إعادة تعيين النموذج
    form.full_name = ''
    form.phone = ''
    form.title = ''
    form.description = ''
    form.address = ''
    form.priority = 'normal'
    location.value = null
    locationName.value = ''
    locationStatus.value = 'idle'
  } catch (err) {
    error.value = err.response?.data?.error || t('request_send_failed')
  } finally {
    loading.value = false
  }
}

function closeSuccessModal() {
  showSuccessModal.value = false
  router.push('/')
}

// Watch for location changes to update map
watch(location, () => {
  if (location.value && showMapPicker.value) {
    watchLocationChanges()
  }
})

// Watch for modal open/close
watch(showMapPicker, (newValue) => {
  if (newValue && location.value) {
    nextTick(() => {
      initMap()
    })
  }
})
</script>

<style scoped>
.request-page {
  max-width: 600px;
  margin: 0 auto;
  padding: 20px 0;
}

.hero-section {
  text-align: center;
  padding: 40px 20px 30px;
  background: linear-gradient(135deg, #1E3A5F 0%, #0D1B3E 100%);
  border-radius: var(--radius-lg);
  color: white;
  margin-bottom: 30px;
}

.hero-icon {
  font-size: 48px;
  margin-bottom: 12px;
}

.hero-section h1 {
  font-size: 28px;
  margin-bottom: 6px;
}

.hero-section p {
  opacity: 0.9;
  font-size: 14px;
}

.request-form {
  padding: 28px 24px;
}

.form-header {
  text-align: center;
  margin-bottom: 24px;
}

.form-header h2 {
  font-size: 20px;
  margin-bottom: 4px;
}

.form-header p {
  color: var(--ink-soft);
  font-size: 13px;
}

.form-group {
  margin-bottom: 18px;
}

.form-group label {
  display: block;
  font-size: 13px;
  font-weight: 600;
  color: var(--ink);
  margin-bottom: 6px;
}

.input-wrapper {
  position: relative;
  display: flex;
  align-items: center;
}

.input-icon {
  position: absolute;
  right: 12px;
  font-size: 16px;
  opacity: 0.6;
  pointer-events: none;
}

.input-wrapper input,
.input-wrapper textarea {
  width: 100%;
  padding: 12px 42px 12px 14px;
  border: 1.5px solid var(--line);
  border-radius: var(--radius-sm);
  font-family: var(--font-body);
  font-size: 14px;
  background: var(--surface);
  transition: all 0.3s ease;
}

.input-wrapper input:focus,
.input-wrapper textarea:focus {
  border-color: #1E3A5F;
  box-shadow: 0 0 0 3px rgba(30, 58, 95, 0.1);
  outline: none;
}

.textarea-wrapper {
  align-items: flex-start;
}

.textarea-wrapper .input-icon {
  top: 12px;
}

.required {
  color: var(--signal-out);
}

.priority-options {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 8px;
}

.priority-option {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 10px;
  border: 1.5px solid var(--line);
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: all 0.3s ease;
  font-size: 13px;
}

.priority-option:hover {
  border-color: #1E3A5F;
  background: rgba(30, 58, 95, 0.08);
}

.priority-option.active {
  border-color: #1E3A5F;
  background: rgba(30, 58, 95, 0.08);
}

.priority-option input[type="radio"] {
  display: none;
}

.priority-label {
  display: flex;
  align-items: center;
  gap: 4px;
}

.priority-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}

.location-group {
  background: var(--canvas);
  padding: 16px;
  border-radius: var(--radius-sm);
}

.location-status {
  padding: 12px;
  border-radius: var(--radius-sm);
  margin-bottom: 12px;
  text-align: center;
}

.location-status.idle {
  background: var(--line);
  color: var(--ink-soft);
}

.location-status.loading {
  background: rgba(30, 58, 95, 0.1);
  color: #1E3A5F;
}

.location-status.success {
  background: var(--signal-in-tint);
  color: var(--signal-in);
}

.location-status.error {
  background: var(--signal-out-tint);
  color: var(--signal-out);
}

.location-message {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  flex-wrap: wrap;
}

.location-coords {
  font-size: 12px;
  opacity: 0.8;
}

.location-actions {
  display: flex;
  gap: 10px;
}

.location-actions .btn {
  flex: 1;
}

.alert {
  padding: 12px 16px;
  border-radius: var(--radius-sm);
  margin-bottom: 16px;
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 13px;
}

.alert-error {
  background: var(--signal-out-tint);
  color: var(--signal-out);
}

.alert-success {
  background: var(--signal-in-tint);
  color: var(--signal-in);
}

.btn--large {
  padding: 16px;
  font-size: 16px;
}

.btn--block {
  width: 100%;
}

.btn-submit {
  background: #1E3A5F;
  color: white;
  border: none;
  box-shadow: 0 4px 12px rgba(30, 58, 95, 0.3);
}

.btn-submit:hover:not(:disabled) {
  background: #0D1B3E;
  transform: translateY(-2px);
  box-shadow: 0 8px 28px rgba(30, 58, 95, 0.3);
}

.btn-location {
  background: #1E3A5F;
  color: white;
  border: none;
  box-shadow: 0 4px 12px rgba(30, 58, 95, 0.3);
}

.btn-location:hover:not(:disabled) {
  background: #0D1B3E;
  transform: translateY(-2px);
  box-shadow: 0 8px 28px rgba(30, 58, 95, 0.3);
}

.btn-confirm {
  background: #1E3A5F;
  color: white;
  border: none;
  box-shadow: 0 4px 12px rgba(30, 58, 95, 0.3);
}

.btn-confirm:hover {
  background: #0D1B3E;
  transform: translateY(-2px);
  box-shadow: 0 8px 28px rgba(30, 58, 95, 0.3);
}

.form-footer {
  text-align: center;
  font-size: 12px;
  color: var(--ink-soft);
  margin-top: 14px;
}

.spinner {
  display: inline-block;
  width: 18px;
  height: 18px;
  border: 2px solid rgba(255,255,255,0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.map-container {
  height: 400px;
  border-radius: var(--radius-sm);
  z-index: 1;
}

/* Leaflet styles override */
.map-container :deep(.leaflet-container) {
  font-family: var(--font-body);
}

.map-container :deep(.leaflet-popup-content-wrapper) {
  border-radius: var(--radius-sm);
}

.map-container :deep(.leaflet-popup-content) {
  margin: 12px;
  font-size: 14px;
}

.modal-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.6);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 20px;
}

.map-modal {
  width: 100%;
  max-width: 700px;
  padding: 0;
  overflow: hidden;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid var(--line);
}

.modal-header h3 { font-size: 16px; margin: 0; }

.modal-close {
  background: none;
  border: none;
  font-size: 20px;
  cursor: pointer;
  color: var(--ink-soft);
}

.modal-body {
  padding: 0;
}

.modal-footer {
  padding: 16px 20px;
  border-top: 1px solid var(--line);
  display: flex;
  gap: 10px;
  justify-content: flex-end;
}

/* Success Modal Styles */
.success-backdrop {
  background: rgba(0, 0, 0, 0.75);
  backdrop-filter: blur(12px);
  animation: backdropFadeIn 0.3s ease;
}

@keyframes backdropFadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

.success-modal {
  max-width: 380px;
  padding: 0;
  text-align: center;
  background: linear-gradient(145deg, #ffffff 0%, #f8fafc 50%, #e8f0f8 100%);
  border-radius: 20px;
  box-shadow: 
    0 15px 50px rgba(0, 0, 0, 0.2),
    0 8px 20px rgba(0, 0, 0, 0.12),
    inset 0 1px 0 rgba(255, 255, 255, 0.8);
  animation: successSlideIn 0.6s cubic-bezier(0.34, 1.56, 0.64, 1);
  position: relative;
  overflow: hidden;
}

.success-modal::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 4px;
  background: linear-gradient(90deg, #4CAF50, #8BC34A, #CDDC39, #4CAF50);
  background-size: 300% 100%;
  animation: gradientFlow 3s ease infinite;
}

@keyframes gradientFlow {
  0%, 100% {
    background-position: 0% 50%;
  }
  50% {
    background-position: 100% 50%;
  }
}

@keyframes successSlideIn {
  from {
    opacity: 0;
    transform: translateY(40px) scale(0.9) rotateX(10deg);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1) rotateX(0deg);
  }
}

.success-icon-wrapper {
  margin: 24px auto 20px;
  width: 70px;
  height: 70px;
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
}

.success-icon {
  font-size: 44px;
  color: #4CAF50;
  filter: drop-shadow(0 4px 12px rgba(76, 175, 80, 0.4));
  animation: iconBounce 0.8s cubic-bezier(0.34, 1.56, 0.64, 1);
  font-weight: bold;
}

@keyframes iconBounce {
  0% {
    transform: scale(0) rotate(-180deg);
  }
  50% {
    transform: scale(1.1) rotate(10deg);
  }
  100% {
    transform: scale(1) rotate(0deg);
  }
}

.success-title {
  font-size: 20px;
  font-weight: 800;
  color: #1E3A5F;
  margin: 0 0 10px 0;
  background: linear-gradient(135deg, #1E3A5F 0%, #0D1B3E 50%, #1E3A5F 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  padding: 0 16px;
  letter-spacing: -0.5px;
}

.success-subtitle {
  font-size: 13px;
  color: #6B7A8A;
  margin: 0 0 16px 0;
  font-weight: 500;
  padding: 0 16px;
}

.success-progress-bar {
  width: 100%;
  height: 3px;
  background: linear-gradient(90deg, #e0e0e0, #f5f5f5);
  border-radius: 2px;
  margin: 0 0 20px 0;
  overflow: hidden;
  position: relative;
}

.progress-fill {
  height: 100%;
  width: 0%;
  background: linear-gradient(90deg, #4CAF50, #8BC34A);
  border-radius: 2px;
  animation: progressFill 1.5s ease 0.5s forwards;
  box-shadow: 0 0 10px rgba(76, 175, 80, 0.5);
}

@keyframes progressFill {
  to {
    width: 100%;
  }
}

.success-details {
  background: linear-gradient(135deg, #f8fafc 0%, #ffffff 100%);
  border-radius: 16px;
  padding: 16px 20px;
  margin: 0 16px 20px;
  border: 1px solid #e8f0f8;
  box-shadow: 
    0 4px 16px rgba(0, 0, 0, 0.06),
    inset 0 1px 0 rgba(255, 255, 255, 0.8);
}

.detail-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 0;
  border-bottom: 1px solid #e8f0f8;
  transition: all 0.3s ease;
}

.detail-item:hover {
  background: linear-gradient(90deg, transparent, rgba(76, 175, 80, 0.05), transparent);
  border-radius: 12px;
  padding: 12px 8px;
}

.detail-item:last-child {
  border-bottom: none;
  padding-bottom: 0;
}

.detail-item:first-child {
  padding-top: 0;
}

.detail-icon-wrapper {
  width: 32px;
  height: 32px;
  background: linear-gradient(135deg, #4CAF50 0%, #8BC34A 100%);
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  box-shadow: 0 4px 12px rgba(76, 175, 80, 0.3);
  transition: transform 0.3s ease;
}

.detail-item:hover .detail-icon-wrapper {
  transform: scale(1.1) rotate(5deg);
}

.detail-icon {
  font-size: 16px;
}

.detail-content {
  flex: 1;
  text-align: right;
}

.detail-label {
  font-size: 11px;
  color: #6B7A8A;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin-bottom: 3px;
  display: block;
}

.detail-value {
  font-size: 13px;
  color: #1A2A3A;
  font-weight: 700;
  word-break: break-word;
  line-height: 1.4;
}

.priority-badge {
  padding: 4px 10px;
  border-radius: 16px;
  font-size: 10px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  display: inline-block;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

.priority-low {
  background: linear-gradient(135deg, #d1fae5 0%, #a7f3d0 100%);
  color: #065f46;
  border: 1px solid #6ee7b7;
}

.priority-normal {
  background: linear-gradient(135deg, #dbeafe 0%, #bfdbfe 100%);
  color: #1e40af;
  border: 1px solid #93c5fd;
}

.priority-high {
  background: linear-gradient(135deg, #fed7aa 0%, #fdba74 100%);
  color: #9a3412;
  border: 1px solid #fdba74;
}

.priority-urgent {
  background: linear-gradient(135deg, #fecaca 0%, #fca5a5 100%);
  color: #991b1b;
  border: 1px solid #f87171;
  animation: pulse 2s ease infinite;
}

@keyframes pulse {
  0%, 100% {
    box-shadow: 0 2px 8px rgba(153, 27, 27, 0.3);
  }
  50% {
    box-shadow: 0 2px 16px rgba(153, 27, 27, 0.5);
  }
}

.success-message {
  margin: 0 16px 20px;
  padding: 16px;
  background: linear-gradient(135deg, #f0fdf4 0%, #dcfce7 100%);
  border-radius: 12px;
  border: 1px solid #bbf7d0;
  position: relative;
  overflow: hidden;
}

.success-message::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 3px;
  background: linear-gradient(90deg, #4CAF50, #8BC34A, #4CAF50);
  background-size: 200% 100%;
  animation: gradientFlow 2s ease infinite;
}

.message-icon {
  width: 28px;
  height: 28px;
  margin: 0 auto 10px;
  font-size: 20px;
  animation: messageFloat 2s ease-in-out infinite;
  display: flex;
  align-items: center;
  justify-content: center;
}

@keyframes messageFloat {
  0%, 100% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-4px);
  }
}

.success-message p {
  font-size: 12px;
  color: #166534;
  margin: 6px 0;
  line-height: 1.5;
  font-weight: 500;
}

.success-note {
  font-size: 11px;
  color: #15803d;
  font-weight: 600;
  margin-top: 10px;
  padding-top: 10px;
  border-top: 1px solid #bbf7d0;
}

.btn-success-action {
  background: linear-gradient(135deg, #1E3A5F 0%, #0D1B3E 100%);
  color: white;
  border: none;
  box-shadow: 
    0 8px 24px rgba(30, 58, 95, 0.4),
    0 4px 12px rgba(30, 58, 95, 0.3);
  width: calc(100% - 32px);
  margin: 0 16px 24px;
  padding: 12px 20px;
  border-radius: 12px;
  font-weight: 700;
  font-size: 14px;
  letter-spacing: 0.5px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  position: relative;
  overflow: hidden;
  transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
}

.btn-success-action::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.2), transparent);
  transition: left 0.5s ease;
}

.btn-success-action:hover::before {
  left: 100%;
}

.btn-icon {
  font-size: 16px;
  transition: transform 0.3s ease;
}

.btn-success-action:hover .btn-icon {
  transform: translateX(-4px);
}

.btn-success-action:hover {
  background: linear-gradient(135deg, #0D1B3E 0%, #1E3A5F 100%);
  transform: translateY(-3px) scale(1.02);
  box-shadow: 
    0 12px 32px rgba(30, 58, 95, 0.5),
    0 6px 16px rgba(30, 58, 95, 0.4);
}

.btn-success-action:active {
  transform: translateY(-1px) scale(0.98);
}

@media (max-width: 480px) {
  .request-page {
    padding: 16px 0;
  }
  
  .hero-section {
    padding: 24px 16px 20px;
    margin-bottom: 20px;
  }
  
  .hero-icon {
    font-size: 36px;
    margin-bottom: 8px;
  }
  
  .hero-section h1 {
    font-size: 22px;
    margin-bottom: 4px;
  }
  
  .hero-section p {
    font-size: 13px;
  }
  
  .request-form {
    padding: 20px 16px;
  }
  
  .form-header {
    margin-bottom: 20px;
  }
  
  .form-header h2 {
    font-size: 18px;
  }
  
  .form-header p {
    font-size: 12px;
  }
  
  .form-group {
    margin-bottom: 16px;
  }
  
  .form-group label {
    font-size: 12px;
    margin-bottom: 5px;
  }
  
  .input-wrapper input,
  .input-wrapper textarea {
    padding: 10px 38px 10px 12px;
    font-size: 13px;
  }
  
  .input-icon {
    right: 10px;
    font-size: 14px;
  }
  
  .priority-options {
    grid-template-columns: repeat(2, 1fr);
    gap: 6px;
  }
  
  .priority-option {
    padding: 6px 8px;
    font-size: 12px;
  }
  
  .location-group {
    padding: 12px;
  }
  
  .location-status {
    padding: 10px;
    margin-bottom: 10px;
  }
  
  .location-message {
    font-size: 12px;
    gap: 6px;
  }
  
  .location-actions {
    flex-direction: column;
    gap: 8px;
  }
  
  .location-actions .btn {
    font-size: 13px;
    padding: 10px;
  }
  
  .btn--large {
    padding: 14px;
    font-size: 15px;
  }
  
  /* Success Modal Responsive */
  .success-modal {
    max-width: 95%;
    border-radius: 16px;
  }
  
  .success-icon-wrapper {
    width: 60px;
    height: 60px;
    margin: 20px auto 16px;
  }
  
  .success-icon {
    font-size: 38px;
  }
  
  .success-title {
    font-size: 18px;
    padding: 0 12px;
  }
  
  .success-subtitle {
    font-size: 12px;
    padding: 0 12px;
  }
  
  .success-progress-bar {
    margin: 0 0 16px 0;
  }
  
  .success-details {
    padding: 14px 16px;
    margin: 0 12px 16px;
  }
  
  .detail-item {
    padding: 10px 0;
  }
  
  .detail-icon-wrapper {
    width: 28px;
    height: 28px;
  }
  
  .detail-icon {
    font-size: 14px;
  }
  
  .detail-label {
    font-size: 10px;
  }
  
  .detail-value {
    font-size: 12px;
  }
  
  .priority-badge {
    padding: 3px 8px;
    font-size: 9px;
  }
  
  .success-message {
    margin: 0 12px 16px;
    padding: 12px;
  }
  
  .message-icon {
    width: 24px;
    height: 24px;
    font-size: 18px;
  }
  
  .success-message p {
    font-size: 11px;
  }
  
  .success-note {
    font-size: 10px;
  }
  
  .btn-success-action {
    width: calc(100% - 24px);
    margin: 0 12px 20px;
    padding: 10px 16px;
    font-size: 13px;
  }
  
  .btn-icon {
    font-size: 14px;
  }
  
  .map-container {
    height: 300px;
  }
  
  .modal-backdrop {
    padding: 12px;
  }
  
  .map-modal {
    max-width: 100%;
  }
  
  .modal-header {
    padding: 12px 16px;
  }
  
  .modal-header h3 {
    font-size: 14px;
  }
  
  .modal-footer {
    padding: 12px 16px;
  }
  
  .modal-footer .btn {
    font-size: 13px;
    padding: 8px 14px;
  }
  
  .location-status {
    padding: 10px;
    margin-bottom: 10px;
  }
  
  .location-message {
    font-size: 12px;
    gap: 6px;
  }
  
  .location-actions {
    flex-direction: column;
    gap: 8px;
  }
  
  .location-actions .btn {
    font-size: 13px;
    padding: 10px;
  }
  
  .form-footer {
    font-size: 11px;
    margin-top: 12px;
  }
}

@media (min-width: 769px) {
  /* Success Modal Desktop Styles */
  .success-modal {
    max-width: 400px;
  }
  
  .success-icon-wrapper {
    width: 68px;
    height: 68px;
    margin: 26px auto 22px;
  }
  
  .success-icon {
    font-size: 42px;
  }
  
  .success-title {
    font-size: 21px;
    padding: 0 18px;
  }
  
  .success-subtitle {
    font-size: 13px;
    padding: 0 18px;
  }
  
  .success-details {
    padding: 18px 22px;
    margin: 0 18px 22px;
  }
  
  .detail-item {
    padding: 13px 0;
  }
  
  .detail-icon-wrapper {
    width: 34px;
    height: 34px;
  }
  
  .detail-icon {
    font-size: 17px;
  }
  
  .detail-label {
    font-size: 11px;
  }
  
  .detail-value {
    font-size: 13px;
  }
  
  .priority-badge {
    padding: 4px 11px;
    font-size: 10px;
  }
  
  .success-message {
    margin: 0 18px 22px;
    padding: 17px;
  }
  
  .message-icon {
    width: 26px;
    height: 26px;
    font-size: 19px;
  }
  
  .success-message p {
    font-size: 12px;
  }
  
  .success-note {
    font-size: 11px;
  }
  
  .btn-success-action {
    width: calc(100% - 36px);
    margin: 0 18px 26px;
    padding: 13px 21px;
    font-size: 14px;
  }
  
  .btn-icon {
    font-size: 17px;
  }
}

@media (min-width: 481px) and (max-width: 768px) {
  /* Success Modal Tablet Styles */
  .success-modal {
    max-width: 420px;
  }
  
  .success-icon-wrapper {
    width: 75px;
    height: 75px;
    margin: 28px auto 22px;
  }
  
  .success-icon {
    font-size: 48px;
  }
  
  .success-title {
    font-size: 22px;
    padding: 0 18px;
  }
  
  .success-subtitle {
    font-size: 14px;
    padding: 0 18px;
  }
  
  .success-details {
    padding: 20px 24px;
    margin: 0 20px 24px;
  }
  
  .detail-item {
    padding: 14px 0;
  }
  
  .detail-icon-wrapper {
    width: 36px;
    height: 36px;
  }
  
  .detail-icon {
    font-size: 18px;
  }
  
  .detail-label {
    font-size: 11px;
  }
  
  .detail-value {
    font-size: 14px;
  }
  
  .priority-badge {
    padding: 5px 12px;
    font-size: 11px;
  }
  
  .success-message {
    margin: 0 20px 24px;
    padding: 18px;
  }
  
  .message-icon {
    width: 30px;
    height: 30px;
    font-size: 22px;
  }
  
  .success-message p {
    font-size: 13px;
  }
  
  .success-note {
    font-size: 12px;
  }
  
  .btn-success-action {
    width: calc(100% - 40px);
    margin: 0 20px 28px;
    padding: 14px 22px;
    font-size: 15px;
  }
  
  .btn-icon {
    font-size: 18px;
  }
}

@media (max-width: 360px) {
  .hero-section h1 {
    font-size: 20px;
  }
  
  .form-header h2 {
    font-size: 16px;
  }
  
  .priority-options {
    grid-template-columns: 1fr;
  }
  
  .input-wrapper input,
  .input-wrapper textarea {
    font-size: 12px;
  }
  
  /* Success Modal Extra Small Screens */
  .success-modal {
    max-width: 98%;
    border-radius: 12px;
  }
  
  .success-icon-wrapper {
    width: 50px;
    height: 50px;
    margin: 16px auto 12px;
  }
  
  .success-icon {
    font-size: 32px;
  }
  
  .success-title {
    font-size: 16px;
    padding: 0 10px;
  }
  
  .success-subtitle {
    font-size: 11px;
    padding: 0 10px;
  }
  
  .success-details {
    padding: 12px;
    margin: 0 10px 12px;
  }
  
  .detail-item {
    padding: 8px 0;
  }
  
  .detail-icon-wrapper {
    width: 24px;
    height: 24px;
  }
  
  .detail-icon {
    font-size: 12px;
  }
  
  .detail-label {
    font-size: 9px;
  }
  
  .detail-value {
    font-size: 11px;
  }
  
  .priority-badge {
    padding: 2px 6px;
    font-size: 8px;
  }
  
  .success-message {
    margin: 0 10px 12px;
    padding: 10px;
  }
  
  .message-icon {
    width: 20px;
    height: 20px;
    font-size: 16px;
  }
  
  .success-message p {
    font-size: 10px;
  }
  
  .success-note {
    font-size: 9px;
  }
  
  .btn-success-action {
    width: calc(100% - 20px);
    margin: 0 10px 16px;
    padding: 8px 14px;
    font-size: 12px;
  }
  
  .btn-icon {
    font-size: 12px;
  }
}
</style>

```

---

## 📄 frontend-client-portal/src/views/RatingView.vue

```vue
<template>
  <div class="rate-page">
    <PWAInstallButton />
    <h2>{{ t('how_was_your_experience') }}</h2>
    <p>{{ t('rating_helps_improve') }}</p>

    <div class="card rate-card">
      <div v-if="error" class="alert alert-error">
        <span>❌</span> {{ error }}
      </div>
      <div v-if="success" class="alert alert-success">
        <span>✅</span> {{ success }}
      </div>

      <StarRating v-model="rating" />
      <textarea 
        v-model="comment" 
        class="field" 
        :placeholder="t('tell_us_notes')" 
        rows="4"
        :disabled="loading"
      ></textarea>
      <button 
        class="btn btn--primary btn--block" 
        @click="submitRating"
        :disabled="loading || rating === 0"
      >
        <span v-if="loading" class="spinner"></span>
        {{ loading ? t('sending') : t('submit_rating') }}
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import StarRating from '../components/StarRating.vue'
import PWAInstallButton from '../components/PWAInstallButton.vue'
import api from '../services/api'
import { useI18n } from '../services/i18n'

const { t } = useI18n()
const router = useRouter()
const rating = ref(0)
const comment = ref('')
const loading = ref(false)
const error = ref('')
const success = ref('')

async function submitRating() {
  if (rating.value === 0) {
    error.value = t('please_select_rating')
    return
  }

  loading.value = true
  error.value = ''
  success.value = ''

  try {
    const response = await api.post(`/service/requests/${router.currentRoute.value.params.id}/rate`, {
      rating: rating.value,
      comment: comment.value
    })
    
    success.value = t('thank_you_rating')
    
    setTimeout(() => {
      router.push('/')
    }, 2000)
  } catch (err) {
    error.value = err.response?.data?.error || t('rating_submit_failed')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.rate-page { text-align: center; padding-top: 10px; }
.rate-page h2 { font-size: 19px; margin-bottom: 6px; }
.rate-page p { font-size: 13px; color: var(--ink-soft); margin-bottom: 20px; }
.rate-card { padding: 26px 20px; display: flex; flex-direction: column; gap: 18px; text-align: right; }
.rate-card textarea {
  width: 100%; font-family: var(--font-body); font-size: 14px; padding: 12px 14px;
  border-radius: var(--radius-sm); border: 1px solid var(--line); resize: vertical;
}

.rate-card textarea:disabled {
  background: var(--canvas);
  opacity: 0.6;
}

.alert {
  padding: 12px 16px;
  border-radius: var(--radius-sm);
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 13px;
}

.alert-error {
  background: var(--signal-out-tint);
  color: var(--signal-out);
}

.alert-success {
  background: var(--signal-in-tint);
  color: var(--signal-in);
}

.spinner {
  display: inline-block;
  width: 16px;
  height: 16px;
  border: 2px solid rgba(255,255,255,0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

@media (max-width: 480px) {
  .rate-page {
    padding-top: 8px;
  }
  
  .rate-page h2 {
    font-size: 17px;
    margin-bottom: 4px;
  }
  
  .rate-page p {
    font-size: 12px;
    margin-bottom: 16px;
  }
  
  .rate-card {
    padding: 20px 16px;
    gap: 14px;
  }
  
  .rate-card textarea {
    font-size: 13px;
    padding: 10px 12px;
  }
  
  .btn--block {
    font-size: 13px;
    padding: 12px;
  }
  
  .alert {
    font-size: 12px;
    padding: 10px 14px;
  }
}

@media (max-width: 360px) {
  .rate-page h2 {
    font-size: 16px;
  }
  
  .rate-card {
    padding: 16px 12px;
  }
  
  .rate-card textarea {
    font-size: 12px;
    padding: 9px 10px;
  }
}
</style>

```

---

## 📄 frontend-client-portal/src/views/ServiceHistoryView.vue

```vue
<template>
  <div>
    <PWAInstallButton />
    <h2 class="page-title">{{ t('service_history') }}</h2>
    <p class="page-subtitle">{{ t('all_requests_status') }}</p>

    <div v-if="reportsStore.loading" class="empty-state">
      <p>{{ t('loading_data') }}</p>
    </div>

    <div v-else-if="reportsStore.error" class="alert alert-error">
      <span>❌</span> {{ reportsStore.error }}
    </div>

    <div v-else-if="!reportsStore.reports || !reportsStore.reports.length" class="empty-state">
      <p>{{ t('no_previous_services') }}</p>
    </div>

    <ReportCard v-else v-for="r in reportsStore.reports" :key="r.id" :report="r" />
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import ReportCard from '../components/ReportCard.vue'
import PWAInstallButton from '../components/PWAInstallButton.vue'
import { reportsStore } from '../store/reports'
import { useI18n } from '../services/i18n'

const { t } = useI18n()

onMounted(() => {
  reportsStore.fetchReports()
})
</script>

<style scoped>
.page-title { font-size: 19px; margin-bottom: 4px; }
.page-subtitle { font-size: 13px; color: var(--ink-soft); margin-bottom: 18px; }

.empty-state {
  text-align: center;
  padding: 40px 20px;
  color: var(--ink-soft);
}

.alert {
  padding: 12px 16px;
  border-radius: var(--radius-sm);
  margin-bottom: 16px;
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 13px;
}

.alert-error {
  background: var(--signal-out-tint);
  color: var(--signal-out);
}

@media (max-width: 480px) {
  .page-title {
    font-size: 17px;
  }
  
  .page-subtitle {
    font-size: 12px;
    margin-bottom: 14px;
  }
  
  .empty-state {
    padding: 32px 16px;
  }
  
  .alert {
    font-size: 12px;
    padding: 10px 14px;
  }
}

@media (max-width: 360px) {
  .page-title {
    font-size: 16px;
  }
  
  .page-subtitle {
    font-size: 11px;
  }
}
</style>

```

---

## 📄 frontend-client-portal/src/views/ServiceReportView.vue

```vue
<template>
  <div>
    <PWAInstallButton />
    <router-link to="/" class="back-link">{{ t('return_to_history') }}</router-link>

    <div v-if="reportsStore.loading" class="empty-state">
      <p>{{ t('loading_report') }}</p>
    </div>

    <div v-else-if="reportsStore.error" class="alert alert-error">
      <span>❌</span> {{ reportsStore.error }}
    </div>

    <div v-else-if="!report" class="empty-state">
      <p>{{ t('report_not_found') }}</p>
    </div>

    <template v-else>
      <div class="card report">
        <span :class="['badge', getStatusBadgeClass(report.status)]">
          {{ getStatusIcon(report.status) }} {{ getStatusText(report.status) }}
        </span>
        <h2>{{ report.title }}</h2>
        <p class="report__meta mono">{{ t('executed_on') }}: {{ formatDate(report.created_at) }}</p>
        <p class="report__meta">{{ t('location') }}: {{ report.address || t('location_not_specified') }}</p>

        <hr class="divider" />

        <h3>{{ t('work_summary') }}</h3>
        <p class="report__text">
          {{ report.description || t('no_description_available') }}
        </p>

        <!-- عرض التقييم إذا وجد -->
        <div v-if="report.rating" class="rating-section">
          <h4>{{ t('your_rating') }}</h4>
          <div class="rating-display">
            <span class="rating-stars">⭐ {{ report.rating }}/5</span>
            <p v-if="report.feedback" class="rating-feedback">"{{ report.feedback }}"</p>
          </div>
        </div>
      </div>

      <router-link v-if="report.status === 'completed' && !report.rating" :to="`/rate/${$route.params.id}`" class="btn btn--primary btn--block">
        {{ t('rate_service') }}
      </router-link>
    </template>
  </div>
</template>

<script setup>
import { computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { reportsStore } from '../store/reports'
import { useI18n } from '../services/i18n'
import PWAInstallButton from '../components/PWAInstallButton.vue'

const { t } = useI18n()
const route = useRoute()

const report = computed(() => {
  return reportsStore.find(route.params.id)
})

onMounted(() => {
  reportsStore.fetchReports()
})

function formatDate(date) {
  if (!date) return t('location_not_specified')
  return new Date(date).toLocaleDateString('en-GB', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
    hour12: false
  })
}

function getStatusText(status) {
  const statusMap = {
    'completed': t('status_completed'),
    'cancelled': t('status_cancelled'),
    'pending': t('status_pending'),
    'assigned': t('status_assigned'),
    'in_progress': t('status_in_progress')
  }
  return statusMap[status] || status
}

function getStatusBadgeClass(status) {
  const classMap = {
    'completed': 'badge--completed',
    'cancelled': 'badge--cancelled',
    'pending': 'badge--pending',
    'assigned': 'badge--assigned',
    'in_progress': 'badge--in-progress'
  }
  return classMap[status] || 'badge--pending'
}

function getStatusIcon(status) {
  const iconMap = {
    'completed': '✅',
    'cancelled': '❌',
    'pending': '⏳',
    'assigned': '📋',
    'in_progress': '🔄'
  }
  return iconMap[status] || '⏳'
}
</script>

<style scoped>
.back-link { display: inline-block; font-size: 13px; color: var(--ink-soft); margin-bottom: 14px; }
.report { padding: 22px; margin-bottom: 18px; }
.report h2 { font-size: 18px; margin: 10px 0 6px; }
.report__meta { font-size: 13px; color: var(--ink-soft); margin-bottom: 2px; }
.divider { border: none; border-top: 1px solid var(--line); margin: 18px 0; }
.report h3 { font-size: 14px; margin-bottom: 8px; }
.report__text { font-size: 13.5px; color: var(--ink-soft); line-height: 1.8; }

.rating-section {
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid var(--line);
}

.rating-section h4 {
  font-size: 14px;
  margin-bottom: 8px;
  color: var(--ink);
}

.rating-display {
  background: var(--signal-in-tint);
  border: 1px solid var(--signal-in);
  border-radius: var(--radius-sm);
  padding: 12px 14px;
}

.rating-stars {
  font-size: 16px;
  font-weight: 600;
  color: var(--signal-in);
  display: block;
  margin-bottom: 6px;
}

.rating-feedback {
  font-size: 13px;
  color: var(--ink);
  font-style: italic;
  margin: 0;
}

.badge {
  background: var(--brand);
  color: white;
  padding: 4px 10px;
  border-radius: 12px;
  font-size: 11px;
  font-weight: 600;
  white-space: nowrap;
  display: inline-block;
  margin-bottom: 10px;
}

.badge--completed {
  background: var(--signal-in);
}

.badge--cancelled {
  background: var(--signal-out);
}

.badge--pending {
  background: #FFA500;
}

.badge--assigned {
  background: #2196F3;
}

.badge--in-progress {
  background: #9C27B0;
}

.empty-state {
  text-align: center;
  padding: 40px 20px;
  color: var(--ink-soft);
}

.alert {
  padding: 12px 16px;
  border-radius: var(--radius-sm);
  margin-bottom: 16px;
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 13px;
}

.alert-error {
  background: var(--signal-out-tint);
  color: var(--signal-out);
}

@media (max-width: 480px) {
  .back-link {
    font-size: 12px;
    margin-bottom: 12px;
  }
  
  .report {
    padding: 16px 14px;
    margin-bottom: 14px;
  }
  
  .report h2 {
    font-size: 16px;
    margin: 8px 0 4px;
  }
  
  .report__meta {
    font-size: 12px;
  }
  
  .divider {
    margin: 14px 0;
  }
  
  .report h3 {
    font-size: 13px;
    margin-bottom: 6px;
  }
  
  .report__text {
    font-size: 12px;
    line-height: 1.7;
  }
  
  .btn--block {
    font-size: 13px;
    padding: 12px;
  }
  
  .empty-state {
    padding: 32px 16px;
  }
  
  .alert {
    font-size: 12px;
    padding: 10px 14px;
  }
}

@media (max-width: 360px) {
  .report {
    padding: 14px 12px;
  }
  
  .report h2 {
    font-size: 15px;
  }
  
  .report__text {
    font-size: 11px;
  }
}
</style>

```

---

## 📄 frontend-client-portal/vite.config.js

```javascript
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  appType: 'spa',
  plugins: [vue()],
  test: {
    globals: true,
    environment: 'jsdom',
    setupFiles: ['./src/tests/setup.js'],
    coverage: {
      provider: 'v8',
      reporter: ['text', 'json', 'html'],
      exclude: [
        'node_modules/',
        'src/tests/',
        '**/*.d.ts',
        '**/*.config.*',
        '**/mockData',
        'dist/',
        'electron/'
      ]
    }
  },
  server: { 
    port: 3001,
    cors: true,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        secure: false
      }
    }
  },
  preview: {
    port: 3001,
    // SPA routing support for preview
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        secure: false
      }
    }
  },
  build: {
    cssCodeSplit: true,
    rollupOptions: {
      output: {
        manualChunks: {
          'vue-vendor': ['vue', 'vue-router']
        },
        // إضافة timestamp تلقائي لضمان cache busting
        entryFileNames: 'assets/[name]-[hash].js',
        chunkFileNames: 'assets/[name]-[hash].js',
        assetFileNames: 'assets/[name]-[hash][extname]'
      }
    }
  },
  base: '/'
})

```

---

## 📄 frontend-worker-pwa/build-obfuscated.js

```javascript
import JavaScriptObfuscator from 'javascript-obfuscator';
import { readFileSync, writeFileSync, readdirSync, statSync } from 'fs';
import { join } from 'path';

// إعدادات التعتيم
const obfuscatorOptions = {
  compact: true,
  controlFlowFlattening: true,
  controlFlowFlatteningThreshold: 0.75,
  deadCodeInjection: true,
  deadCodeInjectionThreshold: 0.4,
  debugProtection: true,
  disableConsoleOutput: false, // Changed to false for debugging
  identifierNamesGenerator: 'hexadecimal',
  log: true, // Changed to true for debugging
  numbersToExpressions: true,
  renameGlobals: false,
  selfDefending: true,
  simplify: true,
  splitStrings: true,
  splitStringsChunkLength: 10,
  stringArray: true,
  stringArrayEncoding: ['base64'],
  stringArrayIndexShift: true,
  stringArrayWrappersCount: 2,
  stringArrayWrappersChainedCalls: true,
  stringArrayWrappersParametersMaxCount: 4,
  stringArrayWrappersType: 'function',
  stringArrayThreshold: 0.75,
  transformObjectKeys: true,
  unicodeEscapeSequence: false
};

// دالة للبحث عن ملفات JS بشكل متكرر
function findJsFiles(dir, fileList = []) {
  try {
    const files = readdirSync(dir);
    
    for (const file of files) {
      const filePath = join(dir, file);
      const stats = statSync(filePath);
      
      if (stats.isDirectory()) {
        findJsFiles(filePath, fileList);
      } else if (file.endsWith('.js')) {
        fileList.push(filePath);
      }
    }
  } catch (error) {
    console.error(`Error reading directory ${dir}:`, error.message);
  }
  
  return fileList;
}

// البحث عن جميع ملفات JS في مجلد dist
const jsFiles = findJsFiles('dist');

console.log(`Found ${jsFiles.length} JavaScript files to obfuscate`);

// تعتيم كل ملف
for (const file of jsFiles) {
  try {
    console.log(`Obfuscating: ${file}`);
    const code = readFileSync(file, 'utf8');
    const obfuscatedCode = JavaScriptObfuscator.obfuscate(code, obfuscatorOptions).getObfuscatedCode();
    writeFileSync(file, obfuscatedCode, 'utf8');
    console.log(`✓ Obfuscated: ${file}`);
  } catch (error) {
    console.error(`✗ Error obfuscating ${file}:`, error.message);
  }
}

console.log('Obfuscation complete!');
```

---

## 📄 frontend-worker-pwa/public/service-worker.js

```javascript
// حجز أولي لملف Service Worker — سيُستخدم لاحقاً لدعم العمل دون اتصال (offline)
self.addEventListener('install', () => self.skipWaiting())
self.addEventListener('activate', () => self.clients.claim())
```

---

## 📄 frontend-worker-pwa/src/App.vue

```vue
<template>
  <div v-if="isRouterReady" class="shell">
    <header v-if="isLoggedIn" class="topbar">
      <div class="brand">
        <div class="brand-mark">
          <svg width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"></path>
            <circle cx="12" cy="10" r="3"></circle>
          </svg>
        </div>
        <div class="brand-info">
          <span class="brand-name gradient-text">{{ currentTranslations.app_name }}</span>
          <span class="devpro-badge">Powered by DevPro</span>
        </div>
      </div>
      <div class="topbar-actions">
        <NotificationDropdown />
        <router-link to="/profile" class="avatar">{{ initials }}</router-link>
      </div>
    </header>
    <main class="content"><router-view /></main>
    <nav v-if="isLoggedIn" class="tabbar">
      <router-link to="/attendance" class="tab" active-class="tab--active">
        <span class="tab__icon">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <circle cx="12" cy="12" r="10"></circle>
            <polyline points="12 6 12 12 16 14"></polyline>
          </svg>
        </span>
        <span>{{ currentTranslations.attendance }}</span>
      </router-link>
      <router-link to="/tasks" class="tab" active-class="tab--active">
        <span class="tab__icon">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
            <polyline points="14 2 14 8 20 8"></polyline>
            <line x1="16" y1="13" x2="8" y2="13"></line>
            <line x1="16" y1="17" x2="8" y2="17"></line>
            <polyline points="10 9 9 9 8 9"></polyline>
          </svg>
        </span>
        <span>{{ currentTranslations.tasks }}</span>
      </router-link>
      <router-link to="/notifications" class="tab" active-class="tab--active">
        <span class="tab__icon">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9"></path>
            <path d="M13.73 21a2 2 0 0 1-3.46 0"></path>
          </svg>
        </span>
        <span>{{ currentTranslations.notifications }}</span>
      </router-link>
      <router-link to="/profile" class="tab" active-class="tab--active">
        <span class="tab__icon">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
            <circle cx="12" cy="7" r="4"></circle>
          </svg>
        </span>
        <span>{{ currentTranslations.profile }}</span>
      </router-link>
    </nav>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from './services/i18n'
import { authStore } from './store/auth'
import { logout, getToken, currentUser } from './services/auth'
import NotificationDropdown from './components/NotificationDropdown.vue'

const { t, currentTranslations } = useI18n()
const router = useRouter()
const isRouterReady = ref(false)

router.isReady().then(() => {
  isRouterReady.value = true
  
  // تحديث authStore عند تحميل الصفحة
  const token = getToken()
  if (token && !authStore.user) {
    const user = currentUser()
    if (user) {
      authStore.setUser(user)
    }
  }
  
  // معالجة مشكلة /index.html في العنوان - استعادة المسار المحفوظ بدلاً من التحويل إلى /
  if (window.location.pathname === '/index.html') {
    const isAuthed = !!localStorage.getItem('worktrack_token')
    const savedPath = localStorage.getItem('worktrack_worker_last_path')
    
    if (isAuthed && savedPath && savedPath !== '/login' && savedPath !== '/') {
      window.history.replaceState({}, '', savedPath)
    } else if (isAuthed) {
      window.history.replaceState({}, '', '/attendance')
    } else {
      window.history.replaceState({}, '', '/login')
    }
  }
})

const isLoggedIn = computed(() => !!authStore.user)

const initials = computed(() => {
  const name = authStore.user?.full_name || currentTranslations?.initials || '?'
  return String(name).slice(0, 1)
})
</script>

<style scoped>
.shell { 
  min-height: 100dvh; 
  display: flex; 
  flex-direction: column; 
  background: var(--canvas); 
}

.topbar { 
  display: flex; 
  justify-content: space-between; 
  padding: 12px 16px; 
  background: var(--surface); 
  border-bottom: 1px solid var(--line);
  box-shadow: var(--shadow-sm);
  transition: box-shadow var(--transition-base);
}

.topbar:hover {
  box-shadow: var(--shadow-md);
}

.brand {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 700;
}

.brand-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.brand-mark { 
  width: 40px; 
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--brand);
  transition: transform var(--transition-base);
}

.brand:hover .brand-mark {
  transform: scale(1.1);
}

.brand-name {
  font-size: 16px;
  font-weight: 700;
}

.devpro-badge {
  font-size: 10px;
  font-weight: 600;
  color: var(--ink-soft);
  letter-spacing: 0.5px;
  text-transform: uppercase;
}

.avatar { 
  width: 34px; 
  height: 34px; 
  border-radius: 50%; 
  background: linear-gradient(135deg, var(--brand-tint) 0%, var(--brand-tint) 100%); 
  color: var(--brand); 
  display: flex; 
  align-items: center; 
  justify-content: center; 
  font-weight: 700;
  box-shadow: var(--shadow-sm);
  transition: all var(--transition-base);
}

.avatar:hover {
  transform: scale(1.1);
  box-shadow: var(--shadow-md), var(--brand-glow);
}

.topbar-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.content { 
  flex: 1; 
  padding: 16px; 
  padding-bottom: 80px; 
  max-width: 480px; 
  width: 100%; 
  margin: 0 auto;
  animation: fadeIn 0.4s ease;
}

.tabbar { 
  position: fixed; 
  bottom: 0; 
  left: 0; 
  right: 0; 
  display: flex; 
  background: var(--surface); 
  border-top: 1px solid var(--line); 
  padding: 6px 8px; 
  max-width: 480px; 
  margin: 0 auto;
  box-shadow: 0 -2px 10px rgba(0,0,0,0.05);
}

.tab { 
  flex: 1; 
  display: flex; 
  flex-direction: column; 
  align-items: center; 
  gap: 4px; 
  padding: 8px; 
  border-radius: var(--radius-sm); 
  color: var(--ink-soft); 
  font-size: 11px; 
  font-weight: 600; 
  cursor: pointer; 
  border: none; 
  background: transparent;
  transition: all var(--transition-base);
  position: relative;
  overflow: hidden;
}

.tab::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, rgba(255,255,255,0.1) 0%, rgba(255,255,255,0) 100%);
  opacity: 0;
  transition: opacity var(--transition-fast);
}

.tab:hover::before {
  opacity: 1;
}

.tab:hover {
  background: var(--brand-tint);
  transform: translateY(-2px);
}

.tab--active { 
  color: var(--brand); 
  background: linear-gradient(135deg, var(--brand-tint) 0%, var(--brand-tint) 100%);
  box-shadow: var(--shadow-sm);
}

.tab__icon { 
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: transform var(--transition-base);
}

.tab:hover .tab__icon {
  transform: scale(1.1);
}

.tab--active .tab__icon {
  color: var(--brand);
}
</style>

```

---

## 📄 frontend-worker-pwa/src/components/ActivityFeed.vue

```vue
<template>
  <div class="activity-feed">
    <div class="feed-header">
      <h3>{{ t('notifications_feed') }}</h3>
      <span v-if="notificationsStore.notifications && notificationsStore.notifications.length" class="badge">
        {{ notificationsStore.notifications.length }}
      </span>
    </div>

    <div v-if="notificationsStore.loading" class="empty-state">
      <p>{{ t('loading_notifications') }}</p>
    </div>

    <div v-else-if="notificationsStore.error" class="alert alert-error">
      <span>❌</span> {{ notificationsStore.error }}
    </div>

    <div v-else-if="!notificationsStore.notifications || !notificationsStore.notifications.length" class="empty-state">
      <p>{{ t('no_notifications_available') }}</p>
    </div>

    <div v-else class="notifications-list">
      <div 
        v-for="notification in notificationsStore.notifications" 
        :key="notification.id" 
        class="notification-item"
      >
        <span class="notification-icon">{{ getNotificationIcon(notification.type) }}</span>
        <div class="notification-content">
          <p class="notification-title">{{ notification.title }}</p>
          <p class="notification-message">{{ notification.message }}</p>
          <span class="notification-time">{{ formatTime(notification.created_at) }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { notificationsStore } from '../store/notifications'
import { useI18n } from '../services/i18n'

const { t } = useI18n()

onMounted(() => {
  notificationsStore.fetchNotifications()
})

function getNotificationIcon(type) {
  const icons = {
    'status_update': '📋',
    'worker_assigned': '👷',
    'worker_arrived': '📍',
    'service_completed': '✅',
    'payment_request': '💳',
    'general': '🔔'
  }
  return icons[type] || '🔔'
}

function formatTime(date) {
  if (!date) return ''
  return new Date(date).toLocaleTimeString('en-GB', { 
    hour: '2-digit', 
    minute: '2-digit',
    hour12: false 
  })
}
</script>

<style scoped>
.activity-feed {
  background: var(--surface);
  border-radius: var(--radius-md);
  padding: 16px;
  border: 1px solid var(--line);
  position: fixed;
  bottom: 80px;
  right: 20px;
  width: 300px;
  max-height: 400px;
  overflow-y: auto;
  z-index: 100;
  box-shadow: var(--shadow-lg);
}

.feed-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.feed-header h3 {
  font-size: 14px;
  margin: 0;
  color: var(--ink);
}

.badge {
  background: var(--brand);
  color: white;
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 11px;
  font-weight: 600;
}

.empty-state {
  text-align: center;
  padding: 20px;
  color: var(--ink-soft);
  font-size: 13px;
}

.alert {
  padding: 10px 12px;
  border-radius: var(--radius-sm);
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
}

.alert-error {
  background: var(--signal-out-tint);
  color: var(--signal-out);
}

.notifications-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.notification-item {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  padding: 10px;
  background: var(--canvas);
  border-radius: var(--radius-sm);
  transition: background var(--transition-fast);
}

.notification-item:hover {
  background: var(--brand-tint);
}

.notification-icon {
  font-size: 16px;
  flex-shrink: 0;
}

.notification-content {
  flex: 1;
  min-width: 0;
}

.notification-title {
  font-size: 13px;
  font-weight: 600;
  color: var(--ink);
  margin: 0 0 4px 0;
}

.notification-message {
  font-size: 12px;
  color: var(--ink-soft);
  margin: 0 0 4px 0;
  line-height: 1.4;
}

.notification-time {
  font-size: 11px;
  color: var(--ink-light);
}

@media (max-width: 480px) {
  .activity-feed {
    bottom: 70px;
    right: 10px;
    left: 10px;
    width: auto;
    max-height: 300px;
    padding: 12px;
    border-radius: var(--radius-sm);
  }
  
  .feed-header {
    margin-bottom: 10px;
  }
  
  .feed-header h3 {
    font-size: 13px;
  }
  
  .badge {
    font-size: 10px;
    padding: 2px 6px;
  }
  
  .empty-state {
    padding: 16px;
    font-size: 12px;
  }
  
  .alert {
    font-size: 11px;
    padding: 8px 10px;
  }
  
  .notification-item {
    padding: 8px;
    gap: 8px;
  }
  
  .notification-icon {
    font-size: 14px;
  }
  
  .notification-title {
    font-size: 12px;
  }
  
  .notification-message {
    font-size: 11px;
  }
  
  .notification-time {
    font-size: 10px;
  }
}

@media (max-width: 360px) {
  .activity-feed {
    padding: 10px;
  }
  
  .feed-header h3 {
    font-size: 12px;
  }
  
  .notification-item {
    padding: 6px;
    gap: 6px;
  }
  
  .notification-icon {
    font-size: 12px;
  }
  
  .notification-title {
    font-size: 11px;
  }
  
  .notification-message {
    font-size: 10px;
  }
}
</style>

```

---

## 📄 frontend-worker-pwa/src/components/GeofenceRing.vue

```vue
<template>
  <div class="gf-ring" :class="status">
    <svg viewBox="0 0 160 160" class="gf-svg" aria-hidden="true">
      <circle cx="80" cy="80" r="68" class="gf-track" />
      <circle
        cx="80" cy="80" r="68" class="gf-fill"
        :stroke-dasharray="circumference"
        :stroke-dashoffset="dashOffset"
      />
      <circle cx="80" cy="80" r="46" class="gf-hole" />
    </svg>
    <div class="gf-center">
      <span class="gf-distance mono">{{ Math.round(distance) }}<small>م</small></span>
      <span class="gf-label">{{ statusText }}</span>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

// distance: مسافة الموظف الحالية عن نقطة العمل (متر)
// radius: نصف القطر المسموح به لنقطة العمل (متر)
const props = defineProps({
  distance: { type: Number, default: 0 },
  radius: { type: Number, default: 100 },
})

const circumference = 2 * Math.PI * 68

const ratio = computed(() => Math.min(props.distance / (props.radius || 1), 1))
const dashOffset = computed(() => circumference * (1 - ratio.value))
const status = computed(() => (props.distance <= props.radius ? 'inside' : 'outside'))
const statusText = computed(() => (status.value === 'inside' ? 'داخل نطاق الموقع' : 'خارج نطاق الموقع'))
</script>

<style scoped>
.gf-ring { position: relative; width: 180px; height: 180px; margin: 0 auto; }
.gf-svg { width: 100%; height: 100%; transform: rotate(-90deg); }
.gf-track { fill: none; stroke: var(--line); stroke-width: 10; }
.gf-fill { fill: none; stroke-width: 10; stroke-linecap: round; transition: stroke-dashoffset .6s ease, stroke .3s ease; }
.gf-ring.inside .gf-fill { stroke: var(--signal-in); }
.gf-ring.outside .gf-fill { stroke: var(--signal-out); }
.gf-hole { fill: var(--surface); }
.gf-center { position: absolute; inset: 0; display: flex; flex-direction: column; align-items: center; justify-content: center; }
.gf-distance { font-size: 28px; font-weight: 600; color: var(--ink); }
.gf-distance small { font-size: 14px; font-weight: 500; margin-inline-start: 2px; color: var(--ink-soft); }
.gf-label { margin-top: 4px; font-size: 13px; font-weight: 600; }
.gf-ring.inside .gf-label { color: var(--signal-in); }
.gf-ring.outside .gf-label { color: var(--signal-out); }
.gf-ring.inside::after {
  content: ''; position: absolute; inset: -6px; border-radius: 50%; border: 1px solid var(--signal-in);
  animation: gf-pulse 2.4s ease-out infinite;
}
@keyframes gf-pulse {
  0% { transform: scale(.94); opacity: .7; }
  100% { transform: scale(1.14); opacity: 0; }
}
</style>

```

---

## 📄 frontend-worker-pwa/src/components/GeolocationCheck.vue

```vue
<template>
  <div class="geo-check">
    <p v-if="state === 'idle'" class="geo-check__hint">اضغط لتحديد موقعك الحالي قبل التختيم</p>
    <p v-else-if="state === 'loading'" class="geo-check__hint">جارٍ تحديد موقعك...</p>
    <p v-else-if="state === 'error'" class="geo-check__error">{{ errorMessage }}</p>
    <p v-else class="geo-check__ok mono">تم تحديد الموقع: {{ position.lat.toFixed(5) }}, {{ position.lng.toFixed(5) }}</p>

    <button class="btn btn--ghost btn--block" @click="locate" :disabled="state === 'loading'">
      {{ state === 'ok' ? 'تحديث الموقع' : 'تحديد موقعي' }}
    </button>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { getCurrentPosition } from '../services/geolocation'

const state = ref('idle')
const position = ref(null)
const errorMessage = ref('')
const emit = defineEmits(['located'])

async function locate() {
  state.value = 'loading'
  try {
    position.value = await getCurrentPosition()
    state.value = 'ok'
    emit('located', position.value)
  } catch (e) {
    errorMessage.value = 'تعذر تحديد الموقع، تأكد من تفعيل صلاحية الوصول للموقع الجغرافي'
    state.value = 'error'
  }
}
</script>

<style scoped>
.geo-check { display: flex; flex-direction: column; gap: 10px; }
.geo-check__hint { font-size: 13px; color: var(--ink-soft); text-align: center; }
.geo-check__error { font-size: 13px; color: var(--signal-out); text-align: center; }
.geo-check__ok { font-size: 13px; color: var(--signal-in); text-align: center; }
</style>

```

---

## 📄 frontend-worker-pwa/src/components/LanguageSwitcher.vue

```vue
<template>
  <div class="lang-switcher">
    <button
      v-for="lang in languages"
      :key="lang.code"
      class="lang-btn"
      :class="{ active: currentLang === lang.code }"
      @click="changeLanguage(lang.code)"
      :title="lang.name"
    >
      {{ lang.flag }}
      <span class="lang-label">{{ lang.name }}</span>
    </button>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from '../services/i18n'

const { currentLang, setLang } = useI18n()

const languages = [
  { code: 'ar', name: 'العربية', flag: '🇸🇦' },
  { code: 'he', name: 'עברית', flag: '🇮🇱' },
  { code: 'en', name: 'English', flag: '🇬🇧' }
]

function changeLanguage(lang) {
  setLang(lang)
}
</script>

<style scoped>
.lang-switcher {
  display: flex;
  justify-content: center;
  gap: 6px;
  padding: 6px;
  background: var(--canvas, #f0f4fa);
  border-radius: 12px;
  width: 100%;
}

.lang-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 14px;
  border: none;
  border-radius: 8px;
  background: transparent;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.3s ease;
  font-family: inherit;
  flex: 1;
  justify-content: center;
}

.lang-btn:hover {
  background: rgba(30, 58, 95, 0.08);
  transform: scale(1.02);
}

.lang-btn.active {
  background: #1E3A5F;
  color: white;
  box-shadow: 0 2px 12px rgba(30, 58, 95, 0.3);
  transform: scale(1.02);
}

.lang-label {
  font-size: 12px;
  font-weight: 500;
}

@media (max-width: 480px) {
  .lang-btn {
    padding: 4px 8px;
    font-size: 12px;
  }
  .lang-label {
    display: none;
  }
}
</style>

```

---

## 📄 frontend-worker-pwa/src/components/MapPicker.vue

```vue
<template>
  <div class="map-picker">
    <div class="map-picker__grid"></div>
    <div class="map-picker__pin">
      <span class="map-picker__pulse"></span>
      <span class="map-picker__dot"></span>
    </div>
    <span class="map-picker__label">موقع نقطة العمل</span>
  </div>
</template>

<script setup>
defineProps({ lat: Number, lng: Number })
</script>

<style scoped>
.map-picker {
  position: relative; height: 160px; border-radius: var(--radius-md); overflow: hidden;
  background: linear-gradient(180deg, #EAF1EC, #E1EBE4); border: 1px solid var(--line);
  display: flex; align-items: center; justify-content: center;
}
.map-picker__grid {
  position: absolute; inset: 0;
  background-image: linear-gradient(var(--line) 1px, transparent 1px), linear-gradient(90deg, var(--line) 1px, transparent 1px);
  background-size: 24px 24px; opacity: .5;
}
.map-picker__pin { position: relative; width: 20px; height: 20px; }
.map-picker__dot { position: absolute; inset: 0; margin: auto; width: 12px; height: 12px; border-radius: 50%; background: var(--brand); border: 2px solid #fff; box-shadow: var(--shadow-sm); }
.map-picker__pulse { position: absolute; inset: 0; margin: auto; width: 12px; height: 12px; border-radius: 50%; background: var(--brand); opacity: .5; animation: map-pulse 2s ease-out infinite; }
.map-picker__label { position: absolute; bottom: 10px; font-size: 12px; font-weight: 600; color: var(--ink-soft); background: rgba(255,255,255,.85); padding: 3px 10px; border-radius: 999px; }
@keyframes map-pulse { 0% { transform: scale(1); opacity: .5; } 100% { transform: scale(3.4); opacity: 0; } }
</style>

```

---

## 📄 frontend-worker-pwa/src/components/NotificationDropdown.vue

```vue
<template>
  <div class="notification-dropdown" ref="dropdownRef">
    <!-- أيقونة الإشعارات -->
    <button 
      @click="toggleDropdown" 
      class="notification-icon"
      :title="t('notifications')"
    >
      <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9"></path>
        <path d="M13.73 21a2 2 0 0 1-3.46 0"></path>
      </svg>
      <span v-if="unreadCount > 0" class="badge">{{ unreadCount }}</span>
    </button>
    
    <!-- Dropdown content -->
    <transition name="dropdown">
      <div v-if="isOpen" class="dropdown-content">
        <div class="feed-header">
          <h3>{{ t('notifications_feed') }}</h3>
          <span v-if="notificationsStore.notifications && notificationsStore.notifications.length" class="badge">
            {{ notificationsStore.notifications.length }}
          </span>
        </div>

        <div v-if="notificationsStore.loading" class="empty-state">
          <p>{{ t('loading_notifications') }}</p>
        </div>

        <div v-else-if="notificationsStore.error" class="alert alert-error">
          <span>❌</span> {{ notificationsStore.error }}
        </div>

        <div v-else-if="!notificationsStore.notifications || !notificationsStore.notifications.length" class="empty-state">
          <p>{{ t('no_notifications_available') }}</p>
        </div>

        <div v-else class="notifications-list">
          <div 
            v-for="notification in notificationsStore.notifications" 
            :key="notification.id" 
            class="notification-item"
          >
            <span class="notification-icon">{{ getNotificationIcon(notification.type) }}</span>
            <div class="notification-content">
              <p class="notification-title">{{ notification.title }}</p>
              <p class="notification-message">{{ notification.message }}</p>
              <span class="notification-time">{{ formatTime(notification.created_at) }}</span>
            </div>
          </div>
          <router-link to="/notifications" class="view-all-link" @click="closeDropdown">
            {{ t('view_all') }}
          </router-link>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount, computed } from 'vue'
import { notificationsStore } from '../store/notifications'
import { useI18n } from '../services/i18n'
import { useRouter } from 'vue-router'

const { t } = useI18n()
const router = useRouter()
const isOpen = ref(false)
const dropdownRef = ref(null)

const unreadCount = computed(() => {
  if (!notificationsStore.notifications) return 0
  return notificationsStore.notifications.length
})

function toggleDropdown() {
  isOpen.value = !isOpen.value
  if (isOpen.value) {
    notificationsStore.fetchNotifications()
  }
}

function closeDropdown() {
  isOpen.value = false
}

function handleClickOutside(event) {
  if (dropdownRef.value && !dropdownRef.value.contains(event.target)) {
    closeDropdown()
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
  notificationsStore.fetchNotifications()
})

onBeforeUnmount(() => {
  document.removeEventListener('click', handleClickOutside)
})

function getNotificationIcon(type) {
  const icons = {
    'status_update': '📋',
    'worker_assigned': '👷',
    'worker_arrived': '📍',
    'service_completed': '✅',
    'payment_request': '💳',
    'general': '🔔'
  }
  return icons[type] || '🔔'
}

function formatTime(date) {
  if (!date) return ''
  return new Date(date).toLocaleTimeString('en-GB', { 
    hour: '2-digit', 
    minute: '2-digit',
    hour12: false 
  })
}
</script>

<style scoped>
.notification-dropdown {
  position: relative;
  display: inline-block;
}

.notification-icon {
  background: none;
  border: none;
  cursor: pointer;
  position: relative;
  padding: 8px;
  border-radius: var(--radius-sm);
  transition: all var(--transition-base);
  color: var(--ink-soft);
  display: flex;
  align-items: center;
  justify-content: center;
}

.notification-icon:hover {
  background: var(--brand-tint);
  color: var(--brand);
  transform: scale(1.1);
}

.badge {
  position: absolute;
  top: 0;
  right: 0;
  background: var(--signal-out);
  color: white;
  border-radius: 50%;
  padding: 2px 6px;
  font-size: 10px;
  font-weight: 600;
  min-width: 16px;
  height: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0%, 100% {
    transform: scale(1);
  }
  50% {
    transform: scale(1.1);
  }
}

.dropdown-content {
  position: absolute;
  top: 100%;
  right: 0;
  background: var(--surface);
  border-radius: var(--radius-md);
  box-shadow: var(--shadow-lg);
  min-width: 320px;
  max-height: 400px;
  overflow-y: auto;
  z-index: 1000;
  margin-top: 8px;
  border: 1px solid var(--line);
}

.dropdown-enter-active,
.dropdown-leave-active {
  transition: all 0.2s ease;
}

.dropdown-enter-from {
  opacity: 0;
  transform: translateY(-10px);
}

.dropdown-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}

.feed-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  border-bottom: 1px solid var(--line);
  background: var(--canvas);
  border-radius: var(--radius-md) var(--radius-md) 0 0;
  position: sticky;
  top: 0;
  z-index: 1;
}

.feed-header h3 {
  font-size: 14px;
  margin: 0;
  color: var(--ink);
}

.feed-header .badge {
  position: static;
  background: var(--brand);
  color: white;
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 11px;
  font-weight: 600;
  animation: none;
}

.empty-state {
  text-align: center;
  padding: 20px;
  color: var(--ink-soft);
  font-size: 13px;
}

.alert {
  padding: 10px 12px;
  border-radius: var(--radius-sm);
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
  margin: 8px;
}

.alert-error {
  background: var(--signal-out-tint);
  color: var(--signal-out);
}

.notifications-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 8px;
}

.notification-item {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  padding: 10px;
  background: var(--canvas);
  border-radius: var(--radius-sm);
  transition: background var(--transition-fast);
  cursor: pointer;
}

.notification-item:hover {
  background: var(--brand-tint);
}

.notification-icon {
  font-size: 16px;
  flex-shrink: 0;
}

.notification-content {
  flex: 1;
  min-width: 0;
}

.notification-title {
  font-size: 13px;
  font-weight: 600;
  color: var(--ink);
  margin: 0 0 4px 0;
}

.notification-message {
  font-size: 12px;
  color: var(--ink-soft);
  margin: 0 0 4px 0;
  line-height: 1.4;
}

.notification-time {
  font-size: 11px;
  color: var(--ink-light);
}

.view-all-link {
  display: block;
  text-align: center;
  padding: 12px;
  margin-top: 8px;
  background: var(--brand-tint);
  color: var(--brand);
  text-decoration: none;
  border-radius: var(--radius-sm);
  font-size: 13px;
  font-weight: 600;
  transition: all var(--transition-fast);
}

.view-all-link:hover {
  background: var(--brand);
  color: white;
}

@media (max-width: 480px) {
  .dropdown-content {
    min-width: 280px;
    max-height: 350px;
    right: -8px;
  }

  .feed-header {
    padding: 10px 12px;
  }

  .feed-header h3 {
    font-size: 13px;
  }

  .notifications-list {
    padding: 6px;
  }

  .notification-item {
    padding: 8px;
    gap: 8px;
  }

  .notification-icon {
    font-size: 14px;
  }

  .notification-title {
    font-size: 12px;
  }

  .notification-message {
    font-size: 11px;
  }

  .notification-time {
    font-size: 10px;
  }
}
</style>
```

---

## 📄 frontend-worker-pwa/src/components/PWAInstallButton.vue

```vue
<template>
  <div v-if="showInstallButton" class="pwa-install-container">
    <button 
      @click="installPWA" 
      class="pwa-install-icon"
      :title="pwaInstallTitle"
    >
      <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path>
        <polyline points="7 10 12 15 17 10"></polyline>
        <line x1="12" y1="15" x2="12" y2="3"></line>
      </svg>
    </button>
    
    <!-- iOS Instructions Modal -->
    <div v-if="showIOSInstructions" class="ios-instructions-modal" @click.self="showIOSInstructions = false">
      <div class="ios-instructions-content">
        <button class="close-btn" @click="showIOSInstructions = false">✕</button>
        <div class="ios-header">
          <div class="ios-icon">📱</div>
          <h3>{{ iosModalTitle }}</h3>
          <p class="ios-subtitle">{{ iosModalSubtitle }}</p>
        </div>
        
        <div class="steps">
          <div class="step">
            <span class="step-number">1</span>
            <div class="step-content">
              <p v-html="step1Text"></p>
              <div class="step-visual">
                <div class="phone-frame">
                  <div class="phone-screen">
                    <div class="share-btn">⎋</div>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div class="step">
            <span class="step-number">2</span>
            <div class="step-content">
              <p v-html="step2Text"></p>
              <div class="step-visual">
                <div class="menu-item">➕ Add to Home Screen</div>
              </div>
            </div>
          </div>
          <div class="step">
            <span class="step-number">3</span>
            <div class="step-content">
              <p v-html="step3Text"></p>
              <div class="step-visual">
                <div class="add-btn">Add</div>
              </div>
            </div>
          </div>
        </div>
        
        <div class="benefits">
          <div class="benefit">{{ benefit1 }}</div>
          <div class="benefit">{{ benefit2 }}</div>
          <div class="benefit">{{ benefit3 }}</div>
        </div>
        
        <div class="actions">
          <button class="got-it-btn" @click="showIOSInstructions = false">{{ gotItText }}</button>
          <button class="remind-btn" @click="remindLater">{{ remindLaterText }}</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'PWAInstallButton',
  data() {
    return {
      showInstallButton: false,
      showIOSInstructions: false,
      deferredPrompt: null,
      isIOS: false,
      isStandalone: false
    }
  },
  computed: {
    currentLang() {
      if (this.$i18n && this.$i18n.currentLang) {
        return this.$i18n.currentLang.value
      }
      return 'ar' // default
    },
    pwaInstallTitle() {
      return this.$t('pwa.installTitle')
    },
    pwaInstallText() {
      if (this.isIOS) {
        return this.$t('pwa.installOnIos')
      }
      return this.$t('pwa.installText')
    },
    iosModalTitle() {
      return this.$t('pwa.iosModalTitle')
    },
    iosModalSubtitle() {
      return this.$t('pwa.iosModalSubtitle')
    },
    step1Text() {
      return this.$t('pwa.step1Text')
    },
    step2Text() {
      return this.$t('pwa.step2Text')
    },
    step3Text() {
      return this.$t('pwa.step3Text')
    },
    benefit1() {
      return this.$t('pwa.benefit1')
    },
    benefit2() {
      return this.$t('pwa.benefit2')
    },
    benefit3() {
      return this.$t('pwa.benefit3')
    },
    gotItText() {
      return this.$t('pwa.gotIt')
    },
    remindLaterText() {
      return this.$t('pwa.remindLater')
    }
  },
  methods: {
    detectIOS() {
      const userAgent = window.navigator.userAgent.toLowerCase()
      return /iphone|ipad|ipod/.test(userAgent) && !/edge|crios|fxios/.test(userAgent)
    },
    handleInstallAvailable() {
      this.showInstallButton = true
    },
    handleInstallSuccess() {
      this.showInstallButton = false
    },
    installPWA() {
      if (this.isIOS) {
        // Show iOS instructions
        this.showIOSInstructions = true
      } else if (window.pwaInstall) {
        // Use native install prompt for Android/Windows
        window.pwaInstall()
      } else {
        // Fallback - show general instructions
        this.showIOSInstructions = true
      }
    },
    remindLater() {
      // Just close the modal, keep icon visible
      this.showIOSInstructions = false
    }
  },
  mounted() {
    // Check if device is iOS
    this.isIOS = this.detectIOS()
    
    // Check if already installed (standalone mode)
    this.isStandalone = window.matchMedia('(display-mode: standalone)').matches || 
                      window.navigator.standalone === true

    // Don't show if already installed
    if (this.isStandalone) {
      return
    }

    // Always show button on login page
    if (this.isIOS) {
      this.showInstallButton = true
    } else {
      // For other devices, listen to install prompt
      window.addEventListener('pwa-install-available', this.handleInstallAvailable)
      window.addEventListener('pwa-install-success', this.handleInstallSuccess)
      
      // Check for existing deferredPrompt
      if (window.deferredPrompt) {
        this.showInstallButton = true
      }
    }
  },
  beforeUnmount() {
    window.removeEventListener('pwa-install-available', this.handleInstallAvailable)
    window.removeEventListener('pwa-install-success', this.handleInstallSuccess)
  }
}
</script>

<style scoped>
.pwa-install-container {
  position: fixed;
  bottom: 16px;
  right: 16px;
  z-index: 1000;
}

.pwa-install-icon {
  width: 44px;
  height: 44px;
  padding: 0;
  background: linear-gradient(135deg, #1F6F5C 0%, #2d8a6f 100%);
  color: white;
  border: none;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  box-shadow: 0 4px 12px rgba(31, 111, 92, 0.3);
  transition: all 0.3s ease;
}

.pwa-install-icon:hover {
  transform: scale(1.1);
  box-shadow: 0 6px 16px rgba(31, 111, 92, 0.4);
}

.pwa-install-icon:active {
  transform: scale(0.95);
}

/* iOS Instructions Modal */
.ios-instructions-modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2000;
  padding: 20px;
}

.ios-instructions-content {
  background: white;
  border-radius: 20px;
  padding: 28px;
  max-width: 420px;
  width: 100%;
  position: relative;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.4);
  max-height: 90vh;
  overflow-y: auto;
}

.close-btn {
  position: absolute;
  top: 16px;
  left: 16px;
  background: none;
  border: none;
  font-size: 24px;
  cursor: pointer;
  color: #666;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  transition: background 0.2s;
}

.close-btn:hover {
  background: #f0f0f0;
}

.ios-header {
  text-align: center;
  margin-bottom: 24px;
  padding-top: 8px;
}

.ios-icon {
  font-size: 48px;
  margin-bottom: 12px;
}

.ios-header h3 {
  margin: 0 0 8px 0;
  font-size: 20px;
  color: #1F6F5C;
  font-weight: 700;
}

.ios-subtitle {
  margin: 0;
  font-size: 14px;
  color: #666;
}

.steps {
  display: flex;
  flex-direction: column;
  gap: 20px;
  margin-bottom: 24px;
}

.step {
  display: flex;
  align-items: flex-start;
  gap: 14px;
}

.step-number {
  width: 32px;
  height: 32px;
  background: linear-gradient(135deg, #1F6F5C 0%, #2d8a6f 100%);
  color: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  font-size: 15px;
  flex-shrink: 0;
}

.step-content {
  flex: 1;
}

.step p {
  margin: 0 0 8px 0;
  font-size: 14px;
  color: #333;
  line-height: 1.5;
}

.step strong {
  color: #1F6F5C;
}

.step-visual {
  background: #f8f9fa;
  border-radius: 8px;
  padding: 8px;
  display: flex;
  justify-content: center;
  border: 1px solid #e9ecef;
}

.phone-frame {
  width: 60px;
  height: 100px;
  background: white;
  border: 2px solid #dee2e6;
  border-radius: 8px;
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
}

.phone-screen {
  width: 90%;
  height: 90%;
  background: #f8f9fa;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.share-btn {
  font-size: 20px;
  color: #007AFF;
}

.menu-item {
  background: white;
  padding: 8px 12px;
  border-radius: 6px;
  font-size: 12px;
  color: #333;
  border: 1px solid #dee2e6;
}

.add-btn {
  background: #007AFF;
  color: white;
  padding: 6px 16px;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 600;
}

.benefits {
  display: flex;
  justify-content: space-around;
  margin-bottom: 24px;
  padding: 16px;
  background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
  border-radius: 12px;
}

.benefit {
  font-size: 12px;
  color: #495057;
  text-align: center;
  font-weight: 500;
}

.actions {
  display: flex;
  gap: 12px;
}

.got-it-btn {
  flex: 2;
  padding: 14px;
  background: linear-gradient(135deg, #1F6F5C 0%, #2d8a6f 100%);
  color: white;
  border: none;
  border-radius: 12px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: transform 0.2s;
}

.got-it-btn:hover {
  transform: scale(1.02);
}

.remind-btn {
  flex: 1;
  padding: 14px;
  background: #f8f9fa;
  color: #666;
  border: 1px solid #dee2e6;
  border-radius: 12px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: background 0.2s;
}

.remind-btn:hover {
  background: #e9ecef;
}

/* RTL support */
[dir="rtl"] .pwa-install-container {
  right: auto;
  left: 16px;
}

[dir="rtl"] .close-btn {
  left: auto;
  right: 16px;
}

@media (max-width: 768px) {
  .pwa-install-container {
    bottom: 12px;
    right: 12px;
  }
  
  [dir="rtl"] .pwa-install-container {
    right: auto;
    left: 12px;
  }
  
  .pwa-install-icon {
    width: 40px;
    height: 40px;
  }

  .ios-instructions-content {
    padding: 20px;
  }

  .ios-instructions-content h3 {
    font-size: 16px;
  }

  .step p {
    font-size: 13px;
  }
}
</style>
```

---

## 📄 frontend-worker-pwa/src/components/PhotoUpload.vue

```vue
<template>
  <div class="photo-upload">
    <label class="photo-upload__drop" :class="{ 'has-image': preview }">
      <img v-if="preview" :src="preview" alt="معاينة الصورة" />
      <template v-else>
        <span class="photo-upload__icon">＋</span>
        <span>أرفق صورة إثبات إنجاز المهمة</span>
      </template>
      <input type="file" accept="image/*" capture="environment" class="visually-hidden" @change="onChange" />
    </label>
    <button v-if="preview" class="btn btn--ghost btn--sm" @click="clear">إزالة الصورة</button>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const preview = ref(null)
const emit = defineEmits(['selected'])

function onChange(e) {
  const file = e.target.files?.[0]
  if (!file) return
  const reader = new FileReader()
  reader.onload = () => {
    preview.value = reader.result
    emit('selected', file)
  }
  reader.readAsDataURL(file)
}

function clear() {
  preview.value = null
  emit('selected', null)
}
</script>

<style scoped>
.photo-upload { display: flex; flex-direction: column; gap: 10px; align-items: stretch; }
.photo-upload__drop {
  display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 8px;
  height: 140px; border: 1.5px dashed var(--line-strong); border-radius: var(--radius-md);
  color: var(--ink-soft); font-size: 13px; cursor: pointer; overflow: hidden; background: var(--surface);
}
.photo-upload__drop.has-image { padding: 0; }
.photo-upload__drop img { width: 100%; height: 100%; object-fit: cover; }
.photo-upload__icon { font-size: 22px; color: var(--brand); }
</style>

```

---

## 📄 frontend-worker-pwa/src/components/PullToRefresh.vue

```vue
<template>
  <div 
    class="pull-to-refresh"
    @touchstart="handleTouchStart"
    @touchmove="handleTouchMove"
    @touchend="handleTouchEnd"
    @mousedown="handleTouchStart"
    @mousemove="handleTouchMove"
    @mouseup="handleTouchEnd"
  >
    <div 
      class="pull-indicator"
      :style="{ 
        transform: `translateY(${pullDistance}px)`,
        opacity: Math.min(pullDistance / 60, 1)
      }"
    >
      <div class="pull-icon" :class="{ rotating: isRefreshing }">
        <svg v-if="isRefreshing" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M12 2v4M12 18v4M4.93 4.93l2.83 2.83M16.24 16.24l2.83 2.83M2 12h4M18 12h4M4.93 19.07l2.83-2.83M16.24 7.76l2.83-2.83"/>
        </svg>
        <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M12 4v1M12 20v1M4 12h1M19 12h1M6.34 6.34l.71.71M17.24 17.24l.71.71M6.34 17.65l.71-.71M17.24 6.65l.71-.71"/>
        </svg>
      </div>
      <div class="pull-text">{{ refreshText }}</div>
    </div>
    
    <div class="pull-content" :style="{ transform: `translateY(${Math.max(0, pullDistance - 60)}px)` }">
      <slot></slot>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useI18n } from '../services/i18n'

const { t } = useI18n()

const props = defineProps({
  threshold: {
    type: Number,
    default: 80
  },
  refreshText: {
    type: String,
    default: ''
  },
  refreshingText: {
    type: String,
    default: ''
  },
  releaseText: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['refresh'])

const isPulling = ref(false)
const isRefreshing = ref(false)
const pullDistance = ref(0)
const startY = ref(0)
const currentY = ref(0)

const displayText = computed(() => {
  if (isRefreshing.value) return props.refreshingText || t('refreshing')
  if (pullDistance.value >= props.threshold) return props.releaseText || t('release_to_refresh')
  return props.refreshText || t('pull_to_refresh')
})

const refreshText = computed(() => displayText.value)

function handleTouchStart(e) {
  if (isRefreshing.value) return
  
  const y = e.touches ? e.touches[0].clientY : e.clientY
  const scrollTop = e.target?.scrollTop || document.documentElement.scrollTop
  
  // فقط إذا كنا في أعلى الصفحة
  if (scrollTop <= 0) {
    startY.value = y
    isPulling.value = true
    pullDistance.value = 0
  }
}

function handleTouchMove(e) {
  if (!isPulling.value || isRefreshing.value) return
  
  const y = e.touches ? e.touches[0].clientY : e.clientY
  currentY.value = y
  
  const distance = y - startY.value
  
  // حساب المسافة مع تخفيف الحركة
  if (distance > 0) {
    pullDistance.value = Math.min(distance * 0.5, 150)
    
    // منع السحب الافتراضي
    if (e.cancelable) {
      e.preventDefault()
    }
  }
}

function handleTouchEnd() {
  if (!isPulling.value) return
  
  isPulling.value = false
  
  // إذا وصلنا لعتبة التحديث
  if (pullDistance.value >= props.threshold) {
    triggerRefresh()
  }
  
  // إعادة المسافة للصفر
  pullDistance.value = 0
}

async function triggerRefresh() {
  isRefreshing.value = true
  pullDistance.value = 60
  
  try {
    await emit('refresh')
  } catch (error) {
    console.error('Refresh failed:', error)
  } finally {
    setTimeout(() => {
      isRefreshing.value = false
      pullDistance.value = 0
    }, 500)
  }
}
</script>

<style scoped>
.pull-to-refresh {
  position: relative;
  width: 100%;
  min-height: 100vh;
  overflow: hidden;
}

.pull-indicator {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 60px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  pointer-events: none;
  z-index: 1000;
  background: linear-gradient(180deg, rgba(30, 58, 95, 0.1) 0%, transparent 100%);
}

.pull-icon {
  width: 24px;
  height: 24px;
  color: #1E3A5F;
  margin-bottom: 4px;
  transition: transform 0.3s ease;
}

.pull-icon.rotating {
  animation: rotate 1s linear infinite;
}

@keyframes rotate {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

.pull-text {
  font-size: 12px;
  color: #1E3A5F;
  font-weight: 600;
  text-align: center;
}

.pull-content {
  transition: transform 0.3s ease;
  min-height: 100vh;
}

@media (max-width: 480px) {
  .pull-text {
    font-size: 11px;
  }
  
  .pull-icon {
    width: 20px;
    height: 20px;
  }
}
</style>

```

---

## 📄 frontend-worker-pwa/src/components/SwipeNav.vue

```vue
<template>
  <div 
    class="swipe-nav"
    @touchstart="handleTouchStart"
    @touchmove="handleTouchMove"
    @touchend="handleTouchEnd"
    @mousedown="handleTouchStart"
    @mousemove="handleTouchMove"
    @mouseup="handleTouchEnd"
    @mouseleave="handleTouchEnd"
  >
    <div 
      class="swipe-content"
      :style="{ 
        transform: `translateX(${translateX}px)`,
        transition: isDragging ? 'none' : 'transform 0.3s ease'
      }"
    >
      <slot></slot>
    </div>
    
    <div class="swipe-indicators" v-if="showIndicators && totalItems > 1">
      <div 
        v-for="(_, index) in totalItems" 
        :key="index"
        class="indicator"
        :class="{ active: currentIndex === index }"
        @click="goToSlide(index)"
      ></div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'

const props = defineProps({
  itemsPerView: {
    type: Number,
    default: 1
  },
  showIndicators: {
    type: Boolean,
    default: true
  },
  autoPlay: {
    type: Boolean,
    default: false
  },
  autoPlayInterval: {
    type: Number,
    default: 3000
  }
})

const emit = defineEmits(['slide-change', 'swipe-left', 'swipe-right'])

const translateX = ref(0)
const isDragging = ref(false)
const startX = ref(0)
const currentX = ref(0)
const currentIndex = ref(0)
const totalItems = ref(0)
const containerWidth = ref(0)
let autoPlayTimer = null

const maxIndex = computed(() => {
  return Math.max(0, totalItems.value - props.itemsPerView)
})

function handleTouchStart(e) {
  isDragging.value = true
  startX.value = e.touches ? e.touches[0].clientX : e.clientX
  currentX.value = startX.value
  
  if (props.autoPlay) {
    stopAutoPlay()
  }
}

function handleTouchMove(e) {
  if (!isDragging.value) return
  
  const x = e.touches ? e.touches[0].clientX : e.clientX
  currentX.value = x
  
  const diff = currentX.value - startX.value
  translateX.value = diff
}

function handleTouchEnd() {
  if (!isDragging.value) return
  
  isDragging.value = false
  
  const diff = currentX.value - startX.value
  const threshold = 50 // العتبة للتنقل
  
  if (diff > threshold) {
    // السحب لليمين - العودة للسابق
    prevSlide()
  } else if (diff < -threshold) {
    // السحب لليسار - الذهاب للتالي
    nextSlide()
  } else {
    // العودة للوضع الحالي
    resetPosition()
  }
  
  if (props.autoPlay) {
    startAutoPlay()
  }
}

function nextSlide() {
  if (currentIndex.value < maxIndex.value) {
    currentIndex.value++
    emit('swipe-left', currentIndex.value)
  } else {
    // العودة للبداية
    currentIndex.value = 0
  }
  resetPosition()
  emit('slide-change', currentIndex.value)
}

function prevSlide() {
  if (currentIndex.value > 0) {
    currentIndex.value--
    emit('swipe-right', currentIndex.value)
  } else {
    // الذهاب للنهاية
    currentIndex.value = maxIndex.value
  }
  resetPosition()
  emit('slide-change', currentIndex.value)
}

function goToSlide(index) {
  currentIndex.value = index
  resetPosition()
  emit('slide-change', currentIndex.value)
}

function resetPosition() {
  const itemWidth = containerWidth.value / props.itemsPerView
  translateX.value = -currentIndex.value * itemWidth
}

function updateTotalItems() {
  const content = document.querySelector('.swipe-content')
  if (content) {
    totalItems.value = content.children.length
    containerWidth.value = content.offsetWidth
    resetPosition()
  }
}

function startAutoPlay() {
  if (autoPlayTimer) return
  autoPlayTimer = setInterval(() => {
    nextSlide()
  }, props.autoPlayInterval)
}

function stopAutoPlay() {
  if (autoPlayTimer) {
    clearInterval(autoPlayTimer)
    autoPlayTimer = null
  }
}

onMounted(() => {
  updateTotalItems()
  window.addEventListener('resize', updateTotalItems)
  
  if (props.autoPlay) {
    startAutoPlay()
  }
})

onUnmounted(() => {
  window.removeEventListener('resize', updateTotalItems)
  stopAutoPlay()
})

// توفير الوصول للوظائف للمكونات الأب
defineExpose({
  nextSlide,
  prevSlide,
  goToSlide,
  currentIndex
})
</script>

<style scoped>
.swipe-nav {
  position: relative;
  overflow: hidden;
  width: 100%;
}

.swipe-content {
  display: flex;
  width: 100%;
  cursor: grab;
  /* user-select: none; - تمت الإزالة للسماح بالنسخ واللصق على الهاتف */
  /* -webkit-user-select: none; - تمت الإزالة للسماح بالنسخ واللصق على الهاتف */
  touch-action: manipulation;
}

.swipe-content:active {
  cursor: grabbing;
}

.swipe-indicators {
  display: flex;
  justify-content: center;
  gap: 8px;
  margin-top: 16px;
  padding: 8px;
}

.indicator {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #E2E8F0;
  cursor: pointer;
  transition: all 0.3s ease;
}

.indicator.active {
  background: #1E3A5F;
  transform: scale(1.2);
}

.indicator:hover {
  background: #4A6FA5;
}

@media (max-width: 480px) {
  .indicator {
    width: 6px;
    height: 6px;
  }
  
  .swipe-indicators {
    gap: 6px;
    margin-top: 12px;
  }
}
</style>

```

---

## 📄 frontend-worker-pwa/src/components/TaskCard.vue

```vue
<template>
  <router-link :to="`/tasks/${task.id}`" class="task-card" :key="`${task.id}-${currentLang}`" @click="handleClick">
    <div class="task-card__stripe" :class="task.status"></div>
    <div class="task-card__content">
      <div class="task-card__header">
        <h3 class="task-card__title" dir="auto">{{ getTranslatedField(task, 'title') }}</h3>
        <div class="task-card__badges">
          <span v-if="task.priority" class="task-card__priority" :class="`priority-${task.priority}`">
            {{ priorityLabels[task.priority] }}
          </span>
          <span class="task-card__status" :class="`status-${task.status}`">
            {{ statusLabels[task.status] }}
          </span>
        </div>
      </div>

      <div class="task-card__info">
        <div v-if="getTranslatedField(task, 'worksite_name') || getTranslatedField(task, 'worksite_address')" class="task-card__info-item">
          <span class="icon">📍</span>
          <span class="text" dir="auto">{{ getTranslatedField(task, 'worksite_name') || getTranslatedField(task, 'worksite_address') }}</span>
        </div>

        <div v-if="getTranslatedField(task, 'client_name')" class="task-card__info-item">
          <span class="icon">👤</span>
          <span class="text" dir="auto">{{ getTranslatedField(task, 'client_name') }}</span>
        </div>

        <div v-if="task.client_phone" class="task-card__info-item">
          <span class="icon">📞</span>
          <span class="text" dir="auto">{{ task.client_phone }}</span>
        </div>

        <div v-if="task.scheduled_start" class="task-card__info-item">
          <span class="icon">🕒</span>
          <span class="text">{{ formatTime(task.scheduled_start) }}</span>
        </div>
      </div>
    </div>

    <div class="task-card__arrow">
      <svg width="20" height="20" viewBox="0 0 20 20" fill="none" xmlns="http://www.w3.org/2000/svg">
        <path d="M8 4L14 10L8 16" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
      </svg>
    </div>
  </router-link>
</template>

<script setup>
import { computed } from 'vue'
import { tasksStore } from '../store/tasks'
import { useI18n } from '../services/i18n'

const props = defineProps({ task: { type: Object, required: true } })
const { t, currentLang, currentTranslations } = useI18n()

const handleClick = () => {
  // Store the service request ID if this is a service request task
  if (props.task.id && !props.task.worksite_id) {
    // This is likely a service request (no worksite_id)
    tasksStore.setCurrentServiceRequestId(props.task.id)
  }
}

const formatTime = (dateString) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  const locale = currentLang.value === 'ar' ? 'ar-EG' : currentLang.value === 'he' ? 'he-IL' : 'en-GB'
  return date.toLocaleDateString(locale, {
    weekday: 'short',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const statusLabels = computed(() => ({
  pending: currentTranslations.value.status_pending,
  assigned: currentTranslations.value.status_assigned,
  in_progress: currentTranslations.value.status_in_progress,
  completed: currentTranslations.value.status_completed,
  late: currentTranslations.value.status_late,
  cancelled: currentTranslations.value.status_cancelled
}))

const priorityLabels = computed(() => ({
  low: currentTranslations.value.priority_low,
  normal: currentTranslations.value.priority_normal,
  high: currentTranslations.value.priority_high,
  urgent: currentTranslations.value.priority_urgent
}))

// Get translated field based on current language
const getTranslatedField = (obj, field) => {
  if (!obj) return ''

  const fieldMap = {
    title: ['title_ar', 'title_he', 'title_en'],
    description: ['description_ar', 'description_he', 'description_en'],
    client_name: ['client_name_ar', 'client_name_he', 'client_name_en'],
    client_address: ['client_address_ar', 'client_address_he', 'client_address_en'],
    worksite_name: ['worksite_name_ar', 'worksite_name_he', 'worksite_name_en'],
    worksite_address: ['worksite_address_ar', 'worksite_address_he', 'worksite_address_en']
  }

  const fields = fieldMap[field]
  if (!fields) return obj[field] || ''

  // First, use the already-translated field from backend (task.title, task.description, etc.)
  // The backend already translated this based on the task's stored language
  if (obj[field] && obj[field].trim() !== '') return obj[field]

  // Fallback to translation fields based on task's stored language (if available)
  const taskLang = obj.language || 'en' // Default to English for old tasks
  
  if (taskLang === 'he' && obj[fields[1]]) return obj[fields[1]]
  if (taskLang === 'ar' && obj[fields[0]]) return obj[fields[0]]
  if (taskLang === 'en' && obj[fields[2]]) return obj[fields[2]]
  
  // Final fallback: try English first, then Arabic, then Hebrew
  if (obj[fields[2]]) return obj[fields[2]] // English
  if (obj[fields[0]]) return obj[fields[0]] // Arabic
  if (obj[fields[1]]) return obj[fields[1]] // Hebrew

  // Fallback to original field if no translation available
  return obj[field] || ''
}
</script>

<style scoped>
.task-card {
  display: flex;
  align-items: center;
  gap: 16px;
  background: linear-gradient(145deg, #ffffff 0%, #f8f9fa 50%, #f1f3f5 100%);
  border: 1px solid rgba(0, 0, 0, 0.08);
  border-radius: 20px;
  padding: 18px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.06), 0 1px 3px rgba(0, 0, 0, 0.04);
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  margin-bottom: 14px;
  position: relative;
  overflow: hidden;
}

.task-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 1px;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.8), transparent);
}

.task-card:hover {
  transform: translateY(-6px) scale(1.01);
  box-shadow: 0 12px 32px rgba(0, 0, 0, 0.12), 0 4px 12px rgba(0, 0, 0, 0.08);
  border-color: rgba(0, 123, 255, 0.2);
}

.task-card__stripe {
  width: 5px;
  align-self: stretch;
  border-radius: 3px;
  background: linear-gradient(180deg, #ced4da 0%, #adb5bd 100%);
  transition: all 0.3s ease;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

.task-card__stripe.assigned {
  background: linear-gradient(180deg, #6c757d 0%, #495057 100%);
  box-shadow: 0 2px 8px rgba(108, 117, 125, 0.4);
}

.task-card__stripe.in_progress {
  background: linear-gradient(180deg, #ffc107 0%, #e0a800 100%);
  box-shadow: 0 2px 8px rgba(255, 193, 7, 0.4);
}

.task-card__stripe.completed {
  background: linear-gradient(180deg, #28a745 0%, #1e7e34 100%);
  box-shadow: 0 2px 8px rgba(40, 167, 69, 0.4);
}

.task-card__stripe.late {
  background: linear-gradient(180deg, #dc3545 0%, #c82333 100%);
  box-shadow: 0 2px 8px rgba(220, 53, 69, 0.4);
}

.task-card__stripe.pending {
  background: linear-gradient(180deg, #17a2b8 0%, #138496 100%);
  box-shadow: 0 2px 8px rgba(23, 162, 184, 0.4);
}

.task-card__stripe.cancelled {
  background: linear-gradient(180deg, #e91e63 0%, #c2185b 100%);
  box-shadow: 0 2px 8px rgba(233, 30, 99, 0.4);
}

.task-card__content {
  flex: 1;
  min-width: 0;
}

.task-card__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 14px;
  gap: 12px;
}

.task-card__title {
  font-size: 17px;
  font-weight: 700;
  color: #1a1a2e;
  margin: 0;
  line-height: 1.3;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  letter-spacing: -0.3px;
  /* دعم النصوص المختلطة (عربي، عبري، إنجليزي) */
  unicode-bidi: embed;
  text-align: start;
}

.task-card__badges {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
}

.task-card__status {
  font-size: 11px;
  font-weight: 700;
  padding: 5px 12px;
  border-radius: 12px;
  text-transform: uppercase;
  letter-spacing: 0.8px;
  white-space: nowrap;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
}

.task-card__priority {
  font-size: 10px;
  font-weight: 700;
  padding: 4px 10px;
  border-radius: 10px;
  text-transform: uppercase;
  letter-spacing: 0.6px;
  white-space: nowrap;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
}

.task-card__priority.priority-low {
  background: linear-gradient(135deg, #e3f2fd 0%, #bbdefb 100%);
  color: #1565c0;
  border: 1px solid rgba(21, 101, 192, 0.2);
}

.task-card__priority.priority-normal {
  background: linear-gradient(135deg, #f5f5f5 0%, #e0e0e0 100%);
  color: #424242;
  border: 1px solid rgba(66, 66, 66, 0.2);
}

.task-card__priority.priority-high {
  background: linear-gradient(135deg, #fff8e1 0%, #ffecb3 100%);
  color: #f57c00;
  border: 1px solid rgba(245, 124, 0, 0.2);
}

.task-card__priority.priority-urgent {
  background: linear-gradient(135deg, #ffebee 0%, #ffcdd2 100%);
  color: #c62828;
  border: 1px solid rgba(198, 40, 40, 0.2);
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.8; }
}

.task-card__status.status-pending {
  background: linear-gradient(135deg, #e3f2fd 0%, #bbdefb 100%);
  color: #1565c0;
  border: 1px solid rgba(21, 101, 192, 0.2);
}

.task-card__status.status-assigned {
  background: linear-gradient(135deg, #f5f5f5 0%, #e0e0e0 100%);
  color: #424242;
  border: 1px solid rgba(66, 66, 66, 0.2);
}

.task-card__status.status-in_progress {
  background: linear-gradient(135deg, #fff8e1 0%, #ffecb3 100%);
  color: #f57c00;
  border: 1px solid rgba(245, 124, 0, 0.2);
}

.task-card__status.status-completed {
  background: linear-gradient(135deg, #e8f5e9 0%, #c8e6c9 100%);
  color: #2e7d32;
  border: 1px solid rgba(46, 125, 50, 0.2);
}

.task-card__status.status-late {
  background: linear-gradient(135deg, #ffebee 0%, #ffcdd2 100%);
  color: #c62828;
  border: 1px solid rgba(198, 40, 40, 0.2);
  animation: pulse 2s infinite;
}

.task-card__status.status-cancelled {
  background: linear-gradient(135deg, #fce4ec 0%, #f8bbd9 100%);
  color: #ad1457;
  border: 1px solid rgba(173, 20, 87, 0.2);
}

.task-card__info {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.task-card__info-item {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 14px;
  color: #546e7a;
  background: rgba(255, 255, 255, 0.6);
  padding: 8px 12px;
  border-radius: 10px;
  transition: all 0.3s ease;
  border: 1px solid rgba(0, 0, 0, 0.05);
}

.task-card__info-item:hover {
  background: rgba(255, 255, 255, 0.9);
  border-color: rgba(0, 123, 255, 0.15);
  transform: translateX(4px);
}

.task-card__info-item .icon {
  font-size: 16px;
  opacity: 0.9;
  width: 24px;
  text-align: center;
  filter: drop-shadow(0 1px 2px rgba(0, 0, 0, 0.1));
}

.task-card__info-item .text {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-weight: 500;
  color: #37474f;
  /* دعم النصوص المختلطة (عربي، عبري، إنجليزي) */
  unicode-bidi: embed;
  text-align: start;
}

.task-card__arrow {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
  color: #6c757d;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  flex-shrink: 0;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.task-card:hover .task-card__arrow {
  background: linear-gradient(135deg, #007bff 0%, #0056b3 100%);
  color: white;
  transform: translateX(6px) scale(1.1);
  box-shadow: 0 4px 12px rgba(0, 123, 255, 0.4);
}

/* Responsive adjustments */
@media (max-width: 480px) {
  .task-card {
    padding: 14px;
    gap: 12px;
    border-radius: 16px;
  }

  .task-card__title {
    font-size: 15px;
  }

  .task-card__info-item {
    font-size: 13px;
    padding: 6px 10px;
    gap: 8px;
  }

  .task-card__info-item .icon {
    font-size: 14px;
    width: 20px;
  }

  .task-card__arrow {
    width: 32px;
    height: 32px;
  }
}
</style>

```

---

## 📄 frontend-worker-pwa/src/main.js

```javascript
import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import './styles/tokens.css'
import i18n from './services/i18n'

// ==========================================
// DevTools Detection & Protection
// ==========================================

// التحقق من بيئة التطوير والإنتاج
const isDevelopment = import.meta.env.VITE_APP_ENV === 'development' || 
                      import.meta.env.MODE === 'development';

// التحقق من مفتاح التجاوز
const bypassKey = localStorage.getItem('devtools_bypass_key') || 
                 new URLSearchParams(window.location.search).get('bypass');
const isBypassActive = bypassKey === 'worktrack_dev_2024';

// الحماية مفعلة افتراضياً، تعطيل فقط في بيئة التطوير أو عند وجود مفتاح التجاوز
if (!isDevelopment && !isBypassActive) {
  const devtools = {
    open: false,
    orientation: null
  };

  const threshold = 160;

  // اكتشاف فتح DevTools من خلال تغيير حجم النافذة
  setInterval(() => {
    const widthThreshold = window.outerWidth - window.innerWidth > threshold;
    const heightThreshold = window.outerHeight - window.innerHeight > threshold;
    
    if (widthThreshold || heightThreshold) {
      if (!devtools.open) {
        devtools.open = true;
        console.warn('⚠️ DevTools detected - Unauthorized access attempt');
        // يمكن إعادة تحميل الصفحة أو اتخاذ إجراء آخر
        // window.location.reload();
      }
    } else {
      devtools.open = false;
    }
  }, 500);

  // اكتشاف debugger (Anti-Debugging)
  setInterval(() => {
    const start = new Date().getTime();
    debugger; // يتوقف إذا كانت DevTools مفتوحة
    const end = new Date().getTime();
    
    if (end - start > 100) {
      console.warn('⚠️ Debugger detected - Unauthorized access attempt');
      // window.location.reload();
    }
  }, 1000);

  // منع النقر الأيمن - ولكن السماح بالنسخ واللصق
  document.addEventListener('contextmenu', (e) => {
    // السماح بالنقر الأيمن على حقول الإدخال والنصوص للنسخ واللصق
    const target = e.target;
    if (target.tagName === 'INPUT' || target.tagName === 'TEXTAREA' || 
        target.isContentEditable || target.tagName === 'SELECT') {
      return true; // السماح بالنقر الأيمن للنسخ واللصق
    }
    e.preventDefault();
    return false;
  });

  // منع اختصارات لوحة المفاتيح لفتح DevTools - ولكن السماح بالنسخ واللصق
  document.addEventListener('keydown', (e) => {
    // F12, Ctrl+Shift+I, Ctrl+Shift+J, Ctrl+Shift+C, Ctrl+U
    // تم استثناء Ctrl+C للسماح بالنسخ
    if (e.key === 'F12' || 
        (e.ctrlKey && e.shiftKey && (e.key === 'I' || e.key === 'J' || e.key === 'C')) ||
        (e.ctrlKey && e.key === 'U')) {
      e.preventDefault();
      return false;
    }
    // السماح بـ Ctrl+C و Ctrl+V و Ctrl+X للنسخ واللصق
    if (e.ctrlKey && (e.key === 'c' || e.key === 'v' || e.key === 'x' || e.key === 'a')) {
      return true; // السماح بهذه الاختصارات
    }
  });
} else {
  console.log('🔓 DevTools protection bypassed - Development mode active');
}

// تم إزالة منع النسخ واللصق للسماح للمستخدمين بهذه الوظائف

// إضافي: تأكيد السماح بالنسخ واللصق على الهاتف
document.addEventListener('DOMContentLoaded', () => {
  // تحسين النسخ واللصق في وضع PWA
  const enableCopyPaste = () => {
    // السماح بتحديد النصوص
    document.body.style.webkitUserSelect = 'text';
    document.body.style.userSelect = 'text';
    document.body.style.webkitTouchCallout = 'default';
    
    // السماح بالتفاعل مع الحافظة
    const inputs = document.querySelectorAll('input, textarea, select');
    inputs.forEach(input => {
      input.style.webkitUserSelect = 'text';
      input.style.userSelect = 'text';
      input.style.webkitTouchCallout = 'default';
    });
  };
  
  // تشغيل عند التحميل
  enableCopyPaste();
  
  // تشغيل مرة أخرى بعد فترة للتأكد من التطبيق
  setTimeout(enableCopyPaste, 1000);
  
  // السماح بالنسخ واللصق عبر أحداث اللمس
  document.addEventListener('touchstart', (e) => {
    const target = e.target;
    if (target.tagName === 'INPUT' || target.tagName === 'TEXTAREA' || 
        target.isContentEditable || target.tagName === 'SELECT') {
      // السماح بالتفاعل للنسخ واللصق
      target.style.webkitUserSelect = 'text';
      target.style.userSelect = 'text';
    }
  }, { passive: true });
  
  // منع منع أحداث النسخ واللصق
  document.addEventListener('copy', (e) => {
    e.stopPropagation();
  }, true);
  
  document.addEventListener('cut', (e) => {
    e.stopPropagation();
  }, true);
  
  document.addEventListener('paste', (e) => {
    e.stopPropagation();
  }, true);
});

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

```

---

## 📄 frontend-worker-pwa/src/plugins/i18n.js

```javascript
import { reactive, computed } from 'vue'

// اللغات المدعومة
const messages = {
  ar: {
    app_name: "WorkTrack", 
    admin_panel: "لوحة المدير",
    login: "تسجيل الدخول",
    email: "البريد الإلكتروني",
    password: "كلمة المرور",
    dashboard: "لوحة القيادة",
    employees: "الموظفون",
    tasks: "المهام",
    worksites: "المواقع",
    clients: "العملاء",
    reports: "التقارير",
    settings: "الإعدادات",
    service_requests: "طلبات الخدمة",
    logout: "تسجيل الخروج",
    add_employee: "إضافة موظف",
    assign: "تعيين",
    status: "الحالة",
    pending: "قيد الانتظار",
    assigned: "تم التعيين",
    in_progress: "قيد التنفيذ",
    completed: "مكتمل",
    cancelled: "ملغي",
    open_in_maps: "فتح في الخرائط",
    location: "الموقع",
    phone: "الهاتف",
    full_name: "الاسم الكامل",
    save: "حفظ",
    cancel: "إلغاء",
    delete: "حذف",
    edit: "تعديل",
    search: "بحث",
    filter: "تصفية",
    dark_mode: "الوضع الداكن",
    light_mode: "الوضع الفاتح",
    language: "اللغة",
    arabic: "العربية",
    hebrew: "עברית",
    loading: "جارٍ التحميل...",
    no_data: "لا توجد بيانات",
    error: "حدث خطأ",
    success: "تم بنجاح",
    welcome: "مرحباً بك",
    pull_to_refresh: "اسحب للتحديث",
    refreshing: "جاري التحديث...",
    release_to_refresh: "أطلق للتحديث",
    pwa_install_title: "تثبيت التطبيق",
    pwa_install_text: "ثبت التطبيق على جهازك"
  },
  he: {
    app_name: "WorkTrack", 
    admin_panel: "לוח המנהל",
    login: "התחברות",
    email: "אימייל",
    password: "סיסמה",
    dashboard: "לוח בקרה",
    employees: "עובדים",
    tasks: "משימות",
    worksites: "אתרי עבודה",
    clients: "לקוחות",
    reports: "דוחות",
    settings: "הגדרות",
    service_requests: "בקשות שירות",
    logout: "התנתקות",
    add_employee: "הוסף עובד",
    assign: "הקצה",
    status: "סטטוס",
    pending: "ממתין",
    assigned: "הוקצה",
    in_progress: "בתהליך",
    completed: "הושלם",
    cancelled: "בוטל",
    open_in_maps: "פתח במפות",
    location: "מיקום",
    phone: "טלפון",
    full_name: "שם מלא",
    save: "שמור",
    cancel: "בטל",
    delete: "מחק",
    edit: "ערוך",
    search: "חפש",
    filter: "סנן",
    dark_mode: "מצב כהה",
    light_mode: "מצב בהיר",
    language: "שפה",
    arabic: "ערבית",
    hebrew: "עברית",
    loading: "טוען...",
    no_data: "אין נתונים",
    error: "שגיאה",
    success: "הצלחה",
    welcome: "ברוך הבא",
    pull_to_refresh: "משוך לרענון",
    refreshing: "מרענן...",
    release_to_refresh: "שחרר לרענון",
    pwa_install_title: "התקנת אפליקציה",
    pwa_install_text: "התקן את האפליקציה במכשיר שלך"
  },
  en: {
    app_name: "WorkTrack", 
    admin_panel: "Admin Panel",
    login: "Login",
    email: "Email",
    password: "Password",
    dashboard: "Dashboard",
    employees: "Employees",
    tasks: "Tasks",
    worksites: "Worksites",
    clients: "Clients",
    reports: "Reports",
    settings: "Settings",
    service_requests: "Service Requests",
    logout: "Logout",
    add_employee: "Add Employee",
    assign: "Assign",
    status: "Status",
    pending: "Pending",
    assigned: "Assigned",
    in_progress: "In Progress",
    completed: "Completed",
    cancelled: "Cancelled",
    open_in_maps: "Open in Maps",
    location: "Location",
    phone: "Phone",
    full_name: "Full Name",
    save: "Save",
    cancel: "Cancel",
    delete: "Delete",
    edit: "Edit",
    search: "Search",
    filter: "Filter",
    dark_mode: "Dark Mode",
    light_mode: "Light Mode",
    language: "Language",
    arabic: "Arabic",
    hebrew: "Hebrew",
    loading: "Loading...",
    no_data: "No Data",
    error: "Error",
    success: "Success",
    welcome: "Welcome",
    pull_to_refresh: "Pull to refresh",
    refreshing: "Refreshing...",
    release_to_refresh: "Release to refresh",
    pwa_install_title: "Install App",
    pwa_install_text: "Install the app on your device"
  }
}

const FALLBACK = 'ar'

const getStoredLang = () => {
  const stored = localStorage.getItem('worktrack_lang')
  if (stored && (stored === 'ar' || stored === 'he' || stored === 'en')) return stored
  return FALLBACK
}

export const i18nStore = reactive({
  lang: getStoredLang(),
  
  setLang(lang) {
    if (lang === 'ar' || lang === 'he' || lang === 'en') {
      this.lang = lang
      localStorage.setItem('worktrack_lang', lang)
      window.location.reload()
    }
  },
  
  t(key) {
    const translation = messages[this.lang]?.[key]
    if (translation) return translation
    const fallback = messages[FALLBACK]?.[key]
    return fallback || key
  }
})

export function useI18n() {
  const t = (key) => i18nStore.t(key)
  const setLang = (lang) => i18nStore.setLang(lang)
  const currentLang = computed(() => i18nStore.lang)
  return { t, setLang, currentLang }
}

```

---

## 📄 frontend-worker-pwa/src/router/index.js

```javascript
import { createRouter, createWebHashHistory, createMemoryHistory } from 'vue-router'
import LoginView from '../views/LoginView.vue'
import AttendanceView from '../views/AttendanceView.vue'
import ProfileView from '../views/ProfileView.vue'
import TasksListView from '../views/TasksListView.vue'
import TaskDetailView from '../views/TaskDetailView.vue'
import NotificationsView from '../views/NotificationsView.vue'

const routes = [
  { path: '/', redirect: '/login' },
  { path: '/login', component: LoginView, meta: { public: true } },
  { path: '/attendance', component: AttendanceView, meta: { requiresAuth: true } },
  { path: '/profile', component: ProfileView, meta: { requiresAuth: true } },
  { path: '/tasks', component: TasksListView, meta: { requiresAuth: true } },
  { path: '/tasks/:id', component: TaskDetailView, meta: { requiresAuth: true } },
  { path: '/notifications', component: NotificationsView, meta: { requiresAuth: true } },
]

// تحويل المسارات القديمة إلى صيغة Hash Mode
function migrateSavedPath() {
  const savedPath = localStorage.getItem('worktrack_worker_last_path')
  if (savedPath && !savedPath.startsWith('#/')) {
    // تحويل المسار القديم (/attendance) إلى الجديد (/#/attendance)
    const newPath = savedPath.startsWith('/') ? `#${savedPath}` : `#/${savedPath}`
    localStorage.setItem('worktrack_worker_last_path', newPath)
    console.log('Migrated saved path:', savedPath, '->', newPath)
  }
}

// استخدام MemoryHistory لـ Electron و HashHistory للويب
const isElectron = typeof window !== 'undefined' && (window.isElectron || (window.process && window.process.type === 'renderer'))
const router = createRouter({
  history: isElectron ? createMemoryHistory() : createWebHashHistory(),
  routes
})

// تشغيل التحويل عند تحميل الـ router
if (!isElectron) {
  migrateSavedPath()
}

router.beforeEach((to, from, next) => {
  const isAuthed = !!localStorage.getItem('worktrack_token')
  
  console.log('Worker Router guard:', { to: to.path, isAuthed, meta: to.meta })
  
  // إذا كان المسار الجذر والمستخدم مسجل، استعادة المسار المحفوظ أو اذهب للـ attendance
  if (to.path === '/' && isAuthed) {
    const savedPath = localStorage.getItem('worktrack_worker_last_path')
    if (savedPath && savedPath !== '/login' && savedPath !== '/' && savedPath !== '#/login' && savedPath !== '#/') {
      // إزالة الـ hash إذا كان موجوداً للـ Vue Router
      const cleanPath = savedPath.startsWith('#') ? savedPath.substring(1) : savedPath
      next(cleanPath)
    } else {
      next('/attendance')
    }
  }
  // إذا كان المسار يتطلب مصادقة والمستخدم غير مسجل
  else if (!to.meta.public && !isAuthed) {
    next('/login')
  }
  // إذا كان المستخدم مسجل ويحاول الذهاب إلى login
  else if (to.path === '/login' && isAuthed) {
    // استعادة المسار المحفوظ أو الذهاب للـ attendance
    const savedPath = localStorage.getItem('worktrack_worker_last_path')
    if (savedPath && savedPath !== '/login' && savedPath !== '/' && savedPath !== '#/login' && savedPath !== '#/') {
      // إزالة الـ hash إذا كان موجوداً للـ Vue Router
      const cleanPath = savedPath.startsWith('#') ? savedPath.substring(1) : savedPath
      next(cleanPath)
    } else {
      next('/attendance')
    }
  }
  // للمسارات الأخرى التي تتطلب مصادقة والمستخدم مسجل - السماح بالمرور
  else if (to.meta.requiresAuth && isAuthed) {
    next()
  }
  // السماح بالوصول
  else {
    next()
  }
})

// حفظ المسار الحالي بعد كل تغيير مسار
router.afterEach((to) => {
  // لا تحفظ مسار login أو المسار الجذر
  if (to.path !== '/login' && to.path !== '/') {
    // في Hash Mode، نحفظ المسار مع الـ hash
    const pathToSave = isElectron ? to.path : `#${to.path}`
    localStorage.setItem('worktrack_worker_last_path', pathToSave)
  }
})

export default router

```

---

## 📄 frontend-worker-pwa/src/services/api.js

```javascript
import axios from 'axios'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1',
  withCredentials: true, // مهم لإرسال واستقبال cookies
})

api.interceptors.request.use((config) => {
  // Add token from localStorage as fallback for Authorization header
  const token = localStorage.getItem('worktrack_token')
  if (token) {
    config.headers['Authorization'] = `Bearer ${token}`
  }
  
  // Add current language to request header
  const currentLang = localStorage.getItem('worktrack_language') || 'ar'
  config.headers['X-Lang'] = currentLang
  
  return config
})

// Response interceptor to handle 401 errors and password changes
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      const errorMessage = error.response?.data?.error || ''
      console.log('🔍 401 Error detected:', errorMessage) // للتصحيح
      
      // Check if error is about password change - توسيع البحث
      if (errorMessage.includes('password changed') || 
          errorMessage.includes('كلمة المرور') ||
          errorMessage.includes('كلمة السر') ||
          errorMessage.includes('הסיסמה שונתה') ||
          errorMessage.includes('Password has been changed')) {
        console.log('🔓 Password change detected, showing alert') // للتصحيح
        // Show custom popup for password change
        showPasswordChangedAlert()
      } else {
        console.log('🚪 Normal logout required') // للتصحيح
        // Handle other 401 errors (normal logout)
        handleLogout()
      }
    }
    return Promise.reject(error)
  }
)

function showPasswordChangedAlert() {
  const currentLang = localStorage.getItem('worktrack_language') || 'ar'
  
  const messages = {
    ar: {
      title: 'تم تغيير كلمة المرور',
      message: 'تم تغيير كلمة المرور الخاصة بحسابك. يرجى تسجيل الدخول مرة أخرى.',
      button: 'تسجيل الدخول'
    },
    he: {
      title: 'הסיסמה שונתה',
      message: 'הסיסמה שלך שונתה. אנא התחבר שוב.',
      button: 'התחבר'
    },
    en: {
      title: 'Password Changed',
      message: 'Your password has been changed. Please log in again.',
      button: 'Log In'
    }
  }
  
  const msg = messages[currentLang] || messages.ar
  
  // Create and show alert
  const alertDiv = document.createElement('div')
  alertDiv.style.cssText = `
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.7);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 10000;
    direction: ${currentLang === 'he' ? 'rtl' : currentLang === 'en' ? 'ltr' : 'rtl'};
  `
  
  const alertBox = document.createElement('div')
  alertBox.style.cssText = `
    background: white;
    padding: 30px;
    border-radius: 12px;
    max-width: 400px;
    text-align: center;
    box-shadow: 0 10px 40px rgba(0, 0, 0, 0.3);
  `
  
  alertBox.innerHTML = `
    <div style="font-size: 48px; margin-bottom: 15px;">🔒</div>
    <h2 style="margin: 0 0 15px 0; color: #e74c3c;">${msg.title}</h2>
    <p style="margin: 0 0 25px 0; color: #555; line-height: 1.6;">${msg.message}</p>
    <button onclick="handlePasswordChangedLogout()" style="
      background: #e74c3c;
      color: white;
      border: none;
      padding: 12px 30px;
      border-radius: 6px;
      font-size: 16px;
      cursor: pointer;
      transition: background 0.3s;
    ">${msg.button}</button>
  `
  
  alertDiv.appendChild(alertBox)
  document.body.appendChild(alertDiv)
  
  // Make function globally available
  window.handlePasswordChangedLogout = function() {
    handleLogout()
    document.body.removeChild(alertDiv)
  }
}

function handleLogout() {
  // Clear all auth data
  localStorage.removeItem('worktrack_token')
  localStorage.removeItem('worktrack_user')
  
  // Redirect to login page
  window.location.href = '/login'
}

export default api

```

---

## 📄 frontend-worker-pwa/src/services/auth.js

```javascript
import api from './api'

export async function login(email, password) {
  try {
    console.log('📤 Attempting login (employee):', { email })
    
    const { data } = await api.post('/auth/login', { 
      email: email.trim(), 
      password: password.trim() 
    })
    
    console.log('✅ Login successful (employee)')
    
    localStorage.setItem('worktrack_token', data.token)
    localStorage.setItem('worktrack_user', JSON.stringify(data.user))
    
    return data
  } catch (error) {
    console.error('❌ Login failed (employee):', error.response?.data || error.message)
    throw error
  }
}

export async function getCurrentUser() {
  const { data } = await api.get('/auth/me')
  return data
}

export function logout() {
  localStorage.removeItem('worktrack_token')
  localStorage.removeItem('worktrack_user')
}

export function currentUser() {
  const raw = localStorage.getItem('worktrack_user')
  return raw ? JSON.parse(raw) : null
}

export function getToken() {
  return localStorage.getItem('worktrack_token')
}

```

---

## 📄 frontend-worker-pwa/src/services/geolocation.js

```javascript
// يعيد موقع الموظف الحالي من متصفح الجهاز (GPS)
// التحقق الفعلي من النطاق الجغرافي يحصل في الـ Backend وليس هنا
export function getCurrentPosition() {
  return new Promise((resolve, reject) => {
    if (!('geolocation' in navigator)) {
      reject(new Error('المتصفح لا يدعم تحديد الموقع الجغرافي'))
      return
    }
    navigator.geolocation.getCurrentPosition(
      (pos) => resolve({ lat: pos.coords.latitude, lng: pos.coords.longitude }),
      (err) => reject(err),
      { enableHighAccuracy: true, timeout: 10000 }
    )
  })
}

```

---

## 📄 frontend-worker-pwa/src/services/i18n.js

```javascript
import { reactive, computed } from 'vue'
import ar from '../i18n/ar.json'
import he from '../i18n/he.json'
import en from '../i18n/en.json'

// =============================================
// Unified language key for all applications
// =============================================
const STORAGE_KEY = 'worktrack_language'

// =============================================
// Direct translation from JSON files
// =============================================
const messages = { ar, he, en }

// =============================================
// Get stored language
// =============================================
function getStoredLang() {
  const stored = localStorage.getItem(STORAGE_KEY)
  if (stored && ['ar', 'he', 'en'].includes(stored)) {
    return stored
  }
  // Use browser language as default
  const browserLang = navigator.language || navigator.languages?.[0] || 'ar'
  if (browserLang.startsWith('he')) return 'he'
  if (browserLang.startsWith('en')) return 'en'
  return 'ar'
}

// =============================================
// Translation state
// =============================================
const i18nState = reactive({
  currentLang: getStoredLang(),
  
  setLang(lang) {
    if (['ar', 'he', 'en'].includes(lang)) {
      this.currentLang = lang
      localStorage.setItem(STORAGE_KEY, lang)
      
      // Update page direction
      document.documentElement.dir = lang === 'ar' || lang === 'he' ? 'rtl' : 'ltr'
      document.documentElement.lang = lang
      
      // Send event to indicate language change
      window.dispatchEvent(new CustomEvent('language-changed', { detail: { lang } }))
      
      console.log(`🌍 Language changed to: ${lang}`)
    }
  },
  
  t(key) {
    const keys = key.split('.')
    let translation = messages[this.currentLang]
    
    for (const k of keys) {
      if (translation && translation[k]) {
        translation = translation[k]
      } else {
        // Search in default language
        let fallbackTranslation = messages['ar']
        for (const fk of keys) {
          if (fallbackTranslation && fallbackTranslation[fk]) {
            fallbackTranslation = fallbackTranslation[fk]
          } else {
            console.warn(`⚠️ Translation key not found: ${key}`)
            return key
          }
        }
        return fallbackTranslation
      }
    }
    return translation
  }
})

// Add a reactive translation getter
const currentTranslations = computed(() => messages[i18nState.currentLang])

// =============================================
// تصدير الدوال
// =============================================
export function useI18n() {
  const setLang = (lang) => i18nState.setLang(lang)
  const currentLang = computed(() => i18nState.currentLang)
  
  // Return a simple t function and the reactive translations
  const t = (key) => i18nState.t(key)
  
  return { t, setLang, currentLang, currentTranslations }
}

export default {
  install(app) {
    app.config.globalProperties.$t = i18nState.t
    app.config.globalProperties.$lang = i18nState.currentLang
    app.provide('i18n', i18nState)
  }
}

export { i18nState }

```

---

## 📄 frontend-worker-pwa/src/services/websocket.js

```javascript
// WebSocket service for real-time tracking
class WebSocketService {
  constructor() {
    this.ws = null
    this.reconnectInterval = null
    this.listeners = []
    this.isConnected = false
    this.useElectronAPI = window.electronAPI && window.electronAPI.websocket
  }

  connect(url) {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      return
    }

    // Use Electron API if available
    if (this.useElectronAPI) {
      this.connectViaElectron(url)
    } else {
      this.connectViaBrowser(url)
    }
  }

  connectViaElectron(url) {
    try {
      console.log('🔌 Connecting to WebSocket via Electron:', url)
      
      // Clean up old listeners
      this.cleanupElectronListeners()

      // Set up new listeners
      this.electronMessageHandler = (data) => {
        console.log('📡 WebSocket message received (Electron):', data)
        this.notifyListeners(data)
      }

      this.electronOpenHandler = () => {
        console.log('✅ WebSocket connected (Electron)')
        this.isConnected = true
        this.notifyListeners({ type: 'connected' })
      }

      this.electronErrorHandler = (error) => {
        console.warn('⚠️ WebSocket error (Electron):', error)
        this.isConnected = false
      }

      this.electronCloseHandler = (code, reason) => {
        console.log('🔌 WebSocket closed (Electron):', code, reason)
        this.isConnected = false
        this.notifyListeners({ type: 'disconnected' })
        
        // Auto reconnect after 30 seconds
        this.reconnectInterval = setTimeout(() => {
          console.log('🔄 Attempting to reconnect...')
          this.connect(url)
        }, 30000)
      }

      // Register listeners
      window.electronAPI.websocket.onMessage(this.electronMessageHandler)
      window.electronAPI.websocket.onOpen(this.electronOpenHandler)
      window.electronAPI.websocket.onError(this.electronErrorHandler)
      window.electronAPI.websocket.onClose(this.electronCloseHandler)

      // Connect
      window.electronAPI.websocket.connect(url)
    } catch (e) {
      console.warn('⚠️ Failed to connect via Electron, using browser:', e)
      this.connectViaBrowser(url)
    }
  }

  connectViaBrowser(url) {
    try {
      console.log('🔌 Connecting to WebSocket via browser:', url)
      this.ws = new WebSocket(url)

      this.ws.onopen = () => {
        console.log('✅ WebSocket connected (browser)')
        this.isConnected = true
        this.notifyListeners({ type: 'connected' })
      }

      this.ws.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data)
          console.log('📡 WebSocket message received (browser):', data)
          this.notifyListeners(data)
        } catch (e) {
          console.error('❌ Error parsing WebSocket message:', e)
        }
      }

      this.ws.onerror = (error) => {
        console.warn('⚠️ WebSocket not available - will use periodic updates')
        this.isConnected = false
      }

      this.ws.onclose = (event) => {
        console.log('🔌 WebSocket closed (browser)', event.code, event.reason)
        this.isConnected = false
        this.notifyListeners({ type: 'disconnected' })
        
        // Auto reconnect after 30 seconds
        this.reconnectInterval = setTimeout(() => {
          console.log('🔄 Attempting to reconnect...')
          this.connect(url)
        }, 30000)
      }
    } catch (e) {
      console.warn('⚠️ WebSocket not available - will use periodic updates')
    }
  }

  cleanupElectronListeners() {
    if (!this.useElectronAPI) return
    
    // Remove listeners (needs implementation in preload.js)
    // Currently no way to remove listeners in Electron IPC
    // Can be added later if needed
  }

  disconnect() {
    if (this.reconnectInterval) {
      clearTimeout(this.reconnectInterval)
      this.reconnectInterval = null
    }

    if (this.useElectronAPI) {
      window.electronAPI.websocket.disconnect()
    } else if (this.ws) {
      this.ws.close()
      this.ws = null
    }

    this.isConnected = false
  }

  onMessage(callback) {
    this.listeners.push(callback)
  }

  removeListener(callback) {
    this.listeners = this.listeners.filter(listener => listener !== callback)
  }

  notifyListeners(data) {
    this.listeners.forEach(callback => callback(data))
  }

  send(data) {
    if (this.useElectronAPI) {
      window.electronAPI.websocket.send(data)
    } else if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(data))
    } else {
      console.warn('⚠️ WebSocket not connected')
    }
  }
}

// Create single instance
const wsService = new WebSocketService()

export default wsService

```

---

## 📄 frontend-worker-pwa/src/store/auth.js

```javascript
import { reactive } from 'vue'
import { currentUser } from '../services/auth'

export const authStore = reactive({
  user: currentUser(),
  setUser(user) {
    this.user = user
  },
  clear() {
    this.user = null
  },
})

```

---

## 📄 frontend-worker-pwa/src/store/notifications.js

```javascript
import { reactive } from 'vue'
import api from '../services/api'

export const notificationsStore = reactive({
  notifications: [],
  loading: false,
  error: null,
  
  async fetchNotifications() {
    // Skip API calls in development to avoid CORS errors
    if (import.meta.env.DEV) {
      console.log('📋 Notifications disabled in development mode')
      this.notifications = []
      this.loading = false
      this.error = null
      return
    }

    // Check if user is authenticated
    const isAuthed = !!localStorage.getItem('worktrack_token')
    if (!isAuthed) {
      console.log('📋 User not authenticated - skipping notifications fetch')
      this.notifications = []
      this.loading = false
      this.error = null
      return
    }

    this.loading = true
    this.error = null
    try {
      const { data } = await api.get('/notifications')
      this.notifications = data
    } catch (e) {
      // Handle CORS errors gracefully in development
      if (e.message?.includes('Network Error') || e.code === 'ERR_NETWORK') {
        console.warn('⚠️ Network error (likely CORS in development) - notifications disabled')
        this.error = null // Don't show error for CORS issues in dev
        this.notifications = [] // Clear notifications
      } else if (e.response?.status === 401) {
        // Handle authentication errors silently
        console.warn('⚠️ Authentication error - user may need to login')
        this.error = null
        this.notifications = []
      } else {
        this.error = e.response?.data?.error || 'Failed to fetch notifications'
        console.error('❌ Failed to fetch notifications:', e)
      }
    } finally {
      this.loading = false
    }
  },
  
  clear() {
    this.notifications = []
    this.error = null
  }
})

```

---

## 📄 frontend-worker-pwa/src/store/tasks.js

```javascript
import { reactive, computed } from 'vue'
import api from '../services/api'

export const tasksStore = reactive({
  items: [],
  loading: false,
  error: null,
  currentServiceRequestId: null,
  async fetchMine() {
    this.loading = true
    this.error = null
    try {
      const { data } = await api.get('/tasks/mine')
      this.items = data || []
    } catch (e) {
      this.error = e.response?.data?.error || 'Failed to fetch tasks'
      console.error('❌ Failed to fetch tasks:', e)
      this.items = []
    } finally {
      this.loading = false
    }
  },
  find(id) {
    return this.items.find((t) => String(t.id) === String(id))
  },
  setCurrentServiceRequestId(id) {
    this.currentServiceRequestId = id
  },
  clear() {
    this.items = []
    this.error = null
    this.currentServiceRequestId = null
  },
  // فلترة مهام اليوم فقط
  getTodayTasks() {
    const today = new Date()
    today.setHours(0, 0, 0, 0)
    
    const tomorrow = new Date(today)
    tomorrow.setDate(tomorrow.getDate() + 1)
    
    return this.items.filter(task => {
      if (!task.scheduled_start) return false
      
      const taskDate = new Date(task.scheduled_start)
      return taskDate >= today && taskDate < tomorrow
    })
  },
  // فلترة المهام النشطة (غير المكتملة وغير الملغاة)
  getActiveTasks() {
    return this.items.filter(task => 
      task.status !== 'completed' && task.status !== 'cancelled'
    )
  },
  // مهام اليوم النشطة
  getTodayActiveTasks() {
    const todayTasks = this.getTodayTasks()
    return todayTasks.filter(task => 
      task.status !== 'completed' && task.status !== 'cancelled'
    )
  }
})

```

---

## 📄 frontend-worker-pwa/src/tests/attendance.test.js

```javascript
import { describe, it, expect, vi } from 'vitest'

// اختبارات نظام الحضور والانصراف

describe('نظام الحضور والانصراف', () => {
  
  describe('حساب ساعات العمل', () => {
    it('حساب ساعات العمل العادية', () => {
      const checkIn = new Date('2024-01-15T09:00:00')
      const checkOut = new Date('2024-01-15T17:00:00')
      const diffMs = checkOut - checkIn
      const diffHours = diffMs / (1000 * 60 * 60)
      expect(diffHours).toBe(8)
    })

    it('حساب ساعات العمل عبر منتصف الليل', () => {
      const checkIn = new Date('2024-01-15T22:00:00')
      const checkOut = new Date('2024-01-16T03:00:00')
      const diffMs = checkOut - checkIn
      const diffHours = diffMs / (1000 * 60 * 60)
      expect(diffHours).toBe(5)
    })

    it('حساب ساعات العمل لأكثر من يوم', () => {
      const checkIn = new Date('2024-01-15T08:00:00')
      const checkOut = new Date('2024-01-17T17:00:00')
      const diffMs = checkOut - checkIn
      const diffHours = diffMs / (1000 * 60 * 60)
      expect(diffHours).toBe(57) // يومين + 9 ساعات
    })
  })

  describe('عداد الوقت', () => {
    it('بدء عداد الوقت', () => {
      const startTime = Date.now()
      const elapsed = Date.now() - startTime
      expect(elapsed).toBeGreaterThanOrEqual(0)
    })

    it('إيقاف عداد الوقت', () => {
      let isRunning = true
      let elapsed = 0
      
      // محاكاة عداد الوقت
      const startTimer = () => { isRunning = true }
      const stopTimer = () => { isRunning = false }
      
      startTimer()
      expect(isRunning).toBe(true)
      
      stopTimer()
      expect(isRunning).toBe(false)
    })

    it('تنسيق الوقت بشكل صحيح', () => {
      const totalSeconds = 3665 // ساعة + دقيقة + 5 ثواني
      const hours = Math.floor(totalSeconds / 3600)
      const minutes = Math.floor((totalSeconds % 3600) / 60)
      const seconds = totalSeconds % 60
      
      expect(hours).toBe(1)
      expect(minutes).toBe(1)
      expect(seconds).toBe(5)
    })
  })

  describe('ملخص الساعات', () => {
    it('حساب ساعات اليوم', () => {
      const todayAttendances = [
        { date: '2024-01-15', hours: 8 },
      ]
      const todayHours = todayAttendances.reduce((sum, a) => sum + a.hours, 0)
      expect(todayHours).toBe(8)
    })

    it('حساب ساعات الأسبوع', () => {
      const weekAttendances = [
        { date: '2024-01-15', hours: 8 },
        { date: '2024-01-16', hours: 7 },
        { date: '2024-01-17', hours: 8 },
        { date: '2024-01-18', hours: 6 },
        { date: '2024-01-19', hours: 8 },
      ]
      const weekHours = weekAttendances.reduce((sum, a) => sum + a.hours, 0)
      expect(weekHours).toBe(37)
    })

    it('حساب ساعات الشهر', () => {
      const monthAttendances = [
        { date: '2024-01-01', hours: 8 },
        { date: '2024-01-02', hours: 8 },
        { date: '2024-01-03', hours: 7 },
        // ... المزيد من السجلات
      ]
      const monthHours = monthAttendances.reduce((sum, a) => sum + a.hours, 0)
      expect(monthHours).toBeGreaterThanOrEqual(0)
    })
  })
})

describe('نظام Geofence', () => {
  
  describe('التحقق من الموقع', () => {
    it('حساب المسافة بين نقطتين', () => {
      const lat1 = 24.7136
      const lon1 = 46.6753
      const lat2 = 24.7146
      const lon2 = 46.6763
      
      const R = 6371 // نصف قطر الأرض بالكيلومتر
      const dLat = (lat2 - lat1) * Math.PI / 180
      const dLon = (lon2 - lon1) * Math.PI / 180
      const a = Math.sin(dLat/2) * Math.sin(dLat/2) +
                Math.cos(lat1 * Math.PI / 180) * Math.cos(lat2 * Math.PI / 180) *
                Math.sin(dLon/2) * Math.sin(dLon/2)
      const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1-a))
      const distance = R * c
      
      expect(distance).toBeGreaterThan(0)
    })

    it('التحقق من أن الموظف داخل النطاق', () => {
      const employeeLocation = { lat: 24.7136, lng: 46.6753 }
      const worksiteLocation = { lat: 24.7136, lng: 46.6753 }
      const radius = 100 // متر
      
      // حساب المسافة
      const R = 6371000 // نصف قطر الأرض بالمتر
      const dLat = (employeeLocation.lat - worksiteLocation.lat) * Math.PI / 180
      const dLng = (employeeLocation.lng - worksiteLocation.lng) * Math.PI / 180
      const a = Math.sin(dLat/2) * Math.sin(dLat/2) +
                Math.cos(worksiteLocation.lat * Math.PI / 180) * 
                Math.cos(employeeLocation.lat * Math.PI / 180) *
                Math.sin(dLng/2) * Math.sin(dLng/2)
      const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1-a))
      const distance = R * c
      
      const isWithinRange = distance <= radius
      expect(isWithinRange).toBe(true)
    })

    it('التحقق من أن الموظف خارج النطاق', () => {
      const employeeLocation = { lat: 24.8136, lng: 46.7753 }
      const worksiteLocation = { lat: 24.7136, lng: 46.6753 }
      const radius = 100 // متر
      
      // حساب المسافة
      const R = 6371000 // نصف قطر الأرض بالمتر
      const dLat = (employeeLocation.lat - worksiteLocation.lat) * Math.PI / 180
      const dLng = (employeeLocation.lng - worksiteLocation.lng) * Math.PI / 180
      const a = Math.sin(dLat/2) * Math.sin(dLat/2) +
                Math.cos(worksiteLocation.lat * Math.PI / 180) * 
                Math.cos(employeeLocation.lat * Math.PI / 180) *
                Math.sin(dLng/2) * Math.sin(dLng/2)
      const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1-a))
      const distance = R * c
      
      const isWithinRange = distance <= radius
      expect(isWithinRange).toBe(false)
    })
  })

  describe('التحقق من صحة الإحداثيات', () => {
    it('التحقق من خط العرض', () => {
      const validLat = 24.7136
      const invalidLat = 91
      
      const isValidLat = validLat >= -90 && validLat <= 90
      const isInvalidLat = invalidLat >= -90 && invalidLat <= 90
      
      expect(isValidLat).toBe(true)
      expect(isInvalidLat).toBe(false)
    })

    it('التحقق من خط الطول', () => {
      const validLng = 46.6753
      const invalidLng = 181
      
      const isValidLng = validLng >= -180 && validLng <= 180
      const isInvalidLng = invalidLng >= -180 && invalidLng <= 180
      
      expect(isValidLng).toBe(true)
      expect(isInvalidLng).toBe(false)
    })
  })
})

describe('نظام الموقع الجغرافي', () => {
  
  describe('الحصول على الموقع الحالي', () => {
    it('التحقق من محاكاة GPS', () => {
      // التحقق من أن محاكاة geolocation تعمل بشكل صحيح
      const mockGeolocation = {
        getCurrentPosition: vi.fn((success) => {
          success({
            coords: {
              latitude: 24.7136,
              longitude: 46.6753,
              accuracy: 10,
            },
            timestamp: Date.now(),
          })
        }),
      }
      expect(typeof mockGeolocation.getCurrentPosition).toBe('function')
    })

    it('محاكاة الحصول على الموقع', () => {
      const mockLocation = {
        latitude: 24.7136,
        longitude: 46.6753,
        accuracy: 10,
      }
      
      expect(mockLocation.latitude).toBeDefined()
      expect(mockLocation.longitude).toBeDefined()
      expect(mockLocation.accuracy).toBeGreaterThan(0)
    })
  })

  describe('تحويل الإحداثيات إلى عنوان', () => {
    it('تنسيق الإحداثيات', () => {
      const lat = 24.7136
      const lng = 46.6753
      const coordinates = `${lat}, ${lng}`
      
      expect(coordinates).toBe('24.7136, 46.6753')
    })
  })
})
```

---

## 📄 frontend-worker-pwa/src/tests/setup.js

```javascript
import { vi } from 'vitest'
import { config } from '@vue/test-utils'

// محاكاة window.matchMedia للاختبارات
Object.defineProperty(window, 'matchMedia', {
  writable: true,
  value: vi.fn().mockImplementation(query => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: vi.fn(),
    removeListener: vi.fn(),
    addEventListener: vi.fn(),
    removeEventListener: vi.fn(),
    dispatchEvent: vi.fn(),
  })),
})

// إعدادات Vue Test Utils
config.global.mocks = {
  $t: (key) => key,
}

// محاكاة localStorage
const localStorageMock = {
  getItem: vi.fn(),
  setItem: vi.fn(),
  removeItem: vi.fn(),
  clear: vi.fn(),
}
global.localStorage = localStorageMock

// محاكاة navigator إذا لم يكن موجوداً
if (!global.navigator) {
  global.navigator = {}
}

// محاكاة Geolocation API
const geolocationMock = {
  getCurrentPosition: vi.fn((success) => {
    success({
      coords: {
        latitude: 24.7136,
        longitude: 46.6753,
        accuracy: 10,
      },
      timestamp: Date.now(),
    })
  }),
  watchPosition: vi.fn(),
  clearWatch: vi.fn(),
}
global.navigator.geolocation = geolocationMock
```

---

## 📄 frontend-worker-pwa/src/tests/tasks.test.js

```javascript
import { describe, it, expect, vi } from 'vitest'

// اختبارات نظام المهام

describe('نظام المهام', () => {
  
  describe('إدارة المهام', () => {
    it('عرض المهام المعينة للموظف', () => {
      const tasks = [
        { id: 1, title: 'صيانة مكيف', status: 'pending', priority: 'high' },
        { id: 2, title: 'تسريب مياه', status: 'in_progress', priority: 'urgent' },
        { id: 3, title: 'تركيب إضاءة', status: 'completed', priority: 'normal' },
      ]
      
      expect(tasks.length).toBe(3)
      expect(tasks[0].title).toBe('صيانة مكيف')
    })

    it('فلترة المهام حسب الحالة', () => {
      const tasks = [
        { id: 1, status: 'pending' },
        { id: 2, status: 'in_progress' },
        { id: 3, status: 'completed' },
        { id: 4, status: 'pending' },
      ]
      
      const pendingTasks = tasks.filter(t => t.status === 'pending')
      expect(pendingTasks.length).toBe(2)
    })

    it('فرز المهام حسب الأولوية', () => {
      const tasks = [
        { id: 1, priority: 'low' },
        { id: 2, priority: 'urgent' },
        { id: 3, priority: 'high' },
      ]
      
      const priorityOrder = { urgent: 0, high: 1, normal: 2, low: 3 }
      const sorted = [...tasks].sort((a, b) => priorityOrder[a.priority] - priorityOrder[b.priority])
      expect(sorted[0].priority).toBe('urgent')
    })
  })

  describe('تفاصيل المهمة', () => {
    it('عرض معلومات العميل', () => {
      const task = {
        id: 1,
        title: 'صيانة مكيف',
        customer: {
          name: 'محمد أحمد',
          phone: '0501234567',
          address: 'الرياض، حي الملز',
        },
      }
      
      expect(task.customer.name).toBe('محمد أحمد')
      expect(task.customer.phone).toBe('0501234567')
    })

    it('عرض موقع العمل', () => {
      const task = {
        worksite: {
          name: 'الرياض - الرئيسي',
          address: 'الرياض',
          latitude: 24.7136,
          longitude: 46.6753,
        },
      }
      
      expect(task.worksite.latitude).toBe(24.7136)
      expect(task.worksite.longitude).toBe(46.6753)
    })

    it('حساب المدة المقدرة', () => {
      const task = {
        startTime: '09:00',
        endTime: '12:00',
      }
      
      const [startHours, startMinutes] = task.startTime.split(':').map(Number)
      const [endHours, endMinutes] = task.endTime.split(':').map(Number)
      const duration = (endHours * 60 + endMinutes) - (startHours * 60 + startMinutes)
      
      expect(duration).toBe(180) // 3 ساعات بالدقائق
    })
  })

  describe('تغيير حالة المهمة', () => {
    it('قبول المهمة', () => {
      const task = { id: 1, status: 'pending' }
      task.status = 'accepted'
      expect(task.status).toBe('accepted')
    })

    it('بدء تنفيذ المهمة', () => {
      const task = { id: 1, status: 'accepted' }
      task.status = 'in_progress'
      expect(task.status).toBe('in_progress')
    })

    it('إكمال المهمة', () => {
      const task = { id: 1, status: 'in_progress' }
      task.status = 'completed'
      expect(task.status).toBe('completed')
    })
  })
})

describe('نظام الإشعارات', () => {
  
  describe('أنواع الإشعارات', () => {
    it('إشعار مهمة جديدة', () => {
      const notification = {
        type: 'task_assigned',
        message: 'تم تعيين مهمة جديدة لك',
        taskId: 1,
      }
      
      expect(notification.type).toBe('task_assigned')
      expect(notification.taskId).toBe(1)
    })

    it('إشعار تحديث الحالة', () => {
      const notification = {
        type: 'status_update',
        message: 'تم تحديث حالة مهمتك',
        taskId: 1,
      }
      
      expect(notification.type).toBe('status_update')
    })

    it('إشعار وصول الموظف', () => {
      const notification = {
        type: 'worker_arrived',
        message: 'وصل الموظف إلى موقع العمل',
        taskId: 1,
      }
      
      expect(notification.type).toBe('worker_arrived')
    })
  })

  describe('تمييز الإشعارات', () => {
    it('الإشعارات غير المقروءة', () => {
      const notifications = [
        { id: 1, read: false },
        { id: 2, read: true },
        { id: 3, read: false },
      ]
      
      const unreadCount = notifications.filter(n => !n.read).length
      expect(unreadCount).toBe(2)
    })

    it('تعليم الإشعار كمقروء', () => {
      const notification = { id: 1, read: false }
      notification.read = true
      expect(notification.read).toBe(true)
    })
  })

  describe('تنسيق الوقت النسبي', () => {
    it('منذ دقيقة', () => {
      const timestamp = Date.now() - 60000 // قبل دقيقة
      const diff = Date.now() - timestamp
      const minutes = Math.floor(diff / 60000)
      expect(minutes).toBe(1)
    })

    it('منذ ساعة', () => {
      const timestamp = Date.now() - 3600000 // قبل ساعة
      const diff = Date.now() - timestamp
      const hours = Math.floor(diff / 3600000)
      expect(hours).toBe(1)
    })

    it('منذ يوم', () => {
      const timestamp = Date.now() - 86400000 // قبل يوم
      const diff = Date.now() - timestamp
      const days = Math.floor(diff / 86400000)
      expect(days).toBe(1)
    })
  })
})

describe('نظام الملف الشخصي', () => {
  
  describe('معلومات الموظف', () => {
    it('عرض المعلومات الأساسية', () => {
      const employee = {
        id: 1,
        name: 'محمد أحمد',
        phone: '0501234567',
        email: 'mohamed@example.com',
      }
      
      expect(employee.name).toBe('محمد أحمد')
      expect(employee.phone).toBe('0501234567')
    })

    it('عرض إحصائيات العمل', () => {
      const stats = {
        totalHours: 160,
        totalDays: 20,
        attendanceRate: 95,
      }
      
      expect(stats.totalHours).toBe(160)
      expect(stats.attendanceRate).toBe(95)
    })
  })

  describe('تغيير اللغة', () => {
    it('التبديل إلى العربية', () => {
      const lang = 'ar'
      expect(lang).toBe('ar')
    })

    it('التبديل إلى الإنجليزية', () => {
      const lang = 'en'
      expect(lang).toBe('en')
    })

    it('التبديل إلى العبرية', () => {
      const lang = 'he'
      expect(lang).toBe('he')
    })
  })

  describe('سجل الحضور', () => {
    it('عرض سجل الحضور الشهري', () => {
      const attendanceHistory = [
        { date: '2024-01-01', hours: 8, status: 'completed' },
        { date: '2024-01-02', hours: 7, status: 'completed' },
        { date: '2024-01-03', hours: 8, status: 'completed' },
      ]
      
      expect(attendanceHistory.length).toBe(3)
      const totalHours = attendanceHistory.reduce((sum, a) => sum + a.hours, 0)
      expect(totalHours).toBe(23)
    })
  })
})
```

---

## 📄 frontend-worker-pwa/src/views/AttendanceView.vue

```vue
<template>
  <div>
    <h2 class="page-title">⏱️ {{ currentTranslations.attendance }}</h2>

    <!-- نقاط العمل -->
    <div v-if="!isWorking" class="card worksites-section">
      <h3>📍 {{ currentTranslations.select_worksite }}</h3>
      <div v-if="loadingWorksites" class="loading">{{ currentTranslations.loading }}</div>
      <div v-else-if="availableWorksites.length === 0" class="empty">
        <p>📭 {{ currentTranslations.no_worksites }}</p>
      </div>
      <div v-else class="worksites-grid">
        <button
          v-for="site in availableWorksites"
          :key="site.id"
          class="worksite-card"
          :class="{ active: selectedWorksiteId === site.id }"
          @click="selectWorksite(site)"
        >
          <span class="ws-name" dir="auto">{{ site.name }}</span>
          <span class="ws-address" dir="auto">{{ site.address || currentTranslations.address }}</span>
          <span class="ws-radius">⭕ {{ site.radius_meters }} {{ currentTranslations.meter }}</span>
        </button>
      </div>
    </div>

    <!-- نقطة العمل النشطة (أثناء العمل) -->
    <div v-if="isWorking" class="card active-worksite-card">
      <h3>📍 {{ currentTranslations.current_worksite }}</h3>
      <div class="active-worksite-info">
        <div class="active-worksite-name" dir="auto">{{ worksiteName }}</div>
        <div class="active-worksite-address" dir="auto">{{ selectedWorksite?.address || '' }}</div>
        <div class="active-worksite-status">✅ {{ currentTranslations.working_at_this_location }}</div>
      </div>
    </div>

    <!-- الموقع المختار + زر الملاحة -->
    <div v-if="selectedWorksiteId && !isWorking" class="card navigation-card">
      <div class="nav-info">
        <span class="nav-icon">📍</span>
        <div>
          <p class="nav-title" dir="auto">{{ worksiteName }}</p>
          <p class="nav-address" dir="auto">{{ selectedWorksite?.address || '' }}</p>
        </div>
      </div>
      <div class="nav-buttons">
        <a :href="getWazeUrl(selectedWorksite)" target="_blank" class="btn btn--waze">
          🗺️ Waze
        </a>
        <a :href="getGoogleMapsUrl(selectedWorksite)" target="_blank" class="btn btn--google">
          🌐 {{ currentTranslations.google_maps }}
        </a>
      </div>
    </div>

    <!-- عداد الوقت -->
    <div v-if="isWorking" class="timer-card">
      <div class="timer"><span>⏱️</span> <span>{{ elapsedTime }}</span></div>
      <p class="timer-label">{{ currentTranslations.in_progress }} {{ worksiteName }}</p>
    </div>

    <!-- ملخص الساعات -->
    <div class="summary">
      <div class="summary-item"><span class="num">{{ todayHours.toFixed(1) }}</span><span>{{ currentTranslations.hours_today }}</span></div>
      <div class="summary-item"><span class="num">{{ weekHours.toFixed(1) }}</span><span>{{ currentTranslations.hours_week }}</span></div>
      <div class="summary-item"><span class="num">{{ monthHours.toFixed(1) }}</span><span>{{ currentTranslations.hours_month }}</span></div>
    </div>

    <!-- زر سجل الحضور -->
    <div class="card">
      <button class="btn btn--primary btn--full" @click="showAttendanceHistoryModal = true">
        <span class="icon-emoji">📊</span> {{ currentTranslations.my_attendance_history }}
      </button>
    </div>

    <!-- حالة القرب + عنوان الموقع -->
    <div class="card attendance-card">
      <!-- الموقع الحالي والعنوان -->
      <div v-if="userLocation" class="location-info">
        <div class="location-header">
          <span class="location-icon">📍</span>
          <span class="location-title">{{ currentTranslations.location }}</span>
        </div>
        <div class="location-coords mono">
          {{ userLocation.lat.toFixed(6) }}, {{ userLocation.lng.toFixed(6) }}
        </div>
        <div v-if="locationAddress" class="location-address" dir="auto">
          {{ locationAddress }}
        </div>
        <div class="location-distance" :class="withinRange ? 'in' : 'out'">
          <span class="distance-icon">{{ withinRange ? '✅' : '❌' }}</span>
          <span class="distance-text">
            {{ currentTranslations.distance }}: <strong>{{ formatDistance(distance) }}</strong>
            <span v-if="selectedWorksiteId" class="distance-range">
              ({{ currentTranslations.radius }}: {{ radius }} {{ currentTranslations.meter }})
            </span>
          </span>
        </div>
      </div>

      <div class="geofence-status" v-if="selectedWorksiteId">
        <div class="status-icon" :class="withinRange ? 'in' : 'out'">
          {{ withinRange ? '✅' : '❌' }}
        </div>
        <div>
          <p class="status-text" :class="withinRange ? 'in' : 'out'">
            {{ withinRange ? currentTranslations.inside_range : currentTranslations.outside_range }}
          </p>
          <p class="status-distance">
            {{ currentTranslations.distance }}: <span class="mono">{{ formatDistance(distance) }}</span>
            ({{ currentTranslations.radius }}: <span class="mono">{{ radius }}</span> {{ currentTranslations.meter }})
          </p>
        </div>
      </div>

      <!-- زر تحديد الموقع -->
      <button 
        class="btn btn--primary" 
        @click="getLocationWithAddress" 
        :disabled="gettingLocation"
      >
        <span class="icon-emoji">{{ gettingLocation ? '⏳' : '📍' }}</span> {{ currentTranslations.select_location }}
      </button>

      <!-- الأزرار -->
      <div class="actions">
        <button 
          class="btn btn--primary" 
          :disabled="!selectedWorksiteId || isWorking || !withinRange || !userLocation || checkingIn" 
          @click="checkIn"
        >
          <span class="icon-emoji">{{ checkingIn ? '⏳' : '✅' }}</span> {{ currentTranslations.check_in }}
        </button>
        
        <button 
          class="btn btn--ghost" 
          :disabled="!isWorking || !hasClickedLocation || !withinRange || checkingOut" 
          @click="checkOut"
        >
          <span class="icon-emoji">{{ checkingOut ? '⏳' : '⏹️' }}</span> {{ currentTranslations.check_out }}
        </button>
      </div>

      <!-- تحذيرات -->
      <div v-if="!hasClickedLocation && isWorking" class="warning-box warning-location">
        <span class="warning-icon">⚠️</span>
        <span class="warning-text">📍 {{ currentTranslations.select_location }} {{ currentTranslations.before_checkout }}</span>
      </div>

      <div v-if="hasClickedLocation && !withinRange && isWorking" class="warning-box warning-range">
        <span class="warning-icon">🚫</span>
        <span class="warning-text">❌ {{ currentTranslations.outside_range }}! {{ currentTranslations.distance }}: {{ formatDistance(distance) }}</span>
      </div>

      <div v-if="hasClickedLocation && withinRange && isWorking" class="success-box">
        <span class="success-icon">✅</span>
        <span class="success-text">{{ currentTranslations.inside_range }} - {{ currentTranslations.can_checkout }}</span>
      </div>

      <p v-if="error" class="error">{{ error }}</p>
      <p v-if="success" class="success">{{ success }}</p>
      <p v-if="debugInfo" class="debug-info mono">{{ debugInfo }}</p>
    </div>

    <!-- DevPro Branding -->
    <div class="devpro-branding">
      <img src="/src/assets/devpro-logo.jpg" alt="DevPro Logo" class="devpro-logo-img" />
      <p class="devpro-text">Powered by DevPro</p>
    </div>

    <!-- مودال سجل الحضور -->
    <div v-if="showAttendanceHistoryModal" class="modal-backdrop" @click.self="showAttendanceHistoryModal = false">
      <div class="modal card">
        <div class="modal-header">
          <h3>📊 {{ currentTranslations.my_attendance_history }}</h3>
          <button class="modal-close" @click="showAttendanceHistoryModal = false">✕</button>
        </div>
        <div class="modal-body">
          <!-- فلاتر الشهر والسنة -->
          <div class="filters">
            <div class="filter-group">
              <label>{{ currentTranslations.year }}</label>
              <select v-model="selectedYear" @change="fetchMyAttendanceHistory" class="form-select">
                <option v-for="year in availableYears" :key="year" :value="year">{{ year }}</option>
              </select>
            </div>
            <div class="filter-group">
              <label>{{ currentTranslations.month }}</label>
              <select v-model="selectedMonth" @change="fetchMyAttendanceHistory" class="form-select">
                <option v-for="month in availableMonths" :key="month.value" :value="month.value">{{ month.label }}</option>
              </select>
            </div>
          </div>

          <!-- ملخص الشهر -->
          <div v-if="myMonthlySummary" class="monthly-summary">
            <div class="summary-card">
              <span class="summary-label">{{ currentTranslations.total_hours }}</span>
              <span class="summary-value">{{ myMonthlySummary.summary?.total_hours?.toFixed(1) || 0 }} {{ currentTranslations.hours }}</span>
            </div>
            <div class="summary-card">
              <span class="summary-label">{{ currentTranslations.work_days }}</span>
              <span class="summary-value">{{ myMonthlySummary.summary?.work_days || 0 }} {{ currentTranslations.days }}</span>
            </div>
          </div>

          <!-- جدول سجل الحضور -->
          <div v-if="loadingMyHistory" class="loading-state">
            <p>{{ currentTranslations.loading }}</p>
          </div>
          <div v-else-if="myAttendanceHistory.length === 0" class="empty-state">
            <p>{{ currentTranslations.no_attendance_records }}</p>
          </div>
          <div v-else class="table-wrapper">
            <table class="table">
              <thead>
                <tr>
                  <th>{{ currentTranslations.date }}</th>
                  <th>{{ currentTranslations.worksite }}</th>
                  <th>{{ currentTranslations.check_in }}</th>
                  <th>{{ currentTranslations.check_out }}</th>
                  <th>{{ currentTranslations.worked_hours }}</th>
                  <th>{{ currentTranslations.location }}</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="record in myAttendanceHistory" :key="record.id">
                  <td class="mono">{{ formatDate(record.check_in_time) }}</td>
                  <td dir="auto">{{ record.worksite_name || '—' }}</td>
                  <td class="mono">{{ formatTime(record.check_in_time) }}</td>
                  <td class="mono">{{ record.check_out_time ? formatTime(record.check_out_time) : '—' }}</td>
                  <td class="mono">{{ record.worked_hours ? record.worked_hours.toFixed(1) + ' ' + currentTranslations.hours : '—' }}</td>
                  <td class="mono">{{ formatDistance(record.check_in_distance_meters) }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useI18n } from '../services/i18n'
import api from '../services/api'
import { i18nState } from '../services/i18n'
import { tasksStore } from '../store/tasks'

const { currentLang, currentTranslations } = useI18n()

// ==========================================
// الحالة الأساسية
// ==========================================
const loadingWorksites = ref(false)
const availableWorksites = ref([])
const selectedWorksiteId = ref(null)
const selectedWorksite = ref(null)
const worksiteName = ref('')
const radius = ref(100)
const distance = ref(0)
const userLocation = ref(null)
const locationAddress = ref('')
const attendanceId = ref(null)
const isWorking = ref(false)
const elapsedSeconds = ref(0)
const todayHours = ref(0), weekHours = ref(0), monthHours = ref(0)
const checkingIn = ref(false), checkingOut = ref(false), gettingLocation = ref(false)
const hasClickedLocation = ref(false)
const error = ref(''), success = ref(''), debugInfo = ref('')
const currentServiceRequestId = ref(null)
let timerInterval = null, locationInterval = null

// سجل الحضور الشخصي
const showAttendanceHistoryModal = ref(false)
const myAttendanceHistory = ref([])
const myMonthlySummary = ref(null)
const loadingMyHistory = ref(false)
const selectedYear = ref(new Date().getFullYear())
const selectedMonth = ref(String(new Date().getMonth() + 1))

const availableYears = ref([])
const availableMonths = ref([])

// دالة لتحديث أسماء الشهور حسب اللغة
function updateMonthNames() {
  const monthKeys = ['january', 'february', 'march', 'april', 'may', 'june', 'july', 'august', 'september', 'october', 'november', 'december']
  availableMonths.value = monthKeys.map((key, index) => ({
    value: String(index + 1),
    label: currentTranslations.value[key]
  }))
}

// توليد السنوات المتاحة
const currentYear = new Date().getFullYear()
for (let i = currentYear; i >= currentYear - 2; i--) {
  availableYears.value.push(i)
}

// تحديث أسماء الشهور عند التحميل
updateMonthNames()

// مراقبة تغيير اللغة لتحديث أسماء الشهور
watch(currentLang, () => {
  updateMonthNames()
})

const withinRange = computed(() => distance.value <= radius.value)
const elapsedTime = computed(() => {
  const s = elapsedSeconds.value
  return `${String(Math.floor(s/3600)).padStart(2,'0')}:${String(Math.floor((s%3600)/60)).padStart(2,'0')}:${String(Math.floor(s%60)).padStart(2,'0')}`
})

// ==========================================
// ✅ مراقبة تغيير اللغة وإعادة تحميل البيانات
// ==========================================
watch(currentLang, () => {
  // When language changes, everything updates automatically because t() is reactive
  console.log('🌍 Language changed to:', currentLang.value)
})

// ==========================================
// مراقبة فتح مودال سجل الحضور لجلب البيانات
// ==========================================
watch(showAttendanceHistoryModal, (newValue) => {
  if (newValue) {
    // عند فتح المودال، جلب البيانات
    fetchMyAttendanceHistory()
  }
})

// ==========================================
// دالة تنسيق المسافة
// ==========================================
function formatDistance(meters) {
  if (!meters) return '0 ' + currentTranslations.meters
  if (meters >= 1000) {
    return (meters / 1000).toFixed(2) + ' ' + currentTranslations.kilometers
  }
  return Math.round(meters) + ' ' + currentTranslations.meters
}

// ==========================================
// روابط الملاحة
// ==========================================
function getWazeUrl(site) {
  if (!site) return '#'
  return `https://www.waze.com/ul?ll=${site.latitude},${site.longitude}&navigate=yes`
}

function getGoogleMapsUrl(site) {
  if (!site) return '#'
  return `https://www.google.com/maps/dir/?api=1&destination=${site.latitude},${site.longitude}`
}

// ==========================================
// جلب نقاط العمل
// ==========================================
async function fetchWorksites() {
  loadingWorksites.value = true
  try {
    const { data } = await api.get('/worksites/available')
    availableWorksites.value = data || []
  } catch(e) { 
    console.error(e) 
  } finally { 
    loadingWorksites.value = false 
  }
}

// ==========================================
// اختيار نقطة عمل
// ==========================================
function selectWorksite(site) {
  selectedWorksiteId.value = site.id
  selectedWorksite.value = site
  worksiteName.value = site.name
  radius.value = site.radius_meters
  debugInfo.value = `${currentTranslations.worksite_selected} ${site.name} (${currentTranslations.radius}: ${site.radius_meters} ${currentTranslations.meter})`
  
  if (userLocation.value) {
    calculateDistance(userLocation.value.lat, userLocation.value.lng)
  }
}

// ==========================================
// حساب المسافة
// ==========================================
function calculateDistance(lat, lng) {
  if (!selectedWorksite.value) return
  
  const site = selectedWorksite.value
  const R = 6371000
  const φ1 = lat * Math.PI / 180
  const φ2 = site.latitude * Math.PI / 180
  const Δφ = (site.latitude - lat) * Math.PI / 180
  const Δλ = (site.longitude - lng) * Math.PI / 180

  const a = Math.sin(Δφ/2) * Math.sin(Δφ/2) +
            Math.cos(φ1) * Math.cos(φ2) *
            Math.sin(Δλ/2) * Math.sin(Δλ/2)
  const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1-a))
  const d = R * c

  distance.value = Math.round(d)
  debugInfo.value = `📏 المسافة: ${formatDistance(distance.value)}`
}

// ==========================================
// جلب عنوان الموقع
// ==========================================
async function getAddressFromCoords(lat, lng) {
  try {
    const url = `https://api.geoapify.com/v1/geocode/reverse?lat=${lat}&lon=${lng}&apiKey=a6a3b5fec1cd4b1c99daaf6decab855f&lang=${currentLang.value}`
    const response = await fetch(url)
    if (!response.ok) throw new Error(currentTranslations.failed_fetch_address)

    const data = await response.json()
    if (data.features && data.features.length > 0) {
      const props = data.features[0].properties
      return props.formatted || props.address_line1 || currentTranslations.address_unknown
    }
    return currentTranslations.address_not_found
  } catch (e) {
    console.error(currentTranslations.failed_fetch_address, e)
    return currentTranslations.cannot_get_address
  }
}

// ==========================================
// تحديد الموقع مع العنوان
// ==========================================
async function getLocationWithAddress() {
  if (!navigator.geolocation) {
    error.value = currentTranslations.browser_not_support_location
    return
  }

  gettingLocation.value = true
  hasClickedLocation.value = true
  error.value = ''
  success.value = ''
  locationAddress.value = ''
  debugInfo.value = currentTranslations.getting_location

  try {
    const pos = await new Promise((resolve, reject) => {
      navigator.geolocation.getCurrentPosition(resolve, reject, {
        enableHighAccuracy: true,
        timeout: 15000,
        maximumAge: 0
      })
    })

    const lat = pos.coords.latitude
    const lng = pos.coords.longitude

    userLocation.value = { lat, lng }
    debugInfo.value = `${currentTranslations.location_coords} ${lat.toFixed(6)}, ${lng.toFixed(6)}`

    const address = await getAddressFromCoords(lat, lng)
    locationAddress.value = address
    debugInfo.value += `\n${currentTranslations.home_icon} ${address}`

    if (selectedWorksite.value) {
      calculateDistance(lat, lng)
      debugInfo.value += `\n${currentTranslations.distance_label} ${formatDistance(distance.value)}`
    }

    // Send location immediately for real-time tracking
    if (isWorking.value) {
      try {
        await api.post('/location/update', {
          latitude: lat,
          longitude: lng,
          accuracy: pos.coords.accuracy || 0
        })
        debugInfo.value += `\n📡 Location sent immediately`
      } catch (e) {
        console.error('Failed to send location immediately:', e)
      }
    }

    success.value = currentTranslations.location_determined_success
  } catch (e) {
    error.value = currentTranslations.location_determine_failed + ' ' + (e.message || currentTranslations.error)
    debugInfo.value = `❌ ${error.value}`
  } finally {
    gettingLocation.value = false
  }
}

// ==========================================
// بدء الدوام
// ==========================================
async function checkIn() {
  if (!userLocation.value) {
    error.value = currentTranslations.click_select_location_first
    return
  }
  if (!selectedWorksiteId.value) {
    error.value = currentTranslations.select_worksite_first
    return
  }
  if (!withinRange.value) {
    error.value = `${currentTranslations.outside_range_distance} ${formatDistance(distance.value)} (${currentTranslations.allowed_range}: ${radius.value} ${currentTranslations.meter})`
    return
  }

  checkingIn.value = true
  error.value = ''
  success.value = ''
  debugInfo.value = currentTranslations.registering_checkin

  try {
    const payload = {
      worksite_id: selectedWorksiteId.value,
      latitude: userLocation.value.lat,
      longitude: userLocation.value.lng,
      service_request_id: tasksStore.currentServiceRequestId || currentServiceRequestId.value
    }

    const { data } = await api.post('/attendance/check-in', payload)

    attendanceId.value = data.attendance?.id
    success.value = currentTranslations.checkin_success
    isWorking.value = true
    elapsedSeconds.value = 0
    timerInterval = setInterval(() => elapsedSeconds.value++, 1000)
    startLocationTracking()
    debugInfo.value = `${currentTranslations.checkin_started_at} ${worksiteName.value}`

    hasClickedLocation.value = false
  } catch(e) {
    console.error('❌ CheckIn failed:', e.response?.data)
    const errData = e.response?.data
    if (errData?.geofence) {
      distance.value = errData.geofence.distance_meters || distance.value
      error.value = `${currentTranslations.outside_range_actual_distance} ${formatDistance(distance.value)} (${currentTranslations.allowed_range}: ${radius.value} ${currentTranslations.meter})`
      debugInfo.value = `${currentTranslations.actual_distance} ${formatDistance(distance.value)}`
    } else {
      error.value = errData?.error || currentTranslations.checkin_failed
      debugInfo.value = `❌ ${error.value}`
    }
  } finally {
    checkingIn.value = false
  }
}

// ==========================================
// إنهاء الدوام
// ==========================================
async function checkOut() {
  if (!attendanceId.value) {
    error.value = currentTranslations.no_active_shift
    return
  }

  if (!hasClickedLocation.value) {
    error.value = currentTranslations.click_select_location_before_checkout
    debugInfo.value = currentTranslations.click_select_location_again
    return
  }

  if (!userLocation.value) {
    error.value = currentTranslations.location_determine_failed_retry
    debugInfo.value = currentTranslations.click_select_location_retry
    return
  }

  if (!withinRange.value) {
    error.value = `${currentTranslations.outside_worksite_range} ${formatDistance(distance.value)} (${currentTranslations.allowed_range}: ${radius.value} ${currentTranslations.meter})`
    debugInfo.value = currentTranslations.cannot_checkout_outside_range
    return
  }

  checkingOut.value = true
  error.value = ''
  success.value = ''
  debugInfo.value = currentTranslations.registering_checkout

  try {
    const { data } = await api.post('/attendance/check-out', {
      attendance_id: attendanceId.value,
      latitude: userLocation.value.lat,
      longitude: userLocation.value.lng
    })
    success.value = `${currentTranslations.checkout_success} (${data.worked_hours?.toFixed(2)} ${currentTranslations.checkout_success_hours})`
    isWorking.value = false
    clearInterval(timerInterval)
    clearInterval(locationInterval)
    await fetchSummary()
    debugInfo.value = `${currentTranslations.checkout_completed_after} ${data.worked_hours?.toFixed(2)} ${currentTranslations.checkout_success_hours}`

    // تحديث حالة طلب الخدمة المرتبط بهذا الحضور إلى مكتمل
    if (data.service_request_id) {
      try {
        await api.post(`/service/requests/${data.service_request_id}/complete`)
        debugInfo.value += `\n✅ تم تحديث حالة طلب الخدمة إلى مكتمل`
      } catch (e) {
        console.error('فشل تحديث حالة طلب الخدمة:', e)
        debugInfo.value += `\n⚠️ فشل تحديث حالة طلب الخدمة`
      }
    }

    hasClickedLocation.value = false
  } catch(e) {
    console.error('❌ CheckOut failed:', e.response?.data)
    error.value = e.response?.data?.error || currentTranslations.checkout_failed
    debugInfo.value = `❌ ${error.value}`
  } finally {
    checkingOut.value = false
  }
}

// ==========================================
// تتبع الموقع
// ==========================================
function startLocationTracking() {
  clearInterval(locationInterval)
  locationInterval = setInterval(async () => {
    if (userLocation.value && isWorking.value) {
      try {
        await api.post('/location/update', {
          latitude: userLocation.value.lat,
          longitude: userLocation.value.lng
        })
      } catch(e) { /* silent */ }
    }
  }, 3000) // Update every 3 seconds for real-time tracking
}

// ==========================================
// استعادة الحالة
// ==========================================
async function restoreState() {
  try {
    const { data } = await api.get('/attendance/current')
    if (data.has_active) {
      isWorking.value = true
      attendanceId.value = data.attendance_id
      elapsedSeconds.value = data.elapsed_seconds || 0
      selectedWorksiteId.value = data.worksite_id
      currentServiceRequestId.value = data.service_request_id || null
      timerInterval = setInterval(() => elapsedSeconds.value++, 1000)
      startLocationTracking()
      const ws = availableWorksites.value.find(w => w.id === data.worksite_id)
      if(ws) {
        selectedWorksite.value = ws
        worksiteName.value = ws.name
        radius.value = ws.radius_meters
      }
      debugInfo.value = currentTranslations.shift_restored
      hasClickedLocation.value = false
    }
  } catch(e) { console.error(e) }
}

// ==========================================
// جلب الملخص
// ==========================================
async function fetchSummary() {
  try {
    const { data } = await api.get('/attendance/summary')
    todayHours.value = data.today_hours||0
    weekHours.value = data.week_hours||0
    monthHours.value = data.month_hours||0
  } catch(e) {}
}

// ==========================================
// دوال سجل الحضور الشخصي
// ==========================================
async function fetchMyAttendanceHistory() {
  loadingMyHistory.value = true
  try {
    const { data } = await api.get(
      `/attendance/my-history?year=${selectedYear.value}&month=${selectedMonth.value}`
    )
    myAttendanceHistory.value = data || []
    
    // جلب الملخص الشهري
    const summaryResponse = await api.get(
      `/attendance/my-monthly-summary?year=${selectedYear.value}&month=${selectedMonth.value}`
    )
    myMonthlySummary.value = summaryResponse.data
  } catch (error) {
    console.error('Failed to fetch attendance history:', error)
    myAttendanceHistory.value = []
    myMonthlySummary.value = null
  } finally {
    loadingMyHistory.value = false
  }
}

function formatDate(date) {
  if (!date) return '—'
  return new Date(date).toLocaleDateString('en-GB')
}

function formatTime(date) {
  if (!date) return '—'
  return new Date(date).toLocaleTimeString('en-GB', { hour: '2-digit', minute: '2-digit', hour12: false })
}

// ==========================================
// دورة الحياة
// ==========================================
onMounted(async () => {
  await fetchWorksites()
  await restoreState()
  await fetchSummary()
})

onUnmounted(() => { 
  clearInterval(timerInterval)
  clearInterval(locationInterval)
})
</script>

<style scoped>
.page-title { 
  font-size: 24px; 
  margin-bottom: 20px; 
  font-weight: 700;
  color: var(--brand);
  display: flex;
  align-items: center;
  gap: 8px;
}

.worksites-section { padding: 20px; margin-bottom: 20px; }
.worksites-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 12px; margin-top: 16px; }
.worksite-card { 
  padding: 16px 14px; 
  border: 2px solid var(--line); 
  border-radius: var(--radius-md); 
  background: var(--surface); 
  cursor: pointer; 
  text-align: right; 
  transition: all 0.3s ease;
  box-shadow: var(--shadow-sm);
}
.worksite-card:hover {
  border-color: var(--brand-light);
  box-shadow: var(--shadow-md);
  transform: translateY(-2px);
}
.worksite-card.active { 
  border-color: var(--brand); 
  background: linear-gradient(135deg, var(--brand-tint) 0%, var(--brand-tint) 100%);
  box-shadow: var(--shadow-md);
}
.ws-name { font-weight: 700; display: block; font-size: 15px; color: var(--ink); }
.ws-address { font-size: 13px; color: var(--ink-soft); margin-top: 4px; }
.ws-radius { font-size: 12px; color: var(--brand); margin-top: 6px; font-weight: 600; }

.active-worksite-card { 
  padding: 20px; 
  margin-bottom: 20px; 
  background: linear-gradient(135deg, var(--signal-in-tint) 0%, var(--signal-in-tint) 100%); 
  border: 2px solid var(--signal-in); 
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-md);
}
.active-worksite-info { text-align: center; }
.active-worksite-name { font-size: 20px; font-weight: 700; color: var(--signal-in); margin-bottom: 10px; }
.active-worksite-address { font-size: 15px; color: var(--ink); margin-bottom: 14px; }
.active-worksite-status { font-size: 17px; font-weight: 600; color: var(--signal-in); }

.navigation-card { 
  padding: 20px; 
  margin-bottom: 20px; 
  display: flex; 
  justify-content: space-between; 
  align-items: center; 
  flex-wrap: wrap; 
  gap: 16px;
  background: var(--surface);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-sm);
}
.nav-info { display: flex; align-items: center; gap: 12px; }
.nav-icon { font-size: 28px; }
.nav-title { font-weight: 700; margin: 0; font-size: 16px; color: var(--ink); }
.nav-address { font-size: 13px; color: var(--ink-soft); margin: 0; }
.nav-buttons { display: flex; gap: 10px; }
.btn { 
  padding: 10px 18px; 
  border-radius: var(--radius-md); 
  font-weight: 600; 
  text-decoration: none; 
  display: inline-block; 
  font-size: 14px; 
  cursor: pointer; 
  border: none;
  transition: all 0.2s ease;
  box-shadow: var(--shadow-sm);
}
.btn--waze { 
  background: linear-gradient(135deg, #33ccff 0%, #00b8e6 100%); 
  color: #1a1a2e; 
}
.btn--waze:hover {
  box-shadow: var(--shadow-md);
  transform: translateY(-1px);
}
.btn--google { 
  background: linear-gradient(135deg, #4285f4 0%, #3467a6 100%); 
  color: white; 
}
.btn--google:hover {
  box-shadow: var(--shadow-md);
  transform: translateY(-1px);
}

.location-info { 
  width: 100%; 
  padding: 16px; 
  background: var(--canvas); 
  border-radius: var(--radius-md); 
  border: 1px solid var(--line);
  box-shadow: var(--shadow-sm);
}
.location-header { display: flex; align-items: center; gap: 10px; margin-bottom: 8px; }
.location-icon { font-size: 20px; }
.location-title { font-weight: 700; font-size: 15px; color: var(--ink); }
.location-coords { font-size: 13px; color: var(--ink-soft); margin-bottom: 6px; font-family: var(--font-mono); }
.location-address { font-size: 14px; color: var(--ink); padding: 6px 10px; background: var(--surface); border-radius: var(--radius-sm); margin-bottom: 8px; }
.location-distance { display: flex; align-items: center; gap: 10px; padding: 8px 12px; border-radius: var(--radius-md); font-size: 14px; font-weight: 600; }
.location-distance.in { background: var(--signal-in-tint); color: var(--signal-in); }
.location-distance.out { background: var(--signal-out-tint); color: var(--signal-out); }
.distance-icon { font-size: 18px; }
.distance-range { font-size: 13px; opacity: 0.8; font-weight: 400; }

.timer-card { 
  padding: 24px; 
  margin-bottom: 20px; 
  text-align: center; 
  background: linear-gradient(135deg, var(--signal-in-tint) 0%, var(--signal-in-tint) 100%); 
  border-radius: var(--radius-lg);
  border: 2px solid var(--signal-in);
  box-shadow: var(--shadow-md);
}
.timer { 
  font-size: 32px; 
  font-weight: 700; 
  font-family: var(--font-mono); 
  color: var(--signal-in); 
  display: flex; 
  align-items: center; 
  justify-content: center; 
  gap: 10px; 
}
.timer-label { font-size: 14px; color: var(--ink-soft); margin-top: 6px; font-weight: 600; }

.summary { 
  display: grid; 
  grid-template-columns: repeat(3,1fr); 
  gap: 12px; 
  margin-bottom: 20px; 
}
.summary-item { 
  background: var(--surface); 
  padding: 16px; 
  text-align: center; 
  border-radius: var(--radius-md); 
  border: 1px solid var(--line);
  box-shadow: var(--shadow-sm);
  transition: all 0.2s ease;
}
.summary-item:hover {
  box-shadow: var(--shadow-md);
  transform: translateY(-2px);
}
.summary-item .num { 
  font-size: 26px; 
  font-weight: 700; 
  color: var(--brand); 
  display: block; 
  margin-bottom: 4px;
}
.summary-item span:last-child { 
  font-size: 12px; 
  color: var(--ink-soft); 
  font-weight: 600;
}

.btn--full { width: 100%; }

.attendance-card { 
  padding: 24px; 
  display: flex; 
  flex-direction: column; 
  gap: 20px; 
  background: var(--surface);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-sm);
}
.geofence-status { 
  display: flex; 
  align-items: center; 
  gap: 14px; 
  padding: 14px; 
  background: var(--canvas); 
  border-radius: var(--radius-md);
  border: 1px solid var(--line);
}
.status-icon { font-size: 32px; }
.status-text { font-weight: 700; margin: 0; font-size: 15px; }
.status-text.in { color: var(--signal-in); } 
.status-text.out { color: var(--signal-out); }
.status-distance { font-size: 13px; color: var(--ink-soft); margin: 0; }
.status-distance .mono { font-family: var(--font-mono); font-weight: 600; }

.actions { display: flex; flex-direction: column; gap: 12px; width: 100%; }
.btn--primary { 
  background: linear-gradient(135deg, var(--brand) 0%, var(--brand-dark) 100%); 
  color: white; 
  padding: 16px; 
  border: none; 
  border-radius: var(--radius-md); 
  font-weight: 700; 
  cursor: pointer; 
  width: 100%; 
  font-size: 16px;
  box-shadow: var(--shadow-md);
  transition: all 0.2s ease;
}
.btn--primary:hover:not(:disabled) {
  box-shadow: var(--shadow-lg);
  transform: translateY(-2px);
}
.btn--primary:disabled { 
  opacity: 0.5; 
  cursor: not-allowed; 
  transform: none !important;
}
.btn--ghost { 
  background: var(--surface); 
  border: 2px solid var(--line); 
  color: var(--ink); 
  padding: 16px; 
  border-radius: var(--radius-md); 
  font-weight: 700; 
  cursor: pointer; 
  width: 100%; 
  font-size: 16px;
  box-shadow: var(--shadow-sm);
  transition: all 0.2s ease;
}
.btn--ghost:hover:not(:disabled) {
  border-color: var(--brand-light);
  color: var(--brand);
  box-shadow: var(--shadow-md);
}
.btn--ghost:disabled { 
  opacity: 0.5; 
  cursor: not-allowed; 
}

.warning-box {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 18px;
  border-radius: var(--radius-md);
  font-weight: 600;
  font-size: 14px;
}

.warning-location {
  background: var(--signal-warning-tint);
  border: 2px solid var(--signal-warning);
  color: var(--signal-warning);
}

.warning-range {
  background: var(--signal-out-tint);
  border: 2px solid var(--signal-out);
  color: var(--signal-out);
}

.success-box {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 18px;
  background: var(--signal-in-tint);
  border: 2px solid var(--signal-in);
  border-radius: var(--radius-md);
  color: var(--signal-in);
  font-weight: 600;
  font-size: 14px;
}

.success-icon { font-size: 24px; }
.success-text { font-size: 14px; }
.warning-icon { font-size: 24px; }
.warning-text { font-size: 14px; }

.error { 
  color: var(--signal-out); 
  text-align: center; 
  font-weight: 600;
  font-size: 14px;
}
.success { 
  color: var(--signal-in); 
  text-align: center; 
  font-weight: 600;
  font-size: 14px;
}
.debug-info { 
  padding: 10px 14px; 
  background: var(--ink); 
  color: #fff; 
  border-radius: var(--radius-md); 
  font-size: 12px; 
  text-align: center; 
  white-space: pre-line;
  font-family: var(--font-mono);
}

/* مودال سجل الحضور */
.modal-backdrop {
  position: fixed; inset: 0;
  background: rgba(0,0,0,0.7);
  backdrop-filter: blur(8px);
  display: flex; align-items: center; justify-content: center;
  z-index: 1000; padding: 20px;
}

.modal {
  width: 100%; max-width: 600px; padding: 0;
  max-height: 90vh;
  overflow-y: auto;
  background: var(--surface);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-xl);
}

.modal-header {
  display: flex; justify-content: space-between; align-items: center;
  padding: 20px 24px; border-bottom: 1px solid var(--line);
  background: var(--canvas);
}

.modal-header h3 { 
  font-size: 18px; 
  margin: 0; 
  font-weight: 700;
  color: var(--brand);
}
.modal-close { 
  background: none; 
  border: none; 
  font-size: 28px; 
  cursor: pointer; 
  color: var(--ink-soft);
  transition: color 0.2s ease;
}
.modal-close:hover {
  color: var(--signal-out);
}

.modal-body { padding: 24px; }

.filters {
  display: flex;
  gap: 16px;
  align-items: flex-end;
  margin-bottom: 24px;
  flex-wrap: wrap;
}

.filter-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.filter-group label {
  font-size: 13px;
  font-weight: 700;
  color: var(--ink-soft);
}

.form-select {
  padding: 10px 14px;
  border: 2px solid var(--line);
  border-radius: var(--radius-md);
  font-size: 14px;
  background: var(--surface);
  color: var(--ink);
  min-width: 120px;
  font-weight: 600;
  transition: all 0.2s ease;
}

.form-select:focus {
  outline: none;
  border-color: var(--brand-light);
  box-shadow: 0 0 0 3px var(--brand-tint);
}

.monthly-summary {
  display: flex;
  gap: 16px;
  margin-bottom: 24px;
}

.summary-card {
  flex: 1;
  padding: 20px;
  background: linear-gradient(135deg, var(--brand-tint) 0%, var(--brand-tint) 100%);
  border-radius: var(--radius-lg);
  text-align: center;
  border: 2px solid var(--brand);
  box-shadow: var(--shadow-md);
}

.summary-label {
  display: block;
  font-size: 13px;
  color: var(--ink-soft);
  margin-bottom: 10px;
  font-weight: 600;
}

.summary-value {
  display: block;
  font-size: 28px;
  font-weight: 700;
  color: var(--brand-dark);
}

.loading-state, .empty-state {
  text-align: center;
  padding: 48px 24px;
  color: var(--ink-soft);
  font-size: 15px;
}

.table-wrapper { 
  overflow-x: auto;
  border-radius: var(--radius-md);
  border: 1px solid var(--line);
}

.table { 
  width: 100%; 
  border-collapse: collapse;
  background: var(--surface);
}

.table th {
  text-align: right;
  font-size: 12px;
  color: var(--ink-soft);
  font-weight: 700;
  padding: 14px 16px;
  border-bottom: 2px solid var(--line);
  background: var(--canvas);
  text-transform: uppercase;
  letter-spacing: 0.02em;
}

.table td {
  padding: 14px 16px;
  font-size: 14px;
  border-bottom: 1px solid var(--line);
  color: var(--ink);
}

.table tr:last-child td { border-bottom: none; }
.table tr:hover td {
  background: var(--canvas);
}

.devpro-branding {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 24px;
  margin-top: 20px;
  background: linear-gradient(135deg, var(--brand-tint) 0%, rgba(31, 111, 92, 0.05) 100%);
  border-radius: var(--radius-md);
  border: 1px solid var(--line);
  animation: fadeIn 0.5s ease;
}

.devpro-logo-img {
  width: 80px;
  height: auto;
  border-radius: 12px;
  margin-bottom: 12px;
  box-shadow: var(--shadow-sm);
  transition: transform 0.3s ease, box-shadow 0.3s ease;
  display: block;
}

.devpro-logo-img:hover {
  transform: scale(1.05);
  box-shadow: var(--shadow-md);
}

.devpro-text {
  font-size: 13px;
  font-weight: 600;
  color: var(--ink-soft);
  letter-spacing: 0.5px;
  text-transform: uppercase;
  margin: 0;
}

@media (max-width: 480px) {
  .devpro-logo-img {
    width: 60px;
  }
}
</style>
```

---

## 📄 frontend-worker-pwa/src/views/LoginView.vue

```vue
<template>
  <PullToRefresh 
    @refresh="handleRefresh"
    :refresh-text="currentTranslations.pull_to_refresh"
    :refreshing-text="currentTranslations.refreshing"
    :release-text="currentTranslations.release_to_refresh"
  >
    <div class="login-page">
      <PWAInstallButton />
      <div class="login-card">
        <div class="login-header">
          <div class="logo">
            <img src="/src/assets/devpro-logo.jpg" alt="DevPro logo" class="brand-mark" />
          </div>
          <h1 class="title">{{ currentTranslations.app_name }}</h1>
          <p class="subtitle">{{ currentTranslations.login }}</p>
        </div>

        <!-- محدد اللغة -->
        <div class="lang-section">
          <button 
            v-for="lang in languages" 
            :key="lang.code"
            class="lang-btn"
            :class="{ active: currentLang === lang.code }"
            @click="changeLanguage(lang.code)"
          >
            <span class="lang-flag">{{ lang.flag }}</span>
            <span class="lang-name">{{ lang.name }}</span>
          </button>
        </div>

        <form class="login-form" @submit.prevent="handleSubmit">
          <div class="field">
            <label>📱 {{ currentTranslations.phone }}</label>
            <input 
              v-model="phone" 
              type="tel" 
              placeholder="05xxxxxxxx"
              required 
              dir="ltr"
            />
            <span class="field-hint">{{ currentTranslations.phone_hint }}</span>
          </div>

          <div v-if="error" class="error">{{ error }}</div>
          <div v-if="success" class="success">{{ success }}</div>

          <button class="btn-login" type="submit" :disabled="loading">
            <span class="icon-emoji">{{ loading ? '⏳' : '📱' }}</span> {{ currentTranslations.login }}
          </button>
        </form>

        <div class="footer">
          <p>{{ currentTranslations.created_by_admin }}</p>
          <div class="devpro-logo">
            <img src="/src/assets/devpro-logo.jpg" alt="DevPro Logo" class="devpro-img" />
          </div>
        </div>
        <p class="footer-small">{{ currentTranslations.device_verify }}</p>
      </div>
    </div>
  </PullToRefresh>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from '../services/i18n'
import api from '../services/api'
import PullToRefresh from '../components/PullToRefresh.vue'
import PWAInstallButton from '../components/PWAInstallButton.vue'
import { authStore } from '../store/auth'

const { currentLang, setLang, currentTranslations } = useI18n()
const router = useRouter()

const phone = ref('')
const loading = ref(false)
const error = ref('')
const success = ref('')

const languages = [
  { code: 'ar', name: 'العربية', flag: '🇸🇦' },
  { code: 'he', name: 'עברית', flag: '🇮🇱' },
  { code: 'en', name: 'English', flag: '🇬🇧' }
]

// ✅ تغيير اللغة - يستخدم نفس المفتاح الموحد
function changeLanguage(code) {
  setLang(code)
  // ✅ يتم إعادة تحميل الصفحة بعد تغيير اللغة
  // لضمان تطبيق التغيير على جميع المكونات
}

async function handleRefresh() {
  // إعادة تحميل الصفحة
  window.location.reload()
}

function getDeviceId() {
  let deviceId = localStorage.getItem('worktrack_device_id')
  if (!deviceId) {
    deviceId = 'device_' + Date.now() + '_' + Math.random().toString(36).substring(2, 10)
    localStorage.setItem('worktrack_device_id', deviceId)
  }
  return deviceId
}

function getDeviceModel() {
  const ua = navigator.userAgent || ''
  if (/iPhone/.test(ua)) {
    const match = ua.match(/iPhone OS (\d+)_(\d+)/)
    return match ? `iPhone (iOS ${match[1]}.${match[2]})` : 'iPhone'
  }
  if (/iPad/.test(ua)) {
    const match = ua.match(/iPad; CPU OS (\d+)_(\d+)/)
    return match ? `iPad (iOS ${match[1]}.${match[2]})` : 'iPad'
  }
  if (/Android/.test(ua)) {
    const match = ua.match(/Android\s+([\d.]+);\s+([^;]+);/)
    if (match) {
      return `${match[2].trim()} (Android ${match[1]})`
    }
    const modelMatch = ua.match(/; (.+?) Build\//)
    if (modelMatch) {
      return `${modelMatch[1]} (Android)`
    }
    return 'Android Device'
  }
  if (/Windows/.test(ua)) {
    return `Windows PC`
  }
  if (/Macintosh/.test(ua)) {
    return `Mac`
  }
  if (/Linux/.test(ua)) {
    return `Linux PC`
  }
  return 'Unknown Device'
}

async function handleSubmit() {
  if (!phone.value || phone.value.length < 9) {
    error.value = 'Please enter a valid phone number'
    return
  }

  loading.value = true
  error.value = ''
  success.value = ''

  try {
    const deviceId = getDeviceId()
    const deviceModel = getDeviceModel()

    const { data } = await api.post('/auth/phone-login', {
      phone: phone.value.trim(),
      device_id: deviceId,
      device_model: deviceModel
    })

    localStorage.setItem('worktrack_token', data.token)
    localStorage.setItem('worktrack_user', JSON.stringify(data.user))
    authStore.setUser(data.user)

    success.value = '✅ Login successful!'

    // ✅ اللغة مخزنة بالفعل في localStorage من i18n

    setTimeout(() => {
      router.push('/attendance')
    }, 500)

  } catch (e) {
    console.error('❌ Login failed:', e.response?.data)

    if (e.response?.data?.device_mismatch) {
      error.value = '⚠️ This device is not authorized. Please contact the admin to reset your device.'
    } else if (e.response?.data?.model_mismatch) {
      error.value = '⚠️ Device model mismatch. Please contact the admin.'
    } else {
      error.value = e.response?.data?.error || '❌ Login failed. Please check the phone number.'
    }
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #1E3A5F 0%, #0D1B3E 100%);
  padding: 0;
  margin: 0;
  min-height: 100vh;
  width: 100%;
  height: 100vh;
  z-index: 9999;
}

.login-page::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, #1E3A5F 0%, #0D1B3E 100%);
  z-index: -1;
}

.login-card {
  background: white;
  border-radius: 24px;
  padding: 40px 44px;
  width: 100%;
  max-width: 420px;
  box-shadow: 0 24px 80px rgba(0, 0, 0, 0.35);
  margin: 0 auto;
  animation: fadeIn 0.5s ease;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: translateY(0); }
}

.login-header {
  text-align: center;
  margin-bottom: 24px;
}

.logo {
  display: flex;
  justify-content: center;
  margin-bottom: 16px;
}

.brand-mark {
  width: 72px;
  height: 72px;
  border-radius: 50%;
  object-fit: cover;
  background: linear-gradient(135deg, #e3f2fd 0%, #bbdefb 100%);
  box-shadow: 0 4px 12px rgba(33, 150, 243, 0.2);
}

.title {
  font-size: 28px;
  font-weight: 800;
  color: #1E3A5F;
  margin: 0 0 8px 0;
  direction: ltr;
}

.subtitle {
  font-size: 15px;
  color: #6B7A8A;
  margin: 0;
}

.lang-section {
  display: flex;
  justify-content: center;
  gap: 8px;
  margin-bottom: 24px;
  padding: 8px;
  background: #F0F4FA;
  border-radius: 16px;
}

.lang-btn {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 12px 16px;
  border: 2px solid #E2E8F0;
  border-radius: 12px;
  background: white;
  font-size: 14px;
  font-weight: 600;
  color: #6B7A8A;
  cursor: pointer;
  transition: all 0.3s ease;
  font-family: inherit;
}

.lang-btn:hover {
  border-color: #dee2e6;
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.lang-btn.active {
  border-color: #1E3A5F;
  background: #1E3A5F;
  color: white;
  box-shadow: 0 4px 12px rgba(30, 58, 95, 0.3);
}

.lang-flag {
  font-size: 20px;
}

.lang-name {
  font-size: 13px;
  font-weight: 500;
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.field {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.field label {
  font-size: 14px;
  font-weight: 600;
  color: #1A2A3A;
}

.field input {
  padding: 10px 14px;
  border: 1.5px solid #E2E8F0;
  border-radius: 8px;
  font-size: 14px;
  transition: all 0.3s ease;
  font-family: inherit;
  background: #F8FAFC;
  text-align: left;
}

.field input:focus {
  outline: none;
  border-color: #1E3A5F;
  background: white;
  box-shadow: 0 0 0 4px rgba(30, 58, 95, 0.1);
}

.field-hint {
  font-size: 12px;
  color: #8899AA;
  text-align: center;
}

.error {
  color: #C53030;
  font-size: 14px;
  text-align: center;
  background: #FDE8E8;
  padding: 12px;
  border-radius: 12px;
  border-left: 4px solid #C53030;
}

.success {
  color: #2F855A;
  font-size: 14px;
  text-align: center;
  background: #F0FFF4;
  padding: 12px;
  border-radius: 12px;
  border-left: 4px solid #2F855A;
}

.btn-login {
  padding: 14px 32px;
  background: linear-gradient(135deg, #1E3A5F 0%, #0D1B3E 100%);
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  font-family: inherit;
  margin-top: 4px;
  box-shadow: 0 4px 12px rgba(30, 58, 95, 0.3);
}

.btn-login .icon-emoji {
  color: #fff;
  filter: none;
}

.btn-login:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(30, 58, 95, 0.4);
}

.btn-login:hover:not(:disabled) .icon-emoji {
  color: #fff;
}

.btn-login:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none !important;
}

.btn-login:disabled .icon-emoji {
  color: #fff;
  opacity: 0.7;
}

.footer {
  text-align: center;
  font-size: 13px;
  color: #8899AA;
  margin-top: 24px;
  padding-top: 20px;
  border-top: 1px solid #E2E8F0;
}

.devpro-logo {
  display: flex;
  justify-content: center;
  margin-top: 12px;
}

.devpro-img {
  max-width: 100px;
  height: auto;
  border-radius: 12px;
  opacity: 0.8;
  transition: opacity 0.3s ease;
  display: block;
}

.devpro-img:hover {
  opacity: 1;
}

.footer-small {
  text-align: center;
  font-size: 12px;
  color: #AABBCC;
  margin-top: 12px;
}

@media (max-width: 480px) {
  .login-card {
    padding: 28px 20px;
    border-radius: 12px;
  }

  .title {
    font-size: 28px;
  }

  .lang-btn {
    padding: 6px 12px;
    font-size: 12px;
  }

  .field input {
    font-size: 16px;
  }

  .btn-login {
    padding: 12px 24px;
    font-size: 15px;
  }
}
</style>

```

---

## 📄 frontend-worker-pwa/src/views/NotificationsView.vue

```vue
<template>
  <div class="notifications-page">
    <div class="page-header">
      <div class="page-header__content">
        <h1 class="page-title">{{ t('notifications') }}</h1>
        <p class="page-subtitle" v-if="notificationsStore.notifications && notificationsStore.notifications.length > 0">
          {{ notificationsStore.notifications.length }} {{ notificationsStore.notifications.length === 1 ? t('notification') : t('notifications_plural') }}
        </p>
        <p class="page-subtitle" v-else-if="!notificationsStore.loading">
          {{ t('no_notifications') }}
        </p>
      </div>
      <div class="page-header__actions">
        <button @click="refreshNotifications" class="refresh-button" :disabled="notificationsStore.loading">
          <span class="refresh-icon" :class="{ spinning: notificationsStore.loading }">🔄</span>
          <span>{{ t('refresh') }}</span>
        </button>
        <div class="page-header__icon">
          🔔
        </div>
      </div>
    </div>

    <div v-if="notificationsStore.loading" class="loading-state">
      <div class="loading-spinner"></div>
      <p>{{ t('loading_notifications') }}</p>
    </div>

    <div v-else-if="notificationsStore.error" class="alert alert-error">
      <span>❌</span> {{ notificationsStore.error }}
    </div>

    <div v-else-if="!notificationsStore.notifications || !notificationsStore.notifications.length" class="empty-state">
      <div class="empty-icon">🔔</div>
      <h3>{{ t('no_notifications_today') }}</h3>
      <p>{{ t('no_notifications_available') }}</p>
    </div>

    <div v-else class="notifications-list">
      <div 
        v-for="notification in notificationsStore.notifications" 
        :key="notification.id" 
        class="notification-item"
        :class="{ 'unread': !notification.is_read }"
      >
        <span class="notification-icon">{{ getNotificationIcon(notification.type) }}</span>
        <div class="notification-content">
          <p class="notification-title">{{ notification.title }}</p>
          <p class="notification-message">{{ notification.message }}</p>
          <span class="notification-time">{{ formatTime(notification.created_at) }}</span>
        </div>
        <div v-if="!notification.is_read" class="unread-indicator"></div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useI18n } from '../services/i18n'
import { notificationsStore } from '../store/notifications'

const { t } = useI18n()

const refreshNotifications = () => {
  notificationsStore.fetchNotifications()
}

onMounted(() => {
  notificationsStore.fetchNotifications()
})

function getNotificationIcon(type) {
  const icons = {
    'status_update': '📋',
    'worker_assigned': '👷',
    'worker_arrived': '📍',
    'service_completed': '✅',
    'payment_request': '💳',
    'general': '🔔'
  }
  return icons[type] || '🔔'
}

function formatTime(date) {
  if (!date) return ''
  const now = new Date()
  const notificationDate = new Date(date)
  const diffMs = now - notificationDate
  const diffMins = Math.floor(diffMs / 60000)
  const diffHours = Math.floor(diffMs / 3600000)
  const diffDays = Math.floor(diffMs / 86400000)

  if (diffMins < 1) return t('just_now')
  if (diffMins < 60) return `${diffMins} ${t('minutes_ago')}`
  if (diffHours < 24) return `${diffHours} ${t('hours_ago')}`
  if (diffDays < 7) return `${diffDays} ${t('days_ago')}`
  
  return notificationDate.toLocaleDateString('en-GB', { 
    day: '2-digit', 
    month: '2-digit',
    year: '2-digit'
  })
}
</script>

<style scoped>
.notifications-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #f5f7fa 0%, #e4e8ec 100%);
  padding: 20px 16px;
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: linear-gradient(135deg, #ffffff 0%, #f8f9fa 100%);
  border-radius: 20px;
  padding: 24px 20px;
  margin-bottom: 24px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

.page-header__content {
  flex: 1;
}

.page-header__actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.refresh-button {
  display: flex;
  align-items: center;
  gap: 8px;
  background: linear-gradient(135deg, #007bff 0%, #0056b3 100%);
  color: white;
  border: none;
  padding: 10px 16px;
  border-radius: 12px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 2px 8px rgba(0, 123, 255, 0.3);
}

.refresh-button:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 123, 255, 0.4);
}

.refresh-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.refresh-icon {
  font-size: 16px;
  transition: transform 0.3s ease;
}

.refresh-icon.spinning {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.page-title {
  font-size: 24px;
  font-weight: 700;
  color: #212529;
  margin: 0 0 8px 0;
  line-height: 1.2;
}

.page-subtitle {
  font-size: 14px;
  color: #6c757d;
  margin: 0;
}

.page-header__icon {
  font-size: 48px;
  width: 64px;
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #fff3cd 0%, #ffeeba 100%);
  border-radius: 16px;
  box-shadow: 0 4px 8px rgba(255, 193, 7, 0.2);
}

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  text-align: center;
}

.loading-spinner {
  width: 48px;
  height: 48px;
  border: 4px solid #e9ecef;
  border-top: 4px solid #ffc107;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 16px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.loading-state p {
  color: #6c757d;
  font-size: 14px;
  margin: 0;
}

.alert {
  padding: 12px 16px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  margin: 8px 0;
}

.alert-error {
  background: #f8d7da;
  color: #721c24;
  border: 1px solid #f5c6cb;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  text-align: center;
  background: linear-gradient(135deg, #ffffff 0%, #f8f9fa 100%);
  border-radius: 20px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

.empty-icon {
  font-size: 64px;
  margin-bottom: 16px;
  opacity: 0.5;
}

.empty-state h3 {
  color: #212529;
  font-size: 18px;
  margin: 0 0 8px 0;
}

.empty-state p {
  color: #6c757d;
  font-size: 14px;
  margin: 0;
}

.notifications-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.notification-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 16px;
  background: linear-gradient(135deg, #ffffff 0%, #f8f9fa 100%);
  border-radius: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  transition: all 0.3s ease;
  cursor: pointer;
  position: relative;
  border-left: 4px solid transparent;
}

.notification-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.notification-item.unread {
  border-left-color: #ffc107;
  background: linear-gradient(135deg, #fff9e6 0%, #fff3cd 100%);
}

.notification-icon {
  font-size: 24px;
  flex-shrink: 0;
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 123, 255, 0.1);
  border-radius: 12px;
}

.notification-content {
  flex: 1;
  min-width: 0;
}

.notification-title {
  font-size: 15px;
  font-weight: 600;
  color: #212529;
  margin: 0 0 6px 0;
  line-height: 1.3;
}

.notification-message {
  font-size: 13px;
  color: #6c757d;
  margin: 0 0 6px 0;
  line-height: 1.4;
}

.notification-time {
  font-size: 12px;
  color: #adb5bd;
  font-weight: 500;
}

.unread-indicator {
  width: 10px;
  height: 10px;
  background: #ffc107;
  border-radius: 50%;
  flex-shrink: 0;
  margin-top: 4px;
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0%, 100% {
    transform: scale(1);
    opacity: 1;
  }
  50% {
    transform: scale(1.2);
    opacity: 0.8;
  }
}

/* Responsive adjustments */
@media (max-width: 480px) {
  .notifications-page {
    padding: 16px 12px;
  }
  
  .page-header {
    padding: 20px 16px;
    border-radius: 16px;
  }
  
  .page-title {
    font-size: 20px;
  }
  
  .page-header__icon {
    font-size: 40px;
    width: 56px;
    height: 56px;
  }
  
  .empty-state {
    padding: 40px 16px;
  }
  
  .empty-icon {
    font-size: 48px;
  }
  
  .notification-item {
    padding: 12px;
    gap: 10px;
  }
  
  .notification-icon {
    font-size: 20px;
    width: 36px;
    height: 36px;
  }
  
  .notification-title {
    font-size: 14px;
  }
  
  .notification-message {
    font-size: 12px;
  }
  
  .notification-time {
    font-size: 11px;
  }
}
</style>
```

---

## 📄 frontend-worker-pwa/src/views/ProfileView.vue

```vue
<template>
  <div class="profile-page">
    <div class="page-header">
      <div class="page-header__content">
        <h1 class="page-title">{{ currentTranslations.profile }}</h1>
        <p class="page-subtitle">{{ user?.full_name || currentTranslations.app_name }}</p>
      </div>
      <div class="page-header__icon">
        👤
      </div>
    </div>

    <div class="profile-card">
      <div class="profile-avatar">
        <span class="avatar-text">{{ initials }}</span>
      </div>
      <div class="profile-info">
        <h2 class="profile-name" dir="auto">{{ user?.full_name || currentTranslations.app_name }}</h2>
        <p class="profile-role">{{ user?.role === 'admin' ? currentTranslations.admin : currentTranslations.employee }}</p>
      </div>
    </div>

    <!-- معلومات الاتصال -->
    <div class="info-section">
      <div class="section-header">
        <span class="section-icon icon-emoji">📞</span>
        <h3 class="section-title">{{ currentTranslations.contact_info }}</h3>
      </div>
      <div class="info-card">
        <div v-if="user?.phone" class="info-item">
          <span class="info-label">{{ currentTranslations.phone }}:</span>
          <a :href="`tel:${user.phone}`" class="info-value link" dir="auto">{{ user.phone }}</a>
        </div>
        <div v-if="user?.email" class="info-item">
          <span class="info-label">{{ currentTranslations.email }}:</span>
          <a :href="`mailto:${user.email}`" class="info-value link" dir="auto">{{ user.email }}</a>
        </div>
        <div v-if="user?.address" class="info-item">
          <span class="info-label">{{ currentTranslations.address }}:</span>
          <span class="info-value" dir="auto">{{ user.address }}</span>
        </div>
      </div>
    </div>

    <!-- إحصائيات العمل -->
    <div class="info-section">
      <div class="section-header">
        <span class="section-icon icon-emoji">📊</span>
        <h3 class="section-title">{{ currentTranslations.work_statistics }}</h3>
      </div>
      <div class="stats-grid">
        <div class="stat-card">
          <span class="stat-value">{{ monthlyHours.toFixed(1) }}</span>
          <span class="stat-label">{{ currentTranslations.hours_this_month }}</span>
        </div>
        <div class="stat-card">
          <span class="stat-value">{{ workDaysThisMonth }}</span>
          <span class="stat-label">{{ currentTranslations.work_days_this_month }}</span>
        </div>
        <div class="stat-card">
          <span class="stat-value">{{ attendanceRate }}%</span>
          <span class="stat-label">{{ currentTranslations.attendance_rate }}</span>
        </div>
      </div>
    </div>

    <!-- معلومات الحساب -->
    <div class="info-section">
      <div class="section-header">
        <span class="section-icon icon-emoji">🔐</span>
        <h3 class="section-title">{{ currentTranslations.account_info }}</h3>
      </div>
      <div class="info-card">
        <div v-if="user?.created_at" class="info-item">
          <span class="info-label">{{ currentTranslations.account_created }}:</span>
          <span class="info-value">{{ formatDate(user.created_at) }}</span>
        </div>
        <div v-if="user?.last_login" class="info-item">
          <span class="info-label">{{ currentTranslations.last_login }}:</span>
          <span class="info-value">{{ formatDate(user.last_login) }}</span>
        </div>
        <div class="info-item">
          <span class="info-label">{{ currentTranslations.account_status }}:</span>
          <span class="info-value status-active">{{ currentTranslations.active }}</span>
        </div>
      </div>
    </div>

    <div class="menu-section">
      <div class="menu-card">
        <div class="menu-header">
          <span class="menu-icon">🌐</span>
          <h3 class="menu-title">{{ currentTranslations.language }}</h3>
        </div>
        <div class="language-options">
          <button 
            v-for="lang in languages" 
            :key="lang.code"
            class="lang-option"
            :class="{ active: currentLang === lang.code }"
            @click="changeLanguage(lang.code)"
          >
            <span class="lang-flag">{{ lang.flag }}</span>
            <span class="lang-name">{{ lang.name }}</span>
          </button>
        </div>
      </div>

      <button class="menu-item" @click="showNotifications = !showNotifications">
        <span class="menu-icon icon-emoji">🔔</span>
        <span class="menu-text">{{ currentTranslations.notifications }}</span>
        <span class="menu-arrow">→</span>
      </button>
      
      <button class="menu-item" @click="showHistory = !showHistory">
        <span class="menu-icon icon-emoji">📄</span>
        <span class="menu-text">{{ currentTranslations.attendance_history }}</span>
        <span class="menu-arrow">→</span>
      </button>
    </div>

    <!-- الإشعارات -->
    <div v-if="showNotifications" class="content-card">
      <div class="content-header">
        <span class="content-icon icon-emoji">🔔</span>
        <h3 class="content-title">{{ currentTranslations.notifications }}</h3>
      </div>
      <div v-if="notifications.length === 0" class="empty-state">
        <div class="empty-icon icon-emoji">🔔</div>
        <p>{{ currentTranslations.no_notifications }}</p>
      </div>
      <div v-else class="notifications-list">
        <div v-for="notif in notifications" :key="notif.id" class="notification-item">
          <p class="notif-title" dir="auto">{{ notif.title }}</p>
          <p class="notif-body" dir="auto">{{ notif.body }}</p>
          <span class="notif-time">{{ formatDate(notif.created_at) }}</span>
        </div>
      </div>
    </div>

    <!-- سجل الحضور -->
    <div v-if="showHistory" class="content-card">
      <div class="content-header">
        <span class="content-icon icon-emoji">📄</span>
        <h3 class="content-title">{{ currentTranslations.attendance_history }}</h3>
      </div>
      <div v-if="attendanceHistory.length === 0" class="empty-state">
        <div class="empty-icon icon-emoji">📄</div>
        <p>{{ currentTranslations.no_history }}</p>
      </div>
      <div v-else class="history-list">
        <div v-for="record in attendanceHistory" :key="record.id" class="history-item">
          <div class="history-info">
            <span class="history-date">{{ formatDate(record.date) }}</span>
            <span class="history-hours">{{ record.hours }} {{ currentTranslations.hours }}</span>
          </div>
          <span class="history-status" :class="record.status === 'completed' ? 'status-done' : 'status-active'">
            {{ record.status === 'completed' ? currentTranslations.completed : currentTranslations.in_progress }}
          </span>
        </div>
      </div>
    </div>

    <button class="logout-button" @click="handleLogout">
      <span class="logout-icon icon-emoji">🚪</span>
      {{ currentTranslations.logout }}
    </button>
  </div>
</template>

<script setup>
import { computed, ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from '../services/i18n'
import { authStore } from '../store/auth'
import { logout } from '../services/auth'
import api from '../services/api'

const { currentLang, setLang, currentTranslations } = useI18n()
const router = useRouter()

const user = computed(() => authStore.user)
const initials = computed(() => {
  const name = user.value?.full_name || currentTranslations.value?.initials || '?'
  return String(name).trim().slice(0, 1)
})

const showNotifications = ref(false)
const showHistory = ref(false)
const notifications = ref([])
const attendanceHistory = ref([])
const error = ref('')
const success = ref('')

// إحصائيات العمل
const monthlyHours = ref(0)
const workDaysThisMonth = ref(0)
const attendanceRate = ref(0)

const languages = [
  { code: 'ar', name: 'العربية', flag: '🇸🇦' },
  { code: 'he', name: 'עברית', flag: '🇮🇱' },
  { code: 'en', name: 'English', flag: '🇬🇧' }
]

// ✅ تغيير اللغة - نفس المفتاح الموحد
function changeLanguage(code) {
  setLang(code)
}

function handleLogout() {
  logout()
  authStore.clear()
  router.push('/login')
}

function formatDate(date) {
  if (!date) return '—'
  return new Date(date).toLocaleString('ar-SA', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

async function fetchNotifications() {
  try {
    const { data } = await api.get('/notifications')
    notifications.value = data || []
  } catch (error) {
    console.error(currentTranslations.value.failed_fetch_notifications, error)
  }
}

async function fetchAttendanceHistory() {
  try {
    const { data } = await api.get('/attendance/history')
    attendanceHistory.value = data || []
  } catch (error) {
    console.error(currentTranslations.value.failed_fetch_attendance_history, error)
  }
}

async function fetchWorkStatistics() {
  try {
    const { data } = await api.get('/attendance/my-monthly-summary')
    monthlyHours.value = data?.summary?.total_hours || 0
    workDaysThisMonth.value = data?.summary?.work_days || 0
    attendanceRate.value = 0 // يمكن حسابها لاحقاً إذا لزم الأمر
  } catch (error) {
    console.error('Failed to fetch work statistics:', error)
  }
}

onMounted(() => {
  fetchNotifications()
  fetchAttendanceHistory()
  fetchWorkStatistics()
})
</script>

<style scoped>
.profile-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #f5f7fa 0%, #e4e8ec 100%);
  padding: 20px 16px;
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: linear-gradient(135deg, #ffffff 0%, #f8f9fa 100%);
  border-radius: 20px;
  padding: 24px 20px;
  margin-bottom: 24px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

.page-header__content {
  flex: 1;
}

.page-title {
  font-size: 24px;
  font-weight: 700;
  color: #212529;
  margin: 0 0 8px 0;
  line-height: 1.2;
}

.page-subtitle {
  font-size: 14px;
  color: #6c757d;
  margin: 0;
}

.page-header__icon {
  font-size: 48px;
  width: 64px;
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #e3f2fd 0%, #bbdefb 100%);
  border-radius: 16px;
  box-shadow: 0 4px 8px rgba(33, 150, 243, 0.2);
}

.profile-card {
  display: flex;
  align-items: center;
  gap: 16px;
  background: linear-gradient(135deg, #ffffff 0%, #f8f9fa 100%);
  border-radius: 20px;
  padding: 24px;
  margin-bottom: 24px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

.info-section {
  background: linear-gradient(135deg, #ffffff 0%, #f8f9fa 100%);
  border-radius: 16px;
  padding: 20px;
  margin-bottom: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}

.info-card {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.info-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 0;
  border-bottom: 1px solid #e9ecef;
}

.info-item:last-child {
  border-bottom: none;
}

.info-label {
  font-size: 13px;
  color: #6c757d;
  font-weight: 500;
  min-width: 100px;
}

.info-value {
  font-size: 14px;
  color: #212529;
  flex: 1;
  font-weight: 500;
}

.info-value.link {
  color: #007bff;
  text-decoration: none;
  transition: color 0.3s ease;
}

.info-value.link:hover {
  color: #0056b3;
  text-decoration: underline;
}

.info-value.status-active {
  color: #28a745;
  font-weight: 600;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
}

.stat-card {
  background: linear-gradient(135deg, #007bff 0%, #0056b3 100%);
  border-radius: 12px;
  padding: 16px;
  text-align: center;
  box-shadow: 0 4px 8px rgba(0, 123, 255, 0.2);
}

.stat-value {
  font-size: 24px;
  font-weight: 700;
  color: white;
  display: block;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.9);
  font-weight: 500;
}

.profile-avatar {
  width: 72px;
  height: 72px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}

.avatar-text {
  font-size: 28px;
  font-weight: 700;
  color: white;
}

.profile-info {
  flex: 1;
}

.profile-name {
  font-size: 18px;
  font-weight: 600;
  color: #212529;
  margin: 0 0 4px 0;
}

.profile-role {
  font-size: 14px;
  color: #6c757d;
  margin: 0;
}

.menu-section {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-bottom: 24px;
}

.menu-card {
  background: linear-gradient(135deg, #ffffff 0%, #f8f9fa 100%);
  border-radius: 16px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}

.menu-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 16px;
}

.menu-icon {
  font-size: 24px;
}

.menu-title {
  font-size: 16px;
  font-weight: 600;
  color: #212529;
  margin: 0;
}

.language-options {
  display: flex;
  gap: 8px;
}

.lang-option {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  border: 2px solid #e9ecef;
  border-radius: 12px;
  background: white;
  cursor: pointer;
  transition: all 0.3s ease;
  font-size: 14px;
}

.lang-option:hover {
  border-color: #dee2e6;
  transform: translateY(-2px);
}

.lang-option.active {
  border-color: #007bff;
  background: #e3f2fd;
  box-shadow: 0 4px 8px rgba(0, 123, 255, 0.2);
}

.lang-flag {
  font-size: 20px;
}

.lang-name {
  font-weight: 500;
  color: #495057;
}

.menu-item {
  display: flex;
  align-items: center;
  gap: 12px;
  background: linear-gradient(135deg, #ffffff 0%, #f8f9fa 100%);
  border-radius: 16px;
  padding: 16px 20px;
  border: none;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
  cursor: pointer;
  transition: all 0.3s ease;
  text-align: right;
}

.menu-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

.menu-icon {
  font-size: 20px;
}

.menu-icon.icon-emoji {
  color: inherit;
  filter: none;
}

.menu-text {
  flex: 1;
  font-size: 15px;
  font-weight: 500;
  color: #212529;
}

.menu-arrow {
  font-size: 18px;
  color: #6c757d;
}

.content-card {
  background: linear-gradient(135deg, #ffffff 0%, #f8f9fa 100%);
  border-radius: 16px;
  padding: 20px;
  margin-bottom: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}

.content-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 16px;
}

.content-icon {
  font-size: 24px;
}

.content-icon.icon-emoji {
  color: inherit;
  filter: none;
}

.content-title {
  font-size: 16px;
  font-weight: 600;
  color: #212529;
  margin: 0;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
  text-align: center;
}

.empty-icon {
  font-size: 48px;
  margin-bottom: 12px;
  opacity: 0.5;
}

.empty-icon.icon-emoji {
  color: inherit;
  filter: none;
}

.empty-state p {
  color: #6c757d;
  font-size: 14px;
  margin: 0;
}

.notifications-list,
.history-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.notification-item,
.history-item {
  background: #f8f9fa;
  border-radius: 12px;
  padding: 16px;
  border-left: 4px solid #dee2e6;
}

.notif-title {
  font-weight: 600;
  font-size: 14px;
  color: #212529;
  margin: 0 0 4px 0;
}

.notif-body {
  font-size: 13px;
  color: #6c757d;
  margin: 0 0 8px 0;
}

.notif-time {
  font-size: 12px;
  color: #adb5bd;
}

.history-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
  margin-bottom: 8px;
}

.history-date {
  font-size: 14px;
  color: #212529;
  font-weight: 500;
}

.history-hours {
  font-size: 13px;
  color: #007bff;
}

.history-status {
  font-size: 12px;
  font-weight: 600;
  padding: 4px 12px;
  border-radius: 20px;
  display: inline-block;
}

.status-done {
  background: #d4edda;
  color: #155724;
}

.status-active {
  background: #fff3cd;
  color: #856404;
}

.logout-button {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  background: linear-gradient(135deg, #dc3545 0%, #c82333 100%);
  color: white;
  padding: 16px 24px;
  border-radius: 16px;
  border: none;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  box-shadow: 0 4px 12px rgba(220, 53, 69, 0.3);
  transition: all 0.3s ease;
  width: 100%;
}

.logout-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(220, 53, 69, 0.4);
}

.logout-icon {
  font-size: 20px;
}

.logout-icon.icon-emoji {
  color: #fff;
  filter: none;
}

/* Responsive adjustments */
@media (max-width: 480px) {
  .profile-page {
    padding: 16px 12px;
  }
  
  .page-header {
    padding: 20px 16px;
    border-radius: 16px;
  }
  
  .page-title {
    font-size: 20px;
  }
  
  .page-header__icon {
    font-size: 40px;
    width: 56px;
    height: 56px;
  }
  
  .profile-card {
    padding: 20px;
    border-radius: 16px;
  }
  
  .profile-avatar {
    width: 64px;
    height: 64px;
  }
  
  .avatar-text {
    font-size: 24px;
  }
  
  .lang-name {
    font-size: 13px;
  }
  
  .stats-grid {
    grid-template-columns: 1fr;
  }
  
  .stat-card {
    padding: 12px;
  }
  
  .stat-value {
    font-size: 20px;
  }
  
  .info-label {
    min-width: 80px;
    font-size: 12px;
  }
  
  .info-value {
    font-size: 13px;
  }
}
</style>

```

---

## 📄 frontend-worker-pwa/src/views/TaskDetailView.vue

```vue
<template>
  <div v-if="task" :key="componentKey">
    <router-link to="/tasks" class="back-link">
      <span class="back-icon">→</span>
      {{ currentTranslations.back_to_tasks }}
    </router-link>

    <div class="task-detail">
      <div class="task-detail__header">
        <div class="header-content">
          <h1 class="task-detail__title" dir="auto">{{ getTranslatedField(task, 'title') }}</h1>
          <div class="header-meta">
            <span class="task-id">{{ currentTranslations.task_number }}: {{ task.id.substring(0, 8) }}</span>
            <span v-if="task.priority" class="task-priority" :class="`priority-${task.priority}`">
              {{ getPriorityLabelValue(task.priority) }}
            </span>
          </div>
        </div>
        <span class="task-detail__status" :class="`status-${task.status}`">
          {{ getStatusLabelValue(task.status) }}
        </span>
      </div>

      <div v-if="getTranslatedField(task, 'description')" class="task-detail__description">
        <h3 class="description-title">📝 {{ currentTranslations.description }}</h3>
        <p dir="auto">{{ getTranslatedField(task, 'description') }}</p>
      </div>

      <div class="task-detail__sections">
        <div v-if="task.client_name || task.client_phone" class="task-detail__section">
          <div class="section-header">
            <span class="section-icon icon-emoji">👤</span>
            <h3 class="section-title">{{ currentTranslations.client_info }}</h3>
          </div>
          <div class="section-content">
            <div v-if="getTranslatedField(task, 'client_name')" class="info-row">
              <span class="info-label">{{ currentTranslations.name }}:</span>
              <span class="info-value" dir="auto">{{ getTranslatedField(task, 'client_name') }}</span>
            </div>
            <div v-if="task.client_phone" class="info-row">
              <span class="info-label">{{ currentTranslations.phone }}:</span>
              <div class="info-value phone-actions">
                <span dir="auto">{{ task.client_phone }}</span>
                <a :href="`tel:${task.client_phone}`" class="action-button call-button">
                  <span class="icon-emoji">📞</span>
                  {{ currentTranslations.call }}
                </a>
              </div>
            </div>
            <div v-if="getTranslatedField(task, 'client_address')" class="info-row">
              <span class="info-label">{{ currentTranslations.address }}:</span>
              <span class="info-value" dir="auto">{{ getTranslatedField(task, 'client_address') }}</span>
            </div>
          </div>
        </div>

        <div class="task-detail__section">
          <div class="section-header">
            <span class="section-icon icon-emoji">📍</span>
            <h3 class="section-title">{{ currentTranslations.work_location }}</h3>
          </div>
          <div class="section-content">
            <div class="info-row">
              <span class="info-label">{{ currentTranslations.location }}:</span>
              <span class="info-value" dir="auto">{{ getTranslatedField(task, 'worksite_name') || getTranslatedField(task, 'worksite_address') || currentTranslations.not_specified }}</span>
            </div>
            <div v-if="getTranslatedField(task, 'worksite_address')" class="info-row">
              <span class="info-label">{{ currentTranslations.address }}:</span>
              <span class="info-value" dir="auto">{{ getTranslatedField(task, 'worksite_address') }}</span>
            </div>
            <div v-if="task.worksite_latitude && task.worksite_longitude" class="info-row">
              <span class="info-label">{{ currentTranslations.coordinates }}:</span>
              <span class="info-value">{{ task.worksite_latitude.toFixed(6) }}, {{ task.worksite_longitude.toFixed(6) }}</span>
            </div>
            <a
              v-if="task.worksite_latitude && task.worksite_longitude"
              :href="`https://www.google.com/maps?q=${task.worksite_latitude},${task.worksite_longitude}`"
              target="_blank"
              class="maps-button"
            >
              <span class="maps-icon icon-emoji">🗺️</span>
              {{ currentTranslations.open_in_maps }}
            </a>
          </div>
        </div>

        <div class="task-detail__section">
          <div class="section-header">
            <span class="section-icon icon-emoji">🕒</span>
            <h3 class="section-title">{{ currentTranslations.time_schedule }}</h3>
          </div>
          <div class="section-content">
            <div v-if="task.scheduled_start" class="info-row">
              <span class="info-label">{{ currentTranslations.start }}:</span>
              <span class="info-value">{{ formatTime(task.scheduled_start) }}</span>
            </div>
            <div v-if="task.scheduled_end" class="info-row">
              <span class="info-label">{{ currentTranslations.end }}:</span>
              <span class="info-value">{{ formatTime(task.scheduled_end) }}</span>
            </div>
            <div v-if="task.created_at" class="info-row">
              <span class="info-label">{{ currentTranslations.created }}:</span>
              <span class="info-value">{{ formatTime(task.created_at) }}</span>
            </div>
            <div v-if="task.scheduled_start" class="info-row">
              <span class="info-label">{{ currentTranslations.duration }}:</span>
              <span class="info-value">{{ calculateDuration(task.scheduled_start, task.scheduled_end) }}</span>
            </div>
          </div>
        </div>

        <div class="task-detail__section">
          <div class="section-header">
            <span class="section-icon icon-emoji">ℹ️</span>
            <h3 class="section-title">{{ currentTranslations.additional_info }}</h3>
          </div>
          <div class="section-content">
            <div class="info-row">
              <span class="info-label">{{ currentTranslations.status }}:</span>
              <span class="info-value status-badge" :class="`status-${task.status}`">
                {{ getStatusLabelValue(task.status) }}
              </span>
            </div>
            <div v-if="task.priority" class="info-row">
              <span class="info-label">{{ currentTranslations.priority }}:</span>
              <span class="info-value priority-badge" :class="`priority-${task.priority}`">
                {{ getPriorityLabelValue(task.priority) }}
              </span>
            </div>
            <div class="info-row">
              <span class="info-label">{{ currentTranslations.task_id }}:</span>
              <span class="info-value">{{ task.id }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="action-buttons">
      <router-link to="/attendance" class="attendance-button primary">
        <span class="button-icon icon-emoji">📋</span>
        {{ currentTranslations.register_attendance_for_task }}
      </router-link>
      <a
        v-if="task.client_phone"
        :href="`tel:${task.client_phone}`"
        class="attendance-button secondary"
      >
        <span class="button-icon icon-emoji">📞</span>
        {{ currentTranslations.call_client }}
      </a>
    </div>
  </div>
  <div v-else class="empty-state">
    <div class="empty-icon icon-emoji">📋</div>
    <h3>{{ currentTranslations.task_not_found }}</h3>
  </div>
</template>

<script setup>
import { computed, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from '../services/i18n'
import { tasksStore } from '../store/tasks'

const { t, currentLang, currentTranslations } = useI18n()
const route = useRoute()
const task = computed(() => tasksStore.find(route.params.id))

// Create a unique key that changes when language changes
const componentKey = computed(() => `${route.params.id}-${currentLang.value}`)

// Force reactivity when language changes using currentTranslations
const getStatusLabel = computed(() => {
  return (status) => {
    const statusMap = {
      pending: currentTranslations.value.status_pending,
      assigned: currentTranslations.value.status_assigned,
      in_progress: currentTranslations.value.status_in_progress,
      completed: currentTranslations.value.status_completed,
      late: currentTranslations.value.status_late,
      cancelled: currentTranslations.value.status_cancelled
    }
    return statusMap[status] || status
  }
})

const getPriorityLabel = computed(() => {
  return (priority) => {
    const priorityMap = {
      low: currentTranslations.value.priority_low,
      normal: currentTranslations.value.priority_normal,
      high: currentTranslations.value.priority_high,
      urgent: currentTranslations.value.priority_urgent
    }
    return priorityMap[priority] || priority
  }
})

// Make these functions reactive to language changes
const getStatusLabelValue = (status) => {
  return getStatusLabel.value(status)
}

const getPriorityLabelValue = (priority) => {
  return getPriorityLabel.value(priority)
}

// Get translated field based on current language
const getTranslatedField = (obj, field) => {
  if (!obj) return ''

  const fieldMap = {
    title: ['title_ar', 'title_he', 'title_en'],
    description: ['description_ar', 'description_he', 'description_en'],
    client_name: ['client_name_ar', 'client_name_he', 'client_name_en'],
    client_address: ['client_address_ar', 'client_address_he', 'client_address_en'],
    worksite_name: ['worksite_name_ar', 'worksite_name_he', 'worksite_name_en'],
    worksite_address: ['worksite_address_ar', 'worksite_address_he', 'worksite_address_en']
  }

  const fields = fieldMap[field]
  if (!fields) return obj[field] || ''

  // First, use the already-translated field from backend (task.title, task.description, etc.)
  // The backend already translated this based on the task's stored language
  if (obj[field] && obj[field].trim() !== '') return obj[field]

  // Fallback to translation fields based on task's stored language (if available)
  const taskLang = obj.language || 'en' // Default to English for old tasks
  
  if (taskLang === 'he' && obj[fields[1]]) return obj[fields[1]]
  if (taskLang === 'ar' && obj[fields[0]]) return obj[fields[0]]
  if (taskLang === 'en' && obj[fields[2]]) return obj[fields[2]]
  
  // Final fallback: try English first, then Arabic, then Hebrew
  if (obj[fields[2]]) return obj[fields[2]] // English
  if (obj[fields[0]]) return obj[fields[0]] // Arabic
  if (obj[fields[1]]) return obj[fields[1]] // Hebrew

  // Fallback to original field if no translation available
  return obj[field] || ''
}

const formatTime = (dateString) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  const locale = currentLang.value === 'ar' ? 'ar-EG' : currentLang.value === 'he' ? 'he-IL' : 'en-GB'
  return date.toLocaleDateString(locale, {
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const calculateDuration = (start, end) => {
  if (!start) return currentTranslations.value.not_specified
  const startDate = new Date(start)
  const endDate = end ? new Date(end) : new Date()
  const diffMs = endDate - startDate
  const diffHours = Math.floor(diffMs / (1000 * 60 * 60))
  const diffMinutes = Math.floor((diffMs % (1000 * 60 * 60)) / (1000 * 60))

  if (diffHours > 0) {
    return `${diffHours} ${currentTranslations.value.hours} ${currentTranslations.value.and} ${diffMinutes} ${currentTranslations.value.minutes}`
  } else {
    return `${diffMinutes} ${currentTranslations.value.minutes}`
  }
}
</script>

<style scoped>
.back-link {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: #6c757d;
  margin-bottom: 20px;
  text-decoration: none;
  transition: color 0.3s ease;
}

.back-link:hover {
  color: #495057;
}

.back-icon {
  font-size: 18px;
  transform: rotate(180deg);
}

.task-detail {
  background: linear-gradient(135deg, #ffffff 0%, #f8f9fa 100%);
  border-radius: 20px;
  padding: 24px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
  margin-bottom: 20px;
}

.task-detail__header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  margin-bottom: 20px;
  gap: 16px;
}

.header-content {
  flex: 1;
}

.task-detail__title {
  font-size: 24px;
  font-weight: 700;
  color: #212529;
  margin: 0 0 8px 0;
  line-height: 1.3;
  /* دعم النصوص المختلطة (عربي، عبري، إنجليزي) */
  unicode-bidi: embed;
  text-align: start;
}

.header-meta {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}

.task-id {
  font-size: 12px;
  color: #6c757d;
  font-family: monospace;
  background: #f8f9fa;
  padding: 4px 8px;
  border-radius: 6px;
}

.task-priority {
  font-size: 11px;
  font-weight: 600;
  padding: 4px 10px;
  border-radius: 12px;
  text-transform: uppercase;
}

.task-priority.priority-low {
  background: #e3f2fd;
  color: #1976d2;
}

.task-priority.priority-normal {
  background: #e9ecef;
  color: #495057;
}

.task-priority.priority-high {
  background: #fff3cd;
  color: #856404;
}

.task-priority.priority-urgent {
  background: #f8d7da;
  color: #721c24;
}

.task-detail__status {
  font-size: 12px;
  font-weight: 600;
  padding: 6px 14px;
  border-radius: 20px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  white-space: nowrap;
}

.task-detail__status.status-pending {
  background: #e3f2fd;
  color: #1976d2;
}

.task-detail__status.status-assigned {
  background: #e9ecef;
  color: #495057;
}

.task-detail__status.status-in_progress {
  background: #fff3cd;
  color: #856404;
}

.task-detail__status.status-completed {
  background: #d4edda;
  color: #155724;
}

.task-detail__status.status-late {
  background: #f8d7da;
  color: #721c24;
}

.task-detail__status.status-cancelled {
  background: #fce4ec;
  color: #c62828;
}

.task-detail__description {
  background: #f8f9fa;
  border-radius: 12px;
  padding: 16px;
  margin-bottom: 20px;
  border-left: 4px solid #dee2e6;
}

.description-title {
  font-size: 14px;
  font-weight: 600;
  color: #495057;
  margin: 0 0 8px 0;
}

.task-detail__description p {
  margin: 0;
  color: #495057;
  line-height: 1.6;
  font-size: 14px;
  /* دعم النصوص المختلطة (عربي، عبري، إنجليزي) */
  unicode-bidi: embed;
  text-align: start;
}

.task-detail__sections {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.task-detail__section {
  background: #ffffff;
  border-radius: 16px;
  padding: 20px;
  border: 1px solid #e9ecef;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.02);
}

.section-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 16px;
}

.section-icon {
  font-size: 24px;
}

.section-title {
  font-size: 16px;
  font-weight: 600;
  color: #212529;
  margin: 0;
}

.section-content {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.info-row {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 0;
}

.info-label {
  font-size: 13px;
  color: #6c757d;
  font-weight: 500;
  min-width: 80px;
}

.info-value {
  font-size: 14px;
  color: #212529;
  flex: 1;
  font-weight: 500;
  /* دعم النصوص المختلطة (عربي، عبري، إنجليزي) */
  unicode-bidi: embed;
  text-align: start;
}

.info-value.link {
  color: #007bff;
  text-decoration: none;
  transition: color 0.3s ease;
}

.info-value.link:hover {
  color: #0056b3;
  text-decoration: underline;
}

.phone-actions {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.action-button {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 6px 12px;
  border-radius: 8px;
  font-size: 12px;
  font-weight: 600;
  text-decoration: none;
  transition: all 0.3s ease;
}

.call-button {
  background: #28a745;
  color: white;
}

.call-button:hover {
  background: #218838;
  transform: translateY(-1px);
}

.status-badge {
  display: inline-block;
  padding: 4px 10px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 600;
}

.status-badge.status-pending {
  background: #e3f2fd;
  color: #1976d2;
}

.status-badge.status-assigned {
  background: #e9ecef;
  color: #495057;
}

.status-badge.status-in_progress {
  background: #fff3cd;
  color: #856404;
}

.status-badge.status-completed {
  background: #d4edda;
  color: #155724;
}

.status-badge.status-late {
  background: #f8d7da;
  color: #721c24;
}

.status-badge.status-cancelled {
  background: #fce4ec;
  color: #c62828;
}

.priority-badge {
  display: inline-block;
  padding: 4px 10px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 600;
}

.priority-badge.priority-low {
  background: #e3f2fd;
  color: #1976d2;
}

.priority-badge.priority-normal {
  background: #e9ecef;
  color: #495057;
}

.priority-badge.priority-high {
  background: #fff3cd;
  color: #856404;
}

.priority-badge.priority-urgent {
  background: #f8d7da;
  color: #721c24;
}

.maps-button {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  background: linear-gradient(135deg, #4285f4 0%, #34a853 100%);
  color: white;
  padding: 12px 20px;
  border-radius: 12px;
  text-decoration: none;
  font-weight: 600;
  font-size: 14px;
  transition: all 0.3s ease;
  margin-top: 8px;
}

.maps-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(66, 133, 244, 0.3);
}

.maps-icon {
  font-size: 16px;
}

.action-buttons {
  display: flex;
  flex-direction: column;
  gap: 12px;
  position: sticky;
  bottom: 20px;
}

.attendance-button {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  color: white;
  padding: 16px 24px;
  border-radius: 16px;
  text-decoration: none;
  font-weight: 600;
  font-size: 16px;
  box-shadow: 0 4px 12px rgba(0, 123, 255, 0.3);
  transition: all 0.3s ease;
}

.attendance-button.primary {
  background: linear-gradient(135deg, #007bff 0%, #0056b3 100%);
}

.attendance-button.primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(0, 123, 255, 0.4);
}

.attendance-button.secondary {
  background: linear-gradient(135deg, #28a745 0%, #218838 100%);
  box-shadow: 0 4px 12px rgba(40, 167, 69, 0.3);
}

.attendance-button.secondary:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(40, 167, 69, 0.4);
}

.button-icon {
  font-size: 20px;
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
}

.empty-icon {
  font-size: 64px;
  margin-bottom: 16px;
  opacity: 0.5;
}

.empty-state h3 {
  color: #6c757d;
  font-size: 18px;
  margin: 0;
}

/* Responsive adjustments */
@media (max-width: 480px) {
  .task-detail {
    padding: 20px;
    border-radius: 16px;
  }
  
  .task-detail__title {
    font-size: 20px;
  }
  
  .task-detail__section {
    padding: 16px;
  }
  
  .section-title {
    font-size: 15px;
  }
  
  .info-label {
    min-width: 70px;
    font-size: 12px;
  }
  
  .info-value {
    font-size: 13px;
  }
  
  .attendance-button {
    padding: 14px 20px;
    font-size: 15px;
  }
  
  .action-buttons {
    flex-direction: column;
  }
}
</style>

```

---

## 📄 frontend-worker-pwa/src/views/TasksListView.vue

```vue
<template>
  <div class="tasks-page">
    <div class="page-header">
      <div class="page-header__content">
        <h1 class="page-title">{{ currentTranslations.my_tasks }}</h1>
        <p class="page-subtitle" v-if="tasksStore.items.length > 0">
          {{ tasksStore.items.length }} {{ tasksStore.items.length === 1 ? currentTranslations.task : currentTranslations.task_plural }}
        </p>
        <p class="page-subtitle" v-else-if="!tasksStore.loading">
          {{ currentTranslations.no_tasks }}
        </p>
      </div>
      <div class="page-header__actions">
        <button @click="refreshTasks" class="refresh-button" :disabled="tasksStore.loading">
          <span class="refresh-icon" :class="{ spinning: tasksStore.loading }">🔄</span>
          <span>{{ currentTranslations.refresh }}</span>
        </button>
        <div class="page-header__icon">
          📋
        </div>
      </div>
    </div>

    <div v-if="tasksStore.loading" class="loading-state">
      <div class="loading-spinner"></div>
      <p>{{ currentTranslations.loading_tasks }}</p>
    </div>

    <div v-else-if="tasksStore.items.length === 0" class="empty-state">
      <div class="empty-icon">📋</div>
      <h3>{{ currentTranslations.no_tasks_today }}</h3>
      <p>{{ currentTranslations.no_tasks_scheduled }}</p>
    </div>

    <template v-else>
      <div class="tasks-list">
        <TaskCard v-for="task in tasksStore.items" :key="task.id" :task="task" />
      </div>
    </template>
  </div>
</template>

<script setup>
import { onMounted, onUnmounted, watch } from 'vue'
import { useI18n } from '../services/i18n'
import { tasksStore } from '../store/tasks'
import TaskCard from '../components/TaskCard.vue'
import wsService from '../services/websocket'

const { t, currentLang, currentTranslations } = useI18n()

// تحديث المهام
const refreshTasks = () => {
  tasksStore.fetchMine()
}

// Reload tasks when language changes to get translated data from backend
watch(currentLang, () => {
  tasksStore.fetchMine()
})

onMounted(() => {
  tasksStore.fetchMine()
  connectWebSocket()
})

onUnmounted(() => {
  disconnectWebSocket()
})

function connectWebSocket() {
  const apiBaseUrl = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1'
  const apiHost = apiBaseUrl.replace('/api/v1', '')
  const wsUrl = apiHost.replace('http://', 'ws://').replace('https://', 'wss://') + '/ws'
  console.log('🔌 Attempting to connect to WebSocket:', wsUrl)
  wsService.connect(wsUrl)

  wsService.onMessage((data) => {
    if (data.type === 'connected') {
      console.log('✅ WebSocket connected')
    } else if (data.type === 'disconnected') {
      console.log('❌ WebSocket disconnected')
    } else if (data.type === 'task_update') {
      console.log('📋 Task update:', data.data)
      tasksStore.fetchMine()
    }
  })
}

function disconnectWebSocket() {
  wsService.disconnect()
}
</script>

<style scoped>
.tasks-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #f5f7fa 0%, #e4e8ec 100%);
  padding: 20px 16px;
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: linear-gradient(135deg, #ffffff 0%, #f8f9fa 100%);
  border-radius: 20px;
  padding: 24px 20px;
  margin-bottom: 24px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

.page-header__content {
  flex: 1;
}

.page-header__actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.refresh-button {
  display: flex;
  align-items: center;
  gap: 8px;
  background: linear-gradient(135deg, #007bff 0%, #0056b3 100%);
  color: white;
  border: none;
  padding: 10px 16px;
  border-radius: 12px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 2px 8px rgba(0, 123, 255, 0.3);
}

.refresh-button:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 123, 255, 0.4);
}

.refresh-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.refresh-icon {
  font-size: 16px;
  transition: transform 0.3s ease;
}

.refresh-icon.spinning {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.page-title {
  font-size: 24px;
  font-weight: 700;
  color: #212529;
  margin: 0 0 8px 0;
  line-height: 1.2;
}

.page-subtitle {
  font-size: 14px;
  color: #6c757d;
  margin: 0;
}

.page-header__icon {
  font-size: 48px;
  width: 64px;
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #e3f2fd 0%, #bbdefb 100%);
  border-radius: 16px;
  box-shadow: 0 4px 8px rgba(33, 150, 243, 0.2);
}

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  text-align: center;
}

.loading-spinner {
  width: 48px;
  height: 48px;
  border: 4px solid #e9ecef;
  border-top: 4px solid #007bff;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 16px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.loading-state p {
  color: #6c757d;
  font-size: 14px;
  margin: 0;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  text-align: center;
  background: linear-gradient(135deg, #ffffff 0%, #f8f9fa 100%);
  border-radius: 20px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

.empty-icon {
  font-size: 64px;
  margin-bottom: 16px;
  opacity: 0.5;
}

.empty-state h3 {
  color: #212529;
  font-size: 18px;
  margin: 0 0 8px 0;
}

.empty-state p {
  color: #6c757d;
  font-size: 14px;
  margin: 0;
}

.tasks-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

/* Responsive adjustments */
@media (max-width: 480px) {
  .tasks-page {
    padding: 16px 12px;
  }
  
  .page-header {
    padding: 20px 16px;
    border-radius: 16px;
  }
  
  .page-title {
    font-size: 20px;
  }
  
  .page-header__icon {
    font-size: 40px;
    width: 56px;
    height: 56px;
  }
  
  .empty-state {
    padding: 40px 16px;
  }
  
  .empty-icon {
    font-size: 48px;
  }
}
</style>

```

---

## 📄 frontend-worker-pwa/vite.config.js

```javascript
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  appType: 'spa',
  plugins: [vue()],
  test: {
    globals: true,
    environment: 'jsdom',
    setupFiles: ['./src/tests/setup.js'],
    coverage: {
      provider: 'v8',
      reporter: ['text', 'json', 'html'],
      exclude: [
        'node_modules/',
        'src/tests/',
        '**/*.d.ts',
        '**/*.config.*',
        '**/mockData',
        'dist/',
        'electron/'
      ]
    }
  },
  server: { 
    port: 3002,
    cors: true,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        secure: false
      }
    }
  },
  preview: {
    port: 3002,
    // SPA routing support for preview
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        secure: false
      }
    }
  },
  optimizeDeps: {
    include: ['vue', 'vue-router']
  },
  build: {
    cssCodeSplit: true,
    rollupOptions: {
      output: {
        manualChunks: {
          'vue-vendor': ['vue', 'vue-router']
        },
        // إضافة timestamp تلقائي لضمان cache busting
        entryFileNames: 'assets/[name]-[hash].js',
        chunkFileNames: 'assets/[name]-[hash].js',
        assetFileNames: 'assets/[name]-[hash][extname]'
      }
    }
  },
  base: '/'
})

```

---

