import { reactive, computed } from 'vue'
import ar from '../i18n/ar.json'
import he from '../i18n/he.json'
import en from '../i18n/en.json'

// =============================================
// المفتاح الموحد للغة في جميع التطبيقات
// =============================================
const STORAGE_KEY = 'worktrack_language'

// =============================================
// ترجمة مباشرة من ملفات JSON
// =============================================
const messages = { ar, he, en }

// =============================================
// الحصول على اللغة المخزنة
// =============================================
function getStoredLang() {
  const stored = localStorage.getItem(STORAGE_KEY)
  if (stored && ['ar', 'he', 'en'].includes(stored)) {
    return stored
  }
  // استخدام لغة المتصفح كافتراضي
  const browserLang = navigator.language || navigator.languages?.[0] || 'ar'
  if (browserLang.startsWith('he')) return 'he'
  if (browserLang.startsWith('en')) return 'en'
  return 'ar'
}

// =============================================
// حالة الترجمة
// =============================================
const i18nState = reactive({
  currentLang: getStoredLang(),
  
  setLang(lang) {
    if (['ar', 'he', 'en'].includes(lang)) {
      this.currentLang = lang
      localStorage.setItem(STORAGE_KEY, lang)
      
      // تحديث اتجاه الصفحة
      document.documentElement.dir = lang === 'ar' || lang === 'he' ? 'rtl' : 'ltr'
      document.documentElement.lang = lang
      
      console.log(`🌍 تم تغيير اللغة إلى: ${lang}`)
    }
  },
  
  t(key) {
    const translation = messages[this.currentLang]?.[key]
    if (translation !== undefined && translation !== null) {
      return translation
    }
    const fallback = messages['ar']?.[key]
    if (fallback !== undefined && fallback !== null) {
      return fallback
    }
    console.warn(`⚠️ مفتاح الترجمة غير موجود: ${key}`)
    return key
  }
})

// =============================================
// تصدير الدوال
// =============================================
export function useI18n() {
  const t = (key) => i18nState.t(key)
  const setLang = (lang) => i18nState.setLang(lang)
  const currentLang = computed(() => i18nState.currentLang)
  return { t, setLang, currentLang }
}

export default {
  install(app) {
    app.config.globalProperties.$t = i18nState.t
    app.config.globalProperties.$lang = i18nState.currentLang
    app.provide('i18n', i18nState)
  }
}

export { i18nState }
