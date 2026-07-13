<template>
  <div class="rate-page">
    <h2>كيف كانت تجربتكم؟</h2>
    <p>تقييمكم يساعدنا على تحسين جودة الخدمة المقدّمة</p>

    <div class="card rate-card">
      <StarRating />
      <textarea class="field" placeholder="أخبرونا بملاحظاتكم (اختياري)" rows="4"></textarea>
      <button class="btn btn--primary btn--block">إرسال التقييم</button>
    </div>
  </div>
</template>

<script setup>
import StarRating from '../components/StarRating.vue'
</script>

<style scoped>
.rate-page { text-align: center; padding-top: 10px; }
.rate-page h2 { font-size: 19px; margin-bottom: 6px; }
.rate-page p { font-size: 13px; color: var(--ink-soft); margin-bottom: 20px; }
.rate-card { padding: 26px 20px; display: flex; flex-direction: column; gap: 18px; text-align: right; }
.rate-card textarea {
  width: 100%; font-family: var(--font-body); font-size: 14px; padding: 12px 14px;
  border-radius: var(--radius-sm); border: 1px solid var(--line); resize: vertical;
}
</style>
