<template>
  <div class="notes-page">
    <!-- Header Section -->
    <div class="notes-header">
      <div class="header-content">
        <div class="header-icon">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
            <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
          </svg>
        </div>
        <div class="header-text">
          <h1>{{ t('notes_title') }}</h1>
          <p>{{ t('notes_desc') }}</p>
        </div>
      </div>
      
      <!-- Stats Cards -->
      <div class="stats-grid">
        <div class="stat-card total">
          <div class="stat-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
              <polyline points="14 2 14 8 20 8"/>
              <line x1="16" y1="13" x2="8" y2="13"/>
              <line x1="16" y1="17" x2="8" y2="17"/>
              <polyline points="10 9 9 9 8 9"/>
            </svg>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ notes.length }}</div>
            <div class="stat-label">{{ t('total_notes') }}</div>
          </div>
        </div>
        <div class="stat-card unread">
          <div class="stat-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"/>
              <line x1="12" y1="8" x2="12" y2="12"/>
              <line x1="12" y1="16" x2="12.01" y2="16"/>
            </svg>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ notes.filter(n => !n.is_read).length }}</div>
            <div class="stat-label">{{ t('unread_notes') }}</div>
          </div>
        </div>
        <div class="stat-card read">
          <div class="stat-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/>
              <polyline points="22 4 12 14.01 9 11.01"/>
            </svg>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ notes.filter(n => n.is_read).length }}</div>
            <div class="stat-label">{{ t('read_notes') }}</div>
          </div>
        </div>
      </div>
    </div>

    <div class="notes-content">
      <!-- Create Note Section -->
      <div class="create-note-card">
        <div class="card-header">
          <div class="card-title">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
              <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
            </svg>
            <h3>{{ t('create_new_note') }}</h3>
          </div>
        </div>
        
        <div class="card-body">
          <div class="form-row">
            <div class="form-group employees-group">
              <label class="form-label">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
                  <circle cx="9" cy="7" r="4"/>
                  <path d="M23 21v-2a4 4 0 0 0-3-3.87"/>
                  <path d="M16 3.13a4 4 0 0 1 0 7.75"/>
                </svg>
                {{ t('select_employees') }}
              </label>
              <div class="employees-selection">
                <div class="select-all-wrapper">
                  <label class="checkbox-label">
                    <input type="checkbox" v-model="selectAll" @change="toggleSelectAll" class="custom-checkbox"/>
                    <span class="checkbox-text">
                      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
                        <circle cx="9" cy="7" r="4"/>
                        <path d="M23 21v-2a4 4 0 0 0-3-3.87"/>
                        <path d="M16 3.13a4 4 0 0 1 0 7.75"/>
                      </svg>
                      {{ t('select_all_employees') }}
                    </span>
                  </label>
                </div>
                <div class="employees-list">
                  <label v-for="emp in employees" :key="emp.id" class="employee-item">
                    <div class="employee-checkbox">
                      <input 
                        type="checkbox" 
                        v-model="selectedEmployees" 
                        :value="emp.id"
                        class="custom-checkbox"
                      />

                      <div class="employee-info">
                        <span class="employee-name">{{ emp.full_name }}</span>
                        <span v-if="emp.phone" class="employee-phone">{{ emp.phone }}</span>
                      </div>
                      <a 
                        v-if="emp.phone" 
                        :href="`https://wa.me/${emp.phone.replace(/[^0-9]/g, '')}`" 
                        target="_blank"
                        class="whatsapp-link"
                        title="تواصل عبر واتساب"
                        @click.stop
                      >
                        <svg viewBox="0 0 24 24" fill="currentColor">
                          <path d="M17.472 14.382c-.297-.149-1.758-.867-2.03-.967-.273-.099-.471-.148-.67.15-.197.297-.767.966-.94 1.164-.173.199-.347.223-.644.075-.297-.15-1.255-.463-2.39-1.475-.883-.788-1.48-1.761-1.653-2.059-.173-.297-.018-.458.13-.606.134-.133.298-.347.446-.52.149-.174.198-.298.298-.497.099-.198.05-.371-.025-.52-.075-.149-.669-1.612-.916-2.207-.242-.579-.487-.5-.669-.51-.173-.008-.371-.01-.57-.01-.198 0-.52.074-.792.372-.272.297-1.04 1.016-1.04 2.479 0 1.462 1.065 2.875 1.213 3.074.149.198 2.096 3.2 5.077 4.487.709.306 1.262.489 1.694.625.712.227 1.36.195 1.871.118.571-.085 1.758-.719 2.006-1.413.248-.694.248-1.289.173-1.413-.074-.124-.272-.198-.57-.347m-5.421 7.403h-.004a9.87 9.87 0 01-5.031-1.378l-.361-.214-3.741.982.998-3.648-.235-.374a9.86 9.86 0 01-1.51-5.26c.001-5.45 4.436-9.884 9.888-9.884 2.64 0 5.122 1.03 6.988 2.898a9.825 9.825 0 012.893 6.994c-.003 5.45-4.437 9.884-9.885 9.884m8.413-18.297A11.815 11.815 0 0012.05 0C5.495 0 .16 5.335.157 11.892c0 2.096.547 4.142 1.588 5.945L.057 24l6.305-1.654a11.882 11.882 0 005.683 1.448h.005c6.554 0 11.89-5.335 11.893-11.893a11.821 11.821 0 00-3.48-8.413z"/>
                        </svg>
                      </a>
                    </div>
                  </label>
                </div>
              </div>
              <div class="selection-summary">
                <div v-if="selectedEmployees.length > 0" class="selected-info">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M16 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
                    <circle cx="8.5" cy="7" r="4"/>
                    <line x1="20" y1="8" x2="20" y2="14"/>
                    <line x1="23" y1="11" x2="17" y2="11"/>
                  </svg>
                  <span>{{ t('selected_employees_count', { count: selectedEmployees.length }) }}</span>
                  <button class="clear-btn" @click="clearSelection">{{ t('clear') }}</button>
                </div>
                <div v-else class="no-selection">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <circle cx="12" cy="12" r="10"/>
                    <line x1="12" y1="8" x2="12" y2="12"/>
                    <line x1="12" y1="16" x2="12.01" y2="16"/>
                  </svg>
                  <span>{{ t('no_employee_selected') }}</span>
                </div>
              </div>
            </div>
          </div>
          
          <div class="form-row">
            <div class="form-group note-group">
              <label class="form-label">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/>
                </svg>
                {{ t('note_text') }}
              </label>
              <textarea
                v-model="newNote.content"
                class="form-textarea"
                :placeholder="t('write_note_here')"
                rows="4"
              ></textarea>
            </div>
          </div>
          
          <div class="form-row">
            <div class="form-group worksite-group">
              <label class="form-label">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"/>
                  <circle cx="12" cy="10" r="3"/>
                </svg>
                {{ t('worksite_optional') }}
              </label>
              <select v-model="newNote.worksite_id" class="form-select">
                <option value="">{{ t('no_worksite') }}</option>
                <option v-for="ws in worksites" :key="ws.id" :value="ws.id">{{ ws.name }}</option>
              </select>
            </div>
          </div>
          
          <div class="form-actions">
            <button
              class="btn btn--primary send-btn"
              @click="createNote"
              :disabled="selectedEmployees.length === 0 || !newNote.content || creating"
            >
              <svg v-if="!creating" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="22" y1="2" x2="11" y2="13"/>
                <polygon points="22 2 15 22 11 13 2 9 22 2"/>
              </svg>
              <svg v-else class="loading-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M12 2v4m0 12v4M4.93 4.93l2.83 2.83m8.48 8.48l2.83 2.83M2 12h4m12 0h4M4.93 19.07l2.83-2.83m8.48-8.48l2.83-2.83"/>
              </svg>
              {{ creating ? 'جاري الإرسال...' : `${t('send_to_employees')} (${selectedEmployees.length})` }}
            </button>
          </div>
        </div>
      </div>

      <!-- Filters Section -->
      <div class="filters-card">
        <div class="filter-header">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polygon points="22 3 2 3 10 12.46 10 19 14 21 14 12.46 22 3"/>
          </svg>
          <h3>{{ t('filter_notes') }}</h3>
        </div>
        <div class="filter-content">
          <div class="filter-group">
            <label>{{ t('filter_by_employee') }}</label>
            <select v-model="filterEmployeeId" @change="fetchNotes" class="form-select">
              <option value="">{{ t('select_all_employees') }}</option>
              <option v-for="emp in employees" :key="emp.id" :value="emp.id">{{ emp.full_name }}</option>
            </select>
          </div>
        </div>
      </div>

      <!-- Notes List -->
      <div v-if="loading" class="loading-state">
        <div class="loading-spinner"></div>
        <p>جارٍ تحميل الملاحظات...</p>
      </div>

      <div v-else-if="notes.length === 0" class="empty-state">
        <div class="empty-icon">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
            <polyline points="14 2 14 8 20 8"/>
            <line x1="12" y1="18" x2="12" y2="12"/>
            <line x1="9" y1="15" x2="15" y2="15"/>
          </svg>
        </div>
        <h3>لا توجد ملاحظات</h3>
      </div>

      <div v-else class="notes-grid">
        <div 
          v-for="note in notes" 
          :key="note.id" 
          class="note-card"
          :class="{ 'unread': !note.is_read }"
        >
          <div class="note-card-header">
            <div class="note-recipient">
              <div class="note-avatar">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
                  <circle cx="12" cy="7" r="4"/>
                </svg>
              </div>
              <div class="note-recipient-info">
                <span class="note-recipient-name">{{ note.employee_name }}</span>
                <span v-if="note.employee_phone" class="note-recipient-phone">{{ note.employee_phone }}</span>
              </div>
            </div>
            <div class="note-actions">
              <button 
                class="action-btn edit-btn" 
                @click="openEditModal(note)"
                title="تعديل"
              >
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
                  <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
                </svg>
              </button>
              <button 
                class="action-btn delete-btn" 
                @click="deleteNote(note.id)"
                title="حذف"
              >
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="3 6 5 6 21 6"/>
                  <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
                </svg>
              </button>
            </div>
          </div>
          
          <div class="note-card-body">
            <div class="note-content">{{ note.content }}</div>
            <div v-if="note.worksite_name" class="note-worksite">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"/>
                <circle cx="12" cy="10" r="3"/>
              </svg>
              {{ note.worksite_name }}
            </div>
          </div>
          
          <div class="note-card-footer">
            <div class="note-date">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="3" y="4" width="18" height="18" rx="2" ry="2"/>
                <line x1="16" y1="2" x2="16" y2="6"/>
                <line x1="8" y1="2" x2="8" y2="6"/>
                <line x1="3" y1="10" x2="21" y2="10"/>
              </svg>
              {{ formatDate(note.created_at) }}
            </div>
            <div v-if="note.updated_at !== note.created_at" class="note-updated">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
                <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
              </svg>
              تم التعديل
            </div>
            <div class="note-read-status" :class="{ 'read': note.is_read }">
              <svg v-if="note.is_read" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/>
                <polyline points="22 4 12 14.01 9 11.01"/>
              </svg>
              <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="10"/>
                <line x1="12" y1="8" x2="12" y2="12"/>
                <line x1="12" y1="16" x2="12.01" y2="16"/>
              </svg>
              {{ note.is_read ? t('read_notes') : t('unread_notes') }}
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Edit Modal -->
    <div v-if="showEditModal" class="modal-backdrop" @click.self="showEditModal = false">
      <div class="modal-card">
        <div class="modal-header">
          <div class="modal-title">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
              <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
            </svg>
            <h3>تعديل الملاحظة</h3>
          </div>
          <button class="modal-close" @click="showEditModal = false">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18"/>
              <line x1="6" y1="6" x2="18" y2="18"/>
            </svg>
          </button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label class="form-label">{{ t('note_text') }}</label>
            <textarea
              v-model="editingNote.content"
              class="form-textarea"
              rows="4"
            ></textarea>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn--ghost" @click="showEditModal = false">إلغاء</button>
          <button 
            class="btn btn--primary" 
            @click="updateNote"
            :disabled="!editingNote.content || updating"
          >
            <svg v-if="!updating" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z"/>
              <polyline points="17 21 17 13 7 13 7 21"/>
              <polyline points="7 3 7 8 15 8"/>
            </svg>
            <svg v-else class="loading-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M12 2v4m0 12v4M4.93 4.93l2.83 2.83m8.48 8.48l2.83 2.83M2 12h4m12 0h4M4.93 19.07l2.83-2.83m8.48-8.48l2.83-2.83"/>
            </svg>
            {{ updating ? 'جاري الحفظ...' : 'حفظ التغييرات' }}
          </button>
        </div>
      </div>
    </div>

    <!-- Success Modal -->
    <div v-if="showSuccessModal" class="modal-backdrop" @click.self="showSuccessModal = false">
      <div class="success-modal-card">
        <div class="success-animation">
          <div class="success-circle">
            <svg class="success-check" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
              <polyline points="20 6 9 17 4 12"/>
            </svg>
          </div>
        </div>
        <h3 class="success-title">تم بنجاح!</h3>
        <p class="success-message">{{ successMessage }}</p>
        <button class="btn btn--primary success-btn" @click="showSuccessModal = false">
          حسناً
        </button>
      </div>
    </div>

    <!-- Delete Confirmation Modal -->
    <div v-if="showDeleteModal" class="modal-overlay" @click.self="cancelDelete">
      <div class="delete-modal-card">
        <div class="delete-icon-wrapper">
          <svg class="delete-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10"/>
            <line x1="15" y1="9" x2="9" y2="15"/>
            <line x1="9" y1="9" x2="15" y2="15"/>
          </svg>
        </div>
        <h3 class="delete-title">{{ t('delete_note_title') }}</h3>
        <p class="delete-message">{{ t('delete_note_message') }}</p>
        <div class="delete-actions">
          <button class="btn btn--ghost cancel-btn" @click="cancelDelete">
            {{ t('cancel') }}
          </button>
          <button class="btn btn--danger delete-confirm-btn" @click="confirmDelete">
            {{ t('delete') }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, watch } from 'vue'
import api from '../services/api'
import { useI18n } from '../services/i18n'
const { t, currentLang } = useI18n()

const loading = ref(false)
const creating = ref(false)
const updating = ref(false)
const notes = ref([])
const employees = ref([])
const worksites = ref([])
const showEditModal = ref(false)
const showSuccessModal = ref(false)
const successMessage = ref('')
const showDeleteModal = ref(false)
const noteToDelete = ref(null)
const filterEmployeeId = ref('')
const selectedEmployees = ref([])
const selectAll = ref(false)

const newNote = reactive({
  content: '',
  worksite_id: ''
})

const editingNote = reactive({
  id: '',
  content: ''
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

function toggleSelectAll() {
  if (selectAll.value) {
    selectedEmployees.value = employees.value.map(emp => emp.id)
  } else {
    selectedEmployees.value = []
  }
}

function clearSelection() {
  selectedEmployees.value = []
  selectAll.value = false
}

// تحديث checkbox "جميع الموظفين" عند تغيير الاختيار الفردي
watch(selectedEmployees, (newVal) => {
  selectAll.value = newVal.length === employees.value.length && employees.value.length > 0
})

async function fetchNotes() {
  loading.value = true
  try {
    const params = new URLSearchParams()
    if (filterEmployeeId.value) {
      params.append('employee_id', filterEmployeeId.value)
    }

    const { data } = await api.get(`/admin/notes?${params.toString()}`)
    notes.value = data || []
  } catch (error) {
    console.error('❌ فشل جلب الملاحظات:', error)
  } finally {
    loading.value = false
  }
}

async function createNote() {
  if (selectedEmployees.value.length === 0) return
  if (!newNote.content) return

  creating.value = true
  try {
    // إرسال للموظفين المختارين
    const promises = selectedEmployees.value.map(empId =>
      api.post('/admin/notes', {
        employee_id: empId,
        content: newNote.content,
        worksite_id: newNote.worksite_id || null
      })
    )

    await Promise.all(promises)
    successMessage.value = `تم إرسال الملاحظة بنجاح لـ ${selectedEmployees.value.length} موظف`
    showSuccessModal.value = true

    // إعادة تعيين النموذج
    selectedEmployees.value = []
    selectAll.value = false
    newNote.content = ''
    newNote.worksite_id = ''

    // جلب الملاحظات المحدثة
    await fetchNotes()
  } catch (error) {
    console.error('❌ فشل إنشاء الملاحظة:', error)
    alert('فشل إنشاء الملاحظة')
  } finally {
    creating.value = false
  }
}

function openEditModal(note) {
  editingNote.id = note.id
  editingNote.content = note.content
  showEditModal.value = true
}

async function updateNote() {
  if (!editingNote.content) return

  updating.value = true
  try {
    await api.put(`/admin/notes/${editingNote.id}`, {
      content: editingNote.content
    })

    showEditModal.value = false
    await fetchNotes()
  } catch (error) {
    console.error('❌ فشل تعديل الملاحظة:', error)
    alert('فشل تعديل الملاحظة')
  } finally {
    updating.value = false
  }
}

async function deleteNote(noteId) {
  showDeleteModal.value = true
  noteToDelete.value = noteId
}

async function confirmDelete() {
  if (!noteToDelete.value) return

  try {
    await api.delete(`/admin/notes/${noteToDelete.value}`)
    await fetchNotes()
    showDeleteModal.value = false
    noteToDelete.value = null
    showSuccessModal.value = true
    successMessage.value = t('delete_success')
  } catch (error) {
    console.error('❌ فشل حذف الملاحظة:', error)
    showDeleteModal.value = false
    noteToDelete.value = null
  }
}

function cancelDelete() {
  showDeleteModal.value = false
  noteToDelete.value = null
}

function formatDate(date) {
  if (!date) return ''
  const localeMap = {
    'ar': 'ar-SA',
    'he': 'he-IL',
    'en': 'en-US'
  }
  const locale = localeMap[currentLang.value] || 'en-US'
  return new Date(date).toLocaleString(locale, {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

onMounted(async () => {
  await Promise.all([
    fetchEmployees(),
    fetchWorksites(),
    fetchNotes()
  ])
})
</script>

<style scoped>
/* Page Layout */
.notes-page {
  min-height: 100vh;
  background: var(--canvas);
  padding: 24px;
}

/* Header Section */
.notes-header {
  margin-bottom: 32px;
}

.header-content {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 24px;
}

.header-icon {
  width: 64px;
  height: 64px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 8px 24px rgba(102, 126, 234, 0.3);
}

.header-icon svg {
  width: 32px;
  height: 32px;
  color: white;
}

.header-text h1 {
  margin: 0 0 4px 0;
  font-size: 28px;
  font-weight: 700;
  color: var(--ink);
}

.header-text p {
  margin: 0;
  font-size: 14px;
  color: var(--ink-soft);
}

/* Stats Grid */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
}

.stat-card {
  background: var(--surface);
  border-radius: 16px;
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 16px;
  box-shadow: var(--shadow-sm);
  transition: all 0.3s ease;
}

.stat-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
}

.stat-card.total .stat-icon {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.stat-card.unread .stat-icon {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.stat-card.read .stat-icon {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.stat-icon svg {
  width: 24px;
  height: 24px;
  color: white;
}

.stat-content {
  flex: 1;
}

.stat-value {
  font-size: 28px;
  font-weight: 700;
  color: var(--ink);
  line-height: 1;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 12px;
  color: var(--ink-soft);
  font-weight: 500;
}

/* Content Section */
.notes-content {
  display: grid;
  gap: 24px;
}

/* Cards */
.create-note-card,
.filters-card {
  background: var(--surface);
  border-radius: 20px;
  box-shadow: var(--shadow-md);
  overflow: hidden;
}

.card-header {
  padding: 20px 24px;
  border-bottom: 1px solid var(--line);
  background: linear-gradient(135deg, var(--canvas) 0%, var(--surface) 100%);
}

.card-title {
  display: flex;
  align-items: center;
  gap: 12px;
}

.card-title svg {
  width: 20px;
  height: 20px;
  color: #667eea;
}

.card-title h3 {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: var(--ink);
}

.card-body {
  padding: 24px;
}

/* Form Elements */
.form-row {
  margin-bottom: 20px;
}

.form-row:last-child {
  margin-bottom: 0;
}

.form-group {
  margin-bottom: 16px;
}

.form-label {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
  font-size: 13px;
  font-weight: 600;
  color: var(--ink);
}

.form-label svg {
  width: 16px;
  height: 16px;
  color: #667eea;
}

.form-select,
.form-textarea {
  width: 100%;
  padding: 12px 16px;
  border: 2px solid var(--line);
  border-radius: 12px;
  font-size: 14px;
  font-family: inherit;
  background: var(--surface);
  color: var(--ink);
  transition: all 0.3s ease;
}

.form-select:focus,
.form-textarea:focus {
  border-color: var(--brand);
  outline: none;
  box-shadow: 0 0 0 4px var(--brand-tint);
}

.form-textarea {
  resize: vertical;
  line-height: 1.6;
}

/* Employees Selection */
.employees-selection {
  border: 2px solid var(--line);
  border-radius: 12px;
  background: var(--canvas);
  max-height: 320px;
  overflow-y: auto;
}

.select-all-wrapper {
  padding: 12px 16px;
  border-bottom: 1px solid var(--line);
  background: linear-gradient(135deg, var(--canvas) 0%, var(--surface) 100%);
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 12px;
  cursor: pointer;
  user-select: none;
}

.checkbox-text {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  font-weight: 600;
  color: var(--ink);
}

.checkbox-text svg {
  width: 18px;
  height: 18px;
  color: #667eea;
}

.custom-checkbox {
  width: 20px;
  height: 20px;
  cursor: pointer;
  accent-color: #667eea;
}

.employees-list {
  padding: 8px;
}

.employee-item {
  display: block;
  padding: 8px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.employee-item:hover {
  background: rgba(102, 126, 234, 0.1);
}

.employee-checkbox {
  display: flex;
  align-items: center;
  gap: 12px;
}

.employee-avatar {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
}

.employee-avatar svg {
  width: 16px;
  height: 16px;
  color: white;
  stroke: white;
}

.employee-info {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.employee-name {
  font-size: 14px;
  font-weight: 600;
  color: var(--ink);
}

.employee-phone {
  font-size: 12px;
  color: var(--ink-soft);
}

.whatsapp-link {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: linear-gradient(135deg, #25D366 0%, #128C7E 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s ease;
  text-decoration: none;
  flex-shrink: 0;
  box-shadow: 0 2px 8px rgba(37, 211, 102, 0.3);
}

.whatsapp-link:hover {
  transform: scale(1.1);
  box-shadow: 0 4px 12px rgba(37, 211, 102, 0.4);
}

.whatsapp-link svg {
  width: 20px;
  height: 20px;
  color: white;
}

/* Selection Summary */
.selection-summary {
  margin-top: 12px;
  padding: 12px 16px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.selected-info {
  background: rgba(102, 126, 234, 0.1);
  color: #667eea;
}

.no-selection {
  background: rgba(239, 68, 68, 0.1);
  color: #ef4444;
}

.selected-info svg,
.no-selection svg {
  width: 16px;
  height: 16px;
}

.clear-btn {
  padding: 4px 12px;
  background: rgba(102, 126, 234, 0.2);
  border: none;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 600;
  color: #667eea;
  cursor: pointer;
  transition: all 0.2s ease;
}

.clear-btn:hover {
  background: rgba(102, 126, 234, 0.3);
}

/* Form Actions */
.form-actions {
  display: flex;
  justify-content: flex-end;
  margin-top: 24px;
}

.send-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 24px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  border-radius: 12px;
  font-size: 14px;
  font-weight: 600;
  color: white;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}

.send-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(102, 126, 234, 0.4);
}

.send-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.send-btn svg {
  width: 18px;
  height: 18px;
}

.loading-icon {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

/* Filters Card */
.filter-header {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px 20px;
  border-bottom: 1px solid var(--line);
}

.filter-header svg {
  width: 18px;
  height: 18px;
  color: var(--brand);
}

.filter-header h3 {
  margin: 0;
  font-size: 14px;
  font-weight: 600;
  color: var(--ink);
}

.filter-content {
  padding: 16px 20px;
}

.filter-group {
  margin-bottom: 0;
}

.filter-group label {
  display: block;
  margin-bottom: 8px;
  font-size: 12px;
  font-weight: 600;
  color: var(--ink-soft);
}

/* Notes Grid */
.notes-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 20px;
}

.note-card {
  background: var(--surface);
  border-radius: 16px;
  box-shadow: var(--shadow-sm);
  transition: all 0.3s ease;
  overflow: hidden;
  border: 2px solid transparent;
}

.note-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
}

.note-card.unread {
  border-color: var(--brand);
  background: linear-gradient(135deg, var(--brand-tint) 0%, var(--surface) 100%);
}

.note-card-header {
  padding: 16px 20px;
  border-bottom: 1px solid var(--line);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.note-recipient {
  display: flex;
  align-items: center;
  gap: 12px;
}

.note-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
}

.note-avatar svg {
  width: 18px;
  height: 18px;
  color: white;
  stroke: white;
}

.note-recipient-info {
  display: flex;
  flex-direction: column;
}

.note-recipient-name {
  font-size: 14px;
  font-weight: 600;
  color: var(--ink);
}

.note-recipient-phone {
  font-size: 12px;
  color: var(--ink-soft);
}

.note-actions {
  display: flex;
  gap: 8px;
}

.action-btn {
  width: 36px;
  height: 36px;
  border-radius: 8px;
  border: none;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s ease;
}

.action-btn svg {
  width: 18px;
  height: 18px;
}

.edit-btn {
  background: rgba(102, 126, 234, 0.1);
  color: #667eea;
}

.edit-btn:hover {
  background: rgba(102, 126, 234, 0.2);
}

.delete-btn {
  background: rgba(239, 68, 68, 0.1);
  color: #ef4444;
}

.delete-btn:hover {
  background: rgba(239, 68, 68, 0.2);
}

.note-card-body {
  padding: 16px 20px;
}

.note-content {
  font-size: 14px;
  line-height: 1.6;
  color: var(--ink);
  white-space: pre-wrap;
  margin-bottom: 12px;
}

.note-worksite {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  background: rgba(59, 130, 246, 0.1);
  border-radius: 8px;
  font-size: 12px;
  color: #3b82f6;
  font-weight: 600;
}

.note-worksite svg {
  width: 14px;
  height: 14px;
}

.note-card-footer {
  padding: 12px 20px;
  border-top: 1px solid var(--line);
  display: flex;
  align-items: center;
  gap: 16px;
  font-size: 12px;
  color: var(--ink-soft);
}

.note-date,
.note-updated,
.note-read-status {
  display: flex;
  align-items: center;
  gap: 4px;
}

.note-date svg,
.note-updated svg,
.note-read-status svg {
  width: 14px;
  height: 14px;
}

.note-updated {
  color: var(--brand);
}

.note-read-status {
  margin-left: auto;
  padding: 4px 10px;
  border-radius: 6px;
  background: var(--canvas);
  font-weight: 600;
}

.note-read-status.read {
  background: rgba(16, 185, 129, 0.1);
  color: #10b981;
}

/* Loading State */
.loading-state {
  text-align: center;
  padding: 60px 20px;
}

.loading-spinner {
  width: 48px;
  height: 48px;
  border: 4px solid var(--line);
  border-top-color: var(--brand);
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto 16px;
}

.loading-state p {
  color: var(--ink-soft);
  font-size: 14px;
}

/* Empty State */
.empty-state {
  text-align: center;
  padding: 60px 20px;
}

.empty-icon {
  width: 80px;
  height: 80px;
  margin: 0 auto 20px;
  background: var(--brand-tint);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.empty-icon svg {
  width: 40px;
  height: 40px;
  color: var(--brand);
}

.empty-state h3 {
  margin: 0 0 8px;
  font-size: 18px;
  color: var(--ink);
}

.empty-state p {
  margin: 0;
  color: var(--ink-soft);
  font-size: 14px;
}

/* Modal */
.modal-backdrop {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 20px;
  backdrop-filter: blur(4px);
}

.modal-card {
  background: var(--surface);
  border-radius: 20px;
  max-width: 500px;
  width: 100%;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: var(--shadow-xl);
  animation: modalSlideIn 0.3s ease-out;
}

@keyframes modalSlideIn {
  from {
    opacity: 0;
    transform: translateY(-20px) scale(0.95);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

.modal-header {
  padding: 20px 24px;
  border-bottom: 1px solid var(--line);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.modal-title {
  display: flex;
  align-items: center;
  gap: 12px;
}

.modal-title svg {
  width: 20px;
  height: 20px;
  color: var(--brand);
}

.modal-title h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: var(--ink);
}

.modal-close {
  width: 36px;
  height: 36px;
  border-radius: 8px;
  border: none;
  background: var(--canvas);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
}

.modal-close:hover {
  background: var(--line);
}

.modal-close svg {
  width: 20px;
  height: 20px;
  color: var(--ink-soft);
}

.modal-body {
  padding: 24px;
}

.modal-footer {
  padding: 16px 24px;
  border-top: 1px solid var(--line);
  display: flex;
  gap: 12px;
  justify-content: flex-end;
}

/* Success Modal */
.success-modal-card {
  background: var(--surface);
  border-radius: 20px;
  max-width: 400px;
  width: 100%;
  text-align: center;
  padding: 40px 32px;
  box-shadow: var(--shadow-xl);
  animation: modalSlideIn 0.3s ease-out;
}

.success-animation {
  margin-bottom: 24px;
  display: flex;
  justify-content: center;
}

.success-circle {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  background: linear-gradient(135deg, #10b981 0%, #059669 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  animation: successPulse 0.5s ease-out;
  box-shadow: 0 10px 30px rgba(16, 185, 129, 0.3);
}

@keyframes successPulse {
  0% {
    transform: scale(0);
  }
  50% {
    transform: scale(1.1);
  }
  100% {
    transform: scale(1);
  }
}

.success-check {
  width: 40px;
  height: 40px;
  color: white;
  animation: checkDraw 0.3s ease-out 0.2s forwards;
  opacity: 0;
}

@keyframes checkDraw {
  from {
    opacity: 0;
    transform: scale(0.5);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}

.success-title {
  margin: 0 0 12px;
  font-size: 24px;
  font-weight: 700;
  color: var(--ink);
}

.success-message {
  margin: 0 0 24px;
  font-size: 15px;
  color: var(--ink-soft);
  line-height: 1.6;
}

.success-btn {
  padding: 12px 32px;
  background: linear-gradient(135deg, #10b981 0%, #059669 100%);
  border: none;
  border-radius: 12px;
  font-size: 14px;
  font-weight: 600;
  color: white;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 4px 12px rgba(16, 185, 129, 0.3);
}

.success-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(16, 185, 129, 0.4);
}

/* Responsive Design */
@media (max-width: 768px) {
  .notes-page {
    padding: 16px;
  }

  .header-content {
    flex-direction: column;
    text-align: center;
  }

  .header-icon {
    width: 56px;
    height: 56px;
  }

  .header-icon svg {
    width: 28px;
    height: 28px;
  }

  .header-text h1 {
    font-size: 24px;
  }

  .stats-grid {
    grid-template-columns: 1fr;
  }

  .notes-grid {
    grid-template-columns: 1fr;
  }

  .note-card-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }

  .note-actions {
    width: 100%;
    justify-content: flex-end;
  }

  .note-card-footer {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }

  .note-read-status {
    margin-left: 0;
  }
}

/* ==========================================================================
   الوضع الليلي - إصلاحات إضافية لقسم الملاحظات
   ========================================================================== */
[data-theme="dark"] .notes-page {
  background: var(--canvas);
}

[data-theme="dark"] .header-icon {
  background: linear-gradient(135deg, var(--brand) 0%, var(--brand-dark) 100%);
  box-shadow: var(--shadow-md), var(--brand-glow);
}

[data-theme="dark"] .stat-card.total .stat-icon {
  background: linear-gradient(135deg, var(--brand) 0%, var(--brand-dark) 100%);
}

[data-theme="dark"] .stat-card.unread .stat-icon {
  background: linear-gradient(135deg, var(--signal-out) 0%, #dc2626 100%);
}

[data-theme="dark"] .stat-card.read .stat-icon {
  background: linear-gradient(135deg, var(--signal-in) 0%, #059669 100%);
}

[data-theme="dark"] .employee-avatar {
  background: linear-gradient(135deg, var(--brand) 0%, var(--brand-dark) 100%);
}

[data-theme="dark"] .note-avatar {
  background: linear-gradient(135deg, var(--brand) 0%, var(--brand-dark) 100%);
}

[data-theme="dark"] .edit-btn {
  background: rgba(59, 130, 246, 0.2);
  color: var(--brand-light);
}

[data-theme="dark"] .edit-btn:hover {
  background: rgba(59, 130, 246, 0.3);
}

[data-theme="dark"] .delete-btn {
  background: rgba(239, 68, 68, 0.2);
  color: var(--signal-out);
}

[data-theme="dark"] .delete-btn:hover {
  background: rgba(239, 68, 68, 0.3);
}

[data-theme="dark"] .note-worksite {
  background: rgba(59, 130, 246, 0.2);
  color: var(--brand-light);
}

[data-theme="dark"] .selected-info {
  background: rgba(59, 130, 246, 0.2);
  color: var(--brand-light);
}

[data-theme="dark"] .no-selection {
  background: rgba(239, 68, 68, 0.2);
  color: var(--signal-out);
}

/* Delete Modal Styles */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  animation: fadeIn 0.2s ease;
}

.delete-modal-card {
  background: var(--surface);
  border-radius: var(--radius-lg);
  padding: 32px;
  max-width: 400px;
  width: 90%;
  box-shadow: var(--shadow-xl);
  animation: scaleIn 0.2s ease;
}

.delete-icon-wrapper {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 64px;
  height: 64px;
  margin: 0 auto 20px;
  border-radius: 50%;
  background: var(--signal-out-tint);
}

.delete-icon {
  width: 32px;
  height: 32px;
  color: var(--signal-out);
}

.delete-title {
  font-size: 20px;
  font-weight: 700;
  color: var(--ink);
  text-align: center;
  margin-bottom: 12px;
}

.delete-message {
  font-size: 14px;
  color: var(--ink-soft);
  text-align: center;
  margin-bottom: 24px;
  line-height: 1.6;
}

.delete-actions {
  display: flex;
  gap: 12px;
  justify-content: center;
}

.cancel-btn {
  flex: 1;
}

.delete-confirm-btn {
  flex: 1;
}

[data-theme="dark"] .delete-icon-wrapper {
  background: rgba(239, 68, 68, 0.2);
}

[data-theme="dark"] .delete-icon {
  color: var(--signal-out);
}

[data-theme="dark"] .delete-title {
  color: var(--ink);
}

[data-theme="dark"] .delete-message {
  color: var(--ink-soft);
}

[data-theme="dark"] .clear-btn:hover {
  background: rgba(59, 130, 246, 0.4);
}

[data-theme="dark"] .send-btn {
  background: linear-gradient(135deg, var(--brand) 0%, var(--brand-dark) 100%);
  box-shadow: var(--shadow-md), var(--brand-glow);
}

[data-theme="dark"] .send-btn:hover:not(:disabled) {
  box-shadow: var(--shadow-lg), var(--shadow-glow);
}

[data-theme="dark"] .success-btn {
  background: linear-gradient(135deg, var(--signal-in) 0%, #059669 100%);
  box-shadow: var(--shadow-md);
}

[data-theme="dark"] .success-btn:hover {
  box-shadow: var(--shadow-lg);
}

[data-theme="dark"] .success-circle {
  background: linear-gradient(135deg, var(--signal-in) 0%, #059669 100%);
  box-shadow: var(--shadow-md);
}

.clear-btn {
  background: rgba(59, 130, 246, 0.1);
  color: var(--brand);
}

.clear-btn:hover {
  background: rgba(59, 130, 246, 0.2);
}

[data-theme="dark"] .clear-btn {
  background: rgba(59, 130, 246, 0.2);
  color: var(--brand-light);
}

[data-theme="dark"] .clear-btn:hover {
  background: rgba(59, 130, 246, 0.4);
}
</style>
