<template>
  <div>
    <div class="page-head">
      <h2>{{ t('my_tasks_today') }}</h2>
      <span class="badge">{{ tasksStore.items.length }} {{ tasksStore.items.length === 1 ? t('task') : t('task_plural') }}</span>
    </div>

    <div v-if="tasksStore.loading" class="empty-state"><p>{{ t('loading_tasks') }}</p></div>

    <div v-else-if="!tasksStore.items.length" class="empty-state">
      <h3>{{ t('no_tasks_currently') }}</h3>
      <p>{{ t('will_notify_new_task') }}</p>
    </div>

    <template v-else>
      <TaskCard v-for="task in tasksStore.items" :key="task.id" :task="task" />
    </template>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useI18n } from '../services/i18n'
import { tasksStore } from '../store/tasks'
import TaskCard from '../components/TaskCard.vue'

const { t } = useI18n()

onMounted(() => tasksStore.fetchMine())
</script>

<style scoped>
.page-head { display: flex; align-items: center; justify-content: space-between; margin-bottom: 16px; }
.page-head h2 { font-size: 19px; }
</style>
