package i18n

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// Detect يحدد لغة الرد المناسبة لطلب معيّن، بترتيب أولوية واضح:
//
//  1. باراميتر صريح في الرابط: /api/v1/tasks?lang=he
//     (مفيد لتجربة الـ API يدوياً أو من تطبيقات لا تتحكم بالهيدرز بسهولة)
//  2. هيدر مخصص يرسله الفرونت إند: X-Lang: he
//     (هذا ما تستخدمه واجهات WorkTrack الثلاث حسب لغة واجهة المستخدم المختارة)
//  3. هيدر المتصفح القياسي: Accept-Language
//  4. العربية كلغة افتراضية أخيرة إن لم يُرسَل أي شيء مما سبق
func Detect(c *gin.Context) Lang {
	if q := c.Query("lang"); q != "" {
		return Normalize(q)
	}

	if h := c.GetHeader("X-Lang"); h != "" {
		return Normalize(h)
	}

	if al := c.GetHeader("Accept-Language"); al != "" {
		// Accept-Language قد يأتي بصيغة معقدة مثل: "he-IL,he;q=0.9,en;q=0.8"
		// نأخذ فقط أول تفضيل ونحذف منه أي كود بلد (مثل IL) لنبقي على "he" فقط
		primary := strings.Split(al, ",")[0]
		primary = strings.Split(primary, ";")[0]
		primary = strings.Split(primary, "-")[0]
		return Normalize(primary)
	}

	return English
}
