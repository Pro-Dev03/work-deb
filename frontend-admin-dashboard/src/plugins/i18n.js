import { reactive, computed } from 'vue'

// تحميل ملفات الترجمة
import ar from '../locales/ar.json'
import he from '../locales/he.json'
import en from '../locales/en.json'

const messages = {
  ar,
  he,
  en
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
