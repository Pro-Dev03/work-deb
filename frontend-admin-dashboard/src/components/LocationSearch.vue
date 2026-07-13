<template>
  <div class="location-search">
    <!-- ========================================== -->
    <!-- 1. اختيار المدينة -->
    <!-- ========================================== -->
    <div class="search-group">
      <label class="search-label">🏙️ المدينة <span class="required">*</span></label>
      <div class="search-wrapper">
        <input
          v-model="cityQuery"
          type="text"
          placeholder="ابحث عن مدينة في إسرائيل..."
          @input="onCitySearch"
          @focus="showCityResults = true"
          class="search-input"
          required
        />
        <span v-if="cityLoading" class="search-loading">⏳</span>
      </div>

      <!-- نتائج المدن -->
      <div v-if="showCityResults && cityResults.length > 0" class="search-results">
        <div
          v-for="city in cityResults"
          :key="city.place_id"
          class="result-item"
          @click="selectCity(city)"
        >
          <span class="result-icon">🏙️</span>
          <div class="result-info">
            <strong>{{ city.display_name.split(',')[0] }}</strong>
            <span class="result-address">{{ city.display_name }}</span>
          </div>
        </div>
      </div>

      <div v-if="showCityResults && cityQuery && cityResults.length === 0 && !cityLoading" class="no-results">
        <span>❌</span>
        <p>لم يتم العثور على مدينة</p>
      </div>
    </div>

    <!-- ========================================== -->
    <!-- 2. اختيار الشارع (يظهر بعد اختيار المدينة) -->
    <!-- ========================================== -->
    <div class="search-group" v-if="selectedCity">
      <label class="search-label">📍 الشارع <span class="required">*</span></label>
      <div class="search-wrapper">
        <input
          v-model="streetQuery"
          type="text"
          :placeholder="`ابحث عن شارع في ${selectedCityName}...`"
          @input="onStreetSearch"
          @focus="showStreetResults = true"
          class="search-input"
          required
        />
        <span v-if="streetLoading" class="search-loading">⏳</span>
      </div>

      <!-- نتائج الشوارع -->
      <div v-if="showStreetResults && streetResults.length > 0" class="search-results">
        <div
          v-for="street in streetResults"
          :key="street.place_id"
          class="result-item"
          @click="selectStreet(street)"
        >
          <span class="result-icon">📍</span>
          <div class="result-info">
            <strong>{{ street.display_name.split(',')[0] }}</strong>
            <span class="result-address">{{ street.display_name }}</span>
          </div>
        </div>
      </div>

      <div v-if="showStreetResults && streetQuery && streetResults.length === 0 && !streetLoading" class="no-results">
        <span>❌</span>
        <p>لم يتم العثور على شارع</p>
      </div>
    </div>

    <!-- ========================================== -->
    <!-- 3. رقم المبنى -->
    <!-- ========================================== -->
    <div class="search-group" v-if="selectedStreet">
      <label class="search-label">🏢 رقم المبنى</label>
      <input
        v-model="buildingNumber"
        type="text"
        placeholder="مثال: 15"
        class="search-input"
      />
      <span class="field-hint">أدخل رقم المبنى (اختياري)</span>
    </div>

    <!-- ========================================== -->
    <!-- 4. الموقع المختار -->
    <!-- ========================================== -->
    <div v-if="selectedLocation" class="selected-location">
      <div class="location-preview">
        <span class="location-icon">📍</span>
        <div>
          <strong>{{ selectedLocation.name }}</strong>
          <p>{{ selectedLocation.address }}</p>
          <p class="mono">خط العرض: {{ selectedLocation.latitude.toFixed(6) }}</p>
          <p class="mono">خط الطول: {{ selectedLocation.longitude.toFixed(6) }}</p>
        </div>
      </div>
    </div>

    <!-- ========================================== -->
    <!-- زر تحديد الموقع الحالي -->
    <!-- ========================================== -->
    <button class="btn btn--ghost btn--block" @click="getCurrentLocation" type="button">
      📍 استخدام موقعي الحالي
    </button>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'

const emit = defineEmits(['select'])

// ==========================================
// حالة البحث عن المدن
// ==========================================
const cityQuery = ref('')
const cityResults = ref([])
const cityLoading = ref(false)
const showCityResults = ref(false)
const selectedCity = ref(null)
const selectedCityName = computed(() => selectedCity.value?.display_name?.split(',')[0] || '')

// ==========================================
// حالة البحث عن الشوارع
// ==========================================
const streetQuery = ref('')
const streetResults = ref([])
const streetLoading = ref(false)
const showStreetResults = ref(false)
const selectedStreet = ref(null)
const buildingNumber = ref('')

// ==========================================
// الموقع المختار
// ==========================================
const selectedLocation = ref(null)

let cityTimeout = null
let streetTimeout = null

// ==========================================
// البحث عن المدن في إسرائيل
// ==========================================
function onCitySearch() {
  clearTimeout(cityTimeout)
  
  if (cityQuery.value.length < 2) {
    cityResults.value = []
    return
  }

  cityLoading.value = true
  cityTimeout = setTimeout(async () => {
    try {
      const response = await fetch(
        `https://nominatim.openstreetmap.org/search?q=${encodeURIComponent(cityQuery.value)}&countrycodes=il&featuretype=city&format=json&limit=10&addressdetails=1`
      )
      
      if (!response.ok) throw new Error('فشل البحث')
      
      const data = await response.json()
      cityResults.value = data.map(item => ({
        place_id: item.place_id,
        display_name: item.display_name,
        lat: parseFloat(item.lat),
        lon: parseFloat(item.lon),
        address: item.address || {}
      }))
      
      showCityResults.value = true
    } catch (error) {
      console.error('خطأ في البحث عن المدن:', error)
      cityResults.value = []
    } finally {
      cityLoading.value = false
    }
  }, 500)
}

// ==========================================
// اختيار مدينة
// ==========================================
function selectCity(city) {
  selectedCity.value = city
  cityQuery.value = city.display_name.split(',')[0]
  showCityResults.value = false
  
  // إعادة تعيين الشارع
  streetQuery.value = ''
  streetResults.value = []
  selectedStreet.value = null
  
  // تحديث الموقع
  updateLocation()
}

// ==========================================
// البحث عن الشوارع في المدينة المختارة
// ==========================================
function onStreetSearch() {
  clearTimeout(streetTimeout)
  
  if (!selectedCity.value || streetQuery.value.length < 2) {
    streetResults.value = []
    return
  }

  streetLoading.value = true
  streetTimeout = setTimeout(async () => {
    try {
      const response = await fetch(
        `https://nominatim.openstreetmap.org/search?q=${encodeURIComponent(streetQuery.value)}&city=${encodeURIComponent(selectedCity.value.display_name.split(',')[0])}&countrycodes=il&featuretype=street&format=json&limit=10&addressdetails=1`
      )
      
      if (!response.ok) throw new Error('فشل البحث')
      
      const data = await response.json()
      streetResults.value = data.map(item => ({
        place_id: item.place_id,
        display_name: item.display_name,
        lat: parseFloat(item.lat),
        lon: parseFloat(item.lon),
        address: item.address || {}
      }))
      
      showStreetResults.value = true
    } catch (error) {
      console.error('خطأ في البحث عن الشوارع:', error)
      streetResults.value = []
    } finally {
      streetLoading.value = false
    }
  }, 500)
}

// ==========================================
// اختيار شارع
// ==========================================
function selectStreet(street) {
  selectedStreet.value = street
  streetQuery.value = street.display_name.split(',')[0]
  showStreetResults.value = false
  
  // تحديث الموقع
  updateLocation()
}

// ==========================================
// تحديث الموقع المختار
// ==========================================
function updateLocation() {
  if (!selectedCity.value) return
  
  const lat = selectedStreet.value?.lat || selectedCity.value.lat
  const lon = selectedStreet.value?.lon || selectedCity.value.lon
  
  const location = {
    name: selectedCity.value.display_name.split(',')[0],
    address: buildAddress(),
    latitude: lat,
    longitude: lon,
    city: selectedCity.value.display_name.split(',')[0],
    street: selectedStreet.value?.display_name?.split(',')[0] || '',
    building_number: buildingNumber.value || '',
    country: 'إسرائيل'
  }
  
  selectedLocation.value = location
  emit('select', location)
}

// ==========================================
// بناء العنوان الكامل
// ==========================================
function buildAddress() {
  let parts = []
  if (buildingNumber.value) parts.push(buildingNumber.value)
  if (selectedStreet.value) parts.push(selectedStreet.value.display_name.split(',')[0])
  if (selectedCity.value) parts.push(selectedCity.value.display_name.split(',')[0])
  return parts.join('، ') || 'إسرائيل'
}

// ==========================================
// مراقبة رقم المبنى
// ==========================================
function updateBuildingNumber() {
  if (selectedLocation.value) {
    selectedLocation.value.building_number = buildingNumber.value
    selectedLocation.value.address = buildAddress()
    emit('select', selectedLocation.value)
  }
}

// ==========================================
// استخدام الموقع الحالي
// ==========================================
function getCurrentLocation() {
  if (!('geolocation' in navigator)) {
    alert('المتصفح لا يدعم تحديد الموقع')
    return
  }

  navigator.geolocation.getCurrentPosition(
    async (pos) => {
      const { latitude, longitude } = pos.coords
      
      try {
        const response = await fetch(
          `https://nominatim.openstreetmap.org/reverse?lat=${latitude}&lon=${longitude}&format=json&zoom=10`
        )
        const data = await response.json()
        
        const location = {
          name: data.address?.city || data.address?.town || data.address?.village || 'موقعي',
          address: data.display_name || `${latitude}, ${longitude}`,
          latitude: latitude,
          longitude: longitude,
          city: data.address?.city || data.address?.town || '',
          street: data.address?.road || '',
          building_number: data.address?.house_number || '',
          country: data.address?.country || 'إسرائيل'
        }
        
        selectedLocation.value = location
        emit('select', location)
        
        // تعبئة الحقول
        cityQuery.value = location.city
        streetQuery.value = location.street
        buildingNumber.value = location.building_number
        
      } catch (error) {
        emit('select', {
          name: 'موقعي الحالي',
          address: `${latitude}, ${longitude}`,
          latitude: latitude,
          longitude: longitude,
          country: 'إسرائيل'
        })
      }
    },
    (err) => {
      alert('فشل تحديد الموقع: ' + err.message)
    },
    { enableHighAccuracy: true }
  )
}

// ==========================================
// إغلاق النتائج عند النقر خارجها
// ==========================================
function closeResults(e) {
  const wrapper = document.querySelector('.location-search')
  if (wrapper && !wrapper.contains(e.target)) {
    showCityResults.value = false
    showStreetResults.value = false
  }
}

document.addEventListener('click', closeResults)

// مراقبة رقم المبنى
import { watch } from 'vue'
watch(buildingNumber, updateBuildingNumber)
</script>

<style scoped>
.location-search {
  width: 100%;
}

/* ==========================================
   مجموعات البحث
   ========================================== */
.search-group {
  margin-bottom: 14px;
  position: relative;
}

.search-label {
  display: block;
  font-size: 13px;
  font-weight: 600;
  color: var(--ink);
  margin-bottom: 4px;
}

.search-wrapper {
  position: relative;
  display: flex;
  align-items: center;
}

.search-input {
  width: 100%;
  padding: 10px 14px;
  padding-left: 40px;
  border: 1.5px solid var(--line);
  border-radius: var(--radius-sm);
  font-size: 14px;
  font-family: var(--font-body);
  background: var(--surface);
  transition: all 0.3s;
}

.search-input:focus {
  outline: none;
  border-color: var(--brand);
  box-shadow: 0 0 0 3px var(--brand-tint);
}

.search-loading {
  position: absolute;
  left: 12px;
  font-size: 16px;
}

.required {
  color: var(--signal-out);
}

.field-hint {
  font-size: 11px;
  color: var(--ink-soft);
  display: block;
  margin-top: 4px;
}

/* ==========================================
   نتائج البحث
   ========================================== */
.search-results {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  background: var(--surface);
  border: 1px solid var(--line);
  border-radius: var(--radius-sm);
  box-shadow: var(--shadow-lg);
  max-height: 250px;
  overflow-y: auto;
  z-index: 1000;
  margin-top: 4px;
}

.result-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 14px;
  cursor: pointer;
  border-bottom: 1px solid var(--line);
  transition: background 0.2s;
}

.result-item:hover {
  background: var(--brand-tint);
}

.result-item:last-child {
  border-bottom: none;
}

.result-icon {
  font-size: 18px;
  flex-shrink: 0;
}

.result-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.result-info strong {
  font-size: 13px;
  color: var(--ink);
}

.result-address {
  font-size: 11px;
  color: var(--ink-soft);
}

/* ==========================================
   عدم وجود نتائج
   ========================================== */
.no-results {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  background: var(--surface);
  border: 1px solid var(--line);
  border-radius: var(--radius-sm);
  padding: 16px;
  text-align: center;
  z-index: 1000;
  margin-top: 4px;
}

.no-results span {
  font-size: 24px;
  display: block;
  margin-bottom: 4px;
}

.no-results p {
  font-size: 13px;
  color: var(--ink-soft);
  margin: 0;
}

/* ==========================================
   الموقع المختار
   ========================================== */
.selected-location {
  background: var(--brand-tint);
  border: 1px solid var(--brand);
  border-radius: var(--radius-sm);
  padding: 12px 16px;
  margin-bottom: 14px;
}

.location-preview {
  display: flex;
  align-items: flex-start;
  gap: 12px;
}

.location-icon {
  font-size: 24px;
}

.location-preview strong {
  font-size: 14px;
  color: var(--brand);
}

.location-preview p {
  font-size: 12px;
  color: var(--ink-soft);
  margin: 2px 0;
}

/* ==========================================
   أزرار
   ========================================== */
.btn--block {
  width: 100%;
}

.btn--ghost {
  margin-top: 4px;
}
</style>
