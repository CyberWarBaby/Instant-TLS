package config

import (
	"log"
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
	env := getEnv("ENV", "development")

	// Render uses PORT
	port := os.Getenv("PORT")
	if port == "" {
		port = getEnv("API_PORT", "8080")
	}

	host := getEnv("API_HOST", "0.0.0.0")

	dbURL := os.Getenv("DATABASE_URL")
	jwt := os.Getenv("JWT_SECRET")
	cors := os.Getenv("CORS_ORIGINS")

	// ✅ In production, these MUST exist (deploy reads them from Render env vars)
	if env == "production" {
		if dbURL == "" {
			log.Fatal("DATABASE_URL is required in production")
		}
		if jwt == "" {
			log.Fatal("JWT_SECRET is required in production")
		}
		if cors == "" {
			log.Fatal("CORS_ORIGINS is required in production")
		}
	}

	// ✅ In development, allow local defaults
	if env != "production" {
		if dbURL == "" {
			dbURL = "postgres://instanttls:instanttls@localhost:5433/instanttls?sslmode=disable"
		}
		if jwt == "" {
			jwt = "dev-secret"
		}
		if cors == "" {
			cors = "http://localhost:3000"
		}
	}

	// Parse origins (trim spaces!)
	rawOrigins := strings.Split(cors, ",")
	origins := make([]string, 0, len(rawOrigins))
	for _, o := range rawOrigins {
		o = strings.TrimSpace(o)
		if o != "" {
			origins = append(origins, o)
		}
	}

	return &Config{
		DatabaseURL: dbURL,
		Port:        port,
		Host:        host,
		JWTSecret:   jwt,
		CORSOrigins: origins,
		Env:         env,
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
