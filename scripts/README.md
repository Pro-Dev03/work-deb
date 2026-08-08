# 🔧 مجلد السكريبتات (Scripts)

هذا المجلد يحتوي على جميع السكريبتات المستخدمة في المشروع.

## 📁 الهيكل

### build/
سكريبتات البناء:
- `build-all-release.sh` - بناء جميع التطبيقات للإصدار الإنتاجي
- `create-keystore.sh` - إنشاء Keystore الرسمي للتطبيقات
- `fix-android-icons.sh` - إصلاح أيقونات Android

### setup/
سكريبتات الإعداد والتهيئة:
- `generate-keystore-auto.sh` - إنشاء Keystore تلقائي
- `setup-codemagic.sh` - إعداد Codemagic
- `setup-codemagic-api.sh` - إعداد Codemagic API
- `change_credentials_sql.sh` - تغيير بيانات الاعتماد في SQL

### python/
سكريبتات بايثون:
- `extract_project_code.py` - استخراج وتنظيم أكواد المشروع
- `convert_html_to_pdf.py` - تحويل HTML إلى PDF
- `extract_project_structure.py` - استخراج هيكل المشروع
- `simple_pdf_converter.py` - محول PDF بسيط

### extraction/
مجلد إخراج السكريبتات:
- يحتوي على نتائج استخراج الأكواد المنظمة

## 🚀 كيفية الاستخدام

### سكريبتات البناء (Bash)
```bash
# بناء جميع التطبيقات
./scripts/build/build-all-release.sh

# إنشاء Keystore
./scripts/build/create-keystore.sh

# إصلاح الأيقونات
./scripts/build/fix-android-icons.sh
```

### سكريبتات الإعداد (Bash)
```bash
# إعداد Codemagic
./scripts/setup/setup-codemagic.sh

# تغيير بيانات الاعتماد
./scripts/setup/change_credentials_sql.sh
```

### سكريبتات بايثون
```bash
# استخراج الأكواد
python3 scripts/python/extract_project_code.py --path /home/dev-bit/project/worktrack

# استخراج مع Mobile
python3 scripts/python/extract_project_code.py --path /home/dev-bit/project/worktrack --mobile
```

## 📝 ملاحظات هامة

- تأكد من منح صلاحيات التنفيذ لسكريبتات Bash: `chmod +x script.sh`
- بعض السكريبتات تتطلب Python 3.x
- اقرأ التعليقات داخل كل سكريبت قبل الاستخدام
- احتفظ بنسخ احتياطية قبل تشغيل سكريبتات حساسة (مثل Keystore)

## ⚠️ تحذيرات

- سكريبتات Keystore حساسة - احتفظ بالملفات الناتجة في مكان آمن
- سكريبتات SQL تغيّر قاعدة البيانات - استخدمها بحذر
- تأكد من اختبار السكريبتات على بيئة staging أولاً