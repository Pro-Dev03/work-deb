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
.page-head { display: flex; align-items: flex-start; justify-content: space-between; margin-bottom: 18px; }
.page-head h2 { font-size: 20px; margin-bottom: 4px; }
.page-head p { font-size: 13px; color: var(--ink-soft); }
.table { width: 100%; border-collapse: collapse; }
.table th { text-align: right; font-size: 12px; color: var(--ink-soft); font-weight: 600; padding: 14px 20px; border-bottom: 1px solid var(--line); }
.table td { padding: 14px 20px; font-size: 14px; border-bottom: 1px solid var(--line); }
.table tr:last-child td { border-bottom: none; }
</style>
