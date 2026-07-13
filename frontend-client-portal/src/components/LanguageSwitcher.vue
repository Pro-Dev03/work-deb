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
      <span class="lang-label">{{ lang.name }}</span>
    </button>
  </div>
</template>

<script setup>
import { computed } from 'vue'
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
  gap: 6px;
  padding: 6px;
  background: var(--canvas, #f0f4fa);
  border-radius: 12px;
  width: 100%;
}

.lang-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 14px;
  border: none;
  border-radius: 8px;
  background: transparent;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.3s ease;
  font-family: inherit;
  flex: 1;
  justify-content: center;
}

.lang-btn:hover {
  background: rgba(30, 58, 95, 0.08);
  transform: scale(1.02);
}

.lang-btn.active {
  background: #1E3A5F;
  color: white;
  box-shadow: 0 2px 12px rgba(30, 58, 95, 0.3);
  transform: scale(1.02);
}

.lang-label {
  font-size: 12px;
  font-weight: 500;
}

@media (max-width: 480px) {
  .lang-btn {
    padding: 4px 8px;
    font-size: 12px;
  }
  .lang-label {
    display: none;
  }
}
</style>
