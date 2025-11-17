package services

import (
	"admin_server/backend/internal/config"
	"admin_server/backend/internal/models"
	"fmt"
	"log"

	"context" // <-- [추가]

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"gopkg.in/yaml.v3"
	"k8s.io/client-go/kubernetes"
)

// RuleService handles rule-related operations
type RuleService struct {
	cfg *config.Config
	// TODO: Add K8s client when implementing actual K8s integration
	clientset kubernetes.Interface
}

func NewRuleService(cfg *config.Config, clientset kubernetes.Interface) *RuleService {
	return &RuleService{
		cfg:       cfg,
		clientset: clientset,
	}
}

// GetRules retrieves the current rules directly from ConfigMap via K8s API
func (s *RuleService) GetRules() (*models.RuleSet, error) {
	// [수정] K8s API를 통해 ConfigMap 데이터를 직접 조회하여 즉각적인 반영을 보장합니다.
	log.Println("Getting rules from Kubernetes ConfigMap via API")

	// 1. K8s API를 통해 ConfigMap의 현재 상태를 가져오기 (파일 읽기 로직 대체)
	configMap, err := s.clientset.CoreV1().ConfigMaps(s.cfg.Namespace).Get(context.TODO(), s.cfg.ConfigMapName, metav1.GetOptions{})
	if err != nil {
		log.Printf("Failed to get ConfigMap %s via API: %v", s.cfg.ConfigMapName, err)
		return nil, fmt.Errorf("failed to get ConfigMap via K8s API: %w", err)
	}

	// 2. ConfigMap의 Data 필드에서 'rule.yaml' 키의 값 추출
	yamlContent, ok := configMap.Data["rule.yaml"]
	if !ok {
		return nil, fmt.Errorf("ConfigMap %s does not contain 'rule.yaml' key", s.cfg.ConfigMapName)
	}

	// 3. YAML 내용을 모델로 언마샬
	var ruleSet models.RuleSet
	err = yaml.Unmarshal([]byte(yamlContent), &ruleSet)
	if err != nil {
		log.Printf("Failed to unmarshal rule YAML from ConfigMap: %v", err)
		return nil, fmt.Errorf("failed to unmarshal rule YAML: %w", err)
	}

	return &ruleSet, nil
}

// UpdateRules updates the rules in ConfigMap
func (s *RuleService) UpdateRules(ruleSet *models.RuleSet) (*models.UpdateRulesResponse, error) {

	// 1. Convert to YAML
	yamlData, err := yaml.Marshal(ruleSet)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal rules to YAML: %w", err)
	}

	log.Printf("Updating ConfigMap '%s' in namespace '%s'", s.cfg.ConfigMapName, s.cfg.Namespace)

	// 2. ConfigMap의 현재 상태를 K8s API에서 가져오기
	// s.clientset을 사용하여 RuleService에 주입된 클라이언트에 접근합니다.
	configMap, err := s.clientset.CoreV1().ConfigMaps(s.cfg.Namespace).Get(context.TODO(), s.cfg.ConfigMapName, metav1.GetOptions{})

	// 3. ConfigMap Get 실패 시 처리
	if err != nil {
		return nil, fmt.Errorf("failed to get ConfigMap %s: %w", s.cfg.ConfigMapName, err)
	}

	// 4. 데이터 업데이트
	configMap.Data["rule.yaml"] = string(yamlData)

	// 5. K8s API로 ConfigMap 업데이트
	_, err = s.clientset.CoreV1().ConfigMaps(s.cfg.Namespace).Update(context.TODO(), configMap, metav1.UpdateOptions{})
	if err != nil {
		// API 호출 실패 시 로그를 남김 (이 로그가 콘솔에 찍히는지 확인해야 함)
		log.Printf("ERROR: Failed to update ConfigMap via K8s API: %v", err)
		return nil, fmt.Errorf("failed to update ConfigMap via K8s API: %w", err)
	}

	// TODO: Trigger rule engine and eBPF generator to reload rules
	// ... (이후 룰 엔진 리로드 로직)

	newVersion := ruleSet.RulesetVersion

	return &models.UpdateRulesResponse{
		Status:     "success",
		Message:    "Rule.yaml ConfigMap updated successfully.",
		NewVersion: newVersion,
	}, nil
}

// ValidateRules validates the rule structure
func (s *RuleService) ValidateRules(ruleSet *models.RuleSet) error {
	if ruleSet.RulesetVersion == "" {
		return fmt.Errorf("ruleset_version is required")
	}
	if len(ruleSet.Rules) == 0 {
		return fmt.Errorf("at least one rule is required")
	}
	for _, rule := range ruleSet.Rules {
		if rule.RuleID == "" {
			return fmt.Errorf("rule_id is required for all rules")
		}
		if len(rule.Conditions) == 0 {
			return fmt.Errorf("at least one condition is required for rule %s", rule.RuleID)
		}
	}
	return nil
}
