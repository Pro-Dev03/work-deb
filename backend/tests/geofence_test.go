package tests

import (
	"testing"

	"worktrack/backend/internal/models"
	"worktrack/backend/internal/services"
)

// TestCheckWithinWorksite_Inside يتحقق أن نقطة قريبة جداً من موقع العمل تُعتبر "داخل النطاق"
func TestCheckWithinWorksite_Inside(t *testing.T) {
	worksite := models.Worksite{
		Latitude:     31.9539,
		Longitude:    35.9106,
		RadiusMeters: 100,
	}

	// نفس الإحداثيات تماماً => المسافة صفر => يجب أن تكون داخل النطاق دائماً
	result := services.CheckWithinWorksite(31.9539, 35.9106, worksite)

	if !result.IsWithinRange {
		t.Errorf("expected point to be within range, got distance=%.2f, radius=%d",
			result.DistanceMeters, result.AllowedRadius)
	}
}

// TestCheckWithinWorksite_Outside يتحقق أن نقطة بعيدة جداً تُعتبر "خارج النطاق" وتُرفض
func TestCheckWithinWorksite_Outside(t *testing.T) {
	worksite := models.Worksite{
		Latitude:     31.9539,
		Longitude:    35.9106,
		RadiusMeters: 100,
	}

	// إحداثيات مدينة أخرى بعيدة تماماً => يجب رفضها بوضوح
	result := services.CheckWithinWorksite(32.0853, 34.7818, worksite)

	if result.IsWithinRange {
		t.Errorf("expected point to be outside range, got distance=%.2f, radius=%d",
			result.DistanceMeters, result.AllowedRadius)
	}
}
