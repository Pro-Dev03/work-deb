package handlers

import (
	"net/http"

	"worktrack/backend/internal/services"

	"github.com/gin-gonic/gin"
)

type GeocodingHandler struct {
	service *services.GeocodingService
}

func NewGeocodingHandler(service *services.GeocodingService) *GeocodingHandler {
	return &GeocodingHandler{service: service}
}

func (h *GeocodingHandler) Autocomplete(c *gin.Context) {
	query := c.Query("q")
	lang := c.Query("lang")
	
	if lang == "" {
		lang = "ar"
	}
	
	if len(query) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "الرجاء إدخال حرفين على الأقل"})
		return
	}

	results, err := h.service.Autocomplete(query, lang)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"results": results,
		"count":   len(results),
	})
}
