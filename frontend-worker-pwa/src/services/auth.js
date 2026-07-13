import api from './api'

export async function login(email, password) {
  try {
    console.log('📤 محاولة تسجيل الدخول (موظف):', { email })
    
    const { data } = await api.post('/auth/login', { 
      email: email.trim(), 
      password: password.trim() 
    })
    
    console.log('✅ تم تسجيل الدخول بنجاح (موظف)')
    
    localStorage.setItem('worktrack_token', data.token)
    localStorage.setItem('worktrack_user', JSON.stringify(data.user))
    
    return data
  } catch (error) {
    console.error('❌ فشل تسجيل الدخول (موظف):', error.response?.data || error.message)
    throw error
  }
}

export function logout() {
  localStorage.removeItem('worktrack_token')
  localStorage.removeItem('worktrack_user')
}

export function currentUser() {
  const raw = localStorage.getItem('worktrack_user')
  return raw ? JSON.parse(raw) : null
}
