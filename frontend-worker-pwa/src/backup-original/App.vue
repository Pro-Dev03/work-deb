<template>
  <div class="shell">
    <header v-if="!isLoginPage" class="topbar">
      <div class="brand">
        <div class="brand-mark">
          <MapPin :size="40" />
        </div>
        <div class="brand-info">
          <span class="brand-name gradient-text">{{ t('app_name') }}</span>
          <span class="devpro-badge">Powered by DevPro</span>
        </div>
      </div>
      <router-link to="/profile" class="avatar">{{ initials }}</router-link>
    </header>
    <main class="content" :class="{ 'content--login': isLoginPage }"><router-view /></main>
    <nav v-if="!isLoginPage" class="tabbar">
      <router-link to="/attendance" class="tab" active-class="tab--active">
        <span class="tab__icon">
          <Clock :size="20" />
        </span>
        <span>{{ t('attendance') }}</span>
      </router-link>
      <router-link to="/notes" class="tab" active-class="tab--active">
        <span class="tab__icon">
          <FileText :size="20" />
        </span>
        <span>{{ t('notes') }}</span>
      </router-link>
      <router-link to="/profile" class="tab" active-class="tab--active">
        <span class="tab__icon">
          <User :size="20" />
        </span>
        <span>{{ t('profile') }}</span>
      </router-link>
    </nav>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from './services/i18n'
import { authStore } from './store/auth'
import { MapPin, Clock, User, FileText } from '@lucide/vue'

const route = useRoute()
const { t } = useI18n()
const initials = computed(() => (authStore.user?.full_name || 'م ع').slice(0, 1))
const isLoginPage = computed(() => route.path === '/login')
</script>

<style scoped>
.shell { 
  min-height: 100dvh; 
  display: flex; 
  flex-direction: column; 
  background: var(--canvas); 
}

.topbar { 
  display: flex; 
  justify-content: space-between; 
  padding: 12px 16px; 
  background: var(--surface); 
  border-bottom: 1px solid var(--line);
  box-shadow: var(--shadow-sm);
  transition: box-shadow var(--transition-base);
}

.topbar:hover {
  box-shadow: var(--shadow-md);
}

.brand {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 700;
}

.brand-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.brand-mark { 
  width: 40px; 
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--brand);
  transition: transform var(--transition-base);
}

.brand:hover .brand-mark {
  transform: scale(1.1);
}

.brand-name {
  font-size: 16px;
  font-weight: 700;
}

.devpro-badge {
  font-size: 10px;
  font-weight: 600;
  color: var(--ink-soft);
  letter-spacing: 0.5px;
  text-transform: uppercase;
}

.avatar { 
  width: 34px; 
  height: 34px; 
  border-radius: 50%; 
  background: linear-gradient(135deg, var(--brand-tint) 0%, var(--brand-tint) 100%); 
  color: var(--brand); 
  display: flex; 
  align-items: center; 
  justify-content: center; 
  font-weight: 700;
  box-shadow: var(--shadow-sm);
  transition: all var(--transition-base);
}

.avatar:hover {
  transform: scale(1.1);
  box-shadow: var(--shadow-md), var(--brand-glow);
}

.content { 
  flex: 1; 
  padding: 16px; 
  padding-bottom: 80px; 
  max-width: 480px; 
  width: 100%; 
  margin: 0 auto;
  animation: fadeIn 0.4s ease;
}

.content--login {
  padding-bottom: 16px;
}

.tabbar { 
  position: fixed; 
  bottom: 0; 
  left: 0; 
  right: 0; 
  display: flex; 
  background: var(--surface); 
  border-top: 1px solid var(--line); 
  padding: 6px 8px; 
  max-width: 480px; 
  margin: 0 auto;
  box-shadow: 0 -2px 10px rgba(0,0,0,0.05);
}

.tab { 
  flex: 1; 
  display: flex; 
  flex-direction: column; 
  align-items: center; 
  gap: 4px; 
  padding: 8px; 
  border-radius: var(--radius-sm); 
  color: var(--ink-soft); 
  font-size: 11px; 
  font-weight: 600; 
  cursor: pointer; 
  border: none; 
  background: transparent;
  transition: all var(--transition-base);
  position: relative;
  overflow: hidden;
}

.tab::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, rgba(255,255,255,0.1) 0%, rgba(255,255,255,0) 100%);
  opacity: 0;
  transition: opacity var(--transition-fast);
}

.tab:hover::before {
  opacity: 1;
}

.tab:hover {
  background: var(--brand-tint);
  transform: translateY(-2px);
}

.tab--active { 
  color: var(--brand); 
  background: linear-gradient(135deg, var(--brand-tint) 0%, var(--brand-tint) 100%);
  box-shadow: var(--shadow-sm);
}

.tab__icon { 
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: transform var(--transition-base);
}

.tab:hover .tab__icon {
  transform: scale(1.1);
}

.tab--active .tab__icon {
  color: var(--brand);
}
</style>
