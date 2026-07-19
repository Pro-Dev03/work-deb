<template>
  <div v-if="isLogin" class="auth-shell">
    <router-view />
  </div>

  <div v-else class="shell">
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
          <img src="/src/assets/company-logo.jpg" alt="WorkTrack logo" class="brand-mark" />
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
          <span class="topbar__name">{{ displayName }}</span>
        </div>
      </header>
      <main class="content"><router-view /></main>
    </div>
  </div>
  
  <!-- PWA Install Button -->
  <PWAInstallButton />
</template>

<script setup>
import { computed, ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { authStore } from './store/auth'
import LanguageSwitcher from './components/LanguageSwitcher.vue'
import PWAInstallButton from './components/PWAInstallButton.vue'
import { useI18n } from './services/i18n'

const { t, currentLang } = useI18n()

const route = useRoute()
const router = useRouter()
const isLogin = computed(() => route.path === '/login')
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
    '/worksites': t('worksites'),
    '/reports': t('reports'),
    '/settings': t('settings'),
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
  const savedTheme = localStorage.getItem('worktrack_theme')
  if (savedTheme === 'dark') {
    isDark.value = true
    document.documentElement.setAttribute('data-theme', 'dark')
  }
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
  gap: 10px;
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
