<template>
  <div>
    <h2 class="page-title">سجل الخدمات</h2>
    <p class="page-subtitle">جميع الخدمات المنفذة لديكم مع تقاريرها</p>

    <ReportCard v-for="r in reports" :key="r.id" :report="r" />
  </div>
</template>

<script setup>
import ReportCard from '../components/ReportCard.vue'

const reports = [
  { id: '1', title: 'صيانة مكيّفات دورية', date: '2026-07-08', status: 'مكتملة' },
  { id: '2', title: 'فحص أنظمة الإنذار', date: '2026-06-22', status: 'مكتملة' },
  { id: '3', title: 'تركيب كاميرات مراقبة', date: '2026-05-30', status: 'مكتملة' },
]
</script>

<style scoped>
.page-title { font-size: 19px; margin-bottom: 4px; }
.page-subtitle { font-size: 13px; color: var(--ink-soft); margin-bottom: 18px; }
</style>
