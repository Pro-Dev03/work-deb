<template>
  <PullToRefresh 
    @refresh="handleRefresh"
    :refresh-text="t('pull_to_refresh')"
    :refreshing-text="t('refreshing')"
    :release-text="t('release_to_refresh')"
  >
    <div class="login-page">
      <div class="login-card">
      <div class="login-header">
        <div class="logo">
          <img src="/src/assets/company-logo.jpg" alt="WorkTrack logo" class="brand-mark" />
        </div>
        <h1 class="title">{{ t('app_name') }} אבן יסודות</h1>
        <p class="subtitle">{{ t('app_name') }} - Employee platform -</p>
        <p class="subtitle-small">{{ t('login') }}</p>
      </div>

      <!-- ✅ محدد اللغة - موحد -->
      <div class="lang-section">
        <button 
          v-for="lang in languages" 
          :key="lang.code"
          class="lang-btn"
          :class="{ active: currentLang === lang.code }"
          @click="changeLanguage(lang.code)"
        >
          {{ lang.flag }} {{ lang.name }}
        </button>
      </div>

      <form class="login-form" @submit.prevent="handleSubmit">
        <div class="field">
          <label>📱 {{ t('phone') }}</label>
          <input 
            v-model="phone" 
            type="tel" 
            placeholder="05xxxxxxxx"
            required 
            dir="ltr"
          />
          <span class="field-hint">{{ t('phone_hint') }}</span>
        </div>

        <div v-if="error" class="error">{{ error }}</div>
        <div v-if="debugInfo" class="debug">{{ debugInfo }}</div>

        <button class="btn-login" type="submit" :disabled="loading">
          {{ loading ? '⏳' : '📱' }} {{ t('login') }}
        </button>
      </form>

      <div class="footer">
        <p>{{ t('created_by_admin') }}</p>
        <div class="devpro-logo">
          <img src="/src/assets/company-logo.jpg" alt="DevPro Logo" class="devpro-img" />
        </div>
      </div>
      <p class="footer-small">{{ t('device_verify') }}</p>
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

const { t, currentLang, setLang } = useI18n()
const router = useRouter()

const phone = ref('')
const loading = ref(false)
const error = ref('')
const debugInfo = ref('')

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
  debugInfo.value = ''

  try {
    const deviceId = getDeviceId()
    const deviceModel = getDeviceModel()
      
    debugInfo.value = `📱 Device: ${deviceModel}\n🆔 ID: ${deviceId}`

    const { data } = await api.post('/auth/phone-login', {
      phone: phone.value.trim(),
      device_id: deviceId,
      device_model: deviceModel
    })

    debugInfo.value = '✅ Verification successful!'
      
    localStorage.setItem('worktrack_token', data.token)
    localStorage.setItem('worktrack_user', JSON.stringify(data.user))
      
    // ✅ اللغة مخزنة بالفعل في localStorage من i18n
      
    setTimeout(() => {
      router.push('/attendance')
    }, 500)

  } catch (e) {
    console.error('❌ Login failed:', e.response?.data)

    if (e.response?.data?.device_mismatch) {
      error.value = '⚠️ This device is not authorized. Please contact the admin to reset your device.'
      debugInfo.value = '🔒 Device not registered'
    } else if (e.response?.data?.model_mismatch) {
      error.value = '⚠️ Device model mismatch. Please contact the admin.'
      debugInfo.value = '🔒 Device model not registered'
    } else {
      error.value = e.response?.data?.error || '❌ Login failed. Please check the phone number.'
      debugInfo.value = '❌ ' + (e.response?.data?.error || 'Unknown error')
    }
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page {
  position: fixed;
  top: 0; left: 0; right: 0; bottom: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #1E3A5F 0%, #0D1B3E 100%);
  padding: 20px;
  margin: 0;
  min-height: 100vh;
  width: 100%;
}

.login-card {
  background: #FFFFFF;
  border-radius: 24px;
  padding: 40px 36px;
  max-width: 400px;
  width: 100%;
  box-shadow: 0 24px 80px rgba(0, 0, 0, 0.35);
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
  margin-bottom: 12px;
}

.brand-mark {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  object-fit: cover;
  background: transparent;
}

.title {
  font-size: 28px;
  font-weight: 800;
  color: #1E3A5F;
  margin: 0 0 4px 0;
  direction: ltr;
}

.subtitle {
  font-size: 15px;
  color: #6B7A8A;
  margin: 0;
}

.subtitle-small {
  font-size: 13px;
  color: #8899AA;
  margin: 4px 0 0;
}

.lang-section {
  display: flex;
  justify-content: center;
  gap: 8px;
  margin-bottom: 24px;
  padding: 6px;
  background: #F0F4FA;
  border-radius: 12px;
}

.lang-btn {
  padding: 8px 16px;
  border: none;
  border-radius: 8px;
  background: transparent;
  font-size: 14px;
  font-weight: 600;
  color: #6B7A8A;
  cursor: pointer;
  transition: all 0.3s ease;
  font-family: inherit;
  flex: 1;
}

.lang-btn:hover {
  background: rgba(30, 58, 95, 0.08);
  color: #1E3A5F;
}

.lang-btn.active {
  background: #1E3A5F;
  color: white;
  box-shadow: 0 4px 12px rgba(30, 58, 95, 0.3);
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.field {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.field label {
  font-size: 14px;
  font-weight: 600;
  color: #1A2A3A;
}

.field input {
  padding: 14px 16px;
  border: 2px solid #E2E8F0;
  border-radius: 10px;
  font-size: 18px;
  transition: all 0.3s ease;
  font-family: inherit;
  background: #F8FAFC;
  text-align: center;
  letter-spacing: 2px;
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
}

.error {
  color: #C53030;
  font-size: 14px;
  text-align: center;
  background: #FDE8E8;
  padding: 12px;
  border-radius: 8px;
}

.debug {
  color: #2B6CB0;
  font-size: 13px;
  text-align: center;
  background: #EBF4FF;
  padding: 10px;
  border-radius: 8px;
  white-space: pre-line;
}

.btn-login {
  padding: 16px;
  background: #1E3A5F;
  color: white;
  border: none;
  border-radius: 10px;
  font-size: 17px;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.3s ease;
  font-family: inherit;
  margin-top: 4px;
}

.btn-login:hover:not(:disabled) {
  background: #0D1B3E;
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(30, 58, 95, 0.3);
}

.btn-login:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none !important;
}

.footer {
  text-align: center;
  font-size: 13px;
  color: #8899AA;
  margin-top: 20px;
  padding-top: 16px;
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
  border-radius: 8px;
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
  margin-top: 8px;
}

@media (max-width: 480px) {
  .login-card {
    padding: 28px 20px;
  }

  .title {
    font-size: 22px;
  }

  .lang-btn {
    padding: 6px 12px;
    font-size: 13px;
  }

  .devpro-img {
    max-width: 80px;
  }
}
</style>
