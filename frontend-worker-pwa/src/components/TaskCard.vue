<template>
  <router-link :to="`/tasks/${task.id}`" class="task-card">
    <div class="task-card__stripe" :class="task.status"></div>
    <div class="task-card__body">
      <h3>{{ task.title }}</h3>
      <p class="task-card__meta">{{ task.worksite }}</p>
      <p class="task-card__time mono">{{ task.time }}</p>
    </div>
    <span class="badge" :class="badgeClass">{{ statusLabels[task.status] }}</span>
  </router-link>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({ task: { type: Object, required: true } })

const statusLabels = {
  pending: 'قيد الانتظار',
  in_progress: 'جارية',
  completed: 'مكتملة',
  late: 'متأخرة',
}
const badgeMap = { pending: '', in_progress: 'badge--gold', completed: 'badge--in', late: 'badge--out' }
const badgeClass = computed(() => badgeMap[props.task.status] || '')
</script>

<style scoped>
.task-card {
  display: flex; align-items: center; gap: 14px; background: var(--surface); border: 1px solid var(--line);
  border-radius: var(--radius-md); padding: 14px 16px; box-shadow: var(--shadow-sm);
  transition: transform .15s ease, box-shadow .15s ease; margin-bottom: 10px;
}
.task-card:hover { transform: translateY(-2px); box-shadow: var(--shadow-md); }
.task-card__stripe { width: 5px; align-self: stretch; border-radius: 4px; background: var(--line-strong); }
.task-card__stripe.in_progress { background: var(--gold); }
.task-card__stripe.completed { background: var(--signal-in); }
.task-card__stripe.late { background: var(--signal-out); }
.task-card__body { flex: 1; min-width: 0; }
.task-card__body h3 { font-size: 15px; margin-bottom: 3px; }
.task-card__meta { font-size: 13px; color: var(--ink-soft); margin-bottom: 3px; }
.task-card__time { font-size: 12px; color: var(--ink-soft); }
</style>
