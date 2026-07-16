<template>
  <div>
    <div class="page-head">
      <div><h2>{{ t('tasks') }}</h2><p>{{ t('tasks_description') }}</p></div>
      <button class="btn btn--primary">{{ t('new_task') }}</button>
    </div>

    <div class="card">
      <!-- جدول للشاشات الكبيرة -->
      <div class="table-wrapper desktop-only">
        <table class="table">
          <thead><tr><th>{{ t('task_title') }}</th><th>{{ t('task_employee') }}</th><th>{{ t('task_worksite') }}</th><th>{{ t('task_status') }}</th></tr></thead>
          <tbody>
            <tr v-for="t in tasks" :key="t.id">
              <td>{{ t.title }}</td>
              <td>{{ t.employee }}</td>
              <td>{{ t.worksite }}</td>
              <td><span class="badge" :class="badgeMap[t.status]">{{ labels[t.status] }}</span></td>
            </tr>
          </tbody>
        </table>
      </div>
      
      <!-- بطاقات للشاشات الصغيرة -->
      <div class="mobile-cards mobile-only">
        <div v-for="t in tasks" :key="t.id" class="task-card">
          <div class="task-card__header">
            <span class="task-card__title">{{ t.title }}</span>
            <span class="badge" :class="badgeMap[t.status]">{{ labels[t.status] }}</span>
          </div>
          <div class="task-card__body">
            <div class="task-card__row">
              <span class="task-card__label">{{ t('task_employee') }}</span>
              <span class="task-card__value">{{ t.employee }}</span>
            </div>
            <div class="task-card__row">
              <span class="task-card__label">{{ t('task_worksite') }}</span>
              <span class="task-card__value">{{ t.worksite }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useI18n } from '../services/i18n'

const { t } = useI18n()

const labels = computed(() => ({ 
  pending: t('status_pending'), 
  in_progress: t('status_in_progress'), 
  completed: t('status_completed'), 
  late: t('status_late') 
}))
const badgeMap = { pending: '', in_progress: 'badge--gold', completed: 'badge--in', late: 'badge--out' }

const tasks = ref([
  { id: 1, title: 'صيانة مكيّفات', employee: 'أحمد ياسين', worksite: 'برج الأمل', status: 'in_progress' },
  { id: 2, title: 'فحص دوري', employee: 'سارة قدورة', worksite: 'الشميساني', status: 'pending' },
  { id: 3, title: 'تركيب كاميرات', employee: 'ليث عودة', worksite: 'طريق المطار', status: 'completed' },
])
</script>

<style scoped>
.page-head { display: flex; align-items: flex-start; justify-content: space-between; margin-bottom: 18px; flex-wrap: wrap; gap: 10px; }
.page-head h2 { font-size: 20px; margin-bottom: 4px; }
.page-head p { font-size: 13px; color: var(--ink-soft); }
.table { width: 100%; border-collapse: collapse; }
.table th { text-align: right; font-size: 12px; color: var(--ink-soft); font-weight: 600; padding: 14px 20px; border-bottom: 1px solid var(--line); }
.table td { padding: 14px 20px; font-size: 14px; border-bottom: 1px solid var(--line); }
.table tr:last-child td { border-bottom: none; }

@media (max-width: 768px) {
  .page-head { flex-direction: column; align-items: flex-start; }
  .page-head .btn { width: 100%; }
  
  .desktop-only { display: none; }
  .mobile-only { display: block; }
}

@media (min-width: 769px) {
  .desktop-only { display: block; }
  .mobile-only { display: none; }
}

/* تصميم بطاقات المهام للهاتف */
.mobile-cards {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.task-card {
  background: var(--surface);
  border: 1px solid var(--line);
  border-radius: var(--radius-md);
  padding: 16px;
  transition: all 0.2s ease;
}

.task-card:hover {
  border-color: var(--brand);
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

.task-card__header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  padding-bottom: 12px;
  border-bottom: 1px solid var(--line);
}

.task-card__title {
  font-weight: 600;
  color: var(--ink);
  font-size: 15px;
}

.task-card__body {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.task-card__row {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.task-card__label {
  font-size: 13px;
  color: var(--ink-soft);
  font-weight: 500;
}

.task-card__value {
  font-size: 14px;
  color: var(--ink);
  font-weight: 500;
}
</style>
