package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ClientHandler struct {
	DB *sql.DB
}

func NewClientHandler(db *sql.DB) *ClientHandler {
	return &ClientHandler{DB: db}
}

func (h *ClientHandler) Create(c *gin.Context) {
	var req struct {
		Name  string `json:"name" binding:"required"`
		Phone string `json:"phone"`
		Email string `json:"email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "بيانات غير صحيحة"})
		return
	}
	id := uuid.NewString()
	_, err := h.DB.Exec(`
		INSERT INTO clients (id, name, phone, email, created_at)
		VALUES ($1, $2, $3, $4, now())`,
		id, req.Name, req.Phone, req.Email,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل الإضافة"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *ClientHandler) List(c *gin.Context) {
	rows, err := h.DB.Query(`SELECT id, name, phone, email FROM clients`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل الجلب"})
		return
	}
	defer rows.Close()
	var clients []gin.H
	for rows.Next() {
		var id, name, phone, email string
		if err := rows.Scan(&id, &name, &phone, &email); err == nil {
			clients = append(clients, gin.H{"id": id, "name": name, "phone": phone, "email": email})
		}
	}
	c.JSON(http.StatusOK, clients)
}
