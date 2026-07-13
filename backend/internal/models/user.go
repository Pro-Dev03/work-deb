package models

import "time"

// User يمثل صف واحد من جدول users. الحقل Role يحدد إن كان "admin" أو "employee"
type User struct {
	ID           string    `json:"id"`
	FullName     string    `json:"full_name"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone"`
	PasswordHash string    `json:"-"` // "-" يعني: لا تُرجَع هذه القيمة أبداً في أي رد JSON
	Role         string    `json:"role"` // admin | employee
	AvatarURL    string    `json:"avatar_url,omitempty"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
}
