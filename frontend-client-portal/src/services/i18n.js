import { reactive, computed } from 'vue'
import ar from '../locales/ar.json'
import he from '../locales/he.json'
import en from '../locales/en.json'

// =============================================
// المفتاح الموحد للغة في جميع التطبيقات
// =============================================
const STORAGE_KEY = 'worktrack_language'

// =============================================
// تُضمَّن الترجمات ضمن حزمة Vite؛ طلب ملفات src عبر fetch لا يعمل بعد النشر.
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
      
      // إرسال event للإشارة بتغيير اللغة
      window.dispatchEvent(new CustomEvent('language-changed', { detail: { lang } }))
      
      console.log(`🌍 تم تغيير اللغة إلى: ${lang}`)
      
      // إعادة تحميل الصفحة لتطبيق التغييرات
      setTimeout(() => {
        window.location.reload()
      }, 300)
    }
  },
  
  t(key) {
    const keys = key.split('.')
    let translation = messages[this.currentLang]

    for (const k of keys) {
      if (translation && translation[k]) {
        translation = translation[k]
      } else {
        // البحث في اللغة الافتراضية
        let fallbackTranslation = messages['ar']
        for (const fk of keys) {
          if (fallbackTranslation && fallbackTranslation[fk]) {
            fallbackTranslation = fallbackTranslation[fk]
          } else {
            console.warn(`⚠️ مفتاح الترجمة غير موجود: ${key}`)
            return key
          }
        }
        return fallbackTranslation
      }
    }
    return translation
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
