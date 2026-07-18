<template>
  <div class="map-container">
    <l-map
      ref="mapRef"
      :key="isDarkMode ? 'dark' : 'light'"
      :zoom="zoom"
      @update:zoom="updateZoom"
      :center="center"
      :options="{ attributionControl: true, zoomControl: true, zoomSnap: 0.5 }"
      :style="{ height: height + 'px', width: '100%' }"
      class="wolt-map"
    >
      <l-tile-layer
        :url="mapTileUrl"
        layer-type="base"
        :name="mapLayerName"
        :attribution="mapAttribution"
      />
      
      <!-- نقاط العمل -->
      <l-marker
        v-for="site in worksites"
        :key="site.id"
        :lat-lng="[site.latitude, site.longitude]"
      >
        <l-icon>
          <div class="worksite-marker">
            <div class="worksite-dot"></div>
            <div class="worksite-ring"></div>
          </div>
        </l-icon>
        <l-popup>
          <div class="popup-content">
            <h4>{{ site.name }}</h4>
            <p v-if="site.address">📍 {{ site.address }}</p>
            <p>⭕ النطاق: {{ site.radius_meters }} متر</p>
            <p>👥 الموظفين: {{ getEmployeeCount(site.id) }}</p>
          </div>
        </l-popup>
      </l-marker>

      <!-- الموظفين -->
      <l-marker
        v-for="emp in employees"
        :key="emp.id"
        :lat-lng="[emp.latitude, emp.longitude]"
      >
        <l-icon>
          <div class="employee-marker" :class="emp.status">
            <div class="employee-pulse" :class="emp.status"></div>
            <div class="employee-dot" :class="emp.status">
              {{ emp.full_name.slice(0, 1) }}
            </div>
          </div>
        </l-icon>
        <l-popup>
          <div class="popup-content employee-popup">
            <h4>{{ emp.full_name }}</h4>
            <p>📍 {{ emp.worksite.name }}</p>
            <p>📏 المسافة: {{ formatDistance(emp.worksite.distance) }}</p>
            <p>⏱️ {{ emp.hours_worked.toFixed(1) }} ساعة</p>
            <p>
              <span class="badge" :class="emp.status === 'inside' ? 'badge--in' : 'badge--out'">
                {{ emp.status_text }}
              </span>
            </p>
            <button
              class="btn btn--sm btn--primary"
              @click="handleShowDetails(emp)"
            >
              عرض التفاصيل
            </button>
          </div>
        </l-popup>
      </l-marker>
    </l-map>

    <div v-if="!employees || employees.length === 0" class="map-overlay">
      <div class="overlay-content">
        <span class="overlay-icon">🗺️</span>
        <h3>لا يوجد موظفين نشطين</h3>
        <p>سيظهر هنا الموظفون الذين بدأوا الدوام</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { LMap, LTileLayer, LMarker, LPopup, LIcon } from '@vue-leaflet/vue-leaflet'
import 'leaflet/dist/leaflet.css'

const props = defineProps({
  employees: { type: Array, default: () => [] },
  worksites: { type: Array, default: () => [] },
  center: { type: Array, default: () => [31.5, 34.8] },
  zoom: { type: Number, default: 7 },
  height: { type: Number, default: 400 }  // ✅ Number وليس String
})

// ✅ تعريف emit بشكل صحيح
const emit = defineEmits(['update:zoom', 'showDetails'])

const mapRef = ref(null)
let watchId = null
let observer = null

// ✅ اكتشاف الوضع الداكن مع مراقبة التغييرات
const isDarkMode = ref(document.documentElement.getAttribute('data-theme') === 'dark')

// ✅ مراقبة تغييرات data-theme
const observeThemeChanges = () => {
  observer = new MutationObserver((mutations) => {
    mutations.forEach((mutation) => {
      if (mutation.type === 'attributes' && mutation.attributeName === 'data-theme') {
        isDarkMode.value = document.documentElement.getAttribute('data-theme') === 'dark'
      }
    })
  })

  observer.observe(document.documentElement, {
    attributes: true,
    attributeFilter: ['data-theme']
  })
}

// ✅ تحديد URL الخريطة حسب الوضع - Wolt-style design
const mapTileUrl = computed(() => {
  const isDark = document.documentElement.getAttribute('data-theme') === 'dark'
  if (isDark) {
    // خريطة داكنة أنيقة من CartoDB - تشبه Wolt
    return 'https://{s}.basemaps.cartocdn.com/dark_all/{z}/{x}/{y}{r}.png'
  } else {
    // خريطة فاتحة أنيقة من CartoDB - تشبه Wolt
    return 'https://{s}.basemaps.cartocdn.com/light_all/{z}/{x}/{y}{r}.png'
  }
})

// ✅ تحديد Attribution حسب الوضع - Wolt-style minimal attribution
const mapAttribution = computed(() => {
  return '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors &copy; <a href="https://carto.com/attributions">CARTO</a>'
})

// ✅ اسم طبقة الخريطة
const mapLayerName = computed(() => {
  const isDark = document.documentElement.getAttribute('data-theme') === 'dark'
  return isDark ? 'CartoDB Dark' : 'CartoDB Light'
})

function updateZoom(newZoom) {
  emit('update:zoom', newZoom)
}

// ✅ دالة handleShowDetails لتوصيل الحدث
function handleShowDetails(employee) {
  console.log('📋 عرض تفاصيل الموظف:', employee.full_name)
  emit('showDetails', employee)
}

function getEmployeeCount(worksiteId) {
  return props.employees.filter(e => e.worksite?.id === worksiteId).length
}

function formatDistance(meters) {
  if (!meters) return '0 م'
  if (meters >= 1000) {
    return (meters / 1000).toFixed(2) + ' كيلومتر'
  }
  return Math.round(meters) + ' متر'
}

function getUserLocation() {
  if (!('geolocation' in navigator)) return

  watchId = navigator.geolocation.watchPosition(
    (pos) => {
      // يمكن استخدام موقع المدير إذا أردت
    },
    (err) => console.error('فشل تحديد الموقع:', err),
    { enableHighAccuracy: true, timeout: 10000, maximumAge: 30000 }
  )
}

onMounted(() => {
  getUserLocation()
  observeThemeChanges()
})

onUnmounted(() => {
  if (watchId) {
    navigator.geolocation.clearWatch(watchId)
  }
  if (observer) {
    observer.disconnect()
  }
})
</script>

<style scoped>
.map-container {
  border-radius: 16px;
  overflow: hidden;
  border: 1px solid rgba(0, 0, 0, 0.08);
  position: relative;
  min-height: 400px;
  background: #f8f9fa;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06), 0 0 0 1px rgba(0, 0, 0, 0.02);
  transition: all 0.3s ease;
}

[data-theme="dark"] .map-container {
  background: #1a1a1a;
  border-color: rgba(255, 255, 255, 0.1);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.4), 0 0 0 1px rgba(255, 255, 255, 0.05);
}

.worksite-marker {
  position: relative;
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.worksite-dot {
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: #333333;
  border: 3px solid white;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
  z-index: 2;
  transition: all 0.3s ease;
}

[data-theme="dark"] .worksite-dot {
  background: #ffffff;
  border-color: #333333;
}

.worksite-ring {
  position: absolute;
  inset: -10px;
  border-radius: 50%;
  border: 2px solid rgba(51, 51, 51, 0.4);
  animation: ringPulse 2.5s ease-out infinite;
}

[data-theme="dark"] .worksite-ring {
  border-color: rgba(255, 255, 255, 0.3);
}

@keyframes ringPulse {
  0% { transform: scale(1); opacity: 1; }
  100% { transform: scale(1.6); opacity: 0; }
}

.employee-marker {
  position: relative;
  width: 44px;
  height: 44px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
}

.employee-dot {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  font-size: 15px;
  color: white;
  border: 3px solid white;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
  z-index: 2;
  transition: all 0.3s ease;
}

.employee-dot.inside {
  background: #22C55E;
  box-shadow: 0 4px 16px rgba(34, 197, 94, 0.5);
}

.employee-dot.outside {
  background: #EF4444;
  box-shadow: 0 4px 16px rgba(239, 68, 68, 0.5);
}

.employee-pulse {
  position: absolute;
  inset: -6px;
  border-radius: 50%;
  z-index: 1;
  animation: employeePulse 2s ease-out infinite;
}

.employee-pulse.inside {
  border: 3px solid rgba(34, 197, 94, 0.6);
}

.employee-pulse.outside {
  border: 3px solid rgba(239, 68, 68, 0.6);
}

@keyframes employeePulse {
  0% { transform: scale(1); opacity: 1; }
  100% { transform: scale(1.7); opacity: 0; }
}

.employee-marker:hover .employee-dot {
  transform: scale(1.15);
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.4);
}

.employee-marker:hover .employee-pulse {
  animation-duration: 1s;
}

.map-overlay {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.92);
  backdrop-filter: blur(12px);
  z-index: 1000;
  padding: 40px;
  border-radius: 16px;
}

[data-theme="dark"] .map-overlay {
  background: rgba(26, 26, 26, 0.92);
}

.overlay-content {
  text-align: center;
  max-width: 400px;
  animation: fadeInUp 0.6s cubic-bezier(0.16, 1, 0.3, 1);
}

@keyframes fadeInUp {
  from { opacity: 0; transform: translateY(24px); }
  to { opacity: 1; transform: translateY(0); }
}

.overlay-icon {
  font-size: 72px;
  display: block;
  margin-bottom: 20px;
  filter: drop-shadow(0 4px 12px rgba(0, 0, 0, 0.1));
}

.overlay-content h3 {
  font-size: 22px;
  font-weight: 700;
  color: var(--ink);
  margin-bottom: 12px;
  letter-spacing: -0.4px;
}

[data-theme="dark"] .overlay-content h3 {
  color: #f1f5f9;
}

.overlay-content p {
  font-size: 15px;
  color: var(--ink-soft);
  margin-bottom: 20px;
  line-height: 1.6;
}

[data-theme="dark"] .overlay-content p {
  color: #cbd5e1;
}

.popup-content h4 {
  margin: 0 0 8px;
  font-size: 16px;
  font-weight: 700;
  color: var(--ink);
  letter-spacing: -0.3px;
}

.popup-content p {
  margin: 4px 0;
  font-size: 13px;
  color: var(--ink-soft);
  line-height: 1.5;
}

[data-theme="dark"] .popup-content h4 {
  color: #f1f5f9;
}

[data-theme="dark"] .popup-content p {
  color: #cbd5e1;
}

.popup-content .btn {
  margin-top: 12px;
  font-size: 13px;
  padding: 6px 16px;
  border-radius: 8px;
  font-weight: 600;
}

.employee-popup {
  min-width: 220px;
  padding: 4px;
}

.badge--in {
  background: rgba(34, 197, 94, 0.15);
  color: #22C55E;
  padding: 4px 12px;
  border-radius: 8px;
  font-size: 12px;
  font-weight: 600;
  display: inline-block;
}

.badge--out {
  background: rgba(239, 68, 68, 0.15);
  color: #EF4444;
  padding: 4px 12px;
  border-radius: 8px;
  font-size: 12px;
  font-weight: 600;
  display: inline-block;
}

/* Wolt-style map controls and attribution */
.wolt-map :deep(.leaflet-control-attribution) {
  background: rgba(255, 255, 255, 0.9) !important;
  backdrop-filter: blur(8px);
  border-radius: 8px;
  padding: 4px 8px;
  font-size: 11px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

[data-theme="dark"] .wolt-map :deep(.leaflet-control-attribution) {
  background: rgba(0, 0, 0, 0.8) !important;
  color: rgba(255, 255, 255, 0.7);
}

.wolt-map :deep(.leaflet-control-attribution a) {
  color: #333 !important;
}

[data-theme="dark"] .wolt-map :deep(.leaflet-control-attribution a) {
  color: rgba(255, 255, 255, 0.8) !important;
}

.wolt-map :deep(.leaflet-control-zoom) {
  border: none !important;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15) !important;
  border-radius: 8px;
  overflow: hidden;
}

.wolt-map :deep(.leaflet-control-zoom a) {
  background: white !important;
  color: #333 !important;
  border: none !important;
  width: 36px !important;
  height: 36px !important;
  line-height: 36px !important;
  font-size: 18px;
  font-weight: 600;
  transition: all 0.2s ease;
}

.wolt-map :deep(.leaflet-control-zoom a:hover) {
  background: #f0f0f0 !important;
  transform: scale(1.05);
}

[data-theme="dark"] .wolt-map :deep(.leaflet-control-zoom a) {
  background: #2a2a2a !important;
  color: white !important;
}

[data-theme="dark"] .wolt-map :deep(.leaflet-control-zoom a:hover) {
  background: #3a3a3a !important;
}

.wolt-map :deep(.leaflet-popup-content-wrapper) {
  border-radius: 12px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
  padding: 0;
  overflow: hidden;
}

.wolt-map :deep(.leaflet-popup-content) {
  margin: 0;
  padding: 16px;
  line-height: 1.5;
}

.wolt-map :deep(.leaflet-popup-tip) {
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
}
</style>
