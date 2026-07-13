<template>
  <div>
    <div class="profile-head">
      <span class="profile-head__avatar">{{ initials }}</span>
      <h2>{{ user?.full_name || t('app_name') }}</h2>
      <p>{{ user?.role === 'admin' ? t('admin') : t('employee') }}</p>
    </div>

    <div class="card profile-menu">
      <!-- ✅ محدد اللغة - موحد -->
      <div class="profile-menu__item language-selector">
        <span class="lang-label">🌐 {{ t('language') }}</span>
        <div class="lang-buttons">
          <button 
            v-for="lang in languages" 
            :key="lang.code"
            class="lang-btn"
            :class="{ active: currentLang === lang.code }"
            @click="changeLanguage(lang.code)"
            :title="lang.name"
          >
            {{ lang.flag }}
            <span class="lang-name">{{ lang.name }}</span>
          </button>
        </div>
      </div>

      <button class="profile-menu__item" @click="showNotifications = !showNotifications">
        🔔 {{ t('notifications') }}
      </button>
      <button class="profile-menu__item" @click="showHistory = !showHistory">
        📄 {{ t('attendance_history') }}
      </button>
    </div>

    <!-- الإشعارات -->
    <div v-if="showNotifications" class="card notifications-card">
      <h4>🔔 {{ t('notifications') }}</h4>
      <div v-if="notifications.length === 0" class="empty-state">
        <p>{{ t('no_notifications') }}</p>
      </div>
      <div v-else v-for="notif in notifications" :key="notif.id" class="notification-item">
        <p class="notif-title">{{ notif.title }}</p>
        <p class="notif-body">{{ notif.body }}</p>
        <span class="notif-time mono">{{ formatDate(notif.created_at) }}</span>
      </div>
    </div>

    <!-- سجل الحضور -->
    <div v-if="showHistory" class="card history-card">
      <h4>📄 {{ t('attendance_history') }}</h4>
      <div v-if="attendanceHistory.length === 0" class="empty-state">
        <p>{{ t('no_history') }}</p>
      </div>
      <div v-else v-for="record in attendanceHistory" :key="record.id" class="history-item">
        <span class="history-date">{{ formatDate(record.date) }}</span>
        <span class="history-hours">{{ record.hours }} {{ t('hours') }}</span>
        <span class="history-status" :class="record.status === 'completed' ? 'status-done' : 'status-active'">
          {{ record.status === 'completed' ? t('completed') : t('in_progress') }}
        </span>
      </div>
    </div>

    <button class="btn btn--danger btn--block" @click="handleLogout">
      🚪 {{ t('logout') }}
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

const { t, currentLang, setLang } = useI18n()
const router = useRouter()

const user = computed(() => authStore.user)
const initials = computed(() => (user.value?.full_name || 'م ع').trim().slice(0, 1))

const showNotifications = ref(false)
const showHistory = ref(false)
const notifications = ref([])
const attendanceHistory = ref([])

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
    console.error('فشل جلب الإشعارات:', error)
  }
}

async function fetchAttendanceHistory() {
  try {
    const { data } = await api.get('/attendance/history')
    attendanceHistory.value = data || []
  } catch (error) {
    console.error('فشل جلب سجل الحضور:', error)
  }
}

onMounted(() => {
  fetchNotifications()
  fetchAttendanceHistory()
})
</script>

<style scoped>
.profile-head {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 24px 0;
}

.profile-head__avatar {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  background: var(--brand-tint);
  color: var(--brand-dark);
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  font-size: 22px;
}

.profile-head h2 {
  font-size: 17px;
}

.profile-head p {
  font-size: 13px;
  color: var(--ink-soft);
}

.profile-menu {
  display: flex;
  flex-direction: column;
  margin: 20px 0;
  overflow: hidden;
}

.profile-menu__item {
  text-align: right;
  padding: 14px 18px;
  border: none;
  background: none;
  font-family: var(--font-body);
  font-size: 14px;
  color: var(--ink);
  border-bottom: 1px solid var(--line);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.profile-menu__item:last-child {
  border-bottom: none;
}

.profile-menu__item:hover {
  background: var(--brand-tint);
}

/* محدد اللغة */
.language-selector {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 18px;
  border-bottom: 1px solid var(--line);
}

.lang-label {
  font-size: 14px;
  color: var(--ink);
}

.lang-buttons {
  display: flex;
  gap: 4px;
  align-items: center;
}

.lang-btn {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 4px 10px;
  border: none;
  border-radius: 6px;
  background: transparent;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.3s ease;
  font-family: inherit;
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
}

.lang-name {
  font-size: 11px;
  font-weight: 500;
}

@media (max-width: 480px) {
  .lang-name {
    display: none;
  }
  .lang-btn {
    padding: 4px 6px;
    font-size: 14px;
  }
}

.notifications-card,
.history-card {
  padding: 16px 18px;
  margin-bottom: 12px;
}

.notifications-card h4,
.history-card h4 {
  font-size: 15px;
  margin-bottom: 10px;
  color: var(--brand);
}

.notification-item,
.history-item {
  padding: 10px 0;
  border-bottom: 1px solid var(--line);
}

.notification-item:last-child,
.history-item:last-child {
  border-bottom: none;
}

.notif-title {
  font-weight: 600;
  font-size: 13px;
  color: var(--ink);
}

.notif-body {
  font-size: 12px;
  color: var(--ink-soft);
  margin: 2px 0;
}

.notif-time {
  font-size: 11px;
  color: var(--ink-light);
}

.history-date {
  font-size: 13px;
  color: var(--ink);
}

.history-hours {
  font-size: 13px;
  color: var(--brand);
  margin: 0 8px;
}

.history-status {
  font-size: 12px;
  font-weight: 600;
  padding: 2px 10px;
  border-radius: 999px;
}

.status-done {
  background: var(--signal-in-tint);
  color: var(--signal-in);
}

.status-active {
  background: var(--signal-warning-tint);
  color: var(--signal-warning);
}

.empty-state {
  text-align: center;
  padding: 20px;
  color: var(--ink-soft);
  font-size: 13px;
}

.btn--block {
  width: 100%;
  justify-content: center;
}
</style>
