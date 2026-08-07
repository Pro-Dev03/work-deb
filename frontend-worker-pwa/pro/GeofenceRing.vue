<template>
  <div class="gf-ring" :class="status">
    <svg viewBox="0 0 160 160" class="gf-svg" aria-hidden="true">
      <circle cx="80" cy="80" r="68" class="gf-track" />
      <circle
        cx="80" cy="80" r="68" class="gf-fill"
        :stroke-dasharray="circumference"
        :stroke-dashoffset="dashOffset"
      />
      <circle cx="80" cy="80" r="46" class="gf-hole" />
    </svg>
    <div class="gf-center">
      <span class="gf-distance mono">{{ Math.round(distance) }}<small>م</small></span>
      <span class="gf-label">{{ statusText }}</span>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

// distance: مسافة الموظف الحالية عن نقطة العمل (متر)
// radius: نصف القطر المسموح به لنقطة العمل (متر)
const props = defineProps({
  distance: { type: Number, default: 0 },
  radius: { type: Number, default: 100 },
})

const circumference = 2 * Math.PI * 68

const ratio = computed(() => Math.min(props.distance / (props.radius || 1), 1))
const dashOffset = computed(() => circumference * (1 - ratio.value))
const status = computed(() => (props.distance <= props.radius ? 'inside' : 'outside'))
const statusText = computed(() => (status.value === 'inside' ? 'داخل نطاق الموقع' : 'خارج نطاق الموقع'))
</script>

<style scoped>
.gf-ring { position: relative; width: 180px; height: 180px; margin: 0 auto; }
.gf-svg { width: 100%; height: 100%; transform: rotate(-90deg); }
.gf-track { fill: none; stroke: var(--border); stroke-width: 10; }
.gf-fill { fill: none; stroke-width: 10; stroke-linecap: round; transition: stroke-dashoffset .6s ease, stroke .3s ease; }
.gf-ring.inside .gf-fill { stroke: var(--success-500); }
.gf-ring.outside .gf-fill { stroke: var(--error-500); }
.gf-hole { fill: var(--surface); }
.gf-center { position: absolute; inset: 0; display: flex; flex-direction: column; align-items: center; justify-content: center; }
.gf-distance { font-size: var(--text-2xl); font-weight: var(--font-semibold); color: var(--text-primary); }
.gf-distance small { font-size: var(--text-sm); font-weight: var(--font-medium); margin-inline-start: 2px; color: var(--text-secondary); }
.gf-label { margin-top: 4px; font-size: var(--text-sm); font-weight: var(--font-semibold); }
.gf-ring.inside .gf-label { color: var(--success-600); }
.gf-ring.outside .gf-label { color: var(--error-600); }
.gf-ring.inside::after {
  content: ''; position: absolute; inset: -6px; border-radius: 50%; border: 1px solid var(--success-500);
  animation: gf-pulse 2.4s ease-out infinite;
}
@keyframes gf-pulse {
  0% { transform: scale(.94); opacity: .7; }
  100% { transform: scale(1.14); opacity: 0; }
}
@media (prefers-reduced-motion: reduce){
  .gf-ring.inside::after{ animation: none; }
}
</style>
