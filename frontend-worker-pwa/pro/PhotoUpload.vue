<template>
  <div class="photo-upload">
    <label class="photo-upload__drop" :class="{ 'has-image': preview }">
      <img v-if="preview" :src="preview" alt="معاينة الصورة" />
      <template v-else>
        <svg class="photo-upload__icon" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M12 5v14M5 12h14"/>
        </svg>
        <span>أرفق صورة إثبات إنجاز المهمة</span>
      </template>
      <input type="file" accept="image/*" capture="environment" class="visually-hidden" @change="onChange" />
    </label>
    <button v-if="preview" class="btn btn--ghost btn--sm" @click="clear">إزالة الصورة</button>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const preview = ref(null)
const emit = defineEmits(['selected'])

function onChange(e) {
  const file = e.target.files?.[0]
  if (!file) return
  const reader = new FileReader()
  reader.onload = () => {
    preview.value = reader.result
    emit('selected', file)
  }
  reader.readAsDataURL(file)
}

function clear() {
  preview.value = null
  emit('selected', null)
}
</script>

<style scoped>
.photo-upload { display: flex; flex-direction: column; gap: 10px; align-items: stretch; }
.photo-upload__drop {
  display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 8px;
  height: 140px; border: 1.5px dashed var(--border-strong); border-radius: var(--radius-md);
  color: var(--text-secondary); font-size: var(--text-sm); cursor: pointer; overflow: hidden; background: var(--surface);
  transition: border-color .2s ease, background .2s ease;
}
.photo-upload__drop:hover { border-color: var(--primary-500); background: var(--primary-50); }
.photo-upload__drop.has-image { padding: 0; }
.photo-upload__drop img { width: 100%; height: 100%; object-fit: cover; }
.photo-upload__icon { color: var(--primary-500); }
</style>
