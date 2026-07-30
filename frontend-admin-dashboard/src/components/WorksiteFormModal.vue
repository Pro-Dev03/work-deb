<template>
  <div class="modal-backdrop" @click.self="$emit('close')">
    <div class="modal card">
      <div class="modal-header">
        <h3>📍 {{ isEditMode ? t('edit_worksite') : t('new_worksite') }}</h3>
        <button class="modal-close" @click="$emit('close')">✕</button>
      </div>

      <p class="modal__hint">{{ t('search_worldwide') }} - {{ t('search_language_hint') }}</p>

      <form class="modal__form" @submit.prevent="handleSubmit">
        <div v-if="error" class="alert alert-error">{{ error }}</div>
        <div v-if="success" class="alert alert-success">{{ success }}</div>

        <div class="form-group">
          <label>📛 {{ t('worksite_name') }} <span class="required">*</span></label>
          <input
            v-model="form.name"
            type="text"
            :placeholder="t('worksite_name_placeholder')"
            required
            class="search-input"
          />
        </div>

        <!-- ========================================== -->
        <!-- البحث العالمي متعدد اللغات -->
        <!-- ========================================== -->
        <div class="form-group">
          <label>📍 {{ t('search_address') }} <span class="required">*</span></label>
          <div class="search-wrapper">
            <input
              v-model="searchQuery"
              type="text"
              :placeholder="t('search_address_placeholder')"
              @input="onSearch"
              @focus="showResults = true"
              class="search-input"
              autocomplete="off"
            />
            <span v-if="loading" class="search-loading">⏳</span>
          </div>

          <!-- نتائج البحث متعددة اللغات -->
          <div v-if="showResults && searchResults.length > 0" class="search-results">
            <div
              v-for="result in searchResults"
              :key="result.id"
              class="result-item"
              @click="selectResult(result)"
            >
              <span class="result-icon">📍</span>
              <div class="result-info">
                <strong>{{ getLocalizedLabel(result) }}</strong>
                <span class="result-address">
                  {{ getLocalizedAddress(result) }}
                </span>
              </div>
              <span class="result-type">{{ getTypeLabel(result.type) }}</span>
            </div>
          </div>

          <div v-if="showResults && searchQuery && searchResults.length === 0 && !loading" class="no-results">
            <span>🔍</span>
            <p>{{ t('no_results_found') }}</p>
          </div>
        </div>

        <!-- ========================================== -->
        <!-- الموقع المختار -->
        <!-- ========================================== -->
        <div v-if="selectedResult" class="selected-location">
          <div class="selected-location__header">
            <span>✅ {{ getLocalizedLabel(selectedResult) }}</span>
            <button type="button" class="selected-clear" @click="clearSelection">✕</button>
          </div>
          <div class="selected-location__details">
            <p><strong>{{ t('city') }}:</strong> {{ getLocalizedCity(selectedResult) }}</p>
            <p><strong>{{ t('street') }}:</strong> {{ getLocalizedStreet(selectedResult) }}</p>
            <p><strong>{{ t('building_number') }}:</strong> {{ selectedResult.house_number || '—' }}</p>
            <p class="mono"><strong>{{ t('coordinates') }}:</strong> {{ selectedResult.latitude }}, {{ selectedResult.longitude }}</p>
          </div>
        </div>

        <div class="form-group">
          <label>⭕ {{ t('allowed_radius') }} <span class="required">*</span></label>
          <input
            v-model.number="form.radius_meters"
            type="number"
            required
            placeholder="100"
            min="10"
            class="search-input"
          />
        </div>

        <div class="form-actions">
          <button type="button" class="btn btn--ghost" @click="$emit('close')">{{ t('cancel') }}</button>
          <button type="submit" class="btn btn--primary" :disabled="loading || !selectedResult">
            {{ loading ? `⏳ ${t('saving')}` : `💾 ${t('save')}` }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import api from '../services/api'
import { useI18n } from '../services/i18n'

const { t, currentLang } = useI18n()
const props = defineProps({
  worksite: {
    type: Object,
    default: null
  }
})
const emit = defineEmits(['close', 'worksite-added', 'worksite-updated'])

const searchQuery = ref('')
const searchResults = ref([])
const loading = ref(false)
const showResults = ref(false)
const selectedResult = ref(null)

const form = ref({
  name: '',
  radius_meters: 100,
  latitude: '',
  longitude: ''
})

const error = ref('')
const success = ref('')
const isSubmitting = ref(false)

const isEditMode = computed(() => !!props.worksite)

// تهيئة النموذج عند التعديل
function initializeForm() {
  if (props.worksite) {
    form.value = {
      name: props.worksite.name || '',
      radius_meters: props.worksite.radius_meters || 100,
      latitude: props.worksite.latitude || '',
      longitude: props.worksite.longitude || ''
    }
    if (props.worksite.latitude && props.worksite.longitude) {
      selectedResult.value = {
        latitude: props.worksite.latitude,
        longitude: props.worksite.longitude,
        label: props.worksite.name,
        address: props.worksite.address
      }
      searchQuery.value = props.worksite.address || props.worksite.name
    }
  } else {
    // إعادة تعيين النموذج عند الإضافة
    form.value = {
      name: '',
      radius_meters: 100,
      latitude: '',
      longitude: ''
    }
    selectedResult.value = null
    searchQuery.value = ''
  }
}

// مراقبة تغيرات props.worksite
watch(() => props.worksite, () => {
  initializeForm()
}, { immediate: true })

onMounted(() => {
  initializeForm()
})

let searchTimeout = null

function getLocalizedLabel(result) {
  const lang = currentLang.value
  if (lang === 'he' && result.label_he) return result.label_he
  if (lang === 'ar' && result.label_ar) return result.label_ar
  return result.label || result.street || result.city
}

function getLocalizedAddress(result) {
  const lang = currentLang.value
  let city = ''
  let street = ''
  
  if (lang === 'he') {
    city = result.city_he || result.city || ''
    street = result.street_he || result.street || ''
  } else if (lang === 'ar') {
    city = result.city_ar || result.city || ''
    street = result.street_ar || result.street || ''
  } else {
    city = result.city || ''
    street = result.street || ''
  }
  
  return `${city} ${street} ${result.house_number || ''}`.trim()
}

function getLocalizedCity(result) {
  const lang = currentLang.value
  if (lang === 'he') return result.city_he || result.city || '—'
  if (lang === 'ar') return result.city_ar || result.city || '—'
  return result.city || '—'
}

function getLocalizedStreet(result) {
  const lang = currentLang.value
  if (lang === 'he') return result.street_he || result.street || '—'
  if (lang === 'ar') return result.street_ar || result.street || '—'
  return result.street || '—'
}

function getTypeLabel(type) {
  const types = {
    'city': t('type_city'),
    'street': t('type_street'),
    'address': t('type_address'),
    'house': t('type_house'),
    'landmark': t('type_landmark'),
    'location': t('type_location')
  }
  return types[type] || type || t('type_location')
}

function onSearch() {
  clearTimeout(searchTimeout)
  
  if (searchQuery.value.length < 2) {
    searchResults.value = []
    return
  }

  loading.value = true
  searchTimeout = setTimeout(async () => {
    try {
      // استخدام اللغة الحالية
      const lang = currentLang.value
      const { data } = await api.get('/geocode/autocomplete', {
        params: {
          q: searchQuery.value.trim(),
          lang: lang
        }
      })
      
      searchResults.value = data.results || []
      showResults.value = true
    } catch (error) {
      console.error('❌ فشل البحث:', error)
      searchResults.value = []
    } finally {
      loading.value = false
    }
  }, 600)
}

function selectResult(result) {
  selectedResult.value = result
  searchQuery.value = getLocalizedLabel(result)
  showResults.value = false
  
  form.value.name = getLocalizedLabel(result)
  form.value.latitude = result.latitude
  form.value.longitude = result.longitude
  
  error.value = ''
}

function clearSelection() {
  selectedResult.value = null
  searchQuery.value = ''
  searchResults.value = []
  form.value.name = ''
  form.value.latitude = ''
  form.value.longitude = ''
}

function closeResults(e) {
  const wrapper = document.querySelector('.modal__form')
  if (wrapper && !wrapper.contains(e.target)) {
    showResults.value = false
  }
}

document.addEventListener('click', closeResults)

async function handleSubmit() {
  if (!selectedResult.value && !isEditMode.value) {
    error.value = t('select_address_required')
    return
  }

  isSubmitting.value = true
  error.value = ''
  success.value = ''

  try {
    const payload = {
      name: form.value.name || (selectedResult.value ? getLocalizedLabel(selectedResult.value) : ''),
      address: selectedResult.value ? getLocalizedLabel(selectedResult.value) : (props.worksite?.address || ''),
      latitude: parseFloat(form.value.latitude || (selectedResult.value ? selectedResult.value.latitude : props.worksite?.latitude)),
      longitude: parseFloat(form.value.longitude || (selectedResult.value ? selectedResult.value.longitude : props.worksite?.longitude)),
      radius_meters: form.value.radius_meters,
      city: selectedResult.value ? getLocalizedCity(selectedResult.value) : '',
      street: selectedResult.value ? getLocalizedStreet(selectedResult.value) : '',
      street_number: selectedResult.value ? (selectedResult.value.house_number || '') : ''
    }

    if (isEditMode.value) {
      // وضع التعديل
      await api.put(`/worksites/${props.worksite.id}`, payload)
      success.value = '✅ تم تعديل نقطة العمل بنجاح'
      
      setTimeout(() => {
        emit('worksite-updated')
        emit('close')
      }, 1500)
    } else {
      // وضع الإضافة
      await api.post('/worksites', payload)
      success.value = '✅ ' + t('worksite_added_successfully')
      
      setTimeout(() => {
        emit('worksite-added')
        emit('close')
      }, 1500)
    }
  } catch (err) {
    error.value = err.response?.data?.error || '❌ ' + t('save_failed')
    console.error('خطأ:', err)
  } finally {
    isSubmitting.value = false
  }
}
</script>

<style scoped>
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

.modal {
  width: 100%;
  max-width: 520px;
  padding: 0;
  max-height: 90vh;
  overflow-y: auto;
  background: var(--surface);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-xl);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid var(--line);
}

.modal-header h3 {
  font-size: 17px;
  margin: 0;
}

.modal-close {
  background: none;
  border: none;
  font-size: 24px;
  cursor: pointer;
  color: var(--ink-soft);
}

.modal__hint {
  padding: 0 20px 16px;
  font-size: 13px;
  color: var(--ink-soft);
  margin: 0;
}

.modal__form {
  padding: 0 20px 20px;
}

.form-group {
  margin-bottom: 14px;
  position: relative;
}

.form-group label {
  display: block;
  font-size: 13px;
  font-weight: 600;
  color: var(--ink);
  margin-bottom: 4px;
}

.search-wrapper {
  position: relative;
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
  top: 50%;
  transform: translateY(-50%);
  font-size: 16px;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: translateY(-50%) rotate(0deg); }
  to { transform: translateY(-50%) rotate(360deg); }
}

.required {
  color: var(--signal-out);
}

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

.result-type {
  font-size: 10px;
  color: var(--brand);
  background: var(--brand-tint);
  padding: 2px 8px;
  border-radius: 999px;
}

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

.selected-location {
  background: var(--brand-tint);
  border: 1px solid var(--brand);
  border-radius: var(--radius-sm);
  padding: 12px 16px;
  margin-bottom: 14px;
}

.selected-location__header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
  font-weight: 600;
  color: var(--brand);
}

.selected-clear {
  background: none;
  border: none;
  color: var(--signal-out);
  cursor: pointer;
  font-size: 18px;
  padding: 0 4px;
}

.selected-location__details p {
  margin: 4px 0;
  font-size: 13px;
  color: var(--ink);
}

.alert {
  padding: 10px 14px;
  border-radius: var(--radius-sm);
  font-size: 13px;
  margin-bottom: 14px;
}

.alert-error {
  background: var(--signal-out-tint);
  color: var(--signal-out);
}

.alert-success {
  background: var(--signal-in-tint);
  color: var(--signal-in);
}

.form-actions {
  display: flex;
  gap: 10px;
  justify-content: flex-end;
  margin-top: 14px;
  padding-top: 14px;
  border-top: 1px solid var(--line);
}
</style>
