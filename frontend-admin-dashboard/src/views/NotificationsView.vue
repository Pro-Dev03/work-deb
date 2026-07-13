<template>
  <div>
    <div class="page-head"><h2>{{ t('notifications') }}</h2></div>
    <div class="card">
      <ActivityFeed :items="notifications" />
    </div>
  </div>
</template>

<script setup>
import ActivityFeed from '../components/ActivityFeed.vue'
import { useI18n } from '../services/i18n'

const { t } = useI18n()

const notifications = [
  { id: 1, text: 'محاولة تختيم مرفوضة خارج النطاق — سارة قدورة', time: 'اليوم، 08:47', tone: 'out' },
  { id: 2, text: 'تم إنشاء حساب موظف جديد — منى الحوراني', time: 'أمس، 17:20', tone: 'neutral' },
  { id: 3, text: 'اكتملت جميع مهام اليوم في برج الأمل', time: 'أمس، 16:40', tone: 'in' },
]
</script>

<style scoped>
.page-head { margin-bottom: 18px; }
.page-head h2 { font-size: 20px; }
.card { padding: 22px; }
</style>
