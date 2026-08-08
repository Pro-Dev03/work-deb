#!/bin/bash
# سكربت shell لتغيير إيميل أو كلمة مرور المستخدم باستخدام psql
# يتطلب: psql (PostgreSQL client)
# نسخة احتياطية تعمل مع قاعدة البيانات مباشرة

# قراءة DATABASE_URL من ملف .env أو استخدام الافتراضي
DEFAULT_DATABASE_URL="postgresql://postgres.ghmdhpikqzilosqzxzfd:KiQeXPQ16Ie0PEzU@aws-0-ap-southeast-1.pooler.supabase.com:5432/postgres?sslmode=require"

if [ -f "backend/.env" ]; then
    ENV_FILE="backend/.env"
elif [ -f ".env" ]; then
    ENV_FILE=".env"
else
    echo "⚠️  ملف .env غير موجود، استخدام DATABASE_URL الافتراضي"
    DATABASE_URL="$DEFAULT_DATABASE_URL"
fi

# استخراج DATABASE_URL من الملف إذا وجد
if [ -n "$ENV_FILE" ]; then
    DATABASE_URL=$(grep "^DATABASE_URL=" "$ENV_FILE" | cut -d '=' -f2-)
    
    if [ -z "$DATABASE_URL" ]; then
        echo "⚠️  DATABASE_URL غير موجود في ملف $ENV_FILE، استخدام القيمة الافتراضية"
        DATABASE_URL="$DEFAULT_DATABASE_URL"
    fi
fi

echo "============================================================"
echo "سكربت تغيير بيانات المستخدم - WorkTrack (نسخة قاعدة البيانات)"
echo "============================================================"
echo ""

while true; do
    echo "الخيارات:"
    echo "1. عرض جميع المستخدمين"
    echo "2. تغيير إيميل مستخدم"
    echo "3. تغيير كلمة مرور مستخدم"
    echo "4. إنشاء حساب أدمن جديد"
    echo "5. حذف حساب"
    echo "6. خروج"
    echo ""
    read -p "اختر رقم (1-6): " choice
    
    case $choice in
        1)
            echo ""
            echo "قائمة المستخدمين:"
            echo "============================================================"
            psql "$DATABASE_URL" -c "
                SELECT 
                    id, 
                    full_name, 
                    email, 
                    role, 
                    is_active 
                FROM users 
                ORDER BY role, full_name;
            "
            echo "============================================================"
            ;;
        2)
            echo ""
            read -p "أدخل ID المستخدم أو الإيميل الحالي: " identifier
            
            # البحث عن المستخدم مع escaping
            escaped_identifier=$(printf '%s' "$identifier" | sed 's/'\''/\\'\''/g')
            user_info=$(psql "$DATABASE_URL" -t -c "SELECT id, full_name, email FROM users WHERE email = '$escaped_identifier' OR id = '$escaped_identifier' LIMIT 1;")
            
            if [ -z "$user_info" ]; then
                echo "❌ المستخدم غير موجود"
                continue
            fi
            
            user_id=$(echo "$user_info" | awk '{print $1}')
            full_name=$(echo "$user_info" | awk '{print $2}')
            current_email=$(echo "$user_info" | awk '{print $3}')
            
            echo "المستخدم المحدد: $full_name ($current_email)"
            read -p "أدخل الإيميل الجديد: " new_email
            
            if [ -z "$new_email" ]; then
                echo "❌ الإيميل لا يمكن أن يكون فارغاً"
                continue
            fi
            
            read -p "تأكيد تغيير الإيميل إلى '$new_email'؟ (n/y): " confirm
            if [ "$confirm" = "y" ]; then
                escaped_email=$(printf '%s' "$new_email" | sed 's/'\''/\\'\''/g')
                escaped_user_id=$(printf '%s' "$user_id" | sed 's/'\''/\\'\''/g')
                
                psql "$DATABASE_URL" -c "UPDATE users SET email = '$escaped_email', updated_at = NOW() WHERE id = '$escaped_user_id';"
                echo "✅ تم تحديث الإيميل بنجاح"
            else
                echo "❌ تم إلغاء العملية"
            fi
            ;;
        3)
            echo ""
            read -p "أدخل ID المستخدم أو الإيميل: " identifier
            
            # البحث عن المستخدم مع escaping
            escaped_identifier=$(printf '%s' "$identifier" | sed 's/'\''/\\'\''/g')
            user_info=$(psql "$DATABASE_URL" -t -c "SELECT id, full_name, email FROM users WHERE email = '$escaped_identifier' OR id = '$escaped_identifier' LIMIT 1;")
            
            if [ -z "$user_info" ]; then
                echo "❌ المستخدم غير موجود"
                continue
            fi
            
            user_id=$(echo "$user_info" | awk '{print $1}')
            full_name=$(echo "$user_info" | awk '{print $2}')
            current_email=$(echo "$user_info" | awk '{print $3}')
            
            echo "المستخدم المحدد: $full_name ($current_email)"
            read -s -p "أدخل كلمة المرور الجديدة: " new_password
            echo ""
            
            if [ -z "$new_password" ]; then
                echo "❌ كلمة المرور لا يمكن أن تكون فارغة"
                continue
            fi
            
            if [ ${#new_password} -lt 6 ]; then
                echo "❌ كلمة المرور يجب أن تكون 6 أحرف على الأقل"
                continue
            fi
            
            read -s -p "أعد إدخال كلمة المرور الجديدة: " confirm_password
            echo ""
            
            if [ "$new_password" != "$confirm_password" ]; then
                echo "❌ كلمات المرور غير متطابقة"
                continue
            fi
            
            read -p "تأكيد تغيير كلمة المرور؟ (n/y): " confirm
            if [ "$confirm" = "y" ]; then
                escaped_password=$(printf '%s' "$new_password" | sed 's/'\''/\\'\''/g')
                escaped_user_id=$(printf '%s' "$user_id" | sed 's/'\''/\\'\''/g')
                
                # استخدام دالة crypt لتشفير كلمة المرور
                psql "$DATABASE_URL" -c "UPDATE users SET password_hash = crypt('$escaped_password', gen_salt('bf')), updated_at = NOW() WHERE id = '$escaped_user_id';"
                echo "✅ تم تحديث كلمة المرور بنجاح"
            else
                echo "❌ تم إلغاء العملية"
            fi
            ;;
        4)
            echo ""
            echo "============================================================"
            echo "إنشاء حساب أدمن جديد"
            echo "============================================================"
            
            read -p "أدخل الاسم الكامل: " full_name
            if [ -z "$full_name" ]; then
                echo "❌ الاسم لا يمكن أن يكون فارغاً"
                continue
            fi
            
            read -p "أدخل الإيميل: " email
            if [ -z "$email" ]; then
                echo "❌ الإيميل لا يمكن أن يكون فارغاً"
                continue
            fi
            
            if [[ ! "$email" =~ *@*.* ]]; then
                echo "❌ الإيميل غير صحيح"
                continue
            fi
            
            read -s -p "أدخل كلمة المرور: " password
            echo ""
            
            if [ -z "$password" ]; then
                echo "❌ كلمة المرور لا يمكن أن تكون فارغة"
                continue
            fi
            
            if [ ${#password} -lt 8 ]; then
                echo "❌ كلمة المرور يجب أن تكون 8 أحرف على الأقل"
                continue
            fi
            
            read -s -p "أعد إدخال كلمة المرور: " confirm_password
            echo ""
            
            if [ "$password" != "$confirm_password" ]; then
                echo "❌ كلمات المرور غير متطابقة"
                continue
            fi
            
            echo ""
            echo "ملخص الحساب الجديد:"
            echo "الاسم: $full_name"
            echo "الإيميل: $email"
            echo "الدور: admin"
            echo ""
            
            read -p "تأكيد إنشاء حساب الأدمن؟ (n/y): " confirm
            if [ "$confirm" = "y" ]; then
                # التحقق من عدم وجود الإيميل مع escaping
                escaped_email_check=$(printf '%s' "$email" | sed 's/'\''/\\'\''/g')
                existing=$(psql "$DATABASE_URL" -t -c "SELECT id FROM users WHERE email = '$escaped_email_check' LIMIT 1;")
                
                if [ -n "$existing" ]; then
                    echo "❌ الإيميل مستخدم بالفعل"
                    continue
                fi
                
                # إنشاء المستخدم باستخدام printf for escaping
                escaped_name=$(printf '%s' "$full_name" | sed 's/'\''/\\'\''/g')
                escaped_email=$(printf '%s' "$email" | sed 's/'\''/\\'\''/g')
                escaped_password=$(printf '%s' "$password" | sed 's/'\''/\\'\''/g')
                
                result=$(psql "$DATABASE_URL" -t -c "INSERT INTO users (id, full_name, email, password_hash, role, is_active, created_at, updated_at) VALUES (uuid_generate_v4(), '$escaped_name', '$escaped_email', crypt('$escaped_password', gen_salt('bf')), 'admin', TRUE, NOW(), NOW()) RETURNING id;")
                
                user_id=$(echo "$result" | xargs)
                echo "✅ تم إنشاء حساب الأدمن بنجاح"
                echo "🆔 معرف المستخدم: $user_id"
            else
                echo "❌ تم إلغاء العملية"
            fi
            ;;
        5)
            echo ""
            echo "============================================================"
            echo "حذف حساب"
            echo "============================================================"
            
            read -p "أدخل ID المستخدم أو الإيميل للحذف: " identifier
            
            # البحث عن المستخدم مع escaping
            escaped_identifier=$(printf '%s' "$identifier" | sed 's/'\''/\\'\''/g')
            user_info=$(psql "$DATABASE_URL" -t -c "SELECT id, full_name, email, role FROM users WHERE email = '$escaped_identifier' OR id = '$escaped_identifier' LIMIT 1;")
            
            if [ -z "$user_info" ]; then
                echo "❌ المستخدم غير موجود"
                continue
            fi
            
            user_id=$(echo "$user_info" | awk '{print $1}')
            full_name=$(echo "$user_info" | awk '{print $2}')
            current_email=$(echo "$user_info" | awk '{print $3}')
            role=$(echo "$user_info" | awk '{print $4}')
            
            echo "المستخدم المحدد:"
            echo "الاسم: $full_name"
            echo "الإيميل: $current_email"
            echo "الدور: $role"
            echo ""
            
            read -p "هل أنت متأكد من حذف هذا الحساب؟ هذا الإجراء لا يمكن التراجع عنه (n/y): " confirm_delete
            if [ "$confirm_delete" = "y" ]; then
                escaped_user_id=$(printf '%s' "$user_id" | sed 's/'\''/\\'\''/g')
                
                # حذف المستخدم
                psql "$DATABASE_URL" -c "DELETE FROM users WHERE id = '$escaped_user_id';"
                echo "✅ تم حذف الحساب بنجاح"
            else
                echo "❌ تم إلغاء عملية الحذف"
            fi
            ;;
        6)
            echo "👋 خروج..."
            break
            ;;
        *)
            echo "❌ اختيار غير صحيح"
            ;;
    esac
    
    echo ""
done

echo "✅ تم إغلاق السكربت"
