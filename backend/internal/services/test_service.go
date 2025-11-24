package services

import (
	"admin_server/backend/internal/config"
	"admin_server/backend/internal/models"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// TestService handles test-related operations
type TestService struct {
	cfg        *config.Config
	httpClient *http.Client
}

func NewTestService(cfg *config.Config) *TestService {
	return &TestService{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// 프론트가 보낸 http 트리거 처리 함수, http 핸들러에서 호출됨, testType에 rule ID 담겨서 오니까 그거로 분리
func (s *TestService) TriggerTest(testType string) (*models.TriggerTestResponse, error) {
	log.Printf("Triggering test: %s", testType)

	baseURL := "http://attacker-service:80"
	var targetURL string

	switch testType {
	case "RULE_A01_HOST_CRITICAL_WRITE":
		targetURL = baseURL + "/attack/write"
	case "RULE_B02_HOST_AUTH_READ":
		targetURL = baseURL + "/attack/read"
	case "RULE_C03_CONTAINER_ESCAPE_PATH":
		targetURL = baseURL + "/attack/read"
	default:
		targetURL = baseURL + "/attack/read"
	}

	//targetURL에 요청 전송함
	resp, err := s.httpClient.Get(targetURL)
	if err != nil {
		log.Printf("ERROR: Failed to trigger test %s: %v", testType, err)
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	responsMsg := string(bodyBytes)
	log.Printf("Received response from attacker service: %s", responsMsg)

	if resp.StatusCode != http.StatusOK {
		log.Printf("ERROR: Attacker service returned non-OK status for test %s: %d", testType, resp.StatusCode)
		return nil, fmt.Errorf("attacker service returned status: %d", resp.StatusCode)
	}

	return &models.TriggerTestResponse{
		Status:  "test_triggered",
		JobName: fmt.Sprintf("http-trigger-%s", testType), // Job 이름 대신 트리거 ID 반환
	}, nil
}
