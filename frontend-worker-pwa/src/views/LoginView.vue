<template>
  <PullToRefresh 
    @refresh="handleRefresh"
    :refresh-text="t('pull_to_refresh')"
    :refreshing-text="t('refreshing')"
    :release-text="t('release_to_refresh')"
  >
    <div class="login-page">
      <div class="login-card">
        <PWAInstallButton />
      <div class="login-header">
        <div class="logo-mark">
          <img src="/src/assets/company-logo.jpg" alt="WorkTrack logo" class="brand-mark" />
        </div>
        <h1>WorkTrack</h1>
        <p>بوابة الموظف</p>
      </div>

      <!-- Theme & Language Toolbar (Above Phone Input) -->
      <div class="settings-toolbar" role="group" aria-label="إعدادات العرض">
        <button 
          class="settings-btn" 
          :class="{ active: isDarkMode }"
          @click="toggleTheme"
          title="الوضع الليلي" 
          aria-label="تبديل الوضع الليلي"
        >
          <svg v-if="!isDarkMode" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79Z"/>
          </svg>
          <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="5"/>
            <path d="M12 1v2M12 21v2M4.22 4.22l1.42 1.42M18.36 18.36l1.42 1.42M1 12h2M21 12h2M4.22 19.78l1.42-1.42M18.36 5.64l1.42-1.42"/>
          </svg>
        </button>
        
        <!-- Language Buttons -->
        <button 
          v-for="lang in languages" 
          :key="lang.code"
          class="settings-btn settings-btn--lang"
          :class="{ active: currentLang === lang.code }"
          @click="changeLanguage(lang.code)"
          :title="lang.name"
          :aria-label="lang.name"
        >
          {{ lang.label }}
        </button>
      </div>

      <form class="login-form" @submit.prevent="handleSubmit">
        <div class="form-group">
          <label class="form-label">
            <span class="icon">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72c.127.96.362 1.903.7 2.81a2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45c.907.338 1.85.573 2.81.7A2 2 0 0 1 22 16.92Z"/>
              </svg>
            </span>
            رقم الهاتف
          </label>
          <input 
            v-model="phone" 
            type="tel" 
            class="form-input" 
            placeholder="05xxxxxxxx" 
            required 
            dir="ltr"
          />
        </div>
        
        <div v-if="error" class="error-message">{{ error }}</div>
        
        <button type="submit" class="btn btn--primary btn--block" :disabled="loading">
          <svg v-if="loading" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="animate-spin">
            <path d="M21 12a9 9 0 1 1-6.219-8.56"/>
          </svg>
          <svg v-else width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M15 3h4a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2h-4"/>
            <polyline points="10 17 15 12 10 7"/>
            <line x1="15" y1="12" x2="3" y2="12"/>
          </svg>
          {{ t('login') }}
        </button>
      </form>
      
      <p class="form-hint">بتسجيل الدخول أنت توافق على شروط الاستخدام</p>
    </div>
    </div>
  </PullToRefresh>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from '../services/i18n'
import api from '../services/api'
import PWAInstallButton from '../components/PWAInstallButton.vue'
import PullToRefresh from '../components/PullToRefresh.vue'

const { t, currentLang, setLang } = useI18n()
const router = useRouter()

const phone = ref('')
const loading = ref(false)
const error = ref('')
const isDarkMode = ref(false)

const languages = [
  { code: 'ar', name: 'العربية', dir: 'rtl', label: 'AR' },
  { code: 'he', name: 'עברית', dir: 'rtl', label: 'HE' },
  { code: 'en', name: 'English', dir: 'ltr', label: 'EN' }
]

onMounted(() => {
  // Load saved theme preference
  const savedTheme = localStorage.getItem('worktrack_theme')
  if (savedTheme === 'dark') {
    isDarkMode.value = true
    document.documentElement.setAttribute('data-theme', 'dark')
  }
})

async function handleRefresh() {
  // إرسال رسالة للـ Service Worker لتنظيف الكاش
  if ('serviceWorker' in navigator) {
    const registration = await navigator.serviceWorker.getRegistration()
    if (registration) {
      registration.active.postMessage({ type: 'CACHE_BUST' })
    }
  }
  
  // إعادة تحميل الصفحة
  window.location.reload()
}

function changeLanguage(code) {
  setLang(code)
  
  // Set direction based on language
  const lang = languages.find(l => l.code === code)
  if (lang) {
    document.documentElement.setAttribute('dir', lang.dir)
    document.documentElement.setAttribute('data-lang', lang.code)
  }
  
  // إعادة تحميل الصفحة بعد تغيير اللغة
  // لضمان تطبيق التغيير على جميع المكونات
}

function toggleTheme() {
  isDarkMode.value = !isDarkMode.value
  const theme = isDarkMode.value ? 'dark' : 'light'
  document.documentElement.setAttribute('data-theme', theme)
  localStorage.setItem('worktrack_theme', theme)
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
    error.value = 'الرجاء إدخال رقم هاتف صحيح'
    return
  }

  loading.value = true
  error.value = ''

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

    setTimeout(() => {
      router.push('/attendance')
    }, 500)

  } catch (e) {
    console.error('❌ Login failed:', e.response?.data)

    if (e.response?.data?.device_mismatch) {
      error.value = '⚠️ Login failed. Please contact the admin.'
    } else if (e.response?.data?.model_mismatch) {
      error.value = '⚠️ Login failed. Please contact the admin.'
    } else {
      error.value = e.response?.data?.error || '❌ Login failed. Please check the phone number.'
    }
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
/* Settings Toolbar (Inside Card) - SaaS Style */
.settings-toolbar {
  display: flex;
  gap: 6px;
  margin-bottom: var(--space-5);
  justify-content: center;
  flex-wrap: wrap;
}

.settings-btn {
  border: none;
  background: transparent;
  cursor: pointer;
  border-radius: var(--radius-md);
  width: 32px;
  height: 32px;
  min-width: 32px;
  min-height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-secondary);
  transition: var(--transition-base) ease;
  font-size: var(--text-xs);
  font-weight: var(--font-medium);
  padding: 0;
}

.settings-btn--lang {
  min-width: 28px;
  font-size: 9px;
}

.settings-btn:hover {
  background: var(--gray-100);
  color: var(--text-primary);
}

.settings-btn.active {
  background: var(--primary-500);
  color: #fff;
}

.login-page {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  padding: var(--space-6);
  background: var(--background);
  position: relative;
}

.login-card {
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
  padding: var(--space-6);
  box-shadow: var(--shadow-sm);
  max-width: 400px;
  width: 100%;
  margin: 0 auto;
  position: relative;
}

.login-header {
  text-align: center;
  margin-bottom: var(--space-5);
}

.logo-mark {
  width: 300px;
  height: 300px;
  margin: 0 auto var(--space-3);
  display: flex;
  align-items: center;
  justify-content: center;
}

.brand-mark {
  width: 300px;
  height: 300px;
  object-fit: cover;
  background: transparent;
}

.login-header h1 {
  font-size: var(--text-xl);
  font-weight: var(--font-semibold);
  letter-spacing: -0.01em;
  color: var(--text-primary);
  margin: 0;
}

.login-header p {
  font-size: var(--text-sm);
  color: var(--text-secondary);
  margin-top: 2px;
}

.animate-spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

/* Settings Toolbar (Inside Card) - SaaS Style */
.settings-toolbar {
  display: flex;
  gap: 6px;
  margin-bottom: var(--space-5);
  justify-content: center;
  flex-wrap: wrap;
}

.settings-btn {
  border: none;
  background: transparent;
  cursor: pointer;
  border-radius: var(--radius-md);
  width: 32px;
  height: 32px;
  min-width: 32px;
  min-height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-secondary);
  transition: var(--transition-base) ease;
  font-size: var(--text-xs);
  font-weight: var(--font-medium);
  padding: 0;
}

.settings-btn--lang {
  min-width: 28px;
  font-size: 9px;
}

.settings-btn:hover {
  background: var(--gray-100);
  color: var(--text-primary);
}

.settings-btn.active {
  background: var(--primary-500);
  color: #fff;
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.form-label {
  font-size: var(--text-sm);
  font-weight: var(--font-medium);
  color: var(--text-primary);
  display: flex;
  align-items: center;
  gap: var(--space-2);
}

.form-input {
  width: 100%;
  padding: 12px 16px;
  border: 2px solid var(--border);
  border-radius: var(--radius-md);
  font-size: 16px;
  background: var(--surface);
  color: var(--text-primary);
  transition: all 0.2s ease;
  font-family: inherit;
}

.form-input:focus {
  outline: none;
  border-color: var(--primary-500);
  box-shadow: 0 0 0 3px var(--primary-50);
}

.error-message {
  color: var(--error-700);
  font-size: var(--text-sm);
  text-align: center;
  background: var(--error-50);
  padding: var(--space-2);
  border-radius: var(--radius-md);
  border: 1px solid var(--error-100);
}

.form-hint {
  font-size: var(--text-xs);
  color: var(--text-tertiary);
  text-align: center;
  margin-top: var(--space-2);
}

/* Responsive Design */
@media (max-width: 480px) {
  .login-page {
    padding: var(--space-4);
  }

  .login-card {
    padding: var(--space-5);
  }

  .login-header h1 {
    font-size: var(--text-lg);
  }
}

@media (min-width: 768px) {
  .login-card {
    padding: var(--space-8);
  }
}
</style>
