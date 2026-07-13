// تنظيف مفاتيح localStorage القديمة
const keysToRemove = [
  'worktrack_lang',
  'worktrack_worker_lang',
  'worktrack_admin_lang',
  'worktrack_client_lang'
];

keysToRemove.forEach(key => {
  localStorage.removeItem(key);
  console.log(`🗑️ تم حذف: ${key}`);
});

// تعيين المفتاح الجديد الموحد
const NEW_LANG_KEY = 'worktrack_language';
localStorage.setItem(NEW_LANG_KEY, 'ar');
console.log(`✅ تم تعيين اللغة الافتراضية: ar`);
