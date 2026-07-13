<template>
  <button class="maps-btn" @click="openInMaps" :title="t('open_in_maps')">
    <span class="maps-icon">🗺️</span>
    <span class="maps-text">{{ t('open_in_maps') }}</span>
  </button>
</template>

<script setup>
import { useI18n } from '../plugins/i18n'

const { t } = useI18n()

const props = defineProps({
  latitude: { type: Number, required: true },
  longitude: { type: Number, required: true },
  label: { type: String, default: 'الموقع' }
})

function openInMaps() {
  const url = `https://www.google.com/maps/dir/?api=1&destination=${props.latitude},${props.longitude}&destination_place_id=${encodeURIComponent(props.label)}`
  
  // فتح في نافذة جديدة
  window.open(url, '_blank')
}
</script>

<style scoped>
.maps-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 14px;
  border: 1px solid var(--brand);
  border-radius: var(--radius-sm);
  background: var(--brand-tint);
  color: var(--brand);
  font-family: var(--font-body);
  font-size: 12px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.maps-btn:hover {
  background: var(--brand);
  color: white;
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
}

.maps-icon {
  font-size: 14px;
}
</style>
