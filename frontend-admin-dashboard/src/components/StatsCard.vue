<template>
  <div class="stat-card">
    <span class="stat-card__label">{{ label }}</span>
    <span class="stat-card__value mono">{{ value }}</span>
    <span v-if="hint" class="stat-card__hint" :class="tone">{{ hint }}</span>
  </div>
</template>

<script setup>
defineProps({
  label: String,
  value: [String, Number],
  hint: String,
  tone: { type: String, default: '' }, // 'up' | 'down'
})
</script>

<style scoped>
.stat-card { background: var(--surface); border: 1px solid var(--line); border-radius: var(--radius-lg); padding: 20px; display: flex; flex-direction: column; gap: 6px; box-shadow: var(--shadow-sm); }
.stat-card__label { font-size: 13px; color: var(--ink-soft); }
.stat-card__value { font-size: 28px; font-weight: 600; }
.stat-card__hint { font-size: 12px; font-weight: 600; color: var(--ink-soft); }
.stat-card__hint.up { color: var(--signal-in); }
.stat-card__hint.down { color: var(--signal-out); }
</style>
