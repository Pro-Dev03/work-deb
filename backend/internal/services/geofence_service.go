package services

import (
	"worktrack/backend/internal/models"
	"worktrack/backend/pkg/utils"
)

// GeofenceCheckResult نتيجة التحقق من موقع الموظف مقابل نقطة العمل
type GeofenceCheckResult struct {
	IsWithinRange  bool    `json:"is_within_range"`
	DistanceMeters float64 `json:"distance_meters"`
	AllowedRadius  int     `json:"allowed_radius_meters"`
}

// CheckWithinWorksite هي الدالة المسؤولة عن القرار الحاسم:
// هل موقع الموظف الحالي (lat, lng) يقع داخل نطاق نقطة العمل المحددة؟
//
// ⚠️ هذه الدالة تُستدعى دائماً من الـ Backend وليس فقط من الواجهة،
// لأن أي تحقق يتم فقط في المتصفح يمكن التلاعب به بسهولة (تعديل JS، إحداثيات مزيّفة...).
// السيرفر هو الحَكَم الوحيد الموثوق لقرار السماح بالتختيم أو رفضه، بغض النظر
// عن لغة الواجهة التي يستخدمها الموظف (عربي أو عبري أو إنجليزي) — القرار
// نفسه لا يتأثر باللغة إطلاقاً، فقط رسالة الرد للمستخدم هي ما تُترجَم لاحقاً.
func CheckWithinWorksite(userLat, userLng float64, worksite models.Worksite) GeofenceCheckResult {
	distance := utils.HaversineDistance(userLat, userLng, worksite.Latitude, worksite.Longitude)

	return GeofenceCheckResult{
		IsWithinRange:  distance <= float64(worksite.RadiusMeters),
		DistanceMeters: distance,
		AllowedRadius:  worksite.RadiusMeters,
	}
}
