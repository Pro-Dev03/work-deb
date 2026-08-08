#!/usr/bin/env python3
# سكربت Python لإدارة مستخدمي WorkTrack عبر قاعدة البيانات مباشرة
# يتطلب: psycopg2-binary, bcrypt

import os
import sys
import psycopg2
import bcrypt
from getpass import getpass
from datetime import datetime, timedelta, timezone

def load_config():
    """قراءة الإعدادات من ملف .env"""
    database_url = ""
    
    env_files = ["backend/.env", ".env", "../backend/.env"]
    for env_file in env_files:
        if os.path.exists(env_file):
            with open(env_file, 'r') as f:
                for line in f:
                    line = line.strip()
                    if line.startswith("DATABASE_URL="):
                        database_url = line.split("=", 1)[1].strip()
                        break
            if database_url:
                break
    
    if not database_url:
        print("❌ DATABASE_URL غير موجود في ملف .env")
        sys.exit(1)
    
    return database_url

def get_db_connection(database_url):
    """الحصول على اتصال بقاعدة البيانات"""
    try:
        conn = psycopg2.connect(database_url)
        conn.autocommit = True
        return conn
    except Exception as e:
        print(f"❌ فشل الاتصال بقاعدة البيانات: {str(e)}")
        sys.exit(1)

def hash_password(password):
    """تشفير كلمة المرور باستخدام bcrypt مع مستوى أمان أعلى"""
    # استخدام bcrypt rounds أعلى لمزيد من الأمان (12 بدلاً من 10 الافتراضي)
    salt = bcrypt.gensalt(rounds=12)
    hashed = bcrypt.hashpw(password.encode('utf-8'), salt)
    return hashed.decode('utf-8')

def list_users(conn):
    """عرض جميع المستخدمين"""
    print("\nقائمة المستخدمين:")
    print("=" * 100)
    
    try:
        cursor = conn.cursor()
        
        # التحقق من الأعمدة الموجودة
        cursor.execute("""
            SELECT column_name 
            FROM information_schema.columns 
            WHERE table_name = 'users'
        """)
        columns = [row[0] for row in cursor.fetchall()]
        
        # بناء الاستعلام بناءً على الأعمدة الموجودة
        query = "SELECT id, full_name, email, role, is_active"
        if 'subscription_status' in columns:
            query += ", subscription_status"
        if 'expires_at' in columns:
            query += ", expires_at"
        query += " FROM users ORDER BY created_at DESC"
        
        cursor.execute(query)
        users = cursor.fetchall()
        
        if not users:
            print("لا يوجد مستخدمين")
            return
        
        # تحديد رؤوس الجدول بناءً على الأعمدة الموجودة
        headers = f"{'ID':<38} {'الاسم':<25} {'الإيميل':<30} {'الدور':<10} {'نشط':<5}"
        if 'subscription_status' in columns:
            headers += f" {'الاشتراك':<10}"
        if 'expires_at' in columns:
            headers += f" {'الانتهاء':<20}"
        
        print(headers)
        print("-" * 120)
        
        for user in users:
            user_id = user[0]
            full_name = user[1]
            email = user[2]
            role = user[3]
            is_active = user[4]
            
            line = f"{str(user_id):<38} {full_name:<25} {email:<30} {role:<10} {'✓' if is_active else '✗':<5}"
            
            index = 5
            if 'subscription_status' in columns:
                subscription_status = user[index]
                line += f" {subscription_status:<10}"
                index += 1
            if 'expires_at' in columns:
                expires_at = user[index]
                if expires_at:
                    expires_str = expires_at.strftime('%Y-%m-%d') if hasattr(expires_at, 'strftime') else str(expires_at)[:10]
                else:
                    expires_str = "N/A"
                line += f" {expires_str:<20}"
            
            print(line)
        
        print("=" * 100)
        return users
        
    except Exception as e:
        print(f"❌ خطأ: {str(e)}")
        return []

def find_user(conn, identifier):
    """البحث عن مستخدم"""
    try:
        cursor = conn.cursor()
        
        # التحقق من الأعمدة الموجودة
        cursor.execute("""
            SELECT column_name 
            FROM information_schema.columns 
            WHERE table_name = 'users'
        """)
        columns = [row[0] for row in cursor.fetchall()]
        
        # البحث بالـ email أولاً
        query = "SELECT id, full_name, email, role, is_active"
        if 'subscription_status' in columns:
            query += ", subscription_status"
        if 'expires_at' in columns:
            query += ", expires_at"
        query += " FROM users WHERE email = %s"
        
        cursor.execute(query, (identifier,))
        result = cursor.fetchone()
        
        if result:
            return result
        
        # إذا لم يتم العثور، البحث بالـ id
        query = "SELECT id, full_name, email, role, is_active"
        if 'subscription_status' in columns:
            query += ", subscription_status"
        if 'expires_at' in columns:
            query += ", expires_at"
        query += " FROM users WHERE id = %s"
        
        cursor.execute(query, (identifier,))
        return cursor.fetchone()
        
    except Exception as e:
        print(f"❌ خطأ في البحث: {str(e)}")
        return None

def update_user_email(conn):
    """تغيير إيميل مستخدم"""
    identifier = input("\nأدخل ID المستخدم أو الإيميل الحالي: ")
    
    user = find_user(conn, identifier)
    if not user:
        print("❌ المستخدم غير موجود")
        return
    
    user_id = user[0]
    full_name = user[1]
    email = user[2]
    print(f"المستخدم المحدد: {full_name} ({email})")
    
    new_email = input("أدخل الإيميل الجديد: ")
    if not new_email:
        print("❌ الإيميل لا يمكن أن يكون فارغاً")
        return
    
    if "@" not in new_email or "." not in new_email:
        print("❌ الإيميل غير صحيح")
        return
    
    confirm = input(f"تأكيد تغيير الإيميل إلى '{new_email}'؟ (n/y): ")
    if confirm.lower() != 'y':
        print("❌ تم إلغاء العملية")
        return
    
    try:
        cursor = conn.cursor()
        cursor.execute("UPDATE users SET email = %s WHERE id = %s", (new_email, user_id))
        print("✅ تم تحديث الإيميل بنجاح")
    except Exception as e:
        print(f"❌ خطأ: {str(e)}")

def update_user_password(conn):
    """تغيير كلمة مرور مستخدم"""
    identifier = input("\nأدخل ID المستخدم أو الإيميل: ")
    
    user = find_user(conn, identifier)
    if not user:
        print("❌ المستخدم غير موجود")
        return
    
    user_id = user[0]
    full_name = user[1]
    email = user[2]
    print(f"المستخدم المحدد: {full_name} ({email})")
    
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
    
    try:
        cursor = conn.cursor()
        hashed_password = hash_password(new_password)
        cursor.execute("UPDATE users SET password_hash = %s WHERE id = %s", (hashed_password, user_id))
        print("✅ تم تحديث كلمة المرور بنجاح")
    except Exception as e:
        print(f"❌ خطأ: {str(e)}")

def create_admin(conn):
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
    
    try:
        cursor = conn.cursor()
        hashed_password = hash_password(password)
        
        # التحقق من الأعمدة الموجودة
        cursor.execute("""
            SELECT column_name 
            FROM information_schema.columns 
            WHERE table_name = 'users'
        """)
        columns = [row[0] for row in cursor.fetchall()]
        
        # بناء الاستعلام بناءً على الأعمدة الموجودة
        if 'subscription_status' in columns:
            cursor.execute("""
                INSERT INTO users (full_name, email, password_hash, role, is_active, subscription_status, created_at)
                VALUES (%s, %s, %s, 'admin', true, 'active', NOW())
                RETURNING id
            """, (full_name, email, hashed_password))
        else:
            cursor.execute("""
                INSERT INTO users (full_name, email, password_hash, role, is_active, created_at)
                VALUES (%s, %s, %s, 'admin', true, NOW())
                RETURNING id
            """, (full_name, email, hashed_password))
        
        user_id = cursor.fetchone()[0]
        print(f"✅ تم إنشاء حساب الأدمن بنجاح")
        print(f"🆔 معرف المستخدم: {user_id}")
        
    except Exception as e:
        print(f"❌ خطأ: {str(e)}")

def delete_user(conn):
    """حذف حساب"""
    print("\n" + "=" * 60)
    print("حذف حساب")
    print("=" * 60)
    
    identifier = input("أدخل ID المستخدم أو الإيميل للحذف: ")
    
    user = find_user(conn, identifier)
    if not user:
        print("❌ المستخدم غير موجود")
        return
    
    user_id = user[0]
    full_name = user[1]
    email = user[2]
    role = user[3]
    
    print("المستخدم المحدد:")
    print(f"الاسم: {full_name}")
    print(f"الإيميل: {email}")
    print(f"الدور: {role}")
    print(f"ID: {user_id}")
    print()
    
    confirm = input("هل أنت متأكد من حذف هذا الحساب؟ هذا الإجراء لا يمكن التراجع عنه (n/y): ")
    if confirm.lower() != 'y':
        print("❌ تم إلغاء عملية الحذف")
        return
    
    try:
        cursor = conn.cursor()
        cursor.execute("DELETE FROM users WHERE id = %s", (user_id,))
        print("✅ تم حذف الحساب بنجاح")
    except Exception as e:
        print(f"❌ خطأ: {str(e)}")

def update_subscription(conn):
    """تحديث مدة الاشتراك"""
    print("\n" + "=" * 60)
    print("تحديث مدة الاشتراك")
    print("=" * 60)
    
    identifier = input("أدخل ID المستخدم أو الإيميل: ")
    
    user = find_user(conn, identifier)
    if not user:
        print("❌ المستخدم غير موجود")
        return
    
    user_id = user[0]
    full_name = user[1]
    email = user[2]
    print(f"المستخدم المحدد: {full_name} ({email})")
    print()
    
    # التحقق من الأعمدة الموجودة
    cursor = conn.cursor()
    cursor.execute("""
        SELECT column_name 
        FROM information_schema.columns 
        WHERE table_name = 'users'
    """)
    columns = [row[0] for row in cursor.fetchall()]
    
    if 'subscription_status' not in columns and 'expires_at' not in columns:
        print("❌ عمود الاشتراك غير موجود في قاعدة البيانات")
        return
    
    # عرض التوقيت العالمي الحالي
    utc_now = datetime.now(timezone.utc)
    print(f"🕐 التوقيت العالمي الحالي (UTC): {utc_now.strftime('%Y-%m-%d %H:%M:%S')}")
    print()
    
    # اختيار نوع التحديث
    print("نوع التحديث:")
    options = []
    if 'subscription_status' in columns:
        print("1. تغيير حالة الاشتراك فقط")
        options.append("1")
    if 'expires_at' in columns:
        print("2. إدخال مدة الاشتراك (بالأيام)" if 'subscription_status' not in columns else "2. إدخال مدة الاشتراك (بالأيام)")
        options.append("2")
    if 'subscription_status' in columns and 'expires_at' in columns:
        print("3. كل ما سبق")
        options.append("3")
    print()
    
    choice = input(f"اختر رقم ({', '.join(options)}): ")
    
    updates = {}
    
    if choice in ["1", "3"] and 'subscription_status' in columns:
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
            updates["subscription_status"] = status_map[status_choice]
    
    if choice in ["2", "3"] and 'expires_at' in columns:
        print("\nأدخل مدة الاشتراك (بالأيام):")
        print("مثال: 7 لسبعة أيام، 30 لشهر، 365 لسنة")
        
        try:
            days = int(input("عدد الأيام: "))
            if days <= 0:
                print("❌ يجب أن يكون عدد الأيام أكبر من صفر")
                return
            
            expires_at = utc_now + timedelta(days=days)
            updates["expires_at"] = expires_at
            
            print(f"\n📅 تاريخ الانتهاء المحسوب (UTC): {expires_at.strftime('%Y-%m-%d %H:%M:%S')}")
            
        except ValueError:
            print("❌ قيمة غير صحيحة، أدخل رقماً")
            return
    
    if not updates:
        print("❌ لم يتم تحديد أي تحديث")
        return
    
    print("\nملخص التحديث:")
    if updates.get("subscription_status"):
        print(f"حالة الاشتراك: {updates['subscription_status']}")
    if updates.get("expires_at"):
        print(f"تاريخ الانتهاء (UTC): {updates['expires_at'].strftime('%Y-%m-%d %H:%M:%S')}")
    print()
    
    confirm = input("تأكيد التحديث؟ (n/y): ")
    if confirm.lower() != 'y':
        print("❌ تم إلغاء العملية")
        return
    
    try:
        cursor = conn.cursor()
        
        if updates.get("subscription_status") and updates.get("expires_at"):
            cursor.execute("""
                UPDATE users 
                SET subscription_status = %s, expires_at = %s 
                WHERE id = %s
            """, (updates["subscription_status"], updates["expires_at"], user_id))
        elif updates.get("subscription_status"):
            cursor.execute("""
                UPDATE users 
                SET subscription_status = %s 
                WHERE id = %s
            """, (updates["subscription_status"], user_id))
        elif updates.get("expires_at"):
            cursor.execute("""
                UPDATE users 
                SET expires_at = %s 
                WHERE id = %s
            """, (updates["expires_at"], user_id))
        
        print("✅ تم تحديث الاشتراك بنجاح")
        
    except Exception as e:
        print(f"❌ خطأ: {str(e)}")

def main():
    """الوظيفة الرئيسية"""
    database_url = load_config()
    
    print("=" * 60)
    print("سكربت إدارة مستخدمي WorkTrack - قاعدة البيانات")
    print("=" * 60)
    print(f"DATABASE URL: {database_url[:50]}...")
    print()
    
    conn = get_db_connection(database_url)
    
    try:
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
                list_users(conn)
            elif choice == "2":
                update_user_email(conn)
            elif choice == "3":
                update_user_password(conn)
            elif choice == "4":
                update_subscription(conn)
            elif choice == "5":
                create_admin(conn)
            elif choice == "6":
                delete_user(conn)
            elif choice == "7":
                print("👋 خروج...")
                break
            else:
                print("❌ اختيار غير صحيح")
            
            print()
            
    finally:
        conn.close()

if __name__ == "__main__":
    try:
        main()
    except KeyboardInterrupt:
        print("\n👋 تم إيقاف السكربت")
        sys.exit(0)
