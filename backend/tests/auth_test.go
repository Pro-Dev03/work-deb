package tests

import (
	"testing"

	"worktrack/backend/internal/services"
)

// TestPasswordHashAndCheck يتحقق أن كلمة المرور المشفّرة تُطابق نفسها عند المقارنة،
// ولا تُطابق كلمة مرور خاطئة
func TestPasswordHashAndCheck(t *testing.T) {
	auth := services.NewAuthService("test_secret")

	hash, err := auth.HashPassword("my-secure-password")
	if err != nil {
		t.Fatalf("unexpected error hashing password: %v", err)
	}

	if !auth.CheckPassword("my-secure-password", hash) {
		t.Error("expected correct password to match hash")
	}

	if auth.CheckPassword("wrong-password", hash) {
		t.Error("expected wrong password to NOT match hash")
	}
}

// TestGenerateAndValidateToken يتحقق من دورة حياة التوكن كاملة: إنشاء ثم تحقق
func TestGenerateAndValidateToken(t *testing.T) {
	auth := services.NewAuthService("test_secret")

	token, err := auth.GenerateToken("user-123", "admin")
	if err != nil {
		t.Fatalf("unexpected error generating token: %v", err)
	}

	claims, err := auth.ValidateToken(token)
	if err != nil {
		t.Fatalf("unexpected error validating token: %v", err)
	}

	if claims.UserID != "user-123" || claims.Role != "admin" {
		t.Errorf("unexpected claims: %+v", claims)
	}
}
