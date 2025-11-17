package services

import (
	"admin_server/backend/internal/config"
	"admin_server/backend/internal/models"
	"fmt"
	"log"
	"os"

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

// GetRules retrieves the current rules from ConfigMap
func (s *RuleService) GetRules() (*models.RuleSet, error) {
	// TODO: Implement actual K8s ConfigMap retrieval
	// For now, return mock data
	log.Println("Getting rules from mounted ConfigMap file")

	// Mock data matching the example in requirements

	yamlData, err := os.ReadFile(s.cfg.RuleYamlPath)
	if err != nil {
		log.Printf("Failed to read rule file (%s): %v", s.cfg.RuleYamlPath, err)
		return nil, fmt.Errorf("failed to read rule file: %w", err)
	}

	var ruleSet models.RuleSet
	err = yaml.Unmarshal(yamlData, &ruleSet)
	if err != nil {
		log.Printf("Failed to unmarshal rule YAML: %v", err)
		return nil, fmt.Errorf("failed to unmarshal rule YAML: %w", err)
	}

	return &ruleSet, nil

	// Actual implementation will look like:
	// clientset, err := kubernetes.NewForConfig(k8sConfig)
	// configMap, err := clientset.CoreV1().ConfigMaps(s.cfg.Namespace).Get(context.TODO(), s.cfg.ConfigMapName, metav1.GetOptions{})
	// yamlContent := configMap.Data["rule.yaml"]
	// var ruleSet models.RuleSet
	// err = yaml.Unmarshal([]byte(yamlContent), &ruleSet)
	// return &ruleSet, nil
}

// UpdateRules updates the rules in ConfigMap
// UpdateRules updates the rules in ConfigMap
func (s *RuleService) UpdateRules(ruleSet *models.RuleSet) (*models.UpdateRulesResponse, error) {

	// Convert to YAML
	yamlData, err := yaml.Marshal(ruleSet)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal rules to YAML: %w", err)
	}

	log.Printf("YAML to be written:\n%s", string(yamlData))

	// 1. K8s API로 ConfigMap 가져오기
	configMap, err := s.clientset.CoreV1().ConfigMaps(s.cfg.Namespace).Get(context.TODO(), s.cfg.ConfigMapName, metav1.GetOptions{})

	// 2. [필수] Get 직후 오류 처리
	if err != nil {
		return nil, fmt.Errorf("failed to get ConfigMap %s: %w", s.cfg.ConfigMapName, err)
	}

	// 3. 데이터 업데이트 (중복 코드 제거)
	configMap.Data["rule.yaml"] = string(yamlData)

	// 4. K8s API로 ConfigMap 업데이트
	_, err = s.clientset.CoreV1().ConfigMaps(s.cfg.Namespace).Update(context.TODO(), configMap, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to update ConfigMap: %w", err)
	}

	// TODO: Trigger rule engine and eBPF generator to reload rules
	// ... (이후 작업)

	// Increment version
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
