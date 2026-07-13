package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"worktrack/backend/internal/i18n"
	"worktrack/backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TaskHandler struct {
	DB *sql.DB
}

func NewTaskHandler(db *sql.DB) *TaskHandler {
	return &TaskHandler{DB: db}
}

// MyTasks يعيد مهام الموظف المسجل دخوله فقط (assigned_user_id = المستخدم الحالي)
func (h *TaskHandler) MyTasks(c *gin.Context) {
	lang := i18n.Detect(c)
	userID, _ := c.Get("user_id")

	rows, err := h.DB.Query(`
		SELECT id, title, description, worksite_id, status, scheduled_start, scheduled_end, created_at
		FROM tasks WHERE assigned_user_id = $1 ORDER BY scheduled_start ASC`, userID)
	if err != nil {
		log.Printf("failed to fetch tasks: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.T(lang, "err_operation_failed")})
		return
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var t models.Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.WorksiteID,
			&t.Status, &t.ScheduledStart, &t.ScheduledEnd, &t.CreatedAt); err == nil {
			tasks = append(tasks, t)
		}
	}

	c.JSON(http.StatusOK, tasks)
}

type createTaskRequest struct {
	Title          string  `json:"title" binding:"required"`
	Description    string  `json:"description"`
	ClientID       *string `json:"client_id"`
	WorksiteID     string  `json:"worksite_id" binding:"required"`
	AssignedUserID *string `json:"assigned_user_id"`
	ScheduledStart string  `json:"scheduled_start"`
	ScheduledEnd   string  `json:"scheduled_end"`
}

// Create ينشئ مهمة جديدة ويربطها بنقطة عمل (worksite) — هذا الربط هو ما يفعّل الـ Geofence لاحقاً
func (h *TaskHandler) Create(c *gin.Context) {
	lang := i18n.Detect(c)

	var req createTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.T(lang, "err_missing_fields")})
		return
	}

	id := uuid.NewString()
	_, err := h.DB.Exec(`
		INSERT INTO tasks (id, title, description, client_id, worksite_id, assigned_user_id,
			status, scheduled_start, scheduled_end, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, 'pending', NULLIF($7,'')::timestamptz, NULLIF($8,'')::timestamptz, now(), now())`,
		id, req.Title, req.Description, req.ClientID, req.WorksiteID, req.AssignedUserID,
		req.ScheduledStart, req.ScheduledEnd,
	)
	if err != nil {
		log.Printf("failed to create task: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.T(lang, "err_operation_failed")})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id, "message": i18n.T(lang, "msg_task_created")})
}

// ListAll يعيد كل المهام (لوحة تحكم المدير)
func (h *TaskHandler) ListAll(c *gin.Context) {
	lang := i18n.Detect(c)

	rows, err := h.DB.Query(`
		SELECT id, title, description, worksite_id, status, scheduled_start, scheduled_end, created_at
		FROM tasks ORDER BY created_at DESC LIMIT 200`)
	if err != nil {
		log.Printf("failed to fetch all tasks: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.T(lang, "err_operation_failed")})
		return
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var t models.Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.WorksiteID,
			&t.Status, &t.ScheduledStart, &t.ScheduledEnd, &t.CreatedAt); err == nil {
			tasks = append(tasks, t)
		}
	}

	c.JSON(http.StatusOK, tasks)
}
