<template>
  <div v-if="isLogin" class="auth-shell">
    <router-view />
  </div>

  <div v-else class="shell">
    <aside class="sidebar">
      <div class="sidebar__brand">
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
            <LayoutDashboard :size="20" />
          </span>
          <span class="nav-label">{{ t('dashboard') }}</span>
        </router-link>
        <router-link to="/employees" class="nav-item" active-class="nav-item--active">
          <span class="nav-icon">
            <Users :size="20" />
          </span>
          <span class="nav-label">{{ t('employees') }}</span>
        </router-link>
        <router-link to="/worksites" class="nav-item" active-class="nav-item--active">
          <span class="nav-icon">
            <MapPin :size="20" />
          </span>
          <span class="nav-label">{{ t('worksites') }}</span>
        </router-link>
        <router-link to="/reports" class="nav-item" active-class="nav-item--active">
          <span class="nav-icon">
            <BarChart3 :size="20" />
          </span>
          <span class="nav-label">{{ t('reports') }}</span>
        </router-link>
        <router-link to="/attendance-management" class="nav-item" active-class="nav-item--active">
          <span class="nav-icon">
            <Clock :size="20" />
          </span>
          <span class="nav-label">{{ t('attendance_management') }}</span>
        </router-link>
        <router-link to="/notes" class="nav-item" active-class="nav-item--active">
          <span class="nav-icon">
            <FileText :size="20" />
          </span>
          <span class="nav-label">{{ t('notes') }}</span>
        </router-link>
        <router-link to="/settings" class="nav-item" active-class="nav-item--active">
          <span class="nav-icon">
            <Settings :size="20" />
          </span>
          <span class="nav-label">{{ t('settings') }}</span>
        </router-link>
      </nav>

      <div class="sidebar__footer">
        <div class="sidebar__footer-top">
          <LanguageSwitcher />
          <button class="theme-toggle" @click="toggleTheme" :title="isDark ? t('theme_light') : t('theme_dark')">
            <Moon v-if="!isDark" :size="18" />
            <Sun v-else :size="18" />
          </button>
        </div>
        <button class="btn-logout" @click="handleLogout">
          <span class="logout-icon">
            <LogOut :size="16" />
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
</template>

<script setup>
import { computed, ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { authStore } from './store/auth'
import LanguageSwitcher from './components/LanguageSwitcher.vue'
import { useI18n } from './services/i18n'
import {
  LayoutDashboard,
  Users,
  MapPin,
  BarChart3,
  Clock,
  FileText,
  Settings,
  Moon,
  Sun,
  LogOut
} from '@lucide/vue'

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
    '/employees': t('employees'),
    '/worksites': t('worksites'),
    '/reports': t('reports'),
    '/settings': t('settings'),
    '/attendance-management': t('attendance_management'),
    '/notes': t('notes'),
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
  background: linear-gradient(180deg, #f8fafc 0%, #e2e8f0 100%);
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

.app-brand {
  display: flex;
  align-items: center;
  gap: 10px;
}

.brand-mark {
  width: 48px;
  height: 48px;
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
  color: #000000;
}

[data-theme="dark"] .nav-item--active .nav-icon {
  color: var(--gold);
}

/* تلوين أيقونات الشريط الجانبي بالأسود */
.nav-icon :deep(svg) {
  color: #000000 !important;
  stroke: #000000 !important;
}

.nav-icon :deep(svg *) {
  stroke: #000000 !important;
}

.nav-item:hover .nav-icon :deep(svg) {
  color: #000000 !important;
  stroke: #000000 !important;
}

.nav-item:hover .nav-icon :deep(svg *) {
  stroke: #000000 !important;
}

/* إزالة الخلفية الزرقاء في الوضع الليلي */
[data-theme="dark"] .nav-item:hover {
  background: transparent !important;
}

[data-theme="dark"] .nav-item--active {
  background: transparent !important;
  color: var(--gold) !important;
  box-shadow: none !important;
}

[data-theme="dark"] .nav-item--active:hover {
  background: transparent !important;
  color: var(--gold) !important;
  box-shadow: none !important;
}

/* تلوين أيقونات الشريط الجانبي بالأسود في الوضع الليلي */
[data-theme="dark"] .nav-icon :deep(svg) {
  color: var(--gold) !important;
  stroke: var(--gold) !important;
}

[data-theme="dark"] .nav-icon :deep(svg *) {
  stroke: var(--gold) !important;
}

[data-theme="dark"] .nav-item:hover .nav-icon :deep(svg) {
  color: var(--gold-light) !important;
  stroke: var(--gold-light) !important;
}

[data-theme="dark"] .nav-item:hover .nav-icon :deep(svg *) {
  stroke: var(--gold-light) !important;
}

.nav-label {
  flex: 1;
}

/* Dark mode sidebar gradient */
[data-theme="dark"] .sidebar {
  background: linear-gradient(180deg, #334155 0%, #1e293b 100%);
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
  background: rgba(212, 175, 55, 0.15);
  cursor: pointer;
  transition: all var(--transition-base);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  color: var(--gold);
}

.theme-toggle:hover {
  background: rgba(212, 175, 55, 0.25);
  border-color: var(--gold);
  transform: scale(1.1);
  color: var(--gold);
  box-shadow: var(--shadow-md), 0 0 15px rgba(212, 175, 55, 0.3);
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
  color: var(--gold);
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
  color: var(--gold);
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
  .brand-sub,
  .footer-brand {
    display: none;
  }

  .sidebar__brand {
    padding: 1px 1px;
    justify-content: flex-end;
  }

  .app-brand {
    justify-content: flex-end;
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
