<template>
  <div class="geo-check">
    <p v-if="state === 'idle'" class="geo-check__hint">اضغط لتحديد موقعك الحالي قبل التختيم</p>
    <p v-else-if="state === 'loading'" class="geo-check__hint">جارٍ تحديد موقعك...</p>
    <p v-else-if="state === 'error'" class="geo-check__error">{{ errorMessage }}</p>
    <p v-else class="geo-check__ok mono">تم تحديد الموقع: {{ position.lat.toFixed(5) }}, {{ position.lng.toFixed(5) }}</p>

    <button class="btn btn--ghost btn--block" @click="locate" :disabled="state === 'loading'">
      {{ state === 'ok' ? 'تحديث الموقع' : 'تحديد موقعي' }}
    </button>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { getCurrentPosition } from '../services/geolocation'

const state = ref('idle')
const position = ref(null)
const errorMessage = ref('')
const emit = defineEmits(['located'])

async function locate() {
  state.value = 'loading'
  try {
    position.value = await getCurrentPosition()
    state.value = 'ok'
    emit('located', position.value)
  } catch (e) {
    errorMessage.value = 'تعذر تحديد الموقع، تأكد من تفعيل صلاحية الوصول للموقع الجغرافي'
    state.value = 'error'
  }
}
</script>

<style scoped>
.geo-check { display: flex; flex-direction: column; gap: 10px; }
.geo-check__hint { font-size: var(--text-sm); color: var(--text-secondary); text-align: center; }
.geo-check__error { font-size: var(--text-sm); color: var(--error-600); text-align: center; }
.geo-check__ok { font-size: var(--text-sm); color: var(--success-600); text-align: center; }
</style>
