<template>
  <div>
    <div class="page-head">
      <h2>{{ t('subscription_title') }}</h2>
      <p class="page-subtitle">{{ t('subscription_intro') }}</p>
    </div>

    <div class="card subscription-card">
      <div v-if="loading" class="subscription-loading">
        {{ t('loading_data') }}
      </div>

      <div v-else-if="error" class="error">
        {{ error }}
      </div>

      <div v-else>
        <div class="field-row">
          <span>{{ t('subscription_status_label') }}</span>
          <strong :class="['status-badge', `status-badge--${statusClass}`]">{{ subscriptionText }}</strong>
        </div>

        <div class="field-row">
          <span>{{ t('subscription_expires_at_label') }}</span>
          <strong>{{ expiresLabel }}</strong>
        </div>

        <div v-if="subscriptionStatus !== 'active'" class="subscription-note">
          {{ t('subscription_expired_message') }}
        </div>

        <div v-if="subscriptionStatus === 'active'" class="subscription-note subscription-note--active">
          {{ t('subscription_active_message') }}
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from '../services/i18n'
import { getCurrentUser } from '../services/auth'

const { t } = useI18n()
const loading = ref(true)
const error = ref('')
const subscriptionStatus = ref('active')
const subscriptionExpiresAt = ref(null)

const statusClass = computed(() => {
  if (subscriptionStatus.value === 'active') return 'active'
  if (subscriptionStatus.value === 'expired') return 'expired'
  if (subscriptionStatus.value === 'canceled') return 'canceled'
  return 'unknown'
})

const subscriptionText = computed(() => {
  if (subscriptionStatus.value === 'active') return t('subscription_active')
  if (subscriptionStatus.value === 'expired') return t('subscription_expired')
  if (subscriptionStatus.value === 'canceled') return t('subscription_canceled')
  return t('undefined_text')
})

const expiresLabel = computed(() => {
  if (!subscriptionExpiresAt.value) {
    return t('undefined_text')
  }
  return new Date(subscriptionExpiresAt.value).toLocaleDateString()
})

onMounted(async () => {
  try {
    const data = await getCurrentUser()
    subscriptionStatus.value = data.subscription_status || 'active'
    subscriptionExpiresAt.value = data.subscription_expires_at || null
  } catch (err) {
    error.value = err.response?.data?.error || t('error')
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.page-head {
  margin-bottom: 18px;
}

.page-head h2 {
  font-size: 18px;
  margin-bottom: 6px;
}

.page-subtitle {
  color: var(--ink-soft);
  margin: 0;
  font-size: 13px;
}

.subscription-card {
  padding: 22px;
  margin-bottom: 16px;
}

.subscription-loading,
.error {
  font-size: 14px;
  color: var(--ink-soft);
}

.field-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  margin-bottom: 14px;
  padding: 14px 0;
  border-bottom: 1px solid var(--line);
}

.field-row:last-child {
  border-bottom: none;
}

.field-row span {
  color: var(--ink-soft);
}

.status-badge {
  padding: 6px 10px;
  border-radius: 999px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.status-badge--active {
  background: #D7F6E9;
  color: #0D7A43;
}

.status-badge--expired,
.status-badge--canceled {
  background: #FFE1E1;
  color: #C21E1E;
}

.subscription-note {
  margin-top: 18px;
  font-size: 14px;
  color: var(--ink-soft);
  line-height: 1.6;
}

.subscription-note--active {
  color: #0D7A43;
}
</style>
