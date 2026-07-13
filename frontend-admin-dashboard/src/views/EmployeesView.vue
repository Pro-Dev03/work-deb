<template>
  <div>
    <div class="page-head">
      <div>
        <h2>{{ t('employees_title') }}</h2>
        <p>{{ t('employees_description') }}</p>
      </div>
      <button class="btn btn--primary" @click="showModal = true">+ {{ t('add_employee') }}</button>
    </div>

    <div v-if="loading" class="empty-state">
      <p>{{ t('loading_employees') }}</p>
    </div>

    <div v-else-if="employees.length === 0" class="empty-state">
      <h3>{{ t('no_employees') }}</h3>
      <p>{{ t('add_new_employee_prompt') }}</p>
    </div>

    <div v-else class="card">
      <div class="table-wrapper">
        <table class="table">
          <thead>
            <tr>
              <th>{{ t('name') }}</th>
              <th>{{ t('email') }}</th>
              <th>{{ t('phone') }}</th>
              <th>{{ t('role') }}</th>
              <th>{{ t('status') }}</th>
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
              <td class="mono">{{ emp.email }}</td>
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
              <td class="mono">{{ formatDate(emp.created_at) }}</td>
              <td>
                <button 
                  class="btn btn--danger btn--sm" 
                  @click="confirmDelete(emp)"
                  :disabled="emp.role === 'admin' && emp.email === 'admin@worktrack.com'"
                >
                  🗑️ {{ t('delete') }}
                </button>
              </td>
            </tr>
          </tbody>
        </table>
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
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '../services/api'
import EmployeeFormModal from '../components/EmployeeFormModal.vue'
import { useI18n } from '../services/i18n'

const { t } = useI18n()
const employees = ref([])
const loading = ref(false)
const showModal = ref(false)
const showDeleteModal = ref(false)
const employeeToDelete = ref(null)
const deleting = ref(false)

async function fetchEmployees() {
  loading.value = true
  try {
    const { data } = await api.get('/admin/employees')
    employees.value = data || []
  } catch (error) {
    console.error('فشل جلب الموظفين:', error)
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
    console.error('فشل حذف الموظف:', error)
    alert(error.response?.data?.error || 'فشل حذف الموظف')
  } finally {
    deleting.value = false
  }
}

function formatDate(date) {
  if (!date) return '—'
  return new Date(date).toLocaleDateString('ar-SA')
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

.table-wrapper { overflow-x: auto; }

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

.modal-header {
  display: flex; justify-content: space-between; align-items: center;
  padding: 16px 20px; border-bottom: 1px solid var(--line);
}

.modal-header h3 { font-size: 17px; margin: 0; }
.modal-close { background: none; border: none; font-size: 24px; cursor: pointer; color: var(--ink-soft); }

.modal-body { padding: 20px; }
.modal-body p { margin-bottom: 8px; }
.text-danger { color: var(--signal-out); font-weight: 600; }

.modal-footer {
  display: flex; gap: 10px;
  justify-content: flex-end;
  padding: 16px 20px;
  border-top: 1px solid var(--line);
}

@media (max-width: 768px) {
  .page-head { flex-direction: column; }
  .page-head .btn { width: 100%; }
}
</style>
