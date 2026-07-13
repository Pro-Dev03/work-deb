<template>
  <div class="live-map">
    <div class="live-map__grid"></div>
    <span v-for="p in pins" :key="p.id" class="live-map__pin" :class="p.status" :style="{ top: p.top + '%', right: p.right + '%' }"></span>
    <div class="live-map__legend">
      <span><i class="in"></i> داخل النطاق</span>
      <span><i class="out"></i> خارج النطاق</span>
    </div>
  </div>
</template>

<script setup>
defineProps({
  pins: {
    type: Array,
    default: () => [
      { id: 1, top: 30, right: 25, status: 'in' },
      { id: 2, top: 55, right: 60, status: 'in' },
      { id: 3, top: 70, right: 40, status: 'out' },
      { id: 4, top: 20, right: 70, status: 'in' },
    ],
  },
})
</script>

<style scoped>
.live-map { position: relative; height: 320px; border-radius: var(--radius-lg); overflow: hidden; background: linear-gradient(180deg, #EAF1EC, #E1EBE4); border: 1px solid var(--line); }
.live-map__grid { position: absolute; inset: 0; background-image: linear-gradient(var(--line) 1px, transparent 1px), linear-gradient(90deg, var(--line) 1px, transparent 1px); background-size: 32px 32px; opacity: .5; }
.live-map__pin { position: absolute; width: 14px; height: 14px; border-radius: 50%; border: 2px solid #fff; box-shadow: var(--shadow-sm); }
.live-map__pin.in { background: var(--signal-in); }
.live-map__pin.out { background: var(--signal-out); }
.live-map__legend { position: absolute; bottom: 14px; right: 14px; display: flex; gap: 14px; background: rgba(255,255,255,.9); padding: 8px 14px; border-radius: 999px; font-size: 12px; color: var(--ink-soft); }
.live-map__legend i { display: inline-block; width: 8px; height: 8px; border-radius: 50%; margin-inline-end: 5px; }
.live-map__legend i.in { background: var(--signal-in); }
.live-map__legend i.out { background: var(--signal-out); }
</style>
