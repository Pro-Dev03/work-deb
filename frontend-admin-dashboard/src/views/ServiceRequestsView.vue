<template>
  <div>
    <div class="page-head">
      <div>
        <h2>📋 {{ t('service_requests') }}</h2>
        <p>{{ t('service_requests_description') }}</p>
      </div>
      <div class="filters">
        <select v-model="filterStatus" @change="fetchRequests">
          <option value="">{{ t('all_statuses') }}</option>
          <option value="pending">{{ t('status_pending') }}</option>
          <option value="assigned">{{ t('status_assigned') }}</option>
          <option value="in_progress">{{ t('status_in_progress') }}</option>
          <option value="completed">{{ t('status_completed') }}</option>
        </select>
      </div>
    </div>

    <div v-if="loading" class="empty-state"><p>⏳ {{ t('loading') }}...</p></div>

    <div v-else-if="requests.length === 0" class="empty-state">
      <h3>📭 {{ t('no_service_requests') }}</h3>
      <p>{{ t('no_service_requests_hint') }}</p>
    </div>

    <div v-else>
      <div v-for="req in filteredRequests" :key="req.id" class="card request-card">
        <div class="request-card__header">
          <div>
            <span class="badge" :class="getPriorityClass(req.priority)">
              {{ getPriorityLabel(req.priority) }}
            </span>
            <span class="badge" :class="getStatusClass(req.status)">
              {{ getStatusLabel(req.status) }}
            </span>
          </div>
          <span class="request-card__time mono">{{ formatDate(req.created_at) }}</span>
        </div>

        <h3>{{ req.title }}</h3>
        <p class="request-card__desc">{{ req.description }}</p>

        <div class="request-card__info">
          <span>👤 {{ req.client_name || t('client') }}</span>
          <span>📞 {{ req.client_phone || req.phone || '—' }}</span>
          <span>📍 {{ req.address || t('no_address') }}</span>
        </div>

        <div class="request-card__location">
          <span class="mono">{{ t('latitude') }}: {{ req.latitude.toFixed(5) }}</span>
          <span class="mono">{{ t('longitude') }}: {{ req.longitude.toFixed(5) }}</span>
        </div>

        <div class="request-card__actions">
          <button v-if="req.status === 'pending'" class="btn btn--primary btn--sm" @click="openAssignModal(req)">
            👤 {{ t('assign_employee') }}
          </button>
          <button v-if="req.status === 'assigned'" class="btn btn--ghost btn--sm">
            ⏳ {{ t('waiting_for_employee') }}
          </button>
          <button v-if="req.status === 'in_progress'" class="btn btn--gold btn--sm">
            🔄 {{ t('in_execution') }}
          </button>
          <button v-if="req.status === 'completed'" class="btn btn--success btn--sm">
            ✅ {{ t('status_completed') }}
          </button>
        </div>
      </div>
    </div>

    <!-- مودال تعيين الموظف -->
    <div v-if="showAssignModal" class="modal-backdrop" @click.self="showAssignModal = false">
      <div class="modal card">
        <h3>👤 {{ t('assign_employee_modal') }}</h3>
        <p>{{ t('assign_employee_hint') }}</p>

        <div v-if="loadingEmployees" class="empty-state"><p>⏳ {{ t('loading_employees') }}...</p></div>

        <div v-else class="employees-list">
          <button
            v-for="emp in employees"
            :key="emp.id"
            class="employee-item"
            @click="assignEmployee(emp.id)"
          >
            <span class="employee-item__avatar">{{ emp.full_name.slice(0, 1) }}</span>
            <div>
              <strong>{{ emp.full_name }}</strong>
            </div>
          </button>
        </div>

        <button class="btn btn--ghost btn--block" @click="showAssignModal = false">{{ t('cancel') }}</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import api from '../services/api'
import { useI18n } from '../services/i18n'

const { t } = useI18n()

const requests = ref([])
const employees = ref([])
const loading = ref(false)
const loadingEmployees = ref(false)
const filterStatus = ref('')
const showAssignModal = ref(false)
const selectedRequest = ref(null)

const filteredRequests = computed(() => {
  if (!filterStatus.value) return requests.value
  return requests.value.filter(r => r.status === filterStatus.value)
})

async function fetchRequests() {
  loading.value = true
  try {
    const { data } = await api.get('/service/requests')
    requests.value = data || []
  } catch (error) {
    console.error(t('failed_to_fetch_requests'), error)
  } finally {
    loading.value = false
  }
}

async function fetchEmployees() {
  loadingEmployees.value = true
  try {
    const { data } = await api.get('/service/employees')
    employees.value = data || []
  } catch (error) {
    console.error(t('failed_to_fetch_employees'), error)
  } finally {
    loadingEmployees.value = false
  }
}

function openAssignModal(req) {
  selectedRequest.value = req
  showAssignModal.value = true
  fetchEmployees()
}

async function assignEmployee(employeeId) {
  try {
    await api.post('/service/assign', {
      request_id: selectedRequest.value.id,
      employee_id: employeeId
    })
    showAssignModal.value = false
    await fetchRequests()
  } catch (error) {
    console.error('فشل التعيين:', error)
  }
}

function getPriorityLabel(priority) {
  const labels = { low: 'منخفضة', normal: 'عادية', high: 'عالية', urgent: 'طارئة' }
  return labels[priority] || priority
}

function getPriorityClass(priority) {
  const classes = { low: 'badge--gray', normal: 'badge--blue', high: 'badge--gold', urgent: 'badge--out' }
  return classes[priority] || ''
}

function getStatusLabel(status) {
  const labels = { pending: 'قيد الانتظار', assigned: 'تم التعيين', in_progress: 'قيد التنفيذ', completed: 'مكتمل' }
  return labels[status] || status
}

function getStatusClass(status) {
  const classes = { pending: 'badge--gold', assigned: 'badge--blue', in_progress: 'badge--info', completed: 'badge--in' }
  return classes[status] || ''
}

function formatDate(date) {
  if (!date) return '—'
  return new Date(date).toLocaleString('ar-SA')
}

onMounted(fetchRequests)
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
.filters select {
  padding: 8px 14px;
  border-radius: var(--radius-sm);
  border: 1px solid var(--line);
  background: var(--surface);
  font-family: var(--font-body);
  min-width: 150px;
}

@media (max-width: 768px) {
  .page-head { flex-direction: column; align-items: flex-start; }
  .filters { width: 100%; }
  .filters select { width: 100%; }
  
  .request-card__location {
    flex-direction: column;
    gap: 8px;
  }
  
  .request-card__actions {
    flex-direction: column;
  }
  
  .request-card__actions .btn {
    width: 100%;
  }
}

.request-card {
  padding: 18px 20px;
  margin-bottom: 12px;
  transition: border-color 0.2s;
}
.request-card:hover { border-color: var(--brand-tint); }

.request-card__header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}
.request-card__header .badge { margin-inline-end: 6px; }
.request-card__time { font-size: 12px; color: var(--ink-soft); }

.request-card h3 { font-size: 16px; margin-bottom: 6px; }
.request-card__desc { font-size: 13px; color: var(--ink-soft); margin-bottom: 10px; }

.request-card__info {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  font-size: 13px;
  color: var(--ink-soft);
  margin-bottom: 8px;
}

.request-card__location {
  display: flex;
  gap: 16px;
  font-size: 12px;
  color: var(--ink-soft);
  background: var(--canvas);
  padding: 8px 12px;
  border-radius: var(--radius-sm);
  margin-bottom: 12px;
}

.request-card__actions {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.badge--gray { background: var(--line); color: var(--ink-soft); }
.badge--blue { background: #E3ECF7; color: #2C6B9E; }
.badge--info { background: #E3F0F7; color: #1A7A8A; }
.badge--success { background: var(--signal-in-tint); color: var(--signal-in); }

.modal-backdrop {
  position: fixed; inset: 0;
  background: rgba(22,35,46,0.6);
  backdrop-filter: blur(4px);
  display: flex; align-items: center; justify-content: center;
  z-index: 50; padding: 20px;
}
.modal {
  width: 100%; max-width: 420px; padding: 24px;
  max-height: 80vh; overflow-y: auto;
}
.modal h3 { font-size: 17px; margin-bottom: 4px; }
.modal p { font-size: 13px; color: var(--ink-soft); margin-bottom: 16px; }

.employees-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 16px;
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
}
.employee-item:hover {
  border-color: var(--brand);
  background: var(--brand-tint);
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
.employee-item div { display: flex; flex-direction: column; }
.employee-item span { font-size: 12px; color: var(--ink-soft); }
</style>
