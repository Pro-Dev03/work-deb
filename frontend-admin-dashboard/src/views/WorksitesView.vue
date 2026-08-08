<template>
  <div>
    <div class="page-head">
      <div>
        <h2>{{ t('worksites_title') }}</h2>
        <p>{{ t('worksites_description') }}</p>
      </div>
      <button class="btn btn--primary" @click="openAddModal">+ {{ t('new_worksite') }}</button>
    </div>

    <div v-if="loading" class="empty-state"><p>{{ t('loading_worksites') }}</p></div>

    <div v-else-if="worksites.length === 0" class="empty-state">
      <h3>{{ t('no_worksites') }}</h3>
      <p>{{ t('add_worksite_prompt') }}</p>
    </div>

    <div v-else class="sites-grid">
      <div v-for="site in worksites" :key="site.id" class="card site-card" :class="{ 'site-card--unassigned': site.is_unassigned }">
        <div class="site-card__header">
          <h3>{{ site.name }}</h3>
          <div class="site-card__actions" v-if="!site.is_unassigned">
            <div class="dropdown" :class="{ 'dropdown--active': activeDropdown === site.id }">
              <button 
                class="dropdown-toggle"
                @click="toggleDropdown(site.id)"
              >
                ⚙️
              </button>
              <div class="dropdown-menu">
                <button 
                  class="dropdown-item"
                  @click.stop="openAssignModal(site)"
                >
                  <User :size="16" />
                  <span>{{ t('assign_employee') }}</span>
                </button>
                <button 
                  class="dropdown-item"
                  @click.stop="openEditModal(site)"
                >
                  <Edit :size="16" />
                  <span>{{ t('edit') }}</span>
                </button>
                <div class="dropdown-divider"></div>
                <button 
                  class="dropdown-item dropdown-item--danger"
                  @click.stop="confirmDelete(site)"
                >
                  <Trash2 :size="16" />
                  <span>{{ t('delete') }}</span>
                </button>
              </div>
            </div>
          </div>
        </div>
        <p class="site-card__address">{{ site.address || t('no_address') }}</p>
        <div class="site-card__details" v-if="!site.is_unassigned">
          <span class="mono">📍 {{ site.latitude?.toFixed(5) }}, {{ site.longitude?.toFixed(5) }}</span>
          <span class="site-card__radius mono">⭕ {{ site.radius_meters || 100 }} {{ t('meters_unit') }}</span>
        </div>
        
        <!-- عرض الموظف المعين -->
        <div class="site-card__assigned" v-if="!site.is_unassigned">
          <span v-if="site.assigned_to?.name" class="badge badge--in">
            👤 {{ site.assigned_to.name }}
          </span>
          <span v-else class="badge badge--out">
            ⚠️ {{ t('unassigned') }}
          </span>
        </div>
        
        <!-- عرض الموظفين العاملين حالياً -->
        <div class="site-card__working" v-if="site.working_employees && site.working_employees.length > 0">
          <div class="site-card__working-label">
            {{ t('currently_working') }} ({{ site.working_employees.length }})
          </div>
          <div class="site-card__working-list">
            <div v-for="emp in site.working_employees" :key="emp.id" class="working-employee-item">
              <span class="badge badge--success badge--compact">
                👤 {{ emp.name }}
              </span>
              <button 
                class="btn btn--warning btn--xs" 
                @click="confirmForceCheckout(emp)"
                :disabled="forceCheckingOut"
              >
                ⏱️ {{ t('end_shift') }}
              </button>
            </div>
          </div>
        </div>
        
        <span class="badge" :class="site.is_active ? 'badge--in' : 'badge--out'" v-if="!site.is_unassigned">
          {{ site.is_active ? t('active_status') : t('inactive_status') }}
        </span>
      </div>
    </div>

    <WorksiteFormModal 
      v-if="showModal" 
      :worksite="worksiteToEdit" 
      @close="showModal = false" 
      @worksite-added="fetchWorksites" 
      @worksite-updated="fetchWorksites" 
    />

    <!-- مودال تأكيد الحذف -->
    <div v-if="showDeleteModal" class="modal-backdrop" @click.self="showDeleteModal = false">
      <div class="modal card">
        <div class="modal-header">
          <h3>⚠️ {{ t('confirm_delete_title') }}</h3>
          <button class="modal-close" @click="showDeleteModal = false">✕</button>
        </div>
        <div class="modal-body">
          <p>{{ t('confirm_delete_message') }} <strong>{{ siteToDelete?.name }}</strong>؟</p>
          <p class="text-danger">{{ t('delete_warning') }}</p>
        </div>
        <div class="modal-footer">
          <button class="btn btn--ghost" @click="showDeleteModal = false">{{ t('cancel') }}</button>
          <button class="btn btn--danger" @click="deleteWorksite" :disabled="deleting">
            {{ deleting ? t('deleting') : t('delete_final') }}
          </button>
        </div>
      </div>
    </div>

    <!-- مودال تعيين موظف - محسّن مع عرض سريع -->
    <div v-if="showAssignModal" class="modal-backdrop" @click.self="showAssignModal = false">
      <div class="modal card assign-modal">
        <div class="modal-header">
          <h3>👤 {{ t('assign_employee_title') }}</h3>
          <button class="modal-close" @click="showAssignModal = false">✕</button>
        </div>
        <div class="modal-body">
          <p>{{ t('choose_employee_to_assign') }} <strong>{{ worksiteToAssign?.name }}</strong></p>
          
          <div v-if="loadingEmployees" class="empty-state">
            <p>{{ t('loading_employees') }}</p>
          </div>
          
          <div v-else-if="employees.length === 0" class="empty-state">
            <p>{{ t('no_available_employees') }}</p>
          </div>
          
          <div v-else class="employees-list">
            <button
              v-for="emp in employees"
              :key="emp.id"
              class="employee-item"
              @click="assignEmployee(emp.id)"
              :disabled="assigning"
            >
              <span class="employee-item__avatar">
                <User :size="16" />
              </span>
              <div>
                <strong>{{ emp.full_name }}</strong>
              </div>
              <span class="employee-item__check">✓</span>
            </button>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn--ghost" @click="showAssignModal = false">{{ t('cancel') }}</button>
        </div>
      </div>
    </div>

    <!-- مودال تأكيد إنهاء الدوام -->
    <div v-if="showForceCheckoutModal" class="modal-backdrop" @click.self="showForceCheckoutModal = false">
      <div class="modal card">
        <div class="modal-header">
          <h3>⏱️ {{ t('force_checkout_title') }}</h3>
          <button class="modal-close" @click="showForceCheckoutModal = false">✕</button>
        </div>
        <div class="modal-body">
          <p>{{ t('force_checkout_message') }} <strong>{{ employeeToForceCheckout?.name }}</strong>؟</p>
          <p class="text-warning">{{ t('force_checkout_warning') }}</p>
        </div>
        <div class="modal-footer">
          <button class="btn btn--ghost" @click="showForceCheckoutModal = false">{{ t('cancel') }}</button>
          <button class="btn btn--warning" @click="forceCheckoutEmployee" :disabled="forceCheckingOut">
            {{ forceCheckingOut ? t('processing') : t('confirm_end_shift') }}
          </button>
        </div>
      </div>
    </div>

    <!-- Success Modal -->
    <div v-if="showSuccessModal" class="modal-backdrop" @click.self="showSuccessModal = false">
      <div class="success-modal-card">
        <div class="success-icon-wrapper">
          <CheckCircle :size="32" />
        </div>
        <h3 class="success-title">{{ t('success') }}</h3>
        <p class="success-message">{{ successMessage }}</p>
        <button class="btn btn--primary success-confirm-btn" @click="showSuccessModal = false">
          {{ t('ok') }}
        </button>
      </div>
    </div>

    <!-- بطاقة إنهاء الدوام الاحترافية -->
    <div v-if="showCompletionCard" class="modal-backdrop" @click.self="showCompletionCard = false">
      <div class="completion-card">
        <div class="completion-header">
          <div class="completion-icon">
            <CheckCircle :size="48" />
          </div>
          <h2 class="completion-title">{{ t('shift_completed') }}</h2>
          <p class="completion-subtitle">{{ t('shift_completed_message') }}</p>
        </div>
        
        <div class="completion-body">
          <div class="completion-info-row">
            <div class="completion-info-item">
              <div class="completion-info-label">{{ t('employee_name') }}</div>
              <div class="completion-info-value">{{ completionData?.employeeName }}</div>
            </div>
            <div class="completion-info-item">
              <div class="completion-info-label">{{ t('worksite') }}</div>
              <div class="completion-info-value">{{ completionData?.worksiteName }}</div>
            </div>
          </div>
          
          <div class="completion-info-row">
            <div class="completion-info-item">
              <div class="completion-info-label">{{ t('check_in_time') }}</div>
              <div class="completion-info-value mono">{{ formatTime(completionData?.checkInTime) }}</div>
            </div>
            <div class="completion-info-item">
              <div class="completion-info-label">{{ t('check_out_time') }}</div>
              <div class="completion-info-value mono">{{ formatTime(completionData?.checkOutTime) }}</div>
            </div>
          </div>
          
          <div class="completion-hours-section">
            <div class="completion-hours-label">{{ t('total_worked_hours') }}</div>
            <div class="completion-hours-value">{{ completionData?.workedHours?.toFixed(1) || '0.0' }} {{ t('hours') }}</div>
          </div>
        </div>
        
        <div class="completion-footer">
          <button class="btn btn--primary completion-confirm-btn" @click="showCompletionCard = false">
            {{ t('close') }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import api, { clearCacheForEndpoint } from '../services/api'
import WorksiteFormModal from '../components/WorksiteFormModal.vue'
import { useI18n } from '../services/i18n'
import { User, CheckCircle, Edit, Trash2 } from '@lucide/vue'

const { t, currentLang } = useI18n()
const worksites = ref([])
const employees = ref([])
const loading = ref(false)
const loadingEmployees = ref(false)
const deleting = ref(false)
const assigning = ref(false)
const forceCheckingOut = ref(false)
const showModal = ref(false)
const showDeleteModal = ref(false)
const showAssignModal = ref(false)
const showForceCheckoutModal = ref(false)
const showSuccessModal = ref(false)
const showCompletionCard = ref(false)
const successMessage = ref('')
const completionData = ref(null)
const siteToDelete = ref(null)
const worksiteToAssign = ref(null)
const employeeToForceCheckout = ref(null)
const worksiteToEdit = ref(null)
const activeDropdown = ref(null)
const employeesCache = ref(null) // تخزين مؤقت للموظفين

async function fetchWorksites() {
  loading.value = true
  try {
    const { data } = await api.get('/worksites')
    // التعامل مع response format الجديد مع pagination
    worksites.value = data.data || data || []
  } catch (error) {
    console.error('❌ ' + t('failed_to_fetch_worksites'), error)
    worksites.value = []
  } finally {
    loading.value = false
  }
}

function confirmDelete(site) {
  siteToDelete.value = site
  showDeleteModal.value = true
}

async function deleteWorksite() {
  if (!siteToDelete.value) return
  
  deleting.value = true
  try {
    await api.delete(`/worksites/${siteToDelete.value.id}`)
    showDeleteModal.value = false
    await fetchWorksites()
  } catch (error) {
    console.error('❌ ' + t('failed_to_delete_worksite'), error)
    alert(error.response?.data?.error || t('failed_to_delete_worksite'))
  } finally {
    deleting.value = false
  }
}

function openAssignModal(site) {
  worksiteToAssign.value = site
  showAssignModal.value = true
  fetchEmployees()
}

function openEditModal(site) {
  worksiteToEdit.value = site
  showModal.value = true
}

function openAddModal() {
  worksiteToEdit.value = null
  showModal.value = true
}

async function fetchEmployees() {
  // استخدام الـ cache إذا كان موجوداً
  if (employeesCache.value && !loadingEmployees.value) {
    employees.value = employeesCache.value
    return
  }

  loadingEmployees.value = true
  try {
    const { data } = await api.get('/worksites/employees')
    employees.value = data || []
    employeesCache.value = data || [] // تخزين في الـ cache
  } catch (error) {
    console.error('❌ ' + t('failed_to_fetch_employees'), error)
    employees.value = []
  } finally {
    loadingEmployees.value = false
  }
}

async function assignEmployee(employeeId) {
  if (!worksiteToAssign.value) return
  
  assigning.value = true
  try {
    const response = await api.post('/worksites/assign', {
      employee_id: employeeId,
      worksite_id: worksiteToAssign.value.id
    })
    
    showAssignModal.value = false
    showSuccessModal.value = true
    successMessage.value = response.data.message || t('employee_assigned_successfully')
    
    // تحديث البيانات
    employeesCache.value = null // مسح الـ cache المحلي للموظفين
    await fetchWorksites()
    await fetchEmployees()
  } catch (error) {
    console.error('❌ ' + t('failed_to_assign_employee'), error)
    const msg = error.response?.data?.error || t('failed_to_assign_employee')
    alert('❌ ' + msg)
  } finally {
    assigning.value = false
  }
}

function confirmForceCheckout(employee) {
  employeeToForceCheckout.value = employee
  showForceCheckoutModal.value = true
}

async function forceCheckoutEmployee() {
  if (!employeeToForceCheckout.value) return
  
  forceCheckingOut.value = true
  try {
    const response = await api.post('/attendance/force-checkout', {
      attendance_id: employeeToForceCheckout.value.attendance_id
    })
    
    showForceCheckoutModal.value = false
    
    // عرض البطاقة الاحترافية مع معلومات الموظف من الاستجابة
    showCompletionCard.value = true
    completionData.value = {
      employeeName: response.data.employee_name || employeeToForceCheckout.value.name,
      worksiteName: response.data.worksite_name || '—',
      checkInTime: response.data.check_in_time,
      checkOutTime: response.data.check_out_time,
      workedHours: response.data.worked_hours,
      message: response.data.message || t('force_checkout_success')
    }
    
    await fetchWorksites()
  } catch (error) {
    console.error('❌ ' + t('force_checkout_failed'), error)
    const msg = error.response?.data?.error || t('force_checkout_failed')
    alert('❌ ' + msg)
  } finally {
    forceCheckingOut.value = false
  }
}

function formatTime(dateStr) {
  if (!dateStr) return '—'
  const date = new Date(dateStr)
  return date.toLocaleTimeString(currentLang.value === 'ar' ? 'ar-SA' : 'en-US', { 
    hour: '2-digit', 
    minute: '2-digit' 
  })
}

function toggleDropdown(siteId) {
  if (activeDropdown.value === siteId) {
    activeDropdown.value = null
  } else {
    activeDropdown.value = siteId
  }
}

function closeDropdown() {
  activeDropdown.value = null
}

function handleClickOutside(event) {
  // Don't close if clicking inside the dropdown menu or toggle
  if (event.target.closest('.dropdown-menu') || event.target.closest('.dropdown-toggle')) {
    return
  }

  // Close if clicking outside
  closeDropdown()
}

onMounted(() => {
  fetchWorksites()
  // تحميل الموظفين مسبقاً لتحسين الأداء عند فتح مودال التعيين
  fetchEmployees()
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
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

.sites-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
}

.site-card {
  padding: 18px 20px;
  transition: all var(--transition-base);
  overflow: visible;
}

.site-card:hover {
  z-index: 10;
}

.site-card__header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 8px;
}

.site-card__header h3 { font-size: 16px; }

.site-card__actions {
  display: flex;
  gap: 6px;
  position: relative;
}

/* =============================================
   Dropdown Menu Styles
   ============================================= */
.dropdown {
  position: relative;
}

.dropdown-toggle {
  width: 40px;
  height: 36px;
  border-radius: var(--radius-sm);
  border: 1px solid var(--line);
  background: var(--surface);
  color: var(--ink-soft);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all var(--transition-base);
  padding: 0;
  font-size: 18px;
}

.dropdown-toggle:hover {
  background: var(--brand-tint);
  color: var(--brand);
  border-color: var(--brand);
}

.dropdown-menu {
  position: static;
  min-width: 200px;
  background: var(--surface);
  border: 1px solid var(--line);
  border-radius: var(--radius-md);
  box-shadow: var(--shadow-xl);
  padding: 8px 0;
  z-index: 99999;
  opacity: 0;
  visibility: hidden;
  height: 0;
  overflow: hidden;
  transition: all var(--transition-base);
}

.dropdown--active .dropdown-menu {
  opacity: 1;
  visibility: visible;
  height: auto;
  overflow: visible;
}

.dropdown-item {
  display: flex;
  align-items: center;
  gap: 10px;
  width: 100%;
  padding: 10px 16px;
  border: none;
  background: transparent;
  color: var(--ink);
  font-size: 14px;
  font-weight: 500;
  text-align: right;
  cursor: pointer;
  transition: all var(--transition-fast);
}

.dropdown-item:hover {
  background: var(--brand-tint);
  color: var(--brand);
}

.dropdown-item:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.dropdown-item :deep(svg) {
  color: inherit;
  stroke: inherit;
  flex-shrink: 0;
}

.dropdown-divider {
  height: 1px;
  background: var(--line);
  margin: 8px 0;
}

.dropdown-item--danger {
  color: var(--signal-out);
}

.dropdown-item--danger:hover {
  background: var(--signal-out-tint);
  color: var(--signal-out);
}

[data-theme="dark"] .dropdown-toggle {
  background: var(--surface);
  border-color: var(--line);
  color: var(--ink-soft);
}

[data-theme="dark"] .dropdown-toggle:hover {
  background: rgba(212, 175, 55, 0.15);
  color: var(--gold);
  border-color: var(--gold);
}

[data-theme="dark"] .dropdown-menu {
  background: var(--surface);
  border-color: var(--line);
}

[data-theme="dark"] .dropdown-item {
  color: var(--ink);
}

[data-theme="dark"] .dropdown-item:hover {
  background: rgba(212, 175, 55, 0.15);
  color: var(--gold);
}

[data-theme="dark"] .dropdown-item--danger {
  color: var(--signal-out);
}

[data-theme="dark"] .dropdown-item--danger:hover {
  background: rgba(239, 68, 68, 0.15);
  color: var(--signal-out);
}

.site-card {
  overflow: visible;
}

.sites-grid {
  overflow: visible;
}

.site-card__address {
  font-size: 13px;
  color: var(--ink-soft);
  margin-bottom: 8px;
}

.site-card__details {
  display: flex;
  gap: 12px;
  font-size: 12px;
  color: var(--ink-soft);
  margin-bottom: 10px;
  flex-wrap: wrap;
}

.site-card__radius {
  background: var(--brand-tint);
  color: var(--brand);
  padding: 2px 10px;
  border-radius: 999px;
}

.site-card__assigned { margin-bottom: 8px; }

.site-card__working {
  margin-top: 12px;
  padding: 10px;
  background: var(--brand-tint);
  border-radius: var(--radius-sm);
  border: 1px solid var(--brand);
}

.site-card__working-label {
  font-size: 12px;
  font-weight: 600;
  color: var(--brand-dark);
  margin-bottom: 6px;
}

.site-card__working-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.working-employee-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.site-card--unassigned {
  border: 2px solid var(--signal-out);
  background: linear-gradient(135deg, var(--surface) 0%, rgba(239, 68, 68, 0.1) 100%);
}

.site-card--unassigned .site-card__header h3 {
  color: var(--signal-out);
}

.modal-backdrop {
  position: fixed; inset: 0;
  background: rgba(0,0,0,0.6);
  backdrop-filter: blur(4px);
  display: flex; align-items: center; justify-content: center;
  z-index: 1000; padding: 20px;
}

.modal {
  width: 100%; max-width: 460px; padding: 0;
  max-height: 90vh;
  overflow-y: auto;
}

.assign-modal { max-width: 500px; }

.modal-header {
  display: flex; justify-content: space-between; align-items: center;
  padding: 16px 20px; border-bottom: 1px solid var(--line);
}

.modal-header h3 { font-size: 17px; margin: 0; }

.modal-close {
  background: none; border: none; font-size: 24px;
  cursor: pointer; color: var(--ink-soft);
}

.modal-body { padding: 20px; }
.modal-body p { margin-bottom: 8px; }

.text-danger { color: var(--signal-out); font-weight: 600; }
.text-warning { color: #f59e0b; font-weight: 600; }

.btn--xs {
  padding: 4px 8px;
  font-size: 11px;
}

.btn--warning {
  background: #f59e0b;
  color: white;
  border: none;
}

.btn--warning:hover {
  background: #d97706;
}

.btn--secondary {
  background: var(--brand);
  color: white;
  border: none;
}

.btn--secondary:hover {
  background: var(--brand-dark);
}

.modal-footer {
  display: flex; gap: 10px;
  justify-content: flex-end;
  padding: 16px 20px;
  border-top: 1px solid var(--line);
}

.employees-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-top: 12px;
}

.employee-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 14px;
  background: var(--surface);
  border: 1px solid var(--line);
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: all 0.2s;
  text-align: right;
  font-family: var(--font-body);
  font-size: 14px;
  width: 100%;
  position: relative;
}

.employee-item:hover {
  border-color: var(--brand);
  background: var(--brand-tint);
}

.employee-item:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.employee-item__avatar {
  width: 36px; height: 36px;
  border-radius: 50%;
  background: var(--brand-tint);
  color: var(--brand-dark);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.employee-item__avatar :deep(svg) {
  color: var(--brand-dark);
  stroke: var(--brand-dark);
}

.employee-item div { 
  display: flex; 
  flex-direction: column; 
  text-align: right; 
  flex: 1;
}

.employee-item span { font-size: 12px; color: var(--ink-soft); }

.employee-item__check {
  opacity: 0;
  color: var(--signal-in);
  font-size: 18px;
  transition: opacity 0.2s;
}

.employee-item:hover .employee-item__check {
  opacity: 1;
}

@media (max-width: 960px) {
  .sites-grid { grid-template-columns: 1fr 1fr; }
}

@media (max-width: 600px) {
  .sites-grid { grid-template-columns: 1fr; }
  
  /* تصغير أيقونات emoji في البطاقات */
  .site-card__details {
    font-size: 11px;
  }
  
  .site-card__radius {
    font-size: 11px;
    padding: 2px 8px;
  }
  
  /* تصغير الأيقونات في badges */
  .badge {
    font-size: 11px;
    padding: 3px 8px;
  }
  
  .badge--compact {
    font-size: 10px;
    padding: 2px 6px;
  }
  
  /* تصغير أزرار الموظفين العاملين */
  .btn--xs {
    font-size: 10px;
    padding: 3px 6px;
  }
  
  /* تحسين عرض عناصر الموظفين العاملين */
  .working-employee-item {
    flex-wrap: wrap;
  }
  
  .working-employee-item .badge {
    flex: 1;
    min-width: 0;
  }
  
  /* Mobile specific adjustments for dropdown */
  .dropdown-menu {
    min-width: 180px;
    right: 0;
  }
}

/* Success Modal Styles */
.modal-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  animation: fadeIn 0.2s ease;
}

.success-modal-card {
  background: var(--surface);
  border-radius: var(--radius-lg);
  padding: 32px;
  max-width: 400px;
  width: 90%;
  box-shadow: var(--shadow-xl);
  animation: scaleIn 0.2s ease;
  text-align: center;
}

.success-icon-wrapper {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 64px;
  height: 64px;
  margin: 0 auto 20px;
  border-radius: 50%;
  background: var(--signal-in-tint);
}

.success-icon-wrapper :deep(svg) {
  color: var(--signal-in);
  stroke: var(--signal-in);
}

.success-title {
  font-size: 20px;
  font-weight: 700;
  color: var(--ink);
  margin-bottom: 12px;
}

.success-message {
  font-size: 14px;
  color: var(--ink-soft);
  margin-bottom: 24px;
  line-height: 1.6;
}

.success-confirm-btn {
  width: 100%;
}

[data-theme="dark"] .success-icon-wrapper {
  background: rgba(16, 185, 129, 0.2);
}

[data-theme="dark"] .success-icon-wrapper :deep(svg) {
  color: var(--signal-in);
  stroke: var(--signal-in);
}

[data-theme="dark"] .success-title {
  color: var(--ink);
}

[data-theme="dark"] .success-message {
  color: var(--ink-soft);
}

/* بطاقة إنهاء الدوام الاحترافية */
.completion-card {
  background: var(--surface);
  border-radius: var(--radius-lg);
  padding: 0;
  max-width: 500px;
  width: 90%;
  box-shadow: var(--shadow-xl);
  animation: scaleIn 0.3s ease;
  overflow: hidden;
  border: 1px solid var(--line);
}

.completion-header {
  background: linear-gradient(135deg, #10b981 0%, #059669 100%);
  padding: 24px 20px;
  text-align: center;
  color: white;
}

.completion-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 64px;
  height: 64px;
  margin: 0 auto 12px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.25);
  backdrop-filter: blur(10px);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
}

.completion-icon :deep(svg) {
  color: white;
  stroke: white;
  stroke-width: 2.5;
}

.completion-title {
  font-size: 22px;
  font-weight: 700;
  margin: 0 0 6px 0;
  color: white;
  font-family: system-ui, -apple-system, 'Segoe UI', Roboto, sans-serif;
  letter-spacing: -0.5px;
}

.completion-subtitle {
  font-size: 14px;
  margin: 0;
  color: rgba(255, 255, 255, 0.95);
  font-weight: 400;
}

.completion-body {
  padding: 20px 20px;
  background: var(--surface);
}

.completion-info-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
  margin-bottom: 16px;
}

.completion-info-item {
  background: var(--surface);
  border: 1px solid var(--line);
  border-radius: var(--radius-md);
  padding: 12px;
  transition: all 0.2s ease;
}

.completion-info-item:hover {
  border-color: var(--brand);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.completion-info-label {
  font-size: 11px;
  color: var(--ink-soft);
  margin-bottom: 4px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.completion-info-value {
  font-size: 14px;
  font-weight: 600;
  color: var(--ink);
  font-family: system-ui, -apple-system, 'Segoe UI', Roboto, sans-serif;
  line-height: 1.4;
}

.completion-info-value.mono {
  font-family: 'SF Mono', 'Monaco', 'Inconsolata', 'Roboto Mono', monospace;
  font-size: 13px;
}

.completion-hours-section {
  background: linear-gradient(135deg, #f0fdf4 0%, #dcfce7 100%);
  border: 2px solid #10b981;
  border-radius: var(--radius-md);
  padding: 18px;
  text-align: center;
  margin-top: 16px;
  box-shadow: 0 4px 12px rgba(16, 185, 129, 0.1);
}

.completion-hours-label {
  font-size: 12px;
  color: #065f46;
  margin-bottom: 6px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.completion-hours-value {
  font-size: 28px;
  font-weight: 800;
  color: #059669;
  font-family: system-ui, -apple-system, 'Segoe UI', Roboto, sans-serif;
  letter-spacing: -1px;
}

.completion-footer {
  padding: 16px 20px;
  border-top: 1px solid var(--line);
  background: var(--surface);
}

.completion-confirm-btn {
  width: 100%;
  padding: 12px;
  font-size: 15px;
  font-weight: 600;
  border-radius: var(--radius-md);
  background: linear-gradient(135deg, #10b981 0%, #059669 100%);
  color: white;
  border: none;
  cursor: pointer;
  transition: all 0.2s ease;
  font-family: system-ui, -apple-system, 'Segoe UI', Roboto, sans-serif;
}

.completion-confirm-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(16, 185, 129, 0.3);
}

.completion-confirm-btn:active {
  transform: translateY(0);
}

[data-theme="dark"] .completion-card {
  background: var(--surface);
  border: 1px solid var(--line);
}

[data-theme="dark"] .completion-header {
  background: linear-gradient(135deg, #059669 0%, #047857 100%);
}

[data-theme="dark"] .completion-info-item {
  background: rgba(255, 255, 255, 0.03);
  border-color: var(--line);
}

[data-theme="dark"] .completion-info-item:hover {
  border-color: var(--brand);
  background: rgba(255, 255, 255, 0.05);
}

[data-theme="dark"] .completion-info-label {
  color: var(--ink-soft);
}

[data-theme="dark"] .completion-info-value {
  color: var(--ink);
}

[data-theme="dark"] .completion-hours-section {
  background: rgba(16, 185, 129, 0.15);
  border-color: #10b981;
}

[data-theme="dark"] .completion-hours-label {
  color: #6ee7b7;
}

[data-theme="dark"] .completion-hours-value {
  color: #34d399;
}

[data-theme="dark"] .completion-footer {
  background: rgba(255, 255, 255, 0.02);
  border-color: var(--line);
}

[data-theme="dark"] .completion-confirm-btn {
  background: linear-gradient(135deg, #059669 0%, #047857 100%);
}

@keyframes scaleIn {
  from {
    transform: scale(0.9);
    opacity: 0;
  }
  to {
    transform: scale(1);
    opacity: 1;
  }
}

@media (max-width: 600px) {
  .completion-info-row {
    grid-template-columns: 1fr;
    gap: 10px;
  }
  
  .completion-card {
    max-width: 95%;
  }
  
  .completion-header {
    padding: 18px 14px;
  }
  
  .completion-icon {
    width: 50px;
    height: 50px;
    margin: 0 auto 10px;
  }
  
  .completion-title {
    font-size: 18px;
  }
  
  .completion-subtitle {
    font-size: 12px;
  }
  
  .completion-body {
    padding: 14px;
  }
  
  .completion-info-item {
    padding: 10px;
  }
  
  .completion-info-label {
    font-size: 10px;
    margin-bottom: 3px;
  }
  
  .completion-info-value {
    font-size: 12px;
  }
  
  .completion-info-value.mono {
    font-size: 11px;
  }
  
  .completion-hours-section {
    padding: 14px;
    margin-top: 12px;
  }
  
  .completion-hours-label {
    font-size: 10px;
    margin-bottom: 4px;
  }
  
  .completion-hours-value {
    font-size: 22px;
  }
  
  .completion-footer {
    padding: 12px 14px;
  }
  
  .completion-confirm-btn {
    padding: 10px;
    font-size: 13px;
  }
}
</style>
