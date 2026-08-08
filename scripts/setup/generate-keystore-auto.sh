#!/bin/bash

# WorkTrack - إنشاء Keystore تلقائياً
# ⚠️  هام: هذا السكريبت يُستخدم لأغراض التطوير فقط
# للإنتاج، استخدم create-keystore.sh لإدخال معلومات حقيقية

echo "🔐 إنشاء Keystore تلقائياً لتطبيقات WorkTrack"
echo "=============================================="
echo ""

# إعدادات افتراضية (للتطوير فقط!)
STORE_PASSWORD="WorkTrackSecure2024!Dev"
KEY_PASSWORD="WorkTrackSecure2024!Dev"
CN="WorkTrack Developer"
OU="Development"
ORG="WorkTrack"
CITY="Riyadh"
STATE="Riyadh Region"
COUNTRY="SA"

KEYSTORE_DIR="android-keystores"

# إنشاء المجلد
if [ ! -d "$KEYSTORE_DIR" ]; then
    mkdir -p "$KEYSTORE_DIR"
    echo "✅ تم إنشاء مجلد android-keystores"
fi

echo "🔧 جاري إنشاء Keystore..."
echo "المعلومات:"
echo "  CN: $CN"
echo "  OU: $OU"
echo "  O: $ORG"
echo "  L: $CITY"
echo "  ST: $STATE"
echo "  C: $COUNTRY"
echo ""

keytool -genkey -v \
  -keystore "$KEYSTORE_DIR/worktrack-release.keystore" \
  -alias worktrack \
  -keyalg RSA \
  -keysize 2048 \
  -validity 10000 \
  -storepass "$STORE_PASSWORD" \
  -keypass "$KEY_PASSWORD" \
  -dname "CN=$CN, OU=$OU, O=$ORG, L=$CITY, ST=$STATE, C=$COUNTRY"

if [ $? -eq 0 ]; then
    echo ""
    echo "✅ تم إنشاء Keystore بنجاح!"
    
    # تأمين الملف
    chmod 600 "$KEYSTORE_DIR/worktrack-release.keystore"
    
    # إضافة إلى .gitignore
    if ! grep -q "android-keystores/" .gitignore 2>/dev/null; then
        echo "android-keystores/" >> .gitignore
        echo "✅ تم إضافة android-keystores/ إلى .gitignore"
    fi
    
    echo ""
    echo "⚠️  تحذيرات أمنية:"
    echo "1. هذا Keystore للتطوير فقط!"
    echo "2. للإنتاج، استخدم create-keystore.sh"
    echo "3. احتفظ بنسخة احتياطية في مكان آمن"
    echo "4. لا ترفع ملف Keystore إلى Git"
    
    echo ""
    echo "📄 جاري تحديث ملفات keystore.properties..."
    
    # تحديث ملفات keystore.properties
    for app in "frontend-worker-pwa" "frontend-admin-dashboard" "frontend-client-portal"; do
        if [ -d "$app/android" ]; then
            cat > "$app/android/keystore.properties" << EOF
RELEASE_STORE_FILE=../../android-keystores/worktrack-release.keystore
RELEASE_KEY_ALIAS=worktrack
RELEASE_STORE_PASSWORD=$STORE_PASSWORD
RELEASE_KEY_PASSWORD=$KEY_PASSWORD
EOF
            echo "✅ تم تحديث keystore.properties لـ $app"
        fi
    done
    
    echo ""
    echo "✅ تم التحضير للبناء بنجاح!"
    echo ""
    echo "🎯 يمكنك الآن تشغيل: ./build-all-release.sh"
    
else
    echo ""
    echo "❌ فشل إنشاء Keystore"
    exit 1
fi
