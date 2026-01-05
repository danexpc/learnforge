package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Port                 string `yaml:"port"`
	Storage              string `yaml:"storage"`
	DatabaseURL          string `yaml:"database_url"`
	AIProvider           string `yaml:"ai_provider"` // "openai" or "gemini"
	AIBaseURL            string `yaml:"ai_base_url"`
	AIApiKey             string `yaml:"ai_api_key"`
	AIModel              string `yaml:"ai_model"`
	LogLevel             string `yaml:"log_level"`
	SlackWebhookURL      string `yaml:"slack_webhook_url"`
	SlackErrorWebhookURL string `yaml:"slack_error_webhook_url"`
	SummaryAPIKey        string `yaml:"summary_api_key"`
	RedisURL             string `yaml:"redis_url"`
}

func Load() (*Config, error) {
	env := getEnv("ENV", "development")
	configPath := getEnv("CONFIG_PATH", "")

	var cfg Config

	if configPath == "" {
		configPath = filepath.Join("config", env+".yml")
	}

	if _, err := os.Stat(configPath); err == nil {
		data, err := os.ReadFile(configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}

		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return nil, fmt.Errorf("failed to parse config file: %w", err)
		}
	}

	if cfg.Port == "" {
		cfg.Port = getEnv("PORT", "8080")
	}
	if cfg.Storage == "" {
		cfg.Storage = getEnv("STORAGE", "inmem")
	}
	if cfg.DatabaseURL == "" {
		cfg.DatabaseURL = getEnv("DATABASE_URL", "")
	}
	if cfg.AIProvider == "" {
		cfg.AIProvider = getEnv("AI_PROVIDER", "openai")
	}
	if cfg.AIBaseURL == "" {
		if cfg.AIProvider == "gemini" {
			cfg.AIBaseURL = "https://generativelanguage.googleapis.com"
		} else {
			cfg.AIBaseURL = getEnv("AI_BASE_URL", "https://api.openai.com")
		}
	}
	if cfg.AIApiKey == "" {
		cfg.AIApiKey = getEnv("AI_API_KEY", "")
	}
	if cfg.AIModel == "" {
		if cfg.AIProvider == "gemini" {
			cfg.AIModel = getEnv("AI_MODEL", "gemini-2.0-flash-exp")
		} else {
			cfg.AIModel = getEnv("AI_MODEL", "gpt-3.5-turbo")
		}
	}
	if cfg.LogLevel == "" {
		cfg.LogLevel = getEnv("LOG_LEVEL", "info")
	}
	if cfg.SlackWebhookURL == "" {
		cfg.SlackWebhookURL = getEnv("SLACK_WEBHOOK_URL", "")
	}
	if cfg.SlackErrorWebhookURL == "" {
		cfg.SlackErrorWebhookURL = getEnv("SLACK_ERROR_WEBHOOK_URL", "")
	}
	if cfg.SummaryAPIKey == "" {
		cfg.SummaryAPIKey = getEnv("SUMMARY_API_KEY", "")
	}
	if cfg.RedisURL == "" {
		cfg.RedisURL = getEnv("REDIS_URL", "")
	}

	return &cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
