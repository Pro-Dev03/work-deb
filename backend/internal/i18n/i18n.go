package i18n

import "strings"

// Lang يمثل رمز اللغة المدعومة في المشروع
type Lang string

const (
	Arabic  Lang = "ar"
	Hebrew  Lang = "he"
	English Lang = "en"
)

// messages هي "قاموس الترجمة" المركزي: مفتاح الرسالة -> نص لكل لغة
// أي رسالة جديدة يحتاجها المشروع تُضاف هنا مرة واحدة فقط بثلاث لغات
var messages = map[string]map[Lang]string{
	// ---------- تسجيل الدخول والحسابات ----------
	"err_missing_login_fields": {
		Arabic:  "الرجاء إدخال البريد الإلكتروني وكلمة المرور",
		Hebrew:  "נא להזין כתובת אימייל וסיסמה",
		English: "Please enter your email and password",
	},
	"err_invalid_credentials": {
		Arabic:  "البريد الإلكتروني أو كلمة المرور غير صحيحة",
		Hebrew:  "כתובת האימייל או הסיסמה שגויות",
		English: "Invalid email or password",
	},
	"err_session_create_failed": {
		Arabic:  "تعذر إنشاء جلسة الدخول",
		Hebrew:  "לא ניתן היה ליצור את הפעלת ההתחברות",
		English: "Could not create a login session",
	},
	"err_missing_fields": {
		Arabic:  "بيانات غير مكتملة",
		Hebrew:  "חסרים פרטים בבקשה",
		English: "Missing required fields",
	},
	"err_invalid_email": {
		Arabic:  "صيغة البريد الإلكتروني غير صحيحة",
		Hebrew:  "כתובת האימייל אינה תקינה",
		English: "Invalid email format",
	},
	"err_password_hash_failed": {
		Arabic:  "فشل تشفير كلمة المرور",
		Hebrew:  "הצפנת הסיסמה נכשלה",
		English: "Failed to hash the password",
	},
	"err_email_in_use": {
		Arabic:  "البريد الإلكتروني مستخدم مسبقاً",
		Hebrew:  "כתובת האימייל כבר קיימת במערכת",
		English: "This email is already registered",
	},
	"msg_account_created": {
		Arabic:  "تم إنشاء الحساب بنجاح",
		Hebrew:  "החשבון נוצר בהצלחה",
		English: "Account created successfully",
	},
	"err_user_not_found": {
		Arabic:  "المستخدم غير موجود",
		Hebrew:  "המשתמש לא נמצא",
		English: "User not found",
	},

	// ---------- الجلسة والصلاحيات ----------
	"err_please_login": {
		Arabic:  "الرجاء تسجيل الدخول",
		Hebrew:  "נא להתחבר תחילה",
		English: "Please log in",
	},
	"err_invalid_session": {
		Arabic:  "جلسة غير صالحة، الرجاء تسجيل الدخول من جديد",
		Hebrew:  "ההתחברות אינה תקפה, נא להתחבר מחדש",
		English: "Invalid session, please log in again",
	},
	"err_subscription_expired": {
		Arabic:  "اشتراكك انتهى أو تم إيقافه، الرجاء التواصل مع الدعم",
		Hebrew:  "המנוי שלך פג או בוטל, נא ליצור קשר עם התמיכה",
		English: "Your subscription has expired or been canceled, please contact support",
	},
	"err_forbidden_role": {
		Arabic:  "ليست لديك صلاحية للوصول لهذا الإجراء",
		Hebrew:  "אין לך הרשאה לבצע פעולה זו",
		English: "You don't have permission to perform this action",
	},
	"err_too_many_requests": {
		Arabic:  "طلبات كثيرة جداً، حاول لاحقاً",
		Hebrew:  "יותר מדי בקשות, נסה שוב מאוחר יותר",
		English: "Too many requests, please try again later",
	},

	// ---------- التختيم والنطاق الجغرافي (Geofence) ----------
	"err_invalid_request_data": {
		Arabic:  "بيانات الطلب غير صحيحة",
		Hebrew:  "נתוני הבקשה אינם תקינים",
		English: "Invalid request data",
	},
	"err_invalid_coordinates": {
		Arabic:  "إحداثيات الموقع غير صالحة",
		Hebrew:  "קואורדינטות המיקום אינן תקינות",
		English: "Invalid location coordinates",
	},
	"err_outside_geofence_checkin": {
		Arabic:  "لا يمكنك بدء المهمة، أنت خارج نطاق موقع العمل المسموح",
		Hebrew:  "לא ניתן להתחיל את המשימה, אתה מחוץ לטווח המותר של אתר העבודה",
		English: "You can't start the task, you are outside the allowed worksite range",
	},
	"err_outside_geofence_checkout": {
		Arabic:  "لا يمكنك إنهاء المهمة، أنت خارج نطاق موقع العمل المسموح",
		Hebrew:  "לא ניתן לסיים את המשימה, אתה מחוץ לטווח המותר של אתר העבודה",
		English: "You can't end the task, you are outside the allowed worksite range",
	},
	"notif_checkin_rejected_title": {
		Arabic:  "محاولة تختيم مرفوضة",
		Hebrew:  "ניסיון החתמה נדחה",
		English: "Check-in attempt rejected",
	},
	"notif_checkin_rejected_body": {
		Arabic:  "لقد حاولت بدء مهمة وأنت خارج نطاق موقع العمل المحدد",
		Hebrew:  "ניסית להתחיל משימה כשאתה מחוץ לטווח אתר העבודה שהוגדר",
		English: "You tried to start a task while outside the defined worksite range",
	},
	"msg_checkin_success": {
		Arabic:  "تم بدء المهمة بنجاح، أنت داخل النطاق المسموح",
		Hebrew:  "המשימה התחילה בהצלחה, אתה בטווח המותר",
		English: "Task started successfully, you are within the allowed range",
	},
	"msg_checkout_success": {
		Arabic:  "تم إنهاء المهمة بنجاح",
		Hebrew:  "המשימה הסתיימה בהצלחה",
		English: "Task ended successfully",
	},

	// ---------- عمليات عامة (نقاط عمل / مهام / عملاء / إشعارات / رفع ملفات) ----------
	"err_operation_failed": {
		Arabic:  "حدث خطأ ما، الرجاء المحاولة مرة أخرى",
		Hebrew:  "משהו השתבש, נא לנסות שוב",
		English: "Something went wrong, please try again",
	},
	"msg_worksite_created": {
		Arabic:  "تم إنشاء نقطة العمل بنجاح",
		Hebrew:  "אתר העבודה נוצר בהצלחה",
		English: "Worksite created successfully",
	},
	"msg_task_created": {
		Arabic:  "تم إنشاء المهمة بنجاح",
		Hebrew:  "המשימה נוצרה בהצלחה",
		English: "Task created successfully",
	},
	"err_no_photo": {
		Arabic:  "لم يتم إرفاق أي صورة",
		Hebrew:  "לא צורפה תמונה",
		English: "No photo was attached",
	},
	"err_file_open_failed": {
		Arabic:  "تعذر فتح الملف",
		Hebrew:  "לא ניתן היה לפתוח את הקובץ",
		English: "Could not open the file",
	},
	"msg_health_ok": {
		Arabic:  "الخدمة تعمل بشكل طبيعي",
		Hebrew:  "השירות פועל כרגיל",
		English: "Service is healthy",
	},
}

// T (Translate) تُرجع نص الرسالة المطابق للمفتاح key بلغة lang
// إذا لم توجد ترجمة لتلك اللغة تحديداً، تُستخدم العربية كخيار احتياطي
// وإذا لم يوجد المفتاح نفسه أصلاً، تُرجَع قيمة المفتاح كما هي (يساعد أثناء التطوير
// على اكتشاف أي رسالة نسينا ترجمتها بسهولة، بدل ظهور رسالة فارغة للمستخدم)
func T(lang Lang, key string) string {
	if translations, ok := messages[key]; ok {
		if text, ok := translations[lang]; ok {
			return text
		}
		if text, ok := translations[Arabic]; ok {
			return text
		}
	}
	return key
}

// Normalize تحوّل أي نص لغة وارد من العميل (query param أو هيدر) إلى Lang معروفة
// وتتجاهل أي قيمة غير مدعومة بإرجاع العربية كافتراضي آمن
func Normalize(raw string) Lang {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "he", "heb", "hebrew", "עברית":
		return Hebrew
	case "en", "eng", "english":
		return English
	case "ar", "ara", "arabic", "العربية":
		return Arabic
	default:
		return Arabic
	}
}
