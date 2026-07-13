package utils

import "strings"

// IsValidEmail تحقق مبسط من صيغة البريد الإلكتروني (ليس تحققاً كاملاً بمعايير RFC،
// لكنه كافٍ لمنع الأخطاء الشائعة عند التسجيل)
func IsValidEmail(email string) bool {
	email = strings.TrimSpace(email)
	return strings.Contains(email, "@") && strings.Contains(email, ".") && len(email) >= 5
}

// IsValidCoordinates تتحقق أن الإحداثيات المرسلة من الموظف ضمن مجال منطقي على الخريطة
// (خط العرض بين -90 و90، وخط الطول بين -180 و180) قبل حتى حساب المسافة
func IsValidCoordinates(lat, lng float64) bool {
	return lat >= -90 && lat <= 90 && lng >= -180 && lng <= 180
}
