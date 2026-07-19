-- حذف عمود ملاحظات إنهاء الدوام
ALTER TABLE attendance DROP COLUMN IF EXISTS check_out_notes;
