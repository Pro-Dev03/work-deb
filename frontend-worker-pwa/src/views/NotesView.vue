<template>
  <div class="notes-view view">
    <!-- Page Header -->
    <div class="page-header">
      <div class="header-icon">
        <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/>
        </svg>
      </div>
      <div class="header-content">
        <div class="title-row">
          <h1>{{ t('notes') }}</h1>
          <button 
            class="whatsapp-float-btn" 
            @click="openWhatsApp"
            title="تواصل عبر واتساب"
          >
            <svg class="whatsapp-icon" viewBox="0 0 32 32" fill="currentColor">
              <path d="M16 2C8.268 2 2 8.268 2 16c0 2.52.666 4.93 1.84 7.094L2.5 29.5l6.562-1.312A13.94 13.94 0 0 0 16 30c7.732 0 14-6.268 14-14S23.732 2 16 2zm0 25.5c-2.234 0-4.39-.586-6.297-1.688l-.453-.266-4.75.95.996-4.625-.281-.469A11.38 11.38 0 0 1 4.5 16c0-6.344 5.156-11.5 11.5-11.5S27.5 9.656 27.5 16 22.344 27.5 16 27.5zm6.344-8.656c-.344-.172-2.031-1-2.344-1.109-.313-.109-.531-.172-.75.172-.219.344-.844 1.109-1.031 1.344-.188.234-.375.266-.719.094-.344-.172-1.453-.531-2.766-1.703-1.031-.906-1.719-2.031-1.906-2.375-.188-.344-.016-.531.156-.688.156-.156.344-.406.516-.609.172-.203.234-.344.344-.578.109-.234.055-.438-.027-.609-.082-.172-.75-1.813-1.031-2.484-.266-.641-.547-.555-.75-.563-.188-.008-.406-.008-.625-.008-.219 0-.578.082-.875.406-.297.328-1.125 1.094-1.125 2.672 0 1.578 1.148 3.094 1.313 3.328.164.234 2.266 3.453 5.484 4.844 2.156.922 2.594.734 3.063.688.469-.047 1.5-.609 1.703-1.203.203-.594.203-1.109.141-1.234-.063-.125-.234-.203-.578-.375z"/>
            </svg>
          </button>
        </div>
        <p>{{ t('notes_subtitle') }}</p>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="loading">{{ t('loading') }}</div>

    <!-- Empty State -->
    <div v-else-if="notes.length === 0" class="empty">
      <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
        <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/>
      </svg>
      <p>{{ t('no_notes') }}</p>
    </div>

    <!-- Notes List -->
    <div v-else class="notes-list stagger">
      <div 
        v-for="note in notes" 
        :key="note.id" 
        class="note-card"
        :class="{ unread: !note.is_read }"
      >
        <div class="note-top">
          <div class="note-sender">
            <div class="avatar avatar--40 admin-avatar">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M16 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
                <circle cx="8.5" cy="7" r="4"/>
                <polyline points="17 11 19 13 23 13"/>
                <path d="M23 16a4 4 0 0 0-4-4"/>
              </svg>
            </div>
            <div class="note-sender-info">
              <div class="note-name">{{ note.admin_name }}</div>
              <div class="note-role">{{ t('admin') }}</div>
            </div>
          </div>
          <div v-if="!note.is_read" class="unread-dot"></div>
        </div>

        <div class="note-content">{{ note.content }}</div>

        <div class="note-foot">
          <div class="note-meta">
            <span v-if="note.worksite_name" class="badge badge--neutral">{{ note.worksite_name }}</span>
            <span class="note-date">{{ formatDate(note.created_at) }}</span>
          </div>
          <button 
            v-if="!note.is_read" 
            class="btn btn--ghost btn--sm" 
            @click="markAsRead(note.id)"
          >
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
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
  } catch(e) {
    console.error('❌ فشل جلب الملاحظات:', e)
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
  } catch(error) {
    console.error('❌ فشل تحديث حالة القراءة:', error)
  }
}

async function markAllAsRead() {
  try {
    const unreadNotes = notes.value.filter(n => !n.is_read)
    if (unreadNotes.length > 0) {
      for (const note of unreadNotes) {
        await api.put(`/notes/${note.id}/read`)
        note.is_read = true
      }
    }
  } catch(error) {
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

function openWhatsApp() {
  const phoneNumber = '0584838136'
  const cleanPhone = phoneNumber.replace(/\D/g, '')
  const whatsappUrl = `https://wa.me/${cleanPhone}`
  window.open(whatsappUrl, '_blank')
}

onMounted(() => {
  fetchNotes()
})
</script>

<style scoped>
.notes-view {
  padding: var(--space-4);
}

.page-header {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  margin-bottom: var(--space-5);
}

.header-icon {
  width: 40px;
  height: 40px;
  border-radius: var(--radius-lg);
  background: var(--primary-100);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--primary-600);
  flex-shrink: 0;
}

.header-content {
  flex: 1;
}

.title-row {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  width: 100%;
}

.header-content h1 {
  font-size: var(--text-xl);
  font-weight: var(--font-semibold);
  color: var(--text-primary);
  margin: 0;
}

.header-content p {
  font-size: var(--text-sm);
  color: var(--text-secondary);
  margin: 2px 0 0 0;
}

.whatsapp-float-btn {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  background: linear-gradient(135deg, #25D366 0%, #128C7E 100%);
  border: none;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 20px rgba(37, 211, 102, 0.4);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  animation: float-pulse 2s ease-in-out infinite;
  margin-left: auto;
}

.whatsapp-float-btn:hover {
  transform: scale(1.1) translateY(-4px);
  box-shadow: 0 8px 30px rgba(37, 211, 102, 0.6);
}

.whatsapp-float-btn:active {
  transform: scale(0.95);
  box-shadow: 0 4px 15px rgba(37, 211, 102, 0.4);
}

.whatsapp-icon {
  width: 32px;
  height: 32px;
  color: white;
  transition: transform 0.3s ease;
}

.whatsapp-float-btn:hover .whatsapp-icon {
  transform: scale(1.1);
}

/* أنيميشن النبض */
@keyframes float-pulse {
  0%, 100% {
    transform: translateY(0);
    box-shadow: 0 4px 20px rgba(37, 211, 102, 0.4);
  }
  50% {
    transform: translateY(-8px);
    box-shadow: 0 8px 30px rgba(37, 211, 102, 0.6);
  }
}

/* تأثير الرنين */
.whatsapp-float-btn::before {
  content: '';
  position: absolute;
  inset: -4px;
  border-radius: 50%;
  border: 2px solid #25D366;
  opacity: 0;
  animation: ripple 2s ease-out infinite;
}

@keyframes ripple {
  0% {
    transform: scale(1);
    opacity: 0.6;
  }
  100% {
    transform: scale(1.5);
    opacity: 0;
  }
}

/* RTL support */
[dir="rtl"] .whatsapp-float-btn {
  margin-left: 0;
  margin-right: auto;
}

.loading, .empty {
  text-align: center;
  padding: var(--space-8);
  color: var(--text-tertiary);
}

.empty svg {
  margin: 0 auto var(--space-4);
  color: var(--text-tertiary);
}

.empty p {
  font-size: var(--text-sm);
}

.notes-list {
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
}

.note-card {
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
  padding: var(--space-4);
  border-radius: var(--radius-lg);
  background: var(--surface);
  border: 1px solid var(--border);
  transition: var(--transition-base) ease;
}

.note-card:hover {
  box-shadow: var(--shadow-sm);
  border-color: var(--border-strong);
}

.note-card.unread {
  border-inline-start: 3px solid var(--primary-500);
  background: linear-gradient(90deg, var(--primary-50) 0%, var(--surface) 15%);
}

[data-theme="dark"] .note-card.unread {
  background: linear-gradient(90deg, rgba(99, 102, 241, 0.12) 0%, var(--surface) 15%);
}

.note-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.note-sender {
  display: flex;
  align-items: center;
  gap: var(--space-3);
}

.note-sender-info {
  display: flex;
  flex-direction: column;
}

.note-name {
  font-weight: var(--font-medium);
  font-size: var(--text-sm);
  color: var(--text-primary);
}

.note-role {
  font-size: var(--text-xs);
  color: var(--text-tertiary);
}

.unread-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--primary-500);
  flex-shrink: 0;
}

.note-content {
  font-size: var(--text-sm);
  color: var(--text-secondary);
  line-height: var(--leading-normal);
}

.note-foot {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--space-3);
}

.note-meta {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  flex-wrap: wrap;
}

.note-date {
  font-size: var(--text-xs);
  color: var(--text-tertiary);
}

.admin-avatar {
  background: linear-gradient(135deg, var(--primary-500) 0%, var(--primary-600) 100%);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
}

.admin-avatar svg {
  color: white;
  stroke: white;
}

@media (min-width: 768px) {
  .notes-list {
    gap: var(--space-4);
  }
  
  .note-card {
    padding: var(--space-5);
  }
}
</style>
