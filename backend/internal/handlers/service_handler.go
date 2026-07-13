package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"worktrack/backend/internal/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ServiceHandler struct {
	DB *sql.DB
}

func NewServiceHandler(db *sql.DB) *ServiceHandler {
	return &ServiceHandler{DB: db}
}

type createServiceRequest struct {
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description"`
	Latitude    float64 `json:"latitude" binding:"required"`
	Longitude   float64 `json:"longitude" binding:"required"`
	Address     string  `json:"address"`
	Phone       string  `json:"phone"`
	Priority    string  `json:"priority"`
	Photos      []string `json:"photos"`
}

func (h *ServiceHandler) CreateRequest(c *gin.Context) {
	lang := i18n.Detect(c)
	userID, _ := c.Get("user_id")

	var req createServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.T(lang, "err_missing_fields")})
		return
	}

	id := uuid.NewString()
	_, err := h.DB.Exec(`
		INSERT INTO service_requests (
			id, client_id, title, description, 
			latitude, longitude, address, phone, 
			priority, status, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, 'pending', now(), now())`,
		id, userID, req.Title, req.Description,
		req.Latitude, req.Longitude, req.Address, req.Phone,
		req.Priority,
	)

	if err != nil {
		log.Printf("failed to create service request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.T(lang, "err_operation_failed")})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":      id,
		"message": "تم إرسال طلب الخدمة بنجاح، سيتم توجيهه لأقرب موظف",
	})
}

func (h *ServiceHandler) ListRequests(c *gin.Context) {
	lang := i18n.Detect(c)

	rows, err := h.DB.Query(`
		SELECT 
			sr.id, sr.title, sr.description, 
			sr.latitude, sr.longitude, sr.address, sr.phone,
			sr.status, sr.priority, sr.created_at,
			u.full_name as client_name, u.phone as client_phone
		FROM service_requests sr
		LEFT JOIN users u ON sr.client_id = u.id
		ORDER BY 
			CASE sr.priority 
				WHEN 'urgent' THEN 1 
				WHEN 'high' THEN 2 
				WHEN 'normal' THEN 3 
				WHEN 'low' THEN 4 
			END,
			sr.created_at DESC
	`)

	if err != nil {
		log.Printf("failed to fetch requests: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.T(lang, "err_operation_failed")})
		return
	}
	defer rows.Close()

	var requests []gin.H
	for rows.Next() {
		var id, title, description, address, phone, status, priority, clientName, clientPhone string
		var latitude, longitude float64
		var createdAt time.Time

		if err := rows.Scan(&id, &title, &description, &latitude, &longitude,
			&address, &phone, &status, &priority, &createdAt,
			&clientName, &clientPhone); err == nil {
			requests = append(requests, gin.H{
				"id":            id,
				"title":         title,
				"description":   description,
				"latitude":      latitude,
				"longitude":     longitude,
				"address":       address,
				"phone":         phone,
				"status":        status,
				"priority":      priority,
				"created_at":    createdAt,
				"client_name":   clientName,
				"client_phone":  clientPhone,
			})
		}
	}

	c.JSON(http.StatusOK, requests)
}

func (h *ServiceHandler) GetRequest(c *gin.Context) {
	lang := i18n.Detect(c)
	id := c.Param("id")

	var req struct {
		ID          string    `json:"id"`
		Title       string    `json:"title"`
		Description string    `json:"description"`
		Latitude    float64   `json:"latitude"`
		Longitude   float64   `json:"longitude"`
		Address     string    `json:"address"`
		Phone       string    `json:"phone"`
		Status      string    `json:"status"`
		Priority    string    `json:"priority"`
		CreatedAt   time.Time `json:"created_at"`
		ClientName  string    `json:"client_name"`
		ClientPhone string    `json:"client_phone"`
	}

	err := h.DB.QueryRow(`
		SELECT 
			sr.id, sr.title, sr.description, 
			sr.latitude, sr.longitude, sr.address, sr.phone,
			sr.status, sr.priority, sr.created_at,
			u.full_name, u.phone
		FROM service_requests sr
		LEFT JOIN users u ON sr.client_id = u.id
		WHERE sr.id = $1`, id,
	).Scan(&req.ID, &req.Title, &req.Description, &req.Latitude, &req.Longitude,
		&req.Address, &req.Phone, &req.Status, &req.Priority, &req.CreatedAt,
		&req.ClientName, &req.ClientPhone)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": i18n.T(lang, "err_operation_failed")})
		return
	}

	c.JSON(http.StatusOK, req)
}

type assignRequest struct {
	RequestID  string `json:"request_id" binding:"required"`
	EmployeeID string `json:"employee_id" binding:"required"`
	AdminNotes string `json:"admin_notes"`
}

func (h *ServiceHandler) AssignEmployee(c *gin.Context) {
	lang := i18n.Detect(c)
	adminID, _ := c.Get("user_id")

	var req assignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.T(lang, "err_missing_fields")})
		return
	}

	tx, err := h.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.T(lang, "err_operation_failed")})
		return
	}
	defer tx.Rollback()

	assignmentID := uuid.NewString()
	_, err = tx.Exec(`
		INSERT INTO assignments (
			id, request_id, employee_id, admin_id, 
			admin_notes, status, assigned_at
		) VALUES ($1, $2, $3, $4, $5, 'assigned', now())`,
		assignmentID, req.RequestID, req.EmployeeID, adminID, req.AdminNotes,
	)

	if err != nil {
		log.Printf("failed to create assignment: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.T(lang, "err_operation_failed")})
		return
	}

	_, err = tx.Exec(`
		UPDATE service_requests 
		SET status = 'assigned', updated_at = now() 
		WHERE id = $1`,
		req.RequestID,
	)

	if err != nil {
		log.Printf("failed to update request status: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.T(lang, "err_operation_failed")})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"message":       "تم تعيين الموظف بنجاح",
		"assignment_id": assignmentID,
	})
}

func (h *ServiceHandler) GetEmployees(c *gin.Context) {
	lang := i18n.Detect(c)

	rows, err := h.DB.Query(`
		SELECT id, full_name, email, phone, is_active
		FROM users 
		WHERE role = 'employee' AND is_active = TRUE
		ORDER BY full_name`)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.T(lang, "err_operation_failed")})
		return
	}
	defer rows.Close()

	var employees []gin.H
	for rows.Next() {
		var id, fullName, email, phone string
		var isActive bool
		if err := rows.Scan(&id, &fullName, &email, &phone, &isActive); err == nil {
			employees = append(employees, gin.H{
				"id":        id,
				"full_name": fullName,
				"email":     email,
				"phone":     phone,
				"is_active": isActive,
			})
		}
	}

	c.JSON(http.StatusOK, employees)
}

func (h *ServiceHandler) UpdateStatus(c *gin.Context) {
	lang := i18n.Detect(c)
	id := c.Param("id")

	var req struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.T(lang, "err_missing_fields")})
		return
	}

	_, err := h.DB.Exec(`
		UPDATE service_requests 
		SET status = $1, updated_at = now() 
		WHERE id = $2`,
		req.Status, id,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.T(lang, "err_operation_failed")})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "تم تحديث حالة الطلب"})
}
