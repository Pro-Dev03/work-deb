<template>
  <div class="stars" role="radiogroup" aria-label="تقييم الخدمة">
    <button
      v-for="n in 5" :key="n" type="button"
      class="stars__item" :class="{ 'stars__item--filled': n <= (hover || rating) }"
      @mouseenter="hover = n" @mouseleave="hover = 0" @click="rating = n"
      :aria-pressed="n <= rating"
    >★</button>
  </div>
</template>

<script setup>
import { ref } from 'vue'
const rating = ref(0)
const hover = ref(0)
</script>

<style scoped>
.stars { display: flex; gap: 6px; justify-content: center; }
.stars__item { font-size: 34px; line-height: 1; background: none; border: none; color: var(--line-strong); cursor: pointer; transition: color .15s ease, transform .1s ease; }
.stars__item:hover { transform: scale(1.08); }
.stars__item--filled { color: var(--gold); }
</style>
