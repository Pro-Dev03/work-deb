<template>
  <div>
    <div class="page-head">
      <div><h2>المهام</h2><p>جميع المهام عبر الموظفين ونقاط العمل</p></div>
      <button class="btn btn--primary">+ مهمة جديدة</button>
    </div>

    <div class="card">
      <table class="table">
        <thead><tr><th>المهمة</th><th>الموظف</th><th>نقطة العمل</th><th>الحالة</th></tr></thead>
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
import { ref } from 'vue'

const labels = { pending: 'قيد الانتظار', in_progress: 'جارية', completed: 'مكتملة', late: 'متأخرة' }
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
