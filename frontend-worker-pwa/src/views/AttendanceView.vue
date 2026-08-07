<template>
  <div class="attendance-view view">
    <!-- Page Header -->
    <div class="page-header">
      <div class="header-icon">
        <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="12" cy="12" r="10"/>
          <polyline points="12 6 12 12 16 14"/>
        </svg>
      </div>
      <div>
        <h1>{{ t('attendance') }}</h1>
        <p>{{ t('attendance_subtitle') }}</p>
      </div>
      
      <!-- Theme & Language Toolbar -->
      <div class="settings-toolbar">
        <button 
          class="settings-btn" 
          :class="{ active: isDarkMode }"
          @click="toggleTheme"
          title="الوضع الليلي" 
          aria-label="تبديل الوضع الليلي"
        >
          <svg v-if="!isDarkMode" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79Z"/>
          </svg>
          <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="5"/>
            <path d="M12 1v2M12 21v2M4.22 4.22l1.42 1.42M18.36 18.36l1.42 1.42M1 12h2M21 12h2M4.22 19.78l1.42-1.42M18.36 5.64l1.42-1.42"/>
          </svg>
        </button>
        
        <!-- Language Buttons -->
        <button 
          v-for="lang in languages" 
          :key="lang.code"
          class="settings-btn settings-btn--lang"
          :class="{ active: currentLang === lang.code }"
          @click="changeLanguage(lang.code)"
          :title="lang.name"
          :aria-label="lang.name"
        >
          {{ lang.label }}
        </button>
      </div>
    </div>

    <!-- Summary Stats -->
    <div class="summary-row stagger">
      <div class="summary-card">
        <div class="s-icon">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10"/>
            <polyline points="12 6 12 12 16 14"/>
          </svg>
        </div>
        <div class="s-value">{{ todayHours.toFixed(1) }}h</div>
        <div class="s-label">{{ t('hours_today') }}</div>
      </div>
      <div class="summary-card">
        <div class="s-icon">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <rect x="3" y="4" width="18" height="18" rx="2"/>
            <path d="M16 2v4M8 2v4M3 10h18"/>
          </svg>
        </div>
        <div class="s-value">{{ weekHours.toFixed(1) }}h</div>
        <div class="s-label">{{ t('hours_week') }}</div>
      </div>
      <div class="summary-card">
        <div class="s-icon">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M22 12h-4l-3 9L9 3l-3 9H2"/>
          </svg>
        </div>
        <div class="s-value">{{ monthHours.toFixed(1) }}h</div>
        <div class="s-label">{{ t('hours_month') }}</div>
      </div>
    </div>

    <!-- Worksite Selection -->
    <div v-if="!isWorking" class="section">
      <div class="section-header">
        <div class="section-label">{{ t('select_worksite') }}</div>
        <button class="refresh-btn" @click="fetchWorksites" :disabled="loadingWorksites">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M23 4v6h-6"/>
            <path d="M20.49 15a9 9 0 1 1-2.12 9 9 0 0 1-2.12-2.12M20.49 9a9 9 0 0 1-2.12-2.12"/>
            <path d="M13.05 13.05a9 9 0 0 1-2.12-2.12"/>
            <path d="M13.05 6.05a9 9 0 0 1 2.12-2.12"/>
            <path d="M16 16l4 4m0-4l-4 4"/>
          </svg>
        </button>
      </div>
      <div v-if="loadingWorksites" class="loading">{{ t('loading') }}</div>
      <div v-else-if="availableWorksites.length === 0" class="empty">
        <p>{{ t('no_worksites') }}</p>
      </div>
      <div v-else class="site-grid">
        <button
          v-for="site in availableWorksites"
          :key="site.id"
          class="site-card"
          :class="{ selected: selectedWorksiteId === site.id }"
          @click="selectWorksite(site)"
        >
          <div class="site-icon">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M6 22V4a2 2 0 0 1 2-2h10a2 2 0 0 1 2 2v18Z"/>
              <path d="M6 12H4a2 2 0 0 0-2 2v6a2 2 0 0 0 2 2h2"/>
              <path d="M18 9h2a2 2 0 0 1 2 2v9a2 2 0 0 1-2 2h-2"/>
              <path d="M10 6h4M10 10h4M10 14h4"/>
            </svg>
          </div>
          <div class="site-info">
            <div class="site-name">{{ site.name }}</div>
            <div class="site-address">{{ site.address || t('address') }}</div>
            <div class="site-distance">
              ⭕ {{ site.radius_meters }} {{ t('meter') }}
            </div>
          </div>
          <div v-if="selectedWorksiteId === site.id" class="check">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#fff" stroke-width="3">
              <polyline points="20 6 9 17 4 12"/>
            </svg>
          </div>
        </button>
      </div>
    </div>

    <!-- Active Worksite Card (When Working) -->
    <div v-if="isWorking" class="card active-worksite-card">
      <h3>📍 {{ t('current_worksite') }}</h3>
      <div class="active-worksite-info">
        <div class="active-worksite-name">{{ worksiteName }}</div>
        <div class="active-worksite-address">{{ selectedWorksite?.address || '' }}</div>
        <div class="active-worksite-status">✅ {{ t('working_at_this_location') }}</div>
      </div>
    </div>

    <!-- Navigation Card (When Worksite Selected) -->
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

    <!-- Timer Section (When Working) -->
    <div v-if="isWorking" class="shift-timer">
      <div class="timer-header">
        <div class="timer-date">{{ formatDate(new Date()) }}</div>
        <div class="timer-day">{{ getDayName(new Date()) }}</div>
      </div>
      
      <div class="ring-wrap">
        <svg width="140" height="140" viewBox="0 0 140 140">
          <circle class="ring-bg" cx="70" cy="70" r="60" fill="none" stroke-width="6"/>
          <circle class="ring-fg" cx="70" cy="70" r="60" fill="none" stroke-width="6"/>
        </svg>
        <div class="ring-center">
          <div class="ring-time">{{ elapsedTime }}</div>
          <div class="ring-label">{{ t('in_progress') }}</div>
        </div>
      </div>
      
      <div class="shift-site">
        <span class="pulse-dot"></span>
        {{ worksiteName }}
      </div>
      
      <button 
        class="btn btn--secondary btn--block" 
        @click="getLocationWithAddress" 
        :disabled="gettingLocation"
      >
        <svg v-if="gettingLocation" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="animate-spin">
          <path d="M21 12a9 9 0 1 1-6.219-8.56"/>
        </svg>
        <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M20 10c0 6-8 12-8 12s-8-6-8-12a8 8 0 0 1 16 0Z"/>
          <circle cx="12" cy="10" r="3"/>
        </svg>
        {{ t('select_location') }}
      </button>
      
      <button 
        class="btn btn--primary btn--block" 
        @click="checkOut"
        :disabled="!userLocation || !withinRange || checkingOut"
      >
        <svg v-if="checkingOut" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="animate-spin">
          <path d="M21 12a9 9 0 1 1-6.219-8.56"/>
        </svg>
        <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <rect x="6" y="4" width="4" height="16" rx="1"/>
          <rect x="14" y="4" width="4" height="16" rx="1"/>
        </svg>
        {{ t('check_out') }}
      </button>
    </div>

    <!-- Location Status -->
    <div v-if="selectedWorksiteId" class="location-card">
      <div class="loc-row">
        <div class="loc-icon" :class="withinRange ? 'success' : 'error'">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M20 10c0 6-8 12-8 12s-8-6-8-12a8 8 0 0 1 16 0Z"/>
            <circle cx="12" cy="10" r="3"/>
          </svg>
        </div>
        <div style="flex:1">
          <div class="loc-title">{{ withinRange ? t('inside_range') : t('outside_range') }}</div>
          <div class="loc-distance">{{ t('distance') }}: {{ formatDistance(distance) }}</div>
        </div>
        <span class="badge" :class="withinRange ? 'badge--success' : 'badge--error'">
          {{ withinRange ? t('active') : t('inactive') }}
        </span>
      </div>
      <div v-if="userLocation" class="loc-meta">
        <div class="loc-coords mono">{{ userLocation.lat.toFixed(6) }}° N, {{ userLocation.lng.toFixed(6) }}° E</div>
        <div v-if="locationAddress" class="loc-address">{{ locationAddress }}</div>
      </div>
    </div>

    <!-- Warning Messages -->
    <div v-if="!hasClickedLocation && isWorking" class="alert alert--warning">
      <span>⚠️</span>
      <span>📍 {{ t('select_location') }} {{ t('before_checkout') }}</span>
    </div>

    <div v-if="hasClickedLocation && !withinRange && isWorking" class="alert alert--error">
      <span>🚫</span>
      <span>❌ {{ t('outside_range') }}! {{ t('distance') }}: {{ formatDistance(distance) }}</span>
    </div>

    <div v-if="hasClickedLocation && withinRange && isWorking" class="alert alert--success">
      <span>✅</span>
      <span>{{ t('inside_range') }} - {{ t('can_checkout') }}</span>
    </div>

    <!-- Action Buttons -->
    <div v-if="!isWorking" class="actions">
      <button 
        class="btn btn--primary btn--block" 
        :disabled="!selectedWorksiteId || !withinRange || !userLocation || checkingIn" 
        @click="checkIn"
      >
        <svg v-if="checkingIn" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="animate-spin">
          <path d="M21 12a9 9 0 1 1-6.219-8.56"/>
        </svg>
        <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M15 3h4a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2h-4"/>
          <polyline points="10 17 15 12 10 7"/>
          <line x1="15" y1="12" x2="3" y2="12"/>
        </svg>
        {{ t('check_in') }}
      </button>
      
      <button 
        class="btn btn--secondary btn--block" 
        @click="getLocationWithAddress" 
        :disabled="gettingLocation"
      >
        <svg v-if="gettingLocation" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="animate-spin">
          <path d="M21 12a9 9 0 1 1-6.219-8.56"/>
        </svg>
        <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M20 10c0 6-8 12-8 12s-8-6-8-12a8 8 0 0 1 16 0Z"/>
          <circle cx="12" cy="10" r="3"/>
        </svg>
        {{ t('select_location') }}
      </button>
      
      <!-- User Location Address -->
      <div v-if="locationAddress" class="user-location-address">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M20 10c0 6-8 12-8 12s-8-6-8-12a8 8 0 0 1 16 0Z"/>
          <circle cx="12" cy="10" r="3"/>
        </svg>
        <span class="address-text">{{ locationAddress }}</span>
      </div>
    </div>

    <!-- Error/Success Messages -->
    <div v-if="error" class="alert alert--error">{{ error }}</div>
    <div v-if="success" class="alert alert--success">{{ success }}</div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useI18n } from '../services/i18n'
import api from '../services/api'

const { t, currentLang, setLang } = useI18n()

// Theme & Language
const isDarkMode = ref(false)

const languages = [
  { code: 'ar', name: 'العربية', dir: 'rtl', label: 'AR' },
  { code: 'he', name: 'עברית', dir: 'rtl', label: 'HE' },
  { code: 'en', name: 'English', dir: 'ltr', label: 'EN' }
]

onMounted(() => {
  const savedTheme = localStorage.getItem('worktrack_theme')
  if (savedTheme === 'dark') {
    isDarkMode.value = true
    document.documentElement.setAttribute('data-theme', 'dark')
  }
})

function toggleTheme() {
  isDarkMode.value = !isDarkMode.value
  const theme = isDarkMode.value ? 'dark' : 'light'
  document.documentElement.setAttribute('data-theme', theme)
  localStorage.setItem('worktrack_theme', theme)
}

function changeLanguage(code) {
  setLang(code)
  const lang = languages.find(l => l.code === code)
  if (lang) {
    document.documentElement.setAttribute('dir', lang.dir)
    document.documentElement.setAttribute('data-lang', lang.code)
  }
}

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
const error = ref(''), success = ref('')
let timerInterval = null, locationInterval = null, attendanceCheckInterval = null

const withinRange = computed(() => distance.value <= radius.value)
const elapsedTime = computed(() => {
  const s = elapsedSeconds.value
  return `${String(Math.floor(s/3600)).padStart(2,'0')}:${String(Math.floor((s%3600)/60)).padStart(2,'0')}:${String(Math.floor(s%60)).padStart(2,'0')}`
})

function formatDistance(meters) {
  if (!meters) return '0 ' + t('meters')
  if (meters >= 1000) {
    return (meters / 1000).toFixed(2) + ' ' + t('kilometers')
  }
  return Math.round(meters) + ' ' + t('meters')
}

// روابط الملاحة
function getWazeUrl(site) {
  if (!site) return '#'
  return `https://www.waze.com/ul?ll=${site.latitude},${site.longitude}&navigate=yes`
}

function getGoogleMapsUrl(site) {
  if (!site) return '#'
  return `https://www.google.com/maps/dir/?api=1&destination=${site.latitude},${site.longitude}`
}

async function fetchWorksites() {
  loadingWorksites.value = true
  try {
    const { data } = await api.get('/worksites/available')
    availableWorksites.value = data || []
    
    // إذا كان هناك نقطة عمل محددة مسبقاً، تحقق من أنها لا تزال متاحة
    if (selectedWorksiteId.value && !availableWorksites.value.find(w => w.id === selectedWorksiteId.value)) {
      // نقطة العمل المحددة لم تعد متاحة
      selectedWorksiteId.value = null
      selectedWorksite.value = null
      worksiteName.value = ''
      if (isWorking.value) {
        // إذا كان الموظف يعمل، طلب منه تسجيل الخروج
        await checkOut()
      }
    }
  } catch(e) { 
    console.error(e) 
  } finally { 
    loadingWorksites.value = false 
  }
}

// التحقق من وجود وردية نشطة (للتحقق من تسجيل الدخول التلقائي من قبل المدير)
async function checkActiveAttendance() {
  try {
    const { data } = await api.get('/attendance/current')
    if (data && data.has_active) {
      // وجدنا وردية نشطة - تم تسجيل الدخول تلقائياً
      attendanceId.value = data.attendance_id
      isWorking.value = true
      
      // إذا كانت هناك نقطة عمل مرتبطة، حددها
      if (data.worksite_id) {
        const worksite = availableWorksites.value.find(w => w.id === data.worksite_id)
        if (worksite) {
          selectedWorksiteId.value = worksite.id
          selectedWorksite.value = worksite
          worksiteName.value = worksite.name
          radius.value = worksite.radius_meters
        }
      }
      
      // حساب الوقت المنقضي من وقت تسجيل الدخول
      if (data.elapsed_seconds) {
        elapsedSeconds.value = data.elapsed_seconds
      } else if (data.check_in_time) {
        const checkInTime = new Date(data.check_in_time)
        const now = new Date()
        elapsedSeconds.value = Math.floor((now - checkInTime) / 1000)
      }
      
      // بدء المؤقت وتتبع الموقع
      timerInterval = setInterval(() => elapsedSeconds.value++, 1000)
      startLocationTracking()
      
      success.value = t('auto_checkin_detected')
    }
  } catch(e) {
    // لا توجد وردية نشطة - هذا طبيعي
    console.log('لا توجد وردية نشطة')
  }
}

function selectWorksite(site) {
  selectedWorksiteId.value = site.id
  selectedWorksite.value = site
  worksiteName.value = site.name
  radius.value = site.radius_meters
  
  if (userLocation.value) {
    calculateDistance(userLocation.value.lat, userLocation.value.lng)
  }
}

// جلب عنوان الموقع
async function getAddressFromCoords(lat, lng) {
  try {
    const url = `https://api.geoapify.com/v1/geocode/reverse?lat=${lat}&lon=${lng}&apiKey=a6a3b5fec1cd4b1c99daaf6decab855f&lang=${currentLang.value}`
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
}

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

    if (selectedWorksite.value) {
      calculateDistance(lat, lng)
    }

    const address = await getAddressFromCoords(lat, lng)
    locationAddress.value = address

    // إرسال الموقع فوراً للتتبع اللحظي
    if (isWorking.value) {
      try {
        await api.post('/location/update', {
          latitude: lat,
          longitude: lng,
          accuracy: pos.coords.accuracy || 0
        })
      } catch (e) {
        console.error('فشل إرسال الموقع الفوري:', e)
      }
    }

    success.value = t('location_determined_success')
  } catch (e) {
    error.value = t('location_determine_failed') + ' ' + (e.message || t('error'))
  } finally {
    gettingLocation.value = false
  }
}

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

    hasClickedLocation.value = false
  } catch(e) {
    console.error('❌ فشل CheckIn:', e.response?.data)
    const errData = e.response?.data
    if (errData?.geofence) {
      distance.value = errData.geofence.distance_meters || distance.value
      error.value = `${t('outside_range_actual_distance')} ${formatDistance(distance.value)} (${t('allowed_range')}: ${radius.value} ${t('meter')})`
    } else {
      error.value = errData?.error || t('checkin_failed')
    }
  } finally {
    checkingIn.value = false
  }
}

// بدء تتبع الموقع
function startLocationTracking() {
  if (!navigator.geolocation) return
  
  locationInterval = setInterval(async () => {
    if (!isWorking.value) {
      stopLocationTracking()
      return
    }
    
    try {
      const pos = await new Promise((resolve, reject) => {
        navigator.geolocation.getCurrentPosition(resolve, reject, {
          enableHighAccuracy: true,
          timeout: 10000,
          maximumAge: 5000
        })
      })
      
      userLocation.value = { lat: pos.coords.latitude, lng: pos.coords.longitude }
      
      if (selectedWorksite.value) {
        calculateDistance(pos.coords.latitude, pos.coords.longitude)
      }
      
      // إرسال الموقع للخادم
      await api.post('/location/update', {
        latitude: pos.coords.latitude,
        longitude: pos.coords.longitude,
        accuracy: pos.coords.accuracy || 0
      })
    } catch (e) {
      console.error('فشل تتبع الموقع:', e)
    }
  }, 30000) // كل 30 ثانية
}

// إيقاف تتبع الموقع
function stopLocationTracking() {
  if (locationInterval) {
    clearInterval(locationInterval)
    locationInterval = null
  }
}

async function checkOut() {
  if (!attendanceId.value) {
    error.value = t('no_active_shift')
    return
  }

  if (!hasClickedLocation.value) {
    error.value = t('click_select_location_before_checkout')
    return
  }

  if (!withinRange.value) {
    error.value = `${t('outside_range_distance')} ${formatDistance(distance.value)}`
    return
  }

  checkingOut.value = true
  error.value = ''
  success.value = ''

  try {
    const payload = {
      attendance_id: attendanceId.value,
      latitude: userLocation.value.lat,
      longitude: userLocation.value.lng
    }

    await api.post('/attendance/check-out', payload)
    
    success.value = t('checkout_success')
    isWorking.value = false
    clearInterval(timerInterval)
    stopLocationTracking()
    attendanceId.value = null
    elapsedSeconds.value = 0
    await fetchMyStats() // Update stats after checkout
  } catch(e) {
    console.error('❌ فشل CheckOut:', e.response?.data)
    error.value = e.response?.data?.error || t('checkout_failed')
  } finally {
    checkingOut.value = false
  }
}

async function fetchMyStats() {
  try {
    const { data } = await api.get('/attendance/summary')
    todayHours.value = data.today_hours || 0
    weekHours.value = data.week_hours || 0
    monthHours.value = data.month_hours || 0
  } catch(e) {
    console.error('Failed to fetch stats:', e)
    // Don't fail the app if stats endpoint doesn't exist
    todayHours.value = 0
    weekHours.value = 0
    monthHours.value = 0
  }
}

// دالة بديلة لاستخدامها أيضاً (مثل الملف الاحتياطي)
async function fetchSummary() {
  try {
    const { data } = await api.get('/attendance/summary')
    todayHours.value = data.today_hours || 0
    weekHours.value = data.week_hours || 0
    monthHours.value = data.month_hours || 0
  } catch(e) {
    console.error('Failed to fetch summary:', e)
  }
}

function formatDate(dateStr) {
  const date = new Date(dateStr)
  return date.toLocaleDateString(currentLang.value === 'ar' ? 'ar-SA' : 'en-US')
}

function getDayName(date) {
  const days = {
    'ar': ['الأحد', 'الاثنين', 'الثلاثاء', 'الأربعاء', 'الخميس', 'الجمعة', 'السبت'],
    'he': ['ראשון', 'שני', 'שלישי', 'רביעי', 'חמישי', 'שישי', 'שבת'],
    'en': ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday']
  }
  const dayIndex = date.getDay()
  return days[currentLang.value]?.[dayIndex] || days['en'][dayIndex]
}

function formatTime(dateStr) {
  const date = new Date(dateStr)
  return date.toLocaleTimeString(currentLang.value === 'ar' ? 'ar-SA' : 'en-US', { hour: '2-digit', minute: '2-digit' })
}

onMounted(() => {
  fetchWorksites()
  checkActiveAttendance() // التحقق من وجود وردية نشطة (تسجيل دخول تلقائي من المدير)
  
  // تحديث دوري كل 30 ثانية للتحقق من حالة الحضور
  const checkInterval = setInterval(() => {
    if (!isWorking.value) {
      checkActiveAttendance()
    }
  }, 30000)
  
  // تخزين الـ interval للتنظيف عند unmount
  attendanceCheckInterval = checkInterval
  fetchMyStats() // Fetch initial stats
})

onUnmounted(() => {
  if (timerInterval) clearInterval(timerInterval)
  if (attendanceCheckInterval) clearInterval(attendanceCheckInterval)
  stopLocationTracking()
})
</script>

<style scoped>
.attendance-view {
  padding: var(--space-4);
}

/* Settings Toolbar in Header */
.settings-toolbar {
  display: flex;
  gap: 4px;
  align-items: center;
  padding: 4px;
  background: var(--surface-elevated);
  border: 1px solid var(--border);
  border-radius: 9999px;
  box-shadow: var(--shadow-xs);
}

.settings-btn {
  border: none;
  background: transparent;
  cursor: pointer;
  border-radius: 50%;
  width: 28px;
  height: 28px;
  min-width: 28px;
  min-height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-secondary);
  transition: var(--transition-base) ease;
  font-size: 10px;
  font-weight: var(--font-semibold);
  padding: 0;
}

.settings-btn--lang {
  font-size: 9px;
  letter-spacing: -0.2px;
}

.settings-btn:hover {
  background: var(--gray-100);
  color: var(--text-primary);
}

.settings-btn.active {
  background: var(--primary-500);
  color: #fff;
  box-shadow: 0 2px 6px rgba(20, 184, 166, 0.25);
}

.summary-row {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: var(--space-3);
  margin-bottom: var(--space-5);
}

.summary-card {
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
  padding: var(--space-3);
  text-align: center;
}

.s-icon {
  width: 28px;
  height: 28px;
  border-radius: var(--radius-md);
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto var(--space-2);
  background: var(--primary-100);
  color: var(--primary-700);
}

.s-value {
  font-size: var(--text-lg);
  font-weight: var(--font-semibold);
  color: var(--text-primary);
}

.s-label {
  font-size: var(--text-xs);
  color: var(--text-tertiary);
}

.section {
  margin-bottom: var(--space-5);
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--space-3);
}

.refresh-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border-radius: var(--radius-md);
  border: 1px solid var(--border);
  background: var(--surface);
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.2s;
}

.refresh-btn:hover:not(:disabled) {
  background: var(--primary-50);
  color: var(--primary-600);
  border-color: var(--primary-200);
}

.refresh-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.refresh-btn svg {
  transition: transform 0.2s;
}

.refresh-btn:hover:not(:disabled) svg {
  transform: rotate(180deg);
}

.site-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: var(--space-3);
}

.site-card {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  cursor: pointer;
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
  padding: var(--space-3);
  background: var(--surface);
  transition: var(--transition-base) ease;
}

.site-card:hover {
  border-color: var(--border-strong);
}

.site-card.selected {
  border-color: var(--primary-500);
  background: var(--primary-50);
}

[data-theme="dark"] .site-card.selected {
  background: rgba(99, 102, 241, 0.12);
}

.site-icon {
  width: 36px;
  height: 36px;
  border-radius: var(--radius-md);
  background: var(--gray-100);
  color: var(--text-secondary);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.site-card.selected .site-icon {
  background: var(--primary-500);
  color: #fff;
}

.site-info {
  flex: 1;
  min-width: 0;
}

.site-name {
  font-weight: var(--font-medium);
  font-size: var(--text-sm);
  color: var(--text-primary);
}

.site-address {
  font-size: var(--text-xs);
  color: var(--text-tertiary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.site-distance {
  font-size: var(--text-xs);
  color: var(--primary-600);
  font-weight: var(--font-medium);
  margin-top: 2px;
}

.check {
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: var(--primary-500);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.location-card {
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
  padding: var(--space-4);
  margin-bottom: var(--space-5);
}

.loc-row {
  display: flex;
  align-items: center;
  gap: var(--space-3);
}

.loc-icon {
  width: 36px;
  height: 36px;
  border-radius: var(--radius-md);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.loc-icon.success {
  background: var(--success-100);
  color: var(--success-700);
}

.loc-icon.error {
  background: var(--error-100);
  color: var(--error-700);
}

.loc-title {
  font-weight: var(--font-medium);
  font-size: var(--text-sm);
  color: var(--text-primary);
}

.loc-distance {
  font-size: var(--text-xs);
  color: var(--text-tertiary);
}

.loc-meta {
  display: flex;
  flex-direction: column;
  gap: var(--space-1);
  margin-top: var(--space-3);
}

.loc-coords {
  font-size: var(--text-xs);
  color: var(--text-tertiary);
}

.loc-address {
  font-size: var(--text-sm);
  color: var(--text-secondary);
}

.actions {
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
  margin-bottom: var(--space-5);
}

.user-location-address {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-3);
  background: linear-gradient(90deg, var(--primary-50) 0%, var(--surface-elevated) 20%);
  border-radius: var(--radius-md);
  font-size: var(--text-sm);
  color: var(--text-secondary);
  border: 1px solid var(--primary-200);
  border-inline-start: 3px solid var(--primary-500);
}

[data-theme="dark"] .user-location-address {
  background: linear-gradient(90deg, rgba(20, 184, 166, 0.12) 0%, var(--surface-elevated) 20%);
  border-color: var(--primary-700);
}

.address-text {
  font-weight: var(--font-medium);
  color: var(--primary-700);
}

[data-theme="dark"] .address-text {
  color: var(--primary-400);
}

.alert {
  padding: var(--space-3);
  border-radius: var(--radius-md);
  font-size: var(--text-sm);
  margin-bottom: var(--space-4);
}

.alert--error {
  background: var(--error-50);
  color: var(--error-700);
  border: 1px solid var(--error-100);
}

.alert--success {
  background: var(--success-50);
  color: var(--success-700);
  border: 1px solid var(--success-100);
}

.loading, .empty {
  text-align: center;
  padding: var(--space-4);
  color: var(--text-tertiary);
  font-size: var(--text-sm);
}

/* Modal */
.modal-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: var(--z-modal);
  padding: var(--space-4);
}

.modal {
  background: var(--surface);
  border-radius: var(--radius-lg);
  width: 100%;
  max-width: 500px;
  max-height: 80vh;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--space-4);
  border-bottom: 1px solid var(--border);
}

.modal-header h3 {
  font-size: var(--text-base);
  font-weight: var(--font-semibold);
  margin: 0;
}

.modal-close {
  border: none;
  background: transparent;
  cursor: pointer;
  padding: var(--space-2);
  border-radius: var(--radius-md);
  color: var(--text-secondary);
  display: flex;
  align-items: center;
  justify-content: center;
}

.modal-close:hover {
  background: var(--gray-100);
  color: var(--text-primary);
}

.modal-body {
  padding: var(--space-4);
  overflow-y: auto;
}

.filters {
  display: flex;
  gap: var(--space-3);
  margin-bottom: var(--space-4);
}

.form-select {
  flex: 1;
  padding: var(--space-2);
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  background: var(--surface);
  color: var(--text-primary);
  font-size: var(--text-sm);
}

.history-list {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.history-item {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  padding: var(--space-3);
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  background: var(--surface);
}

.history-date {
  font-size: var(--text-xs);
  color: var(--text-tertiary);
  min-width: 80px;
}

.history-details {
  flex: 1;
}

.history-worksite {
  font-size: var(--text-sm);
  font-weight: var(--font-medium);
  color: var(--text-primary);
}

.history-times {
  font-size: var(--text-xs);
  color: var(--text-tertiary);
}

.history-hours {
  font-size: var(--text-sm);
  font-weight: var(--font-semibold);
  color: var(--primary-600);
}

.animate-spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

/* Active Worksite Card */
.active-worksite-card {
  background: linear-gradient(135deg, var(--primary-50) 0%, var(--surface) 100%);
  border: 2px solid var(--primary-200);
}

.active-worksite-info {
  padding: var(--space-3);
}

.active-worksite-name {
  font-weight: var(--font-semibold);
  font-size: var(--text-lg);
  color: var(--text-primary);
  margin-bottom: var(--space-1);
}

.active-worksite-address {
  font-size: var(--text-sm);
  color: var(--text-secondary);
  margin-bottom: var(--space-2);
}

.active-worksite-status {
  font-size: var(--text-sm);
  color: var(--primary-600);
  font-weight: var(--font-medium);
}

/* Navigation Card */
.navigation-card {
  background: linear-gradient(135deg, var(--surface) 0%, var(--surface-elevated) 100%);
}

.nav-info {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  margin-bottom: var(--space-3);
}

.nav-icon {
  font-size: 24px;
}

.nav-title {
  font-weight: var(--font-semibold);
  font-size: var(--text-base);
  color: var(--text-primary);
  margin: 0;
}

.nav-address {
  font-size: var(--text-sm);
  color: var(--text-secondary);
  margin: 0;
}

.nav-buttons {
  display: flex;
  gap: var(--space-2);
}

.btn--waze {
  background: linear-gradient(135deg, #33ccff 0%, #0099cc 100%);
  color: white;
  flex: 1;
}

.btn--google {
  background: linear-gradient(135deg, #4285f4 0%, #34a853 100%);
  color: white;
  flex: 1;
}

/* SHIFT TIMER - Professional Design */
.shift-timer {
  position: relative;
  border-radius: var(--radius-lg);
  margin-bottom: var(--space-6);
  background: var(--surface);
  border: 1px solid var(--border);
  padding: var(--space-6);
  box-shadow: var(--shadow-sm);
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  animation: slideUp 0.4s ease-out;
}

.timer-header {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  margin-bottom: var(--space-4);
}

.timer-date {
  font-size: var(--text-sm);
  color: var(--text-tertiary);
  font-weight: var(--font-medium);
}

.timer-day {
  font-size: var(--text-sm);
  color: var(--text-primary);
  font-weight: var(--font-semibold);
}

.ring-wrap {
  position: relative;
  width: 140px;
  height: 140px;
  margin-bottom: var(--space-4);
}

.ring-wrap svg {
  transform: rotate(-90deg);
}

.ring-bg {
  stroke: var(--border);
}

.ring-fg {
  stroke: var(--primary-500);
  stroke-linecap: round;
  stroke-dasharray: 440;
  stroke-dashoffset: 440;
  animation: dashMove 2.2s ease-out forwards 0.3s;
}

.ring-center {
  position: absolute;
  inset: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.ring-time {
  font-size: var(--text-2xl);
  font-weight: var(--font-semibold);
  font-variant-numeric: tabular-nums;
  letter-spacing: -0.01em;
  color: var(--text-primary);
}

.ring-label {
  font-size: var(--text-xs);
  color: var(--text-tertiary);
  margin-top: 2px;
}

.shift-site {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: var(--text-sm);
  color: var(--text-secondary);
  margin-bottom: var(--space-4);
}

.pulse-dot {
  width: 5px;
  height: 5px;
  border-radius: 50%;
  background: var(--success-500);
  box-shadow: 0 0 0 2px rgba(16, 185, 129, 0.2);
  animation: pulseSoft 2s ease-in-out infinite;
}

@media (min-width: 768px) {
  .site-grid {
    grid-template-columns: 1fr 1fr;
  }
  
  .modal {
    max-width: 600px;
  }
}
</style>
