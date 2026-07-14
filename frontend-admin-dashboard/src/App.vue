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
          <span class="nav-icon">📊</span>
          <span class="nav-label">{{ t('dashboard') }}</span>
        </router-link>
        <router-link to="/service-requests" class="nav-item" active-class="nav-item--active">
          <span class="nav-icon">📋</span>
          <span class="nav-label">{{ t('service_requests') }}</span>
        </router-link>
        <router-link to="/employees" class="nav-item" active-class="nav-item--active">
          <span class="nav-icon">👥</span>
          <span class="nav-label">{{ t('employees') }}</span>
        </router-link>
        <router-link to="/worksites" class="nav-item" active-class="nav-item--active">
          <span class="nav-icon">📍</span>
          <span class="nav-label">{{ t('worksites') }}</span>
        </router-link>
        <router-link to="/reports" class="nav-item" active-class="nav-item--active">
          <span class="nav-icon">📈</span>
          <span class="nav-label">{{ t('reports') }}</span>
        </router-link>
        <router-link to="/settings" class="nav-item" active-class="nav-item--active">
          <span class="nav-icon">⚙️</span>
          <span class="nav-label">{{ t('settings') }}</span>
        </router-link>
      </nav>

      <div class="sidebar__footer">
        <div class="sidebar__footer-top">
          <LanguageSwitcher />
          <button class="theme-toggle" @click="toggleTheme" :title="isDark ? t('theme_light') : t('theme_dark')">
            {{ isDark ? '☀️' : '🌙' }}
          </button>
        </div>
        <button class="btn-logout" @click="handleLogout">
          <span class="logout-icon">🚪</span>
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
          <span class="topbar__avatar">{{ initials }}</span>
        </div>
      </header>
      <main class="content"><router-view /></main>
    </div>
  </div>
</template>

<script setup>
import { computed, ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { authStore } from './store/auth'
import LanguageSwitcher from './components/LanguageSwitcher.vue'
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
    '/dashboard': `📊 ${t('dashboard')}`,
    '/service-requests': `📋 ${t('service_requests')}`,
    '/employees': `👥 ${t('employees')}`,
    '/worksites': `📍 ${t('worksites')}`,
    '/reports': `📈 ${t('reports')}`,
    '/settings': `⚙️ ${t('settings')}`,
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
  box-shadow: var(--shadow-sm);
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
  transition: all var(--transition-fast);
}

.nav-item:hover {
  background: var(--brand-tint);
  color: var(--brand);
}

.nav-item--active {
  background: var(--brand);
  color: white;
  box-shadow: var(--shadow-sm);
}

.nav-item--active:hover {
  background: var(--brand-dark);
  color: white;
}

.nav-icon {
  font-size: 18px;
  width: 28px;
  text-align: center;
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
  font-size: 18px;
  cursor: pointer;
  transition: all var(--transition-fast);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.theme-toggle:hover {
  background: var(--brand-tint);
  border-color: var(--brand);
  transform: scale(1.05);
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
  transition: all var(--transition-fast);
  width: 100%;
}

.btn-logout:hover {
  background: var(--signal-out);
  color: white;
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
}

.logout-icon { font-size: 16px; }
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
}

.topbar__title {
  font-size: 18px;
  font-weight: 600;
}

.topbar__user {
  display: flex;
  align-items: center;
  gap: 10px;
}

.topbar__name {
  font-size: 13px;
  color: var(--ink-soft);
}

.topbar__avatar {
  width: 34px;
  height: 34px;
  border-radius: 50%;
  background: var(--brand-tint);
  color: var(--brand);
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  font-size: 14px;
}

.content {
  flex: 1;
  padding: 24px 28px;
  max-width: 1280px;
  width: 100%;
  margin: 0 auto;
}

/* =============================================
   استجابة للشاشات الصغيرة
   ============================================= */
@media (max-width: 768px) {
  .sidebar {
    width: 100px;
    min-width: 100px;
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
