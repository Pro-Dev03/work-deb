#!/bin/bash
# سكربت shell لتغيير إيميل أو كلمة مرور المستخدم باستخدام API
# يتطلب: curl

# قراءة API URL و API key من ملف .env أو استخدام الافتراضي
DEFAULT_API_URL="https://worktrack-v2.onrender.com/api/v1"
DEFAULT_API_KEY="dev-admin-script-key-2024"

if [ -f "backend/.env" ]; then
    ENV_FILE="backend/.env"
elif [ -f ".env" ]; then
    ENV_FILE=".env"
else
    echo "⚠️  ملف .env غير موجود، استخدام القيم الافتراضية"
    API_URL="$DEFAULT_API_URL"
    API_KEY="$DEFAULT_API_KEY"
fi

# استخراج الإعدادات من الملف إذا وجد
if [ -n "$ENV_FILE" ]; then
    # قراءة API URL
    API_URL=$(grep "^API_URL=" "$ENV_FILE" 2>/dev/null | cut -d '=' -f2-)
    if [ -z "$API_URL" ]; then
        API_URL=$(grep "^BACKEND_URL=" "$ENV_FILE" 2>/dev/null | cut -d '=' -f2-)
    fi
    if [ -z "$API_URL" ]; then
        API_URL="$DEFAULT_API_URL"
    fi
    
    # قراءة API Key
    API_KEY=$(grep "^ADMIN_SCRIPT_API_KEY=" "$ENV_FILE" 2>/dev/null | cut -d '=' -f2-)
    if [ -z "$API_KEY" ]; then
        API_KEY="$DEFAULT_API_KEY"
    fi
fi

echo "============================================================"
echo "سكربت تغيير بيانات المستخدم - WorkTrack"
echo "============================================================"
echo "API URL: $API_URL"
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
            
            response=$(curl -s -X GET "$API_URL/admin-script/users" \
                -H "X-API-Key: $API_KEY" \
                -H "Content-Type: application/json")
            
            if echo "$response" | grep -q '"error"'; then
                echo "❌ خطأ: $response"
            else
                echo "$response" | jq -r '.[] | "\(.id) | \(.full_name) | \(.email) | \(.role) | \(.is_active)"' 2>/dev/null || echo "$response"
            fi
            echo "============================================================"
            ;;
        2)
            echo ""
            read -p "أدخل ID المستخدم أو الإيميل الحالي: " identifier
            
            # عرض المستخدم أولاً
            response=$(curl -s -X GET "$API_URL/admin-script/users" \
                -H "X-API-Key: $API_KEY" \
                -H "Content-Type: application/json")
            
            user_info=$(echo "$response" | jq -r --arg id "$identifier" '.[] | select(.id == $id or .email == $id)' 2>/dev/null)
            
            if [ -z "$user_info" ]; then
                echo "❌ المستخدم غير موجود"
                continue
            fi
            
            user_id=$(echo "$user_info" | jq -r '.id' 2>/dev/null)
            full_name=$(echo "$user_info" | jq -r '.full_name' 2>/dev/null)
            current_email=$(echo "$user_info" | jq -r '.email' 2>/dev/null)
            
            echo "المستخدم المحدد: $full_name ($current_email)"
            read -p "أدخل الإيميل الجديد: " new_email
            
            if [ -z "$new_email" ]; then
                echo "❌ الإيميل لا يمكن أن يكون فارغاً"
                continue
            fi
            
            read -p "تأكيد تغيير الإيميل إلى '$new_email'؟ (n/y): " confirm
            if [ "$confirm" = "y" ]; then
                response=$(curl -s -X PUT "$API_URL/admin-script/users/email" \
                    -H "X-API-Key: $API_KEY" \
                    -H "Content-Type: application/json" \
                    -d "{\"identifier\": \"$identifier\", \"new_email\": \"$new_email\"}")
                
                if echo "$response" | grep -q '"error"'; then
                    echo "❌ خطأ: $response"
                else
                    echo "✅ $(echo "$response" | jq -r '.message' 2>/dev/null || echo 'تم التحديث')"
                fi
            else
                echo "❌ تم إلغاء العملية"
            fi
            ;;
        3)
            echo ""
            read -p "أدخل ID المستخدم أو الإيميل: " identifier
            
            # عرض المستخدم أولاً
            response=$(curl -s -X GET "$API_URL/admin-script/users" \
                -H "X-API-Key: $API_KEY" \
                -H "Content-Type: application/json")
            
            user_info=$(echo "$response" | jq -r --arg id "$identifier" '.[] | select(.id == $id or .email == $id)' 2>/dev/null)
            
            if [ -z "$user_info" ]; then
                echo "❌ المستخدم غير موجود"
                continue
            fi
            
            user_id=$(echo "$user_info" | jq -r '.id' 2>/dev/null)
            full_name=$(echo "$user_info" | jq -r '.full_name' 2>/dev/null)
            current_email=$(echo "$user_info" | jq -r '.email' 2>/dev/null)
            
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
                response=$(curl -s -X PUT "$API_URL/admin-script/users/password" \
                    -H "X-API-Key: $API_KEY" \
                    -H "Content-Type: application/json" \
                    -d "{\"identifier\": \"$identifier\", \"new_password\": \"$new_password\"}")
                
                if echo "$response" | grep -q '"error"'; then
                    echo "❌ خطأ: $response"
                else
                    echo "✅ $(echo "$response" | jq -r '.message' 2>/dev/null || echo 'تم التحديث')"
                fi
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
                response=$(curl -s -X POST "$API_URL/admin-script/users/admin" \
                    -H "X-API-Key: $API_KEY" \
                    -H "Content-Type: application/json" \
                    -d "{\"full_name\": \"$full_name\", \"email\": \"$email\", \"password\": \"$password\"}")
                
                if echo "$response" | grep -q '"error"'; then
                    echo "❌ خطأ: $response"
                else
                    echo "✅ $(echo "$response" | jq -r '.message' 2>/dev/null || echo 'تم الإنشاء')"
                    echo "🆔 معرف المستخدم: $(echo "$response" | jq -r '.user_id' 2>/dev/null || echo 'غير معروف')"
                fi
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
            
            # عرض المستخدم أولاً
            response=$(curl -s -X GET "$API_URL/admin-script/users" \
                -H "X-API-Key: $API_KEY" \
                -H "Content-Type: application/json")
            
            user_info=$(echo "$response" | jq -r --arg id "$identifier" '.[] | select(.id == $id or .email == $id)' 2>/dev/null)
            
            if [ -z "$user_info" ]; then
                echo "❌ المستخدم غير موجود"
                continue
            fi
            
            user_id=$(echo "$user_info" | jq -r '.id' 2>/dev/null)
            full_name=$(echo "$user_info" | jq -r '.full_name' 2>/dev/null)
            current_email=$(echo "$user_info" | jq -r '.email' 2>/dev/null)
            role=$(echo "$user_info" | jq -r '.role' 2>/dev/null)
            
            echo "المستخدم المحدد:"
            echo "الاسم: $full_name"
            echo "الإيميل: $current_email"
            echo "الدور: $role"
            echo ""
            
            read -p "هل أنت متأكد من حذف هذا الحساب؟ هذا الإجراء لا يمكن التراجع عنه (n/y): " confirm_delete
            if [ "$confirm_delete" = "y" ]; then
                response=$(curl -s -X DELETE "$API_URL/admin-script/users" \
                    -H "X-API-Key: $API_KEY" \
                    -H "Content-Type: application/json" \
                    -d "{\"identifier\": \"$identifier\"}")
                
                if echo "$response" | grep -q '"error"'; then
                    echo "❌ خطأ: $response"
                else
                    echo "✅ $(echo "$response" | jq -r '.message' 2>/dev/null || echo 'تم الحذف')"
                fi
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