package handlers

import (
	"database/sql"
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
	var completed, inProgress, pending, totalEmployees, waitingEmployees, completedToday int

	// الموظفين الذين لم يبدؤوا العمل اليوم (قيد الانتظار)
	_ = h.DB.QueryRow(`
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

	// الموظفين الذين أكملوا عملهم اليوم (مكتمل)
	_ = h.DB.QueryRow(`
		SELECT COUNT(DISTINCT u.id)
		FROM users u
		JOIN attendance a ON a.user_id = u.id
		WHERE u.role = 'employee' 
		AND u.is_active = TRUE
		AND DATE(a.check_in_time) = CURRENT_DATE
		AND a.status = 'completed'
		AND a.check_out_time IS NOT NULL
	`).Scan(&completedToday)

	// إحصائيات المهام
	_ = h.DB.QueryRow(`SELECT COUNT(*) FROM tasks WHERE status = 'completed' AND created_at::date = CURRENT_DATE`).Scan(&completed)
	_ = h.DB.QueryRow(`SELECT COUNT(*) FROM tasks WHERE status = 'in_progress'`).Scan(&inProgress)
	_ = h.DB.QueryRow(`SELECT COUNT(*) FROM tasks WHERE status = 'pending'`).Scan(&pending)
	_ = h.DB.QueryRow(`SELECT COUNT(*) FROM users WHERE role = 'employee' AND is_active = TRUE`).Scan(&totalEmployees)

	c.JSON(http.StatusOK, gin.H{
		"completed_today":     completed,
		"in_progress":         inProgress,
		"pending":             pending,
		"total_employees":     totalEmployees,
		"waiting_employees":   waitingEmployees,
		"completed_employees": completedToday,
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
			w.name as worksite_name,
			a.check_in_time,
			a.check_out_time,
			EXTRACT(EPOCH FROM (a.check_out_time - a.check_in_time)) / 3600 as hours_worked
		FROM users u
		JOIN attendance a ON a.user_id = u.id
		JOIN worksites w ON a.worksite_id = w.id
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
