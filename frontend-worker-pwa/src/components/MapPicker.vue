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
  background: linear-gradient(180deg, #EAF1EC, #E1EBE4); border: 1px solid var(--line);
  display: flex; align-items: center; justify-content: center;
}
.map-picker__grid {
  position: absolute; inset: 0;
  background-image: linear-gradient(var(--line) 1px, transparent 1px), linear-gradient(90deg, var(--line) 1px, transparent 1px);
  background-size: 24px 24px; opacity: .5;
}
.map-picker__pin { position: relative; width: 20px; height: 20px; }
.map-picker__dot { position: absolute; inset: 0; margin: auto; width: 12px; height: 12px; border-radius: 50%; background: var(--brand); border: 2px solid #fff; box-shadow: var(--shadow-sm); }
.map-picker__pulse { position: absolute; inset: 0; margin: auto; width: 12px; height: 12px; border-radius: 50%; background: var(--brand); opacity: .5; animation: map-pulse 2s ease-out infinite; }
.map-picker__label { position: absolute; bottom: 10px; font-size: 12px; font-weight: 600; color: var(--ink-soft); background: rgba(255,255,255,.85); padding: 3px 10px; border-radius: 999px; }
@keyframes map-pulse { 0% { transform: scale(1); opacity: .5; } 100% { transform: scale(3.4); opacity: 0; } }
</style>
