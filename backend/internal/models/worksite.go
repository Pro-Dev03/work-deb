package models

import "time"

// Worksite تمثل نقطة عمل جغرافية (Geofence Zone)
// radius_meters هو نصف القطر المسموح للموظف أن "يختم" داخله
type Worksite struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Address      string    `json:"address"`
	Latitude     float64   `json:"latitude"`
	Longitude    float64   `json:"longitude"`
	RadiusMeters int       `json:"radius_meters"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
}
