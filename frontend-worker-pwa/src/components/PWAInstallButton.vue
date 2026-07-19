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
        <button class="close-btn" @click="showIOSInstructions = false">✕</button>
        <div class="ios-header">
          <div class="ios-icon">📱</div>
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
                    <div class="share-btn">⎋</div>
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
                <div class="menu-item">➕ Add to Home Screen</div>
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
  bottom: 16px;
  right: 16px;
  z-index: 1000;
}

.pwa-install-icon {
  width: 44px;
  height: 44px;
  padding: 0;
  background: linear-gradient(135deg, #1F6F5C 0%, #2d8a6f 100%);
  color: white;
  border: none;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  box-shadow: 0 4px 12px rgba(31, 111, 92, 0.3);
  transition: all 0.3s ease;
}

.pwa-install-icon:hover {
  transform: scale(1.1);
  box-shadow: 0 6px 16px rgba(31, 111, 92, 0.4);
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
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2000;
  padding: 20px;
}

.ios-instructions-content {
  background: white;
  border-radius: 20px;
  padding: 28px;
  max-width: 420px;
  width: 100%;
  position: relative;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.4);
  max-height: 90vh;
  overflow-y: auto;
}

.close-btn {
  position: absolute;
  top: 16px;
  left: 16px;
  background: none;
  border: none;
  font-size: 24px;
  cursor: pointer;
  color: #666;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  transition: background 0.2s;
}

.close-btn:hover {
  background: #f0f0f0;
}

.ios-header {
  text-align: center;
  margin-bottom: 24px;
  padding-top: 8px;
}

.ios-icon {
  font-size: 48px;
  margin-bottom: 12px;
}

.ios-header h3 {
  margin: 0 0 8px 0;
  font-size: 20px;
  color: #1F6F5C;
  font-weight: 700;
}

.ios-subtitle {
  margin: 0;
  font-size: 14px;
  color: #666;
}

.steps {
  display: flex;
  flex-direction: column;
  gap: 20px;
  margin-bottom: 24px;
}

.step {
  display: flex;
  align-items: flex-start;
  gap: 14px;
}

.step-number {
  width: 32px;
  height: 32px;
  background: linear-gradient(135deg, #1F6F5C 0%, #2d8a6f 100%);
  color: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  font-size: 15px;
  flex-shrink: 0;
}

.step-content {
  flex: 1;
}

.step p {
  margin: 0 0 8px 0;
  font-size: 14px;
  color: #333;
  line-height: 1.5;
}

.step strong {
  color: #1F6F5C;
}

.step-visual {
  background: #f8f9fa;
  border-radius: 8px;
  padding: 8px;
  display: flex;
  justify-content: center;
  border: 1px solid #e9ecef;
}

.phone-frame {
  width: 60px;
  height: 100px;
  background: white;
  border: 2px solid #dee2e6;
  border-radius: 8px;
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
}

.phone-screen {
  width: 90%;
  height: 90%;
  background: #f8f9fa;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.share-btn {
  font-size: 20px;
  color: #007AFF;
}

.menu-item {
  background: white;
  padding: 8px 12px;
  border-radius: 6px;
  font-size: 12px;
  color: #333;
  border: 1px solid #dee2e6;
}

.add-btn {
  background: #007AFF;
  color: white;
  padding: 6px 16px;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 600;
}

.benefits {
  display: flex;
  justify-content: space-around;
  margin-bottom: 24px;
  padding: 16px;
  background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
  border-radius: 12px;
}

.benefit {
  font-size: 12px;
  color: #495057;
  text-align: center;
  font-weight: 500;
}

.actions {
  display: flex;
  gap: 12px;
}

.got-it-btn {
  flex: 2;
  padding: 14px;
  background: linear-gradient(135deg, #1F6F5C 0%, #2d8a6f 100%);
  color: white;
  border: none;
  border-radius: 12px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: transform 0.2s;
}

.got-it-btn:hover {
  transform: scale(1.02);
}

.remind-btn {
  flex: 1;
  padding: 14px;
  background: #f8f9fa;
  color: #666;
  border: 1px solid #dee2e6;
  border-radius: 12px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: background 0.2s;
}

.remind-btn:hover {
  background: #e9ecef;
}

/* RTL support */
[dir="rtl"] .pwa-install-container {
  right: auto;
  left: 16px;
}

[dir="rtl"] .close-btn {
  left: auto;
  right: 16px;
}

@media (max-width: 768px) {
  .pwa-install-container {
    bottom: 12px;
    right: 12px;
  }
  
  [dir="rtl"] .pwa-install-container {
    right: auto;
    left: 12px;
  }
  
  .pwa-install-icon {
    width: 40px;
    height: 40px;
  }

  .ios-instructions-content {
    padding: 20px;
  }

  .ios-instructions-content h3 {
    font-size: 16px;
  }

  .step p {
    font-size: 13px;
  }
}
</style>