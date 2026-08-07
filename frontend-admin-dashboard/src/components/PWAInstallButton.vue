<template>
  <div v-if="showInstallButton" class="pwa-install-container">
    <button 
      @click="installPWA" 
      class="pwa-install-icon"
      :title="pwaInstallTitle"
    >
      <Download :size="20" />
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

<script>
import { Download } from '@lucide/vue'

export default {
  name: 'PWAInstallButton',
  components: {
    Download
  },
  data() {
    return {
      showInstallButton: false,
      showIOSInstructions: false,
      deferredPrompt: null,
      isIOS: false,
      isStandalone: false
    }
  },
  computed: {
    currentLang() {
      if (this.$i18n && this.$i18n.currentLang) {
        return this.$i18n.currentLang.value
      }
      return 'ar' // default
    },
    pwaInstallTitle() {
      return this.$t('pwa.installTitle')
    },
    pwaInstallText() {
      if (this.isIOS) {
        return this.$t('pwa.installOnIos')
      }
      return this.$t('pwa.installText')
    },
    iosModalTitle() {
      return this.$t('pwa.iosModalTitle')
    },
    iosModalSubtitle() {
      return this.$t('pwa.iosModalSubtitle')
    },
    step1Text() {
      return this.$t('pwa.step1Text')
    },
    step2Text() {
      return this.$t('pwa.step2Text')
    },
    step3Text() {
      return this.$t('pwa.step3Text')
    },
    benefit1() {
      return this.$t('pwa.benefit1')
    },
    benefit2() {
      return this.$t('pwa.benefit2')
    },
    benefit3() {
      return this.$t('pwa.benefit3')
    },
    gotItText() {
      return this.$t('pwa.gotIt')
    },
    remindLaterText() {
      return this.$t('pwa.remindLater')
    }
  },
  methods: {
    detectIOS() {
      const userAgent = window.navigator.userAgent.toLowerCase()
      return /iphone|ipad|ipod/.test(userAgent) && !/edge|crios|fxios/.test(userAgent)
    },
    handleInstallAvailable() {
      this.showInstallButton = true
    },
    handleInstallSuccess() {
      this.showInstallButton = false
    },
    installPWA() {
      if (this.isIOS) {
        // Show iOS instructions
        this.showIOSInstructions = true
      } else if (window.pwaInstall) {
        // Use native install prompt for Android/Windows
        window.pwaInstall()
      } else {
        // Fallback - show general instructions
        this.showIOSInstructions = true
      }
    },
    remindLater() {
      // Just close the modal, keep icon visible
      this.showIOSInstructions = false
    }
  },
  mounted() {
    // Check if device is iOS
    this.isIOS = this.detectIOS()
    
    // Check if already installed (standalone mode)
    this.isStandalone = window.matchMedia('(display-mode: standalone)').matches || 
                      window.navigator.standalone === true

    // Don't show if already installed
    if (this.isStandalone) {
      return
    }

    // Always show button on login page
    if (this.isIOS) {
      this.showInstallButton = true
    } else {
      // For other devices, listen to install prompt
      window.addEventListener('pwa-install-available', this.handleInstallAvailable)
      window.addEventListener('pwa-install-success', this.handleInstallSuccess)
      
      // Check for existing deferredPrompt
      if (window.deferredPrompt) {
        this.showInstallButton = true
      }
    }
  },
  beforeUnmount() {
    window.removeEventListener('pwa-install-available', this.handleInstallAvailable)
    window.removeEventListener('pwa-install-success', this.handleInstallSuccess)
  }
}
</script>

<style scoped>
.pwa-install-container {
  position: fixed;
  bottom: 16px;
  left: 16px;
  z-index: 1000;
}

.pwa-install-icon {
  width: 44px;
  height: 44px;
  padding: 0;
  background: linear-gradient(135deg, #1e3a8a 0%, #3b82f6 100%);
  color: white;
  border: none;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  box-shadow: 0 4px 12px rgba(30, 58, 138, 0.3);
  transition: all 0.3s ease;
}

.pwa-install-icon:hover {
  transform: scale(1.1);
  box-shadow: 0 6px 16px rgba(30, 58, 138, 0.4);
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
  color: #1e3a8a;
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
  background: linear-gradient(135deg, #1e3a8a 0%, #3b82f6 100%);
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
  color: #1e3a8a;
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
  background: linear-gradient(135deg, #1e3a8a 0%, #3b82f6 100%);
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
  left: 16px;
  right: auto;
}

[dir="rtl"] .close-btn {
  left: auto;
  right: 16px;
}

@media (max-width: 768px) {
  .pwa-install-container {
    bottom: 12px;
    left: 12px;
  }
  
  [dir="rtl"] .pwa-install-container {
    left: 12px;
    right: auto;
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