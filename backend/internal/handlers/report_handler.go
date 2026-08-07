package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	DB *sql.DB
}

func NewReportHandler(db *sql.DB) *ReportHandler {
	return &ReportHandler{DB: db}
}

// DailySummary - إحصائيات عامة
func (h *ReportHandler) DailySummary(c *gin.Context) {
	var completed, inProgress, pending, late, totalEmployees, waitingEmployees, completedToday int

	// الحصول على معامل الفترة من الاستعلام
	period := c.DefaultQuery("period", "today")

	// تحديد فترة التاريخ حسب المعامل
	var dateCondition string
	switch period {
	case "week":
		dateCondition = "created_at >= CURRENT_DATE - INTERVAL '7 days'"
	case "month":
		dateCondition = "created_at >= CURRENT_DATE - INTERVAL '30 days'"
	default: // today
		dateCondition = "created_at::date = CURRENT_DATE"
	}

	// الموظفين الذين لم يبدؤوا العمل اليوم (قيد الانتظار)
	err := h.DB.QueryRow(`
		SELECT COUNT(DISTINCT u.id)
		FROM users u
		WHERE u.role = 'employee' 
		AND u.is_active = TRUE
		AND NOT EXISTS (
			SELECT 1 FROM attendance a 
			WHERE a.user_id = u.id 
			AND DATE(a.check_in_time) = CURRENT_DATE
		)
	`).Scan(&waitingEmployees)
	if err != nil {
		waitingEmployees = 0
	}

	// الموظفين الذين أكملوا عملهم اليوم (مكتمل)
	err = h.DB.QueryRow(`
		SELECT COUNT(DISTINCT u.id)
		FROM users u
		JOIN attendance a ON a.user_id = u.id
		WHERE u.role = 'employee' 
		AND u.is_active = TRUE
		AND DATE(a.check_in_time) = CURRENT_DATE
		AND a.status = 'completed'
		AND a.check_out_time IS NOT NULL
	`).Scan(&completedToday)
	if err != nil {
		completedToday = 0
	}

	// إحصائيات المهام - جلب المهام حسب الفترة المحددة
	taskQuery := fmt.Sprintf(`SELECT COUNT(*) FROM tasks WHERE status = $1 AND %s`, dateCondition)
	
	err = h.DB.QueryRow(taskQuery, "completed").Scan(&completed)
	if err != nil {
		completed = 0
	}
	
	err = h.DB.QueryRow(taskQuery, "in_progress").Scan(&inProgress)
	if err != nil {
		inProgress = 0
	}
	
	err = h.DB.QueryRow(taskQuery, "pending").Scan(&pending)
	if err != nil {
		pending = 0
	}
	
	err = h.DB.QueryRow(taskQuery, "late").Scan(&late)
	if err != nil {
		late = 0
	}
	
	err = h.DB.QueryRow(`SELECT COUNT(*) FROM users WHERE role = 'employee' AND is_active = TRUE`).Scan(&totalEmployees)
	if err != nil {
		totalEmployees = 0
	}

	c.JSON(http.StatusOK, gin.H{
		"completed_today":     completed,
		"in_progress":         inProgress,
		"pending":             pending,
		"late":                late,
		"total_employees":     totalEmployees,
		"waiting_employees":   waitingEmployees,
		"completed_employees": completedToday,
		"period":              period,
	})
}

// GetPendingEmployees - جلب الموظفين قيد الانتظار
func (h *ReportHandler) GetPendingEmployees(c *gin.Context) {
	rows, err := h.DB.Query(`
		SELECT 
			u.id,
			u.full_name,
			u.email,
			u.phone
		FROM users u
		WHERE u.role = 'employee'
		AND u.is_active = TRUE
		AND NOT EXISTS (
			SELECT 1 FROM attendance a 
			WHERE a.user_id = u.id 
			AND DATE(a.check_in_time) = CURRENT_DATE
		)
		ORDER BY u.full_name
	`)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل جلب الموظفين قيد الانتظار"})
		return
	}
	defer rows.Close()

	var employees []gin.H
	for rows.Next() {
		var id, fullName, email, phone string

		if err := rows.Scan(&id, &fullName, &email, &phone); err != nil {
			continue
		}

		employees = append(employees, gin.H{
			"id":         id,
			"full_name":  fullName,
			"email":      email,
			"phone":      phone,
		})
	}

	c.JSON(http.StatusOK, employees)
}

// GetCompletedEmployees - جلب الموظفين المكتملين اليوم
func (h *ReportHandler) GetCompletedEmployees(c *gin.Context) {
	rows, err := h.DB.Query(`
		SELECT 
			u.id,
			u.full_name,
			u.email,
			u.phone,
			COALESCE(w.name, a.worksite_name_for_history) as worksite_name,
			a.check_in_time,
			a.check_out_time,
			EXTRACT(EPOCH FROM (a.check_out_time - a.check_in_time)) / 3600 as hours_worked
		FROM users u
		JOIN attendance a ON a.user_id = u.id
		LEFT JOIN worksites w ON a.worksite_id = w.id
		WHERE u.role = 'employee'
		AND u.is_active = TRUE
		AND DATE(a.check_in_time) = CURRENT_DATE
		AND a.status = 'completed'
		AND a.check_out_time IS NOT NULL
		ORDER BY a.check_out_time DESC
	`)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل جلب الموظفين المكتملين"})
		return
	}
	defer rows.Close()

	var employees []gin.H
	for rows.Next() {
		var id, fullName, email, phone, worksiteName string
		var checkInTime, checkOutTime string
		var hoursWorked float64

		if err := rows.Scan(&id, &fullName, &email, &phone, &worksiteName, &checkInTime, &checkOutTime, &hoursWorked); err != nil {
			continue
		}

		employees = append(employees, gin.H{
			"id":             id,
			"full_name":      fullName,
			"email":          email,
			"phone":          phone,
			"worksite_name":  worksiteName,
			"check_in_time":  checkInTime,
			"check_out_time": checkOutTime,
			"hours_worked":   hoursWorked,
		})
	}

	c.JSON(http.StatusOK, employees)
}
