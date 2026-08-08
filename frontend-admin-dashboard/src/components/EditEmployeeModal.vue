<template>
  <div class="modal-backdrop" @click.self="$emit('close')">
    <div class="modal card">
      <div class="modal-header">
        <h3>✏️ تعديل بيانات الموظف</h3>
        <button class="modal-close" @click="$emit('close')">✕</button>
      </div>

      <form class="modal__form" @submit.prevent="handleSubmit">
        <div v-if="error" class="alert alert-error">{{ error }}</div>
        <div v-if="success" class="alert alert-success">{{ success }}</div>

        <div class="form-group">
          <label>👤 الاسم الكامل <span class="required">*</span></label>
          <input v-model="form.full_name" type="text" required placeholder="أدخل الاسم الكامل" />
        </div>

        <div class="form-group">
          <label>📱 رقم الهاتف</label>
          <input 
            v-model="form.phone" 
            type="tel" 
            placeholder="أدخل رقم الهاتف"
            dir="ltr"
          />
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
import { ref, reactive, onMounted } from 'vue'
import api from '../services/api'

const props = defineProps({
  employee: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['close', 'employee-updated'])

const loading = ref(false)
const error = ref('')
const success = ref('')

const form = reactive({
  full_name: '',
  phone: ''
})

onMounted(() => {
  if (props.employee) {
    form.full_name = props.employee.full_name
    form.phone = props.employee.phone
  }
})

async function handleSubmit() {
  loading.value = true
  error.value = ''
  success.value = ''

  try {
    const payload = {
      full_name: form.full_name.trim(),
      phone: form.phone.trim()
    }

    console.log('📤 إرسال تعديل بيانات الموظف:', payload)

    const { data } = await api.put(`/admin/employees/${props.employee.id}`, payload)

    success.value = `✅ تم تعديل بيانات الموظف بنجاح!`

    // إرسال الحدث فوراً لإعادة تحميل البيانات
    emit('employee-updated')
    
    // إغلاق المودال بعد فترة قصيرة لعرض رسالة النجاح
    setTimeout(() => {
      emit('close')
    }, 1000)
  } catch (err) {
    console.error('❌ فشل التعديل:', err.response?.data)
    error.value = err.response?.data?.error || '❌ فشل تعديل الموظف'
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
  width: 100%; max-width: 420px; padding: 0;
  max-height: 90vh; overflow-y: auto;
}

.modal-header {
  display: flex; justify-content: space-between; align-items: center;
  padding: 20px 24px; border-bottom: 1px solid var(--line);
}

.modal-header h3 { font-size: 18px; margin: 0; }
.modal-close { background: none; border: none; font-size: 24px; cursor: pointer; color: var(--ink-soft); }

.modal__form { padding: 0 24px 24px; }

.form-group { margin-bottom: 14px; }
.form-group label { display: block; font-size: 13px; font-weight: 600; margin-bottom: 4px; }

.form-group input,
.form-group select {
  width: 100%; padding: 10px 12px;
  border: 1.5px solid var(--line); border-radius: var(--radius-sm);
  font-size: 14px; font-family: inherit; background: var(--surface);
  transition: border-color 0.3s;
}

.form-group input:focus,
.form-group select:focus {
  border-color: var(--brand); outline: none;
  box-shadow: 0 0 0 3px var(--brand-tint);
}

.form-group input:disabled {
  background: var(--line);
  cursor: not-allowed;
  opacity: 0.7;
}

.required { color: var(--signal-out); }
.field-hint { font-size: 11px; color: var(--ink-soft); }

.alert {
  padding: 10px 14px; border-radius: var(--radius-sm);
  font-size: 13px; margin-bottom: 14px; white-space: pre-line;
}
.alert-error { background: var(--signal-out-tint); color: var(--signal-out); }
.alert-success { background: var(--signal-in-tint); color: var(--signal-in); }

.form-actions {
  display: flex; gap: 10px;
  justify-content: flex-end; margin-top: 14px;
  padding-top: 14px; border-top: 1px solid var(--line);
}
</style>
