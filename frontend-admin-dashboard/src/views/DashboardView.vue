<template>
  <div class="dashboard">
    <div class="devpro-watermark">
      <img src="/src/assets/devpro-logo.jpg" alt="DevPro" class="watermark-logo" />
      <span class="watermark-text">{{ t('devpro_watermark') }}</span>
    </div>
    <div class="devpro-watermark">
      <img src="/src/assets/devpro-logo.jpg" alt="DevPro" class="watermark-logo" />
      <span class="watermark-text">{{ t('devpro_watermark') }}</span>
    </div>
    <div class="devpro-watermark">
      <img src="/src/assets/devpro-logo.jpg" alt="DevPro" class="watermark-logo" />
      <span class="watermark-text">{{ t('devpro_watermark') }}</span>
    </div>
    <div v-if="loading" class="loading-state">
      <div class="spinner"></div>
      <p>{{ t('loading_data') }}</p>
    </div>

    <div v-else>
      <!-- إحصائيات سريعة -->
      <div class="stats-grid">
        <div class="stat-card">
          <span class="stat-card__icon">👥</span>
          <div>
            <span class="stat-card__value">{{ stats.total_employees || 0 }}</span>
            <span class="stat-card__label">{{ t('stats_total_employees') }}</span>
          </div>
        </div>
        <div class="stat-card stat-card--waiting">
          <span class="stat-card__icon">⏳</span>
          <div>
            <span class="stat-card__value">{{ stats.waiting_employees || 0 }}</span>
            <span class="stat-card__label">{{ t('stats_waiting_employees') }}</span>
          </div>
        </div>
        <div class="stat-card stat-card--active">
          <span class="stat-card__icon">🟢</span>
          <div>
            <span class="stat-card__value">{{ activeEmployees.length }}</span>
            <span class="stat-card__label">{{ t('stats_active_now') }}</span>
          </div>
        </div>
        <div class="stat-card stat-card--completed">
          <span class="stat-card__icon">✅</span>
          <div>
            <span class="stat-card__value">{{ stats.completed_employees || 0 }}</span>
            <span class="stat-card__label">{{ t('stats_completed_today') }}</span>
          </div>
        </div>
      </div>

      <!-- الخريطة -->
      <div class="card dashboard__map">
        <div class="card-header">
          <h3>{{ t('dashboard_tracking_title') }}</h3>
          <div class="card-header__actions">
            <span class="badge badge--info">🔄 {{ updateCount }} {{ t('update_badge') }}</span>
            <button class="btn btn--sm btn--primary" @click="refreshData">{{ t('refresh_button') }}</button>
          </div>
        </div>
        <RealMap 
          :employees="activeEmployees" 
          :worksites="worksites"
          :height="500"
          @showDetails="showEmployeeDetails"
        />
      </div>

      <!-- قوائم الموظفين -->
      <div class="dashboard__tabs">
        <button 
          class="tab-btn" 
          :class="{ active: activeTab === 'active' }"
          @click="activeTab = 'active'"
        >
          🟢 {{ t('tab_active') }} ({{ activeEmployees.length }})
        </button>
        <button 
          class="tab-btn" 
          :class="{ active: activeTab === 'waiting' }"
          @click="activeTab = 'waiting'"
        >
          ⏳ {{ t('tab_waiting') }} ({{ waitingEmployees.length }})
        </button>
        <button 
          class="tab-btn" 
          :class="{ active: activeTab === 'completed' }"
          @click="activeTab = 'completed'"
        >
          ✅ {{ t('tab_completed') }} ({{ completedEmployees.length }})
        </button>
        <button 
          class="tab-btn tab-btn--alert" 
          :class="{ active: activeTab === 'alerts' }"
          @click="activeTab = 'alerts'"
        >
          🚨 {{ t('tab_alerts') }} ({{ outsideCount }})
        </button>
      </div>

      <!-- محتوى التاب -->
      <div class="dashboard__row">
        <!-- الموظفين النشطين -->
        <div v-if="activeTab === 'active'" class="card dashboard__list">
          <div class="card-header">
            <h3>{{ t('active_employees_title') }}</h3>
            <span class="badge">{{ activeEmployees.length }}</span>
          </div>
          
          <div v-if="activeEmployees.length === 0" class="empty-state">
            <p>{{ t('no_active_employees') }}</p>
          </div>
          
          <div v-else class="employee-list">
            <div 
              v-for="emp in activeEmployees" 
              :key="emp.id" 
              class="employee-item"
              :class="{ 
                'status-inside': emp.status === 'inside',
                'status-outside': emp.status === 'outside'
              }"
            >
              <div class="employee-item__avatar">
                {{ emp.full_name.slice(0, 1) }}
                <span class="status-dot" :class="emp.status"></span>
              </div>
              <div class="employee-item__info">
                <strong>{{ emp.full_name }}</strong>
                <span class="employee-item__worksite">📍 {{ emp.worksite.name }}</span>
                <span class="employee-item__time">
                  🕐 {{ formatTime(emp.check_in_time) }} | 
                  ⏱️ {{ emp.hours_worked.toFixed(1) }} ساعة
                </span>
              </div>
              <div class="employee-item__status">
                <span class="badge" :class="emp.status === 'inside' ? 'badge--in' : 'badge--out'">
                  {{ emp.status_text }}
                </span>
                <span class="employee-item__distance mono">
                  {{ formatDistance(emp.worksite.distance) }}
                </span>
              </div>
              <button 
                class="btn btn--sm btn--ghost" 
                @click="showEmployeeDetails(emp)"
                :title="t('view_details_title')"
              >
                📋
              </button>
            </div>
          </div>
        </div>

        <!-- قيد الانتظار -->
        <div v-if="activeTab === 'waiting'" class="card dashboard__list">
          <div class="card-header">
            <h3>{{ t('waiting_employees_title') }}</h3>
            <span class="badge badge--warning">{{ waitingEmployees.length }}</span>
          </div>
          
          <div v-if="waitingEmployees.length === 0" class="empty-state">
            <p>{{ t('all_employees_started') }}</p>
          </div>
          
          <div v-else class="employee-list">
            <div 
              v-for="emp in waitingEmployees" 
              :key="emp.id" 
              class="employee-item status-waiting"
            >
              <div class="employee-item__avatar" style="background: var(--signal-warning-tint); color: var(--signal-warning);">
                {{ emp.full_name.slice(0, 1) }}
                <span class="status-dot waiting"></span>
              </div>
              <div class="employee-item__info">
                <strong>{{ emp.full_name }}</strong>
                <span class="employee-item__worksite">⏳ لم يبدأ العمل بعد</span>
              </div>
              <div class="employee-item__status">
                <span class="badge badge--warning">⏳ قيد الانتظار</span>
              </div>
            </div>
          </div>
        </div>

        <!-- مكتمل -->
        <div v-if="activeTab === 'completed'" class="card dashboard__list">
          <div class="card-header">
            <h3>{{ t('completed_employees_title') }}</h3>
            <span class="badge badge--in">{{ completedEmployees.length }}</span>
          </div>
          
          <div v-if="completedEmployees.length === 0" class="empty-state">
            <p>{{ t('no_completed_employees') }}</p>
          </div>
          
          <div v-else class="employee-list">
            <div 
              v-for="emp in completedEmployees" 
              :key="emp.id" 
              class="employee-item status-completed"
            >
              <div class="employee-item__avatar" style="background: var(--signal-in-tint); color: var(--signal-in);">
                {{ emp.full_name.slice(0, 1) }}
                <span class="status-dot completed"></span>
              </div>
              <div class="employee-item__info">
                <strong>{{ emp.full_name }}</strong>
                <span class="employee-item__worksite">📍 {{ emp.worksite_name }}</span>
                <span class="employee-item__time">
                  ✅ {{ formatTime(emp.check_out_time) }} | 
                  ⏱️ {{ emp.hours_worked.toFixed(1) }} ساعة
                </span>
              </div>
              <div class="employee-item__status">
                <span class="badge badge--in">✅ مكتمل</span>
              </div>
            </div>
          </div>
        </div>

        <!-- التحذيرات الأمنية -->
        <div v-if="activeTab === 'alerts'" class="card dashboard__alerts">
          <div class="card-header">
            <h3>{{ t('alerts_title') }}</h3>
            <span class="badge badge--out">{{ outsideCount }}</span>
          </div>
         
          <div v-if="outsideCount === 0" class="empty-state">
            <p>{{ t('no_alerts') }}</p>
          </div>
          
          <div v-else class="alerts-list">
            <div 
              v-for="emp in activeEmployees.filter(e => e.status === 'outside')" 
              :key="emp.id" 
              class="alert-item"
            >
              <span class="alert-item__icon">🚨</span>
              <div class="alert-item__info">
                <strong>{{ emp.full_name }}</strong>
                <span class="alert-item__message">
                  {{ t('left_worksite_prefix') }} {{ emp.worksite.name }} ({{ formatDistance(emp.worksite.distance) }})
                </span>
                <span class="alert-item__time mono">{{ formatTime(emp.last_update) }}</span>
              </div>
              <button 
                class="btn btn--sm btn--danger" 
                @click="showEmployeeDetails(emp)"
              >
                📋
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- مودال تفاصيل الموظف -->
    <div v-if="selectedEmployee" class="modal-backdrop" @click.self="selectedEmployee = null">
      <div class="modal card">
        <div class="modal-header">
          <h3>{{ t('employee_details_title') }}</h3>
          <button class="modal-close" @click="selectedEmployee = null">✕</button>
        </div>
        <div class="modal-body">
          <div class="employee-detail">
            <div class="employee-detail__header">
              <span class="employee-detail__avatar">{{ selectedEmployee.full_name.slice(0, 1) }}</span>
              <div>
                <h4>{{ selectedEmployee.full_name }}</h4>
                <p>{{ selectedEmployee.email }} | {{ selectedEmployee.phone }}</p>
              </div>
            </div>
            
            <div class="employee-detail__info">
              <div class="info-row">
                <span class="info-label">📍 {{ t('worksite_label') }}</span>
                <span>{{ selectedEmployee.worksite?.name || t('undefined_text') }}</span>
              </div>
              <div class="info-row">
                <span class="info-label">📏 {{ t('distance_label') }}</span>
                <span :class="selectedEmployee.status === 'inside' ? 'text-success' : 'text-danger'">
                  {{ formatDistance(selectedEmployee.worksite?.distance || 0) }}
                </span>
              </div>
              <div class="info-row">
                <span class="info-label">⏱️ {{ t('working_hours_label') }}</span>
                <span>{{ (selectedEmployee.hours_worked || 0).toFixed(1) }} {{ t('hours') }}</span>
              </div>
              <div class="info-row">
                <span class="info-label">🕐 {{ t('last_update_label') }}</span>
                <span>{{ formatTime(selectedEmployee.last_update) }}</span>
              </div>
              <div class="info-row">
                <span class="info-label">📍 {{ t('location_label') }}</span>
                <span class="mono">{{ (selectedEmployee.latitude || 0).toFixed(6) }}, {{ (selectedEmployee.longitude || 0).toFixed(6) }}</span>
              </div>
            </div>

            <hr class="divider" />

            <h4>{{ t('security_notes_title') }}</h4>
            <div v-if="securityNotes.length === 0" class="empty-state">
              <p>{{ t('no_security_notes') }}</p>
            </div>
            <div v-else class="security-notes">
              <div v-for="note in securityNotes" :key="note.id" class="note-item">
                <span class="note-item__icon">⚠️</span>
                <div>
                  <p class="note-item__title">{{ note.title }}</p>
                  <p class="note-item__body">{{ note.body }}</p>
                  <span class="note-item__time mono">{{ formatTime(note.created_at) }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn--ghost" @click="selectedEmployee = null">{{ t('close') }}</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import api from '../services/api'
import RealMap from '../components/RealMap.vue'
import { useI18n } from '../services/i18n'

const { t } = useI18n()

// ==========================================
// الحالة
// ==========================================
const loading = ref(true)
const activeEmployees = ref([])
const waitingEmployees = ref([])
const completedEmployees = ref([])
const worksites = ref([])
const stats = ref({})
const selectedEmployee = ref(null)
const securityNotes = ref([])
const updateCount = ref(0)
const activeTab = ref('active')
let refreshInterval = null

// ==========================================
// العمليات الحسابية
// ==========================================
const insideCount = computed(() => {
  return activeEmployees.value.filter(e => e.status === 'inside').length
})

const outsideCount = computed(() => {
  return activeEmployees.value.filter(e => e.status === 'outside').length
})

// ==========================================
// دوال مساعدة
// ==========================================
function formatTime(date) {
  if (!date) return '—'
  return new Date(date).toLocaleTimeString('ar-SA', { 
    hour: '2-digit', 
    minute: '2-digit',
    second: '2-digit'
  })
}

function formatDistance(meters) {
  if (!meters) return '0 م'
  if (meters >= 1000) {
    return (meters / 1000).toFixed(2) + ' كيلومتر'
  }
  return Math.round(meters) + ' متر'
}

// ==========================================
// جلب البيانات
// ==========================================
async function fetchData() {
  try {
    const [employeesRes, worksitesRes, statsRes, waitingRes, completedRes] = await Promise.all([
      api.get('/location/active'),
      api.get('/worksites'),
      api.get('/reports/daily-summary'),
      api.get('/reports/pending-employees'),
      api.get('/reports/completed-employees')
    ])
    
    activeEmployees.value = employeesRes.data || []
    worksites.value = worksitesRes.data || []
    stats.value = statsRes.data || {}
    waitingEmployees.value = waitingRes.data || []
    completedEmployees.value = completedRes.data || []
    updateCount.value++
    
  } catch (error) {
    console.error('❌ فشل جلب البيانات:', error)
  } finally {
    loading.value = false
  }
}

async function refreshData() {
  loading.value = true
  await fetchData()
  loading.value = false
}

// ==========================================
// عرض تفاصيل الموظف
// ==========================================
async function showEmployeeDetails(employee) {
  selectedEmployee.value = employee
  
  // جلب الملاحظات الأمنية
  try {
    const { data } = await api.get(`/location/security/${employee.id}`)
    securityNotes.value = data || []
  } catch (error) {
    console.error('❌ فشل جلب الملاحظات الأمنية:', error)
    securityNotes.value = []
  }
}

// ==========================================
// دورة الحياة
// ==========================================
onMounted(async () => {
  await fetchData()
  
  // تحديث كل 10 ثواني
  refreshInterval = setInterval(fetchData, 10000)
})

onUnmounted(() => {
  if (refreshInterval) {
    clearInterval(refreshInterval)
  }
})
</script>

<style scoped>
.dashboard { display: flex; flex-direction: column; gap: 22px; }

.loading-state {
  text-align: center; padding: 60px 20px;
  display: flex; flex-direction: column; align-items: center; gap: 16px;
}

.spinner {
  width: 40px; height: 40px;
  border: 4px solid var(--line);
  border-top-color: var(--brand);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin { to { transform: rotate(360deg); } }

/* ==========================================
   إحصائيات
   ========================================== */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
}

.stat-card {
  background: var(--surface);
  border: 1px solid var(--line);
  border-radius: var(--radius-lg);
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 14px;
  box-shadow: var(--shadow-sm);
  transition: all var(--transition-base);
}

.stat-card:hover {
  box-shadow: var(--shadow-md);
  transform: translateY(-2px);
}

.stat-card--waiting { border-left: 4px solid var(--signal-warning); }
.stat-card--active { border-left: 4px solid var(--signal-in); }
.stat-card--completed { border-left: 4px solid #22C55E; }

.stat-card__icon { font-size: 28px; }
.stat-card__value { font-size: 28px; font-weight: 700; display: block; }
.stat-card__label { font-size: 13px; color: var(--ink-soft); }

/* ==========================================
   الخريطة
   ========================================== */
.dashboard__map { padding: 20px; }

.card-header {
  display: flex; align-items: center; justify-content: space-between;
  margin-bottom: 14px;
  flex-wrap: wrap;
  gap: 10px;
}

.card-header h3 { font-size: 16px; }
.card-header__actions { display: flex; gap: 8px; align-items: center; }

/* ==========================================
   Tabs
   ========================================== */
.dashboard__tabs {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.tab-btn {
  padding: 10px 20px;
  border: 2px solid var(--line);
  border-radius: var(--radius-sm);
  background: var(--surface);
  color: var(--ink-soft);
  font-weight: 600;
  cursor: pointer;
  transition: all var(--transition-fast);
  font-family: var(--font-body);
  font-size: 14px;
}

.tab-btn:hover {
  border-color: var(--brand);
  color: var(--brand);
}

.tab-btn.active {
  border-color: var(--brand);
  background: var(--brand-tint);
  color: var(--brand);
}

.tab-btn--alert {
  border-color: var(--signal-out);
  color: var(--signal-out);
}

.tab-btn--alert:hover {
  background: var(--signal-out-tint);
  border-color: var(--signal-out);
}

.tab-btn--alert.active {
  background: var(--signal-out-tint);
  border-color: var(--signal-out);
  color: var(--signal-out);
}

/* ==========================================
   قوائم الموظفين
   ========================================== */
.dashboard__row {
  display: grid;
  grid-template-columns: 1fr;
  gap: 16px;
}

.dashboard__list { padding: 20px; }

.employee-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  max-height: 500px;
  overflow-y: auto;
}

.employee-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 14px;
  border-radius: var(--radius-sm);
  border: 1px solid var(--line);
  transition: all 0.2s;
  background: var(--surface);
}

.employee-item:hover {
  box-shadow: var(--shadow-sm);
}

.employee-item.status-inside { border-right: 4px solid var(--signal-in); }
.employee-item.status-outside { 
  border-right: 4px solid var(--signal-out);
  background: var(--signal-out-tint);
}
.employee-item.status-waiting {
  border-right: 4px solid var(--signal-warning);
  background: var(--signal-warning-tint);
}
.employee-item.status-completed {
  border-right: 4px solid #22C55E;
  background: #22C55E10;
}

.employee-item__avatar {
  position: relative;
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: var(--brand-tint);
  color: var(--brand);
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  font-size: 16px;
  flex-shrink: 0;
}

.status-dot {
  position: absolute;
  bottom: -2px;
  right: -2px;
  width: 14px;
  height: 14px;
  border-radius: 50%;
  border: 2px solid var(--surface);
}

.status-dot.inside { background: var(--signal-in); }
.status-dot.outside { background: var(--signal-out); }
.status-dot.waiting { background: var(--signal-warning); }
.status-dot.completed { background: #22C55E; }

.employee-item__info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}

.employee-item__info strong {
  font-size: 14px;
  color: var(--ink);
}

.employee-item__worksite { font-size: 12px; color: var(--ink-soft); }
.employee-item__time { font-size: 11px; color: var(--ink-light); }

.employee-item__status {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 4px;
  flex-shrink: 0;
}

.employee-item__distance { font-size: 11px; color: var(--ink-soft); }

/* ==========================================
   التحذيرات
   ========================================== */
.dashboard__alerts { padding: 20px; }

.alerts-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  max-height: 500px;
  overflow-y: auto;
}

.alert-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 14px;
  background: var(--signal-out-tint);
  border-radius: var(--radius-sm);
  border: 1px solid var(--signal-out);
}

.alert-item__icon { font-size: 20px; }

.alert-item__info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.alert-item__info strong {
  font-size: 13px;
  color: var(--signal-out);
}

.alert-item__message { font-size: 12px; color: var(--ink-soft); }
.alert-item__time { font-size: 11px; color: var(--ink-light); }

/* ==========================================
   Badges
   ========================================== */
.badge {
  display: inline-block;
  padding: 2px 10px;
  border-radius: 999px;
  font-size: 11px;
  font-weight: 600;
}

.badge--in { background: #22C55E20; color: #22C55E; }
.badge--out { background: #EF444420; color: #EF4444; }
.badge--warning { background: #F59E0B20; color: #F59E0B; }
.badge--info { background: #3B82F620; color: #3B82F6; }

/* ==========================================
   مودال
   ========================================== */
.modal-backdrop {
  position: fixed; inset: 0;
  background: rgba(0,0,0,0.6);
  backdrop-filter: blur(4px);
  display: flex; align-items: center; justify-content: center;
  z-index: 1000; padding: 20px;
}

.modal {
  width: 100%; max-width: 500px;
  max-height: 90vh; overflow-y: auto;
  padding: 0;
}

.modal-header {
  display: flex; justify-content: space-between; align-items: center;
  padding: 16px 20px; border-bottom: 1px solid var(--line);
}

.modal-header h3 { font-size: 17px; margin: 0; }
.modal-close { background: none; border: none; font-size: 24px; cursor: pointer; color: var(--ink-soft); }

.modal-body { padding: 20px; }
.modal-footer {
  padding: 16px 20px; border-top: 1px solid var(--line);
  display: flex; gap: 10px; justify-content: flex-end;
}

.employee-detail__header {
  display: flex; align-items: center; gap: 14px;
  margin-bottom: 16px;
}

.employee-detail__avatar {
  width: 48px; height: 48px;
  border-radius: 50%;
  background: var(--brand-tint);
  color: var(--brand);
  display: flex; align-items: center; justify-content: center;
  font-weight: 700; font-size: 20px;
}

.employee-detail__info {
  display: flex; flex-direction: column; gap: 6px;
}

.info-row {
  display: flex; gap: 8px;
  font-size: 13px;
}

.info-label {
  font-weight: 600;
  color: var(--ink-soft);
  min-width: 100px;
}

.text-success { color: var(--signal-in); font-weight: 600; }
.text-danger { color: var(--signal-out); font-weight: 600; }

.divider { border: none; border-top: 1px solid var(--line); margin: 16px 0; }

.security-notes {
  display: flex; flex-direction: column; gap: 8px;
}

.note-item {
  display: flex; gap: 10px;
  padding: 10px 12px;
  background: var(--signal-warning-tint);
  border-radius: var(--radius-sm);
  border: 1px solid var(--signal-warning);
}

.note-item__icon { font-size: 18px; }
.note-item__title { font-size: 13px; font-weight: 600; color: var(--signal-warning); }
.note-item__body { font-size: 12px; color: var(--ink-soft); margin: 2px 0; }
.note-item__time { font-size: 11px; color: var(--ink-light); }

.empty-state { text-align: center; padding: 30px 20px; color: var(--ink-soft); }

/* ==========================================
   استجابة
   ========================================== */
@media (max-width: 960px) {
  .stats-grid { grid-template-columns: repeat(2, 1fr); }
}

@media (max-width: 600px) {
  .stats-grid { grid-template-columns: 1fr 1fr; }
  .employee-item { flex-wrap: wrap; }
  .employee-item__status { flex-direction: row; align-items: center; gap: 8px; }
  .dashboard__tabs { flex-direction: column; }
  .tab-btn { width: 100%; text-align: center; }
}
</style>

<style>
/* ==========================================
   خلفية DevPro في Dashboard
   ========================================== */
.dashboard {
  position: relative;
  overflow: hidden;
}

.devpro-watermark {
  position: fixed;
  bottom: 20px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 20px;
  background: rgba(255, 255, 255, 0.85);
  backdrop-filter: blur(10px);
  border-radius: 999px;
  border: 1px solid rgba(30, 58, 95, 0.1);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.06);
  z-index: 5;
  pointer-events: none;
  opacity: 0.6;
  transition: opacity 0.3s ease;
}

.devpro-watermark:hover {
  opacity: 1;
}

.watermark-logo {
  width: 24px;
  height: 24px;
  border-radius: 6px;
  object-fit: contain;
  background: white;
  padding: 2px;
}

.watermark-text {
  font-size: 11px;
  color: #1E3A5F;
  font-weight: 500;
  white-space: nowrap;
}

@media (max-width: 768px) {
  .devpro-watermark {
    padding: 6px 14px;
    bottom: 10px;
  }
  
  .watermark-logo {
    width: 18px;
    height: 18px;
  }
  
  .watermark-text {
    font-size: 9px;
  }
}
</style>

<style>
/* ==========================================
   خلفية DevPro في Dashboard
   ========================================== */
.dashboard {
  position: relative;
  overflow: hidden;
}

.devpro-watermark {
  position: fixed;
  bottom: 20px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 20px;
  background: rgba(255, 255, 255, 0.85);
  backdrop-filter: blur(10px);
  border-radius: 999px;
  border: 1px solid rgba(30, 58, 95, 0.1);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.06);
  z-index: 5;
  pointer-events: none;
  opacity: 0.6;
  transition: opacity 0.3s ease;
}

.devpro-watermark:hover {
  opacity: 1;
}

.watermark-logo {
  width: 24px;
  height: 24px;
  border-radius: 6px;
  object-fit: contain;
  background: white;
  padding: 2px;
}

.watermark-text {
  font-size: 11px;
  color: #1E3A5F;
  font-weight: 500;
  white-space: nowrap;
}

@media (max-width: 768px) {
  .devpro-watermark {
    padding: 6px 14px;
    bottom: 10px;
  }
  
  .watermark-logo {
    width: 18px;
    height: 18px;
  }
  
  .watermark-text {
    font-size: 9px;
  }
}
</style>

<style>
/* ==========================================
   خلفية DevPro في Dashboard
   ========================================== */
.dashboard {
  position: relative;
  overflow: hidden;
}

.devpro-watermark {
  position: fixed;
  bottom: 20px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 20px;
  background: rgba(255, 255, 255, 0.85);
  backdrop-filter: blur(10px);
  border-radius: 999px;
  border: 1px solid rgba(30, 58, 95, 0.1);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.06);
  z-index: 5;
  pointer-events: none;
  opacity: 0.6;
  transition: opacity 0.3s ease;
}

.devpro-watermark:hover {
  opacity: 1;
}

.watermark-logo {
  width: 24px;
  height: 24px;
  border-radius: 6px;
  object-fit: contain;
  background: white;
  padding: 2px;
}

.watermark-text {
  font-size: 11px;
  color: #1E3A5F;
  font-weight: 500;
  white-space: nowrap;
}

@media (max-width: 768px) {
  .devpro-watermark {
    padding: 6px 14px;
    bottom: 10px;
  }
  
  .watermark-logo {
    width: 18px;
    height: 18px;
  }
  
  .watermark-text {
    font-size: 9px;
  }
}
</style>
