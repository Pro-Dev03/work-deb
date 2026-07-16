<template>
  <div>
    <div class="page-head">
      <div>
        <h2>{{ t('clients_title') }}</h2>
        <p>{{ t('clients_description') }}</p>
      </div>
      <button class="btn btn--primary">+ {{ t('new_client') }}</button>
    </div>

    <div class="card">
      <!-- جدول للشاشات الكبيرة -->
      <div class="table-wrapper desktop-only">
        <table class="table">
          <thead>
            <tr>
              <th>{{ t('clients_name') }}</th>
              <th>{{ t('clients_phone') }}</th>
              <th>{{ t('clients_email') }}</th>
              <th>{{ t('clients_service_count') }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="c in clients" :key="c.id">
              <td>{{ c.name }}</td>
              <td class="mono">{{ c.phone }}</td>
              <td class="mono">{{ c.email }}</td>
              <td>{{ c.servicesCount }}</td>
            </tr>
          </tbody>
        </table>
      </div>
      
      <!-- بطاقات للشاشات الصغيرة -->
      <div class="mobile-cards mobile-only">
        <div v-for="c in clients" :key="c.id" class="client-card">
          <div class="client-card__header">
            <span class="client-card__name">{{ c.name }}</span>
            <span class="badge badge--info">{{ c.servicesCount }} {{ t('clients_service_count') }}</span>
          </div>
          <div class="client-card__body">
            <div class="client-card__row">
              <span class="client-card__label">{{ t('clients_phone') }}</span>
              <span class="client-card__value mono">{{ c.phone }}</span>
            </div>
            <div class="client-card__row">
              <span class="client-card__label">{{ t('clients_email') }}</span>
              <span class="client-card__value mono">{{ c.email }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useI18n } from '../services/i18n'

const { t } = useI18n()

const clients = ref([
  { id: 1, name: 'شركة الأفق للعقارات', phone: '079xxxxxxx', email: 'contact@ofoq.com', servicesCount: 12 },
  { id: 2, name: 'مجمع الزهور السكني', phone: '077xxxxxxx', email: 'admin@zohour.com', servicesCount: 5 },
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

/* تصميم بطاقات العملاء للهاتف */
.mobile-cards {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.client-card {
  background: var(--surface);
  border: 1px solid var(--line);
  border-radius: var(--radius-md);
  padding: 16px;
  transition: all 0.2s ease;
}

.client-card:hover {
  border-color: var(--brand);
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

.client-card__header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  padding-bottom: 12px;
  border-bottom: 1px solid var(--line);
}

.client-card__name {
  font-weight: 600;
  color: var(--ink);
  font-size: 15px;
}

.client-card__body {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.client-card__row {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.client-card__label {
  font-size: 13px;
  color: var(--ink-soft);
  font-weight: 500;
}

.client-card__value {
  font-size: 14px;
  color: var(--ink);
  font-weight: 500;
}
</style>
