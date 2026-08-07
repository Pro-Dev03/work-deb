import axios from 'axios'
import { cacheService } from './cache'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1',
  withCredentials: true, // إرسال cookies تلقائياً
})

// قائمة نقاط النهاية التي يمكن تخزينها مؤقتاً
const CACHEABLE_ENDPOINTS = [
  '/worksites',
  '/admin/employees',
  '/reports/daily-summary'
]

// نقاط النهاية التي لا يجب تخزينها مؤقتاً (مثل المصادقة)
const NON_CACHEABLE_ENDPOINTS = [
  '/auth/login',
  '/auth/logout',
  '/auth/refresh',
  '/auth/me'
]

// دالة مساعدة لتوليد مفتاح الـ cache
function getCacheKey(url, params) {
  return `${url}${JSON.stringify(params)}`
}

// إضافة interceptor للطلبات لاستخدام الـ cache
api.interceptors.request.use((config) => {
  const cacheKey = getCacheKey(config.url, config.params)
  
  // عدم تطبيق الـ cache على نقاط النهاية الحساسة
  if (NON_CACHEABLE_ENDPOINTS.some(endpoint => config.url.includes(endpoint))) {
    return config
  }
  
  // التحقق من وجود البيانات في الـ cache للطلبات GET القابلة للتخزين
  if (config.method === 'get' && CACHEABLE_ENDPOINTS.some(endpoint => config.url.includes(endpoint))) {
    const cachedData = cacheService.get(cacheKey)
    if (cachedData) {
      config.adapter = () => Promise.resolve({
        data: cachedData,
        status: 200,
        statusText: 'OK',
        headers: {},
        config
      })
    }
  }
  
  return config
})

// إضافة interceptor للاستجابات لتخزين البيانات في الـ cache والتعامل مع الأخطاء
api.interceptors.response.use(
  (response) => {
    const cacheKey = getCacheKey(response.config.url, response.config.params)
    
    // عدم تخزين نقاط النهاية الحساسة في الـ cache
    if (NON_CACHEABLE_ENDPOINTS.some(endpoint => response.config.url.includes(endpoint))) {
      return response
    }
    
    // تخزين الاستجابة في الـ cache للطلبات GET القابلة للتخزين
    if (response.config.method === 'get' && CACHEABLE_ENDPOINTS.some(endpoint => response.config.url.includes(endpoint))) {
      cacheService.set(cacheKey, response.data)
    }
    
    return response
  },
  async (error) => {
    const originalRequest = error.config

    // إذا كان الخطأ 401 - تنظيف البيانات والتوجيه لتسجيل الدخول مباشرة
    if (error.response?.status === 401) {
      console.error('❌ خطأ 401 - الجلسة منتهية أو غير صالحة')
      
      // تنظيف localStorage
      localStorage.removeItem('worktrack_admin_user')
      
      // إعادة توجيه لصفحة تسجيل الدخول إذا لم نكن فيها
      if (window.location.pathname !== '/login') {
        window.location.href = '/login'
      }
      
      return Promise.reject(error)
    }

    return Promise.reject(error)
  }
)

export default api

// دالة مساعدة لمسح الـ cache عند تغيير البيانات
export function clearCacheForEndpoint(endpoint) {
  // مسح جميع الـ cache keys التي تحتوي على الـ endpoint المحدد
  for (const key of cacheService.cache.keys()) {
    if (key.includes(endpoint)) {
      cacheService.clear(key)
    }
  }
}
