// يعيد موقع الموظف الحالي من متصفح الجهاز (GPS)
// التحقق الفعلي من النطاق الجغرافي يحصل في الـ Backend وليس هنا
export function getCurrentPosition() {
  return new Promise((resolve, reject) => {
    if (!('geolocation' in navigator)) {
      reject(new Error('المتصفح لا يدعم تحديد الموقع الجغرافي'))
      return
    }
    navigator.geolocation.getCurrentPosition(
      (pos) => resolve({ lat: pos.coords.latitude, lng: pos.coords.longitude }),
      (err) => reject(err),
      { enableHighAccuracy: true, timeout: 10000 }
    )
  })
}
