<template>
  <div>
    <h2 class="page-title">⏱️ {{ t('attendance') }}</h2>

    <!-- نقاط العمل -->
    <div class="card worksites-section">
      <h3>📍 {{ t('select_worksite') }}</h3>
      <div v-if="loadingWorksites" class="loading">{{ t('loading') }}</div>
      <div v-else-if="availableWorksites.length === 0" class="empty">
        <p>📭 {{ t('no_worksites') }}</p>
      </div>
      <div v-else class="worksites-grid">
        <button
          v-for="site in availableWorksites"
          :key="site.id"
          class="worksite-card"
          :class="{ active: selectedWorksiteId === site.id }"
          @click="selectWorksite(site)"
        >
          <span class="ws-name">{{ site.name }}</span>
          <span class="ws-address">{{ site.address || t('address') }}</span>
          <span class="ws-radius">⭕ {{ site.radius_meters }} م</span>
        </button>
      </div>
    </div>

    <!-- الموقع المختار + زر الملاحة -->
    <div v-if="selectedWorksiteId" class="card navigation-card">
      <div class="nav-info">
        <span class="nav-icon">📍</span>
        <div>
          <p class="nav-title">{{ worksiteName }}</p>
          <p class="nav-address">{{ selectedWorksite?.address || '' }}</p>
        </div>
      </div>
      <div class="nav-buttons">
        <a :href="getWazeUrl(selectedWorksite)" target="_blank" class="btn btn--waze">
          🗺️ Waze
        </a>
        <a :href="getGoogleMapsUrl(selectedWorksite)" target="_blank" class="btn btn--google">
          🌐 {{ t('google_maps') }}
        </a>
      </div>
    </div>

    <!-- عداد الوقت -->
    <div v-if="isWorking" class="timer-card">
      <div class="timer"><span>⏱️</span> <span>{{ elapsedTime }}</span></div>
      <p class="timer-label">{{ t('in_progress') }} {{ worksiteName }}</p>
    </div>

    <!-- ملخص الساعات -->
    <div class="summary">
      <div class="summary-item"><span class="num">{{ todayHours.toFixed(1) }}</span><span>{{ t('hours_today') }}</span></div>
      <div class="summary-item"><span class="num">{{ weekHours.toFixed(1) }}</span><span>{{ t('hours_week') }}</span></div>
      <div class="summary-item"><span class="num">{{ monthHours.toFixed(1) }}</span><span>{{ t('hours_month') }}</span></div>
    </div>

    <!-- حالة القرب + عنوان الموقع -->
    <div class="card attendance-card">
      <!-- الموقع الحالي والعنوان -->
      <div v-if="userLocation" class="location-info">
        <div class="location-header">
          <span class="location-icon">📍</span>
          <span class="location-title">{{ t('location') }}</span>
        </div>
        <div class="location-coords mono">
          {{ userLocation.lat.toFixed(6) }}, {{ userLocation.lng.toFixed(6) }}
        </div>
        <div v-if="locationAddress" class="location-address">
          {{ locationAddress }}
        </div>
        <div class="location-distance" :class="withinRange ? 'in' : 'out'">
          <span class="distance-icon">{{ withinRange ? '✅' : '❌' }}</span>
          <span class="distance-text">
            {{ t('distance') }}: <strong>{{ formatDistance(distance) }}</strong>
            <span v-if="selectedWorksiteId" class="distance-range">
              ({{ t('radius') }}: {{ radius }} م)
            </span>
          </span>
        </div>
      </div>

      <div class="geofence-status" v-if="selectedWorksiteId">
        <div class="status-icon" :class="withinRange ? 'in' : 'out'">
          {{ withinRange ? '✅' : '❌' }}
        </div>
        <div>
          <p class="status-text" :class="withinRange ? 'in' : 'out'">
            {{ withinRange ? t('inside_range') : t('outside_range') }}
          </p>
          <p class="status-distance">
            {{ t('distance') }}: <span class="mono">{{ formatDistance(distance) }}</span>
            ({{ t('radius') }}: <span class="mono">{{ radius }}</span> م)
          </p>
        </div>
      </div>

      <!-- زر تحديد الموقع -->
      <button 
        class="btn btn--primary" 
        @click="getLocationWithAddress" 
        :disabled="gettingLocation"
      >
        {{ gettingLocation ? '⏳' : '📍' }} {{ t('select_location') }}
      </button>

      <!-- الأزرار -->
      <div class="actions">
        <button 
          class="btn btn--primary" 
          :disabled="!selectedWorksiteId || isWorking || !withinRange || !userLocation || checkingIn" 
          @click="checkIn"
        >
          {{ checkingIn ? '⏳' : '✅' }} {{ t('check_in') }}
        </button>
        
        <button 
          class="btn btn--ghost" 
          :disabled="!isWorking || !hasClickedLocation || !withinRange || checkingOut" 
          @click="checkOut"
        >
          {{ checkingOut ? '⏳' : '⏹️' }} {{ t('check_out') }}
        </button>
      </div>

      <!-- تحذيرات -->
      <div v-if="!hasClickedLocation && isWorking" class="warning-box warning-location">
        <span class="warning-icon">⚠️</span>
        <span class="warning-text">📍 {{ t('select_location') }} {{ t('before_checkout') }}</span>
      </div>

      <div v-if="hasClickedLocation && !withinRange && isWorking" class="warning-box warning-range">
        <span class="warning-icon">🚫</span>
        <span class="warning-text">❌ {{ t('outside_range') }}! {{ t('distance') }}: {{ formatDistance(distance) }}</span>
      </div>

      <div v-if="hasClickedLocation && withinRange && isWorking" class="success-box">
        <span class="success-icon">✅</span>
        <span class="success-text">{{ t('inside_range') }} - {{ t('can_checkout') }}</span>
      </div>

      <p v-if="error" class="error">{{ error }}</p>
      <p v-if="success" class="success">{{ success }}</p>
      <p v-if="debugInfo" class="debug-info mono">{{ debugInfo }}</p>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useI18n } from '../services/i18n'
import api from '../services/api'

const { t, currentLang } = useI18n()

// ==========================================
// الحالة الأساسية
// ==========================================
const loadingWorksites = ref(false)
const availableWorksites = ref([])
const selectedWorksiteId = ref(null)
const selectedWorksite = ref(null)
const worksiteName = ref('')
const radius = ref(100)
const distance = ref(0)
const userLocation = ref(null)
const locationAddress = ref('')
const attendanceId = ref(null)
const isWorking = ref(false)
const elapsedSeconds = ref(0)
const todayHours = ref(0), weekHours = ref(0), monthHours = ref(0)
const checkingIn = ref(false), checkingOut = ref(false), gettingLocation = ref(false)
const hasClickedLocation = ref(false)
const error = ref(''), success = ref(''), debugInfo = ref('')
let timerInterval = null, locationInterval = null

const withinRange = computed(() => distance.value <= radius.value)
const elapsedTime = computed(() => {
  const s = elapsedSeconds.value
  return `${String(Math.floor(s/3600)).padStart(2,'0')}:${String(Math.floor((s%3600)/60)).padStart(2,'0')}:${String(Math.floor(s%60)).padStart(2,'0')}`
})

// ==========================================
// ✅ مراقبة تغيير اللغة وإعادة تحميل البيانات
// ==========================================
watch(currentLang, () => {
  // عند تغيير اللغة، يتم تحديث كل شيء تلقائياً لأن t() تفاعلية
  console.log('🌍 تغيير اللغة إلى:', currentLang.value)
})

// ==========================================
// دالة تنسيق المسافة
// ==========================================
function formatDistance(meters) {
  if (meters >= 1000) {
    return (meters / 1000).toFixed(2) + ' كيلومتر'
  }
  return Math.round(meters) + ' متر'
}

// ==========================================
// روابط الملاحة
// ==========================================
function getWazeUrl(site) {
  if (!site) return '#'
  return `https://www.waze.com/ul?ll=${site.latitude},${site.longitude}&navigate=yes`
}

function getGoogleMapsUrl(site) {
  if (!site) return '#'
  return `https://www.google.com/maps/dir/?api=1&destination=${site.latitude},${site.longitude}`
}

// ==========================================
// جلب نقاط العمل
// ==========================================
async function fetchWorksites() {
  loadingWorksites.value = true
  try {
    const { data } = await api.get('/worksites/available')
    availableWorksites.value = data || []
  } catch(e) { 
    console.error(e) 
  } finally { 
    loadingWorksites.value = false 
  }
}

// ==========================================
// اختيار نقطة عمل
// ==========================================
function selectWorksite(site) {
  selectedWorksiteId.value = site.id
  selectedWorksite.value = site
  worksiteName.value = site.name
  radius.value = site.radius_meters
  debugInfo.value = `تم اختيار: ${site.name} (النطاق: ${site.radius_meters}م)`
  
  if (userLocation.value) {
    calculateDistance(userLocation.value.lat, userLocation.value.lng)
  }
}

// ==========================================
// حساب المسافة
// ==========================================
function calculateDistance(lat, lng) {
  if (!selectedWorksite.value) return
  
  const site = selectedWorksite.value
  const R = 6371000
  const φ1 = lat * Math.PI / 180
  const φ2 = site.latitude * Math.PI / 180
  const Δφ = (site.latitude - lat) * Math.PI / 180
  const Δλ = (site.longitude - lng) * Math.PI / 180

  const a = Math.sin(Δφ/2) * Math.sin(Δφ/2) +
            Math.cos(φ1) * Math.cos(φ2) *
            Math.sin(Δλ/2) * Math.sin(Δλ/2)
  const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1-a))
  const d = R * c

  distance.value = Math.round(d)
  debugInfo.value = `📏 المسافة: ${formatDistance(distance.value)}`
}

// ==========================================
// جلب عنوان الموقع
// ==========================================
async function getAddressFromCoords(lat, lng) {
  try {
    const url = `https://api.geoapify.com/v1/geocode/reverse?lat=${lat}&lon=${lng}&apiKey=a6a3b5fec1cd4b1c99daaf6decab855f&lang=ar`
    const response = await fetch(url)
    if (!response.ok) throw new Error('فشل جلب العنوان')
    
    const data = await response.json()
    if (data.features && data.features.length > 0) {
      const props = data.features[0].properties
      return props.formatted || props.address_line1 || 'عنوان غير معروف'
    }
    return 'لم يتم العثور على عنوان'
  } catch (e) {
    console.error('فشل جلب العنوان:', e)
    return 'تعذر الحصول على العنوان'
  }
}

// ==========================================
// تحديد الموقع مع العنوان
// ==========================================
async function getLocationWithAddress() {
  if (!navigator.geolocation) {
    error.value = 'المتصفح لا يدعم تحديد الموقع'
    return
  }

  gettingLocation.value = true
  hasClickedLocation.value = true
  error.value = ''
  success.value = ''
  locationAddress.value = ''
  debugInfo.value = '⏳ جاري تحديد الموقع...'

  try {
    const pos = await new Promise((resolve, reject) => {
      navigator.geolocation.getCurrentPosition(resolve, reject, {
        enableHighAccuracy: true,
        timeout: 15000,
        maximumAge: 0
      })
    })

    const lat = pos.coords.latitude
    const lng = pos.coords.longitude

    userLocation.value = { lat, lng }
    debugInfo.value = `📍 الموقع: ${lat.toFixed(6)}, ${lng.toFixed(6)}`

    const address = await getAddressFromCoords(lat, lng)
    locationAddress.value = address
    debugInfo.value += `\n🏠 ${address}`

    if (selectedWorksite.value) {
      calculateDistance(lat, lng)
      debugInfo.value += `\n📏 المسافة: ${formatDistance(distance.value)}`
    }

    success.value = '✅ تم تحديد موقعك بنجاح'
  } catch (e) {
    error.value = '❌ فشل تحديد الموقع: ' + (e.message || 'خطأ غير معروف')
    debugInfo.value = `❌ ${error.value}`
  } finally {
    gettingLocation.value = false
  }
}

// ==========================================
// بدء الدوام
// ==========================================
async function checkIn() {
  if (!userLocation.value) { 
    error.value = '📍 اضغط على "تحديد موقعي" أولاً'
    return 
  }
  if (!selectedWorksiteId.value) { 
    error.value = '📌 اختر نقطة العمل أولاً'
    return 
  }
  if (!withinRange.value) {
    error.value = `❌ أنت خارج النطاق! المسافة: ${formatDistance(distance.value)} (المسموح: ${radius.value}م)`
    return
  }
  
  checkingIn.value = true
  error.value = ''
  success.value = ''
  debugInfo.value = '⏳ جاري تسجيل بدء الدوام...'

  try {
    const payload = {
      worksite_id: selectedWorksiteId.value,
      latitude: userLocation.value.lat,
      longitude: userLocation.value.lng
    }
    
    const { data } = await api.post('/attendance/check-in', payload)
    
    attendanceId.value = data.attendance?.id
    success.value = '✅ بدء الدوام بنجاح!'
    isWorking.value = true
    elapsedSeconds.value = 0
    timerInterval = setInterval(() => elapsedSeconds.value++, 1000)
    startLocationTracking()
    debugInfo.value = `✅ بدء الدوام في ${worksiteName.value}`
    
    hasClickedLocation.value = false
  } catch(e) {
    console.error('❌ فشل CheckIn:', e.response?.data)
    const errData = e.response?.data
    if (errData?.geofence) {
      distance.value = errData.geofence.distance_meters || distance.value
      error.value = `❌ خارج النطاق! المسافة: ${formatDistance(distance.value)} (المسموح: ${radius.value}م)`
      debugInfo.value = `📏 المسافة الفعلية: ${formatDistance(distance.value)}`
    } else {
      error.value = errData?.error || '❌ فشل بدء الدوام'
      debugInfo.value = `❌ ${error.value}`
    }
  } finally {
    checkingIn.value = false
  }
}

// ==========================================
// إنهاء الدوام
// ==========================================
async function checkOut() {
  if (!attendanceId.value) { 
    error.value = '⚠️ لا توجد وردية نشطة'
    return 
  }
  
  if (!hasClickedLocation.value) {
    error.value = '📍 يجب الضغط على "تحديد موقعي" أولاً قبل إنهاء الدوام!'
    debugInfo.value = '⚠️ اضغط على "تحديد موقعي" ثم حاول مرة أخرى'
    return
  }

  if (!userLocation.value) {
    error.value = '📍 فشل تحديد الموقع. حاول مرة أخرى'
    debugInfo.value = '⚠️ اضغط على "تحديد موقعي" مرة أخرى'
    return
  }

  if (!withinRange.value) {
    error.value = `❌ أنت خارج نطاق موقع العمل! المسافة: ${formatDistance(distance.value)} (المسموح: ${radius.value}م)`
    debugInfo.value = '⚠️ لا يمكن إنهاء الدوام خارج النطاق المسموح'
    return
  }

  checkingOut.value = true
  error.value = ''
  success.value = ''
  debugInfo.value = '⏳ جاري تسجيل إنهاء الدوام...'

  try {
    const { data } = await api.post('/attendance/check-out', {
      attendance_id: attendanceId.value,
      latitude: userLocation.value.lat,
      longitude: userLocation.value.lng
    })
    success.value = `✅ انتهى الدوام (${data.worked_hours?.toFixed(2)} ساعة)`
    isWorking.value = false
    clearInterval(timerInterval)
    clearInterval(locationInterval)
    await fetchSummary()
    debugInfo.value = `✅ انتهى الدوام بعد ${data.worked_hours?.toFixed(2)} ساعة`
    
    hasClickedLocation.value = false
  } catch(e) {
    console.error('❌ فشل CheckOut:', e.response?.data)
    error.value = e.response?.data?.error || '❌ فشل إنهاء الدوام'
    debugInfo.value = `❌ ${error.value}`
  } finally {
    checkingOut.value = false
  }
}

// ==========================================
// تتبع الموقع
// ==========================================
function startLocationTracking() {
  clearInterval(locationInterval)
  locationInterval = setInterval(async () => {
    if (userLocation.value && isWorking.value) {
      try {
        await api.post('/location/update', {
          latitude: userLocation.value.lat,
          longitude: userLocation.value.lng
        })
      } catch(e) { /* silent */ }
    }
  }, 10000)
}

// ==========================================
// استعادة الحالة
// ==========================================
async function restoreState() {
  try {
    const { data } = await api.get('/attendance/current')
    if (data.has_active) {
      isWorking.value = true
      attendanceId.value = data.attendance_id
      elapsedSeconds.value = data.elapsed_seconds || 0
      selectedWorksiteId.value = data.worksite_id
      timerInterval = setInterval(() => elapsedSeconds.value++, 1000)
      startLocationTracking()
      const ws = availableWorksites.value.find(w => w.id === data.worksite_id)
      if(ws) {
        selectedWorksite.value = ws
        worksiteName.value = ws.name
        radius.value = ws.radius_meters
      }
      debugInfo.value = '🔄 تم استعادة الوردية النشطة'
      hasClickedLocation.value = false
    }
  } catch(e) { console.error(e) }
}

// ==========================================
// جلب الملخص
// ==========================================
async function fetchSummary() {
  try {
    const { data } = await api.get('/attendance/summary')
    todayHours.value = data.today_hours||0
    weekHours.value = data.week_hours||0
    monthHours.value = data.month_hours||0
  } catch(e) {}
}

// ==========================================
// دورة الحياة
// ==========================================
onMounted(async () => {
  await fetchWorksites()
  await restoreState()
  await fetchSummary()
})

onUnmounted(() => { 
  clearInterval(timerInterval)
  clearInterval(locationInterval)
})
</script>

<style scoped>
.page-title { font-size: 20px; margin-bottom: 16px; }

.worksites-section { padding: 16px; margin-bottom: 16px; }
.worksites-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 8px; margin-top: 10px; }
.worksite-card { padding: 12px; border: 2px solid var(--line); border-radius: var(--radius-sm); background: var(--surface); cursor: pointer; text-align: right; transition: 0.2s; }
.worksite-card.active { border-color: var(--brand); background: var(--brand-tint); }
.ws-name { font-weight: 600; display: block; }
.ws-address { font-size: 12px; color: var(--ink-soft); }
.ws-radius { font-size: 11px; color: var(--brand); }

.navigation-card { padding: 16px; margin-bottom: 16px; display: flex; justify-content: space-between; align-items: center; flex-wrap: wrap; gap: 12px; }
.nav-info { display: flex; align-items: center; gap: 10px; }
.nav-icon { font-size: 24px; }
.nav-title { font-weight: 600; margin: 0; }
.nav-address { font-size: 12px; color: var(--ink-soft); margin: 0; }
.nav-buttons { display: flex; gap: 8px; }
.btn { padding: 8px 16px; border-radius: var(--radius-sm); font-weight: 600; text-decoration: none; display: inline-block; font-size: 13px; cursor: pointer; border: none; }
.btn--waze { background: #33ccff; color: #1a1a2e; }
.btn--google { background: #4285f4; color: white; }

.location-info { 
  width: 100%; 
  padding: 14px; 
  background: var(--canvas); 
  border-radius: var(--radius-sm); 
  border: 1px solid var(--line);
}
.location-header { display: flex; align-items: center; gap: 8px; margin-bottom: 6px; }
.location-icon { font-size: 18px; }
.location-title { font-weight: 600; font-size: 14px; }
.location-coords { font-size: 12px; color: var(--ink-soft); margin-bottom: 4px; }
.location-address { font-size: 13px; color: var(--ink); padding: 4px 8px; background: var(--surface); border-radius: var(--radius-sm); margin-bottom: 6px; }
.location-distance { display: flex; align-items: center; gap: 8px; padding: 6px 10px; border-radius: var(--radius-sm); font-size: 13px; }
.location-distance.in { background: var(--signal-in-tint); color: var(--signal-in); }
.location-distance.out { background: var(--signal-out-tint); color: var(--signal-out); }
.distance-icon { font-size: 16px; }
.distance-range { font-size: 12px; opacity: 0.7; }

.timer-card { padding: 16px; margin-bottom: 16px; text-align: center; background: var(--signal-in-tint); border-radius: var(--radius-md); }
.timer { font-size: 28px; font-weight: 700; font-family: var(--font-mono); color: var(--signal-in); display: flex; align-items: center; justify-content: center; gap: 8px; }
.timer-label { font-size: 13px; color: var(--ink-soft); margin-top: 4px; }

.summary { display: grid; grid-template-columns: repeat(3,1fr); gap: 10px; margin-bottom: 16px; }
.summary-item { background: var(--surface); padding: 12px; text-align: center; border-radius: var(--radius-sm); border: 1px solid var(--line); }
.summary-item .num { font-size: 22px; font-weight: 700; color: var(--brand); display: block; }
.summary-item span:last-child { font-size: 11px; color: var(--ink-soft); }

.attendance-card { padding: 20px; display: flex; flex-direction: column; gap: 16px; }
.geofence-status { display: flex; align-items: center; gap: 12px; padding: 12px; background: var(--canvas); border-radius: var(--radius-sm); }
.status-icon { font-size: 28px; }
.status-text { font-weight: 600; margin: 0; }
.status-text.in { color: var(--signal-in); } .status-text.out { color: var(--signal-out); }
.status-distance { font-size: 12px; color: var(--ink-soft); margin: 0; }
.status-distance .mono { font-family: var(--font-mono); font-weight: 600; }

.actions { display: flex; flex-direction: column; gap: 10px; width: 100%; }
.btn--primary { background: var(--brand); color: white; padding: 14px; border: none; border-radius: var(--radius-sm); font-weight: 600; cursor: pointer; width: 100%; font-size: 16px; }
.btn--primary:disabled { opacity: 0.5; cursor: not-allowed; }
.btn--ghost { background: transparent; border: 1.5px solid var(--line); color: var(--ink); padding: 14px; border-radius: var(--radius-sm); font-weight: 600; cursor: pointer; width: 100%; font-size: 16px; }
.btn--ghost:disabled { opacity: 0.5; cursor: not-allowed; }

.warning-box {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 16px;
  border-radius: var(--radius-sm);
  font-weight: 600;
}

.warning-location {
  background: var(--signal-warning-tint);
  border: 2px solid var(--signal-warning);
  color: var(--signal-warning);
}

.warning-range {
  background: var(--signal-out-tint);
  border: 2px solid var(--signal-out);
  color: var(--signal-out);
}

.success-box {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 16px;
  background: var(--signal-in-tint);
  border: 2px solid var(--signal-in);
  border-radius: var(--radius-sm);
  color: var(--signal-in);
  font-weight: 600;
}

.success-icon { font-size: 24px; }
.success-text { font-size: 14px; }
.warning-icon { font-size: 24px; }
.warning-text { font-size: 14px; }

.error { color: var(--signal-out); text-align: center; }
.success { color: var(--signal-in); text-align: center; }
.debug-info { padding: 8px 12px; background: var(--ink); color: #fff; border-radius: var(--radius-sm); font-size: 12px; text-align: center; white-space: pre-line; }
.loading, .empty { text-align: center; padding: 20px; color: var(--ink-soft); }
</style>
