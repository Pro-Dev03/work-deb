<template>
  <div class="shell">
    <header class="topbar">
      <div class="topbar__brand">
        <span class="brand-mark"></span>
        <div class="brand-text">
          <span class="brand-name">WorkTrack</span>
          <span class="brand-sub">بوابة العميل</span>
        </div>
      </div>
      <div class="topbar__actions">
        <LanguageSwitcher />
        <router-link to="/new" class="btn btn--primary btn--sm">طلب خدمة</router-link>
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
}

.topbar__brand {
  display: flex;
  align-items: center;
  gap: 10px;
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
}

.footer {
  text-align: center;
  padding: 16px 20px;
  border-top: 1px solid var(--line);
  background: var(--surface);
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
