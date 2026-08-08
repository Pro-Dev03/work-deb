# 🔒 دليل مراقبة الأمان - WorkTrack

## طرق مراقبة Security Logs

### الطريقة 1: Supabase Dashboard (الأفضل والأسهل)

#### الخطوات:
1. افتح https://supabase.com/dashboard
2. اختر مشروع WorkTrack
3. من القائمة الجانبية، اختر "Table Editor"
4. اختر جدول `security_logs`

#### المزايا:
- ✅ واجهة رسومية سهلة الاستخدام
- ✅ تصفية وترتيب البيانات
- ✅ لا تتطلب تعديلات على الكود
- ✅ تحديث في الوقت الفعلي

#### استعلامات SQL مفيدة:
```sql
-- المحاولات الفاشلة الأخيرة
SELECT * FROM security_logs 
WHERE success = false 
ORDER BY created_at DESC 
LIMIT 20;

-- الأنماط المشبوهة (محاولات كثيرة من نفس IP)
SELECT ip, COUNT(*) as attempts 
FROM security_logs 
WHERE success = false 
AND created_at > NOW() - INTERVAL '1 hour'
GROUP BY ip 
HAVING COUNT(*) > 5;

-- محاولات تسجيل دخول لنفس بريد إلكتروني
SELECT email, COUNT(*) as attempts, 
       MAX(created_at) as last_attempt
FROM security_logs 
WHERE event_type = 'auth_attempt' 
AND success = false 
AND created_at > NOW() - INTERVAL '24 hours'
GROUP BY email 
HAVING COUNT(*) > 3;

-- إحصائيات شاملة
SELECT 
    event_type,
    success,
    COUNT(*) as count,
    COUNT(DISTINCT ip) as unique_ips,
    COUNT(DISTINCT email) as unique_emails
FROM security_logs 
WHERE created_at > NOW() - INTERVAL '24 hours'
GROUP BY event_type, success
ORDER BY count DESC;
```

### الطريقة 2: API Endpoint

#### الاستخدام عبر curl:
```bash
# تسجيل الدخول
curl -X POST https://worktrack-api.onrender.com/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@devpro.com","password":"devproadmin"}' \
  -c cookies.txt

# الحصول على السجلات (آخر 50)
curl -X GET https://worktrack-api.onrender.com/api/v1/security/logs \
  -b cookies.txt
```

#### الاستخدام عبر JavaScript:
```javascript
// في لوحة المدير
const response = await fetch('/api/v1/security/logs', {
  credentials: 'include' // لإرسال cookies
});
const logs = await response.json();
console.log(logs);
```

### الطريقة 3: إضافة واجهة مراقبة في لوحة المدير

#### مثال Vue.js Component:
```vue
<template>
  <div class="security-logs">
    <h2>سجلات الأمان</h2>
    
    <div class="filters">
      <select v-model="filterType">
        <option value="">كل الأنواع</option>
        <option value="auth_attempt">محاولات المصادقة</option>
        <option value="suspicious_activity">نشاط مشبوه</option>
      </select>
      
      <select v-model="filterSuccess">
        <option value="">الكل</option>
        <option value="true">ناجح</option>
        <option value="false">فاشل</option>
      </select>
      
      <button @click="loadLogs">تحديث</button>
    </div>
    
    <div class="logs-table">
      <table>
        <thead>
          <tr>
            <th>الوقت</th>
            <th>النوع</th>
            <th>البريد</th>
            <th>IP</th>
            <th>الحالة</th>
            <th>السبب</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="log in filteredLogs" :key="log.id">
            <td>{{ formatDate(log.created_at) }}</td>
            <td>{{ log.event_type }}</td>
            <td>{{ log.email || '-' }}</td>
            <td>{{ log.ip }}</td>
            <td>{{ log.success ? '✅' : '❌' }}</td>
            <td>{{ log.reason || '-' }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script>
import api from '@/services/api'

export default {
  name: 'SecurityLogs',
  data() {
    return {
      logs: [],
      filterType: '',
      filterSuccess: ''
    }
  },
  computed: {
    filteredLogs() {
      return this.logs.filter(log => {
        if (this.filterType && log.event_type !== this.filterType) return false
        if (this.filterSuccess && String(log.success) !== this.filterSuccess) return false
        return true
      })
    }
  },
  methods: {
    async loadLogs() {
      try {
        const response = await api.get('/security/logs')
        this.logs = response.data
      } catch (error) {
        console.error('فشل تحميل السجلات:', error)
      }
    },
    formatDate(date) {
      return new Date(date).toLocaleString('ar-SA')
    }
  },
  mounted() {
    this.loadLogs()
    // تحديث كل 30 ثانية
    setInterval(this.loadLogs, 30000)
  }
}
</script>

<style scoped>
.security-logs {
  padding: 20px;
}

.filters {
  margin-bottom: 20px;
  display: flex;
  gap: 10px;
}

.logs-table table {
  width: 100%;
  border-collapse: collapse;
}

.logs-table th, .logs-table td {
  padding: 10px;
  border: 1px solid #ddd;
  text-align: right;
}

.logs-table th {
  background-color: #f5f5f5;
}
</style>
```

### الطريقة 4: إشعارات البريد الإلكتروني

يمكنك إضافة نظام إشعارات للمحاولات المشبوهة:

#### في Go Backend:
```go
// في auth_handler.go
func (h *AuthHandler) Login(c *gin.Context) {
    // ... كود موجود ...
    
    if isSuspicious {
        // إرسال إشعار بريد إلكتروني
        h.sendSecurityAlert(ip, email, reason)
        
        c.JSON(http.StatusTooManyRequests, gin.H{
            "error": "محاولات مشبوهة - تم حظر مؤقت",
        })
        return
    }
}

func (h *AuthHandler) sendSecurityAlert(ip, email, reason string) {
    // إرسال إشعار للمديرين
    // يمكن استخدام خدمة بريد إلكتروني مثل SendGrid
}
```

## 📊 استعلامات مراقبة متقدمة

### تحليل يومي:
```sql
SELECT 
    DATE(created_at) as date,
    COUNT(*) as total_attempts,
    COUNT(CASE WHEN success = true THEN 1 END) as successful_logins,
    COUNT(CASE WHEN success = false THEN 1 END) as failed_logins,
    COUNT(DISTINCT ip) as unique_ips
FROM security_logs 
WHERE created_at > NOW() - INTERVAL '7 days'
GROUP BY DATE(created_at)
ORDER BY date DESC;
```

### أكثر IP مشبوهة:
```sql
SELECT 
    ip,
    COUNT(*) as total_attempts,
    COUNT(CASE WHEN success = false THEN 1 END) as failed_attempts,
    MAX(created_at) as last_activity,
    array_agg(DISTINCT email) as attempted_emails
FROM security_logs 
WHERE created_at > NOW() - INTERVAL '24 hours'
GROUP BY ip 
HAVING COUNT(CASE WHEN success = false THEN 1 END) > 10
ORDER BY failed_attempts DESC;
```

### أنماط الوقت:
```sql
SELECT 
    EXTRACT(HOUR FROM created_at) as hour,
    COUNT(*) as attempts
FROM security_logs 
WHERE created_at > NOW() - INTERVAL '7 days'
GROUP BY hour
ORDER BY attempts DESC;
```

## 🚨 التنبيهات المقترحة

### قواعد التنبيه:
1. **أكثر من 5 محاولات فاشلة من نفس IP في 15 دقيقة**
2. **أكثر من 10 محاولات فاشلة لنفس بريد في ساعة**
3. **محاولات من IP من دول غير متوقعة**
4. **محاولات تسجيل دخول خارج أوقات العمل العادية**

### الإجراءات المقترحة:
1. **حظر تلقائي للـ IP المشبوهة**
2. **إرسال إشعار للمديرين**
3. **تطلب التحقق الثنائي (2FA)**
4. **تسجيل الحادثة في سجل تدقيق**

## 🛠️ أدوات إضافية

### لوحة مراقبة مخصصة:
يمكنك إنشاء لوحة مراقبة مخصصة باستخدام:
- Grafana + Supabase
- Metabase
- Retool
- AppSmith

### التكامل مع Slack/Discord:
```go
// إرسال إشعار Slack
func sendSlackAlert(message string) {
    webhookURL := "YOUR_SLACK_WEBHOOK"
    payload := map[string]string{"text": message}
    // إرسال إلى Slack
}
```

## 📈 إعداد Security Logger Service

### 1. تهيئة Security Logger في main.go:
```go
securityLogger := services.NewSecurityLogger(db)
```

### 2. تسجيل محاولات تسجيل الدخول:
```go
// في auth_handler.go
securityLogger.LogAuthAttempt(
    email, 
    phone, 
    c.ClientIP(), 
    c.Request.UserAgent(), 
    success, 
    reason,
)
```

### 3. التحقق من الأنماط المشبوهة:
```go
// قبل السماح بتسجيل الدخول
isSuspicious, reason := securityLogger.CheckForSuspiciousPatterns(ip, email)
if isSuspicious {
    // اتخذ إجراءات أمنية
}
```

### 4. الحصول على السجلات:
```go
// في security handler
logs, err := securityLogger.GetRecentSecurityLogs(50)
if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل جلب السجلات"})
    return
}
c.JSON(http.StatusOK, logs)
```

---

**التوصية:** ابدأ بـ Supabase Dashboard للمراقبة الأساسية، ثم أضف واجهة في لوحة المدير للمراقبة المستمرة.

**تم إنشاء هذا الدليل بواسطة Devin - مساعد الذكاء الاصطناعي**
**التاريخ:** 8 أغسطس 2026
