<template>
  <div class="photo-upload">
    <label class="photo-upload__drop" :class="{ 'has-image': preview }">
      <img v-if="preview" :src="preview" alt="معاينة الصورة" />
      <template v-else>
        <span class="photo-upload__icon">＋</span>
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
  height: 140px; border: 1.5px dashed var(--line-strong); border-radius: var(--radius-md);
  color: var(--ink-soft); font-size: 13px; cursor: pointer; overflow: hidden; background: var(--surface);
}
.photo-upload__drop.has-image { padding: 0; }
.photo-upload__drop img { width: 100%; height: 100%; object-fit: cover; }
.photo-upload__icon { font-size: 22px; color: var(--brand); }
</style>
