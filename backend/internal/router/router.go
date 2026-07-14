package router

import (
	"database/sql"

	"worktrack/backend/internal/config"
	"worktrack/backend/internal/handlers"
	"worktrack/backend/internal/i18n"
	"worktrack/backend/internal/middleware"
	"worktrack/backend/internal/services"

	"github.com/gin-gonic/gin"
)

func Setup(db *sql.DB, cfg *config.Config) *gin.Engine {
	r := gin.Default()

	r.Use(middleware.CORSMiddleware(cfg.AllowedOrigin))
	r.Use(middleware.RateLimiter())

	authService := services.NewAuthService(cfg.JWTSecret)
	attendanceService := services.NewAttendanceService(db)
	notificationService := services.NewNotificationService(db)
	geocodingService := services.NewGeocodingService(cfg.GeoapifyKey)

	authHandler := handlers.NewAuthHandler(db, authService)
	attendanceHandler := handlers.NewAttendanceHandler(attendanceService, notificationService)
	worksiteHandler := handlers.NewWorksiteHandler(db)
	reportHandler := handlers.NewReportHandler(db)
	notificationHandler := handlers.NewNotificationHandler(db)
	serviceHandler := handlers.NewServiceHandler(db)
	locationHandler := handlers.NewLocationHandler(db)
	geocodingHandler := handlers.NewGeocodingHandler(geocodingService)

	r.GET("/health", func(c *gin.Context) {
		lang := i18n.Detect(c)
		c.JSON(200, gin.H{"status": "ok", "message": i18n.T(lang, "msg_health_ok")})
	})

	api := r.Group("/api/v1")

	// مسارات المصادقة (بدون توكن)
	api.POST("/auth/login", authHandler.Login)
	api.POST("/auth/phone-login", authHandler.PhoneLogin)
	api.GET("/geocode/autocomplete", geocodingHandler.Autocomplete)

	// المسارات المحمية (تتطلب توكن)
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(authService), middleware.SubscriptionMiddleware(db))
	{
		// المستخدم
		protected.GET("/auth/me", authHandler.Me)
		protected.GET("/auth/device", authHandler.GetDeviceInfo)

		// نقاط العمل
		protected.GET("/worksites", worksiteHandler.List)
		protected.GET("/worksites/available", worksiteHandler.GetAvailableWorksites)

		// الحضور
		protected.POST("/attendance/check-in", attendanceHandler.CheckIn)
		protected.POST("/attendance/check-out", attendanceHandler.CheckOut)
		protected.GET("/attendance/current", attendanceHandler.GetCurrentAttendance)
		protected.GET("/attendance/summary", attendanceHandler.GetAttendanceSummary)
		protected.GET("/attendance/all-summary", attendanceHandler.GetAllAttendanceSummary)

		// الموقع
		protected.GET("/location/active", locationHandler.GetActiveEmployees)
		protected.GET("/location/track/:id", locationHandler.GetEmployeeTrack)
		protected.GET("/location/security/:id", locationHandler.GetEmployeeSecurityNotes)
		protected.GET("/location/logs", locationHandler.GetLocationLogs)
		protected.POST("/location/update", locationHandler.UpdateLocation)

		// الإشعارات
		protected.GET("/notifications", notificationHandler.List)

		// طلبات الخدمة
		protected.GET("/service/requests", serviceHandler.ListRequests)
		protected.POST("/service/requests", serviceHandler.CreateRequest)
		protected.PUT("/service/requests/:id/status", serviceHandler.UpdateStatus)

		// مسارات المدير (تتطلب دور admin)
		admin := protected.Group("")
		admin.Use(middleware.RequireRole("admin"))
		{
			// إدارة الموظفين
			admin.POST("/auth/employee-phone", authHandler.CreateEmployeePhone)
			admin.GET("/admin/employees", authHandler.ListEmployees)
			admin.DELETE("/admin/employees/:id", authHandler.DeleteEmployee)
			admin.POST("/admin/reset-device", authHandler.ResetDevice)

			// نقاط العمل
			admin.POST("/worksites", worksiteHandler.Create)
			admin.DELETE("/worksites/:id", worksiteHandler.Delete)
			admin.POST("/worksites/assign", worksiteHandler.AssignEmployee)
			admin.GET("/worksites/employees", worksiteHandler.GetAvailableEmployees)

			// التقارير - مسارات جديدة
			admin.GET("/reports/daily-summary", reportHandler.DailySummary)
			admin.GET("/reports/pending-employees", reportHandler.GetPendingEmployees)
			admin.GET("/reports/completed-employees", reportHandler.GetCompletedEmployees)
		}
	}

	return r
}
