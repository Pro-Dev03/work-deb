<template>
  <div class="company-brand" :class="size">
    <img :src="logoUrl" :alt="companyName" class="brand-logo" />
    <div class="brand-text">
      <span class="brand-name">{{ companyName }}</span>
      <span v-if="showSlogan" class="brand-slogan">{{ slogan }}</span>
    </div>
  </div>
</template>

<script setup>
import { companyConfig } from '../config/company'

const props = defineProps({
  size: { type: String, default: 'medium' },
  showSlogan: { type: Boolean, default: true },
})

const companyName = companyConfig.name
const slogan = companyConfig.slogan
const logoUrl = companyConfig.logo
</script>

<style scoped>
.company-brand {
  display: flex;
  align-items: center;
  gap: 12px;
}

.brand-logo {
  object-fit: contain;
  border-radius: 12px;
  background: white;
  padding: 4px;
}

/* الأحجام */
.small .brand-logo { width: 28px; height: 28px; border-radius: 8px; }
.small .brand-name { font-size: 12px; }
.small .brand-slogan { font-size: 9px; }

.medium .brand-logo { width: 40px; height: 40px; }
.medium .brand-name { font-size: 16px; }
.medium .brand-slogan { font-size: 10px; }

.large .brand-logo { width: 60px; height: 60px; }
.large .brand-name { font-size: 20px; }
.large .brand-slogan { font-size: 12px; }

.brand-text {
  display: flex;
  flex-direction: column;
  line-height: 1.2;
}

.brand-name {
  font-weight: 700;
  color: var(--brand);
  letter-spacing: -0.5px;
}

.brand-slogan {
  color: var(--ink-soft);
  font-weight: 400;
  font-size: 10px;
}
</style>
