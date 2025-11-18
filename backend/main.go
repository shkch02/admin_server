package main

import (
	"log"
	"os"

	"admin_server/backend/internal/config"
	"admin_server/backend/internal/handlers"
	"admin_server/backend/internal/services"

	// 컨텍스트 추가
	"context" // 컨텍스트 추가

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// --- [추가 시작]: K8s 클라이언트셋 초기화 ---
	// KubeConfigPath가 비어 있으면 In-Cluster-Config를 사용
	k8sConfig, err := clientcmd.BuildConfigFromFlags("", cfg.KubeConfigPath)
	if err != nil {
		log.Fatalf("Failed to build kubernetes config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(k8sConfig)
	if err != nil {
		log.Fatalf("Failed to create kubernetes clientset: %v", err)
	}
	// --- [추가 끝] ---

	// Initialize CCSL Redis Client (추가)
	ccslRedisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.CCSLRedisAddr,
		Password: cfg.CCSLRedisPassword,
		DB:       0,
	})

	// Check Redis connection
	ctx = context.Background()
	_, err := ccslRedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to CCSL Redis at %s: %v", cfg.CCSLRedisAddr, err)
	}
	log.Println("Successfully connected to CCSL Redis")

	// Initialize services
	ruleService := services.NewRuleService(cfg, clientset)
	syscallService := services.NewSyscallService(cfg, ccslRedisClient)
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
