package handlers

import (
	"admin_server/backend/internal/models"
	"admin_server/backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TestHandler struct {
	service *services.TestService
}

func NewTestHandler(service *services.TestService) *TestHandler {
	return &TestHandler{
		service: service,
	}
}

// TriggerTest handles POST /api/v1/tests/trigger
func (h *TestHandler) TriggerTest(c *gin.Context) {
	var req models.TriggerTestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.TestType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "test_type is required"})
		return
	}

	response, err := h.service.TriggerTest(req.TestType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, response)
}

