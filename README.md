# WorkTrack

منصة لإدارة الموظفين الميدانيين، الحضور، نقاط العمل، وطلبات الخدمة.

## مكوّنات المشروع

- `backend/`: واجهة API بلغة Go وPostgreSQL.
- `frontend-admin-dashboard/`: لوحة المدير (Vue/Vite).
- `frontend-worker-pwa/`: تطبيق الموظف (Vue/Vite/PWA).
- `frontend-client-portal/`: بوابة العميل (Vue/Vite).

## التشغيل محلياً

1. انسخ كل ملف `.env.example` إلى `.env` داخل المكوّن المقابل، وضع القيم المحلية فقط.
2. شغّل PostgreSQL واضبط `DATABASE_URL` في `backend/.env`.
3. شغّل الـ API:

   ```bash
   cd backend && go run ./cmd/api
   ```

4. شغّل أي واجهة:

   ```bash
   cd frontend-admin-dashboard && npm ci && npm run dev
   ```

## النشر على Render

ملف [`render.yaml`](render.yaml) ينشر الـ API من `backend/` كخدمة Docker ويضبط فحص الصحة على `/health`.

1. ارفع المشروع إلى GitHub، ثم اختر **New → Blueprint** في Render وحدد المستودع.
2. أثناء الإنشاء، عيّن `DATABASE_URL` و`ALLOWED_ORIGINS`. يجب أن تحتوي `ALLOWED_ORIGINS` على عناوين HTTPS الدقيقة للواجهات، مفصولة بفواصل، من دون `*`.
3. Render ينشئ `JWT_SECRET` آمناً تلقائياً. لا تضع أي سر في GitHub أو `render.yaml`.
4. ابنِ كل واجهة وانشرها كـ Static Site (على Render أو مزود آخر) مع المتغير `VITE_API_BASE_URL=https://<api-name>.onrender.com/api/v1`.
5. بعد معرفة روابط الواجهات النهائية، حدّث `ALLOWED_ORIGINS` في Render وأعد النشر.

قبل أول نشر، جهّز قاعدة PostgreSQL بمراجعة وتنفيذ migrations الموجودة في `backend/internal/database/migrations/`. لا تستخدم ملفات SQL المحلية المستثناة من Git لإعداد إنتاج لأنها قد تكون مدمّرة للبيانات.

> تنبيه: التخزين المحلي للملفات داخل Render مؤقت. رفع الصور في الكود الحالي يحتاج ربطاً فعلياً بتخزين دائم مثل Cloudflare R2 قبل الاعتماد عليه في الإنتاج.

## قبل الرفع إلى GitHub

`.env` والسجلات وملفات البناء والنسخ الاحتياطية مستثناة في `.gitignore`. لا ترفع كلمات مرور تجريبية أو قواعد بيانات فيها بيانات حقيقية. استخدم ملفات `.env.example` كقوالب فقط.
