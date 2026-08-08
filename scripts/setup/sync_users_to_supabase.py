#!/usr/bin/env python3
# سكربت مزامنة المستخدمين من قاعدة البيانات إلى Supabase Auth

import os
import requests
import json
from getpass import getpass

def load_config():
    """قراءة الإعدادات من ملف .env"""
    supabase_url = ""
    supabase_service_role_key = ""
    
    env_files = [".env", "../backend/.env", "backend/.env"]
    for env_file in env_files:
        if os.path.exists(env_file):
            with open(env_file, 'r') as f:
                for line in f:
                    line = line.strip()
                    if line.startswith("SUPABASE_URL="):
                        supabase_url = line.split("=", 1)[1].strip()
                    elif line.startswith("SUPABASE_SERVICE_ROLE_KEY="):
                        supabase_service_role_key = line.split("=", 1)[1].strip()
            break
    
    return supabase_url, supabase_service_role_key

def get_users_from_api():
    """جلب المستخدمين من API"""
    api_url = "https://worktrack-v2.onrender.com/api/v1/admin-script/users"
    api_key = "dev-admin-script-key-2024"
    
    headers = {
        "X-API-Key": api_key,
        "Content-Type": "application/json"
    }
    
    try:
        response = requests.get(api_url, headers=headers)
        if response.status_code == 200:
            return response.json()
        else:
            print(f"❌ خطأ في جلب المستخدمين: {response.status_code}")
            return []
    except Exception as e:
        print(f"❌ خطأ في الاتصال: {str(e)}")
        return []

def create_user_in_supabase(email, password, user_data, supabase_url, service_role_key):
    """إنشاء مستخدم في Supabase Auth"""
    auth_url = f"{supabase_url}/auth/v1/admin/users"
    headers = {
        "apikey": service_role_key,
        "Authorization": f"Bearer {service_role_key}",
        "Content-Type": "application/json"
    }
    
    user_payload = {
        "email": email,
        "password": password,
        "email_confirm": True,
        "user_metadata": {
            "full_name": user_data.get('full_name', ''),
            "role": user_data.get('role', 'customer')
        }
    }
    
    try:
        response = requests.post(auth_url, headers=headers, json=user_payload)
        if response.status_code >= 200 and response.status_code < 300:
            return response.json(), True
        else:
            return response.text, False
    except Exception as e:
        return str(e), False

def main():
    """الوظيفة الرئيسية"""
    print("=" * 60)
    print("سكربت مزامنة المستخدمين إلى Supabase Auth")
    print("=" * 60)
    
    # تحميل الإعدادات
    supabase_url, supabase_service_role_key = load_config()
    
    if not supabase_url or not supabase_service_role_key:
        print("❌ Supabase credentials not configured")
        print("Please add SUPABASE_URL and SUPABASE_SERVICE_ROLE_KEY to .env file")
        return
    
    print(f"Supabase URL: {supabase_url}")
    print()
    
    # جلب المستخدمين من API
    print("🔄 جاري جلب المستخدمين من قاعدة البيانات...")
    users = get_users_from_api()
    
    if not users:
        print("❌ لا يوجد مستخدمين في قاعدة البيانات")
        return
    
    print(f"✅ تم جلب {len(users)} مستخدم")
    print()
    
    # عرض المستخدمين
    print("قائمة المستخدمين:")
    print(f"{'الاسم':<30} {'الإيميل':<30} {'الدور':<10}")
    print("-" * 70)
    for user in users:
        print(f"{user.get('full_name', ''):<30} {user.get('email', ''):<30} {user.get('role', ''):<10}")
    print()
    
    # طلب كلمة مرور افتراضية
    print("⚠️ سيتم استخدام كلمة مرور افتراضية لجميع المستخدمين")
    default_password = getpass("أدخل كلمة مرور افتراضية: ")
    
    if not default_password or len(default_password) < 6:
        print("❌ كلمة المرور يجب أن تكون 6 أحرف على الأقل")
        return
    
    print()
    confirm = input(f"هل تريد مزامنة {len(users)} مستخدم إلى Supabase Auth؟ (n/y): ")
    if confirm.lower() != 'y':
        print("❌ تم إلغاء العملية")
        return
    
    # مزامنة المستخدمين
    print("🔄 جاري مزامنة المستخدمين...")
    success_count = 0
    error_count = 0
    
    for user in users:
        email = user.get('email')
        if not email:
            print(f"⚠️ المستخدم {user.get('full_name', 'Unknown')} لا يمتلك إيميل - تخطي")
            error_count += 1
            continue
        
        print(f"🔄 جاري مزامنة: {email}...")
        result, success = create_user_in_supabase(
            email, 
            default_password, 
            user, 
            supabase_url, 
            supabase_service_role_key
        )
        
        if success:
            print(f"✅ تم إنشاء: {email}")
            success_count += 1
        else:
            print(f"❌ فشل: {email} - {result}")
            error_count += 1
    
    print()
    print("=" * 60)
    print("ملخص المزامنة:")
    print(f"✅ تمت المزامنة بنجاح: {success_count} مستخدم")
    print(f"❌ فشلت المزامنة: {error_count} مستخدم")
    print("=" * 60)
    
    if success_count > 0:
        print()
        print("⚠️ مهم: يجب على المستخدمين تغيير كلمات المرور الافتراضية عند أول تسجيل دخول")

if __name__ == "__main__":
    try:
        main()
    except KeyboardInterrupt:
        print("\n👋 تم إيقاف السكربت")
    except Exception as e:
        print(f"❌ خطأ: {str(e)}")