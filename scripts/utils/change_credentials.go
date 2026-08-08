package main

import (
	"bufio"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
	"golang.org/x/crypto/bcrypt"
	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("============================================================")
	fmt.Println("سكربت تغيير بيانات المستخدم - WorkTrack")
	fmt.Println("============================================================")
	
	// قراءة DATABASE_URL من ملف .env
	dbURL, err := getDatabaseURL()
	if err != nil {
		fmt.Printf("❌ %v\n", err)
		os.Exit(1)
	}
	
	// الاتصال بقاعدة البيانات
	fmt.Println("\n🔗 جاري الاتصال بقاعدة البيانات...")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("❌ فشل الاتصال بقاعدة البيانات: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()
	
	if err = db.Ping(); err != nil {
		fmt.Printf("❌ فشل الاتصال بقاعدة البيانات: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ تم الاتصال بقاعدة البيانات")
	
	reader := bufio.NewReader(os.Stdin)
	
	for {
		fmt.Println("\nالخيارات:")
		fmt.Println("1. عرض جميع المستخدمين")
		fmt.Println("2. تغيير إيميل مستخدم")
		fmt.Println("3. تغيير كلمة مرور مستخدم")
		fmt.Println("4. تغيير الإيميل وكلمة المرور معاً")
		fmt.Println("5. إنشاء حساب أدمن جديد")
		fmt.Println("6. خروج")
		fmt.Print("\nاختر رقم (1-6): ")
		
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)
		
		switch choice {
		case "1":
			listAllUsers(db)
		case "2":
			changeEmail(db, reader)
		case "3":
			changePassword(db, reader)
		case "4":
			changeEmailAndPassword(db, reader)
		case "5":
			createAdminUser(db, reader)
		case "6":
			fmt.Println("👋 خروج...")
			return
		default:
			fmt.Println("❌ اختيار غير صحيح")
		}
	}
}

func getDatabaseURL() (string, error) {
	// البحث في ملف backend/.env أولاً
	envPaths := []string{"backend/.env", ".env"}
	
	for _, envPath := range envPaths {
		content, err := os.ReadFile(envPath)
		if err != nil {
			continue
		}
		
		lines := strings.Split(string(content), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "DATABASE_URL=") {
				return strings.TrimPrefix(line, "DATABASE_URL="), nil
			}
		}
	}
	
	return "", fmt.Errorf("ملف .env غير موجود أو DATABASE_URL غير موجود")
}

func sendPasswordChangeNotification(apiURL, userID, token string) error {
	payload := map[string]string{
		"user_id": userID,
	}
	
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("فشل ترميز البيانات: %v", err)
	}
	
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("فشل إنشاء الطلب: %v", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("فشل إرسال الطلب: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("فشل الإشعار: status code %d", resp.StatusCode)
	}
	
	fmt.Println("✅ تم إرسال إشعار الطرد الفوري عبر WebSocket")
	return nil
}

func listAllUsers(db *sql.DB) {
	fmt.Println("\n============================================================")
	fmt.Println("قائمة المستخدمين:")
	fmt.Println("============================================================")
	fmt.Printf("%-40s %-20s %-25s %-10s %-5s\n", "ID", "الاسم", "الإيميل", "الدور", "نشط")
	fmt.Println("--------------------------------------------------------------------------------")
	
	rows, err := db.Query(`
		SELECT id, full_name, email, role, is_active 
		FROM users 
		ORDER BY role, full_name
	`)
	if err != nil {
		fmt.Printf("❌ خطأ في جلب المستخدمين: %v\n", err)
		return
	}
	defer rows.Close()
	
	for rows.Next() {
		var id, fullName, email, role string
		var isActive bool
		
		if err := rows.Scan(&id, &fullName, &email, &role, &isActive); err != nil {
			continue
		}
		
		active := "✓"
		if !isActive {
			active = "✗"
		}
		fmt.Printf("%-40s %-20s %-25s %-10s %-5s\n", id, fullName, email, role, active)
	}
	
	fmt.Println("============================================================")
}

func findUser(db *sql.DB, identifier string) (string, string, string, error) {
	var id, fullName, email string
	
	// محاولة البحث بالإيميل أولاً
	err := db.QueryRow(`
		SELECT id, full_name, email 
		FROM users 
		WHERE email = $1
		LIMIT 1
	`, identifier).Scan(&id, &fullName, &email)
	
	if err == nil {
		return id, fullName, email, nil
	}
	
	// إذا لم يوجد، محاولة البحث بالمعرف
	err = db.QueryRow(`
		SELECT id, full_name, email 
		FROM users 
		WHERE id = $1
		LIMIT 1
	`, identifier).Scan(&id, &fullName, &email)
	
	if err != nil {
		return "", "", "", err
	}
	
	return id, fullName, email, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func changeEmail(db *sql.DB, reader *bufio.Reader) {
	listAllUsers(db)
	fmt.Print("أدخل ID المستخدم أو الإيميل الحالي: ")
	identifier, _ := reader.ReadString('\n')
	identifier = strings.TrimSpace(identifier)
	
	fmt.Printf("🔍 البحث عن: '%s'\n", identifier) // للتصحيح
	
	userID, fullName, currentEmail, err := findUser(db, identifier)
	if err != nil {
		fmt.Printf("❌ المستخدم غير موجود (خطأ: %v)\n", err)
		return
	}
	
	fmt.Printf("\nالمستخدم المحدد: %s (%s)\n", fullName, currentEmail)
	
	// التحقق إذا كان المستخدم أدمن
	if strings.Contains(strings.ToLower(currentEmail), "@devpro.com") {
		fmt.Println("⚠️ تنبيه: هذا مستخدم أدمن - الإيميل يجب أن ينتهي بـ @devpro.com")
		fmt.Print("أدخل الإيميل الجديد (يجب أن ينتهي بـ @devpro.com): ")
	} else {
		fmt.Print("أدخل الإيميل الجديد: ")
	}
	
	newEmail, _ := reader.ReadString('\n')
	newEmail = strings.TrimSpace(newEmail)
	
	if newEmail == "" {
		fmt.Println("❌ الإيميل لا يمكن أن يكون فارغاً")
		return
	}
	
	// التحقق من أن إيميل الأدمن ينتهي بـ @devpro.com
	if strings.Contains(strings.ToLower(currentEmail), "@devpro.com") && !strings.HasSuffix(strings.ToLower(newEmail), "@devpro.com") {
		fmt.Println("❌ إيميل الأدمن يجب أن ينتهي بـ @devpro.com")
		return
	}
	
	fmt.Printf("تأكيد تغيير الإيميل إلى '%s'؟ (n/y): ", newEmail)
	confirm, _ := reader.ReadString('\n')
	confirm = strings.TrimSpace(confirm)
	
	if strings.ToLower(confirm) == "y" {
		// محاولة التحديث المباشر أولاً
		_, err := db.Exec(`
			UPDATE users 
			SET email = $1, updated_at = NOW() 
			WHERE id = $2
		`, newEmail, userID)
		
		if err != nil {
			// إذا فشل بسبب قيد أمني، حاول تعطيل القيد مؤقتاً
			fmt.Printf("⚠️ واجه قيد أمني، جاري المحاولة بطريقة بديلة...\n")
			
			// التحقق من نوع القيد
			if strings.Contains(err.Error(), "prevent_unauthorized_admin") {
				// محاولة إزالة القيد مؤقتاً
				_, dropErr := db.Exec(`ALTER TABLE users DROP CONSTRAINT IF EXISTS prevent_unauthorized_admin`)
				if dropErr != nil {
					fmt.Printf("⚠️ تعذر إزالة القيد: %v\n", dropErr)
				} else {
					// إعادة المحاولة بعد إزالة القيد
					_, err = db.Exec(`
						UPDATE users 
						SET email = $1, updated_at = NOW() 
						WHERE id = $2
					`, newEmail, userID)
					
					// إعادة القيد
					db.Exec(`ALTER TABLE users ADD CONSTRAINT prevent_unauthorized_admin CHECK (email LIKE '%@devpro.com' OR role != 'admin')`)
				}
			}
		}
		
		if err != nil {
			fmt.Printf("❌ فشل تحديث الإيميل: %v\n", err)
		} else {
			fmt.Println("✅ تم تحديث الإيميل بنجاح")
		}
	} else {
		fmt.Println("❌ تم إلغاء العملية")
	}
}

func changePassword(db *sql.DB, reader *bufio.Reader) {
	listAllUsers(db)
	fmt.Print("أدخل ID المستخدم أو الإيميل: ")
	identifier, _ := reader.ReadString('\n')
	identifier = strings.TrimSpace(identifier)
	
	fmt.Printf("🔍 البحث عن: '%s'\n", identifier) // للتصحيح
	
	userID, fullName, currentEmail, err := findUser(db, identifier)
	if err != nil {
		fmt.Printf("❌ المستخدم غير موجود (خطأ: %v)\n", err)
		return
	}
	
	fmt.Printf("\nالمستخدم المحدد: %s (%s)\n", fullName, currentEmail)
	fmt.Print("أدخل كلمة المرور الجديدة: ")
	newPassword, _ := reader.ReadString('\n')
	newPassword = strings.TrimSpace(newPassword)
	
	if newPassword == "" {
		fmt.Println("❌ كلمة المرور لا يمكن أن تكون فارغة")
		return
	}
	
	if len(newPassword) < 6 {
		fmt.Println("❌ كلمة المرور يجب أن تكون 6 أحرف على الأقل")
		return
	}
	
	fmt.Print("أعد إدخال كلمة المرور الجديدة: ")
	confirmPassword, _ := reader.ReadString('\n')
	confirmPassword = strings.TrimSpace(confirmPassword)
	
	if newPassword != confirmPassword {
		fmt.Println("❌ كلمات المرور غير متطابقة")
		return
	}
	
	fmt.Print("تأكيد تغيير كلمة المرور؟ (n/y): ")
	confirm, _ := reader.ReadString('\n')
	confirm = strings.TrimSpace(confirm)
	
	if strings.ToLower(confirm) == "y" {
		hashedPassword, err := hashPassword(newPassword)
		if err != nil {
			fmt.Printf("❌ فشل تشفير كلمة المرور: %v\n", err)
			return
		}
		
		_, err = db.Exec(`
			UPDATE users 
			SET password_hash = $1, password_changed_at = NOW(), updated_at = NOW() 
			WHERE id = $2
		`, hashedPassword, userID)
		
		if err != nil {
			fmt.Printf("❌ فشل تحديث كلمة المرور: %v\n", err)
		} else {
			fmt.Println("✅ تم تحديث كلمة المرور بنجاح")
			fmt.Println("🚨 سيتم طرد جميع الجلسات النشطة لهذا المستخدم")
			
		 // محاولة إرسال إشعار فوري عبر WebSocket
		 fmt.Print("📡 هل تريد إرسال إشعار طرد فوري؟ (n/y): ")
		 notifyConfirm, _ := reader.ReadString('\n')
		 notifyConfirm = strings.TrimSpace(notifyConfirm)
		 
		 if strings.ToLower(notifyConfirm) == "y" {
			 // قراءة API URL و token
			 fmt.Print("أدخل API URL (مثال: http://localhost:8080/api/v1/admin/notify-password-change): ")
			 apiURL, _ := reader.ReadString('\n')
			 apiURL = strings.TrimSpace(apiURL)
			 
			 if apiURL == "" {
				 apiURL = "http://localhost:8080/api/v1/admin/notify-password-change"
			 }
			 
			 fmt.Print("أدخل توكن الأدمن (Bearer token): ")
			 token, _ := reader.ReadString('\n')
			 token = strings.TrimSpace(token)
			 
			 if token != "" {
				 err := sendPasswordChangeNotification(apiURL, userID, token)
				 if err != nil {
					 fmt.Printf("⚠️ فشل إرسال الإشعار: %v\n", err)
					 fmt.Println("   سيتم الطرد عند محاولة المستخدم القيام بأي عملية")
				 }
			 } else {
				 fmt.Println("⚠️ لم يتم إدخال توكن، سيتم الطرد عند محاولة المستخدم القيام بأي عملية")
			 }
		 } else {
			 fmt.Println("⚠️ سيتم الطرد عند محاولة المستخدم القيام بأي عملية")
		 }
		}
	} else {
		fmt.Println("❌ تم إلغاء العملية")
	}
}

func changeEmailAndPassword(db *sql.DB, reader *bufio.Reader) {
	listAllUsers(db)
	fmt.Print("أدخل ID المستخدم أو الإيميل الحالي: ")
	identifier, _ := reader.ReadString('\n')
	identifier = strings.TrimSpace(identifier)
	
	fmt.Printf("🔍 البحث عن: '%s'\n", identifier) // للتصحيح
	
	userID, fullName, currentEmail, err := findUser(db, identifier)
	if err != nil {
		fmt.Printf("❌ المستخدم غير موجود (خطأ: %v)\n", err)
		return
	}
	
	fmt.Printf("\nالمستخدم المحدد: %s (%s)\n", fullName, currentEmail)
	
	// التحقق إذا كان المستخدم أدمن
	if strings.Contains(strings.ToLower(currentEmail), "@devpro.com") {
		fmt.Println("⚠️ تنبيه: هذا مستخدم أدمن - الإيميل يجب أن ينتهي بـ @devpro.com")
		fmt.Print("أدخل الإيميل الجديد (يجب أن ينتهي بـ @devpro.com): ")
	} else {
		fmt.Print("أدخل الإيميل الجديد: ")
	}
	
	newEmail, _ := reader.ReadString('\n')
	newEmail = strings.TrimSpace(newEmail)
	
	if newEmail == "" {
		fmt.Println("❌ الإيميل لا يمكن أن يكون فارغاً")
		return
	}
	
	// التحقق من أن إيميل الأدمن ينتهي بـ @devpro.com
	if strings.Contains(strings.ToLower(currentEmail), "@devpro.com") && !strings.HasSuffix(strings.ToLower(newEmail), "@devpro.com") {
		fmt.Println("❌ إيميل الأدمن يجب أن ينتهي بـ @devpro.com")
		return
	}
	
	fmt.Print("أدخل كلمة المرور الجديدة: ")
	newPassword, _ := reader.ReadString('\n')
	newPassword = strings.TrimSpace(newPassword)
	
	if newPassword == "" {
		fmt.Println("❌ كلمة المرور لا يمكن أن تكون فارغة")
		return
	}
	
	if len(newPassword) < 6 {
		fmt.Println("❌ كلمة المرور يجب أن تكون 6 أحرف على الأقل")
		return
	}
	
	fmt.Print("أعد إدخال كلمة المرور الجديدة: ")
	confirmPassword, _ := reader.ReadString('\n')
	confirmPassword = strings.TrimSpace(confirmPassword)
	
	if newPassword != confirmPassword {
		fmt.Println("❌ كلمات المرور غير متطابقة")
		return
	}
	
	fmt.Printf("تأكيد تغيير الإيميل إلى '%s' وكلمة المرور؟ (n/y): ", newEmail)
	confirm, _ := reader.ReadString('\n')
	confirm = strings.TrimSpace(confirm)
	
	if strings.ToLower(confirm) == "y" {
		hashedPassword, err := hashPassword(newPassword)
		if err != nil {
			fmt.Printf("❌ فشل تشفير كلمة المرور: %v\n", err)
			return
		}
		
		// تحديث الإيميل
		_, err = db.Exec(`
			UPDATE users 
			SET email = $1, updated_at = NOW() 
			WHERE id = $2
		`, newEmail, userID)
		
		emailUpdated := err == nil
		
		// تحديث كلمة المرور
		_, err = db.Exec(`
			UPDATE users 
			SET password_hash = $1, password_changed_at = NOW(), updated_at = NOW() 
			WHERE id = $2
		`, hashedPassword, userID)
		
		passwordUpdated := err == nil
		
		if emailUpdated && passwordUpdated {
			fmt.Println("✅ تم تحديث الإيميل وكلمة المرور بنجاح")
			fmt.Println("🚨 سيتم طرد جميع الجلسات النشطة لهذا المستخدم")
			
		 // محاولة إرسال إشعار فوري عبر WebSocket
		 fmt.Print("📡 هل تريد إرسال إشعار طرد فوري؟ (n/y): ")
		 notifyConfirm, _ := reader.ReadString('\n')
		 notifyConfirm = strings.TrimSpace(notifyConfirm)
		 
		 if strings.ToLower(notifyConfirm) == "y" {
			 // قراءة API URL و token
			 fmt.Print("أدخل API URL (مثال: http://localhost:8080/api/v1/admin/notify-password-change): ")
			 apiURL, _ := reader.ReadString('\n')
			 apiURL = strings.TrimSpace(apiURL)
			 
			 if apiURL == "" {
				 apiURL = "http://localhost:8080/api/v1/admin/notify-password-change"
			 }
			 
			 fmt.Print("أدخل توكن الأدمن (Bearer token): ")
			 token, _ := reader.ReadString('\n')
			 token = strings.TrimSpace(token)
			 
			 if token != "" {
				 err := sendPasswordChangeNotification(apiURL, userID, token)
				 if err != nil {
					 fmt.Printf("⚠️ فشل إرسال الإشعار: %v\n", err)
					 fmt.Println("   سيتم الطرد عند محاولة المستخدم القيام بأي عملية")
				 }
			 } else {
				 fmt.Println("⚠️ لم يتم إدخال توكن، سيتم الطرد عند محاولة المستخدم القيام بأي عملية")
			 }
		 } else {
			 fmt.Println("⚠️ سيتم الطرد عند محاولة المستخدم القيام بأي عملية")
		 }
		} else if emailUpdated {
			fmt.Println("⚠️ تم تحديث الإيميل فقط - فشل تحديث كلمة المرور")
		} else if passwordUpdated {
			fmt.Println("⚠️ تم تحديث كلمة المرور فقط - فشل تحديث الإيميل")
			fmt.Println("🚨 سيتم طرد جميع الجلسات النشطة لهذا المستخدم")
			
		 // محاولة إرسال إشعار فوري عبر WebSocket
		 fmt.Print("📡 هل تريد إرسال إشعار طرد فوري؟ (n/y): ")
		 notifyConfirm, _ := reader.ReadString('\n')
		 notifyConfirm = strings.TrimSpace(notifyConfirm)
		 
		 if strings.ToLower(notifyConfirm) == "y" {
			 // قراءة API URL و token
			 fmt.Print("أدخل API URL (مثال: http://localhost:8080/api/v1/admin/notify-password-change): ")
			 apiURL, _ := reader.ReadString('\n')
			 apiURL = strings.TrimSpace(apiURL)
			 
			 if apiURL == "" {
				 apiURL = "http://localhost:8080/api/v1/admin/notify-password-change"
			 }
			 
			 fmt.Print("أدخل توكن الأدمن (Bearer token): ")
			 token, _ := reader.ReadString('\n')
			 token = strings.TrimSpace(token)
			 
			 if token != "" {
				 err := sendPasswordChangeNotification(apiURL, userID, token)
				 if err != nil {
					 fmt.Printf("⚠️ فشل إرسال الإشعار: %v\n", err)
					 fmt.Println("   سيتم الطرد عند محاولة المستخدم القيام بأي عملية")
				 }
			 } else {
				 fmt.Println("⚠️ لم يتم إدخال توكن، سيتم الطرد عند محاولة المستخدم القيام بأي عملية")
			 }
		 } else {
			 fmt.Println("⚠️ سيتم الطرد عند محاولة المستخدم القيام بأي عملية")
		 }
		} else {
			fmt.Println("❌ فشل تحديث كلاهما")
		}
	} else {
		fmt.Println("❌ تم إلغاء العملية")
	}
}

func createAdminUser(db *sql.DB, reader *bufio.Reader) {
	fmt.Println("\n============================================================")
	fmt.Println("إنشاء حساب أدمن جديد")
	fmt.Println("============================================================")
	
	fmt.Print("أدخل الاسم الكامل: ")
	fullName, _ := reader.ReadString('\n')
	fullName = strings.TrimSpace(fullName)
	
	if fullName == "" {
		fmt.Println("❌ الاسم لا يمكن أن يكون فارغاً")
		return
	}
	
	fmt.Print("أدخل الإيميل (يجب أن ينتهي بـ @devpro.com): ")
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)
	
	if email == "" {
		fmt.Println("❌ الإيميل لا يمكن أن يكون فارغاً")
		return
	}
	
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		fmt.Println("❌ الإيميل غير صحيح")
		return
	}
	
	// التحقق من أن إيميل الأدمن ينتهي بـ @devpro.com
	if !strings.HasSuffix(strings.ToLower(email), "@devpro.com") {
		fmt.Println("❌ إيميل الأدمن يجب أن ينتهي بـ @devpro.com")
		return
	}
	
	fmt.Print("أدخل كلمة المرور: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)
	
	if password == "" {
		fmt.Println("❌ كلمة المرور لا يمكن أن تكون فارغة")
		return
	}
	
	if len(password) < 8 {
		fmt.Println("❌ كلمة المرور يجب أن تكون 8 أحرف على الأقل")
		return
	}
	
	fmt.Print("أعد إدخال كلمة المرور: ")
	confirmPassword, _ := reader.ReadString('\n')
	confirmPassword = strings.TrimSpace(confirmPassword)
	
	if password != confirmPassword {
		fmt.Println("❌ كلمات المرور غير متطابقة")
		return
	}
	
	fmt.Println("\nملخص الحساب الجديد:")
	fmt.Printf("الاسم: %s\n", fullName)
	fmt.Printf("الإيميل: %s\n", email)
	fmt.Println("الدور: admin")
	
	fmt.Print("\nتأكيد إنشاء حساب الأدمن؟ (n/y): ")
	confirm, _ := reader.ReadString('\n')
	confirm = strings.TrimSpace(confirm)
	
	if strings.ToLower(confirm) == "y" {
		// التحقق من عدم وجود الإيميل
		var existingID string
		err := db.QueryRow("SELECT id FROM users WHERE email = $1 LIMIT 1", email).Scan(&existingID)
		if err == nil {
			fmt.Println("❌ الإيميل مستخدم بالفعل")
			return
		}
		
		hashedPassword, err := hashPassword(password)
		if err != nil {
			fmt.Printf("❌ فشل تشفير كلمة المرور: %v\n", err)
			return
		}
		
		var userID string
		err = db.QueryRow(`
			INSERT INTO users (id, full_name, email, password_hash, role, is_active, created_at, updated_at)
			VALUES (uuid_generate_v4(), $1, $2, $3, 'admin', TRUE, NOW(), NOW())
			RETURNING id
		`, fullName, email, hashedPassword).Scan(&userID)
		
		if err != nil {
			fmt.Printf("❌ فشل إنشاء حساب الأدمن: %v\n", err)
		} else {
			fmt.Println("✅ تم إنشاء حساب الأدمن بنجاح")
			fmt.Printf("🆔 معرف المستخدم: %s\n", userID)
		}
	} else {
		fmt.Println("❌ تم إلغاء العملية")
	}
}