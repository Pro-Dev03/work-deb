<template>
  <div class="profile-view view">
    <!-- Profile Header -->
    <div class="profile-header">
      <div class="avatar avatar--72">
        <svg class="person-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
          <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
          <circle cx="12" cy="7" r="4"/>
        </svg>
      </div>
      <h2>{{ user?.full_name || t('app_name') }}</h2>
      <p>{{ user?.role === 'admin' ? t('admin') : t('employee') }}</p>
    </div>

    <!-- Menu Section -->
    <div class="section">
      <div class="section-label">{{ t('settings') }}</div>
      <div class="menu-list">
        <button class="menu-item" @click="showNotifications = !showNotifications">
          <div class="m-icon">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M18 8a6 6 0 0 0-12 0c0 7-3 9-3 9h18s-3-2-3-9"/>
              <path d="M13.73 21a2 2 0 0 1-3.46 0"/>
            </svg>
          </div>
          <span>{{ t('notifications') }}</span>
          <svg class="chev" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="15 18 9 12 15 6"/>
          </svg>
        </button>

        <button class="menu-item menu-item--danger" @click="handleLogout">
          <div class="m-icon">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/>
              <polyline points="16 17 21 12 16 7"/>
              <line x1="21" y1="12" x2="9" y2="12"/>
            </svg>
          </div>
          <span>{{ t('logout') }}</span>
        </button>
      </div>
    </div>

    <!-- Notifications Section -->
    <div v-if="showNotifications" class="card">
      <h4>🔔 {{ t('notifications') }}</h4>
      <div v-if="notifications.length === 0" class="empty-state">
        <p>{{ t('no_notifications') }}</p>
      </div>
      <div v-else v-for="notif in notifications" :key="notif.id" class="notification-item">
        <p class="notif-title">{{ notif.title || t('note') }}</p>
        <p class="notif-body">{{ notif.content || notif.body }}</p>
        <span class="notif-time mono">{{ formatDate(notif.created_at) }}</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from '../services/i18n'
import { authStore } from '../store/auth'
import { logout } from '../services/auth'
import api from '../services/api'

const { t, currentLang, setLang } = useI18n()
const router = useRouter()

const user = computed(() => authStore.user)
const initials = computed(() => (user.value?.full_name || t('initials')).trim().slice(0, 1))

const showNotifications = ref(false)
const notifications = ref([])

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
    const { data } = await api.get('/notes')
    notifications.value = data || []
  } catch (error) {
    console.error(t('failed_fetch_notifications'), error)
  }
}

onMounted(() => {
  fetchNotifications()
})
</script>

<style scoped>
.profile-view {
  padding: var(--space-4);
  position: relative;
}

.profile-header {
  text-align: center;
  padding: var(--space-6) var(--space-4);
  background: linear-gradient(160deg, var(--primary-100) 0%, var(--background) 100%);
  border-radius: var(--radius-lg);
  margin-bottom: var(--space-5);
}

[data-theme="dark"] .profile-header {
  background: linear-gradient(160deg, rgba(99, 102, 241, 0.16) 0%, var(--background) 100%);
}

.avatar {
  width: 72px;
  height: 72px;
  border-radius: 50%;
  background: linear-gradient(135deg, var(--primary-500) 0%, var(--primary-600) 100%);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: var(--font-bold);
  font-size: var(--text-xl);
  box-shadow: var(--shadow-md);
  margin: 0 auto var(--space-3);
}

.person-icon {
  width: 36px;
  height: 36px;
  color: white;
}

.profile-header h2 {
  font-size: var(--text-xl);
  font-weight: var(--font-semibold);
  color: var(--text-primary);
  margin: var(--space-3) 0 var(--space-1);
}

.profile-header p {
  font-size: var(--text-sm);
  color: var(--text-secondary);
}

.section {
  margin-bottom: var(--space-5);
}

.menu-item {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  padding: var(--space-3) var(--space-4);
  border-bottom: 1px solid var(--border);
  cursor: pointer;
  transition: var(--transition-base) ease;
  background: transparent;
  border: none;
  width: 100%;
  text-align: start;
}

.menu-item:last-child {
  border-bottom: none;
}

.menu-item:hover {
  background: var(--surface-elevated);
}

.menu-item--danger {
  color: var(--error-600);
}

.menu-item--danger .m-icon {
  background: var(--error-50);
  color: var(--error-600);
}

.m-icon {
  width: 32px;
  height: 32px;
  border-radius: var(--radius-md);
  background: var(--gray-100);
  color: var(--text-secondary);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.menu-item span {
  flex: 1;
  font-size: var(--text-sm);
  font-weight: var(--font-medium);
  color: var(--text-primary);
}

.chev {
  color: var(--text-tertiary);
}

/* Notifications Styles */
.notification-item {
  padding: var(--space-3) 0;
  border-bottom: 1px solid var(--border);
}

.notification-item:last-child {
  border-bottom: none;
}

.notif-title {
  font-weight: var(--font-medium);
  font-size: var(--text-sm);
  color: var(--text-primary);
  margin-bottom: var(--space-1);
}

.notif-body {
  font-size: var(--text-sm);
  color: var(--text-secondary);
  margin-bottom: var(--space-1);
}

.notif-time {
  font-size: var(--text-xs);
  color: var(--text-tertiary);
}

.status-done {
  background: var(--primary-50);
  color: var(--primary-600);
}

.status-active {
  background: var(--warning-50);
  color: var(--warning-600);
}

.empty-state {
  text-align: center;
  padding: var(--space-4);
  color: var(--text-tertiary);
  font-size: var(--text-sm);
}

.card h4 {
  font-size: var(--text-base);
  font-weight: var(--font-semibold);
  color: var(--text-primary);
  margin-bottom: var(--space-3);
}

@media (min-width: 768px) {
  .profile-header {
    padding: var(--space-8) var(--space-6);
  }
  
  .fab-whatsapp {
    bottom: 100px;
    right: 32px;
  }
}
</style>
