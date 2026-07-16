<template>
  <div>
    <div class="page-head">
      <div><h2>{{ t('tasks') }}</h2><p>{{ t('tasks_description') }}</p></div>
      <button class="btn btn--primary">{{ t('new_task') }}</button>
    </div>

    <div class="card">
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
.page-head { display: flex; align-items: flex-start; justify-content: space-between; margin-bottom: 18px; }
.page-head h2 { font-size: 20px; margin-bottom: 4px; }
.page-head p { font-size: 13px; color: var(--ink-soft); }
.table { width: 100%; border-collapse: collapse; }
.table th { text-align: right; font-size: 12px; color: var(--ink-soft); font-weight: 600; padding: 14px 20px; border-bottom: 1px solid var(--line); }
.table td { padding: 14px 20px; font-size: 14px; border-bottom: 1px solid var(--line); }
.table tr:last-child td { border-bottom: none; }
</style>
