import { reactive, computed } from 'vue'

// اللغات المدعومة
const messages = {
  ar: {
    app_name: "WorkTrack", 
    admin_panel: "لوحة المدير",
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
    pull_to_refresh: "اسحب للتحديث",
    refreshing: "جاري التحديث...",
    release_to_refresh: "أطلق للتحديث",
    pwa_install_title: "تثبيت التطبيق",
    pwa_install_text: "ثبت التطبيق على جهازك"
  },
  he: {
    app_name: "WorkTrack", 
    admin_panel: "לוח המנהל",
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
    pull_to_refresh: "משוך לרענון",
    refreshing: "מרענן...",
    release_to_refresh: "שחרר לרענון",
    pwa_install_title: "התקנת אפליקציה",
    pwa_install_text: "התקן את האפליקציה במכשיר שלך"
  },
  en: {
    app_name: "WorkTrack", 
    admin_panel: "Admin Panel",
    login: "Login",
    email: "Email",
    password: "Password",
    dashboard: "Dashboard",
    employees: "Employees",
    tasks: "Tasks",
    worksites: "Worksites",
    clients: "Clients",
    reports: "Reports",
    settings: "Settings",
    service_requests: "Service Requests",
    logout: "Logout",
    add_employee: "Add Employee",
    assign: "Assign",
    status: "Status",
    pending: "Pending",
    assigned: "Assigned",
    in_progress: "In Progress",
    completed: "Completed",
    cancelled: "Cancelled",
    open_in_maps: "Open in Maps",
    location: "Location",
    phone: "Phone",
    full_name: "Full Name",
    save: "Save",
    cancel: "Cancel",
    delete: "Delete",
    edit: "Edit",
    search: "Search",
    filter: "Filter",
    dark_mode: "Dark Mode",
    light_mode: "Light Mode",
    language: "Language",
    arabic: "Arabic",
    hebrew: "Hebrew",
    loading: "Loading...",
    no_data: "No Data",
    error: "Error",
    success: "Success",
    welcome: "Welcome",
    pull_to_refresh: "Pull to refresh",
    refreshing: "Refreshing...",
    release_to_refresh: "Release to refresh",
    pwa_install_title: "Install App",
    pwa_install_text: "Install the app on your device"
  }
}

const FALLBACK = 'ar'

const getStoredLang = () => {
  const stored = localStorage.getItem('worktrack_lang')
  if (stored && (stored === 'ar' || stored === 'he' || stored === 'en')) return stored
  return FALLBACK
}

export const i18nStore = reactive({
  lang: getStoredLang(),
  
  setLang(lang) {
    if (lang === 'ar' || lang === 'he' || lang === 'en') {
      this.lang = lang
      localStorage.setItem('worktrack_lang', lang)
      window.location.reload()
    }
  },
  
  t(key) {
    const translation = messages[this.lang]?.[key]
    if (translation) return translation
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
