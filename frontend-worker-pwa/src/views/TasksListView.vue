<template>
  <div class="tasks-view view">
    <!-- Page Header -->
    <div class="page-header page-header--flex">
      <div style="display:flex;align-items:center;gap:14px;">
        <div class="header-icon">
          <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M9 11l3 3L22 4"/>
            <path d="M21 12v7a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11"/>
          </svg>
        </div>
        <div>
          <h1>{{ t('my_tasks') }}</h1>
          <p>{{ t('tasks_subtitle') }}</p>
        </div>
      </div>
      <span class="badge badge--info">{{ tasksStore.items.length }}</span>
    </div>

    <!-- Loading State -->
    <div v-if="tasksStore.loading" class="loading">{{ t('loading') }}</div>

    <!-- Empty State -->
    <div v-else-if="!tasksStore.items.length" class="empty">
      <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
        <path d="M9 11l3 3L22 4"/>
        <path d="M21 12v7a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11"/>
      </svg>
      <p>{{ t('no_tasks_currently') }}</p>
    </div>

    <!-- Tasks List -->
    <div v-else class="tasks-list stagger">
      <div 
        v-for="task in tasksStore.items" 
        :key="task.id" 
        class="task-card"
        @click="router.push(`/tasks/${task.id}`)"
      >
        <div class="task-stripe" :class="getTaskStatusClass(task.status)"></div>
        <div class="task-body">
          <div class="task-top">
            <span class="task-title">{{ task.title }}</span>
            <span class="badge" :class="getTaskBadgeClass(task.status)">
              {{ getTaskStatusLabel(task.status) }}
            </span>
          </div>
          <div class="task-meta">
            <span class="task-meta-item">
              <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M20 10c0 6-8 12-8 12s-8-6-8-12a8 8 0 0 1 16 0Z"/>
              </svg>
              {{ task.worksite_name || t('general') }}
            </span>
            <span class="task-meta-item">
              <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="10"/>
                <polyline points="12 6 12 12 16 14"/>
              </svg>
              {{ formatDate(task.due_date) }}
            </span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from '../services/i18n'
import { tasksStore } from '../store/tasks'
import wsService from '../services/websocket'

const { t } = useI18n()
const router = useRouter()

onMounted(() => {
  tasksStore.fetchMine()
  connectWebSocket()
})

onUnmounted(() => {
  disconnectWebSocket()
})

function connectWebSocket() {
  const apiBaseUrl = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1'
  const apiHost = apiBaseUrl.replace('/api/v1', '')
  const wsUrl = apiHost.replace('http://', 'ws://').replace('https://', 'wss://') + '/ws'
  wsService.connect(wsUrl)
  
  wsService.onMessage((data) => {
    if (data.type === 'task_update') {
      tasksStore.fetchMine()
    }
  })
}

function disconnectWebSocket() {
  wsService.disconnect()
}

function getTaskStatusClass(status) {
  const statusMap = {
    'pending': 'stripe--pending',
    'in_progress': 'stripe--progress',
    'completed': 'stripe--done',
    'late': 'stripe--late'
  }
  return statusMap[status] || 'stripe--pending'
}

function getTaskBadgeClass(status) {
  const badgeMap = {
    'pending': 'badge--neutral',
    'in_progress': 'badge--warning',
    'completed': 'badge--success',
    'late': 'badge--error'
  }
  return badgeMap[status] || 'badge--neutral'
}

function getTaskStatusLabel(status) {
  const labelMap = {
    'pending': t('pending'),
    'in_progress': t('in_progress'),
    'completed': t('completed'),
    'late': t('late')
  }
  return labelMap[status] || status
}

function formatDate(dateStr) {
  if (!dateStr) return '—'
  const date = new Date(dateStr)
  return date.toLocaleDateString('ar-SA', { month: 'short', day: 'numeric' })
}
</script>

<style scoped>
.tasks-view {
  padding: var(--space-4);
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

.tasks-list {
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
}

.task-card {
  display: flex;
  gap: 0;
  border-radius: var(--radius-lg);
  overflow: hidden;
  background: var(--surface);
  border: 1px solid var(--border);
  cursor: pointer;
  transition: var(--transition-base) ease;
}

.task-card:hover {
  box-shadow: var(--shadow-sm);
  border-color: var(--border-strong);
}

.task-stripe {
  width: 4px;
  flex-shrink: 0;
}

.stripe--pending {
  background: var(--gray-400);
}

.stripe--progress {
  background: var(--warning-500);
}

.stripe--done {
  background: var(--success-500);
}

.stripe--late {
  background: var(--error-500);
}

.task-body {
  padding: var(--space-4);
  flex: 1;
}

.task-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--space-2);
  gap: var(--space-2);
}

.task-title {
  font-weight: var(--font-medium);
  font-size: var(--text-sm);
  color: var(--text-primary);
}

.task-meta {
  display: flex;
  align-items: center;
  gap: var(--space-4);
  font-size: var(--text-xs);
  color: var(--text-tertiary);
}

.task-meta-item {
  display: flex;
  align-items: center;
  gap: 4px;
}

@media (min-width: 768px) {
  .tasks-list {
    gap: var(--space-4);
  }
  
  .task-body {
    padding: var(--space-5);
  }
}
</style>
