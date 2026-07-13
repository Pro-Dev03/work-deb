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
