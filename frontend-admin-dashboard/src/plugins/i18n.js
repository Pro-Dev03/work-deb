import { reactive, computed } from 'vue'

// تحميل ملفات الترجمة
import ar from '../locales/ar.json'
import he from '../locales/he.json'
import en from '../locales/en.json'

// إضافة المفاتيح المفقودة للترجمة
const arExtended = {
  ...ar,
  pull_to_refresh: "اسحب للتحديث",
  refreshing: "جاري التحديث...",
  release_to_refresh: "أطلق للتحديث",
  pwa_install_title: "تثبيت التطبيق",
  pwa_install_text: "ثبت التطبيق على جهازك",
  pwa: {
    installTitle: "تثبيت التطبيق",
    installText: "تثبيت التطبيق"
  }
}

const heExtended = {
  ...he,
  pull_to_refresh: "משוך לרענון",
  refreshing: "מרענן...",
  release_to_refresh: "שחרר לרענון",
  pwa_install_title: "התקנת אפליקציה",
  pwa_install_text: "התקן את האפליקציה במכשיר שלך",
  pwa: {
    installTitle: "התקנת אפליקציה",
    installText: "התקן אפליקציה"
  }
}

const enExtended = {
  ...en,
  pull_to_refresh: "Pull to refresh",
  refreshing: "Refreshing...",
  release_to_refresh: "Release to refresh",
  pwa_install_title: "Install App",
  pwa_install_text: "Install the app on your device",
  pwa: {
    installTitle: "Install App",
    installText: "Install App"
  }
}

const messages = {
  ar: arExtended,
  he: heExtended,
  en: enExtended
}

const FALLBACK = 'ar'

const getStoredLang = () => {
  const stored = localStorage.getItem('worktrack_lang')
  if (stored && (stored === 'ar' || stored === 'he' || stored === 'en')) {
    return stored
  }
  return FALLBACK
}

export const i18nStore = reactive({
  lang: getStoredLang(),
  
  setLang(lang) {
    if (lang === 'ar' || lang === 'he' || lang === 'en') {
      this.lang = lang
      localStorage.setItem('worktrack_lang', lang)
      document.documentElement.dir = lang === 'ar' || lang === 'he' ? 'rtl' : 'ltr'
      document.documentElement.lang = lang

      // إرسال event للإشارة بتغيير اللغة
      window.dispatchEvent(new CustomEvent('language-changed', { detail: { lang } }))
    }
  },
  
  t(key) {
    const keys = key.split('.')
    let translation = messages[this.lang]
    
    for (const k of keys) {
      if (translation && translation[k]) {
        translation = translation[k]
      } else {
        // البحث في اللغة الافتراضية
        let fallbackTranslation = messages[FALLBACK]
        for (const fk of keys) {
          if (fallbackTranslation && fallbackTranslation[fk]) {
            fallbackTranslation = fallbackTranslation[fk]
          } else {
            return key
          }
        }
        return fallbackTranslation
      }
    }
    return translation
  }
})

export function useI18n() {
  const t = (key) => i18nStore.t(key)
  const setLang = (lang) => i18nStore.setLang(lang)
  const currentLang = computed(() => i18nStore.lang)
  
  return { t, setLang, currentLang }
}
