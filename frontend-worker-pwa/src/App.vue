<template>
  <div class="app-shell">
    <!-- Header -->
    <header v-if="!isLoginPage" class="app-header">
      <div class="brand">
        <div class="brand-icon">
          <img src="/src/assets/company-logo.jpg" alt="Company Logo" class="brand-logo" />
        </div>
        <div>
          <div class="brand-name">{{ t('app_name') }}</div>
          <div class="brand-sub">Powered by DevPro</div>
        </div>
      </div>
      <div class="header-actions">
        <router-link 
          to="/notes" 
          class="icon-btn" 
          :class="{ 'has-notifications': hasNotifications }"
          title="الملاحظات" 
          aria-label="الملاحظات"
        >
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M18 8a6 6 0 0 0-12 0c0 7-3 9-3 9h18s-3-2-3-9"/>
            <path d="M13.73 21a2 2 0 0 1-3.46 0"/>
          </svg>
          <span v-if="hasNotifications" class="notification-badge">{{ unreadCount }}</span>
        </router-link>
        <router-link to="/profile" class="avatar avatar--40">
          <svg class="person-icon-header" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
            <circle cx="12" cy="7" r="4"/>
          </svg>
        </router-link>
      </div>
    </header>

    <!-- Main Content -->
    <main class="app-main">
      <router-view />
    </main>

    <!-- Bottom Navigation -->
    <nav v-if="!isLoginPage" class="bottom-nav">
      <router-link to="/attendance" class="nav-item" active-class="nav-item--active">
        <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <circle cx="12" cy="12" r="10"/>
          <polyline points="12 6 12 12 16 14"/>
        </svg>
        <span>{{ t('attendance') }}</span>
      </router-link>
      <router-link to="/notes" class="nav-item" active-class="nav-item--active">
        <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/>
        </svg>
        <span>{{ t('notes') }}</span>
      </router-link>
      <router-link to="/history" class="nav-item" active-class="nav-item--active">
        <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
          <polyline points="14 2 14 8 20 8"/>
          <line x1="16" y1="13" x2="8" y2="13"/>
          <line x1="16" y1="17" x2="8" y2="17"/>
          <polyline points="10 9 9 9 8 9"/>
        </svg>
        <span>{{ t('history') }}</span>
      </router-link>
      <router-link to="/profile" class="nav-item" active-class="nav-item--active">
        <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
          <circle cx="12" cy="7" r="4"/>
        </svg>
        <span>{{ t('profile') }}</span>
      </router-link>
    </nav>
  </div>
</template>

<script setup>
import { computed, ref, onMounted, onUnmounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from './services/i18n'
import { authStore } from './store/auth'
import api from './services/api'

const route = useRoute()
const { t } = useI18n()
const isLoginPage = computed(() => route.path === '/login')
const isNotesPage = computed(() => route.path === '/notes')

const unreadCount = ref(0)
const hasNotifications = computed(() => unreadCount.value > 0 && !isNotesPage.value)
const previousUnreadCount = ref(0)
let currentAudioContext = null
let currentOscillator = null

function playNotificationSound() {
  // Create a simple beep sound using Web Audio API
  const audioContext = new (window.AudioContext || window.webkitAudioContext)()
  const oscillator = audioContext.createOscillator()
  const gainNode = audioContext.createGain()
  
  oscillator.connect(gainNode)
  gainNode.connect(audioContext.destination)
  
  oscillator.frequency.value = 800
  oscillator.type = 'sine'
  
  gainNode.gain.setValueAtTime(0.3, audioContext.currentTime)
  gainNode.gain.exponentialRampToValueAtTime(0.01, audioContext.currentTime + 0.5)
  
  oscillator.start(audioContext.currentTime)
  oscillator.stop(audioContext.currentTime + 0.5)
  
  currentAudioContext = audioContext
  currentOscillator = oscillator
}

function stopNotificationSound() {
  if (currentOscillator) {
    try {
      currentOscillator.stop()
    } catch (e) {
      // Oscillator might already be stopped
    }
    currentOscillator = null
  }
  if (currentAudioContext) {
    try {
      currentAudioContext.close()
    } catch (e) {
      // Context might already be closed
    }
    currentAudioContext = null
  }
}

async function fetchNotifications() {
  try {
    const { data } = await api.get('/notes')
    const notes = Array.isArray(data) ? data : []
    const newUnreadCount = notes.filter(n => !n.is_read).length
    
    // Play sound immediately if new notifications arrived
    if (newUnreadCount > previousUnreadCount.value && previousUnreadCount.value >= 0) {
      playNotificationSound()
    }
    
    unreadCount.value = newUnreadCount
    previousUnreadCount.value = newUnreadCount
  } catch (error) {
    console.error('Failed to fetch notifications:', error)
    unreadCount.value = 0
  }
}

let notificationInterval
let soundInterval

onMounted(() => {
  fetchNotifications()
  // Check for new notifications every 30 seconds
  notificationInterval = setInterval(fetchNotifications, 30000)
  
  // Play notification sound every 5 seconds if there are unread notifications
  soundInterval = setInterval(() => {
    if (hasNotifications.value) {
      playNotificationSound()
    }
  }, 5000)
})

onUnmounted(() => {
  if (notificationInterval) {
    clearInterval(notificationInterval)
  }
  if (soundInterval) {
    clearInterval(soundInterval)
  }
})

// Stop sound immediately when entering notes page
watch(isNotesPage, (newValue) => {
  if (newValue) {
    stopNotificationSound()
  }
})
</script>

<style scoped>
.app-shell {
  min-height: 100dvh;
  display: flex;
  flex-direction: column;
  background: var(--background);
}

/* Header - SaaS Style */
.app-header {
  position: sticky;
  top: 0;
  z-index: var(--z-sticky);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--space-3) var(--space-4);
  background: var(--surface);
  border-bottom: 1px solid var(--border);
}

.brand {
  display: flex;
  align-items: center;
  gap: 8px;
}

.brand-icon {
  width: 32px;
  height: 32px;
  border-radius: var(--radius-md);
  background: var(--primary-500);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
}

.brand-logo {
  width: 32px;
  height: 32px;
  border-radius: var(--radius-md);
  object-fit: cover;
}

.brand-name {
  font-weight: var(--font-semibold);
  font-size: var(--text-base);
  line-height: 1.1;
  color: var(--text-primary);
}

.brand-sub {
  font-size: 10px;
  color: var(--text-tertiary);
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.icon-btn {
  border: none;
  background: transparent;
  cursor: pointer;
  color: var(--text-secondary);
  width: 32px;
  height: 32px;
  min-width: 32px;
  min-height: 32px;
  border-radius: var(--radius-md);
  display: flex;
  align-items: center;
  justify-content: center;
  transition: var(--transition-base) ease;
  position: relative;
}

.icon-btn:hover {
  background: var(--gray-100);
  color: var(--text-primary);
}

.icon-btn.has-notifications {
  color: #FFD700;
  animation: bell-shake 0.5s ease-in-out;
}

.icon-btn.has-notifications:hover {
  color: #FFC700;
  background: rgba(255, 215, 0, 0.1);
}

@keyframes bell-shake {
  0%, 100% {
    transform: rotate(0deg);
  }
  25% {
    transform: rotate(15deg);
  }
  75% {
    transform: rotate(-15deg);
  }
}

.notification-badge {
  position: absolute;
  top: -4px;
  right: -4px;
  background: #FF0000;
  color: white;
  font-size: 10px;
  font-weight: bold;
  width: 18px;
  height: 18px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  animation: badge-pulse 2s ease-in-out infinite;
}

@keyframes badge-pulse {
  0%, 100% {
    transform: scale(1);
    box-shadow: 0 0 0 0 rgba(255, 0, 0, 0.7);
  }
  50% {
    transform: scale(1.1);
    box-shadow: 0 0 0 4px rgba(255, 0, 0, 0.4);
  }
}

.avatar {
  width: 34px;
  height: 34px;
  border-radius: 50%;
  background: linear-gradient(135deg, var(--primary-500) 0%, var(--primary-600) 100%);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  text-decoration: none;
  box-shadow: var(--shadow-sm);
  transition: all var(--transition-base);
}

.avatar:hover {
  transform: scale(1.1);
  box-shadow: var(--shadow-md);
}

.person-icon-header {
  width: 18px;
  height: 18px;
  color: white;
}

/* Main Content */
.app-main {
  flex: 1;
  overflow-y: auto;
  padding: var(--space-4) var(--space-4) 80px;
}

/* Bottom Navigation - Clean & Professional */
.bottom-nav {
  position: sticky;
  bottom: 0;
  z-index: var(--z-sticky);
  display: flex;
  background: var(--surface);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  border-top: 1px solid var(--border);
  box-shadow: 0 -2px 10px rgba(0, 0, 0, 0.05);
  padding: 8px 8px calc(8px + env(safe-area-inset-bottom));
}

[data-theme="dark"] .bottom-nav {
  background: rgba(30, 30, 30, 0.95);
  box-shadow: 0 -2px 10px rgba(0, 0, 0, 0.3);
}

.nav-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  padding: 8px 4px;
  border-radius: 12px;
  cursor: pointer;
  border: none;
  background: transparent;
  color: var(--text-tertiary);
  font-size: 10px;
  font-weight: var(--font-medium);
  transition: all 0.2s ease;
  text-decoration: none;
}

.nav-item svg {
  transition: all 0.2s ease;
}

.nav-item:hover {
  color: var(--text-secondary);
  background: var(--gray-50);
}

.nav-item:hover svg {
  transform: translateY(-1px);
}

[data-theme="dark"] .nav-item:hover {
  background: rgba(255, 255, 255, 0.05);
}

.nav-item--active {
  color: var(--primary-600);
  background: var(--primary-50);
}

.nav-item--active svg {
  transform: translateY(-1px);
}

[data-theme="dark"] .nav-item--active {
  background: rgba(129, 140, 248, 0.15);
  color: var(--primary-400);
}

.nav-item--active span {
  font-weight: var(--font-semibold);
}

/* Responsive Design */
@media (min-width: 768px) {
  .app-main {
    padding: var(--space-6) var(--space-5) 100px;
  }

  .bottom-nav {
    padding: 10px 10px calc(10px + env(safe-area-inset-bottom));
  }

  .nav-item {
    font-size: var(--text-xs);
    padding: 10px 8px;
  }

  .nav-item svg {
    width: 24px;
    height: 24px;
  }
}

@media (min-width: 1024px) {
  .app-main {
    max-width: 1000px;
    margin: 0 auto;
    padding: var(--space-8) var(--space-6) 120px;
  }

  .app-header {
    padding: var(--space-4) var(--space-6);
  }

  .bottom-nav {
    padding: 12px 12px calc(12px + env(safe-area-inset-bottom));
  }

  .nav-item svg {
    width: 26px;
    height: 26px;
  }

  .nav-item {
    padding: 12px 10px;
  }
}
</style>
