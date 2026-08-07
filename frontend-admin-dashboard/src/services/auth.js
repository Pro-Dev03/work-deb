import api from './api'

export async function login(email, password) {
  try {
    console.log('📤 محاولة تسجيل الدخول:', { email })

    const { data } = await api.post('/auth/login', {
      email: email.trim(),
      password: password.trim()
    })

    console.log('✅ تم تسجيل الدخول بنجاح:', data)

    // تخزين بيانات المستخدم والتوكن في localStorage (الحل المزدوج للكوكيز)
    if (data.user) {
      localStorage.setItem('worktrack_admin_user', JSON.stringify(data.user))
      console.log('💾 تم تخزين بيانات المستخدم في localStorage')
    } else {
      console.warn('⚠️ لا توجد بيانات مستخدم في الاستجابة')
    }

    // تخزين التوكن في localStorage لاستخدامه في Authorization header
    if (data.access_token) {
      localStorage.setItem('worktrack_admin_token', data.access_token)
      console.log('💾 تم تخزين التوكن في localStorage')
    } else {
      console.warn('⚠️ لا يوجد توكن في الاستجابة')
    }

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
    localStorage.removeItem('worktrack_admin_token')
  }
}

export function currentUser() {
  const raw = localStorage.getItem('worktrack_admin_user')
  return raw ? JSON.parse(raw) : null
}

export function getToken() {
  // الحل المزدوج: التوكن في localStorage كنسخة احتياطية
  return localStorage.getItem('worktrack_admin_token')
}
