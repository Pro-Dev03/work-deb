# دليل بناء تطبيق Windows - WorkTrack Client Portal

## المتطلبات الأساسية

### 1. تثبيت Node.js و npm
- قم بتحميل وتثبيت [Node.js LTS](https://nodejs.org/) (الإصدار 18 أو أحدث)
- تأكد من تثبيت npm مع Node.js

### 2. تثبيت أدوات بناء Windows
```bash
npm install --global windows-build-tools
```

### 3. تثبيت أداة النشر الإلكتروني
```bash
npm install --global electron-builder
```

## خطوات البناء

### 1. تثبيت الاعتماديات
```bash
cd frontend-client-portal
npm install
```

### 2. بناء التطبيق للإنتاج
```bash
npm run build
```

### 3. بناء تطبيق Windows
```bash
npm run build:electron:win
```

## موقع الملفات المبنية

بعد اكتمال البناء، ستجد الملفات في:
```
frontend-client-portal/dist-electron/
├── WorkTrack-Client-1.0.0-x64.exe
└── builder-effective-config.yaml
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

### المشكلة: WebSocket لا يعمل
**الحل:**
- تم إصلاح هذه المشكلة في آخر التحديثات
- يستخدم التطبيق الآن WebSocket عبر IPC لضمان التوافق مع Windows

## إعدادات المطور

### تشغيل في وضع التطوير
```bash
npm run electron:dev
```

---

**آخر تحديث:** 2026-07-18