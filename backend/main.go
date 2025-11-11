package main

import (
	"log"
	"os"

	"admin_server/backend/internal/config"
	"admin_server/backend/internal/handlers"
	"admin_server/backend/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize services
	ruleService := services.NewRuleService(cfg)
	syscallService := services.NewSyscallService(cfg)
	alertService := services.NewAlertService(cfg)
	testService := services.NewTestService(cfg)

	// Initialize handlers
	ruleHandler := handlers.NewRuleHandler(ruleService)
	syscallHandler := handlers.NewSyscallHandler(syscallService)
	alertHandler := handlers.NewAlertHandler(alertService)
	testHandler := handlers.NewTestHandler(testService)

	// Setup router
	router := gin.Default()

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// API routes
	api := router.Group("/api/v1")
	{
		// Rules endpoints
		api.GET("/rules", ruleHandler.GetRules)
		api.PUT("/rules", ruleHandler.UpdateRules)

		// Syscalls endpoints
		api.GET("/syscalls/callable", syscallHandler.GetCallableSyscalls)

		// Alerts endpoints
		api.GET("/alerts", alertHandler.GetAlerts)
		api.POST("/alerts/webhook", alertHandler.ReceiveWebhook)

		// Test endpoints
		api.POST("/tests/trigger", testHandler.TriggerTest)
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

