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
      {{ lang.flag }}
    </button>
  </div>
</template>

<script setup>
import { useI18n } from '../services/i18n'

const { currentLang, setLang } = useI18n()

const languages = [
  { code: 'ar', name: 'العربية', flag: '🇸🇦' },
  { code: 'he', name: 'עברית', flag: '🇮🇱' },
  { code: 'en', name: 'English', flag: '🇬🇧' }
]

function changeLanguage(lang) {
  setLang(lang)
}
</script>

<style scoped>
.lang-switcher {
  display: flex;
  justify-content: center;
  gap: 2px;
  padding: 2px;
  background: var(--canvas, #f0f4fa);
  border-radius: 8px;
  flex: 1;
}

.lang-btn {
  padding: 4px 6px;
  border: none;
  border-radius: 6px;
  background: transparent;
  font-size: 16px;
  cursor: pointer;
  transition: all 0.3s ease;
  font-family: inherit;
  min-width: 32px;
}

.lang-btn:hover {
  background: rgba(30, 58, 95, 0.08);
  transform: scale(1.05);
}

.lang-btn.active {
  background: #1E3A5F;
  color: white;
  box-shadow: 0 2px 8px rgba(30, 58, 95, 0.3);
  transform: scale(1.05);
  border-radius: 6px;
}

@media (max-width: 480px) {
  .lang-btn {
    font-size: 14px;
    padding: 3px 4px;
    min-width: 28px;
  }
}
</style>
