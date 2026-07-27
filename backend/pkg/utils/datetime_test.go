package utils

import (
	"testing"
	"time"
)

func TestCalculateDayNightHours(t *testing.T) {
	location := JerusalemLocation()

	// اختبار وردية نهارية بالكامل (8 صباحاً - 4 مساءً)
	start := time.Date(2024, 1, 15, 8, 0, 0, 0, location)
	end := time.Date(2024, 1, 15, 16, 0, 0, 0, location)
	result := CalculateDayNightHours(start, end)

	if result.DayHours != 8.0 {
		t.Errorf("Expected 8.0 day hours, got %.2f", result.DayHours)
	}
	if result.NightHours != 0.0 {
		t.Errorf("Expected 0.0 night hours, got %.2f", result.NightHours)
	}

	// اختبار وردية ليلية بالكامل (10 مساءً - 6 صباحاً)
	start = time.Date(2024, 1, 15, 22, 0, 0, 0, location)
	end = time.Date(2024, 1, 16, 6, 0, 0, 0, location)
	result = CalculateDayNightHours(start, end)

	if result.DayHours != 0.0 {
		t.Errorf("Expected 0.0 day hours, got %.2f", result.DayHours)
	}
	if result.NightHours != 8.0 {
		t.Errorf("Expected 8.0 night hours, got %.2f", result.NightHours)
	}

	// اختبار وردية مختلطة (5 مساءً - 1 صباحاً)
	start = time.Date(2024, 1, 15, 17, 0, 0, 0, location)
	end = time.Date(2024, 1, 16, 1, 0, 0, 0, location)
	result = CalculateDayNightHours(start, end)

	expectedDayHours := 5.0 // 5 مساءً - 10 مساءً
	expectedNightHours := 3.0 // 10 مساءً - 1 صباحاً

	if result.DayHours != expectedDayHours {
		t.Errorf("Expected %.2f day hours, got %.2f", expectedDayHours, result.DayHours)
	}
	if result.NightHours != expectedNightHours {
		t.Errorf("Expected %.2f night hours, got %.2f", expectedNightHours, result.NightHours)
	}
}

func TestIsNightShift(t *testing.T) {
	tests := []struct {
		name        string
		nightHours  float64
		totalHours  float64
		expected    bool
	}{
		{"Pure night shift", 8.0, 8.0, true},
		{"Mostly night shift", 5.0, 8.0, true},
		{"Exactly half night", 4.0, 8.0, false},
		{"Mostly day shift", 3.0, 8.0, false},
		{"Pure day shift", 0.0, 8.0, false},
		{"Zero total hours", 4.0, 0.0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsNightShift(tt.nightHours, tt.totalHours)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestSplitShiftAcrossDays(t *testing.T) {
	location := JerusalemLocation()

	// اختبار وردية في نفس اليوم
	start := time.Date(2024, 1, 15, 8, 0, 0, 0, location)
	end := time.Date(2024, 1, 15, 16, 0, 0, 0, location)
	dayOneDate, dayTwoDate, dayOneHours, dayTwoHours := SplitShiftAcrossDays(start, end)

	if dayTwoDate != nil {
		t.Error("Expected dayTwoDate to be nil for same-day shift")
	}
	if dayOneDate == nil {
		t.Error("Expected dayOneDate to be non-nil")
	}
	if dayOneHours != 8.0 {
		t.Errorf("Expected 8.0 dayOneHours, got %.2f", dayOneHours)
	}
	if dayTwoHours != 0.0 {
		t.Errorf("Expected 0.0 dayTwoHours, got %.2f", dayTwoHours)
	}

	// اختبار وردية عبر يومين (10 مساءً - 2 صباحاً)
	start = time.Date(2024, 1, 15, 22, 0, 0, 0, location)
	end = time.Date(2024, 1, 16, 2, 0, 0, 0, location)
	dayOneDate, dayTwoDate, dayOneHours, dayTwoHours = SplitShiftAcrossDays(start, end)

	if dayTwoDate == nil {
		t.Error("Expected dayTwoDate to be non-nil for multi-day shift")
	}
	if dayOneDate == nil {
		t.Error("Expected dayOneDate to be non-nil")
	}
	if dayOneHours != 2.0 {
		t.Errorf("Expected 2.0 dayOneHours, got %.2f", dayOneHours)
	}
	if dayTwoHours != 2.0 {
		t.Errorf("Expected 2.0 dayTwoHours, got %.2f", dayTwoHours)
	}
}
