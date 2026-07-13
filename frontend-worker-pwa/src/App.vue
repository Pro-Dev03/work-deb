<template>
  <div class="shell">
    <header class="topbar">
      <div class="brand">
        <img class="brand-mark" src="./assets/company-logo.jpg" alt="WorkTrack logo" />
        <span>{{ t('app_name') }}</span>
      </div>
      <router-link to="/profile" class="avatar">{{ initials }}</router-link>
    </header>
    <main class="content"><router-view /></main>
    <nav class="tabbar">
      <router-link to="/attendance" class="tab" active-class="tab--active">
        <span class="tab__icon">⏱️</span><span>{{ t('attendance') }}</span>
      </router-link>
      <router-link to="/profile" class="tab" active-class="tab--active">
        <span class="tab__icon">👤</span><span>{{ t('profile') }}</span>
      </router-link>
    </nav>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from './services/i18n'
import { authStore } from './store/auth'

const { t } = useI18n()
const initials = computed(() => (authStore.user?.full_name || 'م ع').slice(0, 1))
</script>

<style scoped>
.shell { min-height: 100dvh; display: flex; flex-direction: column; background: var(--canvas); }
.topbar { display: flex; justify-content: space-between; padding: 12px 16px; background: var(--surface); border-bottom: 1px solid var(--line); }
.brand { display: flex; align-items: center; gap: 8px; font-weight: 700; }
.brand-mark { width: 40px; height: 40px; object-fit: contain; border-radius: 10px; }
.avatar { width: 34px; height: 34px; border-radius: 50%; background: var(--brand-tint); color: var(--brand); display: flex; align-items: center; justify-content: center; font-weight: 700; }
.content { flex: 1; padding: 16px; padding-bottom: 80px; max-width: 480px; width: 100%; margin: 0 auto; }
.tabbar { position: fixed; bottom: 0; left: 0; right: 0; display: flex; background: var(--surface); border-top: 1px solid var(--line); padding: 6px 8px; max-width: 480px; margin: 0 auto; }
.tab { flex: 1; display: flex; flex-direction: column; align-items: center; gap: 2px; padding: 8px; border-radius: var(--radius-sm); color: var(--ink-soft); font-size: 11px; font-weight: 600; cursor: pointer; border: none; background: transparent; }
.tab--active { color: var(--brand); background: var(--brand-tint); }
.tab__icon { font-size: 20px; }
</style>
