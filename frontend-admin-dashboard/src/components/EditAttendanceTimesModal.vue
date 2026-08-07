<template>
  <div class="modal-backdrop" @click.self="$emit('close')">
    <div class="modal card">
      <div class="modal-header">
        <h3>✏️ تعديل أوقات الحضور</h3>
        <button class="modal-close" @click="$emit('close')">✕</button>
      </div>

      <form class="modal__form" @submit.prevent="handleSubmit">
        <div v-if="error" class="alert alert-error">{{ error }}</div>
        <div v-if="success" class="alert alert-success">{{ success }}</div>

        <div class="form-group">
          <label>الموظف</label>
          <input 
            v-model="attendance.employee_name" 
            type="text" 
            disabled
          />
        </div>

        <div class="form-group">
          <label>نقطة العمل</label>
          <input 
            v-model="attendance.worksite_name" 
            type="text" 
            disabled
          />
        </div>

        <div class="form-group">
          <label>🕐 وقت البدء (Check In) <span class="required">*</span></label>
          <input 
            v-model="form.check_in_time" 
            type="datetime-local" 
            required
          />
          <span class="field-hint">الوقت الفعلي للبدء بالعمل</span>
        </div>

        <div class="form-group">
          <label>🕕 وقت الانتهاء (Check Out) <span class="required">*</span></label>
          <input 
            v-model="form.check_out_time" 
            type="datetime-local" 
            required
          />
          <span class="field-hint">الوقت الفعلي لإنهاء العمل</span>
        </div>

        <div class="info-box">
          <span class="info-label">⏱️ الساعات المحسوبة:</span>
          <span class="info-value">{{ calculatedHours.toFixed(1) }} ساعة</span>
        </div>

        <div class="form-actions">
          <button type="button" class="btn btn--ghost" @click="$emit('close')">إلغاء</button>
          <button type="submit" class="btn btn--primary" :disabled="loading">
            {{ loading ? '⏳ جارٍ الحفظ...' : '💾 حفظ التعديلات' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import api from '../services/api'

const props = defineProps({
  attendance: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['close', 'attendance-updated'])

const loading = ref(false)
const error = ref('')
const success = ref('')

const attendance = reactive({
  employee_name: '',
  worksite_name: ''
})

const form = reactive({
  check_in_time: '',
  check_out_time: ''
})

const calculatedHours = computed(() => {
  if (!form.check_in_time || !form.check_out_time) return 0
  
  const checkIn = new Date(form.check_in_time)
  const checkOut = new Date(form.check_out_time)
  
  if (checkOut <= checkIn) return 0
  
  const diffMs = checkOut - checkIn
  const diffHours = diffMs / (1000 * 60 * 60)
  
  return diffHours
})

onMounted(() => {
  if (props.attendance) {
    attendance.employee_name = props.attendance.employee_name
    attendance.worksite_name = props.attendance.worksite_name
    
    // Convert ISO dates to datetime-local format
    if (props.attendance.check_in_time) {
      form.check_in_time = formatDateTimeLocal(props.attendance.check_in_time)
    }
    if (props.attendance.check_out_time) {
      form.check_out_time = formatDateTimeLocal(props.attendance.check_out_time)
    }
  }
})

function formatDateTimeLocal(isoDate) {
  if (!isoDate) return ''
  const date = new Date(isoDate)
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  
  return `${year}-${month}-${day}T${hours}:${minutes}`
}

async function handleSubmit() {
  loading.value = true
  error.value = ''
  success.value = ''

  // Validate times
  if (new Date(form.check_out_time) <= new Date(form.check_in_time)) {
    error.value = '❌ وقت الانتهاء يجب أن يكون بعد وقت البدء'
    loading.value = false
    return
  }

  try {
    // Convert datetime-local to ISO format for API
    const checkInISO = new Date(form.check_in_time).toISOString()
    const checkOutISO = new Date(form.check_out_time).toISOString()

    const payload = {
      check_in_time: checkInISO,
      check_out_time: checkOutISO
    }

    console.log('📤 إرسال تعديل أوقات الحضور:', payload)

    const { data } = await api.put(`/attendance/${props.attendance.id}/times`, payload)

    success.value = `✅ تم تعديل أوقات الحضور بنجاح!`
    success.value += `\n⏱️ الساعات الجديدة: ${calculatedHours.value.toFixed(1)} ساعة`

    setTimeout(() => {
      emit('attendance-updated')
      emit('close')
    }, 1500)
  } catch (err) {
    console.error('❌ فشل التعديل:', err.response?.data)
    error.value = err.response?.data?.error || '❌ فشل تعديل أوقات الحضور'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.modal-backdrop {
  position: fixed; inset: 0;
  background: rgba(0,0,0,0.6);
  backdrop-filter: blur(4px);
  display: flex; align-items: center; justify-content: center;
  z-index: 1000; padding: 20px;
}

.modal {
  width: 100%; max-width: 500px; padding: 0;
  max-height: 90vh; overflow-y: auto;
}

.modal-header {
  display: flex; justify-content: space-between; align-items: center;
  padding: 20px 24px; border-bottom: 1px solid var(--line);
}

.modal-header h3 { font-size: 18px; margin: 0; }
.modal-close { background: none; border: none; font-size: 24px; cursor: pointer; color: var(--ink-soft); }

.modal__form { padding: 0 24px 24px; }

.form-group { margin-bottom: 16px; }
.form-group label { display: block; font-size: 13px; font-weight: 600; margin-bottom: 6px; }

.form-group input {
  width: 100%; padding: 10px 12px;
  border: 1.5px solid var(--line); border-radius: var(--radius-sm);
  font-size: 14px; font-family: inherit; background: var(--surface);
  transition: border-color 0.3s;
}

.form-group input:focus {
  border-color: var(--brand); outline: none;
  box-shadow: 0 0 0 3px var(--brand-tint);
}

.form-group input:disabled {
  background: var(--line);
  cursor: not-allowed;
  opacity: 0.7;
}

.required { color: var(--signal-out); }
.field-hint { font-size: 11px; color: var(--ink-soft); margin-top: 4px; display: block; }

.info-box {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px;
  background: var(--brand-tint);
  border-radius: var(--radius-sm);
  margin-bottom: 16px;
}

.info-label {
  font-size: 13px;
  font-weight: 600;
  color: var(--brand);
}

.info-value {
  font-size: 16px;
  font-weight: 700;
  color: var(--brand);
}

.alert {
  padding: 10px 14px; border-radius: var(--radius-sm);
  font-size: 13px; margin-bottom: 14px; white-space: pre-line;
}
.alert-error { background: var(--signal-out-tint); color: var(--signal-out); }
.alert-success { background: var(--signal-in-tint); color: var(--signal-in); }

.form-actions {
  display: flex; gap: 10px;
  justify-content: flex-end; margin-top: 16px;
  padding-top: 16px; border-top: 1px solid var(--line);
}
</style>
