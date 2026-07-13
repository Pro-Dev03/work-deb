import { reactive, computed } from 'vue'

// ==========================================================================
// اللغات المدعومة: العربية والعبرية فقط
// ==========================================================================
const messages = {
  ar: {
    app_name: "ورك تراك",
    admin_panel: "لوحة المدير",
    worker_app: "تطبيق الموظف",
    login: "تسجيل الدخول",
    email: "البريد الإلكتروني",
    password: "كلمة المرور",
    dashboard: "لوحة القيادة",
    employees: "الموظفون",
    tasks: "المهام",
    worksites: "المواقع",
    clients: "العملاء",
    reports: "التقارير",
    settings: "الإعدادات",
    service_requests: "طلبات الخدمة",
    logout: "تسجيل الخروج",
    add_employee: "إضافة موظف",
    assign: "تعيين",
    status: "الحالة",
    pending: "قيد الانتظار",
    assigned: "تم التعيين",
    in_progress: "قيد التنفيذ",
    completed: "مكتمل",
    cancelled: "ملغي",
    open_in_maps: "فتح في الخرائط",
    location: "الموقع",
    phone: "الهاتف",
    full_name: "الاسم الكامل",
    save: "حفظ",
    cancel: "إلغاء",
    delete: "حذف",
    edit: "تعديل",
    search: "بحث",
    filter: "تصفية",
    dark_mode: "الوضع الداكن",
    light_mode: "الوضع الفاتح",
    language: "اللغة",
    arabic: "العربية",
    hebrew: "עברית",
    loading: "جارٍ التحميل...",
    no_data: "لا توجد بيانات",
    error: "حدث خطأ",
    success: "تم بنجاح",
    welcome: "مرحباً بك",
    login_title: "لوحة تحكم المدير",
    worker_title: "تطبيق الموظف"
  },
  he: {
    app_name: "וורק טראק",
    admin_panel: "לוח המנהל",
    worker_app: "אפליקציית העובד",
    login: "התחברות",
    email: "אימייל",
    password: "סיסמה",
    dashboard: "לוח בקרה",
    employees: "עובדים",
    tasks: "משימות",
    worksites: "אתרי עבודה",
    clients: "לקוחות",
    reports: "דוחות",
    settings: "הגדרות",
    service_requests: "בקשות שירות",
    logout: "התנתקות",
    add_employee: "הוסף עובד",
    assign: "הקצה",
    status: "סטטוס",
    pending: "ממתין",
    assigned: "הוקצה",
    in_progress: "בתהליך",
    completed: "הושלם",
    cancelled: "בוטל",
    open_in_maps: "פתח במפות",
    location: "מיקום",
    phone: "טלפון",
    full_name: "שם מלא",
    save: "שמור",
    cancel: "בטל",
    delete: "מחק",
    edit: "ערוך",
    search: "חפש",
    filter: "סנן",
    dark_mode: "מצב כהה",
    light_mode: "מצב בהיר",
    language: "שפה",
    arabic: "ערבית",
    hebrew: "עברית",
    loading: "טוען...",
    no_data: "אין נתונים",
    error: "שגיאה",
    success: "הצלחה",
    welcome: "ברוך הבא",
    login_title: "לוח המנהל",
    worker_title: "אפליקציית העובד"
  }
}

const FALLBACK = 'ar'

// الحصول على اللغة المخزنة
const getStoredLang = () => {
  const stored = localStorage.getItem('worktrack_lang')
  // التأكد من أن اللغة مدعومة
  if (stored && (stored === 'ar' || stored === 'he')) return stored
  return FALLBACK
}

export const i18nStore = reactive({
  lang: getStoredLang(),
  
  setLang(lang) {
    if (lang === 'ar' || lang === 'he') {
      this.lang = lang
      localStorage.setItem('worktrack_lang', lang)
      // إعادة تحميل الصفحة لتطبيق اللغة الجديدة
      window.location.reload()
    }
  },
  
  t(key) {
    const translation = messages[this.lang]?.[key]
    if (translation) return translation
    // استخدام العربية كافتراضي إذا لم توجد الترجمة
    const fallback = messages[FALLBACK]?.[key]
    return fallback || key
  }
})

export function useI18n() {
  const t = (key) => i18nStore.t(key)
  const setLang = (lang) => i18nStore.setLang(lang)
  const currentLang = computed(() => i18nStore.lang)
  return { t, setLang, currentLang }
}
