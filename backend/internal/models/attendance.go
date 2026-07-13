package models

import "time"

// Attendance يمثل سجل حضور/انصراف واحد لموظف على مهمة معينة
// يحتوي على إثبات الموقع الجغرافي (lat/lng + المسافة المحسوبة) لكل من البدء والانتهاء
type Attendance struct {
	ID         string `json:"id"`
	UserID     string `json:"user_id"`
	TaskID     *string `json:"task_id,omitempty"`
	WorksiteID string `json:"worksite_id"`

	CheckInTime           *time.Time `json:"check_in_time,omitempty"`
	CheckInLat            *float64   `json:"check_in_lat,omitempty"`
	CheckInLng            *float64   `json:"check_in_lng,omitempty"`
	CheckInDistanceMeters *float64   `json:"check_in_distance_meters,omitempty"`

	CheckOutTime           *time.Time `json:"check_out_time,omitempty"`
	CheckOutLat            *float64   `json:"check_out_lat,omitempty"`
	CheckOutLng            *float64   `json:"check_out_lng,omitempty"`
	CheckOutDistanceMeters *float64   `json:"check_out_distance_meters,omitempty"`

	Status    string    `json:"status"` // in_progress | completed
	CreatedAt time.Time `json:"created_at"`
}
