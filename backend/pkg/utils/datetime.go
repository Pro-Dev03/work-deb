package utils

import "time"

// JerusalemLocation يعيد توقيت القدس، يُستخدم في كل مكان يتعلق بتوقيت المهام والحضور
// (المشروع يستهدف مستخدمين في نفس المنطقة الزمنية بغض النظر عن لغة الواجهة)
func JerusalemLocation() *time.Location {
	loc, err := time.LoadLocation("Asia/Jerusalem")
	if err != nil {
		return time.UTC
	}
	return loc
}

func NowInJerusalem() time.Time {
	return time.Now().In(JerusalemLocation())
}

// ShiftTimePeriods يمثل الفترات الزمنية لتقسيم ساعات العمل
type ShiftTimePeriods struct {
	NightHours float64 // ساعات العمل الليلي (10 مساءً - 6 صباحاً)
	DayHours   float64 // ساعات العمل النهاري (6 صباحاً - 10 مساءً)
}

// CalculateDayNightHours يحسب ساعات العمل الليلية والنهارية في فترة زمنية معينة
// العمل الليلي: 22:00 - 06:00
// العمل النهاري: 06:00 - 22:00
func CalculateDayNightHours(start, end time.Time) ShiftTimePeriods {
	result := ShiftTimePeriods{}

	// التأكد من أن وقت النهاية بعد وقت البدء
	if end.Before(start) {
		return result
	}

	// تعريف فترات العمل الليلي والنهاري
	nightStart := time.Date(start.Year(), start.Month(), start.Day(), 22, 0, 0, 0, start.Location())
	nightEnd := time.Date(start.Year(), start.Month(), start.Day()+1, 6, 0, 0, 0, start.Location())
	dayStart := time.Date(start.Year(), start.Month(), start.Day(), 6, 0, 0, 0, start.Location())
	dayEnd := time.Date(start.Year(), start.Month(), start.Day(), 22, 0, 0, 0, start.Location())

	// حساب الفترات الليلية في نفس اليوم
	currentNightStart := nightStart
	currentNightEnd := nightEnd

	for currentNightStart.Before(end) || currentNightStart.Equal(end) {
		periodStart := maxTime(start, currentNightStart)
		periodEnd := minTime(end, currentNightEnd)

		if periodEnd.After(periodStart) {
			duration := periodEnd.Sub(periodStart).Hours()
			result.NightHours += duration
		}

		// الانتقال إلى الليلة التالية
		currentNightStart = currentNightStart.Add(24 * time.Hour)
		currentNightEnd = currentNightEnd.Add(24 * time.Hour)
	}

	// حساب الفترات النهارية
	currentDayStart := dayStart
	currentDayEnd := dayEnd

	for currentDayStart.Before(end) || currentDayStart.Equal(end) {
		periodStart := maxTime(start, currentDayStart)
		periodEnd := minTime(end, currentDayEnd)

		if periodEnd.After(periodStart) {
			duration := periodEnd.Sub(periodStart).Hours()
			result.DayHours += duration
		}

		// الانتقال إلى اليوم التالي
		currentDayStart = currentDayStart.Add(24 * time.Hour)
		currentDayEnd = currentDayEnd.Add(24 * time.Hour)
	}

	return result
}

// IsNightShift يحدد ما إذا كانت الوردية ليلية بشكل أساسي
// تعتبر الوردية ليلية إذا كانت أكثر من 50% من ساعاتها ليلية
func IsNightShift(nightHours, totalHours float64) bool {
	if totalHours == 0 {
		return false
	}
	return (nightHours / totalHours) > 0.5
}

// SplitShiftAcrossDays يقسم الوردية إذا عبرت منتصف الليل
func SplitShiftAcrossDays(start, end time.Time) (dayOneDate, dayTwoDate *time.Time, dayOneHours, dayTwoHours float64) {
	if start.Day() == end.Day() && start.Month() == end.Month() && start.Year() == end.Year() {
		// نفس اليوم، لا تقسيم
		dayOneDate = &start
		dayOneHours = end.Sub(start).Hours()
		return
	}

	// عبرت منتصف الليل
	midnight := time.Date(start.Year(), start.Month(), start.Day()+1, 0, 0, 0, 0, start.Location())

	dayOneDate = &start
	dayOneHours = midnight.Sub(start).Hours()

	dayTwoDate = &midnight
	dayTwoHours = end.Sub(midnight).Hours()

	return
}

// helper functions
func maxTime(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}

func minTime(a, b time.Time) time.Time {
	if a.Before(b) {
		return a
	}
	return b
}
