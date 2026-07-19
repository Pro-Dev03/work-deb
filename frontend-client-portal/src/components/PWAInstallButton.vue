<template>
  <div v-if="showInstallButton" class="pwa-install-container">
    <button 
      @click="installPWA" 
      class="pwa-install-button"
      :title="pwaInstallTitle"
    >
      <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path>
        <polyline points="7 10 12 15 17 10"></polyline>
        <line x1="12" y1="15" x2="12" y2="3"></line>
      </svg>
      <span>{{ pwaInstallText }}</span>
    </button>
  </div>
</template>

<script>
export default {
  name: 'PWAInstallButton',
  data() {
    return {
      showInstallButton: false,
      deferredPrompt: null
    }
  },
  computed: {
    pwaInstallTitle() {
      return this.$t('pwa.installTitle') || 'تثبيت التطبيق'
    },
    pwaInstallText() {
      return this.$t('pwa.installText') || 'تثبيت التطبيق'
    }
  },
  mounted() {
    // الاستماع إلى event توفر التثبيت
    window.addEventListener('pwa-install-available', this.handleInstallAvailable)
    window.addEventListener('pwa-install-success', this.handleInstallSuccess)
    
    // التحقق من وجود deferredPrompt في window
    if (window.deferredPrompt) {
      this.showInstallButton = true
    }
  },
  beforeUnmount() {
    window.removeEventListener('pwa-install-available', this.handleInstallAvailable)
    window.removeEventListener('pwa-install-success', this.handleInstallSuccess)
  },
  methods: {
    handleInstallAvailable() {
      this.showInstallButton = true
    },
    handleInstallSuccess() {
      this.showInstallButton = false
    },
    installPWA() {
      // استخدام الدالة العامة من main.js
      if (window.pwaInstall) {
        window.pwaInstall()
        this.showInstallButton = false
      }
    }
  }
}
</script>

<style scoped>
.pwa-install-container {
  position: fixed;
  bottom: 20px;
  left: 20px;
  z-index: 1000;
}

.pwa-install-button {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 20px;
  background: linear-gradient(135deg, #3880ff 0%, #5aa0ff 100%);
  color: white;
  border: none;
  border-radius: 12px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  box-shadow: 0 4px 15px rgba(56, 128, 255, 0.3);
  transition: all 0.3s ease;
}

.pwa-install-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(56, 128, 255, 0.4);
}

.pwa-install-button:active {
  transform: translateY(0);
}

/* RTL support */
[dir="rtl"] .pwa-install-container {
  left: auto;
  right: 20px;
}

[dir="rtl"] .pwa-install-button {
  flex-direction: row-reverse;
}

@media (max-width: 768px) {
  .pwa-install-container {
    bottom: 16px;
    left: 16px;
    right: 16px;
  }
  
  [dir="rtl"] .pwa-install-container {
    left: 16px;
    right: 16px;
  }
  
  .pwa-install-button {
    width: 100%;
    justify-content: center;
  }
}
</style>