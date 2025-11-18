package config

import (
	"os"
)

type Config struct {
	// K8s configuration
	KubeConfigPath string
	Namespace      string
	ConfigMapName  string
	RuleYamlPath   string

	// CCSL Redis 설정 (추가)
	CCSLRedisAddr     string
	CCSLRedisPassword string
}

func Load() *Config {
	return &Config{
		KubeConfigPath: getEnv("KUBE_CONFIG_PATH", ""),
		Namespace:      getEnv("NAMESPACE", "default"),
		ConfigMapName:  getEnv("CONFIG_MAP_NAME", "rule-yaml"),
		RuleYamlPath:   getEnv("RULE_YAML_FILE_PATH", "/etc/config/rule.yaml"),

		CCSLRedisAddr:     getEnv("CCSL_REDIS_ADDR", "redis-ccsl-svc:6379"),
		CCSLRedisPassword: getEnv("CCSL_REDIS_PASSWORD", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
