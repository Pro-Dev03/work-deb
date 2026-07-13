<template>
  <router-link :to="`/report/${report.id}`" class="report-card">
    <div>
      <h3>{{ report.title }}</h3>
      <p class="report-card__meta mono">{{ report.date }}</p>
    </div>
    <span class="badge badge--in">{{ report.status }}</span>
  </router-link>
</template>

<script setup>
defineProps({ report: { type: Object, required: true } })
</script>

<style scoped>
.report-card { display: flex; align-items: center; justify-content: space-between; background: var(--surface); border: 1px solid var(--line); border-radius: var(--radius-md); padding: 16px; margin-bottom: 10px; box-shadow: var(--shadow-sm); }
.report-card h3 { font-size: 14.5px; margin-bottom: 4px; }
.report-card__meta { font-size: 12px; color: var(--ink-soft); }
</style>
