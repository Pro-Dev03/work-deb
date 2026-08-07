import api from './api'

export async function login(email, password) {
  try {
    console.log('📤 محاولة تسجيل الدخول:', { email })
    
    const { data } = await api.post('/auth/login', { 
      email: email.trim(), 
      password: password.trim() 
    })
    
    console.log('✅ تم تسجيل الدخول بنجاح:', data)
    
    // تخزين بيانات المستخدم فقط (التوكنات في httpOnly cookies)
    if (data.user) {
      localStorage.setItem('worktrack_admin_user', JSON.stringify(data.user))
      console.log('💾 تم تخزين بيانات المستخدم في localStorage')
    } else {
      console.warn('⚠️ لا توجد بيانات مستخدم في الاستجابة')
    }
    
    // التحقق من أن الكوكيز تم تعيينها بشكل صحيح
    console.log('🍪 التحقق من الكوكيز بعد تسجيل الدخول')
    
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
  }
}

export function currentUser() {
  const raw = localStorage.getItem('worktrack_admin_user')
  return raw ? JSON.parse(raw) : null
}

export function getToken() {
  // لم يعد مستخدماً - التوكنات في httpOnly cookies
  return null
}
