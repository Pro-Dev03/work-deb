<template>
  <div>
    <router-link to="/" class="back-link">→ العودة إلى السجل</router-link>

    <div class="card report">
      <span class="badge badge--in">مكتملة</span>
      <h2>صيانة مكيّفات دورية</h2>
      <p class="report__meta mono">تم التنفيذ: 2026-07-08، 10:15 — 11:40</p>
      <p class="report__meta">الموقع: برج الأمل، عمّان</p>

      <hr class="divider" />

      <h3>ملخص العمل المنجز</h3>
      <p class="report__text">
        تم فحص جميع وحدات التكييف وتنظيف الفلاتر واستبدال قطعتين تالفتين، مع التأكد من
        عمل النظام بكفاءة كاملة قبل مغادرة الموقع.
      </p>
    </div>

    <router-link :to="`/rate/${$route.params.id}`" class="btn btn--primary btn--block">
      تقييم الخدمة
    </router-link>
  </div>
</template>

<style scoped>
.back-link { display: inline-block; font-size: 13px; color: var(--ink-soft); margin-bottom: 14px; }
.report { padding: 22px; margin-bottom: 18px; }
.report h2 { font-size: 18px; margin: 10px 0 6px; }
.report__meta { font-size: 13px; color: var(--ink-soft); margin-bottom: 2px; }
.divider { border: none; border-top: 1px solid var(--line); margin: 18px 0; }
.report h3 { font-size: 14px; margin-bottom: 8px; }
.report__text { font-size: 13.5px; color: var(--ink-soft); line-height: 1.8; }
</style>
