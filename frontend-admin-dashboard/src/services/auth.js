import api from './api'

export async function login(email, password) {
  try {
    console.log('📤 محاولة تسجيل الدخول:', { email })

    const { data } = await api.post('/auth/login', {
      email: email.trim(),
      password: password.trim()
    })

    console.log('✅ تم تسجيل الدخول بنجاح:', data)

    // تخزين بيانات المستخدم
    if (data.user) {
      localStorage.setItem('worktrack_admin_user', JSON.stringify(data.user))
      console.log('💾 تم تخزين بيانات المستخدم في localStorage')
    } else {
      console.warn('⚠️ لا توجد بيانات مستخدم في الاستجابة')
    }

    // تخزين CSRF token من استجابة تسجيل الدخول
    if (data.csrf_token) {
      localStorage.setItem('csrf_token', data.csrf_token)
      console.log('🔒 تم تخزين CSRF token من استجابة تسجيل الدخول')
    } else {
      console.warn('⚠️ لا يوجد CSRF token في الاستجابة')
    }

    // تخزين التوكن كنسخة احتياطية (لضمان العمل إذا فشلت cookies)
    if (data.access_token) {
      localStorage.setItem('worktrack_admin_token', data.access_token)
      console.log('💾 تم تخزين التوكن كنسخة احتياطية')
    }

    console.log('🔒 الأمان: نسخة احتياطية + cookies للتوثيق')

    return data
  } catch (error) {
    console.error('❌ فشل تسجيل الدخول:', error.response?.data || error.message)
    throw error
  }
}

export async function getCurrentUser() {
  const { data } = await api.get('/auth/me')
  return data
}

export async function logout() {
  try {
    // استدعاء endpoint تسجيل الخروج لحذف cookies
    await api.post('/auth/logout')
  } catch (error) {
    console.error('⚠️ خطأ في تسجيل الخروج:', error)
  } finally {
    // تنظيف localStorage في كل الأحوال
    localStorage.removeItem('worktrack_admin_user')
    localStorage.removeItem('csrf_token')
    localStorage.removeItem('worktrack_admin_token')
  }
}

export function currentUser() {
  const raw = localStorage.getItem('worktrack_admin_user')
  return raw ? JSON.parse(raw) : null
}

export function getToken() {
  // لم يعد مستخدماً - التوكنات في httpOnly cookies فقط للأمان
  return null
}
