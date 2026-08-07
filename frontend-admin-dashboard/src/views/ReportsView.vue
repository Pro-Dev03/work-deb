<template>
  <div>
    <div class="page-head">
      <div>
        <h2>{{ t('reports') }}</h2>
        <p>{{ t('reports_description') }}</p>
      </div>
      <div class="report-filters">
        <select v-model="selectedPeriod" @change="loadReports" class="period-select">
          <option value="today">{{ t('today') }}</option>
          <option value="week">{{ t('this_week') }}</option>
          <option value="month">{{ t('this_month') }}</option>
        </select>
        <button class="btn btn--primary btn--sm" @click="loadReports">
          🔄 {{ t('refresh') }}
        </button>
      </div>
    </div>

    <div v-if="loading" class="loading-state">
      <div class="spinner"></div>
      <p>{{ t('loading_reports') }}</p>
    </div>

    <div v-else>
      <!-- إحصائيات سريعة -->
      <div class="stats-grid">
        <div class="stat-card stat-card--blue">
          <span class="stat-card__icon">
            <Users :size="24" :stroke-width="2" />
          </span>
          <div>
            <span class="stat-card__value">{{ summary.total_employees || 0 }}</span>
            <span class="stat-card__label">{{ t('total_employees') }}</span>
          </div>
        </div>
        <div class="stat-card stat-card--green">
          <span class="stat-card__icon">
            <CheckSquare :size="24" :stroke-width="2" />
          </span>
          <div>
            <span class="stat-card__value">{{ summary.completed_employees || 0 }}</span>
            <span class="stat-card__label">{{ t('completed_today') }}</span>
          </div>
        </div>
        <div class="stat-card stat-card--yellow">
          <span class="stat-card__icon">
            <Clock :size="24" :stroke-width="2" />
          </span>
          <div>
            <span class="stat-card__value">{{ summary.waiting_employees || 0 }}</span>
            <span class="stat-card__label">{{ t('waiting_employees') }}</span>
          </div>
        </div>
        <div class="stat-card stat-card--purple">
          <span class="stat-card__icon">
            <BarChart3 :size="24" :stroke-width="2" />
          </span>
          <div>
            <span class="stat-card__value">{{ summary.in_progress || 0 }}</span>
            <span class="stat-card__label">{{ t('in_progress_tasks') }}</span>
          </div>
        </div>
      </div>

      <!-- توزيع المهام -->
      <div class="card report-card">
        <div class="card-header">
          <h3>{{ t('task_distribution') }}</h3>
          <div class="header-actions">
            <span class="badge badge--info">{{ selectedPeriodText }}</span>
          </div>
        </div>
        <TaskDistributionChart 
          :completed="summary.completed || 0"
          :in_progress="summary.in_progress || 0"
          :pending="summary.pending || 0"
          :late="summary.late || 0"
        />
      </div>

      <!-- الموظفين المكتملين -->
      <div class="card report-card">
        <div class="card-header">
          <h3>{{ t('completed_employees_title') }}</h3>
          <span class="badge badge--in">{{ completedEmployees.length }}</span>
        </div>
        <div v-if="completedEmployees.length === 0" class="empty-state">
          <p>{{ t('no_completed_employees') }}</p>
        </div>
        <div v-else class="employees-table">
          <div class="table-header">
            <span>{{ t('employee_name') }}</span>
            <span>{{ t('worksite') }}</span>
            <span>{{ t('working_hours') }}</span>
            <span>{{ t('check_out') }}</span>
          </div>
          <div v-for="emp in completedEmployees" :key="emp.id" class="table-row">
            <div class="table-cell">
              <span class="cell-label">{{ t('employee_name') }}</span>
              <span class="cell-value employee-name">{{ emp.full_name }}</span>
            </div>
            <div class="table-cell">
              <span class="cell-label">{{ t('worksite') }}</span>
              <span class="cell-value worksite">{{ emp.worksite_name || '—' }}</span>
            </div>
            <div class="table-cell">
              <span class="cell-label">{{ t('working_hours') }}</span>
              <span class="cell-value hours mono">{{ emp.hours_worked ? emp.hours_worked.toFixed(1) + 'h' : '—' }}</span>
            </div>
            <div class="table-cell">
              <span class="cell-label">{{ t('check_out') }}</span>
              <span class="cell-value time mono">{{ formatTime(emp.check_out_time) }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- الموظفين قيد الانتظار -->
      <div class="card report-card">
        <div class="card-header">
          <h3>{{ t('waiting_employees_title') }}</h3>
          <span class="badge badge--warning">{{ pendingEmployees.length }}</span>
        </div>
        <div v-if="pendingEmployees.length === 0" class="empty-state">
          <p>{{ t('no_waiting_employees') }}</p>
        </div>
        <div v-else class="employees-table">
          <div class="table-header table-header--pending">
            <span>{{ t('employee_name') }}</span>
            <span>{{ t('phone') }}</span>
          </div>
          <div v-for="emp in pendingEmployees" :key="emp.id" class="table-row table-row--pending">
            <div class="table-cell">
              <span class="cell-label">{{ t('employee_name') }}</span>
              <span class="cell-value employee-name">{{ emp.full_name }}</span>
            </div>
            <div class="table-cell">
              <span class="cell-label">{{ t('phone') }}</span>
              <span class="cell-value phone mono">{{ emp.phone || '—' }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from '../services/i18n'
import TaskDistributionChart from '../components/TaskDistributionChart.vue'
import api from '../services/api'
import { Users, CheckSquare, Clock, BarChart3 } from '@lucide/vue'

const { t } = useI18n()

const loading = ref(false)
const selectedPeriod = ref('today')
const summary = ref({
  total_employees: 0,
  completed_employees: 0,
  waiting_employees: 0,
  completed: 0,
  in_progress: 0,
  pending: 0,
  late: 0
})
const completedEmployees = ref([])
const pendingEmployees = ref([])

const selectedPeriodText = computed(() => {
  const texts = {
    today: t('today'),
    week: t('this_week'),
    month: t('this_month')
  }
  return texts[selectedPeriod.value] || t('today')
})

function formatTime(timeString) {
  if (!timeString) return '—'
  try {
    const date = new Date(timeString)
    return date.toLocaleTimeString('ar-SA', { hour: '2-digit', minute: '2-digit' })
  } catch {
    return '—'
  }
}

async function loadReports() {
  loading.value = true
  try {
    // تحميل الملخص اليومي مع معامل الفترة
    const { data: summaryData } = await api.get('/reports/daily-summary', {
      params: { period: selectedPeriod.value }
    })
    summary.value = summaryData

    // تحميل الموظفين المكتملين
    const { data: completedData } = await api.get('/reports/completed-employees')
    completedEmployees.value = completedData || []

    // تحميل الموظفين قيد الانتظار
    const { data: pendingData } = await api.get('/reports/pending-employees')
    pendingEmployees.value = pendingData || []
  } catch (error) {
    console.error('فشل تحميل التقارير:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadReports()
})
</script>

<style scoped>
.page-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  margin-bottom: 18px;
  flex-wrap: wrap;
  gap: 10px;
}

.page-head h2 { font-size: 20px; margin-bottom: 4px; }
.page-head p { font-size: 13px; color: var(--ink-soft); }

.report-filters {
  display: flex;
  gap: 8px;
  align-items: center;
}

.period-select {
  padding: 8px 12px;
  border-radius: var(--radius-sm);
  border: 1px solid var(--line);
  background: var(--surface);
  font-size: 13px;
  cursor: pointer;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
  margin-bottom: 20px;
}

.stat-card {
  background: var(--surface);
  border: 1px solid var(--line);
  border-radius: var(--radius-md);
  padding: 18px 20px;
  display: flex;
  align-items: center;
  gap: 14px;
  transition: all 0.2s ease;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
}

.stat-card__icon {
  width: 44px;
  height: 44px;
  border-radius: var(--radius-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.stat-card__icon :deep(svg) {
  color: inherit;
  stroke: currentColor;
}

.stat-card--blue .stat-card__icon { background: #3b82f620; color: #3b82f6; }
.stat-card--green .stat-card__icon { background: #10b98120; color: #10b981; }
.stat-card--yellow .stat-card__icon { background: #f59e0b20; color: #f59e0b; }
.stat-card--purple .stat-card__icon { background: #8b5cf620; color: #8b5cf6; }

.stat-card__value {
  display: block;
  font-size: 24px;
  font-weight: 700;
  color: var(--ink);
  line-height: 1;
}

.stat-card__label {
  font-size: 12px;
  color: var(--ink-soft);
  margin-top: 4px;
  display: block;
}

.report-card {
  padding: 22px;
  margin-bottom: 16px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.card-header h3 {
  font-size: 15px;
  margin: 0;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.employees-table {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.table-header {
  display: grid;
  grid-template-columns: 2fr 1.5fr 1fr 1fr;
  gap: 12px;
  padding: 12px 16px;
  background: var(--brand-tint);
  border-radius: var(--radius-sm);
  font-size: 12px;
  font-weight: 600;
  color: var(--brand-dark);
}

.table-header--pending {
  grid-template-columns: 2fr 1fr;
}

.table-row {
  display: grid;
  grid-template-columns: 2fr 1.5fr 1fr 1fr;
  gap: 12px;
  padding: 12px 16px;
  background: var(--surface);
  border: 1px solid var(--line);
  border-radius: var(--radius-sm);
  font-size: 13px;
  transition: all 0.2s ease;
}

.table-row--pending {
  grid-template-columns: 2fr 1fr;
}

.table-cell {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.cell-label {
  display: none;
  font-size: 11px;
  font-weight: 600;
  color: var(--ink-soft);
}

.cell-value {
  font-size: 13px;
}

.table-row:hover {
  border-color: var(--brand);
  background: var(--brand-tint);
}

.employee-name {
  font-weight: 500;
  color: var(--ink);
}

.worksite, .phone, .email {
  color: var(--ink-soft);
  font-size: 12px;
}

.hours, .time {
  font-weight: 600;
  color: var(--brand);
}

.mono {
  font-family: monospace;
}

.empty-state {
  text-align: center;
  padding: 40px 20px;
  color: var(--ink-soft);
}

.loading-state {
  text-align: center;
  padding: 60px 20px;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 3px solid var(--line);
  border-top-color: var(--brand);
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto 16px;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

@media (max-width: 1024px) {
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 768px) {
  .page-head {
    flex-direction: column;
    align-items: flex-start;
  }

  .report-filters {
    width: 100%;
  }

  .period-select {
    flex: 1;
  }

  .stats-grid {
    grid-template-columns: 1fr;
  }

  .stat-card {
    padding: 14px 16px;
  }

  .stat-card__icon {
    width: 36px;
    height: 36px;
  }

  .stat-card__value {
    font-size: 20px;
  }

  .stat-card__label {
    font-size: 11px;
  }

  .report-card {
    padding: 16px;
  }

  .card-header h3 {
    font-size: 14px;
  }

  .table-header, .table-row, .table-header--pending, .table-row--pending {
    grid-template-columns: 1fr;
    gap: 8px;
  }

  .table-header {
    display: none;
  }

  .table-row {
    padding: 12px;
    font-size: 13px;
  }

  .table-cell {
    flex-direction: row;
    justify-content: space-between;
    align-items: center;
    padding: 6px 0;
    border-bottom: 1px solid var(--line);
  }

  .table-cell:last-child {
    border-bottom: none;
  }

  .cell-label {
    display: block;
    font-size: 11px;
    font-weight: 600;
    color: var(--ink-soft);
  }

  .cell-value {
    font-size: 13px;
  }

  .employee-name {
    font-size: 14px;
  }

  .worksite, .phone {
    font-size: 12px;
  }

  .hours, .time {
    font-size: 13px;
  }

  .table-cell {
    padding: 4px 0;
  }

  .cell-label {
    font-size: 10px;
  }

  .cell-value {
    font-size: 12px;
  }
}

@media (max-width: 480px) {
  .page-head h2 {
    font-size: 18px;
  }

  .stats-grid {
    gap: 12px;
  }

  .stat-card {
    padding: 12px 14px;
  }

  .stat-card__icon {
    width: 32px;
    height: 32px;
  }

  .stat-card__value {
    font-size: 18px;
  }

  .report-card {
    padding: 14px;
  }

  .table-row {
    padding: 10px;
  }

  .table-cell {
    padding: 3px 0;
  }

  .cell-label {
    font-size: 9px;
  }

  .cell-value {
    font-size: 11px;
  }

  .employee-name {
    font-size: 12px;
  }
}
</style>
