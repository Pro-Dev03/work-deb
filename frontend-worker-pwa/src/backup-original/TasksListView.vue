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
      <SwipeNav :items-per-view="1" :show-indicators="true">
        <TaskCard v-for="task in tasksStore.items" :key="task.id" :task="task" class="swipe-item" />
      </SwipeNav>
    </template>
  </div>
</template>

<script setup>
import { onMounted, onUnmounted } from 'vue'
import { useI18n } from '../services/i18n'
import { tasksStore } from '../store/tasks'
import TaskCard from '../components/TaskCard.vue'
import SwipeNav from '../components/SwipeNav.vue'
import wsService from '../services/websocket'

const { t } = useI18n()

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
  console.log('🔌 Attempting to connect to WebSocket:', wsUrl)
  wsService.connect(wsUrl)
  
  wsService.onMessage((data) => {
    if (data.type === 'connected') {
      console.log('✅ تم الاتصال بـ WebSocket')
    } else if (data.type === 'disconnected') {
      console.log('❌ انقطع الاتصال بـ WebSocket')
    } else if (data.type === 'task_update') {
      console.log('📋 تحديث مهمة:', data.data)
      tasksStore.fetchMine()
    }
  })
}

function disconnectWebSocket() {
  wsService.disconnect()
}
</script>

<style scoped>
.page-head { display: flex; align-items: center; justify-content: space-between; margin-bottom: 16px; }

.swipe-item {
  flex: 0 0 100%;
  width: 100%;
  padding: 0 8px;
  box-sizing: border-box;
}

.page-head h2 { font-size: 19px; }
</style>
