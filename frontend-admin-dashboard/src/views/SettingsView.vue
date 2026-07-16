<template>
  <div>
    <div class="page-head">
      <h2>{{ t('settings') }}</h2>
      <p class="page-subtitle">{{ t('settings_intro') }}</p>
    </div>

    <div class="card subscription-card">
      <h3>{{ t('subscription_title') }}</h3>
      <p class="settings-card__hint">{{ t('subscription_description') }}</p>
      <SubscriptionStatusView />
    </div>

    <div class="card settings-card">
      <h3>{{ t('settings_language_title') }}</h3>
      <div class="settings-card__options">
        <label class="field"><span>{{ t('settings_language_label') }}</span>
          <select v-model="currentLang" @change="changeLanguage">
            <option value="ar">{{ t('ar') }}</option>
            <option value="he">{{ t('he') }}</option>
            <option value="en">{{ t('en') }}</option>
          </select>
        </label>
      </div>
    </div>

    <div class="card settings-card">
      <h3>{{ t('settings_geofence_title') }}</h3>
      <p class="settings-card__hint">{{ t('settings_geofence_hint') }}</p>
      <label class="field" style="max-width: 220px"><span>{{ t('settings_geofence_radius') }}</span><input type="number" value="100" disabled /></label>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from '../services/i18n'
import SubscriptionStatusView from './SubscriptionStatusView.vue'

const { t, setLang, currentLang: computedLang } = useI18n()
const currentLang = ref(computedLang.value)

const changeLanguage = () => {
  setLang(currentLang.value)
}

onMounted(() => {
  currentLang.value = computedLang.value
})
</script>

<style scoped>
.page-head { margin-bottom: 18px; }
.page-head h2 { font-size: 20px; }
.page-subtitle { color: var(--ink-soft); margin: 6px 0 16px; font-size: 14px; }
.settings-card { padding: 22px; margin-bottom: 16px; }
.settings-card h3 { font-size: 15px; margin-bottom: 10px; }
.settings-card__hint { font-size: 13px; color: var(--ink-soft); margin-bottom: 12px; }
.field { display: flex; flex-direction: column; gap: 8px; }
.field span { font-size: 13px; color: var(--ink-soft); }
.field input,
.field select { width: 100%; padding: 10px 12px; border-radius: var(--radius-sm); border: 1px solid var(--line); background: var(--surface); }
</style>
