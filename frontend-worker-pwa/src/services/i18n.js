import { reactive, computed } from 'vue'
import ar from '../i18n/ar.json'
import he from '../i18n/he.json'
import en from '../i18n/en.json'

// =============================================
// Unified language key for all applications
// =============================================
const STORAGE_KEY = 'worktrack_language'

// =============================================
// قائمة اللغات المدعومة
// =============================================
const SUPPORTED_LANGUAGES = ['ar', 'he', 'en', 'fr', 'de', 'es', 'it', 'pt', 'ru', 'zh', 'ja', 'ko', 'tr', 'nl', 'sv', 'pl']

// =============================================
// Direct translation from JSON files
// =============================================
const messages = { ar, he, en }

// =============================================
// تحميل ترجمات اللغات الإضافية
// =============================================
async function loadLanguage(lang) {
  if (messages[lang]) return messages[lang]
  
  try {
    const response = await fetch(`/i18n/${lang}.json`)
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`)
    }
    const translations = await response.json()
    messages[lang] = translations
    return translations
  } catch (error) {
    console.error(`❌ Failed to load language ${lang}:`, error)
    return messages['en'] // Fallback to English
  }
}

// =============================================
// Get stored language
// =============================================
function getStoredLang() {
  const stored = localStorage.getItem(STORAGE_KEY)
  if (stored && SUPPORTED_LANGUAGES.includes(stored)) {
    return stored
  }
  // Use browser language as default
  const browserLang = navigator.language || navigator.languages?.[0] || 'ar'
  if (browserLang.startsWith('he')) return 'he'
  if (browserLang.startsWith('en')) return 'en'
  if (browserLang.startsWith('fr')) return 'fr'
  if (browserLang.startsWith('de')) return 'de'
  if (browserLang.startsWith('es')) return 'es'
  if (browserLang.startsWith('it')) return 'it'
  if (browserLang.startsWith('pt')) return 'pt'
  if (browserLang.startsWith('ru')) return 'ru'
  if (browserLang.startsWith('zh')) return 'zh'
  if (browserLang.startsWith('ja')) return 'ja'
  if (browserLang.startsWith('ko')) return 'ko'
  if (browserLang.startsWith('tr')) return 'tr'
  if (browserLang.startsWith('nl')) return 'nl'
  if (browserLang.startsWith('sv')) return 'sv'
  if (browserLang.startsWith('pl')) return 'pl'
  return 'ar'
}

// =============================================
// Translation state
// =============================================
const i18nState = reactive({
  currentLang: getStoredLang(),
  
  async setLang(lang) {
    if (SUPPORTED_LANGUAGES.includes(lang)) {
      // تحميل الترجمة إذا لم تكن محملة
      if (!messages[lang]) {
        await loadLanguage(lang)
      }
      
      this.currentLang = lang
      localStorage.setItem(STORAGE_KEY, lang)
      
      // Update page direction
      document.documentElement.dir = lang === 'ar' || lang === 'he' ? 'rtl' : 'ltr'
      document.documentElement.lang = lang
      
      // Send event to indicate language change
      window.dispatchEvent(new CustomEvent('language-changed', { detail: { lang } }))
      
      console.log(`🌍 Language changed to: ${lang}`)
    }
  },
  
  t(key) {
    const keys = key.split('.')
    let translation = messages[this.currentLang]
    
    for (const k of keys) {
      if (translation && translation[k]) {
        translation = translation[k]
      } else {
        // Search in default language
        let fallbackTranslation = messages['en']
        for (const fk of keys) {
          if (fallbackTranslation && fallbackTranslation[fk]) {
            fallbackTranslation = fallbackTranslation[fk]
          } else {
            console.warn(`⚠️ Translation key not found: ${key}`)
            return key
          }
        }
        return fallbackTranslation
      }
    }
    return translation
  }
})

// Add a reactive translation getter
const currentTranslations = computed(() => messages[i18nState.currentLang])

// =============================================
// تصدير الدوال
// =============================================
export function useI18n() {
  const setLang = (lang) => i18nState.setLang(lang)
  const currentLang = computed(() => i18nState.currentLang)
  
  // Return a simple t function and the reactive translations
  const t = (key) => i18nState.t(key)
  
  return { t, setLang, currentLang, currentTranslations }
}

export default {
  install(app) {
    app.config.globalProperties.$t = i18nState.t
    app.config.globalProperties.$lang = i18nState.currentLang
    app.provide('i18n', i18nState)
  }
}

export { i18nState }
