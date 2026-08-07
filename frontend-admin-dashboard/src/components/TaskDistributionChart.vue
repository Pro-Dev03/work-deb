<template>
  <div class="dist">
    <div class="dist__bar">
      <span v-for="seg in segments" :key="seg.key" :style="{ width: seg.pct + '%', background: seg.color }"></span>
    </div>
    <ul class="dist__legend">
      <li v-for="seg in segments" :key="seg.key">
        <span class="dist__dot" :style="{ background: seg.color }"></span>
        {{ seg.label }} <b class="mono">{{ seg.value }}</b>
      </li>
    </ul>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  completed: { type: Number, default: 0 },
  in_progress: { type: Number, default: 0 },
  pending: { type: Number, default: 0 },
  late: { type: Number, default: 0 },
})

const segments = computed(() => {
  const total = props.completed + props.in_progress + props.pending + props.late || 1
  const defs = [
    { key: 'completed', label: 'مكتملة', value: props.completed, color: 'var(--signal-in)' },
    { key: 'in_progress', label: 'جارية', value: props.in_progress, color: 'var(--gold)' },
    { key: 'pending', label: 'قيد الانتظار', value: props.pending, color: 'var(--line-strong)' },
    { key: 'late', label: 'متأخرة', value: props.late, color: 'var(--signal-out)' },
  ]
  return defs.map((d) => ({ ...d, pct: total > 0 ? (d.value / total) * 100 : 0 }))
})
</script>

<style scoped>
.dist__bar { display: flex; height: 12px; border-radius: 999px; overflow: hidden; background: var(--line); }
.dist__legend { list-style: none; padding: 0; margin: 16px 0 0; display: flex; flex-wrap: wrap; gap: 14px; font-size: 13px; color: var(--ink-soft); }
.dist__legend li { display: flex; align-items: center; gap: 6px; }
.dist__dot { width: 8px; height: 8px; border-radius: 50%; }
</style>
