package services

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type RefreshTokenService struct {
	DB *sql.DB
}

func NewRefreshTokenService(db *sql.DB) *RefreshTokenService {
	return &RefreshTokenService{DB: db}
}

// GenerateRefreshToken ينشئ refresh token جديد
func (s *RefreshTokenService) GenerateRefreshToken(userID string, userAgent, ipAddress string) (string, error) {
	// إنشاء refresh token عشوائي
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", fmt.Errorf("فشل إنشاء refresh token: %w", err)
	}
	
	token := base64.URLEncoding.EncodeToString(tokenBytes)
	
	// إنشاء hash للـ token (لا نخزن الـ token نفسه)
	hash := sha256.Sum256([]byte(token))
	tokenHash := base64.URLEncoding.EncodeToString(hash[:])
	
	// تحديد وقت انتهاء الصلاحية (7 أيام)
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	
	// إدخال الـ token في قاعدة البيانات
	tokenID := uuid.NewString()
	_, err := s.DB.Exec(`
		INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at, user_agent, ip_address)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, tokenID, userID, tokenHash, expiresAt, userAgent, ipAddress)
	
	if err != nil {
		return "", fmt.Errorf("فشل حفظ refresh token: %w", err)
	}
	
	// تحديث current_refresh_token_id في جدول users
	_, err = s.DB.Exec(`
		UPDATE users SET current_refresh_token_id = $1 WHERE id = $2
	`, tokenID, userID)
	
	if err != nil {
		// ليس خطأ قاتل إذا فشل التحديث
		fmt.Printf("⚠️ فشل تحديث current_refresh_token_id: %v\n", err)
	}
	
	return token, nil
}

// ValidateRefreshToken يتحقق من صحة refresh token
func (s *RefreshTokenService) ValidateRefreshToken(token string) (string, error) {
	// إنشاء hash للـ token
	hash := sha256.Sum256([]byte(token))
	tokenHash := base64.URLEncoding.EncodeToString(hash[:])
	
	var userID string
	var expiresAt time.Time
	var isRevoked bool
	
	err := s.DB.QueryRow(`
		SELECT user_id, expires_at, is_revoked
		FROM refresh_tokens
		WHERE token_hash = $1
	`, tokenHash).Scan(&userID, &expiresAt, &isRevoked)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("refresh token غير صالح")
		}
		return "", fmt.Errorf("فشل التحقق من refresh token: %w", err)
	}
	
	// التحقق من عدم إلغاء الـ token
	if isRevoked {
		return "", fmt.Errorf("refresh token ملغى")
	}
	
	// التحقق من عدم انتهاء الصلاحية
	if time.Now().After(expiresAt) {
		return "", fmt.Errorf("refresh token منتهي الصلاحية")
	}
	
	return userID, nil
}

// RevokeRefreshToken يلغي refresh token
func (s *RefreshTokenService) RevokeRefreshToken(token string) error {
	hash := sha256.Sum256([]byte(token))
	tokenHash := base64.URLEncoding.EncodeToString(hash[:])
	
	_, err := s.DB.Exec(`
		UPDATE refresh_tokens
		SET is_revoked = TRUE, revoked_at = now()
		WHERE token_hash = $1
	`, tokenHash)
	
	if err != nil {
		return fmt.Errorf("فشل إلغاء refresh token: %w", err)
	}
	
	return nil
}

// RevokeAllUserTokens يلغي جميع refresh tokens لمستخدم معين
func (s *RefreshTokenService) RevokeAllUserTokens(userID string) error {
	_, err := s.DB.Exec(`
		UPDATE refresh_tokens
		SET is_revoked = TRUE, revoked_at = now()
		WHERE user_id = $1 AND is_revoked = FALSE
	`, userID)
	
	if err != nil {
		return fmt.Errorf("فشل إلغاء جميع refresh tokens: %w", err)
	}
	
	return nil
}

// CleanupExpiredTokens يحذف الـ tokens المنتهية الصلاحية
func (s *RefreshTokenService) CleanupExpiredTokens() error {
	_, err := s.DB.Exec(`
		DELETE FROM refresh_tokens
		WHERE expires_at < now() OR (is_revoked = TRUE AND revoked_at < now() - INTERVAL '30 days')
	`)
	
	if err != nil {
		return fmt.Errorf("فشل حذف الـ tokens المنتهية: %w", err)
	}
	
	return nil
}

// RotateRefreshToken يدور refresh token (يخلق جديد ويلغي القديم)
func (s *RefreshTokenService) RotateRefreshToken(oldToken string, userAgent, ipAddress string) (string, error) {
	// التحقق من الـ token القديم
	userID, err := s.ValidateRefreshToken(oldToken)
	if err != nil {
		return "", err
	}
	
	// إلغاء الـ token القديم
	if err := s.RevokeRefreshToken(oldToken); err != nil {
		return "", fmt.Errorf("فشل إلغاء الـ token القديم: %w", err)
	}
	
	// إنشاء token جديد
	return s.GenerateRefreshToken(userID, userAgent, ipAddress)
}
