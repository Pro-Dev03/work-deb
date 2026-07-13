package utils

import "math"

const earthRadiusMeters = 6371000

// HaversineDistance تحسب المسافة بالمتر بين نقطتين جغرافيتين (خط عرض/طول)
// هذه هي الدالة الرياضية الأساسية التي يعتمد عليها منع التختيم خارج نطاق العمل:
// كل تحقق Geofence في المشروع يمر من هنا في النهاية
func HaversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	lat1Rad := degreesToRadians(lat1)
	lat2Rad := degreesToRadians(lat2)
	deltaLat := degreesToRadians(lat2 - lat1)
	deltaLon := degreesToRadians(lon2 - lon1)

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLon/2)*math.Sin(deltaLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadiusMeters * c
}

func degreesToRadians(deg float64) float64 {
	return deg * math.Pi / 180
}
