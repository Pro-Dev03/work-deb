#!/bin/bash

# سكريبت لتوليد أسرار جديدة لمشروع WorkTrack
# الاستخدام: ./scripts/utils/generate-secrets.sh

set -e

echo "🔐 توليد أسرار جديدة لـ WorkTrack"
echo "=================================="
echo ""

# توليد JWT Secret
echo "1. توليد JWT Secret (32+ حرف):"
JWT_SECRET=$(openssl rand -base64 32)
echo "JWT_SECRET=$JWT_SECRET"
echo ""

# توليد Database Password
echo "2. توليد كلمة مرور قاعدة البيانات:"
DB_PASSWORD=$(openssl rand -base64 16 | tr -d "=+/" | cut -c1-20)
echo "DB_PASSWORD=$DB_PASSWORD"
echo ""

# توليد API Key بديل
echo "3. توليد API Key عشوائي:"
API_KEY=$(openssl rand -hex 16)
echo "API_KEY=$API_KEY"
echo ""

# حفظ في ملف مؤقت
SECRET_FILE=".secrets-new.txt"
cat > $SECRET_FILE << EOF
# أسرار جديدة تم توليدها في $(date)
# ⚠️ احفظ هذا الملف في مكان آمن ثم احذفه

JWT_SECRET=$JWT_SECRET
DB_PASSWORD=$DB_PASSWORD
API_KEY=$API_KEY

# DATABASE_URL الجديد (استبدل بـ host الحقيقي):
# DATABASE_URL=postgresql://postgres:$DB_PASSWORD@your-host:5432/postgres?sslmode=require
EOF

echo "✅ تم حفظ الأسرار في: $SECRET_FILE"
echo ""
echo "⚠️  مهم:"
echo "1. انسخ هذه الأسرار إلى Render Dashboard"
echo "2. احذف الملف $SECRET_FILE بعد النسخ"
echo "3. قم بإعادة نشر التطبيق بعد تحديث الأسرار"
echo ""
echo "📋 للنسخ السريع:"
echo "----------------"
cat $SECRET_FILE