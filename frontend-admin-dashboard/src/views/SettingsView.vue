<template>
  <div>
    <div class="page-head"><h2>الإعدادات</h2></div>

    <div class="card settings-card">
      <h3>اللغة الافتراضية للواجهات</h3>
      <div class="settings-card__options">
        <label class="field"><span>لوحة المدير وتطبيق الموظف</span>
          <select><option>العربية</option><option>עברית</option><option>English</option></select>
        </label>
      </div>
    </div>

    <div class="card settings-card">
      <h3>نطاق التختيم الافتراضي</h3>
      <p class="settings-card__hint">القيمة المستخدمة عند إنشاء نقطة عمل جديدة إن لم يُحدَّد نصف قطر مخصص</p>
      <label class="field" style="max-width: 220px"><span>نصف القطر (متر)</span><input type="number" value="100" /></label>
    </div>
  </div>
</template>

<style scoped>
.page-head { margin-bottom: 18px; }
.page-head h2 { font-size: 20px; }
.settings-card { padding: 22px; margin-bottom: 16px; }
.settings-card h3 { font-size: 15px; margin-bottom: 10px; }
.settings-card__hint { font-size: 13px; color: var(--ink-soft); margin-bottom: 12px; }
</style>
