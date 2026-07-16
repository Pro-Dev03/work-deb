<template>
  <div>
    <div class="page-head"><h2>{{ t('reports') }}</h2><p>{{ t('reports_description') }}</p></div>

    <div class="card report-card">
      <TaskDistributionChart v-bind="summary" />
    </div>

    <div class="card report-card">
      <h3>{{ t('range_compliance') }}</h3>
      <p class="report-card__hint">{{ t('range_compliance_hint') }}</p>
      <div class="report-bar">
        <span class="report-bar__accepted" style="width: 92%"></span>
      </div>
      <div class="report-bar__legend">
        <span><i class="in"></i> {{ t('accepted') }} 92٪</span>
        <span><i class="out"></i> {{ t('rejected') }} 8٪</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import TaskDistributionChart from '../components/TaskDistributionChart.vue'
import { useI18n } from '../services/i18n'

const { t } = useI18n()
const summary = { completed: 27, in_progress: 6, pending: 9, late: 2 }
</script>

<style scoped>
.page-head { margin-bottom: 18px; }
.page-head h2 { font-size: 20px; margin-bottom: 4px; }
.page-head p { font-size: 13px; color: var(--ink-soft); }
.report-card { padding: 22px; margin-bottom: 16px; }
.report-card h3 { font-size: 15px; margin-bottom: 6px; }
.report-card__hint { font-size: 13px; color: var(--ink-soft); margin-bottom: 14px; }
.report-bar { height: 12px; border-radius: 999px; background: var(--signal-out-tint); overflow: hidden; }
.report-bar__accepted { display: block; height: 100%; background: var(--signal-in); }
.report-bar__legend { display: flex; gap: 16px; margin-top: 12px; font-size: 13px; color: var(--ink-soft); }
.report-bar__legend i { display: inline-block; width: 8px; height: 8px; border-radius: 50%; margin-inline-end: 5px; }
.report-bar__legend i.in { background: var(--signal-in); }
.report-bar__legend i.out { background: var(--signal-out); }

@media (max-width: 768px) {
  .report-bar__legend {
    flex-direction: column;
    gap: 8px;
  }
}
</style>
