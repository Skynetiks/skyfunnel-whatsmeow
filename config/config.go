package config

import (
	"os"
	"strconv"
)

type Config struct {
	DatabaseURL     string
	Port           int
	LogLevel       string
	SessionDir     string
	WhatsMeowLogLevel string
	RedisURL       string
	EnableMetrics  bool
	MetricsPort    int
}

func Load() *Config {
	return &Config{
		DatabaseURL:     getEnv("DATABASE_URL", "postgres://localhost:5432/skyfunnel"),
		Port:           getEnvAsInt("PORT", 8081),
		LogLevel:       getEnv("LOG_LEVEL", "info"),
		SessionDir:     getEnv("WHATSMEOW_SESSION_DIR", "./sessions"),
		WhatsMeowLogLevel: getEnv("WHATSMEOW_LOG_LEVEL", "info"),
		RedisURL:       getEnv("REDIS_URL", ""),
		EnableMetrics:  getEnvAsBool("ENABLE_METRICS", false),
		MetricsPort:    getEnvAsInt("METRICS_PORT", 9090),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
