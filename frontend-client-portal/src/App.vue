<template>
  <div class="shell">
    <header class="topbar">
      <div class="topbar__brand">
        <span class="brand-mark">
          <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"></path>
            <circle cx="12" cy="10" r="3"></circle>
          </svg>
        </span>
        <div class="brand-text">
          <span class="brand-name gradient-text">WorkTrack</span>
          <span class="brand-sub">بوابة العميل</span>
        </div>
      </div>
      <div class="topbar__actions">
        <LanguageSwitcher />
        <router-link to="/new" class="btn btn--primary btn--sm">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M12 5v14M5 12h14"/>
          </svg>
          طلب خدمة
        </router-link>
      </div>
    </header>

    <main class="content">
      <router-view />
    </main>

    <footer class="footer">
      <p>© 2026 WorkTrack - جميع الحقوق محفوظة</p>
    </footer>
  </div>
</template>

<script setup>
import LanguageSwitcher from './components/LanguageSwitcher.vue'
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
  align-items: center;
  justify-content: space-between;
  padding: 12px 20px;
  background: var(--surface);
  border-bottom: 1px solid var(--line);
  position: sticky;
  top: 0;
  z-index: 10;
  box-shadow: var(--shadow-sm);
  transition: box-shadow var(--transition-base);
}

.topbar:hover {
  box-shadow: var(--shadow-md);
}

.topbar__brand {
  display: flex;
  align-items: center;
  gap: 10px;
}

.brand-mark {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--brand);
  transition: transform var(--transition-base);
}

.topbar__brand:hover .brand-mark {
  transform: scale(1.1);
}

.brand-text {
  display: flex;
  flex-direction: column;
  line-height: 1.2;
}

.brand-name {
  font-weight: 700;
  font-size: 16px;
  color: var(--brand);
}

.brand-sub {
  font-size: 10px;
  color: var(--ink-soft);
  font-weight: 500;
}

.topbar__actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.content {
  flex: 1;
  padding: 20px;
  max-width: 600px;
  width: 100%;
  margin: 0 auto;
  animation: fadeIn 0.4s ease;
}

.footer {
  text-align: center;
  padding: 16px 20px;
  border-top: 1px solid var(--line);
  background: var(--surface);
  box-shadow: 0 -2px 10px rgba(0,0,0,0.05);
}

.footer p {
  font-size: 12px;
  color: var(--ink-soft);
}

@media (max-width: 480px) {
  .topbar {
    padding: 10px 14px;
    flex-wrap: wrap;
    gap: 8px;
  }
  .topbar__actions {
    gap: 6px;
  }
  .topbar__actions .btn {
    font-size: 12px;
    padding: 6px 12px;
  }
  .content {
    padding: 14px;
  }
}
</style>
