package models

import "time"

// Task تمثل مهمة عمل واحدة مرتبطة بعميل ونقطة عمل وموظف مكلَّف (اختياري كل منها)
type Task struct {
	ID             string     `json:"id"`
	Title          string     `json:"title"`
	Description    string     `json:"description"`
	ClientID       *string    `json:"client_id,omitempty"`
	WorksiteID     string     `json:"worksite_id"`
	AssignedUserID *string    `json:"assigned_user_id,omitempty"`
	Status         string     `json:"status"` // pending | in_progress | completed | late | cancelled
	ScheduledStart *time.Time `json:"scheduled_start,omitempty"`
	ScheduledEnd   *time.Time `json:"scheduled_end,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
}
