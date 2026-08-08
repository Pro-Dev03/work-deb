#!/usr/bin/env python3
# سكربت Python لإدارة مستخدمي WorkTrack عبر API
# يتطلب: requests

import os
import sys
import json
import requests
from getpass import getpass
from datetime import datetime, timedelta, timezone

# الإعدادات الافتراضية
DEFAULT_API_URL = "https://worktrack-v2.onrender.com/api/v1"
DEFAULT_API_KEY = "dev-admin-script-key-2024"

def load_config():
    """قراءة الإعدادات من ملف .env"""
    api_url = DEFAULT_API_URL
    api_key = DEFAULT_API_KEY
    supabase_url = ""
    supabase_service_role_key = ""
    
    env_files = ["backend/.env", ".env"]
    for env_file in env_files:
        if os.path.exists(env_file):
            with open(env_file, 'r') as f:
                for line in f:
                    line = line.strip()
                    if line.startswith("API_URL="):
                        api_url = line.split("=", 1)[1].strip()
                    elif line.startswith("BACKEND_URL="):
                        api_url = line.split("=", 1)[1].strip()
                    elif line.startswith("ADMIN_SCRIPT_API_KEY="):
                        api_key = line.split("=", 1)[1].strip()
                    elif line.startswith("SUPABASE_URL="):
                        supabase_url = line.split("=", 1)[1].strip()
                    elif line.startswith("SUPABASE_SERVICE_ROLE_KEY="):
                        supabase_service_role_key = line.split("=", 1)[1].strip()
            break
    
    return api_url, api_key, supabase_url, supabase_service_role_key

def make_request(method, endpoint, data=None, api_key=None):
    """إرسال طلب API"""
    headers = {
        "X-API-Key": api_key,
        "Content-Type": "application/json"
    }
    
    url = f"{DEFAULT_API_URL}/admin-script/{endpoint}"
    
    try:
        if method == "GET":
            response = requests.get(url, headers=headers)
        elif method == "POST":
            response = requests.post(url, headers=headers, json=data)
        elif method == "PUT":
            response = requests.put(url, headers=headers, json=data)
        elif method == "DELETE":
            response = requests.delete(url, headers=headers, json=data)
        
        return response.json()
    except requests.exceptions.RequestException as e:
        return {"error": f"خطأ في الاتصال: {str(e)}"}

def delete_from_supabase_auth(user_id, supabase_url, service_role_key):
    """حذف مستخدم من Supabase Auth"""
    if not supabase_url or not service_role_key:
        print("⚠️ Supabase credentials not configured, skipping Supabase Auth deletion")
        return True
    
    auth_url = f"{supabase_url}/auth/v1/admin/users/{user_id}"
    headers = {
        "apikey": service_role_key,
        "Authorization": f"Bearer {service_role_key}",
        "Content-Type": "application/json"
    }
    
    try:
        response = requests.delete(auth_url, headers=headers)
        if response.status_code >= 200 and response.status_code < 300:
            print("✅ تم حذف المستخدم من Supabase Auth بنجاح")
            return True
        else:
            print(f"⚠️ فشل حذف من Supabase Auth (status: {response.status_code})")
            return True  # نستمر بحذف من قاعدة البيانات
    except requests.exceptions.RequestException as e:
        print(f"⚠️ خطأ في الاتصال بـ Supabase Auth: {str(e)}")
        return True  # نستمر بحذف من قاعدة البيانات

def list_users(api_key):
    """عرض جميع المستخدمين"""
    print("\nقائمة المستخدمين:")
    print("=" * 60)
    
    result = make_request("GET", "users", api_key=api_key)
    
    if "error" in result:
        print(f"❌ خطأ: {result['error']}")
        return
    
    users = result if isinstance(result, list) else []
    
    if not users:
        print("لا يوجد مستخدمين")
        return
    
    print(f"{'ID':<40} {'الاسم':<30} {'الإيميل':<30} {'الدور':<10} {'نشط':<5}")
    print("-" * 120)
    
    for user in users:
        user_id = user.get('id', '')
        full_name = user.get('full_name', '')
        email = user.get('email', '')
        role = user.get('role', '')
        is_active = user.get('is_active', False)
        
        print(f"{user_id:<40} {full_name:<30} {email:<30} {role:<10} {'✓' if is_active else '✗':<5}")
    
    print("=" * 60)
    return users

def find_user(users, identifier):
    """البحث عن مستخدم"""
    for user in users:
        if user.get('id') == identifier or user.get('email') == identifier:
            return user
    return None

def update_user_email(api_key, supabase_url, supabase_service_role_key):
    """تغيير إيميل مستخدم"""
    identifier = input("\nأدخل ID المستخدم أو الإيميل الحالي: ")
    
    # جلب جميع المستخدمين
    result = make_request("GET", "users", api_key=api_key)
    if "error" in result:
        print(f"❌ خطأ: {result['error']}")
        return
    
    users = result if isinstance(result, list) else []
    user = find_user(users, identifier)
    
    if not user:
        print("❌ المستخدم غير موجود")
        return
    
    print(f"المستخدم المحدد: {user['full_name']} ({user['email']})")
    new_email = input("أدخل الإيميل الجديد: ")
    
    if not new_email:
        print("❌ الإيميل لا يمكن أن يكون فارغاً")
        return
    
    confirm = input(f"تأكيد تغيير الإيميل إلى '{new_email}'؟ (n/y): ")
    if confirm.lower() != 'y':
        print("❌ تم إلغاء العملية")
        return
    
    data = {
        "identifier": identifier,
        "new_email": new_email
    }
    
    result = make_request("PUT", "users/email", data=data, api_key=api_key)
    
    if "error" in result:
        print(f"❌ خطأ: {result['error']}")
    else:
        print(f"✅ {result.get('message', 'تم التحديث')}")

def update_user_password(api_key):
    """تغيير كلمة مرور مستخدم"""
    identifier = input("\nأدخل ID المستخدم أو الإيميل: ")
    
    # جلب جميع المستخدمين
    result = make_request("GET", "users", api_key=api_key)
    if "error" in result:
        print(f"❌ خطأ: {result['error']}")
        return
    
    users = result if isinstance(result, list) else []
    user = find_user(users, identifier)
    
    if not user:
        print("❌ المستخدم غير موجود")
        return
    
    print(f"المستخدم المحدد: {user['full_name']} ({user['email']})")
    new_password = getpass("أدخل كلمة المرور الجديدة: ")
    
    if not new_password:
        print("❌ كلمة المرور لا يمكن أن تكون فارغة")
        return
    
    if len(new_password) < 6:
        print("❌ كلمة المرور يجب أن تكون 6 أحرف على الأقل")
        return
    
    confirm_password = getpass("أعد إدخال كلمة المرور الجديدة: ")
    
    if new_password != confirm_password:
        print("❌ كلمات المرور غير متطابقة")
        return
    
    confirm = input("تأكيد تغيير كلمة المرور؟ (n/y): ")
    if confirm.lower() != 'y':
        print("❌ تم إلغاء العملية")
        return
    
    data = {
        "identifier": identifier,
        "new_password": new_password
    }
    
    result = make_request("PUT", "users/password", data=data, api_key=api_key)
    
    if "error" in result:
        print(f"❌ خطأ: {result['error']}")
    else:
        print(f"✅ {result.get('message', 'تم التحديث')}")

def create_admin(api_key):
    """إنشاء حساب أدمن جديد"""
    print("\n" + "=" * 60)
    print("إنشاء حساب أدمن جديد")
    print("=" * 60)
    
    full_name = input("أدخل الاسم الكامل: ")
    if not full_name:
        print("❌ الاسم لا يمكن أن يكون فارغاً")
        return
    
    email = input("أدخل الإيميل: ")
    if not email:
        print("❌ الإيميل لا يمكن أن يكون فارغاً")
        return
    
    if "@" not in email or "." not in email:
        print("❌ الإيميل غير صحيح")
        return
    
    password = getpass("أدخل كلمة المرور: ")
    if not password:
        print("❌ كلمة المرور لا يمكن أن تكون فارغة")
        return
    
    if len(password) < 8:
        print("❌ كلمة المرور يجب أن تكون 8 أحرف على الأقل")
        return
    
    confirm_password = getpass("أعد إدخال كلمة المرور: ")
    if password != confirm_password:
        print("❌ كلمات المرور غير متطابقة")
        return
    
    print("\nملخص الحساب الجديد:")
    print(f"الاسم: {full_name}")
    print(f"الإيميل: {email}")
    print("الدور: admin")
    print()
    
    confirm = input("تأكيد إنشاء حساب الأدمن؟ (n/y): ")
    if confirm.lower() != 'y':
        print("❌ تم إلغاء العملية")
        return
    
    data = {
        "full_name": full_name,
        "email": email,
        "password": password
    }
    
    result = make_request("POST", "users/admin", data=data, api_key=api_key)
    
    if "error" in result:
        print(f"❌ خطأ: {result['error']}")
    else:
        print(f"✅ {result.get('message', 'تم الإنشاء')}")
        print(f"🆔 معرف المستخدم: {result.get('user_id', 'غير معروف')}")

def delete_user(api_key, supabase_url, supabase_service_role_key):
    """حذف حساب"""
    print("\n" + "=" * 60)
    print("حذف حساب")
    print("=" * 60)
    
    identifier = input("أدخل ID المستخدم أو الإيميل للحذف: ")
    
    # جلب جميع المستخدمين
    result = make_request("GET", "users", api_key=api_key)
    if "error" in result:
        print(f"❌ خطأ: {result['error']}")
        return
    
    users = result if isinstance(result, list) else []
    user = find_user(users, identifier)
    
    if not user:
        print("❌ المستخدم غير موجود")
        return
    
    print("المستخدم المحدد:")
    print(f"الاسم: {user['full_name']}")
    print(f"الإيميل: {user['email']}")
    print(f"الدور: {user['role']}")
    print(f"ID: {user['id']}")
    print()
    
    confirm = input("هل أنت متأكد من حذف هذا الحساب؟ هذا الإجراء لا يمكن التراجع عنه (n/y): ")
    if confirm.lower() != 'y':
        print("❌ تم إلغاء عملية الحذف")
        return
    
    # حذف من Supabase Auth أولاً
    print("🔄 جاري حذف المستخدم من Supabase Auth...")
    delete_from_supabase_auth(user['id'], supabase_url, supabase_service_role_key)
    
    # استخدام ID المستخدم بدلاً من identifier
    data = {
        "identifier": user['id']  # استخدام ID مباشرة
    }
    
    result = make_request("DELETE", "users", data=data, api_key=api_key)
    
    if "error" in result:
        print(f"❌ خطأ: {result['error']}")
    else:
        print(f"✅ {result.get('message', 'تم الحذف')}")

def update_subscription(api_key):
    """تحديث مدة الاشتراك"""
    print("\n" + "=" * 60)
    print("تحديث مدة الاشتراك")
    print("=" * 60)
    
    identifier = input("أدخل ID المستخدم أو الإيميل: ")
    
    # جلب جميع المستخدمين
    result = make_request("GET", "users", api_key=api_key)
    if "error" in result:
        print(f"❌ خطأ: {result['error']}")
        return
    
    users = result if isinstance(result, list) else []
    user = find_user(users, identifier)
    
    if not user:
        print("❌ المستخدم غير موجود")
        return
    
    print(f"المستخدم المحدد: {user['full_name']} ({user['email']})")
    print()
    
    # عرض التوقيت العالمي الحالي
    utc_now = datetime.now(timezone.utc)
    print(f"🕐 التوقيت العالمي الحالي (UTC): {utc_now.strftime('%Y-%m-%d %H:%M:%S')}")
    print()
    
    # اختيار نوع التحديث
    print("نوع التحديث:")
    print("1. تغيير حالة الاشتراك فقط")
    print("2. إدخال مدة الاشتراك (بالأيام)")
    print("3. كل ما سبق")
    print()
    
    choice = input("اختر رقم (1-3): ")
    
    data = {
        "identifier": user['id']
    }
    
    if choice in ["1", "3"]:
        print("\nحالات الاشتراك:")
        print("1. active (نشط)")
        print("2. expired (منتهي)")
        print("3. canceled (ملغي)")
        status_choice = input("اختر حالة (1-3): ")
        
        status_map = {
            "1": "active",
            "2": "expired", 
            "3": "canceled"
        }
        
        if status_choice in status_map:
            data["subscription_status"] = status_map[status_choice]
    
    if choice in ["2", "3"]:
        # إدخال مباشر لعدد الأيام
        print("\nأدخل مدة الاشتراك (بالأيام):")
        print("مثال: 7 لسبعة أيام، 30 لشهر، 365 لسنة")
        
        try:
            days = int(input("عدد الأيام: "))
            if days <= 0:
                print("❌ يجب أن يكون عدد الأيام أكبر من صفر")
                return
            
            # حساب تاريخ الانتهاء بالتوقيت العالمي
            expires_at = utc_now + timedelta(days=days)
            data["expires_at"] = expires_at.isoformat()
            
            print(f"\n📅 تاريخ الانتهاء المحسوب (UTC): {expires_at.strftime('%Y-%m-%d %H:%M:%S')}")
            
        except ValueError:
            print("❌ قيمة غير صحيحة، أدخل رقماً")
            return
    
    if not data.get("subscription_status") and not data.get("expires_at"):
        print("❌ لم يتم تحديد أي تحديث")
        return
    
    print("\nملخص التحديث:")
    if data.get("subscription_status"):
        print(f"حالة الاشتراك: {data['subscription_status']}")
    if data.get("expires_at"):
        # تحويل ISO format إلى تنسيق مقروء
        try:
            expires_dt = datetime.fromisoformat(data['expires_at'].replace('Z', '+00:00'))
            print(f"تاريخ الانتهاء (UTC): {expires_dt.strftime('%Y-%m-%d %H:%M:%S')}")
        except:
            print(f"تاريخ الانتهاء: {data['expires_at']}")
    print()
    
    confirm = input("تأكيد التحديث؟ (n/y): ")
    if confirm.lower() != 'y':
        print("❌ تم إلغاء العملية")
        return
    
    result = make_request("PUT", "users/subscription", data=data, api_key=api_key)
    
    if "error" in result:
        print(f"❌ خطأ: {result['error']}")
    else:
        print(f"✅ {result.get('message', 'تم التحديث بنجاح')}")

def main():
    """الوظيفة الرئيسية"""
    global DEFAULT_API_URL
    api_url, api_key, supabase_url, supabase_service_role_key = load_config()
    DEFAULT_API_URL = api_url
    
    print("=" * 60)
    print("سكربت تغيير بيانات المستخدم - WorkTrack")
    print("=" * 60)
    print(f"API URL: {api_url}")
    if supabase_url:
        print(f"Supabase URL: {supabase_url}")
        print("✅ Supabase Auth deletion enabled")
    else:
        print("⚠️ Supabase Auth deletion disabled (credentials not configured)")
    print()
    
    while True:
        print("الخيارات:")
        print("1. عرض جميع المستخدمين")
        print("2. تغيير إيميل مستخدم")
        print("3. تغيير كلمة مرور مستخدم")
        print("4. تحديث مدة الاشتراك")
        print("5. إنشاء حساب أدمن جديد")
        print("6. حذف حساب")
        print("7. خروج")
        print()
        
        choice = input("اختر رقم (1-7): ")
        
        if choice == "1":
            list_users(api_key)
        elif choice == "2":
            update_user_email(api_key, supabase_url, supabase_service_role_key)
        elif choice == "3":
            update_user_password(api_key)
        elif choice == "4":
            update_subscription(api_key)
        elif choice == "5":
            create_admin(api_key)
        elif choice == "6":
            delete_user(api_key, supabase_url, supabase_service_role_key)
        elif choice == "7":
            print("👋 خروج...")
            break
        else:
            print("❌ اختيار غير صحيح")
        
        print()

if __name__ == "__main__":
    try:
        main()
    except KeyboardInterrupt:
        print("\n👋 تم إيقاف السكربت")
        sys.exit(0)
