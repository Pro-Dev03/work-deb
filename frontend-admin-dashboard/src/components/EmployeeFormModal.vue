<template>
  <div class="modal-backdrop" @click.self="$emit('close')">
    <div class="modal card">
      <div class="modal-header">
        <h3>➕ إضافة موظف جديد</h3>
        <button class="modal-close" @click="$emit('close')">✕</button>
      </div>

      <p class="modal__hint">📱 سيتم إنشاء حساب الموظف برقم الهاتف فقط</p>
      <p class="modal__hint-small">⚠️ سجل دخول الموظف سيكون برقم هاتفه وجهازه</p>

      <form class="modal__form" @submit.prevent="handleSubmit">
        <div v-if="error" class="alert alert-error">{{ error }}</div>
        <div v-if="success" class="alert alert-success">{{ success }}</div>

        <div class="form-group">
          <label>👤 الاسم الكامل <span class="required">*</span></label>
          <input v-model="form.full_name" type="text" required placeholder="أدخل الاسم الكامل" />
        </div>

        <div class="form-group">
          <label>📱 رقم الهاتف <span class="required">*</span></label>
          <input 
            v-model="form.phone" 
            type="tel" 
            required 
            placeholder="05xxxxxxxx"
            dir="ltr"
          />
          <span class="field-hint">سيستخدم الموظف هذا الرقم لتسجيل الدخول</span>
        </div>

        <div class="form-group">
          <label>🎯 الدور</label>
          <select v-model="form.role">
            <option value="employee">ميداني</option>
            <option value="admin">مدير</option>
          </select>
        </div>

        <div class="form-actions">
          <button type="button" class="btn btn--ghost" @click="$emit('close')">إلغاء</button>
          <button type="submit" class="btn btn--primary" :disabled="loading">
            {{ loading ? '⏳ جارٍ الحفظ...' : '💾 إنشاء الموظف' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import api from '../services/api'

const emit = defineEmits(['close', 'employee-added'])

const loading = ref(false)
const error = ref('')
const success = ref('')

const form = reactive({
  full_name: '',
  phone: '',
  role: 'employee'
})

async function handleSubmit() {
  loading.value = true
  error.value = ''
  success.value = ''

  // التحقق من رقم الهاتف
  if (!form.phone || form.phone.length < 9) {
    error.value = '❌ الرجاء إدخال رقم هاتف صحيح'
    loading.value = false
    return
  }

  try {
    const payload = {
      full_name: form.full_name,
      phone: form.phone.trim(),
      role: form.role || 'employee'
    }

    console.log('📤 إرسال بيانات الموظف:', payload)

    const { data } = await api.post('/auth/employee-phone', payload)

    success.value = `✅ تم إنشاء الموظف "${data.user.full_name}" بنجاح!`
    success.value += `\n📱 رقم الهاتف: ${data.user.phone}`
    success.value += `\n🔑 سيستخدم هذا الرقم لتسجيل الدخول`

    // إرسال الحدث فوراً لإعادة تحميل البيانات
    emit('employee-added')
    
    // إغلاق المودال بعد فترة قصيرة لعرض رسالة النجاح
    setTimeout(() => {
      emit('close')
    }, 1500)
  } catch (err) {
    console.error('❌ فشل الإضافة:', err.response?.data)
    error.value = err.response?.data?.error || '❌ فشل إنشاء الموظف'
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

.modal__hint {
  padding: 0 24px 8px;
  font-size: 14px; font-weight: 600; color: var(--brand); margin: 0;
}

.modal__hint-small {
  padding: 0 24px 16px;
  font-size: 12px; color: var(--ink-soft); margin: 0;
}

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
