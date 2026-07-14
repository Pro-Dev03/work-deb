import api from './api'

export async function login(email, password) {
  try {
    console.log('📤 محاولة تسجيل الدخول:', { email })
    
    const { data } = await api.post('/auth/login', { 
      email: email.trim(), 
      password: password.trim() 
    })
    
    console.log('✅ تم تسجيل الدخول بنجاح')
    
    localStorage.setItem('worktrack_admin_token', data.token)
    localStorage.setItem('worktrack_admin_user', JSON.stringify(data.user))
    
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

export function logout() {
  localStorage.removeItem('worktrack_admin_token')
  localStorage.removeItem('worktrack_admin_user')
}

export function currentUser() {
  const raw = localStorage.getItem('worktrack_admin_user')
  return raw ? JSON.parse(raw) : null
}

export function getToken() {
  return localStorage.getItem('worktrack_admin_token')
}
