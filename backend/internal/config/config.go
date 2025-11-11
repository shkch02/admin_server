package config

import (
	"os"
)

type Config struct {
	// K8s configuration
	KubeConfigPath string
	Namespace      string
	ConfigMapName  string

	// Redis configuration
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int

	// Cluster syscalls Redis
	ClusterSyscallsRedisHost     string
	ClusterSyscallsRedisPort     string
	ClusterSyscallsRedisPassword string
	ClusterSyscallsRedisDB       int
}

func Load() *Config {
	return &Config{
		KubeConfigPath: getEnv("KUBE_CONFIG_PATH", ""),
		Namespace:      getEnv("NAMESPACE", "default"),
		ConfigMapName:  getEnv("CONFIG_MAP_NAME", "rule-yaml"),

		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       0,

		ClusterSyscallsRedisHost:     getEnv("CLUSTER_SYSCALLS_REDIS_HOST", "localhost"),
		ClusterSyscallsRedisPort:     getEnv("CLUSTER_SYSCALLS_REDIS_PORT", "6379"),
		ClusterSyscallsRedisPassword: getEnv("CLUSTER_SYSCALLS_REDIS_PASSWORD", ""),
		ClusterSyscallsRedisDB:       1,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

