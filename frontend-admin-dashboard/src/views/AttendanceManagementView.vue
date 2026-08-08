<template>
  <div class="template-wrapper">
    <div class="page-head">
      <div>
        <h2>📊 {{ t('attendance_management_title') }}</h2>
        <p>{{ t('attendance_management_desc') }}</p>
      </div>
      <div class="page-head-actions">
        <button class="btn btn--primary" @click="exportToExcel" :disabled="attendanceRecords.length === 0">
          <FileSpreadsheet :size="18" />
          {{ t('export') }}
        </button>
        <button class="btn btn--ghost" @click="fetchAttendanceRecords">🔄 {{ t('refresh') }}</button>
      </div>
    </div>

    <!-- الإحصائيات السريعة -->
    <div v-if="stats" class="stats-grid">
      <div class="stat-card">
        <div class="stat-card__icon stat-card__icon--blue">
          <Users :size="24" />
        </div>
        <div class="stat-card__content">
          <span class="stat-card__label">{{ t('total_records') }}</span>
          <span class="stat-card__value">{{ stats.totalRecords }}</span>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-card__icon stat-card__icon--green">
          <Clock :size="24" />
        </div>
        <div class="stat-card__content">
          <span class="stat-card__label">{{ t('total_hours_display') }}</span>
          <span class="stat-card__value">{{ stats.totalHours.toFixed(1) }} {{ t('hour') }}</span>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-card__icon stat-card__icon--purple">
          <CheckCircle :size="24" />
        </div>
        <div class="stat-card__content">
          <span class="stat-card__label">{{ t('completed_records') }}</span>
          <span class="stat-card__value">{{ stats.completedRecords }}</span>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-card__icon stat-card__icon--orange">
          <AlertCircle :size="24" />
        </div>
        <div class="stat-card__content">
          <span class="stat-card__label">{{ t('active_records') }}</span>
          <span class="stat-card__value">{{ stats.activeRecords }}</span>
        </div>
      </div>
    </div>

    <!-- الفلاتر والبحث -->
    <div class="card filters-section">
      <div class="filters-header">
        <div class="search-box">
          <Search :size="18" />
          <input
            v-model="searchQuery"
            type="text"
            :placeholder="t('search_placeholder')"
            class="search-input"
          />
        </div>
        <button
          class="btn btn--ghost btn--sm"
          @click="resetFilters"
          v-if="hasActiveFilters"
        >
          🔄 {{ t('reset_filters') }}
        </button>
      </div>
      <div class="filters-grid">
        <div class="filter-group">
          <label>{{ t('employee') }}</label>
          <select v-model="filters.employee_id" @change="fetchAttendanceRecords" class="form-select">
            <option value="">{{ t('all_employees') }}</option>
            <option v-for="emp in employees" :key="emp?.id || emp?.full_name || Math.random()" :value="emp?.id">{{ emp?.full_name || '—' }}</option>
          </select>
        </div>
        <div class="filter-group">
          <label>{{ t('worksite') }}</label>
          <select v-model="filters.worksite_id" @change="fetchAttendanceRecords" class="form-select">
            <option value="">{{ t('all_worksites') }}</option>
            <option v-for="ws in worksites" :key="ws.id" :value="ws.id">{{ ws.name }}</option>
          </select>
        </div>
        <div class="filter-group">
          <label>{{ t('from_date') }}</label>
          <input v-model="filters.date_from" type="date" @change="fetchAttendanceRecords" class="form-input" />
        </div>
        <div class="filter-group">
          <label>{{ t('to_date') }}</label>
          <input v-model="filters.date_to" type="date" @change="fetchAttendanceRecords" class="form-input" />
        </div>
      </div>
    </div>

    <div v-if="loading" class="empty-state">
      <p>⏳ {{ t('loading_data') }}</p>
    </div>

    <div v-else-if="filteredRecords.length === 0" class="empty-state">
      <h3>📭 {{ t('no_attendance_records_found') }}</h3>
      <p>{{ t('change_filters') }}</p>
    </div>

    <div v-else class="card">
      <!-- جدول متجاوب -->
      <div class="table-responsive">
        <table class="table">
          <thead>
            <tr>
              <th @click="sortBy('employee_name')" class="sortable">
                {{ t('employee_name') }}
                <span class="sort-icon" v-if="sortField === 'employee_name'">
                  {{ sortOrder === 'asc' ? '↑' : '↓' }}
                </span>
              </th>
              <th @click="sortBy('worksite_name')" class="sortable">
                {{ t('worksite') }}
                <span class="sort-icon" v-if="sortField === 'worksite_name'">
                  {{ sortOrder === 'asc' ? '↑' : '↓' }}
                </span>
              </th>
              <th @click="sortBy('check_in_time')" class="sortable">
                {{ t('date') }}
                <span class="sort-icon" v-if="sortField === 'check_in_time'">
                  {{ sortOrder === 'asc' ? '↑' : '↓' }}
                </span>
              </th>
              <th>{{ t('check_in_time') }}</th>
              <th>{{ t('check_out_time') }}</th>
              <th @click="sortBy('worked_hours')" class="sortable">
                {{ t('hours_display') }}
                <span class="sort-icon" v-if="sortField === 'worked_hours'">
                  {{ sortOrder === 'asc' ? '↑' : '↓' }}
                </span>
              </th>
              <th>{{ t('status') }}</th>
              <th>{{ t('actions') }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="record in paginatedRecords" :key="record.id" class="table-row">
              <td>
                <div class="table__person">
                  <span class="table__avatar">
                    <User :size="20" />
                  </span>
                  <span class="table__name">{{ record.employee_name }}</span>
                </div>
              </td>
              <td>
                <span class="table__worksite">{{ record.worksite_name || '—' }}</span>
              </td>
              <td class="mono">{{ formatDate(record.check_in_time) }}</td>
              <td class="mono">{{ formatTime(record.check_in_time) }}</td>
              <td class="mono">{{ record.check_out_time ? formatTime(record.check_out_time) : '—' }}</td>
              <td class="mono">
                <span class="hours-badge">{{ record.worked_hours ? record.worked_hours.toFixed(1) + ' ' + t('hour') : '—' }}</span>
              </td>
              <td>
                <span class="badge" :class="getStatusBadgeClass(record.status)">
                  {{ getStatusText(record.status) }}
                </span>
              </td>
              <td>
                <div class="table__actions">
                  <button
                    class="btn btn--primary btn--sm"
                    @click="openEditModal(record)"
                  >
                    ✏️ {{ t('edit') }}
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- الصفحات -->
      <div v-if="totalPages > 1" class="pagination">
        <button
          class="btn btn--ghost btn--sm"
          @click="currentPage--"
          :disabled="currentPage === 1"
        >
          {{ t('previous') }}
        </button>
        <span class="pagination-info">
          {{ t('page_of') }} {{ currentPage }} {{ t('of') }} {{ totalPages }}
        </span>
        <button
          class="btn btn--ghost btn--sm"
          @click="currentPage++"
          :disabled="currentPage === totalPages"
        >
          {{ t('next') }}
        </button>
      </div>
    </div>

    <!-- مودال تعديل الأوقات -->
    <EditAttendanceTimesModal
      v-if="showEditModal"
      :attendance="selectedRecord"
      @close="showEditModal = false"
      @attendance-updated="fetchAttendanceRecords"
    />
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import api from '../services/api'
import EditAttendanceTimesModal from '../components/EditAttendanceTimesModal.vue'
import { useI18n } from '../services/i18n'
import { Users, Clock, CheckCircle, AlertCircle, Search, User, FileSpreadsheet } from '@lucide/vue'
import * as XLSX from 'xlsx'

const { t, currentLang } = useI18n()

const loading = ref(false)
const attendanceRecords = ref([])
const employees = ref([])
const worksites = ref([])
const showEditModal = ref(false)
const selectedRecord = ref(null)

const filters = reactive({
  employee_id: '',
  worksite_id: '',
  date_from: '',
  date_to: ''
})

const searchQuery = ref('')
const sortField = ref('check_in_time')
const sortOrder = ref('desc')
const currentPage = ref(1)
const itemsPerPage = 10

// الإحصائيات
const stats = computed(() => {
  if (attendanceRecords.value.length === 0) return null
  
  const totalRecords = attendanceRecords.value.length
  const completedRecords = attendanceRecords.value.filter(r => r.status === 'completed').length
  const activeRecords = attendanceRecords.value.filter(r => r.status === 'active').length
  const totalHours = attendanceRecords.value.reduce((sum, r) => sum + (r.worked_hours || 0), 0)
  
  return {
    totalRecords,
    completedRecords,
    activeRecords,
    totalHours
  }
})

// الفلترة والبحث
const filteredRecords = computed(() => {
  let records = [...attendanceRecords.value]
  
  // البحث النصي
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    records = records.filter(record =>
      (record.employee_name || '').toLowerCase().includes(query) ||
      (record.worksite_name || '').toLowerCase().includes(query)
    )
  }
  
  // الترتيب
  if (sortField.value) {
    records.sort((a, b) => {
      let aVal = a[sortField.value]
      let bVal = b[sortField.value]
      
      if (sortField.value === 'worked_hours') {
        aVal = aVal || 0
        bVal = bVal || 0
      }
      
      if (aVal < bVal) return sortOrder.value === 'asc' ? -1 : 1
      if (aVal > bVal) return sortOrder.value === 'asc' ? 1 : -1
      return 0
    })
  }
  
  return records
})

// التحقق من وجود فلاتر نشطة
const hasActiveFilters = computed(() => {
  return searchQuery.value || 
         filters.employee_id || 
         filters.worksite_id || 
         filters.date_from || 
         filters.date_to
})

// إجمالي الصفحات
const totalPages = computed(() => {
  return Math.ceil(filteredRecords.value.length / itemsPerPage)
})

// السجلات في الصفحة الحالية
const paginatedRecords = computed(() => {
  const start = (currentPage.value - 1) * itemsPerPage
  const end = start + itemsPerPage
  return filteredRecords.value.slice(start, end)
})

async function fetchEmployees() {
  try {
    const { data } = await api.get('/admin/employees')
    employees.value = data.data || data || []
  } catch (error) {
    console.error('❌ فشل جلب الموظفين:', error)
  }
}

async function fetchWorksites() {
  try {
    const { data } = await api.get('/worksites')
    worksites.value = data.data || data || []
  } catch (error) {
    console.error('❌ فشل جلب نقاط العمل:', error)
  }
}

async function fetchAttendanceRecords() {
  loading.value = true
  try {
    const params = new URLSearchParams()
    if (filters.employee_id) params.append('employee_id', filters.employee_id)
    if (filters.worksite_id) params.append('worksite_id', filters.worksite_id)
    if (filters.date_from) params.append('date_from', filters.date_from)
    if (filters.date_to) params.append('date_to', filters.date_to)

    const { data } = await api.get(`/attendance/management?${params.toString()}`, { skipCache: true })
    attendanceRecords.value = data || []
  } catch (error) {
    console.error('❌ فشل جلب سجلات الحضور:', error)
    if (error.response?.status === 404) {
      // Endpoint might not exist yet, show empty state
      attendanceRecords.value = []
    } else {
      attendanceRecords.value = []
    }
  } finally {
    loading.value = false
  }
}

function openEditModal(record) {
  selectedRecord.value = record
  showEditModal.value = true
}

function formatDate(date) {
  if (!date) return '—'
  return new Date(date).toLocaleDateString('en-GB')
}

function formatTime(date) {
  if (!date) return '—'
  return new Date(date).toLocaleTimeString('en-GB', { hour: '2-digit', minute: '2-digit', hour12: false })
}

function sortBy(field) {
  if (sortField.value === field) {
    sortOrder.value = sortOrder.value === 'asc' ? 'desc' : 'asc'
  } else {
    sortField.value = field
    sortOrder.value = 'asc'
  }
}

function resetFilters() {
  searchQuery.value = ''
  filters.employee_id = ''
  filters.worksite_id = ''
  filters.date_from = ''
  filters.date_to = ''
  currentPage.value = 1
  fetchAttendanceRecords()
}

function exportToExcel() {
  if (filteredRecords.value.length === 0) return

  // إنشاء ورقة عمل جديدة
  const wb = XLSX.utils.book_new()
  
  // إضافة معلومات الشركة والتاريخ
  const companyName = 'אבן יסודות'
  const reportTitle = t('attendance_management_title')
  const exportDate = new Date().toLocaleDateString(currentLang.value === 'ar' ? 'ar-SA' : 'en-US')
  const reportPeriod = getReportPeriod()

  // معلومات الإحصائيات - منظم بشكل احترافي وواضح
  const summaryData = [
    [t('company_name'), companyName],
    [t('report_title'), reportTitle],
    [t('export_date'), exportDate],
    [t('report_period'), reportPeriod],
    [''],
    [t('summary'), ''],
    [t('total_records'), stats.value?.totalRecords || 0],
    [t('completed_records'), stats.value?.completedRecords || 0],
    [t('active_records'), stats.value?.activeRecords || 0],
    [t('total_hours_display'), `${stats.value?.totalHours?.toFixed(1) || 0} ${t('hour')}`],
    ['']
  ]

  // رأس الجدول المحسن
  const headers = [
    t('no'),
    t('employee_name'),
    t('worksite'),
    t('date'),
    t('check_in_time'),
    t('check_out_time'),
    t('worked_hours'),
    t('status'),
    t('notes')
  ]

  // صفوف البيانات المحسنة
  const rows = filteredRecords.value.map((record, index) => [
    index + 1,
    record.employee_name || '—',
    record.worksite_name || '—',
    formatDate(record.check_in_time),
    formatTime(record.check_in_time),
    record.check_out_time ? formatTime(record.check_out_time) : '—',
    record.worked_hours ? record.worked_hours.toFixed(2) : '0.00',
    getStatusText(record.status),
    record.notes || '—'
  ])

  // دمج البيانات
  const allData = [...summaryData, headers, ...rows]

  // إنشاء ورقة العمل
  const ws = XLSX.utils.aoa_to_sheet(allData)

  // تعيين عرض الأعمدة بشكل احترافي
  ws['!cols'] = [
    { wch: 25 },  // اسم الحقل
    { wch: 30 },  // القيمة
    { wch: 18 },  // الرقم
    { wch: 25 },  // اسم الموظف
    { wch: 25 },  // موقع العمل
    { wch: 15 },  // التاريخ
    { wch: 12 },  // وقت الحضور
    { wch: 12 },  // وقت الانصراف
    { wch: 12 },  // ساعات العمل
    { wch: 15 },  // الحالة
    { wch: 30 }   // الملاحظات
  ]

  // تنسيق الخلايا
  const headerRowIndex = summaryData.length
  const range = XLSX.utils.decode_range(ws['!ref'])
  
  // تنسيق صف العناوين (الصف بعد الملخص)
  for (let col = range.s.c; col <= range.e.c; col++) {
    const cellAddress = XLSX.utils.encode_cell({ r: headerRowIndex, c: col })
    if (ws[cellAddress]) {
      ws[cellAddress].s = {
        font: { bold: true },
        fill: { fgColor: { rgb: "E3F2FD" } },
        alignment: { horizontal: "center" }
      }
    }
  }

  // تنسيق خلايا الملخص (العمود الأول - أسماء الحقول)
  for (let row = 0; row < summaryData.length; row++) {
    const cellAddress = XLSX.utils.encode_cell({ r: row, c: 0 })
    if (ws[cellAddress] && ws[cellAddress].v) {
      ws[cellAddress].s = {
        font: { bold: true },
        fill: { fgColor: { rgb: "F5F5F5" } },
        alignment: { horizontal: "right" }
      }
    }
    
    // تنسيق العمود الثاني (القيم)
    const cellAddress2 = XLSX.utils.encode_cell({ r: row, c: 1 })
    if (ws[cellAddress2] && ws[cellAddress2].v) {
      ws[cellAddress2].s = {
        font: { bold: false },
        alignment: { horizontal: "left" }
      }
    }
  }

  // إضافة الورقة إلى المصنف
  XLSX.utils.book_append_sheet(wb, ws, 'Attendance Report')

  // إنشاء الملف
  const fileName = `WorkTrack_${reportTitle}_${new Date().toISOString().split('T')[0]}.xlsx`
  XLSX.writeFile(wb, fileName)
}

function getReportPeriod() {
  if (filters.date_from && filters.date_to) {
    return `${formatDate(filters.date_from)} - ${formatDate(filters.date_to)}`
  } else if (filters.date_from) {
    return `${t('from')} ${formatDate(filters.date_from)}`
  } else if (filters.date_to) {
    return `${t('to')} ${formatDate(filters.date_to)}`
  } else {
    return t('all_records')
  }
}

function getStatusBadgeClass(status) {
  switch (status) {
    case 'completed':
      return 'badge--in'
    case 'active':
      return 'badge--warning'
    default:
      return 'badge--out'
  }
}

function getStatusText(status) {
  switch (status) {
    case 'completed':
      return t('status_completed')
    case 'active':
      return t('status_active')
    default:
      return status
  }
}

onMounted(async () => {
  await Promise.all([
    fetchEmployees(),
    fetchWorksites(),
    fetchAttendanceRecords()
  ])
})
</script>

<style scoped>
/* الإحصائيات */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 16px;
  margin-bottom: 24px;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px;
  background: var(--surface);
  border: 1px solid var(--line);
  border-radius: var(--radius-md);
  transition: transform 0.2s, box-shadow 0.2s;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

.stat-card__icon {
  width: 48px;
  height: 48px;
  border-radius: var(--radius-md);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  flex-shrink: 0;
}

.stat-card__icon :deep(svg) {
  color: inherit;
  stroke: currentColor;
}

.stat-card__icon--blue {
  background: linear-gradient(135deg, #3b82f6, #2563eb);
}

.stat-card__icon--green {
  background: linear-gradient(135deg, #10b981, #059669);
}

.stat-card__icon--purple {
  background: linear-gradient(135deg, #8b5cf6, #7c3aed);
}

.stat-card__icon--orange {
  background: linear-gradient(135deg, #f59e0b, #d97706);
}

.stat-card__content {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.stat-card__label {
  font-size: 13px;
  color: var(--ink-soft);
  font-weight: 500;
}

.stat-card__value {
  font-size: 24px;
  font-weight: 700;
  color: var(--ink);
}

/* الفلاتر والبحث */
.filters-section {
  margin-bottom: 20px;
}

.filters-header {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;
  align-items: center;
}

.search-box {
  flex: 1;
  position: relative;
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 14px;
  background: var(--surface);
  border: 1.5px solid var(--line);
  border-radius: var(--radius-sm);
  transition: border-color 0.2s, box-shadow 0.2s;
}

.search-box:focus-within {
  border-color: var(--brand);
  box-shadow: 0 0 0 3px var(--brand-tint);
}

.search-box svg, .search-box :deep(svg) {
  color: var(--ink-soft);
  flex-shrink: 0;
}

.search-input {
  flex: 1;
  border: none;
  outline: none;
  background: transparent;
  font-size: 14px;
  font-family: inherit;
  color: var(--ink);
}

.search-input::placeholder {
  color: var(--ink-soft);
}

.filters-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
}

.filter-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.filter-group label {
  font-size: 13px;
  font-weight: 600;
  color: var(--ink-soft);
}

.form-select,
.form-input {
  padding: 10px 12px;
  border: 1.5px solid var(--line);
  border-radius: var(--radius-sm);
  font-size: 14px;
  font-family: inherit;
  background: var(--surface);
}

.form-select:focus,
.form-input:focus {
  border-color: var(--brand);
  outline: none;
  box-shadow: 0 0 0 3px var(--brand-tint);
}

/* الجدول المتجاوب */
.table-responsive {
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
}

.table {
  width: 100%;
  border-collapse: collapse;
  min-width: 800px;
}

.table th {
  padding: 12px 16px;
  text-align: right;
  font-size: 13px;
  font-weight: 600;
  color: var(--ink-soft);
  border-bottom: 2px solid var(--line);
  white-space: nowrap;
}

.table th.sortable {
  cursor: pointer;
  user-select: none;
  transition: color 0.2s;
}

.table th.sortable:hover {
  color: var(--brand);
}

.sort-icon {
  margin-right: 4px;
  font-size: 12px;
}

.table td {
  padding: 12px 16px;
  border-bottom: 1px solid var(--line);
  font-size: 14px;
}

.table-row {
  transition: background-color 0.2s;
}

.table-row:hover {
  background-color: var(--surface-hover);
}

.table__person {
  display: flex;
  align-items: center;
  gap: 10px;
}

.table__avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: var(--brand);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.table__avatar :deep(svg) {
  color: white;
  stroke: white;
}

.table__name {
  font-weight: 500;
  color: var(--ink);
}

.table__worksite {
  display: inline-block;
  padding: 4px 10px;
  background: var(--surface);
  border: 1px solid var(--line);
  border-radius: var(--radius-sm);
  font-size: 13px;
  color: var(--ink-soft);
}

.hours-badge {
  display: inline-block;
  padding: 4px 10px;
  background: rgba(255, 107, 53, 0.15);
  color: #ff6b35;
  border-radius: var(--radius-sm);
  font-size: 13px;
  font-weight: 600;
}

.table__actions {
  display: flex;
  gap: 8px;
  align-items: center;
}

/* إصلاح للوضع الليلي */
[data-theme="dark"] .hours-badge {
  background: rgba(255, 107, 53, 0.25);
  color: #ff8c5a;
}

/* إصلاح النصوص في الوضع الليلي */
[data-theme="dark"] .page-head h2,
[data-theme="dark"] .page-head p {
  color: var(--ink) !important;
}

[data-theme="dark"] .stat-card__label,
[data-theme="dark"] .stat-card__value {
  color: var(--ink) !important;
}

[data-theme="dark"] .filters-header label,
[data-theme="dark"] .filter-group label {
  color: var(--ink) !important;
}

[data-theme="dark"] .search-input,
[data-theme="dark"] .form-select,
[data-theme="dark"] .form-input {
  color: var(--ink) !important;
}

[data-theme="dark"] .table th,
[data-theme="dark"] .table td {
  color: var(--ink) !important;
}

[data-theme="dark"] .pagination-info {
  color: var(--ink) !important;
}

[data-theme="dark"] .empty-state h3,
[data-theme="dark"] .empty-state p {
  color: var(--ink) !important;
}

.mono {
  font-family: 'SF Mono', 'Monaco', 'Inconsolata', 'Fira Code', monospace;
  font-size: 13px;
  color: var(--ink-soft);
}

/* الترقيم */
.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 16px;
  padding: 20px;
  border-top: 1px solid var(--line);
}

.pagination-info {
  font-size: 14px;
  color: var(--ink-soft);
}

/* التجاوب للهواتف */
@media (max-width: 768px) {
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: 12px;
  }

  .stat-card {
    padding: 16px;
    gap: 12px;
  }

  .stat-card__icon {
    width: 40px;
    height: 40px;
  }

  .stat-card__value {
    font-size: 20px;
  }

  .filters-header {
    flex-direction: column;
    align-items: stretch;
  }

  .filters-grid {
    grid-template-columns: 1fr;
  }

  .table-responsive {
    background: var(--surface);
    border-radius: var(--radius-md);
  }

  .pagination {
    flex-direction: column;
    gap: 12px;
  }
}

@media (max-width: 480px) {
  .stats-grid {
    grid-template-columns: 1fr;
  }

  .page-head {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }

  .page-head-actions {
    width: 100%;
    display: flex;
    gap: 8px;
  }

  .page-head-actions .btn {
    flex: 1;
  }
}
</style>
