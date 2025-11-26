package services

import (
	"admin_server/backend/internal/config"
	"admin_server/backend/internal/models"
	"log"
	"sort"
	"time"
)

// AlertService handles alert-related operations
type AlertService struct {
	cfg *config.Config
	// In-memory storage for alerts (TODO: Replace with Redis)
	alerts []models.Alert
	// TODO: Add Redis client when implementing actual Redis integration
	// redisClient *redis.Client
}

func NewAlertService(cfg *config.Config) *AlertService {
	return &AlertService{
		cfg:    cfg,
		alerts: make([]models.Alert, 0),
	}
}

// GetAlerts retrieves alerts with optional filtering
func (s *AlertService) GetAlerts(limit int, since *time.Time) (*models.AlertsResponse, error) {
	// TODO: Implement actual Redis retrieval
	log.Println("Getting alerts (mock implementation)")

	// Filter alerts
	filteredAlerts := s.alerts

	// Filter by time if since is provided
	if since != nil {
		filtered := make([]models.Alert, 0)
		for _, alert := range filteredAlerts {
			alertTime, err := time.Parse(time.RFC3339, alert.Timestamp)
			if err == nil && alertTime.After(*since) {
				filtered = append(filtered, alert)
			}
		}
		filteredAlerts = filtered
	}

	// Sort by timestamp (newest first)
	sort.Slice(filteredAlerts, func(i, j int) bool {
		timeI, _ := time.Parse(time.RFC3339, filteredAlerts[i].Timestamp)
		timeJ, _ := time.Parse(time.RFC3339, filteredAlerts[j].Timestamp)
		return timeI.After(timeJ)
	})

	// Apply limit
	if limit > 0 && limit < len(filteredAlerts) {
		filteredAlerts = filteredAlerts[:limit]
	}

	// Actual implementation will look like:
	// redisClient := redis.NewClient(&redis.Options{
	//     Addr:     fmt.Sprintf("%s:%s", s.cfg.RedisHost, s.cfg.RedisPort),
	//     Password: s.cfg.RedisPassword,
	//     DB:       s.cfg.RedisDB,
	// })
	//
	// // Get alerts from Redis sorted set or list
	// alerts, err := redisClient.ZRangeByScore(context.TODO(), "alerts", &redis.ZRangeBy{
	//     Min: sinceTimestamp,
	//     Max: "+inf",
	// }).Result()

	return &models.AlertsResponse{
		Alerts: filteredAlerts,
	}, nil
}

// ReceiveWebhook receives an alert from the rule engine via webhook
func (s *AlertService) ReceiveWebhook(alert *models.WebhookAlert) error {
	// TODO: This will be called by Kafka consumer
	// For now, store in memory
	log.Printf("Receiving webhook alert: %s", alert.AlertID)

	// Convert WebhookAlert to Alert
	newAlert := models.Alert{
		AlertID:         alert.AlertID,
		Timestamp:       alert.Timestamp,
		RuleID:          alert.RuleID,
		RuleDescription: alert.RuleDescription,
		Severity:        alert.Severity,
		PodName:         alert.PodName,
		Namespace:       alert.Namespace,
		SyscallLog:      alert.SyscallLog,
	}

	s.alerts = append(s.alerts, newAlert)

	// TODO: Store in Redis
	// redisClient := redis.NewClient(&redis.Options{
	//     Addr:     fmt.Sprintf("%s:%s", s.cfg.RedisHost, s.cfg.RedisPort),
	//     Password: s.cfg.RedisPassword,
	//     DB:       s.cfg.RedisDB,
	// })
	//
	// alertJSON, _ := json.Marshal(newAlert)
	// err := redisClient.LPush(context.TODO(), "alerts", alertJSON).Err()
	// or use sorted set with timestamp as score

	// TODO: Send webhook notification (Slack, etc.)
	// sendWebhookNotification(newAlert)

	return nil
}
