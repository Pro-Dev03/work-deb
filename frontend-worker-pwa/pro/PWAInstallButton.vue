<template>
  <div v-if="showInstallButton" class="pwa-install-container">
    <button 
      @click="installPWA" 
      class="pwa-install-icon"
      :title="pwaInstallTitle"
    >
      <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path>
        <polyline points="7 10 12 15 17 10"></polyline>
        <line x1="12" y1="15" x2="12" y2="3"></line>
      </svg>
    </button>
    
    <!-- iOS Instructions Modal -->
    <div v-if="showIOSInstructions" class="ios-instructions-modal" @click.self="showIOSInstructions = false">
      <div class="ios-instructions-content">
        <button class="close-btn" @click="showIOSInstructions = false" aria-label="إغلاق">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
        </button>
        <div class="ios-header">
          <div class="ios-icon">
            <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round"><rect x="7" y="2" width="10" height="20" rx="2"/><line x1="11" y1="18" x2="13" y2="18"/></svg>
          </div>
          <h3>{{ iosModalTitle }}</h3>
          <p class="ios-subtitle">{{ iosModalSubtitle }}</p>
        </div>
        
        <div class="steps">
          <div class="step">
            <span class="step-number">1</span>
            <div class="step-content">
              <p v-html="step1Text"></p>
              <div class="step-visual">
                <div class="phone-frame">
                  <div class="phone-screen">
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M4 12v7a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2v-7"/><polyline points="16 6 12 2 8 6"/><line x1="12" y1="2" x2="12" y2="15"/></svg>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div class="step">
            <span class="step-number">2</span>
            <div class="step-content">
              <p v-html="step2Text"></p>
              <div class="step-visual">
                <div class="menu-item">
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
                  Add to Home Screen
                </div>
              </div>
            </div>
          </div>
          <div class="step">
            <span class="step-number">3</span>
            <div class="step-content">
              <p v-html="step3Text"></p>
              <div class="step-visual">
                <div class="add-btn">Add</div>
              </div>
            </div>
          </div>
        </div>
        
        <div class="benefits">
          <div class="benefit">{{ benefit1 }}</div>
          <div class="benefit">{{ benefit2 }}</div>
          <div class="benefit">{{ benefit3 }}</div>
        </div>
        
        <div class="actions">
          <button class="got-it-btn" @click="showIOSInstructions = false">{{ gotItText }}</button>
          <button class="remind-btn" @click="remindLater">{{ remindLaterText }}</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useI18n } from '../services/i18n'

const { t, currentLang } = useI18n()

const showInstallButton = ref(false)
const showIOSInstructions = ref(false)
const isIOS = ref(false)
const isStandalone = ref(false)

const pwaInstallTitle = computed(() => t('pwa.installTitle'))

const pwaInstallText = computed(() => {
  if (isIOS.value) {
    return t('pwa.installOnIos')
  }
  return t('pwa.installText')
})

const iosModalTitle = computed(() => t('pwa.iosModalTitle'))
const iosModalSubtitle = computed(() => t('pwa.iosModalSubtitle'))
const step1Text = computed(() => t('pwa.step1Text'))
const step2Text = computed(() => t('pwa.step2Text'))
const step3Text = computed(() => t('pwa.step3Text'))
const benefit1 = computed(() => t('pwa.benefit1'))
const benefit2 = computed(() => t('pwa.benefit2'))
const benefit3 = computed(() => t('pwa.benefit3'))
const gotItText = computed(() => t('pwa.gotIt'))
const remindLaterText = computed(() => t('pwa.remindLater'))

function detectIOS() {
  const userAgent = window.navigator.userAgent.toLowerCase()
  return /iphone|ipad|ipod/.test(userAgent) && !/edge|crios|fxios/.test(userAgent)
}

onMounted(() => {
  // Check if device is iOS
  isIOS.value = detectIOS()
  
  // Check if already installed (standalone mode)
  isStandalone.value = window.matchMedia('(display-mode: standalone)').matches || 
                     window.navigator.standalone === true

  // Don't show if already installed
  if (isStandalone.value) {
    return
  }

  // Always show button on login page
  if (isIOS.value) {
    showInstallButton.value = true
  } else {
    // For other devices, listen to install prompt
    window.addEventListener('pwa-install-available', handleInstallAvailable)
    window.addEventListener('pwa-install-success', handleInstallSuccess)
    
    // Check for existing deferredPrompt
    if (window.deferredPrompt) {
      showInstallButton.value = true
    }
  }
})

onUnmounted(() => {
  window.removeEventListener('pwa-install-available', handleInstallAvailable)
  window.removeEventListener('pwa-install-success', handleInstallSuccess)
})

function handleInstallAvailable() {
  showInstallButton.value = true
}

function handleInstallSuccess() {
  showInstallButton.value = false
}

function installPWA() {
  if (isIOS.value) {
    // Show iOS instructions
    showIOSInstructions.value = true
  } else if (window.pwaInstall) {
    // Use native install prompt for Android/Windows
    window.pwaInstall()
  } else {
    // Fallback - show general instructions
    showIOSInstructions.value = true
  }
}

function remindLater() {
  // Just close the modal, keep icon visible
  showIOSInstructions.value = false
}
</script>

<style scoped>
.pwa-install-container {
  position: fixed;
  bottom: var(--space-4);
  right: var(--space-4);
  z-index: 1000;
}

.pwa-install-icon {
  width: 44px;
  height: 44px;
  padding: 0;
  background: linear-gradient(135deg, var(--primary-500) 0%, var(--primary-600) 100%);
  color: var(--text-inverse);
  border: none;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  box-shadow: var(--shadow-lg);
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.pwa-install-icon:hover {
  transform: scale(1.08);
  box-shadow: var(--shadow-xl);
}

.pwa-install-icon:active {
  transform: scale(0.95);
}

/* iOS Instructions Modal */
.ios-instructions-modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(28, 25, 23, 0.6);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2000;
  padding: var(--space-5);
}

.ios-instructions-content {
  background: var(--surface);
  border-radius: var(--radius-2xl);
  padding: var(--space-8);
  max-width: 420px;
  width: 100%;
  position: relative;
  box-shadow: var(--shadow-xl);
  max-height: 90vh;
  overflow-y: auto;
}

.close-btn {
  position: absolute;
  top: var(--space-4);
  left: var(--space-4);
  background: var(--surface-elevated);
  border: none;
  cursor: pointer;
  color: var(--text-secondary);
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  transition: background 0.2s ease, color 0.2s ease;
}

.close-btn:hover {
  background: var(--gray-200);
  color: var(--text-primary);
}

.ios-header {
  text-align: center;
  margin-bottom: var(--space-6);
  padding-top: var(--space-2);
}

.ios-icon {
  width: 56px; height: 56px; margin: 0 auto var(--space-3);
  border-radius: var(--radius-lg);
  background: var(--primary-100); color: var(--primary-700);
  display: flex; align-items: center; justify-content: center;
}

.ios-header h3 {
  margin: 0 0 var(--space-2) 0;
  font-size: var(--text-xl);
  color: var(--text-primary);
  font-weight: var(--font-bold);
}

.ios-subtitle {
  margin: 0;
  font-size: var(--text-sm);
  color: var(--text-secondary);
}

.steps {
  display: flex;
  flex-direction: column;
  gap: var(--space-5);
  margin-bottom: var(--space-6);
}

.step {
  display: flex;
  align-items: flex-start;
  gap: var(--space-3);
}

.step-number {
  width: 32px;
  height: 32px;
  background: linear-gradient(135deg, var(--primary-500) 0%, var(--primary-600) 100%);
  color: var(--text-inverse);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: var(--font-bold);
  font-size: var(--text-sm);
  flex-shrink: 0;
}

.step-content {
  flex: 1;
}

.step p {
  margin: 0 0 var(--space-2) 0;
  font-size: var(--text-sm);
  color: var(--text-primary);
  line-height: var(--leading-normal);
}

.step strong {
  color: var(--primary-600);
}

.step-visual {
  background: var(--surface-elevated);
  border-radius: var(--radius-sm);
  padding: var(--space-2);
  display: flex;
  justify-content: center;
  border: 1px solid var(--border);
  color: var(--text-secondary);
}

.phone-frame {
  width: 60px;
  height: 100px;
  background: var(--surface);
  border: 2px solid var(--border-strong);
  border-radius: var(--radius-sm);
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
}

.phone-screen {
  width: 90%;
  height: 90%;
  background: var(--surface-elevated);
  border-radius: var(--radius-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--primary-500);
}

.menu-item {
  display: flex; align-items: center; gap: 6px;
  background: var(--surface);
  padding: var(--space-2) var(--space-3);
  border-radius: var(--radius-sm);
  font-size: var(--text-xs);
  color: var(--text-primary);
  border: 1px solid var(--border-strong);
}

.add-btn {
  background: var(--primary-500);
  color: var(--text-inverse);
  padding: 6px var(--space-4);
  border-radius: var(--radius-sm);
  font-size: var(--text-xs);
  font-weight: var(--font-semibold);
}

.benefits {
  display: flex;
  justify-content: space-around;
  margin-bottom: var(--space-6);
  padding: var(--space-4);
  background: var(--surface-elevated);
  border-radius: var(--radius-lg);
}

.benefit {
  font-size: var(--text-xs);
  color: var(--text-secondary);
  text-align: center;
  font-weight: var(--font-medium);
}

.actions {
  display: flex;
  gap: var(--space-3);
}

.got-it-btn {
  flex: 2;
  padding: 14px;
  background: linear-gradient(135deg, var(--primary-500) 0%, var(--primary-600) 100%);
  color: var(--text-inverse);
  border: none;
  border-radius: var(--radius-lg);
  font-size: var(--text-base);
  font-weight: var(--font-semibold);
  cursor: pointer;
  box-shadow: var(--shadow-md);
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.got-it-btn:hover {
  box-shadow: var(--shadow-lg);
  transform: translateY(-1px);
}

.remind-btn {
  flex: 1;
  padding: 14px;
  background: var(--surface);
  color: var(--text-secondary);
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
  font-size: var(--text-sm);
  font-weight: var(--font-medium);
  cursor: pointer;
  transition: background 0.2s ease, border-color 0.2s ease;
}

.remind-btn:hover {
  background: var(--surface-elevated);
  border-color: var(--border-strong);
}

/* RTL support */
[dir="rtl"] .pwa-install-container {
  right: auto;
  left: var(--space-4);
}

[dir="rtl"] .close-btn {
  left: auto;
  right: var(--space-4);
}

@media (max-width: 768px) {
  .pwa-install-container {
    bottom: var(--space-3);
    right: var(--space-3);
  }
  
  [dir="rtl"] .pwa-install-container {
    right: auto;
    left: var(--space-3);
  }
  
  .pwa-install-icon {
    width: 40px;
    height: 40px;
  }

  .ios-instructions-content {
    padding: var(--space-5);
  }

  .ios-instructions-content h3 {
    font-size: var(--text-lg);
  }

  .step p {
    font-size: 13px;
  }
}
</style>
