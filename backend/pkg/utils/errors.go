package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorResponse هيكل استجابة الخطأ الموحد
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code,omitempty"`
}

// RespondWithError إرسال استجابة خطأ موحدة
func RespondWithError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, ErrorResponse{
		Error:   http.StatusText(statusCode),
		Message: message,
		Code:    statusCode,
	})
}

// RespondWithSuccess إرسال استجابة نجاح موحدة
func RespondWithSuccess(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, gin.H{
		"success": true,
		"data":    data,
	})
}

// DatabaseErrorHandler معالج أخطاء قاعدة البيانات
func DatabaseErrorHandler(c *gin.Context, err error, operation string) {
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "database_error",
			"message": "فشل العملية على قاعدة البيانات",
			"operation": operation,
		})
		return
	}
}
