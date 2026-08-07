<template>
  <div class="lang-switcher">
    <button
      v-for="lang in languages"
      :key="lang.code"
      class="lang-btn"
      :class="{ active: currentLang === lang.code }"
      @click="changeLanguage(lang.code)"
      :title="lang.name"
    >
      <span class="lang-label">{{ lang.name }}</span>
    </button>
  </div>
</template>

<script setup>
import { useI18n } from '../services/i18n'

const { currentLang, setLang } = useI18n()

// removed emoji flags (🇸🇦 🇮🇱 🇬🇧) — the design system requires SVG icons
// only, no emojis, and flags are a poor proxy for language anyway
// (Arabic/Hebrew aren't tied to one country here)
const languages = [
  { code: 'ar', name: 'العربية' },
  { code: 'he', name: 'עברית' },
  { code: 'en', name: 'English' }
]

function changeLanguage(lang) {
  setLang(lang)
}
</script>

<style scoped>
.lang-switcher {
  display: flex;
  justify-content: center;
  gap: 6px;
  padding: 6px;
  background: var(--surface-elevated);
  border-radius: var(--radius-md);
  width: 100%;
}

.lang-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 8px 14px;
  border: none;
  border-radius: var(--radius-sm);
  background: transparent;
  color: var(--text-secondary);
  font-size: var(--text-sm);
  font-weight: var(--font-medium);
  cursor: pointer;
  transition: background .2s ease, color .2s ease;
  font-family: inherit;
  flex: 1;
}

.lang-btn:hover {
  background: var(--primary-50);
  color: var(--primary-600);
}

.lang-btn.active {
  background: var(--primary-500);
  color: var(--text-inverse);
  box-shadow: var(--shadow-sm);
}

.lang-label {
  font-size: var(--text-xs);
  font-weight: var(--font-semibold);
}

@media (max-width: 480px) {
  .lang-btn {
    padding: 6px 8px;
  }
}
</style>
