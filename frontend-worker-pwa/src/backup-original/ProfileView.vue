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

    <!-- زر واتساب العائم -->
    <button 
      class="whatsapp-float-btn" 
      @click="openWhatsApp"
      title="تواصل عبر واتساب"
    >
      <svg class="whatsapp-icon" viewBox="0 0 32 32" fill="currentColor">
        <path d="M16 2C8.268 2 2 8.268 2 16c0 2.52.666 4.93 1.84 7.094L2.5 29.5l6.562-1.312A13.94 13.94 0 0 0 16 30c7.732 0 14-6.268 14-14S23.732 2 16 2zm0 25.5c-2.234 0-4.39-.586-6.297-1.688l-.453-.266-4.75.95.996-4.625-.281-.469A11.38 11.38 0 0 1 4.5 16c0-6.344 5.156-11.5 11.5-11.5S27.5 9.656 27.5 16 22.344 27.5 16 27.5zm6.344-8.656c-.344-.172-2.031-1-2.344-1.109-.313-.109-.531-.172-.75.172-.219.344-.844 1.109-1.031 1.344-.188.234-.375.266-.719.094-.344-.172-1.453-.531-2.766-1.703-1.031-.906-1.719-2.031-1.906-2.375-.188-.344-.016-.531.156-.688.156-.156.344-.406.516-.609.172-.203.234-.344.344-.578.109-.234.055-.438-.027-.609-.082-.172-.75-1.813-1.031-2.484-.266-.641-.547-.555-.75-.563-.188-.008-.406-.008-.625-.008-.219 0-.578.082-.875.406-.297.328-1.125 1.094-1.125 2.672 0 1.578 1.148 3.094 1.313 3.328.164.234 2.266 3.453 5.484 4.844 2.156.922 2.594.734 3.063.688.469-.047 1.5-.609 1.703-1.203.203-.594.203-1.109.141-1.234-.063-.125-.234-.203-.578-.375z"/>
      </svg>
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
const initials = computed(() => (user.value?.full_name || t('initials')).trim().slice(0, 1))

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

function openWhatsApp() {
  const phoneNumber = '0584838136'
  const cleanPhone = phoneNumber.replace(/\D/g, '')
  const whatsappUrl = `https://wa.me/${cleanPhone}`
  window.open(whatsappUrl, '_blank')
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
    console.error(t('failed_fetch_notifications'), error)
  }
}

async function fetchAttendanceHistory() {
  try {
    const { data } = await api.get('/attendance/history')
    attendanceHistory.value = data || []
  } catch (error) {
    console.error(t('failed_fetch_attendance_history'), error)
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

/* زر واتساب العائم */
.whatsapp-float-btn {
  position: fixed;
  bottom: 24px;
  right: 24px;
  width: 60px;
  height: 60px;
  border-radius: 50%;
  background: linear-gradient(135deg, #25D366 0%, #128C7E 100%);
  border: none;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 20px rgba(37, 211, 102, 0.4);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  z-index: 1000;
  animation: float-pulse 2s ease-in-out infinite;
}

.whatsapp-float-btn:hover {
  transform: scale(1.1) translateY(-4px);
  box-shadow: 0 8px 30px rgba(37, 211, 102, 0.6);
}

.whatsapp-float-btn:active {
  transform: scale(0.95);
  box-shadow: 0 4px 15px rgba(37, 211, 102, 0.4);
}

.whatsapp-icon {
  width: 32px;
  height: 32px;
  color: white;
  transition: transform 0.3s ease;
}

.whatsapp-float-btn:hover .whatsapp-icon {
  transform: scale(1.1);
}

/* أنيميشن النبض */
@keyframes float-pulse {
  0%, 100% {
    transform: translateY(0);
    box-shadow: 0 4px 20px rgba(37, 211, 102, 0.4);
  }
  50% {
    transform: translateY(-8px);
    box-shadow: 0 8px 30px rgba(37, 211, 102, 0.6);
  }
}

/* تأثير الرنين */
.whatsapp-float-btn::before {
  content: '';
  position: absolute;
  inset: -4px;
  border-radius: 50%;
  border: 2px solid #25D366;
  opacity: 0;
  animation: ripple 2s ease-out infinite;
}

@keyframes ripple {
  0% {
    transform: scale(1);
    opacity: 0.6;
  }
  100% {
    transform: scale(1.5);
    opacity: 0;
  }
}

/* استجابة للشاشات الصغيرة */
@media (max-width: 480px) {
  .whatsapp-float-btn {
    width: 50px;
    height: 50px;
    bottom: 20px;
    right: 20px;
  }

  .whatsapp-icon {
    width: 26px;
    height: 26px;
  }
}
</style>
