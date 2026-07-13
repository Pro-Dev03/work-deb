<template>
  <div>
    <div class="page-head">
      <h2>مهامي اليوم</h2>
      <span class="badge">{{ tasksStore.items.length }} مهمة</span>
    </div>

    <div v-if="tasksStore.loading" class="empty-state"><p>جارٍ تحميل المهام...</p></div>

    <div v-else-if="!tasksStore.items.length" class="empty-state">
      <h3>لا توجد مهام حالياً</h3>
      <p>سيتم إشعارك فور تكليفك بمهمة جديدة</p>
    </div>

    <template v-else>
      <TaskCard v-for="task in tasksStore.items" :key="task.id" :task="task" />
    </template>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { tasksStore } from '../store/tasks'
import TaskCard from '../components/TaskCard.vue'

onMounted(() => tasksStore.fetchMine())
</script>

<style scoped>
.page-head { display: flex; align-items: center; justify-content: space-between; margin-bottom: 16px; }
.page-head h2 { font-size: 19px; }
</style>
