package handlers

import (
	"admin_server/backend/internal/models"
	"admin_server/backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RuleHandler struct {
	service *services.RuleService
}

func NewRuleHandler(service *services.RuleService) *RuleHandler {
	return &RuleHandler{
		service: service,
	}
}

// GetRules handles GET /api/v1/rules
func (h *RuleHandler) GetRules(c *gin.Context) {
	rules, err := h.service.GetRules()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rules)
}

// UpdateRules handles PUT /api/v1/rules
func (h *RuleHandler) UpdateRules(c *gin.Context) {
	var ruleSet models.RuleSet
	if err := c.ShouldBindJSON(&ruleSet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate rules
	if err := h.service.ValidateRules(&ruleSet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.UpdateRules(&ruleSet)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

