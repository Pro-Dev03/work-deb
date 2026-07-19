<template>
  <div>
    <div class="page-head">
      <div>
        <h2>{{ t('worksites_title') }}</h2>
        <p>{{ t('worksites_description') }}</p>
      </div>
      <button class="btn btn--primary" @click="showModal = true">+ {{ t('new_worksite') }}</button>
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
            <button class="btn btn--primary btn--sm" @click="openAssignModal(site)">
              👤 {{ t('assign_employee') }}
            </button>
            <button class="btn btn--danger btn--sm" @click="confirmDelete(site)">
              🗑️
            </button>
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

    <WorksiteFormModal v-if="showModal" @close="showModal = false" @worksite-added="fetchWorksites" />

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
              <span class="employee-item__avatar">{{ emp.full_name.slice(0, 1) }}</span>
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
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '../services/api'
import WorksiteFormModal from '../components/WorksiteFormModal.vue'
import { useI18n } from '../services/i18n'

const { t } = useI18n()
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
const siteToDelete = ref(null)
const worksiteToAssign = ref(null)
const employeeToForceCheckout = ref(null)

async function fetchWorksites() {
  loading.value = true
  try {
    const { data } = await api.get('/worksites')
    worksites.value = data || []
  } catch (error) {
    console.error('❌ ' + t('failed_to_fetch_worksites'), error)
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

async function fetchEmployees() {
  loadingEmployees.value = true
  try {
    const { data } = await api.get('/worksites/employees')
    employees.value = data || []
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
    alert('✅ ' + (response.data.message || t('employee_assigned_successfully')))
    await fetchWorksites()
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
    alert('✅ ' + (response.data.message || t('force_checkout_success')))
    await fetchWorksites()
  } catch (error) {
    console.error('❌ ' + t('force_checkout_failed'), error)
    const msg = error.response?.data?.error || t('force_checkout_failed')
    alert('❌ ' + msg)
  } finally {
    forceCheckingOut.value = false
  }
}

onMounted(fetchWorksites)
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
}

.site-card:hover {
  box-shadow: var(--shadow-md);
  transform: translateY(-3px);
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
  font-weight: 700;
  font-size: 16px;
  flex-shrink: 0;
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
}
</style>
