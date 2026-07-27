<template>
  <div>
    <h2 class="page-title">⏱️ {{ t('attendance') }}</h2>

    <!-- نقاط العمل -->
    <div v-if="!isWorking" class="card worksites-section">
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
          <span class="ws-radius">⭕ {{ site.radius_meters }} {{ t('meter') }}</span>
        </button>
      </div>
    </div>

    <!-- نقطة العمل النشطة (أثناء العمل) -->
    <div v-if="isWorking" class="card active-worksite-card">
      <h3>📍 {{ t('current_worksite') }}</h3>
      <div class="active-worksite-info">
        <div class="active-worksite-name">{{ worksiteName }}</div>
        <div class="active-worksite-address">{{ selectedWorksite?.address || '' }}</div>
        <div class="active-worksite-status">✅ {{ t('working_at_this_location') }}</div>
      </div>
    </div>

    <!-- الموقع المختار + زر الملاحة -->
    <div v-if="selectedWorksiteId && !isWorking" class="card navigation-card">
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

    <!-- زر سجل الحضور -->
    <div class="card">
      <button class="btn btn--primary btn--full" @click="showAttendanceHistoryModal = true">
        📊 {{ t('my_attendance_history') }}
      </button>
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
              ({{ t('radius') }}: {{ radius }} {{ t('meter') }})
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
            ({{ t('radius') }}: <span class="mono">{{ radius }}</span> {{ t('meter') }})
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

    <!-- DevPro Branding -->
    <div class="devpro-branding">
      <img src="/src/assets/company-logo.jpg" alt="DevPro Logo" class="devpro-logo-img" />
      <p class="devpro-text">Powered by DevPro</p>
    </div>

    <!-- مودال سجل الحضور -->
    <div v-if="showAttendanceHistoryModal" class="modal-backdrop" @click.self="showAttendanceHistoryModal = false">
      <div class="modal card">
        <div class="modal-header">
          <h3>📊 {{ t('my_attendance_history') }}</h3>
          <button class="modal-close" @click="showAttendanceHistoryModal = false">✕</button>
        </div>
        <div class="modal-body">
          <!-- فلاتر الشهر والسنة -->
          <div class="filters">
            <div class="filter-group">
              <label>{{ t('year') }}</label>
              <select v-model="selectedYear" @change="fetchMyAttendanceHistory" class="form-select">
                <option v-for="year in availableYears" :key="year" :value="year">{{ year }}</option>
              </select>
            </div>
            <div class="filter-group">
              <label>{{ t('month') }}</label>
              <select v-model="selectedMonth" @change="fetchMyAttendanceHistory" class="form-select">
                <option v-for="month in availableMonths" :key="month.value" :value="month.value">{{ month.label }}</option>
              </select>
            </div>
          </div>

          <!-- ملخص الشهر -->
          <div v-if="myMonthlySummary" class="monthly-summary">
            <div class="summary-card">
              <span class="summary-label">{{ t('total_hours') }}</span>
              <span class="summary-value">{{ myMonthlySummary.summary?.total_hours?.toFixed(1) || 0 }} {{ t('hours') }}</span>
            </div>
            <div class="summary-card">
              <span class="summary-label">{{ t('work_days') }}</span>
              <span class="summary-value">{{ myMonthlySummary.summary?.work_days || 0 }} {{ t('days') }}</span>
            </div>
          </div>

          <!-- جدول سجل الحضور -->
          <div v-if="loadingMyHistory" class="loading-state">
            <p>{{ t('loading') }}</p>
          </div>
          <div v-else-if="myAttendanceHistory.length === 0" class="empty-state">
            <p>{{ t('no_attendance_records') }}</p>
          </div>
          <div v-else class="table-wrapper">
            <table class="table">
              <thead>
                <tr>
                  <th>{{ t('date') }}</th>
                  <th>{{ t('worksite') }}</th>
                  <th>{{ t('check_in') }}</th>
                  <th>{{ t('check_out') }}</th>
                  <th>{{ t('worked_hours') }}</th>
                  <th>{{ t('night_hours') }}</th>
                  <th>{{ t('day_hours') }}</th>
                  <th>{{ t('location') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="record in myAttendanceHistory" :key="record.id">
                  <td class="mono">{{ formatDate(record.check_in_time) }}</td>
                  <td>{{ record.worksite_name || '—' }}</td>
                  <td class="mono">{{ formatTime(record.check_in_time) }}</td>
                  <td class="mono">{{ record.check_out_time ? formatTime(record.check_out_time) : '—' }}</td>
                  <td class="mono">{{ formatHours(record.worked_hours) }}</td>
                  <td class="mono">{{ formatHours(record.night_hours) }}</td>
                  <td class="mono">{{ formatHours(record.day_hours) }}</td>
                  <td class="mono">{{ formatDistance(record.check_in_distance_meters) }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useI18n } from '../services/i18n'
import api from '../services/api'
import { i18nState } from '../services/i18n'

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

// سجل الحضور الشخصي
const showAttendanceHistoryModal = ref(false)
const myAttendanceHistory = ref([])
const myMonthlySummary = ref(null)
const loadingMyHistory = ref(false)
const selectedYear = ref(new Date().getFullYear())
const selectedMonth = ref(String(new Date().getMonth() + 1))

const availableYears = ref([])
const availableMonths = ref([])

// دالة لتحديث أسماء الشهور حسب اللغة
function updateMonthNames() {
  const monthKeys = ['january', 'february', 'march', 'april', 'may', 'june', 'july', 'august', 'september', 'october', 'november', 'december']
  availableMonths.value = monthKeys.map((key, index) => ({
    value: String(index + 1),
    label: t(key)
  }))
}

// توليد السنوات المتاحة
const currentYear = new Date().getFullYear()
for (let i = currentYear; i >= currentYear - 2; i--) {
  availableYears.value.push(i)
}

// تحديث أسماء الشهور عند التحميل
updateMonthNames()

// مراقبة تغيير اللغة لتحديث أسماء الشهور
watch(currentLang, () => {
  updateMonthNames()
})

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
// مراقبة فتح مودال سجل الحضور لجلب البيانات
// ==========================================
watch(showAttendanceHistoryModal, (newValue) => {
  if (newValue) {
    // عند فتح المودال، جلب البيانات
    fetchMyAttendanceHistory()
  }
})

// ==========================================
// دالة تنسيق المسافة
// ==========================================
function formatDistance(meters) {
  if (!meters) return '0 ' + t('meters')
  if (meters >= 1000) {
    return (meters / 1000).toFixed(2) + ' ' + t('kilometers')
  }
  return Math.round(meters) + ' ' + t('meters')
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
  debugInfo.value = `${t('worksite_selected')} ${site.name} (${t('radius')}: ${site.radius_meters} ${t('meter')})`
  
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
    if (!response.ok) throw new Error(t('failed_fetch_address'))

    const data = await response.json()
    if (data.features && data.features.length > 0) {
      const props = data.features[0].properties
      return props.formatted || props.address_line1 || t('address_unknown')
    }
    return t('address_not_found')
  } catch (e) {
    console.error(t('failed_fetch_address'), e)
    return t('cannot_get_address')
  }
}

// ==========================================
// تحديد الموقع مع العنوان
// ==========================================
async function getLocationWithAddress() {
  if (!navigator.geolocation) {
    error.value = t('browser_not_support_location')
    return
  }

  gettingLocation.value = true
  hasClickedLocation.value = true
  error.value = ''
  success.value = ''
  locationAddress.value = ''
  debugInfo.value = t('getting_location')

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
    debugInfo.value = `${t('location_coords')} ${lat.toFixed(6)}, ${lng.toFixed(6)}`

    const address = await getAddressFromCoords(lat, lng)
    locationAddress.value = address
    debugInfo.value += `\n${t('home_icon')} ${address}`

    if (selectedWorksite.value) {
      calculateDistance(lat, lng)
      debugInfo.value += `\n${t('distance_label')} ${formatDistance(distance.value)}`
    }

    // إرسال الموقع فوراً للتتبع اللحظي
    if (isWorking.value) {
      try {
        await api.post('/location/update', {
          latitude: lat,
          longitude: lng,
          accuracy: pos.coords.accuracy || 0
        })
        debugInfo.value += `\n📡 تم إرسال الموقع فوراً`
      } catch (e) {
        console.error('فشل إرسال الموقع الفوري:', e)
      }
    }

    success.value = t('location_determined_success')
  } catch (e) {
    error.value = t('location_determine_failed') + ' ' + (e.message || t('error'))
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
    error.value = t('click_select_location_first')
    return
  }
  if (!selectedWorksiteId.value) {
    error.value = t('select_worksite_first')
    return
  }
  if (!withinRange.value) {
    error.value = `${t('outside_range_distance')} ${formatDistance(distance.value)} (${t('allowed_range')}: ${radius.value} ${t('meter')})`
    return
  }

  checkingIn.value = true
  error.value = ''
  success.value = ''
  debugInfo.value = t('registering_checkin')

  try {
    const payload = {
      worksite_id: selectedWorksiteId.value,
      latitude: userLocation.value.lat,
      longitude: userLocation.value.lng
    }

    const { data } = await api.post('/attendance/check-in', payload)

    attendanceId.value = data.attendance?.id
    success.value = t('checkin_success')
    isWorking.value = true
    elapsedSeconds.value = 0
    timerInterval = setInterval(() => elapsedSeconds.value++, 1000)
    startLocationTracking()
    debugInfo.value = `${t('checkin_started_at')} ${worksiteName.value}`

    hasClickedLocation.value = false
  } catch(e) {
    console.error('❌ فشل CheckIn:', e.response?.data)
    const errData = e.response?.data
    if (errData?.geofence) {
      distance.value = errData.geofence.distance_meters || distance.value
      error.value = `${t('outside_range_actual_distance')} ${formatDistance(distance.value)} (${t('allowed_range')}: ${radius.value} ${t('meter')})`
      debugInfo.value = `${t('actual_distance')} ${formatDistance(distance.value)}`
    } else {
      error.value = errData?.error || t('checkin_failed')
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
    error.value = t('no_active_shift')
    return
  }

  if (!hasClickedLocation.value) {
    error.value = t('click_select_location_before_checkout')
    debugInfo.value = t('click_select_location_again')
    return
  }

  if (!userLocation.value) {
    error.value = t('location_determine_failed_retry')
    debugInfo.value = t('click_select_location_retry')
    return
  }

  if (!withinRange.value) {
    error.value = `${t('outside_worksite_range')} ${formatDistance(distance.value)} (${t('allowed_range')}: ${radius.value} ${t('meter')})`
    debugInfo.value = t('cannot_checkout_outside_range')
    return
  }

  checkingOut.value = true
  error.value = ''
  success.value = ''
  debugInfo.value = t('registering_checkout')

  try {
    const { data } = await api.post('/attendance/check-out', {
      attendance_id: attendanceId.value,
      latitude: userLocation.value.lat,
      longitude: userLocation.value.lng
    })
    success.value = `${t('checkout_success')} (${data.worked_hours ? formatHours(data.worked_hours) : '—'})`
    isWorking.value = false
    clearInterval(timerInterval)
    clearInterval(locationInterval)
    await fetchSummary()
    debugInfo.value = `${t('checkout_completed_after')} ${data.worked_hours ? formatHours(data.worked_hours) : '—'}`

    hasClickedLocation.value = false
  } catch(e) {
    console.error('❌ فشل CheckOut:', e.response?.data)
    error.value = e.response?.data?.error || t('checkout_failed')
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
  }, 3000) // تحديث كل 3 ثوانٍ للتتبع اللحظي
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
      debugInfo.value = t('shift_restored')
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
// دوال سجل الحضور الشخصي
// ==========================================
async function fetchMyAttendanceHistory() {
  loadingMyHistory.value = true
  try {
    const { data } = await api.get(
      `/attendance/my-history?year=${selectedYear.value}&month=${selectedMonth.value}`
    )
    myAttendanceHistory.value = data || []
    
    // جلب الملخص الشهري
    const summaryResponse = await api.get(
      `/attendance/my-monthly-summary?year=${selectedYear.value}&month=${selectedMonth.value}`
    )
    myMonthlySummary.value = summaryResponse.data
  } catch (error) {
    console.error('فشل جلب سجل الحضور:', error)
    myAttendanceHistory.value = []
    myMonthlySummary.value = null
  } finally {
    loadingMyHistory.value = false
  }
}

function formatDate(date) {
  if (!date) return '—'
  return new Date(date).toLocaleDateString('en-GB')
}

function formatTime(date) {
  if (!date) return '—'
  return new Date(date).toLocaleTimeString('en-GB', { hour: '2-digit', minute: '2-digit', hour12: false })
}

function formatHours(hours) {
  if (!hours || hours === 0) return '—'
  
  const totalSeconds = Math.floor(hours * 3600)
  const h = Math.floor(totalSeconds / 3600)
  const m = Math.floor((totalSeconds % 3600) / 60)
  const s = totalSeconds % 60
  
  // تنسيق HH:MM:SS
  const pad = (num) => num.toString().padStart(2, '0')
  return `${pad(h)}:${pad(m)}:${pad(s)}`
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
.page-title { 
  font-size: 24px; 
  margin-bottom: 20px; 
  font-weight: 700;
  color: var(--brand);
  display: flex;
  align-items: center;
  gap: 8px;
}

.worksites-section { padding: 20px; margin-bottom: 20px; }
.worksites-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 12px; margin-top: 16px; }
.worksite-card { 
  padding: 16px 14px; 
  border: 2px solid var(--line); 
  border-radius: var(--radius-md); 
  background: var(--surface); 
  cursor: pointer; 
  text-align: right; 
  transition: all 0.3s ease;
  box-shadow: var(--shadow-sm);
}
.worksite-card:hover {
  border-color: var(--brand-light);
  box-shadow: var(--shadow-md);
  transform: translateY(-2px);
}
.worksite-card.active { 
  border-color: var(--brand); 
  background: linear-gradient(135deg, var(--brand-tint) 0%, var(--brand-tint) 100%);
  box-shadow: var(--shadow-md);
}
.ws-name { font-weight: 700; display: block; font-size: 15px; color: var(--ink); }
.ws-address { font-size: 13px; color: var(--ink-soft); margin-top: 4px; }
.ws-radius { font-size: 12px; color: var(--brand); margin-top: 6px; font-weight: 600; }

.active-worksite-card { 
  padding: 20px; 
  margin-bottom: 20px; 
  background: linear-gradient(135deg, var(--signal-in-tint) 0%, var(--signal-in-tint) 100%); 
  border: 2px solid var(--signal-in); 
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-md);
}
.active-worksite-info { text-align: center; }
.active-worksite-name { font-size: 20px; font-weight: 700; color: var(--signal-in); margin-bottom: 10px; }
.active-worksite-address { font-size: 15px; color: var(--ink); margin-bottom: 14px; }
.active-worksite-status { font-size: 17px; font-weight: 600; color: var(--signal-in); }

.navigation-card { 
  padding: 20px; 
  margin-bottom: 20px; 
  display: flex; 
  justify-content: space-between; 
  align-items: center; 
  flex-wrap: wrap; 
  gap: 16px;
  background: var(--surface);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-sm);
}
.nav-info { display: flex; align-items: center; gap: 12px; }
.nav-icon { font-size: 28px; }
.nav-title { font-weight: 700; margin: 0; font-size: 16px; color: var(--ink); }
.nav-address { font-size: 13px; color: var(--ink-soft); margin: 0; }
.nav-buttons { display: flex; gap: 10px; }
.btn { 
  padding: 10px 18px; 
  border-radius: var(--radius-md); 
  font-weight: 600; 
  text-decoration: none; 
  display: inline-block; 
  font-size: 14px; 
  cursor: pointer; 
  border: none;
  transition: all 0.2s ease;
  box-shadow: var(--shadow-sm);
}
.btn--waze { 
  background: linear-gradient(135deg, #33ccff 0%, #00b8e6 100%); 
  color: #1a1a2e; 
}
.btn--waze:hover {
  box-shadow: var(--shadow-md);
  transform: translateY(-1px);
}
.btn--google { 
  background: linear-gradient(135deg, #4285f4 0%, #3467a6 100%); 
  color: white; 
}
.btn--google:hover {
  box-shadow: var(--shadow-md);
  transform: translateY(-1px);
}

.location-info { 
  width: 100%; 
  padding: 16px; 
  background: var(--canvas); 
  border-radius: var(--radius-md); 
  border: 1px solid var(--line);
  box-shadow: var(--shadow-sm);
}
.location-header { display: flex; align-items: center; gap: 10px; margin-bottom: 8px; }
.location-icon { font-size: 20px; }
.location-title { font-weight: 700; font-size: 15px; color: var(--ink); }
.location-coords { font-size: 13px; color: var(--ink-soft); margin-bottom: 6px; font-family: var(--font-mono); }
.location-address { font-size: 14px; color: var(--ink); padding: 6px 10px; background: var(--surface); border-radius: var(--radius-sm); margin-bottom: 8px; }
.location-distance { display: flex; align-items: center; gap: 10px; padding: 8px 12px; border-radius: var(--radius-md); font-size: 14px; font-weight: 600; }
.location-distance.in { background: var(--signal-in-tint); color: var(--signal-in); }
.location-distance.out { background: var(--signal-out-tint); color: var(--signal-out); }
.distance-icon { font-size: 18px; }
.distance-range { font-size: 13px; opacity: 0.8; font-weight: 400; }

.timer-card { 
  padding: 24px; 
  margin-bottom: 20px; 
  text-align: center; 
  background: linear-gradient(135deg, var(--signal-in-tint) 0%, var(--signal-in-tint) 100%); 
  border-radius: var(--radius-lg);
  border: 2px solid var(--signal-in);
  box-shadow: var(--shadow-md);
}
.timer { 
  font-size: 32px; 
  font-weight: 700; 
  font-family: var(--font-mono); 
  color: var(--signal-in); 
  display: flex; 
  align-items: center; 
  justify-content: center; 
  gap: 10px; 
}
.timer-label { font-size: 14px; color: var(--ink-soft); margin-top: 6px; font-weight: 600; }

.summary { 
  display: grid; 
  grid-template-columns: repeat(3,1fr); 
  gap: 12px; 
  margin-bottom: 20px; 
}
.summary-item { 
  background: var(--surface); 
  padding: 16px; 
  text-align: center; 
  border-radius: var(--radius-md); 
  border: 1px solid var(--line);
  box-shadow: var(--shadow-sm);
  transition: all 0.2s ease;
}
.summary-item:hover {
  box-shadow: var(--shadow-md);
  transform: translateY(-2px);
}
.summary-item .num { 
  font-size: 26px; 
  font-weight: 700; 
  color: var(--brand); 
  display: block; 
  margin-bottom: 4px;
}
.summary-item span:last-child { 
  font-size: 12px; 
  color: var(--ink-soft); 
  font-weight: 600;
}

.btn--full { width: 100%; }

.attendance-card { 
  padding: 24px; 
  display: flex; 
  flex-direction: column; 
  gap: 20px; 
  background: var(--surface);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-sm);
}
.geofence-status { 
  display: flex; 
  align-items: center; 
  gap: 14px; 
  padding: 14px; 
  background: var(--canvas); 
  border-radius: var(--radius-md);
  border: 1px solid var(--line);
}
.status-icon { font-size: 32px; }
.status-text { font-weight: 700; margin: 0; font-size: 15px; }
.status-text.in { color: var(--signal-in); } 
.status-text.out { color: var(--signal-out); }
.status-distance { font-size: 13px; color: var(--ink-soft); margin: 0; }
.status-distance .mono { font-family: var(--font-mono); font-weight: 600; }

.actions { display: flex; flex-direction: column; gap: 12px; width: 100%; }
.btn--primary { 
  background: linear-gradient(135deg, var(--brand) 0%, var(--brand-dark) 100%); 
  color: white; 
  padding: 16px; 
  border: none; 
  border-radius: var(--radius-md); 
  font-weight: 700; 
  cursor: pointer; 
  width: 100%; 
  font-size: 16px;
  box-shadow: var(--shadow-md);
  transition: all 0.2s ease;
}
.btn--primary:hover:not(:disabled) {
  box-shadow: var(--shadow-lg);
  transform: translateY(-2px);
}
.btn--primary:disabled { 
  opacity: 0.5; 
  cursor: not-allowed; 
  transform: none !important;
}
.btn--ghost { 
  background: var(--surface); 
  border: 2px solid var(--line); 
  color: var(--ink); 
  padding: 16px; 
  border-radius: var(--radius-md); 
  font-weight: 700; 
  cursor: pointer; 
  width: 100%; 
  font-size: 16px;
  box-shadow: var(--shadow-sm);
  transition: all 0.2s ease;
}
.btn--ghost:hover:not(:disabled) {
  border-color: var(--brand-light);
  color: var(--brand);
  box-shadow: var(--shadow-md);
}
.btn--ghost:disabled { 
  opacity: 0.5; 
  cursor: not-allowed; 
}

.warning-box {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 18px;
  border-radius: var(--radius-md);
  font-weight: 600;
  font-size: 14px;
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
  gap: 12px;
  padding: 14px 18px;
  background: var(--signal-in-tint);
  border: 2px solid var(--signal-in);
  border-radius: var(--radius-md);
  color: var(--signal-in);
  font-weight: 600;
  font-size: 14px;
}

.success-icon { font-size: 24px; }
.success-text { font-size: 14px; }
.warning-icon { font-size: 24px; }
.warning-text { font-size: 14px; }

.error { 
  color: var(--signal-out); 
  text-align: center; 
  font-weight: 600;
  font-size: 14px;
}
.success { 
  color: var(--signal-in); 
  text-align: center; 
  font-weight: 600;
  font-size: 14px;
}
.debug-info { 
  padding: 10px 14px; 
  background: var(--ink); 
  color: #fff; 
  border-radius: var(--radius-md); 
  font-size: 12px; 
  text-align: center; 
  white-space: pre-line;
  font-family: var(--font-mono);
}

/* مودال سجل الحضور */
.modal-backdrop {
  position: fixed; inset: 0;
  background: rgba(0,0,0,0.7);
  backdrop-filter: blur(8px);
  display: flex; align-items: center; justify-content: center;
  z-index: 1000; padding: 20px;
}

.modal {
  width: 100%; max-width: 600px; padding: 0;
  max-height: 90vh;
  overflow-y: auto;
  background: var(--surface);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-xl);
}

.modal-header {
  display: flex; justify-content: space-between; align-items: center;
  padding: 20px 24px; border-bottom: 1px solid var(--line);
  background: var(--canvas);
}

.modal-header h3 { 
  font-size: 18px; 
  margin: 0; 
  font-weight: 700;
  color: var(--brand);
}
.modal-close { 
  background: none; 
  border: none; 
  font-size: 28px; 
  cursor: pointer; 
  color: var(--ink-soft);
  transition: color 0.2s ease;
}
.modal-close:hover {
  color: var(--signal-out);
}

.modal-body { padding: 24px; }

.filters {
  display: flex;
  gap: 16px;
  align-items: flex-end;
  margin-bottom: 24px;
  flex-wrap: wrap;
}

.filter-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.filter-group label {
  font-size: 13px;
  font-weight: 700;
  color: var(--ink-soft);
}

.form-select {
  padding: 10px 14px;
  border: 2px solid var(--line);
  border-radius: var(--radius-md);
  font-size: 14px;
  background: var(--surface);
  color: var(--ink);
  min-width: 120px;
  font-weight: 600;
  transition: all 0.2s ease;
}

.form-select:focus {
  outline: none;
  border-color: var(--brand-light);
  box-shadow: 0 0 0 3px var(--brand-tint);
}

.monthly-summary {
  display: flex;
  gap: 16px;
  margin-bottom: 24px;
}

.summary-card {
  flex: 1;
  padding: 20px;
  background: linear-gradient(135deg, var(--brand-tint) 0%, var(--brand-tint) 100%);
  border-radius: var(--radius-lg);
  text-align: center;
  border: 2px solid var(--brand);
  box-shadow: var(--shadow-md);
}

.summary-label {
  display: block;
  font-size: 13px;
  color: var(--ink-soft);
  margin-bottom: 10px;
  font-weight: 600;
}

.summary-value {
  display: block;
  font-size: 28px;
  font-weight: 700;
  color: var(--brand-dark);
}

.loading-state, .empty-state {
  text-align: center;
  padding: 48px 24px;
  color: var(--ink-soft);
  font-size: 15px;
}

.table-wrapper { 
  overflow-x: auto;
  border-radius: var(--radius-md);
  border: 1px solid var(--line);
}

.table { 
  width: 100%; 
  border-collapse: collapse;
  background: var(--surface);
}

.table th {
  text-align: right;
  font-size: 12px;
  color: var(--ink-soft);
  font-weight: 700;
  padding: 14px 16px;
  border-bottom: 2px solid var(--line);
  background: var(--canvas);
  text-transform: uppercase;
  letter-spacing: 0.02em;
}

.table td {
  padding: 14px 16px;
  font-size: 14px;
  border-bottom: 1px solid var(--line);
  color: var(--ink);
}

.table tr:last-child td { border-bottom: none; }
.table tr:hover td {
  background: var(--canvas);
}

.devpro-branding {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 24px;
  margin-top: 20px;
  background: linear-gradient(135deg, var(--brand-tint) 0%, rgba(31, 111, 92, 0.05) 100%);
  border-radius: var(--radius-md);
  border: 1px solid var(--line);
  animation: fadeIn 0.5s ease;
}

.devpro-logo-img {
  width: 80px;
  height: auto;
  border-radius: 12px;
  margin-bottom: 12px;
  box-shadow: var(--shadow-sm);
  transition: transform 0.3s ease, box-shadow 0.3s ease;
  display: block;
}

.devpro-logo-img:hover {
  transform: scale(1.05);
  box-shadow: var(--shadow-md);
}

.devpro-text {
  font-size: 13px;
  font-weight: 600;
  color: var(--ink-soft);
  letter-spacing: 0.5px;
  text-transform: uppercase;
  margin: 0;
}

@media (max-width: 480px) {
  .devpro-logo-img {
    width: 60px;
  }
}
</style>