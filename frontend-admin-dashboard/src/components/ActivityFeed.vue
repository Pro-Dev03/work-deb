<template>
  <ul class="feed">
    <li v-for="item in items" :key="item.id" class="feed__item">
      <span class="feed__dot" :class="item.tone"></span>
      <div>
        <p class="feed__text">{{ item.text }}</p>
        <p class="feed__time mono">{{ item.time }}</p>
      </div>
    </li>
  </ul>
</template>

<script setup>
defineProps({
  items: {
    type: Array,
    default: () => [
      { id: 1, text: 'أحمد ياسين سجّل بدء الدوام — برج الأمل', time: '09:02', tone: 'in' },
      { id: 2, text: 'محاولة تختيم مرفوضة خارج النطاق — سارة قدورة', time: '08:47', tone: 'out' },
      { id: 3, text: 'تم إنشاء مهمة جديدة — معرض السيارات', time: '08:15', tone: 'neutral' },
    ],
  },
})
</script>

<style scoped>
.feed { list-style: none; margin: 0; padding: 0; display: flex; flex-direction: column; gap: 14px; }
.feed__item { display: flex; align-items: flex-start; gap: 10px; }
.feed__dot { width: 8px; height: 8px; border-radius: 50%; margin-top: 6px; flex-shrink: 0; background: var(--line-strong); }
.feed__dot.in { background: var(--signal-in); }
.feed__dot.out { background: var(--signal-out); }
.feed__text { font-size: 13.5px; }
.feed__time { font-size: 12px; color: var(--ink-soft); margin-top: 2px; }
</style>
