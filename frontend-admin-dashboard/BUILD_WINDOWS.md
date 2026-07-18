# دليل بناء تطبيق Windows - WorkTrack Admin Dashboard

## المتطلبات الأساسية

### 1. تثبيت Node.js و npm
- قم بتحميل وتثبيت [Node.js LTS](https://nodejs.org/) (الإصدار 18 أو أحدث)
- تأكد من تثبيت npm مع Node.js

### 2. تثبيت أدوات بناء Windows
```bash
# تثبيت windows-build-tools (للتطوير على Windows)
npm install --global windows-build-tools
```

### 3. تثبيت أداة النشر الإلكتروني
```bash
npm install --global electron-builder
```

## خطوات البناء

### 1. تثبيت الاعتماديات
```bash
cd frontend-admin-dashboard
npm install
```

### 2. بناء التطبيق للإنتاج
```bash
npm run build
```

### 3. بناء تطبيق Windows
```bash
# بناء كل من NSIS و Portable
npm run build:electron:win

# أو بناء Portable فقط
npm run build:electron:win:portable
```

## موقع الملفات المبنية

بعد اكتمال البناء، ستجد الملفات في:
```
frontend-admin-dashboard/dist-electron/
├── WorkTrack-Admin-1.0.0-x64.exe        # NSIS Installer
├── WorkTrack-Admin-1.0.0-x64-portable.exe # Portable Version
└── builder-effective-config.yaml       # إعدادات البناء
```

## حل المشاكل الشائعة

### المشكلة: التطبيق لا يفتح
**الحل:**
1. تأكد من تثبيت جميع الاعتماديات:
   ```bash
   npm install
   ```

2. قم بتنظيف ملفات البناء السابقة:
   ```bash
   rm -rf dist dist-electron node_modules
   npm install
   npm run build
   npm run build:electron:win
   ```

3. تأكد من وجود ملفات الأيقونات:
   ```bash
   ls public/icon.png
   ```

### المشكلة: أخطاء في البناء
**الحل:**
1. تحديث electron-builder:
   ```bash
   npm install --save-dev electron-builder@latest
   ```

2. تحديث Electron:
   ```bash
   npm install --save-dev electron@latest
   ```

### المشكلة: WebSocket لا يعمل
**الحل:**
- تم إصلاح هذه المشكلة في آخر التحديثات
- يستخدم التطبيق الآن WebSocket عبر IPC لضمان التوافق مع Windows

## إعدادات المطور

### تشغيل في وضع التطوير
```bash
npm run electron:dev
```

### فتح أدوات المطور
- اضغط على `F12` أو استخدم القائمة: عرض → مطور

## اختبار التطبيق قبل البناء

### 1. اختبار محلي
```bash
npm run dev
```
افتح المتصفح على: http://localhost:3001

### 2. اختبار Electron محلي
```bash
npm run electron
```

## التوزيع

### توقيع الكود (للإنتاج)
أضف إلى `.env`:
```
WIN_CSC_LINK=path/to/certificate.pfx
WIN_CSC_KEY_PASSWORD=your_password
```

### النشر على GitHub Releases
```bash
npm run build:electron:win --publish always
```

## ملاحظات مهمة

1. **الأيقونات**: تأكد من وجود `public/icon.png` بالحجم المناسب (512x512)
2. **WebSocket**: يعمل تلقائياً عبر IPC في Electron
3. **CORS**: تمت معالجته في إعدادات Electron
4. **الأمان**: تم تعطيل webSecurity فقط في وضع التطوير

## الدعم الفني

إذا واجهت مشاكل:
1. تحقق من سجلات Electron (افتح DevTools)
2. تأكد من توافق إصدارات Node.js و Electron
3. قم بتحديث جميع الاعتماديات: `npm update`

---

**آخر تحديث:** 2026-07-18
**الإصدار:** 1.0.0