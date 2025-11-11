package handlers

import (
	"admin_server/backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SyscallHandler struct {
	service *services.SyscallService
}

func NewSyscallHandler(service *services.SyscallService) *SyscallHandler {
	return &SyscallHandler{
		service: service,
	}
}

// GetCallableSyscalls handles GET /api/v1/syscalls/callable
func (h *SyscallHandler) GetCallableSyscalls(c *gin.Context) {
	response, err := h.service.GetCallableSyscalls()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

