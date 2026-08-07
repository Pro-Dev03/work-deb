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
// renamed from the old --gold/--in/--out badge names to the shared
// semantic set defined in base.css, so this matches every other badge
// in the app (notes, attendance, profile)
const badgeMap = { pending: 'badge-neutral', in_progress: 'badge-warning', completed: 'badge-success', late: 'badge-error' }
const badgeClass = computed(() => badgeMap[props.task.status] || 'badge-neutral')
</script>

<style scoped>
.task-card {
  display: flex; align-items: center; gap: var(--space-3);
  background: var(--surface); border: 1px solid var(--border);
  border-radius: var(--radius-lg); padding: var(--space-4);
  box-shadow: var(--shadow-sm);
  transition: transform .15s ease, box-shadow .15s ease, border-color .15s ease;
  margin-bottom: 10px;
}
.task-card:hover { transform: translateY(-2px); box-shadow: var(--shadow-md); border-color: var(--border-strong); }
.task-card__stripe { width: 5px; align-self: stretch; border-radius: 4px; background: var(--gray-400); flex-shrink: 0; }
.task-card__stripe.in_progress { background: var(--warning-500); }
.task-card__stripe.completed { background: var(--success-500); }
.task-card__stripe.late { background: var(--error-500); }
.task-card__body { flex: 1; min-width: 0; }
.task-card__body h3 { font-size: var(--text-sm); font-weight: var(--font-semibold); color: var(--text-primary); margin-bottom: 3px; }
.task-card__meta { font-size: var(--text-xs); color: var(--text-tertiary); margin-bottom: 3px; }
.task-card__time { font-size: var(--text-xs); color: var(--text-tertiary); }
</style>
