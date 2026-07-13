package services

import (
	"database/sql"

	"github.com/google/uuid"
)

type NotificationService struct {
	DB *sql.DB
}

func NewNotificationService(db *sql.DB) *NotificationService {
	return &NotificationService{DB: db}
}

// Send ينشئ إشعاراً جديداً لموظف معين (title/body تُمرَّر مُترجَمة مسبقاً من الـ Handler
// عبر i18n.T، حتى يصل الإشعار بنفس لغة واجهة الموظف الذي سيقرأه)
func (s *NotificationService) Send(userID, title, body string) error {
	_, err := s.DB.Exec(`
		INSERT INTO notifications (id, user_id, title, body, is_read, created_at)
		VALUES ($1, $2, $3, $4, FALSE, now())`,
		uuid.NewString(), userID, title, body,
	)
	return err
}
