-- إزالة حقول تحليل الورديات
ALTER TABLE attendance DROP COLUMN IF EXISTS spans_multiple_days;
ALTER TABLE attendance DROP COLUMN IF EXISTS day_one_date;
ALTER TABLE attendance DROP COLUMN IF EXISTS day_two_date;
ALTER TABLE attendance DROP COLUMN IF EXISTS day_one_hours;
ALTER TABLE attendance DROP COLUMN IF EXISTS day_two_hours;
ALTER TABLE attendance DROP COLUMN IF EXISTS night_hours;
ALTER TABLE attendance DROP COLUMN IF EXISTS day_hours;
ALTER TABLE attendance DROP COLUMN IF EXISTS is_night_shift;
