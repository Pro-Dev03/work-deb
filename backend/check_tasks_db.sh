#!/bin/bash

# فحص جدول المهام مباشرة عبر قاعدة البيانات
# استخدام الاتصال المباشر من ملف .env

# قراءة معلومات الاتصال من ملف .env
source /home/dev-bit/worktrack/backend/.env

echo "🔍 فحص جدول المهام..."
psql "$DATABASE_URL" -c "SELECT COUNT(*) as total_tasks FROM tasks;"

echo ""
echo "🔍 فحص جدول الحضور (attendance)..."
psql "$DATABASE_URL" -c "SELECT COUNT(*) as total_attendance FROM attendance;"

echo ""
echo "📊 توزيع حالات الحضور..."
psql "$DATABASE_URL" -c "SELECT status, COUNT(*) as count FROM attendance GROUP BY status ORDER BY status;"

echo ""
echo "📋 عينة من الحضور مع تفاصيل..."
psql "$DATABASE_URL" -c "SELECT id, user_id, task_id, status, check_in_time, check_out_time FROM attendance LIMIT 5;"

echo ""
echo "🏗️ هيكل جدول الحضور..."
psql "$DATABASE_URL" -c "\d attendance"