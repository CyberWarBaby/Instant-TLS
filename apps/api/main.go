package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/CyberWarBaby/Instant-TLS/apps/api/internal/config"
	"github.com/CyberWarBaby/Instant-TLS/apps/api/internal/database"
	"github.com/CyberWarBaby/Instant-TLS/apps/api/internal/handlers"
	"github.com/CyberWarBaby/Instant-TLS/apps/api/internal/middleware"
	"github.com/CyberWarBaby/Instant-TLS/apps/api/internal/migrations"
	"github.com/CyberWarBaby/Instant-TLS/apps/api/internal/seed"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	// Load .env file
	_ = godotenv.Load("../../.env")
	_ = godotenv.Load(".env")

	// Initialize logger
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	sugar := logger.Sugar()

	// Load config
	cfg := config.Load()

	// Handle CLI commands
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "migrate":
			if len(os.Args) > 2 {
				switch os.Args[2] {
				case "up":
					if err := migrations.Up(cfg.DatabaseURL); err != nil {
						log.Fatalf("Migration up failed: %v", err)
					}
					fmt.Println("✅ Migrations applied successfully")
					return
				case "down":
					if err := migrations.Down(cfg.DatabaseURL); err != nil {
						log.Fatalf("Migration down failed: %v", err)
					}
					fmt.Println("✅ Migrations rolled back successfully")
					return
				}
			}
			fmt.Println("Usage: go run . migrate [up|down]")
			return
		case "seed":
			db, err := database.Connect(cfg.DatabaseURL)
			if err != nil {
				log.Fatalf("Failed to connect to database: %v", err)
			}
			if err := seed.SeedDemoUser(db); err != nil {
				log.Fatalf("Seed failed: %v", err)
			}
			fmt.Println("✅ Demo user seeded successfully")
			return
		}
	}

	// Connect to database
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		sugar.Fatalf("Failed to connect to database: %v", err)
	}

	// Seed demo user on startup in development
	if cfg.Env == "development" {
		if err := seed.SeedDemoUser(db); err != nil {
			sugar.Warnf("Failed to seed demo user: %v", err)
		}
	}

	// Setup Gin
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// CORS - Use configured origins from environment
	allowedOrigins := cfg.CORSOrigins
	// Add localhost origins for development
	if cfg.Env == "development" {
		allowedOrigins = append(allowedOrigins, "http://localhost:3000", "https://localhost:3000", "http://127.0.0.1:3000")
	}
	
	corsConfig := cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}
	
	// In development, allow additional origins via function
	if cfg.Env == "development" {
		corsConfig.AllowOriginFunc = func(origin string) bool {
			// Allow any .local domain
			if strings.HasSuffix(origin, ".local") || strings.HasSuffix(origin, ".local:443") {
				return true
			}
			// Allow localhost on any port
			if strings.HasPrefix(origin, "http://localhost") || strings.HasPrefix(origin, "https://localhost") {
				return true
			}
			if strings.HasPrefix(origin, "http://127.0.0.1") || strings.HasPrefix(origin, "https://127.0.0.1") {
				return true
			}
			return false
		}
	}
	
	r.Use(cors.New(corsConfig))

	// Initialize handlers
	h := handlers.New(db, cfg, sugar)

	// Routes
	v1 := r.Group("/v1")
	{
		// Auth routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", h.Register)
			auth.POST("/login", h.Login)
		}

		// Protected routes (PAT auth)
		v1.GET("/me", middleware.PATAuth(db), h.Me)
		v1.GET("/license", middleware.PATAuth(db), h.License)

		// Machine routes (PAT auth)
		machines := v1.Group("/machines")
		machines.Use(middleware.PATAuth(db))
		{
			machines.POST("/ping", h.MachinePing)
		}

		// Token routes (session auth for web)
		tokens := v1.Group("/tokens")
		tokens.Use(middleware.SessionAuth(cfg))
		{
			tokens.GET("", h.ListTokens)
			tokens.POST("", h.CreateToken)
			tokens.DELETE("/:id", h.DeleteToken)
		}

		// User routes (session auth for web)
		user := v1.Group("/user")
		user.Use(middleware.SessionAuth(cfg))
		{
			user.GET("", h.GetUser)
			user.POST("/plan", h.UpdatePlan)
		}
	}

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Start server
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	sugar.Infof("🚀 API server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		sugar.Fatalf("Failed to start server: %v", err)
	}
}
