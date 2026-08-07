// خدمة caching بسيطة لتخزين البيانات مؤقتاً
class CacheService {
  constructor() {
    this.cache = new Map()
    this.defaultTTL = 5 * 60 * 1000 // 5 دقائق بالمللي ثانية
  }

  set(key, data, ttl = this.defaultTTL) {
    this.cache.set(key, {
      data,
      expiresAt: Date.now() + ttl
    })
  }

  get(key) {
    const item = this.cache.get(key)
    if (!item) return null

    // التحقق من انتهاء صلاحية البيانات
    if (Date.now() > item.expiresAt) {
      this.cache.delete(key)
      return null
    }

    return item.data
  }

  has(key) {
    return this.get(key) !== null
  }

  clear(key) {
    if (key) {
      this.cache.delete(key)
    } else {
      this.cache.clear()
    }
  }

  // مسح البيانات المنتهية الصلاحية
  clearExpired() {
    const now = Date.now()
    for (const [key, item] of this.cache.entries()) {
      if (now > item.expiresAt) {
        this.cache.delete(key)
      }
    }
  }
}

export const cacheService = new CacheService()
