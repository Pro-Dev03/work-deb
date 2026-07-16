<template>
  <div v-if="task">
    <router-link to="/tasks" class="back-link">{{ t('back_to_tasks') }}</router-link>

    <div class="card task-detail">
      <span class="badge" :class="badgeClass">{{ statusLabels[task.status] }}</span>
      <h2>{{ task.title }}</h2>
      <p class="task-detail__meta">📍 {{ task.worksite }}</p>
      <p class="task-detail__meta mono">🕒 {{ task.time }}</p>
    </div>

    <router-link to="/attendance" class="btn btn--primary btn--block attendance-cta">
      {{ t('register_attendance_for_task') }}
    </router-link>
  </div>
  <div v-else class="empty-state">
    <h3>{{ t('task_not_found') }}</h3>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from '../services/i18n'
import { tasksStore } from '../store/tasks'

const { t } = useI18n()
const route = useRoute()
const task = computed(() => tasksStore.find(route.params.id))

const statusLabels = computed(() => ({
  pending: t('status_pending'),
  in_progress: t('status_in_progress'),
  completed: t('status_completed'),
  late: t('status_late')
}))

const badgeMap = {
  pending: '',
  in_progress: 'badge--gold',
  completed: 'badge--in',
  late: 'badge--out'
}

const badgeClass = computed(() => (task.value ? badgeMap[task.value.status] : ''))
</script>

<style scoped>
.back-link {
  display: inline-block;
  font-size: 13px;
  color: var(--ink-soft);
  margin-bottom: 14px;
}

.task-detail {
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin-bottom: 14px;
}

.task-detail h2 {
  font-size: 18px;
}

.task-detail__meta {
  font-size: 14px;
  color: var(--ink-soft);
}

.attendance-cta {
  position: sticky;
  bottom: 96px;
}
</style>
