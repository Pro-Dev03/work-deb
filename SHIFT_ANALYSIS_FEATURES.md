# ميزات تحليل الورديات - Shift Analysis Features

## نظرة عامة
تم إضافة ميزات جديدة لتحليل الورديات والتعامل مع الورديات عبر منتصف الليل والتمييز بين العمل الليلي والنهاري.

## الميزات المضافة

### 1. تقسيم الورديات عبر الأيام
- **الهدف**: التعامل مع الورديات التي تبدأ في يوم وتنتهي في اليوم التالي
- **التطبيق**: عندما يعمل الموظف بعد منتصف الليل، يتم تقسيم الوردية إلى يومين
- **الحقول المضافة**:
  - `spans_multiple_days`: Boolean للإشارة ما إذا كانت الوردية عبرت منتصف الليل
  - `day_one_date`: تاريخ اليوم الأول
  - `day_two_date`: تاريخ اليوم الثاني (إذا عبرت منتصف الليل)
  - `day_one_hours`: ساعات العمل في اليوم الأول
  - `day_two_hours`: ساعات العمل في اليوم الثاني

### 2. التمييز بين العمل الليلي والنهاري
- **الهدف**: حساب ساعات العمل الليلية والنهارية بشكل منفصل
- **التعريف**:
  - العمل الليلي: 10 مساءً - 6 صباحاً
  - العمل النهاري: 6 صباحاً - 10 مساءً
- **الحقول المضافة**:
  - `night_hours`: ساعات العمل الليلي
  - `day_hours`: ساعات العمل النهاري
  - `is_night_shift`: Boolean للإشارة ما إذا كانت الوردية ليلية بشكل أساسي (أكثر من 50% ليلي)

## التغييرات التقنية

### Backend (Go)

#### 1. نموذج Attendance
تم تحديث `backend/internal/models/attendance.go` بإضافة الحقول الجديدة:
```go
// حقول تقسيم الورديات عبر منتصف الليل
SpansMultipleDays bool       `json:"spans_multiple_days,omitempty"`
DayOneDate        *time.Time `json:"day_one_date,omitempty"`
DayTwoDate        *time.Time `json:"day_two_date,omitempty"`
DayOneHours       *float64   `json:"day_one_hours,omitempty"`
DayTwoHours       *float64   `json:"day_two_hours,omitempty"`

// حقول التمييز بين العمل الليلي والنهاري
NightHours   *float64 `json:"night_hours,omitempty"`
DayHours     *float64 `json:"day_hours,omitempty"`
IsNightShift bool     `json:"is_night_shift,omitempty"`
```

#### 2. دوال مساعدة للتواريخ
تم تحديث `backend/pkg/utils/datetime.go` بإضافة دوال جديدة:
- `CalculateDayNightHours()`: حساب ساعات العمل الليلية والنهارية
- `IsNightShift()`: تحديد ما إذا كانت الوردية ليلية
- `SplitShiftAcrossDays()`: تقسيم الوردية إذا عبرت منتصف الليل

#### 3. تحديث منطق CheckOut
تم تحديث `backend/internal/services/attendance_service.go`:
- دالة `CheckOut()`: إضافة حساب التقسيم والليل/النهار
- دالة `ForceCheckOut()`: نفس التحديثات للإنهاء الإجباري

#### 4. تحديث الاستعلامات
تم تحديث `backend/internal/handlers/attendance_handler.go`:
- `GetEmployeeAttendanceHistory()`: إضافة الحقول الجديدة للاستعلام
- `GetMyAttendanceHistory()`: نفس التحديثات للسجل الشخصي

#### 5. الترحيلات (Migrations)
تم إنشاء ملفات ترحيل جديدة:
- `backend/internal/database/migrations/000011_add_shift_analysis.up.sql`
- `backend/internal/database/migrations/000011_add_shift_analysis.down.sql`

#### 6. نظام الترحيل
تم إنشاء `backend/internal/database/migrate.go` لنظام ترحيل تلقائي:
- يتم تشغيل الترحيلات تلقائياً عند بدء الخادم
- يتم تتبع الترحيلات المطبقة في جدول `schema_migrations`

### Frontend

#### 1. تحديث ملفات الترجمة
تم تحديث ملفات اللغة في جميع الواجهات:
- `frontend-worker-pwa/src/locales/ar.json`
- `frontend-worker-pwa/src/locales/en.json`
- `frontend-worker-pwa/src/locales/he.json`
- `frontend-admin-dashboard/src/locales/ar.json`
- `frontend-admin-dashboard/src/locales/en.json`
- `frontend-admin-dashboard/src/locales/he.json`

#### 2. تحديث واجهة الموظف
تم تحديث `frontend-worker-pwa/src/views/AttendanceView.vue`:
- إضافة أعمدة جديدة لساعات العمل الليلية والنهارية في جدول سجل الحضور

#### 3. تحديث واجهة المدير
تم تحديث `frontend-admin-dashboard/src/views/EmployeesView.vue`:
- إضافة أعمدة جديدة في جدول سجل الحضور (سطح المكتب والجوال)
- تحديث وظيفة تصدير PDF لتشمل البيانات الجديدة

## الاختبارات

تم إنشاء اختبارات شاملة للدوال الجديدة في `backend/pkg/utils/datetime_test.go`:
- `TestCalculateDayNightHours`: اختبار حساب ساعات الليل/النهار
- `TestIsNightShift`: اختبار تحديد الورديات الليلية
- `TestSplitShiftAcrossDays`: اختبار تقسيم الورديات عبر الأيام

جميع الاختبارات تمر بنجاح.

## مثال على الاستخدام

### مثال وردية عبر منتصف الليل
- الموظف يبدأ العمل: 10:00 مساءً
- الموظف ينتهي العمل: 3:00 صباحاً
- **النتيجة**:
  - `spans_multiple_days`: true
  - `day_one_hours`: 2.0 (10 مساءً - 12 منتصف الليل)
  - `day_two_hours`: 3.0 (12 منتصف الليل - 3 صباحاً)
  - `night_hours`: 5.0 (كلها ساعات ليلية)
  - `day_hours`: 0.0
  - `is_night_shift`: true

### مثال وردية مختلطة
- الموظف يبدأ العمل: 5:00 مساءً
- الموظف ينتهي العمل: 1:00 صباحاً
- **النتيجة**:
  - `spans_multiple_days`: true
  - `day_one_hours`: 7.0 (5 مساءً - 12 منتصف الليل)
  - `day_two_hours`: 1.0 (12 منتصف الليل - 1 صباحاً)
  - `night_hours`: 3.0 (10 مساءً - 1 صباحاً)
  - `day_hours`: 5.0 (5 مساءً - 10 مساءً)
  - `is_night_shift`: false (أقل من 50% ليلي)

## التوافقية

- التغييرات متوافقة مع الإصدارات السابقة
- الحقول الجديدة اختيارية (`omitempty`)
- السجلات القديمة ستحتوي على قيم `null` لهذه الحقول
- نظام الترحيل يضمن تحديث قاعدة البيانات بشكل آمن

## المستقبل

يمكن إضافة ميزات مستقبلية بناءً على هذه الأساس:
- تقارير مفصلة للورديات الليلية
- حساب إضافات العمل الليلي
- رسوم بيانية لتوزيع ساعات العمل
- تنبيهات للورديات الطويلة
