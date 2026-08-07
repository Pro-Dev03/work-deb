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

	// إحصائيات المهام - جلب جميع المهام النشطة بدون فلتر تاريخ
	err = h.DB.QueryRow(`SELECT COUNT(*) FROM tasks WHERE status = $1`, "completed").Scan(&completed)
	if err != nil {
		completed = 0
	}
	
	err = h.DB.QueryRow(`SELECT COUNT(*) FROM tasks WHERE status = $1`, "in_progress").Scan(&inProgress)
	if err != nil {
		inProgress = 0
	}
	
	err = h.DB.QueryRow(`SELECT COUNT(*) FROM tasks WHERE status = $1`, "pending").Scan(&pending)
	if err != nil {
		pending = 0
	}
	
	err = h.DB.QueryRow(`SELECT COUNT(*) FROM tasks WHERE status = $1`, "late").Scan(&late)
	if err != nil {
		late = 0
	}
	
	err = h.DB.QueryRow(`SELECT COUNT(*) FROM users WHERE role = 'employee' AND is_active = TRUE`).Scan(&totalEmployees)
	if err != nil {
		totalEmployees = 0
	}

	// فحص إضافي: استعلام جميع الحالات الموجودة في جدول المهام
	var totalTasks int
	var taskStatuses []string
	statusRows, err := h.DB.Query(`SELECT status, COUNT(*) FROM tasks GROUP BY status`)
	if err == nil {
		defer statusRows.Close()
		for statusRows.Next() {
			var status string
			var count int
			statusRows.Scan(&status, &count)
			taskStatuses = append(taskStatuses, fmt.Sprintf("%s:%d", status, count))
			totalTasks += count
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"completed_today":     completed,
		"in_progress":         inProgress,
		"pending":             pending,
		"late":                late,
		"total_employees":     totalEmployees,
		"waiting_employees":   waitingEmployees,
		"completed_employees": completedToday,
		"_debug": gin.H{
			"tasks_completed": completed,
			"tasks_in_progress": inProgress,
			"tasks_pending": pending,
			"tasks_late": late,
			"total_tasks": totalTasks,
			"task_statuses": taskStatuses,
		},
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

// DiagnosticTasks - فحص تشخيصي لجدول المهام
func (h *ReportHandler) DiagnosticTasks(c *gin.Context) {
	type TaskStatus struct {
		Status string
		Count  int
	}
	
	var taskStatuses []TaskStatus
	var totalTasks int
	
	// فحص جميع الحالات في جدول المهام
	rows, err := h.DB.Query(`SELECT status, COUNT(*) FROM tasks GROUP BY status`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل الاستعلام", "details": err.Error()})
		return
	}
	defer rows.Close()
	
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			continue
		}
		taskStatuses = append(taskStatuses, TaskStatus{Status: status, Count: count})
		totalTasks += count
	}
	
	// فحص إجمالي المهام
	var allTasksCount int
	h.DB.QueryRow(`SELECT COUNT(*) FROM tasks`).Scan(&allTasksCount)
	
	// فحص عينة من المهام
	var sampleTasks []gin.H
	sampleRows, _ := h.DB.Query(`SELECT id, title, status, created_at FROM tasks LIMIT 5`)
	if sampleRows != nil {
		defer sampleRows.Close()
		for sampleRows.Next() {
			var id, title, status, createdAt string
			sampleRows.Scan(&id, &title, &status, &createdAt)
			sampleTasks = append(sampleTasks, gin.H{
				"id": id,
				"title": title,
				"status": status,
				"created_at": createdAt,
			})
		}
	}
	
	c.JSON(http.StatusOK, gin.H{
		"total_tasks": allTasksCount,
		"status_breakdown": taskStatuses,
		"sample_tasks": sampleTasks,
		"table_exists": len(taskStatuses) > 0 || allTasksCount > 0,
	})
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
