package handlers

import (
	"admin_server/backend/internal/models"
	"admin_server/backend/internal/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type AlertHandler struct {
	service *services.AlertService
}

func NewAlertHandler(service *services.AlertService) *AlertHandler {
	return &AlertHandler{
		service: service,
	}
}

// GetAlerts handles GET /api/v1/alerts
func (h *AlertHandler) GetAlerts(c *gin.Context) {
	// Parse query parameters
	limit := 50 // default
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	var since *time.Time
	if sinceStr := c.Query("since"); sinceStr != "" {
		if parsedTime, err := time.Parse(time.RFC3339, sinceStr); err == nil {
			since = &parsedTime
		}
	}

	response, err := h.service.GetAlerts(limit, since)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// ReceiveWebhook handles POST /api/v1/alerts/webhook
func (h *AlertHandler) ReceiveWebhook(c *gin.Context) {
	var alert models.WebhookAlert
	if err := c.ShouldBindJSON(&alert); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.ReceiveWebhook(&alert); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "received"})
}

