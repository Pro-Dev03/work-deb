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
        <div class="company-logo-container">
          <img src="/src/assets/company-logo.jpg" alt="WorkTrack logo" class="company-logo" />
        </div>

        <p class="subtitle">{{ $t('login') }}</p>
      </div>

      <!-- Theme & Language Toolbar -->
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

        <button type="submit" class="btn btn--primary btn--block" :disabled="loading">
          <svg v-if="loading" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="animate-spin">
            <path d="M21 12a9 9 0 1 1-6.219-8.56"/>
          </svg>
          <svg v-else width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M15 3h4a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2h-4"/>
            <polyline points="10 17 15 12 10 7"/>
            <line x1="15" y1="12" x2="3" y2="12"/>
          </svg>
          {{ $t('login') }}
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
import { ref, onMounted } from 'vue'
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

function removeReadonly() {
  if (emailInput.value) {
    emailInput.value.removeAttribute('readonly')
  }
}

function changeLanguage(code) {
  setLang(code)
  
  // Set direction based on language
  const lang = languages.find(l => l.code === code)
  if (lang) {
    document.documentElement.setAttribute('dir', lang.dir)
    document.documentElement.setAttribute('data-lang', lang.code)
  }
}

function toggleTheme() {
  isDarkMode.value = !isDarkMode.value
  const theme = isDarkMode.value ? 'dark' : 'light'
  document.documentElement.setAttribute('data-theme', theme)
  localStorage.setItem('worktrack_theme', theme)
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
/* Settings Toolbar ( Inside Card) - SaaS Style */
.settings-toolbar {
  display: flex;
  gap: 8px;
  margin-bottom: 28px;
  justify-content: center;
  flex-wrap: wrap;
}

.settings-btn {
  border: none;
  cursor: pointer;
  border-radius: var(--radius-sm);
  width: 40px;
  height: 40px;
  min-width: 40px;
  min-height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--ink-soft);
  background: var(--canvas);
  transition: all var(--transition-base);
  font-size: 13px;
  font-weight: var(--font-medium);
  padding: 0;
  border: 1px solid var(--line);
}

.settings-btn--lang {
  min-width: 36px;
  font-size: 11px;
  width: 36px;
  height: 36px;
  border-radius: var(--radius-sm);
}

.settings-btn:hover {
  background: var(--line);
  border-color: var(--line-strong);
  color: var(--ink);
  transform: translateY(-1px);
}

.settings-btn.active {
  background: var(--brand);
  color: white;
  border-color: var(--brand);
  box-shadow: var(--shadow-md), var(--brand-glow);
}

.settings-btn.active:hover {
  background: var(--brand-dark);
  border-color: var(--brand-dark);
  transform: translateY(-1px);
  box-shadow: var(--shadow-lg), var(--shadow-glow);
}

.login-page {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: var(--canvas);
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
  background: rgba(30, 58, 95, 0.4);
  z-index: 0;
}

[data-theme="dark"] .login-page::before {
  background: rgba(0, 0, 0, 0.7);
}

.login-page > * {
  position: relative;
  z-index: 1;
}

.login-page .pwa-install-container {
  position: fixed !important;
  z-index: 9999 !important;
}

.login-card {
  background: rgba(255, 255, 255, 0.6);
  backdrop-filter: blur(25px);
  border-radius: var(--radius-xl);
  padding: 40px 44px;
  max-width: 420px;
  width: 100%;
  box-shadow: var(--shadow-xl);
  animation: fadeIn 0.5s ease;
  border: 1px solid rgba(255, 255, 255, 0.3);
}

[data-theme="dark"] .login-card {
  background: rgba(30, 41, 59, 0.8);
  border: 1px solid rgba(255, 255, 255, 0.1);
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: translateY(0); }
}

.login-header {
  text-align: center;
  margin-bottom: 28px;
}

.company-logo-container {
  display: flex;
  justify-content: center;
  margin-bottom: 24px;
}

.company-logo {
  width: 120px;
  height: 120px;
  object-fit: contain;
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.3);
  padding: 12px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
}

.subtitle {
  font-size: 15px;
  color: var(--ink-soft);
  margin: 0;
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
  color: var(--ink);
}

.field input {
  padding: 14px 16px;
  border: 2px solid var(--line);
  border-radius: var(--radius-md);
  font-size: 15px;
  transition: all var(--transition-base);
  font-family: inherit;
  background: rgba(255, 255, 255, 0.7);
  backdrop-filter: blur(10px);
  width: 100%;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: var(--ink);
}

[data-theme="dark"] .field input {
  background: rgba(30, 41, 59, 0.7);
}

/* Hide Chrome autofill ghost text */
.field input:-webkit-autofill,
.field input:-webkit-autofill:hover,
.field input:-webkit-autofill:focus,
.field input:-webkit-autofill:active {
  -webkit-box-shadow: 0 0 0 30px var(--canvas) inset !important;
  -webkit-text-fill-color: var(--ink) !important;
  transition: background-color 5000s ease-in-out 0s;
}

[data-theme="dark"] .field input:-webkit-autofill,
[data-theme="dark"] .field input:-webkit-autofill:hover,
[data-theme="dark"] .field input:-webkit-autofill:focus,
[data-theme="dark"] .field input:-webkit-autofill:active {
  -webkit-box-shadow: 0 0 0 30px var(--surface) inset !important;
  -webkit-text-fill-color: var(--ink) !important;
}

/* Prevent autofill suggestion ghost text */
.field input[readonly] {
  background: var(--canvas);
  cursor: text;
}

[data-theme="dark"] .field input[readonly] {
  background: var(--surface);
}

.field input:focus {
  outline: none;
  border-color: var(--brand);
  background: var(--surface);
  box-shadow: 0 0 0 4px var(--brand-tint);
}

[data-theme="dark"] .field input:focus {
  background: var(--surface);
}

/* Button Styles */
.btn {
  border: none;
  border-radius: var(--radius-md);
  padding: 12px 24px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all var(--transition-base);
  font-family: inherit;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.btn--primary {
  background: var(--brand);
  backdrop-filter: blur(10px);
  color: white;
  box-shadow: var(--shadow-md), var(--brand-glow);
}

.btn--primary:hover:not(:disabled) {
  background: var(--brand-dark);
  transform: translateY(-2px);
  box-shadow: var(--shadow-lg), var(--shadow-glow);
}

.btn--primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none !important;
}

.btn--block {
  width: 100%;
}

.animate-spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
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
  border-top: 1px solid var(--line);
}

.footer p {
  font-size: 12px;
  color: var(--ink-light);
  margin: 4px 0;
}

.footer-powered {
  font-size: 11px;
  color: var(--ink-light);
}

@media (max-width: 480px) {
  .login-card {
    padding: 24px 18px;
    max-width: 340px;
  }
  
  .company-logo {
    width: 100px;
    height: 100px;
    padding: 10px;
  }
  
  .company-logo-container {
    margin-bottom: 20px;
  }
  
  .title {
    font-size: 28px;
  }
  
  .settings-btn {
    width: 28px;
    height: 28px;
  }
  
  .field input {
    font-size: 16px;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
  }
}
</style>