<template>
  <div>
    <h2 class="page-title">📝 {{ t('notes') }}</h2>

    <!-- قائمة الملاحظات -->
    <div v-if="loading" class="card">
      <div class="loading">{{ t('loading') }}</div>
    </div>

    <div v-else-if="notes.length === 0" class="card">
      <div class="empty">
        <p>📭 {{ t('no_notes') }}</p>
      </div>
    </div>

    <div v-else class="notes-list">
      <div 
        v-for="note in notes" 
        :key="note.id" 
        class="note-card"
        :class="{ 'unread': !note.is_read }"
      >
        <div class="note-card-header">
          <div class="note-sender">
            <div class="note-avatar">{{ note.admin_name?.slice(0, 1) || '?' }}</div>
            <div class="note-sender-info">
              <span class="note-sender-name">{{ note.admin_name }}</span>
              <span class="note-sender-role">{{ t('admin') }}</span>
            </div>
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
          </div>
        </div>

        <div class="note-card-body">
          <div class="note-content">{{ note.content }}</div>
          
          <div v-if="note.worksite_name" class="note-worksite-info">
            <div class="worksite-badge">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"/>
                <circle cx="12" cy="10" r="3"/>
              </svg>
              <span class="worksite-name">{{ note.worksite_name }}</span>
            </div>
          </div>
        </div>

        <div class="note-card-footer">
          <div class="note-meta">
            <div class="note-date">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="3" y="4" width="18" height="18" rx="2" ry="2"/>
                <line x1="16" y1="2" x2="16" y2="6"/>
                <line x1="8" y1="2" x2="8" y2="6"/>
                <line x1="3" y1="10" x2="21" y2="10"/>
              </svg>
              <span>{{ formatDate(note.created_at) }}</span>
            </div>
            <div v-if="note.updated_at !== note.created_at" class="note-updated">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
                <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
              </svg>
              <span>{{ t('updated') }}: {{ formatDate(note.updated_at) }}</span>
            </div>
          </div>
          <button 
            v-if="!note.is_read" 
            class="btn btn--primary btn--sm mark-read-btn" 
            @click="markAsRead(note.id)"
          >
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/>
              <polyline points="22 4 12 14.01 9 11.01"/>
            </svg>
            {{ t('mark_as_read') }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from '../services/i18n'
import api from '../services/api'

const { t } = useI18n()

const loading = ref(false)
const notes = ref([])

async function fetchNotes() {
  loading.value = true
  try {
    const { data } = await api.get('/notes')
    notes.value = data || []
  } catch (error) {
    console.error('❌ فشل جلب الملاحظات:', error)
  } finally {
    loading.value = false
  }
}

async function markAsRead(noteId) {
  try {
    await api.put(`/notes/${noteId}/read`)
    // تحديث الملاحظة محلياً
    const note = notes.value.find(n => n.id === noteId)
    if (note) {
      note.is_read = true
    }
  } catch (error) {
    console.error('❌ فشل تحديث حالة القراءة:', error)
  }
}

function formatDate(date) {
  if (!date) return ''
  return new Date(date).toLocaleString('ar-SA', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

onMounted(() => {
  fetchNotes()
})
</script>

<style scoped>
.page-title {
  font-size: 24px;
  margin-bottom: 20px;
  font-weight: 700;
  color: var(--brand);
  display: flex;
  align-items: center;
  gap: 8px;
}

/* إصلاح الوضع الليلي للعنوان */
[data-theme="dark"] .page-title {
  color: var(--brand-light);
}

.notes-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.note-card {
  background: var(--surface);
  border: 1.5px solid var(--line);
  border-radius: var(--radius-lg);
  overflow: hidden;
  transition: all 0.3s ease;
  box-shadow: var(--shadow-sm);
}

.note-card:hover {
  border-color: var(--brand-light);
  box-shadow: var(--shadow-md);
  transform: translateY(-2px);
}

.note-card.unread {
  border-left: 4px solid var(--brand);
  background: linear-gradient(135deg, var(--brand-tint) 0%, var(--surface) 100%);
}

.note-card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid var(--line);
  background: linear-gradient(135deg, var(--canvas) 0%, var(--surface) 100%);
}

.note-sender {
  display: flex;
  align-items: center;
  gap: 12px;
}

.note-avatar {
  width: 44px;
  height: 44px;
  border-radius: 50%;
  background: linear-gradient(135deg, var(--brand) 0%, var(--brand-dark) 100%);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  font-weight: 700;
  box-shadow: var(--shadow-sm);
}

.note-sender-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.note-sender-name {
  font-weight: 700;
  font-size: 15px;
  color: var(--ink);
}

.note-sender-role {
  font-size: 12px;
  color: var(--ink-soft);
  font-weight: 500;
}

.note-read-status {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: var(--canvas);
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s ease;
}

.note-read-status svg {
  width: 20px;
  height: 20px;
}

.note-read-status.read {
  background: rgba(16, 185, 129, 0.15);
  color: #10b981;
}

.note-read-status:not(.read) {
  background: rgba(245, 158, 11, 0.15);
  color: #f59e0b;
}

.note-card-body {
  padding: 20px;
}

.note-content {
  font-size: 15px;
  line-height: 1.7;
  color: var(--ink);
  white-space: pre-wrap;
  margin-bottom: 16px;
  font-weight: 400;
}

.note-worksite-info {
  margin-top: 12px;
}

.worksite-badge {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  background: linear-gradient(135deg, var(--brand-tint) 0%, rgba(59, 130, 246, 0.1) 100%);
  border: 1px solid var(--brand);
  border-radius: var(--radius-md);
  font-size: 13px;
  color: var(--brand);
  font-weight: 600;
  transition: all 0.3s ease;
}

.worksite-badge:hover {
  background: var(--brand);
  color: white;
  transform: translateY(-1px);
  box-shadow: var(--shadow-sm);
}

.worksite-badge svg {
  width: 16px;
  height: 16px;
}

.worksite-name {
  font-weight: 600;
}

.note-card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 20px;
  border-top: 1px solid var(--line);
  background: var(--canvas);
}

.note-meta {
  display: flex;
  align-items: center;
  gap: 16px;
  flex-wrap: wrap;
}

.note-date,
.note-updated {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: var(--ink-soft);
  font-weight: 500;
}

.note-date svg,
.note-updated svg {
  width: 14px;
  height: 14px;
}

.note-updated {
  color: var(--brand);
}

.mark-read-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  font-size: 13px;
  font-weight: 600;
  border-radius: var(--radius-sm);
  transition: all 0.3s ease;
}

.mark-read-btn svg {
  width: 16px;
  height: 16px;
}

.mark-read-btn:hover {
  transform: translateY(-1px);
  box-shadow: var(--shadow-sm);
}

/* إصلاح الوضع الليلي */
[data-theme="dark"] .note-sender-name {
  color: var(--ink);
}

[data-theme="dark"] .note-sender-role {
  color: var(--ink-soft);
}

[data-theme="dark"] .note-content {
  color: var(--ink);
}

[data-theme="dark"] .note-updated {
  color: var(--brand-light);
}

[data-theme="dark"] .worksite-badge {
  background: rgba(59, 130, 246, 0.2);
  border-color: var(--brand-light);
  color: var(--brand-light);
}

[data-theme="dark"] .worksite-badge:hover {
  background: var(--brand-light);
  color: var(--canvas);
}

[data-theme="dark"] .note-date,
[data-theme="dark"] .note-updated {
  color: var(--ink-soft);
}

.loading {
  text-align: center;
  padding: 40px;
  color: var(--ink-soft);
}

.empty {
  text-align: center;
  padding: 40px;
  color: var(--ink-soft);
}

/* إصلاح الوضع الليلي */
[data-theme="dark"] .loading {
  color: var(--ink-soft);
}

[data-theme="dark"] .empty {
  color: var(--ink-soft);
}

@media (max-width: 768px) {
  .note-card-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }

  .note-card-footer {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }

  .note-meta {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }

  .mark-read-btn {
    width: 100%;
    justify-content: center;
  }

  .worksite-badge {
    width: 100%;
    justify-content: center;
  }
}
</style>
