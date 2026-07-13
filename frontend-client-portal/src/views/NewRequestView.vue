<template>
  <div class="request-page">
    <!-- Hero Section -->
    <div class="hero-section">
      <div class="hero-icon">🔧</div>
      <h1>{{ t('app_name') }}</h1>
      <p>خدمات ميدانية احترافية - نصل إليك أينما كنت</p>
    </div>

    <div class="card request-form">
      <div class="form-header">
        <h2>📝 طلب خدمة جديد</h2>
        <p>املأ البيانات وسنقوم بإرسال أفضل فريق إليك</p>
      </div>

      <form @submit.prevent="submitRequest">
        <!-- رسائل النجاح والخطأ -->
        <div v-if="error" class="alert alert-error">
          <span>❌</span> {{ error }}
        </div>
        <div v-if="success" class="alert alert-success">
          <span>✅</span> {{ success }}
        </div>

        <!-- حقل الاسم -->
        <div class="form-group">
          <label>{{ t('full_name') }} <span class="required">*</span></label>
          <div class="input-wrapper">
            <span class="input-icon">👤</span>
            <input 
              v-model="form.full_name" 
              type="text" 
              required 
              placeholder="أدخل اسمك الكامل"
            />
          </div>
        </div>

        <!-- حقل الهاتف -->
        <div class="form-group">
          <label>{{ t('phone') }} <span class="required">*</span></label>
          <div class="input-wrapper">
            <span class="input-icon">📞</span>
            <input 
              v-model="form.phone" 
              type="tel" 
              required 
              placeholder="05xxxxxxxx"
            />
          </div>
        </div>

        <!-- حقل عنوان الخدمة -->
        <div class="form-group">
          <label>عنوان الخدمة <span class="required">*</span></label>
          <div class="input-wrapper">
            <span class="input-icon">📋</span>
            <input 
              v-model="form.title" 
              type="text" 
              required 
              placeholder="مثال: صيانة مكيف - تركيب كاميرات"
            />
          </div>
        </div>

        <!-- حقل الوصف -->
        <div class="form-group">
          <label>وصف الخدمة <span class="required">*</span></label>
          <div class="input-wrapper textarea-wrapper">
            <span class="input-icon">📝</span>
            <textarea 
              v-model="form.description" 
              required 
              placeholder="وصف تفصيلي للخدمة المطلوبة..."
              rows="4"
            ></textarea>
          </div>
        </div>

        <!-- حقل العنوان -->
        <div class="form-group">
          <label>العنوان التفصيلي</label>
          <div class="input-wrapper">
            <span class="input-icon">📍</span>
            <input 
              v-model="form.address" 
              type="text" 
              placeholder="العنوان الكامل مع المعالم القريبة"
            />
          </div>
        </div>

        <!-- الأولوية -->
        <div class="form-group">
          <label>الأولوية</label>
          <div class="priority-options">
            <label 
              v-for="p in priorities" 
              :key="p.value"
              class="priority-option"
              :class="{ active: form.priority === p.value }"
            >
              <input 
                type="radio" 
                :value="p.value" 
                v-model="form.priority"
                :style="{ accentColor: p.color }"
              />
              <span class="priority-label">
                <span class="priority-dot" :style="{ background: p.color }"></span>
                {{ p.label }}
              </span>
            </label>
          </div>
        </div>

        <!-- تحديد الموقع -->
        <div class="form-group location-group">
          <label>📍 موقعك <span class="required">*</span></label>
          
          <div class="location-status" :class="locationStatusClass">
            <div v-if="locationStatus === 'idle'" class="location-message">
              <span>📍</span>
              <span>اضغط لتحديد موقعك الحالي</span>
            </div>
            <div v-else-if="locationStatus === 'loading'" class="location-message">
              <span class="spinner"></span>
              <span>جارٍ تحديد موقعك...</span>
            </div>
            <div v-else-if="locationStatus === 'success'" class="location-message">
              <span>✅</span>
              <span>تم تحديد موقعك بنجاح!</span>
              <span class="location-coords mono">
                {{ location.lat.toFixed(6) }}, {{ location.lng.toFixed(6) }}
              </span>
            </div>
            <div v-else-if="locationStatus === 'error'" class="location-message">
              <span>❌</span>
              <span>{{ locationError }}</span>
            </div>
          </div>

          <div class="location-actions">
            <button 
              type="button" 
              class="btn btn--primary"
              @click="getLocation"
              :disabled="locationStatus === 'loading'"
            >
              {{ locationStatus === 'loading' ? '⏳ جارٍ...' : '📍 تحديد موقعي' }}
            </button>
            <button 
              type="button" 
              class="btn btn--ghost"
              @click="openMapPicker"
            >
              🗺️ اختيار من الخريطة
            </button>
          </div>
        </div>

        <!-- زر الإرسال -->
        <button 
          type="submit" 
          class="btn btn--primary btn--block btn--large"
          :disabled="loading || !location"
        >
          <span v-if="loading" class="spinner"></span>
          {{ loading ? 'جاري الإرسال...' : '🚀 إرسال الطلب' }}
        </button>

        <p class="form-footer">
          * جميع البيانات محفوظة بشكل آمن
        </p>
      </form>
    </div>

    <!-- خريطة اختيار الموقع (Modal) -->
    <div v-if="showMapPicker" class="modal-backdrop" @click.self="showMapPicker = false">
      <div class="modal map-modal">
        <div class="modal-header">
          <h3>🗺️ اختر موقعك على الخريطة</h3>
          <button class="modal-close" @click="showMapPicker = false">✕</button>
        </div>
        <div class="modal-body">
          <div ref="mapContainer" class="map-container"></div>
        </div>
        <div class="modal-footer">
          <button class="btn btn--ghost" @click="showMapPicker = false">إلغاء</button>
          <button class="btn btn--primary" @click="confirmMapLocation">تأكيد الموقع</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from '../plugins/i18n'
import axios from 'axios'

const { t } = useI18n()
const router = useRouter()

const loading = ref(false)
const error = ref('')
const success = ref('')
const location = ref(null)
const locationStatus = ref('idle')
const locationError = ref('')
const showMapPicker = ref(false)
const mapContainer = ref(null)
let map = null
let mapMarker = null

const form = reactive({
  full_name: '',
  phone: '',
  title: '',
  description: '',
  address: '',
  priority: 'normal'
})

const priorities = [
  { value: 'low', label: 'منخفضة', color: '#4CAF50' },
  { value: 'normal', label: 'عادية', color: '#2196F3' },
  { value: 'high', label: 'عالية', color: '#FF9800' },
  { value: 'urgent', label: 'طارئة', color: '#f44336' }
]

const locationStatusClass = computed(() => {
  return {
    'idle': locationStatus.value === 'idle',
    'loading': locationStatus.value === 'loading',
    'success': locationStatus.value === 'success',
    'error': locationStatus.value === 'error'
  }
})

function getLocation() {
  if (!('geolocation' in navigator)) {
    locationStatus.value = 'error'
    locationError.value = 'المتصفح لا يدعم تحديد الموقع'
    return
  }

  locationStatus.value = 'loading'
  locationError.value = ''

  navigator.geolocation.getCurrentPosition(
    (pos) => {
      location.value = {
        lat: pos.coords.latitude,
        lng: pos.coords.longitude
      }
      locationStatus.value = 'success'
    },
    (err) => {
      locationStatus.value = 'error'
      locationError.value = 'فشل تحديد الموقع: ' + err.message
    },
    { enableHighAccuracy: true, timeout: 15000 }
  )
}

function openMapPicker() {
  showMapPicker.value = true
  nextTick(() => {
    initMap()
  })
}

function initMap() {
  // سيتم ربط Leaflet هنا
}

function confirmMapLocation() {
  showMapPicker.value = false
  locationStatus.value = 'success'
}

async function submitRequest() {
  if (!location.value) {
    error.value = 'الرجاء تحديد موقعك أولاً'
    return
  }

  loading.value = true
  error.value = ''
  success.value = ''

  try {
    // هنا يتم إرسال الطلب للـ API
    const payload = {
      full_name: form.full_name,
      phone: form.phone,
      title: form.title,
      description: form.description,
      address: form.address,
      latitude: location.value.lat,
      longitude: location.value.lng,
      priority: form.priority
    }

    // محاكاة إرسال (سيتم ربطها بالـ API الفعلي)
    await new Promise(resolve => setTimeout(resolve, 1500))
    
    success.value = '✅ تم إرسال طلب الخدمة بنجاح! سنتواصل معك قريباً.'
    
    // إعادة تعيين النموذج
    setTimeout(() => {
      router.push('/')
    }, 3000)
  } catch (err) {
    error.value = err.response?.data?.error || 'فشل إرسال الطلب'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.request-page {
  max-width: 600px;
  margin: 0 auto;
  padding: 20px 0;
}

.hero-section {
  text-align: center;
  padding: 40px 20px 30px;
  background: linear-gradient(135deg, var(--brand), var(--brand-dark));
  border-radius: var(--radius-lg);
  color: white;
  margin-bottom: 30px;
}

.hero-icon {
  font-size: 48px;
  margin-bottom: 12px;
}

.hero-section h1 {
  font-size: 28px;
  margin-bottom: 6px;
}

.hero-section p {
  opacity: 0.9;
  font-size: 14px;
}

.request-form {
  padding: 28px 24px;
}

.form-header {
  text-align: center;
  margin-bottom: 24px;
}

.form-header h2 {
  font-size: 20px;
  margin-bottom: 4px;
}

.form-header p {
  color: var(--ink-soft);
  font-size: 13px;
}

.form-group {
  margin-bottom: 18px;
}

.form-group label {
  display: block;
  font-size: 13px;
  font-weight: 600;
  color: var(--ink);
  margin-bottom: 6px;
}

.input-wrapper {
  position: relative;
  display: flex;
  align-items: center;
}

.input-icon {
  position: absolute;
  right: 12px;
  font-size: 16px;
  opacity: 0.6;
  pointer-events: none;
}

.input-wrapper input,
.input-wrapper textarea {
  width: 100%;
  padding: 12px 42px 12px 14px;
  border: 1.5px solid var(--line);
  border-radius: var(--radius-sm);
  font-family: var(--font-body);
  font-size: 14px;
  background: var(--surface);
  transition: all 0.3s ease;
}

.input-wrapper input:focus,
.input-wrapper textarea:focus {
  border-color: var(--brand);
  box-shadow: 0 0 0 3px var(--brand-tint);
  outline: none;
}

.textarea-wrapper {
  align-items: flex-start;
}

.textarea-wrapper .input-icon {
  top: 12px;
}

.required {
  color: var(--signal-out);
}

.priority-options {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 8px;
}

.priority-option {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 10px;
  border: 1.5px solid var(--line);
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: all 0.3s ease;
  font-size: 13px;
}

.priority-option:hover {
  border-color: var(--brand);
  background: var(--brand-tint);
}

.priority-option.active {
  border-color: var(--brand);
  background: var(--brand-tint);
}

.priority-option input[type="radio"] {
  display: none;
}

.priority-label {
  display: flex;
  align-items: center;
  gap: 4px;
}

.priority-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}

.location-group {
  background: var(--canvas);
  padding: 16px;
  border-radius: var(--radius-sm);
}

.location-status {
  padding: 12px;
  border-radius: var(--radius-sm);
  margin-bottom: 12px;
  text-align: center;
}

.location-status.idle {
  background: var(--line);
  color: var(--ink-soft);
}

.location-status.loading {
  background: var(--brand-tint);
  color: var(--brand);
}

.location-status.success {
  background: var(--signal-in-tint);
  color: var(--signal-in);
}

.location-status.error {
  background: var(--signal-out-tint);
  color: var(--signal-out);
}

.location-message {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  flex-wrap: wrap;
}

.location-coords {
  font-size: 12px;
  opacity: 0.8;
}

.location-actions {
  display: flex;
  gap: 10px;
}

.location-actions .btn {
  flex: 1;
}

.alert {
  padding: 12px 16px;
  border-radius: var(--radius-sm);
  margin-bottom: 16px;
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 13px;
}

.alert-error {
  background: var(--signal-out-tint);
  color: var(--signal-out);
}

.alert-success {
  background: var(--signal-in-tint);
  color: var(--signal-in);
}

.btn--large {
  padding: 16px;
  font-size: 16px;
}

.btn--block {
  width: 100%;
}

.form-footer {
  text-align: center;
  font-size: 12px;
  color: var(--ink-soft);
  margin-top: 14px;
}

.spinner {
  display: inline-block;
  width: 18px;
  height: 18px;
  border: 2px solid rgba(255,255,255,0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.map-container {
  height: 400px;
  border-radius: var(--radius-sm);
}

.modal-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.6);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 20px;
}

.map-modal {
  width: 100%;
  max-width: 700px;
  padding: 0;
  overflow: hidden;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid var(--line);
}

.modal-header h3 { font-size: 16px; margin: 0; }

.modal-close {
  background: none;
  border: none;
  font-size: 20px;
  cursor: pointer;
  color: var(--ink-soft);
}

.modal-body {
  padding: 0;
}

.modal-footer {
  padding: 16px 20px;
  border-top: 1px solid var(--line);
  display: flex;
  gap: 10px;
  justify-content: flex-end;
}

@media (max-width: 480px) {
  .priority-options {
    grid-template-columns: repeat(2, 1fr);
  }
  .location-actions {
    flex-direction: column;
  }
  .hero-section {
    padding: 30px 16px;
  }
}
</style>
