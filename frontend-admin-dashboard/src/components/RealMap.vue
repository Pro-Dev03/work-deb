<template>
  <div class="map-container">
    <l-map
      ref="mapRef"
      :zoom="zoom"
      @update:zoom="updateZoom"
      :center="center"
      :options="{ attributionControl: true, zoomControl: true }"
      :style="{ height: height + 'px', width: '100%' }"
    >
      <l-tile-layer
        :url="mapTileUrl"
        layer-type="base"
        :name="isDarkMode ? 'CartoDB Dark' : 'OpenStreetMap'"
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
            <h4>🏢 {{ site.name }}</h4>
            <p>{{ site.address || 'لا يوجد عنوان' }}</p>
            <p>⭕ النطاق: {{ site.radius_meters }} متر</p>
            <p>👥 عدد الموظفين: {{ getEmployeeCount(site.id) }}</p>
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
            <h4>👤 {{ emp.full_name }}</h4>
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
              📋 عرض التفاصيل
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

// ✅ اكتشاف الوضع الداكن
const isDarkMode = computed(() => {
  return document.documentElement.getAttribute('data-theme') === 'dark'
})

// ✅ تحديد URL الخريطة حسب الوضع
const mapTileUrl = computed(() => {
  if (isDarkMode.value) {
    // خريطة داكنة من CartoDB
    return 'https://{s}.basemaps.cartocdn.com/dark_all/{z}/{x}/{y}{r}.png'
  } else {
    // خريطة عادية من OpenStreetMap
    return 'https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png'
  }
})

// ✅ تحديد Attribution حسب الوضع
const mapAttribution = computed(() => {
  if (isDarkMode.value) {
    return '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors &copy; <a href="https://carto.com/attributions">CARTO</a>'
  } else {
    return '&copy; OpenStreetMap contributors'
  }
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
})

onUnmounted(() => {
  if (watchId) {
    navigator.geolocation.clearWatch(watchId)
  }
})
</script>

<style scoped>
.map-container {
  border-radius: var(--radius-lg);
  overflow: hidden;
  border: 1px solid var(--line);
  position: relative;
  min-height: 400px;
  background: #E8EDF2;
}

[data-theme="dark"] .map-container {
  background: #0f172a;
  border-color: #334155;
}

.worksite-marker {
  position: relative;
  width: 30px;
  height: 30px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.worksite-dot {
  width: 16px;
  height: 16px;
  border-radius: 50%;
  background: #2563EB;
  border: 2px solid white;
  box-shadow: 0 2px 8px rgba(37, 99, 235, 0.4);
  z-index: 2;
}

.worksite-ring {
  position: absolute;
  inset: -8px;
  border-radius: 50%;
  border: 2px solid rgba(37, 99, 235, 0.3);
  animation: ringPulse 2s ease-out infinite;
}

@keyframes ringPulse {
  0% { transform: scale(1); opacity: 1; }
  100% { transform: scale(1.5); opacity: 0; }
}

.employee-marker {
  position: relative;
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
}

.employee-dot {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  font-size: 14px;
  color: white;
  border: 2px solid white;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.25);
  z-index: 2;
  transition: all 0.3s ease;
}

.employee-dot.inside {
  background: #22C55E;
  box-shadow: 0 0 20px rgba(34, 197, 94, 0.4);
}

.employee-dot.outside {
  background: #EF4444;
  box-shadow: 0 0 20px rgba(239, 68, 68, 0.4);
}

.employee-pulse {
  position: absolute;
  inset: -4px;
  border-radius: 50%;
  z-index: 1;
  animation: employeePulse 1.5s ease-out infinite;
}

.employee-pulse.inside {
  border: 2px solid #22C55E;
}

.employee-pulse.outside {
  border: 2px solid #EF4444;
}

@keyframes employeePulse {
  0% { transform: scale(1); opacity: 1; }
  100% { transform: scale(1.6); opacity: 0; }
}

.employee-marker:hover .employee-dot {
  transform: scale(1.1);
}

.employee-marker:hover .employee-pulse {
  animation-duration: 0.8s;
}

.map-overlay {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.85);
  backdrop-filter: blur(4px);
  z-index: 1000;
  padding: 30px;
  border-radius: var(--radius-lg);
}

[data-theme="dark"] .map-overlay {
  background: rgba(15, 23, 42, 0.85);
}

.overlay-content {
  text-align: center;
  max-width: 400px;
  animation: fadeInUp 0.5s ease;
}

@keyframes fadeInUp {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: translateY(0); }
}

.overlay-icon {
  font-size: 64px;
  display: block;
  margin-bottom: 16px;
}

.overlay-content h3 {
  font-size: 20px;
  color: var(--ink);
  margin-bottom: 8px;
}

.overlay-content p {
  font-size: 14px;
  color: var(--ink-soft);
  margin-bottom: 16px;
}

.popup-content h4 {
  margin: 0 0 4px;
  font-size: 14px;
  color: var(--ink);
}

.popup-content p {
  margin: 2px 0;
  font-size: 12px;
  color: var(--ink-soft);
}

[data-theme="dark"] .popup-content h4 {
  color: #f1f5f9;
}

[data-theme="dark"] .popup-content p {
  color: #cbd5e1;
}

.popup-content .btn {
  margin-top: 8px;
  font-size: 12px;
  padding: 4px 12px;
}

.employee-popup {
  min-width: 200px;
}

.badge--in {
  background: #22C55E20;
  color: #22C55E;
  padding: 2px 10px;
  border-radius: 999px;
  font-size: 11px;
  font-weight: 600;
}

.badge--out {
  background: #EF444420;
  color: #EF4444;
  padding: 2px 10px;
  border-radius: 999px;
  font-size: 11px;
  font-weight: 600;
}
</style>
