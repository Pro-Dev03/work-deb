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
          <img src="/src/assets/company-logo.jpg" alt="WorkTrack logo" class="brand-mark" />
          <h1 class="title">{{ $t('app_name') }}</h1>
        </div>
        <p class="subtitle">{{ $t('login') }}</p>
      </div>

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
          <label>{{ $t('email') }}</label>
          <input v-model="email" type="email" :placeholder="$t('email_placeholder')" required autocomplete="username" readonly @focus="removeReadonly" ref="emailInput" />
        </div>
         
        <div class="field">
          <label>{{ $t('password') }}</label>
          <input v-model="password" type="password" :placeholder="$t('password_placeholder')" required autocomplete="current-password" />
        </div>

        <div v-if="error" class="error">{{ error }}</div>
        <div v-if="debugInfo" class="debug">{{ debugInfo }}</div>

        <button class="btn-login" type="submit" :disabled="loading">
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
      
      if (status === 401) {
        debugInfo.value += '\n\n' + t('login_check')
        debugInfo.value += '\n1. ' + t('login_check_user_exists')
        debugInfo.value += '\n2. ' + t('login_check_password_correct')
        debugInfo.value += '\n3. ' + t('login_check_account_active')
      }
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
  background-color: #f0f8f0;
  background-image: url('../assets/company-logo.jpg');
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
  background-attachment: fixed;
  padding: 20px;
  margin: 0;
  min-height: 100vh;
  width: 100%;
}

.login-page::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(30, 58, 95, 0.85);
  z-index: 0;
}

.login-page > * {
  position: relative;
  z-index: 1;
}

.login-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(15px);
  border-radius: 24px;
  padding: 40px 44px;
  max-width: 420px;
  width: 100%;
  box-shadow: 0 24px 80px rgba(0, 0, 0, 0.35);
  animation: fadeIn 0.5s ease;
  border: 1px solid rgba(255, 255, 255, 0.5);
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
  background: linear-gradient(135deg, #1E3A5F08, #1E3A5F12);
  border-radius: 12px;
  border: 1px solid #1E3A5F15;
  margin-bottom: 20px;
}

.powered-logo {
  width: 44px;
  height: 44px;
  border-radius: 10px;
  object-fit: contain;
  background: white;
  padding: 4px;
}

.powered-text {
  text-align: right;
  font-size: 11px;
  color: #6B7A8A;
  line-height: 1.4;
}

.powered-text strong {
  font-size: 14px;
  color: #1E3A5F;
  font-weight: 800;
}

.powered-slogan {
  display: block;
  font-size: 9px;
  color: #8899AA;
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
  color: #1E3A5F;
  margin: 0;
  letter-spacing: -0.5px;
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
  margin-bottom: 28px;
  padding: 6px;
  background: #F0F4FA;
  border-radius: 14px;
}

.lang-btn {
  padding: 10px 20px;
  border: none;
  border-radius: 10px;
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
  box-shadow: 0 4px 16px rgba(30, 58, 95, 0.3);
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 18px;
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
  border-radius: 12px;
  font-size: 15px;
  transition: all 0.3s ease;
  font-family: inherit;
  background: #F8FAFC;
  width: 100%;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #1A2A3A;
}

/* Hide Chrome autofill ghost text */
.field input:-webkit-autofill,
.field input:-webkit-autofill:hover,
.field input:-webkit-autofill:focus,
.field input:-webkit-autofill:active {
  -webkit-box-shadow: 0 0 0 30px #F8FAFC inset !important;
  -webkit-text-fill-color: #1A2A3A !important;
  transition: background-color 5000s ease-in-out 0s;
}

/* Prevent autofill suggestion ghost text */
.field input[readonly] {
  background: #F8FAFC;
  cursor: text;
}

.field input:focus {
  outline: none;
  border-color: #1E3A5F;
  background: white;
  box-shadow: 0 0 0 4px rgba(30, 58, 95, 0.1);
}

.btn-login {
  padding: 16px;
  background: #1E3A5F;
  color: white;
  border: none;
  border-radius: 12px;
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
  box-shadow: 0 8px 28px rgba(30, 58, 95, 0.3);
}

.btn-login:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none !important;
}

.error {
  color: #C53030;
  font-size: 14px;
  text-align: center;
  background: #FDE8E8;
  padding: 12px;
  border-radius: 10px;
}

.debug {
  color: #2B6CB0;
  font-size: 13px;
  text-align: center;
  background: #EBF4FF;
  padding: 10px;
  border-radius: 10px;
  white-space: pre-line;
}

.footer {
  text-align: center;
  margin-top: 24px;
  padding-top: 20px;
  border-top: 1px solid #E2E8F0;
}

.footer p {
  font-size: 12px;
  color: #8899AA;
  margin: 4px 0;
}

.footer-powered {
  font-size: 11px;
  color: #AABBCC;
}

@media (max-width: 480px) {
  .login-card {
    padding: 28px 20px;
  }
  
  .title {
    font-size: 28px;
  }
  
  .lang-btn {
    padding: 8px 12px;
    font-size: 13px;
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
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
  }
}
</style> 