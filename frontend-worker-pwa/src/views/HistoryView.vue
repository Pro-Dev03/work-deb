<template>
  <div class="history-view view">
    <!-- Page Header -->
    <div class="page-header">
      <div class="header-icon">
        <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
          <polyline points="14 2 14 8 20 8"/>
          <line x1="16" y1="13" x2="8" y2="13"/>
          <line x1="16" y1="17" x2="8" y2="17"/>
          <polyline points="10 9 9 9 8 9"/>
        </svg>
      </div>
      <div>
        <h1>{{ t('my_attendance_history') }}</h1>
        <p>{{ t('attendance_history') }}</p>
      </div>
    </div>

    <!-- Filters -->
    <div class="filters">
      <select v-model="selectedYear" @change="fetchMyAttendanceHistory" class="form-select">
        <option v-for="year in availableYears" :key="year" :value="year">{{ year }}</option>
      </select>
      <select v-model="selectedMonth" @change="fetchMyAttendanceHistory" class="form-select">
        <option v-for="month in availableMonths" :key="month.value" :value="month.value">{{ month.label }}</option>
      </select>
    </div>

    <!-- Monthly Summary -->
    <div v-if="myAttendanceHistory.length > 0" class="summary-row">
      <div class="summary-card">
        <div class="s-icon">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <rect x="3" y="4" width="18" height="18" rx="2"/>
            <path d="M16 2v4M8 2v4M3 10h18"/>
          </svg>
        </div>
        <div class="s-value">{{ totalHours.toFixed(1) }}h</div>
        <div class="s-label">{{ t('total_hours') }}</div>
      </div>
      <div class="summary-card">
        <div class="s-icon">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10"/>
            <polyline points="12 6 12 12 16 14"/>
          </svg>
        </div>
        <div class="s-value">{{ workDays }}</div>
        <div class="s-label">{{ t('work_days') }}</div>
      </div>
    </div>

    <!-- History List -->
    <div v-if="loadingMyHistory" class="loading">{{ t('loading') }}</div>
    <div v-else-if="myAttendanceHistory.length === 0" class="empty">{{ t('no_attendance_records') }}</div>
    <div v-else class="history-list">
      <div v-for="record in myAttendanceHistory" :key="record.id" class="history-item">
        <div class="history-date">{{ formatDate(record.check_in_time) }}</div>
        <div class="history-details">
          <div class="history-worksite">{{ record.worksite_name || '—' }}</div>
          <div class="history-times mono">
            {{ formatTime(record.check_in_time) }} - {{ record.check_out_time ? formatTime(record.check_out_time) : '—' }}
          </div>
        </div>
        <div class="history-hours mono">{{ record.worked_hours ? record.worked_hours.toFixed(1) + 'h' : '—' }}</div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useI18n } from '../services/i18n'
import api from '../services/api'

const { t, currentLang } = useI18n()

const myAttendanceHistory = ref([])
const loadingMyHistory = ref(false)
const selectedYear = ref(new Date().getFullYear())
const selectedMonth = ref(String(new Date().getMonth() + 1))
const availableYears = ref([])
const availableMonths = ref([])

const totalHours = computed(() => {
  return myAttendanceHistory.value.reduce((sum, record) => sum + (record.worked_hours || 0), 0)
})

const workDays = computed(() => {
  return myAttendanceHistory.value.length
})

function updateMonthNames() {
  const monthKeys = ['january', 'february', 'march', 'april', 'may', 'june', 'july', 'august', 'september', 'october', 'november', 'december']
  availableMonths.value = monthKeys.map((key, index) => ({
    value: String(index + 1),
    label: t(key)
  }))
}

const currentYear = new Date().getFullYear()
for (let i = currentYear; i >= currentYear - 2; i--) {
  availableYears.value.push(i)
}

updateMonthNames()

watch(currentLang, () => {
  updateMonthNames()
})

async function fetchMyAttendanceHistory() {
  loadingMyHistory.value = true
  try {
    const response = await api.get('/attendance/history', {
      params: {
        year: selectedYear.value,
        month: selectedMonth.value
      }
    })
    const data = response.data
    // الباك-أند يُرجع { data: [...], pagination: {...} }
    const records = Array.isArray(data) ? data : (data?.data || data?.records || [])
    myAttendanceHistory.value = records.filter(record => record !== null)
  } catch (error) {
    console.error('Failed to fetch attendance history:', error)
    myAttendanceHistory.value = []
  } finally {
    loadingMyHistory.value = false
  }
}

function formatDate(dateStr) {
  if (!dateStr) return '—'
  const date = new Date(dateStr)
  return date.toLocaleDateString(currentLang.value === 'ar' ? 'ar-SA' : 'en-US')
}

function formatTime(dateStr) {
  if (!dateStr) return '—'
  const date = new Date(dateStr)
  return date.toLocaleTimeString(currentLang.value === 'ar' ? 'ar-SA' : 'en-US', { hour: '2-digit', minute: '2-digit' })
}

onMounted(() => {
  fetchMyAttendanceHistory()
})
</script>

<style scoped>
.history-view {
  padding: var(--space-4);
}

.page-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: var(--space-4);
}

.header-icon {
  width: 40px;
  height: 40px;
  border-radius: var(--radius-md);
  background: var(--primary-100);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--primary-600);
}

.page-header h1 {
  font-size: var(--text-lg);
  font-weight: var(--font-semibold);
  color: var(--text-primary);
  margin: 0;
}

.page-header p {
  font-size: var(--text-sm);
  color: var(--text-secondary);
  margin: 4px 0 0 0;
}

.filters {
  display: flex;
  gap: var(--space-2);
  margin-bottom: var(--space-4);
}

.form-select {
  flex: 1;
  padding: var(--space-2) var(--space-3);
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  background: var(--surface);
  color: var(--text-primary);
  font-size: var(--text-sm);
  cursor: pointer;
}

.summary-row {
  display: flex;
  gap: var(--space-2);
  margin-bottom: var(--space-4);
}

.summary-card {
  flex: 1;
  padding: var(--space-3);
  border-radius: var(--radius-md);
  background: var(--surface);
  border: 1px solid var(--border);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.s-icon {
  color: var(--primary-500);
}

.s-value {
  font-size: var(--text-lg);
  font-weight: var(--font-semibold);
  color: var(--text-primary);
}

.s-label {
  font-size: var(--text-xs);
  color: var(--text-secondary);
}

.loading {
  text-align: center;
  padding: var(--space-8);
  color: var(--text-secondary);
}

.empty {
  text-align: center;
  padding: var(--space-8);
  color: var(--text-secondary);
}

.history-list {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.history-item {
  padding: var(--space-3);
  border-radius: var(--radius-md);
  background: var(--surface);
  border: 1px solid var(--border);
  display: flex;
  align-items: center;
  gap: var(--space-3);
}

.history-date {
  font-size: var(--text-sm);
  font-weight: var(--font-medium);
  color: var(--text-primary);
  min-width: 80px;
}

.history-details {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.history-worksite {
  font-size: var(--text-sm);
  color: var(--text-primary);
}

.history-times {
  font-size: var(--text-xs);
  color: var(--text-secondary);
}

.history-hours {
  font-size: var(--text-sm);
  font-weight: var(--font-semibold);
  color: var(--primary-600);
  min-width: 50px;
  text-align: right;
}

.mono {
  font-family: monospace;
}
</style>
