package services

import (
	"admin_server/backend/internal/config"
	"admin_server/backend/internal/models"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
)

// RuleService handles rule-related operations
type RuleService struct {
	cfg *config.Config
	// TODO: Add K8s client when implementing actual K8s integration
	// clientset kubernetes.Interface
}

func NewRuleService(cfg *config.Config) *RuleService {
	return &RuleService{
		cfg: cfg,
	}
}

// GetRules retrieves the current rules from ConfigMap
func (s *RuleService) GetRules() (*models.RuleSet, error) {
	// TODO: Implement actual K8s ConfigMap retrieval
	// For now, return mock data
	log.Println("Getting rules from ConfigMap (mock implementation)")

	// Mock data matching the example in requirements
	mockRules := &models.RuleSet{
		RulesetVersion: "1.0.0",
		Description:    "eBPF System Call based Security Violation Detection Rules",
		Rules: []models.Rule{
			{
				RuleID:      "RULE_A01_HOST_CRITICAL_WRITE",
				Description: "Detection of write access to critical host files (/etc/passwd, /etc/shadow, /etc/hosts)",
				Conditions: []models.Condition{
					{Field: "syscall_name", Operator: "equals", Value: "openat"},
					{Field: "flags", Operator: "contains_any", Value: []string{"O_WRONLY", "O_RDWR"}},
					{Field: "file_path", Operator: "starts_with_any", Value: []string{"/etc/passwd", "/etc/shadow", "/etc/hosts"}},
				},
			},
			{
				RuleID:      "RULE_B02_HOST_AUTH_READ",
				Description: "Detection of read access to host authentication files (.ssh/id_rsa, .kube/config)",
				Conditions: []models.Condition{
					{Field: "syscall_name", Operator: "equals", Value: "openat"},
					{Field: "flags", Operator: "not_contains_any", Value: []string{"O_WRONLY", "O_RDWR"}},
					{Field: "file_path", Operator: "ends_with_any", Value: []string{"/id_rsa", "/known_hosts", "/.ssh/config", "/.kube/config"}},
				},
			},
			{
				RuleID:      "RULE_C03_CONTAINER_ESCAPE_PATH",
				Description: "Detection of access attempts to common container escape paths (/proc/sys, /sys/kernel)",
				Conditions: []models.Condition{
					{Field: "syscall_name", Operator: "equals", Value: "openat"},
					{Field: "file_path", Operator: "starts_with_any", Value: []string{"/proc/sys/kernel/", "/sys/kernel/"}},
				},
			},
		},
	}

	return mockRules, nil

	// Actual implementation will look like:
	// clientset, err := kubernetes.NewForConfig(k8sConfig)
	// configMap, err := clientset.CoreV1().ConfigMaps(s.cfg.Namespace).Get(context.TODO(), s.cfg.ConfigMapName, metav1.GetOptions{})
	// yamlContent := configMap.Data["rule.yaml"]
	// var ruleSet models.RuleSet
	// err = yaml.Unmarshal([]byte(yamlContent), &ruleSet)
	// return &ruleSet, nil
}

// UpdateRules updates the rules in ConfigMap
func (s *RuleService) UpdateRules(ruleSet *models.RuleSet) (*models.UpdateRulesResponse, error) {
	// TODO: Implement actual K8s ConfigMap update
	log.Println("Updating rules in ConfigMap (mock implementation)")

	// Convert to YAML
	yamlData, err := yaml.Marshal(ruleSet)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal rules to YAML: %w", err)
	}

	log.Printf("YAML to be written:\n%s", string(yamlData))

	// TODO: Update ConfigMap via K8s API
	// clientset, err := kubernetes.NewForConfig(k8sConfig)
	// configMap, err := clientset.CoreV1().ConfigMaps(s.cfg.Namespace).Get(context.TODO(), s.cfg.ConfigMapName, metav1.GetOptions{})
	// configMap.Data["rule.yaml"] = string(yamlData)
	// _, err = clientset.CoreV1().ConfigMaps(s.cfg.Namespace).Update(context.TODO(), configMap, metav1.UpdateOptions{})

	// TODO: Trigger rule engine and eBPF generator to reload rules
	// This can be done via:
	// 1. Sending a signal to rule engine pod
	// 2. Updating an annotation on the rule engine deployment
	// 3. Using a webhook/event system

	// Increment version
	newVersion := "1.0.1" // TODO: Calculate actual new version

	return &models.UpdateRulesResponse{
		Status:    "success",
		Message:   "Rule.yaml ConfigMap updated successfully.",
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

