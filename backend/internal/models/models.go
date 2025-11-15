package models

// RuleSet represents the entire ruleset configuration
type RuleSet struct {
	RulesetVersion string  `json:"ruleset_version" yaml:"ruleset_version"`
	Description    string  `json:"description" yaml:"description"`
	Rules          []Rule  `json:"rules" yaml:"rules"`
}

// Rule represents a single security rule
type Rule struct {
	RuleID      string      `json:"rule_id" yaml:"rule_id"`
	Description string      `json:"description" yaml:"description"`
	Conditions  []Condition `json:"conditions" yaml:"conditions"`
}

// Condition represents a condition in a rule
type Condition struct {
	Field    string      `json:"field" yaml:"field"`
	Operator string      `json:"operator" yaml:"operator"`
	Value    interface{} `json:"value" yaml:"value"`
}

// UpdateRulesResponse represents the response for updating rules
type UpdateRulesResponse struct {
	Status    string `json:"status"`
	Message   string `json:"message"`
	NewVersion string `json:"new_version,omitempty"`
}

// SyscallArg represents a syscall argument
type SyscallArg struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

// Syscall represents a system call with its arguments
type Syscall struct {
	Name string       `json:"name"`
	Args []SyscallArg `json:"args"`
}

// CallableSyscallsResponse represents the response for callable syscalls
type CallableSyscallsResponse struct {
	TotalCount int       `json:"total_count"`
	Syscalls   []Syscall `json:"syscalls"`
}

// Alert represents a security alert
type Alert struct {
	AlertID         string                 `json:"alert_id"`
	Timestamp       string                 `json:"timestamp"`
	RuleID          string                 `json:"rule_id"`
	RuleDescription string                 `json:"rule_description"`
	Severity        string                 `json:"severity"`
	PodName         string                 `json:"pod_name"`
	Namespace       string                 `json:"namespace"`
	SyscallLog      map[string]interface{} `json:"syscall_log"`
}

// AlertsResponse represents the response for alerts
type AlertsResponse struct {
	Alerts []Alert `json:"alerts"`
}

// TriggerTestRequest represents the request for triggering a test
type TriggerTestRequest struct {
	TestType string `json:"test_type"`
}

// TriggerTestResponse represents the response for triggering a test
type TriggerTestResponse struct {
	Status  string `json:"status"`
	JobName string `json:"job_name"`
}

// WebhookAlert represents an alert received via webhook
type WebhookAlert struct {
	AlertID         string                 `json:"alert_id"`
	Timestamp       string                 `json:"timestamp"`
	RuleID          string                 `json:"rule_id"`
	RuleDescription string                 `json:"rule_description"`
	Severity        string                 `json:"severity"`
	PodName         string                 `json:"pod_name"`
	Namespace       string                 `json:"namespace"`
	SyscallLog      map[string]interface{} `json:"syscall_log"`
}



