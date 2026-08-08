package services

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

// SecurityLogger يسجل أحداث الأمان والمحاولات المشبوهة
type SecurityLogger struct {
	DB *sql.DB
}

func NewSecurityLogger(db *sql.DB) *SecurityLogger {
	return &SecurityLogger{DB: db}
}

// LogAuthAttempt يسجل محاولة مصادقة
func (sl *SecurityLogger) LogAuthAttempt(email, phone, ip, userAgent string, success bool, reason string) {
	// التحقق من وجود الجدول أولاً
	var tableExists bool
	err := sl.DB.QueryRow(`
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = 'security_logs'
		)
	`).Scan(&tableExists)
	
	if err != nil || !tableExists {
		// الجدول غير موجود، نتخطى التسجيل
		return
	}
	
	_, err = sl.DB.Exec(`
		INSERT INTO security_logs (event_type, email, phone, ip, user_agent, success, reason, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, "auth_attempt", email, phone, ip, userAgent, success, reason, time.Now())
	
	if err != nil {
		log.Printf("⚠️ فشل تسجيل محاولة مصادقة: %v", err)
	}
}

// LogSuspiciousActivity يسجل نشاط مشبوه
func (sl *SecurityLogger) LogSuspiciousActivity(userID, ip, userAgent, activityType, details string) {
	// التحقق من وجود الجدول أولاً
	var tableExists bool
	err := sl.DB.QueryRow(`
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = 'security_logs'
		)
	`).Scan(&tableExists)
	
	if err != nil || !tableExists {
		// الجدول غير موجود، نتخطى التسجيل
		return
	}
	
	_, err = sl.DB.Exec(`
		INSERT INTO security_logs (event_type, user_id, ip, user_agent, details, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, activityType, userID, ip, userAgent, details, time.Now())
	
	if err != nil {
		log.Printf("⚠️ فشل تسجيل نشاط مشبوه: %v", err)
	}
}

// CheckForSuspiciousPatterns يتحقق من أنماط مشبوهة في محاولات الدخول
func (sl *SecurityLogger) CheckForSuspiciousPatterns(ip, email string) (bool, string) {
	// التحقق من وجود الجدول أولاً
	var tableExists bool
	err := sl.DB.QueryRow(`
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = 'security_logs'
		)
	`).Scan(&tableExists)
	
	if err != nil || !tableExists {
		// الجدول غير موجود، لا نتحقق من الأنماط
		return false, ""
	}
	
	// التحقق من محاولات فاشلة متعددة من نفس IP (خلال 15 دقيقة)
	var failedAttempts int
	fifteenMinutesAgo := time.Now().Add(-15 * time.Minute)
	err = sl.DB.QueryRow(`
		SELECT COUNT(*) FROM security_logs 
		WHERE ip = $1 AND event_type = 'auth_attempt' AND success = false 
		AND created_at > $2
	`, ip, fifteenMinutesAgo).Scan(&failedAttempts)
	
	if err == nil && failedAttempts > 5 {
		return true, fmt.Sprintf("محاولات فاشلة متعددة من IP %s (%d محاولة)", ip, failedAttempts)
	}
	
	// التحقق من محاولات فاشلة متعددة لنفس البريد (خلال ساعة واحدة)
	oneHourAgo := time.Now().Add(-1 * time.Hour)
	err = sl.DB.QueryRow(`
		SELECT COUNT(*) FROM security_logs 
		WHERE email = $1 AND event_type = 'auth_attempt' AND success = false 
		AND created_at > $2
	`, email, oneHourAgo).Scan(&failedAttempts)
	
	if err == nil && failedAttempts > 10 {
		return true, fmt.Sprintf("محاولات فاشلة متعددة للبريد %s (%d محاولة)", email, failedAttempts)
	}
	
	return false, ""
}

// GetRecentSecurityLogs يجلب سجلات الأمان الأخيرة
func (sl *SecurityLogger) GetRecentSecurityLogs(limit int) ([]map[string]interface{}, error) {
	// التحقق من وجود الجدول أولاً
	var tableExists bool
	err := sl.DB.QueryRow(`
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = 'security_logs'
		)
	`).Scan(&tableExists)
	
	if err != nil || !tableExists {
		// الجدول غير موجود، نرجع مصفوفة فارغة
		return []map[string]interface{}{}, nil
	}
	
	rows, err := sl.DB.Query(`
		SELECT id, event_type, email, phone, ip, user_agent, success, reason, details, created_at
		FROM security_logs
		ORDER BY created_at DESC
		LIMIT $1
	`, limit)
	
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var logs []map[string]interface{}
	for rows.Next() {
		var id, eventType, email, phone, ip, userAgent, reason, details string
		var success bool
		var createdAt time.Time
		
		err := rows.Scan(&id, &eventType, &email, &phone, &ip, &userAgent, &success, &reason, &details, &createdAt)
		if err != nil {
			continue
		}
		
		log := map[string]interface{}{
			"id":          id,
			"event_type":  eventType,
			"email":       email,
			"phone":       phone,
			"ip":          ip,
			"user_agent":  userAgent,
			"success":     success,
			"reason":      reason,
			"details":     details,
			"created_at":  createdAt,
		}
		
		logs = append(logs, log)
	}
	
	return logs, nil
}
