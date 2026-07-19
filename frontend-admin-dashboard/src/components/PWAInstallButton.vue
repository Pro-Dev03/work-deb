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

<script>
export default {
  name: 'PWAInstallButton',
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
      if (this.$t) {
        return this.$t('pwa.installTitle') || this.t('pwa_install_title')
      } else if (window.i18nStore && window.i18nStore.t) {
        return window.i18nStore.t('pwa.installTitle') || this.t('pwa_install_title')
      }
      return this.t('pwa_install_title')
    },
    pwaInstallText() {
      if (this.isIOS) {
        return this.t('install_on_ios')
      }
      if (this.$t) {
        return this.$t('pwa.installText') || this.t('pwa_install_text')
      } else if (window.i18nStore && window.i18nStore.t) {
        return window.i18nStore.t('pwa.installText') || this.t('pwa_install_text')
      }
      return this.t('pwa_install_text')
    },
    iosModalTitle() {
      return this.t('ios_modal_title')
    },
    iosModalSubtitle() {
      return this.t('ios_modal_subtitle')
    },
    step1Text() {
      return this.t('step1_text')
    },
    step2Text() {
      return this.t('step2_text')
    },
    step3Text() {
      return this.t('step3_text')
    },
    benefit1() {
      return this.t('benefit1')
    },
    benefit2() {
      return this.t('benefit2')
    },
    benefit3() {
      return this.t('benefit3')
    },
    gotItText() {
      return this.t('got_it')
    },
    remindLaterText() {
      return this.t('remind_later')
    }
  },
  methods: {
    t(key) {
      const translations = {
        ar: {
          pwa_install_title: 'تثبيت التطبيق',
          pwa_install_text: 'تثبيت التطبيق',
          install_on_ios: 'تثبيت على iPhone',
          ios_modal_title: 'تثبيت التطبيق على iPhone',
          ios_modal_subtitle: 'احصل على تجربة تطبيق أفضل!',
          step1_text: 'اضغط على زر Share ⎋ في أسفل الشاشة',
          step2_text: 'مرر لأسفل واضغط على Add to Home Screen',
          step3_text: 'اضغط على Add في الزاوية العلوية',
          benefit1: '⚡ تشغيل أسرع',
          benefit2: '📱 أيقونة على الشاشة الرئيسية',
          benefit3: '🎨 تصميم شبيه بالتطبيقات',
          got_it: 'فهمت ✓',
          remind_later: 'ذكرني لاحقاً'
        },
        he: {
          pwa_install_title: 'התקנת אפליקציה',
          pwa_install_text: 'התקן אפליקציה',
          install_on_ios: 'התקנה באייפון',
          ios_modal_title: 'התקנת האפליקציה באייפון',
          ios_modal_subtitle: 'קבל חוויית אפליקציה טובה יותר!',
          step1_text: 'לחץ על כפתור Share ⎋ בתחתית המסך',
          step2_text: 'גרור למטה ולחץ על Add to Home Screen',
          step3_text: 'לחץ על Add בפינה העליונה',
          benefit1: '⚡ טעינה מהירה יותר',
          benefit2: '📱 אייקון במסך הבית',
          benefit3: '🎨 עיצוב דמוי אפליקציה',
          got_it: 'הבנתי ✓',
          remind_later: 'תזכיר לי מאוחר יותר'
        },
        en: {
          pwa_install_title: 'Install App',
          pwa_install_text: 'Install App',
          install_on_ios: 'Install on iPhone',
          ios_modal_title: 'Install App on iPhone',
          ios_modal_subtitle: 'Get a better app experience!',
          step1_text: 'Tap the Share ⎋ button at the bottom',
          step2_text: 'Scroll down and tap Add to Home Screen',
          step3_text: 'Tap Add in the top corner',
          benefit1: '⚡ Faster loading',
          benefit2: '📱 Icon on home screen',
          benefit3: '🎨 App-like design',
          got_it: 'Got it ✓',
          remind_later: 'Remind me later'
        }
      }
      return translations[this.currentLang]?.[key] || translations['ar'][key]
    },
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
        this.showInstallButton = false
      } else {
        // Fallback - show general instructions
        this.showIOSInstructions = true
      }
    },
    remindLater() {
      // Hide for 24 hours (change this number to adjust duration)
      const hoursToHide = 24
      localStorage.setItem('pwa-install-dismissed', Date.now().toString())
      localStorage.setItem('pwa-install-dismiss-duration', hoursToHide.toString())
      this.showIOSInstructions = false
      this.showInstallButton = false
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

    // Check if already dismissed
    const dismissed = localStorage.getItem('pwa-install-dismissed')
    const duration = parseInt(localStorage.getItem('pwa-install-dismiss-duration') || '24')
    if (dismissed && Date.now() - parseInt(dismissed) < duration * 60 * 60 * 1000) {
      return
    }

    // For iOS, always show button
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
  right: 16px;
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