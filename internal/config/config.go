package config

import (
	"os"
	"time"
)

type AppConfig struct {
	Okta   *OktaConfig
	Server *ServerConfig
}

type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type OktaConfig struct {
	Domain   string
	APIToken string
	Issuer   string
	Audience string
}

func Load() *AppConfig {
	cfg := &AppConfig{
		Server: &ServerConfig{
			Port:         getEnvOrDefault("SERVER_PORT", "8080"),
			ReadTimeout:  getDurationOrDefault("SERVER_READ_TIMEOUT", "15s"),
			IdleTimeout:  getDurationOrDefault("SERVER_IDLE_TIMEOUT", "60s"),
			WriteTimeout: getDurationOrDefault("SERVER_WRITE_TIMEOUT", "15s"),
		},
		Okta: &OktaConfig{
			Domain:   os.Getenv("OKTA_DOMAIN"),
			Issuer:   os.Getenv("OKTA_ISSUER"),
			Audience: os.Getenv("OKTA_AUDIENCE"),
			APIToken: os.Getenv("OKTA_API_TOKEN"),
		},
	}
	return cfg
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getDurationOrDefault(key, defaultValue string) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	duration, _ := time.ParseDuration(defaultValue)
	return duration
}
