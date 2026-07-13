package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"worktrack/backend/internal/i18n"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	DB *sql.DB
}

func NewNotificationHandler(db *sql.DB) *NotificationHandler {
	return &NotificationHandler{DB: db}
}

// List يعيد إشعارات المستخدم الحالي فقط. عناوين ونصوص الإشعارات نفسها كانت
// قد خُزّنت مسبقاً بلغة الموظف وقت إنشائها (راجع notif_checkin_rejected_* في i18n)
func (h *NotificationHandler) List(c *gin.Context) {
	lang := i18n.Detect(c)
	userID, _ := c.Get("user_id")

	rows, err := h.DB.Query(`
		SELECT id, title, body, is_read, created_at FROM notifications
		WHERE user_id = $1 ORDER BY created_at DESC LIMIT 50`, userID)
	if err != nil {
		log.Printf("failed to fetch notifications: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.T(lang, "err_operation_failed")})
		return
	}
	defer rows.Close()

	var notifications []gin.H
	for rows.Next() {
		var id, title, body, createdAt string
		var isRead bool
		if err := rows.Scan(&id, &title, &body, &isRead, &createdAt); err == nil {
			notifications = append(notifications, gin.H{
				"id": id, "title": title, "body": body, "is_read": isRead, "created_at": createdAt,
			})
		}
	}

	c.JSON(http.StatusOK, notifications)
}
