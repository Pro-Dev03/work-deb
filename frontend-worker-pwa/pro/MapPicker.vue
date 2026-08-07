<template>
  <div class="map-picker">
    <div class="map-picker__grid"></div>
    <div class="map-picker__pin">
      <span class="map-picker__pulse"></span>
      <span class="map-picker__dot"></span>
    </div>
    <span class="map-picker__label">موقع نقطة العمل</span>
  </div>
</template>

<script setup>
defineProps({ lat: Number, lng: Number })
</script>

<style scoped>
.map-picker {
  position: relative; height: 160px; border-radius: var(--radius-md); overflow: hidden;
  background: linear-gradient(180deg, var(--primary-50), var(--primary-100)); border: 1px solid var(--border);
  display: flex; align-items: center; justify-content: center;
}
.map-picker__grid {
  position: absolute; inset: 0;
  background-image: linear-gradient(var(--border) 1px, transparent 1px), linear-gradient(90deg, var(--border) 1px, transparent 1px);
  background-size: 24px 24px; opacity: .5;
}
.map-picker__pin { position: relative; width: 20px; height: 20px; }
.map-picker__dot { position: absolute; inset: 0; margin: auto; width: 12px; height: 12px; border-radius: 50%; background: var(--primary-500); border: 2px solid var(--surface); box-shadow: var(--shadow-sm); }
.map-picker__pulse { position: absolute; inset: 0; margin: auto; width: 12px; height: 12px; border-radius: 50%; background: var(--primary-500); opacity: .5; animation: map-pulse 2s ease-out infinite; }
.map-picker__label { position: absolute; bottom: 10px; font-size: var(--text-xs); font-weight: var(--font-semibold); color: var(--text-secondary); background: rgba(255,255,255,.85); padding: 3px 10px; border-radius: var(--radius-full); }
@keyframes map-pulse { 0% { transform: scale(1); opacity: .5; } 100% { transform: scale(3.4); opacity: 0; } }
@media (prefers-reduced-motion: reduce){
  .map-picker__pulse{ animation: none; }
}
</style>
