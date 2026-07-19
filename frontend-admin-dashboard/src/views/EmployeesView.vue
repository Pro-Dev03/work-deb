<template>
  <div>
    <div class="page-head">
      <div>
        <h2>{{ t('employees_title') }}</h2>
        <p>{{ t('employees_description') }}</p>
      </div>
      <div class="page-head-actions">
        <button class="btn btn--primary" @click="showModal = true">+ {{ t('add_employee') }}</button>
        <button class="btn btn--danger" @click="cleanupOldRecords" :disabled="cleaning">
          {{ cleaning ? '⏳' : '🗑️' }} {{ t('cleanup_old_records') }}
        </button>
      </div>
    </div>

    <div v-if="loading" class="empty-state">
      <p>{{ t('loading_employees') }}</p>
    </div>

    <div v-else-if="employees.length === 0" class="empty-state">
      <h3>{{ t('no_employees') }}</h3>
      <p>{{ t('add_new_employee_prompt') }}</p>
    </div>

    <div v-else class="card">
      <!-- جدول للشاشات الكبيرة -->
      <div class="table-wrapper desktop-only">
        <table class="table">
          <thead>
            <tr>
              <th>{{ t('name') }}</th>
              <th>{{ t('phone') }}</th>
              <th>{{ t('role') }}</th>
              <th>{{ t('status') }}</th>
              <th>{{ t('current_worksite') }}</th>
              <th>{{ t('created_at') }}</th>
              <th>{{ t('actions') }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="emp in employees" :key="emp.id">
              <td>
                <div class="table__person">
                  <span class="table__avatar">{{ emp.full_name?.slice(0, 1) || '?' }}</span>
                  {{ emp.full_name }}
                </div>
              </td>
              <td class="mono">{{ emp.phone || '—' }}</td>
              <td>
                <span class="badge" :class="emp.role === 'admin' ? 'badge--gold' : ''">
                  {{ emp.role === 'admin' ? t('admin_role') : t('field_employee') }}
                </span>
              </td>
              <td>
                <span class="badge" :class="emp.is_active ? 'badge--in' : 'badge--out'">
                  {{ emp.is_active ? t('active_status') : t('suspended_status') }}
                </span>
              </td>
              <td>
                <span v-if="emp.current_worksite" class="badge badge--success">
                  📍 {{ emp.current_worksite }}
                </span>
                <span v-else class="text-muted">—</span>
              </td>
              <td class="mono">{{ formatDate(emp.created_at) }}</td>
              <td>
                <div class="table-actions">
                  <button 
                    class="btn btn--primary btn--sm" 
                    @click="showAttendanceHistory(emp)"
                  >
                    📊 {{ t('attendance_history') }}
                  </button>
                  <button 
                    class="btn btn--danger btn--sm" 
                    @click="confirmDelete(emp)"
                    :disabled="emp.role === 'admin' && emp.email === 'admin@worktrack.com'"
                  >
                    🗑️ {{ t('delete') }}
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      
      <!-- بطاقات للشاشات الصغيرة -->
      <div class="mobile-cards mobile-only">
        <div v-for="emp in employees" :key="emp.id" class="employee-card">
          <div class="employee-card__header">
            <div class="employee-card__person">
              <span class="table__avatar">{{ emp.full_name?.slice(0, 1) || '?' }}</span>
              <div class="employee-card__info">
                <span class="employee-card__name">{{ emp.full_name }}</span>
              </div>
            </div>
            <div class="employee-card__badges">
              <span class="badge badge--compact" :class="emp.role === 'admin' ? 'badge--gold' : ''">
                {{ emp.role === 'admin' ? t('admin_role') : t('field_employee') }}
              </span>
              <span class="badge badge--compact" :class="emp.is_active ? 'badge--in' : 'badge--out'">
                {{ emp.is_active ? t('active_status') : t('suspended_status') }}
              </span>
            </div>
          </div>
          <div class="employee-card__body">
            <div class="employee-card__row">
              <span class="employee-card__label">{{ t('phone') }}</span>
              <span class="employee-card__value mono">{{ emp.phone || '—' }}</span>
            </div>
            <div class="employee-card__row">
              <span class="employee-card__label">{{ t('current_worksite') }}</span>
              <span v-if="emp.current_worksite" class="employee-card__value">
                <span class="badge badge--success badge--compact">📍 {{ emp.current_worksite }}</span>
              </span>
              <span v-else class="employee-card__value text-muted">—</span>
            </div>
            <div class="employee-card__row">
              <span class="employee-card__label">{{ t('created_at') }}</span>
              <span class="employee-card__value mono">{{ formatDate(emp.created_at) }}</span>
            </div>
          </div>
          <div class="employee-card__actions">
            <button 
              class="btn btn--primary btn--sm btn--compact" 
              @click="showAttendanceHistory(emp)"
            >
              <span class="btn-icon">📊</span>
              <span class="btn-text">{{ t('attendance_history') }}</span>
            </button>
            <button 
              class="btn btn--danger btn--sm btn--compact" 
              @click="confirmDelete(emp)"
              :disabled="emp.role === 'admin' && emp.email === 'admin@worktrack.com'"
            >
              <span class="btn-icon">🗑️</span>
              <span class="btn-text">{{ t('delete') }}</span>
            </button>
          </div>
        </div>
      </div>
    </div>

    <EmployeeFormModal 
      v-if="showModal" 
      @close="showModal = false" 
      @employee-added="fetchEmployees"
    />

    <!-- مودال تأكيد الحذف -->
    <div v-if="showDeleteModal" class="modal-backdrop" @click.self="showDeleteModal = false">
      <div class="modal card">
        <div class="modal-header">
          <h3>⚠️ {{ t('confirm_delete_title') }}</h3>
          <button class="modal-close" @click="showDeleteModal = false">✕</button>
        </div>
        <div class="modal-body">
          <p>{{ t('confirm_delete_message') }} <strong>{{ employeeToDelete?.full_name }}</strong>؟</p>
          <p class="text-danger">{{ t('delete_irreversible') }}</p>
        </div>
        <div class="modal-footer">
          <button class="btn btn--ghost" @click="showDeleteModal = false">{{ t('cancel') }}</button>
          <button class="btn btn--danger" @click="deleteEmployee" :disabled="deleting">
            {{ deleting ? t('deleting') : t('delete_final') }}
          </button>
        </div>
      </div>
    </div>

    <!-- مودال سجل الحضور -->
    <div v-if="showAttendanceModal" class="modal-backdrop" @click.self="showAttendanceModal = false">
      <div class="modal modal-lg card">
        <div class="modal-header">
          <h3>📊 {{ t('attendance_history') }} - {{ selectedEmployee?.full_name }}</h3>
          <button class="modal-close" @click="showAttendanceModal = false">✕</button>
        </div>
        <div class="modal-body">
          <!-- فلاتر الشهر والسنة -->
          <div class="filters">
            <div class="filter-group">
              <label>{{ t('year') }}</label>
              <select v-model="selectedYear" @change="fetchAttendanceHistory" class="form-select">
                <option v-for="year in availableYears" :key="year" :value="year">{{ year }}</option>
              </select>
            </div>
            <div class="filter-group">
              <label>{{ t('month') }}</label>
              <select v-model="selectedMonth" @change="fetchAttendanceHistory" class="form-select">
                <option v-for="month in availableMonths" :key="month.value" :value="month.value">{{ month.label }}</option>
              </select>
            </div>
            <button class="btn btn--primary" @click="exportToPDF" :disabled="loadingHistory || !attendanceHistory.length">
              📄 {{ t('export_pdf') }}
            </button>
          </div>

          <!-- ملخص الشهر -->
          <div v-if="monthlySummary" class="monthly-summary">
            <div class="summary-card">
              <span class="summary-label">{{ t('total_hours') }}</span>
              <span class="summary-value">{{ monthlySummary.summary?.total_hours?.toFixed(1) || 0 }} {{ t('hours') }}</span>
            </div>
            <div class="summary-card">
              <span class="summary-label">{{ t('work_days') }}</span>
              <span class="summary-value">{{ monthlySummary.summary?.work_days || 0 }} {{ t('days') }}</span>
            </div>
          </div>

          <!-- جدول سجل الحضور -->
          <div v-if="loadingHistory" class="loading-state">
            <p>{{ t('loading') }}</p>
          </div>
          <div v-else-if="attendanceHistory.length === 0" class="empty-state">
            <p>{{ t('no_attendance_records') }}</p>
          </div>
          <div v-else>
            <!-- جدول للشاشات الكبيرة -->
            <div class="table-wrapper desktop-only">
              <table class="table">
                <thead>
                  <tr>
                    <th>{{ t('date') }}</th>
                    <th>{{ t('worksite') }}</th>
                    <th>{{ t('check_in') }}</th>
                    <th>{{ t('check_out') }}</th>
                    <th>{{ t('worked_hours') }}</th>
                    <th>{{ t('location') }}</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="record in attendanceHistory" :key="record.id">
                    <td class="mono">{{ formatDate(record.check_in_time) }}</td>
                    <td>{{ record.worksite_name || '—' }}</td>
                    <td class="mono">{{ formatTime(record.check_in_time) }}</td>
                    <td class="mono">{{ record.check_out_time ? formatTime(record.check_out_time) : '—' }}</td>
                    <td class="mono">{{ record.worked_hours ? record.worked_hours.toFixed(1) + ' ' + t('hours') : '—' }}</td>
                    <td class="mono">{{ formatDistance(record.check_in_distance_meters) }}</td>
                  </tr>
                </tbody>
              </table>
            </div>
            
            <!-- بطاقات للشاشات الصغيرة -->
            <div class="mobile-cards mobile-only">
              <div v-for="record in attendanceHistory" :key="record.id" class="attendance-card">
                <div class="attendance-card__header">
                  <span class="attendance-card__date">{{ formatDate(record.check_in_time) }}</span>
                  <span class="badge badge--info">{{ record.worked_hours ? record.worked_hours.toFixed(1) + ' ' + t('hours') : '—' }}</span>
                </div>
                <div class="attendance-card__body">
                  <div class="attendance-card__row">
                    <span class="attendance-card__label">{{ t('worksite') }}</span>
                    <span class="attendance-card__value">{{ record.worksite_name || '—' }}</span>
                  </div>
                  <div class="attendance-card__row">
                    <span class="attendance-card__label">{{ t('check_in') }}</span>
                    <span class="attendance-card__value mono">{{ formatTime(record.check_in_time) }}</span>
                  </div>
                  <div class="attendance-card__row">
                    <span class="attendance-card__label">{{ t('check_out') }}</span>
                    <span class="attendance-card__value mono">{{ record.check_out_time ? formatTime(record.check_out_time) : '—' }}</span>
                  </div>
                  <div class="attendance-card__row">
                    <span class="attendance-card__label">{{ t('location') }}</span>
                    <span class="attendance-card__value mono">{{ formatDistance(record.check_in_distance_meters) }}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import api from '../services/api'
import EmployeeFormModal from '../components/EmployeeFormModal.vue'
import { useI18n } from '../services/i18n'

const { t, currentLang } = useI18n()
const employees = ref([])
const loading = ref(false)
const showModal = ref(false)
const showDeleteModal = ref(false)
const employeeToDelete = ref(null)
const deleting = ref(false)
const cleaning = ref(false)

// سجل الحضور
const showAttendanceModal = ref(false)
const selectedEmployee = ref(null)
const attendanceHistory = ref([])
const monthlySummary = ref(null)
const loadingHistory = ref(false)
const selectedYear = ref(new Date().getFullYear())
const selectedMonth = ref(String(new Date().getMonth() + 1))

const availableYears = ref([])
const availableMonths = ref([])

// دالة لتحديث أسماء الشهور حسب اللغة
function updateMonthNames() {
  const monthKeys = ['january', 'february', 'march', 'april', 'may', 'june', 'july', 'august', 'september', 'october', 'november', 'december']
  availableMonths.value = monthKeys.map((key, index) => ({
    value: String(index + 1),
    label: t(key)
  }))
}

// توليد السنوات المتاحة
const currentYear = new Date().getFullYear()
for (let i = currentYear; i >= currentYear - 5; i--) {
  availableYears.value.push(i)
}

// تحديث أسماء الشهور عند التحميل
updateMonthNames()

// مراقبة تغيير اللغة لتحديث أسماء الشهور
watch(currentLang, () => {
  updateMonthNames()
})

async function fetchEmployees() {
  loading.value = true
  try {
    const { data } = await api.get('/admin/employees')
    employees.value = data || []
  } catch (error) {
    console.error('❌ ' + t('failed_to_fetch_employees'), error)
    employees.value = []
  } finally {
    loading.value = false
  }
}

function confirmDelete(emp) {
  employeeToDelete.value = emp
  showDeleteModal.value = true
}

async function deleteEmployee() {
  if (!employeeToDelete.value) return
  
  deleting.value = true
  try {
    await api.delete(`/admin/employees/${employeeToDelete.value.id}`)
    showDeleteModal.value = false
    await fetchEmployees()
  } catch (error) {
    console.error('❌ ' + t('failed_to_delete_employee'), error)
    alert(error.response?.data?.error || t('failed_to_delete_employee'))
  } finally {
    deleting.value = false
  }
}

function formatDate(date) {
  if (!date) return '—'
  return new Date(date).toLocaleDateString('en-GB')
}

function formatTime(date) {
  if (!date) return '—'
  return new Date(date).toLocaleTimeString('en-GB', { hour: '2-digit', minute: '2-digit', hour12: false })
}

function formatDistance(meters) {
  if (!meters) return '—'
  if (meters >= 1000) {
    return (meters / 1000).toFixed(2) + ' ' + t('kilometers')
  }
  return Math.round(meters) + ' ' + t('meters')
}

// دوال سجل الحضور
async function showAttendanceHistory(employee) {
  selectedEmployee.value = employee
  showAttendanceModal.value = true
  await fetchAttendanceHistory()
  
  // تمرير سلس للمحتوى بعد التحميل
  setTimeout(() => {
    const attendanceContent = document.querySelector('.modal-body')
    if (attendanceContent) {
      attendanceContent.scrollTo({
        top: 0,
        behavior: 'smooth'
      })
    }
  }, 100)
}

async function fetchAttendanceHistory() {
  if (!selectedEmployee.value) return
  
  loadingHistory.value = true
  try {
    const { data } = await api.get(
      `/attendance/employee/${selectedEmployee.value.id}/history?year=${selectedYear.value}&month=${selectedMonth.value}`
    )
    attendanceHistory.value = data || []
    
    // جلب الملخص الشهري
    const summaryResponse = await api.get(
      `/attendance/employee/${selectedEmployee.value.id}/monthly-summary?year=${selectedYear.value}&month=${selectedMonth.value}`
    )
    monthlySummary.value = summaryResponse.data
  } catch (error) {
    console.error('❌ ' + t('failed_to_fetch_attendance_history'), error)
    attendanceHistory.value = []
    monthlySummary.value = null
  } finally {
    loadingHistory.value = false
  }
}

async function exportToPDF() {
  if (!selectedEmployee.value || !attendanceHistory.value.length) return

  try {
    // الحصول على اللغة الحالية
    const currentLangValue = currentLang.value

    // محتوى الترجمة حسب اللغة
    const translations = {
      ar: {
        title: 'سجل الحضور',
        employeeLabel: 'الموظف',
        periodLabel: 'الفترة',
        totalHoursLabel: 'إجمالي الساعات',
        workDaysLabel: 'أيام العمل',
        dateLabel: 'التاريخ',
        worksiteLabel: 'نقطة العمل',
        checkInLabel: 'بداية العمل',
        checkOutLabel: 'نهاية العمل',
        workedHoursLabel: 'ساعات العمل',
        distanceLabel: 'المسافة',
        hoursUnit: 'ساعة',
        daysUnit: 'يوم',
        footer: 'تم إنشاء هذا التقرير من WorkTrack',
        direction: 'rtl'
      },
      en: {
        title: 'Attendance History',
        employeeLabel: 'Employee',
        periodLabel: 'Period',
        totalHoursLabel: 'Total Hours',
        workDaysLabel: 'Work Days',
        dateLabel: 'Date',
        worksiteLabel: 'Worksite',
        checkInLabel: 'Check In',
        checkOutLabel: 'Check Out',
        workedHoursLabel: 'Worked Hours',
        distanceLabel: 'Distance',
        hoursUnit: 'hours',
        daysUnit: 'days',
        footer: 'Report generated from WorkTrack',
        direction: 'ltr'
      },
      he: {
        title: 'היסטוריית נוכחות',
        employeeLabel: 'עובד',
        periodLabel: 'תקופה',
        totalHoursLabel: 'סה"כ שעות',
        workDaysLabel: 'ימי עבודה',
        dateLabel: 'תאריך',
        worksiteLabel: 'אתר עבודה',
        checkInLabel: 'כניסה',
        checkOutLabel: 'יציאה',
        workedHoursLabel: 'שעות עבודה',
        distanceLabel: 'מרחק',
        hoursUnit: 'שעות',
        daysUnit: 'ימים',
        footer: 'הדוח נוצר מ-WorkTrack',
        direction: 'rtl'
      }
    }

    const trans = translations[currentLangValue] || translations.ar

    // الحصول على اسم الشهر حسب اللغة
    const monthName = Array.isArray(availableMonths)
      ? (availableMonths.find(m => m.value === selectedMonth.value)?.label || selectedMonth.value)
      : selectedMonth.value

    // إنشاء محتوى HTML للطباعة
    const htmlContent = `
      <html dir="${trans.direction}">
      <head>
        <meta charset="UTF-8">
        <title>${trans.title} - ${selectedEmployee.value.full_name}</title>
        <style>
          body { font-family: Arial, sans-serif; padding: 20px; direction: ${trans.direction}; }
          h1 { text-align: center; color: #333; }
          .summary { display: flex; justify-content: space-around; margin: 20px 0; padding: 15px; background: #f5f5f5; border-radius: 8px; }
          .summary-item { text-align: center; }
          .summary-label { font-size: 14px; color: #666; }
          .summary-value { font-size: 24px; font-weight: bold; color: #333; }
          table { width: 100%; border-collapse: collapse; margin-top: 20px; }
          th, td { padding: 12px; border: 1px solid #ddd; }
          th { background: #4CAF50; color: white; }
          th { text-align: ${trans.direction === 'rtl' ? 'right' : 'left'}; }
          td { text-align: ${trans.direction === 'rtl' ? 'right' : 'left'}; }
          tr:nth-child(even) { background: #f9f9f9; }
          .footer { margin-top: 30px; text-align: center; color: #666; font-size: 12px; }
        </style>
      </head>
      <body>
        <h1>${trans.title}</h1>
        <h2 style="text-align: center; color: #666;">${trans.employeeLabel}: ${selectedEmployee.value.full_name}</h2>
        <p style="text-align: center; color: #666;">${trans.periodLabel}: ${monthName} ${selectedYear.value}</p>

        <div class="summary">
          <div class="summary-item">
            <div class="summary-label">${trans.totalHoursLabel}</div>
            <div class="summary-value">${monthlySummary.value?.summary?.total_hours?.toFixed(1) || 0} ${trans.hoursUnit}</div>
          </div>
          <div class="summary-item">
            <div class="summary-label">${trans.workDaysLabel}</div>
            <div class="summary-value">${monthlySummary.value?.summary?.work_days || 0} ${trans.daysUnit}</div>
          </div>
        </div>

        <table>
          <thead>
            <tr>
              <th>${trans.dateLabel}</th>
              <th>${trans.worksiteLabel}</th>
              <th>${trans.checkInLabel}</th>
              <th>${trans.checkOutLabel}</th>
              <th>${trans.workedHoursLabel}</th>
              <th>${trans.distanceLabel}</th>
            </tr>
          </thead>
          <tbody>
            ${attendanceHistory.value.map(record => `
              <tr>
                <td>${formatDate(record.check_in_time)}</td>
                <td>${record.worksite_name || '—'}</td>
                <td>${formatTime(record.check_in_time)}</td>
                <td>${record.check_out_time ? formatTime(record.check_out_time) : '—'}</td>
                <td>${record.worked_hours ? record.worked_hours.toFixed(1) + ' ' + trans.hoursUnit : '—'}</td>
                <td>${formatDistance(record.check_in_distance_meters)}</td>
              </tr>
            `).join('')}
          </tbody>
        </table>

        <div class="footer">
          <p>${trans.footer}</p>
          <p>${new Date().toLocaleDateString('en-GB')} - ${new Date().toLocaleTimeString('en-GB', { hour: '2-digit', minute: '2-digit', hour12: false })}</p>
        </div>
      </body>
      </html>
    `
    
    // إنشاء نافذة جديدة للطباعة
    const printWindow = window.open('', '_blank')
    printWindow.document.write(htmlContent)
    printWindow.document.close()
    
    // انتظار تحميل المحتوى ثم الطباعة
    printWindow.onload = function() {
      printWindow.print()
    }
  } catch (error) {
    console.error('❌ ' + t('failed_to_export_pdf'), error)
    alert(t('failed_to_export_report'))
  }
}

async function cleanupOldRecords() {
  if (!confirm(t('confirm_cleanup_old_records'))) {
    return
  }

  cleaning.value = true
  try {
    const { data } = await api.post('/attendance/cleanup-old-records')
    alert(`${t('cleanup_success')}: ${data.deleted_count} records`)
  } catch (error) {
    console.error('❌ ' + t('cleanup_failed'), error)
    alert(error.response?.data?.error || t('cleanup_failed'))
  } finally {
    cleaning.value = false
  }
}

onMounted(fetchEmployees)
</script>

<style scoped>
.page-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 18px;
  flex-wrap: wrap;
  gap: 10px;
}

.page-head h2 { font-size: 20px; margin-bottom: 4px; }
.page-head p { font-size: 13px; color: var(--ink-soft); }

.page-head-actions {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.table-wrapper { overflow-x: auto; }

/* إخفاء افتراضي للعناصر المحمولة والسطحية */
.desktop-only { display: none; }
.mobile-only { display: none; }

/* إظهار الجدول للشاشات الكبيرة فقط */
@media (min-width: 769px) {
  .desktop-only { display: block !important; }
  .mobile-only { display: none !important; }
}

/* إظهار البطاقات للشاشات الصغيرة فقط */
@media (max-width: 768px) {
  .desktop-only { display: none !important; }
  .mobile-only { display: block !important; }
}

.table { width: 100%; border-collapse: collapse; }

.table th {
  text-align: right;
  font-size: 12px;
  color: var(--ink-soft);
  font-weight: 600;
  padding: 12px 14px;
  border-bottom: 2px solid var(--line);
}

.table td {
  padding: 12px 14px;
  font-size: 14px;
  border-bottom: 1px solid var(--line);
}

.table tr:last-child td { border-bottom: none; }

.table__person { display: flex; align-items: center; gap: 10px; }

.table__avatar {
  width: 32px; height: 32px;
  border-radius: 50%;
  background: var(--brand-tint);
  color: var(--brand-dark);
  font-size: 13px;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.table-actions {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.modal-backdrop {
  position: fixed; inset: 0;
  background: rgba(0,0,0,0.6);
  backdrop-filter: blur(4px);
  display: flex; align-items: center; justify-content: center;
  z-index: 1000; padding: 20px;
}

.modal {
  width: 100%; max-width: 420px; padding: 0;
}

.modal-lg {
  max-width: 900px;
}

.modal-header {
  display: flex; justify-content: space-between; align-items: center;
  padding: 16px 20px; border-bottom: 1px solid var(--line);
}

.modal-header h3 { font-size: 17px; margin: 0; }
.modal-close { background: none; border: none; font-size: 24px; cursor: pointer; color: var(--ink-soft); }

.modal-body { 
  padding: 20px; 
  overflow-y: auto;
  max-height: 70vh;
  scroll-behavior: smooth;
}
.modal-body p { margin-bottom: 8px; }
.text-danger { color: var(--signal-out); font-weight: 600; }

.modal-footer {
  display: flex; gap: 10px;
  justify-content: flex-end;
  padding: 16px 20px;
  border-top: 1px solid var(--line);
}

/* سجل الحضور */
.filters {
  display: flex;
  gap: 16px;
  align-items: flex-end;
  margin-bottom: 20px;
  flex-wrap: wrap;
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

.form-select {
  padding: 8px 12px;
  border: 1px solid var(--line);
  border-radius: var(--radius-sm);
  font-size: 14px;
  background: var(--surface);
  color: var(--ink);
  min-width: 120px;
}

.monthly-summary {
  display: flex;
  gap: 16px;
  margin-bottom: 20px;
}

.summary-card {
  flex: 1;
  padding: 16px;
  background: var(--brand-tint);
  border-radius: var(--radius-md);
  text-align: center;
  border: 1px solid var(--brand);
}

.summary-label {
  display: block;
  font-size: 13px;
  color: var(--ink-soft);
  margin-bottom: 8px;
}

.summary-value {
  display: block;
  font-size: 24px;
  font-weight: 700;
  color: var(--brand-dark);
}

.loading-state, .empty-state {
  text-align: center;
  padding: 40px 20px;
  color: var(--ink-soft);
}

@media (max-width: 768px) {
  .page-head { flex-direction: column; }
  .page-head .btn { width: 100%; }
  
  .filters {
    flex-direction: column;
    align-items: stretch;
  }
  
  .form-select {
    width: 100%;
  }
  
  .monthly-summary {
    flex-direction: column;
  }
  
  .modal-lg {
    max-width: 100%;
  }
  
  .card {
    border-radius: var(--radius-md);
  }
  
  .mobile-cards {
    padding: 0 4px;
  }
}

/* تصميم بطاقات سجل الحضور للهاتف */
.mobile-cards {
  display: flex;
  flex-direction: column;
  gap: 12px;
  width: 100%;
  box-sizing: border-box;
}

.attendance-card {
  background: var(--surface);
  border: 1px solid var(--line);
  border-radius: var(--radius-md);
  padding: 16px;
  transition: all 0.2s ease;
  width: 100%;
  box-sizing: border-box;
}

.attendance-card:hover {
  border-color: var(--brand);
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

.attendance-card__header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  padding-bottom: 12px;
  border-bottom: 1px solid var(--line);
}

.attendance-card__date {
  font-weight: 600;
  color: var(--ink);
  font-size: 14px;
}

.attendance-card__body {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.attendance-card__row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 8px;
}

.attendance-card__label {
  font-size: 13px;
  color: var(--ink-soft);
  font-weight: 500;
  flex-shrink: 0;
}

.attendance-card__value {
  font-size: 14px;
  color: var(--ink);
  font-weight: 500;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 180px;
}

@media (max-width: 380px) {
  .attendance-card__value {
    max-width: 140px;
  }
}

/* تصميم بطاقات الموظفين للهاتف */
.employee-card {
  background: var(--surface);
  border: 1px solid var(--line);
  border-radius: var(--radius-md);
  padding: 16px;
  transition: all 0.2s ease;
  width: 100%;
  box-sizing: border-box;
}

.employee-card:hover {
  border-color: var(--brand);
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

.employee-card__header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 12px;
  padding-bottom: 12px;
  border-bottom: 1px solid var(--line);
  gap: 8px;
  flex-wrap: wrap;
}

.employee-card__person {
  display: flex;
  align-items: center;
  gap: 10px;
  flex: 1;
  min-width: 0; /* مهم للنصوص الطويلة */
}

.employee-card__info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.employee-card__name {
  font-weight: 600;
  color: var(--ink);
  font-size: 15px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 200px;
}

.employee-card__badges {
  display: flex;
  flex-direction: column;
  gap: 6px;
  align-items: flex-end;
  flex-shrink: 0;
}

.badge--compact {
  font-size: 11px;
  padding: 4px 8px;
  white-space: nowrap;
}

.badge--success {
  background-color: #10b981;
  color: white;
}

.text-muted {
  color: var(--ink-soft);
}

.employee-card__body {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 12px;
}

.employee-card__row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 8px;
}

.employee-card__label {
  font-size: 13px;
  color: var(--ink-soft);
  font-weight: 500;
  flex-shrink: 0;
}

.employee-card__value {
  font-size: 14px;
  color: var(--ink);
  font-weight: 500;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 180px;
}

.employee-card__actions {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.employee-card__actions .btn {
  flex: 1;
  min-width: 110px;
  font-size: 13px;
  padding: 8px 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
}

.btn--compact {
  font-size: 12px;
  padding: 6px 10px;
}

.btn--compact .btn-icon {
  font-size: 14px;
  flex-shrink: 0;
}

.btn--compact .btn-text {
  display: inline;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  flex: 1;
  min-width: 0;
}

@media (max-width: 380px) {
  .employee-card__name {
    max-width: 150px;
  }
  
  .employee-card__value {
    max-width: 130px;
  }
  
  .employee-card__actions .btn {
    min-width: 95px;
    font-size: 11px;
    padding: 6px 8px;
    gap: 4px;
  }
  
  .btn--compact .btn-icon {
    font-size: 12px;
  }
  
  .badge--compact {
    font-size: 10px;
    padding: 3px 6px;
  }
  
  .employee-card {
    padding: 12px;
  }
  
  .employee-card__header {
    gap: 6px;
  }
  
  .employee-card__body {
    gap: 6px;
  }
  
  .mobile-cards {
    padding: 0 2px;
  }
}

@media (max-width: 340px) {
  .employee-card__actions {
    flex-direction: column;
    gap: 6px;
  }
  
  .employee-card__actions .btn {
    min-width: 100%;
    font-size: 12px;
    padding: 8px 12px;
  }
  
  .btn--compact .btn-icon {
    font-size: 14px;
  }
  
  .employee-card {
    padding: 10px;
  }
  
  .employee-card__header {
    flex-direction: column;
    align-items: flex-start;
  }
  
  .employee-card__badges {
    align-items: flex-start;
    flex-direction: row;
    flex-wrap: wrap;
  }
  
  .employee-card__person {
    width: 100%;
  }
  
  .employee-card__name {
    max-width: 100%;
  }
  
  .employee-card__value {
    max-width: 100%;
  }
  
  .mobile-cards {
    padding: 0;
  }
  
  .card {
    border-radius: var(--radius-sm);
  }
}

/* تحسين شريط التمرير */
.modal-body::-webkit-scrollbar {
  width: 8px;
}

.modal-body::-webkit-scrollbar-track {
  background: var(--canvas);
  border-radius: 4px;
}

.modal-body::-webkit-scrollbar-thumb {
  background: var(--line);
  border-radius: 4px;
}

.modal-body::-webkit-scrollbar-thumb:hover {
  background: var(--ink-soft);
}

/* تحسين التمرير على الهاتف */
@media (max-width: 768px) {
  .modal-body {
    max-height: 60vh;
    -webkit-overflow-scrolling: touch;
  }
  
  .modal-body::-webkit-scrollbar {
    width: 4px;
  }
  
  .table-wrapper {
    overflow-x: auto;
    -webkit-overflow-scrolling: touch;
  }
  
  .employee-card {
    margin: 0 4px 8px 4px;
  }
  
  .attendance-card {
    margin: 0 4px 8px 4px;
  }
}
</style>