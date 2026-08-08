# Security & Obfuscation

## تم الاستخراج في: 2026-08-06 02:09:35

## عدد الملفات: 10

---

## 📄 KEYSTORE_SETUP.md

```markdown
# دليل إعداد Signed Keystore لتطبيقات WorkTrack

## ما هو Keystore؟

Keystore هو ملف أمان يحتوي على المفاتيح الرقمية المستخدمة لتوقيع تطبيقات Android. هذا التوقيع مطلوب لنشر التطبيقات على Google Play Store.

## إنشاء Keystore

### 1. تثبيت Java JDK (إذا لم يكن مثبتاً)

```bash
# التحقق من تثبيت Java
java -version

# إذا لم يكن مثبتاً، ثبته:
sudo apt update
sudo apt install default-jdk
```

### 2. إنشاء Keystore

```bash
# الانتقال إلى مجلد المشروع
cd /home/dev-bit/project/worktrack

# إنشاء مجلد للـ keystores
mkdir -p android-keystores

# إنشاء keystore للتطبيقات
keytool -genkey -v -keystore android-keystores/worktrack-release.keystore \
  -alias worktrack \
  -keyalg RSA \
  -keysize 2048 \
  -validity 10000 \
  -storepass your_secure_password_here \
  -keypass your_secure_password_here \
  -dname "CN=WorkTrack, OU=Development, O=WorkTrack, L=YourCity, ST=YourState, C=YourCountry"
```

### 3. معلومات المطلوبة عند الإنشاء

عند تشغيل الأمر، سيطلب منك المعلومات التالية:

- **Keystore password**: كلمة مرور قوية للملف (احفظها بأمان!)
- **Key password**: كلمة مرور للمفتاح (نفس كلمة مرور الـ keystore أو مختلفة)
- **First and Last Name**: اسم المطور أو الشركة
- **Organizational Unit**: القسم (مثلاً: Development)
- **Organization**: اسم الشركة
- **City or Locality**: المدينة
- **State or Province**: الولاية/المقاطعة
- **Country Code**: رمز الدولة (مثلاً: SA, US, GB)

### 4. تأمين Keystore

```bash
# تغيير الأذونات
chmod 600 android-keystores/worktrack-release.keystore

# إضافة الملف إلى .gitignore
echo "android-keystores/" >> .gitignore
```

## تكوين Gradle

### 1. إنشاء ملف keystore.properties

```bash
# إنشاء الملف في مجلد كل تطبيق
cd frontend-worker-pwa/android
cat > keystore.properties << 'EOF'
RELEASE_STORE_FILE=../../android-keystores/worktrack-release.keystore
RELEASE_KEY_ALIAS=worktrack
RELEASE_STORE_PASSWORD=your_secure_password_here
RELEASE_KEY_PASSWORD=your_secure_password_here
EOF
```

### 2. إضافة keystore.properties إلى .gitignore

```bash
cd /home/dev-bit/project/worktrack
echo "*/android/keystore.properties" >> .gitignore
```

### 3. تحديث build.gradle للتطبيق

لكل تطبيق (worker, admin, client)، أضف التكوين التالي:

```gradle
// في android/app/build.gradle

android {
    // ... التكوين الحالي ...

    signingConfigs {
        release {
            if (project.hasProperty('RELEASE_STORE_FILE')) {
                storeFile file(RELEASE_STORE_FILE)
                storePassword RELEASE_STORE_PASSWORD
                keyAlias RELEASE_KEY_ALIAS
                keyPassword RELEASE_KEY_PASSWORD
            }
        }
    }

    buildTypes {
        release {
            minifyEnabled true
            proguardFiles getDefaultProguardFile('proguard-android.txt'), 'proguard-rules.pro'
            signingConfig signingConfigs.release
        }
    }
}
```

### 4. تحميل keystore.properties في build.gradle

أضف هذا في بداية `android/app/build.gradle`:

```gradle
// تحميل keystore.properties
def keystorePropertiesFile = rootProject.file("keystore.properties")
def keystoreProperties = new Properties()
if (keystorePropertiesFile.exists()) {
    keystoreProperties.load(new FileInputStream(keystorePropertiesFile))
}

android {
    // ... rest of config ...
}
```

## بناء Release APK

### 1. بناء Release APK

```bash
# تطبيق الموظف
cd frontend-worker-pwa
npm run build:obfuscated
npx cap sync android
cd android
./gradlew assembleRelease

# تطبيق الإدارة
cd ../../frontend-admin-dashboard
npm run build:obfuscated
npx cap sync android
cd android
./gradlew assembleRelease

# تطبيق العميل
cd ../../frontend-client-portal
npm run build:obfuscated
npx cap sync android
cd android
./gradlew assembleRelease
```

### 2. موقع ملفات APK

بعد البناء، ستجد ملفات APK في:

```
frontend-worker-pwa/android/app/build/outputs/apk/release/app-release.apk
frontend-admin-dashboard/android/app/build/outputs/apk/release/app-release.apk
frontend-client-portal/android/app/build/outputs/apk/release/app-release.apk
```

## الأمان الهام

### ⚠️ تحذيرات أمنية:

1. **لا تضيف keystore إلى Git** - أبداً!
2. **احتفظ بنسخة احتياطية** في مكان آمن
3. **استخدم كلمات مرور قوية** - على الأقل 16 حرف
4. **لا تشارك Keystore** - أبداً مع أي شخص
5. **احتفظ بالمعلومات** - اكتب كلمات المرور واحفظها في مكان آمن

### إذا فقدت Keystore:

- **لن تتمكن** من تحديث التطبيق على Google Play
- ستحتاج **نشر تطبيق جديد** باسم package مختلف
- ستفقد **المستخدمين الحاليين** والتقييمات

## النسخ الاحتياطي للـ Keystore

```bash
# إنشاء نسخة احتياطية
cp android-keystores/worktrack-release.keystore android-keystores/worktrack-release.keystore.backup

# تشفير النسخة الاحتياطية (اختياري)
gpg --symmetric --cipher-algo AES256 android-keystores/worktrack-release.keystore.backup

# تخزين النسخة المشفرة في مكان آمن (حسب السحابة، USB، إلخ)
```

## التحقق من التوقيع

```bash
# التحقق من توقيع APK
jarsigner -verify -verbose -certs android/app/build/outputs/apk/release/app-release.apk

# عرض معلومات التوقيع
keytool -printcert -jarfile android/app/build/outputs/apk/release/app-release.apk
```

## إنشاء keystore.properties لكل تطبيق

لتسهيل العملية، أنشئ ملفات keystore.properties لكل تطبيق:

```bash
# تطبيق الموظف
cd frontend-worker-pwa/android
cat > keystore.properties << 'EOF'
RELEASE_STORE_FILE=../../android-keystores/worktrack-release.keystore
RELEASE_KEY_ALIAS=worktrack
RELEASE_STORE_PASSWORD=your_secure_password_here
RELEASE_KEY_PASSWORD=your_secure_password_here
EOF

# تطبيق الإدارة
cd ../../frontend-admin-dashboard/android
cat > keystore.properties << 'EOF'
RELEASE_STORE_FILE=../../android-keystores/worktrack-release.keystore
RELEASE_KEY_ALIAS=worktrack
RELEASE_STORE_PASSWORD=your_secure_password_here
RELEASE_KEY_PASSWORD=your_secure_password_here
EOF

# تطبيق العميل
cd ../../frontend-client-portal/android
cat > keystore.properties << 'EOF'
RELEASE_STORE_FILE=../../android-keystores/worktrack-release.keystore
RELEASE_KEY_ALIAS=worktrack
RELEASE_STORE_PASSWORD=your_secure_password_here
RELEASE_KEY_PASSWORD=your_secure_password_here
EOF
```

## ملاحظات هامة

1. **استخدم نفس Keystore** لجميع التطبيقات الثلاثة أو أنشئ keystore منفصل لكل تطبيق
2. **تأكد من صحة المسارات** في keystore.properties
3. **اختبر البناء** قبل الرفع إلى Google Play
4. **احتفظ بكلمات المرور** في مكان آمن جداً
5. **راجع معلومات المطور** في Google Play Console قبل الرفع

## المشاكل الشائعة والحلول

### مشكلة: Keystore not found
**الحل**: تأكد من صحة المسار في keystore.properties

### مشكلة: Wrong password
**الحل**: تحقق من كلمات المرور في keystore.properties

### مشكلة: Build failed
**الحل**: تأكد من أن keystore.properties محمل بشكل صحيح في build.gradle

### مشكلة: Signing config not found
**الحل**: تأكد من إضافة signingConfigs في build.gradle
```

---

## 📄 PRIVACY_POLICY.md

```markdown
# سياسة الخصوصية لتطبيق WorkTrack

## مقدمة
تحترم WorkTrack خصوصيتك وتلتزم بحماية بياناتك الشخصية. توضح هذه السياسة كيفية جمع واستخدام وحماية بياناتك عند استخدام تطبيقات WorkTrack (تطبيق الموظف، تطبيق الإدارة، تطبيق العميل).

## البيانات التي نجمعها

### 1. بيانات الموقع الجغرافي (Location Data)
- **ما نجمعه**: الموقع الجغرافي الحالي، المسار التاريخي للموظفين أثناء ساعات العمل، سرعة الحركة، الاتجاه
- **طريقة الجمع**: GPS، Wi-Fi، شبكات الهاتف
- **الغرض**: 
  - تتبع الحضور والانصراف بدقة
  - التحقق من وجود الموظف في موقع العمل المحدد (Geofencing)
  - تحسين سلامة الموظفين في الميدان
  - إدارة توزيع الموظفين على نقاط العمل
- **الأساس القانوني**: ضروري لأداء العقد (اتفاقية العمل) وتحسين سلامة الموظفين
- **الدقة**: موقع دقيق حتى 10 أمتار (عند تفعيل GPS عالي الدقة)

### 2. البيانات الشخصية (Personal Information)
- **ما نجمعه**: 
  - الاسم الكامل (الكامل الأول والأخير)
  - البريد الإلكتروني
  - رقم الهاتف
  - صورة الملف الشخصي
  - العنوان (اختياري)
  - معلومات التوظيف (المنصب، القسم، تاريخ البدء)
- **الغرض**: 
  - إدارة الحساب والتوثيق
  - التواصل بشأن المهام والتحديثات
  - التعريف في التطبيق وتنظيم الفرق
  - إدارة شؤون الموظفين
- **الأساس القانوني**: موافقة المستخدم وضروري لأداء العقد

### 3. بيانات الحضور والعمل (Attendance & Work Data)
- **ما نجمعه**: 
  - سجلات الحضور والانصراف (أوقات الدخول والخروج)
  - إجمالي ساعات العمل اليومية والأسبوعية والشهرية
  - نقاط العمل المزارة والمدة في كل نقطة
  - المهام المخصصة وحالة إكمالها
  - الملاحظات والتقارير
  - تاريخ ووقت كل نشاط
- **الغرض**: 
  - إدارة شؤون الموظفين وحساب الرواتب
  - تتبع الإنتاجية وتقييم الأداء
  - التخطيط وتوزيع الموارد
  - التقارير الإدارية والتحليلية
- **الأساس القانوني**: ضروري لأداء العقد والامتثال الضريبي

### 4. الصور والملفات (Photos & Files)
- **ما نجمعه**: 
  - صور الملفات الشخصية
  - صور التحقق من الحضور (Attendance verification photos)
  - صور المهام المكتملة
  - المستندات المتعلقة بالعمل
  - لقطات الشاشة والتقارير
- **الغرض**: 
  - التحقق من الهوية والحضور
  - توثيق إكمال المهام
  - إدارة المستندات
  - دعم عمليات التدقيق
- **الأساس القانوني**: موافقة المستخدم وضروري لأداء العقد

### 5. بيانات الجهاز (Device Data)
- **ما نجمعه**: 
  - نوع الجهاز والطراز
  - نظام التشغيل والإصدار
  - معرف الجهاز الفريد (Device ID)
  - معلومات البطارية والاتصال
  - دقة الشاشة والأبعاد
- **الغرض**: 
  - تحسين أداء التطبيق وتوافقه
  - الدعم الفني وتشخيص المشاكل
  - تحليلات الاستخدام لتحسين الميزات
  - اكتشاف الأخطاء وإصلاحها
- **الأساس القانوني**: مصلحة مشروعة لتحسين الخدمة

### 6. بيانات الاستخدام (Usage Data)
- **ما نجمعه**: 
  - تسجيلات الدخول والخروج
  - الأوقات التي يقضيها المستخدم في التطبيق
  - الميزات المستخدمة وتكرار استخدامها
  - أخطاء التطبيق وتعطل النظام
  - إحصائيات الأداء (وقت الاستجابة، سرعة التحميل)
- **الغرض**: 
  - تحسين تجربة المستخدم
  - تحديد الميزات الأكثر أهمية
  - اكتشاف وإصلاح المشاكل
  - تطوير ميزات جديدة
- **الأساس القانوني**: مصلحة مشروعة لتحسين الخدمة

## كيفية استخدام بياناتك

### الأغراض الرئيسية:

#### 1. إدارة شؤون الموظفين
- تتبع الحضور والانصراف بدقة
- حساب ساعات العمل والرواتب
- إدارة الإجازات والمهام
- تقييم الأداء والإنتاجية

#### 2. التحقق من الموقع والأمان
- التأكد من وجود الموظف في موقع العمل المحدد
- إنشاء مناطق جغرافية (Geofencing) لنقاط العمل
- تحسين سلامة الموظفين في الميدان
- تتبع الحركة في حالات الطوارئ

#### 3. إدارة المهام والمشاريع
- توزيع المهام على الموظفين
- تتبع حالة إكمال المهام
- التنسيق بين الفرق
- إدارة العملاء وطلبات الخدمة

#### 4. التواصل والإشعارات
- إرسال إشعارات المهام الجديدة
- تنبيهات الهامة والعاجلة
- تحديثات الحالة والتذكيرات
- التواصل مع الإدارة

#### 5. تحسين التطبيق
- تحليلات الاستخدام لتحسين الميزات
- اكتشاف وإصلاح الأخطاء
- تطوير ميزات جديدة بناءً على الاحتياجات
- تحسين الأداء والاستقرار

#### 6. الامتثال القانوني
- الاحتفاظ بالسجلات للفترات القانونية المطلوبة
- الامتثال لقوانين العمل والضرائب
- دعم عمليات التدقيق والمراجعة

## مشاركة البيانات

### لا نشارك بياناتك مع:
- ❌ أي طرف ثالث لأغراض تسويقية
- ❌ شركات الإعلانات أو الشبكات الاجتماعية
- ❌ موردي البيانات الخارجيين
- ❌ أي طرف بدون موافقتك الصريحة

### قد نشارك البيانات مع:

#### 1. المشرفين والإدارة
- **البيانات**: بيانات الحضور، الموقع أثناء العمل، حالة المهام
- **الغرض**: الإدارة اليومية، تقييم الأداء، التنسيق
- **القيود**: فقط البيانات الضرورية لأداء وظائفهم

#### 2. إدارة الموارد البشرية
- **البيانات**: البيانات الشخصية، سجلات الحضور، بيانات الرواتب
- **الغرض**: إدارة شؤون الموظفين، الرواتب، الإجازات
- **القيود**: الوصول المصرح به فقط

#### 3. مزودي الخدمات التقنية
- **البيانات**: بيانات مشفرة، معلومات تقنية محدودة
- **الأمثلة**: 
  - Render.com (استضافة السيرفرات)
  - Geoapify (خدمات الخرائط والعناوين)
- **الغرض**: تشغيل التطبيق وتوفير الخدمات
- **القيود**: فقط البيانات الضرورية، مشفرة، مع عقود حماية

#### 4. السلطات القانونية
- **البيانات**: عند الطلب القانوني المسبق
- **الغرض**: الامتثال للقوانين والأنظمة
- **القيود**: فقط عند الطلب القانوني المسبق، الحد الأدنى المطلوب

## أمن البيانات

### التدابير الأمنية المطبقة:

#### 1. التشفير (Encryption)
- **أثناء النقل**: جميع البيانات تُشفّر باستخدام TLS 1.3/HTTPS
- **أثناء التخزين**: البيانات الحساسة تُشفّر في قاعدة البيانات
- **المفاتيح**: تدوير المفاتيح الدوري وإدارة آمنة

#### 2. التحكم في الوصول (Access Control)
- **المصادقة**: نظام مصادقة قوي بـ JWT tokens
- **التفويض**: التحكم في الصلاحيات حسب الدور (Admin, Worker, Client)
- **المراجعة**: سجلات كاملة لجميع عمليات الوصول
- **الحدود**: قيود على محاولات تسجيل الدخول

#### 3. الحماية من الهجمات
- **DDoS Protection**: حماية من هجمات الحرمان من الخدمة
- **SQL Injection Prevention**: حماية من حقن SQL
- **XSS Protection**: حماية من هجمات XSS
- **Rate Limiting**: حدود على الطلبات لمنع الاستغلال

#### 4. المراقبة والاكتشاف
- **المراقبة المستمرة**: مراقبة 24/7 للأنشطة المشبوهة
- **التسجيل**: سجلات مفصلة لجميع العمليات
- **التنبيهات**: تنبيهات فورية للأنشطة غير المعتادة
- **التدقيق**: تدقيق دوري للأنظمة والإجراءات

#### 5. النسخ الاحتياطي (Backup)
- **النسخ الاحتياطي**: نسخ احتياطي يومي تلقائي
- **التشفير**: النسخ الاحتياطية مشفرة
- **التخزين**: تخزين في مواقع متعددة
- **الاختبار**: اختبار دوري لعملية الاستعادة

#### 6. التحديثات
- **التحديثات الأمنية**: تحديثات فورية للثغرات
- **التصحيح**: إصلاح دوري للمشاكل المكتشفة
- **الاختبار**: اختبار شامل قبل التحديثات

## الاحتفاظ بالبيانات (Data Retention)

### فترات الاحتفاظ:

#### 1. بيانات الموقع
- **الفترة**: 6 أشهر (يمكن تمديدها حسب سياسة الشركة)
- **السبب**: تتبع الأداء، تحسين السلامة، الامتثال القانوني
- **الحذف**: حذف آمن بعد انتهاء الفترة

#### 2. سجلات الحضور
- **الفترة**: 5 سنوات
- **السبب**: متطلبات ضريبية وقانونية
- **الحذف**: حذف آمن بعد انتهاء الفترة القانونية

#### 3. البيانات الشخصية
- **الفترة**: طوال فترة التوظيف + سنتين بعد انتهائه
- **السبب**: متطلبات قانونية، مراجع التوظيف
- **الحذف**: حذف آمن بعد انتهاء الفترة

#### 4. الصور والملفات
- **الفترة**: طوال فترة الاحتفاظ بالحساب
- **السبب**: التوثيق، التدقيق، الأرشفة
- **الحذف**: حذف آمن عند حذف الحساب

#### 5. بيانات الاستخدام
- **الفترة**: 12 شهر
- **السبب**: تحليل الأداء، تحسين الخدمة
- **الحذف**: حذف آمن بعد انتهاء الفترة

### عند حذف الحساب أو إنهاء التوظيف:
- يتم حذف البيانات الشخصية فوراً
- تُحذف بيانات الموقع خلال 30 يوماً
- تُحذف الصور والملفات خلال 30 يوماً
- تُحفظ السجلات القانونية للفترة المطلوبة فقط
- يتم إرسال تأكيد بالحذف

## حقوقك (Your Rights)

### حقوق المستخدم بموجب GDPR وCCPA:

#### 1. حق الوصول (Right to Access)
- يمكنك طلب نسخة كاملة من بياناتك في أي وقت
- سنقدمها بصيغة مقروءة آلياً (JSON/CSV)
- سنقدمها خلال 30 يوماً من الطلب

#### 2. حق التصحيح (Right to Rectification)
- يمكنك تصحيح بياناتك غير الدقيقة أو غير المكتملة
- سنقوم بالتحديث فوراً
- سنرسل تأكيداً بالتحديث

#### 3. حق الحذف (Right to Erasure / Right to be Forgotten)
- يمكنك طلب حذف حسابك وبياناتك
- سنحذف البيانات الشخصية فوراً
- سنحذف البيانات الأخرى حسب سياسة الاحتفاظ
- سنرسل تأكيداً بالحذف

#### 4. حق الإلغاء (Right to Withdraw Consent)
- يمكنك سحب موافقتك في أي وقت
- سيؤثر ذلك على الخدمات التي تتطلب الموافقة
- سنشرح التأثيرات قبل الإلغاء

#### 5. حق التقييد (Right to Restrict Processing)
- يمكنك طلب تقييد معالجة بياناتك
- سنحتفظ بالبيانات لكن لن نعالجها
- سنشرح التأثيرات على الخدمات

#### 6. حق النقل (Right to Data Portability)
- يمكنك طلب نقل بياناتك إلى مزود آخر
- سنقدمها بصيغة منظمة وقابلة للنقل
- سنساعد في عملية النقل

#### 7. حق الاعتراض (Right to Object)
- يمكنك الاعتراض على معالجة بياناتك
- سنوقف المعالجة ما لم تكن ضرورية قانونياً
- سنشرح البدائل المتاحة

#### 8. حق الشكوى (Right to Complain)
- يمكنك تقديم شكوى للسلطات المختصة
- سنقدم معلومات حول كيفية القيام بذلك
- لن ننتقم منك لتقديم شكوى

### كيفية ممارسة حقوقك:

#### للاتصال بنا:
- **البريد الإلكتروني**: privacy@worktrack.com
- **الهاتف**: [أضف رقم هاتف الشركة]
- **العنوان**: [أضف عنوان الشركة]
- **النموذج**: نموذج طلب حقوق المستخدم على موقعنا

#### العملية:
1. تقديم الطلب عبر البريد الإلكتروني أو النموذج
2. التحقق من الهوية (لحماية بياناتك)
3. معالجة الطلب خلال 30 يوماً
4. إرسال الرد والنتيجة

## البيانات التي نجمعها تلقائياً

### عبر التطبيق:
- **سجلات الاستخدام**: تسجيل الدخول، النشاط، الأوقات
- **بيانات التتبع**: أخطاء التطبيق، الأداء، الاستقرار
- **الإحصائيات**: استخدام الميزات، التردد، المدة
- **بيانات الموقع**: عند تفعيل تتبع الموقع

### عبر Cookies والتقنيات المشابهة:
- **Authentication tokens**: للبقاء مسجلاً للدخول
- **Preferences**: إعدادات اللغة والإشعارات
- **Cache**: لتحسين أداء التطبيق
- **Analytics**: لتحليل الاستخدام (بدون تعريف شخصي)

## الأطفال (Children)

### سياستنا:
- التطبيق مخصص للموظفين البالغين (18+)
- لا نجمع عمداً بيانات الأطفال
- إذا اكتشفنا بيانات أطفال، سنحذفها فوراً

### للآباء والأوصياء:
- إذا كان لديك طفل يستخدم التطبيق، اتصل بنا فوراً
- سنحذف بيانات الطفل فوراً
- سنقدم معلومات حول كيفية حدوث ذلك

## التغييرات على هذه السياسة

### متى نحدث السياسة:
- عند إضافة ميزات جديدة
- عند تغيير كيفية استخدام البيانات
- عند تغيير القوانين المطبقة
- عند تحسين إجراءات الأمان

### كيف نخبرك بالتغييرات:
- إشعار في التطبيق
- بريد إلكتروني
- تحديث في هذه الصفحة
- نشر في موقعنا

### تاريخ النفاذ:
- التغييرات الكبيرة: 30 يوماً قبل النفاذ
- التغييرات الصغيرة: فور النشر
- الاستمرار في استخدام التطبيق يعني الموافقة على التغييرات

## معلومات الاتصال

### للاستفسارات حول الخصوصية:
- **البريد الإلكتروني**: privacy@worktrack.com
- **الهاتف**: [أضف رقم هاتف الشركة]
- **العنوان**: [أضف عنوان الشركة الكامل]
- **ساعات العمل**: [أضف ساعات العمل]

### لاستفسارات الدعم الفني:
- **البريد الإلكتروني**: support@worktrack.com
- **الهاتف**: [أضف رقم هاتف الدعم]
- **الموقع**: [أضف موقع الويب]

### للشكاوى والاستئناف:
- **البريد الإلكتروني**: complaints@worktrack.com
- **الهاتف**: [أضف رقم هاتف الشكاوى]
- **العنوان**: [أضف عنوان الشركة]

## التشريعات المطبقة

### نلتزم بـ:
- **GDPR** (اللائحة العامة لحماية البيانات) - للمستخدمين في الاتحاد الأوروبي
- **CCPA** (قانون خصوصية المستهلكين في كاليفورنيا) - للمستخدمين في كاليفورنيا
- **قوانين حماية البيانات المحلية** - حسب الدولة/المنطقة
- **قوانين العمل والضرائب** - للسجلات المطلوبة

### ممثل البيانات:
- للاتحاد الأوروبي: [أضف معلومات الممثل]
- لكاليفورنيا: [أضف معلومات الممثل]

## التدقيق والمراجعة

### التدقيق الداخلي:
- مراجعة دورية للإجراءات الأمنية
- اختبار دوري للاستجابة للحوادث
- مراجعة سياسات الاحتفاظ بالبيانات

### التدقيق الخارجي:
- تدقيق دوري من طرف ثالث مستقل
- اختبار الاختراق المعتمد
- مراجعة الامتثال القانوني

## إجراءات اختراق البيانات

### في حالة اختراق البيانات:
- إشعار المستخدمين المتأثرين خلال 72 ساعة
- إشعار السلطات المختصة حسب القانون
- توضيح ما حدث وما البيانات المتأثرة
- شرح الخطوات المتخذة لحماية المستخدمين
- تقديم دعم ومساعدة للمستخدمين المتأثرين

## الروابط ذات الصلة

- **شروط الخدمة**: [رابط شروط الخدمة]
- **سياسة ملفات تعريف الارتباط**: [رابط سياسة Cookies]
- **إعدادات الخصوصية في التطبيق**: [شرح كيفية الوصول]

## آخر تحديث
**التاريخ**: 5 أغسطس 2026
**الإصدار**: 1.0

---

## أسئلة شائعة (FAQ)

### س1: هل يمكنني رفض تتبع الموقع؟
**ج**: نعم، لكن قد يؤثر ذلك على قدرتك على استخدام بعض ميزات التطبيق مثل التحقق من الحضور. تحدث مع مشرفك لمعرفة البدائل.

### س2: هل يمكنني طلب نسخة من بياناتي؟
**ج**: نعم، يمكنك طلب نسخة كاملة من بياناتك في أي وقت عبر البريد الإلكتروني privacy@worktrack.com

### س3: ماذا يحدث لبياناتي عند مغادرة الشركة؟
**ج**: يتم حذف بياناتك الشخصية فوراً، مع الاحتفاظ بالسجلات القانونية للفترة المطلوبة فقط.

### س4: هل تشاركون بياناتي مع أطراف ثالثة؟
**ج**: لا، لا نشارك بياناتك مع أي طرف ثالث لأغراض تسويقية. قد نشاركها مع مزودي الخدمات التقنية فقط عند الضرورة.

### س5: كيف تحمون بياناتي؟
**ج**: نستخدم تشفيراً قوياً، تحكماً صارماً في الوصول، مراقبة مستمرة، ونسخاً احتياطياً مشفراً.

### س6: كم مدة الاحتفاظ ببياناتي؟
**ج**: تختلف حسب نوع البيانات: بيانات الموقع (6 أشهر)، سجلات الحضور (5 سنوات)، البيانات الشخصية (حتى نهاية التوظيف + سنتين).

### س7: هل يمكنني حذف حسابي؟
**ج**: نعم، يمكنك طلب حذف حسابك في أي وقت. سنحذف بياناتك الشخصية فوراً مع الاحتفاظ بالسجلات القانونية فقط.

---

## ملاحظات للمطور:

### عند استخدام هذا المستند:
1. **استبدل [النص بين الأقواس]** بمعلوماتك الحقيقية
2. **أضف معلومات الاتصال الحقيقية** للشركة
3. **راجع القوانين المحلية** في منطقتك
4. **حدّث التاريخ** عند أي تغييرات
5. **ارفع السياسة** على موقع ويب أو GitHub Pages

### النشر على GitHub Pages:
```bash
# إنشاء repository للسياسة
git init
git add PRIVACY_POLICY.md
git commit -m "Add Privacy Policy"
git branch -M main
git remote add origin https://github.com/yourusername/worktrack-privacy-policy.git
git push -u origin main

# تفعيل GitHub Pages من Settings
```

### إضافة رابط السياسة في التطبيق:
```javascript
// في صفحة الإعدادات أو التسجيل
<a href="https://yourusername.github.io/worktrack-privacy-policy/" target="_blank">
  سياسة الخصوصية
</a>
```
```

---

## 📄 PRIVACY_POLICY_TEMPLATE.md

```markdown
# سياسة الخصوصية لتطبيق WorkTrack

## مقدمة
تحترم WorkTrack خصوصيتك وتلتزم بحماية بياناتك الشخصية. توضح هذه السياسة كيفية جمع واستخدام وحماية بياناتك عند استخدام تطبيقات WorkTrack.

## البيانات التي نجمعها

### 1. بيانات الموقع
- **الموقع الجغرافي**: نجمع بيانات الموقع الجغرافي الحالي والمسار التاريخي للموظفين أثناء ساعات العمل
- **الغرض**: تتبع الحضور والانصراف، التحقق من وجود الموظف في موقع العمل، تحسين سلامة الموظفين
- **الأساس القانوني**: ضروري لأداء العقد (اتفاقية العمل)

### 2. البيانات الشخصية
- **معلومات الحساب**: الاسم الكامل، البريد الإلكتروني، رقم الهاتف، صورة الملف الشخصي
- **الغرض**: إدارة الحساب، التواصل، التعريف في التطبيق
- **الأساس القانوني**: موافقة المستخدم

### 3. بيانات الحضور والعمل
- **سجلات الحضور**: أوقات الدخول والخروج، إجمالي ساعات العمل، نقاط العمل
- **البيانات الوظيفية**: المهام المخصصة، حالة إكمال المهام، الملاحظات
- **الغرض**: إدارة شؤون الموظفين، حساب الرواتب، تتبع الإنتاجية
- **الأساس القانوني**: ضروري لأداء العقد

### 4. الصور والملفات
- **الصور الشخصية**: صور الملفات الشخصية، صور التحقق من الحضور
- **الملفات المرفقة**: المستندات المتعلقة بالعمل
- **الغرض**: التحقق من الهوية، توثيق الحضور، إدارة المهام
- **الأساس القانوني**: موافقة المستخدم

### 5. بيانات الجهاز
- **معلومات الجهاز**: نوع الجهاز، نظام التشغيل، معرف الجهاز الفريد
- **الغرض**: تحسين الأداء، دعم فني، تحليلات الاستخدام
- **الأساس القانوني**: مصلحة مشروعة

## كيفية استخدام بياناتك

### الأغراض الرئيسية:
1. **إدارة شؤون الموظفين**: تتبع الحضور، حساب ساعات العمل، إدارة المهام
2. **التحقق من الموقع**: التأكد من وجود الموظف في مواقع العمل المحددة
3. **إدارة الأداء**: تقييم الإنتاجية، تحسين توزيع الموارد
4. **الأمان**: حماية سلامة الموظفين من خلال تتبع الموقع
5. **التواصل**: إرسال إشعارات، تحديثات المهام، تنبيهات هامة
6. **تحسين التطبيق**: تحليلات الاستخدام، تحسين الأداء

## مشاركة البيانات

### لا نشارك بياناتك مع:
- أي طرف ثالث لأغراض تسويقية
- شركات الإعلانات
- شبكات اجتماعية

### قد نشارك البيانات مع:
- **المشرف المباشر**: للإدارة اليومية للموظفين
- **إدارة الموارد البشرية**: لإدارة شؤون الموظفين والرواتب
- **مزودي الخدمات**: فقط عند الضرورة لتشغيل التطبيق (مثل Render.com للاستضافة)
- **السلطات القانونية**: عند الطلب القانوني المسبق

## أمن البيانات

### التدابير الأمنية:
- **التشفير**: جميع البيانات تُشفّر أثناء النقل (HTTPS) والتخزين
- **التحكم في الوصول**: الوصول المحدود والمصرح به فقط
- **المراقبة**: مراقبة مستمرة للأنشطة المشبوهة
- **النسخ الاحتياطي**: نسخ احتياطي منتظم للبيانات
- **التحديثات**: تحديثات أمنية دورية

## الاحتفاظ بالبيانات

### فترات الاحتفاظ:
- **بيانات الموقع**: 6 أشهر (يمكن تمديدها حسب سياسة الشركة)
- **سجلات الحضور**: 5 سنوات (لأغراض ضريبية وقانونية)
- **البيانات الشخصية**: طوال فترة التوظيف + سنتين بعد انتهائه
- **الصور والملفات**: طوال فترة الاحتفاظ بالحساب

### عند حذف الحساب:
- يتم حذف البيانات الشخصية فوراً
- تُحذف بيانات الموقع خلال 30 يوماً
- تُحذف الصور والملفات خلال 30 يوماً
- تُحفظ السجلات القانونية للفترة المطلوبة

## حقوقك

### حقوق المستخدم:
1. **الوصول**: يمكنك طلب نسخة من بياناتك في أي وقت
2. **التصحيح**: يمكنك تصحيح بياناتك غير الدقيقة
3. **الحذف**: يمكنك طلب حذف حسابك وبياناتك
4. **الإلغاء**: يمكنك سحب موافقتك في أي وقت
5. **الشكوى**: يمكنك تقديم شكوى للسلطات المختصة

### كيفية ممارسة حقوقك:
- الاتصال بنا: privacy@worktrack.com
- أو عبر: support@worktrack.com
- سنرد خلال 30 يوماً

## البيانات التي نجمعها تلقائياً

### عبر التطبيق:
- **سجلات الاستخدام**: تسجيل الدخول، النشاط، الأوقات
- **بيانات التتبع**: أخطاء التطبيق، الأداء، الاستقرار
- **الإحصائيات**: استخدام الميزات، التردد، المدة

### عبر Cookies والتقنيات المشابهة:
- **Authentication tokens**: للبقاء مسجلاً للدخول
- **Preferences**: إعدادات اللغة والإشعارات
- **Cache**: لتحسين أداء التطبيق

## الأطفال

التطبيق مخصص للموظفين البالغين (18+). لا نجمع عمداً بيانات الأطفال.

## التغييرات على هذه السياسة

قد نحدث هذه السياسة periodically. سنخبرك بالتغييرات الكبيرة عبر:
- إشعار في التطبيق
- بريد إلكتروني
- تحديث في هذه الصفحة

## معلومات الاتصال

### للاستفسارات حول الخصوصية:
- **البريد الإلكتروني**: privacy@worktrack.com
- **الهاتف**: [رقم هاتف الشركة]
- **العنوان**: [عنوان الشركة]

### لاستفسارات الدعم الفني:
- **البريد الإلكتروني**: support@worktrack.com
- **الهاتف**: [رقم هاتف الدعم]

## التشريعات المطبقة

نلتزم بـ:
- GDPR (للمستخدمين في الاتحاد الأوروبي)
- CCPA (للمستخدمين في كاليفورنيا)
- قوانين حماية البيانات المحلية الأخرى

## آخر تحديث
[التاريخ الحالي]

---

## ملاحظات للمطور:

### عند استخدام هذا القالب:
1. **تخصيص المعلومات**: استبدل [النص بين الأقواس] بمعلوماتك الحقيقية
2. **إضافة عناوين الشركة**: أضف معلومات الاتصال الحقيقية
3. **مراجعة القوانين المحلية**: تأكد من الامتثال للقوانين في منطقتك
4. **تحديث التاريخ**: حدّث تاريخ آخر تحديث
5. **رفع السياسة**: انشرها على موقع ويب أو GitHub Pages

### النشر على GitHub Pages:
```bash
# إنشاء repository للسياسة
git init
git add PRIVACY_POLICY.md
git commit -m "Add Privacy Policy"
git branch -M main
git remote add origin https://github.com/yourusername/privacy-policy.git
git push -u origin main

# تفعيل GitHub Pages
# Settings > Pages > Select main branch
```

### إضافة رابط السياسة في التطبيق:
```javascript
// في صفحة الإعدادات أو التسجيل
<a href="https://yourusername.github.io/privacy-policy/" target="_blank">
  سياسة الخصوصية
</a>
```
```

---

## 📄 SECURITY_DEPLOYMENT_GUIDE.md

```markdown
# 🚀 دليل تطبيق الإصلاحات الأمنية - WorkTrack

## 📋 نظرة عامة

هذا الدليل يشرح كيفية تطبيق الإصلاحات الأمنية التي تم رفعها على GitHub.

## 🔒 الإصلاحات المنفذة

### 1. JWT Storage Security (حرج)
- ✅ تحويل من localStorage إلى HttpOnly cookies
- ✅ إضافة logout endpoint
- ✅ تحديث جميع الواجهات الأمامية

### 2. Security Headers (حرج)
- ✅ إضافة CSP, X-Frame-Options, X-XSS-Protection
- ✅ تفعيل رؤوس أمان شاملة

### 3. Rate Limiting Enhancement (عالي)
- ✅ تقليل الحد العام إلى 150 طلب/دقيقة
- ✅ حد صارم للمصادقة (10 محاولات)
- ✅ حظر تلقائي للـ IP المشبوهة

### 4. Security Logging System (عالي)
- ✅ نظام تسجيل شامل للأحداث الأمنية
- ✅ كشف الأنماط المشبوهة
- ✅ endpoint جديد لمراقبة السجلات

## 🚀 خطوات التطبيق

### الخطوة 1: تطبيق قاعدة البيانات

**من الطرفية (Terminal):**

```bash
cd /home/dev-bit/project/worktrack/backend

# باستخدام psql مباشرة
psql "postgresql://postgres.ghmdhpikqzilosqzxzfd:KiQeXPQ16Ie0PEzU@aws-0-ap-southeast-1.pooler.supabase.com:5432/postgres?sslmode=require" -f internal/database/migrations/000020_add_security_logs.up.sql
```

**أو عبر Supabase Dashboard:**
1. افتح Supabase Dashboard
2. اذهب إلى SQL Editor
3. انسخ محتوى `backend/internal/database/migrations/000020_add_security_logs.up.sql`
4. نفذ الـ SQL

**التحقق من التطبيق:**
```sql
-- التحقق من إنشاء الجدول
SELECT table_name FROM information_schema.tables 
WHERE table_name = 'security_logs';

-- التحقق من الفهارس
SELECT indexname FROM pg_indexes 
WHERE tablename = 'security_logs';
```

### الخطوة 2: إعادة بناء الواجهة الخلفية

**إذا كان Go مثبتاً:**
```bash
cd /home/dev-bit/project/worktrack/backend
go build ./cmd/api
```

**على Render (الطريقة الموصى بها):**
1. سيتم بناء التطبيق تلقائياً عند الدفع
2. تأكد من أن `render.yaml` موجود في `backend/`
3. سيقوم Render ببناء Docker image تلقائياً

### الخطوة 3: إعادة بناء الواجهات الأمامية

**لوحة المدير:**
```bash
cd /home/dev-bit/project/worktrack/frontend-admin-dashboard
npm ci
npm run build
```

**تطبيق الموظف:**
```bash
cd /home/dev-bit/project/worktrack/frontend-worker-pwa
npm ci
npm run build
```

**بوابة العميل:**
```bash
cd /home/dev-bit/project/worktrack/frontend-client-portal
npm ci
npm run build
```

### الخطوة 4: إعادة النشر

**على Render:**
1. تم دفع التغييرات بالفعل إلى GitHub
2. سيقوم Render ببناء الخدمات تلقائياً
3. راقب بناء الخدمات في Render Dashboard

**تأكيد النشر:**
- تحقق من logs في Render Dashboard
- تأكد من أن جميع الخدمات تعمل
- اختبر تسجيل الدخول

## 🧪 الاختبار

### اختبار 1: تسجيل الدخول
```bash
# اختبر تسجيل دخول المدير
curl -X POST https://your-api.onrender.com/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"password"}' \
  -c cookies.txt

# التحقق من تعيين cookie
cat cookies.txt
```

### اختبار 2: رؤوس الأمان
```bash
# التحقق من رؤوس الأمان
curl -I https://your-api.onrender.com/health

# يجب أن ترى:
# Content-Security-Policy
# X-Frame-Options: DENY
# X-Content-Type-Options: nosniff
# X-XSS-Protection
```

### اختبار 3: Rate Limiting
```bash
# اختبر الحد (ستحصل على 429 بعد عدة محاولات)
for i in {1..20}; do
  curl -X POST https://your-api.onrender.com/api/v1/auth/login \
    -H "Content-Type: application/json" \
    -d '{"email":"test@example.com","password":"wrong"}'
done
```

### اختبار 4: Security Logs
```bash
# اختبر endpoint سجلات الأمان (يتطلب توكن مدير)
curl -X GET https://your-api.onrender.com/api/v1/security/logs \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## ⚠️ ملاحظات هامة

### HTTPS Required
- الكوكيات تعمل بـ `Secure: true`
- يتطلب HTTPS في الإنتاج
- في التطوير المحلي، قد تحتاج لتعديل هذا

### التوافقية
- النظام يدعم cookies و Authorization header
- عملاء API القديمون سيستمر في العمل
- يُنصح بالترقية إلى cookies

### التراجع عن التغييرات
إذا واجهت مشاكل، يمكنك التراجع:

```bash
# التراجع عن migration
psql "postgresql://..." -f internal/database/migrations/000020_add_security_logs.down.sql

# التراجع عن الكود
git revert 30c986c
git push origin main
```

## 📊 المراقبة

بعد النشر، راقب:

1. **Logs في Render Dashboard**
   - تحقق من أخطاء البناء
   - راقب سجلات الأمان

2. **Security Logs**
   - زر `/api/v1/security/logs` (للمدراء)
   - راقب الأنماط المشبوهة

3. **Performance**
   - تأكد من أن Rate Limiting لا يؤثر على المستخدمين الشرعيين
   - راقب أوقات الاستجابة

## 🆘 استكشاف الأخطاء

### مشكلة: "Cookie not being set"
**الحل:**
- تأكد من استخدام HTTPS
- تحقق من إعدادات CORS
- تأكد من `withCredentials: true` في الواجهة الأمامية

### مشكلة: "Migration failed"
**الحل:**
- تحقق من صلاحيات قاعدة البيانات
- تأكد من أن الجدول غير موجود مسبقاً
- راقب أخطاء SQL

### مشكلة: "Build failed"
**الحل:**
- تحقق من logs في Render
- تأكد من أن جميع التبعيات موجودة
- راقب أخطاء Go compilation

## 📞 الدعم

إذا واجهت مشاكل:
1. راقب logs في Render Dashboard
2. تحقق من security logs في قاعدة البيانات
3. راجع هذا الدليل
4. افتح issue على GitHub

---

**تم إنشاء هذا الدليل بواسطة Devin - مساعد الذكاء الاصطناعي**
**التاريخ:** 27 يوليو 2026
```

---

## 📄 SECURITY_MONITORING_GUIDE.md

```markdown
# 🔒 دليل مراقبة الأمان - WorkTrack

## طرق مراقبة Security Logs

### الطريقة 1: Supabase Dashboard (الأفضل والأسهل)

#### الخطوات:
1. افتح https://supabase.com/dashboard
2. اختر مشروع WorkTrack
3. من القائمة الجانبية، اختر "Table Editor"
4. اختر جدول `security_logs`

#### المزايا:
- ✅ واجهة رسومية سهلة الاستخدام
- ✅ تصفية وترتيب البيانات
- ✅ لا تتطلب تعديلات على الكود
- ✅ تحديث في الوقت الفعلي

#### استعلامات SQL مفيدة:
```sql
-- المحاولات الفاشلة الأخيرة
SELECT * FROM security_logs 
WHERE success = false 
ORDER BY created_at DESC 
LIMIT 20;

-- الأنماط المشبوهة (محاولات كثيرة من نفس IP)
SELECT ip, COUNT(*) as attempts 
FROM security_logs 
WHERE success = false 
AND created_at > NOW() - INTERVAL '1 hour'
GROUP BY ip 
HAVING COUNT(*) > 5;

-- محاولات تسجيل دخول لنفس بريد إلكتروني
SELECT email, COUNT(*) as attempts, 
       MAX(created_at) as last_attempt
FROM security_logs 
WHERE event_type = 'auth_attempt' 
AND success = false 
AND created_at > NOW() - INTERVAL '24 hours'
GROUP BY email 
HAVING COUNT(*) > 3;

-- إحصائيات شاملة
SELECT 
    event_type,
    success,
    COUNT(*) as count,
    COUNT(DISTINCT ip) as unique_ips,
    COUNT(DISTINCT email) as unique_emails
FROM security_logs 
WHERE created_at > NOW() - INTERVAL '24 hours'
GROUP BY event_type, success
ORDER BY count DESC;
```

### الطريقة 2: API Endpoint

#### الاستخدام عبر curl:
```bash
# تسجيل الدخول
curl -X POST https://worktrack-v2-api.onrender.com/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@devpro.com","password":"devproadmin"}' \
  -c cookies.txt

# الحصول على السجلات (آخر 50)
curl -X GET https://worktrack-v2-api.onrender.com/api/v1/security/logs \
  -b cookies.txt
```

#### الاستخدام عبر JavaScript:
```javascript
// في لوحة المدير
const response = await fetch('/api/v1/security/logs', {
  credentials: 'include' // لإرسال cookies
});
const logs = await response.json();
console.log(logs);
```

### الطريقة 3: إضافة واجهة مراقبة في لوحة المدير

#### مثال Vue.js Component:
```vue
<template>
  <div class="security-logs">
    <h2>سجلات الأمان</h2>
    
    <div class="filters">
      <select v-model="filterType">
        <option value="">كل الأنواع</option>
        <option value="auth_attempt">محاولات المصادقة</option>
        <option value="suspicious_activity">نشاط مشبوه</option>
      </select>
      
      <select v-model="filterSuccess">
        <option value="">الكل</option>
        <option value="true">ناجح</option>
        <option value="false">فاشل</option>
      </select>
      
      <button @click="loadLogs">تحديث</button>
    </div>
    
    <div class="logs-table">
      <table>
        <thead>
          <tr>
            <th>الوقت</th>
            <th>النوع</th>
            <th>البريد</th>
            <th>IP</th>
            <th>الحالة</th>
            <th>السبب</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="log in filteredLogs" :key="log.id">
            <td>{{ formatDate(log.created_at) }}</td>
            <td>{{ log.event_type }}</td>
            <td>{{ log.email || '-' }}</td>
            <td>{{ log.ip }}</td>
            <td>{{ log.success ? '✅' : '❌' }}</td>
            <td>{{ log.reason || '-' }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script>
import api from '@/services/api'

export default {
  name: 'SecurityLogs',
  data() {
    return {
      logs: [],
      filterType: '',
      filterSuccess: ''
    }
  },
  computed: {
    filteredLogs() {
      return this.logs.filter(log => {
        if (this.filterType && log.event_type !== this.filterType) return false
        if (this.filterSuccess && String(log.success) !== this.filterSuccess) return false
        return true
      })
    }
  },
  methods: {
    async loadLogs() {
      try {
        const response = await api.get('/security/logs')
        this.logs = response.data
      } catch (error) {
        console.error('فشل تحميل السجلات:', error)
      }
    },
    formatDate(date) {
      return new Date(date).toLocaleString('ar-SA')
    }
  },
  mounted() {
    this.loadLogs()
    // تحديث كل 30 ثانية
    setInterval(this.loadLogs, 30000)
  }
}
</script>

<style scoped>
.security-logs {
  padding: 20px;
}

.filters {
  margin-bottom: 20px;
  display: flex;
  gap: 10px;
}

.logs-table table {
  width: 100%;
  border-collapse: collapse;
}

.logs-table th, .logs-table td {
  padding: 10px;
  border: 1px solid #ddd;
  text-align: right;
}

.logs-table th {
  background-color: #f5f5f5;
}
</style>
```

### الطريقة 4: إشعارات البريد الإلكتروني

يمكنك إضافة نظام إشعارات للمحاولات المشبوهة:

#### في Go Backend:
```go
// في auth_handler.go
func (h *AuthHandler) Login(c *gin.Context) {
    // ... كود موجود ...
    
    if isSuspicious {
        // إرسال إشعار بريد إلكتروني
        h.sendSecurityAlert(ip, email, reason)
        
        c.JSON(http.StatusTooManyRequests, gin.H{
            "error": "محاولات مشبوهة - تم حظر مؤقت",
        })
        return
    }
}

func (h *AuthHandler) sendSecurityAlert(ip, email, reason string) {
    // إرسال إشعار للمديرين
    // يمكن استخدام خدمة بريد إلكتروني مثل SendGrid
}
```

## 📊 استعلامات مراقبة متقدمة

### تحليل يومي:
```sql
SELECT 
    DATE(created_at) as date,
    COUNT(*) as total_attempts,
    COUNT(CASE WHEN success = true THEN 1 END) as successful_logins,
    COUNT(CASE WHEN success = false THEN 1 END) as failed_logins,
    COUNT(DISTINCT ip) as unique_ips
FROM security_logs 
WHERE created_at > NOW() - INTERVAL '7 days'
GROUP BY DATE(created_at)
ORDER BY date DESC;
```

### أكثر IP مشبوهة:
```sql
SELECT 
    ip,
    COUNT(*) as total_attempts,
    COUNT(CASE WHEN success = false THEN 1 END) as failed_attempts,
    MAX(created_at) as last_activity,
    array_agg(DISTINCT email) as attempted_emails
FROM security_logs 
WHERE created_at > NOW() - INTERVAL '24 hours'
GROUP BY ip 
HAVING COUNT(CASE WHEN success = false THEN 1 END) > 10
ORDER BY failed_attempts DESC;
```

### أنماط الوقت:
```sql
SELECT 
    EXTRACT(HOUR FROM created_at) as hour,
    COUNT(*) as attempts
FROM security_logs 
WHERE created_at > NOW() - INTERVAL '7 days'
GROUP BY hour
ORDER BY attempts DESC;
```

## 🚨 التنبيهات المقترحة

### قواعد التنبيه:
1. **أكثر من 5 محاولات فاشلة من نفس IP في 15 دقيقة**
2. **أكثر من 10 محاولات فاشلة لنفس بريد في ساعة**
3. **محاولات من IP من دول غير متوقعة**
4. **محاولات تسجيل دخول خارج أوقات العمل العادية**

### الإجراءات المقترحة:
1. **حظر تلقائي للـ IP المشبوهة**
2. **إرسال إشعار للمديرين**
3. **تطلب التحقق الثنائي (2FA)**
4. **تسجيل الحادثة في سجل تدقيق**

## 🛠️ أدوات إضافية

### لوحة مراقبة مخصصة:
يمكنك إنشاء لوحة مراقبة مخصصة باستخدام:
- Grafana + Supabase
- Metabase
- Retool
- AppSmith

### التكامل مع Slack/Discord:
```go
// إرسال إشعار Slack
func sendSlackAlert(message string) {
    webhookURL := "YOUR_SLACK_WEBHOOK"
    payload := map[string]string{"text": message}
    // إرسال إلى Slack
}
```

---

**التوصية:** ابدأ بـ Supabase Dashboard للمراقبة الأساسية، ثم أضف واجهة في لوحة المدير للمراقبة المستمرة.
```

---

## 📄 android/app/proguard-rules.pro

```prolog
# Add project specific ProGuard rules here.
# You can control the set of applied configuration files using the
# proguardFiles setting in build.gradle.
#
# For more details, see
#   http://developer.android.com/guide/developing/tools/proguard.html

# If your project uses WebView with JS, uncomment the following
# and specify the fully qualified class name to the JavaScript interface
# class:
#-keepclassmembers class fqcn.of.javascript.interface.for.webview {
#   public *;
#}

# Uncomment this to preserve the line number information for
# debugging stack traces.
#-keepattributes SourceFile,LineNumberTable

# If you keep the line number information, uncomment this to
# hide the original source file name.
#-renamesourcefileattribute SourceFile

```

---

## 📄 extracted_code/security_&_obfuscation.md

```markdown

```

---

## 📄 frontend-admin-dashboard/android/app/proguard-rules.pro

```prolog
# WorkTrack Admin App - ProGuard Configuration

# Keep Capacitor classes
-keep class com.getcapacitor.** { *; }
-dontwarn com.getcapacitor.**

# Keep Cordova classes
-keep class org.apache.cordova.** { *; }
-dontwarn org.apache.cordova.**

# Keep webview classes
-keep class android.webkit.** { *; }
-dontwarn android.webkit.**

# Keep Kotlin classes
-keep class kotlin.** { *; }
-dontwarn kotlin.**

# Keep Coroutines
-keepnames class kotlinx.coroutines.internal.MainDispatcherFactory {}
-keepnames class kotlinx.coroutines.CoroutineExceptionHandler {}
-keepclassmembers class kotlinx.coroutines.** {
    <fields>;
    <methods>;
}

# Keep JSON classes
-keepattributes Signature
-keepattributes *Annotation*
-dontwarn sun.misc.**
-keep class com.google.gson.** { *; }
-keep class * implements com.google.gson.TypeAdapter
-keep class * implements com.google.gson.TypeAdapterFactory
-keep class * implements com.google.gson.JsonSerializer
-keep class * implements com.google.gson.JsonDeserializer

# Keep OkHttp
-dontwarn okhttp3.**
-dontwarn okio.**
-keepattributes Signature
-keepattributes Exceptions
-keep class okhttp3.** { *; }
-keep interface okhttp3.** { *; }

# Keep Retrofit
-dontwarn retrofit2.**
-keep class retrofit2.** { *; }
-keepattributes Signature
-keepattributes Exceptions
-keepclasseswithmembers class * {
    @retrofit2.http.* <methods>;
}

# Keep model classes
-keep class com.worktrack.admin.model.** { *; }

# Keep location services
-keep class android.location.** { *; }
-keep class com.google.android.gms.location.** { *; }
-dontwarn com.google.android.gms.**

# Keep Firebase classes
-keep class com.google.firebase.** { *; }
-dontwarn com.google.firebase.**

# Keep Maps
-keep class com.google.android.gms.maps.** { *; }
-dontwarn com.google.android.gms.maps.**

# Remove logging
-assumenosideeffects class android.util.Log {
    public static *** d(...);
    public static *** v(...);
    public static *** i(...);
}

# Optimization
-optimizationpasses 5
-dontusemixedcaseclassnames
-dontskipnonpubliclibraryclasses
-dontpreverify
-verbose

```

---

## 📄 frontend-client-portal/android/app/proguard-rules.pro

```prolog
# WorkTrack Client App - ProGuard Configuration

# Keep Capacitor classes
-keep class com.getcapacitor.** { *; }
-dontwarn com.getcapacitor.**

# Keep Cordova classes
-keep class org.apache.cordova.** { *; }
-dontwarn org.apache.cordova.**

# Keep webview classes
-keep class android.webkit.** { *; }
-dontwarn android.webkit.**

# Keep Kotlin classes
-keep class kotlin.** { *; }
-dontwarn kotlin.**

# Keep Coroutines
-keepnames class kotlinx.coroutines.internal.MainDispatcherFactory {}
-keepnames class kotlinx.coroutines.CoroutineExceptionHandler {}
-keepclassmembers class kotlinx.coroutines.** {
    <fields>;
    <methods>;
}

# Keep JSON classes
-keepattributes Signature
-keepattributes *Annotation*
-dontwarn sun.misc.**
-keep class com.google.gson.** { *; }
-keep class * implements com.google.gson.TypeAdapter
-keep class * implements com.google.gson.TypeAdapterFactory
-keep class * implements com.google.gson.JsonSerializer
-keep class * implements com.google.gson.JsonDeserializer

# Keep OkHttp
-dontwarn okhttp3.**
-dontwarn okio.**
-keepattributes Signature
-keepattributes Exceptions
-keep class okhttp3.** { *; }
-keep interface okhttp3.** { *; }

# Keep Retrofit
-dontwarn retrofit2.**
-keep class retrofit2.** { *; }
-keepattributes Signature
-keepattributes Exceptions
-keepclasseswithmembers class * {
    @retrofit2.http.* <methods>;
}

# Keep model classes
-keep class com.worktrack.client.model.** { *; }

# Keep location services
-keep class android.location.** { *; }
-keep class com.google.android.gms.location.** { *; }
-dontwarn com.google.android.gms.**

# Keep Firebase classes
-keep class com.google.firebase.** { *; }
-dontwarn com.google.firebase.**

# Keep Maps
-keep class com.google.android.gms.maps.** { *; }
-dontwarn com.google.android.gms.maps.**

# Remove logging
-assumenosideeffects class android.util.Log {
    public static *** d(...);
    public static *** v(...);
    public static *** i(...);
}

# Optimization
-optimizationpasses 5
-dontusemixedcaseclassnames
-dontskipnonpubliclibraryclasses
-dontpreverify
-verbose

```

---

## 📄 frontend-worker-pwa/android/app/proguard-rules.pro

```prolog
# Add project specific ProGuard rules here.
# You can control the set of applied configuration files using the
# proguardFiles setting in build.gradle.
#
# For more details, see
#   http://developer.android.com/guide/developing/tools/proguard.html

# If your project uses WebView with JS, uncomment the following
# and specify the fully qualified class name to the JavaScript interface
# class:
#-keepclassmembers class fqcn.of.javascript.interface.for.webview {
#   public *;
#}

# Uncomment this to preserve the line number information for
# debugging stack traces.
#-keepattributes SourceFile,LineNumberTable

# If you keep the line number information, uncomment this to
# hide the original source file name.
#-renamesourcefileattribute SourceFile

```

---

