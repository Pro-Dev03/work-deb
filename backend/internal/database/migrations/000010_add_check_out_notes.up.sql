-- إضافة عمود لملاحظات إنهاء الدوام
ALTER TABLE attendance ADD COLUMN IF NOT EXISTS check_out_notes TEXT;
