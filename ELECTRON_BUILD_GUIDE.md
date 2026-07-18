# دليل بناء تطبيقات سطح المكتب والموبايل

## 📱 المحتوى
- [PWA (تطبيق الموظف)](#pwa-تطبيق-الموظف)
- [Electron (سطح المكتب)](#electron-سطح-المكتب)
- [بناء التطبيقات](#بناء-التطبيقات)

---

## 🌐 PWA (تطبيق الموظف)

تم تحسين تطبيق الموظف PWA بشكل كامل ليعمل على الموبايل والويب.

### الميزات المضافة:
- ✅ Service Worker متقدم مع دعم offline
- ✅ Manifest.json محسّن مع أيقونات واختصارات
- ✅ دعم RTL للغة العربية
- ✅ تثبيت التطبيق على الموبايل

### تعليمات التثبيت على الموبايل:
1. افتح التطبيق في المتصفح (Chrome/Safari)
2. سيظهر زر "Add to Home Screen"
3. اضغط عليه واتبع التعليمات

---

## 💻 Electron (سطح المكتب)

تم إضافة Electron لجميع الواجهات الثلاث:
- `frontend-admin-dashboard` - لوحة تحكم الإدارة
- `frontend-worker-pwa` - تطبيق الموظف
- `frontend-client-portal` - بوابة العملاء

### الملفات المضافة:
- `electron/main.js` - ملف الرئيسي للإطار
- `electron/preload.js` - ملف preload للأمان
- تحديثات `package.json` مع scripts جديدة

### أوامر التطوير:

```bash
# لوحة تحكم الإدارة
cd frontend-admin-dashboard
npm run electron:dev

# تطبيق الموظف
cd frontend-worker-pwa
npm run electron:dev

# بوابة العملاء
cd frontend-client-portal
npm run electron:dev
```

---

## 🏗️ بناء التطبيقات

### للويب (العادي):
```bash
cd frontend-admin-dashboard
npm run build
```

### لسطح المكتب (Windows):
```bash
cd frontend-admin-dashboard
npm run build:electron:win
```

### لسطح المكتب (Mac):
```bash
cd frontend-admin-dashboard
npm run build:electron:mac
```

### لسطح المكتب (Linux):
```bash
cd frontend-admin-dashboard
npm run build:electron:linux
```

### لجميع المنصات:
```bash
cd frontend-admin-dashboard
npm run build:electron
```

---

## 📦 الملفات الناتجة

بعد البناء، ستجد الملفات في مجلد `dist-electron/`:

### Windows:
- `WorkTrack Admin Setup.exe` - ملف التثبيت

### Mac:
- `WorkTrack Admin.dmg` - ملف DMG

### Linux:
- `WorkTrack Admin.AppImage` - ملف AppImage
- `worktrack-admin_1.0.0_amd64.deb` - ملف DEB

---

## 🔧 إعدادات مخصصة

يمكنك تعديل الإعدادات في `package.json`:

```json
"build": {
  "appId": "com.worktrack.admin",
  "productName": "WorkTrack Admin",
  "directories": {
    "output": "dist-electron"
  }
}
```

---

## 🎨 تخصيص الأيقونات

استبدل `public/favicon.ico` بأيقونة مخصصة للحصول على مظهر أفضل.

---

## 📝 ملاحظات مهمة

1. **أمان**: تم تعطيل `nodeIntegration` وتمكين `contextIsolation` للأمان
2. **الأداء**: يتم استخدام Vite للبناء السريع
3. **الحجم**: التطبيقات قد تكون كبيرة بسبب Electron (حوالي 100-150MB)
4. **المنافذ**: تم تعيين منافذ مختلفة لكل واجهة:
   - Admin: 3001
   - Worker: 3000
   - Client: 3002

---

## 🚀 الخطوات التالية

لتحويل هذا إلى تطبيقات موبايل أصلية (Native):

1. إضافة Capacitor:
```bash
npm install @capacitor/core @capacitor/cli
npx cap init
npm install @capacitor/android @capacitor/ios
npx cap add android
npx cap add ios
```

2. البناء والمزامنة:
```bash
npm run build
npx cap sync
npx cap open android
```

---

## 🐛 استكشاف الأخطاء

### المشكلة: التطبيق لا يفتح
- تأكد من بناء المشروع أولاً: `npm run build`
- تحقق من مسار الملفات في `electron/main.js`

### المشكلة: المنفذ مستخدم
- غيّر المنفذ في `vite.config.js` و `electron/main.js`

### المشكلة: Service Worker لا يعمل
- تأكد من وجود `public/service-worker.js`
- تحقق من التسجيل في `main.js`

---

## 📞 الدعم

للمزيد من المساعدة، راجع:
- [Electron Documentation](https://www.electronjs.org/docs)
- [Capacitor Documentation](https://capacitorjs.com/docs)
- [PWA Documentation](https://web.dev/progressive-web-apps/)