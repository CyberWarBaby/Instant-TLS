package config

import (
	"os"
	"strings"
)

type Config struct {
	DatabaseURL string
	Port        string
	Host        string
	JWTSecret   string
	CORSOrigins []string
	Env         string
}

func Load() *Config {
	corsOrigins := os.Getenv("CORS_ORIGINS")
	if corsOrigins == "" {
		corsOrigins = "http://localhost:3000"
	}

	return &Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgres://instanttls:instanttls@localhost:5432/instanttls?sslmode=disable"),
		Port:        getEnv("API_PORT", "8081"),
		Host:        getEnv("API_HOST", "0.0.0.0"),
		JWTSecret:   getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-in-production"),
		CORSOrigins: strings.Split(corsOrigins, ","),
		Env:         getEnv("ENV", "development"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
